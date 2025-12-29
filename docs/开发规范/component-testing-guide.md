# 组件测试指南

## 概述

本文档提供ZKER前端组件测试的完整指南，包括测试工具、测试策略和示例。

## 1. 测试工具配置

### 1.1 安装依赖

```bash
npm install --save-dev @testing-library/react @testing-library/jest-dom @testing-library/user-event vitest jsdom
```

### 1.2 Vitest配置

```typescript
// vitest.config.ts
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    globals: true,
  },
});
```

### 1.3 测试设置文件

```typescript
// src/test/setup.ts
import '@testing-library/jest-dom';
import { expect, afterEach, vi } from 'vitest';
import { cleanup } from '@testing-library/react';

// 每个测试后清理
afterEach(() => {
  cleanup();
});
```

## 2. 基础UI组件测试

### 2.1 Button组件测试

```typescript
// arch/ui-components/src/components/Button/Button.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Button } from './Button';

describe('Button', () => {
  it('renders children correctly', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick when clicked', () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click me</Button>);

    fireEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('shows loading spinner when loading', () => {
    render(<Button loading>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
    expect(document.querySelector('.spinner')).toBeInTheDocument();
  });

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Click me</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('applies correct variant class', () => {
    const { container: container1 } = render(<Button variant="primary">Primary</Button>);
    const { container: container2 } = render(<Button variant="danger">Danger</Button>);

    expect(container1.querySelector('.button')).toHaveClass({ 'variant-primary': true });
    expect(container2.querySelector('.button')).toHaveClass({ 'variant-danger': true });
  });
});
```

### 2.2 Input组件测试

```typescript
// arch/ui-components/src/components/Input/Input.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Input } from './Input';

describe('Input', () => {
  it('renders input element', () => {
    render(<Input placeholder="Enter text" />);
    expect(screen.getByPlaceholderText('Enter text')).toBeInTheDocument();
  });

  it('calls onChange when value changes', () => {
    const handleChange = vi.fn();
    render(<Input onChange={handleChange} />);

    const input = screen.getByRole('textbox');
    fireEvent.change(input, { target: { value: 'test' } });

    expect(handleChange).toHaveBeenCalled();
  });

  it('shows error state when error prop is true', () => {
    render(<Input error />);
    const input = screen.getByRole('textbox');
    expect(input).toHaveClass({ 'error': true });
  });

  it('displays prefix and suffix icons', () => {
    render(
      <Input
        prefix={<span>$</span>}
        suffix={<span>USD</span>}
      />
    );

    expect(screen.getByText('$')).toBeInTheDocument();
    expect(screen.getByText('USD')).toBeInTheDocument();
  });
});
```

## 3. 业务组件测试

### 3.1 QuotaIndicator测试

```typescript
// arch/business-components/src/components/QuotaIndicator/QuotaIndicator.test.tsx
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { QuotaIndicator } from './QuotaIndicator';

describe('QuotaIndicator', () => {
  it('displays correct usage percentage', () => {
    render(
      <QuotaIndicator
        used={75}
        max={100}
        resourceType="API调用"
        showPercentage
      />
    );

    expect(screen.getByText('75.0%')).toBeInTheDocument();
  });

  it('shows over limit status when used exceeds max', () => {
    const { container } = render(
      <QuotaIndicator
        used={120}
        max={100}
        resourceType="API调用"
      />
    );

    expect(container.querySelector('.overLimit')).toBeInTheDocument();
    expect(screen.getByText('已超限')).toBeInTheDocument();
  });

  it('shows near limit status when usage >= 80%', () => {
    const { container } = render(
      <QuotaIndicator
        used={85}
        max={100}
        resourceType="API调用"
      />
    );

    expect(container.querySelector('.nearLimit')).toBeInTheDocument();
    expect(screen.getByText('即将超限')).toBeInTheDocument();
  });
});
```

### 3.2 PermissionTree测试

```typescript
// arch/business-components/src/components/PermissionTree/PermissionTree.test.tsx
import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { PermissionTree } from './PermissionTree';

const mockData = [
  {
    key: 'bot',
    title: '机器人管理',
    children: [
      { key: 'bot.view', title: '查看' },
      { key: 'bot.create', title: '创建' },
    ],
  },
];

describe('PermissionTree', () => {
  it('renders permission tree', () => {
    render(<PermissionTree data={mockData} />);

    expect(screen.getByText('机器人管理')).toBeInTheDocument();
    expect(screen.getByText('查看')).toBeInTheDocument();
    expect(screen.getByText('创建')).toBeInTheDocument();
  });

  it('calls onChange when checkbox is clicked', () => {
    const handleChange = vi.fn();
    render(
      <PermissionTree
        data={mockData}
        onChange={handleChange}
        checkable
        showCheckbox
      />
    );

    const checkboxes = screen.getAllByRole('checkbox');
    fireEvent.click(checkboxes[0]);

    expect(handleChange).toHaveBeenCalled();
  });

  it('expands and collapses nodes', () => {
    render(<PermissionTree data={mockData} />);

    const expandIcon = screen.getAllByText('▶')[0];
    fireEvent.click(expandIcon);

    // 检查子节点是否显示
    expect(screen.getByText('查看')).toBeVisible();
  });
});
```

## 4. 页面组件测试

### 4.1 TenantList页面测试

```typescript
// studio/pages/tenant/TenantList/TenantList.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { TenantList } from './TenantList';

// Mock API
vi.mock('@coze-studio/api-client', () => ({
  useTenantList: () => ({
    data: {
      tenants: [
        {
          tenant_id: '1',
          tenant_name: 'Test Tenant',
          tenant_type: 'enterprise',
          status: 'active',
        },
      ],
      total_count: 1,
    },
    loading: false,
  }),
}));

describe('TenantList', () => {
  it('renders page title', async () => {
    render(<TenantList />);

    expect(screen.getByText('租户管理')).toBeInTheDocument();
  });

  it('displays tenant data', async () => {
    render(<TenantList />);

    await waitFor(() => {
      expect(screen.getByText('Test Tenant')).toBeInTheDocument();
      expect(screen.getByText('企业')).toBeInTheDocument();
    });
  });

  it('shows create button', () => {
    render(<TenantList />);

    expect(screen.getByText('创建租户')).toBeInTheDocument();
  });
});
```

## 5. 测试覆盖率目标

- **基础UI组件**：≥ 80%
- **业务组件**：≥ 70%
- **页面组件**：≥ 60%

## 6. 运行测试

```bash
# 运行所有测试
npm test

# 运行特定文件的测试
npm test Button.test.tsx

# 生成覆盖率报告
npm run test:cov
```

## 7. 最佳实践

1. **测试用户行为**：测试用户如何使用组件，而不是实现细节
2. **保持简单**：每个测试应该只测试一个功能点
3. **使用描述性名称**：测试名称应该清晰描述测试内容
4. **Mock外部依赖**：Mock API调用、路由等
5. **测试边界情况**：空数据、错误状态、极端值等

## 8. 快速测试清单

- [ ] 组件能正常渲染
- [ ] 用户交互正常工作
- [ ] Props变化时组件正确更新
- [ ] 错误状态正确显示
- [ ] 边界情况得到处理

---

**相关文档**：
- [Vitest官方文档](https://vitest.dev/)
- [Testing Library文档](https://testing-library.com/)
