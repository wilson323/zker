// backend/types/errno/cache.go
package errno

import "net/http"

// 缓存模块错误码(CACHE前缀)
// 遵循ZKER统一错误码定义规范

var (
	// 缓存连接错误 (500)
	CACHE500001 = &BaseErrorCode{
		code:       "CACHE500001",
		message:    "Cache connection failed",
		messageZH:  "缓存连接失败",
		messageEN:  "Cache connection failed",
		httpStatus: http.StatusInternalServerError,
	}
	CACHE500002 = &BaseErrorCode{
		code:       "CACHE500002",
		message:    "Cache operation failed",
		messageZH:  "缓存操作失败",
		messageEN:  "Cache operation failed",
		httpStatus: http.StatusInternalServerError,
	}
	CACHE500003 = &BaseErrorCode{
		code:       "CACHE500003",
		message:    "Cache avalanche",
		messageZH:  "缓存雪崩",
		messageEN:  "Cache avalanche",
		httpStatus: http.StatusInternalServerError,
	}
	CACHE500004 = &BaseErrorCode{
		code:       "CACHE500004",
		message:    "Cache penetration",
		messageZH:  "缓存穿透",
		messageEN:  "Cache penetration",
		httpStatus: http.StatusInternalServerError,
	}
	CACHE500005 = &BaseErrorCode{
		code:       "CACHE500005",
		message:    "Cache breakdown",
		messageZH:  "缓存击穿",
		messageEN:  "Cache breakdown",
		httpStatus: http.StatusInternalServerError,
	}
	CACHE500006 = &BaseErrorCode{
		code:       "CACHE500006",
		message:    "Cache serialization failed",
		messageZH:  "缓存序列化失败",
		messageEN:  "Cache serialization failed",
		httpStatus: http.StatusInternalServerError,
	}
	CACHE500007 = &BaseErrorCode{
		code:       "CACHE500007",
		message:    "Cache deserialization failed",
		messageZH:  "缓存反序列化失败",
		messageEN:  "Cache deserialization failed",
		httpStatus: http.StatusInternalServerError,
	}

	// 缓存参数错误 (400)
	CACHE400001 = &BaseErrorCode{
		code:       "CACHE400001",
		message:    "Invalid cache key",
		messageZH:  "缓存键无效",
		messageEN:  "Invalid cache key",
		httpStatus: http.StatusBadRequest,
	}
	CACHE400002 = &BaseErrorCode{
		code:       "CACHE400002",
		message:    "Invalid cache value",
		messageZH:  "缓存值无效",
		messageEN:  "Invalid cache value",
		httpStatus: http.StatusBadRequest,
	}
	CACHE400003 = &BaseErrorCode{
		code:       "CACHE400003",
		message:    "Invalid TTL parameter",
		messageZH:  "TTL参数无效",
		messageEN:  "Invalid TTL parameter",
		httpStatus: http.StatusBadRequest,
	}

	// 缓存不存在 (404)
	CACHE404001 = &BaseErrorCode{
		code:       "CACHE404001",
		message:    "Cache key not found",
		messageZH:  "缓存键不存在",
		messageEN:  "Cache key not found",
		httpStatus: http.StatusNotFound,
	}

	// ===== 便捷错误变量（向后兼容） =====
	ErrCacheMiss = &BaseErrorCode{
		code:       "CACHE404001",
		message:    "Cache key not found",
		messageZH:  "缓存键不存在",
		messageEN:  "Cache key not found",
		httpStatus: http.StatusNotFound,
	}
	ErrCacheGetFailed = &BaseErrorCode{
		code:       "CACHE500002",
		message:    "Cache operation failed",
		messageZH:  "缓存获取失败",
		messageEN:  "Cache operation failed",
		httpStatus: http.StatusInternalServerError,
	}
	ErrCacheDecodeFailed = &BaseErrorCode{
		code:       "CACHE500007",
		message:    "Cache deserialization failed",
		messageZH:  "缓存解码失败",
		messageEN:  "Cache deserialization failed",
		httpStatus: http.StatusInternalServerError,
	}
	ErrCacheEncodeFailed = &BaseErrorCode{
		code:       "CACHE500006",
		message:    "Cache serialization failed",
		messageZH:  "缓存编码失败",
		messageEN:  "Cache serialization failed",
		httpStatus: http.StatusInternalServerError,
	}

	// 其他缓存错误码
	ErrCacheSetFailed = &BaseErrorCode{
		code:       "CACHE500002",
		message:    "Cache operation failed",
		messageZH:  "缓存设置失败",
		messageEN:  "Cache operation failed",
		httpStatus: http.StatusInternalServerError,
	}
	ErrCacheDeleteFailed = &BaseErrorCode{
		code:       "CACHE500002",
		message:    "Cache operation failed",
		messageZH:  "缓存删除失败",
		messageEN:  "Cache operation failed",
		httpStatus: http.StatusInternalServerError,
	}
	ErrCacheCheckFailed = &BaseErrorCode{
		code:       "CACHE500002",
		message:    "Cache operation failed",
		messageZH:  "缓存检查失败",
		messageEN:  "Cache operation failed",
		httpStatus: http.StatusInternalServerError,
	}
	ErrCacheExpireFailed = &BaseErrorCode{
		code:       "CACHE500002",
		message:    "Cache operation failed",
		messageZH:  "缓存过期时间设置失败",
		messageEN:  "Cache operation failed",
		httpStatus: http.StatusInternalServerError,
	}
)

// ErrNotImplemented 未实现错误（使用通用错误码）
var ErrNotImplemented = InternalError
