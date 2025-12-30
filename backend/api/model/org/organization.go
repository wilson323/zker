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

// CreateOrganizationRequest 创建组织请求
type CreateOrganizationRequest struct {
	TenantID     string `json:"tenant_id" binding:"required"`
	OrgName      string `json:"org_name" binding:"required,min=1,max=200"`
	OrgType      string `json:"org_type" binding:"required,oneof=company division department project"`
	ParentID     *string `json:"parent_id,omitempty"`
	OrgCode      string `json:"org_code" binding:"required,min=1,max=50"`
	LeaderID     *string `json:"leader_id,omitempty"`
	Description  string `json:"description,omitempty"`
	SortOrder    int    `json:"sort_order,omitempty"`
}

// CreateOrganizationResponse 创建组织响应
type CreateOrganizationResponse struct {
	Code    int32             `json:"code"`
	Message string            `json:"msg"`
	Data    *OrganizationData `json:"data,omitempty"`
}

// OrganizationData 组织数据
type OrganizationData struct {
	OrgID        string `json:"org_id"`
	TenantID     string `json:"tenant_id"`
	OrgName      string `json:"org_name"`
	OrgType      string `json:"org_type"`
	ParentID     *string `json:"parent_id,omitempty"`
	OrgCode      string `json:"org_code"`
	Level        int    `json:"level"`
	Path         string `json:"path"`
	SortOrder    int    `json:"sort_order"`
	Status       string `json:"status"`
	Description  string `json:"description"`
	LeaderID     *string `json:"leader_id,omitempty"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

// GetOrganizationResponse 获取组织响应
type GetOrganizationResponse struct {
	Code    int32             `json:"code"`
	Message string            `json:"msg"`
	Data    *OrganizationData `json:"data,omitempty"`
}

// UpdateOrganizationRequest 更新组织请求
type UpdateOrganizationRequest struct {
	OrgID       string  `json:"org_id" binding:"required"`
	OrgName     string  `json:"org_name" binding:"omitempty,min=1,max=200"`
	LeaderID    *string `json:"leader_id,omitempty"`
	Description string  `json:"description,omitempty"`
	SortOrder   int     `json:"sort_order,omitempty"`
	Status      string  `json:"status" binding:"omitempty,oneof=active inactive frozen"`
}

// UpdateOrganizationResponse 更新组织响应
type UpdateOrganizationResponse struct {
	Code    int32             `json:"code"`
	Message string            `json:"msg"`
	Data    *OrganizationData `json:"data,omitempty"`
}

// DeleteOrganizationResponse 删除组织响应
type DeleteOrganizationResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"msg"`
}

// MoveOrganizationRequest 移动组织请求
type MoveOrganizationRequest struct {
	OrgID       string  `json:"org_id" binding:"required"`
	NewParentID *string `json:"new_parent_id,omitempty"`
}

// MoveOrganizationResponse 移动组织响应
type MoveOrganizationResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"msg"`
}

// GetOrganizationTreeResponse 获取组织树响应
type GetOrganizationTreeResponse struct {
	Code    int32               `json:"code"`
	Message string              `json:"msg"`
	Data    []*OrganizationData `json:"data,omitempty"`
}

// ListOrganizationsRequest 列出组织请求
type ListOrganizationsRequest struct {
	TenantID  string `json:"tenant_id" binding:"required"`
	OrgType   string `json:"org_type,omitempty" binding:"omitempty,oneof=company division department project"`
	Status    string `json:"status,omitempty" binding:"omitempty,oneof=active inactive frozen"`
	ParentID  string `json:"parent_id,omitempty"`
	Keyword   string `json:"keyword,omitempty"`
	Page      int    `json:"page,omitempty" binding:"omitempty,min=1"`
	PageSize  int    `json:"page_size,omitempty" binding:"omitempty,min=1,max=100"`
	SortBy    string `json:"sort_by,omitempty" binding:"omitempty,oneof=created_at updated_at sort_order org_name"`
	SortOrder string `json:"sort_order,omitempty" binding:"omitempty,oneof=asc desc"`
}

// ListOrganizationsResponse 列出组织响应
type ListOrganizationsResponse struct {
	Code      int32               `json:"code"`
	Message   string              `json:"msg"`
	Data      []*OrganizationData `json:"data,omitempty"`
	Total     int64               `json:"total,omitempty"`
	Page      int                 `json:"page,omitempty"`
	PageSize  int                 `json:"page_size,omitempty"`
}

// GetChildrenResponse 获取子组织响应
type GetChildrenResponse struct {
	Code    int32               `json:"code"`
	Message string              `json:"msg"`
	Data    []*OrganizationData `json:"data,omitempty"`
}

// GetAncestorsResponse 获取祖先组织响应
type GetAncestorsResponse struct {
	Code    int32               `json:"code"`
	Message string              `json:"msg"`
	Data    []*OrganizationData `json:"data,omitempty"`
}

// GetDescendantsResponse 获取后代组织响应
type GetDescendantsResponse struct {
	Code    int32               `json:"code"`
	Message string              `json:"msg"`
	Data    []*OrganizationData `json:"data,omitempty"`
}
