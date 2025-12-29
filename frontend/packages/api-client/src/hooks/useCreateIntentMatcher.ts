// frontend/packages/api-client/src/hooks/useCreateIntentMatcher.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { IntentMatcherDTO, CreateIntentMatcherRequest } from '../types';

/**
 * 创建意图匹配器
 */
export function useCreateIntentMatcher(
  options?: Omit<
    UseMutationOptions<IntentMatcherDTO, Error, CreateIntentMatcherRequest>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async (data: CreateIntentMatcherRequest) => {
      const response = await httpClient.post<IntentMatcherDTO>(
        API_ENDPOINTS.INTENT_MATCHER.CREATE,
        data
      );
      return response.data;
    },
    ...options,
  });
}
