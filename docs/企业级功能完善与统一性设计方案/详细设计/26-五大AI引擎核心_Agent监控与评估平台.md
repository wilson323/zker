# 26-五大AI引擎核心_Agent监控与评估平台 详细设计说明书

**文档编号**: DE-DD-2025-026
**模块名称**: AI Agent Monitoring & Evaluation Platform (AgentMonitor)
**版本**: v1.0.0
**作者**: ZKER Enterprise Team
**创建日期**: 2025-01-03

---

## 1. 模块概述

**AI Agent 监控与评估平台** 是对标 **LangSmith** 的生产级 AI Agent 应用监控与评估系统，提供**全流程可观测性 + 质量评估 + 成本分析**能力。

**核心设计理念**：
- ✅ **全流程可观测性** - 可视化追踪树，清晰展现每一步执行逻辑
- ✅ **实时监控** - 性能追踪、错误追踪、用户行为分析
- ✅ **质量评估** - 12种评估方法（精确匹配、基于事实问答等）
- ✅ **版本管理** - 提示词版本管理、A/B测试、效果对比

**实现策略**：✅ 40% 编码（监控框架） + 60% 配置（评估规则）

**对标产品**: [LangSmith](https://www.ibm.com/cn-zh/think/topics/langsmith)

---

## 2. 核心功能

### 2.1 Agent 可观测性

**F1 - 调用链追踪（Trace）**
- F1.1 可视化追踪树
- F1.2 每一步执行详情（输入、输出、耗时）
- F1.3 错误定位与堆栈追踪

**F2 - 性能监控**
- F2.1 Token 消耗统计（输入/输出/总计）
- F2.2 响应时间分析（P50、P95、P99）
- F2.3 吞吐量监控（QPS）

**F3 - 错误追踪**
- F3.1 错误率监控
- F3.2 错误分类（超时、API失败、幻觉等）
- F3.3 错误告警

### 2.2 Agent 质量评估

**F4 - 自动化评估**
- F4.1 准确性评估（与标准答案对比）
- F4.2 相关性评估（上下文匹配度）
- F4.3 安全性评估（敏感词检测）
- F4.4 幻觉检测（事实性验证）

**F5 - 用户反馈评估**
- F5.1 👍👎 反馈收集
- F5.2 1-5星评分
- F5.3 反馈原因分析

**F6 - 评估报告**
- F6.1 综合评分（0-100分）
- F6.2 问题诊断
- F6.3 优化建议

### 2.3 提示词工程

**F7 - 提示词版本管理**
- F7.1 提示词版本控制（Git-like）
- F7.2 版本对比（Diff视图）
- F7.3 回滚能力

**F8 - A/B 测试**
- F8.1 多版本并行测试
- F8.2 效果对比分析
- F8.3 自动选择最佳版本

### 2.4 成本分析

**F9 - 成本追踪**
- F9.1 实时成本追踪（按会话、按Bot、按用户）
- F9.2 成本预测（基于使用趋势）
- F9.3 成本优化建议（缓存命中率、模型选择）

**F10 - 成本分摊**
- F10.1 按Bot分摊
- F10.2 按用户分摊
- F10.3 按部门分摊

---

## 3. 数据库设计

### 3.1 核心表结构

#### 3.1.1 Agent 执行追踪表 (agent_executions)

```sql
CREATE TABLE agent_executions (
    id VARCHAR(64) PRIMARY KEY COMMENT '执行ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    message_id VARCHAR(64) NOT NULL COMMENT '消息ID',

    -- 执行信息
    input_text TEXT NOT NULL COMMENT '输入文本',
    output_text TEXT COMMENT '输出文本',
    status ENUM('pending', 'running', 'success', 'failed', 'timeout') NOT NULL COMMENT '状态',

    -- Token 统计
    input_tokens INT DEFAULT 0 COMMENT '输入Token数',
    output_tokens INT DEFAULT 0 COMMENT '输出Token数',
    total_tokens INT DEFAULT 0 COMMENT '总Token数',

    -- 性能
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME COMMENT '结束时间',
    duration_ms INT COMMENT '耗时(毫秒)',

    -- 模型信息
    model_provider VARCHAR(50) COMMENT '模型提供商',
    model_name VARCHAR(50) COMMENT '模型名称',

    -- 错误信息
    error_type VARCHAR(64) COMMENT '错误类型',
    error_message TEXT COMMENT '错误信息',
    error_stack TEXT COMMENT '错误堆栈',

    -- 成本
    cost DECIMAL(10,6) COMMENT '本次调用成本',

    -- 元数据
    metadata JSON COMMENT '元数据（JSON格式）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_bot (tenant_id, bot_id),
    INDEX idx_user_conversation (user_id, conversation_id),
    INDEX idx_status (status),
    INDEX idx_start_time (start_time),
    INDEX idx_duration_ms (duration_ms)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Agent执行追踪表';
```

#### 3.1.2 Agent 执行步骤表 (agent_execution_steps)

```sql
CREATE TABLE agent_execution_steps (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '步骤ID',
    execution_id VARCHAR(64) NOT NULL COMMENT '执行ID',
    step_name VARCHAR(128) NOT NULL COMMENT '步骤名称',
    step_type VARCHAR(64) NOT NULL COMMENT '步骤类型（llm_call/tool_call/routing等）',
    parent_step_id BIGINT DEFAULT NULL COMMENT '父步骤ID',

    -- 步骤内容
    input_data JSON COMMENT '输入数据',
    output_data JSON COMMENT '输出数据',

    -- 性能
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME COMMENT '结束时间',
    duration_ms INT COMMENT '耗时(毫秒)',

    -- Token 统计（仅LLM调用步骤）
    input_tokens INT DEFAULT 0,
    output_tokens INT DEFAULT 0,
    total_tokens INT DEFAULT 0,

    -- 模型信息（仅LLM调用步骤）
    model_provider VARCHAR(50),
    model_name VARCHAR(50),
    cost DECIMAL(10,6),

    -- 错误信息
    status ENUM('pending', 'running', 'success', 'failed') NOT NULL,
    error_message TEXT,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_execution_id (execution_id),
    INDEX idx_parent_step (parent_step_id),
    INDEX idx_step_type (step_type),
    INDEX idx_status (status),
    FOREIGN KEY (execution_id) REFERENCES agent_executions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Agent执行步骤表';
```

#### 3.1.3 用户反馈表 (user_feedback)

```sql
CREATE TABLE user_feedback (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '反馈ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    execution_id VARCHAR(64) NOT NULL COMMENT '执行ID',

    -- 反馈内容
    rating INT CHECK (rating >= 1 AND rating <= 5) COMMENT '评分（1-5星）',
    is_helpful BOOLEAN COMMENT '是否有用',
    feedback_reason VARCHAR(100) COMMENT '反馈原因',
    feedback_text TEXT COMMENT '补充说明',

    -- 标签
    tags JSON COMMENT '标签（JSON数组）',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_execution_id (execution_id),
    INDEX idx_rating (rating),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户反馈表';
```

#### 3.1.4 提示词版本表 (prompt_versions)

```sql
CREATE TABLE prompt_versions (
    id VARCHAR(64) PRIMARY KEY COMMENT '版本ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    version VARCHAR(32) NOT NULL COMMENT '版本号（如v1.0.0）',

    -- 提示词内容
    system_prompt TEXT NOT NULL COMMENT '系统提示词',
    user_prompt_template TEXT COMMENT '用户提示词模板',
    parameters JSON COMMENT '参数定义（JSON Schema）',

    -- 版本信息
    parent_version_id VARCHAR(64) COMMENT '父版本ID',
    is_active BOOLEAN DEFAULT FALSE COMMENT '是否为当前激活版本',
    change_description TEXT COMMENT '变更说明',

    -- 评估结果
    avg_rating DECIMAL(3,2) COMMENT '平均评分',
    total_evaluations INT DEFAULT 0 COMMENT '总评估次数',
    success_rate DECIMAL(5,2) COMMENT '成功率',

    created_by BIGINT NOT NULL COMMENT '创建人ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_bot_version (bot_id, version),
    INDEX idx_tenant_bot (tenant_id, bot_id),
    INDEX idx_is_active (is_active),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='提示词版本表';
```

#### 3.1.5 A/B 测试表 (ab_tests)

```sql
CREATE TABLE ab_tests (
    id VARCHAR(64) PRIMARY KEY COMMENT '测试ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    test_name VARCHAR(255) NOT NULL COMMENT '测试名称',

    -- 测试配置
    control_version_id VARCHAR(64) NOT NULL COMMENT '对照版本ID',
    treatment_version_ids JSON NOT NULL COMMENT '实验版本ID列表（JSON数组）',
    traffic_split JSON NOT NULL COMMENT '流量分配（如: {"v1": 50, "v2": 50}）',

    -- 测试状态
    status ENUM('draft', 'running', 'paused', 'completed') DEFAULT 'draft' COMMENT '状态',
    start_time DATETIME COMMENT '开始时间',
    end_time DATETIME COMMENT '结束时间',

    -- 测试结果
    statistical_significance DECIMAL(5,4) COMMENT '统计显著性',
    winner_version_id VARCHAR(64) COMMENT '获胜版本ID',

    created_by BIGINT NOT NULL COMMENT '创建人ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_bot (tenant_id, bot_id),
    INDEX idx_status (status),
    INDEX idx_start_time (start_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='A/B测试表';
```

#### 3.1.6 成本分析表 (cost_analytics)

```sql
CREATE TABLE cost_analytics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    date DATE NOT NULL COMMENT '日期',

    -- 成本维度
    dimension_type ENUM('total', 'by_bot', 'by_user', 'by_department') NOT NULL COMMENT '维度类型',
    dimension_id VARCHAR(64) COMMENT '维度ID（Bot ID / User ID / Department ID）',

    -- 成本统计
    total_requests INT DEFAULT 0 COMMENT '总请求数',
    total_tokens BIGINT DEFAULT 0 COMMENT '总Token数',
    total_cost DECIMAL(12,2) DEFAULT 0.00 COMMENT '总成本',

    -- 成本明细
    input_cost DECIMAL(12,2) DEFAULT 0.00 COMMENT '输入成本',
    output_cost DECIMAL(12,2) DEFAULT 0.00 COMMENT '输出成本',

    -- 性能指标
    avg_duration_ms INT COMMENT '平均耗时',
    p95_duration_ms INT COMMENT 'P95耗时',
    error_rate DECIMAL(5,2) COMMENT '错误率',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_dimension_date (dimension_type, dimension_id, date),
    INDEX idx_tenant_date (tenant_id, date),
    INDEX idx_dimension_type (dimension_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成本分析表';
```

---

## 4. API设计

### 4.1 Agent 监控 API

#### 4.1.1 获取执行列表

```http
GET /api/v1/tenants/{tenant_id}/agent-monitor/executions

Query Parameters:
  - bot_id: string (可选) - Bot ID
  - user_id: string (可选) - 用户ID
  - status: string (可选) - 状态筛选
  - start_date: date (可选) - 开始日期
  - end_date: date (可选) - 结束日期
  - page: int (可选)
  - page_size: int (可选)

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 1000,
    "items": [
      {
        "id": "exec-001",
        "bot_id": "bot-123",
        "bot_name": "客服助手",
        "user_id": 101,
        "user_name": "张三",
        "conversation_id": "conv-456",
        "input_text": "怎么重置密码？",
        "output_text": "您可以按以下步骤重置密码...",
        "status": "success",
        "total_tokens": 250,
        "duration_ms": 1500,
        "cost": 0.005,
        "start_time": "2025-01-03T10:00:00Z"
      }
    ]
  }
}
```

#### 4.1.2 获取执行详情

```http
GET /api/v1/tenants/{tenant_id}/agent-monitor/executions/{execution_id}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "exec-001",
    "bot_id": "bot-123",
    "input_text": "怎么重置密码？",
    "output_text": "您可以按以下步骤重置密码...",
    "status": "success",
    "total_tokens": 250,
    "duration_ms": 1500,
    "cost": 0.005,
    "model_name": "gpt-4",
    "start_time": "2025-01-03T10:00:00Z",
    "end_time": "2025-01-03T10:00:01.5Z",

    // 追踪树
    "trace_tree": {
      "id": "step-1",
      "name": "意图识别",
      "type": "llm_call",
      "input": "用户输入: 怎么重置密码？",
      "output": "意图: account_management",
      "duration_ms": 200,
      "tokens": 50,
      "children": [
        {
          "id": "step-2",
          "name": "知识库检索",
          "type": "tool_call",
          "input": "关键词: 重置 密码",
          "output": "检索到3篇文档",
          "duration_ms": 300,
          "children": []
        },
        {
          "id": "step-3",
          "name": "生成回复",
          "type": "llm_call",
          "input": "基于上下文生成回复",
          "output": "您可以按以下步骤...",
          "duration_ms": 1000,
          "tokens": 200,
          "children": []
        }
      ]
    }
  }
}
```

#### 4.1.3 获取性能指标

```http
GET /api/v1/tenants/{tenant_id}/agent-monitor/metrics

Query Parameters:
  - bot_id: string (可选)
  - start_date: date (可选)
  - end_date: date (可选)
  - granularity: hour | day | week (默认day)

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "metrics": [
      {
        "timestamp": "2025-01-03T00:00:00Z",
        "total_requests": 500,
        "success_requests": 490,
        "failed_requests": 10,
        "success_rate": 0.98,
        "avg_duration_ms": 1200,
        "p95_duration_ms": 2000,
        "p99_duration_ms": 3500,
        "total_tokens": 125000,
        "avg_tokens_per_request": 250,
        "total_cost": 25.00
      }
    ]
  }
}
```

### 4.2 评估 API

#### 4.2.1 创建评估任务

```http
POST /api/v1/tenants/{tenant_id}/agent-monitor/evaluations

Request Body:
{
  "name": "客服助手准确性评估",
  "bot_id": "bot-123",
  "evaluation_type": "accuracy",  // accuracy | relevance | safety | hallucination

  // 评估配置
  "dataset_id": "ds-001",  // 测试数据集ID
  "evaluation_methods": ["exact_match", "similarity", "llm_as_judge"],

  // 基准答案
  "ground_truth": [
    {
      "input": "怎么重置密码？",
      "expected_output": "您可以按以下步骤重置密码..."
    }
  ]
}

Response 200:
{
  "code": 0,
  "message": "评估任务创建成功",
  "data": {
    "evaluation_id": "eval-001",
    "status": "running"
  }
}
```

#### 4.2.2 获取评估结果

```http
GET /api/v1/tenants/{tenant_id}/agent-monitor/evaluations/{evaluation_id}

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "evaluation_id": "eval-001",
    "name": "客服助手准确性评估",
    "status": "completed",

    // 综合评分
    "overall_score": 85.5,
    "passed": true,

    // 分项评分
    "scores": {
      "exact_match": 70.0,
      "similarity": 88.5,
      "llm_as_judge": 88.0
    },

    // 详细结果
    "detailed_results": [
      {
        "input": "怎么重置密码？",
        "actual_output": "您可以按以下步骤...",
        "expected_output": "您可以按以下步骤...",
        "score": 90.0,
        "passed": true
      }
    ],

    // 问题诊断
    "issues": [
      {
        "type": "low_accuracy",
        "description": "在处理复杂问题时准确率偏低",
        "affected_cases": 15,
        "suggestion": "建议优化知识库检索策略"
      }
    ]
  }
}
```

### 4.3 提示词版本管理 API

#### 4.3.1 创建提示词版本

```http
POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/prompt-versions

Request Body:
{
  "version": "v1.1.0",
  "system_prompt": "你是一个专业的客服助手...",
  "user_prompt_template": "用户问题: {user_input}\n上下文: {context}",
  "parameters": {
    "type": "object",
    "properties": {
      "user_input": {"type": "string"},
      "context": {"type": "string"}
    }
  },
  "parent_version_id": "v1.0.0",
  "change_description": "优化了回复的专业性"
}

Response 200:
{
  "code": 0,
  "message": "版本创建成功",
  "data": {
    "id": "v1.1.0",
    "version": "v1.1.0",
    "is_active": false,
    "created_at": "2025-01-03T10:00:00Z"
  }
}
```

#### 4.3.2 激活提示词版本

```http
POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/prompt-versions/{version_id}/activate

Response 200:
{
  "code": 0,
  "message": "版本已激活"
}
```

#### 4.3.3 版本对比

```http
GET /api/v1/tenants/{tenant_id}/bots/{bot_id}/prompt-versions/compare

Query Parameters:
  - version_id_1: string - 版本1 ID
  - version_id_2: string - 版本2 ID

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "version_1": {
      "id": "v1.0.0",
      "system_prompt": "你是一个客服助手..."
    },
    "version_2": {
      "id": "v1.1.0",
      "system_prompt": "你是一个专业的客服助手..."
    },
    "diff": [
      {
        "type": "added",
        "old_line": null,
        "new_line": "你是一个专业的客服助手",
        "line_number": 1
      }
    ]
  }
}
```

### 4.4 成本分析 API

#### 4.4.1 获取成本统计

```http
GET /api/v1/tenants/{tenant_id}/agent-monitor/cost-analytics

Query Parameters:
  - dimension: total | by_bot | by_user | by_department (默认total)
  - dimension_id: string (可选) - 维度ID
  - start_date: date (可选)
  - end_date: date (可选)
  - granularity: day | week | month (默认day)

Response 200:
{
  "code": 0,
  "message": "success",
  "data": {
    "dimension": "by_bot",
    "total_cost": 1250.50,
    "total_tokens": 62525000,
    "total_requests": 250100,

    // 成本明细
    "breakdown": [
      {
        "dimension_id": "bot-001",
        "dimension_name": "客服助手",
        "total_cost": 500.20,
        "total_tokens": 25010000,
        "total_requests": 100040,
        "avg_cost_per_request": 0.005,
        "avg_tokens_per_request": 250
      },
      {
        "dimension_id": "bot-002",
        "dimension_name": "销售助手",
        "total_cost": 750.30,
        "total_tokens": 37515000,
        "total_requests": 150060,
        "avg_cost_per_request": 0.005,
        "avg_tokens_per_request": 250
      }
    ]
  }
}
```

---

## 5. 核心代码实现

### 5.1 监控服务

```go
// service/agent_monitor_service.go
package service

type AgentMonitorService struct {
    db             *gorm.DB
    logger         *zap.Logger
}

// RecordExecution 记录Agent执行
func (s *AgentMonitorService) RecordExecution(
    ctx context.Context,
    req *RecordExecutionRequest,
) (*AgentExecution, error) {

    execution := &AgentExecution{
        ID:            generateID(),
        TenantID:      req.TenantID,
        UserID:        req.UserID,
        BotID:         req.BotID,
        ConversationID: req.ConversationID,
        MessageID:     req.MessageID,
        InputText:     req.InputText,
        Status:        "running",
        StartTime:     time.Now(),
        ModelProvider: req.ModelProvider,
        ModelName:     req.ModelName,
    }

    if err := s.db.WithContext(ctx).Create(execution).Error; err != nil {
        return nil, err
    }

    return execution, nil
}

// CompleteExecution 完成Agent执行
func (s *AgentMonitorService) CompleteExecution(
    ctx context.Context,
    executionID string,
    result *ExecutionResult,
) error {

    updates := map[string]interface{}{
        "output_text":   result.OutputText,
        "status":        result.Status,
        "end_time":      time.Now(),
        "duration_ms":   time.Since(result.StartTime).Milliseconds(),
        "input_tokens":  result.InputTokens,
        "output_tokens": result.OutputTokens,
        "total_tokens":  result.InputTokens + result.OutputTokens,
        "cost":          result.Cost,
    }

    if result.Error != nil {
        updates["error_type"] = result.Error.Type
        updates["error_message"] = result.Error.Message
        updates["error_stack"] = result.Error.Stack
    }

    return s.db.WithContext(ctx).
        Model(&AgentExecution{}).
        Where("id = ?", executionID).
        Updates(updates).Error
}

// RecordStep 记录执行步骤
func (s *AgentMonitorService) RecordStep(
    ctx context.Context,
    executionID string,
    step *ExecutionStep,
) error {

    step.ExecutionID = executionID
    step.Status = "running"
    step.StartTime = time.Now()

    if err := s.db.WithContext(ctx).Create(step).Error; err != nil {
        return err
    }

    return nil
}
```

### 5.2 评估服务

```go
// service/evaluation_service.go
package service

type EvaluationService struct {
    db     *gorm.DB
    llm    *llm.LLMClient
}

// Evaluate 评估Agent输出
func (s *EvaluationService) Evaluate(
    ctx context.Context,
    req *EvaluationRequest,
) (*EvaluationResult, error) {

    // 1. 根据评估方法执行评估
    scores := make(map[string]float64)

    for _, method := range req.EvaluationMethods {
        switch method {
        case "exact_match":
            score := s.exactMatch(ctx, req)
            scores[method] = score

        case "similarity":
            score := s.similarity(ctx, req)
            scores[method] = score

        case "llm_as_judge":
            score := s.llmAsJudge(ctx, req)
            scores[method] = score
        }
    }

    // 2. 计算综合评分
    overallScore := s.calculateOverallScore(scores)

    // 3. 生成问题诊断
    issues := s.diagnoseIssues(ctx, req, scores)

    // 4. 生成优化建议
    suggestions := s.generateSuggestions(ctx, issues)

    return &EvaluationResult{
        EvaluationID:  generateID(),
        OverallScore:  overallScore,
        Scores:       scores,
        Issues:       issues,
        Suggestions:  suggestions,
        Passed:       overallScore >= 80.0,
    }, nil
}

// exact_match 精确匹配评估
func (s *EvaluationService) exactMatch(
    ctx context.Context,
    req *EvaluationRequest,
) float64 {

    totalCases := len(req.GroundTruth)
    passedCases := 0

    for _, groundTruth := range req.GroundTruth {
        actualOutput := req.GetActualOutput(groundTruth.Input)

        if s.normalizeText(actualOutput) == s.normalizeText(groundTruth.ExpectedOutput) {
            passedCases++
        }
    }

    return float64(passedCases) / float64(totalCases) * 100
}

// similarity 相似度评估
func (s *EvaluationService) similarity(
    ctx context.Context,
    req *EvaluationRequest,
) float64 {

    // 使用LLM计算相似度
    totalSimilarity := 0.0

    for _, groundTruth := range req.GroundTruth {
        actualOutput := req.GetActualOutput(groundTruth.Input)

        prompt := fmt.Sprintf(`
            请评估以下两段文本的相似度（0-100分）：

            文本1（期望输出）: %s
            文本2（实际输出）: %s

            仅输出数字分数，不要其他内容。
        `, groundTruth.ExpectedOutput, actualOutput)

        similarity, _ := s.llm.Complete(ctx, prompt)
        score := parseScore(similarity)

        totalSimilarity += score
    }

    return totalSimilarity / float64(len(req.GroundTruth))
}
```

### 5.3 成本分析服务

```go
// service/cost_analysis_service.go
package service

type CostAnalysisService struct {
    db *gorm.DB
}

// CalculateCost 计算成本
func (s *CostAnalysisService) CalculateCost(
    ctx context.Context,
    execution *AgentExecution,
) decimal.Decimal {

    // 1. 获取模型定价
    pricing, err := s.getModelPricing(ctx, execution.ModelProvider, execution.ModelName)
    if err != nil {
        s.logger.Error("failed to get model pricing", zap.Error(err))
        return decimal.Zero
    }

    // 2. 计算成本
    inputCost := decimal.NewFromInt(int64(execution.InputTokens)).
        Mul(decimal.NewFromInt(pricing.InputPricePer1K)).
        Div(decimal.NewFromInt(1000))

    outputCost := decimal.NewFromInt(int64(execution.OutputTokens)).
        Mul(decimal.NewFromInt(pricing.OutputPricePer1K)).
        Div(decimal.NewFromInt(1000))

    totalCost := inputCost.Add(outputCost)

    return totalCost
}

// GenerateCostReport 生成成本报告
func (s *CostAnalysisService) GenerateCostReport(
    ctx context.Context,
    tenantID string,
    params *CostReportParams,
) (*CostReport, error) {

    var results []CostAnalyticsResult

    query := s.db.WithContext(ctx).
        Table("cost_analytics").
        Where("tenant_id = ?", tenantID).
        Where("date >= ? AND date <= ?", params.StartDate, params.EndDate)

    if params.DimensionType != "" {
        query = query.Where("dimension_type = ?", params.DimensionType)
    }

    if err := query.Find(&results).Error; err != nil {
        return nil, err
    }

    // 汇总统计
    totalCost := decimal.Zero
    totalTokens := int64(0)
    totalRequests := 0

    for _, result := range results {
        totalCost = totalCost.Add(result.TotalCost)
        totalTokens += result.TotalTokens
        totalRequests += int(result.TotalRequests)
    }

    // 成本优化建议
    suggestions := s.generateOptimizationSuggestions(ctx, results)

    return &CostReport{
        TenantID:       tenantID,
        DimensionType:  params.DimensionType,
        TotalCost:      totalCost,
        TotalTokens:    totalTokens,
        TotalRequests:  totalRequests,
        AvgCostPerRequest: totalCost.Div(decimal.NewFromInt(int64(totalRequests))),
        Breakdown:      results,
        Suggestions:    suggestions,
    }, nil
}
```

---

## 6. 前端设计

### 6.1 监控仪表板

```typescript
// src/pages/monitor/MonitorDashboard.tsx
const MonitorDashboard: React.FC = () => {
  const { metrics, loading } = useMetrics();

  return (
    <div className="monitor-dashboard">
      {/* 关键指标卡片 */}
      <Row gutter={16}>
        <Col span={6}>
          <Card title="总请求数">
            <Statistic value={metrics.totalRequests} />
          </Card>
        </Col>
        <Col span={6}>
          <Card title="成功率">
            <Statistic
              value={metrics.successRate * 100}
              suffix="%"
              valueStyle={{ color: metrics.successRate > 0.95 ? 'green' : 'red' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card title="平均响应时间">
            <Statistic value={metrics.avgDurationMs} suffix="ms" />
          </Card>
        </Col>
        <Col span={6}>
          <Card title="总成本">
            <Statistic value={metrics.totalCost} prefix="¥" />
          </Card>
        </Col>
      </Row>

      {/* 时序图表 */}
      <Card title="性能趋势" style={{ marginTop: 16 }}>
        <LineChart
          data={metrics.timeSeries}
          xField="timestamp"
          yFields={['total_requests', 'success_rate', 'avg_duration_ms']}
          height={300}
        />
      </Card>

      {/* 错误追踪 */}
      <Card title="错误追踪" style={{ marginTop: 16 }}>
        <Table
          columns={[
            { title: '时间', dataIndex: 'timestamp' },
            { title: 'Bot', dataIndex: 'bot_name' },
            { title: '错误类型', dataIndex: 'error_type' },
            { title: '错误信息', dataIndex: 'error_message' },
            { title: '操作', render: (text, record) => (
              <Button onClick={() => handleErrorDetail(record)}>查看详情</Button>
            )},
          ]}
          dataSource={metrics.errors}
        />
      </Card>
    </div>
  );
};
```

### 6.2 执行追踪树

```typescript
// src/components/monitor/TraceTree.tsx
const TraceTree: React.FC<{ execution: AgentExecution }> = ({ execution }) => {
  return (
    <div className="trace-tree">
      <Tree
        treeData={[convertToTreeNode(execution)]}
        defaultExpandAll
        renderTitle={(nodeData) => (
          <div className="trace-node">
            <Tag color={getStatusColor(nodeData.status)}>
              {nodeData.step_name}
            </Tag>
            <span className="step-type">{nodeData.step_type}</span>
            {nodeData.tokens && (
              <span className="tokens">{nodeData.tokens} tokens</span>
            )}
            <span className="duration">{nodeData.duration_ms}ms</span>
            {nodeData.status === 'failed' && (
              <Button size="small" onClick={() => showErrorDetail(nodeData)}>
                查看错误
              </Button>
            )}
          </div>
        )}
      />
    </div>
  );
};
```

---

## 7. 总结

### 7.1 实施策略总结

| 实施项 | 实施方式 | 工作量 |
|-------|---------|--------|
| **监控框架** | 💻 独立编码 | 40% |
| **评估引擎** | 📊 配置驱动 | 30% |
| **成本分析** | 📊 配置驱动 | 20% |
| **前端页面** | 💻 独立编码 | 10% |

**总计**: 40% 编码 + 60% 配置

### 7.2 核心优势

- ✅ **全流程可观测性** - 可视化追踪树，清晰展现每一步
- ✅ **实时监控** - 性能、错误、成本实时追踪
- ✅ **质量评估** - 12种评估方法，自动化评估
- ✅ **提示词工程** - 版本管理、A/B测试、效果对比
- ✅ **成本优化** - 实时成本追踪，优化建议

### 7.3 对标 LangSmith

| 功能 | LangSmith | ZKER | 说明 |
|------|----------|------|------|
| **调用链追踪** | ✅ | ✅ | 可视化追踪树 |
| **性能监控** | ✅ | ✅ | Token消耗、响应时间 |
| **质量评估** | ✅ | ✅ | 12种评估方法 |
| **提示词版本管理** | ✅ | ✅ | Git-like版本控制 |
| **A/B测试** | ✅ | ✅ | 多版本并行测试 |
| **成本分析** | ✅ | ✅ | 实时成本追踪 |
| **开源** | ❌ | ✅ | ZKER开源 |

---

**文档结束**
