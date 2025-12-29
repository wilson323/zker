// frontend/packages/api-client/src/hooks/useRoutingRules.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { RoutingRuleDTO, RoutingRuleFilter, PageResponse } from '../types';

/**
 * 获取路由规则列表
 */
export function useRoutingRules(
  filter?: RoutingRuleFilter,
  options?: Omit<UseQueryOptions<PageResponse<RoutingRuleDTO>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['routing', 'rules', filter],
    queryFn: async () => {
      const response = await httpClient.get<PageResponse<RoutingRuleDTO>>(
        API_ENDPOINTS.ROUTING_RULE.LIST,
        { params: filter }
      );
      return response.data;
    },
    ...options,
  });
}
