// frontend/packages/api-client/src/hooks/useUpdateRole.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { RoleDTO, UpdateRoleRequest } from '../types';

/**
 * 更新角色
 */
export function useUpdateRole(
  options?: Omit<
    UseMutationOptions<RoleDTO, Error, { roleId: string; data: UpdateRoleRequest }>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async ({ roleId, data }: { roleId: string; data: UpdateRoleRequest }) => {
      const response = await httpClient.put<RoleDTO>(
        API_ENDPOINTS.ROLE.UPDATE(roleId),
        data
      );
      return response.data;
    },
    ...options,
  });
}
