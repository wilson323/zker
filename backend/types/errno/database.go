// backend/types/errno/database.go
package errno

import "net/http"

// 数据库模块错误码(DB前缀)
// 遵循ZKER统一错误码定义规范

var (
	// 数据库连接错误 (500)
	DB500001 = &BaseErrorCode{
		code:       "DB500001",
		message:    "Database connection failed",
		messageZH:  "数据库连接失败",
		messageEN:  "Database connection failed",
		httpStatus: http.StatusInternalServerError,
	}
	DB500002 = &BaseErrorCode{
		code:       "DB500002",
		message:    "SQL execution failed",
		messageZH:  "SQL执行失败",
		messageEN:  "SQL execution failed",
		httpStatus: http.StatusInternalServerError,
	}
	DB500003 = &BaseErrorCode{
		code:       "DB500003",
		message:    "Transaction failed",
		messageZH:  "事务失败",
		messageEN:  "Transaction failed",
		httpStatus: http.StatusInternalServerError,
	}
	DB500004 = &BaseErrorCode{
		code:       "DB500004",
		message:    "Database deadlock",
		messageZH:  "死锁错误",
		messageEN:  "Database deadlock",
		httpStatus: http.StatusInternalServerError,
	}
	DB500005 = &BaseErrorCode{
		code:       "DB500005",
		message:    "Connection pool exhausted",
		messageZH:  "连接池耗尽",
		messageEN:  "Connection pool exhausted",
		httpStatus: http.StatusInternalServerError,
	}
	DB500006 = &BaseErrorCode{
		code:       "DB500006",
		message:    "Database operation timeout",
		messageZH:  "数据库超时",
		messageEN:  "Database operation timeout",
		httpStatus: http.StatusInternalServerError,
	}
	DB500007 = &BaseErrorCode{
		code:       "DB500007",
		message:    "Database migration failed",
		messageZH:  "数据库迁移失败",
		messageEN:  "Database migration failed",
		httpStatus: http.StatusInternalServerError,
	}
	DB500008 = &BaseErrorCode{
		code:       "DB500008",
		message:    "Batch operation failed",
		messageZH:  "批量操作失败",
		messageEN:  "Batch operation failed",
		httpStatus: http.StatusInternalServerError,
	}

	// 数据库查询错误 (400/404)
	DB400001 = &BaseErrorCode{
		code:       "DB400001",
		message:    "Invalid query parameters",
		messageZH:  "无效的查询参数",
		messageEN:  "Invalid query parameters",
		httpStatus: http.StatusBadRequest,
	}
	DB400002 = &BaseErrorCode{
		code:       "DB400002",
		message:    "Query syntax error",
		messageZH:  "查询语法错误",
		messageEN:  "Query syntax error",
		httpStatus: http.StatusBadRequest,
	}
	DB400003 = &BaseErrorCode{
		code:       "DB400003",
		message:    "Invalid sort parameters",
		messageZH:  "排序参数无效",
		messageEN:  "Invalid sort parameters",
		httpStatus: http.StatusBadRequest,
	}
	DB400004 = &BaseErrorCode{
		code:       "DB400004",
		message:    "Invalid pagination parameters",
		messageZH:  "分页参数无效",
		messageEN:  "Invalid pagination parameters",
		httpStatus: http.StatusBadRequest,
	}
	DB400005 = &BaseErrorCode{
		code:       "DB400005",
		message:    "Invalid field list",
		messageZH:  "字段列表无效",
		messageEN:  "Invalid field list",
		httpStatus: http.StatusBadRequest,
	}

	// 数据库约束错误 (409)
	DB409001 = &BaseErrorCode{
		code:       "DB409001",
		message:    "Unique constraint violation",
		messageZH:  "唯一约束冲突",
		messageEN:  "Unique constraint violation",
		httpStatus: http.StatusConflict,
	}
	DB409002 = &BaseErrorCode{
		code:       "DB409002",
		message:    "Foreign key constraint violation",
		messageZH:  "外键约束冲突",
		messageEN:  "Foreign key constraint violation",
		httpStatus: http.StatusConflict,
	}
	DB409003 = &BaseErrorCode{
		code:       "DB409003",
		message:    "Data conflict",
		messageZH:  "数据冲突",
		messageEN:  "Data conflict",
		httpStatus: http.StatusConflict,
	}
)
