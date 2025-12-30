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

// ProfileRepository 员工画像仓储接口
type ProfileRepository interface {
	// Create 创建员工画像
	Create(ctx context.Context, profile *entity.EmployeeProfile) error

	// Update 更新员工画像
	Update(ctx context.Context, profile *entity.EmployeeProfile) error

	// GetByID 根据ID获取员工画像
	GetByID(ctx context.Context, employeeID string) (*entity.EmployeeProfile, error)

	// GetByTenantID 根据租户ID列出员工画像
	GetByTenantID(ctx context.Context, tenantID string, role *entity.EmployeeRole, status *entity.EmployeeStatus, page, pageSize int) ([]*entity.EmployeeProfile, int64, error)

	// Delete 删除员工画像（软删除）
	Delete(ctx context.Context, employeeID string) error

	// GetByBotID 根据Bot ID获取员工画像
	GetByBotID(ctx context.Context, botID string) (*entity.EmployeeProfile, error)

	// ExistsByName 检查名称是否存在（同一租户下）
	ExistsByName(ctx context.Context, tenantID, name string) (bool, error)
}
