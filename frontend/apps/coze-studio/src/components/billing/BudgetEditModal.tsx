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
 * 预算编辑弹窗组件
 *
 * 提供预算金额、告警阈值、通知配置等编辑功能
 */

import React, { useState, useEffect } from 'react';
import {
  Modal,
  Form,
  InputNumber,
  Select,
  Switch,
  Slider,
  Button,
  Space,
  Tag,
  message,
} from '@coze-studio/ui-components';
import { useTranslation } from '@coze-studio/i18n';
import type { BudgetSettingsDTO, UpdateBudgetRequest, NotificationChannel } from '@coze-studio/common/types/billing';

interface BudgetEditModalProps {
  visible: boolean;
  budget: BudgetSettingsDTO | null;
  onSave: (data: UpdateBudgetRequest) => Promise<void>;
  onCancel: () => void;
}

/**
 * 预算编辑弹窗组件
 */
export const BudgetEditModal: React.FC<BudgetEditModalProps> = ({
  visible,
  budget,
  onSave,
  onCancel,
}) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [notificationChannels, setNotificationChannels] = useState<NotificationChannel[]>([]);
  const [recipientInput, setRecipientInput] = useState('');
  const [recipients, setRecipients] = useState<Array<{ id: string; address: string }>>([]);

  // 初始化表单数据
  useEffect(() => {
    if (visible && budget) {
      form.setFieldsValue({
        budget_type: budget.budget_type,
        budget_amount: budget.budget_amount,
        alert_threshold_1: budget.alert_threshold_1,
        alert_threshold_2: budget.alert_threshold_2,
        hard_cap_enabled: budget.hard_cap_enabled,
        hard_cap_amount: budget.hard_cap_amount,
        auto_downgrade_enabled: budget.auto_downgrade_enabled,
        downgrade_original_provider: budget.downgrade_config?.original_provider,
        downgrade_original_model: budget.downgrade_config?.original_model,
        downgrade_downgrade_provider: budget.downgrade_config?.downgrade_provider,
        downgrade_downgrade_model: budget.downgrade_config?.downgrade_model,
      });
      setNotificationChannels(budget.notification_channels || []);
      setRecipients(
        budget.notification_recipients?.map(r => ({
          id: r.recipient_id,
          address: r.recipient_address,
        })) || [],
      );
    }
  }, [visible, budget, form]);

  // 处理保存
  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      const updateData: UpdateBudgetRequest = {
        budget_type: values.budget_type,
        budget_amount: values.budget_amount,
        alert_threshold_1: values.alert_threshold_1,
        alert_threshold_2: values.alert_threshold_2,
        hard_cap_enabled: values.hard_cap_enabled,
        hard_cap_amount: values.hard_cap_amount,
        auto_downgrade_enabled: values.auto_downgrade_enabled,
        notification_channels,
        notification_recipients: recipients.map(r => ({
          recipient_id: r.id,
          recipient_type: 'user' as const,
          recipient_address: r.address,
          channels: notificationChannels,
        })),
      };

      if (values.auto_downgrade_enabled) {
        updateData.downgrade_config = {
          original_provider: values.downgrade_original_provider,
          original_model: values.downgrade_original_model,
          downgrade_provider: values.downgrade_downgrade_provider,
          downgrade_model: values.downgrade_downgrade_model,
        };
      }

      await onSave(updateData);
      message.success(t('common.success'));
      form.resetFields();
      setRecipients([]);
      setNotificationChannels([]);
    } catch (error) {
      console.error('Failed to save budget:', error);
      message.error(t('common.failed'));
    } finally {
      setLoading(false);
    }
  };

  // 处理取消
  const handleCancel = () => {
    form.resetFields();
    setRecipients([]);
    setNotificationChannels([]);
    onCancel();
  };

  // 添加接收人
  const handleAddRecipient = () => {
    if (recipientInput && !recipients.find(r => r.address === recipientInput)) {
      setRecipients([...recipients, { id: `recipient-${Date.now()}`, address: recipientInput }]);
      setRecipientInput('');
    }
  };

  // 删除接收人
  const handleRemoveRecipient = (id: string) => {
    setRecipients(recipients.filter(r => r.id !== id));
  };

  return (
    <Modal
      title={t('billing.budget.editBudget')}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      width={720}
      footer={[
        <Button key="cancel" onClick={handleCancel}>
          {t('common.cancel')}
        </Button>,
        <Button key="submit" type="primary" loading={loading} onClick={handleOk}>
          {t('common.save')}
        </Button>,
      ]}
    >
      <Form form={form} layout="vertical" style={{ marginTop: '24px' }}>
        <Form.Item
          label={t('billing.budget.budgetType')}
          name="budget_type"
          rules={[{ required: true, message: t('form.requiredError') }]}
        >
          <Select>
            <Select.Option value="monthly">{t('billing.budget.budgetTypes.monthly')}</Select.Option>
            <Select.Option value="quarterly">
              {t('billing.budget.budgetTypes.quarterly')}
            </Select.Option>
            <Select.Option value="yearly">{t('billing.budget.budgetTypes.yearly')}</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          label={t('billing.budget.budgetAmount')}
          name="budget_amount"
          rules={[{ required: true, message: t('form.requiredError') }]}
        >
          <InputNumber
            style={{ width: '100%' }}
            min={0}
            precision={2}
            prefix="$"
            placeholder={t('billing.budget.budgetAmount')}
          />
        </Form.Item>

        <Form.Item
          label={`${t('billing.budget.alertThreshold1')} (%)`}
          name="alert_threshold_1"
          rules={[{ required: true, message: t('form.requiredError') }]}
        >
          <Slider min={0} max={100} marks={{ 0: '0%', 50: '50%', 100: '100%' }} />
        </Form.Item>

        <Form.Item
          label={`${t('billing.budget.alertThreshold2')} (%)`}
          name="alert_threshold_2"
          rules={[{ required: true, message: t('form.requiredError') }]}
        >
          <Slider min={0} max={100} marks={{ 0: '0%', 50: '50%', 100: '100%' }} />
        </Form.Item>

        <Form.Item
          label={t('billing.budget.hardCapEnabled')}
          name="hard_cap_enabled"
          valuePropName="checked"
        >
          <Switch />
        </Form.Item>

        <Form.Item noStyle shouldUpdate={(prev, curr) => prev.hard_cap_enabled !== curr.hard_cap_enabled}>
          {({ getFieldValue }) =>
            getFieldValue('hard_cap_enabled') ? (
              <Form.Item
                label={t('billing.budget.hardCapAmount')}
                name="hard_cap_amount"
                rules={[{ required: true, message: t('form.requiredError') }]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  precision={2}
                  prefix="$"
                  placeholder={t('billing.budget.hardCapAmount')}
                />
              </Form.Item>
            ) : null
          }
        </Form.Item>

        <Form.Item label={t('billing.budget.autoDowngradeEnabled')} name="auto_downgrade_enabled" valuePropName="checked">
          <Switch />
        </Form.Item>

        <Form.Item noStyle shouldUpdate={(prev, curr) => prev.auto_downgrade_enabled !== curr.auto_downgrade_enabled}>
          {({ getFieldValue }) =>
            getFieldValue('auto_downgrade_enabled') ? (
              <>
                <Form.Item
                  label={t('billing.budget.originalModel')}
                  name="downgrade_original_provider"
                  rules={[{ required: true, message: t('form.requiredError') }]}
                >
                  <InputNumber style={{ width: '100%' }} placeholder="Provider/Model" />
                </Form.Item>
                <Form.Item
                  label={t('billing.budget.downgradeModel')}
                  name="downgrade_downgrade_provider"
                  rules={[{ required: true, message: t('form.requiredError') }]}
                >
                  <InputNumber style={{ width: '100%' }} placeholder="Provider/Model" />
                </Form.Item>
              </>
            ) : null
          }
        </Form.Item>

        <Form.Item label={t('billing.budget.notificationChannels')}>
          <Select
            mode="multiple"
            value={notificationChannels}
            onChange={setNotificationChannels}
            placeholder={t('billing.budget.notificationChannels')}
            style={{ width: '100%' }}
          >
            <Select.Option value="email">{t('billing.budget.channels.email')}</Select.Option>
            <Select.Option value="sms">{t('billing.budget.channels.sms')}</Select.Option>
            <Select.Option value="webhook">{t('billing.budget.channels.webhook')}</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label={t('billing.budget.recipients')}>
          <Space direction="vertical" style={{ width: '100%' }}>
            <Space.Compact style={{ width: '100%' }}>
              <InputNumber
                style={{ flex: 1 }}
                value={recipientInput}
                onChange={setRecipientInput}
                placeholder="email@example.com or +1234567890"
                onPressEnter={handleAddRecipient}
              />
              <Button type="primary" onClick={handleAddRecipient}>
                {t('billing.budget.addRecipient')}
              </Button>
            </Space.Compact>
            <div>
              {recipients.map(recipient => (
                <Tag
                  key={recipient.id}
                  closable
                  onClose={() => handleRemoveRecipient(recipient.id)}
                  style={{ marginBottom: '8px' }}
                >
                  {recipient.address}
                </Tag>
              ))}
            </div>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default BudgetEditModal;
