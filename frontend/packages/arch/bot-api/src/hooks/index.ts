/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * ZKER 企业级 React Hooks 统一导出
 *
 * 提供租户、配额、权限管理的完整 React Hooks
 *
 * @module enterprise-hooks
 * @example
 * ```typescript
 * import { useTenant, useQuotas, useRoles } from '@coze-arch/bot-api';
 *
 * function TenantPage() {
 *   const { tenant, loading, error } = useTenant('tenant_001');
 *
 *   if (loading) return <Spin />;
 *   if (error) return <Alert message={error.message} />;
 *
 *   return <div>{tenant.tenant_name}</div>;
 * }
 * ```
 */

// ==================== 租户管理 Hooks ====================
export {
  useTenant,
  useTenants,
  useCreateTenant,
  useUpdateTenant,
  useDeleteTenant,
  useSubscription,
  useUpgradeSubscription,
  useTenantStats,
} from './useTenant';

// ==================== 配额管理 Hooks ====================
export {
  useQuotas,
  useQuotaCheck,
  useQuotaUsage,
  useUpdateQuotaLimit,
  useQuotaWarning,
  useQuotaStats,
} from './useQuota';

// ==================== 权限管理 Hooks ====================
export {
  useRoles,
  useRole,
  useCreateRole,
  useUpdateRole,
  useDeleteRole,
  useAssignRole,
  useRevokeRole,
  useDepartments,
  useDataPermission,
  useFieldPermission,
  usePermissionStats,
} from './usePermission';
