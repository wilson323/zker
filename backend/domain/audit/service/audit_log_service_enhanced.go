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
	"encoding/json"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/audit/entity"
	auditrepo "github.com/coze-dev/coze-studio/backend/domain/audit/repository"
	"github.com/coze-dev/coze-studio/backend/domain/security/service"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// AuditLogServiceEnhanced 增强的审计日志服务
type AuditLogServiceEnhanced struct {
	db             *gorm.DB
	repo           auditrepo.AuditLogRepository
	maskingSvc     *service.DataMaskingService
	esClient       *elastic.Client
	logger         *zap.Logger
	indexPrefix    string
	enableES       bool
}

// ElasticsearchIndexer Elasticsearch索引器接口
type ElasticsearchIndexer interface {
	Index(index string, data interface{}) (*elastic.IndexResponse, error)
	BulkIndex(index string, items []elastic.BulkIndexRequest) (*elastic.BulkResponse, error)
	Search(index string, query elastic.Query) (*elastic.SearchResult, error)
	Delete(index string, id string) (*elastic.DeleteResponse, error)
}

// NewAuditLogServiceEnhanced 创建增强的审计日志服务实例
func NewAuditLogServiceEnhanced(
	db *gorm.DB,
	repo auditrepo.AuditLogRepository,
	maskingSvc *service.DataMaskingService,
	esClient *elastic.Client,
	logger *zap.Logger,
	enableES bool,
) *AuditLogServiceEnhanced {
	return &AuditLogServiceEnhanced{
		db:          db,
		repo:        repo,
		maskingSvc:  maskingSvc,
		esClient:    esClient,
		logger:      logger,
		indexPrefix: "zker_audit_logs",
		enableES:    enableES,
	}
}

// LogAction 记录审计日志
func (s *AuditLogServiceEnhanced) LogAction(ctx context.Context, action *AuditAction) error {
	s.logger.Info("Logging audit action",
		zap.String("tenant_id", action.TenantID),
		zap.String("action", string(action.Action)),
		zap.String("resource", string(action.Resource)),
	)

	// 构建审计日志实体
	log := &entity.AuditLog{
		LogID:         generateLogID(),
		TenantID:      action.TenantID,
		UserID:        action.UserID,
		Username:      action.Username,
		UserEmail:     action.UserEmail,
		Action:        action.Action,
		Resource:      action.Resource,
		ResourceID:    action.ResourceID,
		ResourceName:  action.ResourceName,
		RequestMethod: action.RequestMethod,
		RequestPath:   action.RequestPath,
		RequestIP:     action.RequestIP,
		UserAgent:     action.UserAgent,
		Status:        entity.AuditStatusSuccess,
		ErrorCode:     action.ErrorCode,
		ErrorMsg:      action.ErrorMessage,
		SessionID:     action.SessionID,
		TraceID:       action.TraceID,
		CreatedAt:     time.Now(),
	}

	// 设置状态
	if action.ErrorCode != "" {
		log.Status = entity.AuditStatusFailed
	}

	// 序列化请求和响应数据
	if action.RequestData != nil {
		data, err := json.Marshal(action.RequestData)
		if err != nil {
			s.logger.Warn("Failed to marshal request data", zap.Error(err))
		} else {
			log.RequestData = string(data)
		}
	}

	if action.ResponseData != nil {
		data, err := json.Marshal(action.ResponseData)
		if err != nil {
			s.logger.Warn("Failed to marshal response data", zap.Error(err))
		} else {
			log.ResponseData = string(data)
		}
	}

	// 序列化变更内容
	if action.Changes != nil {
		data, err := json.Marshal(action.Changes)
		if err != nil {
			s.logger.Warn("Failed to marshal changes", zap.Error(err))
		} else {
			log.Changes = string(data)
		}
	}

	// 序列化元数据
	if action.Metadata != nil {
		data, err := json.Marshal(action.Metadata)
		if err != nil {
			s.logger.Warn("Failed to marshal metadata", zap.Error(err))
		} else {
			log.Metadata = string(data)
		}
	}

	// 脱敏敏感数据
	if err := s.maskSensitiveData(log); err != nil {
		return errorx.Wrapf(err, errno.DataMaskingFailed)
	}

	// 生成数字签名
	signature, err := s.generateSignature(log)
	if err != nil {
		return errorx.Wrapf(err, errno.SignatureGenerationFailed)
	}
	log.Signature = signature

	// 保存到数据库
	if err := s.repo.Create(ctx, log); err != nil {
		s.logger.Error("Failed to save audit log to database", zap.Error(err))
		return errorx.Wrapf(err, errno.AuditLogSaveFailed)
	}

	// 异步索引到Elasticsearch
	if s.enableES && s.esClient != nil {
		go s.indexToElasticsearch(log)
	}

	return nil
}

// QueryLogs 查询审计日志
func (s *AuditLogServiceEnhanced) QueryLogs(ctx context.Context, req *QueryLogsRequest) (*AuditLogsResponse, error) {
	s.logger.Info("Querying audit logs",
		zap.String("tenant_id", req.TenantID),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize),
	)

	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, errorx.Wrapf(err, errno.InvalidRequest)
	}

	// 构建过滤器
	filter := &entity.LogFilter{
		TenantID:      req.TenantID,
		UserID:        req.UserID,
		UserEmail:     req.UserEmail,
		Actions:       req.Actions,
		Resources:     req.Resources,
		ResourceID:    req.ResourceID,
		Status:        req.Status,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		RequestIP:     req.RequestIP,
		SessionID:     req.SessionID,
		TraceID:       req.TraceID,
		Page:          req.Page,
		PageSize:      req.PageSize,
		OrderBy:       req.OrderBy,
		OrderDesc:     req.OrderDesc,
	}

	// 查询
	logs, total, err := s.repo.Query(ctx, filter)
	if err != nil {
		return nil, errorx.Wrapf(err, errno.AuditLogQueryFailed)
	}

	response := &AuditLogsResponse{
		Logs:       logs,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: (total + int64(req.PageSize) - 1) / int64(req.PageSize),
	}

	return response, nil
}

// QueryLogsByUser 按用户查询
func (s *AuditLogServiceEnhanced) QueryLogsByUser(
	ctx context.Context,
	userID string,
	timeRange *TimeRange,
	page, pageSize int,
) ([]*entity.AuditLog, error) {
	s.logger.Info("Querying audit logs by user",
		zap.String("user_id", userID),
		zap.Time("start", timeRange.StartTime),
		zap.Time("end", timeRange.EndTime),
	)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 20
	}

	filter := &entity.LogFilter{
		UserID:    userID,
		StartTime: &timeRange.StartTime,
		EndTime:   &timeRange.EndTime,
		Page:      page,
		PageSize:  pageSize,
		OrderBy:   "created_at",
		OrderDesc: true,
	}

	logs, _, err := s.repo.Query(ctx, filter)
	if err != nil {
		return nil, errorx.Wrapf(err, errno.AuditLogQueryFailed)
	}

	return logs, nil
}

// QueryLogsByResource 按资源查询
func (s *AuditLogServiceEnhanced) QueryLogsByResource(
	ctx context.Context,
	resourceType string,
	resourceID string,
	page, pageSize int,
) ([]*entity.AuditLog, error) {
	s.logger.Info("Querying audit logs by resource",
		zap.String("resource_type", resourceType),
		zap.String("resource_id", resourceID),
	)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 20
	}

	filter := &entity.LogFilter{
		Resource:   entity.AuditResource(resourceType),
		ResourceID: resourceID,
		Page:       page,
		PageSize:   pageSize,
		OrderBy:    "created_at",
		OrderDesc:  true,
	}

	logs, _, err := s.repo.Query(ctx, filter)
	if err != nil {
		return nil, errorx.Wrapf(err, errno.AuditLogQueryFailed)
	}

	return logs, nil
}

// QueryLogsByTimeRange 按时间范围查询
func (s *AuditLogServiceEnhanced) QueryLogsByTimeRange(
	ctx context.Context,
	start, end time.Time,
	page, pageSize int,
) ([]*entity.AuditLog, error) {
	s.logger.Info("Querying audit logs by time range",
		zap.Time("start", start),
		zap.Time("end", end),
	)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 20
	}

	filter := &entity.LogFilter{
		StartTime: &start,
		EndTime:   &end,
		Page:      page,
		PageSize:  pageSize,
		OrderBy:   "created_at",
		OrderDesc: true,
	}

	logs, _, err := s.repo.Query(ctx, filter)
	if err != nil {
		return nil, errorx.Wrapf(err, errno.AuditLogQueryFailed)
	}

	return logs, nil
}

// ExportLogs 导出审计日志
func (s *AuditLogServiceEnhanced) ExportLogs(ctx context.Context, req *ExportRequest) (*ExportResponse, error) {
	s.logger.Info("Exporting audit logs",
		zap.String("tenant_id", req.TenantID),
		zap.String("format", req.Format),
	)

	// 构建过滤器
	filter := &entity.LogFilter{
		TenantID:  req.TenantID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Page:      1,
		PageSize:  req.MaxRows,
	}

	// 查询日志
	logs, total, err := s.repo.Query(ctx, filter)
	if err != nil {
		return nil, errorx.Wrapf(err, errno.AuditLogQueryFailed)
	}

	// 检查导出行数限制
	if total > int64(req.MaxRows) {
		return nil, errorx.New(errno.ExportLimitExceeded, fmt.Sprintf("Export limit exceeded: %d > %d", total, req.MaxRows))
	}

	// 根据格式生成导出文件
	var data []byte
	var fileName string
	var contentType string

	switch req.Format {
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
		return nil, errorx.New(errno.InvalidExportFormat, "Unsupported export format: "+req.Format)
	}

	if err != nil {
		return nil, err
	}

	// 生成导出结果
	result := &ExportResponse{
		FileID:       generateFileID(),
		FileName:     fileName,
		FileSize:     int64(len(data)),
		RecordCount:  len(logs),
		Format:       req.Format,
		ContentType:  contentType,
		Data:         data,
		DownloadURL:  fmt.Sprintf("/api/audit/download/%s", result.FileID),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	return result, nil
}

// GetAuditStatistics 获取审计统计
func (s *AuditLogServiceEnhanced) GetAuditStatistics(
	ctx context.Context,
	tenantID string,
	timeRange *TimeRange,
) (*AuditStatistics, error) {
	s.logger.Info("Getting audit statistics",
		zap.String("tenant_id", tenantID),
	)

	stats, err := s.repo.GetStatistics(ctx, tenantID, timeRange.StartTime, timeRange.EndTime)
	if err != nil {
		return nil, errorx.Wrapf(err, errno.AuditStatisticsQueryFailed)
	}

	return &AuditStatistics{
		TenantID:        tenantID,
		TimeRange:       *timeRange,
		TotalLogs:       stats.TotalLogs,
		SuccessLogs:     stats.SuccessLogs,
		FailedLogs:      stats.FailedLogs,
		ActionStats:     stats.ActionStats,
		ResourceStats:   stats.ResourceStats,
		UserStats:       stats.UserStats,
		DailyStats:      stats.DailyStats,
		GeneratedAt:     time.Now(),
	}, nil
}

// QueryLogsFullText 全文搜索审计日志(使用Elasticsearch)
func (s *AuditLogServiceEnhanced) QueryLogsFullText(
	ctx context.Context,
	req *FullTextSearchRequest,
) (*AuditLogsResponse, error) {
	if !s.enableES || s.esClient == nil {
		return nil, errorx.New(errno.ElasticsearchNotAvailable, "Elasticsearch is not enabled")
	}

	s.logger.Info("Full text searching audit logs",
		zap.String("query", req.Query),
		zap.String("tenant_id", req.TenantID),
	)

	// 构建Elasticsearch查询
	boolQuery := elastic.NewBoolQuery()

	// 租户过滤
	boolQuery.Must(elastic.NewTermQuery("tenant_id", req.TenantID))

	// 全文搜索
	if req.Query != "" {
		multiMatchQuery := elastic.NewMultiMatchQuery(req.Query, "username", "user_email", "resource_name", "error_msg")
		boolQuery.Must(multiMatchQuery)
	}

	// 时间范围过滤
	if req.StartTime != nil && req.EndTime != nil {
		rangeQuery := elastic.NewRangeQuery("created_at").
			Gte(req.StartTime.Format(time.RFC3339)).
			Lte(req.EndTime.Format(time.RFC3339))
		boolQuery.Must(rangeQuery)
	}

	// 执行搜索
	searchResult, err := s.esClient.Search().
		Index(s.getIndexName(req.TenantID)).
		Query(boolQuery).
		From((req.Page - 1) * req.PageSize).
		Size(req.PageSize).
		Sort("created_at", false).
		Do(ctx)

	if err != nil {
		s.logger.Error("Elasticsearch search failed", zap.Error(err))
		return nil, errorx.Wrapf(err, errno.AuditLogQueryFailed)
	}

	// 解析结果
	logs := make([]*entity.AuditLog, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		var log entity.AuditLog
		if err := json.Unmarshal(hit.Source, &log); err != nil {
			s.logger.Warn("Failed to unmarshal audit log", zap.Error(err))
			continue
		}
		logs = append(logs, &log)
	}

	response := &AuditLogsResponse{
		Logs:       logs,
		Total:      searchResult.TotalHits(),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: int((searchResult.TotalHits() + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}

	return response, nil
}

// ============================================================
// 私有方法
// ============================================================

// maskSensitiveData 脱敏敏感数据
func (s *AuditLogServiceEnhanced) maskSensitiveData(log *entity.AuditLog) error {
	if log.RequestData != "" {
		log.RequestData = s.maskingSvc.SanitizeForLog(log.RequestData)
	}
	if log.ResponseData != "" {
		log.ResponseData = s.maskingSvc.SanitizeForLog(log.ResponseData)
	}
	if log.Changes != "" {
		log.Changes = s.maskingSvc.SanitizeForLog(log.Changes)
	}
	log.UserEmail = s.maskingSvc.MaskEmail(log.UserEmail)
	return nil
}

// generateSignature 生成数字签名
func (s *AuditLogServiceEnhanced) generateSignature(log *entity.AuditLog) (string, error) {
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

// indexToElasticsearch 索引到Elasticsearch
func (s *AuditLogServiceEnhanced) indexToElasticsearch(log *entity.AuditLog) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexName := s.getIndexName(log.TenantID)

	_, err := s.esClient.Index().
		Index(indexName).
		Id(log.LogID).
		BodyJson(log).
		Do(ctx)

	if err != nil {
		s.logger.Error("Failed to index audit log to Elasticsearch",
			zap.String("log_id", log.LogID),
			zap.Error(err),
		)
	}
}

// getIndexName 获取索引名称
func (s *AuditLogServiceEnhanced) getIndexName(tenantID string) string {
	date := time.Now().Format("2006-01")
	return fmt.Sprintf("%s-%s-%s", s.indexPrefix, tenantID, date)
}

// exportToCSV 导出到CSV
func (s *AuditLogServiceEnhanced) exportToCSV(logs []*entity.AuditLog) ([]byte, string, error) {
	csv := "Time,User,Action,Resource,Resource Name,Status,Error,IP,Path\n"

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

// exportToExcel 导出到Excel
func (s *AuditLogServiceEnhanced) exportToExcel(logs []*entity.AuditLog) ([]byte, string, error) {
	// 简化实现，返回CSV
	return s.exportToCSV(logs)
}

// exportToJSON 导出到JSON
func (s *AuditLogServiceEnhanced) exportToJSON(logs []*entity.AuditLog) ([]byte, string, error) {
	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, "", err
	}

	fileName := fmt.Sprintf("audit_logs_%s.json", time.Now().Format("20060102_150405"))
	return data, fileName, nil
}

// ============================================================
// 辅助函数
// ============================================================

func generateLogID() string {
	return fmt.Sprintf("audit_%d", time.Now().UnixNano())
}

func generateFileID() string {
	return fmt.Sprintf("export_%d", time.Now().UnixNano())
}

// ============================================================
// 数据类型定义
// ============================================================

// AuditAction 审计操作
type AuditAction struct {
	TenantID      string                 `json:"tenant_id"`
	UserID        string                 `json:"user_id"`
	Username      string                 `json:"username"`
	UserEmail     string                 `json:"user_email"`
	Action        entity.AuditAction     `json:"action"`
	Resource      entity.AuditResource   `json:"resource"`
	ResourceID    string                 `json:"resource_id"`
	ResourceName  string                 `json:"resource_name"`
	RequestMethod string                 `json:"request_method"`
	RequestPath   string                 `json:"request_path"`
	RequestIP     string                 `json:"request_ip"`
	UserAgent     string                 `json:"user_agent"`
	ErrorCode     string                 `json:"error_code,omitempty"`
	ErrorMessage  string                 `json:"error_message,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	TraceID       string                 `json:"trace_id,omitempty"`
	RequestData   map[string]interface{} `json:"request_data,omitempty"`
	ResponseData  map[string]interface{} `json:"response_data,omitempty"`
	Changes       map[string]interface{} `json:"changes,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// QueryLogsRequest 查询审计日志请求
type QueryLogsRequest struct {
	TenantID    string                `json:"tenant_id" binding:"required"`
	UserID      string                `json:"user_id,omitempty"`
	UserEmail   string                `json:"user_email,omitempty"`
	Actions     []entity.AuditAction  `json:"actions,omitempty"`
	Resources   []entity.AuditResource `json:"resources,omitempty"`
	ResourceID  string                `json:"resource_id,omitempty"`
	Status      entity.AuditStatus    `json:"status,omitempty"`
	StartTime   *time.Time            `json:"start_time,omitempty"`
	EndTime     *time.Time            `json:"end_time,omitempty"`
	RequestIP   string                `json:"request_ip,omitempty"`
	SessionID   string                `json:"session_id,omitempty"`
	TraceID     string                `json:"trace_id,omitempty"`
	Page        int                   `json:"page" binding:"min=1"`
	PageSize    int                   `json:"page_size" binding:"min=1,max=1000"`
	OrderBy     string                `json:"order_by,omitempty"`
	OrderDesc   bool                  `json:"order_desc,omitempty"`
}

// Validate 验证请求
func (r *QueryLogsRequest) Validate() error {
	if r.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if r.Page < 1 {
		r.Page = 1
	}
	if r.PageSize < 1 || r.PageSize > 1000 {
		r.PageSize = 20
	}
	if r.OrderBy == "" {
		r.OrderBy = "created_at"
		r.OrderDesc = true
	}
	return nil
}

// AuditLogsResponse 审计日志响应
type AuditLogsResponse struct {
	Logs       []*entity.AuditLog `json:"logs"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// TimeRange 时间范围
type TimeRange struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// ExportRequest 导出请求
type ExportRequest struct {
	TenantID  string     `json:"tenant_id" binding:"required"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Format    string     `json:"format" binding:"required,oneof=csv excel json"`
	MaxRows   int        `json:"max_rows" binding:"min=1,max=100000"`
}

// ExportResponse 导出响应
type ExportResponse struct {
	FileID       string    `json:"file_id"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	RecordCount  int       `json:"record_count"`
	Format       string    `json:"format"`
	ContentType  string    `json:"content_type"`
	Data         []byte    `json:"data,omitempty"`
	DownloadURL  string    `json:"download_url"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuditStatistics 审计统计
type AuditStatistics struct {
	TenantID      string              `json:"tenant_id"`
	TimeRange     TimeRange           `json:"time_range"`
	TotalLogs     int64               `json:"total_logs"`
	SuccessLogs   int64               `json:"success_logs"`
	FailedLogs    int64               `json:"failed_logs"`
	ActionStats   map[string]int64    `json:"action_stats"`
	ResourceStats map[string]int64    `json:"resource_stats"`
	UserStats     map[string]int64    `json:"user_stats"`
	DailyStats    []entity.DailyStat  `json:"daily_stats"`
	GeneratedAt   time.Time           `json:"generated_at"`
}

// FullTextSearchRequest 全文搜索请求
type FullTextSearchRequest struct {
	TenantID   string    `json:"tenant_id" binding:"required"`
	Query      string    `json:"query" binding:"required"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	Page       int       `json:"page" binding:"min=1"`
	PageSize   int       `json:"page_size" binding:"min=1,max=1000"`
}

import (
	"crypto/sha256"
	"encoding/hex"
	"gorm.io/gorm"
)
