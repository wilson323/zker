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

import { request } from '@/utils/request';
import type {
  TenantOverview,
  QPSMetrics,
  ResponseTimeMetrics,
  ErrorRateMetrics,
  HealthScore,
  AlertHistory,
  AlertRule,
  AlertFilter,
  AlertStatistics,
  CreateAlertRuleRequest,
  AcknowledgeAlertRequest,
  ResolveAlertRequest,
  SilenceAlertRequest,
  TimeRangeParams,
  RecordMetricRequest,
  APIResponse,
  AgentMetrics,
  AgentPerformanceReport,
} from '@/types/monitoring';

/**
 * 监控模块API客户端
 */
export const monitoringApi = {
  // ============================================================
  // 租户监控API
  // ============================================================

  /**
   * 获取租户概览
   */
  getTenantOverview: (tenantId: string) =>
    request.get<APIResponse<TenantOverview>>(
      `/api/v1/monitoring/tenants/${tenantId}/overview`,
    ),

  /**
   * 获取QPS指标
   */
  getQPSMetrics: (tenantId: string, params: TimeRangeParams) =>
    request.get<APIResponse<QPSMetrics>>(
      `/api/v1/monitoring/tenants/${tenantId}/metrics/qps`,
      { params },
    ),

  /**
   * 获取响应时间指标
   */
  getResponseTimeMetrics: (tenantId: string, params: TimeRangeParams) =>
    request.get<APIResponse<ResponseTimeMetrics>>(
      `/api/v1/monitoring/tenants/${tenantId}/metrics/response_time`,
      { params },
    ),

  /**
   * 获取错误率指标
   */
  getErrorRateMetrics: (tenantId: string, params: TimeRangeParams) =>
    request.get<APIResponse<ErrorRateMetrics>>(
      `/api/v1/monitoring/tenants/${tenantId}/metrics/error_rate`,
      { params },
    ),

  /**
   * 获取健康评分
   */
  getHealthScore: (tenantId: string) =>
    request.get<APIResponse<HealthScore>>(
      `/api/v1/monitoring/tenants/${tenantId}/health_score`,
    ),

  /**
   * 获取监控Dashboard数据
   */
  getDashboardData: (tenantId: string) =>
    request.get<APIResponse<TenantOverview>>(
      '/api/v1/monitoring/dashboard',
      { params: { tenant_id: tenantId } },
    ),

  // ============================================================
  // 告警管理API
  // ============================================================

  /**
   * 创建告警规则
   */
  createAlertRule: (data: CreateAlertRuleRequest) =>
    request.post<APIResponse<AlertRule>>(
      '/api/v1/monitoring/alert_rules',
      data,
    ),

  /**
   * 获取告警历史
   */
  getAlertHistory: (filter: AlertFilter) =>
    request.get<APIResponse<{ alerts: AlertHistory[]; total: number }>>(
      '/api/v1/monitoring/alerts/history',
      { params: filter },
    ),

  /**
   * 确认告警
   */
  acknowledgeAlert: (alertId: string, data: AcknowledgeAlertRequest) =>
    request.post<APIResponse<void>>(
      `/api/v1/monitoring/alerts/${alertId}/acknowledge`,
      data,
    ),

  /**
   * 解决告警
   */
  resolveAlert: (alertId: string, data: ResolveAlertRequest) =>
    request.post<APIResponse<void>>(
      `/api/v1/monitoring/alerts/${alertId}/resolve`,
      data,
    ),

  /**
   * 静默告警
   */
  silenceAlert: (alertId: string, data: SilenceAlertRequest) =>
    request.post<APIResponse<void>>(
      `/api/v1/monitoring/alerts/${alertId}/silence`,
      data,
    ),

  /**
   * 获取告警统计
   */
  getAlertStatistics: (params: { tenant_id?: string; start_time?: string; end_time?: string }) =>
    request.get<APIResponse<AlertStatistics>>(
      '/api/v1/monitoring/alerts/statistics',
      { params },
    ),

  // ============================================================
  // 指标记录API
  // ============================================================

  /**
   * 记录单个指标
   */
  recordMetric: (data: RecordMetricRequest) =>
    request.post<APIResponse<void>>(
      '/api/v1/monitoring/metrics/record',
      data,
    ),

  /**
   * 批量记录指标
   */
  recordMetricBatch: (data: RecordMetricRequest[]) =>
    request.post<APIResponse<{ count: number }>>(
      '/api/v1/monitoring/metrics/batch',
      data,
    ),

  // ============================================================
  // Agent监控API
  // ============================================================

  /**
   * 获取Agent指标
   */
  getAgentMetrics: (agentId: string, metricType: string, params: TimeRangeParams) =>
    request.get<APIResponse<AgentMetrics[]>>(
      `/api/v1/monitoring/agents/${agentId}/metrics`,
      { params: { ...params, metric_type: metricType } },
    ),

  /**
   * 获取Agent性能报告
   */
  getAgentPerformanceReport: (agentId: string, params: TimeRangeParams) =>
    request.get<APIResponse<AgentPerformanceReport>>(
      `/api/v1/monitoring/agents/${agentId}/performance_report`,
      { params },
    ),

  /**
   * 对比Agent性能
   */
  compareAgents: (tenantId: string, agentIds: string[], params: TimeRangeParams) =>
    request.post<APIResponse<any[]>>(
      '/api/v1/monitoring/agents/compare',
      { tenant_id: tenantId, agent_ids: agentIds, ...params },
    ),
};

export default monitoringApi;
