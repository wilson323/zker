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

package humaninloop

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/api/handler/coze"
	"github.com/coze-dev/coze-studio/backend/api/model/humaninloop"
	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
	hilservice "github.com/coze-dev/coze-studio/backend/domain/humaninloop/service"
)

// CollaborationHandler 协同任务处理器
type CollaborationHandler struct {
	orchestrator hilservice.CollaborationOrchestrator
	queueService hilservice.QueueService
	protocol     hilservice.ProtocolService
}

// NewCollaborationHandler 创建协同任务处理器
func NewCollaborationHandler(
	orchestrator hilservice.CollaborationOrchestrator,
	queueService hilservice.QueueService,
	protocol hilservice.ProtocolService,
) *CollaborationHandler {
	return &CollaborationHandler{
		orchestrator: orchestrator,
		queueService: queueService,
		protocol:     protocol,
	}
}

// ==================== 任务管理接口 ====================

// CreateTask 创建协同任务
// @router /api/human-in-loop/tasks [POST]
func CreateTask(ctx context.Context, c *app.RequestContext) {
	var req humaninloop.CreateTaskRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.invalidParamRequestResponse(c, err.Error())
		return
	}

	// 获取租户ID
	tenantID := coze.GetTenantIDFromContext(ctx)

	// 构建上下文
	var contextData map[string]interface{}
	if err := json.Unmarshal(req.Context, &contextData); err != nil {
		coze.invalidParamRequestResponse(c, "invalid context format")
		return
	}

	// 确定优先级
	priority := req.Priority
	if priority == "" {
		priority = entity.PriorityMedium
	}

	// 确定来源
	source := req.Source
	if source == "" {
		source = entity.SourceAI
	}

	// 创建任务请求
	createReq := &hilservice.CreateTaskRequest{
		TenantID:       tenantID,
		TaskType:       req.TaskType,
		Source:         source,
		Priority:       priority,
		ConversationID: req.ConversationID,
		MessageID:      req.MessageID,
		BotID:          req.BotID,
		WorkflowID:     req.WorkflowID,
		Context:        contextData,
		TriggerReason:  req.TriggerReason,
		AssignTo:       req.AssignTo,
		SLAMinutes:     req.SLAMinutes,
	}

	// 调用服务
	task, err := handler.orchestrator.CreateTask(ctx, createReq)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &humaninloop.TaskResponse{
		Code:    0,
		Message: "success",
		Data:    task,
	})
}

// GetTask 获取任务详情
// @router /api/human-in-loop/tasks/:task_id [GET]
func GetTask(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("task_id")
	if taskID == "" {
		coze.invalidParamRequestResponse(c, "task_id is required")
		return
	}

	task, err := handler.orchestrator.GetTask(ctx, taskID)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &humaninloop.TaskResponse{
		Code:    0,
		Message: "success",
		Data:    task,
	})
}

// ListTasks 获取任务列表
// @router /api/human-in-loop/tasks [GET]
func ListTasks(ctx context.Context, c *app.RequestContext) {
	var req humaninloop.ListTasksRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.invalidParamRequestResponse(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 获取租户ID
	tenantID := coze.GetTenantIDFromContext(ctx)

	// 构建过滤器
	filter := &hilservice.QueueFilter{
		Status:    req.Status,
		Priority:  req.Priority,
		TaskType:  req.TaskType,
		IsOverdue: req.IsOverdue,
		Page:      req.Page,
		PageSize:  req.PageSize,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	// 调用服务
	tasks, total, err := handler.queueService.GetQueue(ctx, tenantID, "", filter)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, &humaninloop.TaskListResponse{
		Code:    0,
		Message: "success",
		Data: &humaninloop.TaskListData{
			Tasks:      tasks,
			Total:      total,
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalPages: totalPages,
		},
	})
}

// AssignTask 分配任务
// @router /api/human-in-loop/tasks/:task_id/assign [POST]
func AssignTask(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("task_id")
	if taskID == "" {
		coze.invalidParamRequestResponse(c, "task_id is required")
		return
	}

	var req humaninloop.AssignTaskRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.invalidParamRequestResponse(c, err.Error())
		return
	}

	err := handler.orchestrator.AssignTask(ctx, taskID, req.AssignTo, req.Reason)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, humaninloop.SuccessResponse("Task assigned successfully"))
}

// ReviewTask 审核任务
// @router /api/human-in-loop/tasks/:task_id/review [POST]
func ReviewTask(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("task_id")
	if taskID == "" {
		coze.invalidParamRequestResponse(c, "task_id is required")
		return
	}

	var req humaninloop.ReviewTaskRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.invalidParamRequestResponse(c, err.Error())
		return
	}

	// 解析人工决策
	var humanDecision map[string]interface{}
	if err := json.Unmarshal(req.HumanDecision, &humanDecision); err != nil {
		coze.invalidParamRequestResponse(c, "invalid human_decision format")
		return
	}

	// 构建提交请求
	resultReq := &hilservice.SubmitResultRequest{
		ReviewerID:      coze.GetUserIDFromContext(ctx),
		Decision:        req.Decision,
		HumanDecision:   humanDecision,
		DecisionReason:  req.DecisionReason,
		ModifiedContent: req.ModifiedContent,
	}

	err := handler.orchestrator.SubmitResult(ctx, taskID, resultReq)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, humaninloop.SuccessResponse("Task reviewed successfully"))
}

// EscalateTask 升级任务
// @router /api/human-in-loop/tasks/:task_id/escalate [POST]
func EscalateTask(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("task_id")
	if taskID == "" {
		coze.invalidParamRequestResponse(c, "task_id is required")
		return
	}

	var req humaninloop.EscalateTaskRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.invalidParamRequestResponse(c, err.Error())
		return
	}

	err := handler.orchestrator.EscalateTask(ctx, taskID, req.EscalateTo, req.Reason)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, humaninloop.SuccessResponse("Task escalated successfully"))
}

// ==================== 历史记录接口 ====================

// GetTaskHistory 获取任务历史
// @router /api/human-in-loop/tasks/:task_id/history [GET]
func GetTaskHistory(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("task_id")
	if taskID == "" {
		coze.invalidParamRequestResponse(c, "task_id is required")
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	var pageInt, pageSizeInt int
	if _, err := fmt.Sscanf(page, "%d", &pageInt); err != nil || pageInt <= 0 {
		pageInt = 1
	}
	if _, err := fmt.Sscanf(pageSize, "%d", &pageSizeInt); err != nil || pageSizeInt <= 0 {
		pageSizeInt = 20
	}

	histories, total, err := handler.protocol.GetHistoryWithPagination(ctx, taskID, pageInt, pageSizeInt)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &humaninloop.HistoryResponse{
		Code:    0,
		Message: "success",
		Data: &humaninloop.HistoryData{
			TaskID:    taskID,
			Histories: histories,
			Total:     int(total),
		},
	})
}

// ==================== 统计分析接口 ====================

// GetQueueStats 获取队列统计
// @router /api/human-in-loop/stats/queue [GET]
func GetQueueStats(ctx context.Context, c *app.RequestContext) {
	tenantID := coze.GetTenantIDFromContext(ctx)

	stats, err := handler.queueService.GetQueueStats(ctx, tenantID)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &humaninloop.QueueStatsResponse{
		Code:    0,
		Message: "success",
		Data:    stats,
	})
}

// MonitorSLA 监控SLA
// @router /api/human-in-loop/monitor/sla [GET]
func MonitorSLA(ctx context.Context, c *app.RequestContext) {
	overdueTasks, err := handler.orchestrator.MonitorSLA(ctx)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为响应格式
	overdueInfos := make([]*humaninloop.OverdueTask, len(overdueTasks))
	for i, task := range overdueTasks {
		overdueInfos[i] = &humaninloop.OverdueTask{
			TaskID:          task.TaskID,
			TenantID:        task.TenantID,
			TaskType:        task.TaskType,
			Priority:        task.Priority,
			Status:          task.Status,
			AssignedTo:      task.AssignedTo,
			OverdueDuration: task.OverdueDuration,
		}
	}

	c.JSON(http.StatusOK, &humaninloop.OverdueTaskResponse{
		Code:    0,
		Message: "success",
		Data:    overdueInfos,
	})
}

// GetAnalytics 获取分析数据
// @router /api/human-in-loop/analytics [GET]
func GetAnalytics(ctx context.Context, c *app.RequestContext) {
	// TODO: 实现分析数据查询
	c.JSON(http.StatusOK, &humaninloop.AnalyticsResponse{
		Code:    0,
		Message: "success",
		Data:    &humaninloop.AnalyticsData{},
	})
}

// GetTaskMetrics 获取任务性能指标
// @router /api/human-in-loop/tasks/:task_id/metrics [GET]
func GetTaskMetrics(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("task_id")
	if taskID == "" {
		coze.invalidParamRequestResponse(c, "task_id is required")
		return
	}

	metrics, err := handler.protocol.CalculateMetrics(ctx, taskID)
	if err != nil {
		coze.internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &humaninloop.BaseResponse{
		Code:    0,
		Message: "success",
		Data:    metrics,
	})
}
