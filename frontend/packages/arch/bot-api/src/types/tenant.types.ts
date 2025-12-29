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
 * 租户状态枚举
 */
export enum TenantStatus {
  /** 草稿 - 未激活 */
  DRAFT = 'draft',
  /** 活跃 - 正常使用 */
  ACTIVE = 'active',
  /** 暂停 - 欠费或违规 */
  SUSPENDED = 'suspended',
  /** 已删除 - 软删除 */
  DELETED = 'deleted',
}

/**
 * 订阅套餐级别
 */
export enum SubscriptionTier {
  /** 免费版 */
  FREE = 'free',
  /** 基础版 */
  BASIC = 'basic',
  /** 专业版 */
  PROFESSIONAL = 'professional',
  /** 企业版 */
  ENTERPRISE = 'enterprise',
}

/**
 * 计费周期
 */
export enum BillingCycle {
  /** 按月付费 */
  MONTHLY = 'monthly',
  /** 按年付费 */
  YEARLY = 'yearly',
}

/**
 * 订阅状态
 */
export enum SubscriptionStatus {
  /** 活跃 */
  ACTIVE = 'active',
  /** 过期 */
  EXPIRED = 'expired',
  /** 取消 */
  CANCELLED = 'cancelled',
  /** 暂停 */
  SUSPENDED = 'suspended',
}

/**
 * 租户信息
 */
export interface Tenant {
  /** 租户ID */
  tenant_id: string;
  /** 租户名称 */
  tenant_name: string;
  /** 租户描述 */
  description?: string;
  /** 联系人邮箱 */
  contact_email: string;
  /** 联系人电话 */
  contact_phone?: string;
  /** 公司/组织名称 */
  company_name?: string;
  /** 租户状态 */
  status: TenantStatus;
  /** 创建时间 */
  created_at: string;
  /** 更新时间 */
  updated_at: string;
  /** 删除时间 */
  deleted_at?: string;
}

/**
 * 订阅信息
 */
export interface Subscription {
  /** 订阅ID */
  subscription_id: string;
  /** 租户ID */
  tenant_id: string;
  /** 套餐级别 */
  plan_tier: SubscriptionTier;
  /** 计费周期 */
  billing_cycle: BillingCycle;
  /** 订阅状态 */
  status: SubscriptionStatus;
  /** 开始时间 */
  started_at: string;
  /** 过期时间 */
  expires_at: string;
  /** 是否自动续费 */
  auto_renew: boolean;
  /** 月度价格（分） */
  monthly_price_cents?: number;
  /** 年度价格（分） */
  yearly_price_cents?: number;
  /** 折扣百分比 (0-100) */
  discount_percentage?: number;
}

/**
 * 套餐配额
 */
export interface PlanQuota {
  /** 资源类型 */
  resource_type: string;
  /** 配额限制 */
  limit: number;
  /** 单位 */
  unit: string;
}

/**
 * ==================== 租户管理 API ====================
 */

/**
 * 创建租户请求
 */
export interface CreateTenantRequest {
  /** 租户名称（必填，2-100字符） */
  tenant_name: string;
  /** 租户描述（可选，最大500字符） */
  description?: string;
  /** 联系人邮箱（必填） */
  contact_email: string;
  /** 联系人电话（可选） */
  contact_phone?: string;
  /** 公司/组织名称（可选） */
  company_name?: string;
}

/**
 * 创建租户响应
 */
export interface CreateTenantResponse {
  /** 租户信息 */
  tenant: Tenant;
  /** 初始订阅信息 */
  subscription: Subscription;
}

/**
 * 获取租户详情请求
 */
export interface GetTenantRequest {
  /** 租户ID */
  tenant_id: string;
}

/**
 * 获取租户详情响应
 */
export interface GetTenantResponse {
  /** 租户信息 */
  tenant: Tenant;
  /** 订阅信息 */
  subscription: Subscription;
}

/**
 * 更新租户请求
 */
export interface UpdateTenantRequest {
  /** 租户ID */
  tenant_id: string;
  /** 租户名称（可选） */
  tenant_name?: string;
  /** 租户描述（可选） */
  description?: string;
  /** 联系人邮箱（可选） */
  contact_email?: string;
  /** 联系人电话（可选） */
  contact_phone?: string;
  /** 公司/组织名称（可选） */
  company_name?: string;
}

/**
 * 更新租户响应
 */
export interface UpdateTenantResponse {
  /** 更新后的租户信息 */
  tenant: Tenant;
}

/**
 * 列出租户请求（分页）
 */
export interface ListTenantsRequest {
  /** 页码（从1开始） */
  page?: number;
  /** 每页数量（1-100） */
  page_size?: number;
  /** 租户状态过滤 */
  status?: TenantStatus;
  /** 搜索关键词（搜索租户名称或公司名称） */
  search?: string;
}

/**
 * 列出租户响应
 */
export interface ListTenantsResponse {
  /** 租户列表 */
  tenants: Tenant[];
  /** 总数 */
  total: number;
  /** 当前页 */
  page: number;
  /** 每页数量 */
  page_size: number;
}

/**
 * 删除租户请求
 */
export interface DeleteTenantRequest {
  /** 租户ID */
  tenant_id: string;
}

/**
 * 删除租户响应
 */
export interface DeleteTenantResponse {
  /** 是否成功 */
  success: boolean;
  /** 消息 */
  message: string;
}

/**
 * ==================== 订阅管理 API ====================
 */

/**
 * 获取订阅请求
 */
export interface GetSubscriptionRequest {
  /** 租户ID */
  tenant_id: string;
}

/**
 * 获取订阅响应
 */
export interface GetSubscriptionResponse {
  /** 订阅信息 */
  subscription: Subscription;
  /** 套餐配额 */
  quotas: PlanQuota[];
}

/**
 * 更新订阅请求
 */
export interface UpdateSubscriptionRequest {
  /** 租户ID */
  tenant_id: string;
  /** 是否自动续费 */
  auto_renew?: boolean;
  /** 订阅状态（仅允许 Active → Cancelled） */
  status?: SubscriptionStatus;
}

/**
 * 更新订阅响应
 */
export interface UpdateSubscriptionResponse {
  /** 更新后的订阅信息 */
  subscription: Subscription;
}

/**
 * 升级订阅请求
 */
export interface UpgradeSubscriptionRequest {
  /** 租户ID */
  tenant_id: string;
  /** 目标套餐级别 */
  target_plan_tier: SubscriptionTier;
  /** 目标计费周期 */
  target_billing_cycle?: BillingCycle;
  /** 折扣码（可选） */
  discount_code?: string;
}

/**
 * 升级订阅响应
 */
export interface UpgradeSubscriptionResponse {
  /** 升级后的订阅信息 */
  subscription: Subscription;
  /** 新套餐配额 */
  quotas: PlanQuota[];
  /** 订单信息（如有） */
  order?: {
    /** 订单ID */
    order_id: string;
    /** 订单金额（分） */
    amount_cents: number;
    /** 货币 */
    currency: string;
  };
}
