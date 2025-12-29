// frontend/packages/api-client/src/hooks/useTenantList.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { TenantDTO, TenantFilter, PageResponse } from '../types';

/**
 * 获取租户列表
 */
export function useTenantList(
  filter?: TenantFilter,
  options?: Omit<UseQueryOptions<PageResponse<TenantDTO>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['tenants', 'list', filter],
    queryFn: async () => {
      const response = await httpClient.get<PageResponse<TenantDTO>>(
        API_ENDPOINTS.TENANT.LIST,
        { params: filter }
      );
      return response.data;
    },
    ...options,
  });
}
