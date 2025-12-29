// frontend/packages/api-client/src/hooks/index.ts

/**
 * API Client Hooks - 统一导出
 */

// 租户管理
export * from './useTenantList';
export * from './useTenant';
export * from './useCreateTenant';
export * from './useUpdateTenant';
export * from './useDeleteTenant';
export * from './useTenantStats';

// 角色管理
export * from './useRoleList';
export * from './useRole';
export * from './useCreateRole';
export * from './useUpdateRole';
export * from './useDeleteRole';

// 权限管理
export * from './usePermissionList';
export * from './useCheckPermission';

// 路由管理
export * from './useRoutingRules';
export * from './useRoutingRule';
export * from './useCreateRoutingRule';
export * from './useUpdateRoutingRule';
export * from './useDeleteRoutingRule';

// 意图匹配器
export * from './useIntentMatchers';
export * from './useIntentMatcher';
export * from './useCreateIntentMatcher';
export * from './useUpdateIntentMatcher';
export * from './useDeleteIntentMatcher';

// 配额管理
export * from './useQuota';
export * from './useQuotaUsage';
export * from './useQuotaUsageDetails';
