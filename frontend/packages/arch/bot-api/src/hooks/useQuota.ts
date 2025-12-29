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

import { useCallback, useEffect, useState } from 'react';

import { quotaApi } from '@coze-arch/bot-api';
import type {
  Quota,
  QuotaUsageDetail,
  ResourceType,
  CheckQuotaRequest,
  UpdateQuotaLimitRequest,
  GetQuotaUsageRequest,
} from '@coze-arch/bot-api';

/**
 * 配额列表 Hook
 */
export const useQuotas = (tenantId?: string, resourceTypes?: ResourceType[]) => {
  const [quotas, setQuotas] = useState<Quota[]>([]);
  const [overallUsage, setOverallUsage] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchQuotas = useCallback(async () => {
    if (!tenantId) return;

    setLoading(true);
    setError(null);

    try {
      const response = await quotaApi.getQuotas({
        tenant_id: tenantId,
        resource_types: resourceTypes,
      });
      setQuotas(response.quotas);
      setOverallUsage(response.overall_usage_percentage);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [tenantId, resourceTypes]);

  useEffect(() => {
    fetchQuotas();
  }, [fetchQuotas]);

  return {
    quotas,
    overallUsage,
    loading,
    error,
    refetch: fetchQuotas,
  };
};

/**
 * 配额检查 Hook
 */
export const useQuotaCheck = () => {
  const [checking, setChecking] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const checkQuota = useCallback(async (request: CheckQuotaRequest) => {
    setChecking(true);
    setError(null);

    try {
      const response = await quotaApi.checkQuota(request);
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setChecking(false);
    }
  }, []);

  return {
    checkQuota,
    checking,
    error,
  };
};

/**
 * 配额使用详情 Hook
 */
export const useQuotaUsage = (
  tenantId?: string,
  resourceType?: ResourceType,
  from?: string,
  to?: string,
) => {
  const [usage, setUsage] = useState<QuotaUsageDetail | null>(null);
  const [trend, setTrend] = useState<{
    direction: 'up' | 'down' | 'stable';
    change_percentage: number;
  } | null>(null);
  const [estimatedExhaustion, setEstimatedExhaustion] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchUsage = useCallback(async () => {
    if (!tenantId || !resourceType) return;

    setLoading(true);
    setError(null);

    try {
      const response = await quotaApi.getQuotaUsage({
        tenant_id: tenantId,
        resource_type: resourceType,
        from,
        to,
      });
      setUsage(response.usage);
      setTrend(response.trend);
      setEstimatedExhaustion(response.estimated_exhaustion_at || null);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [tenantId, resourceType, from, to]);

  useEffect(() => {
    fetchUsage();
  }, [fetchUsage]);

  return {
    usage,
    trend,
    estimatedExhaustion,
    loading,
    error,
    refetch: fetchUsage,
  };
};

/**
 * 更新配额限制 Hook
 */
export const useUpdateQuotaLimit = () => {
  const [updating, setUpdating] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const updateLimit = useCallback(async (request: UpdateQuotaLimitRequest) => {
    setUpdating(true);
    setError(null);

    try {
      const response = await quotaApi.updateQuotaLimit(request);
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setUpdating(false);
    }
  }, []);

  return {
    updateLimit,
    updating,
    error,
  };
};

/**
 * 配额警告 Hook
 * 当配额使用率超过阈值时触发警告
 */
export const useQuotaWarning = (quotas: Quota[], warningThreshold = 80, criticalThreshold = 90) => {
  const warnings = quotas.filter(q => q.usage_percentage >= warningThreshold);
  const criticals = quotas.filter(q => q.usage_percentage >= criticalThreshold);

  return {
    hasWarning: warnings.length > 0,
    hasCritical: criticals.length > 0,
    warnings,
    criticals,
    warningCount: warnings.length,
    criticalCount: criticals.length,
  };
};

/**
 * 配额统计 Hook
 */
export const useQuotaStats = (quotas: Quota[]) => {
  const stats = {
    total: quotas.length,
    healthy: quotas.filter(q => q.usage_percentage < 80).length,
    warning: quotas.filter(q => q.usage_percentage >= 80 && q.usage_percentage < 90).length,
    critical: quotas.filter(q => q.usage_percentage >= 90).length,
    averageUsage: quotas.length > 0
      ? quotas.reduce((sum, q) => sum + q.usage_percentage, 0) / quotas.length
      : 0,
  };

  return stats;
};
