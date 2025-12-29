# API接口文档：问数ChatBI模块

**模块名称**: 问数ChatBI (ChatBI)
**设计文档**: 02-ZKER前台_问数_ChatBI.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 数据源管理API](#2-数据源管理api)
- [3. 智能查询API](#3-智能查询api)
- [4. 流式查询API](#4-流式查询api)
- [5. 查询历史API](#5-查询历史api)
- [6. 可视化API](#6-可视化api)
- [7. 数据模型](#7-数据模型)
- [8. 后端代码示例](#8-后端代码示例)
- [9. 前端代码示例](#9-前端代码示例)
- [10. 错误码定义](#10-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

问数ChatBI是ZKER企业级SaaS平台的智能数据分析前台模块，提供自然语言查询数据库并生成可视化报表的能力：

- ✅ **自然语言查询** - 中文问答式数据查询，无需SQL知识
- ✅ **智能SQL生成** - Text2SQL引擎，支持多数据库方言
- ✅ **多数据源支持** - MySQL、PostgreSQL、ClickHouse、Excel、CSV、API
- ✅ **智能可视化** - 自动推荐最佳图表类型
- ✅ **流式执行** - 实时返回查询进度和结果
- ✅ **查询历史** - 历史记录、收藏、标签管理
- ✅ **数据导出** - 导出为Excel、CSV、图片

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- NLU引擎: 通义千问/GPT-4 + spaCy/HanLP
- Text2SQL: Fine-tuned Text2SQL模型
- 数据库: MySQL 8.4.5、PostgreSQL、ClickHouse
- 缓存: Redis 8.0
- ORM: GORM + gorm.io/gen

**前端技术栈**:
- 框架: React 18 + TypeScript
- UI库: Semi Design
- 图表库: ECharts 5.x、D3.js
- Markdown渲染: react-markdown

### 1.3 性能指标

| 指标 | 目标值 | 说明 |
|------|--------|------|
| NLU准确率 | > 90% | 意图识别和实体提取准确率 |
| SQL生成准确率 | > 85% | 生成的SQL可执行且语义正确 |
| 查询响应P95 | < 3s | 95%的查询在3秒内返回 |
| 并发能力 | > 100 QPS | 单实例并发查询能力 |
| 流式首字节 | < 500ms | 流式查询首字节返回时间 |

---

## 2. 数据源管理API

### 2.1 创建数据源

**接口地址**: `POST /api/v1/chatbi/datasources`

**功能说明**: 创建新的数据源连接

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "name": "销售数据库",
  "type": "mysql",
  "config": {
    "host": "192.168.1.100",
    "port": 3306,
    "database": "sales_db",
    "username": "readonly_user",
    "password": "encrypted_password",
    "ssl_mode": "disabled",
    "max_open_conns": 10,
    "max_idle_conns": 5,
    "conn_max_lifetime": 3600
  },
  "description": "2025年销售数据主库"
}
```

**数据源类型枚举**:
- `mysql`: MySQL数据库
- `postgresql`: PostgreSQL数据库
- `clickhouse`: ClickHouse分析型数据库
- `excel`: Excel文件
- `csv`: CSV文件
- `api`: HTTP API接口

**响应示例**:
```json
{
  "code": 0,
  "message": "数据源创建成功",
  "data": {
    "id": "ds-001",
    "tenant_id": "tenant-123",
    "name": "销售数据库",
    "type": "mysql",
    "config": {
      "host": "192.168.1.100",
      "port": 3306,
      "database": "sales_db"
    },
    "status": "active",
    "created_at": "2025-01-03T10:00:00Z",
    "test_result": {
      "success": true,
      "latency_ms": 45,
      "database_version": "8.4.5"
    }
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "trace-abc-123"
}
```

### 2.2 测试数据源连接

**接口地址**: `POST /api/v1/chatbi/datasources/{datasource_id}/test`

**功能说明**: 测试数据源连接是否可用

**响应示例**:
```json
{
  "code": 0,
  "message": "连接测试成功",
  "data": {
    "success": true,
    "latency_ms": 45,
    "database_version": "MySQL 8.4.5",
    "tables_count": 127,
    "sample_tables": [
      {"name": "orders", "rows": 1500000},
      {"name": "customers", "rows": 50000},
      {"name": "products", "rows": 10000}
    ]
  }
}
```

### 2.3 查询数据源列表

**接口地址**: `GET /api/v1/chatbi/datasources`

**功能说明**: 获取当前租户的数据源列表

**查询参数**:
- page: 页码，默认1
- page_size: 每页数量，默认20
- type: 数据源类型过滤
- status: 状态过滤 (active/inactive)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "id": "ds-001",
        "name": "销售数据库",
        "type": "mysql",
        "status": "active",
        "tables_count": 127,
        "last_tested_at": "2025-01-03T09:30:00Z",
        "created_at": "2025-01-03T10:00:00Z"
      }
    ]
  }
}
```

### 2.4 获取数据源详情

**接口地址**: `GET /api/v1/chatbi/datasources/{datasource_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "ds-001",
    "name": "销售数据库",
    "type": "mysql",
    "config": {
      "host": "192.168.1.100",
      "port": 3306,
      "database": "sales_db"
    },
    "status": "active",
    "tables": [
      {
        "name": "orders",
        "columns": [
          {"name": "id", "type": "bigint", "primary_key": true},
          {"name": "customer_id", "type": "bigint"},
          {"name": "order_date", "type": "datetime"},
          {"name": "amount", "type": "decimal(10,2)"}
        ],
        "rows": 1500000
      }
    ],
    "created_at": "2025-01-03T10:00:00Z"
  }
}
```

### 2.5 更新数据源

**接口地址**: `PUT /api/v1/chatbi/datasources/{datasource_id}`

**请求参数**: 同创建数据源

### 2.6 删除数据源

**接口地址**: `DELETE /api/v1/chatbi/datasources/{datasource_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "数据源删除成功"
}
```

---

## 3. 智能查询API

### 3.1 执行自然语言查询

**接口地址**: `POST /api/v1/chatbi/query`

**功能说明**: 将自然语言转换为SQL并执行

**请求参数**:
```json
{
  "query": "2025年1月各地区的销售额统计",
  "datasource_id": "ds-001",
  "chart_type": "auto",
  "options": {
    "include_sql": true,
    "include_explanation": true,
    "max_rows": 1000
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "query_id": "q-001",
    "natural_query": "2025年1月各地区的销售额统计",
    "nlu_result": {
      "intent": "aggregate_query",
      "entities": {
        "time_range": {
          "start": "2025-01-01",
          "end": "2025-01-31"
        },
        "dimensions": ["region"],
        "metrics": ["sales_amount"],
        "aggregation": "sum"
      }
    },
    "generated_sql": "SELECT region, SUM(amount) as total_sales FROM orders WHERE order_date BETWEEN '2025-01-01' AND '2025-01-31' GROUP BY region ORDER BY total_sales DESC",
    "sql_explanation": "该SQL查询2025年1月各地区的销售总额，按销售额降序排列",
    "execution_result": {
      "rows": [
        {"region": "华东", "total_sales": 1500000.00},
        {"region": "华南", "total_sales": 1200000.00},
        {"region": "华北", "total_sales": 980000.00}
      ],
      "rows_count": 3,
      "execution_time_ms": 245
    },
    "visualization": {
      "recommended_chart": "bar",
      "chart_config": {
        "type": "bar",
        "x_axis": "region",
        "y_axis": "total_sales",
        "title": "2025年1月各地区销售额统计",
        "series": [
          {
            "name": "销售额",
            "data": [1500000, 1200000, 980000],
            "itemStyle": {"color": "#5470C6"}
          }
        ]
      }
    },
    "statistics": {
      "confidence_score": 0.92,
      "sql_accuracy_score": 0.95,
      "total_latency_ms": 1234
    }
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "trace-xyz-789"
}
```

### 3.2 SQL查询模式

**接口地址**: `POST /api/v1/chatbi/query/sql`

**功能说明**: 直接执行SQL查询

**请求参数**:
```json
{
  "sql": "SELECT region, SUM(amount) as total_sales FROM orders WHERE order_date >= '2025-01-01' GROUP BY region",
  "datasource_id": "ds-001",
  "chart_type": "auto"
}
```

### 3.3 获取查询建议

**接口地址**: `GET /api/v1/chatbi/query/suggestions`

**功能说明**: 基于数据源结构获取查询建议

**查询参数**:
- datasource_id: 数据源ID
- query: 用户输入的部分查询文本

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "suggestions": [
      "2025年1月各地区的销售额统计",
      "2025年1月各产品类别的销售数量",
      "2025年1月Top10客户的销售额排名"
    ],
    "example_queries": [
      "查询2025年1月的总销售额",
      "按产品类别统计销售额",
      "销售额Top10的客户有哪些"
    ]
  }
}
```

---

## 4. 流式查询API

### 4.1 WebSocket连接建立

**WebSocket URL**:
```
wss://api.example.com/v1/chatbi/query/stream?token={jwt_token}
```

### 4.2 客户端消息格式

#### 发起流式查询
```json
{
  "type": "query.execute",
  "data": {
    "query": "2025年1月各地区的销售额统计",
    "datasource_id": "ds-001",
    "query_id": "q-client-123"
  }
}
```

#### 停止查询
```json
{
  "type": "query.cancel",
  "data": {
    "query_id": "q-client-123"
  }
}
```

### 4.3 服务端消息格式

#### NLU解析完成
```json
{
  "type": "nlu.complete",
  "data": {
    "query_id": "q-client-123",
    "temp_id": "q-server-456",
    "nlu_result": {
      "intent": "aggregate_query",
      "entities": {
        "time_range": {"start": "2025-01-01", "end": "2025-01-31"},
        "dimensions": ["region"],
        "metrics": ["sales_amount"]
      }
    },
    "confidence_score": 0.92
  }
}
```

#### SQL生成完成
```json
{
  "type": "sql.generated",
  "data": {
    "query_id": "q-server-456",
    "sql": "SELECT region, SUM(amount) as total_sales FROM orders WHERE order_date BETWEEN '2025-01-01' AND '2025-01-31' GROUP BY region ORDER BY total_sales DESC",
    "explanation": "该SQL查询2025年1月各地区的销售总额"
  }
}
```

#### 开始执行查询
```json
{
  "type": "execution.start",
  "data": {
    "query_id": "q-server-456",
    "estimated_rows": 3
  }
}
```

#### 查询结果片段
```json
{
  "type": "execution.row",
  "data": {
    "query_id": "q-server-456",
    "row": {"region": "华东", "total_sales": 1500000.00},
    "row_index": 0
  }
}
```

#### 查询完成
```json
{
  "type": "execution.complete",
  "data": {
    "query_id": "q-server-456",
    "rows_count": 3,
    "execution_time_ms": 245,
    "visualization": {
      "recommended_chart": "bar",
      "chart_config": {...}
    }
  }
}
```

#### 查询失败
```json
{
  "type": "execution.error",
  "data": {
    "query_id": "q-server-456",
    "error": {
      "code": 20001,
      "message": "SQL执行错误: Unknown column 'invalid_column'",
      "details": "Table 'orders' does not have column 'invalid_column'"
    }
  }
}
```

---

## 5. 查询历史API

### 5.1 获取查询历史

**接口地址**: `GET /api/v1/chatbi/history`

**查询参数**:
- page: 页码
- page_size: 每页数量
- datasource_id: 数据源ID过滤
- is_favorited: 是否收藏过滤

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 50,
    "items": [
      {
        "id": "q-001",
        "natural_query": "2025年1月各地区的销售额统计",
        "generated_sql": "SELECT region, SUM(amount)...",
        "datasource_id": "ds-001",
        "datasource_name": "销售数据库",
        "execution_count": 5,
        "is_favorited": true,
        "tags": ["销售分析", "地区统计"],
        "last_executed_at": "2025-01-03T10:00:00Z",
        "created_at": "2025-01-03T09:00:00Z"
      }
    ]
  }
}
```

### 5.2 收藏查询

**接口地址**: `POST /api/v1/chatbi/history/{query_id}/favorite`

**响应示例**:
```json
{
  "code": 0,
  "message": "收藏成功"
}
```

### 5.3 取消收藏

**接口地址**: `DELETE /api/v1/chatbi/history/{query_id}/favorite`

### 5.4 添加标签

**接口地址**: `POST /api/v1/chatbi/history/{query_id}/tags`

**请求参数**:
```json
{
  "tags": ["销售分析", "月报", "重要"]
}
```

### 5.5 删除查询历史

**接口地址**: `DELETE /api/v1/chatbi/history/{query_id}`

---

## 6. 可视化API

### 6.1 获取图表配置

**接口地址**: `GET /api/v1/chatbi/visualization/chart-config`

**查询参数**:
- chart_type: 图表类型 (bar/line/pie/table/scatter/area/funnel)
- data_preview: 数据预览（JSON数组）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "chart_type": "bar",
    "echarts_option": {
      "title": {"text": "2025年1月各地区销售额统计"},
      "tooltip": {},
      "legend": {"data": ["销售额"]},
      "xAxis": {
        "type": "category",
        "data": ["华东", "华南", "华北"]
      },
      "yAxis": {"type": "value"},
      "series": [
        {
          "name": "销售额",
          "type": "bar",
          "data": [1500000, 1200000, 980000],
          "itemStyle": {"color": "#5470C6"}
        }
      ]
    }
  }
}
```

### 6.2 导出图表

**接口地址**: `POST /api/v1/chatbi/visualization/export`

**请求参数**:
```json
{
  "query_id": "q-001",
  "export_format": "png",
  "chart_config": {...},
  "options": {
    "width": 1920,
    "height": 1080,
    "dpi": 300
  }
}
```

**导出格式**: png、jpg、pdf、svg

**响应示例**:
```json
{
  "code": 0,
  "message": "导出成功",
  "data": {
    "download_url": "https://minio.example.com/exports/chart-q-001-1735900800.png",
    "expires_at": "2025-01-10T10:00:00Z",
    "file_size": 245678
  }
}
```

### 6.3 导出数据

**接口地址**: `POST /api/v1/chatbi/query/export`

**请求参数**:
```json
{
  "query_id": "q-001",
  "export_format": "xlsx",
  "options": {
    "include_sql": true,
    "include_explanation": true
  }
}
```

**导出格式**: xlsx、csv、json

---

## 7. 数据模型

### 7.1 DataSource（数据源）
```typescript
interface DataSource {
  id: string;
  tenant_id: string;
  name: string;
  type: 'mysql' | 'postgresql' | 'clickhouse' | 'excel' | 'csv' | 'api';
  config: DataSourceConfig;
  status: 'active' | 'inactive' | 'error';
  tables_count?: number;
  last_tested_at?: string;
  created_at: string;
  updated_at: string;
}

interface DataSourceConfig {
  // MySQL/PostgreSQL配置
  host?: string;
  port?: number;
  database?: string;
  username?: string;
  password?: string;
  ssl_mode?: string;

  // Excel/CSV配置
  file_url?: string;
  sheet_name?: string;

  // API配置
  api_url?: string;
  api_method?: 'GET' | 'POST';
  api_headers?: Record<string, string>;

  // 连接池配置
  max_open_conns?: number;
  max_idle_conns?: number;
  conn_max_lifetime?: number;
}
```

### 7.2 Query（查询）
```typescript
interface Query {
  id: string;
  tenant_id: string;
  user_id: number;
  datasource_id: string;
  natural_query: string;
  nlu_result: NLUResult;
  generated_sql: string;
  sql_explanation?: string;
  execution_result?: ExecutionResult;
  visualization?: VisualizationConfig;
  statistics: QueryStatistics;
  execution_count: number;
  is_favorited: boolean;
  tags: string[];
  last_executed_at: string;
  created_at: string;
}

interface NLUResult {
  intent: string;
  entities: {
    time_range?: { start: string; end: string };
    dimensions?: string[];
    metrics?: string[];
    aggregation?: string;
    filters?: Record<string, any>;
  };
  confidence_score: number;
}

interface ExecutionResult {
  rows: Record<string, any>[];
  rows_count: number;
  execution_time_ms: number;
}

interface VisualizationConfig {
  recommended_chart: 'bar' | 'line' | 'pie' | 'table' | 'scatter' | 'area' | 'funnel';
  chart_config: any;
}

interface QueryStatistics {
  confidence_score: number;
  sql_accuracy_score: number;
  total_latency_ms: number;
}
```

---

## 8. 后端代码示例

### 8.1 Text2SQL引擎实现

```go
package service

import (
	"context"
	"fmt"
)

// Text2SQLEngine Text2SQL引擎
type Text2SQLEngine struct {
	llmClient     *LLMClient
	dbSchema      *DatabaseSchema
	sqlValidator  *SQLValidator
}

// GenerateSQL 将自然语言转换为SQL
func (e *Text2SQLEngine) GenerateSQL(
	ctx context.Context,
	req *Text2SQLRequest,
) (*Text2SQLResponse, error) {
	// 1. NLU解析：意图识别和实体提取
	nluResult, err := e.parseNL(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("NLU解析失败: %w", err)
	}

	// 2. 获取数据库Schema
	schema, err := e.dbSchema.GetSchema(ctx, req.DataSourceID)
	if err != nil {
		return nil, fmt.Errorf("获取Schema失败: %w", err)
	}

	// 3. SQL生成（使用LLM）
	sql, err := e.llmClient.GenerateSQL(ctx, &SQLGenerationPrompt{
		NaturalQuery: req.Query,
		Intent:       nluResult.Intent,
		Entities:     nluResult.Entities,
		Schema:       schema,
		Dialect:      req.DBType,
	})
	if err != nil {
		return nil, fmt.Errorf("SQL生成失败: %w", err)
	}

	// 4. SQL验证和优化
	validatedSQL, err := e.sqlValidator.ValidateAndOptimize(ctx, &SQLValidateRequest{
		SQL:         sql,
		Schema:      schema,
		MaxRows:     req.MaxRows,
		TimeoutSecs: 30,
	})
	if err != nil {
		return nil, fmt.Errorf("SQL验证失败: %w", err)
	}

	// 5. 生成SQL解释
	explanation := e.generateExplanation(ctx, validatedSQL, nluResult)

	return &Text2SQLResponse{
		SQL:                validatedSQL,
		Explanation:        explanation,
		ConfidenceScore:    nluResult.ConfidenceScore,
		SQLAccuracyScore:   0.95,
	}, nil
}

// parseNL 自然语言解析
func (e *Text2SQLEngine) parseNL(
	ctx context.Context,
	query string,
) (*NLUResult, error) {
	// 1. 意图识别
	intent, err := e.llmClient.ClassifyIntent(ctx, query)
	if err != nil {
		return nil, err
	}

	// 2. 实体提取（时间、维度、指标、过滤条件）
	entities, err := e.llmClient.ExtractEntities(ctx, query, intent)
	if err != nil {
		return nil, err
	}

	// 3. 置信度计算
	confidenceScore := e.calculateConfidence(intent, entities)

	return &NLUResult{
		Intent:          intent,
		Entities:        entities,
		ConfidenceScore: confidenceScore,
	}, nil
}

// calculateConfidence 计算置信度
func (e *Text2SQLEngine) calculateConfidence(
	intent string,
	entities map[string]interface{},
) float64 {
	score := 0.0

	// 意图清晰度
	if intent != "" {
		score += 0.3
	}

	// 实体完整性
	if _, ok := entities["time_range"]; ok {
		score += 0.2
	}
	if _, ok := entities["dimensions"]; ok {
		score += 0.2
	}
	if _, ok := entities["metrics"]; ok {
		score += 0.2
	}
	if _, ok := entities["filters"]; ok {
		score += 0.1
	}

	return score
}
```

### 8.2 查询执行器

```go
package service

import (
	"context"
	"database/sql"
	"time"
)

// QueryExecutor 查询执行器
type QueryExecutor struct {
	dbManager *DatabaseManager
	limiter   *QueryLimiter
	cache     *CacheManager
}

// Execute 执行SQL查询
func (e *QueryExecutor) Execute(
	ctx context.Context,
	req *ExecuteRequest,
) (*ExecuteResponse, error) {
	// 1. 检查缓存（相同查询5分钟内缓存）
	cacheKey := e.generateCacheKey(req.SQL, req.DataSourceID)
	if cached := e.cache.Get(ctx, cacheKey); cached != nil {
		return cached.(*ExecuteResponse), nil
	}

	// 2. 限流检查（单用户最多10个并发查询）
	if err := e.limiter.Acquire(ctx, req.UserID); err != nil {
		return nil, fmt.Errorf("查询并发数超限: %w", err)
	}
	defer e.limiter.Release(req.UserID)

	// 3. 获取数据库连接
	db, err := e.dbManager.GetConnection(ctx, req.DataSourceID)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 4. 执行查询（带超时控制）
	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	startTime := time.Now()
	rows, err := db.QueryContext(queryCtx, req.SQL)
	if err != nil {
		return nil, fmt.Errorf("SQL执行失败: %w", err)
	}
	defer rows.Close()

	// 5. 解析结果集
	result, err := e.parseRows(rows, req.MaxRows)
	if err != nil {
		return nil, fmt.Errorf("解析结果集失败: %w", err)
	}

	executionTime := time.Since(startTime).Milliseconds()

	// 6. 缓存结果
	response := &ExecuteResponse{
		Rows:             result.Rows,
		RowsCount:        result.RowsCount,
		ExecutionTimeMs:  executionTime,
	}
	e.cache.Set(ctx, cacheKey, response, 5*time.Minute)

	return response, nil
}

// parseRows 解析查询结果
func (e *QueryExecutor) parseRows(
	rows *sql.Rows,
	maxRows int,
) (*QueryResult, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := &QueryResult{
		Rows: make([]map[string]interface{}, 0),
	}

	for rows.Next() {
		// 限制返回行数
		if len(result.Rows) >= maxRows {
			break
		}

		// 扫描行数据
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// 构建行Map
		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		result.Rows = append(result.Rows, row)
	}

	result.RowsCount = len(result.Rows)
	return result, nil
}
```

### 8.3 可视化推荐器

```go
package service

import (
	"context"
)

// Visualizer 可视化推荐器
type Visualizer struct {
	ruleEngine *ChartRuleEngine
}

// RecommendChart 推荐最佳图表类型
func (v *Visualizer) RecommendChart(
	ctx context.Context,
	req *RecommendRequest,
) (*RecommendResponse, error) {
	// 1. 分析数据特征
	features := v.analyzeDataFeatures(req.QueryResult)

	// 2. 规则匹配
	recommendations := v.ruleEngine.MatchRules(features)

	// 3. 评分排序
	bestChart := v.scoreAndSort(recommendations, req.UserPreference)

	// 4. 生成图表配置
	chartConfig := v.generateChartConfig(bestChart, req.QueryResult)

	return &RecommendResponse{
		RecommendedChart: bestChart.Type,
		ChartConfig:      chartConfig,
		ConfidenceScore:  bestChart.Score,
	}, nil
}

// analyzeDataFeatures 分析数据特征
func (v *Visualizer) analyzeDataFeatures(
	result *QueryResult,
) *DataFeatures {
	features := &DataFeatures{
		RowCount:    result.RowsCount,
		ColumnCount: len(result.Columns),
	}

	// 检测维度和指标
	for _, col := range result.Columns {
		if isNumericType(col.Type) {
			features.Metrics = append(features.Metrics, col.Name)
		} else {
			features.Dimensions = append(features.Dimensions, col.Name)
		}
	}

	// 检测时间序列
	if features.RowCount > 1 {
		if hasTimeColumn(result.Columns) {
			features.IsTimeSeries = true
		}
	}

	// 检测分类数据
	if len(features.Dimensions) > 0 {
		features.IsCategorical = true
	}

	return features
}

// ChartRuleEngine 图表规则引擎
type ChartRuleEngine struct {
	rules []ChartRule
}

// ChartRule 图表规则
type ChartRule struct {
	ChartType   string
	Condition   func(*DataFeatures) bool
	Score       float64
	Description string
}

// MatchRules 匹配规则
func (e *ChartRuleEngine) MatchRules(
	features *DataFeatures,
) []*ChartRecommendation {
	var recommendations []*ChartRecommendation

	for _, rule := range e.rules {
		if rule.Condition(features) {
			recommendations = append(recommendations, &ChartRecommendation{
				Type:        rule.ChartType,
				Score:       rule.Score,
				Description: rule.Description,
			})
		}
	}

	return recommendations
}

// 初始化规则
func NewChartRuleEngine() *ChartRuleEngine {
	return &ChartRuleEngine{
		rules: []ChartRule{
			{
				ChartType: "bar",
				Condition: func(f *DataFeatures) bool {
					return f.IsCategorical && len(f.Metrics) == 1
				},
				Score:       0.95,
				Description: "分类数据对比，适合柱状图",
			},
			{
				ChartType: "line",
				Condition: func(f *DataFeatures) bool {
					return f.IsTimeSeries && len(f.Metrics) == 1
				},
				Score:       0.95,
				Description: "时间序列趋势，适合折线图",
			},
			{
				ChartType: "pie",
				Condition: func(f *DataFeatures) bool {
					return f.IsCategorical && len(f.Metrics) == 1 && f.RowCount <= 10
				},
				Score:       0.90,
				Description: "少量分类占比，适合饼图",
			},
			{
				ChartType: "table",
				Condition: func(f *DataFeatures) bool {
					return f.ColumnCount > 5 || f.RowCount > 20
				},
				Score:       0.85,
				Description: "多维数据展示，适合表格",
			},
		},
	}
}
```

---

## 9. 前端代码示例

### 9.1 ChatBI查询组件

```typescript
// src/modules/chatbi/components/ChatBIQueryPanel.tsx
import React, { useState, useRef } from 'react';
import { Button, Input, Card, Spin, Toast } from '@douyinfe/semi-ui';
import { IconSend, IconStop } from '@douyinfe/semi-icons';
import { useChatBIQuery } from '../hooks/useChatBIQuery';
import { ResultPanel } from './ResultPanel';
import { StreamingResult } from './StreamingResult';
import styles from './ChatBIQueryPanel.module.scss';

interface ChatBIQueryPanelProps {
  datasourceId: string;
  datasourceName: string;
}

export const ChatBIQueryPanel: React.FC<ChatBIQueryPanelProps> = ({
  datasourceId,
  datasourceName,
}) => {
  const [query, setQuery] = useState('');
  const [isStreaming, setIsStreaming] = useState(false);
  const [queryResult, setQueryResult] = useState<QueryResult | null>(null);

  const wsRef = useRef<WebSocket | null>(null);
  const { executeQuery, cancelQuery } = useChatBIQuery();

  // 发起查询
  const handleExecute = async () => {
    if (!query.trim()) {
      Toast.warning('请输入查询内容');
      return;
    }

    setIsStreaming(true);
    setQueryResult(null);

    try {
      // 建立WebSocket连接
      const ws = new WebSocket(
        `wss://api.example.com/v1/chatbi/query/stream?token=${getJWTToken()}`
      );
      wsRef.current = ws;

      // 监听消息
      ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        handleStreamMessage(message);
      };

      // 发送查询请求
      ws.onopen = () => {
        ws.send(JSON.stringify({
          type: 'query.execute',
          data: {
            query,
            datasource_id: datasourceId,
            query_id: generateUniqueId(),
          },
        }));
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        Toast.error('连接失败，请重试');
        setIsStreaming(false);
      };

      ws.onclose = () => {
        setIsStreaming(false);
      };
    } catch (error) {
      console.error('Query error:', error);
      Toast.error('查询失败');
      setIsStreaming(false);
    }
  };

  // 处理流式消息
  const handleStreamMessage = (message: StreamMessage) => {
    switch (message.type) {
      case 'nlu.complete':
        console.log('NLU解析完成:', message.data.nlu_result);
        break;

      case 'sql.generated':
        console.log('SQL生成:', message.data.sql);
        break;

      case 'execution.row':
        setQueryResult((prev) => ({
          ...prev,
          rows: [...(prev?.rows || []), message.data.row],
        }));
        break;

      case 'execution.complete':
        setQueryResult({
          queryId: message.data.query_id,
          rows: queryResult?.rows || [],
          rowsCount: message.data.rows_count,
          executionTimeMs: message.data.execution_time_ms,
          visualization: message.data.visualization,
        });
        Toast.success('查询完成');
        break;

      case 'execution.error':
        Toast.error(`查询失败: ${message.data.error.message}`);
        setIsStreaming(false);
        break;
    }
  };

  // 停止查询
  const handleCancel = () => {
    if (wsRef.current) {
      wsRef.current.send(JSON.stringify({
        type: 'query.cancel',
        data: { query_id: queryResult?.queryId },
      }));
      wsRef.current.close();
    }
    setIsStreaming(false);
  };

  return (
    <div className={styles.queryPanel}>
      <div className={styles.queryInput}>
        <Input
          placeholder={`输入查询内容，例如："2025年1月各地区的销售额统计"`}
          value={query}
          onChange={setQuery}
          onEnterPress={handleExecute}
          showClear
          multiline
          maxRows={4}
          disabled={isStreaming}
        />
        <div className={styles.actionButtons}>
          {!isStreaming ? (
            <Button
              theme="solid"
              type="primary"
              icon={<IconSend />}
              onClick={handleExecute}
            >
              执行查询
            </Button>
          ) : (
            <Button
              theme="solid"
              type="danger"
              icon={<IconStop />}
              onClick={handleCancel}
            >
              停止查询
            </Button>
          )}
        </div>
      </div>

      {/* 流式结果展示 */}
      {isStreaming && (
        <Card className={styles.streamingCard}>
          <StreamingResult
            datasourceName={datasourceName}
            query={query}
            rows={queryResult?.rows || []}
          />
        </Card>
      )}

      {/* 最终结果展示 */}
      {!isStreaming && queryResult && (
        <ResultPanel
          datasourceName={datasourceName}
          query={query}
          result={queryResult}
        />
      )}
    </div>
  );
};
```

### 9.2 ECharts可视化组件

```typescript
// src/modules/chatbi/components/ChartVisualization.tsx
import React, { useEffect, useRef } from 'react';
import * as echarts from 'echarts';
import { Card } from '@douyinfe/semi-ui';
import styles from './ChartVisualization.module.scss';

interface ChartVisualizationProps {
  chartType: string;
  chartConfig: any;
  title: string;
  width?: number;
  height?: number;
}

export const ChartVisualization: React.FC<ChartVisualizationProps> = ({
  chartType,
  chartConfig,
  title,
  width = 800,
  height = 500,
}) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    // 初始化图表实例
    chartInstance.current = echarts.init(chartRef.current);

    // 设置图表配置
    chartInstance.current.setOption({
      title: {
        text: title,
        left: 'center',
        textStyle: {
          fontSize: 18,
          fontWeight: 'bold',
        },
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow',
        },
      },
      legend: {
        bottom: 10,
      },
      ...chartConfig,
    });

    // 响应式调整
    const handleResize = () => {
      chartInstance.current?.resize();
    };
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      chartInstance.current?.dispose();
    };
  }, [chartType, chartConfig, title]);

  // 导出图表
  const handleExport = () => {
    if (!chartInstance.current) return;

    const url = chartInstance.current.getDataURL({
      type: 'png',
      pixelRatio: 2,
      backgroundColor: '#fff',
    });

    const link = document.createElement('a');
    link.href = url;
    link.download = `${title}.png`;
    link.click();
  };

  return (
    <Card
      className={styles.chartCard}
      style={{ width, height }}
      title={title}
      headerExtraContent={
        <Button size="small" onClick={handleExport}>
          导出图片
        </Button>
      }
    >
      <div ref={chartRef} style={{ width: '100%', height: '100%' }} />
    </Card>
  );
};
```

### 9.3 查询历史组件

```typescript
// src/modules/chatbi/components/QueryHistory.tsx
import React, { useEffect, useState } from 'react';
import { List, Tag, Button, Empty, Toast } from '@douyinfe/semi-ui';
import { IconStar, IconDelete, IconEye } from '@douyinfe/semi-icons';
import { useChatBIHistory } from '../hooks/useChatBIHistory';
import styles from './QueryHistory.module.scss';

export const QueryHistory: React.FC = () => {
  const [history, setHistory] = useState<QueryHistoryItem[]>([]);
  const [filter, setFilter] = useState<'all' | 'favorited'>('all');

  const { getHistory, favoriteQuery, deleteQuery } = useChatBIHistory();

  useEffect(() => {
    loadHistory();
  }, [filter]);

  const loadHistory = async () => {
    const result = await getHistory({
      is_favorited: filter === 'favorited' ? true : undefined,
    });
    setHistory(result.items);
  };

  const handleFavorite = async (queryId: string) => {
    await favoriteQuery(queryId);
    Toast.success('收藏成功');
    loadHistory();
  };

  const handleDelete = async (queryId: string) => {
    await deleteQuery(queryId);
    Toast.success('删除成功');
    loadHistory();
  };

  const handleReExecute = (query: string) => {
    // 跳转到查询页面并填充查询内容
    window.location.href = `/chatbi?query=${encodeURIComponent(query)}`;
  };

  return (
    <div className={styles.historyContainer}>
      <div className={styles.filterTabs}>
        <Button
          theme={filter === 'all' ? 'solid' : 'light'}
          onClick={() => setFilter('all')}
        >
          全部
        </Button>
        <Button
          theme={filter === 'favorited' ? 'solid' : 'light'}
          onClick={() => setFilter('favorited')}
        >
          已收藏
        </Button>
      </div>

      {history.length === 0 ? (
        <Empty
          title="暂无查询历史"
          description="开始查询数据后，历史记录将显示在这里"
        />
      ) : (
        <List
          dataSource={history}
          renderItem={(item) => (
            <List.Item
              className={styles.historyItem}
              main={
                <div className={styles.itemContent}>
                  <div className={styles.queryText}>
                    {item.natural_query}
                  </div>
                  <div className={styles.queryMeta}>
                    <span>数据源: {item.datasource_name}</span>
                    <span>执行次数: {item.execution_count}</span>
                    <span>
                      最后执行: {formatDate(item.last_executed_at)}
                    </span>
                  </div>
                  <div className={styles.tags}>
                    {item.tags.map((tag) => (
                      <Tag key={tag} color="blue">
                        {tag}
                      </Tag>
                    ))}
                  </div>
                </div>
              }
              extra={
                <div className={styles.itemActions}>
                  <Button
                    icon={<IconEye />}
                    onClick={() => handleReExecute(item.natural_query)}
                  >
                    重新执行
                  </Button>
                  <Button
                    icon={<IconStar />}
                    theme={item.is_favorited ? 'solid' : 'light'}
                    onClick={() => handleFavorite(item.id)}
                  >
                    {item.is_favorited ? '已收藏' : '收藏'}
                  </Button>
                  <Button
                    icon={<IconDelete />}
                    type="danger"
                    onClick={() => handleDelete(item.id)}
                  >
                    删除
                  </Button>
                </div>
              }
            />
          )}
        />
      )}
    </div>
  );
};
```

---

## 10. 错误码定义

| 错误码 | HTTP状态码 | 说明 | 处理建议 |
|--------|-----------|------|----------|
| 20001 | 400 | NLU解析失败 | 用户查询表述不清，建议重新表述 |
| 20002 | 400 | SQL生成失败 | 查询过于复杂或缺少必要信息，建议简化查询 |
| 20003 | 400 | SQL执行错误 | 生成的SQL有语法错误，请尝试其他表述方式 |
| 20004 | 404 | 数据源不存在 | 检查数据源ID是否正确 |
| 20005 | 503 | 数据源连接失败 | 检查数据源配置是否正确，网络是否可达 |
| 20006 | 400 | SQL验证失败 | 生成的SQL不符合安全规范，请联系管理员 |
| 20007 | 429 | 查询并发数超限 | 等待当前查询完成后重试 |
| 20008 | 400 | 查询结果超出行数限制 | 增加过滤条件或调整时间范围 |
| 20009 | 500 | 可视化生成失败 | 数据格式不支持该图表类型，尝试其他图表 |
| 20010 | 400 | 导出失败 | 文件生成或上传失败，请重试 |
| 20011 | 401 | 未授权访问 | 请登录后重试 |
| 20012 | 403 | 权限不足 | 无权限访问该数据源 |

**文档结束**
