// +build ignore

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// 测试配置
type TestConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// 测试结果
type TestResult struct {
	Name    string
	Passed  bool
	Message string
}

func main() {
	config := parseTestConfig()
	results := runTests(config)

	printResults(results)

	// 如果有测试失败，退出码为1
	for _, result := range results {
		if !result.Passed {
			os.Exit(1)
		}
	}
}

func parseTestConfig() TestConfig {
	config := TestConfig{}
	flag.StringVar(&config.Host, "host", "localhost", "MySQL host")
	flag.IntVar(&config.Port, "port", 3306, "MySQL port")
	flag.StringVar(&config.User, "user", "root", "MySQL user")
	flag.StringVar(&config.Password, "password", "", "MySQL password")
	flag.StringVar(&config.Database, "database", "opencoze", "Database name")
	flag.Parse()

	if config.Password == "" {
		config.Password = os.Getenv("MYSQL_PASSWORD")
	}

	return config
}

func runTests(config TestConfig) []TestResult {
	var results []TestResult

	// 连接数据库
	db, err := connectDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 运行测试
	results = append(results, testTableExists(db, "tenants"))
	results = append(results, testTableExists(db, "subscriptions"))
	results = append(results, testTableExists(db, "quotas"))
	results = append(results, testTableExists(db, "quota_usage_log"))

	results = append(results, testTriggers(db))
	results = append(results, testViews(db))
	results = append(results, testSystemTenant(db))
	results = append(results, testIndexes(db))
	results = append(results, testForeignKeys(db))
	results = append(results, testDataIntegrity(db))

	return results
}

func connectDB(config TestConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func testTableExists(db *sql.DB, tableName string) TestResult {
	query := `
		SELECT COUNT(*)
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
	`

	var count int
	err := db.QueryRow(query, "opencoze", tableName).Scan(&count)

	if err != nil {
		return TestResult{
			Name:    fmt.Sprintf("Table: %s exists", tableName),
			Passed:  false,
			Message: fmt.Sprintf("Error: %v", err),
		}
	}

	if count == 0 {
		return TestResult{
			Name:    fmt.Sprintf("Table: %s exists", tableName),
			Passed:  false,
			Message: "Table does not exist",
		}
	}

	return TestResult{
		Name:    fmt.Sprintf("Table: %s exists", tableName),
		Passed:  true,
		Message: "OK",
	}
}

func testTriggers(db *sql.DB) TestResult {
	query := `
		SELECT COUNT(*)
		FROM information_schema.TRIGGERS
		WHERE TRIGGER_SCHEMA = ? AND TRIGGER_NAME LIKE 'trg_%'
	`

	var count int
	err := db.QueryRow(query, "opencoze").Scan(&count)

	if err != nil {
		return TestResult{
			Name:    "Triggers exist",
			Passed:  false,
			Message: fmt.Sprintf("Error: %v", err),
		}
	}

	if count < 2 {
		return TestResult{
			Name:    "Triggers exist",
			Passed:  false,
			Message: fmt.Sprintf("Expected 2 triggers, found %d", count),
		}
	}

	return TestResult{
		Name:    "Triggers exist",
		Passed:  true,
		Message: fmt.Sprintf("OK (%d triggers)", count),
	}
}

func testViews(db *sql.DB) TestResult {
	query := `
		SELECT COUNT(*)
		FROM information_schema.VIEWS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME LIKE 'v_%'
	`

	var count int
	err := db.QueryRow(query, "opencoze").Scan(&count)

	if err != nil {
		return TestResult{
			Name:    "Views exist",
			Passed:  false,
			Message: fmt.Sprintf("Error: %v", err),
		}
	}

	if count < 3 {
		return TestResult{
			Name:    "Views exist",
			Passed:  false,
			Message: fmt.Sprintf("Expected 3 views, found %d", count),
		}
	}

	return TestResult{
		Name:    "Views exist",
		Passed:  true,
		Message: fmt.Sprintf("OK (%d views)", count),
	}
}

func testSystemTenant(db *sql.DB) TestResult {
	query := `
		SELECT COUNT(*)
		FROM tenants
		WHERE tenant_id = 'system-default'
	`

	var count int
	err := db.QueryRow(query).Scan(&count)

	if err != nil {
		return TestResult{
			Name:    "System tenant exists",
			Passed:  false,
			Message: fmt.Sprintf("Error: %v", err),
		}
	}

	if count == 0 {
		return TestResult{
			Name:    "System tenant exists",
			Passed:  false,
			Message: "System tenant not found",
		}
	}

	return TestResult{
		Name:    "System tenant exists",
		Passed:  true,
		Message: "OK",
	}
}

func testIndexes(db *sql.DB) TestResult {
	expectedIndexes := []string{
		"uk_tenant_name",
		"idx_tenant_type",
		"idx_status",
		"uk_tenant_plan",
		"uk_tenant_resource",
	}

	query := `
		SELECT COUNT(DISTINCT INDEX_NAME)
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ? AND INDEX_NAME IN (?` + string(makeRune(len(expectedIndexes)-1, ",?")) + `)
	`

	args := []interface{}{"opencoze"}
	for _, idx := range expectedIndexes {
		args = append(args, idx)
	}

	var count int
	err := db.QueryRow(query, args...).Scan(&count)

	if err != nil {
		return TestResult{
			Name:    "Indexes exist",
			Passed:  false,
			Message: fmt.Sprintf("Error: %v", err),
		}
	}

	if count < len(expectedIndexes) {
		return TestResult{
			Name:    "Indexes exist",
			Passed:  false,
			Message: fmt.Sprintf("Expected %d indexes, found %d", len(expectedIndexes), count),
		}
	}

	return TestResult{
		Name:    "Indexes exist",
		Passed:  true,
		Message: fmt.Sprintf("OK (%d indexes)", count),
	}
}

func testForeignKeys(db *sql.DB) TestResult {
	query := `
		SELECT COUNT(*)
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = ?
		  AND REFERENCED_TABLE_NAME = 'tenants'
		  AND TABLE_NAME IN ('subscriptions', 'quotas', 'quota_usage_log')
	`

	var count int
	err := db.QueryRow(query, "opencoze").Scan(&count)

	if err != nil {
		return TestResult{
			Name:    "Foreign keys exist",
			Passed:  false,
			Message: fmt.Sprintf("Error: %v", err),
		}
	}

	if count < 3 {
		return TestResult{
			Name:    "Foreign keys exist",
			Passed:  false,
			Message: fmt.Sprintf("Expected 3 foreign keys, found %d", count),
		}
	}

	return TestResult{
		Name:    "Foreign keys exist",
		Passed:  true,
		Message: fmt.Sprintf("OK (%d foreign keys)", count),
	}
}

func testDataIntegrity(db *sql.DB) TestResult {
	// 测试系统租户的订阅和配额
	query := `
		SELECT
			(SELECT COUNT(*) FROM subscriptions WHERE tenant_id = 'system-default') as subs,
			(SELECT COUNT(*) FROM quotas WHERE tenant_id = 'system-default') as quotas
	`

	var subs, quotas int
	err := db.QueryRow(query).Scan(&subs, &quotas)

	if err != nil {
		return TestResult{
			Name:    "Data integrity",
			Passed:  false,
			Message: fmt.Sprintf("Error: %v", err),
		}
	}

	if subs == 0 || quotas == 0 {
		return TestResult{
			Name:    "Data integrity",
			Passed:  false,
			Message: fmt.Sprintf("System tenant missing data: subs=%d, quotas=%d", subs, quotas),
		}
	}

	return TestResult{
		Name:    "Data integrity",
		Passed:  true,
		Message: fmt.Sprintf("OK (subs=%d, quotas=%d)", subs, quotas),
	}
}

func printResults(results []TestResult) {
	fmt.Println("\n========================================")
	fmt.Println("租户表DDL集成测试结果")
	fmt.Println("========================================\n")

	passed := 0
	failed := 0

	for _, result := range results {
		status := "✓"
		if !result.Passed {
			status = "✗"
			failed++
		} else {
			passed++
		}

		fmt.Printf("%s [%s] %s\n", status, result.Name, result.Message)
	}

	fmt.Println("\n========================================")
	fmt.Printf("总计: %d 通过, %d 失败\n", passed, failed)
	fmt.Println("========================================")
}

func makeRune(count int, char string) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}
