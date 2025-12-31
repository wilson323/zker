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

package security

import (
	"strings"
	"testing"
)

// TestEscapeHTML 测试HTML转义功能
func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "转义script标签",
			input:    "<script>alert('XSS')</script>",
			expected: "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
		},
		{
			name:     "转义img标签onerror事件",
			input:    "<img src=x onerror=alert('XSS')>",
			expected: "&lt;img src=x onerror=alert(&#39;XSS&#39;)&gt;",
		},
		{
			name:     "转义特殊字符",
			input:    "<div>&\"'hello'</div>",
			expected: "&lt;div&gt;&amp;&#34;&#39;hello&#39;&lt;/div&gt;",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "纯文本",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeHTML(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeHTML() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSanitizeHTML 测试HTML清理功能
func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantSafe bool
	}{
		{
			name:     "移除script标签",
			input:    "<script>alert('XSS')</script>Hello",
			wantSafe: true,
		},
		{
			name:     "移除iframe标签",
			input:    "<iframe src='http://evil.com'></iframe>",
			wantSafe: true,
		},
		{
			name:     "移除onclick事件",
			input:    "<div onclick='alert(1)'>Click</div>",
			wantSafe: true,
		},
		{
			name:     "移除javascript:协议",
			input:    "<a href='javascript:alert(1)'>Click</a>",
			wantSafe: true,
		},
		{
			name:     "保留安全的HTML",
			input:    "<p>Hello, <b>World</b>!</p>",
			wantSafe: true,
		},
		{
			name:     "纯文本",
			input:    "Just plain text",
			wantSafe: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeHTML(tt.input)
			// 检查结果中是否还包含危险内容
			if strings.Contains(strings.ToLower(result), "<script") ||
				strings.Contains(strings.ToLower(result), "<iframe") ||
				strings.Contains(strings.ToLower(result), "javascript:") {
				t.Errorf("SanitizeHTML() failed to sanitize: %v -> %v", tt.input, result)
			}
		})
	}
}

// TestValidateInput 测试输入验证功能
func TestValidateInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValid   bool
		wantReason  string
	}{
		{
			name:       "检测script标签",
			input:      "<script>alert(1)</script>",
			wantValid:  false,
			wantReason: "input contains dangerous HTML tags",
		},
		{
			name:       "检测javascript:协议",
			input:      "javascript:alert(1)",
			wantValid:  false,
			wantReason: "input contains JavaScript code",
		},
		{
			name:       "检测SQL注入",
			input:      "admin' OR '1'='1",
			wantValid:  false,
			wantReason: "input contains SQL injection patterns",
		},
		{
			name:       "检测命令注入",
			input:      "hello; rm -rf /",
			wantValid:  false,
			wantReason: "input contains command injection patterns",
		},
		{
			name:       "安全的文本输入",
			input:      "Hello, World!",
			wantValid:  true,
			wantReason: "",
		},
		{
			name:       "空字符串",
			input:      "",
			wantValid:  true,
			wantReason: "",
		},
		{
			name:       "带引号的文本",
			input:      "It's a beautiful day",
			wantValid:  true,
			wantReason: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, reason := ValidateInput(tt.input)
			if valid != tt.wantValid {
				t.Errorf("ValidateInput() valid = %v, want %v", valid, tt.wantValid)
			}
			if tt.wantReason != "" && !strings.Contains(reason, tt.wantReason) {
				t.Errorf("ValidateInput() reason = %v, want包含 %v", reason, tt.wantReason)
			}
		})
	}
}

// TestStripTags 测试移除HTML标签功能
func TestStripTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "移除所有HTML标签",
			input:    "<p>Hello, <b>World</b>!</p>",
			expected: "Hello, World!",
		},
		{
			name:     "移除script标签",
			input:    "<script>alert(1)</script>Hello",
			expected: "Hello",
		},
		{
			name:     "纯文本",
			input:    "Just plain text",
			expected: "Just plain text",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "嵌套标签",
			input:    "<div><p><span>Nested</span></p></div>",
			expected: "Nested",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripTags(tt.input)
			if result != tt.expected {
				t.Errorf("StripTags() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestTruncateText 测试文本截断功能
func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "截断长文本",
			input:    "This is a very long text that should be truncated",
			maxLen:   10,
			expected: "This is a ",
		},
		{
			name:     "不截断短文本",
			input:    "Short",
			maxLen:   10,
			expected: "Short",
		},
		{
			name:     "空字符串",
			input:    "",
			maxLen:   10,
			expected: "",
		},
		{
			name:     "中文字符",
			input:    "这是一段很长的中文文本需要被截断",
			maxLen:   10,
			expected: "这是一段很长的中文",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateText(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsValidURL 测试URL验证功能
func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "有效的HTTPS URL",
			url:      "https://example.com",
			expected: true,
		},
		{
			name:     "有效的HTTP URL",
			url:      "http://example.com",
			expected: true,
		},
		{
			name:     "javascript:协议",
			url:      "javascript:alert(1)",
			expected: false,
		},
		{
			name:     "data:协议",
			url:      "data:text/html,<script>alert(1)</script>",
			expected: false,
		},
		{
			name:     "file:协议",
			url:      "file:///etc/passwd",
			expected: false,
		},
		{
			name:     "空字符串",
			url:      "",
			expected: false,
		},
		{
			name:     "没有协议",
			url:      "example.com",
			expected: false,
		},
		{
			name:     "带路径的URL",
			url:      "https://example.com/path/to/page?query=value",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsValidURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestXSSAttackVectors 测试常见的XSS攻击向量
func TestXSSAttackVectors(t *testing.T) {
	attackVectors := []string{
		// Script注入
		"<script>alert('XSS')</script>",
		"<SCRIPT>alert('XSS')</SCRIPT>",
		"<img src=x onerror=alert('XSS')>",
		"<svg onload=alert('XSS')>",

		// 事件处理器
		"<div onmouseover='alert(1)'>Hover</div>",
		"<input onfocus='alert(1)'>",
		"<body onload='alert(1)'>",

		// JavaScript协议
		"<a href='javascript:alert(1)'>Click</a>",
		"<a href='JAVASCRIPT:alert(1)'>Click</a>",

		// Style注入
		"<div style='background:url(javascript:alert(1))'>",
		"<style>@import 'javascript:alert(1)'</style>",

		// 混淆攻击
		"<img src=x onerror=&quot;alert('XSS')&quot;>",
		"<script\\x20type='text/javascript'>alert(1)</script>",
	}

	for _, payload := range attackVectors {
		t.Run(payload, func(t *testing.T) {
			// 测试EscapeHTML
			escaped := EscapeHTML(payload)
			if strings.Contains(escaped, "<script>") ||
				strings.Contains(escaped, "<img src") {
				t.Errorf("EscapeHTML failed to sanitize: %s -> %s", payload, escaped)
			}

			// 测试ValidateInput
			valid, _ := ValidateInput(payload)
			if valid {
				t.Errorf("ValidateInput failed to detect attack: %s", payload)
			}

			// 测试SanitizeHTML
			sanitized := SanitizeHTML(payload)
			if strings.Contains(strings.ToLower(sanitized), "javascript:") {
				t.Errorf("SanitizeHTML failed to remove javascript: protocol: %s -> %s", payload, sanitized)
			}
		})
	}
}

// BenchmarkEscapeHTML 性能测试
func BenchmarkEscapeHTML(b *testing.B) {
	input := "<script>alert('XSS')</script><div>Hello</div>"
	for i := 0; i < b.N; i++ {
		EscapeHTML(input)
	}
}

// BenchmarkValidateInput 性能测试
func BenchmarkValidateInput(b *testing.B) {
	input := "Normal text without HTML"
	for i := 0; i < b.N; i++ {
		ValidateInput(input)
	}
}
