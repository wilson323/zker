/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package routing

import (
	"go.uber.org/zap"

	routingapp "github.com/coze-dev/coze-studio/backend/application/routing"
)

// RoutingHandlers 路由模块的Handler容器
type RoutingHandlers struct {
	IntentRecognition  *IntentRecognitionHandler
	RoutingOptimizer    *RoutingOptimizerHandler
	RoutingRuleConfig  *RoutingRuleConfigHandler
	ABTest             *ABTestHandler
	RoutingLearning    *RoutingLearningHandler
	RoutingMonitoring  *RoutingMonitoringHandler
}

// NewRoutingHandlers 创建路由模块的所有Handler实例
// 使用应用层初始化的全局服务实例
func NewRoutingHandlers(logger *zap.Logger) *RoutingHandlers {
	return &RoutingHandlers{
		IntentRecognition:  NewIntentRecognitionHandler(routingapp.IntentRecognitionSVC, logger),
		RoutingOptimizer:    NewRoutingOptimizerHandler(routingapp.RoutingOptimizerSVC, logger),
		RoutingRuleConfig:  NewRoutingRuleConfigHandler(routingapp.RoutingRuleConfigSVC, logger),
		ABTest:             NewABTestHandler(routingapp.ABTestSVC, logger),
		RoutingLearning:    NewRoutingLearningHandler(routingapp.RoutingLearningSVC, logger),
		RoutingMonitoring:  NewRoutingMonitoringHandler(routingapp.Handler, routingapp.ABTestSVC, routingapp.RoutingLearningSVC, routingapp.LoadMonitor, routingapp.AdvancedLoadBalancer),
	}
}
