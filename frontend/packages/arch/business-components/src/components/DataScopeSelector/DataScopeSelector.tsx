// frontend/packages/arch/business-components/src/components/DataScopeSelector/DataScopeSelector.tsx

import React from 'react';
import { Radio } from '@coze-studio/ui-components';
import { useStyles } from './DataScopeSelector.styles';

export type DataScope = 'ALL' | 'DEPARTMENT' | 'OWN' | 'CUSTOM' | 'NONE';

export interface DataScopeOption {
  value: DataScope;
  label: string;
  description: string;
  icon: string;
}

export interface DataScopeSelectorProps {
  value?: DataScope;
  onChange?: (scope: DataScope) => void;
  resourceType?: string;
  disabled?: boolean;
  className?: string;
}

const scopeOptions: DataScopeOption[] = [
  {
    value: 'ALL',
    label: '全部数据',
    description: '可以访问所有数据，包括其他部门创建的数据',
    icon: '🌐',
  },
  {
    value: 'DEPARTMENT',
    label: '本部门数据',
    description: '只能访问本部门成员创建的数据',
    icon: '👥',
  },
  {
    value: 'OWN',
    label: '仅自己',
    description: '只能访问自己创建的数据',
    icon: '👤',
  },
  {
    value: 'CUSTOM',
    label: '自定义',
    description: '根据自定义条件访问数据',
    icon: '⚙️',
  },
  {
    value: 'NONE',
    label: '无权限',
    description: '没有任何访问权限',
    icon: '🚫',
  },
];

export const DataScopeSelector: React.FC<DataScopeSelectorProps> = ({
  value = 'OWN',
  onChange,
  resourceType = '数据',
  disabled = false,
  className,
}) => {
  const classes = useStyles();

  return (
    <div className={`${classes.container} ${className || ''}`}>
      <div className={classes.header}>
        <h3 className={classes.title}>{resourceType}访问权限</h3>
        <p className={classes.subtitle}>选择该角色可以访问的数据范围</p>
      </div>

      <div className={classes.options}>
        {scopeOptions.map(option => (
          <label
            key={option.value}
            className={`${classes.option} ${
              value === option.value ? classes.optionSelected : ''
            }`}
          >
            <Radio
              checked={value === option.value}
              onChange={() => onChange?.(option.value)}
              disabled={disabled}
            />

            <div className={classes.optionContent}>
              <div className={classes.optionHeader}>
                <span className={classes.icon}>{option.icon}</span>
                <span className={classes.label}>{option.label}</span>
              </div>
              <p className={classes.description}>{option.description}</p>
            </div>
          </label>
        ))}
      </div>
    </div>
  );
};
