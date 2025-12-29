// frontend/packages/api-client/src/hooks/useCheckPermission.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { CheckPermissionRequest, CheckPermissionResponse } from '../types';

/**
 * 检查权限
 */
export function useCheckPermission(
  options?: Omit<
    UseMutationOptions<CheckPermissionResponse, Error, CheckPermissionRequest>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async (data: CheckPermissionRequest) => {
      const response = await httpClient.post<CheckPermissionResponse>(
        API_ENDPOINTS.PERMISSION.CHECK,
        data
      );
      return response.data;
    },
    ...options,
  });
}
