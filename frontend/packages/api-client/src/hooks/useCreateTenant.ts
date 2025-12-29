// frontend/packages/api-client/src/hooks/useCreateTenant.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { TenantDTO, CreateTenantRequest } from '../types';

/**
 * 创建租户
 */
export function useCreateTenant(
  options?: Omit<
    UseMutationOptions<TenantDTO, Error, CreateTenantRequest>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async (data: CreateTenantRequest) => {
      const response = await httpClient.post<TenantDTO>(
        API_ENDPOINTS.TENANT.CREATE,
        data
      );
      return response.data;
    },
    ...options,
  });
}
