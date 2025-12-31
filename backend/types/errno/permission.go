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

package errno

import (
	"net/http"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx/code"
)

// Permission: 108 000 000 ~ 108 999 999
const (
	// 通用错误 (108 000 000 ~ 108 009 999)
	ErrPermissionInvalidParamCode        = 108000001 // 无效参数
	ErrPermissionCheckFailedCode         = 108000002 // 权限检查失败
	ErrPermissionDeniedCode              = 108000003 // 权限拒绝
	ErrResourceTypeNotSupportedCode      = 108000004 // 不支持的资源类型

	// 角色相关错误 (108 010 000 ~ 108 019 999)
	ErrRoleNotFoundCode                  = 108010001 // 角色不存在
	ErrRoleAlreadyExistsCode             = 108010002 // 角色已存在
	ErrInvalidRoleCode                   = 108010003 // 无效角色
	ErrRoleNameInvalidCode               = 108010004 // 角色名称无效
	ErrRoleDescriptionInvalidCode        = 108010005 // 角色描述无效
	ErrCannotDeleteSystemRoleCode        = 108010006 // 不能删除系统角色
	ErrRoleInUseCode                     = 108010007 // 角色使用中
	ErrCircularRoleDependencyCode        = 108010008 // 循环角色依赖

	// 角色分配错误 (108 020 000 ~ 108 029 999)
	ErrRoleAssignedCode                  = 108020001 // 角色已分配
	ErrDuplicateRoleAssignmentCode       = 108020002 // 重复角色分配
	ErrUserRoleLimitExceededCode         = 108020003 // 用户角色数量超限
	ErrInvalidRoleAssignmentCode         = 108020004 // 无效角色分配

	// 数据权限错误 (108 030 000 ~ 108 039 999)
	ErrDataPermissionDeniedCode          = 108030001 // 数据权限拒绝
	ErrDataPermissionNotFoundCode        = 108030002 // 数据权限配置不存在
	ErrDataPermissionScopeInvalidCode    = 108030003 // 数据权限范围无效
	ErrInvalidPermissionFilterCode       = 108030004 // 无效权限过滤器
	ErrCustomFilterSyntaxErrorCode       = 108030005 // 自定义过滤器语法错误

	// 字段权限错误 (108 040 000 ~ 108 049 999)
	ErrFieldPermissionDeniedCode         = 108040001 // 字段权限拒绝
	ErrFieldPermissionNotFoundCode       = 108040002 // 字段权限配置不存在
	ErrFieldPermissionInvalidCode        = 108040003 // 字段权限无效
	ErrFieldNotFoundCode                 = 108040004 // 字段不存在
	ErrFieldPermissionLevelInvalidCode   = 108040005 // 字段权限级别无效

	// 部门权限错误 (108 050 000 ~ 108 059 999)
	ErrDepartmentNotFoundCode            = 108050001 // 部门不存在
	ErrUserNotInDepartmentCode           = 108050002 // 用户不在部门中
	ErrDepartmentAccessDeniedCode        = 108050003 // 部门访问拒绝
	ErrInvalidDepartmentStructureCode    = 108050004 // 无效部门结构

	// 租户用户错误 (108 060 000 ~ 108 069 999)
	ErrUserNotInTenantCode               = 108060001 // 用户不在租户中
	ErrTenantUserNotFoundCode            = 108060002 // 租户用户不存在
	ErrTenantUserAlreadyExistsCode       = 108060003 // 租户用户已存在

	// 权限继承错误 (108 070 000 ~ 108 079 999)
	ErrPermissionInheritanceCycleCode    = 108070001 // 权限继承循环
	ErrPermissionInheritanceDepthCode    = 108070002 // 权限继承深度超限
	ErrBasePermissionNotFoundCode        = 108070003 // 基础权限不存在

	// 批量操作错误 (108 080 000 ~ 108 089 999)
	ErrBatchOperationFailedCode          = 108080001 // 批量操作失败
	ErrBatchOperationPartialFailureCode  = 108080002 // 批量操作部分失败
)

func init() {
	// 通用错误
	code.Register(
		ErrPermissionInvalidParamCode,
		"invalid parameter: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrPermissionCheckFailedCode,
		"permission check failed: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrPermissionDeniedCode,
		"permission denied: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrResourceTypeNotSupportedCode,
		"resource type not supported: {msg}",
		code.WithAffectStability(false),
	)

	// 角色相关错误
	code.Register(
		ErrRoleNotFoundCode,
		"role not found: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrRoleAlreadyExistsCode,
		"role already exists: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrInvalidRoleCode,
		"invalid role: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrRoleNameInvalidCode,
		"role name is invalid: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrRoleDescriptionInvalidCode,
		"role description is invalid: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrCannotDeleteSystemRoleCode,
		"cannot delete system role: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrRoleInUseCode,
		"role is in use and cannot be deleted: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrCircularRoleDependencyCode,
		"circular role dependency detected: {msg}",
		code.WithAffectStability(false),
	)

	// 角色分配错误
	code.Register(
		ErrRoleAssignedCode,
		"role already assigned to user: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrDuplicateRoleAssignmentCode,
		"duplicate role assignment: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrUserRoleLimitExceededCode,
		"user role limit exceeded: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrInvalidRoleAssignmentCode,
		"invalid role assignment: {msg}",
		code.WithAffectStability(false),
	)

	// 数据权限错误
	code.Register(
		ErrDataPermissionDeniedCode,
		"data permission denied: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrDataPermissionNotFoundCode,
		"data permission configuration not found: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrDataPermissionScopeInvalidCode,
		"invalid data permission scope: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrInvalidPermissionFilterCode,
		"invalid permission filter: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrCustomFilterSyntaxErrorCode,
		"custom filter syntax error: {msg}",
		code.WithAffectStability(false),
	)

	// 字段权限错误
	code.Register(
		ErrFieldPermissionDeniedCode,
		"field permission denied: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrFieldPermissionNotFoundCode,
		"field permission configuration not found: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrFieldPermissionInvalidCode,
		"invalid field permission: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrFieldNotFoundCode,
		"field not found: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrFieldPermissionLevelInvalidCode,
		"invalid field permission level: {msg}",
		code.WithAffectStability(false),
	)

	// 部门权限错误
	code.Register(
		ErrDepartmentNotFoundCode,
		"department not found: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrUserNotInDepartmentCode,
		"user is not in department: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrDepartmentAccessDeniedCode,
		"department access denied: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrInvalidDepartmentStructureCode,
		"invalid department structure: {msg}",
		code.WithAffectStability(false),
	)

	// 租户用户错误
	code.Register(
		ErrUserNotInTenantCode,
		"user is not in tenant: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrTenantUserNotFoundCode,
		"tenant user not found: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrTenantUserAlreadyExistsCode,
		"tenant user already exists: {msg}",
		code.WithAffectStability(false),
	)

	// 权限继承错误
	code.Register(
		ErrPermissionInheritanceCycleCode,
		"permission inheritance cycle detected: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrPermissionInheritanceDepthCode,
		"permission inheritance depth exceeded: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrBasePermissionNotFoundCode,
		"base permission not found: {msg}",
		code.WithAffectStability(false),
	)

	// 批量操作错误
	code.Register(
		ErrBatchOperationFailedCode,
		"batch operation failed: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrBatchOperationPartialFailureCode,
		"batch operation partial failure: {msg}",
		code.WithAffectStability(false),
	)
}

// Permission 便捷错误码变量（用于直接使用）
var (
	// 通用错误
	ErrPermissionInvalidParam = &BaseErrorCode{
		code:       "PERMISSION_INVALID_PARAM",
		message:    "Invalid parameter",
		messageZH:  "无效参数",
		messageEN:  "Invalid parameter",
		httpStatus: http.StatusBadRequest,
	}
	ErrPermissionCheckFailed = &BaseErrorCode{
		code:       "PERMISSION_CHECK_FAILED",
		message:    "Permission check failed",
		messageZH:  "权限检查失败",
		messageEN:  "Permission check failed",
		httpStatus: http.StatusInternalServerError,
	}
	// Note: ErrPermissionDenied is defined in botstore.go to avoid duplication
	ErrResourceTypeNotSupported = &BaseErrorCode{
		code:       "RESOURCE_TYPE_NOT_SUPPORTED",
		message:    "Resource type not supported",
		messageZH:  "不支持的资源类型",
		messageEN:  "Resource type not supported",
		httpStatus: http.StatusBadRequest,
	}

	// 角色相关错误
	ErrRoleNotFound = &BaseErrorCode{
		code:       "ROLE_NOT_FOUND",
		message:    "Role not found",
		messageZH:  "角色不存在",
		messageEN:  "Role not found",
		httpStatus: http.StatusNotFound,
	}
	ErrRoleAlreadyExists = &BaseErrorCode{
		code:       "ROLE_ALREADY_EXISTS",
		message:    "Role already exists",
		messageZH:  "角色已存在",
		messageEN:  "Role already exists",
		httpStatus: http.StatusConflict,
	}
	ErrInvalidRole = &BaseErrorCode{
		code:       "INVALID_ROLE",
		message:    "Invalid role",
		messageZH:  "无效角色",
		messageEN:  "Invalid role",
		httpStatus: http.StatusBadRequest,
	}
	ErrRoleNameInvalid = &BaseErrorCode{
		code:       "ROLE_NAME_INVALID",
		message:    "Invalid role name",
		messageZH:  "角色名称无效",
		messageEN:  "Invalid role name",
		httpStatus: http.StatusBadRequest,
	}
	ErrRoleDescriptionInvalid = &BaseErrorCode{
		code:       "ROLE_DESCRIPTION_INVALID",
		message:    "Invalid role description",
		messageZH:  "角色描述无效",
		messageEN:  "Invalid role description",
		httpStatus: http.StatusBadRequest,
	}
	ErrCannotDeleteSystemRole = &BaseErrorCode{
		code:       "CANNOT_DELETE_SYSTEM_ROLE",
		message:    "Cannot delete system role",
		messageZH:  "不能删除系统角色",
		messageEN:  "Cannot delete system role",
		httpStatus: http.StatusForbidden,
	}
	ErrRoleInUse = &BaseErrorCode{
		code:       "ROLE_IN_USE",
		message:    "Role is in use",
		messageZH:  "角色使用中",
		messageEN:  "Role is in use",
		httpStatus: http.StatusConflict,
	}
	ErrCircularRoleDependency = &BaseErrorCode{
		code:       "CIRCULAR_ROLE_DEPENDENCY",
		message:    "Circular role dependency detected",
		messageZH:  "循环角色依赖",
		messageEN:  "Circular role dependency detected",
		httpStatus: http.StatusBadRequest,
	}

	// 角色分配错误
	ErrRoleAssigned = &BaseErrorCode{
		code:       "ROLE_ASSIGNED",
		message:    "Role already assigned",
		messageZH:  "角色已分配",
		messageEN:  "Role already assigned",
		httpStatus: http.StatusConflict,
	}
	ErrDuplicateRoleAssignment = &BaseErrorCode{
		code:       "DUPLICATE_ROLE_ASSIGNMENT",
		message:    "Duplicate role assignment",
		messageZH:  "重复角色分配",
		messageEN:  "Duplicate role assignment",
		httpStatus: http.StatusConflict,
	}
	ErrUserRoleLimitExceeded = &BaseErrorCode{
		code:       "USER_ROLE_LIMIT_EXCEEDED",
		message:    "User role limit exceeded",
		messageZH:  "用户角色数量超限",
		messageEN:  "User role limit exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrInvalidRoleAssignment = &BaseErrorCode{
		code:       "INVALID_ROLE_ASSIGNMENT",
		message:    "Invalid role assignment",
		messageZH:  "无效角色分配",
		messageEN:  "Invalid role assignment",
		httpStatus: http.StatusBadRequest,
	}

	// 数据权限错误
	ErrDataPermissionDenied = &BaseErrorCode{
		code:       "DATA_PERMISSION_DENIED",
		message:    "Data permission denied",
		messageZH:  "数据权限拒绝",
		messageEN:  "Data permission denied",
		httpStatus: http.StatusForbidden,
	}
	ErrDataPermissionNotFound = &BaseErrorCode{
		code:       "DATA_PERMISSION_NOT_FOUND",
		message:    "Data permission configuration not found",
		messageZH:  "数据权限配置不存在",
		messageEN:  "Data permission configuration not found",
		httpStatus: http.StatusNotFound,
	}
	ErrDataPermissionScopeInvalid = &BaseErrorCode{
		code:       "DATA_PERMISSION_SCOPE_INVALID",
		message:    "Invalid data permission scope",
		messageZH:  "数据权限范围无效",
		messageEN:  "Invalid data permission scope",
		httpStatus: http.StatusBadRequest,
	}
	ErrInvalidPermissionFilter = &BaseErrorCode{
		code:       "INVALID_PERMISSION_FILTER",
		message:    "Invalid permission filter",
		messageZH:  "无效权限过滤器",
		messageEN:  "Invalid permission filter",
		httpStatus: http.StatusBadRequest,
	}
	ErrCustomFilterSyntaxError = &BaseErrorCode{
		code:       "CUSTOM_FILTER_SYNTAX_ERROR",
		message:    "Custom filter syntax error",
		messageZH:  "自定义过滤器语法错误",
		messageEN:  "Custom filter syntax error",
		httpStatus: http.StatusBadRequest,
	}

	// 字段权限错误
	ErrFieldPermissionDenied = &BaseErrorCode{
		code:       "FIELD_PERMISSION_DENIED",
		message:    "Field permission denied",
		messageZH:  "字段权限拒绝",
		messageEN:  "Field permission denied",
		httpStatus: http.StatusForbidden,
	}
	ErrFieldPermissionNotFound = &BaseErrorCode{
		code:       "FIELD_PERMISSION_NOT_FOUND",
		message:    "Field permission configuration not found",
		messageZH:  "字段权限配置不存在",
		messageEN:  "Field permission configuration not found",
		httpStatus: http.StatusNotFound,
	}
	ErrFieldPermissionInvalid = &BaseErrorCode{
		code:       "FIELD_PERMISSION_INVALID",
		message:    "Invalid field permission",
		messageZH:  "字段权限无效",
		messageEN:  "Invalid field permission",
		httpStatus: http.StatusBadRequest,
	}
	ErrFieldNotFound = &BaseErrorCode{
		code:       "FIELD_NOT_FOUND",
		message:    "Field not found",
		messageZH:  "字段不存在",
		messageEN:  "Field not found",
		httpStatus: http.StatusNotFound,
	}
	ErrFieldPermissionLevelInvalid = &BaseErrorCode{
		code:       "FIELD_PERMISSION_LEVEL_INVALID",
		message:    "Invalid field permission level",
		messageZH:  "字段权限级别无效",
		messageEN:  "Invalid field permission level",
		httpStatus: http.StatusBadRequest,
	}

	// 部门权限错误
	ErrDepartmentNotFound = &BaseErrorCode{
		code:       "DEPARTMENT_NOT_FOUND",
		message:    "Department not found",
		messageZH:  "部门不存在",
		messageEN:  "Department not found",
		httpStatus: http.StatusNotFound,
	}
	ErrUserNotInDepartment = &BaseErrorCode{
		code:       "USER_NOT_IN_DEPARTMENT",
		message:    "User is not in department",
		messageZH:  "用户不在部门中",
		messageEN:  "User is not in department",
		httpStatus: http.StatusBadRequest,
	}
	ErrDepartmentAccessDenied = &BaseErrorCode{
		code:       "DEPARTMENT_ACCESS_DENIED",
		message:    "Department access denied",
		messageZH:  "部门访问拒绝",
		messageEN:  "Department access denied",
		httpStatus: http.StatusForbidden,
	}
	ErrInvalidDepartmentStructure = &BaseErrorCode{
		code:       "INVALID_DEPARTMENT_STRUCTURE",
		message:    "Invalid department structure",
		messageZH:  "无效部门结构",
		messageEN:  "Invalid department structure",
		httpStatus: http.StatusBadRequest,
	}

	// 租户用户错误
	ErrUserNotInTenant = &BaseErrorCode{
		code:       "USER_NOT_IN_TENANT",
		message:    "User is not in tenant",
		messageZH:  "用户不在租户中",
		messageEN:  "User is not in tenant",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantUserNotFound = &BaseErrorCode{
		code:       "TENANT_USER_NOT_FOUND",
		message:    "Tenant user not found",
		messageZH:  "租户用户不存在",
		messageEN:  "Tenant user not found",
		httpStatus: http.StatusNotFound,
	}
	ErrTenantUserAlreadyExists = &BaseErrorCode{
		code:       "TENANT_USER_ALREADY_EXISTS",
		message:    "Tenant user already exists",
		messageZH:  "租户用户已存在",
		messageEN:  "Tenant user already exists",
		httpStatus: http.StatusConflict,
	}

	// 权限继承错误
	ErrPermissionInheritanceCycle = &BaseErrorCode{
		code:       "PERMISSION_INHERITANCE_CYCLE",
		message:    "Permission inheritance cycle detected",
		messageZH:  "权限继承循环",
		messageEN:  "Permission inheritance cycle detected",
		httpStatus: http.StatusBadRequest,
	}
	ErrPermissionInheritanceDepth = &BaseErrorCode{
		code:       "PERMISSION_INHERITANCE_DEPTH",
		message:    "Permission inheritance depth exceeded",
		messageZH:  "权限继承深度超限",
		messageEN:  "Permission inheritance depth exceeded",
		httpStatus: http.StatusBadRequest,
	}
	ErrBasePermissionNotFound = &BaseErrorCode{
		code:       "BASE_PERMISSION_NOT_FOUND",
		message:    "Base permission not found",
		messageZH:  "基础权限不存在",
		messageEN:  "Base permission not found",
		httpStatus: http.StatusNotFound,
	}

	// 批量操作错误
	ErrBatchOperationFailed = &BaseErrorCode{
		code:       "BATCH_OPERATION_FAILED",
		message:    "Batch operation failed",
		messageZH:  "批量操作失败",
		messageEN:  "Batch operation failed",
		httpStatus: http.StatusInternalServerError,
	}
	ErrBatchOperationPartialFailure = &BaseErrorCode{
		code:       "BATCH_OPERATION_PARTIAL_FAILURE",
		message:    "Batch operation partial failure",
		messageZH:  "批量操作部分失败",
		messageEN:  "Batch operation partial failure",
		httpStatus: http.StatusMultiStatus,
	}
)
