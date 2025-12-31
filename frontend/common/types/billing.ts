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
 * Token计量和预算管理相关类型定义
 *
 * 包含Token使用统计、预算设置、告警等所有类型定义
 */

// ================================================================================
// Token计量相关类型
// ================================================================================

/**
 * Token使用记录请求
 */
export interface RecordTokenUsageRequest {
  /** 用户ID */
  user_id: number;
  /** Bot ID */
  bot_id: string;
  /** 模型提供商 */
  model_provider: string;
  /** 模型名称 */
  model_name: string;
  /** 输入Token数 */
  input_tokens: number;
  /** 输出Token数 */
  output_tokens: number;
  /** 总Token数 */
  total_tokens: number;
  /** 请求类型（chat/completion等） */
  request_type: string;
}

/**
 * Token使用记录响应
 */
export interface RecordTokenUsageResponse {
  /** 记录ID */
  record_id: string;
  /** 租户ID */
  tenant_id: string;
  /** 用户ID */
  user_id: number;
  /** Bot ID */
  bot_id: string;
  /** 模型提供商 */
  model_provider: string;
  /** 模型名称 */
  model_name: string;
  /** 输入Token数 */
  input_tokens: number;
  /** 输出Token数 */
  output_tokens: number;
  /** 总Token数 */
  total_tokens: number;
  /** 请求类型 */
  request_type: string;
  /** 成本（USD） */
  cost_usd: number;
  /** 创建时间（毫秒时间戳） */
  created_at: number;
}

/**
 * 批量记录响应
 */
export interface BatchRecordResponse {
  /** 成功记录数 */
  success_count: number;
  /** 失败记录数 */
  failure_count: number;
  /** 失败记录索引列表 */
  failed_indices: number[];
  /** 错误消息 */
  error?: string;
}

/**
 * 使用统计过滤条件
 */
export interface UsageStatsFilter {
  /** 开始时间（毫秒时间戳） */
  start_time?: number;
  /** 结束时间（毫秒时间戳） */
  end_time?: number;
  /** Bot ID过滤 */
  bot_id?: string;
  /** 模型提供商过滤 */
  model_provider?: string;
  /** 模型名称过滤 */
  model_name?: string;
  /** 请求类型过滤 */
  request_type?: string;
}

/**
 * 使用统计响应
 */
export interface UsageStatsResponse {
  /** 总Token数 */
  total_tokens: number;
  /** 输入Token总数 */
  total_input_tokens: number;
  /** 输出Token总数 */
  total_output_tokens: number;
  /** 总成本（USD） */
  total_cost: number;
  /** 总请求数 */
  total_requests: number;
  /** 缓存命中请求数 */
  cached_requests: number;
  /** 平均响应时间（毫秒） */
  avg_response_time: number;
  /** 模型使用统计 */
  model_usage: ModelUsage[];
  /** Bot使用统计 */
  bot_usage: BotUsage[];
}

/**
 * 模型使用统计
 */
export interface ModelUsage {
  /** 模型提供商 */
  model_provider: string;
  /** 模型名称 */
  model_name: string;
  /** 总Token数 */
  total_tokens: number;
  /** 输入Token数 */
  input_tokens: number;
  /** 输出Token数 */
  output_tokens: number;
  /** 总成本（USD） */
  total_cost: number;
  /** 请求数 */
  request_count: number;
  /** 平均成本（每1000 Token） */
  avg_cost_per_1k: number;
}

/**
 * Bot使用统计
 */
export interface BotUsage {
  /** Bot ID */
  bot_id: string;
  /** Bot名称 */
  bot_name: string;
  /** 总Token数 */
  total_tokens: number;
  /** 总成本（USD） */
  total_cost: number;
  /** 请求数 */
  request_count: number;
  /** 平均成本（每请求） */
  avg_cost_per_request: number;
}

/**
 * 每日使用统计响应
 */
export interface DailyUsageStatsResponse {
  /** 每日统计数据 */
  daily_stats: DailyUsage[];
  /** 统计天数 */
  days: number;
  /** 总Token数 */
  total_tokens: number;
  /** 总成本（USD） */
  total_cost: number;
  /** 总请求数 */
  total_requests: number;
}

/**
 * 每日使用统计
 */
export interface DailyUsage {
  /** 日期（YYYY-MM-DD） */
  date: string;
  /** 时间戳（毫秒） */
  timestamp: number;
  /** 总Token数 */
  total_tokens: number;
  /** 输入Token数 */
  input_tokens: number;
  /** 输出Token数 */
  output_tokens: number;
  /** 总成本（USD） */
  total_cost: number;
  /** 请求数 */
  request_count: number;
  /** 缓存命中数 */
  cached_count: number;
  /** 平均响应时间（毫秒） */
  avg_response_time: number;
}

/**
 * 模型使用统计响应
 */
export interface ModelUsageStatsResponse {
  /** 模型统计列表 */
  model_stats: ModelUsageStat[];
  /** 总Token数 */
  total_tokens: number;
  /** 总成本（USD） */
  total_cost: number;
  /** 总模型数 */
  total_models: number;
}

/**
 * 模型使用统计详情
 */
export interface ModelUsageStat {
  /** 模型提供商 */
  model_provider: string;
  /** 模型名称 */
  model_name: string;
  /** 总Token数 */
  total_tokens: number;
  /** 输入Token数 */
  input_tokens: number;
  /** 输出Token数 */
  output_tokens: number;
  /** 总成本（USD） */
  total_cost: number;
  /** 请求数 */
  request_count: number;
  /** 缓存命中数 */
  cached_count: number;
  /** 平均响应时间（毫秒） */
  avg_response_time: number;
  /** 平均成本（每1000 Token） */
  avg_cost_per_1k: number;
  /** 使用占比（%） */
  usage_percentage: number;
}

// ================================================================================
// 预算管理相关类型
// ================================================================================

/**
 * 预算类型
 */
export enum BudgetType {
  /** 月度预算 */
  MONTHLY = 'monthly',
  /** 季度预算 */
  QUARTERLY = 'quarterly',
  /** 年度预算 */
  YEARLY = 'yearly',
}

/**
 * 预算设置DTO
 */
export interface BudgetSettingsDTO {
  /** 租户ID */
  tenant_id: string;
  /** 预算类型 */
  budget_type: BudgetType;
  /** 预算金额（USD） */
  budget_amount: number;
  /** 一级告警阈值（%） */
  alert_threshold_1: number;
  /** 二级告警阈值（%） */
  alert_threshold_2: number;
  /** 是否启用硬性上限 */
  hard_cap_enabled: boolean;
  /** 硬性上限金额 */
  hard_cap_amount?: number;
  /** 是否启用自动降级 */
  auto_downgrade_enabled: boolean;
  /** 降级模型配置 */
  downgrade_config?: DowngradeConfig;
  /** 通知渠道 */
  notification_channels: NotificationChannel[];
  /** 通知接收人列表 */
  notification_recipients: NotificationRecipient[];
  /** 创建时间（毫秒时间戳） */
  created_at: number;
  /** 更新时间（毫秒时间戳） */
  updated_at: number;
}

/**
 * 降级配置
 */
export interface DowngradeConfig {
  /** 原模型提供商 */
  original_provider: string;
  /** 原模型名称 */
  original_model: string;
  /** 降级模型提供商 */
  downgrade_provider: string;
  /** 降级模型名称 */
  downgrade_model: string;
}

/**
 * 通知渠道
 */
export enum NotificationChannel {
  /** 邮件 */
  EMAIL = 'email',
  /** 短信 */
  SMS = 'sms',
  /** Webhook */
  WEBHOOK = 'webhook',
}

/**
 * 通知接收人
 */
export interface NotificationRecipient {
  /** 接收人ID */
  recipient_id: string;
  /** 接收人类型（user/team） */
  recipient_type: 'user' | 'team';
  /** 接收人地址（邮箱/手机号/Webhook URL） */
  recipient_address: string;
  /** 接收渠道 */
  channels: NotificationChannel[];
}

/**
 * 创建预算请求
 */
export interface CreateBudgetRequest {
  /** 预算类型 */
  budget_type: BudgetType;
  /** 预算金额（USD） */
  budget_amount: number;
  /** 一级告警阈值（%） */
  alert_threshold_1: number;
  /** 二级告警阈值（%） */
  alert_threshold_2: number;
  /** 是否启用硬性上限 */
  hard_cap_enabled: boolean;
  /** 硬性上限金额 */
  hard_cap_amount?: number;
  /** 是否启用自动降级 */
  auto_downgrade_enabled: boolean;
  /** 降级配置 */
  downgrade_config?: DowngradeConfig;
  /** 通知渠道 */
  notification_channels: NotificationChannel[];
  /** 通知接收人列表 */
  notification_recipients: NotificationRecipient[];
}

/**
 * 更新预算请求
 */
export interface UpdateBudgetRequest {
  /** 预算类型 */
  budget_type?: BudgetType;
  /** 预算金额（USD） */
  budget_amount?: number;
  /** 一级告警阈值（%） */
  alert_threshold_1?: number;
  /** 二级告警阈值（%） */
  alert_threshold_2?: number;
  /** 是否启用硬性上限 */
  hard_cap_enabled?: boolean;
  /** 硬性上限金额 */
  hard_cap_amount?: number;
  /** 是否启用自动降级 */
  auto_downgrade_enabled?: boolean;
  /** 降级配置 */
  downgrade_config?: DowngradeConfig;
  /** 通知渠道 */
  notification_channels?: NotificationChannel[];
  /** 通知接收人列表 */
  notification_recipients?: NotificationRecipient[];
}

/**
 * 预算使用DTO
 */
export interface BudgetUsageDTO {
  /** 租户ID */
  tenant_id: string;
  /** 预算金额（USD） */
  budget_amount: number;
  /** 已使用金额（USD） */
  used_amount: number;
  /** 剩余金额（USD） */
  remaining_amount: number;
  /** 使用率（%） */
  usage_percent: number;
  /** 当前周期开始时间（毫秒时间戳） */
  period_start: number;
  /** 当前周期结束时间（毫秒时间戳） */
  period_end: number;
  /** 总Token数 */
  total_tokens: number;
  /** 总请求数 */
  total_requests: number;
  /** 告警触发次数 */
  alert_count: number;
  /** 是否将超出预算 */
  will_exceed_budget: boolean;
  /** 预计超出时间（毫秒时间戳） */
  estimated_exceed_time?: number;
  /** 使用预测 */
  usage_prediction?: UsagePrediction;
}

/**
 * 使用预测
 */
export interface UsagePrediction {
  /** 当前日均使用量（USD） */
  current_daily_avg: number;
  /** 预测周期结束使用量（USD） */
  predicted_end_usage: number;
  /** 预测超出量（USD） */
  predicted_overage: number;
  /** 预计超出时间（毫秒时间戳） */
  estimated_exceed_time?: number;
  /** 建议措施 */
  recommendations: string[];
}

/**
 * 预算检查结果
 */
export interface BudgetCheckResult {
  /** 租户ID */
  tenant_id: string;
  /** 是否允许请求 */
  allowed: boolean;
  /** 拒绝原因 */
  reason?: string;
  /** 当前使用率（%） */
  usage_percent: number;
  /** 是否触发告警 */
  alert_triggered: boolean;
  /** 告警级别 */
  alert_level?: 'warning' | 'critical' | 'exceeded';
  /** 剩余预算（USD） */
  remaining_budget: number;
  /** 是否已启用硬性上限 */
  hard_cap_enabled: boolean;
}

/**
 * 告警过滤条件
 */
export interface AlertFilter {
  /** 告警级别 */
  alert_level?: 'warning' | 'critical' | 'exceeded';
  /** 开始时间（毫秒时间戳） */
  start_time?: number;
  /** 结束时间（毫秒时间戳） */
  end_time?: number;
  /** 分页游标 */
  cursor?: string;
  /** 每页数量 */
  limit?: number;
}

/**
 * 预算告警DTO
 */
export interface BudgetAlertDTO {
  /** 告警ID */
  alert_id: string;
  /** 租户ID */
  tenant_id: string;
  /** 告警级别 */
  alert_level: 'warning' | 'critical' | 'exceeded';
  /** 告警类型 */
  alert_type: string;
  /** 当前使用率（%） */
  usage_percent: number;
  /** 告警消息 */
  alert_message: string;
  /** 是否已发送 */
  is_sent: boolean;
  /** 发送时间（毫秒时间戳） */
  sent_at?: number;
  /** 创建时间（毫秒时间戳） */
  created_at: number;
}

// ================================================================================
// 分页相关类型
// ================================================================================

/**
 * Token记录列表响应
 */
export interface TokenRecordListResponse {
  /** 记录列表 */
  records: RecordTokenUsageResponse[];
  /** 总数 */
  total: number;
  /** 当前页游标 */
  cursor?: string;
  /** 是否有下一页 */
  has_more: boolean;
}

/**
 * 告警列表响应
 */
export interface AlertListResponse {
  /** 告警列表 */
  alerts: BudgetAlertDTO[];
  /** 总数 */
  total: number;
  /** 下一页游标 */
  next_cursor?: string;
  /** 是否有下一页 */
  has_more: boolean;
}
