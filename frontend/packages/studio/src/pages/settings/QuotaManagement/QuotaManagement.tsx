// frontend/packages/studio/src/pages/settings/QuotaManagement/QuotaManagement.tsx

import React, { useState } from 'react';
import { Button, Card } from '@coze-studio/ui-components';
import { QuotaEditor } from '@coze-studio/business-components';
import { useTenantStore } from '@coze-studio/stores';
import { useQuota, useQuotaUsageDetails } from '@coze-studio/api-client';
import { QuotaResourceType, QuotaUnit, QuotaLimit } from '@coze-studio/api-client';
import './QuotaManagement.styles.css';

export const QuotaManagement: React.FC = () => {
  const { currentTenant } = useTenantStore();
  const [isEditing, setIsEditing] = useState(false);

  const { data: quotas, isLoading: quotasLoading } = useQuota(
    currentTenant?.tenant_id || '',
    {
      enabled: !!currentTenant?.tenant_id,
    }
  );

  const { data: usageDetails, isLoading: usageLoading } = useQuotaUsageDetails(
    currentTenant?.tenant_id || '',
    {
      enabled: !!currentTenant?.tenant_id,
    }
  );

  const handleSaveQuotas = (updatedQuotas: QuotaLimit[]) => {
    // TODO: 实现保存逻辑
    console.log('Save quotas:', updatedQuotas);
    setIsEditing(false);
  };

  if (quotasLoading || usageLoading) {
    return <div className="quota-loading">加载中...</div>;
  }

  if (!quotas || !usageDetails) {
    return <div className="quota-error">未找到配额信息</div>;
  }

  return (
    <div className="quota-management">
      <div className="quota-header">
        <h1>配额管理</h1>
        <p>查看和管理您的资源配额使用情况</p>
      </div>

      {/* 配额使用概览 */}
      <div className="quota-overview">
        <h2>配额使用概览</h2>
        <div className="overview-cards">
          {usageDetails.map((detail) => (
            <Card key={detail.resource_type} className="overview-card">
              <div className="card-header">
                <h3>{getResourceTypeLabel(detail.resource_type)}</h3>
                <div className="usage-percentage">
                  <span className={`${detail.is_over_limit ? 'over-limit' : ''}`}>
                    {detail.usage_percentage.toFixed(1)}%
                  </span>
                </div>
              </div>

              <div className="usage-info">
                <div className="usage-bar">
                  <div className="bar-fill" style={{ width: `${detail.usage_percentage}%` }} />
                </div>
                <div className="usage-text">
                  已使用：{detail.current_usage} / {detail.limit}
                  {detail.is_over_limit && (
                    <span className="over-limit-badge">已超限</span>
                  )}
                </div>
              </div>

              {detail.is_over_limit && (
                <div className="overage-fee">
                  超量费用：¥{detail.overage_fee || '计算中'}
                </div>
              )}

              <div className="reset-info">
                重置时间：{new Date(detail.reset_at).toLocaleString()}
              </div>
            </Card>
          ))}
        </div>
      </div>

      {/* 配额编辑 */}
      <div className="quota-editor-section">
        <div className="section-header">
          <h2>配额设置</h2>
          {!isEditing && (
            <Button onClick={() => setIsEditing(true)}>编辑配额</Button>
          )}
        </div>

        {isEditing ? (
          <div className="editing-mode">
            <QuotaEditor
              quotas={quotas.map((q) => ({
                resource_type: q.resource_type as QuotaResourceType,
                limit: q.limit,
                unit: q.unit as QuotaUnit,
                is_soft_limit: q.is_soft_limit,
                overage_fee: q.overage_fee,
              }))}
              onChange={handleSaveQuotas}
              showPrice
            />
            <div className="edit-actions">
              <Button onClick={() => setIsEditing(false)}>取消</Button>
              <Button onClick={() => {/* 保存逻辑已在上面的 onChange 中处理 */}}>
                保存更改
              </Button>
            </div>
          </div>
        ) : (
          <div className="quota-list">
            {quotas.map((quota) => (
              <Card key={quota.quota_id} className="quota-item">
                <div className="quota-info">
                  <h3>{getResourceTypeLabel(quota.resource_type)}</h3>
                  <p>限制：{quota.limit} 单位</p>
                  <p>软限制：{quota.is_soft_limit ? '是' : '否'}</p>
                  {quota.is_soft_limit && <p>超量费用：¥{quota.overage_fee}</p>}
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>

      {/* 配额预警规则 */}
      <div className="quota-alerts">
        <h2>配额预警</h2>
        <Card className="alert-card">
          <p>当配额使用率达到阈值时，系统将发送预警通知。</p>
          <div className="alert-rules">
            <div className="alert-rule-item">
              <span className="rule-label">机器人数量预警</span>
              <span className="rule-value">80%</span>
            </div>
            <div className="alert-rule-item">
              <span className="rule-label">API调用预警</span>
              <span className="rule-value">90%</span>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

// 辅助函数：获取资源类型标签
function getResourceTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    BOT: '机器人数量',
    USER: '用户数量',
    API_CALL: 'API调用次数',
    MESSAGE: '消息数量',
    TOKEN: 'Token使用量',
    STORAGE: '存储空间',
    KNOWLEDGE_BASE: '知识库数量',
    WORKFLOW: '工作流数量',
  };
  return labels[type] || type;
}

export default QuotaManagement;
