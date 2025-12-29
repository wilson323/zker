// frontend/packages/studio/src/pages/settings/SubscriptionManagement/SubscriptionManagement.tsx

import React, { useState } from 'react';
import { Button } from '@coze-studio/ui-components';
import { useTenantStore } from '@coze-studio/stores';
import { useTenantSubscription, useSubscriptionPlans } from '@coze-studio/api-client';
import './SubscriptionManagement.styles.css';

export const SubscriptionManagement: React.FC = () => {
  const { currentTenant } = useTenantStore();
  const [isUpgrading, setIsUpgrading] = useState(false);

  const { data: subscription, isLoading } = useTenantSubscription(
    currentTenant?.tenant_id || '',
    {
      enabled: !!currentTenant?.tenant_id,
    }
  );

  const { data: plans } = useSubscriptionPlans();

  const handleUpgrade = (planId: string) => {
    setIsUpgrading(true);
    // TODO: 实现升级逻辑
    console.log('Upgrade to plan:', planId);
    setTimeout(() => setIsUpgrading(false), 1000);
  };

  const handleCancel = () => {
    // TODO: 实现取消订阅逻辑
    console.log('Cancel subscription');
  };

  if (isLoading) {
    return <div className="subscription-loading">加载中...</div>;
  }

  if (!subscription) {
    return <div className="subscription-error">未找到订阅信息</div>;
  }

  return (
    <div className="subscription-management">
      <div className="subscription-header">
        <h1>订阅管理</h1>
        <p>管理您的订阅方案和账单信息</p>
      </div>

      <div className="subscription-current">
        <h2>当前订阅方案</h2>
        <div className="current-plan-card">
          <div className="plan-info">
            <h3>{subscription.plan_name}</h3>
            <div className="plan-tier">
              <span className={`tier-badge tier-${subscription.tier}`}>
                {subscription.tier.toUpperCase()}
              </span>
            </div>
          </div>

          <div className="plan-pricing">
            <div className="price">
              ¥{subscription.price}
              <span className="period">/{subscription.billing_cycle === 'monthly' ? '月' : '年'}</span>
            </div>
            <div className="next-billing">
              下次计费日期：{new Date(subscription.current_period_end).toLocaleDateString()}
            </div>
          </div>

          <div className="plan-actions">
            {subscription.status === 'active' && !subscription.cancel_at_period_end && (
              <Button onClick={() => handleCancel(subscription.subscription_id)}>
                取消订阅
              </Button>
            )}
            {subscription.cancel_at_period_end && (
              <div className="cancel-notice">
                ⚠️ 订阅将在 {new Date(subscription.current_period_end).toLocaleDateString()}{' '}
                后取消
              </div>
            )}
          </div>
        </div>
      </div>

      {plans && plans.length > 0 && (
        <div className="subscription-plans">
          <h2>升级订阅方案</h2>
          <div className="plans-grid">
            {plans
              .filter((plan) => plan.is_active && plan.plan_id !== subscription.plan_id)
              .map((plan) => (
                <div key={plan.plan_id} className="plan-card">
                  <div className="plan-header">
                    <h3>{plan.plan_name}</h3>
                    <span className={`tier-badge tier-${plan.tier}`}>
                      {plan.tier.toUpperCase()}
                    </span>
                  </div>

                  <div className="plan-pricing">
                    <div className="price">
                      ¥{plan.price_monthly}
                      <span className="period">/月</span>
                    </div>
                    {plan.price_yearly && (
                      <div className="yearly-price">
                        年付：¥{plan.price_yearly}
                        <span className="discount">
                          省¥{plan.price_monthly * 12 - plan.price_yearly}
                        </span>
                      </div>
                    )}
                  </div>

                  <div className="plan-features">
                    <h4>包含功能</h4>
                    <ul>
                      {plan.features.map((feature, index) => (
                        <li key={index}>✓ {feature}</li>
                      ))}
                    </ul>
                  </div>

                  <Button
                    onClick={() => handleUpgrade(plan.plan_id)}
                    disabled={isUpgrading}
                    className="upgrade-button"
                  >
                    {isUpgrading ? '处理中...' : '立即升级'}
                  </Button>
                </div>
              ))}
          </div>
        </div>
      )}

      <div className="subscription-invoices">
        <h2>发票历史</h2>
        <div className="invoices-list">
          <div className="empty-state">暂无发票记录</div>
        </div>
      </div>
    </div>
  );
};

export default SubscriptionManagement;
