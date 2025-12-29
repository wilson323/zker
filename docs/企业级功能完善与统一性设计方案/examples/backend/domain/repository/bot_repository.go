// Package repository 定义领域仓储接口
//
// 本文件展示Bot仓储接口的完整实现，包括：
// - CRUD操作接口定义
// - 查询方法接口定义
// - 分页和排序支持
// - 多租户隔离
package repository

import (
	"context"
	"time"

	"github.com/coze-studio/domain/entity"
)

// BotRepository Bot仓储接口
//
// 遵循依赖倒置原则（DIP），领域层定义接口，基础设施层实现
// 这样领域层不依赖具体的数据访问技术，便于测试和切换实现
type BotRepository interface {
	// ============ 基础 CRUD 操作 ============

	// Save 保存Bot（创建或更新）
	Save(ctx context.Context, bot *entity.Bot) error

	// FindByID 根据ID查找Bot
	FindByID(ctx context.Context, tenantID entity.TenantID, botID string) (*entity.Bot, error)

	// Delete 删除Bot（软删除）
	Delete(ctx context.Context, tenantID entity.TenantID, botID string) error

	// ============ 批量操作 ============

	// FindByTenantID 查找租户下的所有Bot
	FindByTenantID(ctx context.Context, tenantID entity.TenantID) ([]*entity.Bot, error)

	// FindByTenantIDWithPaging 分页查找租户下的Bot
	FindByTenantIDWithPaging(
		ctx context.Context,
		tenantID entity.TenantID,
		req *ListBotsRequest,
	) (*BotListResult, error)

	// FindByIDs 根据ID列表批量查找Bot
	FindByIDs(ctx context.Context, tenantID entity.TenantID, botIDs []string) ([]*entity.Bot, error)

	// ============ 查询方法 ============

	// FindByName 根据名称查找Bot
	FindByName(ctx context.Context, tenantID entity.TenantID, name string) (*entity.Bot, error)

	// FindPublished 查找已发布的Bot
	FindPublished(ctx context.Context, tenantID entity.TenantID, published bool) ([]*entity.Bot, error)

	// FindByTag 根据标签查找Bot
	FindByTag(ctx context.Context, tenantID entity.TenantID, tag string) ([]*entity.Bot, error)

	// FindByCreator 根据创建者查找Bot
	FindByCreator(ctx context.Context, tenantID entity.TenantID, creatorID string) ([]*entity.Bot, error)

	// ============ 统计方法 ============

	// CountByTenantID 统计租户下的Bot数量
	CountByTenantID(ctx context.Context, tenantID entity.TenantID) (int64, error)

	// CountByStatus 根据状态统计Bot数量
	CountByStatus(ctx context.Context, tenantID entity.TenantID, status entity.BotStatus) (int64, error)

	// ============ 事务支持 ============

	// WithTx 在事务中执行操作
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// ============ 请求和响应 DTO ============

// ListBotsRequest 列表查询请求
type ListBotsRequest struct {
	// TenantID 租户ID（从context中获取）
	TenantID entity.TenantID `json:"-"`

	// Page 页码（从1开始）
	Page int `json:"page" validate:"min=1"`
	// PageSize 每页大小
	PageSize int `json:"page_size" validate:"min=1,max=100"`
	// SortBy 排序字段（id, name, created_at, updated_at）
	SortBy string `json:"sort_by" validate:"oneof=id name created_at updated_at"`
	// SortOrder 排序方向（asc, desc）
	SortOrder string `json:"sort_order" validate:"oneof=asc desc"`

	// Name 名称模糊搜索
	Name string `json:"name"`
	// Status 状态过滤
	Status *entity.BotStatus `json:"status"`
	// Published 是否发布过滤
	Published *bool `json:"published"`
	// Tag 标签过滤
	Tag string `json:"tag"`
	// CreatorID 创建者过滤
	CreatorID string `json:"creator_id"`
	// StartTime 创建时间范围-开始
	StartTime *time.Time `json:"start_time"`
	// EndTime 创建时间范围-结束
	EndTime *time.Time `json:"end_time"`
}

// BotListResult Bot列表结果
type BotListResult struct {
	// Bots Bot列表
	Bots []*entity.Bot `json:"bots"`
	// Total 总数
	Total int64 `json:"total"`
	// Page 当前页码
	Page int `json:"page"`
	// PageSize 每页大小
	PageSize int `json:"page_size"`
	// TotalPages 总页数
	TotalPages int64 `json:"total_pages"`
}

// NewListBotsRequest 创建列表查询请求（带默认值）
func NewListBotsRequest() *ListBotsRequest {
	return &ListBotsRequest{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// Validate 验证请求参数
func (r *ListBotsRequest) Validate() error {
	// 这里可以使用验证库如 go-playground/validator
	// 为了简洁，这里省略具体实现
	return nil
}

// GetOffset 计算偏移量
func (r *ListBotsRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

// HasFilters 检查是否有过滤条件
func (r *ListBotsRequest) HasFilters() bool {
	return r.Name != "" ||
		r.Status != nil ||
		r.Published != nil ||
		r.Tag != "" ||
		r.CreatorID != "" ||
		r.StartTime != nil ||
		r.EndTime != nil
}
