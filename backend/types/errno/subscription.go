// backend/types/errno/subscription.go
package errno

import (
	"net/http"
)

// 订阅相关错误码常量
const (
	ErrSubscriptionNotFoundCode      = 4000001
	ErrSubscriptionExpiredCode       = 4003001
	ErrSubscriptionInactiveCode      = 4003002
	ErrSubscriptionSuspendedCode     = 4003003
	ErrSubscriptionPaymentFailedCode = 4005001
	ErrSubscriptionPaymentRequiredCode = 4005002
)

var (
	SUBSCRIPTION402001 = &BaseErrorCode{
		code:       "SUBSCRIPTION402001",
		message:    "Subscription has expired",
		messageZH:  "订阅已过期",
		messageEN:  "Subscription has expired",
		httpStatus: http.StatusPaymentRequired,
	}
	SUBSCRIPTION402002 = &BaseErrorCode{
		code:       "SUBSCRIPTION402002",
		message:    "Subscription quota exhausted",
		messageZH:  "订阅配额已用完",
		messageEN:  "Subscription quota exhausted",
		httpStatus: http.StatusPaymentRequired,
	}
	SUBSCRIPTION404001 = &BaseErrorCode{
		code:       "SUBSCRIPTION404001",
		message:    "Subscription not found",
		messageZH:  "订阅不存在",
		messageEN:  "Subscription not found",
		httpStatus: http.StatusNotFound,
	}
	SUBSCRIPTION409001 = &BaseErrorCode{
		code:       "SUBSCRIPTION409001",
		message:    "Subscription already exists",
		messageZH:  "订阅已存在",
		messageEN:  "Subscription already exists",
		httpStatus: http.StatusConflict,
	}
	SUBSCRIPTION500001 = &BaseErrorCode{
		code:       "SUBSCRIPTION500001",
		message:    "Failed to create subscription",
		messageZH:  "订阅创建失败",
		messageEN:  "Failed to create subscription",
		httpStatus: http.StatusInternalServerError,
	}
)
