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

package search

import (
	intelligenceCommon "github.com/coze-dev/coze-studio/backend/api/model/app/intelligence/common"
	workflowAPI "github.com/coze-dev/coze-studio/backend/api/model/workflow"
	searchVO "github.com/coze-dev/coze-studio/backend/domain/search/entity/vo"
	workflowVO "github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
)

// ConvertIntelligenceTypeFromAPI converts API layer IntelligenceType to domain VO IntelligenceType
func ConvertIntelligenceTypeFromAPI(typ intelligenceCommon.IntelligenceType) searchVO.IntelligenceType {
	return searchVO.IntelligenceType(typ)
}

// ConvertIntelligenceTypeToAPI converts domain VO IntelligenceType to API layer IntelligenceType
func ConvertIntelligenceTypeToAPI(typ searchVO.IntelligenceType) intelligenceCommon.IntelligenceType {
	return intelligenceCommon.IntelligenceType(typ)
}

// ConvertIntelligenceStatusFromAPI converts API layer IntelligenceStatus to domain VO IntelligenceStatus
func ConvertIntelligenceStatusFromAPI(status intelligenceCommon.IntelligenceStatus) searchVO.IntelligenceStatus {
	return searchVO.IntelligenceStatus(status)
}

// ConvertIntelligenceStatusToAPI converts domain VO IntelligenceStatus to API layer IntelligenceStatus
func ConvertIntelligenceStatusToAPI(status searchVO.IntelligenceStatus) intelligenceCommon.IntelligenceStatus {
	return intelligenceCommon.IntelligenceStatus(status)
}

// ConvertWorkflowModeFromAPI converts API layer WorkflowMode to domain VO WorkflowMode
func ConvertWorkflowModeFromAPI(mode workflowAPI.WorkflowMode) workflowVO.WorkflowMode {
	return workflowVO.WorkflowMode(mode)
}

// ConvertWorkflowModeToAPI converts domain VO WorkflowMode to API layer WorkflowMode
func ConvertWorkflowModeToAPI(mode workflowVO.WorkflowMode) workflowAPI.WorkflowMode {
	return workflowAPI.WorkflowMode(mode)
}
