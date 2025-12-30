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
	"database/sql"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/application/memory"
)

// ==================== Mock对话服务（带记忆功能） ====================

// MockConversationServiceWithMemory Mock对话服务（带记忆）
// 职责: 提供带记忆功能的Mock对话服务
type MockConversationServiceWithMemory struct {
	mu            sync.RWMutex
	db            *sql.DB
	memoryService memory.ConversationMemoryService
	llmClient     *MockLLMClient
	conversations map[string]*Conversation
	messages      map[string][]*Message
}

// Message 消息
type Message struct {
	MessageID   string
	ConversationID string
	Role        string
	Content     string
	CreatedAt   time.Time
}

// NewMockConversationServiceWithMemory 创建带记忆的Mock对话服务
func NewMockConversationServiceWithMemory(db *sql.DB, memSvc memory.ConversationMemoryService, llmClient *MockLLMClient) *MockConversationServiceWithMemory {
	return &MockConversationServiceWithMemory{
		db:            db,
		memoryService: memSvc,
		llmClient:     llmClient,
		conversations: make(map[string]*Conversation),
		messages:      make(map[string][]*Message),
	}
}

// CreateConversation 创建会话
func (s *MockConversationServiceWithMemory) CreateConversation(ctx context.Context, conv *Conversation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.conversations[conv.ConversationID] = conv
	s.messages[conv.ConversationID] = make([]*Message, 0)
	return nil
}

// Chat 对话（带记忆注入和提取）
func (s *MockConversationServiceWithMemory) Chat(ctx context.Context, req interface{}) (*ChatResponse, error) {
	// 简化实现
	reqData, ok := req.(interface {
		GetConversationID() string
		GetUserID() string
		GetMessage() string
		GetEnableMemory() bool
	})

	if !ok {
		return &ChatResponse{
			Content: "Mock response",
			Status:  "success",
		}, nil
	}

	conversationID := reqData.GetConversationID()
	userID := reqData.GetUserID()
	message := reqData.GetMessage()
	enableMemory := reqData.GetEnableMemory()

	// 1. 如果启用记忆，检索相关记忆
	var memories string
	if enableMemory {
		// 这里简化为直接设置，实际应该调用memoryService
		memories = "用户喜欢简洁的回答"
	}

	// 2. 构造带记忆的Prompt
	prompt := message
	if memories != "" {
		prompt = fmt.Sprintf("[用户记忆: %s]\n用户问题: %s", memories, message)
	}

	// 3. 调用LLM生成响应
	response, err := s.llmClient.Chat(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// 4. 保存消息记录
	s.mu.Lock()
	s.messages[conversationID] = append(s.messages[conversationID], &Message{
		MessageID:   fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		ConversationID: conversationID,
		Role:        "user",
		Content:     message,
		CreatedAt:   time.Now(),
	})
	s.mu.Unlock()

	// 5. 如果启用记忆，异步提取新记忆
	if enableMemory {
		go func() {
			// 提取实体
			entities := s.llmClient.GetEntityExtraction()
			for _, entity := range entities {
				// 存储记忆
				// s.memoryService.StoreMemory(ctx, ...)
			}
		}()
	}

	return &ChatResponse{
		Content: response,
		Status:  "success",
	}, nil
}

// ==================== Mock Memory Service ====================

// MockMemoryService Mock记忆服务
// 职责: 提供Mock的记忆操作功能
type MockMemoryService struct {
	mu       sync.RWMutex
	memories map[string][]*MockMemory
}

// MockMemory 记忆
type MockMemory struct {
	MemoryID     string
	UserID       string
	TenantID     string
	MemoryType   string
	Content      string
	Importance   float64
	LastAccessed time.Time
	ExpiresAt    time.Time
}

// NewMockMemoryService 创建Mock记忆服务
func NewMockMemoryService() *MockMemoryService {
	return &MockMemoryService{
		memories: make(map[string][]*MockMemory),
	}
}

// StoreMemory 存储记忆
func (s *MockMemoryService) StoreMemory(ctx context.Context, memory *MockMemory) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	userMemories := s.memories[memory.UserID]
	userMemories = append(userMemories, memory)
	s.memories[memory.UserID] = userMemories
	return nil
}

// RetrieveMemories 检索记忆
func (s *MockMemoryService) RetrieveMemories(ctx context.Context, userID, query string, limit int) ([]*MockMemory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	memories, exists := s.memories[userID]
	if !exists {
		return []*MockMemory{}, nil
	}

	// 简单返回所有记忆
	if len(memories) > limit {
		return memories[:limit], nil
	}
	return memories, nil
}

// CleanupExpiredMemories 清理过期记忆
func (s *MockMemoryService) CleanupExpiredMemories(ctx context.Context, userID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	memories, exists := s.memories[userID]
	if !exists {
		return 0, nil
	}

	count := 0
	now := time.Now()
	filtered := make([]*MockMemory, 0)

	for _, mem := range memories {
		if mem.ExpiresAt.Before(now) {
			count++
		} else {
			filtered = append(filtered, mem)
		}
	}

	s.memories[userID] = filtered
	return count, nil
}

// ApplyImportanceDecay 应用重要性衰减
func (s *MockMemoryService) ApplyImportanceDecay(ctx context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	memories, exists := s.memories[userID]
	if !exists {
		return nil
	}

	now := time.Now()
	for _, mem := range memories {
		// 根据最后访问时间计算衰减
		daysSinceAccess := int(now.Sub(mem.LastAccessed).Hours() / 24)
		decayFactor := 1.0 - float64(daysSinceAccess)*0.01 // 每天衰减1%

		if decayFactor < 0 {
			decayFactor = 0
		}

		mem.Importance = mem.Importance * decayFactor
	}

	return nil
}

// UpdateOrMergeMemory 更新或合并记忆
func (s *MockMemoryService) UpdateOrMergeMemory(ctx context.Context, userID, content, updatedContent string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	memories, exists := s.memories[userID]
	if !exists {
		return nil
	}

	// 查找相似记忆并更新
	for _, mem := range memories {
		if contains(mem.Content, content) {
			mem.Content = updatedContent
			mem.Importance = mem.Importance + 0.1 // 提升重要性
			return nil
		}
	}

	// 没找到，创建新记忆
	newMemory := &MockMemory{
		MemoryID:     fmt.Sprintf("memory-%d", time.Now().UnixNano()),
		UserID:       userID,
		Content:      updatedContent,
		Importance:   0.8,
		LastAccessed: time.Now(),
	}

	s.memories[userID] = append(memories, newMemory)
	return nil
}

// ==================== Mock权限服务 ====================

// MockAuthorizationService Mock权限服务
// 职责: 提供Mock的权限检查功能
type MockAuthorizationService struct {
	mu          sync.RWMutex
	permissions map[string]*MockDataPermission
}

// MockDataPermission 数据权限
type MockDataPermission struct {
	UserID       string
	TenantID     string
	ResourceType string
	Action       string
	Allowed      bool
}

// NewMockAuthorizationService 创建Mock权限服务
func NewMockAuthorizationService() *MockAuthorizationService {
	return &MockAuthorizationService{
		permissions: make(map[string]*MockDataPermission),
	}
}

// CheckAuthz 检查权限
func (s *MockAuthorizationService) CheckAuthz(ctx context.Context, req interface{}) (interface{}, error) {
	// 简化实现，总是允许
	return map[string]interface{}{
		"decision": "allow",
	}, nil
}

// BuildDataPermissionFilter 构建数据权限过滤器
func (s *MockAuthorizationService) BuildDataPermissionFilter(ctx context.Context, req interface{}) (*MockDataPermissionFilter, error) {
	// 简化实现
	return &MockDataPermissionFilter{
		TenantIDFilter: "tenant-001",
	}, nil
}

// MockDataPermissionFilter 数据权限过滤器
type MockDataPermissionFilter struct {
	TenantIDFilter string
}

// ==================== Mock角色服务 ====================

// MockRoleService Mock角色服务
// 职责: 提供Mock的角色管理功能
type MockRoleService struct {
	mu    sync.RWMutex
	roles map[string]*MockRole
}

// MockRole 角色
type MockRole struct {
	RoleID      string
	TenantID    string
	Name        string
	Description string
}

// NewMockRoleService 创建Mock角色服务
func NewMockRoleService() *MockRoleService {
	return &MockRoleService{
		roles: make(map[string]*MockRole),
	}
}

// CreateRole 创建角色
func (s *MockRoleService) CreateRole(ctx context.Context, role *MockRole) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.roles[role.RoleID] = role
	return nil
}

// ListRoles 列出角色
func (s *MockRoleService) ListRoles(ctx context.Context, tenantID string) ([]*MockRole, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roles := make([]*MockRole, 0)
	for _, role := range s.roles {
		if role.TenantID == tenantID {
			roles = append(roles, role)
		}
	}
	return roles, nil
}

// GetRolePermissions 获取角色权限
func (s *MockRoleService) GetRolePermissions(ctx context.Context, roleID string) ([]interface{}, error) {
	// 简化实现
	return []interface{}{}, nil
}

// CheckFieldPermission 检查字段权限
func (s *MockRoleService) CheckFieldPermission(ctx context.Context, req interface{}) (bool, error) {
	// 简化实现，默认允许
	return true, nil
}

// ==================== 辅助接口 ====================

// ChatRequest 对话请求接口
type ChatRequest interface {
	GetConversationID() string
	GetUserID() string
	GetMessage() string
	GetEnableMemory() bool
}

// ==================== 类型辅助函数 ====================

// 包含检查
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || indexOf(s, substr) >= 0)
}

// 查找子串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
