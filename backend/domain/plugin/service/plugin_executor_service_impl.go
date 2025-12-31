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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/repository"
)

const (
	// DefaultTimeout 默认超时时间(秒)
	DefaultTimeout = 30
	// DefaultMaxRetries 默认最大重试次数
	DefaultMaxRetries = 3
	// MaxConcurrentExecutions 最大并发执行数
	MaxConcurrentExecutions = 100
)

// NewPluginExecutorService 创建插件执行引擎服务
func NewPluginExecutorService(
	pluginRepo repository.PluginRepository,
	db *gorm.DB,
) PluginExecutorService {
	return &pluginExecutorServiceImpl{
		pluginRepo:      pluginRepo,
		db:              db,
		httpClient:      &http.Client{Timeout: 90 * time.Second},
		executionCache:  make(map[string]*entity.PluginExecution),
		executionMutexes: make(map[string]*sync.Mutex),
		cacheMutex:      &sync.RWMutex{},
	}
}

type pluginExecutorServiceImpl struct {
	pluginRepo      repository.PluginRepository
	db              *gorm.DB
	httpClient      *http.Client
	executionCache  map[string]*entity.PluginExecution
	executionMutexes map[string]*sync.Mutex
	cacheMutex      *sync.RWMutex
}

// Execute 执行插件(异步)
func (s *pluginExecutorServiceImpl) Execute(
	ctx context.Context,
	req *entity.ExecutePluginRequest,
) (*entity.ExecutePluginResponse, error) {
	// 1. 验证请求
	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 2. 获取插件配置
	plugin, exist, err := s.pluginRepo.GetOnlinePlugin(ctx, 0) // TODO: 根据name查询
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin: %w", err)
	}
	if !exist {
		return nil, fmt.Errorf("plugin not found: %s", req.PluginName)
	}

	// 3. 创建执行记录
	execution := s.createExecution(req, plugin)

	// 4. 存储到数据库
	if err := s.db.Create(execution).Error; err != nil {
		return nil, fmt.Errorf("failed to create execution record: %w", err)
	}

	// 5. 缓存执行记录
	s.cacheExecution(execution)

	// 6. 异步执行
	go s.executePlugin(context.Background(), plugin, execution)

	// 7. 返回响应
	return &entity.ExecutePluginResponse{
		ExecutionID: execution.ExecutionID,
		Status:      entity.ExecutionStatusRunning,
	}, nil
}

// ExecuteSync 同步执行插件
func (s *pluginExecutorServiceImpl) ExecuteSync(
	ctx context.Context,
	req *entity.ExecutePluginRequest,
) (*entity.ExecutePluginResponse, error) {
	// 1. 验证请求
	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 2. 获取插件配置
	plugin, exist, err := s.pluginRepo.GetOnlinePlugin(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin: %w", err)
	}
	if !exist {
		return nil, fmt.Errorf("plugin not found: %s", req.PluginName)
	}

	// 3. 创建执行记录
	execution := s.createExecution(req, plugin)

	// 4. 存储到数据库
	if err := s.db.Create(execution).Error; err != nil {
		return nil, fmt.Errorf("failed to create execution record: %w", err)
	}

	// 5. 同步执行
	result, err := s.executeWithRetry(ctx, plugin, execution)

	// 6. 更新数据库
	s.db.Save(execution)

	// 7. 返回结果
	if err != nil {
		return &entity.ExecutePluginResponse{
			ExecutionID: execution.ExecutionID,
			Status:      execution.Status,
			Error:       err.Error(),
		}, err
	}

	return &entity.ExecutePluginResponse{
		ExecutionID: execution.ExecutionID,
		Status:      execution.Status,
		Result:      execution.Result,
	}, nil
}

// GetExecution 获取执行记录
func (s *pluginExecutorServiceImpl) GetExecution(
	ctx context.Context,
	req *entity.GetExecutionRequest,
) (*entity.GetExecutionResponse, error) {
	// 1. 从缓存获取
	if execution := s.getCachedExecution(req.ExecutionID); execution != nil {
		return &entity.GetExecutionResponse{
			Execution: execution,
			Result:    execution.Result,
			Error:     execution.Error,
		}, nil
	}

	// 2. 从数据库获取
	var execution entity.PluginExecution
	err := s.db.Where("execution_id = ? AND tenant_id = ?", req.ExecutionID, req.TenantID).
		First(&execution).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}

	// 3. 缓存结果
	s.cacheExecution(&execution)

	return &entity.GetExecutionResponse{
		Execution: &execution,
		Result:    execution.Result,
		Error:     execution.Error,
	}, nil
}

// CancelExecution 取消执行
func (s *pluginExecutorServiceImpl) CancelExecution(
	ctx context.Context,
	executionID string,
) error {
	// 1. 获取执行记录
	execution := s.getCachedExecution(executionID)
	if execution == nil {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	// 2. 检查状态
	if execution.IsCompleted() {
		return fmt.Errorf("execution already completed: %s", executionID)
	}

	// 3. 更新状态
	execution.Status = entity.ExecutionStatusCancelled
	now := time.Now()
	execution.CompletedAt = &now
	execution.CalculateDuration()

	// 4. 保存到数据库
	err := s.db.Save(execution).Error
	if err != nil {
		return fmt.Errorf("failed to cancel execution: %w", err)
	}

	return nil
}

// BatchExecute 批量执行插件
func (s *pluginExecutorServiceImpl) BatchExecute(
	ctx context.Context,
	reqs []*entity.ExecutePluginRequest,
) ([]*entity.ExecutePluginResponse, error) {
	responses := make([]*entity.ExecutePluginResponse, len(reqs))
	errors := make([]error, len(reqs))

	// 使用WaitGroup并发执行
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, MaxConcurrentExecutions)

	for i, req := range reqs {
		wg.Add(1)
		go func(index int, r *entity.ExecutePluginRequest) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			resp, err := s.Execute(ctx, r)
			responses[index] = resp
			errors[index] = err
		}(i, req)
	}

	wg.Wait()

	// 检查错误
	for _, err := range errors {
		if err != nil {
			return responses, fmt.Errorf("batch execution partially failed: %w", err)
		}
	}

	return responses, nil
}

// RetryExecution 重试执行
func (s *pluginExecutorServiceImpl) RetryExecution(
	ctx context.Context,
	executionID string,
) (*entity.ExecutePluginResponse, error) {
	// 1. 获取原执行记录
	var execution entity.PluginExecution
	err := s.db.Where("execution_id = ?", executionID).First(&execution).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}

	// 2. 检查是否可以重试
	if !execution.CanRetry() {
		return nil, fmt.Errorf("execution cannot be retried: %s", executionID)
	}

	// 3. 构建新请求
	req := &entity.ExecutePluginRequest{
		PluginName:    execution.PluginName,
		PluginVersion: execution.PluginVersion,
		Parameters:    execution.Parameters,
		Timeout:       execution.Timeout,
		MaxRetries:    execution.MaxRetries,
		TenantID:      execution.TenantID,
		UserID:        execution.UserID,
	}

	// 4. 执行
	return s.ExecuteSync(ctx, req)
}

// validateRequest 验证请求
func (s *pluginExecutorServiceImpl) validateRequest(req *entity.ExecutePluginRequest) error {
	if req.PluginName == "" {
		return fmt.Errorf("plugin_name is required")
	}
	if req.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if req.Parameters == nil {
		req.Parameters = make(map[string]interface{})
	}
	if req.Timeout == 0 {
		req.Timeout = DefaultTimeout
	}
	if req.MaxRetries == 0 {
		req.MaxRetries = DefaultMaxRetries
	}
	return nil
}

// createExecution 创建执行记录
func (s *pluginExecutorServiceImpl) createExecution(
	req *entity.ExecutePluginRequest,
	plugin *entity.PluginInfo,
) *entity.PluginExecution {
	now := time.Now()
	return &entity.PluginExecution{
		ExecutionID:    uuid.New().String(),
		PluginID:       plugin.PluginInfo.ID,
		PluginName:     req.PluginName,
		PluginVersion:  req.PluginVersion,
		Parameters:     req.Parameters,
		Status:         entity.ExecutionStatusRunning,
		RetryCount:     0,
		MaxRetries:     req.MaxRetries,
		Timeout:        req.Timeout,
		StartedAt:      now,
		CreatedAt:      now,
		UpdatedAt:      now,
		TenantID:       req.TenantID,
		UserID:         req.UserID,
	}
}

// executePlugin 执行插件(异步)
func (s *pluginExecutorServiceImpl) executePlugin(
	ctx context.Context,
	plugin *entity.PluginInfo,
	execution *entity.PluginExecution,
) {
	// 获取执行互斥锁
	mutex := s.getExecutionMutex(execution.ExecutionID)
	mutex.Lock()
	defer mutex.Unlock()

	// 执行并重试
	result, err := s.executeWithRetry(ctx, plugin, execution)

	// 更新数据库
	s.db.Save(execution)

	// 清理缓存
	s.removeCachedExecution(execution.ExecutionID)
}

// executeWithRetry 执行并重试
func (s *pluginExecutorServiceImpl) executeWithRetry(
	ctx context.Context,
	plugin *entity.PluginInfo,
	execution *entity.PluginExecution,
) (map[string]interface{}, error) {
	var lastErr error

	for attempt := 0; attempt <= execution.MaxRetries; attempt++ {
		if attempt > 0 {
			// 指数退避
			backoffTime := time.Duration(attempt) * time.Second
			time.Sleep(backoffTime)
			execution.IncrementRetry()
		}

		// 检查是否已取消
		if execution.Status == entity.ExecutionStatusCancelled {
			return nil, fmt.Errorf("execution cancelled")
		}

		// 执行插件
		result, err := s.executeOnce(ctx, plugin, execution)
		if err == nil {
			execution.SetResult(result)
			return result, nil
		}

		lastErr = err
	}

	// 所有重试都失败
	execution.SetError(lastErr)
	return nil, lastErr
}

// executeOnce 执行一次插件
func (s *pluginExecutorServiceImpl) executeOnce(
	ctx context.Context,
	plugin *entity.PluginInfo,
	execution *entity.PluginExecution,
) (map[string]interface{}, error) {
	// 1. 构建请求URL
	endpoint := s.getPluginEndpoint(plugin)

	// 2. 设置超时
	timeout := time.Duration(execution.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 3. 构建请求体
	reqBody, err := json.Marshal(execution.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	// 4. 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 5. 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Execution-ID", execution.ExecutionID)
	httpReq.Header.Set("X-Tenant-ID", execution.TenantID)
	httpReq.Header.Set("X-User-ID", execution.UserID)

	// 6. 发送请求
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 7. 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 8. 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plugin returned error: %s", string(body))
	}

	// 9. 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

// getPluginEndpoint 获取插件端点
func (s *pluginExecutorServiceImpl) getPluginEndpoint(plugin *entity.PluginInfo) string {
	// 简化实现:使用ServerURL
	if plugin.GetServerURL() != "" {
		return plugin.GetServerURL()
	}
	return "http://localhost:8080/execute"
}

// cacheExecution 缓存执行记录
func (s *pluginExecutorServiceImpl) cacheExecution(execution *entity.PluginExecution) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.executionCache[execution.ExecutionID] = execution
}

// getCachedExecution 获取缓存的执行记录
func (s *pluginExecutorServiceImpl) getCachedExecution(executionID string) *entity.PluginExecution {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()
	return s.executionCache[executionID]
}

// removeCachedExecution 移除缓存的执行记录
func (s *pluginExecutorServiceImpl) removeCachedExecution(executionID string) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	delete(s.executionCache, executionID)
}

// getExecutionMutex 获取执行互斥锁
func (s *pluginExecutorServiceImpl) getExecutionMutex(executionID string) *sync.Mutex {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	if _, exists := s.executionMutexes[executionID]; !exists {
		s.executionMutexes[executionID] = &sync.Mutex{}
	}

	return s.executionMutexes[executionID]
}
