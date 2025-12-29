# 创建 Repository Skill

## 技能描述

创建符合 Coze Studio DDD 架构的 Repository 实现，包括接口定义、GORM 实现、错误处理和测试。

## 适用场景

- 需要持久化领域实体
- 需要复杂查询逻辑
- 需要事务支持
- 需要缓存数据

## 工作流程

### 1. 分析仓储需求

明确以下信息：
- **实体类型**: 管理哪种实体
- **存储介质**: MySQL、Redis、Elasticsearch 等
- **查询需求**: 需要哪些查询方法
- **性能要求**: 是否需要缓存、索引优化
- **事务需求**: 是否需要事务支持

### 2. 定义仓储接口

**文件位置**: `domain/{domain}/repository/interface.go`

```go
package repository

import (
	"context"
	"github.com/coze-dev/coze-studio/backend/domain/{domain}/entity"
)

// {Entity}Repository {entity}仓储接口
type {Entity}Repository interface {
	// 基础 CRUD 操作

	// Create 创建{entity}
	// 参数：
	//   ctx - 请求上下文
	//   {entity} - {entity}实体
	//
	// 返回：
	//   int64 - 创建的{entity} ID
	//   error - 错误信息
	Create(ctx context.Context, {entity} *entity.{Entity}) (int64, error)

	// FindByID 根据ID查询{entity}
	// 返回：
	//   *entity.{Entity} - {entity}实体
	//   bool - 是否存在
	//   error - 错误信息
	FindByID(ctx context.Context, id int64) (*entity.{Entity}, bool, error)

	// Update 更新{entity}
	Update(ctx context.Context, {entity} *entity.{Entity}) error

	// Delete 删除{entity}（软删除）
	Delete(ctx context.Context, id int64) error

	// 查询方法

	// FindByEmail 根据邮箱查询
	FindByEmail(ctx context.Context, email string) (*entity.{Entity}, bool, error)

	// FindByStatus 根据状态查询
	FindByStatus(ctx context.Context, status entity.{Entity}Status) ([]*entity.{Entity}, error)

	// FindBySpaceID 查询空间下的{entity}列表
	FindBySpaceID(ctx context.Context, spaceID int64, params *QueryParams) ([]*entity.{Entity}, error)

	// 存在性检查

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByID 检查ID是否存在
	ExistsByID(ctx context.Context, id int64) (bool, error)

	// 分页查询

	// FindWithPagination 分页查询{entity}
	// 返回：
	//   []*entity.{Entity} - {entity}列表
	//   int64 - 总数
	//   error - 错误信息
	FindWithPagination(ctx context.Context, params *PaginationParams) ([]*entity.{Entity}, int64, error)

	// 统计方法

	// CountBySpaceID 统计空间下的{entity}数量
	CountBySpaceID(ctx context.Context, spaceID int64) (int64, error)

	// 批量操作

	// BatchCreate 批量创建{entity}
	BatchCreate(ctx context.Context, {entities} []*entity.{Entity}) error

	// BatchDelete 批量删除{entity}
	BatchDelete(ctx context.Context, ids []int64) error

	// 事务操作

	// Transaction 事务执行
	Transaction(ctx context.Context, fn func(repo {Entity}Repository) error) error
}

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int    // 页码（从1开始）
	PageSize int    // 每页数量
	OrderBy  string // 排序字段
	Order    string // 排序方向：ASC/DESC
}

// QueryParams 查询参数
type QueryParams struct {
	Status   string             // 状态过滤
	Keyword  string             // 关键词搜索
	DateFrom *time.Time         // 起始日期
	DateTo   *time.Time         // 结束日期
	Filters  map[string]interface{} // 其他过滤条件
}
```

### 3. 实现 Repository

**文件位置**: `domain/{domain}/repository/repository_impl.go` 或 `infra/database/repository/{entity}_repository.go`

```go
package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/coze-dev/coze-studio/backend/domain/{domain}/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// {entity}RepositoryImpl {entity}仓储实现
type {entity}RepositoryImpl struct {
	db *gorm.DB
}

// New{Entity}Repository 创建{entity}仓储
func New{Entity}Repository(db *gorm.DB) {Entity}Repository {
	return &{entity}RepositoryImpl{
		db: db,
	}
}

// Create 创建{entity}
func (r *{entity}RepositoryImpl) Create(ctx context.Context, {entity} *entity.{Entity}) (int64, error) {
	if err := r.db.WithContext(ctx).Create({entity}).Error; err != nil {
		return 0, errorx.Wrapf(err, "create {entity} failed")
	}
	return {entity}.ID, nil
}

// FindByID 根据ID查询{entity}
func (r *{entity}RepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.{Entity}, bool, error) {
	var {entity} entity.{Entity}
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&{entity}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errorx.Wrapf(err, "find {entity} by id failed, id=%d", id)
	}

	return &{entity}, true, nil
}

// Update 更新{entity}
func (r *{entity}RepositoryImpl) Update(ctx context.Context, {entity} *entity.{Entity}) error {
	result := r.db.WithContext(ctx).
		Model(&entity.{Entity}{}).
		Where("id = ? AND deleted_at IS NULL", {entity}.ID).
		Updates({entity})

	if result.Error != nil {
		return errorx.Wrapf(result.Error, "update {entity} failed, id=%d", {entity}.ID)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete 删除{entity}（软删除）
func (r *{entity}RepositoryImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).
		Select("deleted_at").
		Delete(&entity.{Entity}{}, id)

	if result.Error != nil {
		return errorx.Wrapf(result.Error, "delete {entity} failed, id=%d", id)
	}

	return nil
}

// FindByEmail 根据邮箱查询
func (r *{entity}RepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.{Entity}, bool, error) {
	var {entity} entity.{Entity}
	err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&{entity}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errorx.Wrapf(err, "find {entity} by email failed, email=%s", email)
	}

	return &{entity}, true, nil
}

// FindByStatus 根据状态查询
func (r *{entity}RepositoryImpl) FindByStatus(ctx context.Context, status entity.{Entity}Status) ([]*entity.{Entity}, error) {
	var {entities} []*entity.{Entity}
	err := r.db.WithContext(ctx).
		Where("status = ? AND deleted_at IS NULL", status).
		Find(&{entities}).Error

	if err != nil {
		return nil, errorx.Wrapf(err, "find {entities} by status failed, status=%s", status)
	}

	return {entities}, nil
}

// FindBySpaceID 查询空间下的{entity}列表
func (r *{entity}RepositoryImpl) FindBySpaceID(
	ctx context.Context,
	spaceID int64,
	params *QueryParams,
) ([]*entity.{Entity}, error) {
	query := r.db.WithContext(ctx).
		Model(&entity.{Entity}{}).
		Where("space_id = ? AND deleted_at IS NULL", spaceID)

	// 动态条件
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+params.Keyword+"%")
	}
	if params.DateFrom != nil {
		query = query.Where("created_at >= ?", params.DateFrom)
	}
	if params.DateTo != nil {
		query = query.Where("created_at <= ?", params.DateTo)
	}

	// 其他过滤条件
	for key, value := range params.Filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	var {entities} []*entity.{Entity}
	if err := query.Find(&{entities}).Error; err != nil {
		return nil, errorx.Wrapf(err, "find {entities} by space_id failed, spaceID=%d", spaceID)
	}

	return {entities}, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *{entity}RepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.{Entity}{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error

	if err != nil {
		return false, errorx.Wrapf(err, "check email existence failed, email=%s", email)
	}

	return count > 0, nil
}

// ExistsByID 检查ID是否存在
func (r *{entity}RepositoryImpl) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.{Entity}{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error

	if err != nil {
		return false, errorx.Wrapf(err, "check id existence failed, id=%d", id)
	}

	return count > 0, nil
}

// FindWithPagination 分页查询
func (r *{entity}RepositoryImpl) FindWithPagination(
	ctx context.Context,
	params *PaginationParams,
) ([]*entity.{Entity}, int64, error) {
	var {entities} []*entity.{Entity}
	var total int64

	// 构建查询
	query := r.db.WithContext(ctx).
		Model(&entity.{Entity}{}).
		Where("deleted_at IS NULL")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errorx.Wrapf(err, "count {entities} failed")
	}

	// 排序
	orderBy := "created_at"
	if params.OrderBy != "" {
		orderBy = params.OrderBy
	}

	order := "DESC"
	if params.Order != "" {
		order = params.Order
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	err := query.
		Order(fmt.Sprintf("%s %s", orderBy, order)).
		Offset(offset).
		Limit(params.PageSize).
		Find(&{entities}).Error

	if err != nil {
		return nil, 0, errorx.Wrapf(err, "find {entities} with pagination failed")
	}

	return {entities}, total, nil
}

// CountBySpaceID 统计空间下的{entity}数量
func (r *{entity}RepositoryImpl) CountBySpaceID(ctx context.Context, spaceID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.{Entity}{}).
		Where("space_id = ? AND deleted_at IS NULL", spaceID).
		Count(&count).Error

	if err != nil {
		return 0, errorx.Wrapf(err, "count {entities} failed, spaceID=%d", spaceID)
	}

	return count, nil
}

// BatchCreate 批量创建
func (r *{entity}RepositoryImpl) BatchCreate(ctx context.Context, {entities} []*entity.{Entity}) error {
	if len({entities}) == 0 {
		return nil
	}

	// 使用 ON DUPLICATE KEY UPDATE 提高性能
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches({entities}, 100).
		Error; err != nil {
		return errorx.Wrapf(err, "batch create {entities} failed, count=%d", len({entities}))
	}

	return nil
}

// BatchDelete 批量删除
func (r *{entity}RepositoryImpl) BatchDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Delete(&entity.{Entity}{}, ids)

	if result.Error != nil {
		return errorx.Wrapf(result.Error, "batch delete {entities} failed, count=%d", len(ids))
	}

	return nil
}

// Transaction 事务执行
func (r *{entity}RepositoryImpl) Transaction(ctx context.Context, fn func(repo {Entity}Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &{entity}RepositoryImpl{db: tx}
		return fn(txRepo)
	})
}
```

### 4. Repository 测试

**测试文件**: `repository_test.go`

```go
package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/{domain}/entity"
	"github.com/coze-dev/coze-studio/backend/domain/{domain}/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&entity.{Entity}{})
	require.NoError(t, err)

	return db
}

func Test{Entity}Repository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New{Entity}Repository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		{entity} := &entity.{Entity}{
			Name:    "Test {Entity}",
			Email:   "test@example.com",
			SpaceID: 1,
			Status:  entity.{Entity}StatusActive,
		}

		id, err := repo.Create(ctx, {entity})

		require.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("duplicate email", func(t *testing.T) {
		{entity}1 := &entity.{Entity}{
			Name:    "Test 1",
			Email:   "duplicate@example.com",
			SpaceID: 1,
		}
		{entity}2 := &entity.{Entity}{
			Name:    "Test 2",
			Email:   "duplicate@example.com",
			SpaceID: 1,
		}

		_, err := repo.Create(ctx, {entity}1)
		require.NoError(t, err)

		_, err = repo.Create(ctx, {entity}2)
		assert.Error(t, err)
	})
}

func Test{Entity}Repository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New{Entity}Repository(db)
	ctx := context.Background()

	// 创建测试数据
	{entity} := &entity.{Entity}{
		Name:    "Test",
		Email:   "test@example.com",
		SpaceID: 1,
	}
	id, _ := repo.Create(ctx, {entity})

	t.Run("found", func(t *testing.T) {
		found, exist, err := repo.FindByID(ctx, id)

		require.NoError(t, err)
		assert.True(t, exist)
		assert.NotNil(t, found)
		assert.Equal(t, "Test", found.Name)
	})

	t.Run("not found", func(t *testing.T) {
		found, exist, err := repo.FindByID(ctx, 99999)

		require.NoError(t, err)
		assert.False(t, exist)
		assert.Nil(t, found)
	})
}

func Test{Entity}Repository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New{Entity}Repository(db)
	ctx := context.Background()

	{entity} := &entity.{Entity}{
		Name:    "Original",
		Email:   "original@example.com",
		SpaceID: 1,
	}
	id, _ := repo.Create(ctx, {entity})

	t.Run("success", func(t *testing.T) {
		{entity}.Name = "Updated"
		err := repo.Update(ctx, {entity})

		require.NoError(t, err)

		updated, exist, _ := repo.FindByID(ctx, id)
		assert.True(t, exist)
		assert.Equal(t, "Updated", updated.Name)
	})

	t.Run("not found", func(t *testing.T) {
		nonExistent := &entity.{Entity}{ID: 99999}
		err := repo.Update(ctx, nonExistent)

		assert.Error(t, err)
	})
}

func Test{Entity}Repository_ExistsByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New{Entity}Repository(db)
	ctx := context.Background()

	// 创建测试数据
	{entity} := &entity.{Entity}{
		Email:   "exists@example.com",
		SpaceID: 1,
	}
	repo.Create(ctx, {entity})

	t.Run("exists", func(t *testing.T) {
		exists, err := repo.ExistsByEmail(ctx, "exists@example.com")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not exists", func(t *testing.T) {
		exists, err := repo.ExistsByEmail(ctx, "notexists@example.com")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func Test{Entity}Repository_FindWithPagination(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New{Entity}Repository(db)
	ctx := context.Background()

	// 创建测试数据
	for i := 1; i <= 25; i++ {
		{entity} := &entity.{Entity}{
			Name:    fmt.Sprintf("Entity %d", i),
			Email:   fmt.Sprintf("entity%d@example.com", i),
			SpaceID: 1,
		}
		repo.Create(ctx, {entity})
	}

	t.Run("first page", func(t *testing.T) {
		params := &repository.PaginationParams{
			Page:     1,
			PageSize: 10,
		}

		entities, total, err := repo.FindWithPagination(ctx, params)

		require.NoError(t, err)
		assert.Equal(t, int64(25), total)
		assert.Len(t, entities, 10)
	})

	t.Run("second page", func(t *testing.T) {
		params := &repository.PaginationParams{
			Page:     2,
			PageSize: 10,
		}

		entities, total, err := repo.FindWithPagination(ctx, params)

		require.NoError(t, err)
		assert.Equal(t, int64(25), total)
		assert.Len(t, entities, 10)
	})
}
```

## 必须遵循的规范

### ✅ 命名规范

```go
// ✅ 正确
type UserRepository interface {}
type userRepositoryImpl struct {}
func NewUserRepository(db *gorm.DB) UserRepository {}

// ❌ 错误
type IUserRepository interface {}           ← 不使用 I 前缀
type UserRepositoryImpl struct {}           ← 内部实现应小写
func NewRepo(db *gorm.DB) UserRepository {}   ← 过于通用
```

### ✅ 方法命名规范

```go
// ✅ 正确的方法命名
Create(ctx, entity) (int64, error)
FindByID(ctx, id) (*Entity, bool, error)
FindByEmail(ctx, email) (*Entity, bool, error)
Update(ctx, entity) error
Delete(ctx, id) error
ExistsByEmail(ctx, email) (bool, error)
CountBySpaceID(ctx, spaceID) (int64, error)

// ❌ 错误的方法命名
Get(ctx, id) (*Entity, error)              ← 应使用 FindByID
Save(ctx, entity) error                    ← 应使用 Create/Update
Remove(ctx, id) error                     ← 应使用 Delete
CheckEmail(ctx, email) (bool, error)       ← 应使用 ExistsByEmail
```

### ✅ GORM 使用规范

```go
// ✅ 正确的 GORM 使用
db.Where("id = ? AND deleted_at IS NULL", id).First(&entity)
db.Where("status = ?", status).Find(&entities)
db.Create(&entity)
db.Updates(&entity)
db.Delete(&entity)

// ❌ 错误的 GORM 使用
db.Where("id = ?", id).First(&entity)      ← 缺少软删除检查
db.Raw("SELECT * FROM users")              ← 应使用 ORM 方法
db.Exec("DELETE FROM users")               ← 应使用 Delete 方法
```

### ✅ 错误处理规范

```go
// ✅ 正确的错误处理
if err != nil {
    if err == gorm.ErrRecordNotFound {
        return nil, false, nil
    }
    return nil, false, errorx.Wrapf(err, "operation failed, id=%d", id)
}

// ❌ 错误的错误处理
if err != nil {
    return nil, err                         ← 应添加上下文
}
return nil, err, nil                        ← 未区分错误类型
```

### ✅ 软删除规范

```go
// ✅ 正确的软删除
db.Where("id = ? AND deleted_at IS NULL", id).Delete(&entity)

// ✅ 正确的查询
db.Where("deleted_at IS NULL").Find(&entities)

// ❌ 错误的查询
db.Where("id = ?", id).Delete(&entity)     ← 未检查软删除状态
db.Find(&entities)                         ← 未过滤已删除数据
```

## 输出检查清单

创建完成后，必须确保：

- [ ] 仓储接口定义清晰
- [ ] 仓储实现命名规范
- [ ] 方法签名完整
- [ ] 使用 GORM ORM 方法
- [ ] 正确处理软删除
- [ ] 使用 WithContext(ctx)
- [ ] 错误使用 errorx.Wrapf 包装
- [ ] 包含完整的单元测试
- [ ] 测试覆盖正常和异常场景
- [ ] 通过 gofmt 检查
- [ ] 通过 golangci-lint 检查

## 注意事项

1. **软删除优先**: 默认使用软删除，保留数据
2. **索引优化**: 根据查询需求添加索引
3. **批量操作**: 使用批量方法提高性能
4. **事务处理**: 正确使用事务保证一致性
5. **查询优化**: 避免 N+1 查询，使用预加载
6. **连接池**: 合理配置数据库连接池
7. **超时控制**: 设置合理的查询超时时间
