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

import "time"

// TokenUsage Token使用记录
type TokenUsage struct {
	ID              string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID        string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID          string    `json:"user_id" gorm:"type:varchar(36);not null"`
	BotID           string    `json:"bot_id" gorm:"type:varchar(36);not null;index:idx_bot_id"`
	ModelID         string    `json:"model_id" gorm:"type:varchar(100);not null"`
	PromptTokens    int       `json:"prompt_tokens" gorm:"not null"`
	CompletionTokens int      `json:"completion_tokens" gorm:"not null"`
	TotalTokens     int       `json:"total_tokens" gorm:"not null"`
	CostUSD         float64   `json:"cost_usd" gorm:"type:decimal(10,4);not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"not null"`
}

// TokenBudget Token预算
type TokenBudget struct {
	ID            string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string    `json:"tenant_id" gorm:"type:varchar(36);not null;uniqueIndex:uk_tenant"`

	// 预算配置
	MaxTokens                int     `json:"max_tokens" gorm:"not null;default:0;comment:最大Token数"`
	ResetPeriod              string  `json:"reset_period" gorm:"type:varchar(20);default:'monthly';comment:重置周期"`
	AlertThresholdPercent    int     `json:"alert_threshold_percent" gorm:"default:80;comment:告警阈值(百分比)"`
	AutoDegradationModel     string  `json:"auto_degradation_model" gorm:"type:varchar(100);comment:自动降级模型"`

	// 遗留字段（兼容旧版本）
	BudgetLimit   int       `json:"budget_limit" gorm:"not null;default:0"`
	BudgetUsed    int       `json:"budget_used" gorm:"not null;default:0"`
	LastResetAt   time.Time `json:"last_reset_at"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"not null"`
}
