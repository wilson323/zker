# ZKER 文档门户

**版本**: v3.0.0 | **更新**: 2025-01-03 | **状态**: 企业级功能完善阶段

---

## 📚 文档导航

### 快速入口

| 入口 | 说明 | 适用对象 |
|------|------|---------|
| **[项目概览](../01-OVERVIEW/README.md)** | 项目定位、技术栈、架构概览 | 所有人 |
| **[开发规范](../02-SPECS/README.md)** | 命名、API、数据库、测试规范 | 开发人员 |
| **[设计文档](../03-DESIGN/README.md)** | 多租户、RBAC、路由设计 | 架构师/开发 |
| **[实施指南](../04-IMPLEMENTATION/README.md)** | 开发指南、部署流程 | 开发/运维 |
| **[运维手册](../05-OPERATIONS/README.md)** | 监控、故障排查、回滚 | 运维人员 |

---

## 🎯 按角色查找文档

### 架构师

必读文档：
1. [03-DESIGN/multi-tenant/](../03-DESIGN/multi-tenant/) - 多租户架构设计
2. [03-DESIGN/rbac/](../03-DESIGN/rbac/) - RBAC 权限设计
3. [03-DESIGN/routing/](../03-DESIGN/routing/) - 智能路由设计
4. [01-OVERVIEW/architecture.md](../01-OVERVIEW/architecture.md) - 系统架构

### 后端开发

必读文档：
1. [02-SPECS/backend-dev-guide.md](../02-SPECS/backend-dev-guide.md) - 后端开发规范
2. [02-SPECS/naming-conventions.md](../02-SPECS/naming-conventions.md) - 命名规范
3. [02-SPECS/database-design.md](../02-SPECS/database-design.md) - 数据库设计
4. [02-SPECS/error-handling.md](../02-SPECS/error-handling.md) - 错误处理

### 前端开发

必读文档：
1. [02-SPECS/frontend-dev-guide.md](../02-SPECS/frontend-dev-guide.md) - 前端开发规范
2. [02-SPECS/component-design.md](../02-SPECS/component-design.md) - 组件设计
3. [02-SPECS/performance-optimization.md](../02-SPECS/performance-optimization.md) - 性能优化

### DevOps

必读文档：
1. [04-IMPLEMENTATION/deployment.md](../04-IMPLEMENTATION/deployment.md) - 部署指南
2. [05-OPERATIONS/monitoring.md](../05-OPERATIONS/monitoring.md) - 监控配置
3. [05-OPERATIONS/troubleshooting.md](../05-OPERATIONS/troubleshooting.md) - 故障排查
4. [05-OPERATIONS/rollback.md](../05-OPERATIONS/rollback.md) - 回滚方案

---

## 📂 文档结构说明

```
docs/
├── 00-META/              # 元文档（本目录）
│   ├── README.md         # 文档门户
│   ├── VERSION.md        # 版本索引
│   └── CHANGELOG.md      # 变更日志
│
├── 01-OVERVIEW/          # 概览文档
│   ├── README.md         # 项目概览
│   ├── architecture.md   # 架构设计
│   └── tech-stack.md     # 技术栈
│
├── 02-SPECS/             # 规范文档
│   ├── README.md         # 规范索引
│   ├── naming-conventions.md
│   ├── backend-dev-guide.md
│   ├── frontend-dev-guide.md
│   ├── api-design.md
│   ├── database-design.md
│   ├── error-handling.md
│   └── testing-guide.md
│
├── 03-DESIGN/            # 设计文档
│   ├── README.md         # 设计索引
│   ├── multi-tenant/     # 多租户设计
│   ├── rbac/             # RBAC 设计
│   └── routing/          # 智能路由设计
│
├── 04-IMPLEMENTATION/    # 实施指南
│   ├── README.md
│   ├── backend-dev-guide.md
│   ├── frontend-dev-guide.md
│   └── deployment.md
│
├── 05-OPERATIONS/        # 运维文档
│   ├── README.md
│   ├── monitoring.md
│   ├── troubleshooting.md
│   └── rollback.md
│
└── archive/              # 归档（旧版本）
    └── enterprise-level-v1/
```

---

## 🔄 文档版本管理

### 当前版本: v3.0.0

**发布日期**: 2025-01-03

**主要变更**:
- ✅ 重构文档目录结构
- ✅ 精简 CLAUDE.md（714 行 → 247 行）
- ✅ 统一开发规范文档
- ✅ 增强企业级 Skills

### 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| v3.0.0 | 2025-01-03 | 文档结构重构，规范统一 |
| v2.0.0 | 2024-12-29 | 企业级功能完善 |
| v1.0.0 | 2024-12-01 | 初始版本 |

详见: [VERSION.md](VERSION.md)

---

## 📖 文档贡献指南

### 文档编写规范

1. **标题规范**:
   - 一级标题: `# 文档名称`
   - 二级标题: `## 主要章节`
   - 三级标题: `### 子章节`

2. **代码规范**:
   - 使用语法高亮: ```go ```typescript ```sql
   - 注释清晰: `// Good` `// Bad`
   - 完整示例: 可直接运行的代码

3. **链接规范**:
   - 相对链接: `[文档名](../path/to/file.md)`
   - 绝对链接: `[文档名](/docs/path/to/file.md)`
   - 外部链接: 确保可访问

### 文档审核流程

1. 编写 → 2. 自查 → 3. PR → 4. 审查 → 5. 合并

**审核清单**:
- [ ] 内容准确，无错误
- [ ] 格式统一，符合规范
- [ ] 链接有效，无死链
- [ ] 示例代码可运行

---

## 🔍 文档搜索技巧

### 按关键词搜索

```bash
# 搜索多租户相关文档
grep -r "多租户\|multi-tenant\|tenant" docs/

# 搜索 RBAC 相关文档
grep -r "RBAC\|权限\|permission" docs/

# 搜索错误码相关文档
grep -r "错误码\|errno\|error" docs/
```

### 按文件类型搜索

```bash
# 查找所有规范文档
find docs/02-SPECS/ -name "*.md"

# 查找所有设计文档
find docs/03-DESIGN/ -name "*.md"

# 查找所有运维文档
find docs/05-OPERATIONS/ -name "*.md"
```

---

## 💡 常见问题

### Q: 找不到需要的文档？

A: 请按以下步骤查找：
1. 检查本文档的"按角色查找文档"部分
2. 使用关键词搜索
3. 查看归档目录（docs/archive/）
4. 提交 Issue 申请添加新文档

### Q: 文档内容过时？

A: 请：
1. 查看 [VERSION.md](VERSION.md) 确认当前版本
2. 对比代码实现，更新文档
3. 提交 PR 修复

### Q: 如何贡献文档？

A: 请：
1. 阅读 [贡献指南](#-文档贡献指南)
2. Fork 项目，创建分支
3. 编写文档，遵循规范
4. 提交 PR，等待审核

---

## 📞 联系方式

- **文档问题**: [GitHub Issues](https://github.com/coze-dev/coze-studio/issues)
- **贡献指南**: [CONTRIBUTING.md](../../CONTRIBUTING.md)
- **团队联系**: zker-team@example.com

---

**🎯 目标**: 打造清晰、准确、易用的企业级文档体系！
