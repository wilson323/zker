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
	"fmt"
)

// ModelMatcher 模型意图匹配器
type ModelMatcher struct {
	enabled    bool
	modelEndpoint string // 模型API端点
	apiKey      string
}

// NewModelMatcher 创建模型匹配器
func NewModelMatcher() *ModelMatcher {
	return &ModelMatcher{
		enabled: false, // 默认禁用，需要配置模型API
	}
}

// Match 模型匹配
func (m *ModelMatcher) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
	if !m.enabled {
		return nil, fmt.Errorf("model matcher is not enabled")
	}

	// TODO: 实现模型调用
	// 1. 构造模型请求
	// 2. 调用模型API进行意图识别
	// 3. 解析模型响应
	// 4. 转换为MatchOutput格式

	return nil, fmt.Errorf("model matcher not implemented yet")
}

// Enable 启用模型匹配器
func (m *ModelMatcher) Enable(endpoint, apiKey string) {
	m.enabled = true
	m.modelEndpoint = endpoint
	m.apiKey = apiKey
}

// Disable 禁用模型匹配器
func (m *ModelMatcher) Disable() {
	m.enabled = false
}

// IsEnabled 检查是否启用
func (m *ModelMatcher) IsEnabled() bool {
	return m.enabled
}
