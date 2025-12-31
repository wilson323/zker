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
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
	"github.com/coze-dev/coze-studio/backend/domain/developer/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/security"
)

// APIKeyManagementService API密钥管理服务接口
type APIKeyManagementService interface {
	// CreateAPIKey 创建API密钥
	CreateAPIKey(ctx context.Context, req *CreateAPIKeyRequest) (*entity.APIKey, string, error)

	// GetAPIKey 获取API密钥详情
	GetAPIKey(ctx context.Context, keyID string) (*entity.APIKey, error)

	// ListAPIKeys 查询API密钥列表
	ListAPIKeys(ctx context.Context, filter *APIKeyListFilter) ([]*entity.APIKey, int64, error)

	// RevokeAPIKey 撤销API密钥
	RevokeAPIKey(ctx context.Context, keyID string) error

	// DeleteAPIKey 删除API密钥
	DeleteAPIKey(ctx context.Context, keyID string) error

	// ValidateAPIKey 验证API密钥
	ValidateAPIKey(ctx context.Context, keySecret string) (*entity.APIKey, error)

	// RegenerateAPIKey 重新生成API密钥
	RegenerateAPIKey(ctx context.Context, keyID string) (*entity.APIKey, string, error)

	// GetProjectAPIKeys 获取项目的API密钥列表
	GetProjectAPIKeys(ctx context.Context, projectID string) ([]*entity.APIKey, error)

	// UpdateLastUsed 更新最后使用时间
	UpdateLastUsed(ctx context.Context, keyID string) error

	// GetExpiringKeys 获取即将过期的密钥
	GetExpiringKeys(ctx context.Context, tenantID string) ([]*entity.APIKey, error)
}

// CreateAPIKeyRequest 创建API密钥请求
type CreateAPIKeyRequest struct {
	TenantID  string                `json:"tenant_id" validate:"required"`
	ProjectID string                `json:"project_id" validate:"required"`
	KeyName   string                `json:"key_name" validate:"required,max=100"`
	KeyPrefix entity.APIKeyPrefix   `json:"key_prefix" validate:"required"`
	Scopes    entity.APIKeyScopes   `json:"scopes" validate:"required"`
	ExpiresIn *int                  `json:"expires_in"` // 过期时间（天数），nil表示永不过期
}

// APIKeyListFilter API密钥列表过滤器
type APIKeyListFilter struct {
	TenantID  string                 `validate:"required"`
	ProjectID string
	KeyPrefix entity.APIKeyPrefix
	Status    entity.APIKeyStatus
	PageToken string
	PageSize  int
}

// apiKeyManagementService API密钥管理服务实现
type apiKeyManagementService struct {
	apiKeyRepo  repository.APIKeyRepository
	projectRepo repository.ProjectRepository
	encryptor   security.EncryptionService
}

// NewAPIKeyManagementService 创建API密钥管理服务实例
func NewAPIKeyManagementService(
	apiKeyRepo repository.APIKeyRepository,
	projectRepo repository.ProjectRepository,
	secretKey string,
) APIKeyManagementService {
	// 创建AES-256-GCM加密器
	encryptor, err := security.NewAESGCMEncryptor(secretKey)
	if err != nil {
		// 如果加密器创建失败，使用环境变量中的密钥
		envKey := os.Getenv("API_KEY_ENCRYPTION_KEY")
		if envKey == "" {
			// 如果环境变量也没有，生成新密钥
			generatedKey, _ := security.GenerateEncryptionKey()
			envKey = generatedKey
			fmt.Printf("⚠️  Generated new encryption key: %s (please save it to environment variable API_KEY_ENCRYPTION_KEY)\n", envKey)
		}
		encryptor, _ = security.NewAESGCMEncryptor(envKey)
	}

	return &apiKeyManagementService{
		apiKeyRepo:  apiKeyRepo,
		projectRepo: projectRepo,
		encryptor:   encryptor,
	}
}

// CreateAPIKey 创建API密钥
func (s *apiKeyManagementService) CreateAPIKey(ctx context.Context, req *CreateAPIKeyRequest) (*entity.APIKey, string, error) {
	// 验证项目是否存在
	project, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, "", err
	}
	if project == nil {
		return nil, "", errors.New("project not found")
	}

	// 验证项目属于该租户
	if project.TenantID != req.TenantID {
		return nil, "", errors.New("project does not belong to tenant")
	}

	// 生成API密钥
	keySecret := s.generateKeySecret(req.KeyPrefix)

	// 使用AES-256-GCM加密密钥
	encryptedSecret, err := s.encryptor.Encrypt(keySecret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt api key: %w", err)
	}

	// 计算过期时间
	var expiresAt *int64
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		expiry := time.Now().AddDate(0, 0, *req.ExpiresIn).UnixMilli()
		expiresAt = &expiry
	}

	// 生成脱敏密钥（用于显示）
	keyMasked := s.maskKeySecret(keySecret)

	apiKey := &entity.APIKey{
		KeyID:     generateID("key"),
		TenantID:  req.TenantID,
		ProjectID: req.ProjectID,
		KeyName:   req.KeyName,
		KeySecret: encryptedSecret,
		KeyPrefix: req.KeyPrefix,
		KeyMasked: keyMasked,
		Scopes:    req.Scopes,
		ExpiresAt: expiresAt,
		Status:    entity.APIKeyStatusActive,
	}

	if err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, "", err
	}

	return apiKey, keySecret, nil
}

// GetAPIKey 获取API密钥详情
func (s *apiKeyManagementService) GetAPIKey(ctx context.Context, keyID string) (*entity.APIKey, error) {
	return s.apiKeyRepo.GetByID(ctx, keyID)
}

// ListAPIKeys 查询API密钥列表
func (s *apiKeyManagementService) ListAPIKeys(ctx context.Context, filter *APIKeyListFilter) ([]*entity.APIKey, int64, error) {
	repoFilter := &repository.APIKeyFilter{
		TenantID:  filter.TenantID,
		ProjectID: filter.ProjectID,
		KeyPrefix: filter.KeyPrefix,
		Status:    filter.Status,
		PageToken: filter.PageToken,
		PageSize:  filter.PageSize,
	}

	return s.apiKeyRepo.List(ctx, repoFilter)
}

// RevokeAPIKey 撤销API密钥
func (s *apiKeyManagementService) RevokeAPIKey(ctx context.Context, keyID string) error {
	return s.apiKeyRepo.Revoke(ctx, keyID)
}

// DeleteAPIKey 删除API密钥
func (s *apiKeyManagementService) DeleteAPIKey(ctx context.Context, keyID string) error {
	return s.apiKeyRepo.Delete(ctx, keyID)
}

// ValidateAPIKey 验证API密钥
func (s *apiKeyManagementService) ValidateAPIKey(ctx context.Context, keySecret string) (*entity.APIKey, error) {
	// 获取所有API密钥并逐个解密验证（性能优化：可添加哈希索引）
	// 注意：这里简化处理，生产环境建议在KeySecret字段添加哈希索引
	apiKeys, err := s.apiKeyRepo.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}

	for _, apiKey := range apiKeys {
		// 解密存储的密钥
		decryptedSecret, err := s.encryptor.Decrypt(apiKey.KeySecret)
		if err != nil {
			// 解密失败，跳过该密钥
			continue
		}

		// 验证密钥是否匹配
		if decryptedSecret == keySecret {
			// 检查密钥是否激活
			if !apiKey.IsActive() {
				return nil, errors.New("api key is not active")
			}

			// 检查是否过期
			if apiKey.IsExpired() {
				// 自动标记为过期
				_ = s.apiKeyRepo.Expire(ctx, apiKey.KeyID)
				return nil, errors.New("api key has expired")
			}

			return apiKey, nil
		}
	}

	return nil, errors.New("invalid api key")
}

// RegenerateAPIKey 重新生成API密钥
func (s *apiKeyManagementService) RegenerateAPIKey(ctx context.Context, keyID string) (*entity.APIKey, string, error) {
	// 获取现有密钥
	apiKey, err := s.apiKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return nil, "", err
	}
	if apiKey == nil {
		return nil, "", errors.New("api key not found")
	}

	// 生成新密钥
	newKeySecret := s.generateKeySecret(apiKey.KeyPrefix)

	// 使用AES-256-GCM加密新密钥
	encryptedSecret, err := s.encryptor.Encrypt(newKeySecret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt api key: %w", err)
	}

	// 更新密钥
	apiKey.KeySecret = encryptedSecret
	apiKey.KeyMasked = s.maskKeySecret(newKeySecret)

	if err := s.apiKeyRepo.Update(ctx, apiKey); err != nil {
		return nil, "", err
	}

	return apiKey, newKeySecret, nil
}

// GetProjectAPIKeys 获取项目的API密钥列表
func (s *apiKeyManagementService) GetProjectAPIKeys(ctx context.Context, projectID string) ([]*entity.APIKey, error) {
	return s.apiKeyRepo.GetByProjectID(ctx, projectID)
}

// UpdateLastUsed 更新最后使用时间
func (s *apiKeyManagementService) UpdateLastUsed(ctx context.Context, keyID string) error {
	return s.apiKeyRepo.UpdateLastUsedAt(ctx, keyID)
}

// GetExpiringKeys 获取即将过期的密钥
func (s *apiKeyManagementService) GetExpiringKeys(ctx context.Context, tenantID string) ([]*entity.APIKey, error) {
	return s.apiKeyRepo.GetExpiringKeys(ctx, tenantID)
}

// generateKeySecret 生成API密钥
func (s *apiKeyManagementService) generateKeySecret(prefix entity.APIKeyPrefix) string {
	// 生成32字节随机数据
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选
		randomBytes = []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	}

	// 编码为base64
	encoded := base64.StdEncoding.EncodeToString(randomBytes)

	// 格式：{prefix}_{encoded}
	return fmt.Sprintf("%s_%s", prefix, encoded)
}

// maskKeySecret 生成脱敏密钥（用于显示）
func (s *apiKeyManagementService) maskKeySecret(keySecret string) string {
	if len(keySecret) <= 8 {
		return keySecret[:2] + "****"
	}

	// 保留前缀和后4位
	prefixEnd := 3
	if len(keySecret) < prefixEnd+4 {
		prefixEnd = len(keySecret) - 4
	}

	return keySecret[:prefixEnd] + "****" + keySecret[len(keySecret)-4:]
}
