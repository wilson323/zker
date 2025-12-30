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

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
)

// RoutingLogService 路由日志服务
type RoutingLogService struct {
	logRepo repository.RoutingLogRepository
}

// NewRoutingLogService 创建路由日志服务
func NewRoutingLogService(logRepo repository.RoutingLogRepository) *RoutingLogService {
	return &RoutingLogService{
		logRepo: logRepo,
	}
}

// LogRouting 记录路由决策
func (s *RoutingLogService) LogRouting(
	ctx context.Context,
	tenantID string,
	userInput string,
	decision *RoutingDecision,
) error {
	// 构建路由日志实体
	log := &entity.RoutingLog{
		TenantID:         tenantID,
		UserInput:        userInput,
		MatchedRuleID:    &decision.RuleID,
		MatchedBotID:     &decision.BotID,
		MatchedWorkflowID: &decision.WorkflowID,
		Confidence:       &decision.Confidence,
		RoutingScore:    &decision.Score,
		MatchType:        &decision.MatchType,
	}

	return s.logRepo.Create(ctx, log)
}

// GetRoutingLogs 获取路由日志
func (s *RoutingLogService) GetRoutingLogs(
	ctx context.Context,
	tenantID string,
	limit int,
	offset int,
) ([]*entity.RoutingLog, error) {
	return s.logRepo.ListByTenant(ctx, tenantID, limit, offset)
}

// GetRoutingStats 获取路由统计
func (s *RoutingLogService) GetRoutingStats(
	ctx context.Context,
	tenantID string,
) (*RoutingStats, error) {
	logs, err := s.logRepo.ListByTenant(ctx, tenantID, 1000, 0)
	if err != nil {
		return nil, err
	}

	stats := &RoutingStats{
		TotalRoutings:   len(logs),
		BotRoutings:     0,
		WorkflowRoutings: 0,
		AvgConfidence:   0,
		AvgScore:        0,
	}

	totalConfidence := 0.0
	totalScore := 0.0

	for _, log := range logs {
		if log.MatchedBotID != nil && *log.MatchedBotID != "" {
			stats.BotRoutings++
		}
		if log.MatchedWorkflowID != nil && *log.MatchedWorkflowID != "" {
			stats.WorkflowRoutings++
		}
		if log.Confidence != nil {
			totalConfidence += *log.Confidence
		}
		if log.RoutingScore != nil {
			totalScore += *log.RoutingScore
		}
	}

	if len(logs) > 0 {
		stats.AvgConfidence = totalConfidence / float64(len(logs))
		stats.AvgScore = totalScore / float64(len(logs))
	}

	return stats, nil
}

// RoutingStats 路由统计信息
type RoutingStats struct {
	TotalRoutings    int
	BotRoutings      int
	WorkflowRoutings int
	AvgConfidence    float64
	AvgScore         float64
}
