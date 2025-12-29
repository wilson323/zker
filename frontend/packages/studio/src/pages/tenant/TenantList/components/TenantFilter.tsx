// frontend/packages/studio/src/pages/tenant/TenantList/components/TenantFilter.tsx

import React, { useState } from 'react';
import { Input } from '@coze-studio/ui-components';
import { useStyles } from './TenantFilter.styles';

/**
 * 筛选条件接口
 */
export interface TenantFilterValue {
  keyword?: string;
  tenantType?: 'individual' | 'team' | 'enterprise';
  status?: 'active' | 'suspended' | 'deleted';
}

/**
 * TenantFilter组件Props接口
 */
export interface TenantFilterProps {
  /** 当前筛选条件 */
  filter: TenantFilterValue;
  /** 筛选变化回调 */
  onChange: (filter: TenantFilterValue) => void;
}

/**
 * TenantFilter 租户筛选器
 *
 * 用于筛选租户列表
 *
 * @example
 * ```tsx
 * <TenantFilter
 *   filter={filter}
 *   onChange={setFilter}
 * />
 * ```
 */
export const TenantFilter: React.FC<TenantFilterProps> = ({ filter, onChange }) => {
  const classes = useStyles();

  const [keyword, setKeyword] = useState(filter.keyword || '');

  const handleKeywordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setKeyword(value);
    onChange({ ...filter, keyword: value || undefined });
  };

  const handleTenantTypeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    onChange({
      ...filter,
      tenantType: value === 'all' ? undefined : (value as any),
    });
  };

  const handleStatusChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    onChange({
      ...filter,
      status: value === 'all' ? undefined : (value as any),
    });
  };

  const handleReset = () => {
    setKeyword('');
    onChange({});
  };

  return (
    <div className={classes.container}>
      <div className={classes.filterGroup}>
        <div className={classes.filterItem}>
          <label className={classes.label}>搜索</label>
          <Input
            placeholder="搜索租户名称"
            value={keyword}
            onChange={handleKeywordChange}
            block
          />
        </div>

        <div className={classes.filterItem}>
          <label className={classes.label}>租户类型</label>
          <select
            className={classes.select}
            value={filter.tenantType || 'all'}
            onChange={handleTenantTypeChange}
          >
            <option value="all">全部</option>
            <option value="individual">个人</option>
            <option value="team">团队</option>
            <option value="enterprise">企业</option>
          </select>
        </div>

        <div className={classes.filterItem}>
          <label className={classes.label}>状态</label>
          <select
            className={classes.select}
            value={filter.status || 'all'}
            onChange={handleStatusChange}
          >
            <option value="all">全部</option>
            <option value="active">正常</option>
            <option value="suspended">暂停</option>
            <option value="deleted">已删除</option>
          </select>
        </div>

        <div className={classes.filterItem}>
          <label className={classes.label}>&nbsp;</label>
          <button className={classes.resetButton} onClick={handleReset}>
            重置
          </button>
        </div>
      </div>
    </div>
  );
};

export default TenantFilter;
