// frontend/packages/api-client/src/hooks/useUpdateTenant.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { TenantDTO, UpdateTenantRequest } from '../types';

/**
 * 更新租户
 */
export function useUpdateTenant(
  options?: Omit<
    UseMutationOptions<TenantDTO, Error, { tenantId: string; data: UpdateTenantRequest }>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async ({ tenantId, data }: { tenantId: string; data: UpdateTenantRequest }) => {
      const response = await httpClient.put<TenantDTO>(
        API_ENDPOINTS.TENANT.UPDATE(tenantId),
        data
      );
      return response.data;
    },
    ...options,
  });
}
