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
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/human_ai_collaboration/entity"
	"github.com/coze-dev/coze-studio/backend/domain/human_ai_collaboration/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// CollaborationTriggerService 人机协同触发服务
// 职责：实现多种协同触发机制(置信度阈值、用户主动、异常情况)
type CollaborationTriggerService struct {
	collabRepo   repository.CollaborationRepository
	logger       *zap.Logger

	// 触发配置
	confidenceThreshold     float64       // 默认置信度阈值
	autoTriggerEnabled      bool          // 是否启用自动触发
	userTriggerEnabled      bool          // 是否启用用户触发
	exceptionTriggerEnabled bool          // 是否启用异常触发
	timeout                 time.Duration // 超时时间
}

// NewCollaborationTriggerService 创建协同触发服务实例
func NewCollaborationTriggerService(
	collabRepo repository.CollaborationRepository,
	logger *zap.Logger,
) *CollaborationTriggerService {
	return &CollaborationTriggerService{
		collabRepo:              collabRepo,
		logger:                  logger,
		confidenceThreshold:     0.7,
		autoTriggerEnabled:      true,
		userTriggerEnabled:      true,
		exceptionTriggerEnabled: true,
		timeout:                 30 * time.Minute,
	}
}

// SetConfig 设置触发配置
func (s *CollaborationTriggerService) SetConfig(
	confidenceThreshold float64,
	autoTrigger, userTrigger, exceptionTrigger bool,
	timeout time.Duration,
) {
	s.confidenceThreshold = confidenceThreshold
	s.autoTriggerEnabled = autoTrigger
	s.userTriggerEnabled = userTrigger
	s.exceptionTriggerEnabled = exceptionTrigger
	s.timeout = timeout
}

// CheckAndTrigger 检查并触发协同
func (s *CollaborationTriggerService) CheckAndTrigger(
	ctx context.Context,
	req *CollaborationRequest,
) (*entity.CollaborationSession, error) {
	// 1. 检查是否应该触发协同
	shouldTrigger, triggerType, reason := s.shouldTriggerCollaboration(ctx, req)
	if !shouldTrigger {
		return nil, nil
	}

	// 2. 创建协同会话
	session := &entity.CollaborationSession{
		SessionID:      uuid.New().String(),
		TenantID:       req.TenantID,
		UserID:         req.UserID,
		ConversationID: req.ConversationID,
		MessageID:      req.MessageID,
		AgentID:        req.AgentID,
		TriggerType:    triggerType,
		TriggerReason:  reason,
		Confidence:     req.Confidence,
		Status:         "pending",
		Priority:       s.calculatePriority(req),
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(s.timeout),
	}

	// 3. 保存协同会话
	if err := s.collabRepo.CreateSession(ctx, session); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrCollaborationCreateFailedCode, errorx.KV("error", err.Error()))
	}

	s.logger.Info("collaboration triggered",
		zap.String("session_id", session.SessionID),
		zap.String("trigger_type", triggerType),
		zap.String("reason", reason))

	return session, nil
}

// shouldTriggerCollaboration 判断是否应该触发协同
func (s *CollaborationTriggerService) shouldTriggerCollaboration(
	ctx context.Context,
	req *CollaborationRequest,
) (bool, string, string) {
	// 1. 检查用户主动触发
	if s.userTriggerEnabled && req.UserInitiated {
		return true, "user_initiated", "用户主动请求人工介入"
	}

	// 2. 检查置信度阈值触发
	if s.autoTriggerEnabled && req.Confidence < s.confidenceThreshold {
		reason := fmt.Sprintf("AI置信度低于阈值(%.2f < %.2f)", req.Confidence, s.confidenceThreshold)
		return true, "low_confidence", reason
	}

	// 3. 检查异常情况触发
	if s.exceptionTriggerEnabled {
		// 检测是否有异常
		if s.hasException(req) {
			return true, "exception_detected", "检测到异常情况，需要人工审核"
		}

		// 检测是否有敏感内容
		if s.hasSensitiveContent(req) {
			return true, "sensitive_content", "检测到敏感内容，需要人工审核"
		}

		// 检测是否有安全风险
		if s.hasSecurityRisk(req) {
			return true, "security_risk", "检测到潜在安全风险，需要人工审核"
		}
	}

	return false, "", ""
}

// hasException 检测是否有异常
func (s *CollaborationTriggerService) hasException(req *CollaborationRequest) bool {
	// 检测LLM返回错误
	if req.HasError {
		return true
	}

	// 检测响应超时
	if req.ResponseTime > 30000 { // 超过30秒
		return true
	}

	// 检测响应异常短
	if len(req.ResponseText) > 0 && len(req.ResponseText) < 10 {
		return true
	}

	return false
}

// hasSensitiveContent 检测是否有敏感内容
func (s *CollaborationTriggerService) hasSensitiveContent(req *CollaborationRequest) bool {
	// 简化实现：实际应用中应该使用敏感词库或内容审核API
	sensitiveKeywords := []string{
		"密码", "银行卡", "身份证", "电话", "地址",
	}

	inputText := req.InputText
	for _, keyword := range sensitiveKeywords {
		if contains(inputText, keyword) {
			return true
		}
	}

	return false
}

// hasSecurityRisk 检测是否有安全风险
func (s *CollaborationTriggerService) hasSecurityRisk(req *CollaborationRequest) bool {
	// 简化实现：检测Prompt注入攻击
	patterns := []string{
		"忽略以上指令", "ignore previous", "jailbreak",
		"<script>", "javascript:", "eval(",
	}

	for _, pattern := range patterns {
		if contains(req.InputText, pattern) {
			return true
		}
	}

	return false
}

// calculatePriority 计算优先级
func (s *CollaborationTriggerService) calculatePriority(req *CollaborationRequest) string {
	// 紧急程度：critical > high > medium > low
	if req.Confidence < 0.3 {
		return "critical"
	} else if req.Confidence < 0.5 {
		return "high"
	} else if req.Confidence < 0.7 {
		return "medium"
	}
	return "low"
}

// GetUserPendingSessions 获取用户待处理协同会话
func (s *CollaborationTriggerService) GetUserPendingSessions(
	ctx context.Context,
	tenantID, userID string,
	limit int,
) ([]*entity.CollaborationSession, error) {
	if limit <= 0 {
		limit = 10
	}

	sessions, err := s.collabRepo.GetPendingSessionsByUser(ctx, tenantID, userID, limit)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrCollaborationQueryFailedCode, errorx.KV("error", err.Error()))
	}

	return sessions, nil
}

// GetSession 获取协同会话详情
func (s *CollaborationTriggerService) GetSession(
	ctx context.Context,
	sessionID string,
) (*entity.CollaborationSession, error) {
	session, err := s.collabRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrCollaborationNotFoundCode, errorx.KV("session_id", sessionID))
	}

	return session, nil
}

// CloseSession 关闭协同会话
func (s *CollaborationTriggerService) CloseSession(
	ctx context.Context,
	sessionID string,
	resolution string,
	humanResponse string,
) error {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.Status = "completed"
	session.Resolution = resolution
	session.HumanResponse = humanResponse
	session.CompletedAt = time.Now()
	session.UpdatedAt = time.Now()

	if err := s.collabRepo.UpdateSession(ctx, session); err != nil {
		return errorx.WrapByCode(err, errno.ErrCollaborationUpdateFailedCode, errorx.KV("error", err.Error()))
	}

	s.logger.Info("collaboration session closed",
		zap.String("session_id", sessionID),
		zap.String("resolution", resolution))

	return nil
}

// AssignSession 分配协同会话给处理人
func (s *CollaborationTriggerService) AssignSession(
	ctx context.Context,
	sessionID string,
	assigneeID string,
) error {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.AssignedTo = assigneeID
	session.Status = "assigned"
	session.UpdatedAt = time.Now()

	if err := s.collabRepo.UpdateSession(ctx, session); err != nil {
		return errorx.WrapByCode(err, errno.ErrCollaborationUpdateFailedCode, errorx.KV("error", err.Error()))
	}

	s.logger.Info("collaboration session assigned",
		zap.String("session_id", sessionID),
		zap.String("assignee_id", assigneeID))

	return nil
}

// GetCollaborationStatistics 获取协同统计信息
func (s *CollaborationTriggerService) GetCollaborationStatistics(
	ctx context.Context,
	tenantID string,
	startTime, endTime time.Time,
) (*CollaborationStatistics, error) {
	sessions, err := s.collabRepo.GetSessionsByTimeRange(ctx, tenantID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	stats := &CollaborationStatistics{
		TotalSessions:      len(sessions),
		PendingSessions:    0,
		AssignedSessions:   0,
		CompletedSessions:  0,
		ExpiredSessions:    0,
		TriggerTypeDistribution: make(map[string]int),
		ResolutionDistribution:  make(map[string]int),
	}

	for _, session := range sessions {
		// 统计状态
		switch session.Status {
		case "pending":
			stats.PendingSessions++
		case "assigned":
			stats.AssignedSessions++
		case "completed":
			stats.CompletedSessions++
		case "expired":
			stats.ExpiredSessions++
		}

		// 统计触发类型
		stats.TriggerTypeDistribution[session.TriggerType]++

		// 统计处理结果
		if session.Resolution != "" {
			stats.ResolutionDistribution[session.Resolution]++
		}
	}

	return stats, nil
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (
		s[:len(substr)] == substr ||
		s[len(s)-len(substr):] == substr ||
		containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// CollaborationRequest 协同请求
type CollaborationRequest struct {
	TenantID       string
	UserID         string
	ConversationID string
	MessageID      string
	AgentID        string
	InputText      string
	ResponseText   string
	Confidence     float64
	UserInitiated  bool
	HasError       bool
	ResponseTime   int64 // 毫秒
	Metadata       map[string]string
}

// CollaborationStatistics 协同统计信息
type CollaborationStatistics struct {
	TotalSessions           int            `json:"total_sessions"`
	PendingSessions         int            `json:"pending_sessions"`
	AssignedSessions        int            `json:"assigned_sessions"`
	CompletedSessions       int            `json:"completed_sessions"`
	ExpiredSessions         int            `json:"expired_sessions"`
	TriggerTypeDistribution map[string]int `json:"trigger_type_distribution"`
	ResolutionDistribution  map[string]int `json:"resolution_distribution"`
}
