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

package fixtures

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
)

// ==================== 租户相关 ====================

// CreateTestTenant 创建测试租户
// 职责: 创建测试用的租户数据
func CreateTestTenant(t *testing.T, db *sql.DB, name string) *entity.Tenant {
	tenant := &entity.Tenant{
		TenantID:          fmt.Sprintf("tenant-%d", time.Now().UnixNano()),
		TenantName:        name,
		Status:            entity.TenantStatusActive,
		SubscriptionTier:  entity.SubscriptionTierFree,
		IsolationStrategy: entity.IsolationStrategyRowLevel,
		ContactEmail:      fmt.Sprintf("%s@example.com", name),
		CreatedAt:         time.Now().Unix(),
		UpdatedAt:         time.Now().Unix(),
	}

	query := `
		INSERT INTO tenants (
			tenant_id, tenant_name, tenant_type, subdomain, status, subscription_tier, isolation_strategy,
			contact_email, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query,
		tenant.TenantID, tenant.TenantName, "individual", "test-"+name, tenant.Status,
		tenant.SubscriptionTier, tenant.IsolationStrategy, tenant.ContactEmail,
		tenant.CreatedAt, tenant.UpdatedAt,
	)
	assert.NoError(t, err, "创建测试租户失败")

	return tenant
}

// ==================== 用户相关 ====================

// CreateTestUser 创建测试用户
// 职责: 创建测试用的用户数据
func CreateTestUser(t *testing.T, db *sql.DB, tenantID, username string) *User {
	userID := fmt.Sprintf("user-%d", time.Now().UnixNano())

	user := &User{
		UserID:        userID,
		TenantID:      tenantID,
		Username:      username,
		Email:         fmt.Sprintf("%s@example.com", username),
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	query := `
		INSERT INTO users (
			user_id, tenant_id, username, email, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query,
		user.UserID, user.TenantID, user.Username, user.Email,
		user.Status, user.CreatedAt, user.UpdatedAt,
	)
	assert.NoError(t, err, "创建测试用户失败")

	return user
}

// User 简化的用户结构
type User struct {
	UserID    string
	TenantID  string
	Username  string
	Email     string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ==================== Bot相关 ====================

// CreateTestBot 创建测试Bot
// 职责: 创建测试用的Bot数据
func CreateTestBot(t *testing.T, db *sql.DB, tenantID, name string) *Bot {
	bot := &Bot{
		BotID:      fmt.Sprintf("bot-%d", time.Now().UnixNano()),
		TenantID:   tenantID,
		Name:       name,
		Status:     "draft",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	query := `
		INSERT INTO bots (
			bot_id, tenant_id, name, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query,
		bot.BotID, bot.TenantID, bot.Name, bot.Status,
		bot.CreatedAt, bot.UpdatedAt,
	)
	assert.NoError(t, err, "创建测试Bot失败")

	return bot
}

// Bot 简化的Bot结构
type Bot struct {
	BotID      string
	TenantID   string
	Name       string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ListBotsByTenant 列出租户的所有Bot
// 职责: 查询指定租户的Bot列表
func ListBotsByTenant(t *testing.T, db *sql.DB, tenantID string) ([]*Bot, error) {
	query := `
		SELECT bot_id, tenant_id, name, status, created_at, updated_at
		FROM bots
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bots []*Bot
	for rows.Next() {
		var bot Bot
		err := rows.Scan(&bot.BotID, &bot.TenantID, &bot.Name, &bot.Status, &bot.CreatedAt, &bot.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bots = append(bots, &bot)
	}

	return bots, nil
}

// ==================== 角色相关 ====================

// CreateTestRole 创建测试角色
// 职责: 创建测试用的角色数据
func CreateTestRole(t *testing.T, db *sql.DB, tenantID, name string) *Role {
	role := &Role{
		RoleID:      fmt.Sprintf("role-%d", time.Now().UnixNano()),
		TenantID:    tenantID,
		Name:        name,
		Description: fmt.Sprintf("测试角色: %s", name),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO roles (
			role_id, tenant_id, name, description, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query,
		role.RoleID, role.TenantID, role.Name, role.Description,
		role.CreatedAt, role.UpdatedAt,
	)
	assert.NoError(t, err, "创建测试角色失败")

	return role
}

// Role 简化的角色结构
type Role struct {
	RoleID      string
	TenantID    string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AssignRoleToUser 分配角色给用户
// 职责: 为用户分配角色
func AssignRoleToUser(t *testing.T, db *sql.DB, userID, roleID string) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES (?, ?, ?)
	`
	_, err := db.Exec(query, userID, roleID, time.Now())
	return err
}

// ==================== 权限相关 ====================

// CreateDataPermission 创建数据权限
// 职责: 创建数据权限记录
func CreateDataPermission(t *testing.T, db *sql.DB, perm interface{}) error {
	// 使用类型断言处理不同类型的权限
	switch p := perm.(type) {
	case *DataPermission:
		query := `
			INSERT INTO data_permissions (
				permission_id, role_id, resource_type, action, scope, created_at
			) VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err := db.Exec(query, p.PermissionID, p.RoleID, p.ResourceType, p.Action, p.Scope, time.Now())
		return err
	default:
		return fmt.Errorf("unknown permission type")
	}
}

// DataPermission 数据权限
type DataPermission struct {
	PermissionID string
	RoleID       string
	ResourceType string
	Action       string
	Scope        string
}

// CreateFieldPermission 创建字段权限
// 职责: 创建字段权限记录
func CreateFieldPermission(t *testing.T, db *sql.DB, perm interface{}) error {
	switch p := perm.(type) {
	case *FieldPermission:
		query := `
			INSERT INTO field_permissions (
				permission_id, role_id, resource_type, field_name, is_allowed, permission_type, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err := db.Exec(query,
			p.PermissionID, p.RoleID, p.ResourceType, p.FieldName,
			p.IsAllowed, p.PermissionType, time.Now(),
		)
		return err
	default:
		return fmt.Errorf("unknown field permission type")
	}
}

// FieldPermission 字段权限
type FieldPermission struct {
	PermissionID   string
	RoleID         string
	ResourceType   string
	FieldName      string
	IsAllowed      bool
	PermissionType string
}

// ==================== 知识库相关 ====================

// CreateTestKnowledge 创建测试知识库
// 职责: 创建测试用的知识库数据
func CreateTestKnowledge(t *testing.T, db *sql.DB, tenantID, botID, name string) *Knowledge {
	knowledge := &Knowledge{
		KnowledgeID: fmt.Sprintf("knowledge-%d", time.Now().UnixNano()),
		TenantID:    tenantID,
		BotID:       botID,
		Name:        name,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO knowledge (
			knowledge_id, tenant_id, bot_id, name, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query,
		knowledge.KnowledgeID, knowledge.TenantID, knowledge.BotID,
		knowledge.Name, knowledge.Status, knowledge.CreatedAt, knowledge.UpdatedAt,
	)
	assert.NoError(t, err, "创建测试知识库失败")

	return knowledge
}

// Knowledge 简化的知识库结构
type Knowledge struct {
	KnowledgeID string
	TenantID    string
	BotID       string
	Name        string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ==================== 会话相关 ====================

// CreateTestConversation 创建测试会话
// 职责: 创建测试用的会话数据
func CreateTestConversation(t *testing.T, db *sql.DB, tenantID, userID, botID string) *Conversation {
	conv := &Conversation{
		ConversationID: fmt.Sprintf("conv-%d", time.Now().UnixNano()),
		TenantID:       tenantID,
		UserID:         userID,
		BotID:          botID,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `
		INSERT INTO conversations (
			conversation_id, tenant_id, user_id, bot_id, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query,
		conv.ConversationID, conv.TenantID, conv.UserID, conv.BotID,
		conv.Status, conv.CreatedAt, conv.UpdatedAt,
	)
	assert.NoError(t, err, "创建测试会话失败")

	return conv
}

// Conversation 简化的会话结构
type Conversation struct {
	ConversationID string
	TenantID       string
	UserID         string
	BotID          string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ==================== 通用辅助函数 ====================

// CleanupTable 清理表数据
// 职责: 清理指定表的所有数据
func CleanupTable(t *testing.T, db *sql.DB, tableName string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s", tableName)
	_, err := db.Exec(query)
	return err
}

// WaitForOperation 等待操作完成
// 职责: 轮询检查操作是否完成
func WaitForOperation(ctx context.Context, checkFunc func() (bool, error), maxAttempts int, interval time.Duration) error {
	for i := 0; i < maxAttempts; i++ {
		done, err := checkFunc()
		if err != nil {
			return err
		}
		if done {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
	return fmt.Errorf("operation did not complete after %d attempts", maxAttempts)
}
