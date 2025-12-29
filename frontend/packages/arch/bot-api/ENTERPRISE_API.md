# ZKER 企业级前端 API 使用指南

**版本**: v1.0.0
**最后更新**: 2025-01-01

---

## 📋 目录

- [概述](#概述)
- [安装与导入](#安装与导入)
- [租户管理 API](#租户管理-api)
- [配额管理 API](#配额管理-api)
- [权限管理 API](#权限管理-api)
- [错误处理](#错误处理)
- [TypeScript 类型定义](#typescript-类型定义)
- [最佳实践](#最佳实践)

---

## 🎯 概述

ZKER 企业级前端 API 层提供完整的 TypeScript 类型定义和 API 调用封装，支持：

- ✅ **租户管理**: 租户 CRUD、订阅管理、套餐升级
- ✅ **配额管理**: 配额查询、检查、限额设置、使用详情
- ✅ **权限管理**: 角色 CRUD、部门管理、权限检查

### API 列表

| 模块 | API 服务 | 文件位置 |
|------|---------|---------|
| 租户管理 | `tenantApi` | `tenant-api.ts` |
| 配额管理 | `quotaApi` | `quota-api.ts` |
| 权限管理 | `permissionApi` | `permission-api.ts` |

---

## 📦 安装与导入

### 导入方式

```typescript
// 方式一：从 bot-api 包导入
import { tenantApi, quotaApi, permissionApi } from '@coze-arch/bot-api';

// 方式二：从 enterprise-api 模块导入
import {
  tenantApi,
  quotaApi,
  permissionApi,
  // 类型定义
  type Tenant,
  type Subscription,
  type Role,
  type ResourceType,
} from '@coze-arch/bot-api';
```

### 依赖说明

```json
{
  "dependencies": {
    "@coze-arch/bot-api": "workspace:*",
    "@coze-arch/bot-http": "workspace:*"
  }
}
```

---

## 🏢 租户管理 API

### API 方法

| 方法 | 说明 |
|------|------|
| `createTenant()` | 创建租户 |
| `getTenant()` | 获取租户详情 |
| `updateTenant()` | 更新租户信息 |
| `listTenants()` | 列出租户（分页） |
| `deleteTenant()` | 删除租户（软删除） |
| `getSubscription()` | 获取订阅信息 |
| `updateSubscription()` | 更新订阅信息 |
| `upgradeSubscription()` | 升级订阅套餐 |

### 使用示例

```typescript
import { tenantApi, TenantStatus } from '@coze-arch/bot-api';

// 1. 创建租户
const newTenant = await tenantApi.createTenant({
  tenant_name: '示例公司',
  contact_email: 'admin@example.com',
  contact_phone: '+86-138-0000-0000',
  company_name: '示例公司',
});

console.log('创建成功:', newTenant.tenant);

// 2. 获取租户详情
const tenantDetail = await tenantApi.getTenant({
  tenant_id: 'tenant_001',
});

console.log('租户信息:', tenantDetail.tenant);
console.log('订阅信息:', tenantDetail.subscription);

// 3. 列出租户（分页）
const tenantList = await tenantApi.listTenants({
  page: 1,
  page_size: 20,
  status: TenantStatus.ACTIVE,
  search: '示例',
});

console.log('租户列表:', tenantList.tenants);
console.log('总数:', tenantList.total);

// 4. 更新租户
const updatedTenant = await tenantApi.updateTenant({
  tenant_id: 'tenant_001',
  tenant_name: '新名称',
  contact_email: 'new@example.com',
});

// 5. 升级订阅套餐
const upgradedSubscription = await tenantApi.upgradeSubscription({
  tenant_id: 'tenant_001',
  target_plan_tier: 'professional',
  target_billing_cycle: 'yearly',
  discount_code: 'SAVE20',
});

console.log('升级后订阅:', upgradedSubscription.subscription);
console.log('新配额:', upgradedSubscription.quotas);
```

---

## 📊 配额管理 API

### API 方法

| 方法 | 说明 |
|------|------|
| `getQuotas()` | 获取租户配额 |
| `checkQuota()` | 检查配额（操作前） |
| `updateQuotaLimit()` | 更新配额限制 |
| `getQuotaUsage()` | 获取配额使用详情 |
| `batchCheckQuota()` | 批量检查配额 |

### 使用示例

```typescript
import { quotaApi, ResourceType } from '@coze-arch/bot-api';

// 1. 获取租户配额
const quotas = await quotaApi.getQuotas({
  tenant_id: 'tenant_001',
  resource_types: [
    ResourceType.BOTS,
    ResourceType.WORKFLOWS,
    ResourceType.MESSAGES,
  ],
});

console.log('Bot 配额:', quotas.quotas.find(q => q.resource_type === ResourceType.BOTS));
console.log('整体使用率:', quotas.overall_usage_percentage);

// 2. 检查配额（创建 Bot 前检查）
const quotaCheck = await quotaApi.checkQuota({
  tenant_id: 'tenant_001',
  resource_type: ResourceType.BOTS,
  required_count: 1,
  operation_description: '创建新 Bot',
});

if (quotaCheck.allowed) {
  console.log('配额充足，可以创建');
  console.log('剩余配额:', quotaCheck.remaining);
} else {
  console.log('配额不足:', quotaCheck.reason);
  if (quotaCheck.suggested_upgrade) {
    console.log('建议升级到:', quotaCheck.suggested_upgrade.target_plan_tier);
  }
}

// 3. 获取配额使用详情
const usage = await quotaApi.getQuotaUsage({
  tenant_id: 'tenant_001',
  resource_type: ResourceType.MESSAGES,
  from: '2025-01-01T00:00:00Z',
  to: '2025-01-31T23:59:59Z',
});

console.log('当前使用量:', usage.usage.current_usage);
console.log('使用趋势:', usage.usage.trend);
console.log('预计耗尽时间:', usage.estimated_exhaustion_at);

// 4. 批量检查配额
const batchResults = await quotaApi.batchCheckQuota(
  'tenant_001',
  [
    { resource_type: ResourceType.BOTS, required_count: 1 },
    { resource_type: ResourceType.WORKFLOWS, required_count: 2 },
  ],
);

batchResults.forEach(result => {
  console.log(`${result.resource_type}: ${result.allowed ? '✓' : '✗'}`);
});
```

---

## 🔐 权限管理 API

### API 方法

#### 角色管理
| 方法 | 说明 |
|------|------|
| `createRole()` | 创建角色 |
| `getRole()` | 获取角色详情 |
| `updateRole()` | 更新角色 |
| `listRoles()` | 列出角色 |
| `deleteRole()` | 删除角色 |
| `assignRole()` | 分配角色 |
| `revokeRole()` | 撤销角色 |
| `getUserRoles()` | 获取用户角色 |

#### 部门管理
| 方法 | 说明 |
|------|------|
| `createDepartment()` | 创建部门 |
| `getDepartment()` | 获取部门详情 |
| `updateDepartment()` | 更新部门 |
| `listDepartments()` | 列出部门 |
| `deleteDepartment()` | 删除部门 |
| `addDepartmentMember()` | 添加部门成员 |
| `removeDepartmentMember()` | 移除部门成员 |

#### 权限检查
| 方法 | 说明 |
|------|------|
| `checkDataPermission()` | 检查数据权限 |
| `checkFieldPermission()` | 检查字段权限 |
| `batchCheckDataPermission()` | 批量检查数据权限 |

### 使用示例

```typescript
import {
  permissionApi,
  PermissionScope,
  FieldPermissionLevel,
} from '@coze-arch/bot-api';

// 1. 创建角色
const newRole = await permissionApi.createRole({
  tenant_id: 'tenant_001',
  role_name: '开发者',
  role_code: 'developer',
  description: '可以创建和管理 Bot',
  permissions: [
    'bots.create',
    'bots.view',
    'bots.edit',
    'bots.delete',
  ],
  data_permissions: [
    {
      resource_type: 'bots',
      scope: PermissionScope.OWN,
    },
  ],
  field_permissions: [
    {
      resource_type: 'bots',
      field_name: 'api_key',
      permission_level: FieldPermissionLevel.HIDDEN,
    },
  ],
});

// 2. 列出角色
const roles = await permissionApi.listRoles({
  tenant_id: 'tenant_001',
  page: 1,
  page_size: 20,
});

// 3. 分配角色
const assigned = await permissionApi.assignRole({
  tenant_id: 'tenant_001',
  user_id: 'user_001',
  role_id: 'role_developer',
  is_primary: true,
});

// 4. 获取用户角色
const userRoles = await permissionApi.getUserRoles({
  tenant_id: 'tenant_001',
  user_id: 'user_001',
});

console.log('用户角色:', userRoles.roles);

// 5. 创建部门
const department = await permissionApi.createDepartment({
  tenant_id: 'tenant_001',
  department_name: '研发部',
  parent_department_id: 'root_dept',
});

// 6. 添加部门成员
await permissionApi.addDepartmentMember({
  tenant_id: 'tenant_001',
  department_id: 'dept_001',
  user_id: 'user_001',
  is_leader: true,
});

// 7. 检查数据权限
const dataPermission = await permissionApi.checkDataPermission({
  tenant_id: 'tenant_001',
  user_id: 'user_001',
  resource_type: 'bots',
  action: 'edit',
  resource_ids: ['bot_001', 'bot_002', 'bot_003'],
});

console.log('允许访问:', dataPermission.allowed_resource_ids);
console.log('拒绝访问:', dataPermission.denied_resource_ids);

// 8. 检查字段权限
const fieldPermission = await permissionApi.checkFieldPermission({
  tenant_id: 'tenant_001',
  user_id: 'user_001',
  resource_type: 'bots',
  action: 'view',
  field_names: ['name', 'api_key', 'config'],
});

console.log('字段权限:', fieldPermission.field_permissions);
```

---

## ⚠️ 错误处理

### 错误响应格式

所有 API 调用都可能抛出错误，建议使用 try-catch 处理：

```typescript
import { tenantApi } from '@coze-arch/bot-api';

try {
  const tenant = await tenantApi.getTenant({ tenant_id: 'xxx' });
  console.log('成功:', tenant);
} catch (error) {
  console.error('错误:', error);

  // 错误通常包含以下字段：
  // - code: 错误码
  // - msg: 错误消息
  // - details: 详细信息（可选）

  if (error.code === 108000001) {
    console.error('参数错误');
  } else if (error.code === 108000002) {
    console.error('租户不存在');
  }
}
```

### 常见错误码

| 错误码 | 说明 |
|--------|------|
| 108000001 | 参数错误 |
| 108000002 | 租户不存在 |
| 108000003 | 订阅已过期 |
| 108010001 | 配额不足 |
| 108020001 | 角色不存在 |
| 108020002 | 权限不足 |

### 禁用错误提示

某些场景下可能需要禁用自动错误提示（Toast）：

```typescript
import { tenantApi } from '@coze-arch/bot-api';

const tenant = await tenantApi.getTenant(
  { tenant_id: 'xxx' },
  {
    __disableErrorToast: true, // 禁用错误提示
  },
);
```

---

## 🔷 TypeScript 类型定义

### 主要类型

```typescript
import type {
  // 租户相关
  Tenant,
  Subscription,
  TenantStatus,
  SubscriptionTier,

  // 配额相关
  Quota,
  QuotaUsageDetail,
  ResourceType,
  QuotaPeriod,

  // 权限相关
  Role,
  DataPermission,
  FieldPermission,
  Department,
  PermissionScope,
  FieldPermissionLevel,
} from '@coze-arch/bot-api';
```

### 类型使用示例

```typescript
import type { Tenant, SubscriptionTier } from '@coze-arch/bot-api';

const processTenant = (tenant: Tenant) => {
  console.log(`租户: ${tenant.tenant_name}`);
  console.log(`状态: ${tenant.status}`);
};

const upgradeTier: SubscriptionTier = SubscriptionTier.PROFESSIONAL;
```

---

## ✨ 最佳实践

### 1. 统一错误处理

```typescript
// utils/api.ts
import { tenantApi } from '@coze-arch/bot-api';

export const getTenantWithErrorHandling = async (tenantId: string) => {
  try {
    return await tenantApi.getTenant({ tenant_id: tenantId });
  } catch (error) {
    if (error.code === 108000002) {
      console.error('租户不存在');
    } else {
      console.error('未知错误:', error);
    }
    throw error;
  }
};
```

### 2. 请求重试

```typescript
import { quotaApi } from '@coze-arch/bot-api';

const checkQuotaWithRetry = async (
  tenantId: string,
  resourceType: ResourceType,
  requiredCount: number,
  maxRetries = 3,
) => {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await quotaApi.checkQuota({
        tenant_id: tenantId,
        resource_type: resourceType,
        required_count: requiredCount,
      });
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
    }
  }
};
```

### 3. 批量操作

```typescript
import { permissionApi } from '@coze-arch/bot-api';

const batchAssignRole = async (
  tenantId: string,
  roleIds: string[],
  userIds: string[],
) => {
  const results = await Promise.allSettled(
    userIds.map(userId =>
      permissionApi.assignRole({
        tenant_id: tenantId,
        user_id: userId,
        role_id: roleIds[0],
      }),
    ),
  );

  const succeeded = results.filter(r => r.status === 'fulfilled').length;
  const failed = results.filter(r => r.status === 'rejected').length;

  console.log(`成功: ${succeeded}, 失败: ${failed}`);
  return results;
};
```

### 4. React Hooks 集成

```typescript
// hooks/useTenant.ts
import { useState, useEffect } from 'react';
import { tenantApi } from '@coze-arch/bot-api';
import type { Tenant } from '@coze-arch/bot-api';

export const useTenant = (tenantId: string) => {
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchTenant = async () => {
      try {
        setLoading(true);
        const response = await tenantApi.getTenant({ tenant_id: tenantId });
        setTenant(response.tenant);
      } catch (err) {
        setError(err as Error);
      } finally {
        setLoading(false);
      }
    };

    fetchTenant();
  }, [tenantId]);

  return { tenant, loading, error };
};
```

---

## 📚 相关文档

- [后端 API 文档](../../../../../backend/api/tenant/v1/tenant.proto)
- [类型定义](./types/tenant.types.ts)
- [企业级开发规范](../../../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

---

**更新日志**：

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，完整的租户、配额、权限管理 API | Claude AI |
