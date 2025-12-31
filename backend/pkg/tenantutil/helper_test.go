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

package tenantutil

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/api/middleware"
	tenantentity "github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
)

type TestModel struct {
	ID       uint   `gorm:"primarykey"`
	TenantID string `gorm:"index"`
	Name     string
}

func TestWithTenantFilter(t *testing.T) {
	// 设置内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// 创建测试表
	db.AutoMigrate(&TestModel{})

	// 插入测试数据
	testTenantID := "tenant-123"
	db.Create(&TestModel{TenantID: testTenantID, Name: "Test1"})
	db.Create(&TestModel{TenantID: "other-tenant", Name: "Test2"})

	tests := []struct {
		name          string
		setupCtx      func() context.Context
		expectedCount int
		shouldError   bool
	}{
		{
			name: "valid tenant from string",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), "tenant_id", testTenantID)
			},
			expectedCount: 1,
			shouldError:   false,
		},
		{
			name: "valid tenant from entity",
			setupCtx: func() context.Context {
				tenant := &tenantentity.Tenant{TenantID: testTenantID}
				return context.WithValue(context.Background(), middleware.TenantContextKey, tenant)
			},
			expectedCount: 1,
			shouldError:   false,
		},
		{
			name: "empty tenant returns empty result",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), "tenant_id", "")
			},
			expectedCount: 0,
			shouldError:   false,
		},
		{
			name: "missing tenant returns empty result",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expectedCount: 0,
			shouldError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()

			var results []TestModel
			query := WithTenantFilter(db, ctx)
			err := query.Find(&results).Error

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(results) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(results))
			}

			if len(results) > 0 && results[0].TenantID != testTenantID {
				t.Errorf("expected tenant_id %s, got %s", testTenantID, results[0].TenantID)
			}
		})
	}
}

func TestWithTenantFilterAndDeleted(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	type SoftDeleteModel struct {
		ID        uint   `gorm:"primarykey"`
		TenantID  string `gorm:"index"`
		Name      string
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	db.AutoMigrate(&SoftDeleteModel{})

	testTenantID := "tenant-123"
	db.Create(&SoftDeleteModel{TenantID: testTenantID, Name: "Active"})
	db.Create(&SoftDeleteModel{TenantID: testTenantID, Name: "Deleted"}).Delete(&SoftDeleteModel{})

	ctx := context.WithValue(context.Background(), "tenant_id", testTenantID)

	var results []SoftDeleteModel
	query := WithTenantFilterAndDeleted(db, ctx)
	err = query.Find(&results).Error

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result (active only), got %d", len(results))
	}

	if len(results) > 0 && results[0].Name != "Active" {
		t.Errorf("expected Active record, got %s", results[0].Name)
	}
}

func TestGetTenantID(t *testing.T) {
	tests := []struct {
		name        string
		setupCtx    func() context.Context
		expectError bool
		expectedID  string
	}{
		{
			name: "get from string value",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), "tenant_id", "tenant-123")
			},
			expectError: false,
			expectedID:  "tenant-123",
		},
		{
			name: "get from entity value",
			setupCtx: func() context.Context {
				tenant := &tenantentity.Tenant{TenantID: "tenant-456"}
				return context.WithValue(context.Background(), middleware.TenantContextKey, tenant)
			},
			expectError: false,
			expectedID:  "tenant-456",
		},
		{
			name:        "missing tenant_id",
			setupCtx:    func() context.Context { return context.Background() },
			expectError: true,
			expectedID:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			tenantID, err := GetTenantID(ctx)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if err != ErrTenantNotFound {
					t.Errorf("expected ErrTenantNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tenantID != tt.expectedID {
					t.Errorf("expected tenant_id %s, got %s", tt.expectedID, tenantID)
				}
			}
		})
	}
}

func TestCheckTenantAccess(t *testing.T) {
	tests := []struct {
		name             string
		ctxTenantID      string
		resourceTenantID string
		shouldError      bool
	}{
		{
			name:             "access granted",
			ctxTenantID:      "tenant-123",
			resourceTenantID: "tenant-123",
			shouldError:      false,
		},
		{
			name:             "access denied",
			ctxTenantID:      "tenant-123",
			resourceTenantID: "tenant-456",
			shouldError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "tenant_id", tt.ctxTenantID)
			err := CheckTenantAccess(ctx, tt.resourceTenantID)

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestIsTenantOwned(t *testing.T) {
	tests := []struct {
		name             string
		ctxTenantID      string
		resourceTenantID string
		expected         bool
	}{
		{
			name:             "owned",
			ctxTenantID:      "tenant-123",
			resourceTenantID: "tenant-123",
			expected:         true,
		},
		{
			name:             "not owned",
			ctxTenantID:      "tenant-123",
			resourceTenantID: "tenant-456",
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "tenant_id", tt.ctxTenantID)
			result := IsTenantOwned(ctx, tt.resourceTenantID)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestWithTenantFilterByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	db.AutoMigrate(&TestModel{})

	testTenantID := "tenant-123"
	db.Create(&TestModel{TenantID: testTenantID, Name: "Test1"})
	db.Create(&TestModel{TenantID: "other-tenant", Name: "Test2"})

	var results []TestModel
	query := WithTenantFilterByID(db, testTenantID)
	err = query.Find(&results).Error

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	if results[0].TenantID != testTenantID {
		t.Errorf("expected tenant_id %s, got %s", testTenantID, results[0].TenantID)
	}
}

func TestMustGetTenantID(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic but didn't get one")
		}
	}()

	// 应该panic
	ctx := context.Background()
	_ = MustGetTenantID(ctx)
}

// 基准测试
func BenchmarkWithTenantFilter(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&TestModel{})

	ctx := context.WithValue(context.Background(), "tenant_id", "tenant-123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := WithTenantFilter(db, ctx)
		_ = query.Find(&[]TestModel{})
	}
}
