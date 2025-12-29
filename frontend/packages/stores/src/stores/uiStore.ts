// frontend/packages/stores/src/stores/uiStore.ts

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

/**
 * 主题模式
 */
export type ThemeMode = 'light' | 'dark' | 'auto';

/**
 * 侧边栏状态
 */
export interface SidebarState {
  collapsed: boolean;
  width: number;
}

/**
 * UI 状态
 */
interface UIState {
  // 主题
  themeMode: ThemeMode;

  // 侧边栏
  sidebar: SidebarState;

  // 全局加载状态
  globalLoading: boolean;

  // Toast 消息队列
  toasts: {
    id: string;
    message: string;
    type: 'success' | 'info' | 'warning' | 'error';
    duration?: number;
  }[];

  // Modal 状态
  modalStack: string[];

  // Actions
  setThemeMode: (mode: ThemeMode) => void;
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  setGlobalLoading: (loading: boolean) => void;

  // Toast 管理
  addToast: (
    message: string,
    type: 'success' | 'info' | 'warning' | 'error',
    duration?: number
  ) => void;
  removeToast: (id: string) => void;
  clearToasts: () => void;

  // Modal 管理
  pushModal: (modalId: string) => void;
  popModal: () => void;
  clearModals: () => void;
}

/**
 * UI Store
 */
export const useUIStore = create<UIState>()(
  persist(
    (set, get) => ({
      // 初始状态
      themeMode: 'light',
      sidebar: {
        collapsed: false,
        width: 240,
      },
      globalLoading: false,
      toasts: [],
      modalStack: [],

      // Actions
      setThemeMode: (mode) => set({ themeMode: mode }),

      toggleSidebar: () =>
        set((state) => ({
          sidebar: {
            ...state.sidebar,
            collapsed: !state.sidebar.collapsed,
          },
        })),

      setSidebarCollapsed: (collapsed) =>
        set((state) => ({
          sidebar: {
            ...state.sidebar,
            collapsed,
          },
        })),

      setGlobalLoading: (loading) => set({ globalLoading: loading }),

      // Toast 管理
      addToast: (message, type, duration = 3000) => {
        const id = Date.now().toString();
        set((state) => ({
          toasts: [...state.toasts, { id, message, type, duration }],
        }));

        // 自动移除
        if (duration > 0) {
          setTimeout(() => {
            get().removeToast(id);
          }, duration);
        }
      },

      removeToast: (id) =>
        set((state) => ({
          toasts: state.toasts.filter((t) => t.id !== id),
        })),

      clearToasts: () => set({ toasts: [] }),

      // Modal 管理
      pushModal: (modalId) =>
        set((state) => ({
          modalStack: [...state.modalStack, modalId],
        })),

      popModal: () =>
        set((state) => ({
          modalStack: state.modalStack.slice(0, -1),
        })),

      clearModals: () => set({ modalStack: [] }),
    }),
    {
      name: 'ui-storage',
      partialize: (state) => ({
        themeMode: state.themeMode,
        sidebar: state.sidebar,
      }),
    }
  )
);
