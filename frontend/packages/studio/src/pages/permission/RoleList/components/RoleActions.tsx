// frontend/packages/studio/src/pages/permission/RoleList/components/RoleActions.tsx

import React from 'react';
import { Button } from '@coze-studio/ui-components';

/**
 * 角色接口（简化版）
 */
export interface Role {
  role_id: string;
  role_name: string;
  role_code: string;
  role_type: 'system' | 'custom';
  description?: string;
}

/**
 * RoleActions组件Props接口
 */
export interface RoleActionsProps {
  /** 角色数据 */
  role: Role;
  /** 查看详情回调 */
  onViewDetail?: () => void;
  /** 编辑回调 */
  onEdit?: () => void;
  /** 删除回调 */
  onDelete?: () => void;
  /** 复制回调 */
  onCopy?: () => void;
}

/**
 * RoleActions 角色操作按钮
 *
 * 用于角色列表的操作列
 *
 * @example
 * ```tsx
 * <RoleActions
 *   role={role}
 *   onViewDetail={() => navigate(`/permissions/roles/${role.role_id}`)}
 *   onEdit={() => openEditModal(role)}
 *   onDelete={() => deleteRole(role.role_id)}
 * />
 * ```
 */
export const RoleActions: React.FC<RoleActionsProps> = ({
  role,
  onViewDetail,
  onEdit,
  onDelete,
  onCopy,
}) => {
  const isSystemRole = role.role_type === 'system';

  return (
    <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
      {onViewDetail && (
        <Button
          size="sm"
          variant="outline"
          onClick={onViewDetail}
        >
          查看详情
        </Button>
      )}

      {onEdit && !isSystemRole && (
        <Button
          size="sm"
          variant="outline"
          onClick={onEdit}
        >
          编辑
        </Button>
      )}

      {onCopy && (
        <Button
          size="sm"
          variant="text"
          onClick={onCopy}
        >
          复制
        </Button>
      )}

      {onDelete && !isSystemRole && (
        <Button
          size="sm"
          variant="text"
          onClick={onDelete}
          style={{ color: '#F5222D' }}
        >
          删除
        </Button>
      )}
    </div>
  );
};

export default RoleActions;
