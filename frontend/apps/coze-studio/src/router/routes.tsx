// frontend/apps/coze-studio/src/router/routes.tsx

import { RouteObject } from 'react-router-dom';

/**
 * 应用路由配置
 */
export const routes: RouteObject[] = [
  {
    path: '/',
    lazy: async () => {
      const { default: Layout } = await import('../components/Layout/Layout');
      return { Component: Layout };
    },
    children: [
      // 首页
      {
        index: true,
        lazy: async () => {
          const { default: Home } = await import('../pages/Home/Home');
          return { Component: Home };
        },
      },

      // 租户管理
      {
        path: 'tenants',
        lazy: async () => {
          const { default: TenantList } = await import(
            '@coze-studio/studio/src/pages/tenant/TenantList/TenantList'
          );
          return { Component: TenantList };
        },
      },
      {
        path: 'tenants/:tenantId',
        lazy: async () => {
          const { default: TenantDetail } = await import(
            '@coze-studio/studio/src/pages/tenant/TenantDetail/TenantDetail'
          );
          return { Component: TenantDetail };
        },
      },

      // 权限管理 - 角色
      {
        path: 'permissions/roles',
        lazy: async () => {
          const { default: RoleList } = await import(
            '@coze-studio/studio/src/pages/permission/RoleList/RoleList'
          );
          return { Component: RoleList };
        },
      },
      {
        path: 'permissions/roles/:roleId',
        lazy: async () => {
          const { default: RoleDetail } = await import(
            '@coze-studio/studio/src/pages/permission/RoleDetail/RoleDetail'
          );
          return { Component: RoleDetail };
        },
      },

      // 路由配置
      {
        path: 'routing/rules',
        lazy: async () => {
          const { default: RoutingRules } = await import(
            '@coze-studio/studio/src/pages/routing/RoutingRules/RoutingRules'
          );
          return { Component: RoutingRules };
        },
      },
      {
        path: 'routing/intents',
        lazy: async () => {
          const { default: IntentMatcher } = await import(
            '@coze-studio/studio/src/pages/routing/IntentMatcher/IntentMatcher'
          );
          return { Component: IntentMatcher };
        },
      },

      // 设置
      {
        path: 'settings/subscription',
        lazy: async () => {
          const { default: SubscriptionManagement } = await import(
            '@coze-studio/studio/src/pages/settings/SubscriptionManagement/SubscriptionManagement'
          );
          return { Component: SubscriptionManagement };
        },
      },
      {
        path: 'settings/quotas',
        lazy: async () => {
          const { default: QuotaManagement } = await import(
            '@coze-studio/studio/src/pages/settings/QuotaManagement/QuotaManagement'
          );
          return { Component: QuotaManagement };
        },
      },
      {
        path: 'settings/profile',
        lazy: async () => {
          const { default: UserProfile } = await import(
            '@coze-studio/studio/src/pages/settings/UserProfile/UserProfile'
          );
          return { Component: UserProfile };
        },
      },
      {
        path: 'settings/organization',
        lazy: async () => {
          const { default: OrganizationManagement } = await import(
            '@coze-studio/studio/src/pages/settings/OrganizationManagement/OrganizationManagement'
          );
          return { Component: OrganizationManagement };
        },
      },

      // 404页面
      {
        path: '*',
        lazy: async () => {
          const { default: NotFound } = await import('../pages/NotFound/NotFound');
          return { Component: NotFound };
        },
      },
    ],
  },

  // 登录页（无布局）
  {
    path: '/login',
    lazy: async () => {
      const { default: Login } = await import('../pages/Login/Login');
      return { Component: Login };
    },
  },
];
