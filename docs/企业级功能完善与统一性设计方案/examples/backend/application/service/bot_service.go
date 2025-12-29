// Package service 提供应用服务
//
// 应用服务层负责：
// - 用例编排（协调领域对象完成业务用例）
// - DTO转换（领域实体 <-> DTO）
// - 事务管理
// - 权限检查（委托给领域服务）
// - 缓存管理（可选）
package service

import (
	"context"
	"fmt"

	"github.com/coze-studio/application/dto"
	"github.com/coze-studio/domain/entity"
	"github.com/coze-studio/domain/repository"
	"github.com/coze-studio/crossdomain/auth"
	"github.com/coze-studio/crossdomain/permission"
)

// BotApplicationService Bot应用服务
//
// 应用服务封装用例，协调领域对象和基础设施
// 本服务不包含业务逻辑，业务逻辑在领域层
type BotApplicationService struct {
	botRepo          repository.BotRepository
	conversationRepo repository.ConversationRepository
	permissionSvc    *permission.PermissionService
	logger           Logger
}

// Logger 日志接口（依赖倒置）
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

// NewBotApplicationService 创建Bot应用服务
func NewBotApplicationService(
	botRepo repository.BotRepository,
	conversationRepo repository.ConversationRepository,
	permissionSvc *permission.PermissionService,
	logger Logger,
) *BotApplicationService {
	return &BotApplicationService{
		botRepo:          botRepo,
		conversationRepo: conversationRepo,
		permissionSvc:    permissionSvc,
		logger:           logger,
	}
}

// ============ 用例：创建Bot ============

// CreateBot 创建Bot用例
//
// 本方法展示完整的创建流程：
// 1. 参数验证
// 2. 权限检查
// 3. 实体创建（使用工厂方法）
// 4. 持久化
// 5. 事件发布
// 6. 返回结果
func (s *BotApplicationService) CreateBot(
	ctx context.Context,
	req *dto.CreateBotRequest,
	currentUser *auth.UserContext,
) (*dto.BotResponse, error) {
	// 步骤1: 验证请求参数
	if err := s.validateCreateBotRequest(ctx, req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 步骤2: 权限检查
	if err := s.permissionSvc.CheckPermission(ctx, currentUser, "bot.create"); err != nil {
		return nil, fmt.Errorf("permission denied: %w", err)
	}

	// 步骤3: 检查名称重复
	existingBot, err := s.botRepo.FindByName(ctx, currentUser.TenantID, req.Name)
	if err == nil && existingBot != nil {
		return nil, fmt.Errorf("bot with name '%s' already exists", req.Name)
	}

	// 步骤4: 创建领域实体（使用工厂方法）
	botID := entity.GenerateBotID()
	bot, err := entity.NewBot(
		botID,
		currentUser.TenantID,
		req.Name,
		req.Description,
		req.Avatar,
		req.SystemPrompt,
		req.KnowledgeBaseID,
		req.LLMConfig,
		req.Tags,
		req.Public,
		currentUser.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot entity: %w", err)
	}

	// 步骤5: 持久化
	if err := s.botRepo.Save(ctx, bot); err != nil {
		s.logger.Error("failed to save bot",
			"bot_id", botID,
			"error", err,
		)
		return nil, fmt.Errorf("failed to save bot: %w", err)
	}

	// 步骤6: 获取领域事件并发布
	events := bot.GetEvents()
	if err := s.publishEvents(ctx, events); err != nil {
		s.logger.Warn("failed to publish events", "error", err)
		// 事件发布失败不影响主流程
	}

	// 步骤7: 转换为DTO并返回
	response := &dto.BotResponse{}
	response.FromEntity(bot)

	s.logger.Info("bot created successfully",
		"bot_id", botID,
		"name", req.Name,
		"user_id", currentUser.UserID,
	)

	return response, nil
}

// ============ 用例：更新Bot ============

// UpdateBot 更新Bot用例
func (s *BotApplicationService) UpdateBot(
	ctx context.Context,
	botID string,
	req *dto.UpdateBotRequest,
	currentUser *auth.UserContext,
) (*dto.BotResponse, error) {
	// 步骤1: 加载Bot
	bot, err := s.botRepo.FindByID(ctx, currentUser.TenantID, botID)
	if err != nil {
		return nil, fmt.Errorf("bot not found: %w", err)
	}

	// 步骤2: 权限检查（必须是创建者或有编辑权限）
	if err := s.permissionSvc.CheckBotOwnership(ctx, currentUser, bot); err != nil {
		return nil, fmt.Errorf("permission denied: %w", err)
	}

	// 步骤3: 应用更新
	if req.Name != nil {
		if err := bot.UpdateName(*req.Name); err != nil {
			return nil, fmt.Errorf("invalid name: %w", err)
		}
	}

	if req.Description != nil {
		bot.UpdateDescription(*req.Description)
	}

	if req.Avatar != nil {
		bot.UpdateAvatar(*req.Avatar)
	}

	if req.SystemPrompt != nil {
		bot.UpdateSystemPrompt(*req.SystemPrompt)
	}

	if req.LLMConfig != nil {
		bot.UpdateLLMConfig(req.LLMConfig)
	}

	if req.Tags != nil {
		bot.UpdateTags(req.Tags)
	}

	if req.Public != nil {
		bot.SetPublic(*req.Public)
	}

	if req.KnowledgeBaseID != nil {
		bot.UpdateKnowledgeBaseID(*req.KnowledgeBaseID)
	}

	// 步骤4: 持久化
	if err := s.botRepo.Save(ctx, bot); err != nil {
		return nil, fmt.Errorf("failed to update bot: %w", err)
	}

	// 步骤5: 转换为DTO并返回
	response := &dto.BotResponse{}
	response.FromEntity(bot)

	return response, nil
}

// ============ 用例：发布Bot ============

// PublishBot 发布Bot用例
func (s *BotApplicationService) PublishBot(
	ctx context.Context,
	botID string,
	req *dto.PublishBotRequest,
	currentUser *auth.UserContext,
) (*dto.BotResponse, error) {
	// 步骤1: 加载Bot
	bot, err := s.botRepo.FindByID(ctx, currentUser.TenantID, botID)
	if err != nil {
		return nil, fmt.Errorf("bot not found: %w", err)
	}

	// 步骤2: 权限检查
	if err := s.permissionSvc.CheckBotOwnership(ctx, currentUser, bot); err != nil {
		return nil, fmt.Errorf("permission denied: %w", err)
	}

	// 步骤3: 发布Bot
	version := "auto"
	if req.Version != nil {
		version = *req.Version
	}

	if err := bot.Publish(version, req.ReleaseNotes); err != nil {
		return nil, fmt.Errorf("failed to publish bot: %w", err)
	}

	// 步骤4: 持久化
	if err := s.botRepo.Save(ctx, bot); err != nil {
		return nil, fmt.Errorf("failed to save bot: %w", err)
	}

	// 步骤5: 获取领域事件并发布
	events := bot.GetEvents()
	if err := s.publishEvents(ctx, events); err != nil {
		s.logger.Warn("failed to publish events", "error", err)
	}

	// 步骤6: 转换为DTO并返回
	response := &dto.BotResponse{}
	response.FromEntity(bot)

	s.logger.Info("bot published successfully",
		"bot_id", botID,
		"version", bot.Version(),
		"user_id", currentUser.UserID,
	)

	return response, nil
}

// ============ 用例：查询Bot列表 ============

// ListBots 查询Bot列表用例
func (s *BotApplicationService) ListBots(
	ctx context.Context,
	req *dto.ListBotsRequest,
	currentUser *auth.UserContext,
) (*dto.BotListResponse, error) {
	// 步骤1: 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	// 步骤2: 权限检查（只能查看自己有权限的Bot）
	// 这里可以添加数据权限过滤逻辑

	// 步骤3: 查询列表
	listReq := &repository.ListBotsRequest{
		TenantID:   currentUser.TenantID,
		Page:       req.Page,
		PageSize:   req.PageSize,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Name:       req.Name,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}

	// 应用状态过滤
	if req.Status != nil {
		status := entity.BotStatus(*req.Status)
		listReq.Status = &status
	}

	// 应用发布过滤
	if req.Published != nil {
		listReq.Published = req.Published
	}

	result, err := s.botRepo.FindByTenantIDWithPaging(ctx, listReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list bots: %w", err)
	}

	// 步骤4: 转换为DTO
	bots := make([]*dto.BotResponse, 0, len(result.Bots))
	for _, bot := range result.Bots {
		botDTO := &dto.BotResponse{}
		botDTO.FromEntity(bot)
		bots = append(bots, botDTO)
	}

	response := &dto.BotListResponse{
		Bots:       bots,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}

	return response, nil
}

// ============ 用例：删除Bot ============

// DeleteBot 删除Bot用例
func (s *BotApplicationService) DeleteBot(
	ctx context.Context,
	botID string,
	currentUser *auth.UserContext,
) error {
	// 步骤1: 加载Bot
	bot, err := s.botRepo.FindByID(ctx, currentUser.TenantID, botID)
	if err != nil {
		return fmt.Errorf("bot not found: %w", err)
	}

	// 步骤2: 权限检查
	if err := s.permissionSvc.CheckBotOwnership(ctx, currentUser, bot); err != nil {
		return fmt.Errorf("permission denied: %w", err)
	}

	// 步骤3: 检查是否有正在进行的对话
	conversationCount, err := s.conversationRepo.CountActiveByBotID(ctx, currentUser.TenantID, botID)
	if err != nil {
		return fmt.Errorf("failed to check conversations: %w", err)
	}

	if conversationCount > 0 {
		return fmt.Errorf("cannot delete bot with active conversations (count: %d)", conversationCount)
	}

	// 步骤4: 删除（软删除）
	if err := s.botRepo.Delete(ctx, currentUser.TenantID, botID); err != nil {
		return fmt.Errorf("failed to delete bot: %w", err)
	}

	s.logger.Info("bot deleted successfully",
		"bot_id", botID,
		"user_id", currentUser.UserID,
	)

	return nil
}

// ============ 私有辅助方法 ============

// validateCreateBotRequest 验证创建请求
func (s *BotApplicationService) validateCreateBotRequest(
	ctx context.Context,
	req *dto.CreateBotRequest,
) error {
	// 参数验证
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(req.Name) < 2 || len(req.Name) > 100 {
		return fmt.Errorf("name must be between 2 and 100 characters")
	}

	if req.Description != "" && len(req.Description) > 500 {
		return fmt.Errorf("description must not exceed 500 characters")
	}

	if req.SystemPrompt != "" && len(req.SystemPrompt) > 10000 {
		return fmt.Errorf("system prompt must not exceed 10000 characters")
	}

	if req.Tags != nil && len(req.Tags) > 10 {
		return fmt.Errorf("tags must not exceed 10")
	}

	// LLM配置验证
	if req.LLMConfig != nil {
		if req.LLMConfig.Provider == "" {
			return fmt.Errorf("LLM provider is required")
		}
		if req.LLMConfig.Model == "" {
			return fmt.Errorf("LLM model is required")
		}
	}

	return nil
}

// publishEvents 发布领域事件
func (s *BotApplicationService) publishEvents(
	ctx context.Context,
	events []interface{},
) error {
	// 这里调用事件发布器
	// 实际实现中会使用消息队列（如NSQ）或事件总线
	for _, event := range events {
		s.logger.Info("publishing event", "event_type", fmt.Sprintf("%T", event))
		// eventBus.Publish(ctx, event)
	}
	return nil
}
