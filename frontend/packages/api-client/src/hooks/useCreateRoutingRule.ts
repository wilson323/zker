// frontend/packages/api-client/src/hooks/useCreateRoutingRule.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { RoutingRuleDTO, CreateRoutingRuleRequest } from '../types';

/**
 * 创建路由规则
 */
export function useCreateRoutingRule(
  options?: Omit<
    UseMutationOptions<RoutingRuleDTO, Error, CreateRoutingRuleRequest>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async (data: CreateRoutingRuleRequest) => {
      const response = await httpClient.post<RoutingRuleDTO>(
        API_ENDPOINTS.ROUTING_RULE.CREATE,
        data
      );
      return response.data;
    },
    ...options,
  });
}
