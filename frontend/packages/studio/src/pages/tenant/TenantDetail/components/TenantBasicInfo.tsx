// frontend/packages/studio/src/pages/tenant/TenantDetail/components/TenantBasicInfo.tsx

import React from 'react';
import { Input } from '@coze-studio/ui-components';
import { useStyles } from './TenantBasicInfo.styles';

/**
 * 租户数据接口
 */
export interface TenantData {
  tenant_id: string;
  tenant_name: string;
  tenant_type: 'individual' | 'team' | 'enterprise';
  status: 'active' | 'suspended' | 'deleted';
  contact_email?: string;
  contact_phone?: string;
  description?: string;
}

/**
 * TenantBasicInfo组件Props接口
 */
export interface TenantBasicInfoProps {
  /** 租户数据 */
  tenant: TenantData;
  /** 更新回调 */
  onUpdate?: (data: Partial<TenantData>) => void;
}

/**
 * TenantBasicInfo 租户基本信息
 *
 * 显示和编辑租户的基本信息
 */
export const TenantBasicInfo: React.FC<TenantBasicInfoProps> = ({
  tenant,
  onUpdate,
}) => {
  const classes = useStyles();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // 这里应该调用更新API
    console.log('Update tenant basic info');
    onUpdate?.(tenant);
  };

  return (
    <div className={classes.container}>
      <h3 className={classes.title}>基本信息</h3>

      <form className={classes.form} onSubmit={handleSubmit}>
        <div className={classes.formRow}>
          <label className={classes.label}>租户ID</label>
          <input
            type="text"
            className={classes.readonlyInput}
            value={tenant.tenant_id}
            disabled
          />
        </div>

        <div className={classes.formRow}>
          <label className={classes.label}>租户名称 *</label>
          <Input
            placeholder="请输入租户名称"
            defaultValue={tenant.tenant_name}
            block
          />
        </div>

        <div className={classes.formRow}>
          <label className={classes.label}>租户类型 *</label>
          <select className={classes.select} defaultValue={tenant.tenant_type}>
            <option value="individual">个人</option>
            <option value="team">团队</option>
            <option value="enterprise">企业</option>
          </select>
        </div>

        <div className={classes.formRow}>
          <label className={classes.label}>状态</label>
          <select className={classes.select} defaultValue={tenant.status}>
            <option value="active">正常</option>
            <option value="suspended">暂停</option>
            <option value="deleted">已删除</option>
          </select>
        </div>

        <div className={classes.formRow}>
          <label className={classes.label}>联系邮箱</label>
          <Input
            type="email"
            placeholder="请输入联系邮箱"
            defaultValue={tenant.contact_email}
            block
          />
        </div>

        <div className={classes.formRow}>
          <label className={classes.label}>联系电话</label>
          <Input
            type="tel"
            placeholder="请输入联系电话"
            defaultValue={tenant.contact_phone}
            block
          />
        </div>

        <div className={classes.formRow}>
          <label className={classes.label}>描述</label>
          <textarea
            className={classes.textarea}
            placeholder="请输入租户描述"
            rows={4}
            defaultValue={tenant.description}
          />
        </div>

        <div className={classes.footer}>
          <button type="submit" className={classes.submitButton}>
            保存
          </button>
        </div>
      </form>
    </div>
  );
};

export default TenantBasicInfo;
