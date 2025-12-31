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
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/dal"
)

// NewProfileRepository 创建员工画像Repository
func NewProfileRepository(db *gorm.DB) ProfileRepository {
	// dal.ProfileDAL 隐式实现了 repository.ProfileRepository
	return dal.NewProfileDAL(db)
}

// NewTaskRepository 创建任务分配Repository
func NewTaskRepository(db *gorm.DB) TaskAssignmentRepository {
	return dal.NewTaskAssignmentDAL(db)
}

// NewPerformanceRepository 创建绩效Repository
func NewPerformanceRepository(db *gorm.DB) PerformanceRepository {
	return dal.NewPerformanceDAL(db)
}

// NewRepositories 创建所有repository实例
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Profile:     NewProfileRepository(db),
		Task:        NewTaskRepository(db),
		Performance: NewPerformanceRepository(db),
	}
}

// Repositories repository集合
type Repositories struct {
	Profile     ProfileRepository
	Task        TaskAssignmentRepository
	Performance PerformanceRepository
}
