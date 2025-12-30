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
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/internal/dal"
)

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	// 使用 testcontainers 模拟 MySQL 8.4.5
	// 注意：实际使用时需要安装 testcontainers-go
	dsn := "testuser:testpass@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(
		&entity.CollaborationTask{},
		&entity.CollaborationHistory{},
		&entity.CollaborationConfig{},
	)
	require.NoError(t, err)

	return db
}

// TestOrchestrator_CreateTask 测试创建任务
func TestOrchestrator_CreateTask(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 创建依赖
	taskRepo := dal.NewCollaborationRepository(db)
	historyRepo := dal.NewHistoryRepository(db)
	configRepo := dal.NewConfigRepository(db)
	assigner := NewTaskAssigner()
	notifier := NewNotificationService()
	logger := zap.NewNop()

	orchestrator := NewCollaborationOrchestrator(
		taskRepo,
		historyRepo,
		configRepo,
		assigner,
		notifier,
		logger,
	)

	ctx := context.Background()

	// 创建配置
	config := &entity.CollaborationConfig{
		TenantID:             "tenant-001",
		AutoReviewThreshold:   0.70,
		AutoEscalateThreshold: 0.30,
		DefaultSLAMinutes:     60,
		UrgentSLAMinutes:      15,
		NotificationEnabled:   true,
	}
	err := configRepo.Create(ctx, config)
	require.NoError(t, err)

	// 测试数据
	contextData := map[string]interface{}{
		"bot_id":   "bot-001",
		"bot_name": "Test Bot",
		"content":  "This is a test content",
	}

	req := &CreateTaskRequest{
		TenantID:      "tenant-001",
		TaskType:      entity.TaskTypeReview,
		Source:        entity.SourceAI,
		Priority:      entity.PriorityMedium,
		Context:       contextData,
		TriggerReason: "Low confidence",
	}

	// 执行
	task, err := orchestrator.CreateTask(ctx, req)

	// 断言
	require.NoError(t, err)
	assert.NotNil(t, task)
	assert.NotEmpty(t, task.TaskID)
	assert.Equal(t, "tenant-001", task.TenantID)
	assert.Equal(t, entity.TaskTypeReview, task.TaskType)
	assert.Equal(t, entity.PriorityMedium, task.Priority)
	assert.Equal(t, entity.TaskStatusPending, task.Status)
	assert.True(t, task.SLADeadline.After(time.Now()))
}

// TestOrchestrator_AssignTask 测试分配任务
func TestOrchestrator_AssignTask(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 创建依赖
	taskRepo := dal.NewCollaborationRepository(db)
	historyRepo := dal.NewHistoryRepository(db)
	configRepo := dal.NewConfigRepository(db)
	assigner := NewTaskAssigner()
	notifier := NewNotificationService()
	logger := zap.NewNop()

	orchestrator := NewCollaborationOrchestrator(
		taskRepo,
		historyRepo,
		configRepo,
		assigner,
		notifier,
		logger,
	)

	ctx := context.Background()

	// 1. 先创建任务
	config := &entity.CollaborationConfig{
		TenantID:           "tenant-001",
		DefaultSLAMinutes:  60,
		UrgentSLAMinutes:   15,
	}
	err := configRepo.Create(ctx, config)
	require.NoError(t, err)

	task, err := orchestrator.CreateTask(ctx, &CreateTaskRequest{
		TenantID: "tenant-001",
		TaskType: entity.TaskTypeReview,
		Source:   entity.SourceAI,
		Priority: entity.PriorityMedium,
		Context:  map[string]interface{}{"test": "data"},
	})
	require.NoError(t, err)

	// 2. 分配任务
	err = orchestrator.AssignTask(ctx, task.TaskID, "user-001", "测试分配")
	require.NoError(t, err)

	// 3. 验证
	updatedTask, err := orchestrator.GetTask(ctx, task.TaskID)
	require.NoError(t, err)
	assert.Equal(t, entity.TaskStatusAssigned, updatedTask.Status)
	assert.NotNil(t, updatedTask.AssignedTo)
	assert.Equal(t, "user-001", *updatedTask.AssignedTo)
}

// TestOrchestrator_SubmitResult 测试提交结果
func TestOrchestrator_SubmitResult(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 创建依赖
	taskRepo := dal.NewCollaborationRepository(db)
	historyRepo := dal.NewHistoryRepository(db)
	configRepo := dal.NewConfigRepository(db)
	assigner := NewTaskAssigner()
	notifier := NewNotificationService()
	logger := zap.NewNop()

	orchestrator := NewCollaborationOrchestrator(
		taskRepo,
		historyRepo,
		configRepo,
		assigner,
		notifier,
		logger,
	)

	ctx := context.Background()

	// 1. 创建并分配任务
	config := &entity.CollaborationConfig{
		TenantID:          "tenant-001",
		DefaultSLAMinutes: 60,
		UrgentSLAMinutes:  15,
	}
	err := configRepo.Create(ctx, config)
	require.NoError(t, err)

	task, err := orchestrator.CreateTask(ctx, &CreateTaskRequest{
		TenantID: "tenant-001",
		TaskType: entity.TaskTypeReview,
		Source:   entity.SourceAI,
		Priority: entity.PriorityMedium,
		Context:  map[string]interface{}{"test": "data"},
	})
	require.NoError(t, err)

	err = orchestrator.AssignTask(ctx, task.TaskID, "user-001", "测试分配")
	require.NoError(t, err)

	// 2. 提交审核结果
	result := &SubmitResultRequest{
		ReviewerID: "user-001",
		Decision:   "approved",
		HumanDecision: map[string]interface{}{
			"approved":    true,
			"comment":     "符合要求",
			"suggestion":  "无需修改",
		},
		DecisionReason: "内容准确",
	}

	err = orchestrator.SubmitResult(ctx, task.TaskID, result)
	require.NoError(t, err)

	// 3. 验证
	completedTask, err := orchestrator.GetTask(ctx, task.TaskID)
	require.NoError(t, err)
	assert.Equal(t, entity.TaskStatusCompleted, completedTask.Status)
	assert.NotNil(t, completedTask.Result)
}

// TestOrchestrator_EscalateTask 测试升级任务
func TestOrchestrator_EscalateTask(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 创建依赖
	taskRepo := dal.NewCollaborationRepository(db)
	historyRepo := dal.NewHistoryRepository(db)
	configRepo := dal.NewConfigRepository(db)
	assigner := NewTaskAssigner()
	notifier := NewNotificationService()
	logger := zap.NewNop()

	orchestrator := NewCollaborationOrchestrator(
		taskRepo,
		historyRepo,
		configRepo,
		assigner,
		notifier,
		logger,
	)

	ctx := context.Background()

	// 1. 创建任务
	config := &entity.CollaborationConfig{
		TenantID:          "tenant-001",
		DefaultSLAMinutes: 60,
		UrgentSLAMinutes:  15,
	}
	err := configRepo.Create(ctx, config)
	require.NoError(t, err)

	task, err := orchestrator.CreateTask(ctx, &CreateTaskRequest{
		TenantID: "tenant-001",
		TaskType: entity.TaskTypeReview,
		Source:   entity.SourceAI,
		Priority: entity.PriorityMedium,
		Context:  map[string]interface{}{"test": "data"},
	})
	require.NoError(t, err)

	// 2. 升级任务
	err = orchestrator.EscalateTask(ctx, task.TaskID, "admin-001", "需要高级审核")
	require.NoError(t, err)

	// 3. 验证
	escalatedTask, err := orchestrator.GetTask(ctx, task.TaskID)
	require.NoError(t, err)
	assert.Equal(t, entity.TaskStatusEscalated, escalatedTask.Status)
	assert.Equal(t, "admin-001", *escalatedTask.AssignedTo)
}

// TestOrchestrator_MonitorSLA 测试SLA监控
func TestOrchestrator_MonitorSLA(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 创建依赖
	taskRepo := dal.NewCollaborationRepository(db)
	historyRepo := dal.NewHistoryRepository(db)
	configRepo := dal.NewConfigRepository(db)
	assigner := NewTaskAssigner()
	notifier := NewNotificationService()
	logger := zap.NewNop()

	orchestrator := NewCollaborationOrchestrator(
		taskRepo,
		historyRepo,
		configRepo,
		assigner,
		notifier,
		logger,
	)

	ctx := context.Background()

	// 1. 创建一个已超时的任务
	config := &entity.CollaborationConfig{
		TenantID:          "tenant-001",
		DefaultSLAMinutes: 60,
		UrgentSLAMinutes:  15,
	}
	err := configRepo.Create(ctx, config)
	require.NoError(t, err)

	task, err := orchestrator.CreateTask(ctx, &CreateTaskRequest{
		TenantID:   "tenant-001",
		TaskType:   entity.TaskTypeReview,
		Source:     entity.SourceAI,
		Priority:   entity.PriorityMedium,
		Context:    map[string]interface{}{"test": "data"},
		SLAMinutes: 1, // 1分钟SLA
	})
	require.NoError(t, err)

	// 等待超时
	time.Sleep(2 * time.Second)

	// 2. 监控SLA
	overdueTasks, err := orchestrator.MonitorSLA(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(overdueTasks), 1)
}

// BenchmarkOrchestrator_CreateTask 性能测试
func BenchmarkOrchestrator_CreateTask(b *testing.B) {
	db := setupTestDB(&testing.T{})
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	taskRepo := dal.NewCollaborationRepository(db)
	historyRepo := dal.NewHistoryRepository(db)
	configRepo := dal.NewConfigRepository(db)
	assigner := NewTaskAssigner()
	notifier := NewNotificationService()
	logger := zap.NewNop()

	orchestrator := NewCollaborationOrchestrator(
		taskRepo,
		historyRepo,
		configRepo,
		assigner,
		notifier,
		logger,
	)

	ctx := context.Background()

	config := &entity.CollaborationConfig{
		TenantID:          "tenant-001",
		DefaultSLAMinutes: 60,
		UrgentSLAMinutes:  15,
	}
	configRepo.Create(ctx, config)

	contextData := map[string]interface{}{"test": "data"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := orchestrator.CreateTask(ctx, &CreateTaskRequest{
			TenantID: "tenant-001",
			TaskType: entity.TaskTypeReview,
			Source:   entity.SourceAI,
			Priority: entity.PriorityMedium,
			Context:  contextData,
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestQueueService_SortByPriority 测试队列排序
func TestQueueService_SortByPriority(t *testing.T) {
	// 创建测试任务
	tasks := []*entity.CollaborationTask{
		{
			TaskID:      "task-001",
			Priority:    entity.PriorityLow,
			SLADeadline: time.Now().Add(2 * time.Hour),
		},
		{
			TaskID:      "task-002",
			Priority:    entity.PriorityUrgent,
			SLADeadline: time.Now().Add(30 * time.Minute),
		},
		{
			TaskID:      "task-003",
			Priority:    entity.PriorityHigh,
			SLADeadline: time.Now().Add(1 * time.Hour),
		},
	}

	// 执行排序
	queueSvc := NewQueueService(nil, nil, nil)
	sortedTasks := queueSvc.SortByPriority(tasks)

	// 断言：紧急任务应该排在前面
	assert.Equal(t, entity.PriorityUrgent, sortedTasks[0].Priority)
	assert.Equal(t, entity.PriorityHigh, sortedTasks[1].Priority)
	assert.Equal(t, entity.PriorityLow, sortedTasks[2].Priority)
}

// TestProtocolService_CalculateMetrics 测试性能指标计算
func TestProtocolService_CalculateMetrics(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	taskRepo := dal.NewCollaborationRepository(db)
	historyRepo := dal.NewHistoryRepository(db)

	protocol := NewProtocolService(taskRepo, historyRepo)

	ctx := context.Background()

	// 创建并完成任务
	task := &entity.CollaborationTask{
		TaskID:      "task-001",
		TenantID:    "tenant-001",
		TaskType:    entity.TaskTypeReview,
		Source:      entity.SourceAI,
		Priority:    entity.PriorityMedium,
		Status:      entity.TaskStatusPending,
		Context:     entity.JSON([]byte(`{"test": "data"}`)),
		SLADeadline: time.Now().Add(1 * time.Hour),
		CreatedAt:   time.Now().Unix()*1000 - 3600000, // 1小时前创建
	}
	err := taskRepo.Create(ctx, task)
	require.NoError(t, err)

	// 记录历史
	histories := []*entity.CollaborationHistory{
		{
			TaskID:   "task-001",
			TenantID: "tenant-001",
			ActorID:  "system",
			Action:   entity.ActionCreated,
			Details:  entity.JSON([]byte(`{}`)),
			Timestamp: task.CreatedAt,
		},
		{
			TaskID:   "task-001",
			TenantID: "tenant-001",
			ActorID:  "user-001",
			Action:   entity.ActionAssigned,
			Details:  entity.JSON([]byte(`{}`)),
			Timestamp: task.CreatedAt + 60000, // 1分钟后
		},
		{
			TaskID:   "task-001",
			TenantID: "tenant-001",
			ActorID:  "user-001",
			Action:   entity.ActionCompleted,
			Details:  entity.JSON([]byte(`{}`)),
			Timestamp: task.CreatedAt + 1800000, // 30分钟后
		},
	}
	for _, h := range histories {
		err = historyRepo.Create(ctx, h)
		require.NoError(t, err)
	}

	// 计算指标
	metrics, err := protocol.CalculateMetrics(ctx, "task-001")
	require.NoError(t, err)

	// 断言
	assert.Equal(t, "task-001", metrics.TaskID)
	assert.Equal(t, 3, metrics.ActionsCount)
	assert.Greater(t, metrics.WaitTime, float64(0))
	assert.Greater(t, metrics.HandleTime, float64(0))
}
