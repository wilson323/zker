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

package model

import (
	"time"
)

// MFAConfigDO MFA配置数据对象
type MFAConfigDO struct {
	ConfigID       string        `gorm:"column:config_id;primaryKey;type:varchar(36)" json:"config_id"`
	TenantID       string        `gorm:"column:tenant_id;type:varchar(36);not null;index:idx_tenant_id" json:"tenant_id"`
	UserID         string        `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_id" json:"user_id"`

	Type           string        `gorm:"column:type;type:enum('totp','sms','email','hardware');not null" json:"type"`
	Status         string        `gorm:"column:status;type:enum('enabled','disabled','locked');not null;default:'disabled'" json:"status"`
	Secret         string        `gorm:"column:secret;type:varchar(255);not null" json:"-"`
	BackupCodes    string        `gorm:"column:backup_codes;type:text" json:"-"`

	Issuer         string        `gorm:"column:issuer;type:varchar(100)" json:"issuer"`
	Account        string        `gorm:"column:account;type:varchar(100)" json:"account"`
	Algorithm      string        `gorm:"column:algorithm;type:varchar(20);default:'SHA1'" json:"algorithm"`
	Digits         int           `gorm:"column:digits;not null;default:6" json:"digits"`
	Period         int           `gorm:"column:period;not null;default:30" json:"period"`

	BackupType     string        `gorm:"column:backup_type;type:enum('totp','sms','email','hardware')" json:"backup_type"`
	BackupValue    string        `gorm:"column:backup_value;type:varchar(255)" json:"-"`

	UsedCount      int           `gorm:"column:used_count;not null;default:0" json:"used_count"`
	LastUsedAt     *time.Time    `gorm:"column:last_used_at" json:"last_used_at,omitempty"`
	FailedAttempts int           `gorm:"column:failed_attempts;not null;default:0" json:"failed_attempts"`
	LockedUntil    *time.Time    `gorm:"column:locked_until" json:"locked_until,omitempty"`

	CreatedAt      time.Time     `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      *time.Time    `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (MFAConfigDO) TableName() string {
	return "mfa_configs"
}

// RecoveryCodeDO 恢复码数据对象
type RecoveryCodeDO struct {
	CodeID    string     `gorm:"column:code_id;primaryKey;type:varchar(36)" json:"code_id"`
	ConfigID  string     `gorm:"column:config_id;type:varchar(36);not null;index:idx_config_id" json:"config_id"`
	UserID    string     `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_id" json:"user_id"`
	TenantID  string     `gorm:"column:tenant_id;type:varchar(36);not null;index:idx_tenant_id" json:"tenant_id"`

	CodeHash  string     `gorm:"column:code_hash;type:varchar(255);not null;unique" json:"-"`
	Used      bool       `gorm:"column:used;type:tinyint(1);not null;default:0;index:idx_used" json:"used"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"used_at,omitempty"`

	CreatedAt time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt *time.Time `gorm:"column:expires_at" json:"expires_at,omitempty"`
}

// TableName 指定表名
func (RecoveryCodeDO) TableName() string {
	return "recovery_codes"
}

// MFAVerificationDO MFA验证记录数据对象
type MFAVerificationDO struct {
	VerifyID       string    `gorm:"column:verify_id;primaryKey;type:varchar(36)" json:"verify_id"`
	TenantID       string    `gorm:"column:tenant_id;type:varchar(36);not null;index:idx_tenant_id" json:"tenant_id"`
	UserID         string    `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_id" json:"user_id"`

	Type           string    `gorm:"column:type;type:enum('totp','sms','email','hardware');not null" json:"type"`
	Success        bool      `gorm:"column:success;type:tinyint(1);not null" json:"success"`
	FailureReason  string    `gorm:"column:failure_reason;type:varchar(255)" json:"failure_reason,omitempty"`

	RequestIP      string    `gorm:"column:request_ip;type:varchar(45)" json:"request_ip"`
	UserAgent      string    `gorm:"column:user_agent;type:varchar(500)" json:"user_agent"`

	CreatedAt      time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index:idx_created_at" json:"created_at"`
}

// TableName 指定表名
func (MFAVerificationDO) TableName() string {
	return "mfa_verifications"
}

// MFALoginSessionDO MFA登录会话数据对象
type MFALoginSessionDO struct {
	SessionID           string     `gorm:"column:session_id;primaryKey;type:varchar(36)" json:"session_id"`
	TenantID            string     `gorm:"column:tenant_id;type:varchar(36);not null;index:idx_tenant_id" json:"tenant_id"`
	UserID              string     `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_id" json:"user_id"`

	FirstFactorPassed   bool       `gorm:"column:first_factor_passed;type:tinyint(1);not null;default:0" json:"first_factor_passed"`
	FirstFactorAt       *time.Time `gorm:"column:first_factor_at" json:"first_factor_at,omitempty"`

	SecondFactorPassed  bool       `gorm:"column:second_factor_passed;type:tinyint(1);not null;default:0" json:"second_factor_passed"`
	SecondFactorAt      *time.Time `gorm:"column:second_factor_at" json:"second_factor_at,omitempty"`

	Completed           bool       `gorm:"column:completed;type:tinyint(1);not null;default:0;index:idx_completed" json:"completed"`
	ExpiresAt           time.Time  `gorm:"column:expires_at;not null;index:idx_expires_at" json:"expires_at"`

	CreatedAt           time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 指定表名
func (MFALoginSessionDO) TableName() string {
	return "mfa_login_sessions"
}
