// Package repository 提供仓储接口的持久化实现
//
// 本文件展示Bot仓储的MySQL实现，包括：
// - CRUD操作的SQL执行
// - 复杂查询的构建
// - 分页和排序
// - 事务支持
// - 数据库实体的映射
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/coze-studio/domain/entity"
	"github.com/coze-studio/domain/repository"
)

// botRepositoryImpl Bot仓储实现
type botRepositoryImpl struct {
	db *sql.DB
}

// NewBotRepository 创建Bot仓储实例
func NewBotRepository(db *sql.DB) repository.BotRepository {
	return &botRepositoryImpl{db: db}
}

// ============ CRUD 操作实现 ============

// Save 保存Bot（创建或更新）
func (r *botRepositoryImpl) Save(ctx context.Context, bot *entity.Bot) error {
	// 判断是创建还是更新
	if bot.IsNew() {
		return r.create(ctx, bot)
	}
	return r.update(ctx, bot)
}

// create 创建新Bot
func (r *botRepositoryImpl) create(ctx context.Context, bot *entity.Bot) error {
	query := `
		INSERT INTO bots (
			id, tenant_id, name, description, avatar,
			system_prompt, knowledge_base_id, llm_provider, llm_model,
			llm_config, tags, status, published, version, public,
			creator_id, created_at, updated_at, published_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	llmConfigJSON, _ := bot.LLMConfigJSON()
	tagsJSON, _ := bot.TagsJSON()

	_, err := r.db.ExecContext(
		ctx,
		query,
		bot.ID().String(),
		bot.TenantID().String(),
		bot.Name(),
		bot.Description(),
		bot.Avatar(),
		bot.SystemPrompt(),
		bot.KnowledgeBaseID(),
		bot.LLMProvider(),
		bot.LLMModel(),
		llmConfigJSON,
		tagsJSON,
		string(bot.Status()),
		bot.IsPublished(),
		bot.Version(),
		bot.IsPublic(),
		bot.CreatorID(),
		bot.CreatedAt(),
		bot.UpdatedAt(),
		bot.PublishedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert bot: %w", err)
	}

	return nil
}

// update 更新Bot
func (r *botRepositoryImpl) update(ctx context.Context, bot *entity.Bot) error {
	query := `
		UPDATE bots SET
			name = ?,
			description = ?,
			avatar = ?,
			system_prompt = ?,
			knowledge_base_id = ?,
			llm_provider = ?,
			llm_model = ?,
			llm_config = ?,
			tags = ?,
			status = ?,
			published = ?,
			version = ?,
			public = ?,
			updated_at = ?,
			published_at = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	llmConfigJSON, _ := bot.LLMConfigJSON()
	tagsJSON, _ := bot.TagsJSON()

	result, err := r.db.ExecContext(
		ctx,
		query,
		bot.Name(),
		bot.Description(),
		bot.Avatar(),
		bot.SystemPrompt(),
		bot.KnowledgeBaseID(),
		bot.LLMProvider(),
		bot.LLMModel(),
		llmConfigJSON,
		tagsJSON,
		string(bot.Status()),
		bot.IsPublished(),
		bot.Version(),
		bot.IsPublic(),
		bot.UpdatedAt(),
		bot.PublishedAt(),
		bot.ID().String(),
		bot.TenantID().String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update bot: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FindByID 根据ID查找Bot
func (r *botRepositoryImpl) FindByID(
	ctx context.Context,
	tenantID entity.TenantID,
	botID string,
) (*entity.Bot, error) {
	query := `
		SELECT
			id, tenant_id, name, description, avatar,
			system_prompt, knowledge_base_id, llm_provider, llm_model,
			llm_config, tags, status, published, version, public,
			creator_id, created_at, updated_at, published_at
		FROM bots
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, botID, tenantID.String())

	bot, err := r.scanBot(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bot not found")
		}
		return nil, fmt.Errorf("failed to scan bot: %w", err)
	}

	return bot, nil
}

// Delete 删除Bot（软删除）
func (r *botRepositoryImpl) Delete(
	ctx context.Context,
	tenantID entity.TenantID,
	botID string,
) error {
	query := `
		UPDATE bots SET
			deleted_at = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), botID, tenantID.String())
	if err != nil {
		return fmt.Errorf("failed to delete bot: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ============ 批量操作实现 ============

// FindByTenantID 查找租户下的所有Bot
func (r *botRepositoryImpl) FindByTenantID(
	ctx context.Context,
	tenantID entity.TenantID,
) ([]*entity.Bot, error) {
	query := `
		SELECT
			id, tenant_id, name, description, avatar,
			system_prompt, knowledge_base_id, llm_provider, llm_model,
			llm_config, tags, status, published, version, public,
			creator_id, created_at, updated_at, published_at
		FROM bots
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query bots: %w", err)
	}
	defer rows.Close()

	bots := make([]*entity.Bot, 0)
	for rows.Next() {
		bot, err := r.scanBot(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bot: %w", err)
		}
		bots = append(bots, bot)
	}

	return bots, nil
}

// FindByTenantIDWithPaging 分页查找租户下的Bot
func (r *botRepositoryImpl) FindByTenantIDWithPaging(
	ctx context.Context,
	tenantID entity.TenantID,
	req *repository.ListBotsRequest,
) (*repository.BotListResult, error) {
	// 构建查询条件
	whereConditions := []string{
		"tenant_id = ?",
		"deleted_at IS NULL",
	}
	args := []interface{}{tenantID.String()}

	// 应用过滤条件
	if req.Name != "" {
		whereConditions = append(whereConditions, "name LIKE ?")
		args = append(args, "%"+req.Name+"%")
	}

	if req.Status != nil {
		whereConditions = append(whereConditions, "status = ?")
		args = append(args, string(*req.Status))
	}

	if req.Published != nil {
		whereConditions = append(whereConditions, "published = ?")
		args = append(args, *req.Published)
	}

	if req.Tag != "" {
		whereConditions = append(whereConditions, "JSON_CONTAINS(tags, ?)")
		args = append(args, fmt.Sprintf(`"%s"`, req.Tag))
	}

	if req.CreatorID != "" {
		whereConditions = append(whereConditions, "creator_id = ?")
		args = append(args, req.CreatorID)
	}

	if req.StartTime != nil {
		whereConditions = append(whereConditions, "created_at >= ?")
		args = append(args, *req.StartTime)
	}

	if req.EndTime != nil {
		whereConditions = append(whereConditions, "created_at <= ?")
		args = append(args, *req.EndTime)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM bots WHERE %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count bots: %w", err)
	}

	// 查询列表
	offset := req.GetOffset()
	listQuery := fmt.Sprintf(`
		SELECT
			id, tenant_id, name, description, avatar,
			system_prompt, knowledge_base_id, llm_provider, llm_model,
			llm_config, tags, status, published, version, public,
			creator_id, created_at, updated_at, published_at
		FROM bots
		WHERE %s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereClause, req.SortBy, req.SortOrder)

	args = append(args, req.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bots: %w", err)
	}
	defer rows.Close()

	bots := make([]*entity.Bot, 0)
	for rows.Next() {
		bot, err := r.scanBot(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bot: %w", err)
		}
		bots = append(bots, bot)
	}

	totalPages := int64(0)
	if req.PageSize > 0 {
		totalPages = (total + int64(req.PageSize) - 1) / int64(req.PageSize)
	}

	result := &repository.BotListResult{
		Bots:       bots,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	return result, nil
}

// FindByIDs 根据ID列表批量查找Bot
func (r *botRepositoryImpl) FindByIDs(
	ctx context.Context,
	tenantID entity.TenantID,
	botIDs []string,
) ([]*entity.Bot, error) {
	if len(botIDs) == 0 {
		return []*entity.Bot{}, nil
	}

	placeholders := strings.Repeat("?,", len(botIDs))
	placeholders = placeholders[:len(placeholders)-1] // 移除最后一个逗号

	query := fmt.Sprintf(`
		SELECT
			id, tenant_id, name, description, avatar,
			system_prompt, knowledge_base_id, llm_provider, llm_model,
			llm_config, tags, status, published, version, public,
			creator_id, created_at, updated_at, published_at
		FROM bots
		WHERE tenant_id = ? AND id IN (%s) AND deleted_at IS NULL
	`, placeholders)

	args := make([]interface{}, 0, len(botIDs)+1)
	args = append(args, tenantID.String())
	for _, id := range botIDs {
		args = append(args, id)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bots: %w", err)
	}
	defer rows.Close()

	bots := make([]*entity.Bot, 0)
	for rows.Next() {
		bot, err := r.scanBot(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bot: %w", err)
		}
		bots = append(bots, bot)
	}

	return bots, nil
}

// ============ 查询方法实现 ============

// FindByName 根据名称查找Bot
func (r *botRepositoryImpl) FindByName(
	ctx context.Context,
	tenantID entity.TenantID,
	name string,
) (*entity.Bot, error) {
	query := `
		SELECT
			id, tenant_id, name, description, avatar,
			system_prompt, knowledge_base_id, llm_provider, llm_model,
			llm_config, tags, status, published, version, public,
			creator_id, created_at, updated_at, published_at
		FROM bots
		WHERE tenant_id = ? AND name = ? AND deleted_at IS NULL
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, tenantID.String(), name)

	bot, err := r.scanBot(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan bot: %w", err)
	}

	return bot, nil
}

// FindPublished 查找已发布的Bot
func (r *botRepositoryImpl) FindPublished(
	ctx context.Context,
	tenantID entity.TenantID,
	published bool,
) ([]*entity.Bot, error) {
	query := `
		SELECT
			id, tenant_id, name, description, avatar,
			system_prompt, knowledge_base_id, llm_provider, llm_model,
			llm_config, tags, status, published, version, public,
			creator_id, created_at, updated_at, published_at
		FROM bots
		WHERE tenant_id = ? AND published = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID.String(), published)
	if err != nil {
		return nil, fmt.Errorf("failed to query bots: %w", err)
	}
	defer rows.Close()

	bots := make([]*entity.Bot, 0)
	for rows.Next() {
		bot, err := r.scanBot(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bot: %w", err)
		}
		bots = append(bots, bot)
	}

	return bots, nil
}

// FindByTag 根据标签查找Bot
func (r *botRepositoryImpl) FindByTag(
	ctx context.Context,
	tenantID entity.TenantID,
	tag string,
) ([]*entity.Bot, error) {
	query := `
		SELECT
			id, tenant_id, name, description, avatar,
			system_prompt, knowledge_base_id, llm_provider, llm_model,
			llm_config, tags, status, published, version, public,
			creator_id, created_at, updated_at, published_at
		FROM bots
		WHERE tenant_id = ? AND JSON_CONTAINS(tags, ?) AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID.String(), fmt.Sprintf(`"%s"`, tag))
	if err != nil {
		return nil, fmt.Errorf("failed to query bots: %w", err)
	}
	defer rows.Close()

	bots := make([]*entity.Bot, 0)
	for rows.Next() {
		bot, err := r.scanBot(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bot: %w", err)
		}
		bots = append(bots, bot)
	}

	return bots, nil
}

// FindByCreator 根据创建者查找Bot
func (r *botRepositoryImpl) FindByCreator(
	ctx context.Context,
	tenantID entity.TenantID,
	creatorID string,
) ([]*entity.Bot, error) {
	query := `
		SELECT
			id, tenant_id, name, description, avatar,
			system_prompt, knowledge_base_id, llm_provider, llm_model,
			llm_config, tags, status, published, version, public,
			creator_id, created_at, updated_at, published_at
		FROM bots
		WHERE tenant_id = ? AND creator_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID.String(), creatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query bots: %w", err)
	}
	defer rows.Close()

	bots := make([]*entity.Bot, 0)
	for rows.Next() {
		bot, err := r.scanBot(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bot: %w", err)
		}
		bots = append(bots, bot)
	}

	return bots, nil
}

// ============ 统计方法实现 ============

// CountByTenantID 统计租户下的Bot数量
func (r *botRepositoryImpl) CountByTenantID(
	ctx context.Context,
	tenantID entity.TenantID,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM bots
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	var count int64
	if err := r.db.QueryRowContext(ctx, query, tenantID.String()).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count bots: %w", err)
	}

	return count, nil
}

// CountByStatus 根据状态统计Bot数量
func (r *botRepositoryImpl) CountByStatus(
	ctx context.Context,
	tenantID entity.TenantID,
	status entity.BotStatus,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM bots
		WHERE tenant_id = ? AND status = ? AND deleted_at IS NULL
	`

	var count int64
	if err := r.db.QueryRowContext(ctx, query, tenantID.String(), string(status)).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count bots: %w", err)
	}

	return count, nil
}

// ============ 事务支持 ============

// WithTx 在事务中执行操作
func (r *botRepositoryImpl) WithTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 确保事务被回滚或提交
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 执行业务逻辑
	if err := fn(ctx); err != nil {
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ============ 私有辅助方法 ============

// scanBot 扫描Bot行数据
func (r *botRepositoryImpl) scanBot(
	scanner interface {
		Scan(dest ...interface{}) error
	},
) (*entity.Bot, error) {
	var (
		id             string
		tenantID       string
		name           string
		description    sql.NullString
		avatar         sql.NullString
		systemPrompt   sql.NullString
		knowledgeBaseID sql.NullString
		llmProvider    string
		llmModel       string
		llmConfigJSON  string
		tagsJSON       string
		status         string
		published      bool
		version        int
		public         bool
		creatorID      string
		createdAt      time.Time
		updatedAt      time.Time
		publishedAt    sql.NullTime
	)

	err := scanner.Scan(
		&id, &tenantID, &name, &description, &avatar,
		&systemPrompt, &knowledgeBaseID, &llmProvider, &llmModel,
		&llmConfigJSON, &tagsJSON, &status, &published, &version, &public,
		&creatorID, &createdAt, &updatedAt, &publishedAt,
	)
	if err != nil {
		return nil, err
	}

	// 从数据库行重建领域实体
	bot, err := entity.ReconstructBot(
		id, tenantID, name, description.String, avatar.String,
		systemPrompt.String, knowledgeBaseID.String,
		llmProvider, llmModel, llmConfigJSON, tagsJSON,
		entity.BotStatus(status), published, version, public,
		creatorID, createdAt, updatedAt, publishedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct bot: %w", err)
	}

	return bot, nil
}
