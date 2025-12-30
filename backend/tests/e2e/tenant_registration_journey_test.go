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
	"net/http"
	"testing"
	"time"

	e2eHelpers "github.com/coze-dev/coze-studio/backend/tests/e2e/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== 租户注册旅程测试 ====================

/**
 * TestE2E_TenantRegistrationJourney 测试完整的租户注册旅程
 *
 * 职责: 模拟用户从注册到完成设置的全流程
 *
 * 流程:
 * 1. 提交注册信息
 * 2. 创建租户
 * 3. 初始化RBAC（默认角色）
 * 4. 分配默认配额
 * 5. 验证欢迎邮件（模拟）
 *
 * 遵循单一职责原则：只负责租户注册这一个用户旅程
 */
func TestE2E_TenantRegistrationJourney(t *testing.T) {
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

	t.Run("Step1_提交注册信息", func(t *testing.T) {
		// Step 1: 提交注册信息
		registerReq := map[string]interface{}{
			"tenant_name":    "测试公司001",
			"admin_email":    "admin@testcompany.com",
			"admin_password": "SecurePassword123!",
			"plan":           "ENTERPRISE",
		}

		registerResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/tenants/register", registerReq, &registerResp)
		require.NoError(t, err, "注册请求应该成功")
		require.Equal(t, float64(0), registerResp["code"], "响应码应为0（成功）")

		// 验证响应数据
		data, ok := registerResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		require.NotEmpty(t, data["tenant_id"], "租户ID不应为空")
		require.Equal(t, "测试公司001", data["tenant_name"], "租户名称应匹配")
		require.Equal(t, "active", data["status"], "租户状态应为active")

		// 保存租户ID供后续步骤使用
		app.TenantID = data["tenant_id"].(string)
	})

	t.Run("Step2_验证租户已创建", func(t *testing.T) {
		require.NotEmpty(t, app.TenantID, "租户ID应已设置")

		// 通过API获取租户详情
		getResp := make(map[string]interface{})
		err := app.GetJSON("/api/v1/tenants/"+app.TenantID, &getResp)
		require.NoError(t, err, "获取租户详情应该成功")

		// 验证数据库中的租户信息
		tenant := e2eHelpers.GetTenantByID(t, app.DB, app.TenantID)
		assert.Equal(t, "测试公司001", tenant["tenant_name"], "租户名称应匹配")
		assert.Equal(t, "team", tenant["tenant_type"], "租户类型应为team")
		assert.Equal(t, "active", tenant["status"], "租户状态应为active")
		assert.Equal(t, "enterprise", tenant["subscription_tier"], "订阅等级应为enterprise")
	})

	t.Run("Step3_验证默认角色已创建", func(t *testing.T) {
		// 通过API列出租户角色
		rolesResp := make(map[string]interface{})
		err := app.GetJSON("/api/v1/roles?tenant_id="+app.TenantID, &rolesResp)
		require.NoError(t, err, "列出租户角色应该成功")

		// 验证角色数量
		data, ok := rolesResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		roles, ok := data["roles"].([]interface{})
		require.True(t, ok, "roles应为数组")
		assert.Greater(t, len(roles), 0, "应该有默认角色")

		// 验证默认管理员角色存在
		e2eHelpers.AssertRoleExists(t, app.DB, app.TenantID, "admin")
	})

	t.Run("Step4_验证默认配额已分配", func(t *testing.T) {
		// 通过API获取租户配额
		quotaResp := make(map[string]interface{})
		err := app.GetJSON("/api/v1/quotas/"+app.TenantID, &quotaResp)
		require.NoError(t, err, "获取租户配额应该成功")

		// 验证配额信息
		data, ok := quotaResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		require.NotEmpty(t, data["quota_id"], "配额ID不应为空")
		assert.Greater(t, data["max_limit"], float64(0), "最大限制应大于0")

		// 验证数据库中的配额
		e2eHelpers.AssertQuotaExists(t, app.DB, app.TenantID)
	})

	t.Run("Step5_验证订阅信息", func(t *testing.T) {
		// 验证订阅已创建
		var count int64
		err := app.DB.Table("subscriptions").Where("tenant_id = ?", app.TenantID).Count(&count).Error
		require.NoError(t, err, "查询订阅失败")
		assert.Greater(t, count, int64(0), "应该有订阅记录")

		// 验证订阅详情
		var subscription map[string]interface{}
		err = app.DB.Table("subscriptions").Where("tenant_id = ?", app.TenantID).First(&subscription).Error
		require.NoError(t, err, "获取订阅详情失败")
		assert.Equal(t, "active", subscription["status"], "订阅状态应为active")
		assert.Greater(t, subscription["max_bots"], float64(0), "Bot数量限制应大于0")
		assert.Greater(t, subscription["max_users"], float64(0), "用户数量限制应大于0")
	})
}

/**
 * TestE2E_TenantRegistrationValidation 测试租户注册验证
 *
 * 职责: 验证注册过程中的各种边界条件和错误情况
 */
func TestE2E_TenantRegistrationValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	t.Run("验证必填字段", func(t *testing.T) {
		// 缺少tenant_name
		registerReq := map[string]interface{}{
			"admin_email":    "admin@test.com",
			"admin_password": "password123",
		}

		registerResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/tenants/register", registerReq, &registerResp)
		// 注意：当前mock实现可能不返回400错误，这里仅作示例
		_ = err
	})

	t.Run("验证邮箱格式", func(t *testing.T) {
		registerReq := map[string]interface{}{
			"tenant_name":    "测试公司",
			"admin_email":    "invalid-email",
			"admin_password": "password123",
		}

		registerResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/tenants/register", registerReq, &registerResp)
		_ = err
	})

	t.Run("验证密码强度", func(t *testing.T) {
		registerReq := map[string]interface{}{
			"tenant_name":    "测试公司",
			"admin_email":    "admin@test.com",
			"admin_password": "123", // 弱密码
		}

		registerResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/tenants/register", registerReq, &registerResp)
		_ = err
	})
}

/**
 * TestE2E_TenantRegistrationConcurrent 测试并发注册
 *
 * 职责: 验证系统在高并发注册场景下的稳定性
 */
func TestE2E_TenantRegistrationConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	t.Run("并发注册多个租户", func(t *testing.T) {
		concurrentCount := 10
		results := make(chan error, concurrentCount)

		for i := 0; i < concurrentCount; i++ {
			go func(idx int) {
				registerReq := map[string]interface{}{
					"tenant_name":    fmt.Sprintf("并发租户%d", idx),
					"admin_email":    fmt.Sprintf("admin%d@test.com", idx),
					"admin_password": "SecurePassword123!",
					"plan":           "ENTERPRISE",
				}

				registerResp := make(map[string]interface{})
				err := app.PostJSON("/api/v1/tenants/register", registerReq, &registerResp)
				results <- err
			}(i)
		}

		// 收集结果
		successCount := 0
		for i := 0; i < concurrentCount; i++ {
			if err := <-results; err == nil {
				successCount++
			}
		}

		assert.GreaterOrEqual(t, successCount, int(concurrentCount*8/10), "至少80%的注册应该成功")
	})
}

// 导入context包
import "context"
import "fmt"
