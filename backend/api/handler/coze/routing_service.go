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

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/model/routing"
	routingapp "github.com/coze-dev/coze-studio/backend/application/routing"
)

// ==================== 路由规则管理接口 ====================

// CreateRoutingRule 创建路由规则
// @router /api/routing/rules [POST]
func CreateRoutingRule(ctx context.Context, c *app.RequestContext) {
	var err error
	var req routing.CreateRoutingRuleRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := routingapp.RoutingAppSVC.CreateRoutingRule(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetRoutingRule 获取路由规则
// @router /api/routing/rules/:rule_id [GET]
func GetRoutingRule(ctx context.Context, c *app.RequestContext) {
	ruleID := c.Param("rule_id")
	if ruleID == "" {
		invalidParamRequestResponse(c, "rule_id is required")
		return
	}

	resp, err := routingapp.RoutingAppSVC.GetRoutingRule(ctx, ruleID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateRoutingRule 更新路由规则
// @router /api/routing/rules/:rule_id [PUT]
func UpdateRoutingRule(ctx context.Context, c *app.RequestContext) {
	ruleID := c.Param("rule_id")
	if ruleID == "" {
		invalidParamRequestResponse(c, "rule_id is required")
		return
	}

	var err error
	var req routing.UpdateRoutingRuleRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := routingapp.RoutingAppSVC.UpdateRoutingRule(ctx, ruleID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteRoutingRule 删除路由规则
// @router /api/routing/rules/:rule_id [DELETE]
func DeleteRoutingRule(ctx context.Context, c *app.RequestContext) {
	ruleID := c.Param("rule_id")
	if ruleID == "" {
		invalidParamRequestResponse(c, "rule_id is required")
		return
	}

	err := routingapp.RoutingAppSVC.DeleteRoutingRule(ctx, ruleID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
}

// ListRoutingRules 列出路由规则
// @router /api/routing/rules [GET]
func ListRoutingRules(ctx context.Context, c *app.RequestContext) {
	var req routing.ListRoutingRulesRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := routingapp.RoutingAppSVC.ListRoutingRules(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== 路由决策接口 ====================

// ExecuteRouting 执行路由
// @router /api/routing/route [POST]
func ExecuteRouting(ctx context.Context, c *app.RequestContext) {
	var err error
	var req routing.ExecuteRoutingRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := routingapp.RoutingAppSVC.ExecuteRouting(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== 路由日志接口 ====================

// ListRoutingLogs 列出路由日志
// @router /api/routing/logs [GET]
func ListRoutingLogs(ctx context.Context, c *app.RequestContext) {
	var req routing.ListRoutingLogsRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := routingapp.RoutingAppSVC.ListRoutingLogs(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== 配置管理接口 ====================

// SetMatcherWeights 设置匹配器权重
// @router /api/routing/matcher/weights [POST]
func SetMatcherWeights(ctx context.Context, c *app.RequestContext) {
	var err error
	var req routing.SetMatcherWeightsRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = routingapp.RoutingAppSVC.SetMatcherWeights(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &routing.SetWeightsResponse{
		Code:    0,
		Message: "success",
	})
}

// SetScorerWeights 设置评分权重
// @router /api/routing/scorer/weights [POST]
func SetScorerWeights(ctx context.Context, c *app.RequestContext) {
	var err error
	var req routing.SetScorerWeightsRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = routingapp.RoutingAppSVC.SetScorerWeights(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &routing.SetWeightsResponse{
		Code:    0,
		Message: "success",
	})
}

// SetSimilarityThreshold 设置相似度阈值
// @router /api/routing/similarity/threshold [POST]
func SetSimilarityThreshold(ctx context.Context, c *app.RequestContext) {
	var err error
	var req routing.SetSimilarityThresholdRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = routingapp.RoutingAppSVC.SetSimilarityThreshold(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &routing.SetWeightsResponse{
		Code:    0,
		Message: "success",
	})
}

// GenerateBotEmbedding 生成Bot向量嵌入
// @router /api/bots/:bot_id/embedding [POST]
func GenerateBotEmbedding(ctx context.Context, c *app.RequestContext) {
	botID := c.Param("bot_id")
	if botID == "" {
		invalidParamRequestResponse(c, "bot_id is required")
		return
	}

	var err error
	var req routing.GenerateBotEmbeddingRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = routingapp.RoutingAppSVC.GenerateBotEmbedding(ctx, botID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &routing.GenerateBotEmbeddingResponse{
		Code:    0,
		Message: "success",
	})
}

// ==================== 服务健康检查接口 ====================

// HealthCheck 健康检查
// @router /api/routing/health/check [POST]
func HealthCheck(ctx context.Context, c *app.RequestContext) {
	var err error
	var req routing.HealthCheckRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := routingapp.RoutingAppSVC.HealthCheck(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &routing.HealthCheckResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// ==================== 负载均衡接口 ====================

// ConfigureLoadBalance 配置负载均衡
// @router /api/routing/load-balance/config [POST]
func ConfigureLoadBalance(ctx context.Context, c *app.RequestContext) {
	var err error
	var req routing.LoadBalanceConfigRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = routingapp.RoutingAppSVC.ConfigureLoadBalance(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &routing.SetWeightsResponse{
		Code:    0,
		Message: "success",
	})
}

// ==================== 熔断器接口 ====================

// ConfigureCircuitBreaker 配置熔断器
// @router /api/routing/circuit-breaker/config [POST]
func ConfigureCircuitBreaker(ctx context.Context, c *app.RequestContext) {
	var err error
	var req routing.CircuitBreakerConfigRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = routingapp.RoutingAppSVC.ConfigureCircuitBreaker(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &routing.SetWeightsResponse{
		Code:    0,
		Message: "success",
	})
}

// ==================== 路由统计接口 ====================

// GetRoutingStats 获取路由统计
// @router /api/routing/stats [GET]
func GetRoutingStats(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Query("tenant_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	resp, err := routingapp.RoutingAppSVC.GetRoutingStats(ctx, tenantID, startDate, endDate)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &routing.RoutingStatsResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}
