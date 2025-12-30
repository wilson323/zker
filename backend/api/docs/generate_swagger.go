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
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SwaggerConfig Swagger配置
type SwaggerConfig struct {
	InputDir  string `json:"input_dir"`  // 输入目录（Handler文件）
	OutputDir string `json:"output_dir"` // 输出目录
	Format    string `json:"format"`    // 输出格式（yaml/json）
	Version   string `json:"version"`   // API版本
	BasePath  string `json:"base_path"`  // API基础路径
}

// OpenAPISpec OpenAPI规范结构
type OpenAPISpec struct {
	OpenAPI    string                 `yaml:"openapi"`
	Info       Info                   `yaml:"info"`
	Servers    []Server               `yaml:"servers"`
	Tags       []Tag                  `yaml:"tags"`
	Security   []SecurityRequirement  `yaml:"security"`
	Paths      map[string]PathItem    `yaml:"paths"`
	Components Components             `yaml:"components"`
}

type Info struct {
	Title          string   `yaml:"title"`
	Version        string   `yaml:"version"`
	Description    string   `yaml:"description"`
	TermsOfService string   `yaml:"termsOfService,omitempty"`
	Contact        *Contact `yaml:"contact,omitempty"`
	License        *License `yaml:"license,omitempty"`
}

type Contact struct {
	Name  string `yaml:"name,omitempty"`
	Email string `yaml:"email,omitempty"`
	URL   string `yaml:"url,omitempty"`
}

type License struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url,omitempty"`
}

type Server struct {
	URL         string            `yaml:"url"`
	Description string            `yaml:"description"`
	Variables   map[string]Variable `yaml:"variables,omitempty"`
}

type Variable struct {
	Default     interface{} `yaml:"default"`
	Description string      `yaml:"description,omitempty"`
	Enum        []string    `yaml:"enum,omitempty"`
}

type Tag struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type SecurityRequirement struct {
	AuthType []string `yaml:"-"`
}

type PathItem struct {
	Ref     string    `yaml:"$ref,omitempty"`
	Get     *Operation `yaml:"get,omitempty"`
	Put     *Operation `yaml:"put,omitempty"`
	Post    *Operation `yaml:"post,omitempty"`
	Delete  *Operation `yaml:"delete,omitempty"`
	Options *Operation `yaml:"options,omitempty"`
	Head    *Operation `yaml:"head,omitempty"`
	Patch   *Operation `yaml:"patch,omitempty"`
}

type Operation struct {
	Tags        []string            `yaml:"tags"`
	Summary     string              `yaml:"summary"`
	Description string              `yaml:"description,omitempty"`
	OperationID string              `yaml:"operationId,omitempty"`
	Parameters  []Parameter         `yaml:"parameters,omitempty"`
	RequestBody *RequestBody        `yaml:"requestBody,omitempty"`
	Responses   map[string]Response `yaml:"responses"`
	Security    []map[string][]string `yaml:"security,omitempty"`
}

type Parameter struct {
	Name            string      `yaml:"name"`
	In              string      `yaml:"in"`
	Description     string      `yaml:"description,omitempty"`
	Required        bool        `yaml:"required"`
	Deprecated      bool        `yaml:"deprecated,omitempty"`
	AllowEmptyValue bool        `yaml:"allowEmptyValue,omitempty"`
	Style           string      `yaml:"style,omitempty"`
	Explode         bool        `yaml:"explode,omitempty"`
	AllowReserved   bool        `yaml:"allowReserved,omitempty"`
	Schema          *Schema     `yaml:"schema,omitempty"`
	Example         interface{} `yaml:"example,omitempty"`
	Examples        map[string]Example `yaml:"examples,omitempty"`
	Content         map[string]MediaType `yaml:"content,omitempty"`
}

type Schema struct {
	Ref                  string                 `yaml:"$ref,omitempty"`
	Type                 string                 `yaml:"type,omitempty"`
	Format               string                 `yaml:"format,omitempty"`
	Title                string                 `yaml:"title,omitempty"`
	Description          string                 `yaml:"description,omitempty"`
	Default              interface{}            `yaml:"default,omitempty"`
	MultipleOf           float64                `yaml:"multipleOf,omitempty"`
	Maximum              *float64               `yaml:"maximum,omitempty"`
	ExclusiveMaximum     bool                   `yaml:"exclusiveMaximum,omitempty"`
	Minimum              *float64               `yaml:"minimum,omitempty"`
	ExclusiveMinimum     bool                   `yaml:"exclusiveMinimum,omitempty"`
	MaxLength            *int                   `yaml:"maxLength,omitempty"`
	MinLength            *int                   `yaml:"minLength,omitempty"`
	Pattern              string                 `yaml:"pattern,omitempty"`
	MaxItems             *int                   `yaml:"maxItems,omitempty"`
	MinItems             *int                   `yaml:"minItems,omitempty"`
	UniqueItems          bool                   `yaml:"uniqueItems,omitempty"`
	MaxProperties        *int                   `yaml:"maxProperties,omitempty"`
	MinProperties        *int                   `yaml:"minProperties,omitempty"`
	Required             []string               `yaml:"required,omitempty"`
	Enum                 []interface{}          `yaml:"enum,omitempty"`
	AllOf                []Schema               `yaml:"allOf,omitempty"`
	AnyOf                []Schema               `yaml:"anyOf,omitempty"`
	OneOf                []Schema               `yaml:"oneOf,omitempty"`
	Not                  *Schema                `yaml:"not,omitempty"`
	Items                *Schema                `yaml:"items,omitempty"`
	Properties           map[string]Schema      `yaml:"properties,omitempty"`
	AdditionalProperties *Schema                `yaml:"additionalProperties,omitempty"`
	Discriminator        *Discriminator         `yaml:"discriminator,omitempty"`
	ReadOnly             bool                   `yaml:"readOnly,omitempty"`
	WriteOnly            bool                   `yaml:"writeOnly,omitempty"`
	XML                  *XML                   `yaml:"xml,omitempty"`
	ExternalDocs         *ExternalDocumentation `yaml:"externalDocs,omitempty"`
	Example              interface{}            `yaml:"example,omitempty"`
}

type MediaType struct {
	Schema   *Schema            `yaml:"schema,omitempty"`
	Example  interface{}        `yaml:"example,omitempty"`
	Examples map[string]Example `yaml:"examples,omitempty"`
	Encoding map[string]Encoding `yaml:"encoding,omitempty"`
}

type RequestBody struct {
	Description string                `yaml:"description,omitempty"`
	Required    bool                  `yaml:"required"`
	Content     map[string]MediaType  `yaml:"content"`
}

type Response struct {
	Ref         string               `yaml:"$ref,omitempty"`
	Description string               `yaml:"description"`
	Headers     map[string]Header    `yaml:"headers,omitempty"`
	Content     map[string]MediaType `yaml:"content,omitempty"`
	Links       map[string]Link      `yaml:"links,omitempty"`
}

type Header struct {
	Description string  `yaml:"description,omitempty"`
	Required    bool    `yaml:"required"`
	Deprecated  bool    `yaml:"deprecated,omitempty"`
	AllowEmptyValue bool `yaml:"allowEmptyValue,omitempty"`
	Schema      *Schema `yaml:"schema,omitempty"`
	Example     interface{} `yaml:"example,omitempty"`
}

type Encoding struct {
	ContentType   string            `yaml:"contentType,omitempty"`
	Headers       map[string]Header `yaml:"headers,omitempty"`
	Style         string            `yaml:"style,omitempty"`
	Explode       bool              `yaml:"explode,omitempty"`
	AllowReserved bool              `yaml:"allowReserved,omitempty"`
}

type Link struct {
	Href     string            `yaml:"href"`
	Ref      string            `yaml:"$ref,omitempty"`
	OperationRef string        `yaml:"operationRef,omitempty"`
	OperationID  string        `yaml:"operationId,omitempty"`
	Parameters   map[string]interface{} `yaml:"parameters,omitempty"`
	RequestBody  interface{}   `yaml:"requestBody,omitempty"`
	Server       *Server       `yaml:"server,omitempty"`
	Description  string        `yaml:"description,omitempty"`
}

type Example struct {
	Summary     string      `yaml:"summary,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Value       interface{} `yaml:"value"`
	ExternalValue string    `yaml:"externalValue,omitempty"`
}

type Discriminator struct {
	PropertyName string            `yaml:"propertyName"`
	Mapping      map[string]string `yaml:"mapping,omitempty"`
}

type XML struct {
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
	Prefix    string `yaml:"prefix,omitempty"`
	Attribute bool   `yaml:"attribute,omitempty"`
	Wrapped   bool   `yaml:"wrapped,omitempty"`
}

type ExternalDocumentation struct {
	Description string `yaml:"description,omitempty"`
	URL         string `yaml:"url"`
}

type Components struct {
	Schemas         map[string]Schema           `yaml:"schemas,omitempty"`
	Responses       map[string]Response         `yaml:"responses,omitempty"`
	Parameters      map[string]Parameter        `yaml:"parameters,omitempty"`
	Examples        map[string]Example          `yaml:"examples,omitempty"`
	RequestBodies   map[string]RequestBody     `yaml:"requestBodies,omitempty"`
	SecuritySchemes map[string]SecurityScheme   `yaml:"securitySchemes,omitempty"`
	Links           map[string]Link             `yaml:"links,omitempty"`
	Callbacks       map[string]Callback         `yaml:"callbacks,omitempty"`
}

type SecurityScheme struct {
	Type             string            `yaml:"type"`
	Description      string            `yaml:"description,omitempty"`
	Name             string            `yaml:"name,omitempty"`
	In               string            `yaml:"in,omitempty"`
	Scheme           string            `yaml:"scheme,omitempty"`
	BearerFormat     string            `yaml:"bearerFormat,omitempty"`
	Flows            map[string]OAuthFlow `yaml:"flows,omitempty"`
	OpenIdConnectUrl string            `yaml:"openIdConnectUrl,omitempty"`
}

type OAuthFlow struct {
	AuthorizationUrl string            `yaml:"authorizationUrl"`
	TokenUrl         string            `yaml:"tokenUrl"`
	RefreshUrl       string            `yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `yaml:"scopes"`
}

type Callback struct {
	Expression string              `yaml:"expression"`
	Paths      map[string]PathItem `yaml:"-$key,omitempty"`
}

func main() {
	// 默认配置（使用绝对路径或相对于当前工作目录的路径）
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	config := SwaggerConfig{
		InputDir:  filepath.Join(filepath.Dir(wd), "handler", "coze", "org"),
		OutputDir: filepath.Join(wd, "swagger"),
		Format:    "yaml",
		Version:   "v1.0.0",
		BasePath:  "/api",
	}

	// 从命令行参数读取配置
	if len(os.Args) > 1 {
		configFile := os.Args[1]
		data, err := os.ReadFile(configFile)
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			os.Exit(1)
		}

		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Printf("Error parsing config file: %v\n", err)
			os.Exit(1)
		}
	}

	// 生成OpenAPI规范
	fmt.Printf("📂 Scanning handler files in: %s\n", config.InputDir)
	spec, err := generateOpenAPISpec(config)
	if err != nil {
		fmt.Printf("❌ Error generating OpenAPI spec: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Found %d API paths\n", len(spec.Paths))

	// 创建输出目录
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// 写入主格式文件
	outputFile := filepath.Join(config.OutputDir, "swagger."+config.Format)
	if err := writeSpec(spec, outputFile, config.Format); err != nil {
		fmt.Printf("Error writing spec file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ OpenAPI spec generated successfully: %s\n", outputFile)

	// 同时生成另一种格式
	otherFormat := "yaml"
	if config.Format == "yaml" || config.Format == "yml" {
		otherFormat = "json"
	}
	otherFile := filepath.Join(config.OutputDir, "swagger."+otherFormat)
	if err := writeSpec(spec, otherFile, otherFormat); err != nil {
		fmt.Printf("⚠️  Warning: could not generate %s format: %v\n", otherFormat, err)
	} else {
		fmt.Printf("✅ Also generated: %s\n", otherFile)
	}

	// 验证规范
	fmt.Println("\n🔍 Validating OpenAPI spec...")
	if err := validateSpec(spec); err != nil {
		fmt.Printf("⚠️  Warning: spec validation failed: %v\n", err)
		// 统计有多少路径没有操作
		emptyPaths := 0
		for path, pathItem := range spec.Paths {
			if pathItem.Get == nil && pathItem.Post == nil && pathItem.Put == nil &&
				pathItem.Delete == nil && pathItem.Patch == nil {
				emptyPaths++
				if emptyPaths <= 5 { // 只显示前5个
					fmt.Printf("  - Empty path: %s\n", path)
				}
			}
		}
		if emptyPaths > 5 {
			fmt.Printf("  ... and %d more empty paths\n", emptyPaths-5)
		}
	} else {
		fmt.Println("✅ Spec validation passed")
	}

	fmt.Println("\n📖 Next steps:")
	fmt.Printf("  1. View generated spec: %s\n", outputFile)
	if config.Format == "yaml" || config.Format == "yml" {
		fmt.Printf("  2. Convert to JSON for tools: Use any YAML-to-JSON converter\n")
	}
	fmt.Printf("  3. Import to Swagger UI: https://editor.swagger.io/\n")
	fmt.Printf("  4. View in Redoc: https://redocly.github.io/redoc/\n")
}

// generateOpenAPISpec 生成OpenAPI规范
func generateOpenAPISpec(config SwaggerConfig) (*OpenAPISpec, error) {
	spec := &OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       "组织中心管理API",
			Version:     strings.TrimPrefix(config.Version, "v"),
			Description: "企业级组织管理系统的RESTful API文档",
			Contact: &Contact{
				Name:  "API Support",
				Email: "api-support@coze.com",
			},
			License: &License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		Servers: []Server{
			{
				URL:         "http://localhost:8080",
				Description: "开发环境",
			},
			{
				URL:         "https://api-test.coze.com",
				Description: "测试环境",
			},
			{
				URL:         "https://api.coze.com",
				Description: "生产环境",
			},
		},
		Tags: []Tag{
			{Name: "组织管理", Description: "组织架构管理"},
			{Name: "部门管理", Description: "部门层级关系维护"},
			{Name: "员工管理", Description: "员工全生命周期管理"},
			{Name: "岗位管理", Description: "岗位定义和职级管理"},
			{Name: "通讯录服务", Description: "组织目录查询和员工搜索"},
		},
		Security: []SecurityRequirement{
			{AuthType: []string{"BearerAuth"}},
		},
		Paths:      make(map[string]PathItem),
		Components: generateComponents(),
	}

	// 扫描Handler文件并生成Paths
	if err := scanHandlerFiles(config.InputDir, spec); err != nil {
		return nil, fmt.Errorf("error scanning handler files: %w", err)
	}

	return spec, nil
}

// generateComponents 生成组件定义
func generateComponents() Components {
	return Components{
		SecuritySchemes: map[string]SecurityScheme{
			"BearerAuth": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "JWT Token认证",
			},
			"ApiKeyAuth": {
				Type:        "apiKey",
				In:          "header",
				Name:        "X-API-Key",
				Description: "API Key认证",
			},
		},
		Schemas: generateSchemas(),
		Responses: generateCommonResponses(),
	}
}

// generateSchemas 生成Schema定义
func generateSchemas() map[string]Schema {
	return map[string]Schema{
		"ErrorResponse": {
			Type: "object",
			Required: []string{"code", "message"},
			Properties: map[string]Schema{
				"code": {
					Type:        "integer",
					Format:      "int32",
					Description: "错误码",
				},
				"message": {
					Type:        "string",
					Description: "错误信息（中文）",
				},
				"message_en": {
					Type:        "string",
					Description: "错误信息（英文）",
				},
				"data": {
					Type:        "object",
					Description: "额外错误详情",
				},
			},
		},
		"OrganizationData": {
			Type: "object",
			Required: []string{"org_id", "tenant_id", "org_name", "org_type", "org_code"},
			Properties: map[string]Schema{
				"org_id": {
					Type:        "string",
					Description: "组织ID",
				},
				"tenant_id": {
					Type:        "string",
					Description: "租户ID",
				},
				"org_name": {
					Type:        "string",
					Description: "组织名称",
				},
				"org_type": {
					Type:        "string",
					Enum:        []interface{}{"company", "division", "department", "project"},
					Description: "组织类型",
				},
				"parent_id": {
					Type:        "string",
					Description: "父组织ID",
				},
				"org_code": {
					Type:        "string",
					Description: "组织编码",
				},
				"level": {
					Type:        "integer",
					Description: "组织层级",
				},
				"path": {
					Type:        "string",
					Description: "组织路径",
				},
				"status": {
					Type:        "string",
					Enum:        []interface{}{"active", "inactive", "frozen"},
					Description: "组织状态",
				},
				"description": {
					Type:        "string",
					Description: "组织描述",
				},
				"leader_id": {
					Type:        "string",
					Description: "负责人ID",
				},
				"created_at": {
					Type:        "integer",
					Format:      "int64",
					Description: "创建时间",
				},
				"updated_at": {
					Type:        "integer",
					Format:      "int64",
					Description: "更新时间",
				},
			},
		},
		// 可以继续添加其他Schema定义...
	}
}

// generateCommonResponses 生成通用响应定义
func generateCommonResponses() map[string]Response {
	return map[string]Response{
		"BadRequest": {
			Description: "请求参数错误",
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{
						Ref: "#/components/schemas/ErrorResponse",
					},
				},
			},
		},
		"Unauthorized": {
			Description: "未授权",
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{
						Ref: "#/components/schemas/ErrorResponse",
					},
				},
			},
		},
		"Forbidden": {
			Description: "权限不足",
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{
						Ref: "#/components/schemas/ErrorResponse",
					},
				},
			},
		},
		"NotFound": {
			Description: "资源不存在",
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{
						Ref: "#/components/schemas/ErrorResponse",
					},
				},
			},
		},
		"InternalServerError": {
			Description: "服务器内部错误",
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{
						Ref: "#/components/schemas/ErrorResponse",
					},
				},
			},
		},
	}
}

// scanHandlerFiles 扫描Handler文件并生成API路径定义
func scanHandlerFiles(inputDir string, spec *OpenAPISpec) error {
	// 遍历Handler文件
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和测试文件
		if info.IsDir() || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// 只处理.go文件
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// 读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		// 解析注释中的API定义
		if err := parseAPIAnnotations(string(content), spec); err != nil {
			fmt.Printf("Warning: error parsing %s: %v\n", path, err)
		}

		return nil
	})

	return err
}

// parseAPIAnnotations 解析文件中的API注释
func parseAPIAnnotations(content string, spec *OpenAPISpec) error {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		// 查找@router注释
		if strings.Contains(line, "@router") {
			routerParts := strings.Fields(line)
			if len(routerParts) < 3 {
				continue
			}

			method := routerParts[1]
			path := strings.Trim(routerParts[2], `"[]`)

			// 解析后续的注释行
			var operation Operation
			for j := i + 1; j < len(lines) && j < i+20; j++ {
				nextLine := strings.TrimSpace(lines[j])
				if !strings.HasPrefix(nextLine, "//") {
					break
				}

				commentLine := strings.TrimPrefix(nextLine, "//")
				commentLine = strings.TrimSpace(commentLine)

				if strings.HasPrefix(commentLine, "@summary") {
					operation.Summary = strings.TrimSpace(strings.TrimPrefix(commentLine, "@summary"))
				} else if strings.HasPrefix(commentLine, "@description") {
					operation.Description = strings.TrimSpace(strings.TrimPrefix(commentLine, "@description"))
				} else if strings.HasPrefix(commentLine, "@tags") {
					tags := strings.Fields(strings.TrimPrefix(commentLine, "@tags"))
					operation.Tags = tags
				} else if strings.HasPrefix(commentLine, "@param") {
					// TODO: 解析参数定义
				} else if strings.HasPrefix(commentLine, "@success") {
					// TODO: 解析成功响应
				} else if strings.HasPrefix(commentLine, "@failure") {
					// TODO: 解析失败响应
				}
			}

			// 确保有默认值
			if operation.Tags == nil {
				operation.Tags = []string{"其他"}
			}
			if operation.Summary == "" {
				operation.Summary = method + " " + path
			}
			if operation.Responses == nil {
				operation.Responses = map[string]Response{
					"200": {
						Description: "成功",
						Content: map[string]MediaType{
							"application/json": {
								Schema: &Schema{
									Type: "object",
								},
							},
						},
					},
				}
			}

			// 添加到Paths
			pathItem := spec.Paths[path]
			switch strings.ToUpper(method) {
			case "GET":
				pathItem.Get = &operation
			case "POST":
				pathItem.Post = &operation
			case "PUT":
				pathItem.Put = &operation
			case "DELETE":
				pathItem.Delete = &operation
			case "PATCH":
				pathItem.Patch = &operation
			}
			spec.Paths[path] = pathItem
		}
	}

	return nil
}

// writeSpec 将OpenAPI规范写入文件
func writeSpec(spec *OpenAPISpec, filename, format string) error {
	var data []byte
	var err error

	if format == "yaml" || format == "yml" {
		data, err = yaml.Marshal(spec)
		if err != nil {
			return fmt.Errorf("error marshaling to YAML: %w", err)
		}
	} else {
		data, err = json.MarshalIndent(spec, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling to JSON: %w", err)
		}
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

// validateSpec 验证OpenAPI规范
func validateSpec(spec *OpenAPISpec) error {
	// 基本验证
	if spec.OpenAPI == "" {
		return fmt.Errorf("missing openapi version")
	}

	if spec.Info.Title == "" {
		return fmt.Errorf("missing API title")
	}

	if spec.Info.Version == "" {
		return fmt.Errorf("missing API version")
	}

	// 验证至少有一个路径
	if len(spec.Paths) == 0 {
		return fmt.Errorf("no API paths defined")
	}

	// 验证每个路径都有操作
	for path, pathItem := range spec.Paths {
		if pathItem.Get == nil && pathItem.Post == nil && pathItem.Put == nil &&
			pathItem.Delete == nil && pathItem.Patch == nil {
			return fmt.Errorf("path %s has no operations", path)
		}
	}

	return nil
}
