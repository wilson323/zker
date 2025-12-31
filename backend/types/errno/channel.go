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

// Channel: 205 000 000 ~ 205 999 999
// 渠道适配器相关错误码
const (
	// 通用错误 (205 000 000 ~ 205 009 999)
	ErrChannelInvalidParamCode       = 205000001 // 无效参数
	ErrChannelConfigInvalidCode      = 205000002 // 配置无效
	ErrChannelNotFoundCode           = 205000003 // 渠道不存在
	ErrChannelAlreadyExistsCode      = 205000004 // 渠道已存在
	ErrChannelPublishFailedCode      = 205000005 // 发布失败
	ErrChannelUnpublishFailedCode    = 205000006 // 取消发布失败

	// 消息处理错误 (205 010 000 ~ 205 019 999)
	ErrChannelMessageParseFailed     = 205010001 // 消息解析失败
	ErrChannelMessageTypeNotSupported = 205010002 // 不支持的消息类型
	ErrChannelMessageSendFailed      = 205010003 // 消息发送失败
	ErrChannelMessageTooLarge        = 205010004 // 消息过大

	// 签名验证错误 (205 020 000 ~ 205 029 999)
	ErrChannelSignatureInvalid       = 205020001 // 签名无效
	ErrChannelTimestampExpired       = 205020002 // 时间戳过期
	ErrChannelNonceReused            = 205020003 // Nonce重复使用

	// Webhook错误 (205 030 000 ~ 205 039 999)
	ErrChannelWebhookFailed          = 205030001 // Webhook处理失败
	ErrChannelWebhookTimeout         = 205030002 // Webhook超时

	// 第三方API错误 (205 040 000 ~ 205 049 999)
	ErrChannelAPIRequestFailed       = 205040001 // API请求失败
	ErrChannelAPIRateLimitExceeded   = 205040002 // API频率限制
	ErrChannelAPITokenExpired        = 205040003 // API令牌过期

	// 微信公众号专用错误 (205 100 000 ~ 205 109 999)
	ErrChannelWechatAppIDInvalid     = 205100001 // AppID无效
	ErrChannelWechatAppSecretInvalid = 205100002 // AppSecret无效
	ErrChannelWechatAccessTokenFailed = 205100003 // 获取AccessToken失败
	ErrChannelWechatMenuCreateFailed = 205100004 // 创建菜单失败
	ErrChannelWechatMediaUploadFailed = 205100005 // 上传媒体文件失败
)

func init() {
	// 通用错误
	code.Register(
		ErrChannelInvalidParamCode,
		"channel invalid parameter: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelConfigInvalidCode,
		"channel config invalid: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelNotFoundCode,
		"channel not found: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelAlreadyExistsCode,
		"channel already exists: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelPublishFailedCode,
		"channel publish failed: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelUnpublishFailedCode,
		"channel unpublish failed: {msg}",
		code.WithAffectStability(false),
	)

	// 消息处理错误
	code.Register(
		ErrChannelMessageParseFailed,
		"channel message parse failed: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelMessageTypeNotSupported,
		"channel message type not supported: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelMessageSendFailed,
		"channel message send failed: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelMessageTooLarge,
		"channel message too large: {msg}",
		code.WithAffectStability(false),
	)

	// 签名验证错误
	code.Register(
		ErrChannelSignatureInvalid,
		"channel signature invalid: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelTimestampExpired,
		"channel timestamp expired: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelNonceReused,
		"channel nonce reused: {msg}",
		code.WithAffectStability(false),
	)

	// Webhook错误
	code.Register(
		ErrChannelWebhookFailed,
		"channel webhook failed: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelWebhookTimeout,
		"channel webhook timeout: {msg}",
		code.WithAffectStability(false),
	)

	// 第三方API错误
	code.Register(
		ErrChannelAPIRequestFailed,
		"channel API request failed: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelAPIRateLimitExceeded,
		"channel API rate limit exceeded: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelAPITokenExpired,
		"channel API token expired: {msg}",
		code.WithAffectStability(false),
	)

	// 微信公众号专用错误
	code.Register(
		ErrChannelWechatAppIDInvalid,
		"wechat app id invalid: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelWechatAppSecretInvalid,
		"wechat app secret invalid: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelWechatAccessTokenFailed,
		"wechat get access token failed: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelWechatMenuCreateFailed,
		"wechat create menu failed: {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrChannelWechatMediaUploadFailed,
		"wechat upload media failed: {msg}",
		code.WithAffectStability(false),
	)
}
