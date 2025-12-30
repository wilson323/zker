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

// AuditLogDO 审计日志数据对象
type AuditLogDO struct {
	LogID         string        `gorm:"column:log_id;primaryKey;type:varchar(36)" json:"log_id"`
	TenantID      string        `gorm:"column:tenant_id;type:varchar(36);not null;index:idx_tenant_id" json:"tenant_id"`
	UserID        string        `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_id" json:"user_id"`
	Username      string        `gorm:"column:username;type:varchar(100);not null" json:"username"`
	UserEmail     string        `gorm:"column:user_email;type:varchar(255);not null" json:"user_email"`

	Action        string        `gorm:"column:action;type:varchar(50);not null;index:idx_action" json:"action"`
	Resource      string        `gorm:"column:resource;type:varchar(50);not null;index:idx_resource" json:"resource"`
	ResourceID    string        `gorm:"column:resource_id;type:varchar(36);index:idx_resource_id" json:"resource_id"`
	ResourceName  string        `gorm:"column:resource_name;type:varchar(255)" json:"resource_name"`

	RequestMethod string        `gorm:"column:request_method;type:varchar(10)" json:"request_method"`
	RequestPath   string        `gorm:"column:request_path;type:varchar(500);index:idx_request_path" json:"request_path"`
	RequestIP     string        `gorm:"column:request_ip;type:varchar(45);index:idx_request_ip" json:"request_ip"`
	UserAgent     string        `gorm:"column:user_agent;type:varchar(500)" json:"user_agent"`

	Status        string        `gorm:"column:status;type:enum('success','failed','pending');not null;index:idx_status" json:"status"`
	ErrorCode     string        `gorm:"column:error_code;type:varchar(50)" json:"error_code"`
	ErrorMsg      string        `gorm:"column:error_msg;type:text" json:"error_msg"`

	RequestData   string        `gorm:"column:request_data;type:longtext" json:"request_data"`
	ResponseData  string        `gorm:"column:response_data;type:longtext" json:"response_data"`
	Changes       string        `gorm:"column:changes;type:longtext" json:"changes"`
	Metadata      string        `gorm:"column:metadata;type:json" json:"metadata"`

	SessionID     string        `gorm:"column:session_id;type:varchar(36);index:idx_session_id" json:"session_id"`
	TraceID       string        `gorm:"column:trace_id;type:varchar(36);index:idx_trace_id" json:"trace_id"`

	CreatedAt     time.Time     `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index:idx_created_at" json:"created_at"`
	ArchivedAt    *time.Time    `gorm:"column:archived_at;index" json:"archived_at"`

	Signature     string        `gorm:"column:signature;type:varchar(64);index" json:"signature"`
	IsArchived    bool          `gorm:"column:is_archived;type:tinyint(1);not null;default:0;index:idx_is_archived" json:"is_archived"`
}

// TableName 指定表名
func (AuditLogDO) TableName() string {
	return "audit_logs"
}
