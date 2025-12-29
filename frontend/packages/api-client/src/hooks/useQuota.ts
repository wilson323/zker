// frontend/packages/api-client/src/hooks/useQuota.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { QuotaDTO, TenantID } from '../types';

/**
 * 获取租户配额列表
 */
export function useQuota(
  tenantId: TenantID,
  options?: Omit<UseQueryOptions<QuotaDTO[]>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['quotas', tenantId],
    queryFn: async () => {
      const response = await httpClient.get<QuotaDTO[]>(
        API_ENDPOINTS.QUOTA.LIST,
        { params: { tenant_id: tenantId } }
      );
      return response.data;
    },
    enabled: !!tenantId,
    ...options,
  });
}
