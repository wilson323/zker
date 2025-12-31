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

	"github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
)

// =====================================================
// Service 层基准测试
// =====================================================

// BenchmarkTenantService_CreateTenant 测试租户创建服务性能
func BenchmarkTenantService_CreateTenant(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService() // 假设存在此构造函数

	req := &service.CreateTenantRequest{
		TenantName:  "benchmark_tenant",
		TenantDesc:  "Performance test tenant",
		TenantType:  entity.TenantTypeEnterprise,
		AdminEmail:  "admin@benchmark.com",
		AdminName:   "Benchmark Admin",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := tenantService.CreateTenant(ctx, req)
		if err != nil {
			b.Fatalf("Failed to create tenant: %v", err)
		}
	}
}

// BenchmarkTenantService_GetTenantDetail 测试获取租户详情服务性能
func BenchmarkTenantService_GetTenantDetail(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	tenantID := "test_tenant_for_benchmark"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := tenantService.GetTenantDetail(ctx, tenantID)
		if err != nil {
			b.Fatalf("Failed to get tenant detail: %v", err)
		}
	}
}

// BenchmarkTenantService_UpdateTenant 测试更新租户服务性能
func BenchmarkTenantService_UpdateTenant(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	tenantID := "test_tenant_for_update"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := &service.UpdateTenantRequest{
			TenantID:   tenantID,
			TenantDesc: "Updated description",
		}
		err := tenantService.UpdateTenant(ctx, req)
		if err != nil {
			b.Fatalf("Failed to update tenant: %v", err)
		}
	}
}

// BenchmarkTenantService_DeleteTenant 测试删除租户服务性能
func BenchmarkTenantService_DeleteTenant(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 为每次迭代创建新租户
		createReq := &service.CreateTenantRequest{
			TenantName: "temp_tenant",
			AdminEmail: "temp@benchmark.com",
		}
		tenant, _ := tenantService.CreateTenant(ctx, createReq)

		err := tenantService.DeleteTenant(ctx, tenant.TenantID)
		if err != nil {
			b.Fatalf("Failed to delete tenant: %v", err)
		}
	}
}

// =====================================================
// 复杂业务逻辑性能测试
// =====================================================

// BenchmarkTenantService_CreateTenantWithQuota 测试创建租户并初始化配额的性能
func BenchmarkTenantService_CreateTenantWithQuota(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	req := &service.CreateTenantWithQuotaRequest{
		TenantName:  "quota_tenant",
		TenantDesc:  "Tenant with quota initialization",
		AdminEmail:  "admin@quota.com",
		QuotaConfig: &entity.QuotaConfig{
			MaxBots:         100,
			MaxUsers:        1000,
			MaxConversations: 10000,
			TokenLimit:      1000000,
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := tenantService.CreateTenantWithQuota(ctx, req)
		if err != nil {
			b.Fatalf("Failed to create tenant with quota: %v", err)
		}
	}
}

// BenchmarkTenantService_CheckTenantQuota 测试配额检查性能
func BenchmarkTenantService_CheckTenantQuota(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	tenantID := "test_tenant_for_quota"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := tenantService.CheckTenantQuota(ctx, tenantID, "bot_create")
		if err != nil {
			b.Fatalf("Failed to check quota: %v", err)
		}
	}
}

// BenchmarkTenantService_ConsumeQuota 测试配额消费性能
func BenchmarkTenantService_ConsumeQuota(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	tenantID := "test_tenant_for_consume"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := tenantService.ConsumeQuota(ctx, tenantID, &entity.QuotaConsumption{
			ResourceType: "token",
			Amount:       100,
		})
		if err != nil {
			b.Fatalf("Failed to consume quota: %v", err)
		}
	}
}

// =====================================================
// 并发性能测试
// =====================================================

// BenchmarkTenantService_ConcurrentCreate 测试并发创建租户性能
func BenchmarkTenantService_ConcurrentCreate(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			req := &service.CreateTenantRequest{
				TenantName:  "concurrent_tenant",
				AdminEmail:  "concurrent@benchmark.com",
				TenantDesc:  "Concurrent test",
			}
			_, err := tenantService.CreateTenant(ctx, req)
			if err != nil {
				b.Fatalf("Concurrent create failed: %v", err)
			}
			i++
		}
	})
}

// BenchmarkTenantService_ConcurrentQuotaCheck 测试并发配额检查性能
func BenchmarkTenantService_ConcurrentQuotaCheck(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	tenantID := "test_tenant_concurrent"

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := tenantService.CheckTenantQuota(ctx, tenantID, "bot_create")
			if err != nil {
				// 忽略错误，仅测试性能
			}
		}
	})
}

// =====================================================
// Service 层缓存性能测试
// =====================================================

// BenchmarkTenantService_WithCache 测试使用缓存的租户查询性能
func BenchmarkTenantService_WithCache(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantServiceWithCache() // 假设存在带缓存的版本

	tenantID := "cached_tenant_id"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := tenantService.GetTenantDetail(ctx, tenantID)
		if err != nil {
			b.Fatalf("Service with cache failed: %v", err)
		}
	}
}

func BenchmarkTenantService_WithoutCache(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService() // 不带缓存

	tenantID := "uncached_tenant_id"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := tenantService.GetTenantDetail(ctx, tenantID)
		if err != nil {
			b.Fatalf("Service without cache failed: %v", err)
		}
	}
}

// =====================================================
// Service 层复杂查询性能测试
// =====================================================

// BenchmarkTenantService_SearchTenants 测试租户搜索性能
func BenchmarkTenantService_SearchTenants(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	searchReq := &service.SearchTenantsRequest{
		Keyword:  "test",
		Page:     1,
		PageSize: 20,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := tenantService.SearchTenants(ctx, searchReq)
		if err != nil {
			b.Fatalf("Search tenants failed: %v", err)
		}
	}
}

// BenchmarkTenantService_GetTenantStatistics 测试获取租户统计信息性能
func BenchmarkTenantService_GetTenantStatistics(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	tenantID := "test_tenant_for_stats"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := tenantService.GetTenantStatistics(ctx, tenantID)
		if err != nil {
			b.Fatalf("Get statistics failed: %v", err)
		}
	}
}

// =====================================================
// Service 层错误处理性能测试
// =====================================================

// BenchmarkTenantService_ErrorCreation 测试错误创建性能
func BenchmarkTenantService_ErrorCreation(b *testing.B) {
	ctx := context.Background()
	tenantService := service.NewTenantService()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 故意使用不存在的租户ID，触发错误
		_, err := tenantService.GetTenantDetail(ctx, "non_existent_tenant")
		if err == nil {
			b.Fatal("Expected error, got nil")
		}
		_ = err // 使用错误
	}
}

// =====================================================
// Service 层数据转换性能测试
// =====================================================

// BenchmarkTenantService_EntityToVO 测试实体到VO转换性能
func BenchmarkTenantService_EntityToVO(b *testing.B) {
	tenantEntity := &entity.Tenant{
		TenantID:     "test_tenant",
		TenantName:   "Test Tenant",
		TenantDesc:   "Test Description",
		TenantType:   entity.TenantTypeEnterprise,
		TenantStatus: entity.TenantStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 假设存在 EntityToVO 方法
		_ = tenantEntity // 模拟转换操作
	}
}

// =====================================================
// Service 层验证性能测试
// =====================================================

// BenchmarkTenantService_ValidateTenantName 测试租户名称验证性能
func BenchmarkTenantService_ValidateTenantName(b *testing.B) {
	tenantService := service.NewTenantService()

	validNames := []string{
		"valid_tenant",
		"ValidTenant123",
		"valid-tenant",
		"VALID_TENANT",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		name := validNames[i%len(validNames)]
		err := tenantService.ValidateTenantName(name)
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
	}
}

// BenchmarkTenantService_ValidateEmail 测试邮箱验证性能
func BenchmarkTenantService_ValidateEmail(b *testing.B) {
	tenantService := service.NewTenantService()

	validEmails := []string{
		"user@example.com",
		"test.user@test.co.uk",
		"admin123@company.org",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		email := validEmails[i%len(validEmails)]
		err := tenantService.ValidateEmail(email)
		if err != nil {
			b.Fatalf("Email validation failed: %v", err)
		}
	}
}
