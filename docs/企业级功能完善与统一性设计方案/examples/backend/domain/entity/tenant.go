// Package domain 提供领域层实体定义
//
// 本文件展示租户（Tenant）实体的完整实现，包括：
// - 值对象（TenantID, TenantStatus）
// - 聚合根（Tenant）
// - 领域事件（TenantCreated, TenantUpgraded）
// - 工厂方法
// - 业务规则验证
package domain

import (
	"errors"
	"time"
)

var (
	// ErrInvalidTenantID 表示租户ID无效
	ErrInvalidTenantID = errors.New("invalid tenant ID")
	// ErrTenantNameRequired 表示租户名称必填
	ErrTenantNameRequired = errors.New("tenant name is required")
	// ErrTenantNameTooLong 表示租户名称过长
	ErrTenantNameTooLong = errors.New("tenant name must not exceed 100 characters")
	// ErrInvalidTenantStatus 表示租户状态无效
	ErrInvalidTenantStatus = errors.New("invalid tenant status")
	// ErrQuotaExceeded 表示配额超限
	ErrQuotaExceeded = errors.New("quota exceeded")
	// ErrTenantExpired 表示租户已过期
	ErrTenantExpired = errors.New("tenant has expired")
	// ErrTenantSuspended 表示租户已暂停
	ErrTenantSuspended = errors.New("tenant is suspended")
)

// TenantStatus 租户状态值对象
type TenantStatus string

const (
	// TenantStatusTrial 试用期
	TenantStatusTrial TenantStatus = "trial"
	// TenantStatusActive 正常活跃
	TenantStatusActive TenantStatus = "active"
	// TenantStatusSuspended 已暂停
	TenantStatusSuspended TenantStatus = "suspended"
	// TenantStatusExpired 已过期
	TenantStatusExpired TenantStatus = "expired"
)

// IsValid 验证租户状态是否有效
func (s TenantStatus) IsValid() bool {
	switch s {
	case TenantStatusTrial, TenantStatusActive, TenantStatusSuspended, TenantStatusExpired:
		return true
	default:
		return false
	}
}

// IsActive 检查租户是否处于活跃状态
func (s TenantStatus) IsActive() bool {
	return s == TenantStatusActive || s == TenantStatusTrial
}

// TenantID 租户ID值对象
type TenantID struct {
	value string
}

// NewTenantID 创建租户ID值对象
func NewTenantID(value string) (TenantID, error) {
	if value == "" {
		return TenantID{}, ErrInvalidTenantID
	}
	return TenantID{value: value}, nil
}

// String 返回租户ID的字符串表示
func (id TenantID) String() string {
	return id.value
}

// Subscription 订阅信息值对象
type Subscription struct {
	Plan          string    `json:"plan"`           // 套餐类型：free, pro, enterprise
	Quota         int64     `json:"quota"`          // Token配额
	Used          int64     `json:"used"`           // 已使用量
	ExpiresAt     time.Time `json:"expires_at"`     // 过期时间
	AutoRenew     bool      `json:"auto_renew"`     // 自动续费
	PaymentMethod string    `json:"payment_method"` // 支付方式
}

// IsExpired 检查订阅是否已过期
func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// HasQuota 检查是否还有配额
func (s *Subscription) HasQuota(tokens int64) bool {
	return (s.Used + tokens) <= s.Quota
}

// UseQuota 使用配额
func (s *Subscription) UseQuota(tokens int64) error {
	if !s.HasQuota(tokens) {
		return ErrQuotaExceeded
	}
	s.Used += tokens
	return nil
}

// Tenant 租户聚合根
//
// Tenant 是系统的核心聚合根，负责：
// - 维护租户的完整性和业务规则
// - 管理订阅和配额
// - 生成领域事件
type Tenant struct {
	id           TenantID
	name         string
	status       TenantStatus
	subscription Subscription
	contactEmail string
	contactPhone string
	createdAt    time.Time
	updatedAt    time.Time
	// 领域事件列表（用于事件溯源）
	events []interface{}
}

// NewTenant 创建新租户（工厂方法）
//
// 参数：
//   - id: 租户ID
//   - name: 租户名称
//   - contactEmail: 联系邮箱
//
// 返回：
//   - *Tenant: 租户聚合根
//   - error: 错误信息
func NewTenant(id TenantID, name, contactEmail string) (*Tenant, error) {
	// 业务规则验证
	if name == "" {
		return nil, ErrTenantNameRequired
	}
	if len(name) > 100 {
		return nil, ErrTenantNameTooLong
	}

	tenant := &Tenant{
		id:     id,
		name:   name,
		status: TenantStatusTrial, // 默认试用状态
		subscription: Subscription{
			Plan:      "free",
			Quota:     1000000, // 100万tokens免费额度
			Used:      0,
			ExpiresAt: time.Now().AddDate(0, 1, 0), // 试用期1个月
			AutoRenew: false,
		},
		contactEmail: contactEmail,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
		events:       []interface{}{},
	}

	// 记录领域事件
	tenant.recordEvent(TenantCreatedEvent{
		TenantID:    id.String(),
		Name:        name,
		OccurredAt:  time.Now(),
	})

	return tenant, nil
}

// Reconstruct 从持久化状态重建租户（用于仓储恢复）
func Reconstruct(
	id TenantID,
	name string,
	status TenantStatus,
	subscription Subscription,
	contactEmail, contactPhone string,
	createdAt, updatedAt time.Time,
) *Tenant {
	return &Tenant{
		id:           id,
		name:         name,
		status:       status,
		subscription: subscription,
		contactEmail: contactEmail,
		contactPhone: contactPhone,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		events:       []interface{}{},
	}
}

// Getters
func (t *Tenant) ID() TenantID              { return t.id }
func (t *Tenant) Name() string               { return t.name }
func (t *Tenant) Status() TenantStatus       { return t.status }
func (t *Tenant) Subscription() Subscription { return t.subscription }
func (t *Tenant) ContactEmail() string       { return t.contactEmail }
func (t *Tenant) ContactPhone() string       { return t.contactPhone }
func (t *Tenant) CreatedAt() time.Time       { return t.createdAt }
func (t *Tenant) UpdatedAt() time.Time       { return t.updatedAt }

// UpdateName 更新租户名称
func (t *Tenant) UpdateName(name string) error {
	if name == "" {
		return ErrTenantNameRequired
	}
	if len(name) > 100 {
		return ErrTenantNameTooLong
	}

	t.name = name
	t.updatedAt = time.Now()
	return nil
}

// UpdateContact 更新联系方式
func (t *Tenant) UpdateContact(email, phone string) {
	t.contactEmail = email
	t.contactPhone = phone
	t.updatedAt = time.Now()
}

// Activate 激活租户
func (t *Tenant) Activate() error {
	if t.status == TenantStatusActive {
		return nil // 已经是活跃状态
	}

	t.status = TenantStatusActive
	t.updatedAt = time.Now()

	return nil
}

// Suspend 暂停租户
func (t *Tenant) Suspend() error {
	if t.status == TenantStatusSuspended {
		return nil // 已经是暂停状态
	}

	t.status = TenantStatusSuspended
	t.updatedAt = time.Now()

	return nil
}

// CheckIsActive 检查租户是否处于活跃状态
func (t *Tenant) CheckIsActive() error {
	if !t.status.IsActive() {
		if t.status == TenantStatusExpired {
			return ErrTenantExpired
		}
		if t.status == TenantStatusSuspended {
			return ErrTenantSuspended
		}
		return ErrInvalidTenantStatus
	}

	if t.subscription.IsExpired() {
		return ErrTenantExpired
	}

	return nil
}

// UseTokens 使用Token配额
func (t *Tenant) UseTokens(tokens int64) error {
	// 先检查租户状态
	if err := t.CheckIsActive(); err != nil {
		return err
	}

	// 使用配额
	if err := t.subscription.UseQuota(tokens); err != nil {
		return err
	}

	t.updatedAt = time.Now()

	return nil
}

// UpgradePlan 升级套餐
func (t *Tenant) UpgradePlan(plan string, quota int64, expiresAt time.Time) error {
	if err := t.CheckIsActive(); err != nil {
		return err
	}

	oldPlan := t.subscription.Plan
	t.subscription.Plan = plan
	t.subscription.Quota = quota
	t.subscription.ExpiresAt = expiresAt
	t.updatedAt = time.Now()

	// 记录领域事件
	t.recordEvent(TenantUpgradedEvent{
		TenantID:    t.id.String(),
		OldPlan:     oldPlan,
		NewPlan:     plan,
		OccurredAt:  time.Now(),
	})

	return nil
}

// recordEvent 记录领域事件
func (t *Tenant) recordEvent(event interface{}) {
	t.events = append(t.events, event)
}

// GetEvents 获取并清空领域事件
func (t *Tenant) GetEvents() []interface{} {
	events := t.events
	t.events = []interface{}{} // 清空事件列表
	return events
}

// ============ 领域事件定义 ============

// TenantCreatedEvent 租户创建事件
type TenantCreatedEvent struct {
	TenantID   string    `json:"tenant_id"`
	Name       string    `json:"name"`
	OccurredAt time.Time `json:"occurred_at"`
}

// TenantUpgradedEvent 租户升级事件
type TenantUpgradedEvent struct {
	TenantID   string    `json:"tenant_id"`
	OldPlan    string    `json:"old_plan"`
	NewPlan    string    `json:"new_plan"`
	OccurredAt time.Time `json:"occurred_at"`
}

// TenantSuspendedEvent 租户暂停事件
type TenantSuspendedEvent struct {
	TenantID   string    `json:"tenant_id"`
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
}
