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
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/coze-dev/coze-studio/backend/infra/dynconf"
)

// Provider etcd配置提供者
type Provider struct {
	endpoints []string
	// 可选：TLS配置
	tlsEnabled bool
	// 可选：认证信息
	username string
	password string
}

// NewProvider 创建etcd provider
func NewProvider(endpoints []string) *Provider {
	return &Provider{
		endpoints: endpoints,
	}
}

// WithTLS 启用TLS
func (p *Provider) WithTLS() *Provider {
	p.tlsEnabled = true
	return p
}

// WithAuth 设置认证信息
func (p *Provider) WithAuth(username, password string) *Provider {
	p.username = username
	p.password = password
	return p
}

// Initialize 初始化etcd客户端
func (p *Provider) Initialize(ctx context.Context, namespace, group string, opts ...dynconf.Option) (dynconf.DynamicClient, error) {
	// 构建etcd客户端配置
	cliCfg := clientv3.Config{
		Endpoints:   p.endpoints,
		DialTimeout: 5 * time.Second,
	}

	// TLS配置（如果启用）
	if p.tlsEnabled {
		// TODO: 添加TLS配置
		// cliCfg.TLS = &tls.Config{...}
	}

	// 认证配置
	if p.username != "" && p.password != "" {
		cliCfg.Username = p.username
		cliCfg.Password = p.password
	}

	// 创建etcd客户端
	cli, err := clientv3.New(cliCfg)
	if err != nil {
		return nil, fmt.Errorf("create etcd client failed: %w", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	_, err = cli.Get(ctx, "test")
	if err != nil && err != context.DeadlineExceeded {
		// 忽略这个特定的错误，因为key不存在是正常的
		cli.Close()
		return nil, fmt.Errorf("connect to etcd failed: %w", err)
	}

	// 创建动态客户端
	client := &EtcdClient{
		cli:       cli,
		namespace: namespace,
		group:     group,
	}

	return client, nil
}
