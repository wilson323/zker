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
	"testing"

	intelligenceCommon "github.com/coze-dev/coze-studio/backend/api/model/app/intelligence/common"
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	agentVO "github.com/coze-dev/coze-studio/backend/domain/agent/singleagent/entity/vo"
	searchVO "github.com/coze-dev/coze-studio/backend/domain/search/entity/vo"
	"github.com/stretchr/testify/assert"
)

// TestConvertIntelligenceStatusFromAPI 测试状态转换
func TestConvertIntelligenceStatusFromAPI(t *testing.T) {
	tests := []struct {
		name     string
		input    intelligenceCommon.IntelligenceStatus
		expected searchVO.IntelligenceStatus
	}{
		{
			name:     "Draft status",
			input:    intelligenceCommon.IntelligenceStatus_Draft,
			expected: searchVO.IntelligenceStatus(1),
		},
		{
			name:     "Online status",
			input:    intelligenceCommon.IntelligenceStatus_Online,
			expected: searchVO.IntelligenceStatus(2),
		},
		{
			name:     "Published status",
			input:    intelligenceCommon.IntelligenceStatus_Published,
			expected: searchVO.IntelligenceStatus(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertIntelligenceStatusFromAPI(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestConvertIntelligenceTypeFromAPI 测试类型转换
func TestConvertIntelligenceTypeFromAPI(t *testing.T) {
	tests := []struct {
		name     string
		input    intelligenceCommon.IntelligenceType
		expected searchVO.IntelligenceType
	}{
		{
			name:     "Bot type",
			input:    intelligenceCommon.IntelligenceType_Bot,
			expected: searchVO.IntelligenceType(1),
		},
		{
			name:     "Project type",
			input:    intelligenceCommon.IntelligenceType_Project,
			expected: searchVO.IntelligenceType(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertIntelligenceTypeFromAPI(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestConvertDraftBotDisplayInfoFromAPI 测试Bot显示信息转换 (API → Domain)
func TestConvertDraftBotDisplayInfoFromAPI(t *testing.T) {
	t.Run("nil输入返回nil", func(t *testing.T) {
		result := ConvertDraftBotDisplayInfoFromAPI(nil)
		assert.Nil(t, result)
	})

	t.Run("完整转换", func(t *testing.T) {
		tabStatus := developer_api.TabStatus(1)
		apiData := &developer_api.DraftBotDisplayInfoData{
			TabDisplayInfo: &developer_api.TabDisplayItems{
				PluginTabStatus:  &tabStatus,
				WorkflowTabStatus: &tabStatus,
				KnowledgeTabStatus: &tabStatus,
			},
		}

		result := ConvertDraftBotDisplayInfoFromAPI(apiData)

		assert.NotNil(t, result)
		assert.NotNil(t, result.TabDisplayInfo)
		assert.NotNil(t, result.TabDisplayInfo.PluginTabStatus)
		assert.Equal(t, agentVO.TabStatus(1), *result.TabDisplayInfo.PluginTabStatus)
		assert.Equal(t, agentVO.TabStatus(1), *result.TabDisplayInfo.WorkflowTabStatus)
		assert.Equal(t, agentVO.TabStatus(1), *result.TabDisplayInfo.KnowledgeTabStatus)
	})

	t.Run("所有TabStatus字段转换", func(t *testing.T) {
		status := developer_api.TabStatus(2)
		apiData := &developer_api.DraftBotDisplayInfoData{
			TabDisplayInfo: &developer_api.TabDisplayItems{
				PluginTabStatus:             &status,
				WorkflowTabStatus:           &status,
				KnowledgeTabStatus:          &status,
				DatabaseTabStatus:           &status,
				VariableTabStatus:           &status,
				OpeningDialogTabStatus:      &status,
				ScheduledTaskTabStatus:      &status,
				SuggestionTabStatus:         &status,
				TtsTabStatus:                &status,
				FileboxTabStatus:            &status,
				LongTermMemoryTabStatus:     &status,
				AnswerActionTabStatus:       &status,
				ImageflowTabStatus:          &status,
				BackgroundImageTabStatus:   &status,
				ShortcutTabStatus:           &status,
				KnowledgeTableTabStatus:     &status,
				KnowledgeTextTabStatus:      &status,
				KnowledgePhotoTabStatus:     &status,
				HookInfoTabStatus:           &status,
				DefaultUserInputTabStatus:   &status,
			},
		}

		result := ConvertDraftBotDisplayInfoFromAPI(apiData)

		assert.NotNil(t, result)
		assert.NotNil(t, result.TabDisplayInfo)
		// 验证所有字段都已转换
		assert.NotNil(t, result.TabDisplayInfo.PluginTabStatus)
		assert.NotNil(t, result.TabDisplayInfo.WorkflowTabStatus)
		assert.NotNil(t, result.TabDisplayInfo.KnowledgeTabStatus)
		assert.NotNil(t, result.TabDisplayInfo.DatabaseTabStatus)
		assert.NotNil(t, result.TabDisplayInfo.VariableTabStatus)
		assert.Equal(t, agentVO.TabStatus(2), *result.TabDisplayInfo.PluginTabStatus)
	})
}

// TestConvertDraftBotDisplayInfoToAPI 测试Bot显示信息转换 (Domain → API)
func TestConvertDraftBotDisplayInfoToAPI(t *testing.T) {
	t.Run("nil输入返回nil", func(t *testing.T) {
		result := ConvertDraftBotDisplayInfoToAPI(nil)
		assert.Nil(t, result)
	})

	t.Run("完整转换", func(t *testing.T) {
		tabStatus := agentVO.TabStatus(1)
		domainData := &agentVO.DraftBotDisplayInfoData{
			TabDisplayInfo: &agentVO.TabDisplayItems{
				PluginTabStatus:  &tabStatus,
				WorkflowTabStatus: &tabStatus,
				KnowledgeTabStatus: &tabStatus,
			},
		}

		result := ConvertDraftBotDisplayInfoToAPI(domainData)

		assert.NotNil(t, result)
		assert.NotNil(t, result.TabDisplayInfo)
		assert.NotNil(t, result.TabDisplayInfo.PluginTabStatus)
		assert.Equal(t, developer_api.TabStatus(1), *result.TabDisplayInfo.PluginTabStatus)
		assert.Equal(t, developer_api.TabStatus(1), *result.TabDisplayInfo.WorkflowTabStatus)
		assert.Equal(t, developer_api.TabStatus(1), *result.TabDisplayInfo.KnowledgeTabStatus)
	})

	t.Run("双向转换一致性", func(t *testing.T) {
		// API → Domain → API 应该保持一致
		originalStatus := developer_api.TabStatus(3)
		apiData := &developer_api.DraftBotDisplayInfoData{
			TabDisplayInfo: &developer_api.TabDisplayItems{
				PluginTabStatus: &originalStatus,
			},
		}

		// API → Domain
		domainData := ConvertDraftBotDisplayInfoFromAPI(apiData)
		assert.NotNil(t, domainData)

		// Domain → API
		result := ConvertDraftBotDisplayInfoToAPI(domainData)
		assert.NotNil(t, result)
		assert.NotNil(t, result.TabDisplayInfo)
		assert.NotNil(t, result.TabDisplayInfo.PluginTabStatus)
		assert.Equal(t, originalStatus, *result.TabDisplayInfo.PluginTabStatus)
	})
}

// TestConvertTabStatus 测试TabStatus转换逻辑
func TestConvertTabStatus(t *testing.T) {
	t.Run("nil输入返回nil", func(t *testing.T) {
		// 测试 convertTabStatus (私有函数，通过公开函数间接测试)
		apiData := &developer_api.DraftBotDisplayInfoData{
			TabDisplayInfo: &developer_api.TabDisplayItems{
				PluginTabStatus: nil,
			},
		}

		result := ConvertDraftBotDisplayInfoFromAPI(apiData)
		assert.NotNil(t, result)
		assert.NotNil(t, result.TabDisplayInfo)
		assert.Nil(t, result.TabDisplayInfo.PluginTabStatus)
	})

	t.Run("正确转换TabStatus", func(t *testing.T) {
		statusValues := []developer_api.TabStatus{0, 1, 2, 3, 4, 5}

		for _, status := range statusValues {
			apiData := &developer_api.DraftBotDisplayInfoData{
				TabDisplayInfo: &developer_api.TabDisplayItems{
					PluginTabStatus: &status,
				},
			}

			result := ConvertDraftBotDisplayInfoFromAPI(apiData)
			assert.NotNil(t, result)
			assert.NotNil(t, result.TabDisplayInfo.PluginTabStatus)
			assert.Equal(t, agentVO.TabStatus(status), *result.TabDisplayInfo.PluginTabStatus)
		}
	})
}

// TestTypeConverterNilSafety 测试类型转换器的nil安全性
func TestTypeConverterNilSafety(t *testing.T) {
	t.Run("ConvertDraftBotDisplayInfoFromAPI - nil安全性", func(t *testing.T) {
		result := ConvertDraftBotDisplayInfoFromAPI(nil)
		assert.Nil(t, result)
	})

	t.Run("ConvertDraftBotDisplayInfoToAPI - nil安全性", func(t *testing.T) {
		result := ConvertDraftBotDisplayInfoToAPI(nil)
		assert.Nil(t, result)
	})

	t.Run("嵌套nil TabDisplayInfo", func(t *testing.T) {
		// API层nil TabDisplayInfo
		apiData := &developer_api.DraftBotDisplayInfoData{
			TabDisplayInfo: nil,
		}
		result := ConvertDraftBotDisplayInfoFromAPI(apiData)
		assert.NotNil(t, result)
		assert.Nil(t, result.TabDisplayInfo)

		// Domain层nil TabDisplayInfo
		domainData := &agentVO.DraftBotDisplayInfoData{
			TabDisplayInfo: nil,
		}
		result = ConvertDraftBotDisplayInfoToAPI(domainData)
		assert.NotNil(t, result)
		assert.Nil(t, result.TabDisplayInfo)
	})
}

// TestTypeConverterConsistency 测试类型转换器的一致性
func TestTypeConverterConsistency(t *testing.T) {
	t.Run("多次转换结果一致", func(t *testing.T) {
		status := intelligenceCommon.IntelligenceStatus_Published
		expected := searchVO.IntelligenceStatus(status)

		// 多次调用应该返回相同结果
		result1 := ConvertIntelligenceStatusFromAPI(status)
		result2 := ConvertIntelligenceStatusFromAPI(status)
		result3 := ConvertIntelligenceStatusFromAPI(status)

		assert.Equal(t, expected, result1)
		assert.Equal(t, result1, result2)
		assert.Equal(t, result2, result3)
	})

	t.Run("类型转换一致性", func(t *testing.T) {
		typ := intelligenceCommon.IntelligenceType_Bot
		expected := searchVO.IntelligenceType(typ)

		result1 := ConvertIntelligenceTypeFromAPI(typ)
		result2 := ConvertIntelligenceTypeFromAPI(typ)

		assert.Equal(t, expected, result1)
		assert.Equal(t, result1, result2)
	})
}

// BenchmarkConvertIntelligenceStatus 基准测试状态转换
func BenchmarkConvertIntelligenceStatus(b *testing.B) {
	status := intelligenceCommon.IntelligenceStatus_Published

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvertIntelligenceStatusFromAPI(status)
	}
}

// BenchmarkConvertDraftBotDisplayInfo 基准测试Bot显示信息转换
func BenchmarkConvertDraftBotDisplayInfo(b *testing.B) {
	tabStatus := developer_api.TabStatus(1)
	apiData := &developer_api.DraftBotDisplayInfoData{
		TabDisplayInfo: &developer_api.TabDisplayItems{
			PluginTabStatus:    &tabStatus,
			WorkflowTabStatus:  &tabStatus,
			KnowledgeTabStatus: &tabStatus,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvertDraftBotDisplayInfoFromAPI(apiData)
	}
}
