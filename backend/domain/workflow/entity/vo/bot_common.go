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

// PluginFrom 插件来源（领域值对象）
// 从 api/model/app/bot_common.PluginFrom 迁移过来
type PluginFrom int64

const (
	PluginFrom_Default  PluginFrom = 0
	PluginFrom_FromSaas PluginFrom = 1
)

// Validate 验证值对象
func (p PluginFrom) Validate() error {
	if p != PluginFrom_Default && p != PluginFrom_FromSaas {
		return ErrInvalidPluginFrom
	}
	return nil
}
