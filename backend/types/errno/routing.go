// backend/types/errno/routing.go
package errno

// 路由模块错误码 (209 xxx xxx)
// 209001-209009: 资源不存在
// 209010-209019: 资源已存在
// 209020-209039: 参数无效
// 209040-209049: 权限不足
// 209050-209089: 操作失败

var (
	// 路由规则相关 (209001-209009)
	ROUTING209001 = &BaseErrorCode{"ROUTING209001", "路由规则不存在", "Routing rule not found", "Route not found", 404}
	ROUTING209002 = &BaseErrorCode{"ROUTING209002", "意图不存在", "Intent not found", "Intent not found", 404}
	ROUTING209003 = &BaseErrorCode{"ROUTING209003", "A/B测试不存在", "A/B test not found", "A/B test not found", 404}

	// 路由规则相关 (209010-209019)
	ROUTING209010 = &BaseErrorCode{"ROUTING209010", "路由规则已存在", "Routing rule already exists", "Routing rule already exists", 409}
	ROUTING209011 = &BaseErrorCode{"ROUTING209011", "意图已存在", "Intent already exists", "Intent already exists", 409}

	// 路由规则相关 (209020-209039)
	ROUTING201001 = &BaseErrorCode{"ROUTING201001", "路由规则格式错误", "路由规则格式错误", "Invalid routing rule format", 400}
	ROUTING209020 = &BaseErrorCode{"ROUTING209020", "路由参数无效", "Invalid routing parameters", "Invalid routing parameters", 400}
	ROUTING209021 = &BaseErrorCode{"ROUTING209021", "意图参数无效", "Invalid intent parameters", "Invalid intent parameters", 400}
	ROUTING209022 = &BaseErrorCode{"ROUTING209022", "A/B测试参数无效", "Invalid A/B test parameters", "Invalid A/B test parameters", 400}
	ROUTING209023 = &BaseErrorCode{"ROUTING209023", "路由策略无效", "Invalid routing strategy", "Invalid routing strategy", 400}

	// 路由操作相关 (209050-209089)
	ROUTING500001 = &BaseErrorCode{"ROUTING500001", "路由决策失败", "路由决策失败", "Routing decision failed", 500}
	ROUTING500002 = &BaseErrorCode{"ROUTING500002", "意图识别失败", "意图识别失败", "Intent recognition failed", 500}
	ROUTING209050 = &BaseErrorCode{"ROUTING209050", "路由优化失败", "Routing optimization failed", "Routing optimization failed", 500}
	ROUTING209051 = &BaseErrorCode{"ROUTING209051", "实体提取失败", "Entity extraction failed", "Entity extraction failed", 500}
	ROUTING209052 = &BaseErrorCode{"ROUTING209052", "向量搜索失败", "Vector search failed", "Vector search failed", 500}
	ROUTING209053 = &BaseErrorCode{"ROUTING209053", "A/B测试执行失败", "A/B test execution failed", "A/B test execution failed", 500}
	ROUTING209054 = &BaseErrorCode{"ROUTING209054", "路由学习失败", "Routing learning failed", "Routing learning failed", 500}
	ROUTING209055 = &BaseErrorCode{"ROUTING209055", "评分计算失败", "Score calculation failed", "Score calculation failed", 500}

	// 兼容旧错误码
	ROUTING404001 = ROUTING209001
)
