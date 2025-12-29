// Consul服务发现
package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	ID                string            `json:"id"`                 // 服务实例ID
	Name              string            `json:"name"`               // 服务名称
	Address           string            `json:"address"`            // 服务地址
	Port              int               `json:"port"`               // 服务端口
	Tags              []string          `json:"tags"`               // 服务标签
	Meta              map[string]string `json:"meta"`               // 元数据
	HealthCheckURL    string            `json:"health_check_url"`   // 健康检查URL
	HealthCheckInterval string          `json:"health_check_interval"` // 健康检查间隔
	HealthCheckTimeout  string          `json:"health_check_timeout"`  // 健康检查超时
}

// ConsulClient Consul客户端
type ConsulClient struct {
	client *api.Client
	config *ConsulConfig
}

// ConsulConfig Consul配置
type ConsulConfig struct {
	Address    string `json:"address"`     // Consul地址
	Scheme     string `json:"scheme"`      // 协议(http/https)
	Token      string `json:"token"`       // ACL Token
	Datacenter string `json:"datacenter"`  // 数据中心
	Namespace  string `json:"namespace"`   // 命名空间
}

// NewConsulClient 创建Consul客户端
func NewConsulClient(cfg *ConsulConfig) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = cfg.Address
	config.Scheme = cfg.Scheme
	config.Token = cfg.Token
	config.Datacenter = cfg.Datacenter

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &ConsulClient{
		client: client,
		config: cfg,
	}, nil
}

// RegisterService 注册服务
func (c *ConsulClient) RegisterService(ctx context.Context, info *ServiceInfo) error {
	registration := &api.AgentServiceRegistration{
		ID:      info.ID,
		Name:    info.Name,
		Address: info.Address,
		Port:    info.Port,
		Tags:    info.Tags,
		Meta:    info.Meta,
		Check: &api.AgentServiceCheck{
			CheckID:        fmt.Sprintf("%s-check", info.ID),
			HTTP:           info.HealthCheckURL,
			Interval:       info.HealthCheckInterval,
			Timeout:        info.HealthCheckTimeout,
			DeregisterCriticalServiceAfter: "30s", // 30s后注销失败服务
		},
	}

	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	return nil
}

// DeregisterService 注销服务
func (c *ConsulClient) DeregisterService(ctx context.Context, serviceID string) error {
	if err := c.client.Agent().ServiceDeregister(serviceID); err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	return nil
}

// DiscoverService 发现服务实例
func (c *ConsulClient) DiscoverService(ctx context.Context, serviceName string, tags []string) ([]*ServiceInfo, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, &api.QueryOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to discover service: %w", err)
	}

	var result []*ServiceInfo
	for _, service := range services {
		// 标签过滤
		if tags != nil && len(tags) > 0 {
			if !containsAllTags(service.Service.Tags, tags) {
				continue
			}
		}

		result = append(result, &ServiceInfo{
			ID:      service.Service.ID,
			Name:    service.Service.Service,
			Address: service.Service.Address,
			Port:    service.Service.Port,
			Tags:    service.Service.Tags,
			Meta:    service.Service.Meta,
		})
	}

	return result, nil
}

// GetOneService 获取一个健康的实例（简单负载均衡）
func (c *ConsulClient) GetOneService(ctx context.Context, serviceName string, tags []string) (*ServiceInfo, error) {
	services, err := c.DiscoverService(ctx, serviceName, tags)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	// 简单的随机选择（生产环境应使用更复杂的负载均衡策略）
	idx := time.Now().UnixNano() % int64(len(services))
	return services[idx], nil
}

// WatchServices 监控服务变化
func (c *ConsulClient) WatchServices(ctx context.Context, serviceName string, callback func([]*ServiceInfo)) error {
	ticker := time.NewTicker(5 * time.Second) // 每5秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			services, err := c.DiscoverService(ctx, serviceName, nil)
			if err != nil {
				continue
			}
			callback(services)
		}
	}
}

// GetServiceConfiguration 获取服务配置（KV存储）
func (c *ConsulClient) GetServiceConfiguration(ctx context.Context, key string) (map[string]interface{}, error) {
	kvPair, _, err := c.client.KV().Get(key, &api.QueryOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	if kvPair == nil {
		return make(map[string]interface{}), nil
	}

	var config map[string]interface{}
	if err := json.Unmarshal(kvPair.Value, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return config, nil
}

// SetServiceConfiguration 设置服务配置
func (c *ConsulClient) SetServiceConfiguration(ctx context.Context, key string, config map[string]interface{}) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	kvPair := &api.KVPair{
		Key:   key,
		Value: data,
	}

	if _, err := c.client.KV().Put(kvPair, &api.WriteOptions{}); err != nil {
		return fmt.Errorf("failed to set configuration: %w", err)
	}

	return nil
}

// WatchConfiguration 监控配置变化
func (c *ConsulClient) WatchConfiguration(ctx context.Context, key string, callback func(map[string]interface{})) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	lastIndex := uint64(0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			kvPair, meta, err := c.client.KV().Get(key, &api.QueryOptions{
				WaitIndex: lastIndex,
			})
			if err != nil {
				continue
			}

			if meta.LastIndex <= lastIndex {
				continue
			}

			lastIndex = meta.LastIndex

			if kvPair != nil {
				var config map[string]interface{}
				if err := json.Unmarshal(kvPair.Value, &config); err == nil {
					callback(config)
				}
			}
		}
	}
}

// CreateSession 创建会话（用于分布式锁）
func (c *ConsulClient) CreateSession(ctx context.Context, name string, ttl string) (string, error) {
	session := &api.SessionEntry{
		Name:     name,
		TTL:      ttl,
		Behavior: api.SessionBehaviorDelete,
	}

	id, _, err := c.client.Session().Create(session, &api.WriteOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return id, nil
}

// DestroySession 销毁会话
func (c *ConsulClient) DestroySession(ctx context.Context, sessionID string) error {
	_, err := c.client.Session().Destroy(sessionID, &api.WriteOptions{})
	return err
}

// AcquireLock 获取分布式锁
func (c *ConsulClient) AcquireLock(ctx context.Context, key string, sessionID string) (bool, error) {
	kvPair := &api.KVPair{
		Key:     key,
		Value:   []byte(sessionID),
		Session: sessionID,
	}

	acquired, _, err := c.client.KV().Acquire(kvPair, &api.WriteOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return acquired, nil
}

// ReleaseLock 释放分布式锁
func (c *ConsulClient) ReleaseLock(ctx context.Context, key string, sessionID string) error {
	kvPair := &api.KVPair{
		Key:     key,
		Value:   []byte(sessionID),
		Session: sessionID,
	}

	_, err := c.client.KV().Release(kvPair, &api.WriteOptions{})
	return err
}

// 辅助函数：检查是否包含所有标签
func containsAllTags(serviceTags, requiredTags []string) bool {
	tagMap := make(map[string]string)
	for _, tag := range serviceTags {
		tagMap[tag] = tag
	}

	for _, required := range requiredTags {
		if _, exists := tagMap[required]; !exists {
			return false
		}
	}

	return true
}
