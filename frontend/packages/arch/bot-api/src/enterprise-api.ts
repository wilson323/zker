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
 * ZKER 企业级 API 统一导出
 *
 * 提供租户管理、配额管理、权限管理等企业级功能的 API 调用接口
 *
 * @module enterprise-api
 * @example
 * ```typescript
 * import { tenantApi, quotaApi, permissionApi } from '@coze-arch/bot-api';
 *
 * // 获取租户列表
 * const tenants = await tenantApi.listTenants({
 *   page: 1,
 *   page_size: 20,
 * });
 *
 * // 检查配额
 * const quota = await quotaApi.checkQuota({
 *   tenant_id: 'xxx',
 *   resource_type: ResourceType.BOTS,
 *   required_count: 1,
 * });
 *
 * // 获取用户角色
 * const roles = await permissionApi.getUserRoles({
 *   tenant_id: 'xxx',
 *   user_id: 'xxx',
 * });
 * ```
 */

// ==================== 租户管理 API ====================
export { tenantApi } from './tenant-api';
export type {
  // 租户管理类型
  CreateTenantRequest,
  CreateTenantResponse,
  GetTenantRequest,
  GetTenantResponse,
  UpdateTenantRequest,
  UpdateTenantResponse,
  ListTenantsRequest,
  ListTenantsResponse,
  DeleteTenantRequest,
  DeleteTenantResponse,
  // 订阅管理类型
  GetSubscriptionRequest,
  GetSubscriptionResponse,
  UpdateSubscriptionRequest,
  UpdateSubscriptionResponse,
  UpgradeSubscriptionRequest,
  UpgradeSubscriptionResponse,
  // 枚举和接口
  Tenant,
  Subscription,
  TenantStatus,
  SubscriptionTier,
  BillingCycle,
  SubscriptionStatus,
  PlanQuota,
} from './types/tenant.types';

// ==================== 配额管理 API ====================
export { quotaApi } from './quota-api';
export type {
  GetQuotasRequest,
  GetQuotasResponse,
  CheckQuotaRequest,
  CheckQuotaResponse,
  UpdateQuotaLimitRequest,
  UpdateQuotaLimitResponse,
  GetQuotaUsageRequest,
  GetQuotaUsageResponse,
  // 枚举和接口
  Quota,
  QuotaUsageDetail,
  UsageHistoryItem,
  ResourceType,
  QuotaPeriod,
} from './types/quota.types';

// ==================== 权限管理 API ====================
export { permissionApi } from './permission-api';
export type {
  // 角色管理类型
  CreateRoleRequest,
  CreateRoleResponse,
  GetRoleRequest,
  GetRoleResponse,
  UpdateRoleRequest,
  UpdateRoleResponse,
  ListRolesRequest,
  ListRolesResponse,
  DeleteRoleRequest,
  DeleteRoleResponse,
  AssignRoleRequest,
  AssignRoleResponse,
  RevokeRoleRequest,
  RevokeRoleResponse,
  GetUserRolesRequest,
  GetUserRolesResponse,
  // 部门管理类型
  CreateDepartmentRequest,
  CreateDepartmentResponse,
  GetDepartmentRequest,
  GetDepartmentResponse,
  UpdateDepartmentRequest,
  UpdateDepartmentResponse,
  ListDepartmentsRequest,
  ListDepartmentsResponse,
  DeleteDepartmentRequest,
  DeleteDepartmentResponse,
  AddDepartmentMemberRequest,
  AddDepartmentMemberResponse,
  RemoveDepartmentMemberRequest,
  RemoveDepartmentMemberResponse,
  // 权限检查类型
  CheckDataPermissionRequest,
  CheckDataPermissionResponse,
  CheckFieldPermissionRequest,
  CheckFieldPermissionResponse,
  // 枚举和接口
  Role,
  DataPermission,
  FieldPermission,
  UserRole,
  Department,
  DepartmentMember,
  PermissionScope,
  FieldPermissionLevel,
} from './types/permission.types';
