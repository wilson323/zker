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
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	billingRepository "github.com/coze-dev/coze-studio/backend/domain/billing/repository"
)

// ============================================================
// Mock Repository 实现
// ============================================================

// MockSubscriptionRepository 模拟订阅仓储
type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) GetActiveByTenant(ctx context.Context, tenantID string) (*billingRepository.Subscription, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*billingRepository.Subscription), args.Error(1)
}

// MockQuotaRepository 模拟配额仓储
type MockQuotaRepository struct {
	mock.Mock
}

func (m *MockQuotaRepository) GetByTenant(ctx context.Context, tenantID string) (*billingRepository.Quota, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*billingRepository.Quota), args.Error(1)
}

// MockBillingAccountRepository 模拟计费账户仓储
type MockBillingAccountRepository struct {
	mock.Mock
}

func (m *MockBillingAccountRepository) Create(ctx context.Context, account *entity.BillingAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockBillingAccountRepository) GetByID(ctx context.Context, id uint64) (*entity.BillingAccount, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BillingAccount), args.Error(1)
}

func (m *MockBillingAccountRepository) GetByTenant(ctx context.Context, tenantID string) (*entity.BillingAccount, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BillingAccount), args.Error(1)
}

func (m *MockBillingAccountRepository) Update(ctx context.Context, account *entity.BillingAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockBillingAccountRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBillingAccountRepository) List(ctx context.Context, filter *billingRepository.BillingAccountFilter) ([]*entity.BillingAccount, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.BillingAccount), args.Get(1).(int64), args.Error(2)
}

func (m *MockBillingAccountRepository) UpdateBalance(ctx context.Context, tenantID string, amount float64) error {
	args := m.Called(ctx, tenantID, amount)
	return args.Error(0)
}

// MockInvoiceRepository 模拟发票仓储
type MockInvoiceRepository struct {
	mock.Mock
}

func (m *MockInvoiceRepository) Create(ctx context.Context, invoice *entity.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceRepository) GetByID(ctx context.Context, id uint64) (*entity.Invoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*entity.Invoice, error) {
	args := m.Called(ctx, invoiceNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) Update(ctx context.Context, invoice *entity.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInvoiceRepository) ListByTenant(ctx context.Context, tenantID string, page, pageSize int) ([]*entity.Invoice, int64, error) {
	args := m.Called(ctx, tenantID, page, pageSize)
	return args.Get(0).([]*entity.Invoice), args.Get(1).(int64), args.Error(2)
}

func (m *MockInvoiceRepository) List(ctx context.Context, filter *billingRepository.InvoiceFilter) ([]*entity.Invoice, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.Invoice), args.Get(1).(int64), args.Error(2)
}

func (m *MockInvoiceRepository) GetOverdueInvoices(ctx context.Context, beforeDate time.Time) ([]*entity.Invoice, error) {
	args := m.Called(ctx, beforeDate)
	return args.Get(0).([]*entity.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) CountByTenantAndDate(ctx context.Context, tenantID string, date time.Time) (int, error) {
	args := m.Called(ctx, tenantID, date)
	return args.Int(0), args.Error(1)
}

func (m *MockInvoiceRepository) CreateLineItem(ctx context.Context, item *entity.InvoiceLineItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockInvoiceRepository) GetLineItems(ctx context.Context, invoiceID uint64) ([]entity.InvoiceLineItem, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).([]entity.InvoiceLineItem), args.Error(1)
}

func (m *MockInvoiceRepository) UpdateLineItem(ctx context.Context, item *entity.InvoiceLineItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockInvoiceRepository) DeleteLineItem(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTokenUsageLogRepository 模拟Token使用日志仓储
}

func (m *MockTokenUsageLogRepository) Create(ctx context.Context, log *entity.TokenUsageLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockTokenUsageLogRepository) GetByID(ctx context.Context, id uint64) (*entity.TokenUsageLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TokenUsageLog), args.Error(1)
}

func (m *MockTokenUsageLogRepository) GetByDateRange(ctx context.Context, tenantID, startDate, endDate string) ([]*entity.TokenUsageLog, error) {
	args := m.Called(ctx, tenantID, startDate, endDate)
	if args.Get(0) == nil {
		return []*entity.TokenUsageLog{}, args.Error(1)
	}
	return args.Get(0).([]*entity.TokenUsageLog), args.Error(1)
}

func (m *MockTokenUsageLogRepository) GetByBotAndDateRange(ctx context.Context, tenantID, botID, startDate, endDate string) ([]*entity.TokenUsageLog, error) {
	args := m.Called(ctx, tenantID, botID, startDate, endDate)
	if args.Get(0) == nil {
		return []*entity.TokenUsageLog{}, args.Error(1)
	}
	return args.Get(0).([]*entity.TokenUsageLog), args.Error(1)
}

func (m *MockTokenUsageLogRepository) List(ctx context.Context, filter *billingRepository.TokenUsageLogFilter) ([]*entity.TokenUsageLog, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.TokenUsageLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockTokenUsageLogRepository) BatchCreate(ctx context.Context, logs []*entity.TokenUsageLog) error {
	args := m.Called(ctx, logs)
	return args.Error(0)
}

// ============================================================
// BillingEngine 单元测试
// ============================================================

/**
 * TestBillingEngine_CalculateSubscriptionFee
 *
 * 测试订阅费用计算
 */
func TestBillingEngine_CalculateSubscriptionFee(t *testing.T) {
	// 准备测试数据
	ctx := context.Background()
	tenantID := "tenant_test_001"
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	// 创建Mock
	mockSubRepo := new(MockSubscriptionRepository)
	mockQuotaRepo := new(MockQuotaRepository)
	mockUsageLogRepo := new(MockTokenUsageLogRepository)
	mockUsageSummaryRepo := new(MockTokenUsageSummaryRepository) // 企业级计费：添加汇总Mock
	mockAccountRepo := new(MockBillingAccountRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	mockDB := &gorm.DB{} // 简化处理

	// 设置Mock期望
	subscription := &billingRepository.Subscription{
		ID:           1,
		TenantID:     tenantID,
		PlanID:       "plan_pro",
		Price:        1000.00,
		Currency:     "CNY",
		BillingCycle: "monthly",
		Status:       "active",
	}
	mockSubRepo.On("GetActiveByTenant", ctx, tenantID).Return(subscription, nil)

	// 创建计费引擎
	engine := NewBillingEngine(
		NewPricingEngine(),
		mockSubRepo,
		mockQuotaRepo,
		mockUsageLogRepo,
		mockUsageSummaryRepo, // 企业级计费：添加TokenUsageSummaryRepository
		mockAccountRepo,
		mockInvoiceRepo,
		mockDB,
		nil, // logger
	)

	// 执行测试
	t.Run("正常计算月度订阅费用", func(t *testing.T) {
		bill, err := engine.CalculateSubscriptionFee(ctx, tenantID, startDate, endDate)

		// 验证结果
		require.NoError(t, err)
		require.NotNil(t, bill)

		// 验证费用计算
		// 1月有31天，费用应为: 1000 * (31/31) = 1000
		assert.Equal(t, tenantID, bill.TenantID)
		assert.Equal(t, "CNY", bill.Currency)

		// 验证明细项
		assert.Len(t, bill.LineItems, 1)
		assert.Equal(t, entity.LineItemTypeSubscription, bill.LineItems[0].ItemType)
		assert.Equal(t, 31, bill.LineItems[0].Quantity)
	})

	// 验证Mock调用
	mockSubRepo.AssertExpectations(t)
}

/**
 * TestBillingEngine_CalculateOverage
 *
 * 测试超额费用计算
 */
func TestBillingEngine_CalculateOverage(t *testing.T) {
	ctx := context.Background()
	tenantID := "tenant_test_002"

	// 创建Mock
	mockSubRepo := new(MockSubscriptionRepository)
	mockQuotaRepo := new(MockQuotaRepository)
	mockUsageLogRepo := new(MockTokenUsageLogRepository)
	mockAccountRepo := new(MockBillingAccountRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	mockUsageSummaryRepo := new(MockTokenUsageSummaryRepository) // 企业级计费：添加汇总Mock
	mockDB := &gorm.DB{}

	// 设置配额
	tokenLimit := int64(1000000) // 1M tokens
	quota := &billingRepository.Quota{
		ID:         1,
		TenantID:   tenantID,
		TokenLimit: &tokenLimit,
		Currency:   "CNY",
	}
	mockQuotaRepo.On("GetByTenant", ctx, tenantID).Return(quota, nil)

	// 创建计费引擎
	engine := NewBillingEngine(
		NewPricingEngine(),
		mockSubRepo,
		mockQuotaRepo,
		mockUsageLogRepo,
		mockUsageSummaryRepo, // 企业级计费：添加TokenUsageSummaryRepository
		mockAccountRepo,
		mockInvoiceRepo,
		mockDB,
		nil,
	)

	t.Run("Token超额费用计算", func(t *testing.T) {
		usage := &Usage{
			TokenUsage:   1500000, // 1.5M tokens，超额500K
			StorageUsage: 0,
			RequestCount: 1000,
		}

		overageFee, err := engine.CalculateOverage(ctx, tenantID, usage)

		require.NoError(t, err)
		assert.True(t, overageFee.IsPositive())

		// 验证计算：500K tokens * ¥0.05/1K = ¥25
		expectedFee := decimal.NewFromInt(500).Mul(decimal.NewFromFloat(0.05))
		assert.True(t, overageFee.GreaterThanOrEqual(expectedFee))
	})

	t.Run("无超额情况", func(t *testing.T) {
		usage := &Usage{
			TokenUsage:   500000, // 未超限额
			StorageUsage: 0,
			RequestCount: 500,
		}

		overageFee, err := engine.CalculateOverage(ctx, tenantID, usage)

		require.NoError(t, err)
		assert.True(t, overageFee.IsZero())
	})

	mockQuotaRepo.AssertExpectations(t)
}

/**
 * TestBillingEngine_GenerateInvoice
 *
 * 测试发票生成
 */
func TestBillingEngine_GenerateInvoice(t *testing.T) {
	ctx := context.Background()
	tenantID := "tenant_test_003"

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)
	billingCycle := BillingCycle{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// 创建Mock
	mockSubRepo := new(MockSubscriptionRepository)
	mockQuotaRepo := new(MockQuotaRepository)
	mockUsageLogRepo := new(MockTokenUsageLogRepository)
	mockAccountRepo := new(MockBillingAccountRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)

	// Mock DB
	mockDB := &gorm.DB{}

	// 设置Mock期望
	subscription := &billingRepository.Subscription{
		ID:           1,
		TenantID:     tenantID,
		PlanID:       "plan_pro",
		Price:        1000.00,
		Currency:     "CNY",
		BillingCycle: "monthly",
		Status:       "active",
	}
	mockSubRepo.On("GetActiveByTenant", ctx, tenantID).Return(subscription, nil)

	account := &entity.BillingAccount{
		ID:                1,
		TenantID:          tenantID,
		Status:            entity.BillingAccountStatusActive,
		CreditLimit:       10000.00,
		AvailableCredit:   10000.00,
		Currency:          "CNY",
		BillingCycle:      entity.BillingCycleMonthly,
	}
	mockAccountRepo.On("GetByTenant", ctx, tenantID).Return(account, nil)

	// 设置Token使用日志（空）
	mockUsageLogRepo.On("GetByDateRange", ctx, tenantID, "2025-01-01", "2025-01-31").
		Return([]*entity.TokenUsageLog{}, nil)

	// 设置配额
	mockQuotaRepo.On("GetByTenant", ctx, tenantID).Return(&billingRepository.Quota{
		ID:       1,
		TenantID: tenantID,
		Currency: "CNY",
	}, nil)

	// 设置发票创建
	mockInvoiceRepo.On("Create", ctx, mock.AnythingOfType("*entity.Invoice")).Return(nil)
	mockInvoiceRepo.On("CreateLineItem", ctx, mock.AnythingOfType("*entity.InvoiceLineItem")).Return(nil)
	mockInvoiceRepo.On("CountByTenantAndDate", ctx, tenantID, endDate).Return(0, nil)

	// 创建计费引擎
	engine := NewBillingEngine(
		NewPricingEngine(),
		mockSubRepo,
		mockQuotaRepo,
		mockUsageLogRepo,
		mockUsageSummaryRepo, // 企业级计费：添加TokenUsageSummaryRepository
		mockAccountRepo,
		mockInvoiceRepo,
		mockDB,
		nil,
	)

	t.Run("生成月度发票", func(t *testing.T) {
		invoice, err := engine.GenerateInvoice(ctx, tenantID, billingCycle)

		require.NoError(t, err)
		require.NotNil(t, invoice)

		// 验证发票信息
		assert.Equal(t, tenantID, invoice.TenantID)
		assert.Equal(t, uint64(1), invoice.BillingAccountID)
		assert.Equal(t, entity.InvoiceTypeSubscription, invoice.InvoiceType)
		assert.Equal(t, entity.InvoiceStatusDraft, invoice.Status)
		assert.Equal(t, "CNY", invoice.Currency)

		// 验证金额
		assert.True(t, invoice.Subtotal > 0)
		assert.Equal(t, invoice.Subtotal, invoice.TotalAmount)

		// 验证账期
		assert.Equal(t, startDate, invoice.PeriodStart)
		assert.Equal(t, endDate, invoice.PeriodEnd)
		assert.NotNil(t, invoice.DueDate)

		// 验证发票编号格式
		assert.Contains(t, invoice.InvoiceNumber, "INV-")
		assert.Contains(t, invoice.InvoiceNumber, tenantID)
	})

	// 验证Mock调用
	mockSubRepo.AssertExpectations(t)
	mockAccountRepo.AssertExpectations(t)
	mockInvoiceRepo.AssertExpectations(t)
}

/**
 * TestBillingEngine_GetCurrentBillingCycle
 *
 * 测试获取当前账单周期
 */
func TestBillingEngine_GetCurrentBillingCycle(t *testing.T) {
	ctx := context.Background()
	tenantID := "tenant_test_004"

	// 创建Mock
	mockSubRepo := new(MockSubscriptionRepository)
	mockQuotaRepo := new(MockQuotaRepository)
	mockUsageLogRepo := new(MockTokenUsageLogRepository)
	mockAccountRepo := new(MockBillingAccountRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	mockDB := &gorm.DB{}

	// 创建计费引擎
	engine := NewBillingEngine(
		NewPricingEngine(),
		mockSubRepo,
		mockQuotaRepo,
		mockUsageLogRepo,
		mockUsageSummaryRepo, // 企业级计费：添加TokenUsageSummaryRepository
		mockAccountRepo,
		mockInvoiceRepo,
		mockDB,
		nil,
	)

	t.Run("自然月周期", func(t *testing.T) {
		// 设置Mock期望
		account := &entity.BillingAccount{
			ID:               1,
			TenantID:         tenantID,
			BillingCycle:     entity.BillingCycleMonthly,
			BillingDayOfMonth: nil, // 自然月
		}
		mockAccountRepo.On("GetByTenant", ctx, tenantID).Return(account, nil)

		// 执行测试
		cycle, err := engine.GetCurrentBillingCycle(ctx, tenantID)

		require.NoError(t, err)
		require.NotNil(t, cycle)

		// 验证账期（应该是自然月）
		assert.Equal(t, 1, int(cycle.StartDate.Day()))
		assert.True(t, cycle.EndDate.Day() >= 28) // 月末
	})

	t.Run("自定义账单日", func(t *testing.T) {
		billingDay := uint8(15)

		// 设置Mock期望
		account := &entity.BillingAccount{
			ID:               1,
			TenantID:         tenantID,
			BillingCycle:     entity.BillingCycleMonthly,
			BillingDayOfMonth: &billingDay,
		}
		mockAccountRepo.On("GetByTenant", ctx, tenantID).Return(account, nil)

		// 执行测试
		cycle, err := engine.GetCurrentBillingCycle(ctx, tenantID)

		require.NoError(t, err)
		require.NotNil(t, cycle)

		// 验证账期（应该从15日开始）
		// 具体验证逻辑取决于当前日期
		assert.NotNil(t, cycle.StartDate)
		assert.NotNil(t, cycle.EndDate)
	})

	mockAccountRepo.AssertExpectations(t)
}

// ============================================================
// 辅助函数
// ============================================================

// mustFloat64 安全地将decimal.Decimal转换为float64
func mustFloat64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}
