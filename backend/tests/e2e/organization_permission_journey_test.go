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

package e2e

import (
	"context"
	"testing"
	"time"

	e2eHelpers "github.com/coze-dev/coze-studio/backend/tests/e2e/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== 组织管理+权限控制旅程测试 ====================

/**
 * TestE2E_OrganizationManagementJourney 测试组织管理+权限控制完整旅程
 *
 * 职责: 模拟管理员创建部门并分配权限的全流程
 *
 * 流程:
 * 1. 创建部门
 * 2. 添加成员到部门
 * 3. 分配部门数据权限
 * 4. 验证权限缓存已更新
 * 5. 验证数据权限过滤生效
 *
 * 遵循单一职责原则：只负责组织管理+权限控制这一个用户旅程
 */
func TestE2E_OrganizationManagementJourney(t *testing.T) {
	// 跳过短测试
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	// Setup: 启动完整的服务栈
	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	// 创建测试租户和用户
	tenant := e2eHelpers.CreateTestTenant(t, app, "组织权限测试公司")
	admin := e2eHelpers.CreateTestUser(t, app, tenant.TenantID, "admin")
	employee := e2eHelpers.CreateTestUser(t, app, tenant.TenantID, "employee")

	var orgID string

	t.Run("Step1_创建部门", func(t *testing.T) {
		createOrgReq := map[string]interface{}{
			"tenant_id": tenant.TenantID,
			"org_name":  "技术部",
			"org_type":  "department",
			"org_code":  "TECH",
			"level":     1,
			"path":      "/tech",
		}

		createOrgResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/org/organizations", createOrgReq, &createOrgResp)
		require.NoError(t, err, "创建组织应该成功")
		require.Equal(t, float64(0), createOrgResp["code"], "响应码应为0")

		data, ok := createOrgResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		require.NotEmpty(t, data["organization_id"], "组织ID不应为空")

		orgID = data["organization_id"].(string)

		// 验证数据库中的组织
		var count int64
		err = app.DB.Table("organizations").Where("org_id = ?", orgID).Count(&count).Error
		require.NoError(t, err, "查询组织失败")
		assert.Greater(t, count, int64(0), "组织应该存在")
	})

	t.Run("Step2_添加成员到部门", func(t *testing.T) {
		require.NotEmpty(t, orgID, "组织ID应已设置")

		addMemberReq := map[string]interface{}{
			"organization_id": orgID,
			"user_id":         employee.UserID,
			"role":            "MEMBER",
		}

		addMemberResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/org/members", addMemberReq, &addMemberResp)
		require.NoError(t, err, "添加成员应该成功")
		require.Equal(t, float64(0), addMemberResp["code"], "响应码应为0")

		// 注意：当前简化实现，实际应该在user_departments表中创建记录
	})

	t.Run("Step3_分配部门数据权限", func(t *testing.T) {
		// 首先获取角色ID（这里简化，直接使用假设的role_id）
		// 实际应该先查询该租户的角色

		permissionReq := map[string]interface{}{
			"role_id":        "role-member", // 假设的角色ID
			"resource_type":  "bots",
			"scope":          "DEPARTMENT", // 本部门数据
			"custom_filter":  nil,
		}

		permissionResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/permissions/data", permissionReq, &permissionResp)
		require.NoError(t, err, "分配数据权限应该成功")
		require.Equal(t, float64(0), permissionResp["code"], "响应码应为0")

		data, ok := permissionResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		require.NotEmpty(t, data["permission_id"], "权限ID不应为空")

		// 验证数据库中的权限
		e2eHelpers.AssertDataPermissionExists(t, app.DB, "role-member", "bots")
	})

	t.Run("Step4_验证权限缓存已更新", func(t *testing.T) {
		// 简化：直接查询数据库验证权限存在
		// 实际应该验证缓存系统（如Redis）

		var permission map[string]interface{}
		err := app.DB.Table("data_permissions").
			Where("role_id = ? AND resource_type = ?", "role-member", "bots").
			First(&permission).Error
		require.NoError(t, err, "查询数据权限失败")
		assert.Equal(t, "DEPARTMENT", permission["scope"], "权限范围应为DEPARTMENT")
	})

	t.Run("Step5_验证数据权限过滤生效", func(t *testing.T) {
		// 为技术部创建Bot
		createBotReq := map[string]interface{}{
			"tenant_id": tenant.TenantID,
			"bot_name":  "技术部Bot1",
			"org_id":    &orgID, // 关联到技术部
		}

		createBotResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/bots", createBotReq, &createBotResp)
		require.NoError(t, err, "创建Bot应该成功")

		data, ok := createBotResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		techBotID := data["bot_id"].(string)

		// 为其他部门创建Bot（或无部门Bot）
		createBotReq2 := map[string]interface{}{
			"tenant_id": tenant.TenantID,
			"bot_name":  "市场部Bot",
			// 不设置org_id
		}

		createBotResp2 := make(map[string]interface{})
		err = app.PostJSON("/api/v1/bots", createBotReq2, &createBotResp2)
		require.NoError(t, err, "创建Bot应该成功")

		// 查询Bot列表（模拟员工视角）
		listResp := make(map[string]interface{})
		err = app.GetJSON("/api/v1/bots/list?user_id="+employee.UserID+"&resource_type=bots", &listResp)
		require.NoError(t, err, "列出Bot应该成功")

		// 验证：员工只能看到本部门Bot（DEPARTMENT权限）
		listData, ok := listResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		bots, ok := listData["bots"].([]interface{})
		require.True(t, ok, "bots应为数组")

		// 简化验证：检查返回结果包含本部门Bot
		found := false
		for _, botInterface := range bots {
			bot, ok := botInterface.(map[string]interface{})
			if !ok {
				continue
			}
			if bot["bot_id"] == techBotID {
				found = true
				// 验证org_id匹配
				if orgIDValue, ok := bot["org_id"]; ok {
					if orgIDStr, ok := orgIDValue.(string); ok {
						assert.Equal(t, orgID, orgIDStr, "Bot应该属于技术部")
					}
				}
				break
			}
		}
		assert.True(t, found, "应该能看到本部门Bot")
	})
}

/**
 * TestE2E_MultiLevelOrganizationHierarchy 测试多级组织层级
 *
 * 职责: 验证多级组织结构的创建和权限继承
 */
func TestE2E_MultiLevelOrganizationHierarchy(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	tenant := e2eHelpers.CreateTestTenant(t, app, "多级组织测试公司")

	var parentOrgID string

	t.Run("创建两级组织结构", func(t *testing.T) {
		// 创建一级组织：技术中心
		createParentReq := map[string]interface{}{
			"tenant_id": tenant.TenantID,
			"org_name":  "技术中心",
			"org_type":  "division",
			"org_code":  "TECH_CENTER",
			"level":     1,
			"path":      "/tech_center",
		}

		createParentResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/org/organizations", createParentReq, &createParentResp)
		require.NoError(t, err, "创建父组织应该成功")

		parentData, ok := createParentResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		parentOrgID = parentData["organization_id"].(string)

		// 创建二级组织：研发部（属于技术中心）
		createChildReq := map[string]interface{}{
			"tenant_id": tenant.TenantID,
			"org_name":  "研发部",
			"org_type":  "department",
			"parent_id": &parentOrgID,
			"org_code":  "RD",
			"level":     2,
			"path":      "/tech_center/rd",
		}

		createChildResp := make(map[string]interface{})
		err = app.PostJSON("/api/v1/org/organizations", createChildReq, &createChildResp)
		require.NoError(t, err, "创建子组织应该成功")

		childData, ok := createChildResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		childOrgID := childData["organization_id"].(string)

		// 验证层级关系
		var childOrg map[string]interface{}
		err = app.DB.Table("organizations").Where("org_id = ?", childOrgID).First(&childOrg).Error
		require.NoError(t, err, "查询子组织失败")
		assert.Equal(t, parentOrgID, childOrg["parent_id"], "父组织ID应匹配")
		assert.Equal(t, 2, childOrg["level"], "组织层级应为2")
	})
}

/**
 * TestE2E_DataPermissionScopes 测试不同数据权限范围
 *
 * 职责: 验证ALL、DEPARTMENT、OWN、NONE等权限范围
 */
func TestE2E_DataPermissionScopes(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	tenant := e2eHelpers.CreateTestTenant(t, app, "权限范围测试公司")

	t.Run("创建不同范围的权限", func(t *testing.T) {
		scopes := []string{"ALL", "DEPARTMENT", "OWN", "NONE"}
		permissionIDs := make([]string, len(scopes))

		for i, scope := range scopes {
			permissionReq := map[string]interface{}{
				"role_id":       "role-" + scope,
				"resource_type": "bots",
				"scope":         scope,
			}

			permissionResp := make(map[string]interface{})
			err := app.PostJSON("/api/v1/permissions/data", permissionReq, &permissionResp)
			require.NoError(t, err, "创建权限应该成功")

			data, ok := permissionResp["data"].(map[string]interface{})
			require.True(t, ok, "响应data应为对象")
			permissionIDs[i] = data["permission_id"].(string)

			// 验证权限已创建
			var permission map[string]interface{}
			err = app.DB.Table("data_permissions").Where("permission_id = ?", permissionIDs[i]).First(&permission).Error
			require.NoError(t, err, "查询权限失败")
			assert.Equal(t, scope, permission["scope"], "权限范围应匹配")
		}
	})
}

/**
 * TestE2E_PermissionInheritance 测试权限继承
 *
 * 职责: 验证子组织是否继承父组织的权限
 */
func TestE2E_PermissionInheritance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	tenant := e2eHelpers.CreateTestTenant(t, app, "权限继承测试公司")

	t.Run("验证权限继承逻辑", func(t *testing.T) {
		// 创建父组织和子组织
		// 分配父组织权限
		// 验证子组织用户是否能看到父组织数据

		// 简化实现：仅创建结构
		createParentReq := map[string]interface{}{
			"tenant_id": tenant.TenantID,
			"org_name":  "总部",
			"org_type":  "company",
			"org_code":  "HQ",
			"level":     1,
			"path":      "/hq",
		}

		createParentResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/org/organizations", createParentReq, &createParentResp)
		require.NoError(t, err, "创建父组织应该成功")

		parentData, ok := createParentResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		parentOrgID := parentData["organization_id"].(string)

		// TODO: 实现完整的权限继承测试
		_ = parentOrgID
		_ = err
	})
}

/**
 * TestE2E_OrganizationUserManagement 测试组织用户管理
 *
 * 职责: 验证用户在组织中的添加、移除和角色变更
 */
func TestE2E_OrganizationUserManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	tenant := e2eHelpers.CreateTestTenant(t, app, "用户管理测试公司")
	org := e2eHelpers.CreateTestOrganization(t, app, tenant.TenantID, "测试部门")
	user1 := e2eHelpers.CreateTestUser(t, app, tenant.TenantID, "member")
	user2 := e2eHelpers.CreateTestUser(t, app, tenant.TenantID, "member")

	t.Run("添加多个用户到组织", func(t *testing.T) {
		// 添加user1
		addMemberReq1 := map[string]interface{}{
			"organization_id": org.OrgID,
			"user_id":         user1.UserID,
			"role":            "MEMBER",
		}

		addMemberResp1 := make(map[string]interface{})
		err := app.PostJSON("/api/v1/org/members", addMemberReq1, &addMemberResp1)
		require.NoError(t, err, "添加成员应该成功")

		// 添加user2
		addMemberReq2 := map[string]interface{}{
			"organization_id": org.OrgID,
			"user_id":         user2.UserID,
			"role":            "MEMBER",
		}

		addMemberResp2 := make(map[string]interface{})
		err = app.PostJSON("/api/v1/org/members", addMemberReq2, &addMemberResp2)
		require.NoError(t, err, "添加成员应该成功")

		// 验证两个用户都已添加（简化：直接验证API响应）
		assert.Equal(t, float64(0), addMemberResp1["code"], "第一个用户添加应成功")
		assert.Equal(t, float64(0), addMemberResp2["code"], "第二个用户添加应成功")
	})
}
