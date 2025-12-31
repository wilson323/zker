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

package coze

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ValidateTenantID 验证租户ID格式
func ValidateTenantID(tenantID string) error {
	if tenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if len(tenantID) < 1 || len(tenantID) > 36 {
		return fmt.Errorf("tenant_id length must be between 1 and 36 characters")
	}
	// 可选：添加UUID格式验证
	uuidRegex := regexp.MustCompile(`^[a-fA-F0-9\-]{1,36}$`)
	if !uuidRegex.MatchString(tenantID) {
		return fmt.Errorf("tenant_id contains invalid characters")
	}
	return nil
}

// ValidateUUID 验证UUID格式
func ValidateUUID(uuid string) error {
	if uuid == "" {
		return fmt.Errorf("uuid is required")
	}
	uuidRegex := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	if !uuidRegex.MatchString(uuid) {
		return fmt.Errorf("invalid UUID format")
	}
	return nil
}

// ValidateDate 验证日期格式（YYYY-MM-DD）
func ValidateDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %s, expected YYYY-MM-DD", dateStr)
	}
	// 检查日期是否合理
	if t.Year() < 2000 || t.Year() > 2100 {
		return time.Time{}, fmt.Errorf("date year out of range: %d", t.Year())
	}
	return t, nil
}

// ValidateDateTime 验证日期时间格式（RFC3339）
func ValidateDateTime(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid datetime format: %s, expected RFC3339", dateStr)
	}
	return t, nil
}

// ValidateLimit 验证分页limit参数
func ValidateLimit(limitStr string, defaultValue int) (int, error) {
	if limitStr == "" {
		if defaultValue < 1 || defaultValue > 100 {
			return 20, nil // 确保默认值合理
		}
		return defaultValue, nil
	}
	limit := 0
	_, err := fmt.Sscanf(limitStr, "%d", &limit)
	if err != nil {
		return 0, fmt.Errorf("invalid limit format: %s", limitStr)
	}
	if limit < 1 {
		return 0, fmt.Errorf("limit must be at least 1, got: %d", limit)
	}
	if limit > 100 {
		return 0, fmt.Errorf("limit must not exceed 100, got: %d", limit)
	}
	return limit, nil
}

// ValidateOffset 验证分页offset参数
func ValidateOffset(offsetStr string) (int, error) {
	if offsetStr == "" {
		return 0, nil
	}
	offset := 0
	_, err := fmt.Sscanf(offsetStr, "%d", &offset)
	if err != nil {
		return 0, fmt.Errorf("invalid offset format: %s", offsetStr)
	}
	if offset < 0 {
		return 0, fmt.Errorf("offset must be non-negative, got: %d", offset)
	}
	if offset > 10000 {
		return 0, fmt.Errorf("offset must not exceed 10000, got: %d", offset)
	}
	return offset, nil
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if len(email) > 254 { // RFC 5321
		return fmt.Errorf("email length exceeds maximum 254 characters")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}
	return nil
}

// ValidatePhoneNumber 验证手机号格式（中国大陆）
func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number is required")
	}
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format: %s", phone)
	}
	return nil
}

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("username length must be between 3 and 32 characters")
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, hyphens and underscores")
	}
	return nil
}

// ValidateStringRange 验证字符串长度范围
func ValidateStringRange(value, fieldName string, minLen, maxLen int) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	length := len([]rune(value)) // 使用rune计算UTF-8字符数
	if length < minLen {
		return fmt.Errorf("%s length must be at least %d characters", fieldName, minLen)
	}
	if length > maxLen {
		return fmt.Errorf("%s length must not exceed %d characters", fieldName, maxLen)
	}
	return nil
}

// ValidateIntRange 验证整数范围
func ValidateIntRange(value int, fieldName string, minVal, maxVal int) error {
	if value < minVal {
		return fmt.Errorf("%s must be at least %d, got: %d", fieldName, minVal, value)
	}
	if value > maxVal {
		return fmt.Errorf("%s must not exceed %d, got: %d", fieldName, maxVal, value)
	}
	return nil
}

// ValidateStatus 验证状态值
func ValidateStatus(status, fieldName string, validStatuses []string) error {
	if status == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return fmt.Errorf("invalid %s: %s, valid values are: %s", fieldName, status, strings.Join(validStatuses, ", "))
}

// ValidateSortBy 验证排序字段（防止SQL注入）
func ValidateSortBy(sortBy, fieldName string, allowedFields []string) error {
	if sortBy == "" {
		return nil // 可选参数
	}
	for _, allowedField := range allowedFields {
		if sortBy == allowedField {
			return nil
		}
	}
	return fmt.Errorf("invalid %s: %s, allowed fields are: %s", fieldName, sortBy, strings.Join(allowedFields, ", "))
}

// ValidateSortOrder 验证排序方向
func ValidateSortOrder(sortOrder string) error {
	if sortOrder == "" {
		return nil // 可选参数，默认asc
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		return fmt.Errorf("invalid sort_order: %s, must be 'asc' or 'desc'", sortOrder)
	}
	return nil
}

// ValidatePositiveInteger 验证正整数
func ValidatePositiveInteger(valueStr, fieldName string) (int, error) {
	if valueStr == "" {
		return 0, fmt.Errorf("%s is required", fieldName)
	}
	value := 0
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: %s", fieldName, valueStr)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s must be positive, got: %d", fieldName, value)
	}
	return value, nil
}

// ValidateNonNegativeInteger 验证非负整数
func ValidateNonNegativeInteger(valueStr, fieldName string) (int, error) {
	if valueStr == "" {
		return 0, nil
	}
	value := 0
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: %s", fieldName, valueStr)
	}
	if value < 0 {
		return 0, fmt.Errorf("%s must be non-negative, got: %d", fieldName, value)
	}
	return value, nil
}

// SanitizeString 清理字符串输入（去除前后空格，防止XSS）
func SanitizeString(input string) string {
	input = strings.TrimSpace(input)
	// 移除可能的恶意字符（基本清理）
	input = strings.ReplaceAll(input, "\x00", "")
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\n", " ")
	return input
}

// ValidateTags 验证标签列表
func ValidateTags(tags []string) error {
	if len(tags) > 20 {
		return fmt.Errorf("too many tags, maximum is 20")
	}
	for i, tag := range tags {
		tag = strings.TrimSpace(tag)
		if len(tag) == 0 {
			return fmt.Errorf("tag at index %d is empty", i)
		}
		if len(tag) > 50 {
			return fmt.Errorf("tag at index %d exceeds maximum length of 50 characters", i)
		}
	}
	return nil
}

// ValidateEnum 验证枚举值
func ValidateEnum(value, fieldName string, allowedValues []string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}
	return fmt.Errorf("invalid %s: %s, allowed values are: %s", fieldName, value, strings.Join(allowedValues, ", "))
}
