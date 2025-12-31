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

type User struct {
	UserID       int64
	TenantID     string // 租户ID（多租户隔离）

	Name         string // nickname
	UniqueName   string // unique name
	Email        string // email
	Password     string // password (encrypted)
	Description  string // user description
	IconURI      string // avatar URI
	IconURL      string // avatar URL
	UserVerified bool   // Is the user authenticated?
	Locale       string
	SessionKey   string // session key
	IsEnabled    bool   // is user enabled (for RBAC)

	CreatedAt int64 // creation time
	UpdatedAt int64 // update time
	DeletedAt *int64 // deletion time (soft delete)
}

// GetTenantID 获取租户ID
func (u *User) GetTenantID() string {
	return u.TenantID
}

// HasTenantID 检查是否有租户ID
func (u *User) HasTenantID() bool {
	return u.TenantID != ""
}

type UserBenefit struct {
	UserID        int64
	UserLevel     UserLevel
	UsedCount     int32
	TotalCount    int32
	IsUnlimited   bool
	ResetDatetime int64
	CallQPS       int32
}

type SaasUserData struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	NickName  string `json:"nick_name"`
	AvatarURL string `json:"avatar_url"`
}
