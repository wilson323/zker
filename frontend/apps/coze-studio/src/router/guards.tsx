// frontend/apps/coze-studio/src/router/guards.tsx

import React, { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';

/**
 * 路由守卫组件
 */

// 受保护的路由（需要登录）
export const ProtectedRoute: React.FC<{ children: ReactNode }> = ({ children }) => {
  const token = localStorage.getItem('auth_token');
  const location = useLocation();

  if (!token) {
    // 未登录，重定向到登录页
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
};

// 权限路由（需要特定权限）
export interface PermissionRouteProps {
  children: ReactNode;
  resource: string;
  action: string;
  fallback?: ReactNode;
}

export const PermissionRoute: React.FC<PermissionRouteProps> = ({
  children,
  resource,
  action,
  fallback,
}) => {
  // TODO: 从权限 store 或 API 检查权限
  // 这里先简化处理，实际需要从 usePermissionStore 获取

  const hasPermission = true; // 临时实现

  if (!hasPermission) {
    return <>{fallback || <div>您没有权限访问此页面</div>}</>;
  }

  return <>{children}</>;
};

// 租户路由（需要在特定租户上下文中）
export interface TenantRouteProps {
  children: ReactNode;
  requiredTenantId?: string;
  fallback?: ReactNode;
}

export const TenantRoute: React.FC<TenantRouteProps> = ({
  children,
  requiredTenantId,
  fallback,
}) => {
  const currentTenantId = localStorage.getItem('current_tenant_id');

  if (!currentTenantId) {
    return <>{fallback || <div>请先选择租户</div>}</>;
  }

  if (requiredTenantId && currentTenantId !== requiredTenantId) {
    return <>{fallback || <div>租户不匹配</div>}</>;
  }

  return <>{children}</>;
};
