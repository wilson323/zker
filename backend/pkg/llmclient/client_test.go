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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				APIKey:  "test-api-key",
				BaseURL: "https://api.example.com",
				Model:   "gpt-3.5-turbo",
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
			errMsg:  "config is required",
		},
		{
			name: "empty api key",
			config: &Config{
				APIKey:  "",
				BaseURL: "https://api.example.com",
				Model:   "gpt-3.5-turbo",
			},
			wantErr: true,
			errMsg:  "api_key is required",
		},
		{
			name: "empty base url",
			config: &Config{
				APIKey:  "test-api-key",
				BaseURL: "",
				Model:   "gpt-3.5-turbo",
			},
			wantErr: true,
			errMsg:  "base_url is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, client)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
				defer client.Close()
			}
		})
	}
}

func TestChat(t *testing.T) {
	config := &Config{
		APIKey:  "test-api-key",
		BaseURL: "https://api.example.com",
		Model:   "gpt-3.5-turbo",
	}

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	req := &ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp, err := client.Chat(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// 验证响应字段
	assert.NotEmpty(t, resp.ID)
	assert.NotEmpty(t, resp.Object)
	assert.NotEmpty(t, resp.Model)
	assert.NotEmpty(t, resp.Content) // 验证便捷字段
	assert.Greater(t, resp.Created, int64(0))

	// 验证 Choices
	assert.NotEmpty(t, resp.Choices)
	choice := resp.Choices[0]
	assert.Equal(t, 0, choice.Index)
	assert.Equal(t, "assistant", choice.Message.Role)
	assert.NotEmpty(t, choice.Message.Content)
	assert.Equal(t, "stop", choice.FinishReason)

	// 验证 Usage
	assert.Greater(t, resp.Usage.PromptTokens, 0)
	assert.Greater(t, resp.Usage.CompletionTokens, 0)
	assert.Greater(t, resp.Usage.TotalTokens, 0)
}

func TestChatWithMultipleMessages(t *testing.T) {
	config := &Config{
		APIKey:  "test-api-key",
		BaseURL: "https://api.example.com",
		Model:   "gpt-3.5-turbo",
	}

	client, err := NewClient(config)
	require.NoError(t, err)
	defer client.Close()

	req := &ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*Message{
			{Role: "system", Content: "You are a helpful assistant"},
			{Role: "user", Content: "What is the capital of France?"},
		},
		Temperature: 0.7,
		MaxTokens:   100,
	}

	resp, err := client.Chat(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.NotEmpty(t, resp.Choices)
	assert.Equal(t, 1, len(resp.Choices))
}

func TestChatStream(t *testing.T) {
	config := &Config{
		APIKey:  "test-api-key",
		BaseURL: "https://api.example.com",
		Model:   "gpt-3.5-turbo",
	}

	client, err := NewClient(config)
	require.NoError(t, err)
	defer client.Close()

	req := &ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*Message{
			{Role: "user", Content: "Hello stream"},
		},
	}

	chunkChan, err := client.ChatStream(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, chunkChan)

	// 读取流式响应
	chunkCount := 0
	for chunk := range chunkChan {
		chunkCount++
		assert.NotNil(t, chunk)
		assert.NotEmpty(t, chunk.ID)
		assert.NotEmpty(t, chunk.Model)

		if len(chunk.Choices) > 0 {
			choice := chunk.Choices[0]
			assert.NotNil(t, choice.Delta)
		}
	}

	assert.Greater(t, chunkCount, 0, "应该至少收到一个chunk")
}

func TestClose(t *testing.T) {
	config := &Config{
		APIKey:  "test-api-key",
		BaseURL: "https://api.example.com",
		Model:   "gpt-3.5-turbo",
	}

	client, err := NewClient(config)
	require.NoError(t, err)

	// 关闭客户端
	err = client.Close()
	assert.NoError(t, err)

	// 再次关闭应该不会报错
	err = client.Close()
	assert.NoError(t, err)
}
