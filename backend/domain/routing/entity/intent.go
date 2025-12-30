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
	"encoding/json"
	"time"
)

// Intent 意图实体
type Intent struct {
	IntentID    string    `json:"intent_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID    string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id,priority:1"`
	IntentName  string    `json:"intent_name" gorm:"type:varchar(100);not null;index:idx_tenant_name,priority:2"`
	Description string    `json:"description" gorm:"type:text"`
	Examples    string    `json:"examples" gorm:"type:json"` // JSON数组：示例句子
	AgentID     string    `json:"agent_id" gorm:"type:varchar(36);not null;index:idx_agent_id"` // 路由目标
	WorkflowID  string    `json:"workflow_id" gorm:"type:varchar(36);index:idx_workflow_id"`     // 路由目标（可选）
	Confidence  float64   `json:"confidence" gorm:"type:decimal(5,4);not null;default:0.8000"` // 置信度阈值
	IsActive    bool      `json:"is_active" gorm:"default:true;index:idx_is_active"`
	CreatedAt   int64     `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt   int64     `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt   *int64    `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Intent) TableName() string {
	return "intents"
}

// IsIntentActive 意图是否激活
func (i *Intent) IsIntentActive() bool {
	return i.IsActive && i.DeletedAt == nil
}

// GetExamples 解析示例句子
func (i *Intent) GetExamples() ([]string, error) {
	if i.Examples == "" {
		return []string{}, nil
	}
	var examples []string
	err := json.Unmarshal([]byte(i.Examples), &examples)
	return examples, err
}

// SetExamples 设置示例句子
func (i *Intent) SetExamples(examples []string) error {
	data, err := json.Marshal(examples)
	if err != nil {
		return err
	}
	i.Examples = string(data)
	return nil
}

// GetTargetType 获取目标类型（agent或workflow）
func (i *Intent) GetTargetType() string {
	if i.WorkflowID != "" {
		return "workflow"
	}
	return "agent"
}

// IntentSample 意图样本实体（用于训练）
type IntentSample struct {
	SampleID   string    `json:"sample_id" gorm:"primaryKey;type:varchar(36)"`
	IntentID   string    `json:"intent_id" gorm:"type:varchar(36);not null;index:idx_intent_id"`
	Text       string    `json:"text" gorm:"type:text;not null"` // 训练样本文本
	Embedding  string    `json:"embedding" gorm:"type:json"`     // 向量嵌入（JSON数组）
	CreateTime int64     `json:"create_time" gorm:"not null;default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
}

// TableName 指定表名
func (IntentSample) TableName() string {
	return "intent_samples"
}

// GetEmbedding 解析向量嵌入
func (s *IntentSample) GetEmbedding() ([]float64, error) {
	if s.Embedding == "" {
		return []float64{}, nil
	}
	var embedding []float64
	err := json.Unmarshal([]byte(s.Embedding), &embedding)
	return embedding, err
}

// SetEmbedding 设置向量嵌入
func (s *IntentSample) SetEmbedding(embedding []float64) error {
	data, err := json.Marshal(embedding)
	if err != nil {
		return err
	}
	s.Embedding = string(data)
	return nil
}

// Entity 实体提取结果
type Entity struct {
	EntityName  string `json:"entity_name"`  // 实体名称（如location、action、target）
	EntityValue string `json:"entity_value"` // 实体值（如"北京"、"查询"、"天气"）
	Confidence  float64 `json:"confidence"`  // 置信度
	StartPosition int  `json:"start_position"` // 在文本中的起始位置
	EndPosition   int  `json:"end_position"`   // 在文本中的结束位置
}

// IntentRecognitionResult 意图识别结果
type IntentRecognitionResult struct {
	IntentID    string    `json:"intent_id"`
	IntentName  string    `json:"intent_name"`
	Confidence  float64   `json:"confidence"`
	AgentID     string    `json:"agent_id"`
	WorkflowID  string    `json:"workflow_id"`
	Entities    []*Entity `json:"entities"`
	MatchMethod string    `json:"match_method"` // llm, vector, rule, hybrid
}

// RoutingOptimizeLog 路由优化日志（用于学习）
type RoutingOptimizeLog struct {
	LogID          string  `json:"log_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID       string  `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID         string  `json:"user_id" gorm:"type:varchar(36);not null"`
	Query          string  `json:"query" gorm:"type:text;not null"`
	Intent         string  `json:"intent" gorm:"type:varchar(100);index:idx_intent"`
	AgentID        string  `json:"agent_id" gorm:"type:varchar(36);not null;index:idx_agent_id"`
	Strategy       string  `json:"strategy" gorm:"type:varchar(100);not null"` // 使用的路由策略
	ResponseTime   int     `json:"response_time" gorm:"not null"`             // ms
	UserRating     int     `json:"user_rating"`                                // 1-5
	Feedback       string  `json:"feedback" gorm:"type:text"`                 // 用户反馈
	CreatedAt      int64   `json:"created_at" gorm:"not null;index:idx_created_at"`
	DeletedAt      *int64  `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (RoutingOptimizeLog) TableName() string {
	return "routing_optimize_logs"
}

// RoutingDecision 路由决策（增强版）
type RoutingDecision struct {
	AgentID      string     `json:"agent_id"`
	WorkflowID   string     `json:"workflow_id"`
	Confidence   float64    `json:"confidence"`
	Score        float64    `json:"score"`
	Reasons      []string   `json:"reasons"`
	MatchType    string     `json:"match_type"`
	RuleID       string     `json:"rule_id"`
	IntentID     string     `json:"intent_id"`
	Entities     []*Entity  `json:"entities"`
	Strategy     string     `json:"strategy"`     // 使用的路由策略
	PredictedLoad float64   `json:"predicted_load"` // 预测负载
}

// IsDecisionValid 决策是否有效
func (d *RoutingDecision) IsDecisionValid() bool {
	return d.AgentID != "" || d.WorkflowID != ""
}

// GetTargetID 获取目标ID
func (d *RoutingDecision) GetTargetID() string {
	if d.AgentID != "" {
		return d.AgentID
	}
	return d.WorkflowID
}

// GetTargetType 获取目标类型
func (d *RoutingDecision) GetTargetType() string {
	if d.AgentID != "" {
		return "agent"
	}
	return "workflow"
}
