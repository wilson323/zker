# ZKER 前端用户手册 v1.0

> **适用对象**: 研发 C - 前端工程师、前端开发团队
> **最后更新**: 2025-01-01
> **版本**: v1.0

---

## 📖 目录

1. [快速开始](#快速开始)
2. [组件库使用指南](#组件库使用指南)
3. [业务页面开发](#业务页面开发)
4. [国际化使用](#国际化使用)
5. [性能优化指南](#性能优化指南)
6. [常见问题](#常见问题)

---

## 快速开始

### 环境准备

```bash
# 1. 安装依赖
rush update

# 2. 启动开发服务器
cd frontend/apps/coze-studio
npm run dev

# 3. 访问应用
open http://localhost:8888
```

### 项目结构

```
frontend/packages/
├── arch/
│   ├── ui-components/        # 基础UI组件库
│   └── business-components/  # 业务组件库
├── common/                   # 共享组件和工具
│   ├── i18n/                # 国际化
│   ├── themes/              # 主题系统
│   └── components/          # 通用组件
└── studio/                  # Studio应用
    └── pages/               # 页面组件
```

---

## 组件库使用指南

### 基础UI组件

#### Button 按钮

```typescript
import { Button } from '@coze-studio/ui-components';

// 基础用法
<Button onClick={handleClick}>点击我</Button>

// 不同变体
<Button variant="primary">主要按钮</Button>
<Button variant="danger">删除</Button>
<Button variant="outline">取消</Button>

// 不同尺寸
<Button size="sm">小按钮</Button>
<Button size="md">中按钮</Button>
<Button size="lg">大按钮</Button>

// 加载状态
<Button loading loading>提交中...</Button>

// 带图标
<Button icon="🚀">启动</Button>

// 禁用状态
<Button disabled>不可点击</Button>
```

#### Input 输入框

```typescript
import { Input } from '@coze-studio/ui-components';

// 基础用法
<Input placeholder="请输入内容" onChange={handleChange} />

// 错误状态
<Input error placeholder="错误的输入框" />

// 带前后缀
<Input prefix="$" suffix="USD" placeholder="金额" />

// 不同尺寸
<Input size="sm" />
<Input size="md" />
<Input size="lg" />
```

#### Select 下拉选择器

```typescript
import { Select } from '@coze-studio/ui-components';

const options = [
  { label: '选项1', value: '1' },
  { label: '选项2', value: '2' },
];

<Select
  value={value}
  onChange={setValue}
  options={options}
  placeholder="请选择"
/>
```

#### Form 表单

```typescript
import { Form, FormItem, Input, Button } from '@coze-studio/ui-components';

const handleSubmit = (values) => {
  console.log(values);
};

<Form initialValues={{ name: '' }} onSubmit={handleSubmit}>
  <FormItem name="name" label="姓名" required>
    <Input placeholder="请输入姓名" />
  </FormItem>

  <Button type="submit">提交</Button>
</Form>
```

### 业务组件

#### TenantSelector 租户选择器

```typescript
import { TenantSelector } from '@coze-studio/business-components';

<TenantSelector
  value={tenantId}
  onChange={setTenantId}
  placeholder="选择租户"
  allowClear
/>
```

#### QuotaIndicator 配额指示器

```typescript
import { QuotaIndicator } from '@coze-studio/business-components';

<QuotaIndicator
  used={750}
  max={1000}
  resourceType="API调用"
  showPercentage
/>
```

#### PermissionTree 权限树

```typescript
import { PermissionTree } from '@coze-studio/business-components';

const permissions = [
  {
    key: 'bot',
    title: '机器人管理',
    children: [
      { key: 'bot.view', title: '查看' },
      { key: 'bot.create', title: '创建' },
    ],
  },
];

<PermissionTree
  data={permissions}
  onChange={handleChange}
  checkable
  showCheckbox
/>
```

#### SubscriptionSelector 订阅选择器

```typescript
import { SubscriptionSelector } from '@coze-studio/business-components';

<SubscriptionSelector
  value="pro"
  onChange={setSubscription}
/>
```

#### DataScopeSelector 数据权限选择器

```typescript
import { DataScopeSelector } from '@coze-studio/business-components';

<DataScopeSelector
  value="DEPARTMENT"
  onChange={setDataScope}
  resourceType="机器人"
/>
```

---

## 业务页面开发

### 租户管理页面

```typescript
// 使用TenantList页面
import { TenantList } from '@coze-studio/studio';

// 路由配置
<Route path="/tenants" element={<TenantList />} />
<Route path="/tenants/:tenantId" element={<TenantDetail />} />
```

### 权限管理页面

```typescript
import { RoleList } from '@coze-studio/studio';

// 路由配置
<Route path="/permissions/roles" element={<RoleList />} />
<Route path="/permissions/roles/:roleId" element={<RoleDetail />} />
```

### 路由配置页面

```typescript
import { RoutingRules, IntentMatcher } from '@coze-studio/studio';

// 路由配置
<Route path="/routing/rules" element={<RoutingRules />} />
<Route path="/routing/intents" element={<IntentMatcher />} />
```

---

## 国际化使用

### 在组件中使用i18n

```typescript
import { useTranslation } from 'react-i18next';

export const MyComponent = () => {
  const { t } = useTranslation();

  return (
    <div>
      <h1>{t('tenant.title')}</h1>
      <Button>{t('common.confirm')}</Button>
    </div>
  );
};
```

### 切换语言

```typescript
import { LanguageSwitcher } from '@coze-studio/common/components';

<LanguageSwitcher showLabel />
```

---

## 性能优化指南

### 代码分割

```typescript
import { lazy, Suspense } from 'react';

const TenantList = lazy(() => import('./pages/tenant/TenantList'));

<Suspense fallback={<LoadingWrapper loading />}>
  <TenantList />
</Suspense>
```

### React性能优化

```typescript
import { memo, useMemo, useCallback } from 'react';

// 使用React.memo
export const MyComponent = memo(({ data }) => {
  // 使用useMemo缓存计算结果
  const sortedData = useMemo(() => {
    return data.sort((a, b) => a.createdAt - b.createdAt);
  }, [data]);

  // 使用useCallback稳定函数引用
  const handleClick = useCallback(() => {
    console.log('clicked');
  }, []);

  return <div>{/* ... */}</div>;
});
```

---

## 常见问题

### Q: 如何添加新的UI组件？

A: 在 `arch/ui-components/src/components/` 下创建新组件文件夹，包含：
- `ComponentName.tsx` - 组件实现
- `ComponentName.styles.ts` - 样式定义
- `index.ts` - 导出文件
- `ComponentName.test.tsx` - 测试文件
- `ComponentName.stories.tsx` - Storybook故事

### Q: 如何添加新的业务页面？

A: 在 `studio/src/pages/` 下创建新页面文件夹，参考现有页面的结构：
- 页面主组件
- 子组件（如需要）
- 样式文件
- 索引文件

### Q: 如何添加新的翻译？

A: 编辑以下文件：
- `frontend/packages/common/i18n/locales/zh-CN.json`
- `frontend/packages/common/i18n/locales/en-US.json`

遵循命名规范：`模块.子模块.具体项`

### Q: 如何使用主题Token？

A: 使用 `useTheme` Hook：

```typescript
import { useTheme } from '@coze-studio/common/themes';

export const MyComponent = () => {
  const theme = useTheme();

  return (
    <div style={{ color: theme.colors.primary[500] }}>
      Hello
    </div>
  );
};
```

---

## 📚 更多资源

- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)
- [全局一致性检查清单](./ZKER-全局一致性检查清单_v1.0.md)
- [组件测试指南](../开发规范/component-testing-guide.md)
- [性能优化指南](../开发规范/performance-optimization-guide.md)

---

**文档维护**: 研发 C - 前端工程师
**反馈渠道**: 技术负责人
