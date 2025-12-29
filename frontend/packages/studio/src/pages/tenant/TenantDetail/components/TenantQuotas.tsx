// frontend/packages/studio/src/pages/tenant/TenantDetail/components/TenantQuotas.tsx

import React from 'react';
import { QuotaIndicator } from '@coze-studio/business-components';
import { useStyles } from './TenantQuotas.styles';

/**
 * 配额数据接口
 */
export interface QuotaData {
  resource_type: string;
  used: number;
  max: number;
  unit: string;
}

/**
 * TenantQuotas组件Props接口
 */
export interface TenantQuotasProps {
  /** 租户ID */
  tenantId: string;
}

/**
 * TenantQuotas 租户配额管理
 *
 * 显示和管理租户的配额使用情况
 */
export const TenantQuotas: React.FC<TenantQuotasProps> = ({ tenantId }) => {
  const classes = useStyles();

  // 模拟配额数据
  const quotas: QuotaData[] = [
    {
      resource_type: 'API调用次数',
      used: 750,
      max: 1000,
      unit: '次/月',
    },
    {
      resource_type: '机器人数量',
      used: 5,
      max: 10,
      unit: '个',
    },
    {
      resource_type: '知识库数量',
      used: 3,
      max: 20,
      unit: '个',
    },
    {
      resource_type: '存储空间',
      used: 2.5,
      max: 5,
      unit: 'GB',
    },
    {
      resource_type: '工作流数量',
      used: 8,
      max: 50,
      unit: '个',
    },
  ];

  return (
    <div className={classes.container}>
      <h3 className={classes.title}>配额管理</h3>

      <div className={classes.quotaList}>
        {quotas.map((quota, index) => (
          <div key={index} className={classes.quotaItem}>
            <QuotaIndicator
              used={quota.used}
              max={quota.max}
              resourceType={`${quota.resource_type}（${quota.unit}）`}
              showLabel
              showPercentage
            />
          </div>
        ))}
      </div>

      <div className={classes.footer}>
        <button className={classes.adjustButton}>
          调整配额
        </button>
      </div>
    </div>
  );
};

export default TenantQuotas;
