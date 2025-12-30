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

package coze

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"

	billingModel "github.com/coze-dev/coze-studio/backend/api/model/billing"
	billingService "github.com/coze-dev/coze-studio/backend/domain/billing/service"
)

// ============================================================
// Token计量API Handler
// ============================================================

var (
	tokenMeteringService *billingService.TokenMeteringService
)

// InitTokenMeteringHandler 初始化Token计量Handler
func InitTokenMeteringHandler(svc *billingService.TokenMeteringService) {
	tokenMeteringService = svc
}

// ==================== Token计量核心接口 ====================

// RecordTokenUsage 记录Token使用
// @router /api/v1/tenants/:tenant_id/tokens/usage [POST]
// @summary 记录Token使用
// @description 记录单次AI模型调用的Token使用情况和成本
// @tags Token计量
// @accept json
// @produce json
// @param tenant_id path string true "租户ID"
// @param request body billingModel.RecordTokenUsageRequest true "Token使用记录"
// @success 200 {object} billingModel.RecordTokenUsageResponse
// @failure 400 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func RecordTokenUsage(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 绑定并验证请求
	var req billingModel.RecordTokenUsageRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 3. 设置租户ID
	req.TenantID = tenantID

	// 4. 调用服务
	resp, err := tokenMeteringService.RecordTokenUsage(ctx, &billingService.RecordTokenUsageRequest{
		TenantID:       req.TenantID,
		UserID:         req.UserID,
		BotID:          req.BotID,
		ConversationID: req.ConversationID,
		MessageID:      req.MessageID,
		ModelProvider:  req.ModelProvider,
		ModelName:      req.ModelName,
		ModelVersion:   req.ModelVersion,
		InputTokens:    req.InputTokens,
		OutputTokens:   req.OutputTokens,
		TotalTokens:    req.TotalTokens,
		RequestType:    req.RequestType,
		ResponseTimeMs: req.ResponseTimeMs,
		LatencyMs:      req.LatencyMs,
		IsCached:       req.IsCached,
		Metadata:       req.Metadata,
	})

	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 5. 返回响应
	c.JSON(http.StatusOK, resp)
}

// BatchRecordTokenUsage 批量记录Token使用
// @router /api/v1/tenants/:tenant_id/tokens/usage/batch [POST]
// @summary 批量记录Token使用
// @description 批量记录多次AI模型调用的Token使用情况（最多1000条）
// @tags Token计量
// @accept json
// @produce json
// @param tenant_id path string true "租户ID"
// @param request body billingModel.BatchRecordTokenUsageRequest true "批量Token使用记录"
// @success 200 {object} billingModel.BatchRecordTokenUsageResponse
// @failure 400 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func BatchRecordTokenUsage(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 绑定并验证请求
	var req billingModel.BatchRecordTokenUsageRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 3. 转换为服务层请求
	records := make([]*billingService.RecordTokenUsageRequest, 0, len(req.Records))
	for _, record := range req.Records {
		records = append(records, &billingService.RecordTokenUsageRequest{
			TenantID:       tenantID,
			UserID:         record.UserID,
			BotID:          record.BotID,
			ConversationID: record.ConversationID,
			MessageID:      record.MessageID,
			ModelProvider:  record.ModelProvider,
			ModelName:      record.ModelName,
			ModelVersion:   record.ModelVersion,
			InputTokens:    record.InputTokens,
			OutputTokens:   record.OutputTokens,
			TotalTokens:    record.TotalTokens,
			RequestType:    record.RequestType,
			ResponseTimeMs: record.ResponseTimeMs,
			LatencyMs:      record.LatencyMs,
			IsCached:       record.IsCached,
			Metadata:       record.Metadata,
		})
	}

	// 4. 调用服务
	svcResp, err := tokenMeteringService.BatchRecordTokenUsage(ctx, records)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 5. 转换为API响应
	resp := &billingModel.BatchRecordTokenUsageResponse{
		SuccessCount: svcResp.SuccessCount,
		FailedCount:  svcResp.FailedCount,
		LogIDs:       svcResp.LogIDs,
		TotalCost:    svcResp.TotalCost,
	}

	c.JSON(http.StatusOK, resp)
}

// GetUsageStats 获取使用统计
// @router /api/v1/tenants/:tenant_id/tokens/stats [GET]
// @summary 获取Token使用统计
// @description 获取租户的Token使用统计汇总数据
// @tags Token计量
// @accept json
// @produce json
// @param tenant_id path string true "租户ID"
// @param bot_id query string false "Bot ID（可选，用于过滤特定Bot）"
// @param model_name query string false "模型名称（可选，用于过滤特定模型）"
// @param start_date query string false "开始日期（格式: 2006-01-02，默认: 7天前）"
// @param end_date query string false "结束日期（格式: 2006-01-02，默认: 今天）"
// @param granularity query string false "粒度（daily/hourly，默认: daily）"
// @success 200 {object} billingModel.UsageStatsResponse
// @failure 400 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func GetUsageStats(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 绑定查询参数
	var req billingModel.GetUsageStatsRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 3. 转换为服务层请求
	filter := &billingService.UsageStatsFilter{
		BotID:       req.BotID,
		ModelName:   req.ModelName,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Granularity: req.Granularity,
	}

	// 4. 调用服务
	svcResp, err := tokenMeteringService.GetUsageStats(ctx, tenantID, filter)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 5. 转换为API响应
	resp := &billingModel.UsageStatsResponse{
		TenantID:          svcResp.TenantID,
		TotalInputTokens:  svcResp.TotalInputTokens,
		TotalOutputTokens: svcResp.TotalOutputTokens,
		TotalTokens:       svcResp.TotalTokens,
		TotalCost:         svcResp.TotalCost,
		TotalRequests:     svcResp.TotalRequests,
		CachedRequests:    svcResp.CachedRequests,
		AvgResponseTime:   svcResp.AvgResponseTime,
		StartDate:         svcResp.StartDate,
		EndDate:           svcResp.EndDate,
	}

	c.JSON(http.StatusOK, resp)
}

// GetDailyUsageStats 获取每日使用趋势
// @router /api/v1/tenants/:tenant_id/tokens/daily [GET]
// @summary 获取每日Token使用趋势
// @description 获取租户每日的Token使用趋势数据
// @tags Token计量
// @accept json
// @produce json
// @param tenant_id path string true "租户ID"
// @param bot_id query string false "Bot ID（可选，用于过滤特定Bot）"
// @param start_date query string false "开始日期（格式: 2006-01-02，默认: 30天前）"
// @param end_date query string false "结束日期（格式: 2006-01-02，默认: 今天）"
// @success 200 {object} billingModel.DailyUsageStatsResponse
// @failure 400 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func GetDailyUsageStats(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 绑定查询参数
	botID := c.Query("bot_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// 3. 设置默认值（最近30天）
	if startDate == "" || endDate == "" {
		now := time.Now()
		endDate = now.Format("2006-01-02")
		startDate = now.AddDate(0, 0, -30).Format("2006-01-02")
	}

	// 4. 验证日期格式
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		invalidParamRequestResponse(c, "invalid start_date format, expected: 2006-01-02")
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		invalidParamRequestResponse(c, "invalid end_date format, expected: 2006-01-02")
		return
	}

	// 5. 调用服务
	var botIDPtr *string
	if botID != "" {
		botIDPtr = &botID
	}

	svcResp, err := tokenMeteringService.GetDailyUsageStats(ctx, tenantID, startDate, endDate, botIDPtr)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 6. 计算汇总统计
	totalTokens := int64(0)
	totalCost := 0.0
	totalRequests := int64(0)
	topModel := ""
	topModelCost := 0.0

	for _, stat := range svcResp {
		totalTokens += stat.TotalTokens
		totalCost += stat.TotalCost
		totalRequests += stat.TotalRequests

		// 找出成本最高的日期
		if stat.TotalCost > topModelCost {
			topModelCost = stat.TotalCost
			topModel = stat.Date
		}
	}

	// 7. 转换为API响应
	dailyStats := make([]*billingModel.DailyUsageStats, 0, len(svcResp))
	for _, stat := range svcResp {
		dailyStats = append(dailyStats, &billingModel.DailyUsageStats{
			Date:            stat.Date,
			TotalTokens:     stat.TotalTokens,
			TotalCost:       stat.TotalCost,
			TotalRequests:   stat.TotalRequests,
			CachedRequests:  stat.CachedRequests,
			AvgResponseTime: stat.AvgResponseTime,
		})
	}

	resp := &billingModel.DailyUsageStatsResponse{
		DailyStats: dailyStats,
		Total: &billingModel.AggregateStats{
			TotalTokens:  totalTokens,
			TotalCost:    totalCost,
			TotalRequests: totalRequests,
			TopModel:     topModel,
			TopModelCost: topModelCost,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// GetModelUsageStats 获取模型使用统计
// @router /api/v1/tenants/:tenant_id/tokens/models [GET]
// @summary 获取模型使用统计
// @description 按模型维度统计Token使用情况和成本分布
// @tags Token计量
// @accept json
// @produce json
// @param tenant_id path string true "租户ID"
// @param start_date query string false "开始日期（格式: 2006-01-02，默认: 7天前）"
// @param end_date query string false "结束日期（格式: 2006-01-02，默认: 今天）"
// @success 200 {object} billingModel.ModelUsageStatsResponse
// @failure 400 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func GetModelUsageStats(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 绑定查询参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// 3. 设置默认值（最近7天）
	if startDate == "" || endDate == "" {
		now := time.Now()
		endDate = now.Format("2006-01-02")
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
	}

	// 4. 验证日期格式
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		invalidParamRequestResponse(c, "invalid start_date format, expected: 2006-01-02")
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		invalidParamRequestResponse(c, "invalid end_date format, expected: 2006-01-02")
		return
	}

	// 5. 调用服务
	svcResp, err := tokenMeteringService.GetModelUsageStats(ctx, tenantID, startDate, endDate)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 6. 计算汇总统计
	totalTokens := int64(0)
	totalCost := 0.0
	totalRequests := int64(0)
	topModel := ""
	topModelCost := 0.0

	modelStats := make([]*billingModel.ModelUsageStats, 0, len(svcResp))
	for _, stat := range svcResp {
		totalTokens += stat.TotalTokens
		totalCost += stat.TotalCost
		totalRequests += stat.TotalRequests

		// 找出成本最高的模型
		if stat.TotalCost > topModelCost {
			topModelCost = stat.TotalCost
			topModel = stat.ModelName
		}

		// 计算成本占比
		costPercentage := 0.0
		if totalCost > 0 {
			costPercentage = (stat.TotalCost / totalCost) * 100
		}

		modelStats = append(modelStats, &billingModel.ModelUsageStats{
			ModelProvider:     stat.ModelProvider,
			ModelName:         stat.ModelName,
			TotalTokens:       stat.TotalTokens,
			TotalCost:         stat.TotalCost,
			TotalRequests:     stat.TotalRequests,
			AvgCostPer1KToken: stat.AvgCostPer1KToken,
			CostPercentage:   costPercentage,
		})
	}

	// 7. 构建响应
	resp := &billingModel.ModelUsageStatsResponse{
		ModelStats: modelStats,
		Total: &billingModel.AggregateStats{
			TotalTokens:  totalTokens,
			TotalCost:    totalCost,
			TotalRequests: totalRequests,
			TopModel:     topModel,
			TopModelCost: topModelCost,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== 辅助函数 ====================

// invalidParamRequestResponse 返回参数错误响应
func invalidParamRequestResponse(c *app.RequestContext, message string) {
	c.JSON(http.StatusBadRequest, &billingModel.ErrorResponse{
		Code:    "INVALID_PARAMETER",
		Message: "Invalid request parameter",
		Details: message,
	})
}

// internalServerErrorResponse 返回服务器错误响应
func internalServerErrorResponse(ctx context.Context, c *app.RequestContext, err error) {
	logger := logrus.StandardLogger()
	logger.WithError(err).WithField("path", string(c.Request.URI().Path())).Error("[TokenMetering] Internal server error")

	c.JSON(http.StatusInternalServerError, &billingModel.ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
		Details: err.Error(),
	})
}
