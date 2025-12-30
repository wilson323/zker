# 组织中心路由注册文档

## 概述

本目录包含组织中心的两个路由注册文件，用于注册63个API端点：

1. **org_router.go** - 组织、部门、岗位、员工、通讯录路由 (48个)
2. **hr_router.go** - HR生命周期管理路由 (15个)

## 路由统计

### org_router.go (48个路由)

| 模块 | 路由数量 | 路径前缀 |
|------|---------|---------|
| 组织管理 | 10 | `/api/organizations/*` |
| 部门管理 | 11 | `/api/org/departments/*` |
| 岗位管理 | 9 | `/api/org/positions/*` |
| 员工管理 | 12 | `/api/org/employees/*` |
| 通讯录 | 6 | `/api/org/directory/*` |

### hr_router.go (15个路由)

| 模块 | 路由数量 | 路径前缀 |
|------|---------|---------|
| 合同管理 | 5 | `/api/hr/contracts/*` |
| 调岗管理 | 3 | `/api/hr/employees/*`, `/api/hr/transfers/*` |
| 离职管理 | 6 | `/api/hr/resignations/*` |
| 试用期管理 | 1 | `/api/hr/employees/probation/*` |

## 使用方法

### 1. 初始化Handler

在应用启动时，需要初始化所有Handler并注入依赖：

```go
package main

import (
    "github.com/coze-dev/coze-studio/backend/api/router/coze/org"
    "github.com/coze-dev/coze-studio/backend/api/router/coze/hr"
    orgdomain "github.com/coze-dev/coze-studio/backend/domain/org/service"
    "github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
    r := server.New()

    // 1. 初始化领域服务
    orgService := orgdomain.NewOrganizationService(...)
    deptService := orgdomain.NewDepartmentService(...)
    positionService := orgdomain.NewPositionService(...)
    hrService := orgdomain.NewHRLifecycleService(...)

    // 2. 初始化Handler
    orgHandlers := &org.OrgHandlers{
        Organization: org.NewOrganizationHandler(orgService),
        Department:   org.NewDepartmentHandler(deptService),
        Position:     org.NewPositionHandler(positionService),
        // Employee和Directory使用函数式Handler，已在路由注册时适配
    }

    // 3. 注册路由
    org.RegisterOrgRoutes(r, orgHandlers)
    hr.RegisterHRRoutes(r, hrHandler)

    // 4. 启动服务器
    r.Spin()
}
```

### 2. 租户隔离

所有路由都自动应用 `TenantIsolationMiddleware()` 中间件，实现租户隔离：

- **自动提取tenant_id**: 从HTTP Header (`X-Tenant-ID`)、Session、JWT中自动提取
- **自动注入context**: tenant_id自动注入到context中，Handler可通过 `middleware.GetTenantIDFromContext(ctx)` 获取
- **自动验证租户**: 验证租户存在且激活，防止跨租户访问
- **自动数据过滤**: Service层根据context中的tenant_id自动过滤数据

### 3. API端点列表

#### 组织管理 (10个路由)

```
POST   /api/organizations              创建组织
GET    /api/organizations/:id          获取组织详情
PUT    /api/organizations/:id          更新组织
DELETE /api/organizations/:id          删除组织
GET    /api/organizations/tree         获取组织树
GET    /api/organizations              分页查询组织列表
POST   /api/organizations/:id/move     移动组织
GET    /api/organizations/:id/children 获取子组织
GET    /api/organizations/:id/ancestors 获取祖先组织
GET    /api/organizations/:id/descendants 获取后代组织
```

#### 部门管理 (11个路由)

```
POST   /api/org/departments              创建部门
GET    /api/org/departments/:id          获取部门详情
PUT    /api/org/departments/:id          更新部门
DELETE /api/org/departments/:id          删除部门
GET    /api/org/departments/tree         获取部门树
GET    /api/org/departments              分页查询部门列表
POST   /api/org/departments/:id/move     移动部门
GET    /api/org/departments/:id/children 获取子部门
GET    /api/org/departments/:id/ancestors 获取祖先部门
GET    /api/org/departments/:id/descendants 获取后代部门
GET    /api/org/departments/org/:org_id  获取组织的所有部门
```

#### 岗位管理 (9个路由)

```
POST   /api/org/positions              创建岗位
GET    /api/org/positions/:id          获取岗位详情
PUT    /api/org/positions/:id          更新岗位
DELETE /api/org/positions/:id          删除岗位
GET    /api/org/positions              分页查询岗位列表
GET    /api/org/positions/code/:code   按岗位编码查询
GET    /api/org/positions/dept/:dept_id 按部门查询岗位
GET    /api/org/positions/level/:level 按职级查询岗位
GET    /api/org/positions/category/:category 按类别查询岗位
```

#### 员工管理 (12个路由)

```
POST   /api/org/employees              创建员工
GET    /api/org/employees/:id          获取员工详情
PUT    /api/org/employees/:id          更新员工
DELETE /api/org/employees/:id          删除员工
GET    /api/org/employees              分页查询员工列表
GET    /api/org/employees/code/:code   按员工编号查询
GET    /api/org/employees/search       搜索员工
GET    /api/org/employees/dept/:dept_id 按部门查询员工
GET    /api/org/employees/position/:position_id 按岗位查询员工
POST   /api/org/employees/:id/activate   激活员工
POST   /api/org/employees/:id/deactivate 停用员工
GET    /api/org/employees/:id/status    获取员工状态
```

#### 通讯录 (6个路由)

```
GET    /api/org/directory           获取通讯录列表
GET    /api/org/directory/tree      获取通讯录树
GET    /api/org/directory/org/:org_id 按组织查询通讯录
GET    /api/org/directory/dept/:dept_id 按部门查询通讯录
GET    /api/org/directory/export    导出通讯录
POST   /api/org/directory/import    导入通讯录
```

#### HR生命周期管理 (15个路由)

```
# 合同管理
POST   /api/hr/contracts                      创建合同
GET    /api/hr/contracts/:id                  获取合同详情
POST   /api/hr/contracts/:id/sign             签署合同
GET    /api/hr/employees/:emp_id/contracts    获取员工的所有合同
GET    /api/hr/employees/:emp_id/contracts/active 获取员工的生效合同

# 调岗管理
POST   /api/hr/employees/:emp_id/transfer     调岗
GET    /api/hr/employees/:emp_id/transfers    获取员工的调岗记录
GET    /api/hr/transfers/:id                  获取调岗记录详情

# 离职管理
POST   /api/hr/employees/:emp_id/resign       员工离职
GET    /api/hr/employees/:emp_id/resignation  获取员工的离职记录
GET    /api/hr/resignations/:id               获取离职记录详情
POST   /api/hr/resignations/:id/approve       审批离职
PUT    /api/hr/resignations/:id/handover      更新交接状态
GET    /api/hr/resignations/pending           获取待审批的离职列表

# 试用期管理
GET    /api/hr/employees/probation/upcoming   获取即将结束试用期的员工列表
```

## 中间件配置

所有路由已自动配置以下中间件：

1. **TenantIsolationMiddleware** - 租户隔离中间件
   - 从请求中提取tenant_id
   - 验证租户有效性
   - 注入context供后续使用

2. **认证中间件** - 需要在路由组外部配置
   ```go
   apiGroup := r.Group("/api", middleware.AuthMiddleware(), middleware.TenantIsolationMiddleware())
   ```

3. **日志中间件** - 需要在路由组外部配置
   ```go
   apiGroup := r.Group("/api", middleware.LogMiddleware(), middleware.TenantIsolationMiddleware())
   ```

## 错误处理

所有Handler使用统一的错误响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

错误码定义在 `backend/types/errno` 包中，遵循企业级错误码规范。

## 测试

可以使用curl或Postman测试API：

```bash
# 创建组织
curl -X POST http://localhost:8888/api/organizations \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_123" \
  -d '{
    "tenant_id": "tenant_123",
    "org_name": "技术部",
    "org_type": "department",
    "org_code": "TECH",
    "description": "负责技术研发"
  }'

# 获取组织详情
curl -X GET http://localhost:8888/api/organizations/org_123 \
  -H "X-Tenant-ID: tenant_123"

# 获取组织树
curl -X GET "http://localhost:8888/api/organizations/tree?tenant_id=tenant_123" \
  -H "X-Tenant-ID: tenant_123"
```

## 注意事项

1. **租户隔离**: 所有API必须在请求Header中携带 `X-Tenant-ID`，或在Session/JWT中包含tenant_id
2. **权限控制**: 部分接口需要特定权限（如管理员权限），需在应用层配置权限中间件
3. **分页参数**: 列表查询接口支持 `page` 和 `page_size` 参数，默认值为 `page=1&page_size=20`
4. **软删除**: 删除操作为软删除，数据不会物理删除，只是标记为已删除
5. **级联操作**: 删除组织/部门时，会检查是否有子组织/子部门/员工，需要先处理关联数据

## 全局一致性

遵循企业级开发规范：

- ✅ RESTful API设计规范
- ✅ 统一错误码系统
- ✅ 租户隔离中间件
- ✅ 统一响应格式
- ✅ 完整的文档注释
- ✅ 清晰的代码结构

## 相关文档

- [企业级开发规范手册](../../../../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [统一错误码定义规范](../../../../../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
- [多租户SaaS架构](../../../../../../docs/企业级功能完善与统一性设计方案/zker_MultiTenant_SaaS_完整架构设计文档.md)
