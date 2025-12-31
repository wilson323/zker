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

package ctxutil

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

func GetUserSessionFromCtx(ctx context.Context) *entity.Session {
	data, ok := ctxcache.Get[*entity.Session](ctx, consts.SessionDataKeyInCtx)
	if !ok {
		return nil
	}

	return data
}

// GetUIDFromCtx 获取用户ID，如果未登录返回error
// 这是新的推荐方法，返回error而不是panic
func GetUIDFromCtxWithError(ctx context.Context) (int64, error) {
	sessionData := GetUserSessionFromCtx(ctx)
	if sessionData == nil {
		return 0, fmt.Errorf("session data is required")
	}
	return sessionData.UserID, nil
}

// GetUIDFromCtx 获取用户ID的指针，未登录返回nil
// 保留向后兼容
func GetUIDFromCtx(ctx context.Context) *int64 {
	sessionData := GetUserSessionFromCtx(ctx)
	if sessionData == nil {
		return nil
	}
	return &sessionData.UserID
}

// MustGetUIDFromCtx 获取用户ID，失败则panic
// Deprecated: 使用 GetUIDFromCtxWithError 替代以获得更好的错误处理
func MustGetUIDFromCtx(ctx context.Context) int64 {
	userID, err := GetUIDFromCtxWithError(ctx)
	if err != nil {
		// 记录日志后panic（向后兼容）
		// 生产环境应该使用 GetUIDFromCtxWithError 并处理error
		panic(fmt.Errorf("MustGetUIDFromCtx: %w", err))
	}
	return userID
}
