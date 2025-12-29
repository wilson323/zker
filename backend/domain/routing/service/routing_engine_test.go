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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coze-studio/backend/domain/routing/entity"
)

// MockIntentMatcher Mock意图匹配器
type MockIntentMatcher struct {
	mock.Mock
}

func (m *MockIntentMatcher) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MatchOutput), args.Error(1)
}

// MockServiceRegistry Mock服务注册中心
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

// MockLoadMonitor Mock负载监控器
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

// MockRoutingLogRepository Mock路由日志仓储
type MockRoutingLogRepository struct {
	mock.Mock
}

func (m *MockRoutingLogRepository) Create(ctx context.Context, log *entity.RoutingLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockRoutingLogRepository) List(ctx context.Context, filter interface{}) ([]*entity.RoutingLog, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.RoutingLog), args.Get(1).(int64), args.Error(2)
}

// TestScoreBasedRouter_Route_Success 测试路由决策成功
func TestScoreBasedRouter_Route_Success(t *testing.T) {
	// Arrange
	mockMatcher := new(MockIntentMatcher)
	mockRegistry := new(MockServiceRegistry)
	mockLoadMonitor := new(MockLoadMonitor)
	mockLogRepo := new(MockRoutingLogRepository)

	router := NewScoreBasedRouter(mockMatcher, mockRegistry, mockLoadMonitor, mockLogRepo)

	input := &MatchInput{
		UserInput: "我想创建一个聊天机器人",
		TenantID:  "tenant_123",
		Context:   make(map[string]interface{}),
	}

	candidates := []*MatchOutput{
		{
			BotID:      "bot_1",
			Confidence: 0.9,
			MatchType:  "keyword",
			RuleID:     "rule_1",
		},
		{
			BotID:      "bot_2",
			Confidence: 0.7,
			MatchType:  "similarity",
			RuleID:     "rule_2",
		},
	}

	mockMatcher.On("Match", mock.Anything, input).
		Return(candidates, nil)

	mockRegistry.On("GetBotHealth", mock.Anything, "bot_1").
		Return(&ServiceHealth{
			ServiceID:   "bot_1",
			IsHealthy:    true,
			SuccessRate:  0.95,
			AvgLatency:   100 * time.Millisecond,
			CurrentLoad:  5,
			MaxCapacity:  10,
		}, nil)

	mockLoadMonitor.On("GetServiceLoad", mock.Anything, "bot_1").
		Return(&ServiceLoad{
			ServiceID:   "bot_1",
			CurrentLoad: 5,
			MaxCapacity: 10,
		}, nil)

	mockLogRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil)

	// Act
	decision, err := router.Route(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "bot_1", decision.BotID)
	assert.True(t, decision.Score > 0)
	assert.NotEmpty(t, decision.Reasons)
	mockMatcher.AssertExpectations(t)
	mockRegistry.AssertExpectations(t)
	mockLoadMonitor.AssertExpectations(t)
	// logRepo.Create是异步调用的，所以不会被mock框架捕获
}

// TestScoreBasedRouter_Route_NoCandidates 测试无候选结果
func TestScoreBasedRouter_Route_NoCandidates(t *testing.T) {
	// Arrange
	mockMatcher := new(MockIntentMatcher)
	mockRegistry := new(MockServiceRegistry)
	mockLoadMonitor := new(MockLoadMonitor)
	mockLogRepo := new(MockRoutingLogRepository)

	router := NewScoreBasedRouter(mockMatcher, mockRegistry, mockLoadMonitor, mockLogRepo)

	input := &MatchInput{
		UserInput: "未知请求",
		TenantID:  "tenant_123",
		Context:   make(map[string]interface{}),
	}

	mockMatcher.On("Match", mock.Anything, input).
		Return([]*MatchOutput{}, nil)

	// Act
	decision, err := router.Route(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, decision)
	assert.Contains(t, err.Error(), "no matching candidates")
	mockMatcher.AssertExpectations(t)
}

// TestScoreBasedRouter_Route_Workflow 测试工作流路由
func TestScoreBasedRouter_Route_Workflow(t *testing.T) {
	// Arrange
	mockMatcher := new(MockIntentMatcher)
	mockRegistry := new(MockServiceRegistry)
	mockLoadMonitor := new(MockLoadMonitor)
	mockLogRepo := new(MockRoutingLogRepository)

	router := NewScoreBasedRouter(mockMatcher, mockRegistry, mockLoadMonitor, mockLogRepo)

	input := &MatchInput{
		UserInput: "执行工作流",
		TenantID:  "tenant_123",
		Context:   make(map[string]interface{}),
	}

	candidates := []*MatchOutput{
		{
			WorkflowID:  "workflow_1",
			Confidence:  0.85,
			MatchType:   "intent",
			RuleID:      "rule_3",
		},
	}

	mockMatcher.On("Match", mock.Anything, input).
		Return(candidates, nil)

	mockRegistry.On("GetWorkflowHealth", mock.Anything, "workflow_1").
		Return(&ServiceHealth{
			ServiceID:   "workflow_1",
			IsHealthy:    true,
			SuccessRate:  0.9,
			AvgLatency:   200 * time.Millisecond,
			CurrentLoad:  3,
			MaxCapacity:  5,
		}, nil)

	mockLoadMonitor.On("GetServiceLoad", mock.Anything, "workflow_1").
		Return(&ServiceLoad{
			ServiceID:   "workflow_1",
			CurrentLoad: 3,
			MaxCapacity: 5,
		}, nil)

	mockLogRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil)

	// Act
	decision, err := router.Route(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "workflow_1", decision.WorkflowID)
	assert.Empty(t, decision.BotID)
}

// TestScoreBasedRouter_calculateScore_UnhealthyService 测试不健康服务评分
func TestScoreBasedRouter_calculateScore_UnhealthyService(t *testing.T) {
	// Arrange
	mockMatcher := new(MockIntentMatcher)
	mockRegistry := new(MockServiceRegistry)
	mockLoadMonitor := new(MockLoadMonitor)
	mockLogRepo := new(MockRoutingLogRepository)

	router := NewScoreBasedRouter(mockMatcher, mockRegistry, mockLoadMonitor, mockLogRepo)

	input := &MatchInput{
		UserInput: "测试请求",
		TenantID:  "tenant_123",
		Context:   make(map[string]interface{}),
	}

	candidates := []*MatchOutput{
		{
			BotID:      "bot_unhealthy",
			Confidence: 0.9,
			MatchType:  "keyword",
			RuleID:     "rule_1",
		},
	}

	mockMatcher.On("Match", mock.Anything, input).
		Return(candidates, nil)

	mockRegistry.On("GetBotHealth", mock.Anything, "bot_unhealthy").
		Return(&ServiceHealth{
			ServiceID:  "bot_unhealthy",
			IsHealthy:   false,
			SuccessRate: 0.5,
		}, nil)

	mockLoadMonitor.On("GetServiceLoad", mock.Anything, "bot_unhealthy").
		Return(&ServiceLoad{
			ServiceID:   "bot_unhealthy",
			CurrentLoad: 8,
			MaxCapacity: 10,
		}, nil)

	mockLogRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil)

	// Act
	decision, err := router.Route(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	// 不健康的服务得分应该较低
	assert.True(t, decision.Score < 0.5)
}

// TestScoreBasedRouter_SetWeights 测试设置权重
func TestScoreBasedRouter_SetWeights(t *testing.T) {
	// Arrange
	router := NewScoreBasedRouter(nil, nil, nil, nil)

	// Act
	router.SetWeights(0.5, 0.2, 0.1, 0.1, 0.1)

	// Assert
	assert.Equal(t, 0.5, router.intentWeight)
	assert.Equal(t, 0.2, router.healthWeight)
	assert.Equal(t, 0.1, router.loadWeight)
	assert.Equal(t, 0.1, router.costWeight)
	assert.Equal(t, 0.1, router.regionWeight)
}
