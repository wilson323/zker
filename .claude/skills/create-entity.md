# Entity 开发 Skill

## 技能描述

创建符合 DDD 架构和 GORM 规范的领域实体（Entity）和数据模型（Model）。

## 适用场景

- 需要创建新的领域实体
- 需要定义数据模型和数据库映射
- 需要实现实体业务规则和状态转换
- 需要添加查询选项（Where Option）

## 实体设计原则

### 1. 实体分层

**Domain Entity (领域实体)**:
- 位置：`backend/domain/{domain}/entity/{entity}.go`
- 职责：表达业务概念，封装业务规则
- 特点：不包含 GORM 标签，纯领域逻辑

**Crossdomain Model (跨领域模型)**:
- 位置：`backend/crossdomain/{domain}/model/{model}.go`
- 职责：数据库映射，持久化关注
- 特点：包含 GORM 标签，数据库表映射

### 2. 实体命名规范

```
Domain Entity:
- 文件：{entity}.go (lowercase)
- 类型：{Entity} (PascalCase)
- 示例：user.go → type User struct

Crossdomain Model:
- 文件：{model}.go (lowercase)
- 类型：{Model} (PascalCase)
- 示例：plugin.go → type PluginInfo struct
```

### 3. 字段命名规范

```go
// ✅ 良好的字段命名
type User struct {
    ID          int64     // 主键
    Name        string    // 名称
    Email       string    // 邮箱
    SpaceID     int64     // 空间ID (外键)
    CreatorID   int64     // 创建者ID
    CreatedAt   int64     // 创建时间
    UpdatedAt   int64     // 更新时间
    DeletedAt   *int64    // 删除时间 (软删除，指针类型)
}

// ❌ 避免的字段命名
type User struct {
    user_id     int64     // 应使用驼峰命名
    UserName    string    // 应统一命名风格
    created_at  int64     // 应导出字段
}
```

## 实体模板

### 1. Domain Entity 模板

```go
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

package entity

import (
    "github.com/coze-dev/coze-studio/backend/crossdomain/{domain}/model"
)

// EntityName 实体功能描述
//
// 实体的职责、业务规则和状态转换。
//
// 状态转换：
//   - StatusA → StatusB: 触发条件
//   - StatusB → StatusC: 触发条件
//
// 业务规则：
//   - 规则1
//   - 规则2
type EntityName struct {
    *model.{ModelName}
}

// WhereEntityNameOption 查询选项
//
// 用于构建灵活的查询条件。
type WhereEntityNameOption struct {
    IDs         []int64
    SpaceID     *int64
    AppID       *int64
    UserID      *int64
    Status      []int32
    Name        *string  // 精确匹配
    Query       *string  // 模糊匹配
    Page        *int
    PageSize    *int
    Order       *Order
    OrderType   *OrderType
}

// Order 排序字段
type Order int32

const (
    OrderCreatedAt Order = 1  // 创建时间
    OrderUpdatedAt Order = 2  // 更新时间
    OrderName      Order = 3  // 名称
)

// OrderType 排序方向
type OrderType int32

const (
    OrderTypeAsc  OrderType = 1  // 升序
    OrderTypeDesc OrderType = 2  // 降序
)

// StatusType 状态类型
type StatusType int32

const (
    StatusTypeDraft     StatusType = 0  // 草稿
    StatusTypePublished StatusType = 1  // 已发布
    StatusTypeArchived  StatusType = 2  // 已归档
)

// GetBasic 获取基本信息
//
// 返回实体的基本信息，用于跨层传递。
func (e *EntityName) GetBasic() *EntityNameBasic {
    return &EntityNameBasic{
        ID:        e.ID,
        SpaceID:   e.SpaceID,
        Name:      e.Name,
        Status:    e.Status,
        CreatedAt: e.CreatedAt,
    }
}

// CanPublish 检查是否可以发布
//
// 验证实体是否满足发布条件。
func (e *EntityName) CanPublish() error {
    if e.Status != StatusTypeDraft {
        return errors.New("only draft can be published")
    }
    if e.Name == "" {
        return errors.New("name is required")
    }
    return nil
}

// CanDelete 检查是否可以删除
//
// 验证实体是否满足删除条件。
func (e *EntityName) CanDelete() error {
    if e.Status == StatusTypePublished {
        return errors.New("published entity cannot be deleted")
    }
    return nil
}

// EntityNameBasic 实体基本信息
type EntityNameBasic struct {
    ID        int64
    SpaceID   int64
    Name      string
    Status    StatusType
    CreatedAt int64
}
```

### 2. Crossdomain Model 模板

```go
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

package model

import "time"

// EntityName 数据库实体
//
// 包含数据库映射标签和业务字段。
//
// 表名：entity_names
// 索引：
//   - idx_space_id: 空间ID索引
//   - idx_space_id_status: 空间和状态联合索引
//   - idx_name: 名称索引
type EntityName struct {
    // ID 主键
    ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

    // SpaceID 空间ID
    SpaceID int64 `gorm:"column:space_id;not null;index:idx_space_id,priority:1" json:"space_id"`

    // AppID 应用ID (可选)
    AppID *int64 `gorm:"column:app_id;index:idx_app_id" json:"app_id,omitempty"`

    // Name 名称
    Name string `gorm:"column:name;type:varchar(255);not null" json:"name"`

    // Description 描述
    Description string `gorm:"column:description;type:text" json:"description"`

    // IconURI 图标URI
    IconURI string `gorm:"column:icon_uri;type:varchar(512)" json:"icon_uri"`

    // IconURL 图标URL
    IconURL string `gorm:"column:icon_url;type:varchar(512)" json:"icon_url"`

    // Status 状态
    Status Status `gorm:"column:status;not null;default:0;index:idx_space_id_status,priority:2" json:"status"`

    // CreatorID 创建者ID
    CreatorID int64 `gorm:"column:creator_id;not null" json:"creator_id"`

    // CreatedAt 创建时间 (毫秒时间戳)
    CreatedAt int64 `gorm:"column:created_at;not null;index:idx_space_id_status,priority:3;autoCreateTime:milli" json:"created_at"`

    // UpdatedAt 更新时间 (毫秒时间戳)
    UpdatedAt int64 `gorm:"column:updated_at;not null;index:idx_created_at" json:"updated_at"`

    // DeletedAt 删除时间 (软删除)
    DeletedAt *int64 `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// Status 状态枚举
type Status int32

const (
    // StatusDraft 草稿
    StatusDraft Status = 0
    // StatusPublished 已发布
    StatusPublished Status = 1
    // StatusArchived 已归档
    StatusArchived Status = 2
)

// TableName 指定表名
func (EntityName) TableName() string {
    return "entity_names"
}

// BeforeCreate GORM 钩子：创建前
func (e *EntityName) BeforeCreate(tx *gorm.DB) error {
    // 设置默认值
    if e.Status == 0 {
        e.Status = StatusDraft
    }
    return nil
}

// BeforeUpdate GORM 钩子：更新前
func (e *EntityName) BeforeUpdate(tx *gorm.DB) error {
    // 更新时间戳
    e.UpdatedAt = time.Now().UnixMilli()
    return nil
}
```

## GORM 标签规范

### 1. 基础标签

```go
// 主键
ID int64 `gorm:"column:id;primaryKey;autoIncrement"`

// 必填字段
Name string `gorm:"column:name;not null"`

// 可选字段 (使用指针)
Description *string `gorm:"column:description"`

// 默认值
Status int32 `gorm:"column:status;not null;default:0"`

// 唯一索引
Email string `gorm:"column:email;uniqueIndex:idx_email"`

// 普通索引
SpaceID int64 `gorm:"column:space_id;index:idx_space_id"`

// 联合索引
SpaceID int64 `gorm:"column:space_id;index:idx_space_status,priority:1"`
Status   int32 `gorm:"column:status;index:idx_space_status,priority:2"`
```

### 2. 类型标签

```go
// 字符串类型
Name      string `gorm:"column:name;type:varchar(255)"`
Content   string `gorm:"column:content;type:text"`

// 数字类型
Count     int32  `gorm:"column:count;type:int"`
Amount    float64 `gorm:"column:amount;type:decimal(10,2)"`

// 时间类型
CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime:milli"`

// JSON 类型
Metadata  string `gorm:"column:metadata;type:json"`
```

### 3. 软删除

```go
// 使用指针类型的 DeletedAt
DeletedAt *int64 `gorm:"column:deleted_at;index"`

// 查询时自动过滤已删除记录
db.Where("space_id = ?", spaceID).Find(&entities)
```

## 实体方法规范

### 1. 构造方法

```go
// NewEntityName 创建新实体
func NewEntityName(spaceID, creatorID int64, name string) *EntityName {
    return &EntityName{
        SpaceID:   spaceID,
        CreatorID: creatorID,
        Name:      name,
        Status:    StatusDraft,
    }
}
```

### 2. 业务方法

```go
// Publish 发布实体
func (e *EntityName) Publish() error {
    if err := e.CanPublish(); err != nil {
        return err
    }
    e.Status = StatusPublished
    return nil
}

// Archive 归档实体
func (e *EntityName) Archive() error {
    if e.Status != StatusPublished {
        return errors.New("only published entity can be archived")
    }
    e.Status = StatusArchived
    return nil
}

// UpdateName 更新名称
func (e *EntityName) UpdateName(name string) error {
    if name == "" {
        return errors.New("name cannot be empty")
    }
    e.Name = name
    return nil
}
```

### 3. 查询方法

```go
// GetStatusText 获取状态文本
func (e *EntityName) GetStatusText() string {
    switch e.Status {
    case StatusDraft:
        return "Draft"
    case StatusPublished:
        return "Published"
    case StatusArchived:
        return "Archived"
    default:
        return "Unknown"
    }
}

// IsPublished 是否已发布
func (e *EntityName) IsPublished() bool {
    return e.Status == StatusPublished
}

// IsDeleted 是否已删除
func (e *EntityName) IsDeleted() bool {
    return e.DeletedAt != nil
}
```

## 值对象（Value Object）

```go
// IDVersionPair ID-版本对
type IDVersionPair struct {
    ID      int64
    Version string
}

// EntityReference 实体引用
type EntityReference struct {
    ID       int64
    Name     string
    IconURL  string
    SpaceID  int64
}

// ValidateIssue 验证问题
type ValidateIssue struct {
    Code    string
    Message string
    Level   SeverityLevel
}

type SeverityLevel int32

const (
    SeverityLevelInfo    SeverityLevel = 0
    SeverityLevelWarning SeverityLevel = 1
    SeverityLevelError   SeverityLevel = 2
)
```

## 实体关系

### 1. 一对一

```go
type User struct {
    ID       int64
    Profile  *Profile  // 一对一
}

type Profile struct {
    ID     int64
    UserID int64  // 外键
    User   *User  // 反向引用
}
```

### 2. 一对多

```go
type Space struct {
    ID    int64
    Name  string
    Apps  []*App  // 一对多
}

type App struct {
    ID       int64
    SpaceID  int64   // 外键
    Space    *Space  // 反向引用
}
```

### 3. 多对多

```go
type App struct {
    ID     int64
    Tags   []*Tag         `gorm:"many2many:app_tags;"`
}

type Tag struct {
    ID    int64
    Apps  []*App         `gorm:"many2many:app_tags;"`
}
```

## 生成流程

### 步骤 1: 分析需求

1. 识别业务概念和属性
2. 确定实体职责
3. 定义状态和状态转换
4. 识别实体关系

### 步骤 2: 设计 Model

1. 定义数据库表结构
2. 添加 GORM 标签
3. 定义索引
4. 实现 GORM 钩子

### 步骤 3: 设计 Entity

1. 组合 Model
2. 定义查询选项
3. 实现业务方法
4. 添加验证逻辑

### 步骤 4: 编写测试

1. 测试实体创建
2. 测试状态转换
3. 测试业务规则
4. 测试边界条件

## 输出检查清单

创建实体后，确保：

- [ ] 文件命名符合规范
- [ ] 类型命名符合规范 (PascalCase)
- [ ] 字段命名使用驼峰式
- [ ] 主键正确定义
- [ ] 外键正确定义
- [ ] 索引合理设置
- [ ] GORM 标签完整
- [ ] 软删除正确实现
- [ ] 状态常量定义
- [ ] 业务方法实现
- [ ] 验证逻辑完整
- [ ] TableName 方法实现
- [ ] GORM 钩子实现
- [ ] 完整的注释文档

## 真实代码示例

### 示例 1: User Entity

```go
// backend/domain/user/entity/user.go
package entity

type User struct {
    UserID      int64
    Name        string  // nickname
    UniqueName  string  // unique name
    Email       string  // email
    Description string  // user description
    IconURI     string  // avatar URI
    IconURL     string  // avatar URL
    UserVerified bool   // Is the user authenticated?
    Locale      string
    SessionKey  string  // session key
    CreatedAt   int64   // creation time
    UpdatedAt   int64   // update time
}

type UserLevel int32

const (
    UserLevelFree    UserLevel = 0
    UserLevelPro     UserLevel = 1
    UserLevelPremium UserLevel = 2
)

type UserBenefit struct {
    UserID         int64
    UserLevel      UserLevel
    UsedCount      int32
    TotalCount     int32
    IsUnlimited    bool
    ResetDatetime  int64
    CallQPS        int32
}
```

### 示例 2: Space Entity

```go
// backend/domain/user/entity/space.go
package entity

type SpaceType int32

const (
    SpaceTypePersonal SpaceType = 1
    SpaceTypeTeam     SpaceType = 2
)

type Space struct {
    ID          int64
    Name        string
    Description string
    IconURL     string
    SpaceType   SpaceType
    OwnerID     int64
    CreatorID   int64
    CreatedAt   int64
    UpdatedAt   int64
}

// IsPersonal 是否为个人空间
func (s *Space) IsPersonal() bool {
    return s.SpaceType == SpaceTypePersonal
}

// IsTeam 是否为团队空间
func (s *Space) IsTeam() bool {
    return s.SpaceType == SpaceTypeTeam
}
```

### 示例 3: Knowledge Entity

```go
// backend/domain/knowledge/entity/knowledge.go
package entity

import model "github.com/coze-dev/coze-studio/backend/crossdomain/knowledge/model"

type Knowledge struct {
    *model.Knowledge
}

type WhereKnowledgeOption struct {
    KnowledgeIDs []int64
    AppID        *int64
    SpaceID      *int64
    Name         *string  // Exact match
    Status       []int32
    UserID       *int64
    Query        *string  // fuzzy match
    Page         *int
    PageSize     *int
    Order        *Order
    OrderType    *OrderType
    FormatType   *int64
}

type OrderType int32

const (
    OrderTypeAsc  OrderType = 1
    OrderTypeDesc OrderType = 2
)

type Order int32

const (
    OrderCreatedAt Order = 1
    OrderUpdatedAt Order = 2
)
```

## 注意事项

1. **职责单一**: 实体只负责封装业务规则，不涉及持久化细节
2. **不可变性**: 尽量保持实体的不可变性，通过方法修改状态
3. **验证前置**: 在状态转换前进行验证，拒绝非法状态
4. **类型安全**: 使用类型常量而非魔法数字
5. **关系清晰**: 明确定义实体间的关系，避免循环依赖
6. **测试覆盖**: 为所有业务方法和状态转换编写测试
