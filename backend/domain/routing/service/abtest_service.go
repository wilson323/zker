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
	"time"

	"github.com/google/uuid"
	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"go.uber.org/zap"
)

// ABTestService A/B测试服务
// 职责：管理路由策略的A/B测试，包括测试创建、执行、统计分析
type ABTestService struct {
	abTestRepo     repository.ABTestRepository
	abRecordRepo   repository.ABTestRecordRepository
	logger         *zap.Logger
}

// NewABTestService 创建A/B测试服务实例
func NewABTestService(
	abTestRepo repository.ABTestRepository,
	abRecordRepo repository.ABTestRecordRepository,
	logger *zap.Logger,
) *ABTestService {
	return &ABTestService{
		abTestRepo:   abTestRepo,
		abRecordRepo: abRecordRepo,
		logger:       logger,
	}
}

// CreateABTest 创建A/B测试
func (s *ABTestService) CreateABTest(ctx context.Context, test *entity.ABTest) error {
	// 1. 验证测试参数
	if err := s.validateABTest(test); err != nil {
		return errorx.WrapByCode(err, errno.ErrRouteFormatInvalidCode)
	}

	// 2. 设置默认值
	if test.Status == "" {
		test.Status = entity.ABTestStatusRunning
	}
	if test.TrafficSplit == 0 {
		test.TrafficSplit = 50 // 默认50:50分配
	}
	if test.SampleSize == 0 {
		test.SampleSize = 1000 // 默认1000样本
	}

	// 3. 设置时间
	now := time.Now()
	test.CreatedAt = now
	test.UpdatedAt = now
	if test.Status == entity.ABTestStatusRunning && test.StartTime == nil {
		test.StartTime = &now
	}

	// 4. 创建测试
	if err := s.abTestRepo.Create(ctx, test); err != nil {
		s.logger.Error("failed to create ab test",
			zap.String("tenant_id", test.TenantID),
			zap.Error(err))
		return errorx.WrapByCode(err, errno.ErrRoutingDecisionFailedCode)
	}

	return nil
}

// ExecuteABTest 执行A/B测试路由
func (s *ABTestService) ExecuteABTest(ctx context.Context, testID string, req *RoutingRequest) (*entity.RoutingDecision, error) {
	// 1. 获取测试信息
	test, err := s.abTestRepo.GetByID(ctx, testID)
	if err != nil {
		return nil, errorx.New(errno.ErrABTestNotFoundCode, errorx.KV("test_id", testID))
	}

	// 2. 验证测试状态
	if !test.IsTestActive() {
		return nil, errorx.New(errno.ErrRouteFormatInvalidCode, errorx.KV("reason", "test is not active"))
	}

	// 3. 根据流量分配选择策略
	trafficValue := int(hash(req.UserID+req.Query) % 100)
	strategy := test.GetStrategyForTraffic(trafficValue)

	// 4. 记录测试执行
	record := &entity.ABTestRecord{
		RecordID:  uuid.New().String(),
		TestID:    testID,
		Strategy:  strategy,
		UserID:    req.UserID,
		AgentID:   "", // 将由外部填充
		CreatedAt: time.Now(),
	}

	// 5. 异步保存记录
	go func() {
		_ = s.abRecordRepo.Create(context.Background(), record)
	}()

	// 6. 返回路由决策（包含策略信息）
	decision := &entity.RoutingDecision{
		Strategy: strategy,
		Reasons:  []string{fmt.Sprintf("A/B Test: %s", strategy)},
	}

	return decision, nil
}

// GetABTestResults 获取A/B测试结果
func (s *ABTestService) GetABTestResults(ctx context.Context, testID string) (*entity.ABTestResult, error) {
	// 1. 获取测试信息
	test, err := s.abTestRepo.GetByID(ctx, testID)
	if err != nil {
		return nil, errorx.New(errno.ErrABTestNotFoundCode, errorx.KV("test_id", testID))
	}

	// 2. 获取测试记录
	records, err := s.abRecordRepo.GetByTestID(ctx, testID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrRoutingDecisionFailedCode, errorx.KV("error", err.Error()))
	}

	// 3. 统计结果
	result := &entity.ABTestResult{
		TestID:     testID,
		StrategyAStats: s.calculateStats(records, test.StrategyA),
		StrategyBStats: s.calculateStats(records, test.StrategyB),
		Winner:     test.Winner,
		Confidence: test.Confidence,
	}

	// 4. 判断统计显著性
	result.IsStatisticallySignificant = s.isStatisticallySignificant(result)

	// 5. 生成推荐
	result.Recommendation = s.generateRecommendation(result)

	return result, nil
}

// ConcludeABTest 结束A/B测试并选择获胜者
func (s *ABTestService) ConcludeABTest(ctx context.Context, testID string, winner string) error {
	// 1. 获取测试信息
	test, err := s.abTestRepo.GetByID(ctx, testID)
	if err != nil {
		return errorx.New(errno.ErrABTestNotFoundCode, errorx.KV("reason", "test not found"), errorx.KV("test_id", testID))
	}

	// 2. 验证测试状态
	if test.Status != entity.ABTestStatusRunning {
		return errorx.New(errno.ErrRouteFormatInvalidCode, errorx.KV("reason", "test is not running"))
	}

	// 3. 获取测试结果
	result, err := s.GetABTestResults(ctx, testID)
	if err != nil {
		return err
	}

	// 4. 更新测试状态
	now := time.Now()
	test.Status = entity.ABTestStatusCompleted
	test.Winner = winner
	test.Confidence = result.Confidence
	test.EndTime = &now
	test.UpdatedAt = now

	if err := s.abTestRepo.Update(ctx, test); err != nil {
		s.logger.Error("failed to update ab test",
			zap.String("test_id", testID),
			zap.Error(err))
		return errorx.WrapByCode(err, errno.ErrRoutingDecisionFailedCode)
	}

	return nil
}

// validateABTest 验证A/B测试参数
func (s *ABTestService) validateABTest(test *entity.ABTest) error {
	if test.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if test.TestName == "" {
		return fmt.Errorf("test_name is required")
	}
	if test.StrategyA == "" || test.StrategyB == "" {
		return fmt.Errorf("strategy_a and strategy_b are required")
	}
	if test.StrategyA == test.StrategyB {
		return fmt.Errorf("strategy_a and strategy_b must be different")
	}
	if test.TrafficSplit < 0 || test.TrafficSplit > 100 {
		return fmt.Errorf("traffic_split must be between 0 and 100")
	}
	return nil
}

// calculateStats 计算单个策略的统计数据
func (s *ABTestService) calculateStats(records []*entity.ABTestRecord, strategy string) *entity.TestStrategyStats {
	stats := &entity.TestStrategyStats{
		Strategy:    strategy,
		SampleCount: 0,
	}

	// 过滤该策略的记录
	strategyRecords := make([]*entity.ABTestRecord, 0)
	for _, record := range records {
		if record.Strategy == strategy {
			strategyRecords = append(strategyRecords, record)
		}
	}

	if len(strategyRecords) == 0 {
		return stats
	}

	stats.SampleCount = len(strategyRecords)

	// 计算平均响应时间
	totalTime := 0
	for _, record := range strategyRecords {
		totalTime += record.ResponseTime
	}
	stats.AvgResponseTime = float64(totalTime) / float64(len(strategyRecords))

	// 计算平均用户评分
	totalRating := 0
	ratingCount := 0
	for _, record := range strategyRecords {
		if record.UserRating > 0 {
			totalRating += record.UserRating
			ratingCount++
		}
	}
	if ratingCount > 0 {
		stats.AvgUserRating = float64(totalRating) / float64(ratingCount)
	}

	// 计算成功率（用户评分>=3视为成功）
	successCount := 0
	for _, record := range strategyRecords {
		if record.UserRating >= 3 {
			successCount++
		}
	}
	stats.SuccessRate = float64(successCount) / float64(len(strategyRecords))

	// 错误率（用户评分<=2视为错误）
	errorCount := 0
	for _, record := range strategyRecords {
		if record.UserRating > 0 && record.UserRating <= 2 {
			errorCount++
		}
	}
	stats.ErrorRate = float64(errorCount) / float64(len(strategyRecords))

	return stats
}

// isStatisticallySignificant 判断统计显著性
// 使用Z检验计算两个策略的显著性差异
func (s *ABTestService) isStatisticallySignificant(result *entity.ABTestResult) bool {
	if result.StrategyAStats == nil || result.StrategyBStats == nil {
		return false
	}

	// 样本量检查
	if result.StrategyAStats.SampleCount < 30 || result.StrategyBStats.SampleCount < 30 {
		return false
	}

	// 计算Z分数
	p1 := result.StrategyAStats.SuccessRate
	p2 := result.StrategyBStats.SuccessRate
	n1 := float64(result.StrategyAStats.SampleCount)
	n2 := float64(result.StrategyBStats.SampleCount)

	// 合并比例
	pooledP := (p1*n1 + p2*n2) / (n1 + n2)

	// 标准误差
	se := math.Sqrt(pooledP * (1 - pooledP) * (1/n1 + 1/n2))

	if se == 0 {
		return false
	}

	// Z分数
	z := (p1 - p2) / se

	// 95%置信度，Z分数应大于1.96
	return math.Abs(z) >= 1.96
}

// generateRecommendation 生成推荐
func (s *ABTestService) generateRecommendation(result *entity.ABTestResult) string {
	if !result.IsStatisticallySignificant {
		return "样本量不足或差异不显著，建议继续测试"
	}

	if result.StrategyAStats == nil || result.StrategyBStats == nil {
		return "数据不足，无法生成推荐"
	}

	// 比较成功率
	if result.StrategyAStats.SuccessRate > result.StrategyBStats.SuccessRate {
		return fmt.Sprintf("推荐策略A: 成功率较高 (%.2f%% vs %.2f%%)",
			result.StrategyAStats.SuccessRate*100,
			result.StrategyBStats.SuccessRate*100)
	} else if result.StrategyBStats.SuccessRate > result.StrategyAStats.SuccessRate {
		return fmt.Sprintf("推荐策略B: 成功率较高 (%.2f%% vs %.2f%%)",
			result.StrategyBStats.SuccessRate*100,
			result.StrategyAStats.SuccessRate*100)
	}

	return "两个策略表现相似，建议结合其他指标综合考虑"
}

// hash 简单哈希函数（用于流量分配）
func hash(s string) uint32 {
	h := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(s); i++ {
		h *= prime32
		h ^= uint32(s[i])
	}
	return h
}

// PauseABTest 暂停A/B测试
func (s *ABTestService) PauseABTest(ctx context.Context, testID string) error {
	test, err := s.abTestRepo.GetByID(ctx, testID)
	if err != nil {
		return errorx.New(errno.ErrABTestNotFoundCode, errorx.KV("reason", "test not found"), errorx.KV("test_id", testID))
	}

	if test.Status != entity.ABTestStatusRunning {
		return errorx.New(errno.ErrRouteFormatInvalidCode, errorx.KV("reason", "test is not running"))
	}

	test.Status = entity.ABTestStatusPaused
	test.UpdatedAt = time.Now()

	if err := s.abTestRepo.Update(ctx, test); err != nil {
		return errorx.WrapByCode(err, errno.ErrRoutingDecisionFailedCode)
	}

	return nil
}

// ResumeABTest 恢复A/B测试
func (s *ABTestService) ResumeABTest(ctx context.Context, testID string) error {
	test, err := s.abTestRepo.GetByID(ctx, testID)
	if err != nil {
		return errorx.New(errno.ErrABTestNotFoundCode, errorx.KV("reason", "test not found"), errorx.KV("test_id", testID))
	}

	if test.Status != entity.ABTestStatusPaused {
		return errorx.New(errno.ErrRouteFormatInvalidCode, errorx.KV("reason", "test is not paused"))
	}

	test.Status = entity.ABTestStatusRunning
	test.UpdatedAt = time.Now()

	if err := s.abTestRepo.Update(ctx, test); err != nil {
		return errorx.WrapByCode(err, errno.ErrRoutingDecisionFailedCode)
	}

	return nil
}
