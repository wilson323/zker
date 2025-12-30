package saga

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRepository Mock仓储
type MockRepository struct {
	sagas      map[string]*Saga
	executions map[string]*SagaExecution
	saveError  error
	findError  error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		sagas:      make(map[string]*Saga),
		executions: make(map[string]*SagaExecution),
	}
}

func (m *MockRepository) SaveSaga(ctx context.Context, saga *Saga) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.sagas[saga.ID] = saga
	return nil
}

func (m *MockRepository) FindSagaByID(ctx context.Context, id string) (*Saga, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	saga, ok := m.sagas[id]
	if !ok {
		return nil, errors.New("saga not found")
	}
	return saga, nil
}

func (m *MockRepository) FindSagaByName(ctx context.Context, name string) (*Saga, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	for _, saga := range m.sagas {
		if saga.Name == name {
			return saga, nil
		}
	}
	return nil, errors.New("saga not found")
}

func (m *MockRepository) ListSagas(ctx context.Context) ([]*Saga, error) {
	sagas := make([]*Saga, 0, len(m.sagas))
	for _, saga := range m.sagas {
		sagas = append(sagas, saga)
	}
	return sagas, nil
}

func (m *MockRepository) SaveExecution(ctx context.Context, execution *SagaExecution) error {
	m.executions[execution.ID] = execution
	return nil
}

func (m *MockRepository) UpdateExecution(ctx context.Context, execution *SagaExecution) error {
	m.executions[execution.ID] = execution
	return nil
}

func (m *MockRepository) FindExecutionByID(ctx context.Context, id string) (*SagaExecution, error) {
	execution, ok := m.executions[id]
	if !ok {
		return nil, errors.New("execution not found")
	}
	return execution, nil
}

func (m *MockRepository) FindExecutionsBySagaID(ctx context.Context, sagaID string) ([]*SagaExecution, error) {
	executions := make([]*SagaExecution, 0)
	for _, exec := range m.executions {
		if exec.SagaID == sagaID {
			executions = append(executions, exec)
		}
	}
	return executions, nil
}

// MockStep Mock步骤
type MockStep struct {
	name        string
	timeout     time.Duration
	executeFunc func(ctx context.Context, data interface{}) (interface{}, error)
}

func (m *MockStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, data)
	}
	return data, nil
}

func (m *MockStep) Name() string {
	return m.name
}

func (m *MockStep) Timeout() time.Duration {
	return m.timeout
}

// MockCompensationStep Mock补偿步骤
type MockCompensationStep struct {
	name           string
	timeout        time.Duration
	compensateFunc func(ctx context.Context, data interface{}) error
}

func (m *MockCompensationStep) Compensate(ctx context.Context, data interface{}) error {
	if m.compensateFunc != nil {
		return m.compensateFunc(ctx, data)
	}
	return nil
}

func (m *MockCompensationStep) Name() string {
	return m.name
}

func (m *MockCompensationStep) Timeout() time.Duration {
	return m.timeout
}

// TestNewSagaOrchestrator 测试创建Saga协调器
func TestNewSagaOrchestrator(t *testing.T) {
	repo := NewMockRepository()
	logger := zap.NewNop()

	orchestrator := NewSagaOrchestrator(repo, logger)

	assert.NotNil(t, orchestrator)
	assert.Equal(t, repo, orchestrator.repo)
	assert.Equal(t, logger, orchestrator.logger)
}

// TestDefineSaga 测试定义Saga
func TestDefineSaga(t *testing.T) {
	repo := NewMockRepository()
	logger := zap.NewNop()
	orchestrator := NewSagaOrchestrator(repo, logger)

	saga := &Saga{
		ID:          "test-saga",
		Name:        "测试Saga",
		Description: "测试Saga定义",
		Steps: []SagaStep{
			&MockStep{name: "step1", timeout: 5 * time.Second},
		},
		Compensations: []CompensationStep{
			&MockCompensationStep{name: "comp1", timeout: 5 * time.Second},
		},
		Timeout:     30 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	err := orchestrator.DefineSaga(saga)

	assert.NoError(t, err)
	assert.Contains(t, repo.sagas, saga.ID)
}

// TestExecuteSaga_Success 测试Saga执行成功
func TestExecuteSaga_Success(t *testing.T) {
	repo := NewMockRepository()
	logger := zap.NewNop()
	orchestrator := NewSagaOrchestrator(repo, logger)

	// 定义Saga
	saga := &Saga{
		ID:          "test-saga",
		Name:        "测试Saga",
		Description: "测试Saga执行",
		Steps: []SagaStep{
			&MockStep{
				name:    "step1",
				timeout: 5 * time.Second,
				executeFunc: func(ctx context.Context, data interface{}) (interface{}, error) {
					return "step1-output", nil
				},
			},
			&MockStep{
				name:    "step2",
				timeout: 5 * time.Second,
				executeFunc: func(ctx context.Context, data interface{}) (interface{}, error) {
					return "step2-output", nil
				},
			},
		},
		Compensations: []CompensationStep{
			&MockCompensationStep{name: "comp1", timeout: 5 * time.Second},
			&MockCompensationStep{name: "comp2", timeout: 5 * time.Second},
		},
		Timeout:     30 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	err := orchestrator.DefineSaga(saga)
	require.NoError(t, err)

	// 执行Saga
	ctx := context.Background()
	execution, err := orchestrator.ExecuteSaga(ctx, "测试Saga", "input-data")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, execution)
	assert.Equal(t, SagaStatusCompleted, execution.Status)
	assert.Equal(t, 2, len(execution.StepExecutions))
	assert.Equal(t, 2, execution.CurrentStep)

	// 验证每个步骤都成功
	for _, stepExec := range execution.StepExecutions {
		assert.Equal(t, StepStatusCompleted, stepExec.Status)
		assert.Nil(t, stepExec.Error)
		assert.NotNil(t, stepExec.CompletedAt)
	}
}

// TestExecuteSaga_StepFailure 测试Saga步骤失败并补偿
func TestExecuteSaga_StepFailure(t *testing.T) {
	repo := NewMockRepository()
	logger := zap.NewNop()
	orchestrator := NewSagaOrchestrator(repo, logger)

	// 定义Saga
	step2Executed := false
	saga := &Saga{
		ID:          "test-saga-fail",
		Name:        "测试Saga失败",
		Description: "测试Saga失败和补偿",
		Steps: []SagaStep{
			&MockStep{
				name:    "step1",
				timeout: 5 * time.Second,
				executeFunc: func(ctx context.Context, data interface{}) (interface{}, error) {
					return "step1-output", nil
				},
			},
			&MockStep{
				name:    "step2",
				timeout: 5 * time.Second,
				executeFunc: func(ctx context.Context, data interface{}) (interface{}, error) {
					step2Executed = true
					return nil, errors.New("step2 failed")
				},
			},
			&MockStep{
				name:    "step3",
				timeout: 5 * time.Second,
				executeFunc: func(ctx context.Context, data interface{}) (interface{}, error) {
					return "step3-output", nil
				},
			},
		},
		Compensations: []CompensationStep{
			&MockCompensationStep{
				name:    "comp1",
				timeout: 5 * time.Second,
				compensateFunc: func(ctx context.Context, data interface{}) error {
					return nil
				},
			},
			&MockCompensationStep{
				name:    "comp2",
				timeout: 5 * time.Second,
				compensateFunc: func(ctx context.Context, data interface{}) error {
					return nil
				},
			},
			&MockCompensationStep{
				name:    "comp3",
				timeout: 5 * time.Second,
				compensateFunc: func(ctx context.Context, data interface{}) error {
					return nil
				},
			},
		},
		Timeout:     30 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	err := orchestrator.DefineSaga(saga)
	require.NoError(t, err)

	// 执行Saga
	ctx := context.Background()
	execution, err := orchestrator.ExecuteSaga(ctx, "测试Saga失败", "input-data")

	// 验证结果
	assert.Error(t, err)
	assert.NotNil(t, execution)
	assert.Equal(t, SagaStatusCompensated, execution.Status)
	assert.True(t, step2Executed)

	// 验证步骤执行状态
	assert.Equal(t, 2, len(execution.StepExecutions))
	assert.Equal(t, StepStatusCompleted, execution.StepExecutions[0].Status)
	assert.Equal(t, StepStatusFailed, execution.StepExecutions[1].Status)
	assert.Equal(t, 1, execution.CurrentStep)
}

// TestExecuteSaga_Retry 测试步骤重试
func TestExecuteSaga_Retry(t *testing.T) {
	repo := NewMockRepository()
	logger := zap.NewNop()
	orchestrator := NewSagaOrchestrator(repo, logger)

	// 定义重试策略
	retryPolicy := &RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	attempts := 0
	saga := &Saga{
		ID:          "test-saga-retry",
		Name:        "测试Saga重试",
		Description: "测试步骤重试机制",
		Steps: []SagaStep{
			&MockStep{
				name:    "step1",
				timeout: 5 * time.Second,
				executeFunc: func(ctx context.Context, data interface{}) (interface{}, error) {
					attempts++
					if attempts < 2 {
						return nil, errors.New("step1 temporary failed")
					}
					return "step1-output", nil
				},
			},
		},
		Compensations: []CompensationStep{
			&MockCompensationStep{name: "comp1", timeout: 5 * time.Second},
		},
		Timeout:     30 * time.Second,
		RetryPolicy: retryPolicy,
	}

	err := orchestrator.DefineSaga(saga)
	require.NoError(t, err)

	// 执行Saga
	ctx := context.Background()
	execution, err := orchestrator.ExecuteSaga(ctx, "测试Saga重试", "input-data")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, execution)
	assert.Equal(t, SagaStatusCompleted, execution.Status)
	assert.Equal(t, 2, attempts) // 第1次失败,第2次成功
}

// TestGetStatus 测试获取Saga状态
func TestGetStatus(t *testing.T) {
	repo := NewMockRepository()
	logger := zap.NewNop()
	orchestrator := NewSagaOrchestrator(repo, logger)

	// 定义并执行Saga
	saga := &Saga{
		ID:          "test-saga-status",
		Name:        "测试Saga状态",
		Description: "测试获取Saga状态",
		Steps: []SagaStep{
			&MockStep{name: "step1", timeout: 5 * time.Second},
		},
		Compensations: []CompensationStep{
			&MockCompensationStep{name: "comp1", timeout: 5 * time.Second},
		},
		Timeout:     30 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	err := orchestrator.DefineSaga(saga)
	require.NoError(t, err)

	ctx := context.Background()
	execution, err := orchestrator.ExecuteSaga(ctx, "测试Saga状态", "input-data")
	require.NoError(t, err)

	// 获取状态
	status, err := orchestrator.GetStatus(ctx, execution.ID)

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, execution.ID, status.ID)
	assert.Equal(t, SagaStatusCompleted, status.Status)
}

// TestValidateSaga 测试Saga验证
func TestValidateSaga(t *testing.T) {
	repo := NewMockRepository()
	logger := zap.NewNop()
	orchestrator := NewSagaOrchestrator(repo, logger)

	tests := []struct {
		name        string
		saga        *Saga
		expectError bool
		errorMsg    string
	}{
		{
			name: "有效Saga",
			saga: &Saga{
				ID:          "valid-saga",
				Name:        "有效Saga",
				Description: "测试有效的Saga",
				Steps: []SagaStep{
					&MockStep{name: "step1", timeout: 5 * time.Second},
				},
				Compensations: []CompensationStep{
					&MockCompensationStep{name: "comp1", timeout: 5 * time.Second},
				},
				Timeout:     30 * time.Second,
				RetryPolicy: DefaultRetryPolicy(),
			},
			expectError: false,
		},
		{
			name: "空ID",
			saga: &Saga{
				ID:   "",
				Name: "空ID Saga",
				Steps: []SagaStep{
					&MockStep{name: "step1", timeout: 5 * time.Second},
				},
				Compensations: []CompensationStep{
					&MockCompensationStep{name: "comp1", timeout: 5 * time.Second},
				},
				Timeout:     30 * time.Second,
				RetryPolicy: DefaultRetryPolicy(),
			},
			expectError: true,
			errorMsg:    "Saga ID不能为空",
		},
		{
			name: "步骤和补偿数量不匹配",
			saga: &Saga{
				ID:          "mismatch-saga",
				Name:        "数量不匹配",
				Description: "步骤和补偿数量不匹配",
				Steps: []SagaStep{
					&MockStep{name: "step1", timeout: 5 * time.Second},
				},
				Compensations: []CompensationStep{
					&MockCompensationStep{name: "comp1", timeout: 5 * time.Second},
					&MockCompensationStep{name: "comp2", timeout: 5 * time.Second},
				},
				Timeout:     30 * time.Second,
				RetryPolicy: DefaultRetryPolicy(),
			},
			expectError: true,
			errorMsg:    "步骤数量与补偿数量不匹配",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := orchestrator.DefineSaga(tt.saga)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
