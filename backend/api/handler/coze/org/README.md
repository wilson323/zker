# Organization Handler 使用指南

## 概述

`OrganizationHandler` 是组织管理的HTTP处理器，实现了10个RESTful API端点，严格遵循DDD架构和SOLID原则。

## 架构设计

### 分层架构

```
┌─────────────────────────────────────┐
│   HTTP Layer (OrganizationHandler)  │  ← 处理HTTP请求/响应
├─────────────────────────────────────┤
│   Application Layer (Service)       │  ← 业务逻辑编排
├─────────────────────────────────────┤
│   Domain Layer (Entity/Repository)  │  ← 核心业务逻辑
└─────────────────────────────────────┘
```

### SOLID原则

- **单一职责**: Handler仅处理HTTP层，不包含业务逻辑
- **开闭原则**: 通过依赖注入OrganizationService，易于扩展
- **里氏替换**: 依赖接口而非具体实现
- **接口隔离**: 每个API端点职责单一
- **依赖倒置**: Handler依赖service抽象，不依赖具体实现

## API端点

### 基础CRUD操作

#### 1. 创建组织

```http
POST /api/organizations
Content-Type: application/json

{
  "tenant_id": "tenant_001",
  "org_name": "技术部",
  "org_type": "department",
  "parent_id": "org_parent_001",
  "org_code": "TECH",
  "leader_id": "user_001",
  "description": "负责技术研发",
  "sort_order": 1
}
```

**响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "org_id": "org_001",
    "tenant_id": "tenant_001",
    "org_name": "技术部",
    "org_type": "department",
    "org_code": "TECH",
    "level": 2,
    "path": "/org_parent_001/org_001",
    "status": "active"
  }
}
```

#### 2. 获取组织详情

```http
GET /api/organizations/:id
```

#### 3. 更新组织

```http
PUT /api/organizations/:id
Content-Type: application/json

{
  "org_name": "技术研发中心",
  "leader_id": "user_002",
  "description": "负责核心技术研发",
  "sort_order": 2,
  "status": "active"
}
```

#### 4. 删除组织

```http
DELETE /api/organizations/:id
```

### 组织树操作

#### 5. 获取组织树

```http
GET /api/organizations/tree?tenant_id=tenant_001
```

**响应**: 返回完整的树形结构

#### 6. 分页查询组织列表

```http
GET /api/organizations?tenant_id=tenant_001&page=1&page_size=20&sort_by=created_at&sort_order=desc
```

**查询参数**:
- `tenant_id` (必填): 租户ID
- `org_type` (可选): 组织类型 (company/division/department/project)
- `status` (可选): 状态 (active/inactive/frozen)
- `parent_id` (可选): 父组织ID
- `keyword` (可选): 搜索关键词
- `page` (可选): 页码，默认1
- `page_size` (可选): 每页数量，默认20，最大100
- `sort_by` (可选): 排序字段 (created_at/updated_at/sort_order/org_name)
- `sort_order` (可选): 排序方向 (asc/desc)

#### 7. 移动组织

```http
POST /api/organizations/:id/move
Content-Type: application/json

{
  "org_id": "org_001",
  "new_parent_id": "org_new_parent_001"
}
```

### 组织关系查询

#### 8. 获取子组织

```http
GET /api/organizations/:id/children
```

#### 9. 获取祖先组织

```http
GET /api/organizations/:id/ancestors
```

#### 10. 获取后代组织

```http
GET /api/organizations/:id/descendants
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 参数验证失败 |
| 404 | 组织不存在 |
| 4004001 | 组织编码已存在 |
| 4004002 | 父组织不存在 |
| 4004003 | 租户不匹配 |
| 4004004 | 组织类型无效 |
| 4004005 | 有子组织，无法删除 |
| 4004006 | 根组织无法删除 |
| 4004007 | 根组织必须为公司类型 |
| 4004008 | 无法移动到自己的后代 |
| 4004009 | 无法移动根组织 |
| 500 | 内部服务器错误 |

## 使用示例

### 初始化Handler

```go
package main

import (
    "github.com/coze-dev/coze-studio/backend/api/handler/coze/org"
    orgservice "github.com/coze-dev/coze-studio/backend/domain/org/service"
)

func main() {
    // 创建服务实例
    orgService := orgservice.NewOrganizationService(
        orgRepository,
        treeRepository,
        db,
    )

    // 创建Handler实例
    handler := org.NewOrganizationHandler(orgService)

    // 注册路由
    r := gin.Default()
    r.POST("/api/organizations", handler.CreateOrganization)
    r.GET("/api/organizations/:id", handler.GetOrganization)
    // ... 其他路由
}
```

### 在Hertz中注册路由

```go
package main

import (
    "github.com/cloudwego/hertz/pkg/app/server"
    orghandler "github.com/coze-dev/coze-studio/backend/api/handler/coze/org"
)

func main() {
    h := server.Default()

    orgHandler := orghandler.NewOrganizationHandler(orgService)

    organization := h.Group("/api/organizations")
    {
        organization.POST("", orgHandler.CreateOrganization)
        organization.GET("/:id", orgHandler.GetOrganization)
        organization.PUT("/:id", orgHandler.UpdateOrganization)
        organization.DELETE("/:id", orgHandler.DeleteOrganization)
        organization.GET("/tree", orgHandler.GetOrganizationTree)
        organization.GET("", orgHandler.ListOrganizations)
        organization.POST("/:id/move", orgHandler.MoveOrganization)
        organization.GET("/:id/children", orgHandler.GetChildren)
        organization.GET("/:id/ancestors", orgHandler.GetAncestors)
        organization.GET("/:id/descendants", orgHandler.GetDescendants)
    }

    h.Spin()
}
```

## 测试

### 单元测试示例

```go
func TestCreateOrganization(t *testing.T) {
    // Mock service
    mockService := new(MockOrganizationService)
    handler := org.NewOrganizationHandler(mockService)

    // 准备请求数据
    req := &org.CreateOrganizationRequest{
        TenantID:    "tenant_001",
        OrgName:     "测试组织",
        OrgType:     "department",
        OrgCode:     "TEST",
        Description: "测试描述",
    }

    // 执行测试
    ctx := context.Background()
    c := &app.RequestContext{}

    // 设置请求体
    body, _ := json.Marshal(req)
    c.Request.SetBody(body)

    // 调用Handler
    handler.CreateOrganization(ctx, c)

    // 验证响应
    assert.Equal(t, 200, c.Response.StatusCode())
}
```

## 数据一致性

### 事务处理

所有涉及多表操作的方法都使用数据库事务：

- `CreateOrganization`: 创建组织 + 闭包表路径
- `DeleteOrganization`: 软删除组织 + 删除闭包表路径
- `MoveOrganization`: 更新组织 + 更新闭包表路径

### 租户隔离

所有查询都自动添加`tenant_id`过滤：

```go
// 在Repository层自动实现
WHERE tenant_id = ? AND deleted_at IS NULL
```

### 软删除

使用`deleted_at`字段实现软删除，数据不会物理删除：

```sql
UPDATE organizations SET deleted_at = NOW() WHERE org_id = ?
```

## 性能优化

### 索引建议

```sql
-- 租户隔离索引
CREATE INDEX idx_tenant_id ON organizations(tenant_id, deleted_at);

-- 编码唯一索引
CREATE UNIQUE INDEX uk_tenant_code ON organizations(tenant_id, org_code);

-- 状态查询索引
CREATE INDEX idx_status ON organizations(status, deleted_at);

-- 排序索引
CREATE INDEX idx_sort ON organizations(tenant_id, sort_order, created_at);

-- 父组织查询索引
CREATE INDEX idx_parent_id ON organizations(parent_id, deleted_at);
```

### 查询优化

- 使用闭包表实现高效的树形查询
- 避免N+1查询，使用JOIN或Preload
- 分页查询使用游标分页而非OFFSET

## 监控指标

建议监控以下指标：

- API响应时间 (P50, P95, P99)
- API错误率
- 数据库查询时间
- 事务成功率
- 租户级别的组织数量

## 安全建议

1. **权限验证**: 在中间件层添加租户权限验证
2. **输入验证**: 使用`binding` tag进行参数验证
3. **SQL注入**: 使用GORM参数化查询
4. **敏感信息**: 日志中不要输出组织详情的敏感字段

## 常见问题

### Q: 如何处理大量组织数据的查询？

A: 使用分页查询 + 游标分页，避免深度分页性能问题。

### Q: 组织层级有限制吗？

A: 理论上无限制，但建议不超过5层，避免查询复杂度过高。

### Q: 如何实现组织权限继承？

A: 在闭包表中存储权限信息，查询时自动获取所有祖先的权限。

## 更新日志

### v1.0.0 (2025-01-01)

- 实现10个RESTful API端点
- 支持组织树的CRUD操作
- 支持组织移动和层级查询
- 完整的参数验证和错误处理
- 遵循企业级开发规范
