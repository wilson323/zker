// frontend/packages/arch/business-components/src/index.ts

/**
 * ZKER 业务组件库
 *
 * 提供租户、权限、配额、订阅相关的业务组件
 */

// 租户相关
export { TenantSelector } from './components/TenantSelector';
export type { TenantSelectorProps, Tenant } from './components/TenantSelector';

export { TenantStats } from './components/TenantStats';
export type { TenantStatsProps } from './components/TenantStats';

// 权限相关
export { PermissionTree } from './components/PermissionTree';
export type { PermissionTreeProps, PermissionNode } from './components/PermissionTree';

export { RoleMemberSelector } from './components/RoleMemberSelector';
export type { RoleMemberSelectorProps } from './components/RoleMemberSelector';

// 配额相关
export { QuotaIndicator } from './components/QuotaIndicator';
export type { QuotaIndicatorProps } from './components/QuotaIndicator';

export { QuotaEditor } from './components/QuotaEditor';
export type { QuotaEditorProps, QuotaLimit } from './components/QuotaEditor';

// 订阅相关
export { SubscriptionSelector } from './components/SubscriptionSelector';
export type { SubscriptionSelectorProps, SubscriptionTierOption, SubscriptionTier } from './components/SubscriptionSelector';

export { DataScopeSelector } from './components/DataScopeSelector';
export type { DataScopeSelectorProps, DataScope, DataScopeOption } from './components/DataScopeSelector';
