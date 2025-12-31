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
	"strings"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
)

// MockAPIKeyRepository Mock API密钥仓储
type MockAPIKeyRepository struct {
	apiKeys map[string]*entity.APIKey
}

func NewMockAPIKeyRepository() *MockAPIKeyRepository {
	return &MockAPIKeyRepository{
		apiKeys: make(map[string]*entity.APIKey),
	}
}

func (m *MockAPIKeyRepository) Create(ctx context.Context, apiKey *entity.APIKey) error {
	m.apiKeys[apiKey.KeyID] = apiKey
	return nil
}

func (m *MockAPIKeyRepository) GetByID(ctx context.Context, keyID string) (*entity.APIKey, error) {
	return m.apiKeys[keyID], nil
}

func (m *MockAPIKeyRepository) GetByProjectID(ctx context.Context, projectID string) ([]*entity.APIKey, error) {
	var result []*entity.APIKey
	for _, key := range m.apiKeys {
		if key.ProjectID == projectID {
			result = append(result, key)
		}
	}
	return result, nil
}

func (m *MockAPIKeyRepository) GetAllActive(ctx context.Context) ([]*entity.APIKey, error) {
	var result []*entity.APIKey
	for _, key := range m.apiKeys {
		if key.Status == entity.APIKeyStatusActive {
			result = append(result, key)
		}
	}
	return result, nil
}

func (m *MockAPIKeyRepository) Update(ctx context.Context, apiKey *entity.APIKey) error {
	m.apiKeys[apiKey.KeyID] = apiKey
	return nil
}

func (m *MockAPIKeyRepository) Delete(ctx context.Context, keyID string) error {
	delete(m.apiKeys, keyID)
	return nil
}

func (m *MockAPIKeyRepository) List(ctx context.Context, filter interface{}) ([]*entity.APIKey, int64, error) {
	return nil, 0, nil
}

func (m *MockAPIKeyRepository) Revoke(ctx context.Context, keyID string) error {
	if key, ok := m.apiKeys[keyID]; ok {
		key.Status = entity.APIKeyStatusRevoked
	}
	return nil
}

func (m *MockAPIKeyRepository) Expire(ctx context.Context, keyID string) error {
	if key, ok := m.apiKeys[keyID]; ok {
		key.Status = entity.APIKeyStatusExpired
	}
	return nil
}

func (m *MockAPIKeyRepository) UpdateLastUsedAt(ctx context.Context, keyID string) error {
	now := time.Now().UnixMilli()
	if key, ok := m.apiKeys[keyID]; ok {
		key.LastUsedAt = &now
	}
	return nil
}

func (m *MockAPIKeyRepository) GetExpiringKeys(ctx context.Context, tenantID string) ([]*entity.APIKey, error) {
	return nil, nil
}

func (m *MockAPIKeyRepository) CountByProjectID(ctx context.Context, projectID string) (int64, error) {
	return 0, nil
}

// MockProjectRepository Mock项目仓储
type MockProjectRepository struct {
	projects map[string]*entity.Project
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects: make(map[string]*entity.Project),
	}
}

func (m *MockProjectRepository) Create(ctx context.Context, project *entity.Project) error {
	m.projects[project.ProjectID] = project
	return nil
}

func (m *MockProjectRepository) GetByID(ctx context.Context, projectID string) (*entity.Project, error) {
	return m.projects[projectID], nil
}

func (m *MockProjectRepository) Update(ctx context.Context, project *entity.Project) error {
	m.projects[project.ProjectID] = project
	return nil
}

func (m *MockProjectRepository) Delete(ctx context.Context, projectID string) error {
	return nil
}

func (m *MockProjectRepository) List(ctx context.Context, filter interface{}) ([]*entity.Project, int64, error) {
	return nil, 0, nil
}

// TestCreateAPIKey 测试创建API密钥
func TestCreateAPIKey(t *testing.T) {
	apiKeyRepo := NewMockAPIKeyRepository()
	projectRepo := NewMockProjectRepository()

	// 添加测试项目
	project := &entity.Project{
		ProjectID: "proj_1",
		TenantID:  "tenant_1",
		ProjectName: "Test Project",
	}
	projectRepo.projects[project.ProjectID] = project

	service := NewAPIKeyManagementService(apiKeyRepo, projectRepo, "test_secret_key_32_bytes_long!")

	req := &CreateAPIKeyRequest{
		TenantID:  "tenant_1",
		ProjectID: "proj_1",
		KeyName:   "Test Key",
		KeyPrefix: entity.APIKeyPrefixSecret,
		Scopes:    entity.APIKeyScopes{entity.APIKeyScopeRead, entity.APIKeyScopeWrite},
		ExpiresIn: nil,
	}

	apiKey, keySecret, err := service.CreateAPIKey(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateAPIKey failed: %v", err)
	}

	// 验证返回的密钥
	if apiKey == nil {
		t.Fatal("apiKey should not be nil")
	}

	if keySecret == "" {
		t.Error("keySecret should not be empty")
	}

	if !strings.HasPrefix(keySecret, string(entity.APIKeyPrefixSecret)+"_") {
		t.Errorf("keySecret should have prefix 'sk_', got: %s", keySecret)
	}

	// 验证密钥已加密存储（应该不等于原始密钥）
	if apiKey.KeySecret == keySecret {
		t.Error("KeySecret should be encrypted, not stored in plain text")
	}

	// 验证脱敏密钥格式
	if !strings.Contains(apiKey.KeyMasked, "****") {
		t.Errorf("KeyMasked should contain '****', got: %s", apiKey.KeyMasked)
	}

	// 验证其他字段
	if apiKey.KeyName != req.KeyName {
		t.Errorf("KeyName mismatch: got %s, want %s", apiKey.KeyName, req.KeyName)
	}

	if apiKey.KeyPrefix != req.KeyPrefix {
		t.Errorf("KeyPrefix mismatch: got %s, want %s", apiKey.KeyPrefix, req.KeyPrefix)
	}

	if apiKey.Status != entity.APIKeyStatusActive {
		t.Errorf("Status should be active, got: %s", apiKey.Status)
	}
}

// TestValidateAPIKey 测试验证API密钥
func TestValidateAPIKey(t *testing.T) {
	apiKeyRepo := NewMockAPIKeyRepository()
	projectRepo := NewMockProjectRepository()

	// 添加测试项目
	project := &entity.Project{
		ProjectID: "proj_1",
		TenantID:  "tenant_1",
		ProjectName: "Test Project",
	}
	projectRepo.projects[project.ProjectID] = project

	service := NewAPIKeyManagementService(apiKeyRepo, projectRepo, "test_secret_key_32_bytes_long!")

	// 创建API密钥
	req := &CreateAPIKeyRequest{
		TenantID:  "tenant_1",
		ProjectID: "proj_1",
		KeyName:   "Test Key",
		KeyPrefix: entity.APIKeyPrefixSecret,
		Scopes:    entity.APIKeyScopes{entity.APIKeyScopeRead},
		ExpiresIn: nil,
	}

	_, keySecret, err := service.CreateAPIKey(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateAPIKey failed: %v", err)
	}

	// 验证正确的密钥
	apiKey, err := service.ValidateAPIKey(context.Background(), keySecret)
	if err != nil {
		t.Errorf("ValidateAPIKey with correct key failed: %v", err)
	}

	if apiKey == nil {
		t.Fatal("Validated apiKey should not be nil")
	}

	if apiKey.Status != entity.APIKeyStatusActive {
		t.Errorf("Validated key should be active, got: %s", apiKey.Status)
	}

	// 验证错误的密钥
	_, err = service.ValidateAPIKey(context.Background(), "invalid_key")
	if err == nil {
		t.Error("ValidateAPIKey with invalid key should return error")
	}

	if !strings.Contains(err.Error(), "invalid api key") {
		t.Errorf("Error should mention 'invalid api key', got: %v", err)
	}
}

// TestRevokeAPIKey 测试撤销API密钥
func TestRevokeAPIKey(t *testing.T) {
	apiKeyRepo := NewMockAPIKeyRepository()
	projectRepo := NewMockProjectRepository()

	// 添加测试项目
	project := &entity.Project{
		ProjectID: "proj_1",
		TenantID:  "tenant_1",
		ProjectName: "Test Project",
	}
	projectRepo.projects[project.ProjectID] = project

	service := NewAPIKeyManagementService(apiKeyRepo, projectRepo, "test_secret_key_32_bytes_long!")

	// 创建API密钥
	req := &CreateAPIKeyRequest{
		TenantID:  "tenant_1",
		ProjectID: "proj_1",
		KeyName:   "Test Key",
		KeyPrefix: entity.APIKeyPrefixSecret,
		Scopes:    entity.APIKeyScopes{entity.APIKeyScopeRead},
		ExpiresIn: nil,
	}

	apiKey, _, err := service.CreateAPIKey(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateAPIKey failed: %v", err)
	}

	// 撤销密钥
	err = service.RevokeAPIKey(context.Background(), apiKey.KeyID)
	if err != nil {
		t.Errorf("RevokeAPIKey failed: %v", err)
	}

	// 验证密钥已被撤销
	revokedKey, err := apiKeyRepo.GetByID(context.Background(), apiKey.KeyID)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if revokedKey.Status != entity.APIKeyStatusRevoked {
		t.Errorf("Key status should be revoked, got: %s", revokedKey.Status)
	}

	// 验证已撤销的密钥无法使用
	_, err = service.ValidateAPIKey(context.Background(), "test_key")
	if err == nil {
		t.Error("Revoked key should not be valid")
	}
}

// TestRegenerateAPIKey 测试重新生成API密钥
func TestRegenerateAPIKey(t *testing.T) {
	apiKeyRepo := NewMockAPIKeyRepository()
	projectRepo := NewMockProjectRepository()

	// 添加测试项目
	project := &entity.Project{
		ProjectID: "proj_1",
		TenantID:  "tenant_1",
		ProjectName: "Test Project",
	}
	projectRepo.projects[project.ProjectID] = project

	service := NewAPIKeyManagementService(apiKeyRepo, projectRepo, "test_secret_key_32_bytes_long!")

	// 创建API密钥
	req := &CreateAPIKeyRequest{
		TenantID:  "tenant_1",
		ProjectID: "proj_1",
		KeyName:   "Test Key",
		KeyPrefix: entity.APIKeyPrefixSecret,
		Scopes:    entity.APIKeyScopes{entity.APIKeyScopeRead},
		ExpiresIn: nil,
	}

	apiKey, originalSecret, err := service.CreateAPIKey(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateAPIKey failed: %v", err)
	}

	// 重新生成密钥
	newApiKey, newSecret, err := service.RegenerateAPIKey(context.Background(), apiKey.KeyID)
	if err != nil {
		t.Errorf("RegenerateAPIKey failed: %v", err)
	}

	if newSecret == originalSecret {
		t.Error("New secret should be different from old secret")
	}

	if newApiKey.KeyID != apiKey.KeyID {
		t.Error("KeyID should remain the same after regeneration")
	}

	if newApiKey.KeyName != apiKey.KeyName {
		t.Error("KeyName should remain the same after regeneration")
	}
}

// TestMaskKeySecret 测试密钥脱敏
func TestMaskKeySecret(t *testing.T) {
	service := &apiKeyManagementService{}

	testCases := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "Short key",
			key:      "sk_abc",
			expected: "sk****",
		},
		{
			name:     "Normal key",
			key:      "sk_1234567890abcdef",
			expected: "sk_****cdef",
		},
		{
			name:     "Long key",
			key:      "sk_" + strings.Repeat("a", 50),
			expected: "sk_****" + strings.Repeat("a", 4),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			masked := service.maskKeySecret(tc.key)
			if masked != tc.expected {
				t.Errorf("maskKeySecret(%q) = %q, want %q", tc.key, masked, tc.expected)
			}

			// 验证脱敏后的密钥不包含完整的原始密钥
			if strings.Contains(masked, strings.Replace(tc.key, "****", "", -1)) {
				t.Error("Masked key should not contain full original key")
			}
		})
	}
}

// TestGenerateKeySecret 测试密钥生成
func TestGenerateKeySecret(t *testing.T) {
	service := &apiKeyManagementService{}

	prefixes := []entity.APIKeyPrefix{
		entity.APIKeyPrefixSecret,
		entity.APIKeyPrefixPublic,
		entity.APIKeyPrefixTest,
	}

	for _, prefix := range prefixes {
		t.Run(string(prefix), func(t *testing.T) {
			key1 := service.generateKeySecret(prefix)
			key2 := service.generateKeySecret(prefix)

			// 生成的密钥应该不同
			if key1 == key2 {
				t.Error("Generated keys should be different")
			}

			// 密钥应该有正确的前缀
			if !strings.HasPrefix(key1, string(prefix)+"_") {
				t.Errorf("Key should have prefix '%s_', got: %s", prefix, key1)
			}

			// 密钥长度应该合理
			if len(key1) < 10 {
				t.Errorf("Key length should be at least 10, got: %d", len(key1))
			}
		})
	}
}
