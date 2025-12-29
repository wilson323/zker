// frontend/packages/studio/src/pages/routing/RoutingRules/RoutingRules.tsx

import React, { useState } from 'react';
import { Table, Button, Modal, message, Tag } from '@coze-studio/ui-components';
import { useStyles } from './RoutingRules.styles';

export interface RoutingRule {
  rule_id: string;
  rule_name: string;
  priority: number;
  intent_matcher_id: string;
  target_service: string;
  target_endpoint: string;
  condition: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export const RoutingRules: React.FC = () => {
  const classes = useStyles();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);
  const [keyword, setKeyword] = useState('');

  // Mock data - 应该从 @coze-studio/api-client 获取
  const [data, setData] = useState<RoutingRule[]>([
    {
      rule_id: '1',
      rule_name: '高价值客户路由',
      priority: 100,
      intent_matcher_id: 'intent-001',
      target_service: 'premium-service',
      target_endpoint: '/api/premium/handle',
      condition: 'customer_tier == "enterprise" && value > 10000',
      is_active: true,
      created_at: '2025-01-01T10:00:00Z',
      updated_at: '2025-01-01T10:00:00Z',
    },
    {
      rule_id: '2',
      rule_name: '普通客户路由',
      priority: 50,
      intent_matcher_id: 'intent-002',
      target_service: 'standard-service',
      target_endpoint: '/api/standard/handle',
      condition: 'customer_tier == "individual"',
      is_active: true,
      created_at: '2025-01-01T11:00:00Z',
      updated_at: '2025-01-01T11:00:00Z',
    },
    {
      rule_id: '3',
      rule_name: '测试路由规则',
      priority: 10,
      intent_matcher_id: 'intent-003',
      target_service: 'test-service',
      target_endpoint: '/api/test/handle',
      condition: 'environment == "test"',
      is_active: false,
      created_at: '2025-01-01T12:00:00Z',
      updated_at: '2025-01-01T12:00:00Z',
    },
  ]);

  const handleDelete = (ruleId: string, ruleName: string) => {
    Modal.confirm({
      visible: true,
      title: '确认删除',
      onClose: () => {},
      children: `确定要删除路由规则"${ruleName}"吗？`,
      onOk: () => {
        setData(data.filter(rule => rule.rule_id !== ruleId));
        message.success('删除成功');
      },
    });
  };

  const handleToggleActive = (ruleId: string) => {
    setData(data.map(rule =>
      rule.rule_id === ruleId
        ? { ...rule, is_active: !rule.is_active }
        : rule
    ));
    message.success('状态更新成功');
  };

  const filteredData = data.filter(rule =>
    rule.rule_name.toLowerCase().includes(keyword.toLowerCase()) ||
    rule.target_service.toLowerCase().includes(keyword.toLowerCase())
  );

  const columns = [
    {
      title: '规则名称',
      dataIndex: 'rule_name',
      key: 'rule_name',
      sorter: true,
      render: (name: string, record: RoutingRule) => (
        <div>
          <div>{name}</div>
          <div style={{ fontSize: '12px', color: '#8C8C8C' }}>
            {record.rule_id}
          </div>
        </div>
      ),
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      key: 'priority',
      sorter: true,
      render: (priority: number) => (
        <Tag color={priority >= 80 ? 'red' : priority >= 50 ? 'orange' : 'blue'}>
          {priority}
        </Tag>
      ),
    },
    {
      title: '目标服务',
      dataIndex: 'target_service',
      key: 'target_service',
      render: (service: string, record: RoutingRule) => (
        <div>
          <div>{service}</div>
          <div style={{ fontSize: '12px', color: '#8C8C8C' }}>
            {record.target_endpoint}
          </div>
        </div>
      ),
    },
    {
      title: '条件',
      dataIndex: 'condition',
      key: 'condition',
      ellipsis: true,
      render: (condition: string) => (
        <code
          style={{
            fontSize: '12px',
            padding: '2px 6px',
            backgroundColor: '#F5F5F5',
            borderRadius: '4px',
          }}
        >
          {condition}
        </code>
      ),
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'success' : 'default'}>
          {isActive ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: RoutingRule) => (
        <div style={{ display: 'flex', gap: '8px' }}>
          <Button
            size="sm"
            variant="outline"
            onClick={() => handleToggleActive(record.rule_id)}
          >
            {record.is_active ? '禁用' : '启用'}
          </Button>
          <Button size="sm" variant="outline">
            编辑
          </Button>
          <Button
            size="sm"
            variant="text"
            style={{ color: '#F5222D' }}
            onClick={() => handleDelete(record.rule_id, record.rule_name)}
          >
            删除
          </Button>
        </div>
      ),
    },
  ];

  return (
    <div className={classes.container}>
      <div className={classes.header}>
        <h1 className={classes.title}>路由规则管理</h1>
        <Button variant="primary">创建规则</Button>
      </div>

      <div className={classes.filter}>
        <input
          type="text"
          placeholder="搜索规则名称或目标服务"
          value={keyword}
          onChange={e => setKeyword(e.target.value)}
          className={classes.searchInput}
        />
        <Button variant="outline">重置</Button>
      </div>

      <Table
        className={classes.table}
        columns={columns}
        dataSource={filteredData}
        pagination={{
          current: currentPage,
          pageSize,
          total: filteredData.length,
          onChange: page => setCurrentPage(page),
        }}
      />
    </div>
  );
};
