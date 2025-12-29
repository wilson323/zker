// frontend/packages/studio/src/pages/permission/RoleDetail/components/DataPermissionConfig.tsx

import React, { useState } from 'react';
import { PermissionTree, PermissionNode } from '@coze-studio/business-components';
import { useStyles } from './DataPermissionConfig.styles';

/**
 * 权限配置接口
 */
interface PermissionConfig {
  resourceType: string;
  scope: 'ALL' | 'DEPARTMENT' | 'OWN' | 'CUSTOM' | 'NONE';
}

/**
 * DataPermissionConfig组件Props接口
 */
export interface DataPermissionConfigProps {
  /** 配置变化回调 */
  onChange?: (configs: PermissionConfig[]) => void;
}

/**
 * DataPermissionConfig 数据权限配置
 *
 * 用于配置角色的数据权限范围
 *
 * @example
 * ```tsx
 * <DataPermissionConfig
 *   onChange={(configs) => console.log(configs)}
 * />
 * ```
 */
export const DataPermissionConfig: React.FC<DataPermissionConfigProps> = ({
  onChange,
}) => {
  const classes = useStyles();

  // 初始权限配置
  const [permissions, setPermissions] = useState<PermissionConfig[]>([
    { resourceType: 'bots', scope: 'ALL' },
    { resourceType: 'conversations', scope: 'OWN' },
    { resourceType: 'knowledge', scope: 'DEPARTMENT' },
    { resourceType: 'workflows', scope: 'ALL' },
    { resourceType: 'plugins', scope: 'ALL' },
  ]);

  // 资源类型选项
  const resourceTypeOptions: { label: string; value: string }[] = [
    { label: '机器人', value: 'bots' },
    { label: '对话', value: 'conversations' },
    { label: '知识库', value: 'knowledge' },
    { label: '工作流', value: 'workflows' },
    { label: '插件', value: 'plugins' },
  ];

  // 数据范围选项
  const scopeOptions: { label: string; value: PermissionConfig['scope']; description: string }[] = [
    { label: '全部数据', value: 'ALL', description: '可以访问所有数据' },
    { label: '本部门数据', value: 'DEPARTMENT', description: '只能访问本部门的数据' },
    { label: '仅自己', value: 'OWN', description: '只能访问自己创建的数据' },
    { label: '自定义', value: 'CUSTOM', description: '根据自定义条件访问数据' },
    { label: '无权限', value: 'NONE', description: '没有任何访问权限' },
  ];

  const handleScopeChange = (index: number, scope: PermissionConfig['scope']) => {
    const newPermissions = [...permissions];
    newPermissions[index].scope = scope;
    setPermissions(newPermissions);
    onChange?.(newPermissions);
  };

  return (
    <div className={classes.container}>
      <h3 className={classes.title}>数据权限配置</h3>
      <p className={classes.description}>
        为不同类型的资源配置数据访问范围，控制角色可以看到哪些数据
      </p>

      <div className={classes.permissionList}>
        {permissions.map((perm, index) => {
          const resourceLabel = resourceTypeOptions.find(
            (opt) => opt.value === perm.resourceType
          )?.label;

          return (
            <div key={index} className={classes.permissionItem}>
              <div className={classes.permissionHeader}>
                <span className={classes.resourceType}>{resourceLabel}</span>
              </div>

              <div className={classes.scopeOptions}>
                {scopeOptions.map((option) => (
                  <label
                    key={option.value}
                    className={`${classes.scopeOption} ${
                      perm.scope === option.value ? classes.selected : ''
                    }`}
                  >
                    <input
                      type="radio"
                      name={`scope-${index}`}
                      value={option.value}
                      checked={perm.scope === option.value}
                      onChange={() => handleScopeChange(index, option.value)}
                    />
                    <div className={classes.optionContent}>
                      <span className={classes.optionLabel}>{option.label}</span>
                      <span className={classes.optionDescription}>
                        {option.description}
                      </span>
                    </div>
                  </label>
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default DataPermissionConfig;
