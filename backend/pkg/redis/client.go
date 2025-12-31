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

// Package redis provides a unified Redis client wrapper for the application
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client is a unified Redis client wrapper
type Client struct {
	client *redis.Client
}

// NewClient creates a new Redis client
func NewClient(addr, password string, db int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 100,
	})

	return &Client{
		client: rdb,
	}
}

// Close closes the Redis client
func (c *Client) Close() error {
	return c.client.Close()
}

// Set stores a value in Redis with expiration
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value from Redis
func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("redis: key not found")
		}
		return fmt.Errorf("failed to get from redis: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Del deletes a key from Redis
func (c *Client) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	return count > 0, err
}

// Expire sets an expiration time for a key
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL returns the time to live for a key
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// Incr increments a counter
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// Decr decrements a counter
func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

// HSet sets a field in a hash
func (c *Client) HSet(ctx context.Context, key, field string, value interface{}) error {
	return c.client.HSet(ctx, key, field, value).Err()
}

// HGet gets a field from a hash
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return c.client.HGet(ctx, key, field).Result()
}

// HGetAll gets all fields from a hash
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HDel deletes a field from a hash
func (c *Client) HDel(ctx context.Context, key, field string) error {
	return c.client.HDel(ctx, key, field).Err()
}

// ZAdd adds a member to a sorted set
func (c *Client) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return c.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// ZRange retrieves a range of members from a sorted set
func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRem removes a member from a sorted set
func (c *Client) ZRem(ctx context.Context, key, member string) error {
	return c.client.ZRem(ctx, key, member).Err()
}

// Pipeline returns a Redis pipeline
func (c *Client) Pipeline() redis.Pipeliner {
	return c.client.Pipeline()
}

// GetRawClient returns the underlying redis.Client for advanced usage
func (c *Client) GetRawClient() *redis.Client {
	return c.client
}
