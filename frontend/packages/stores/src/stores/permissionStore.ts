// frontend/packages/stores/src/stores/permissionStore.ts

import { create } from 'zustand';
import { PermissionDTO, PermissionResource, PermissionAction } from '@coze-studio/api-client';

/**
 * 权限缓存项
 */
interface PermissionCacheItem {
  hasPermission: boolean;
  timestamp: number;
}

/**
 * 权限状态
 */
interface PermissionState {
  // 用户权限列表
  permissions: PermissionDTO[];

  // 权限检查缓存
  permissionCache: Record<string, PermissionCacheItem>;

  // 缓存有效期（毫秒）- 默认 5 分钟
  cacheExpiry: number;

  // Actions
  setPermissions: (permissions: PermissionDTO[]) => void;
  clearPermissions: () => void;

  // 检查权限
  checkPermission: (resource: PermissionResource, action: PermissionAction) => boolean;
  checkPermissionSync: (
    resource: PermissionResource,
    action: PermissionAction
  ) => Promise<boolean>;

  // 缓存管理
  setCacheExpiry: (expiry: number) => void;
  clearCache: () => void;
}

/**
 * 生成权限缓存键
 */
function getPermissionCacheKey(resource: string, action: string): string {
  return `${resource}:${action}`;
}

/**
 * 权限 Store
 */
export const usePermissionStore = create<PermissionState>((set, get) => ({
  // 初始状态
  permissions: [],
  permissionCache: {},
  cacheExpiry: 5 * 60 * 1000, // 5 分钟

  // Actions
  setPermissions: (permissions) => set({ permissions }),

  clearPermissions: () => set({ permissions: [], permissionCache: {} }),

  checkPermission: (resource, action) => {
    const { permissions, permissionCache, cacheExpiry } = get();
    const cacheKey = getPermissionCacheKey(resource, action);
    const now = Date.now();

    // 检查缓存
    const cached = permissionCache[cacheKey];
    if (cached && now - cached.timestamp < cacheExpiry) {
      return cached.hasPermission;
    }

    // 检查权限
    const hasPermission = permissions.some(
      (p) => p.resource === resource && p.action === action
    );

    // 更新缓存
    set((state) => ({
      permissionCache: {
        ...state.permissionCache,
        [cacheKey]: {
          hasPermission,
          timestamp: now,
        },
      },
    }));

    return hasPermission;
  },

  checkPermissionSync: async (resource, action) => {
    // TODO: 实际应该调用 API 检查权限
    // 这里简化处理，直接返回缓存的权限检查结果
    return get().checkPermission(resource, action);
  },

  setCacheExpiry: (expiry) => set({ cacheExpiry: expiry }),

  clearCache: () => set({ permissionCache: {} }),
}));
