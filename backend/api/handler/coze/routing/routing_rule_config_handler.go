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

// RoutingRuleConfigHandler 路由规则配置Handler
// 职责：处理路由规则配置相关的HTTP请求
type RoutingRuleConfigHandler struct {
	ruleService *routingapp.RoutingRuleConfigService
	logger       *zap.Logger
}

// NewRoutingRuleConfigHandler 创建路由规则配置Handler实例
func NewRoutingRuleConfigHandler(
	ruleService *routingapp.RoutingRuleConfigService,
	logger *zap.Logger,
) *RoutingRuleConfigHandler {
	return &RoutingRuleConfigHandler{
		ruleService: ruleService,
		logger:       logger,
	}
}

// CreateRuleRequest 创建规则请求
type CreateRuleRequest struct {
	TenantID      string                 `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	RuleName      string                 `json:"rule_name" binding:"required,min=1,max=100"`
	Description   string                 `json:"description" binding:"omitempty,max=500"`
	RuleType      string                 `json:"rule_type" binding:"required,oneof=intent keyword regex hybrid"`
	Priority      int                    `json:"priority" binding:"required,min=1,max=1000"`
	Conditions    []RuleCondition        `json:"conditions" binding:"required,min=1"`
	TargetAgentID string                 `json:"target_agent_id" binding:"required"`
	TargetType    string                 `json:"target_type" binding:"required,oneof=agent workflow"`
	IsActive      bool                   `json:"is_active" binding:"omitempty"`
	Metadata      map[string]interface{} `json:"metadata" binding:"omitempty"`
}

// RuleCondition 规则条件
type RuleCondition struct {
	ConditionType string                 `json:"condition_type" binding:"required,oneof=intent keyword regex entity"`
	Field         string                 `json:"field" binding:"required"`
	Operator      string                 `json:"operator" binding:"required,oneof=contains equals matches starts_with ends_with"`
	Value         interface{}            `json:"value" binding:"required"`
	Weight        float64                `json:"weight" binding:"omitempty,min=0,max=1"`
}

// CreateRuleResponse 创建规则响应
type CreateRuleResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    *RuleInfo    `json:"data,omitempty"`
}

// RuleInfo 规则信息
type RuleInfo struct {
	RuleID        string                 `json:"rule_id"`
	TenantID      string                 `json:"tenant_id"`
	RuleName      string                 `json:"rule_name"`
	Description   string                 `json:"description"`
	RuleType      string                 `json:"rule_type"`
	Priority      int                    `json:"priority"`
	Conditions    []RuleCondition        `json:"conditions"`
	TargetAgentID string                 `json:"target_agent_id"`
	TargetType    string                 `json:"target_type"`
	IsActive      bool                   `json:"is_active"`
	HitCount      int64                  `json:"hit_count"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     int64                  `json:"created_at"`
	UpdatedAt     int64                  `json:"updated_at"`
}

// CreateRule 创建规则
// @router /api/routing/rules [POST]
func (h *RoutingRuleConfigHandler) CreateRule(ctx context.Context, c *app.RequestContext) {
	var req CreateRuleRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if !req.IsActive {
		req.IsActive = true
	}

	// 转换条件
	conditions := make([]routingapp.RuleCondition, len(req.Conditions))
	for i, cond := range req.Conditions {
		conditions[i] = routingapp.RuleCondition{
			ConditionType: cond.ConditionType,
			Field:         cond.Field,
			Operator:      cond.Operator,
			Value:         cond.Value,
			Weight:        cond.Weight,
		}
	}

	// 调用Service层创建规则
	rule, err := h.ruleService.CreateRule(ctx, req.TenantID, req.RuleName, req.Description, req.RuleType,
		req.Priority, conditions, req.TargetAgentID, req.TargetType, req.IsActive, req.Metadata)
	if err != nil {
		h.logger.Error("failed to create rule",
			zap.String("tenant_id", req.TenantID),
			zap.String("rule_name", req.RuleName),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	responseConditions := make([]RuleCondition, len(rule.Conditions))
	for i, cond := range rule.Conditions {
		responseConditions[i] = RuleCondition{
			ConditionType: cond.ConditionType,
			Field:         cond.Field,
			Operator:      cond.Operator,
			Value:         cond.Value,
			Weight:        cond.Weight,
		}
	}

	c.JSON(http.StatusOK, CreateRuleResponse{
		Code:    0,
		Message: "success",
		Data: &RuleInfo{
			RuleID:        rule.RuleID,
			TenantID:      rule.TenantID,
			RuleName:      rule.RuleName,
			Description:   rule.Description,
			RuleType:      rule.RuleType,
			Priority:      rule.Priority,
			Conditions:    responseConditions,
			TargetAgentID: rule.TargetAgentID,
			TargetType:    rule.TargetType,
			IsActive:      rule.IsActive,
			HitCount:      rule.HitCount,
			Metadata:      rule.Metadata,
			CreatedAt:     rule.CreatedAt,
			UpdatedAt:     rule.UpdatedAt,
		},
	})
}

// UpdateRuleRequest 更新规则请求
type UpdateRuleRequest struct {
	RuleName      *string                `json:"rule_name" binding:"omitempty,min=1,max=100"`
	Description   *string                `json:"description" binding:"omitempty,max=500"`
	Priority      *int                   `json:"priority" binding:"omitempty,min=1,max=1000"`
	Conditions    []RuleCondition        `json:"conditions" binding:"omitempty,min=1"`
	TargetAgentID *string                `json:"target_agent_id" binding:"omitempty"`
	IsActive      *bool                  `json:"is_active" binding:"omitempty"`
	Metadata      map[string]interface{} `json:"metadata" binding:"omitempty"`
}

// UpdateRuleResponse 更新规则响应
type UpdateRuleResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    *RuleInfo `json:"data,omitempty"`
}

// UpdateRule 更新规则
// @router /api/routing/rules/:id [PUT]
func (h *RoutingRuleConfigHandler) UpdateRule(ctx context.Context, c *app.RequestContext) {
	ruleID := c.Param("id")
	if ruleID == "" {
		httputil.BadRequest(c, "rule_id is required")
		return
	}

	var req UpdateRuleRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 转换条件
	var conditions []routingapp.RuleCondition
	if req.Conditions != nil {
		conditions = make([]routingapp.RuleCondition, len(req.Conditions))
		for i, cond := range req.Conditions {
			conditions[i] = routingapp.RuleCondition{
				ConditionType: cond.ConditionType,
				Field:         cond.Field,
				Operator:      cond.Operator,
				Value:         cond.Value,
				Weight:        cond.Weight,
			}
		}
	}

	// 调用Service层更新规则
	rule, err := h.ruleService.UpdateRule(ctx, ruleID, req.RuleName, req.Description, req.Priority,
		conditions, req.TargetAgentID, req.IsActive, req.Metadata)
	if err != nil {
		h.logger.Error("failed to update rule",
			zap.String("rule_id", ruleID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	responseConditions := make([]RuleCondition, len(rule.Conditions))
	for i, cond := range rule.Conditions {
		responseConditions[i] = RuleCondition{
			ConditionType: cond.ConditionType,
			Field:         cond.Field,
			Operator:      cond.Operator,
			Value:         cond.Value,
			Weight:        cond.Weight,
		}
	}

	c.JSON(http.StatusOK, UpdateRuleResponse{
		Code:    0,
		Message: "success",
		Data: &RuleInfo{
			RuleID:        rule.RuleID,
			TenantID:      rule.TenantID,
			RuleName:      rule.RuleName,
			Description:   rule.Description,
			RuleType:      rule.RuleType,
			Priority:      rule.Priority,
			Conditions:    responseConditions,
			TargetAgentID: rule.TargetAgentID,
			TargetType:    rule.TargetType,
			IsActive:      rule.IsActive,
			HitCount:      rule.HitCount,
			Metadata:      rule.Metadata,
			CreatedAt:     rule.CreatedAt,
			UpdatedAt:     rule.UpdatedAt,
		},
	})
}

// EnableRuleRequest 启用规则请求
type EnableRuleRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=200"`
}

// EnableRuleResponse 启用规则响应
type EnableRuleResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// EnableRule 启用规则
// @router /api/routing/rules/:id/enable [POST]
func (h *RoutingRuleConfigHandler) EnableRule(ctx context.Context, c *app.RequestContext) {
	ruleID := c.Param("id")
	if ruleID == "" {
		httputil.BadRequest(c, "rule_id is required")
		return
	}

	var req EnableRuleRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层启用规则
	if err := h.ruleService.EnableRule(ctx, ruleID, req.Reason); err != nil {
		h.logger.Error("failed to enable rule",
			zap.String("rule_id", ruleID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, EnableRuleResponse{
		Code:    0,
		Message: "success",
	})
}

// DisableRuleRequest 禁用规则请求
type DisableRuleRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=200"`
}

// DisableRuleResponse 禁用规则响应
type DisableRuleResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DisableRule 禁用规则
// @router /api/routing/rules/:id/disable [POST]
func (h *RoutingRuleConfigHandler) DisableRule(ctx context.Context, c *app.RequestContext) {
	ruleID := c.Param("id")
	if ruleID == "" {
		httputil.BadRequest(c, "rule_id is required")
		return
	}

	var req DisableRuleRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层禁用规则
	if err := h.ruleService.DisableRule(ctx, ruleID, req.Reason); err != nil {
		h.logger.Error("failed to disable rule",
			zap.String("rule_id", ruleID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, DisableRuleResponse{
		Code:    0,
		Message: "success",
	})
}

// ListRulesRequest 列出规则请求
type ListRulesRequest struct {
	TenantID   string  `form:"tenant_id" binding:"required" vd:"len($) > 0"`
	RuleType   string  `form:"rule_type" binding:"omitempty,oneof=intent keyword regex hybrid"`
	IsActive   *bool   `form:"is_active"`
	TargetType string  `form:"target_type" binding:"omitempty,oneof=agent workflow"`
	Page       int     `form:"page" binding:"omitempty,min=1"`
	PageSize   int     `form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy     string  `form:"sort_by" binding:"omitempty,oneof=priority hit_count created_at updated_at"`
	SortOrder  string  `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// ListRulesResponse 列出规则响应
type ListRulesResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    *RuleListData   `json:"data,omitempty"`
}

// RuleListData 规则列表数据
type RuleListData struct {
	Rules    []RuleInfo `json:"rules"`
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// ListRules 列出规则
// @router /api/routing/rules [GET]
func (h *RoutingRuleConfigHandler) ListRules(ctx context.Context, c *app.RequestContext) {
	var req ListRulesRequest
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
	if req.SortBy == "" {
		req.SortBy = "priority"
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	// 调用Service层列出规则
	rules, total, err := h.ruleService.ListRules(ctx, req.TenantID, req.RuleType, req.IsActive,
		req.TargetType, req.Page, req.PageSize, req.SortBy, req.SortOrder)
	if err != nil {
		h.logger.Error("failed to list rules",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	ruleList := make([]RuleInfo, len(rules))
	for i, rule := range rules {
		conditions := make([]RuleCondition, len(rule.Conditions))
		for j, cond := range rule.Conditions {
			conditions[j] = RuleCondition{
				ConditionType: cond.ConditionType,
				Field:         cond.Field,
				Operator:      cond.Operator,
				Value:         cond.Value,
				Weight:        cond.Weight,
			}
		}

		ruleList[i] = RuleInfo{
			RuleID:        rule.RuleID,
			TenantID:      rule.TenantID,
			RuleName:      rule.RuleName,
			Description:   rule.Description,
			RuleType:      rule.RuleType,
			Priority:      rule.Priority,
			Conditions:    conditions,
			TargetAgentID: rule.TargetAgentID,
			TargetType:    rule.TargetType,
			IsActive:      rule.IsActive,
			HitCount:      rule.HitCount,
			Metadata:      rule.Metadata,
			CreatedAt:     rule.CreatedAt,
			UpdatedAt:     rule.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, ListRulesResponse{
		Code:    0,
		Message: "success",
		Data: &RuleListData{
			Rules:    ruleList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	})
}

// ValidateRuleRequest 验证规则请求
type ValidateRuleRequest struct {
	RuleName      string                 `json:"rule_name" binding:"required"`
	RuleType      string                 `json:"rule_type" binding:"required,oneof=intent keyword regex hybrid"`
	Conditions    []RuleCondition        `json:"conditions" binding:"required,min=1"`
	Priority      int                    `json:"priority" binding:"required"`
	TargetAgentID string                 `json:"target_agent_id" binding:"required"`
	TargetType    string                 `json:"target_type" binding:"required,oneof=agent workflow"`
}

// ValidateRuleResponse 验证规则响应
type ValidateRuleResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *ValidationResult `json:"data,omitempty"`
}

// ValidationResult 验证结果
type ValidationResult struct {
	IsValid       bool     `json:"is_valid"`
	Errors        []string `json:"errors,omitempty"`
	Warnings      []string `json:"warnings,omitempty"`
	Conflicts     []string `json:"conflicts,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}

// ValidateRule 验证规则
// @router /api/routing/rules/validate [POST]
func (h *RoutingRuleConfigHandler) ValidateRule(ctx context.Context, c *app.RequestContext) {
	var req ValidateRuleRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 转换条件
	conditions := make([]routingapp.RuleCondition, len(req.Conditions))
	for i, cond := range req.Conditions {
		conditions[i] = routingapp.RuleCondition{
			ConditionType: cond.ConditionType,
			Field:         cond.Field,
			Operator:      cond.Operator,
			Value:         cond.Value,
			Weight:        cond.Weight,
		}
	}

	// 调用Service层验证规则
	result, err := h.ruleService.ValidateRule(ctx, req.RuleName, req.RuleType, conditions,
		req.Priority, req.TargetAgentID, req.TargetType)
	if err != nil {
		h.logger.Error("failed to validate rule",
			zap.String("rule_name", req.RuleName),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, ValidateRuleResponse{
		Code:    0,
		Message: "success",
		Data: &ValidationResult{
			IsValid:         result.IsValid,
			Errors:          result.Errors,
			Warnings:        result.Warnings,
			Conflicts:       result.Conflicts,
			Recommendations: result.Recommendations,
		},
	})
}

// TestRuleRequest 测试规则请求
type TestRuleRequest struct {
	RuleID    string                 `json:"rule_id" binding:"required"`
	TestCases []RuleTestCase        `json:"test_cases" binding:"required,min=1,max=20"`
}

// RuleTestCase 规则测试用例
type RuleTestCase struct {
	Input       string                 `json:"input" binding:"required"`
	ExpectedMatch bool                 `json:"expected_match" binding:"required"`
	Context     map[string]interface{} `json:"context" binding:"omitempty"`
}

// TestRuleResponse 测试规则响应
type TestRuleResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    *TestResultData `json:"data,omitempty"`
}

// TestResultData 测试结果数据
type TestResultData struct {
	RuleID       string           `json:"rule_id"`
	RuleName     string           `json:"rule_name"`
	TotalTests   int              `json:"total_tests"`
	PassedTests  int              `json:"passed_tests"`
	FailedTests  int              `json:"failed_tests"`
	PassRate     float64          `json:"pass_rate"`
	TestResults  []TestCaseResult `json:"test_results"`
}

// TestCaseResult 测试用例结果
type TestCaseResult struct {
	Input          string `json:"input"`
	ExpectedMatch  bool   `json:"expected_match"`
	ActualMatch    bool   `json:"actual_match"`
	Passed         bool   `json:"passed"`
	MatchScore     float64 `json:"match_score"`
	Reason         string `json:"reason,omitempty"`
}

// TestRule 测试规则
// @router /api/routing/rules/:id/test [POST]
func (h *RoutingRuleConfigHandler) TestRule(ctx context.Context, c *app.RequestContext) {
	ruleID := c.Param("id")
	if ruleID == "" {
		httputil.BadRequest(c, "rule_id is required")
		return
	}

	var req TestRuleRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	req.RuleID = ruleID

	// 转换测试用例
	testCases := make([]routingapp.RuleTestCase, len(req.TestCases))
	for i, tc := range req.TestCases {
		testCases[i] = routingapp.RuleTestCase{
			Input:         tc.Input,
			ExpectedMatch: tc.ExpectedMatch,
			Context:       tc.Context,
		}
	}

	// 调用Service层测试规则
	result, err := h.ruleService.TestRule(ctx, req.RuleID, testCases)
	if err != nil {
		h.logger.Error("failed to test rule",
			zap.String("rule_id", req.RuleID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	testResults := make([]TestCaseResult, len(result.TestResults))
	for i, tr := range result.TestResults {
		testResults[i] = TestCaseResult{
			Input:         tr.Input,
			ExpectedMatch: tr.ExpectedMatch,
			ActualMatch:   tr.ActualMatch,
			Passed:        tr.Passed,
			MatchScore:    tr.MatchScore,
			Reason:        tr.Reason,
		}
	}

	c.JSON(http.StatusOK, TestRuleResponse{
		Code:    0,
		Message: "success",
		Data: &TestResultData{
			RuleID:       result.RuleID,
			RuleName:     result.RuleName,
			TotalTests:   result.TotalTests,
			PassedTests:  result.PassedTests,
			FailedTests:  result.FailedTests,
			PassRate:     result.PassRate,
			TestResults:  testResults,
		},
	})
}
