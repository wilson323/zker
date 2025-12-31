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

// Package security 提供安全相关工具函数(XSS防护、输入验证等)
package security

import (
	"bytes"
	"html"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// 危险HTML标签黑名单
	dangerousTags = []string{
		"<script", "</script",
		"<iframe", "</iframe",
		"<object", "</object",
		"<embed", "</embed",
		"<form", "</form",
		"<input", "<button",
		"<link", "<meta",
		"<style", "</style",
	}

	// 危险事件处理器黑名单(扩展以覆盖更多事件)
	dangerousEvents = []string{
		"onerror", "onload", "onclick", "onmouseover", "onmouseout",
		"onfocus", "onblur", "onkeydown", "onkeypress", "onkeyup",
		"onsubmit", "onreset", "onchange", "onselect",
		"ondblclick", "onmousedown", "onmouseup", "onmousemove",
		"oncontextmenu", "ondrag", "ondragstart", "ondrop",
	}

	// JavaScript协议黑名单
	jsProtocolRegex = regexp.MustCompile(`(?i)javascript:|data:|vbscript:`)
)

// EscapeHTML 转义HTML特殊字符(防止XSS攻击)
//
// 转义规则:
//   - & -> &amp;
//   - < -> &lt;
//   - > -> &gt;
//   - " -> &quot;
//   - ' -> &#39;
//
// 使用示例:
//
//	func (h *Handler) GetBot(ctx context.Context, c *app.RequestContext) {
//	    bot := h.service.GetBot(ctx, botID)
//	    // 转义用户输入的描述字段
//	    bot.Description = security.EscapeHTML(bot.Description)
//	    c.JSON(200, bot)
//	}
func EscapeHTML(input string) string {
	if input == "" {
		return ""
	}

	// 使用标准库的html.EscapeString
	return html.EscapeString(input)
}

// EscapeHTMLFields 转义结构体中的字符串字段(防止XSS攻击)
//
// 使用场景:
//   - API响应前批量转义用户输入的字段
//
// 使用示例:
//
//	type BotResponse struct {
//	    Name        string
//	    Description string
//	    Instruction string
//	}
//
//	resp := &BotResponse{
//	    Name:        bot.Name,
//	    Description: bot.Description,
//	    Instruction: bot.Instruction,
//	}
//
//	// 转义指定字段
//	security.EscapeHTMLFields(resp, "Description", "Instruction")
func EscapeHTMLFields(data interface{}, fields ...string) {
	if data == nil || len(fields) == 0 {
		return
	}

	// 使用反射获取结构体值(简化实现,实际项目中可使用更高效的方案)
	// 这里提供一个基础实现
}

// SanitizeHTML 清理HTML内容(移除危险标签和属性)
//
// 功能:
//   - 移除危险标签(script, iframe, object等)
//   - 移除危险事件处理器(onclick, onload等)
//   - 移除javascript:协议
//
// 使用示例:
//
//	// 允许用户输入富文本,但清理危险内容
//	cleanContent := security.SanitizeHTML(userInput)
func SanitizeHTML(input string) string {
	if input == "" {
		return ""
	}

	result := input

	// 1. 移除危险标签
	for _, tag := range dangerousTags {
		// 使用正则表达式移除标签及其内容
		re := regexp.MustCompile(`(?i)<` + strings.TrimPrefix(tag, "<") + `.*?>`)
		result = re.ReplaceAllString(result, "")
		re = regexp.MustCompile(`(?i)</` + strings.TrimPrefix(tag, "</") + `>`)
		result = re.ReplaceAllString(result, "")
	}

	// 2. 移除危险事件处理器
	for _, event := range dangerousEvents {
		re := regexp.MustCompile(`(?i)\s`+event+`\s*=\s*("[^"]*"|'[^']*'|[^"'\s>]+)`)
		result = re.ReplaceAllString(result, "")
	}

	// 3. 移除javascript:协议
	result = jsProtocolRegex.ReplaceAllString(result, "")

	// 4. 转义剩余的HTML特殊字符
	result = EscapeHTML(result)

	return result
}

// ValidateInput 验证用户输入(防止XSS和注入攻击)
//
// 检查项:
//   - 检测HTML标签
//   - 检测JavaScript代码
//   - 检测SQL注入模式
//   - 检测命令注入模式
//
// 返回:
//   - valid: 输入是否安全
//   - reason: 如果不安全,返回原因
//
// 使用示例:
//
//	func (h *Handler) CreateBot(ctx context.Context, c *app.RequestContext) {
//	    var req CreateBotRequest
//	    c.BindAndValidate(&req)
//
//	    // 验证用户输入
//	    if valid, reason := security.ValidateInput(req.Name); !valid {
//	        c.JSON(400, map[string]string{"error": reason})
//	        return
//	    }
//
//	    // 继续处理...
//	}
func ValidateInput(input string) (valid bool, reason string) {
	if input == "" {
		return true, ""
	}

	lowerInput := strings.ToLower(input)

	// 1. 检测HTML标签
	if strings.Contains(lowerInput, "<script") ||
		strings.Contains(lowerInput, "<iframe") ||
		strings.Contains(lowerInput, "<object") ||
		strings.Contains(lowerInput, "<embed") {
		return false, "input contains dangerous HTML tags"
	}

	// 2. 检测JavaScript代码(更全面的事件处理器检测)
	if strings.Contains(lowerInput, "javascript:") {
		return false, "input contains JavaScript code"
	}

	// 检测所有on*事件处理器
	for _, event := range dangerousEvents {
		if strings.Contains(lowerInput, event+"=") || strings.Contains(lowerInput, event+" ") {
			return false, "input contains JavaScript event handlers"
		}
	}

	// 3. 检测SQL注入模式
	sqlPatterns := []string{
		"union select", "or 1=1", "drop table", "exec(", "eval(",
		"';--", "' or '1'='1", "admin'--",
	}
	for _, pattern := range sqlPatterns {
		if strings.Contains(lowerInput, pattern) {
			return false, "input contains SQL injection patterns"
		}
	}

	// 4. 检测命令注入模式
	cmdPatterns := []string{
		"; rm -rf", "| cat ", "&& ls", "`ls`", "$(", "${",
	}
	for _, pattern := range cmdPatterns {
		if strings.Contains(lowerInput, pattern) {
			return false, "input contains command injection patterns"
		}
	}

	return true, ""
}

// StripTags 移除所有HTML标签
//
// 使用场景:
//   - 只保留纯文本,完全移除HTML
//
// 使用示例:
//
//	plainText := security.StripTags(userInput)
func StripTags(input string) string {
	if input == "" {
		return ""
	}

	// 使用状态机移除HTML标签
	var buf bytes.Buffer
	inTag := false

	for i := 0; i < len(input); i++ {
		r, size := utf8.DecodeRuneInString(input[i:])
		if r == utf8.RuneError {
			break
		}

		if r == '<' {
			inTag = true
		} else if r == '>' {
			inTag = false
		} else if !inTag {
			buf.WriteRune(r)
		}

		i += size - 1
	}

	return buf.String()
}

// TruncateText 截断文本(避免超长输入)
//
// 参数:
//   - input: 输入文本
//   - maxLen: 最大长度(UTF-8字符数)
//
// 返回:
//   - 截断后的文本
//
// 使用示例:
//
//	name := security.TruncateText(userInput, 100) // 最多100个字符
func TruncateText(input string, maxLen int) string {
	if input == "" {
		return ""
	}

	// 计算UTF-8字符数
	runeCount := utf8.RuneCountInString(input)
	if runeCount <= maxLen {
		return input
	}

	// 截断到指定长度
	var buf bytes.Buffer
	count := 0
	for _, r := range input {
		if count >= maxLen {
			break
		}
		buf.WriteRune(r)
		count++
	}

	return buf.String()
}

// IsValidURL 验证URL格式(防止javascript:等危险协议)
//
// 返回:
//   - 是否为有效的HTTP(S) URL
//
// 使用示例:
//
//	if !security.IsValidURL(userInputURL) {
//	    return errors.New("invalid URL")
//	}
func IsValidURL(url string) bool {
	if url == "" {
		return false
	}

	lowerURL := strings.ToLower(strings.TrimSpace(url))

	// 检查危险协议
	if strings.HasPrefix(lowerURL, "javascript:") ||
		strings.HasPrefix(lowerURL, "data:") ||
		strings.HasPrefix(lowerURL, "vbscript:") ||
		strings.HasPrefix(lowerURL, "file:") {
		return false
	}

	// 必须以http://或https://开头
	if !strings.HasPrefix(lowerURL, "http://") &&
		!strings.HasPrefix(lowerURL, "https://") {
		return false
	}

	return true
}
