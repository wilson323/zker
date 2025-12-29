// frontend/packages/api-client/src/hooks/useIntentMatchers.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { IntentMatcherDTO, IntentMatcherFilter, PageResponse } from '../types';

/**
 * 获取意图匹配器列表
 */
export function useIntentMatchers(
  filter?: IntentMatcherFilter,
  options?: Omit<UseQueryOptions<PageResponse<IntentMatcherDTO>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['routing', 'matchers', filter],
    queryFn: async () => {
      const response = await httpClient.get<PageResponse<IntentMatcherDTO>>(
        API_ENDPOINTS.INTENT_MATCHER.LIST,
        { params: filter }
      );
      return response.data;
    },
    ...options,
  });
}
