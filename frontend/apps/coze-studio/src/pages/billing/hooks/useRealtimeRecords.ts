/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { useEffect, useState } from 'react';
import type { BillingService } from '@coze-studio/api-client';
import type { TokenRecordListResponse } from '@coze-studio/common/types/billing';

/**
 * 实时Token记录Hook
 */
export const useRealtimeRecords = (
  billingService: Promise<BillingService>,
  tenantId: string,
) => {
  const [data, setData] = useState<TokenRecordListResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);
      const service = await billingService;
      const response = await service.getTokenRecords(tenantId, { limit: 50 });
      setData(response.data);
    } catch (err) {
      setError(err as Error);
      console.error('Failed to fetch token records:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [billingService, tenantId]);

  return {
    data,
    loading,
    error,
    refetch: fetchData,
  };
};
