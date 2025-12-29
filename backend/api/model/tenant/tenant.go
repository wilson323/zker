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

package tenant

// ==================== Request DTOs ====================

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	TenantName   string `json:"tenant_name" binding:"required,min=1,max=200"`
	TenantType   string `json:"tenant_type" binding:"required,oneof=individual team enterprise"`
	ContactEmail string `json:"contact_email" binding:"required,email"`
	ContactPhone string `json:"contact_phone"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	TenantName   *string `json:"tenant_name,omitempty" binding:"omitempty,min=1,max=200"`
	ContactEmail *string `json:"contact_email,omitempty" binding:"omitempty,email"`
	ContactPhone *string `json:"contact_phone,omitempty"`
}

// UpgradeSubscriptionRequest 升级订阅请求
type UpgradeSubscriptionRequest struct {
	TargetTier string `json:"target_tier" binding:"required,oneof=free pro enterprise"`
}

// UpdateSubscriptionRequest 更新订阅请求
type UpdateSubscriptionRequest struct {
	BillingCycle *string `json:"billing_cycle,omitempty" binding:"omitempty,oneof=monthly yearly"`
	AutoRenew    *bool   `json:"auto_renew,omitempty"`
}

// CheckQuotaRequest 检查配额请求
type CheckQuotaRequest struct {
	ResourceType string `json:"resource_type" binding:"required,oneof=bots workflows messages storage team_members"`
	RequiredCount int   `json:"required_count" binding:"required,min=1"`
}

// ConsumeQuotaRequest 消费配额请求
type ConsumeQuotaRequest struct {
	ResourceType string `json:"resource_type" binding:"required,oneof=bots workflows messages storage team_members"`
	Count        int    `json:"count" binding:"required,min=1"`
}

// RollbackQuotaRequest 回滚配额请求
type RollbackQuotaRequest struct {
	ResourceType string `json:"resource_type" binding:"required,oneof=bots workflows messages storage team_members"`
	Count        int    `json:"count" binding:"required,min=1"`
}

// RecordUsageRequest 记录使用量请求
type RecordUsageRequest struct {
	ResourceType string                 `json:"resource_type" binding:"required,oneof=bots workflows messages storage team_members"`
	Action       string                 `json:"action" binding:"required,oneof=create update delete query"`
	Quantity     int                    `json:"quantity" binding:"required,min=1"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// GenerateInvoiceRequest 生成账单请求
type GenerateInvoiceRequest struct {
	StartDate string `json:"start_date" binding:"required" example:"2024-01-01"`
	EndDate   string `json:"end_date" binding:"required" example:"2024-01-31"`
}

// ListTenantsRequest 列出租户请求
type ListTenantsRequest struct {
	Status           *string `form:"status" binding:"omitempty,oneof=active suspended deleted"`
	TenantType       *string `form:"tenant_type" binding:"omitempty,oneof=individual team enterprise"`
	SubscriptionTier *string `form:"subscription_tier" binding:"omitempty,oneof=free pro enterprise"`
	PageSize         int     `form:"page_size" binding:"omitempty,min=1,max=100"`
	PageToken        string  `form:"page_token" binding:"omitempty"`
}

// GetUsageSummaryRequest 获取使用量汇总请求
type GetUsageSummaryRequest struct {
	StartDate string `form:"start_date" binding:"required" example:"2024-01-01"`
	EndDate   string `form:"end_date" binding:"required" example:"2024-01-31"`
}

// ==================== Response DTOs ====================

// TenantInfo 租户信息
type TenantInfo struct {
	TenantID         string `json:"tenant_id"`
	TenantName       string `json:"tenant_name"`
	TenantType       string `json:"tenant_type"`
	Status           string `json:"status"`
	SubscriptionTier string `json:"subscription_tier"`
	CreatedAt        int64  `json:"created_at"`
	UpdatedAt        int64  `json:"updated_at"`
}

// TenantDetail 租户详情
type TenantDetail struct {
	TenantID         string        `json:"tenant_id"`
	TenantName       string        `json:"tenant_name"`
	TenantType       string        `json:"tenant_type"`
	Status           string        `json:"status"`
	SubscriptionTier string        `json:"subscription_tier"`
	ContactEmail     string        `json:"contact_email"`
	ContactPhone     string        `json:"contact_phone,omitempty"`
	Quotas           []QuotaInfo   `json:"quotas,omitempty"`
	CreatedAt        int64         `json:"created_at"`
	UpdatedAt        int64         `json:"updated_at"`
}

// CreateTenantResponse 创建租户响应
type CreateTenantResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    TenantInfo  `json:"data"`
}

// UpdateTenantResponse 更新租户响应
type UpdateTenantResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    TenantInfo  `json:"data"`
}

// GetTenantResponse 获取租户响应
type GetTenantResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    TenantDetail  `json:"data"`
}

// ListTenantsResponse 列出租户响应
type ListTenantsResponse struct {
	Code        int          `json:"code"`
	Message     string       `json:"message"`
	Data        ListTenantsData `json:"data"`
}

// ListTenantsData 列出租户数据
type ListTenantsData struct {
	Tenants      []TenantInfo `json:"tenants"`
	TotalCount   int          `json:"total_count"`
	NextPageToken string       `json:"next_page_token,omitempty"`
}

// SubscriptionInfo 订阅信息
type SubscriptionInfo struct {
	SubscriptionID  string `json:"subscription_id"`
	TenantID        string `json:"tenant_id"`
	PlanTier        string `json:"plan_tier"`
	BillingCycle    string `json:"billing_cycle"`
	Status          string `json:"status"`
	StartedAt       int64  `json:"started_at"`
	ExpiresAt       int64  `json:"expires_at"`
	AutoRenew       bool   `json:"auto_renew"`
}

// SubscriptionResponse 订阅响应
type SubscriptionResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    SubscriptionInfo  `json:"data"`
}

// QuotaInfo 配额信息
type QuotaInfo struct {
	ResourceType  string  `json:"resource_type"`
	MaxLimit      int     `json:"max_limit"`
	UsedCount     int     `json:"used_count"`
	Remaining     int     `json:"remaining"`
	UsagePercent  float64 `json:"usage_percent"`
	ResetCycle    string  `json:"reset_cycle"`
	LastResetAt   int64   `json:"last_reset_at"`
	AlertLevel    string  `json:"alert_level"`
}

// QuotasResponse 配额响应
type QuotasResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    QuotasData   `json:"data"`
}

// QuotasData 配额数据
type QuotasData struct {
	Quotas []QuotaInfo `json:"quotas"`
}

// CheckQuotaResponse 检查配额响应
type CheckQuotaResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    CheckQuotaData   `json:"data"`
}

// CheckQuotaData 检查配额数据
type CheckQuotaData struct {
	Allowed       bool `json:"allowed"`
	CurrentUsage  int  `json:"current_usage"`
	MaxLimit      int  `json:"max_limit"`
	Remaining     int  `json:"remaining"`
}

// ConsumeQuotaResponse 消费配额响应
type ConsumeQuotaResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    ConsumeQuotaData  `json:"data"`
}

// ConsumeQuotaData 消费配额数据
type ConsumeQuotaData struct {
	ResourceType   string `json:"resource_type"`
	PreviousCount  int    `json:"previous_count"`
	CurrentCount   int    `json:"current_count"`
	Remaining      int    `json:"remaining"`
}

// RecordUsageResponse 记录使用量响应
type RecordUsageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// UsageSummary 使用量汇总
type UsageSummary struct {
	ResourceType    string                 `json:"resource_type"`
	TotalQuantity   int64                  `json:"total_quantity"`
	ActionCounts    map[string]int64       `json:"action_counts"`
	FirstUsage      string                 `json:"first_usage"`
	LastUsage       string                 `json:"last_usage"`
}

// GetUsageSummaryResponse 获取使用量汇总响应
type GetUsageSummaryResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    GetUsageSummaryData  `json:"data"`
}

// GetUsageSummaryData 获取使用量汇总数据
type GetUsageSummaryData struct {
	Summaries []UsageSummary `json:"summaries"`
}

// InvoiceInfo 账单信息
type InvoiceInfo struct {
	InvoiceID     string  `json:"invoice_id"`
	TenantID      string  `json:"tenant_id"`
	BillingCycle  string  `json:"billing_cycle"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	TotalUsage    int64   `json:"total_usage"`
	TotalAmount   float64 `json:"total_amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	DueDate       string  `json:"due_date"`
	CreatedAt     string  `json:"created_at"`
	PaidAt        string  `json:"paid_at,omitempty"`
}

// GenerateInvoiceResponse 生成账单响应
type GenerateInvoiceResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    InvoiceInfo   `json:"data"`
}

// ListInvoicesResponse 列出账单响应
type ListInvoicesResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    ListInvoicesData `json:"data"`
}

// ListInvoicesData 列出账单数据
type ListInvoicesData struct {
	Invoices    []InvoiceInfo `json:"invoices"`
	TotalCount  int           `json:"total_count"`
}

// PayInvoiceResponse 支付账单响应
type PayInvoiceResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    PayInvoiceData `json:"data"`
}

// PayInvoiceData 支付账单数据
type PayInvoiceData struct {
	InvoiceID string `json:"invoice_id"`
	Status    string `json:"status"`
	PaidAt    int64  `json:"paid_at"`
}

// MonitoringStatus 配额监控状态
type MonitoringStatus struct {
	TenantID     string            `json:"tenant_id"`
	Quotas       []QuotaInfo       `json:"quotas"`
	Alerts       []AlertInfo       `json:"alerts,omitempty"`
}

// AlertInfo 告警信息
type AlertInfo struct {
	ResourceType string `json:"resource_type"`
	AlertType    string `json:"alert_type"`
	Message      string `json:"message"`
}

// MonitoringResponse 监控响应
type MonitoringResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    MonitoringStatus   `json:"data"`
}
