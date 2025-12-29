// Package dto 提供数据传输对象
//
// DTO用于API层和应用层之间的数据传输，职责包括：
// - 数据验证
// - 数据转换（领域实体 <-> DTO）
// - API请求和响应的格式化
package dto

import (
	"time"

	"github.com/coze-studio/domain/entity"
)

// ============ 请求 DTO ============

// CreateBotRequest 创建Bot请求
type CreateBotRequest struct {
	// Name Bot名称（必填，2-100字符）
	Name string `json:"name" validate:"required,min=2,max=100"`
	// Description Bot描述（可选，最多500字符）
	Description string `json:"description" validate:"max=500"`
	// Avatar 头像URL
	Avatar string `json:"avatar" validate:"url"`
	// SystemPrompt 系统提示词
	SystemPrompt string `json:"system_prompt" validate:"max=10000"`
	// KnowledgeBaseID 关联的知识库ID
	KnowledgeBaseID *string `json:"knowledge_base_id"`
	// LLMConfig LLM配置
	LLMConfig *LLMConfigDTO `json:"llm_config"`
	// Tags 标签列表
	Tags []string `json:"tags" validate:"max=10"`
	// Public 是否公开（公开Bot可以被其他租户使用）
	Public bool `json:"public"`
}

// UpdateBotRequest 更新Bot请求
type UpdateBotRequest struct {
	// Name Bot名称
	Name *string `json:"name" validate:"omitempty,min=2,max=100"`
	// Description Bot描述
	Description *string `json:"description" validate:"omitempty,max=500"`
	// Avatar 头像URL
	Avatar *string `json:"avatar" validate:"omitempty,url"`
	// SystemPrompt 系统提示词
	SystemPrompt *string `json:"system_prompt" validate:"omitempty,max=10000"`
	// KnowledgeBaseID 关联的知识库ID
	KnowledgeBaseID *string `json:"knowledge_base_id"`
	// LLMConfig LLM配置
	LLMConfig *LLMConfigDTO `json:"llm_config"`
	// Tags 标签列表
	Tags []string `json:"tags" validate:"max=10"`
	// Public 是否公开
	Public *bool `json:"public"`
}

// PublishBotRequest 发布Bot请求
type PublishBotRequest struct {
	// Version 版本号（可选，不提供则自动递增）
	Version *string `json:"version"`
	// ReleaseNotes 发布说明
	ReleaseNotes string `json:"release_notes" validate:"max=1000"`
}

// ListBotsRequest 查询Bot列表请求
type ListBotsRequest struct {
	// Page 页码（默认1）
	Page int `json:"page" validate:"min=1"`
	// PageSize 每页大小（默认20，最大100）
	PageSize int `json:"page_size" validate:"min=1,max=100"`
	// SortBy 排序字段（默认created_at）
	SortBy string `json:"sort_by" validate:"oneof=id name created_at updated_at"`
	// SortOrder 排序方向（默认desc）
	SortOrder string `json:"sort_order" validate:"oneof=asc,desc"`

	// Name 名称模糊搜索
	Name string `json:"name"`
	// Status 状态过滤
	Status *string `json:"status" validate:"oneof=draft published archived"`
	// Published 是否发布过滤
	Published *bool `json:"published"`
	// Tag 标签过滤
	Tag string `json:"tag"`
	// CreatorID 创建者过滤
	CreatorID string `json:"creator_id"`
	// StartTime 创建时间范围-开始
	StartTime *time.Time `json:"start_time"`
	// EndTime 创建时间范围-结束
	EndTime *time.Time `json:"end_time"`
}

// ============ 响应 DTO ============

// BotResponse Bot响应
type BotResponse struct {
	// ID Bot ID
	ID string `json:"id"`
	// Name Bot名称
	Name string `json:"name"`
	// Description Bot描述
	Description string `json:"description"`
	// Avatar 头像URL
	Avatar string `json:"avatar"`
	// SystemPrompt 系统提示词
	SystemPrompt string `json:"system_prompt"`
	// KnowledgeBaseID 关联的知识库ID
	KnowledgeBaseID *string `json:"knowledge_base_id"`
	// LLMConfig LLM配置
	LLMConfig *LLMConfigDTO `json:"llm_config"`
	// Tags 标签列表
	Tags []string `json:"tags"`
	// Status 状态
	Status string `json:"status"`
	// Published 是否已发布
	Published bool `json:"published"`
	// Version 版本号
	Version int `json:"version"`
	// Public 是否公开
	Public bool `json:"public"`
	// Statistics 统计信息
	Statistics *BotStatisticsDTO `json:"statistics"`
	// CreatorID 创建者ID
	CreatorID string `json:"creator_id"`
	// CreatorName 创建者名称
	CreatorName string `json:"creator_name"`
	// TenantID 租户ID
	TenantID string `json:"tenant_id"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
	// PublishedAt 发布时间
	PublishedAt *time.Time `json:"published_at"`
}

// BotListResponse Bot列表响应
type BotListResponse struct {
	// Bots Bot列表
	Bots []*BotResponse `json:"bots"`
	// Total 总数
	Total int64 `json:"total"`
	// Page 当前页码
	Page int `json:"page"`
	// PageSize 每页大小
	PageSize int `json:"page_size"`
	// TotalPages 总页数
	TotalPages int64 `json:"total_pages"`
}

// BotStatisticsDTO Bot统计信息
type BotStatisticsDTO struct {
	// TotalConversations 总对话数
	TotalConversations int64 `json:"total_conversations"`
	// TotalMessages 总消息数
	TotalMessages int64 `json:"total_messages"`
	// TotalTokensUsed 总Token使用量
	TotalTokensUsed int64 `json:"total_tokens_used"`
	// AvgConversationLength 平均对话长度
	AvgConversationLength float64 `json:"avg_conversation_length"`
	// LastUsedAt 最后使用时间
	LastUsedAt *time.Time `json:"last_used_at"`
}

// LLMConfigDTO LLM配置
type LLMConfigDTO struct {
	// Provider 提供商（openai, claude, qwen等）
	Provider string `json:"provider" validate:"required"`
	// Model 模型名称
	Model string `json:"model" validate:"required"`
	// Temperature 温度参数（0-2）
	Temperature *float64 `json:"temperature" validate:"omitempty,min=0,max=2"`
	// MaxTokens 最大Token数
	MaxTokens *int `json:"max_tokens" validate:"omitempty,min=1,max=128000"`
	// TopP Top-P采样参数
	TopP *float64 `json:"top_p" validate:"omitempty,min=0,max=1"`
	// FrequencyPenalty 频率惩罚
	FrequencyPenalty *float64 `json:"frequency_penalty" validate:"omitempty,min=-2,max=2"`
	// PresencePenalty 存在惩罚
	PresencePenalty *float64 `json:"presence_penalty" validate:"omitempty,min=-2,max=2"`
	// Stop 停止序列
	Stop []string `json:"stop" validate:"max=4"`
}

// ============ DTO转换方法 ============

// FromEntity 从领域实体转换为DTO
func (dto *BotResponse) FromEntity(bot *entity.Bot) {
	dto.ID = bot.ID().String()
	dto.Name = bot.Name()
	dto.Description = bot.Description()
	dto.Avatar = bot.Avatar()
	dto.SystemPrompt = bot.SystemPrompt()
	dto.KnowledgeBaseID = bot.KnowledgeBaseID()
	dto.LLMConfig = &LLMConfigDTO{
		Provider: bot.LLMProvider(),
		Model:    bot.LLMModel(),
		// 其他LLM配置字段...
	}
	dto.Tags = bot.Tags()
	dto.Status = string(bot.Status())
	dto.Published = bot.IsPublished()
	dto.Version = bot.Version()
	dto.Public = bot.IsPublic()
	dto.Statistics = &BotStatisticsDTO{
		TotalConversations:   bot.TotalConversations(),
		TotalMessages:        bot.TotalMessages(),
		TotalTokensUsed:      bot.TotalTokensUsed(),
		AvgConversationLength: bot.AvgConversationLength(),
		LastUsedAt:           bot.LastUsedAt(),
	}
	dto.CreatorID = bot.CreatorID()
	dto.TenantID = bot.TenantID().String()
	dto.CreatedAt = bot.CreatedAt()
	dto.UpdatedAt = bot.UpdatedAt()
	dto.PublishedAt = bot.PublishedAt()
}

// ToEntityForCreate 从创建请求转换为领域实体（部分字段）
func (req *CreateBotRequest) ToEntityForCreate(
	tenantID entity.TenantID,
	botID string,
	creatorID string,
) *entity.Bot {
	// 这里需要调用实体构造函数
	// 实际实现中，实体应该在应用服务中创建
	return nil
}

// ============ 错误响应 DTO ============

// ErrorResponse 错误响应
type ErrorResponse struct {
	// Code 错误码
	Code string `json:"code"`
	// Message 错误消息
	Message string `json:"message"`
	// Details 详细信息（可选）
	Details interface{} `json:"details,omitempty"`
	// RequestID 请求追踪ID
	RequestID string `json:"request_id"`
}

// ValidationErrorDetail 验证错误详情
type ValidationErrorDetail struct {
	// Field 字段名
	Field string `json:"field"`
	// Message 错误消息
	Message string `json:"message"`
	// Value 当前值
	Value interface{} `json:"value,omitempty"`
}
