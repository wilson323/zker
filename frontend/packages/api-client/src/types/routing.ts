// frontend/packages/api-client/src/types/routing.ts

import { PageRequest, EntityID, TenantID, Timestamp } from './common';

/**
 * 路由规则DTO
 */
export interface RoutingRuleDTO {
  rule_id: EntityID;
  tenant_id: TenantID;
  rule_name: string;
  priority: number;
  intent_matcher_id: EntityID;
  target_service: string;
  target_endpoint: string;
  condition?: string; // 条件表达式
  timeout_ms?: number;
  retry_count?: number;
  is_active: boolean;
  metadata?: Record<string, any>;
  created_at: Timestamp;
  updated_at: Timestamp;
  created_by: string;
}

/**
 * 路由规则过滤条件
 */
export interface RoutingRuleFilter extends PageRequest {
  tenant_id?: TenantID;
  rule_name?: string;
  is_active?: boolean;
  priority_min?: number;
  priority_max?: number;
  target_service?: string;
}

/**
 * 创建路由规则请求
 */
export interface CreateRoutingRuleRequest {
  tenant_id: TenantID;
  rule_name: string;
  priority: number;
  intent_matcher_id: EntityID;
  target_service: string;
  target_endpoint: string;
  condition?: string;
  timeout_ms?: number;
  retry_count?: number;
  is_active?: boolean;
  metadata?: Record<string, any>;
}

/**
 * 更新路由规则请求
 */
export interface UpdateRoutingRuleRequest {
  rule_name?: string;
  priority?: number;
  intent_matcher_id?: EntityID;
  target_service?: string;
  target_endpoint?: string;
  condition?: string;
  timeout_ms?: number;
  retry_count?: number;
  is_active?: boolean;
  metadata?: Record<string, any>;
}

/**
 * 意图匹配器类型
 */
export enum IntentMatcherType {
  KEYWORD = 'keyword', // 关键词匹配
  REGEX = 'regex', // 正则表达式匹配
  ML = 'ml', // 机器学习模型匹配
  HYBRID = 'hybrid', // 混合匹配
}

/**
 * 意图匹配器DTO
 */
export interface IntentMatcherDTO {
  matcher_id: EntityID;
  tenant_id: TenantID;
  matcher_name: string;
  matcher_type: IntentMatcherType;
  config: {
    // 关键词配置
    keywords?: string[];
    keyword_match_type?: 'exact' | 'contains' | 'fuzzy';

    // 正则配置
    pattern?: string;
    flags?: string;

    // ML模型配置
    model_id?: string;
    model_version?: string;
    threshold?: number; // 置信度阈值 0-1

    // 混合匹配配置
    matchers?: {
      matcher_id: EntityID;
      weight: number;
    }[];
  };
  score_weight: number; // 评分权重
  is_active: boolean;
  description?: string;
  created_at: Timestamp;
  updated_at: Timestamp;
  created_by: string;
}

/**
 * 意图匹配器过滤条件
 */
export interface IntentMatcherFilter extends PageRequest {
  tenant_id?: TenantID;
  matcher_name?: string;
  matcher_type?: IntentMatcherType;
  is_active?: boolean;
}

/**
 * 创建意图匹配器请求
 */
export interface CreateIntentMatcherRequest {
  tenant_id: TenantID;
  matcher_name: string;
  matcher_type: IntentMatcherType;
  config: Record<string, any>;
  score_weight: number;
  is_active?: boolean;
  description?: string;
}

/**
 * 更新意图匹配器请求
 */
export interface UpdateIntentMatcherRequest {
  matcher_name?: string;
  config?: Record<string, any>;
  score_weight?: number;
  is_active?: boolean;
  description?: string;
}

/**
 * 路由测试请求
 */
export interface RoutingTestRequest {
  tenant_id: TenantID;
  input_text: string;
  context?: Record<string, any>;
}

/**
 * 路由测试响应
 */
export interface RoutingTestResponse {
  matched: boolean;
  confidence: number;
  matched_rules: {
    rule_id: EntityID;
    rule_name: string;
    score: number;
    target_service: string;
    target_endpoint: string;
  }[];
  match_details: {
    matcher_id: EntityID;
    matcher_name: string;
    matcher_type: IntentMatcherType;
    score: number;
  }[];
}

/**
 * 路由统计DTO
 */
export interface RoutingStatsDTO {
  tenant_id: TenantID;
  total_requests: number;
  successful_routes: number;
  failed_routes: number;
  avg_confidence: number;
  most_used_rules: {
    rule_id: EntityID;
    rule_name: string;
    usage_count: number;
  }[];
  period_start: Timestamp;
  period_end: Timestamp;
}
