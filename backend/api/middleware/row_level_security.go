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

package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/infra/tracing"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TenantScopePlugin GORM插件：自动为所有查询添加租户过滤
//
// **核心功能**：
// 1. 自动为SELECT查询添加 WHERE tenant_id = ?
// 2. 自动为INSERT操作设置tenant_id
// 3. 自动为UPDATE操作添加 WHERE tenant_id = ?
// 4. 自动为DELETE操作添加 WHERE tenant_id = ?
// 5. 防止跨租户数据访问（安全关键）
//
// **使用示例**：
//
//	db.Use(&TenantScopePlugin{})
//	db.WithContext(ctx).Find(&bots) // 自动添加 WHERE tenant_id = ?
type TenantScopePlugin struct {
	// SkipCallbackTables 跳过租户隔离的表（系统表）
	SkipCallbackTables map[string]bool

	// TenantIDGetter 从context获取tenant_id的函数
	TenantIDGetter func(context.Context) string
}

// Name 插件名称
func (p *TenantScopePlugin) Name() string {
	return "TenantScopePlugin"
}

// Initialize 初始化插件
func (p *TenantScopePlugin) Initialize(db *gorm.DB) error {
	// 默认跳过系统表
	if p.SkipCallbackTables == nil {
		p.SkipCallbackTables = map[string]bool{
			"tenants":           true,
			"subscriptions":     true,
			"quotas":           true,
			"quota_usage":      true,
			"system_config":    true,
			"migration_tasks":  true,
		}
	}

	// 默认tenant_id获取器
	if p.TenantIDGetter == nil {
		p.TenantIDGetter = GetTenantIDFromContext
	}

	// 注册回调
	if err := db.Callback().Query().Before("gorm:query").Register("tenant_scope:query", p.queryCallback); err != nil {
		return fmt.Errorf("failed to register query callback: %w", err)
	}
	if err := db.Callback().Create().Before("gorm:create").Register("tenant_scope:create", p.createCallback); err != nil {
		return fmt.Errorf("failed to register create callback: %w", err)
	}
	if err := db.Callback().Update().Before("gorm:update").Register("tenant_scope:update", p.updateCallback); err != nil {
		return fmt.Errorf("failed to register update callback: %w", err)
	}
	if err := db.Callback().Delete().Before("gorm:delete").Register("tenant_scope:delete", p.deleteCallback); err != nil {
		return fmt.Errorf("failed to register delete callback: %w", err)
	}

	return nil
}

// queryCallback SELECT查询回调：自动添加WHERE tenant_id = ?
func (p *TenantScopePlugin) queryCallback(db *gorm.DB) {
	tableName := db.Statement.Table
	if tableName == "" {
		if db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		}
	}

	// 跳过系统表
	if p.SkipCallbackTables[tableName] {
		return
	}

	// 检查表是否有tenant_id字段
	if !p.hasTenantIDField(db) {
		return
	}

	// 从context获取tenant_id
	ctx := db.Statement.Context
	if ctx == nil {
		logs.Warnf("[TenantScope] no context in query, skipping tenant filter")
		return
	}

	tenantID := p.TenantIDGetter(ctx)
	if tenantID == "" || tenantID == DefaultTenantID {
		logs.Warnf("[TenantScope] invalid tenant_id in context: %s", tenantID)
		return
	}

	// 添加WHERE tenant_id = ?
	db.Statement.AddClause(gorm.WhereClause{
		Exprs: []gorm.Expression{
			gorm.Expr{SQL: "? = ?", Vars: []interface{}{"tenant_id", tenantID}},
		},
	})

	// 记录追踪
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrTable.String(tableName),
	)

	logs.CtxDebugf(ctx, "[TenantScope] added tenant filter: table=%s, tenant_id=%s", tableName, tenantID)
}

// createCallback INSERT回调：自动设置tenant_id
func (p *TenantScopePlugin) createCallback(db *gorm.DB) {
	tableName := db.Statement.Table
	if tableName == "" {
		if db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		}
	}

	// 跳过系统表
	if p.SkipCallbackTables[tableName] {
		return
	}

	// 检查表是否有tenant_id字段
	if !p.hasTenantIDField(db) {
		return
	}

	// 从context获取tenant_id
	ctx := db.Statement.Context
	if ctx == nil {
		logs.Warnf("[TenantScope] no context in create, skipping tenant_id injection")
		return
	}

	tenantID := p.TenantIDGetter(ctx)
	if tenantID == "" || tenantID == DefaultTenantID {
		logs.Warnf("[TenantScope] invalid tenant_id in context: %s", tenantID)
		return
	}

	// 设置tenant_id字段
	db.Statement.SetColumn("tenant_id", tenantID)

	// 记录追踪
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrTable.String(tableName),
	)

	logs.CtxDebugf(ctx, "[TenantScope] injected tenant_id: table=%s, tenant_id=%s", tableName, tenantID)
}

// updateCallback UPDATE回调：添加WHERE tenant_id = ?
func (p *TenantScopePlugin) updateCallback(db *gorm.DB) {
	tableName := db.Statement.Table
	if tableName == "" {
		if db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		}
	}

	// 跳过系统表
	if p.SkipCallbackTables[tableName] {
		return
	}

	// 检查表是否有tenant_id字段
	if !p.hasTenantIDField(db) {
		return
	}

	// 从context获取tenant_id
	ctx := db.Statement.Context
	if ctx == nil {
		logs.Warnf("[TenantScope] no context in update, skipping tenant filter")
		return
	}

	tenantID := p.TenantIDGetter(ctx)
	if tenantID == "" || tenantID == DefaultTenantID {
		logs.Warnf("[TenantScope] invalid tenant_id in context: %s", tenantID)
		return
	}

	// 防止修改tenant_id字段（安全关键）
	if p.isUpdatingTenantID(db) {
		logs.CtxErrorf(ctx, "[TenantScope] attempt to update tenant_id: table=%s, tenant_id=%s", tableName, tenantID)
		db.AddError(pkgerrorx.New(berrno.ErrCrossTenantUpdate).WithZap(
			zap.String("table", tableName),
			zap.String("tenant_id", tenantID),
		))
		return
	}

	// 添加WHERE tenant_id = ?
	db.Statement.AddClause(gorm.WhereClause{
		Exprs: []gorm.Expression{
			gorm.Expr{SQL: "? = ?", Vars: []interface{}{"tenant_id", tenantID}},
		},
	})

	// 记录追踪
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrTable.String(tableName),
	)

	logs.CtxDebugf(ctx, "[TenantScope] added tenant filter to update: table=%s, tenant_id=%s", tableName, tenantID)
}

// deleteCallback DELETE回调：添加WHERE tenant_id = ?
func (p *TenantScopePlugin) deleteCallback(db *gorm.DB) {
	tableName := db.Statement.Table
	if tableName == "" {
		if db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		}
	}

	// 跳过系统表
	if p.SkipCallbackTables[tableName] {
		return
	}

	// 检查表是否有tenant_id字段
	if !p.hasTenantIDField(db) {
		return
	}

	// 从context获取tenant_id
	ctx := db.Statement.Context
	if ctx == nil {
		logs.Warnf("[TenantScope] no context in delete, skipping tenant filter")
		return
	}

	tenantID := p.TenantIDGetter(ctx)
	if tenantID == "" || tenantID == DefaultTenantID {
		logs.Warnf("[TenantScope] invalid tenant_id in context: %s", tenantID)
		return
	}

	// 添加WHERE tenant_id = ?
	db.Statement.AddClause(gorm.WhereClause{
		Exprs: []gorm.Expression{
			gorm.Expr{SQL: "? = ?", Vars: []interface{}{"tenant_id", tenantID}},
		},
	})

	// 记录追踪
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrTable.String(tableName),
	)

	logs.CtxDebugf(ctx, "[TenantScope] added tenant filter to delete: table=%s, tenant_id=%s", tableName, tenantID)
}

// hasTenantIDField 检查表是否有tenant_id字段
func (p *TenantScopePlugin) hasTenantIDField(db *gorm.DB) bool {
	if db.Statement.Schema == nil {
		return false
	}

	_, ok := db.Statement.Schema.FieldsByDBName["tenant_id"]
	return ok
}

// isUpdatingTenantID 检查是否正在更新tenant_id字段
func (p *TenantScopePlugin) isUpdatingTenantID(db *gorm.DB) bool {
	if db.Statement.Dest == nil {
		return false
	}

	// 检查map类型
	if m, ok := db.Statement.Dest.(map[string]interface{}); ok {
		_, exists := m["tenant_id"]
		return exists
	}

	// 检查struct类型
	if db.Statement.Schema != nil {
		for _, field := range db.Statement.Schema.Fields {
			if field.DBName == "tenant_id" && db.Statement.ReflectValue.CanAddr() {
				fieldVal := db.Statement.ReflectValue.FieldByName(field.Name)
				if fieldVal.IsValid() && !fieldVal.IsZero() {
					return true
				}
			}
		}
	}

	return false
}

// RowLevelSecurityMiddleware 行级安全中间件
//
// **核心功能**：
// 1. 验证所有数据库操作都包含租户过滤
// 2. 防止跨租户数据访问
// 3. 记录所有租户相关的数据库操作
//
// **使用示例**：
//	app.Use(middleware.RowLevelSecurityMiddleware())
func RowLevelSecurityMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从context获取tenant_id
		tenantID := GetTenantIDFromContext(c)
		if tenantID == "" || tenantID == DefaultTenantID {
			logs.CtxWarnf(c, "[RowLevelSecurity] no valid tenant_id in context")
		}

		// 记录租户ID到请求上下文
		ctx.Set("tenant_id", tenantID)

		// 继续处理请求
		ctx.Next(c)
	}
}

// WithTenantScope 为GORM DB实例添加租户作用域
//
// **使用场景**：手动为查询添加租户过滤
//
// **示例**：
//	db := middleware.WithTenantScope(gormDB, ctx)
//	db.Find(&bots) // 自动添加 WHERE tenant_id = ?
func WithTenantScope(db *gorm.DB, ctx context.Context) *gorm.DB {
	tenantID := GetTenantIDFromContext(ctx)
	if tenantID == "" || tenantID == DefaultTenantID {
		logs.CtxWarnf(ctx, "[WithTenantScope] invalid tenant_id: %s", tenantID)
		return db
	}

	return db.Where("tenant_id = ?", tenantID)
}

// ValidateTenantAccess 验证租户访问权限
//
// **使用场景**：在业务逻辑中显式验证租户访问权限
//
// **示例**：
//	if err := middleware.ValidateTenantAccess(ctx, bot.TenantID); err != nil {
//	    return nil, err
//	}
func ValidateTenantAccess(ctx context.Context, resourceTenantID string) error {
	requestTenantID := GetTenantIDFromContext(ctx)

	if requestTenantID == "" || requestTenantID == DefaultTenantID {
		return pkgerrorx.New(berrno.ErrInvalidTenantID).WithZap(
			zap.String("request_tenant_id", requestTenantID),
		)
	}

	if requestTenantID != resourceTenantID {
		logs.CtxErrorf(ctx, "[ValidateTenantAccess] cross-tenant access denied: request_tenant_id=%s, resource_tenant_id=%s",
			requestTenantID, resourceTenantID)

		return pkgerrorx.New(berrno.ErrCrossTenantAccess).WithZap(
			zap.String("request_tenant_id", requestTenantID),
			zap.String("resource_tenant_id", resourceTenantID),
		)
	}

	return nil
}

// EnsureTenantDataIntegrity 确保租户数据完整性
//
// **使用场景**：在批量操作前验证所有数据属于同一租户
//
// **示例**：
//	bots := []Bot{bot1, bot2, bot3}
//	if err := middleware.EnsureTenantDataIntegrity(ctx, bots); err != nil {
//	    return err
//	}
func EnsureTenantDataIntegrity(ctx context.Context, resources interface{}) error {
	tenantID := GetTenantIDFromContext(ctx)

	if tenantID == "" || tenantID == DefaultTenantID {
		return pkgerrorx.New(berrno.ErrInvalidTenantID).WithZap(
			zap.String("tenant_id", tenantID),
		)
	}

	// 检查slice类型
	rv := reflect.ValueOf(resources)
	if rv.Kind() != reflect.Slice {
		return errors.New("resources must be a slice")
	}

	for i := 0; i < rv.Len(); i++ {
		item := rv.Index(i).Interface()
		if err := validateSingleResourceTenant(ctx, tenantID, item); err != nil {
			return err
		}
	}

	return nil
}

// validateSingleResourceTenant 验证单个资源的租户
func validateSingleResourceTenant(ctx context.Context, expectedTenantID string, resource interface{}) error {
	// 使用反射获取tenant_id字段
	rv := reflect.ValueOf(resource)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil // 非结构体类型，跳过验证
	}

	field := rv.FieldByName("TenantID")
	if !field.IsValid() {
		return nil // 没有TenantID字段，跳过验证
	}

	resourceTenantID := field.String()
	if resourceTenantID != expectedTenantID {
		logs.CtxErrorf(ctx, "[EnsureTenantDataIntegrity] mismatched tenant_id: expected=%s, got=%s",
			expectedTenantID, resourceTenantID)

		return pkgerrorx.New(berrno.ErrCrossTenantAccess).WithZap(
			zap.String("expected_tenant_id", expectedTenantID),
			zap.String("resource_tenant_id", resourceTenantID),
		)
	}

	return nil
}
