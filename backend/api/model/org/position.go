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

package org

// CreatePositionRequest 创建岗位请求
type CreatePositionRequest struct {
	TenantID         string  `json:"tenant_id" binding:"required"`
	DeptID           *string `json:"dept_id,omitempty"`
	PositionName     string  `json:"position_name" binding:"required,min=1,max=100"`
	PositionCode     string  `json:"position_code" binding:"required,min=1,max=50"`
	Level            int     `json:"level" binding:"required,min=1,max=10"`
	Category         string  `json:"category" binding:"required,min=1,max=50"`
	Responsibilities string  `json:"responsibilities,omitempty"`
	Requirements     string  `json:"requirements,omitempty"`
	SortOrder        int     `json:"sort_order,omitempty"`
}

// CreatePositionResponse 创建岗位响应
type CreatePositionResponse struct {
	Code    int32         `json:"code"`
	Message string        `json:"msg"`
	Data    *PositionData `json:"data,omitempty"`
}

// PositionData 岗位数据
type PositionData struct {
	PositionID       string  `json:"position_id"`
	TenantID         string  `json:"tenant_id"`
	DeptID           *string `json:"dept_id,omitempty"`
	PositionName     string  `json:"position_name"`
	PositionCode     string  `json:"position_code"`
	Level            int     `json:"level"`
	Category         string  `json:"category"`
	Responsibilities string  `json:"responsibilities"`
	Requirements     string  `json:"requirements"`
	Status           string  `json:"status"`
	SortOrder        int     `json:"sort_order"`
	CreatedAt        int64   `json:"created_at"`
	UpdatedAt        int64   `json:"updated_at"`
}

// GetPositionResponse 获取岗位响应
type GetPositionResponse struct {
	Code    int32         `json:"code"`
	Message string        `json:"msg"`
	Data    *PositionData `json:"data,omitempty"`
}

// UpdatePositionRequest 更新岗位请求
type UpdatePositionRequest struct {
	PositionID       string `json:"position_id" binding:"required"`
	PositionName     string `json:"position_name" binding:"omitempty,min=1,max=100"`
	Level            int    `json:"level" binding:"omitempty,min=1,max=10"`
	Category         string `json:"category" binding:"omitempty,min=1,max=50"`
	Responsibilities string `json:"responsibilities,omitempty"`
	Requirements     string `json:"requirements,omitempty"`
	SortOrder        int    `json:"sort_order,omitempty"`
	Status           string `json:"status" binding:"omitempty,oneof=active inactive frozen"`
}

// UpdatePositionResponse 更新岗位响应
type UpdatePositionResponse struct {
	Code    int32         `json:"code"`
	Message string        `json:"msg"`
	Data    *PositionData `json:"data,omitempty"`
}

// DeletePositionResponse 删除岗位响应
type DeletePositionResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"msg"`
}

// ListPositionsRequest 列出岗位请求
type ListPositionsRequest struct {
	TenantID  string `json:"tenant_id" binding:"required"`
	DeptID    string `json:"dept_id,omitempty"`
	Level     *int   `json:"level,omitempty"`
	Category  string `json:"category,omitempty"`
	Status    string `json:"status,omitempty" binding:"omitempty,oneof=active inactive frozen"`
	Keyword   string `json:"keyword,omitempty"`
	Page      int    `json:"page,omitempty" binding:"omitempty,min=1"`
	PageSize  int    `json:"page_size,omitempty" binding:"omitempty,min=1,max=100"`
	SortBy    string `json:"sort_by,omitempty" binding:"omitempty,oneof=created_at updated_at sort_order position_name"`
	SortOrder string `json:"sort_order_dir,omitempty" binding:"omitempty,oneof=asc desc"`
}

// ListPositionsResponse 列出岗位响应
type ListPositionsResponse struct {
	Code     int32           `json:"code"`
	Message  string          `json:"msg"`
	Data     []*PositionData `json:"data,omitempty"`
	Total    int64           `json:"total,omitempty"`
	Page     int             `json:"page,omitempty"`
	PageSize int             `json:"page_size,omitempty"`
}

// GetPositionByCodeResponse 根据编码获取岗位响应
type GetPositionByCodeResponse struct {
	Code    int32         `json:"code"`
	Message string        `json:"msg"`
	Data    *PositionData `json:"data,omitempty"`
}

// GetPositionsByDepartmentResponse 获取部门的所有岗位响应
type GetPositionsByDepartmentResponse struct {
	Code    int32           `json:"code"`
	Message string          `json:"msg"`
	Data    []*PositionData `json:"data,omitempty"`
}

// GetPositionsByLevelResponse 按职级获取岗位响应
type GetPositionsByLevelResponse struct {
	Code    int32           `json:"code"`
	Message string          `json:"msg"`
	Data    []*PositionData `json:"data,omitempty"`
}

// GetPositionsByCategoryResponse 按类别获取岗位响应
type GetPositionsByCategoryResponse struct {
	Code    int32           `json:"code"`
	Message string          `json:"msg"`
	Data    []*PositionData `json:"data,omitempty"`
}
