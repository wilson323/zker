// frontend/packages/api-client/src/hooks/usePermissionList.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { PermissionDTO } from '../types';

/**
 * 获取权限列表
 */
export function usePermissionList(
  options?: Omit<UseQueryOptions<PermissionDTO[]>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['permissions', 'list'],
    queryFn: async () => {
      const response = await httpClient.get<PermissionDTO[]>(
        API_ENDPOINTS.PERMISSION.LIST
      );
      return response.data;
    },
    ...options,
  });
}
