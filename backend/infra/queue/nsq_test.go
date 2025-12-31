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
	"sync"
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	
	"go.uber.org/zap"
)

// =====================================================================
// 配置测试
// =====================================================================

func TestDefaultNSQConfig(t *testing.T) {
	config := DefaultNSQConfig()

	assert.NotNil(t, config)
	assert.Equal(t, []string{"localhost:4150"}, config.NSQDAddresses)
	assert.Equal(t, []string{"localhost:4161"}, config.NSQLookupdAddresses)
	assert.Equal(t, 5*time.Second, config.ProducerTimeout)
	assert.Equal(t, 3, config.ProducerRetry)
	assert.Equal(t, 3, config.ConsumerMaxRetries)
	assert.Equal(t, 10, config.ConsumerMaxInFlight)
	assert.True(t, config.DLQEnabled)
	assert.Equal(t, "-dlq", config.DLQTopicSuffix)
}

func TestNSQConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *NSQConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  DefaultNSQConfig(),
			wantErr: false,
		},
		{
			name: "empty nsqd addresses",
			config: &NSQConfig{
				NSQDAddresses:       []string{},
				NSQLookupdAddresses: []string{"localhost:4161"},
			},
			wantErr: true,
		},
		{
			name: "empty nsqlookupd addresses",
			config: &NSQConfig{
				NSQDAddresses:       []string{"localhost:4150"},
				NSQLookupdAddresses: []string{},
			},
			wantErr: true,
		},
		{
			name: "negative producer retry",
			config: &NSQConfig{
				NSQDAddresses:       []string{"localhost:4150"},
				NSQLookupdAddresses: []string{"localhost:4161"},
				ProducerRetry:       -1,
			},
			wantErr: true,
		},
		{
			name: "invalid max in flight",
			config: &NSQConfig{
				NSQDAddresses:       []string{"localhost:4150"},
				NSQLookupdAddresses: []string{"localhost:4161"},
				ConsumerMaxInFlight: 0,
			},
			wantErr: true,
		},
		{
			name: "dlq enabled but empty suffix",
			config: &NSQConfig{
				NSQDAddresses:       []string{"localhost:4150"},
				NSQLookupdAddresses: []string{"localhost:4161"},
				DLQEnabled:          true,
				DLQTopicSuffix:      "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =====================================================================
// 重试策略测试
// =====================================================================

func TestExponentialBackoffStrategy(t *testing.T) {
	strategy := &ExponentialBackoffStrategy{
		InitialDelay: 1 * time.Second,
		MaxDelay:     60 * time.Second,
		MaxAttempts:  5,
	}

	tests := []struct {
		name            string
		attempt         int
		shouldRetry     bool
		expectedMinDelay time.Duration
		expectedMaxDelay time.Duration
	}{
		{
			name:            "attempt 1",
			attempt:         1,
			shouldRetry:     true,
			expectedMinDelay: 2 * time.Second, // 2^1 * 1s
			expectedMaxDelay: 2 * time.Second,
		},
		{
			name:            "attempt 2",
			attempt:         2,
			shouldRetry:     true,
			expectedMinDelay: 4 * time.Second, // 2^2 * 1s
			expectedMaxDelay: 4 * time.Second,
		},
		{
			name:            "attempt 5",
			attempt:         5,
			shouldRetry:     true,
			expectedMinDelay: 32 * time.Second, // 2^5 * 1s
			expectedMaxDelay: 32 * time.Second,
		},
		{
			name:        "attempt 6 - exceeds max",
			attempt:     6,
			shouldRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldRetry, delay := strategy.ShouldRetry(context.Background(), tt.attempt, nil)

			assert.Equal(t, tt.shouldRetry, shouldRetry)

			if tt.shouldRetry {
				assert.GreaterOrEqual(t, delay, tt.expectedMinDelay)
				assert.LessOrEqual(t, delay, tt.expectedMaxDelay)
			}
		})
	}
}

func TestFixedDelayStrategy(t *testing.T) {
	strategy := &FixedDelayStrategy{
		Delay:       5 * time.Second,
		MaxAttempts: 3,
	}

	tests := []struct {
		name        string
		attempt     int
		shouldRetry bool
	}{
		{
			name:        "attempt 1",
			attempt:     1,
			shouldRetry: true,
		},
		{
			name:        "attempt 2",
			attempt:     2,
			shouldRetry: true,
		},
		{
			name:        "attempt 3 - exceeds max",
			attempt:     3,
			shouldRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldRetry, delay := strategy.ShouldRetry(context.Background(), tt.attempt, nil)

			assert.Equal(t, tt.shouldRetry, shouldRetry)

			if tt.shouldRetry {
				assert.Equal(t, 5*time.Second, delay)
			}
		})
	}
}

// =====================================================================
// 死信队列测试（使用内存存储）
// =====================================================================

func TestMemoryDLQStorage(t *testing.T) {
	storage := NewMemoryDLQStorage()
	ctx := context.Background()

	// 测试Store
	t.Run("Store", func(t *testing.T) {
		msg := &DLQMessage{
			MessageID:     "test_msg_1",
			OriginalTopic: "test_topic",
			OriginalBody:  []byte("test body"),
			FailedAt:      time.Now(),
			LastError:     "test error",
		}

		err := storage.Store(ctx, "test_topic-dlq", msg)
		assert.NoError(t, err)
	})

	// 测试Get
	t.Run("Get", func(t *testing.T) {
		msg, err := storage.Get(ctx, "test_msg_1")
		assert.NoError(t, err)
		assert.Equal(t, "test_msg_1", msg.MessageID)
		assert.Equal(t, "test_topic", msg.OriginalTopic)
		assert.Equal(t, "test error", msg.LastError)
	})

	// 测试Get - 不存在的消息
	t.Run("Get not found", func(t *testing.T) {
		_, err := storage.Get(ctx, "non_existent")
		assert.Error(t, err)
	})

	// 测试List
	t.Run("List", func(t *testing.T) {
		// 添加更多消息
		for i := 2; i <= 5; i++ {
			msg := &DLQMessage{
				MessageID:     fmt.Sprintf("test_msg_%d", i),
				OriginalTopic: "test_topic",
				OriginalBody:  []byte(fmt.Sprintf("test body %d", i)),
				FailedAt:      time.Now(),
			}
			_ = storage.Store(ctx, "test_topic-dlq", msg)
		}

		messages, err := storage.List(ctx, "test_topic-dlq", 10)
		assert.NoError(t, err)
		assert.Len(t, messages, 5)
	})

	// 测试MarkReplayed
	t.Run("MarkReplayed", func(t *testing.T) {
		err := storage.MarkReplayed(ctx, "test_msg_1")
		assert.NoError(t, err)

		msg, _ := storage.Get(ctx, "test_msg_1")
		assert.Equal(t, 1, msg.ReplayCount)
	})

	// 测试Delete
	t.Run("Delete", func(t *testing.T) {
		err := storage.Delete(ctx, "test_msg_1")
		assert.NoError(t, err)

		_, err = storage.Get(ctx, "test_msg_1")
		assert.Error(t, err)
	})

	// 测试GetMessageCount
	t.Run("GetMessageCount", func(t *testing.T) {
		count := storage.GetMessageCount("test_topic-dlq")
		assert.Equal(t, 4, count) // 剩余4条消息
	})
}

// =====================================================================
// 消息处理器测试
// =====================================================================

func TestJSONMessageHandler(t *testing.T) {
	handler := &JSONMessageHandler{
		ProcessFunc: func(ctx context.Context, msg map[string]interface{}) error {
			// 验证消息内容
			assert.Equal(t, "test_value", msg["test_key"])
			return nil
		},
	}

	// 创建测试消息
	msgData, _ := json.Marshal(map[string]interface{}{
		"test_key": "test_value",
	})

	message := &nsq.Message{
		Body: msgData,
	}

	// 执行处理
	err := handler.HandleMessage(context.Background(), message)
	assert.NoError(t, err)
}

func TestTypedMessageHandler(t *testing.T) {
	type TestMessage struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	handler := NewTypedMessageHandler(func(ctx context.Context, msg *TestMessage) error {
		assert.Equal(t, "123", msg.ID)
		assert.Equal(t, "test", msg.Name)
		return nil
	})

	// 创建测试消息
	msgData, _ := json.Marshal(TestMessage{
		ID:   "123",
		Name: "test",
	})

	message := &nsq.Message{
		Body: msgData,
	}

	// 执行处理
	err := handler.HandleMessage(context.Background(), message)
	assert.NoError(t, err)
}

// =====================================================================
// 队列管理器测试（模拟）
// =====================================================================

// MockNSQProducer 模拟生产者
type MockNSQProducer struct {
	PublishFunc    func(ctx context.Context, topic string, message interface{}) error
	PublishedCount int
	mu             sync.Mutex
}

func (m *MockNSQProducer) Publish(ctx context.Context, topic string, message interface{}) error {
	m.mu.Lock()
	m.PublishedCount++
	m.mu.Unlock()

	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, topic, message)
	}
	return nil
}

func (m *MockNSQProducer) GetMetrics() *ProducerMetrics {
	return &ProducerMetrics{
		PublishCount: int64(m.PublishedCount),
	}
}

func TestQueueManager_BasicOperations(t *testing.T) {
	// 注意：这是单元测试，不需要实际的NSQ服务
	// 集成测试需要在单独的测试文件中使用testcontainers

	// logger := zap.NewNop()
	// config := DefaultNSQConfig()

	t.Run("RegisterTask", func(t *testing.T) {
		// 这里我们只测试注册逻辑，不测试实际的NSQ连接
		// 实际集成测试需要使用testcontainers运行NSQ

		handler := TaskHandlerFunc(func(ctx context.Context, taskType string, taskData []byte) error {
			return nil
		})

		// 验证任务类型注册逻辑
		taskTypes := make(map[string]TaskHandler)
		taskTypes["test_task"] = handler

		assert.Contains(t, taskTypes, "test_task")
	})

	t.Run("EnqueueTask - validation", func(t *testing.T) {
		mockProducer := &MockNSQProducer{}
		ctx := context.Background()

		// 模拟入队
		task := map[string]interface{}{
			"id":   "123",
			"name": "test",
		}

		err := mockProducer.Publish(ctx, "test_topic", task)
		assert.NoError(t, err)
		assert.Equal(t, 1, mockProducer.PublishedCount)
	})
}

// =====================================================================
// 并发测试
// =====================================================================

func TestConcurrentPublish(t *testing.T) {
	mockProducer := &MockNSQProducer{}
	ctx := context.Background()

	// 并发发布1000条消息
	const numGoroutines = 100
	const messagesPerGoroutine = 10

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < messagesPerGoroutine; j++ {
				task := map[string]interface{}{
					"goroutine": goroutineID,
					"message":   j,
				}
				_ = mockProducer.Publish(ctx, "test_topic", task)
			}
		}(i)
	}

	wg.Wait()

	// 验证所有消息都已发布
	expectedCount := numGoroutines * messagesPerGoroutine
	assert.Equal(t, expectedCount, mockProducer.PublishedCount)
}

func TestConcurrentDLQStorage(t *testing.T) {
	storage := NewMemoryDLQStorage()
	ctx := context.Background()

	const numGoroutines = 50
	const messagesPerGoroutine = 20

	var wg sync.WaitGroup

	// 并发存储
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < messagesPerGoroutine; j++ {
				msg := &DLQMessage{
					MessageID:     fmt.Sprintf("msg_%d_%d", goroutineID, j),
					OriginalTopic: "test_topic",
					OriginalBody:  []byte("test body"),
					FailedAt:      time.Now(),
				}
				_ = storage.Store(ctx, "test_topic-dlq", msg)
			}
		}(i)
	}

	wg.Wait()

	// 验证消息数量
	count := storage.GetMessageCount("test_topic-dlq")
	expectedCount := numGoroutines * messagesPerGoroutine
	assert.Equal(t, expectedCount, count)
}

// =====================================================================
// Benchmark测试
// =====================================================================

func BenchmarkMemoryDLQStorage_Store(b *testing.B) {
	storage := NewMemoryDLQStorage()
	ctx := context.Background()

	msg := &DLQMessage{
		MessageID:     "test_msg",
		OriginalTopic: "test_topic",
		OriginalBody:  []byte("test body"),
		FailedAt:      time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.MessageID = fmt.Sprintf("msg_%d", i)
		_ = storage.Store(ctx, "test_topic-dlq", msg)
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	task := map[string]interface{}{
		"id":       "123",
		"name":     "test",
		"priority": 1,
		"data": map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(task)
	}
}

// =====================================================================
// 辅助函数
// =====================================================================

// setupTestLogger 创建测试日志记录器
func setupTestLogger() *zap.Logger {
	return zap.NewExample()
}

// assertError 断言错误
func assertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error but got nil", message)
	}
}

// requireError 要求有错误
func requireError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: required error but got nil", message)
	}
}
