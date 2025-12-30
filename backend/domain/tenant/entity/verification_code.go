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

// VerificationCode 验证码实体
type VerificationCode struct {
	Email      string    `json:"email" gorm:"primaryKey;type:varchar(255)"`
	Code       string    `json:"code" gorm:"type:varchar(10);not null"`
	ExpiresAt  int64     `json:"expires_at" gorm:"not null"`
	CreatedAt  int64     `json:"created_at" gorm:"not null"`
	VerifiedAt *int64    `json:"verified_at,omitempty"`
}

// TableName 指定表名
func (VerificationCode) TableName() string {
	return "verification_codes"
}

// IsExpired 检查验证码是否已过期
func (v *VerificationCode) IsExpired() bool {
	return time.Now().UnixMilli() > v.ExpiresAt
}

// IsVerified 检查验证码是否已验证
func (v *VerificationCode) IsVerified() bool {
	return v.VerifiedAt != nil
}

// GetExpiresAtAsTime 获取过期时间
func (v *VerificationCode) GetExpiresAtAsTime() time.Time {
	return time.Unix(v.ExpiresAt/1000, 0)
}
