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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/coze-studio/coze-studio/backend/types/errno"
)

// ErrorInfo 错误信息结构
type ErrorInfo struct {
	Code       int32  `json:"code"`
	Message    string `json:"message"`
	MessageZH  string `json:"message_zh,omitempty"`
	HTTPStatus int    `json:"http_status"`
	Category   string `json:"category"`
}

// ErrorCategory 错误分类
type ErrorCategory struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Errors      []ErrorInfo   `json:"errors"`
}

func main() {
	// 生成错误码文档
	generateErrorDocumentation()

	// 生成错误码JSON文件(供前端使用)
	generateErrorCodesJSON()

	// 生成错误码测试用例
	generateErrorTestCases()

	fmt.Println("✅ 错误码文档生成完成!")
	fmt.Println("📄 生成的文件:")
	fmt.Println("  - docs/errno/error_codes.md")
	fmt.Println("  - docs/errno/error_codes.json")
	fmt.Println("  - docs/errno/error_codes_zh.json")
	fmt.Println("  - types/errno/error_codes_test.go")
}

// generateErrorDocumentation 生成错误码文档
func generateErrorDocumentation() {
	// 收集所有错误码
	categories := collectErrorCategories()

	// 生成Markdown文档
	generateMarkdownDoc(categories)

	// 生成JSON文档
	generateJSONDoc(categories)
}

// collectErrorCategories 收集错误码分类
func collectErrorCategories() []ErrorCategory {
	categories := []ErrorCategory{
		{
			Name:        "租户错误 (Tenant Errors)",
			Description: "2xx开头,与租户管理、配额、订阅相关的错误",
			Errors: []ErrorInfo{
				{Code: errno.ErrTenantNotFoundCode, Message: "Tenant not found", MessageZH: "租户不存在", HTTPStatus: 404, Category: "tenant"},
				{Code: errno.ErrTenantAlreadyExistsCode, Message: "Tenant already exists", MessageZH: "租户已存在", HTTPStatus: 409, Category: "tenant"},
				{Code: errno.ErrTenantSuspendedCode, Message: "Tenant is suspended", MessageZH: "租户已暂停", HTTPStatus: 403, Category: "tenant"},
				{Code: errno.ErrTenantQuotaExceededCode, Message: "Tenant quota exceeded", MessageZH: "租户配额已超限", HTTPStatus: 403, Category: "quota"},
				{Code: errno.ErrTenantSubscriptionExpiredCode, Message: "Subscription expired", MessageZH: "订阅已过期", HTTPStatus: 403, Category: "subscription"},
			},
		},
		{
			Name:        "配额错误 (Quota Errors)",
			Description: "3xx开头,与配额检查、配额限制相关的错误",
			Errors: []ErrorInfo{
				{Code: errno.ErrQuotaExceededCode, Message: "Quota exceeded", MessageZH: "配额已超限", HTTPStatus: 403, Category: "quota"},
				{Code: errno.ErrQuotaBotExceededCode, Message: "Bot quota exceeded", MessageZH: "Bot数量已超限", HTTPStatus: 403, Category: "quota"},
				{Code: errno.ErrQuotaKnowledgeExceededCode, Message: "Knowledge base quota exceeded", MessageZH: "知识库配额已超限", HTTPStatus: 403, Category: "quota"},
				{Code: errno.ErrQuotaWorkflowExceededCode, Message: "Workflow quota exceeded", MessageZH: "工作流配额已超限", HTTPStatus: 403, Category: "quota"},
				{Code: errno.ErrQuotaAPICallExceededCode, Message: "API call quota exceeded", MessageZH: "API调用次数已超限", HTTPStatus: 429, Category: "quota"},
				{Code: errno.ErrQuotaStorageExceededCode, Message: "Storage quota exceeded", MessageZH: "存储空间已超限", HTTPStatus: 403, Category: "quota"},
			},
		},
		{
			Name:        "订阅错误 (Subscription Errors)",
			Description: "4xx开头,与订阅管理、计费、支付相关的错误",
			Errors: []ErrorInfo{
				{Code: errno.ErrSubscriptionNotFoundCode, Message: "Subscription not found", MessageZH: "订阅不存在", HTTPStatus: 404, Category: "subscription"},
				{Code: errno.ErrSubscriptionExpiredCode, Message: "Subscription expired", MessageZH: "订阅已过期", HTTPStatus: 403, Category: "subscription"},
				{Code: errno.ErrSubscriptionPaymentFailedCode, Message: "Payment failed", MessageZH: "支付失败", HTTPStatus: 402, Category: "subscription"},
				{Code: errno.ErrSubscriptionPaymentRequiredCode, Message: "Payment required", MessageZH: "需要支付", HTTPStatus: 402, Category: "subscription"},
			},
		},
		{
			Name:        "用户错误 (User Errors)",
			Description: "7xx开头,与用户认证、授权相关的错误",
			Errors: []ErrorInfo{
				{Code: errno.ErrUserAuthenticationFailed, Message: "Authentication failed", MessageZH: "认证失败", HTTPStatus: 401, Category: "user"},
				{Code: errno.ErrUserEmailAlreadyExistCode, Message: "Email already exists", MessageZH: "邮箱已存在", HTTPStatus: 409, Category: "user"},
				{Code: errno.ErrUserInfoInvalidateCode, Message: "Invalid email or password", MessageZH: "邮箱或密码无效", HTTPStatus: 401, Category: "user"},
				{Code: errno.ErrUserResourceNotFound, Message: "Resource not found", MessageZH: "资源不存在", HTTPStatus: 404, Category: "user"},
				{Code: errno.ErrUserPermissionCode, Message: "Permission denied", MessageZH: "权限不足", HTTPStatus: 403, Category: "user"},
			},
		},
	}

	return categories
}

// generateMarkdownDoc 生成Markdown文档
func generateMarkdownDoc(categories []ErrorCategory) {
	// 确保目录存在
	os.MkdirAll("docs/errno", 0755)

	file, err := os.Create("docs/errno/error_codes.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 写入文档头部
	fmt.Fprintln(file, "# ZKER 错误码完整清单")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "> **最后更新**: 2025-01-01")
	fmt.Fprintln(file, "> **版本**: v1.0")
	fmt.Fprintln(file, "> **维护者**: 研发B - 后端工程师")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "---")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "## 📋 目录")
	fmt.Fprintln(file, "")

	for _, cat := range categories {
		anchor := strings.ToLower(strings.ReplaceAll(cat.Name, " ", "-"))
		fmt.Fprintf(file, "- [%s](#%s)\n", cat.Name, anchor)
	}

	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "---")
	fmt.Fprintln(file, "")

	// 写入每个分类的错误码
	for _, cat := range categories {
		anchor := strings.ToLower(strings.ReplaceAll(cat.Name, " ", "-"))
		fmt.Fprintf(file, "## %s {#%s}\n", cat.Name, anchor)
		fmt.Fprintln(file, "")
		fmt.Fprintln(file, cat.Description)
		fmt.Fprintln(file, "")
		fmt.Fprintln(file, "| 错误码 | 英文消息 | 中文消息 | HTTP 状态码 | 分类 |")
		fmt.Fprintln(file, "|--------|---------|---------|------------|------|")

		for _, err := range cat.Errors {
			fmt.Fprintf(file, "| %d | %s | %s | %d | %s |\n",
				err.Code, err.Message, err.MessageZH, err.HTTPStatus, err.Category)
		}

		fmt.Fprintln(file, "")
	}

	fmt.Fprintln(file, "---")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "## 🔍 错误码命名规范")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "### 编码规则")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "- **1xx**: 通用错误")
	fmt.Fprintln(file, "- **2xx**: 租户错误")
	fmt.Fprintln(file, "- **3xx**: 配额错误")
	fmt.Fprintln(file, "- **4xx**: 订阅错误")
	fmt.Fprintln(file, "- **5xx**: Bot错误")
	fmt.Fprintln(file, "- **6xx**: 工作流错误")
	fmt.Fprintln(file, "- **7xx**: 用户错误")
	fmt.Fprintln(file, "- **8xx**: 路由错误")
	fmt.Fprintln(file, "- **9xx**: 计费错误")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "### HTTP状态码映射")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "| 错误类型 | HTTP状态码 | 说明 |")
	fmt.Fprintln(file, "|---------|-----------|------|")
	fmt.Fprintln(file, "| 资源不存在 | 404 Not Found | Tenant/Bot/Subscription等资源不存在 |")
	fmt.Fprintln(file, "| 资源冲突 | 409 Conflict | 资源已存在 |")
	fmt.Fprintln(file, "| 权限拒绝 | 403 Forbidden | 租户暂停、配额超限、权限不足 |")
	fmt.Fprintln(file, "| 认证失败 | 401 Unauthorized | 登录失败、token无效 |")
	fmt.Fprintln(file, "| 参数错误 | 400 Bad Request | 请求参数验证失败 |")
	fmt.Fprintln(file, "| 请求过多 | 429 Too Many Requests | API调用频率超限 |")
	fmt.Fprintln(file, "| 需要支付 | 402 Payment Required | 订阅过期、需要续费 |")
	fmt.Fprintln(file, "| 服务器错误 | 500 Internal Server Error | 服务器内部错误 |")
	fmt.Fprintln(file, "")

	fmt.Println("✅ Markdown文档生成完成: docs/errno/error_codes.md")
}

// generateJSONDoc 生成JSON文档
func generateJSONDoc(categories []ErrorCategory) {
	// 生成英文版JSON
	enData := make(map[int32]ErrorInfo)
	for _, cat := range categories {
		for _, err := range cat.Errors {
			enData[err.Code] = err
		}
	}

	enJSON, _ := json.MarshalIndent(enData, "", "  ")
	os.WriteFile("docs/errno/error_codes.json", enJSON, 0644)

	// 生成中文版JSON
	zhData := make(map[int32]map[string]string)
	for _, cat := range categories {
		for _, err := range cat.Errors {
			zhData[err.Code] = map[string]string{
				"message_zh": err.MessageZH,
				"http_status": fmt.Sprintf("%d", err.HTTPStatus),
				"category":    err.Category,
			}
		}
	}

	zhJSON, _ := json.MarshalIndent(zhData, "", "  ")
	os.WriteFile("docs/errno/error_codes_zh.json", zhJSON, 0644)

	fmt.Println("✅ JSON文档生成完成:")
	fmt.Println("  - docs/errno/error_codes.json (英文版)")
	fmt.Println("  - docs/errno/error_codes_zh.json (中文版)")
}

// generateErrorCodesJSON 生成错误码JSON文件(供前端使用)
func generateErrorCodesJSON() {
	categories := collectErrorCategories()

	allErrors := []ErrorInfo{}
	for _, cat := range categories {
		allErrors = append(allErrors, cat.Errors...)
	}

	// 生成完整错误码列表
	jsonData, _ := json.MarshalIndent(allErrors, "", "  ")
	os.WriteFile("docs/errno/all_error_codes.json", jsonData, 0644)

	fmt.Println("✅ 错误码JSON文件生成完成: docs/errno/all_error_codes.json")
}

// generateErrorTestCases 生成错误码测试用例
func generateErrorTestCases() {
	testTemplate := `/*
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

package errno

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name       string
		code       int32
		message    string
		messageZH  string
		httpStatus int
	}{
		{{range .Categories}}{{range .Errors}}
		{
			name:       "{{.Message}}",
			code:       {{.Code}},
			message:    "{{.Message}}",
			messageZH:  "{{.MessageZH}}",
			httpStatus: {{.HTTPStatus}},
		},{{end}}{{end}}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建增强错误
			err := NewEnhancedError(tt.code, tt.message, tt.messageZH)

			// 验证错误码
			assert.Equal(t, tt.code, err.Code)
			assert.Equal(t, tt.message, err.Message)
			assert.Equal(t, tt.messageZH, err.MessageZH)
			assert.Equal(t, tt.httpStatus, err.GetHTTPStatus())

			// 验证JSON序列化
			jsonData := err.ToJSON()
			assert.NotNil(t, jsonData)

			// 验证Error()方法
			errorStr := err.Error()
			assert.Contains(t, errorStr, tt.message)
		})
	}
}

func TestHTTPStatusMapping(t *testing.T) {
	tests := []struct {
		name       string
		code       int32
		httpStatus int
	}{
		{{range .Categories}}{{range .Errors}}
		{
			name:       "{{.Message}}",
			code:       {{.Code}},
			httpStatus: {{.HTTPStatus}},
		},{{end}}{{end}}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewEnhancedError(tt.code, "test", "测试")
			assert.Equal(t, tt.httpStatus, err.GetHTTPStatus())
		})
	}
}
`

	// 准备模板数据
	data := struct {
		Categories []ErrorCategory
	}{
		Categories: collectErrorCategories(),
	}

	// 解析模板
	tmpl, err := template.New("test").Parse(testTemplate)
	if err != nil {
		panic(err)
	}

	// 生成测试文件
	file, err := os.Create("types/errno/error_codes_test.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ 错误码测试用例生成完成: types/errno/error_codes_test.go")
}
