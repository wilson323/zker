// frontend/packages/studio/src/pages/permission/RoleDetail/RoleDetail.tsx

import React, { useState } from 'react';
import { Button, Input, PermissionTree, PermissionNode } from '@coze-studio/ui-components';
import { DataPermissionConfig } from './components/DataPermissionConfig';
import { useStyles } from './RoleDetail.styles';

/**
 * RoleDetail 角色详情页
 *
 * 显示角色的详细信息，包括基本信息、功能权限、数据权限
 *
 * @example
 * ```tsx
 * <RoleDetail roleId="123" />
 * ```
 */
export const RoleDetail: React.FC<{ roleId?: string }> = ({ roleId }) => {
  const classes = useStyles();
  const [activeTab, setActiveTab] = useState<'basic' | 'permission' | 'data'>('basic');

  // 模拟角色数据
  const mockRole = {
    role_id: roleId || '1',
    role_name: '管理员',
    role_code: 'admin',
    role_type: 'system' as const,
    description: '拥有大部分权限',
  };

  // 模拟权限树数据
  const permissionTreeData: PermissionNode[] = [
    {
      key: 'bot',
      title: '机器人管理',
      children: [
        { key: 'bot.view', title: '查看' },
        { key: 'bot.create', title: '创建' },
        { key: 'bot.edit', title: '编辑' },
        { key: 'bot.delete', title: '删除' },
      ],
    },
    {
      key: 'knowledge',
      title: '知识库管理',
      children: [
        { key: 'knowledge.view', title: '查看' },
        { key: 'knowledge.create', title: '创建' },
        { key: 'knowledge.edit', title: '编辑' },
        { key: 'knowledge.delete', title: '删除' },
      ],
    },
    {
      key: 'conversation',
      title: '对话管理',
      children: [
        { key: 'conversation.view', title: '查看' },
        { key: 'conversation.export', title: '导出' },
      ],
    },
  ];

  const [roleName, setRoleName] = useState(mockRole.role_name);
  const [roleCode, setRoleCode] = useState(mockRole.role_code);
  const [description, setDescription] = useState(mockRole.description || '');

  const handleBack = () => {
    window.history.back();
  };

  const handleSave = () => {
    console.log('Save role:', { roleName, roleCode, description });
    // 这里应该调用保存API
  };

  const tabs = [
    { key: 'basic' as const, label: '基本信息' },
    { key: 'permission' as const, label: '功能权限' },
    { key: 'data' as const, label: '数据权限' },
  ];

  return (
    <div className={classes.container}>
      {/* 页面头部 */}
      <div className={classes.header}>
        <div className={classes.headerLeft}>
          <button className={classes.backButton} onClick={handleBack}>
            ← 返回
          </button>
          <h1 className={classes.title}>{roleName}</h1>
        </div>
        <div className={classes.headerRight}>
          <Button variant="outline" onClick={handleBack}>
            取消
          </Button>
          <Button variant="primary" onClick={handleSave}>
            保存
          </Button>
        </div>
      </div>

      {/* 标签页导航 */}
      <div className={classes.tabs}>
        {tabs.map((tab) => (
          <button
            key={tab.key}
            className={`${classes.tab} ${activeTab === tab.key ? classes.activeTab : ''}`}
            onClick={() => setActiveTab(tab.key)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* 标签页内容 */}
      <div className={classes.tabContent}>
        {activeTab === 'basic' && (
          <div className={classes.basicInfo}>
            <h3 className={classes.sectionTitle}>基本信息</h3>

            <div className={classes.form}>
              <div className={classes.formRow}>
                <label className={classes.label}>角色ID</label>
                <input
                  type="text"
                  className={classes.readonlyInput}
                  value={mockRole.role_id}
                  disabled
                />
              </div>

              <div className={classes.formRow}>
                <label className={classes.label}>角色名称 *</label>
                <Input
                  placeholder="请输入角色名称"
                  value={roleName}
                  onChange={(e) => setRoleName(e.target.value)}
                  block
                />
              </div>

              <div className={classes.formRow}>
                <label className={classes.label}>角色编码 *</label>
                <Input
                  placeholder="请输入角色编码"
                  value={roleCode}
                  onChange={(e) => setRoleCode(e.target.value)}
                  block
                  disabled={mockRole.role_type === 'system'}
                />
              </div>

              <div className={classes.formRow}>
                <label className={classes.label}>描述</label>
                <textarea
                  className={classes.textarea}
                  placeholder="请输入角色描述"
                  rows={4}
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
              </div>
            </div>
          </div>
        )}

        {activeTab === 'permission' && (
          <div className={classes.permissions}>
            <h3 className={classes.sectionTitle}>功能权限</h3>
            <p className={classes.sectionDescription}>
              选择该角色可以访问的功能模块和操作权限
            </p>
            <PermissionTree
              data={permissionTreeData}
              defaultCheckedKeys={['bot.view', 'bot.create', 'knowledge.view']}
              checkable
              showCheckbox
              onChange={(keys, nodes) => {
                console.log('Permission changed:', keys, nodes);
              }}
            />
          </div>
        )}

        {activeTab === 'data' && (
          <div className={classes.dataPermissions}>
            <h3 className={classes.sectionTitle}>数据权限</h3>
            <p className={classes.sectionDescription}>
              配置该角色可以访问的数据范围
            </p>
            <DataPermissionConfig
              onChange={(configs) => {
                console.log('Data permission changed:', configs);
              }}
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default RoleDetail;
