// frontend/apps/coze-studio/src/router/guards.tsx

import React, { useEffect } from 'react';
import { Navigate, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@coze-studio/stores';
import type { RouteMeta } from './types';

/**
 * 需要登录的路由守卫
 */
export const ProtectedRoute: React.FC<{ meta?: RouteMeta }> = ({ meta }) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const location = useLocation();

  if (!meta?.requireAuth) {
    return <Outlet />;
  }

  if (!isAuthenticated) {
    return (
      <Navigate
        to="/login"
        state={{ from: location.pathname + location.search }}
        replace
      />
    );
  }

  return <Outlet />;
};
