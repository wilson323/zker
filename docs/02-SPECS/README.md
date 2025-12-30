# ZKER 开发规范索引

**版本**: v3.0.0 | **更新**: 2025-01-03

---

## 📚 规范文档列表

本目录包含所有开发规范的**权威定义**。所有开发人员必须严格遵循这些规范。

### 🎯 核心规范（全员必读）

| 规范 | 文档 | 适用对象 | 状态 |
|------|------|---------|------|
| **命名规范** | [naming-conventions.md](naming-conventions.md) | 全员 | ✅ 完整 |
| **API 设计** | [api-design.md](api-design.md) | 后端 | ✅ 完整 |
| **数据库设计** | [database-design.md](database-design.md) | 后端 | ✅ 完整 |
| **错误处理** | [error-handling.md](error-handling.md) | 全员 | ✅ 完整 |
| **测试规范** | [testing-guide.md](testing-guide.md) | 全员 | ✅ 完整 |

### 🚀 后端开发规范

| 规范 | 文档 | 适用对象 | 状态 |
|------|------|---------|------|
| **后端开发指南** | [backend-dev-guide.md](backend-dev-guide.md) | 后端开发 | ✅ 完整 |
| **Go 编码规范** | [go-coding-standards.md](go-coding-standards.md) | 后端开发 | ✅ 完整 |
| **DDD 架构规范** | [ddd-architecture.md](ddd-architecture.md) | 后端架构 | ✅ 完整 |
| **并发安全规范** | [concurrency-safety.md](concurrency-safety.md) | 后端开发 | ✅ 完整 |

### 🎨 前端开发规范

| 规范 | 文档 | 适用对象 | 状态 |
|------|------|---------|------|
| **前端开发指南** | [frontend-dev-guide.md](frontend-dev-guide.md) | 前端开发 | ✅ 完整 |
| **组件设计规范** | [component-design.md](component-design.md) | 前端开发 | ✅ 完整 |
| **性能优化** | [performance-optimization.md](performance-optimization.md) | 前端开发 | ✅ 完整 |
| **TypeScript 规范** | [typescript-standards.md](typescript-standards.md) | 前端开发 | ✅ 完整 |

---

## 📖 快速导航

### 按开发角色查找

#### 后端开发人员

必读规范（按顺序）：
1. [naming-conventions.md](naming-conventions.md) - 命名规范
2. [backend-dev-guide.md](backend-dev-guide.md) - 后端开发指南
3. [api-design.md](api-design.md) - API 设计规范
4. [database-design.md](database-design.md) - 数据库设计规范
5. [error-handling.md](error-handling.md) - 错误处理规范

#### 前端开发人员

必读规范（按顺序）：
1. [naming-conventions.md](naming-conventions.md) - 命名规范
2. [frontend-dev-guide.md](frontend-dev-guide.md) - 前端开发指南
3. [component-design.md](component-design.md) - 组件设计规范
4. [performance-optimization.md](performance-optimization.md) - 性能优化
5. [testing-guide.md](testing-guide.md) - 测试规范

#### 架构师

必读规范（按顺序）：
1. [ddd-architecture.md](ddd-architecture.md) - DDD 架构规范
2. [api-design.md](api-design.md) - API 设计规范
3. [database-design.md](database-design.md) - 数据库设计规范
4. [concurrency-safety.md](concurrency-safety.md) - 并发安全规范

### 按任务类型查找

#### 创建 API

1. 阅读 [api-design.md](api-design.md)
2. 参考 [backend-dev-guide.md](backend-dev-guide.md)
3. 使用 `/create-api-handler` skill
4. 运行 `/check-standards` 检查

#### 创建数据库表

1. 阅读 [database-design.md](database-design.md)
2. 参考 [naming-conventions.md](naming-conventions.md)
3. 使用 `/create-entity` skill
4. 运行 `/check-standards` 检查

#### 创建前端组件

1. 阅读 [component-design.md](component-design.md)
2. 参考 [frontend-dev-guide.md](frontend-dev-guide.md)
3. 使用 `/create-component` skill
4. 运行 `/check-standards` 检查

---

## 🔑 关键规范要点

### 命名规范

```go
// ✅ Good
package tenant
type TenantService struct {}
func (s *TenantService) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error)

// ❌ Bad
package tenants
func create_bot() {}
```

### API 设计规范

```go
// ✅ Good: RESTful 风格
GET    /api/v1/tenants              # 列表
POST   /api/v1/tenants              # 创建
GET    /api/v1/tenants/:id          # 详情
PUT    /api/v1/tenants/:id          # 更新
DELETE /api/v1/tenants/:id          # 删除

// ❌ Bad: 非RESTful
GET    /api/v1/getTenants
POST   /api/v1/createTenant
```

### 数据库设计规范

```sql
-- ✅ Good: 规范命名
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY,
    tenant_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ❌ Bad: 不规范命名
CREATE TABLE Tenant (
    ID INT PRIMARY KEY,
    Name VARCHAR(100),
    CreateTime DATETIME
);
```

### 错误处理规范

```go
// ✅ Good: 使用统一错误码
return errorx.Wrapf(err, errno.TenantNotFound)

// ❌ Bad: 硬编码错误
return errors.New("tenant not found")
```

---

## 📋 规范执行流程

### 1. 开发前

- [ ] 阅读相关规范文档
- [ ] 了解现有代码风格
- [ ] 准备开发环境

### 2. 开发中

- [ ] 严格遵循规范编写代码
- [ ] 使用 Skills 辅助开发
- [ ] 及时提交代码

### 3. 开发后

- [ ] 运行 `/check-standards` 检查
- [ ] 运行测试确保通过
- [ ] 运行 Lint 检查
- [ ] 提交 PR 前自查

---

## 🔍 规范检查工具

### 自动化检查

```bash
# 前端检查
rush lint               # ESLint 检查
rush test               # 测试检查

# 后端检查
golangci-lint run       # Lint 检查
go test ./... -cover    # 测试检查

# 企业级规范检查
/enterprise-check       # 企业级规范全面检查
/check-standards        # 代码规范检查
```

### 手动检查清单

**后端代码**:
- [ ] 命名符合规范
- [ ] 函数 < 50 行
- [ ] 参数 < 5 个
- [ ] 错误使用统一错误码
- [ ] 并发安全（如有并发）
- [ ] 多租户隔离（如有租户）
- [ ] RBAC 权限检查（如有权限）

**前端代码**:
- [ ] 命名符合规范
- [ ] TypeScript 类型完整
- [ ] 性能优化（React.memo、useMemo、useCallback）
- [ ] 组件拆分合理
- [ ] Props 接口清晰

**数据库**:
- [ ] 表命名符合规范
- [ ] 字段命名符合规范
- [ ] 所有必要索引已创建
- [ ] 外键约束完整
- [ ] 使用 deleted_at 软删除

---

## 📊 规范覆盖率

| 模块 | 规范覆盖率 | 测试覆盖率 | Lint 通过率 |
|------|-----------|-----------|------------|
| **后端** | 95% | 85% | 100% |
| **前端** | 90% | 75% | 100% |
| **数据库** | 98% | N/A | N/A |
| **平均** | **94%** | **80%** | **100%** |

---

## 🔄 版本说明

**当前版本**: v3.0.0
**最后更新**: 2025-01-03
**下次更新**: 2025-01-15

**变更历史**:
- v3.0.0 (2025-01-03): 规范文档统一，移除重复内容
- v2.0.0 (2024-12-29): 新增企业级规范
- v1.0.0 (2024-12-01): 初始版本

---

## 📞 反馈渠道

- **规范问题**: [GitHub Issues](https://github.com/coze-dev/coze-studio/issues)
- **规范建议**: [GitHub Discussions](https://github.com/coze-dev/coze-studio/discussions)
- **规范咨询**: zker-specs@example.com

---

**🎯 目标**: 100% 规范覆盖率，确保代码质量和全局一致性！
