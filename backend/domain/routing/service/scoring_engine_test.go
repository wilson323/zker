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
	"testing"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockLoadMonitor 模拟负载监控器
type MockLoadMonitor struct {
	mock.Mock
}

func (m *MockLoadMonitor) GetServiceLoad(ctx context.Context, serviceID string) (*ServiceLoad, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ServiceLoad), args.Error(1)
}

// MockServiceRegistry 模拟服务注册中心
type MockServiceRegistry struct {
	mock.Mock
}

func (m *MockServiceRegistry) GetBotHealth(ctx context.Context, botID string) (*ServiceHealth, error) {
	args := m.Called(ctx, botID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ServiceHealth), args.Error(1)
}

func (m *MockServiceRegistry) GetWorkflowHealth(ctx context.Context, workflowID string) (*ServiceHealth, error) {
	args := m.Called(ctx, workflowID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ServiceHealth), args.Error(1)
}

// MockRoutingOptimizeLogRepository 模拟路由优化日志仓储
type MockRoutingOptimizeLogRepository struct {
	mock.Mock
}

func (m *MockRoutingOptimizeLogRepository) Create(ctx context.Context, log *entity.RoutingOptimizeLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockRoutingOptimizeLogRepository) CreateBatch(ctx context.Context, logs []*entity.RoutingOptimizeLog) error {
	args := m.Called(ctx, logs)
	return args.Error(0)
}

func (m *MockRoutingOptimizeLogRepository) GetByTimeRange(ctx context.Context, tenantID string, startTime, endTime int64) ([]*entity.RoutingOptimizeLog, error) {
	args := m.Called(ctx, tenantID, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.RoutingOptimizeLog), args.Error(1)
}

func (m *MockRoutingOptimizeLogRepository) GetByIntent(ctx context.Context, tenantID, intent string, limit int) ([]*entity.RoutingOptimizeLog, error) {
	args := m.Called(ctx, tenantID, intent, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.RoutingOptimizeLog), args.Error(1)
}

func (m *MockRoutingOptimizeLogRepository) GetByAgent(ctx context.Context, tenantID, agentID string, limit int) ([]*entity.RoutingOptimizeLog, error) {
	args := m.Called(ctx, tenantID, agentID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.RoutingOptimizeLog), args.Error(1)
}

func (m *MockRoutingOptimizeLogRepository) DeleteOldLogs(ctx context.Context, beforeDate int64) error {
	args := m.Called(ctx, beforeDate)
	return args.Error(0)
}

// TestScoringEngine_ScoreAgent 测试Agent评分
func TestScoringEngine_ScoreAgent(t *testing.T) {
	// 创建模拟对象
	mockLoadMonitor := new(MockLoadMonitor)
	mockServiceRegistry := new(MockServiceRegistry)
	mockOptLogRepo := new(MockRoutingOptimizeLogRepository)
	logger := zap.NewNop()

	// 创建评分引擎
	engine := NewScoringEngine(mockLoadMonitor, mockServiceRegistry, mockOptLogRepo, logger)

	// 准备测试数据
	ctx := context.Background()
	agentID := "test-agent-123"
	routingCtx := &RoutingContext{
		Request: &RoutingRequest{
			TenantID: "tenant-123",
			UserID:   "user-123",
		},
		UserHistory: []string{agentID}, // 用户有使用历史
	}

	// 设置模拟返回值
	mockLoadMonitor.On("GetServiceLoad", ctx, agentID).Return(&ServiceLoad{
		ServiceID:   agentID,
		CurrentLoad: 30,
		MaxCapacity: 100,
	}, nil)

	mockServiceRegistry.On("GetBotHealth", ctx, agentID).Return(&ServiceHealth{
		ServiceID:  agentID,
		IsHealthy:  true,
		SuccessRate: 0.95,
	}, nil)

	mockOptLogRepo.On("GetByAgent", ctx, mock.Anything, agentID, 100).Return([]*entity.RoutingOptimizeLog{
		{
			LogID:        "log-1",
			AgentID:      agentID,
			ResponseTime: 1500,
			UserRating:   4,
		},
		{
			LogID:        "log-2",
			AgentID:      agentID,
			ResponseTime: 2000,
			UserRating:   5,
		},
	}, nil)

	// 执行评分
	score, err := engine.ScoreAgent(ctx, agentID, routingCtx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, score)
	assert.Equal(t, agentID, score.AgentID)
	assert.Greater(t, score.TotalScore, 0.0)
	assert.LessOrEqual(t, score.TotalScore, 1.0)
	assert.NotEmpty(t, score.Reasons)

	// 验证模拟对象被正确调用
	mockLoadMonitor.AssertExpectations(t)
	mockServiceRegistry.AssertExpectations(t)
	mockOptLogRepo.AssertExpectations(t)
}

// TestScoringEngine_ScoreAgentWithError 测试Agent评分错误场景
func TestScoringEngine_ScoreAgentWithError(t *testing.T) {
	// 创建模拟对象
	mockLoadMonitor := new(MockLoadMonitor)
	mockServiceRegistry := new(MockServiceRegistry)
	mockOptLogRepo := new(MockRoutingOptimizeLogRepository)
	logger := zap.NewNop()

	engine := NewScoringEngine(mockLoadMonitor, mockServiceRegistry, mockOptLogRepo, logger)

	ctx := context.Background()
	agentID := "test-agent-123"
	routingCtx := &RoutingContext{
		Request: &RoutingRequest{
			TenantID: "tenant-123",
			UserID:   "user-123",
		},
	}

	// 设置模拟返回错误
	mockLoadMonitor.On("GetServiceLoad", ctx, agentID).Return(nil, assert.AnError)
	mockServiceRegistry.On("GetBotHealth", ctx, agentID).Return(nil, assert.AnError)
	mockOptLogRepo.On("GetByAgent", ctx, mock.Anything, agentID, 100).Return(nil, assert.AnError)

	// 执行评分
	score, err := engine.ScoreAgent(ctx, agentID, routingCtx)

	// 验证结果（错误场景应返回默认分数）
	assert.NoError(t, err) // 即使失败，也应该返回默认评分
	assert.NotNil(t, score)
	assert.Greater(t, score.TotalScore, 0.0)
}

// TestScoringEngine_BatchScoreAgents 测试批量评分
func TestScoringEngine_BatchScoreAgents(t *testing.T) {
	mockLoadMonitor := new(MockLoadMonitor)
	mockServiceRegistry := new(MockServiceRegistry)
	mockOptLogRepo := new(MockRoutingOptimizeLogRepository)
	logger := zap.NewNop()

	engine := NewScoringEngine(mockLoadMonitor, mockServiceRegistry, mockOptLogRepo, logger)

	ctx := context.Background()
	agentIDs := []string{"agent-1", "agent-2", "agent-3"}
	routingCtx := &RoutingContext{
		Request: &RoutingRequest{
			TenantID: "tenant-123",
		},
	}

	// 设置模拟返回值
	for _, agentID := range agentIDs {
		mockLoadMonitor.On("GetServiceLoad", ctx, agentID).Return(&ServiceLoad{
			ServiceID:   agentID,
			CurrentLoad: 50,
			MaxCapacity: 100,
		}, nil)

		mockServiceRegistry.On("GetBotHealth", ctx, agentID).Return(&ServiceHealth{
			ServiceID:  agentID,
			IsHealthy:  true,
			SuccessRate: 0.9,
		}, nil)

		mockOptLogRepo.On("GetByAgent", ctx, mock.Anything, agentID, 100).Return([]*entity.RoutingOptimizeLog{
			{
				LogID:        "log-1",
				AgentID:      agentID,
				ResponseTime: 1000,
				UserRating:   4,
			},
		}, nil)
	}

	// 执行批量评分
	scores, err := engine.BatchScoreAgents(ctx, agentIDs, routingCtx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, scores)
	assert.Len(t, scores, 3)

	for _, score := range scores {
		assert.NotEmpty(t, score.AgentID)
		assert.Greater(t, score.TotalScore, 0.0)
	}
}

// TestScoringEngine_scoreLoad 测试负载评分
func TestScoringEngine_scoreLoad(t *testing.T) {
	tests := []struct {
		name           string
		currentLoad    int
		maxCapacity    int
		expectedScore  float64
	}{
		{
			name:          "低负载",
			currentLoad:   10,
			maxCapacity:   100,
			expectedScore: 0.9,
		},
		{
			name:          "中等负载",
			currentLoad:   50,
			maxCapacity:   100,
			expectedScore: 0.5,
		},
		{
			name:          "高负载",
			currentLoad:   90,
			maxCapacity:   100,
			expectedScore: 0.1,
		},
		{
			name:          "满负载",
			currentLoad:   100,
			maxCapacity:   100,
			expectedScore: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLoadMonitor := new(MockLoadMonitor)
			mockServiceRegistry := new(MockServiceRegistry)
			mockOptLogRepo := new(MockRoutingOptimizeLogRepository)
			logger := zap.NewNop()

			engine := NewScoringEngine(mockLoadMonitor, mockServiceRegistry, mockOptLogRepo, logger)

			ctx := context.Background()
			agentID := "test-agent"

			mockLoadMonitor.On("GetServiceLoad", ctx, agentID).Return(&ServiceLoad{
				ServiceID:   agentID,
				CurrentLoad: tt.currentLoad,
				MaxCapacity: tt.maxCapacity,
			}, nil)

			score, err := engine.scoreLoad(ctx, agentID)

			assert.NoError(t, err)
			assert.InDelta(t, tt.expectedScore, score, 0.01)
		})
	}
}

// TestScoringEngine_scoreHealth 测试健康度评分
func TestScoringEngine_scoreHealth(t *testing.T) {
	tests := []struct {
		name          string
		isHealthy     bool
		successRate   float64
		expectedScore float64
	}{
		{
			name:          "健康且高成功率",
			isHealthy:     true,
			successRate:   0.95,
			expectedScore: 0.95,
		},
		{
			name:          "健康但低成功率",
			isHealthy:     true,
			successRate:   0.7,
			expectedScore: 0.7,
		},
		{
			name:          "不健康",
			isHealthy:     false,
			successRate:   0.5,
			expectedScore: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLoadMonitor := new(MockLoadMonitor)
			mockServiceRegistry := new(MockServiceRegistry)
			mockOptLogRepo := new(MockRoutingOptimizeLogRepository)
			logger := zap.NewNop()

			engine := NewScoringEngine(mockLoadMonitor, mockServiceRegistry, mockOptLogRepo, logger)

			ctx := context.Background()
			agentID := "test-agent"

			mockServiceRegistry.On("GetBotHealth", ctx, agentID).Return(&ServiceHealth{
				ServiceID:   agentID,
				IsHealthy:   tt.isHealthy,
				SuccessRate: tt.successRate,
			}, nil)

			score, err := engine.scoreHealth(ctx, agentID)

			assert.NoError(t, err)
			assert.InDelta(t, tt.expectedScore, score, 0.01)
		})
	}
}
