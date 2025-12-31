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

// Package mocks provides mock utilities for testing
package mocks

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/pkg/contextutil"
)

// MockContext creates a mock context with tenant and user information
func MockContext(tenantID string, userID string) context.Context {
	ctx := context.Background()
	ctx = contextutil.WithTenantID(ctx, tenantID)
	ctx = contextutil.WithUserID(ctx, userID)
	return ctx
}

// MockContextWithUserID creates a mock context with only user ID
func MockContextWithUserID(userID string) context.Context {
	ctx := context.Background()
	ctx = contextutil.WithUserID(ctx, userID)
	return ctx
}

// MockContextWithTenantID creates a mock context with only tenant ID
func MockContextWithTenantID(tenantID string) context.Context {
	ctx := context.Background()
	ctx = contextutil.WithTenantID(ctx, tenantID)
	return ctx
}

// MockContextWithRequestID creates a mock context with request ID
func MockContextWithRequestID(requestID string) context.Context {
	ctx := context.Background()
	ctx = contextutil.WithRequestID(ctx, requestID)
	return ctx
}
