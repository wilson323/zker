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

package middleware

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// 定义敏感字段列表
var sensitiveFields = []string{
	"password", "passwd", "pwd",
	"token", "secret", "api_key", "apikey",
	"credit_card", "cc_number", "ssn",
	"private_key", "access_token", "refresh_token",
	"authorization", "auth_token",
}

// 正则预编译（性能优化）
var sensitiveFieldPatterns = make([]*regexp.Regexp, 0, len(sensitiveFields))

func init() {
	for _, field := range sensitiveFields {
		pattern := regexp.MustCompile(
			fmt.Sprintf(`"%s"\s*:\s*"[^"]*"`, regexp.QuoteMeta(field)),
		)
		sensitiveFieldPatterns = append(sensitiveFieldPatterns, pattern)
	}
}

// SanitizeRequest 过滤请求体中的敏感字段
func SanitizeRequest(body []byte) []byte {
	result := body
	for _, pattern := range sensitiveFieldPatterns {
		result = pattern.ReplaceAll(result, []byte(`"$1":"***"`))
	}
	return result
}

// SanitizeLog 过滤日志中的敏感信息
func SanitizeLog(logMsg string) string {
	// 移除密码等敏感字段
	for _, field := range sensitiveFields {
		pattern := regexp.MustCompile(
			regexp.QuoteMeta(field)+`[:=]\s*[^\s,}]*`,
		)
		logMsg = pattern.ReplaceAllString(logMsg, field+"=***")
	}
	return logMsg
}

// SanitizeMap 过滤map中的敏感字段
func SanitizeMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		if isSensitiveField(k) {
			result[k] = "***"
		} else if subMap, ok := v.(map[string]interface{}); ok {
			result[k] = SanitizeMap(subMap)
		} else if subSlice, ok := v.([]interface{}); ok {
			result[k] = sanitizeSlice(subSlice)
		} else {
			result[k] = v
		}
	}
	return result
}

// SanitizeError 过滤错误消息中的敏感信息
func SanitizeError(err error) string {
	if err == nil {
		return ""
	}
	errMsg := err.Error()

	// 移除可能的SQL语句（防止泄露表结构）
	sqlPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)SELECT\s+.*?FROM\s+[\w_]+`),
		regexp.MustCompile(`(?i)INSERT\s+INTO\s+[\w_]+`),
		regexp.MustCompile(`(?i)UPDATE\s+[\w_]+`),
		regexp.MustCompile(`(?i)DELETE\s+FROM\s+[\w_]+`),
	}

	for _, pattern := range sqlPatterns {
		errMsg = pattern.ReplaceAllString(errMsg, "[SQL QUERY REDACTED]")
	}

	// 过滤敏感字段
	errMsg = SanitizeLog(errMsg)

	return errMsg
}

// SanitizeAuditLog 审计日志专用过滤
func SanitizeAuditLog(auditLog map[string]interface{}) map[string]interface{} {
	// 审计日志需要记录"有什么字段被修改"但不记录具体值
	result := make(map[string]interface{})

	for k, v := range auditLog {
		lowerKey := strings.ToLower(k)

		// 检查是否是敏感字段
		isSensitive := false
		for _, sensitiveField := range sensitiveFields {
			if strings.Contains(lowerKey, sensitiveField) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			// 记录"该字段存在"但不记录值
			result[k] = "[REDACTED]"
		} else {
			result[k] = v
		}
	}

	return result
}

// SanitizeResponse 过滤API响应中的敏感字段
func SanitizeResponse(responseData interface{}) interface{} {
	switch data := responseData.(type) {
	case map[string]interface{}:
		return SanitizeMap(data)
	case []interface{}:
		return sanitizeSlice(data)
	default:
		return responseData
	}
}

// isSensitiveField 检查字段名是否敏感
func isSensitiveField(fieldName string) bool {
	lowerField := strings.ToLower(fieldName)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(lowerField, sensitive) {
			return true
		}
	}
	return false
}

// sanitizeSlice 递归过滤slice
func sanitizeSlice(slice []interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		if subMap, ok := v.(map[string]interface{}); ok {
			result[i] = SanitizeMap(subMap)
		} else if subSlice, ok := v.([]interface{}); ok {
			result[i] = sanitizeSlice(subSlice)
		} else {
			result[i] = v
		}
	}
	return result
}

// SanitizeJSONString 过滤JSON字符串中的敏感字段
func SanitizeJSONString(jsonStr string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", err
	}

	sanitized := SanitizeResponse(data)
	result, err := json.Marshal(sanitized)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// MaskString 部分遮蔽字符串（如：显示前4位和后4位）
func MaskString(value string, showFirst, showLast int) string {
	if len(value) <= showFirst+showLast {
		return "***"
	}
	first := value[:showFirst]
	last := value[len(value)-showLast:]
	maskLen := len(value) - showFirst - showLast
	if maskLen < 0 {
		maskLen = 0
	}
	return fmt.Sprintf("%s%s%s", first, strings.Repeat("*", maskLen), last)
}

// MaskEmail 遮蔽邮箱（如：ab***@example.com）
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***"
	}
	username := parts[0]
	domain := parts[1]

	if len(username) <= 2 {
		username = "***"
	} else {
		username = username[:2] + "***"
	}

	return fmt.Sprintf("%s@%s", username, domain)
}

// MaskPhone 遮蔽手机号（如：138****5678）
func MaskPhone(phone string) string {
	if len(phone) != 11 {
		return "***"
	}
	return phone[:3] + "****" + phone[7:]
}
