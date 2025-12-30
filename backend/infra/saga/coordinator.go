// Saga分布式事务协调器
// 用于管理跨多个微服务的长事务，确保最终一致性
package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SagaTransaction Saga事务
type SagaTransaction struct {
	TransactionID   string                 `json:"transaction_id"`   // 事务ID
	TransactionType string                 `json:"transaction_type"` // 事务类型
	Status          StepStatus             `json:"status"`           // 事务状态
	Steps           []*TransactionStep     `json:"steps"`            // 事务步骤
	Context         map[string]interface{} `json:"context"`          // 上下文数据
	CreatedAt       time.Time              `json:"created_at"`       // 创建时间
	UpdatedAt       time.Time              `json:"updated_at"`       // 更新时间
	mu              sync.RWMutex           `json:"-"`
}

// TransactionStep 事务步骤
type TransactionStep struct {
	StepID         string                 `json:"step_id"`          // 步骤ID
	Service        string                 `json:"service"`          // 目标服务
	Action         string                 `json:"action"`           // 操作名称
	Status         StepStatus             `json:"status"`           // 步骤状态
	RequestPayload map[string]interface{} `json:"request_payload"`  // 请求参数
	ResponsePayload map[string]interface{} `json:"response_payload"` // 响应数据
	ErrorMessage   string                 `json:"error_message"`    // 错误信息
	ExecuteAt      time.Time              `json:"execute_at"`       // 执行时间
	CompletedAt    *time.Time             `json:"completed_at"`     // 完成时间
}

// StepExecutor 步骤执行器接口
type StepExecutor interface {
	Execute(ctx context.Context, step *TransactionStep) error
	Compensate(ctx context.Context, step *TransactionStep) error
}

// Coordinator Saga协调器
type Coordinator struct {
	executors    map[string]StepExecutor // 服务 -> 执行器
	transactions map[string]*SagaTransaction // 事务ID -> 事务
	storage      TransactionStorage         // 持久化存储
	eventBus     EventBus                  // 事件总线
	mu           sync.RWMutex
}

// TransactionStorage 事务存储接口
type TransactionStorage interface {
	Save(ctx context.Context, tx *SagaTransaction) error
	Get(ctx context.Context, transactionID string) (*SagaTransaction, error)
	List(ctx context.Context, status StepStatus) ([]*SagaTransaction, error)
	Delete(ctx context.Context, transactionID string) error
}

// EventBus 事件总线接口
type EventBus interface {
	Publish(ctx context.Context, event string, data interface{}) error
	Subscribe(event string, handler func(data interface{})) error
}

// NewCoordinator 创建Saga协调器
func NewCoordinator(storage TransactionStorage, eventBus EventBus) *Coordinator {
	return &Coordinator{
		executors:    make(map[string]StepExecutor),
		transactions: make(map[string]*SagaTransaction),
		storage:      storage,
		eventBus:     eventBus,
	}
}

// RegisterExecutor 注册服务执行器
func (c *Coordinator) RegisterExecutor(service string, executor StepExecutor) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.executors[service] = executor
}

// CreateTransaction 创建Saga事务
func (c *Coordinator) CreateTransaction(ctx context.Context, transactionType string, steps []*TransactionStep) (*SagaTransaction, error) {
	tx := &SagaTransaction{
		TransactionID:   uuid.New().String(),
		TransactionType: transactionType,
		Status:          StepStatusPending,
		Steps:           steps,
		Context:         make(map[string]interface{}),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 初始化步骤状态
	for _, step := range steps {
		step.StepID = uuid.New().String()
		step.Status = StepStatusPending
	}

	// 持久化
	if err := c.storage.Save(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// 内存缓存
	c.mu.Lock()
	c.transactions[tx.TransactionID] = tx
	c.mu.Unlock()

	// 发布事件
	c.eventBus.Publish(ctx, "saga.transaction.created", tx)

	return tx, nil
}

// ExecuteTransaction 执行Saga事务
func (c *Coordinator) ExecuteTransaction(ctx context.Context, transactionID string) error {
	c.mu.RLock()
	tx, exists := c.transactions[transactionID]
	c.mu.RUnlock()

	if !exists {
		return fmt.Errorf("transaction not found: %s", transactionID)
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.Status != StepStatusPending {
		return fmt.Errorf("transaction is not in pending status: %s", tx.Status)
	}

	tx.Status = StepStatusRunning
	tx.UpdatedAt = time.Now()
	c.storage.Save(ctx, tx)

	// 顺序执行所有步骤
	for i, step := range tx.Steps {
		if err := c.executeStep(ctx, tx, step); err != nil {
			// 步骤失败，执行补偿
			tx.Status = StepStatusFailed
			c.storage.Save(ctx, tx)
			return c.compensateTransaction(ctx, tx, i)
		}
	}

	// 所有步骤成功
	tx.Status = StepStatusCompleted
	tx.UpdatedAt = time.Now()
	c.storage.Save(ctx, tx)

	c.eventBus.Publish(ctx, "saga.transaction.completed", tx)
	return nil
}

// executeStep 执行单个步骤
func (c *Coordinator) executeStep(ctx context.Context, tx *SagaTransaction, step *TransactionStep) error {
	c.mu.RLock()
	executor, exists := c.executors[step.Service]
	c.mu.RUnlock()

	if !exists {
		return fmt.Errorf("executor not found for service: %s", step.Service)
	}

	step.Status = StepStatusRunning
	step.ExecuteAt = time.Now()
	tx.UpdatedAt = time.Now()
	c.storage.Save(ctx, tx)

	// 执行步骤
	if err := executor.Execute(ctx, step); err != nil {
		step.Status = StepStatusFailed
		step.ErrorMessage = err.Error()
		tx.UpdatedAt = time.Now()
		c.storage.Save(ctx, tx)
		return fmt.Errorf("step execution failed: %w", err)
	}

	step.Status = StepStatusCompleted
	now := time.Now()
	step.CompletedAt = &now
	tx.UpdatedAt = time.Now()
	c.storage.Save(ctx, tx)

	// 发布步骤完成事件
	c.eventBus.Publish(ctx, fmt.Sprintf("saga.step.%s.completed", step.Action), step)

	return nil
}

// compensateTransaction 补偿事务
func (c *Coordinator) compensateTransaction(ctx context.Context, tx *SagaTransaction, failedStepIndex int) error {
	tx.Status = StepStatusCompensating
	c.storage.Save(ctx, tx)

	// 逆序补偿已执行的步骤
	for i := failedStepIndex - 1; i >= 0; i-- {
		step := tx.Steps[i]
		if err := c.compensateStep(ctx, step); err != nil {
			// 补偿失败，记录错误但继续
			step.ErrorMessage = fmt.Sprintf("compensation failed: %v", err)
			c.storage.Save(ctx, tx)
		}
	}

	tx.Status = StepStatusCompensated
	tx.UpdatedAt = time.Now()
	c.storage.Save(ctx, tx)

	c.eventBus.Publish(ctx, "saga.transaction.compensated", tx)
	return fmt.Errorf("transaction compensated after failure at step %d", failedStepIndex)
}

// compensateStep 补偿单个步骤
func (c *Coordinator) compensateStep(ctx context.Context, step *TransactionStep) error {
	c.mu.RLock()
	executor, exists := c.executors[step.Service]
	c.mu.RUnlock()

	if !exists {
		return fmt.Errorf("executor not found for service: %s", step.Service)
	}

	step.Status = StepStatusCompensating

	if err := executor.Compensate(ctx, step); err != nil {
		return err
	}

	step.Status = StepStatusCompensated
	return nil
}

// GetTransaction 获取事务状态
func (c *Coordinator) GetTransaction(ctx context.Context, transactionID string) (*SagaTransaction, error) {
	c.mu.RLock()
	tx, exists := c.transactions[transactionID]
	c.mu.RUnlock()

	if exists {
		return tx, nil
	}

	// 从存储加载
	return c.storage.Get(ctx, transactionID)
}

// RetryFailedTransactions 重试失败的事务
func (c *Coordinator) RetryFailedTransactions(ctx context.Context) error {
	txs, err := c.storage.List(ctx, StepStatusFailed)
	if err != nil {
		return fmt.Errorf("failed to list failed transactions: %w", err)
	}

	for _, tx := range txs {
		// 重置失败步骤的状态
		for _, step := range tx.Steps {
			if step.Status == StepStatusFailed {
				step.Status = StepStatusPending
				step.ErrorMessage = ""
				step.ExecuteAt = time.Time{}
				step.CompletedAt = nil
			}
		}

		tx.Status = StepStatusPending
		tx.UpdatedAt = time.Now()

		// 重新执行
		if err := c.ExecuteTransaction(ctx, tx.TransactionID); err != nil {
			// 记录日志，继续下一个
			continue
		}
	}

	return nil
}

// JSON序列化辅助方法
func (tx *SagaTransaction) ToJSON() ([]byte, error) {
	tx.mu.RLock()
	defer tx.mu.RUnlock()
	return json.Marshal(tx)
}

func (tx *SagaTransaction) FromJSON(data []byte) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	return json.Unmarshal(data, tx)
}
