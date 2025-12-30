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

func TestSpaceTenantID(t *testing.T) {
	tests := []struct {
		name     string
		space    *Space
		expected string
		hasValue bool
	}{
		{
			name: "Space with valid tenant_id",
			space: &Space{
				ID:        1,
				TenantID:  "tenant-001",
				Name:      "Test Space",
				SpaceType: SpaceTypePersonal,
			},
			expected: "tenant-001",
			hasValue: true,
		},
		{
			name: "Space with empty tenant_id",
			space: &Space{
				ID:        2,
				TenantID:  "",
				Name:      "Test Space 2",
				SpaceType: SpaceTypeTeam,
			},
			expected: "",
			hasValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.space.GetTenantID())
			assert.Equal(t, tt.hasValue, tt.space.HasTenantID())
		})
	}
}

func TestSpaceGetTenantID(t *testing.T) {
	space := &Space{
		ID:        1,
		TenantID:  "tenant-123",
		Name:      "My Space",
		SpaceType: SpaceTypePersonal,
	}

	assert.Equal(t, "tenant-123", space.GetTenantID())
}

func TestSpaceHasTenantID(t *testing.T) {
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
			space := &Space{
				ID:       1,
				TenantID: tt.tenantID,
			}
			assert.Equal(t, tt.expected, space.HasTenantID())
		})
	}
}

func TestSessionTenantID(t *testing.T) {
	tests := []struct {
		name     string
		session  *Session
		expected string
		hasValue bool
	}{
		{
			name: "Session with valid tenant_id",
			session: &Session{
				UserID:    1,
				TenantID:  "tenant-001",
				UserEmail: "test@example.com",
			},
			expected: "tenant-001",
			hasValue: true,
		},
		{
			name: "Session with empty tenant_id",
			session: &Session{
				UserID:    2,
				TenantID:  "",
				UserEmail: "test2@example.com",
			},
			expected: "",
			hasValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.session.GetTenantID())
			assert.Equal(t, tt.hasValue, tt.session.HasTenantID())
		})
	}
}

func TestSessionGetTenantID(t *testing.T) {
	session := &Session{
		UserID:    1,
		TenantID:  "tenant-123",
		UserEmail: "john@example.com",
	}

	assert.Equal(t, "tenant-123", session.GetTenantID())
}

func TestSessionHasTenantID(t *testing.T) {
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
			session := &Session{
				UserID:   1,
				TenantID: tt.tenantID,
			}
			assert.Equal(t, tt.expected, session.HasTenantID())
		})
	}
}
