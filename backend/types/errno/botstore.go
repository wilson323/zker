// backend/types/errno/botstore.go
package errno

import "github.com/coze-dev/coze-studio/backend/pkg/errorx"

// =====================================================
// 统一错误码规范 - BotStore模块
// 错误码段: 206 000 000 ~ 206 999 999
// =====================================================

const (
	// BotStore通用错误 (206 000 000 ~ 206 009 999)
	ErrBotStoreInvalidParamCode = 206000001
	ErrBotStoreNotFoundCode     = 206000002
	ErrBotStorePermissionCode   = 206000003
	ErrBotStorePublishCode      = 206000004
	ErrBotStoreUpdateCode       = 206000005
	ErrBotStoreDeleteCode       = 206000006
	ErrBotStoreReviewCode       = 206000007

	// 通用错误 (206 010 000 ~ 206 019 999)
	// Note: ErrBotStoreInvalidParam was moved to var section to avoid duplication
	ErrBotStoreInvalidBotIDCode         = 206010002 // Bot ID无效
	ErrBotStoreInvalidBotNameCode       = 206010003 // Bot名称无效
	ErrBotStoreInvalidTenantIDCode      = 206010004 // 租户ID无效
	ErrBotStoreInvalidPublisherIDCode   = 206010005 // 发布者ID无效
	ErrBotStoreInvalidPriceCode         = 206010006 // 价格无效
	ErrBotStoreInvalidCategoryCode      = 206010007 // 分类无效
	ErrBotStoreInvalidRatingCode        = 206010008 // 评分无效
	ErrBotStoreCommentTooLongCode       = 206010009 // 评论内容过长
	ErrBotStoreInvalidItemIDCode        = 206010010 // 商品ID无效

	// 权限错误 (206 020 000 ~ 206 029 999)
	ErrBotStorePermissionDeniedCode        = 206020001 // 无权访问Bot商店项目
	ErrBotStoreNoModifyPermissionCode      = 206020002 // 无权修改Bot商店项目
	ErrBotStoreNoDeletePermissionCode      = 206020003 // 无权删除Bot商店项目
	ErrBotStoreNoReviewPermissionCode      = 206020004 // 无权审核Bot商店项目
	ErrBotStoreAlreadyPublishedCode        = 206020005 // Bot已发布
	ErrBotStorePendingReviewCode           = 206020006 // Bot待审核中
	ErrBotStoreCannotModifyPublishedItemCode = 206020007 // 无法修改已发布的项目
	ErrBotStoreNoModifyReviewPermissionCode = 206020008 // 无权修改此评论
	ErrBotStoreNoDeleteReviewPermissionCode = 206020009 // 无权删除此评论
	ErrBotStoreAlreadyReviewedCode         = 206020010 // 已经评论过此Bot

	// 资源不存在错误 (206 030 000 ~ 206 039 999)
	ErrBotStoreItemNotFoundCode = 206030001 // Bot商店项目不存在
	ErrBotStoreCategoryNotFoundCode     = 206030002 // 分类不存在
	ErrBotStoreReviewNotFoundCode       = 206030003 // 评论不存在

	// 状态错误 (206 040 000 ~ 206 049 999)
	ErrBotStoreInvalidStatusCode        = 206040001 // Bot状态无效
	ErrBotStoreRejectReasonRequiredCode = 206040002 // 拒绝原因必填

	// 服务器错误 (206 050 000 ~ 206 059 999)
	ErrBotStoreCreateBotStoreItemFailedCode  = 206050001 // 创建Bot商店项目失败
	ErrBotStoreUpdateBotStoreItemFailedCode  = 206050002 // 更新Bot商店项目失败
	ErrBotStoreDeleteBotStoreItemFailedCode  = 206050003 // 删除Bot商店项目失败
	ErrBotStoreGetBotStoreItemFailedCode     = 206050004 // 获取Bot商店项目失败
	ErrBotStoreReviewBotStoreItemFailedCode  = 206050005 // 审核Bot商店项目失败
	ErrBotStoreCreateReviewFailedCode        = 206050006 // 创建评论失败
	ErrBotStoreUpdateReviewFailedCode        = 206050007 // 更新评论失败
	ErrBotStoreDeleteReviewFailedCode        = 206050008 // 删除评论失败
	ErrBotStoreGetReviewFailedCode           = 206050009 // 获取评论失败
	ErrBotStoreUpdateRatingStatsFailedCode   = 206050010 // 更新评分统计失败
)

// 便捷错误变量（用于错误处理）
var (
	// BotStore通用错误
	ErrBotStoreInvalidParam = errorx.New(ErrBotStoreInvalidParamCode)
	ErrBotStoreNotFound     = errorx.New(ErrBotStoreNotFoundCode)
	ErrBotStorePermission   = errorx.New(ErrBotStorePermissionCode)
	ErrBotStorePublish      = errorx.New(ErrBotStorePublishCode)
	ErrBotStoreUpdate       = errorx.New(ErrBotStoreUpdateCode)
	ErrBotStoreDelete       = errorx.New(ErrBotStoreDeleteCode)
	ErrBotStoreReview       = errorx.New(ErrBotStoreReviewCode)

	// 参数验证错误
	ErrInvalidBotID       = errorx.New(ErrBotStoreInvalidBotIDCode)
	ErrInvalidBotName     = errorx.New(ErrBotStoreInvalidBotNameCode)
	ErrInvalidTenantID    = errorx.New(ErrBotStoreInvalidTenantIDCode)
	ErrInvalidPublisherID = errorx.New(ErrBotStoreInvalidPublisherIDCode)
	ErrInvalidPrice       = errorx.New(ErrBotStoreInvalidPriceCode)
	ErrInvalidCategory    = errorx.New(ErrBotStoreInvalidCategoryCode)
	ErrInvalidRating      = errorx.New(ErrBotStoreInvalidRatingCode)
	ErrCommentTooLong     = errorx.New(ErrBotStoreCommentTooLongCode)
	ErrInvalidItemID      = errorx.New(ErrBotStoreInvalidItemIDCode)

	// 权限错误
	ErrPermissionDenied          = errorx.New(ErrBotStorePermissionDeniedCode)
	ErrNoModifyPermission        = errorx.New(ErrBotStoreNoModifyPermissionCode)
	ErrNoDeletePermission        = errorx.New(ErrBotStoreNoDeletePermissionCode)
	ErrNoReviewPermission        = errorx.New(ErrBotStoreNoReviewPermissionCode)
	ErrAlreadyPublished          = errorx.New(ErrBotStoreAlreadyPublishedCode)
	ErrPendingReview             = errorx.New(ErrBotStorePendingReviewCode)
	ErrCannotModifyPublishedItem = errorx.New(ErrBotStoreCannotModifyPublishedItemCode)
	ErrNoModifyReviewPermission  = errorx.New(ErrBotStoreNoModifyReviewPermissionCode)
	ErrNoDeleteReviewPermission  = errorx.New(ErrBotStoreNoDeleteReviewPermissionCode)
	ErrAlreadyReviewed           = errorx.New(ErrBotStoreAlreadyReviewedCode)

	// 资源不存在错误
	ErrBotStoreItemNotFound = errorx.New(ErrBotStoreItemNotFoundCode)
	ErrCategoryNotFound     = errorx.New(ErrBotStoreCategoryNotFoundCode)
	ErrReviewNotFound       = errorx.New(ErrBotStoreReviewNotFoundCode)

	// 状态错误
	ErrInvalidStatus        = errorx.New(ErrBotStoreInvalidStatusCode)
	ErrRejectReasonRequired = errorx.New(ErrBotStoreRejectReasonRequiredCode)

	// 服务器错误
	ErrCreateBotStoreItemFailed  = errorx.New(ErrBotStoreCreateBotStoreItemFailedCode)
	ErrUpdateBotStoreItemFailed  = errorx.New(ErrBotStoreUpdateBotStoreItemFailedCode)
	ErrDeleteBotStoreItemFailed  = errorx.New(ErrBotStoreDeleteBotStoreItemFailedCode)
	ErrGetBotStoreItemFailed     = errorx.New(ErrBotStoreGetBotStoreItemFailedCode)
	ErrReviewBotStoreItemFailed  = errorx.New(ErrBotStoreReviewBotStoreItemFailedCode)
	ErrCreateReviewFailed        = errorx.New(ErrBotStoreCreateReviewFailedCode)
	ErrUpdateReviewFailed        = errorx.New(ErrBotStoreUpdateReviewFailedCode)
	ErrDeleteReviewFailed        = errorx.New(ErrBotStoreDeleteReviewFailedCode)
	ErrGetReviewFailed           = errorx.New(ErrBotStoreGetReviewFailedCode)
	ErrUpdateRatingStatsFailed   = errorx.New(ErrBotStoreUpdateRatingStatsFailedCode)
)

// 向后兼容：保留旧的BSTORE错误码常量（已废弃）
// TODO: 逐步迁移到206段错误码，移除旧段定义
const (
	BSTORE400001 = 206010001
	BSTORE400002 = 206010002
	BSTORE400003 = 206010003
	BSTORE400004 = 206010004
	BSTORE400005 = 206010005
	BSTORE400006 = 206010006
	BSTORE400007 = 206010007
	BSTORE400008 = 206010008
	BSTORE400009 = 206010009
	BSTORE400010 = 206010010

	BSTORE403001 = 206020001
	BSTORE403002 = 206020002
	BSTORE403003 = 206020003
	BSTORE403004 = 206020004
	BSTORE403005 = 206020005
	BSTORE403006 = 206020006
	BSTORE403007 = 206020007
	BSTORE403008 = 206020008
	BSTORE403009 = 206020009
	BSTORE403010 = 206020010

	BSTORE404001 = 206030001
	BSTORE404002 = 206030002
	BSTORE404003 = 206030003

	BSTORE409001 = 206040001
	BSTORE409002 = 206040002

	BSTORE500001 = 206050001
	BSTORE500002 = 206050002
	BSTORE500003 = 206050003
	BSTORE500004 = 206050004
	BSTORE500005 = 206050005
	BSTORE500006 = 206050006
	BSTORE500007 = 206050007
	BSTORE500008 = 206050008
	BSTORE500009 = 206050009
	BSTORE500010 = 206050010
)
