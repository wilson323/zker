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

/**
 * 预算管理页面
 *
 * 展示预算配置、告警设置、使用预测、告警历史
 */

import React, { useState, useMemo, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Card,
  Layout,
  Button,
  Progress,
  Table,
  Tag,
  Space,
  Typography,
  Alert,
  Spin,
  Descriptions,
  Switch,
  message,
  Popconfirm,
} from '@coze-studio/ui-components';
import { useTranslation } from '@coze-studio/i18n';
import { getBillingService } from '@coze-studio/api-client';
import type {
  BudgetSettingsDTO,
  BudgetUsageDTO,
  BudgetAlertDTO,
  UpdateBudgetRequest,
  AlertListResponse,
} from '@coze-studio/common/types/billing';
import BudgetEditModal from '../../components/billing/BudgetEditModal';

const { Header, Content } = Layout;
const { Title, Text } = Typography;

/**
 * 预算管理页面组件
 */
export const BudgetManagementPage: React.FC = () => {
  const { tenantId = 'current' } = useParams<{ tenantId: string }>();
  const { t } = useTranslation();

  // 预算设置
  const [budget, setBudget] = useState<BudgetSettingsDTO | null>(null);
  const [budgetLoading, setBudgetLoading] = useState(true);

  // 预算使用情况
  const [budgetUsage, setBudgetUsage] = useState<BudgetUsageDTO | null>(null);
  const [usageLoading, setUsageLoading] = useState(true);

  // 告警历史
  const [alerts, setAlerts] = useState<BudgetAlertDTO[]>([]);
  const [alertsLoading, setAlertsLoading] = useState(true);

  // 编辑弹窗
  const [editModalVisible, setEditModalVisible] = useState(false);

  // 获取BillingService
  const billingService = useMemo(() => getBillingService(), []);

  // 获取预算设置
  const fetchBudget = async () => {
    try {
      setBudgetLoading(true);
      const service = await billingService;
      const response = await service.getBudget(tenantId);
      setBudget(response.data);
    } catch (error) {
      console.error('Failed to fetch budget:', error);
      message.error(t('common.failed'));
    } finally {
      setBudgetLoading(false);
    }
  };

  // 获取预算使用情况
  const fetchBudgetUsage = async () => {
    try {
      setUsageLoading(true);
      const service = await billingService;
      const response = await service.getBudgetUsage(tenantId);
      setBudgetUsage(response.data);
    } catch (error) {
      console.error('Failed to fetch budget usage:', error);
      message.error(t('common.failed'));
    } finally {
      setUsageLoading(false);
    }
  };

  // 获取告警历史
  const fetchAlerts = async () => {
    try {
      setAlertsLoading(true);
      const service = await billingService;
      const response = await service.getBudgetAlerts(tenantId, { limit: 20 });
      setAlerts(response.data.alerts);
    } catch (error) {
      console.error('Failed to fetch alerts:', error);
      message.error(t('common.failed'));
    } finally {
      setAlertsLoading(false);
    }
  };

  // 初始化加载数据
  useEffect(() => {
    fetchBudget();
    fetchBudgetUsage();
    fetchAlerts();
  }, [billingService, tenantId]);

  // 保存预算设置
  const handleSaveBudget = async (data: UpdateBudgetRequest) => {
    try {
      const service = await billingService;
      await service.updateBudget(tenantId, data);
      await fetchBudget();
      setEditModalVisible(false);
    } catch (error) {
      console.error('Failed to update budget:', error);
      throw error;
    }
  };

  // 删除预算
  const handleDeleteBudget = async () => {
    try {
      const service = await billingService;
      await service.deleteBudget(tenantId);
      message.success(t('common.success'));
      setBudget(null);
    } catch (error) {
      console.error('Failed to delete budget:', error);
      message.error(t('common.failed'));
    }
  };

  // 获取进度条颜色
  const getProgressColor = (percent: number) => {
    if (percent < 80) return '#52c41a';
    if (percent < 95) return '#faad14';
    return '#f5222d';
  };

  // 获取告警级别标签颜色
  const getAlertLevelColor = (level: string) => {
    switch (level) {
      case 'warning':
        return 'warning';
      case 'critical':
        return 'error';
      case 'exceeded':
        return 'error';
      default:
        return 'default';
    }
  };

  // 告警表格列定义
  const alertColumns = [
    {
      title: t('billing.budget.alertTime'),
      dataIndex: 'created_at',
      key: 'created_at',
      render: (timestamp: number) => new Date(timestamp).toLocaleString(),
      width: 200,
    },
    {
      title: t('billing.budget.alertLevel'),
      dataIndex: 'alert_level',
      key: 'alert_level',
      render: (level: string) => (
        <Tag color={getAlertLevelColor(level)}>{t(`billing.budget.alertLevels.${level}`)}</Tag>
      ),
      width: 100,
    },
    {
      title: t('billing.budget.usagePercent'),
      dataIndex: 'usage_percent',
      key: 'usage_percent',
      render: (percent: number) => `${percent.toFixed(1)}%`,
      width: 100,
    },
    {
      title: t('billing.budget.alertMessage'),
      dataIndex: 'alert_message',
      key: 'alert_message',
      ellipsis: true,
    },
    {
      title: t('common.status'),
      dataIndex: 'is_sent',
      key: 'is_sent',
      render: (sent: boolean) => (
        <Tag color={sent ? 'success' : 'default'}>{sent ? t('common.success') : 'Pending'}</Tag>
      ),
      width: 100,
    },
  ];

  return (
    <Layout style={{ minHeight: '100vh', background: '#f0f2f5' }}>
      <Header
        style={{
          background: '#fff',
          padding: '16px 24px',
          borderBottom: '1px solid #f0f0f0',
        }}
      >
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Title level={3} style={{ margin: 0 }}>
            {t('billing.budget.title')}
          </Title>
          <Space>
            <Button onClick={fetchBudget}>{t('common.refresh')}</Button>
            {budget && (
              <>
                <Button type="primary" onClick={() => setEditModalVisible(true)}>
                  {t('billing.budget.editBudget')}
                </Button>
                <Popconfirm
                  title={t('billing.budget.deleteConfirm')}
                  onConfirm={handleDeleteBudget}
                  okText={t('common.confirm')}
                  cancelText={t('common.cancel')}
                >
                  <Button danger>{t('billing.budget.deleteBudget')}</Button>
                </Popconfirm>
              </>
            )}
          </Space>
        </div>
      </Header>

      <Content style={{ padding: '24px' }}>
        {!budget ? (
          <Card>
            <div style={{ textAlign: 'center', padding: '60px 0' }}>
              <Title level={4}>{t('billing.budget.createBudget')}</Title>
              <Button type="primary" size="large" onClick={() => setEditModalVisible(true)}>
                {t('billing.budget.createBudget')}
              </Button>
            </div>
          </Card>
        ) : (
          <>
            {/* 预算概览卡片 */}
            <Spin spinning={budgetLoading || usageLoading}>
              <Card
                title={t('billing.budget.currentBudget')}
                style={{ marginBottom: '24px' }}
                extra={
                  <Tag color={budgetUsage?.will_exceed_budget ? 'error' : 'success'}>
                    {budgetUsage?.will_exceed_budget
                      ? t('billing.budget.willExceed')
                      : t('common.status.active')}
                  </Tag>
                }
              >
                <Descriptions column={3} bordered>
                  <Descriptions.Item label={t('billing.budget.budgetType')}>
                    {t(`billing.budget.budgetTypes.${budget.budget_type}`)}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('billing.budget.budgetAmount')}>
                    ${budget.budget_amount.toFixed(2)}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('billing.budget.usedAmount')}>
                    ${budgetUsage?.used_amount.toFixed(2) || '0.00'}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('billing.budget.remainingAmount')}>
                    ${budgetUsage?.remaining_amount.toFixed(2) || '0.00'}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('billing.budget.usagePercent')} span={2}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                      <Progress
                        percent={parseFloat((budgetUsage?.usage_percent || 0).toFixed(1))}
                        strokeColor={getProgressColor(budgetUsage?.usage_percent || 0)}
                        style={{ flex: 1 }}
                      />
                      <Text strong>{(budgetUsage?.usage_percent || 0).toFixed(1)}%</Text>
                    </div>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('billing.tokenUsage.totalTokens')} span={1}>
                    {budgetUsage?.total_tokens?.toLocaleString() || 0}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('billing.tokenUsage.totalRequests')} span={2}>
                    {budgetUsage?.total_requests?.toLocaleString() || 0}
                  </Descriptions.Item>
                </Descriptions>
              </Card>
            </Spin>

            {/* 使用预测 */}
            {budgetUsage?.usage_prediction && (
              <Card title={t('billing.budget.usagePrediction')} style={{ marginBottom: '24px' }}>
                <Alert
                  message={
                    budgetUsage.will_exceed_budget ? (
                      <>
                        <Text strong>{t('billing.budget.willExceed')}</Text>
                        <br />
                        {t('billing.budget.estimatedExceedTime')}:{' '}
                        {budgetUsage.usage_prediction.estimated_exceed_time
                          ? new Date(budgetUsage.usage_prediction.estimated_exceed_time).toLocaleString()
                          : 'N/A'}
                      </>
                    ) : (
                      <Text type="success">{t('common.status.normal')}</Text>
                    )
                  }
                  description={
                    <ul style={{ marginTop: '12px', marginBottom: 0, paddingLeft: '20px' }}>
                      <li>
                        {t('billing.budget.currentDailyAvg')}: $
                        {budgetUsage.usage_prediction.current_daily_avg.toFixed(2)}
                      </li>
                      <li>
                        {t('billing.budget.predictedEndUsage')}: $
                        {budgetUsage.usage_prediction.predicted_end_usage.toFixed(2)}
                      </li>
                      {budgetUsage.usage_prediction.predicted_overage > 0 && (
                        <li>
                          {t('billing.budget.predictedOverage')}: $
                          {budgetUsage.usage_prediction.predicted_overage.toFixed(2)}
                        </li>
                      )}
                      {budgetUsage.usage_prediction.recommendations?.length > 0 && (
                        <li>
                          {t('billing.budget.recommendations')}:
                          <ul style={{ marginTop: '8px' }}>
                            {budgetUsage.usage_prediction.recommendations.map((rec, idx) => (
                              <li key={idx}>{rec}</li>
                            ))}
                          </ul>
                        </li>
                      )}
                    </ul>
                  }
                  type={budgetUsage.will_exceed_budget ? 'warning' : 'success'}
                  showIcon
                />
              </Card>
            )}

            {/* 告警配置 */}
            <Card title={t('billing.budget.alertConfig')} style={{ marginBottom: '24px' }}>
              <Descriptions column={2}>
                <Descriptions.Item label={t('billing.budget.alertThreshold1')}>
                  {budget.alert_threshold_1}%
                </Descriptions.Item>
                <Descriptions.Item label={t('billing.budget.alertThreshold2')}>
                  {budget.alert_threshold_2}%
                </Descriptions.Item>
                <Descriptions.Item label={t('billing.budget.hardCap')}>
                  <Switch checked={budget.hard_cap_enabled} disabled />
                  {budget.hard_cap_enabled && ` ($${budget.hard_cap_amount?.toFixed(2)})`}
                </Descriptions.Item>
                <Descriptions.Item label={t('billing.budget.autoDowngrade')}>
                  <Switch checked={budget.auto_downgrade_enabled} disabled />
                  {budget.auto_downgrade_enabled &&
                    ` (${budget.downgrade_config?.original_provider}/${budget.downgrade_config?.original_model} → ${budget.downgrade_config?.downgrade_provider}/${budget.downgrade_config?.downgrade_model})`}
                </Descriptions.Item>
                <Descriptions.Item label={t('billing.budget.notificationChannels')} span={2}>
                  <Space>
                    {budget.notification_channels.map(channel => (
                      <Tag key={channel}>{t(`billing.budget.channels.${channel}`)}</Tag>
                    ))}
                  </Space>
                </Descriptions.Item>
              </Descriptions>
            </Card>

            {/* 告警历史 */}
            <Card title={t('billing.budget.alertHistory')}>
              <Spin spinning={alertsLoading}>
                <Table
                  columns={alertColumns}
                  dataSource={alerts}
                  rowKey="alert_id"
                  pagination={{
                    pageSize: 10,
                    showSizeChanger: true,
                    showTotal: total => t('table.total', { total }),
                  }}
                />
              </Spin>
            </Card>
          </>
        )}
      </Content>

      {/* 编辑弹窗 */}
      <BudgetEditModal
        visible={editModalVisible}
        budget={budget}
        onSave={handleSaveBudget}
        onCancel={() => setEditModalVisible(false)}
      />
    </Layout>
  );
};

export default BudgetManagementPage;
