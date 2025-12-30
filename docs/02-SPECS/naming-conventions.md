# ZKER 命名规范

**版本**: v3.0.0 | **更新**: 2025-01-03 | **状态**: 强制执行

---

## 目录

- [通用原则](#通用原则)
- [Go 命名规范](#go-命名规范)
- [TypeScript 命名规范](#typescript-命名规范)
- [数据库命名规范](#数据库命名规范)
- [API 命名规范](#api-命名规范)
- [检查清单](#检查清单)

---

## 通用原则

### 核心原则

| 原则 | 说明 | 示例 |
|------|------|------|
| **清晰优先** | 名字要能清晰表达意图 | `GetUserByID` ✅ vs `Get` ❌ |
| **一致性** | 同类事物使用相同命名模式 | `CreateXxx`, `UpdateXxx`, `DeleteXxx` |
| **避免缩写** | 除非是业界通用缩写 | `UserID` ✅ vs `Uid` ❌ |
| **避免拼写错误** | 正确使用英语 | `Message` ✅ vs `Mesage` ❌ |
| **避免数字后缀** | 使用更有意义的名字 | `MessageV2` ❌ vs `MessageWithAttachment` ✅ |

### 通用缩写表

| 缩写 | 全称 | 使用场景 |
|------|------|---------|
| **ID** | Identifier | `user_id`, `userID` |
| **API** | Application Programming Interface | `api_handler.go`, `API_BASE_URL` |
| **URL** | Uniform Resource Locator | `avatar_url`, `url` |
| **HTTP** | HyperText Transfer Protocol | `http_status`, `httpClient` |
| **JSON** | JavaScript Object Notation | `json_data`, `json.Marshal` |
| **SQL** | Structured Query Language | `sql_query`, `sql.DB` |
| **DB** | Database | `db_client`, `mysqlDB` |
| **DTO** | Data Transfer Object | `UserDTO`, `LoginRequestDTO` |
| **RBAC** | Role-Based Access Control | `rbac_service.go`, `RBACMiddleware` |
| **JWT** | JSON Web Token | `jwt_token`, `jwtClaims` |
| **FAQ** | Frequently Asked Questions | `faq_page`, `FAQComponent` |

---

## Go 命名规范

### 包命名 (Package)

#### 规则

1. **全小写，单个单词**
2. **简短且有意义**
3. **不要使用下划线或驼峰**
4. **不要使用复数形式**

#### 示例

```go
// ✅ Good
package tenant
package permission
package routing
package conversation

// ❌ Bad
package tenants        // 不要使用复数
package tenantService  // 不要使用驼峰
package tenant_service // 不要使用下划线
package Tls            // 不要使用缩写（除非业界通用）
```

### 文件命名

#### 规则

1. **全小写，使用下划线分隔**
2. **以实现的核心内容命名**
3. **测试文件以 `_test.go` 结尾**

#### 示例

```
// ✅ Good
domain/tenant/
├── entity/
│   └── tenant.go              # 租户实体
├── repository/
│   ├── tenant_repository.go   # 租户仓储接口
│   └── tenant_repository_impl.go # 租户仓储实现
├── service/
│   ├── tenant_service.go      # 租户服务接口
│   └── tenant_service_impl.go # 租户服务实现
└── tenant_test.go             # 租户测试

// ❌ Bad
domain/tenant/
├── Tenant.go                  # 不要大写开头
├── tenantRepo.go              # 不要使用驼峰
├── TenantService.go           # 不要大写开头
└── tenant_service_test_test.go # 不要重复 _test
```

### 变量命名

#### 局部变量

```go
// ✅ Good: 驼峰命名，清晰表达意图
func (s *Service) CreateTenant(ctx context.Context, req *CreateRequest) error {
    tenantID := req.TenantID
    tenantName := req.TenantName
    isActive := req.Status == "active"

    // 使用有意义的名字
    tenant, err := s.repo.FindByID(ctx, tenantID)
    if err != nil {
        return err
    }

    // ...
}

// ❌ Bad: 过于简短，没有意义
func (s *Service) CreateTenant(ctx context.Context, req *CreateRequest) error {
    tid := req.TenantID
    tn := req.TenantName
    f := req.Status == "active"

    t, err := s.repo.FindByID(ctx, tid)
    if err != nil {
        return err
    }

    // ...
}
```

#### 常量

```go
// ✅ Good: 全大写，使用下划线分隔
const (
    MAX_TENANTS = 10000
    DEFAULT_PAGE_SIZE = 20
    TOKEN_EXPIRE_DURATION = 24 * time.Hour
)

// ❌ Bad
const (
    maxTenants = 10000        // 不要使用驼峰
    defaultPageSize = 20      // 不要使用驼峰
    tokenExpireDuration = 24 * time.Hour // 不要使用驼峰
)
```

### 函数命名

#### 规则

1. **驼峰命名（PascalCase for exported, camelCase for unexported）**
2. **以动词开头**
3. **清晰表达功能**
4. **参数 ≤ 5 个**

#### 示例

```go
// ✅ Good
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateRequest) (*Tenant, error)
func (s *TenantService) GetTenantByID(ctx context.Context, id string) (*Tenant, error)
func (s *TenantService) UpdateTenant(ctx context.Context, id string, req *UpdateRequest) (*Tenant, error)
func (s *TenantService) DeleteTenant(ctx context.Context, id string) error
func (s *TenantService) ListTenants(ctx context.Context, req *ListRequest) ([]*Tenant, int64, error)

// 内部函数使用 camelCase
func (s *TenantService) validateTenant(tenant *Tenant) error
func (s *TenantService) sanitizeInput(input string) string

// ❌ Bad
func (s *TenantService) create(ctx context.Context, req *CreateRequest) (*Tenant, error) // 导出函数要大写开头
func (s *TenantService) Tenant(id string) (*Tenant, error) // 要有动词
func (s *TenantService) GetTenantByIDWithPermissionsAndQuota(ctx context.Context, id string, includeDeleted bool, fetchRoles bool, loadQuota bool) (*Tenant, error) // 参数太多
```

### 接口命名

#### 规则

1. **接口名以 `-er` 结尾**
2. **单方法接口通常以方法名命名**

#### 示例

```go
// ✅ Good
type TenantRepository interface {
    Create(ctx context.Context, tenant *Tenant) error
    FindByID(ctx context.Context, id string) (*Tenant, error)
    Update(ctx context.Context, tenant *Tenant) error
    Delete(ctx context.Context, id string) error
}

type TenantService interface {
    CreateTenant(ctx context.Context, req *CreateRequest) (*Tenant, error)
    GetTenantByID(ctx context.Context, id string) (*Tenant, error)
}

// 单方法接口
type Stringer interface {
    String() string
}

type Reader interface {
    Read(p []byte) (n int, err error)
}

// ❌ Bad
type ITenantRepository interface {} // 不要使用 I 前缀
type TenantRepositoryInterface interface {} // 不要使用 Interface 后缀
type TenantRepo interface {} // 不要使用缩写（除非业界通用）
```

### 结构体命名

#### 规则

1. **PascalCase（导出）或 camelCase（不导出）**
2. **清晰表达类型**

#### 示例

```go
// ✅ Good
type Tenant struct {
    TenantID   string    `json:"tenant_id"`
    TenantName string    `json:"tenant_name"`
    Status     string    `json:"status"`
    CreatedAt  time.Time `json:"created_at"`
}

type CreateTenantRequest struct {
    TenantName string `json:"tenant_name" validate:"required,min=3,max=100"`
    PlanType   string `json:"plan_type" validate:"required,oneof=free pro enterprise"`
}

type TenantResponse struct {
    TenantID   string `json:"tenant_id"`
    TenantName string `json:"tenant_name"`
    PlanType   string `json:"plan_type"`
    Quota      Quota  `json:"quota"`
}

// ❌ Bad
type tenant struct {} // 导出类型要大写开头
type tenantStruct struct {} // 不要使用 Struct 后缀
type TenantData struct {} // 不要使用 Data 后缀（除非必要）
```

### 错误命名

#### 规则

1. **使用统一错误码**
2. **错误变量以 `Err` 开头**

#### 示例

```go
// ✅ Good
import "backend/types/errno"

// 使用统一错误码
return errorx.Wrapf(err, errno.TenantNotFound)

// 自定义错误
var (
    ErrTenantNotFound = errors.New("tenant not found")
    ErrTenantAlreadyExists = errors.New("tenant already exists")
    ErrInvalidTenantName = errors.New("invalid tenant name")
)

// ❌ Bad
return errors.New("tenant not found") // 不要硬编码错误信息
return fmt.Errorf("tenant %s not found", id) // 尽量使用统一错误码
```

---

## TypeScript 命名规范

### 文件命名

#### 规则

1. **组件文件：PascalCase**
2. **工具文件：camelCase**
3. **类型文件：camelCase**
4. **常量文件：camelCase**
5. **Hooks 文件：use 前缀**

#### 示例

```
// ✅ Good
components/
├── TenantList.tsx          # 组件：PascalCase
├── TenantForm.tsx          # 组件：PascalCase
└── TenantCard.tsx          # 组件：PascalCase

hooks/
├── useTenants.ts           # Hook：use 前缀
├── useTenantForm.ts        # Hook：use 前缀
└── useModal.ts             # Hook：use 前缀

utils/
├── formatUtils.ts          # 工具：camelCase
├── validationUtils.ts      # 工具：camelCase
└── dateUtils.ts            # 工具：camelCase

types/
├── tenantTypes.ts          # 类型：camelCase + Types 后缀
├── apiTypes.ts             # 类型：camelCase + Types 后缀
└── commonTypes.ts          # 类型：camelCase + Types 后缀

constants/
├── apiConstants.ts         # 常量：camelCase + Constants 后缀
├── tenantConstants.ts      # 常量：camelCase + Constants 后缀
└── routeConstants.ts       # 常量：camelCase + Constants 后缀

// ❌ Bad
components/
├── tenantList.tsx          # 组件要大写开头
├── tenant-form.tsx         # 不要使用连字符
└── tenant_card.tsx         # 不要使用下划线

hooks/
├── tenantHook.ts           # Hook 要使用 use 前缀
├── getTenants.ts           # Hook 要使用 use 前缀
└── use_tenant.ts           # 不要使用下划线
```

### 变量命名

#### 规则

1. **camelCase**
2. **常量：UPPER_SNAKE_CASE**
3. **类型/接口：PascalCase**
4. **枚举：PascalCase**

#### 示例

```typescript
// ✅ Good
// 变量：camelCase
const tenantList: Tenant[] = [];
const currentPage: number = 1;
const isLoading: boolean = false;

// 常量：UPPER_SNAKE_CASE
const MAX_TENANTS = 10000;
const DEFAULT_PAGE_SIZE = 20;
const API_BASE_URL = 'http://localhost:8888';

// 类型/接口：PascalCase
interface Tenant {
  tenant_id: string;
  tenant_name: string;
  status: string;
  quota: Quota;
}

type TenantStatus = 'active' | 'inactive' | 'suspended';

interface CreateTenantRequest {
  tenant_name: string;
  plan_type: 'free' | 'pro' | 'enterprise';
}

// 枚举：PascalCase
enum TenantStatus {
  Active = 'active',
  Inactive = 'inactive',
  Suspended = 'suspended',
}

// ❌ Bad
const tenant_list: Tenant[] = [];      // 不要使用下划线
const TenantList: Tenant[] = [];       // 变量不要大写开头
const maxTenants = 10000;              // 常量要全大写
interface tenant {}                    // 接口要大写开头
type tenant_status = 'active';         // 类型不要使用下划线
```

### 函数命名

#### 规则

1. **camelCase**
2. **以动词开头**
3. **清晰表达功能**

#### 示例

```typescript
// ✅ Good
// API 调用
export const getTenants = async (params: GetTenantsParams) => { ... };
export const getTenantByID = async (id: string) => { ... };
export const createTenant = async (data: CreateTenantRequest) => { ... };
export const updateTenant = async (id: string, data: UpdateTenantRequest) => { ... };
export const deleteTenant = async (id: string) => { ... };

// 工具函数
export const formatDate = (date: Date) => { ... };
export const validateEmail = (email: string) => { ... };
export const sanitizeInput = (input: string) => { ... };

// 事件处理
export const handleTenantChange = (tenant: Tenant) => { ... };
export const handleFormSubmit = (values: FormValues) => { ... };
export const handleModalClose = () => { ... };

// ❌ Bad
export const Tenants = async () => { ... };               // 要有动词
export const tenant = async (id: string) => { ... };      // 要有动词
export const get_all_tenants = async () => { ... };       // 不要使用下划线
export const GetTenants = async () => { ... };            // 函数不要大写开头
```

### 组件命名

#### 规则

1. **PascalCase**
2. **清晰表达用途**

#### 示例

```typescript
// ✅ Good
export const TenantList: React.FC<TenantListProps> = () => { ... };
export const TenantForm: React.FC<TenantFormProps> = () => { ... };
export const TenantCard: React.FC<TenantCardProps> = () => { ... };
export const TenantModal: React.FC<TenantModalProps> = () => { ... };

// ❌ Bad
export const tenantList: React.FC = () => { ... };    // 组件要大写开头
export const Tenant_List: React.FC = () => { ... };   // 不要使用下划线
export const tenantlist: React.FC = () => { ... };    // 要使用 PascalCase
```

### Props 接口命名

#### 规则

1. **PascalCase + Props 后缀**

#### 示例

```typescript
// ✅ Good
interface TenantListProps {
  tenants: Tenant[];
  loading: boolean;
  onTenantClick: (tenant: Tenant) => void;
}

export const TenantList: React.FC<TenantListProps> = (props) => { ... };

// ❌ Bad
interface TenantList {}           // 要有 Props 后缀
interface tenantListProps {}      // 接口要大写开头
interface TenantListPropsInterface {} // 不要使用 Interface 后缀
```

---

## 数据库命名规范

### 表命名

#### 规则

1. **全小写**
2. **使用下划线分隔**
3. **使用复数形式**
4. **以业务含义命名**

#### 示例

```sql
-- ✅ Good
CREATE TABLE tenants (...);
CREATE TABLE users (...);
CREATE TABLE roles (...);
CREATE TABLE permissions (...);
CREATE TABLE user_roles (...);

-- ❌ Bad
CREATE TABLE Tenant (...);              -- 不要大写开头
CREATE TABLE tenant (...);               -- 要使用复数
CREATE TABLE TenantInfo (...);           -- 不要使用驼峰
CREATE TABLE tenant_info (...);          -- 不要使用 Info 后缀（除非必要）
CREATE TABLE tenant_data (...);          -- 不要使用 Data 后缀（除非必要）
```

### 字段命名

#### 规则

1. **全小写**
2. **使用下划线分隔**
3. **使用单数形式**
4. **以表名缩写作为前缀（可选）**

#### 示例

```sql
-- ✅ Good
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY,
    tenant_name VARCHAR(100) NOT NULL,
    plan_type VARCHAR(20) NOT NULL DEFAULT 'free',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    quota_max_tokens BIGINT NOT NULL DEFAULT 1000000,
    quota_used_tokens BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_tenant_name (tenant_name),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ❌ Bad
CREATE TABLE tenants (
    TenantID VARCHAR(36) PRIMARY KEY,              -- 不要使用驼峰
    tenant_name VARCHAR(100) NOT NULL,             -- 命名不一致
    PlanType VARCHAR(20) NOT NULL,                 -- 不要使用驼峰
    max_tokens BIGINT NOT NULL DEFAULT 1000000,    -- 要有前缀（多表时有歧义）
    CreateTime TIMESTAMP NOT NULL                  -- 不要使用驼峰
);
```

### 索引命名

#### 规则

1. **主键索引：`pk_表名`**
2. **唯一索引：`uk_表名_字段名`**
3. **普通索引：`idx_表名_字段名`**
4. **全文索引：`ft_表名_字段名`**

#### 示例

```sql
-- ✅ Good
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY,                           -- 主键：自动命名为 PRIMARY
    tenant_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    UNIQUE KEY uk_tenants_tenant_name (tenant_name),             -- 唯一索引
    UNIQUE KEY uk_tenants_email (email),                         -- 唯一索引
    INDEX idx_tenants_status (status),                           -- 普通索引
    INDEX idx_tenants_created_at (created_at)                    -- 普通索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ❌ Bad
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY,
    tenant_name VARCHAR(100) NOT NULL,
    UNIQUE KEY tenant_name (tenant_name),        -- 要有 uk_ 前缀
    INDEX status (status),                       -- 要有 idx_ 前缀
    INDEX created (created_at)                   -- 要有 idx_ 前缀
);
```

### 外键命名

#### 规则

1. **`fk_当前表_关联表_字段`**

#### 示例

```sql
-- ✅ Good
CREATE TABLE user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY fk_user_roles_users_user_id (user_id) REFERENCES users(user_id),
    FOREIGN KEY fk_user_roles_roles_role_id (role_id) REFERENCES roles(role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ❌ Bad
CREATE TABLE user_roles (
    ...
    FOREIGN KEY (user_id) REFERENCES users(user_id),           -- 要有 fk_ 前缀
    FOREIGN KEY (role_id) REFERENCES roles(role_id)            -- 要有 fk_ 前缀
);
```

---

## API 命名规范

### 路由命名

#### 规则

1. **全小写**
2. **使用连字符（kebab-case）**
3. **使用复数形式**
4. **RESTful 风格**

#### 示例

```go
// ✅ Good
GET    /api/v1/tenants              # 获取租户列表
POST   /api/v1/tenants              # 创建租户
GET    /api/v1/tenants/:id          # 获取租户详情
PUT    /api/v1/tenants/:id          # 更新租户
DELETE /api/v1/tenants/:id          # 删除租户
GET    /api/v1/tenants/:id/users    # 获取租户的用户列表
POST   /api/v1/tenants/:id/users    # 为租户创建用户

// ❌ Bad
GET    /api/v1/tenant               // 要使用复数
POST   /api/v1/createTenant         // 要使用 RESTful 风格
GET    /api/v1/tenants/:tenantId    // 参数名使用小写 + 下划线
GET    /api/v1/getTenants           // 不要在路径中使用动词
GET    /api/v1/Tenants              // 不要大写
```

### 查询参数命名

#### 规则

1. **全小写**
2. **使用下划线分隔**

#### 示例

```
# ✅ Good
GET /api/v1/tenants?page=1&page_size=20&status=active&sort_by=created_at&order=desc

# ❌ Bad
GET /api/v1/tenants?page=1&pageSize=20&Status=active&sortBy=createdAt    // 不要使用驼峰或大写
```

---

## 检查清单

### Go 代码检查清单

- [ ] 包名全小写，单个单词
- [ ] 文件名全小写，使用下划线分隔
- [ ] 导出的类型/函数使用 PascalCase
- [ ] 未导出的使用 camelCase
- [ ] 常量使用 UPPER_SNAKE_CASE
- [ ] 接口名以 -er 结尾
- [ ] 错误使用统一错误码

### TypeScript 代码检查清单

- [ ] 组件文件使用 PascalCase
- [ ] 工具/Hooks 文件使用 camelCase
- [ ] 变量使用 camelCase
- [ ] 常量使用 UPPER_SNAKE_CASE
- [ ] 类型/接口使用 PascalCase
- [ ] Props 接口以 Props 结尾

### 数据库检查清单

- [ ] 表名全小写，使用下划线分隔，复数形式
- [ ] 字段名全小写，使用下划线分隔，单数形式
- [ ] 索引命名：`pk_`, `uk_`, `idx_`
- [ ] 外键命名：`fk_`
- [ ] 使用 `deleted_at` 软删除

### API 检查清单

- [ ] 路由全小写，使用连字符分隔
- [ ] 使用复数形式
- [ ] RESTful 风格
- [ ] 查询参数使用下划线分隔

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [api-design.md](api-design.md) | API 设计规范 |
| [database-design.md](database-design.md) | 数据库设计规范 |
| [backend-dev-guide.md](backend-dev-guide.md) | 后端开发指南 |
| [frontend-dev-guide.md](frontend-dev-guide.md) | 前端开发指南 |

---

**🎯 目标**: 统一命名风格，提高代码可读性和可维护性！
