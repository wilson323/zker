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

// HybridIntentMatcherTestSuite 测试套件
type HybridIntentMatcherTestSuite struct {
	suite.Suite
	matcher             *HybridIntentMatcher
	mockRuleMatcher     *MockRuleMatcherForHybrid
	mockSimilarityMatcher *MockSimilarityMatcherForHybrid
	ctx                 context.Context
}

// MockRuleMatcherForHybrid 规则匹配器 Mock（用于混合匹配器测试）
type MockRuleMatcherForHybrid struct {
	mock.Mock
}

func (m *MockRuleMatcherForHybrid) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MatchOutput), args.Error(1)
}

// MockSimilarityMatcherForHybrid 相似度匹配器 Mock（用于混合匹配器测试）
type MockSimilarityMatcherForHybrid struct {
	mock.Mock
}

func (m *MockSimilarityMatcherForHybrid) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MatchOutput), args.Error(1)
}

func (s *HybridIntentMatcherTestSuite) SetupTest() {
	s.mockRuleMatcher = new(MockRuleMatcherForHybrid)
	s.mockSimilarityMatcher = new(MockSimilarityMatcherForHybrid)
	s.matcher = NewHybridIntentMatcher(s.mockRuleMatcher, s.mockSimilarityMatcher)
	s.ctx = context.Background()
}

// TestMatch_BothSucceed 测试两个匹配器都成功
func (s *HybridIntentMatcherTestSuite) TestMatch_BothSucceed() {
	// Arrange
	tenantID := "tenant-123"
	input := &MatchInput{
		UserInput: "今天天气怎么样？",
		TenantID:  tenantID,
	}

	// 规则匹配器结果
	ruleResults := []*MatchOutput{
		{BotID: "bot-1", Confidence: 0.8, MatchType: "keyword"},
		{BotID: "bot-2", Confidence: 0.6, MatchType: "regex"},
	}

	// 相似度匹配器结果
	simResults := []*MatchOutput{
		{BotID: "bot-1", Confidence: 0.9, MatchType: "similarity"},
		{BotID: "bot-3", Confidence: 0.7, MatchType: "similarity"},
	}

	s.mockRuleMatcher.On("Match", s.ctx, input).Return(ruleResults, nil)
	s.mockSimilarityMatcher.On("Match", s.ctx, input).Return(simResults, nil)

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), results)

	// bot-1应该被聚合（两个匹配器都匹配到了）
	foundBot1 := false
	for _, r := range results {
		if r.BotID == "bot-1" {
			foundBot1 = true
			// 置信度应该被聚合（权重调整后）
			expectedConfidence := (0.8 * 0.5) + (0.9 * 0.5) // 默认权重各50%
			assert.InDelta(s.T(), expectedConfidence, r.Confidence, 0.01)
			// 匹配类型应该被合并
			assert.Contains(s.T(), r.MatchType, "keyword")
			assert.Contains(s.T(), r.MatchType, "similarity")
		}
	}
	assert.True(s.T(), foundBot1, "bot-1应该被聚合")

	s.mockRuleMatcher.AssertExpectations(s.T())
	s.mockSimilarityMatcher.AssertExpectations(s.T())
}

// TestMatch_RuleMatcherFails 测试规则匹配器失败
func (s *HybridIntentMatcherTestSuite) TestMatch_RuleMatcherFails() {
	// Arrange
	tenantID := "tenant-123"
	input := &MatchInput{
		UserInput: "今天天气怎么样？",
		TenantID:  tenantID,
	}

	simResults := []*MatchOutput{
		{BotID: "bot-1", Confidence: 0.9, MatchType: "similarity"},
	}

	s.mockRuleMatcher.On("Match", s.ctx, input).Return(
		([]*MatchOutput)(nil), errors.New("rule matcher error"))
	s.mockSimilarityMatcher.On("Match", s.ctx, input).Return(simResults, nil)

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert - 相似度匹配器成功，应该返回其结果
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 1)
	assert.Equal(s.T(), "bot-1", results[0].BotID)

	s.mockRuleMatcher.AssertExpectations(s.T())
	s.mockSimilarityMatcher.AssertExpectations(s.T())
}

// TestMatch_SimilarityMatcherFails 测试相似度匹配器失败
func (s *HybridIntentMatcherTestSuite) TestMatch_SimilarityMatcherFails() {
	// Arrange
	tenantID := "tenant-123"
	input := &MatchInput{
		UserInput: "今天天气怎么样？",
		TenantID:  tenantID,
	}

	ruleResults := []*MatchOutput{
		{BotID: "bot-1", Confidence: 0.8, MatchType: "keyword"},
	}

	s.mockRuleMatcher.On("Match", s.ctx, input).Return(ruleResults, nil)
	s.mockSimilarityMatcher.On("Match", s.ctx, input).Return(
		([]*MatchOutput)(nil), errors.New("similarity matcher error"))

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert - 规则匹配器成功，应该返回其结果
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 1)
	assert.Equal(s.T(), "bot-1", results[0].BotID)

	s.mockRuleMatcher.AssertExpectations(s.T())
	s.mockSimilarityMatcher.AssertExpectations(s.T())
}

// TestMatch_BothFail 测试两个匹配器都失败
func (s *HybridIntentMatcherTestSuite) TestMatch_BothFail() {
	// Arrange
	tenantID := "tenant-123"
	input := &MatchInput{
		UserInput: "今天天气怎么样？",
		TenantID:  tenantID,
	}

	s.mockRuleMatcher.On("Match", s.ctx, input).Return(
		([]*MatchOutput)(nil), errors.New("rule matcher error"))
	s.mockSimilarityMatcher.On("Match", s.ctx, input).Return(
		([]*MatchOutput)(nil), errors.New("similarity matcher error"))

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert - 两个都失败，应该返回空结果而不是错误
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 0)

	s.mockRuleMatcher.AssertExpectations(s.T())
	s.mockSimilarityMatcher.AssertExpectations(s.T())
}

// TestMatch_Aggregation 测试结果聚合
func (s *HybridIntentMatcherTestSuite) TestMatch_Aggregation() {
	// Arrange
	tenantID := "tenant-123"
	input := &MatchInput{
		UserInput: "今天天气怎么样？",
		TenantID:  tenantID,
	}

	// 规则匹配器返回bot-1
	ruleResults := []*MatchOutput{
		{BotID: "bot-1", Confidence: 0.5, MatchType: "keyword"},
	}

	// 相似度匹配器也返回bot-1
	simResults := []*MatchOutput{
		{BotID: "bot-1", Confidence: 0.5, MatchType: "similarity"},
	}

	s.mockRuleMatcher.On("Match", s.ctx, input).Return(ruleResults, nil)
	s.mockSimilarityMatcher.On("Match", s.ctx, input).Return(simResults, nil)

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 1) // bot-1应该被聚合成一个结果
	assert.Equal(s.T(), "bot-1", results[0].BotID)
	assert.InDelta(s.T(), 0.5, results[0].Confidence, 0.01) // 0.5*0.5 + 0.5*0.5 = 0.5

	s.mockRuleMatcher.AssertExpectations(s.T())
	s.mockSimilarityMatcher.AssertExpectations(s.T())
}

// TestMatch_Top3Limit 测试Top-3限制
func (s *HybridIntentMatcherTestSuite) TestMatch_Top3Limit() {
	// Arrange
	tenantID := "tenant-123"
	input := &MatchInput{
		UserInput: "今天天气怎么样？",
		TenantID:  tenantID,
	}

	// 规则匹配器返回5个Bot
	ruleResults := []*MatchOutput{
		{BotID: "bot-1", Confidence: 1.0, MatchType: "keyword"},
		{BotID: "bot-2", Confidence: 0.9, MatchType: "keyword"},
		{BotID: "bot-3", Confidence: 0.8, MatchType: "keyword"},
		{BotID: "bot-4", Confidence: 0.7, MatchType: "keyword"},
		{BotID: "bot-5", Confidence: 0.6, MatchType: "keyword"},
	}

	s.mockRuleMatcher.On("Match", s.ctx, input).Return(ruleResults, nil)
	s.mockSimilarityMatcher.On("Match", s.ctx, input).Return([]*MatchOutput{}, nil)

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.LessOrEqual(s.T(), len(results), 3) // 最多返回3个

	s.mockRuleMatcher.AssertExpectations(s.T())
	s.mockSimilarityMatcher.AssertExpectations(s.T())
}

// TestSetWeights 测试设置权重
func (s *HybridIntentMatcherTestSuite) TestSetWeights() {
	// 设置规则权重为70%，相似度权重为30%
	s.matcher.SetWeights(0.7, 0.3)

	assert.InDelta(s.T(), 0.7, s.matcher.ruleWeight, 0.01)
	assert.InDelta(s.T(), 0.3, s.matcher.similarityWeight, 0.01)
}

// TestSetWeights_Normalization 测试权重归一化
func (s *HybridIntentMatcherTestSuite) TestSetWeights_Normalization() {
	// 设置权重为200和100，应该被归一化为0.67和0.33
	s.matcher.SetWeights(200, 100)

	assert.InDelta(s.T(), 0.67, s.matcher.ruleWeight, 0.01)
	assert.InDelta(s.T(), 0.33, s.matcher.similarityWeight, 0.01)
}

// TestAggregateResults 测试结果聚合
func (s *HybridIntentMatcherTestSuite) TestAggregateResults() {
	// Arrange
	results := []*MatchOutput{
		{BotID: "bot-1", Confidence: 0.5, MatchType: "keyword"},
		{BotID: "bot-1", Confidence: 0.3, MatchType: "similarity"}, // 同一个Bot
		{BotID: "bot-2", Confidence: 0.4, MatchType: "keyword"},
	}

	// Act
	aggregated := s.matcher.aggregateResults(results)

	// Assert
	assert.Len(s.T(), aggregated, 2) // 应该聚合成2个结果

	// 验证bot-1的聚合
	var bot1Result *MatchOutput
	for _, r := range aggregated {
		if r.BotID == "bot-1" {
			bot1Result = r
			break
		}
	}
	assert.NotNil(s.T(), bot1Result)
	assert.InDelta(s.T(), 0.8, bot1Result.Confidence, 0.01) // 0.5 + 0.3
	assert.Contains(s.T(), bot1Result.MatchType, "keyword")
	assert.Contains(s.T(), bot1Result.MatchType, "similarity")
}

// TestAggregateResults_WorkflowID 测试WorkflowID聚合
func (s *HybridIntentMatcherTestSuite) TestAggregateResults_WorkflowID() {
	// Arrange
	results := []*MatchOutput{
		{WorkflowID: "workflow-1", Confidence: 0.5, MatchType: "keyword"},
		{WorkflowID: "workflow-1", Confidence: 0.3, MatchType: "similarity"},
	}

	// Act
	aggregated := s.matcher.aggregateResults(results)

	// Assert
	assert.Len(s.T(), aggregated, 1)
	assert.Equal(s.T(), "workflow-1", aggregated[0].WorkflowID)
}

// TestAggregateResults_EmptyBotIDAndWorkflowID 测试空BotID和WorkflowID
func (s *HybridIntentMatcherTestSuite) TestAggregateResults_EmptyBotIDAndWorkflowID() {
	// Arrange
	results := []*MatchOutput{
		{BotID: "", WorkflowID: "", Confidence: 0.5, MatchType: "keyword"}, // 无效结果
		{BotID: "bot-1", Confidence: 0.3, MatchType: "similarity"},
	}

	// Act
	aggregated := s.matcher.aggregateResults(results)

	// Assert - 空ID的结果应该被跳过
	assert.Len(s.T(), aggregated, 1)
	assert.Equal(s.T(), "bot-1", aggregated[0].BotID)
}

// 运行测试套件
func TestHybridIntentMatcherSuite(t *testing.T) {
	suite.Run(t, new(HybridIntentMatcherTestSuite))
}

// TestHybridIntentMatcher_ZeroWeights 测试零权重
func TestHybridIntentMatcher_ZeroWeights(t *testing.T) {
	mockRuleMatcher := new(MockRuleMatcherForHybrid)
	mockSimilarityMatcher := new(MockSimilarityMatcherForHybrid)
	matcher := NewHybridIntentMatcher(mockRuleMatcher, mockSimilarityMatcher)

	// 设置权重为0和0，应该保持默认权重
	matcher.SetWeights(0, 0)

	// 应该保持默认值或合理值
	assert.Greater(t, matcher.ruleWeight, 0.0)
	assert.Greater(t, matcher.similarityWeight, 0.0)
}
