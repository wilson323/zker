// frontend/packages/studio/src/pages/tenant/TenantList/TenantList.tsx

import React, { useState } from 'react';
import { Button, Table, Modal } from '@coze-studio/ui-components';
import { TenantFilter, TenantFilterValue } from './components/TenantFilter';
import { TenantActions, Tenant } from './components/TenantActions';
import { useStyles } from './TenantList.styles';

/**
 * TenantList 租户列表页
 *
 * 展示所有租户的列表，支持筛选、分页、操作
 *
 * @example
 * ```tsx
 * <TenantList />
 * ```
 */
export const TenantList: React.FC = () => {
  const classes = useStyles();

  // 筛选条件
  const [filter, setFilter] = useState<TenantFilterValue>({});

  // 分页
  const [pageToken, setPageToken] = useState('');
  const [pageSize] = useState(20);
  const [currentPage, setCurrentPage] = useState(1);

  // 模拟数据
  const mockTenants: Tenant[] = [
    {
      tenant_id: '1',
      tenant_name: '示例租户1',
      tenant_type: 'enterprise',
      status: 'active',
    },
    {
      tenant_id: '2',
      tenant_name: '示例租户2',
      tenant_type: 'team',
      status: 'active',
    },
  ];

  const [data, setData] = useState(mockTenants);
  const [loading, setLoading] = useState(false);

  // 处理筛选
  const handleFilterChange = (newFilter: TenantFilterValue) => {
    setFilter(newFilter);
    // 这里应该调用API重新获取数据
    console.log('Filter changed:', newFilter);
  };

  // 处理删除
  const handleDelete = (tenantId: string, tenantName: string) => {
    Modal.confirm({
      visible: true,
      title: '确认删除',
      onClose: () => {},
      children: `确定要删除租户"${tenantName}"吗？此操作不可恢复。`,
      onOk: () => {
        // 这里应该调用删除API
        console.log('Delete tenant:', tenantId);
        setData(data.filter((t) => t.tenant_id !== tenantId));
      },
    });
  };

  // 处理分页
  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    // 这里应该调用API获取对应页的数据
    console.log('Page changed:', page);
  };

  // 表格列定义
  const columns = [
    {
      title: '租户名称',
      dataIndex: 'tenant_name' as keyof Tenant,
      key: 'tenant_name',
    },
    {
      title: '租户类型',
      dataIndex: 'tenant_type' as keyof Tenant,
      key: 'tenant_type',
      render: (type: Tenant['tenant_type']) => {
        const typeMap = {
          individual: '个人',
          team: '团队',
          enterprise: '企业',
        };
        return typeMap[type] || type;
      },
    },
    {
      title: '状态',
      dataIndex: 'status' as keyof Tenant,
      key: 'status',
      render: (status: Tenant['status']) => {
        const statusMap = {
          active: '正常',
          suspended: '暂停',
          deleted: '已删除',
        };
        const colorMap = {
          active: '#52C41A',
          suspended: '#FAAD14',
          deleted: '#8C8C8C',
        };
        return (
          <span style={{ color: colorMap[status] }}>
            {statusMap[status] || status}
          </span>
        );
      },
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Tenant) => (
        <TenantActions
          tenant={record}
          onViewDetail={() => console.log('View detail:', record.tenant_id)}
          onEdit={() => console.log('Edit tenant:', record.tenant_id)}
          onDelete={() => handleDelete(record.tenant_id, record.tenant_name)}
          onSuspend={() => console.log('Suspend tenant:', record.tenant_id)}
        />
      ),
    },
  ];

  return (
    <div className={classes.container}>
      {/* 页面头部 */}
      <div className={classes.header}>
        <h1 className={classes.title}>租户管理</h1>
        <Button variant="primary">创建租户</Button>
      </div>

      {/* 筛选器 */}
      <TenantFilter filter={filter} onChange={handleFilterChange} />

      {/* 表格 */}
      <Table
        columns={columns}
        dataSource={data}
        loading={loading}
        pagination={{
          current: currentPage,
          pageSize,
          total: 100, // 模拟总数
          onChange: handlePageChange,
        }}
        bordered
      />
    </div>
  );
};

export default TenantList;
