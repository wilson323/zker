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
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// SearchCache 搜索缓存
// 使用LRU策略缓存搜索结果,提升重复查询的性能
type SearchCache struct {
	mu       sync.RWMutex
	cache    map[string]*cacheEntry
	lruList  *lruList
	maxSize  int
	ttl      time.Duration
	hits     int64
	misses   int64
	evictions int64
}

// cacheEntry 缓存条目
type cacheEntry struct {
	key        string
	value      *SearchResult
	expireTime time.Time
	prev, next *cacheEntry
}

// lruList LRU双向链表
type lruList struct {
	head, tail *cacheEntry
	size       int
}

// NewSearchCache 创建搜索缓存
func NewSearchCache(maxSize int, ttl time.Duration) *SearchCache {
	if maxSize <= 0 {
		maxSize = 1000 // 默认缓存1000条
	}
	if ttl <= 0 {
		ttl = 5 * time.Minute // 默认TTL 5分钟
	}

	cache := &SearchCache{
		cache:   make(map[string]*cacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
		lruList: &lruList{},
	}

	// 初始化LRU链表
	cache.lruList.head = &cacheEntry{}
	cache.lruList.tail = &cacheEntry{}
	cache.lruList.head.next = cache.lruList.tail
	cache.lruList.tail.prev = cache.lruList.head

	return cache
}

// Get 获取缓存
func (c *SearchCache) Get(ctx context.Context, query string, knowledgeID int64, topK int) *SearchResult {
	key := c.buildCacheKey(query, knowledgeID, topK)

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cache[key]
	if !exists {
		c.misses++
		return nil
	}

	// 检查是否过期
	if time.Now().After(entry.expireTime) {
		c.removeEntry(entry)
		c.misses++
		return nil
	}

	// 更新LRU:移到链表头部
	c.lruList.moveToFront(entry)

	c.hits++
	logs.CtxInfof(ctx, "[SearchCache] cache hit for query: %s", query)

	return entry.value
}

// Set 设置缓存
func (c *SearchCache) Set(ctx context.Context, query string, knowledgeID int64, topK int, result *SearchResult) {
	key := c.buildCacheKey(query, knowledgeID, topK)

	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果已存在,更新值并移到头部
	if entry, exists := c.cache[key]; exists {
		entry.value = result
		entry.expireTime = time.Now().Add(c.ttl)
		c.lruList.moveToFront(entry)
		return
	}

	// 检查缓存是否已满
	if c.lruList.size >= c.maxSize {
		// 淘汰最久未使用的条目
		c.evictLRU()
	}

	// 创建新条目
	entry := &cacheEntry{
		key:        key,
		value:      result,
		expireTime: time.Now().Add(c.ttl),
	}

	// 添加到缓存和LRU链表头部
	c.cache[key] = entry
	c.lruList.addToFront(entry)

	logs.CtxInfof(ctx, "[SearchCache] cached result for query: %s (size=%d)", query, c.lruList.size)
}

// Delete 删除缓存
func (c *SearchCache) Delete(ctx context.Context, query string, knowledgeID int64, topK int) {
	key := c.buildCacheKey(query, knowledgeID, topK)

	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, exists := c.cache[key]; exists {
		c.removeEntry(entry)
		logs.CtxInfof(ctx, "[SearchCache] deleted cache for query: %s", query)
	}
}

// Clear 清空缓存
func (c *SearchCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 清空map
	c.cache = make(map[string]*cacheEntry)

	// 重置LRU链表
	c.lruList.head.next = c.lruList.tail
	c.lruList.tail.prev = c.lruList.head
	c.lruList.size = 0

	logs.CtxInfof(context.Background(), "[SearchCache] cache cleared")
}

// GetStats 获取缓存统计信息
func (c *SearchCache) GetStats() *CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(c.hits) / float64(total)
	}

	return &CacheStats{
		Size:      c.lruList.size,
		MaxSize:   c.maxSize,
		Hits:      c.hits,
		Misses:    c.misses,
		HitRate:   hitRate,
		Evictions: c.evictions,
		TTL:       c.ttl,
	}
}

// buildCacheKey 构建缓存key
func (c *SearchCache) buildCacheKey(query string, knowledgeID int64, topK int) string {
	// 使用SHA256哈希生成唯一key
	data := fmt.Sprintf("%s:%d:%d", query, knowledgeID, topK)
	hash := sha256.Sum256([]byte(data))
	// 取前16字节编码为base64,缩短key长度
	return base64.URLEncoding.EncodeToString(hash[:16])
}

// evictLRU 淘汰最久未使用的条目
func (c *SearchCache) evictLRU() {
	// 获取链表尾部(最久未使用)
	lru := c.lruList.tail.prev
	if lru == c.lruList.head {
		return
	}

	c.removeEntry(lru)
	c.evictions++
}

// removeEntry 移除条目
func (c *SearchCache) removeEntry(entry *cacheEntry) {
	if entry == nil {
		return
	}

	// 从map中删除
	delete(c.cache, entry.key)

	// 从LRU链表中删除
	c.lruList.remove(entry)
}

// LRU链表操作
func (l *lruList) addToFront(entry *cacheEntry) {
	entry.prev = l.head
	entry.next = l.head.next
	l.head.next.prev = entry
	l.head.next = entry
	l.size++
}

func (l *lruList) remove(entry *cacheEntry) {
	entry.prev.next = entry.next
	entry.next.prev = entry.prev
	l.size--
}

func (l *lruList) moveToFront(entry *cacheEntry) {
	l.remove(entry)
	l.addToFront(entry)
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Size      int           // 当前缓存条目数
	MaxSize   int           // 最大缓存条目数
	Hits      int64         // 命中次数
	Misses    int64         // 未命中次数
	HitRate   float64       // 命中率
	Evictions int64         // 淘汰次数
	TTL       time.Duration // 过期时间
}

// String 返回统计信息的字符串表示
func (s *CacheStats) String() string {
	statsJSON, _ := json.Marshal(map[string]interface{}{
		"size":       s.Size,
		"max_size":   s.MaxSize,
		"hits":       s.Hits,
		"misses":     s.Misses,
		"hit_rate":   fmt.Sprintf("%.2f%%", s.HitRate*100),
		"evictions":  s.Evictions,
		"ttl_seconds": s.TTL.Seconds(),
	})
	return string(statsJSON)
}

// CachedSearchEngine 带缓存的搜索引擎
type CachedSearchEngine struct {
	engine *HybridSearchService
	cache  *SearchCache
}

// NewCachedSearchEngine 创建带缓存的搜索引擎
func NewCachedSearchEngine(engine *HybridSearchService, cache *SearchCache) *CachedSearchEngine {
	return &CachedSearchEngine{
		engine: engine,
		cache:  cache,
	}
}

// Search 执行搜索(带缓存)
func (e *CachedSearchEngine) Search(ctx context.Context, req *SearchRequest) (*SearchResult, error) {
	// 尝试从缓存获取
	if req.EnableCache && e.cache != nil {
		if cached := e.cache.Get(ctx, req.Query, req.KnowledgeID, req.TopK); cached != nil {
			cached.SearchMeta.CacheHit = true
			logs.CtxInfof(ctx, "[CachedSearch] cache hit for query: %s", req.Query)
			return cached, nil
		}
	}

	// 执行搜索
	result, err := e.engine.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	if req.EnableCache && e.cache != nil {
		e.cache.Set(ctx, req.Query, req.KnowledgeID, req.TopK, result)
	}

	return result, nil
}

// InvalidateCache 使缓存失效
func (e *CachedSearchEngine) InvalidateCache(ctx context.Context, query string, knowledgeID int64, topK int) {
	if e.cache != nil {
		e.cache.Delete(ctx, query, knowledgeID, topK)
	}
}

// ClearCache 清空缓存
func (e *CachedSearchEngine) ClearCache() {
	if e.cache != nil {
		e.cache.Clear()
	}
}

// GetCacheStats 获取缓存统计信息
func (e *CachedSearchEngine) GetCacheStats() *CacheStats {
	if e.cache != nil {
		return e.cache.GetStats()
	}
	return nil
}
