// frontend/packages/api-client/src/types/quota.ts

import { TenantID, Timestamp } from './common';

/**
 * 配额资源类型
 */
export enum QuotaResourceType {
  BOT = 'bot', // 机器人数量
  USER = 'user', // 用户数量
  API_CALL = 'api_call', // API调用次数
  MESSAGE = 'message', // 消息数量
  TOKEN = 'token', // Token使用量
  STORAGE = 'storage', // 存储空间(MB)
  KNOWLEDGE_BASE = 'knowledge_base', // 知识库数量
  WORKFLOW = 'workflow', // 工作流数量
}

/**
 * 配额单位
 */
export enum QuotaUnit {
  COUNT = 'count', // 计数
  MB = 'mb', // 兆字节
  GB = 'gb', // 吉字节
  TIMES = 'times', // 次数
}

/**
 * 配额DTO
 */
export interface QuotaDTO {
  quota_id: string;
  tenant_id: TenantID;
  resource_type: QuotaResourceType;
  limit: number;
  unit: QuotaUnit;
  period?: string; // 限制周期：daily, monthly, yearly
  is_soft_limit: boolean; // 软限制：可超过但需付费
  overage_fee?: number; // 超量单价
  created_at: Timestamp;
  updated_at: Timestamp;
}

/**
 * 配额使用量DTO
 */
export interface QuotaUsageDTO {
  usage_id: string;
  tenant_id: TenantID;
  resource_type: QuotaResourceType;
  usage: number;
  limit: number;
  usage_percentage: number;
  period_start: Timestamp;
  period_end: Timestamp;
  reset_at: Timestamp;
}

/**
 * 配额使用详情DTO
 */
export interface QuotaUsageDetailDTO {
  resource_type: QuotaResourceType;
  current_usage: number;
  limit: number;
  usage_percentage: number;
  remaining: number;
  is_over_limit: boolean;
  period_start: Timestamp;
  period_end: Timestamp;
  reset_at: Timestamp;
  overage_fee?: number;
}

/**
 * 配额历史DTO
 */
export interface QuotaHistoryDTO {
  history_id: string;
  tenant_id: TenantID;
  resource_type: QuotaResourceType;
  action: 'increment' | 'decrement' | 'reset';
  amount: number;
  description?: string;
  created_at: Timestamp;
}

/**
 * 配额预警规则DTO
 */
export interface QuotaAlertRuleDTO {
  rule_id: string;
  tenant_id: TenantID;
  resource_type: QuotaResourceType;
  threshold_percentage: number; // 阈值百分比
  is_enabled: boolean;
  alert_channels: string[]; // 邮件、短信、webhook等
  created_at: Timestamp;
}

/**
 * 更新配额请求
 */
export interface UpdateQuotaRequest {
  limit: number;
  is_soft_limit?: boolean;
  overage_fee?: number;
}

/**
 * 订阅方案DTO
 */
export interface SubscriptionPlanDTO {
  plan_id: string;
  plan_name: string;
  tier: 'free' | 'pro' | 'enterprise';
  price_monthly: number;
  price_yearly: number;
  currency: string;
  quotas: {
    resource_type: QuotaResourceType;
    limit: number;
    unit: QuotaUnit;
  }[];
  features: string[];
  is_active: boolean;
}

/**
 * 租户订阅DTO
 */
export interface TenantSubscriptionDTO {
  subscription_id: string;
  tenant_id: TenantID;
  plan_id: string;
  plan_name: string;
  tier: string;
  status: 'active' | 'suspended' | 'cancelled';
  billing_cycle: 'monthly' | 'yearly';
  price: number;
  currency: string;
  current_period_start: Timestamp;
  current_period_end: Timestamp;
  cancel_at_period_end: boolean;
  created_at: Timestamp;
  updated_at: Timestamp;
}
