// backend/scripts/migration/monitor.go
// 迁移监控面板 - 实时进度、性能指标、数据可视化
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MonitorServer 监控服务器
type MonitorServer struct {
	db      *gorm.DB
	httpSrv *http.Server
	monitor *ProgressMonitor
}

// ProgressMonitor 进度监控器
type ProgressMonitor struct {
	migrations map[string]*MigrationStatus
	mu         chan struct{}
}

// MigrationStatus 迁移状态
type MigrationStatus struct {
	MigrationID     string                 `json:"migration_id"`
	Status          string                 `json:"status"` // pending, running, completed, failed, rolled_back
	CurrentTable    string                 `json:"current_table"`
	ProgressPercent float64                `json:"progress_percent"`
	TotalTables     int                    `json:"total_tables"`
	CompletedTables int                    `json:"completed_tables"`
	TotalRecords    int64                  `json:"total_records"`
	MigratedRecords int64                  `json:"migrated_records"`
	FailedRecords   int64                  `json:"failed_records"`
	StartedAt       time.Time              `json:"started_at"`
	FinishedAt      *time.Time             `json:"finished_at,omitempty"`
	Duration        time.Duration          `json:"duration"`
	QPS             float64                `json:"qps"`
	ETA             time.Duration          `json:"eta"`
	Error           string                 `json:"error,omitempty"`
	Metrics         *MigrationMetrics      `json:"metrics,omitempty"`
	TableStatus     map[string]*TableProgress `json:"table_status,omitempty"`
}

// MigrationMetrics 迁移指标
type MigrationMetrics struct {
	CPUUsage        float64   `json:"cpu_usage"`
	MemoryUsage     float64   `json:"memory_usage"`
	DiskIO          float64   `json:"disk_io"`
	NetworkIO       float64   `json:"network_io"`
	DBConnections   int       `json:"db_connections"`
	ActiveThreads   int       `json:"active_threads"`
	LockWaits       int64     `json:"lock_waits"`
	Deadlocks       int64     `json:"deadlocks"`
	SlowQueries     int64     `json:"slow_queries"`
	Timestamp       time.Time `json:"timestamp"`
}

// TableProgress 表进度
type TableProgress struct {
	TableName       string        `json:"table_name"`
	Status          string        `json:"status"`
	TotalRecords    int64         `json:"total_records"`
	MigratedRecords int64         `json:"migrated_records"`
	ProgressPercent float64       `json:"progress_percent"`
	StartedAt       time.Time     `json:"started_at"`
	FinishedAt      *time.Time    `json:"finished_at,omitempty"`
	Duration        time.Duration `json:"duration"`
	QPS             float64       `json:"qps"`
}

// DashboardData 仪表板数据
type DashboardData struct {
	Timestamp        time.Time                      `json:"timestamp"`
	ActiveMigrations int                            `json:"active_migrations"`
	CompletedMigrations int                         `json:"completed_migrations"`
	FailedMigrations int                            `json:"failed_migrations"`
	MigrationStatus  map[string]*MigrationStatus    `json:"migration_status"`
	SystemMetrics    *SystemMetrics                 `json:"system_metrics"`
	RecentLogs       []LogEntry                     `json:"recent_logs"`
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryUsage   float64 `json:"memory_usage"`
	DiskUsage     float64 `json:"disk_usage"`
	NetworkRX     float64 `json:"network_rx"`
	NetworkTX     float64 `json:"network_tx"`
	DBConnections int     `json:"db_connections"`
	Timestamp     time.Time `json:"timestamp"`
}

// LogEntry 日志条目
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
}

var activeMonitor *MonitorServer

func main() {
	config := parseMonitorFlags()

	var err error
	activeMonitor, err = NewMonitorServer(config)
	if err != nil {
		log.Fatalf("创建监控服务器失败: %v", err)
	}

	log.Printf("启动迁移监控面板 http://localhost:%d\n", config.Port)
	log.Fatal(activeMonitor.Start())
}

// Config 配置
type Config struct {
	DBHost    string
	DBPort    int
	DBUser    string
	DBPassword string
	DBName    string
	Port      int
}

// parseMonitorFlags 解析命令行参数
func parseMonitorFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.DBHost, "host", "localhost", "MySQL主机地址")
	flag.IntVar(&config.DBPort, "port", 3306, "MySQL端口")
	flag.StringVar(&config.DBUser, "user", "root", "MySQL用户名")
	flag.StringVar(&config.DBPassword, "password", "", "MySQL密码")
	flag.StringVar(&config.DBName, "database", "coze_studio", "数据库名称")
	flag.IntVar(&config.Port, "web-port", 8080, "Web面板端口")

	flag.Parse()

	return config
}

// NewMonitorServer 创建监控服务器
func NewMonitorServer(config *Config) (*MonitorServer, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	monitor := &MonitorServer{
		db: db,
		httpSrv: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.Port),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		monitor: &ProgressMonitor{
			migrations: make(map[string]*MigrationStatus),
			mu:         make(chan struct{}),
		},
	}

	// 注册路由
	monitor.registerRoutes()

	return monitor, nil
}

// registerRoutes 注册路由
func (m *MonitorServer) registerRoutes() {
	http.HandleFunc("/", m.handleDashboard)
	http.HandleFunc("/api/migrations", m.handleListMigrations)
	http.HandleFunc("/api/migration/", m.handleGetMigration)
	http.HandleFunc("/api/metrics", m.handleMetrics)
	http.HandleFunc("/api/logs", m.handleLogs)
	http.HandleFunc("/health", m.handleHealth)
}

// Start 启动服务器
func (m *MonitorServer) Start() error {
	return m.httpSrv.ListenAndServe()
}

// handleDashboard 处理仪表板页面
func (m *MonitorServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>ZKER 迁移监控面板</title>
    <meta charset="utf-8">
    <meta http-equiv="refresh" content="5">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #007bff; padding-bottom: 10px; }
        .metrics { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin: 20px 0; }
        .metric-card { background: #f8f9fa; padding: 15px; border-radius: 5px; border-left: 4px solid #007bff; }
        .metric-card h3 { margin: 0 0 10px 0; font-size: 14px; color: #666; }
        .metric-card .value { font-size: 24px; font-weight: bold; color: #333; }
        .status-running { border-left-color: #28a745; }
        .status-completed { border-left-color: #007bff; }
        .status-failed { border-left-color: #dc3545; }
        .migration-list { margin-top: 20px; }
        .migration-item { padding: 15px; margin-bottom: 10px; background: #f8f9fa; border-radius: 5px; }
        .progress-bar { width: 100%; height: 20px; background: #e9ecef; border-radius: 10px; overflow: hidden; margin-top: 10px; }
        .progress-fill { height: 100%; background: linear-gradient(90deg, #007bff, #28a745); transition: width 0.3s; }
        table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #dee2e6; }
        th { background: #f8f9fa; font-weight: 600; }
        .badge { padding: 4px 8px; border-radius: 4px; font-size: 12px; font-weight: 600; }
        .badge-running { background: #d4edda; color: #155724; }
        .badge-completed { background: #cce5ff; color: #004085; }
        .badge-failed { background: #f8d7da; color: #721c24; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 ZKER 数据迁移监控面板</h1>
        <p>最后更新: <span id="timestamp"></span></p>

        <div class="metrics">
            <div class="metric-card">
                <h3>活跃迁移</h3>
                <div class="value" id="active-migrations">-</div>
            </div>
            <div class="metric-card">
                <h3>已完成</h3>
                <div class="value" id="completed-migrations">-</div>
            </div>
            <div class="metric-card">
                <h3>失败</h3>
                <div class="value" id="failed-migrations">-</div>
            </div>
            <div class="metric-card">
                <h3>CPU使用率</h3>
                <div class="value" id="cpu-usage">-</div>
            </div>
            <div class="metric-card">
                <h3>内存使用率</h3>
                <div class="value" id="memory-usage">-</div>
            </div>
            <div class="metric-card">
                <h3>数据库连接</h3>
                <div class="value" id="db-connections">-</div>
            </div>
        </div>

        <div class="migration-list">
            <h2>迁移任务列表</h2>
            <div id="migration-list"></div>
        </div>
    </div>

    <script>
        // 更新时间戳
        document.getElementById('timestamp').textContent = new Date().toLocaleString('zh-CN');

        // 加载数据
        fetch('/api/dashboard')
            .then(res => res.json())
            .then(data => {
                document.getElementById('active-migrations').textContent = data.active_migrations;
                document.getElementById('completed-migrations').textContent = data.completed_migrations;
                document.getElementById('failed-migrations').textContent = data.failed_migrations;

                if (data.system_metrics) {
                    document.getElementById('cpu-usage').textContent = data.system_metrics.cpu_usage.toFixed(1) + '%';
                    document.getElementById('memory-usage').textContent = data.system_metrics.memory_usage.toFixed(1) + '%';
                    document.getElementById('db-connections').textContent = data.system_metrics.db_connections;
                }

                const listEl = document.getElementById('migration-list');
                if (Object.keys(data.migration_status).length === 0) {
                    listEl.innerHTML = '<p style="text-align: center; color: #666; padding: 20px;">暂无活跃迁移任务</p>';
                } else {
                    let html = '';
                    for (const [id, status] of Object.entries(data.migration_status)) {
                        const statusClass = status.status === 'running' ? 'running' :
                                          status.status === 'completed' ? 'completed' : 'failed';
                        html += '<div class="migration-item status-' + statusClass + '">';
                        html += '<h3>' + status.migration_id + '</h3>';
                        html += '<p><strong>状态:</strong> <span class="badge badge-' + statusClass + '">' + status.status + '</span></p>';
                        html += '<p><strong>当前表:</strong> ' + status.current_table + '</p>';
                        html += '<p><strong>进度:</strong> ' + status.progress_percent.toFixed(2) + '% (' +
                                status.migrated_records + '/' + status.total_records + ')</p>';
                        html += '<div class="progress-bar"><div class="progress-fill" style="width: ' +
                                status.progress_percent + '%"></div></div>';
                        html += '<p><strong>QPS:</strong> ' + status.qps.toFixed(2) + ' | <strong>ETA:</strong> ' +
                                formatDuration(status.eta) + '</p>';

                        if (status.table_status) {
                            html += '<table><thead><tr><th>表名</th><th>进度</th><th>QPS</th><th>状态</th></tr></thead><tbody>';
                            for (const [tableName, tableStatus] of Object.entries(status.table_status)) {
                                html += '<tr>';
                                html += '<td>' + tableName + '</td>';
                                html += '<td>' + tableStatus.progress_percent.toFixed(2) + '%</td>';
                                html += '<td>' + tableStatus.qps.toFixed(2) + '</td>';
                                html += '<td>' + tableStatus.status + '</td>';
                                html += '</tr>';
                            }
                            html += '</tbody></table>';
                        }

                        html += '</div>';
                    }
                    listEl.innerHTML = html;
                }
            })
            .catch(err => console.error('加载数据失败:', err));

        function formatDuration(ns) {
            const seconds = ns / 1e9;
            if (seconds < 60) return seconds.toFixed(0) + '秒';
            if (seconds < 3600) return (seconds / 60).toFixed(0) + '分钟';
            return (seconds / 3600).toFixed(1) + '小时';
        }
    </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// handleListMigrations 处理列出所有迁移
func (m *MonitorServer) handleListMigrations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 从数据库查询迁移记录
	var migrations []struct {
		MigrationID     string    `gorm:"column:migration_id"`
		Status          string    `gorm:"column:status"`
		ProgressPercent float64   `gorm:"column:progress_percent"`
		TotalRecords    int64     `gorm:"column:total_records"`
		MigratedRecords int64     `gorm:"column:migrated_records"`
		FailedRecords   int64     `gorm:"column:failed_records"`
		StartedAt       time.Time `gorm:"column:started_at"`
		FinishedAt      *time.Time `gorm:"column:finished_at"`
		Error           string    `gorm:"column:error"`
	}

	err := m.db.Table("migration_records").
		Order("started_at DESC").
		Find(&migrations).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(migrations)
}

// handleGetMigration 处理获取单个迁移详情
func (m *MonitorServer) handleGetMigration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 从URL路径提取migration_id
	// /api/migration/{id}
	path := r.URL.Path[len("/api/migration/"):]
	if path == "" {
		http.Error(w, "missing migration id", http.StatusBadRequest)
		return
	}

	// 查询迁移记录
	var record struct {
		MigrationID     string    `gorm:"column:migration_id"`
		Status          string    `gorm:"column:status"`
		ProgressPercent float64   `gorm:"column:progress_percent"`
		TotalRecords    int64     `gorm:"column:total_records"`
		MigratedRecords int64     `gorm:"column:migrated_records"`
		FailedRecords   int64     `gorm:"column:failed_records"`
		StartedAt       time.Time `gorm:"column:started_at"`
		FinishedAt      *time.Time `gorm:"column:finished_at"`
		Error           string    `gorm:"column:error"`
	}

	err := m.db.Table("migration_records").
		Where("migration_id = ?", path).
		First(&record).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(record)
}

// handleMetrics 处理系统指标
func (m *MonitorServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	metrics := m.collectSystemMetrics()
	json.NewEncoder(w).Encode(metrics)
}

// handleLogs 处理日志
func (m *MonitorServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: 从日志文件或数据库读取
	logs := []LogEntry{}

	json.NewEncoder(w).Encode(logs)
}

// handleHealth 处理健康检查
func (m *MonitorServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 检查数据库连接
	sqlDB, err := m.db.DB()
	if err != nil {
		http.Error(w, "database connection failed", http.StatusInternalServerError)
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		http.Error(w, "database ping failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}

// collectSystemMetrics 收集系统指标
func (m *MonitorServer) collectSystemMetrics() *SystemMetrics {
	// TODO: 实际应该从系统监控库获取
	// 这里简化为返回示例数据

	return &SystemMetrics{
		CPUUsage:      45.2,
		MemoryUsage:   62.8,
		DiskUsage:     78.5,
		NetworkRX:     1024.5,
		NetworkTX:     2048.3,
		DBConnections: 15,
		Timestamp:     time.Now(),
	}
}
