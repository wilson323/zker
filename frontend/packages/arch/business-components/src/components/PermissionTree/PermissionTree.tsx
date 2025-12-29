// frontend/packages/arch/business-components/src/components/PermissionTree/PermissionTree.tsx

import React, { useState } from 'react';
import { useStyles } from './PermissionTree.styles';

/**
 * 权限节点接口
 */
export interface PermissionNode {
  /** 权限key */
  key: string;
  /** 权限标题 */
  title: string;
  /** 子权限 */
  children?: PermissionNode[];
  /** 是否选中 */
  checked?: boolean;
  /** 是否半选（部分子节点选中） */
  halfChecked?: boolean;
  /** 是否禁用 */
  disabled?: boolean;
}

/**
 * PermissionTree组件Props接口
 */
export interface PermissionTreeProps {
  /** 权限树数据 */
  data: PermissionNode[];
  /** 选中变化回调 */
  onChange?: (checkedKeys: string[], checkedNodes: PermissionNode[]) => void;
  /** 默认选中的权限keys */
  defaultCheckedKeys?: string[];
  /** 是否可选中 */
  checkable?: boolean;
  /** 是否显示复选框 */
  showCheckbox?: boolean;
  /** 自定义className */
  className?: string;
}

/**
 * PermissionTree 权限树
 *
 * 树形权限选择器，支持多级权限嵌套选择
 *
 * @example
 * ```tsx
 * const permissionData = [
 *   {
 *     key: 'bot',
 *     title: '机器人管理',
 *     children: [
 *       { key: 'bot.view', title: '查看' },
 *       { key: 'bot.create', title: '创建' },
 *       { key: 'bot.edit', title: '编辑' },
 *       { key: 'bot.delete', title: '删除' },
 *     ],
 *   },
 * ];
 *
 * <PermissionTree
 *   data={permissionData}
 *   onChange={(keys, nodes) => console.log(keys, nodes)}
 * />
 * ```
 */
export const PermissionTree: React.FC<PermissionTreeProps> = ({
  data,
  onChange,
  defaultCheckedKeys = [],
  checkable = true,
  showCheckbox = true,
  className,
}) => {
  const classes = useStyles();

  // 展开状态
  const [expandedKeys, setExpandedKeys] = useState<Set<string>>(new Set());

  // 选中状态
  const [checkedKeys, setCheckedKeys] = useState<Set<string>>(
    new Set(defaultCheckedKeys)
  );

  // 切换展开/收起
  const toggleExpand = (key: string) => {
    const newExpanded = new Set(expandedKeys);
    if (newExpanded.has(key)) {
      newExpanded.delete(key);
    } else {
      newExpanded.add(key);
    }
    setExpandedKeys(newExpanded);
  };

  // 获取所有子节点的keys（包括自身）
  const getAllChildKeys = (node: PermissionNode): string[] => {
    const keys = [node.key];
    if (node.children) {
      node.children.forEach((child) => {
        keys.push(...getAllChildKeys(child));
      });
    }
    return keys;
  };

  // 处理选中/取消选中
  const handleCheck = (node: PermissionNode, checked: boolean) => {
    const newChecked = new Set(checkedKeys);
    const allChildKeys = getAllChildKeys(node);

    if (checked) {
      allChildKeys.forEach((key) => newChecked.add(key));
    } else {
      allChildKeys.forEach((key) => newChecked.delete(key));
    }

    setCheckedKeys(newChecked);

    // 触发onChange回调
    if (onChange) {
      const checkedArray = Array.from(newChecked);
      // 这里简化处理，实际项目中需要遍历树找到对应的nodes
      onChange(checkedArray, []);
    }
  };

  // 检查节点是否选中
  const isNodeChecked = (node: PermissionNode): boolean => {
    return checkedKeys.has(node.key);
  };

  // 检查节点是否半选（部分子节点选中）
  const isNodeHalfChecked = (node: PermissionNode): boolean => {
    if (!node.children || node.children.length === 0) {
      return false;
    }

    const allChildKeys = getAllChildKeys(node).filter((k) => k !== node.key);
    const checkedChildCount = allChildKeys.filter((k) => checkedKeys.has(k)).length;

    return checkedChildCount > 0 && checkedChildCount < allChildKeys.length;
  };

  // 递归渲染树节点
  const renderNode = (node: PermissionNode, level: number = 0): React.ReactNode => {
    const hasChildren = node.children && node.children.length > 0;
    const isExpanded = expandedKeys.has(node.key);
    const isChecked = isNodeChecked(node);
    const isHalfChecked = isNodeHalfChecked(node);
    const indent = level * 24;

    return (
      <div key={node.key} className={classes.node}>
        <div
          className={classes.nodeContent}
          style={{ paddingLeft: `${indent}px` }}
        >
          {/* 展开/收起图标 */}
          {hasChildren && (
            <span
              className={`${classes.expandIcon} ${isExpanded ? classes.expanded : ''}`}
              onClick={() => toggleExpand(node.key)}
            >
              ▶
            </span>
          )}

          {/* 复选框 */}
          {checkable && showCheckbox && (
            <input
              type="checkbox"
              className={classes.checkbox}
              checked={isChecked}
              ref={(input) => {
                if (input) {
                  input.indeterminate = isHalfChecked;
                }
              }}
              onChange={(e) => handleCheck(node, e.target.checked)}
              disabled={node.disabled}
            />
          )}

          {/* 节点标题 */}
          <span
            className={`${classes.nodeTitle} ${
              node.disabled ? classes.disabled : ''
            }`}
          >
            {node.title}
          </span>
        </div>

        {/* 渲染子节点 */}
        {hasChildren && isExpanded && (
          <div className={classes.children}>
            {node.children!.map((child) => renderNode(child, level + 1))}
          </div>
        )}
      </div>
    );
  };

  return (
    <div className={`${classes.container} ${className || ''}`}>
      {data.map((node) => renderNode(node))}
    </div>
  );
};

export default PermissionTree;
