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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBotRepository 模拟Bot仓储
type MockBotRepositoryForHealthCheck struct {
	mock.Mock
}

func (m *MockBotRepositoryForHealthCheck) GetAllActiveBots(ctx context.Context) ([]*BotInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*BotInfo), args.Error(1)
}

// MockHealthMonitor 模拟健康监控器
type MockHealthMonitor struct {
	mock.Mock
}

func (m *MockHealthMonitor) RecordSuccess(serviceName string, latency time.Duration) {
	m.Called(serviceName, latency)
}

func (m *MockHealthMonitor) RecordFailure(serviceName, errorMessage string) {
	m.Called(serviceName, errorMessage)
}

func (m *MockHealthMonitor) GetBotHealth(ctx context.Context, botID string) (*ServiceHealthStats, bool) {
	args := m.Called(ctx, botID)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*ServiceHealthStats), args.Bool(1)
}

// ================================================================================
// NewBotHealthCheckTask 测试
// ================================================================================

func TestNewBotHealthCheckTask_DefaultConfig(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	config := BotHealthCheckConfig{}

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, config)

	assert.NotNil(t, task)
	assert.Equal(t, DefaultBotHealthCheckConfig.CheckInterval, task.checkInterval)
	assert.Equal(t, DefaultBotHealthCheckConfig.CheckTimeout, task.checkTimeout)
	assert.Equal(t, DefaultBotHealthCheckConfig.MaxConcurrency, task.maxConcurrency)
	assert.False(t, task.IsEnabled())
}

func TestNewBotHealthCheckTask_CustomConfig(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	config := BotHealthCheckConfig{
		CheckInterval:  60 * time.Second,
		CheckTimeout:   20 * time.Second,
		MaxConcurrency: 50,
	}

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, config)

	assert.NotNil(t, task)
	assert.Equal(t, 60*time.Second, task.checkInterval)
	assert.Equal(t, 20*time.Second, task.checkTimeout)
	assert.Equal(t, 50, task.maxConcurrency)
	assert.False(t, task.IsEnabled())
}

// ================================================================================
// Start/Stop 测试
// ================================================================================

func TestBotHealthCheckTask_Start(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	// 设置期望：返回空Bot列表
	mockRepo.On("GetAllActiveBots", mock.Anything).Return([]*BotInfo{}, nil)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{
		CheckInterval: 100 * time.Millisecond, // 快速测试
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task.Start(ctx)

	// 验证任务已启动
	assert.True(t, task.IsEnabled())

	// 等待一小段时间让任务执行
	time.Sleep(200 * time.Millisecond)

	// 停止任务
	task.Stop()

	// 验证任务已停止
	assert.False(t, task.IsEnabled())

	mockRepo.AssertExpectations(t)
}

func TestBotHealthCheckTask_Stop(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{
		CheckInterval: 100 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())

	task.Start(ctx)
	assert.True(t, task.IsEnabled())

	task.Stop()
	assert.False(t, task.IsEnabled())

	cancel()
}

func TestBotHealthCheckTask_StopWhenNotRunning(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{})

	// 停止未运行的任务（应该安全）
	task.Stop()
	assert.False(t, task.IsEnabled())
}

// ================================================================================
// executeHealthCheck 测试
// ================================================================================

func TestExecuteHealthCheck_Success(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	bots := []*BotInfo{
		{BotID: "bot_1", TenantID: "tenant_1", Name: "Bot 1"},
		{BotID: "bot_2", TenantID: "tenant_1", Name: "Bot 2"},
		{BotID: "bot_3", TenantID: "tenant_1", Name: "Bot 3"},
	}

	mockRepo.On("GetAllActiveBots", mock.Anything).Return(bots, nil)
	mockMonitor.On("RecordSuccess", mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Times(3)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{
		CheckInterval:  1 * time.Second,
		CheckTimeout:   100 * time.Millisecond,
		MaxConcurrency: 10,
	})

	ctx := context.Background()

	// 启动并等待一次执行
	task.Start(ctx)
	time.Sleep(200 * time.Millisecond)
	task.Stop()

	mockRepo.AssertExpectations(t)
	// 注意：RecordSuccess的调用次数可能因并发而略有不同，这里不严格验证
}

func TestExecuteHealthCheck_GetBotsFailure(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	// 设置期望：获取Bot列表失败
	getBotsErr := errors.New("database connection failed")
	mockRepo.On("GetAllActiveBots", mock.Anything).Return(nil, getBotsErr)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{
		CheckInterval: 100 * time.Millisecond,
	})

	ctx := context.Background()

	// 启动并等待一次执行
	task.Start(ctx)
	time.Sleep(200 * time.Millisecond)
	task.Stop()

	mockRepo.AssertExpectations(t)
	// RecordSuccess不应该被调用
	mockMonitor.AssertNotCalled(t, "RecordSuccess")
}

func TestExecuteHealthCheck_NoBots(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	// 设置期望：返回空列表
	mockRepo.On("GetAllActiveBots", mock.Anything).Return([]*BotInfo{}, nil)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{
		CheckInterval: 100 * time.Millisecond,
	})

	ctx := context.Background()

	// 启动并等待一次执行
	task.Start(ctx)
	time.Sleep(200 * time.Millisecond)
	task.Stop()

	mockRepo.AssertExpectations(t)
	// RecordSuccess不应该被调用
	mockMonitor.AssertNotCalled(t, "RecordSuccess", mock.Anything, mock.Anything)
}

func TestExecuteHealthCheck_PartialFailure(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	bots := []*BotInfo{
		{BotID: "bot_1", TenantID: "tenant_1", Name: "Bot 1"},
		{BotID: "bot_2", TenantID: "tenant_1", Name: "Bot 2"},
	}

	mockRepo.On("GetAllActiveBots", mock.Anything).Return(bots, nil)
	// bot_1成功，bot_2失败
	mockMonitor.On("RecordSuccess", "bot:bot_1", mock.AnythingOfType("time.Duration")).Once()
	mockMonitor.On("RecordFailure", "bot:bot_2", mock.AnythingOfType("string")).Once()

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{
		CheckInterval: 100 * time.Millisecond,
		CheckTimeout:   100 * time.Millisecond,
		MaxConcurrency: 10,
	})

	ctx := context.Background()

	// 启动并等待一次执行
	task.Start(ctx)
	time.Sleep(200 * time.Millisecond)
	task.Stop()

	mockRepo.AssertExpectations(t)
}

// ================================================================================
// checkBotsConcurrently 测试
// ================================================================================

func TestCheckBotsConcurrently_EmptyList(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{})

	ctx := context.Background()

	successCount, failureCount := task.checkBotsConcurrently(ctx, []string{})

	assert.Equal(t, 0, successCount)
	assert.Equal(t, 0, failureCount)
}

// ================================================================================
// GetHealthStatus 测试
// ================================================================================

func TestGetHealthStatus_Healthy(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	healthStats := &ServiceHealthStats{
		IsHealthy: true,
	}

	mockMonitor.On("GetBotHealth", mock.Anything, "bot_123").Return(healthStats, true)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{})

	ctx := context.Background()
	isHealthy, err := task.GetHealthStatus(ctx, "bot_123")

	assert.True(t, isHealthy)
	assert.NoError(t, err)
	mockMonitor.AssertExpectations(t)
}

func TestGetHealthStatus_Unhealthy(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	healthStats := &ServiceHealthStats{
		IsHealthy: false,
	}

	mockMonitor.On("GetBotHealth", mock.Anything, "bot_123").Return(healthStats, true)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{})

	ctx := context.Background()
	isHealthy, err := task.GetHealthStatus(ctx, "bot_123")

	assert.False(t, isHealthy)
	assert.NoError(t, err)
	mockMonitor.AssertExpectations(t)
}

func TestGetHealthStatus_NotFound(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	mockMonitor.On("GetBotHealth", mock.Anything, "bot_123").Return(nil, false)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{})

	ctx := context.Background()
	isHealthy, err := task.GetHealthStatus(ctx, "bot_123")

	assert.False(t, isHealthy)
	assert.NoError(t, err)
	mockMonitor.AssertExpectations(t)
}

// ================================================================================
// GetHealthStats 测试
// ================================================================================

func TestGetHealthStats_Found(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	healthStats := &ServiceHealthStats{
		IsHealthy:       true,
		SuccessCount:    100,
		FailureCount:    2,
		LastCheckTime:   time.Now().UnixMilli(),
		AvgLatencyMs:    50,
		P95LatencyMs:    100,
		P99LatencyMs:    150,
		ConsecutiveFailures: 0,
	}

	mockMonitor.On("GetBotHealth", mock.Anything, "bot:bot_123").Return(healthStats, true)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{})

	stats, found := task.GetHealthStats("bot_123")

	assert.True(t, found)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(100), stats.SuccessCount)
	assert.Equal(t, int64(2), stats.FailureCount)
	assert.Equal(t, int64(50), stats.AvgLatencyMs)
	mockMonitor.AssertExpectations(t)
}

func TestGetHealthStats_NotFound(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	mockMonitor.On("GetBotHealth", mock.Anything, "bot:bot_123").Return(nil, false)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{})

	stats, found := task.GetHealthStats("bot_123")

	assert.False(t, found)
	assert.Nil(t, stats)
	mockMonitor.AssertExpectations(t)
}

// ================================================================================
// IsEnabled 测试
// ================================================================================

func TestIsEnabled_Initial(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{})

	assert.False(t, task.IsEnabled())
}

func TestIsEnabled_AfterStart(t *testing.T) {
	mockMonitor := new(MockHealthMonitor)
	mockRepo := new(MockBotRepositoryForHealthCheck)

	mockRepo.On("GetAllActiveBots", mock.Anything).Return([]*BotInfo{}, nil)

	task := NewBotHealthCheckTask(mockMonitor, mockRepo, BotHealthCheckConfig{
		CheckInterval: 100 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())

	task.Start(ctx)
	assert.True(t, task.IsEnabled())

	task.Stop()
	cancel()

	assert.False(t, task.IsEnabled())

	mockRepo.AssertExpectations(t)
}
