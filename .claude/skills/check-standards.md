# 代码规范检查 Skill

## 技能描述

全面检查代码是否符合 Coze Studio 开发规范，包括命名规范、代码结构、错误处理等各个方面。

## 适用场景

- 代码提交前的自检
- Pull Request 代码审查
- 代码重构验证
- 新人代码质量检查

## 检查流程

### 1. 前端代码检查

#### 1.1 文件命名检查

```typescript
// ✅ 正确
user-profile.tsx
use-auth-form.ts
user.types.ts
api.constants.ts

// ❌ 错误（需要报告）
UserProfile.tsx           → 应改为 user-profile.tsx
useAuthForm.ts            → 应改为 use-auth-form.ts
userTypes.ts              → 应改为 user.types.ts
```

**检查要点**：
- [ ] 组件文件：kebab-case.tsx
- [ ] Hooks 文件：use-{purpose}.ts
- [ ] 类型文件：{entity}.types.ts
- [ ] 常量文件：{module}.constants.ts
- [ ] 测试文件：{filename}.test.ts

#### 1.2 变量命名检查

```typescript
// ✅ 正确
const MAX_RETRY_COUNT = 3;
const isActive = true;
const userId = 123;

// ❌ 错误（需要报告）
const max_retry_count = 3;        → 应使用 MAX_RETRY_COUNT
const active = true;               → 应使用 isActive
const UserId = 123;                → 应使用 userId
```

**检查要点**：
- [ ] 常量：UPPER_SNAKE_CASE
- [ ] 变量：camelCase
- [ ] 类/接口：PascalCase
- [ ] 布尔值：is/has/can 前缀
- [ ] 私有变量：_ 前缀

#### 1.3 函数命名检查

```typescript
// ✅ 正确
function getUserById(id: number): User {}
async function fetchUserData(id: number): Promise<User> {}
function handleSubmit(event: FormEvent): void {}
function isValidUser(user: User): boolean {}

// ❌ 错误（需要报告）
function user(id: number) {}                    → 应改为 getUserById
function data(id: number) {}                    → 应改为 fetchData
function submit() {}                             → 应改为 handleSubmit
function valid(user: User) {}                   → 应改为 isValidUser
```

**检查要点**：
- [ ] 获取数据：get/fetch 前缀
- [ ] 设置数据：set/update 前缀
- [ ] 创建/删除：create/delete 前缀
- [ ] 事件处理：handle 前缀
- [ ] 回调函数：on 前缀
- [ ] 布尔判断：is/has/can 前缀

#### 1.4 组件结构检查

```typescript
// ✅ 正确的结构
export function UserProfile({ userId, onEdit }: UserProfileProps) {
  // 1. Hooks
  const [state, setState] = useState();

  // 2. 派生状态
  const value = useMemo(() => compute(state), [state]);

  // 3. 副作用
  useEffect(() => {}, []);

  // 4. 事件处理
  const handler = useCallback(() => {}, []);

  // 5. 渲染辅助
  function helper() {}

  // 6. 条件渲染
  if (loading) return <Spinner />;

  // 7. 返回 JSX
  return <div>...</div>;
}
```

**检查要点**：
- [ ] Hooks 按顺序调用
- [ ] useMemo 用于计算
- [ ] useEffect 用于副作用
- [ ] useCallback 用于事件处理
- [ ] 早期返回用于条件渲染
- [ ] JSX 在最后返回

#### 1.5 TypeScript 类型检查

```typescript
// ✅ 正确
interface UserProfileProps {
  userId: number;
  onEdit?: () => void;
}

function getUserById(id: number): User | null {}

// ❌ 错误（需要报告）
function getUserById(id) {}                 → 缺少类型注解
interface Props {}                          → 过于通用
```

**检查要点**：
- [ ] 所有 Props 必须定义接口
- [ ] 函数参数必须有类型
- [ ] 函数返回值必须有类型
- [ ] 避免使用 any 类型
- [ ] 泛型使用描述性名称

### 2. 后端代码检查

#### 2.1 文件命名检查

```go
// ✅ 正确
user.go
user_service.go
user_repository.go

// ❌ 错误（需要报告）
User.go                      → 应改为 user.go
user_service.go              ✅ 正确
UserService.go               → 应改为 user_service.go
```

**检查要点**：
- [ ] Go 文件：lowercase.go
- [ ] 测试文件：{filename}_test.go
- [ ] 包目录：小写单词

#### 2.2 包命名检查

```go
// ✅ 正确
package service
package repository
package entity

// ❌ 错误（需要报告）
package UserService          → 应改为 service
package user_repository      → 应改为 repository
package UserManagement       → 应改为 user
```

**检查要点**：
- [ ] 包名：小写单词
- [ ] 包名：简洁明了
- [ ] 避免使用 pkg、common

#### 2.3 变量命名检查

```go
// ✅ 正确
const MaxRetryCount = 3
var userID int64
var isActive bool

// ❌ 错误（需要报告）
const max_retry_count = 3       → 应改为 MaxRetryCount
var UserId int64                 → 应改为 userID
var Active bool                  → 应改为 isActive
```

**检查要点**：
- [ ] 导出常量：PascalCase
- [ ] 内部常量：camelCase 或 PascalCase
- [ ] 导出变量：PascalCase
- [ ] 内部变量：camelCase
- [ ] 缩写词保持一致（ID, URL, HTTP）

#### 2.4 函数/方法命名检查

```go
// ✅ 正确
func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {}
func (s *Service) GetUserByID(ctx context.Context, id int64) (*User, error) {}
func (s *Service) ValidateUserEmail(email string) error {}
func (s *Service) IsUserActive(userID int64) (bool, error) {}

// ❌ 错误（需要报告）
func (s *Service) user(id int64) {}                      → 应改为 GetUserByID
func (s *Service) data(id int64) {}                      → 应改为 GetDataByID
func (s *Service) validate() {}                          → 应改为 ValidateUser
```

**检查要点**：
- [ ] 导出函数：PascalCase
- [ ] 内部函数：camelCase
- [ ] 方法接收者：单字母简写
- [ ] 构造函数：New 前缀
- [ ] Get/Create/Update/Delete 前缀

#### 2.5 接口定义检查

```go
// ✅ 正确
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    GetUserByID(ctx context.Context, id int64) (*User, error)
}

type userServiceImpl struct {
    repo UserRepository
}

// ❌ 错误（需要报告）
type IUserService interface {}                           → 不使用 I 前缀
type UserServiceImpl struct {}                          → 内部实现应小写
```

**检查要点**：
- [ ] 接口命名：PascalCase，不使用 I 前缀
- [ ] 实现命名：{interface}Impl，小写开头
- [ ] 接口方法：清晰描述功能
- [ ] 接口隔离：职责单一

#### 2.6 错误处理检查

```go
// ✅ 正确
func (s *service) CreateUser(ctx context.Context, req *Request) (*User, error) {
    if err := s.validate(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, errorx.Wrapf(err, "create failed, id=%d", req.ID)
    }

    return user, nil
}

// ❌ 错误（需要报告）
user, _ := s.repo.Create(ctx, req)                      → 不应忽略错误
return nil, err                                         → 应包装错误
```

**检查要点**：
- [ ] 所有错误必须检查
- [ ] 错误必须处理或传递
- [ ] 使用 fmt.Errorf 包装
- [ ] 使用 errorx.Wrapf 添加上下文
- [ ] 错误信息清晰明确

#### 2.7 DDD 架构检查

```
✅ 正确的分层
domain/{domain}/
├── entity/           # 实体
├── repository/       # 仓储接口
└── service/          # 领域服务

application/{domain}/
└── usecase/          # 用例

api/handler/          # API 处理器
infra/                # 基础设施
```

**检查要点**：
- [ ] 领域层不依赖基础设施
- [ ] 应用层编排领域逻辑
- [ ] API 层仅处理 HTTP
- [ ] 依赖方向正确

### 3. 通用检查

#### 3.1 导入顺序检查

```typescript
// ✅ 正确
import path from 'path';
import { Button } from '@douyinfe/semi-ui';
import { useAuth } from '@coze-studio/auth';
import { formatDate } from './utils/format-date';

// ❌ 错误（需要报告）
import { Button } from '@douyinfe/semi-ui';
import path from 'path';                    → 顺序错误
```

**检查要点**：
- [ ] Node.js 内置模块在前
- [ ] 第三方库按字母排序
- [ ] 内部包按字母排序
- [ ] 相对路径在最后

#### 3.2 注释检查

```go
// ✅ 正确
// CreateUser 创建新用户
//
// 该方法会验证用户数据的合法性，检查邮箱是否已存在，
// 然后将用户信息持久化到数据库。
//
// 参数:
//   ctx - 请求上下文
//   req - 创建用户请求
//
// 返回:
//   *User - 创建的用户实体
//   error - 错误信息（如果创建失败）
func (s *service) CreateUser(ctx context.Context, req *Request) (*User, error) {}
```

**检查要点**：
- [ ] 导出函数必须有注释
- [ ] 复杂逻辑必须有说明
- [ ] 参数说明完整
- [ ] 返回值说明完整
- [ ] 错误情况说明清楚

#### 3.3 测试检查

```typescript
// ✅ 正确
describe('ComponentName', () => {
  describe('when loading', () => {
    it('should show spinner', () => {});
  });

  describe('when loaded', () => {
    it('should display data', () => {});
  });
});
```

**检查要点**：
- [ ] 测试文件命名正确
- [ ] 使用 describe/it 结构
- [ ] 测试描述清晰
- [ ] 覆盖率符合要求
  - Level 1: 80%
  - Level 2: 30%
  - Level 3-4: 灵活

## 输出格式

### 检查报告模板

```markdown
# 代码规范检查报告

## 检查概述
- **检查时间**: 2025-01-01 12:00:00
- **检查文件**: {文件路径}
- **总体评分**: ⭐⭐⭐⭐☆ (4/5)

## 检查结果

### ✅ 通过项 (12)
- ✅ 文件命名符合规范
- ✅ 变量命名符合规范
- ✅ 函数命名符合规范
- ...

### ❌ 失败项 (3)

#### 1. 函数命名不规范
- **位置**: {file}:{line}
- **问题**: 函数名 `user` 不符合规范
- **建议**: 改为 `getUserById`
- **代码**:
  ```typescript
  function user(id: number) {}
  ```
- **修复**:
  ```typescript
  function getUserById(id: number) {}
  ```

#### 2. 缺少错误处理
- **位置**: {file}:{line}
- **问题**: 未检查函数返回的错误
- **建议**: 添加错误检查
- **代码**:
  ```go
  user, _ := repo.Create(ctx, req)
  ```
- **修复**:
  ```go
  user, err := repo.Create(ctx, req)
  if err != nil {
      return nil, fmt.Errorf("create failed: %w", err)
  }
  ```

#### 3. 缺少类型注解
- **位置**: {file}:{line}
- **问题**: 函数参数缺少类型
- **建议**: 添加类型注解
- **代码**:
  ```typescript
  function getUser(id) {}
  ```
- **修复**:
  ```typescript
  function getUser(id: number): User {}
  ```

## 改进建议

1. **命名规范**: 统一使用 camelCase 和 PascalCase
2. **错误处理**: 所有函数调用必须检查错误
3. **类型注解**: 所有函数参数和返回值必须有类型

## 总体评价

代码整体质量较好，但需要改进以下方面：
- 函数命名需要更明确
- 错误处理需要完善
- 类型注解需要补充
```

## 检查清单

### 前端检查清单
- [ ] 文件命名：kebab-case
- [ ] 组件名：PascalCase
- [ ] 变量：camelCase
- [ ] 常量：UPPER_SNAKE_CASE
- [ ] 函数：动词-名词模式
- [ ] 组件结构：7 步骤
- [ ] 类型定义：完整
- [ ] 测试覆盖：符合要求

### 后端检查清单
- [ ] 文件命名：lowercase.go
- [ ] 包名：小写单词
- [ ] 导出：PascalCase
- [ ] 内部：camelCase
- [ ] 接口：无 I 前缀
- [ ] 实现：Impl 后缀
- [ ] 错误：全部检查
- [ ] DDD：层次清晰

## 注意事项

1. **严格性**: 严格按照规范执行
2. **建议性**: 对于模糊的规范，给出建议
3. **优先级**: 严重错误优先报告
4. **可操作性**: 提供具体的修复方案
5. **全局性**: 考虑整体代码一致性
