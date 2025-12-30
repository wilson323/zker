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

package entity

import (
	"time"
)

// ABTestStatus A/B测试状态
type ABTestStatus string

const (
	ABTestStatusRunning   ABTestStatus = "running"   // 运行中
	ABTestStatusPaused    ABTestStatus = "paused"    // 暂停
	ABTestStatusCompleted ABTestStatus = "completed" // 已完成
)

// ABTest A/B测试实体
type ABTest struct {
	TestID          string       `json:"test_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID        string       `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	TestName        string       `json:"test_name" gorm:"type:varchar(200);not null"`
	Description     string       `json:"description" gorm:"type:text"`
	StrategyA       string       `json:"strategy_a" gorm:"type:varchar(100);not null"` // 路由策略A
	StrategyB       string       `json:"strategy_b" gorm:"type:varchar(100);not null"` // 路由策略B
	TrafficSplit    int          `json:"traffic_split" gorm:"default:50"`              // 流量分配A:B (50:50)
	Status          ABTestStatus `json:"status" gorm:"type:varchar(20);not null;index:idx_status"`
	Winner          string       `json:"winner" gorm:"type:varchar(100)"`              // 获胜策略
	Confidence      float64      `json:"confidence" gorm:"type:decimal(5,4)"`          // 统计显著性（0-1）
	SampleSize      int          `json:"sample_size"`                                  // 所需样本量
	CurrentSampleA  int          `json:"current_sample_a"`                             // 当前A组样本数
	CurrentSampleB  int          `json:"current_sample_b"`                             // 当前B组样本数
	StartTime       *time.Time   `json:"start_time"`
	EndTime         *time.Time   `json:"end_time"`
	CreatedAt       time.Time    `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time    `json:"updated_at" gorm:"not null"`
	DeletedAt       *time.Time   `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (ABTest) TableName() string {
	return "ab_tests"
}

// IsTestActive 测试是否激活
func (t *ABTest) IsTestActive() bool {
	return t.Status == ABTestStatusRunning && t.DeletedAt == nil
}

// GetStrategyForTraffic 根据流量分配获取策略
func (t *ABTest) GetStrategyForTraffic(trafficValue int) string {
	// trafficValue: 0-99 的随机数
	if trafficValue < t.TrafficSplit {
		return t.StrategyA
	}
	return t.StrategyB
}

// ABTestRecord A/B测试记录实体
type ABTestRecord struct {
	RecordID     string    `json:"record_id" gorm:"primaryKey;type:varchar(36)"`
	TestID       string    `json:"test_id" gorm:"type:varchar(36);not null;index:idx_test_id"`
	Strategy     string    `json:"strategy" gorm:"type:varchar(100);not null;index:idx_strategy"` // A or B
	UserID       string    `json:"user_id" gorm:"type:varchar(36);not null"`
	AgentID      string    `json:"agent_id" gorm:"type:varchar(36);not null"` // 路由到的Agent
	WorkflowID   string    `json:"workflow_id" gorm:"type:varchar(36)"`       // 路由到的Workflow（可选）
	ResponseTime int       `json:"response_time"` // 响应时间(ms)
	UserRating   int       `json:"user_rating"`   // 用户评分1-5
	Feedback     string    `json:"feedback" gorm:"type:text"` // 用户反馈
	CreatedAt    time.Time `json:"created_at" gorm:"not null;index:idx_created_at"`
}

// TableName 指定表名
func (ABTestRecord) TableName() string {
	return "ab_test_records"
}

// ABTestResult A/B测试结果统计
type ABTestResult struct {
	TestID              string  `json:"test_id"`
	StrategyAStats      *TestStrategyStats `json:"strategy_a_stats"`
	StrategyBStats      *TestStrategyStats `json:"strategy_b_stats"`
	Winner              string  `json:"winner"`
	Confidence          float64 `json:"confidence"`
	IsStatisticallySignificant bool `json:"is_statistically_significant"`
	Recommendation      string  `json:"recommendation"`
}

// TestStrategyStats 单个策略的统计数据
type TestStrategyStats struct {
	Strategy         string  `json:"strategy"`
	SampleCount      int     `json:"sample_count"`
	AvgResponseTime  float64 `json:"avg_response_time"` // ms
	AvgUserRating    float64 `json:"avg_user_rating"`   // 1-5
	SuccessRate      float64 `json:"success_rate"`      // 0-1
	ErrorRate        float64 `json:"error_rate"`        // 0-1
	ConversionRate   float64 `json:"conversion_rate"`   // 0-1
}

// RoutingStrategy 路由策略配置
type RoutingStrategy struct {
	StrategyID       string                 `json:"strategy_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID         string                 `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	StrategyName     string                 `json:"strategy_name" gorm:"type:varchar(100);not null;unique:idx_tenant_name"`
	Description      string                 `json:"description" gorm:"type:text"`
	Config           string                 `json:"config" gorm:"type:json;not null"` // 策略配置（JSON）
	IsActive         bool                   `json:"is_active" gorm:"default:true;index:idx_is_active"`
	Priority         int                    `json:"priority" gorm:"default:0"` // 优先级
	CreatedAt        time.Time              `json:"created_at" gorm:"not null"`
	UpdatedAt        time.Time              `json:"updated_at" gorm:"not null"`
	DeletedAt        *time.Time             `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (RoutingStrategy) TableName() string {
	return "routing_strategies"
}

// GetConfig 解析配置
func (s *RoutingStrategy) GetConfig() (map[string]interface{}, error) {
	if s.Config == "" {
		return make(map[string]interface{}), nil
	}
	var config map[string]interface{}
	// 这里需要import encoding/json
	// 简化处理，实际应在service层解析
	return config, nil
}
