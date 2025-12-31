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
 * 组织管理页面
 *
 * 功能：
 * - 组织列表展示（支持分页、排序、筛选）
 * - 组织树形结构展示（Tree组件，支持展开/折叠）
 * - 创建组织对话框（表单验证）
 * - 编辑组织对话框
 * - 删除组织确认对话框（带子组织检查）
 * - 移动组织对话框（选择新父组织，带循环验证）
 * - 组织详情查看（祖先链、子组织列表、统计信息）
 */

import React, {
  useState,
  useMemo,
  useEffect,
  useCallback,
  useRef,
} from 'react';
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
  Modal,
  Form,
  message,
  Popconfirm,
  Tree,
  Descriptions,
  Row,
  Col,
  Statistic,
  Tooltip,
  Divider,
  Breadcrumb,
  Empty,
} from '@coze-studio/ui-components';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  MoveOutlined,
  EyeOutlined,
  ReloadOutlined,
  SearchOutlined,
  ApartmentOutlined,
  TeamOutlined,
  RobotOutlined,
  BranchesOutlined,
} from '@coze-studio/ui-icons';
import { useTranslation } from '@coze-studio/i18n';
import { organizationApi } from '@coze-studio/api-client';
import type {
  OrganizationDTO,
  OrganizationTreeNodeDTO,
  OrganizationDetailDTO,
  OrganizationStatsDTO,
  OrganizationType,
  OrganizationStatus,
  CreateOrganizationRequest,
  UpdateOrganizationRequest,
} from '@coze-studio/api-client';

const { Header, Content, Sider } = Layout;
const { Title, Text, Paragraph } = Typography;
const { Option } = Select;
const { Search } = Input;
const { DirectoryTree } = Tree;

// ==================== 类型定义 ====================

/**
 * Tree节点数据结构
 */
interface TreeNode {
  title: string;
  key: string;
  children?: TreeNode[];
  org: OrganizationDTO;
}

/**
 * 表格排序类型
 */
interface TableSorter {
  field?: string;
  order?: 'ascend' | 'descend';
}

// ==================== 主组件 ====================

/**
 * 组织管理页面组件
 */
export const OrganizationPage: React.FC = () => {
  const { tenantId = 'current' } = useParams<{ tenantId: string }>();
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [moveForm] = Form.useForm();

  // ==================== 状态管理 ====================

  // 组织列表
  const [organizations, setOrganizations] = useState<OrganizationDTO[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [tableSorter, setTableSorter] = useState<TableSorter>({});

  // 组织树
  const [orgTree, setOrgTree] = useState<OrganizationTreeNodeDTO[]>([]);
  const [treeLoading, setTreeLoading] = useState(false);
  const [selectedOrgId, setSelectedOrgId] = useState<string>();
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);
  const [autoExpandParent, setAutoExpandParent] = useState(true);

  // 筛选状态
  const [filterType, setFilterType] = useState<OrganizationType | undefined>();
  const [filterStatus, setFilterStatus] = useState<OrganizationStatus | undefined>();
  const [searchKeyword, setSearchKeyword] = useState('');

  // 对话框状态
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [moveModalVisible, setMoveModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [currentOrg, setCurrentOrg] = useState<OrganizationDTO | null>(null);
  const [orgDetail, setOrgDetail] = useState<OrganizationDetailDTO | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);

  // 统计数据
  const [stats, setStats] = useState<OrganizationStatsDTO | null>(null);
  const [statsLoading, setStatsLoading] = useState(false);

  // 可选的父组织列表（用于移动和创建）
  const [availableParents, setAvailableParents] = useState<OrganizationDTO[]>([]);

  // ==================== 数据获取 ====================

  /**
   * 获取组织列表
   */
  const fetchOrganizations = useCallback(async () => {
    try {
      setLoading(true);
      const response = await organizationApi.listOrganizations({
        tenant_id: tenantId,
        parent_id: selectedOrgId,
        org_type: filterType,
        status: filterStatus,
        keyword: searchKeyword || undefined,
        page,
        page_size: pageSize,
        sort_by: tableSorter.field,
        sort_order: tableSorter.order === 'ascend' ? 'asc' : 'desc',
      });
      setOrganizations(response.orgs);
      setTotal(response.total);
    } catch (error) {
      console.error('Failed to fetch organizations:', error);
      message.error(t('common.failed'));
    } finally {
      setLoading(false);
    }
  }, [tenantId, selectedOrgId, filterType, filterStatus, searchKeyword, page, pageSize, tableSorter, t]);

  /**
   * 获取组织树
   */
  const fetchOrgTree = useCallback(async () => {
    try {
      setTreeLoading(true);
      const tree = await organizationApi.getOrganizationTree({
        tenant_id: tenantId,
        max_depth: -1,
        include_inactive: true,
      });
      setOrgTree(tree);
    } catch (error) {
      console.error('Failed to fetch organization tree:', error);
      message.error(t('common.failed'));
    } finally {
      setTreeLoading(false);
    }
  }, [tenantId, t]);

  /**
   * 获取统计数据
   */
  const fetchStats = useCallback(async () => {
    try {
      setStatsLoading(true);
      const data = await organizationApi.getOrganizationStats(tenantId);
      setStats(data);
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    } finally {
      setStatsLoading(false);
    }
  }, [tenantId]);

  /**
   * 获取组织详情
   */
  const fetchOrgDetail = useCallback(async (orgId: string) => {
    try {
      setDetailLoading(true);
      const detail = await organizationApi.getOrganizationDetail(orgId);
      setOrgDetail(detail);
    } catch (error) {
      console.error('Failed to fetch organization detail:', error);
      message.error(t('common.failed'));
    } finally {
      setDetailLoading(false);
    }
  }, [t]);

  /**
   * 获取可用的父组织列表
   */
  const fetchAvailableParents = useCallback(async (excludeOrgId?: string) => {
    try {
      const allOrgs = await organizationApi.listOrganizations({
        tenant_id: tenantId,
        page: 1,
        page_size: 1000,
        status: 'active',
      });
      // 排除当前组织和其所有后代组织
      const excludeIds = new Set([excludeOrgId]);
      if (excludeOrgId) {
        const descendants = await organizationApi.getOrganizationDescendants(excludeOrgId);
        const collectIds = (nodes: OrganizationTreeNodeDTO[]) => {
          nodes.forEach(node => {
            excludeIds.add(node.org_id);
            if (node.children) {
              collectIds(node.children);
            }
          });
        };
        collectIds(descendants);
      }
      setAvailableParents(
        allOrgs.orgs.filter(org => !excludeIds.has(org.org_id))
      );
    } catch (error) {
      console.error('Failed to fetch available parents:', error);
    }
  }, [tenantId]);

  // ==================== 初始化加载 ====================

  useEffect(() => {
    fetchOrganizations();
    fetchOrgTree();
    fetchStats();
  }, [fetchOrgTree, fetchStats]); // 注意：fetchOrganizations依赖较多，单独处理

  useEffect(() => {
    fetchOrganizations();
  }, [fetchOrganizations]);

  // ==================== 事件处理 ====================

  /**
   * 树节点选择
   */
  const handleTreeSelect = useCallback(
    async (selectedKeys: React.Key[], info: any) => {
      if (selectedKeys.length > 0) {
        const orgId = selectedKeys[0] as string;
        setSelectedOrgId(orgId);
        setPage(1);
      } else {
        setSelectedOrgId(undefined);
        setPage(1);
      }
    },
    []
  );

  /**
   * 树节点展开
   */
  const handleTreeExpand = useCallback((expandedKeys: React.Key[]) => {
    setExpandedKeys(expandedKeys);
    setAutoExpandParent(false);
  }, []);

  /**
   * 创建组织
   */
  const handleCreate = useCallback(async () => {
    try {
      const values = await form.validateFields();
      const request: CreateOrganizationRequest = {
        tenant_id: tenantId,
        ...values,
      };
      await organizationApi.createOrganization(request);
      message.success(t('common.success'));
      setEditModalVisible(false);
      form.resetFields();
      fetchOrgTree();
      fetchOrganizations();
      fetchStats();
    } catch (error) {
      console.error('Failed to create organization:', error);
      if (error !== 'validation failed') {
        message.error(t('common.failed'));
      }
    }
  }, [form, tenantId, t, fetchOrgTree, fetchOrganizations, fetchStats]);

  /**
   * 编辑组织
   */
  const handleEdit = useCallback(
    async (org: OrganizationDTO) => {
      setCurrentOrg(org);
      form.setFieldsValue({
        org_name: org.org_name,
        org_type: org.org_type,
        org_code: org.org_code,
        parent_id: org.parent_id,
        sort_order: org.sort_order,
        description: org.description,
        leader_id: org.leader_id,
      });
      setEditModalVisible(true);
    },
    [form]
  );

  /**
   * 更新组织
   */
  const handleUpdate = useCallback(async () => {
    if (!currentOrg) return;

    try {
      const values = await form.validateFields();
      const request: UpdateOrganizationRequest = {
        ...values,
      };
      await organizationApi.updateOrganization(currentOrg.org_id, request);
      message.success(t('common.success'));
      setEditModalVisible(false);
      setCurrentOrg(null);
      form.resetFields();
      fetchOrgTree();
      fetchOrganizations();
    } catch (error) {
      console.error('Failed to update organization:', error);
      if (error !== 'validation failed') {
        message.error(t('common.failed'));
      }
    }
  }, [currentOrg, form, t, fetchOrgTree, fetchOrganizations]);

  /**
   * 删除组织
   */
  const handleDelete = useCallback(
    async (org: OrganizationDTO) => {
      try {
        // 检查是否有子组织
        if (org.child_count && org.child_count > 0) {
          Modal.confirm({
            title: t('org.deleteWarning'),
            content: t('org.deleteWarningDesc'),
            okText: t('common.confirm'),
            cancelText: t('common.cancel'),
            okButtonProps: { danger: true },
            onOk: async () => {
              await organizationApi.deleteOrganization(org.org_id, true);
              message.success(t('common.success'));
              fetchOrgTree();
              fetchOrganizations();
              fetchStats();
            },
          });
        } else {
          await organizationApi.deleteOrganization(org.org_id);
          message.success(t('common.success'));
          fetchOrgTree();
          fetchOrganizations();
          fetchStats();
        }
      } catch (error) {
        console.error('Failed to delete organization:', error);
        message.error(t('common.failed'));
      }
    },
    [t, fetchOrgTree, fetchOrganizations, fetchStats]
  );

  /**
   * 打开移动对话框
   */
  const handleMoveModalOpen = useCallback(
    async (org: OrganizationDTO) => {
      setCurrentOrg(org);
      moveForm.setFieldsValue({
        new_parent_id: org.parent_id || undefined,
        new_sort_order: org.sort_order,
      });
      await fetchAvailableParents(org.org_id);
      setMoveModalVisible(true);
    },
    [moveForm, fetchAvailableParents]
  );

  /**
   * 移动组织
   */
  const handleMove = useCallback(async () => {
    if (!currentOrg) return;

    try {
      const values = await moveForm.validateFields();
      await organizationApi.moveOrganization(currentOrg.org_id, {
        new_parent_id: values.new_parent_id,
        new_sort_order: values.new_sort_order,
      });
      message.success(t('common.success'));
      setMoveModalVisible(false);
      setCurrentOrg(null);
      moveForm.resetFields();
      fetchOrgTree();
      fetchOrganizations();
    } catch (error) {
      console.error('Failed to move organization:', error);
      if (error !== 'validation failed') {
        message.error(t('common.failed'));
      }
    }
  }, [currentOrg, moveForm, t, fetchOrgTree, fetchOrganizations]);

  /**
   * 查看详情
   */
  const handleViewDetail = useCallback(
    async (org: OrganizationDTO) => {
      setCurrentOrg(org);
      setDetailModalVisible(true);
      await fetchOrgDetail(org.org_id);
    },
    [fetchOrgDetail]
  );

  // ==================== 辅助函数 ====================

  /**
   * 转换树数据
   */
  const convertTreeData = useCallback((nodes: OrganizationTreeNodeDTO[]): TreeNode[] => {
    return nodes.map(node => ({
      title: (
        <Space>
          <Text>{node.org_name}</Text>
          <Tag color={getStatusColor(node.status)}>{t(`org.status.${node.status}`)}</Tag>
        </Space>
      ),
      key: node.org_id,
      children: node.children ? convertTreeData(node.children) : undefined,
      org: node,
    }));
  }, [t]);

  /**
   * 获取状态颜色
   */
  const getStatusColor = (status: OrganizationStatus) => {
    switch (status) {
      case 'active':
        return 'success';
      case 'inactive':
        return 'default';
      case 'frozen':
        return 'error';
      default:
        return 'default';
    }
  };

  /**
   * 获取类型标签颜色
   */
  const getTypeColor = (type: OrganizationType) => {
    switch (type) {
      case 'company':
        return 'blue';
      case 'division':
        return 'cyan';
      case 'department':
        return 'green';
      case 'project':
        return 'orange';
      default:
        return 'default';
    }
  };

  /**
   * 获取类型图标
   */
  const getTypeIcon = (type: OrganizationType) => {
    switch (type) {
      case 'company':
        return <ApartmentOutlined />;
      case 'division':
        return <BranchesOutlined />;
      case 'department':
        return <TeamOutlined />;
      case 'project':
        return <RobotOutlined />;
      default:
        return null;
    }
  };

  // ==================== 计算属性 ====================

  /**
   * 表格列定义
   */
  const columns = useMemo(() => [
    {
      title: t('org.orgName'),
      dataIndex: 'org_name',
      key: 'org_name',
      sorter: true,
      render: (text: string, record: OrganizationDTO) => (
        <Space>
          {getTypeIcon(record.org_type)}
          <Text strong>{text}</Text>
          {record.status !== 'active' && (
            <Tag color={getStatusColor(record.status)}>
              {t(`org.status.${record.status}`)}
            </Tag>
          )}
        </Space>
      ),
    },
    {
      title: t('org.orgCode'),
      dataIndex: 'org_code',
      key: 'org_code',
      render: (text: string) => <Text code>{text}</Text>,
    },
    {
      title: t('org.orgType'),
      dataIndex: 'org_type',
      key: 'org_type',
      sorter: true,
      filters: [
        { text: t('org.type.company'), value: 'company' },
        { text: t('org.type.division'), value: 'division' },
        { text: t('org.type.department'), value: 'department' },
        { text: t('org.type.project'), value: 'project' },
      ],
      render: (type: OrganizationType) => (
        <Tag color={getTypeColor(type)} icon={getTypeIcon(type)}>
          {t(`org.type.${type}`)}
        </Tag>
      ),
    },
    {
      title: t('org.level'),
      dataIndex: 'level',
      key: 'level',
      sorter: true,
      width: 80,
      render: (level: number) => <Tag>Level {level}</Tag>,
    },
    {
      title: t('org.parent'),
      dataIndex: 'parent_id',
      key: 'parent_id',
      render: (parentId: string | undefined) =>
        parentId ? (
          <Text ellipsis={{ tooltip: parentId }}>{parentId}</Text>
        ) : (
          <Tag color="blue">{t('org.root')}</Tag>
        ),
    },
    {
      title: t('org.sortOrder'),
      dataIndex: 'sort_order',
      key: 'sort_order',
      sorter: true,
      width: 100,
    },
    {
      title: t('org.description'),
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text: string) => (
        <Text ellipsis={{ tooltip: text }}>{text || '-'}</Text>
      ),
    },
    {
      title: t('common.createdAt'),
      dataIndex: 'created_at',
      key: 'created_at',
      sorter: true,
      width: 180,
      render: (value: number) => new Date(value * 1000).toLocaleString(),
    },
    {
      title: t('common.actions'),
      key: 'actions',
      fixed: 'right' as const,
      width: 200,
      render: (_: unknown, record: OrganizationDTO) => (
        <Space size="small">
          <Tooltip title={t('common.detail')}>
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={() => handleViewDetail(record)}
            />
          </Tooltip>
          <Tooltip title={t('common.edit')}>
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          <Tooltip title={t('org.move')}>
            <Button
              type="text"
              icon={<MoveOutlined />}
              onClick={() => handleMoveModalOpen(record)}
            />
          </Tooltip>
          <Popconfirm
            title={t('common.deleteConfirm')}
            onConfirm={() => handleDelete(record)}
            okText={t('common.confirm')}
            cancelText={t('common.cancel')}
            okButtonProps={{ danger: true }}
          >
            <Button type="text" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ], [t, getTypeIcon, getStatusColor, getTypeColor, handleViewDetail, handleEdit, handleMoveModalOpen, handleDelete]);

  /**
   * 树数据
   */
  const treeData = useMemo(() => convertTreeData(orgTree), [orgTree, convertTreeData]);

  // ==================== 渲染 ====================

  return (
    <Layout style={{ minHeight: '100vh', background: '#f0f2f5' }}>
      {/* 头部 */}
      <Header
        style={{
          background: '#fff',
          padding: '16px 24px',
          borderBottom: '1px solid #f0f0f0',
        }}
      >
        <Row justify="space-between" align="middle">
          <Col>
            <Title level={3} style={{ margin: 0 }}>
              <ApartmentOutlined /> {t('org.title')}
            </Title>
          </Col>
          <Col>
            <Space>
              <Button icon={<ReloadOutlined />} onClick={fetchOrganizations}>
                {t('common.refresh')}
              </Button>
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => {
                  setCurrentOrg(null);
                  form.resetFields();
                  setEditModalVisible(true);
                }}
              >
                {t('org.create')}
              </Button>
            </Space>
          </Col>
        </Row>
      </Header>

      <Content style={{ padding: '24px' }}>
        <Row gutter={24}>
          {/* 左侧：统计卡片 */}
          <Col span={24}>
            <Spin spinning={statsLoading}>
              <Row gutter={16} style={{ marginBottom: 24 }}>
                <Col span={6}>
                  <Card>
                    <Statistic
                      title={t('org.totalOrgs')}
                      value={stats?.total_orgs || 0}
                      prefix={<ApartmentOutlined />}
                    />
                  </Card>
                </Col>
                <Col span={6}>
                  <Card>
                    <Statistic
                      title={t('org.activeOrgs')}
                      value={stats?.active_orgs || 0}
                      valueStyle={{ color: '#3f8600' }}
                    />
                  </Card>
                </Col>
                <Col span={6}>
                  <Card>
                    <Statistic
                      title={t('org.maxLevel')}
                      value={stats?.max_level || 0}
                      suffix={t('org.levels')}
                    />
                  </Card>
                </Col>
                <Col span={6}>
                  <Card>
                    <Statistic
                      title={t('org.rootOrgs')}
                      value={stats?.root_orgs || 0}
                    />
                  </Card>
                </Col>
              </Row>
            </Spin>
          </Col>

          {/* 左侧：组织树 */}
          <Col span={6}>
            <Card
              title={t('org.orgTree')}
              extra={<Tag color="blue">{orgTree.length} {t('org.rootOrgs')}</Tag>}
              bodyStyle={{ padding: '12px', maxHeight: 'calc(100vh - 350px)', overflow: 'auto' }}
            >
              <Spin spinning={treeLoading}>
                {orgTree.length > 0 ? (
                  <DirectoryTree
                    multiple={false}
                    defaultExpandAll={false}
                    expandedKeys={expandedKeys}
                    autoExpandParent={autoExpandParent}
                    onExpand={handleTreeExpand}
                    onSelect={handleTreeSelect}
                    treeData={treeData}
                  />
                ) : (
                  <Empty description={t('org.noOrgs')} />
                )}
              </Spin>
            </Card>
          </Col>

          {/* 右侧：组织列表 */}
          <Col span={18}>
            <Card>
              <Space direction="vertical" size="large" style={{ width: '100%' }}>
                {/* 筛选栏 */}
                <Row justify="space-between" align="middle" gutter={16}>
                  <Col flex="auto">
                    <Space wrap>
                      <Search
                        placeholder={t('org.searchPlaceholder')}
                        allowClear
                        style={{ width: 300 }}
                        onSearch={setSearchKeyword}
                        onChange={e => setSearchKeyword(e.target.value)}
                        enterButton={<SearchOutlined />}
                      />
                      <Select
                        style={{ width: 150 }}
                        value={filterType}
                        onChange={setFilterType}
                        allowClear
                        placeholder={t('org.selectType')}
                      >
                        <Option value="company">{t('org.type.company')}</Option>
                        <Option value="division">{t('org.type.division')}</Option>
                        <Option value="department">{t('org.type.department')}</Option>
                        <Option value="project">{t('org.type.project')}</Option>
                      </Select>
                      <Select
                        style={{ width: 150 }}
                        value={filterStatus}
                        onChange={setFilterStatus}
                        allowClear
                        placeholder={t('common.status')}
                      >
                        <Option value="active">{t('org.status.active')}</Option>
                        <Option value="inactive">{t('org.status.inactive')}</Option>
                        <Option value="frozen">{t('org.status.frozen')}</Option>
                      </Select>
                      {selectedOrgId && (
                        <Button
                          onClick={() => {
                            setSelectedOrgId(undefined);
                            setPage(1);
                          }}
                        >
                          {t('org.clearFilter')}
                        </Button>
                      )}
                    </Space>
                  </Col>
                </Row>

                {/* 当前筛选提示 */}
                {selectedOrgId && (
                  <Alert
                    message={
                      <Space>
                        <Text>
                          {t('org.currentFilter')}:{' '}
                          {organizations.find(o => o.org_id === selectedOrgId)?.org_name}
                        </Text>
                        <Button
                          type="link"
                          size="small"
                          onClick={() => {
                            setSelectedOrgId(undefined);
                            setPage(1);
                          }}
                        >
                          {t('common.clear')}
                        </Button>
                      </Space>
                    }
                    type="info"
                    showIcon
                    closable
                    onClose={() => setSelectedOrgId(undefined)}
                  />
                )}

                {/* 表格 */}
                <Table
                  columns={columns}
                  dataSource={organizations}
                  loading={loading}
                  rowKey="org_id"
                  scroll={{ x: 1400 }}
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
                  onChange={(pagination, filters, sorter) => {
                    if (sorter && typeof sorter === 'object' && 'field' in sorter) {
                      setTableSorter({
                        field: sorter.field as string,
                        order: sorter.order as 'ascend' | 'descend' | undefined,
                      });
                    }
                  }}
                />
              </Space>
            </Card>
          </Col>
        </Row>
      </Content>

      {/* 创建/编辑对话框 */}
      <Modal
        title={
          currentOrg
            ? t('org.editOrg')
            : t('org.createOrg')
        }
        open={editModalVisible}
        onOk={currentOrg ? handleUpdate : handleCreate}
        onCancel={() => {
          setEditModalVisible(false);
          setCurrentOrg(null);
          form.resetFields();
        }}
        width={700}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            sort_order: 0,
          }}
        >
          <Form.Item
            name="org_name"
            label={t('org.orgName')}
            rules={[
              { required: true, message: t('org.orgNameRequired') },
              { max: 100, message: t('org.orgNameTooLong') },
            ]}
          >
            <Input placeholder={t('org.orgNamePlaceholder')} />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="org_type"
                label={t('org.orgType')}
                rules={[{ required: true, message: t('org.orgTypeRequired') }]}
              >
                <Select placeholder={t('org.selectType')}>
                  <Option value="company">{t('org.type.company')}</Option>
                  <Option value="division">{t('org.type.division')}</Option>
                  <Option value="department">{t('org.type.department')}</Option>
                  <Option value="project">{t('org.type.project')}</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="org_code"
                label={t('org.orgCode')}
                rules={[
                  { required: true, message: t('org.orgCodeRequired') },
                  { pattern: /^[A-Z0-9_-]+$/, message: t('org.orgCodeInvalid') },
                ]}
              >
                <Input placeholder={t('org.orgCodePlaceholder')} />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="parent_id"
                label={t('org.parentOrg')}
              >
                <Select
                  placeholder={t('org.selectParent')}
                  allowClear
                  showSearch
                  optionFilterProp="children"
                >
                  {availableParents.map(org => (
                    <Option key={org.org_id} value={org.org_id}>
                      {`${'  '.repeat(org.level)}${org.org_name} (${org.org_code})`}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="sort_order"
                label={t('org.sortOrder')}
                rules={[{ required: true, type: 'number', message: t('org.sortOrderRequired') }]}
              >
                <Input type="number" placeholder={t('org.sortOrderPlaceholder')} />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="description"
            label={t('org.description')}
          >
            <Input.TextArea
              rows={3}
              placeholder={t('org.descriptionPlaceholder')}
            />
          </Form.Item>

          <Form.Item
            name="leader_id"
            label={t('org.leader')}
          >
            <Input placeholder={t('org.leaderPlaceholder')} />
          </Form.Item>
        </Form>
      </Modal>

      {/* 移动对话框 */}
      <Modal
        title={t('org.moveOrg')}
        open={moveModalVisible}
        onOk={handleMove}
        onCancel={() => {
          setMoveModalVisible(false);
          setCurrentOrg(null);
          moveForm.resetFields();
        }}
        width={600}
        destroyOnClose
      >
        <Alert
          message={t('org.moveTip')}
          description={t('org.moveTipDesc')}
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
        />

        {currentOrg && (
          <Descriptions column={1} bordered size="small" style={{ marginBottom: 16 }}>
            <Descriptions.Item label={t('org.currentOrg')}>
              {currentOrg.org_name} ({currentOrg.org_code})
            </Descriptions.Item>
            <Descriptions.Item label={t('org.currentParent')}>
              {currentOrg.parent_id || t('org.root')}
            </Descriptions.Item>
          </Descriptions>
        )}

        <Form form={moveForm} layout="vertical">
          <Form.Item
            name="new_parent_id"
            label={t('org.newParentOrg')}
          >
            <Select
              placeholder={t('org.selectParent')}
              allowClear
              showSearch
              optionFilterProp="children"
            >
              <Option value={undefined}>{t('org.rootLevel')}</Option>
              {availableParents.map(org => (
                <Option key={org.org_id} value={org.org_id}>
                  {`${'  '.repeat(org.level)}${org.org_name} (${org.org_code})`}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="new_sort_order"
            label={t('org.newSortOrder')}
            rules={[{ required: true, type: 'number' }]}
          >
            <Input type="number" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 详情对话框 */}
      <Modal
        title={t('org.orgDetail')}
        open={detailModalVisible}
        onCancel={() => {
          setDetailModalVisible(false);
          setCurrentOrg(null);
          setOrgDetail(null);
        }}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            {t('common.close')}
          </Button>,
        ]}
        width={900}
        destroyOnClose
      >
        <Spin spinning={detailLoading}>
          {currentOrg && (
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              {/* 基本信息 */}
              <Card title={t('org.basicInfo')} size="small">
                <Descriptions column={2} bordered size="small">
                  <Descriptions.Item label={t('org.orgName')} span={2}>
                    <Space>
                      {getTypeIcon(currentOrg.org_type)}
                      {currentOrg.org_name}
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('org.orgCode')}>
                    {currentOrg.org_code}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('org.orgType')}>
                    <Tag color={getTypeColor(currentOrg.org_type)}>
                      {t(`org.type.${currentOrg.org_type}`)}
                    </Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('org.status')}>
                    <Tag color={getStatusColor(currentOrg.status)}>
                      {t(`org.status.${currentOrg.status}`)}
                    </Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('org.level')}>
                    {currentOrg.level}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('org.sortOrder')} span={2}>
                    {currentOrg.sort_order}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('org.description')} span={2}>
                    {currentOrg.description || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('common.createdAt')}>
                    {new Date(currentOrg.created_at * 1000).toLocaleString()}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('common.updatedAt')}>
                    {new Date(currentOrg.updated_at * 1000).toLocaleString()}
                  </Descriptions.Item>
                </Descriptions>
              </Card>

              {/* 组织路径 */}
              {orgDetail?.ancestor_chain && orgDetail.ancestor_chain.length > 0 && (
                <Card title={t('org.orgPath')} size="small">
                  <Breadcrumb>
                    {orgDetail.ancestor_chain.map((ancestor, index) => (
                      <Breadcrumb.Item key={ancestor.org_id}>
                        <Space>
                          {getTypeIcon(ancestor.org_type)}
                          {ancestor.org_name}
                        </Space>
                      </Breadcrumb.Item>
                    ))}
                    <Breadcrumb.Item>
                      <Space>
                        {getTypeIcon(currentOrg.org_type)}
                        <Text strong>{currentOrg.org_name}</Text>
                      </Space>
                    </Breadcrumb.Item>
                  </Breadcrumb>
                </Card>
              )}

              {/* 统计信息 */}
              {orgDetail && (
                <Card title={t('org.statistics')} size="small">
                  <Row gutter={16}>
                    <Col span={6}>
                      <Statistic
                        title={t('org.childCount')}
                        value={orgDetail.child_count}
                        prefix={<BranchesOutlined />}
                      />
                    </Col>
                    <Col span={6}>
                      <Statistic
                        title={t('org.totalDescendants')}
                        value={orgDetail.total_descendant_count}
                      />
                    </Col>
                    <Col span={6}>
                      <Statistic
                        title={t('org.userCount')}
                        value={orgDetail.user_count}
                        prefix={<TeamOutlined />}
                      />
                    </Col>
                    <Col span={6}>
                      <Statistic
                        title={t('org.botCount')}
                        value={orgDetail.bot_count}
                        prefix={<RobotOutlined />}
                      />
                    </Col>
                  </Row>
                </Card>
              )}
            </Space>
          )}
        </Spin>
      </Modal>
    </Layout>
  );
};

export default OrganizationPage;
