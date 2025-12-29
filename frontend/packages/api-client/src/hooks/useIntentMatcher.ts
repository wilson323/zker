// frontend/packages/api-client/src/hooks/useIntentMatcher.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { IntentMatcherDTO } from '../types';

/**
 * 获取单个意图匹配器
 */
export function useIntentMatcher(
  matcherId: string,
  options?: Omit<UseQueryOptions<IntentMatcherDTO>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['routing', 'matchers', matcherId],
    queryFn: async () => {
      const response = await httpClient.get<IntentMatcherDTO>(
        API_ENDPOINTS.INTENT_MATCHER.GET(matcherId)
      );
      return response.data;
    },
    enabled: !!matcherId,
    ...options,
  });
}
