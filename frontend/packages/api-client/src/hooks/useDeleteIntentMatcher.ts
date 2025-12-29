// frontend/packages/api-client/src/hooks/useDeleteIntentMatcher.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';

/**
 * 删除意图匹配器
 */
export function useDeleteIntentMatcher(
  options?: Omit<
    UseMutationOptions<void, Error, string>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async (matcherId: string) => {
      await httpClient.delete(API_ENDPOINTS.INTENT_MATCHER.DELETE(matcherId));
    },
    ...options,
  });
}
