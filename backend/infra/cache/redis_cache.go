// backend/infra/cache/redis_cache.go
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/redis/go-redis/v9"
)

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(addr, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 100,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logs.Errorf("Failed to connect to Redis: %v", err)
		panic(err)
	}

	logs.Infof("Redis cache initialized: addr=%s, db=%d, ttl=%s", addr, db, ttl)

	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

// Set 设置缓存
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.buildKey(key), data, c.ttl).Err()
}

// Get 获取缓存
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, c.buildKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheNotFound
		}
		return err
	}

	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.buildKey(key)).Err()
}

// SetBatch 批量设置缓存
func (c *RedisCache) SetBatch(ctx context.Context, items map[string]interface{}) error {
	pipe := c.client.Pipeline()

	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			logs.Errorf("Failed to marshal cache value: key=%s, error=%v", key, err)
			continue
		}

		pipe.Set(ctx, c.buildKey(key), data, c.ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetBatch 批量获取缓存
func (c *RedisCache) GetBatch(ctx context.Context, keys []string) (map[string]interface{}, error) {
	pipe := c.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))

	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, c.buildKey(key))
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	results := make(map[string]interface{})
	for i, cmd := range cmds {
		var value interface{}
		if err := json.Unmarshal([]byte(cmd.Val()), &value); err == nil {
			results[keys[i]] = value
		}
	}

	return results, nil
}

// DeleteByPattern 根据模式删除缓存
func (c *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, c.buildKey(pattern), 0).Iterator()
	keys := []string{}

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists 检查缓存是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) bool {
	count, err := c.client.Exists(ctx, c.buildKey(key)).Result()
	if err != nil {
		return false
	}
	return count > 0
}

// buildKey 构建缓存键
func (c *RedisCache) buildKey(key string) string {
	return "zker:" + key
}

// Close 关闭连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// =====================================================================
// 错误定义
// =====================================================================

var (
	ErrCacheNotFound = errors.New("cache not found")
)

// =====================================================================
// Bot缓存示例（临时类型定义，实际应从domain/entity导入）
// =====================================================================

// Bot Bot实体（临时定义）
type Bot struct {
	BotID   string `json:"bot_id"`
	Name    string `json:"name"`
	TenantID string `json:"tenant_id"`
}

// BotCache Bot缓存
type BotCache struct {
	cache *RedisCache
}

// NewBotCache 创建Bot缓存
func NewBotCache(redisCache *RedisCache) *BotCache {
	return &BotCache{cache: redisCache}
}

// GetBot 获取Bot缓存
func (c *BotCache) GetBot(ctx context.Context, botID string) (*Bot, error) {
	var bot Bot
	err := c.cache.Get(ctx, fmt.Sprintf("bot:%s", botID), &bot)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// SetBot 设置Bot缓存
func (c *BotCache) SetBot(ctx context.Context, bot *Bot) error {
	return c.cache.Set(ctx, fmt.Sprintf("bot:%s", bot.BotID), bot)
}

// GetBotList 获取Bot列表缓存
func (c *BotCache) GetBotList(ctx context.Context, tenantID string, page int) ([]*Bot, error) {
	key := fmt.Sprintf("bots:%s:page:%d", tenantID, page)
	var bots []*Bot
	err := c.cache.Get(ctx, key, &bots)
	if err != nil {
		return nil, err
	}
	return bots, nil
}

// SetBotList 设置Bot列表缓存
func (c *BotCache) SetBotList(ctx context.Context, tenantID string, page int, bots []*Bot) error {
	key := fmt.Sprintf("bots:%s:page:%d", tenantID, page)
	return c.cache.Set(ctx, key, bots)
}

// InvalidateBot 使Bot缓存失效
func (c *BotCache) InvalidateBot(ctx context.Context, botID string) error {
	return c.cache.Delete(ctx, fmt.Sprintf("bot:%s", botID))
}

// InvalidateBotList 使Bot列表缓存失效
func (c *BotCache) InvalidateBotList(ctx context.Context, tenantID string) error {
	pattern := fmt.Sprintf("bots:%s:*", tenantID)
	return c.cache.DeleteByPattern(ctx, pattern)
}
