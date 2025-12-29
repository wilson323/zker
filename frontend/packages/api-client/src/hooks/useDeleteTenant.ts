// frontend/packages/api-client/src/hooks/useDeleteTenant.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { TenantDTO } from '../types';

/**
 * 删除租户
 */
export function useDeleteTenant(
  options?: Omit<
    UseMutationOptions<void, Error, string>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async (tenantId: string) => {
      await httpClient.delete(API_ENDPOINTS.TENANT.DELETE(tenantId));
    },
    ...options,
  });
}
