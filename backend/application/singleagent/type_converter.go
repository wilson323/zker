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

package singleagent

import (
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	intelligenceCommon "github.com/coze-dev/coze-studio/backend/api/model/app/intelligence/common"
	agentVO "github.com/coze-dev/coze-studio/backend/domain/agent/singleagent/entity/vo"
	searchVO "github.com/coze-dev/coze-studio/backend/domain/search/entity/vo"
)

// ConvertIntelligenceStatusFromAPI converts API layer IntelligenceStatus to domain VO IntelligenceStatus
func ConvertIntelligenceStatusFromAPI(status intelligenceCommon.IntelligenceStatus) searchVO.IntelligenceStatus {
	return searchVO.IntelligenceStatus(status)
}

// ConvertIntelligenceTypeFromAPI converts API layer IntelligenceType to domain VO IntelligenceType
func ConvertIntelligenceTypeFromAPI(typ intelligenceCommon.IntelligenceType) searchVO.IntelligenceType {
	return searchVO.IntelligenceType(typ)
}

// ConvertDraftBotDisplayInfoFromAPI converts API layer DraftBotDisplayInfoData to domain VO DraftBotDisplayInfoData
func ConvertDraftBotDisplayInfoFromAPI(apiData *developer_api.DraftBotDisplayInfoData) *agentVO.DraftBotDisplayInfoData {
	if apiData == nil {
		return nil
	}

	return &agentVO.DraftBotDisplayInfoData{
		TabDisplayInfo: convertTabDisplayItemsFromAPI(apiData.TabDisplayInfo),
	}
}

// convertTabDisplayItemsFromAPI converts API layer TabDisplayItems to domain VO TabDisplayItems
func convertTabDisplayItemsFromAPI(apiItems *developer_api.TabDisplayItems) *agentVO.TabDisplayItems {
	if apiItems == nil {
		return nil
	}

	return &agentVO.TabDisplayItems{
		PluginTabStatus:           convertTabStatus(apiItems.PluginTabStatus),
		WorkflowTabStatus:         convertTabStatus(apiItems.WorkflowTabStatus),
		KnowledgeTabStatus:        convertTabStatus(apiItems.KnowledgeTabStatus),
		DatabaseTabStatus:         convertTabStatus(apiItems.DatabaseTabStatus),
		VariableTabStatus:         convertTabStatus(apiItems.VariableTabStatus),
		OpeningDialogTabStatus:    convertTabStatus(apiItems.OpeningDialogTabStatus),
		ScheduledTaskTabStatus:    convertTabStatus(apiItems.ScheduledTaskTabStatus),
		SuggestionTabStatus:       convertTabStatus(apiItems.SuggestionTabStatus),
		TtsTabStatus:              convertTabStatus(apiItems.TtsTabStatus),
		FileboxTabStatus:          convertTabStatus(apiItems.FileboxTabStatus),
		LongTermMemoryTabStatus:   convertTabStatus(apiItems.LongTermMemoryTabStatus),
		AnswerActionTabStatus:     convertTabStatus(apiItems.AnswerActionTabStatus),
		ImageflowTabStatus:        convertTabStatus(apiItems.ImageflowTabStatus),
		BackgroundImageTabStatus:  convertTabStatus(apiItems.BackgroundImageTabStatus),
		ShortcutTabStatus:         convertTabStatus(apiItems.ShortcutTabStatus),
		KnowledgeTableTabStatus:   convertTabStatus(apiItems.KnowledgeTableTabStatus),
		KnowledgeTextTabStatus:    convertTabStatus(apiItems.KnowledgeTextTabStatus),
		KnowledgePhotoTabStatus:   convertTabStatus(apiItems.KnowledgePhotoTabStatus),
		HookInfoTabStatus:         convertTabStatus(apiItems.HookInfoTabStatus),
		DefaultUserInputTabStatus: convertTabStatus(apiItems.DefaultUserInputTabStatus),
	}
}

// convertTabStatus converts API layer TabStatus to domain VO TabStatus
func convertTabStatus(status *developer_api.TabStatus) *agentVO.TabStatus {
	if status == nil {
		return nil
	}
	s := agentVO.TabStatus(*status)
	return &s
}

// ConvertDraftBotDisplayInfoToAPI converts domain VO DraftBotDisplayInfoData to API layer DraftBotDisplayInfoData
func ConvertDraftBotDisplayInfoToAPI(domainData *agentVO.DraftBotDisplayInfoData) *developer_api.DraftBotDisplayInfoData {
	if domainData == nil {
		return nil
	}

	return &developer_api.DraftBotDisplayInfoData{
		TabDisplayInfo: convertTabDisplayItemsToAPI(domainData.TabDisplayInfo),
	}
}

// convertTabDisplayItemsToAPI converts domain VO TabDisplayItems to API layer TabDisplayItems
func convertTabDisplayItemsToAPI(domainItems *agentVO.TabDisplayItems) *developer_api.TabDisplayItems {
	if domainItems == nil {
		return nil
	}

	return &developer_api.TabDisplayItems{
		PluginTabStatus:           convertTabStatusToAPI(domainItems.PluginTabStatus),
		WorkflowTabStatus:         convertTabStatusToAPI(domainItems.WorkflowTabStatus),
		KnowledgeTabStatus:        convertTabStatusToAPI(domainItems.KnowledgeTabStatus),
		DatabaseTabStatus:         convertTabStatusToAPI(domainItems.DatabaseTabStatus),
		VariableTabStatus:         convertTabStatusToAPI(domainItems.VariableTabStatus),
		OpeningDialogTabStatus:    convertTabStatusToAPI(domainItems.OpeningDialogTabStatus),
		ScheduledTaskTabStatus:    convertTabStatusToAPI(domainItems.ScheduledTaskTabStatus),
		SuggestionTabStatus:       convertTabStatusToAPI(domainItems.SuggestionTabStatus),
		TtsTabStatus:              convertTabStatusToAPI(domainItems.TtsTabStatus),
		FileboxTabStatus:          convertTabStatusToAPI(domainItems.FileboxTabStatus),
		LongTermMemoryTabStatus:   convertTabStatusToAPI(domainItems.LongTermMemoryTabStatus),
		AnswerActionTabStatus:     convertTabStatusToAPI(domainItems.AnswerActionTabStatus),
		ImageflowTabStatus:        convertTabStatusToAPI(domainItems.ImageflowTabStatus),
		BackgroundImageTabStatus:  convertTabStatusToAPI(domainItems.BackgroundImageTabStatus),
		ShortcutTabStatus:         convertTabStatusToAPI(domainItems.ShortcutTabStatus),
		KnowledgeTableTabStatus:   convertTabStatusToAPI(domainItems.KnowledgeTableTabStatus),
		KnowledgeTextTabStatus:    convertTabStatusToAPI(domainItems.KnowledgeTextTabStatus),
		KnowledgePhotoTabStatus:   convertTabStatusToAPI(domainItems.KnowledgePhotoTabStatus),
		HookInfoTabStatus:         convertTabStatusToAPI(domainItems.HookInfoTabStatus),
		DefaultUserInputTabStatus: convertTabStatusToAPI(domainItems.DefaultUserInputTabStatus),
	}
}

// convertTabStatusToAPI converts domain VO TabStatus to API layer TabStatus
func convertTabStatusToAPI(status *agentVO.TabStatus) *developer_api.TabStatus {
	if status == nil {
		return nil
	}
	s := developer_api.TabStatus(*status)
	return &s
}
