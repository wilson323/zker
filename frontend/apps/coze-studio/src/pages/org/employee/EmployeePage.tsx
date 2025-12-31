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
 * 员工管理页面
 *
 * 功能：
 * - 员工列表展示（表格，支持复杂筛选）
 * - 创建/编辑/删除员工
 * - 员工详情查看（完整信息卡片）
 * - 员工搜索（姓名、工号、手机号、邮箱）
 * - 员工状态管理（在职、试用、停职、离职）
 * - 批量操作（批量导入、批量导出）
 * - 邮箱/手机号格式验证
 * - 身份证号格式验证
 * - 状态转换按钮（试用→转正、在职→离职）
 * - 拼音搜索支持
 * - 员工头像显示
 */

import React, { useState, useMemo, useCallback, useEffect } from 'react';
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
  DatePicker,
  Modal,
  Form,
  message,
  Popconfirm,
  Tabs,
  Descriptions,
  Row,
  Col,
  Upload,
  Avatar,
  Radio,
  Divider,
  Tooltip,
  Dropdown,
  Menu,
  Badge,
} from '@coze-studio/ui-components';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
  DownloadOutlined,
  UploadOutlined,
  UserOutlined,
  ExportOutlined,
  CheckCircleOutlined,
  StopOutlined,
  RollbackOutlined,
  FileTextOutlined,
  PhoneOutlined,
  MailOutlined,
  IdcardOutlined,
  TeamOutlined,
  ApartmentOutlined,
  SafetyOutlined,
  EnvironmentOutlined,
  ContactsOutlined,
  CalendarOutlined,
  StarOutlined,
} from '@coze-studio/ui-icons';
import { useTranslation } from '@coze-studio/i18n';
import type { UploadProps } from '@coze-studio/ui-components';
import type { MenuProps } from '@coze-studio/ui-components';
import dayjs from 'dayjs';

const { Header, Content } = Layout;
const { Title, Text, Paragraph } = Typography;
const { Option } = Select;
const { TextArea } = Input;
const { TabPane } = Tabs;

/**
 * 员工状态枚举
 */
enum EmployeeStatus {
  PROBATION = 'probation',     // 试用
  ACTIVE = 'active',           // 在职
  SUSPENDED = 'suspended',     // 停职
  RESIGNED = 'resigned',       // 离职
}

/**
 * 合同类型枚举
 */
enum ContractType {
  LABOR = 'labor',             // 劳动合同
  INTERNSHIP = 'internship',   // 实习协议
  SERVICE = 'service',         // 劳务协议
  OUTSOURCING = 'outsourcing', // 外包协议
}

/**
 * 性别枚举
 */
enum Gender {
  MALE = 'male',
  FEMALE = 'female',
  OTHER = 'other',
}

/**
 * 员工信息接口
 */
interface Employee {
  employee_id: string;
  employee_number: string;          // 工号
  name: string;                     // 姓名
  name_pinyin?: string;             // 拼音
  gender: Gender;
  phone: string;                    // 手机号
  email: string;                    // 邮箱
  avatar_url?: string;              // 头像

  // 职位信息
  department_id: string;            // 所属部门ID
  department_name: string;          // 所属部门名称
  position: string;                 // 岗位
  level: string;                    // 职级
  hire_date: string;                // 就职日期

  // 合同信息
  contract_type: ContractType;      // 合同类型
  contract_number: string;          // 合同编号
  contract_start_date: string;      // 合同开始日期
  contract_end_date?: string;       // 合同结束日期

  // 试用期
  is_probation: boolean;            // 是否试用期
  probation_end_date?: string;      // 试用期结束日期

  // 状态
  status: EmployeeStatus;           // 员工状态

  // 其他信息
  id_card_number?: string;          // 身份证号
  address?: string;                 // 地址
  emergency_contact?: string;       // 紧急联系人
  emergency_phone?: string;         // 紧急联系电话

  // 系统字段
  created_at: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
}

/**
 * 部门信息接口
 */
interface Department {
  department_id: string;
  department_name: string;
  parent_id?: string;
}

/**
 * 表单数据接口
 */
interface EmployeeFormData {
  // 基本信息
  name: string;
  gender: Gender;
  phone: string;
  email: string;
  avatar_url?: string;

  // 职位信息
  department_id: string;
  position: string;
  level: string;
  hire_date: string;

  // 合同信息
  contract_type: ContractType;
  contract_number: string;
  contract_start_date: string;
  contract_end_date?: string;

  // 试用期
  is_probation: boolean;
  probation_end_date?: string;

  // 其他信息
  id_card_number?: string;
  address?: string;
  emergency_contact?: string;
  emergency_phone?: string;
}

/**
 * 筛选条件接口
 */
interface FilterParams {
  keyword?: string;
  department_id?: string;
  status?: EmployeeStatus;
  position?: string;
  level?: string;
  hire_date_start?: string;
  hire_date_end?: string;
}

/**
 * 员工管理页面组件
 */
export const EmployeePage: React.FC = () => {
  const { t } = useTranslation();
  const { orgId } = useParams<{ orgId: string }>();

  // 状态管理
  const [loading, setLoading] = useState(false);
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [total, setTotal] = useState(0);

  // 分页
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
  });

  // 筛选条件
  const [filters, setFilters] = useState<FilterParams>({});

  // 模态框状态
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [currentEmployee, setCurrentEmployee] = useState<Employee | null>(null);

  // 表单
  const [createForm] = Form.useForm<EmployeeFormData>();
  const [editForm] = Form.useForm<EmployeeFormData>();

  // 加载员工列表
  const loadEmployees = useCallback(async () => {
    setLoading(true);
    try {
      // TODO: 调用实际API
      // const response = await employeeApi.list({
      //   org_id: orgId,
      //   page: pagination.current,
      //   page_size: pagination.pageSize,
      //   ...filters,
      // });

      // 模拟数据
      const mockData: Employee[] = [
        {
          employee_id: '1',
          employee_number: 'EMP001',
          name: '张三',
          name_pinyin: 'zhang san',
          gender: Gender.MALE,
          phone: '13800138001',
          email: 'zhangsan@example.com',
          avatar_url: '',
          department_id: 'dept1',
          department_name: '技术部',
          position: '高级工程师',
          level: 'P7',
          hire_date: '2023-01-15',
          contract_type: ContractType.LABOR,
          contract_number: 'CT202301001',
          contract_start_date: '2023-01-15',
          contract_end_date: '2026-01-14',
          is_probation: false,
          status: EmployeeStatus.ACTIVE,
          id_card_number: '110101199001011234',
          address: '北京市朝阳区',
          emergency_contact: '李四',
          emergency_phone: '13900139001',
          created_at: '2023-01-15T10:00:00Z',
          updated_at: '2023-12-01T10:00:00Z',
          created_by: 'admin',
          updated_by: 'admin',
        },
        {
          employee_id: '2',
          employee_number: 'EMP002',
          name: '李四',
          name_pinyin: 'li si',
          gender: Gender.FEMALE,
          phone: '13800138002',
          email: 'lisi@example.com',
          avatar_url: '',
          department_id: 'dept2',
          department_name: '产品部',
          position: '产品经理',
          level: 'P6',
          hire_date: '2023-06-01',
          contract_type: ContractType.LABOR,
          contract_number: 'CT202306001',
          contract_start_date: '2023-06-01',
          contract_end_date: '2026-05-31',
          is_probation: true,
          probation_end_date: '2023-12-01',
          status: EmployeeStatus.PROBATION,
          id_card_number: '110101199002022345',
          address: '北京市海淀区',
          emergency_contact: '王五',
          emergency_phone: '13900139002',
          created_at: '2023-06-01T10:00:00Z',
          updated_at: '2023-12-01T10:00:00Z',
          created_by: 'admin',
          updated_by: 'admin',
        },
      ];

      setEmployees(mockData);
      setTotal(mockData.length);
    } catch (error) {
      message.error(t('employee.load_failed'));
    } finally {
      setLoading(false);
    }
  }, [orgId, pagination, filters, t]);

  // 加载部门列表
  const loadDepartments = useCallback(async () => {
    try {
      // TODO: 调用实际API
      // const response = await departmentApi.list({ org_id: orgId });

      // 模拟数据
      const mockDepartments: Department[] = [
        { department_id: 'dept1', department_name: '技术部' },
        { department_id: 'dept2', department_name: '产品部' },
        { department_id: 'dept3', department_name: '设计部' },
        { department_id: 'dept4', department_name: '运营部' },
        { department_id: 'dept5', department_name: '市场部' },
      ];
      setDepartments(mockDepartments);
    } catch (error) {
      message.error(t('employee.load_departments_failed'));
    }
  }, [orgId, t]);

  // 初始化
  useEffect(() => {
    loadEmployees();
    loadDepartments();
  }, [loadEmployees, loadDepartments]);

  // 处理表格变化
  const handleTableChange = (newPagination: any) => {
    setPagination({
      current: newPagination.current,
      pageSize: newPagination.pageSize,
    });
  };

  // 处理筛选变化
  const handleFilterChange = (key: string, value: any) => {
    setFilters(prev => ({
      ...prev,
      [key]: value,
    }));
    setPagination({ ...pagination, current: 1 });
  };

  // 重置筛选
  const handleResetFilters = () => {
    setFilters({});
    setPagination({ ...pagination, current: 1 });
  };

  // 打开创建员工对话框
  const handleOpenCreateModal = () => {
    createForm.resetFields();
    setCreateModalVisible(true);
  };

  // 创建员工
  const handleCreateEmployee = async () => {
    try {
      const values = await createForm.validateFields();

      // TODO: 调用实际API
      // await employeeApi.create({
      //   org_id: orgId,
      //   ...values,
      // });

      message.success(t('employee.create_success'));
      setCreateModalVisible(false);
      createForm.resetFields();
      loadEmployees();
    } catch (error) {
      message.error(t('employee.create_failed'));
    }
  };

  // 打开编辑员工对话框
  const handleOpenEditModal = (employee: Employee) => {
    setCurrentEmployee(employee);
    editForm.setFieldsValue({
      name: employee.name,
      gender: employee.gender,
      phone: employee.phone,
      email: employee.email,
      avatar_url: employee.avatar_url,
      department_id: employee.department_id,
      position: employee.position,
      level: employee.level,
      hire_date: dayjs(employee.hire_date),
      contract_type: employee.contract_type,
      contract_number: employee.contract_number,
      contract_start_date: dayjs(employee.contract_start_date),
      contract_end_date: employee.contract_end_date ? dayjs(employee.contract_end_date) : undefined,
      is_probation: employee.is_probation,
      probation_end_date: employee.probation_end_date ? dayjs(employee.probation_end_date) : undefined,
      id_card_number: employee.id_card_number,
      address: employee.address,
      emergency_contact: employee.emergency_contact,
      emergency_phone: employee.emergency_phone,
    });
    setEditModalVisible(true);
  };

  // 编辑员工
  const handleEditEmployee = async () => {
    if (!currentEmployee) return;

    try {
      const values = await editForm.validateFields();

      // TODO: 调用实际API
      // await employeeApi.update({
      //   org_id: orgId,
      //   employee_id: currentEmployee.employee_id,
      //   ...values,
      // });

      message.success(t('employee.update_success'));
      setEditModalVisible(false);
      setCurrentEmployee(null);
      editForm.resetFields();
      loadEmployees();
    } catch (error) {
      message.error(t('employee.update_failed'));
    }
  };

  // 查看员工详情
  const handleViewDetail = (employee: Employee) => {
    setCurrentEmployee(employee);
    setDetailModalVisible(true);
  };

  // 删除员工
  const handleDeleteEmployee = async (employeeId: string) => {
    try {
      // TODO: 调用实际API
      // await employeeApi.delete({
      //   org_id: orgId,
      //   employee_id: employeeId,
      // });

      message.success(t('employee.delete_success'));
      loadEmployees();
    } catch (error) {
      message.error(t('employee.delete_failed'));
    }
  };

  // 更改员工状态
  const handleChangeStatus = async (employeeId: string, newStatus: EmployeeStatus) => {
    try {
      // TODO: 调用实际API
      // await employeeApi.updateStatus({
      //   org_id: orgId,
      //   employee_id: employeeId,
      //   status: newStatus,
      // });

      message.success(t('employee.status_update_success'));
      loadEmployees();
    } catch (error) {
      message.error(t('employee.status_update_failed'));
    }
  };

  // 批量导出
  const handleBatchExport = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning(t('employee.select_to_export'));
      return;
    }

    try {
      // TODO: 调用实际API
      // const response = await employeeApi.export({
      //   org_id: orgId,
      //   employee_ids: selectedRowKeys,
      // });

      message.success(t('employee.export_success'));
      setSelectedRowKeys([]);
    } catch (error) {
      message.error(t('employee.export_failed'));
    }
  };

  // 批量导入
  const handleBatchImport: UploadProps['onChange'] = async info => {
    const { file } = info;
    try {
      // TODO: 调用实际API
      // await employeeApi.import({
      //   org_id: orgId,
      //   file: file.originFileObj,
      // });

      message.success(t('employee.import_success'));
      loadEmployees();
    } catch (error) {
      message.error(t('employee.import_failed'));
    }
  };

  // 表格列定义
  const columns = [
    {
      title: t('employee.avatar'),
      dataIndex: 'avatar_url',
      key: 'avatar',
      width: 80,
      render: (avatar: string, record: Employee) => (
        <Avatar
          size={40}
          src={avatar}
          icon={<UserOutlined />}
          alt={record.name}
        >
          {record.name?.charAt(0)}
        </Avatar>
      ),
    },
    {
      title: t('employee.name'),
      dataIndex: 'name',
      key: 'name',
      width: 120,
      render: (name: string, record: Employee) => (
        <Space direction="vertical" size={0}>
          <Text strong>{name}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.employee_number}
          </Text>
        </Space>
      ),
    },
    {
      title: t('employee.contact'),
      key: 'contact',
      width: 180,
      render: (_, record: Employee) => (
        <Space direction="vertical" size={0}>
          <Text type="secondary" style={{ fontSize: 12 }}>
            <PhoneOutlined /> {record.phone}
          </Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            <MailOutlined /> {record.email}
          </Text>
        </Space>
      ),
    },
    {
      title: t('employee.department'),
      dataIndex: 'department_name',
      key: 'department',
      width: 120,
    },
    {
      title: t('employee.position'),
      key: 'position',
      width: 150,
      render: (_, record: Employee) => (
        <Space direction="vertical" size={0}>
          <Text>{record.position}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.level}
          </Text>
        </Space>
      ),
    },
    {
      title: t('employee.contract'),
      key: 'contract',
      width: 150,
      render: (_, record: Employee) => (
        <Space direction="vertical" size={0}>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {t(`employee.contract_type_${record.contract_type}`)}
          </Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.contract_number}
          </Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {dayjs(record.contract_start_date).format('YYYY-MM-DD')} ~
            {record.contract_end_date
              ? dayjs(record.contract_end_date).format('YYYY-MM-DD')
              : t('employee.unlimited')}
          </Text>
        </Space>
      ),
    },
    {
      title: t('employee.status'),
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: EmployeeStatus, record: Employee) => {
        const statusConfig = {
          [EmployeeStatus.PROBATION]: {
            color: 'blue',
            icon: <StarOutlined />,
            text: t('employee.status_probation'),
          },
          [EmployeeStatus.ACTIVE]: {
            color: 'green',
            icon: <CheckCircleOutlined />,
            text: t('employee.status_active'),
          },
          [EmployeeStatus.SUSPENDED]: {
            color: 'orange',
            icon: <StopOutlined />,
            text: t('employee.status_suspended'),
          },
          [EmployeeStatus.RESIGNED]: {
            color: 'red',
            icon: <RollbackOutlined />,
            text: t('employee.status_resigned'),
          },
        };
        const config = statusConfig[status];
        return (
          <Tag color={config.color} icon={config.icon}>
            {config.text}
          </Tag>
        );
      },
    },
    {
      title: t('employee.hire_date'),
      dataIndex: 'hire_date',
      key: 'hire_date',
      width: 120,
      render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: t('employee.actions'),
      key: 'actions',
      width: 200,
      fixed: 'right' as const,
      render: (_, record: Employee) => {
        const statusMenuItems: MenuProps['items'] = [
          {
            key: 'probation',
            label: t('employee.status_probation'),
            icon: <StarOutlined />,
            disabled: record.status === EmployeeStatus.PROBATION,
            onClick: () => handleChangeStatus(record.employee_id, EmployeeStatus.PROBATION),
          },
          {
            key: 'active',
            label: t('employee.status_active'),
            icon: <CheckCircleOutlined />,
            disabled: record.status === EmployeeStatus.ACTIVE,
            onClick: () => handleChangeStatus(record.employee_id, EmployeeStatus.ACTIVE),
          },
          {
            type: 'divider',
          },
          {
            key: 'suspended',
            label: t('employee.status_suspended'),
            icon: <StopOutlined />,
            disabled: record.status === EmployeeStatus.SUSPENDED,
            onClick: () => handleChangeStatus(record.employee_id, EmployeeStatus.SUSPENDED),
          },
          {
            key: 'resigned',
            label: t('employee.status_resigned'),
            icon: <RollbackOutlined />,
            danger: true,
            disabled: record.status === EmployeeStatus.RESIGNED,
            onClick: () => handleChangeStatus(record.employee_id, EmployeeStatus.RESIGNED),
          },
        ];

        return (
          <Space size="small">
            <Tooltip title={t('employee.view_detail')}>
              <Button
                type="text"
                icon={<FileTextOutlined />}
                onClick={() => handleViewDetail(record)}
              />
            </Tooltip>
            <Tooltip title={t('employee.edit')}>
              <Button
                type="text"
                icon={<EditOutlined />}
                onClick={() => handleOpenEditModal(record)}
              />
            </Tooltip>
            <Dropdown menu={{ items: statusMenuItems }}>
              <Button type="text">
                {t('employee.change_status')} <span className="caret" />
              </Button>
            </Dropdown>
            <Popconfirm
              title={t('employee.delete_confirm')}
              onConfirm={() => handleDeleteEmployee(record.employee_id)}
              okText={t('common.confirm')}
              cancelText={t('common.cancel')}
            >
              <Tooltip title={t('employee.delete')}>
                <Button type="text" danger icon={<DeleteOutlined />} />
              </Tooltip>
            </Popconfirm>
          </Space>
        );
      },
    },
  ];

  // 状态标签配置
  const statusOptions = [
    { label: t('employee.status_probation'), value: EmployeeStatus.PROBATION },
    { label: t('employee.status_active'), value: EmployeeStatus.ACTIVE },
    { label: t('employee.status_suspended'), value: EmployeeStatus.SUSPENDED },
    { label: t('employee.status_resigned'), value: EmployeeStatus.RESIGNED },
  ];

  // 合同类型选项
  const contractTypeOptions = [
    { label: t('employee.contract_type_labor'), value: ContractType.LABOR },
    { label: t('employee.contract_type_internship'), value: ContractType.INTERNSHIP },
    { label: t('employee.contract_type_service'), value: ContractType.SERVICE },
    { label: t('employee.contract_type_outsourcing'), value: ContractType.OUTSOURCING },
  ];

  return (
    <Layout className="employee-page">
      <Header>
        <div className="header-content">
          <Title level={4}>{t('employee.title')}</Title>
          <Space>
            <Button icon={<DownloadOutlined />} onClick={handleBatchExport}>
              {t('employee.batch_export')}
            </Button>
            <Upload
              accept=".xlsx,.xls,.csv"
              showUploadList={false}
              onChange={handleBatchImport}
            >
              <Button icon={<UploadOutlined />}>
                {t('employee.batch_import')}
              </Button>
            </Upload>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreateModal}>
              {t('employee.create')}
            </Button>
          </Space>
        </div>
      </Header>

      <Content className="content">
        <Card>
          {/* 筛选区域 */}
          <Space direction="vertical" size="middle" style={{ width: '100%', marginBottom: 16 }}>
            <Row gutter={16}>
              <Col span={6}>
                <Input
                  placeholder={t('employee.search_placeholder')}
                  prefix={<SearchOutlined />}
                  value={filters.keyword}
                  onChange={e => handleFilterChange('keyword', e.target.value)}
                  allowClear
                />
              </Col>
              <Col span={4}>
                <Select
                  placeholder={t('employee.filter_by_department')}
                  value={filters.department_id}
                  onChange={value => handleFilterChange('department_id', value)}
                  allowClear
                  style={{ width: '100%' }}
                >
                  {departments.map(dept => (
                    <Option key={dept.department_id} value={dept.department_id}>
                      {dept.department_name}
                    </Option>
                  ))}
                </Select>
              </Col>
              <Col span={4}>
                <Select
                  placeholder={t('employee.filter_by_status')}
                  value={filters.status}
                  onChange={value => handleFilterChange('status', value)}
                  allowClear
                  style={{ width: '100%' }}
                >
                  {statusOptions.map(option => (
                    <Option key={option.value} value={option.value}>
                      {option.label}
                    </Option>
                  ))}
                </Select>
              </Col>
              <Col span={4}>
                <Input
                  placeholder={t('employee.filter_by_position')}
                  value={filters.position}
                  onChange={e => handleFilterChange('position', e.target.value)}
                  allowClear
                />
              </Col>
              <Col span={4}>
                <Input
                  placeholder={t('employee.filter_by_level')}
                  value={filters.level}
                  onChange={e => handleFilterChange('level', e.target.value)}
                  allowClear
                />
              </Col>
              <Col span={2}>
                <Button onClick={handleResetFilters}>
                  {t('common.reset')}
                </Button>
              </Col>
            </Row>

            <Row gutter={16}>
              <Col span={6}>
                <DatePicker.RangePicker
                  style={{ width: '100%' }}
                  placeholder={[
                    t('employee.hire_date_start'),
                    t('employee.hire_date_end'),
                  ]}
                  onChange={(dates) => {
                    if (dates && dates[0] && dates[1]) {
                      setFilters(prev => ({
                        ...prev,
                        hire_date_start: dates[0]!.format('YYYY-MM-DD'),
                        hire_date_end: dates[1]!.format('YYYY-MM-DD'),
                      }));
                    } else {
                      setFilters(prev => ({
                        ...prev,
                        hire_date_start: undefined,
                        hire_date_end: undefined,
                      }));
                    }
                  }}
                />
              </Col>
            </Row>
          </Space>

          {/* 统计信息 */}
          <Row gutter={16} style={{ marginBottom: 16 }}>
            <Col span={6}>
              <Card size="small">
                <Space>
                  <Badge count={total} showZero>
                    <Avatar shape="square" icon={<TeamOutlined />} />
                  </Badge>
                  <Space direction="vertical" size={0}>
                    <Text type="secondary">{t('employee.total')}</Text>
                    <Text strong>{total}</Text>
                  </Space>
                </Space>
              </Card>
            </Col>
            <Col span={6}>
              <Card size="small">
                <Space>
                  <Badge
                    count={employees.filter(e => e.status === EmployeeStatus.ACTIVE).length}
                    showZero
                  >
                    <Avatar shape="square" icon={<CheckCircleOutlined />} style={{ backgroundColor: '#52c41a' }} />
                  </Badge>
                  <Space direction="vertical" size={0}>
                    <Text type="secondary">{t('employee.status_active')}</Text>
                    <Text strong>
                      {employees.filter(e => e.status === EmployeeStatus.ACTIVE).length}
                    </Text>
                  </Space>
                </Space>
              </Card>
            </Col>
            <Col span={6}>
              <Card size="small">
                <Space>
                  <Badge
                    count={employees.filter(e => e.status === EmployeeStatus.PROBATION).length}
                    showZero
                  >
                    <Avatar shape="square" icon={<StarOutlined />} style={{ backgroundColor: '#1890ff' }} />
                  </Badge>
                  <Space direction="vertical" size={0}>
                    <Text type="secondary">{t('employee.status_probation')}</Text>
                    <Text strong>
                      {employees.filter(e => e.status === EmployeeStatus.PROBATION).length}
                    </Text>
                  </Space>
                </Space>
              </Card>
            </Col>
            <Col span={6}>
              <Card size="small">
                <Space>
                  <Badge
                    count={employees.filter(e => e.status === EmployeeStatus.RESIGNED).length}
                    showZero
                  >
                    <Avatar shape="square" icon={<RollbackOutlined />} style={{ backgroundColor: '#ff4d4f' }} />
                  </Badge>
                  <Space direction="vertical" size={0}>
                    <Text type="secondary">{t('employee.status_resigned')}</Text>
                    <Text strong>
                      {employees.filter(e => e.status === EmployeeStatus.RESIGNED).length}
                    </Text>
                  </Space>
                </Space>
              </Card>
            </Col>
          </Row>

          {/* 员工列表 */}
          <Table
            rowKey="employee_id"
            columns={columns}
            dataSource={employees}
            loading={loading}
            scroll={{ x: 1500 }}
            pagination={{
              current: pagination.current,
              pageSize: pagination.pageSize,
              total,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: total => t('employee.pagination_total', { total }),
            }}
            onChange={handleTableChange}
            rowSelection={{
              selectedRowKeys,
              onChange: setSelectedRowKeys,
            }}
          />
        </Card>
      </Content>

      {/* 创建员工对话框 */}
      <Modal
        title={t('employee.create_title')}
        open={createModalVisible}
        onOk={handleCreateEmployee}
        onCancel={() => {
          setCreateModalVisible(false);
          createForm.resetFields();
        }}
        width={800}
        okText={t('common.confirm')}
        cancelText={t('common.cancel')}
      >
        <Form form={createForm} layout="vertical">
          <Tabs defaultActiveKey="basic">
            <TabPane tab={t('employee.tab_basic')} key="basic">
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.name')}
                    name="name"
                    rules={[
                      { required: true, message: t('employee.name_required') },
                      { min: 2, max: 20, message: t('employee.name_length_error') },
                    ]}
                  >
                    <Input placeholder={t('employee.name_placeholder')} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.gender')}
                    name="gender"
                    rules={[{ required: true, message: t('employee.gender_required') }]}
                  >
                    <Radio.Group>
                      <Radio value={Gender.MALE}>{t('employee.gender_male')}</Radio>
                      <Radio value={Gender.FEMALE}>{t('employee.gender_female')}</Radio>
                      <Radio value={Gender.OTHER}>{t('employee.gender_other')}</Radio>
                    </Radio.Group>
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.phone')}
                    name="phone"
                    rules={[
                      { required: true, message: t('employee.phone_required') },
                      {
                        pattern: /^1[3-9]\d{9}$/,
                        message: t('employee.phone_format_error'),
                      },
                    ]}
                  >
                    <Input prefix={<PhoneOutlined />} placeholder={t('employee.phone_placeholder')} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.email')}
                    name="email"
                    rules={[
                      { required: true, message: t('employee.email_required') },
                      {
                        type: 'email',
                        message: t('employee.email_format_error'),
                      },
                    ]}
                  >
                    <Input prefix={<MailOutlined />} placeholder={t('employee.email_placeholder')} />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item
                label={t('employee.avatar')}
                name="avatar_url"
              >
                <Upload
                  listType="picture-card"
                  maxCount={1}
                  accept="image/*"
                >
                  <div>
                    <PlusOutlined />
                    <div style={{ marginTop: 8 }}>{t('employee.upload_avatar')}</div>
                  </div>
                </Upload>
              </Form.Item>
            </TabPane>

            <TabPane tab={t('employee.tab_position')} key="position">
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.department')}
                    name="department_id"
                    rules={[{ required: true, message: t('employee.department_required') }]}
                  >
                    <Select placeholder={t('employee.select_department')}>
                      {departments.map(dept => (
                        <Option key={dept.department_id} value={dept.department_id}>
                          {dept.department_name}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.position')}
                    name="position"
                    rules={[{ required: true, message: t('employee.position_required') }]}
                  >
                    <Input placeholder={t('employee.position_placeholder')} />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.level')}
                    name="level"
                    rules={[{ required: true, message: t('employee.level_required') }]}
                  >
                    <Select placeholder={t('employee.select_level')}>
                      <Option value="P4">P4</Option>
                      <Option value="P5">P5</Option>
                      <Option value="P6">P6</Option>
                      <Option value="P7">P7</Option>
                      <Option value="P8">P8</Option>
                      <Option value="P9">P9</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.hire_date')}
                    name="hire_date"
                    rules={[{ required: true, message: t('employee.hire_date_required') }]}
                  >
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
              </Row>
            </TabPane>

            <TabPane tab={t('employee.tab_contract')} key="contract">
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.contract_type')}
                    name="contract_type"
                    rules={[{ required: true, message: t('employee.contract_type_required') }]}
                  >
                    <Select placeholder={t('employee.select_contract_type')}>
                      {contractTypeOptions.map(option => (
                        <Option key={option.value} value={option.value}>
                          {option.label}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.contract_number')}
                    name="contract_number"
                    rules={[{ required: true, message: t('employee.contract_number_required') }]}
                  >
                    <Input placeholder={t('employee.contract_number_placeholder')} />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.contract_start_date')}
                    name="contract_start_date"
                    rules={[{ required: true, message: t('employee.contract_start_date_required') }]}
                  >
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.contract_end_date')}
                    name="contract_end_date"
                  >
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item
                label={t('employee.is_probation')}
                name="is_probation"
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>

              <Form.Item
                noStyle
                shouldUpdate={(prevValues, currentValues) =>
                  prevValues.is_probation !== currentValues.is_probation
                }
              >
                {({ getFieldValue }) =>
                  getFieldValue('is_probation') ? (
                    <Form.Item
                      label={t('employee.probation_end_date')}
                      name="probation_end_date"
                      rules={[{ required: true, message: t('employee.probation_end_date_required') }]}
                    >
                      <DatePicker style={{ width: '100%' }} />
                    </Form.Item>
                  ) : null
                }
              </Form.Item>
            </TabPane>

            <TabPane tab={t('employee.tab_other')} key="other">
              <Form.Item
                label={t('employee.id_card_number')}
                name="id_card_number"
                rules={[
                  {
                    pattern: /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/,
                    message: t('employee.id_card_format_error'),
                  },
                ]}
              >
                <Input prefix={<IdcardOutlined />} placeholder={t('employee.id_card_placeholder')} />
              </Form.Item>

              <Form.Item
                label={t('employee.address')}
                name="address"
              >
                <TextArea
                  rows={3}
                  prefix={<EnvironmentOutlined />}
                  placeholder={t('employee.address_placeholder')}
                />
              </Form.Item>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.emergency_contact')}
                    name="emergency_contact"
                  >
                    <Input
                      prefix={<ContactsOutlined />}
                      placeholder={t('employee.emergency_contact_placeholder')}
                    />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.emergency_phone')}
                    name="emergency_phone"
                    rules={[
                      {
                        pattern: /^1[3-9]\d{9}$/,
                        message: t('employee.phone_format_error'),
                      },
                    ]}
                  >
                    <Input
                      prefix={<PhoneOutlined />}
                      placeholder={t('employee.emergency_phone_placeholder')}
                    />
                  </Form.Item>
                </Col>
              </Row>
            </TabPane>
          </Tabs>
        </Form>
      </Modal>

      {/* 编辑员工对话框 */}
      <Modal
        title={t('employee.edit_title')}
        open={editModalVisible}
        onOk={handleEditEmployee}
        onCancel={() => {
          setEditModalVisible(false);
          setCurrentEmployee(null);
          editForm.resetFields();
        }}
        width={800}
        okText={t('common.confirm')}
        cancelText={t('common.cancel')}
      >
        <Form form={editForm} layout="vertical">
          <Tabs defaultActiveKey="basic">
            <TabPane tab={t('employee.tab_basic')} key="basic">
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.name')}
                    name="name"
                    rules={[
                      { required: true, message: t('employee.name_required') },
                      { min: 2, max: 20, message: t('employee.name_length_error') },
                    ]}
                  >
                    <Input placeholder={t('employee.name_placeholder')} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.gender')}
                    name="gender"
                    rules={[{ required: true, message: t('employee.gender_required') }]}
                  >
                    <Radio.Group>
                      <Radio value={Gender.MALE}>{t('employee.gender_male')}</Radio>
                      <Radio value={Gender.FEMALE}>{t('employee.gender_female')}</Radio>
                      <Radio value={Gender.OTHER}>{t('employee.gender_other')}</Radio>
                    </Radio.Group>
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.phone')}
                    name="phone"
                    rules={[
                      { required: true, message: t('employee.phone_required') },
                      {
                        pattern: /^1[3-9]\d{9}$/,
                        message: t('employee.phone_format_error'),
                      },
                    ]}
                  >
                    <Input prefix={<PhoneOutlined />} placeholder={t('employee.phone_placeholder')} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.email')}
                    name="email"
                    rules={[
                      { required: true, message: t('employee.email_required') },
                      {
                        type: 'email',
                        message: t('employee.email_format_error'),
                      },
                    ]}
                  >
                    <Input prefix={<MailOutlined />} placeholder={t('employee.email_placeholder')} />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item
                label={t('employee.avatar')}
                name="avatar_url"
              >
                <Upload
                  listType="picture-card"
                  maxCount={1}
                  accept="image/*"
                >
                  <div>
                    <PlusOutlined />
                    <div style={{ marginTop: 8 }}>{t('employee.upload_avatar')}</div>
                  </div>
                </Upload>
              </Form.Item>
            </TabPane>

            <TabPane tab={t('employee.tab_position')} key="position">
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.department')}
                    name="department_id"
                    rules={[{ required: true, message: t('employee.department_required') }]}
                  >
                    <Select placeholder={t('employee.select_department')}>
                      {departments.map(dept => (
                        <Option key={dept.department_id} value={dept.department_id}>
                          {dept.department_name}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.position')}
                    name="position"
                    rules={[{ required: true, message: t('employee.position_required') }]}
                  >
                    <Input placeholder={t('employee.position_placeholder')} />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.level')}
                    name="level"
                    rules={[{ required: true, message: t('employee.level_required') }]}
                  >
                    <Select placeholder={t('employee.select_level')}>
                      <Option value="P4">P4</Option>
                      <Option value="P5">P5</Option>
                      <Option value="P6">P6</Option>
                      <Option value="P7">P7</Option>
                      <Option value="P8">P8</Option>
                      <Option value="P9">P9</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.hire_date')}
                    name="hire_date"
                    rules={[{ required: true, message: t('employee.hire_date_required') }]}
                  >
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
              </Row>
            </TabPane>

            <TabPane tab={t('employee.tab_contract')} key="contract">
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.contract_type')}
                    name="contract_type"
                    rules={[{ required: true, message: t('employee.contract_type_required') }]}
                  >
                    <Select placeholder={t('employee.select_contract_type')}>
                      {contractTypeOptions.map(option => (
                        <Option key={option.value} value={option.value}>
                          {option.label}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.contract_number')}
                    name="contract_number"
                    rules={[{ required: true, message: t('employee.contract_number_required') }]}
                  >
                    <Input placeholder={t('employee.contract_number_placeholder')} />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.contract_start_date')}
                    name="contract_start_date"
                    rules={[{ required: true, message: t('employee.contract_start_date_required') }]}
                  >
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.contract_end_date')}
                    name="contract_end_date"
                  >
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item
                label={t('employee.is_probation')}
                name="is_probation"
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>

              <Form.Item
                noStyle
                shouldUpdate={(prevValues, currentValues) =>
                  prevValues.is_probation !== currentValues.is_probation
                }
              >
                {({ getFieldValue }) =>
                  getFieldValue('is_probation') ? (
                    <Form.Item
                      label={t('employee.probation_end_date')}
                      name="probation_end_date"
                      rules={[{ required: true, message: t('employee.probation_end_date_required') }]}
                    >
                      <DatePicker style={{ width: '100%' }} />
                    </Form.Item>
                  ) : null
                }
              </Form.Item>
            </TabPane>

            <TabPane tab={t('employee.tab_other')} key="other">
              <Form.Item
                label={t('employee.id_card_number')}
                name="id_card_number"
                rules={[
                  {
                    pattern: /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/,
                    message: t('employee.id_card_format_error'),
                  },
                ]}
              >
                <Input prefix={<IdcardOutlined />} placeholder={t('employee.id_card_placeholder')} />
              </Form.Item>

              <Form.Item
                label={t('employee.address')}
                name="address"
              >
                <TextArea
                  rows={3}
                  prefix={<EnvironmentOutlined />}
                  placeholder={t('employee.address_placeholder')}
                />
              </Form.Item>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.emergency_contact')}
                    name="emergency_contact"
                  >
                    <Input
                      prefix={<ContactsOutlined />}
                      placeholder={t('employee.emergency_contact_placeholder')}
                    />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label={t('employee.emergency_phone')}
                    name="emergency_phone"
                    rules={[
                      {
                        pattern: /^1[3-9]\d{9}$/,
                        message: t('employee.phone_format_error'),
                      },
                    ]}
                  >
                    <Input
                      prefix={<PhoneOutlined />}
                      placeholder={t('employee.emergency_phone_placeholder')}
                    />
                  </Form.Item>
                </Col>
              </Row>
            </TabPane>
          </Tabs>
        </Form>
      </Modal>

      {/* 员工详情对话框 */}
      <Modal
        title={t('employee.detail_title')}
        open={detailModalVisible}
        onCancel={() => {
          setDetailModalVisible(false);
          setCurrentEmployee(null);
        }}
        footer={null}
        width={900}
      >
        {currentEmployee && (
          <div>
            {/* 员工头像和基本信息 */}
            <div style={{ textAlign: 'center', marginBottom: 24 }}>
              <Avatar size={100} src={currentEmployee.avatar_url} icon={<UserOutlined />}>
                {currentEmployee.name?.charAt(0)}
              </Avatar>
              <Title level={4} style={{ marginTop: 16 }}>
                {currentEmployee.name}
              </Title>
              <Space>
                <Tag color="blue">{currentEmployee.employee_number}</Tag>
                <Tag color="green">{t(`employee.gender_${currentEmployee.gender}`)}</Tag>
                <Tag color="purple">{currentEmployee.level}</Tag>
              </Space>
            </div>

            <Divider />

            <Tabs defaultActiveKey="info">
              <TabPane tab={t('employee.tab_info')} key="info">
                <Descriptions bordered column={2}>
                  <Descriptions.Item label={t('employee.name')}>
                    {currentEmployee.name}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.gender')}>
                    {t(`employee.gender_${currentEmployee.gender}`)}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.phone')}>
                    <Space>
                      <PhoneOutlined />
                      {currentEmployee.phone}
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.email')}>
                    <Space>
                      <MailOutlined />
                      {currentEmployee.email}
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.id_card_number')} span={2}>
                    <Space>
                      <IdcardOutlined />
                      {currentEmployee.id_card_number || '-'}
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.address')} span={2}>
                    <Space>
                      <EnvironmentOutlined />
                      {currentEmployee.address || '-'}
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.emergency_contact')}>
                    <Space>
                      <ContactsOutlined />
                      {currentEmployee.emergency_contact || '-'}
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.emergency_phone')}>
                    <Space>
                      <PhoneOutlined />
                      {currentEmployee.emergency_phone || '-'}
                    </Space>
                  </Descriptions.Item>
                </Descriptions>
              </TabPane>

              <TabPane tab={t('employee.tab_position')} key="position">
                <Descriptions bordered column={2}>
                  <Descriptions.Item label={t('employee.department')}>
                    <Space>
                      <ApartmentOutlined />
                      {currentEmployee.department_name}
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.position')}>
                    {currentEmployee.position}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.level')}>
                    <Tag color="purple">{currentEmployee.level}</Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.status')}>
                    <Tag
                      color={
                        currentEmployee.status === EmployeeStatus.ACTIVE
                          ? 'green'
                          : currentEmployee.status === EmployeeStatus.PROBATION
                          ? 'blue'
                          : currentEmployee.status === EmployeeStatus.SUSPENDED
                          ? 'orange'
                          : 'red'
                      }
                    >
                      {t(`employee.status_${currentEmployee.status}`)}
                    </Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.hire_date')}>
                    <Space>
                      <CalendarOutlined />
                      {dayjs(currentEmployee.hire_date).format('YYYY-MM-DD')}
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.work_years')}>
                    {Math.floor(dayjs().diff(dayjs(currentEmployee.hire_date), 'month', true) / 12)}{' '}
                    {t('employee.years')}
                  </Descriptions.Item>
                </Descriptions>
              </TabPane>

              <TabPane tab={t('employee.tab_contract')} key="contract">
                <Descriptions bordered column={2}>
                  <Descriptions.Item label={t('employee.contract_type')}>
                    {t(`employee.contract_type_${currentEmployee.contract_type}`)}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.contract_number')}>
                    <Space>
                      <SafetyOutlined />
                      {currentEmployee.contract_number}
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.contract_start_date')}>
                    {dayjs(currentEmployee.contract_start_date).format('YYYY-MM-DD')}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.contract_end_date')}>
                    {currentEmployee.contract_end_date
                      ? dayjs(currentEmployee.contract_end_date).format('YYYY-MM-DD')
                      : t('employee.unlimited')}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.contract_duration')} span={2}>
                    {currentEmployee.contract_end_date
                      ? `${Math.floor(
                          dayjs(currentEmployee.contract_end_date).diff(
                            dayjs(currentEmployee.contract_start_date),
                            'month',
                            true,
                          ),
                        )} ${t('employee.months')}`
                      : t('employee.unlimited')}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.is_probation')}>
                    <Tag color={currentEmployee.is_probation ? 'blue' : 'default'}>
                      {currentEmployee.is_probation ? t('common.yes') : t('common.no')}
                    </Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.probation_end_date')}>
                    {currentEmployee.probation_end_date
                      ? dayjs(currentEmployee.probation_end_date).format('YYYY-MM-DD')
                      : '-'}
                  </Descriptions.Item>
                </Descriptions>
              </TabPane>

              <TabPane tab={t('employee.tab_system')} key="system">
                <Descriptions bordered column={2}>
                  <Descriptions.Item label={t('employee.employee_id')}>
                    {currentEmployee.employee_id}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('employee.employee_number')}>
                    {currentEmployee.employee_number}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('common.created_at')}>
                    {dayjs(currentEmployee.created_at).format('YYYY-MM-DD HH:mm:ss')}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('common.updated_at')}>
                    {dayjs(currentEmployee.updated_at).format('YYYY-MM-DD HH:mm:ss')}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('common.created_by')}>
                    {currentEmployee.created_by}
                  </Descriptions.Item>
                  <Descriptions.Item label={t('common.updated_by')}>
                    {currentEmployee.updated_by}
                  </Descriptions.Item>
                </Descriptions>
              </TabPane>
            </Tabs>

            <Divider />

            {/* 操作按钮 */}
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button
                icon={<EditOutlined />}
                onClick={() => {
                  setDetailModalVisible(false);
                  handleOpenEditModal(currentEmployee);
                }}
              >
                {t('employee.edit')}
              </Button>
              {currentEmployee.status === EmployeeStatus.PROBATION && (
                <Button
                  type="primary"
                  icon={<CheckCircleOutlined />}
                  onClick={() => handleChangeStatus(currentEmployee.employee_id, EmployeeStatus.ACTIVE)}
                >
                  {t('employee.become_regular')}
                </Button>
              )}
              {currentEmployee.status === EmployeeStatus.ACTIVE && (
                <Popconfirm
                  title={t('employee.resign_confirm')}
                  onConfirm={() => handleChangeStatus(currentEmployee.employee_id, EmployeeStatus.RESIGNED)}
                  okText={t('common.confirm')}
                  cancelText={t('common.cancel')}
                >
                  <Button danger icon={<RollbackOutlined />}>
                    {t('employee.resign')}
                  </Button>
                </Popconfirm>
              )}
            </Space>
          </div>
        )}
      </Modal>

      <style jsx>{`
        .employee-page {
          padding: 0;
          min-height: 100vh;
          background: #f0f2f5;
        }

        .employee-page .header {
          background: #fff;
          padding: 16px 24px;
          border-bottom: 1px solid #f0f0f0;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .employee-page .header .header-content {
          width: 100%;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .employee-page .content {
          padding: 24px;
        }

        .caret {
          display: inline-block;
          width: 0;
          height: 0;
          margin-left: 4px;
          vertical-align: middle;
          border-top: 4px solid;
          border-right: 4px solid transparent;
          border-left: 4px solid transparent;
        }
      `}</style>
    </Layout>
  );
};

export default EmployeePage;
