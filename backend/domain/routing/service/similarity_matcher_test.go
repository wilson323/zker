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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockEmbeddingClient 向量嵌入客户端 Mock
type MockEmbeddingClient struct {
	mock.Mock
}

func (m *MockEmbeddingClient) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]float32), args.Error(1)
}

// MockBotRepository Bot仓储 Mock
type MockBotRepository struct {
	mock.Mock
}

func (m *MockBotRepository) ListByTenant(ctx context.Context, tenantID string, filter interface{}) ([]*BotInfo, error) {
	args := m.Called(ctx, tenantID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*BotInfo), args.Error(1)
}

// SimilarityMatcherTestSuite 测试套件
type SimilarityMatcherTestSuite struct {
	suite.Suite
	matcher             *SimilarityMatcher
	mockEmbeddingClient *MockEmbeddingClient
	mockBotRepo         *MockBotRepository
	ctx                 context.Context
}

func (s *SimilarityMatcherTestSuite) SetupTest() {
	s.mockEmbeddingClient = new(MockEmbeddingClient)
	s.mockBotRepo = new(MockBotRepository)
	s.matcher = NewSimilarityMatcher(s.mockBotRepo, s.mockEmbeddingClient, 0.7)
	s.ctx = context.Background()
}

// TestMatch_Success 测试匹配成功
func (s *SimilarityMatcherTestSuite) TestMatch_Success() {
	// Arrange
	tenantID := "tenant-123"
	userInput := "今天天气怎么样？"

	// 用户输入的向量
	userEmbedding := []float32{0.1, 0.2, 0.3, 0.4, 0.5}

	// Bot列表
	bots := []*BotInfo{
		{
			BotID:     "bot-1",
			Name:      "天气助手",
			Desc:      "专门回答天气相关问题",
			Embedding: []float32{0.11, 0.19, 0.31, 0.39, 0.51}, // 相似向量
			IsActive:  true,
		},
		{
			BotID:     "bot-2",
			Name:      "新闻助手",
			Desc:      "获取最新新闻",
			Embedding: []float32{0.9, 0.8, 0.7, 0.6, 0.5}, // 不相似向量
			IsActive:  true,
		},
		{
			BotID:     "bot-3",
			Name:      "未激活Bot",
			Desc:      "这个Bot未激活",
			Embedding: []float32{0.11, 0.19, 0.31, 0.39, 0.51},
			IsActive:  false, // 未激活
		},
		{
			BotID:    "bot-4",
			Name:     "无向量Bot",
			Desc:     "这个Bot没有向量",
			IsActive: true,
			// Embedding 为 nil
		},
	}

	s.mockEmbeddingClient.On("GetEmbedding", s.ctx, userInput).Return(userEmbedding, nil)
	s.mockBotRepo.On("ListByTenant", s.ctx, tenantID, nil).Return(bots, nil)

	input := &MatchInput{
		UserInput: userInput,
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	// bot-1应该匹配（高相似度），bot-2可能不匹配（低相似度），bot-3和bot-4不应该匹配
	assert.Greater(s.T(), len(results), 0)

	// 验证bot-1在结果中
	foundBot1 := false
	for _, r := range results {
		if r.BotID == "bot-1" {
			foundBot1 = true
			assert.Equal(s.T(), "similarity", r.MatchType)
		}
	}
	assert.True(s.T(), foundBot1, "bot-1应该匹配成功")

	s.mockEmbeddingClient.AssertExpectations(s.T())
	s.mockBotRepo.AssertExpectations(s.T())
}

// TestMatch_EmbeddingError 测试获取向量失败
func (s *SimilarityMatcherTestSuite) TestMatch_EmbeddingError() {
	// Arrange
	tenantID := "tenant-123"

	s.mockEmbeddingClient.On("GetEmbedding", s.ctx, "test input").Return(
		([]float32)(nil), errors.New("embedding service error"))

	input := &MatchInput{
		UserInput: "test input",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), results)
	s.mockEmbeddingClient.AssertExpectations(s.T())
}

// TestMatch_BotRepositoryError 测试Bot仓储错误
func (s *SimilarityMatcherTestSuite) TestMatch_BotRepositoryError() {
	// Arrange
	tenantID := "tenant-123"
	userEmbedding := []float32{0.1, 0.2, 0.3}

	s.mockEmbeddingClient.On("GetEmbedding", s.ctx, "test input").Return(userEmbedding, nil)
	s.mockBotRepo.On("ListByTenant", s.ctx, tenantID, nil).Return(
		([]*BotInfo)(nil), errors.New("database error"))

	input := &MatchInput{
		UserInput: "test input",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), results)
	s.mockEmbeddingRepo.AssertExpectations(s.T())
	s.mockBotRepo.AssertExpectations(s.T())
}

// TestMatch_NoBots 测试没有Bot
func (s *SimilarityMatcherTestSuite) TestMatch_NoBots() {
	// Arrange
	tenantID := "tenant-123"
	userEmbedding := []float32{0.1, 0.2, 0.3}

	s.mockEmbeddingClient.On("GetEmbedding", s.ctx, "test input").Return(userEmbedding, nil)
	s.mockBotRepo.On("ListByTenant", s.ctx, tenantID, nil).Return([]*BotInfo{}, nil)

	input := &MatchInput{
		UserInput: "test input",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 0)
	s.mockEmbeddingClient.AssertExpectations(s.T())
	s.mockBotRepo.AssertExpectations(s.T())
}

// TestMatch_Top3Limit 测试Top-3限制
func (s *SimilarityMatcherTestSuite) TestMatch_Top3Limit() {
	// Arrange
	tenantID := "tenant-123"
	userInput := "test"

	userEmbedding := []float32{0.1, 0.2, 0.3}

	// 创建5个高相似度的Bot
	bots := make([]*BotInfo, 5)
	for i := 0; i < 5; i++ {
		bots[i] = &BotInfo{
			BotID:     "bot-" + string(rune('1'+i)),
			Name:      "Bot " + string(rune('1'+i)),
			Desc:      "Description",
			Embedding: []float32{0.11, 0.19, 0.31}, // 高相似度
			IsActive:  true,
		}
	}

	s.mockEmbeddingClient.On("GetEmbedding", s.ctx, userInput).Return(userEmbedding, nil)
	s.mockBotRepo.On("ListByTenant", s.ctx, tenantID, nil).Return(bots, nil)

	input := &MatchInput{
		UserInput: userInput,
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.LessOrEqual(s.T(), len(results), 3) // 最多返回3个
	s.mockEmbeddingClient.AssertExpectations(s.T())
	s.mockBotRepo.AssertExpectations(s.T())
}

// TestCosineSimilarity 测试余弦相似度计算
func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float64
	}{
		{
			name:     "Identical vectors",
			a:        []float32{1, 2, 3},
			b:        []float32{1, 2, 3},
			expected: 1.0,
		},
		{
			name:     "Orthogonal vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "Opposite vectors",
			a:        []float32{1, 2, 3},
			b:        []float32{-1, -2, -3},
			expected: -1.0,
		},
		{
			name:     "Different length",
			a:        []float32{1, 2, 3},
			b:        []float32{1, 2},
			expected: 0.0,
		},
		{
			name:     "Zero vector",
			a:        []float32{0, 0, 0},
			b:        []float32{1, 2, 3},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cosineSimilarity(tt.a, tt.b)
			assert.InDelta(t, tt.expected, result, 0.001)
		})
	}
}

// TestSimilarityMatcher_DefaultThreshold 测试默认阈值
func TestSimilarityMatcher_DefaultThreshold(t *testing.T) {
	mockBotRepo := new(MockBotRepository)
	mockEmbeddingClient := new(MockEmbeddingClient)

	// 设置阈值为0
	matcher := NewSimilarityMatcher(mockBotRepo, mockEmbeddingClient, 0)
	assert.Equal(t, float32(0.7), matcher.similarityThreshold, "应该使用默认阈值0.7")
}

// 运行测试套件
func TestSimilarityMatcherSuite(t *testing.T) {
	suite.Run(t, new(SimilarityMatcherTestSuite))
}
