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
	"reflect"
	"strings"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// CustomFilterEngine 自定义过滤器引擎
type CustomFilterEngine struct {
	// 可集成表达式引擎，如 govaluate
}

// FilterCondition 过滤条件
type FilterCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // eq, ne, gt, lt, gte, lte, in, contains, starts_with, ends_with
	Value    interface{} `json:"value"`
	Logic    string      `json:"logic"` // AND, OR
}

// NewCustomFilterEngine 创建自定义过滤器引擎实例
func NewCustomFilterEngine() *CustomFilterEngine {
	return &CustomFilterEngine{}
}

// Match 检查资源是否匹配自定义过滤条件
func (e *CustomFilterEngine) Match(ctx context.Context, customFilterJSON string, resource map[string]interface{}) (bool, error) {
	if customFilterJSON == "" {
		return true, nil // 空过滤器默认匹配所有
	}

	var conditions []FilterCondition
	if err := json.Unmarshal([]byte(customFilterJSON), &conditions); err != nil {
		return false, errorx.WrapByCode(err, errno.ErrPermissionInvalidParamCode,
                errorx.KV("operation", "parse custom filter"),
            )
	}

	return e.evaluateConditions(resource, conditions), nil
}

// MatchWithJSON 检查JSON格式的资源是否匹配
func (e *CustomFilterEngine) MatchWithJSON(ctx context.Context, customFilterJSON string, resourceJSON string) (bool, error) {
	if customFilterJSON == "" {
		return true, nil
	}

	var resource map[string]interface{}
	if err := json.Unmarshal([]byte(resourceJSON), &resource); err != nil {
		return false, errorx.WrapByCode(err, errno.ErrPermissionInvalidParamCode,
                errorx.KV("operation", "parse resource JSON"),
            )
	}

	return e.Match(ctx, customFilterJSON, resource)
}

// evaluateConditions 评估条件列表
func (e *CustomFilterEngine) evaluateConditions(resource map[string]interface{}, conditions []FilterCondition) bool {
	if len(conditions) == 0 {
		return true
	}

	result := true
	currentLogic := "AND"

	for _, condition := range conditions {
		matched := e.evaluateCondition(resource, condition)

		// 更新逻辑操作符
		if condition.Logic != "" {
			currentLogic = strings.ToUpper(condition.Logic)
		}

		// 根据逻辑操作符合并结果
		if currentLogic == "AND" {
			result = result && matched
		} else if currentLogic == "OR" {
			result = result || matched
		}

		// 优化：如果AND已经是false，或者OR已经是true，可以提前退出
		if (currentLogic == "AND" && !result) || (currentLogic == "OR" && result) {
			// 不提前退出，因为可能需要检查所有条件的有效性
		}
	}

	return result
}

// evaluateCondition 评估单个条件
func (e *CustomFilterEngine) evaluateCondition(resource map[string]interface{}, condition FilterCondition) bool {
	fieldValue, exists := resource[condition.Field]
	if !exists {
		return false
	}

	switch condition.Operator {
	case "eq":
		return e.compareEqual(fieldValue, condition.Value)
	case "ne":
		return !e.compareEqual(fieldValue, condition.Value)
	case "gt":
		return e.compareNumbers(fieldValue, condition.Value) > 0
	case "lt":
		return e.compareNumbers(fieldValue, condition.Value) < 0
	case "gte":
		return e.compareNumbers(fieldValue, condition.Value) >= 0
	case "lte":
		return e.compareNumbers(fieldValue, condition.Value) <= 0
	case "in":
		return e.isInArray(fieldValue, condition.Value)
	case "not_in":
		return !e.isInArray(fieldValue, condition.Value)
	case "contains":
		return e.containsString(fieldValue, condition.Value)
	case "not_contains":
		return !e.containsString(fieldValue, condition.Value)
	case "starts_with":
		return e.startsWith(fieldValue, condition.Value)
	case "ends_with":
		return e.endsWith(fieldValue, condition.Value)
	case "is_null":
		return e.isNull(fieldValue)
	case "is_not_null":
		return !e.isNull(fieldValue)
	default:
		return false
	}
}

// compareEqual 比较相等
func (e *CustomFilterEngine) compareEqual(fieldValue, conditionValue interface{}) bool {
	return fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", conditionValue)
}

// compareNumbers 比较数字
func (e *CustomFilterEngine) compareNumbers(fieldValue, conditionValue interface{}) int {
	fieldFloat, err1 := toFloat64(fieldValue)
	conditionFloat, err2 := toFloat64(conditionValue)

	if err1 != nil || err2 != nil {
		return 0
	}

	if fieldFloat < conditionFloat {
		return -1
	} else if fieldFloat > conditionFloat {
		return 1
	}
	return 0
}

// isInArray 检查值是否在数组中
func (e *CustomFilterEngine) isInArray(fieldValue, conditionValue interface{}) bool {
	// 将条件值转换为数组
	conditionArray, ok := conditionValue.([]interface{})
	if !ok {
		// 尝试JSON解析
		if strVal, ok := conditionValue.(string); ok {
			var arr []interface{}
			if err := json.Unmarshal([]byte(strVal), &arr); err == nil {
				conditionArray = arr
			} else {
				return false
			}
		} else {
			return false
		}
	}

	fieldStr := fmt.Sprintf("%v", fieldValue)
	for _, item := range conditionArray {
		if fieldStr == fmt.Sprintf("%v", item) {
			return true
		}
	}

	return false
}

// containsString 检查字符串包含
func (e *CustomFilterEngine) containsString(fieldValue, conditionValue interface{}) bool {
	fieldStr := fmt.Sprintf("%v", fieldValue)
	conditionStr := fmt.Sprintf("%v", conditionValue)
	return strings.Contains(fieldStr, conditionStr)
}

// startsWith 检查字符串前缀
func (e *CustomFilterEngine) startsWith(fieldValue, conditionValue interface{}) bool {
	fieldStr := fmt.Sprintf("%v", fieldValue)
	conditionStr := fmt.Sprintf("%v", conditionValue)
	return strings.HasPrefix(fieldStr, conditionStr)
}

// endsWith 检查字符串后缀
func (e *CustomFilterEngine) endsWith(fieldValue, conditionValue interface{}) bool {
	fieldStr := fmt.Sprintf("%v", fieldValue)
	conditionStr := fmt.Sprintf("%v", conditionValue)
	return strings.HasSuffix(fieldStr, conditionStr)
}

// isNull 检查是否为null
func (e *CustomFilterEngine) isNull(fieldValue interface{}) bool {
	if fieldValue == nil {
		return true
	}

	// 检查空字符串、空数组、空map
	v := reflect.ValueOf(fieldValue)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	}

	return false
}

// toFloat64 转换为float64
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		// 尝试解析字符串
		var f float64
		_, err := fmt.Sscanf(v, "%f", &f)
		return f, err
	default:
		return 0, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "cannot convert to float64"),
                errorx.KV("value_type", fmt.Sprintf("%T", value)),
            )
	}
}

// BuildFilterCondition 构建过滤条件（辅助函数）
func BuildFilterCondition(field, operator string, value interface{}, logic string) FilterCondition {
	return FilterCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
		Logic:    logic,
	}
}

// ValidateConditions 验证过滤条件
func (e *CustomFilterEngine) ValidateConditions(conditions []FilterCondition) error {
	validOperators := map[string]bool{
		"eq":          true,
		"ne":          true,
		"gt":          true,
		"lt":          true,
		"gte":         true,
		"lte":         true,
		"in":          true,
		"not_in":      true,
		"contains":    true,
		"not_contains": true,
		"starts_with":  true,
		"ends_with":    true,
		"is_null":      true,
		"is_not_null":  true,
	}

	for i, condition := range conditions {
		// 验证操作符
		if !validOperators[condition.Operator] {
			return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "invalid operator"),
                errorx.KV("operator", condition.Operator),
                errorx.KV("condition_index", fmt.Sprintf("%d", i)),
            )
		}

		// 验证逻辑操作符
		if condition.Logic != "" && condition.Logic != "AND" && condition.Logic != "OR" {
			return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "invalid logic operator"),
                errorx.KV("logic", condition.Logic),
                errorx.KV("condition_index", fmt.Sprintf("%d", i)),
            )
		}

		// 验证字段名
		if condition.Field == "" {
			return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "empty field name"),
                errorx.KV("condition_index", fmt.Sprintf("%d", i)),
            )
		}
	}

	return nil
}
