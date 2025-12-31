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

package service

import (
	"crypto/md5"
	"encoding/binary"
	"hash/fnv"
	"math"
	"math/rand"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
)

// TrafficSplitter 流量分割器
type TrafficSplitter struct {
	strategies map[string]SplitStrategy
	rng        *rand.Rand
}

// SplitStrategy 流量分割策略接口
type SplitStrategy interface {
	Split(userID string, variants []Variant) *Variant
	Name() string
}

// Variant 实验变体定义
type Variant struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Weight    float64                `json:"weight"` // 流量权重 0-100
	Config    map[string]interface{} `json:"config"`
	IsControl bool                   `json:"is_control"` // 是否为对照组
}

// NewTrafficSplitter 创建流量分割器
func NewTrafficSplitter() *TrafficSplitter {
	// 使用种子初始化随机数生成器
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	splitter := &TrafficSplitter{
		strategies: make(map[string]SplitStrategy),
		rng:        rng,
	}

	// 注册内置策略
	splitter.RegisterStrategy(&WeightedSplitStrategy{rng: rng})
	splitter.RegisterStrategy(&ConsistentHashSplitStrategy{})
	splitter.RegisterStrategy(&RoundRobinSplitStrategy{})
	splitter.RegisterStrategy(&SessionBasedSplitStrategy{})

	return splitter
}

// RegisterStrategy 注册分割策略
func (s *TrafficSplitter) RegisterStrategy(strategy SplitStrategy) {
	s.strategies[strategy.Name()] = strategy
}

// GetStrategy 获取分割策略
func (s *TrafficSplitter) GetStrategy(name string) (SplitStrategy, bool) {
	strategy, exists := s.strategies[name]
	return strategy, exists
}

// Split 使用指定策略分割流量
func (s *TrafficSplitter) Split(userID string, variants []Variant, strategyName string) (*Variant, error) {
	strategy, exists := s.GetStrategy(strategyName)
	if !exists {
		// 默认使用加权随机策略
		strategy = s.strategies["weighted"]
	}

	return strategy.Split(userID, variants), nil
}

// WeightedSplitStrategy 权重分割策略（基于随机数）
type WeightedSplitStrategy struct {
	rng *rand.Rand
}

func (s *WeightedSplitStrategy) Name() string {
	return "weighted"
}

func (s *WeightedSplitStrategy) Split(userID string, variants []Variant) *Variant {
	if len(variants) == 0 {
		return nil
	}

	// 验证权重总和
	totalWeight := 0.0
	for _, v := range variants {
		totalWeight += v.Weight
	}

	// 生成0-100的随机数
	r := s.rng.Float64() * 100.0

	// 根据权重分配
	cumulative := 0.0
	for _, variant := range variants {
		// 归一化权重
		normalizedWeight := (variant.Weight / totalWeight) * 100.0
		cumulative += normalizedWeight
		if r <= cumulative {
			return &variant
		}
	}

	// 默认返回最后一个
	return &variants[len(variants)-1]
}

// ConsistentHashSplitStrategy 一致性哈希分割策略
// 确保同一用户总是分配到同一变体
type ConsistentHashSplitStrategy struct{}

func (s *ConsistentHashSplitStrategy) Name() string {
	return "consistent_hash"
}

func (s *ConsistentHashSplitStrategy) Split(userID string, variants []Variant) *Variant {
	if len(variants) == 0 {
		return nil
	}

	// 使用FNV哈希算法
	hash := fnv.New32a()
	hash.Write([]byte(userID))
	hashValue := hash.Sum32()

	// 使用哈希值选择变体
	index := int(hashValue) % len(variants)
	return &variants[index]
}

// RoundRobinSplitStrategy 轮询分割策略
type RoundRobinSplitStrategy struct {
	counter uint64
}

func (s *RoundRobinSplitStrategy) Name() string {
	return "round_robin"
}

func (s *RoundRobinSplitStrategy) Split(userID string, variants []Variant) *Variant {
	if len(variants) == 0 {
		return nil
	}

	// 原子递增计数器
	index := int(s.counter % uint64(len(variants)))
	s.counter++

	return &variants[index]
}

// SessionBasedSplitStrategy 基于会话的分割策略
// 使用用户ID的MD5哈希值进行分割
type SessionBasedSplitStrategy struct{}

func (s *SessionBasedSplitStrategy) Name() string {
	return "session_based"
}

func (s *SessionBasedSplitStrategy) Split(userID string, variants []Variant) *Variant {
	if len(variants) == 0 {
		return nil
	}

	// 计算MD5哈希
	hash := md5.Sum([]byte(userID))

	// 将前8个字节转换为uint64
	hashValue := binary.BigEndian.Uint64(hash[:8])

	// 使用哈希值选择变体
	index := int(hashValue) % len(variants)
	return &variants[index]
}

// PercentageSplitStrategy 百分比分割策略（精确控制流量分配）
type PercentageSplitStrategy struct {
	buckets []Bucket
}

// Bucket 流量桶
type Bucket struct {
	RangeStart float64
	RangeEnd   float64
	Variant    *Variant
}

func NewPercentageSplitStrategy(variants []Variant) *PercentageSplitStrategy {
	s := &PercentageSplitStrategy{
		buckets: make([]Bucket, 0, len(variants)),
	}

	// 构建流量桶
	cumulative := 0.0
	for i := range variants {
		bucket := Bucket{
			RangeStart: cumulative,
			RangeEnd:   cumulative + variants[i].Weight,
			Variant:    &variants[i],
		}
		s.buckets = append(s.buckets, bucket)
		cumulative += variants[i].Weight
	}

	return s
}

func (s *PercentageSplitStrategy) Name() string {
	return "percentage"
}

func (s *PercentageSplitStrategy) Split(userID string, variants []Variant) *Variant {
	// 生成0-100的随机数
	r := rand.Float64() * 100.0

	// 查找对应的桶
	for _, bucket := range s.buckets {
		if r >= bucket.RangeStart && r < bucket.RangeEnd {
			return bucket.Variant
		}
	}

	// 默认返回第一个
	return s.buckets[0].Variant
}

// EvaluateSplitStrategy 评估分割策略的准确性
type EvaluateSplitStrategy struct{}

// SplitEvaluation 分割评估结果
type SplitEvaluation struct {
	StrategyName   string             `json:"strategy_name"`
	TotalSamples   int                `json:"total_samples"`
	VariantCounts  map[string]int     `json:"variant_counts"`
	VariantPercent map[string]float64 `json:"variant_percent"`
	ExpectedPercent map[string]float64 `json:"expected_percent"`
	ErrorMargin    float64            `json:"error_margin"`
	IsAcceptable   bool               `json:"is_acceptable"`
}

// EvaluateSplit 评估分割策略的准确性
func (e *EvaluateSplitStrategy) EvaluateSplit(strategy SplitStrategy, variants []Variant, sampleSize int) *SplitEvaluation {
	// 模拟分割
	variantCounts := make(map[string]int)
	for i := 0; i < sampleSize; i++ {
		userID := generateTestUserID(i)
		variant := strategy.Split(userID, variants)
		variantCounts[variant.ID]++
	}

	// 计算百分比
	variantPercent := make(map[string]float64)
	expectedPercent := make(map[string]float64)

	totalWeight := 0.0
	for _, v := range variants {
		totalWeight += v.Weight
	}

	for _, v := range variants {
		variantPercent[v.ID] = float64(variantCounts[v.ID]) / float64(sampleSize) * 100.0
		expectedPercent[v.ID] = (v.Weight / totalWeight) * 100.0
	}

	// 计算误差范围
	maxError := 0.0
	for _, v := range variants {
		error := math.Abs(variantPercent[v.ID] - expectedPercent[v.ID])
		if error > maxError {
			maxError = error
		}
	}

	// 判断是否可接受（误差在5%以内）
	isAcceptable := maxError <= 5.0

	return &SplitEvaluation{
		StrategyName:    strategy.Name(),
		TotalSamples:    sampleSize,
		VariantCounts:   variantCounts,
		VariantPercent:  variantPercent,
		ExpectedPercent: expectedPercent,
		ErrorMargin:     maxError,
		IsAcceptable:    isAcceptable,
	}
}

// generateTestUserID 生成测试用户ID
func generateTestUserID(index int) string {
	return "test_user_" + string(rune(index))
}

// ConvertABTestToVariants 将AB测试实体转换为变体列表
func ConvertABTestToVariants(test *entity.ABTest) []Variant {
	// 计算权重
	weightA := float64(test.TrafficSplit)
	weightB := 100.0 - weightA

	return []Variant{
		{
			ID:        "strategy_a",
			Name:      "Strategy A",
			Weight:    weightA,
			IsControl: true, // A组作为对照组
		},
		{
			ID:        "strategy_b",
			Name:      "Strategy B",
			Weight:    weightB,
			IsControl: false,
		},
	}
}
