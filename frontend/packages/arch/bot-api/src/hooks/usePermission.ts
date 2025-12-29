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

import { useCallback, useEffect, useState } from 'react';

import { permissionApi } from '@coze-arch/bot-api';
import type {
  Role,
  Department,
  CreateRoleRequest,
  UpdateRoleRequest,
  ListRolesRequest,
  AssignRoleRequest,
  ListDepartmentsRequest,
  CheckDataPermissionRequest,
  CheckFieldPermissionRequest,
} from '@coze-arch/bot-api';

/**
 * 角色列表 Hook
 */
export const useRoles = (tenantId?: string, params?: Omit<ListRolesRequest, 'tenant_id'>) => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchRoles = useCallback(async () => {
    if (!tenantId) return;

    setLoading(true);
    setError(null);

    try {
      const response = await permissionApi.listRoles({
        tenant_id: tenantId,
        ...params,
      });
      setRoles(response.roles);
      setTotal(response.total);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [tenantId, params]);

  useEffect(() => {
    fetchRoles();
  }, [fetchRoles]);

  return {
    roles,
    total,
    loading,
    error,
    refetch: fetchRoles,
  };
};

/**
 * 角色详情 Hook
 */
export const useRole = (tenantId?: string, roleId?: string) => {
  const [role, setRole] = useState<Role | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchRole = useCallback(async () => {
    if (!tenantId || !roleId) return;

    setLoading(true);
    setError(null);

    try {
      const response = await permissionApi.getRole({
        tenant_id: tenantId,
        role_id: roleId,
      });
      setRole(response.role);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [tenantId, roleId]);

  useEffect(() => {
    fetchRole();
  }, [fetchRole]);

  return {
    role,
    loading,
    error,
    refetch: fetchRole,
  };
};

/**
 * 创建角色 Hook
 */
export const useCreateRole = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const createRole = useCallback(async (data: CreateRoleRequest) => {
    setLoading(true);
    setError(null);

    try {
      const response = await permissionApi.createRole(data);
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    createRole,
    loading,
    error,
  };
};

/**
 * 更新角色 Hook
 */
export const useUpdateRole = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const updateRole = useCallback(async (tenantId: string, roleId: string, data: UpdateRoleRequest) => {
    setLoading(true);
    setError(null);

    try {
      const response = await permissionApi.updateRole({
        tenant_id: tenantId,
        role_id: roleId,
        ...data,
      });
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    updateRole,
    loading,
    error,
  };
};

/**
 * 删除角色 Hook
 */
export const useDeleteRole = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const deleteRole = useCallback(async (tenantId: string, roleId: string) => {
    setLoading(true);
    setError(null);

    try {
      const response = await permissionApi.deleteRole({
        tenant_id: tenantId,
        role_id: roleId,
      });
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    deleteRole,
    loading,
    error,
  };
};

/**
 * 分配角色 Hook
 */
export const useAssignRole = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const assignRole = useCallback(async (data: AssignRoleRequest) => {
    setLoading(true);
    setError(null);

    try {
      const response = await permissionApi.assignRole(data);
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    assignRole,
    loading,
    error,
  };
};

/**
 * 撤销角色 Hook
 */
export const useRevokeRole = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const revokeRole = useCallback(async (tenantId: string, userId: string, roleId: string) => {
    setLoading(true);
    setError(null);

    try {
      const response = await permissionApi.revokeRole({
        tenant_id: tenantId,
        user_id: userId,
        role_id: roleId,
      });
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    revokeRole,
    loading,
    error,
  };
};

/**
 * 部门列表 Hook
 */
export const useDepartments = (
  tenantId?: string,
  params?: Omit<ListDepartmentsRequest, 'tenant_id'>,
) => {
  const [departments, setDepartments] = useState<Department[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchDepartments = useCallback(async () => {
    if (!tenantId) return;

    setLoading(true);
    setError(null);

    try {
      const response = await permissionApi.listDepartments({
        tenant_id: tenantId,
        ...params,
      });
      setDepartments(response.departments);
      setTotal(response.total);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [tenantId, params]);

  useEffect(() => {
    fetchDepartments();
  }, [fetchDepartments]);

  return {
    departments,
    total,
    loading,
    error,
    refetch: fetchDepartments,
  };
};

/**
 * 数据权限检查 Hook
 */
export const useDataPermission = () => {
  const [checking, setChecking] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const checkDataPermission = useCallback(async (request: CheckDataPermissionRequest) => {
    setChecking(true);
    setError(null);

    try {
      const response = await permissionApi.checkDataPermission(request);
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setChecking(false);
    }
  }, []);

  return {
    checkDataPermission,
    checking,
    error,
  };
};

/**
 * 字段权限检查 Hook
 */
export const useFieldPermission = () => {
  const [checking, setChecking] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const checkFieldPermission = useCallback(async (request: CheckFieldPermissionRequest) => {
    setChecking(true);
    setError(null);

    try {
      const response = await permissionApi.checkFieldPermission(request);
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setChecking(false);
    }
  }, []);

  return {
    checkFieldPermission,
    checking,
    error,
  };
};

/**
 * 权限统计 Hook
 */
export const usePermissionStats = (roles: Role[]) => {
  const stats = {
    total: roles.length,
    system: roles.filter(r => r.is_system).length,
    custom: roles.filter(r => !r.is_system).length,
  };

  return stats;
};
