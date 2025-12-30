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

// OrganizationTree 组织树闭包表（用于高效查询组织层级关系）
// 使用闭包表（Closure Table）模式存储组织树的所有路径
type OrganizationTree struct {
	ID         int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID   string `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	AncestorID string `json:"ancestor_id" gorm:"type:varchar(36);not null;index:idx_ancestor,index:idx_ancestor_depth"` // 祖先节点ID
	DescendantID string `json:"descendant_id" gorm:"type:varchar(36);not null;index:idx_descendant,index:idx_ancestor_depth"` // 后代节点ID
	Depth      int    `json:"depth" gorm:"type:int;not null;index:idx_ancestor_depth"` // 层级深度：0表示自己，1表示子节点，2表示孙节点
}

// TableName 指定表名
func (OrganizationTree) TableName() string {
	return "organization_trees"
}

// DepartmentTree 部门树闭包表（用于高效查询部门层级关系）
type DepartmentTree struct {
	ID         int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID   string `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	AncestorID string `json:"ancestor_id" gorm:"type:varchar(36);not null;index:idx_ancestor,index:idx_ancestor_depth"` // 祖先节点ID
	DescendantID string `json:"descendant_id" gorm:"type:varchar(36);not null;index:idx_descendant,index:idx_ancestor_depth"` // 后代节点ID
	Depth      int    `json:"depth" gorm:"type:int;not null;index:idx_ancestor_depth"` // 层级深度：0表示自己，1表示子节点，2表示孙节点
}

// TableName 指定表名
func (DepartmentTree) TableName() string {
	return "department_trees"
}
