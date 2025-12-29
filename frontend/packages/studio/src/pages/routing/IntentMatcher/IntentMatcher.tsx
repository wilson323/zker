// frontend/packages/studio/src/pages/routing/IntentMatcher/IntentMatcher.tsx

import React, { useState } from 'react';
import { Table, Button, Modal, message, Tag, Tabs } from '@coze-studio/ui-components';
import { useStyles } from './IntentMatcher.styles';

export interface IntentMatcher {
  matcher_id: string;
  matcher_name: string;
  matcher_type: 'keyword' | 'regex' | 'ml';
  config: {
    keywords?: string[];
    pattern?: string;
    model_id?: string;
    threshold?: number;
  };
  score_weight: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export const IntentMatcher: React.FC = () => {
  const classes = useStyles();
  const [activeTab, setActiveTab] = useState<'list' | 'test'>('list');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);
  const [keyword, setKeyword] = useState('');

  // Mock data - 应该从 @coze-studio/api-client 获取
  const [data, setData] = useState<IntentMatcher[]>([
    {
      matcher_id: 'intent-001',
      matcher_name: '企业客户意图识别',
      matcher_type: 'keyword',
      config: {
        keywords: ['企业版', '定制', '专属客服', '批量'],
      },
      score_weight: 0.8,
      is_active: true,
      created_at: '2025-01-01T10:00:00Z',
      updated_at: '2025-01-01T10:00:00Z',
    },
    {
      matcher_id: 'intent-002',
      matcher_name: '技术支持意图',
      matcher_type: 'regex',
      config: {
        pattern: '(bug|错误|问题|故障|无法使用)',
      },
      score_weight: 0.9,
      is_active: true,
      created_at: '2025-01-01T11:00:00Z',
      updated_at: '2025-01-01T11:00:00Z',
    },
    {
      matcher_id: 'intent-003',
      matcher_name: 'ML智能意图识别',
      matcher_type: 'ml',
      config: {
        model_id: 'nlu-intent-model-v1',
        threshold: 0.75,
      },
      score_weight: 1.0,
      is_active: true,
      created_at: '2025-01-01T12:00:00Z',
      updated_at: '2025-01-01T12:00:00Z',
    },
    {
      matcher_id: 'intent-004',
      matcher_name: '销售咨询意图',
      matcher_type: 'keyword',
      config: {
        keywords: ['价格', '购买', '付费', '试用', '升级'],
      },
      score_weight: 0.7,
      is_active: false,
      created_at: '2025-01-01T13:00:00Z',
      updated_at: '2025-01-01T13:00:00Z',
    },
  ]);

  const handleDelete = (matcherId: string, matcherName: string) => {
    Modal.confirm({
      visible: true,
      title: '确认删除',
      onClose: () => {},
      children: `确定要删除意图匹配器"${matcherName}"吗？`,
      onOk: () => {
        setData(data.filter(matcher => matcher.matcher_id !== matcherId));
        message.success('删除成功');
      },
    });
  };

  const filteredData = data.filter(matcher =>
    matcher.matcher_name.toLowerCase().includes(keyword.toLowerCase())
  );

  const columns = [
    {
      title: '匹配器名称',
      dataIndex: 'matcher_name',
      key: 'matcher_name',
      sorter: true,
      render: (name: string, record: IntentMatcher) => (
        <div>
          <div>{name}</div>
          <div style={{ fontSize: '12px', color: '#8C8C8C' }}>
            {record.matcher_id}
          </div>
        </div>
      ),
    },
    {
      title: '类型',
      dataIndex: 'matcher_type',
      key: 'matcher_type',
      render: (type: string) => {
        const typeMap = {
          keyword: { label: '关键词', color: 'blue' },
          regex: { label: '正则表达式', color: 'green' },
          ml: { label: '机器学习', color: 'purple' },
        };
        const config = typeMap[type as keyof typeof typeMap];
        return <Tag color={config.color}>{config.label}</Tag>;
      },
    },
    {
      title: '配置',
      key: 'config',
      render: (_: any, record: IntentMatcher) => (
        <div>
          {record.matcher_type === 'keyword' && (
            <div>
              {record.config.keywords?.map((kw, i) => (
                <Tag key={i} style={{ marginBottom: '4px' }}>
                  {kw}
                </Tag>
              ))}
            </div>
          )}
          {record.matcher_type === 'regex' && (
            <code
              style={{
                fontSize: '12px',
                padding: '2px 6px',
                backgroundColor: '#F5F5F5',
                borderRadius: '4px',
              }}
            >
              {record.config.pattern}
            </code>
          )}
          {record.matcher_type === 'ml' && (
            <div>
              <div>模型: {record.config.model_id}</div>
              <div style={{ fontSize: '12px', color: '#8C8C8C' }}>
                阈值: {record.config.threshold}
              </div>
            </div>
          )}
        </div>
      ),
    },
    {
      title: '权重',
      dataIndex: 'score_weight',
      key: 'score_weight',
      render: (weight: number) => (
        <Tag color={weight >= 0.8 ? 'red' : weight >= 0.5 ? 'orange' : 'blue'}>
          {weight}
        </Tag>
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
      render: (_: any, record: IntentMatcher) => (
        <div style={{ display: 'flex', gap: '8px' }}>
          <Button size="sm" variant="outline">
            编辑
          </Button>
          <Button
            size="sm"
            variant="text"
            style={{ color: '#F5222D' }}
            onClick={() => handleDelete(record.matcher_id, record.matcher_name)}
          >
            删除
          </Button>
        </div>
      ),
    },
  ];

  const tabs = [
    { key: 'list' as const, label: '匹配器列表' },
    { key: 'test' as const, label: '测试工具' },
  ];

  return (
    <div className={classes.container}>
      <div className={classes.header}>
        <h1 className={classes.title}>意图匹配器配置</h1>
        <Button variant="primary">创建匹配器</Button>
      </div>

      <div className={classes.tabs}>
        {tabs.map(tab => (
          <button
            key={tab.key}
            className={activeTab === tab.key ? classes.activeTab : classes.tab}
            onClick={() => setActiveTab(tab.key)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'list' && (
        <>
          <div className={classes.filter}>
            <input
              type="text"
              placeholder="搜索匹配器名称"
              value={keyword}
              onChange={e => setKeyword(e.target.value)}
              className={classes.searchInput}
            />
            <select className={classes.typeSelect}>
              <option value="">全部类型</option>
              <option value="keyword">关键词</option>
              <option value="regex">正则表达式</option>
              <option value="ml">机器学习</option>
            </select>
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
        </>
      )}

      {activeTab === 'test' && (
        <div className={classes.testTool}>
          <h3 className={classes.testTitle}>意图匹配测试</h3>
          <div className={classes.testInput}>
            <textarea
              placeholder="输入测试文本，测试意图匹配效果..."
              className={classes.testTextarea}
              rows={4}
            />
            <Button variant="primary">测试匹配</Button>
          </div>
          <div className={classes.testResult}>
            <h4>匹配结果</h4>
            <div className={classes.resultEmpty}>
              请输入测试文本后点击"测试匹配"按钮
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
