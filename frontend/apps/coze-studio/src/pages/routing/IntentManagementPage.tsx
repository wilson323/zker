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
 * 意图管理页面
 *
 * 功能：
 * - 意图列表展示（支持筛选、搜索、分页）
 * - 意图创建/编辑对话框
 * - 示例句子管理
 * - 意图测试工具
 * - 批量识别测试
 */

import React, { useState, useMemo, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import {
  Card,
  Layout,
  Button,
  Table,
  Tag,
  Space,
  Typography,
  Alert,
  Spin,
  Input,
  Select,
  Switch,
  Modal,
  Form,
  message,
  Popconfirm,
  Tabs,
  Descriptions,
  Progress,
  Statistic,
  Row,
  Col,
  Tooltip,
} from '@coze-studio/ui-components';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  TestOutlined,
  RocketOutlined,
  BulbOutlined,
} from '@coze-studio/ui-icons';
import { useTranslation } from '@coze-studio/i18n';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip as RechartsTooltip, Legend, ResponsiveContainer } from 'recharts';

const { Header, Content } = Layout;
const { Title, Text, Paragraph } = Typography;
const { Option } = Select;
const { TabPane } = Tabs;

/**
 * 意图信息接口
 */
interface Intent {
  intent_id: string;
  intent_name: string;
  description: string;
  agent_id: string;
  workflow_id?: string;
  confidence: number;
  is_active: boolean;
  examples: string[];
  created_at: number;
  updated_at: number;
}

/**
 * 意图识别结果接口
 */
interface RecognitionResult {
  intent_id: string;
  intent_name: string;
  confidence: number;
  agent_id: string;
  workflow_id?: string;
  match_method: 'llm' | 'vector' | 'hybrid';
  entities: Record<string, string[]>;
}

/**
 * 意图管理页面组件
 */
export const IntentManagementPage: React.FC = () => {
  const { tenantId = 'current' } = useParams<{ tenantId: string }>();
  const { t } = useTranslation();

  // 状态管理
  const [intents, setIntents] = useState<Intent[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  // 筛选状态
  const [filterActive, setFilterActive] = useState<boolean | undefined>();
  const [searchKeyword, setSearchKeyword] = useState('');

  // 编辑/创建对话框
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [currentIntent, setCurrentIntent] = useState<Intent | null>(null);
  const [form] = Form.useForm();

  // 测试对话框
  const [testModalVisible, setTestModalVisible] = useState(false);
  const [testInput, setTestInput] = useState('');
  const [testResults, setTestResults] = useState<RecognitionResult | null>(null);
  const [testLoading, setTestLoading] = useState(false);

  // 训练对话框
  const [trainModalVisible, setTrainModalVisible] = useState(false);
  const [trainSamples, setTrainSamples] = useState<string[]>([]);
  const [trainProgress, setTrainProgress] = useState(0);
  const [trainLoading, setTrainLoading] = useState(false);

  // 统计数据
  const [stats, setStats] = useState({
    totalIntents: 0,
    activeIntents: 0,
    totalSamples: 0,
    avgAccuracy: 0,
  });

  /**
   * 获取意图列表
   */
  const fetchIntents = useCallback(async () => {
    try {
      setLoading(true);
      // TODO: 调用API获取意图列表
      // const response = await routingService.getIntents(tenantId, {
      //   is_active: filterActive,
      //   page,
      //   page_size: pageSize,
      // });
      // setIntents(response.data.intents);
      // setTotal(response.data.total);

      // 模拟数据
      setIntents([]);
      setTotal(0);
    } catch (error) {
      console.error('Failed to fetch intents:', error);
      message.error(t('common.failed'));
    } finally {
      setLoading(false);
    }
  }, [tenantId, filterActive, page, pageSize, t]);

  /**
   * 获取统计数据
   */
  const fetchStats = useCallback(async () => {
    try {
      // TODO: 调用API获取统计数据
      // const response = await routingService.getIntentStats(tenantId);
      // setStats(response.data);
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    }
  }, [tenantId]);

  /**
   * 初始化加载
   */
  useEffect(() => {
    fetchIntents();
    fetchStats();
  }, [fetchIntents, fetchStats]);

  /**
   * 创建意图
   */
  const handleCreate = () => {
    setCurrentIntent(null);
    form.resetFields();
    setEditModalVisible(true);
  };

  /**
   * 编辑意图
   */
  const handleEdit = (intent: Intent) => {
    setCurrentIntent(intent);
    form.setFieldsValue({
      intent_name: intent.intent_name,
      description: intent.description,
      agent_id: intent.agent_id,
      workflow_id: intent.workflow_id,
      confidence: intent.confidence,
      is_active: intent.is_active,
      examples: intent.examples,
    });
    setEditModalVisible(true);
  };

  /**
   * 删除意图
   */
  const handleDelete = async (intentId: string) => {
    try {
      // TODO: 调用API删除意图
      // await routingService.deleteIntent(intentId);
      message.success(t('common.success'));
      fetchIntents();
      fetchStats();
    } catch (error) {
      console.error('Failed to delete intent:', error);
      message.error(t('common.failed'));
    }
  };

  /**
   * 保存意图
   */
  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      // TODO: 调用API保存意图
      // if (currentIntent) {
      //   await routingService.updateIntent(currentIntent.intent_id, values);
      // } else {
      //   await routingService.createIntent(tenantId, values);
      // }
      message.success(t('common.success'));
      setEditModalVisible(false);
      fetchIntents();
      fetchStats();
    } catch (error) {
      console.error('Failed to save intent:', error);
      message.error(t('common.failed'));
    }
  };

  /**
   * 测试意图识别
   */
  const handleTest = async () => {
    if (!testInput.trim()) {
      message.warning(t('routing.pleaseEnterTestText'));
      return;
    }

    try {
      setTestLoading(true);
      // TODO: 调用API测试识别
      // const response = await routingService.recognizeIntent(tenantId, testInput);
      // setTestResults(response.data);
    } catch (error) {
      console.error('Failed to test intent:', error);
      message.error(t('common.failed'));
    } finally {
      setTestLoading(false);
    }
  };

  /**
   * 训练意图模型
   */
  const handleTrain = async () => {
    if (trainSamples.length < 5) {
      message.warning(t('routing.atLeast5SamplesRequired'));
      return;
    }

    try {
      setTrainLoading(true);
      setTrainProgress(0);

      // TODO: 调用API训练模型
      // const response = await routingService.trainIntentModel(currentIntent.intent_id, {
      //   tenant_id: tenantId,
      //   samples: trainSamples,
      // });
      // 模拟进度
      const interval = setInterval(() => {
        setTrainProgress(prev => {
          if (prev >= 100) {
            clearInterval(interval);
            setTrainLoading(false);
            message.success(t('routing.trainSuccess'));
            return 100;
          }
          return prev + 10;
        });
      }, 200);
    } catch (error) {
      console.error('Failed to train intent:', error);
      message.error(t('common.failed'));
      setTrainLoading(false);
    }
  };

  /**
   * 表格列定义
   */
  const columns = useMemo(() => [
    {
      title: t('routing.intentName'),
      dataIndex: 'intent_name',
      key: 'intent_name',
      sorter: true,
      render: (text: string, record: Intent) => (
        <Space>
          <Text strong>{text}</Text>
          {!record.is_active && <Tag color="red">{t('common.inactive')}</Tag>}
        </Space>
      ),
    },
    {
      title: t('routing.description'),
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text: string) => <Text ellipsis={{ tooltip: text }}>{text || '-'}</Text>,
    },
    {
      title: t('routing.targetAgent'),
      dataIndex: 'agent_id',
      key: 'agent_id',
      render: (text: string) => <Text code>{text}</Text>,
    },
    {
      title: t('routing.confidence'),
      dataIndex: 'confidence',
      key: 'confidence',
      sorter: true,
      render: (value: number) => (
        <Progress
          percent={Math.round(value * 100)}
          size="small"
          status={value >= 0.8 ? 'success' : value >= 0.6 ? 'normal' : 'exception'}
        />
      ),
    },
    {
      title: t('routing.exampleCount'),
      dataIndex: 'examples',
      key: 'examples',
      render: (examples: string[]) => <Tag color="blue">{examples?.length || 0}</Tag>,
    },
    {
      title: t('common.updatedAt'),
      dataIndex: 'updated_at',
      key: 'updated_at',
      sorter: true,
      render: (value: number) => new Date(value * 1000).toLocaleString(),
    },
    {
      title: t('common.actions'),
      key: 'actions',
      fixed: 'right' as const,
      width: 200,
      render: (_: unknown, record: Intent) => (
        <Space>
          <Tooltip title={t('routing.test')}>
            <Button
              type="text"
              icon={<TestOutlined />}
              onClick={() => {
                setCurrentIntent(record);
                setTestInput('');
                setTestResults(null);
                setTestModalVisible(true);
              }}
            />
          </Tooltip>
          <Tooltip title={t('routing.train')}>
            <Button
              type="text"
              icon={<RocketOutlined />}
              onClick={() => {
                setCurrentIntent(record);
                setTrainSamples(record.examples || []);
                setTrainProgress(0);
                setTrainModalVisible(true);
              }}
            />
          </Tooltip>
          <Tooltip title={t('common.edit')}>
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          <Popconfirm
            title={t('common.deleteConfirm')}
            onConfirm={() => handleDelete(record.intent_id)}
            okText={t('common.confirm')}
            cancelText={t('common.cancel')}
          >
            <Button type="text" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ], [t]);

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#fff', padding: '0 24px', borderBottom: '1px solid #f0f0f0' }}>
        <Row justify="space-between" align="middle">
          <Col>
            <Title level={3} style={{ margin: 0 }}>
              <BulbOutlined /> {t('routing.intentManagement')}
            </Title>
          </Col>
          <Col>
            <Space>
              <Button
                icon={<TestOutlined />}
                onClick={() => {
                  setTestInput('');
                  setTestResults(null);
                  setTestModalVisible(true);
                }}
              >
                {t('routing.batchTest')}
              </Button>
              <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
                {t('routing.createIntent')}
              </Button>
            </Space>
          </Col>
        </Row>
      </Header>

      <Content style={{ padding: '24px', background: '#f0f2f5' }}>
        {/* 统计卡片 */}
        <Row gutter={16} style={{ marginBottom: 24 }}>
          <Col span={6}>
            <Card>
              <Statistic
                title={t('routing.totalIntents')}
                value={stats.totalIntents}
                prefix={<BulbOutlined />}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title={t('routing.activeIntents')}
                value={stats.activeIntents}
                valueStyle={{ color: '#3f8600' }}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title={t('routing.totalSamples')}
                value={stats.totalSamples}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title={t('routing.avgAccuracy')}
                value={stats.avgAccuracy}
                suffix="%"
                precision={2}
              />
            </Card>
          </Col>
        </Row>

        {/* 意图列表 */}
        <Card>
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            {/* 筛选栏 */}
            <Row justify="space-between" align="middle">
              <Col>
                <Space>
                  <Input.Search
                    placeholder={t('routing.searchIntents')}
                    allowClear
                    style={{ width: 300 }}
                    onSearch={setSearchKeyword}
                    onChange={e => setSearchKeyword(e.target.value)}
                  />
                  <Select
                    style={{ width: 150 }}
                    value={filterActive}
                    onChange={setFilterActive}
                    allowClear
                    placeholder={t('common.status')}
                  >
                    <Option value={true}>{t('common.active')}</Option>
                    <Option value={false}>{t('common.inactive')}</Option>
                  </Select>
                </Space>
              </Col>
            </Row>

            {/* 表格 */}
            <Table
              columns={columns}
              dataSource={intents}
              loading={loading}
              rowKey="intent_id"
              pagination={{
                current: page,
                pageSize,
                total,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: total => t('common.total', { count: total }),
                onChange: (page, pageSize) => {
                  setPage(page);
                  setPageSize(pageSize || 20);
                },
              }}
              scroll={{ x: 1200 }}
            />
          </Space>
        </Card>
      </Content>

      {/* 编辑/创建对话框 */}
      <Modal
        title={currentIntent ? t('routing.editIntent') : t('routing.createIntent')}
        open={editModalVisible}
        onOk={handleSave}
        onCancel={() => setEditModalVisible(false)}
        width={800}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            confidence: 0.8,
            is_active: true,
            examples: [],
          }}
        >
          <Form.Item
            name="intent_name"
            label={t('routing.intentName')}
            rules={[{ required: true, message: t('routing.intentNameRequired') }]}
          >
            <Input placeholder={t('routing.intentNamePlaceholder')} />
          </Form.Item>

          <Form.Item
            name="description"
            label={t('routing.description')}
          >
            <Input.TextArea
              rows={3}
              placeholder={t('routing.descriptionPlaceholder')}
            />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="agent_id"
                label={t('routing.targetAgent')}
                rules={[{ required: true, message: t('routing.targetAgentRequired') }]}
              >
                <Select placeholder={t('routing.selectAgent')}>
                  {/* TODO: 从API加载Agent列表 */}
                  <Option value="agent1">Agent 1</Option>
                  <Option value="agent2">Agent 2</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="confidence"
                label={t('routing.confidenceThreshold')}
                rules={[{ required: true }]}
              >
                <Input
                  type="number"
                  min={0}
                  max={1}
                  step={0.01}
                  addonAfter="%"
                />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="examples"
            label={t('routing.exampleSentences')}
            rules={[{ type: 'array' }]}
          >
            <Select
              mode="tags"
              placeholder={t('routing.enterExampleSentences')}
              tokenSeparators={[', '\n']}
            >
              {form.getFieldValue('examples')?.map((example: string) => (
                <Option key={example} value={example}>
                  {example}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="is_active"
            label={t('common.status')}
            valuePropName="checked"
          >
            <Switch checkedChildren={t('common.active')} unCheckedChildren={t('common.inactive')} />
          </Form.Item>
        </Form>
      </Modal>

      {/* 测试对话框 */}
      <Modal
        title={t('routing.testIntentRecognition')}
        open={testModalVisible}
        onCancel={() => setTestModalVisible(false)}
        footer={null}
        width={800}
      >
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <Input.TextArea
            rows={4}
            placeholder={t('routing.enterTestText')}
            value={testInput}
            onChange={e => setTestInput(e.target.value)}
          />
          <Button
            type="primary"
            icon={<TestOutlined />}
            loading={testLoading}
            onClick={handleTest}
            block
          >
            {t('routing.test')}
          </Button>

          {testResults && (
            <Card title={t('routing.testResult')} size="small">
              <Descriptions column={2} bordered size="small">
                <Descriptions.Item label={t('routing.intentName')}>
                  {testResults.intent_name}
                </Descriptions.Item>
                <Descriptions.Item label={t('routing.confidence')}>
                  <Progress
                    percent={Math.round(testResults.confidence * 100)}
                    size="small"
                  />
                </Descriptions.Item>
                <Descriptions.Item label={t('routing.matchMethod')}>
                  <Tag color="blue">{testResults.match_method}</Tag>
                </Descriptions.Item>
                <Descriptions.Item label={t('routing.targetAgent')}>
                  {testResults.agent_id}
                </Descriptions.Item>
              </Descriptions>

              {Object.keys(testResults.entities || {}).length > 0 && (
                <div style={{ marginTop: 16 }}>
                  <Text strong>{t('routing.extractedEntities')}:</Text>
                  <div style={{ marginTop: 8 }}>
                    {Object.entries(testResults.entities).map(([key, values]) => (
                      <Tag key={key} color="green">
                        {key}: {values.join(', ')}
                      </Tag>
                    ))}
                  </div>
                </div>
              )}
            </Card>
          )}
        </Space>
      </Modal>

      {/* 训练对话框 */}
      <Modal
        title={t('routing.trainIntentModel')}
        open={trainModalVisible}
        onCancel={() => setTrainModalVisible(false)}
        footer={[
          <Button key="cancel" onClick={() => setTrainModalVisible(false)}>
            {t('common.cancel')}
          </Button>,
          <Button
            key="submit"
            type="primary"
            loading={trainLoading}
            onClick={handleTrain}
            disabled={trainSamples.length < 5}
          >
            {t('routing.startTrain')}
          </Button>,
        ]}
        width={800}
      >
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <Alert
            message={t('routing.trainTip')}
            description={t('routing.trainTipDescription')}
            type="info"
            showIcon
          />

          <div>
            <Text strong>{t('routing.trainingSamples')} ({trainSamples.length})</Text>
            <Select
              mode="tags"
              style={{ width: '100%', marginTop: 8 }}
              placeholder={t('routing.enterTrainingSamples')}
              value={trainSamples}
              onChange={setTrainSamples}
              tokenSeparators={['\n']}
            >
              {trainSamples.map(sample => (
                <Option key={sample} value={sample}>
                  {sample}
                </Option>
              ))}
            </Select>
          </div>

          {trainLoading && (
            <div>
              <Text strong>{t('routing.trainingProgress')}</Text>
              <Progress percent={trainProgress} />
            </div>
          )}
        </Space>
      </Modal>
    </Layout>
  );
};

export default IntentManagementPage;
