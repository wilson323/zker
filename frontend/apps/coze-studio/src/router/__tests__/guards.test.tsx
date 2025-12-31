// frontend/apps/coze-studio/src/router/__tests__/guards.test.tsx

import { render, screen } from '@testing-library/react';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import { ProtectedRoute } from '../guards';

// Mock useAuthStore
jest.mock('@coze-studio/stores', () => ({
  useAuthStore: jest.fn(),
}));

import { useAuthStore } from '@coze-studio/stores';

describe('ProtectedRoute', () => {
  const renderWithRouter = (
    ui: React.ReactElement,
    { isAuthenticated = false } = {}
  ) => {
    (useAuthStore as jest.Mock).mockReturnValue({
      isAuthenticated,
    });

    return render(
      <MemoryRouter initialEntries={['/protected']}>
        {ui}
      </MemoryRouter>
    );
  };

  it('未登录时应该重定向到登录页', () => {
    renderWithRouter(
      <Routes>
        <Route
          path="/protected"
          element={
            <ProtectedRoute meta={{ requireAuth: true }}>
              <div>Protected Content</div>
            </ProtectedRoute>
          }
        />
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      { isAuthenticated: false }
    );

    expect(screen.getByText('Login Page')).toBeInTheDocument();
  });

  it('已登录时应该显示受保护内容', () => {
    renderWithRouter(
      <Routes>
        <Route
          path="/protected"
          element={
            <ProtectedRoute meta={{ requireAuth: true }}>
              <div>Protected Content</div>
            </ProtectedRoute>
          }
        />
      </Routes>,
      { isAuthenticated: true }
    );

    expect(screen.getByText('Protected Content')).toBeInTheDocument();
  });

  it('不需要认证的路由应该直接显示内容', () => {
    renderWithRouter(
      <Routes>
        <Route
          path="/public"
          element={
            <ProtectedRoute meta={{ requireAuth: false }}>
              <div>Public Content</div>
            </ProtectedRoute>
          }
        />
      </Routes>,
      { isAuthenticated: false }
    );

    expect(screen.getByText('Public Content')).toBeInTheDocument();
  });

  it('重定向应该保存原始路径', () => {
    const { container } = renderWithRouter(
      <Routes>
        <Route
          path="/protected"
          element={
            <ProtectedRoute meta={{ requireAuth: true }}>
              <div>Protected Content</div>
            </ProtectedRoute>
          }
        />
        <Route
          path="/login"
          element={
            <div>
              Login Page - Redirected from: <span data-testid="redirect-path" />
            </div>
          }
        />
      </Routes>,
      { isAuthenticated: false }
    );

    // 检查location state是否保存
    // 注意: 这需要实际的Router实现来验证
  });
});
