// Package saga 提供Bot创建Saga测试
package saga

import (
	"context"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/infra/saga"
)

// TestNewBotCreationSaga 测试Saga定义创建
func TestNewBotCreationSaga(t *testing.T) {
	sagaDef := NewBotCreationSaga()

	// 验证基本属性
	if sagaDef.ID != "saga-bot-creation" {
		t.Errorf("期望 Saga ID 为 'saga-bot-creation', 实际为 '%s'", sagaDef.ID)
	}

	if sagaDef.Name != "bot-creation-saga" {
		t.Errorf("期望 Saga Name 为 'bot-creation-saga', 实际为 '%s'", sagaDef.Name)
	}

	// 验证步骤数量
	if len(sagaDef.Steps) != 3 {
		t.Errorf("期望 3 个步骤, 实际为 %d", len(sagaDef.Steps))
	}

	if len(sagaDef.Compensations) != 3 {
		t.Errorf("期望 3 个补偿步骤, 实际为 %d", len(sagaDef.Compensations))
	}

	// 验证超时设置
	if sagaDef.Timeout != 5*time.Minute {
		t.Errorf("期望超时时间为 5 分钟, 实际为 %v", sagaDef.Timeout)
	}

	// 验证重试策略
	if sagaDef.RetryPolicy == nil {
		t.Error("重试策略不应为空")
	} else {
		if sagaDef.RetryPolicy.MaxAttempts != 3 {
			t.Errorf("期望最大重试次数为 3, 实际为 %d", sagaDef.RetryPolicy.MaxAttempts)
		}
	}

	t.Logf("✓ Saga定义验证通过")
	t.Logf("  - ID: %s", sagaDef.ID)
	t.Logf("  - Name: %s", sagaDef.Name)
	t.Logf("  - Description: %s", sagaDef.Description)
	t.Logf("  - Steps: %d", len(sagaDef.Steps))
	t.Logf("  - Compensations: %d", len(sagaDef.Compensations))
	t.Logf("  - Timeout: %v", sagaDef.Timeout)
}

// TestCreateBotStep 测试创建Bot步骤
func TestCreateBotStep(t *testing.T) {
	step := &createBotStep{timeout: 30 * time.Second}

	// 验证步骤名称
	if step.Name() != "create_bot" {
		t.Errorf("期望步骤名称为 'create_bot', 实际为 '%s'", step.Name())
	}

	// 验证超时时间
	if step.Timeout() != 30*time.Second {
		t.Errorf("期望超时时间为 30秒, 实际为 %v", step.Timeout())
	}

	// 执行步骤
	cmd := &BotCreationCommand{
		TenantID:    "tenant-123",
		Name:        "测试Bot",
		Description: "这是一个测试Bot",
		Avatar:      "avatar.png",
		Type:        "chatbot",
		CreatorID:   "user-456",
	}

	ctx := context.Background()
	result, err := step.Execute(ctx, cmd)

	if err != nil {
		t.Fatalf("执行步骤失败: %v", err)
	}

	bot, ok := result.(*Bot)
	if !ok {
		t.Fatalf("返回结果类型错误, 期望 *Bot, 实际为 %T", result)
	}

	// 验证Bot属性
	if bot.ID == "" {
		t.Error("Bot ID 不应为空")
	}

	if bot.TenantID != cmd.TenantID {
		t.Errorf("期望 TenantID 为 '%s', 实际为 '%s'", cmd.TenantID, bot.TenantID)
	}

	if bot.Name != cmd.Name {
		t.Errorf("期望 Name 为 '%s', 实际为 '%s'", cmd.Name, bot.Name)
	}

	if bot.Status != "draft" {
		t.Errorf("期望 Status 为 'draft', 实际为 '%s'", bot.Status)
	}

	if bot.CreatedBy != cmd.CreatorID {
		t.Errorf("期望 CreatedBy 为 '%s', 实际为 '%s'", cmd.CreatorID, bot.CreatedBy)
	}

	t.Logf("✓ 创建Bot步骤执行成功")
	t.Logf("  - Bot ID: %s", bot.ID)
	t.Logf("  - Bot Name: %s", bot.Name)
	t.Logf("  - Bot Status: %s", bot.Status)
}

// TestCreateBotCompensation 测试创建Bot补偿步骤
func TestCreateBotCompensation(t *testing.T) {
	comp := &createBotCompensation{timeout: 30 * time.Second}

	// 验证补偿步骤名称
	if comp.Name() != "create_bot_compensation" {
		t.Errorf("期望补偿步骤名称为 'create_bot_compensation', 实际为 '%s'", comp.Name())
	}

	// 执行补偿
	bot := &Bot{
		ID:        saga.GenerateID(),
		TenantID:  "tenant-123",
		Name:      "测试Bot",
		Status:    "draft",
		CreatedBy: "user-456",
	}

	ctx := context.Background()
	err := comp.Compensate(ctx, bot)

	if err != nil {
		t.Fatalf("执行补偿失败: %v", err)
	}

	t.Logf("✓ 创建Bot补偿执行成功")
	t.Logf("  - 已删除 Bot ID: %s", bot.ID)
}

// TestCreateKnowledgeBaseStep 测试创建知识库步骤
func TestCreateKnowledgeBaseStep(t *testing.T) {
	step := &createKnowledgeBaseStep{timeout: 30 * time.Second}

	bot := &Bot{
		ID:        saga.GenerateID(),
		TenantID:  "tenant-123",
		Name:      "测试Bot",
		CreatedBy: "user-456",
	}

	ctx := context.Background()
	result, err := step.Execute(ctx, bot)

	if err != nil {
		t.Fatalf("执行步骤失败: %v", err)
	}

	kb, ok := result.(*KnowledgeBase)
	if !ok {
		t.Fatalf("返回结果类型错误, 期望 *KnowledgeBase, 实际为 %T", result)
	}

	// 验证知识库属性
	if kb.ID == "" {
		t.Error("知识库 ID 不应为空")
	}

	if kb.BotID != bot.ID {
		t.Errorf("期望 BotID 为 '%s', 实际为 '%s'", bot.ID, kb.BotID)
	}

	if kb.TenantID != bot.TenantID {
		t.Errorf("期望 TenantID 为 '%s', 实际为 '%s'", bot.TenantID, kb.TenantID)
	}

	if kb.Name != bot.Name+"-默认知识库" {
		t.Errorf("期望 Name 为 '%s', 实际为 '%s'", bot.Name+"-默认知识库", kb.Name)
	}

	t.Logf("✓ 创建知识库步骤执行成功")
	t.Logf("  - 知识库 ID: %s", kb.ID)
	t.Logf("  - 知识库名称: %s", kb.Name)
}

// TestAssignPermissionsStep 测试分配权限步骤
func TestAssignPermissionsStep(t *testing.T) {
	step := &assignPermissionsStep{timeout: 30 * time.Second}

	kb := &KnowledgeBase{
		ID:        saga.GenerateID(),
		BotID:     saga.GenerateID(),
		TenantID:  "tenant-123",
		Name:      "测试Bot-默认知识库",
		CreatedBy: "user-456",
	}

	ctx := context.Background()
	result, err := step.Execute(ctx, kb)

	if err != nil {
		t.Fatalf("执行步骤失败: %v", err)
	}

	if result != nil {
		t.Logf("✓ 分配权限步骤执行成功，返回值: %v", result)
	}
}

// TestGenerateID 测试ID生成
func TestGenerateID(t *testing.T) {
	id1 := saga.GenerateID()
	id2 := saga.GenerateID()

	// 验证ID不为空
	if id1 == "" || id2 == "" {
		t.Error("生成的ID不应为空")
	}

	// 验证ID唯一性
	if id1 == id2 {
		t.Error("两次生成的ID应该不同")
	}

	// 验证ID格式（UUID格式：36字符，包含4个连字符）
	if len(id1) != 36 {
		t.Errorf("UUID长度应为 36, 实际为 %d", len(id1))
	}

	t.Logf("✓ ID生成验证通过")
	t.Logf("  - ID1: %s", id1)
	t.Logf("  - ID2: %s", id2)
}

// BenchmarkGenerateID 性能测试：ID生成
func BenchmarkGenerateID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = saga.GenerateID()
	}
}
