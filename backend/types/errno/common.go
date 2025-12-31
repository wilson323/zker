// backend/types/errno/common.go
package errno

import (
	"net/http"
)

// ErrorCode 错误码接口
type ErrorCode interface {
	Code() string
	Message() string
	MessageZH() string
	MessageEN() string
	HTTPStatus() int
}

// BaseErrorCode 基础错误码实现
type BaseErrorCode struct {
	code       string
	message    string
	messageZH  string
	messageEN  string
	httpStatus int
}

func (e *BaseErrorCode) Code() string       { return e.code }
func (e *BaseErrorCode) Message() string    { return e.message }
func (e *BaseErrorCode) MessageZH() string  { return e.messageZH }
func (e *BaseErrorCode) MessageEN() string  { return e.messageEN }
func (e *BaseErrorCode) HTTPStatus() int    { return e.httpStatus }

// Int32Code 返回int32格式的错误码
// 这个方法将字符串错误码（如"CACHE500001"）转换为int32
// 注意：这里简单地返回HTTP状态码，实际应用中需要更复杂的映射逻辑
func (e *BaseErrorCode) Int32Code() int32 {
	// 将HTTP状态码作为默认错误码
	return int32(e.httpStatus)
}

// 通用错误码（20个）
var (
	Success = &BaseErrorCode{
		code:       "SUCCESS",
		message:    "Operation successful",
		messageZH:  "操作成功",
		messageEN:  "Operation successful",
		httpStatus: http.StatusOK,
	}

	// 客户端错误 (2xx)
	InvalidParams = &BaseErrorCode{
		code:       "COMMON201001",
		message:    "Invalid parameters",
		messageZH:  "参数错误",
		messageEN:  "Invalid parameters",
		httpStatus: http.StatusBadRequest,
	}
	MissingParam = &BaseErrorCode{
		code:       "COMMON201002",
		message:    "Missing required parameter",
		messageZH:  "缺少必要参数",
		messageEN:  "Missing required parameter",
		httpStatus: http.StatusBadRequest,
	}
	InvalidFormat = &BaseErrorCode{
		code:       "COMMON201003",
		message:    "Invalid format",
		messageZH:  "格式错误",
		messageEN:  "Invalid format",
		httpStatus: http.StatusBadRequest,
	}

	// 认证错误 (4xx)
	Unauthorized = &BaseErrorCode{
		code:       "COMMON401001",
		message:    "Unauthorized",
		messageZH:  "未认证",
		messageEN:  "Unauthorized",
		httpStatus: http.StatusUnauthorized,
	}
	TokenExpired = &BaseErrorCode{
		code:       "COMMON401002",
		message:    "Token expired",
		messageZH:  "Token已过期",
		messageEN:  "Token expired",
		httpStatus: http.StatusUnauthorized,
	}
	InvalidToken = &BaseErrorCode{
		code:       "COMMON401003",
		message:    "Invalid token",
		messageZH:  "无效的Token",
		messageEN:  "Invalid token",
		httpStatus: http.StatusUnauthorized,
	}

	// 权限错误 (4xx)
	Forbidden = &BaseErrorCode{
		code:       "COMMON403001",
		message:    "Permission denied",
		messageZH:  "权限不足",
		messageEN:  "Permission denied",
		httpStatus: http.StatusForbidden,
	}

	// 资源错误 (4xx)
	NotFound = &BaseErrorCode{
		code:       "COMMON404001",
		message:    "Resource not found",
		messageZH:  "资源不存在",
		messageEN:  "Resource not found",
		httpStatus: http.StatusNotFound,
	}

	// 服务器错误 (5xx)
	InternalError = &BaseErrorCode{
		code:       "COMMON500001",
		message:    "Internal server error",
		messageZH:  "服务器内部错误",
		messageEN:  "Internal server error",
		httpStatus: http.StatusInternalServerError,
	}
	DatabaseError = &BaseErrorCode{
		code:       "COMMON500002",
		message:    "Database error",
		messageZH:  "数据库错误",
		messageEN:  "Database error",
		httpStatus: http.StatusInternalServerError,
	}
	IDGenError = &BaseErrorCode{
		code:       "COMMON500003",
		message:    "Failed to generate ID",
		messageZH:  "ID生成失败",
		messageEN:  "Failed to generate ID",
		httpStatus: http.StatusInternalServerError,
	}

	// 缓存错误
	CacheError = &BaseErrorCode{
		code:       "COMMON500004",
		message:    "Cache operation failed",
		messageZH:  "缓存操作失败",
		messageEN:  "Cache operation failed",
		httpStatus: http.StatusInternalServerError,
	}
)

// 通用错误码int32常量（用于errorx.New等需要int32的场景）
const (
	// 成功
	SuccessCode = 200000000

	// 客户端错误 (2xx)
	InvalidParamsCode = 201000001
	MissingParamCode  = 201000002
	InvalidFormatCode = 201000003

	// 认证错误 (4xx)
	UnauthorizedCode = 401000001
	TokenExpiredCode  = 401000002
	InvalidTokenCode  = 401000003

	// 权限错误 (4xx)
	ForbiddenCode = 403000001

	// 资源错误 (4xx)
	NotFoundCode = 404000001

	// 服务器错误 (5xx)
	InternalErrorCode = 500000001
	DatabaseErrorCode = 500000002
	IDGenErrorCode    = 500000003
	CacheErrorCode    = 500000004
	ErrRateLimitExceededCode       = 429000001
	ErrConcurrencyLimitExceededCode = 429000002
)

// 便捷别名（不带Err前缀，用于直接使用int32常量）
const (
	SuccessValue      = SuccessCode
	InvalidParam      = InvalidParamsCode
	ParamMissing      = MissingParamCode
	FormatInvalid     = InvalidFormatCode
	AuthRequired      = UnauthorizedCode
	TokenExpiredValue = TokenExpiredCode
	TokenInvalid      = InvalidTokenCode
	AccessDenied      = ForbiddenCode
	ResourceNotFound  = NotFoundCode
	InternalErr       = InternalErrorCode
	DatabaseErr       = DatabaseErrorCode
)

// Err前缀别名（用于中间件统一命名）
const (
	ErrUnauthorizedCode         = UnauthorizedCode
	ErrForbiddenCode            = ForbiddenCode
	ErrNotFoundCode             = NotFoundCode
	ErrInvalidParamsCode        = InvalidParamsCode
	ErrInternalErrorCode        = InternalErrorCode
)
