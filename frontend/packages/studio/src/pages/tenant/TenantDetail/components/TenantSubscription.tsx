// frontend/packages/studio/src/pages/tenant/TenantDetail/components/TenantSubscription.tsx

import React from 'react';
import { QuotaIndicator } from '@coze-studio/business-components';
import { useStyles } from './TenantSubscription.styles';

/**
 * 订阅数据接口
 */
export interface SubscriptionData {
  subscription_tier: 'free' | 'pro' | 'enterprise';
  start_date: string;
  end_date?: string;
  auto_renew: boolean;
}

/**
 * TenantSubscription组件Props接口
 */
export interface TenantSubscriptionProps {
  /** 订阅数据 */
  subscription: SubscriptionData;
  /** 更新回调 */
  onUpdate?: (data: Partial<SubscriptionData>) => void;
}

/**
 * TenantSubscription 租户订阅管理
 *
 * 显示和管理租户的订阅信息
 */
export const TenantSubscription: React.FC<TenantSubscriptionProps> = ({
  subscription,
  onUpdate,
}) => {
  const classes = useStyles();

  const tierLabels = {
    free: '免费版',
    pro: '专业版',
    enterprise: '企业版',
  };

  const tierColors = {
    free: '#8C8C8C',
    pro: '#1890FF',
    enterprise: '#722ED1',
  };

  const handleChangeTier = (newTier: SubscriptionData['subscription_tier']) => {
    onUpdate?.({ subscription_tier: newTier });
  };

  const handleToggleAutoRenew = () => {
    onUpdate?.({ auto_renew: !subscription.auto_renew });
  };

  return (
    <div className={classes.container}>
      <h3 className={classes.title}>订阅管理</h3>

      <div className={classes.content}>
        {/* 当前订阅等级 */}
        <div className={classes.section}>
          <label className={classes.label}>当前订阅等级</label>
          <div className={classes.tierBadge} style={{ backgroundColor: tierColors[subscription.subscription_tier] }}>
            {tierLabels[subscription.subscription_tier]}
          </div>
        </div>

        {/* 订阅时间 */}
        <div className={classes.section}>
          <label className={classes.label}>订阅开始时间</label>
          <div className={classes.value}>
            {new Date(subscription.start_date).toLocaleString('zh-CN')}
          </div>
        </div>

        {subscription.end_date && (
          <div className={classes.section}>
            <label className={classes.label}>订阅结束时间</label>
            <div className={classes.value}>
              {new Date(subscription.end_date).toLocaleString('zh-CN')}
            </div>
          </div>
        )}

        {/* 自动续费 */}
        <div className={classes.section}>
          <label className={classes.label}>自动续费</label>
          <label className={classes.switch}>
            <input
              type="checkbox"
              checked={subscription.auto_renew}
              onChange={handleToggleAutoRenew}
            />
            <span className={classes.slider}></span>
          </label>
        </div>

        {/* 配额使用情况 */}
        <div className={classes.quotaSection}>
          <h4 className={classes.quotaTitle}>配额使用情况</h4>

          <QuotaIndicator
            used={750}
            max={1000}
            resourceType="API调用次数（次/月）"
            showLabel
            showPercentage
          />

          <QuotaIndicator
            used={5}
            max={10}
            resourceType="机器人数量（个）"
            showLabel
            showPercentage
          />

          <QuotaIndicator
            used={2.5}
            max={5}
            resourceType="存储空间（GB）"
            showLabel
            showPercentage
          />
        </div>

        {/* 升级订阅按钮 */}
        <div className={classes.footer}>
          <button className={classes.upgradeButton}>
            升级订阅
          </button>
        </div>
      </div>
    </div>
  );
};

export default TenantSubscription;
