# 国际化(i18n)迁移指南

## 概述

本文档指导如何将现有组件改造为支持多语言。

## 1. 在组件中使用useTranslation Hook

### 示例：改造Button组件

```tsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import { useStyles } from './Button.styles';

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  loading = false,
  children,
  ...rest
}) => {
  // 添加useTranslation Hook
  const { t } = useTranslation();

  const classes = useStyles({ variant, size });

  // 使用t函数获取翻译文本
  const loadingText = loading ? t('common.loading') : '';

  return (
    <button className={classes.button} {...rest}>
      {loading && <span className={classes.spinner} />}
      {children}
    </button>
  );
};
```

### 示例：改造TenantList页面

```tsx
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Table, Button, Modal } from '@coze-studio/ui-components';

export const TenantList: React.FC = () => {
  const { t } = useTranslation();

  const columns = [
    {
      title: t('tenant.name'), // 使用翻译key
      dataIndex: 'tenant_name',
      key: 'tenant_name',
    },
    {
      title: t('tenant.type'),
      dataIndex: 'tenant_type',
      key: 'tenant_type',
      render: (type: string) => t(`tenant.types.${type}`), // 动态翻译key
    },
  ];

  return (
    <div>
      <h1>{t('tenant.title')}</h1>
      <Button>{t('tenant.create')}</Button>
    </div>
  );
};
```

## 2. 在配置文件中初始化i18n

在应用入口文件（如App.tsx或main.tsx）中：

```tsx
import i18n from '@coze-studio/common/i18n';
import { I18nextProvider } from 'react-i18next';
import App from './App';

const Root = () => (
  <I18nextProvider i18n={i18n}>
    <App />
  </I18nextProvider>
);

export default Root;
```

## 3. 翻译Key命名规范

### 规范

- 使用点号(.)分隔命名空间
- 结构：`模块.子模块.具体项`
- 复用通用翻译：`common.*`

### 示例

```json
{
  "common": {
    "confirm": "确认",
    "cancel": "取消",
    "save": "保存"
  },
  "tenant": {
    "title": "租户管理",
    "name": "租户名称",
    "types": {
      "individual": "个人",
      "enterprise": "企业"
    }
  }
}
```

## 4. 动态内容翻译

使用插值语法：

```tsx
// 组件中
const message = t('tenant.deleteConfirm', { name: tenantName });

// 翻译文件中
"deleteConfirm": "确定要删除租户\"{{name}}\"吗？"
```

## 5. 日期和数字格式化

```tsx
import { useTranslation } from 'react-i18next';

export const MyComponent = () => {
  const { t, i18n } = useTranslation();

  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat(i18n.language).format(date);
  };

  const formatNumber = (num: number) => {
    return new Intl.NumberFormat(i18n.language).format(num);
  };
};
```

## 6. 常见模式

### 条件翻译

```tsx
const status = t(`tenant.statuses.${status}`); // 动态key
```

### 带默认值的翻译

```tsx
const text = t('custom.key', { defaultValue: '默认文本' });
```

### 批量翻译

```tsx
const texts = {
  title: t('module.title'),
  description: t('module.description'),
  submit: t('common.submit'),
};
```

## 7. 最佳实践

1. **尽早国际化**：从项目开始就使用i18n，避免后期大量重构
2. **保持同步**：添加新功能时同步更新中英文翻译
3. **使用命名空间**：避免翻译key冲突
4. **提取公共文本**：将"确认"、"取消"等放在common下
5. **测试多语言**：切换语言测试所有页面

## 8. 已完成的翻译内容

✅ 通用文本 (common.*)
✅ 租户管理 (tenant.*)
✅ 权限管理 (permission.*)
✅ 配额管理 (quota.*)
✅ 按钮组件 (button.*)
✅ 表格组件 (table.*)
✅ 模态框 (modal.*)
✅ 表单 (form.*)
✅ 错误页面 (error.*)
✅ 语言切换 (language.*)
✅ 资源类型 (resources.*)
✅ 操作文本 (operations.*)

## 9. 组件改造优先级

1. **高优先级**：用户直接看到的页面（TenantList, TenantDetail, RoleList等）
2. **中优先级**：业务组件（TenantSelector, QuotaIndicator等）
3. **低优先级**：基础UI组件（Button, Input等）

## 10. 快速检查清单

- [ ] 组件导入了useTranslation Hook
- [ ] 所有硬编码文本替换为t()函数
- [ ] 动态内容使用插值语法
- [ ] 翻译key已添加到zh-CN.json和en-US.json
- [ ] 测试中英文切换正常

---

**相关文档**：
- [ZKER-企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)
- [i18next官方文档](https://www.i18next.com/)
