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

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
)

// NotificationService 通知服务接口
type NotificationService interface {
	// NotifyTaskAssigned 通知任务分配
	NotifyTaskAssigned(ctx context.Context, task *entity.CollaborationTask, assigneeID string)

	// NotifyTaskCompleted 通知任务完成
	NotifyTaskCompleted(ctx context.Context, task *entity.CollaborationTask, decision string)

	// NotifyTaskEscalated 通知任务升级
	NotifyTaskEscalated(ctx context.Context, task *entity.CollaborationTask, escalateTo string, reason string)

	// NotifySLAWarning 通知SLA预警
	NotifySLAWarning(ctx context.Context, task *entity.CollaborationTask)
}

// notificationService 通知服务实现
type notificationService struct {
	// 可以集成邮件、短信、WebSocket等通知渠道
}

// NewNotificationService 创建通知服务
func NewNotificationService() NotificationService {
	return &notificationService{}
}

// NotifyTaskAssigned 通知任务分配
func (s *notificationService) NotifyTaskAssigned(ctx context.Context, task *entity.CollaborationTask, assigneeID string) {
	// TODO: 实现通知逻辑
	// 1. 查询审核人的通知偏好
	// 2. 根据偏好发送通知（邮件、短信、WebSocket等）
}

// NotifyTaskCompleted 通知任务完成
func (s *notificationService) NotifyTaskCompleted(ctx context.Context, task *entity.CollaborationTask, decision string) {
	// TODO: 实现通知逻辑
}

// NotifyTaskEscalated 通知任务升级
func (s *notificationService) NotifyTaskEscalated(ctx context.Context, task *entity.CollaborationTask, escalateTo string, reason string) {
	// TODO: 实现通知逻辑
}

// NotifySLAWarning 通知SLA预警
func (s *notificationService) NotifySLAWarning(ctx context.Context, task *entity.CollaborationTask) {
	// TODO: 实现通知逻辑
}
