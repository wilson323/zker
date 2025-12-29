// frontend/packages/arch/business-components/src/components/TenantSelector/TenantSelector.tsx

import React from 'react';
import { useStyles } from './TenantSelector.styles';

/**
 * 租户类型定义
 */
export interface Tenant {
  tenant_id: string;
  tenant_name: string;
  tenant_type: 'individual' | 'team' | 'enterprise';
  status: 'active' | 'suspended' | 'deleted';
}

/**
 * TenantSelector组件Props接口
 */
export interface TenantSelectorProps {
  /** 当前选中的租户ID */
  value?: string;
  /** 选择变化回调 */
  onChange?: (tenantId: string) => void;
  /** 是否禁用 */
  disabled?: boolean;
  /** 占位文本 */
  placeholder?: string;
  /** 是否允许清除 */
  allowClear?: boolean;
  /** 租户列表数据 */
  tenants?: Tenant[];
  /** 是否加载中 */
  loading?: boolean;
}

/**
 * TenantSelector 租户选择器
 *
 * 用于选择租户的业务组件
 *
 * @example
 * ```tsx
 * <TenantSelector
 *   value={selectedTenantId}
 *   onChange={setSelectedTenantId}
 *   tenants={tenantList}
 *   loading={isLoading}
 * />
 * ```
 */
export const TenantSelector: React.FC<TenantSelectorProps> = ({
  value,
  onChange,
  disabled = false,
  placeholder = '请选择租户',
  allowClear = true,
  tenants = [],
  loading = false,
}) => {
  const classes = useStyles();

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const tenantId = e.target.value;
    onChange?.(tenantId);
  };

  const handleClear = () => {
    onChange?.('');
  };

  const getTenantTypeLabel = (type: Tenant['tenant_type']) => {
    const typeMap = {
      individual: '个人',
      team: '团队',
      enterprise: '企业',
    };
    return typeMap[type] || type;
  };

  const getTenantTypeColor = (type: Tenant['tenant_type']) => {
    const colorMap = {
      individual: classes.tagBlue,
      team: classes.tagGreen,
      enterprise: classes.tagPurple,
    };
    return colorMap[type] || classes.tagGray;
  };

  return (
    <div className={classes.container}>
      <select
        className={classes.select}
        value={value}
        onChange={handleChange}
        disabled={disabled || loading}
      >
        <option value="" disabled>
          {placeholder}
        </option>
        {tenants.map((tenant) => (
          <option key={tenant.tenant_id} value={tenant.tenant_id}>
            {tenant.tenant_name} - {getTenantTypeLabel(tenant.tenant_type)}
          </option>
        ))}
      </select>
      {allowClear && value && !disabled && (
        <button
          className={classes.clearButton}
          onClick={handleClear}
          type="button"
          aria-label="清除"
        >
          ×
        </button>
      )}
      {loading && <div className={classes.loading} />}
    </div>
  );
};

export default TenantSelector;
