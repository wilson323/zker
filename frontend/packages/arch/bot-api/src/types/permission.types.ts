/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * 数据权限范围枚举
 */
export enum PermissionScope {
  /** 全部数据 */
  ALL = 'ALL',
  /** 本部门及下级部门 */
  DEPARTMENT = 'DEPARTMENT',
  /** 仅自己创建的数据 */
  OWN = 'OWN',
  /** 自定义过滤 */
  CUSTOM = 'CUSTOM',
  /** 无权限 */
  NONE = 'NONE',
}

/**
 * 字段权限级别枚举
 */
export enum FieldPermissionLevel {
  /** 隐藏 - 不显示 */
  HIDDEN = 'hidden',
  /** 只读 - 可见但不可编辑 */
  READONLY = 'readonly',
  /** 可编辑 - 可见且可编辑 */
  EDITABLE = 'editable',
}

/**
 * ==================== 角色相关类型 ====================
 */

/**
 * 角色信息
 */
export interface Role {
  /** 角色ID */
  role_id: string;
  /** 租户ID */
  tenant_id: string;
  /** 角色名称 */
  role_name: string;
  /** 角色代码（唯一标识） */
  role_code: string;
  /** 角色描述 */
  description?: string;
  /** 是否系统角色 */
  is_system: boolean;
  /** 功能权限列表 */
  permissions: string[];
  /** 创建时间 */
  created_at: string;
  /** 更新时间 */
  updated_at: string;
}

/**
 * 数据权限配置
 */
export interface DataPermission {
  /** 权限ID */
  permission_id: string;
  /** 角色ID */
  role_id: string;
  /** 资源类型 */
  resource_type: string;
  /** 权限范围 */
  scope: PermissionScope;
  /** 自定义过滤器（scope=CUSTOM时） */
  custom_filters?: string[];
}

/**
 * 字段权限配置
 */
export interface FieldPermission {
  /** 权限ID */
  permission_id: string;
  /** 角色ID */
  role_id: string;
  /** 资源类型 */
  resource_type: string;
  /** 字段名 */
  field_name: string;
  /** 权限级别 */
  permission_level: FieldPermissionLevel;
}

/**
 * 用户角色关联
 */
export interface UserRole {
  /** 关联ID */
  id: string;
  /** 租户ID */
  tenant_id: string;
  /** 用户ID */
  user_id: string;
  /** 角色ID */
  role_id: string;
  /** 是否为主角色 */
  is_primary: boolean;
  /** 分配时间 */
  assigned_at: string;
  /** 分配人ID */
  assigned_by: string;
}

/**
 * 创建角色请求
 */
export interface CreateRoleRequest {
  /** 租户ID */
  tenant_id: string;
  /** 角色名称（必填，2-50字符） */
  role_name: string;
  /** 角色代码（必填，2-50字符，字母数字下划线） */
  role_code: string;
  /** 角色描述（可选，最大500字符） */
  description?: string;
  /** 功能权限列表（必填） */
  permissions: string[];
  /** 数据权限配置（可选） */
  data_permissions?: Omit<DataPermission, 'permission_id' | 'role_id'>[];
  /** 字段权限配置（可选） */
  field_permissions?: Omit<FieldPermission, 'permission_id' | 'role_id'>[];
}

/**
 * 创建角色响应
 */
export interface CreateRoleResponse {
  /** 创建的角色 */
  role: Role;
  /** 数据权限 */
  data_permissions?: DataPermission[];
  /** 字段权限 */
  field_permissions?: FieldPermission[];
}

/**
 * 获取角色详情请求
 */
export interface GetRoleRequest {
  /** 租户ID */
  tenant_id: string;
  /** 角色ID */
  role_id: string;
}

/**
 * 获取角色详情响应
 */
export interface GetRoleResponse {
  /** 角色信息 */
  role: Role;
  /** 数据权限列表 */
  data_permissions: DataPermission[];
  /** 字段权限列表 */
  field_permissions: FieldPermission[];
}

/**
 * 更新角色请求
 */
export interface UpdateRoleRequest {
  /** 租户ID */
  tenant_id: string;
  /** 角色ID */
  role_id: string;
  /** 角色名称（可选） */
  role_name?: string;
  /** 角色描述（可选） */
  description?: string;
  /** 功能权限列表（可选） */
  permissions?: string[];
  /** 数据权限配置（可选） */
  data_permissions?: Omit<DataPermission, 'permission_id' | 'role_id'>[];
  /** 字段权限配置（可选） */
  field_permissions?: Omit<FieldPermission, 'permission_id' | 'role_id'>[];
}

/**
 * 更新角色响应
 */
export interface UpdateRoleResponse {
  /** 更新后的角色 */
  role: Role;
  /** 数据权限列表 */
  data_permissions: DataPermission[];
  /** 字段权限列表 */
  field_permissions: FieldPermission[];
}

/**
 * 列出角色请求
 */
export interface ListRolesRequest {
  /** 租户ID */
  tenant_id: string;
  /** 页码（从1开始） */
  page?: number;
  /** 每页数量（1-100） */
  page_size?: number;
  /** 是否仅系统角色 */
  is_system?: boolean;
  /** 搜索关键词（角色名称或代码） */
  search?: string;
}

/**
 * 列出角色响应
 */
export interface ListRolesResponse {
  /** 角色列表 */
  roles: Role[];
  /** 总数 */
  total: number;
  /** 当前页 */
  page: number;
  /** 每页数量 */
  page_size: number;
}

/**
 * 删除角色请求
 */
export interface DeleteRoleRequest {
  /** 租户ID */
  tenant_id: string;
  /** 角色ID */
  role_id: string;
}

/**
 * 删除角色响应
 */
export interface DeleteRoleResponse {
  /** 是否成功 */
  success: boolean;
  /** 消息 */
  message: string;
}

/**
 * 分配角色请求
 */
export interface AssignRoleRequest {
  /** 租户ID */
  tenant_id: string;
  /** 用户ID */
  user_id: string;
  /** 角色ID */
  role_id: string;
  /** 是否为主角色 */
  is_primary?: boolean;
}

/**
 * 分配角色响应
 */
export interface AssignRoleResponse {
  /** 用户角色关联 */
  user_role: UserRole;
  /** 角色信息 */
  role: Role;
}

/**
 * 撤销角色请求
 */
export interface RevokeRoleRequest {
  /** 租户ID */
  tenant_id: string;
  /** 用户ID */
  user_id: string;
  /** 角色ID */
  role_id: string;
}

/**
 * 撤销角色响应
 */
export interface RevokeRoleResponse {
  /** 是否成功 */
  success: boolean;
  /** 消息 */
  message: string;
}

/**
 * 获取用户角色请求
 */
export interface GetUserRolesRequest {
  /** 租户ID */
  tenant_id: string;
  /** 用户ID */
  user_id: string;
}

/**
 * 获取用户角色响应
 */
export interface GetUserRolesResponse {
  /** 用户角色列表 */
  user_roles: UserRole[];
  /** 角色详情列表 */
  roles: Role[];
}

/**
 * ==================== 部门相关类型 ====================
 */

/**
 * 部门信息
 */
export interface Department {
  /** 部门ID */
  department_id: string;
  /** 租户ID */
  tenant_id: string;
  /** 部门名称 */
  department_name: string;
  /** 父部门ID */
  parent_department_id?: string;
  /** 部门层级（0为根部门） */
  level: number;
  /** 部门路径（逗号分隔的部门ID） */
  path: string;
  /** 部门描述 */
  description?: string;
  /** 创建时间 */
  created_at: string;
  /** 更新时间 */
  updated_at: string;
}

/**
 * 部门成员
 */
export interface DepartmentMember {
  /** 用户ID */
  user_id: string;
  /** 部门ID */
  department_id: string;
  /** 是否为部门负责人 */
  is_leader: boolean;
  /** 加入时间 */
  joined_at: string;
}

/**
 * 创建部门请求
 */
export interface CreateDepartmentRequest {
  /** 租户ID */
  tenant_id: string;
  /** 部门名称（必填，2-50字符） */
  department_name: string;
  /** 父部门ID（可选，不填则为根部门） */
  parent_department_id?: string;
  /** 部门描述（可选，最大500字符） */
  description?: string;
}

/**
 * 创建部门响应
 */
export interface CreateDepartmentResponse {
  /** 创建的部门 */
  department: Department;
}

/**
 * 获取部门详情请求
 */
export interface GetDepartmentRequest {
  /** 租户ID */
  tenant_id: string;
  /** 部门ID */
  department_id: string;
}

/**
 * 获取部门详情响应
 */
export interface GetDepartmentResponse {
  /** 部门信息 */
  department: Department;
  /** 父部门信息 */
  parent_department?: Department;
  /** 子部门列表 */
  child_departments: Department[];
  /** 部门成员列表 */
  members: DepartmentMember[];
}

/**
 * 更新部门请求
 */
export interface UpdateDepartmentRequest {
  /** 租户ID */
  tenant_id: string;
  /** 部门ID */
  department_id: string;
  /** 部门名称（可选） */
  department_name?: string;
  /** 部门描述（可选） */
  description?: string;
}

/**
 * 更新部门响应
 */
export interface UpdateDepartmentResponse {
  /** 更新后的部门 */
  department: Department;
}

/**
 * 列出部门请求
 */
export interface ListDepartmentsRequest {
  /** 租户ID */
  tenant_id: string;
  /** 父部门ID（可选，不填则返回所有根部门） */
  parent_department_id?: string;
  /** 是否递归返回所有子部门 */
  recursive?: boolean;
  /** 搜索关键词（部门名称） */
  search?: string;
}

/**
 * 列出部门响应
 */
export interface ListDepartmentsResponse {
  /** 部门列表 */
  departments: Department[];
  /** 总数 */
  total: number;
}

/**
 * 删除部门请求
 */
export interface DeleteDepartmentRequest {
  /** 租户ID */
  tenant_id: string;
  /** 部门ID */
  department_id: string;
}

/**
 * 删除部门响应
 */
export interface DeleteDepartmentResponse {
  /** 是否成功 */
  success: boolean;
  /** 消息 */
  message: string;
}

/**
 * 添加部门成员请求
 */
export interface AddDepartmentMemberRequest {
  /** 租户ID */
  tenant_id: string;
  /** 部门ID */
  department_id: string;
  /** 用户ID */
  user_id: string;
  /** 是否为部门负责人 */
  is_leader?: boolean;
}

/**
 * 添加部门成员响应
 */
export interface AddDepartmentMemberResponse {
  /** 部门成员 */
  member: DepartmentMember;
}

/**
 * 移除部门成员请求
 */
export interface RemoveDepartmentMemberRequest {
  /** 租户ID */
  tenant_id: string;
  /** 部门ID */
  department_id: string;
  /** 用户ID */
  user_id: string;
}

/**
 * 移除部门成员响应
 */
export interface RemoveDepartmentMemberResponse {
  /** 是否成功 */
  success: boolean;
  /** 消息 */
  message: string;
}

/**
 * ==================== 权限检查相关类型 ====================
 */

/**
 * 检查数据权限请求
 */
export interface CheckDataPermissionRequest {
  /** 租户ID */
  tenant_id: string;
  /** 用户ID */
  user_id: string;
  /** 资源类型（bots, workflows, conversations等） */
  resource_type: string;
  /** 操作类型（create, view, edit, delete等） */
  action: string;
  /** 资源ID列表 */
  resource_ids: string[];
}

/**
 * 检查数据权限响应
 */
export interface CheckDataPermissionResponse {
  /** 允许访问的资源ID列表 */
  allowed_resource_ids: string[];
  /** 拒绝访问的资源ID列表 */
  denied_resource_ids: string[];
  /** 权限详情 */
  details: Array<{
    resource_id: string;
    allowed: boolean;
    reason?: string;
  }>;
}

/**
 * 检查字段权限请求
 */
export interface CheckFieldPermissionRequest {
  /** 租户ID */
  tenant_id: string;
  /** 用户ID */
  user_id: string;
  /** 资源类型 */
  resource_type: string;
  /** 操作类型（view, edit） */
  action: string;
  /** 字段名列表 */
  field_names: string[];
}

/**
 * 检查字段权限响应
 */
export interface CheckFieldPermissionResponse {
  /** 字段权限映射 */
  field_permissions: Record<
    string,
    {
      /** 字段名 */
      field_name: string;
      /** 权限级别 */
      permission_level: FieldPermissionLevel;
    }
  >;
}
