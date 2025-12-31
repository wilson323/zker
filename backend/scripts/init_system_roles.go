/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
)

// ==================== 系统预置角色定义 ====================

// SystemRoleDefinition 系统角色定义
type SystemRoleDefinition struct {
	RoleName        string              // 角色名称
	RoleCode        string              // 角色代码
	Description     string              // 角色描述
	IsSystem        bool                // 是否系统角色
	Permissions     []string            // 权限列表
	DataPermissions []DataPermissionDef // 数据权限配置
	FieldPermissions []FieldPermissionDef // 字段权限配置
}

// DataPermissionDef 数据权限定义
type DataPermissionDef struct {
	ResourceType string // 资源类型
	Scope        string // 权限范围: ALL, DEPARTMENT, OWN, CUSTOM, NONE
	CustomFilter string // 自定义过滤器（scope=CUSTOM时）
}

// FieldPermissionDef 字段权限定义
type FieldPermissionDef struct {
	ResourceType string                          // 资源类型
	Fields       map[string]FieldPermissionLevel // 字段名 -> 权限级别
}

// FieldPermissionLevel 字段权限级别定义
type FieldPermissionLevel string

const (
	FieldLevelHidden   FieldPermissionLevel = "hidden"
	FieldLevelReadonly FieldPermissionLevel = "readonly"
	FieldLevelEditable FieldPermissionLevel = "editable"
)

// ==================== 预置角色配置 ====================

// SystemRoles 系统预置角色列表
var SystemRoles = []SystemRoleDefinition{
	{
		RoleName:    "超级管理员",
		RoleCode:    "super_admin",
		Description: "拥有所有权限，可以管理系统的所有资源",
		IsSystem:    true,
		Permissions: []string{
			"*", // 通配符，表示所有权限
		},
		DataPermissions: []DataPermissionDef{
			{
				ResourceType: "*",
				Scope:        "ALL",
			},
		},
		FieldPermissions: []FieldPermissionDef{}, // 超级管理员无字段限制
	},
	{
		RoleName:    "管理员",
		RoleCode:    "admin",
		Description: "拥有租户内所有资源的完整管理权限",
		IsSystem:    true,
		Permissions: []string{
			"bots.*",
			"workflows.*",
			"conversations.*",
			"knowledge_bases.*",
			"users.*",
			"roles.*",
			"departments.*",
			"tenants.view",
		},
		DataPermissions: []DataPermissionDef{
			{
				ResourceType: "bots",
				Scope:        "ALL",
			},
			{
				ResourceType: "workflows",
				Scope:        "ALL",
			},
			{
				ResourceType: "conversations",
				Scope:        "ALL",
			},
			{
				ResourceType: "knowledge_bases",
				Scope:        "ALL",
			},
			{
				ResourceType: "users",
				Scope:        "ALL",
			},
		},
		FieldPermissions: []FieldPermissionDef{}, // 管理员无字段限制
	},
	{
		RoleName:    "开发者",
		RoleCode:    "developer",
		Description: "可以创建和管理Bot、工作流，查看自己的数据",
		IsSystem:    true,
		Permissions: []string{
			"bots.create",
			"bots.view",
			"bots.edit",
			"bots.delete",
			"workflows.create",
			"workflows.view",
			"workflows.edit",
			"workflows.delete",
			"conversations.view",
			"knowledge_bases.create",
			"knowledge_bases.view",
			"knowledge_bases.edit",
			"knowledge_bases.delete",
		},
		DataPermissions: []DataPermissionDef{
			{
				ResourceType: "bots",
				Scope:        "OWN",
			},
			{
				ResourceType: "workflows",
				Scope:        "OWN",
			},
			{
				ResourceType: "conversations",
				Scope:        "OWN",
			},
			{
				ResourceType: "knowledge_bases",
				Scope:        "OWN",
			},
		},
		FieldPermissions: []FieldPermissionDef{
			{
				ResourceType: "bots",
				Fields: map[string]FieldPermissionLevel{
					"name":              FieldLevelEditable,
					"description":       FieldLevelEditable,
					"config":            FieldLevelEditable,
					"api_key":           FieldLevelHidden,   // API密钥隐藏
					"webhook_secret":    FieldLevelHidden,   // Webhook密钥隐藏
					"created_at":        FieldLevelReadonly,
					"created_by":        FieldLevelReadonly,
				},
			},
		},
	},
	{
		RoleName:    "查看者",
		RoleCode:    "viewer",
		Description: "只能查看数据，不能修改或删除",
		IsSystem:    true,
		Permissions: []string{
			"bots.view",
			"workflows.view",
			"conversations.view",
			"knowledge_bases.view",
		},
		DataPermissions: []DataPermissionDef{
			{
				ResourceType: "bots",
				Scope:        "DEPARTMENT",
			},
			{
				ResourceType: "workflows",
				Scope:        "DEPARTMENT",
			},
			{
				ResourceType: "conversations",
				Scope:        "OWN",
			},
			{
				ResourceType: "knowledge_bases",
				Scope:        "DEPARTMENT",
			},
		},
		FieldPermissions: []FieldPermissionDef{
			{
				ResourceType: "bots",
				Fields: map[string]FieldPermissionLevel{
					"name":              FieldLevelReadonly,
					"description":       FieldLevelReadonly,
					"config":            FieldLevelReadonly,
					"api_key":           FieldLevelHidden,
					"webhook_secret":    FieldLevelHidden,
					"created_at":        FieldLevelReadonly,
					"created_by":        FieldLevelReadonly,
					"updated_at":        FieldLevelReadonly,
				},
			},
		},
	},
	{
		RoleName:    "数据分析员",
		RoleCode:    "analyst",
		Description: "可以查看所有数据并进行统计分析",
		IsSystem:    true,
		Permissions: []string{
			"bots.view",
			"workflows.view",
			"conversations.view",
			"knowledge_bases.view",
			"statistics.view",
			"reports.view",
			"exports.create",
		},
		DataPermissions: []DataPermissionDef{
			{
				ResourceType: "bots",
				Scope:        "ALL",
			},
			{
				ResourceType: "workflows",
				Scope:        "ALL",
			},
			{
				ResourceType: "conversations",
				Scope:        "ALL",
			},
		},
		FieldPermissions: []FieldPermissionDef{
			{
				ResourceType: "conversations",
				Fields: map[string]FieldPermissionLevel{
					"conversation_id": FieldLevelReadonly,
					"messages":        FieldLevelReadonly,
					"user_id":         FieldLevelReadonly,
					"created_at":      FieldLevelReadonly,
					// 隐藏敏感信息
					"user_ip":         FieldLevelHidden,
					"user_agent":      FieldLevelHidden,
				},
			},
		},
	},
}

// ==================== 主程序 ====================

func main() {
	// 解析命令行参数
	tenantID := flag.String("tenant", "", "租户ID (必填)")
	dsn := flag.String("dsn", "", "数据库连接字符串 (必填)")
	dryRun := flag.Bool("dry-run", false, "演练模式（不执行实际操作）")
	flag.Parse()

	// 验证参数
	if *tenantID == "" {
		log.Fatal("❌ 错误: 必须指定 -tenant 参数")
	}
	if *dsn == "" {
		log.Fatal("❌ 错误: 必须指定 -dsn 参数")
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	log.Printf("✅ 数据库连接成功")

	// 演练模式
	if *dryRun {
		log.Printf("🔍 [DRY RUN] 演练模式，将执行以下操作:")
		log.Printf("   租户ID: %s", *tenantID)
		log.Printf("   角色数量: %d", len(SystemRoles))
		for _, role := range SystemRoles {
			log.Printf("   - %s (%s)", role.RoleName, role.RoleCode)
		}
		log.Printf("✅ [DRY RUN] 演练完成")
		return
	}

	// 执行初始化
	if err := InitializeSystemRoles(db, *tenantID); err != nil {
		log.Fatalf("❌ 初始化系统角色失败: %v", err)
	}

	log.Printf("✅ 系统角色初始化完成！")
}

// InitializeSystemRoles 初始化系统角色
func InitializeSystemRoles(db *gorm.DB, tenantID string) error {
	log.Printf("🚀 开始初始化系统角色...")

	for _, roleDef := range SystemRoles {
		log.Printf("📝 正在创建角色: %s (%s)", roleDef.RoleName, roleDef.RoleCode)

		// 1. 创建角色
		role := &entity.Role{
			RoleID:      generateRoleID(tenantID, roleDef.RoleCode),
			TenantID:    tenantID,
			RoleName:    roleDef.RoleName,
			Description: roleDef.Description,
			IsSystem:    roleDef.IsSystem,
		}

		if err := db.Where("role_id = ?", role.RoleID).First(role).Error; err == gorm.ErrRecordNotFound {
			// 角色不存在，创建新角色
			if err := db.Create(role).Error; err != nil {
				return fmt.Errorf("创建角色 %s 失败: %w", roleDef.RoleName, err)
			}
			log.Printf("   ✅ 角色创建成功: %s", role.RoleName)
		} else if err != nil {
			return fmt.Errorf("查询角色 %s 失败: %w", roleDef.RoleName, err)
		} else {
			// 角色已存在，更新
			if err := db.Save(role).Error; err != nil {
				return fmt.Errorf("更新角色 %s 失败: %w", roleDef.RoleName, err)
			}
			log.Printf("   ♻️  角色已更新: %s", role.RoleName)
		}

		// 2. 创建数据权限
		for _, dpDef := range roleDef.DataPermissions {
			dataPerm := &entity.DataPermission{
				PermissionID: generatePermissionID(role.RoleID, dpDef.ResourceType),
				RoleID:       role.RoleID,
				ResourceType: entity.ResourceType(dpDef.ResourceType),
				Scope:        entity.DataPermissionScope(dpDef.Scope),
				CustomFilter: dpDef.CustomFilter,
			}

			if err := db.Where("permission_id = ?", dataPerm.PermissionID).First(dataPerm).Error; err == gorm.ErrRecordNotFound {
				if err := db.Create(dataPerm).Error; err != nil {
					return fmt.Errorf("创建数据权限失败: %w", err)
				}
			} else if err == nil {
				if err := db.Save(dataPerm).Error; err != nil {
					return fmt.Errorf("更新数据权限失败: %w", err)
				}
			}
		}

		// 3. 创建字段权限
		for _, fpDef := range roleDef.FieldPermissions {
			for fieldName, permLevel := range fpDef.Fields {
				fieldPerm := &entity.FieldPermission{
					PermissionID:    generateFieldPermID(role.RoleID, fpDef.ResourceType, fieldName),
					RoleID:          role.RoleID,
					ResourceType:    fpDef.ResourceType,
					FieldName:       fieldName,
					PermissionLevel: entity.FieldPermissionLevel(permLevel),
				}

				if err := db.Where("permission_id = ?", fieldPerm.PermissionID).First(fieldPerm).Error; err == gorm.ErrRecordNotFound {
					if err := db.Create(fieldPerm).Error; err != nil {
						return fmt.Errorf("创建字段权限失败: %w", err)
					}
				} else if err == nil {
					if err := db.Save(fieldPerm).Error; err != nil {
						return fmt.Errorf("更新字段权限失败: %w", err)
					}
				}
			}
		}
	}

	log.Printf("✅ 系统角色初始化完成！共创建/更新 %d 个角色", len(SystemRoles))
	return nil
}

// ==================== ID生成辅助函数 ====================

// generateRoleID 生成角色ID
func generateRoleID(tenantID, roleCode string) string {
	return fmt.Sprintf("%s:%s", tenantID, roleCode)
}

// generatePermissionID 生成数据权限ID
func generatePermissionID(roleID, resourceType string) string {
	return fmt.Sprintf("%s:dp:%s", roleID, resourceType)
}

// generateFieldPermID 生成字段权限ID
func generateFieldPermID(roleID, resourceType, fieldName string) string {
	return fmt.Sprintf("%s:fp:%s:%s", roleID, resourceType, fieldName)
}

// ==================== 使用说明 ====================

/*
使用方法：

1. 编译程序：
   cd backend/scripts
   go build -o init_system_roles init_system_roles.go

2. 执行初始化：
   ./init_system_roles -tenant="your_tenant_id" -dsn="user:password@tcp(localhost:3306)/database"

3. 演练模式（不实际执行）：
   ./init_system_roles -tenant="your_tenant_id" -dsn="..." -dry-run

4. 使用环境变量：
   export MYSQL_DSN="user:password@tcp(localhost:3306)/database"
   ./init_system_roles -tenant="your_tenant_id"

注意事项：
- 执行前请确保数据库已正确配置
- 系统角色创建后不能删除（只能修改）
- 建议在创建租户后立即执行此脚本
- 可以重复执行，会更新已存在的角色

系统角色列表：
1. super_admin - 超级管理员（所有权限）
2. admin - 管理员（租户内所有权限）
3. developer - 开发者（创建和管理Bot、工作流）
4. viewer - 查看者（只读权限）
5. analyst - 数据分析员（查看和导出数据）
*/
