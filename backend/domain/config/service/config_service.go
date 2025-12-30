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
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/config/entity"
	"github.com/coze-dev/coze-studio/backend/domain/config/repository"
	"github.com/google/uuid"
)

// ConfigService 配置服务
type ConfigService struct {
	configRepo     repository.ConfigRepository
	historyRepo    repository.ConfigHistoryRepository
	dynamicClient  DynamicConfigClient // 可选：etcd等动态配置客户端
}

// DynamicConfigClient 动态配置客户端接口
type DynamicConfigClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
}

// NewConfigService 创建配置服务
func NewConfigService(
	configRepo repository.ConfigRepository,
	historyRepo repository.ConfigHistoryRepository,
	dynamicClient DynamicConfigClient,
) *ConfigService {
	return &ConfigService{
		configRepo:    configRepo,
		historyRepo:   historyRepo,
		dynamicClient: dynamicClient,
	}
}

// CreateConfig 创建配置
func (s *ConfigService) CreateConfig(
	ctx context.Context,
	tenantID,
	key,
	value,
	configType,
	description string,
	updatedBy string,
) (*entity.ConfigItem, error) {
	// 检查key是否已存在
	existing, err := s.configRepo.GetByKey(ctx, tenantID, key)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("config key already exists: %s", key)
	}

	config := &entity.ConfigItem{
		ConfigID:    uuid.New().String(),
		TenantID:    tenantID,
		ConfigKey:   key,
		ConfigValue: value,
		ConfigType:  configType,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UpdatedBy:   updatedBy,
	}

	if err := s.configRepo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("create config failed: %w", err)
	}

	// 同步到配置中心
	if s.dynamicClient != nil {
		dyncKey := s.buildDynamicKey(tenantID, key)
		if err := s.dynamicClient.Set(ctx, dyncKey, value); err != nil {
			// 记录错误但不中断流程
			fmt.Printf("sync to config center failed: %v\n", err)
		}
	}

	return config, nil
}

// GetConfig 获取配置
func (s *ConfigService) GetConfig(ctx context.Context, tenantID, key string) (*entity.ConfigItem, error) {
	// 优先从配置中心获取（如果启用）
	if s.dynamicClient != nil {
		dyncKey := s.buildDynamicKey(tenantID, key)
		value, err := s.dynamicClient.Get(ctx, dyncKey)
		if err == nil {
			return &entity.ConfigItem{
				TenantID:    tenantID,
				ConfigKey:   key,
				ConfigValue: value,
			}, nil
		}
		// 如果配置中心获取失败，回退到数据库
	}

	// 从数据库获取
	return s.configRepo.GetByKey(ctx, tenantID, key)
}

// UpdateConfig 更新配置
func (s *ConfigService) UpdateConfig(
	ctx context.Context,
	tenantID,
	key,
	newValue,
	changeReason,
	updatedBy string,
) error {
	// 获取现有配置
	existing, err := s.configRepo.GetByKey(ctx, tenantID, key)
	if err != nil {
		return fmt.Errorf("get existing config failed: %w", err)
	}

	oldValue := existing.ConfigValue

	// 记录历史
	if err := s.recordHistory(ctx, existing, oldValue, newValue, changeReason, updatedBy); err != nil {
		return fmt.Errorf("record history failed: %w", err)
	}

	// 更新配置
	existing.ConfigValue = newValue
	existing.UpdatedAt = time.Now()
	existing.UpdatedBy = updatedBy

	if err := s.configRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("update config failed: %w", err)
	}

	// 同步到配置中心
	if s.dynamicClient != nil {
		dyncKey := s.buildDynamicKey(tenantID, key)
		if err := s.dynamicClient.Set(ctx, dyncKey, newValue); err != nil {
			fmt.Printf("sync to config center failed: %v\n", err)
		}
	}

	return nil
}

// DeleteConfig 删除配置
func (s *ConfigService) DeleteConfig(ctx context.Context, tenantID, key string) error {
	// 获取配置
	existing, err := s.configRepo.GetByKey(ctx, tenantID, key)
	if err != nil {
		return fmt.Errorf("get config failed: %w", err)
	}

	// 从数据库删除
	if err := s.configRepo.Delete(ctx, existing.ConfigID); err != nil {
		return fmt.Errorf("delete config failed: %w", err)
	}

	// 从配置中心删除
	if s.dynamicClient != nil {
		dyncKey := s.buildDynamicKey(tenantID, key)
		if err := s.dynamicClient.Delete(ctx, dyncKey); err != nil {
			fmt.Printf("delete from config center failed: %v\n", err)
		}
	}

	return nil
}

// ListConfigs 列出租户所有配置
func (s *ConfigService) ListConfigs(ctx context.Context, tenantID string) ([]*entity.ConfigItem, error) {
	return s.configRepo.ListByTenant(ctx, tenantID)
}

// GetConfigHistory 获取配置变更历史
func (s *ConfigService) GetConfigHistory(
	ctx context.Context,
	configID string,
	limit,
	offset int,
) ([]*entity.ConfigHistory, error) {
	return s.historyRepo.ListByConfigID(ctx, configID, limit, offset)
}

// RollbackConfig 回滚配置到指定版本
func (s *ConfigService) RollbackConfig(
	ctx context.Context,
	tenantID,
	key string,
	version int,
	updatedBy string,
) error {
	// 获取配置历史
	histories, err := s.historyRepo.ListByConfigID(ctx, key, 100, 0)
	if err != nil {
		return fmt.Errorf("get history failed: %w", err)
	}

	// 查找目标版本
	var targetHistory *entity.ConfigHistory
	for _, h := range histories {
		if h.VersionNumber == version {
			targetHistory = h
			break
		}
	}

	if targetHistory == nil {
		return fmt.Errorf("version not found: %d", version)
	}

	// 更新配置为旧值
	return s.UpdateConfig(
		ctx,
		tenantID,
		key,
		targetHistory.OldValue,
		fmt.Sprintf("rollback to version %d", version),
		updatedBy,
	)
}

// GetConfigAsJSON 获取JSON类型配置并解析
func (s *ConfigService) GetConfigAsJSON(ctx context.Context, tenantID, key string, target interface{}) error {
	config, err := s.GetConfig(ctx, tenantID, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(config.ConfigValue), target)
}

// recordHistory 记录配置变更历史
func (s *ConfigService) recordHistory(
	ctx context.Context,
	config *entity.ConfigItem,
	oldValue,
	newValue,
	reason,
	updatedBy string,
) error {
	// 获取最新版本号
	latestVersion, err := s.historyRepo.GetLatestVersion(ctx, config.ConfigID)
	if err != nil {
		latestVersion = 0
	}

	history := &entity.ConfigHistory{
		HistoryID:     uuid.New().String(),
		ConfigID:      config.ConfigID,
		TenantID:      config.TenantID,
		ConfigKey:     config.ConfigKey,
		OldValue:      oldValue,
		NewValue:      newValue,
		ChangeReason:  reason,
		ChangedBy:     updatedBy,
		ChangedAt:     time.Now(),
		VersionNumber: latestVersion + 1,
	}

	return s.historyRepo.Create(ctx, history)
}

// buildDynamicKey 构建动态配置key
func (s *ConfigService) buildDynamicKey(tenantID, key string) string {
	return fmt.Sprintf("/tenant/%s/config/%s", tenantID, key)
}
