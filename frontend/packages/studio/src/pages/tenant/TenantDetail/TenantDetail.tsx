// frontend/packages/studio/src/pages/tenant/TenantDetail/TenantDetail.tsx

import React, { useState } from 'react';
import { Button } from '@coze-studio/ui-components';
import { TenantBasicInfo, TenantData } from './components/TenantBasicInfo';
import { TenantSubscription, SubscriptionData } from './components/TenantSubscription';
import { TenantQuotas } from './components/TenantQuotas';
import { useStyles } from './TenantDetail.styles';

/**
 * TenantDetail 租户详情页
 *
 * 显示租户的详细信息，支持多标签页切换
 *
 * @example
 * ```tsx
 * <TenantDetail tenantId="123" />
 * ```
 */
export const TenantDetail: React.FC<{ tenantId?: string }> = ({ tenantId }) => {
  const classes = useStyles();
  const [activeTab, setActiveTab] = useState<'basic' | 'subscription' | 'quotas'>('basic');

  // 模拟租户数据
  const mockTenant: TenantData = {
    tenant_id: tenantId || '1',
    tenant_name: '示例租户',
    tenant_type: 'enterprise',
    status: 'active',
    contact_email: 'admin@example.com',
    contact_phone: '13800138000',
    description: '这是一个示例租户',
  };

  const mockSubscription: SubscriptionData = {
    subscription_tier: 'pro',
    start_date: '2024-01-01T00:00:00Z',
    end_date: '2025-01-01T00:00:00Z',
    auto_renew: true,
  };

  const [tenant, setTenant] = useState<TenantData>(mockTenant);
  const [subscription, setSubscription] = useState<SubscriptionData>(mockSubscription);

  const handleBack = () => {
    // 返回列表页
    window.history.back();
  };

  const handleUpdateTenant = (data: Partial<TenantData>) => {
    setTenant({ ...tenant, ...data });
    console.log('Update tenant:', data);
  };

  const handleUpdateSubscription = (data: Partial<SubscriptionData>) => {
    setSubscription({ ...subscription, ...data });
    console.log('Update subscription:', data);
  };

  const tabs = [
    { key: 'basic' as const, label: '基本信息' },
    { key: 'subscription' as const, label: '订阅管理' },
    { key: 'quotas' as const, label: '配额管理' },
  ];

  return (
    <div className={classes.container}>
      {/* 页面头部 */}
      <div className={classes.header}>
        <div className={classes.headerLeft}>
          <button className={classes.backButton} onClick={handleBack}>
            ← 返回
          </button>
          <h1 className={classes.title}>{tenant.tenant_name}</h1>
        </div>
        <div className={classes.headerRight}>
          <span className={`${classes.status} ${classes.statusActive}`}>
            正常
          </span>
        </div>
      </div>

      {/* 标签页导航 */}
      <div className={classes.tabs}>
        {tabs.map((tab) => (
          <button
            key={tab.key}
            className={`${classes.tab} ${activeTab === tab.key ? classes.activeTab : ''}`}
            onClick={() => setActiveTab(tab.key)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* 标签页内容 */}
      <div className={classes.tabContent}>
        {activeTab === 'basic' && (
          <TenantBasicInfo tenant={tenant} onUpdate={handleUpdateTenant} />
        )}

        {activeTab === 'subscription' && (
          <TenantSubscription subscription={subscription} onUpdate={handleUpdateSubscription} />
        )}

        {activeTab === 'quotas' && (
          <TenantQuotas tenantId={tenantId || ''} />
        )}
      </div>
    </div>
  );
};

export default TenantDetail;
