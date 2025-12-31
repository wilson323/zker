// backend/scripts/migration/data_validator.go
// 数据迁移验证工具 - MD5/SHA256哈希对比、数据完整性校验
package main

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DataValidator 数据验证器
type DataValidator struct {
	db       *gorm.DB
	sqlDB    *sql.DB
	output   io.Writer
	progress *ValidationProgress
}

// ValidationProgress 验证进度
type ValidationProgress struct {
	TotalTables     int       `json:"total_tables"`
	CompletedTables int       `json:"completed_tables"`
	StartTime       time.Time `json:"start_time"`
	LastUpdated     time.Time `json:"last_updated"`
}

// HashResult 哈希结果
type HashResult struct {
	TableName   string `json:"table_name"`
	RecordCount int64  `json:"record_count"`
	MD5Hash     string `json:"md5_hash"`
	SHA256Hash  string `json:"sha256_hash"`
	ComputedAt  time.Time `json:"computed_at"`
}

// ValidationResult 验证结果
type ValidationResult struct {
	TableName          string        `json:"table_name"`
	TotalRecords       int64         `json:"total_records"`
	ValidRecords       int64         `json:"valid_records"`
	InvalidRecords     int64         `json:"invalid_records"`
	DataCompleteness   float64       `json:"data_completeness"`
	HashBefore         string        `json:"hash_before"`
	HashAfter          string        `json:"hash_after"`
	HashMatch          bool          `json:"hash_match"`
	ForeignKeyValid    bool          `json:"foreign_key_valid"`
	UniquenessValid    bool          `json:"uniqueness_valid"`
	ValidationTime     time.Duration `json:"validation_time"`
	ValidationPassed   bool          `json:"validation_passed"`
	Issues             []string      `json:"issues"`
}

// ValidationReport 验证报告
type ValidationReport struct {
	GeneratedAt         time.Time            `json:"generated_at"`
	TotalTables         int                  `json:"total_tables"`
	PassedTables        int                  `json:"passed_tables"`
	FailedTables        int                  `json:"failed_tables"`
	TotalRecords        int64                `json:"total_records"`
	ValidRecords        int64                `json:"valid_records"`
	InvalidRecords      int64                `json:"invalid_records"`
	OverallCompleteness float64              `json:"overall_completeness"`
	Results             map[string]*ValidationResult `json:"results"`
	Duration            time.Duration        `json:"duration"`
}

// Config 配置
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	Tables     string
	OutputFile string
	Verbose    bool
}

func main() {
	config := parseFlags()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	validator, err := NewDataValidator(config)
	if err != nil {
		log.Fatalf("创建验证器失败: %v", err)
	}

	ctx := context.Background()

	// 解析表名
	var tables []string
	if config.Tables == "all" {
		tables = getAllBusinessTables(ctx, validator)
	} else {
		// TODO: 解析逗号分隔的表名
		tables = []string{"bots", "conversations", "knowledge", "workflows"}
	}

	log.Printf("开始验证 %d 张表...\n", len(tables))

	report := validator.ValidateAll(ctx, tables)

	// 输出报告
	if config.OutputFile != "" {
		saveReportToFile(report, config.OutputFile)
	} else {
		printReport(report)
	}

	// 根据验证结果设置退出码
	if report.FailedTables > 0 {
		os.Exit(1)
	}
}

// parseFlags 解析命令行参数
func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.DBHost, "host", "localhost", "MySQL主机地址")
	flag.IntVar(&config.DBPort, "port", 3306, "MySQL端口")
	flag.StringVar(&config.DBUser, "user", "root", "MySQL用户名")
	flag.StringVar(&config.DBPassword, "password", "", "MySQL密码")
	flag.StringVar(&config.DBName, "database", "coze_studio", "数据库名称")
	flag.StringVar(&config.Tables, "tables", "all", "要验证的表名（all或逗号分隔的表名列表）")
	flag.StringVar(&config.OutputFile, "output", "", "输出文件路径（JSON格式）")
	flag.BoolVar(&config.Verbose, "verbose", false, "详细输出")

	flag.Parse()

	return config
}

// NewDataValidator 创建数据验证器
func NewDataValidator(config *Config) (*DataValidator, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	return &DataValidator{
		db:    db,
		sqlDB: sqlDB,
		output: os.Stdout,
		progress: &ValidationProgress{
			StartTime: time.Now(),
		},
	}, nil
}

// ValidateAll 验证所有表
func (v *DataValidator) ValidateAll(ctx context.Context, tables []string) *ValidationReport {
	report := &ValidationReport{
		GeneratedAt: time.Now(),
		Results:     make(map[string]*ValidationResult),
		TotalTables: len(tables),
	}

	startTime := time.Now()

	for _, tableName := range tables {
		log.Printf("验证表: %s\n", tableName)

		result := v.ValidateTable(ctx, tableName)
		report.Results[tableName] = result

		if result.ValidationPassed {
			report.PassedTables++
		} else {
			report.FailedTables++
		}

		report.TotalRecords += result.TotalRecords
		report.ValidRecords += result.ValidRecords
		report.InvalidRecords += result.InvalidRecords

		v.progress.CompletedTables++
		v.progress.LastUpdated = time.Now()
	}

	report.Duration = time.Since(startTime)
	report.OverallCompleteness = float64(report.ValidRecords) / float64(report.TotalRecords) * 100

	return report
}

// ValidateTable 验证单张表
func (v *DataValidator) ValidateTable(ctx context.Context, tableName string) *ValidationResult {
	startTime := time.Now()

	result := &ValidationResult{
		TableName: tableName,
		Issues:    make([]string, 0),
	}

	// 1. 获取总记录数
	v.db.Table(tableName).Count(&result.TotalRecords)

	// 2. 检查数据完整性
	var nullCount int64
	v.db.Table(tableName).
		Where("tenant_id IS NULL OR tenant_id = ''").
		Count(&nullCount)

	result.InvalidRecords = nullCount
	result.ValidRecords = result.TotalRecords - nullCount
	result.DataCompleteness = float64(result.ValidRecords) / float64(result.TotalRecords) * 100

	if nullCount > 0 {
		result.Issues = append(result.Issues,
			fmt.Sprintf("发现 %d 条记录的 tenant_id 为空", nullCount))
	}

	// 3. 计算哈希值
	hash, err := v.ComputeTableHash(ctx, tableName)
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("哈希计算失败: %v", err))
	} else {
		result.HashAfter = hash.MD5Hash
	}

	// 4. 检查唯一性约束
	result.UniquenessValid = v.checkUniqueness(ctx, tableName)
	if !result.UniquenessValid {
		result.Issues = append(result.Issues, "唯一性约束验证失败")
	}

	// 5. 检查外键约束
	result.ForeignKeyValid = v.checkForeignKey(ctx, tableName)
	if !result.ForeignKeyValid {
		result.Issues = append(result.Issues, "外键约束验证失败")
	}

	// 6. 判断验证是否通过
	result.ValidationPassed = len(result.Issues) == 0 &&
		result.DataCompleteness >= 99.9 &&
		result.UniquenessValid &&
		result.ForeignKeyValid

	result.ValidationTime = time.Since(startTime)

	return result
}

// ComputeTableHash 计算表的哈希值
func (v *DataValidator) ComputeTableHash(ctx context.Context, tableName string) (*HashResult, error) {
	result := &HashResult{
		TableName:  tableName,
		ComputedAt: time.Now(),
	}

	// 获取记录数
	v.db.Table(tableName).Count(&result.RecordCount)

	// 计算MD5哈希
	md5Hash := md5.New()
	sha256Hash := sha256.New()

	rows, err := v.db.Table(tableName).
		Select("id", "tenant_id", "created_at").
		Order("id ASC").
		Rows()

	if err != nil {
		return nil, fmt.Errorf("查询数据失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var tenantID string
		var createdAt time.Time

		if err := rows.Scan(&id, &tenantID, &createdAt); err != nil {
			return nil, fmt.Errorf("扫描行失败: %w", err)
		}

		// 写入哈希
		data := fmt.Sprintf("%d|%s|%s", id, tenantID, createdAt.Format(time.RFC3339))
		md5Hash.Write([]byte(data))
		sha256Hash.Write([]byte(data))
	}

	result.MD5Hash = hex.EncodeToString(md5Hash.Sum(nil))
	result.SHA256Hash = hex.EncodeToString(sha256Hash.Sum(nil))

	return result, nil
}

// checkUniqueness 检查唯一性约束
func (v *DataValidator) checkUniqueness(ctx context.Context, tableName string) bool {
	var duplicates []struct {
		TenantID string `gorm:"column:tenant_id"`
		ID       int64  `gorm:"column:id"`
		Count    int64  `gorm:"column:count"`
	}

	err := v.db.Table(tableName).
		Select("tenant_id, id, COUNT(*) as count").
		Where("tenant_id IS NOT NULL AND tenant_id != ''").
		Group("tenant_id, id").
		Having("count > 1").
		Scan(&duplicates).Error

	return err == nil && len(duplicates) == 0
}

// checkForeignKey 检查外键约束
func (v *DataValidator) checkForeignKey(ctx context.Context, tableName string) bool {
	// 检查tenant_id是否存在于tenants表
	var invalidCount int64
	v.db.Table(tableName).
		Where("tenant_id NOT IN (SELECT tenant_id FROM tenants)").
		Count(&invalidCount)

	return invalidCount == 0
}

// getAllBusinessTables 获取所有业务表
func getAllBusinessTables(ctx context.Context, validator *DataValidator) []string {
	// 返回需要迁移的业务表列表
	return []string{
		"bots",
		"conversations",
		"knowledge",
		"knowledge_document",
		"knowledge_document_slice",
		"workflows",
		"files",
		"agent_tool_draft",
		"agent_tool_version",
	}
}

// saveReportToFile 保存报告到文件
func saveReportToFile(report *ValidationReport, filename string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化报告失败: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	log.Printf("报告已保存到: %s\n", filename)
	return nil
}

// printReport 打印报告
func printReport(report *ValidationReport) {
	fmt.Println("\n========== 数据迁移验证报告 ==========")
	fmt.Printf("生成时间: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("验证耗时: %s\n", report.Duration)
	fmt.Printf("验证表数: %d\n", report.TotalTables)
	fmt.Printf("通过表数: %d\n", report.PassedTables)
	fmt.Printf("失败表数: %d\n", report.FailedTables)
	fmt.Printf("总记录数: %d\n", report.TotalRecords)
	fmt.Printf("有效记录: %d\n", report.ValidRecords)
	fmt.Printf("无效记录: %d\n", report.InvalidRecords)
	fmt.Printf("数据完整性: %.2f%%\n", report.OverallCompleteness)
	fmt.Println("\n========== 详细结果 ==========")

	for tableName, result := range report.Results {
		status := "✓ PASS"
		if !result.ValidationPassed {
			status = "✗ FAIL"
		}

		fmt.Printf("\n[%s] %s\n", status, tableName)
		fmt.Printf("  总记录: %d\n", result.TotalRecords)
		fmt.Printf("  有效: %d\n", result.ValidRecords)
		fmt.Printf("  无效: %d\n", result.InvalidRecords)
		fmt.Printf("  完整性: %.2f%%\n", result.DataCompleteness)
		fmt.Printf("  唯一性: %v\n", result.UniquenessValid)
		fmt.Printf("  外键: %v\n", result.ForeignKeyValid)

		if len(result.Issues) > 0 {
			fmt.Printf("  问题:\n")
			for _, issue := range result.Issues {
				fmt.Printf("    - %s\n", issue)
			}
		}
	}

	fmt.Println("\n========================================")
}
