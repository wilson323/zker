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

package service

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
)

// TestGenerateID_Concurrency 测试ID生成的并发安全性
func TestGenerateID_Concurrency(t *testing.T) {
	t.Run("并发生成1000个ID无重复", func(t *testing.T) {
		const concurrency = 100
		const iterations = 10

		var wg sync.WaitGroup
		idMap := make(map[string]bool)
		idMutex := sync.Mutex{}

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(workerIndex int) {
				defer wg.Done()

				for j := 0; j < iterations; j++ {
					prefix := fmt.Sprintf("worker_%d", workerIndex)
					id := generateID(prefix)

					// 检查ID是否唯一
					idMutex.Lock()
					if idMap[id] {
						t.Errorf("Duplicate ID generated: %s", id)
					}
					idMap[id] = true
					idMutex.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// 验证生成的ID数量
		expectedCount := concurrency * iterations
		actualCount := len(idMap)
		assert.Equal(t, expectedCount, actualCount, "Should generate %d unique IDs", expectedCount)
	})

	t.Run("快速连续生成ID不重复", func(t *testing.T) {
		const count = 1000
		ids := make([]string, 0, count)
		idMap := make(map[string]bool)

		// 快速生成1000个ID
		for i := 0; i < count; i++ {
			id := generateID("test")
			ids = append(ids, id)

			if idMap[id] {
				t.Errorf("Duplicate ID at index %d: %s", i, id)
			}
			idMap[id] = true
		}

		assert.Equal(t, count, len(idMap), "All IDs should be unique")

		// 验证ID格式
		for _, id := range ids {
			assert.Contains(t, id, "test_", "ID should start with prefix")
			assert.Contains(t, id, "_", "ID should contain separators")
		}
	})
}

// TestTenantRegistrationService_ConcurrentRegistrations 测试并发注册
func TestTenantRegistrationService_ConcurrentRegistrations(t *testing.T) {
	// Skip测试：需要真实数据库环境
	t.Skip("Skipping concurrent registration test - requires database")

	/*
		db := setupTestDB(t)
		tenantRepo := repository.NewTenantRepositoryImpl(db)
		codeSvc := &mockVerificationCodeService{}
		subscriptionRepo := repository.NewSubscriptionRepositoryImpl(db)
		quotaRepo := repository.NewQuotaRepositoryImpl(db)

		svc := NewTenantRegistrationService(db, tenantRepo, codeSvc, subscriptionRepo, quotaRepo)

		const concurrency = 10

		var wg sync.WaitGroup
		errors := make([]error, concurrency)
		results := make([]*RegisterTenantResult, concurrency)

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				req := &RegisterTenantRequest{
					Email:            fmt.Sprintf("user%d@example.com", index),
					VerificationCode: "123456",
					CompanyName:      fmt.Sprintf("Company%d", index),
					Subdomain:        fmt.Sprintf("company%d", index),
					Industry:         "technology",
					AdminUsername:    fmt.Sprintf("admin%d", index),
					AdminPassword:    "SecurePassword123!",
					AdminEmail:       fmt.Sprintf("admin%d@example.com", index),
				}

				result, err := svc.RegisterTenant(context.Background(), req)
				errors[index] = err
				results[index] = result
			}(i)
		}

		wg.Wait()

		// 验证所有注册都成功
		for i, err := range errors {
			assert.NoError(t, err, "Registration %d should succeed", i)
			assert.NotNil(t, results[i], "Registration %d should return result", i)
		}
	*/
}

// TestTenantRegistrationService_TransactionRollback 测试事务回滚
func TestTenantRegistrationService_TransactionRollback(t *testing.T) {
	// Skip测试：需要真实数据库环境
	t.Skip("Skipping transaction rollback test - requires database")

	/*
		db := setupTestDB(t)
		tenantRepo := &mockTenantRepository{
			createFunc: func(ctx context.Context, tenant *entity.Tenant) error {
				return nil
			},
			getByNameFunc: func(ctx context.Context, name string) (*entity.Tenant, error) {
				return nil, repository.ErrRecordNotFound
			},
			getBySubdomainFunc: func(ctx context.Context, subdomain string) (*entity.Tenant, error) {
				return nil, repository.ErrRecordNotFound
			},
		}

		codeSvc := &mockVerificationCodeService{
			verifyCodeFunc: func(ctx context.Context, email, code string) (bool, error) {
				return true, nil
			},
		}

		subscriptionRepo := &mockSubscriptionRepository{
			createFunc: func(ctx context.Context, subscription *entity.Subscription) error {
				// 模拟订阅创建失败
				return fmt.Errorf("subscription creation failed")
			},
		}

		quotaRepo := &mockQuotaRepository{}

		svc := NewTenantRegistrationService(db, tenantRepo, codeSvc, subscriptionRepo, quotaRepo)

		req := &RegisterTenantRequest{
			Email:            "test@example.com",
			VerificationCode: "123456",
			CompanyName:      "TestCompany",
			Subdomain:        "testcompany",
			Industry:         "technology",
			AdminUsername:    "admin",
			AdminPassword:    "SecurePassword123!",
			AdminEmail:       "admin@example.com",
		}

		// 注册应该失败
		result, err := svc.RegisterTenant(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, result)

		// 验证租户未创建（事务已回滚）
		_, err = tenantRepo.GetByName(context.Background(), "TestCompany")
		assert.Equal(t, repository.ErrRecordNotFound, err)
	*/
}

// TestTenantValidationService_ConcurrentChecks 测试并发验证
func TestTenantValidationService_ConcurrentChecks(t *testing.T) {
	// Skip测试：需要真实数据库环境
	t.Skip("Skipping concurrent validation test - requires database")

	/*
		db := setupTestDB(t)
		tenantRepo := repository.NewTenantRepositoryImpl(db)
		svc := NewTenantValidationService(tenantRepo)

		const concurrency = 50

		var wg sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				// 检查不同的公司名称
				companyName := fmt.Sprintf("Company%d", index%10) // 10个不同的名称
				available, _, err := svc.CheckCompanyNameAvailability(context.Background(), companyName)

				assert.NoError(t, err)
				assert.True(t, available) // 应该都可用（数据库为空）
			}(i)
		}

		wg.Wait()
	*/
}

// TestTenantManagementService_ConcurrentQuotaUpdates 测试并发配额更新
func TestTenantManagementService_ConcurrentQuotaUpdates(t *testing.T) {
	// Skip测试：需要真实数据库环境
	t.Skip("Skipping concurrent quota update test - requires database")

	/*
		db := setupTestDB(t)
		tenantRepo := repository.NewTenantRepositoryImpl(db)
		quotaRepo := repository.NewQuotaRepositoryImpl(db)
		subscriptionRepo := repository.NewSubscriptionRepositoryImpl(db)

		svc := NewTenantManagementService(db, tenantRepo, quotaRepo, subscriptionRepo)

		// 创建测试租户和配额
		tenantID := "test_tenant_concurrent"
		quota := &entity.Quota{
			QuotaID:      generateID("quota"),
			TenantID:     tenantID,
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    0,
			MaxLimit:     100,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
			CreatedAt:    time.Now().UnixMilli(),
			UpdatedAt:    time.Now().UnixMilli(),
		}
		db.Create(quota)

		const concurrency = 10
		const incrementsPerWorker = 10

		var wg sync.WaitGroup
		errors := make([]error, concurrency)

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(workerIndex int) {
				defer wg.Done()

				for j := 0; j < incrementsPerWorker; j++ {
					// 每个worker增加配额10次，每次1个
					checker := NewQuotaChecker(db, quotaRepo)
					success, err := checker.CheckAndConsume(context.Background(), tenantID, entity.ResourceTypeBots, 1)
					if err != nil {
						errors[workerIndex] = err
						return
					}
					if !success {
						errors[workerIndex] = fmt.Errorf("quota check failed")
						return
					}
				}
			}(i)
		}

		wg.Wait()

		// 验证没有错误
		for i, err := range errors {
			assert.NoError(t, err, "Worker %d should succeed", i)
		}

		// 验证最终配额计数
		finalQuota, err := quotaRepo.GetByTenantAndResource(context.Background(), tenantID, entity.ResourceTypeBots)
		assert.NoError(t, err)

		expectedCount := concurrency * incrementsPerWorker
		assert.Equal(t, expectedCount, finalQuota.UsedCount, "Total quota count should be %d", expectedCount)
	*/
}

// BenchmarkGenerateID 性能基准测试
func BenchmarkGenerateID(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = generateID("bench")
		}
	})
}

// BenchmarkCheckCompanyNameAvailability 性能基准测试
func BenchmarkCheckCompanyNameAvailability(b *testing.B) {
	// Skip：需要真实数据库
	b.Skip("Skipping benchmark - requires database")

	/*
		db := setupTestDB(b)
		tenantRepo := repository.NewTenantRepositoryImpl(db)
		svc := NewTenantValidationService(tenantRepo)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = svc.CheckCompanyNameAvailability(context.Background(), fmt.Sprintf("Company%d", i))
		}
	*/
}

// Mock实现示例（用于单元测试）
type mockVerificationCodeService struct {
	verifyCodeFunc func(ctx context.Context, email, code string) (bool, error)
	sendCodeFunc   func(ctx context.Context, email, purpose string) error
}

func (m *mockVerificationCodeService) SendCode(ctx context.Context, email, purpose string) error {
	if m.sendCodeFunc != nil {
		return m.sendCodeFunc(ctx, email, purpose)
	}
	return nil
}

func (m *mockVerificationCodeService) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	if m.verifyCodeFunc != nil {
		return m.verifyCodeFunc(ctx, email, code)
	}
	return true, nil
}

type mockTenantRepository struct {
	repository.TenantRepository
	createFunc           func(ctx context.Context, tenant *entity.Tenant) error
	getByNameFunc        func(ctx context.Context, name string) (*entity.Tenant, error)
	getBySubdomainFunc   func(ctx context.Context, subdomain string) (*entity.Tenant, error)
}

func (m *mockTenantRepository) Create(ctx context.Context, tenant *entity.Tenant) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, tenant)
	}
	return nil
}

func (m *mockTenantRepository) GetByName(ctx context.Context, name string) (*entity.Tenant, error) {
	if m.getByNameFunc != nil {
		return m.getByNameFunc(ctx, name)
	}
	return nil, repository.ErrRecordNotFound
}

func (m *mockTenantRepository) GetBySubdomain(ctx context.Context, subdomain string) (*entity.Tenant, error) {
	if m.getBySubdomainFunc != nil {
		return m.getBySubdomainFunc(ctx, subdomain)
	}
	return nil, repository.ErrRecordNotFound
}

type mockSubscriptionRepository struct {
	repository.SubscriptionRepository
	createFunc func(ctx context.Context, subscription *entity.Subscription) error
}

func (m *mockSubscriptionRepository) Create(ctx context.Context, subscription *entity.Subscription) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, subscription)
	}
	return nil
}

type mockQuotaRepository struct {
	repository.QuotaRepository
}

// 辅助函数：创建测试数据库连接
func setupTestDB(t testing.TB) *gorm.DB {
	// 在实际实现中，这里应该：
	// 1. 使用内存SQLite数据库
	// 2. 或使用Docker启动临时MySQL容器
	// 3. 或使用testify/mock进行Mock
	return nil
}
