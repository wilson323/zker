// frontend/packages/api-client/src/index.ts

/**
 * API Client - 统一导出
 */

// HTTP客户端
export { httpClient, createHTTPClient } from './http/client';
export { requestInterceptor } from './http/interceptor';
export { errorHandler, type APIError } from './http/errorHandler';

// 类型定义
export * from './types';

// 工具
export { API_ENDPOINTS } from './utils/apiEndpoints';

// React Query Hooks
export * from './hooks';
