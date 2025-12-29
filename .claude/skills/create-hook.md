# 创建自定义 Hook Skill

## 技能描述

创建符合 Coze Studio 规范的自定义 React Hook，基于项目中的实际 Hook 模式（如 useInitialValue、useResponsive）。

## 适用场景

- 需要复用状态逻辑
- 需要封装副作用逻辑
- 需要抽象数据获取逻辑
- 需要封装响应式逻辑

## 工作流程

### 1. 分析 Hook 需求

明确以下信息：
- **Hook 功能**: 主要功能描述
- **输入参数**: Hook 接收的参数
- **返回值**: Hook 返回的数据和方法
- **依赖逻辑**: 副作用、订阅等
- **复用性**: 是否需要在多处复用

### 2. 确定 Hook 类型

根据功能选择 Hook 类型：

**类型 1: 值保持 Hook**（如 useInitialValue）
```typescript
// 保持初始值不变
export function useInitialValue<T>(value: T): T {
  const ref = useRef<T>(value);
  return ref.current;
}
```

**类型 2: 响应式 Hook**（如 useResponsive）
```typescript
// 响应式断点逻辑
export const useIsResponsiveByRouteConfig = () => {
  const { responsive } = useRouteConfig();
  const matches = useMediaQuery(/* ... */);
  return shouldResponsive && isResponsive;
};
```

**类型 3: 数据获取 Hook**
```typescript
// 数据获取和状态管理
export function useUserData(userId: number) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    // 获取逻辑
  }, [userId]);

  return { user, loading, error, refetch };
}
```

**类型 4: 事件处理 Hook**
```typescript
// 事件处理和回调管理
export function useHandleSubmit(onSubmit: (data: FormData) => Promise<void>) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = useCallback(async (data: FormData) => {
    setIsSubmitting(true);
    setError(null);
    try {
      await onSubmit(data);
    } catch (err) {
      setError(err);
    } finally {
      setIsSubmitting(false);
    }
  }, [onSubmit]);

  return { handleSubmit, isSubmitting, error };
}
```

### 3. 创建 Hook 文件

**文件位置**: `hooks/use-{purpose}.ts`

```typescript
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

/**
 * use{HookName} - Hook 功能描述
 *
 * @description 详细描述 Hook 的功能和使用场景
 *
 * @template T - 泛型参数说明（如果适用）
 * @param param1 - 参数1说明
 * @param param2 - 参数2说明
 * @returns 返回值说明
 *
 * @example
 * ```typescript
 * const result = use{HookName}('param');
 * ```
 */
export function use{HookName}<T>(
  value: T,
  options?: { option1?: boolean }
): {
  value: T;
  setValue: (value: T) => void;
  reset: () => void;
} {
  // 1. 使用 useRef 保持可变值
  const ref = useRef<T>(value);

  // 2. 使用 useState 管理状态
  const [state, setState] = useState<T>(value);

  // 3. 使用 useMemo 缓存计算值
  const computedValue = useMemo(() => {
    return processValue(state);
  }, [state]);

  // 4. 使用 useCallback 包装回调
  const setValue = useCallback((newValue: T) => {
    setState(newValue);
  }, []);

  const reset = useCallback(() => {
    setState(value);
  }, [value]);

  // 5. 使用 useEffect 处理副作用
  useEffect(() => {
    // 副作用逻辑
    ref.current = state;

    // 清理函数
    return () => {
      // 清理逻辑
    };
  }, [state]);

  // 6. 返回值（对象或元组）
  return {
    value: computedValue,
    setValue,
    reset,
  };
}
```

### 4. Hook 类型定义模式

#### 模式 1: 简单 Hook
```typescript
export function useInitialValue<T>(value: T): T {
  const ref = useRef<T>(value);
  return ref.current;
}
```

#### 模式 2: 状态 + 操作返回
```typescript
export function useToggle(initialValue = false): [boolean, () => void] {
  const [value, setValue] = useState(initialValue);
  const toggle = useCallback(() => setValue(v => !v), []);
  return [value, toggle];
}
```

#### 模式 3: 对象返回
```typescript
export function useLocalStorage<T>(key: string, initialValue: T) {
  const [storedValue, setStoredValue] = useState<T>(() => {
    // 初始化逻辑
  });

  const setValue = useCallback((value: T) => {
    // 设置逻辑
  }, [key]);

  const removeValue = useCallback(() => {
    // 删除逻辑
  }, [key]);

  return {
    value: storedValue,
    setValue,
    removeValue,
  };
}
```

### 5. 创建 Hook 测试

**测试文件**: `hooks/use-{hook-name}.test.ts`

```typescript
import { renderHook, act, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { use{HookName} } from './use-{hook-name}';

describe('use{HookName}', () => {
  // 1. 初始化测试
  describe('initialization', () => {
    it('should initialize with default value', () => {
      const { result } = renderHook(() => use{HookName}());
      expect(result.current.value).toBeDefined();
    });

    it('should initialize with provided value', () => {
      const { result } = renderHook(() => use{HookName}('test'));
      expect(result.current.value).toBe('test');
    });
  });

  // 2. 状态更新测试
  describe('state updates', () => {
    it('should update value when setValue is called', () => {
      const { result } = renderHook(() => use{HookName}('initial'));

      act(() => {
        result.current.setValue('updated');
      });

      expect(result.current.value).toBe('updated');
    });

    it('should reset to initial value when reset is called', () => {
      const { result } = renderHook(() => use{HookName}('initial'));

      act(() => {
        result.current.setValue('updated');
      });

      act(() => {
        result.current.reset();
      });

      expect(result.current.value).toBe('initial');
    });
  });

  // 3. 副作用测试
  describe('side effects', () => {
    it('should call callback on mount', () => {
      const callback = vi.fn();
      renderHook(() => use{HookName}('test', { onMount: callback }));
      expect(callback).toHaveBeenCalledTimes(1);
    });

    it('should call callback on unmount', () => {
      const callback = vi.fn();
      const { unmount } = renderHook(() =>
        use{HookName}('test', { onUnmount: callback })
      );

      unmount();
      expect(callback).toHaveBeenCalledTimes(1);
    });

    it('should call callback when value changes', () => {
      const callback = vi.fn();
      const { result } = renderHook(() =>
        use{HookName}('test', { onChange: callback })
      );

      act(() => {
        result.current.setValue('updated');
      });

      expect(callback).toHaveBeenCalledWith('updated');
    });
  });

  // 4. 异步操作测试
  describe('async operations', () => {
    it('should handle async operations correctly', async () => {
      const { result } = renderHook(() => use{HookName}());

      let resolvedValue;
      await act(async () => {
        resolvedValue = await result.current.asyncOperation();
      });

      expect(resolvedValue).toBeDefined();
    });
  });

  // 5. 内存泄漏测试
  describe('cleanup', () => {
    it('should clean up subscriptions on unmount', () => {
      const unsubscribe = vi.fn();
      const { unmount } = renderHook(() =>
        use{HookName}({ subscribe: () => unsubscribe })
      );

      unmount();
      expect(unsubscribe).toHaveBeenCalledTimes(1);
    });
  });
});
```

### 6. Hook 目录结构

```
hooks/
├── use-{hook-name}.ts          # Hook 实现
└── __tests__/
    └── use-{hook-name}.test.ts # Hook 测试
```

或在组件内部：

```
components/
└── {component-name}/
    └── hooks/
        ├── use-{hook-1}.ts
        ├── use-{hook-2}.ts
        └── __tests__/
            ├── use-{hook-1}.test.ts
            └── use-{hook-2}.test.ts
```

## 必须遵循的规范

### ✅ 命名规范

```typescript
// ✅ 正确
useInitialValue
useResponsive
useAuthForm
useMediaQuery
useLocalStorage
useDebounce

// ❌ 错误
useInitialvalue          ← 应使用 useInitialValue
useauth-form            ← 应使用 useAuthForm
getUseData             ← 不应使用 get 前缀
use_data               ← 应使用驼峰命名
```

### ✅ 文件命名规范

```bash
# ✅ 正确
use-initial-value.ts
use-auth-form.ts
use-responsive.ts

# ❌ 错误
useInitialValue.ts       ← 应使用 kebab-case
use_auth_form.ts         ← 不应使用下划线
useAuthForm.ts           ← 应使用 kebab-case
```

### ✅ Hook 规则

1. **只在顶层调用 Hook**
```typescript
// ✅ 正确
export function MyComponent() {
  const value = useMyHook();
  return <div>{value}</div>;
}

// ❌ 错误
export function MyComponent() {
  if (condition) {
    const value = useMyHook();  ← 不应在条件语句中使用
  }
}
```

2. **只在 React 函数中调用 Hook**
```typescript
// ✅ 正确
export function useMyHook() {
  const [state, setState] = useState();
  useEffect(() => {}, []);
}

// ❌ 错误
export function normalFunction() {
  const [state, setState] = useState();  ← 不应在普通函数中使用
}
```

3. **使用 useCallback 包装回调函数**
```typescript
// ✅ 正确
export function useMyHook() {
  const handleClick = useCallback(() => {
    // 处理逻辑
  }, []);

  return { handleClick };
}

// ❌ 错误
export function useMyHook() {
  const handleClick = () => {  ← 未使用 useCallback，每次渲染都会创建新函数
    // 处理逻辑
  };

  return { handleClick };
}
```

4. **使用 useMemo 缓存计算值**
```typescript
// ✅ 正确
export function useMyHook(items) {
  const sortedItems = useMemo(() => {
    return items.sort();
  }, [items]);

  return { sortedItems };
}

// ❌ 错误
export function useMyHook(items) {
  const sortedItems = items.sort();  ← 每次渲染都会重新计算
  return { sortedItems };
}
```

### ✅ 返回值规范

**元组返回**（类似 useState）:
```typescript
export function useToggle(initialValue = false): [boolean, () => void] {
  const [value, setValue] = useState(initialValue);
  const toggle = useCallback(() => setValue(v => !v), []);
  return [value, toggle];
}

// 使用
const [isOpen, toggle] = useToggle();
```

**对象返回**:
```typescript
export function useLocalStorage(key: string) {
  const [value, setValue] = useState();

  return {
    value,
    setValue,
    removeValue: () => {},
  };
}

// 使用
const { value, setValue, removeValue } = useLocalStorage('key');
```

### ✅ TypeScript 类型规范

```typescript
// ✅ 正确: 使用泛型
export function useInitialValue<T>(value: T): T {
  const ref = useRef<T>(value);
  return ref.current;
}

// ✅ 正确: 明确的返回类型
export function useToggle(): [boolean, () => void] {
  // 实现
}

// ✅ 正确: 对象返回类型
export function useAuth() {
  return {
    user: User | null,
    login: (credentials: Credentials) => Promise<void>,
    logout: () => void,
  };
}

// ❌ 错误: 缺少类型
export function useInitialValue(value) {  ← 应添加类型注解
  return value;
}
```

## 输出检查清单

创建完成后，必须确保：

- [ ] Hook 文件名使用 kebab-case
- [ ] Hook 函数名使用 use 前缀 + camelCase
- [ ] 只在顶层调用其他 Hooks
- [ ] 使用 useCallback 包装回调函数
- [ ] 使用 useMemo 缓存计算值
- [ ] 使用 useEffect 处理副作用（包含清理函数）
- [ ] 返回值类型明确（元组或对象）
- [ ] 泛型参数有意义（使用 T、TData 等）
- [ ] 包含完整的 JSDoc 注释
- [ ] 包含测试文件
- [ ] 测试覆盖所有场景
- [ ] 通过 ESLint 检查
- [ ] 通过 TypeScript 编译

## 示例输出

**用户输入**：
```
创建一个防抖 Hook，用于延迟执行函数
```

**AI 输出**：

1. **文件**: `hooks/use-debounce.ts`

```typescript
import { useEffect, useRef } from 'react';

/**
 * useDebounce - 防抖 Hook
 *
 * 延迟执行函数，在指定时间内多次调用只执行最后一次
 *
 * @template T - 函数类型
 * @param func - 需要防抖的函数
 * @param delay - 延迟时间（毫秒）
 * @returns 防抖后的函数
 *
 * @example
 * ```typescript
 * const debouncedSearch = useDebounce(searchFunction, 300);
 * debouncedSearch('query');
 * ```
 */
export function useDebounce<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (...args: Parameters<T>) => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    timeoutRef.current = setTimeout(() => {
      func(...args);
    }, delay);
  };
}
```

2. **测试文件**: `hooks/use-debounce.test.ts`

```typescript
import { renderHook, act, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useDebounce } from './use-debounce';

describe('useDebounce', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should delay function execution', () => {
    const callback = vi.fn();
    const { result } = renderHook(() => useDebounce(callback, 300));

    act(() => {
      result.current('test');
    });

    expect(callback).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(300);
    });

    expect(callback).toHaveBeenCalledWith('test');
  });

  it('should reset timer on multiple calls', () => {
    const callback = vi.fn();
    const { result } = renderHook(() => useDebounce(callback, 300));

    act(() => {
      result.current('test1');
      result.current('test2');
      result.current('test3');
    });

    act(() => {
      vi.advanceTimersByTime(300);
    });

    expect(callback).toHaveBeenCalledTimes(1);
    expect(callback).toHaveBeenCalledWith('test3');
  });

  it('should clean up timeout on unmount', () => {
    const callback = vi.fn();
    const { unmount } = renderHook(() => useDebounce(callback, 300));

    act(() => {
      const timeout = result.current('test');
      unmount();
    });

    // Timeout should be cleared
    act(() => {
      vi.advanceTimersByTime(300);
    });

    expect(callback).not.toHaveBeenCalled();
  });
});
```

## 注意事项

1. **单一职责**: 每个 Hook 只做一件事
2. **可复用性**: Hook 应该可以在多处复用
3. **性能优化**: 使用 useCallback、useMemo 避免不必要的重渲染
4. **内存管理**: 确保 useEffect 的清理函数正确实现
5. **类型安全**: 充分利用 TypeScript 的类型系统
6. **测试完整**: 测试应覆盖所有使用场景
