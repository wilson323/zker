// Package saga 提供Bot创建相关的Saga实现
//
// 本包实现了企业级分布式事务系统，使用Saga模式保证Bot创建过程中的最终一致性。
// 业务流程:
//   1. 创建Bot记录
//   2. 创建默认知识库
//   3. 分配Bot权限
//
// 任何步骤失败都会自动触发补偿事务，回滚已执行的步骤。
package saga

import (
	"context"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/infra/saga"
)

// BotCreationCommand Bot创建命令
type BotCreationCommand struct {
	TenantID    string
	Name        string
	Description string
	Avatar      string
	Type        string
	CreatorID   string
}

// Bot Bot实体(简化版)
type Bot struct {
	ID          string
	TenantID    string
	Name        string
	Description string
	Avatar      string
	Type        string
	Status      string
	CreatedBy   string
	CreatedAt   time.Time
}

// KnowledgeBase 知识库实体(简化版)
type KnowledgeBase struct {
	ID        string
	TenantID  string
	BotID     string
	Name      string
	Type      string
	Status    string
	CreatedBy string
}

// Permission 权限实体(简化版)
type Permission struct {
	ResourceID   string
	ResourceType string
	UserID       string
	Role         string
	Actions      []string
}

// createBotStep 创建Bot步骤
type createBotStep struct {
	timeout time.Duration
}

func (s *createBotStep) Name() string {
	return "create_bot"
}

func (s *createBotStep) Timeout() time.Duration {
	return s.timeout
}

func (s *createBotStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	cmd := data.(*BotCreationCommand)
	bot := &Bot{
		ID:          saga.GenerateID(),
		TenantID:    cmd.TenantID,
		Name:        cmd.Name,
		Description: cmd.Description,
		Avatar:      cmd.Avatar,
		Type:        cmd.Type,
		Status:      "draft",
		CreatedBy:   cmd.CreatorID,
		CreatedAt:   time.Now(),
	}
	// 实际场景: botRepo.Save(ctx, bot)
	fmt.Printf("[执行] 创建Bot: %s (ID: %s)\n", bot.Name, bot.ID)
	return bot, nil
}

// createBotCompensation 创建Bot补偿步骤
type createBotCompensation struct {
	timeout time.Duration
}

func (c *createBotCompensation) Name() string {
	return "create_bot_compensation"
}

func (c *createBotCompensation) Timeout() time.Duration {
	return c.timeout
}

func (c *createBotCompensation) Compensate(ctx context.Context, data interface{}) error {
	bot := data.(*Bot)
	// 实际场景: botRepo.Delete(ctx, bot.ID)
	fmt.Printf("[补偿] 删除Bot: %s (ID: %s)\n", bot.Name, bot.ID)
	return nil
}

// createKnowledgeBaseStep 创建知识库步骤
type createKnowledgeBaseStep struct {
	timeout time.Duration
}

func (s *createKnowledgeBaseStep) Name() string {
	return "create_knowledge_base"
}

func (s *createKnowledgeBaseStep) Timeout() time.Duration {
	return s.timeout
}

func (s *createKnowledgeBaseStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	bot := data.(*Bot)
	kb := &KnowledgeBase{
		ID:        saga.GenerateID(),
		TenantID:  bot.TenantID,
		BotID:     bot.ID,
		Name:      bot.Name + "-默认知识库",
		Type:      "private",
		Status:    "active",
		CreatedBy: bot.CreatedBy,
	}
	// 实际场景: kbRepo.Save(ctx, kb)
	fmt.Printf("[执行] 创建知识库: %s (ID: %s)\n", kb.Name, kb.ID)
	return kb, nil
}

// createKnowledgeBaseCompensation 创建知识库补偿步骤
type createKnowledgeBaseCompensation struct {
	timeout time.Duration
}

func (c *createKnowledgeBaseCompensation) Name() string {
	return "create_knowledge_base_compensation"
}

func (c *createKnowledgeBaseCompensation) Timeout() time.Duration {
	return c.timeout
}

func (c *createKnowledgeBaseCompensation) Compensate(ctx context.Context, data interface{}) error {
	kb := data.(*KnowledgeBase)
	// 实际场景: kbRepo.Delete(ctx, kb.ID)
	fmt.Printf("[补偿] 删除知识库: %s (ID: %s)\n", kb.Name, kb.ID)
	return nil
}

// assignPermissionsStep 分配权限步骤
type assignPermissionsStep struct {
	timeout time.Duration
}

func (s *assignPermissionsStep) Name() string {
	return "assign_permissions"
}

func (s *assignPermissionsStep) Timeout() time.Duration {
	return s.timeout
}

func (s *assignPermissionsStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	kb := data.(*KnowledgeBase)
	permissions := []Permission{
		{
			ResourceID:   kb.BotID,
			ResourceType: "bot",
			UserID:       kb.CreatedBy,
			Role:         "owner",
			Actions:      []string{"read", "write", "delete", "share"},
		},
	}
	// 实际场景: permissionRepo.BatchGrant(ctx, permissions)
	fmt.Printf("[执行] 分配Bot权限: %d条 (BotID: %s)\n", len(permissions), kb.BotID)
	return nil, nil
}

// assignPermissionsCompensation 分配权限补偿步骤
type assignPermissionsCompensation struct {
	timeout time.Duration
}

func (c *assignPermissionsCompensation) Name() string {
	return "assign_permissions_compensation"
}

func (c *assignPermissionsCompensation) Timeout() time.Duration {
	return c.timeout
}

func (c *assignPermissionsCompensation) Compensate(ctx context.Context, data interface{}) error {
	kb := data.(*KnowledgeBase)
	// 实际场景: permissionRepo.Revoke(ctx, kb.BotID, kb.CreatedBy)
	fmt.Printf("[补偿] 撤销Bot权限: BotID=%s, UserID=%s\n", kb.BotID, kb.CreatedBy)
	return nil
}

// NewBotCreationSaga 创建Bot创建Saga定义
//
// 返回一个完整的Saga定义，包含3个执行步骤和对应的3个补偿步骤。
// 所有步骤默认超时时间为30秒，使用默认的重试策略。
//
// 使用示例:
//
//	sagaDef := NewBotCreationSaga()
//	orchestrator := saga.NewSagaOrchestrator(repo, logger)
//	if err := orchestrator.DefineSaga(sagaDef); err != nil {
//	    log.Fatal("定义Saga失败", err)
//	}
//
//	cmd := &BotCreationCommand{
//	    TenantID:  "tenant-123",
//	    Name:      "测试Bot",
//	    CreatorID: "user-456",
//	}
//	execution, err := orchestrator.ExecuteSaga(ctx, "bot-creation-saga", cmd)
func NewBotCreationSaga() *saga.Saga {
	stepTimeout := 30 * time.Second

	return &saga.Saga{
		ID:          "saga-bot-creation",
		Name:        "bot-creation-saga",
		Description: "Bot创建流程：创建Bot记录 → 创建默认知识库 → 分配权限",
		Steps: []saga.SagaStep{
			&createBotStep{timeout: stepTimeout},
			&createKnowledgeBaseStep{timeout: stepTimeout},
			&assignPermissionsStep{timeout: stepTimeout},
		},
		Compensations: []saga.CompensationStep{
			&createBotCompensation{timeout: stepTimeout},
			&createKnowledgeBaseCompensation{timeout: stepTimeout},
			&assignPermissionsCompensation{timeout: stepTimeout},
		},
		Timeout:     5 * time.Minute, // 整个Saga超时时间
		RetryPolicy: saga.DefaultRetryPolicy(),
	}
}
