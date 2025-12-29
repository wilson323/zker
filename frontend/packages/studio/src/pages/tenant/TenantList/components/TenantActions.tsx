// frontend/packages/studio/src/pages/tenant/TenantList/components/TenantActions.tsx

import React from 'react';
import { Button } from '@coze-studio/ui-components';

/**
 * 租户接口（简化版）
 */
export interface Tenant {
  tenant_id: string;
  tenant_name: string;
  tenant_type: 'individual' | 'team' | 'enterprise';
  status: 'active' | 'suspended' | 'deleted';
}

/**
 * TenantActions组件Props接口
 */
export interface TenantActionsProps {
  /** 租户数据 */
  tenant: Tenant;
  /** 查看详情回调 */
  onViewDetail?: () => void;
  /** 编辑回调 */
  onEdit?: () => void;
  /** 删除回调 */
  onDelete?: () => void;
  /** 暂停回调 */
  onSuspend?: () => void;
  /** 激活回调 */
  onActivate?: () => void;
}

/**
 * TenantActions 租户操作按钮
 *
 * 用于租户列表的操作列
 *
 * @example
 * ```tsx
 * <TenantActions
 *   tenant={tenant}
 *   onViewDetail={() => navigate(`/tenants/${tenant.tenant_id}`)}
 *   onEdit={() => openEditModal(tenant)}
 *   onDelete={() => deleteTenant(tenant.tenant_id)}
 * />
 * ```
 */
export const TenantActions: React.FC<TenantActionsProps> = ({
  tenant,
  onViewDetail,
  onEdit,
  onDelete,
  onSuspend,
  onActivate,
}) => {
  const isDeleted = tenant.status === 'deleted';
  const isSuspended = tenant.status === 'suspended';
  const isActive = tenant.status === 'active';

  return (
    <div style={{ display: 'flex', gap: '8px' }}>
      {onViewDetail && (
        <Button
          size="sm"
          variant="outline"
          onClick={onViewDetail}
        >
          查看详情
        </Button>
      )}

      {onEdit && !isDeleted && (
        <Button
          size="sm"
          variant="outline"
          onClick={onEdit}
        >
          编辑
        </Button>
      )}

      {onSuspend && isActive && (
        <Button
          size="sm"
          variant="text"
          onClick={onSuspend}
        >
          暂停
        </Button>
      )}

      {onActivate && isSuspended && (
        <Button
          size="sm"
          variant="text"
          onClick={onActivate}
        >
          激活
        </Button>
      )}

      {onDelete && !isDeleted && (
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

export default TenantActions;
