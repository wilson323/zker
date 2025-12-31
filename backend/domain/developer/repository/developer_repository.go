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

package repository

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
)

// DeveloperRepository 开发者仓储接口
type DeveloperRepository interface {
	// Create 创建开发者
	Create(ctx context.Context, developer *entity.Developer) error

	// GetByID 根据ID获取开发者
	GetByID(ctx context.Context, developerID string) (*entity.Developer, error)

	// GetByTenantIDAndUserID 根据租户ID和用户ID获取开发者
	GetByTenantIDAndUserID(ctx context.Context, tenantID, userID string) (*entity.Developer, error)

	// GetByEmail 根据邮箱获取开发者
	GetByEmail(ctx context.Context, email string) (*entity.Developer, error)

	// Update 更新开发者
	Update(ctx context.Context, developer *entity.Developer) error

	// Delete 软删除开发者
	Delete(ctx context.Context, developerID string) error

	// List 分页查询开发者列表
	List(ctx context.Context, filter *DeveloperFilter) ([]*entity.Developer, int64, error)

	// UpdateStatus 更新开发者状态
	UpdateStatus(ctx context.Context, developerID string, status entity.DeveloperStatus) error

	// Exists 检查开发者是否存在
	Exists(ctx context.Context, tenantID, userID string) (bool, error)
}

// DeveloperFilter 开发者查询过滤器
type DeveloperFilter struct {
	TenantID string
	UserID   string
	Status   entity.DeveloperStatus
	PageToken string
	PageSize  int
}

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	// Create 创建项目
	Create(ctx context.Context, project *entity.Project) error

	// GetByID 根据ID获取项目
	GetByID(ctx context.Context, projectID string) (*entity.Project, error)

	// GetByTenantIDAndName 根据租户ID和名称获取项目
	GetByTenantIDAndName(ctx context.Context, tenantID, name string) (*entity.Project, error)

	// GetByDeveloperID 根据开发者ID获取项目列表
	GetByDeveloperID(ctx context.Context, developerID string) ([]*entity.Project, error)

	// GetByTenantID 根据租户ID获取项目列表
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Project, error)

	// Update 更新项目
	Update(ctx context.Context, project *entity.Project) error

	// Delete 软删除项目
	Delete(ctx context.Context, projectID string) error

	// List 分页查询项目列表
	List(ctx context.Context, filter *ProjectFilter) ([]*entity.Project, int64, error)

	// UpdateStatus 更新项目状态
	UpdateStatus(ctx context.Context, projectID string, status entity.ProjectStatus) error

	// UpdateConfig 更新项目配置
	UpdateConfig(ctx context.Context, projectID string, config string) error

	// CountByDeveloperID 统计开发者的项目数量
	CountByDeveloperID(ctx context.Context, developerID string) (int64, error)

	// CountByTenantID 统计租户的项目数量
	CountByTenantID(ctx context.Context, tenantID string) (int64, error)
}

// ProjectFilter 项目查询过滤器
type ProjectFilter struct {
	TenantID     string
	DeveloperID  string
	ProjectType  entity.ProjectType
	ProjectStatus entity.ProjectStatus
	PageToken    string
	PageSize     int
}

// APIKeyRepository API密钥仓储接口
type APIKeyRepository interface {
	// Create 创建API密钥
	Create(ctx context.Context, apiKey *entity.APIKey) error

	// GetByID 根据ID获取API密钥
	GetByID(ctx context.Context, keyID string) (*entity.APIKey, error)

	// GetByProjectID 根据项目ID获取API密钥列表
	GetByProjectID(ctx context.Context, projectID string) ([]*entity.APIKey, error)

	// GetByKeySecret 根据密钥获取API密钥（用于验证）
	GetByKeySecret(ctx context.Context, keySecret string) (*entity.APIKey, error)

	// Update 更新API密钥
	Update(ctx context.Context, apiKey *entity.APIKey) error

	// Delete 软删除API密钥
	Delete(ctx context.Context, keyID string) error

	// List 分页查询API密钥列表
	List(ctx context.Context, filter *APIKeyFilter) ([]*entity.APIKey, int64, error)

	// Revoke 撤销API密钥
	Revoke(ctx context.Context, keyID string) error

	// UpdateLastUsedAt 更新最后使用时间
	UpdateLastUsedAt(ctx context.Context, keyID string) error

	// Expire 过期API密钥
	Expire(ctx context.Context, keyID string) error

	// CountByProjectID 统计项目的API密钥数量
	CountByProjectID(ctx context.Context, projectID string) (int64, error)

	// GetExpiringKeys 获取即将过期的密钥（7天内）
	// GetAllActive 获取所有激活的API密钥（用于密钥验证）
	GetAllActive(ctx context.Context) ([]*entity.APIKey, error)

	GetExpiringKeys(ctx context.Context, tenantID string) ([]*entity.APIKey, error)
}

// APIKeyFilter API密钥查询过滤器
type APIKeyFilter struct {
	TenantID   string
	ProjectID  string
	KeyPrefix  entity.APIKeyPrefix
	Status     entity.APIKeyStatus
	PageToken  string
	PageSize   int
}

// WebhookRepository Webhook仓储接口
type WebhookRepository interface {
	// Create 创建Webhook
	Create(ctx context.Context, webhook *entity.Webhook) error

	// GetByID 根据ID获取Webhook
	GetByID(ctx context.Context, webhookID string) (*entity.Webhook, error)

	// GetByProjectID 根据项目ID获取Webhook列表
	GetByProjectID(ctx context.Context, projectID string) ([]*entity.Webhook, error)

	// GetByTenantID 根据租户ID获取Webhook列表
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Webhook, error)

	// Update 更新Webhook
	Update(ctx context.Context, webhook *entity.Webhook) error

	// Delete 软删除Webhook
	Delete(ctx context.Context, webhookID string) error

	// List 分页查询Webhook列表
	List(ctx context.Context, filter *WebhookFilter) ([]*entity.Webhook, int64, error)

	// UpdateStatus 更新Webhook状态
	UpdateStatus(ctx context.Context, webhookID string, status entity.WebhookStatus) error

	// UpdateStats 更新统计信息
	UpdateStats(ctx context.Context, webhookID string, success bool) error

	// UpdateLastTriggerAt 更新最后触发时间
	UpdateLastTriggerAt(ctx context.Context, webhookID string) error

	// CountByProjectID 统计项目的Webhook数量
	CountByProjectID(ctx context.Context, projectID string) (int64, error)
}

// WebhookFilter Webhook查询过滤器
type WebhookFilter struct {
	TenantID string
	ProjectID string
	Status   entity.WebhookStatus
	PageToken string
	PageSize  int
}

// SDKRepository SDK仓储接口
type SDKRepository interface {
	// Create 创建SDK
	Create(ctx context.Context, sdk *entity.SDK) error

	// GetByID 根据ID获取SDK
	GetByID(ctx context.Context, sdkID string) (*entity.SDK, error)

	// GetByProjectID 根据项目ID获取SDK列表
	GetByProjectID(ctx context.Context, projectID string) ([]*entity.SDK, error)

	// GetByProjectIDAndLanguage 根据项目ID和语言获取SDK
	GetByProjectIDAndLanguage(ctx context.Context, projectID string, language entity.SDKLanguage) (*entity.SDK, error)

	// Update 更新SDK
	Update(ctx context.Context, sdk *entity.SDK) error

	// Delete 软删除SDK
	Delete(ctx context.Context, sdkID string) error

	// List 分页查询SDK列表
	List(ctx context.Context, filter *SDKFilter) ([]*entity.SDK, int64, error)

	// UpdateStatus 更新SDK状态
	UpdateStatus(ctx context.Context, sdkID string, status entity.SDKStatus) error

	// IncrementDownloadCount 增加下载次数
	IncrementDownloadCount(ctx context.Context, sdkID string) error

	// GetLatestVersion 获取最新版本的SDK
	GetLatestVersion(ctx context.Context, projectID string, language entity.SDKLanguage) (*entity.SDK, error)
}

// SDKFilter SDK查询过滤器
type SDKFilter struct {
	TenantID  string
	ProjectID string
	Language  entity.SDKLanguage
	Status    entity.SDKStatus
	PageToken string
	PageSize  int
}

// WebhookLogRepository Webhook日志仓储接口
type WebhookLogRepository interface {
	// Create 创建Webhook日志
	Create(ctx context.Context, log *entity.WebhookLog) error

	// GetByWebhookID 根据Webhook ID获取日志列表
	GetByWebhookID(ctx context.Context, webhookID string, limit, offset int) ([]*entity.WebhookLog, error)

	// DeleteOldLogs 删除旧日志（保留N天）
	DeleteOldLogs(ctx context.Context, days int) (int64, error)

	// GetStats 获取统计信息
	GetStats(ctx context.Context, webhookID string, startTime, endTime int64) (*WebhookStats, error)
}

// WebhookStats Webhook统计信息
type WebhookStats struct {
	TotalCalls   int64   `json:"total_calls"`
	SuccessCalls int64   `json:"success_calls"`
	FailureCalls int64   `json:"failure_calls"`
	SuccessRate  float64 `json:"success_rate"`
	AvgDuration  float64 `json:"avg_duration"`
}

// WebhookDLQRepository Webhook死信队列仓储接口
type WebhookDLQRepository interface {
	// Create 创建死信队列条目
	Create(ctx context.Context, entry *entity.WebhookDLQEntry) error

	// GetByID 根据ID获取死信队列条目
	GetByID(ctx context.Context, dlqID string) (*entity.WebhookDLQEntry, error)

	// GetByWebhookID 根据Webhook ID获取死信队列条目列表
	GetByWebhookID(ctx context.Context, webhookID string, limit, offset int) ([]*entity.WebhookDLQEntry, error)

	// GetExpired 获取过期的死信队列条目（用于清理）
	GetExpired(ctx context.Context) ([]*entity.WebhookDLQEntry, error)

	// GetRetryable 获取可重试的条目
	GetRetryable(ctx context.Context, maxRetries int) ([]*entity.WebhookDLQEntry, error)

	// UpdateRetryCount 更新重试次数
	UpdateRetryCount(ctx context.Context, dlqID string, retryCount int) error

	// Delete 根据ID删除死信队列条目
	Delete(ctx context.Context, dlqID string) error

	// DeleteExpired 删除过期的死信队列条目
	DeleteExpired(ctx context.Context) (int64, error)

	// CountByWebhookID 统计Webhook的死信队列条目数量
	CountByWebhookID(ctx context.Context, webhookID string) (int64, error)
}
