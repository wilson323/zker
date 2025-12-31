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

package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/infra/queue"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

// TaskType 任务类型常量
const (
	TaskTypeKnowledgeDocument = "knowledge_document"
	TaskTypeWorkflowExecution = "workflow_execution"
	TaskTypeBotPublish        = "bot_publish"
	TaskTypeReportGeneration  = "report_generation"
	TaskTypeTokenUsage        = "token_usage"
	TaskTypeEmailNotification = "email_notification"
)

// ============================================================
// 知识库文档处理任务
// ============================================================

// KnowledgeDocumentProcessTask 知识库文档处理任务
type KnowledgeDocumentProcessTask struct {
	DocumentID string `json:"document_id"`
	KnowledgeID string `json:"knowledge_id"`
	Action     string `json:"action"` // "parse", "index", "delete", "reindex"
	Priority   int    `json:"priority"`
	TenantID   string `json:"tenant_id"`
	UserID     string `json:"user_id"`
}

// KnowledgeDocumentHandler 知识库文档处理器
type KnowledgeDocumentHandler struct {
	knowledgeService KnowledgeService
	logger           *zap.Logger
}

// KnowledgeService 知识库服务接口
type KnowledgeService interface {
	ParseDocument(ctx context.Context, documentID string) error
	IndexDocument(ctx context.Context, documentID string) error
	DeleteDocument(ctx context.Context, documentID string) error
	ReindexDocument(ctx context.Context, documentID string) error
}

// NewKnowledgeDocumentHandler 创建知识库文档处理器
func NewKnowledgeDocumentHandler(svc KnowledgeService, logger *zap.Logger) *KnowledgeDocumentHandler {
	return &KnowledgeDocumentHandler{
		knowledgeService: svc,
		logger:           logger,
	}
}

// HandleTask 实现TaskHandler接口
func (h *KnowledgeDocumentHandler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
	var task KnowledgeDocumentProcessTask
	if err := json.Unmarshal(taskData, &task); err != nil {
		return fmt.Errorf("failed to unmarshal task: %w", err)
	}

	h.logger.Info("processing knowledge document task",
		zap.String("task_type", taskType),
		zap.String("document_id", task.DocumentID),
		zap.String("action", task.Action),
	)

	switch task.Action {
	case "parse":
		return h.knowledgeService.ParseDocument(ctx, task.DocumentID)
	case "index":
		return h.knowledgeService.IndexDocument(ctx, task.DocumentID)
	case "delete":
		return h.knowledgeService.DeleteDocument(ctx, task.DocumentID)
	case "reindex":
		return h.knowledgeService.ReindexDocument(ctx, task.DocumentID)
	default:
		return fmt.Errorf("unknown action: %s", task.Action)
	}
}

// ============================================================
// 工作流执行任务
// ============================================================

// WorkflowExecutionTask 工作流执行任务
type WorkflowExecutionTask struct {
	WorkflowID  string                 `json:"workflow_id"`
	VersionID   string                 `json:"version_id,omitempty"`
	Input       map[string]interface{} `json:"input"`
	TriggeredBy  string                 `json:"triggered_by"` // "user", "schedule", "event"
	TenantID    string                 `json:"tenant_id"`
	UserID      string                 `json:"user_id"`
	Priority    int                    `json:"priority"`
	ScheduledAt string                 `json:"scheduled_at,omitempty"` // ISO 8601格式
}

// WorkflowExecutionHandler 工作流执行处理器
type WorkflowExecutionHandler struct {
	workflowService WorkflowService
	logger          *zap.Logger
}

// WorkflowService 工作流服务接口
type WorkflowService interface {
	ExecuteWorkflow(ctx context.Context, workflowID string, input map[string]interface{}) (string, error)
	ExecuteWorkflowVersion(ctx context.Context, versionID string, input map[string]interface{}) (string, error)
	GetWorkflowStatus(ctx context.Context, executionID string) (*WorkflowExecutionStatus, error)
}

// WorkflowExecutionStatus 工作流执行状态
type WorkflowExecutionStatus struct {
	ExecutionID string
	Status      string // "running", "completed", "failed"
	Progress    float64
	Result      map[string]interface{}
}

// NewWorkflowExecutionHandler 创建工作流执行处理器
func NewWorkflowExecutionHandler(svc WorkflowService, logger *zap.Logger) *WorkflowExecutionHandler {
	return &WorkflowExecutionHandler{
		workflowService: svc,
		logger:          logger,
	}
}

// HandleTask 实现TaskHandler接口
func (h *WorkflowExecutionHandler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
	var task WorkflowExecutionTask
	if err := json.Unmarshal(taskData, &task); err != nil {
		return fmt.Errorf("failed to unmarshal task: %w", err)
	}

	h.logger.Info("executing workflow task",
		zap.String("task_type", taskType),
		zap.String("workflow_id", task.WorkflowID),
		zap.String("triggered_by", task.TriggeredBy),
	)

	// 执行工作流
	var executionID string
	var err error

	if task.VersionID != "" {
		executionID, err = h.workflowService.ExecuteWorkflowVersion(ctx, task.VersionID, task.Input)
	} else {
		executionID, err = h.workflowService.ExecuteWorkflow(ctx, task.WorkflowID, task.Input)
	}

	if err != nil {
		return fmt.Errorf("failed to execute workflow: %w", err)
	}

	h.logger.Info("workflow execution started",
		zap.String("execution_id", executionID),
	)

	return nil
}

// ============================================================
// Bot发布任务
// ============================================================

// BotPublishTask Bot发布任务
type BotPublishTask struct {
	BotID      string   `json:"bot_id"`
	VersionID  string   `json:"version_id"`
	PublishTo  []string `json:"publish_to"` // ["weixin", "feishu", "api"]
	Priority   int      `json:"priority"`
	TenantID   string   `json:"tenant_id"`
	UserID     string   `json:"user_id"`
	AutoEnable bool     `json:"auto_enable"` // 是否自动启用
}

// BotPublishHandler Bot发布处理器
type BotPublishHandler struct {
	botService BotService
	logger     *zap.Logger
}

// BotService Bot服务接口
type BotService interface {
	PublishBot(ctx context.Context, botID string, publishTo []string) error
	EnableBot(ctx context.Context, botID string) error
	DisableBot(ctx context.Context, botID string) error
	GetPublishStatus(ctx context.Context, botID string) (*BotPublishStatus, error)
}

// BotPublishStatus Bot发布状态
type BotPublishStatus struct {
	BotID      string
	Status     string // "publishing", "published", "failed"
	PublishedTo []string
	Errors     map[string]string
}

// NewBotPublishHandler 创建Bot发布处理器
func NewBotPublishHandler(svc BotService, logger *zap.Logger) *BotPublishHandler {
	return &BotPublishHandler{
		botService: svc,
		logger:     logger,
	}
}

// HandleTask 实现TaskHandler接口
func (h *BotPublishHandler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
	var task BotPublishTask
	if err := json.Unmarshal(taskData, &task); err != nil {
		return fmt.Errorf("failed to unmarshal task: %w", err)
	}

	h.logger.Info("publishing bot task",
		zap.String("task_type", taskType),
		zap.String("bot_id", task.BotID),
		zap.Strings("publish_to", task.PublishTo),
	)

	// 发布Bot
	if err := h.botService.PublishBot(ctx, task.BotID, task.PublishTo); err != nil {
		return fmt.Errorf("failed to publish bot: %w", err)
	}

	// 如果需要自动启用
	if task.AutoEnable {
		if err := h.botService.EnableBot(ctx, task.BotID); err != nil {
			h.logger.Error("failed to auto-enable bot",
				zap.String("bot_id", task.BotID),
				zap.Error(err),
			)
			// 不返回错误，因为发布已成功
		}
	}

	h.logger.Info("bot published successfully",
		zap.String("bot_id", task.BotID),
	)

	return nil
}

// ============================================================
// 报告生成任务
// ============================================================

// ReportGenerationTask 报告生成任务
type ReportGenerationTask struct {
	ReportID   string `json:"report_id"`
	ReportType string `json:"report_type"` // "usage", "performance", "error"
	StartDate  string `json:"start_date"`  // ISO 8601格式
	EndDate    string `json:"end_date"`    // ISO 8601格式
	TenantID   string `json:"tenant_id"`
	UserID     string `json:"user_id"`
	Format     string `json:"format"`      // "pdf", "xlsx", "csv"
	EmailTo    string `json:"email_to,omitempty"`
}

// ReportGenerationHandler 报告生成处理器
type ReportGenerationHandler struct {
	reportService ReportService
	logger        *zap.Logger
}

// ReportService 报告服务接口
type ReportService interface {
	GenerateUsageReport(ctx context.Context, reportID, startDate, endDate string) error
	GeneratePerformanceReport(ctx context.Context, reportID, startDate, endDate string) error
	GenerateErrorReport(ctx context.Context, reportID, startDate, endDate string) error
	SendReportByEmail(ctx context.Context, reportID, emailTo string) error
	GetReportStatus(ctx context.Context, reportID string) (*ReportStatus, error)
}

// ReportStatus 报告状态
type ReportStatus struct {
	ReportID   string
	Status     string // "generating", "completed", "failed"
	FileURL    string
	FileSize   int64
	GeneratedAt string
}

// NewReportGenerationHandler 创建报告生成处理器
func NewReportGenerationHandler(svc ReportService, logger *zap.Logger) *ReportGenerationHandler {
	return &ReportGenerationHandler{
		reportService: svc,
		logger:        logger,
	}
}

// HandleTask 实现TaskHandler接口
func (h *ReportGenerationHandler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
	var task ReportGenerationTask
	if err := json.Unmarshal(taskData, &task); err != nil {
		return fmt.Errorf("failed to unmarshal task: %w", err)
	}

	h.logger.Info("generating report task",
		zap.String("task_type", taskType),
		zap.String("report_id", task.ReportID),
		zap.String("report_type", task.ReportType),
	)

	// 生成报告
	var err error
	switch task.ReportType {
	case "usage":
		err = h.reportService.GenerateUsageReport(ctx, task.ReportID, task.StartDate, task.EndDate)
	case "performance":
		err = h.reportService.GeneratePerformanceReport(ctx, task.ReportID, task.StartDate, task.EndDate)
	case "error":
		err = h.reportService.GenerateErrorReport(ctx, task.ReportID, task.StartDate, task.EndDate)
	default:
		return fmt.Errorf("unknown report type: %s", task.ReportType)
	}

	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// 如果指定了邮箱，发送报告
	if task.EmailTo != "" {
		if err := h.reportService.SendReportByEmail(ctx, task.ReportID, task.EmailTo); err != nil {
			h.logger.Error("failed to send report by email",
				zap.String("report_id", task.ReportID),
				zap.String("email_to", task.EmailTo),
				zap.Error(err),
			)
			// 不返回错误，因为报告已生成
		}
	}

	h.logger.Info("report generated successfully",
		zap.String("report_id", task.ReportID),
	)

	return nil
}

// ============================================================
// Token计量任务
// ============================================================

// TokenUsageTask Token计量任务
type TokenUsageTask struct {
	TenantID         string  `json:"tenant_id"`
	UserID           string  `json:"user_id"`
	BotID            string  `json:"bot_id"`
	ModelID          string  `json:"model_id"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	CostUSD          float64 `json:"cost_usd"`
	ConversationID   string  `json:"conversation_id"`
	MessageID        string  `json:"message_id"`
}

// TokenUsageHandler Token计量处理器
type TokenUsageHandler struct {
	tokenService TokenUsageService
	logger       *zap.Logger
}

// TokenUsageService Token计量服务接口
type TokenUsageService interface {
	RecordTokenUsage(ctx context.Context, usage *TokenUsageTask) error
	CheckQuota(ctx context.Context, tenantID string) (bool, error)
	GetUsageStats(ctx context.Context, tenantID, startDate, endDate string) (*UsageStats, error)
}

// UsageStats 使用统计
type UsageStats struct {
	TenantID       string
	TotalTokens    int64
	TotalCost      float64
	MessageCount   int
	BotUsage       map[string]int64
	ModelUsage     map[string]int64
}

// NewTokenUsageHandler 创建Token计量处理器
func NewTokenUsageHandler(svc TokenUsageService, logger *zap.Logger) *TokenUsageHandler {
	return &TokenUsageHandler{
		tokenService: svc,
		logger:       logger,
	}
}

// HandleTask 实现TaskHandler接口
func (h *TokenUsageHandler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
	var task TokenUsageTask
	if err := json.Unmarshal(taskData, &task); err != nil {
		return fmt.Errorf("failed to unmarshal task: %w", err)
	}

	h.logger.Debug("recording token usage",
		zap.String("task_type", taskType),
		zap.String("tenant_id", task.TenantID),
		zap.String("model_id", task.ModelID),
		zap.Int("total_tokens", task.TotalTokens),
	)

	if err := h.tokenService.RecordTokenUsage(ctx, &task); err != nil {
		return fmt.Errorf("failed to record token usage: %w", err)
	}

	return nil
}

// ============================================================
// 邮件通知任务
// ============================================================

// EmailNotificationTask 邮件通知任务
type EmailNotificationTask struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Template string `json:"template"` // "welcome", "reset_password", "report_ready"
	Data     map[string]interface{} `json:"data"`
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

// EmailNotificationHandler 邮件通知处理器
type EmailNotificationHandler struct {
	emailService EmailService
	logger       *zap.Logger
}

// EmailService 邮件服务接口
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, template string, data map[string]interface{}) error
	SendRawEmail(ctx context.Context, to, subject, body string) error
}

// NewEmailNotificationHandler 创建邮件通知处理器
func NewEmailNotificationHandler(svc EmailService, logger *zap.Logger) *EmailNotificationHandler {
	return &EmailNotificationHandler{
		emailService: svc,
		logger:       logger,
	}
}

// HandleTask 实现TaskHandler接口
func (h *EmailNotificationHandler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
	var task EmailNotificationTask
	if err := json.Unmarshal(taskData, &task); err != nil {
		return fmt.Errorf("failed to unmarshal task: %w", err)
	}

	h.logger.Info("sending email notification",
		zap.String("task_type", taskType),
		zap.String("to", task.To),
		zap.String("template", task.Template),
	)

	if err := h.emailService.SendEmail(ctx, task.To, task.Subject, task.Template, task.Data); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// ============================================================
// 便捷函数：将NSQ MessageHandler转换为TaskHandler
// ============================================================

// MessageHandlerToTaskHandler 将NSQ MessageHandler转换为TaskHandler
func MessageHandlerToTaskHandler(handler queue.MessageHandler) queue.TaskHandler {
	return queue.TaskHandlerFunc(func(ctx context.Context, taskType string, taskData []byte) error {
		msg := &nsq.Message{
			Body: taskData,
		}
		return handler.HandleMessage(ctx, msg)
	})
}
