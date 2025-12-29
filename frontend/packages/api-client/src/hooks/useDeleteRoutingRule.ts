// frontend/packages/api-client/src/hooks/useDeleteRoutingRule.ts

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';

/**
 * 删除路由规则
 */
export function useDeleteRoutingRule(
  options?: Omit<
    UseMutationOptions<void, Error, string>,
    'mutationFn'
  >
) {
  return useMutation({
    mutationFn: async (ruleId: string) => {
      await httpClient.delete(API_ENDPOINTS.ROUTING_RULE.DELETE(ruleId));
    },
    ...options,
  });
}
