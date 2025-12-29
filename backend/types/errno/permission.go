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

import "github.com/coze-dev/coze-studio/backend/pkg/errorx/code"

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
