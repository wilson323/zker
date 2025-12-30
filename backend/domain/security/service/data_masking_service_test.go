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
	"strings"
	"testing"
)

func TestDataMaskingService_MaskEmail(t *testing.T) {
	svc := NewDataMaskingService(MaskingLevelMedium)

	tests := []struct {
		name     string
		email    string
		level    MaskingLevel
		expected string
	}{
		{
			name:     "中等脱敏-正常邮箱",
			email:    "user@example.com",
			level:    MaskingLevelMedium,
			expected: "u***@example.com",
		},
		{
			name:     "低级脱敏",
			email:    "username@example.com",
			level:    MaskingLevelLow,
			expected: "user****@example.com",
		},
		{
			name:     "高级脱敏",
			email:    "user@example.com",
			level:    MaskingLevelHigh,
			expected: "***@******.com",
		},
		{
			name:     "无脱敏",
			email:    "user@example.com",
			level:    MaskingLevelNone,
			expected: "user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.MaskEmailWithLevel(tt.email, tt.level)
			if result != tt.expected {
				t.Errorf("MaskEmail() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDataMaskingService_MaskPhone(t *testing.T) {
	svc := NewDataMaskingService(MaskingLevelMedium)

	tests := []struct {
		name     string
		phone    string
		level    MaskingLevel
		expected string
	}{
		{
			name:     "中等脱敏-正常手机号",
			phone:    "13812345678",
			level:    MaskingLevelMedium,
			expected: "138****5678",
		},
		{
			name:     "带区号",
			phone:    "+8613812345678",
			level:    MaskingLevelMedium,
			expected: "138****5678",
		},
		{
			name:     "高级脱敏",
			phone:    "13812345678",
			level:    MaskingLevelHigh,
			expected: "*******5678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.MaskPhoneWithLevel(tt.phone, tt.level)
			if result != tt.expected {
				t.Errorf("MaskPhone() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDataMaskingService_MaskIDCard(t *testing.T) {
	svc := NewDataMaskingService(MaskingLevelMedium)

	tests := []struct {
		name     string
		idCard   string
		expected string
	}{
		{
			name:     "18位身份证",
			idCard:   "110101199001011234",
			expected: "110101********1234",
		},
		{
			name:     "15位身份证",
			idCard:   "110101900101123",
			expected: "1101**********123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.MaskIDCard(tt.idCard)
			if result != tt.expected {
				t.Errorf("MaskIDCard() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDataMaskingService_MaskBankCard(t *testing.T) {
	svc := NewDataMaskingService(MaskingLevelMedium)

	tests := []struct {
		name     string
		card     string
		expected string
	}{
		{
			name:     "16位银行卡",
			card:     "6222021234567890",
			expected: "6222**********7890",
		},
		{
			name:     "19位银行卡",
			card:     "6222021234567890123",
			expected: "6222*************0123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.MaskBankCard(tt.card)
			if result != tt.expected {
				t.Errorf("MaskBankCard() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDataMaskingService_MaskStringInString(t *testing.T) {
	svc := NewDataMaskingService(MaskingLevelMedium)

	tests := []struct {
		name     string
		content  string
		contains []string
		notContains []string
	}{
		{
			name:    "包含手机号",
			content: "请联系13812345678",
			contains: []string{"138****"},
			notContains: []string{"13812345678"},
		},
		{
			name:    "包含邮箱",
			content: "邮箱是user@example.com",
			contains: []string{"u***@"},
			notContains: []string{"user@example.com"},
		},
		{
			name:    "包含身份证",
			content: "身份证110101199001011234",
			contains: []string{"110101********"},
			notContains: []string{"110101199001011234"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.MaskStringInString(tt.content)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("MaskStringInString() should contain %v, got %v", expected, result)
				}
			}

			for _, unexpected := range tt.notContains {
				if strings.Contains(result, unexpected) {
					t.Errorf("MaskStringInString() should not contain %v, got %v", unexpected, result)
				}
			}
		})
	}
}

func TestDataMaskingService_SanitizeForLog(t *testing.T) {
	svc := NewDataMaskingService(MaskingLevelMedium)

	tests := []struct {
		name  string
		input string
		check func(string) bool
	}{
		{
			name:  "移除JWT",
			input: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			check: func(s string) bool {
				return !strings.Contains(s, "eyJ") && strings.Contains(s, "***JWT***")
			},
		},
		{
			name:  "移除API Key",
			input: "api_key: sk-1234567890abcdefghijklmnopqrstuvwxyz",
			check: func(s string) bool {
				return !strings.Contains(s, "sk-1234567890") && strings.Contains(s, "***API_KEY***")
			},
		},
		{
			name:  "脱敏手机号",
			input: "phone: 13812345678",
			check: func(s string) bool {
				return strings.Contains(s, "****") && !strings.Contains(s, "13812345678")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.SanitizeForLog(tt.input)
			if !tt.check(result) {
				t.Errorf("SanitizeForLog() = %v, check failed", result)
			}
		})
	}
}

func TestDataMaskingService_MaskJSON(t *testing.T) {
	svc := NewDataMaskingService(MaskingLevelMedium)

	jsonData := `{"email":"user@example.com","phone":"13812345678","name":"张三"}`

	result, err := svc.MaskJSON([]byte(jsonData), "email", "phone", "name")
	if err != nil {
		t.Fatalf("MaskJSON() error = %v", err)
	}

	resultStr := string(result)

	// 验证脱敏
	if strings.Contains(resultStr, "user@example.com") {
		t.Error("email should be masked")
	}
	if strings.Contains(resultStr, "13812345678") {
		t.Error("phone should be masked")
	}
	if !strings.Contains(resultStr, "***") {
		t.Error("should contain masking characters")
	}
}

func TestGenerateSecret(t *testing.T) {
	secret := generateSecret()
	if len(secret) != 32 { // Base32编码后160bits = 32字符
		t.Errorf("generateSecret() length = %v, want 32", len(secret))
	}
}

func TestGenerateRecoveryCode(t *testing.T) {
	code := generateRecoveryCode()

	// 格式: XXXX-XXXX (8位+1个横杠)
	if len(code) != 9 {
		t.Errorf("generateRecoveryCode() length = %v, want 9", len(code))
	}

	if !strings.Contains(code, "-") {
		t.Error("recovery code should contain '-'")
	}

	parts := strings.Split(code, "-")
	if len(parts) != 2 || len(parts[0]) != 4 || len(parts[1]) != 4 {
		t.Error("recovery code format should be XXXX-XXXX")
	}
}

func BenchmarkMaskEmail(b *testing.B) {
	svc := NewDataMaskingService(MaskingLevelMedium)
	email := "user@example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.MaskEmail(email)
	}
}

func BenchmarkMaskPhone(b *testing.B) {
	svc := NewDataMaskingService(MaskingLevelMedium)
	phone := "13812345678"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.MaskPhone(phone)
	}
}

func BenchmarkSanitizeForLog(b *testing.B) {
	svc := NewDataMaskingService(MaskingLevelMedium)
	content := `{"email":"user@example.com","phone":"13812345678","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.SanitizeForLog(content)
	}
}
