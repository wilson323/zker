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

// TabStatus 标签页状态
type TabStatus int64

const (
	TabStatus_Default TabStatus = 0
	TabStatus_Open    TabStatus = 1
	TabStatus_Close   TabStatus = 2
	TabStatus_Hide    TabStatus = 3
)

// TabDisplayItems 标签页显示项
// 从 api/model/app/developer_api 迁移到领域层
type TabDisplayItems struct {
	PluginTabStatus           *TabStatus `json:"plugin_tab_status,omitempty"`
	WorkflowTabStatus         *TabStatus `json:"workflow_tab_status,omitempty"`
	KnowledgeTabStatus        *TabStatus `json:"knowledge_tab_status,omitempty"`
	DatabaseTabStatus         *TabStatus `json:"database_tab_status,omitempty"`
	VariableTabStatus         *TabStatus `json:"variable_tab_status,omitempty"`
	OpeningDialogTabStatus    *TabStatus `json:"opening_dialog_tab_status,omitempty"`
	ScheduledTaskTabStatus    *TabStatus `json:"scheduled_task_tab_status,omitempty"`
	SuggestionTabStatus       *TabStatus `json:"suggestion_tab_status,omitempty"`
	TtsTabStatus              *TabStatus `json:"tts_tab_status,omitempty"`
	FileboxTabStatus          *TabStatus `json:"filebox_tab_status,omitempty"`
	LongTermMemoryTabStatus   *TabStatus `json:"long_term_memory_tab_status,omitempty"`
	AnswerActionTabStatus     *TabStatus `json:"answer_action_tab_status,omitempty"`
	ImageflowTabStatus        *TabStatus `json:"imageflow_tab_status,omitempty"`
	BackgroundImageTabStatus  *TabStatus `json:"background_image_tab_status,omitempty"`
	ShortcutTabStatus         *TabStatus `json:"shortcut_tab_status,omitempty"`
	KnowledgeTableTabStatus   *TabStatus `json:"knowledge_table_tab_status,omitempty"`
	KnowledgeTextTabStatus    *TabStatus `json:"knowledge_text_tab_status,omitempty"`
	KnowledgePhotoTabStatus   *TabStatus `json:"knowledge_photo_tab_status,omitempty"`
	HookInfoTabStatus         *TabStatus `json:"hook_info_tab_status,omitempty"`
	DefaultUserInputTabStatus *TabStatus `json:"default_user_input_tab_status,omitempty"`
}

// DraftBotDisplayInfoData Draft Bot显示信息数据（领域值对象）
// 从 api/model/app/developer_api.DraftBotDisplayInfoData 迁移过来
type DraftBotDisplayInfoData struct {
	TabDisplayInfo *TabDisplayItems `json:"tab_display_info,omitempty"`
}

// Validate 验证值对象
func (d *DraftBotDisplayInfoData) Validate() error {
	// TODO: 添加验证逻辑
	return nil
}
