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

package monitoring

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	monitoringEntity "github.com/coze-dev/coze-studio/backend/domain/monitoring/entity"
	monitoringService "github.com/coze-dev/coze-studio/backend/domain/monitoring/service"
)

// MonitoringHandler 监控API处理器
type MonitoringHandler struct {
	tenantMonitoringSvc *monitoringService.TenantMonitoringService
	alertSvc            *monitoringService.AlertService
	logger              *logrus.Logger
}

// NewMonitoringHandler 创建监控处理器
func NewMonitoringHandler(
	tenantMonitoringSvc *monitoringService.TenantMonitoringService,
	alertSvc *monitoringService.AlertService,
	logger *logrus.Logger,
) *MonitoringHandler {
	return &MonitoringHandler{
		tenantMonitoringSvc: tenantMonitoringSvc,
		alertSvc:            alertSvc,
		logger:              logger,
	}
}

// ============================================================
// 统一响应结构
// ============================================================

// APIResponse 统一API响应
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ============================================================
// 租户监控相关API
// ============================================================

// GetTenantOverview 获取租户概览
// GET /api/v1/monitoring/tenants/:tenant_id/overview
func (h *MonitoringHandler) GetTenantOverview(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Param("tenant_id")
	tenantName := c.Query("tenant_name")
	if tenantName == "" {
		tenantName = tenantID
	}

	overview, err := h.tenantMonitoringSvc.GetTenantOverview(ctx, tenantID, tenantName)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get tenant overview")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to get tenant overview",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Success",
		Data:    overview,
	})
}

// GetQPSMetrics 获取QPS指标
// GET /api/v1/monitoring/tenants/:tenant_id/metrics/qps
func (h *MonitoringHandler) GetQPSMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Param("tenant_id")

	// 解析时间范围
	timeRange, err := h.parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	metrics, err := h.tenantMonitoringSvc.CollectQPS(ctx, tenantID, timeRange)
	if err != nil {
		h.logger.WithError(err).Error("Failed to collect QPS metrics")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to collect QPS metrics",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Success",
		Data:    metrics,
	})
}

// GetResponseTimeMetrics 获取响应时间指标
// GET /api/v1/monitoring/tenants/:tenant_id/metrics/response_time
func (h *MonitoringHandler) GetResponseTimeMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Param("tenant_id")

	timeRange, err := h.parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	metrics, err := h.tenantMonitoringSvc.CollectResponseTime(ctx, tenantID, timeRange)
	if err != nil {
		h.logger.WithError(err).Error("Failed to collect response time metrics")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to collect response time metrics",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Success",
		Data:    metrics,
	})
}

// GetErrorRateMetrics 获取错误率指标
// GET /api/v1/monitoring/tenants/:tenant_id/metrics/error_rate
func (h *MonitoringHandler) GetErrorRateMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Param("tenant_id")

	timeRange, err := h.parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	metrics, err := h.tenantMonitoringSvc.CollectErrorRate(ctx, tenantID, timeRange)
	if err != nil {
		h.logger.WithError(err).Error("Failed to collect error rate metrics")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to collect error rate metrics",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Success",
		Data:    metrics,
	})
}

// GetHealthScore 获取健康评分
// GET /api/v1/monitoring/tenants/:tenant_id/health_score
func (h *MonitoringHandler) GetHealthScore(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Param("tenant_id")

	healthScore, err := h.tenantMonitoringSvc.CalculateTenantHealthScore(ctx, tenantID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to calculate health score")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to calculate health score",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Success",
		Data:    healthScore,
	})
}

// ============================================================
// 告警相关API
// ============================================================

// CreateAlertRuleRequest 创建告警规则请求
type CreateAlertRuleRequest struct {
	TenantID             string                          `json:"tenant_id" binding:"required"`
	RuleName             string                          `json:"rule_name" binding:"required"`
	Description          string                          `json:"description"`
	MetricType           string                          `json:"metric_type" binding:"required"`
	Threshold            float64                         `json:"threshold" binding:"required"`
	Comparison           monitoringEntity.AlertComparisonOperator `json:"comparison" binding:"required"`
	Severity             monitoringEntity.AlertSeverity  `json:"severity" binding:"required"`
	NotificationChannels monitoringEntity.NotificationChannels `json:"notification_channels"`
	NotificationConfig   string                          `json:"notification_config"`
	EvaluationInterval   int                             `json:"evaluation_interval"`
	ForDuration          int                             `json:"for_duration"`
}

// CreateAlertRule 创建告警规则
// POST /api/v1/monitoring/alert_rules
func (h *MonitoringHandler) CreateAlertRule(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateAlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	rule := monitoringEntity.NewAlertRule(
		req.TenantID,
		req.RuleName,
		req.MetricType,
		req.Threshold,
		req.Comparison,
		req.Severity,
	)

	rule.Description = req.Description
	rule.NotificationChannels = req.NotificationChannels
	rule.NotificationConfig = req.NotificationConfig
	if req.EvaluationInterval > 0 {
		rule.EvaluationInterval = req.EvaluationInterval
	}
	if req.ForDuration > 0 {
		rule.ForDuration = req.ForDuration
	}

	if err := h.alertSvc.CreateAlertRule(ctx, rule); err != nil {
		h.logger.WithError(err).Error("Failed to create alert rule")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to create alert rule",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    201,
		Message: "Alert rule created successfully",
		Data:    rule,
	})
}

// GetAlertHistoryRequest 获取告警历史请求
type GetAlertHistoryRequest struct {
	TenantID  string                     `json:"tenant_id"`
	Severity  monitoringEntity.AlertSeverity `json:"severity,omitempty"`
	Status    monitoringEntity.AlertStatus   `json:"status,omitempty"`
	StartTime string                     `json:"start_time,omitempty"`
	EndTime   string                     `json:"end_time,omitempty"`
	PageToken string                     `json:"page_token,omitempty"`
	PageSize  int                        `json:"page_size,omitempty"`
}

// GetAlertHistory 获取告警历史
// GET /api/v1/monitoring/alerts/history
func (h *MonitoringHandler) GetAlertHistory(c *gin.Context) {
	ctx := c.Request.Context()

	filter := &monitoringEntity.AlertFilter{
		TenantID:  c.Query("tenant_id"),
		PageToken: c.Query("page_token"),
		PageSize:  h.parseInt(c.Query("page_size"), 20),
	}

	if severity := c.Query("severity"); severity != "" {
		filter.Severity = monitoringEntity.AlertSeverity(severity)
	}

	if status := c.Query("status"); status != "" {
		filter.Status = monitoringEntity.AlertStatus(status)
	}

	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filter.StartTime = &t
		}
	}

	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			filter.EndTime = &t
		}
	}

	history, total, err := h.alertSvc.GetAlertHistory(ctx, filter)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get alert history")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to get alert history",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Success",
		Data: gin.H{
			"alerts": history,
			"total":  total,
		},
	})
}

// AcknowledgeAlertRequest 确认告警请求
type AcknowledgeAlertRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// AcknowledgeAlert 确认告警
// POST /api/v1/monitoring/alerts/:alert_id/acknowledge
func (h *MonitoringHandler) AcknowledgeAlert(c *gin.Context) {
	ctx := c.Request.Context()
	alertID := c.Param("alert_id")

	var req AcknowledgeAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	if err := h.alertSvc.AcknowledgeAlert(ctx, alertID, req.UserID); err != nil {
		h.logger.WithError(err).Error("Failed to acknowledge alert")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to acknowledge alert",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Alert acknowledged successfully",
	})
}

// ResolveAlertRequest 解决告警请求
type ResolveAlertRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Note   string `json:"note"`
}

// ResolveAlert 解决告警
// POST /api/v1/monitoring/alerts/:alert_id/resolve
func (h *MonitoringHandler) ResolveAlert(c *gin.Context) {
	ctx := c.Request.Context()
	alertID := c.Param("alert_id")

	var req ResolveAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	if err := h.alertSvc.ResolveAlert(ctx, alertID, req.UserID, req.Note); err != nil {
		h.logger.WithError(err).Error("Failed to resolve alert")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to resolve alert",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Alert resolved successfully",
	})
}

// SilenceAlertRequest 静默告警请求
type SilenceAlertRequest struct {
	DurationMinutes int `json:"duration_minutes" binding:"required,min=1"`
}

// SilenceAlert 静默告警
// POST /api/v1/monitoring/alerts/:alert_id/silence
func (h *MonitoringHandler) SilenceAlert(c *gin.Context) {
	ctx := c.Request.Context()
	alertID := c.Param("alert_id")

	var req SilenceAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	duration := time.Duration(req.DurationMinutes) * time.Minute
	if err := h.alertSvc.SilenceAlert(ctx, alertID, duration); err != nil {
		h.logger.WithError(err).Error("Failed to silence alert")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to silence alert",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Alert silenced successfully",
	})
}

// GetAlertStatistics 获取告警统计
// GET /api/v1/monitoring/alerts/statistics
func (h *MonitoringHandler) GetAlertStatistics(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Query("tenant_id")

	// 默认查询最近24小时
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	if startStr := c.Query("start_time"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = t
		}
	}

	if endStr := c.Query("end_time"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			endTime = t
		}
	}

	stats, err := h.alertSvc.GetAlertStatistics(ctx, tenantID, startTime, endTime)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get alert statistics")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to get alert statistics",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Success",
		Data:    stats,
	})
}

// ============================================================
// 辅助方法
// ============================================================

// parseTimeRange 解析时间范围
func (h *MonitoringHandler) parseTimeRange(c *gin.Context) (*monitoringEntity.TimeRange, error) {
	// 支持预设时间范围：5m, 15m, 1h, 6h, 24h, 7d, 30d
	if preset := c.Query("preset"); preset != "" {
		defaultRanges := monitoringEntity.GetDefaultTimeRanges()
		if duration, ok := defaultRanges[preset]; ok {
			return monitoringEntity.NewTimeRangeFromDuration(duration), nil
		}
	}

	// 支持自定义时间范围
	startTimeStr := c.DefaultQuery("start_time", "")
	endTimeStr := c.DefaultQuery("end_time", "")

	if startTimeStr == "" || endTimeStr == "" {
		// 默认最近1小时
		return monitoringEntity.NewTimeRangeFromDuration(1 * time.Hour), nil
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return nil, err
	}

	return monitoringEntity.NewTimeRange(startTime, endTime)
}

// parseInt 解析整数，提供默认值
func (h *MonitoringHandler) parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

// RecordMetricRequest 记录指标请求
type RecordMetricRequest struct {
	TenantID    string                       `json:"tenant_id" binding:"required"`
	MetricType  monitoringEntity.MetricType  `json:"metric_type" binding:"required"`
	Value       float64                      `json:"value" binding:"required"`
	Tags        monitoringEntity.MetricTags  `json:"tags"`
}

// RecordMetric 记录指标
// POST /api/v1/monitoring/metrics/record
func (h *MonitoringHandler) RecordMetric(c *gin.Context) {
	ctx := c.Request.Context()

	var req RecordMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	if err := h.tenantMonitoringSvc.RecordMetric(ctx, req.TenantID, req.MetricType, req.Value, req.Tags); err != nil {
		h.logger.WithError(err).Error("Failed to record metric")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to record metric",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    201,
		Message: "Metric recorded successfully",
	})
}

// RecordMetricBatch 批量记录指标
// POST /api/v1/monitoring/metrics/batch
func (h *MonitoringHandler) RecordMetricBatch(c *gin.Context) {
	ctx := c.Request.Context()

	var requests []RecordMetricRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	metricsList := make([]*monitoringEntity.TenantMetrics, 0, len(requests))
	for _, req := range requests {
		metric := monitoringEntity.NewTenantMetrics(req.TenantID, req.MetricType, req.Value, req.Tags)
		metricsList = append(metricsList, metric)
	}

	if err := h.tenantMonitoringSvc.RecordMetricBatch(ctx, metricsList); err != nil {
		h.logger.WithError(err).Error("Failed to record metrics batch")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to record metrics batch",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    201,
		Message: "Metrics recorded successfully",
		Data: gin.H{
			"count": len(metricsList),
		},
	})
}

// GetDashboardData 获取监控大盘数据
// GET /api/v1/monitoring/dashboard
func (h *MonitoringHandler) GetDashboardData(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "tenant_id is required",
		})
		return
	}

	// 并行获取各类数据
	type result struct {
		overview    *monitoringEntity.TenantOverview
		qpsMetrics  *monitoringEntity.QPSMetrics
		rtMetrics   *monitoringEntity.ResponseTimeMetrics
		erMetrics   *monitoringEntity.ErrorRateMetrics
		healthScore *monitoringEntity.HealthScore
		err         error
	}

	ch := make(chan result, 1)

	go func() {
		overview, err := h.tenantMonitoringSvc.GetTenantOverview(ctx, tenantID, tenantID)
		ch <- result{overview: overview, err: err}
	}()

	r := <-ch
	if r.err != nil {
		h.logger.WithError(r.err).Error("Failed to get dashboard data")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to get dashboard data",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Success",
		Data:    r.overview,
	})
}
