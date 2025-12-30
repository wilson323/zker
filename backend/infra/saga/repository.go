package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Repository Saga仓储接口
type Repository interface {
	// Saga定义管理
	SaveSaga(ctx context.Context, saga *Saga) error
	FindSagaByID(ctx context.Context, id string) (*Saga, error)
	FindSagaByName(ctx context.Context, name string) (*Saga, error)
	ListSagas(ctx context.Context) ([]*Saga, error)

	// Saga执行管理
	SaveExecution(ctx context.Context, execution *SagaExecution) error
	UpdateExecution(ctx context.Context, execution *SagaExecution) error
	FindExecutionByID(ctx context.Context, id string) (*SagaExecution, error)
	FindExecutionsBySagaID(ctx context.Context, sagaID string) ([]*SagaExecution, error)
}

// MySQLSagaRepository MySQL仓储实现
type MySQLSagaRepository struct {
	db *gorm.DB
}

// NewMySQLSagaRepository 创建MySQL仓储
func NewMySQLSagaRepository(db *gorm.DB) *MySQLSagaRepository {
	return &MySQLSagaRepository{db: db}
}

// SaveSaga 保存Saga定义
func (r *MySQLSagaRepository) SaveSaga(ctx context.Context, saga *Saga) error {
	// 序列化步骤
	stepsJSON, err := json.Marshal(saga.Steps)
	if err != nil {
		return fmt.Errorf("序列化步骤失败: %w", err)
	}

	// 序列化补偿步骤
	compensationsJSON, err := json.Marshal(saga.Compensations)
	if err != nil {
		return fmt.Errorf("序列化补偿步骤失败: %w", err)
	}

	// 序列化重试策略
	retryPolicyJSON, _ := json.Marshal(saga.RetryPolicy)

	record := &SagaRecord{
		ID:             saga.ID,
		Name:           saga.Name,
		Description:    saga.Description,
		Steps:          string(stepsJSON),
		Compensations:  string(compensationsJSON),
		RetryPolicy:    string(retryPolicyJSON),
		TimeoutSeconds: int64(saga.Timeout.Seconds()),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return r.db.WithContext(ctx).Create(record).Error
}

// FindSagaByID 根据ID查询Saga
func (r *MySQLSagaRepository) FindSagaByID(ctx context.Context, id string) (*Saga, error) {
	var record SagaRecord
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&record).
		Error

	if err != nil {
		return nil, err
	}

	return record.ToSaga()
}

// FindSagaByName 根据名称查询Saga
func (r *MySQLSagaRepository) FindSagaByName(ctx context.Context, name string) (*Saga, error) {
	var record SagaRecord
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&record).
		Error

	if err != nil {
		return nil, err
	}

	return record.ToSaga()
}

// ListSagas 查询所有Saga
func (r *MySQLSagaRepository) ListSagas(ctx context.Context) ([]*Saga, error) {
	var records []SagaRecord
	err := r.db.WithContext(ctx).Find(&records).Error
	if err != nil {
		return nil, err
	}

	sagas := make([]*Saga, 0, len(records))
	for _, record := range records {
		saga, err := record.ToSaga()
		if err != nil {
			return nil, err
		}
		sagas = append(sagas, saga)
	}

	return sagas, nil
}

// SaveExecution 保存Saga执行记录
func (r *MySQLSagaRepository) SaveExecution(
	ctx context.Context,
	execution *SagaExecution,
) error {
	// 序列化输入数据
	inputJSON, err := json.Marshal(execution.InputData)
	if err != nil {
		return fmt.Errorf("序列化输入数据失败: %w", err)
	}

	// 序列化步骤执行记录
	stepsJSON, err := json.Marshal(execution.StepExecutions)
	if err != nil {
		return fmt.Errorf("序列化步骤执行记录失败: %w", err)
	}

	var outputJSON []byte
	if execution.OutputData != nil {
		outputJSON, _ = json.Marshal(execution.OutputData)
	}

	record := &SagaExecutionRecord{
		ID:              execution.ID,
		SagaID:          execution.SagaID,
		Status:          string(execution.Status),
		CurrentStep:     execution.CurrentStep,
		InputData:       string(inputJSON),
		OutputData:      string(outputJSON),
		StepExecutions:  string(stepsJSON),
		StartedAt:       execution.StartedAt,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if execution.CompletedAt != nil {
		completedAt := *execution.CompletedAt
		record.CompletedAt = &completedAt
	}

	if execution.Error != nil {
		record.Error = execution.Error.Error()
	}

	return r.db.WithContext(ctx).Create(record).Error
}

// UpdateExecution 更新Saga执行记录
func (r *MySQLSagaRepository) UpdateExecution(
	ctx context.Context,
	execution *SagaExecution,
) error {
	// 序列化输入数据
	inputJSON, err := json.Marshal(execution.InputData)
	if err != nil {
		return fmt.Errorf("序列化输入数据失败: %w", err)
	}

	// 序列化步骤执行记录
	stepsJSON, err := json.Marshal(execution.StepExecutions)
	if err != nil {
		return fmt.Errorf("序列化步骤执行记录失败: %w", err)
	}

	var outputJSON []byte
	if execution.OutputData != nil {
		outputJSON, _ = json.Marshal(execution.OutputData)
	}

	updates := map[string]interface{}{
		"status":          string(execution.Status),
		"current_step":    execution.CurrentStep,
		"input_data":      string(inputJSON),
		"output_data":     string(outputJSON),
		"step_executions": string(stepsJSON),
		"updated_at":      time.Now(),
	}

	if execution.CompletedAt != nil {
		updates["completed_at"] = *execution.CompletedAt
	}

	if execution.Error != nil {
		updates["error"] = execution.Error.Error()
	} else {
		updates["error"] = nil
	}

	return r.db.WithContext(ctx).
		Model(&SagaExecutionRecord{}).
		Where("id = ?", execution.ID).
		Updates(updates).
		Error
}

// FindExecutionByID 根据ID查询执行记录
func (r *MySQLSagaRepository) FindExecutionByID(
	ctx context.Context,
	id string,
) (*SagaExecution, error) {
	var record SagaExecutionRecord
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&record).
		Error

	if err != nil {
		return nil, err
	}

	return record.ToExecution()
}

// FindExecutionsBySagaID 根据SagaID查询执行记录
func (r *MySQLSagaRepository) FindExecutionsBySagaID(
	ctx context.Context,
	sagaID string,
) ([]*SagaExecution, error) {
	var records []SagaExecutionRecord
	err := r.db.WithContext(ctx).
		Where("saga_id = ?", sagaID).
		Order("started_at DESC").
		Find(&records).
		Error

	if err != nil {
		return nil, err
	}

	executions := make([]*SagaExecution, 0, len(records))
	for _, record := range records {
		execution, err := record.ToExecution()
		if err != nil {
			return nil, err
		}
		executions = append(executions, execution)
	}

	return executions, nil
}

// ============ 数据库记录定义 ============

// SagaRecord Saga定义表
type SagaRecord struct {
	ID             string `gorm:"primaryKey;size:64"`
	Name           string `gorm:"unique;not null;size:128"`
	Description    string `gorm:"type:text"`
	Steps          string `gorm:"type:text"`
	Compensations  string `gorm:"type:text"`
	RetryPolicy    string `gorm:"type:text"`
	TimeoutSeconds int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TableName 指定表名
func (SagaRecord) TableName() string {
	return "saga_definitions"
}

// ToSaga 转换为Saga对象
func (r *SagaRecord) ToSaga() (*Saga, error) {
	saga := &Saga{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Timeout:     time.Duration(r.TimeoutSeconds) * time.Second,
	}

	// 反序列化步骤
	if err := json.Unmarshal([]byte(r.Steps), &saga.Steps); err != nil {
		return nil, fmt.Errorf("反序列化步骤失败: %w", err)
	}

	// 反序列化补偿步骤
	if err := json.Unmarshal([]byte(r.Compensations), &saga.Compensations); err != nil {
		return nil, fmt.Errorf("反序列化补偿步骤失败: %w", err)
	}

	// 反序列化重试策略
	if r.RetryPolicy != "" {
		if err := json.Unmarshal([]byte(r.RetryPolicy), &saga.RetryPolicy); err != nil {
			return nil, fmt.Errorf("反序列化重试策略失败: %w", err)
		}
	} else {
		saga.RetryPolicy = DefaultRetryPolicy()
	}

	return saga, nil
}

// SagaExecutionRecord Saga执行记录表
type SagaExecutionRecord struct {
	ID             string     `gorm:"primaryKey;size:64"`
	SagaID         string     `gorm:"not null;index:idx_saga_id;size:64"`
	Status         string     `gorm:"not null;index:idx_status;size:32"`
	CurrentStep    int
	InputData      string `gorm:"type:text"`
	OutputData     string `gorm:"type:text"`
	Error          string `gorm:"type:text"`
	StepExecutions string `gorm:"type:text"`
	StartedAt      time.Time `gorm:"index:idx_started_at"`
	CompletedAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TableName 指定表名
func (SagaExecutionRecord) TableName() string {
	return "saga_executions"
}

// ToExecution 转换为SagaExecution对象
func (r *SagaExecutionRecord) ToExecution() (*SagaExecution, error) {
	execution := &SagaExecution{
		ID:          r.ID,
		SagaID:      r.SagaID,
		Status:      SagaStatus(r.Status),
		CurrentStep: r.CurrentStep,
		StartedAt:   r.StartedAt,
		CompletedAt: r.CompletedAt,
	}

	// 反序列化输入数据
	if r.InputData != "" {
		if err := json.Unmarshal([]byte(r.InputData), &execution.InputData); err != nil {
			return nil, fmt.Errorf("反序列化输入数据失败: %w", err)
		}
	}

	// 反序列化输出数据
	if r.OutputData != "" {
		if err := json.Unmarshal([]byte(r.OutputData), &execution.OutputData); err != nil {
			return nil, fmt.Errorf("反序列化输出数据失败: %w", err)
		}
	}

	// 反序列化步骤执行记录
	if r.StepExecutions != "" {
		if err := json.Unmarshal([]byte(r.StepExecutions), &execution.StepExecutions); err != nil {
			return nil, fmt.Errorf("反序列化步骤执行记录失败: %w", err)
		}
	}

	// 设置错误
	if r.Error != "" {
		execution.Error = fmt.Errorf("%s", r.Error)
	}

	return execution, nil
}

// AutoMigrate 自动迁移表结构
func (r *MySQLSagaRepository) AutoMigrate() error {
	return r.db.AutoMigrate(
		&SagaRecord{},
		&SagaExecutionRecord{},
	)
}
