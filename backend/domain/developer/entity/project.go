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

package entity

import (
	"time"
)

// ProjectType 项目类型
type ProjectType string

const (
	ProjectTypeBot         ProjectType = "bot"         // Bot项目
	ProjectTypeWorkflow    ProjectType = "workflow"    // 工作流项目
	ProjectTypeIntegration ProjectType = "integration" // 集成项目
)

// ProjectStatus 项目状态
type ProjectStatus string

const (
	ProjectStatusDevelopment ProjectStatus = "development" // 开发中
	ProjectStatusProduction  ProjectStatus = "production"  // 生产环境
	ProjectStatusArchived    ProjectStatus = "archived"    // 已归档
)

// Project 项目实体
type Project struct {
	ProjectID    string        `json:"project_id" gorm:"primaryKey;type:varchar(36);comment:项目ID"`
	TenantID     string        `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id;comment:租户ID"`
	DeveloperID  string        `json:"developer_id" gorm:"type:varchar(36);not null;index:idx_developer_id;comment:开发者ID"`
	ProjectName  string        `json:"project_name" gorm:"type:varchar(100);not null;comment:项目名称"`
	ProjectType  ProjectType   `json:"project_type" gorm:"type:enum('bot','workflow','integration');not null;index:idx_type;comment:项目类型"`
	Description  string        `json:"description" gorm:"type:text;comment:项目描述"`
	Status       ProjectStatus `json:"status" gorm:"type:enum('development','production','archived');default:'development';not null;index:idx_status;comment:状态"`
	Config       string        `json:"config,omitempty" gorm:"type:json;comment:项目配置"` // JSON格式的配置
	CreatedAt    int64         `json:"created_at" gorm:"not null;default:0;comment:创建时间"`
	UpdatedAt    int64         `json:"updated_at" gorm:"not null;default:0;comment:更新时间"`
	DeletedAt    *int64        `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"`

	// 关联
	Developer Developer `json:"developer,omitempty" gorm:"foreignKey:DeveloperID;references:DeveloperID"`
	APIKeys   []APIKey  `json:"api_keys,omitempty" gorm:"foreignKey:ProjectID;references:ProjectID"`
	Webhooks  []Webhook `json:"webhooks,omitempty" gorm:"foreignKey:ProjectID;references:ProjectID"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "developer_projects"
}

// IsActive 是否激活（非归档且未删除）
func (p *Project) IsActive() bool {
	return p.Status != ProjectStatusArchived && p.DeletedAt == nil
}

// IsProduction 是否为生产环境
func (p *Project) IsProduction() bool {
	return p.Status == ProjectStatusProduction
}

// IsDeleted 是否已删除
func (p *Project) IsDeleted() bool {
	return p.DeletedAt != nil
}

// GetCreatedAtAsTime 获取创建时间
func (p *Project) GetCreatedAtAsTime() time.Time {
	return time.Unix(p.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (p *Project) GetUpdatedAtAsTime() time.Time {
	return time.Unix(p.UpdatedAt/1000, 0)
}

// ProjectConfig 项目配置
type ProjectConfig struct {
	EnableWebhook     bool     `json:"enable_webhook,omitempty"`
	EnableRateLimit   bool     `json:"enable_rate_limit,omitempty"`
	RateLimitPerMin   int      `json:"rate_limit_per_min,omitempty"`
	AllowedOrigins    []string `json:"allowed_origins,omitempty"`
	EnableCaching     bool     `json:"enable_caching,omitempty"`
	CacheTTL          int      `json:"cache_ttl,omitempty"`
	EnableAnalytics   bool     `json:"enable_analytics,omitempty"`
	EnableMonitoring  bool     `json:"enable_monitoring,omitempty"`
	CustomDomain      string   `json:"custom_domain,omitempty"`
	WebhookSecret     string   `json:"webhook_secret,omitempty"`
	RetentionDays     int      `json:"retention_days,omitempty"`
}
