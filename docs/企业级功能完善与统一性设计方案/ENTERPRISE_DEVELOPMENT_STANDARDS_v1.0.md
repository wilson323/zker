# ZKER企业级开发规范与质量标准手册

**版本**: v1.0
**日期**: 2025-01-03
**适用范围**: 前端 + 后端 + 全栈
**目标**: 企业级质量标准

---

## 📋 目录

1. [开发规范](#开发规范)
2. [架构设计规范](#架构设计规范)
3. [代码质量标准](#代码质量标准)
4. [测试规范](#测试规范)
5. [文档规范](#文档规范)
6. [Git工作流规范](#git工作流规范)
7. [CI/CD规范](#cicd规范)
8. [安全规范](#安全规范)
9. [性能优化规范](#性能优化规范)
10. [可访问性规范](#可访问性规范)

---

## 📐 开发规范

### 1. 命名规范

#### TypeScript/JavaScript命名

**文件命名**:
```typescript
// ✅ Good: PascalCase for components and types
InvoiceManagement.tsx
InvoiceManagement.test.tsx
IInvoiceRepository.ts

// ✅ Good: camelCase for utilities and hooks
useInvoiceManagement.ts
invoiceUtils.ts

// ❌ Bad: 不使用kebab-case或snake_case
invoice-management.tsx  // kebab-case
invoice_management.tsx // snake_case
```

**组件命名**:
```typescript
// ✅ Good: PascalCase
const InvoiceManagement: React.FC<Props> = () => {}

// ❌ Bad: camelCase
const invoiceManagement: React.FC<Props> = () => {}
```

**变量/函数命名**:
```typescript
// ✅ Good: camelCase
const invoiceList = []
const fetchInvoices = () => {}

// ❌ Bad: snake_case或PascalCase
const invoice_list = []       // snake_case
const FetchInvoices = () => {} // PascalCase
```

**常量命名**:
```typescript
// ✅ Good: UPPER_SNAKE_CASE
const MAX_PAGE_SIZE = 100
const API_BASE_URL = 'https://api.example.com'

// ❌ Bad: camelCase
const maxPageSize = 100
```

**接口/类型命名**:
```typescript
// ✅ Good: PascalCase, I prefix for interfaces
interface IInvoiceService {}
type InvoiceStatus = 'paid' | 'pending'

// ❌ Bad: camelCase或no prefix
interface invoiceService {} // camelCase
type invoiceStatus = 'paid' | 'pending' // camelCase
```

**枚举命名**:
```typescript
// ✅ Good: PascalCase for enum, UPPER_SNAKE_CASE for values
enum InvoiceStatus {
  PAID = 'PAID',
  PENDING = 'PENDING',
  CANCELLED = 'CANCELLED',
}

// ❌ Bad: camelCase
enum invoiceStatus {
  paid = 'paid',
  pending = 'pending',
}
```

---

#### Go语言命名（后端）

**文件命名**:
```go
// ✅ Good: snake_case
invoice_service.go
invoice_repository.go

// ❌ Bad: camelCase或kebab-case
invoiceService.go  // camelCase
invoice-service.go // kebab-case
```

**包命名**:
```go
// ✅ Good: 简短、小写、单数
package service
package repository
package model

// ❌ Bad: 长名称、复数、大写
package services      // 复数
package models        // 复数
package Service       // 大写
```

**函数命名**:
```go
// ✅ Good: PascalCase（导出）
func GetInvoice(id string) (*Invoice, error) {}
func CreateInvoice(req *CreateInvoiceRequest) (*Invoice, error) {}

// ✅ Good: camelCase（私有）
func validateInvoice(invoice *Invoice) error {}

// ❌ Bad: snake_case
func get_invoice(id string) (*Invoice, error) {} // 私有函数应该camelCase
```

**接口命名**:
```go
// ✅ Good: I prefix + er suffix
type IInvoiceService interface {}
type IInvoiceRepository interface {}

// ✅ Good: 纯接口命名
type InvoiceService interface {}
type InvoiceRepository interface {}

// ❌ Bad: 不一致
type invoiceService interface {} // 小写
type Invoice_Service interface {} // 下划线
```

**常量命名**:
```go
// ✅ Good: PascalCase或UPPER_SNAKE_CASE
const MaxPageSize = 100
const API_BASE_URL = "https://api.example.com"

// ❌ Bad: camelCase
const maxPageSize = 100
```

---

### 2. 代码组织规范

#### 前端项目结构

```typescript
// ✅ 标准页面结构
frontend/packages/studio/src/pages/billing/InvoiceManagement/
├── InvoiceManagement.tsx          // 页面入口
├── components/                    // 页面专属组件
│   ├── InvoiceList.tsx
│   ├── InvoiceDetail.tsx
│   └── InvoiceFilter.tsx
├── hooks/                         // 页面专属hooks
│   └── useInvoiceManagement.ts
├── types/                         // 页面类型定义
│   └── invoice.types.ts
├── constants/                     // 页面常量
│   └── invoice.constants.ts
├── utils/                         // 页面工具函数
│   └── invoice.utils.ts
├── __tests__/                     // 测试文件
│   ├── InvoiceManagement.test.tsx
│   └── __mocks__/
│       └── api.ts
├── index.ts                       // 导出
└── styles.module.less             // 样式文件
```

**导入顺序**:
```typescript
// ✅ Good: 标准导入顺序
// 1. React相关
import React, { useState, useEffect } from 'react'
import { useHistory } from 'react-router-dom'

// 2. 第三方库
import { Button, Table } from '@douyinfe/semi-ui'
import { debounce } from 'lodash-es'

// 3. 项目内部模块（绝对路径）
import { useInvoices } from '@/packages/api-client'
import { InvoiceList } from './components'

// 4. 类型导入
import type { Invoice } from './types'

// 5. 样式文件
import styles from './styles.module.less'

// ❌ Bad: 混乱无序
import styles from './styles.module.less'
import { Button } from '@douyinfe/semi-ui'
import React from 'react'
import { useInvoices } from '@/packages/api-client'
```

---

#### 后端项目结构

```go
// ✅ 标准DDD结构
backend/domain/
├── invoice/                      // 领域
│   ├── entity/                   // 实体
│   │   └── invoice.go
│   ├── repository/               // 仓储接口
│   │   └── invoice_repository.go
│   └── service/                  // 领域服务
│       └── invoice_service.go

backend/application/
├── invoice/                      // 应用服务
│   └── invoice_app_service.go

backend/api/handler/coze/
├── invoice_service.go            // HTTP处理器

backend/infra/
├── repository/                   // 仓储实现
│   └── invoice_repository_impl.go
└── cache/                        // 基础设施
    └── redis_cache.go

backend/tests/
├── invoice_service_test.go       // 测试
```

**导入顺序**:
```go
// ✅ Good: 标准导入顺序
// 1. 标准库
import (
	"context"
	"fmt"
	"time"
)

// 2. 项目内部模块
import (
	"github.com/coze-dev/coze-studio/backend/domain/invoice/entity"
	"github.com/coze-dev/coze-studio/backend/domain/invoice/repository"
)

// 3. 第三方库
import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/samber/lo/v3"
)

// ❌ Bad: 混乱无序
import (
	"github.com/cloudwego/hertz/pkg/app"
	"context"
	"fmt"
	"github.com/samber/lo/v3"
	"github.com/coze-dev/coze-studio/backend/domain/invoice/entity"
)
```

---

### 3. 注释规范

#### JSDoc注释（前端）

```typescript
/**
 * 发票管理页面
 *
 * @description 用于管理和查看发票信息
 *
 * @author Frontend Team
 * @since 2025-01-03
 *
 * @example
 * ```tsx
 * import { InvoiceManagement } from '@/pages/billing/InvoiceManagement'
 *
 * <InvoiceManagement />
 * ```
 *
 * @features
 * - 发票列表展示
 * - 高级筛选
 * - 批量导出
 *
 * @permissions
 * - billing.invoices.read - 查看发票
 * - billing.invoices.export - 导出发票
 */
export const InvoiceManagement: React.FC<InvoiceManagementProps> = () => {
  // ...
}

/**
 * 获取发票列表
 *
 * @param params - 查询参数
 * @param params.page - 页码（从1开始）
 * @param params.pageSize - 每页数量
 * @param params.status - 发票状态
 * @returns 发票列表响应
 *
 * @throws {InvoiceError} 当API调用失败时
 *
 * @example
 * ```typescript
 * const invoices = await fetchInvoices({
 *   page: 1,
 *   pageSize: 20,
 *   status: 'paid'
 * })
 * ```
 */
export const fetchInvoices = async (
  params: InvoiceListParams
): Promise<InvoiceListResponse> => {
  // ...
}
```

---

#### Go注释（后端）

```go
// Package service 提供发票领域服务
//
// 该包实现了发票相关的业务逻辑，包括：
// - 发票创建
// - 发票查询
// - 发票状态更新
//
// 示例：
//
//	svc := service.NewInvoiceService(repo)
//	invoice, err := svc.CreateInvoice(ctx, req)
package service

// InvoiceService 发票服务接口
//
// 定义了发票相关的业务操作方法
type InvoiceService interface {
	// CreateInvoice 创建发票
	//
	// 参数：
	//   ctx - 上下文
	//   req - 创建发票请求
	//
	// 返回：
	//   *entity.Invoice - 创建的发票实体
	//   error - 错误信息
	//
	// 错误：
	//   errno.ErrInvalidParam - 参数无效
	//   errno.ErrInvoiceDuplicate - 发票重复
	CreateInvoice(ctx context.Context, req *CreateInvoiceRequest) (*entity.Invoice, error)

	// GetInvoice 获取发票详情
	GetInvoice(ctx context.Context, id string) (*entity.Invoice, error)
}

// NewInvoiceService 创建发票服务实例
func NewInvoiceService(repo repository.IInvoiceRepository) InvoiceService {
	return &invoiceServiceImpl{
		repo: repo,
	}
}
```

---

## 🏗️ 架构设计规范

### 1. DDD分层架构

#### 四层架构

```typescript
// ✅ 标准DDD分层
┌─────────────────────────────────────┐
│   API Layer (Handler)               │  ← HTTP接口、参数验证
│   - controller/                     │
│   - middleware/                     │
├─────────────────────────────────────┤
│   Application Layer (UseCase)       │  ← 用例编排、事务管理
│   - service/                        │
│   - dto/                            │
├─────────────────────────────────────┤
│   Domain Layer (Entity)             │  ← 核心业务、领域规则
│   - entity/                         │
│   - valueobject/                    │
│   - repository/                     │
├─────────────────────────────────────┤
│   Infrastructure Layer (Impl)        │  ← 技术实现、外部依赖
│   - repository/impl/                │
│   - cache/                          │
│   - mq/                             │
└─────────────────────────────────────┘
```

**依赖规则**:
- ✅ 上层可以依赖下层
- ❌ 下层不能依赖上层
- ✅ 同层之间不直接依赖（通过接口）

---

### 2. SOLID原则

#### S - 单一职责原则（Single Responsibility）

```typescript
// ✅ Good: 单一职责
class InvoiceList {
  // 只负责展示发票列表
}

class InvoiceFilter {
  // 只负责筛选功能
}

// ❌ Bad: 多职责
class InvoiceManagement {
  // 既展示列表，又筛选，又创建，又编辑...
  // 违反单一职责原则
}
```

#### O - 开放封闭原则（Open-Closed）

```typescript
// ✅ Good: 对扩展开放，对修改封闭
interface INotificationSender {
  send(message: string): void
}

class EmailNotificationSender implements INotificationSender {
  send(message: string): void {
    // Email发送逻辑
  }
}

class SmsNotificationSender implements INotificationSender {
  send(message: string): void {
    // SMS发送逻辑
  }
}

// 新增通知方式只需实现接口，无需修改现有代码
class PushNotificationSender implements INotificationSender {
  send(message: string): void {
    // Push发送逻辑
  }
}

// ❌ Bad: 每次新增都需要修改现有代码
class NotificationSender {
  send(type: string, message: string): void {
    if (type === 'email') {
      // Email逻辑
    } else if (type === 'sms') {
      // SMS逻辑
    } else if (type === 'push') {
      // 新增需要修改这里
    }
  }
}
```

#### L - 里氏替换原则（Liskov Substitution）

```typescript
// ✅ Good: 子类可以替换父类
class Rectangle {
  setWidth(width: number): void {}
  setHeight(height: number): void {}
  getArea(): number { return 0 }
}

class Square extends Rectangle {
  // 正方形可以替换矩形
}

// ❌ Bad: 子类不能随意替换父类
class Bird {
  fly(): void {
    console.log('Flying')
  }
}

class Penguin extends Bird {
  fly(): void {
    throw new Error('Penguins cannot fly')
    // 企鹅不能飞，违反里氏替换原则
  }
}
```

#### I - 接口隔离原则（Interface Segregation）

```typescript
// ✅ Good: 接口职责单一
interface IReadable {
  read(): void
}

interface IWritable {
  write(data: string): void
}

class File implements IReadable, IWritable {
  read(): void {}
  write(data: string): void {}
}

// ❌ Bad: 胖接口（Fat Interface）
interface IFileOperation {
  read(): void
  write(data: string): void
  delete(): void
  copy(): void
  move(): void
  compress(): void
  encrypt(): void
  // ... 太多方法
}

class ReadOnlyFile implements IFileOperation {
  read(): void {}
  write(data: string): void {
    throw new Error('Read-only file')
    // 只读文件不需要写方法，但必须实现
  }
  delete(): void { throw new Error('Read-only file') }
  // ... 很多不需要的方法
}
```

#### D - 依赖倒置原则（Dependency Inversion）

```typescript
// ✅ Good: 依赖抽象而非具体实现
interface IInvoiceRepository {
  findById(id: string): Promise<Invoice>
  save(invoice: Invoice): Promise<void>
}

class InvoiceService {
  constructor(
    private readonly repo: IInvoiceRepository  // 依赖接口
  ) {}

  async getInvoice(id: string): Promise<Invoice> {
    return this.repo.findById(id)
  }
}

// ❌ Bad: 依赖具体实现
class InvoiceService {
  constructor(
    private readonly repo: MongoInvoiceRepository  // 依赖具体实现
  ) {}

  // 如果要切换数据库，需要修改这个类
}
```

---

### 3. KISS原则（Keep It Simple, Stupid）

```typescript
// ✅ Good: 简单直接
const formatDate = (date: Date): string => {
  return date.toLocaleDateString('zh-CN')
}

// ❌ Bad: 过度复杂
const formatDate = (date: Date): string => {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()

  // 不必要的复杂性
  const yearStr = year.toString()
  const monthStr = month < 10 ? `0${month}` : month.toString()
  const dayStr = day < 10 ? `0${day}` : day.toString()

  return `${yearStr}-${monthStr}-${dayStr}`
}
```

---

### 4. DRY原则（Don't Repeat Yourself）

```typescript
// ❌ Bad: 重复代码
const getUserById = async (id: string) => {
  const response = await fetch(`/api/users/${id}`)
  const data = await response.json()
  return data
}

const getInvoiceById = async (id: string) => {
  const response = await fetch(`/api/invoices/${id}`)
  const data = await response.json()
  return data
}

// ✅ Good: 提取可复用逻辑
const fetchById = async (resource: string, id: string) => {
  const response = await fetch(`/api/${resource}/${id}`)
  const data = await response.json()
  return data
}

const getUserById = (id: string) => fetchById('users', id)
const getInvoiceById = (id: string) => fetchById('invoices', id)
```

---

### 5. YAGNI原则（You Aren't Gonna Need It）

```typescript
// ❌ Bad: 提前实现不需要的功能
class InvoiceService {
  async createInvoice(req: CreateInvoiceRequest): Promise<Invoice> {
    // ...
  }

  // 这些功能现在不需要，但被提前实现了
  async archiveInvoice(id: string): Promise<void> {}
  async restoreInvoice(id: string): Promise<void> {}
  async duplicateInvoice(id: string): Promise<Invoice> {}
  async mergeInvoices(ids: string[]): Promise<Invoice> {}
}

// ✅ Good: 只实现当前需要的功能
class InvoiceService {
  async createInvoice(req: CreateInvoiceRequest): Promise<Invoice> {
    // ...
  }

  // 需要时再添加
}
```

---

## 📊 代码质量标准

### 1. 复杂度控制

#### 圈复杂度（Cyclomatic Complexity）

```typescript
// ✅ Good: 圈复杂度 ≤10
function calculateInvoiceTotal(items: InvoiceItem[]): number {
  let total = 0

  for (const item of items) {
    if (item.type === 'product') {
      total += item.price * item.quantity
    } else if (item.type === 'service') {
      total += item.price
    } else if (item.type === 'discount') {
      total -= item.price
    }
  }

  return total
}

// ❌ Bad: 圈复杂度 >10
function calculateInvoiceTotal(items: InvoiceItem[]): number {
  let total = 0

  for (const item of items) {
    if (item.type === 'product') {
      if (item.category === 'electronics') {
        if (item.brand === 'apple') {
          if (item.model === 'iphone') {
            // 太多层嵌套
          }
        }
      }
    } else if (item.type === 'service') {
      // ...
    } else if (item.type === 'discount') {
      // ...
    }
  }

  return total
}
```

---

### 2. 代码行数限制

```typescript
// ✅ Good: 函数 ≤50行
const processInvoice = (invoice: Invoice): ProcessedInvoice => {
  // 1. 验证
  if (!validateInvoice(invoice)) {
    throw new Error('Invalid invoice')
  }

  // 2. 计算总额
  const total = calculateTotal(invoice.items)

  // 3. 应用折扣
  const discounted = applyDiscount(total, invoice.discount)

  // 4. 返回结果
  return {
    ...invoice,
    total,
    discounted
  }
}

// ❌ Bad: 函数 >50行
const processInvoice = (invoice: Invoice): ProcessedInvoice => {
  // 100多行代码...
  // 难以阅读和维护
}
```

---

### 3. 魔法值消除

```typescript
// ❌ Bad: 魔法值
if (invoice.status === 1) {
  // ...
}

if (user.role === 3) {
  // ...
}

const pageSize = 20

// ✅ Good: 使用常量
enum InvoiceStatus {
  DRAFT = 1,
  PENDING = 2,
  PAID = 3,
}

enum UserRole {
  ADMIN = 1,
  USER = 2,
  GUEST = 3,
}

const DEFAULT_PAGE_SIZE = 20
const MAX_PAGE_SIZE = 100

if (invoice.status === InvoiceStatus.PAID) {
  // ...
}

if (user.role === UserRole.ADMIN) {
  // ...
}

const pageSize = DEFAULT_PAGE_SIZE
```

---

## 🧪 测试规范

### 1. 单元测试

#### 测试覆盖率要求

| 代码类型 | 覆盖率要求 | 说明 |
|---------|-----------|------|
| 业务逻辑 | ≥90% | 关键业务逻辑必须高覆盖 |
| 工具函数 | ≥80% | 可复用工具函数 |
| 组件 | ≥80% | UI组件 |
| API层 | ≥70% | HTTP接口 |

---

#### 测试命名规范

```typescript
// ✅ Good: 清晰的测试命名
describe('InvoiceService', () => {
  describe('createInvoice', () => {
    it('should create invoice when valid request', async () => {
      // Given
      const req = createValidRequest()

      // When
      const result = await service.createInvoice(req)

      // Then
      expect(result).toBeDefined()
      expect(result.status).toBe(InvoiceStatus.PENDING)
    })

    it('should throw error when invalid request', async () => {
      // Given
      const req = createInvalidRequest()

      // When & Then
      await expect(
        service.createInvoice(req)
      ).rejects.toThrow('Invalid request')
    })
  })
})

// ❌ Bad: 不清晰的测试命名
describe('InvoiceService', () => {
  it('test1', async () => {
    // test1是什么？
  })

  it('test2', async () => {
    // test2是什么？
  })
})
```

---

### 2. 集成测试

```typescript
// ✅ Good: 集成测试示例
describe('InvoiceManagement Integration', () => {
  it('should create and display invoice', async () => {
    // 1. 渲染页面
    render(<InvoiceManagement />)

    // 2. 点击创建按钮
    fireEvent.click(screen.getByText('创建发票'))

    // 3. 填写表单
    fillInvoiceForm({
      customer: 'Test Customer',
      amount: 1000,
    })

    // 4. 提交表单
    fireEvent.click(screen.getByText('提交'))

    // 5. 验证结果
    await waitFor(() => {
      expect(screen.getByText('创建成功')).toBeInTheDocument()
      expect(screen.getByText('Test Customer')).toBeInTheDocument()
    })
  })
})
```

---

### 3. E2E测试

```typescript
// ✅ Good: E2E测试示例
test('完整发票创建流程', async ({ page }) => {
  // 1. 登录
  await page.goto('http://localhost:8888/login')
  await page.fill('#username', 'admin')
  await page.fill('#password', 'password')
  await page.click('button[type="submit"]')

  // 2. 导航到发票管理
  await page.click('text=发票管理')

  // 3. 创建发票
  await page.click('text=创建发票')
  await page.fill('#customer', 'Test Customer')
  await page.fill('#amount', '1000')
  await page.click('button:has-text("提交")')

  // 4. 验证
  await expect(page.locator('text=创建成功')).toBeVisible()
  await expect(page.locator('text=Test Customer')).toBeVisible()
})
```

---

## 📖 文档规范

### 1. API文档

```typescript
/**
 * 获取发票列表
 *
 * @description 根据查询参数获取发票列表，支持分页和筛选
 *
 * @param {InvoiceListParams} params - 查询参数
 * @param {number} params.page - 页码（从1开始）
 * @param {number} params.pageSize - 每页数量（1-100）
 * @param {InvoiceStatus} params.status - 发票状态筛选
 * @param {string} params.startDate - 开始日期（YYYY-MM-DD）
 * @param {string} params.endDate - 结束日期（YYYY-MM-DD）
 *
 * @returns {Promise<InvoiceListResponse>} 发票列表响应
 * @returns {Invoice[]} InvoiceListResponse.invoices - 发票列表
 * @returns {number} InvoiceListResponse.total - 总数
 * @returns {number} InvoiceListResponse.page - 当前页码
 *
 * @throws {InvoiceError} 当API调用失败时抛出
 * @throws {NetworkError} 当网络错误时抛出
 *
 * @example
 * ```typescript
 * const result = await getInvoices({
 *   page: 1,
 *   pageSize: 20,
 *   status: InvoiceStatus.PAID
 * })
 *
 * console.log(result.invoices) // 发票列表
 * console.log(result.total)    // 总数
 * ```
 *
 * @see {@link https://api.example.com/docs/invoices} API文档
 * @since 2025-01-03
 */
export const getInvoices = async (
  params: InvoiceListParams
): Promise<InvoiceListResponse> => {
  // implementation
}
```

---

### 2. 组件文档

```typescript
/**
 * InvoiceManagement 组件
 *
 * @description 发票管理页面，提供发票列表展示、筛选、导出等功能
 *
 * @component
 *
 * @example
 * ```tsx
 * import { InvoiceManagement } from '@/pages/billing/InvoiceManagement'
 *
 * export default function App() {
 *   return <InvoiceManagement />
 * }
 * ```
 *
 * @features
 * - 📋 发票列表展示（表格/卡片视图）
 * - 🔍 高级筛选（日期范围、状态、金额）
 * - 📤 批量导出（CSV、Excel、PDF）
 * - 📄 发票详情查看
 * - ⬇️ 发票下载（PDF）
 *
 * @permissions
 * - billing.invoices.read - 查看发票
 * - billing.invoices.export - 导出发票
 *
 * @apis
 * - GET /api/billing/invoices - 获取发票列表
 * - GET /api/billing/invoices/:id - 获取发票详情
 * - POST /api/billing/invoices/export - 导出发票
 *
 * @author Frontend Team
 * @since 2025-01-03
 * @version 1.0.0
 */
export const InvoiceManagement: React.FC<InvoiceManagementProps> = () => {
  // implementation
}
```

---

## 🌳 Git工作流规范

### 1. 分支策略

```bash
main          # 主分支（生产环境）
├── develop   # 开发分支
│   ├── feature/billing-invoice-management    # 功能分支
│   ├── feature/audit-log-search             # 功能分支
│   ├── feature/org-management               # 功能分支
│   ├── bugfix/invoice-calculation-error     # 修复分支
│   └── hotfix/security-patch               # 紧急修复分支
```

**分支命名规范**:
- `feature/功能名称` - 新功能开发
- `bugfix/问题描述` - Bug修复
- `hotfix/紧急问题` - 生产环境紧急修复
- `refactor/重构内容` - 代码重构
- `docs/文档内容` - 文档更新
- `test/测试内容` - 测试相关
- `chore/杂项` - 构建、工具等

---

### 2. Commit规范

#### Commit Message格式

```bash
<type>(<scope>): <subject>

<body>

<footer>
```

**Type类型**:
- `feat` - 新功能
- `fix` - Bug修复
- `docs` - 文档更新
- `style` - 代码格式调整（不影响功能）
- `refactor` - 重构（既不是新功能也不是修复）
- `perf` - 性能优化
- `test` - 测试相关
- `chore` - 构建、工具等
- `ci` - CI/CD相关

**示例**:
```bash
# ✅ Good: 清晰的Commit Message
feat(billing): add invoice management page

- Implement invoice list with pagination
- Add invoice detail view
- Add invoice export functionality (CSV, Excel, PDF)
- Add invoice filter (date range, status, amount)

Closes #123

# ❌ Bad: 不清晰的Commit Message
add page
fix bug
update
```

---

## 🔒 安全规范

### 1. 前端安全

```typescript
// ✅ Good: XSS防护
const UserInput: React.FC = () => {
  const [content, setContent] = useState('')

  return (
    <div>
      {/* React自动转义，防止XSS */}
      <p>{content}</p>
    </div>
  )
}

// ❌ Bad: 危险的HTML插入
const UserInput: React.FC = () => {
  const [content, setContent] = useState('')

  return (
    <div>
      {/* 危险：直接插入HTML */}
      <div dangerouslySetInnerHTML={{ __html: content }} />
    </div>
  )
}

// ✅ Good: 使用DOMPurify清理HTML
import DOMPurify from 'dompurify'

const SafeHTML: React.FC = ({ html }) => {
  const cleanHtml = DOMPurify.sanitize(html)
  return <div dangerouslySetInnerHTML={{ __html: cleanHtml }} />
}
```

---

### 2. 后端安全

```go
// ✅ Good: 参数化查询，防止SQL注入
func GetInvoice(id string) (*Invoice, error) {
	var invoice Invoice
	err := db.Where("id = ?", id).First(&invoice).Error
	return &invoice, err
}

// ❌ Bad: 字符串拼接，SQL注入风险
func GetInvoice(id string) (*Invoice, error) {
	var invoice Invoice
	err := db.Raw("SELECT * FROM invoices WHERE id = '" + id + "'").First(&invoice).Error
	return &invoice, err
}

// ✅ Good: 输入验证
func CreateInvoice(req *CreateInvoiceRequest) error {
	// 验证输入
	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if req.CustomerID == "" {
		return errors.New("customer_id is required")
	}

	// ... 创建逻辑
	return nil
}
```

---

## ⚡ 性能优化规范

### 1. 前端性能

```typescript
// ✅ Good: 使用React.memo避免不必要的重渲染
const InvoiceItem = React.memo<InvoiceItemProps>(({ invoice, onSelect }) => {
  return (
    <div onClick={() => onSelect(invoice.id)}>
      {invoice.name}
    </div>
  )
})

// ✅ Good: 使用useMemo缓存计算结果
const InvoiceList: React.FC = ({ invoices }) => {
  const total = useMemo(() => {
    return invoices.reduce((sum, inv) => sum + inv.amount, 0)
  }, [invoices])

  return <div>Total: {total}</div>
}

// ✅ Good: 使用useCallback缓存函数
const InvoiceList: React.FC = () => {
  const handleSelect = useCallback((id: string) => {
    console.log('Selected:', id)
  }, [])

  return <InvoiceItem id="1" onSelect={handleSelect} />
}

// ✅ Good: 代码分割
const InvoiceDetail = React.lazy(() => import('./InvoiceDetail'))

const App: React.FC = () => {
  return (
    <Suspense fallback={<Loading />}>
      <InvoiceDetail />
    </Suspense>
  )
}
```

---

### 2. 后端性能

```go
// ✅ Good: 使用缓存
func (s *invoiceService) GetInvoice(ctx context.Context, id string) (*Invoice, error) {
	// 1. 先查缓存
	if cached, found := s.cache.Get(id); found {
		return cached.(*Invoice), nil
	}

	// 2. 查数据库
	invoice, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. 写入缓存
	s.cache.Set(id, invoice, 5*time.Minute)

	return invoice, nil
}

// ✅ Good: 批量查询，避免N+1
func (s *invoiceService) GetInvoices(ctx context.Context, ids []string) ([]*Invoice, error) {
	// ❌ Bad: N+1查询
	// for _, id := range ids {
	//   invoice, _ := s.repo.FindByID(ctx, id)
	//   invoices = append(invoices, invoice)
	// }

	// ✅ Good: 批量查询
	return s.repo.FindByIDs(ctx, ids)
}
```

---

## ♿ 可访问性规范

### 1. 键盘导航

```typescript
// ✅ Good: 所有交互支持键盘
const InvoiceButton: React.FC = () => {
  const handleClick = () => {
    console.log('Clicked')
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault()
      handleClick()
    }
  }

  return (
    <button
      onClick={handleClick}
      onKeyDown={handleKeyDown}
      tabIndex={0}  // 可聚焦
    >
      点击我
    </button>
  )
}
```

---

### 2. ARIA标签

```typescript
// ✅ Good: 完整的ARIA标签
const InvoiceModal: React.FC = ({ isOpen, onClose }) => {
  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
      aria-describedby="modal-description"
    >
      <h2 id="modal-title">发票详情</h2>
      <p id="modal-description">查看发票详细信息</p>

      <button
        onClick={onClose}
        aria-label="关闭对话框"
      >
        关闭
      </button>
    </div>
  )
}
```

---

### 3. 颜色对比度

```less
// ✅ Good: 符合WCAG AA标准（4.5:1）
.invoice-text {
  color: #333333;  // 深灰色文字
  background: #ffffff;  // 白色背景
  // 对比度: 12.6:1 ✅
}

.invoice-button {
  color: #ffffff;  // 白色文字
  background: #0066cc;  // 蓝色背景
  // 对比度: 5.9:1 ✅
}

// ❌ Bad: 不符合标准
.low-contrast {
  color: #cccccc;  // 浅灰色文字
  background: #eeeeee;  // 浅灰色背景
  // 对比度: 1.6:1 ❌
}
```

---

## 🎯 质量门禁

### 提交前检查（Pre-commit）

```json
{
  "scripts": {
    "precommit": "npm run lint && npm run type-check && npm run test -- --passWithNoTests",
    "lint": "eslint . --ext .ts,.tsx --fix",
    "type-check": "tsc --noEmit",
    "test": "jest --coverage",
    "format": "prettier --write \"**/*.{ts,tsx,json,md}\""
  }
}
```

---

### CI/CD质量门禁

```yaml
# .github/workflows/quality-gate.yml
name: Quality Gate

on:
  pull_request:
    branches: [main, develop]

jobs:
  quality:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install dependencies
        run: npm ci

      # 质量门禁1: Lint
      - name: Lint
        run: npm run lint
        continue-on-error: false  # 必须通过

      # 质量门禁2: TypeScript
      - name: Type check
        run: npm run type-check
        continue-on-error: false  # 必须通过

      # 质量门禁3: 测试覆盖率
      - name: Test with coverage
        run: npm run test:coverage
        continue-on-error: false  # 必须通过

      # 质量门禁4: 构建
      - name: Build
        run: npm run build
        continue-on-error: false  # 必须通过
```

---

## 📊 质量评分卡

### 每个页面的质量评分

| 维度 | 满分 | 评分标准 | 实际得分 |
|------|------|----------|----------|
| **代码质量** | 20 | Lint 0错误，TS 0错误 | __/20 |
| **测试覆盖** | 20 | 覆盖率≥80% | __/20 |
| **性能** | 15 | LCP<2.5s, FID<100ms, CLS<0.1 | __/15 |
| **可访问性** | 15 | WCAG 2.1 AA级 | __/15 |
| **文档完整** | 10 | JSDoc 100% | __/10 |
| **国际化** | 10 | i18n 100% | __/10 |
| **安全性** | 10 | 无安全漏洞 | __/10 |
| **总计** | **100** | **≥90分通过** | **__/100** |

---

## 🎊 总结

### 核心原则

1. **SOLID原则**:
   - S: 单一职责
   - O: 开放封闭
   - L: 里氏替换
   - I: 接口隔离
   - D: 依赖倒置

2. **KISS原则**: 保持简单，避免过度设计

3. **DRY原则**: 消除重复，提取可复用代码

4. **YAGNI原则**: 只实现当前需要的功能

### 质量标准

- ✅ 代码质量: Lint 0错误，TypeScript 0错误
- ✅ 测试覆盖: ≥80%
- ✅ 性能: LCP<2.5s, FID<100ms, CLS<0.1
- ✅ 可访问性: WCAG 2.1 AA级
- ✅ 安全性: 无已知漏洞

### 成功关键

- 🎯 **质量优先**: 不牺牲质量换取速度
- 📚 **文档先行**: 代码与文档同步更新
- 🧪 **测试驱动**: 编写单元测试
- 🔒 **安全第一**: 遵循安全最佳实践
- ⚡ **性能为王**: 持续优化性能

---

**最后更新**: 2025-01-03
**执行开始**: 立即生效
**适用范围**: 全体开发人员
