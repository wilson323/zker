// frontend/apps/coze-studio/src/providers/QueryProvider.tsx

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { ReactNode } from 'react';

/**
 * 创建QueryClient实例
 *
 * 配置说明:
 * - staleTime: 数据保持新鲜的时间(5分钟)
 * - gcTime: 垃圾回收时间(10分钟)
 * - retry: 失败重试次数(1次)
 * - refetchOnWindowFocus: 窗口聚焦时不自动重新获取数据
 */
const createQueryClient = (): QueryClient => {
  return new QueryClient({
    defaultOptions: {
      queries: {
        // 数据保持新鲜时间: 5分钟
        staleTime: 5 * 60 * 1000,

        // 未使用缓存保留时间: 10分钟
        gcTime: 10 * 60 * 1000,

        // 失败重试次数
        retry: 1,

        // 窗口聚焦时不自动重新获取
        refetchOnWindowFocus: false,

        // 组件挂载时不自动重新获取(如果有缓存)
        refetchOnMount: false,

        // 失败时不自动重试(除非是特定错误)
        retryOnMount: false,
      },
      mutations: {
        // 变更失败重试次数
        retry: 1,

        // 变更失败时不抛出错误(由调用方处理)
        throwOnError: false,
      },
    },
    // 查询缓存配置
    queryCache: {
      // 配置全局错误处理
      onError: (error) => {
        console.error('[Query Cache Error]', error);
      },
    },
    // 变更缓存配置
    mutationCache: {
      onError: (error) => {
        console.error('[Mutation Cache Error]', error);
      },
    },
  });
};

/**
 * QueryClient实例(单例)
 */
let queryClientInstance: QueryClient | null = null;

export const getQueryClient = (): QueryClient => {
  if (!queryClientInstance) {
    queryClientInstance = createQueryClient();
  }
  return queryClientInstance;
};

/**
 * QueryProvider组件属性
 */
export interface QueryProviderProps {
  children: ReactNode;
  /**
   * 是否显示ReactQueryDevtools
   * @default process.env.NODE_ENV === 'development'
   */
  showDevtools?: boolean;
}

/**
 * React Query Provider组件
 *
 * 使用示例:
 * ```tsx
 * <QueryProvider>
 *   <App />
 * </QueryProvider>
 * ```
 */
export const QueryProvider: React.FC<QueryProviderProps> = ({
  children,
  showDevtools = process.env.NODE_ENV === 'development',
}) => {
  const queryClient = getQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      {showDevtools && (
        <ReactQueryDevtools
          initialIsOpen={false}
          position="bottom-right"
        />
      )}
    </QueryClientProvider>
  );
};
