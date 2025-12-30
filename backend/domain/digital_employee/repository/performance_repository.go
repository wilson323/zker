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

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/entity"
)

// PerformanceRepository 绩效仓储接口
type PerformanceRepository interface {
	// GetOrCreate 获取或创建绩效记录
	GetOrCreate(ctx context.Context, employeeID string, tenantID string, period entity.PerformancePeriod, date string) (*entity.EmployeePerformance, error)

	// Update 更新绩效记录
	Update(ctx context.Context, performance *entity.EmployeePerformance) error

	// GetByEmployeeID 获取员工绩效列表
	GetByEmployeeID(ctx context.Context, employeeID string, period entity.PerformancePeriod, page, pageSize int) ([]*entity.EmployeePerformance, int64, error)

	// GetByTenantID 获取租户下所有员工绩效
	GetByTenantID(ctx context.Context, tenantID string, period entity.PerformancePeriod, page, pageSize int) ([]*entity.EmployeePerformance, int64, error)

	// BatchUpdateStats 批量更新统计数据（定时任务调用）
	BatchUpdateStats(ctx context.Context) error
}
