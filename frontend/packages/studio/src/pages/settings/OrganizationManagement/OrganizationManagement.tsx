// frontend/packages/studio/src/pages/settings/OrganizationManagement/OrganizationManagement.tsx

import React, { useState } from 'react';
import { Button, Input, Modal, Form, Select } from '@coze-studio/ui-components';
import { useTenantStore } from '@coze-studio/stores';
import './OrganizationManagement.styles.css';

interface Department {
  id: string;
  name: string;
  parent_id: string | null;
  children?: Department[];
}

interface Position {
  id: string;
  name: string;
  description: string;
}

export const OrganizationManagement: React.FC = () => {
  const { currentTenant } = useTenantStore();
  const [activeTab, setActiveTab] = useState<'structure' | 'departments' | 'positions'>('structure');

  // 模拟数据 - 实际应从 API 获取
  const [departments, setDepartments] = useState<Department[]>([
    {
      id: '1',
      name: '技术部',
      parent_id: null,
      children: [
        { id: '1-1', name: '前端开发组', parent_id: '1' },
        { id: '1-2', name: '后端开发组', parent_id: '1' },
      ],
    },
    {
      id: '2',
      name: '产品部',
      parent_id: null,
      children: [
        { id: '2-1', name: '产品设计组', parent_id: '2' },
      ],
    },
  ]);

  const [positions, setPositions] = useState<Position[]>([
    { id: '1', name: '工程师', description: '技术研发岗位' },
    { id: '2', name: '产品经理', description: '产品规划和设计' },
  ]);

  // Modal 状态
  const [isDepartmentModalOpen, setIsDepartmentModalOpen] = useState(false);
  const [isPositionModalOpen, setIsPositionModalOpen] = useState(false);

  // 表单数据
  const [departmentForm, setDepartmentForm] = useState({
    name: '',
    parent_id: '',
  });

  const [positionForm, setPositionForm] = useState({
    name: '',
    description: '',
  });

  const handleAddDepartment = () => {
    if (!departmentForm.name) {
      alert('请输入部门名称');
      return;
    }

    const newDepartment: Department = {
      id: Date.now().toString(),
      name: departmentForm.name,
      parent_id: departmentForm.parent_id || null,
    };

    setDepartments([...departments, newDepartment]);
    setDepartmentForm({ name: '', parent_id: '' });
    setIsDepartmentModalOpen(false);
  };

  const handleAddPosition = () => {
    if (!positionForm.name) {
      alert('请输入岗位名称');
      return;
    }

    const newPosition: Position = {
      id: Date.now().toString(),
      name: positionForm.name,
      description: positionForm.description,
    };

    setPositions([...positions, newPosition]);
    setPositionForm({ name: '', description: '' });
    setIsPositionModalOpen(false);
  };

  const renderDepartmentTree = (deps: Department[], level = 0) => {
    return deps.map((dept) => (
      <div key={dept.id} className="department-tree-item" style={{ paddingLeft: `${level * 24}px` }}>
        <div className="department-node">
          <span className="department-icon">📁</span>
          <span className="department-name">{dept.name}</span>
          <div className="department-actions">
            <Button size="small">编辑</Button>
            <Button size="small" type="danger">删除</Button>
          </div>
        </div>
        {dept.children && renderDepartmentTree(dept.children, level + 1)}
      </div>
    ));
  };

  return (
    <div className="organization-management">
      <div className="org-header">
        <h1>组织管理</h1>
        <p>管理公司的组织架构、部门和岗位信息</p>
      </div>

      <div className="org-tabs">
        <button
          className={`tab ${activeTab === 'structure' ? 'active' : ''}`}
          onClick={() => setActiveTab('structure')}
        >
          组织架构
        </button>
        <button
          className={`tab ${activeTab === 'departments' ? 'active' : ''}`}
          onClick={() => setActiveTab('departments')}
        >
          部门管理
        </button>
        <button
          className={`tab ${activeTab === 'positions' ? 'active' : ''}`}
          onClick={() => setActiveTab('positions')}
        >
          岗位管理
        </button>
      </div>

      {activeTab === 'structure' && (
        <div className="org-section">
          <div className="section-header">
            <h2>组织架构树</h2>
            <Button onClick={() => setIsDepartmentModalOpen(true)}>添加部门</Button>
          </div>
          <div className="organization-tree">
            {renderDepartmentTree(departments)}
          </div>
        </div>
      )}

      {activeTab === 'departments' && (
        <div className="org-section">
          <div className="section-header">
            <h2>部门列表</h2>
            <Button onClick={() => setIsDepartmentModalOpen(true)}>添加部门</Button>
          </div>
          <div className="departments-list">
            {departments.map((dept) => (
              <div key={dept.id} className="department-card">
                <div className="card-header">
                  <h3>{dept.name}</h3>
                  <div className="card-actions">
                    <Button size="small">编辑</Button>
                    <Button size="small" type="danger">删除</Button>
                  </div>
                </div>
                {dept.children && dept.children.length > 0 && (
                  <div className="sub-departments">
                    <h4>下级部门：</h4>
                    {dept.children.map((child) => (
                      <span key={child.id} className="sub-dept-tag">
                        {child.name}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === 'positions' && (
        <div className="org-section">
          <div className="section-header">
            <h2>岗位列表</h2>
            <Button onClick={() => setIsPositionModalOpen(true)}>添加岗位</Button>
          </div>
          <div className="positions-grid">
            {positions.map((position) => (
              <div key={position.id} className="position-card">
                <h3>{position.name}</h3>
                <p>{position.description}</p>
                <div className="card-actions">
                  <Button size="small">编辑</Button>
                  <Button size="small" type="danger">删除</Button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* 添加部门 Modal */}
      <Modal
        open={isDepartmentModalOpen}
        title="添加部门"
        onClose={() => setIsDepartmentModalOpen(false)}
      >
        <Form>
          <Form.Item label="部门名称">
            <Input
              value={departmentForm.name}
              onChange={(e) => setDepartmentForm({ ...departmentForm, name: e.target.value })}
              placeholder="请输入部门名称"
            />
          </Form.Item>
          <Form.Item label="上级部门">
            <Select
              value={departmentForm.parent_id}
              onChange={(value) => setDepartmentForm({ ...departmentForm, parent_id: value })}
              options={[
                { label: '无（顶级部门）', value: '' },
                ...departments.map((d) => ({ label: d.name, value: d.id })),
              ]}
              placeholder="请选择上级部门"
            />
          </Form.Item>
        </Form>
        <div className="modal-actions">
          <Button onClick={() => setIsDepartmentModalOpen(false)}>取消</Button>
          <Button onClick={handleAddDepartment}>确定</Button>
        </div>
      </Modal>

      {/* 添加岗位 Modal */}
      <Modal
        open={isPositionModalOpen}
        title="添加岗位"
        onClose={() => setIsPositionModalOpen(false)}
      >
        <Form>
          <Form.Item label="岗位名称">
            <Input
              value={positionForm.name}
              onChange={(e) => setPositionForm({ ...positionForm, name: e.target.value })}
              placeholder="请输入岗位名称"
            />
          </Form.Item>
          <Form.Item label="岗位描述">
            <textarea
              value={positionForm.description}
              onChange={(e) => setPositionForm({ ...positionForm, description: e.target.value })}
              placeholder="请输入岗位描述"
              rows={4}
              className="form-textarea"
            />
          </Form.Item>
        </Form>
        <div className="modal-actions">
          <Button onClick={() => setIsPositionModalOpen(false)}>取消</Button>
          <Button onClick={handleAddPosition}>确定</Button>
        </div>
      </Modal>
    </div>
  );
};

export default OrganizationManagement;
