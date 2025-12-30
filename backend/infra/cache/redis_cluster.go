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

package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClusterConfig Redis集群配置
type RedisClusterConfig struct {
	// 集群节点地址
	Addrs []string `toml:"addrs"` // ["localhost:7001", "localhost:7002", ...]

	// 连接配置
	MaxRetries      int           `toml:"max_retries"`
	PoolSize        int           `toml:"pool_size"`
	MinIdleConns    int           `toml:"min_idle_conns"`
	MaxIdleConns    int           `toml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `toml:"conn_max_idle_time"`

	// 只读配置
	ReadOnly         bool `toml:"read_only"`          // 允许从节点读取
	RouteByLatency   bool `toml:"route_by_latency"`   // 按延迟路由
	RouteRandomly    bool `toml:"route_randomly"`     // 随机路由

	// 超时配置
	DialTimeout  time.Duration `toml:"dial_timeout"`
	ReadTimeout  time.Duration `toml:"read_timeout"`
	WriteTimeout time.Duration `toml:"write_timeout"`

	// 密码
	Password string `toml:"password"`
}

// NewRedisCluster 创建Redis集群客户端
func NewRedisCluster(cfg *RedisClusterConfig) (*redis.ClusterClient, error) {
	if len(cfg.Addrs) == 0 {
		return nil, fmt.Errorf("redis cluster addrs cannot be empty")
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           cfg.Addrs,
		Password:        cfg.Password,
		MaxRetries:      cfg.MaxRetries,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ReadOnly:        cfg.ReadOnly,
		RouteByLatency:  cfg.RouteByLatency,
		RouteRandomly:   cfg.RouteRandomly,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,

		// 故障转移配置
		MaxRedirects: 3,

		// 命令超时
		NewClient: func(opt *redis.Options) *redis.Client {
			return redis.NewClient(opt)
		},
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping集群
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis cluster: %w", err)
	}

	// 检查集群状态
	result, err := client.ClusterInfo(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}

	if result == "" {
		return nil, fmt.Errorf("redis cluster is not initialized")
	}

	return client, nil
}

// RedisClusterClient Redis集群客户端封装
type RedisClusterClient struct {
	client *redis.ClusterClient
}

// NewRedisClusterClient 创建Redis集群客户端封装
func NewRedisClusterClient(client *redis.ClusterClient) *RedisClusterClient {
	return &RedisClusterClient{client: client}
}

// Set 设置键值
func (r *RedisClusterClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *RedisClusterClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del 删除键
func (r *RedisClusterClient) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisClusterClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *RedisClusterClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (r *RedisClusterClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// HSet 设置哈希字段
func (r *RedisClusterClient) HSet(ctx context.Context, key, field string, value interface{}) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希字段
func (r *RedisClusterClient) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (r *RedisClusterClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (r *RedisClusterClient) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// SAdd 添加到集合
func (r *RedisClusterClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合成员
func (r *RedisClusterClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// ZAdd 添加到有序集合
func (r *RedisClusterClient) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return r.client.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围
func (r *RedisClusterClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}

// Incr 自增
func (r *RedisClusterClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Decr 自减
func (r *RedisClusterClient) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// Pipeline 管道操作
func (r *RedisClusterClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// TxPipeline 事务管道
func (r *RedisClusterClient) TxPipeline() redis.Pipeliner {
	return r.client.TxPipeline()
}

// ========== 使用示例 ==========

// ExampleSetGet 设置和获取示例
func ExampleSetGet(client *RedisClusterClient) error {
	ctx := context.Background()

	// 设置键值 (自动路由到正确的分片)
	err := client.Set(ctx, "tenant:123", "tenant_data", 1*time.Hour)
	if err != nil {
		return err
	}

	// 获取值 (自动路由到主节点或从节点)
	val, err := client.Get(ctx, "tenant:123")
	if err != nil {
		return err
	}

	fmt.Printf("Value: %s\n", val)
	return nil
}

// ExampleHashOperation 哈希操作示例
func ExampleHashOperation(client *RedisClusterClient) error {
	ctx := context.Background()

	// 设置哈希字段
	err := client.HSet(ctx, "bot:456", "name", "TestBot")
	if err != nil {
		return err
	}

	// 获取哈希字段
	name, err := client.HGet(ctx, "bot:456", "name")
	if err != nil {
		return err
	}

	fmt.Printf("Bot name: %s\n", name)
	return nil
}
