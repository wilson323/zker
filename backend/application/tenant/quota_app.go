// backend/application/tenant/quota_app.go
package tenant

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// BaseResponse 基础响应
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// CheckQuotaRequest 配额检查请求
type CheckQuotaRequest struct {
	TenantID     string `json:"tenant_id"`
	ResourceType string `json:"resource_type"`
}

// CheckQuotaResponse 配额检查响应
type CheckQuotaResponse struct {
	BaseResponse
	Data struct {
		Allowed      bool `json:"allowed"`
		CurrentUsage int  `json:"current_usage"`
		MaxLimit     int  `json:"max_limit"`
	} `json:"data"`
}

// QuotaAppService 配额应用服务
// 负责配额相关的应用层业务逻辑
type QuotaAppService struct {
	quotaRepo repository.QuotaRepository
}

// NewQuotaAppService 创建配额应用服务
func NewQuotaAppService(quotaRepo repository.QuotaRepository) *QuotaAppService {
	return &QuotaAppService{
		quotaRepo: quotaRepo,
	}
}

// CheckQuota 检查配额是否充足
// 返回值: (allowed, currentUsage, maxLimit, error)
func (s *QuotaAppService) CheckQuota(ctx context.Context, req *CheckQuotaRequest) (*CheckQuotaResponse, error) {
	// 1. 参数验证
	if req.TenantID == "" {
		logs.CtxWarnf(ctx, "[QuotaApp] tenant_id is empty")
		return nil, fmt.Errorf("tenant_id cannot be empty")
	}

	if req.ResourceType == "" {
		logs.CtxWarnf(ctx, "[QuotaApp] resource_type is empty")
		return nil, fmt.Errorf("resource_type cannot be empty")
	}

	// 2. 转换资源类型
	resourceType := entity.ResourceType(req.ResourceType)

	// 3. 获取配额信息
	quota, err := s.quotaRepo.GetByTenantAndResource(ctx, req.TenantID, resourceType)
	if err != nil {
		logs.CtxErrorf(ctx, "[QuotaApp] failed to get quota: tenant=%s resource=%s error=%v",
			req.TenantID, req.ResourceType, err)
		return nil, fmt.Errorf("failed to get quota: %w", err)
	}

	if quota == nil {
		logs.CtxWarnf(ctx, "[QuotaApp] quota not found: tenant=%s resource=%s",
			req.TenantID, req.ResourceType)
		// 配额不存在时，默认允许（兼容旧数据）
		return &CheckQuotaResponse{
			BaseResponse: BaseResponse{Code: 0, Message: "success"},
			Data: struct {
				Allowed      bool `json:"allowed"`
				CurrentUsage int  `json:"current_usage"`
				MaxLimit     int  `json:"max_limit"`
			}{
				Allowed:      true,
				CurrentUsage: 0,
				MaxLimit:     0,
			},
		}, nil
	}

	// 4. 检查是否无限制
	if quota.IsUnlimited() {
		logs.CtxDebugf(ctx, "[QuotaApp] unlimited quota: tenant=%s resource=%s",
			req.TenantID, req.ResourceType)
		return &CheckQuotaResponse{
			BaseResponse: BaseResponse{Code: 0, Message: "success"},
			Data: struct {
				Allowed      bool `json:"allowed"`
				CurrentUsage int  `json:"current_usage"`
				MaxLimit     int  `json:"max_limit"`
			}{
				Allowed:      true,
				CurrentUsage: quota.UsedCount,
				MaxLimit:     quota.MaxLimit,
			},
		}, nil
	}

	// 5. 检查配额是否充足
	allowed := quota.UsedCount < quota.MaxLimit

	logs.CtxInfof(ctx, "[QuotaApp] quota check: tenant=%s resource=%s allowed=%v usage=%d/%d",
		req.TenantID, req.ResourceType, allowed, quota.UsedCount, quota.MaxLimit)

	return &CheckQuotaResponse{
		BaseResponse: BaseResponse{Code: 0, Message: "success"},
		Data: struct {
			Allowed      bool `json:"allowed"`
			CurrentUsage int  `json:"current_usage"`
			MaxLimit     int  `json:"max_limit"`
		}{
			Allowed:      allowed,
			CurrentUsage: quota.UsedCount,
			MaxLimit:     quota.MaxLimit,
		},
	}, nil
}

// GetQuotaUsage 获取配额使用情况
func (s *QuotaAppService) GetQuotaUsage(ctx context.Context, tenantID, resourceType string) (*entity.Quota, error) {
	quota, err := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, entity.ResourceType(resourceType))
	if err != nil {
		return nil, err
	}
	return quota, nil
}

// GetAllQuotas 获取租户的所有配额
func (s *QuotaAppService) GetAllQuotas(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	quotas, err := s.quotaRepo.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return quotas, nil
}
