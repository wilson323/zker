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

package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserTenantID(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		expected string
		hasValue bool
	}{
		{
			name: "User with valid tenant_id",
			user: &User{
				UserID:   1,
				TenantID: "tenant-001",
				Name:     "Test User",
				Email:    "test@example.com",
			},
			expected: "tenant-001",
			hasValue: true,
		},
		{
			name: "User with empty tenant_id",
			user: &User{
				UserID:   2,
				TenantID: "",
				Name:     "Test User 2",
				Email:    "test2@example.com",
			},
			expected: "",
			hasValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.user.GetTenantID())
			assert.Equal(t, tt.hasValue, tt.user.HasTenantID())
		})
	}
}

func TestUserGetTenantID(t *testing.T) {
	user := &User{
		UserID:   1,
		TenantID: "tenant-123",
		Name:     "John Doe",
		Email:    "john@example.com",
	}

	assert.Equal(t, "tenant-123", user.GetTenantID())
}

func TestUserHasTenantID(t *testing.T) {
	tests := []struct {
		name     string
		tenantID string
		expected bool
	}{
		{
			name:     "Has tenant ID",
			tenantID: "tenant-001",
			expected: true,
		},
		{
			name:     "Empty tenant ID",
			tenantID: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				UserID:   1,
				TenantID: tt.tenantID,
			}
			assert.Equal(t, tt.expected, user.HasTenantID())
		})
	}
}
