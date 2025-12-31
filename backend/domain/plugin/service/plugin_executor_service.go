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

	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
)

// PluginExecutorService 插件执行引擎服务接口
type PluginExecutorService interface {
	// Execute 执行插件(异步)
	Execute(
		ctx context.Context,
		req *entity.ExecutePluginRequest,
	) (*entity.ExecutePluginResponse, error)

	// ExecuteSync 同步执行插件
	ExecuteSync(
		ctx context.Context,
		req *entity.ExecutePluginRequest,
	) (*entity.ExecutePluginResponse, error)

	// GetExecution 获取执行记录
	GetExecution(
		ctx context.Context,
		req *entity.GetExecutionRequest,
	) (*entity.GetExecutionResponse, error)

	// CancelExecution 取消执行
	CancelExecution(
		ctx context.Context,
		executionID string,
	) error

	// BatchExecute 批量执行插件
	BatchExecute(
		ctx context.Context,
		reqs []*entity.ExecutePluginRequest,
	) ([]*entity.ExecutePluginResponse, error)

	// RetryExecution 重试执行
	RetryExecution(
		ctx context.Context,
		executionID string,
	) (*entity.ExecutePluginResponse, error)
}
