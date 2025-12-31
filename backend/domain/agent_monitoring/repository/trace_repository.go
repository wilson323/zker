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

	"github.com/coze-dev/coze-studio/backend/domain/agent_monitoring/entity"
)

// TraceRepository 链路追踪仓储接口
type TraceRepository interface {
	// Create 创建链路追踪记录
	Create(ctx context.Context, trace *entity.Trace) error

	// GetByID 根据ID获取链路追踪
	GetByID(ctx context.Context, traceID string) (*entity.Trace, error)

	// Query 查询链路追踪列表
	Query(ctx context.Context, filter *entity.TraceFilter) ([]*entity.Trace, int64, error)

	// GetSpansByTraceID 获取链路的所有跨度
	GetSpansByTraceID(ctx context.Context, traceID string) ([]*entity.TraceSpan, error)

	// Delete 删除链路追踪
	Delete(ctx context.Context, traceID string) error

	// BatchDelete 批量删除链路追踪
	BatchDelete(ctx context.Context, traceIDs []string) error
}
