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

package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/audit/entity"
	auditmodel "github.com/coze-dev/coze-studio/backend/domain/audit/internal/dal/model"
)

// AuditLogRepository 审计日志仓储接口
type AuditLogRepository interface {
	// Create 创建审计日志
	Create(ctx context.Context, log *entity.AuditLog) error

	// BatchCreate 批量创建审计日志
	BatchCreate(ctx context.Context, logs []*entity.AuditLog) error

	// GetByID 根据ID获取审计日志
	GetByID(ctx context.Context, logID string) (*entity.AuditLog, error)

	// Query 查询审计日志
	Query(ctx context.Context, filter *entity.LogFilter) ([]*entity.AuditLog, int64, error)

	// GetByUserID 根据用户ID获取日志
	GetByUserID(ctx context.Context, tenantID, userID string, limit int, offset int) ([]*entity.AuditLog, error)

	// GetByResourceID 根据资源ID获取日志
	GetByResourceID(ctx context.Context, tenantID string, resource entity.AuditResource, resourceID string) ([]*entity.AuditLog, error)

	// GetSensitiveOperations 获取敏感操作日志
	GetSensitiveOperations(ctx context.Context, tenantID string, startTime, endTime time.Time) ([]*entity.AuditLog, error)

	// GetStatistics 获取审计统计
	GetStatistics(ctx context.Context, tenantID string, startTime, endTime time.Time) (*entity.AuditStats, error)

	// Archive 归档日志
	Archive(ctx context.Context, beforeDate time.Time) (*entity.ArchiveStats, error)

	// Delete 删除日志(仅管理员)
	Delete(ctx context.Context, logID string) error

	// DeleteOldLogs 删除旧日志(定期清理任务)
	DeleteOldLogs(ctx context.Context, beforeDate time.Time) (int64, error)
}

// auditLogRepository 审计日志仓储实现
type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository 创建审计日志仓储实例
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Create 创建审计日志
func (r *auditLogRepository) Create(ctx context.Context, log *entity.AuditLog) error {
	do := r.entityToDO(log)
	if err := r.db.WithContext(ctx).Create(do).Error; err != nil {
		return err
	}
	return nil
}

// BatchCreate 批量创建审计日志
func (r *auditLogRepository) BatchCreate(ctx context.Context, logs []*entity.AuditLog) error {
	if len(logs) == 0 {
		return nil
	}

	dos := make([]*auditmodel.AuditLogDO, len(logs))
	for i, log := range logs {
		dos[i] = r.entityToDO(log)
	}

	return r.db.WithContext(ctx).CreateInBatches(dos, 100).Error
}

// GetByID 根据ID获取审计日志
func (r *auditLogRepository) GetByID(ctx context.Context, logID string) (*entity.AuditLog, error) {
	var do auditmodel.AuditLogDO
	err := r.db.WithContext(ctx).Where("log_id = ?", logID).First(&do).Error
	if err != nil {
		return nil, err
	}
	return r.doToEntity(&do), nil
}

// Query 查询审计日志
func (r *auditLogRepository) Query(ctx context.Context, filter *entity.LogFilter) ([]*entity.AuditLog, int64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&auditmodel.AuditLogDO{})

	// 应用过滤条件
	query = r.applyFilters(query, filter)

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderBY := filter.OrderBy
	if filter.OrderDesc {
		orderBY += " DESC"
	} else {
		orderBY += " ASC"
	}
	query = query.Order(orderBY)

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	// 查询
	var dos []*auditmodel.AuditLogDO
	if err := query.Find(&dos).Error; err != nil {
		return nil, 0, err
	}

	// 转换
	logs := make([]*entity.AuditLog, len(dos))
	for i, do := range dos {
		logs[i] = r.doToEntity(do)
	}

	return logs, total, nil
}

// GetByUserID 根据用户ID获取日志
func (r *auditLogRepository) GetByUserID(ctx context.Context, tenantID, userID string, limit int, offset int) ([]*entity.AuditLog, error) {
	var dos []*auditmodel.AuditLogDO
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&dos).Error
	if err != nil {
		return nil, err
	}

	logs := make([]*entity.AuditLog, len(dos))
	for i, do := range dos {
		logs[i] = r.doToEntity(do)
	}
	return logs, nil
}

// GetByResourceID 根据资源ID获取日志
func (r *auditLogRepository) GetByResourceID(ctx context.Context, tenantID string, resource entity.AuditResource, resourceID string) ([]*entity.AuditLog, error) {
	var dos []*auditmodel.AuditLogDO
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND resource = ? AND resource_id = ?", tenantID, resource, resourceID).
		Order("created_at DESC").
		Find(&dos).Error
	if err != nil {
		return nil, err
	}

	logs := make([]*entity.AuditLog, len(dos))
	for i, do := range dos {
		logs[i] = r.doToEntity(do)
	}
	return logs, nil
}

// GetSensitiveOperations 获取敏感操作日志
func (r *auditLogRepository) GetSensitiveOperations(ctx context.Context, tenantID string, startTime, endTime time.Time) ([]*entity.AuditLog, error) {
	sensitiveActions := []entity.AuditAction{
		entity.AuditActionUserDelete,
		entity.AuditActionRoleDelete,
		entity.AuditActionDataDelete,
		entity.AuditActionDataExport,
		entity.AuditActionUserPassword,
		entity.AuditActionConfigUpdate,
		entity.AuditActionConfigDelete,
		entity.AuditActionSystemRestore,
	}

	var dos []*auditmodel.AuditLogDO
	query := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("action IN ?", sensitiveActions)

	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	err := query.Order("created_at DESC").Find(&dos).Error
	if err != nil {
		return nil, err
	}

	logs := make([]*entity.AuditLog, len(dos))
	for i, do := range dos {
		logs[i] = r.doToEntity(do)
	}
	return logs, nil
}

// GetStatistics 获取审计统计
func (r *auditLogRepository) GetStatistics(ctx context.Context, tenantID string, startTime, endTime time.Time) (*entity.AuditStats, error) {
	stats := &entity.AuditStats{
		ActionStats:   make(map[string]int64),
		ResourceStats: make(map[string]int64),
		UserStats:     make(map[string]int64),
	}

	query := r.db.WithContext(ctx).Model(&auditmodel.AuditLogDO{}).Where("tenant_id = ?", tenantID)

	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	// 总数统计
	if err := query.Count(&stats.TotalLogs).Error; err != nil {
		return nil, err
	}

	// 成功/失败统计
	if err := query.Where("status = ?", entity.AuditStatusSuccess).Count(&stats.SuccessLogs).Error; err != nil {
		return nil, err
	}
	stats.FailedLogs = stats.TotalLogs - stats.SuccessLogs

	// 操作类型统计
	var actionStats []struct {
		Action string
		Count  int64
	}
	if err := query.Select("action, count(*) as count").
		Group("action").
		Scan(&actionStats).Error; err != nil {
		return nil, err
	}
	for _, stat := range actionStats {
		stats.ActionStats[stat.Action] = stat.Count
	}

	// 资源类型统计
	var resourceStats []struct {
		Resource string
		Count    int64
	}
	if err := query.Select("resource, count(*) as count").
		Group("resource").
		Scan(&resourceStats).Error; err != nil {
		return nil, err
	}
	for _, stat := range resourceStats {
		stats.ResourceStats[stat.Resource] = stat.Count
	}

	// 用户统计(Top 10)
	var userStats []struct {
		UserEmail string
		Count     int64
	}
	if err := query.Select("user_email, count(*) as count").
		Group("user_email").
		Order("count DESC").
		Limit(10).
		Scan(&userStats).Error; err != nil {
		return nil, err
	}
	for _, stat := range userStats {
		stats.UserStats[stat.UserEmail] = stat.Count
	}

	// 每日统计(最近30天)
	var dailyStats []struct {
		Date  string
		Count int64
	}
	if err := r.db.WithContext(ctx).
		Select("DATE(created_at) as date, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyStats).Error; err != nil {
		return nil, err
	}
	stats.DailyStats = make([]entity.DailyStat, len(dailyStats))
	for i, stat := range dailyStats {
		stats.DailyStats[i] = entity.DailyStat{
			Date:  stat.Date,
			Count: stat.Count,
		}
	}

	return stats, nil
}

// Archive 归档日志
func (r *auditLogRepository) Archive(ctx context.Context, beforeDate time.Time) (*entity.ArchiveStats, error) {
	stats := &entity.ArchiveStats{
		ArchiveDate: time.Now(),
	}

	// 统计需要归档的日志数
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&auditmodel.AuditLogDO{}).
		Where("created_at < ? AND is_archived = ?", beforeDate, false).
		Count(&count).Error; err != nil {
		return nil, err
	}
	stats.TotalLogs = int(count)

	// 标记为已归档(软删除，不物理删除)
	result := r.db.WithContext(ctx).
		Model(&auditmodel.AuditLogDO{}).
		Where("created_at < ? AND is_archived = ?", beforeDate, false).
		Updates(map[string]interface{}{
			"is_archived": true,
			"archived_at": time.Now(),
		})

	if err := result.Error; err != nil {
		return nil, err
	}

	stats.ArchivedLogs = int(result.RowsAffected)
	return stats, nil
}

// Delete 删除日志
func (r *auditLogRepository) Delete(ctx context.Context, logID string) error {
	return r.db.WithContext(ctx).
		Where("log_id = ?", logID).
		Delete(&auditmodel.AuditLogDO{}).Error
}

// DeleteOldLogs 删除旧日志(物理删除，仅限管理员操作)
func (r *auditLogRepository) DeleteOldLogs(ctx context.Context, beforeDate time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("created_at < ?", beforeDate).
		Delete(&auditmodel.AuditLogDO{})
	if err := result.Error; err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

// applyFilters 应用过滤条件
func (r *auditLogRepository) applyFilters(query *gorm.DB, filter *entity.LogFilter) *gorm.DB {
	query = query.Where("tenant_id = ?", filter.TenantID)

	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.UserEmail != "" {
		query = query.Where("user_email LIKE ?", "%"+filter.UserEmail+"%")
	}
	if len(filter.Actions) > 0 {
		query = query.Where("action IN ?", filter.Actions)
	}
	if len(filter.Resources) > 0 {
		query = query.Where("resource IN ?", filter.Resources)
	}
	if filter.ResourceID != "" {
		query = query.Where("resource_id = ?", filter.ResourceID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.RequestIP != "" {
		query = query.Where("request_ip = ?", filter.RequestIP)
	}
	if filter.SessionID != "" {
		query = query.Where("session_id = ?", filter.SessionID)
	}
	if filter.TraceID != "" {
		query = query.Where("trace_id = ?", filter.TraceID)
	}
	if !filter.StartTime.IsZero() {
		query = query.Where("created_at >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		query = query.Where("created_at <= ?", filter.EndTime)
	}

	return query
}

// entityToDO 实体转DO
func (r *auditLogRepository) entityToDO(log *entity.AuditLog) *auditmodel.AuditLogDO {
	return &auditmodel.AuditLogDO{
		LogID:         log.LogID,
		TenantID:      log.TenantID,
		UserID:        log.UserID,
		Username:      log.Username,
		UserEmail:     log.UserEmail,
		Action:        string(log.Action),
		Resource:      string(log.Resource),
		ResourceID:    log.ResourceID,
		ResourceName:  log.ResourceName,
		RequestMethod: log.RequestMethod,
		RequestPath:   log.RequestPath,
		RequestIP:     log.RequestIP,
		UserAgent:     log.UserAgent,
		Status:        string(log.Status),
		ErrorCode:     log.ErrorCode,
		ErrorMsg:      log.ErrorMsg,
		RequestData:   log.RequestData,
		ResponseData:  log.ResponseData,
		Changes:       log.Changes,
		Metadata:      log.Metadata,
		SessionID:     log.SessionID,
		TraceID:       log.TraceID,
		CreatedAt:     log.CreatedAt,
		ArchivedAt:    log.ArchivedAt,
		Signature:     log.Signature,
		IsArchived:    log.IsArchived,
	}
}

// doToEntity DO转实体
func (r *auditLogRepository) doToEntity(do *auditmodel.AuditLogDO) *entity.AuditLog {
	return &entity.AuditLog{
		LogID:         do.LogID,
		TenantID:      do.TenantID,
		UserID:        do.UserID,
		Username:      do.Username,
		UserEmail:     do.UserEmail,
		Action:        entity.AuditAction(do.Action),
		Resource:      entity.AuditResource(do.Resource),
		ResourceID:    do.ResourceID,
		ResourceName:  do.ResourceName,
		RequestMethod: do.RequestMethod,
		RequestPath:   do.RequestPath,
		RequestIP:     do.RequestIP,
		UserAgent:     do.UserAgent,
		Status:        entity.AuditStatus(do.Status),
		ErrorCode:     do.ErrorCode,
		ErrorMsg:      do.ErrorMsg,
		RequestData:   do.RequestData,
		ResponseData:  do.ResponseData,
		Changes:       do.Changes,
		Metadata:      do.Metadata,
		SessionID:     do.SessionID,
		TraceID:       do.TraceID,
		CreatedAt:     do.CreatedAt,
		ArchivedAt:    do.ArchivedAt,
		Signature:     do.Signature,
		IsArchived:    do.IsArchived,
	}
}
