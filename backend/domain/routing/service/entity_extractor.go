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
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"go.uber.org/zap"
)

// EntityExtractor 实体提取器
// 职责：从用户输入中提取关键实体（如地点、时间、人物等）
type EntityExtractor struct {
	llmClient LLMClient
	logger    *zap.Logger

	// 提取策略
	strategy string // llm, rule, hybrid
}

// NewEntityExtractor 创建实体提取器实例
func NewEntityExtractor(llmClient LLMClient, logger *zap.Logger) *EntityExtractor {
	return &EntityExtractor{
		llmClient: llmClient,
		logger:    logger,
		strategy:  "hybrid", // 默认使用混合策略
	}
}

// ExtractEntities 提取实体
func (e *EntityExtractor) ExtractEntities(ctx context.Context, text string) ([]*entity.Entity, error) {
	if text == "" {
		return []*entity.Entity{}, nil
	}

	switch e.strategy {
	case "llm":
		return e.extractByLLM(ctx, text)
	case "rule":
		return e.extractByRule(ctx, text)
	case "hybrid":
		return e.extractByHybrid(ctx, text)
	default:
		return nil, errorx.New(errno.ErrRouteFormatInvalidCode, errorx.KV("reason", "unknown extraction strategy"))
	}
}

// extractByLLM 使用LLM提取实体
func (e *EntityExtractor) extractByLLM(ctx context.Context, text string) ([]*entity.Entity, error) {
	prompt := e.buildLLMPrompt(text)

	response, err := e.llmClient.GenerateText(ctx, prompt)
	if err != nil {
		return nil, errorx.New(errno.ErrEntityExtractionFailedCode, errorx.KV("error", err.Error()))
	}

	var llmResult struct {
		Entities []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
			Start int    `json:"start"`
			End   int    `json:"end"`
		} `json:"entities"`
	}

	if err := json.Unmarshal([]byte(response), &llmResult); err != nil {
		return nil, errorx.New(errno.ErrEntityExtractionFailedCode, errorx.KV("error", "failed to parse llm response"))
	}

	entities := make([]*entity.Entity, 0, len(llmResult.Entities))
	for _, e := range llmResult.Entities {
		entities = append(entities, &entity.Entity{
			EntityName:    e.Name,
			EntityValue:   e.Value,
			Confidence:    0.9, // LLM提取的置信度默认较高
			StartPosition: e.Start,
			EndPosition:   e.End,
		})
	}

	return entities, nil
}

// extractByRule 使用规则提取实体
func (e *EntityExtractor) extractByRule(ctx context.Context, text string) ([]*entity.Entity, error) {
	entities := make([]*entity.Entity, 0)

	// 1. 提取地点实体
	locations := e.extractLocation(text)
	entities = append(entities, locations...)

	// 2. 提取时间实体
	times := e.extractTime(text)
	entities = append(entities, times...)

	// 3. 提取数字实体
	numbers := e.extractNumber(text)
	entities = append(entities, numbers...)

	// 4. 提取动作实体
	actions := e.extractAction(text)
	entities = append(entities, actions...)

	return entities, nil
}

// extractByHybrid 使用混合策略提取实体
func (e *EntityExtractor) extractByHybrid(ctx context.Context, text string) ([]*entity.Entity, error) {
	// 先使用规则提取
	ruleEntities, err := e.extractByRule(ctx, text)
	if err != nil {
		e.logger.Warn("rule extraction failed", zap.Error(err))
		ruleEntities = []*entity.Entity{}
	}

	// 如果规则提取到实体，直接返回
	if len(ruleEntities) > 0 {
		return ruleEntities, nil
	}

	// 否则使用LLM提取
	return e.extractByLLM(ctx, text)
}

// buildLLMPrompt 构建LLM提示词
func (e *EntityExtractor) buildLLMPrompt(text string) string {
	prompt := fmt.Sprintf(`你是一个实体提取助手。请从用户输入中提取关键实体。

用户输入: %s

请提取以下类型的实体：
- location: 地点（如北京、上海）
- time: 时间（如今天、明天、2025-01-01）
- person: 人物（如张三、李四）
- number: 数字（如10、100）
- action: 动作（如查询、预订）
- target: 目标（如天气、机票）

请以JSON格式返回提取结果:
{
  "entities": [
    {
      "name": "实体类型",
      "value": "实体值",
      "start": 0,
      "end": 10
    }
  ]
}
`, text)

	return prompt
}

// extractLocation 提取地点实体
func (e *EntityExtractor) extractLocation(text string) []*entity.Entity {
	// 简化实现：使用常见城市名列表
	cities := []string{
		"北京", "上海", "广州", "深圳", "杭州", "南京", "成都", "重庆",
		"武汉", "西安", "天津", "苏州", "长沙", "郑州", "青岛", "大连",
		"Beijing", "Shanghai", "Guangzhou", "Shenzhen", "Hangzhou",
	}

	entities := make([]*entity.Entity, 0)
	for _, city := range cities {
		idx := indexOf(text, city)
		if idx != -1 {
			entities = append(entities, &entity.Entity{
				EntityName:    "location",
				EntityValue:   city,
				Confidence:    0.95,
				StartPosition: idx,
				EndPosition:   idx + len([]rune(city)),
			})
		}
	}

	return entities
}

// extractTime 提取时间实体
func (e *EntityExtractor) extractTime(text string) []*entity.Entity {
	entities := make([]*entity.Entity, 0)

	// 匹配相对时间
	timePatterns := []struct {
		pattern string
		value   string
	}{
		{`今天`, "today"},
		{`明天`, "tomorrow"},
		{`后天`, "day_after_tomorrow"},
		{`昨天`, "yesterday"},
		{`本周`, "this_week"},
		{`下周`, "next_week"},
		{`本月`, "this_month"},
		{`下月`, "next_month"},
	}

	for _, tp := range timePatterns {
		re := regexp.MustCompile(tp.pattern)
		matches := re.FindStringIndex(text)
		if matches != nil {
			entities = append(entities, &entity.Entity{
				EntityName:    "time",
				EntityValue:   tp.value,
				Confidence:    0.9,
				StartPosition: matches[0],
				EndPosition:   matches[1],
			})
		}
	}

	// 匹配日期格式（YYYY-MM-DD）
	datePattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	matches := datePattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		entities = append(entities, &entity.Entity{
			EntityName:    "time",
			EntityValue:   text[match[0]:match[1]],
			Confidence:    0.95,
			StartPosition: match[0],
			EndPosition:   match[1],
		})
	}

	return entities
}

// extractNumber 提取数字实体
func (e *EntityExtractor) extractNumber(text string) []*entity.Entity {
	entities := make([]*entity.Entity, 0)

	// 匹配整数和浮点数
	pattern := regexp.MustCompile(`\d+(\.\d+)?`)
	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		entities = append(entities, &entity.Entity{
			EntityName:    "number",
			EntityValue:   text[match[0]:match[1]],
			Confidence:    0.95,
			StartPosition: match[0],
			EndPosition:   match[1],
		})
	}

	return entities
}

// extractAction 提取动作实体
func (e *EntityExtractor) extractAction(text string) []*entity.Entity {
	actions := []string{
		"查询", "搜索", "找", "查找",
		"预订", "预定", "买", "订购",
		"取消", "退订", "删除",
		"修改", "更改", "更新",
		"发送", "传递", "转发",
		"create", "update", "delete", "search", "query",
	}

	entities := make([]*entity.Entity, 0)
	for _, action := range actions {
		idx := indexOf(text, action)
		if idx != -1 {
			entities = append(entities, &entity.Entity{
				EntityName:    "action",
				EntityValue:   action,
				Confidence:    0.9,
				StartPosition: idx,
				EndPosition:   idx + len([]rune(action)),
			})
		}
	}

	return entities
}

// indexOf 查找子字符串位置（支持中文）
func indexOf(s, substr string) int {
	runes := []rune(s)
	subrunes := []rune(substr)

	for i := 0; i <= len(runes)-len(subrunes); i++ {
		match := true
		for j := 0; j < len(subrunes); j++ {
			if runes[i+j] != subrunes[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}

	return -1
}

// BatchExtractEntities 批量提取实体
func (e *EntityExtractor) BatchExtractEntities(ctx context.Context, texts []string) ([][]*entity.Entity, error) {
	results := make([][]*entity.Entity, 0, len(texts))

	for _, text := range texts {
		entities, err := e.ExtractEntities(ctx, text)
		if err != nil {
			e.logger.Warn("failed to extract entities",
				zap.String("text", text),
				zap.Error(err))
			entities = []*entity.Entity{}
		}
		results = append(results, entities)
	}

	return results, nil
}
