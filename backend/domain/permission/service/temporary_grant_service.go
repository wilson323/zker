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
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
)

// TemporaryGrantService 临时授权服务
type TemporaryGrantService struct {
	grantRepo     repository.TemporaryGrantRepository
	historyRepo   repository.TemporaryGrantHistoryRepository
	userRoleRepo  repository.UserRoleRepository
}

// NewTemporaryGrantService 创建临时授权服务实例
func NewTemporaryGrantService(
	grantRepo repository.TemporaryGrantRepository,
	historyRepo repository.TemporaryGrantHistoryRepository,
	userRoleRepo repository.UserRoleRepository,
) *TemporaryGrantService {
	return &TemporaryGrantService{
		grantRepo:    grantRepo,
		historyRepo:  historyRepo,
		userRoleRepo: userRoleRepo,
	}
}

// CreateTemporaryGrantRequest 创建临时授权请求
type CreateTemporaryGrantRequest struct {
	TenantID        string                     `json:"tenant_id" validate:"required"`
	GranteeID       string                     `json:"grantee_id" validate:"required"`    // 被授权人用户ID
	GrantorID       string                     `json:"grantor_id" validate:"required"`    // 授权人用户ID
	PermissionType  entity.PermissionType      `json:"permission_type" validate:"required"` // 权限类型
	PermissionData  entity.PermissionData      `json:"permission_data" validate:"required"` // 权限数据（JSON）
	Duration        time.Duration              `json:"duration" validate:"required,min=1m,max=720h"` // 有效期（例如：24*time.Hour）
	Reason          string                     `json:"reason" validate:"omitempty,max=500"`
	RequestID       string                     `json:"request_id" validate:"omitempty,uuid"`
}

// CreateTemporaryGrant 创建临时授权
func (s *TemporaryGrantService) CreateTemporaryGrant(
	ctx context.Context,
	req *CreateTemporaryGrantRequest,
) (*entity.TemporaryGrant, error) {
	// 1. 生成授权码
	grantCode, err := s.generateGrantCode()
	if err != nil {
		return nil, fmt.Errorf("生成授权码失败: %w", err)
	}

	// 2. 序列化权限数据
	permDataJSON, err := json.Marshal(req.PermissionData)
	if err != nil {
		return nil, fmt.Errorf("序列化权限数据失败: %w", err)
	}

	// 3. 计算过期时间
	expiresAt := time.Now().Add(req.Duration).UnixMilli()

	// 4. 创建临时授权实体
	grant := &entity.TemporaryGrant{
		GrantCode:      grantCode,
		GranteeID:      req.GranteeID,
		GrantorID:      req.GrantorID,
		TenantID:       req.TenantID,
		PermissionType: req.PermissionType,
		PermissionData: entity.PermissionData{},
		ExpiresAt:      expiresAt,
		Reason:         req.Reason,
		RequestID:      req.RequestID,
	}

	// 反序列化权限数据到实体
	if err := json.Unmarshal(permDataJSON, &grant.PermissionData); err != nil {
		return nil, fmt.Errorf("反序列化权限数据失败: %w", err)
	}

	// 5. 保存到数据库
	if err := s.grantRepo.Create(ctx, grant); err != nil {
		return nil, fmt.Errorf("保存临时授权失败: %w", err)
	}

	// 6. 记录历史
	s.recordHistory(ctx, grant, entity.ActionTypeCreated, req.GrantorID, "", "创建临时授权: "+req.Reason)

	return grant, nil
}

// UseTemporaryGrant 使用临时授权码
func (s *TemporaryGrantService) UseTemporaryGrant(
	ctx context.Context,
	grantCode string,
	userID string,
) (*entity.UserRole, error) {
	// 1. 查询临时授权
	grant, err := s.grantRepo.GetByGrantCode(ctx, grantCode)
	if err != nil {
		return nil, fmt.Errorf("查询授权码失败: %w", err)
	}
	if grant == nil {
		return nil, fmt.Errorf("授权码不存在")
	}

	// 2. 验证授权有效性
	if !grant.IsValid() {
		if grant.IsExpired() {
			return nil, fmt.Errorf("授权码已过期")
		}
		if grant.IsUsed {
			return nil, fmt.Errorf("授权码已被使用")
		}
		if grant.IsRevoked {
			return nil, fmt.Errorf("授权码已被撤销")
		}
	}

	// 3. 验证被授权人
	if grant.GranteeID != userID {
		return nil, fmt.Errorf("授权码不属于当前用户")
	}

	// 4. 根据权限类型创建UserRole
	switch grant.PermissionType {
	case entity.PermissionTypeRole:
		return s.useRoleGrant(ctx, grant)

	case entity.PermissionTypeDataPermission:
		return nil, fmt.Errorf("数据权限临时授予功能暂未实现")

	case entity.PermissionTypeFieldPermission:
		return nil, fmt.Errorf("字段权限临时授予功能暂未实现")

	default:
		return nil, fmt.Errorf("无效的权限类型: %s", grant.PermissionType)
	}
}

// useRoleGrant 使用角色授权
func (s *TemporaryGrantService) useRoleGrant(
	ctx context.Context,
	grant *entity.TemporaryGrant,
) (*entity.UserRole, error) {
	// 1. 解析权限数据
	roleID, ok := grant.PermissionData["role_id"].(string)
	if !ok || roleID == "" {
		return nil, fmt.Errorf("无效的角色ID")
	}

	// 2. 检查是否已经分配该角色
	exists, err := s.userRoleRepo.Exists(ctx, grant.GranteeID, grant.TenantID, roleID)
	if err != nil {
		return nil, fmt.Errorf("检查角色分配失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("用户已拥有该角色")
	}

	// 3. 创建UserRole（带过期时间）
	userRole := &entity.UserRole{
		UserID:    grant.GranteeID,
		TenantID:  grant.TenantID,
		RoleID:    roleID,
		ExpiresAt: &grant.ExpiresAt, // 设置过期时间
		CreatedAt: time.Now().UnixMilli(),
	}

	// 4. 保存UserRole
	if err := s.userRoleRepo.Create(ctx, userRole); err != nil {
		return nil, fmt.Errorf("分配角色失败: %w", err)
	}

	// 5. 标记授权码为已使用
	grant.MarkAsUsed()
	if err := s.grantRepo.Update(ctx, grant); err != nil {
		return nil, fmt.Errorf("更新授权码状态失败: %w", err)
	}

	// 6. 记录历史
	s.recordHistory(ctx, grant, entity.ActionTypeUsed, grant.GranteeID, "", "使用授权码")

	return userRole, nil
}

// RevokeTemporaryGrant 撤销临时授权
func (s *TemporaryGrantService) RevokeTemporaryGrant(
	ctx context.Context,
	grantCode string,
	operatorID string,
	reason string,
) error {
	// 1. 查询临时授权
	grant, err := s.grantRepo.GetByGrantCode(ctx, grantCode)
	if err != nil {
		return fmt.Errorf("查询授权码失败: %w", err)
	}
	if grant == nil {
		return fmt.Errorf("授权码不存在")
	}

	// 2. 检查是否已使用
	if grant.IsUsed {
		return fmt.Errorf("授权码已被使用，无法撤销")
	}

	// 3. 检查是否已撤销
	if grant.IsRevoked {
		return fmt.Errorf("授权码已被撤销")
	}

	// 4. 标记为已撤销
	grant.MarkAsRevoked(reason)
	if err := s.grantRepo.Update(ctx, grant); err != nil {
		return fmt.Errorf("撤销授权码失败: %w", err)
	}

	// 5. 记录历史
	s.recordHistory(ctx, grant, entity.ActionTypeRevoked, operatorID, "", reason)

	return nil
}

// CleanupExpiredGrants 清理过期的临时授权
func (s *TemporaryGrantService) CleanupExpiredGrants(ctx context.Context) (int64, error) {
	now := time.Now().UnixMilli()

	// 1. 查询所有过期的授权
	expiredGrants, err := s.grantRepo.ListExpired(ctx, now)
	if err != nil {
		return 0, fmt.Errorf("查询过期授权失败: %w", err)
	}

	count := int64(len(expiredGrants))

	// 2. 处理每个过期授权
	for _, grant := range expiredGrants {
		// 删除过期的UserRole（如果是角色类型且已使用）
		// TODO: 实现删除过期UserRole的功能
		// if grant.PermissionType == entity.PermissionTypeRole && grant.IsUsed {
		//     roleID, ok := grant.PermissionData["role_id"].(string)
		//     if ok {
		//         s.userRoleRepo.DeleteExpiredUserRole(ctx, grant.GranteeID, grant.TenantID, roleID)
		//     }
		// }

		// 记录历史
		s.recordHistory(ctx, grant, entity.ActionTypeExpired, "system", "", "自动清理过期授权")
	}

	// 3. 删除过期的临时授权
	deletedCount, err := s.grantRepo.DeleteExpired(ctx, now)
	if err != nil {
		return count, fmt.Errorf("删除过期授权失败: %w", err)
	}

	return count + deletedCount, nil
}

// GetTemporaryGrant 获取临时授权信息
func (s *TemporaryGrantService) GetTemporaryGrant(
	ctx context.Context,
	grantCode string,
) (*entity.TemporaryGrant, error) {
	grant, err := s.grantRepo.GetByGrantCode(ctx, grantCode)
	if err != nil {
		return nil, fmt.Errorf("查询授权码失败: %w", err)
	}
	if grant == nil {
		return nil, fmt.Errorf("授权码不存在")
	}

	return grant, nil
}

// ListTemporaryGrants 查询临时授权列表
func (s *TemporaryGrantService) ListTemporaryGrants(
	ctx context.Context,
	filter *repository.TemporaryGrantFilter,
) ([]*entity.TemporaryGrant, int64, error) {
	grants, total, err := s.grantRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("查询临时授权列表失败: %w", err)
	}

	return grants, total, nil
}

// GetGrantHistory 获取授权历史
func (s *TemporaryGrantService) GetGrantHistory(
	ctx context.Context,
	grantCode string,
) ([]*entity.TemporaryGrantHistory, error) {
	histories, err := s.historyRepo.GetByGrantCode(ctx, grantCode)
	if err != nil {
		return nil, fmt.Errorf("查询授权历史失败: %w", err)
	}

	return histories, nil
}

// CleanupExpiredHistory 清理过期的历史记录
func (s *TemporaryGrantService) CleanupExpiredHistory(ctx context.Context, before int64) (int64, error) {
	count, err := s.historyRepo.DeleteExpired(ctx, before)
	if err != nil {
		return 0, fmt.Errorf("清理过期历史失败: %w", err)
	}

	return count, nil
}

// generateGrantCode 生成授权码
// 格式: XXXX-XXXX-XXXX-XXXX（便于人工输入）
func (s *TemporaryGrantService) generateGrantCode() (string, error) {
	// 生成随机字节（16字节 = 128位）
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// Base32编码（便于人工输入，避免歧义字符）
	encoded := base32HexEncoding.EncodeToString(b)

	// 格式化为 XXXX-XXXX-XXXX-XXXX
	var codeBuilder strings.Builder
	for i := 0; i < 4; i++ {
		start := i * 4
		end := start + 4
		if end > len(encoded) {
			end = len(encoded)
		}
		codeBuilder.WriteString(encoded[start:end])
		if i < 3 {
			codeBuilder.WriteString("-")
		}
	}

	return codeBuilder.String(), nil
}

// recordHistory 记录历史
func (s *TemporaryGrantService) recordHistory(
	ctx context.Context,
	grant *entity.TemporaryGrant,
	actionType entity.ActionType,
	operatorID, operatorName, reason string,
) {
	// 创建快照
	snapshot := make(map[string]interface{})
	snapshot["grant_id"] = grant.GrantID
	snapshot["grant_code"] = grant.GrantCode
	snapshot["permission_type"] = grant.PermissionType
	snapshot["expires_at"] = grant.ExpiresAt
	snapshot["is_used"] = grant.IsUsed
	snapshot["is_revoked"] = grant.IsRevoked

	history := &entity.TemporaryGrantHistory{
		GrantID:       grant.GrantID,
		GrantCode:     grant.GrantCode,
		ActionType:    actionType,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
		SnapshotData:  snapshot,
		Reason:        reason,
	}

	// 异步记录历史（不阻塞主流程）
	go func() {
		_ = s.historyRepo.Create(context.Background(), history)
	}()
}

// base32HexEncoding 自定义Base32编码（去除歧义字符）
var base32HexEncoding = base64.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").WithPadding(base64.NoPadding)
