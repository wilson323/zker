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
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// MaskingLevel 脱敏级别
type MaskingLevel int

const (
	// MaskingLevelNone 不脱敏(仅管理员可见)
	MaskingLevelNone MaskingLevel = iota
	// MaskingLevelLow 低级别脱敏(保留大部分信息)
	MaskingLevelLow
	// MaskingLevelMedium 中级别脱敏(保留部分信息)
	MaskingLevelMedium
	// MaskingLevelHigh 高级别脱敏(仅保留格式)
	MaskingLevelHigh
)

// DataMaskingService 数据脱敏服务
type DataMaskingService struct {
	defaultLevel MaskingLevel
}

// NewDataMaskingService 创建数据脱敏服务实例
func NewDataMaskingService(defaultLevel MaskingLevel) *DataMaskingService {
	return &DataMaskingService{
		defaultLevel: defaultLevel,
	}
}

// MaskEmail 邮箱脱敏
// 示例: user***@example.com
func (s *DataMaskingService) MaskEmail(email string) string {
	return s.MaskEmailWithLevel(email, s.defaultLevel)
}

// MaskEmailWithLevel 按级别脱敏邮箱
func (s *DataMaskingService) MaskEmailWithLevel(email string, level MaskingLevel) string {
	if email == "" {
		return ""
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		// 格式不合法，全部脱敏
		return "***@***"
	}

	username := parts[0]
	domain := parts[1]

	switch level {
	case MaskingLevelNone:
		return email
	case MaskingLevelLow:
		// user***@example.com
		if len(username) > 4 {
			return username[:4] + strings.Repeat("*", len(username)-4) + "@" + domain
		}
		return strings.Repeat("*", len(username)) + "@" + domain
	case MaskingLevelMedium:
		// u***@example.com
		if len(username) > 1 {
			return string(username[0]) + strings.Repeat("*", len(username)-1) + "@" + domain
		}
		return "*" + "@" + domain
	case MaskingLevelHigh:
		// ***@***.com
		return strings.Repeat("*", len(username)) + "@" + s.maskDomain(domain)
	default:
		return s.MaskEmailWithLevel(email, MaskingLevelMedium)
	}
}

// maskDomain 脱敏域名(保留顶级域)
func (s *DataMaskingService) maskDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return strings.Repeat("*", len(domain))
	}

	// 保留顶级域
	tld := parts[len(parts)-1]
	maskedParts := make([]string, len(parts))
	for i := 0; i < len(parts)-1; i++ {
		maskedParts[i] = strings.Repeat("*", len(parts[i]))
	}
	maskedParts[len(parts)-1] = tld

	return strings.Join(maskedParts, ".")
}

// MaskPhone 手机号脱敏
// 示例: 138****5678
func (s *DataMaskingService) MaskPhone(phone string) string {
	return s.MaskPhoneWithLevel(phone, s.defaultLevel)
}

// MaskPhoneWithLevel 按级别脱敏手机号
func (s *DataMaskingService) MaskPhoneWithLevel(phone string, level MaskingLevel) string {
	if phone == "" {
		return ""
	}

	// 移除非数字字符
	phone = s.cleanPhone(phone)

	if len(phone) < 7 {
		// 长度不足，全部脱敏
		return strings.Repeat("*", len(phone))
	}

	switch level {
	case MaskingLevelNone:
		return phone
	case MaskingLevelLow:
		// 1381234****78
		if len(phone) >= 11 {
			return phone[:7] + strings.Repeat("*", len(phone)-7-2) + phone[len(phone)-2:]
		}
		return phone[:4] + strings.Repeat("*", len(phone)-6) + phone[len(phone)-2:]
	case MaskingLevelMedium:
		// 138****5678
		if len(phone) >= 11 {
			return phone[:3] + strings.Repeat("*", 4) + phone[len(phone)-4:]
		}
		return phone[:3] + strings.Repeat("*", len(phone)-5) + phone[len(phone)-2:]
	case MaskingLevelHigh:
		// *******5678
		return strings.Repeat("*", len(phone)-4) + phone[len(phone)-4:]
	default:
		return s.MaskPhoneWithLevel(phone, MaskingLevelMedium)
	}
}

// cleanPhone 清理手机号(移除+86、空格、横线等)
func (s *DataMaskingService) cleanPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.TrimPrefix(phone, "+86")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	return phone
}

// MaskIDCard 身份证号脱敏
// 示例: 110101********1234
func (s *DataMaskingService) MaskIDCard(idCard string) string {
	return s.MaskIDCardWithLevel(idCard, s.defaultLevel)
}

// MaskIDCardWithLevel 按级别脱敏身份证号
func (s *DataMaskingService) MaskIDCardWithLevel(idCard string, level MaskingLevel) string {
	if idCard == "" {
		return ""
	}

	// 移除空格
	idCard = strings.TrimSpace(idCard)

	if len(idCard) < 8 {
		return strings.Repeat("*", len(idCard))
	}

	switch level {
	case MaskingLevelNone:
		return idCard
	case MaskingLevelLow:
		// 110101199001************4
		if len(idCard) == 18 {
			return idCard[:6] + idCard[6:10] + strings.Repeat("*", len(idCard)-10-2) + idCard[len(idCard)-2:]
		}
		return idCard[:4] + strings.Repeat("*", len(idCard)-6) + idCard[len(idCard)-2:]
	case MaskingLevelMedium:
		// 110101********1234
		return idCard[:6] + strings.Repeat("*", 8) + idCard[len(idCard)-4:]
	case MaskingLevelHigh:
		// ************1234
		return strings.Repeat("*", len(idCard)-4) + idCard[len(idCard)-4:]
	default:
		return s.MaskIDCardWithLevel(idCard, MaskingLevelMedium)
	}
}

// MaskBankCard 银行卡号脱敏
// 示例: 6222***********1234
func (s *DataMaskingService) MaskBankCard(card string) string {
	return s.MaskBankCardWithLevel(card, s.defaultLevel)
}

// MaskBankCardWithLevel 按级别脱敏银行卡号
func (s *DataMaskingService) MaskBankCardWithLevel(card string, level MaskingLevel) string {
	if card == "" {
		return ""
	}

	// 移除空格
	card = strings.TrimSpace(card)
	card = strings.ReplaceAll(card, " ", "")

	if len(card) < 8 {
		return strings.Repeat("*", len(card))
	}

	switch level {
	case MaskingLevelNone:
		return card
	case MaskingLevelLow:
		// 622202123456****1234
		if len(card) >= 16 {
			return card[:10] + strings.Repeat("*", len(card)-10-4) + card[len(card)-4:]
		}
		return card[:6] + strings.Repeat("*", len(card)-10) + card[len(card)-4:]
	case MaskingLevelMedium:
		// 6222***********1234
		return card[:4] + strings.Repeat("*", len(card)-8) + card[len(card)-4:]
	case MaskingLevelHigh:
		// ****1234
		return strings.Repeat("*", 4) + card[len(card)-4:]
	default:
		return s.MaskBankCardWithLevel(card, MaskingLevelMedium)
	}
}

// MaskString 通用字符串脱敏
// 保留前prefix和后suffix个字符
func (s *DataMaskingService) MaskString(str string, prefix, suffix int) string {
	if str == "" {
		return ""
	}

	length := len(str)
	if length <= prefix+suffix {
		return strings.Repeat("*", length)
	}

	maskedLength := length - prefix - suffix
	return str[:prefix] + strings.Repeat("*", maskedLength) + str[length-suffix:]
}

// MaskPassword 密码脱敏(完全隐藏)
func (s *DataMaskingService) MaskPassword(password string) string {
	if password == "" {
		return ""
	}
	return "******"
}

// MaskName 姓名脱敏
// 示例: 张三 -> 张*; 张小三 -> 张小*
func (s *DataMaskingService) MaskName(name string) string {
	if name == "" {
		return ""
	}

	runes := []rune(name)
	if len(runes) == 1 {
		return "*"
	}

	// 保留姓氏，脱敏名字
	return string(runes[:1]) + strings.Repeat("*", len(runes)-1)
}

// MaskAddress 地址脱敏
// 示例: 北京市朝阳区***号
func (s *DataMaskingService) MaskAddress(address string) string {
	if address == "" {
		return ""
	}

	runes := []rune(address)
	if len(runes) <= 6 {
		return strings.Repeat("*", len(runes))
	}

	// 保留前6个字符，其余脱敏
	return string(runes[:6]) + strings.Repeat("*", len(runes)-6)
}

// MaskJSON 脱敏JSON中的敏感字段
func (s *DataMaskingService) MaskJSON(data []byte, sensitiveFields ...string) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// 脱敏敏感字段
	for _, field := range sensitiveFields {
		if value, ok := obj[field]; ok {
			if strValue, ok := value.(string); ok {
				obj[field] = s.maskFieldByType(field, strValue)
			}
		}
	}

	return json.Marshal(obj)
}

// maskFieldByType 根据字段名自动判断脱敏类型
func (s *DataMaskingService) maskFieldByType(fieldName, value string) string {
	fieldName = strings.ToLower(fieldName)

	switch {
	case strings.Contains(fieldName, "email"), strings.Contains(fieldName, "mail"):
		return s.MaskEmail(value)
	case strings.Contains(fieldName, "phone"), strings.Contains(fieldName, "mobile"), strings.Contains(fieldName, "tel"):
		return s.MaskPhone(value)
	case strings.Contains(fieldName, "idcard"), strings.Contains(fieldName, "id_card"):
		return s.MaskIDCard(value)
	case strings.Contains(fieldName, "bank"), strings.Contains(fieldName, "card"):
		return s.MaskBankCard(value)
	case strings.Contains(fieldName, "password"), strings.Contains(fieldName, "passwd"), strings.Contains(fieldName, "secret"):
		return s.MaskPassword(value)
	case strings.Contains(fieldName, "name") && !strings.Contains(fieldName, "username"):
		return s.MaskName(value)
	case strings.Contains(fieldName, "address"):
		return s.MaskAddress(value)
	default:
		// 未知类型，使用中等脱敏
		if len(value) > 8 {
			return s.MaskString(value, 2, 2)
		}
		return strings.Repeat("*", len(value))
	}
}

// MaskStringInString 在字符串中脱敏敏感信息(使用正则)
func (s *DataMaskingService) MaskStringInString(content string) string {
	if content == "" {
		return content
	}

	// 手机号正则
	phoneRegex := regexp.MustCompile(`1[3-9]\d{9}`)
	content = phoneRegex.ReplaceAllStringFunc(content, func(phone string) string {
		return s.MaskPhone(phone)
	})

	// 邮箱正则
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	content = emailRegex.ReplaceAllStringFunc(content, func(email string) string {
		return s.MaskEmail(email)
	})

	// 身份证正则
	idCardRegex := regexp.MustCompile(`\d{17}[\dXx]`)
	content = idCardRegex.ReplaceAllStringFunc(content, func(idCard string) string {
		return s.MaskIDCard(idCard)
	})

	// 银行卡正则(13-19位数字)
	bankCardRegex := regexp.MustCompile(`\d{13,19}`)
	content = bankCardRegex.ReplaceAllStringFunc(content, func(card string) string {
		// 避免脱敏身份证(18位)和纯数字
		if len(card) == 17 || len(card) == 18 {
			// 可能是身份证，跳过
			return card
		}
		return s.MaskBankCard(card)
	})

	return content
}

// MaskStruct 脱敏结构体中的敏感字段(通过标签)
// 支持的标签: mask:"type" 其中type可以是email, phone, idcard, bankcard, password, name, address
func (s *DataMaskingService) MaskStruct(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}

	// 使用JSON序列化/反序列化进行脱敏
	data, err := json.Marshal(obj)
	if err != nil {
		return obj
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return obj
	}

	// 递归脱敏
	s.maskMapRecursively(result)

	return result
}

// maskMapRecursively 递归脱敏map
func (s *DataMaskingService) maskMapRecursively(data map[string]interface{}) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			data[key] = s.maskFieldByType(key, v)
		case map[string]interface{}:
			s.maskMapRecursively(v)
		case []interface{}:
			s.maskSliceRecursively(v)
		}
	}
}

// maskSliceRecursively 递归脱敏slice
func (s *DataMaskingService) maskSliceRecursively(data []interface{}) {
	for i, value := range data {
		switch v := value.(type) {
		case string:
			// 数组中的字符串，如果是敏感格式则脱敏
			if s.isEmail(v) {
				data[i] = s.MaskEmail(v)
			} else if s.isPhone(v) {
				data[i] = s.MaskPhone(v)
			}
		case map[string]interface{}:
			s.maskMapRecursively(v)
		case []interface{}:
			s.maskSliceRecursively(v)
		}
	}
}

// isEmail 判断是否是邮箱
func (s *DataMaskingService) isEmail(str string) bool {
	return strings.Contains(str, "@") && strings.Contains(str, ".")
}

// isPhone 判断是否是手机号
func (s *DataMaskingService) isPhone(str string) bool {
	if len(str) != 11 {
		return false
	}
	for _, r := range str {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return str[0] == '1'
}

// SanitizeForLog 清理日志内容(移除所有敏感信息)
func (s *DataMaskingService) SanitizeForLog(content string) string {
	if content == "" {
		return content
	}

	// 先脱敏字符串中的敏感信息
	content = s.MaskStringInString(content)

	// 移除可能的JSON Web Token
	jwtRegex := regexp.MustCompile(`eyJ[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+`)
	content = jwtRegex.ReplaceAllString(content, "***JWT***")

	// 移除API Key
	apiKeyRegex := regexp.MustCompile(`[a-zA-Z0-9]{32,}`)
	content = apiKeyRegex.ReplaceAllString(content, "***API_KEY***")

	return content
}

// MaskingConfig 脱敏配置
type MaskingConfig struct {
	Enabled       bool                 `json:"enabled"`
	DefaultLevel  MaskingLevel         `json:"default_level"`
	FieldRules    map[string]MaskingLevel `json:"field_rules"`    // 字段级别的脱敏规则
	SensitiveData []string             `json:"sensitive_data"`   // 敏感数据模式列表
}

// ApplyConfig 应用脱敏配置
func (s *DataMaskingService) ApplyConfig(config *MaskingConfig) {
	if config == nil || !config.Enabled {
		return
	}

	s.defaultLevel = config.DefaultLevel
	// 可以在这里添加更多的配置应用逻辑
}

// NewMaskingConfig 创建默认脱敏配置
func NewMaskingConfig() *MaskingConfig {
	return &MaskingConfig{
		Enabled:      true,
		DefaultLevel: MaskingLevelMedium,
		FieldRules: map[string]MaskingLevel{
			"email":    MaskingLevelMedium,
			"phone":    MaskingLevelMedium,
			"idcard":   MaskingLevelHigh,
			"bankcard": MaskingLevelHigh,
			"password": MaskingLevelHigh,
		},
		SensitiveData: []string{
			"password", "secret", "token", "key", "auth",
		},
	}
}

// MaskReader 脱敏io.Reader内容
func (s *DataMaskingService) MaskReader(reader *bytes.Reader) (*bytes.Reader, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	content := buf.String()
	masked := s.MaskStringInString(content)

	return bytes.NewReader([]byte(masked)), nil
}

// GetMaskingPercentage 计算脱敏百分比(用于统计)
func (s *DataMaskingService) GetMaskingPercentage(original, masked string) float64 {
	if len(original) == 0 {
		return 0.0
	}

	originalRunes := []rune(original)
	maskedRunes := []rune(masked)

	maskedCount := 0
	for i := 0; i < len(originalRunes) && i < len(maskedRunes); i++ {
		if maskedRunes[i] == '*' && originalRunes[i] != '*' {
			maskedCount++
		}
	}

	return float64(maskedCount) / float64(len(originalRunes)) * 100.0
}
