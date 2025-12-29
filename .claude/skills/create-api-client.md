# API Client 开发 Skill

## 技能描述

创建符合项目规范的 API Client，用于前端调用后端服务接口。

## 适用场景

- 需要调用新的后端 API
- 需要封装 HTTP 请求逻辑
- 需要统一错误处理
- 需要添加请求拦截器

## API Client 架构

### 1. 目录结构

```
frontend/packages/arch/bot-api/src/
├── axios.ts                    # Axios 实例配置
├── basic-api.ts               # 基础 API 客户端
├── knowledge-api.ts           # 知识库 API 客户端
├── workflow-api.ts            # 工作流 API 客户端
├── coze-space-api.ts          # 空间 API 客户端
└── idl/                       # IDL 生成的类型定义
    ├── basic_api/
    ├── knowledge/
    └── workflow/
```

### 2. 命名规范

```
文件命名：
- 格式：{service}-api.ts (kebab-case)
- 示例：basic-api.ts, knowledge-api.ts, workflow-api.ts

导出命名：
- 格式：{service}Api (camelCase)
- 示例：basicApi, knowledgeApi, workflowApi
```

## API Client 模板

### 1. 基础 API Client

```typescript
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

import BasicApiService from '@coze-arch/idl/basic_api';

import { axiosInstance, type BotAPIRequestConfig } from './axios';

/**
 * 基础服务 API 客户端
 *
 * 封装用户、空间等基础服务的 API 调用。
 */
export const basicApi = new BasicApiService<BotAPIRequestConfig>({
  request: (params, config = {}) =>
    axiosInstance.request({
      ...params,
      ...config,
      headers: {
        ...params.headers,
        ...config.headers,
        'Agw-Js-Conv': 'str',
      },
    }),
});
```

### 2. 完整 API Client 模板

```typescript
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

import { ServiceName } from './idl/{service_name}';
import { axiosInstance, type BotAPIRequestConfig } from './axios';

/**
 * {ServiceDisplayName} API 客户端
 *
 * 提供 {domain} 相关的 API 调用能力。
 *
 * @example
 * ```typescript
 * import { {serviceName}Api } from '@coze-arch/bot-api';
 *
 * // 调用 API
 * const result = await {serviceName}Api.{method}({ param1: 'value1' });
 * ```
 */
// eslint-disable-next-line @typescript-eslint/naming-convention
export const {serviceName}Api = new ServiceName<BotAPIRequestConfig>({
  request: (params, config = {}) => {
    const { headers } = config;
    const reqHeaders = {
      ...headers,
      'Agw-Js-Conv': 'str',
    };
    return axiosInstance.request({
      ...params,
      ...config,
      headers: reqHeaders,
    });
  },
});
```

### 3. 自定义 API Client

```typescript
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

import { axiosInstance, type BotAPIRequestConfig } from './axios';
import type { AxiosRequestConfig } from 'axios';

/**
 * 自定义 API 请求配置
 */
interface CustomRequestConfig extends BotAPIRequestConfig {
  skipErrorHandler?: boolean;  // 跳过全局错误处理
  timeout?: number;             // 请求超时时间
}

/**
 * {ServiceDisplayName} API 客户端
 *
 * 提供 {domain} 相关的 API 调用能力。
 * 包含自定义错误处理和请求拦截。
 */
export const {serviceName}Api = {
  /**
   * 创建资源
   *
   * @param data - 创建数据
   * @param config - 请求配置
   * @returns 创建的资源
   */
  async create<T>(
    data: T,
    config?: CustomRequestConfig,
  ): Promise<{ id: string }> {
    const response = await axiosInstance.request<
      unknown,
      { id: string },
      T
    >({
      method: 'POST',
      url: '/api/v1/{resources}',
      data,
      ...config,
    });
    return response.data;
  },

  /**
   * 获取资源详情
   *
   * @param id - 资源ID
   * @param config - 请求配置
   * @returns 资源详情
   */
  async getById<T>(
    id: string,
    config?: CustomRequestConfig,
  ): Promise<T> {
    const response = await axiosInstance.request<
      unknown,
      T,
      unknown
    >({
      method: 'GET',
      url: `/api/v1/{resources}/${id}`,
      ...config,
    });
    return response.data;
  },

  /**
   * 更新资源
   *
   * @param id - 资源ID
   * @param data - 更新数据
   * @param config - 请求配置
   * @returns 更新后的资源
   */
  async update<T>(
    id: string,
    data: Partial<T>,
    config?: CustomRequestConfig,
  ): Promise<T> {
    const response = await axiosInstance.request<
      unknown,
      T,
      Partial<T>
    >({
      method: 'PUT',
      url: `/api/v1/{resources}/${id}`,
      data,
      ...config,
    });
    return response.data;
  },

  /**
   * 删除资源
   *
   * @param id - 资源ID
   * @param config - 请求配置
   * @returns 删除结果
   */
  async delete(
    id: string,
    config?: CustomRequestConfig,
  ): Promise<void> {
    await axiosInstance.request({
      method: 'DELETE',
      url: `/api/v1/{resources}/${id}`,
      ...config,
    });
  },

  /**
   * 列表查询
   *
   * @param params - 查询参数
   * @param config - 请求配置
   * @returns 资源列表
   */
  async list<T>(
    params?: {
      page?: number;
      pageSize?: number;
      query?: string;
      sortBy?: string;
      sortOrder?: 'asc' | 'desc';
    },
    config?: CustomRequestConfig,
  ): Promise<{ data: T[]; total: number }> {
    const response = await axiosInstance.request<
      unknown,
      { data: T[]; total: number },
      unknown
    >({
      method: 'GET',
      url: '/api/v1/{resources}',
      params,
      ...config,
    });
    return response.data;
  },
};
```

## Axios 配置规范

### 1. Axios 实例配置

```typescript
// axios.ts
import axios from 'axios';

/**
 * Bot API 请求配置接口
 */
export interface BotAPIRequestConfig extends AxiosRequestConfig {
  // 自定义配置
  skipErrorHandler?: boolean;
  retryTimes?: number;
}

/**
 * Axios 实例
 *
 * 统一的 HTTP 客户端实例，配置了基础 URL、超时、拦截器等。
 */
export const axiosInstance = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
axiosInstance.interceptors.request.use(
  (config) => {
    // 添加认证 token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // 添加时间戳防止缓存
    if (config.method === 'get') {
      config.params = {
        ...config.params,
        _t: Date.now(),
      };
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// 响应拦截器
axiosInstance.interceptors.response.use(
  (response) => {
    // 统一响应格式
    return response.data;
  },
  (error) => {
    // 统一错误处理
    if (error.response) {
      const { status, data } = error.response;

      switch (status) {
        case 401:
          // 未授权，跳转登录
          window.location.href = '/login';
          break;
        case 403:
          // 无权限
          console.error('Permission denied');
          break;
        case 500:
          // 服务器错误
          console.error('Server error:', data.message);
          break;
        default:
          console.error('Request failed:', data.message);
      }
    }

    return Promise.reject(error);
  },
);
```

### 2. 错误处理

```typescript
/**
 * API 错误类
 */
export class APIError extends Error {
  constructor(
    public code: number,
    public message: string,
    public details?: unknown,
  ) {
    super(message);
    this.name = 'APIError';
  }
}

/**
 * 处理 API 响应错误
 */
function handleResponseError(error: unknown): never {
  if (axios.isAxiosError(error)) {
    const response = error.response;

    if (response) {
      throw new APIError(
        response.status,
        response.data?.message || 'Request failed',
        response.data,
      );
    }

    if (error.request) {
      throw new APIError(0, 'Network error: No response received');
    }
  }

  throw new APIError(0, error instanceof Error ? error.message : 'Unknown error');
}
```

## 类型定义规范

### 1. 请求类型

```typescript
/**
 * 创建{Resource}请求
 */
export interface Create{Resource}Request {
  name: string;
  description?: string;
  spaceId: number;
  metadata?: Record<string, unknown>;
}

/**
 * 更新{Resource}请求
 */
export interface Update{Resource}Request {
  name?: string;
  description?: string;
  metadata?: Record<string, unknown>;
}

/**
 * 查询{Resource}请求
 */
export interface List{Resource}Request {
  page?: number;
  pageSize?: number;
  query?: string;
  status?: string;
  sortBy?: 'createdAt' | 'updatedAt';
  sortOrder?: 'asc' | 'desc';
}
```

### 2. 响应类型

```typescript
/**
 * {Resource}响应
 */
export interface {Resource} {
  id: string;
  name: string;
  description: string;
  spaceId: number;
  status: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * {Resource}列表响应
 */
export interface List{Resource}Response {
  data: {Resource}[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * 分页信息
 */
export interface PaginationInfo {
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  hasNext: boolean;
  hasPrevious: boolean;
}
```

## 使用方式

### 1. 基础使用

```typescript
import { knowledgeApi } from '@coze-arch/bot-api';

// 获取知识库列表
const listKnowledge = async () => {
  try {
    const result = await knowledgeApi.ListKnowledge({
      spaceId: 123,
      page: 1,
      pageSize: 20,
    });
    console.log('Knowledge list:', result);
  } catch (error) {
    console.error('Failed to list knowledge:', error);
  }
};
```

### 2. React Hook 封装

```typescript
import { useState, useEffect } from 'react';
import { knowledgeApi } from '@coze-arch/bot-api';
import type { Knowledge } from './types';

/**
 * 使用知识库列表 Hook
 */
export function useKnowledgeList(spaceId: number) {
  const [data, setData] = useState<Knowledge[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchList = async () => {
    setLoading(true);
    setError(null);

    try {
      const result = await knowledgeApi.ListKnowledge({
        spaceId,
        page: 1,
        pageSize: 20,
      });
      setData(result.data || []);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchList();
  }, [spaceId]);

  return {
    data,
    loading,
    error,
    refetch: fetchList,
  };
}
```

### 3. 错误处理

```typescript
import { knowledgeApi } from '@coze-arch/bot-api';
import { APIError } from './axios';

const createKnowledge = async (data: CreateKnowledgeRequest) => {
  try {
    const result = await knowledgeApi.CreateKnowledge(data);
    return result;
  } catch (error) {
    if (error instanceof APIError) {
      // 处理业务错误
      switch (error.code) {
        case 40001:
          console.error('Invalid parameters');
          break;
        case 40003:
          console.error('Resource already exists');
          break;
        default:
          console.error('Unknown error:', error.message);
      }
    } else {
      // 处理网络错误等
      console.error('Network error:', error);
    }
    throw error;
  }
};
```

## 高级功能

### 1. 请求取消

```typescript
import { useRef, useEffect } from 'react';
import { knowledgeApi } from '@coze-arch/bot-api';

export function useKnowledgeList(spaceId: number) {
  const abortControllerRef = useRef<AbortController | null>(null);

  useEffect(() => {
    // 取消之前的请求
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }

    // 创建新的 AbortController
    abortControllerRef.current = new AbortController();

    const fetchList = async () => {
      try {
        const result = await knowledgeApi.ListKnowledge(
          { spaceId },
          {
            signal: abortControllerRef.current?.signal,
          },
        );
        // 处理结果
      } catch (error) {
        if (error.name !== 'AbortError') {
          console.error('Failed to fetch:', error);
        }
      }
    };

    fetchList();

    return () => {
      abortControllerRef.current?.abort();
    };
  }, [spaceId]);
}
```

### 2. 请求重试

```typescript
/**
 * 带重试的请求
 */
async function requestWithRetry<T>(
  fn: () => Promise<T>,
  maxRetries = 3,
  delay = 1000,
): Promise<T> {
  let lastError: Error | null = null;

  for (let i = 0; i < maxRetries; i++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error as Error;

      // 如果是最后一次重试，直接抛出错误
      if (i === maxRetries - 1) {
        throw error;
      }

      // 等待后重试
      await new Promise(resolve => setTimeout(resolve, delay * (i + 1)));
    }
  }

  throw lastError;
}

// 使用
const result = await requestWithRetry(
  () => knowledgeApi.GetKnowledgeById({ id: '123' }),
  3,
  1000,
);
```

### 3. 请求缓存

```typescript
const cache = new Map<string, { data: unknown; timestamp: number }>();
const CACHE_TTL = 5 * 60 * 1000; // 5分钟

/**
 * 带缓存的请求
 */
async function requestWithCache<T>(
  key: string,
  fn: () => Promise<T>,
): Promise<T> {
  const cached = cache.get(key);

  if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
    return cached.data as T;
  }

  const data = await fn();
  cache.set(key, { data, timestamp: Date.now() });

  return data;
}

// 使用
const result = await requestWithCache(
  `knowledge:${id}`,
  () => knowledgeApi.GetKnowledgeById({ id }),
);
```

## 测试规范

### 1. 单元测试

```typescript
import { knowledgeApi } from '@coze-arch/bot-api';
import { MockAxios } from '../__mocks__/axios';

describe('Knowledge API', () => {
  let mockAxios: MockAxios;

  beforeEach(() => {
    mockAxios = new MockAxios();
  });

  afterEach(() => {
    mockAxios.restore();
  });

  test('should list knowledge successfully', async () => {
    const mockData = {
      data: [
        { id: '1', name: 'Knowledge 1' },
        { id: '2', name: 'Knowledge 2' },
      ],
      total: 2,
    };

    mockAxios.mockResponse(mockData);

    const result = await knowledgeApi.ListKnowledge({
      spaceId: 123,
    });

    expect(result.data).toHaveLength(2);
    expect(result.total).toBe(2);
  });

  test('should handle error correctly', async () => {
    mockAxios.mockError(new Error('Network error'));

    await expect(
      knowledgeApi.ListKnowledge({ spaceId: 123 }),
    ).rejects.toThrow('Network error');
  });
});
```

## 生成流程

### 步骤 1: 定义接口

1. 确定需要调用的后端 API
2. 分析请求和响应结构
3. 定义 TypeScript 类型

### 步骤 2: 创建 Client

1. 从 IDL 生成或手动创建 Service 类
2. 配置请求拦截器
3. 配置响应拦截器
4. 实现错误处理

### 步骤 3: 封装 Hook

1. 创建自定义 Hook 封装 API 调用
2. 处理加载状态
3. 处理错误状态
4. 添加缓存机制（可选）

### 步骤 4: 编写测试

1. Mock Axios 请求
2. 测试成功场景
3. 测试错误场景
4. 测试边界条件

## 输出检查清单

创建 API Client 后，确保：

- [ ] 文件命名符合规范
- [ ] 导出命名符合规范
- [ ] 使用统一的 axiosInstance
- [ ] 配置正确的请求头
- [ ] 实现错误处理
- [ ] 类型定义完整
- [ ] 添加完整的注释
- [ ] 编写使用示例
- [ ] 编写单元测试
- [ ] 测试覆盖主要场景

## 真实代码示例

### 示例 1: Basic API

```typescript
// frontend/packages/arch/bot-api/src/basic-api.ts
import BasicApiService from '@coze-arch/idl/basic_api';

import { axiosInstance, type BotAPIRequestConfig } from './axios';

export const basicApi = new BasicApiService<BotAPIRequestConfig>({
  request: (params, config = {}) =>
    axiosInstance.request({
      ...params,
      ...config,
      headers: {
        ...params.headers,
        ...config.headers,
        'Agw-Js-Conv': 'str',
      },
    }),
});
```

### 示例 2: Knowledge API

```typescript
// frontend/packages/arch/bot-api/src/knowledge-api.ts
import KnowledgeService from './idl/knowledge';
import { axiosInstance, type BotAPIRequestConfig } from './axios';

// eslint-disable-next-line @typescript-eslint/naming-convention
export const KnowledgeApi = new KnowledgeService<BotAPIRequestConfig>({
  request: (params, config = {}) => {
    const { headers } = config;
    const reqHeaders = {
      ...headers,
      'Agw-Js-Conv': 'str',
    };
    return axiosInstance.request({
      ...params,
      ...config,
      headers: reqHeaders,
    });
  },
});
```

## 注意事项

1. **统一实例**: 使用统一的 axiosInstance，便于全局配置和拦截
2. **类型安全**: 充分利用 TypeScript 类型检查
3. **错误处理**: 统一的错误处理机制
4. **请求取消**: 组件卸载时取消未完成的请求
5. **性能优化**: 合理使用缓存、防抖、节流
6. **测试覆盖**: 为关键 API 调用编写测试
