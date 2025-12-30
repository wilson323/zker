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
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
	"github.com/coze-dev/coze-studio/backend/domain/developer/repository"
)

// SDKGeneratorService SDK生成服务接口
type SDKGeneratorService interface {
	// GenerateSDK 为项目生成SDK
	GenerateSDK(ctx context.Context, req *GenerateSDKRequest) (*entity.SDK, error)

	// GetSDK 获取SDK详情
	GetSDK(ctx context.Context, sdkID string) (*entity.SDK, error)

	// ListSDKs 查询SDK列表
	ListSDKs(ctx context.Context, filter *SDKListFilter) ([]*entity.SDK, int64, error)

	// PublishSDK 发布SDK
	PublishSDK(ctx context.Context, sdkID string) error

	// DeleteSDK 删除SDK
	DeleteSDK(ctx context.Context, sdkID string) error

	// DownloadSDK 下载SDK包
	DownloadSDK(ctx context.Context, sdkID string) (*entity.SDKPackage, error)

	// GetProjectSDKs 获取项目的SDK列表
	GetProjectSDKs(ctx context.Context, projectID string) ([]*entity.SDK, error)

	// GenerateSDKCode 生成SDK代码
	GenerateSDKCode(ctx context.Context, language entity.SDKLanguage, project *entity.Project) (string, error)
}

// GenerateSDKRequest 生成SDK请求
type GenerateSDKRequest struct {
	TenantID    string                `json:"tenant_id" validate:"required"`
	ProjectID   string                `json:"project_id" validate:"required"`
	SDKName     string                `json:"sdk_name" validate:"required,max=100"`
	Language    entity.SDKLanguage    `json:"language" validate:"required"`
	Version     string                `json:"version" validate:"required"`
	Description string                `json:"description"`
}

// SDKListFilter SDK列表过滤器
type SDKListFilter struct {
	TenantID  string              `validate:"required"`
	ProjectID string
	Language  entity.SDKLanguage
	Status    entity.SDKStatus
	PageToken string
	PageSize  int
}

// sdkGeneratorService SDK生成服务实现
type sdkGeneratorService struct {
	sdkRepo     repository.SDKRepository
	projectRepo repository.ProjectRepository
	templates   map[entity.SDKLanguage]*entity.SDKTemplate
}

// NewSDKGeneratorService 创建SDK生成服务实例
func NewSDKGeneratorService(
	sdkRepo repository.SDKRepository,
	projectRepo repository.ProjectRepository,
) SDKGeneratorService {
	service := &sdkGeneratorService{
		sdkRepo:     sdkRepo,
		projectRepo: projectRepo,
		templates:   make(map[entity.SDKLanguage]*entity.SDKTemplate),
	}

	// 初始化模板
	service.initTemplates()

	return service
}

// initTemplates 初始化SDK模板
func (s *sdkGeneratorService) initTemplates() {
	// Python模板
	s.templates[entity.SDKLanguagePython] = &entity.SDKTemplate{
		Language:     entity.SDKLanguagePython,
		TemplateName: "python",
		Dependencies: []string{"requests>=2.28.0", "pydantic>=1.10.0"},
	}

	// JavaScript/TypeScript模板
	s.templates[entity.SDKLanguageJavaScript] = &entity.SDKTemplate{
		Language:     entity.SDKLanguageJavaScript,
		TemplateName: "typescript",
		Dependencies: []string{"axios@^1.4.0", "@types/node@^20.0.0"},
	}

	// Go模板
	s.templates[entity.SDKLanguageGo] = &entity.SDKTemplate{
		Language:     entity.SDKLanguageGo,
		TemplateName: "go",
		Dependencies: []string{"github.com/valyala/fasthttp"},
	}

	// Java模板
	s.templates[entity.SDKLanguageJava] = &entity.SDKTemplate{
		Language:     entity.SDKLanguageJava,
		TemplateName: "java",
		Dependencies: []string{"com.squareup.okhttp3:okhttp:4.11.0"},
	}
}

// GenerateSDK 为项目生成SDK
func (s *sdkGeneratorService) GenerateSDK(ctx context.Context, req *GenerateSDKRequest) (*entity.SDK, error) {
	// 获取项目信息
	project, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	// 检查项目是否属于该租户
	if project.TenantID != req.TenantID {
		return nil, fmt.Errorf("project does not belong to tenant")
	}

	// 生成SDK代码
	code, err := s.GenerateSDKCode(ctx, req.Language, project)
	if err != nil {
		return nil, fmt.Errorf("failed to generate sdk code: %w", err)
	}

	// 生成README
	readme := s.generateReadme(req.Language, project)

	sdk := &entity.SDK{
		SDKID:       generateID("sdk"),
		TenantID:    req.TenantID,
		ProjectID:   req.ProjectID,
		SDKName:     req.SDKName,
		Language:    req.Language,
		Version:     req.Version,
		Description: req.Description,
		Code:        code,
		Readme:      readme,
		Status:      entity.SDKStatusDraft,
	}

	if err := s.sdkRepo.Create(ctx, sdk); err != nil {
		return nil, err
	}

	return sdk, nil
}

// GetSDK 获取SDK详情
func (s *sdkGeneratorService) GetSDK(ctx context.Context, sdkID string) (*entity.SDK, error) {
	return s.sdkRepo.GetByID(ctx, sdkID)
}

// ListSDKs 查询SDK列表
func (s *sdkGeneratorService) ListSDKs(ctx context.Context, filter *SDKListFilter) ([]*entity.SDK, int64, error) {
	repoFilter := &repository.SDKFilter{
		TenantID:  filter.TenantID,
		ProjectID: filter.ProjectID,
		Language:  filter.Language,
		Status:    filter.Status,
		PageToken: filter.PageToken,
		PageSize:  filter.PageSize,
	}

	return s.sdkRepo.List(ctx, repoFilter)
}

// PublishSDK 发布SDK
func (s *sdkGeneratorService) PublishSDK(ctx context.Context, sdkID string) error {
	return s.sdkRepo.UpdateStatus(ctx, sdkID, entity.SDKStatusPublished)
}

// DeleteSDK 删除SDK
func (s *sdkGeneratorService) DeleteSDK(ctx context.Context, sdkID string) error {
	return s.sdkRepo.Delete(ctx, sdkID)
}

// DownloadSDK 下载SDK包
func (s *sdkGeneratorService) DownloadSDK(ctx context.Context, sdkID string) (*entity.SDKPackage, error) {
	sdk, err := s.sdkRepo.GetByID(ctx, sdkID)
	if err != nil {
		return nil, err
	}
	if sdk == nil {
		return nil, fmt.Errorf("sdk not found")
	}

	// 增加下载次数
	_ = s.sdkRepo.IncrementDownloadCount(ctx, sdkID)

	// TODO: 打包SDK代码并生成下载链接
	pkg := &entity.SDKPackage{
		SDKID:       sdkID,
		FileName:    fmt.Sprintf("%s-%s-%s.zip", sdk.SDKName, sdk.Language, sdk.Version),
		FileSize:    int64(len(sdk.Code)),
		DownloadURL: fmt.Sprintf("/api/developer/sdk/%s/download", sdkID),
	}

	return pkg, nil
}

// GetProjectSDKs 获取项目的SDK列表
func (s *sdkGeneratorService) GetProjectSDKs(ctx context.Context, projectID string) ([]*entity.SDK, error) {
	return s.sdkRepo.GetByProjectID(ctx, projectID)
}

// GenerateSDKCode 生成SDK代码
func (s *sdkGeneratorService) GenerateSDKCode(ctx context.Context, language entity.SDKLanguage, project *entity.Project) (string, error) {
	template, ok := s.templates[language]
	if !ok {
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	// 根据语言生成代码
	switch language {
	case entity.SDKLanguagePython:
		return s.generatePythonSDK(project, template)
	case entity.SDKLanguageJavaScript:
		return s.generateJavaScriptSDK(project, template)
	case entity.SDKLanguageGo:
		return s.generateGoSDK(project, template)
	case entity.SDKLanguageJava:
		return s.generateJavaSDK(project, template)
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
}

// generatePythonSDK 生成Python SDK
func (s *sdkGeneratorService) generatePythonSDK(project *entity.Project, template *entity.SDKTemplate) (string, error) {
	return fmt.Sprintf(`"""
%s SDK for %s
Auto-generated on %s
"""

import requests
from typing import Optional, Dict, Any
from pydantic import BaseModel

class CozeClient:
    """Coze Studio API Client"""

    def __init__(self, api_key: str, base_url: str = "https://api.coze.com"):
        self.api_key = api_key
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json"
        })

    def create_bot(self, name: str, description: Optional[str] = None) -> Dict[str, Any]:
        """Create a new bot"""
        url = f"{self.base_url}/api/bot/create"
        payload = {"name": name}
        if description:
            payload["description"] = description

        response = self.session.post(url, json=payload)
        response.raise_for_status()
        return response.json()

    def get_bot(self, bot_id: str) -> Dict[str, Any]:
        """Get bot details"""
        url = f"{self.base_url}/api/bot/get"
        response = self.session.post(url, json={"bot_id": bot_id})
        response.raise_for_status()
        return response.json()

    def send_message(self, bot_id: str, message: str) -> Dict[str, Any]:
        """Send message to bot"""
        url = f"{self.base_url}/api/message/send"
        response = self.session.post(url, json={
            "bot_id": bot_id,
            "message": message
        })
        response.raise_for_status()
        return response.json()
`, project.ProjectName, project.ProjectName, time.Now().Format("2006-01-02")), nil
}

// generateJavaScriptSDK 生成JavaScript/TypeScript SDK
func (s *sdkGeneratorService) generateJavaScriptSDK(project *entity.Project, template *entity.SDKTemplate) (string, error) {
	return fmt.Sprintf(`/**
 * %s SDK for %s
 * Auto-generated on %s
 */

import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

export interface CozeClientConfig {
  apiKey: string;
  baseURL?: string;
}

export class CozeClient {
  private client: AxiosInstance;

  constructor(config: CozeClientConfig) {
    const { apiKey, baseURL = 'https://api.coze.com' } = config;

    this.client = axios.create({
      baseURL,
      headers: {
        'Authorization': \`Bearer \${apiKey}\`,
        'Content-Type': 'application/json',
      },
    });
  }

  async createBot(params: { name: string; description?: string }): Promise<any> {
    const response = await this.client.post('/api/bot/create', params);
    return response.data;
  }

  async getBot(botId: string): Promise<any> {
    const response = await this.client.post('/api/bot/get', { bot_id: botId });
    return response.data;
  }

  async sendMessage(params: { botId: string; message: string }): Promise<any> {
    const response = await this.client.post('/api/message/send', {
      bot_id: params.botId,
      message: params.message,
    });
    return response.data;
  }
}

export default CozeClient;
`, project.ProjectName, project.ProjectName, time.Now().Format("2006-01-02")), nil
}

// generateGoSDK 生成Go SDK
func (s *sdkGeneratorService) generateGoSDK(project *entity.Project, template *entity.SDKTemplate) (string, error) {
	return fmt.Sprintf(`// %s SDK for %s
// Auto-generated on %s

package coze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultBaseURL = "https://api.coze.com"
)

type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		client:  &http.Client{},
	}
}

func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

type CreateBotRequest struct {
	Name        string \`json:"name"\`
	Description string \`json:"description,omitempty"\`
}

type CreateBotResponse struct {
	BotID string \`json:"bot_id"\`
}

func (c *Client) CreateBot(req *CreateBotRequest) (*CreateBotResponse, error) {
	url := c.baseURL + "/api/bot/create"

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CreateBotResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetBot(botID string) (map[string]interface{}, error) {
	url := c.baseURL + "/api/bot/get"

	reqBody := map[string]string{"bot_id": botID}
	body, _ := json.Marshal(reqBody)

	httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
`, project.ProjectName, project.ProjectName, time.Now().Format("2006-01-02")), nil
}

// generateJavaSDK 生成Java SDK
func (s *sdkGeneratorService) generateJavaSDK(project *entity.Project, template *entity.SDKTemplate) (string, error) {
	return fmt.Sprintf(`/**
 * %s SDK for %s
 * Auto-generated on %s
 */

package com.coze.sdk;

import okhttp3.*;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import java.io.IOException;
import java.util.Map;
import java.util.HashMap;

public class CozeClient {
    private final String apiKey;
    private final String baseUrl;
    private final OkHttpClient client;
    private final Gson gson;

    public CozeClient(String apiKey) {
        this(apiKey, "https://api.coze.com");
    }

    public CozeClient(String apiKey, String baseUrl) {
        this.apiKey = apiKey;
        this.baseUrl = baseUrl;
        this.client = new OkHttpClient();
        this.gson = new GsonBuilder().create();
    }

    public CreateBotResponse createBot(CreateBotRequest request) throws IOException {
        String json = gson.toJson(request);

        Request httpRequest = new Request.Builder()
            .url(baseUrl + "/api/bot/create")
            .post(RequestBody.create(json, MediaType.parse("application/json")))
            .addHeader("Authorization", "Bearer " + apiKey)
            .addHeader("Content-Type", "application/json")
            .build();

        try (Response response = client.newCall(httpRequest).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Unexpected code " + response);
            }

            String responseBody = response.body().string();
            return gson.fromJson(responseBody, CreateBotResponse.class);
        }
    }

    public Map<String, Object> getBot(String botId) throws IOException {
        Map<String, String> requestBody = new HashMap<>();
        requestBody.put("bot_id", botId);

        String json = gson.toJson(requestBody);

        Request httpRequest = new Request.Builder()
            .url(baseUrl + "/api/bot/get")
            .post(RequestBody.create(json, MediaType.parse("application/json")))
            .addHeader("Authorization", "Bearer " + apiKey)
            .addHeader("Content-Type", "application/json")
            .build();

        try (Response response = client.newCall(httpRequest).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Unexpected code " + response);
            }

            String responseBody = response.body().string();
            return gson.fromJson(responseBody, Map.class);
        }
    }

    public static class CreateBotRequest {
        private String name;
        private String description;

        public CreateBotRequest(String name) {
            this.name = name;
        }

        public void setDescription(String description) {
            this.description = description;
        }
    }

    public static class CreateBotResponse {
        private String botId;

        public String getBotId() {
            return botId;
        }
    }
}
`, project.ProjectName, project.ProjectName, time.Now().Format("2006-01-02")), nil
}

// generateReadme 生成README文档
func (s *sdkGeneratorService) generateReadme(language entity.SDKLanguage, project *entity.Project) string {
	switch language {
	case entity.SDKLanguagePython:
		return fmt.Sprintf(`# %s Python SDK

## Installation

\`\`\`bash
pip install coze-studio-%s
\`\`\`

## Quick Start

\`\`\`python
from coze_studio import CozeClient

client = CozeClient(api_key="your-api-key")

# Create a bot
bot = client.create_bot(name="My Bot", description="A helpful bot")
print(f"Bot created: {bot['bot_id']}")

# Send a message
response = client.send_message(bot_id=bot['bot_id'], message="Hello!")
print(response)
\`\`\`

## API Reference

See the [full API documentation](https://docs.coze.com).
`, project.ProjectName, project.ProjectName)

	case entity.SDKLanguageJavaScript:
		return fmt.Sprintf(`# %s TypeScript SDK

## Installation

\`\`\`bash
npm install @coze-studio/%s
\`\`\`

## Quick Start

\`\`\`typescript
import { CozeClient } from '@coze-studio/%s';

const client = new CozeClient({
  apiKey: 'your-api-key',
});

// Create a bot
const bot = await client.createBot({ name: 'My Bot', description: 'A helpful bot' });
console.log(\`Bot created: \${bot.bot_id}\`);

// Send a message
const response = await client.sendMessage({
  botId: bot.bot_id,
  message: 'Hello!',
});
console.log(response);
\`\`\`

## API Reference

See the [full API documentation](https://docs.coze.com).
`, project.ProjectName, project.ProjectName, project.ProjectName)

	case entity.SDKLanguageGo:
		return fmt.Sprintf(`# %s Go SDK

## Installation

\`\`\`bash
go get github.com/coze-studio/%s
\`\`\`

## Quick Start

\`\`\`go
package main

import (
	"fmt"
	"github.com/coze-studio/%s"
)

func main() {
	client := coze.NewClient("your-api-key")

	// Create a bot
	resp, err := client.CreateBot(&coze.CreateBotRequest{
		Name: "My Bot",
		Description: "A helpful bot",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bot created: %s\\n", resp.BotID)
}
\`\`\`

## API Reference

See the [full API documentation](https://docs.coze.com).
`, project.ProjectName, project.ProjectName, project.ProjectName)

	case entity.SDKLanguageJava:
		return fmt.Sprintf(`# %s Java SDK

## Installation

Add this to your \`pom.xml\`:

\`\`\`xml
<dependency>
    <groupId>com.coze-studio</groupId>
    <artifactId>%s-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
\`\`\`

## Quick Start

\`\`\`java
import com.coze.sdk.CozeClient;
import com.coze.sdk.CreateBotRequest;

public class Main {
    public static void main(String[] args) throws IOException {
        CozeClient client = new CozeClient("your-api-key");

        // Create a bot
        CreateBotRequest request = new CreateBotRequest("My Bot");
        request.setDescription("A helpful bot");

        CreateBotResponse response = client.createBot(request);
        System.out.println("Bot created: " + response.getBotId());
    }
}
\`\`\`

## API Reference

See the [full API documentation](https://docs.coze.com).
`, project.ProjectName, project.ProjectName)

	default:
		return ""
	}
}
