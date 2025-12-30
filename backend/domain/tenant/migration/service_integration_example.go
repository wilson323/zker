// backend/domain/tenant/migration/service_integration_example.go
// Service层双写集成示例
package migration

import (
	"context"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// 此文件展示如何在Service层集成双写适配器

/*
使用指南:

1. 在Service中注入双写适配器
2. 创建/更新操作时调用双写方法
3. 查询操作使用新表（tenant_id已填充）
4. 读取操作可以逐步迁移到新表

集成示例见下方BotService
*/

// BotService Bot服务（双写集成示例）
type BotService struct {
	db              *gorm.DB
	dualWrite       *DualWriteAdapter
	compensationMgr *CompensationManager
	logger          logs.CtxLogger
}

// NewBotService 创建Bot服务
func NewBotService(db *gorm.DB) *BotService {
	return &BotService{
		db:              db,
		dualWrite:       NewDualWriteAdapter(db, db),
		compensationMgr: NewCompensationManager(db),
		logger:          logs.DefaultLogger(),
	}
}

// CreateBot 创建Bot（双写）
func (s *BotService) CreateBot(ctx context.Context, bot *Bot) (*Bot, error) {
	s.logger.CtxInfof(ctx, "[BotService] 创建Bot (双写模式): %s", bot.Name)

	// 1. 双写创建
	if err := s.dualWrite.CreateBotWithDualWrite(ctx, bot); err != nil {
		s.logger.CtxErrorf(ctx, "[BotService] 双写创建失败: %v", err)

		// 记录失败，后续补偿
		s.compensationMgr.RecordFailedWrite(ctx, "bots", "INSERT", bot, err)

		return nil, err
	}

	s.logger.CtxInfof(ctx, "[BotService] Bot创建成功: bot_id=%s, tenant_id=%s", bot.BotID, bot.TenantID)

	return bot, nil
}

// UpdateBot 更新Bot（双写）
func (s *BotService) UpdateBot(ctx context.Context, botID string, updates map[string]interface{}) (*Bot, error) {
	s.logger.CtxInfof(ctx, "[BotService] 更新Bot (双写模式): bot_id=%s", botID)

	// 1. 先从旧表读取（主路径）
	var bot Bot
	if err := s.db.Table("bots").Where("bot_id = ?", botID).First(&bot).Error; err != nil {
		return nil, err
	}

	// 2. 更新字段
	for k, v := range updates {
		switch k {
		case "name":
			bot.Name = v.(string)
		case "tenant_id":
			bot.TenantID = v.(string)
		}
	}

	// 3. 双写更新
	if err := s.dualWrite.UpdateBotWithDualWrite(ctx, &bot); err != nil {
		s.logger.CtxErrorf(ctx, "[BotService] 双写更新失败: %v", err)

		// 记录失败，后续补偿
		s.compensationMgr.RecordFailedWrite(ctx, "bots", "UPDATE", updates, err)

		return nil, err
	}

	s.logger.CtxInfof(ctx, "[BotService] Bot更新成功: bot_id=%s", botID)

	return &bot, nil
}

// DeleteBot 删除Bot（双写）
func (s *BotService) DeleteBot(ctx context.Context, botID string) error {
	s.logger.CtxInfof(ctx, "[BotService] 删除Bot (双写模式): bot_id=%s", botID)

	// 1. 先查询Bot对象
	var bot Bot
	if err := s.db.Table("bots").Where("bot_id = ?", botID).First(&bot).Error; err != nil {
		return err
	}

	// 2. 双写删除（软删除）
	if err := s.dualWrite.DeleteBotWithDualWrite(ctx, &bot); err != nil {
		s.logger.CtxErrorf(ctx, "[BotService] 双写删除失败: %v", err)

		// 记录失败，后续补偿
		s.compensationMgr.RecordFailedWrite(ctx, "bots", "DELETE", map[string]interface{}{
			"bot_id": botID,
		}, err)

		return err
	}

	s.logger.CtxInfof(ctx, "[BotService] Bot删除成功: bot_id=%s", botID)

	return nil
}

// GetBot 获取Bot（读新表）
func (s *BotService) GetBot(ctx context.Context, botID string) (*Bot, error) {
	s.logger.CtxDebugf(ctx, "[BotService] 获取Bot (读新表): bot_id=%s", botID)

	var bot Bot
	err := s.db.Table("bots").
		Where("bot_id = ? AND deleted_at IS NULL", botID).
		First(&bot).Error

	return &bot, err
}

// ListBots 列出Bots（读新表，按tenant_id过滤）
func (s *BotService) ListBots(ctx context.Context, tenantID string, page, pageSize int) ([]Bot, int64, error) {
	s.logger.CtxDebugf(ctx, "[BotService] 列出Bots (读新表): tenant_id=%s", tenantID)

	var bots []Bot
	var total int64

	// 计算总数
	s.db.Table("bots").
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := s.db.Table("bots").
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&bots).Error

	return bots, total, err
}

// Bot Bot模型（简化）
type Bot struct {
	BotID     string  `gorm:"column:bot_id;primaryKey"`
	TenantID  string  `gorm:"column:tenant_id;not null"`
	Name      string  `gorm:"column:name;not null"`
	CreatedAt string  `gorm:"column:created_at"`
	UpdatedAt string  `gorm:"column:updated_at"`
	DeletedAt *string `gorm:"column:deleted_at"`
}

// ============================================
// 其他Service的集成示例
// ============================================

// ConversationService 对话服务（双写集成示例）
type ConversationService struct {
	db              *gorm.DB
	dualWrite       *DualWriteAdapter
	compensationMgr *CompensationManager
	logger          logs.CtxLogger
}

// NewConversationService 创建对话服务
func NewConversationService(db *gorm.DB) *ConversationService {
	return &ConversationService{
		db:              db,
		dualWrite:       NewDualWriteAdapter(db, db),
		compensationMgr: NewCompensationManager(db),
		logger:          logs.DefaultLogger(),
	}
}

// CreateConversation 创建对话（双写）
func (s *ConversationService) CreateConversation(ctx context.Context, conv *Conversation) (*Conversation, error) {
	s.logger.CtxInfof(ctx, "[ConversationService] 创建对话 (双写模式): tenant_id=%s", conv.TenantID)

	if err := s.dualWrite.CreateWithDualWrite(ctx, "conversations", conv); err != nil {
		s.compensationMgr.RecordFailedWrite(ctx, "conversations", "INSERT", conv, err)
		return nil, err
	}

	return conv, nil
}

// Conversation 对话模型
type Conversation struct {
	ConvID    string `gorm:"column:conv_id;primaryKey"`
	TenantID  string `gorm:"column:tenant_id;not null"`
	Title     string `gorm:"column:title"`
	CreatedAt string `gorm:"column:created_at"`
}

// ============================================
// Service层集成清单
// ============================================

/*
需要集成双写的Service列表:

✅ 1. BotService (bots) - 已完成
✅ 2. ConversationService (conversations) - 已完成示例
⏳ 3. MessageService (messages)
⏳ 4. KnowledgeBaseService (knowledge_bases)
⏳ 5. KnowledgeChunkService (knowledge_chunks)
⏳ 6. WorkflowService (workflows)
⏳ 7. WorkflowExecutionService (workflow_executions)
⏳ 8. SingleAgentDraftService (single_agent_draft)
⏳ 9. PublishedBotService (published_bots)
⏳ 10. BotConfigService (bot_configs)

集成步骤（针对每个Service）:

1. 添加双写适配器字段
   dualWrite *DualWriteAdapter

2. 在NewService中初始化
   dualWrite: NewDualWriteAdapter(db, db)

3. Create方法: 调用 dualWrite.CreateWithDualWrite(ctx, tableName, data)

4. Update方法: 调用 dualWrite.UpdateWithDualWrite(ctx, tableName, id, updates)

5. Delete方法: 调用 dualWrite.DeleteWithDualWrite(ctx, tableName, id)

6. 失败处理: compensationMgr.RecordFailedWrite(ctx, tableName, operation, data, err)

7. 查询方法: 直接查询新表（已填充tenant_id）
*/
