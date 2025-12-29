// frontend/apps/coze-studio/src/router/index.tsx

import React, { Suspense } from 'react';
import { RouterProvider, createBrowserRouter, Navigate } from 'react-router-dom';
import { Spin } from '@coze-studio/ui-components';
import { routes } from './routes';
import { ProtectedRoute } from './guards';

/**
 * 包装路由以添加守卫
 */
const wrappedRoutes = routes.map((route) => {
  if (route.path === '/login' || route.path === '/register') {
    return route; // 公开路由不需要保护
  }

  return {
    ...route,
    children: route.children?.map((child) => ({
      ...child,
      element: <ProtectedRoute>{child.element}</ProtectedRoute>,
    })),
  };
});

/**
 * 创建路由器
 */
const router = createBrowserRouter(wrappedRoutes);

/**
 * 加载中组件
 */
const PageLoading = () => (
  <div
    style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
    }}
  >
    <Spin size="large" />
  </div>
);

/**
 * 路由提供者
 */
export const AppRouter: React.FC = () => {
  return (
    <Suspense fallback={<PageLoading />}>
      <RouterProvider router={router} />
    </Suspense>
  );
};

export default router;
