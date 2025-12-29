# Package 开发 Skill

## 技能描述

创建符合 Rush.js Monorepo 规范的完整前端 Package，包括配置文件、目录结构、测试和文档。

## 适用场景

- 需要创建新的前端功能包
- 需要封装共享组件和工具
- 需要创建适配器或接口层
- 需要创建业务模块

## Package 层级规范

### 依赖层级

Coze Studio 前端采用严格的 4 层依赖结构：

**Level 1 - Architecture Layer (@coze-arch)**
- 基础设施和核心工具
- 不依赖任何其他业务包
- 示例：logger, bot-utils, bot-hooks-base

**Level 2 - Common Layer (@coze-studio/common)**
- 共享组件和工具
- 可依赖 Level 1 包
- 示例：biz-components, auth, editor-plugins

**Level 3 - Domain Layer (@coze-studio/{domain})**
- 业务领域模块
- 可依赖 Level 1-2 包
- 示例：agent-ide, workflow, studio

**Level 4 - Application Layer (@coze-studio/apps)**
- 应用入口
- 可依赖所有层级
- 示例：coze-studio, open-chat

### Package 类型

**1. 基础设施包 (Infrastructure)**
- 位置：`packages/arch/`
- 用途：提供底层能力
- 命名：@coze-arch/{name}
- 示例：logger, utils, bot-api

**2. 共享组件包 (Components)**
- 位置：`packages/common/`
- 用途：共享 UI 组件
- 命名：@coze-studio/{name}
- 示例：biz-components, editor-plugins

**3. 业务模块包 (Domain)**
- 位置：`packages/{domain}/`
- 用途：业务功能模块
- 命名：@coze-studio/{domain}
- 示例：agent-ide, workflow

**4. 适配器包 (Adapter)**
- 位置：任意层级
- 用途：解耦和适配
- 命名：{name}-adapter
- 示例：bot-env-adapter, auth-adapter

**5. 接口包 (Interface)**
- 位置：任意层级
- 用途：定义接口契约
- 命名：{name}-interface
- 示例：slardar-interface, uploader-interface

## Package 目录结构

### 标准结构

```
{package-name}/
├── src/                          # 源代码目录
│   ├── index.ts                  # 包入口，导出公共 API
│   ├── types/                    # 类型定义
│   │   └── index.ts
│   ├── components/               # 组件目录（可选）
│   ├── hooks/                    # Hooks 目录（可选）
│   ├── utils/                    # 工具函数目录（可选）
│   ├── logger/                   # 功能模块目录（可选）
│   └── global.d.ts               # 全局类型声明（可选）
│
├── __tests__/                    # 测试目录
│   ├── unit/                     # 单元测试
│   ├── integration/              # 集成测试
│   └── __mocks__/                # Mock 文件
│
├── config/                       # 配置文件目录（可选）
│
├── .gitignore                    # Git 忽略文件
├── eslint.config.js             # ESLint 配置
├── package.json                  # NPM 包配置
├── tsconfig.json                 # TypeScript 基础配置
├── tsconfig.build.json           # TypeScript 构建配置
├── tsconfig.misc.json            # TypeScript 其他配置
├── vitest.config.ts             # Vitest 测试配置
└── README.md                     # 包说明文档
```

### 最小结构

```
{package-name}/
├── src/
│   └── index.ts
├── __tests__/
│   └── index.test.ts
├── .gitignore
├── eslint.config.js
├── package.json
├── tsconfig.json
├── tsconfig.build.json
└── README.md
```

## 配置文件模板

### 1. package.json

```json
{
  "name": "@coze-arch/{package-name}",
  "version": "0.0.1",
  "author": "your-email@example.com",
  "main": "./src/index.ts",
  "scripts": {
    "build": "exit 0",
    "lint": "eslint ./ --cache",
    "test": "vitest --run --passWithNoTests",
    "test:cov": "npm run test -- --coverage"
  },
  "dependencies": {
    "@coze-arch/bot-env": "workspace:*",
    "@coze-arch/bot-typings": "workspace:*",
    "react": "~18.2.0"
  },
  "devDependencies": {
    "@coze-arch/eslint-config": "workspace:*",
    "@coze-arch/ts-config": "workspace:*",
    "@coze-arch/vitest-config": "workspace:*",
    "@types/react": "18.2.37",
    "@types/node": "^18",
    "@vitest/coverage-v8": "~3.0.5",
    "vitest": "~3.0.5"
  }
}
```

**字段说明：**
- `name`: 包名，遵循命名规范
- `main`: 入口文件，必须是 `./src/index.ts`
- `scripts`: 标准脚本命令
- `dependencies`: 使用 `workspace:*` 引用内部包
- `devDependencies`: 共享配置包和工具

### 2. tsconfig.json

```json
{
  "extends": "@coze-arch/ts-config/tsconfig.json",
  "compilerOptions": {
    "composite": true,
    "outDir": "./dist",
    "rootDir": "./src"
  },
  "include": [
    "src/**/*.ts",
    "src/**/*.tsx"
  ],
  "exclude": [
    "node_modules",
    "dist",
    "__tests__"
  ]
}
```

### 3. tsconfig.build.json

```json
{
  "extends": "./tsconfig.json",
  "exclude": [
    "**/*.test.ts",
    "**/*.test.tsx",
    "**/*.spec.ts",
    "**/*.spec.tsx",
    "__tests__",
    "node_modules",
    "dist"
  ]
}
```

### 4. tsconfig.misc.json

```json
{
  "extends": "./tsconfig.json",
  "compilerOptions": {
    "composite": false
  },
  "include": [
    "eslint.config.js",
    "vitest.config.ts",
    "config/**/*.ts"
  ]
}
```

### 5. eslint.config.js

```javascript
module.exports = {
  root: true,
  extends: ['@coze-arch/eslint-config'],
  parserOptions: {
    project: './tsconfig.misc.json',
  },
};
```

### 6. vitest.config.ts

```typescript
import { defineConfig } from 'vitest/config';
import { path as rootPath } from '@coze-arch/ts-config';

export default defineConfig({
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: [],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        '__tests__/',
        '**/*.test.ts',
        '**/*.test.tsx',
        '**/*.spec.ts',
        '**/*.spec.tsx',
        '**/types/**',
      ],
    },
  },
  resolve: {
    alias: {
      '@': rootPath(),
    },
  },
});
```

### 7. .gitignore

```
node_modules/
dist/
*.log
.DS_Store
*.tsbuildinfo
coverage/
.nyc_output/
```

## 源代码组织

### 1. 入口文件 (src/index.ts)

```typescript
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

// 导出类型
export type * from './types';

// 导出工具函数
export * from './utils';

// 导出组件（如果有）
export { ComponentName } from './components';

// 导出 Hooks（如果有）
export * from './hooks';

// 导出主类或实例
export { MainClass, mainInstance } from './main';
```

### 2. 类型定义 (src/types/index.ts)

```typescript
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

/**
 * 选项配置接口
 */
export interface Options {
  /** 配置项1 */
  option1?: string;
  /** 配置项2 */
  option2?: number;
}

/**
 * 数据结构接口
 */
export interface Data {
  id: string;
  name: string;
}

/**
 * 回调函数类型
 */
export type Callback = (data: Data) => void;
```

### 3. 工具函数 (src/utils/index.ts)

```typescript
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

import type { Options, Data } from '../types';

/**
 * 工具函数说明
 *
 * @param param1 - 参数1说明
 * @param param2 - 参数2说明
 * @returns 返回值说明
 *
 * @example
 * ```typescript
 * const result = utilFunction('param1', 123);
 * console.log(result);
 * ```
 */
export function utilFunction(
  param1: string,
  param2: number,
): Data {
  // 实现
  return {
    id: param1,
    name: String(param2),
  };
}

/**
 * 异步工具函数
 */
export async function asyncUtilFunction(
  options: Options,
): Promise<Data> {
  // 实现
  return {
    id: '123',
    name: options.option1 || 'default',
  };
}
```

## 测试规范

### 1. 单元测试结构

```
__tests__/
├── unit/
│   ├── utils.test.ts
│   ├── components.test.tsx
│   └── hooks.test.ts
├── integration/
│   └── features.test.ts
└── __mocks__/
    └── external-module.ts
```

### 2. 单元测试模板

```typescript
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

import { describe, it, expect, vi } from 'vitest';
import { utilFunction } from '../src/utils';

describe('utilFunction', () => {
  it('should return correct data', () => {
    const result = utilFunction('test', 123);
    expect(result).toEqual({
      id: 'test',
      name: '123',
    });
  });

  it('should handle empty input', () => {
    const result = utilFunction('', 0);
    expect(result.id).toBe('');
    expect(result.name).toBe('0');
  });

  it('should throw error on invalid input', () => {
    expect(() => utilFunction(null as any, 123)).toThrow();
  });
});
```

### 3. 组件测试模板

```typescript
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ComponentName } from '../src/components';

describe('ComponentName', () => {
  it('should render correctly', () => {
    render(<ComponentName title="Test" />);
    expect(screen.getByText('Test')).toBeInTheDocument();
  });

  it('should call callback on click', () => {
    const handleClick = vi.fn();
    render(<ComponentName title="Test" onClick={handleClick} />);

    screen.getByRole('button').click();
    expect(handleClick).toHaveBeenCalledTimes(1);
  });
});
```

## README 文档规范

### 完整 README 模板

```markdown
# @coze-arch/package-name

包功能简述（一句话）。

## 功能特性

- 特性1：详细说明
- 特性2：详细说明
- 特性3：详细说明

## 安装

\`\`\`bash
# Rush.js workspace
rush add -p @coze-arch/package-name --dev
rush update
\`\`\`

## 使用方式

### 基本使用

\`\`\`typescript
import { utilFunction } from '@coze-arch/package-name';

const result = utilFunction('param1', 123);
console.log(result);
\`\`\`

### 高级用法

\`\`\`typescript
// 高级用法示例
\`\`\`

## API 文档

### utilFunction

工具函数说明。

**参数：**
- `param1` (string): 参数1说明
- `param2` (number): 参数2说明

**返回：** `Data` - 返回值说明

**示例：**
\`\`\`typescript
const result = utilFunction('test', 123);
\`\`\`

## 开发指南

### 运行测试

\`\`\`bash
# 单元测试
rushx test

# 测试覆盖率
rushx test:cov
\`\`\`

### 代码检查

\`\`\`bash
# Lint
rushx lint

# 自动修复
rushx lint --fix
\`\`\`

### 类型检查

\`\`\`bash
rushx ts-check
\`\`\`

## 依赖说明

### 生产依赖

- **@coze-arch/bot-env**: 环境配置
- **@coze-arch/bot-typings**: 共享类型

### 开发依赖

- **@coze-arch/eslint-config**: ESLint 配置
- **@coze-arch/ts-config**: TypeScript 配置
- **@coze-arch/vitest-config**: Vitest 配置

## 注意事项

1. 注意事项1
2. 注意事项2

## 常见问题

### 问题1

**问题描述**: 问题描述

**解决方案**: 解决方案

## 许可证

Apache License 2.0
```

## Package 类型规范

### 1. 基础设施包

**特点：**
- 纯工具库，无业务逻辑
- 不依赖其他业务包
- 高度可复用

**示例：**
```
@coze-arch/logger         - 日志工具
@coze-arch/utils          - 通用工具
@coze-arch/bot-api        - API 客户端
```

### 2. 组件包

**特点：**
- 包含 UI 组件
- 依赖 React 和 UI 库
- 可依赖 Level 1 包

**示例：**
```
@coze-studio/biz-components  - 业务组件
@coze-studio/editor-plugins  - 编辑器插件
```

### 3. 适配器包

**特点：**
- 用于解耦和适配
- 实现统一接口
- 便于替换实现

**示例：**
```
bot-env-adapter          - 环境适配器
auth-adapter             - 认证适配器
uploader-adapter         - 上传适配器
```

**适配器模板：**
```typescript
// src/index.ts
import type { TargetInterface } from '@coze-arch/target-interface';

/**
 * 适配器实现
 */
export class AdapterImplementation implements TargetInterface {
  async method(params: Params): Result {
    // 适配逻辑
    return result;
  }
}
```

### 4. 接口包

**特点：**
- 只定义接口，不实现
- 被适配器包实现
- 被业务包依赖

**示例：**
```
slardar-interface        - Slardar 接口
uploader-interface       - 上传接口
tea-interface           - TEA 接口
```

**接口包模板：**
```typescript
// src/index.ts
/**
 * 接口定义
 */
export interface TargetInterface {
  method(params: Params): Promise<Result>;
}

/**
 * 接口类型
 */
export type Params = {
  // 类型定义
};

export type Result = {
  // 类型定义
};
```

## 生成流程

### 步骤 1: 确定包类型

1. 确定包的层级（Level 1-4）
2. 确定包的类型（基础设施/组件/适配器/接口）
3. 确定包的命名（@coze-arch 或 @coze-studio）

### 步骤 2: 创建目录结构

1. 在对应层级创建包目录
2. 创建标准目录结构
3. 创建配置文件

### 步骤 3: 配置依赖

1. 配置 package.json
2. 配置 TypeScript
3. 配置 ESLint
4. 配置 Vitest

### 步骤 4: 实现功能

1. 实现类型定义
2. 实现工具函数
3. 实现组件（如果需要）
4. 实现导出

### 步骤 5: 编写测试

1. 编写单元测试
2. 编写集成测试
3. 配置 Mock
4. 运行测试

### 步骤 6: 编写文档

1. 编写 README
2. 编写 JSDoc
3. 提供示例代码

### 步骤 7: 更新 Rush

1. 在 rush.json 中注册包
2. 运行 `rush update`
3. 验证依赖关系

## 输出检查清单

创建 Package 后，确保：

- [ ] 包名符合命名规范
- [ ] 目录结构完整
- [ ] 配置文件完整且正确
- [ ] 依赖关系正确
- [ ] 使用 workspace:* 引用内部包
- [ ] src/index.ts 正确导出
- [ ] 类型定义完整
- [ ] 工具函数实现
- [ ] 单元测试存在
- [ ] 测试覆盖率达标
- [ ] ESLint 配置正确
- [ ] TypeScript 配置正确
- [ ] README 文档完整
- [ ] 包含使用示例
- [ ] rush.json 已更新

## 真实代码示例

### 示例 1: @coze-arch/logger

**目录结构：**
```
logger/
├── src/
│   ├── index.ts
│   ├── logger/
│   ├── reporter/
│   ├── slardar/
│   └── types/
├── __tests__/
├── config/
├── package.json
├── tsconfig.json
└── README.md
```

**package.json：**
```json
{
  "name": "@coze-arch/logger",
  "main": "./src/index.ts",
  "dependencies": {
    "@coze-arch/bot-env": "workspace:*",
    "@coze-arch/bot-typings": "workspace:*",
    "react": "~18.2.0",
    "react-error-boundary": "^4.0.9"
  }
}
```

**导出结构：**
```typescript
// src/index.ts
export * from './logger';
export * from './reporter';
export * from './slardar';
export * from './types';
export { default as ErrorBoundary } from './components/error-boundary';
export * from './hooks';
```

### 示例 2: 适配器包结构

```
bot-env-adapter/
├── src/
│   ├── index.ts
│   ├── env-adapter.ts
│   └── types.ts
├── __tests__/
└── package.json
```

**实现：**
```typescript
// src/env-adapter.ts
import type { EnvInterface } from '@coze-arch/bot-env-interface';

export class EnvAdapter implements EnvInterface {
  getEnv(): string {
    // 适配逻辑
    return process.env.NODE_ENV || 'development';
  }
}
```

## 注意事项

1. **层级依赖**: 严格遵守依赖层级，禁止反向依赖
2. **Workspace**: 内部包依赖必须使用 `workspace:*`
3. **类型导出**: 所有公共类型必须导出
4. **测试覆盖**: 确保 Level 1-2 包有足够测试覆盖率
5. **文档完整**: README 必须包含使用示例
6. **命名规范**: 包名、文件名、变量名都要符合规范
7. **配置共享**: 使用共享的 ESLint、TS、Vitest 配置
8. **构建脚本**: build 脚本对于纯 TS 包可以使用 `exit 0`
