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
// 账单生成API Model定义
// ============================================================

// GenerateBillRequest 生成账单请求
type GenerateBillRequest struct {
	TenantID      string `json:"tenant_id" binding:"required"`
	BillingCycle  string `json:"billing_cycle" binding:"required"` // YYYY-MM格式
	BillType      string `json:"bill_type" binding:"required,oneof=subscription usage onetime"`
	Description   string `json:"description,omitempty"`
}

// BillDetail 账单详情
type BillDetail struct {
	BillID        string  `json:"bill_id"`
	TenantID      string  `json:"tenant_id"`
	BillingCycle  string  `json:"billing_cycle"`
	BillType      string  `json:"bill_type"`
	TotalAmount   float64 `json:"total_amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"` // pending, paid, overdue, cancelled
	Description   string  `json:"description,omitempty"`
	DueDate       *int64  `json:"due_date,omitempty"`
	PaidAt        *int64  `json:"paid_at,omitempty"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
}

// BillSummary 账单摘要（用于列表）
type BillSummary struct {
	BillID        string  `json:"bill_id"`
	BillingCycle  string  `json:"billing_cycle"`
	BillType      string  `json:"bill_type"`
	TotalAmount   float64 `json:"total_amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	Description   string  `json:"description,omitempty"`
	CreatedAt     int64   `json:"created_at"`
	DueDate       int64   `json:"due_date"`
}

// PayBillRequest 支付账单请求
type PayBillRequest struct {
	TenantID      string  `json:"tenant_id" binding:"required"`
	BillID        string  `json:"bill_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay wechat bank_transfer"`
	ReturnURL     string  `json:"return_url,omitempty"`
}

// PaymentInfo 支付信息
type PaymentInfo struct {
	PaymentID     string `json:"payment_id"`
	PaymentNumber string `json:"payment_number"`
	Amount        float64 `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	PaymentURL    string `json:"payment_url,omitempty"`
	QRCodeURL     string `json:"qr_code_url,omitempty"`
	ExpiresAt     int64  `json:"expires_at"`
	CreatedAt     int64  `json:"created_at"`
}

// InvoiceDetail 发票详情
type InvoiceDetail struct {
	InvoiceID     string  `json:"invoice_id"`
	InvoiceNumber string  `json:"invoice_number"`
	BillID        string  `json:"bill_id"`
	TenantID      string  `json:"tenant_id"`
	TotalAmount   float64 `json:"total_amount"`
	Currency      string  `json:"currency"`
	TaxAmount     float64 `json:"tax_amount,omitempty"`
	Status        string  `json:"status"` // draft, issued, paid, cancelled
	IssuedAt      int64   `json:"issued_at"`
	DueDate       int64   `json:"due_date"`
	PaidAt        *int64  `json:"paid_at,omitempty"`
	Notes         string  `json:"notes,omitempty"`
}
