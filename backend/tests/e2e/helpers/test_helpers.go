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

package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ==================== 单一职责：每个函数只负责一个创建或验证任务 ====================

/**
 * createTestTenant 创建测试租户
 *
 * 职责: 创建E2E测试用的租户
 * 遵循单一职责原则：只负责创建租户
 */
func CreateTestTenant(t *testing.T, app *e2e.TestApplication, name string) TestTenant {
	req := map[string]interface{}{
		"tenant_name": name,
		"tenant_type": "team",
		"subdomain":   fmt.Sprintf("tenant-%d", time.Now().UnixNano()),
		"status":      "active",
	}

	resp := make(map[string]interface{})
	err := app.PostJSON("/api/v1/tenants", req, &resp)
	require.NoError(t, err, "创建租户失败")
	require.Equal(t, float64(0), resp["code"], "创建租户响应码应为0")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应data应为对象")
	require.NotEmpty(t, data["tenant_id"], "租户ID不应为空")

	// 设置当前租户ID
	app.TenantID = data["tenant_id"].(string)

	return TestTenant{
		TenantID:   data["tenant_id"].(string),
		TenantName: name,
	}
}

/**
 * TestTenant 测试租户结构
 */
type TestTenant struct {
	TenantID   string
	TenantName string
}

/**
 * createTestUser 创建测试用户
 *
 * 职责: 创建E2E测试用的用户
 * 遵循单一职责原则：只负责创建用户
 */
func CreateTestUser(t *testing.T, app *e2e.TestApplication, tenantID, role string) TestUser {
	req := map[string]interface{}{
		"tenant_id": tenantID,
		"username":  fmt.Sprintf("user_%d", time.Now().UnixNano()),
		"role":      role,
	}

	resp := make(map[string]interface{})
	err := app.PostJSON("/api/v1/users", req, &resp)
	require.NoError(t, err, "创建用户失败")
	require.Equal(t, float64(0), resp["code"], "创建用户响应码应为0")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应data应为对象")
	require.NotEmpty(t, data["user_id"], "用户ID不应为空")

	return TestUser{
		UserID:   data["user_id"].(string),
		Username: req["username"].(string),
		Role:     role,
	}
}

/**
 * TestUser 测试用户结构
 */
type TestUser struct {
	UserID   string
	Username string
	Role     string
}

/**
 * createTestBot 创建测试Bot
 *
 * 职责: 创建E2E测试用的Bot
 * 遵循单一职责原则：只负责创建Bot
 */
func CreateTestBot(t *testing.T, app *e2e.TestApplication, tenantID, name string, aiConfidence float64) TestBot {
	req := map[string]interface{}{
		"tenant_id":    tenantID,
		"bot_name":     name,
		"ai_confidence": aiConfidence,
	}

	resp := make(map[string]interface{})
	err := app.PostJSON("/api/v1/bots", req, &resp)
	require.NoError(t, err, "创建Bot失败")
	require.Equal(t, float64(0), resp["code"], "创建Bot响应码应为0")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应data应为对象")
	require.NotEmpty(t, data["bot_id"], "Bot ID不应为空")

	return TestBot{
		BotID: data["bot_id"].(string),
		Name:  name,
		Status: data["status"].(string),
	}
}

/**
 * TestBot 测试Bot结构
 */
type TestBot struct {
	BotID  string
	Name   string
	Status string
}

/**
 * createTestOrganization 创建测试组织
 *
 * 职责: 创建E2E测试用的组织
 * 遵循单一职责原则：只负责创建组织
 */
func CreateTestOrganization(t *testing.T, app *e2e.TestApplication, tenantID, name string) TestOrganization {
	req := map[string]interface{}{
		"tenant_id": tenantID,
		"org_name":  name,
		"org_type":  "department",
		"org_code":  fmt.Sprintf("org-%d", time.Now().UnixNano()),
		"level":     1,
		"path":      fmt.Sprintf("/org-%d", time.Now().UnixNano()),
	}

	resp := make(map[string]interface{})
	err := app.PostJSON("/api/v1/org/organizations", req, &resp)
	require.NoError(t, err, "创建组织失败")
	require.Equal(t, float64(0), resp["code"], "创建组织响应码应为0")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应data应为对象")
	require.NotEmpty(t, data["organization_id"], "组织ID不应为空")

	return TestOrganization{
		OrgID:   data["organization_id"].(string),
		OrgName: name,
	}
}

/**
 * TestOrganization 测试组织结构
 */
type TestOrganization struct {
	OrgID   string
	OrgName string
}

// ==================== 验证辅助函数 ====================

/**
 * assertTenantExists 验证租户存在
 *
 * 职责: 验证数据库中存在指定租户
 * 遵循单一职责原则：只负责验证租户存在性
 */
func AssertTenantExists(t *testing.T, db *gorm.DB, tenantID string) {
	var count int64
	err := db.Table("tenants").Where("tenant_id = ?", tenantID).Count(&count).Error
	require.NoError(t, err, "查询租户失败")
	assert.Greater(t, count, int64(0), "租户应该存在")
}

/**
 * assertRoleExists 验证角色存在
 *
 * 职责: 验证数据库中存在指定角色
 */
func AssertRoleExists(t *testing.T, db *gorm.DB, tenantID, roleCode string) {
	var count int64
	err := db.Table("roles").Where("tenant_id = ? AND role_code = ?", tenantID, roleCode).Count(&count).Error
	require.NoError(t, err, "查询角色失败")
	assert.Greater(t, count, int64(0), "角色应该存在")
}

/**
 * assertQuotaExists 验证配额存在
 *
 * 职责: 验证数据库中存在指定配额
 */
func AssertQuotaExists(t *testing.T, db *gorm.DB, tenantID string) {
	var count int64
	err := db.Table("quotas").Where("tenant_id = ?", tenantID).Count(&count).Error
	require.NoError(t, err, "查询配额失败")
	assert.Greater(t, count, int64(0), "配额应该存在")
}

/**
 * assertBotStatus 验证Bot状态
 *
 * 职责: 验证Bot的状态是否符合预期
 */
func AssertBotStatus(t *testing.T, db *gorm.DB, botID, expectedStatus string) {
	var status string
	err := db.Table("bots").Where("bot_id = ?", botID).Select("status").Scan(&status).Error
	require.NoError(t, err, "查询Bot状态失败")
	assert.Equal(t, expectedStatus, status, "Bot状态应该匹配")
}

/**
 * assertReviewTaskExists 验证审核任务存在
 *
 * 职责: 验证数据库中存在指定状态的审核任务
 */
func AssertReviewTaskExists(t *testing.T, db *gorm.DB, tenantID, status string) int64 {
	var count int64
	err := db.Table("review_tasks").Where("tenant_id = ? AND status = ?", tenantID, status).Count(&count).Error
	require.NoError(t, err, "查询审核任务失败")
	assert.Greater(t, count, int64(0), "审核任务应该存在")
	return count
}

/**
 * assertMemoryExists 验证记忆存在
 *
 * 职责: 验证数据库中存在指定用户的记忆
 */
func AssertMemoryExists(t *testing.T, db *gorm.DB, userID, entityType string) int64 {
	var count int64
	query := db.Table("entity_memories").Where("user_id = ?", userID)
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}
	err := query.Count(&count).Error
	require.NoError(t, err, "查询记忆失败")
	assert.Greater(t, count, int64(0), "记忆应该存在")
	return count
}

/**
 * assertDataPermissionExists 验证数据权限存在
 *
 * 职责: 验证数据库中存在指定角色的数据权限
 */
func AssertDataPermissionExists(t *testing.T, db *gorm.DB, roleID, resourceType string) {
	var count int64
	err := db.Table("data_permissions").Where("role_id = ? AND resource_type = ?", roleID, resourceType).Count(&count).Error
	require.NoError(t, err, "查询数据权限失败")
	assert.Greater(t, count, int64(0), "数据权限应该存在")
}

/**
 * getTenantByID 通过ID获取租户
 *
 * 职责: 从数据库获取租户信息
 */
func GetTenantByID(t *testing.T, db *gorm.DB, tenantID string) map[string]interface{} {
	var tenant map[string]interface{}
	err := db.Table("tenants").Where("tenant_id = ?", tenantID).First(&tenant).Error
	require.NoError(t, err, "获取租户失败")
	return tenant
}

/**
 * getBotByID 通过ID获取Bot
 *
 * 职责: 从数据库获取Bot信息
 */
func GetBotByID(t *testing.T, db *gorm.DB, botID string) map[string]interface{} {
	var bot map[string]interface{}
	err := db.Table("bots").Where("bot_id = ?", botID).First(&bot).Error
	require.NoError(t, err, "获取Bot失败")
	return bot
}

/**
 * listRolesByTenantID 列出租户的角色
 *
 * 职责: 获取租户的所有角色
 */
func ListRolesByTenantID(t *testing.T, db *gorm.DB, tenantID string) []map[string]interface{} {
	var roles []map[string]interface{}
	err := db.Table("roles").Where("tenant_id = ?", tenantID).Find(&roles).Error
	require.NoError(t, err, "列出租户角色失败")
	return roles
}

/**
 * cleanupTenantData 清理租户数据
 *
 * 职责: 清理指定租户的所有测试数据
 * 遵循测试隔离原则：每个测试后清理数据
 */
func CleanupTenantData(t *testing.T, ctx context.Context, db *gorm.DB, tenantID string) {
	// 按依赖顺序删除表
	tables := []string{
		"entity_memories",
		"messages",
		"conversations",
		"review_tasks",
		"bots",
		"data_permissions",
		"user_roles",
		"roles",
		"organizations",
		"users",
		"quotas",
		"subscriptions",
		"tenants",
	}

	for _, table := range tables {
		err := db.Exec("DELETE FROM ? WHERE tenant_id = ?", table, tenantID).Error
		if err != nil {
			t.Logf("清理表 %s 失败: %v", table, err)
		}
	}
}

/**
 * waitForCondition 等待条件满足
 *
 * 职责: 轮询等待某个条件满足（用于异步操作）
 */
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, checkInterval time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(checkInterval)
	}
	return false
}

/**
 * retry 重试函数
 *
 * 职责: 重试执行函数直到成功或达到最大次数
 */
func Retry(t *testing.T, maxAttempts int, fn func() error) error {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
			t.Logf("重试 %d/%d: %v", i+1, maxAttempts, err)
			time.Sleep(time.Second)
		}
	}
	return lastErr
}
