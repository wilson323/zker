# API接口文档：人机协同引擎模块

**模块名称**: 人机协同引擎 (HumanInLoop)
**设计文档**: 15-五大AI引擎核心_人机协同引擎.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 协同任务API](#2-协同任务api)
- [3. 规则配置API](#3-规则配置api)
- [4. 审核工作台API](#4-审核工作台api)
- [5. 反馈学习API](#5-反馈学习api)
- [6. 数据模型](#6-数据模型)
- [7. 错误码定义](#7-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

人机协同引擎是五大AI引擎的核心，通过AI决策+人工审核的协同机制，在关键决策点引入人工干预：

- ✅ **智能触发** - 自动识别需要人工干预的场景
- ✅ **协同工作流** - 人工审核、修正、反馈
- ✅ **持续学习** - 人工反馈用于模型微调
- ✅ **规则配置** - 灵活配置触发规则
- ✅ **任务分配** - 自动分配审核人

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 消息队列: NSQ（异步通知）

### 1.3 触发场景

| 触发类型 | 说明 | 示例 |
|---------|------|------|
| **置信度低** | AI置信度低于阈值 | 置信度<0.7 |
| **敏感内容** | 包含敏感关键词 | 涉及薪资、合同 |
| **复杂业务** | 特定业务类别 | HR、法务、财务 |
| **异常行为** | 检测到异常模式 | 异常登录、频繁操作 |

---

## 2. 协同任务API

### 2.1 创建协同任务

**接口地址**: `POST /api/v1/human-in-loop/tasks`

**请求参数**:
```json
{
  "task_type": "content_review",
  "conversation_id": "conv-001",
  "message_id": "msg-001",
  "ai_output": "AI生成的回复内容...",
  "confidence": 0.65,
  "trigger_reason": "置信度低于阈值0.7",
  "priority": "high",
  "assign_to": 1001
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "协同任务创建成功",
  "data": {
    "id": "task-001",
    "status": "pending",
    "createdAt": "2025-01-03T11:00:00Z"
  }
}
```

### 2.2 获取任务列表

**接口地址**: `GET /api/v1/human-in-loop/tasks`

**查询参数**:
- page: 页码
- page_size: 每页数量
- status: 状态过滤
- priority: 优先级过滤

### 2.3 审核任务

**接口地址**: `POST /api/v1/human-in-loop/tasks/{task_id}/review`

**请求参数**:
```json
{
  "decision": "approved",
  "human_decision": "同意AI的回复",
  "decision_reason": "回复内容符合要求"
}
```

---

## 3. 规则配置API

### 3.1 获取规则列表

**接口地址**: `GET /api/v1/human-in-loop/rules`

### 3.2 创建规则

**接口地址**: `POST /api/v1/human-in-loop/rules`

---

## 4. 审核工作台API

### 4.1 获取审核统计

**接口地址**: `GET /api/v1/human-in-loop/dashboard/stats`

### 4.2 获取我的任务

**接口地址**: `GET /api/v1/human-in-loop/my-tasks`

---

## 5. 反馈学习API

### 5.1 导出审核数据

**接口地址**: `GET /api/v1/human-in-loop/feedback/export`

---

## 6. 数据模型

### 6.1 HumanLoopTask
```typescript
interface HumanLoopTask {
  id: string;
  task_type: 'content_review' | 'decision_approval' | 'error_correction';
  ai_output: string;
  confidence?: number;
  status: 'pending' | 'in_review' | 'approved' | 'rejected' | 'modified';
  priority: 'low' | 'medium' | 'high' | 'urgent';
}
```

---

## 7. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 70001 | 404 | 任务不存在 |
| 70002 | 400 | 任务状态不允许此操作 |
| 70101 | 400 | 规则配置无效 |

---

**文档结束**
