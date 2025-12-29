# 05-ZKER前台_待办_TaskCenter 详细设计说明书

**文档编号**: DE-DD-2025-005
**模块名称**: 待办_TaskCenter (TaskCenter)
**版本**: v1.0.0
**作者**: ZKER Enterprise Team
**创建日期**: 2025-01-03

---

## 1. 模块概述

**待办_TaskCenter** 是 ZKER 企业级 SaaS 平台的前台模块，通过**通用任务引擎 + 任务类型配置**的方式，提供待办事项管理能力。

**核心设计理念**：
- ✅ **通用引擎**：20% 通用任务处理框架
- ✅ **类型配置**：80% 任务类型通过数据库配置
- ✅ **灵活扩展**：支持自定义任务类型和字段

**实现策略**：✅ 20% 编码（通用引擎） + 80% 配置（任务类型）

---

## 2. 核心功能

### 2.1 用户端功能

**F1 - 待办管理**
- F1.1 创建待办（标题、描述、截止日期、优先级）
- F1.2 编辑待办
- F1.3 完成待办
- F1.4 删除待办
- F1.5 待办分类（工作、个人、学习等）

**F2 - 任务类型**
- F2.1 预置任务类型（审批、报告、会议等）
- F2.2 企业自定义任务类型

---

## 3. 数据库设计

### 3.1 核心表结构

#### 3.1.1 待办事项表 (tasks)

```sql
CREATE TABLE tasks (
    id VARCHAR(64) PRIMARY KEY COMMENT '任务ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    type_id BIGINT COMMENT '任务类型ID',

    title VARCHAR(200) NOT NULL COMMENT '任务标题',
    description TEXT COMMENT '任务描述',
    status ENUM('pending', 'in_progress', 'completed', 'cancelled') DEFAULT 'pending',
    priority ENUM('low', 'medium', 'high', 'urgent') DEFAULT 'medium',

    due_date DATETIME COMMENT '截止日期',
    completed_at DATETIME COMMENT '完成时间',

    -- 扩展字段（JSON）
    custom_fields JSON COMMENT '自定义字段',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,

    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_status (status),
    INDEX idx_due_date (due_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='待办事项表';
```

#### 3.1.2 任务类型配置表 (task_types)

**核心设计**：通过数据库配置不同任务类型的字段和规则

```sql
CREATE TABLE task_types (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '类型ID',
    tenant_id VARCHAR(64) COMMENT '租户ID (NULL表示平台预置)',
    name VARCHAR(50) NOT NULL COMMENT '类型名称',
    description VARCHAR(200) COMMENT '类型描述',
    icon VARCHAR(10) COMMENT '图标 emoji',
    color VARCHAR(20) COMMENT '主题颜色',

    -- 字段配置
    fields JSON NOT NULL COMMENT '字段定义',
    validation_rules JSON COMMENT '验证规则',

    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_tenant_name (tenant_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='任务类型配置表';
```

**字段配置 JSON 示例**（审批任务）：

```json
{
  "fields": [
    {
      "name": "applicant",
      "label": "申请人",
      "type": "text",
      "required": true
    },
    {
      "name": "approval_type",
      "label": "审批类型",
      "type": "select",
      "required": true,
      "options": [
        {"value": "leave", "label": "请假审批"},
        {"value": "expense", "label": "费用审批"},
        {"value": "purchase", "label": "采购审批"}
      ]
    },
    {
      "name": "amount",
      "label": "金额",
      "type": "number",
      "required": false
    }
  ]
}
```

---

## 4. API 设计

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | /api/v1/tasks | 创建待办 |
| GET | /api/v1/tasks | 列出待办 |
| PUT | /api/v1/tasks/:id | 更新待办 |
| DELETE | /api/v1/tasks/:id | 删除待办 |
| POST | /api/v1/tasks/:id/complete | 完成待办 |
| GET | /api/v1/task-types | 获取任务类型 |

---

## 5. 配置系统设计

### 5.1 预置任务类型

```sql
-- 平台预置任务类型
INSERT INTO task_types (tenant_id, name, description, icon, color, fields) VALUES
-- 审批类
(NULL, '请假审批', '员工请假申请审批', '📅', '#1890ff',
'{"fields": [{"name":"applicant","label":"申请人","type":"text","required":true},{"name":"days","label":"请假天数","type":"number","required":true}]}'),

-- 报告类
(NULL, '周报', '工作周报提交', '📝', '#52c41a',
'{"fields": [{"name":"week_start","label":"周开始日期","type":"date","required":true},{"name":"achievements","label":"本周成果","type":"textarea","required":true}]}'),

-- 会议类
(NULL, '会议准备', '会议相关准备事项', '📞', '#faad14',
'{"fields": [{"name":"meeting_title","label":"会议主题","type":"text","required":true},{"name":"attendees","label":"参会人员","type":"textarea","required":false}]}');
```

---

## 6. 核心代码实现

### 6.1 通用任务引擎

```go
package task

// TaskEngine 通用任务引擎
type TaskEngine struct {
    taskTypeRepo  repository.TaskTypeRepository
    taskRepo      repository.TaskRepository
    validator     *Validator
}

// CreateTask 创建任务（支持动态字段）
func (e *TaskEngine) CreateTask(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
    // 1. 加载任务类型
    taskType, err := e.taskTypeRepo.GetByID(ctx, req.TypeID)
    if err != nil {
        return nil, err
    }

    // 2. 验证自定义字段
    if err := e.validator.ValidateCustomFields(taskType.Fields, req.CustomFields); err != nil {
        return nil, err
    }

    // 3. 创建任务
    task := &Task{
        ID:          generateTaskID(),
        TenantID:    getTenantID(ctx),
        UserID:      getUserID(ctx),
        TypeID:      req.TypeID,
        Title:       req.Title,
        Description: req.Description,
        Status:      "pending",
        Priority:    req.Priority,
        DueDate:     req.DueDate,
        CustomFields: req.CustomFields,
    }

    if err := e.taskRepo.Create(ctx, task); err != nil {
        return nil, err
    }

    return task, nil
}
```

---

## 7. 前端设计

### 7.1 动态表单组件

```tsx
import React, { useEffect, useState } from 'react';
import { Form, Input, Select, DatePicker } from '@douyinfe/semi-ui';

export const DynamicTaskForm: React.FC<{ typeId: number }> = ({ typeId }) => {
  const [taskType, setTaskType] = useState<TaskType | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadTaskType(typeId);
  }, [typeId]);

  const loadTaskType = async (typeId: number) => {
    const resp = await taskAPI.getTaskType(typeId);
    setTaskType(resp.data);
  };

  const renderField = (field: FieldConfig) => {
    switch (field.type) {
      case 'text':
        return <Input placeholder={field.label} />;
      case 'textarea':
        return <Input.TextArea rows={4} placeholder={field.label} />;
      case 'select':
        return (
          <Select
            placeholder={field.label}
            optionList={field.options}
          />
        );
      case 'date':
        return <DatePicker type="date" />;
      case 'number':
        return <Input type="number" placeholder={field.label} />;
      default:
        return null;
    }
  };

  return (
    <Form form={form} onSubmit={handleSubmit}>
      <Form.Input field="title" label="任务标题" required />
      <Form.TextArea field="description" label="描述" rows={3} />

      {/* 动态字段 */}
      {taskType?.fields.map(field => (
        <Form.Input
          key={field.name}
          field={`customFields.${field.name}`}
          label={field.label}
          required={field.required}
        >
          {renderField(field)}
        </Form.Input>
      ))}

      <Button type="primary" htmlType="submit">
        创建任务
      </Button>
    </Form>
  );
};
```

---

## 8. 总结

### 8.1 实施策略总结

| 实施项 | 实施方式 | 工作量 |
|-------|---------|--------|
| **通用任务引擎** | 💻 独立编码 | 20% |
| **任务类型配置** | 📊 数据库驱动 | 80% |

**总计**：20% 编码 + 80% 配置

### 8.2 核心优势

- ✅ **灵活配置**：80% 通过数据库配置任务类型
- ✅ **扩展性强**：新增任务类型无需编码
- ✅ **通用引擎**：一次开发，长期复用

---

**文档结束**
