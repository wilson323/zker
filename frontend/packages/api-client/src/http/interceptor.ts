// frontend/packages/api-client/src/http/interceptor.ts

import { InternalAxiosRequestConfig } from 'axios';

/**
 * 请求拦截器配置
 */
export const requestInterceptor = {
  onRequest: (config: InternalAxiosRequestConfig) => {
    // 添加认证Token
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // 添加租户ID
    const tenantId = localStorage.getItem('current_tenant_id');
    if (tenantId) {
      config.headers['X-Tenant-ID'] = tenantId;
    }

    // 添加请求ID
    config.headers['X-Request-ID'] = generateRequestId();

    return config;
  },

  onRequestError: (error: any) => {
    console.error('请求错误:', error);
    return Promise.reject(error);
  },
};

/**
 * 生成请求ID
 */
function generateRequestId(): string {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}
