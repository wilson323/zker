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
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
	routingservice "github.com/coze-dev/coze-studio/backend/domain/routing/service"
)

var (
	RoutingAppSVC           *RoutingApplicationService
	IntentRecognitionSVC   *IntentRecognitionService
	RoutingOptimizerSVC    *RoutingOptimizerService
	RoutingRuleConfigSVC   *RoutingRuleConfigService
	ABTestSVC              *ABTestService
	RoutingLearningSVC     *RoutingLearningService
)

// ServiceComponents 路由应用服务组件
type ServiceComponents struct {
	DB *gorm.DB
}

// InitService 初始化路由应用服务
func InitService(c *ServiceComponents) (*RoutingApplicationService, error) {
	// 1. 初始化领域仓储
	ruleRepo := repository.NewRoutingRuleRepository(c.DB)
	logRepo := repository.NewRoutingLogRepository(c.DB)

	// 2. 初始化领域服务
	ruleSVC := routingservice.NewRuleService(ruleRepo)
	logSVC := routingservice.NewRoutingLogService(logRepo)

	// 3. 初始化匹配器
	ruleMatcher := routingservice.NewRuleBasedMatcher(ruleRepo)
	similarityMatcher := routingservice.NewSimilarityMatcher(0.75) // 默认阈值0.75

	// 4. 初始化混合意图匹配器（只需2个参数）
	intentMatcher := routingservice.NewHybridIntentMatcher(
		ruleMatcher,
		similarityMatcher,
	)

	// 设置匹配器权重（可选）
	intentMatcher.SetWeights(0.5, 0.5)

	// 5. 初始化服务健康监控
	healthMonitor := routingservice.NewServiceHealthMonitor()
	loadMonitor := routingservice.NewLoadMonitor()

	// 6. 初始化评分路由引擎
	routingEngine := routingservice.NewRoutingEngine(
		intentMatcher,
		healthMonitor,
		loadMonitor,
		logRepo,
		0.3, // intent_weight
		0.2, // health_weight
		0.2, // load_weight
		0.1, // cost_weight
		0.1, // region_weight
	)

	// 7. 初始化负载均衡器和熔断器
	loadBalancer := routingservice.NewLoadBalancer("least_loaded", 80)
	circuitBreaker := routingservice.NewCircuitBreakerService()

	// 8. 初始化应用服务
	routingAppSvc := NewRoutingApplicationService(
		routingEngine,
		ruleSVC,
		logSVC,
		healthMonitor,
		loadBalancer,
		circuitBreaker,
	)

	// 9. 设置全局变量
	RoutingAppSVC = routingAppSvc

	// 10. 初始化Handler层所需的服务适配器
	IntentRecognitionSVC = NewIntentRecognitionService(routingAppSvc)
	RoutingOptimizerSVC = NewRoutingOptimizerService(routingAppSvc)
	RoutingRuleConfigSVC = NewRoutingRuleConfigService(routingAppSvc)
	ABTestSVC = NewABTestService(routingAppSvc)
	RoutingLearningSVC = NewRoutingLearningService(routingAppSvc)

	return routingAppSvc, nil
}
