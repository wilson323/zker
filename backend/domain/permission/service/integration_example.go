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

package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
)

// IntegrationExample 集成示例
type IntegrationExample struct {
	db                *gorm.DB
	dataPermChecker   DataPermissionChecker
	fieldPermChecker  FieldPermissionChecker
	roleRepo          repository.RoleRepository
	dataPermRepo      repository.DataPermissionRepository
	fieldPermRepo     repository.FieldPermissionRepository
	userRoleRepo      repository.UserRoleRepository
}

// NewIntegrationExample 创建集成示例
func NewIntegrationExample(
	db *gorm.DB,
	dataPermChecker DataPermissionChecker,
	fieldPermChecker FieldPermissionChecker,
	roleRepo repository.RoleRepository,
	dataPermRepo repository.DataPermissionRepository,
	fieldPermRepo repository.FieldPermissionRepository,
	userRoleRepo repository.UserRoleRepository,
) *IntegrationExample {
	return &IntegrationExample{
		db:               db,
		dataPermChecker:  dataPermChecker,
		fieldPermChecker: fieldPermChecker,
		roleRepo:         roleRepo,
		dataPermRepo:     dataPermRepo,
		fieldPermRepo:    fieldPermRepo,
		userRoleRepo:     userRoleRepo,
	}
}

// Example1_SimplePermissionCheck 示例1：简单的权限检查
func (e *IntegrationExample) Example1_SimplePermissionCheck(ctx context.Context, userID, botID string) (bool, error) {
	// 场景：检查用户是否有权限访问指定Bot

	// 方法1：使用CheckAccess
	hasAccess, err := e.dataPermChecker.CheckAccess(ctx, userID, botID, SELF, "bots")
	if err != nil {
		return false, fmt.Errorf("check access failed: %w", err)
	}

	if !hasAccess {
		return false, nil
	}

	// 方法2：生成过滤条件（用于列表查询）
	filter, err := e.dataPermChecker.Filter(ctx, userID, SELF, "bots")
	if err != nil {
		return false, fmt.Errorf("generate filter failed: %w", err)
	}

	// 应用过滤条件
	var count int64
	if err := e.db.Table("bots").Where(filter.WhereClause, filter.Args...).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Example2_FieldPermissionMask 示例2：字段权限脱敏
func (e *IntegrationExample) Example2_FieldPermissionMask(ctx context.Context, userID string) error {
	// 场景：获取Bot列表并应用字段权限脱敏

	// 1. 获取Bot数据（原始数据）
	var bots []map[string]interface{}
	if err := e.db.Table("bots").Limit(10).Find(&bots).Error; err != nil {
		return err
	}

	// 2. 对每个Bot应用字段权限脱敏
	for _, bot := range bots {
		maskedBot, err := e.fieldPermChecker.MaskSensitiveFields(ctx, userID, bot)
		if err != nil {
			return fmt.Errorf("mask fields failed: %w", err)
		}

		fmt.Printf("Original: %+v\n", bot)
		fmt.Printf("Masked:  %+v\n", maskedBot)
	}

	return nil
}

// Example3_CreateRoleWithPermissions 示例3：创建角色并分配权限
func (e *IntegrationExample) Example3_CreateRoleWithPermissions(ctx context.Context, tenantID string) error {
	// 场景：创建一个"部门经理"角色，并配置数据权限和字段权限

	// 1. 创建角色
	role := &entity.Role{
		RoleID:      generateID(),
		TenantID:    tenantID,
		RoleName:    "部门经理",
		RoleCode:    "dept_manager",
		RoleType:    entity.RoleTypeCustom,
		Description: "管理部门及子部门的数据",
		Priority:    100,
		IsSystem:    false,
	}

	if err := e.roleRepo.Create(ctx, role); err != nil {
		return fmt.Errorf("create role failed: %w", err)
	}

	// 2. 配置Bot数据权限（本部门及子部门）
	botDataPerm := &entity.DataPermission{
		PermissionID: generateID(),
		RoleID:       role.RoleID,
		ResourceType: entity.ResourceTypeBots,
		Scope:        entity.DataPermissionScopeDepartmentAndSub,
	}

	if err := e.dataPermRepo.Create(ctx, botDataPerm); err != nil {
		return fmt.Errorf("create bot data permission failed: %w", err)
	}

	// 3. 配置字段权限
	fieldPerms := []*entity.FieldPermission{
		{
			PermissionID:    generateID(),
			RoleID:          role.RoleID,
			ResourceType:    "bots",
			FieldName:       "api_key",
			PermissionLevel: entity.FieldPermissionLevelReadonly, // 只读
		},
		{
			PermissionID:    generateID(),
			RoleID:          role.RoleID,
			ResourceType:    "bots",
			FieldName:       "webhook_url",
			PermissionLevel: entity.FieldPermissionLevelHidden, // 隐藏
		},
		{
			PermissionID:    generateID(),
			RoleID:          role.RoleID,
			ResourceType:    "bots",
			FieldName:       "name",
			PermissionLevel: entity.FieldPermissionLevelEditable, // 可编辑
		},
	}

	for _, perm := range fieldPerms {
		if err := e.fieldPermRepo.Create(ctx, perm); err != nil {
			return fmt.Errorf("create field permission failed: %w", err)
		}
	}

	fmt.Printf("角色 %s 创建成功，权限配置完成\n", role.RoleName)
	return nil
}

// Example4_AssignRoleToUser 示例4：分配角色给用户
func (e *IntegrationExample) Example4_AssignRoleToUser(ctx context.Context, userID, tenantID, roleID string) error {
	// 场景：将角色分配给用户

	// 创建用户角色关联
	userRole := &entity.UserRole{
		UserID:    userID,
		TenantID:  tenantID,
		RoleID:    roleID,
		GrantedBy: "system",
		GrantedReason: "系统初始化",
	}

	if err := e.userRoleRepo.Create(ctx, userRole); err != nil {
		return fmt.Errorf("assign role failed: %w", err)
	}

	fmt.Printf("角色分配成功: user=%s, role=%s\n", userID, roleID)
	return nil
}

// Example5_CustomDataPermission 示例5：自定义数据权限过滤
func (e *IntegrationExample) Example5_CustomDataPermission(ctx context.Context, roleID string) error {
	// 场景：创建自定义数据权限（只看已发布的Bot）

	// 自定义过滤器JSON（修正拼写错误：logci -> logic）
	customFilter := `{
		"conditions": [
			{"field": "status", "operator": "=", "value": "published"},
			{"field": "created_at", "operator": ">", "value": "1704067200"}
		],
		"logic": "AND"
	}`

	// 创建自定义数据权限
	customPerm := &entity.DataPermission{
		PermissionID: generateID(),
		RoleID:       roleID,
		ResourceType: entity.ResourceTypeBots,
		Scope:        entity.DataPermissionScopeCustom,
		CustomFilter: customFilter,
	}

	if err := e.dataPermRepo.Create(ctx, customPerm); err != nil {
		return fmt.Errorf("create custom permission failed: %w", err)
	}

	fmt.Println("自定义数据权限创建成功")
	return nil
}

// Example6_CompleteWorkflow 示例6：完整的工作流程
func (e *IntegrationExample) Example6_CompleteWorkflow(ctx context.Context, userID, tenantID string) error {
	// 场景：完整的RBAC权限工作流程
	// 1. 创建角色
	// 2. 配置数据权限和字段权限
	// 3. 分配角色给用户
	// 4. 验证权限

	// 1. 创建角色
	if err := e.Example3_CreateRoleWithPermissions(ctx, tenantID); err != nil {
		return err
	}

	// 获取角色ID
	role, err := e.roleRepo.GetByCode(ctx, tenantID, "dept_manager")
	if err != nil {
		return fmt.Errorf("get role failed: %w", err)
	}

	// 2. 分配角色给用户
	if err := e.Example4_AssignRoleToUser(ctx, userID, tenantID, role.RoleID); err != nil {
		return err
	}

	// 3. 验证权限
	hasAccess, err := e.Example1_SimplePermissionCheck(ctx, userID, "test_bot_id")
	if err != nil {
		return err
	}

	fmt.Printf("权限验证结果: %v\n", hasAccess)

	// 4. 测试字段权限脱敏
	if err := e.Example2_FieldPermissionMask(ctx, userID); err != nil {
		return err
	}

	return nil
}

// Example7_PerformanceTest 示例7：性能测试
func (e *IntegrationExample) Example7_PerformanceTest(ctx context.Context, userID string) error {
	// 场景：测试权限检查性能

	// 测试1：数据权限过滤性能
	start := time.Now()
	filter, err := e.dataPermChecker.Filter(ctx, userID, DEPARTMENT, "bots")
	if err != nil {
		return err
	}
	duration := time.Since(start)
	fmt.Printf("数据权限过滤耗时: %v\n", duration)
	fmt.Printf("过滤条件: %s, 参数: %v\n", filter.WhereClause, filter.Args)

	// 测试2：字段权限脱敏性能
	testData := map[string]interface{}{
		"bot_id":      "bot123",
		"name":        "Test Bot",
		"api_key":     "sk-1234567890",
		"webhook_url": "https://example.com/webhook",
	}

	start = time.Now()
	maskedData, err := e.fieldPermChecker.MaskSensitiveFields(ctx, userID, testData)
	if err != nil {
		return err
	}
	duration = time.Since(start)
	fmt.Printf("字段权限脱敏耗时: %v\n", duration)
	fmt.Printf("脱敏后数据: %+v\n", maskedData)

	return nil
}

// 辅助函数：生成ID
func generateID() string {
	return fmt.Sprintf("id_%d", time.Now().UnixNano())
}
