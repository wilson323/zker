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

import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Tag,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Switch,
  message,
  Popconfirm,
  Tooltip,
  Badge,
  Row,
  Col,
  Statistic,
} from 'antd';
import {
  PlusOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  BellOutlined,
  DeleteOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { monitoringApi } from '@/api/monitoring';
import type {
  AlertHistory,
  AlertRule,
  AlertFilter,
  AlertStatistics,
} from '@/types/monitoring';
import type { ColumnsType } from 'antd/es/table';

/**
 * 告警管理页面
 * 提供告警规则配置、告警历史查询、告警确认/解决等功能
 */
export const AlertManagementPage: React.FC = () => {
  // 状态管理
  const [loading, setLoading] = useState(false);
  const [alerts, setAlerts] = useState<AlertHistory[]>([]);
  const [total, setTotal] = useState(0);
  const [statistics, setStatistics] = useState<AlertStatistics | null>(null);
  const [filter, setFilter] = useState<AlertFilter>({
    page_size: 20,
    page_token: '',
  });

  // Modal状态
  const [ruleModalVisible, setRuleModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [selectedAlert, setSelectedAlert] = useState<AlertHistory | null>(null);
  const [form] = Form.useForm();

  // 加载告警历史
  const loadAlerts = async () => {
    setLoading(true);
    try {
      const res = await monitoringApi.getAlertHistory(filter);
      setAlerts(res.data.alerts);
      setTotal(res.data.total);
    } catch (err: any) {
      message.error('加载告警历史失败: ' + (err.message || '未知错误'));
    } finally {
      setLoading(false);
    }
  };

  // 加载告警统计
  const loadStatistics = async () => {
    try {
      const res = await monitoringApi.getAlertStatistics({
        tenant_id: filter.tenant_id,
      });
      setStatistics(res.data);
    } catch (err: any) {
      console.error('Failed to load alert statistics:', err);
    }
  };

  // 初始加载
  useEffect(() => {
    loadAlerts();
    loadStatistics();
  }, [filter]);

  // 确认告警
  const handleAcknowledge = async (alertId: string) => {
    try {
      await monitoringApi.acknowledgeAlert(alertId, {
        user_id: 'current-user-id', // 从上下文获取
      });
      message.success('告警已确认');
      loadAlerts();
      loadStatistics();
    } catch (err: any) {
      message.error('确认告警失败: ' + (err.message || '未知错误'));
    }
  };

  // 解决告警
  const handleResolve = async (alertId: string, note: string) => {
    try {
      await monitoringApi.resolveAlert(alertId, {
        user_id: 'current-user-id',
        note,
      });
      message.success('告警已解决');
      loadAlerts();
      loadStatistics();
    } catch (err: any) {
      message.error('解决告警失败: ' + (err.message || '未知错误'));
    }
  };

  // 静默告警
  const handleSilence = async (alertId: string, duration: number) => {
    try {
      await monitoringApi.silenceAlert(alertId, {
        duration_minutes: duration,
      });
      message.success('告警已静默');
      loadAlerts();
    } catch (err: any) {
      message.error('静默告警失败: ' + (err.message || '未知错误'));
    }
  };

  // 创建告警规则
  const handleCreateRule = async (values: any) => {
    try {
      await monitoringApi.createAlertRule(values);
      message.success('告警规则创建成功');
      setRuleModalVisible(false);
      form.resetFields();
    } catch (err: any) {
      message.error('创建告警规则失败: ' + (err.message || '未知错误'));
    }
  };

  // 查看告警详情
  const handleViewDetail = (alert: AlertHistory) => {
    setSelectedAlert(alert);
    setDetailModalVisible(true);
  };

  // 渲染严重级别标签
  const renderSeverityTag = (severity: string) => {
    const config: Record<string, { color: string; text: string }> = {
      emergency: { color: 'red', text: '紧急' },
      critical: { color: 'orange', text: '严重' },
      warning: { color: 'gold', text: '警告' },
    };
    const { color, text } = config[severity] || { color: 'default', text: severity };
    return <Tag color={color}>{text}</Tag>;
  };

  // 渲染状态标签
  const renderStatusTag = (status: string) => {
    const config: Record<string, { color: string; text: string }> = {
      pending: { color: 'red', text: '待处理' },
      acknowledged: { color: 'orange', text: '已确认' },
      resolved: { color: 'green', text: '已解决' },
      silenced: { color: 'default', text: '已静默' },
    };
    const { color, text } = config[status] || { color: 'default', text: status };
    return <Tag color={color}>{text}</Tag>;
  };

  // 表格列定义
  const columns: ColumnsType<AlertHistory> = [
    {
      title: '告警ID',
      dataIndex: 'id',
      key: 'id',
      width: 200,
      ellipsis: true,
    },
    {
      title: '规则名称',
      dataIndex: 'alert_rule_name',
      key: 'alert_rule_name',
      width: 150,
    },
    {
      title: '严重级别',
      dataIndex: 'severity',
      key: 'severity',
      width: 100,
      render: renderSeverityTag,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: renderStatusTag,
    },
    {
      title: '告警消息',
      dataIndex: 'alert_message',
      key: 'alert_message',
      ellipsis: true,
    },
    {
      title: '触发时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="查看详情">
            <Button
              type="link"
              icon={<EyeOutlined />}
              onClick={() => handleViewDetail(record)}
            />
          </Tooltip>
          {record.status === 'pending' && (
            <>
              <Tooltip title="确认告警">
                <Button
                  type="link"
                  icon={<CheckCircleOutlined />}
                  onClick={() => handleAcknowledge(record.id)}
                />
              </Tooltip>
              <Popconfirm
                title="确认要静默此告警吗？"
                description="静默后将不再发送通知，直到静默期结束"
                onConfirm={() => handleSilence(record.id, 60)}
              >
                <Button type="link" icon={<BellOutlined />} />
              </Popconfirm>
              <Popconfirm
                title="确认要解决此告警吗？"
                onConfirm={() => handleResolve(record.id, '已手动解决')}
              >
                <Button type="link" icon={<CloseCircleOutlined />} danger />
              </Popconfirm>
            </>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      {/* 告警统计卡片 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总告警数"
              value={statistics?.total_count || 0}
              prefix={<BellOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="待处理"
              value={statistics?.pending_count || 0}
              valueStyle={{ color: '#ff4d4f' }}
              prefix={<Badge status="error" />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="已确认"
              value={statistics?.acknowledged_count || 0}
              valueStyle={{ color: '#faad14' }}
              prefix={<Badge status="warning" />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="已解决"
              value={statistics?.resolved_count || 0}
              valueStyle={{ color: '#52c41a' }}
              prefix={<Badge status="success" />}
            />
          </Card>
        </Col>
      </Row>

      {/* 告警列表 */}
      <Card
        title="告警历史"
        extra={
          <Space>
            <Select
              placeholder="筛选状态"
              style={{ width: 120 }}
              allowClear
              onChange={(value) => setFilter({ ...filter, status: value })}
            >
              <Select.Option value="pending">待处理</Select.Option>
              <Select.Option value="acknowledged">已确认</Select.Option>
              <Select.Option value="resolved">已解决</Select.Option>
              <Select.Option value="silenced">已静默</Select.Option>
            </Select>
            <Select
              placeholder="筛选级别"
              style={{ width: 120 }}
              allowClear
              onChange={(value) => setFilter({ ...filter, severity: value })}
            >
              <Select.Option value="warning">警告</Select.Option>
              <Select.Option value="critical">严重</Select.Option>
              <Select.Option value="emergency">紧急</Select.Option>
            </Select>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setRuleModalVisible(true)}
            >
              创建规则
            </Button>
          </Space>
        }
      >
        <Table
          columns={columns}
          dataSource={alerts}
          rowKey="id"
          loading={loading}
          pagination={{
            total,
            pageSize: filter.page_size,
            current: parseInt(filter.page_token || '1') + 1,
            onChange: (page, pageSize) => {
              setFilter({
                ...filter,
                page_token: String(page - 1),
                page_size: pageSize,
              });
            },
          }}
        />
      </Card>

      {/* 创建告警规则Modal */}
      <Modal
        title="创建告警规则"
        open={ruleModalVisible}
        onCancel={() => setRuleModalVisible(false)}
        onOk={() => form.submit()}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreateRule}
        >
          <Form.Item
            name="tenant_id"
            label="租户ID"
            rules={[{ required: true, message: '请输入租户ID' }]}
          >
            <Input placeholder="请输入租户ID" />
          </Form.Item>

          <Form.Item
            name="rule_name"
            label="规则名称"
            rules={[{ required: true, message: '请输入规则名称' }]}
          >
            <Input placeholder="例如：高错误率告警" />
          </Form.Item>

          <Form.Item
            name="metric_type"
            label="监控指标"
            rules={[{ required: true, message: '请选择监控指标' }]}
          >
            <Select placeholder="请选择监控指标">
              <Select.Option value="qps">QPS</Select.Option>
              <Select.Option value="response_time">响应时间</Select.Option>
              <Select.Option value="error_rate">错误率</Select.Option>
            </Select>
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="threshold"
                label="阈值"
                rules={[{ required: true, message: '请输入阈值' }]}
              >
                <InputNumber style={{ width: '100%' }} placeholder="阈值" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="comparison"
                label="比较方式"
                rules={[{ required: true, message: '请选择比较方式' }]}
              >
                <Select placeholder="请选择">
                  <Select.Option value="gt">大于 (>)</Select.Option>
                  <Select.Option value="lt">小于 (&lt;)</Select.Option>
                  <Select.Option value="eq">等于 (=)</Select.Option>
                  <Select.Option value="gte">大于等于 (≥)</Select.Option>
                  <Select.Option value="lte">小于等于 (≤)</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="severity"
            label="告警级别"
            rules={[{ required: true, message: '请选择告警级别' }]}
          >
            <Select placeholder="请选择告警级别">
              <Select.Option value="warning">警告</Select.Option>
              <Select.Option value="critical">严重</Select.Option>
              <Select.Option value="emergency">紧急</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="notification_channels"
            label="通知渠道"
            rules={[{ required: true, message: '请选择通知渠道' }]}
          >
            <Select mode="multiple" placeholder="请选择通知渠道">
              <Select.Option value="email">邮件</Select.Option>
              <Select.Option value="sms">短信</Select.Option>
              <Select.Option value="webhook">Webhook</Select.Option>
              <Select.Option value="slack">Slack</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="is_enabled"
            label="启用规则"
            valuePropName="checked"
            initialValue={true}
          >
            <Switch />
          </Form.Item>
        </Form>
      </Modal>

      {/* 告警详情Modal */}
      <Modal
        title="告警详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          selectedAlert?.status === 'pending' && (
            <Button
              key="resolve"
              type="primary"
              danger
              onClick={() => {
                if (selectedAlert) {
                  handleResolve(selectedAlert.id, '已手动解决');
                  setDetailModalVisible(false);
                }
              }}
            >
              标记为已解决
            </Button>
          ),
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            关闭
          </Button>,
        ]}
        width={700}
      >
        {selectedAlert && (
          <div>
            <Row gutter={16}>
              <Col span={12}>
                <p>
                  <strong>告警ID:</strong> {selectedAlert.id}
                </p>
              </Col>
              <Col span={12}>
                <p>
                  <strong>规则名称:</strong> {selectedAlert.alert_rule_name}
                </p>
              </Col>
            </Row>
            <Row gutter={16}>
              <Col span={12}>
                <p>
                  <strong>严重级别:</strong> {renderSeverityTag(selectedAlert.severity)}
                </p>
              </Col>
              <Col span={12}>
                <p>
                  <strong>状态:</strong> {renderStatusTag(selectedAlert.status)}
                </p>
              </Col>
            </Row>
            <p>
              <strong>告警消息:</strong>
            </p>
            <p>{selectedAlert.alert_message}</p>
            <p>
              <strong>告警数据:</strong>
            </p>
            <pre style={{ background: '#f5f5f5', padding: 16 }}>
              {selectedAlert.alert_data}
            </pre>
            <p>
              <strong>触发时间:</strong>{' '}
              {new Date(selectedAlert.created_at).toLocaleString('zh-CN')}
            </p>
            {selectedAlert.acknowledged_at && (
              <p>
                <strong>确认时间:</strong>{' '}
                {new Date(selectedAlert.acknowledged_at).toLocaleString('zh-CN')}
              </p>
            )}
            {selectedAlert.resolved_at && (
              <p>
                <strong>解决时间:</strong>{' '}
                {new Date(selectedAlert.resolved_at).toLocaleString('zh-CN')}
              </p>
            )}
            {selectedAlert.resolution_note && (
              <p>
                <strong>解决备注:</strong> {selectedAlert.resolution_note}
              </p>
            )}
          </div>
        )}
      </Modal>
    </div>
  );
};

export default AlertManagementPage;
