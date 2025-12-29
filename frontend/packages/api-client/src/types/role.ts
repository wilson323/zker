// frontend/packages/api-client/src/types/role.ts

import { PageRequest, PageResponse, EntityID, TenantID, Timestamp } from './common';

/**
 * 权限操作
 */
export enum PermissionAction {
  CREATE = 'create',
  READ = 'read',
  UPDATE = 'update',
  DELETE = 'delete',
  EXECUTE = 'execute',
  MANAGE = 'manage',
}

/**
 * 权限资源
 */
export enum PermissionResource {
  BOT = 'bot',
  KNOWLEDGE = 'knowledge',
  WORKFLOW = 'workflow',
  CONVERSATION = 'conversation',
  USER = 'user',
  ROLE = 'role',
  TENANT = 'tenant',
}

/**
 * 数据权限范围
 */
export enum DataScope {
  ALL = 'all', // 全部数据
  DEPARTMENT = 'department', // 本部门及下级部门
  OWN_DEPARTMENT = 'own_department', // 仅本部门
  OWN = 'own', // 仅自己
  CUSTOM = 'custom', // 自定义
}

/**
 * 权限DTO
 */
export interface PermissionDTO {
  permission_id: EntityID;
  permission_name: string;
  permission_code: string;
  resource: PermissionResource;
  action: PermissionAction;
  description?: string;
  created_at: Timestamp;
}

/**
 * 角色DTO
 */
export interface RoleDTO {
  role_id: EntityID;
  tenant_id: TenantID;
  role_name: string;
  role_code: string;
  description?: string;
  is_system_role: boolean;
  permissions: PermissionDTO[];
  user_count: number;
  created_at: Timestamp;
  updated_at: Timestamp;
  created_by: string;
}

/**
 * 角色过滤条件
 */
export interface RoleFilter extends PageRequest {
  tenant_id?: TenantID;
  role_name?: string;
  role_code?: string;
  is_system_role?: boolean;
}

/**
 * 创建角色请求
 */
export interface CreateRoleRequest {
  tenant_id: TenantID;
  role_name: string;
  role_code: string;
  description?: string;
  permission_ids: EntityID[];
}

/**
 * 更新角色请求
 */
export interface UpdateRoleRequest {
  role_name?: string;
  description?: string;
  permission_ids?: EntityID[];
}

/**
 * 数据权限DTO
 */
export interface DataPermissionDTO {
  permission_id: EntityID;
  tenant_id: TenantID;
  role_id: EntityID;
  resource: PermissionResource;
  scope: DataScope;
  scope_config?: {
    department_ids?: EntityID[];
    user_ids?: EntityID[];
    org_unit_ids?: EntityID[];
  };
  created_at: Timestamp;
}

/**
 * 字段权限DTO
 */
export interface FieldPermissionDTO {
  permission_id: EntityID;
  tenant_id: TenantID;
  role_id: EntityID;
  resource: PermissionResource;
  field_name: string;
  is_readable: boolean;
  is_writable: boolean;
  created_at: Timestamp;
}

/**
 * 角色成员DTO
 */
export interface RoleMemberDTO {
  user_id: EntityID;
  username: string;
  email: string;
  display_name?: string;
  avatar_url?: string;
  joined_at: Timestamp;
}

/**
 * 权限检查请求
 */
export interface CheckPermissionRequest {
  resource: PermissionResource;
  action: PermissionAction;
  resource_id?: EntityID;
}

/**
 * 权限检查响应
 */
export interface CheckPermissionResponse {
  has_permission: boolean;
  reason?: string;
}
