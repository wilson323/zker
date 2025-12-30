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

package fixtures

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/infra/saga"
)

// ==================== Mock Bot Repository ====================

// MockBotRepository Mock Bot仓储
// 职责: 提供Mock的Bot数据操作，避免真实数据库调用
type MockBotRepository struct {
	mu   sync.RWMutex
	bots map[string]*Bot
}

// NewMockBotRepository 创建Mock Bot仓储
func NewMockBotRepository() *MockBotRepository {
	return &MockBotRepository{
		bots: make(map[string]*Bot),
	}
}

// Create 创建Bot
func (r *MockBotRepository) Create(ctx context.Context, bot *Bot) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bots[bot.BotID]; exists {
		return fmt.Errorf("bot already exists: %s", bot.BotID)
	}

	r.bots[bot.BotID] = bot
	return nil
}

// GetByID 根据ID获取Bot
func (r *MockBotRepository) GetByID(ctx context.Context, botID string) (*Bot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bot, exists := r.bots[botID]
	if !exists {
		return nil, nil
	}
	return bot, nil
}

// Delete 删除Bot
func (r *MockBotRepository) Delete(ctx context.Context, botID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.bots, botID)
	return nil
}

// Clear 清空所有Bot
func (r *MockBotRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bots = make(map[string]*Bot)
}

// ==================== Mock Knowledge Repository ====================

// MockKnowledgeRepository Mock 知识库仓储
// 职责: 提供Mock的知识库数据操作
type MockKnowledgeRepository struct {
	mu         sync.RWMutex
	knowledge  map[string]*Knowledge
}

// NewMockKnowledgeRepository 创建Mock知识库仓储
func NewMockKnowledgeRepository() *MockKnowledgeRepository {
	return &MockKnowledgeRepository{
		knowledge: make(map[string]*Knowledge),
	}
}

// Create 创建知识库
func (r *MockKnowledgeRepository) Create(ctx context.Context, knowledge *Knowledge) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.knowledge[knowledge.KnowledgeID]; exists {
		return fmt.Errorf("knowledge already exists: %s", knowledge.KnowledgeID)
	}

	r.knowledge[knowledge.KnowledgeID] = knowledge
	return nil
}

// GetByID 根据ID获取知识库
func (r *MockKnowledgeRepository) GetByID(ctx context.Context, knowledgeID string) (*Knowledge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	knowledge, exists := r.knowledge[knowledgeID]
	if !exists {
		return nil, nil
	}
	return knowledge, nil
}

// Delete 删除知识库
func (r *MockKnowledgeRepository) Delete(ctx context.Context, knowledgeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.knowledge, knowledgeID)
	return nil
}

// Clear 清空所有知识库
func (r *MockKnowledgeRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.knowledge = make(map[string]*Knowledge)
}

// ==================== Failing Mock Knowledge Repository ====================

// FailingMockKnowledgeRepository 总是失败的知识库仓储
// 职责: 模拟知识库创建失败的场景
type FailingMockKnowledgeRepository struct {
	*MockKnowledgeRepository
}

// NewFailingMockKnowledgeRepository 创建总是失败的知识库仓储
func NewFailingMockKnowledgeRepository() *FailingMockKnowledgeRepository {
	return &FailingMockKnowledgeRepository{
		MockKnowledgeRepository: NewMockKnowledgeRepository(),
	}
}

// Create 创建知识库（总是失败）
func (r *FailingMockKnowledgeRepository) Create(ctx context.Context, knowledge *Knowledge) error {
	return fmt.Errorf("knowledge creation failed")
}

// ==================== Retry Mock Knowledge Repository ====================

// RetryMockKnowledgeRepository 重试后成功的知识库仓储
// 职责: 模拟前N次失败，第N+1次成功的场景
type RetryMockKnowledgeRepository struct {
	*MockKnowledgeRepository
	callCounts map[string]int
	mu         sync.Mutex
	maxRetries int
}

// NewRetryMockKnowledgeRepository 创建重试仓储
func NewRetryMockKnowledgeRepository(successOnAttempt int) *RetryMockKnowledgeRepository {
	return &RetryMockKnowledgeRepository{
		MockKnowledgeRepository: NewMockKnowledgeRepository(),
		callCounts:              make(map[string]int),
		maxRetries:              successOnAttempt,
	}
}

// Create 创建知识库（前N次失败，第N+1次成功）
func (r *RetryMockKnowledgeRepository) Create(ctx context.Context, knowledge *Knowledge) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	count := r.callCounts[knowledge.KnowledgeID]
	count++
	r.callCounts[knowledge.KnowledgeID] = count

	if count < r.maxRetries {
		return fmt.Errorf("knowledge creation failed (attempt %d/%d)", count, r.maxRetries)
	}

	// 成功
	return r.MockKnowledgeRepository.Create(ctx, knowledge)
}

// GetCreateCallCount 获取创建调用次数
func (r *RetryMockKnowledgeRepository) GetCreateCallCount(knowledgeID string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.callCounts[knowledgeID]
}

// ==================== Timeout Mock Knowledge Repository ====================

// TimeoutMockKnowledgeRepository 超时的知识库仓储
// 职责: 模拟操作超时的场景
type TimeoutMockKnowledgeRepository struct {
	*MockKnowledgeRepository
}

// NewTimeoutMockKnowledgeRepository 创建超时仓储
func NewTimeoutMockKnowledgeRepository() *TimeoutMockKnowledgeRepository {
	return &TimeoutMockKnowledgeRepository{
		MockKnowledgeRepository: NewMockKnowledgeRepository(),
	}
}

// Create 创建知识库（超时）
func (r *TimeoutMockKnowledgeRepository) Create(ctx context.Context, knowledge *Knowledge) error {
	// 等待上下文取消
	<-ctx.Done()
	return ctx.Err()
}

// ==================== Mock LLM Client ====================

// MockLLMClient Mock LLM客户端
// 职责: 提供Mock的LLM响应，避免真实API调用
type MockLLMClient struct {
	mu         sync.RWMutex
	confidence float64
	responses  map[string]string
	entities   []map[string]interface{}
}

// NewMockLLMClient 创建Mock LLM客户端
func NewMockLLMClient() *MockLLMClient {
	return &MockLLMClient{
		confidence: 0.8,
		responses:  make(map[string]string),
		entities:   make([]map[string]interface{}, 0),
	}
}

// Chat 聊天
func (m *MockLLMClient) Chat(ctx context.Context, prompt string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回预设响应或默认响应
	if response, ok := m.responses[prompt]; ok {
		return response, nil
	}

	return "Mock LLM response", nil
}

// SetConfidence 设置AI置信度
func (m *MockLLMClient) SetConfidence(confidence float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.confidence = confidence
}

// GetConfidence 获取AI置信度
func (m *MockLLMClient) GetConfidence() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.confidence
}

// SetResponse 设置预设响应
func (m *MockLLMClient) SetResponse(prompt, response string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[prompt] = response
}

// SetEntityExtraction 设置实体提取结果
func (m *MockLLMClient) SetEntityExtraction(entities []map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entities = entities
}

// GetEntityExtraction 获取实体提取结果
func (m *MockLLMClient) GetEntityExtraction() []map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.entities
}

// Clear 清空数据
func (m *MockLLMClient) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = make(map[string]string)
	m.entities = make([]map[string]interface{}, 0)
}

// ==================== Mock Milvus Client ====================

// MockMilvusClient Mock Milvus客户端
// 职责: 提供Mock的向量数据库操作
type MockMilvusClient struct {
	mu      sync.RWMutex
	vectors map[string][]float32
}

// NewMockMilvusClient 创建Mock Milvus客户端
func NewMockMilvusClient() *MockMilvusClient {
	return &MockMilvusClient{
		vectors: make(map[string][]float32),
	}
}

// Insert 插入向量
func (m *MockMilvusClient) Insert(ctx context.Context, collectionName string, ids []string, vectors [][]float32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, id := range ids {
		m.vectors[id] = vectors[i]
	}
	return nil
}

// Search 搜索向量
func (m *MockMilvusClient) Search(ctx context.Context, collectionName string, vector []float32, topK int) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 简单返回所有ID
	result := make([]string, 0, len(m.vectors))
	for id := range m.vectors {
		result = append(result, id)
		if len(result) >= topK {
			break
		}
	}
	return result, nil
}

// Delete 删除向量
func (m *MockMilvusClient) Delete(ctx context.Context, collectionName string, ids []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, id := range ids {
		delete(m.vectors, id)
	}
	return nil
}

// Clear 清空向量
func (m *MockMilvusClient) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.vectors = make(map[string][]float32)
}

// ==================== Saga步骤实现 ====================

// CreateBotStep 创建Bot步骤
type CreateBotStep struct {
	BotRepo   *MockBotRepository
	BotID     string
	BotName   string
	BotDesc   string
	TenantID  string
	CreatorID string
}

// Execute 执行步骤
func (s *CreateBotStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	bot := &Bot{
		BotID:     s.BotID,
		TenantID:  s.TenantID,
		Name:      s.BotName,
		Status:    "draft",
		CreatedAt: time.Now(),
	}

	err := s.BotRepo.Create(ctx, bot)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

// Name 步骤名称
func (s *CreateBotStep) Name() string {
	return "CreateBot"
}

// Timeout 超时时间
func (s *CreateBotStep) Timeout() time.Duration {
	return 10 * time.Second
}

// DeleteBotStep 删除Bot补偿步骤
type DeleteBotStep struct {
	BotRepo *MockBotRepository
}

// Compensate 执行补偿
func (s *DeleteBotStep) Compensate(ctx context.Context, data interface{}) error {
	bot, ok := data.(*Bot)
	if !ok || bot == nil {
		return nil // 没有数据需要补偿
	}
	return s.BotRepo.Delete(ctx, bot.BotID)
}

// Name 补偿步骤名称
func (s *DeleteBotStep) Name() string {
	return "DeleteBot"
}

// Timeout 超时时间
func (s *DeleteBotStep) Timeout() time.Duration {
	return 5 * time.Second
}

// CreateKnowledgeStep 创建知识库步骤
type CreateKnowledgeStep struct {
	KnowledgeRepo *MockKnowledgeRepository
	KnowledgeID   string
	KnowledgeName string
	BotID         string
	TenantID      string
}

// Execute 执行步骤
func (s *CreateKnowledgeStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	knowledge := &Knowledge{
		KnowledgeID: s.KnowledgeID,
		TenantID:    s.TenantID,
		BotID:       s.BotID,
		Name:        s.KnowledgeName,
		Status:      "active",
		CreatedAt:   time.Now(),
	}

	err := s.KnowledgeRepo.Create(ctx, knowledge)
	if err != nil {
		return nil, err
	}

	return knowledge, nil
}

// Name 步骤名称
func (s *CreateKnowledgeStep) Name() string {
	return "CreateKnowledge"
}

// Timeout 超时时间
func (s *CreateKnowledgeStep) Timeout() time.Duration {
	return 10 * time.Second
}

// DeleteKnowledgeStep 删除知识库补偿步骤
type DeleteKnowledgeStep struct {
	KnowledgeRepo *MockKnowledgeRepository
}

// Compensate 执行补偿
func (s *DeleteKnowledgeStep) Compensate(ctx context.Context, data interface{}) error {
	knowledge, ok := data.(*Knowledge)
	if !ok || knowledge == nil {
		return nil
	}
	return s.KnowledgeRepo.Delete(ctx, knowledge.KnowledgeID)
}

// Name 补偿步骤名称
func (s *DeleteKnowledgeStep) Name() string {
	return "DeleteKnowledge"
}

// Timeout 超时时间
func (s *DeleteKnowledgeStep) Timeout() time.Duration {
	return 5 * time.Second
}

// ==================== Mock人机协同服务 ====================

// MockCollaborationOrchestrator Mock人机协同编排器
// 职责: 提供Mock的人机协同功能
type MockCollaborationOrchestrator struct {
	mu    sync.RWMutex
	db    interface{} // 可以是*sql.DB或其他类型
	tasks map[string]*ReviewTask
	logger *zap.Logger
}

// ReviewTask 审核任务
type ReviewTask struct {
	TaskID        string
	ResourceType  string
	ResourceID    string
	TenantID      string
	Status        string
	Reason        string
	Comment       string
	Priority      string
	Version       int
	CreatedAt     time.Time
}

// NewMockCollaborationOrchestrator 创建Mock人机协同编排器
func NewMockCollaborationOrchestrator(db interface{}, logger *zap.Logger) *MockCollaborationOrchestrator {
	return &MockCollaborationOrchestrator{
		db:     db,
		tasks:  make(map[string]*ReviewTask),
		logger: logger,
	}
}

// CreateTask 创建审核任务
func (m *MockCollaborationOrchestrator) CreateTask(ctx context.Context, task *ReviewTask) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task.TaskID = fmt.Sprintf("task-%d", time.Now().UnixNano())
	task.Status = "pending"
	task.CreatedAt = time.Now()
	task.Version = 1

	m.tasks[task.TaskID] = task
	return nil
}

// ListTasks 列出审核任务
func (m *MockCollaborationOrchestrator) ListTasks(ctx context.Context, filter interface{}) ([]*ReviewTask, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*ReviewTask, 0, len(m.tasks))
	for _, task := range m.tasks {
		result = append(result, task)
	}
	return result, nil
}

// GetTask 获取审核任务
func (m *MockCollaborationOrchestrator) GetTask(ctx context.Context, taskID string) (*ReviewTask, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	return task, nil
}

// SubmitReview 提交审核结果
func (m *MockCollaborationOrchestrator) SubmitReview(ctx context.Context, result interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 简化实现，实际需要解析result
	// task.Status = result.Decision
	// task.Comment = result.Comment
	return nil
}

// Clear 清空所有任务
func (m *MockCollaborationOrchestrator) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tasks = make(map[string]*ReviewTask)
}

// ==================== Mock对话服务 ====================

// MockConversationService Mock对话服务
// 职责: 提供Mock的对话功能
type MockConversationService struct {
	mu            sync.RWMutex
	db            interface{}
	conversations map[string]*Conversation
}

// NewMockConversationService 创建Mock对话服务
func NewMockConversationService(db interface{}) *MockConversationService {
	return &MockConversationService{
		db:            db,
		conversations: make(map[string]*Conversation),
	}
}

// CreateConversation 创建会话
func (s *MockConversationService) CreateConversation(ctx context.Context, conv *Conversation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.conversations[conv.ConversationID] = conv
	return nil
}

// Chat 对话
func (s *MockConversationService) Chat(ctx context.Context, req interface{}) (interface{}, error) {
	return &ChatResponse{
		Content: "Mock response",
		Status:  "success",
	}, nil
}

// ChatResponse 对话响应
type ChatResponse struct {
	Content string
	Status  string
}

// ==================== Mock Bot Service ====================

// MockBotService Mock Bot服务
// 职责: 提供Mock的Bot操作功能
type MockBotService struct {
	mu              sync.RWMutex
	db              interface{}
	collabService   *MockCollaborationOrchestrator
	llmClient       *MockLLMClient
	bots            map[string]*BotEntity
}

// BotEntity Bot实体
type BotEntity struct {
	BotID       string
	TenantID    string
	Name        string
	Description string
	Status      string
	IsComplex   bool
	Complexity  string
	AIConfidence float64
}

// NewMockBotService 创建Mock Bot服务
func NewMockBotService(db interface{}, collabService *MockCollaborationOrchestrator, llmClient *MockLLMClient) *MockBotService {
	return &MockBotService{
		db:            db,
		collabService: collabService,
		llmClient:     llmClient,
		bots:          make(map[string]*BotEntity),
	}
}

// CreateBotWithReview 创建Bot（带审核）
func (s *MockBotService) CreateBotWithReview(ctx context.Context, req interface{}) (*BotEntity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 简化实现，实际需要解析req
	bot := &BotEntity{
		BotID:    fmt.Sprintf("bot-%d", time.Now().UnixNano()),
		Status:   "draft",
	}

	s.bots[bot.BotID] = bot
	return bot, nil
}

// GetBotByID 获取Bot
func (s *MockBotService) GetBotByID(ctx context.Context, botID string) (*BotEntity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bot, exists := s.bots[botID]
	if !exists {
		return nil, fmt.Errorf("bot not found: %s", botID)
	}
	return bot, nil
}

// UpdateBotAndResubmit 更新Bot并重新提交审核
func (s *MockBotService) UpdateBotAndResubmit(ctx context.Context, req interface{}) (*BotEntity, error) {
	// 简化实现
	return nil, nil
}
