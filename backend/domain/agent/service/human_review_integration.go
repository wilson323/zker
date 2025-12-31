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
	"fmt"

	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
	hilservice "github.com/coze-dev/coze-studio/backend/domain/humaninloop/service"
)

// HumanReviewService 人工审核集成服务
// 这个服务展示了如何在现有业务中集成人机协同功能
type HumanReviewService struct {
	orchestrator hilservice.CollaborationOrchestrator
	botRepo      BotRepository     // Bot仓储，需要注入
	configRepo   ConfigRepository  // 配置仓储，需要注入
}

// NewHumanReviewService 创建人工审核服务
func NewHumanReviewService(
	orchestrator hilservice.CollaborationOrchestrator,
	botRepo BotRepository,
	configRepo ConfigRepository,
) *HumanReviewService {
	return &HumanReviewService{
		orchestrator: orchestrator,
		botRepo:      botRepo,
		configRepo:   configRepo,
	}
}

// ==================== 示例1: Bot创建需要人工审核 ====================

// CreateBotWithReview 创建Bot并触发人工审核
func (s *HumanReviewService) CreateBotWithReview(
	ctx context.Context,
	req *CreateBotRequest,
	aiConfidence float64,
) (*Bot, error) {
	// 1. 创建Bot（草稿状态）
	bot := &Bot{
		BotID:   uuid.New().String(),
		Name:    req.Name,
		Config:  req.Config,
		Status:  "draft",
		// ... 其他字段
	}

	if err := s.botRepo.Create(ctx, bot); err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	// 2. 检查是否需要人工审核
	if s.shouldTriggerReview(ctx, bot, aiConfidence) {
		// 3. 创建人机协同任务
		task, err := s.orchestrator.CreateTask(ctx, &hilservice.CreateTaskRequest{
			TenantID:      getTenantID(ctx),
			TaskType:      entity.TaskTypeReview,
			Source:        entity.SourceAI,
			Priority:      s.determinePriority(bot, aiConfidence),
			BotID:         &bot.BotID,
			Context: map[string]interface{}{
				"bot_id":   bot.BotID,
				"bot_name": bot.Name,
				"config":   bot.Config,
				"reason":   "Bot创建需要审核",
			},
			TriggerReason: fmt.Sprintf("AI置信度 %.2f 低于阈值", aiConfidence),
		})

		if err != nil {
			return nil, fmt.Errorf("failed to create review task: %w", err)
		}

		// 4. 记录任务ID到Bot
		bot.ReviewTaskID = &task.TaskID
		s.botRepo.Update(ctx, bot)
	}

	return bot, nil
}

// ==================== 示例2: 敏感内容检测触发人工审核 ====================

// CheckAndReviewSensitiveContent 检查敏感内容并触发审核
func (s *HumanReviewService) CheckAndReviewSensitiveContent(
	ctx context.Context,
	conversationID, messageID string,
	content string,
) error {
	// 1. 检测敏感内容
	if s.containsSensitiveContent(content) {
		// 2. 创建审核任务
		_, err := s.orchestrator.CreateTask(ctx, &hilservice.CreateTaskRequest{
			TenantID:       getTenantID(ctx),
			TaskType:       entity.TaskTypeReview,
			Source:         entity.SourceSystem,
			Priority:       entity.PriorityHigh,
			ConversationID: &conversationID,
			MessageID:      &messageID,
			Context: map[string]interface{}{
				"conversation_id": conversationID,
				"message_id":      messageID,
				"content":         content,
				"detected_keywords": s.getSensitiveKeywords(content),
			},
			TriggerReason: "检测到敏感内容",
			SLAMinutes:    30, // 30分钟内处理
		})

		if err != nil {
			return fmt.Errorf("failed to create sensitive content review: %w", err)
		}
	}

	return nil
}

// ==================== 示例3: AI决策置信度低触发人工验证 ====================

// ValidateAIDecision 验证AI决策（当置信度低时）
func (s *HumanReviewService) ValidateAIDecision(
	ctx context.Context,
	botID string,
	aiDecision map[string]interface{},
	confidence float64,
) (map[string]interface{}, error) {
	// 1. 置信度足够，直接返回
	if confidence >= 0.8 {
		return aiDecision, nil
	}

	// 2. 置信度低，需要人工验证
	task, err := s.orchestrator.CreateTask(ctx, &hilservice.CreateTaskRequest{
		TenantID:  getTenantID(ctx),
		TaskType:  entity.TaskTypeValidation,
		Source:    entity.SourceAI,
		Priority:  entity.PriorityMedium,
		BotID:     &botID,
		Context: map[string]interface{}{
			"bot_id":     botID,
			"decision":   aiDecision,
			"confidence": confidence,
		},
		TriggerReason: fmt.Sprintf("AI决策置信度 %.2f 低于阈值 0.8", confidence),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create validation task: %w", err)
	}

	// 3. 等待人工审核结果（同步等待示例）
	// 注意：实际使用时可能需要异步处理
	result, err := s.waitForReview(ctx, task.TaskID, 5*60) // 最多等待5分钟
	if err != nil {
		return nil, fmt.Errorf("review timeout or failed: %w", err)
	}

	return result, nil
}

// ==================== 示例4: 工作流错误修正 ====================

// CorrectWorkflowError 工作流错误修正
func (s *HumanReviewService) CorrectWorkflowError(
	ctx context.Context,
	workflowID string,
	errorMsg string,
	errorContext map[string]interface{},
) error {
	// 创建错误修正任务
	_, err := s.orchestrator.CreateTask(ctx, &hilservice.CreateTaskRequest{
		TenantID:  getTenantID(ctx),
		TaskType:  entity.TaskTypeCorrection,
		Source:    entity.SourceSystem,
		Priority:  entity.PriorityUrgent,
		WorkflowID: &workflowID,
		Context: map[string]interface{}{
			"workflow_id":   workflowID,
			"error_message": errorMsg,
			"error_context": errorContext,
		},
		TriggerReason: fmt.Sprintf("工作流执行错误: %s", errorMsg),
		SLAMinutes:    15, // 紧急，15分钟内处理
	})

	return err
}

// ==================== 示例5: 自动处理审核结果 ====================

// ProcessReviewResult 处理审核结果（异步回调）
func (s *HumanReviewService) ProcessReviewResult(
	ctx context.Context,
	taskID string,
	result map[string]interface{},
) error {
	// 1. 获取任务
	task, err := s.orchestrator.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	// 2. 根据任务类型处理结果
	switch task.TaskType {
	case entity.TaskTypeReview:
		return s.processBotReviewResult(ctx, task, result)

	case entity.TaskTypeCorrection:
		return s.processWorkflowCorrection(ctx, task, result)

	case entity.TaskTypeValidation:
		return s.processValidationResult(ctx, task, result)

	default:
		return fmt.Errorf("unknown task type: %s", task.TaskType)
	}
}

// ==================== 辅助方法 ====================

// shouldTriggerReview 判断是否需要触发审核
func (s *HumanReviewService) shouldTriggerReview(
	ctx context.Context,
	bot *Bot,
	aiConfidence float64,
) bool {
	// 获取租户配置
	config, err := s.configRepo.GetByTenantID(ctx, getTenantID(ctx))
	if err != nil {
		return true // 配置获取失败，默认需要审核
	}

	// 置信度低于阈值，需要审核
	return aiConfidence < config.AutoReviewThreshold
}

// determinePriority 确定任务优先级
func (s *HumanReviewService) determinePriority(
	bot *Bot,
	aiConfidence float64,
) entity.TaskPriority {
	if aiConfidence < 0.3 {
		return entity.PriorityUrgent
	} else if aiConfidence < 0.5 {
		return entity.PriorityHigh
	} else if aiConfidence < 0.7 {
		return entity.PriorityMedium
	}
	return entity.PriorityLow
}

// containsSensitiveContent 检测敏感内容
func (s *HumanReviewService) containsSensitiveContent(content string) bool {
	sensitiveKeywords := []string{"薪资", "合同", "开除", "辞职", "机密"}
	for _, keyword := range sensitiveKeywords {
		if contains(content, keyword) {
			return true
		}
	}
	return false
}

// getSensitiveKeywords 获取检测到的敏感关键词
func (s *HumanReviewService) getSensitiveKeywords(content string) []string {
	sensitiveKeywords := []string{"薪资", "合同", "开除", "辞职", "机密"}
	found := []string{}
	for _, keyword := range sensitiveKeywords {
		if contains(content, keyword) {
			found = append(found, keyword)
		}
	}
	return found
}

// waitForReview 等待审核结果
func (s *HumanReviewService) waitForReview(
	ctx context.Context,
	taskID string,
	timeoutSeconds int,
) (map[string]interface{}, error) {
	// TODO: 实现异步等待逻辑
	// 可以使用轮询、WebSocket、消息队列等方式
	return nil, fmt.Errorf("not implemented")
}

// processBotReviewResult 处理Bot审核结果
func (s *HumanReviewService) processBotReviewResult(
	ctx context.Context,
	task *entity.CollaborationTask,
	result map[string]interface{},
) error {
	if task.BotID == nil {
		return fmt.Errorf("bot_id is nil")
	}

	// 解析审核结果
	decision, _ := result["decision"].(string)
	if decision == "approved" {
		// 审核通过，发布Bot
		return s.botRepo.UpdateStatus(ctx, *task.BotID, "published")
	} else if decision == "rejected" {
		// 审核拒绝，标记为 rejected
		return s.botRepo.UpdateStatus(ctx, *task.BotID, "rejected")
	} else if decision == "modified" {
		// 审核修改，应用修改内容
		modifiedContent, _ := result["modified_content"].(string)
		return s.applyBotModifications(ctx, *task.BotID, modifiedContent)
	}

	return nil
}

// processWorkflowCorrection 处理工作流错误修正
func (s *HumanReviewService) processWorkflowCorrection(
	ctx context.Context,
	task *entity.CollaborationTask,
	result map[string]interface{},
) error {
	// TODO: 实现工作流修正逻辑
	return nil
}

// processValidationResult 处理验证结果
func (s *HumanReviewService) processValidationResult(
	ctx context.Context,
	task *entity.CollaborationTask,
	result map[string]interface{},
) error {
	// TODO: 实现验证结果处理逻辑
	return nil
}

// applyBotModifications 应用Bot修改
func (s *HumanReviewService) applyBotModifications(
	ctx context.Context,
	botID, modifiedContent string,
) error {
	// TODO: 解析并应用修改内容
	return nil
}

// contains 字符串包含辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (
		s[:len(substr)] == substr ||
		s[len(s)-len(substr):] == substr ||
		indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// ==================== 数据结构 ====================

// CreateBotRequest 创建Bot请求
type CreateBotRequest struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
	// ... 其他字段
}

// Bot Bot实体
type Bot struct {
	BotID        string  `json:"bot_id"`
	Name         string  `json:"name"`
	Config       map[string]interface{} `json:"config"`
	Status       string  `json:"status"`
	ReviewTaskID *string `json:"review_task_id,omitempty"`
	// ... 其他字段
}

// BotRepository Bot仓储接口（示例）
type BotRepository interface {
	Create(ctx context.Context, bot *Bot) error
	Update(ctx context.Context, bot *Bot) error
	UpdateStatus(ctx context.Context, botID, status string) error
}

// ConfigRepository 配置仓储接口（示例）
type ConfigRepository interface {
	GetByTenantID(ctx context.Context, tenantID string) (*TenantConfig, error)
}

// TenantConfig 租户配置（示例）
type TenantConfig struct {
	TenantID             string  `json:"tenant_id"`
	AutoReviewThreshold  float64 `json:"auto_review_threshold"` // 自动审核阈值
	// ... 其他配置字段
}

// getTenantID 从上下文获取租户ID（示例）
func getTenantID(ctx context.Context) string {
	// TODO: 从上下文中获取租户ID
	return "tenant-001"
}
