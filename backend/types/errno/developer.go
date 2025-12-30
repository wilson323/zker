// backend/types/errno/developer.go
package errno

import (
	"net/http"
)

// 开发者平台相关错误码（DEV 开头）
var (
	// 项目管理错误（DEV2xx）
	ErrProjectNotFound      = &BaseErrorCode{"DEV201001", "项目不存在", "Project not found", "Project not found", http.StatusNotFound}
	ErrProjectNameExists    = &BaseErrorCode{"DEV201002", "项目名称已存在", "项目名称已存在", "Project name already exists", http.StatusConflict}
	ErrProjectTypeInvalid   = &BaseErrorCode{"DEV201003", "无效的项目类型", "无效的项目类型", "Invalid project type", http.StatusBadRequest}
	ErrProjectStatusInvalid = &BaseErrorCode{"DEV201004", "无效的项目状态", "无效的项目状态", "Invalid project status", http.StatusBadRequest}
	ErrProjectLimitExceeded = &BaseErrorCode{"DEV201005", "项目数量超限", "项目数量已达上限", "Project limit exceeded", http.StatusForbidden}
	ErrProjectInProduction  = &BaseErrorCode{"DEV201006", "项目正在生产环境", "无法修改生产环境的项目", "Project is in production", http.StatusForbidden}

	// API密钥管理错误（DEV4xx）
	ErrAPIKeyNotFound          = &BaseErrorCode{"DEV401001", "API密钥不存在", "API密钥不存在", "API key not found", http.StatusNotFound}
	ErrAPIKeyInvalid           = &BaseErrorCode{"DEV401002", "无效的API密钥", "无效的API密钥", "Invalid API key", http.StatusUnauthorized}
	ErrAPIKeyExpired           = &BaseErrorCode{"DEV401003", "API密钥已过期", "API密钥已过期", "API key has expired", http.StatusUnauthorized}
	ErrAPIKeyRevoked           = &BaseErrorCode{"DEV401004", "API密钥已撤销", "API密钥已撤销", "API key has been revoked", http.StatusUnauthorized}
	ErrAPIKeyLimitExceeded     = &BaseErrorCode{"DEV401005", "API密钥数量超限", "API密钥数量已达上限", "API key limit exceeded", http.StatusForbidden}
	ErrAPIKeyPermissionDenied  = &BaseErrorCode{"DEV401006", "API密钥权限不足", "API密钥权限不足", "API key permission denied", http.StatusForbidden}
	ErrAPIKeyNameInvalid       = &BaseErrorCode{"DEV401007", "无效的密钥名称", "无效的密钥名称", "Invalid key name", http.StatusBadRequest}
	ErrAPIKeyScopeInvalid      = &BaseErrorCode{"DEV401008", "无效的密钥权限", "无效的密钥权限", "Invalid key scopes", http.StatusBadRequest}

	// Webhook错误（DEV4xx）
	ErrWebhookNotFound      = &BaseErrorCode{"DEV402001", "Webhook不存在", "Webhook不存在", "Webhook not found", http.StatusNotFound}
	ErrWebhookURLInvalid    = &BaseErrorCode{"DEV402002", "无效的Webhook URL", "无效的Webhook URL", "Invalid webhook URL", http.StatusBadRequest}
	ErrWebhookSignatureInvalid = &BaseErrorCode{"DEV402003", "Webhook签名验证失败", "Webhook签名验证失败", "Invalid webhook signature", http.StatusUnauthorized}
	ErrWebhookEventInvalid  = &BaseErrorCode{"DEV402004", "无效的Webhook事件", "无效的Webhook事件", "Invalid webhook event", http.StatusBadRequest}
	ErrWebhookLimitExceeded = &BaseErrorCode{"DEV402005", "Webhook数量超限", "Webhook数量已达上限", "Webhook limit exceeded", http.StatusForbidden}
	ErrWebhookTriggerFailed = &BaseErrorCode{"DEV402006", "Webhook触发失败", "Webhook触发失败", "Failed to trigger webhook", http.StatusInternalServerError}

	// SDK错误（DEV4xx）
	ErrSDKNotFound         = &BaseErrorCode{"DEV403001", "SDK不存在", "SDK不存在", "SDK not found", http.StatusNotFound}
	ErrSDKLanguageInvalid  = &BaseErrorCode{"DEV403002", "不支持的SDK语言", "不支持的SDK语言", "Unsupported SDK language", http.StatusBadRequest}
	ErrSDKVersionInvalid   = &BaseErrorCode{"DEV403003", "无效的SDK版本", "无效的SDK版本", "Invalid SDK version", http.StatusBadRequest}
	ErrSDKGenerationFailed = &BaseErrorCode{"DEV403004", "SDK生成失败", "SDK生成失败", "Failed to generate SDK", http.StatusInternalServerError}
	ErrSDKDownloadFailed   = &BaseErrorCode{"DEV403005", "SDK下载失败", "SDK下载失败", "Failed to download SDK", http.StatusInternalServerError}

	// 开发者错误（DEV4xx）
	ErrDeveloperNotFound      = &BaseErrorCode{"DEV404001", "开发者不存在", "开发者不存在", "Developer not found", http.StatusNotFound}
	ErrDeveloperEmailInvalid  = &BaseErrorCode{"DEV404002", "无效的开发者邮箱", "无效的开发者邮箱", "Invalid developer email", http.StatusBadRequest}
	ErrDeveloperExists        = &BaseErrorCode{"DEV404003", "开发者已存在", "开发者已存在", "Developer already exists", http.StatusConflict}
	ErrDeveloperSuspended     = &BaseErrorCode{"DEV404004", "开发者账户已暂停", "开发者账户已暂停", "Developer account suspended", http.StatusForbidden}
)
