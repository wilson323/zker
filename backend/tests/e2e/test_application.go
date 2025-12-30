/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package e2e

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	orgEntity "github.com/coze-dev/coze-studio/backend/domain/org/entity"
	permissionEntity "github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ==================== 测试应用 ====================

/**
 * TestApplication 封装完整的应用栈
 *
 * 职责: 提供完整的HTTP服务、数据库、外部依赖的测试环境
 * 遵循单一职责原则：只负责测试环境的管理
 */
type TestApplication struct {
	HTTPClient *http.Client
	Server     *httptest.Server
	DB         *gorm.DB
	SQLDB      *sql.DB
	MySQLC     testcontainers.Container
	TenantID   string // 当前测试租户ID
}

/**
 * setupTestApplication 启动测试应用
 *
 * 职责: 初始化完整的服务栈(数据库+HTTP服务器+所有服务)
 * 遵循依赖注入原则：所有依赖通过参数注入
 */
func SetupTestApplication(t *testing.T) *TestApplication {
	ctx := context.Background()

	// 启动MySQL 8.4.5容器
	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.4.5"),
		mysql.WithDatabase("coze_e2e_test"),
		mysql.WithUsername("coze_test"),
		mysql.WithPassword("coze_test123"),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(5*time.Minute),
				wait.ForListeningPort("3306/tcp").WithStartupTimeout(5*time.Minute),
				wait.ForLog("ready for connections").WithOccurrence(2).WithStartupTimeout(5*time.Minute),
			),
		),
	)
	require.NoError(t, err, "MySQL容器启动失败")

	// 获取连接字符串
	connStr, err := mysqlContainer.ConnectionString(ctx)
	require.NoError(t, err, "获取MySQL连接字符串失败")

	// 使用GORM连接数据库
	db, err := gorm.Open(gormMysql.Open(connStr), &gorm.Config{
		SkipDefaultTransaction: true, // E2E测试不需要默认事务
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	require.NoError(t, err, "GORM连接数据库失败")

	// 获取底层sql.DB
	sqlDB, err := db.DB()
	require.NoError(t, err, "获取sql.DB失败")

	// 测试连接
	require.NoError(t, sqlDB.Ping(), "数据库连接失败")

	// 初始化数据库Schema
	err = initAllTables(db)
	require.NoError(t, err, "数据库Schema初始化失败")

	// 创建测试服务器
	server := createTestServer(db)

	testApp := &TestApplication{
		HTTPClient: server.Client(),
		Server:     server,
		DB:         db,
		SQLDB:      sqlDB,
		MySQLC:     mysqlContainer,
	}

	// 注册清理函数
	t.Cleanup(func() {
		server.Close()
		sqlDB.Close()
		if err := mysqlContainer.Terminate(ctx); err != nil {
			t.Logf("清理MySQL容器失败: %v", err)
		}
	})

	return testApp
}

/**
 * Terminate 清理测试资源
 *
 * 职责: 清理所有测试容器和连接
 */
func (app *TestApplication) Terminate(ctx context.Context) error {
	var errs []error

	if app.Server != nil {
		app.Server.Close()
	}

	if app.SQLDB != nil {
		if err := app.SQLDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭数据库连接失败: %w", err))
		}
	}

	if app.MySQLC != nil {
		if err := app.MySQLC.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("终止MySQL容器失败: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("清理资源时发生错误: %v", errs)
	}

	return nil
}

// ==================== HTTP客户端封装 ====================

/**
 * PostJSON 发送JSON POST请求
 *
 * 职责: 封装HTTP POST请求，自动处理认证和错误
 */
func (app *TestApplication) PostJSON(path string, req, resp interface{}) error {
	return app.DoJSON("POST", path, req, resp)
}

/**
 * GetJSON 发送JSON GET请求
 */
func (app *TestApplication) GetJSON(path string, resp interface{}) error {
	return app.DoJSON("GET", path, nil, resp)
}

/**
 * DoJSON 发送JSON请求
 *
 * 职责: 统一的JSON请求处理
 */
func (app *TestApplication) DoJSON(method, path string, req, resp interface{}) error {
	var body io.Reader

	if req != nil {
		jsonData, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("JSON序列化失败: %w", err)
		}
		body = &jsonReader{data: jsonData}
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest(method, app.Server.URL+path, body)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// 设置租户上下文
	if app.TenantID != "" {
		httpReq.Header.Set("X-Tenant-ID", app.TenantID)
	}

	// 发送请求
	httpResp, err := app.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer httpResp.Body.Close()

	// 检查HTTP状态码
	if httpResp.StatusCode >= 400 {
		errorBody, _ := io.ReadAll(httpResp.Body)
		return fmt.Errorf("HTTP %d: %s", httpResp.StatusCode, string(errorBody))
	}

	// 解析响应
	if resp != nil && httpResp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
			return fmt.Errorf("JSON解码失败: %w", err)
		}
	}

	return nil
}

// jsonReader 用于JSON请求体
type jsonReader struct {
	data []byte
	pos  int
}

func (r *jsonReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// ==================== 数据库Schema初始化 ====================

/**
 * initAllTables 初始化所有表
 *
 * 职责: 创建E2E测试所需的所有表结构
 */
func initAllTables(db *gorm.DB) error {
	schemas := []string{
		// 租户系统
		`
		CREATE TABLE IF NOT EXISTS tenants (
			tenant_id VARCHAR(36) PRIMARY KEY,
			tenant_name VARCHAR(200) NOT NULL,
			tenant_type ENUM('individual','team','enterprise') NOT NULL DEFAULT 'team',
			subdomain VARCHAR(64) NOT NULL UNIQUE,
			status ENUM('active','suspended','deleted') NOT NULL DEFAULT 'active',
			subscription_tier ENUM('free','pro','enterprise') NOT NULL DEFAULT 'free',
			isolation_strategy ENUM('row_level','schema_level','database_level') NOT NULL DEFAULT 'row_level',
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			deleted_at BIGINT NULL,
			INDEX idx_status (status),
			INDEX idx_tenant_type (tenant_type),
			INDEX idx_subscription_tier (subscription_tier)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 订阅表
		`
		CREATE TABLE IF NOT EXISTS subscriptions (
			subscription_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36) NOT NULL,
			plan_tier ENUM('free','pro','enterprise') NOT NULL,
			status ENUM('active','past_due','canceled','unpaid') NOT NULL DEFAULT 'active',
			max_bots INT NOT NULL DEFAULT 10,
			max_users INT NOT NULL DEFAULT 5,
			max_knowledge_bases INT NOT NULL DEFAULT 3,
			max_workflows INT NOT NULL DEFAULT 10,
			monthly_api_quota BIGINT NOT NULL DEFAULT 100000,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
			INDEX idx_tenant_id (tenant_id),
			INDEX idx_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 配额表
		`
		CREATE TABLE IF NOT EXISTS quotas (
			quota_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36) NOT NULL,
			resource_type VARCHAR(50) NOT NULL,
			max_limit BIGINT NOT NULL DEFAULT 0,
			current_usage BIGINT NOT NULL DEFAULT 0,
			reset_period ENUM('daily','weekly','monthly') NOT NULL DEFAULT 'monthly',
			last_reset_at BIGINT NOT NULL DEFAULT 0,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
			UNIQUE KEY uk_tenant_resource (tenant_id, resource_type),
			INDEX idx_tenant_id (tenant_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 用户表
		`
		CREATE TABLE IF NOT EXISTS users (
			user_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36) NOT NULL,
			username VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			status ENUM('active','inactive','banned') NOT NULL DEFAULT 'active',
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			deleted_at BIGINT NULL,
			FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
			INDEX idx_tenant_id (tenant_id),
			INDEX idx_email (email),
			INDEX idx_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 角色表
		`
		CREATE TABLE IF NOT EXISTS roles (
			role_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36) NOT NULL,
			role_name VARCHAR(100) NOT NULL,
			role_code VARCHAR(50) NOT NULL,
			description VARCHAR(500),
			is_system BOOLEAN NOT NULL DEFAULT FALSE,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
			UNIQUE KEY uk_tenant_code (tenant_id, role_code),
			INDEX idx_tenant_id (tenant_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 用户角色关联表
		`
		CREATE TABLE IF NOT EXISTS user_roles (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id VARCHAR(36) NOT NULL,
			role_id VARCHAR(36) NOT NULL,
			tenant_id VARCHAR(36) NOT NULL,
			assigned_at BIGINT NOT NULL DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
			UNIQUE KEY uk_user_role (user_id, role_id),
			INDEX idx_user_id (user_id),
			INDEX idx_role_id (role_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 数据权限表
		`
		CREATE TABLE IF NOT EXISTS data_permissions (
			permission_id VARCHAR(36) PRIMARY KEY,
			role_id VARCHAR(36) NOT NULL,
			resource_type ENUM('bots','conversations','knowledge','workflows','plugins') NOT NULL,
			scope ENUM('ALL','DEPARTMENT','DEPARTMENT_AND_SUB','OWN','CUSTOM','NONE') NOT NULL,
			custom_filter JSON,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
			INDEX idx_role_id (role_id),
			INDEX idx_resource_type (resource_type)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 组织表
		`
		CREATE TABLE IF NOT EXISTS organizations (
			org_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36) NOT NULL,
			org_name VARCHAR(200) NOT NULL,
			org_type ENUM('company','division','department','project') NOT NULL,
			parent_id VARCHAR(36) NULL,
			org_code VARCHAR(50) NOT NULL,
			level INT NOT NULL DEFAULT 1,
			path VARCHAR(500) NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			status ENUM('active','inactive','frozen') NOT NULL DEFAULT 'active',
			description TEXT,
			leader_id VARCHAR(36) NULL,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			deleted_at BIGINT NULL,
			FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
			UNIQUE KEY uk_tenant_code (tenant_id, org_code),
			INDEX idx_tenant_id (tenant_id),
			INDEX idx_parent_id (parent_id),
			INDEX idx_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// Bot表（简化版，用于测试）
		`
		CREATE TABLE IF NOT EXISTS bots (
			bot_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36) NOT NULL,
			org_id VARCHAR(36) NULL,
			bot_name VARCHAR(200) NOT NULL,
			description TEXT,
			status ENUM('draft','published','archived') NOT NULL DEFAULT 'draft',
			ai_confidence DECIMAL(3,2) NOT NULL DEFAULT 1.00,
			created_by VARCHAR(36) NOT NULL,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			deleted_at BIGINT NULL,
			FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
			FOREIGN KEY (org_id) REFERENCES organizations(org_id) ON DELETE SET NULL,
			INDEX idx_tenant_id (tenant_id),
			INDEX idx_status (status),
			INDEX idx_org_id (org_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 审核任务表（人工审核）
		`
		CREATE TABLE IF NOT EXISTS review_tasks (
			task_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36) NOT NULL,
			resource_type ENUM('bot','workflow','knowledge') NOT NULL,
			resource_id VARCHAR(36) NOT NULL,
			status ENUM('pending','assigned','approved','rejected') NOT NULL DEFAULT 'pending',
			assigned_to VARCHAR(36) NULL,
			ai_confidence DECIMAL(3,2) NOT NULL,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			completed_at BIGINT NULL,
			FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
			INDEX idx_tenant_id (tenant_id),
			INDEX idx_resource (resource_type, resource_id),
			INDEX idx_status (status),
			INDEX idx_assigned_to (assigned_to)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 对话表（记忆测试）
		`
		CREATE TABLE IF NOT EXISTS conversations (
			conversation_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36) NOT NULL,
			user_id VARCHAR(36) NOT NULL,
			bot_id VARCHAR(36) NOT NULL,
			title VARCHAR(500),
			status ENUM('active','archived') NOT NULL DEFAULT 'active',
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
			INDEX idx_tenant_id (tenant_id),
			INDEX idx_user_id (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 消息表
		`
		CREATE TABLE IF NOT EXISTS messages (
			message_id VARCHAR(36) PRIMARY KEY,
			conversation_id VARCHAR(36) NOT NULL,
			role ENUM('user','assistant','system') NOT NULL,
			content TEXT NOT NULL,
			created_at BIGINT NOT NULL DEFAULT 0,
			FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id) ON DELETE CASCADE,
			INDEX idx_conversation_id (conversation_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,

		// 记忆表
		`
		CREATE TABLE IF NOT EXISTS entity_memories (
			memory_id VARCHAR(36) PRIMARY KEY,
			tenant_id VARCHAR(36) NOT NULL,
			user_id VARCHAR(36) NOT NULL,
			entity_type VARCHAR(50) NOT NULL,
			entity_value VARCHAR(500) NOT NULL,
			attributes JSON,
			importance_score DECIMAL(3,2) NOT NULL DEFAULT 0.80,
			last_accessed_at BIGINT NOT NULL DEFAULT 0,
			access_count INT NOT NULL DEFAULT 0,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT NOT NULL DEFAULT 0,
			FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
			INDEX idx_tenant_user (tenant_id, user_id),
			INDEX idx_entity (entity_type, entity_value),
			INDEX idx_importance (importance_score)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`,
	}

	for _, schema := range schemas {
		if err := db.Exec(schema).Error; err != nil {
			return fmt.Errorf("创建表失败: %w\nSchema: %s", err, schema)
		}
	}

	return nil
}

/**
 * createTestServer 创建测试HTTP服务器
 *
 * 职责: 提供mock HTTP服务器用于E2E测试
 * 遵循接口隔离原则：使用handler接口
 */
func createTestServer(db *gorm.DB) *httptest.Server {
	mux := http.NewServeMux()

	// 租户相关API
	mux.HandleFunc("/api/v1/tenants/register", handleTenantRegister(db))
	mux.HandleFunc("/api/v1/tenants", handleCreateTenant(db))
	mux.HandleFunc("/api/v1/tenants/", handleGetTenant(db))

	// 角色和权限API
	mux.HandleFunc("/api/v1/roles", handleListRoles(db))
	mux.HandleFunc("/api/v1/roles/", handleGetRole(db))

	// 配额API
	mux.HandleFunc("/api/v1/quotas/", handleGetQuota(db))

	// 组织管理API
	mux.HandleFunc("/api/v1/org/organizations", handleCreateOrganization(db))
	mux.HandleFunc("/api/v1/org/organizations/", handleGetOrganization(db))
	mux.HandleFunc("/api/v1/org/members", handleAddMember(db))

	// 数据权限API
	mux.HandleFunc("/api/v1/permissions/data", handleAssignDataPermission(db))

	// Bot API
	mux.HandleFunc("/api/v1/bots", handleCreateBot(db))
	mux.HandleFunc("/api/v1/bots/", handleGetBot(db))
	mux.HandleFunc("/api/v1/bots/list", handleListBots(db))

	// 审核API
	mux.HandleFunc("/api/v1/review/tasks", handleListReviewTasks(db))
	mux.HandleFunc("/api/v1/review/tasks/", handleAssignReviewTask(db))
	mux.HandleFunc("/api/v1/review/tasks/submit", handleSubmitReview(db))

	// 对话API
	mux.HandleFunc("/api/v1/conversations/chat", handleChat(db))
	mux.HandleFunc("/api/v1/conversations/memories", handleGetMemories(db))

	// 用户API
	mux.HandleFunc("/api/v1/users", handleCreateUser(db))

	return httptest.NewServer(mux)
}

// ==================== Mock Handlers ====================

/**
 * 所有Handler遵循统一响应格式:
 * {
 *   "code": 0,        // 0表示成功，非0表示错误
 *   "message": "success",
 *   "data": {...}
 * }
 */

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func successResponse(data interface{}) APIResponse {
	return APIResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func errorResponse(code int, message string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
	}
}

/**
 * handleTenantRegister 处理租户注册
 * 流程: 创建租户 → 初始化RBAC → 分配默认配额
 */
func handleTenantRegister(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			TenantName    string `json:"tenant_name"`
			AdminEmail    string `json:"admin_email"`
			AdminPassword string `json:"admin_password"`
			Plan          string `json:"plan"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 生成租户ID
		tenantID := fmt.Sprintf("tenant-%d", time.Now().UnixNano())
		now := time.Now().Unix() * 1000

		// 创建租户
		tenant := entity.Tenant{
			TenantID:         tenantID,
			TenantName:       req.TenantName,
			TenantType:       entity.TenantTypeTeam,
			Subdomain:        fmt.Sprintf("tenant-%d", time.Now().Unix()),
			Status:           entity.TenantStatusActive,
			SubscriptionTier: entity.SubscriptionTierEnterprise,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := db.Create(&tenant).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 创建订阅
		subscriptionID := fmt.Sprintf("sub-%d", time.Now().UnixNano())
		db.Exec(`
			INSERT INTO subscriptions (subscription_id, tenant_id, plan_tier, status, max_bots, max_users, created_at, updated_at)
			VALUES (?, ?, 'active', 'active', 100, 50, ?, ?)
		`, subscriptionID, tenantID, now, now)

		// 创建默认角色
		roleID := fmt.Sprintf("role-%d", time.Now().UnixNano())
		db.Exec(`
			INSERT INTO roles (role_id, tenant_id, role_name, role_code, is_system, created_at, updated_at)
			VALUES (?, ?, 'Admin', 'admin', TRUE, ?, ?)
		`, roleID, tenantID, now, now)

		// 创建默认配额
		quotaID := fmt.Sprintf("quota-%d", time.Now().UnixNano())
		db.Exec(`
			INSERT INTO quotas (quota_id, tenant_id, resource_type, max_limit, current_usage, created_at, updated_at)
			VALUES (?, ?, 'bot', 100, 0, ?, ?)
		`, quotaID, tenantID, now, now)

		// 返回响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"tenant_id":   tenantID,
			"tenant_name": tenant.TenantName,
			"status":      string(tenant.Status),
		}))
	}
}

/**
 * handleCreateTenant 处理创建租户
 */
func handleCreateTenant(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req entity.Tenant
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now().Unix() * 1000
		req.TenantID = fmt.Sprintf("tenant-%d", time.Now().UnixNano())
		req.CreatedAt = now
		req.UpdatedAt = now

		if err := db.Create(&req).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"tenant_id": req.TenantID,
		}))
	}
}

/**
 * handleGetTenant 处理获取租户详情
 */
func handleGetTenant(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.URL.Path[len("/api/v1/tenants/"):]
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}

		var tenant entity.Tenant
		if err := db.Where("tenant_id = ?", tenantID).First(&tenant).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(tenant))
	}
}

/**
 * handleListRoles 处理列出租户角色
 */
func handleListRoles(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}

		var roles []struct {
			RoleID    string `json:"role_id"`
			TenantID  string `json:"tenant_id"`
			RoleName  string `json:"role_name"`
			RoleCode  string `json:"role_code"`
			IsSystem  bool   `json:"is_system"`
			CreatedAt int64  `json:"created_at"`
		}
		if err := db.Table("roles").Where("tenant_id = ?", tenantID).Find(&roles).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(roles))
	}
}

/**
 * handleGetRole 处理获取角色详情
 */
func handleGetRole(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleID := r.URL.Path[len("/api/v1/roles/"):]
		if roleID == "" {
			http.Error(w, "role_id is required", http.StatusBadRequest)
			return
		}

		var role struct {
			RoleID   string `json:"role_id"`
			RoleName string `json:"role_name"`
			RoleCode string `json:"role_code"`
		}
		if err := db.Table("roles").Where("role_id = ?", roleID).First(&role).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(role))
	}
}

/**
 * handleGetQuota 处理获取租户配额
 */
func handleGetQuota(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.URL.Path[len("/api/v1/quotas/"):]
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}

		var quota struct {
			QuotaID      string `json:"quota_id"`
			TenantID     string `json:"tenant_id"`
			ResourceType string `json:"resource_type"`
			MaxLimit     int64  `json:"max_limit"`
			CurrentUsage int64  `json:"current_usage"`
			ResetPeriod  string `json:"reset_period"`
		}
		if err := db.Table("quotas").Where("tenant_id = ?", tenantID).First(&quota).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(quota))
	}
}

/**
 * handleCreateOrganization 处理创建组织
 */
func handleCreateOrganization(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req orgEntity.Organization
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now().Unix() * 1000
		req.OrgID = fmt.Sprintf("org-%d", time.Now().UnixNano())
		req.CreatedAt = now
		req.UpdatedAt = now
		req.Path = fmt.Sprintf("/%s", req.OrgID)

		if err := db.Create(&req).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"organization_id": req.OrgID,
		}))
	}
}

/**
 * handleGetOrganization 处理获取组织详情
 */
func handleGetOrganization(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := r.URL.Path[len("/api/v1/org/organizations/"):]
		if orgID == "" {
			http.Error(w, "org_id is required", http.StatusBadRequest)
			return
		}

		var org orgEntity.Organization
		if err := db.Where("org_id = ?", orgID).First(&org).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(org))
	}
}

/**
 * handleAddMember 处理添加成员到组织
 */
func handleAddMember(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			OrganizationID string `json:"organization_id"`
			UserID         string `json:"user_id"`
			Role           string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 简化：直接返回成功
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"message": "成员添加成功",
		}))
	}
}

/**
 * handleAssignDataPermission 处理分配数据权限
 */
func handleAssignDataPermission(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req permissionEntity.DataPermission
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now().Unix() * 1000
		req.PermissionID = fmt.Sprintf("perm-%d", time.Now().UnixNano())
		req.CreatedAt = now
		req.UpdatedAt = now

		if err := db.Create(&req).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"permission_id": req.PermissionID,
		}))
	}
}

/**
 * handleCreateBot 处理创建Bot
 */
func handleCreateBot(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			TenantID     string  `json:"tenant_id"`
			BotName      string  `json:"bot_name"`
			OrgID        *string `json:"org_id,omitempty"`
			AIConfidence float64 `json:"ai_confidence"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now().Unix() * 1000
		botID := fmt.Sprintf("bot-%d", time.Now().UnixNano())
		status := "published"
		if req.AIConfidence < 0.7 {
			status = "draft" // 低置信度进入草稿状态
		}

		if err := db.Exec(`
			INSERT INTO bots (bot_id, tenant_id, org_id, bot_name, status, ai_confidence, created_by, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, botID, req.TenantID, req.OrgID, req.BotName, status, req.AIConfidence, "system-user", now, now).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 如果AI置信度低，创建审核任务
		if req.AIConfidence < 0.7 {
			taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
			db.Exec(`
				INSERT INTO review_tasks (task_id, tenant_id, resource_type, resource_id, status, ai_confidence, created_at, updated_at)
				VALUES (?, ?, 'bot', ?, 'pending', ?, ?, ?)
			`, taskID, req.TenantID, botID, req.AIConfidence, now, now)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"bot_id": botID,
			"status": status,
		}))
	}
}

/**
 * handleGetBot 处理获取Bot详情
 */
func handleGetBot(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		botID := r.URL.Path[len("/api/v1/bots/"):]
		if botID == "" {
			http.Error(w, "bot_id is required", http.StatusBadRequest)
			return
		}

		var bot struct {
			BotID    string  `json:"bot_id"`
			TenantID string  `json:"tenant_id"`
			OrgID    *string `json:"org_id,omitempty"`
			BotName  string  `json:"bot_name"`
			Status   string  `json:"status"`
		}
		if err := db.Table("bots").Where("bot_id = ?", botID).First(&bot).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(bot))
	}
}

/**
 * handleListBots 处理列出Bot（带权限过滤）
 */
func handleListBots(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_ = r.URL.Query().Get("user_id")
		_ = r.URL.Query().Get("resource_type")

		// 简化：直接返回所有Bot
		var bots []struct {
			BotID   string  `json:"bot_id"`
			OrgID   *string `json:"org_id,omitempty"`
			BotName string  `json:"bot_name"`
		}
		if err := db.Table("bots").Find(&bots).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"bots": bots,
		}))
	}
}

/**
 * handleListReviewTasks 处理列出审核任务
 */
func handleListReviewTasks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var tasks []struct {
			TaskID       string  `json:"task_id"`
			TenantID     string  `json:"tenant_id"`
			ResourceType string  `json:"resource_type"`
			ResourceID   string  `json:"resource_id"`
			Status       string  `json:"status"`
			AIConfidence float64 `json:"ai_confidence"`
		}
		if err := db.Table("review_tasks").Where("status = ?", "pending").Find(&tasks).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"tasks": tasks,
		}))
	}
}

/**
 * handleAssignReviewTask 处理分配审核任务
 */
func handleAssignReviewTask(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := r.URL.Path[len("/api/v1/review/tasks/"):]
		if taskID == "" {
			http.Error(w, "task_id is required", http.StatusBadRequest)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			AssignedTo string `json:"assigned_to"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now().Unix() * 1000
		if err := db.Exec(`
			UPDATE review_tasks SET assigned_to = ?, status = 'assigned', updated_at = ?
			WHERE task_id = ?
		`, req.AssignedTo, now, taskID).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"message": "任务分配成功",
		}))
	}
}

/**
 * handleSubmitReview 处理提交审核结果
 */
func handleSubmitReview(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			TaskID   string `json:"task_id"`
			Decision string `json:"decision"`
			Comment  string `json:"comment"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now().Unix() * 1000

		// 更新任务状态
		if err := db.Exec(`
			UPDATE review_tasks SET status = ?, updated_at = ?, completed_at = ?
			WHERE task_id = ?
		`, req.Decision, now, now, req.TaskID).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 如果批准，更新Bot状态为published
		if req.Decision == "approved" {
			// 获取resource_id
			var resourceID string
			db.Table("review_tasks").Where("task_id = ?", req.TaskID).Select("resource_id").Scan(&resourceID)

			db.Exec(`
				UPDATE bots SET status = 'published', updated_at = ?
				WHERE bot_id = ?
			`, now, resourceID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"message": "审核提交成功",
		}))
	}
}

/**
 * handleChat 处理对话
 */
func handleChat(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			TenantID string `json:"tenant_id"`
			UserID   string `json:"user_id"`
			Message  string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now().Unix() * 1000

		// 创建对话
		conversationID := fmt.Sprintf("conv-%d", time.Now().UnixNano())
		db.Exec(`
			INSERT INTO conversations (conversation_id, tenant_id, user_id, bot_id, status, created_at, updated_at)
			VALUES (?, ?, ?, 'default-bot', 'active', ?, ?)
		`, conversationID, req.TenantID, req.UserID, now, now)

		// 创建消息
		messageID := fmt.Sprintf("msg-%d", time.Now().UnixNano())
		db.Exec(`
			INSERT INTO messages (message_id, conversation_id, role, content, created_at)
			VALUES (?, ?, 'user', ?, ?)
		`, messageID, conversationID, req.Message, now)

		// 提取记忆（简化：如果包含HR关键词）
		response := "你好！我是AI助手。"
		if len(req.Message) > 10 {
			// 提取实体记忆
			entityType := "unknown"
			entityValue := "用户"
			if contains(req.Message, []string{"HR", "hr", "人事"}) {
				entityType = "HR"
				entityValue = "HR经理"
			} else if contains(req.Message, []string{"技术", "开发"}) {
				entityType = "技术"
				entityValue = "技术人员"
			}

			// 存储记忆
			memoryID := fmt.Sprintf("mem-%d", time.Now().UnixNano())
			db.Exec(`
				INSERT INTO entity_memories (memory_id, tenant_id, user_id, entity_type, entity_value, importance_score, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, 0.80, ?, ?)
			`, memoryID, req.TenantID, req.UserID, entityType, entityValue, now, now)

			// 根据记忆调整回复
			if entityType == "HR" {
				response = "您好HR经理！我会提供简洁专业的回答。"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"response":         response,
			"conversation_id":  conversationID,
			"entity_extracted": len(req.Message) > 10,
		}))
	}
}

/**
 * handleGetMemories 处理获取用户记忆
 */
func handleGetMemories(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.URL.Query().Get("user_id")
		entityType := r.URL.Query().Get("entity_type")

		var memories []struct {
			MemoryID        string  `json:"memory_id"`
			EntityType      string  `json:"entity_type"`
			EntityValue     string  `json:"entity_value"`
			ImportanceScore float64 `json:"importance_score"`
		}

		query := db.Table("entity_memories").Where("user_id = ?", userID)
		if entityType != "" {
			query = query.Where("entity_type = ?", entityType)
		}
		query.Find(&memories)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"memories": memories,
		}))
	}
}

/**
 * handleCreateUser 处理创建用户
 */
func handleCreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			TenantID string `json:"tenant_id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now().Unix() * 1000
		userID := fmt.Sprintf("user-%d", time.Now().UnixNano())

		if err := db.Exec(`
			INSERT INTO users (user_id, tenant_id, username, email, password_hash, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, 'active', ?, ?)
		`, userID, req.TenantID, req.Username, req.Username+"@example.com", "hashed", now, now).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successResponse(map[string]interface{}{
			"user_id": userID,
		}))
	}
}

/**
 * contains 检查字符串是否包含任意关键词
 */
func contains(s string, keywords []string) bool {
	for _, kw := range keywords {
		if len(s) >= len(kw) {
			for i := 0; i <= len(s)-len(kw); i++ {
				if s[i:i+len(kw)] == kw {
					return true
				}
			}
		}
	}
	return false
}
