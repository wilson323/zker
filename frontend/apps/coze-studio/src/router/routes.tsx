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
        path: 'tenants/create',
        lazy: async () => {
          const { default: TenantCreate } = await import(
            '@coze-studio/studio/src/pages/tenant/TenantCreate/TenantCreate'
          );
          return { Component: TenantCreate };
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
      {
        path: 'tenants/:tenantId/edit',
        lazy: async () => {
          const { default: TenantEdit } = await import(
            '@coze-studio/studio/src/pages/tenant/TenantEdit/TenantEdit'
          );
          return { Component: TenantEdit };
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

      // 计费管理
      {
        path: 'billing/token-usage',
        lazy: async () => {
          const { default: TokenUsagePage } = await import(
            '../pages/billing/TokenUsagePage'
          );
          return { Component: TokenUsagePage };
        },
      },
      {
        path: 'billing/budget',
        lazy: async () => {
          const { default: BudgetManagementPage } = await import(
            '../pages/billing/BudgetManagementPage'
          );
          return { Component: BudgetManagementPage };
        },
      },

      // 组织管理
      {
        path: 'org/virtual',
        lazy: async () => {
          const { default: VirtualOrganizationPage } = await import(
            '@coze-studio/studio/src/pages/org/VirtualOrganization'
          );
          return { Component: VirtualOrganizationPage };
        },
      },
      {
        path: 'org/delegate',
        lazy: async () => {
          const { default: OrgDelegatePage } = await import(
            '@coze-studio/studio/src/pages/org/OrgDelegate'
          );
          return { Component: OrgDelegatePage };
        },
      },

      // 审计日志
      {
        path: 'audit/logs',
        lazy: async () => {
          const { default: AuditLogsPage } = await import(
            '@coze-studio/studio/src/pages/audit/AuditLogs'
          );
          return { Component: AuditLogsPage };
        },
      },

      // 审核工作台
      {
        path: 'audit/queue',
        lazy: async () => {
          const { default: ReviewQueue } = await import(
            '@coze-studio/studio/src/pages/audit/ReviewQueue/ReviewQueue'
          );
          return { Component: ReviewQueue };
        },
      },
      {
        path: 'audit/tasks/:taskId',
        lazy: async () => {
          const { default: TaskDetail } = await import(
            '@coze-studio/studio/src/pages/audit/TaskDetail/TaskDetail'
          );
          return { Component: TaskDetail };
        },
      },
      {
        path: 'audit/history',
        lazy: async () => {
          const { default: ReviewHistory } = await import(
            '@coze-studio/studio/src/pages/audit/ReviewHistory/ReviewHistory'
          );
          return { Component: ReviewHistory };
        },
      },

      // 性能监控
      {
        path: 'monitoring/performance',
        lazy: async () => {
          const { default: PerformancePage } = await import(
            '@coze-studio/studio/src/pages/monitoring/Performance'
          );
          return { Component: PerformancePage };
        },
      },

      // 智能路由管理
      {
        path: 'routing/management',
        lazy: async () => {
          const { default: RoutingManagementPage } = await import(
            '@coze-studio/studio/src/pages/routing/Management'
          );
          return { Component: RoutingManagementPage };
        },
      },

      // A/B测试
      {
        path: 'routing/abtest',
        lazy: async () => {
          const { default: ABTestList } = await import(
            '@coze-studio/studio/src/pages/routing/ABTest/ABTestList'
          );
          return { Component: ABTestList };
        },
      },

      // 开发者平台 - 用户管理
      {
        path: 'developer/users',
        lazy: async () => {
          const { default: Users } = await import(
            '@coze-studio/studio/src/pages/developer/Users/Users'
          );
          return { Component: Users };
        },
      },

      // 开发者平台 - CLI下载
      {
        path: 'developer/cli',
        lazy: async () => {
          const { default: CLIDownload } = await import(
            '@coze-studio/studio/src/pages/developer/CLI/CLIDownload'
          );
          return { Component: CLIDownload };
        },
      },

      // 开发者平台 - SDK示例
      {
        path: 'developer/sdk',
        lazy: async () => {
          const { default: SDKExamples } = await import(
            '@coze-studio/studio/src/pages/developer/SDK/SDKExamples'
          );
          return { Component: SDKExamples };
        },
      },

      // 开发者平台
      {
        path: 'developer',
        lazy: async () => {
          const { default: DeveloperPlatformPage } = await import(
            '@coze-studio/developer-platform'
          );
          return { Component: DeveloperPlatformPage };
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
