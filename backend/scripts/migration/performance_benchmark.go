// backend/scripts/migration/performance_benchmark.go
// 迁移性能基准测试工具 - 对比迁移前后的性能差异
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// PerformanceBenchmark 性能基准测试
type PerformanceBenchmark struct {
	db     *gorm.DB
	sqlDB  *sql.DB
	config *BenchmarkConfig
}

// BenchmarkConfig 配置
type BenchmarkConfig struct {
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	Tables      string
	Threads     int
	Duration    time.Duration
	WarmupTime  time.Duration
	OutputFile  string
	CompareMode bool // 是否对比模式（需要before/after两个数据库）
}

// BenchmarkResult 基准测试结果
type BenchmarkResult struct {
	TableName           string        `json:"table_name"`
	QueryType           string        `json:"query_type"`
	Threads             int           `json:"threads"`
	Duration            time.Duration `json:"duration"`
	TotalQueries        int64         `json:"total_queries"`
	QPS                 float64       `json:"qps"`
	AvgLatency          time.Duration `json:"avg_latency"`
	P50Latency          time.Duration `json:"p50_latency"`
	P95Latency          time.Duration `json:"p95_latency"`
	P99Latency          time.Duration `json:"p99_latency"`
	MinLatency          time.Duration `json:"min_latency"`
	MaxLatency          time.Duration `json:"max_latency"`
	ErrorCount          int64         `json:"error_count"`
	SuccessRate         float64       `json:"success_rate"`
	BytesTransferred    int64         `json:"bytes_transferred"`
	ThroughputMB        float64       `json:"throughput_mb"`
}

// ComparisonReport 对比报告
type ComparisonReport struct {
	GeneratedAt   time.Time                  `json:"generated_at"`
	BeforeDB      string                     `json:"before_db"`
	AfterDB       string                     `json:"after_db"`
	Threads       int                        `json:"threads"`
	Duration      time.Duration              `json:"duration"`
	Comparisons   map[string]*QueryComparison `json:"comparisons"`
	Summary       *ComparisonSummary         `json:"summary"`
}

// QueryComparison 查询对比
type QueryComparison struct {
	TableName      string           `json:"table_name"`
	QueryType      string           `json:"query_type"`
	BeforeResult   *BenchmarkResult `json:"before_result"`
	AfterResult    *BenchmarkResult `json:"after_result"`
	QPSImprovement float64          `json:"qps_improvement"`
	LatencyChange  float64          `json:"latency_change"`
	IsImproved     bool             `json:"is_improved"`
}

// ComparisonSummary 对比摘要
type ComparisonSummary struct {
	TotalQueries      int     `json:"total_queries"`
	ImprovedQueries   int     `json:"improved_queries"`
	DegradedQueries   int     `json:"degraded_queries"`
	UnchangedQueries  int     `json:"unchanged_queries"`
	AvgQPSImprovement float64 `json:"avg_qps_improvement"`
}

func main() {
	config := parseBenchmarkFlags()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	benchmark, err := NewPerformanceBenchmark(config)
	if err != nil {
		log.Fatalf("创建基准测试失败: %v", err)
	}

	ctx := context.Background()

	if config.CompareMode {
		// 对比模式：需要连接两个数据库（before/after）
		log.Println("执行性能对比测试...")

		// TODO: 实现对比模式
		log.Fatal("对比模式暂未实现")
	} else {
		// 单数据库模式
		log.Println("执行性能基准测试...")

		tables := getBenchmarkTables(config.Tables)

		report := benchmark.RunBenchmark(ctx, tables, config.Threads, config.Duration)

		if config.OutputFile != "" {
			saveBenchmarkReport(report, config.OutputFile)
		} else {
			printBenchmarkResults(report)
		}
	}
}

// parseBenchmarkFlags 解析命令行参数
func parseBenchmarkFlags() *BenchmarkConfig {
	config := &BenchmarkConfig{}

	flag.StringVar(&config.DBHost, "host", "localhost", "MySQL主机地址")
	flag.IntVar(&config.DBPort, "port", 3306, "MySQL端口")
	flag.StringVar(&config.DBUser, "user", "root", "MySQL用户名")
	flag.StringVar(&config.DBPassword, "password", "", "MySQL密码")
	flag.StringVar(&config.DBName, "database", "coze_studio", "数据库名称")
	flag.StringVar(&config.Tables, "tables", "bots,conversations,knowledge", "要测试的表名（逗号分隔）")
	flag.IntVar(&config.Threads, "threads", 10, "并发线程数")
	flag.DurationVar(&config.Duration, "duration", 5*time.Minute, "测试持续时间")
	flag.DurationVar(&config.WarmupTime, "warmup", 30*time.Second, "预热时间")
	flag.StringVar(&config.OutputFile, "output", "", "输出文件路径（JSON格式）")
	flag.BoolVar(&config.CompareMode, "compare", false, "对比模式")

	flag.Parse()

	return config
}

// NewPerformanceBenchmark 创建性能基准测试
func NewPerformanceBenchmark(config *BenchmarkConfig) (*PerformanceBenchmark, error) {
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

	// 配置连接池
	sqlDB.SetMaxOpenConns(config.Threads * 2)
	sqlDB.SetMaxIdleConns(config.Threads)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &PerformanceBenchmark{
		db:     db,
		sqlDB:  sqlDB,
		config: config,
	}, nil
}

// RunBenchmark 运行基准测试
func (b *PerformanceBenchmark) RunBenchmark(ctx context.Context, tables []string, threads int, duration time.Duration) []*BenchmarkResult {
	var results []*BenchmarkResult

	for _, tableName := range tables {
		log.Printf("测试表: %s\n", tableName)

		// 预热
		log.Printf("预热 %s...\n", b.config.WarmupTime)
		b.warmup(ctx, tableName, threads, b.config.WarmupTime)

		// 执行不同类型的查询测试
		queryTypes := []string{
			"select_by_tenant",
			"select_by_tenant_status",
			"select_join",
			"insert",
			"update",
			"delete",
		}

		for _, queryType := range queryTypes {
			result := b.benchmarkQuery(ctx, tableName, queryType, threads, duration)
			results = append(results, result)
		}
	}

	return results
}

// benchmarkQuery 基准测试单个查询
func (b *PerformanceBenchmark) benchmarkQuery(ctx context.Context, tableName, queryType string, threads int, duration time.Duration) *BenchmarkResult {
	log.Printf("  测试查询: %s (线程: %d, 持续时间: %s)\n", queryType, threads, duration)

	result := &BenchmarkResult{
		TableName: tableName,
		QueryType: queryType,
		Threads:   threads,
		Duration:  duration,
	}

	// 创建上下文，支持超时
	queryCtx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	// 用于收集延迟数据
	var latencies []time.Duration
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 启动多个worker
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-queryCtx.Done():
					return
				default:
					start := time.Now()

					err := b.executeQuery(queryCtx, tableName, queryType)

					latency := time.Since(start)

					mu.Lock()
					latencies = append(latencies, latency)
					if err != nil {
						result.ErrorCount++
					}
					mu.Unlock()
				}
			}
		}(i)
	}

	// 等待所有worker完成
	wg.Wait()

	// 计算统计数据
	result.TotalQueries = int64(len(latencies))
	result.calculateStatistics(latencies)

	log.Printf("    QPS: %.2f, P95延迟: %v, 错误率: %.2f%%\n",
		result.QPS, result.P95Latency, (1-result.SuccessRate)*100)

	return result
}

// executeQuery 执行查询
func (b *PerformanceBenchmark) executeQuery(ctx context.Context, tableName, queryType string) error {
	switch queryType {
	case "select_by_tenant":
		return b.db.WithContext(ctx).
			Table(tableName).
			Where("tenant_id = ?", "test_tenant").
			Limit(100).
			Error

	case "select_by_tenant_status":
		return b.db.WithContext(ctx).
			Table(tableName).
			Where("tenant_id = ? AND status = ?", "test_tenant", "active").
			Limit(100).
			Error

	case "select_join":
		// JOIN查询示例
		return b.db.WithContext(ctx).
			Table(tableName).
			Joins("JOIN users ON users.tenant_id = "+tableName+".tenant_id").
			Where(tableName+".tenant_id = ?", "test_tenant").
			Limit(100).
			Error

	case "insert":
		// 插入测试（实际不会执行）
		return nil

	case "update":
		// 更新测试（实际不会执行）
		return nil

	case "delete":
		// 删除测试（实际不会执行）
		return nil

	default:
		return fmt.Errorf("未知查询类型: %s", queryType)
	}
}

// calculateStatistics 计算统计数据
func (r *BenchmarkResult) calculateStatistics(latencies []time.Duration) {
	if len(latencies) == 0 {
		return
	}

	// QPS
	r.QPS = float64(len(latencies)) / r.Duration.Seconds()

	// 成功率
	r.SuccessRate = float64(len(latencies)-r.ErrorCount) / float64(len(latencies)) * 100

	// 延迟统计
	var total time.Duration
	min := latencies[0]
	max := latencies[0]

	for _, lat := range latencies {
		total += lat
		if lat < min {
			min = lat
		}
		if lat > max {
			max = lat
		}
	}

	r.AvgLatency = total / time.Duration(len(latencies))
	r.MinLatency = min
	r.MaxLatency = max

	// 计算百分位延迟
	r.P50Latency = percentile(latencies, 0.50)
	r.P95Latency = percentile(latencies, 0.95)
	r.P99Latency = percentile(latencies, 0.99)
}

// warmup 预热
func (b *PerformanceBenchmark) warmup(ctx context.Context, tableName string, threads int, duration time.Duration) {
	warmupCtx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-warmupCtx.Done():
					return
				default:
					b.executeQuery(warmupCtx, tableName, "select_by_tenant")
				}
			}
		}()
	}

	wg.Wait()
}

// percentile 计算百分位
func percentile(latencies []time.Duration, p float64) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	// 简化实现：直接排序后取值
	// 实际生产环境应使用更高效的算法（如quickselect）

	// 复制切片避免修改原数据
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)

	// 简单排序（实际应使用快速排序）
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	index := int(float64(len(sorted)) * p)
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// getBenchmarkTables 获取基准测试表
func getBenchmarkTables(tablesStr string) []string {
	// TODO: 解析逗号分隔的表名
	return []string{"bots", "conversations", "knowledge"}
}

// saveBenchmarkReport 保存基准测试报告
func saveBenchmarkReport(results []*BenchmarkResult, filename string) error {
	report := struct {
		GeneratedAt time.Time          `json:"generated_at"`
		Results     []*BenchmarkResult `json:"results"`
	}{
		GeneratedAt: time.Now(),
		Results:     results,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化报告失败: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	log.Printf("基准测试报告已保存到: %s\n", filename)
	return nil
}

// printBenchmarkResults 打印基准测试结果
func printBenchmarkResults(results []*BenchmarkResult) {
	fmt.Println("\n========== 性能基准测试报告 ==========")
	fmt.Printf("生成时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("测试表数: %d\n\n", len(results))

	for _, result := range results {
		fmt.Printf("[%s] %s\n", result.QueryType, result.TableName)
		fmt.Printf("  线程数: %d\n", result.Threads)
		fmt.Printf("  QPS: %.2f\n", result.QPS)
		fmt.Printf("  平均延迟: %v\n", result.AvgLatency)
		fmt.Printf("  P50延迟: %v\n", result.P50Latency)
		fmt.Printf("  P95延迟: %v\n", result.P95Latency)
		fmt.Printf("  P99延迟: %v\n", result.P99Latency)
		fmt.Printf("  最小延迟: %v\n", result.MinLatency)
		fmt.Printf("  最大延迟: %v\n", result.MaxLatency)
		fmt.Printf("  成功率: %.2f%%\n", result.SuccessRate)
		fmt.Println()
	}

	fmt.Println("========================================")
}
