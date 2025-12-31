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
 * 监控模块TypeScript类型定义
 */

// ============================================================
// 指标类型
// ============================================================

export type MetricType =
  | 'qps'
  | 'response_time'
  | 'error_rate'
  | 'concurrency'
  | 'cpu_usage'
  | 'memory_usage'
  | 'disk_io'
  | 'network_io';

// ============================================================
// QPS指标
// ============================================================

export interface QPSMetrics {
  count: number;
  average: number;
  min: number;
  max: number;
  p50: number;
  p95: number;
  p99: number;
  trend: 'up' | 'down' | 'stable';
  trend_change: number;
}

// ============================================================
// 响应时间指标
// ============================================================

export interface ResponseTimeMetrics {
  average: number;
  min: number;
  max: number;
  p50: number;
  p95: number;
  p99: number;
  trend: 'up' | 'down' | 'stable';
  trend_change: number;
}

// ============================================================
// 错误率指标
// ============================================================

export interface ErrorRateMetrics {
  total_requests: number;
  error_requests: number;
  error_rate: number;
  average: number;
  max: number;
  trend: 'up' | 'down' | 'stable';
  trend_change: number;
  error_breakdown: Record<string, number>;
}

// ============================================================
// 并发指标
// ============================================================

export interface ConcurrencyMetrics {
  current: number;
  average: number;
  peak: number;
  peak_time?: string;
  trend: 'up' | 'down' | 'stable';
  trend_change: number;
}

// ============================================================
// 健康评分
// ============================================================

export interface HealthScore {
  score: number;
  level: 'excellent' | 'good' | 'warning' | 'critical';
  qps_score: number;
  response_score: number;
  error_score: number;
  resource_score: number;
  trend: 'improving' | 'stable' | 'degrading';
  assessed_at: string;
}

// ============================================================
// 租户概览
// ============================================================

export interface TenantOverview {
  tenant_id: string;
  tenant_name: string;
  health_score: HealthScore;
  current_qps: number;
  avg_response_time: number;
  error_rate: number;
  active_agents: number;
  total_requests: number;
  cpu_usage: number;
  memory_usage: number;
  disk_usage: number;
  alert_count: number;
  critical_alert_count: number;
  subscription_tier: string;
  quota_usage_percent: number;
  last_updated: string;
}

// ============================================================
// 告警相关类型
// ============================================================

export type AlertSeverity = 'warning' | 'critical' | 'emergency';
export type AlertStatus = 'pending' | 'acknowledged' | 'resolved' | 'silenced';
export type AlertComparisonOperator = 'gt' | 'lt' | 'eq' | 'gte' | 'lte';
export type NotificationChannel = 'email' | 'sms' | 'webhook' | 'slack' | 'dingtalk' | 'wechat';

export interface AlertRule {
  id: string;
  tenant_id: string;
  rule_name: string;
  description?: string;
  metric_type: string;
  threshold: number;
  comparison: AlertComparisonOperator;
  severity: AlertSeverity;
  notification_channels: NotificationChannel[];
  notification_config?: string;
  evaluation_interval: number;
  for_duration: number;
  is_enabled: boolean;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface AlertHistory {
  id: string;
  tenant_id: string;
  alert_rule_id: string;
  alert_rule_name: string;
  severity: AlertSeverity;
  alert_message: string;
  alert_data: string;
  agent_id?: string;
  status: AlertStatus;
  acknowledged_by?: string;
  acknowledged_at?: string;
  resolved_at?: string;
  silenced_until?: string;
  resolved_by?: string;
  resolution_note?: string;
  created_at: string;
  updated_at: string;
}

export interface AlertFilter {
  tenant_id?: string;
  alert_rule_id?: string;
  severity?: AlertSeverity;
  status?: AlertStatus;
  agent_id?: string;
  start_time?: string;
  end_time?: string;
  page_token?: string;
  page_size?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

export interface AlertStatistics {
  total_count: number;
  pending_count: number;
  acknowledged_count: number;
  resolved_count: number;
  severity_distribution: Record<AlertSeverity, number>;
  trend_data: TrendDataPoint[];
}

export interface TrendDataPoint {
  timestamp: string;
  count: number;
  severity: AlertSeverity;
}

// ============================================================
// 性能报告相关类型
// ============================================================

export type ReportType = 'daily' | 'weekly' | 'monthly';

export interface PerformanceReport {
  id: string;
  tenant_id: string;
  agent_id?: string;
  report_type: ReportType;
  report_date: string;
  qps_metrics?: string;
  response_time_metrics?: string;
  error_rate: number;
  satisfaction?: number;
  token_usage?: number;
  estimated_cost?: number;
  recommendations?: string;
  summary?: string;
  top_issues?: string;
  improvements?: string;
  generated_by: string;
  created_at: string;
}

export interface ReportFilter {
  tenant_id: string;
  agent_id?: string;
  report_type?: ReportType;
  start_date?: string;
  end_date?: string;
  page_token?: string;
  page_size?: number;
}

// ============================================================
// Agent监控相关类型
// ============================================================

export type AgentMetricType =
  | 'qps'
  | 'response_time'
  | 'error_rate'
  | 'satisfaction'
  | 'token_usage'
  | 'cost'
  | 'conversation'
  | 'active_users';

export interface AgentMetrics {
  id: string;
  tenant_id: string;
  agent_id: string;
  agent_name: string;
  metric_type: AgentMetricType;
  metric_value: number;
  metric_timestamp: string;
  metadata?: string;
  created_at: string;
  updated_at: string;
}

export interface AgentQPSMetrics {
  agent_id: string;
  agent_name: string;
  count: number;
  average: number;
  min: number;
  max: number;
  p50: number;
  p95: number;
  p99: number;
  trend: 'up' | 'down' | 'stable';
  trend_change: number;
}

export interface AgentResponseTimeMetrics {
  agent_id: string;
  agent_name: string;
  average: number;
  min: number;
  max: number;
  p50: number;
  p95: number;
  p99: number;
  trend: 'up' | 'down' | 'stable';
  trend_change: number;
}

export interface AgentErrorRateMetrics {
  agent_id: string;
  agent_name: string;
  total_requests: number;
  error_requests: number;
  error_rate: number;
  average: number;
  max: number;
  trend: 'up' | 'down' | 'stable';
  trend_change: number;
  error_breakdown: Record<string, number>;
}

export interface AgentPerformanceReport {
  agent_id: string;
  agent_name: string;
  report_date: string;
  time_range: {
    start_time: string;
    end_time: string;
  };
  qps_metrics?: AgentQPSMetrics;
  response_time_metrics?: AgentResponseTimeMetrics;
  error_rate_metrics?: AgentErrorRateMetrics;
  satisfaction_metrics?: any;
  token_usage_metrics?: any;
  overall_score: number;
  performance_level: string;
  recommendations: string[];
  compared_with_period: string;
  improvements: string[];
  concerns: string[];
  generated_at: string;
}

export interface AgentComparison {
  agent_id: string;
  agent_name: string;
  qps_ranking: number;
  response_time_ranking: number;
  error_rate_ranking: number;
  satisfaction_ranking: number;
  overall_ranking: number;
  qps_percentile: number;
  response_time_percentile: number;
  error_rate_percentile: number;
  satisfaction_percentile: number;
}

export interface AgentHealthStatus {
  agent_id: string;
  agent_name: string;
  status: 'healthy' | 'warning' | 'critical' | 'offline';
  health_score: number;
  qps_status: string;
  response_time_status: string;
  error_rate_status: string;
  last_active_time: string;
  issues: string[];
  alert_count: number;
}

// ============================================================
// API请求和响应类型
// ============================================================

export interface TimeRangeParams {
  start_time: string;
  end_time: string;
  preset?: string;
}

export interface CreateAlertRuleRequest {
  tenant_id: string;
  rule_name: string;
  description?: string;
  metric_type: string;
  threshold: number;
  comparison: AlertComparisonOperator;
  severity: AlertSeverity;
  notification_channels: NotificationChannel[];
  notification_config?: string;
  evaluation_interval?: number;
  for_duration?: number;
}

export interface AcknowledgeAlertRequest {
  user_id: string;
}

export interface ResolveAlertRequest {
  user_id: string;
  note?: string;
}

export interface SilenceAlertRequest {
  duration_minutes: number;
}

export interface RecordMetricRequest {
  tenant_id: string;
  metric_type: MetricType;
  value: number;
  tags?: Record<string, string>;
}

// ============================================================
// API响应类型
// ============================================================

export interface APIResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page_token?: string;
}
