// frontend/packages/api-client/src/hooks/useRoleList.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { RoleDTO, RoleFilter, PageResponse } from '../types';

/**
 * 获取角色列表
 */
export function useRoleList(
  filter?: RoleFilter,
  options?: Omit<UseQueryOptions<PageResponse<RoleDTO>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['roles', 'list', filter],
    queryFn: async () => {
      const response = await httpClient.get<PageResponse<RoleDTO>>(
        API_ENDPOINTS.ROLE.LIST,
        { params: filter }
      );
      return response.data;
    },
    ...options,
  });
}
