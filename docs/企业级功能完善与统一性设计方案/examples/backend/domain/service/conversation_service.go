// Package service 提供领域服务
//
// 本文件展示对话领域服务的实现，包括：
// - 复杂业务逻辑封装
// - 跨聚合根的协调
// - 领域规则的执行
package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coze-studio/domain/entity"
	"github.com/coze-studio/domain/repository"
)

// ConversationDomainService 对话领域服务
//
// 领域服务封装不属于任何特定实体的复杂业务逻辑
// 当一个操作涉及多个聚合根或需要复杂的规则时，使用领域服务
type ConversationDomainService struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	botRepo          repository.BotRepository
	knowledgeRepo    repository.KnowledgeRepository
}

// NewConversationDomainService 创建对话领域服务
func NewConversationDomainService(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	botRepo repository.BotRepository,
	knowledgeRepo repository.KnowledgeRepository,
) *ConversationDomainService {
	return &ConversationDomainService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		botRepo:          botRepo,
		knowledgeRepo:    knowledgeRepo,
	}
}

// ProcessUserMessage 处理用户消息（核心领域逻辑）
//
// 本方法展示完整的对话处理流程：
// 1. 验证租户状态和Bot可用性
// 2. 检查Token配额
// 3. 构建对话上下文（历史消息）
// 4. 知识库检索（如果启用）
// 5. 调用LLM生成回复
// 6. 保存消息记录
// 7. 更新Token使用量
// 8. 触发相关事件
func (s *ConversationDomainService) ProcessUserMessage(
	ctx context.Context,
	tenantID entity.TenantID,
	conversationID string,
	botID string,
	userMessage string,
	userID string,
) (*entity.Message, *entity.Bot, error) {
	// 步骤1: 加载并验证Bot
	bot, err := s.botRepo.FindByID(ctx, tenantID, botID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find bot: %w", err)
	}

	if !bot.IsPublished() {
		return nil, nil, fmt.Errorf("bot is not published")
	}

	// 步骤2: 加载对话（如果存在）
	var conversation *entity.Conversation
	if conversationID != "" {
		conversation, err = s.conversationRepo.FindByID(ctx, tenantID, conversationID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to find conversation: %w", err)
		}
	} else {
		// 创建新对话
		conversation = entity.NewConversation(tenantID, botID, userID)
	}

	// 步骤3: 检查Token配额
	tokenEstimate := s.estimateTokens(userMessage, 100) // 估算Token使用量
	if err := bot.UseTokens(tokenEstimate); err != nil {
		return nil, nil, fmt.Errorf("token quota exceeded: %w", err)
	}

	// 步骤4: 构建对话上下文
	messages, err := s.messageRepo.FindByConversationID(ctx, tenantID, conversationID, 10)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load conversation history: %w", err)
	}

	// 步骤5: 知识库检索（如果启用）
	var knowledgeContext []string
	if bot.KnowledgeBaseID != "" {
		knowledgeContext, err = s.retrieveKnowledge(ctx, tenantID, bot.KnowledgeBaseID, userMessage)
		if err != nil {
			// 知识库检索失败不影响主流程，记录日志即可
			// logger.Warn("failed to retrieve knowledge", "error", err)
		}
	}

	// 步骤6: 构建LLM请求
	llmMessages := s.buildLLMMessages(messages, userMessage, knowledgeContext, bot.SystemPrompt)

	// 步骤7: 调用LLM（这里使用接口，实际实现在基础设施层）
	llmResponse, tokensUsed, err := s.callLLM(ctx, bot, llmMessages)
	if err != nil {
		return nil, nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// 步骤8: 创建用户消息
	userMsg := entity.NewMessage(
		conversation.ID(),
		entity.MessageRoleUser,
		userMessage,
		userID,
	)

	// 步骤9: 创建助手消息
	assistantMsg := entity.NewMessage(
		conversation.ID(),
		entity.MessageRoleAssistant,
		llmResponse,
		"", // 助手消息不需要UserID
	)
	assistantMsg.SetTokenCount(tokensUsed)

	// 如果有知识库引用，添加到消息元数据
	if len(knowledgeContext) > 0 {
		assistantMsg.AddMetadata("knowledge_used", true)
		assistantMsg.AddMetadata("knowledge_count", len(knowledgeContext))
	}

	// 步骤10: 保存消息
	if err := s.messageRepo.Save(ctx, userMsg); err != nil {
		return nil, nil, fmt.Errorf("failed to save user message: %w", err)
	}

	if err := s.messageRepo.Save(ctx, assistantMsg); err != nil {
		return nil, nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	// 步骤11: 更新对话元数据
	conversation.UpdateLastMessageAt()
	conversation.IncrementMessageCount()
	if err := s.conversationRepo.Save(ctx, conversation); err != nil {
		return nil, nil, fmt.Errorf("failed to update conversation: %w", err)
	}

	// 步骤12: 更新Bot Token使用量
	bot.UseTokens(tokensUsed)
	if err := s.botRepo.Save(ctx, bot); err != nil {
		return nil, nil, fmt.Errorf("failed to update bot token usage: %w", err)
	}

	return assistantMsg, bot, nil
}

// retrieveKnowledge 知识库检索（混合检索）
//
// 本方法展示RAG检索流程：
// 1. 向量检索（语义相似）
// 2. 关键词检索（精确匹配）
// 3. 结果融合
// 4. 重排序（Rerank）
func (s *ConversationDomainService) retrieveKnowledge(
	ctx context.Context,
	tenantID entity.TenantID,
	knowledgeBaseID string,
	query string,
) ([]string, error) {
	// 步骤1: 向量检索（Top-K=10）
	vectorResults, err := s.knowledgeRepo.VectorSearch(ctx, tenantID, knowledgeBaseID, query, 10)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 步骤2: 关键词检索（Top-K=10）
	keywordResults, err := s.knowledgeRepo.KeywordSearch(ctx, tenantID, knowledgeBaseID, query, 10)
	if err != nil {
		return nil, fmt.Errorf("keyword search failed: %w", err)
	}

	// 步骤3: 结果融合（去重并合并）
	fusedResults := s.fuseResults(vectorResults, keywordResults)

	// 步骤4: 重排序（使用reranker模型）
	rerankedResults, err := s.knowledgeRepo.Rerank(ctx, tenantID, knowledgeBaseID, query, fusedResults)
	if err != nil {
		// Rerank失败，返回融合结果
		return s.extractContent(fusedResults[:3]), nil
	}

	// 步骤5: 返回Top-3结果
	return s.extractContent(rerankedResults[:3]), nil
}

// fuseResults 融合检索结果
//
// 使用Reciprocal Rank Fusion (RRF)算法：
// score = sum(1 / (k + rank_i))
// 其中 k=60（经验值）
func (s *ConversationDomainService) fuseResults(
	vectorResults []entity.KnowledgeChunk,
	keywordResults []entity.KnowledgeChunk,
) []entity.KnowledgeChunk {
	k := 60
	scores := make(map[string]float64)
	chunkMap := make(map[string]entity.KnowledgeChunk)

	// 计算向量检索得分
	for rank, chunk := range vectorResults {
		id := chunk.ID()
		scores[id] += 1.0 / float64(k+rank+1)
		chunkMap[id] = chunk
	}

	// 计算关键词检索得分
	for rank, chunk := range keywordResults {
		id := chunk.ID()
		scores[id] += 1.0 / float64(k+rank+1)
		if _, exists := chunkMap[id]; !exists {
			chunkMap[id] = chunk
		}
	}

	// 按得分排序
	type scoredChunk struct {
		chunk entity.KnowledgeChunk
		score float64
	}

	var result []scoredChunk
	for id, score := range scores {
		result = append(result, scoredChunk{
			chunk: chunkMap[id],
			score: score,
		})
	}

	// 降序排序
	// sort.Slice(result, func(i, j int) bool {
	// 	return result[i].score > result[j].score
	// })

	var sorted []entity.KnowledgeChunk
	for _, sc := range result {
		sorted = append(sorted, sc.chunk)
	}

	return sorted
}

// extractContent 提取知识块内容
func (s *ConversationDomainService) extractContent(chunks []entity.KnowledgeChunk) []string {
	contents := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		contents = append(contents, chunk.Content())
	}
	return contents
}

// buildLLMMessages 构建LLM消息列表
func (s *ConversationDomainService) buildLLMMessages(
	messages []*entity.Message,
	userMessage string,
	knowledgeContext []string,
	systemPrompt string,
) []map[string]string {
	llmMessages := []map[string]string{}

	// 系统提示
	if systemPrompt != "" {
		llmMessages = append(llmMessages, map[string]string{
			"role":    "system",
			"content": s.buildSystemPrompt(systemPrompt, knowledgeContext),
		})
	}

	// 历史消息
	for _, msg := range messages {
		llmMessages = append(llmMessages, map[string]string{
			"role":    string(msg.Role()),
			"content": msg.Content(),
		})
	}

	// 当前用户消息
	llmMessages = append(llmMessages, map[string]string{
		"role":    "user",
		"content": userMessage,
	})

	return llmMessages
}

// buildSystemPrompt 构建系统提示（包含知识库上下文）
func (s *ConversationDomainService) buildSystemPrompt(basePrompt string, knowledgeContext []string) string {
	if len(knowledgeContext) == 0 {
		return basePrompt
	}

	knowledgeSection := "\n\n=== 知识库内容 ===\n"
	for i, ctx := range knowledgeContext {
		knowledgeSection += fmt.Sprintf("[知识%d] %s\n", i+1, ctx)
	}
	knowledgeSection += "\n请基于上述知识库内容回答用户问题。如果知识库中没有相关信息，请基于你的训练数据回答。\n"

	return basePrompt + knowledgeSection
}

// callLLM 调用LLM接口（实现在基础设施层）
func (s *ConversationDomainService) callLLM(
	ctx context.Context,
	bot *entity.Bot,
	messages []map[string]string,
) (string, int64, error) {
	// 这里调用基础设施层的LLM服务
	// 返回：响应内容、使用的Token数、错误

	// 简化实现：返回模拟数据
	return "这是一个模拟的LLM响应", 100, nil
}

// estimateTokens 估算Token使用量
func (s *ConversationDomainService) estimateTokens(text string, overhead int) int64 {
	// 粗略估算：中文约1.5字符=1token，英文约4字符=1token
	// 这里使用简化估算：1字符=0.5token
	tokens := int64(float64(len(text))*0.5) + int64(overhead)
	return tokens
}

// ============ 辅助方法 ============

// KnowledgeChunk 知识块（简化定义）
type KnowledgeChunk struct {
	id      string
	content string
	score   float64
}

func (kc *KnowledgeChunk) ID() string        { return kc.id }
func (kc *KnowledgeChunk) Content() string    { return kc.content }
func (kc *KnowledgeChunk) Score() float64     { return kc.score }

// Message 消息实体（简化定义）
type Message struct {
	id             string
	conversationID string
	role           string
	content        string
	userID         string
	tokenCount     int64
	metadata       map[string]interface{}
}

func (m *Message) ID() string           { return m.id }
func (m *Message) Role() string         { return m.role }
func (m *Message) Content() string      { return m.content }
func (m *Message) SetTokenCount(count int64) {
	m.tokenCount = count
}
func (m *Message) AddMetadata(key string, value interface{}) {
	if m.metadata == nil {
		m.metadata = make(map[string]interface{})
	}
	m.metadata[key] = value
}

// Conversation 对话实体（简化定义）
type Conversation struct {
	id           string
	tenantID     string
	botID        string
	userID       string
	messageCount int
	lastMsgAt    *time.Time
}

func (c *Conversation) ID() string { return c.id }
func (c *Conversation) UpdateLastMessageAt() {
	now := time.Now()
	c.lastMsgAt = &now
}
func (c *Conversation) IncrementMessageCount() {
	c.messageCount++
}

// Bot Bot实体（简化定义）
type Bot struct {
	id              string
	knowledgeBaseID string
	systemPrompt    string
	published       bool
}

func (b *Bot) ID() string { return b.id }
func (b *Bot) IsPublished() bool { return b.published }
func (b *Bot) UseTokens(tokens int64) error { return nil }

// 为了避免导入错误，这里使用简化的类型定义
// 实际项目中应该导入正确的包

import (
	"context"
	"time"
)

type TenantID string

type BotStatus string

const (
	BotStatusPublished BotStatus = "published"
	BotStatusDraft     BotStatus = "draft"
)

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)
