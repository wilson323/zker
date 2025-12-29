# 创建 Zustand Store Skill

## 技能描述

创建符合 Coze Studio 规范的 Zustand Store，基于项目中的实际 Store 模式（如 ProjectAuthStore）。

## 适用场景

- 需要全局状态管理
- 需要跨组件共享状态
- 需要复杂状态逻辑
- 需要状态持久化

## 工作流程

### 1. 分析 Store 需求

明确以下信息：
- **Store 用途**: 管理什么数据
- **状态结构**: 状态字段有哪些
- **Actions**: 需要哪些操作方法
- **持久化**: 是否需要持久化到 localStorage
- **DevTools**: 是否需要 Redux DevTools 集成

### 2. 定义 Store 接口

```typescript
// {domain}-store.types.ts

/**
 * {Domain}State - 状态接口
 */
interface {Domain}State {
  // 数据状态
  data: DataType;
  items: Item[];

  // UI 状态
  isLoading: boolean;
  isModalOpen: boolean;
  selectedId: string | null;

  // 错误状态
  error: Error | null;
  errorMessage: string;
}

/**
 * {Domain}Actions - 操作接口
 */
interface {Domain}Actions {
  // 数据操作
  fetchData: () => Promise<void>;
  setData: (data: DataType) => void;
  addItem: (item: Item) => void;
  updateItem: (id: string, data: Partial<Item>) => void;
  removeItem: (id: string) => void;

  // UI 操作
  setIsLoading: (loading: boolean) => void;
  openModal: () => void;
  closeModal: () => void;
  setSelectedId: (id: string | null) => void;

  // 错误处理
  setError: (error: Error | null) => void;
  clearError: () => void;
}

/**
 * {Domain}Store - 完整 Store 类型
 */
type {Domain}Store = {Domain}State & {Domain}Actions;
```

### 3. 创建 Store 实现

**文件位置**: `store.ts` 或 `{domain}-store.ts`

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

import { devtools, persist } from 'zustand/middleware';
import { create } from 'zustand';

import type { {Domain}Store } from './{domain}-store.types';

/**
 * {Domain}Store - {domain}数据管理
 *
 * 管理{domain}相关的状态和操作
 */
export const use{Domain}Store = create<{Domain}Store>()(
  devtools(
    (set, get) => ({
      // ==================== 初始状态 ====================

      // 数据状态
      data: null,
      items: [],

      // UI 状态
      isLoading: false,
      isModalOpen: false,
      selectedId: null,

      // 错误状态
      error: null,
      errorMessage: '',

      // ==================== 数据操作 ====================

      /**
       * 获取数据
       */
      fetchData: async () => {
        set({ isLoading: true, error: null });
        try {
          const data = await api.getData();
          set({ data, isLoading: false });
        } catch (error) {
          set({
            error: error as Error,
            errorMessage: error.message,
            isLoading: false,
          });
        }
      },

      /**
       * 设置数据
       */
      setData: (data) => {
        set({ data });
      },

      /**
       * 添加项目
       */
      addItem: (item) => {
        set((state) => ({
          items: [...state.items, item],
        }));
      },

      /**
       * 更新项目
       */
      updateItem: (id, data) => {
        set((state) => ({
          items: state.items.map((item) =>
            item.id === id ? { ...item, ...data } : item
          ),
        }));
      },

      /**
       * 删除项目
       */
      removeItem: (id) => {
        set((state) => ({
          items: state.items.filter((item) => item.id !== id),
        }));
      },

      // ==================== UI 操作 ====================

      /**
       * 设置加载状态
       */
      setIsLoading: (isLoading) => {
        set({ isLoading });
      },

      /**
       * 打开模态框
       */
      openModal: () => {
        set({ isModalOpen: true });
      },

      /**
       * 关闭模态框
       */
      closeModal: () => {
        set({ isModalOpen: false });
      },

      /**
       * 设置选中项 ID
       */
      setSelectedId: (id) => {
        set({ selectedId: id });
      },

      // ==================== 错误处理 ====================

      /**
       * 设置错误
       */
      setError: (error) => {
        set({
          error,
          errorMessage: error?.message || '',
        });
      },

      /**
       * 清除错误
       */
      clearError: () => {
        set({ error: null, errorMessage: '' });
      },
    }),
    {
      enabled: process.env.NODE_ENV === 'development',
      name: '{domain}Store',
    }
  )
);

// ==================== 选择器 Hooks ====================

/**
 * 使用{domain}数据
 */
export const use{Domain}Data = () => use{Domain}Store((state) => state.data);

/**
 * 使用{domain}项目列表
 */
export const use{Domain}Items = () => use{Domain}Store((state) => state.items);

/**
 * 使用{domain}加载状态
 */
export const use{Domain}Loading = () => use{Domain}Store((state) => state.isLoading);

/**
 * 使用{domain}操作
 */
export const use{Domain}Actions = () =>
  use{Domain}Store((state) => ({
    fetchData: state.fetchData,
    setData: state.setData,
    addItem: state.addItem,
    updateItem: state.updateItem,
    removeItem: state.removeItem,
    setIsLoading: state.setIsLoading,
    openModal: state.openModal,
    closeModal: state.closeModal,
    setSelectedId: state.setSelectedId,
    setError: state.setError,
    clearError: state.clearError,
  }));
```

### 4. 持久化 Store（可选）

```typescript
import { devtools, persist } from 'zustand/middleware';

export const use{Domain}Store = create<{Domain}Store>()(
  devtools(
    persist(
      (set, get) => ({
        // Store 实现
      }),
      {
        name: '{domain}-storage', // localStorage key
        partialize: (state) => ({
          // 只持久化部分状态
          data: state.data,
          items: state.items,
          // 不持久化 UI 状态和错误状态
        }),
      }
    ),
    {
      enabled: process.env.NODE_ENV === 'development',
      name: '{domain}Store',
    }
  )
);
```

### 5. Store 测试

**测试文件**: `{domain}-store.test.ts`

```typescript
import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { use{Domain}Store } from './{domain}-store';

describe('{Domain}Store', () => {
  beforeEach(() => {
    // 重置 Store 状态
    use{Domain}Store.setState({
      data: null,
      items: [],
      isLoading: false,
      isModalOpen: false,
      selectedId: null,
      error: null,
      errorMessage: '',
    });
  });

  describe('initial state', () => {
    it('should have correct initial state', () => {
      const { result } = renderHook(() => use{Domain}Store());

      expect(result.current.data).toBeNull();
      expect(result.current.items).toEqual([]);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.isModalOpen).toBe(false);
      expect(result.current.selectedId).toBeNull();
      expect(result.current.error).toBeNull();
    });
  });

  describe('setData', () => {
    it('should set data', () => {
      const { result } = renderHook(() => use{Domain}Store());
      const testData = { id: '1', name: 'Test' };

      act(() => {
        result.current.setData(testData);
      });

      expect(result.current.data).toEqual(testData);
    });
  });

  describe('addItem', () => {
    it('should add item to items array', () => {
      const { result } = renderHook(() => use{Domain}Store());
      const newItem = { id: '1', name: 'Item 1' };

      act(() => {
        result.current.addItem(newItem);
      });

      expect(result.current.items).toHaveLength(1);
      expect(result.current.items[0]).toEqual(newItem);
    });

    it('should add multiple items', () => {
      const { result } = renderHook(() => use{Domain}Store());

      act(() => {
        result.current.addItem({ id: '1', name: 'Item 1' });
        result.current.addItem({ id: '2', name: 'Item 2' });
      });

      expect(result.current.items).toHaveLength(2);
    });
  });

  describe('updateItem', () => {
    it('should update item by id', () => {
      const { result } = renderHook(() => use{Domain}Store());

      act(() => {
        result.current.addItem({ id: '1', name: 'Item 1' });
        result.current.updateItem('1', { name: 'Updated Item' });
      });

      expect(result.current.items[0].name).toBe('Updated Item');
    });

    it('should not update other items', () => {
      const { result } = renderHook(() => use{Domain}Store());

      act(() => {
        result.current.addItem({ id: '1', name: 'Item 1' });
        result.current.addItem({ id: '2', name: 'Item 2' });
        result.current.updateItem('1', { name: 'Updated Item' });
      });

      expect(result.current.items[0].name).toBe('Updated Item');
      expect(result.current.items[1].name).toBe('Item 2');
    });
  });

  describe('removeItem', () => {
    it('should remove item by id', () => {
      const { result } = renderHook(() => use{Domain}Store());

      act(() => {
        result.current.addItem({ id: '1', name: 'Item 1' });
        result.current.addItem({ id: '2', name: 'Item 2' });
        result.current.removeItem('1');
      });

      expect(result.current.items).toHaveLength(1);
      expect(result.current.items[0].id).toBe('2');
    });
  });

  describe('UI operations', () => {
    it('should open modal', () => {
      const { result } = renderHook(() => use{Domain}Store());

      act(() => {
        result.current.openModal();
      });

      expect(result.current.isModalOpen).toBe(true);
    });

    it('should close modal', () => {
      const { result } = renderHook(() => use{Domain}Store());

      act(() => {
        result.current.openModal();
        result.current.closeModal();
      });

      expect(result.current.isModalOpen).toBe(false);
    });

    it('should set selected id', () => {
      const { result } = renderHook(() => use{Domain}Store());

      act(() => {
        result.current.setSelectedId('123');
      });

      expect(result.current.selectedId).toBe('123');
    });
  });

  describe('error handling', () => {
    it('should set error', () => {
      const { result } = renderHook(() => use{Domain}Store());
      const error = new Error('Test error');

      act(() => {
        result.current.setError(error);
      });

      expect(result.current.error).toEqual(error);
      expect(result.current.errorMessage).toBe('Test error');
    });

    it('should clear error', () => {
      const { result } = renderHook(() => use{Domain}Store());

      act(() => {
        result.current.setError(new Error('Test error'));
        result.current.clearError();
      });

      expect(result.current.error).toBeNull();
      expect(result.current.errorMessage).toBe('');
    });
  });

  describe('async operations', () => {
    it('should set loading state during fetch', async () => {
      const { result } = renderHook(() => use{Domain}Store());

      act(() => {
        result.current.fetchData();
      });

      expect(result.current.isLoading).toBe(true);
    });

    it('should handle fetch success', async () => {
      const { result } = renderHook(() => use{Domain}Store());

      await act(async () => {
        await result.current.fetchData();
      });

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it('should handle fetch error', async () => {
      const { result } = renderHook(() => use{Domain}Store>());

      await act(async () => {
        try {
          await result.current.fetchData();
        } catch (error) {
          // 预期会抛出错误
        }
      });

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeTruthy();
    });
  });
});
```

## 必须遵循的规范

### ✅ 命名规范

```typescript
// ✅ 正确的 Store 命名
useProjectAuthStore
useUserStore
useAgentStore
useWorkflowStore

// ❌ 错误的 Store 命名
projectAuthStore         ← 应使用 useProjectAuthStore
use_project_auth_store   ← 应使用驼峰命名
UseProjectStore          ← 不应使用大写开头
```

### ✅ State 字段命名规范

```typescript
// ✅ 正确的状态命名
interface State {
  // 数据状态 - 名词
  users: User[];
  agents: Agent[];

  // UI 状态 - is/has/show 前缀
  isLoading: boolean;
  isModalOpen: boolean;
  hasError: boolean;
  showSidebar: boolean;

  // 选中状态 - selected 前缀
  selectedUserId: number | null;
  selectedItems: string[];

  // 错误状态
  error: Error | null;
  errorMessage: string;
}

// ❌ 错误的状态命名
interface State {
  data: User[];          ← 过于通用，应使用 users
  visible: boolean;      ← 应使用 isModalVisible
  id: number;            ← 应使用 selectedUserId
  form: FormData;        ← 应使用 userForm
}
```

### ✅ Action 命名规范

```typescript
// ✅ 正确的 Action 命名
interface Actions {
  // CRUD 操作
  fetchUsers: () => Promise<void>;
  createUser: (data: CreateUserData) => Promise<void>;
  updateUser: (id: number, data: UpdateUserData) => Promise<void>;
  deleteUser: (id: number) => Promise<void>;

  // 设置操作
  setCurrentUser: (user: User) => void;
  setSelectedUserId: (id: number) => void;

  // UI 操作
  openModal: () => void;
  closeModal: () => void;
  toggleSidebar: () => void;

  // 错误处理
  setError: (error: Error) => void;
  clearError: () => void;
}

// ❌ 错误的 Action 命名
interface Actions {
  users: () => void;             ← 应为 fetchUsers
  add: (user: User) => void;     ← 应为 addUser
  remove: (id: number) => void;  ← 应为 deleteUser
  modal: (open: boolean) => void; ← 应拆分为 openModal/closeModal
}
```

### ✅ 中间件使用顺序

```typescript
// ✅ 正确的中间件顺序
export const useStore = create()(
  devtools(                    // 外层：DevTools
    persist(                  // 中间层：持久化
      (set, get) => ({        // 内层：Store 实现
        // ...
      }),
      {
        name: 'app-storage',
        partialize: (state) => ({ /* ... */ }),
      }
    ),
    {
      enabled: true,
      name: 'appStore',
    }
  )
);

// ❌ 错误的中间件顺序
export const useStore = create()(
  persist(                    ← persist 应在 devtools 内部
    devtools(
      (set, get) => ({ /* ... */ })
    )
  )
);
```

### ✅ 选择器 Hooks 规范

```typescript
// ✅ 正确的选择器命名
export const useUserData = () => useStore((state) => state.data);
export const useUserList = () => useStore((state) => state.users);
export const useUserLoading = () => useStore((state) => state.isLoading);
export const useUserActions = () => useStore((state) => ({
  fetchUsers: state.fetchUsers,
  addUser: state.addUser,
}));

// ❌ 错误的选择器命名
export const useData = () => useStore((state) => state.data);  ← 过于通用
export const useStore = () => useStore();                       ← 与 Store 重名
export const getUsers = () => useStore((state) => state.users); ← 应使用 use 前缀
```

## 输出检查清单

创建完成后，必须确保：

- [ ] Store 文件名使用 kebab-case 或 camelCase
- [ ] Store Hook 使用 use 前缀 + PascalCase
- [ ] State 接口定义完整
- [ ] Actions 接口定义完整
- [ ] 使用 devtools 中间件
- [ ] 使用 persist 中间件（如需持久化）
- [ ] Actions 使用动词-名词模式
- [ ] 状态更新使用 set 函数
- [ ] 提供选择器 Hooks
- [ ] 包含完整的测试
- [ ] 通过 ESLint 检查
- [ ] 通过 TypeScript 编译

## 示例输出

**用户输入**：
```
创建一个用户管理 Store，支持 CRUD 操作和加载状态
```

**AI 输出**：

**文件**: `user-store.ts`

```typescript
import { devtools } from 'zustand/middleware';
import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

interface User {
  id: string;
  name: string;
  email: string;
}

interface UserState {
  users: User[];
  selectedUserId: string | null;
  isLoading: boolean;
  error: Error | null;
}

interface UserActions {
  fetchUsers: () => Promise<void>;
  addUser: (user: Omit<User, 'id'>) => Promise<void>;
  updateUser: (id: string, data: Partial<User>) => void;
  deleteUser: (id: string) => void;
  setSelectedUserId: (id: string | null) => void;
  clearError: () => void;
}

type UserStore = UserState & UserActions;

export const useUserStore = create<UserStore>()(
  devtools(
    (set, get) => ({
      users: [],
      selectedUserId: null,
      isLoading: false,
      error: null,

      fetchUsers: async () => {
        set({ isLoading: true, error: null });
        try {
          const users = await userAPI.getUsers();
          set({ users, isLoading: false });
        } catch (error) {
          set({ error: error as Error, isLoading: false });
        }
      },

      addUser: async (userData) => {
        const user = await userAPI.createUser(userData);
        set((state) => ({
          users: [...state.users, user],
        }));
      },

      updateUser: (id, data) => {
        set((state) => ({
          users: state.users.map((user) =>
            user.id === id ? { ...user, ...data } : user
          ),
        }));
      },

      deleteUser: (id) => {
        set((state) => ({
          users: state.users.filter((user) => user.id !== id),
        }));
      },

      setSelectedUserId: (id) => {
        set({ selectedUserId: id });
      },

      clearError: () => {
        set({ error: null });
      },
    }),
    {
      enabled: process.env.NODE_ENV === 'development',
      name: 'userStore',
    }
  )
);

// 选择器
export const useUsers = () => useUserStore((state) => state.users);
export const useUserLoading = () => useUserStore((state) => state.isLoading);
export const useUserActions = () => useUserStore((state) => ({
  fetchUsers: state.fetchUsers,
  addUser: state.addUser,
  updateUser: state.updateUser,
  deleteUser: state.deleteUser,
}));
```

## 注意事项

1. **状态扁平化**: 避免深层嵌套的状态结构
2. **派生状态**: 使用选择器计算派生状态，而不是存储
3. **性能优化**: 使用选择器避免不必要的重渲染
4. **错误边界**: 在 Store 中处理错误，而不是在组件中
5. **类型安全**: 充分利用 TypeScript 的类型系统
6. **测试覆盖**: 测试所有 Actions 和状态变化
