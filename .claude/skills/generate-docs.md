# 文档生成 Skill

## 技能描述

为代码生成完整的文档，包括 Go Doc、JSDoc 注释、README 文件和 API 文档。

## 适用场景

- 新功能完成后需要补充文档
- 代码审查发现文档不完整
- API 接口需要生成文档
- README 需要更新

## 工作流程

### 1. Go 代码文档

#### 函数/方法注释模板

```go
// FunctionName 函数功能描述（简短）
//
// 详细描述函数的功能、使用场景和注意事项。
// 可以包含多行说明。
//
// 参数：
//   ctx - 请求上下文，用于超时控制和取消
//   param1 - 参数1说明
//   param2 - 参数2说明
//
// 返回：
//   *ReturnType - 返回值说明
//   error - 错误信息（如果操作失败）
//
// 错误码：
//   ErrCodeNotFound - 资源未找到时返回
//   ErrCodeInvalidParams - 参数不合法时返回
//
// 示例：
//   // 使用示例
//   result, err := FunctionName(ctx, req)
//   if err != nil {
//       // 处理错误
//   }
//
// 注意：
//   - 注意事项1
//   - 注意事项2
func (s *Service) FunctionName(ctx context.Context, param1 string, param2 int) (*ReturnType, error) {
	// 实现
}
```

#### 接口注释模板

```go
// ServiceName 服务功能描述
//
// 提供哪些能力，负责哪些业务逻辑。
//
// 设计理念：
//   - 为什么这样设计
//   - 设计考虑
//
// 使用方式：
//   service := NewService(components)
//   result, err := service.Method(ctx, req)
type ServiceName interface {
	// MethodName 方法功能描述
	//
	// 参数：
	//   ctx - 上下文
	//   req - 请求参数
	//
	// 返回：
	//   *Response - 响应数据
	//   error - 错误信息
	MethodName(ctx context.Context, req *Request) (*Response, error)
}
```

#### 实体注释模板

```go
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
//
// 索引说明：
//   - idx_email: 邮箱唯一索引
//   - idx_space_id_status: 空间和状态联合索引
type EntityName struct {
	// ID 主键
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// Name 名称
	// 字段详细说明
	Name string `gorm:"column:name;not null" json:"name"`

	// ...
}

// TableName 指定表名
func (EntityName) TableName() string {
	return "entity_names"
}

// BeforeCreate GORM 钩子：创建前
func (e *EntityName) BeforeCreate(tx *gorm.DB) error {
	// 钩子逻辑
	return nil
}
```

### 2. TypeScript 代码文档

#### 函数注释模板

```typescript
/**
 * 函数功能描述（简短）
 *
 * 详细描述函数的功能、使用场景和注意事项。
 *
 * @template T - 泛型参数说明（如果适用）
 * @param param1 - 参数1说明
 * @param param2 - 参数2说明
 * @returns 返回值说明
 * @throws {ErrorType} 可能抛出的错误类型和条件
 *
 * @example
 * ```typescript
 * // 使用示例
 * const result = functionName('param1', 'param2');
 * console.log(result);
 * ```
 *
 * @see {@link https://link} 相关链接
 * @since 1.0.0
 * @deprecated 如果已废弃，说明替代方案
 */
export function functionName<T>(
	param1: string,
	param2: number
): ReturnType {
	// 实现
}
```

#### 组件注释模板

```typescript
/**
 * ComponentName 组件功能描述
 *
 * @description 详细描述组件的功能、用途和使用场景。
 *
 * **Props:**
 * - `prop1` (string): 属性1说明
 * - `prop2` (number, optional): 属性2说明
 *
 * **Usage:**
 * ```tsx
 * <ComponentName prop1="value1" prop2={123} />
 * ```
 *
 * **Features:**
 * - 特性1
 * - 特性2
 *
 * @param props - 组件属性
 * @returns React 元素
 *
 * @example
 * ```tsx
 * // 基本使用
 * <ComponentName title="Hello" />
 *
 * // 高级用法
 * <ComponentName title="Hello" onAction={handleAction} />
 * ```
 */
export function ComponentName(props: ComponentNameProps): JSX.Element {
	// 实现
}
```

#### 接口/类型注释模板

```typescript
/**
 * InterfaceName 接口功能描述
 *
 * 详细描述接口的用途、使用场景和实现要求。
 *
 * @property prop1 - 属性1说明
 * @property prop2 - 属性2说明
 *
 * @example
 * ```typescript
 * const obj: InterfaceName = {
 *   prop1: 'value1',
 *   prop2: 'value2',
 * };
 * ```
 */
export interface InterfaceName {
	/** 属性1说明 */
	prop1: string;

	/** 属性2说明 */
	prop2: number;
}
```

### 3. README 文档生成

#### 包 README 模板

```markdown
# @coze-studio/{package-name}

包功能描述（一句话）。

## 功能特性

- 特性1：详细说明
- 特性2：详细说明
- 特性3：详细说明

## 安装

\`\`\`bash
rush add {package-name}
\`\`\`

## 使用方式

### 基本使用

\`\`\`typescript
import { ComponentName } from '@coze-studio/{package-name}';

function App() {
  return <ComponentName prop="value" />;
}
\`\`\`

### 高级用法

\`\`\`typescript
// 高级用法示例
\`\`\`

## API 文档

详见 [API.md](./docs/API.md)

## 开发指南

### 环境要求
- Node.js >= 18
- React >= 18

### 本地开发
\`\`\`bash
# 安装依赖
rush update

# 启动开发服务器
npm run dev

# 运行测试
npm run test

# 构建
npm run build
\`\`\`

## 注意事项

1. 注意事项1
2. 注意事项2

## 常见问题

### 问题1

**问题描述**: 问题描述

**解决方案**: 解决方案

## 贡献指南

欢迎贡献！请查看 [CONTRIBUTING.md](./CONTRIBUTING.md)

## 许可证

Apache License 2.0
```

### 4. API 文档生成

```markdown
# API 文档：{Service}服务

## 概述

{Service}服务提供{domain}相关的API接口。

## 基础信息

- **基础路径**: `/api/v1/{resources}`
- **认证方式**: Bearer Token
- **响应格式**: JSON

## API 接口

### 1. 创建{Resource}

创建新的{resource}。

**请求**:
- **方法**: POST
- **路径**: `/api/v1/{resources}`
- **Content-Type**: application/json

**请求参数**:
\`\`\`json
{
  "name": "string (必填)",
  "email": "string (必填)",
  "spaceId": "number (必填)"
}
\`\`\`

**响应**:
\`\`\`json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "name": "string",
    "email": "string"
  }
}
\`\`\`

**错误响应**:
\`\`\`json
{
  "code": 20001,
  "message": "{resource}已存在"
}
\`\`\`

**示例**:
\`\`\`bash
curl -X POST https://api.example.com/api/v1/{resources} \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_TOKEN" \\
  -d '{
    "name": "Test",
    "email": "test@example.com",
    "spaceId": 1
  }'
\`\`\`

### 2. 获取{Resource}详情

根据ID获取{resource}详情。

**请求**:
- **方法**: GET
- **路径**: `/api/v1/{resources}/:id`

**路径参数**:
- `id` (number): {resource} ID

**响应**:
\`\`\`json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "name": "string",
    "email": "string",
    "createdAt": "2025-01-01T00:00:00Z"
  }
}
\`\`\`

### 3. 更新{Resource}

更新{resource}信息。

**请求**:
- **方法**: PUT
- **路径**: `/api/v1/{resources}/:id`

**请求参数**:
\`\`\`json
{
  "name": "string (可选)",
  "email": "string (可选)"
}
\`\`\`

### 4. 删除{Resource}

删除{resource}。

**请求**:
- **方法**: DELETE
- **路径**: `/api/v1/{resources}/:id`

**响应**:
\`\`\`json
{
  "code": 0,
  "message": "success"
}
\`\`\`

## 错误码

| 错误码 | 说明 |
|-------|------|
| 10001 | 请求参数不合法 |
| 10002 | 未授权 |
| 20001 | {resource}不存在 |
| 20002 | {resource}已存在 |
```

## 生成流程

### 步骤 1: 分析代码

1. 识别文件类型（Go/TypeScript/React）
2. 识别代码结构（函数/类/接口）
3. 识别依赖关系

### 步骤 2: 提取信息

1. 提取函数/类签名
2. 分析参数和返回值
3. 识别业务逻辑
4. 找出关键注意事项

### 步骤 3: 生成文档

1. 生成注释文档
2. 生成 API 文档
3. 生成 README 文档
4. 生成示例代码

### 步骤 4: 格式化输出

1. 格式化 Markdown
2. 高亮代码块
3. 添加目录和链接

## 输出检查清单

生成文档后，确保：

- [ ] 注释完整准确
- [ ] 参数说明详细
- [ ] 返回值说明清楚
- [ ] 包含使用示例
- [ ] 错误场景说明
- [ ] 注意事项清晰
- [ ] 格式规范统一
- [ ] 代码示例正确
- [ ] 文档结构清晰
- [ ] 更新相关链接

## 注意事项

1. **准确性**: 文档必须与代码保持一致
2. **完整性**: 覆盖所有公开接口
3. **清晰性**: 表达清楚，易于理解
4. **示例性**: 提供实用的示例代码
5. **维护性**: 代码变更时及时更新文档
6. **格式化**: 使用标准格式和模板
