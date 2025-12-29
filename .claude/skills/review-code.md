# 代码审查 Skill

## 技能描述

对代码进行深入审查，评估功能性、规范性、可维护性、性能、安全性和测试覆盖。

## 适用场景

- Pull Request 审查
- 代码重构评审
- 架构设计评审
- 代码质量检查

## 审查维度

### 1. 功能性审查

**检查内容**:
- ✅ 是否实现了需求文档中的所有功能
- ✅ 边界条件是否得到处理
- ✅ 错误场景是否被覆盖
- ✅ 业务规则是否正确实现

**审查示例**:

```typescript
// ✅ 良好的边界处理
function divide(a: number, b: number): number {
  if (b === 0) {
    throw new Error('Division by zero');
  }
  return a / b;
}

// ❌ 缺少边界处理
function divide(a: number, b: number): number {
  return a / b;  // 没有处理除零情况
}
```

**审查要点**:
```
□ 功能完整性
□ 边界条件处理
□ 错误场景覆盖
□ 业务规则正确性
```

### 2. 规范性审查

**检查内容**:
- ✅ 命名规范（文件、变量、函数）
- ✅ 代码结构（组件分层、DDD 分层）
- ✅ 代码风格（缩进、空行、注释）
- ✅ 导入顺序
- ✅ TypeScript/Go 类型使用

**审查示例**:

```typescript
// ✅ 符合规范
function getUserById(id: number): User | null {}
const MAX_RETRY_COUNT = 3;
interface UserProfileProps {}

// ❌ 不符合规范
function user(id: number) {}         ← 应使用 getUserById
const max_retry_count = 3;          ← 应使用 UPPER_SNAKE_CASE
interface Props {}                    ← 过于通用
```

**审查要点**:
```
□ 文件命名规范
□ 变量命名规范
□ 函数命名规范
□ 代码结构规范
□ 导入顺序规范
□ 类型注解完整性
```

### 3. 可维护性审查

**检查内容**:
- ✅ 代码是否清晰易懂
- ✅ 是否遵循 DRY 原则（避免重复）
- ✅ 函数是否职责单一
- ✅ 是否有适当的注释
- ✅ 是否易于测试

**审查示例**:

```typescript
// ✅ 易于维护
function getUserById(id: number): User | null {
  // 清晰的函数名和实现
  return userRepository.findById(id);
}

// ❌ 难以维护
function process(data) {  ← 函数名不明确
  const result = [];
  for (let i = 0; i < data.length; i++) {
    result.push(data[i].value * 2);  ← 魔法数字
  }
  return result;
}
```

**审查要点**:
```
□ 代码可读性
□ 函数职责单一
□ 避免代码重复
□ 有意义的变量名
□ 适当的注释
□ 易于测试
```

### 4. 性能审查

**检查内容**:
- ✅ 是否有不必要的重渲染
- ✅ 是否有内存泄漏
- ✅ 是否有 N+1 查询
- ✅ 是否使用了缓存
- ✅ 是否有大循环或递归

**审查示例**:

```typescript
// ✅ 性能优化
const memoizedValue = useMemo(() => {
  return expensiveComputation(data);
}, [data]);

const handleClick = useCallback(() => {
  onClick();
}, [onClick]);

// ❌ 性能问题
function Component() {
  const value = expensiveComputation(data);  ← 每次渲染都计算
  return <div onClick={() => onClick()}>Click</div>;  ← 每次渲染创建新函数
}
```

**审查要点**:
```
□ 使用 useMemo 缓存计算
□ 使用 useCallback 包装回调
□ 避免 N+1 查询
□ 列表渲染使用 key
□ 懒加载或代码分割
□ 防抖/节流使用
```

### 5. 安全性审查

**检查内容**:
- ✅ 是否有 SQL 注入风险
- ✅ 是否有 XSS 风险
- ✅ 是否验证用户输入
- ✅ 是否有权限检查
- ✅ 是否有敏感信息泄露

**审查示例**:

```go
// ✅ 安全的查询
db.Where("name = ?", name).Find(&users)  ← 使用参数化查询

// ❌ SQL 注入风险
db.Raw(fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name))  ← 不应拼接 SQL

// ✅ 输入验证
if req.Name == "" {
    return errors.New("name is required")
}

// ❌ 缺少验证
user.Name = req.Name  ← 没有验证
```

**审查要点**:
```
□ 参数化查询
□ 输入验证
□ 权限检查
□ 敏感数据处理
□ 错误信息不泄露细节
□ HTTPS 使用
```

### 6. 测试审查

**检查内容**:
- ✅ 是否有单元测试
- ✅ 测试覆盖率是否达标
- ✅ 是否测试边界情况
- ✅ 是否测试错误场景
- ✅ 测试是否易于维护

**审查要点**:
```
□ 单元测试存在
□ 测试覆盖率达标（Level 1: 80%, Level 2: 30%）
□ 正常场景测试
□ 边界情况测试
□ 错误场景测试
□ Mock 使用正确
□ 测试独立可重复
```

### 7. 文档审查

**检查内容**:
- ✅ 是否有函数/方法注释
- ✅ 是否有参数说明
- ✅ 是否有返回值说明
- ✅ 是否有使用示例
- ✅ README 是否完整

**审查要点**:
```
□ 函数注释完整
□ 参数说明清晰
□ 返回值说明完整
□ 复杂逻辑有注释
□ API 文档完整
□ README 准确
```

## 审查流程

### 步骤 1: 自动化检查

```bash
# 前端检查
npm run lint
npm run test
npm run test:cov

# 后端检查
gofmt ./...
go test ./...
golangci-lint run
```

### 步骤 2: 代码阅读

1. 快速浏览文件结构
2. 阅读关键函数/方法
3. 检查命名规范
4. 检查代码结构

### 步骤 3: 深入分析

1. 分析代码逻辑
2. 检查错误处理
3. 评估性能影响
4. 评估安全性

### 步骤 4: 生成报告

根据审查结果生成详细报告。

## 审查报告格式

```markdown
# 代码审查报告

## 审查概述
- **审查对象**: {文件路径或 PR 编号}
- **审查时间**: 2025-01-01 12:00:00
- **审查人**: AI Reviewer
- **总体评分**: ⭐⭐⭐⭐☆ (4/5)
- **审查结论**: ✅ 通过 / ⚠️ 需要改进 / ❌ 拒绝

## 审查详情

### ✅ 优点 (5)
1. 命名规范，易于理解
2. 错误处理完善
3. 代码结构清晰
4. 有完整的单元测试
5. 性能优化得当

### ⚠️ 需要改进 (3)

#### 1. 函数命名不够明确
- **位置**: `{file}:{line}`
- **问题**: 函数名 `process` 不够明确
- **建议**: 改为 `processUserData`
- **优先级**: 中
- **代码**:
  ```typescript
  function process(data) {
    // ...
  }
  ```
- **修复**:
  ```typescript
  function processUserData(data: UserData) {
    // ...
  }
  ```

#### 2. 缺少错误处理
- **位置**: `{file}:{line}`
- **问题**: 未检查函数返回的错误
- **建议**: 添加错误检查
- **优先级**: 高
- **代码**:
  ```go
  user, _ := repo.Create(ctx, req)
  ```
- **修复**:
  ```go
  user, err := repo.Create(ctx, req)
  if err != nil {
      return fmt.Errorf("create user failed: %w", err)
  }
  ```

#### 3. 缺少类型注解
- **位置**: `{file}:{line}`
- **问题**: 函数参数缺少类型
- **建议**: 添加类型注解
- **优先级**: 低
- **代码**:
  ```typescript
  function getUser(id) {
    return users.find(u => u.id === id);
  }
  ```
- **修复**:
  ```typescript
  function getUser(id: number): User | undefined {
    return users.find(u => u.id === id);
  }
  ```

### ❌ 严重问题 (0)

无严重问题。

## 分项评分

| 维度 | 评分 | 说明 |
|-----|------|------|
| 功能性 | 5/5 | 功能完整，边界处理得当 |
| 规范性 | 4/5 | 整体符合规范，少量命名需改进 |
| 可维护性 | 4/5 | 代码清晰，部分函数可拆分 |
| 性能 | 5/5 | 性能优化良好，无性能问题 |
| 安全性 | 5/5 | 安全措施完善，无明显漏洞 |
| 测试 | 4/5 | 测试覆盖率高，边界测试可补充 |
| 文档 | 3/5 | 注释完整，但缺少使用示例 |

## 改进建议

### 短期改进（本次 PR）
1. 修复函数命名问题（3 处）
2. 补充错误处理（2 处）
3. 添加类型注解（5 处）

### 长期改进
1. 考虑拆分复杂函数
2. 补充集成测试
3. 添加使用文档

## 审查结论

**✅ 通过，建议修改**

代码整体质量良好，有少量需要改进的地方。建议作者根据反馈进行修改后合并。

**优先级**:
- 🔴 高: 必须修改
- 🟡 中: 建议修改
- 🔵 低: 可选修改
```

## 审查检查清单

```markdown
## 功能性
- [ ] 实现了所有需求功能
- [ ] 边界条件处理正确
- [ ] 错误场景覆盖完整
- [ ] 业务规则实现正确

## 规范性
- [ ] 文件命名符合规范
- [ ] 变量命名符合规范
- [ ] 函数命名符合规范
- [ ] 代码结构符合规范
- [ ] 导入顺序正确
- [ ] 类型注解完整

## 可维护性
- [ ] 代码清晰易懂
- [ ] 函数职责单一
- [ ] 避免代码重复
- [ ] 有适当注释
- [ ] 易于测试

## 性能
- [ ] 使用 useMemo 缓存
- [ ] 使用 useCallback 包装
- [ ] 避免 N+1 查询
- [ ] 列表使用 key
- [ ] 无内存泄漏

## 安全性
- [ ] 参数化查询
- [ ] 输入验证
- [ ] 权限检查
- [ ] 敏感数据保护
- [ ] 错误信息安全

## 测试
- [ ] 有单元测试
- [ ] 覆盖率达标
- [ ] 测试正常场景
- [ ] 测试边界情况
- [ ] 测试错误场景

## 文档
- [ ] 函数有注释
- [ ] 参数有说明
- [ ] 返回值有说明
- [ ] 复杂逻辑有注释
- [ ] API 文档完整
```

## 注意事项

1. **建设性**: 审查意见应具有建设性
2. **具体性**: 问题描述和修复建议具体
3. **优先级**: 明确问题优先级
4. **尊重**: 尊重作者，友好沟通
5. **学习**: 审查过程也是学习过程
