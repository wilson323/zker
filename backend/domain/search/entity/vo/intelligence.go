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

package vo

// IntelligenceStatus 智能体状态（领域值对象）
// 从 api/model/app/intelligence/common.IntelligenceStatus 迁移过来
type IntelligenceStatus int64

const (
	IntelligenceStatus_Using      IntelligenceStatus = 1
	IntelligenceStatus_Deleted    IntelligenceStatus = 2
	IntelligenceStatus_Banned     IntelligenceStatus = 3
	IntelligenceStatus_MoveFailed IntelligenceStatus = 4
	IntelligenceStatus_Copying    IntelligenceStatus = 5
	IntelligenceStatus_CopyFailed IntelligenceStatus = 6
)

// IntelligenceType 智能体类型（领域值对象）
// 从 api/model/app/intelligence/common.IntelligenceType 迁移过来
type IntelligenceType int64

const (
	IntelligenceType_Bot     IntelligenceType = 1
	IntelligenceType_Project IntelligenceType = 2
)

// Validate 验证值对象
func (t IntelligenceType) Validate() error {
	if t != IntelligenceType_Bot && t != IntelligenceType_Project {
		return ErrInvalidIntelligenceType
	}
	return nil
}

// Validate 验证值对象
func (s IntelligenceStatus) Validate() error {
	if s < IntelligenceStatus_Using || s > IntelligenceStatus_CopyFailed {
		return ErrInvalidIntelligenceStatus
	}
	return nil
}
