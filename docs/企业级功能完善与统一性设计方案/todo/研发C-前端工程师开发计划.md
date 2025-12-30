# 研发C - 前端工程师开发计划

**负责人**: 研发C（前端工程师）
**开发周期**: 5周（4人并行）
**总工作量**: 23人天
**核心职责**: API调用层、状态管理、路由系统、管理页面、UI组件
**最后更新**: 2025-12-30

---

## 📋 工作量总览

| 周次 | 主要任务 | 工作量 | 优先级 |
|------|---------|--------|--------|
| **Week 1-2** | 前端API调用层 + 类型定义 | 6人天 | **P1** ⭐ |
| **Week 2** | 前端状态管理 | 4人天 | **P1** ⭐ |
| **Week 2** | 前端路由系统 | 3人天 | **P1** ⭐ |
| **Week 3-4** | 前端管理页面 | 11人天 | **P1** ⭐ |
| **Week 5** | UI组件补充 + 单元测试 | 7人天 | **P2** |

**📌 注意**: 前端实现完成度仅**35%**，你是用户可见界面开发的主力！

---

## ⚠️ 核心注意事项

### 1. 代码隔离原则

**你的职责范围**:
- ✅ 前端API调用层（TypeScript）
- ✅ 前端状态管理（Zustand + React Query）
- ✅ 前端路由系统（React Router）
- ✅ 前端管理页面（React组件）
- ✅ UI组件补充（Semi Design扩展）
- ✅ 前端单元测试

**不要触碰**:
- ❌ 后端代码（研发A、研发B负责）
- ❌ CI/CD配置（研发D负责）
- ❌ 监控大盘配置（研发D负责）
- ❌ 数据库相关（研发B负责）

### 2. 前端开发红线

**必须遵守**:
1. **所有组件使用TypeScript**，禁止使用any
2. **组件必须是函数式组件** + Hooks，不使用类组件
3. **样式使用Tailwind CSS**，不使用内联样式
4. **状态管理使用Zustand**（全局状态）+ React Query（服务端状态）
5. **API调用必须统一**，不允许直接使用fetch/axios
6. **所有组件要有PropTypes**或TypeScript接口定义
7. **组件文件不超过300行**

**禁止行为**:
1. 使用类组件
2. 使用any类型
3. 直接使用fetch/axios
4. 硬编码API URL
5. 不写TypeScript类型定义

---

## 🎯 Week 1-2: 前端API调用层 + 类型定义 (P1)

### 任务1.1: 创建API Client包结构 (Day 1)

**目录结构**:

```bash
frontend/packages/api-client/
├── src/
│   ├── http/
│   │   ├── client.ts           # Axios封装
│   │   ├── interceptor.ts      # 请求/响应拦截器
│   │   └── errorHandler.ts     # 统一错误处理
│   ├── hooks/
│   │   ├── useTenant.ts        # 租户相关Hooks
│   │   ├── useRole.ts          # 角色相关Hooks
│   │   ├── usePermission.ts    # 权限相关Hooks
│   │   ├── useRouting.ts       # 路由相关Hooks
│   │   └── useQuota.ts         # 配额相关Hooks
│   ├── types/
│   │   ├── tenant.ts           # 租户类型定义
│   │   ├── role.ts             # 角色类型定义
│   │   ├── permission.ts       # 权限类型定义
│   │   ├── routing.ts          # 路由类型定义
│   │   └── quota.ts            # 配额类型定义
│   ├── utils/
│   │   ├── apiEndpoints.ts     # API端点定义
│   │   └── request.ts          # 请求工具函数
│   └── index.ts                # 统一导出
├── package.json
└── tsconfig.json
```

**初始化命令**:

```bash
cd frontend/packages
mkdir -p api-client/src/{http,hooks,types,utils}
cd api-client
npm init -y
npm install axios @tanstack/react-query
npm install -D @types/node typescript
```

**⚠️ 注意事项**:
1. **包名使用@coze-studio/api-client**
2. **TypeScript严格模式开启**
3. **使用workspace:协议引用内部包**

**📖 开发规范**:
- 目录结构清晰
- 依赖管理规范
- TypeScript配置正确

**🔗 设计文档链接**:
- [前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md)

---

### 任务1.2: HTTP客户端封装 (Day 1-2)

**文件**: `frontend/packages/api-client/src/http/client.ts`

```typescript
// frontend/packages/api-client/src/http/client.ts
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

// 创建axios实例
const createClient = (): AxiosInstance => {
  const client = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8001',
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // 请求拦截器
  client.interceptors.request.use(
    (config) => {
      // 添加token
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }

      // 添加request_id
      config.headers['X-Request-ID'] = generateRequestId();

      // 添加tenant_id
      const tenantId = localStorage.getItem('tenant_id');
      if (tenantId) {
        config.headers['X-Tenant-ID'] = tenantId;
      }

      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // 响应拦截器
  client.interceptors.response.use(
    (response: AxiosResponse) => {
      return response.data;
    },
    (error) => {
      // 统一错误处理
      return Promise.reject(handleError(error));
    }
  );

  return client;
};

// 生成request_id
const generateRequestId = (): string => {
  return `req_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
};

// 导出单例
export const apiClient = createClient();

// 通用请求方法
export const request = {
  get: <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    return apiClient.get(url, config);
  },

  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    return apiClient.post(url, data, config);
  },

  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    return apiClient.put(url, data, config);
  },

  delete: <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    return apiClient.delete(url, config);
  },

  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    return apiClient.patch(url, data, config);
  },
};
```

**文件**: `frontend/packages/api-client/src/http/errorHandler.ts`

```typescript
// frontend/packages/api-client/src/http/errorHandler.ts
import { AxiosError } from 'axios';

interface ErrorResponse {
  code: string;
  message: string;
  message_zh: string;
  message_en: string;
  details?: any;
  request_id: string;
}

// 统一错误处理
export const handleError = (error: AxiosError): Error => {
  if (error.response) {
    const { status, data } = error.response;
    const errorData = data as ErrorResponse;

    // 根据状态码处理
    switch (status) {
      case 401:
        return new Error('未登录或登录已过期');
      case 403:
        return new Error(errorData.message_zh || '权限不足');
      case 404:
        return new Error(errorData.message_zh || '资源不存在');
      case 402:
        return new Error(errorData.message_zh || '配额已用完，请升级订阅');
      case 500:
        return new Error('服务器内部错误，请稍后重试');
      default:
        return new Error(errorData.message_zh || error.message);
    }
  }

  if (error.request) {
    return new Error('网络错误，请检查您的网络连接');
  }

  return new Error(error.message);
};

// 显示错误提示
export const showError = (error: Error) => {
  // TODO: 集成Toast组件
  console.error('[API Error]', error.message);
};
```

**⚠️ 注意事项**:
1. **token和tenant_id从localStorage获取**
2. **所有请求自动添加request_id**
3. **错误处理统一，显示友好提示**
4. **支持TypeScript泛型**

**📖 开发规范**:
- 单一职责
- TypeScript类型完整
- 错误处理优雅

**🔗 设计文档链接**:
- [前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md) - 任务1章节

---

### 任务1.3: TypeScript类型定义 (Day 2)

**文件**: `frontend/packages/api-client/src/types/tenant.ts`

```typescript
// frontend/packages/api-client/src/types/tenant.ts

/**
 * 租户类型
 */
export type TenantType = 'individual' | 'team' | 'enterprise';

/**
 * 租户状态
 */
export type TenantStatus = 'active' | 'suspended' | 'deleted';

/**
 * 订阅层级
 */
export type SubscriptionTier = 'free' | 'pro' | 'enterprise';

/**
 * 租户信息
 */
export interface Tenant {
  tenant_id: string;
  tenant_name: string;
  tenant_type: TenantType;
  status: TenantStatus;
  subscription_tier: SubscriptionTier;
  created_at: string;
  updated_at: string;
}

/**
 * 创建租户请求
 */
export interface CreateTenantRequest {
  tenant_name: string;
  tenant_type: TenantType;
  contact_email: string;
  contact_phone?: string;
}

/**
 * 更新租户请求
 */
export interface UpdateTenantRequest {
  tenant_name?: string;
  contact_email?: string;
  contact_phone?: string;
}

/**
 * 租户列表查询参数
 */
export interface TenantListParams {
  status?: TenantStatus;
  tenant_type?: TenantType;
  subscription_tier?: SubscriptionTier;
  page?: number;
  page_size?: number;
}

/**
 * 租户列表响应
 */
export interface TenantListResponse {
  tenants: Tenant[];
  total: number;
  page: number;
  page_size: number;
}

/**
 * 订阅信息
 */
export interface Subscription {
  subscription_id: string;
  tenant_id: string;
  plan_tier: SubscriptionTier;
  status: 'active' | 'past_due' | 'canceled' | 'expired';
  start_date: string;
  end_date: string;
  quota_bots: number;
  quota_messages_per_month: number;
  quota_storage_gb: number;
}

/**
 * 配额信息
 */
export interface QuotaInfo {
  resource_type: string;
  max_limit: number;
  used_count: number;
  remaining: number;
  usage_percent: number;
  alert_level: 'normal' | 'warning' | 'critical' | 'exceeded';
  reset_cycle: string;
}

/**
 * 检查配额请求
 */
export interface CheckQuotaRequest {
  resource_type: string;
  count: number;
}

/**
 * 检查配额响应
 */
export interface CheckQuotaResponse {
  allowed: boolean;
  current_usage: number;
  max_limit: number;
  reason?: string;
}

/**
 * 消费配额请求
 */
export interface ConsumeQuotaRequest {
  resource_type: string;
  count: number;
  description?: string;
}

/**
 * 消费配额响应
 */
export interface ConsumeQuotaResponse {
  success: boolean;
  current_usage: number;
  remaining: number;
}
```

**文件**: `frontend/packages/api-client/src/types/role.ts`

```typescript
// frontend/packages/api-client/src/types/role.ts

/**
 * 角色类型
 */
export type RoleType = 'system' | 'custom';

/**
 * 数据权限范围
 */
export type DataPermissionScope =
  | 'ALL'                 // 所有数据
  | 'DEPARTMENT_SUB'      // 部门及以下
  | 'DEPARTMENT'          // 仅本部门
  | 'OWN'                 // 仅自己的数据
  | 'CUSTOM'              // 自定义
  | 'NONE';               // 无权限

/**
 * 字段权限级别
 */
export type FieldPermissionLevel = 'hidden' | 'readonly' | 'editable';

/**
 * 角色信息
 */
export interface Role {
  role_id: string;
  tenant_id: string;
  role_name: string;
  role_code: string;
  role_type: RoleType;
  description: string;
  created_at: string;
  updated_at: string;
}

/**
 * 创建角色请求
 */
export interface CreateRoleRequest {
  role_name: string;
  role_code: string;
  description?: string;
}

/**
 * 更新角色请求
 */
export interface UpdateRoleRequest {
  role_name?: string;
  description?: string;
}

/**
 * 数据权限
 */
export interface DataPermission {
  permission_id: string;
  role_id: string;
  resource_type: string;
  scope: DataPermissionScope;
  custom_filter?: Record<string, any>;
}

/**
 * 字段权限
 */
export interface FieldPermission {
  permission_id: string;
  role_id: string;
  resource_type: string;
  field_name: string;
  permission_level: FieldPermissionLevel;
}

/**
 * 部门信息
 */
export interface Department {
  department_id: string;
  tenant_id: string;
  parent_id?: string;
  department_name: string;
  description?: string;
  depth: number;
  path: string;
  created_at: string;
  updated_at: string;
}

/**
 * 检查数据权限请求
 */
export interface CheckDataPermissionRequest {
  user_id: string;
  tenant_id: string;
  resource_type: string;
  resource_id: string;
  action: string;
}

/**
 * 检查数据权限响应
 */
export interface CheckDataPermissionResponse {
  allowed: boolean;
  reason?: string;
}
```

**⚠️ 注意事项**:
1. **所有接口都要有JSDoc注释**
2. **类型定义要完整**，不要使用any
3. **枚举类型使用联合类型**（string literals）
4. **可选字段使用?标记**

**📖 开发规范**:
- TypeScript类型完整
- 注释清晰
- 命名规范

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md) - 前端开发规范章节

---

### 任务1.4: API调用函数实现 (Day 3-6)

**文件**: `frontend/packages/api-client/src/hooks/useTenant.ts`

```typescript
// frontend/packages/api-client/src/hooks/useTenant.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { request } from '../http/client';
import type {
  Tenant,
  CreateTenantRequest,
  UpdateTenantRequest,
  TenantListParams,
  TenantListResponse,
  Subscription,
  QuotaInfo,
  CheckQuotaRequest,
  CheckQuotaResponse,
  ConsumeQuotaRequest,
  ConsumeQuotaResponse,
} from '../types/tenant';

// ==================== 租户管理 ====================

/**
 * 获取租户列表
 */
export const useTenantList = (params?: TenantListParams) => {
  return useQuery({
    queryKey: ['tenants', params],
    queryFn: () => request.get<TenantListResponse>('/api/tenants', { params }),
  });
};

/**
 * 获取租户详情
 */
export const useTenant = (tenantId: string) => {
  return useQuery({
    queryKey: ['tenant', tenantId],
    queryFn: () => request.get<Tenant>(`/api/tenants/${tenantId}`),
    enabled: !!tenantId,
  });
};

/**
 * 创建租户
 */
export const useCreateTenant = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTenantRequest) =>
      request.post<Tenant>('/api/tenants', data),
    onSuccess: () => {
      // 刷新租户列表
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
    },
  });
};

/**
 * 更新租户
 */
export const useUpdateTenant = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ tenantId, data }: { tenantId: string; data: UpdateTenantRequest }) =>
      request.put<Tenant>(`/api/tenants/${tenantId}`, data),
    onSuccess: (_, variables) => {
      // 刷新租户详情
      queryClient.invalidateQueries({ queryKey: ['tenant', variables.tenantId] });
    },
  });
};

/**
 * 删除租户
 */
export const useDeleteTenant = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (tenantId: string) =>
      request.delete(`/api/tenants/${tenantId}`),
    onSuccess: () => {
      // 刷新租户列表
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
    },
  });
};

// ==================== 订阅管理 ====================

/**
 * 获取订阅信息
 */
export const useSubscription = (tenantId: string) => {
  return useQuery({
    queryKey: ['subscription', tenantId],
    queryFn: () => request.get<Subscription>(`/api/tenants/${tenantId}/subscription`),
    enabled: !!tenantId,
  });
};

/**
 * 升级订阅
 */
export const useUpgradeSubscription = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ tenantId, data }: { tenantId: string; data: any }) =>
      request.post<Subscription>(`/api/tenants/${tenantId}/upgrade`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['subscription', variables.tenantId] });
    },
  });
};

// ==================== 配额管理 ====================

/**
 * 获取配额列表
 */
export const useQuotas = (tenantId: string) => {
  return useQuery({
    queryKey: ['quotas', tenantId],
    queryFn: () => request.get<QuotaInfo[]>(`/api/tenants/${tenantId}/quotas`),
    enabled: !!tenantId,
  });
};

/**
 * 检查配额
 */
export const useCheckQuota = () => {
  return useMutation({
    mutationFn: ({ tenantId, data }: { tenantId: string; data: CheckQuotaRequest }) =>
      request.post<CheckQuotaResponse>(`/api/tenants/${tenantId}/quotas/check`, data),
  });
};

/**
 * 消费配额
 */
export const useConsumeQuota = () => {
  return useMutation({
    mutationFn: ({ tenantId, data }: { tenantId: string; data: ConsumeQuotaRequest }) =>
      request.post<ConsumeQuotaResponse>(`/api/tenants/${tenantId}/quotas/consume`, data),
  });
};

/**
 * 回滚配额
 */
export const useRollbackQuota = () => {
  return useMutation({
    mutationFn: ({ tenantId, data }: { tenantId: string; data: any }) =>
      request.post(`/api/tenants/${tenantId}/quotas/rollback`, data),
  });
};
```

**文件**: `frontend/packages/api-client/src/hooks/useRole.ts`

```typescript
// frontend/packages/api-client/src/hooks/useRole.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { request } from '../http/client';
import type {
  Role,
  CreateRoleRequest,
  UpdateRoleRequest,
  DataPermission,
  FieldPermission,
  Department,
  CheckDataPermissionRequest,
  CheckDataPermissionResponse,
} from '../types/role';

// ==================== 角色管理 ====================

export const useRoleList = (params?: { tenant_id?: string; role_type?: string }) => {
  return useQuery({
    queryKey: ['roles', params],
    queryFn: () => request.get('/api/roles', { params }),
  });
};

export const useRole = (roleId: string) => {
  return useQuery({
    queryKey: ['role', roleId],
    queryFn: () => request.get<Role>(`/api/roles/${roleId}`),
    enabled: !!roleId,
  });
};

export const useCreateRole = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateRoleRequest) =>
      request.post<Role>('/api/roles', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] });
    },
  });
};

export const useUpdateRole = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ roleId, data }: { roleId: string; data: UpdateRoleRequest }) =>
      request.put<Role>(`/api/roles/${roleId}`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['role', variables.roleId] });
    },
  });
};

export const useDeleteRole = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (roleId: string) =>
      request.delete(`/api/roles/${roleId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] });
    },
  });
};

// ==================== 用户角色 ====================

export const useAssignRole = () => {
  return useMutation({
    mutationFn: ({ userId, data }: { userId: string; data: any }) =>
      request.post(`/api/users/${userId}/roles`, data),
  });
};

export const useRevokeRole = () => {
  return useMutation({
    mutationFn: ({ userId, tenantId, roleId }: { userId: string; tenantId: string; roleId: string }) =>
      request.delete(`/api/users/${userId}/roles/${roleId}?tenant_id=${tenantId}`),
  });
};

export const useUserRoles = (userId: string, tenantId: string) => {
  return useQuery({
    queryKey: ['userRoles', userId, tenantId],
    queryFn: () =>
      request.get(`/api/users/${userId}/roles?tenant_id=${tenantId}`),
    enabled: !!userId && !!tenantId,
  });
};

// ==================== 数据权限 ====================

export const useSetDataPermission = () => {
  return useMutation({
    mutationFn: ({ roleId, data }: { roleId: string; data: any }) =>
      request.post(`/api/roles/${roleId}/data-permissions`, data),
  });
};

export const useCheckDataPermission = () => {
  return useMutation({
    mutationFn: (data: CheckDataPermissionRequest) =>
      request.post<CheckDataPermissionResponse>('/api/permissions/check-data', data),
  });
};

// ==================== 字段权限 ====================

export const useSetFieldPermission = () => {
  return useMutation({
    mutationFn: ({ roleId, data }: { roleId: string; data: any }) =>
      request.post(`/api/roles/${roleId}/field-permissions`, data),
  });
};

export const useFieldPermissions = (userId: string, tenantId: string, resourceType: string) => {
  return useQuery({
    queryKey: ['fieldPermissions', userId, tenantId, resourceType],
    queryFn: () =>
      request.get<FieldPermission[]>(
        `/api/users/${userId}/field-permissions?tenant_id=${tenantId}&resource_type=${resourceType}`
      ),
    enabled: !!userId && !!tenantId && !!resourceType,
  });
};

// ==================== 部门管理 ====================

export const useCreateDepartment = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: any) =>
      request.post<Department>('/api/departments', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['departments'] });
    },
  });
};

export const useDepartmentList = (tenantId: string) => {
  return useQuery({
    queryKey: ['departments', tenantId],
    queryFn: () =>
      request.get(`/api/departments?tenant_id=${tenantId}`),
    enabled: !!tenantId,
  });
};

// ... 其他API Hooks
```

**⚠️ 注意事项**:
1. **所有Hook都以use开头**
2. **使用React Query的缓存机制**，queryKey要唯一
3. **mutation成功后invalidate相关queries**
4. **类型参数使用泛型**

**📖 开发规范**:
- React Hooks规范
- React Query最佳实践
- TypeScript类型完整

**🔗 设计文档链接**:
- [前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md) - 任务1章节

---

## 🎯 Week 2: 前端状态管理 + 路由系统 (P1)

### 任务2.1: Zustand状态管理 (Day 1-2)

**文件**: `frontend/packages/stores/src/authStore.ts`

```typescript
// frontend/packages/stores/src/authStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
  user_id: string;
  username: string;
  email: string;
  avatar?: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  tenantId: string | null;
  isAuthenticated: boolean;

  // Actions
  login: (user: User, token: string) => void;
  logout: () => void;
  setTenantId: (tenantId: string) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      tenantId: null,
      isAuthenticated: false,

      login: (user, token) =>
        set({
          user,
          token,
          isAuthenticated: true,
        }),

      logout: () =>
        set({
          user: null,
          token: null,
          tenantId: null,
          isAuthenticated: false,
        }),

      setTenantId: (tenantId) => set({ tenantId }),
    }),
    {
      name: 'auth-storage', // localStorage key
    }
  )
);
```

**文件**: `frontend/packages/stores/src/tenantStore.ts`

```typescript
// frontend/packages/stores/src/tenantStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Tenant } from '@coze-studio/api-client';

interface TenantState {
  currentTenant: Tenant | null;
  tenants: Tenant[];

  // Actions
  setCurrentTenant: (tenant: Tenant) => void;
  setTenants: (tenants: Tenant[]) => void;
  updateTenant: (tenantId: string, updates: Partial<Tenant>) => void;
}

export const useTenantStore = create<TenantState>()(
  persist(
    (set) => ({
      currentTenant: null,
      tenants: [],

      setCurrentTenant: (tenant) => set({ currentTenant: tenant }),

      setTenants: (tenants) => set({ tenants }),

      updateTenant: (tenantId, updates) =>
        set((state) => ({
          tenants: state.tenants.map((t) =>
            t.tenant_id === tenantId ? { ...t, ...updates } : t
          ),
          currentTenant:
            state.currentTenant?.tenant_id === tenantId
              ? { ...state.currentTenant, ...updates }
              : state.currentTenant,
        })),
    }),
    {
      name: 'tenant-storage',
    }
  )
);
```

**文件**: `frontend/packages/stores/src/permissionStore.ts`

```typescript
// frontend/packages/stores/src/permissionStore.ts
import { create } from 'zustand';

interface PermissionState {
  userRoles: string[];
  permissions: Record<string, boolean>;

  // Actions
  setUserRoles: (roles: string[]) => void;
  setPermissions: (permissions: Record<string, boolean>) => void;
  hasPermission: (resource: string, action: string) => boolean;
}

export const usePermissionStore = create<PermissionState>((set, get) => ({
  userRoles: [],
  permissions: {},

  setUserRoles: (roles) => set({ userRoles: roles }),

  setPermissions: (permissions) => set({ permissions }),

  hasPermission: (resource, action) => {
    const key = `${resource}:${action}`;
    return get().permissions[key] || false;
  },
}));
```

**⚠️ 注意事项**:
1. **使用persist中间件持久化到localStorage**
2. **State更新不可变**，使用展开运算符
3. **类型定义完整**
4. **Action命名清晰**

**📖 开发规范**:
- Zustand最佳实践
- 状态不可变性
- 类型安全

**🔗 设计文档链接**:
- [前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md) - 任务3章节

---

### 任务2.2: React Query配置 (Day 2)

**文件**: `frontend/apps/coze-studio/src/utils/queryClient.ts`

```typescript
// frontend/apps/coze-studio/src/utils/queryClient.ts
import { QueryClient } from '@tanstack/react-query';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5分钟
      cacheTime: 10 * 60 * 1000, // 10分钟
      retry: 1,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 1,
    },
  },
});
```

**集成到App.tsx**:

```typescript
// frontend/apps/coze-studio/src/App.tsx
import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { queryClient } from './utils/queryClient';

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      {/* ... 应用内容 ... */}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
```

---

### 任务2.3: 路由系统 (Day 3)

**文件**: `frontend/apps/coze-studio/src/router/routes.tsx`

```typescript
// frontend/apps/coze-studio/src/router/routes.tsx
import { createBrowserRouter } from 'react-router-dom';
import { ProtectedRoute } from './guards';
import { PermissionRoute } from './guards';

// 页面组件导入
import { LoginPage } from '../pages/LoginPage';
import { TenantList } from '../pages/tenant/TenantList';
import { TenantDetail } from '../pages/tenant/TenantDetail';
import { RoleList } from '../pages/permission/RoleList';
import { QuotaManagement } from '../pages/tenant/QuotaManagement';
import { RoutingManagement } from '../pages/routing/RoutingManagement';

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/',
    element: <ProtectedRoute />,
    children: [
      {
        index: true,
        element: <div>Dashboard</div>,
      },
      // 租户管理路由
      {
        path: 'tenants',
        children: [
          {
            index: true,
            element: <TenantList />,
          },
          {
            path: ':tenantId',
            element: <TenantDetail />,
          },
          {
            path: ':tenantId/quotas',
            element: <QuotaManagement />,
          },
        ],
      },
      // 权限管理路由
      {
        path: 'permissions',
        element: <PermissionRoute requiredPermission="permission:read" />,
        children: [
          {
            path: 'roles',
            element: <RoleList />,
          },
        ],
      },
      // 路由管理路由
      {
        path: 'routing',
        element: <PermissionRoute requiredPermission="routing:manage" />,
        children: [
          {
            path: 'rules',
            element: <RoutingManagement />,
          },
        ],
      },
    ],
  },
]);
```

**文件**: `frontend/apps/coze-studio/src/router/guards.tsx`

```typescript
// frontend/apps/coze-studio/src/router/guards.tsx
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuthStore } from '@coze-studio/stores';

/**
 * 需要登录的路由守卫
 */
export const ProtectedRoute = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const location = useLocation();

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <Outlet />;
};

/**
 * 需要特定权限的路由守卫
 */
interface PermissionRouteProps {
  children?: React.ReactNode;
  requiredPermission: string;
}

export const PermissionRoute = ({ children, requiredPermission }: PermissionRouteProps) => {
  const { hasPermission } = usePermissionStore();
  const hasAccess = hasPermission(requiredPermission);

  if (!hasAccess) {
    return <div>权限不足</div>;
  }

  return <>{children}</>;
};

/**
 * 租户级别路由守卫
 */
interface TenantRouteProps {
  children?: React.ReactNode;
}

export const TenantRoute = ({ children }: TenantRouteProps) => {
  const tenantId = useAuthStore((state) => state.tenantId);

  if (!tenantId) {
    return <div>请先选择租户</div>;
  }

  return <>{children}</>;
};
```

**⚠️ 注意事项**:
1. **路由守卫使用Outlet渲染子路由**
2. **权限检查使用store**
3. **未授权跳转到登录页**
4. **保持路由层级不超过3层**

**📖 开发规范**:
- React Router v6规范
- 路由守卫清晰
- 权限检查统一

**🔗 设计文档链接**:
- [前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md) - 任务2章节

---

## 🎯 Week 3-4: 前端管理页面 (P1)

### 任务3.1: 配额管理页面 (Day 1-3, 5人天)

**文件**: `frontend/apps/coze-studio/src/pages/tenant/QuotaManagement.tsx`

```typescript
// frontend/apps/coze-studio/src/pages/tenant/QuotaManagement.tsx
import React, { useState, useEffect } from 'react';
import { Table, Card, Button, Progress, Alert, Modal } from '@douyinfe/semi-ui';
import { useQuotas, useUpgradeSubscription } from '@coze-studio/api-client';
import { useTenantContext } from '../../context/tenant';
import type { QuotaInfo } from '@coze-studio/api-client';

export const QuotaManagement: React.FC = () => {
  const { tenant } = useTenantContext();
  const { data: quotas, loading } = useQuotas(tenant?.tenant_id || '');
  const upgradeSubscription = useUpgradeSubscription();

  const [showUpgradeModal, setShowUpgradeModal] = useState(false);

  const getAlertColor = (level: string) => {
    switch (level) {
      case 'normal':
        return 'green';
      case 'warning':
        return 'orange';
      case 'critical':
      case 'exceeded':
        return 'red';
      default:
        return 'green';
    }
  };

  const handleUpgrade = () => {
    upgradeSubscription.mutate({
      tenantId: tenant?.tenant_id || '',
      data: { plan_tier: 'pro' },
    });
    setShowUpgradeModal(false);
  };

  if (loading) {
    return <div>加载中...</div>;
  }

  const hasCritical = quotas?.some((q) => q.alert_level === 'critical' || q.alert_level === 'exceeded');

  return (
    <div className="quota-management">
      <Card title="配额使用情况">
        {hasCritical && (
          <Alert
            type="warning"
            message="部分资源配额即将用完或已超限，建议升级订阅或清理数据"
            description={quotas
              ?.filter((q) => q.alert_level === 'critical' || q.alert_level === 'exceeded')
              .map((q) => q.resource_type)
              .join(', ')}
          />
        )}

        <Table
          dataSource={quotas || []}
          columns={[
            { title: '资源类型', dataIndex: 'resource_type', key: 'resource_type' },
            { title: '已使用', dataIndex: 'used_count', key: 'used_count' },
            { title: '配额上限', dataIndex: 'max_limit', key: 'max_limit' },
            {
              title: '使用率',
              dataIndex: 'usage_percent',
              key: 'usage_percent',
              render: (percent: number, record: QuotaInfo) => (
                <Progress
                  percent={percent}
                  showInfo={true}
                  stroke={getAlertColor(record.alert_level)}
                />
              ),
            },
            { title: '状态', dataIndex: 'alert_level', key: 'alert_level' },
            { title: '重置周期', dataIndex: 'reset_cycle', key: 'reset_cycle' },
          ]}
        />

        <div style={{ marginTop: 16 }}>
          <Button theme="solid" type="primary" onClick={() => setShowUpgradeModal(true)}>
            升级订阅
          </Button>
        </div>
      </Card>

      <Modal
        title="升级订阅"
        visible={showUpgradeModal}
        onOk={handleUpgrade}
        onCancel={() => setShowUpgradeModal(false)}
      >
        <p>确定要升级到Pro版本吗？</p>
      </Modal>
    </div>
  );
};
```

**⚠️ 注意事项**:
1. **使用Semi Design组件**
2. **使用Tailwind CSS工具类**
3. **配额不足显示Alert警告**
4. **升级订阅要有Modal确认**
5. **Loading状态显示**

**📖 开发规范**:
- 组件职责单一
- Props接口完整
- 样式使用Tailwind

**🔗 设计文档链接**:
- [前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md) - 任务6章节

---

### 任务3.2: 权限管理页面 (Day 4-7, 6人天)

**文件**: `frontend/apps/coze-studio/src/pages/permission/RoleList.tsx`

```typescript
// frontend/apps/coze-studio/src/pages/permission/RoleList.tsx
import React, { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Space, Popconfirm } from '@douyinfe/semi-ui';
import { useRoleList, useCreateRole, useUpdateRole, useDeleteRole } from '@coze-studio/api-client';
import type { Role, CreateRoleRequest, UpdateRoleRequest } from '@coze-studio/api-client';

export const RoleList: React.FC = () => {
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [currentRole, setCurrentRole] = useState<Role | null>(null);

  const { data, loading, refetch } = useRoleList({});
  const createRole = useCreateRole();
  const updateRole = useUpdateRole();
  const deleteRole = useDeleteRole();

  const [form] = Form.useForm();

  const handleCreate = async (values: any) => {
    await createRole.mutateAsync(values as CreateRoleRequest);
    setShowCreateModal(false);
    form.reset();
    refetch();
  };

  const handleEdit = async (values: any) => {
    if (currentRole) {
      await updateRole.mutateAsync({
        roleId: currentRole.role_id,
        data: values as UpdateRoleRequest,
      });
      setShowEditModal(false);
      refetch();
    }
  };

  const handleDelete = async (roleId: string) => {
    await deleteRole.mutateAsync(roleId);
    refetch();
  };

  const columns = [
    { title: '角色名称', dataIndex: 'role_name', key: 'role_name' },
    { title: '角色代码', dataIndex: 'role_code', key: 'role_code' },
    { title: '角色类型', dataIndex: 'role_type', key: 'role_type' },
    { title: '描述', dataIndex: 'description', key: 'description' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Role) => (
        <Space>
          <Button
            size="small"
            onClick={() => {
              setCurrentRole(record);
              setShowEditModal(true);
              form.setValues(record);
            }}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个角色吗？"
            onConfirm={() => handleDelete(record.role_id)}
          >
            <Button size="small" type="danger">
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="role-list">
      <div className="mb-4 flex justify-between">
        <h2 className="text-2xl font-bold">角色管理</h2>
        <Button theme="solid" type="primary" onClick={() => setShowCreateModal(true)}>
          创建角色
        </Button>
      </div>

      <Table
        loading={loading}
        dataSource={data?.roles || []}
        columns={columns}
        rowKey="role_id"
      />

      {/* 创建角色Modal */}
      <Modal
        title="创建角色"
        visible={showCreateModal}
        onOk={() => form.submit()}
        onCancel={() => setShowCreateModal(false)}
      >
        <Form
          form={form}
          onSubmit={handleCreate}
          labelPosition="left"
          labelWidth="100px"
        >
          <Form.Input
            field="role_name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          />
          <Form.Input
            field="role_code"
            label="角色代码"
            rules={[{ required: true, message: '请输入角色代码' }]}
          />
          <Form.TextArea
            field="description"
            label="描述"
            placeholder="请输入角色描述"
          />
        </Form>
      </Modal>

      {/* 编辑角色Modal */}
      <Modal
        title="编辑角色"
        visible={showEditModal}
        onOk={() => form.submit()}
        onCancel={() => setShowEditModal(false)}
      >
        <Form
          form={form}
          onSubmit={handleEdit}
          labelPosition="left"
          labelWidth="100px"
        >
          <Form.Input
            field="role_name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          />
          <Form.TextArea
            field="description"
            label="描述"
            placeholder="请输入角色描述"
          />
        </Form>
      </Modal>
    </div>
  );
};
```

**⚠️ 注意事项**:
1. **CRUD操作完整**：创建、读取、更新、删除
2. **删除操作要有Popconfirm确认**
3. **表单验证完整**
4. **操作成功后刷新列表**
5. **Loading和Error状态处理**

**📖 开发规范**:
- 组件职责单一
- 表单验证完整
- 用户操作友好

**🔗 设计文档链接**:
- [前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md) - 任务6章节

---

### 任务3.3: 路由管理页面 (Day 8-11, 6人天)

**文件**: `frontend/apps/coze-studio/src/pages/routing/RoutingManagement.tsx`

```typescript
// frontend/apps/coze-studio/src/pages/routing/RoutingManagement.tsx
import React, { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Tag } from '@douyinfe/semi-ui';
import { useRoutingRules } from '@coze-studio/api-client';
import type { RoutingRule } from '@coze-studio/api-client';

export const RoutingManagement: React.FC = () => {
  const [showCreateModal, setShowCreateModal] = useState(false);
  const { data, loading, refetch } = useRoutingRules({});

  const columns = [
    { title: '规则名称', dataIndex: 'rule_name', key: 'rule_name' },
    {
      title: '规则类型',
      dataIndex: 'rule_type',
      key: 'rule_type',
      render: (type: string) => {
        const typeMap: Record<string, { text: string; color: string }> = {
          keyword: { text: '关键词', color: 'blue' },
          regex: { text: '正则', color: 'green' },
          intent: { text: '意图', color: 'orange' },
          category: { text: '分类', color: 'purple' },
        };
        const config = typeMap[type] || { text: type, color: 'default' };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    { title: 'Bot ID', dataIndex: 'bot_id', key: 'bot_id' },
    { title: '优先级', dataIndex: 'priority', key: 'priority' },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (active: boolean) => (
        <Tag color={active ? 'green' : 'red'}>
          {active ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: RoutingRule) => (
        <Space>
          <Button size="small">编辑</Button>
          <Button size="small" type="danger">
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="routing-management">
      <div className="mb-4 flex justify-between">
        <h2 className="text-2xl font-bold">路由规则管理</h2>
        <Button theme="solid" type="primary" onClick={() => setShowCreateModal(true)}>
          创建规则
        </Button>
      </div>

      <Table
        loading={loading}
        dataSource={data?.rules || []}
        columns={columns}
        rowKey="rule_id"
      />

      {/* 创建规则Modal */}
      <Modal
        title="创建路由规则"
        visible={showCreateModal}
        onOk={() => {/* TODO */}}
        onCancel={() => setShowCreateModal(false)}
      >
        <Form labelPosition="left" labelWidth="100px">
          <Form.Input
            field="rule_name"
            label="规则名称"
            rules={[{ required: true, message: '请输入规则名称' }]}
          />
          <Form.Select
            field="rule_type"
            label="规则类型"
            rules={[{ required: true, message: '请选择规则类型' }]}
          >
            <Select.Option value="keyword">关键词</Select.Option>
            <Select.Option value="regex">正则</Select.Option>
            <Select.Option value="intent">意图</Select.Option>
            <Select.Option value="category">分类</Select.Option>
          </Form.Select>
          <Form.Input field="bot_id" label="Bot ID" />
          <Form.InputNumber field="priority" label="优先级" />
        </Form>
      </Modal>
    </div>
  );
};
```

**⚠️ 注意事项**:
1. **表格列使用Tag显示状态**
2. **操作按钮使用Space分隔**
3. **表单验证完整**
4. **规则类型使用下拉选择**

**📖 开发规范**:
- UI组件规范
- 数据展示清晰
- 交互友好

**🔗 设计文档链接**:
- [前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md) - 任务6章节

---

## 🎯 Week 5: UI组件补充 + 单元测试 (P2)

### 任务4.1: UI组件补充 (Day 1-5, 7人天)

由于篇幅限制，这里只列出需要补充的组件清单：

**必补组件**:
1. **DatePicker** (2-3天) - 日期选择器
2. **Upload** (2-3天) - 文件上传
3. **Dropdown** (1天) - 下拉菜单
4. **Tooltip** (1天) - 工具提示
5. **Alert/Message** (1-2天) - 警告提示
6. **Progress** (1天) - 进度条
7. **Spin** (1天) - 加载中

**详细实施步骤**请参考：
[前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md) - 任务4章节

---

### 任务4.2: 前端单元测试 (Day 6-7)

**测试覆盖率要求**:
- UI组件测试 ≥ 80%
- 业务组件测试 ≥ 70%
- 页面组件测试 ≥ 60%

**测试工具**: Vitest + @testing-library/react

**示例测试**:

```typescript
// QuotaManagement.test.tsx
import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { QuotaManagement } from './QuotaManagement';

describe('QuotaManagement', () => {
  it('should render quota list', () => {
    render(<QuotaManagement />);
    expect(screen.getByText('配额使用情况')).toBeInTheDocument();
  });

  it('should show alert when quota exceeded', () => {
    // TODO: 测试配额超限时的Alert显示
  });
});
```

**⚠️ 注意事项**:
1. **测试关注行为而非实现**
2. **使用mock数据**
3. **异步操作使用waitFor**
4. **测试用户交互**

**📖 开发规范**:
- Vitest最佳实践
- Testing Library规范
- 覆盖率要求

**🔗 设计文档链接**:
- [前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md) - 任务7章节

---

## 📊 关键文档链接汇总

### 必读文档（开发前必读）

1. **[前端开发遗漏工作-立即执行清单.md](../前端开发遗漏工作-立即执行清单.md)** ⭐⭐⭐
   - 完整的前端开发任务
   - 三周快速路径

2. **[ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md)** ⭐⭐⭐
   - 前端开发规范章节
   - React + TypeScript规范

3. **[ZKER-全局一致性检查清单_v1.0.md](../ZKER-全局一致性检查清单_v1.0.md)** ⭐⭐
   - L1个人检查清单
   - 前端代码一致性检查

4. **[前端组件库完整索引_v1.0.md](../前端组件库完整索引_v1.0.md)** ⭐
   - 现有组件索引
   - 组件使用示例

### 参考文档

5. **[API接口文档_租户系统.md](../API接口文档_租户系统.md)** - 租户API
6. **[API接口文档_用户管理RBAC.md](../API接口文档_用户管理RBAC.md)** - RBAC API
7. **[API接口文档_智能路由引擎.md](../API接口文档_智能路由引擎.md)** - 路由API

---

## ⚠️ 核心注意事项（重申）

### 1. 代码隔离原则

**你的职责范围**:
- ✅ 前端代码全部（TypeScript + React）
- ✅ API调用层
- ✅ 状态管理
- ✅ 路由系统
- ✅ UI组件

**不要触碰**:
- ❌ 后端代码（Go）
- ❌ 数据库（SQL）
- ❌ CI/CD配置
- ❌ 监控配置

### 2. 前端开发红线（重申）

**必须遵守**:
1. 所有组件使用TypeScript，禁止any
2. 函数式组件 + Hooks，不使用类组件
3. 样式使用Tailwind CSS
4. 状态管理使用Zustand + React Query
5. API调用统一，不直接用fetch/axios
6. 所有组件有类型定义

**禁止行为**:
1. 使用类组件
2. 使用any类型
3. 直接使用fetch/axios
4. 硬编码API URL
5. 不写TypeScript类型

### 3. Git提交规范

**Commit Message格式**:
```
feat(frontend): implement quota management page

- Add QuotaManagement component with Semi Design
- Integrate useQuotas hook for data fetching
- Add upgrade subscription modal
- Display quota usage with Progress bars

Refs: #123
```

---

## 📅 每日工作检查清单

### 开发前
- [ ] 阅读相关设计文档
- [ ] 理解API契约和数据模型
- [ ] 确认组件职责
- [ ] 拉取最新代码

### 开发中
- [ ] 遵循SOLID、KISS、DRY、YAGNI原则
- [ ] TypeScript类型完整
- [ ] 组件职责单一（<300行）
- [ ] 样式使用Tailwind
- [ ] 边开发边写测试

### 提交前
- [ ] ESLint检查通过
- [ ] TypeScript编译通过
- [ ] 代码符合规范
- [ ] 自我Code Review
- [ ] 单元测试通过

---

## 🎯 成功标准

### Week 1-2结束时（API层）
- [ ] 55+个API函数全部实现
- [ ] 所有TypeScript类型定义完成
- [ ] HTTP客户端封装完成
- [ ] 错误处理统一

### Week 2结束时（状态+路由）
- [ ] 3个Zustand Store创建完成
- [ ] React Query配置完成
- [ ] 路由系统配置完成
- [ ] 路由守卫实现完成

### Week 3-4结束时（管理页面）
- [ ] 4个管理页面全部实现
- [ ] 所有页面集成API调用
- [ ] 表格支持分页、筛选、排序
- [ ] 用户操作有确认对话框

### Week 5结束时（组件+测试）
- [ ] 7个UI组件全部实现
- [ ] 单元测试覆盖率 ≥ 70%
- [ ] 所有测试通过
- [ ] 文档完整交付

---

**🎉 你是用户可见界面开发的主力，前端质量直接影响用户体验！**
**遇到问题随时查阅设计文档或与团队讨论UI/UX设计。**
