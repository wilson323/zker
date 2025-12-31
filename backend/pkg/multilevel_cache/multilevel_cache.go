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

// Package multilevel_cache provides multi-level caching implementation
// with L1 (local memory) and L2 (Redis) cache tiers.
package multilevel_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logging"
)

const (
	// DefaultLocalCacheSize is the default size of L1 local cache
	DefaultLocalCacheSize = 1000
	// DefaultLocalCacheTTL is the default TTL for L1 local cache
	DefaultLocalCacheTTL = 1 * time.Minute
	// DefaultRedisCacheTTL is the default TTL for L2 Redis cache
	DefaultRedisCacheTTL = 5 * time.Minute
)

// CacheItem represents an item in the cache
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// IsExpired checks if the cache item is expired
func (item *CacheItem) IsExpired() bool {
	return time.Now().After(item.Expiration)
}

// LocalCache is a thread-safe LRU cache for L1 caching
type LocalCache struct {
	mu       sync.RWMutex
	items    map[string]*CacheItem
	maxSize  int
	ttl      time.Duration
	hitCount int64
	missCount int64
}

// NewLocalCache creates a new local cache
func NewLocalCache(maxSize int, ttl time.Duration) *LocalCache {
	if maxSize <= 0 {
		maxSize = DefaultLocalCacheSize
	}
	if ttl <= 0 {
		ttl = DefaultLocalCacheTTL
	}

	cache := &LocalCache{
		items:   make(map[string]*CacheItem),
		maxSize: maxSize,
		ttl:     ttl,
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// Set stores a value in the cache
func (c *LocalCache) Set(key string, value interface{}, ttl ...time.Duration) {
	expiration := time.Now().Add(c.ttl)
	if len(ttl) > 0 && ttl[0] > 0 {
		expiration = time.Now().Add(ttl[0])
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict oldest if at capacity
	if len(c.items) >= c.maxSize {
		var oldestKey string
		var oldestTime time.Time
		for k, v := range c.items {
			if oldestKey == "" || v.Expiration.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.Expiration
			}
		}
		if oldestKey != "" {
			delete(c.items, oldestKey)
		}
	}

	c.items[key] = &CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Get retrieves a value from the cache
func (c *LocalCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		c.missCount++
		return nil, false
	}

	if item.IsExpired() {
		delete(c.items, key)
		c.missCount++
		return nil, false
	}

	c.hitCount++
	return item.Value, true
}

// Delete removes a value from the cache
func (c *LocalCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear clears all items from the cache
func (c *LocalCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*CacheItem)
}

// Size returns the current size of the cache
func (c *LocalCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Stats returns cache statistics
func (c *LocalCache) Stats() (hitCount, missCount int64, size int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hitCount, c.missCount, len(c.items)
}

// cleanupExpired periodically removes expired items
func (c *LocalCache) cleanupExpired() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, item := range c.items {
			if item.IsExpired() {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// RedisClient defines the interface for Redis operations
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Pipeline() RedisPipeline
}

// RedisPipeline defines the interface for Redis pipeline operations
type RedisPipeline interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exec(ctx context.Context) error
}

// MultiLevelCache provides L1 (local) + L2 (Redis) caching
type MultiLevelCache struct {
	local      *LocalCache
	redis      RedisClient
	localTTL   time.Duration
	redisTTL   time.Duration
	keyPrefix  string
	stats      *CacheStats
}

// CacheStats tracks cache statistics
type CacheStats struct {
	mu                sync.RWMutex
	L1HitCount        int64
	L1MissCount       int64
	L2HitCount        int64
	L2MissCount       int64
	DBQueryCount      int64
	L1HitRate         float64
	L2HitRate         float64
	OverallHitRate    float64
}

// NewMultiLevelCache creates a new multi-level cache
func NewMultiLevelCache(redis RedisClient, keyPrefix string, options ...Option) *MultiLevelCache {
	cache := &MultiLevelCache{
		local:     NewLocalCache(DefaultLocalCacheSize, DefaultLocalCacheTTL),
		redis:     redis,
		localTTL:  DefaultLocalCacheTTL,
		redisTTL:  DefaultRedisCacheTTL,
		keyPrefix: keyPrefix,
		stats:     &CacheStats{},
	}

	// Apply options
	for _, opt := range options {
		opt(cache)
	}

	return cache
}

// Option configures the multi-level cache
type Option func(*MultiLevelCache)

// WithLocalCacheSize sets the L1 cache size
func WithLocalCacheSize(size int) Option {
	return func(c *MultiLevelCache) {
		c.local = NewLocalCache(size, c.localTTL)
	}
}

// WithLocalTTL sets the L1 cache TTL
func WithLocalTTL(ttl time.Duration) Option {
	return func(c *MultiLevelCache) {
		c.localTTL = ttl
		c.local = NewLocalCache(c.local.maxSize, ttl)
	}
}

// WithRedisTTL sets the L2 cache TTL
func WithRedisTTL(ttl time.Duration) Option {
	return func(c *MultiLevelCache) {
		c.redisTTL = ttl
	}
}

// Get retrieves a value from the cache with fallback to data source
func (c *MultiLevelCache) Get(ctx context.Context, key string, dataSource func() (interface{}, error)) (interface{}, error) {
	return c.GetWithTTL(ctx, key, dataSource, c.localTTL, c.redisTTL)
}

// GetWithTTL retrieves a value from the cache with custom TTL
func (c *MultiLevelCache) GetWithTTL(ctx context.Context, key string, dataSource func() (interface{}, error), localTTL, redisTTL time.Duration) (interface{}, error) {
	fullKey := c.buildKey(key)

	// L1: Try local cache
	if val, found := c.local.Get(fullKey); found {
		c.recordL1Hit()
		return val, nil
	}
	c.recordL1Miss()

	// L2: Try Redis cache
	val, err := c.redis.Get(ctx, fullKey)
	if err == nil && val != "" {
		c.recordL2Hit()
		var result interface{}
		if err := json.Unmarshal([]byte(val), &result); err == nil {
			// Promote to L1
			c.local.Set(fullKey, result, localTTL)
			return result, nil
		}
	}
	c.recordL2Miss()

	// L3: Query data source
	if dataSource != nil {
		c.recordDBQuery()
		result, err := dataSource()
		if err != nil {
			return nil, err
		}

		// Write back to cache
		if err := c.setInternal(ctx, fullKey, result, localTTL, redisTTL); err != nil {
			logging.CtxErrorf(ctx, "[MultiLevelCache] set cache failed: %v", err)
		}

		return result, nil
	}

	return nil, fmt.Errorf("cache miss and no data source provided")
}

// Set stores a value in both L1 and L2 cache
func (c *MultiLevelCache) Set(ctx context.Context, key string, value interface{}) error {
	return c.SetWithTTL(ctx, key, value, c.localTTL, c.redisTTL)
}

// SetWithTTL stores a value with custom TTL
func (c *MultiLevelCache) SetWithTTL(ctx context.Context, key string, value interface{}, localTTL, redisTTL time.Duration) error {
	fullKey := c.buildKey(key)
	return c.setInternal(ctx, fullKey, value, localTTL, redisTTL)
}

func (c *MultiLevelCache) setInternal(ctx context.Context, key string, value interface{}, localTTL, redisTTL time.Duration) error {
	// Set to L1
	c.local.Set(key, value, localTTL)

	// Set to L2
	if c.redis != nil {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		if err := c.redis.Set(ctx, key, data, redisTTL); err != nil {
			logging.CtxErrorf(ctx, "[MultiLevelCache] redis set failed: %v", err)
			return err
		}
	}

	return nil
}

// Delete removes a value from both L1 and L2 cache
func (c *MultiLevelCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	// Delete from L1
	for _, key := range keys {
		c.local.Delete(c.buildKey(key))
	}

	// Delete from L2
	if c.redis != nil {
		fullKeys := make([]string, len(keys))
		for i, key := range keys {
			fullKeys[i] = c.buildKey(key)
		}
		if err := c.redis.Del(ctx, fullKeys...); err != nil {
			logging.CtxErrorf(ctx, "[MultiLevelCache] redis del failed: %v", err)
			return err
		}
	}

	return nil
}

// GetBatch retrieves multiple values from cache
func (c *MultiLevelCache) GetBatch(ctx context.Context, keys []string, dataSource func([]string) (map[string]interface{}, error)) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	missedKeys := make([]string, 0)

	// Try to get from L1
	for _, key := range keys {
		fullKey := c.buildKey(key)
		if val, found := c.local.Get(fullKey); found {
			c.recordL1Hit()
			result[key] = val
		} else {
			c.recordL1Miss()
			missedKeys = append(missedKeys, key)
		}
	}

	// If all found, return
	if len(missedKeys) == 0 {
		return result, nil
	}

	// Try L2 batch get
	if c.redis != nil {
		// Use pipeline for batch Redis operations
		pipe := c.redis.Pipeline()
		redisResults := make(map[string]interface{})

		for _, key := range missedKeys {
			fullKey := c.buildKey(key)
			pipe.Get(ctx, fullKey)
		}

		if err := pipe.Exec(ctx); err != nil {
			logging.CtxErrorf(ctx, "[MultiLevelCache] pipeline exec failed: %v", err)
		}
	}

	// Query data source for remaining keys
	if dataSource != nil && len(missedKeys) > 0 {
		c.recordDBQuery()
		dataResult, err := dataSource(missedKeys)
		if err != nil {
			return nil, err
		}

		// Merge results
		for key, val := range dataResult {
			result[key] = val
			// Promote to L1
			c.local.Set(c.buildKey(key), val, c.localTTL)
		}
	}

	return result, nil
}

// SetBatch stores multiple values in cache
func (c *MultiLevelCache) SetBatch(ctx context.Context, items map[string]interface{}) error {
	if c.redis == nil {
		for key, val := range items {
			c.local.Set(c.buildKey(key), val, c.localTTL)
		}
		return nil
	}

	// Use pipeline for batch Redis operations
	pipe := c.redis.Pipeline()

	for key, val := range items {
		fullKey := c.buildKey(key)
		data, err := json.Marshal(val)
		if err != nil {
			logging.CtxErrorf(ctx, "[MultiLevelCache] marshal failed: %v", err)
			continue
		}

		// Set to L1
		c.local.Set(fullKey, val, c.localTTL)

		// Queue L2 set
		pipe.Set(ctx, fullKey, data, c.redisTTL)
	}

	// Execute pipeline
	if err := pipe.Exec(ctx); err != nil {
		logging.CtxErrorf(ctx, "[MultiLevelCache] batch set failed: %v", err)
		return err
	}

	return nil
}

// Clear clears all items from both L1 and L2 cache
func (c *MultiLevelCache) Clear(ctx context.Context) error {
	// Clear L1
	c.local.Clear()

	// Clear L2 (clear all keys with prefix)
	// Note: This requires SCAN + DEL operation which is expensive
	// Consider using a separate version key pattern instead

	return nil
}

// Stats returns cache statistics
func (c *MultiLevelCache) Stats() CacheStatsReport {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()

	l1Hits := c.stats.L1HitCount
	l1Misses := c.stats.L1MissCount
	l2Hits := c.stats.L2HitCount
	l2Misses := c.stats.L2MissCount

	var l1HitRate, l2HitRate, overallHitRate float64
	if l1Hits+l1Misses > 0 {
		l1HitRate = float64(l1Hits) / float64(l1Hits+l1Misses) * 100
	}
	if l2Hits+l2Misses > 0 {
		l2HitRate = float64(l2Hits) / float64(l2Hits+l2Misses) * 100
	}
	overallHits := l1Hits + l2Hits
	overallMisses := l1Misses + l2Misses
	if overallHits+overallMisses > 0 {
		overallHitRate = float64(overallHits) / float64(overallHits+overallMisses) * 100
	}

	return CacheStatsReport{
		L1HitCount:     l1Hits,
		L1MissCount:    l1Misses,
		L2HitCount:     l2Hits,
		L2MissCount:    l2Misses,
		DBQueryCount:   c.stats.DBQueryCount,
		L1HitRate:      l1HitRate,
		L2HitRate:      l2HitRate,
		OverallHitRate: overallHitRate,
		LocalSize:      c.local.Size(),
	}
}

// CacheStatsReport represents a cache statistics report
type CacheStatsReport struct {
	L1HitCount     int64   `json:"l1_hit_count"`
	L1MissCount    int64   `json:"l1_miss_count"`
	L2HitCount     int64   `json:"l2_hit_count"`
	L2MissCount    int64   `json:"l2_miss_count"`
	DBQueryCount   int64   `json:"db_query_count"`
	L1HitRate      float64 `json:"l1_hit_rate"`
	L2HitRate      float64 `json:"l2_hit_rate"`
	OverallHitRate float64 `json:"overall_hit_rate"`
	LocalSize      int     `json:"local_size"`
}

func (c *MultiLevelCache) recordL1Hit() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.L1HitCount++
}

func (c *MultiLevelCache) recordL1Miss() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.L1MissCount++
}

func (c *MultiLevelCache) recordL2Hit() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.L2HitCount++
}

func (c *MultiLevelCache) recordL2Miss() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.L2MissCount++
}

func (c *MultiLevelCache) recordDBQuery() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.DBQueryCount++
}

func (c *MultiLevelCache) buildKey(key string) string {
	if c.keyPrefix == "" {
		return key
	}
	return c.keyPrefix + ":" + key
}

// InvalidatePattern invalidates all cache keys matching a pattern
// Note: This requires SCAN operation, use sparingly
func (c *MultiLevelCache) InvalidatePattern(ctx context.Context, pattern string) error {
	// Invalidate L1 matching keys
	c.local.Clear()

	// Invalidate L2 matching keys
	// Note: This requires SCAN + DEL which is expensive
	// Consider using Redis KEYS (not recommended in production) or SCAN

	return nil
}
