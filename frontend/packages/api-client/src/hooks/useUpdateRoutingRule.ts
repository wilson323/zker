// frontend/packages/api-client/src/hooks/useUpdateRoutingRule.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { RoutingRuleDTO, UpdateRoutingRuleRequest } from '../types';

/**
 * 更新路由规则
 */
export function useUpdateRoutingRule(
  options?: Omit<
    UseMutationOptions<RoutingRuleDTO, Error, { ruleId: string; data: UpdateRoutingRuleRequest }>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async ({ ruleId, data }: { ruleId: string; data: UpdateRoutingRuleRequest }) => {
      const response = await httpClient.put<RoutingRuleDTO>(
        API_ENDPOINTS.ROUTING_RULE.UPDATE(ruleId),
        data
      );
      return response.data;
    },
    ...options,
  });
}
