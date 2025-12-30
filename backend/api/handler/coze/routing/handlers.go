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
}

// NewRoutingHandlers 创建路由模块的所有Handler实例
func NewRoutingHandlers(
	intentService *routingapp.IntentRecognitionService,
	optimizerService *routingapp.RoutingOptimizerService,
	ruleConfigService *routingapp.RoutingRuleConfigService,
	abTestService *routingapp.ABTestService,
	learningService *routingapp.RoutingLearningService,
	logger *zap.Logger,
) *RoutingHandlers {
	return &RoutingHandlers{
		IntentRecognition:  NewIntentRecognitionHandler(intentService, logger),
		RoutingOptimizer:    NewRoutingOptimizerHandler(optimizerService, logger),
		RoutingRuleConfig:  NewRoutingRuleConfigHandler(ruleConfigService, logger),
		ABTest:             NewABTestHandler(abTestService, logger),
		RoutingLearning:    NewRoutingLearningHandler(learningService, logger),
	}
}
