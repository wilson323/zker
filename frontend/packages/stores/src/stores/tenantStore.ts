// frontend/packages/stores/src/stores/tenantStore.ts

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { TenantDTO } from '@coze-studio/api-client';

/**
 * 租户状态
 */
interface TenantState {
  // 当前选中的租户
  currentTenant: TenantDTO | null;

  // 用户可访问的租户列表
  accessibleTenants: TenantDTO[];

  // Actions
  setCurrentTenant: (tenant: TenantDTO | null) => void;
  setAccessibleTenants: (tenants: TenantDTO[]) => void;
  addAccessibleTenant: (tenant: TenantDTO) => void;
  removeAccessibleTenant: (tenantId: string) => void;
  clear: () => void;
}

/**
 * 租户 Store
 */
export const useTenantStore = create<TenantState>()(
  persist(
    (set) => ({
      // 初始状态
      currentTenant: null,
      accessibleTenants: [],

      // Actions
      setCurrentTenant: (tenant) => {
        set({ currentTenant: tenant });
        // 同步到 localStorage（供 HTTP 拦截器使用）
        if (tenant) {
          localStorage.setItem('current_tenant_id', tenant.tenant_id);
        } else {
          localStorage.removeItem('current_tenant_id');
        }
      },

      setAccessibleTenants: (tenants) => set({ accessibleTenants: tenants }),

      addAccessibleTenant: (tenant) =>
        set((state) => ({
          accessibleTenants: [...state.accessibleTenants, tenant],
        })),

      removeAccessibleTenant: (tenantId) =>
        set((state) => ({
          accessibleTenants: state.accessibleTenants.filter((t) => t.tenant_id !== tenantId),
          currentTenant:
            state.currentTenant?.tenant_id === tenantId ? null : state.currentTenant,
        })),

      clear: () =>
        set({
          currentTenant: null,
          accessibleTenants: [],
        }),
    }),
    {
      name: 'tenant-storage',
      partialize: (state) => ({
        currentTenant: state.currentTenant,
        accessibleTenants: state.accessibleTenants,
      }),
    }
  )
);
