// frontend/packages/api-client/src/types/tenant.ts

import { PageRequest, PageResponse, EntityID, TenantID, Timestamp } from './common';

/**
 * 租户状态
 */
export enum TenantStatus {
  ACTIVE = 'active',
  SUSPENDED = 'suspended',
  DELETED = 'deleted',
}

/**
 * 租户订阅计划
 */
export enum SubscriptionTier {
  FREE = 'free',
  PRO = 'pro',
  ENTERPRISE = 'enterprise',
}

/**
 * 租户DTO
 */
export interface TenantDTO {
  tenant_id: TenantID;
  tenant_name: string;
  tenant_code: string;
  status: TenantStatus;
  subscription_tier: SubscriptionTier;
  contact_email: string;
  contact_phone?: string;
  max_bots: number;
  max_users: number;
  metadata?: Record<string, any>;
  created_at: Timestamp;
  updated_at: Timestamp;
  deleted_at?: Timestamp;
}

/**
 * 租户过滤条件
 */
export interface TenantFilter extends PageRequest {
  tenant_id?: TenantID;
  tenant_name?: string;
  tenant_code?: string;
  status?: TenantStatus;
  subscription_tier?: SubscriptionTier;
  created_at_start?: Timestamp;
  created_at_end?: Timestamp;
}

/**
 * 创建租户请求
 */
export interface CreateTenantRequest {
  tenant_name: string;
  tenant_code: string;
  contact_email: string;
  contact_phone?: string;
  subscription_tier?: SubscriptionTier;
  max_bots?: number;
  max_users?: number;
  metadata?: Record<string, any>;
}

/**
 * 更新租户请求
 */
export interface UpdateTenantRequest {
  tenant_name?: string;
  contact_email?: string;
  contact_phone?: string;
  status?: TenantStatus;
  subscription_tier?: SubscriptionTier;
  max_bots?: number;
  max_users?: number;
  metadata?: Record<string, any>;
}

/**
 * 租户统计DTO
 */
export interface TenantStatsDTO {
  tenant_id: TenantID;
  total_bots: number;
  active_bots: number;
  total_users: number;
  active_users: number;
  total_messages: number;
  total_tokens_used: number;
  storage_used_mb: number;
}
