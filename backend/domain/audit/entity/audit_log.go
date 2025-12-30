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

// AuditAction 审计操作类型
type AuditAction string

const (
	// 用户操作
	AuditActionUserLogin       AuditAction = "user.login"        // 用户登录
	AuditActionUserLogout      AuditAction = "user.logout"       // 用户登出
	AuditActionUserCreate      AuditAction = "user.create"       // 创建用户
	AuditActionUserUpdate      AuditAction = "user.update"       // 更新用户
	AuditActionUserDelete      AuditAction = "user.delete"       // 删除用户
	AuditActionUserPassword    AuditAction = "user.password"     // 修改密码
	AuditActionUserMFAEnable   AuditAction = "user.mfa.enable"   // 启用MFA
	AuditActionUserMFADisable  AuditAction = "user.mfa.disable"  // 禁用MFA

	// 权限操作
	AuditActionRoleCreate      AuditAction = "role.create"       // 创建角色
	AuditActionRoleUpdate      AuditAction = "role.update"       // 更新角色
	AuditActionRoleDelete      AuditAction = "role.delete"       // 删除角色
	AuditActionRoleAssign      AuditAction = "role.assign"       // 分配角色
	AuditActionRoleRevoke      AuditAction = "role.revoke"       // 撤销角色

	// 数据操作
	AuditActionDataCreate      AuditAction = "data.create"       // 创建数据
	AuditActionDataUpdate      AuditAction = "data.update"       // 更新数据
	AuditActionDataDelete      AuditAction = "data.delete"       // 删除数据
	AuditActionDataQuery       AuditAction = "data.query"        // 查询数据
	AuditActionDataExport      AuditAction = "data.export"       // 导出数据
	AuditActionDataImport      AuditAction = "data.import"       // 导入数据

	// 配置操作
	AuditActionConfigUpdate    AuditAction = "config.update"     // 更新配置
	AuditActionConfigDelete    AuditAction = "config.delete"     // 删除配置

	// 系统操作
	AuditActionSystemBackup    AuditAction = "system.backup"     // 系统备份
	AuditActionSystemRestore   AuditAction = "system.restore"    // 系统恢复
	AuditActionSystemUpgrade   AuditAction = "system.upgrade"    // 系统升级
)

// AuditResource 审计资源类型
type AuditResource string

const (
	AuditResourceUser       AuditResource = "user"        // 用户
	AuditResourceRole       AuditResource = "role"        // 角色
	AuditResourceTenant     AuditResource = "tenant"      // 租户
	AuditResourceBot        AuditResource = "bot"         // Bot
	AuditResourceConversation AuditResource = "conversation" // 对话
	AuditResourceKnowledge  AuditResource = "knowledge"   // 知识库
	AuditResourceWorkflow   AuditResource = "workflow"    // 工作流
	AuditResourceDatabase   AuditResource = "database"    // 数据库
	AuditResourceConfig     AuditResource = "config"      // 配置
	AuditResourceSystem     AuditResource = "system"      // 系统
)

// AuditStatus 审计状态
type AuditStatus string

const (
	AuditStatusSuccess AuditStatus = "success" // 成功
	AuditStatusFailed  AuditStatus = "failed"  // 失败
	AuditStatusPending AuditStatus = "pending" // 进行中
)

// AuditLog 审计日志实体
type AuditLog struct {
	LogID         string        `json:"log_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string        `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID        string        `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	Username      string        `json:"username" gorm:"type:varchar(100);not null"`
	UserEmail     string        `json:"user_email" gorm:"type:varchar(255);not null"`

	// 操作信息
	Action        AuditAction   `json:"action" gorm:"type:varchar(50);not null;index:idx_action"`
	Resource      AuditResource `json:"resource" gorm:"type:varchar(50);not null;index:idx_resource"`
	ResourceID    string        `json:"resource_id" gorm:"type:varchar(36);index:idx_resource_id"`
	ResourceName  string        `json:"resource_name" gorm:"type:varchar(255)"`

	// 请求信息
	RequestMethod string        `json:"request_method" gorm:"type:varchar(10)"`
	RequestPath   string        `json:"request_path" gorm:"type:varchar(500);index:idx_request_path"`
	RequestIP     string        `json:"request_ip" gorm:"type:varchar(45);index:idx_request_ip"`
	UserAgent     string        `json:"user_agent" gorm:"type:varchar(500)"`

	// 状态和结果
	Status        AuditStatus   `json:"status" gorm:"type:enum('success','failed','pending');not null;index:idx_status"`
	ErrorCode     string        `json:"error_code" gorm:"type:varchar(50)"`
	ErrorMsg      string        `json:"error_msg" gorm:"type:text"`

	// 详细信息
	RequestData   string        `json:"request_data" gorm:"type:longtext"`  // 敏感数据已脱敏
	ResponseData  string        `json:"response_data" gorm:"type:longtext"` // 敏感数据已脱敏
	Changes       string        `json:"changes" gorm:"type:longtext"`       // 变更内容(JSON格式)
	Metadata      string        `json:"metadata" gorm:"type:json"`          // 额外元数据(JSON格式)

	// 合规字段
	SessionID     string        `json:"session_id" gorm:"type:varchar(36);index:idx_session_id"`
	TraceID       string        `json:"trace_id" gorm:"type:varchar(36);index:idx_trace_id"`

	// 时间戳
	CreatedAt     time.Time     `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_created_at"`
	ArchivedAt    *time.Time    `json:"archived_at,omitempty" gorm:"index"`

	// 数字签名(用于防篡改)
	Signature     string        `json:"signature" gorm:"type:varchar(64);index"` // SHA-256签名
	IsArchived    bool          `json:"is_archived" gorm:"not null;default:false;index:idx_is_archived"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// IsSuccessful 是否成功
func (a *AuditLog) IsSuccessful() bool {
	return a.Status == AuditStatusSuccess
}

// IsFailed 是否失败
func (a *AuditLog) IsFailed() bool {
	return a.Status == AuditStatusFailed
}

// IsSensitiveOperation 是否是敏感操作
func (a *AuditLog) IsSensitiveOperation() bool {
	sensitiveActions := map[AuditAction]bool{
		AuditActionUserDelete:      true,
		AuditActionRoleDelete:      true,
		AuditActionDataDelete:      true,
		AuditActionDataExport:      true,
		AuditActionUserPassword:    true,
		AuditActionConfigUpdate:    true,
		AuditActionConfigDelete:    true,
		AuditActionSystemRestore:   true,
	}
	return sensitiveActions[a.Action]
}

// GetTimestamp 获取时间戳(毫秒)
func (a *AuditLog) GetTimestamp() int64 {
	return a.CreatedAt.UnixMilli()
}

// LogFilter 审计日志查询过滤器
type LogFilter struct {
	TenantID      string        `json:"tenant_id"`
	UserID        string        `json:"user_id,omitempty"`
	UserEmail     string        `json:"user_email,omitempty"`
	Actions       []AuditAction `json:"actions,omitempty"`
	Resources     []AuditResource `json:"resources,omitempty"`
	ResourceID    string        `json:"resource_id,omitempty"`
	Status        AuditStatus   `json:"status,omitempty"`
	StartTime     *time.Time    `json:"start_time,omitempty"`
	EndTime       *time.Time    `json:"end_time,omitempty"`
	RequestIP     string        `json:"request_ip,omitempty"`
	SessionID     string        `json:"session_id,omitempty"`
	TraceID       string        `json:"trace_id,omitempty"`

	// 分页参数
	Page          int           `json:"page"`
	PageSize      int           `json:"page_size"`

	// 排序
	OrderBy       string        `json:"order_by"` // created_at, user_email, action, etc.
	OrderDesc     bool          `json:"order_desc"`
}

// Validate 验证过滤器
func (f *LogFilter) Validate() error {
	if f.TenantID == "" {
		return ErrTenantIDRequired
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 1000 {
		f.PageSize = 20
	}
	if f.OrderBy == "" {
		f.OrderBy = "created_at"
		f.OrderDesc = true
	}
	return nil
}

// ExportResult 导出结果
type ExportResult struct {
	FileID       string    `json:"file_id"`        // 导出文件ID
	FileName     string    `json:"file_name"`      // 导出文件名
	FileSize     int64     `json:"file_size"`      // 文件大小(字节)
	RecordCount  int       `json:"record_count"`   // 记录数
	Format       string    `json:"format"`         // 导出格式(csv, excel, json)
	ContentType  string    `json:"content_type"`   // MIME类型(用于HTTP响应头)
	DownloadURL  string    `json:"download_url"`   // 下载链接
	ExpiresAt    time.Time `json:"expires_at"`     // 过期时间
	CreatedAt    time.Time `json:"created_at"`     // 创建时间
}

// ArchiveStats 归档统计
type ArchiveStats struct {
	TotalLogs        int       `json:"total_logs"`         // 总日志数
	ArchivedLogs     int       `json:"archived_logs"`      // 已归档日志数
	ArchiveDate      time.Time `json:"archive_date"`       // 归档日期
	ArchiveSize      int64     `json:"archive_size"`       // 归档大小(字节)
	ArchiveFilePath  string    `json:"archive_file_path"`  // 归档文件路径
}

// AuditStats 审计统计
type AuditStats struct {
	TotalLogs      int64              `json:"total_logs"`       // 总日志数
	SuccessLogs    int64              `json:"success_logs"`     // 成功日志数
	FailedLogs     int64              `json:"failed_logs"`      // 失败日志数
	ActionStats    map[string]int64   `json:"action_stats"`     // 操作统计
	ResourceStats  map[string]int64   `json:"resource_stats"`   // 资源统计
	UserStats      map[string]int64   `json:"user_stats"`       // 用户统计
	DailyStats     []DailyStat        `json:"daily_stats"`      // 每日统计
}

// DailyStat 每日统计
type DailyStat struct {
	Date   string `json:"date"`   // 日期(YYYY-MM-DD)
	Count  int64  `json:"count"`  // 日志数
}

// 错误定义
var (
	ErrTenantIDRequired   = &ValidationError{Field: "tenant_id", Message: "租户ID必填"}
	ErrInvalidDateRange   = &ValidationError{Field: "date_range", Message: "日期范围无效"}
	ErrInvalidExportFormat = &ValidationError{Field: "format", Message: "不支持的导出格式"}
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
