// frontend/packages/api-client/src/hooks/useQuotaUsageDetails.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { QuotaUsageDetailDTO, TenantID } from '../types';

/**
 * 获取租户配额使用详情
 */
export function useQuotaUsageDetails(
  tenantId: TenantID,
  options?: Omit<UseQueryOptions<QuotaUsageDetailDTO[]>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['quotas', tenantId, 'usage', 'details'],
    queryFn: async () => {
      const response = await httpClient.get<QuotaUsageDetailDTO[]>(
        API_ENDPOINTS.QUOTA.USAGE_DETAILS(tenantId)
      );
      return response.data;
    },
    enabled: !!tenantId,
    ...options,
  });
}
