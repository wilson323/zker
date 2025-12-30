// backend/domain/tenant/migration/migrate_test.go
package migration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 创建测试表结构
	err = db.Exec(`
		CREATE TABLE tenants (
			tenant_id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(100)
		)
	`).Error
	assert.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE users (
			user_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36),
			username VARCHAR(100)
		)
	`).Error
	assert.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE bots (
			id VARCHAR(36) PRIMARY KEY,
			creator_id VARCHAR(36),
			tenant_id VARCHAR(36),
			name VARCHAR(100)
		)
	`).Error
	assert.NoError(t, err)

	// 插入测试数据
	db.Exec("INSERT INTO tenants (tenant_id, name) VALUES (?, ?)", "tenant_1", "Test Tenant")
	db.Exec("INSERT INTO users (user_id, tenant_id, username) VALUES (?, ?, ?)", "user_1", "tenant_1", "testuser")

	for i := 1; i <= 100; i++ {
		db.Exec("INSERT INTO bots (id, creator_id, name) VALUES (?, ?, ?)",
			fmt.Sprintf("bot_%d", i), "user_1", fmt.Sprintf("Bot %d", i))
	}

	return db
}

// TestMigrateTable 测试表迁移功能
func TestMigrateTable(t *testing.T) {
	db := setupTestDB(t)
	migrator := NewTenantIDMigrator(db)
	ctx := context.Background()

	// 1. 执行迁移
	err := migrator.MigrateTable(ctx, "bots", "creator_id")
	assert.NoError(t, err)

	// 2. 验证迁移结果
	var count int64
	db.Raw("SELECT COUNT(*) FROM bots WHERE tenant_id IS NOT NULL").Scan(&count)
	assert.Equal(t, int64(100), count, "所有记录都应该已迁移")

	// 3. 验证tenant_id值正确
	var tenantID string
	db.Raw("SELECT tenant_id FROM bots LIMIT 1").Scan(&tenantID)
	assert.Equal(t, "tenant_1", tenantID, "tenant_id应该正确")
}

// TestValidateMigration 测试迁移验证功能
func TestValidateMigration(t *testing.T) {
	db := setupTestDB(t)
	migrator := NewTenantIDMigrator(db)
	ctx := context.Background()

	// 1. 执行迁移
	err := migrator.MigrateTable(ctx, "bots", "creator_id")
	assert.NoError(t, err)

	// 2. 验证应该通过
	err = migrator.ValidateMigration(ctx, "bots")
	assert.NoError(t, err)

	// 3. 添加一条无效数据
	db.Exec("INSERT INTO bots (id, creator_id, tenant_id) VALUES (?, ?, ?)",
		"bot_invalid", "user_1", "invalid_tenant")

	// 4. 验证应该失败
	err = migrator.ValidateMigration(ctx, "bots")
	assert.Error(t, err)
}

// TestRollbackMigration 测试回滚功能
func TestRollbackMigration(t *testing.T) {
	db := setupTestDB(t)
	migrator := NewTenantIDMigrator(db)
	ctx := context.Background()

	// 1. 执行迁移
	err := migrator.MigrateTable(ctx, "bots", "creator_id")
	assert.NoError(t, err)

	// 2. 验证已迁移
	var count int64
	db.Raw("SELECT COUNT(*) FROM bots WHERE tenant_id IS NOT NULL").Scan(&count)
	assert.Equal(t, int64(100), count)

	// 3. 执行回滚
	err = migrator.RollbackMigration(ctx, "bots")
	assert.NoError(t, err)

	// 4. 验证已回滚
	db.Raw("SELECT COUNT(*) FROM bots WHERE tenant_id IS NULL").Scan(&count)
	assert.Equal(t, int64(100), count, "所有记录应该已回滚")
}

// TestGetMigrationStats 测试统计功能
func TestGetMigrationStats(t *testing.T) {
	db := setupTestDB(t)
	migrator := NewTenantIDMigrator(db)
	ctx := context.Background()

	// 1. 迁移50条数据
	db.Exec("UPDATE bots SET tenant_id = ? WHERE id <= ?", "tenant_1", 50)

	// 2. 获取统计信息
	stats, err := migrator.GetMigrationStats(ctx, "bots")
	assert.NoError(t, err)

	// 3. 验证统计数据
	assert.Equal(t, int64(100), stats.TotalRecords)
	assert.Equal(t, int64(50), stats.MigratedRecords)
	assert.Equal(t, int64(50), stats.UnmigratedRecords)
	assert.Equal(t, 50.0, stats.Progress)
}

// BenchmarkMigrateTable 性能基准测试
func BenchmarkMigrateTable(b *testing.B) {
	db := setupTestDB(&testing.T{})
	migrator := NewTenantIDMigrator(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 重置数据
		db.Exec("UPDATE bots SET tenant_id = NULL")

		// 执行迁移
		_ = migrator.MigrateTable(ctx, "bots", "creator_id")
	}
}
