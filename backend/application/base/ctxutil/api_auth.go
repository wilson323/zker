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

	"github.com/coze-dev/coze-studio/backend/domain/openauth/openapiauth/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

func GetApiAuthFromCtx(ctx context.Context) *entity.ApiKey {
	data, ok := ctxcache.Get[*entity.ApiKey](ctx, consts.OpenapiAuthKeyInCtx)

	if !ok {
		return nil
	}
	return data
}

// GetUIDFromApiAuthCtx 从API Auth上下文获取用户ID
// 这是新的推荐方法，返回error而不是panic
func GetUIDFromApiAuthCtx(ctx context.Context) (int64, error) {
	apiKeyInfo := GetApiAuthFromCtx(ctx)
	if apiKeyInfo == nil {
		return 0, fmt.Errorf("api auth info is required")
	}
	return apiKeyInfo.UserID, nil
}

// MustGetUIDFromApiAuthCtx 获取用户ID，失败则panic
// Deprecated: 使用 GetUIDFromApiAuthCtx 替代以获得更好的错误处理
func MustGetUIDFromApiAuthCtx(ctx context.Context) int64 {
	userID, err := GetUIDFromApiAuthCtx(ctx)
	if err != nil {
		// 记录日志后panic（向后兼容）
		// 生产环境应该使用 GetUIDFromApiAuthCtx 并处理error
		panic(fmt.Errorf("MustGetUIDFromApiAuthCtx: %w", err))
	}
	return userID
}
