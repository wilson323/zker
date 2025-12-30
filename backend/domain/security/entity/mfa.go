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

// MFAType MFA类型
type MFAType string

const (
	MFATypeTOTP        MFAType = "totp"        // 基于时间的一次性密码
	MFATypeSMS         MFAType = "sms"         // 短信验证码
	MFATypeEmail       MFAType = "email"       // 邮箱验证码
	MFATypeHardwareKey MFAType = "hardware"    // 硬件密钥(YubiKey等)
)

// MFAStatus MFA状态
type MFAStatus string

const (
	MFAStatusEnabled  MFAStatus = "enabled"  // 已启用
	MFAStatusDisabled MFAStatus = "disabled" // 已禁用
	MFAStatusLocked   MFAStatus = "locked"   // 已锁定
)

// MFAConfig MFA配置实体
type MFAConfig struct {
	ConfigID      string     `json:"config_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string     `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID        string     `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`

	// MFA配置
	Type          MFAType    `json:"type" gorm:"type:enum('totp','sms','email','hardware');not null"`
	Status        MFAStatus  `json:"status" gorm:"type:enum('enabled','disabled','locked');not null;default:'disabled'"`
	Secret        string     `json:"-" gorm:"type:varchar(255);not null"` // TOTP密钥(加密存储)
	BackupCodes   string     `json:"-" gorm:"type:text"`                   // 恢复码(加密存储,JSON数组)

	// TOTP配置
	Issuer        string     `json:"issuer" gorm:"type:varchar(100)"`      // 发行者名称
	Account       string     `json:"account" gorm:"type:varchar(100)"`     // 账户名称
	Algorithm     string     `json:"algorithm" gorm:"type:varchar(20);default:'SHA1'"` // 算法(SHA1, SHA256, SHA512)
	Digits        int        `json:"digits" gorm:"not null;default:6"`      // 位数(6或8)
	Period        int        `json:"period" gorm:"not null;default:30"`     // 时间步长(秒)

	// 备用验证方式
	BackupType    MFAType    `json:"backup_type" gorm:"type:enum('totp','sms','email','hardware')"`
	BackupValue   string     `json:"-" gorm:"type:varchar(255)"`           // 备用方式值(如手机号、邮箱)

	// 统计信息
	UsedCount     int        `json:"used_count" gorm:"not null;default:0"`  // 使用次数
	LastUsedAt    *time.Time `json:"last_used_at,omitempty"`                // 最后使用时间
	FailedAttempts int       `json:"failed_attempts" gorm:"not null;default:0"` // 失败次数
	LockedUntil   *time.Time `json:"locked_until,omitempty"`                // 锁定到期时间

	// 时间戳
	CreatedAt     time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (MFAConfig) TableName() string {
	return "mfa_configs"
}

// IsEnabled 是否已启用
func (m *MFAConfig) IsEnabled() bool {
	return m.Status == MFAStatusEnabled
}

// IsLocked 是否已锁定
func (m *MFAConfig) IsLocked() bool {
	return m.Status == MFAStatusLocked
}

// ShouldLock 是否应该锁定(失败次数超过阈值)
func (m *MFAConfig) ShouldLock(maxAttempts int) bool {
	return m.FailedAttempts >= maxAttempts
}

// IsLockExpired 锁定是否已过期
func (m *MFAConfig) IsLockExpired() bool {
	if !m.IsLocked() || m.LockedUntil == nil {
		return false
	}
	return time.Now().After(*m.LockedUntil)
}

// RecoveryCode 恢复码实体
type RecoveryCode struct {
	CodeID      string     `json:"code_id" gorm:"primaryKey;type:varchar(36)"`
	ConfigID    string     `json:"config_id" gorm:"type:varchar(36);not null;index:idx_config_id"`
	UserID      string     `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	TenantID    string     `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`

	CodeHash    string     `json:"-" gorm:"type:varchar(255);not null;unique"` // 恢复码哈希值
	Used        bool       `json:"used" gorm:"type:tinyint(1);not null;default:0;index:idx_used"`
	UsedAt      *time.Time `json:"used_at,omitempty"`

	CreatedAt   time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (RecoveryCode) TableName() string {
	return "recovery_codes"
}

// IsExpired 是否已过期
func (r *RecoveryCode) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.ExpiresAt)
}

// IsValid 是否有效(未使用且未过期)
func (r *RecoveryCode) IsValid() bool {
	return !r.Used && !r.IsExpired()
}

// MFAVerification MFA验证记录
type MFAVerification struct {
	VerifyID    string     `json:"verify_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID    string     `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID      string     `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`

	Type        MFAType    `json:"type" gorm:"type:enum('totp','sms','email','hardware');not null"`
	Success     bool       `json:"success" gorm:"type:tinyint(1);not null"`
	FailureReason string    `json:"failure_reason,omitempty" gorm:"type:varchar(255)"`

	RequestIP   string     `json:"request_ip" gorm:"type:varchar(45)"`
	UserAgent   string     `json:"user_agent" gorm:"type:varchar(500)"`

	CreatedAt   time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_created_at"`
}

// TableName 指定表名
func (MFAVerification) TableName() string {
	return "mfa_verifications"
}

// MFALoginSession MFA登录会话(用于多步骤验证)
type MFALoginSession struct {
	SessionID   string     `json:"session_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID    string     `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID      string     `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`

	// 第一步认证信息(密码/邮箱验证)
	FirstFactorPassed bool       `json:"first_factor_passed" gorm:"type:tinyint(1);not null;default:0"`
	FirstFactorAt     *time.Time `json:"first_factor_at,omitempty"`

	// 第二步认证信息(MFA)
	SecondFactorPassed bool      `json:"second_factor_passed" gorm:"type:tinyint(1);not null;default:0"`
	SecondFactorAt     *time.Time `json:"second_factor_at,omitempty"`

	// 会话状态
	Completed   bool       `json:"completed" gorm:"type:tinyint(1);not null;default:0;index:idx_completed"`
	ExpiresAt   time.Time  `json:"expires_at" gorm:"not null;index:idx_expires_at"`

	CreatedAt   time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (MFALoginSession) TableName() string {
	return "mfa_login_sessions"
}

// IsExpired 是否已过期
func (s *MFALoginSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsComplete 是否已完成两步验证
func (s *MFALoginSession) IsComplete() bool {
	return s.FirstFactorPassed && s.SecondFactorPassed
}

// MFAConfigFilter MFA配置查询过滤器
type MFAConfigFilter struct {
	TenantID   string    `json:"tenant_id"`
	UserID     string    `json:"user_id,omitempty"`
	Type       MFAType   `json:"type,omitempty"`
	Status     MFAStatus `json:"status,omitempty"`
}

// Validate 验证过滤器
func (f *MFAConfigFilter) Validate() error {
	if f.TenantID == "" {
		return ErrTenantIDRequired
	}
	return nil
}

// 错误定义
var (
	ErrTenantIDRequired = &ValidationError{Field: "tenant_id", Message: "租户ID必填"}
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
