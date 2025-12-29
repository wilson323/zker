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
	"github.com/coze-dev/coze-studio/backend/pkg/errorx/code"
)

// Tenant: 200 000 000 ~ 200 999 999
const (
	// 租户基础错误 (200 000 000 ~ 200 009 999)
	ErrTenantNotFoundCode            = 200000001 // 租户不存在
	ErrTenantAlreadyExistsCode       = 200000002 // 租户已存在
	ErrTenantSuspendedCode           = 200000003 // 租户已暂停
	ErrTenantDeletedCode             = 200000004 // 租户已删除
	ErrTenantInvalidParamCode        = 200000005 // 租户参数无效
	ErrTenantNameTooLongCode         = 200000006 // 租户名称过长
	ErrTenantTypeInvalidCode         = 200000007 // 租户类型无效
	ErrTenantStatusInvalidCode       = 200000008 // 租户状态无效
	ErrTenantOperationNotAllowedCode = 200000009 // 操作不允许
	ErrInvalidTenantIDCode           = 200000010 // 租户ID无效

	// 租户配额相关 (200 010 000 ~ 200 019 999)
	ErrTenantQuotaExceededCode     = 200010001 // 租户配额已超限
	ErrTenantQuotaInvalidCode      = 200010002 // 配额无效
	ErrTenantQuotaCheckFailedCode  = 200010003 // 配额检查失败
	ErrTenantQuotaResetFailedCode  = 200010004 // 配额重置失败
	ErrTenantQuotaNotFoundCode     = 200010005 // 配额不存在
	ErrTenantQuotaUpdateFailedCode = 200010006 // 配额更新失败

	// 租户订阅相关 (200 020 000 ~ 200 029 999)
	ErrTenantSubscriptionNotFoundCode     = 200020001 // 订阅不存在
	ErrTenantSubscriptionExpiredCode      = 200020002 // 订阅已过期
	ErrTenantSubscriptionInvalidCode      = 200020003 // 订阅无效
	ErrTenantSubscriptionUpgradeFailedCode = 200020004 // 订阅升级失败
	ErrTenantSubscriptionDowngradeFailedCode = 200020005 // 订阅降级失败
	ErrTenantSubscriptionCancelFailedCode = 200020006 // 订阅取消失败
	ErrTenantSubscriptionActiveCode       = 200020007 // 订阅已激活
)

// 便捷错误变量（直接使用，避免重复创建）
var (
	ErrTenantNotFound      = errorx.New(ErrTenantNotFoundCode)
	ErrTenantAlreadyExists = errorx.New(ErrTenantAlreadyExistsCode)
	ErrTenantSuspended     = errorx.New(ErrTenantSuspendedCode)
	ErrTenantDeleted       = errorx.New(ErrTenantDeletedCode)
	ErrInvalidTenantID     = errorx.New(ErrInvalidTenantIDCode)
)

func init() {
	// 租户基础错误注册
	code.Register(
		ErrTenantNotFoundCode,
		"Tenant not found: {tenant_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantAlreadyExistsCode,
		"Tenant already exists: {tenant_name}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantSuspendedCode,
		"Tenant is suspended: {tenant_id}, reason: {reason}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTenantDeletedCode,
		"Tenant is deleted: {tenant_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantInvalidParamCode,
		"Invalid tenant parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantNameTooLongCode,
		"Tenant name too long: max {max_length} characters, got {actual_length}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantTypeInvalidCode,
		"Invalid tenant type: {tenant_type}, must be one of: individual, team, enterprise",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantStatusInvalidCode,
		"Invalid tenant status: {status}, must be one of: active, suspended, deleted",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantOperationNotAllowedCode,
		"Operation not allowed for tenant: {tenant_id}, operation: {operation}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrInvalidTenantIDCode,
		"Invalid tenant_id: {tenant_id}, reason: {reason}",
		code.WithAffectStability(false),
	)

	// 租户配额错误注册
	code.Register(
		ErrTenantQuotaExceededCode,
		"Tenant quota exceeded: tenant_id={tenant_id}, resource_type={resource_type}, current={current_usage}, limit={quota_limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTenantQuotaInvalidCode,
		"Invalid quota: tenant_id={tenant_id}, resource_type={resource_type}, reason: {reason}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantQuotaCheckFailedCode,
		"Failed to check quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTenantQuotaResetFailedCode,
		"Failed to reset quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTenantQuotaNotFoundCode,
		"Quota not found: tenant_id={tenant_id}, resource_type={resource_type}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantQuotaUpdateFailedCode,
		"Failed to update quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	// 租户订阅错误注册
	code.Register(
		ErrTenantSubscriptionNotFoundCode,
		"Subscription not found: tenant_id={tenant_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantSubscriptionExpiredCode,
		"Subscription expired: tenant_id={tenant_id}, expired_at={expired_at}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTenantSubscriptionInvalidCode,
		"Invalid subscription: tenant_id={tenant_id}, tier={tier}, reason: {reason}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTenantSubscriptionUpgradeFailedCode,
		"Failed to upgrade subscription: tenant_id={tenant_id}, from_tier={from_tier}, to_tier={to_tier}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTenantSubscriptionDowngradeFailedCode,
		"Failed to downgrade subscription: tenant_id={tenant_id}, from_tier={from_tier}, to_tier={to_tier}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTenantSubscriptionCancelFailedCode,
		"Failed to cancel subscription: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTenantSubscriptionActiveCode,
		"Subscription already active: tenant_id={tenant_id}, tier={tier}",
		code.WithAffectStability(false),
	)
}
