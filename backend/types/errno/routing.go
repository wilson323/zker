// backend/types/errno/routing.go
package errno

import "github.com/coze-dev/coze-studio/backend/pkg/errorx"

// =====================================================
// 统一错误码规范 - Routing模块
// 错误码段: 205 000 000 ~ 205 999 999
// =====================================================

const (
	// Routing基础错误 (205 000 000 ~ 205 009 999)
	ErrRoutingNotFoundCode     = 205000001 // 路由不存在
	ErrRoutingAlreadyExistsCode = 205000002 // 路由已存在
	ErrRoutingInvalidParamCode  = 205000003 // 路由参数无效
	ErrRoutingPermissionCode    = 205000004 // 路由权限不足
	ErrRoutingFailedCode        = 205000005 // 路由操作失败

	// 资源不存在 (205 010 000 ~ 205 019 999)
	ErrRouteNotFoundCode       = 205010001 // 路由规则不存在
	ErrIntentNotFoundCode       = 205010002 // 意图不存在
	ErrABTestNotFoundCode       = 205010003 // A/B测试不存在

	// 资源已存在 (205 020 000 ~ 205 029 999)
	ErrRouteAlreadyExistsCode = 205020001 // 路由规则已存在
	ErrIntentAlreadyExistsCode = 205020002 // 意图已存在

	// 参数无效 (205 030 000 ~ 205 039 999)
	ErrRouteFormatInvalidCode = 205030001 // 路由规则格式错误
	ErrRouteParamInvalidCode  = 205030002 // 路由参数无效
	ErrIntentParamInvalidCode = 205030003 // 意图参数无效
	ErrABTestParamInvalidCode = 205030004 // A/B测试参数无效
	ErrRoutingStrategyInvalidCode = 205030005 // 路由策略无效

	// 操作失败 (205 050 000 ~ 205 059 999)
	ErrRoutingDecisionFailedCode  = 205050001 // 路由决策失败
	ErrIntentRecognitionFailedCode = 205050002 // 意图识别失败
	ErrRoutingOptimizationFailedCode = 205050003 // 路由优化失败
	ErrEntityExtractionFailedCode   = 205050004 // 实体提取失败
	ErrVectorSearchFailedCode       = 205050005 // 向量搜索失败
	ErrABTestExecutionFailedCode    = 205050006 // A/B测试执行失败
	ErrRoutingLearningFailedCode    = 205050007 // 路由学习失败
	ErrScoreCalculationFailedCode   = 205050008 // 评分计算失败
	ErrRouteMatchFailedCode         = 205050009 // 路由匹配失败
	ErrRouteScoreFailedCode         = 205050010 // 路由评分失败
	ErrRouteSearchFailedCode        = 205050011 // 路由搜索失败
	ErrIntentIndexingFailedCode     = 205050012 // 意图索引失败
)

// 便捷错误变量（用于错误处理）
var (
	// 基础错误
	ErrRouteNotFound       = errorx.New(ErrRouteNotFoundCode)
	ErrIntentNotFound       = errorx.New(ErrIntentNotFoundCode)
	ErrABTestNotFound       = errorx.New(ErrABTestNotFoundCode)
	ErrRouteAlreadyExists   = errorx.New(ErrRouteAlreadyExistsCode)
	ErrIntentAlreadyExists  = errorx.New(ErrIntentAlreadyExistsCode)
	ErrRouteFormatInvalid   = errorx.New(ErrRouteFormatInvalidCode)
	ErrRouteParamInvalid    = errorx.New(ErrRouteParamInvalidCode)
	ErrIntentParamInvalid   = errorx.New(ErrIntentParamInvalidCode)
	ErrABTestParamInvalid   = errorx.New(ErrABTestParamInvalidCode)
	ErrRoutingStrategyInvalid = errorx.New(ErrRoutingStrategyInvalidCode)
	ErrRoutingDecisionFailed  = errorx.New(ErrRoutingDecisionFailedCode)
	ErrIntentRecognitionFailed = errorx.New(ErrIntentRecognitionFailedCode)
	ErrRoutingOptimizationFailed = errorx.New(ErrRoutingOptimizationFailedCode)
	ErrEntityExtractionFailed   = errorx.New(ErrEntityExtractionFailedCode)
	ErrVectorSearchFailed       = errorx.New(ErrVectorSearchFailedCode)
	ErrABTestExecutionFailed    = errorx.New(ErrABTestExecutionFailedCode)
	ErrRoutingLearningFailed    = errorx.New(ErrRoutingLearningFailedCode)
	ErrScoreCalculationFailed   = errorx.New(ErrScoreCalculationFailedCode)
	ErrRouteMatchFailed         = errorx.New(ErrRouteMatchFailedCode)
	ErrRouteScoreFailed         = errorx.New(ErrRouteScoreFailedCode)
	ErrRouteSearchFailed        = errorx.New(ErrRouteSearchFailedCode)
	ErrIntentIndexingFailed     = errorx.New(ErrIntentIndexingFailedCode)
)

// 向后兼容：保留旧的ROUTING错误码常量（已废弃）
// TODO: 逐步迁移到205段错误码，移除旧段定义
const (
	// 209段 → 205段映射
	ROUTING209001 = 205010001
	ROUTING209002 = 205010002
	ROUTING209003 = 205010003
	ROUTING209010 = 205020001
	ROUTING209011 = 205020002
	ROUTING201001 = 205030001
	ROUTING209020 = 205030002
	ROUTING209021 = 205030003
	ROUTING209022 = 205030004
	ROUTING209023 = 205030005
	ROUTING500001 = 205050001
	ROUTING500002 = 205050002
	ROUTING209050 = 205050003
	ROUTING209051 = 205050004
	ROUTING209052 = 205050005
	ROUTING209053 = 205050006
	ROUTING209054 = 205050007
	ROUTING209055 = 205050008

	// 向后兼容别名
	ROUTING404001 = 205010001
)
