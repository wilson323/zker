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
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
	"go.uber.org/zap"
)

// FeatureExtractor 特征提取器
// 职责：从用户查询和上下文中提取路由决策所需的多维特征
type FeatureExtractor struct {
	intentRepo   repository.IntentRepository
	routingLogRepo repository.RoutingOptimizeLogRepository
	logger       *zap.Logger
}

// NewFeatureExtractor 创建特征提取器
func NewFeatureExtractor(
	intentRepo repository.IntentRepository,
	routingLogRepo repository.RoutingOptimizeLogRepository,
	logger *zap.Logger,
) *FeatureExtractor {
	return &FeatureExtractor{
		intentRepo:     intentRepo,
		routingLogRepo: routingLogRepo,
		logger:         logger,
	}
}

// RoutingFeature 路由特征集合
type RoutingFeature struct {
	QueryFeatures    QueryFeatures    `json:"query_features"`
	BotFeatures      BotFeatures      `json:"bot_features"`
	TemporalFeatures TemporalFeatures `json:"temporal_features"`
	ContextFeatures  ContextFeatures  `json:"context_features"`
	UserFeatures     UserFeatures     `json:"user_features"`
}

// QueryFeatures 查询文本特征
type QueryFeatures struct {
	Length          int     `json:"length"`
	WordCount       int     `json:"word_count"`
	CharCount       int     `json:"char_count"`
	SentenceCount   int     `json:"sentence_count"`
	AvgWordLength   float64 `json:"avg_word_length"`
	IntentScore     float64 `json:"intent_score"`
	Category        string  `json:"category"`
	Complexity      float64 `json:"complexity"`       // 0-1, 查询复杂度
	HasNumbers      bool    `json:"has_numbers"`
	HasDates        bool    `json:"has_dates"`
	HasEntities     bool    `json:"has_entities"`
	QuestionType    string  `json:"question_type"`    // what/where/when/how/why/who
	urgencyLevel    int     `json:"urgency_level"`    // 1-5, 紧急程度
	EmotionalTone   string  `json:"emotional_tone"`   // positive/neutral/negative
}

// BotFeatures Bot性能特征
type BotFeatures struct {
	BotID            string  `json:"bot_id"`
	AvgResponseTime  float64 `json:"avg_response_time"`  // ms
	SuccessRate      float64 `json:"success_rate"`       // 0-1
	UserSatisfaction float64 `json:"user_satisfaction"`  // 0-5
	ResourceUsage    float64 `json:"resource_usage"`     // 0-1
	ErrorRate        float64 `json:"error_rate"`         // 0-1
	TotalRequests    int64   `json:"total_requests"`
	LastUsedTime     int64   `json:"last_used_time"`
	Specialization   string  `json:"specialization"`     // Bot专长领域
	Capabilities     []string `json:"capabilities"`      // Bot能力列表
}

// TemporalFeatures 时间特征
type TemporalFeatures struct {
	HourOfDay        int       `json:"hour_of_day"`
	DayOfWeek        int       `json:"day_of_week"`
	IsWeekend        bool      `json:"is_weekend"`
	IsPeakHour       bool      `json:"is_peak_hour"`
	IsBusinessHour   bool      `json:"is_business_hour"`
	TimeZone         string    `json:"time_zone"`
	Season           string    `json:"season"`         // spring/summer/autumn/winter
	Timestamp        int64     `json:"timestamp"`
}

// ContextFeatures 上下文特征
type ContextFeatures struct {
	TenantID         string   `json:"tenant_id"`
	UserSegment      string   `json:"user_segment"`       // 用户群组
	SessionLength    int      `json:"session_length"`     // 会话长度
	PreviousRoutes   []string `json:"previous_routes"`    // 历史路由
	ConversationDepth int     `json:"conversation_depth"` // 对话深度
	HasContext       bool     `json:"has_context"`        // 是否有上下文
	Channel          string   `json:"channel"`            // 接入渠道
	Platform         string   `json:"platform"`           // 平台
}

// UserFeatures 用户特征
type UserFeatures struct {
	UserID           string  `json:"user_id"`
	UserType         string  `json:"user_type"`         // new/active/loyal/churned
	RegistrationDays int     `json:"registration_days"`  // 注册天数
	TotalQueries     int     `json:"total_queries"`
	AvgRating        float64 `json:"avg_rating"`         // 平均评分
	PreferredBots    []string `json:"preferred_bots"`    // 偏好的Bot
	Language         string  `json:"language"`           // 语言
}

// Extract 提取路由特征
func (e *FeatureExtractor) Extract(
	ctx context.Context,
	query string,
	tenantID string,
	userID string,
) (*RoutingFeature, error) {
	// 1. 提取查询特征
	queryFeatures := e.extractQueryFeatures(ctx, query)

	// 2. 提取时间特征
	temporalFeatures := e.extractTemporalFeatures()

	// 3. 提取上下文特征
	contextFeatures := e.extractContextFeatures(ctx, tenantID, userID)

	// 4. 提取用户特征
	userFeatures := e.extractUserFeatures(ctx, userID, tenantID)

	// 5. Bot特征需要动态获取（由外部填充）
	botFeatures := BotFeatures{}

	return &RoutingFeature{
		QueryFeatures:    queryFeatures,
		BotFeatures:      botFeatures,
		TemporalFeatures: temporalFeatures,
		ContextFeatures:  contextFeatures,
		UserFeatures:     userFeatures,
	}, nil
}

// extractQueryFeatures 提取查询文本特征
func (e *FeatureExtractor) extractQueryFeatures(ctx context.Context, query string) QueryFeatures {
	features := QueryFeatures{}

	// 基础统计
	features.Length = len([]rune(query)) // 字符数（考虑Unicode）
	features.CharCount = len(query)

	// 分词统计
	words := strings.Fields(query)
	features.WordCount = len(words)

	// 句子统计
	sentences := regexp.MustCompile(`[.!?]+`).Split(query, -1)
	features.SentenceCount = len(sentences)

	// 平均词长
	if features.WordCount > 0 {
		totalWordLen := 0
		for _, word := range words {
			totalWordLen += len([]rune(word))
		}
		features.AvgWordLength = float64(totalWordLen) / float64(features.WordCount)
	}

	// 检测数字
	features.HasNumbers = regexp.MustCompile(`\d`).MatchString(query)

	// 检测日期
	features.HasDates = regexp.MustCompile(`\d{4}[-/]\d{1,2}[-/]\d{1,2}|`+
		`\d{1,2}[-/]\d{1,2}[-/]\d{4}|`+
		`今天|明天|昨天`).MatchString(query)

	// 检测实体（简化版）
	features.HasEntities = e.detectEntities(query)

	// 问题类型
	features.QuestionType = e.detectQuestionType(query)

	// 复杂度评分
	features.Complexity = e.calculateComplexity(query, features)

	// 紧急程度
	features.urgencyLevel = e.detectUrgency(query)

	// 情感倾向
	features.EmotionalTone = e.detectEmotionalTone(query)

	// 意图分类（简化版）
	features.IntentScore = e.calculateIntentScore(query)
	features.Category = e.categorizeQuery(query)

	return features
}

// extractTemporalFeatures 提取时间特征
func (e *FeatureExtractor) extractTemporalFeatures() TemporalFeatures {
	now := time.Now()

	features := TemporalFeatures{
		HourOfDay:     now.Hour(),
		DayOfWeek:     int(now.Weekday()),
		IsWeekend:     now.Weekday() == time.Saturday || now.Weekday() == time.Sunday,
		IsPeakHour:    now.Hour() >= 9 && now.Hour() <= 18,
		IsBusinessHour: now.Hour() >= 8 && now.Hour() <= 20,
		TimeZone:      "Local",
		Season:        e.getSeason(now),
		Timestamp:     now.Unix(),
	}

	return features
}

// extractContextFeatures 提取上下文特征
func (e *FeatureExtractor) extractContextFeatures(
	ctx context.Context,
	tenantID string,
	userID string,
) ContextFeatures {
	features := ContextFeatures{
		TenantID:        tenantID,
		UserSegment:     e.determineUserSegment(ctx, userID),
		PreviousRoutes:  []string{}, // 需要从会话上下文获取
		ConversationDepth: 0,         // 需要从会话上下文获取
		HasContext:      false,       // 需要从会话上下文获取
		Channel:         "web",       // 默认值
		Platform:        "desktop",   // 默认值
	}

	return features
}

// extractUserFeatures 提取用户特征
func (e *FeatureExtractor) extractUserFeatures(
	ctx context.Context,
	userID string,
	tenantID string,
) UserFeatures {
	features := UserFeatures{
		UserID:           userID,
		UserType:         "new",       // 默认新用户
		RegistrationDays: 0,
		TotalQueries:     0,
		AvgRating:        0.0,
		PreferredBots:    []string{},
		Language:         "zh-CN",     // 默认中文
	}

	// TODO: 从数据库获取用户历史数据
	// 这里可以查询用户表、历史日志等

	return features
}

// detectEntities 检测实体（简化版）
func (e *FeatureExtractor) detectEntities(query string) bool {
	// 检测常见实体类型
	entityPatterns := []string{
		`\b[A-Z][a-z]+\s[A-Z][a-z]+\b`,     // 人名模式
		`\b\d{11}\b`,                        // 手机号
		`\b[\w._%+-]+@[\w.-]+\.[A-Z]{2,}\b`, // 邮箱
		`\b[A-Z]{2,}\b`,                     // 缩写/机构名
	}

	for _, pattern := range entityPatterns {
		if regexp.MustCompile(pattern).MatchString(query) {
			return true
		}
	}

	return false
}

// detectQuestionType 检测问题类型
func (e *FeatureExtractor) detectQuestionType(query string) string {
	lowerQuery := strings.ToLower(strings.TrimSpace(query))

	questionWords := map[string]string{
		"what":   "what",
		"where":  "where",
		"when":   "when",
		"how":    "how",
		"why":    "why",
		"who":    "who",
		"which":  "which",
		"什么":   "what",
		"哪里":   "where",
		"何时":   "when",
		"怎么":   "how",
		"为什么":  "why",
		"谁":     "who",
		"哪个":   "which",
	}

	for word, qType := range questionWords {
		if strings.HasPrefix(lowerQuery, word) || strings.Contains(lowerQuery, word+" ") {
			return qType
		}
	}

	// 检查是否有问号
	if strings.Contains(query, "?") || strings.Contains(query, "？") {
		return "general"
	}

	return "statement"
}

// calculateComplexity 计算查询复杂度
func (e *FeatureExtractor) calculateComplexity(query string, features QueryFeatures) float64 {
	complexity := 0.0

	// 1. 长度因子（0-0.3）
	if features.Length > 0 {
		lengthScore := math.Min(float64(features.Length)/200.0, 1.0) * 0.3
		complexity += lengthScore
	}

	// 2. 词数因子（0-0.2）
	if features.WordCount > 0 {
		wordScore := math.Min(float64(features.WordCount)/30.0, 1.0) * 0.2
		complexity += wordScore
	}

	// 3. 句子数因子（0-0.2）
	if features.SentenceCount > 1 {
		sentenceScore := math.Min(float64(features.SentenceCount)/5.0, 1.0) * 0.2
		complexity += sentenceScore
	}

	// 4. 特殊字符因子（0-0.1）
	specialChars := 0
	for _, r := range query {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r) {
			specialChars++
		}
	}
	if specialChars > 0 {
		specialScore := math.Min(float64(specialChars)/10.0, 1.0) * 0.1
		complexity += specialScore
	}

	// 5. 实体和数字因子（0-0.2）
	entityScore := 0.0
	if features.HasNumbers {
		entityScore += 0.1
	}
	if features.HasDates {
		entityScore += 0.05
	}
	if features.HasEntities {
		entityScore += 0.05
	}
	complexity += entityScore

	return math.Min(complexity, 1.0)
}

// detectUrgency 检测紧急程度（1-5）
func (e *FeatureExtractor) detectUrgency(query string) int {
	lowerQuery := strings.ToLower(query)

	urgentKeywords := map[string]int{
		"紧急":  5,
		"马上":  5,
		"立即":  5,
		"urgent": 5,
		"asap":   5,
		"尽快":   4,
		"麻烦":   3,
		"急":     3,
		"help":   3,
	}

	maxUrgency := 1 // 默认最低紧急度
	for keyword, level := range urgentKeywords {
		if strings.Contains(lowerQuery, keyword) && level > maxUrgency {
			maxUrgency = level
		}
	}

	return maxUrgency
}

// detectEmotionalTone 检测情感倾向
func (e *FeatureExtractor) detectEmotionalTone(query string) string {
	lowerQuery := strings.ToLower(query)

	// 积极词汇
	positiveWords := []string{
		"好", "棒", "优秀", "满意", "喜欢", "谢谢", "感谢",
		"good", "great", "excellent", "thanks", "thank you",
	}

	// 消极词汇
	negativeWords := []string{
		"差", "坏", "糟糕", "不满", "讨厌", "投诉", "问题", "错误",
		"bad", "terrible", "awful", "hate", "complaint", "issue", "error",
	}

	positiveCount := 0
	negativeCount := 0

	for _, word := range positiveWords {
		if strings.Contains(lowerQuery, word) {
			positiveCount++
		}
	}

	for _, word := range negativeWords {
		if strings.Contains(lowerQuery, word) {
			negativeCount++
		}
	}

	if positiveCount > negativeCount {
		return "positive"
	} else if negativeCount > positiveCount {
		return "negative"
	}

	return "neutral"
}

// calculateIntentScore 计算意图得分
func (e *FeatureExtractor) calculateIntentScore(query string) float64 {
	// 简化版意图得分计算
	// 实际应用中应使用机器学习模型

	score := 0.5 // 默认中等置信度

	// 如果包含明确的意图词，提高得分
	intentKeywords := []string{
		"查询", "搜索", "找", "问",
		"query", "search", "find", "ask",
	}

	for _, keyword := range intentKeywords {
		if strings.Contains(query, keyword) {
			score += 0.1
		}
	}

	return math.Min(score, 1.0)
}

// categorizeQuery 对查询分类
func (e *FeatureExtractor) categorizeQuery(query string) string {
	lowerQuery := strings.ToLower(query)

	categories := map[string][]string{
		"technical": {"代码", "编程", "bug", "开发", "api", "code", "programming", "bug", "development"},
		"business":  {"订单", "支付", "价格", "产品", "服务", "order", "payment", "price", "product", "service"},
		"support":   {"问题", "帮助", "客服", "使用", "problem", "help", "support", "how to"},
		"general":   {"什么", "怎么", "如何", "what", "how", "why"},
	}

	maxMatches := 0
	bestCategory := "general"

	for category, keywords := range categories {
		matches := 0
		for _, keyword := range keywords {
			if strings.Contains(lowerQuery, keyword) {
				matches++
			}
		}

		if matches > maxMatches {
			maxMatches = matches
			bestCategory = category
		}
	}

	return bestCategory
}

// getSeason 获取季节
func (e *FeatureExtractor) getSeason(t time.Time) string {
	month := t.Month()
	switch {
	case month >= 3 && month <= 5:
		return "spring"
	case month >= 6 && month <= 8:
		return "summer"
	case month >= 9 && month <= 11:
		return "autumn"
	default:
		return "winter"
	}
}

// determineUserSegment 确定用户群组
func (e *FeatureExtractor) determineUserSegment(ctx context.Context, userID string) string {
	// TODO: 实际应该从数据库获取用户分组信息
	// 这里返回默认值
	return "default"
}
