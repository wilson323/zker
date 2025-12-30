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

package routing

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	routingapp "github.com/coze-dev/coze-studio/backend/application/routing"
)

// ABTestHandler A/B测试Handler
// 职责：处理A/B测试相关的HTTP请求
type ABTestHandler struct {
	abTestService *routingapp.ABTestService
	logger        *zap.Logger
}

// NewABTestHandler 创建A/B测试Handler实例
func NewABTestHandler(
	abTestService *routingapp.ABTestService,
	logger *zap.Logger,
) *ABTestHandler {
	return &ABTestHandler{
		abTestService: abTestService,
		logger:        logger,
	}
}

// CreateABTestRequest 创建A/B测试请求
type CreateABTestRequest struct {
	TenantID       string                 `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	TestName       string                 `json:"test_name" binding:"required,min=1,max=100"`
	Description    string                 `json:"description" binding:"omitempty,max=500"`
	StrategyA      string                 `json:"strategy_a" binding:"required"` // 控制组策略ID
	StrategyB      string                 `json:"strategy_b" binding:"required"` // 实验组策略ID
	TrafficSplit   int                    `json:"traffic_split" binding:"omitempty,min=1,max=99"` // A组流量百分比
	SampleSize     int                    `json:"sample_size" binding:"omitempty,min=100,max=100000"`
	SuccessMetrics []string               `json:"success_metrics" binding:"required,min=1"` // 成功指标列表
	Metadata       map[string]interface{} `json:"metadata" binding:"omitempty"`
}

// CreateABTestResponse 创建A/B测试响应
type CreateABTestResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *ABTestInfo `json:"data,omitempty"`
}

// ABTestInfo A/B测试信息
type ABTestInfo struct {
	TestID         string                 `json:"test_id"`
	TenantID       string                 `json:"tenant_id"`
	TestName       string                 `json:"test_name"`
	Description    string                 `json:"description"`
	StrategyA      string                 `json:"strategy_a"`
	StrategyB      string                 `json:"strategy_b"`
	TrafficSplit   int                    `json:"traffic_split"`
	SampleSize     int                    `json:"sample_size"`
	SuccessMetrics []string               `json:"success_metrics"`
	Status         string                 `json:"status"` // running, paused, completed
	StartTime      *int64                 `json:"start_time,omitempty"`
	EndTime        *int64                 `json:"end_time,omitempty"`
	Results        *ABTestResultsInfo     `json:"results,omitempty"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      int64                  `json:"created_at"`
	UpdatedAt      int64                  `json:"updated_at"`
}

// ABTestResultsInfo A/B测试结果信息
type ABTestResultsInfo struct {
	TotalRequests      int      `json:"total_requests"`
	GroupARequests     int      `json:"group_a_requests"`
	GroupBRequests     int      `json:"group_b_requests"`
	GroupASuccessRate  float64  `json:"group_a_success_rate"`
	GroupBSuccessRate  float64  `json:"group_b_success_rate"`
	Improvement        float64  `json:"improvement"`
	PValue             float64  `json:"p_value"`
	IsSignificant      bool     `json:"is_significant"`
	ConfidenceInterval []float64 `json:"confidence_interval"`
	Recommendation     string   `json:"recommendation"`
}

// CreateABTest 创建A/B测试
// @router /api/routing/abtests [POST]
func (h *ABTestHandler) CreateABTest(ctx context.Context, c *app.RequestContext) {
	var req CreateABTestRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.TrafficSplit == 0 {
		req.TrafficSplit = 50
	}
	if req.SampleSize == 0 {
		req.SampleSize = 1000
	}

	// 调用Service层创建测试
	test, err := h.abTestService.CreateABTest(ctx, req.TenantID, req.TestName, req.Description,
		req.StrategyA, req.StrategyB, req.TrafficSplit, req.SampleSize, req.SuccessMetrics, req.Metadata)
	if err != nil {
		h.logger.Error("failed to create ab test",
			zap.String("tenant_id", req.TenantID),
			zap.String("test_name", req.TestName),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	var results *ABTestResultsInfo
	if test.Results != nil {
		results = &ABTestResultsInfo{
			TotalRequests:      test.Results.TotalRequests,
			GroupARequests:     test.Results.GroupARequests,
			GroupBRequests:     test.Results.GroupBRequests,
			GroupASuccessRate:  test.Results.GroupASuccessRate,
			GroupBSuccessRate:  test.Results.GroupBSuccessRate,
			Improvement:        test.Results.Improvement,
			PValue:             test.Results.PValue,
			IsSignificant:      test.Results.IsSignificant,
			ConfidenceInterval: test.Results.ConfidenceInterval,
			Recommendation:     test.Results.Recommendation,
		}
	}

	c.JSON(http.StatusOK, CreateABTestResponse{
		Code:    0,
		Message: "success",
		Data: &ABTestInfo{
			TestID:         test.TestID,
			TenantID:       test.TenantID,
			TestName:       test.TestName,
			Description:    test.Description,
			StrategyA:      test.StrategyA,
			StrategyB:      test.StrategyB,
			TrafficSplit:   test.TrafficSplit,
			SampleSize:     test.SampleSize,
			SuccessMetrics: test.SuccessMetrics,
			Status:         test.Status,
			StartTime:      test.StartTime,
			EndTime:        test.EndTime,
			Results:        results,
			Metadata:       test.Metadata,
			CreatedAt:      test.CreatedAt,
			UpdatedAt:      test.UpdatedAt,
		},
	})
}

// ExecuteABTestRequest 执行A/B测试请求
type ExecuteABTestRequest struct {
	TestID    string                 `json:"test_id" binding:"required"`
	UserInput string                 `json:"user_input" binding:"required"`
	Context   map[string]interface{} `json:"context" binding:"omitempty"`
	UserID    string                 `json:"user_id" binding:"omitempty"`
}

// ExecuteABTestResponse 执行A/B测试响应
type ExecuteABTestResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    *ABTestExecutionInfo `json:"data,omitempty"`
}

// ABTestExecutionInfo A/B测试执行信息
type ABTestExecutionInfo struct {
	TestID       string `json:"test_id"`
	AssignedGroup string `json:"assigned_group"` // A or B
	StrategyID   string `json:"strategy_id"`
	RoutingDecision *RoutingDecisionInfo `json:"routing_decision"`
}

// RoutingDecisionInfo 路由决策信息
type RoutingDecisionInfo struct {
	AgentID      string `json:"agent_id"`
	AgentName    string `json:"agent_name"`
	Reason       string `json:"reason"`
	Confidence   float64 `json:"confidence"`
}

// ExecuteABTest 执行A/B测试
// @router /api/routing/abtests/:id/execute [POST]
func (h *ABTestHandler) ExecuteABTest(ctx context.Context, c *app.RequestContext) {
	testID := c.Param("id")
	if testID == "" {
		httputil.BadRequest(c, "test_id is required")
		return
	}

	var req ExecuteABTestRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	req.TestID = testID

	// 调用Service层执行测试
	result, err := h.abTestService.ExecuteABTest(ctx, req.TestID, req.UserInput, req.Context, req.UserID)
	if err != nil {
		h.logger.Error("failed to execute ab test",
			zap.String("test_id", req.TestID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, ExecuteABTestResponse{
		Code:    0,
		Message: "success",
		Data: &ABTestExecutionInfo{
			TestID:        result.TestID,
			AssignedGroup: result.AssignedGroup,
			StrategyID:    result.StrategyID,
			RoutingDecision: &RoutingDecisionInfo{
				AgentID:    result.RoutingDecision.AgentID,
				AgentName:  result.RoutingDecision.AgentName,
				Reason:     result.RoutingDecision.Reason,
				Confidence: result.RoutingDecision.Confidence,
			},
		},
	})
}

// GetABTestResultsResponse 获取A/B测试结果响应
type GetABTestResultsResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *ABTestResultsInfo `json:"data,omitempty"`
}

// GetABTestResults 获取A/B测试结果
// @router /api/routing/abtests/:id/results [GET]
func (h *ABTestHandler) GetABTestResults(ctx context.Context, c *app.RequestContext) {
	testID := c.Param("id")
	if testID == "" {
		httputil.BadRequest(c, "test_id is required")
		return
	}

	// 调用Service层获取测试结果
	result, err := h.abTestService.GetABTestResults(ctx, testID)
	if err != nil {
		h.logger.Error("failed to get ab test results",
			zap.String("test_id", testID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, GetABTestResultsResponse{
		Code:    0,
		Message: "success",
		Data: &ABTestResultsInfo{
			TotalRequests:      result.TotalRequests,
			GroupARequests:     result.GroupARequests,
			GroupBRequests:     result.GroupBRequests,
			GroupASuccessRate:  result.GroupASuccessRate,
			GroupBSuccessRate:  result.GroupBSuccessRate,
			Improvement:        result.Improvement,
			PValue:             result.PValue,
			IsSignificant:      result.IsSignificant,
			ConfidenceInterval: result.ConfidenceInterval,
			Recommendation:     result.Recommendation,
		},
	})
}

// ConcludeABTestRequest 结束A/B测试请求
type ConcludeABTestRequest struct {
	WinningStrategy string `json:"winning_strategy" binding:"omitempty,oneof=A B none"`
	Reason          string `json:"reason" binding:"omitempty,max=500"`
	ApplyWinner     bool   `json:"apply_winner" binding:"omitempty"` // 是否应用获胜策略
}

// ConcludeABTestResponse 结束A/B测试响应
type ConcludeABTestResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *ABTestInfo `json:"data,omitempty"`
}

// ConcludeABTest 结束A/B测试
// @router /api/routing/abtests/:id/conclude [POST]
func (h *ABTestHandler) ConcludeABTest(ctx context.Context, c *app.RequestContext) {
	testID := c.Param("id")
	if testID == "" {
		httputil.BadRequest(c, "test_id is required")
		return
	}

	var req ConcludeABTestRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层结束测试
	test, err := h.abTestService.ConcludeABTest(ctx, testID, req.WinningStrategy, req.Reason, req.ApplyWinner)
	if err != nil {
		h.logger.Error("failed to conclude ab test",
			zap.String("test_id", testID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	var results *ABTestResultsInfo
	if test.Results != nil {
		results = &ABTestResultsInfo{
			TotalRequests:      test.Results.TotalRequests,
			GroupARequests:     test.Results.GroupARequests,
			GroupBRequests:     test.Results.GroupBRequests,
			GroupASuccessRate:  test.Results.GroupASuccessRate,
			GroupBSuccessRate:  test.Results.GroupBSuccessRate,
			Improvement:        test.Results.Improvement,
			PValue:             test.Results.PValue,
			IsSignificant:      test.Results.IsSignificant,
			ConfidenceInterval: test.Results.ConfidenceInterval,
			Recommendation:     test.Results.Recommendation,
		}
	}

	c.JSON(http.StatusOK, ConcludeABTestResponse{
		Code:    0,
		Message: "success",
		Data: &ABTestInfo{
			TestID:         test.TestID,
			TenantID:       test.TenantID,
			TestName:       test.TestName,
			Description:    test.Description,
			StrategyA:      test.StrategyA,
			StrategyB:      test.StrategyB,
			TrafficSplit:   test.TrafficSplit,
			SampleSize:     test.SampleSize,
			SuccessMetrics: test.SuccessMetrics,
			Status:         test.Status,
			StartTime:      test.StartTime,
			EndTime:        test.EndTime,
			Results:        results,
			Metadata:       test.Metadata,
			CreatedAt:      test.CreatedAt,
			UpdatedAt:      test.UpdatedAt,
		},
	})
}

// ListABTestsRequest 列出A/B测试请求
type ListABTestsRequest struct {
	TenantID string  `form:"tenant_id" binding:"required" vd:"len($) > 0"`
	Status   string  `form:"status" binding:"omitempty,oneof=running paused completed"`
	Page     int     `form:"page" binding:"omitempty,min=1"`
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// ListABTestsResponse 列出A/B测试响应
type ListABTestsResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *ABTestListData  `json:"data,omitempty"`
}

// ABTestListData A/B测试列表数据
type ABTestListData struct {
	Tests    []ABTestInfo `json:"tests"`
	Total    int          `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

// ListABTests 列出A/B测试
// @router /api/routing/abtests [GET]
func (h *ABTestHandler) ListABTests(ctx context.Context, c *app.RequestContext) {
	var req ListABTestsRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认分页
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 调用Service层列出测试
	tests, total, err := h.abTestService.ListABTests(ctx, req.TenantID, req.Status, req.Page, req.PageSize)
	if err != nil {
		h.logger.Error("failed to list ab tests",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	testList := make([]ABTestInfo, len(tests))
	for i, test := range tests {
		var results *ABTestResultsInfo
		if test.Results != nil {
			results = &ABTestResultsInfo{
				TotalRequests:      test.Results.TotalRequests,
				GroupARequests:     test.Results.GroupARequests,
				GroupBRequests:     test.Results.GroupBRequests,
				GroupASuccessRate:  test.Results.GroupASuccessRate,
				GroupBSuccessRate:  test.Results.GroupBSuccessRate,
				Improvement:        test.Results.Improvement,
				PValue:             test.Results.PValue,
				IsSignificant:      test.Results.IsSignificant,
				ConfidenceInterval: test.Results.ConfidenceInterval,
				Recommendation:     test.Results.Recommendation,
			}
		}

		testList[i] = ABTestInfo{
			TestID:         test.TestID,
			TenantID:       test.TenantID,
			TestName:       test.TestName,
			Description:    test.Description,
			StrategyA:      test.StrategyA,
			StrategyB:      test.StrategyB,
			TrafficSplit:   test.TrafficSplit,
			SampleSize:     test.SampleSize,
			SuccessMetrics: test.SuccessMetrics,
			Status:         test.Status,
			StartTime:      test.StartTime,
			EndTime:        test.EndTime,
			Results:        results,
			Metadata:       test.Metadata,
			CreatedAt:      test.CreatedAt,
			UpdatedAt:      test.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, ListABTestsResponse{
		Code:    0,
		Message: "success",
		Data: &ABTestListData{
			Tests:    testList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	})
}
