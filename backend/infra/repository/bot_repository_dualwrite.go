// backend/infra/repository/bot_repository_dualwrite.go
package repository

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/agent/singleagent/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// BotRepositoryWithDualWrite 支持双写的Bot Repository
// 用于迁移期间同时读写新旧字段
type BotRepositoryWithDualWrite struct {
	db              *gorm.DB
	dualWriteEnabled bool // 双写开关（从配置文件读取）
}

// NewBotRepositoryWithDualWrite 创建支持双写的Bot Repository
func NewBotRepositoryWithDualWrite(
	db *gorm.DB,
	dualWriteEnabled bool,
) *BotRepositoryWithDualWrite {
	return &BotRepositoryWithDualWrite{
		db:              db,
		dualWriteEnabled: dualWriteEnabled,
	}
}

// Create 创建Bot（支持双写）
func (r *BotRepositoryWithDualWrite) Create(
	ctx context.Context,
	bot *entity.SingleAgent,
) error {
	// 创建记录
	err := r.db.WithContext(ctx).Create(bot).Error
	if err != nil {
		return fmt.Errorf("创建Bot失败: %w", err)
	}

	logs.CtxInfof(ctx, "[BotRepo] bot created: id=%d, space_id=%d, dual_write=%v",
		bot.AgentID, bot.SpaceID, r.dualWriteEnabled)

	return nil
}

// Update 更新Bot（支持双写）
func (r *BotRepositoryWithDualWrite) Update(
	ctx context.Context,
	bot *entity.SingleAgent,
) error {
	// 更新记录
	err := r.db.WithContext(ctx).Save(bot).Error
	if err != nil {
		return fmt.Errorf("更新Bot失败: %w", err)
	}

	logs.CtxInfof(ctx, "[BotRepo] bot updated: id=%d, space_id=%d, dual_write=%v",
		bot.AgentID, bot.SpaceID, r.dualWriteEnabled)

	return nil
}

// GetByID 根据ID查询Bot（支持懒加载）
func (r *BotRepositoryWithDualWrite) GetByID(
	ctx context.Context,
	id int64,
) (*entity.SingleAgent, error) {
	var bot entity.SingleAgent
	err := r.db.WithContext(ctx).Where("agent_id = ?", id).First(&bot).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Bot不存在")
		}
		return nil, fmt.Errorf("查询Bot失败: %w", err)
	}

	return &bot, nil
}

// ListBySpaceID 根据space_id查询Bot列表
func (r *BotRepositoryWithDualWrite) ListBySpaceID(
	ctx context.Context,
	spaceID int64,
	pageNum, pageSize int,
) ([]*entity.SingleAgent, int64, error) {
	var bots []*entity.SingleAgent
	var total int64

	// 查询总数
	err := r.db.WithContext(ctx).Model(&entity.SingleAgent{}).
		Where("space_id = ?", spaceID).
		Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询Bot总数失败: %w", err)
	}

	// 分页查询
	offset := (pageNum - 1) * pageSize
	err = r.db.WithContext(ctx).Where("space_id = ?", spaceID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&bots).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询Bot列表失败: %w", err)
	}

	return bots, total, nil
}

// Delete 删除Bot
func (r *BotRepositoryWithDualWrite) Delete(
	ctx context.Context,
	id int64,
) error {
	err := r.db.WithContext(ctx).Where("agent_id = ?", id).Delete(&entity.SingleAgent{}).Error
	if err != nil {
		return fmt.Errorf("删除Bot失败: %w", err)
	}

	logs.CtxInfof(ctx, "[BotRepo] bot deleted: id=%d", id)

	return nil
}

// GetAllActiveBots 获取所有活跃的Bot（未删除）
func (r *BotRepositoryWithDualWrite) GetAllActiveBots(
	ctx context.Context,
) ([]*entity.SingleAgent, error) {
	var bots []*entity.SingleAgent

	// 查询所有未删除的Bot（deleted_at IS NULL）
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&bots).Error
	if err != nil {
		return nil, fmt.Errorf("查询活跃Bot列表失败: %w", err)
	}

	logs.CtxInfof(ctx, "[BotRepo] found %d active bots", len(bots))

	return bots, nil
}
