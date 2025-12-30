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
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestTenantScopePlugin_Query 测试SELECT查询租户过滤
func TestTenantScopePlugin_Query(t *testing.T) {
	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 创建测试表
	type TestModel struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Name     string
	}

	err = db.AutoMigrate(&TestModel{})
	require.NoError(t, err)

	// 插入测试数据
	tenant1 := "tenant-001"
	tenant2 := "tenant-002"

	db.Create(&TestModel{ID: "1", TenantID: tenant1, Name: "Test1"})
	db.Create(&TestModel{ID: "2", TenantID: tenant1, Name: "Test2"})
	db.Create(&TestModel{ID: "3", TenantID: tenant2, Name: "Test3"})

	// 安装插件
	plugin := &TenantScopePlugin{
		TenantIDGetter: func(ctx context.Context) string {
			return tenant1
		},
	}
	err = plugin.Initialize(db)
	require.NoError(t, err)

	// 测试1: 正常查询（应该只返回tenant1的数据）
	t.Run("QueryWithTenantFilter", func(t *testing.T) {
		ctx := context.Background()
		var results []TestModel

		queryDB := db.WithContext(ctx)
		err = queryDB.Find(&results).Error
		require.NoError(t, err)

		// 应该只返回tenant1的数据
		assert.Equal(t, 2, len(results))
		for _, result := range results {
			assert.Equal(t, tenant1, result.TenantID)
		}
	})

	// 测试2: 系统表查询（应该不添加租户过滤）
	t.Run("QuerySystemTable", func(t *testing.T) {
		plugin.SkipCallbackTables["test_models"] = true

		ctx := context.Background()
		var results []TestModel

		err = db.WithContext(ctx).Find(&results).Error
		require.NoError(t, err)

		// 应该返回所有数据
		assert.Equal(t, 3, len(results))

		// 清理
		delete(plugin.SkipCallbackTables, "test_models")
	})

	// 测试3: 无租户ID的context（应该跳过过滤）
	t.Run("QueryWithoutTenantID", func(t *testing.T) {
		plugin.TenantIDGetter = func(ctx context.Context) string {
			return ""
		}

		ctx := context.Background()
		var results []TestModel

		err = db.WithContext(ctx).Find(&results).Error
		require.NoError(t, err)

		// 应该返回所有数据（因为没有租户ID）
		assert.Equal(t, 3, len(results))
	})
}

// TestTenantScopePlugin_Create 测试INSERT租户注入
func TestTenantScopePlugin_Create(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type TestModel struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Name     string
	}

	err = db.AutoMigrate(&TestModel{})
	require.NoError(t, err)

	// 安装插件
	tenantID := "tenant-001"
	plugin := &TenantScopePlugin{
		TenantIDGetter: func(ctx context.Context) string {
			return tenantID
		},
	}
	err = plugin.Initialize(db)
	require.NoError(t, err)

	// 测试1: 创建记录（应该自动注入tenant_id）
	t.Run("CreateWithTenantInjection", func(t *testing.T) {
		ctx := context.Background()
		record := TestModel{ID: "1", Name: "Test"}

		err = db.WithContext(ctx).Create(&record).Error
		require.NoError(t, err)

		// 验证tenant_id已自动注入
		assert.Equal(t, tenantID, record.TenantID)

		// 从数据库读取验证
		var saved TestModel
		err = db.First(&saved, "id = ?", "1").Error
		require.NoError(t, err)
		assert.Equal(t, tenantID, saved.TenantID)
	})

	// 测试2: 创建时手动指定tenant_id（应该被覆盖）
	t.Run("CreateWithManualTenantID", func(t *testing.T) {
		ctx := context.Background()
		record := TestModel{ID: "2", Name: "Test2", TenantID: "manual-tenant"}

		err = db.WithContext(ctx).Create(&record).Error
		require.NoError(t, err)

		// 验证tenant_id被覆盖为context中的值
		assert.Equal(t, tenantID, record.TenantID)
	})
}

// TestTenantScopePlugin_Update 测试UPDATE租户保护
func TestTenantScopePlugin_Update(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type TestModel struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Name     string
	}

	err = db.AutoMigrate(&TestModel{})
	require.NoError(t, err)

	// 插入测试数据
	tenantID := "tenant-001"
	db.Create(&TestModel{ID: "1", TenantID: tenantID, Name: "Test"})

	// 安装插件
	plugin := &TenantScopePlugin{
		TenantIDGetter: func(ctx context.Context) string {
			return tenantID
		},
	}
	err = plugin.Initialize(db)
	require.NoError(t, err)

	// 测试1: 正常更新（应该成功）
	t.Run("UpdateNormalFields", func(t *testing.T) {
		ctx := context.Background()

		err = db.WithContext(ctx).
			Model(&TestModel{}).
			Where("id = ?", "1").
			Update("name", "Updated").Error
		require.NoError(t, err)

		// 验证更新成功
		var result TestModel
		err = db.First(&result, "id = ?", "1").Error
		require.NoError(t, err)
		assert.Equal(t, "Updated", result.Name)
		assert.Equal(t, tenantID, result.TenantID)
	})

	// 测试2: 尝试修改tenant_id（应该失败）
	t.Run("UpdateTenantID", func(t *testing.T) {
		ctx := context.Background()

		err = db.WithContext(ctx).
			Model(&TestModel{}).
			Where("id = ?", "1").
			Update("tenant_id", "other-tenant").Error

		// 应该返回错误
		assert.Error(t, err)
	})
}

// TestTenantScopePlugin_Delete 测试DELETE租户过滤
func TestTenantScopePlugin_Delete(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type TestModel struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Name     string
	}

	err = db.AutoMigrate(&TestModel{})
	require.NoError(t, err)

	// 插入测试数据
	tenant1 := "tenant-001"
	tenant2 := "tenant-002"
	db.Create(&TestModel{ID: "1", TenantID: tenant1, Name: "Test1"})
	db.Create(&TestModel{ID: "2", TenantID: tenant2, Name: "Test2"})

	// 安装插件
	plugin := &TenantScopePlugin{
		TenantIDGetter: func(ctx context.Context) string {
			return tenant1
		},
	}
	err = plugin.Initialize(db)
	require.NoError(t, err)

	// 测试: 删除记录（应该只删除当前租户的数据）
	t.Run("DeleteWithTenantFilter", func(t *testing.T) {
		ctx := context.Background()

		err = db.WithContext(ctx).
			Where(&TestModel{TenantID: tenant1}).
			Delete(&TestModel{}).Error
		require.NoError(t, err)

		// 验证tenant1的数据已删除
		var count1 int64
		db.Model(&TestModel{}).Where("tenant_id = ?", tenant1).Count(&count1)
		assert.Equal(t, int64(0), count1)

		// 验证tenant2的数据仍然存在
		var count2 int64
		db.Model(&TestModel{}).Where("tenant_id = ?", tenant2).Count(&count2)
		assert.Equal(t, int64(1), count2)
	})
}

// TestValidateTenantAccess 测试租户访问权限验证
func TestValidateTenantAccess(t *testing.T) {
	tenant1 := "tenant-001"
	tenant2 := "tenant-002"

	// 测试1: 相同租户（应该成功）
	t.Run("SameTenant", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, tenant1)
		err := ValidateTenantAccess(ctx, tenant1)
		assert.NoError(t, err)
	})

	// 测试2: 不同租户（应该失败）
	t.Run("DifferentTenant", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, tenant1)
		err := ValidateTenantAccess(ctx, tenant2)
		assert.Error(t, err)
	})

	// 测试3: 无效的request_tenant_id（应该失败）
	t.Run("InvalidRequestTenantID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, "")
		err := ValidateTenantAccess(ctx, tenant1)
		assert.Error(t, err)
	})
}

// TestEnsureTenantDataIntegrity 测试租户数据完整性验证
func TestEnsureTenantDataIntegrity(t *testing.T) {
	type TestResource struct {
		TenantID string
		Name     string
	}

	tenantID := "tenant-001"

	// 测试1: 所有资源属于同一租户（应该成功）
	t.Run("AllSameTenant", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, tenantID)

		resources := []TestResource{
			{TenantID: tenantID, Name: "R1"},
			{TenantID: tenantID, Name: "R2"},
			{TenantID: tenantID, Name: "R3"},
		}

		err := EnsureTenantDataIntegrity(ctx, resources)
		assert.NoError(t, err)
	})

	// 测试2: 资源属于不同租户（应该失败）
	t.Run("MixedTenant", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, tenantID)

		resources := []TestResource{
			{TenantID: tenantID, Name: "R1"},
			{TenantID: "other-tenant", Name: "R2"},
		}

		err := EnsureTenantDataIntegrity(ctx, resources)
		assert.Error(t, err)
	})

	// 测试3: 指针类型资源
	t.Run("PointerResources", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, tenantID)

		resources := []*TestResource{
			{TenantID: tenantID, Name: "R1"},
			{TenantID: tenantID, Name: "R2"},
		}

		err := EnsureTenantDataIntegrity(ctx, resources)
		assert.NoError(t, err)
	})

	// 测试4: 非slice类型（应该失败）
	t.Run("NotSlice", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, tenantID)

		resource := TestResource{TenantID: tenantID, Name: "R1"}

		err := EnsureTenantDataIntegrity(ctx, resource)
		assert.Error(t, err)
	})
}

// TestGetTenantIDFromContext 测试从context获取tenant_id
func TestGetTenantIDFromContext(t *testing.T) {
	tenantID := "tenant-001"

	// 测试1: 从Hertz context获取
	t.Run("FromHertzContext", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, tenantID)
		result := GetTenantIDFromContext(ctx)
		assert.Equal(t, tenantID, result)
	})

	// 测试2: 没有tenant_id（返回默认值）
	t.Run("NoTenantID", func(t *testing.T) {
		ctx := context.Background()
		result := GetTenantIDFromContext(ctx)
		assert.Equal(t, DefaultTenantID, result)
	})
}

// BenchmarkTenantScopePlugin_Query 性能测试
func BenchmarkTenantScopePlugin_Query(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	type TestModel struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Name     string
	}

	db.AutoMigrate(&TestModel{})

	// 插入1000条数据
	for i := 0; i < 1000; i++ {
		db.Create(&TestModel{
			ID:       string(rune(i)),
			TenantID: "tenant-001",
			Name:     "Test",
		})
	}

	// 安装插件
	plugin := &TenantScopePlugin{
		TenantIDGetter: func(ctx context.Context) string {
			return "tenant-001"
		},
	}
	plugin.Initialize(db)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var results []TestModel
		db.WithContext(ctx).Find(&results)
	}
}

// BenchmarkTenantScopePlugin_Insert 性能测试
func BenchmarkTenantScopePlugin_Insert(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	type TestModel struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Name     string
	}

	db.AutoMigrate(&TestModel{})

	// 安装插件
	plugin := &TenantScopePlugin{
		TenantIDGetter: func(ctx context.Context) string {
			return "tenant-001"
		},
	}
	plugin.Initialize(db)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		record := TestModel{
			ID:   string(rune(i)),
			Name: "Test",
		}
		db.WithContext(ctx).Create(&record)
	}
}
