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
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
)

// MockRoutingRuleRepository 路由规则仓储 Mock
type MockRoutingRuleRepository struct {
	mock.Mock
}

func (m *MockRoutingRuleRepository) Create(ctx context.Context, rule *entity.RoutingRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockRoutingRuleRepository) GetByID(ctx context.Context, ruleID string) (*entity.RoutingRule, error) {
	args := m.Called(ctx, ruleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RoutingRule), args.Error(1)
}

func (m *MockRoutingRuleRepository) Update(ctx context.Context, rule *entity.RoutingRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockRoutingRuleRepository) Delete(ctx context.Context, ruleID string) error {
	args := m.Called(ctx, ruleID)
	return args.Error(0)
}

func (m *MockRoutingRuleRepository) GetActiveRulesByTenant(ctx context.Context, tenantID string) ([]*entity.RoutingRule, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.RoutingRule), args.Error(1)
}

func (m *MockRoutingRuleRepository) List(ctx context.Context, filter *repository.RoutingRuleFilter) ([]*entity.RoutingRule, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.RoutingRule), args.Get(1).(int64), args.Error(2)
}

// RuleBasedMatcherTestSuite 测试套件
type RuleBasedMatcherTestSuite struct {
	suite.Suite
	matcher     *RuleBasedMatcher
	mockRuleRepo *MockRoutingRuleRepository
	ctx         context.Context
}

func (s *RuleBasedMatcherTestSuite) SetupTest() {
	s.mockRuleRepo = new(MockRoutingRuleRepository)
	s.matcher = NewRuleBasedMatcher(s.mockRuleRepo)
	s.ctx = context.Background()
}

// TestMatch_KeywordMatch 测试关键词匹配
func (s *RuleBasedMatcherTestSuite) TestMatch_KeywordMatch() {
	// Arrange
	tenantID := "tenant-123"
	keywordRule := &entity.RoutingRule{
		RuleID:    "rule-keyword",
		TenantID:  tenantID,
		RuleName:  "Keyword Rule",
		RuleType:  entity.RuleTypeKeyword,
		Priority:  10,
		IsActive:  true,
		TargetBotID: func() *string { s := "bot-keyword"; return &s }(),
	}

	// 构建条件
	condition := map[string]interface{}{
		"keywords": []string{"天气", "温度", "气温"},
	}
	conditionJSON, _ := json.Marshal(condition)
	keywordRule.Condition = string(conditionJSON)

	s.mockRuleRepo.On("GetActiveRulesByTenant", s.ctx, tenantID).Return([]*entity.RoutingRule{keywordRule}, nil)

	input := &MatchInput{
		UserInput: "今天北京天气怎么样？温度多少？",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 1)
	assert.Equal(s.T(), "bot-keyword", results[0].BotID)
	assert.Greater(s.T(), results[0].Confidence, 0.5) // 应该有较高的置信度（匹配了2个关键词）
	assert.Equal(s.T(), string(entity.RuleTypeKeyword), results[0].MatchType)
	s.mockRuleRepo.AssertExpectations(s.T())
}

// TestMatch_RegexMatch 测试正则表达式匹配
func (s *RuleBasedMatcherTestSuite) TestMatch_RegexMatch() {
	// Arrange
	tenantID := "tenant-123"
	regexRule := &entity.RoutingRule{
		RuleID:    "rule-regex",
		TenantID:  tenantID,
		RuleName:  "Regex Rule",
		RuleType:  entity.RuleTypeRegex,
		Priority:  10,
		IsActive:  true,
		TargetBotID: func() *string { s := "bot-regex"; return &s }(),
	}

	condition := map[string]interface{}{
		"pattern": `^(\d{11})$`, // 匹配11位手机号
	}
	conditionJSON, _ := json.Marshal(condition)
	regexRule.Condition = string(conditionJSON)

	s.mockRuleRepo.On("GetActiveRulesByTenant", s.ctx, tenantID).Return([]*entity.RoutingRule{regexRule}, nil)

	input := &MatchInput{
		UserInput: "13812345678",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 1)
	assert.Equal(s.T(), "bot-regex", results[0].BotID)
	assert.Equal(s.T(), 1.0, results[0].Confidence) // 正则匹配的置信度是1.0
	s.mockRuleRepo.AssertExpectations(s.T())
}

// TestMatch_RegexMatch_Failed 测试正则表达式匹配失败
func (s *RuleBasedMatcherTestSuite) TestMatch_RegexMatch_Failed() {
	// Arrange
	tenantID := "tenant-123"
	regexRule := &entity.RoutingRule{
		RuleID:    "rule-regex",
		TenantID:  tenantID,
		RuleName:  "Regex Rule",
		RuleType:  entity.RuleTypeRegex,
		Priority:  10,
		IsActive:  true,
		TargetBotID: func() *string { s := "bot-regex"; return &s }(),
	}

	condition := map[string]interface{}{
		"pattern": `^(\d{11})$`,
	}
	conditionJSON, _ := json.Marshal(condition)
	regexRule.Condition = string(conditionJSON)

	s.mockRuleRepo.On("GetActiveRulesByTenant", s.ctx, tenantID).Return([]*entity.RoutingRule{regexRule}, nil)

	input := &MatchInput{
		UserInput: "这不是手机号",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 0) // 没有匹配
	s.mockRuleRepo.AssertExpectations(s.T())
}

// TestMatch_IntentMatch 测试意图匹配
func (s *RuleBasedMatcherTestSuite) TestMatch_IntentMatch() {
	// Arrange
	tenantID := "tenant-123"
	intentRule := &entity.RoutingRule{
		RuleID:    "rule-intent",
		TenantID:  tenantID,
		RuleName:  "Intent Rule",
		RuleType:  entity.RuleTypeIntent,
		Priority:  10,
		IsActive:  true,
		TargetBotID: func() *string { s := "bot-intent"; return &s }(),
	}

	condition := map[string]interface{}{
		"intent": "book_flight",
	}
	conditionJSON, _ := json.Marshal(condition)
	intentRule.Condition = string(conditionJSON)

	s.mockRuleRepo.On("GetActiveRulesByTenant", s.ctx, tenantID).Return([]*entity.RoutingRule{intentRule}, nil)

	input := &MatchInput{
		UserInput: "我想订机票",
		TenantID:  tenantID,
		Context: map[string]interface{}{
			"intent": "book_flight",
		},
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 1)
	assert.Equal(s.T(), "bot-intent", results[0].BotID)
	assert.Equal(s.T(), 1.0, results[0].Confidence)
	s.mockRuleRepo.AssertExpectations(s.T())
}

// TestMatch_CategoryMatch 测试分类匹配
func (s *RuleBasedMatcherTestSuite) TestMatch_CategoryMatch() {
	// Arrange
	tenantID := "tenant-123"
	categoryRule := &entity.RoutingRule{
		RuleID:    "rule-category",
		TenantID:  tenantID,
		RuleName:  "Category Rule",
		RuleType:  entity.RuleTypeCategory,
		Priority:  10,
		IsActive:  true,
		TargetBotID: func() *string { s := "bot-category"; return &s }(),
	}

	condition := map[string]interface{}{
		"category": "customer_service",
	}
	conditionJSON, _ := json.Marshal(condition)
	categoryRule.Condition = string(conditionJSON)

	s.mockRuleRepo.On("GetActiveRulesByTenant", s.ctx, tenantID).Return([]*entity.RoutingRule{categoryRule}, nil)

	input := &MatchInput{
		UserInput: "我要投诉",
		TenantID:  tenantID,
		Context: map[string]interface{}{
			"category": "customer_service",
		},
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 1)
	assert.Equal(s.T(), "bot-category", results[0].BotID)
	assert.Equal(s.T(), 1.0, results[0].Confidence)
	s.mockRuleRepo.AssertExpectations(s.T())
}

// TestMatch_MultipleRules 测试多个规则匹配
func (s *RuleBasedMatcherTestSuite) TestMatch_MultipleRules() {
	// Arrange
	tenantID := "tenant-123"
	bot1 := "bot-1"
	bot2 := "bot-2"

	rules := []*entity.RoutingRule{
		{
			RuleID:    "rule-1",
			TenantID:  tenantID,
			RuleName:  "Rule 1",
			RuleType:  entity.RuleTypeKeyword,
			Priority:  10,
			IsActive:  true,
			TargetBotID: &bot1,
		},
		{
			RuleID:    "rule-2",
			TenantID:  tenantID,
			RuleName:  "Rule 2",
			RuleType:  entity.RuleTypeKeyword,
			Priority:  5,
			IsActive:  true,
			TargetBotID: &bot2,
		},
	}

	// 规则1的条件
	condition1 := map[string]interface{}{
		"keywords": []string{"关键词1"},
	}
	conditionJSON1, _ := json.Marshal(condition1)
	rules[0].Condition = string(conditionJSON1)

	// 规则2的条件
	condition2 := map[string]interface{}{
		"keywords": []string{"关键词2", "关键词3"},
	}
	conditionJSON2, _ := json.Marshal(condition2)
	rules[1].Condition = string(conditionJSON2)

	s.mockRuleRepo.On("GetActiveRulesByTenant", s.ctx, tenantID).Return(rules, nil)

	input := &MatchInput{
		UserInput: "包含关键词1和关键词2的内容",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 2) // 两个规则都匹配
	// 应该按置信度排序
	s.mockRuleRepo.AssertExpectations(s.T())
}

// TestMatch_NoMatch 测试没有匹配
func (s *RuleBasedMatcherTestSuite) TestMatch_NoMatch() {
	// Arrange
	tenantID := "tenant-123"
	keywordRule := &entity.RoutingRule{
		RuleID:    "rule-keyword",
		TenantID:  tenantID,
		RuleName:  "Keyword Rule",
		RuleType:  entity.RuleTypeKeyword,
		Priority:  10,
		IsActive:  true,
		TargetBotID: func() *string { s := "bot-keyword"; return &s }(),
	}

	condition := map[string]interface{}{
		"keywords": []string{"不存在的关键词"},
	}
	conditionJSON, _ := json.Marshal(condition)
	keywordRule.Condition = string(conditionJSON)

	s.mockRuleRepo.On("GetActiveRulesByTenant", s.ctx, tenantID).Return([]*entity.RoutingRule{keywordRule}, nil)

	input := &MatchInput{
		UserInput: "完全不相关的内容",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 0)
	s.mockRuleRepo.AssertExpectations(s.T())
}

// TestMatch_RepositoryError 测试仓储错误
func (s *RuleBasedMatcherTestSuite) TestMatch_RepositoryError() {
	// Arrange
	tenantID := "tenant-123"
	s.mockRuleRepo.On("GetActiveRulesByTenant", s.ctx, tenantID).Return(
		([]*entity.RoutingRule)(nil), errors.New("database error"))

	input := &MatchInput{
		UserInput: "测试内容",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), results)
	s.mockRuleRepo.AssertExpectations(s.T())
}

// TestMatch_Top3Limit 测试Top-3限制
func (s *RuleBasedMatcherTestSuite) TestMatch_Top3Limit() {
	// Arrange
	tenantID := "tenant-123"
	bot1 := "bot-1"
	bot2 := "bot-2"
	bot3 := "bot-3"
	bot4 := "bot-4"
	bot5 := "bot-5"

	rules := make([]*entity.RoutingRule, 5)
	for i := 0; i < 5; i++ {
		botID := []string{bot1, bot2, bot3, bot4, bot5}[i]
		rules[i] = &entity.RoutingRule{
			RuleID:    "rule-" + string(rune('1'+i)),
			TenantID:  tenantID,
			RuleName:  "Rule",
			RuleType:  entity.RuleTypeKeyword,
			Priority:  10,
			IsActive:  true,
			TargetBotID: &botID,
		}

		condition := map[string]interface{}{
			"keywords": []string{"关键词"},
		}
		conditionJSON, _ := json.Marshal(condition)
		rules[i].Condition = string(conditionJSON)
	}

	s.mockRuleRepo.On("GetActiveRulesByTenant", s.ctx, tenantID).Return(rules, nil)

	input := &MatchInput{
		UserInput: "包含关键词的内容",
		TenantID:  tenantID,
	}

	// Act
	results, err := s.matcher.Match(s.ctx, input)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 3) // 应该只返回Top-3
	s.mockRuleRepo.AssertExpectations(s.T())
}

// TestMatchKeyword_CaseInsensitive 测试关键词大小写不敏感
func (s *RuleBasedMatcherTestSuite) TestMatchKeyword_CaseInsensitive() {
	// Arrange
	matcher := NewRuleBasedMatcher(nil)
	condition := map[string]interface{}{
		"keywords": []string{"Weather", "Temperature"},
	}
	input := "WEATHER and temperature"

	// Act
	matched, confidence := matcher.matchKeyword(condition, input)

	// Assert
	assert.True(s.T(), matched)
	assert.Greater(s.T(), confidence, 0.0)
}

// TestMatchRegex_InvalidPattern 测试无效正则表达式
func (s *RuleBasedMatcherTestSuite) TestMatchRegex_InvalidPattern() {
	// Arrange
	matcher := NewRuleBasedMatcher(nil)
	condition := map[string]interface{}{
		"pattern": "[invalid(", // 无效的正则表达式
	}
	input := "test"

	// Act
	matched, confidence := matcher.matchRegex(condition, input)

	// Assert
	assert.False(s.T(), matched)
	assert.Equal(s.T(), 0.0, confidence)
}

// 运行测试套件
func TestRuleBasedMatcherSuite(t *testing.T) {
	suite.Run(t, new(RuleBasedMatcherTestSuite))
}

// TestRuleBasedMatcher_MatchIntent_NoContext 测试无上下文的意图匹配
func TestRuleBasedMatcher_MatchIntent_NoContext(t *testing.T) {
	matcher := NewRuleBasedMatcher(nil)
	condition := map[string]interface{}{
		"intent": "test_intent",
	}

	// 无上下文
	matched, confidence := matcher.matchIntent(condition, nil)

	assert.False(t, matched)
	assert.Equal(t, 0.0, confidence)
}

// TestRuleBasedMatcher_MatchCategory_NoContext 测试无上下文的分类匹配
func TestRuleBasedMatcher_MatchCategory_NoContext(t *testing.T) {
	matcher := NewRuleBasedMatcher(nil)
	condition := map[string]interface{}{
		"category": "test_category",
	}

	// 无上下文
	matched, confidence := matcher.matchCategory(condition, nil)

	assert.False(t, matched)
	assert.Equal(t, 0.0, confidence)
}
