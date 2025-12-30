package saga

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Orchestrator Saga协调器接口
type Orchestrator interface {
	// DefineSaga 定义Saga
	DefineSaga(saga *Saga) error

	// ExecuteSaga 执行Saga
	ExecuteSaga(ctx context.Context, sagaName string, input interface{}) (*SagaExecution, error)

	// GetStatus 获取Saga状态
	GetStatus(ctx context.Context, executionID string) (*SagaExecution, error)
}

// SagaOrchestrator Saga协调器实现
type SagaOrchestrator struct {
	repo   Repository
	logger *zap.Logger
}

// NewSagaOrchestrator 创建Saga协调器
func NewSagaOrchestrator(repo Repository, logger *zap.Logger) *SagaOrchestrator {
	return &SagaOrchestrator{
		repo:   repo,
		logger: logger,
	}
}

// DefineSaga 定义Saga
func (o *SagaOrchestrator) DefineSaga(saga *Saga) error {
	if err := o.validateSaga(saga); err != nil {
		return fmt.Errorf("无效的Saga定义: %w", err)
	}

	ctx := context.Background()
	if err := o.repo.SaveSaga(ctx, saga); err != nil {
		return fmt.Errorf("保存Saga定义失败: %w", err)
	}

	o.logger.Info("Saga定义已注册",
		zap.String("saga_id", saga.ID),
		zap.String("saga_name", saga.Name))

	return nil
}

// ExecuteSaga 执行Saga
func (o *SagaOrchestrator) ExecuteSaga(
	ctx context.Context,
	sagaName string,
	input interface{},
) (*SagaExecution, error) {
	// 1. 加载Saga定义
	saga, err := o.repo.FindSagaByName(ctx, sagaName)
	if err != nil {
		return nil, &SagaError{
			Code:    ErrSagaNotFound.Code,
			Message: fmt.Sprintf("Saga不存在: %s", sagaName),
			Cause:   err,
			Recoverable: ErrSagaNotFound.Recoverable,
		}
	}

	// 2. 创建Saga执行记录
	execution := &SagaExecution{
		ID:             generateID(),
		SagaID:         saga.ID,
		Status:         SagaStatusRunning,
		CurrentStep:    0,
		InputData:      input,
		StartedAt:      time.Now(),
		StepExecutions: make([]StepExecution, 0, len(saga.Steps)),
	}

	// 3. 保存执行记录
	if err := o.repo.SaveExecution(ctx, execution); err != nil {
		return nil, fmt.Errorf("保存Saga执行记录失败: %w", err)
	}

	o.logger.Info("开始执行Saga",
		zap.String("execution_id", execution.ID),
		zap.String("saga_name", sagaName))

	// 4. 设置Saga超时
	sagaCtx, cancel := context.WithTimeout(ctx, saga.Timeout)
	defer cancel()

	// 5. 依次执行每个步骤
	var lastOutput interface{}
	var compensationData []interface{}

	for i, step := range saga.Steps {
		o.logger.Info("执行Saga步骤",
			zap.String("step_name", step.Name()),
			zap.Int("step_index", i))

		// 准备输入数据
		stepInput := input
		if i > 0 {
			stepInput = lastOutput
		}

		// 执行步骤(带重试)
		output, err := o.executeStepWithRetry(
			sagaCtx,
			step,
			stepInput,
			saga.RetryPolicy,
		)

		// 记录步骤执行
		stepExec := StepExecution{
			StepName:  step.Name(),
			Status:    StepStatusCompleted,
			Input:     stepInput,
			Output:    output,
			StartedAt: time.Now(),
		}

		if err != nil {
			// 步骤执行失败,开始补偿
			o.logger.Error("Saga步骤执行失败",
				zap.String("step_name", step.Name()),
				zap.Error(err))

			completedAt := time.Now()
			stepExec.Status = StepStatusFailed
			stepExec.Error = err
			stepExec.CompletedAt = &completedAt
			execution.StepExecutions = append(execution.StepExecutions, stepExec)

			// 更新执行状态
			execution.Status = SagaStatusCompensating
			execution.Error = err

			if updateErr := o.repo.UpdateExecution(ctx, execution); updateErr != nil {
				o.logger.Error("更新Saga执行记录失败", zap.Error(updateErr))
			}

			// 执行补偿
			if compErr := o.compensate(ctx, execution, saga, compensationData); compErr != nil {
				o.logger.Error("Saga补偿失败", zap.Error(compErr))
				execution.Status = SagaStatusFailed
			} else {
				execution.Status = SagaStatusCompensated
			}

			// 更新最终状态
			if updateErr := o.repo.UpdateExecution(ctx, execution); updateErr != nil {
				o.logger.Error("更新Saga执行记录失败", zap.Error(updateErr))
			}

			return execution, &SagaError{
				Code:        ErrStepFailed.Code,
				Message:     fmt.Sprintf("Saga执行失败: %s", step.Name()),
				SagaID:      execution.SagaID,
				StepName:    step.Name(),
				Cause:       err,
				Recoverable: ErrStepFailed.Recoverable,
			}
		}

		// 步骤执行成功
		completedAt := time.Now()
		stepExec.CompletedAt = &completedAt
		execution.StepExecutions = append(execution.StepExecutions, stepExec)

		// 保存补偿数据
		compensationData = append(compensationData, output)

		// 更新最后输出
		lastOutput = output
		execution.CurrentStep = i + 1

		// 更新执行记录
		if updateErr := o.repo.UpdateExecution(ctx, execution); updateErr != nil {
			o.logger.Error("更新Saga执行记录失败", zap.Error(updateErr))
		}
	}

	// 所有步骤执行成功
	execution.Status = SagaStatusCompleted
	execution.OutputData = lastOutput
	completedAt := time.Now()
	execution.CompletedAt = &completedAt

	if err := o.repo.UpdateExecution(ctx, execution); err != nil {
		o.logger.Error("更新Saga执行记录失败", zap.Error(err))
	}

	o.logger.Info("Saga执行成功",
		zap.String("execution_id", execution.ID),
		zap.String("saga_name", sagaName))

	return execution, nil
}

// GetStatus 获取Saga状态
func (o *SagaOrchestrator) GetStatus(ctx context.Context, executionID string) (*SagaExecution, error) {
	execution, err := o.repo.FindExecutionByID(ctx, executionID)
	if err != nil {
		return nil, fmt.Errorf("查询Saga执行记录失败: %w", err)
	}

	return execution, nil
}

// executeStepWithRetry 带重试的步骤执行
func (o *SagaOrchestrator) executeStepWithRetry(
	ctx context.Context,
	step SagaStep,
	input interface{},
	retryPolicy *RetryPolicy,
) (interface{}, error) {
	var lastErr error

	maxAttempts := retryPolicy.MaxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			o.logger.Info("重试Saga步骤",
				zap.String("step_name", step.Name()),
				zap.Int("attempt", attempt))

			// 计算退避时间
			backoff := o.calculateBackoff(attempt, retryPolicy)

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		// 设置步骤超时
		stepCtx, cancel := context.WithTimeout(ctx, step.Timeout())

		output, err := step.Execute(stepCtx, input)
		cancel()

		if err == nil {
			return output, nil
		}

		lastErr = err
		o.logger.Error("步骤执行失败",
			zap.String("step_name", step.Name()),
			zap.Int("attempt", attempt+1),
			zap.Error(err))
	}

	return nil, fmt.Errorf("步骤执行失败,已重试%d次: %w", maxAttempts, lastErr)
}

// compensate 执行补偿
func (o *SagaOrchestrator) compensate(
	ctx context.Context,
	execution *SagaExecution,
	saga *Saga,
	compensationData []interface{},
) error {
	o.logger.Info("开始Saga补偿",
		zap.String("execution_id", execution.ID),
		zap.Int("current_step", execution.CurrentStep))

	// 反向执行补偿步骤
	for i := execution.CurrentStep - 1; i >= 0; i-- {
		if i >= len(saga.Compensations) {
			o.logger.Warn("补偿步骤不存在",
				zap.Int("step_index", i),
				zap.Int("compensations_count", len(saga.Compensations)))
			continue
		}

		compensation := saga.Compensations[i]
		data := compensationData[i]

		o.logger.Info("执行补偿步骤",
			zap.String("compensation_name", compensation.Name()),
			zap.Int("step_index", i))

		// 设置补偿超时
		compCtx, cancel := context.WithTimeout(ctx, compensation.Timeout())

		if err := compensation.Compensate(compCtx, data); err != nil {
			cancel()
			o.logger.Error("补偿步骤执行失败",
				zap.String("compensation_name", compensation.Name()),
				zap.Error(err))
			// 继续执行后续补偿
			continue
		}
		cancel()

		o.logger.Info("补偿步骤执行成功",
			zap.String("compensation_name", compensation.Name()))
	}

	return nil
}

// calculateBackoff 计算退避时间
func (o *SagaOrchestrator) calculateBackoff(attempt int, policy *RetryPolicy) time.Duration {
	if policy == nil {
		return 1 * time.Second
	}

	// 指数退避
	backoff := time.Duration(float64(policy.InitialInterval) *
		float64(uint(1)<<uint(attempt-1)) * policy.Multiplier)

	// 限制最大退避时间
	if backoff > policy.MaxInterval {
		return policy.MaxInterval
	}

	return backoff
}

// validateSaga 验证Saga定义
func (o *SagaOrchestrator) validateSaga(saga *Saga) error {
	if saga.ID == "" {
		return fmt.Errorf("Saga ID不能为空")
	}

	if saga.Name == "" {
		return fmt.Errorf("Saga名称不能为空")
	}

	if len(saga.Steps) == 0 {
		return fmt.Errorf("Saga步骤不能为空")
	}

	if len(saga.Steps) != len(saga.Compensations) {
		return fmt.Errorf("步骤数量与补偿数量不匹配")
	}

	for i, step := range saga.Steps {
		if step == nil {
			return fmt.Errorf("第%d个步骤为空", i)
		}
		if step.Name() == "" {
			return fmt.Errorf("第%d个步骤名称为空", i)
		}
		if step.Timeout() <= 0 {
			return fmt.Errorf("第%d个步骤超时时间无效", i)
		}
	}

	for i, comp := range saga.Compensations {
		if comp == nil {
			return fmt.Errorf("第%d个补偿步骤为空", i)
		}
		if comp.Name() == "" {
			return fmt.Errorf("第%d个补偿步骤名称为空", i)
		}
		if comp.Timeout() <= 0 {
			return fmt.Errorf("第%d个补偿步骤超时时间无效", i)
		}
	}

	return nil
}
