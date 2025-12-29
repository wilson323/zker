# ZKER 企业级AI智能体工作台 - 问数(ChatBI)模块详细设计

> **文档编号**: ZKER-DD-2025-002
> **模块名称**: 问数(ChatBI) - 自然语言数据分析
> **项目名称**: ZKER Enterprise Edition
> **文档版本**: v1.0
> **编写日期**: 2025-12-29
> **密级**: 内部公开

---

## 文档修订历史

| 版本 | 日期 | 作者 | 修订说明 |
|------|------|------|----------|
| v1.0 | 2025-12-29 | AI助手 | 初始版本 |

---

## 目录

1. [模块概述](#1-模块概述)
2. [功能需求分析](#2-功能需求分析)
3. [系统架构设计](#3-系统架构设计)
4. [数据库设计](#4-数据库设计)
5. [接口设计](#5-接口设计)
6. [前端设计](#6-前端设计)
7. [后端设计](#7-后端设计)
8. [性能优化](#8-性能优化)
9. [测试用例](#9-测试用例)

---

## 1. 模块概述

### 1.1 功能定义

**ZKER 问数(ChatBI)** 是企业级自然语言数据分析模块，让业务人员通过对话方式查询和分析企业数据，无需编写SQL或使用复杂的BI工具。

### 1.2 核心价值

**对企业**：
- ✅ 降低数据分析门槛，业务人员自助分析
- ✅ 提升决策效率，实时获取数据洞察
- ✅ 减少IT部门压力，减少报表开发需求

**对员工**：
- ✅ 零学习成本，像聊天一样查询数据
- ✅ 快速获取洞察，无需等待IT部门
- ✅ 灵活分析，支持多维度探索

### 1.3 业务场景

```
场景1: 销售数据查询
用户: "上个月华东区的销售额是多少?"
ZKER: [生成SQL] → [执行查询] → [生成图表]
    "上个月华东区销售额为 ¥1,234,567,同比增长 15%"

场景2: 趋势分析
用户: "按月份统计今年的销售趋势"
ZKER: [生成SQL] → [查询数据] → [生成折线图]
    [展示折线图 + 趋势分析]

场景3: 对比分析
用户: "对比各地区的销售业绩"
ZKER: [生成SQL] → [查询数据] → [生成柱状图]
    [展示柱状图 + 排名]

场景4: 数据导出
用户: "把这个报表导出Excel"
ZKER: [生成Excel文件] → [自动下载]
```

### 1.4 技术目标

| 指标 | 目标值 | 测量方法 |
|------|--------|----------|
| **自然语言理解准确率** | > 90% | 测试集评估 |
| **SQL生成正确率** | > 85% | 测试集评估 |
| **查询响应时间** | P95 < 3s | APM监控 |
| **图表生成时间** | < 1s | 性能测试 |
| **并发查询能力** | > 100 QPS | 压力测试 |

### 1.5 技术栈

**前端**:
- React 18 + TypeScript
- ECharts / D3.js (图表库)
- Monaco Editor (SQL编辑器)

**后端**:
- Go 1.23 + Hertz框架
- NLG (Natural Language Generation) 模型
- SQL解析器 (SQLGlot)
- 数据库驱动 (MySQL / PostgreSQL / ClickHouse)

**AI模型**:
- 通义千问 (Qwen) / GPT-4 (文本转SQL)
- Text-to-SQL (Text2SQL) 微调模型

---

## 2. 功能需求分析

### 2.1 功能全景图

```
ZKER 问数(ChatBI) 模块
│
├── 2.1 自然语言查询
│   ├── 文本转SQL (Text2SQL)
│   ├── 意图识别
│   ├── 实体提取
│   └── SQL生成与优化
│
├── 2.2 数据源管理
│   ├── 数据源连接
│   │   ├── MySQL
│   │   ├── PostgreSQL
│   │   ├── ClickHouse
│   │   └── Excel/CSV
│   ├── 数据源权限控制
│   └── 数据源元数据管理
│
├── 2.3 查询执行引擎
│   ├── SQL执行
│   ├── 结果缓存
│   ├── 流式返回
│   └── 错误处理
│
├── 2.4 可视化呈现
│   ├── 图表自动生成
│   │   ├── 柱状图
│   │   ├── 折线图
│   │   ├── 饼图
│   │   ├── 散点图
│   │   └── 表格
│   ├── 智能图表推荐
│   └── 交互式探索
│
├── 2.5 数据导出
│   ├── Excel导出
│   ├── CSV导出
│   ├── PDF报告
│   └── 图片导出
│
└── 2.6 查询历史与收藏
    ├── 历史查询记录
    ├── 常用查询收藏
    └── 查询分享
```

### 2.2 核心功能规格

#### 2.2.1 自然语言查询 (P0)

**功能描述**: 用户使用自然语言描述数据需求，系统自动转换为SQL并执行。

**输入**:
```typescript
interface QueryRequest {
  query: string;           // 自然语言查询
  dataSourceId?: string;   // 数据源ID (可选,默认使用主数据源)
  userId: string;          // 用户ID
  tenantId: string;        // 租户ID
}
```

**处理流程**:
```go
type ChatBIService struct {
    nlgEngine       *NLGEngine          // 自然语言生成引擎
    sqlParser       *SQLParser          // SQL解析器
    queryExecutor   *QueryExecutor      // 查询执行器
    visualizer      *Visualizer         // 可视化引擎
    permissionCheck *PermissionChecker  // 权限检查
}

func (s *ChatBIService) ProcessQuery(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
    // 1. 权限检查 (用户是否有权限访问该数据源)
    if err := s.permissionCheck.CheckDataSourceAccess(ctx, req.UserId, req.DataSourceId); err != nil {
        return nil, err
    }

    // 2. 意图识别 (识别查询类型: 聚合、对比、趋势等)
    intent := s.nlgEngine.RecognizeIntent(ctx, req.Query)

    // 3. 实体提取 (提取关键实体: 时间、地区、指标等)
    entities := s.nlgEngine.ExtractEntities(ctx, req.Query, intent)

    // 4. SQL生成 (根据意图和实体生成SQL)
    sql, err := s.nlgEngine.GenerateSQL(ctx, intent, entities, req.DataSourceId)
    if err != nil {
        return nil, err
    }

    // 5. SQL优化 (优化SQL性能)
    sql = s.sqlParser.Optimize(sql)

    // 6. 执行查询
    result, err := s.queryExecutor.Execute(ctx, sql, req.DataSourceId)
    if err != nil {
        return nil, err
    }

    // 7. 数据可视化 (选择最佳图表类型)
    chart := s.visualizer.GenerateChart(ctx, result, intent)

    // 8. 自然语言解释 (生成数据洞察)
    insight := s.nlgEngine.GenerateInsight(ctx, result, intent)

    return &QueryResponse{
        SQL:      sql,
        Result:   result,
        Chart:    chart,
        Insight:   insight,
        Metadata: result.Metadata,
    }, nil
}
```

**输出**:
```typescript
interface QueryResponse {
  sql: string;              // 生成的SQL语句
  result: QueryResult;       // 查询结果
  chart: ChartConfig;        // 图表配置
  insight: string;           // 数据洞察
  metadata: {
    executionTime: number;   // 执行时间(ms)
    rowCount: number;        // 返回行数
    cached: boolean;         // 是否命中缓存
  };
}
```

**验收标准**:
- ✅ 自然语言理解准确率 > 90%
- ✅ SQL生成正确率 > 85%
- ✅ 查询响应时间 P95 < 3s
- ✅ 支持复杂查询 (JOIN、聚合、子查询)

#### 2.2.2 数据源管理 (P0)

**功能描述**: 管理企业数据源连接,支持多种数据源类型。

**数据源类型支持**:
```typescript
enum DataSourceType {
  MySQL = 'mysql',
  PostgreSQL = 'postgresql',
  ClickHouse = 'clickhouse',
  Excel = 'excel',
  CSV = 'csv',
  API = 'api'
}

interface DataSourceConfig {
  id: string;
  name: string;
  type: DataSourceType;
  connection: {
    host?: string;
    port?: number;
    database?: string;
    username?: string;
    password?: string;  // 加密存储
    // Excel/CSV特定配置
    fileUrl?: string;
    sheetName?: string;
  };
  permissions: {
    userIds: string[];      // 有权限的用户ID列表
    roleIds: string[];       // 有权限的角色ID列表
  };
  metadata: {
    tableName: string;
    schema: TableSchema;
    sampleData?: any[];
  };
  status: 'active' | 'inactive' | 'error';
  createdAt: Date;
  updatedAt: Date;
}
```

**数据源连接测试**:
```go
func (s *DataSourceService) TestConnection(ctx context.Context, config *DataSourceConfig) error {
    switch config.Type {
    case DataSourceTypeMySQL:
        return s.testMySQLConnection(ctx, config.Connection)
    case DataSourceTypePostgreSQL:
        return s.testPostgreSQLConnection(ctx, config.Connection)
    case DataSourceTypeClickHouse:
        return s.testClickHouseConnection(ctx, config.Connection)
    case DataSourceTypeExcel:
        return s.testExcelConnection(ctx, config.Connection)
    default:
        return errors.New("unsupported data source type")
    }
}
```

**验收标准**:
- ✅ 支持5种以上数据源
- ✅ 连接测试准确率 100%
- ✅ 权限控制粒度到表级别
- ✅ 连接信息加密存储

#### 2.2.3 可视化呈现 (P0)

**功能描述**: 自动选择最佳图表类型并生成交互式图表。

**图表类型映射**:
```typescript
type ChartType =
  | 'bar'       // 柱状图 - 对比分析
  | 'line'      // 折线图 - 趋势分析
  | 'pie'       // 饼图 - 占比分析
  | 'table'     // 表格 - 详细数据
  | 'scatter'   // 散点图 - 相关性分析
  | 'area'      // 面积图 - 累积趋势
  | 'funnel'    // 漏斗图 - 转化分析

interface ChartRecommendation {
  chartType: ChartType;
  confidence: number;       // 推荐置信度 (0-1)
  reason: string;           // 推荐理由
}

// 智能图表推荐算法
func (v *Visualizer) RecommendChart(queryIntent *Intent, data *QueryResult) *ChartRecommendation {
    // 根据查询意图和数据特征推荐图表
    switch queryIntent.Type {
    case IntentTypeComparison:
        return &ChartRecommendation{
            ChartType:  'bar',
            Confidence: 0.95,
            Reason:     "对比分析适合使用柱状图",
        }
    case IntentTypeTrend:
        return &ChartRecommendation{
            ChartType:  'line',
            Confidence: 0.98,
            Reason:     "趋势分析适合使用折线图",
        }
    case IntentTypeProportion:
        return &ChartRecommendation{
            ChartType:  'pie',
            Confidence: 0.90,
            Reason:     "占比分析适合使用饼图",
        }
    default:
        return &ChartRecommendation{
            ChartType:  'table',
            Confidence: 0.80,
            Reason:     "默认使用表格展示",
        }
    }
}
```

**图表配置生成** (ECharts格式):
```typescript
function generateBarChart(data: QueryResult): EChartsOption {
    return {
        title: {
            text: data.title,
            left: 'center'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            }
        },
        xAxis: {
            type: 'category',
            data: data.rows.map(row => row.category)
        },
        yAxis: {
            type: 'value',
            name: data.yAxisName
        },
        series: [{
            type: 'bar',
            data: data.rows.map(row => row.value),
            itemStyle: {
                color: '#5470c6'
            }
        }],
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        }
    };
}
```

**验收标准**:
- ✅ 支持6种以上图表类型
- ✅ 图表推荐准确率 > 85%
- ✅ 图表渲染时间 < 1s
- ✅ 支持图表交互 (缩放、筛选、导出)

#### 2.2.4 查询历史与收藏 (P1)

**功能描述**: 保存查询历史，支持收藏常用查询。

**数据模型**:
```sql
CREATE TABLE chatbi_queries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    query_text TEXT NOT NULL,                -- 自然语言查询
    generated_sql TEXT NOT NULL,            -- 生成的SQL
    query_intent VARCHAR(50),                -- 查询意图类型
    data_source_id VARCHAR(64),              -- 数据源ID
    is_favorite BOOLEAN DEFAULT FALSE,       -- 是否收藏
    favorite_name VARCHAR(255),              -- 收藏名称
    execution_count INT DEFAULT 1,           -- 执行次数
    last_executed_at DATETIME,               -- 最后执行时间
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_favorite (is_favorite),
    INDEX idx_last_executed (last_executed_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='问数查询历史表';
```

**验收标准**:
- ✅ 历史查询完整保存
- ✅ 支持收藏、重命名、删除
- ✅ 支持搜索和筛选
- ✅ 支持一键重新执行

---

## 3. 系统架构设计

### 3.1 模块架构图

```
┌─────────────────────────────────────────────────────────┐
│                   ZKER 问数(ChatBI) 模块                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ 前端界面      │  │ NLG引擎      │  │ 查询引擎     │  │
│  │ - 查询输入框  │  │ - 意图识别   │  │ - SQL执行    │  │
│  │ - 图表展示    │  │ - SQL生成    │  │ - 结果缓存   │  │
│  │ - 数据源管理  │  │ - 洞察生成   │  │ - 流式返回   │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │              数据访问层 (Repository)              │  │
│  │  - DataSourceRepository                          │  │
│  │  - QueryHistoryRepository                        │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │              数据存储层                           │  │
│  │  - MySQL (元数据、配置)                           │  │
│  │  - Redis (查询结果缓存)                          │  │
│  │  - ClickHouse (分析数据仓库)                     │  │
│  │  - 外部数据源 (MySQL/PostgreSQL/Excel)           │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 3.2 核心组件设计

#### 3.2.1 NLG引擎 (Natural Language Generation)

**职责**: 将自然语言转换为SQL

**核心类**:
```go
type NLGEngine struct {
    llmClient    *LLMClient             // 大语言模型客户端
    fewShots     *FewShots               // Few-shot示例
    sqlValidator *SQLValidator           // SQL验证器
}

type Intent struct {
    Type        IntentType             // 查询类型
    Metrics     []string               // 指标
    Dimensions  []string               // 维度
    Filters     []Filter               // 过滤条件
    TimeRange   *TimeRange             // 时间范围
}

type IntentType string

const (
    IntentTypeAggregation    IntentType = "aggregation"    // 聚合查询
    IntentTypeComparison     IntentType = "comparison"     // 对比查询
    IntentTypeTrend          IntentType = "trend"          // 趋势查询
    IntentTypeProportion     IntentType = "proportion"     // 占比查询
    IntentTypeList          IntentType = "list"           // 列表查询
)

func (e *NLGEngine) RecognizeIntent(ctx context.Context, query string) (*Intent, error) {
    // 1. 构造Prompt
    prompt := e.buildIntentPrompt(query)

    // 2. 调用LLM
    response, err := e.llmClient.Complete(ctx, prompt)
    if err != nil {
        return nil, err
    }

    // 3. 解析响应
    intent := &Intent{}
    if err := json.Unmarshal([]byte(response), intent); err != nil {
        return nil, err
    }

    return intent, nil
}

func (e *NLGEngine) GenerateSQL(ctx context.Context, intent *Intent, dataSourceID string) (string, error) {
    // 1. 获取数据源元数据
    metadata, err := e.getDataSourceMetadata(ctx, dataSourceID)
    if err != nil {
        return "", err
    }

    // 2. 构造Few-shot提示
    prompt := e.fewShots.BuildSQLPrompt(intent, metadata)

    // 3. 调用LLM生成SQL
    sql, err := e.llmClient.Complete(ctx, prompt)
    if err != nil {
        return "", err
    }

    // 4. 验证SQL
    if err := e.sqlValidator.Validate(sql, metadata.Schema); err != nil {
        // SQL验证失败,尝试修正
        sql, err = e.fixSQL(ctx, sql, err, metadata)
        if err != nil {
            return "", err
        }
    }

    return sql, nil
}

func (e *NLGEngine) GenerateInsight(ctx context.Context, result *QueryResult, intent *Intent) (string, error) {
    // 构造洞察生成Prompt
    prompt := fmt.Sprintf(`
基于以下数据生成简洁的洞察分析:
查询意图: %s
数据行数: %d
数据摘要: %v

请生成3-5条关键洞察,每条洞察包含:
1. 数据趋势
2. 异常值
3. 业务建议
`, intent.Type, len(result.Rows), result.Summary)

    insight, err := e.llmClient.Complete(ctx, prompt)
    return insight, err
}
```

#### 3.2.2 查询执行引擎

**职责**: 执行SQL并返回结果

**核心类**:
```go
type QueryExecutor struct {
    connections map[string]Connection // 数据源连接池
    cache       *Cache                // 结果缓存
    monitor     *Monitor              // 性能监控
}

type Connection interface {
    Execute(ctx context.Context, sql string) (*QueryResult, error)
    Close() error
}

func (e *QueryExecutor) Execute(ctx context.Context, sql string, dataSourceID string) (*QueryResult, error) {
    // 1. 检查缓存
    cacheKey := fmt.Sprintf("chatbi:%s:%s", dataSourceID, sql)
    if cached, found := e.cache.Get(ctx, cacheKey); found {
        return cached.(*QueryResult), nil
    }

    // 2. 获取数据源连接
    conn, err := e.getConnection(dataSourceID)
    if err != nil {
        return nil, err
    }

    // 3. 执行SQL
    startTime := time.Now()
    result, err := conn.Execute(ctx, sql)
    executionTime := time.Since(startTime)

    if err != nil {
        e.monitor.RecordError(dataSourceID, sql, err)
        return nil, err
    }

    // 4. 记录性能指标
    e.monitor.RecordQuery(dataSourceID, sql, executionTime, len(result.Rows))

    // 5. 缓存结果 (缓存5分钟)
    e.cache.Set(ctx, cacheKey, result, 5*time.Minute)

    return result, nil
}

// MySQL连接实现
type MySQLConnection struct {
    db *gorm.DB
}

func NewMySQLConnection(dsn string) (*MySQLConnection, error) {
    db, err := gorm.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    return &MySQLConnection{db: db}, nil
}

func (c *MySQLConnection) Execute(ctx context.Context, sql string) (*QueryResult, error) {
    var rows []map[string]interface{}
    err := c.db.WithContext(ctx).Raw(sql).Scan(&rows).Error
    if err != nil {
        return nil, err
    }

    return &QueryResult{
        Rows:     rows,
        Columns:  extractColumns(rows),
        Metadata: extractMetadata(rows),
    }, nil
}
```

#### 3.2.3 可视化引擎

**职责**: 将查询结果转换为图表配置

**核心类**:
```go
type Visualizer struct {
    chartRegistry *ChartRegistry // 图表注册表
}

type ChartRegistry struct {
    generators map[ChartType]ChartGenerator
}

func NewVisualizer() *Visualizer {
    registry := &ChartRegistry{
        generators: make(map[ChartType]ChartGenerator),
    }

    // 注册图表生成器
    registry.Register(ChartTypeBar, &BarChartGenerator{})
    registry.Register(ChartTypeLine, &LineChartGenerator{})
    registry.Register(ChartTypePie, &PieChartGenerator{})
    registry.Register(ChartTypeTable, &TableGenerator{})
    registry.Register(ChartTypeScatter, &ScatterChartGenerator{})

    return &Visualizer{chartRegistry: registry}
}

func (v *Visualizer) GenerateChart(ctx context.Context, result *QueryResult, intent *Intent) (*Chart, error) {
    // 1. 推荐图表类型
    recommendation := v.RecommendChart(intent, result)

    // 2. 获取图表生成器
    generator := v.chartRegistry.Get(recommendation.ChartType)

    // 3. 生成图表配置
    config, err := generator.Generate(result, intent)
    if err != nil {
        return nil, err
    }

    return &Chart{
        Type:    recommendation.ChartType,
        Config:  config,
        Reason:  recommendation.Reason,
    }, nil
}

// 柱状图生成器
type BarChartGenerator struct{}

func (g *BarChartGenerator) Generate(result *QueryResult, intent *Intent) (interface{}, error) {
    return generateEChartsBarChart(result), nil
}

func generateEChartsBarChart(result *QueryResult) map[string]interface{} {
    return map[string]interface{}{
        "title": map[string]interface{}{
            "text":  result.Title,
            "left":  "center",
        },
        "tooltip": map[string]interface{}{
            "trigger": "axis",
            "axisPointer": map[string]interface{}{
                "type": "shadow",
            },
        },
        "xAxis": map[string]interface{}{
            "type":     "category",
            "data":     extractCategories(result.Rows),
        },
        "yAxis": map[string]interface{}{
            "type":  "value",
            "name":  result.YAxisName,
        },
        "series": []map[string]interface{}{
            {
                "type": "bar",
                "data": extractValues(result.Rows),
                "itemStyle": map[string]interface{}{
                    "color": "#5470c6",
                },
            },
        },
    }
}
```

### 3.3 数据流图

```
用户输入自然语言查询
    ↓
前端发送查询请求
    ↓
┌─────────────────────────────────────────┐
│         ChatBI Controller               │
│  - 接收查询请求                           │
│  - 调用NLG引擎                           │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│         NLG Engine                      │
│  - 意图识别                               │
│  - 实体提取                               │
│  - SQL生成                               │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│      SQL Validator & Optimizer          │
│  - SQL验证                               │
│  - SQL优化                               │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│       Query Executor                     │
│  - 检查缓存                               │
│  - 连接数据源                             │
│  - 执行SQL                                │
│  - 返回结果                               │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│         Visualizer                       │
│  - 推荐图表类型                           │
│  - 生成图表配置                           │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│      Insight Generator                   │
│  - 生成数据洞察                           │
│  - 生成业务建议                           │
└─────────────────────────────────────────┘
    ↓
返回查询结果给前端
    ↓
前端渲染图表和洞察
```

---

## 4. 数据库设计

### 4.1 核心表结构

#### 4.1.1 数据源表 (chatbi_data_sources)

```sql
CREATE TABLE chatbi_data_sources (
    id VARCHAR(64) PRIMARY KEY COMMENT '数据源ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    name VARCHAR(255) NOT NULL COMMENT '数据源名称',
    type ENUM('mysql', 'postgresql', 'clickhouse', 'excel', 'csv', 'api') NOT NULL COMMENT '数据源类型',
    connection JSON NOT NULL COMMENT '连接配置 (加密)',
    metadata JSON COMMENT '元数据 (表结构、字段信息等)',
    permissions JSON COMMENT '权限配置',
    status ENUM('active', 'inactive', 'error') NOT NULL DEFAULT 'active' COMMENT '状态',
    last_test_at DATETIME COMMENT '最后测试时间',
    last_test_result JSON COMMENT '最后测试结果',
    created_by BIGINT NOT NULL COMMENT '创建者ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_tenant_name (tenant_id, name),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='问数数据源表';
```

#### 4.1.2 查询历史表 (chatbi_queries)

```sql
CREATE TABLE chatbi_queries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '查询ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    query_text TEXT NOT NULL COMMENT '自然语言查询',
    generated_sql TEXT NOT NULL COMMENT '生成的SQL',
    query_intent VARCHAR(50) COMMENT '查询意图类型',
    data_source_id VARCHAR(64) COMMENT '数据源ID',
    chart_type VARCHAR(50) COMMENT '图表类型',
    execution_result JSON COMMENT '执行结果摘要',
    is_favorite BOOLEAN DEFAULT FALSE COMMENT '是否收藏',
    favorite_name VARCHAR(255) COMMENT '收藏名称',
    execution_count INT DEFAULT 1 COMMENT '执行次数',
    execution_time_avg INT COMMENT '平均执行时间(ms)',
    last_executed_at DATETIME COMMENT '最后执行时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_data_source (data_source_id),
    INDEX idx_favorite (is_favorite),
    INDEX idx_last_executed (last_executed_at),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (data_source_id) REFERENCES chatbi_data_sources(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='问数查询历史表';
```

#### 4.1.3 查询执行日志表 (chatbi_query_logs)

```sql
CREATE TABLE chatbi_query_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    query_id BIGINT NOT NULL COMMENT '查询ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    sql_text TEXT NOT NULL COMMENT '执行的SQL',
    execution_time INT COMMENT '执行时间(ms)',
    row_count INT COMMENT '返回行数',
    cache_hit BOOLEAN DEFAULT FALSE COMMENT '是否命中缓存',
    error_message TEXT COMMENT '错误信息',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_query_id (query_id),
    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='问数查询执行日志表';
```

### 4.2 索引设计

| 表名 | 索引名 | 类型 | 字段 | 用途 |
|------|--------|------|------|------|
| chatbi_data_sources | uk_tenant_name | UNIQUE | tenant_id, name | 租户内名称唯一 |
| chatbi_data_sources | idx_tenant_id | INDEX | tenant_id | 按租户查询 |
| chatbi_data_sources | idx_type | INDEX | type | 按类型筛选 |
| chatbi_queries | idx_tenant_user | INDEX | tenant_id, user_id | 查询用户历史 |
| chatbi_queries | idx_favorite | INDEX | is_favorite | 收藏查询 |
| chatbi_queries | idx_last_executed | INDEX | last_executed_at | 最近执行 |
| chatbi_query_logs | idx_query_id | INDEX | query_id | 查询日志 |
| chatbi_query_logs | idx_created_at | INDEX | created_at | 时间范围查询 |

---

## 5. 接口设计

### 5.1 REST API

#### 5.1.1 数据源管理API

```go
// 数据源管理路由
chatbiGroup := app.Group("/api/chatbi")
{
    // 数据源管理
    chatbiGroup.POST("/datasources", h.CreateDataSource)          // 创建数据源
    chatbiGroup.GET("/datasources", h.ListDataSources)            // 数据源列表
    chatbiGroup.GET("/datasources/:id", h.GetDataSource)          // 数据源详情
    chatbiGroup.PUT("/datasources/:id", h.UpdateDataSource)       // 更新数据源
    chatbiGroup.DELETE("/datasources/:id", h.DeleteDataSource)    // 删除数据源
    chatbiGroup.POST("/datasources/:id/test", h.TestDataSource)    // 测试连接

    // 查询执行
    chatbiGroup.POST("/query", h.ExecuteQuery)                   // 执行查询
    chatbiGroup.POST("/query/stream", h.ExecuteQueryStream)     // 流式查询

    // 历史管理
    chatbiGroup.GET("/history", h.GetQueryHistory)               // 查询历史
    chatbiGroup.POST("/history/:id/favorite", h.FavoriteQuery)   // 收藏查询
    chatbiGroup.DELETE("/history/:id", h.DeleteQuery)           // 删除历史
}
```

**API详细定义**:

**POST /api/chatbi/query - 执行查询**
```typescript
// 请求
interface ExecuteQueryRequest {
  query: string;              // 自然语言查询
  dataSourceId?: string;       // 数据源ID (可选)
  chartType?: ChartType;      // 图表类型 (可选)
}

// 响应
interface ExecuteQueryResponse {
  queryId: number;            // 查询ID
  sql: string;                // 生成的SQL
  result: {
    columns: string[];        // 列名
    rows: any[];             // 数据行
    rowCount: number;        // 行数
  };
  chart: {
    type: ChartType;         // 图表类型
    config: any;             // ECharts配置
    reason: string;          // 推荐理由
  };
  insight: string;            // 数据洞察
  metadata: {
    executionTime: number;   // 执行时间(ms)
    cached: boolean;         // 是否缓存
  };
}
```

**示例请求**:
```json
POST /api/chatbi/query
{
  "query": "上个月华东区的销售额是多少？",
  "dataSourceId": "ds_123"
}
```

**示例响应**:
```json
{
  "queryId": 1001,
  "sql": "SELECT region, SUM(sales) as total_sales FROM sales WHERE date >= '2025-11-01' AND date <= '2025-11-30' AND region = '华东' GROUP BY region",
  "result": {
    "columns": ["region", "total_sales"],
    "rows": [
      {"region": "华东", "total_sales": 1234567}
    ],
    "rowCount": 1
  },
  "chart": {
    "type": "bar",
    "config": {
      "title": {"text": "华东区销售额", "left": "center"},
      "xAxis": {"type": "category", "data": ["华东"]},
      "yAxis": {"type": "value", "name": "销售额"},
      "series": [{"type": "bar", "data": [1234567]}]
    },
    "reason": "单值对比适合使用柱状图"
  },
  "insight": "华东区上月销售额为 ¥1,234,567",
  "metadata": {
    "executionTime": 1250,
    "cached": false
  }
}
```

#### 5.1.2 流式查询API

**POST /api/chatbi/query/stream - 流式查询**
```typescript
// 请求
interface ExecuteQueryStreamRequest {
  query: string;
  dataSourceId?: string;
}

// 响应 (Server-Sent Events)
event: progress
data: {
  step: "recognize_intent" | "generate_sql" | "execute_query" | "generate_chart" | "done"
  progress: number;  // 0-100
  message: string;
  data?: any;
}

event: result
data: {
  chart: any;
  insight: string;
  metadata: any;
}
```

**前端使用示例**:
```typescript
const eventSource = new EventSource(
  `/api/chatbi/query/stream?query=${encodeURIComponent(query)}`
);

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);

  switch (data.step) {
    case 'recognize_intent':
      updateProgress('正在理解您的查询...', 20);
      break;
    case 'generate_sql':
      updateProgress('正在生成查询语句...', 40);
      break;
    case 'execute_query':
      updateProgress('正在查询数据...', 70);
      break;
    case 'generate_chart':
      updateProgress('正在生成图表...', 90);
      break;
    case 'done':
      renderResult(data.data);
      break;
  }
};
```

### 5.2 WebSocket API

**连接**: `wss://{subdomain}.zker.com/ws/chatbi`

**消息协议**:
```typescript
// 客户端发送查询
{
  "event": "query",
  "data": {
    "query": "上个月销售额TOP5的地区",
    "dataSourceId": "ds_123"
  }
}

// 服务端推送进度
{
  "event": "progress",
  "data": {
    "step": "execute_query",
    "progress": 60,
    "message": "正在查询数据..."
  }
}

// 服务端推送结果
{
  "event": "result",
  "data": {
    "chart": {...},
    "insight": "..."
  }
}
```

---

## 6. 前端设计

### 6.1 页面结构

```
frontend/apps/zker-workspace/src/pages/chatbi/
├── index.tsx                 # 问数主页
├── components/
│   ├── QueryInput.tsx        # 查询输入框
│   ├── ChartDisplay.tsx      # 图表展示
│   ├── InsightPanel.tsx     # 洞察面板
│   └── DataSourceManager.tsx # 数据源管理
├── hooks/
│   ├── useQuery.ts           # 查询Hook
│   ├── useStreamQuery.ts    # 流式查询Hook
│   └── useDataSource.ts     # 数据源Hook
└── types/
    ├── query.ts              # 查询类型定义
    ├── chart.ts              # 图表类型定义
    └── datasource.ts         # 数据源类型定义
```

### 6.2 核心组件

#### 6.2.1 QueryInput 组件

```typescript
interface QueryInputProps {
  onSubmit: (query: string) => void;
  dataSourceId?: string;
  loading?: boolean;
}

const QueryInput: React.FC<QueryInputProps> = ({ onSubmit, dataSourceId, loading }) => {
  const [value, setValue] = useState('');
  const [suggestions, setSuggestions] = useState<string[]>([]);

  const handleSubmit = () => {
    if (value.trim()) {
      onSubmit(value);
    }
  };

  const handleInputChange = (input: string) => {
    setValue(input);
    // TODO: 实现查询建议
    // fetchSuggestions(input);
  };

  return (
    <div className="query-input-container">
      <TextArea
        value={value}
        onChange={(e) => handleInputChange(e.target.value)}
        placeholder="请输入您的数据查询，例如：上个月华东区的销售额是多少？"
        autoSize={{ minRows: 2, maxRows: 6 }}
        onPressEnter={handleSubmit}
      />
      <Button
        type="primary"
        onClick={handleSubmit}
        loading={loading}
        icon={<SearchOutlined />}
      >
        查询
      </Button>
    </div>
  );
};
```

#### 6.2.2 ChartDisplay 组件

```typescript
interface ChartDisplayProps {
  chartConfig: any;
  chartType: ChartType;
  loading?: boolean;
}

const ChartDisplay: React.FC<ChartDisplayProps> = ({ chartConfig, chartType, loading }) => {
  const chartRef = useRef<ReactEChartsInstance>(null);

  useEffect(() => {
    if (chartRef.current && chartConfig) {
      chartRef.current.setOption(chartConfig);
    }
  }, [chartConfig]);

  const onEvents = {
    click: (params: any) => {
      console.log('Chart clicked:', params);
      // TODO: 实现图表交互 (下钻、筛选等)
    },
    legendselectchanged: (params: any) => {
      console.log('Legend select changed:', params);
      // TODO: 实现图例筛选
    },
  };

  if (loading) {
    return <Skeleton active />;
  }

  return (
    <div className="chart-display">
      <ReactECharts
        ref={chartRef}
        option={chartConfig}
        onEvents={onEvents}
        style={{ height: '400px' }}
        notMerge={true}
        lazyUpdate={true}
      />
    </div>
  );
};
```

### 6.3 状态管理

```typescript
// stores/chatbiStore.ts
interface ChatBIState {
  currentQuery: string;
  queryResult: QueryResult | null;
  chartConfig: any;
  insight: string;
  loading: boolean;
  error: string | null;
  dataSourceId: string | null;
  queryHistory: QueryHistoryItem[];
}

interface ChatBIActions {
  executeQuery: (query: string) => Promise<void>;
  favoriteQuery: (queryId: number) => Promise<void>;
  setDataSource: (dataSourceId: string) => void;
}

export const useChatBIStore = create<ChatBIState & ChatBIActions>((set) => ({
  currentQuery: '',
  queryResult: null,
  chartConfig: null,
  insight: '',
  loading: false,
  error: null,
  dataSourceId: null,
  queryHistory: [],

  executeQuery: async (query: string) => {
    set({ loading: true, currentQuery: query, error: null });

    try {
      const response = await api.post('/api/chatbi/query', { query });

      set({
        queryResult: response.data.result,
        chartConfig: response.data.chart.config,
        insight: response.data.insight,
        loading: false,
      });

      // 添加到历史记录
      set((state) => ({
        queryHistory: [
          {
            id: response.data.queryId,
            query: query,
            timestamp: new Date(),
            favorite: false,
          },
          ...state.queryHistory,
        ].slice(0, 50), // 保留最近50条
      }));
    } catch (error: any) {
      set({
        error: error.message,
        loading: false,
      });
    }
  },

  favoriteQuery: async (queryId: number) => {
    await api.post(`/api/chatbi/history/${queryId}/favorite`);
    set((state) => ({
      queryHistory: state.queryHistory.map((item) =>
        item.id === queryId ? { ...item, favorite: true } : item
      ),
    }));
  },

  setDataSource: (dataSourceId: string) => {
    set({ dataSourceId });
  },
}));
```

---

## 7. 后端设计

### 7.1 服务层设计

#### 7.1.1 ChatBIService

```go
type ChatBIService struct {
    nlgEngine       *NLGEngine
    queryExecutor   *QueryExecutor
    visualizer      *Visualizer
    insightEngine   *InsightEngine
    dataSourceRepo  *DataSourceRepository
    queryHistoryRepo *QueryHistoryRepository
    permissionCheck *PermissionChecker
}

func NewChatBIService(
    nlgEngine *NLGEngine,
    queryExecutor *QueryExecutor,
    visualizer *Visualizer,
    insightEngine *InsightEngine,
    dataSourceRepo *DataSourceRepository,
    queryHistoryRepo *QueryHistoryRepository,
    permissionCheck *PermissionChecker,
) *ChatBIService {
    return &ChatBIService{
        nlgEngine:       nlgEngine,
        queryExecutor:   queryExecutor,
        visualizer:      visualizer,
        insightEngine:   insightEngine,
        dataSourceRepo:  dataSourceRepo,
        queryHistoryRepo: queryHistoryRepo,
        permissionCheck: permissionCheck,
    }
}

func (s *ChatBIService) ExecuteQuery(ctx context.Context, req *ExecuteQueryRequest) (*ExecuteQueryResponse, error) {
    // 1. 权限检查
    dataSourceID := req.DataSourceID
    if dataSourceID == "" {
        // 使用默认数据源
        defaultDS, err := s.dataSourceRepo.GetDefaultDataSource(ctx, req.TenantID)
        if err != nil {
            return nil, err
        }
        dataSourceID = defaultDS.ID
    }

    // 检查数据源访问权限
    if err := s.permissionCheck.CheckDataSourceAccess(ctx, req.UserID, dataSourceID); err != nil {
        return nil, status.Error(codes.PermissionDenied, "无数据源访问权限")
    }

    // 2. 意图识别
    intent, err := s.nlgEngine.RecognizeIntent(ctx, req.Query)
    if err != nil {
        return nil, err
    }

    // 3. 实体提取
    entities := s.nlgEngine.ExtractEntities(ctx, req.Query, intent)

    // 4. SQL生成
    sql, err := s.nlgEngine.GenerateSQL(ctx, intent, entities, dataSourceID)
    if err != nil {
        return nil, err
    }

    // 5. SQL优化
    sql = s.optimizeSQL(sql)

    // 6. 执行查询
    startTime := time.Now()
    result, err := s.queryExecutor.Execute(ctx, sql, dataSourceID)
    if err != nil {
        return nil, err
    }
    executionTime := time.Since(startTime).Milliseconds()

    // 7. 生成图表
    chart, err := s.visualizer.GenerateChart(ctx, result, intent)
    if err != nil {
        return nil, err
    }

    // 8. 生成洞察
    insight, err := s.insightEngine.Generate(ctx, result, intent)
    if err != nil {
        insight = "查询完成"
    }

    // 9. 保存查询历史
    queryHistory := &QueryHistory{
        TenantID:       req.TenantID,
        UserID:         req.UserID,
        QueryText:      req.Query,
        GeneratedSQL:   sql,
        QueryIntent:    intent.Type.String(),
        DataSourceID:   dataSourceID,
        ChartType:      chart.Type,
        ExecutionTime:  executionTime,
        RowCount:       result.RowCount,
    }
    s.queryHistoryRepo.Save(ctx, queryHistory)

    return &ExecuteQueryResponse{
        QueryID:     queryHistory.ID,
        SQL:         sql,
        Result:      result,
        Chart:       chart,
        Insight:     insight,
        Metadata: Metadata{
            ExecutionTime: executionTime,
            Cached:        false,
        },
    }, nil
}
```

#### 7.1.2 DataSourceService

```go
type DataSourceService struct {
    repo           *DataSourceRepository
    permissionSvc  *PermissionService
    validator      *DataSourceValidator
    encryptor      *Encryptor
}

func (s *DataSourceService) CreateDataSource(ctx context.Context, req *CreateDataSourceRequest) (*DataSource, error) {
    // 1. 验证连接信息
    if err := s.validator.ValidateConnection(req.Type, req.Connection); err != nil {
        return nil, err
    }

    // 2. 加密敏感信息
    encryptedConn, err := s.encryptor.Encrypt(req.Connection)
    if err != nil {
        return nil, err
    }

    // 3. 提取元数据
    metadata, err := s.extractMetadata(ctx, req)
    if err != nil {
        return nil, err
    }

    // 4. 创建数据源
    dataSource := &DataSource{
        ID:          uuid.New().String(),
        TenantID:    req.TenantID,
        Name:        req.Name,
        Type:        req.Type,
        Connection:  encryptedConn,
        Metadata:    metadata,
        Permissions: req.Permissions,
        Status:      "active",
        CreatedBy:   req.CreatedBy,
    }

    if err := s.repo.Create(ctx, dataSource); err != nil {
        return nil, err
    }

    return dataSource, nil
}

func (s *DataSourceService) TestConnection(ctx context.Context, dataSourceID string) error {
    // 1. 获取数据源
    dataSource, err := s.repo.GetByID(ctx, dataSourceID)
    if err != nil {
        return err
    }

    // 2. 解密连接信息
    conn, err := s.encryptor.Decrypt(dataSource.Connection)
    if err != nil {
        return err
    }

    // 3. 测试连接
    var err error
    switch dataSource.Type {
    case DataSourceTypeMySQL:
        err = s.testMySQLConnection(ctx, conn)
    case DataSourceTypePostgreSQL:
        err = s.testPostgreSQLConnection(ctx, conn)
    case DataSourceTypeClickHouse:
        err = s.testClickHouseConnection(ctx, conn)
    default:
        return errors.New("unsupported data source type")
    }

    // 4. 更新测试结果
    dataSource.LastTestAt = time.Now()
    dataSource.LastTestResult = map[string]interface{}{
        "success": err == nil,
        "error":   err,
    }
    if err == nil {
        dataSource.Status = "active"
    } else {
        dataSource.Status = "error"
    }

    s.repo.Update(ctx, dataSource)

    return err
}
```

### 7.2 中间件设计

#### 7.2.1 权限检查中间件

```go
type DataSourcePermissionMiddleware struct {
    permissionSvc *PermissionService
}

func (m *DataSourcePermissionMiddleware) CheckPermission() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        userID := ctx.Value("user_id").(string)
        tenantID := ctx.Value("tenant_id").(string)
        dataSourceID := c.Query("dataSourceId")

        // 检查数据源权限
        hasAccess, err := m.permissionSvc.CheckDataSourceAccess(ctx, userID, tenantID, dataSourceID)
        if err != nil || !hasAccess {
            c.JSON(403, map[string]string{
                "error": "无数据源访问权限",
            })
            c.Abort()
            return
        }

        c.Next(ctx)
    }
}
```

---

## 8. 性能优化

### 8.1 缓存策略

#### 8.1.1 查询结果缓存

```go
type QueryCache struct {
    redis      *redis.Client
    ttl        time.Duration
}

const (
    cacheKeyPrefix = "chatbi:query:"
)

func (c *QueryCache) Get(ctx context.Context, sql string, dataSourceID string) (*QueryResult, error) {
    key := fmt.Sprintf("%s%s:%s", cacheKeyPrefix, dataSourceID, sqlHash(sql))

    cached, err := c.redis.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }

    var result QueryResult
    if err := json.Unmarshal([]byte(cached), &result); err != nil {
        return nil, err
    }

    return &result, nil
}

func (c *QueryCache) Set(ctx context.Context, sql string, dataSourceID string, result *QueryResult) error {
    key := fmt.Sprintf("%s%s:%s", cacheKeyPrefix, dataSourceID, sqlHash(sql))

    data, err := json.Marshal(result)
    if err != nil {
        return err
    }

    return c.redis.Set(ctx, key, data, c.ttl).Err()
}

func sqlHash(sql string) string {
    h := sha256.Sum256([]byte(sql))
    return hex.EncodeToString(h[:])
}
```

**缓存策略**:
- 相同SQL查询 → 缓存命中,直接返回
- 缓存有效期: 5分钟
- 缓存容量: 10000个查询结果
- 缓存淘汰: LRU算法

#### 8.1.2 元数据缓存

```go
type MetadataCache struct {
    redis *redis.Client
    ttl   time.Duration
}

func (c *MetadataCache) GetDataSourceMetadata(ctx context.Context, dataSourceID string) (*DataSourceMetadata, error) {
    key := fmt.Sprintf("chatbi:metadata:%s", dataSourceID)

    cached, err := c.redis.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }

    var metadata DataSourceMetadata
    if err := json.Unmarshal([]byte(cached), &metadata); err != nil {
        return nil, err
    }

    return &metadata, nil
}

func (c *MetadataCache) SetDataSourceMetadata(ctx context.Context, dataSourceID string, metadata *DataSourceMetadata) error {
    key := fmt.Sprintf("chatbi:metadata:%s", dataSourceID)

    data, err := json.Marshal(metadata)
    if err != nil {
        return err
    }

    return c.redis.Set(ctx, key, data, c.ttl).Err()
}
```

### 8.2 SQL优化

#### 8.2.1 SQL优化器

```go
type SQLOptimizer struct {
    rules []SQLOptimizationRule
}

type SQLOptimizationRule interface {
    Optimize(sql string) (string, error)
    Name() string
}

// Limit优化规则
type LimitOptimizationRule struct{}

func (r *LimitOptimizationRule) Optimize(sql string) (string, error) {
    // 自动添加LIMIT限制,防止查询返回过多数据
    if !strings.Contains(strings.ToUpper(sql), "LIMIT") {
        sql += " LIMIT 10000"
    }
    return sql, nil
}

func (r *LimitOptimizationRule) Name() string {
    return "LimitOptimization"
}

// 索引提示优化规则
type IndexHintOptimizationRule struct{}

func (r *IndexHintOptimizationRule) Optimize(sql string) (string, error) {
    // 分析SQL,添加索引提示
    // SELECT * FROM users WHERE name = 'xxx'
    // → SELECT * FROM users USE INDEX (idx_name) WHERE name = 'xxx'
    return sql, nil
}

func (r *IndexHintOptimizationRule) Name() string {
    return "IndexHintOptimization"
}
```

#### 8.2.2 查询并行化

```go
// 对大数据量查询,使用并行执行
func (e *QueryExecutor) executeParallel(ctx context.Context, queries []string, dataSourceID string) (*QueryResult, error) {
    var wg sync.WaitGroup
    results := make(chan *QueryResult, len(queries))
    errors := make(chan error, len(queries))

    for _, query := range queries {
        wg.Add(1)
        go func(sql string) {
            defer wg.Done()
            result, err := e.Execute(ctx, sql, dataSourceID)
            if err != nil {
                errors <- err
            } else {
                results <- result
            }
        }(query)
    }

    go func() {
        wg.Wait()
        close(results)
        close(errors)
    }()

    // 收集结果
    var allRows []map[string]interface{}
    for result := range results {
        allRows = append(allRows, result.Rows...)
    }

    // 检查错误
    if len(errors) > 0 {
        return nil, <-errors
    }

    return &QueryResult{
        Rows: allRows,
    }, nil
}
```

---

## 9. 测试用例

### 9.1 单元测试

#### 9.1.1 NLG引擎测试

```go
func TestNLGEngine_RecognizeIntent(t *testing.T) {
    engine := NewNLGEngine(testLLMClient)

    tests := []struct {
        name     string
        query    string
        expected IntentType
    }{
        {
            name:     "聚合查询",
            query:    "上个月的总销售额是多少？",
            expected: IntentTypeAggregation,
        },
        {
            name:     "对比查询",
            query:    "对比各地区的销售业绩",
            expected: IntentTypeComparison,
        },
        {
            name:     "趋势查询",
            query:    "显示今年的销售趋势",
            expected: IntentTypeTrend,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            intent, err := engine.RecognizeIntent(context.Background(), tt.query)
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, intent.Type)
        })
    }
}

func TestNLGEngine_GenerateSQL(t *testing.T) {
    engine := NewNLGEngine(testLLMClient, testSQLParser)

    intent := &Intent{
        Type:       IntentTypeAggregation,
        Metrics:    []string{"sales"},
        Dimensions: []string{"region"},
        TimeRange:  &TimeRange{
            Start: parseDate("2025-11-01"),
            End:   parseDate("2025-11-30"),
        },
    }

    metadata := &DataSourceMetadata{
        Schema: &TableSchema{
            TableName: "sales",
            Columns: []Column{
                {Name: "region", Type: "string"},
                {Name: "sales", Type: "decimal"},
                {Name: "date", Type: "date"},
            },
        },
    }

    sql, err := engine.GenerateSQL(context.Background(), intent, metadata, "ds_123")
    assert.NoError(t, err)

    expectedSQL := `SELECT region, SUM(sales) as total_sales FROM sales WHERE date >= '2025-11-01' AND date <= '2025-11-30' GROUP BY region`
    assert.Equal(t, expectedSQL, sql)
}
```

#### 9.1.2 查询执行器测试

```go
func TestQueryExecutor_Execute(t *testing.T) {
    executor := NewQueryExecutor(testConnections, testCache)

    sql := `SELECT region, SUM(amount) as total FROM sales GROUP BY region`

    result, err := executor.Execute(context.Background(), sql, "ds_mysql")
    assert.NoError(t, err)

    assert.NotNil(t, result)
    assert.Greater(t, len(result.Rows), 0)
}

func TestQueryExecutor_Cache(t *testing.T) {
    executor := NewQueryExecutor(testConnections, testCache)

    sql := `SELECT COUNT(*) FROM users`

    // 第一次执行 - 未命中缓存
    result1, err := executor.Execute(context.Background(), sql, "ds_mysql")
    assert.NoError(t, err)

    startTime := time.Now()
    // 第二次执行 - 命中缓存
    result2, err := executor.Execute(context.Background(), sql, "ds_mysql")
    duration := time.Since(startTime)

    assert.NoError(t, err)
    assert.Equal(t, result1.RowCount, result2.RowCount)
    assert.Less(t, duration.Milliseconds(), 50) // 缓存响应应该 < 50ms
}
```

### 9.2 集成测试

```go
func TestChatBIService_ExecuteQuery_EndToEnd(t *testing.T) {
    service := setupTestService()

    req := &ExecuteQueryRequest{
        Query:       "上个月华东区的销售额是多少？",
        DataSourceID: "ds_mysql",
        TenantID:    "tenant-001",
        UserID:      123,
    }

    response, err := service.ExecuteQuery(context.Background(), req)
    assert.NoError(t, err)

    // 验证响应
    assert.NotEmpty(t, response.SQL)
    assert.NotNil(t, response.Result)
    assert.NotEmpty(t, response.Chart)
    assert.NotEmpty(t, response.Insight)
}
```

### 9.3 性能测试

```go
func BenchmarkQueryExecution(b *testing.B) {
    executor := NewQueryExecutor(testConnections, testCache)
    sql := `SELECT region, SUM(amount) FROM sales WHERE date >= '2025-01-01' GROUP BY region`

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := executor.Execute(context.Background(), sql, "ds_mysql")
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkNLGEngine_GenerateSQL(b *testing.B) {
    engine := NewNLGEngine(testLLMClient)
    intent := &Intent{Type: IntentTypeAggregation}
    metadata := testMetadata

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := engine.GenerateSQL(context.Background(), intent, metadata, "ds_123")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

---

## 附录

### 附录A: 错误码

| 错误码 | 错误信息 | HTTP状态码 | 说明 |
|--------|---------|-----------|------|
| CHATBI_001 | 数据源不存在 | 404 | 指定的数据源ID不存在 |
| CHATBI_002 | 无数据源访问权限 | 403 | 用户无权访问该数据源 |
| CHATBI_003 | 自然语言理解失败 | 400 | 无法理解用户的查询意图 |
| CHATBI_004 | SQL生成失败 | 500 | 无法生成有效的SQL语句 |
| CHATBI_005 | SQL执行失败 | 500 | SQL执行出错 |
| CHATBI_006 | 数据源连接失败 | 500 | 无法连接到数据源 |
| CHATBI_007 | 查询结果为空 | 200 | 查询返回0行数据 |

### 附录B: 配置示例

**数据源配置示例 (MySQL)**:
```json
{
  "name": "业务数据库",
  "type": "mysql",
  "connection": {
    "host": "mysql.example.com",
    "port": 3306,
    "database": "business_db",
    "username": "readonly_user",
    "password": "encrypted_password"
  },
  "permissions": {
    "userIds": ["user-001", "user-002"],
    "roleIds": ["role-analyst"]
  }
}
```

**数据源配置示例 (ClickHouse)**:
```json
{
  "name": "数据分析仓库",
  "type": "clickhouse",
  "connection": {
    "host": "clickhouse.example.com",
    "port": 8123,
    "database": "analytics",
    "username": "analyst",
    "password": "encrypted_password"
  },
  "permissions": {
    "userIds": [],
    "roleIds": ["role-data-analyst"]
  }
}
```

---

**文档结束**

> 本文档详细定义了 ZKER 问数(ChatBI) 模块的完整设计,包含功能规格、架构设计、数据库设计、接口设计、前端设计、后端设计、性能优化和测试用例,为开发和实施提供了完整指导。
