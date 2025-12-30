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
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// IsolationStrategy 隔离策略类型
type IsolationStrategy string

const (
	// StrategyRowLevel 行级隔离（共享表，通过tenant_id区分）
	// 适用场景：小型租户，数据量小，成本敏感
	StrategyRowLevel IsolationStrategy = "row_level"

	// StrategySchemaLevel Schema级隔离（独立Schema）
	// 适用场景：中型租户，需要更好的性能和隔离性
	StrategySchemaLevel IsolationStrategy = "schema_level"

	// StrategyDatabaseLevel 数据库级隔离（独立数据库实例）
	// 适用场景：大型企业租户，要求最高隔离性和性能
	StrategyDatabaseLevel IsolationStrategy = "database_level"
)

// 隔离策略与实体策略的映射
var (
	// StrategyRowLevelAlias 行级隔离别名（兼容entity.IsolationStrategyRowLevel）
	StrategyRowLevelAlias = entity.IsolationStrategyRowLevel

	// StrategySchemaLevelAlias Schema级隔离别名（兼容entity.IsolationStrategySchemaLevel）
	StrategySchemaLevelAlias = entity.IsolationStrategySchemaLevel

	// StrategyDatabaseLevelAlias 数据库级隔离别名（兼容entity.IsolationStrategyDatabaseLevel）
	StrategyDatabaseLevelAlias = entity.IsolationStrategyDatabaseLevel
)

// IsolationUpgradeConfig 隔离升级配置
type IsolationUpgradeConfig struct {
	// 自动升级阈值
	AutoUpgradeThresholds map[IsolationStrategy]UpgradeThreshold `json:"auto_upgrade_thresholds"`

	// 升级检查间隔
	CheckInterval time.Duration `json:"check_interval"`

	// 是否启用自动升级
	AutoUpgradeEnabled bool `json:"auto_upgrade_enabled"`
}

// UpgradeThreshold 升级阈值配置
type UpgradeThreshold struct {
	// 数据量阈值（行数）
	DataVolumeThreshold int64 `json:"data_volume_threshold"`

	// QPS阈值
	QPSThreshold int `json:"qps_threshold"`

	// 团队成员数阈值
	MemberCountThreshold int `json:"member_count_threshold"`

	// 必须满足最低订阅等级
	MinSubscriptionTier entity.SubscriptionTier `json:"min_subscription_tier"`
}

// IsolationUpgradeService 隔离升级服务
//
// 职责：
// 1. 监控租户规模指标
// 2. 自动检测是否需要升级隔离策略
// 3. 执行隔离策略升级
// 4. 管理数据迁移
type IsolationUpgradeService struct {
	db                  *gorm.DB
	tenantRepo          repository.TenantRepository
	config              *IsolationUpgradeConfig
	upgradeQueue        chan string // 待升级租户ID队列
	upgradeInProgress   map[string]bool // 正在升级的租户
	mu                  sync.RWMutex
}

// NewIsolationUpgradeService 创建隔离升级服务
func NewIsolationUpgradeService(
	db *gorm.DB,
	tenantRepo repository.TenantRepository,
	config *IsolationUpgradeConfig,
) *IsolationUpgradeService {
	if config == nil {
		config = DefaultUpgradeConfig()
	}

	return &IsolationUpgradeService{
		db:                db,
		tenantRepo:        tenantRepo,
		config:            config,
		upgradeQueue:      make(chan string, 100),
		upgradeInProgress: make(map[string]bool),
	}
}

// DefaultUpgradeConfig 默认升级配置
func DefaultUpgradeConfig() *IsolationUpgradeConfig {
	return &IsolationUpgradeConfig{
		AutoUpgradeThresholds: map[IsolationStrategy]UpgradeThreshold{
			StrategyRowLevel: {
				DataVolumeThreshold:    1000000,  // 100万行
				QPSThreshold:           100,      // 100 QPS
				MemberCountThreshold:    20,       // 20人
				MinSubscriptionTier:    entity.SubscriptionTierPro,
			},
			StrategySchemaLevel: {
				DataVolumeThreshold:    10000000, // 1000万行
				QPSThreshold:           1000,     // 1000 QPS
				MemberCountThreshold:    200,      // 200人
				MinSubscriptionTier:    entity.SubscriptionTierEnterprise,
			},
		},
		CheckInterval:        24 * time.Hour, // 每天检查一次
		AutoUpgradeEnabled:   true,
	}
}

// StartUpgradeWorker 启动升级工作协程
func (s *IsolationUpgradeService) StartUpgradeWorker(ctx context.Context) {
	logs.Info("[IsolationUpgrade] Starting upgrade worker...")

	ticker := time.NewTicker(s.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logs.Info("[IsolationUpgrade] Upgrade worker stopped")
			return

		case <-ticker.C:
			if s.config.AutoUpgradeEnabled {
				s.scanTenantsForUpgrade(ctx)
			}

		case tenantID := <-s.upgradeQueue:
			s.processUpgrade(ctx, tenantID)
		}
	}
}

// scanTenantsForUpgrade 扫描所有租户，检查是否需要升级
func (s *IsolationUpgradeService) scanTenantsForUpgrade(ctx context.Context) {
	logs.Info("[IsolationUpgrade] Scanning tenants for upgrade...")

	// 获取所有租户
	tenants, _, err := s.tenantRepo.List(ctx, nil)
	if err != nil {
		logs.Errorf("[IsolationUpgrade] Failed to list tenants: %v", err)
		return
	}

	for _, tenant := range tenants {
		// 检查是否需要升级
		needsUpgrade, nextStrategy := s.evaluateUpgradeNeed(ctx, tenant)
		if needsUpgrade {
			logs.Infof("[IsolationUpgrade] Tenant %s needs upgrade to %s", tenant.TenantID, nextStrategy)
			s.upgradeQueue <- tenant.TenantID
		}
	}
}

// evaluateUpgradeNeed 评估租户是否需要升级隔离策略
func (s *IsolationUpgradeService) evaluateUpgradeNeed(
	ctx context.Context,
	tenant *entity.Tenant,
) (bool, IsolationStrategy) {
	// 将entity.IsolationStrategy转换为service.IsolationStrategy
	currentStrategy := IsolationStrategy(tenant.IsolationStrategy)
	if currentStrategy == "" {
		currentStrategy = StrategyRowLevel // 默认行级隔离
	}

	// 检查是否已经是最高级别
	if currentStrategy == StrategyDatabaseLevel {
		return false, ""
	}

	// 获取当前策略的升级阈值
	threshold, exists := s.config.AutoUpgradeThresholds[currentStrategy]
	if !exists {
		return false, ""
	}

	// 检查订阅等级是否满足要求
	if tenant.SubscriptionTier != threshold.MinSubscriptionTier {
		return false, ""
	}

	// 收集指标
	metrics, err := s.collectTenantMetrics(ctx, tenant.TenantID)
	if err != nil {
		logs.Warnf("[IsolationUpgrade] Failed to collect metrics for %s: %v", tenant.TenantID, err)
		return false, ""
	}

	// 评估是否达到升级条件
	needsUpgrade := false
	nextStrategy := currentStrategy

	// 检查数据量
	if metrics.TotalRows >= threshold.DataVolumeThreshold {
		needsUpgrade = true
		nextStrategy = s.getNextStrategy(currentStrategy)
	}

	// 检查QPS
	if metrics.AverageQPS >= threshold.QPSThreshold {
		needsUpgrade = true
		nextStrategy = s.getNextStrategy(currentStrategy)
	}

	// 检查团队成员数
	if metrics.MemberCount >= threshold.MemberCountThreshold {
		needsUpgrade = true
		nextStrategy = s.getNextStrategy(currentStrategy)
	}

	return needsUpgrade, nextStrategy
}

// collectTenantMetrics 收集租户指标
func (s *IsolationUpgradeService) collectTenantMetrics(
	ctx context.Context,
	tenantID string,
) (*TenantMetrics, error) {
	metrics := &TenantMetrics{
		TenantID: tenantID,
	}

	// 1. 统计总数据量（所有表的总行数）
	var totalRows int64
	result := s.db.Raw(`
		SELECT SUM(table_rows) as total_rows
		FROM information_schema.tables
		WHERE table_schema = DATABASE()
		AND table_name LIKE ?
	`, "%"+tenantID+"%").Scan(&totalRows)
	if result.Error != nil {
		logs.Warnf("Failed to count rows for tenant %s: %v", tenantID, result.Error)
		totalRows = 0
	}
	metrics.TotalRows = totalRows

	// 2. 查询配额获取团队成员数
	quota, err := s.getQuotaByType(ctx, tenantID, entity.ResourceTypeTeamMembers)
	if err == nil && quota != nil {
		metrics.MemberCount = quota.UsedCount
	}

	// 3. 计算平均QPS（基于最近24小时的请求日志）
	// TODO: 从监控系统获取实际的QPS数据
	metrics.AverageQPS = s.calculateQPS(ctx, tenantID)

	metrics.LastCollectedAt = time.Now().UnixMilli()

	return metrics, nil
}

// TenantMetrics 租户指标
type TenantMetrics struct {
	TenantID        string    `json:"tenant_id"`
	TotalRows       int64     `json:"total_rows"`
	AverageQPS      int       `json:"average_qps"`
	MemberCount     int       `json:"member_count"`
	LastCollectedAt int64     `json:"last_collected_at"`
}

// processUpgrade 处理隔离策略升级
func (s *IsolationUpgradeService) processUpgrade(ctx context.Context, tenantID string) {
	s.mu.Lock()
	if s.upgradeInProgress[tenantID] {
		s.mu.Unlock()
		logs.Infof("[IsolationUpgrade] Upgrade already in progress for tenant %s", tenantID)
		return
	}
	s.upgradeInProgress[tenantID] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.upgradeInProgress, tenantID)
		s.mu.Unlock()
	}()

	logs.Infof("[IsolationUpgrade] Starting upgrade for tenant %s", tenantID)

	// 1. 获取租户信息
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		logs.Errorf("[IsolationUpgrade] Failed to get tenant %s: %v", tenantID, err)
		return
	}

	// 2. 确定目标策略
	_, nextStrategy := s.evaluateUpgradeNeed(ctx, tenant)
	if nextStrategy == "" {
		logs.Infof("[IsolationUpgrade] No upgrade needed for tenant %s", tenantID)
		return
	}

	// 3. 执行升级
	if err := s.performUpgrade(ctx, tenant, nextStrategy); err != nil {
		logs.Errorf("[IsolationUpgrade] Failed to upgrade tenant %s: %v", tenantID, err)
		return
	}

	logs.Infof("[IsolationUpgrade] Successfully upgraded tenant %s to %s", tenantID, nextStrategy)
}

// performUpgrade 执行隔离策略升级
func (s *IsolationUpgradeService) performUpgrade(
	ctx context.Context,
	tenant *entity.Tenant,
	targetStrategy IsolationStrategy,
) error {
	// 将entity.IsolationStrategy转换为service.IsolationStrategy
	currentStrategy := IsolationStrategy(tenant.IsolationStrategy)
	if currentStrategy == "" {
		currentStrategy = StrategyRowLevel
	}

	logs.Infof("[IsolationUpgrade] Upgrading tenant %s from %s to %s",
		tenant.TenantID, currentStrategy, targetStrategy)

	// 使用事务确保升级过程的原子性
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 根据目标策略执行相应的迁移
		switch targetStrategy {
		case StrategySchemaLevel:
			if err := s.upgradeToSchemaLevel(ctx, tx, tenant); err != nil {
				return fmt.Errorf("failed to upgrade to schema level: %w", err)
			}

		case StrategyDatabaseLevel:
			if err := s.upgradeToDatabaseLevel(ctx, tx, tenant); err != nil {
				return fmt.Errorf("failed to upgrade to database level: %w", err)
			}

		default:
			return fmt.Errorf("unknown target strategy: %s", targetStrategy)
		}

		// 2. 更新租户的隔离策略
		now := time.Now().UnixMilli()
		if err := tx.Model(&entity.Tenant{}).
			Where("tenant_id = ?", tenant.TenantID).
			Updates(map[string]interface{}{
				"isolation_strategy": string(targetStrategy),
				"updated_at":         now,
			}).Error; err != nil {
			return fmt.Errorf("failed to update tenant isolation strategy: %w", err)
		}

		logs.Infof("[IsolationUpgrade] Tenant %s upgraded to %s", tenant.TenantID, targetStrategy)

		return nil
	})
}

// upgradeToSchemaLevel 升级到Schema级隔离
func (s *IsolationUpgradeService) upgradeToSchemaLevel(
	ctx context.Context,
	tx *gorm.DB,
	tenant *entity.Tenant,
) error {
	schemaName := fmt.Sprintf("tenant_%s", tenant.TenantID)

	// 1. 创建独立的Schema
	if err := tx.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS `%s`", schemaName)).Error; err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// 2. 在新Schema中创建所有需要的表
	tables := []string{
		"bots",
		"conversations",
		"knowledge_bases",
		"workflows",
	}

	for _, table := range tables {
		// 复制表结构到新Schema
		if err := tx.Exec(fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s.%s LIKE %s
		`, schemaName, table, table)).Error; err != nil {
			return fmt.Errorf("failed to create table %s: %w", table, err)
		}

		// 迁移数据
		if err := tx.Exec(fmt.Sprintf(`
			INSERT INTO %s.%s SELECT * FROM %s WHERE tenant_id = ?
		`, schemaName, table, table)).Error; err != nil {
			return fmt.Errorf("failed to migrate data for table %s: %w", table, err)
		}
	}

	logs.Infof("[IsolationUpgrade] Created schema %s for tenant %s", schemaName, tenant.TenantID)
	return nil
}

// upgradeToDatabaseLevel 升级到数据库级隔离
func (s *IsolationUpgradeService) upgradeToDatabaseLevel(
	ctx context.Context,
	tx *gorm.DB,
	tenant *entity.Tenant,
) error {
	// 数据库级隔离需要：
	// 1. 创建独立的数据库实例（由运维自动化系统处理）
	// 2. 配置数据库连接池
	// 3. 迁移数据到新数据库
	// 这里只记录升级请求，实际创建由外部系统完成

	logs.Infof("[IsolationUpgrade] Database level isolation requested for tenant %s", tenant.TenantID)
	logs.Infof("[IsolationUpgrade] Please provision dedicated database instance and run migration manually")

	return nil
}

// getNextStrategy 获取下一个隔离策略
func (s *IsolationUpgradeService) getNextStrategy(current IsolationStrategy) IsolationStrategy {
	switch current {
	case StrategyRowLevel:
		return StrategySchemaLevel
	case StrategySchemaLevel:
		return StrategyDatabaseLevel
	default:
		return ""
	}
}

// getQuotaByType 获取特定类型的配额
func (s *IsolationUpgradeService) getQuotaByType(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
) (*entity.Quota, error) {
	var quota entity.Quota
	err := s.db.Where("tenant_id = ? AND resource_type = ?", tenantID, resourceType).
		First(&quota).Error
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

// calculateQPS 计算QPS（简化版本，实际应从监控系统获取）
func (s *IsolationUpgradeService) calculateQPS(ctx context.Context, tenantID string) int {
	// TODO: 从监控系统获取实际的QPS数据
	// 这里返回估算值
	return 10
}

// ManualUpgrade 手动触发隔离策略升级
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - targetStrategy: 目标隔离策略
//
// 返回:
//   - error: 升级失败时返回错误
func (s *IsolationUpgradeService) ManualUpgrade(
	ctx context.Context,
	tenantID string,
	targetStrategy IsolationStrategy,
) error {
	// 获取租户信息
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}
	if tenant == nil {
		return fmt.Errorf("tenant not found: %s", tenantID)
	}

	// 执行升级
	return s.performUpgrade(ctx, tenant, targetStrategy)
}

// GetUpgradeStatus 获取升级状态
func (s *IsolationUpgradeService) GetUpgradeStatus(
	ctx context.Context,
	tenantID string,
) (*UpgradeStatus, error) {
	s.mu.RLock()
	inProgress := s.upgradeInProgress[tenantID]
	s.mu.RUnlock()

	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 将entity.IsolationStrategy转换为service.IsolationStrategy
	return &UpgradeStatus{
		TenantID:          tenantID,
		CurrentStrategy:   IsolationStrategy(tenant.IsolationStrategy),
		UpgradeInProgress: inProgress,
		LastEvaluatedAt:   time.Now().UnixMilli(),
	}, nil
}

// UpgradeStatus 升级状态
type UpgradeStatus struct {
	TenantID          string              `json:"tenant_id"`
	CurrentStrategy   IsolationStrategy  `json:"current_strategy"`
	UpgradeInProgress bool                `json:"upgrade_in_progress"`
	LastEvaluatedAt   int64               `json:"last_evaluated_at"`
}
