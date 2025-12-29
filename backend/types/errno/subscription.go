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

// Subscription: 400 000 000 ~ 400 999 999
const (
	// 订阅基础错误 (400 000 000 ~ 400 009 999)
	ErrSubscriptionNotFoundCode        = 400000001 // 订阅不存在
	ErrSubscriptionAlreadyExistsCode   = 400000002 // 订阅已存在
	ErrSubscriptionInvalidCode         = 400000003 // 订阅无效
	ErrSubscriptionExpiredCode         = 400000004 // 订阅已过期
	ErrSubscriptionInactiveCode        = 400000005 // 订阅未激活
	ErrSubscriptionSuspendedCode       = 400000006 // 订阅已暂停
	ErrSubscriptionCreateFailedCode    = 400000007 // 订阅创建失败
	ErrSubscriptionUpdateFailedCode    = 400000008 // 订阅更新失败
	ErrSubscriptionCancelFailedCode    = 400000009 // 订阅取消失败
	ErrSubscriptionDeleteFailedCode    = 400000010 // 订阅删除失败
	ErrSubscriptionInvalidTierCode     = 400000011 // 无效的订阅等级
	ErrSubscriptionInvalidStatusCode   = 400000012 // 无效的订阅状态

	// 订阅等级错误 (400 010 000 ~ 400 019 999)
	ErrSubscriptionTierUpgradeFailedCode  = 400010001 // 订阅升级失败
	ErrSubscriptionTierDowngradeFailedCode = 400010002 // 订阅降级失败
	ErrSubscriptionTierInvalidCode        = 400010003 // 无效的订阅等级
	ErrSubscriptionTierNotAllowedCode     = 400010004 // 订阅等级不允许
	ErrSubscriptionTierTransitionFailedCode = 400010005 // 订阅等级转换失败

	// 订阅计费错误 (400 020 000 ~ 400 029 999)
	ErrSubscriptionBillingFailedCode         = 400020001 // 计费失败
	ErrSubscriptionPaymentFailedCode         = 400020002 // 支付失败
	ErrSubscriptionPaymentRequiredCode       = 400020003 // 需要支付
	ErrSubscriptionInvoiceNotFoundCode       = 400020004 // 发票不存在
	ErrSubscriptionInvoiceGenerationFailedCode = 400020005 // 发票生成失败
	ErrSubscriptionRefundFailedCode          = 400020006 // 退款失败
	ErrSubscriptionBillingCycleInvalidCode   = 400020007 // 无效的计费周期

	// 订阅配额错误 (400 030 000 ~ 400 039 999)
	ErrSubscriptionQuotaInsufficientCode = 400030001 // 订阅配额不足
	ErrSubscriptionQuotaExceededCode     = 400030002 // 订阅配额超限
	ErrSubscriptionQuotaInvalidCode      = 400030003 // 订阅配额无效
	ErrSubscriptionQuotaUpdateFailedCode = 400030004 // 订阅配额更新失败

	// 订阅试用错误 (400 040 000 ~ 400 049 999)
	ErrSubscriptionTrialExpiredCode     = 400040001 // 试用期已过期
	ErrSubscriptionTrialNotAvailableCode = 400040002 // 试用期不可用
	ErrSubscriptionTrialAlreadyUsedCode = 400040003 // 试用期已使用
	ErrSubscriptionTrialStartFailedCode = 400040004 // 试用开始失败

	// 订阅续费错误 (400 050 000 ~ 400 059 999)
	ErrSubscriptionRenewalFailedCode     = 400050001 // 续费失败
	ErrSubscriptionRenewalNotAllowedCode = 400050002 // 不允许续费
	ErrSubscriptionAutoRenewalFailedCode = 400050003 // 自动续费失败
	ErrSubscriptionAutoRenewalDisabledCode = 400050004 // 自动续费已禁用
)

func init() {
	// 订阅基础错误注册
	code.Register(
		ErrSubscriptionNotFoundCode,
		"Subscription not found: tenant_id={tenant_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionAlreadyExistsCode,
		"Subscription already exists: tenant_id={tenant_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionInvalidCode,
		"Invalid subscription: tenant_id={tenant_id}, reason: {reason}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionExpiredCode,
		"Subscription expired: tenant_id={tenant_id}, expired_at={expired_at}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionInactiveCode,
		"Subscription inactive: tenant_id={tenant_id}, status={status}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionSuspendedCode,
		"Subscription suspended: tenant_id={tenant_id}, reason: {reason}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionCreateFailedCode,
		"Failed to create subscription: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionUpdateFailedCode,
		"Failed to update subscription: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionCancelFailedCode,
		"Failed to cancel subscription: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionDeleteFailedCode,
		"Failed to delete subscription: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionInvalidTierCode,
		"Invalid subscription tier: {tier}, must be one of: free, pro, enterprise",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionInvalidStatusCode,
		"Invalid subscription status: {status}, must be one of: active, inactive, suspended, expired",
		code.WithAffectStability(false),
	)

	// 订阅等级错误注册
	code.Register(
		ErrSubscriptionTierUpgradeFailedCode,
		"Failed to upgrade subscription tier: tenant_id={tenant_id}, from_tier={from_tier}, to_tier={to_tier}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionTierDowngradeFailedCode,
		"Failed to downgrade subscription tier: tenant_id={tenant_id}, from_tier={from_tier}, to_tier={to_tier}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionTierInvalidCode,
		"Invalid subscription tier: {tier}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionTierNotAllowedCode,
		"Subscription tier not allowed: tenant_id={tenant_id}, tier={tier}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionTierTransitionFailedCode,
		"Failed to transition subscription tier: tenant_id={tenant_id}, from_tier={from_tier}, to_tier={to_tier}, error: {error}",
		code.WithAffectStability(true),
	)

	// 订阅计费错误注册
	code.Register(
		ErrSubscriptionBillingFailedCode,
		"Billing failed: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionPaymentFailedCode,
		"Payment failed: tenant_id={tenant_id}, payment_id={payment_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionPaymentRequiredCode,
		"Payment required: tenant_id={tenant_id}, amount={amount}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionInvoiceNotFoundCode,
		"Invoice not found: tenant_id={tenant_id}, invoice_id={invoice_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionInvoiceGenerationFailedCode,
		"Failed to generate invoice: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionRefundFailedCode,
		"Refund failed: tenant_id={tenant_id}, payment_id={payment_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionBillingCycleInvalidCode,
		"Invalid billing cycle: {cycle}, must be one of: monthly, yearly",
		code.WithAffectStability(false),
	)

	// 订阅配额错误注册
	code.Register(
		ErrSubscriptionQuotaInsufficientCode,
		"Subscription quota insufficient: tenant_id={tenant_id}, resource_type={resource_type}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionQuotaExceededCode,
		"Subscription quota exceeded: tenant_id={tenant_id}, resource_type={resource_type}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionQuotaInvalidCode,
		"Invalid subscription quota: tenant_id={tenant_id}, resource_type={resource_type}, reason: {reason}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionQuotaUpdateFailedCode,
		"Failed to update subscription quota: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// 订阅试用错误注册
	code.Register(
		ErrSubscriptionTrialExpiredCode,
		"Trial period expired: tenant_id={tenant_id}, expired_at={expired_at}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionTrialNotAvailableCode,
		"Trial period not available: tenant_id={tenant_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionTrialAlreadyUsedCode,
		"Trial period already used: tenant_id={tenant_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionTrialStartFailedCode,
		"Failed to start trial period: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// 订阅续费错误注册
	code.Register(
		ErrSubscriptionRenewalFailedCode,
		"Failed to renew subscription: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionRenewalNotAllowedCode,
		"Subscription renewal not allowed: tenant_id={tenant_id}, reason: {reason}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSubscriptionAutoRenewalFailedCode,
		"Auto-renewal failed: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSubscriptionAutoRenewalDisabledCode,
		"Auto-renewal disabled: tenant_id={tenant_id}",
		code.WithAffectStability(false),
	)
}
