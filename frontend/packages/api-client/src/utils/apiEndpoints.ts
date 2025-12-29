// frontend/packages/api-client/src/utils/apiEndpoints.ts

/**
 * API 端点常量
 */
export const API_ENDPOINTS = {
  // 租户管理
  TENANT: {
    LIST: '/api/v1/tenants',
    GET: (tenantId: string) => `/api/v1/tenants/${tenantId}`,
    CREATE: '/api/v1/tenants',
    UPDATE: (tenantId: string) => `/api/v1/tenants/${tenantId}`,
    DELETE: (tenantId: string) => `/api/v1/tenants/${tenantId}`,
    STATS: (tenantId: string) => `/api/v1/tenants/${tenantId}/stats`,
  },

  // 角色管理
  ROLE: {
    LIST: '/api/v1/roles',
    GET: (roleId: string) => `/api/v1/roles/${roleId}`,
    CREATE: '/api/v1/roles',
    UPDATE: (roleId: string) => `/api/v1/roles/${roleId}`,
    DELETE: (roleId: string) => `/api/v1/roles/${roleId}`,
    MEMBERS: (roleId: string) => `/api/v1/roles/${roleId}/members`,
    ADD_MEMBER: (roleId: string) => `/api/v1/roles/${roleId}/members`,
    REMOVE_MEMBER: (roleId: string, userId: string) =>
      `/api/v1/roles/${roleId}/members/${userId}`,
  },

  // 权限管理
  PERMISSION: {
    LIST: '/api/v1/permissions',
    CHECK: '/api/v1/permissions/check',
    DATA_PERMISSIONS: '/api/v1/permissions/data',
    FIELD_PERMISSIONS: '/api/v1/permissions/field',
  },

  // 路由规则管理
  ROUTING_RULE: {
    LIST: '/api/v1/routing/rules',
    GET: (ruleId: string) => `/api/v1/routing/rules/${ruleId}`,
    CREATE: '/api/v1/routing/rules',
    UPDATE: (ruleId: string) => `/api/v1/routing/rules/${ruleId}`,
    DELETE: (ruleId: string) => `/api/v1/routing/rules/${ruleId}`,
    TEST: '/api/v1/routing/rules/test',
    STATS: '/api/v1/routing/rules/stats',
    ENABLE: (ruleId: string) => `/api/v1/routing/rules/${ruleId}/enable`,
    DISABLE: (ruleId: string) => `/api/v1/routing/rules/${ruleId}/disable`,
  },

  // 意图匹配器管理
  INTENT_MATCHER: {
    LIST: '/api/v1/routing/matchers',
    GET: (matcherId: string) => `/api/v1/routing/matchers/${matcherId}`,
    CREATE: '/api/v1/routing/matchers',
    UPDATE: (matcherId: string) => `/api/v1/routing/matchers/${matcherId}`,
    DELETE: (matcherId: string) => `/api/v1/routing/matchers/${matcherId}`,
    TEST: '/api/v1/routing/matchers/test',
  },

  // 配额管理
  QUOTA: {
    LIST: '/api/v1/quotas',
    GET: (quotaId: string) => `/api/v1/quotas/${quotaId}`,
    UPDATE: (quotaId: string) => `/api/v1/quotas/${quotaId}`,
    USAGE: (tenantId: string) => `/api/v1/quotas/usage/${tenantId}`,
    USAGE_DETAILS: (tenantId: string) => `/api/v1/quotas/usage/${tenantId}/details`,
    HISTORY: '/api/v1/quotas/history',
    ALERT_RULES: (tenantId: string) => `/api/v1/quotas/alerts/${tenantId}`,
  },

  // 订阅管理
  SUBSCRIPTION: {
    GET: (tenantId: string) => `/api/v1/subscriptions/${tenantId}`,
    UPDATE: (subscriptionId: string) => `/api/v1/subscriptions/${subscriptionId}`,
    UPGRADE: (subscriptionId: string) => `/api/v1/subscriptions/${subscriptionId}/upgrade`,
    DOWNGRADE: (subscriptionId: string) => `/api/v1/subscriptions/${subscriptionId}/downgrade`,
    CANCEL: (subscriptionId: string) => `/api/v1/subscriptions/${subscriptionId}/cancel`,
    PLANS: '/api/v1/subscriptions/plans',
    INVOICE_HISTORY: (subscriptionId: string) =>
      `/api/v1/subscriptions/${subscriptionId}/invoices`,
  },

  // 用户管理
  USER: {
    PROFILE: '/api/v1/users/profile',
    UPDATE_PROFILE: '/api/v1/users/profile',
    CHANGE_PASSWORD: '/api/v1/users/password',
    TENANTS: '/api/v1/users/tenants',
    SWITCH_TENANT: '/api/v1/users/switch-tenant',
  },

  // 组织管理
  ORGANIZATION: {
    TREE: (tenantId: string) => `/api/v1/organizations/${tenantId}/tree`,
    MEMBERS: (tenantId: string) => `/api/v1/organizations/${tenantId}/members`,
    DEPARTMENTS: (tenantId: string) => `/api/v1/organizations/${tenantId}/departments`,
    POSITIONS: (tenantId: string) => `/api/v1/organizations/${tenantId}/positions`,
  },
} as const;
