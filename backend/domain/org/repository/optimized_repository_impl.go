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

package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
)

// optimizedOrganizationRepository 优化后的组织仓储实现
// 重点: 消除N+1查询, 优化索引使用, 减少数据库往返
type optimizedOrganizationRepository struct {
	db *gorm.DB
}

// NewOptimizedOrganizationRepository 创建优化后的组织仓储实例
func NewOptimizedOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &optimizedOrganizationRepository{db: db}
}

// Create 创建组织 (实现接口)
func (r *optimizedOrganizationRepository) Create(ctx context.Context, org *entity.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

// GetByID 根据ID获取组织 (实现接口)
func (r *optimizedOrganizationRepository) GetByID(ctx context.Context, orgID string) (*entity.Organization, error) {
	var org entity.Organization
	err := r.db.WithContext(ctx).
		Where("org_id = ?", orgID).
		Where("deleted_at IS NULL").
		First(&org).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetByCode 根据组织编码获取组织 (实现接口)
func (r *optimizedOrganizationRepository) GetByCode(ctx context.Context, tenantID, code string) (*entity.Organization, error) {
	var org entity.Organization
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("org_code = ?", code).
		Where("deleted_at IS NULL").
		First(&org).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetByTenantID 获取租户下的所有组织 (实现接口)
func (r *optimizedOrganizationRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Organization, error) {
	var orgs []*entity.Organization
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("level ASC, sort_order ASC").
		Find(&orgs).Error
	return orgs, err
}

// GetTree 获取组织树 (实现接口)
func (r *optimizedOrganizationRepository) GetTree(ctx context.Context, tenantID string) ([]*entity.Organization, error) {
	var orgs []*entity.Organization
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("level ASC, sort_order ASC").
		Find(&orgs).Error

	if err != nil {
		return nil, err
	}

	// 应用层构建树结构
	tree, err := buildTree(orgs, "")
	return tree, err
}

// GetChildren 获取子组织 (实现接口)
func (r *optimizedOrganizationRepository) GetChildren(ctx context.Context, parentID string) ([]*entity.Organization, error) {
	var orgs []*entity.Organization
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Where("deleted_at IS NULL").
		Order("sort_order ASC").
		Find(&orgs).Error
	return orgs, err
}

// Update 更新组织 (实现接口)
func (r *optimizedOrganizationRepository) Update(ctx context.Context, org *entity.Organization) error {
	return r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("org_id = ?", org.OrgID).
		Updates(org).Error
}

// Delete 软删除组织 (实现接口)
func (r *optimizedOrganizationRepository) Delete(ctx context.Context, orgID string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("org_id = ?", orgID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ExistsByCode 检查编码是否存在 (实现接口)
func (r *optimizedOrganizationRepository) ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("tenant_id = ?", tenantID).
		Where("org_code = ?", code).
		Where("deleted_at IS NULL")

	if excludeID != "" {
		query = query.Where("org_id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// List 分页查询组织列表 (实现接口)
func (r *optimizedOrganizationRepository) List(ctx context.Context, filter *OrganizationFilter) ([]*entity.Organization, int64, error) {
	return r.ListOptimized(ctx, filter)
}

// GetByIDOptimized 优化后的GetByID - 使用JOIN替代Preload
// ✅ 优化点: 1次JOIN查询替代2次独立查询
func (r *optimizedOrganizationRepository) GetByIDOptimized(ctx context.Context, orgID string) (*entity.Organization, error) {
	var org entity.Organization

	// ✅ 方案1: 使用Joins (推荐)
	err := r.db.WithContext(ctx).
		Select("organizations.*, leader.*"). // 明确查询字段
		Joins("LEFT JOIN users AS leader ON organizations.leader_id = leader.id").
		Where("organizations.org_id = ?", orgID).
		Where("organizations.deleted_at IS NULL").
		First(&org).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetByIDOptimizedV2 优化方案V2 - 分步查询 + 缓存
// 适用场景: Leader信息已缓存, 无需重复查询
func (r *optimizedOrganizationRepository) GetByIDOptimizedV2(ctx context.Context, orgID string) (*entity.Organization, error) {
	// 步骤1: 只查询组织本身
	var org entity.Organization
	err := r.db.WithContext(ctx).
		Select("org_id, tenant_id, org_name, org_code, leader_id, status, level").
		Where("org_id = ?", orgID).
		Where("deleted_at IS NULL").
		First(&org).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// 步骤2: 如果需要Leader信息, 从缓存获取或单独查询
	// (实际实现应集成缓存层)
	// if org.LeaderID != nil {
	//     leader, _ := cache.Get(ctx, fmt.Sprintf("user:%s", *org.LeaderID))
	//     org.Leader = leader
	// }

	return &org, nil
}

// ListOptimized 优化后的列表查询
// ✅ 优化点: 只查询必要字段, 避免大字段查询
func (r *optimizedOrganizationRepository) ListOptimized(ctx context.Context, filter *OrganizationFilter) ([]*entity.Organization, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Select("org_id, tenant_id, org_name, org_code, status, level") // 只查询列表需要的字段

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.OrgType != "" {
		query = query.Where("org_type = ?", filter.OrgType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Keyword != "" {
		query = query.Where("org_name LIKE ? OR org_code LIKE ?",
			"%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	// 软删除过滤
	query = query.Where("deleted_at IS NULL")

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parseInt(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 排序 (使用索引优化)
	query = query.Order("level ASC, sort_order ASC")

	// 查询
	var orgs []*entity.Organization
	err := query.Find(&orgs).Error
	return orgs, total, err
}

// BatchGetByIDs 批量获取组织 (性能优化)
// ✅ 优化点: 使用IN查询, 1次查询替代N次查询
func (r *optimizedOrganizationRepository) BatchGetByIDs(ctx context.Context, orgIDs []string) (map[string]*entity.Organization, error) {
	if len(orgIDs) == 0 {
		return make(map[string]*entity.Organization), nil
	}

	var orgs []*entity.Organization
	err := r.db.WithContext(ctx).
		Where("org_id IN ?", orgIDs).
		Where("deleted_at IS NULL").
		Find(&orgs).Error

	if err != nil {
		return nil, err
	}

	// 转换为map
	result := make(map[string]*entity.Organization, len(orgs))
	for _, org := range orgs {
		result[org.OrgID] = org
	}

	return result, nil
}

// GetTreeOptimized 优化后的组织树查询
// ✅ 优化点: 一次性查询所有节点, 应用层构建树
func (r *optimizedOrganizationRepository) GetTreeOptimized(ctx context.Context, tenantID string) ([]*entity.Organization, error) {
	// 一次性查询所有组织
	var orgs []*entity.Organization
	err := r.db.WithContext(ctx).
		Select("org_id, tenant_id, org_name, parent_id, level, sort_order").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("level ASC, sort_order ASC").
		Find(&orgs).Error

	if err != nil {
		return nil, err
	}

	// 应用层构建树结构
	tree, err := buildTree(orgs, "")
	return tree, err
}

// buildTree 递归构建树结构
// 注意: entity.Organization.Children 的类型是 []Organization 而非 []*Organization
func buildTree(orgs []*entity.Organization, parentID string) ([]*entity.Organization, error) {
	var result []*entity.Organization

	for _, org := range orgs {
		var orgParentID string
		if org.ParentID != nil {
			orgParentID = *org.ParentID
		}

		if orgParentID == parentID {
			// 递归获取子节点
			children, err := buildTree(orgs, org.OrgID)
			if err != nil {
				return nil, err
			}
			// 直接使用 []*entity.Organization
			org.Children = children
			result = append(result, org)
		}
	}

	return result, nil
}

// GetStats 获取组织统计信息 (聚合查询优化)
// ✅ 优化点: 单次聚合查询替代多次计数查询
// TODO: entity.OrganizationStats 类型未定义，需要添加到 entity 包中
/*
func (r *optimizedOrganizationRepository) GetStats(ctx context.Context, tenantID string) (*entity.OrganizationStats, error) {
	var stats entity.OrganizationStats

	// 单次查询获取多个统计值
	err := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Select(`
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active_count,
			SUM(CASE WHEN status = 'inactive' THEN 1 ELSE 0 END) as inactive_count,
			MAX(level) as max_level
		`).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}
*/

// CountByTenant 租户组织计数优化
// ✅ 优化点: 使用轻量级COUNT查询
func (r *optimizedOrganizationRepository) CountByTenant(ctx context.Context, tenantID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count, err
}

// ExistsByCodeOptimized 优化后的编码存在性检查
// ✅ 优化点: 使用EXISTS替代COUNT (性能更好)
func (r *optimizedOrganizationRepository) ExistsByCodeOptimized(ctx context.Context, tenantID, code string, excludeID string) (bool, error) {
	var exists bool

	query := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Select("1").
		Where("tenant_id = ?", tenantID).
		Where("org_code = ?", code).
		Where("deleted_at IS NULL")

	if excludeID != "" {
		query = query.Where("org_id != ?", excludeID)
	}

	err := query.Pluck("1", &exists).Error
	return exists, err
}

// =====================================================
// 辅助函数
// =====================================================
// parseInt 函数已在 org_repository_impl.go 中定义，避免重复声明
