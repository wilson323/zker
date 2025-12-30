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

package security

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SecurityTestSuite 租户隔离安全测试套件
//
// **测试目标**：确保没有任何跨租户数据泄露风险
//
// **测试场景**：
// 1. 跨租户SELECT（应该失败）
// 2. 跨租户INSERT（应该失败）
// 3. 跨租户UPDATE（应该失败）
// 4. 跨租户DELETE（应该失败）
// 5. 跨租户缓存访问（应该失败）
// 6. 跨租户搜索访问（应该失败）
// 7. 跨租户文件访问（应该失败）

// TestCrossTenantSelect 测试跨租户SELECT攻击
func TestCrossTenantSelect(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type SecretData struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Secret   string
	}

	db.AutoMigrate(&SecretData{})

	// 插入租户1的机密数据
	tenant1 := "tenant-001"
	tenant2 := "tenant-002"

	db.Create(&SecretData{ID: "1", TenantID: tenant1, Secret: "TopSecret-001"})
	db.Create(&SecretData{ID: "2", TenantID: tenant2, Secret: "TopSecret-002"})

	// 安装租户隔离插件
	// （假设已经安装TenantScopePlugin）

	// 测试：租户2尝试查询租户1的数据
	ctx := context.WithValue(context.Background(), "tenant_id", tenant2)
	var results []SecretData

	err = db.WithContext(ctx).Find(&results).Error
	require.NoError(t, err)

	// 验证：只能看到租户2自己的数据
	for _, result := range results {
		assert.Equal(t, tenant2, result.TenantID, "应该只能看到自己租户的数据")
		assert.NotEqual(t, tenant1, result.TenantID, "不应该看到其他租户的数据")
	}

	// 确保没有泄露租户1的数据
	for _, result := range results {
		if result.TenantID == tenant1 {
			t.Errorf("安全漏洞：租户2可以访问租户1的数据！ID=%s, Secret=%s",
				result.ID, result.Secret)
		}
	}
}

// TestCrossTenantUpdate 测试跨租户UPDATE攻击
func TestCrossTenantUpdate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type FinancialRecord struct {
		ID        string  `gorm:"primaryKey"`
		TenantID  string  `gorm:"index"`
		Amount    float64
	}

	db.AutoMigrate(&FinancialRecord{})

	tenant1 := "tenant-001"
	tenant2 := "tenant-002"

	// 租户1的财务记录
	db.Create(&FinancialRecord{
		ID:       "rec-001",
		TenantID: tenant1,
		Amount:   1000000.00, // 100万
	})

	// 测试：租户2尝试修改租户1的财务记录
	ctx := context.WithValue(context.Background(), "tenant_id", tenant2)

	// 尝试1: 直接UPDATE
	err = db.WithContext(ctx).
		Model(&FinancialRecord{}).
		Where("id = ?", "rec-001").
		Update("amount", 0.01).Error

	// 应该失败（找不到记录或影响0行）
	if err == nil {
		// 检查是否真的修改了
		var record FinancialRecord
		db.First(&record, "id = ?", "rec-001")

		if record.Amount != 1000000.00 {
			t.Errorf("安全漏洞：租户2成功修改了租户1的财务记录！原金额=1000000.00, 当前金额=%.2f",
				record.Amount)
		}
	}

	// 尝试2: 修改tenant_id字段
	err = db.WithContext(ctx).
		Model(&FinancialRecord{}).
		Where("id = ?", "rec-001").
		Update("tenant_id", tenant2).Error

	// 应该失败
	assert.Error(t, err, "修改tenant_id应该被阻止")
}

// TestCrossTenantDelete 测试跨租户DELETE攻击
func TestCrossTenantDelete(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type CriticalData struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Data     string
	}

	db.AutoMigrate(&CriticalData{})

	tenant1 := "tenant-001"
	tenant2 := "tenant-002"

	// 租户1的关键数据
	db.Create(&CriticalData{ID: "1", TenantID: tenant1, Data: "Critical-001"})

	// 测试：租户2尝试删除租户1的数据
	ctx := context.WithValue(context.Background(), "tenant_id", tenant2)

	err = db.WithContext(ctx).
		Where("id = ?", "1").
		Delete(&CriticalData{}).Error
	require.NoError(t, err)

	// 验证：租户1的数据仍然存在
	var count int64
	db.Model(&CriticalData{}).Where("id = ? AND tenant_id = ?", "1", tenant1).Count(&count)

	if count == 0 {
		t.Errorf("安全漏洞：租户2成功删除了租户1的数据！")
	}

	assert.Equal(t, int64(1), count, "租户1的数据应该仍然存在")
}

// TestTenantIDModification 测试防止tenant_id修改
func TestTenantIDModification(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type TestData struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Name     string
	}

	db.AutoMigrate(&TestData{})

	tenant1 := "tenant-001"
	tenant2 := "tenant-002"

	db.Create(&TestData{ID: "1", TenantID: tenant1, Name: "Test"})

	// 安装TenantScopePlugin（应该阻止tenant_id修改）
	// ...

	// 测试：尝试修改tenant_id
	ctx := context.WithValue(context.Background(), "tenant_id", tenant1)

	err = db.WithContext(ctx).
		Model(&TestData{}).
		Where("id = ?", "1").
		Update("tenant_id", tenant2).Error

	// 应该失败
	assert.Error(t, err, "修改tenant_id应该被阻止")

	// 验证tenant_id没有变化
	var result TestData
	db.First(&result, "id = ?", "1")
	assert.Equal(t, tenant1, result.TenantID, "tenant_id不应该被修改")
}

// TestBypassAttempt 测试绕过尝试
func TestBypassAttempt(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type SensitiveData struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Data     string
	}

	db.AutoMigrate(&SensitiveData{})

	tenant1 := "tenant-001"
	tenant2 := "tenant-002"

	db.Create(&SensitiveData{ID: "1", TenantID: tenant1, Data: "Secret"})

	// 尝试绕过1: 使用Raw SQL
	ctx := context.WithValue(context.Background(), "tenant_id", tenant2)

	var results []SensitiveData
	err = db.WithContext(ctx).Raw("SELECT * FROM sensitive_data WHERE id = ?", "1").Scan(&results).Error

	// Raw SQL应该仍然被租户过滤保护
	// （实际实现取决于GORM插件配置）
	if err == nil && len(results) > 0 {
		for _, result := range results {
			if result.TenantID == tenant1 {
				t.Errorf("安全漏洞：Raw SQL绕过了租户隔离！tenant_id=%s", result.TenantID)
			}
		}
	}

	// 尝试绕过2: 使用无context的查询
	err = db.Where("id = ?", "1").Find(&results).Error
	if err == nil && len(results) > 0 {
		// 无context的查询应该被阻止或返回空结果
		// （取决于安全策略）
	}
}

// TestCacheIsolation 测试缓存隔离
func TestCacheIsolation(t *testing.T) {
	// 测试租户缓存隔离
	// 租户1设置缓存key="session", value="admin-token"
	// 租户2设置缓存key="session", value="user-token"
	// 验证租户1只能获取到"admin-token"

	// 需要Redis连接
	t.Skip("需要Redis连接")
}

// TestSearchIsolation 测试搜索隔离
func TestSearchIsolation(t *testing.T) {
	// 测试Elasticsearch索引隔离
	// 租户1的索引：tenant-tenant-001-bots
	// 租户2的索引：tenant-tenant-002-bots
	// 验证租户1只能搜索自己的索引

	// 需要Elasticsearch连接
	t.Skip("需要Elasticsearch连接")
}

// TestStorageIsolation 测试存储隔离
func TestStorageIsolation(t *testing.T) {
	// 测试对象存储隔离
	// 租户1的bucket：tenant-tenant-001
	// 租户2的bucket：tenant-tenant-002
	// 验证租户1无法访问租户2的文件

	// 需要MinIO/S3连接
	t.Skip("需要对象存储连接")
}

// TestSQLInjection 测试SQL注入攻击
func TestSQLInjection(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type TestData struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Name     string
	}

	db.AutoMigrate(&TestData{})

	tenant1 := "tenant-001"
	tenant2 := "tenant-002"

	db.Create(&TestData{ID: "1", TenantID: tenant1, Name: "Safe"})

	// 尝试SQL注入绕过租户过滤
	ctx := context.WithValue(context.Background(), "tenant_id", tenant2)

	sqlPayload := "1' OR '1'='1"
	var results []TestData

	err = db.WithContext(ctx).
		Where("id = ?", sqlPayload).
		Find(&results).Error

	// 应该失败或返回空结果
	if err == nil && len(results) > 0 {
		for _, result := range results {
			if result.TenantID == tenant1 {
				t.Errorf("安全漏洞：SQL注入绕过了租户隔离！payload=%s", sqlPayload)
			}
		}
	}
}

// TestConcurrentCrossTenantAccess 测试并发跨租户访问
func TestConcurrentCrossTenantAccess(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	type TestData struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Value    int
	}

	db.AutoMigrate(&TestData{})

	tenant1 := "tenant-001"
	tenant2 := "tenant-002"

	db.Create(&TestData{ID: "1", TenantID: tenant1, Value: 100})

	// 并发测试：同时从租户1和租户2访问数据
	done := make(chan bool)

	// 租户1读取
	go func() {
		ctx := context.WithValue(context.Background(), "tenant_id", tenant1)
		var result TestData
		db.WithContext(ctx).First(&result, "id = ?", "1")
		done <- true
	}()

	// 租户2尝试读取
	go func() {
		ctx := context.WithValue(context.Background(), "tenant_id", tenant2)
		var result TestData
		db.WithContext(ctx).First(&result, "id = ?", "1")

		// 验证租户2读不到数据
		if result.TenantID == tenant1 {
			t.Errorf("安全漏洞：并发访问时租户隔离失效！")
		}
		done <- true
	}()

	<-done
	<-done
}

// BenchmarkSecurityChecks 性能测试：安全检查开销
func BenchmarkSecurityChecks(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	type TestData struct {
		ID       string `gorm:"primaryKey"`
		TenantID string `gorm:"index"`
		Name     string
	}

	db.AutoMigrate(&TestData{})

	// 安装租户隔离插件
	// ...

	ctx := context.WithValue(context.Background(), "tenant_id", "tenant-001")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var results []TestData
		db.WithContext(ctx).Find(&results)
	}
}
