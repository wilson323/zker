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

import { axiosInstance, type BotAPIRequestConfig } from './axios';
import type {
  // Role APIs
  CreateRoleRequest,
  CreateRoleResponse,
  GetRoleRequest,
  GetRoleResponse,
  UpdateRoleRequest,
  UpdateRoleResponse,
  ListRolesRequest,
  ListRolesResponse,
  DeleteRoleRequest,
  DeleteRoleResponse,
  AssignRoleRequest,
  AssignRoleResponse,
  RevokeRoleRequest,
  RevokeRoleResponse,
  GetUserRolesRequest,
  GetUserRolesResponse,
  // Department APIs
  CreateDepartmentRequest,
  CreateDepartmentResponse,
  GetDepartmentRequest,
  GetDepartmentResponse,
  UpdateDepartmentRequest,
  UpdateDepartmentResponse,
  ListDepartmentsRequest,
  ListDepartmentsResponse,
  DeleteDepartmentRequest,
  DeleteDepartmentResponse,
  AddDepartmentMemberRequest,
  AddDepartmentMemberResponse,
  RemoveDepartmentMemberRequest,
  RemoveDepartmentMemberResponse,
  // Permission Check APIs
  CheckDataPermissionRequest,
  CheckDataPermissionResponse,
  CheckFieldPermissionRequest,
  CheckFieldPermissionResponse,
} from './types/permission.types';

/**
 * 权限管理 API
 * 提供角色、部门、权限检查等功能
 */
class PermissionApiService {
  private readonly baseURL = '/api/v1';

  /**
   * ==================== 角色管理 API ====================
   */

  /**
   * 创建角色
   */
  async createRole(
    request: CreateRoleRequest,
    config?: BotAPIRequestConfig,
  ): Promise<CreateRoleResponse> {
    return axiosInstance.request(
      {
        method: 'POST',
        url: `${this.baseURL}/roles`,
        data: request,
      },
      config,
    );
  }

  /**
   * 获取角色详情
   */
  async getRole(
    request: GetRoleRequest,
    config?: BotAPIRequestConfig,
  ): Promise<GetRoleResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: `${this.baseURL}/roles/${request.role_id}`,
        params: { tenant_id: request.tenant_id },
      },
      config,
    );
  }

  /**
   * 更新角色
   */
  async updateRole(
    request: UpdateRoleRequest,
    config?: BotAPIRequestConfig,
  ): Promise<UpdateRoleResponse> {
    const { role_id, tenant_id, ...updateData } = request;
    return axiosInstance.request(
      {
        method: 'PUT',
        url: `${this.baseURL}/roles/${role_id}`,
        params: { tenant_id },
        data: updateData,
      },
      config,
    );
  }

  /**
   * 列出角色
   */
  async listRoles(
    request: ListRolesRequest,
    config?: BotAPIRequestConfig,
  ): Promise<ListRolesResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: `${this.baseURL}/roles`,
        params: request,
      },
      config,
    );
  }

  /**
   * 删除角色
   */
  async deleteRole(
    request: DeleteRoleRequest,
    config?: BotAPIRequestConfig,
  ): Promise<DeleteRoleResponse> {
    return axiosInstance.request(
      {
        method: 'DELETE',
        url: `${this.baseURL}/roles/${request.role_id}`,
        params: { tenant_id: request.tenant_id },
      },
      config,
    );
  }

  /**
   * 分配角色
   */
  async assignRole(
    request: AssignRoleRequest,
    config?: BotAPIRequestConfig,
  ): Promise<AssignRoleResponse> {
    return axiosInstance.request(
      {
        method: 'POST',
        url: `${this.baseURL}/users/${request.user_id}/roles`,
        data: request,
      },
      config,
    );
  }

  /**
   * 撤销角色
   */
  async revokeRole(
    request: RevokeRoleRequest,
    config?: BotAPIRequestConfig,
  ): Promise<RevokeRoleResponse> {
    return axiosInstance.request(
      {
        method: 'DELETE',
        url: `${this.baseURL}/users/${request.user_id}/roles/${request.role_id}`,
        params: { tenant_id: request.tenant_id },
      },
      config,
    );
  }

  /**
   * 获取用户角色
   */
  async getUserRoles(
    request: GetUserRolesRequest,
    config?: BotAPIRequestConfig,
  ): Promise<GetUserRolesResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: `${this.baseURL}/users/${request.user_id}/roles`,
        params: { tenant_id: request.tenant_id },
      },
      config,
    );
  }

  /**
   * ==================== 部门管理 API ====================
   */

  /**
   * 创建部门
   */
  async createDepartment(
    request: CreateDepartmentRequest,
    config?: BotAPIRequestConfig,
  ): Promise<CreateDepartmentResponse> {
    return axiosInstance.request(
      {
        method: 'POST',
        url: `${this.baseURL}/departments`,
        data: request,
      },
      config,
    );
  }

  /**
   * 获取部门详情
   */
  async getDepartment(
    request: GetDepartmentRequest,
    config?: BotAPIRequestConfig,
  ): Promise<GetDepartmentResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: `${this.baseURL}/departments/${request.department_id}`,
        params: { tenant_id: request.tenant_id },
      },
      config,
    );
  }

  /**
   * 更新部门
   */
  async updateDepartment(
    request: UpdateDepartmentRequest,
    config?: BotAPIRequestConfig,
  ): Promise<UpdateDepartmentResponse> {
    const { department_id, tenant_id, ...updateData } = request;
    return axiosInstance.request(
      {
        method: 'PUT',
        url: `${this.baseURL}/departments/${department_id}`,
        params: { tenant_id },
        data: updateData,
      },
      config,
    );
  }

  /**
   * 列出部门
   */
  async listDepartments(
    request: ListDepartmentsRequest,
    config?: BotAPIRequestConfig,
  ): Promise<ListDepartmentsResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: `${this.baseURL}/departments`,
        params: request,
      },
      config,
    );
  }

  /**
   * 删除部门
   */
  async deleteDepartment(
    request: DeleteDepartmentRequest,
    config?: BotAPIRequestConfig,
  ): Promise<DeleteDepartmentResponse> {
    return axiosInstance.request(
      {
        method: 'DELETE',
        url: `${this.baseURL}/departments/${request.department_id}`,
        params: { tenant_id: request.tenant_id },
      },
      config,
    );
  }

  /**
   * 添加部门成员
   */
  async addDepartmentMember(
    request: AddDepartmentMemberRequest,
    config?: BotAPIRequestConfig,
  ): Promise<AddDepartmentMemberResponse> {
    const { department_id, tenant_id, user_id } = request;
    return axiosInstance.request(
      {
        method: 'POST',
        url: `${this.baseURL}/departments/${department_id}/members/${user_id}`,
        params: { tenant_id },
      },
      config,
    );
  }

  /**
   * 移除部门成员
   */
  async removeDepartmentMember(
    request: RemoveDepartmentMemberRequest,
    config?: BotAPIRequestConfig,
  ): Promise<RemoveDepartmentMemberResponse> {
    const { department_id, tenant_id, user_id } = request;
    return axiosInstance.request(
      {
        method: 'DELETE',
        url: `${this.baseURL}/departments/${department_id}/members/${user_id}`,
        params: { tenant_id },
      },
      config,
    );
  }

  /**
   * ==================== 权限检查 API ====================
   */

  /**
   * 检查数据权限
   */
  async checkDataPermission(
    request: CheckDataPermissionRequest,
    config?: BotAPIRequestConfig,
  ): Promise<CheckDataPermissionResponse> {
    return axiosInstance.request(
      {
        method: 'POST',
        url: `${this.baseURL}/permissions/check-data`,
        data: request,
      },
      config,
    );
  }

  /**
   * 检查字段权限
   */
  async checkFieldPermission(
    request: CheckFieldPermissionRequest,
    config?: BotAPIRequestConfig,
  ): Promise<CheckFieldPermissionResponse> {
    return axiosInstance.request(
      {
        method: 'POST',
        url: `${this.baseURL}/permissions/check-field`,
        data: request,
      },
      config,
    );
  }

  /**
   * 批量检查数据权限
   */
  async batchCheckDataPermission(
    tenantId: string,
    userId: string,
    checks: Array<{
      resource_type: string;
      action: string;
      resource_ids: string[];
    }>,
    config?: BotAPIRequestConfig,
  ): Promise<
    Array<{
      resource_type: string;
      action: string;
      allowed_resource_ids: string[];
      denied_resource_ids: string[];
    }>
  > {
    return Promise.all(
      checks.map(check =>
        this.checkDataPermission(
          {
            tenant_id: tenantId,
            user_id: userId,
            ...check,
          },
          config,
        ),
      ),
    );
  }
}

/**
 * 导出单例实例
 */
export const permissionApi = new PermissionApiService();

/**
 * 导出类型
 */
export type {
  // Role APIs
  CreateRoleRequest,
  CreateRoleResponse,
  GetRoleRequest,
  GetRoleResponse,
  UpdateRoleRequest,
  UpdateRoleResponse,
  ListRolesRequest,
  ListRolesResponse,
  DeleteRoleRequest,
  DeleteRoleResponse,
  AssignRoleRequest,
  AssignRoleResponse,
  RevokeRoleRequest,
  RevokeRoleResponse,
  GetUserRolesRequest,
  GetUserRolesResponse,
  // Department APIs
  CreateDepartmentRequest,
  CreateDepartmentResponse,
  GetDepartmentRequest,
  GetDepartmentResponse,
  UpdateDepartmentRequest,
  UpdateDepartmentResponse,
  ListDepartmentsRequest,
  ListDepartmentsResponse,
  DeleteDepartmentRequest,
  DeleteDepartmentResponse,
  AddDepartmentMemberRequest,
  AddDepartmentMemberResponse,
  RemoveDepartmentMemberRequest,
  RemoveDepartmentMemberResponse,
  // Permission Check APIs
  CheckDataPermissionRequest,
  CheckDataPermissionResponse,
  CheckFieldPermissionRequest,
  CheckFieldPermissionResponse,
};
