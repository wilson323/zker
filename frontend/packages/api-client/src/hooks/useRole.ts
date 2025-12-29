// frontend/packages/api-client/src/hooks/useRole.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { RoleDTO } from '../types';

/**
 * 获取单个角色
 */
export function useRole(
  roleId: string,
  options?: Omit<UseQueryOptions<RoleDTO>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['roles', roleId],
    queryFn: async () => {
      const response = await httpClient.get<RoleDTO>(
        API_ENDPOINTS.ROLE.GET(roleId)
      );
      return response.data;
    },
    enabled: !!roleId,
    ...options,
  });
}
