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

package llmclient

import (
	"context"
	"fmt"
	"time"
)

// LLMClient LLM客户端接口
type LLMClient interface {
	// Chat 发送聊天请求
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// ChatStream 发送流式聊天请求
	ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error)

	// Close 关闭客户端
	Close() error
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string    `json:"model"`                  // 模型名称
	Messages    []*Message `json:"messages"`              // 消息列表
	Temperature float64   `json:"temperature,omitempty"`  // 温度参数
	MaxTokens   int       `json:"max_tokens,omitempty"`   // 最大token数
}

// Message 消息
type Message struct {
	Role    string `json:"role"`    // 角色：system/user/assistant
	Content string `json:"content"` // 内容
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Content string   `json:"content"` // 响应内容（便捷字段）
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatChunk 流式聊天响应块
type ChatChunk struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Created int64       `json:"created"`
	Model   string      `json:"model"`
	Choices []ChunkChoice `json:"choices"`
}

// ChunkChoice 流式选择
type ChunkChoice struct {
	Index        int           `json:"index"`
	Delta        MessageDelta  `json:"delta"`
	FinishReason *string       `json:"finish_reason"`
}

// MessageDelta 消息增量
type MessageDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// Config LLM客户端配置
type Config struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	Model   string `json:"model"`
}

// NewClient 创建LLM客户端（工厂函数）
func NewClient(config *Config) (LLMClient, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if config.APIKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base_url is required")
	}

	// TODO: 根据实际需求实现具体的客户端
	// 例如：OpenAI、Claude、本地模型等
	return &mockClient{config: config}, nil
}

// mockClient 模拟客户端（用于编译通过）
type mockClient struct {
	config *Config
}

func (m *mockClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 模拟实现
	content := "This is a mock response. Please implement the actual LLM client."
	if len(req.Messages) > 0 {
		content = fmt.Sprintf("Mock response to: %s", req.Messages[len(req.Messages)-1].Content)
	}

	return &ChatResponse{
		ID:      "mock-id",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Content: content, // 便捷字段
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}

func (m *mockClient) ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error) {
	// 模拟实现
	ch := make(chan *ChatChunk, 1)
	ch <- &ChatChunk{
		ID:      "mock-id",
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []ChunkChoice{
			{
				Index: 0,
				Delta: MessageDelta{
					Role:    "assistant",
					Content: "This is a mock stream response.",
				},
			},
		},
	}
	close(ch)
	return ch, nil
}

func (m *mockClient) Close() error {
	return nil
}
