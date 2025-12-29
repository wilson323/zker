// frontend/packages/stores/src/stores/authStore.ts

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

/**
 * 用户信息
 */
export interface UserInfo {
  user_id: string;
  username: string;
  email: string;
  display_name?: string;
  avatar_url?: string;
  is_active: boolean;
}

/**
 * 认证状态
 */
interface AuthState {
  // 状态
  isAuthenticated: boolean;
  token: string | null;
  userInfo: UserInfo | null;
  loading: boolean;

  // Actions
  login: (token: string, userInfo: UserInfo) => void;
  logout: () => void;
  updateUserInfo: (userInfo: Partial<UserInfo>) => void;
  setToken: (token: string) => void;
  setLoading: (loading: boolean) => void;
}

/**
 * 认证 Store
 */
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      // 初始状态
      isAuthenticated: false,
      token: null,
      userInfo: null,
      loading: false,

      // Actions
      login: (token, userInfo) =>
        set({
          token,
          userInfo,
          isAuthenticated: true,
          loading: false,
        }),

      logout: () =>
        set({
          token: null,
          userInfo: null,
          isAuthenticated: false,
          loading: false,
        }),

      updateUserInfo: (partialUserInfo) =>
        set((state) => ({
          userInfo: state.userInfo ? { ...state.userInfo, ...partialUserInfo } : null,
        })),

      setToken: (token) => set({ token }),

      setLoading: (loading) => set({ loading }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        token: state.token,
        userInfo: state.userInfo,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
