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

// Package contextutil provides context utility functions that are independent
// and can be safely used by any layer without creating circular dependencies.
package contextutil

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

const (
	// TenantIDKey is the key for tenant ID in context
	TenantIDKey = "tenant_id"

	// DefaultTenantID is the default tenant ID when no tenant ID is found in context
	DefaultTenantID = "default"
)

// GetTenantIDFromContext retrieves tenant ID from context
// This function is moved to pkg/contextutil to break circular dependency:
// - application/user → bizpkg/config → api/middleware → application/user
//
// Priority:
// 1. Direct context value (TenantIDKey)
// 2. Cached context value (consts.TenantIDKeyInCtx)
// 3. Default tenant ID (fallback)
func GetTenantIDFromContext(c context.Context) string {
	// Priority 1: Get from direct context value
	if tenantID, ok := c.Value(TenantIDKey).(string); ok {
		return tenantID
	}

	// Priority 2: Get from ctxcache
	if tenantID, ok := ctxcache.Get[string](c, consts.TenantIDKeyInCtx); ok {
		return tenantID
	}

	// Priority 3: Fallback to default tenant ID
	logs.CtxWarnf(c, "[GetTenantIDFromContext] no tenant_id in context, using default")
	return DefaultTenantID
}

// MustGetTenantIDFromContext retrieves tenant ID from context, panics if not found
func MustGetTenantIDFromContext(c context.Context) string {
	tenantID := GetTenantIDFromContext(c)
	if tenantID == "" || tenantID == DefaultTenantID {
		logs.CtxErrorf(c, "[MustGetTenantIDFromContext] tenant_id is empty or default")
	}
	return tenantID
}
