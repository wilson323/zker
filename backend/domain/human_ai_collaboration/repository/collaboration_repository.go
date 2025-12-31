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

	"github.com/coze-dev/coze-studio/backend/domain/human_ai_collaboration/entity"
)

// CollaborationRepository 人机协同仓储接口
type CollaborationRepository interface {
	// CreateSession 创建协同会话
	CreateSession(ctx context.Context, session *entity.CollaborationSession) error

	// GetSessionByID 根据ID获取协同会话
	GetSessionByID(ctx context.Context, sessionID string) (*entity.CollaborationSession, error)

	// GetPendingSessionsByUser 获取用户的待处理会话
	GetPendingSessionsByUser(ctx context.Context, tenantID, userID string, limit int) ([]*entity.CollaborationSession, error)

	// GetSessionsByTimeRange 获取时间范围内的会话
	GetSessionsByTimeRange(ctx context.Context, tenantID string, startTime, endTime time.Time) ([]*entity.CollaborationSession, error)

	// UpdateSession 更新协同会话
	UpdateSession(ctx context.Context, session *entity.CollaborationSession) error

	// DeleteSession 删除协同会话
	DeleteSession(ctx context.Context, sessionID string) error

	// CreateLog 创建协同日志
	CreateLog(ctx context.Context, log *entity.CollaborationLog) error

	// GetLogsBySessionID 获取会话的所有日志
	GetLogsBySessionID(ctx context.Context, sessionID string) ([]*entity.CollaborationLog, error)
}
