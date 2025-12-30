// backend/types/errno/botstore.go
package errno

var (
	// Bot商店通用错误 (BSTORE)
	BSTORE400001 = &BaseErrorCode{"BSTORE400001", "请求参数无效", "请求参数无效", "Invalid request parameter", 400}
	BSTORE400002 = &BaseErrorCode{"BSTORE400002", "Bot ID无效", "Bot ID无效", "Invalid bot ID", 400}
	BSTORE400003 = &BaseErrorCode{"BSTORE400003", "Bot名称无效", "Bot名称无效", "Invalid bot name", 400}
	BSTORE400004 = &BaseErrorCode{"BSTORE400004", "租户ID无效", "租户ID无效", "Invalid tenant ID", 400}
	BSTORE400005 = &BaseErrorCode{"BSTORE400005", "发布者ID无效", "发布者ID无效", "Invalid publisher ID", 400}
	BSTORE400006 = &BaseErrorCode{"BSTORE400006", "价格无效", "价格无效", "Invalid price", 400}
	BSTORE400007 = &BaseErrorCode{"BSTORE400007", "分类无效", "分类无效", "Invalid category", 400}

	// Bot商店权限错误 (403)
	BSTORE403001 = &BaseErrorCode{"BSTORE403001", "无权访问Bot商店项目", "无权访问Bot商店项目", "No permission to access bot store item", 403}
	BSTORE403002 = &BaseErrorCode{"BSTORE403002", "无权修改Bot商店项目", "无权修改Bot商店项目", "No permission to modify bot store item", 403}
	BSTORE403003 = &BaseErrorCode{"BSTORE403003", "无权删除Bot商店项目", "无权删除Bot商店项目", "No permission to delete bot store item", 403}
	BSTORE403004 = &BaseErrorCode{"BSTORE403004", "无权审核Bot商店项目", "无权审核Bot商店项目", "No permission to review bot store item", 403}
	BSTORE403005 = &BaseErrorCode{"BSTORE403005", "Bot已发布", "Bot已发布", "Bot already published", 403}
	BSTORE403006 = &BaseErrorCode{"BSTORE403006", "Bot待审核中", "Bot待审核中", "Bot is pending review", 403}
	BSTORE403007 = &BaseErrorCode{"BSTORE403007", "无法修改已发布的项目", "无法修改已发布的项目", "Cannot modify published item", 403}

	// Bot商店资源不存在错误 (404)
	BSTORE404001 = &BaseErrorCode{"BSTORE404001", "Bot商店项目不存在", "Bot商店项目不存在", "Bot store item not found", 404}
	BSTORE404002 = &BaseErrorCode{"BSTORE404002", "分类不存在", "分类不存在", "Category not found", 404}

	// Bot商店状态错误 (409)
	BSTORE409001 = &BaseErrorCode{"BSTORE409001", "Bot状态无效", "Bot状态无效", "Invalid bot status", 409}
	BSTORE409002 = &BaseErrorCode{"BSTORE409002", "拒绝原因必填", "拒绝原因必填", "Reject reason is required", 409}

	// Bot商店服务器错误 (500)
	BSTORE500001 = &BaseErrorCode{"BSTORE500001", "创建Bot商店项目失败", "创建Bot商店项目失败", "Failed to create bot store item", 500}
	BSTORE500002 = &BaseErrorCode{"BSTORE500002", "更新Bot商店项目失败", "更新Bot商店项目失败", "Failed to update bot store item", 500}
	BSTORE500003 = &BaseErrorCode{"BSTORE500003", "删除Bot商店项目失败", "删除Bot商店项目失败", "Failed to delete bot store item", 500}
	BSTORE500004 = &BaseErrorCode{"BSTORE500004", "获取Bot商店项目失败", "获取Bot商店项目失败", "Failed to get bot store item", 500}
	BSTORE500005 = &BaseErrorCode{"BSTORE500005", "审核Bot商店项目失败", "审核Bot商店项目失败", "Failed to review bot store item", 500}
)

// 便捷错误变量
var (
	ErrInvalidBotID              = BSTORE400002
	ErrInvalidBotName            = BSTORE400003
	ErrInvalidTenantID           = BSTORE400004
	ErrInvalidPublisherID        = BSTORE400005
	ErrInvalidPrice              = BSTORE400006
	ErrInvalidCategory           = BSTORE400007
	ErrPermissionDenied          = BSTORE403001
	ErrAlreadyPublished          = BSTORE403005
	ErrPendingReview             = BSTORE403006
	ErrCannotModifyPublishedItem = BSTORE403007
	ErrBotStoreItemNotFound      = BSTORE404001
	ErrInvalidStatus             = BSTORE409001
	ErrRejectReasonRequired      = BSTORE409002
)
