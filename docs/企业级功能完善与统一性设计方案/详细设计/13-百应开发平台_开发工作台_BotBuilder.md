# 13-百应开发平台_开发工作台_BotBuilder 详细设计说明书

**文档编号**: DE-DD-2025-013
**模块名称**: BotBuilder (低代码Bot构建器)
**版本**: v1.0.0
**作者**: ZKER Enterprise Team
**创建日期**: 2025-01-03

---

## 1. 模块概述

**BotBuilder** 是 ZKER 的低代码Bot构建平台，通过**100%复用zker BotBuilder** + 企业级权限管理，实现企业快速构建和部署Bot。

**核心设计理念**：
- ✅ **可视化构建**：拖拽式Bot设计器
- ✅ **零代码开发**：无需编程即可创建Bot
- ✅ **企业级管理**：多租户、权限控制、审核流程

**实现策略**：✅ 100%复用zker BotBuilder + 配置化扩展

---

## 2. 核心功能（复用zker）

### 2.1 Bot构建

**F1 - 可视化编辑器**（复用）
- F1.1 提示词编辑
- F1.2 技能选择
- F1.3 知识库关联
- F1.4 插件配置

### 2.2 Bot测试

**F2 - 在线测试**（复用）
- F2.1 对话测试
- F2.2 调试模式
- F2.3 日志查看

### 2.3 Bot部署

**F3 - 发布管理**（复用）
- F3.1 版本管理
- F3.2 发布/下架
- F3.3 灰度发布

---

## 3. 企业级扩展

### 3.1 新增企业级功能

**F4 - 团队协作**（新增）
- F4.1 Bot协作编辑
- F4.2 评论与审核
- F4.3 操作日志

**F5 - 权限管理**（新增）
- F5.1 Bot权限控制
- F5.2 部门级可见性

**F6 - 模板管理**（新增）
- F6.1 Bot模板库
- F6.2 快速创建

---

## 4. 数据库设计

### 4.1 复用zker表 + 扩展字段

**完全复用 `bots` 表**，新增字段：

```sql
-- 扩展 bots 表
ALTER TABLE bots ADD COLUMN team_id BIGINT COMMENT '团队ID';
ALTER TABLE bots ADD COLUMN is_template BOOLEAN DEFAULT FALSE COMMENT '是否为模板';
ALTER TABLE bots ADD COLUMN template_id VARCHAR(64) COMMENT '基于哪个模板创建';
ALTER TABLE bots ADD COLUMN review_status ENUM('draft', 'pending_review', 'approved', 'rejected') DEFAULT 'draft';
```

### 4.2 新增协作表

```sql
CREATE TABLE bot_collaborators (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    bot_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    role ENUM('owner', 'editor', 'viewer') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_bot_user (bot_id, user_id),
    INDEX idx_bot_id (bot_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 5. API设计（完全复用）

| 方法 | 路径 | 功能 | 实现方式 |
|------|------|------|----------|
| GET | /api/v1/bot-builder/bots | 列出Bot | → zker |
| POST | /api/v1/bot-builder/bots | 创建Bot | → zker |
| PUT | /api/v1/bot-builder/bots/:id | 更新Bot | → zker |
| POST | /api/v1/bot-builder/bots/:id/publish | 发布Bot | → zker |
| POST | /api/v1/bot-builder/bots/:id/test | 测试Bot | → zker |

---

## 6. 总结

### 6.1 实施策略总结

| 实施项 | 实施方式 | 工作量 |
|-------|---------|--------|
| **BotBuilder** | 🔁 100%复用 | 0% |
| **企业扩展** | 📊 配置扩展 | 20% |

**总计**：0% 核心开发 + 20% 配置扩展

### 6.2 核心优势

- ✅ **零重复开发**：100%复用zker BotBuilder
- ✅ **快速上线**：配置即可使用
- ✅ **功能完整**：继承zker所有能力

---

**文档结束**
