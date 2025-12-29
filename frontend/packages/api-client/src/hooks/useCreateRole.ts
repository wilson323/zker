// frontend/packages/api-client/src/hooks/useCreateRole.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { RoleDTO, CreateRoleRequest } from '../types';

/**
 * 创建角色
 */
export function useCreateRole(
  options?: Omit<
    UseMutationOptions<RoleDTO, Error, CreateRoleRequest>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async (data: CreateRoleRequest) => {
      const response = await httpClient.post<RoleDTO>(
        API_ENDPOINTS.ROLE.CREATE,
        data
      );
      return response.data;
    },
    ...options,
  });
}
