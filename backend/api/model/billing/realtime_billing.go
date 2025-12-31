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

package billing

// ============================================================
// 实时计费API Model定义
// ============================================================

// ChargeRequest 实时扣费请求
type ChargeRequest struct {
	TenantID    string  `json:"tenant_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,oneof=CNY USD EUR"`
	Description string  `json:"description" binding:"required,max=500"`
	Metadata    string  `json:"metadata,omitempty"`
}

// BalanceResponse 余额响应
type BalanceResponse struct {
	TenantID        string  `json:"tenant_id"`
	CurrentBalance  float64 `json:"current_balance"`
	Currency        string  `json:"currency"`
	AvailableAmount float64 `json:"available_amount"`
	FrozenAmount    float64 `json:"frozen_amount"`
	CreditLimit     float64 `json:"credit_limit"`
	UpdatedAt       int64   `json:"updated_at"`
}

// ============================================================
// 配额管理API Model定义
// ============================================================

// ConsumeQuotaRequest 消费配额请求
type ConsumeQuotaRequest struct {
	TenantID     string `json:"tenant_id" binding:"required"`
	ResourceType string `json:"resource_type" binding:"required,oneof=bots knowledge workflows api_calls storage concurrent"`
	Quantity     int    `json:"quantity" binding:"required,min=1"`
	OperationID  string `json:"operation_id" binding:"required"`
	Description  string `json:"description,omitempty"`
}

// ConsumeQuotaResponse 消费配额响应
type ConsumeQuotaResponse struct {
	Success     bool  `json:"success"`
	NewQuantity int   `json:"new_quantity"`
	Remaining   int   `json:"remaining"`
}

// CheckQuotaResponse 检查配额响应
type CheckQuotaResponse struct {
	TenantID     string `json:"tenant_id"`
	ResourceType string `json:"resource_type"`
	Sufficient   bool   `json:"sufficient"`
	Current      int    `json:"current"`
	Limit        int    `json:"limit"`
	Remaining    int    `json:"remaining"`
}

// ResetQuotaRequest 重置配额请求
type ResetQuotaRequest struct {
	TenantID     string  `json:"tenant_id" binding:"required"`
	ResourceType string  `json:"resource_type" binding:"required,oneof=bots knowledge workflows api_calls storage concurrent"`
	NewLimit     int     `json:"new_limit" binding:"required,min=0"`
	Reason       string `json:"reason" binding:"required,max=500"`
}

// RollbackQuotaRequest 回滚配额请求
type RollbackQuotaRequest struct {
	TenantID    string `json:"tenant_id" binding:"required"`
	OperationID string `json:"operation_id" binding:"required"`
	Reason      string `json:"reason" binding:"required,max=500"`
}
