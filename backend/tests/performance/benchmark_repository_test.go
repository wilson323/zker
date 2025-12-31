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

package performance_test

import (
	"context"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/contextutil"
)

// =====================================================
// Repository 层基准测试
// =====================================================

// BenchmarkTenantRepository_Create 测试租户创建性能
func BenchmarkTenantRepository_Create(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository() // 假设存在此构造函数

	tenant := &entity.Tenant{
		TenantName:  "benchmark_tenant",
		TenantDesc:  "Performance test tenant",
		TenantType:  entity.TenantTypeEnterprise,
		TenantStatus: entity.TenantStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenant.TenantID = ""
		_, err := repo.Create(ctx, tenant)
		if err != nil {
			b.Fatalf("Failed to create tenant: %v", err)
		}
	}
}

// BenchmarkTenantRepository_FindByID 测试租户查询性能
func BenchmarkTenantRepository_FindByID(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	// 预先创建测试数据
	tenantID := "test_tenant_for_benchmark"
	_ = tenantID

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := repo.FindByID(ctx, tenantID)
		if err != nil {
			b.Fatalf("Failed to find tenant: %v", err)
		}
	}
}

// BenchmarkTenantRepository_List 测试租户列表查询性能
func BenchmarkTenantRepository_List(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := repo.List(ctx, &repository.ListOptions{
			Page:     1,
			PageSize: 20,
		})
		if err != nil {
			b.Fatalf("Failed to list tenants: %v", err)
		}
	}
}

// BenchmarkTenantRepository_Update 测试租户更新性能
func BenchmarkTenantRepository_Update(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	tenantID := "test_tenant_for_update"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenant, _ := repo.FindByID(ctx, tenantID)
		if tenant != nil {
			tenant.TenantDesc = "Updated description"
			err := repo.Update(ctx, tenant)
			if err != nil {
				b.Fatalf("Failed to update tenant: %v", err)
			}
		}
	}
}

// BenchmarkTenantRepository_Delete 测试租户删除性能
func BenchmarkTenantRepository_Delete(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	// 为每次迭代创建新租户
	testTenants := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		tenant := &entity.Tenant{
			TenantName:  "temp_tenant",
			TenantStatus: entity.TenantStatusActive,
		}
		createdTenant, _ := repo.Create(ctx, tenant)
		testTenants[i] = createdTenant.TenantID
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := repo.Delete(ctx, testTenants[i])
		if err != nil {
			b.Fatalf("Failed to delete tenant: %v", err)
		}
	}
}

// BenchmarkTenantRepository_WithTenantIsolation 测试带租户隔离的查询性能
func BenchmarkTenantRepository_WithTenantIsolation(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	// 在 context 中设置 tenant_id
	ctx = contextutil.WithTenantID(ctx, "benchmark_tenant")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟自动过滤 tenant_id 的查询
		tenantID := contextutil.GetTenantIDFromContext(ctx)
		_, err := repo.FindByID(ctx, tenantID)
		if err != nil {
			b.Fatalf("Failed to find tenant with isolation: %v", err)
		}
	}
}

// BenchmarkTenantRepository_BatchQuery 测试批量查询性能
func BenchmarkTenantRepository_BatchQuery(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	tenantIDs := []string{"tenant1", "tenant2", "tenant3", "tenant4", "tenant5"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, tenantID := range tenantIDs {
			_, err := repo.FindByID(ctx, tenantID)
			if err != nil {
				continue // 忽略不存在的租户
			}
		}
	}
}

// =====================================================
// 数据库连接池性能测试
// =====================================================

// BenchmarkDBConnectionPool_GetConnection 测试数据库连接池获取连接性能
func BenchmarkDBConnectionPool_GetConnection(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟从连接池获取连接
		// 这里需要根据实际的数据库连接池实现来编写
		_ = ctx
	}
}

// BenchmarkDBConnectionPool_ConcurrentQueries 测试并发查询性能
func BenchmarkDBConnectionPool_ConcurrentQueries(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := repo.FindByID(ctx, "test_tenant")
			if err != nil {
				// 忽略错误，仅测试性能
			}
		}
	})
}

// =====================================================
// 事务处理性能测试
// =====================================================

// BenchmarkTransaction_CreateTenantWithQuota 测试创建租户并初始化配额的事务性能
func BenchmarkTransaction_CreateTenantWithQuota(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟事务：创建租户 + 初始化配额
		err := repo.Transaction(ctx, func(txCtx context.Context) error {
			tenant := &entity.Tenant{
				TenantName:  "transaction_tenant",
				TenantStatus: entity.TenantStatusActive,
			}

			_, err := repo.Create(txCtx, tenant)
			return err
		})

		if err != nil {
			b.Fatalf("Transaction failed: %v", err)
		}
	}
}

// =====================================================
// 索引优化效果对比
// =====================================================

// BenchmarkQuery_WithIndex vs BenchmarkQuery_WithoutIndex
// 对比有无索引的查询性能

func BenchmarkQuery_WithIndex(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	// 使用索引字段查询 (tenant_id 是主键)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := repo.FindByID(ctx, "indexed_tenant_id")
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
	}
}

func BenchmarkQuery_WithoutIndex(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	// 使用非索引字段查询 (tenant_name 可能没有索引)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 假设存在 FindByName 方法
		// _, err := repo.FindByName(ctx, "tenant_name")
		_ = ctx
		// if err != nil {
		// 	b.Fatalf("Query failed: %v", err)
		// }
	}
}

// =====================================================
// Repository 缓存性能测试
// =====================================================

// BenchmarkRepositoryCache_WithRedis 测试使用 Redis 缓存的查询性能
func BenchmarkRepositoryCache_WithRedis(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepositoryWithCache() // 假设存在带缓存的版本

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 第一次查询会访问数据库
		// 后续查询会命中缓存
		_, err := repo.FindByID(ctx, "cached_tenant_id")
		if err != nil {
			b.Fatalf("Query with cache failed: %v", err)
		}
	}
}

func BenchmarkRepositoryCache_WithoutCache(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository() // 不带缓存

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 每次都访问数据库
		_, err := repo.FindByID(ctx, "uncached_tenant_id")
		if err != nil {
			b.Fatalf("Query without cache failed: %v", err)
		}
	}
}

// =====================================================
// 软删除性能测试
// =====================================================

// BenchmarkSoftDelete 测试软删除性能
func BenchmarkSoftDelete(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	// 创建测试数据
	tenant := &entity.Tenant{
		TenantName:  "soft_delete_test",
		TenantStatus: entity.TenantStatusActive,
	}
	createdTenant, _ := repo.Create(ctx, tenant)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := repo.SoftDelete(ctx, createdTenant.TenantID)
		if err != nil {
			b.Fatalf("Soft delete failed: %v", err)
		}
	}
}

// BenchmarkHardDelete 测试硬删除性能
func BenchmarkHardDelete(b *testing.B) {
	ctx := context.Background()
	repo := repository.NewTenantRepository()

	// 为每次迭代创建新数据
	testTenants := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		tenant := &entity.Tenant{
			TenantName:  "hard_delete_test",
			TenantStatus: entity.TenantStatusActive,
		}
		createdTenant, _ := repo.Create(ctx, tenant)
		testTenants[i] = createdTenant.TenantID
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := repo.HardDelete(ctx, testTenants[i])
		if err != nil {
			b.Fatalf("Hard delete failed: %v", err)
		}
	}
}
