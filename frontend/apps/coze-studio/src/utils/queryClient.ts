// frontend/apps/coze-studio/src/utils/queryClient.ts

import { QueryClient } from '@tanstack/react-query';

/**
 * React Query Client 配置
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // 数据保持新鲜时间（5分钟）
      staleTime: 5 * 60 * 1000,

      // 缓存时间（10分钟）
      gcTime: 10 * 60 * 1000,

      // 重试次数
      retry: 1,

      // 重试延迟
      retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),

      // 窗口重新获得焦点时重新获取
      refetchOnWindowFocus: false,

      // 组件挂载时重新获取
      refetchOnMount: false,

      // 重连时重新获取
      refetchOnReconnect: true,
    },
    mutations: {
      // 变更重试次数
      retry: 1,

      // 变更重试延迟
      retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
    },
  },
});
