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

package coze

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
)

// InvalidParamRequestResponse 返回无效参数错误响应
func InvalidParamRequestResponse(c *app.RequestContext, errMsg string) {
	httputil.BadRequest(c, errMsg)
}

// InternalServerErrorResponse 返回内部服务器错误响应
func InternalServerErrorResponse(ctx context.Context, c *app.RequestContext, err error) {
	httputil.InternalError(ctx, c, err)
}

// 内部使用（保持向后兼容）
func invalidParamRequestResponse(c *app.RequestContext, errMsg string) {
	InvalidParamRequestResponse(c, errMsg)
}

func internalServerErrorResponse(ctx context.Context, c *app.RequestContext, err error) {
	InternalServerErrorResponse(ctx, c, err)
}
