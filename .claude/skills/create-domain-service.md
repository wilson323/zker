# 创建领域服务 Skill

## 技能描述

创建符合 DDD 架构的 Go 领域服务，包括服务接口、实现、错误处理和测试。

## 适用场景

- 需要创建新的领域服务
- 需要实现复杂业务逻辑
- 需要跨多个实体的业务操作

## 工作流程

### 1. 分析领域需求

在创建服务之前，必须明确：
- **领域边界**：服务属于哪个领域（user、agent、workflow 等）
- **业务能力**：服务提供哪些业务功能
- **依赖资源**：需要哪些 Repository、外部服务
- **业务规则**：需要验证哪些业务规则

### 2. 定义服务接口

**文件位置**: `domain/{domain}/service/service.go`

```go
package service

import (
    "context"
    "github.com/coze-dev/coze-studio/backend/domain/{domain}/entity"
)

// {Domain}Service 定义{domain}领域服务接口
type {Domain}Service interface {
    // Create{Entity} 创建{entity}
    // 参数：
    //   ctx - 请求上下文
    //   req - 创建请求
    //
    // 返回：
    //   *entity.{Entity} - 创建的实体
    //   error - 错误信息
    //
    // 错误码：
    //   ErrCode{Entity}AlreadyExists - {entity}已存在
    Create{Entity}(ctx context.Context, req *Create{Entity}Request) (*entity.{Entity}, error)

    // Get{Entity}ByID 根据ID获取{entity}
    Get{Entity}ByID(ctx context.Context, id int64) (*entity.{Entity}, error)

    // Update{Entity} 更新{entity}
    Update{Entity}(ctx context.Context, req *Update{Entity}Request) error

    // Delete{Entity} 删除{entity}
    Delete{Entity}(ctx context.Context, id int64) error
}
```

### 3. 定义请求/响应结构

```go
// Create{Entity}Request 创建{entity}请求
type Create{Entity}Request struct {
    Name     string
    Email    string
    SpaceID  int64
    OwnerID  int64
}

// Update{Entity}Request 更新{entity}请求
type Update{Entity}Request struct {
    {Entity}ID int64
    Name       *string
    Email      *string
}

// Publish{Entity}Request 发布{entity}请求
type Publish{Entity}Request struct {
    {Entity}ID int64
    Version    string
}
```

### 4. 实现服务

**文件位置**: `domain/{domain}/service/service_impl.go`

```go
package service

import (
    "context"
    "fmt"

    "gorm.io/gorm"

    "github.com/coze-dev/coze-studio/backend/domain/{domain}/entity"
    "github.com/coze-dev/coze-studio/backend/domain/{domain}/repository"
    "github.com/coze-dev/coze-studio/backend/pkg/errorx"
    "github.com/coze-dev/coze-studio/backend/types/errno"
)

// Components 服务依赖组件
type Components struct {
    DB    *gorm.DB
    {Entity}Repo repository.{Entity}Repository
}

// NewService 创建服务实例
func NewService(components *Components) {Domain}Service {
    return &{domain}ServiceImpl{
        Components: components,
    }
}

type {domain}ServiceImpl struct {
    *Components
}

// Create{Entity} 创建{entity}
func (s *{domain}ServiceImpl) Create{Entity}(ctx context.Context, req *Create{Entity}Request) (*entity.{Entity}, error) {
    // 1. 参数验证
    if err := s.validateCreateRequest(ctx, req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. 构建领域对象
    {entity} := &entity.{Entity}{
        Name:     req.Name,
        Email:    req.Email,
        SpaceID:  req.SpaceID,
        OwnerID:  req.OwnerID,
        Status:   entity.{Entity}StatusActive,
    }

    // 3. 业务规则验证
    exists, err := s.{Entity}Repo.ExistsByEmail(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("check existence failed: %w", err)
    }
    if exists {
        return nil, errorx.New(errno.Err{Entity}AlreadyExists)
    }

    // 4. 持久化
    {entity}ID, err := s.{Entity}Repo.Create(ctx, {entity})
    if err != nil {
        return nil, errorx.Wrapf(err, "create {entity} failed, spaceID=%d", req.SpaceID)
    }

    {entity}.ID = {entity}ID

    // 5. 发布领域事件（如果需要）
    // s.eventPublisher.Publish(ctx, &{Entity}CreatedEvent{...})

    return {entity}, nil
}

// Get{Entity}ByID 获取{entity}
func (s *{domain}ServiceImpl) Get{Entity}ByID(ctx context.Context, id int64) (*entity.{Entity}, error) {
    {entity}, exist, err := s.{Entity}Repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if !exist {
        return nil, errorx.New(errno.Err{Entity}NotFound)
    }

    return {entity}, nil
}

// Update{Entity} 更新{entity}
func (s *{domain}ServiceImpl) Update{Entity}(ctx context.Context, req *Update{Entity}Request) error {
    // 1. 获取实体
    {entity}, exist, err := s.{Entity}Repo.FindByID(ctx, req.{Entity}ID)
    if err != nil {
        return err
    }
    if !exist {
        return errorx.New(errno.Err{Entity}NotFound)
    }

    // 2. 更新字段
    if req.Name != nil {
        {entity}.Name = *req.Name
    }
    if req.Email != nil {
        {entity}.Email = *req.Email
    }

    // 3. 业务验证
    if err := {entity}.Validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // 4. 持久化
    if err := s.{Entity}Repo.Update(ctx, {entity}); err != nil {
        return errorx.Wrapf(err, "update {entity} failed, id=%d", req.{Entity}ID)
    }

    return nil
}

// Delete{Entity} 删除{entity}
func (s *{domain}ServiceImpl) Delete{Entity}(ctx context.Context, id int64) error {
    if err := s.{Entity}Repo.Delete(ctx, id); err != nil {
        return errorx.Wrapf(err, "delete {entity} failed, id=%d", id)
    }
    return nil
}
```

### 5. 创建实体

**文件位置**: `domain/{domain}/entity/{entity}.go`

```go
package entity

import "time"

// {Entity} {entity}实体
type {Entity} struct {
    // 主键
    ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

    // 业务字段
    Name   string `gorm:"column:name;not null" json:"name"`
    Email  string `gorm:"column:email;uniqueIndex;not null" json:"email"`
    Status {Entity}Status `gorm:"column:status;not null" json:"status"`

    // 审计字段
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
    DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

    // 外键
    SpaceID  int64 `gorm:"column:space_id;index;not null" json:"space_id"`
    CreatorID int64 `gorm:"column:creator_id;index;not null" json:"creator_id"`

    // 版本（乐观锁）
    Version int `gorm:"column:version;default:0" json:"version"`
}

// {Entity}Status {entity}状态
type {Entity}Status string

const (
    {Entity}StatusActive   {Entity}Status = "active"
    {Entity}StatusInactive {Entity}Status = "inactive"
)

// TableName 指定表名
func ({Entity}) TableName() string {
    return "{entity}s"
}

// Validate 验证{entity}
func (e *{Entity}) Validate() error {
    if e.Name == "" {
        return ErrInvalidName
    }
    if !isValidEmail(e.Email) {
        return ErrInvalidEmail
    }
    return nil
}

// IsActive 判断是否激活
func (e *{Entity}) IsActive() bool {
    return e.Status == {Entity}StatusActive
}
```

### 6. 创建仓储接口

**文件位置**: `domain/{domain}/repository/interface.go`

```go
package repository

import (
    "context"
    "github.com/coze-dev/coze-studio/backend/domain/{domain}/entity"
)

// {Entity}Repository {entity}仓储接口
type {Entity}Repository interface {
    // 基础 CRUD
    Create(ctx context.Context, {entity} *entity.{Entity}) (int64, error)
    FindByID(ctx context.Context, id int64) (*entity.{Entity}, bool, error)
    Update(ctx context.Context, {entity} *entity.{Entity}) error
    Delete(ctx context.Context, id int64) error

    // 查询方法
    FindByEmail(ctx context.Context, email string) (*entity.{Entity}, bool, error)
    FindByStatus(ctx context.Context, status entity.{Entity}Status) ([]*entity.{Entity}, error)

    // 存在性检查
    ExistsByEmail(ctx context.Context, email string) (bool, error)

    // 分页查询
    FindWithPagination(ctx context.Context, params *PaginationParams) ([]*entity.{Entity}, int64, error)
}
```

### 7. 编写测试

**文件位置**: `domain/{domain}/service/service_test.go`

```go
package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/coze-dev/coze-studio/backend/domain/{domain}/entity"
    "github.com/coze-dev/coze-studio/backend/domain/{domain}/service"
)

func TestCreate{Entity}(t *testing.T) {
    tests := []struct {
        name    string
        req     *service.Create{Entity}Request
        want    *entity.{Entity}
        wantErr bool
        errCode int
    }{
        {
            name: "valid request",
            req: &service.Create{Entity}Request{
                Name:    "Test {Entity}",
                Email:   "test@example.com",
                SpaceID: 1,
                OwnerID: 1,
            },
            want: &entity.{Entity}{
                Name:  "Test {Entity}",
                Email: "test@example.com",
            },
            wantErr: false,
        },
        {
            name: "duplicate email",
            req: &service.Create{Entity}Request{
                Name:    "Test {Entity}",
                Email:   "duplicate@example.com",
                SpaceID: 1,
                OwnerID: 1,
            },
            wantErr: true,
            errCode: errno.Err{Entity}AlreadyExists,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            ctx := context.Background()
            s := setupTestService(t)

            // Act
            got, err := s.Create{Entity}(ctx, tt.req)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
                assert.Equal(t, tt.errCode, errorx.Code(err))
            } else {
                require.NoError(t, err)
                assert.NotNil(t, got)
                assert.Equal(t, tt.want.Name, got.Name)
                assert.Equal(t, tt.want.Email, got.Email)
            }
        })
    }
}
```

## 必须遵循的规范

### ✅ 命名规范
- **服务接口**: `{Domain}Service`（PascalCase）
- **服务实现**: `{domain}ServiceImpl`（小写开头）
- **方法名**: PascalCase + 动词前缀（Create/Get/Update/Delete）
- **请求结构**: `{Action}{Entity}Request`
- **错误变量**: `Err{Description}`

### ✅ DDD 分层规范
```
domain/{domain}/
├── entity/           # 实体
│   └── {entity}.go
├── repository/       # 仓储接口
│   ├── interface.go
│   └── repository.go
└── service/          # 领域服务
    ├── service.go           # 接口
    ├── service_impl.go      # 实现
    └── service_test.go      # 测试
```

### ✅ 错误处理规范
```go
// 1. 使用预定义错误
var (
    Err{Entity}NotFound = errors.New("{entity} not found")
    ErrInvalidEmail    = errors.New("invalid email")
)

// 2. 使用 errorx 包装错误
return nil, errorx.Wrapf(err, "operation failed, id=%d", id)

// 3. 返回错误码
return nil, errorx.New(errno.Err{Entity}NotFound)
```

### ✅ 方法实现模式
```go
func (s *serviceImpl) Method(ctx context.Context, req *Request) (*Response, error) {
    // 1. 参数验证
    if err := s.validateRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. 业务规则检查
    if exists, err := s.repo.Exists(ctx, req.ID); err != nil {
        return nil, err
    } else if exists {
        return nil, ErrAlreadyExists
    }

    // 3. 构建实体
    entity := &Entity{...}

    // 4. 持久化
    if err := s.repo.Create(ctx, entity); err != nil {
        return nil, errorx.Wrapf(err, "create failed")
    }

    // 5. 返回结果
    return entity, nil
}
```

## 输出检查清单

创建完成后，必须确保：

- [ ] 服务接口定义清晰
- [ ] 服务实现包含完整注释
- [ ] 方法遵循命名规范
- [ ] 错误处理完善
- [ ] 包含单元测试
- [ ] 通过 gofmt 检查
- [ ] 通过 golangci-lint 检查
- [ ] 遵循 DDD 分层架构
- [ ] Repository 接口定义完整
- [ ] 实体包含必要的方法

## 示例输出

**用户输入**：
```
创建用户管理服务，支持创建、查询、更新、删除用户
```

**AI 输出**：

1. **文件结构**：
```
domain/user/
├── entity/
│   └── user.go
├── repository/
│   ├── interface.go
│   └── repository.go
└── service/
    ├── service.go
    ├── service_impl.go
    └── service_test.go
```

2. **服务接口**：
```go
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    GetUserByID(ctx context.Context, id int64) (*User, error)
    UpdateUser(ctx context.Context, req *UpdateUserRequest) error
    DeleteUser(ctx context.Context, id int64) error
}
```

3. **服务实现**：
```go
func (s *userServiceImpl) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 1. 参数验证
    // 2. 业务规则检查
    // 3. 创建实体
    // 4. 持久化
    // 5. 返回结果
}
```

## 注意事项

1. **保持纯粹**：领域服务不应包含基础设施细节
2. **接口隔离**：接口应职责单一
3. **依赖注入**：通过构造函数注入依赖
4. **错误处理**：始终返回错误，不要忽略
5. **测试覆盖**：关键业务逻辑必须有测试
