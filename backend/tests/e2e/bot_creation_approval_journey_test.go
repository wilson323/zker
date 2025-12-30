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
	"fmt"
	"net/http"
	"testing"
	"time"

	e2eHelpers "github.com/coze-dev/coze-studio/backend/tests/e2e/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Bot创建+审核旅程测试 ====================

/**
 * TestE2E_BotCreationWithApprovalJourney 测试Bot创建+审核完整旅程
 *
 * 职责: 模拟从Bot创建到审核通过的全流程
 *
 * 流程:
 * 1. 创建Bot(低AI置信度，触发审核)
 * 2. 验证Bot为草稿状态
 * 3. 验证审核任务已创建
 * 4. 审核人接取任务
 * 5. 提交审核结果(批准)
 * 6. 验证Bot已发布
 * 7. 验证审核历史已记录
 *
 * 遵循单一职责原则：只负责Bot创建+审核这一个用户旅程
 */
func TestE2E_BotCreationWithApprovalJourney(t *testing.T) {
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
	tenant := e2eHelpers.CreateTestTenant(t, app, "Bot审核测试公司")
	reviewer := e2eHelpers.CreateTestUser(t, app, tenant.TenantID, "reviewer")

	var botID string
	var taskID string

	t.Run("Step1_创建Bot_低AI置信度", func(t *testing.T) {
		createBotReq := map[string]interface{}{
			"tenant_id":    tenant.TenantID,
			"bot_name":     "客服机器人",
			"ai_confidence": 0.5, // 低置信度，触发审核
		}

		createBotResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/bots", createBotReq, &createBotResp)
		require.NoError(t, err, "创建Bot应该成功")

		// 验证响应
		require.Equal(t, float64(0), createBotResp["code"], "响应码应为0（成功）")
		data, ok := createBotResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		require.NotEmpty(t, data["bot_id"], "Bot ID不应为空")

		// 验证Bot状态为draft（草稿）
		assert.Equal(t, "draft", data["status"], "Bot应该是草稿状态")
		botID = data["bot_id"].(string)
	})

	t.Run("Step2_验证Bot为草稿状态", func(t *testing.T) {
		require.NotEmpty(t, botID, "Bot ID应已设置")

		// 通过API获取Bot详情
		getResp := make(map[string]interface{})
		err := app.GetJSON("/api/v1/bots/"+botID, &getResp)
		require.NoError(t, err, "获取Bot详情应该成功")

		// 验证状态
		data, ok := getResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		assert.Equal(t, "draft", data["status"], "Bot状态应为draft")

		// 验证数据库
		e2eHelpers.AssertBotStatus(t, app.DB, botID, "draft")
	})

	t.Run("Step3_验证审核任务已创建", func(t *testing.T) {
		// 通过API列出审核任务
		tasksResp := make(map[string]interface{})
		err := app.GetJSON("/api/v1/review/tasks", &tasksResp)
		require.NoError(t, err, "列出审核任务应该成功")

		// 验证有待审核任务
		data, ok := tasksResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		tasks, ok := data["tasks"].([]interface{})
		require.True(t, ok, "tasks应为数组")
		assert.Greater(t, len(tasks), 0, "应该有待审核任务")

		// 获取第一个任务ID
		firstTask, ok := tasks[0].(map[string]interface{})
		require.True(t, ok, "任务应为对象")
		taskID = firstTask["task_id"].(string)
		assert.Equal(t, "pending", firstTask["status"], "任务状态应为pending")

		// 验证数据库
		e2eHelpers.AssertReviewTaskExists(t, app.DB, tenant.TenantID, "pending")
	})

	t.Run("Step4_审核人接取任务", func(t *testing.T) {
		require.NotEmpty(t, taskID, "任务ID应已设置")

		assignReq := map[string]interface{}{
			"assigned_to": reviewer.UserID,
		}

		assignResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/review/tasks/"+taskID, assignReq, &assignResp)
		require.NoError(t, err, "分配任务应该成功")
		require.Equal(t, float64(0), assignResp["code"], "响应码应为0")

		// 验证任务状态已更新
		var status string
		err = app.DB.Table("review_tasks").Where("task_id = ?", taskID).Select("status").Scan(&status).Error
		require.NoError(t, err, "查询任务状态失败")
		assert.Equal(t, "assigned", status, "任务状态应为assigned")
	})

	t.Run("Step5_提交审核结果_批准", func(t *testing.T) {
		submitReq := map[string]interface{}{
			"task_id":  taskID,
			"decision": "approved",
			"comment":  "Bot配置合理，可以发布",
		}

		submitResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/review/tasks/submit", submitReq, &submitResp)
		require.NoError(t, err, "提交审核应该成功")
		require.Equal(t, float64(0), submitResp["code"], "响应码应为0")
	})

	t.Run("Step6_验证Bot已发布", func(t *testing.T) {
		// 通过API获取Bot详情
		getResp := make(map[string]interface{})
		err := app.GetJSON("/api/v1/bots/"+botID, &getResp)
		require.NoError(t, err, "获取Bot详情应该成功")

		// 验证状态
		data, ok := getResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		assert.Equal(t, "published", data["status"], "Bot状态应为published")

		// 验证数据库
		e2eHelpers.AssertBotStatus(t, app.DB, botID, "published")
	})

	t.Run("Step7_验证审核历史已记录", func(t *testing.T) {
		// 验证审核任务已完成
		var completedAt interface{}
		err := app.DB.Table("review_tasks").Where("task_id = ?", taskID).Select("completed_at").Scan(&completedAt).Error
		require.NoError(t, err, "查询审核任务失败")
		assert.NotNil(t, completedAt, "完成时间不应为空")

		// 验证任务状态
		var status string
		err = app.DB.Table("review_tasks").Where("task_id = ?", taskID).Select("status").Scan(&status).Error
		require.NoError(t, err, "查询任务状态失败")
		assert.Equal(t, "approved", status, "任务状态应为approved")
	})
}

/**
 * TestE2E_BotCreationHighConfidence 测试高置信度Bot直接发布
 *
 * 职责: 验证高AI置信度的Bot不需要审核直接发布
 */
func TestE2E_BotCreationHighConfidence(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	tenant := e2eHelpers.CreateTestTenant(t, app, "高置信度测试公司")

	t.Run("创建高置信度Bot", func(t *testing.T) {
		createBotReq := map[string]interface{}{
			"tenant_id":     tenant.TenantID,
			"bot_name":      "智能助手",
			"ai_confidence": 0.95, // 高置信度，不需要审核
		}

		createBotResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/bots", createBotReq, &createBotResp)
		require.NoError(t, err, "创建Bot应该成功")

		// 验证Bot直接发布（不需要审核）
		data, ok := createBotResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		assert.Equal(t, "published", data["status"], "高置信度Bot应该直接发布")

		// 验证没有创建审核任务
		var count int64
		app.DB.Table("review_tasks").Where("resource_id = ?", data["bot_id"]).Count(&count)
		assert.Equal(t, int64(0), count, "高置信度Bot不应该创建审核任务")
	})
}

/**
 * TestE2E_BotApprovalRejection 测试Bot审核拒绝
 *
 * 职责: 验证审核拒绝流程
 */
func TestE2E_BotApprovalRejection(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	tenant := e2eHelpers.CreateTestTenant(t, app, "拒绝测试公司")
	reviewer := e2eHelpers.CreateTestUser(t, app, tenant.TenantID, "reviewer")

	var botID string
	var taskID string

	t.Run("创建Bot并拒绝", func(t *testing.T) {
		// 创建Bot
		createBotReq := map[string]interface{}{
			"tenant_id":     tenant.TenantID,
			"bot_name":      "待拒绝Bot",
			"ai_confidence": 0.5,
		}

		createBotResp := make(map[string]interface{})
		err := app.PostJSON("/api/v1/bots", createBotReq, &createBotResp)
		require.NoError(t, err, "创建Bot应该成功")

		data, ok := createBotResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		botID = data["bot_id"].(string)

		// 获取审核任务
		tasksResp := make(map[string]interface{})
		err = app.GetJSON("/api/v1/review/tasks", &tasksResp)
		require.NoError(t, err, "列出审核任务应该成功")

		tasksData, ok := tasksResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		tasks, ok := tasksData["tasks"].([]interface{})
		require.True(t, ok, "tasks应为数组")
		require.Greater(t, len(tasks), 0, "应该有待审核任务")

		firstTask := tasks[0].(map[string]interface{})
		taskID = firstTask["task_id"].(string)

		// 拒绝审核
		submitReq := map[string]interface{}{
			"task_id":  taskID,
			"decision": "rejected",
			"comment":  "Bot配置不符合要求",
		}

		submitResp := make(map[string]interface{})
		err = app.PostJSON("/api/v1/review/tasks/submit", submitReq, &submitResp)
		require.NoError(t, err, "提交审核应该成功")

		// 验证Bot状态保持draft（拒绝后不会发布）
		e2eHelpers.AssertBotStatus(t, app.DB, botID, "draft")

		// 验证任务状态
		var status string
		err = app.DB.Table("review_tasks").Where("task_id = ?", taskID).Select("status").Scan(&status).Error
		require.NoError(t, err, "查询任务状态失败")
		assert.Equal(t, "rejected", status, "任务状态应为rejected")
	})
}

/**
 * TestE2E_BotApprovalConcurrent 测试并发审核
 *
 * 职责: 验证多个Bot并发审核场景
 */
func TestE2E_BotApprovalConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试（使用 -short 标志）")
	}

	app := e2e.SetupTestApplication(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.Terminate(ctx)
	}()

	tenant := e2eHelpers.CreateTestTenant(t, app, "并发审核测试公司")
	reviewer := e2eHelpers.CreateTestUser(t, app, tenant.TenantID, "reviewer")

	t.Run("并发创建多个Bot", func(t *testing.T) {
		concurrentCount := 5
		botIDs := make([]string, concurrentCount)
		taskIDs := make([]string, concurrentCount)

		// 并发创建Bot
		for i := 0; i < concurrentCount; i++ {
			createBotReq := map[string]interface{}{
				"tenant_id":     tenant.TenantID,
				"bot_name":      fmt.Sprintf("并发Bot%d", i),
				"ai_confidence": 0.5,
			}

			createBotResp := make(map[string]interface{})
			err := app.PostJSON("/api/v1/bots", createBotReq, &createBotResp)
			require.NoError(t, err, "创建Bot应该成功")

			data, ok := createBotResp["data"].(map[string]interface{})
			require.True(t, ok, "响应data应为对象")
			botIDs[i] = data["bot_id"].(string)
		}

		// 获取所有审核任务
		tasksResp := make(map[string]interface{})
		err := app.GetJSON("/api/v1/review/tasks", &tasksResp)
		require.NoError(t, err, "列出审核任务应该成功")

		tasksData, ok := tasksResp["data"].(map[string]interface{})
		require.True(t, ok, "响应data应为对象")
		tasks, ok := tasksData["tasks"].([]interface{})
		require.True(t, ok, "tasks应为数组")
		require.GreaterOrEqual(t, len(tasks), concurrentCount, "应该有足够的审核任务")

		// 并发审核
		for i := 0; i < concurrentCount; i++ {
			task := tasks[i].(map[string]interface{})
			taskIDs[i] = task["task_id"].(string)

			assignReq := map[string]interface{}{
				"assigned_to": reviewer.UserID,
			}

			assignResp := make(map[string]interface{})
			err := app.PostJSON("/api/v1/review/tasks/"+taskIDs[i], assignReq, &assignResp)
			require.NoError(t, err, "分配任务应该成功")

			submitReq := map[string]interface{}{
				"task_id":  taskIDs[i],
				"decision": "approved",
				"comment":  fmt.Sprintf("审核通过Bot%d", i),
			}

			submitResp := make(map[string]interface{})
			err = app.PostJSON("/api/v1/review/tasks/submit", submitReq, &submitResp)
			require.NoError(t, err, "提交审核应该成功")
		}

		// 验证所有Bot都已发布
		for _, botID := range botIDs {
			e2eHelpers.AssertBotStatus(t, app.DB, botID, "published")
		}
	})
}
