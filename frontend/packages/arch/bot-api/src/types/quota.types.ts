/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * 资源类型枚举
 */
export enum ResourceType {
  /** Bot 数量 */
  BOTS = 'bots',
  /** 工作流数量 */
  WORKFLOWS = 'workflows',
  /** 对话数量 */
  CONVERSATIONS = 'conversations',
  /** 消息数量 */
  MESSAGES = 'messages',
  /** 知识库数量 */
  KNOWLEDGE_BASES = 'knowledge_bases',
  /** 文档数量 */
  DOCUMENTS = 'documents',
  /** 存储空间 (bytes) */
  STORAGE = 'storage',
  /** API 调用次数 */
  API_CALLS = 'api_calls',
  /** Token 使用量 */
  TOKENS = 'tokens',
  /** 自定义资源 */
  CUSTOM = 'custom',
}

/**
 * 计费周期
 */
export enum QuotaPeriod {
  /** 每小时 */
  HOURLY = 'hourly',
  /** 每天 */
  DAILY = 'daily',
  /** 每周 */
  WEEKLY = 'weekly',
  /** 每月 */
  MONTHLY = 'monthly',
  /** 每年 */
  YEARLY = 'yearly',
  /** 永久 */
  FOREVER = 'forever',
}

/**
 * 配额信息
 */
export interface Quota {
  /** 资源类型 */
  resource_type: ResourceType;
  /** 配额限制 */
  limit: number;
  /** 当前使用量 */
  usage: number;
  /** 使用率百分比 (0-100) */
  usage_percentage: number;
  /** 单位 */
  unit: string;
  /** 重置周期 */
  reset_period: QuotaPeriod;
  /** 下次重置时间 */
  reset_at?: string;
}

/**
 * 配额使用详情
 */
export interface QuotaUsageDetail {
  /** 资源类型 */
  resource_type: ResourceType;
  /** 当前使用量 */
  current_usage: number;
  /** 配额限制 */
  limit: number;
  /** 使用率百分比 */
  usage_percentage: number;
  /** 使用历史（按时间统计） */
  usage_history: UsageHistoryItem[];
}

/**
 * 使用历史项
 */
export interface UsageHistoryItem {
  /** 时间点 */
  timestamp: string;
  /** 使用量 */
  usage: number;
}

/**
 * ==================== 配额管理 API ====================
 */

/**
 * 获取配额请求
 */
export interface GetQuotasRequest {
  /** 租户ID */
  tenant_id: string;
  /** 资源类型过滤（可选，多个） */
  resource_types?: ResourceType[];
}

/**
 * 获取配额响应
 */
export interface GetQuotasResponse {
  /** 配额列表 */
  quotas: Quota[];
  /** 总体使用率 */
  overall_usage_percentage: number;
  /** 警告阈值 */
  warning_threshold: number;
  /** 危险阈值 */
  critical_threshold: number;
}

/**
 * 检查配额请求
 */
export interface CheckQuotaRequest {
  /** 租户ID */
  tenant_id: string;
  /** 资源类型 */
  resource_type: ResourceType;
  /** 所需数量 */
  required_count: number;
  /** 操作描述（可选，用于审计） */
  operation_description?: string;
}

/**
 * 检查配额响应
 */
export interface CheckQuotaResponse {
  /** 是否允许 */
  allowed: boolean;
  /** 当前使用量 */
  current_usage: number;
  /** 配额限制 */
  limit: number;
  /** 剩余配额 */
  remaining: number;
  /** 拒绝原因（如果不允许） */
  reason?: string;
  /** 建议升级套餐（如果配额不足） */
  suggested_upgrade?: {
    /** 目标套餐 */
    target_plan_tier: string;
    /** 新配额限制 */
    new_limit: number;
  };
}

/**
 * 更新配额限制请求
 * 仅超级管理员可调用
 */
export interface UpdateQuotaLimitRequest {
  /** 租户ID */
  tenant_id: string;
  /** 资源类型 */
  resource_type: ResourceType;
  /** 新的配额限制 */
  limit: number;
  /** 更新原因（可选） */
  reason?: string;
}

/**
 * 更新配额限制响应
 */
export interface UpdateQuotaLimitResponse {
  /** 更新后的配额信息 */
  quota: Quota;
  /** 是否成功 */
  success: boolean;
  /** 消息 */
  message: string;
}

/**
 * 获取配额使用详情请求
 */
export interface GetQuotaUsageRequest {
  /** 租户ID */
  tenant_id: string;
  /** 资源类型 */
  resource_type: ResourceType;
  /** 开始时间（ISO 8601） */
  from?: string;
  /** 结束时间（ISO 8601） */
  to?: string;
}

/**
 * 获取配额使用详情响应
 */
export interface GetQuotaUsageResponse {
  /** 配额使用详情 */
  usage: QuotaUsageDetail;
  /** 趋势数据 */
  trend: {
    /** 增长趋势 (up/down/stable) */
    direction: 'up' | 'down' | 'stable';
    /** 变化百分比 */
    change_percentage: number;
  };
  /** 预计耗尽时间（如果趋势持续） */
  estimated_exhaustion_at?: string;
}
