// frontend/packages/studio/src/pages/permission/RoleList/RoleList.tsx

import React, { useState } from 'react';
import { Button, Table, Modal, Input } from '@coze-studio/ui-components';
import { RoleActions, Role } from './components/RoleActions';
import { useStyles } from './RoleList.styles';

/**
 * RoleList 角色列表页
 *
 * 展示所有角色的列表，支持搜索、分页、操作
 *
 * @example
 * ```tsx
 * <RoleList />
 * ```
 */
export const RoleList: React.FC = () => {
  const classes = useStyles();

  // 搜索关键词
  const [keyword, setKeyword] = useState('');

  // 分页
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);

  // 模拟数据
  const mockRoles: Role[] = [
    {
      role_id: '1',
      role_name: '超级管理员',
      role_code: 'super_admin',
      role_type: 'system',
      description: '拥有所有权限',
    },
    {
      role_id: '2',
      role_name: '管理员',
      role_code: 'admin',
      role_type: 'system',
      description: '拥有大部分权限',
    },
    {
      role_id: '3',
      role_name: '普通用户',
      role_code: 'user',
      role_type: 'system',
      description: '拥有基本权限',
    },
    {
      role_id: '4',
      role_name: '客服专员',
      role_code: 'customer_service',
      role_type: 'custom',
      description: '负责客户服务',
    },
  ];

  const [data, setData] = useState(mockRoles);
  const [loading, setLoading] = useState(false);

  // 过滤数据
  const filteredData = data.filter((role) =>
    role.role_name.toLowerCase().includes(keyword.toLowerCase()) ||
    role.role_code.toLowerCase().includes(keyword.toLowerCase())
  );

  // 处理搜索
  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setKeyword(e.target.value);
  };

  // 处理删除
  const handleDelete = (roleId: string, roleName: string) => {
    Modal.confirm({
      visible: true,
      title: '确认删除',
      onClose: () => {},
      children: `确定要删除角色"${roleName}"吗？此操作不可恢复。`,
      onOk: () => {
        // 这里应该调用删除API
        console.log('Delete role:', roleId);
        setData(data.filter((r) => r.role_id !== roleId));
      },
    });
  };

  // 处理复制
  const handleCopy = (role: Role) => {
    console.log('Copy role:', role);
    // 这里应该调用复制API
  };

  // 处理分页
  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    console.log('Page changed:', page);
  };

  // 表格列定义
  const columns = [
    {
      title: '角色名称',
      dataIndex: 'role_name' as keyof Role,
      key: 'role_name',
    },
    {
      title: '角色编码',
      dataIndex: 'role_code' as keyof Role,
      key: 'role_code',
    },
    {
      title: '角色类型',
      dataIndex: 'role_type' as keyof Role,
      key: 'role_type',
      render: (type: Role['role_type']) => {
        const typeMap = {
          system: '系统角色',
          custom: '自定义角色',
        };
        const colorMap = {
          system: '#722ED1',
          custom: '#1890FF',
        };
        return (
          <span style={{ color: colorMap[type] }}>
            {typeMap[type] || type}
          </span>
        );
      },
    },
    {
      title: '描述',
      dataIndex: 'description' as keyof Role,
      key: 'description',
      render: (text: string) => text || '-',
    },
    {
      title: '用户数',
      key: 'user_count',
      render: () => Math.floor(Math.random() * 100), // 模拟用户数
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Role) => (
        <RoleActions
          role={record}
          onViewDetail={() => console.log('View detail:', record.role_id)}
          onEdit={() => console.log('Edit role:', record.role_id)}
          onCopy={() => handleCopy(record)}
          onDelete={() => handleDelete(record.role_id, record.role_name)}
        />
      ),
    },
  ];

  return (
    <div className={classes.container}>
      {/* 页面头部 */}
      <div className={classes.header}>
        <h1 className={classes.title}>角色管理</h1>
        <Button variant="primary">创建角色</Button>
      </div>

      {/* 搜索栏 */}
      <div className={classes.searchBar}>
        <Input
          placeholder="搜索角色名称或编码"
          value={keyword}
          onChange={handleSearchChange}
          block
        />
      </div>

      {/* 表格 */}
      <Table
        columns={columns}
        dataSource={filteredData}
        loading={loading}
        pagination={{
          current: currentPage,
          pageSize,
          total: filteredData.length,
          onChange: handlePageChange,
        }}
        bordered
      />
    </div>
  );
};

export default RoleList;
