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
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/repository"
)

// ============================================================
// 共享Mock实现（用于单元测试）
// ============================================================

// MockTokenUsageLogRepository Mock日志仓储
type MockTokenUsageLogRepository struct {
	mock.Mock
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

func (m *MockTokenUsageLogRepository) GetByDateRange(ctx context.Context, tenantID string, startDate, endDate string) ([]*entity.TokenUsageLog, error) {
	args := m.Called(ctx, tenantID, startDate, endDate)
	if args.Get(0) == nil {
		return []*entity.TokenUsageLog{}, args.Error(1)
	}
	return args.Get(0).([]*entity.TokenUsageLog), args.Error(1)
}

func (m *MockTokenUsageLogRepository) GetByBotAndDateRange(ctx context.Context, tenantID, botID string, startDate, endDate string) ([]*entity.TokenUsageLog, error) {
	args := m.Called(ctx, tenantID, botID, startDate, endDate)
	if args.Get(0) == nil {
		return []*entity.TokenUsageLog{}, args.Error(1)
	}
	return args.Get(0).([]*entity.TokenUsageLog), args.Error(1)
}

func (m *MockTokenUsageLogRepository) List(ctx context.Context, filter *repository.TokenUsageLogFilter) ([]*entity.TokenUsageLog, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return []*entity.TokenUsageLog{}, int64(0), args.Error(2)
	}
	return args.Get(0).([]*entity.TokenUsageLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockTokenUsageLogRepository) BatchCreate(ctx context.Context, logs []*entity.TokenUsageLog) error {
	args := m.Called(ctx, logs)
	return args.Error(0)
}

// MockTokenUsageSummaryRepository Mock汇总仓储
type MockTokenUsageSummaryRepository struct {
	mock.Mock
}

func (m *MockTokenUsageSummaryRepository) Upsert(ctx context.Context, summary *entity.TokenUsageSummary) error {
	args := m.Called(ctx, summary)
	return args.Error(0)
}

func (m *MockTokenUsageSummaryRepository) GetByTenantAndDate(ctx context.Context, tenantID string, summaryDate string, summaryHour *uint8) (*entity.TokenUsageSummary, error) {
	args := m.Called(ctx, tenantID, summaryDate, summaryHour)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TokenUsageSummary), args.Error(1)
}

func (m *MockTokenUsageSummaryRepository) GetByDateRange(ctx context.Context, tenantID string, botID *string, startDate, endDate string) ([]*entity.TokenUsageSummary, error) {
	args := m.Called(ctx, tenantID, botID, startDate, endDate)
	if args.Get(0) == nil {
		return []*entity.TokenUsageSummary{}, args.Error(1)
	}
	return args.Get(0).([]*entity.TokenUsageSummary), args.Error(1)
}

func (m *MockTokenUsageSummaryRepository) GetTotalCost(ctx context.Context, tenantID string, startDate, endDate time.Time) (float64, error) {
	args := m.Called(ctx, tenantID, startDate, endDate)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTokenUsageSummaryRepository) List(ctx context.Context, filter *repository.TokenUsageSummaryFilter) ([]*entity.TokenUsageSummary, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return []*entity.TokenUsageSummary{}, int64(0), args.Error(2)
	}
	return args.Get(0).([]*entity.TokenUsageSummary), args.Get(1).(int64), args.Error(2)
}

// MockBudgetSettingsRepository Mock预算配置仓储
type MockBudgetSettingsRepository struct {
	mock.Mock
}

func (m *MockBudgetSettingsRepository) Create(ctx context.Context, budget *entity.BudgetSettings) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockBudgetSettingsRepository) GetByTenantID(ctx context.Context, tenantID string) (*entity.BudgetSettings, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BudgetSettings), args.Error(1)
}

func (m *MockBudgetSettingsRepository) Update(ctx context.Context, budget *entity.BudgetSettings) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockBudgetSettingsRepository) Delete(ctx context.Context, tenantID string) error {
	args := m.Called(ctx, tenantID)
	return args.Error(0)
}

func (m *MockBudgetSettingsRepository) List(ctx context.Context, filter *repository.BudgetSettingsFilter) ([]*entity.BudgetSettings, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return []*entity.BudgetSettings{}, int64(0), args.Error(2)
	}
	return args.Get(0).([]*entity.BudgetSettings), args.Get(1).(int64), args.Error(2)
}

// MockBudgetAlertRepository Mock告警仓储
type MockBudgetAlertRepository struct {
	mock.Mock
}

func (m *MockBudgetAlertRepository) Create(ctx context.Context, alert *entity.BudgetAlert) error {
	args := m.Called(ctx, alert)
	return args.Error(0)
}

func (m *MockBudgetAlertRepository) GetByID(ctx context.Context, id uint64) (*entity.BudgetAlert, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BudgetAlert), args.Error(1)
}

func (m *MockBudgetAlertRepository) GetByTenant(ctx context.Context, tenantID string, limit int) ([]*entity.BudgetAlert, error) {
	args := m.Called(ctx, tenantID, limit)
	if args.Get(0) == nil {
		return []*entity.BudgetAlert{}, args.Error(1)
	}
	return args.Get(0).([]*entity.BudgetAlert), args.Error(1)
}

func (m *MockBudgetAlertRepository) GetByTenantAndType(ctx context.Context, tenantID string, alertType string, limit int) ([]*entity.BudgetAlert, error) {
	args := m.Called(ctx, tenantID, alertType, limit)
	if args.Get(0) == nil {
		return []*entity.BudgetAlert{}, args.Error(1)
	}
	return args.Get(0).([]*entity.BudgetAlert), args.Error(1)
}

func (m *MockBudgetAlertRepository) UpdateNotificationStatus(ctx context.Context, id uint64, status string, sentAt *time.Time, errorMessage *string) error {
	args := m.Called(ctx, id, status, sentAt, errorMessage)
	return args.Error(0)
}

func (m *MockBudgetAlertRepository) CheckAlertSentToday(ctx context.Context, tenantID string, alertType string, today string) (bool, error) {
	args := m.Called(ctx, tenantID, alertType, today)
	return args.Bool(0), args.Error(1)
}

func (m *MockBudgetAlertRepository) List(ctx context.Context, filter *repository.BudgetAlertFilter) ([]*entity.BudgetAlert, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return []*entity.BudgetAlert{}, int64(0), args.Error(2)
	}
	return args.Get(0).([]*entity.BudgetAlert), args.Get(1).(int64), args.Error(2)
}

// ============================================================
// 辅助函数
// ============================================================

func uint64Ptr(v uint64) *uint64 {
	return &v
}

func strPtr(v string) *string {
	return &v
}

func intPtr(v int) *int {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}
