// frontend/packages/api-client/src/hooks/useDeleteRole.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { RoleDTO } from '../types';

/**
 * 删除角色
 */
export function useDeleteRole(
  options?: Omit<
    UseMutationOptions<void, Error, string>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async (roleId: string) => {
      await httpClient.delete(API_ENDPOINTS.ROLE.DELETE(roleId));
    },
    ...options,
  });
}
