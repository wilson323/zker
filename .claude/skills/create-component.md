# 创建 React 组件 Skill

## 技能描述

创建符合 Coze Studio 开发规范的 React 组件，包括组件文件、类型定义、样式文件和测试文件。

## 适用场景

- 需要创建新的 UI 组件
- 需要重构现有组件
- 需要添加复杂交互逻辑的组件

## 工作流程

### 1. 收集需求信息

在开始之前，必须明确：
- **组件名称**：PascalCase（如 UserProfile）
- **组件功能**：描述组件的主要功能
- **Props 接口**：组件需要接收哪些属性
- **状态管理**：是否需要内部状态
- **事件处理**：需要处理哪些用户交互
- **依赖资源**：是否需要调用 API 或使用 Store

### 2. 确定组件结构

根据复杂度选择组件结构：

**简单组件**（单文件）：
```
component-name.tsx
```

**复杂组件**（多文件）：
```
component-name/
├── index.ts                    # 导出入口
├── component-name.tsx          # 主组件
├── component-name.types.ts     # 类型定义
├── component-name.styles.ts    # 样式
└── component-name.test.tsx     # 测试
```

### 3. 创建组件文件

#### 3.1 类型定义文件（`component-name.types.ts`）

```typescript
// ✅ 必须遵循的命名规范
export interface ComponentNameProps {
  // 必需属性
  id: string | number;

  // 可选属性（使用 ? 标记）
  title?: string;
  description?: string;

  // 回调函数（on 前缀）
  onClick?: () => void;
  onSubmit?: (data: FormData) => void;

  // 布尔属性（is/has/show 前缀）
  isLoading?: boolean;
  hasPermission?: boolean;
  showActions?: boolean;

  // 渲染函数
  renderHeader?: (data: Data) => React.ReactNode;

  // children
  children?: React.ReactNode;
}

// ✅ 导出默认值
export const defaultProps: Partial<ComponentNameProps> = {
  isLoading: false,
  showActions: true,
};
```

#### 3.2 主组件文件（`component-name.tsx`）

```typescript
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import type { ComponentNameProps } from './component-name.types';

/**
 * ComponentName 组件
 * @description 组件的详细描述
 */
export function ComponentName({ id, onClick, isLoading }: ComponentNameProps) {
  // 1. Hooks 调用（必须按顺序）
  const [state, setState] = useState(null);
  const [loading, setLoading] = useState(false);

  // 2. 派生状态（useMemo）
  const computedValue = useMemo(() => {
    return computeValue(state);
  }, [state]);

  // 3. 副作用（useEffect）
  useEffect(() => {
    // 副作用逻辑
    return () => {
      // 清理函数
    };
  }, [id]);

  // 4. 事件处理函数（useCallback）
  const handleClick = useCallback(() => {
    onClick?.();
  }, [onClick]);

  const handleSubmit = useCallback((data: FormData) => {
    // 提交逻辑
  }, []);

  // 5. 渲染辅助函数（纯函数）
  function renderHeader() {
    if (!hasHeader) return null;
    return <header>{/* Header 内容 */}</header>;
  }

  // 6. 条件渲染（早期返回）
  if (loading) return <Spinner />;
  if (error) return <ErrorMessage />;

  // 7. 返回 JSX
  return (
    <div className="component-name">
      {renderHeader()}
      {/* 组件内容 */}
    </div>
  );
}

// ✅ 设置默认 Props
ComponentName.defaultProps = defaultProps;
```

#### 3.3 导出文件（`index.ts`）

```typescript
export { ComponentName } from './component-name';
export type { ComponentNameProps } from './component-name.types';
```

### 4. 样式规范

**CSS Modules**（推荐用于复杂组件）：
```typescript
import styles from './component-name.module.css';

export function ComponentName() {
  return <div className={styles.container}>...</div>;
}
```

**Tailwind CSS**（推荐用于简单组件）：
```typescript
import { cn } from '@coze-studio/utils';

export function ComponentName({ className }) {
  return <div className={cn('flex items-center', className)}>...</div>;
}
```

### 5. 测试文件

```typescript
import { render, screen } from '@testing-library/react';
import { ComponentName } from './component-name';

describe('ComponentName', () => {
  describe('when loading', () => {
    it('should show spinner', () => {
      render(<ComponentName isLoading />);
      expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });
  });

  describe('when data is loaded', () => {
    it('should display content', () => {
      render(<ComponentName id="123" />);
      expect(screen.getByText(/content/i)).toBeInTheDocument();
    });
  });
});
```

## 必须遵循的规范

### ✅ 命名规范
- **组件名**: PascalCase（UserProfile）
- **文件名**: kebab-case（user-profile.tsx）
- **Props 接口**: ComponentNameProps
- **事件处理**: handle 前缀（handleClick）
- **回调 Props**: on 前缀（onClick）

### ✅ 组件结构规范
```typescript
// 严格遵循以下顺序
export function Component(props) {
  // 1. Hooks 调用
  const [state, setState] = useState();

  // 2. 派生状态（useMemo）
  const value = useMemo(() => {}, []);

  // 3. 副作用（useEffect）
  useEffect(() => {}, []);

  // 4. 事件处理函数（useCallback）
  const handler = useCallback(() => {}, []);

  // 5. 渲染辅助函数
  function helper() {}

  // 6. 条件渲染（早期返回）
  if (condition) return <Component />;

  // 7. 返回 JSX
  return <div>...</div>;
}
```

### ✅ Props 定义规范
- 必需属性放在前面
- 可选属性使用 `?` 标记
- 回调函数使用 `on` 前缀
- 布尔值使用 `is/has/show` 前缀
- 提供 defaultProps

### ✅ 性能优化规范
- 使用 `useCallback` 包装事件处理函数
- 使用 `useMemo` 缓存计算结果
- 避免在渲染中创建新对象/数组
- 合理使用 `React.memo`

### ✅ TypeScript 规范
- 所有 Props 必须定义类型
- 使用严格的类型检查
- 避免使用 `any` 类型
- 泛型使用描述性名称

## 输出检查清单

创建完成后，必须确保：

- [ ] 组件名使用 PascalCase
- [ ] 文件名使用 kebab-case
- [ ] Props 接口定义完整
- [ ] 遵循组件结构规范
- [ ] 事件处理函数使用 useCallback
- [ ] 计算值使用 useMemo
- [ ] 包含 PropTypes 或 TypeScript 类型
- [ ] 包含测试文件
- [ ] 通过 ESLint 检查
- [ ] 通过 TypeScript 编译
- [ ] 添加必要的注释

## 示例输出

**用户输入**：
```
创建一个用户头像组件，支持点击事件，可以显示加载状态
```

**AI 输出**：

1. **文件结构**：
```
user-avatar/
├── index.ts
├── user-avatar.tsx
├── user-avatar.types.ts
└── user-avatar.test.tsx
```

2. **类型定义**：
```typescript
export interface UserAvatarProps {
  src: string;
  alt?: string;
  size?: 'sm' | 'md' | 'lg';
  onClick?: () => void;
  isLoading?: boolean;
}
```

3. **组件实现**：
```typescript
export function UserAvatar({ src, alt, size = 'md', onClick, isLoading }: UserAvatarProps) {
  const handleClick = useCallback(() => {
    onClick?.();
  }, [onClick]);

  if (isLoading) {
    return <AvatarSkeleton size={size} />;
  }

  return (
    <img
      src={src}
      alt={alt}
      className={cn('avatar', `avatar-${size}`)}
      onClick={handleClick}
    />
  );
}
```

## 注意事项

1. **避免过度拆分**：简单组件使用单文件即可
2. **保持纯粹**：组件不应包含业务逻辑
3. **可复用性**：组件应设计为可复用的
4. **性能优先**：合理使用 memo、useCallback、useMemo
5. **可测试性**：组件应易于测试
