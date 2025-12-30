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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/infra/saga"
	"github.com/coze-dev/coze-studio/backend/tests/integration/fixtures"
)

// SagaBusinessIntegrationSuite Saga+业务系统集成测试套件
// 职责: 测试Saga编排引擎与业务系统的集成（Bot创建、知识库创建等）
type SagaBusinessIntegrationSuite struct {
	suite.Suite
	ctx           context.Context
	db            *sql.DB
	logger        *zap.Logger
	orchestrator  saga.Orchestrator
	botRepo       *fixtures.MockBotRepository
	knowledgeRepo *fixtures.MockKnowledgeRepository
}

// SetupSuite 测试套件初始化
func (s *SagaBusinessIntegrationSuite) SetupSuite() {
	s.ctx = context.Background()
	s.logger = zap.NewNop()

	// 创建MySQL测试容器
	testCtx, err := SetupMySQLTestContainer(s.T())
	assert.NoError(s.T(), err)
	s.db = testCtx.DB

	// 初始化Saga协调器
	sagaRepo := saga.NewRepository(s.db)
	s.orchestrator = saga.NewSagaOrchestrator(sagaRepo, s.logger)

	// 初始化Mock Repository
	s.botRepo = fixtures.NewMockBotRepository()
	s.knowledgeRepo = fixtures.NewMockKnowledgeRepository()
}

// TearDownSuite 测试套件清理
func (s *SagaBusinessIntegrationSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

// SetupTest 每个测试前的设置
func (s *SagaBusinessIntegrationSuite) SetupTest() {
	// 清理Mock数据
	s.botRepo.Clear()
	s.knowledgeRepo.Clear()
}

// TestSagaBotCreationFlow 测试Bot创建Saga完整流程
// 职责: 验证Saga编排引擎在Bot创建场景下的完整执行
func (s *SagaBusinessIntegrationSuite) TestSagaBotCreationFlow() {
	t := s.T()

	// Setup: 定义Bot创建Saga
	botID := fmt.Sprintf("bot-%d", time.Now().UnixNano())
	knowledgeID := fmt.Sprintf("knowledge-%d", time.Now().UnixNano())

	botCreationSaga := &saga.Saga{
		ID:          fmt.Sprintf("saga-bot-%d", time.Now().UnixNano()),
		Name:        "bot_creation",
		Description: "Bot创建Saga",
		Timeout:     30 * time.Second,
		RetryPolicy: saga.DefaultRetryPolicy(),
		Steps: []saga.SagaStep{
			&fixtures.CreateBotStep{
				BotRepo:    s.botRepo,
				BotID:      botID,
				BotName:    "测试Bot",
				BotDesc:    "这是一个测试Bot",
				TenantID:   "tenant-001",
				CreatorID:  "user-001",
			},
			&fixtures.CreateKnowledgeStep{
				KnowledgeRepo: s.knowledgeRepo,
				KnowledgeID:   knowledgeID,
				KnowledgeName: "测试知识库",
				BotID:         botID,
				TenantID:      "tenant-001",
			},
		},
		Compensations: []saga.CompensationStep{
			&fixtures.DeleteBotStep{
				BotRepo: s.botRepo,
			},
			&fixtures.DeleteKnowledgeStep{
				KnowledgeRepo: s.knowledgeRepo,
			},
		},
	}

	// Execute: 注册并执行Saga
	err := s.orchestrator.DefineSaga(botCreationSaga)
	assert.NoError(t, err, "定义Saga应该成功")

	execution, err := s.orchestrator.ExecuteSaga(s.ctx, "bot_creation", nil)
	assert.NoError(t, err, "执行Saga应该成功")
	assert.NotNil(t, execution)

	// Verify: 验证Saga执行状态
	assert.Equal(t, saga.SagaStatusCompleted, execution.Status, "Saga应该完成")
	assert.NotNil(t, execution.CompletedAt, "应该有完成时间")
	assert.Equal(t, 2, len(execution.StepExecutions), "应该执行2个步骤")

	// Verify: 验证所有步骤都成功
	for _, stepExec := range execution.StepExecutions {
		assert.Equal(t, saga.StepStatusCompleted, stepExec.Status, "步骤应该成功")
	}

	// Verify: 验证Bot已创建
	bot, err := s.botRepo.GetByID(s.ctx, botID)
	assert.NoError(t, err)
	assert.NotNil(t, bot, "Bot应该被创建")
	assert.Equal(t, botID, bot.BotID)
	assert.Equal(t, "测试Bot", bot.Name)
	assert.Equal(t, "tenant-001", bot.TenantID)

	// Verify: 验证知识库已创建
	knowledge, err := s.knowledgeRepo.GetByID(s.ctx, knowledgeID)
	assert.NoError(t, err)
	assert.NotNil(t, knowledge, "知识库应该被创建")
	assert.Equal(t, knowledgeID, knowledge.KnowledgeID)
	assert.Equal(t, "测试知识库", knowledge.Name)
	assert.Equal(t, botID, knowledge.BotID)
}

// TestSagaBotCreationCompensation 测试Bot创建Saga失败时的补偿
// 职责: 验证Saga步骤失败时自动补偿机制
func (s *SagaBusinessIntegrationSuite) TestSagaBotCreationCompensation() {
	t := s.T()

	// Setup: 创建会失败的知识库Repository
	failingKnowledgeRepo := fixtures.NewFailingMockKnowledgeRepository()

	botID := fmt.Sprintf("bot-%d", time.Now().UnixNano())
	knowledgeID := fmt.Sprintf("knowledge-%d", time.Now().UnixNano())

	// Setup: 定义Bot创建Saga（第二步会失败）
	botCreationSaga := &saga.Saga{
		ID:          fmt.Sprintf("saga-bot-fail-%d", time.Now().UnixNano()),
		Name:        "bot_creation_with_failure",
		Description: "Bot创建Saga（第二步失败）",
		Timeout:     30 * time.Second,
		RetryPolicy: saga.DefaultRetryPolicy(),
		Steps: []saga.SagaStep{
			&fixtures.CreateBotStep{
				BotRepo:    s.botRepo,
				BotID:      botID,
				BotName:    "测试Bot",
				BotDesc:    "这个Bot会被补偿删除",
				TenantID:   "tenant-001",
				CreatorID:  "user-001",
			},
			&fixtures.CreateKnowledgeStep{
				KnowledgeRepo: failingKnowledgeRepo, // 这个会失败
				KnowledgeID:   knowledgeID,
				KnowledgeName: "会失败的知识库",
				BotID:         botID,
				TenantID:      "tenant-001",
			},
		},
		Compensations: []saga.CompensationStep{
			&fixtures.DeleteBotStep{
				BotRepo: s.botRepo,
			},
			&fixtures.DeleteKnowledgeStep{
				KnowledgeRepo: failingKnowledgeRepo,
			},
		},
	}

	// Execute: 注册并执行Saga
	err := s.orchestrator.DefineSaga(botCreationSaga)
	assert.NoError(t, err)

	execution, err := s.orchestrator.ExecuteSaga(s.ctx, "bot_creation_with_failure", nil)

	// Verify: 验证Saga执行失败
	assert.Error(t, err, "执行Saga应该失败")
	assert.NotNil(t, execution)
	assert.Equal(t, saga.SagaStatusCompensated, execution.Status, "Saga应该被补偿")

	// Verify: 验证第一步成功，第二步失败
	assert.Equal(t, 2, len(execution.StepExecutions))
	assert.Equal(t, saga.StepStatusCompleted, execution.StepExecutions[0].Status, "第一步应该成功")
	assert.Equal(t, saga.StepStatusFailed, execution.StepExecutions[1].Status, "第二步应该失败")

	// Verify: 验证Bot已被删除（补偿）
	bot, err := s.botRepo.GetByID(s.ctx, botID)
	assert.NoError(t, err)
	assert.Nil(t, bot, "Bot应该被补偿删除")

	// Verify: 验证知识库未创建（因为失败了）
	knowledge, err := failingKnowledgeRepo.GetByID(s.ctx, knowledgeID)
	assert.NoError(t, err)
	assert.Nil(t, knowledge, "知识库不应该被创建")
}

// TestSagaRetryMechanism 测试Saga重试机制
// 职责: 验证Saga步骤失败时的重试机制
func (s *SagaBusinessIntegrationSuite) TestSagaRetryMechanism() {
	t := s.T()

	// Setup: 创建会重试成功的Repository
	retryKnowledgeRepo := fixtures.NewRetryMockKnowledgeRepository(2) // 第2次成功

	botID := fmt.Sprintf("bot-%d", time.Now().UnixNano())
	knowledgeID := fmt.Sprintf("knowledge-%d", time.Now().UnixNano())

	// Setup: 定义Bot创建Saga
	botCreationSaga := &saga.Saga{
		ID:          fmt.Sprintf("saga-bot-retry-%d", time.Now().UnixNano()),
		Name:        "bot_creation_with_retry",
		Description: "Bot创建Saga（带重试）",
		Timeout:     30 * time.Second,
		RetryPolicy: &saga.RetryPolicy{
			MaxAttempts:     3,
			InitialInterval: 100 * time.Millisecond,
			MaxInterval:     1 * time.Second,
			Multiplier:      2.0,
		},
		Steps: []saga.SagaStep{
			&fixtures.CreateBotStep{
				BotRepo:    s.botRepo,
				BotID:      botID,
				BotName:    "测试Bot",
				BotDesc:    "这个Bot会成功",
				TenantID:   "tenant-001",
				CreatorID:  "user-001",
			},
			&fixtures.CreateKnowledgeStep{
				KnowledgeRepo: retryKnowledgeRepo, // 第2次会成功
				KnowledgeID:   knowledgeID,
				KnowledgeName: "重试成功知识库",
				BotID:         botID,
				TenantID:      "tenant-001",
			},
		},
		Compensations: []saga.CompensationStep{
			&fixtures.DeleteBotStep{
				BotRepo: s.botRepo,
			},
			&fixtures.DeleteKnowledgeStep{
				KnowledgeRepo: retryKnowledgeRepo,
			},
		},
	}

	// Execute: 注册并执行Saga
	err := s.orchestrator.DefineSaga(botCreationSaga)
	assert.NoError(t, err)

	execution, err := s.orchestrator.ExecuteSaga(s.ctx, "bot_creation_with_retry", nil)

	// Verify: 验证Saga执行成功
	assert.NoError(t, err, "重试后Saga应该成功")
	assert.NotNil(t, execution)
	assert.Equal(t, saga.SagaStatusCompleted, execution.Status, "Saga应该完成")

	// Verify: 验证Bot和知识库都已创建
	bot, err := s.botRepo.GetByID(s.ctx, botID)
	assert.NoError(t, err)
	assert.NotNil(t, bot)

	knowledge, err := retryKnowledgeRepo.GetByID(s.ctx, knowledgeID)
	assert.NoError(t, err)
	assert.NotNil(t, knowledge)

	// Verify: 验证重试次数
	assert.Equal(t, 2, retryKnowledgeRepo.GetCreateCallCount(knowledgeID), "应该重试2次")
}

// TestSagaTimeout 测试Saga超时机制
// 职责: 验证Saga步骤超时时的处理
func (s *SagaBusinessIntegrationSuite) TestSagaTimeout() {
	t := s.T()

	// Setup: 创建超时的Repository
	timeoutKnowledgeRepo := fixtures.NewTimeoutMockKnowledgeRepository()

	botID := fmt.Sprintf("bot-%d", time.Now().UnixNano())
	knowledgeID := fmt.Sprintf("knowledge-%d", time.Now().UnixNano())

	// Setup: 定义Bot创建Saga（知识库创建会超时）
	botCreationSaga := &saga.Saga{
		ID:          fmt.Sprintf("saga-bot-timeout-%d", time.Now().UnixNano()),
		Name:        "bot_creation_with_timeout",
		Description: "Bot创建Saga（超时）",
		Timeout:     30 * time.Second,
		RetryPolicy: saga.DefaultRetryPolicy(),
		Steps: []saga.SagaStep{
			&fixtures.CreateBotStep{
				BotRepo:    s.botRepo,
				BotID:      botID,
				BotName:    "测试Bot",
				BotDesc:    "这个Bot会被补偿",
				TenantID:   "tenant-001",
				CreatorID:  "user-001",
			},
			&fixtures.CreateKnowledgeStep{
				KnowledgeRepo: timeoutKnowledgeRepo, // 会超时
				KnowledgeID:   knowledgeID,
				KnowledgeName: "超时知识库",
				BotID:         botID,
				TenantID:      "tenant-001",
			},
		},
		Compensations: []saga.CompensationStep{
			&fixtures.DeleteBotStep{
				BotRepo: s.botRepo,
			},
			&fixtures.DeleteKnowledgeStep{
				KnowledgeRepo: timeoutKnowledgeRepo,
			},
		},
	}

	// Execute: 注册并执行Saga
	err := s.orchestrator.DefineSaga(botCreationSaga)
	assert.NoError(t, err)

	execution, err := s.orchestrator.ExecuteSaga(s.ctx, "bot_creation_with_timeout", nil)

	// Verify: 验证Saga超时失败
	assert.Error(t, err, "Saga应该超时失败")
	assert.NotNil(t, execution)
	assert.Equal(t, saga.SagaStatusCompensated, execution.Status, "Saga应该被补偿")

	// Verify: 验证Bot已被删除
	bot, err := s.botRepo.GetByID(s.ctx, botID)
	assert.NoError(t, err)
	assert.Nil(t, bot, "Bot应该被补偿删除")
}

// TestSagaParallelExecution 测试Saga并行执行（如果支持）
// 职责: 验证Saga步骤并行执行的能力
func (s *SagaBusinessIntegrationSuite) TestSagaParallelExecution() {
	t := s.T()

	// Setup: 创建多个独立的Bot
	botIDs := []string{
		fmt.Sprintf("bot-1-%d", time.Now().UnixNano()),
		fmt.Sprintf("bot-2-%d", time.Now().UnixNano()),
		fmt.Sprintf("bot-3-%d", time.Now().UnixNano()),
	}

	// Setup: 定义并行创建多个Bot的Saga
	parallelBotCreationSaga := &saga.Saga{
		ID:          fmt.Sprintf("saga-parallel-bot-%d", time.Now().UnixNano()),
		Name:        "parallel_bot_creation",
		Description: "并行创建多个Bot",
		Timeout:     30 * time.Second,
		RetryPolicy: saga.DefaultRetryPolicy(),
		Steps: []saga.SagaStep{
			&fixtures.CreateBotStep{
				BotRepo:   s.botRepo,
				BotID:     botIDs[0],
				BotName:   "并行Bot1",
				BotDesc:   "第一个并行Bot",
				TenantID:  "tenant-001",
				CreatorID: "user-001",
			},
			&fixtures.CreateBotStep{
				BotRepo:   s.botRepo,
				BotID:     botIDs[1],
				BotName:   "并行Bot2",
				BotDesc:   "第二个并行Bot",
				TenantID:  "tenant-001",
				CreatorID: "user-001",
			},
			&fixtures.CreateBotStep{
				BotRepo:   s.botRepo,
				BotID:     botIDs[2],
				BotName:   "并行Bot3",
				BotDesc:   "第三个并行Bot",
				TenantID:  "tenant-001",
				CreatorID: "user-001",
			},
		},
		Compensations: []saga.CompensationStep{
			&fixtures.DeleteBotStep{BotRepo: s.botRepo},
			&fixtures.DeleteBotStep{BotRepo: s.botRepo},
			&fixtures.DeleteBotStep{BotRepo: s.botRepo},
		},
	}

	// Execute: 注册并执行Saga
	err := s.orchestrator.DefineSaga(parallelBotCreationSaga)
	assert.NoError(t, err)

	startTime := time.Now()
	execution, err := s.orchestrator.ExecuteSaga(s.ctx, "parallel_bot_creation", nil)
	duration := time.Since(startTime)

	// Verify: 验证Saga执行成功
	assert.NoError(t, err, "并行Saga应该成功")
	assert.Equal(t, saga.SagaStatusCompleted, execution.Status, "Saga应该完成")

	// Verify: 验证所有Bot都已创建
	for _, botID := range botIDs {
		bot, err := s.botRepo.GetByID(s.ctx, botID)
		assert.NoError(t, err)
		assert.NotNil(t, bot, "Bot应该被创建")
	}

	// Note: 并行执行的性能验证
	// 如果串行执行，3个Bot创建耗时约为 3 * 单个Bot创建时间
	// 如果并行执行，耗时应该约为单个Bot创建时间
	t.Logf("并行执行耗时: %v", duration)
}

// TestSagaStatusQuery 测试Saga状态查询
// 职责: 验证Saga执行状态的查询功能
func (s *SagaBusinessIntegrationSuite) TestSagaStatusQuery() {
	t := s.T()

	botID := fmt.Sprintf("bot-%d", time.Now().UnixNano())

	// Setup: 定义Bot创建Saga
	botCreationSaga := &saga.Saga{
		ID:          fmt.Sprintf("saga-bot-status-%d", time.Now().UnixNano()),
		Name:        "bot_creation_for_status",
		Description: "测试状态查询",
		Timeout:     30 * time.Second,
		RetryPolicy: saga.DefaultRetryPolicy(),
		Steps: []saga.SagaStep{
			&fixtures.CreateBotStep{
				BotRepo:   s.botRepo,
				BotID:     botID,
				BotName:   "测试Bot",
				BotDesc:   "用于状态查询",
				TenantID:  "tenant-001",
				CreatorID: "user-001",
			},
		},
		Compensations: []saga.CompensationStep{
			&fixtures.DeleteBotStep{BotRepo: s.botRepo},
		},
	}

	// Execute: 注册并执行Saga
	err := s.orchestrator.DefineSaga(botCreationSaga)
	assert.NoError(t, err)

	execution, err := s.orchestrator.ExecuteSaga(s.ctx, "bot_creation_for_status", nil)
	assert.NoError(t, err)
	assert.NotNil(t, execution)

	// Execute: 查询执行状态
	status, err := s.orchestrator.GetStatus(s.ctx, execution.ID)
	assert.NoError(t, err, "查询状态应该成功")
	assert.NotNil(t, status)

	// Verify: 验证状态信息
	assert.Equal(t, execution.ID, status.ID, "执行ID应该一致")
	assert.Equal(t, execution.SagaID, status.SagaID, "Saga ID应该一致")
	assert.Equal(t, saga.SagaStatusCompleted, status.Status, "状态应该是已完成")
	assert.Equal(t, 1, len(status.StepExecutions), "应该有1个步骤执行记录")

	// Verify: 验证步骤执行详情
	stepExec := status.StepExecutions[0]
	assert.Equal(t, "CreateBot", stepExec.StepName)
	assert.Equal(t, saga.StepStatusCompleted, stepExec.Status)
	assert.NotNil(t, stepExec.CompletedAt, "应该有完成时间")
}

// TestSagaValidation 测试Saga定义验证
// 职责: 验证无效的Saga定义会被拒绝
func (s *SagaBusinessIntegrationSuite) TestSagaValidation() {
	t := s.T()

	tests := []struct {
		name        string
		saga        *saga.Saga
		expectError string
	}{
		{
			name: "空ID",
			saga: &saga.Saga{
				ID:    "",
				Name:  "invalid_saga",
				Steps: []saga.SagaStep{},
			},
			expectError: "Saga ID不能为空",
		},
		{
			name: "空名称",
			saga: &saga.Saga{
				ID:    "saga-001",
				Name:  "",
				Steps: []saga.SagaStep{},
			},
			expectError: "Saga名称不能为空",
		},
		{
			name: "步骤与补偿数量不匹配",
			saga: &saga.Saga{
				ID:    "saga-002",
				Name:  "mismatch_saga",
				Steps: []saga.SagaStep{
					&fixtures.CreateBotStep{},
				},
				Compensations: []saga.CompensationStep{},
			},
			expectError: "步骤数量与补偿数量不匹配",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.orchestrator.DefineSaga(tt.saga)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}
}

// RunTestSuite 运行测试套件
func TestSagaBusinessIntegrationSuite(t *testing.T) {
	suite.Run(t, new(SagaBusinessIntegrationSuite))
}
