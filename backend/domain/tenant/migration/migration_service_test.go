// backend/domain/tenant/migration/migration_service_test.go
package migration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestMigrationService_DoubleWriteMode 测试双写模式
func TestMigrationService_DoubleWriteMode(t *testing.T) {
	// 使用testcontainers启动MySQL
	db, cleanup := setupMySQLTestContainer(t)
	defer cleanup()

	// 1. 准备测试数据
	ctx := context.Background()
	setupTestData(t, ctx, db)

	// 2. 创建迁移服务
	service := NewMigrationService(db)

	// 3. 启动双写模式
	config := &MigrationConfig{
		TableNames: []string{"bots", "conversations"},
		Mode:       "double_write",
	}

	err := service.StartMigration(ctx, config)
	require.NoError(t, err, "启动双写模式应该成功")

	// 4. 验证双写正确性
	verifier := NewDualWriteVerifier(db)

	// 等待一段时间，让双写生效
	time.Sleep(2 * time.Second)

	// 验证数据一致性
	err = verifier.VerifyDualWrite(ctx, "bots", 1*time.Hour)
	assert.NoError(t, err, "数据一致性应该通过")

	// 5. 检查迁移进度
	progress, err := service.GetProgress(ctx, config.MigrationID)
	assert.NoError(t, err)
	assert.Equal(t, float64(100), progress.ProgressPercent, "迁移应该100%完成")
}

// TestMigrationService_Rollback 测试回滚机制
func TestMigrationService_Rollback(t *testing.T) {
	db, cleanup := setupMySQLTestContainer(t)
	defer cleanup()

	ctx := context.Background()
	setupTestData(t, ctx, db)

	service := NewMigrationService(db)

	// 1. 启动迁移
	config := &MigrationConfig{
		TableNames: []string{"bots"},
		Mode:       "double_write",
	}

	err := service.StartMigration(ctx, config)
	require.NoError(t, err)

	// 2. 执行回滚
	err = service.Rollback(ctx, config.MigrationID)
	assert.NoError(t, err, "回滚应该成功")

	// 3. 验证回滚后数据一致性
	// tenant_id字段应该被删除
	var columnExists int64
	err = db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'bots'
		  AND COLUMN_NAME = 'tenant_id'
	`).Scan(&columnExists).Error

	assert.NoError(t, err)
	assert.Equal(t, int64(0), columnExists, "tenant_id字段应该被删除")
}

// TestMigrationService_Validation 测试迁移验证
func TestMigrationService_Validation(t *testing.T) {
	db, cleanup := setupMySQLTestContainer(t)
	defer cleanup()

	ctx := context.Background()
	setupTestData(t, ctx, db)

	service := NewMigrationService(db)

	// 启动迁移
	config := &MigrationConfig{
		TableNames: []string{"bots", "conversations"},
		Mode:       "double_write",
	}

	err := service.StartMigration(ctx, config)
	require.NoError(t, err)

	// 等待迁移完成
	time.Sleep(3 * time.Second)

	// 验证迁移结果
	err = service.Validate(ctx, config.MigrationID)
	assert.NoError(t, err, "迁移验证应该通过")
}

// TestMigrationService_ConcurrentMigrations 测试并发迁移
func TestMigrationService_ConcurrentMigrations(t *testing.T) {
	db, cleanup := setupMySQLTestContainer(t)
	defer cleanup()

	ctx := context.Background()
	setupTestData(t, ctx, db)

	service := NewMigrationService(db)

	// 并发启动多个迁移任务
	tables := []string{"bots", "conversations", "knowledge", "workflows"}
	var migrationIDs []string

	for _, table := range tables {
		config := &MigrationConfig{
			TableNames: []string{table},
			Mode:       "double_write",
		}

		err := service.StartMigration(ctx, config)
		require.NoError(t, err)

		migrationIDs = append(migrationIDs, config.MigrationID)
	}

	// 等待所有迁移完成
	time.Sleep(5 * time.Second)

	// 验证所有迁移都成功
	for _, id := range migrationIDs {
		progress, err := service.GetProgress(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, float64(100), progress.ProgressPercent, fmt.Sprintf("迁移 %s 应该100%%完成", id))
	}
}

// ================================================================================
// 辅助函数
// ================================================================================

// setupMySQLTestContainer 设置MySQL测试容器
func setupMySQLTestContainer(t *testing.T) (*gorm.DB, func()) {
	// 使用testcontainers-for-go启动MySQL容器
	// 这里简化为连接本地MySQL（实际应该使用testcontainers）

	dsn := "root:password@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err, "连接数据库应该成功")

	// 清理函数
	cleanup := func() {
		// 清理测试数据
		db.Exec("DROP TABLE IF EXISTS bots")
		db.Exec("DROP TABLE IF EXISTS conversations")
		db.Exec("DROP TABLE IF EXISTS knowledge")
		db.Exec("DROP TABLE IF EXISTS workflows")
		db.Exec("DROP TABLE IF EXISTS users")
		db.Exec("DROP TABLE IF EXISTS tenants")
	}

	return db, cleanup
}

// setupTestData 准备测试数据
func setupTestData(t *testing.T, ctx context.Context, db *gorm.DB) {
	// 1. 创建租户表
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			tenant_id VARCHAR(36) PRIMARY KEY,
			tenant_name VARCHAR(200) NOT NULL,
			tenant_type VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL,
			subscription_tier VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`).Error
	require.NoError(t, err)

	// 2. 创建用户表
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id VARCHAR(36) NOT NULL UNIQUE,
			tenant_id VARCHAR(36) NOT NULL DEFAULT 'system_tenant',
			username VARCHAR(100) NOT NULL,
			email VARCHAR(200) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_tenant_id (tenant_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`).Error
	require.NoError(t, err)

	// 3. 创建bots表
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bots (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			bot_id VARCHAR(36) NOT NULL UNIQUE,
			tenant_id VARCHAR(36) NOT NULL DEFAULT 'system_tenant',
			creator_id VARCHAR(36) NOT NULL,
			name VARCHAR(200) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_tenant_id (tenant_id),
			INDEX idx_creator_id (creator_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`).Error
	require.NoError(t, err)

	// 4. 创建conversations表
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS conversations (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			conv_id VARCHAR(36) NOT NULL UNIQUE,
			tenant_id VARCHAR(36) NOT NULL DEFAULT 'system_tenant',
			creator_id VARCHAR(36) NOT NULL,
			bot_id VARCHAR(36) NOT NULL,
			title VARCHAR(200),
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_tenant_id (tenant_id),
			INDEX idx_creator_id (creator_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`).Error
	require.NoError(t, err)

	// 5. 创建其他表...
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS knowledge (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			knowledge_id VARCHAR(36) NOT NULL UNIQUE,
			tenant_id VARCHAR(36) NOT NULL DEFAULT 'system_tenant',
			creator_id VARCHAR(36) NOT NULL,
			name VARCHAR(200) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_tenant_id (tenant_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`).Error
	require.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS workflows (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			workflow_id VARCHAR(36) NOT NULL UNIQUE,
			tenant_id VARCHAR(36) NOT NULL DEFAULT 'system_tenant',
			creator_id VARCHAR(36) NOT NULL,
			name VARCHAR(200) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_tenant_id (tenant_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`).Error
	require.NoError(t, err)

	// 6. 插入测试数据
	// 插入租户
	err = db.Exec(`
		INSERT INTO tenants (tenant_id, tenant_name, tenant_type, status, subscription_tier)
		VALUES
			('tenant_1', 'Test Tenant 1', 'individual', 'active', 'free'),
			('tenant_2', 'Test Tenant 2', 'enterprise', 'active', 'pro')
	`).Error
	require.NoError(t, err)

	// 插入用户
	err = db.Exec(`
		INSERT INTO users (user_id, tenant_id, username, email)
		VALUES
			('user_1', 'tenant_1', 'user1', 'user1@example.com'),
			('user_2', 'tenant_2', 'user2', 'user2@example.com')
	`).Error
	require.NoError(t, err)

	// 插入bots
	for i := 1; i <= 100; i++ {
		tenantID := "tenant_1"
		creatorID := "user_1"
		if i > 50 {
			tenantID = "tenant_2"
			creatorID = "user_2"
		}

		err = db.Exec(`
			INSERT INTO bots (bot_id, tenant_id, creator_id, name, status)
			VALUES (?, ?, ?, ?, ?)
		`, fmt.Sprintf("bot_%d", i), tenantID, creatorID, fmt.Sprintf("Bot %d", i), "active").Error
		require.NoError(t, err)
	}

	// 插入conversations
	for i := 1; i <= 100; i++ {
		tenantID := "tenant_1"
		creatorID := "user_1"
		if i > 50 {
			tenantID = "tenant_2"
			creatorID = "user_2"
		}

		err = db.Exec(`
			INSERT INTO conversations (conv_id, tenant_id, creator_id, bot_id, title, status)
			VALUES (?, ?, ?, ?, ?, ?)
		`, fmt.Sprintf("conv_%d", i), tenantID, creatorID, "bot_1", fmt.Sprintf("Conversation %d", i), "active").Error
		require.NoError(t, err)
	}

	// 插入knowledge
	for i := 1; i <= 50; i++ {
		err = db.Exec(`
			INSERT INTO knowledge (knowledge_id, tenant_id, creator_id, name)
			VALUES (?, ?, ?, ?)
		`, fmt.Sprintf("knowledge_%d", i), "tenant_1", "user_1", fmt.Sprintf("Knowledge %d", i)).Error
		require.NoError(t, err)
	}

	// 插入workflows
	for i := 1; i <= 50; i++ {
		err = db.Exec(`
			INSERT INTO workflows (workflow_id, tenant_id, creator_id, name)
			VALUES (?, ?, ?, ?)
		`, fmt.Sprintf("workflow_%d", i), "tenant_1", "user_1", fmt.Sprintf("Workflow %d", i)).Error
		require.NoError(t, err)
	}

	logs.CtxInfof(ctx, "测试数据准备完成")
}

// ================================================================================
// 基准测试
// ================================================================================

// BenchmarkMigration_DoubleWrite 双写模式基准测试
func BenchmarkMigration_DoubleWrite(b *testing.B) {
	db, cleanup := setupMySQLTestContainer(&testing.T{})
	defer cleanup()

	ctx := context.Background()
	setupTestData(&testing.T{}, ctx, db)

	service := NewMigrationService(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &MigrationConfig{
			TableNames: []string{"bots"},
			Mode:       "double_write",
		}

		err := service.StartMigration(ctx, config)
		if err != nil {
			b.Fatalf("迁移失败: %v", err)
		}
	}
}

// BenchmarkDualWrite_Verification 双写验证基准测试
func BenchmarkDualWrite_Verification(b *testing.B) {
	db, cleanup := setupMySQLTestContainer(&testing.T{})
	defer cleanup()

	ctx := context.Background()
	setupTestData(&testing.T{}, ctx, db)

	verifier := NewDualWriteVerifier(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := verifier.VerifyDualWrite(ctx, "bots", 1*time.Hour)
		if err != nil {
			b.Fatalf("验证失败: %v", err)
		}
	}
}
