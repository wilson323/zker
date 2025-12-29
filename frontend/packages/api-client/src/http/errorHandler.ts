// frontend/packages/api-client/src/http/errorHandler.ts

import { AxiosError } from 'axios';

/**
 * API错误响应类型
 */
export interface APIError {
  code: string;
  message: string;
  message_en?: string;
  details?: any;
  request_id?: string;
  timestamp: string;
}

/**
 * 错误处理器
 */
export const errorHandler = {
  onResponseError: (error: AxiosError) => {
    if (error.response) {
      const { status, data } = error.response;
      const apiError: APIError = data as APIError;

      // 处理特定错误码
      switch (status) {
        case 401:
          // 未授权 - 跳转登录
          handleUnauthorized();
          break;
        case 403:
          // 无权限
          handleForbidden(apiError);
          break;
        case 404:
          // 资源不存在
          handleNotFound(apiError);
          break;
        case 500:
          // 服务器错误
          handleServerError(apiError);
          break;
        default:
          handleOtherError(apiError);
      }

      return Promise.reject(apiError);
    }

    if (error.request) {
      // 请求已发出但无响应
      console.error('网络错误:', error.message);
      return Promise.reject({
        code: 'NETWORK_ERROR',
        message: '网络连接失败，请检查网络设置',
        timestamp: new Date().toISOString(),
      });
    }

    // 请求配置错误
    console.error('请求配置错误:', error.message);
    return Promise.reject({
      code: 'REQUEST_CONFIG_ERROR',
      message: '请求配置错误',
      timestamp: new Date().toISOString(),
    });
  },
};

/**
 * 处理401未授权
 */
function handleUnauthorized() {
  // 清除本地存储
  localStorage.removeItem('auth_token');
  localStorage.removeItem('user_info');

  // 跳转登录页
  window.location.href = '/login';
}

/**
 * 处理403无权限
 */
function handleForbidden(error: APIError) {
  console.error('权限不足:', error);
  // 可以显示全局提示
}

/**
 * 处理404资源不存在
 */
function handleNotFound(error: APIError) {
  console.error('资源不存在:', error);
}

/**
 * 处理500服务器错误
 */
function handleServerError(error: APIError) {
  console.error('服务器错误:', error);
  // 可以显示全局错误提示
}

/**
 * 处理其他错误
 */
function handleOtherError(error: APIError) {
  console.error('API错误:', error);
}
