// backend/types/errno/user.go
package errno

import (
	"net/http"
)

// 用户相关错误码常量
const (
	ErrUserAuthenticationFailed  = 7000001
	ErrUserEmailAlreadyExistCode = 7000010
	ErrUserInfoInvalidateCode    = 7000002
	ErrUserResourceNotFound      = 7000003
	ErrUserInvalidParamCode      = 7000020
	ErrUserPermissionCode        = 7000040
)

var (
	USER201001 = &BaseErrorCode{
		code:       "USER201001",
		message:    "Invalid username format",
		messageZH:  "用户名格式错误",
		messageEN:  "Invalid username format",
		httpStatus: http.StatusBadRequest,
	}
	USER201002 = &BaseErrorCode{
		code:       "USER201002",
		message:    "Email already exists",
		messageZH:  "邮箱已存在",
		messageEN:  "Email already exists",
		httpStatus: http.StatusConflict,
	}
	USER401001 = &BaseErrorCode{
		code:       "USER401001",
		message:    "User not found",
		messageZH:  "用户不存在",
		messageEN:  "User not found",
		httpStatus: http.StatusNotFound,
	}
	USER403001 = &BaseErrorCode{
		code:       "USER403001",
		message:    "Account has been disabled",
		messageZH:  "账号已禁用",
		messageEN:  "Account has been disabled",
		httpStatus: http.StatusForbidden,
	}
	USER500001 = &BaseErrorCode{
		code:       "USER500001",
		message:    "Failed to create user",
		messageZH:  "用户创建失败",
		messageEN:  "Failed to create user",
		httpStatus: http.StatusInternalServerError,
	}
	USER500002 = &BaseErrorCode{
		code:       "USER500002",
		message:    "Failed to update user",
		messageZH:  "用户更新失败",
		messageEN:  "Failed to update user",
		httpStatus: http.StatusInternalServerError,
	}
)
