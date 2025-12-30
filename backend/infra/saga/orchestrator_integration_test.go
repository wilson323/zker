//go:build integration
// +build integration

package saga

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupMySQLTestContainer 启动MySQL测试容器
func setupMySQLTestContainer(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	// 启动MySQL容器
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mysql:8.4.5",
			ExposedPorts: []string{"3306/tcp", "33060/tcp"},
			Env: map[string]string{
				"MYSQL_ROOT_PASSWORD": "root",
				"MYSQL_DATABASE":      "test_saga",
			},
			WaitingFor: wait.ForLog("port: 3306  MySQL Community Server"),
		},
		Started: true,
	})
	require.NoError(t, err)

	// 获取数据库连接信息
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "3306")
	require.NoError(t, err)

	// 连接数据库
	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/test_saga?charset=utf8mb4&parseTime=True&loc=Local",
		host, port.Port())

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移表结构
	repo := NewMySQLSagaRepository(db)
	err = repo.AutoMigrate()
	require.NoError(t, err)

	// 返回清理函数
	cleanup := func() {
		db.Exec("DROP TABLE IF EXISTS saga_executions")
		db.Exec("DROP TABLE IF EXISTS saga_definitions")
		container.Terminate(ctx)
	}

	return db, cleanup
}

// TestMySQLSagaRepository_Integration 测试MySQL仓储集成
func TestMySQLSagaRepository_Integration(t *testing.T) {
	db, cleanup := setupMySQLTestContainer(t)
	defer cleanup()

	repo := NewMySQLSagaRepository(db)
	ctx := context.Background()

	t.Run("保存和查询Saga定义", func(t *testing.T) {
		// 创建测试Saga
		saga := &Saga{
			ID:          "integration-test-saga",
			Name:        "集成测试Saga",
			Description: "测试MySQL仓储",
			Steps: []SagaStep{
				&MockStep{name: "step1", timeout: 5 * time.Second},
			},
			Compensations: []CompensationStep{
				&MockCompensationStep{name: "comp1", timeout: 5 * time.Second},
			},
			Timeout:     30 * time.Second,
			RetryPolicy: DefaultRetryPolicy(),
		}

		// 保存Saga
		err := repo.SaveSaga(ctx, saga)
		require.NoError(t, err)

		// 根据ID查询
		found, err := repo.FindSagaByID(ctx, saga.ID)
		require.NoError(t, err)
		assert.Equal(t, saga.ID, found.ID)
		assert.Equal(t, saga.Name, found.Name)
		assert.Equal(t, 1, len(found.Steps))
		assert.Equal(t, 1, len(found.Compensations))

		// 根据名称查询
		foundByName, err := repo.FindSagaByName(ctx, saga.Name)
		require.NoError(t, err)
		assert.Equal(t, saga.ID, foundByName.ID)

		// 列出所有Saga
		sagas, err := repo.ListSagas(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(sagas), 1)
	})

	t.Run("保存和更新执行记录", func(t *testing.T) {
		// 创建测试执行记录
		execution := &SagaExecution{
			ID:          generateID(),
			SagaID:      "integration-test-saga",
			Status:      SagaStatusRunning,
			CurrentStep: 0,
			InputData:   "test-input",
			StartedAt:   time.Now(),
			StepExecutions: []StepExecution{
				{
					StepName:  "step1",
					Status:    StepCompleted,
					StartedAt: time.Now(),
				},
			},
		}

		// 保存执行记录
		err := repo.SaveExecution(ctx, execution)
		require.NoError(t, err)

		// 查询执行记录
		found, err := repo.FindExecutionByID(ctx, execution.ID)
		require.NoError(t, err)
		assert.Equal(t, execution.ID, found.ID)
		assert.Equal(t, execution.SagaID, found.SagaID)
		assert.Equal(t, SagaStatusRunning, found.Status)

		// 更新执行记录
		execution.Status = SagaStatusCompleted
		execution.CurrentStep = 1
		completedAt := time.Now()
		execution.CompletedAt = &completedAt

		err = repo.UpdateExecution(ctx, execution)
		require.NoError(t, err)

		// 验证更新
		updated, err := repo.FindExecutionByID(ctx, execution.ID)
		require.NoError(t, err)
		assert.Equal(t, SagaStatusCompleted, updated.Status)
		assert.Equal(t, 1, updated.CurrentStep)
		assert.NotNil(t, updated.CompletedAt)

		// 根据SagaID查询执行记录
		executions, err := repo.FindExecutionsBySagaID(ctx, execution.SagaID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(executions), 1)
	})
}

// TestSagaOrchestrator_Integration 测试Saga协调器集成
func TestSagaOrchestrator_Integration(t *testing.T) {
	db, cleanup := setupMySQLTestContainer(t)
	defer cleanup()

	repo := NewMySQLSagaRepository(db)
	logger := zap.NewNop()
	orchestrator := NewSagaOrchestrator(repo, logger)
	ctx := context.Background()

	t.Run("完整Saga执行流程", func(t *testing.T) {
		// 定义Saga
		saga := &Saga{
			ID:          "integration-execution-saga",
			Name:        "集成执行测试Saga",
			Description: "测试完整的Saga执行流程",
			Steps: []SagaStep{
				&MockStep{
					name:    "create_tenant",
					timeout: 5 * time.Second,
					executeFunc: func(ctx context.Context, data interface{}) (interface{}, error) {
						return map[string]interface{}{
							"tenant_id": generateID(),
							"name":      "测试租户",
						}, nil
					},
				},
				&MockStep{
					name:    "create_organization",
					timeout: 5 * time.Second,
					executeFunc: func(ctx context.Context, data interface{}) (interface{}, error) {
						tenantData := data.(map[string]interface{})
						return map[string]interface{}{
							"org_id":    generateID(),
							"tenant_id": tenantData["tenant_id"],
							"name":      "默认组织",
						}, nil
					},
				},
			},
			Compensations: []CompensationStep{
				&MockCompensationStep{
					name:    "delete_tenant",
					timeout: 5 * time.Second,
					compensateFunc: func(ctx context.Context, data interface{}) error {
						fmt.Printf("[补偿] 删除租户: %v\n", data)
						return nil
					},
				},
				&MockCompensationStep{
					name:    "delete_organization",
					timeout: 5 * time.Second,
					compensateFunc: func(ctx context.Context, data interface{}) error {
						fmt.Printf("[补偿] 删除组织: %v\n", data)
						return nil
					},
				},
			},
			Timeout:     30 * time.Second,
			RetryPolicy: DefaultRetryPolicy(),
		}

		// 注册Saga
		err := orchestrator.DefineSaga(saga)
		require.NoError(t, err)

		// 执行Saga
		input := map[string]interface{}{
			"enterprise_name": "测试公司",
			"contact_email":   "test@example.com",
		}

		execution, err := orchestrator.ExecuteSaga(ctx, "集成执行测试Saga", input)
		require.NoError(t, err)

		// 验证执行结果
		assert.Equal(t, SagaStatusCompleted, execution.Status)
		assert.Equal(t, 2, len(execution.StepExecutions))
		assert.Equal(t, 2, execution.CurrentStep)

		// 验证每个步骤
		for i, stepExec := range execution.StepExecutions {
			assert.Equal(t, StepCompleted, stepExec.Status)
			assert.NotNil(t, stepExec.CompletedAt)
			assert.Nil(t, stepExec.Error)

			if i == 0 {
				assert.Equal(t, "create_tenant", stepExec.StepName)
			} else if i == 1 {
				assert.Equal(t, "create_organization", stepExec.StepName)
			}
		}

		// 验证可以从数据库查询执行记录
		found, err := orchestrator.GetStatus(ctx, execution.ID)
		require.NoError(t, err)
		assert.Equal(t, execution.ID, found.ID)
		assert.Equal(t, SagaStatusCompleted, found.Status)
	})

	t.Run("Saga失败和补偿", func(t *testing.T) {
		// 定义会失败的Saga
		saga := &Saga{
			ID:          "integration-fail-saga",
			Name:        "集成失败测试Saga",
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
						return nil, fmt.Errorf("step2 failed")
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
						fmt.Printf("[补偿] 补偿步骤1: %v\n", data)
						return nil
					},
				},
				&MockCompensationStep{
					name:    "comp2",
					timeout: 5 * time.Second,
					compensateFunc: func(ctx context.Context, data interface{}) error {
						fmt.Printf("[补偿] 补偿步骤2\n")
						return nil
					},
				},
				&MockCompensationStep{
					name:    "comp3",
					timeout: 5 * time.Second,
					compensateFunc: func(ctx context.Context, data interface{}) error {
						fmt.Printf("[补偿] 补偿步骤3\n")
						return nil
					},
				},
			},
			Timeout:     30 * time.Second,
			RetryPolicy: &RetryPolicy{
				MaxAttempts:     2,
				InitialInterval: 10 * time.Millisecond,
				MaxInterval:     50 * time.Millisecond,
				Multiplier:      2.0,
			},
		}

		// 注册Saga
		err := orchestrator.DefineSaga(saga)
		require.NoError(t, err)

		// 执行Saga
		execution, err := orchestrator.ExecuteSaga(ctx, "集成失败测试Saga", "input")

		// 验证失败和补偿
		assert.Error(t, err)
		assert.NotNil(t, execution)
		assert.Equal(t, SagaStatusCompensated, execution.Status)
		assert.Equal(t, 2, len(execution.StepExecutions))

		// 验证步骤状态
		assert.Equal(t, StepCompleted, execution.StepExecutions[0].Status)
		assert.Equal(t, StepFailed, execution.StepExecutions[1].Status)
		assert.Equal(t, 1, execution.CurrentStep)

		// 验证可以从数据库查询补偿后的状态
		found, err := orchestrator.GetStatus(ctx, execution.ID)
		require.NoError(t, err)
		assert.Equal(t, SagaStatusCompensated, found.Status)
	})
}
