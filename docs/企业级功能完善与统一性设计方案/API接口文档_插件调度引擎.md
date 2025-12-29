# API接口文档：插件调度引擎模块

**模块名称**: 插件调度引擎 (PluginScheduler)
**设计文档**: 19-五大AI引擎核心_插件调度引擎.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 插件调度API](#2-插件调度api)
- [3. 插件管理API](#3-插件管理api)
- [4. 执行监控API](#4-执行监控api)
- [5. 数据模型](#5-数据模型)
- [6. 后端实现示例](#6-后端实现示例)
- [7. 前端调用示例](#7-前端调用示例)
- [8. 错误码定义](#8-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

插件调度引擎负责智能调用和管理插件/技能，通过意图识别、参数提取和插件执行的流程，实现AI能力的扩展：

- ✅ **智能调度** - 根据意图自动选择合适的插件
- ✅ **并行执行** - 支持多插件并行调用
- ✅ **参数提取** - 自动从用户输入中提取插件参数
- ✅ **结果格式化** - 整合和格式化插件输出
- ✅ **插件管理** - 动态注册、健康检查、限流
- ✅ **执行监控** - 日志记录、性能监控

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5 (复用zker plugins表 + 新增plugin_executions表)
- 意图识别: 内置IntentEngine
- HTTP客户端: 标准库 + 超时控制

### 1.3 调度流程

```
用户输入 → 意图识别 → 插件匹配 → 参数提取 → 插件执行 → 结果格式化 → 返回用户
          (NLU)      (规则/ML)      (LLM)      (HTTP API)     (模板)
```

---

## 2. 插件调度API

### 2.1 调度插件执行

**接口地址**: `POST /api/v1/plugin/schedule`

**请求参数**:
```json
{
  "query": "帮我查询一下北京的天气",
  "conversation_id": "conv-001",
  "message_id": "msg-001",
  "context": {
    "user_location": "上海"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "插件调度成功",
  "data": {
    "execution_id": "exec-001",
    "plugin_id": "weather-plugin",
    "plugin_name": "天气查询插件",
    "intent": "query_weather",
    "parameters": {
      "location": "北京"
    },
    "result": {
      "temperature": "15°C",
      "condition": "晴",
      "humidity": "45%"
    },
    "status": "success",
    "duration_ms": 1250,
    "executed_at": "2025-01-03T11:00:00Z"
  },
  "timestamp": "2025-01-03T11:00:00Z",
  "trace_id": "trace-001"
}
```

### 2.2 并行调度多个插件

**接口地址**: `POST /api/v1/plugin/schedule/batch`

**请求参数**:
```json
{
  "query": "帮我查询北京天气并推荐适合的穿搭",
  "conversation_id": "conv-001",
  "message_id": "msg-001",
  "execution_mode": "parallel"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "批量调度成功",
  "data": {
    "executions": [
      {
        "execution_id": "exec-001",
        "plugin_id": "weather-plugin",
        "plugin_name": "天气查询插件",
        "result": {
          "temperature": "15°C",
          "condition": "晴"
        },
        "status": "success"
      },
      {
        "execution_id": "exec-002",
        "plugin_id": "outfit-plugin",
        "plugin_name": "穿搭推荐插件",
        "result": {
          "recommendation": "适合穿轻薄外套"
        },
        "status": "success"
      }
    ],
    "total": 2,
    "aggregated_result": "北京今天15°C，晴天，湿度45%。适合穿轻薄外套，搭配牛仔裤。"
  }
}
```

### 2.3 手动指定插件执行

**接口地址**: `POST /api/v1/plugin/execute`

**请求参数**:
```json
{
  "plugin_id": "weather-plugin",
  "parameters": {
    "location": "北京",
    "unit": "celsius"
  },
  "conversation_id": "conv-001",
  "message_id": "msg-001"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "插件执行成功",
  "data": {
    "execution_id": "exec-001",
    "result": {
      "temperature": "15°C",
      "condition": "晴"
    }
  }
}
```

---

## 3. 插件管理API

### 3.1 注册新插件

**接口地址**: `POST /api/v1/plugin/register`

**请求参数**:
```json
{
  "name": "天气查询插件",
  "description": "查询各地天气情况",
  "category": "utility",
  "intents": ["query_weather", "get_temperature"],
  "config": {
    "api_url": "https://api.weather.com/v1",
    "api_key": "sk-xxx",
    "action": "get_weather",
    "timeout": 30,
    "rate_limit": 100
  },
  "parameter_schema": {
    "type": "object",
    "properties": {
      "location": {
        "type": "string",
        "description": "城市名称",
        "required": true
      },
      "unit": {
        "type": "string",
        "enum": ["celsius", "fahrenheit"],
        "default": "celsius"
      }
    }
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "插件注册成功",
  "data": {
    "plugin_id": "plugin-001",
    "name": "天气查询插件",
    "status": "active",
    "created_at": "2025-01-03T11:00:00Z"
  }
}
```

### 3.2 更新插件配置

**接口地址**: `PUT /api/v1/plugin/{plugin_id}`

**请求参数**:
```json
{
  "config": {
    "api_url": "https://api.weather.com/v2",
    "timeout": 45
  }
}
```

### 3.3 获取插件详情

**接口地址**: `GET /api/v1/plugin/{plugin_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "plugin_id": "plugin-001",
    "name": "天气查询插件",
    "description": "查询各地天气情况",
    "category": "utility",
    "intents": ["query_weather"],
    "status": "active",
    "config": {
      "api_url": "https://api.weather.com/v1",
      "timeout": 30
    },
    "parameter_schema": {},
    "created_at": "2025-01-03T11:00:00Z",
    "updated_at": "2025-01-03T11:00:00Z"
  }
}
```

### 3.4 列出所有插件

**接口地址**: `GET /api/v1/plugin/list`

**查询参数**:
- category: 类别过滤
- status: 状态过滤
- page: 页码
- page_size: 每页数量

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "plugins": [
      {
        "plugin_id": "plugin-001",
        "name": "天气查询插件",
        "category": "utility",
        "status": "active"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

### 3.5 启用/禁用插件

**接口地址**: `PUT /api/v1/plugin/{plugin_id}/status`

**请求参数**:
```json
{
  "status": "disabled",
  "reason": "API接口维护中"
}
```

### 3.6 删除插件

**接口地址**: `DELETE /api/v1/plugin/{plugin_id}`

### 3.7 插件健康检查

**接口地址**: `POST /api/v1/plugin/{plugin_id}/health-check`

**响应示例**:
```json
{
  "code": 0,
  "message": "健康检查完成",
  "data": {
    "plugin_id": "plugin-001",
    "status": "healthy",
    "response_time_ms": 120,
    "last_check": "2025-01-03T11:00:00Z",
    "details": {
      "api_reachable": true,
      "auth_valid": true,
      "error_rate": 0.01
    }
  }
}
```

---

## 4. 执行监控API

### 4.1 获取执行日志

**接口地址**: `GET /api/v1/plugin/executions`

**查询参数**:
- plugin_id: 插件ID过滤
- status: 状态过滤
- start_time: 开始时间
- end_time: 结束时间
- page: 页码
- page_size: 每页数量

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "executions": [
      {
        "execution_id": "exec-001",
        "plugin_id": "plugin-001",
        "plugin_name": "天气查询插件",
        "status": "success",
        "duration_ms": 1250,
        "start_time": "2025-01-03T11:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 4.2 获取执行详情

**接口地址**: `GET /api/v1/plugin/executions/{execution_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "execution_id": "exec-001",
    "tenant_id": "tenant-001",
    "user_id": 1001,
    "conversation_id": "conv-001",
    "message_id": "msg-001",
    "plugin_id": "plugin-001",
    "plugin_name": "天气查询插件",
    "input_params": {
      "location": "北京"
    },
    "output_data": {
      "temperature": "15°C",
      "condition": "晴"
    },
    "status": "success",
    "error_message": null,
    "start_time": "2025-01-03T11:00:00Z",
    "end_time": "2025-01-03T11:00:01.25Z",
    "duration_ms": 1250
  }
}
```

### 4.3 获取插件统计

**接口地址**: `GET /api/v1/plugin/{plugin_id}/stats`

**查询参数**:
- start_time: 开始时间
- end_time: 结束时间

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "plugin_id": "plugin-001",
    "plugin_name": "天气查询插件",
    "total_executions": 1500,
    "successful_executions": 1425,
    "failed_executions": 75,
    "success_rate": 0.95,
    "avg_duration_ms": 1200,
    "p50_duration_ms": 1100,
    "p95_duration_ms": 1800,
    "p99_duration_ms": 2500,
    "last_execution": "2025-01-03T11:00:00Z"
  }
}
```

### 4.4 获取执行趋势

**接口地址**: `GET /api/v1/plugin/executions/trends`

**查询参数**:
- plugin_id: 插件ID
- granularity: 时间粒度 (hour/day/week)
- start_time: 开始时间
- end_time: 结束时间

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "trends": [
      {
        "timestamp": "2025-01-03T10:00:00Z",
        "total": 120,
        "success": 115,
        "failed": 5,
        "avg_duration_ms": 1180
      }
    ]
  }
}
```

---

## 5. 数据模型

### 5.1 Plugin

```typescript
interface Plugin {
  plugin_id: string;
  name: string;
  description?: string;
  category: string;
  intents: string[];
  status: 'active' | 'disabled' | 'deprecated';
  config: PluginConfig;
  parameter_schema: JSONSchema;
  created_at: Date;
  updated_at: Date;
}

interface PluginConfig {
  api_url: string;
  api_key?: string;
  action?: string;
  timeout?: number;
  rate_limit?: number;
  [key: string]: any;
}
```

### 5.2 PluginExecution

```typescript
interface PluginExecution {
  execution_id: string;
  tenant_id: string;
  user_id: number;
  conversation_id?: string;
  message_id?: string;

  plugin_id: string;
  plugin_name: string;

  input_params?: Record<string, any>;
  output_data?: Record<string, any>;
  status: 'pending' | 'running' | 'success' | 'failed' | 'timeout';
  error_message?: string;

  start_time: Date;
  end_time?: Date;
  duration_ms?: number;
}
```

### 5.3 ScheduleRequest

```typescript
interface ScheduleRequest {
  query: string;
  conversation_id: string;
  message_id: string;
  context?: Record<string, any>;
}
```

### 5.4 ScheduleResult

```typescript
interface ScheduleResult {
  execution_id: string;
  plugin_id: string;
  plugin_name: string;
  intent: string;
  parameters: Record<string, any>;
  result: Record<string, any>;
  status: string;
  duration_ms: number;
  executed_at: Date;
}
```

---

## 6. 后端实现示例

### 6.1 PluginScheduler 核心实现

```go
package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type PluginScheduler struct {
	pluginRepo    PluginRepository
	executionRepo ExecutionRepository
	intentEngine  *IntentEngine
	httpClient    *http.Client
}

// SchedulePlugin 调度插件
func (s *PluginScheduler) SchedulePlugin(
	ctx context.Context,
	req *ScheduleRequest,
) (*ScheduleResult, error) {
	// 1. 意图识别
	intent, err := s.intentEngine.RecognizeIntent(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("意图识别失败: %w", err)
	}

	// 2. 匹配插件
	plugin, err := s.pluginRepo.GetByIntent(ctx, intent.Intent)
	if err != nil {
		return nil, fmt.Errorf("未找到匹配的插件: %s", intent.Intent)
	}

	if plugin.Status != "active" {
		return nil, fmt.Errorf("插件未启用: %s", plugin.ID)
	}

	// 3. 参数提取
	params, err := s.extractParameters(ctx, plugin, intent.Entities)
	if err != nil {
		return nil, fmt.Errorf("参数提取失败: %w", err)
	}

	// 4. 创建执行记录
	execution := &PluginExecution{
		ID:             generateExecutionID(),
		TenantID:       getTenantID(ctx),
		UserID:         getUserID(ctx),
		ConversationID: req.ConversationID,
		MessageID:      req.MessageID,
		PluginID:       plugin.ID,
		PluginName:     plugin.Name,
		InputParams:    params,
		Status:         "pending",
		StartTime:      time.Now(),
	}

	err = s.executionRepo.Create(ctx, execution)
	if err != nil {
		return nil, err
	}

	// 5. 执行插件
	result, err := s.executePlugin(ctx, plugin, params)
	endTime := time.Now()
	durationMs := int(endTime.Sub(execution.StartTime).Milliseconds())

	execution.EndTime = &endTime
	execution.DurationMs = &durationMs

	if err != nil {
		execution.Status = "failed"
		execution.ErrorMessage = err.Error()
		s.executionRepo.Update(ctx, execution)
		return nil, err
	}

	// 6. 保存成功结果
	execution.Status = "success"
	execution.OutputData = result
	s.executionRepo.Update(ctx, execution)

	return &ScheduleResult{
		ExecutionID: execution.ID,
		PluginID:    plugin.ID,
		PluginName:  plugin.Name,
		Intent:      intent.Intent,
		Parameters:  params,
		Result:      result,
		Status:      "success",
		DurationMs:  durationMs,
		ExecutedAt:  endTime,
	}, nil
}

// ScheduleBatchPlugins 并行调度多个插件
func (s *PluginScheduler) ScheduleBatchPlugins(
	ctx context.Context,
	req *ScheduleRequest,
) (*BatchScheduleResult, error) {
	// 1. 意图识别
	intent, err := s.intentEngine.RecognizeIntent(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	// 2. 查找所有相关插件
	plugins, err := s.pluginRepo.GetByIntentMatch(ctx, intent.Intent)
	if err != nil {
		return nil, err
	}

	// 过滤活跃插件
	activePlugins := make([]*Plugin, 0)
	for _, plugin := range plugins {
		if plugin.Status == "active" {
			activePlugins = append(activePlugins, plugin)
		}
	}

	// 3. 并行执行
	type resultPair struct {
		Index   int
		Result  *ScheduleResult
		Errored error
	}

	resultChan := make(chan resultPair, len(activePlugins))

	for i, plugin := range activePlugins {
		go func(index int, p *Plugin) {
			params, _ := s.extractParameters(ctx, p, intent.Entities)
			execResult, err := s.executePlugin(ctx, p, params)

			resultChan <- resultPair{
				Index:   index,
				Result:  execResult,
				Errored: err,
			}
		}(i, plugin)
	}

	// 4. 收集结果
	results := make([]*ScheduleResult, len(activePlugins))
	errors := make([]error, 0)

	for range activePlugins {
		pair := <-resultChan
		if pair.Errored != nil {
			errors = append(errors, pair.Errored)
		} else {
			results[pair.Index] = pair.Result
		}
	}

	// 5. 聚合结果
	aggregated := s.aggregateResults(results)

	return &BatchScheduleResult{
		Executions:       results,
		Total:            len(results),
		AggregatedResult: aggregated,
	}, nil
}

// executePlugin 执行插件
func (s *PluginScheduler) executePlugin(
	ctx context.Context,
	plugin *Plugin,
	params map[string]interface{},
) (map[string]interface{}, error) {
	// 1. 构建请求
	apiURL := plugin.Config["api_url"].(string)
	reqBody := map[string]interface{}{
		"action": plugin.Config["action"],
		"params": params,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 2. 创建HTTP请求
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		apiURL,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if apiKey, ok := plugin.Config["api_key"].(string); ok {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	// 3. 设置超时并执行
	timeout := 30 * time.Second
	if t, ok := plugin.Config["timeout"].(int); ok {
		timeout = time.Duration(t) * time.Second
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4. 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("插件返回错误: %v", result)
	}

	return result, nil
}

// extractParameters 提取参数
func (s *PluginScheduler) extractParameters(
	ctx context.Context,
	plugin *Plugin,
	entities map[string]string,
) (map[string]interface{}, error) {
	// 根据插件的parameter_schema和entities提取参数
	params := make(map[string]interface{})

	// 简化实现：直接使用entities
	for key, value := range entities {
		params[key] = value
	}

	return params, nil
}

// aggregateResults 聚合结果
func (s *PluginScheduler) aggregateResults(
	results []*ScheduleResult,
) string {
	// 简化实现：拼接所有结果
	var aggregated string
	for _, result := range results {
		if result != nil {
			aggregated += fmt.Sprintf("%v ", result.Result)
		}
	}
	return aggregated
}
```

### 6.2 HTTP Handler

```go
// SchedulePluginHandler 调度插件
func SchedulePluginHandler(ctx context.Context, c *app.RequestContext) {
	var req ScheduleRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	result, err := pluginScheduler.SchedulePlugin(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "插件调度失败",
			Detail:  err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, SuccessResponse{
		Code:    0,
		Message: "插件调度成功",
		Data:    result,
	})
}

// RegisterPluginHandler 注册插件
func RegisterPluginHandler(ctx context.Context, c *app.RequestContext) {
	var req RegisterPluginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	plugin := &Plugin{
		ID:          generatePluginID(),
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Intents:     req.Intents,
		Config:      req.Config,
		Status:      "active",
		CreatedAt:   time.Now(),
	}

	err := pluginRepo.Create(ctx, plugin)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "插件注册失败",
		})
		return
	}

	c.JSON(consts.StatusOK, SuccessResponse{
		Code:    0,
		Message: "插件注册成功",
		Data:    plugin,
	})
}
```

---

## 7. 前端调用示例

### 7.1 调度插件

```typescript
import { request } from '@/utils/request';

export async function schedulePlugin(params: {
  query: string;
  conversationId: string;
  messageId: string;
  context?: Record<string, any>;
}) {
  return request.post<ScheduleResult>('/api/v1/plugin/schedule', {
    query: params.query,
    conversation_id: params.conversationId,
    message_id: params.messageId,
    context: params.context,
  });
}

// 使用示例
const result = await schedulePlugin({
  query: '帮我查询北京的天气',
  conversationId: 'conv-001',
  messageId: 'msg-001',
});

console.log('插件执行结果:', result.data.result);
```

### 7.2 注册插件

```typescript
export async function registerPlugin(params: {
  name: string;
  description: string;
  category: string;
  intents: string[];
  config: Record<string, any>;
  parameterSchema: JSONSchema;
}) {
  return request.post<Plugin>('/api/v1/plugin/register', {
    name: params.name,
    description: params.description,
    category: params.category,
    intents: params.intents,
    config: params.config,
    parameter_schema: params.parameterSchema,
  });
}

// 使用示例
const plugin = await registerPlugin({
  name: '天气查询插件',
  description: '查询各地天气情况',
  category: 'utility',
  intents: ['query_weather'],
  config: {
    api_url: 'https://api.weather.com/v1',
    api_key: 'sk-xxx',
  },
  parameterSchema: {
    type: 'object',
    properties: {
      location: { type: 'string' },
    },
  },
});
```

### 7.3 获取执行日志

```typescript
export async function getExecutionLogs(params: {
  pluginId?: string;
  status?: string;
  startTime?: Date;
  endTime?: Date;
  page?: number;
  pageSize?: number;
}) {
  return request.get<{
    executions: PluginExecution[];
    total: number;
  }>('/api/v1/plugin/executions', {
    params: {
      plugin_id: params.pluginId,
      status: params.status,
      start_time: params.startTime?.toISOString(),
      end_time: params.endTime?.toISOString(),
      page: params.page || 1,
      page_size: params.pageSize || 20,
    },
  });
}

// 使用示例
const logs = await getExecutionLogs({
  pluginId: 'plugin-001',
  status: 'success',
  page: 1,
});

console.log('执行日志:', logs.data.executions);
```

### 7.4 React Hook 示例

```typescript
import { useState } from 'react';
import { schedulePlugin } from '@/api/plugin-scheduler';

export function usePluginScheduler() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const schedule = async (
    query: string,
    conversationId: string,
    messageId: string,
  ) => {
    setLoading(true);
    setError(null);

    try {
      const result = await schedulePlugin({
        query,
        conversationId,
        messageId,
      });
      return result.data;
    } catch (err) {
      setError(err as Error);
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { loading, error, schedule };
}

// 使用示例
function ChatComponent() {
  const { loading, schedule } = usePluginScheduler();

  const handleSendMessage = async (message: string) => {
    const result = await schedule(
      message,
      'conv-001',
      generateMessageId(),
    );

    if (result) {
      console.log('插件结果:', result.result);
    }
  };

  return (
    <div>
      <button onClick={() => handleSendMessage('查询北京天气')}>
        查询天气
      </button>
      {loading && <div>执行中...</div>}
    </div>
  );
}
```

---

## 8. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 10001 | 404 | 插件不存在 |
| 10002 | 400 | 插件未启用 |
| 10003 | 400 | 插件配置无效 |
| 10004 | 400 | 参数提取失败 |
| 10005 | 400 | 参数验证失败 |
| 10006 | 500 | 插件执行超时 |
| 10007 | 500 | 插件API调用失败 |
| 10008 | 429 | 插件调用限流 |
| 10009 | 500 | 意图识别失败 |
| 10010 | 404 | 未找到匹配的插件 |
| 10011 | 500 | 结果聚合失败 |
| 10101 | 400 | 插件注册失败 |
| 10102 | 403 | 插件名称重复 |
| 10201 | 404 | 执行记录不存在 |

---

**文档结束**
