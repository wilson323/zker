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

package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/audit/entity"
	"github.com/coze-dev/coze-studio/backend/domain/audit/repository"
	"github.com/coze-dev/coze-studio/backend/domain/security/service"
)

// AuditLogService 审计日志服务接口
type AuditLogService interface {
	// LogOperation 记录操作日志
	LogOperation(ctx context.Context, log *entity.AuditLog) error

	// LogOperationAsync 异步记录操作日志(通过NSQ)
	LogOperationAsync(ctx context.Context, log *entity.AuditLog) error

	// QueryLogs 查询审计日志
	QueryLogs(ctx context.Context, filter *entity.LogFilter) ([]*entity.AuditLog, int64, error)

	// GetLogByID 根据ID获取日志
	GetLogByID(ctx context.Context, logID string) (*entity.AuditLog, error)

	// GetUserLogs 获取用户日志
	GetUserLogs(ctx context.Context, tenantID, userID string, page, pageSize int) ([]*entity.AuditLog, int64, error)

	// GetResourceLogs 获取资源日志
	GetResourceLogs(ctx context.Context, tenantID string, resource entity.AuditResource, resourceID string) ([]*entity.AuditLog, error)

	// GetSensitiveOperations 获取敏感操作日志
	GetSensitiveOperations(ctx context.Context, tenantID string, days int) ([]*entity.AuditLog, error)

	// GetStatistics 获取审计统计
	GetStatistics(ctx context.Context, tenantID string, days int) (*entity.AuditStats, error)

	// ExportLogs 导出审计日志
	ExportLogs(ctx context.Context, filter *entity.LogFilter, format string) (*entity.ExportResult, error)

	// ArchiveLogs 归档审计日志
	ArchiveLogs(ctx context.Context, beforeDate time.Time) (*entity.ArchiveStats, error)

	// DeleteLog 删除日志(管理员)
	DeleteLog(ctx context.Context, logID string) error

	// CreateExportTask 创建导出任务
	CreateExportTask(ctx context.Context, tenantID, userID string, filter *entity.LogFilter) (*entity.ExportResult, error)
}

// NSQProducer NSQ生产者接口
type NSQProducer interface {
	Publish(topic string, data []byte) error
}

// auditLogService 审计日志服务实现
type auditLogService struct {
	db            *gorm.DB
	repo          repository.AuditLogRepository
	maskingSvc    *service.DataMaskingService
	nsqProducer   NSQProducer
	enableAsync   bool
}

// NewAuditLogService 创建审计日志服务实例
func NewAuditLogService(
	db *gorm.DB,
	repo repository.AuditLogRepository,
	maskingSvc *service.DataMaskingService,
	nsqProducer NSQProducer,
	enableAsync bool,
) AuditLogService {
	return &auditLogService{
		db:          db,
		repo:        repo,
		maskingSvc:  maskingSvc,
		nsqProducer: nsqProducer,
		enableAsync: enableAsync,
	}
}

// LogOperation 记录操作日志
func (s *auditLogService) LogOperation(ctx context.Context, log *entity.AuditLog) error {
	// 1. 验证必填字段
	if err := s.validateLog(log); err != nil {
		return fmt.Errorf("invalid audit log: %w", err)
	}

	// 2. 脱敏敏感数据
	if err := s.maskSensitiveData(log); err != nil {
		return fmt.Errorf("failed to mask sensitive data: %w", err)
	}

	// 3. 生成数字签名(防篡改)
	signature, err := s.generateSignature(log)
	if err != nil {
		return fmt.Errorf("failed to generate signature: %w", err)
	}
	log.Signature = signature

	// 4. 设置创建时间
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	// 5. 保存到数据库
	return s.repo.Create(ctx, log)
}

// LogOperationAsync 异步记录操作日志
func (s *auditLogService) LogOperationAsync(ctx context.Context, log *entity.AuditLog) error {
	if !s.enableAsync || s.nsqProducer == nil {
		// 如果异步未启用，回退到同步记录
		return s.LogOperation(ctx, log)
	}

	// 1. 验证必填字段
	if err := s.validateLog(log); err != nil {
		return err
	}

	// 2. 脱敏敏感数据
	if err := s.maskSensitiveData(log); err != nil {
		return err
	}

	// 3. 序列化日志
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal audit log: %w", err)
	}

	// 4. 发布到NSQ
	topic := "audit_logs"
	if err := s.nsqProducer.Publish(topic, data); err != nil {
		// NSQ发布失败，回退到同步记录
		return s.LogOperation(ctx, log)
	}

	return nil
}

// QueryLogs 查询审计日志
func (s *auditLogService) QueryLogs(ctx context.Context, filter *entity.LogFilter) ([]*entity.AuditLog, int64, error) {
	return s.repo.Query(ctx, filter)
}

// GetLogByID 根据ID获取日志
func (s *auditLogService) GetLogByID(ctx context.Context, logID string) (*entity.AuditLog, error) {
	log, err := s.repo.GetByID(ctx, logID)
	if err != nil {
		return nil, err
	}

	// 验证签名完整性
	if err := s.verifySignature(log); err != nil {
		return nil, fmt.Errorf("log signature verification failed: %w", err)
	}

	return log, nil
}

// GetUserLogs 获取用户日志
func (s *auditLogService) GetUserLogs(ctx context.Context, tenantID, userID string, page, pageSize int) ([]*entity.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	logs, err := s.repo.GetByUserID(ctx, tenantID, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	filter := &entity.LogFilter{
		TenantID: tenantID,
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
	}
	_, total, err := s.repo.Query(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetResourceLogs 获取资源日志
func (s *auditLogService) GetResourceLogs(ctx context.Context, tenantID string, resource entity.AuditResource, resourceID string) ([]*entity.AuditLog, error) {
	return s.repo.GetByResourceID(ctx, tenantID, resource, resourceID)
}

// GetSensitiveOperations 获取敏感操作日志
func (s *auditLogService) GetSensitiveOperations(ctx context.Context, tenantID string, days int) ([]*entity.AuditLog, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)
	return s.repo.GetSensitiveOperations(ctx, tenantID, startTime, endTime)
}

// GetStatistics 获取审计统计
func (s *auditLogService) GetStatistics(ctx context.Context, tenantID string, days int) (*entity.AuditStats, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)
	return s.repo.GetStatistics(ctx, tenantID, startTime, endTime)
}

// ExportLogs 导出审计日志
func (s *auditLogService) ExportLogs(ctx context.Context, filter *entity.LogFilter, format string) (*entity.ExportResult, error) {
	// 1. 查询日志
	logs, _, err := s.repo.Query(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 2. 根据格式生成导出文件
	var data []byte
	var fileName string
	var contentType string

	switch format {
	case "csv":
		data, fileName, err = s.exportToCSV(logs)
		contentType = "text/csv"
	case "excel":
		data, fileName, err = s.exportToExcel(logs)
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "json":
		data, fileName, err = s.exportToJSON(logs)
		contentType = "application/json"
	default:
		return nil, entity.ErrInvalidExportFormat
	}

	if err != nil {
		return nil, err
	}

	// 3. 保存导出结果(这里简化处理，实际应该上传到对象存储)
	result := &entity.ExportResult{
		FileID:      fmt.Sprintf("audit_export_%s_%d", filter.TenantID, time.Now().Unix()),
		FileName:    fileName,
		FileSize:    int64(len(data)),
		RecordCount: len(logs),
		Format:      format,
		ContentType: contentType, // 设置内容类型，用于HTTP响应头
		DownloadURL: fmt.Sprintf("/api/audit/download/%s", fileName),
		ExpiresAt:   time.Now().Add(24 * time.Hour), // 24小时后过期
		CreatedAt:   time.Now(),
	}

	return result, nil
}

// ArchiveLogs 归档审计日志
func (s *auditLogService) ArchiveLogs(ctx context.Context, beforeDate time.Time) (*entity.ArchiveStats, error) {
	return s.repo.Archive(ctx, beforeDate)
}

// DeleteLog 删除日志
func (s *auditLogService) DeleteLog(ctx context.Context, logID string) error {
	return s.repo.Delete(ctx, logID)
}

// CreateExportTask 创建导出任务
func (s *auditLogService) CreateExportTask(ctx context.Context, tenantID, userID string, filter *entity.LogFilter) (*entity.ExportResult, error) {
	// 1. 验证过滤条件
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	// 2. 创建导出任务(异步处理)
	result := &entity.ExportResult{
		FileID:      fmt.Sprintf("audit_task_%s_%d", tenantID, time.Now().UnixNano()),
		FileName:    fmt.Sprintf("audit_logs_%s_%s.csv", tenantID, time.Now().Format("20060102_150405")),
		Format:      "csv",
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour), // 7天后过期
	}

	// 这里应该创建异步任务处理导出，简化处理直接返回
	return result, nil
}

// validateLog 验证日志
func (s *auditLogService) validateLog(log *entity.AuditLog) error {
	if log.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if log.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if log.UserEmail == "" {
		return fmt.Errorf("user_email is required")
	}
	if log.Action == "" {
		return fmt.Errorf("action is required")
	}
	if log.Resource == "" {
		return fmt.Errorf("resource is required")
	}
	return nil
}

// maskSensitiveData 脱敏敏感数据
func (s *auditLogService) maskSensitiveData(log *entity.AuditLog) error {
	// 脱敏请求数据
	if log.RequestData != "" {
		log.RequestData = s.maskingSvc.SanitizeForLog(log.RequestData)
	}

	// 脱敏响应数据
	if log.ResponseData != "" {
		log.ResponseData = s.maskingSvc.SanitizeForLog(log.ResponseData)
	}

	// 脱敏变更内容
	if log.Changes != "" {
		log.Changes = s.maskingSvc.SanitizeForLog(log.Changes)
	}

	// 脱敏邮箱
	log.UserEmail = s.maskingSvc.MaskEmail(log.UserEmail)

	return nil
}

// generateSignature 生成数字签名
func (s *auditLogService) generateSignature(log *entity.AuditLog) (string, error) {
	// 构建签名数据(排除signature字段本身)
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%d",
		log.TenantID,
		log.UserID,
		log.Action,
		log.Resource,
		log.ResourceID,
		log.RequestPath,
		log.Status,
		log.CreatedAt.UnixNano(),
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}

// verifySignature 验证数字签名
func (s *auditLogService) verifySignature(log *entity.AuditLog) error {
	expected, err := s.generateSignature(log)
	if err != nil {
		return err
	}

	if log.Signature != expected {
		return fmt.Errorf("signature mismatch: expected %s, got %s", expected, log.Signature)
	}

	return nil
}

// exportToCSV 导出到CSV
func (s *auditLogService) exportToCSV(logs []*entity.AuditLog) ([]byte, string, error) {
	// CSV header
	csv := "Time,User,Action,Resource,Resource Name,Status,Error,IP,Path\n"

	// CSV rows
	for _, log := range logs {
		csv += fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			log.CreatedAt.Format("2006-01-02 15:04:05"),
			log.UserEmail,
			log.Action,
			log.Resource,
			log.ResourceName,
			log.Status,
			log.ErrorMsg,
			log.RequestIP,
			log.RequestPath,
		)
	}

	fileName := fmt.Sprintf("audit_logs_%s.csv", time.Now().Format("20060102_150405"))
	return []byte(csv), fileName, nil
}

// exportToExcel 导出到Excel(简化实现，返回CSV)
func (s *auditLogService) exportToExcel(logs []*entity.AuditLog) ([]byte, string, error) {
	// 简化实现，实际应该使用excelize库
	return s.exportToCSV(logs)
}

// exportToJSON 导出到JSON
func (s *auditLogService) exportToJSON(logs []*entity.AuditLog) ([]byte, string, error) {
	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, "", err
	}

	fileName := fmt.Sprintf("audit_logs_%s.json", time.Now().Format("20060102_150405"))
	return data, fileName, nil
}
