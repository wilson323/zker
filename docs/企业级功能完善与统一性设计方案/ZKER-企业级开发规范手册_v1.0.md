# ZKER 企业级开发规范手册

**文档版本**: v1.0
**创建日期**: 2025-01-01
**最后更新**: 2025-01-01
**适用范围**: 所有 ZKER 项目开发人员

---

## 📋 文档概述

### 规范目标

本开发规范手册确保所有开发人员：

**核心目标**:
- ✅ **代码一致性**: 全局统一的代码风格和结构
- ✅ **高质量标准**: 企业级代码质量要求
- ✅ **可维护性**: 清晰的代码结构和注释
- ✅ **可测试性**: 完善的单元测试和集成测试
- ✅ **文档化**: 代码即文档，完善的 API 文档

### 适用人员

**必须遵守规范的人员**:
- 后端开发工程师（Go）
- 前端开发工程师（React + TypeScript）
- 数据库管理员（DBA）
- 测试工程师
- DevOps 工程师

**规范执行机制**:
- ✅ **自动化检查**: ESLint、golangci-lint、Pre-commit Hooks
- ✅ **CI/CD 验证**: 自动运行代码质量检查
- ✅ **代码审查**: PR 必须通过规范审查
- ✅ **定期评审**: 每季度评审并更新规范

---

## 🎯 一、全局开发原则

### 1.1 KISS 原则（Keep It Simple, Stupid）

**原则说明**: 保持代码和设计的极致简洁

**实践指南**:
```go
// ❌ 不好的写法：过度复杂的抽象
type AbstractFactoryFactory interface {
    CreateFactoryFactory() FactoryFactory
}

// ✅ 好的写法：简单直接
type BotService interface {
    CreateBot(ctx context.Context, bot *Bot) (*Bot, error)
}
```

**检查清单**:
- [ ] 函数长度不超过 50 行
- [ ] 嵌套层级不超过 3 层
- [ ] 避免过度设计的抽象
- [ ] 优先选择最直观的解决方案

### 1.2 DRY 原则（Don't Repeat Yourself）

**原则说明**: 避免代码重复，提取公共逻辑

**实践指南**:
```go
// ❌ 不好的写法：重复的数据库查询
func GetBotByID(ctx context.Context, botID string) (*Bot, error) {
    var bot Bot
    err := db.Where("bot_id = ?", botID).First(&bot).Error
    return &bot, err
}

func GetBotByName(ctx context.Context, name string) (*Bot, error) {
    var bot Bot
    err := db.Where("name = ?", name).First(&bot).Error
    return &bot, err
}

// ✅ 好的写法：提取公共逻辑
func getBotByCondition(ctx context.Context, condition interface{}) (*Bot, error) {
    var bot Bot
    err := db.Where(condition).First(&bot).Error
    return &bot, err
}

func GetBotByID(ctx context.Context, botID string) (*Bot, error) {
    return getBotByCondition(ctx, map[string]interface{}{"bot_id": botID})
}
```

**检查清单**:
- [ ] 相同逻辑出现 3 次以上必须提取
- [ ] 复制粘贴代码必须有注释说明原因
- [ ] 使用工具检测重复代码（如 SonarQube）

### 1.3 YAGNI 原则（You Aren't Gonna Need It）

**原则说明**: 只实现当前明确所需的功能

**实践指南**:
```go
// ❌ 不好的写法：预留大量未使用的字段
type Bot struct {
    BotID       string
    Name        string
    // 以下字段预留但从未使用
    FutureField1 string
    FutureField2 int
    FutureField3 bool
}

// ✅ 好的写法：只定义当前需要的字段
type Bot struct {
    BotID    string
    Name     string
    Status   string
}
```

**检查清单**:
- [ ] 不为"可能以后会用到"的功能开发
- [ ] 不过度设计扩展点
- [ ] 及时删除未使用的代码和依赖

### 1.4 SOLID 原则

**单一职责原则 (SRP)**:
```go
// ✅ 每个模块只负责一件事
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, userID string) (*User, error)
    Update(ctx context.Context, user *User) error
}

type BotRepository interface {
    Create(ctx context.Context, bot *Bot) error
    FindByID(ctx context.Context, botID string) (*Bot, error)
    Update(ctx context.Context, bot *Bot) error
}
```

**开闭原则 (OCP)**:
```go
// ✅ 对扩展开放，对修改关闭
type PermissionChecker interface {
    Check(ctx context.Context, req *CheckRequest) (*CheckResult, error)
}

// 可以添加新的权限检查器，无需修改现有代码
type RBACChecker struct{}
type ABACChecker struct{}
```

**依赖倒置原则 (DIP)**:
```go
// ✅ 依赖抽象而非具体实现
type BotService struct {
    repo BotRepository  // 依赖接口，不依赖具体实现
}

func NewBotService(repo BotRepository) *BotService {
    return &BotService{repo: repo}
}
```

---

## 💻 二、后端开发规范（Go）

### 2.1 项目结构规范

#### DDD 分层结构

```
backend/
├── api/                    # API 层（接口层）
│   └── handler/           # 请求处理器
├── application/           # 应用层（业务编排）
│   └── app/
├── domain/                # 领域层（核心业务逻辑）
│   ├── {aggregate}/       # 聚合根
│   │   ├── entity/       # 实体
│   │   ├── service/      # 领域服务
│   │   └── repository/   # 仓储接口
│   └── internal/
│       └── dal/          # 数据访问层实现
├── crossdomain/           # 跨域层（领域间共享）
│   └── {domain}/
│       └── model/
├── infra/                # 基础设施层（技术实现）
│   ├── cache/
│   ├── eventbus/
│   └── orm/
└── bizpkg/               # 业务包
```

#### 包命名规范

```go
// ✅ 好的包命名：小写，单数，描述性
package permission       // 权限领域
package rbac_permission  // RBAC 权限服务
package middleware      // 中间件

// ❌ 不好的包命名
package permissions    // 复数
package rbacPerm        // 缩写
package util            // 过于通用
```

#### 文件命名规范

```
// ✅ 好的文件命名
user_repository.go        // 仓储实现
user_repository_test.go   // 仓储测试
rbac_permission.go        // 领域服务
rbac_permission_test.go   // 领域服务测试

// ❌ 不好的文件命名
user.go                   // 过于通用
UserRepo.go              // 大驼峰（Go 文件名应小写）
rbac.go                  // 过于简短
```

### 2.2 代码风格规范

#### 命名规范

**常量命名**:
```go
// ✅ 好的写法：大驼峰或全大写 + 下划线
const MaxRetries = 3
const DEFAULT_TIMEOUT = 30 * time.Second
const (
    StatusDraft    = "draft"
    StatusPublished = "published"
)

// ❌ 不好的写法
const maxRetries = 3      // 小驼峰
const default_timeout = 30 // 驼峰
```

**变量命名**:
```go
// ✅ 好的写法：大驼峰（导出）或小驼峰（私有）
var UserService *userService
type BotService struct {}
func GetUserByID() {}

// ❌ 不好的写法
var user_service *UserService  // 下划线
type bot_service struct {}      // 下划线
func get_user_by_id() {}       // 下划线
```

**接口命名**:
```go
// ✅ 好的写法：大驼峰 + er 后缀
type UserRepository interface {}
type BotService interface {}
type CacheManager interface {}

// ❌ 不好的写法
type IUserInterface interface {}  // 过度前缀
type repo interface {}             // 缩写
```

#### 函数设计规范

**函数长度**:
```go
// ✅ 好的写法：函数职责单一，长度 < 50 行
func (s *BotService) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error) {
    // 1. 参数验证
    if err := s.validateCreateRequest(req); err != nil {
        return nil, err
    }

    // 2. 业务逻辑
    bot := s.buildBotFromRequest(req)
    if err := s.repo.Create(ctx, bot); err != nil {
        return nil, err
    }

    // 3. 返回结果
    return bot, nil
}

// ❌ 不好的写法：函数过长，逻辑混杂
func (s *BotService) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error) {
    // 100+ 行代码，包含验证、逻辑、持久化、通知、日志...
}
```

**参数数量**:
```go
// ✅ 好的写法：参数 < 5 个，使用结构体
func CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {}

// ❌ 不好的写法：参数过多
func CreateUser(
    ctx context.Context,
    username string,
    email string,
    password string,
    nickname string,
    avatar string,
    status string,
) (*User, error) {}
```

**返回值规范**:
```go
// ✅ 好的写法：返回值 + error
func GetUserByID(ctx context.Context, userID string) (*User, error) {
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    return user, nil
}

// ❌ 不好的写法：忽略错误
func GetUserByID(ctx context.Context, userID string) *User {
    user, _ := s.repo.FindByID(ctx, userID)
    return user
}
```

#### 错误处理规范

**错误定义**:
```go
// ✅ 好的写法：使用自定义错误码
var (
    ErrUserNotFound = errno.NewError(
        "USER404",
        "User not found",
        "用户不存在",
        404,
    )

    ErrInvalidEmail = errno.NewError(
        "USER400",
        "Invalid email format",
        "邮箱格式无效",
        400,
    )
)

// 使用错误码
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    return user, nil
}
```

**错误传递**:
```go
// ✅ 好的写法：使用 fmt.Errorf 和 %w 保留错误链
func (s *BotService) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error) {
    if err := s.repo.Create(ctx, bot); err != nil {
        return nil, fmt.Errorf("failed to create bot: %w", err)
    }
    return bot, nil
}

// ❌ 不好的写法：忽略原始错误
func (s *BotService) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error) {
    if err := s.repo.Create(ctx, bot); err != nil {
        return nil, errors.New("create bot failed")
    }
}
```

### 2.3 并发安全规范

#### Goroutine 使用

**✅ 好的实践**:
```go
// 1. 传递 context
func (s *BotService) ProcessAsync(ctx context.Context, botID string) error {
    go func() {
        // 使用 context 传递取消信号
        select {
        case <-ctx.Done():
            log.Printf("Context cancelled: %v", ctx.Err())
            return
        default:
            s.processBot(ctx, botID)
        }
    }()
    return nil
}

// 2. WaitGroup 等待
func (s *BotService) ProcessBatch(ctx context.Context, botIDs []string) error {
    var wg sync.WaitGroup
    for _, botID := range botIDs {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            s.processBot(ctx, id)
        }(botID)
    }
    wg.Wait()
    return nil
}

// 3. 限制并发数
func (s *BotService) ProcessBatch(ctx context.Context, botIDs []string) error {
    sem := make(chan struct{}, 10) // 最多 10 个并发
    var wg sync.WaitGroup

    for _, botID := range botIDs {
        wg.Add(1)
        sem <- struct{}{} // 获取信号量
        go func(id string) {
            defer wg.Done()
            defer func() { <-sem }() // 释放信号量
            s.processBot(ctx, id)
        }(botID)
    }

    wg.Wait()
    return nil
}
```

#### Channel 使用

**✅ 好的实践**:
```go
// 1. 关闭 channel
func producer(ctx context.Context) chan int {
    ch := make(chan int, 10)
    go func() {
        defer close(ch)
        for i := 0; i < 100; i++ {
            select {
            case ch <- i:
            case <-ctx.Done():
                return
            }
        }
    }()
    return ch
}

// 2. 读取 channel
func consumer(ch <-chan int) {
    for val := range ch {
        fmt.Println(val)
    }
}

// 3. 防止 goroutine 泄漏
func (s *BotService) StartWorkers(ctx context.Context) {
    for i := 0; i < 10; i++ {
        go func() {
            for {
                select {
                case <-ctx.Done():
                    return  // context 取消时退出
                case task := <-s.taskQueue:
                    s.processTask(task)
                }
            }
        }()
    }
}
```

### 2.4 数据库操作规范

#### GORM 使用规范

**✅ 好的实践**:
```go
// 1. 使用事务
func (s *BotService) CreateBotWithConfig(ctx context.Context, bot *Bot, config *BotConfig) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 创建 Bot
        if err := tx.Create(bot).Error; err != nil {
            return err
        }

        // 创建配置
        config.BotID = bot.BotID
        if err := tx.Create(config).Error; err != nil {
            return err
        }

        return nil
    })
}

// 2. 使用 Preload 预加载
func (s *BotService) GetBotWithConfig(ctx context.Context, botID string) (*Bot, error) {
    var bot Bot
    err := s.db.
        Preload("Config").
        Preload("Creator").
        Where("bot_id = ?", botID).
        First(&bot).Error

    if err != nil {
        return nil, err
    }
    return &bot, nil
}

// 3. 使用分页
func (s *BotService) ListBots(ctx context.Context, page, pageSize int) ([]*Bot, int64, error) {
    var bots []*Bot
    var total int64

    if err := s.db.Model(&Bot{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * pageSize
    if err := s.db.
        Offset(offset).
        Limit(pageSize).
        Find(&bots).Error; err != nil {
        return nil, 0, err
    }

    return bots, total, nil
}

// 4. 使用索引查询
func (s *BotService) GetBotByTenant(ctx context.Context, tenantID string, botID string) (*Bot, error) {
    var bot Bot
    // 确保 WHERE 条件使用索引字段
    err := s.db.
        Where("tenant_id = ? AND bot_id = ?", tenantID, botID).
        First(&bot).Error

    if err != nil {
        return nil, err
    }
    return &bot, nil
}
```

#### SQL 注入防护

```go
// ✅ 好的写法：使用参数化查询
func (s *BotService) SearchBots(ctx context.Context, keyword string) ([]*Bot, error) {
    var bots []*Bot
    // GORM 自动参数化
    err := s.db.
        Where("name LIKE ?", "%"+keyword+"%").
        Find(&bots).Error
    return bots, err
}

// ❌ 不好的写法：字符串拼接（SQL 注入风险）
func (s *BotService) SearchBots(ctx context.Context, keyword string) ([]*Bot, error) {
    var bots []*Bot
    err := s.db.
        Where("name LIKE '%" + keyword + "%'").  // 危险！
        Find(&bots).Error
    return bots, err
}
```

### 2.5 API 接口规范

#### Hertz Handler 规范

**✅ 好的实践**:
```go
// api/handler/bot_handler.go
package handler

import (
    "context"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

type BotHandler struct {
    botService *service.BotService
}

func NewBotHandler(botService *service.BotService) *BotHandler {
    return &BotHandler{botService: botService}
}

// CreateBot 创建 Bot
// @router /api/v1/bots [POST]
func (h *BotHandler) CreateBot(ctx context.Context, c *app.RequestContext) {
    // 1. 参数绑定
    var req CreateBotRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(consts.StatusBadRequest, map[string]interface{}{
            "code":       "BOT400",
            "message":    "Invalid request parameters",
            "message_zh": "请求参数无效",
            "request_id": getRequestID(ctx),
        })
        return
    }

    // 2. 调用服务
    bot, err := h.botService.CreateBot(ctx, &req)
    if err != nil {
        // 错误已在 service 层转换为错误码
        errno := err.(errno.ErrorCode)
        c.JSON(errno.HTTPStatus(), map[string]interface{}{
            "code":       errno.Code(),
            "message":    errno.Message(),
            "message_zh": errno.MessageZH(),
            "request_id": getRequestID(ctx),
        })
        return
    }

    // 3. 返回成功响应
    c.JSON(consts.StatusCreated, map[string]interface{}{
        "code":       "SUCCESS",
        "message":    "Bot created successfully",
        "message_zh": "Bot 创建成功",
        "data": map[string]interface{}{
            "bot_id": bot.BotID,
            "name":   bot.Name,
        },
        "request_id": getRequestID(ctx),
    })
}

// ListBots 获取 Bot 列表
// @router /api/v1/bots [GET]
func (h *BotHandler) ListBots(ctx context.Context, c *app.RequestContext) {
    // 1. 参数解析
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("page_size", "20")

    // 2. 调用服务
    bots, total, err := h.botService.ListBots(ctx, page, pageSize)
    if err != nil {
        // 错误处理...
    }

    // 3. 返回响应
    c.JSON(consts.StatusOK, map[string]interface{}{
        "code":       "SUCCESS",
        "message":    "Success",
        "message_zh": "成功",
        "data": map[string]interface{}{
            "bots":       bots,
            "total":      total,
            "page":       page,
            "page_size":  pageSize,
        },
        "request_id": getRequestID(ctx),
    })
}

func getRequestID(ctx context.Context) string {
    // 从 context 或 header 获取 request_id
    if reqID := ctx.Value("request_id"); reqID != nil {
        return reqID.(string)
    }
    return generateRequestID()
}
```

### 2.6 测试规范

#### 单元测试

**✅ 好的实践**:
```go
// domain/permission/service/rbac_permission_test.go
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock 仓储
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, userID string) (*User, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, errors.New("not found")
    }
    return args.Get(0).(*User), args.Error(1)
}

// 表格驱动测试
func TestRBACPermissionService_CheckPermission(t *testing.T) {
    tests := []struct {
        name           string
        userID         string
        resourceID     string
        action         string
        expectedResult bool
        expectError    bool
    }{
        {
            name:           "Owner has full permission",
            userID:         "user1",
            resourceID:     "bot1",
            action:         "read",
            expectedResult: true,
            expectError:    false,
        },
        {
            name:           "Member has no permission",
            userID:         "user2",
            resourceID:     "bot1",
            action:         "delete",
            expectedResult: false,
            expectError:    false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            service := setupTestService()
            ctx := context.Background()

            // Act
            allowed, err := service.CheckPermission(ctx, tt.userID, tt.resourceID, tt.action)

            // Assert
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedResult, allowed)
            }
        })
    }
}

// 测试覆盖率要求：> 80%
```

#### 集成测试

```go
// test/integration/api/bot_api_test.go
package integration

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "encoding/json"
)

func TestBotAPI_CreateBot(t *testing.T) {
    // 1. 启动测试服务器
    router := setupTestRouter()
    server := httptest.NewServer(router)
    defer server.Close()

    // 2. 构造请求
    reqBody := CreateBotRequest{
        Name:        "Test Bot",
        Description: "Test Description",
    }
    body, _ := json.Marshal(reqBody)

    req, _ := http.NewRequest("POST", server.URL+"/api/v1/bots", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+getTestToken())

    // 3. 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    assert.NoError(t, err)

    // 4. 验证响应
    assert.Equal(t, 201, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    assert.Equal(t, "SUCCESS", result["code"])
}
```

---

## 🎨 三、前端开发规范（React + TypeScript）

### 3.1 项目结构规范

#### 组件目录结构

```
packages/
├── arch/                      # 架构层（基础设施）
│   └── bot-components/       # Bot 相关组件
│       ├── index.ts          # 导出入口
│       ├── BotCard/          # 组件目录
│       │   ├── index.tsx     # 组件实现
│       │   ├── index.mdx     # 组件文档
│       │   ├── style.ts      # 组件样式
│       │   └── __tests__/    # 测试文件
│       └── types.ts          # 类型定义
│
├── common/                    # 通用组件层
│   ├── ui-components/        # 基础 UI 组件
│   └── biz-components/       # 业务组件
│
└── foundation/               # 基础包
    └── account-*/            # 账户相关
```

#### 组件命名规范

```tsx
// ✅ 好的组件命名：大驼峰
export const Button: React.FC<ButtonProps> = () => {}
export const MessageBox: React.FC<MessageBoxProps> = () => {}

// ✅ 文件命名：组件同名
// MessageBox/index.tsx
// MessageBox/style.ts
// MessageBox/__tests__/MessageBox.test.tsx

// ❌ 不好的组件命名
export const button: React.FC = () => {}  // 小驼峰
export const message_box: React.FC = () => {}  // 下划线
```

### 3.2 TypeScript 规范

#### 类型定义

**✅ 好的实践**:
```typescript
// 1. 接口定义：大驼峰 + Props 后缀
interface ButtonProps extends SemiButtonProps {
  type?: 'primary' | 'secondary' | 'tertiary';
  size?: 'small' | 'medium' | 'large';
  disabled?: boolean;
  loading?: boolean;
  onClick?: (event: React.MouseEvent) => void;
  children?: React.ReactNode;
}

// 2. 类型导入
import type { ButtonProps } from './types';

// 3. 泛型使用
interface ApiResponse<T> {
  code: string;
  message: string;
  data: T;
}

function fetchData<T>(url: string): Promise<ApiResponse<T>> {
  // ...
}

// 4. 枚举定义
enum BotStatus {
  Draft = 'draft',
  Published = 'published',
  Archived = 'archived',
}

// 5. 联合类型
type BotState = {
  status: BotStatus;
  lastPublishedAt: Date;
} & (
  | { status: BotStatus.Draft; draftId: string }
  | { status: BotStatus.Published; publishedId: string }
  | { status: BotStatus.Archived; archivedAt: Date }
);
```

#### 类型导出

```typescript
// ✅ 好的写法：统一导出类型
// types/index.ts
export type { ButtonProps } from './button';
export type { InputProps } from './input';
export type { ModalProps } from './modal';
export type { SelectProps } from './select';

// 使用时导入
import type { ButtonProps, InputProps } from '@arch/ui-components';
```

### 3.3 React 组件规范

#### 函数组件

```tsx
// ✅ 好的写法：使用函数组件 + Hooks
import React, { useState, useEffect } from 'react';

interface MessageBoxProps {
  message: {
    id: string;
    role: 'user' | 'assistant';
    content: string;
  };
}

export const MessageBox: React.FC<MessageBoxProps> = ({ message }) => {
  const [isExpanded, setIsExpanded] = useState(false);

  useEffect(() => {
    // 组件挂载时执行
    console.log('Message mounted:', message.id);
  }, [message.id]);

  return (
    <div className="message-box">
      <div className="message-role">{message.role}</div>
      <div className="message-content">{message.content}</div>
    </div>
  );
};

// ❌ 不好的写法：使用类组件（不推荐）
class MessageBox extends React.Component {
  // ...
}
```

#### Props 解构

```tsx
// ✅ 好的写法：直接解构 props
export const Button: React.FC<ButtonProps> = ({
  type = 'primary',
  size = 'medium',
  disabled = false,
  loading = false,
  onClick,
  children,
}) => {
  return (
    <SemiButton
      type={type}
      size={size}
      disabled={disabled || loading}
      onClick={onClick}
    >
      {loading ? '加载中...' : children}
    </SemiButton>
  );
};
```

#### 状态管理

```typescript
// ✅ 使用 Zustand
import create from 'zustand';

interface BotStore {
  bots: Bot[];
  loading: boolean;
  fetchBots: () => Promise<void>;
  createBot: (bot: Partial<Bot>) => Promise<void>;
}

export const useBotStore = create<BotStore>((set, get) => ({
  bots: [],
  loading: false,
  fetchBots: async () => {
    set({ loading: true });
    try {
      const response = await botApi.list();
      set({ bots: response.data.bots, loading: false });
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },
  createBot: async (bot) => {
    const response = await botApi.create(bot);
    set((state) => ({ bots: [...state.bots, response.data] }));
  },
}));

// 组件中使用
export const BotList: React.FC = () => {
  const { bots, loading, fetchBots } = useBotStore();

  useEffect(() => {
    fetchBots();
  }, [fetchBots]);

  if (loading) {
    return <Spin />;
  }

  return <div>{bots.map(bot => <BotCard key={bot.id} bot={bot} />)}</div>;
};
```

#### 事件处理

```tsx
// ✅ 好的写法：事件处理器命名以 handle 开头
export const BotForm: React.FC = () => {
  const [formData, setFormData] = useState<BotFormData>({
    name: '',
    description: '',
  });

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    try {
      await botApi.create(formData);
      Toast.success('Bot 创建成功');
    } catch (error) {
      Toast.error('Bot 创建失败');
    }
  };

  const handleNameChange = (value: string) => {
    setFormData(prev => ({ ...prev, name: value }));
  };

  return (
    <form onSubmit={handleSubmit}>
      <Input
        value={formData.name}
        onChange={handleNameChange}
        placeholder="请输入 Bot 名称"
      />
      <Button type="primary" htmlType="submit">提交</Button>
    </form>
  );
};
```

### 3.4 样式规范

#### Tailwind CSS 使用

```tsx
// ✅ 好的写法：优先使用 Tailwind 类名
export const Card: React.FC = () => {
  return (
    <div className="p-4 bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow">
      <h2 className="text-xl font-semibold text-gray-900 mb-2">
        卡片标题
      </h2>
      <p className="text-gray-600">
        卡片内容
      </p>
    </div>
  );
};

// ✅ 响应式设计
export const ResponsiveCard: React.FC = () => {
  return (
    <div className="p-4 md:p-6 lg:p-8 bg-white rounded-lg">
      <h2 className="text-lg md:text-xl lg:text-2xl">
        响应式标题
      </h2>
    </div>
  );
};
```

#### CSS Modules 使用

```tsx
// ✅ 复杂样式使用 CSS Modules
import styles from './CustomComponent.module.scss';

export const CustomComponent: React.FC = () => {
  return (
    <div className={styles.container}>
      <div className={styles.header} />
      <div className={styles.content} />
    </div>
  );
};

// CustomComponent.module.scss
.container {
  display: flex;
  flex-direction: column;

  .header {
    padding: 16px;
    background: #f5f5f5;
  }

  .content {
    flex: 1;
  }
}
```

### 3.5 性能优化规范

#### React.memo 使用

```tsx
// ✅ 好的写法：对纯展示组件使用 memo
export const BotCard: React.FC<BotCardProps> = React.memo(({ bot }) => {
  return (
    <div className="bot-card">
      <h3>{bot.name}</h3>
      <p>{bot.description}</p>
    </div>
  );
}, (prevProps, nextProps) => {
  // 自定义比较函数
  return prevProps.bot.id === nextProps.bot.id
    && prevProps.bot.status === nextProps.bot.status;
});
```

#### useMemo 使用

```tsx
// ✅ 好的写法：缓存计算结果
export const ExpensiveComponent: React.FC = ({ data }) => {
  const sortedData = useMemo(() => {
    return data.sort((a, b) => a.createdAt - b.createdAt);
  }, [data]);

  const filteredData = useMemo(() => {
    return sortedData.filter(item => item.status === 'active');
  }, [sortedData]);

  return <div>{filteredData.map(item => <Item key={item.id} data={item} />)}</div>;
};
```

#### useCallback 使用

```tsx
// ✅ 好的写法：缓存回调函数
export const BotForm: React.FC = () => {
  const [formData, setFormData] = useState<BotFormData>({});

  const handleSubmit = useCallback(async (event: React.FormEvent) => {
    event.preventDefault();
    await botApi.create(formData);
    Toast.success('Bot 创建成功');
  }, [formData]);

  const handleChange = useCallback((field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  }, []);

  return (
    <form onSubmit={handleSubmit}>
      <Input value={formData.name} onChange={(v) => handleChange('name', v)} />
    </form>
  );
};
```

### 3.6 测试规范

#### Vitest 单元测试

```typescript
// __tests__/Button.test.tsx
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from './index';

describe('Button', () => {
  it('should render button with text', () => {
    render(<Button type="primary">点击我</Button>);
    expect(screen.getByText('点击我')).toBeInTheDocument();
  });

  it('should call onClick when clicked', () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>点击我</Button>);

    fireEvent.click(screen.getByText('点击我'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('should be disabled when disabled prop is true', () => {
    render(<Button disabled>点击我</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });
});
```

---

## 🗄️ 四、数据库开发规范

### 4.1 表设计规范

#### 表命名

```sql
-- ✅ 好的表命名：小写 + 下划线 + 复数
CREATE TABLE users (
    user_id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ❌ 不好的表命名
CREATE TABLE UserInfo (  -- 大驼峰
    ID BIGINT,               -- 不一致的命名
    UserName VARCHAR(100)    -- 驼峰命名
) ENGINE=InnoDB;
```

#### 字段命名

```sql
-- ✅ 好的字段命名
CREATE TABLE bots (
    -- 主键：{entity}_id
    bot_id VARCHAR(36) PRIMARY KEY,

    -- 外键：{referenced_entity}_id
    tenant_id VARCHAR(36) NOT NULL,
    creator_id VARCHAR(36) NOT NULL,

    -- 布尔值：is_{property}
    is_public BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,

    -- 时间戳：{action}_at
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_creator_id (creator_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 索引设计

```sql
-- ✅ 好的索引设计
-- 1. 外键字段建索引
INDEX idx_tenant_id (tenant_id)

-- 2. 查询条件建索引
INDEX idx_status_created (status, created_at)

-- 3. 唯一索引
UNIQUE KEY uk_user_email (email)

-- 4. 复合索引遵循最左前缀原则
INDEX idx_tenant_status_created (tenant_id, status, created_at)
-- 可以使用:
-- WHERE tenant_id = ?
-- WHERE tenant_id = ? AND status = ?
-- WHERE tenant_id = ? AND status = ? AND created_at = ?
```

### 4.2 查询优化规范

#### 避免 SELECT *

```sql
-- ❌ 不好的写法：SELECT *
SELECT * FROM bots WHERE tenant_id = 'xxx';

-- ✅ 好的写法：只查询需要的字段
SELECT bot_id, name, status FROM bots WHERE tenant_id = 'xxx';
```

#### 分页查询

```sql
-- ✅ 好的分页查询
-- 方法 1: 使用 OFFSET
SELECT bot_id, name, status
FROM bots
WHERE tenant_id = 'xxx'
ORDER BY created_at DESC
LIMIT 20 OFFSET 0;  -- 第 1 页

-- 方法 2: 使用游标分页（性能更好）
-- 第 1 次查询
SELECT bot_id, name, status
FROM bots
WHERE tenant_id = 'xxx'
ORDER BY created_at DESC
LIMIT 21;  -- 多查 1 条

-- 第 2 次查询（使用上次最后一条记录的 ID）
SELECT bot_id, name, status
FROM bots
WHERE tenant_id = 'xxx' AND created_at < :last_created_at
ORDER BY created_at DESC
LIMIT 21;
```

#### 避免 N+1 查询

```sql
-- ❌ 不好的写法：N+1 查询
-- 先查询 Bot 列表
SELECT * FROM bots WHERE tenant_id = 'xxx';
-- 然后循环查询每个 Bot 的创建者信息
SELECT * FROM users WHERE user_id IN (...);

-- ✅ 好的写法：使用 JOIN
SELECT
    b.bot_id,
    b.name AS bot_name,
    u.username AS creator_name
FROM bots b
LEFT JOIN users u ON b.creator_id = u.user_id
WHERE b.tenant_id = 'xxx';
```

### 4.3 事务使用规范

```go
// ✅ 好的实践：事务范围尽可能小
func (s *BotService) CreateBotWithConfig(ctx context.Context, bot *Bot, config *BotConfig) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 创建 Bot
        if err := tx.Create(bot).Error; err != nil {
            return fmt.Errorf("failed to create bot: %w", err)
        }

        // 2. 创建配置
        config.BotID = bot.BotID
        if err := tx.Create(config).Error; err != nil {
            return fmt.Errorf("failed to create config: %w", err)
        }

        // 3. 更新配额使用
        if err := s.updateQuotaUsage(ctx, tx, bot.TenantID, "bots"); err != nil {
            return fmt.Errorf("failed to update quota: %w", err)
        }

        return nil
    })
}
```

---

## 🔧 五、API 开发规范

### 5.1 RESTful API 设计

#### URL 设计

```yaml
# ✅ 好的 URL 设计
GET    /api/v1/tenants                    # 租户列表
GET    /api/v1/tenants/:tenant_id         # 租户详情
POST   /api/v1/tenants                    # 创建租户
PUT    /api/v1/tenants/:tenant_id         # 更新租户
DELETE /api/v1/tenants/:tenant_id         # 删除租户

# 嵌套资源
GET    /api/v1/tenants/:tenant_id/bots    # 租户的 Bot 列表
POST   /api/v1/tenants/:tenant_id/bots    # 为租户创建 Bot

# ❌ 不好的 URL 设计
GET    /api/v1/getTenants                  # 动词在 URL 中
GET    /api/v1/tenant/:id                  # 资源名单数
POST   /api/v1/tenant/create              # 动词在 URL 中
```

#### HTTP 方法使用

```yaml
# 安全性和幂等性
GET    /api/v1/bots/:bot_id     # 安全、幂等
POST   /api/v1/bots             # 不安全、不幂等
PUT    /api/v1/bots/:bot_id     # 不安全、幂等
PATCH  /api/v1/bots/:bot_id     # 不安全、不幂等
DELETE /api/v1/bots/:bot_id     # 不安全、幂等
```

#### 状态码使用

```yaml
# 成功响应
200 OK              # GET、PATCH 成功
201 Created         # POST 创建成功
204 No Content      # DELETE 成功

# 客户端错误
400 Bad Request    # 请求参数错误
401 Unauthorized    # 未认证
403 Forbidden       # 无权限
404 Not Found       # 资源不存在
409 Conflict        # 资源冲突
422 Unprocessable  # 请求格式正确但语义错误
429 Too Many Requests # 请求过于频繁

# 服务端错误
500 Internal Server Error  # 服务器内部错误
502 Bad Gateway           # 网关错误
503 Service Unavailable    # 服务不可用
```

### 5.2 响应格式规范

#### 成功响应

```json
{
  "code": "SUCCESS",
  "message": "Operation successful",
  "message_zh": "操作成功",
  "message_en": "Operation successful",
  "message_ja": "操作成功",
  "data": {
    "bot_id": "bot_123",
    "name": "Test Bot",
    "status": "published"
  },
  "request_id": "req_1234567890",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

#### 错误响应

```json
{
  "code": "BOT404",
  "message": "Bot not found",
  "message_zh": "Bot 不存在",
  "message_en": "Bot not found",
  "message_ja": "Botが存在しません",
  "request_id": "req_1234567890",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

---

## 📝 六、Git 工作流规范

### 6.1 分支策略

```yaml
main (生产)
  ↑
develop (开发主分支)
  ↑
feature/tenant-system      # 功能分支（新功能）
bugfix/permission-error   # 缺陷修复分支
hotfix/security-patch     # 紧急修复分支
release/v1.0.0             # 发布分支
```

### 6.2 提交信息规范

#### Conventional Commits

```bash
# 格式：<type>(<scope>): <subject>

# type 类型：
feat:     新功能
fix:      Bug 修复
docs:     文档更新
style:    代码格式（不影响功能）
refactor:  重构（不是新功能也不是修复）
perf:     性能优化
test:     测试相关
chore:    构建过程或辅助工具的变动
revert:   回退提交

# 示例：
feat(auth): add tenant_id field to user table
fix(permission): resolve role check bug
docs(api): update authentication API documentation
test(bot): add unit tests for bot service
refactor(rbac): simplify permission checking logic
```

### 6.3 Pull Request 规范

#### PR 标题

```yaml
# ✅ 好的 PR 标题
feat(rbac): implement 5-level data permission control

# ❌ 不好的 PR 标题
implement rbac  # 缺少类型和范围
fix bugs        # 过于模糊
update code     # 信息不足
```

#### PR 描述模板

```markdown
## 变更类型
- [ ] feat: 新功能
- [ ] fix: Bug 修复
- [ ] docs: 文档更新
- [ ] style: 代码格式
- [ ] refactor: 重构
- [ ] perf: 性能优化
- [ ] test: 测试
- [ ] chore: 构建/工具

## 变更说明
<!-- 详细描述本次变更的内容 -->

## 相关 Issue
Closes #123
Related to #456

## 测试
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 手动测试通过

## 检查清单
- [ ] 代码风格符合规范
- [ ] 添加了必要的注释
- [ ] 更新了相关文档
- [ ] 没有引入新的警告
```

---

## ✅ 七、代码审查清单

### 7.1 后端代码审查

#### 功能性
- [ ] 代码实现了需求的所有功能
- [ ] 边界情况处理完善
- [ ] 错误处理覆盖全面

#### 可读性
- [ ] 代码结构清晰，易于理解
- [ ] 变量命名准确描述用途
- [ ] 复杂逻辑有注释说明

#### 性能
- [ ] 数据库查询使用索引
- [ ] 避免 N+1 查询
- [ ] 合理使用缓存
- [ ] 并发安全

#### 安全性
- [ ] 输入验证完善
- [ ] SQL 注入防护
- [ ] 权限检查完整
- [ ] 敏感数据脱敏

#### 测试
- [ ] 单元测试覆盖率 > 80%
- [ ] 关键逻辑有集成测试
- [ ] 测试用例覆盖正常和异常情况

### 7.2 前端代码审查

#### 组件设计
- [ ] 组件职责单一
- [ ] Props 接口合理
- [ ] 状态管理正确

#### 性能
- [ ] 避免不必要的重渲染
- [ ] 使用 React.memo 优化
- [ ] 大列表使用虚拟化
- [ ] 图片和资源优化

#### 可访问性
- [ ] 支持键盘导航
- [ ] 语义化 HTML
- [ ] ARIA 标签完整
- [ ] 颜色对比度符合标准

#### 测试
- [ ] 单元测试覆盖组件
- [ ] 关键交互有 E2E 测试
- [ ] 测试覆盖多种场景

---

## 🔍 八、全局一致性检查清单

### 8.1 代码一致性

- [ ] **Go 代码**: 通过 golangci-lint 检查
- [ ] **TypeScript 代码**: 通过 ESLint 检查
- [ ] **命名规范**: 遵循命名约定
- [ ] **格式化**: 统一使用 go fmt / prettier

### 8.2 架构一致性

- [ ] **DDD 分层**: 遵循 Controller → Service → Repository → DAO
- [ ] **依赖方向**: 上层依赖下层
- [ ] **模块边界**: 跨模块通信通过定义接口
- [ ] **事件驱动**: 跨模块异步通信使用 EventBus

### 8.3 数据一致性

- [ ] **租户隔离**: 所有业务表包含 `tenant_id`
- [ ] **审计字段**: `created_at`, `updated_at`, `created_by`, `updated_by`
- [ ] **软删除**: 使用 `deleted_at`
- [ ] **事务管理**: 跨表操作使用事务

### 8.4 API 一致性

- [ ] **URL 规范**: `/api/v{version}/{resource}`
- [ ] **响应格式**: 统一的 `{code, message, data, request_id, timestamp}`
- [ ] **错误码**: 使用统一的错误码定义
- [ ] **认证**: Bearer Token 方式

### 8.5 文档一致性

- [ ] **API 文档**: 所有 API 有 OpenAPI 规范
- [ ] **组件文档**: 所有组件有使用示例
- [ ] **变更日志**: 每次发布更新 CHANGELOG
- [ ] **架构文档**: 关键决策有 ADR

---

## 📞 九、规范执行与监督

### 9.1 自动化检查

**Pre-commit Hooks**:
```yaml
# .husky/pre-commit
#!/bin/sh
# 运行 linter
rush lint

# 运行类型检查
rush typecheck

# 运行格式化
rush prettier

# 运行单元测试
rush test --changed
```

**CI Pipeline**:
```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]
jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Lint
        run: rush lint

      - name: Type Check
        run: rush typecheck

      - name: Unit Tests
        run: rush test --coverage

      - name: Security Scan
        run: trivy scan . --security-checks vuln,config
```

### 9.2 代码审查流程

**审查者职责**:
1. **功能审查**: 验证代码实现需求
2. **质量审查**: 检查代码质量和规范
3. **安全审查**: 识别潜在安全风险
4. **性能审查**: 评估性能影响

**审查响应时间**:
- 简单 PR: < 4 小时
- 复杂 PR: < 1 天
- 紧急修复: < 1 小时

### 9.3 定期评审

**每周代码质量会议**:
- 时间: 每周五下午 3:00
- 参与者: 所有开发人员
- 内容:
  - 审查本周代码质量报告
  - 讨论规范执行情况
  - 识别改进机会

**每季度规范评审**:
- 时间: 每季度最后一周
- 内容:
  - 评估规范执行效果
  - 收集开发人员反馈
  - 更新规范内容

---

## 📚 十、参考文档

### 10.1 已生成的 ZKER 文档

1. [ZKER-统一错误码定义规范.md](./ZKER-统一错误码定义规范.md)
2. [ZKER-技术组件清单与使用指南(完整版).md](./ZKER-技术组件清单与使用指南(完整版).md)
3. [openapi/zker-api-v1-core-modules.yaml](./openapi/zker-api-v1-core-modules.yaml)
4. [ZKER-核心算法实现指南.md](./ZKER-核心算法实现指南.md)
5. [ZKER-开发快速入门指南.md](./ZKER-开发快速入门指南.md)
6. [ZKER-性能测试计划_v1.0.md](./ZKER-性能测试计划_v1.0.md)
7. [ZKER-数据迁移方案_v1.0.md](./ZKER-数据迁移方案_v1.0.md)
8. [ZKER-灰度发布策略_v1.0.md](./ZKER-灰度发布策略_v1.0.md)
9. [ZKER-一键回滚方案_v1.0.md](./ZKER-一键回滚方案_v1.0.md)
10. [ZKER-故障排查手册_v1.0.md](./ZKER-故障排查手册_v1.0.md)
11. [ZKER-实现差距分析与研发计划_v1.0.md](./ZKER-实现差距分析与研发计划_v1.0.md)

### 10.2 外部参考资料

**Go 代码规范**:
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide)

**React 规范**:
- [React Docs](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [Airbnb React/JSX Style Guide](https://github.com/airbnb/javascript/tree/master/react)

**数据库规范**:
- [MySQL 8.0 Reference Manual](https://dev.mysql.com/doc/refman/8.0/en/)
- [SQL Style Guide](https://www.sqlstyle.guide/)

---

**文档维护**:
- 版本: v1.0
- 最后更新: 2025-01-01
- 下次审查: 每季度
- 维护团队: 技术架构委员会

---

**© 2025 ZKER Project. All rights reserved.**
