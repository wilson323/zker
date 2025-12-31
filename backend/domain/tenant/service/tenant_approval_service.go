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
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// ApprovalStatus 审核状态
type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "pending"  // 待审核
	ApprovalStatusApproved ApprovalStatus = "approved" // 已通过
	ApprovalStatusRejected ApprovalStatus = "rejected" // 已拒绝
)

// TenantApprovalRequest 租户审核请求
type TenantApprovalRequest struct {
	TenantID      string          `json:"tenant_id" validate:"required"`
	Status        ApprovalStatus  `json:"status" validate:"required,oneof=pending approved rejected"`
	ReviewerID    string          `json:"reviewer_id" validate:"required"`
	Reason        string          `json:"reason,omitempty"`
	AutoApprove   bool            `json:"auto_approve"`   // 是否自动审核
	ApprovalRules *ApprovalRules  `json:"approval_rules,omitempty"` // 自动审核规则
}

// ApprovalRules 自动审核规则
type ApprovalRules struct {
	CheckEmailDomain    bool     `json:"check_email_domain"`     // 检查邮箱域名
	AllowedDomains      []string `json:"allowed_domains"`        // 白名单域名
	CheckCompanyName    bool     `json:"check_company_name"`     // 检查企业名称
	BannedKeywords      []string `json:"banned_keywords"`        // 禁用关键词
	RequireVerification bool     `json:"require_verification"`  // 需要人工验证
}

// TenantApprovalRecord 租户审核记录
type TenantApprovalRecord struct {
	ID           int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID     string          `json:"tenant_id" gorm:"type:varchar(36);not null;index"`
	Status       ApprovalStatus  `json:"status" gorm:"type:enum('pending','approved','rejected');not null"`
	ReviewerID   string          `json:"reviewer_id" gorm:"type:varchar(36)"`
	Reason       string          `json:"reason" gorm:"type:text"`
	ReviewedAt   *time.Time      `json:"reviewed_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"not null"`
}

// TableName 指定表名
func (TenantApprovalRecord) TableName() string {
	return "tenant_approval_records"
}

// TenantApprovalService 租户审核服务
// 负责处理租户注册审核,包括自动审核和人工审核
//
// 使用示例:
//
//	approvalSvc := service.NewTenantApprovalService(db, tenantRepo, emailSvc)
//	err := approvalSvc.ApproveTenant(ctx, &service.TenantApprovalRequest{...})
type TenantApprovalService struct {
	db        *gorm.DB
	tenantRepo repository.TenantRepository
	emailSvc  *EmailService
}

// NewTenantApprovalService 创建租户审核服务实例
func NewTenantApprovalService(
	db *gorm.DB,
	tenantRepo repository.TenantRepository,
	emailSvc *EmailService,
) *TenantApprovalService {
	return &TenantApprovalService{
		db:         db,
		tenantRepo: tenantRepo,
		emailSvc:   emailSvc,
	}
}

// ApproveTenant 审核租户
//
// 功能:
// 1. 自动审核(根据规则)
// 2. 人工审核(管理员操作)
// 3. 审核通过后激活租户
// 4. 审核拒绝后标记租户状态
// 5. 发送审核结果通知
//
// 参数:
//   - ctx: 上下文
//   - req: 审核请求
//
// 返回:
//   - error: 审核失败时返回错误
func (s *TenantApprovalService) ApproveTenant(ctx context.Context, req *TenantApprovalRequest) error {
	// 1. 参数验证
	if err := s.validateRequest(req); err != nil {
		return fmt.Errorf("request validation failed: %w", err)
	}

	// 2. 查询租户信息
	tenant, err := s.tenantRepo.GetByID(ctx, req.TenantID)
	if err != nil {
		return fmt.Errorf("tenant not found: %w", err)
	}

	// 3. 检查是否已经审核过
	hasApproved, err := s.hasApproved(ctx, req.TenantID)
	if err != nil {
		return fmt.Errorf("failed to check approval status: %w", err)
	}
	if hasApproved {
		return fmt.Errorf("tenant has already been reviewed")
	}

	// 4. 执行审核(自动审核或人工审核)
	approvalStatus := req.Status
	reviewerID := req.ReviewerID
	reviewedAt := time.Now()
	reason := req.Reason

	if req.AutoApprove {
		// 自动审核
		approvalStatus, reason = s.autoApprove(ctx, tenant, req.ApprovalRules)
		if approvalStatus == ApprovalStatusPending {
			// 自动审核无法决策,转为人工审核
			logs.Infof("Auto-approval pending for tenant %s, requires manual review", req.TenantID)
			return s.createPendingRecord(ctx, req.TenantID)
		}
		reviewerID = "system"
	}

	// 5. 使用事务更新租户状态和创建审核记录
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 5.1 更新租户状态
		if approvalStatus == ApprovalStatusApproved {
			tenant.Status = entity.TenantStatusActive
			logs.Infof("Tenant %s approved and activated", req.TenantID)
		} else if approvalStatus == ApprovalStatusRejected {
			tenant.Status = entity.TenantStatusSuspended
			logs.Infof("Tenant %s rejected and suspended", req.TenantID)
		} else {
			// 待审核,不更新状态
			logs.Infof("Tenant %s pending review", req.TenantID)
		}

		if approvalStatus != ApprovalStatusPending {
			tenant.UpdatedAt = time.Now().UnixMilli()
			if err := tx.Save(tenant).Error; err != nil {
				logs.Errorf("Failed to update tenant status: %v", err)
				return fmt.Errorf("failed to update tenant status: %w", err)
			}
		}

		// 5.2 创建审核记录
		record := &TenantApprovalRecord{
			TenantID:   req.TenantID,
			Status:     approvalStatus,
			ReviewerID: reviewerID,
			Reason:     reason,
			ReviewedAt: &reviewedAt,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := tx.Create(record).Error; err != nil {
			logs.Errorf("Failed to create approval record: %v", err)
			return fmt.Errorf("failed to create approval record: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 6. 发送审核结果通知(异步)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if approvalStatus == ApprovalStatusApproved {
			if err := s.emailSvc.SendApprovalEmail(ctx, tenant.ContactEmail, tenant.TenantName); err != nil {
				logs.Warnf("Failed to send approval email to %s: %v", tenant.ContactEmail, err)
			}
		} else if approvalStatus == ApprovalStatusRejected {
			if err := s.emailSvc.SendRejectionEmail(ctx, tenant.ContactEmail, tenant.TenantName, reason); err != nil {
				logs.Warnf("Failed to send rejection email to %s: %v", tenant.ContactEmail, err)
			}
		}
	}()

	logs.Infof("Tenant approval completed: %s -> %s", req.TenantID, approvalStatus)
	return nil
}

// autoApprove 自动审核
func (s *TenantApprovalService) autoApprove(ctx context.Context, tenant *entity.Tenant, rules *ApprovalRules) (ApprovalStatus, string) {
	if rules == nil {
		// 默认规则:全部转人工审核
		return ApprovalStatusPending, "pending manual review"
	}

	// 1. 检查禁用关键词
	if rules.CheckCompanyName {
		for _, keyword := range rules.BannedKeywords {
			if strings.Contains(tenant.TenantName, keyword) {
				return ApprovalStatusRejected, fmt.Sprintf("企业名称包含禁用词: %s", keyword)
			}
		}
	}

	// 2. 检查邮箱域名白名单
	if rules.CheckEmailDomain && len(rules.AllowedDomains) > 0 {
		domain := extractDomain(tenant.ContactEmail)
		if !contains(rules.AllowedDomains, domain) {
			return ApprovalStatusRejected, fmt.Sprintf("邮箱域名不在白名单中: %s", domain)
		}
	}

	// 3. 如果需要人工验证,转人工审核
	if rules.RequireVerification {
		return ApprovalStatusPending, "requires manual verification"
	}

	// 4. 自动通过
	return ApprovalStatusApproved, "auto-approved"
}

// GetApprovalStatus 查询审核状态
func (s *TenantApprovalService) GetApprovalStatus(ctx context.Context, tenantID string) (*TenantApprovalRecord, error) {
	var record TenantApprovalRecord
	err := s.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		First(&record).Error

	if err != nil {
		return nil, fmt.Errorf("approval record not found: %w", err)
	}

	return &record, nil
}

// ListPendingApprovals 列出待审核租户
func (s *TenantApprovalService) ListPendingApprovals(ctx context.Context, page, pageSize int) ([]*TenantApprovalRecord, int64, error) {
	var records []*TenantApprovalRecord
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	err := s.db.WithContext(ctx).
		Model(&TenantApprovalRecord{}).
		Where("status = ?", ApprovalStatusPending).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = s.db.WithContext(ctx).
		Where("status = ?", ApprovalStatusPending).
		Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// validateRequest 验证审核请求
func (s *TenantApprovalService) validateRequest(req *TenantApprovalRequest) error {
	if req == nil {
		return fmt.Errorf("请求不能为空")
	}

	if req.TenantID == "" {
		return fmt.Errorf("租户ID不能为空")
	}

	if req.Status == "" {
		return fmt.Errorf("审核状态不能为空")
	}

	if req.Status != ApprovalStatusPending &&
	   req.Status != ApprovalStatusApproved &&
	   req.Status != ApprovalStatusRejected {
		return fmt.Errorf("无效的审核状态: %s", req.Status)
	}

	if !req.AutoApprove && req.ReviewerID == "" {
		return fmt.Errorf("人工审核必须指定审核人ID")
	}

	return nil
}

// hasApproved 检查是否已经审核过
func (s *TenantApprovalService) hasApproved(ctx context.Context, tenantID string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&TenantApprovalRecord{}).
		Where("tenant_id = ? AND status IN (?)", tenantID, []ApprovalStatus{ApprovalStatusApproved, ApprovalStatusRejected}).
		Count(&count).Error

	return count > 0, err
}

// createPendingRecord 创建待审核记录
func (s *TenantApprovalService) createPendingRecord(ctx context.Context, tenantID string) error {
	record := &TenantApprovalRecord{
		TenantID:  tenantID,
		Status:    ApprovalStatusPending,
		ReviewerID: "system",
		Reason:    "auto-approval pending, requires manual review",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.db.WithContext(ctx).Create(record).Error
}

// 辅助函数

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func extractDomain(email string) string {
	for i := len(email) - 1; i >= 0; i-- {
		if email[i] == '@' {
			return email[i+1:]
		}
	}
	return ""
}
