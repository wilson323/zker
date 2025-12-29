// frontend/packages/api-client/src/hooks/useRoutingRule.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { RoutingRuleDTO } from '../types';

/**
 * 获取单个路由规则
 */
export function useRoutingRule(
  ruleId: string,
  options?: Omit<UseQueryOptions<RoutingRuleDTO>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['routing', 'rules', ruleId],
    queryFn: async () => {
      const response = await httpClient.get<RoutingRuleDTO>(
        API_ENDPOINTS.ROUTING_RULE.GET(ruleId)
      );
      return response.data;
    },
    enabled: !!ruleId,
    ...options,
  });
}
