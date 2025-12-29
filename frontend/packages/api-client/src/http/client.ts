// frontend/packages/api-client/src/http/client.ts

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { requestInterceptor } from './interceptor';
import { errorHandler } from './errorHandler';

/**
 * 创建HTTP客户端实例
 */
export const createHTTPClient = (baseURL: string): AxiosInstance => {
  const client = axios.create({
    baseURL,
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // 请求拦截器
  client.interceptors.request.use(
    requestInterceptor.onRequest,
    requestInterceptor.onRequestError
  );

  // 响应拦截器
  client.interceptors.response.use(
    (response) => response,
    errorHandler.onResponseError
  );

  return client;
};

/**
 * 默认客户端实例
 */
export const httpClient = createHTTPClient(
  process.env.API_BASE_URL || 'http://localhost:8888'
);
