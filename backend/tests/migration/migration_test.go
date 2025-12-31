// backend/tests/migration/migration_test.go
// 数据迁移自动化测试套件
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MigrationTestSuite 迁移测试套件
type MigrationTestSuite struct {
	suite.Suite
	db              *gorm.DB
	sqlDB           *sql.DB
	testDBName      string
	migrationTool   *TestMigrationTool
	validator       *TestDataValidator
	benchmark       *TestBenchmark
}

// TestMigrationTool 测试迁移工具
type TestMigrationTool struct {
	db *gorm.DB
}

// TestDataValidator 测试数据验证器
type TestDataValidator struct {
	db *gorm.DB
}

// TestBenchmark 测试基准测试
type TestBenchmark struct {
	db *gorm.DB
}

// TestConfig 测试配置
type TestConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	TestDBName string
	Verbose    bool
}

var config *TestConfig

func TestMain(m *testing.M) {
	// 解析命令行参数
	config = parseTestFlags()

	// 运行测试
	code := m.Run()

	// 清理测试数据库
	if config.TestDBName != "" {
		cleanupTestDatabase()
	}

	os.Exit(code)
}

// parseTestFlags 解析测试参数
func parseTestFlags() *TestConfig {
	cfg := &TestConfig{}

	flag.StringVar(&cfg.DBHost, "host", "localhost", "MySQL主机地址")
	flag.IntVar(&cfg.DBPort, "port", 3306, "MySQL端口")
	flag.StringVar(&cfg.DBUser, "user", "root", "MySQL用户名")
	flag.StringVar(&cfg.DBPassword, "password", "", "MySQL密码")
	flag.StringVar(&cfg.TestDBName, "database", "test_migration", "测试数据库名称")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "详细输出")

	flag.Parse()

	return cfg
}

// SetupSuite 测试套件初始化
func (s *MigrationTestSuite) SetupSuite() {
	require.NoError(s.T(), setupTestDatabase(config.TestDBName))

	db, err := gorm.Open(mysql.Open(formatDSN("", config.TestDBName)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(s.T(), err)

	s.db = db
	s.sqlDB, err = db.DB()
	require.NoError(s.T(), err)

	s.migrationTool = &TestMigrationTool{db: db}
	s.validator = &TestDataValidator{db: db}
	s.benchmark = &TestBenchmark{db: db}

	log.Printf("测试套件初始化完成: %s", config.TestDBName)
}

// TearDownSuite 测试套件清理
func (MigrationTestSuite) TearDownSuite() {
	log.Println("测试套件清理完成")
}

// SetupTest 每个测试前的准备
func (s *MigrationTestSuite) SetupTest() {
	// 清理测试数据
	s.cleanupTestData()

	// 准备测试数据
	s.prepareTestData()
}

// TearDownTest 每个测试后的清理
func (s *MigrationTestSuite) TearDownTest() {
	// 清理测试数据
	s.cleanupTestData()
}

// ============================================================================
// 测试用例
// ============================================================================

// TestMigration_AddTenantIDColumn 测试添加tenant_id字段
func (s *MigrationTestSuite) TestMigration_AddTenantIDColumn() {
	ctx := context.Background()

	// 执行迁移
	err := s.migrationTool.AddTenantIDColumn(ctx, "test_bots")
	require.NoError(s.T(), err, "添加字段应该成功")

	// 验证字段存在
	var columnExists int64
	s.db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ?
		  AND TABLE_NAME = 'test_bots'
		  AND COLUMN_NAME = 'tenant_id'
	`, config.TestDBName).Scan(&columnExists)

	assert.Equal(s.T(), int64(1), columnExists, "tenant_id字段应该存在")

	// 验证索引存在
	var indexExists int64
	s.db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ?
		  AND TABLE_NAME = 'test_bots'
		  AND INDEX_NAME = 'idx_tenant_id'
	`, config.TestDBName).Scan(&indexExists)

	assert.Equal(s.T(), int64(1), indexExists, "idx_tenant_id索引应该存在")
}

// TestMigration_BackfillTenantID 测试数据回填
func (s *MigrationTestSuite) TestMigration_BackfillTenantID() {
	ctx := context.Background()

	// 1. 添加字段
	err := s.migrationTool.AddTenantIDColumn(ctx, "test_bots")
	require.NoError(s.T(), err)

	// 2. 回填数据
	result, err := s.migrationTool.BackfillTenantID(ctx, "test_bots", "test_tenant", 100)
	require.NoError(s.T(), err, "数据回填应该成功")

	// 3. 验证回填结果
	assert.Equal(s.T(), int64(100), result.MigratedCount, "应该迁移100条记录")
	assert.Equal(s.T(), int64(0), result.FailedCount, "不应该有失败的记录")

	// 4. 验证数据完整性
	validation := s.validator.ValidateTenantID(ctx, "test_bots")
	assert.Equal(s.T(), int64(100), validation.TotalRecords, "总记录数应该是100")
	assert.Equal(s.T(), int64(100), validation.ValidRecords, "有效记录数应该是100")
	assert.Equal(s.T(), int64(0), validation.NullRecords, "NULL记录数应该是0")
	assert.True(s.T(), validation.AllValid, "所有记录都应该有效")
}

// TestMigration_DoubleWriteMode 测试双写模式
func (s *MigrationTestSuite) TestMigration_DoubleWriteMode() {
	ctx := context.Background()

	// 1. 添加字段
	err := s.migrationTool.AddTenantIDColumn(ctx, "test_bots")
	require.NoError(s.T(), err)

	// 2. 启动双写
	err = s.migrationTool.EnableDualWrite(ctx, "test_bots")
	require.NoError(s.T(), err)

	// 3. 插入测试数据（应该双写）
	for i := 0; i < 10; i++ {
		bot := map[string]interface{}{
			"name":       fmt.Sprintf("bot_%d", i),
			"creator_id": "test_user",
			"tenant_id":  "test_tenant", // 应用层写入
			"status":     "active",
		}
		s.db.Table("test_bots").Create(bot)
	}

	// 4. 验证双写正确性
	var count int64
	s.db.Table("test_bots").Where("tenant_id = ?", "test_tenant").Count(&count)
	assert.Equal(s.T(), int64(10), count, "应该有10条记录的tenant_id为test_tenant")
}

// TestMigration_Performance 测试迁移性能
func (s *MigrationTestSuite) TestMigration_Performance() {
	ctx := context.Background()

	// 准备大规模测试数据（10000条）
	s.prepareLargeScaleData("test_bots", 10000)

	// 添加字段
	err := s.migrationTool.AddTenantIDColumn(ctx, "test_bots")
	require.NoError(s.T(), err)

	// 测试回填性能
	startTime := time.Now()
	result, err := s.migrationTool.BackfillTenantID(ctx, "test_bots", "test_tenant", 1000)
	duration := time.Since(startTime)

	require.NoError(s.T(), err)

	// 验证性能指标
	qps := float64(result.MigratedCount) / duration.Seconds()
	log.Printf("迁移性能: %d 条记录, 耗时: %v, QPS: %.2f", result.MigratedCount, duration, qps)

	assert.Greater(s.T(), qps, 1000.0, "QPS应该大于1000")
	assert.Equal(s.T(), int64(10000), result.MigratedCount, "应该迁移所有记录")
}

// TestMigration_Rollback 测试回滚
func (s *MigrationTestSuite) TestMigration_Rollback() {
	ctx := context.Background()

	// 1. 添加字段
	err := s.migrationTool.AddTenantIDColumn(ctx, "test_bots")
	require.NoError(s.T(), err)

	// 2. 回填数据
	_, err = s.migrationTool.BackfillTenantID(ctx, "test_bots", "test_tenant", 100)
	require.NoError(s.T(), err)

	// 3. 执行回滚
	err = s.migrationTool.Rollback(ctx, "test_bots")
	require.NoError(s.T(), err, "回滚应该成功")

	// 4. 验证字段已删除
	var columnExists int64
	s.db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ?
		  AND TABLE_NAME = 'test_bots'
		  AND COLUMN_NAME = 'tenant_id'
	`, config.TestDBName).Scan(&columnExists)

	assert.Equal(s.T(), int64(0), columnExists, "tenant_id字段应该已删除")
}

// TestMigration_ConcurrentMigration 测试并发迁移
func (s *MigrationTestSuite) TestMigration_ConcurrentMigration() {
	ctx := context.Background()

	// 准备多个测试表
	tables := []string{"test_bots", "test_conversations", "test_knowledge"}
	for _, tableName := range tables {
		s.createTestTable(tableName)
		s.prepareTestData(tableName, 1000)
	}

	// 并发迁移
	results := make(map[string]*MigrationResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, tableName := range tables {
		wg.Add(1)
		go func(tbl string) {
			defer wg.Done()

			// 添加字段
			err := s.migrationTool.AddTenantIDColumn(ctx, tbl)
			assert.NoError(s.T(), err)

			// 回填数据
			result, err := s.migrationTool.BackfillTenantID(ctx, tbl, "test_tenant", 100)
			assert.NoError(s.T(), err)

			mu.Lock()
			results[tbl] = result
			mu.Unlock()
		}(tableName)
	}

	wg.Wait()

	// 验证所有表都迁移成功
	for _, tableName := range tables {
		result, exists := results[tableName]
		assert.True(s.T(), exists, "表 %s 应该有迁移结果", tableName)
		assert.Equal(s.T(), int64(1000), result.MigratedCount, "表 %s 应该迁移1000条记录", tableName)
		assert.Equal(s.T(), int64(0), result.FailedCount, "表 %s 不应该有失败记录", tableName)
	}
}

// TestMigration_DataIntegrity 测试数据完整性
func (s *MigrationTestSuite) TestMigration_DataIntegrity() {
	ctx := context.Background()

	// 1. 准备测试数据
	initialCount := int64(1000)
	s.prepareTestData("test_bots", initialCount)

	// 2. 执行迁移
	err := s.migrationTool.AddTenantIDColumn(ctx, "test_bots")
	require.NoError(s.T(), err)

	_, err = s.migrationTool.BackfillTenantID(ctx, "test_bots", "test_tenant", 100)
	require.NoError(s.T(), err)

	// 3. 验证记录数一致
	var finalCount int64
	s.db.Table("test_bots").Count(&finalCount)
	assert.Equal(s.T(), initialCount, finalCount, "迁移前后记录数应该一致")

	// 4. 验证数据内容
	var bots []map[string]interface{}
	s.db.Table("test_bots").Find(&bots)

	for _, bot := range bots {
		assert.NotEmpty(s.T(), bot["tenant_id"], "tenant_id不应该为空")
		assert.Equal(s.T(), "test_tenant", bot["tenant_id"], "tenant_id应该是test_tenant")
	}
}

// TestMigration_ErrorHandling 测试错误处理
func (s *MigrationTestSuite) TestMigration_ErrorHandling() {
	ctx := context.Background()

	// 1. 测试重复添加字段（应该失败）
	err := s.migrationTool.AddTenantIDColumn(ctx, "test_bots")
	require.NoError(s.T(), err)

	err = s.migrationTool.AddTenantIDColumn(ctx, "test_bots")
	assert.Error(s.T(), err, "重复添加字段应该失败")

	// 2. 测试不存在的表
	err = s.migrationTool.AddTenantIDColumn(ctx, "non_existent_table")
	assert.Error(s.T(), err, "不存在的表应该返回错误")

	// 3. 测试空数据库（不应该有影响）
	s.cleanupTestData()
	result, err := s.migrationTool.BackfillTenantID(ctx, "test_bots", "test_tenant", 100)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(0), result.MigratedCount, "空数据库迁移记录数应该是0")
}

// ============================================================================
// 辅助方法
// ============================================================================

// setupTestDatabase 创建测试数据库
func setupTestDatabase(dbName string) error {
	// 连接到MySQL（不指定数据库）
	dsn := formatDSN("", "")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	// 创建数据库
	db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)).Error
	if err != nil {
		return fmt.Errorf("创建数据库失败: %w", err)
	}

	log.Printf("测试数据库创建成功: %s", dbName)
	return nil
}

// cleanupTestDatabase 清理测试数据库
func cleanupTestDatabase() {
	if config.TestDBName == "" {
		return
	}

	dsn := formatDSN("", "")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("清理测试数据库失败: %v", err)
		return
	}

	db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.TestDBName))
	log.Printf("测试数据库已清理: %s", config.TestDBName)
}

// formatDSN 格式化数据源名称
func formatDSN(user, dbName string) string {
	if user == "" {
		user = config.DBUser
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, config.DBPassword, config.DBHost, config.DBPort, dbName)
}

// createTestTable 创建测试表
func (s *MigrationTestSuite) createTestTable(tableName string) {
	s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(200) NOT NULL,
			creator_id VARCHAR(100) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_creator (creator_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`, tableName))
}

// prepareTestData 准备测试数据
func (s *MigrationTestSuite) prepareTestData(tableName string, count int) {
	s.createTestTable(tableName)

	for i := 0; i < count; i++ {
		s.db.Table(tableName).Create(map[string]interface{}{
			"name":       fmt.Sprintf("test_%d", i),
			"creator_id": "test_user",
			"status":     "active",
		})
	}
}

// prepareLargeScaleData 准备大规模测试数据
func (s *MigrationTestSuite) prepareLargeScaleData(tableName string, count int) {
	s.createTestTable(tableName)

	// 批量插入
	batchSize := 1000
	for i := 0; i < count; i += batchSize {
		end := i + batchSize
		if end > count {
			end = count
		}

		records := make([]map[string]interface{}, 0, batchSize)
		for j := i; j < end; j++ {
			records = append(records, map[string]interface{}{
				"name":       fmt.Sprintf("test_%d", j),
				"creator_id": "test_user",
				"status":     "active",
			})
		}

		s.db.Table(tableName).Create(&records)
	}
}

// cleanupTestData 清理测试数据
func (s *MigrationTestSuite) cleanupTestData() {
	tables := []string{"test_bots", "test_conversations", "test_knowledge"}
	for _, tableName := range tables {
		s.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))
	}
}

// ============================================================================
// 运行测试套件
// ============================================================================

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}

// ============================================================================
// 辅助类型定义
// ============================================================================

// MigrationResult 迁移结果
type MigrationResult struct {
	MigratedCount int64
	FailedCount   int64
	StartTime     time.Time
	EndTime       time.Time
	Duration      time.Duration
}

// TenantIDValidation tenant_id验证结果
type TenantIDValidation struct {
	TotalRecords int64
	ValidRecords int64
	NullRecords  int64
	AllValid     bool
}

// 需要添加 sync 导入
import "sync"
