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

package etcd

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// EtcdClient etcd动态配置客户端
type EtcdClient struct {
	cli       *clientv3.Client
	namespace string
	group     string

	// 监听器管理
	mu        sync.RWMutex
	watchers  map[string]*watcherInfo
	cancelCtx context.CancelFunc
}

// watcherInfo 监听器信息
type watcherInfo struct {
	callback func(value string, err error)
	cancel   context.CancelFunc
}

// buildKey 构建完整的配置key
func (c *EtcdClient) buildKey(key string) string {
	// 格式：/namespace/group/key
	parts := []string{"/"}
	if c.namespace != "" {
		parts = append(parts, c.namespace)
	}
	if c.group != "" {
		parts = append(parts, c.group)
	}
	parts = append(parts, key)
	return path.Join(parts...)
}

// AddListener 添加配置变更监听器
func (c *EtcdClient) AddListener(key string, callback func(value string, err error)) error {
	fullKey := c.buildKey(key)

	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查是否已经存在监听器
	if _, exists := c.watchers[fullKey]; exists {
		return fmt.Errorf("listener already exists for key: %s", key)
	}

	// 创建带取消的上下文
	watchCtx, cancel := context.WithCancel(context.Background())

	// 启动监听goroutine
	go c.watchKey(watchCtx, fullKey, callback)

	// 保存监听器信息
	if c.watchers == nil {
		c.watchers = make(map[string]*watcherInfo)
	}
	c.watchers[fullKey] = &watcherInfo{
		callback: callback,
		cancel:   cancel,
	}

	logs.Infof("[etcd] added listener for key: %s", fullKey)
	return nil
}

// watchKey 监听key变化
func (c *EtcdClient) watchKey(ctx context.Context, fullKey string, callback func(value string, err error)) {
	watcher := c.cli.Watch(ctx, fullKey)
	for {
		select {
		case <-ctx.Done():
			logs.Infof("[etcd] watch stopped for key: %s", fullKey)
			return
		case wresp, ok := <-watcher:
			if !ok {
				// watcher channel关闭
				logs.Warnf("[etcd] watcher channel closed for key: %s", fullKey)
				return
			}
			if wresp.Err() != nil {
				// 监听出错，通知回调
				callback("", fmt.Errorf("watch error: %w", wresp.Err()))
				continue
			}

			// 处理每个事件
			for _, event := range wresp.Events {
				if event.Type == clientv3.EventTypeDelete {
					// key被删除
					callback("", fmt.Errorf("key deleted"))
				} else {
					// key被修改或创建
					callback(string(event.Kv.Value), nil)
				}
			}
		}
	}
}

// RemoveListener 移除配置变更监听器
func (c *EtcdClient) RemoveListener(key string) error {
	fullKey := c.buildKey(key)

	c.mu.Lock()
	defer c.mu.Unlock()

	watcher, exists := c.watchers[fullKey]
	if !exists {
		return fmt.Errorf("listener not found for key: %s", key)
	}

	// 取消监听
	watcher.cancel()
	delete(c.watchers, fullKey)

	logs.Infof("[etcd] removed listener for key: %s", fullKey)
	return nil
}

// Get 获取配置值
func (c *EtcdClient) Get(ctx context.Context, key string) (string, error) {
	fullKey := c.buildKey(key)

	resp, err := c.cli.Get(ctx, fullKey)
	if err != nil {
		return "", fmt.Errorf("get from etcd failed: %w", err)
	}

	if resp.Count == 0 {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return string(resp.Kvs[0].Value), nil
}

// GetWithPrefix 获取指定前缀的所有配置
func (c *EtcdClient) GetWithPrefix(ctx context.Context, prefix string) (map[string]string, error) {
	fullPrefix := c.buildKey(prefix)

	resp, err := c.cli.Get(ctx, fullPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("get with prefix from etcd failed: %w", err)
	}

	result := make(map[string]string)
	for _, kv := range resp.Kvs {
		// 去掉前缀，返回相对路径
		relKey := strings.TrimPrefix(string(kv.Key), fullPrefix)
		relKey = strings.TrimPrefix(relKey, "/")
		if relKey == "" {
			relKey = prefix
		}
		result[relKey] = string(kv.Value)
	}

	return result, nil
}

// Set 设置配置值
func (c *EtcdClient) Set(ctx context.Context, key string, value string) error {
	fullKey := c.buildKey(key)

	_, err := c.cli.Put(ctx, fullKey, value)
	if err != nil {
		return fmt.Errorf("put to etcd failed: %w", err)
	}

	logs.Debugf("[etcd] set key: %s", fullKey)
	return nil
}

// Delete 删除配置
func (c *EtcdClient) Delete(ctx context.Context, key string) error {
	fullKey := c.buildKey(key)

	_, err := c.cli.Delete(ctx, fullKey)
	if err != nil {
		return fmt.Errorf("delete from etcd failed: %w", err)
	}

	logs.Infof("[etcd] deleted key: %s", fullKey)
	return nil
}

// Close 关闭客户端
func (c *EtcdClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 取消所有监听器
	for key, watcher := range c.watchers {
		watcher.cancel()
		delete(c.watchers, key)
		logs.Infof("[etcd] removed listener for key: %s", key)
	}

	// 关闭etcd客户端
	if c.cli != nil {
		return c.cli.Close()
	}

	return nil
}
