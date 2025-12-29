// frontend/packages/api-client/src/hooks/useUpdateIntentMatcher.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { IntentMatcherDTO, UpdateIntentMatcherRequest } from '../types';

/**
 * 更新意图匹配器
 */
export function useUpdateIntentMatcher(
  options?: Omit<
    UseMutationOptions<IntentMatcherDTO, Error, { matcherId: string; data: UpdateIntentMatcherRequest }>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async ({ matcherId, data }: { matcherId: string; data: UpdateIntentMatcherRequest }) => {
      const response = await httpClient.put<IntentMatcherDTO>(
        API_ENDPOINTS.INTENT_MATCHER.UPDATE(matcherId),
        data
      );
      return response.data;
    },
    ...options,
  });
}
