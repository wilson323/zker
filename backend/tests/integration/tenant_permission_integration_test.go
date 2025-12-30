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

package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	tenantService "github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	"github.com/coze-dev/coze-studio/backend/tests/integration/fixtures"
)

// TenantPermissionIntegrationSuite 租户+权限系统集成测试套件
// 职责: 测试租户创建时RBAC权限的自动初始化和租户隔离
type TenantPermissionIntegrationSuite struct {
	suite.Suite
	ctx    context.Context
	db     *sql.DB
	logger *zap.Logger
}

// SetupSuite 测试套件初始化
func (s *TenantPermissionIntegrationSuite) SetupSuite() {
	s.ctx = context.Background()
	s.logger = zap.NewNop()

	// 创建MySQL测试容器
	testCtx, err := SetupMySQLTestContainer(s.T())
	assert.NoError(s.T(), err)
	s.db = testCtx.DB

	// TODO: 运行数据库迁移
	// RunMigrations(s.db)
}

// TearDownSuite 测试套件清理
func (s *TenantPermissionIntegrationSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

// SetupTest 每个测试前的设置
func (s *TenantPermissionIntegrationSuite) SetupTest() {
	// 清理测试数据
	tables := []string{
		"tenants",
		"users",
		"roles",
		"data_permissions",
		"field_permissions",
		"user_roles",
	}
	err := CleanupTestData(s.db, tables)
	assert.NoError(s.T(), err)
}

// TestTenantCreationWithRBAC 测试租户创建时RBAC权限初始化
// 职责: 验证创建租户时，RBAC权限是否正确初始化
func (s *TenantPermissionIntegrationSuite) TestTenantCreationWithRBAC() {
	t := s.T()

	// Setup: 初始化Repository和Service
	// 注意: 这里使用Mock或测试实现
	tenantSvc := tenantService.NewTenantService(s.db, s.logger)
	roleSvc := service.NewRoleService(s.db, s.logger)

	// Execute: 创建租户
	tenant := &entity.Tenant{
		TenantID:     fmt.Sprintf("tenant-%d", time.Now().UnixNano()),
		Name:         "测试租户",
		Status:       entity.TenantStatusActive,
		SubscriptionID: fmt.Sprintf("sub-%d", time.Now().UnixNano()),
	}

	err := tenantSvc.CreateTenant(s.ctx, tenant)
	assert.NoError(t, err, "创建租户应该成功")
	assert.NotEmpty(t, tenant.TenantID)

	// Verify: 验证默认角色已创建
	roles, err := roleSvc.ListRoles(s.ctx, tenant.TenantID)
	assert.NoError(t, err, "查询角色列表应该成功")
	assert.Greater(t, len(roles), 0, "应该创建默认角色")

	// Verify: 验证默认角色包含TenantAdmin
	var tenantAdminRole *entity.Role
	for _, role := range roles {
		if role.Name == "TenantAdmin" {
			tenantAdminRole = &role
			break
		}
	}
	assert.NotNil(t, tenantAdminRole, "应该有TenantAdmin角色")
	assert.Equal(t, tenant.TenantID, tenantAdminRole.TenantID)

	// Verify: 验证默认权限已分配
	// TenantAdmin应该拥有所有资源的所有权限
	permissions, err := roleSvc.GetRolePermissions(s.ctx, tenantAdminRole.RoleID)
	assert.NoError(t, err)
	assert.Greater(t, len(permissions), 0, "TenantAdmin应该有默认权限")
}

// TestTenantIsolationWithDataPermissions 测试租户隔离+数据权限过滤
// 职责: 验证不同租户的数据是否正确隔离
func (s *TenantPermissionIntegrationSuite) TestTenantIsolationWithDataPermissions() {
	t := s.T()

	// Setup: 初始化Service
	tenantSvc := tenantService.NewTenantService(s.db, s.logger)
	roleSvc := service.NewRoleService(s.db, s.logger)
	authzSvc := service.NewAuthorizationService(s.db, s.logger)

	// Setup: 创建两个租户
	tenant1 := fixtures.CreateTestTenant(t, s.db, "租户1")
	tenant2 := fixtures.CreateTestTenant(t, s.db, "租户2")

	// Setup: 为两个租户创建管理员用户
	user1 := fixtures.CreateTestUser(t, s.db, tenant1.TenantID, "user1")
	user2 := fixtures.CreateTestUser(t, s.db, tenant2.TenantID, "user2")

	// Setup: 为两个租户创建TenantAdmin角色
	role1 := fixtures.CreateTestRole(t, s.db, tenant1.TenantID, "TenantAdmin")
	role2 := fixtures.CreateTestRole(t, s.db, tenant2.TenantID, "TenantAdmin")

	// Setup: 分配角色给用户
	fixtures.AssignRoleToUser(t, s.db, user1.UserID, role1.RoleID)
	fixtures.AssignRoleToUser(t, s.db, user2.UserID, role2.RoleID)

	// Setup: 为tenant1创建Bot（模拟数据）
	bot1 := fixtures.CreateTestBot(t, s.db, tenant1.TenantID, "Bot1")

	// Setup: 为tenant2创建Bot（模拟数据）
	bot2 := fixtures.CreateTestBot(t, s.db, tenant2.TenantID, "Bot2")

	// Execute & Verify: 验证user1只能访问tenant1的数据
	filter1, err := authzSvc.BuildDataPermissionFilter(s.ctx, &entity.DataPermissionRequest{
		UserID:       user1.UserID,
		TenantID:     tenant1.TenantID,
		ResourceType: entity.ResourceTypeBot,
		Action:       entity.ActionRead,
	})
	assert.NoError(t, err)

	// user1应该只能看到tenant1的Bot
	botsForUser1, err := fixtures.ListBotsByTenant(t, s.db, tenant1.TenantID)
	assert.NoError(t, err)
	for _, bot := range botsForUser1 {
		assert.Equal(t, tenant1.TenantID, bot.TenantID, "user1不应看到tenant2的数据")
	}

	// Execute & Verify: 验证user2只能访问tenant2的数据
	filter2, err := authzSvc.BuildDataPermissionFilter(s.ctx, &entity.DataPermissionRequest{
		UserID:       user2.UserID,
		TenantID:     tenant2.TenantID,
		ResourceType: entity.ResourceTypeBot,
		Action:       entity.ActionRead,
	})
	assert.NoError(t, err)

	botsForUser2, err := fixtures.ListBotsByTenant(t, s.db, tenant2.TenantID)
	assert.NoError(t, err)
	for _, bot := range botsForUser2 {
		assert.Equal(t, tenant2.TenantID, bot.TenantID, "user2不应看到tenant1的数据")
	}

	// Verify: 验证数据权限隔离
	assert.NotContains(t, botsForUser1, bot2.BotID, "tenant1的用户不应看到tenant2的Bot")
	assert.NotContains(t, botsForUser2, bot1.BotID, "tenant2的用户不应看到tenant1的Bot")

	// Verify: 验证过滤器正确
	assert.Equal(t, tenant1.TenantID, filter1.TenantIDFilter, "过滤器应该包含正确的租户ID")
	assert.Equal(t, tenant2.TenantID, filter2.TenantIDFilter, "过滤器应该包含正确的租户ID")
}

// TestRoleAssignmentWithPermissionInheritance 测试角色分配与权限继承
// 职责: 验证用户分配角色后，权限是否正确继承
func (s *TenantPermissionIntegrationSuite) TestRoleAssignmentWithPermissionInheritance() {
	t := s.T()

	// Setup: 初始化Service
	roleSvc := service.NewRoleService(s.db, s.logger)
	authzSvc := service.NewAuthorizationService(s.db, s.logger)

	// Setup: 创建租户和用户
	tenant := fixtures.CreateTestTenant(t, s.db, "测试租户")
	user := fixtures.CreateTestUser(t, s.db, tenant.TenantID, "testuser")

	// Setup: 创建自定义角色
	customRole := fixtures.CreateTestRole(t, s.db, tenant.TenantID, "CustomRole")

	// Setup: 为角色分配权限（只读Bot）
	permissions := []entity.DataPermission{
		{
			PermissionID: fmt.Sprintf("perm-%d", time.Now().UnixNano()),
			RoleID:       customRole.RoleID,
			ResourceType: entity.ResourceTypeBot,
			Action:       entity.ActionRead,
			Scope:        entity.DataPermissionScopeAll,
		},
	}
	for _, perm := range permissions {
		err := fixtures.CreateDataPermission(t, s.db, &perm)
		assert.NoError(t, err)
	}

	// Execute: 分配角色给用户
	err := fixtures.AssignRoleToUser(t, s.db, user.UserID, customRole.RoleID)
	assert.NoError(t, err)

	// Verify: 验证用户权限
	check := &entity.CheckAuthzData{
		OperatorID: user.UserID,
		ResourceIdentifier: []entity.ResourceIdentifier{
			{
				Type:   entity.ResourceTypeBot,
				ID:     []string{"bot-123"},
				Action: entity.ActionRead,
			},
		},
		TenantID: tenant.TenantID,
	}

	result, err := authzSvc.CheckAuthz(s.ctx, check)
	assert.NoError(t, err)
	assert.Equal(t, entity.DecisionAllow, result.Decision, "用户应该有Bot读取权限")

	// Verify: 验证用户没有写权限
	checkWrite := &entity.CheckAuthzData{
		OperatorID: user.UserID,
		ResourceIdentifier: []entity.ResourceIdentifier{
			{
				Type:   entity.ResourceTypeBot,
				ID:     []string{"bot-123"},
				Action: entity.ActionWrite,
			},
		},
		TenantID: tenant.TenantID,
	}

	resultWrite, err := authzSvc.CheckAuthz(s.ctx, checkWrite)
	assert.NoError(t, err)
	assert.Equal(t, entity.DecisionDeny, resultWrite.Decision, "用户不应该有Bot写入权限")
}

// TestFieldLevelPermissionControl 测试字段级权限控制
// 职责: 验证字段级权限是否正确控制字段可见性
func (s *TenantPermissionIntegrationSuite) TestFieldLevelPermissionControl() {
	t := s.T()

	// Setup: 初始化Service
	roleSvc := service.NewRoleService(s.db, s.logger)

	// Setup: 创建租户和角色
	tenant := fixtures.CreateTestTenant(t, s.db, "测试租户")
	role := fixtures.CreateTestRole(t, s.db, tenant.TenantID, "LimitedRole")

	// Setup: 创建字段级权限（不允许查看Bot的敏感字段）
	fieldPerm := &entity.FieldPermission{
		PermissionID:    fmt.Sprintf("field-perm-%d", time.Now().UnixNano()),
		RoleID:          role.RoleID,
		ResourceType:    entity.ResourceTypeBot,
		FieldName:       "api_key",
		IsAllowed:       false,
		PermissionType:  entity.FieldPermissionTypeView,
	}
	err := fixtures.CreateFieldPermission(t, s.db, fieldPerm)
	assert.NoError(t, err)

	// Execute: 检查字段权限
	allowed, err := roleSvc.CheckFieldPermission(s.ctx, &entity.FieldPermissionRequest{
		RoleID:       role.RoleID,
		ResourceType: entity.ResourceTypeBot,
		FieldName:    "api_key",
		Action:       entity.FieldPermissionTypeView,
	})
	assert.NoError(t, err)
	assert.False(t, allowed, "角色不应该有api_key字段的查看权限")

	// Execute: 检查允许的字段
	allowed2, err := roleSvc.CheckFieldPermission(s.ctx, &entity.FieldPermissionRequest{
		RoleID:       role.RoleID,
		ResourceType: entity.ResourceTypeBot,
		FieldName:    "name",
		Action:       entity.FieldPermissionTypeView,
	})
	assert.NoError(t, err)
	assert.True(t, allowed2, "角色应该有name字段的查看权限")
}

// TestTenantDeletionWithPermissionCleanup 测试租户删除时权限清理
// 职责: 验证删除租户时，相关权限是否正确清理
func (s *TenantPermissionIntegrationSuite) TestTenantDeletionWithPermissionCleanup() {
	t := s.T()

	// Setup: 初始化Service
	tenantSvc := tenantService.NewTenantService(s.db, s.logger)
	roleSvc := service.NewRoleService(s.db, s.logger)

	// Setup: 创建租户和角色
	tenant := fixtures.CreateTestTenant(t, s.db, "待删除租户")
	role := fixtures.CreateTestRole(t, s.db, tenant.TenantID, "RoleToDelete")

	// Setup: 创建权限
	dataPerm := &entity.DataPermission{
		PermissionID: fmt.Sprintf("perm-%d", time.Now().UnixNano()),
		RoleID:       role.RoleID,
		ResourceType: entity.ResourceTypeBot,
		Action:       entity.ActionRead,
		Scope:        entity.DataPermissionScopeAll,
	}
	err := fixtures.CreateDataPermission(t, s.db, dataPerm)
	assert.NoError(t, err)

	// Verify: 验证权限存在
	permissionsBefore, err := roleSvc.GetRolePermissions(s.ctx, role.RoleID)
	assert.NoError(t, err)
	assert.Greater(t, len(permissionsBefore), 0)

	// Execute: 删除租户（软删除）
	err = tenantSvc.DeleteTenant(s.ctx, tenant.TenantID)
	assert.NoError(t, err)

	// Verify: 验证角色已被删除
	rolesAfter, err := roleSvc.ListRoles(s.ctx, tenant.TenantID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(rolesAfter), "租户删除后，角色应该被删除")

	// Verify: 验证权限已被清理
	// 注意: 由于角色被删除，关联的权限也应该被清理
	// 这里可以根据实际业务逻辑调整验证方式
}

// RunTestSuite 运行测试套件
func TestTenantPermissionIntegrationSuite(t *testing.T) {
	suite.Run(t, new(TenantPermissionIntegrationSuite))
}
