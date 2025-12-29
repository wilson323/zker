# 研发 C - 前端工程师 8 周开发计划 v1.0

> **角色定位**: 前端工程师 - 负责企业级前端组件库、租户/权限/路由管理页面、UX优化、多语言支持
>
> **工作目标**: 构建高质量、可复用、易维护的前端组件和页面系统，确保出色的用户体验
>
> **核心原则**: 组件化、一致性、性能优先、用户友好

---

## 📋 个人职责概述

### 核心负责模块

```
frontend/packages/
├── arch/ui-components/       ✅ 专属负责 - 基础UI组件库
│   ├── src/components/
│   │   ├── Button/
│   │   ├── Input/
│   │   ├── Modal/
│   │   └── Table/
├── arch/business-components/ ✅ 专属负责 - 业务组件库
│   ├── src/components/
│   │   ├── TenantSelector/
│   │   ├── PermissionTree/
│   │   └── QuotaIndicator/
├── studio/pages/tenant/       ✅ 专属负责 - 租户管理页面
│   ├── TenantList/
│   ├── TenantDetail/
│   └── Subscription/
├── studio/pages/permission/   ✅ 专属负责 - 权限管理页面
│   ├── RoleList/
│   ├── RoleDetail/
│   └── DataPermission/
├── studio/pages/routing/      ✅ 专属负责 - 路由配置页面
│   ├── RoutingRules/
│   └── IntentMatcher/
├── common/i18n/              ✅ 专属负责 - 国际化
│   ├── locales/
│   │   ├── zh-CN.json
│   │   └── en-US.json
└── common/themes/            ✅ 专属负责 - 主题系统
    └── tokens/
```

### 协作接口

| 协作对象 | 协作内容 | 接口定义位置 | 依赖关系 |
|---------|---------|------------|---------|
| **研发 A** | API调用、数据格式、前端契约 | `frontend/packages/api-client/` | 研发 A 提供 API → 研发 C 集成 |
| **研发 B** | 错误码映射、监控数据展示 | `frontend/packages/error-handler/` | 研发 B 提供错误码 → 研发 C 展示 |
| **研发 D** | 前端构建配置、静态资源部署 | `frontend/rush.json`、`docker/` | 研发 C 提供构建产物 → 研发 D 部署 |

---

## 🎯 8 周详细开发计划

### Week 1-2: 组件库 + 开发规范

#### Week 1: 基础UI组件库

**目标**: 建立统一的UI组件库，确保视觉一致性

##### Day 1-3: 设计令牌系统

**设计令牌定义** (`common/themes/tokens/index.ts`):
```typescript
// common/themes/tokens/index.ts

// 颜色系统
export const colors = {
  // 主色调
  primary: {
    50: '#E6F7FF',
    100: '#BAE7FF',
    200: '#91D5FF',
    300: '#69C0FF',
    400: '#40A9FF',
    500: '#1890FF',  // 主色
    600: '#096DD9',
    700: '#0050B3',
    800: '#003A8C',
    900: '#002766',
  },

  // 中性色
  gray: {
    50: '#FAFAFA',
    100: '#F5F5F5',
    200: '#E8E8E8',
    300: '#D9D9D9',
    400: '#BFBFBF',
    500: '#8C8C8C',
    600: '#595959',
    700: '#434343',
    800: '#262626',
    900: '#1F1F1F',
  },

  // 语义色
  success: {
    light: '#95DE64',
    main: '#52C41A',
    dark: '#389E0D',
  },
  warning: {
    light: '#FFD666',
    main: '#FAAD14',
    dark: '#D48806',
  },
  error: {
    light: '#FF7875',
    main: '#F5222D',
    dark: '#CF1322',
  },
  info: {
    light: '#91D5FF',
    main: '#1890FF',
    dark: '#0050B3',
  },
};

// 间距系统 (4的倍数)
export const spacing = {
  0: '0',
  1: '0.25rem',   // 4px
  2: '0.5rem',    // 8px
  3: '0.75rem',   // 12px
  4: '1rem',      // 16px
  5: '1.25rem',   // 20px
  6: '1.5rem',    // 24px
  8: '2rem',      // 32px
  10: '2.5rem',   // 40px
  12: '3rem',     // 48px
  16: '4rem',     // 64px
};

// 字体系统
export const typography = {
  fontFamily: {
    base: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
    mono: '"SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace',
  },
  fontSize: {
    xs: '0.75rem',    // 12px
    sm: '0.875rem',   // 14px
    base: '1rem',     // 16px
    lg: '1.125rem',   // 18px
    xl: '1.25rem',    // 20px
    '2xl': '1.5rem',  // 24px
    '3xl': '1.875rem', // 30px
    '4xl': '2.25rem',  // 36px
  },
  fontWeight: {
    normal: 400,
    medium: 500,
    semibold: 600,
    bold: 700,
  },
  lineHeight: {
    tight: 1.25,
    normal: 1.5,
    relaxed: 1.75,
  },
};

// 圆角系统
export const borderRadius = {
  none: '0',
  sm: '0.125rem',   // 2px
  base: '0.25rem',  // 4px
  md: '0.375rem',   // 6px
  lg: '0.5rem',     // 8px
  xl: '0.75rem',    // 12px
  '2xl': '1rem',    // 16px
  full: '9999px',
};

// 阴影系统
export const boxShadow = {
  sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
  base: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)',
  md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
  lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
  xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
};

// 断点系统
export const breakpoints = {
  sm: '640px',
  md: '768px',
  lg: '1024px',
  xl: '1280px',
  '2xl': '1536px',
};

// 过渡动画
export const transitions = {
  fast: '150ms cubic-bezier(0.4, 0, 0.2, 1)',
  base: '200ms cubic-bezier(0.4, 0, 0.2, 1)',
  slow: '300ms cubic-bezier(0.4, 0, 0.2, 1)',
};

// Z-index系统
export const zIndex = {
  dropdown: 1000,
  sticky: 1020,
  fixed: 1030,
  modalBackdrop: 1040,
  modal: 1050,
  popover: 1060,
  tooltip: 1070,
};
```

**主题Hook** (`common/themes/hooks/useTheme.ts`):
```typescript
// common/themes/hooks/useTheme.ts
import { useMemo } from 'react';
import { colors, spacing, typography, borderRadius, boxShadow, breakpoints } from '../tokens';

export interface Theme {
  colors: typeof colors;
  spacing: typeof spacing;
  typography: typeof typography;
  borderRadius: typeof borderRadius;
  boxShadow: typeof boxShadow;
  breakpoints: typeof breakpoints;
}

export const useTheme = (): Theme => {
  return useMemo(() => ({
    colors,
    spacing,
    typography,
    borderRadius,
    boxShadow,
    breakpoints,
  }), []);
};
```

##### Day 4-5: 基础组件

**Button组件** (`arch/ui-components/src/components/Button/Button.tsx`):
```typescript
// Button/Button.tsx
import React from 'react';
import { ButtonHTMLAttributes } from 'react';
import { useStyles } from './Button.styles';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'text' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  disabled?: boolean;
  icon?: React.ReactNode;
  block?: boolean;
  children: React.ReactNode;
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  loading = false,
  disabled = false,
  icon,
  block = false,
  children,
  className,
  ...rest
}) => {
  const classes = useStyles({ variant, size, block });

  return (
    <button
      className={`${classes.button} ${className || ''}`}
      disabled={disabled || loading}
      {...rest}
    >
      {loading && <span className={classes.spinner} />}
      {icon && !loading && <span className={classes.icon}>{icon}</span>}
      {children}
    </button>
  );
};
```

**Button样式** (`arch/ui-components/src/components/Button/Button.styles.ts`):
```typescript
// Button/Button.styles.ts
import { css } from '@emotion/react';
import { colors, spacing, borderRadius, transitions, typography } from '@coze-studio/common/themes';

interface ButtonStyleProps {
  variant: 'primary' | 'secondary' | 'outline' | 'text' | 'danger';
  size: 'sm' | 'md' | 'lg';
  block: boolean;
}

export const useStyles = (props: ButtonStyleProps) => {
  const { variant, size, block } = props;

  // 尺寸样式
  const sizeStyles = {
    sm: css`
      padding: ${spacing[2]} ${spacing[3]};
      ${typography.fontSize.sm};
    `,
    md: css`
      padding: ${spacing[3]} ${spacing[4]};
      ${typography.fontSize.base};
    `,
    lg: css`
      padding: ${spacing[4]} ${spacing[6]};
      ${typography.fontSize.lg};
    `,
  };

  // 变体样式
  const variantStyles = {
    primary: css`
      background-color: ${colors.primary[500]};
      color: white;
      &:hover:not(:disabled) {
        background-color: ${colors.primary[600]};
      }
      &:active:not(:disabled) {
        background-color: ${colors.primary[700]};
      }
    `,
    secondary: css`
      background-color: ${colors.gray[200]};
      color: ${colors.gray[800]};
      &:hover:not(:disabled) {
        background-color: ${colors.gray[300]};
      }
    `,
    outline: css`
      background-color: transparent;
      border: 1px solid ${colors.gray[300]};
      color: ${colors.gray[700]};
      &:hover:not(:disabled) {
        border-color: ${colors.primary[500]};
        color: ${colors.primary[500]};
      }
    `,
    text: css`
      background-color: transparent;
      color: ${colors.primary[500]};
      &:hover:not(:disabled) {
        background-color: ${colors.primary[50]};
      }
    `,
    danger: css`
      background-color: ${colors.error.main};
      color: white;
      &:hover:not(:disabled) {
        background-color: ${colors.error.dark};
      }
    `,
  };

  return {
    button: css`
      display: ${block ? 'block' : 'inline-flex'};
      width: ${block ? '100%' : 'auto'};
      align-items: center;
      justify-content: center;
      gap: ${spacing[2]};
      border: none;
      border-radius: ${borderRadius.base};
      font-weight: ${typography.fontWeight.medium};
      cursor: pointer;
      transition: all ${transitions.base};

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      ${sizeStyles[size]}
      ${variantStyles[variant]}
    `,
    spinner: css`
      width: 1em;
      height: 1em;
      border: 2px solid currentColor;
      border-top-color: transparent;
      border-radius: 50%;
      animation: spin 0.6s linear infinite;

      @keyframes spin {
        to { transform: rotate(360deg); }
      }
    `,
    icon: css`
      display: inline-flex;
      align-items: center;
    `,
  };
};
```

#### Week 2: 业务组件库

**目标**: 开发租户、权限、路由相关的业务组件

##### Day 1-3: 租户选择器

**TenantSelector组件** (`arch/business-components/src/components/TenantSelector/TenantSelector.tsx`):
```typescript
// TenantSelector/TenantSelector.tsx
import React, { useState, useEffect } from 'react';
import { useTenantList } from '@coze-studio/api-client';
import { Select } from '@coze-studio/ui-components';
import { useStyles } from './TenantSelector.styles';

export interface TenantSelectorProps {
  value?: string;
  onChange?: (tenantId: string) => void;
  disabled?: boolean;
  placeholder?: string;
  allowClear?: boolean;
}

export const TenantSelector: React.FC<TenantSelectorProps> = ({
  value,
  onChange,
  disabled = false,
  placeholder = '请选择租户',
  allowClear = true,
}) => {
  const classes = useStyles();
  const { data, loading } = useTenantList();

  const tenants = data?.tenants || [];

  const handleChange = (tenantId: string) => {
    onChange?.(tenantId);
  };

  return (
    <Select
      className={classes.selector}
      value={value}
      onChange={handleChange}
      disabled={disabled}
      loading={loading}
      placeholder={placeholder}
      allowClear={allowClear}
      options={tenants.map((tenant) => ({
        label: (
          <div className={classes.option}>
            <span className={classes.name}>{tenant.tenant_name}</span>
            <span className={classes.tag}>{tenant.tenant_type}</span>
          </div>
        ),
        value: tenant.tenant_id,
      }))}
    />
  );
};
```

##### Day 4-5: 配额指示器

**QuotaIndicator组件** (`arch/business-components/src/components/QuotaIndicator/QuotaIndicator.tsx`):
```typescript
// QuotaIndicator/QuotaIndicator.tsx
import React from 'react';
import { useStyles } from './QuotaIndicator.styles';

export interface QuotaIndicatorProps {
  used: number;
  max: number;
  resourceType: string;
  showLabel?: boolean;
}

export const QuotaIndicator: React.FC<QuotaIndicatorProps> = ({
  used,
  max,
  resourceType,
  showLabel = true,
}) => {
  const classes = useStyles({ used, max });

  const percentage = max > 0 ? (used / max) * 100 : 0;
  const isOverLimit = percentage >= 100;
  const isNearLimit = percentage >= 80 && percentage < 100;

  return (
    <div className={classes.container}>
      {showLabel && (
        <div className={classes.header}>
          <span className={classes.label}>{resourceType}</span>
          <span className={classes.count}>
            {used} / {max === -1 ? '∞' : max}
          </span>
        </div>
      )}
      <div className={classes.barContainer}>
        <div
          className={`${classes.bar} ${
            isOverLimit ? classes.overLimit : isNearLimit ? classes.nearLimit : ''
          }`}
          style={{ width: `${Math.min(percentage, 100)}%` }}
        />
      </div>
      {showLabel && (
        <div className={classes.footer}>
          <span className={classes.percentage}>{percentage.toFixed(1)}%</span>
        </div>
      )}
    </div>
  );
};
```

---

### Week 3-4: 业务页面开发

#### Week 3: 租户管理页面

**目标**: 完成租户列表、详情、订阅管理页面

##### Day 1-3: 租户列表页

**TenantList页面** (`studio/pages/tenant/TenantList/TenantList.tsx`):
```typescript
// TenantList/TenantList.tsx
import React, { useState } from 'react';
import { useTenantList, useDeleteTenant } from '@coze-studio/api-client';
import { Table, Button, Modal, message } from '@coze-studio/ui-components';
import { TenantFilter } from './components/TenantFilter';
import { TenantActions } from './components/TenantActions';
import { useStyles } from './TenantList.styles';

export const TenantList: React.FC = () => {
  const classes = useStyles();
  const [filter, setFilter] = useState({});
  const [pageToken, setPageToken] = useState('');
  const [pageSize] = useState(20);

  const { data, loading, refetch } = useTenantList({
    filter,
    pageToken,
    pageSize,
  });

  const { mutate: deleteTenant } = useDeleteTenant({
    onSuccess: () => {
      message.success('删除成功');
      refetch();
    },
  });

  const handleDelete = (tenantId: string, tenantName: string) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除租户"${tenantName}"吗？`,
      onOk: () => deleteTenant(tenantId),
    });
  };

  const columns = [
    {
      title: '租户名称',
      dataIndex: 'tenant_name',
      key: 'tenant_name',
      sorter: true,
    },
    {
      title: '租户类型',
      dataIndex: 'tenant_type',
      key: 'tenant_type',
      render: (type: string) => {
        const typeMap = {
          individual: '个人',
          team: '团队',
          enterprise: '企业',
        };
        return typeMap[type] || type;
      },
    },
    {
      title: '订阅等级',
      dataIndex: 'subscription_tier',
      key: 'subscription_tier',
      render: (tier: string) => {
        const tierMap = {
          free: '免费版',
          pro: '专业版',
          enterprise: '企业版',
        };
        return tierMap[tier] || tier;
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusMap = {
          active: '正常',
          suspended: '暂停',
          deleted: '已删除',
        };
        return statusMap[status] || status;
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <TenantActions
          tenant={record}
          onDelete={() => handleDelete(record.tenant_id, record.tenant_name)}
        />
      ),
    },
  ];

  return (
    <div className={classes.container}>
      <div className={classes.header}>
        <h1 className={classes.title}>租户管理</h1>
        <Button variant="primary">创建租户</Button>
      </div>

      <TenantFilter filter={filter} onChange={setFilter} />

      <Table
        className={classes.table}
        columns={columns}
        dataSource={data?.tenants || []}
        loading={loading}
        pagination={{
          current: pageToken ? 2 : 1,
          pageSize,
          total: data?.total_count || 0,
          onChange: () => {},
        }}
      />
    </div>
  );
};
```

##### Day 4-5: 租户详情页

**TenantDetail页面** (`studio/pages/tenant/TenantDetail/TenantDetail.tsx`):
```typescript
// TenantDetail/TenantDetail.tsx
import React from 'react';
import { useParams } from 'react-router-dom';
import { useTenant, useUpdateTenant } from '@coze-studio/api-client';
import { Tabs, Button, message } from '@coze-studio/ui-components';
import { TenantBasicInfo } from './components/TenantBasicInfo';
import { TenantSubscription } from './components/TenantSubscription';
import { TenantQuotas } from './components/TenantQuotas';
import { useStyles } from './TenantDetail.styles';

export const TenantDetail: React.FC = () => {
  const classes = useStyles();
  const { tenantId } = useParams<{ tenantId: string }>();
  const { data, loading, refetch } = useTenant(tenantId!);

  const { mutate: updateTenant } = useUpdateTenant({
    onSuccess: () => {
      message.success('更新成功');
      refetch();
    },
  });

  const handleUpdate = (values: any) => {
    updateTenant({ tenantId: tenantId!, data: values });
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className={classes.container}>
      <div className={classes.header}>
        <h1 className={classes.title}>{data?.tenant_name}</h1>
        <Button variant="outline">返回列表</Button>
      </div>

      <Tabs
        items={[
          {
            key: 'basic',
            label: '基本信息',
            children: (
              <TenantBasicInfo tenant={data} onUpdate={handleUpdate} />
            ),
          },
          {
            key: 'subscription',
            label: '订阅管理',
            children: (
              <TenantSubscription tenant={data} onUpdate={handleUpdate} />
            ),
          },
          {
            key: 'quotas',
            label: '配额管理',
            children: <TenantQuotas tenantId={tenantId!} />,
          },
        ]}
      />
    </div>
  );
};
```

#### Week 4: 权限管理页面

**目标**: 完成角色列表、角色详情、数据权限配置页面

##### Day 1-3: 角色列表页

**RoleList页面** (`studio/pages/permission/RoleList/RoleList.tsx`):
```typescript
// RoleList/RoleList.tsx
import React, { useState } from 'react';
import { useRoleList, useDeleteRole } from '@coze-studio/api-client';
import { Table, Button, Modal, message } from '@coze-studio/ui-components';
import { useStyles } from './RoleList.styles';

export const RoleList: React.FC = () => {
  const classes = useStyles();
  const [pageToken, setPageToken] = useState('');
  const [pageSize] = useState(20);

  const { data, loading, refetch } = useRoleList({
    pageToken,
    pageSize,
  });

  const { mutate: deleteRole } = useDeleteRole({
    onSuccess: () => {
      message.success('删除成功');
      refetch();
    },
  });

  const columns = [
    {
      title: '角色名称',
      dataIndex: 'role_name',
      key: 'role_name',
    },
    {
      title: '角色编码',
      dataIndex: 'role_code',
      key: 'role_code',
    },
    {
      title: '角色类型',
      dataIndex: 'role_type',
      key: 'role_type',
      render: (type: string) => (type === 'system' ? '系统角色' : '自定义角色'),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <div className={classes.actions}>
          <Button size="sm" variant="outline">
            编辑
          </Button>
          <Button
            size="sm"
            variant="text"
            onClick={() => handleDelete(record.role_id, record.role_name)}
          >
            删除
          </Button>
        </div>
      ),
    },
  ];

  return (
    <div className={classes.container}>
      <div className={classes.header}>
        <h1 className={classes.title}>角色管理</h1>
        <Button variant="primary">创建角色</Button>
      </div>

      <Table
        columns={columns}
        dataSource={data?.roles || []}
        loading={loading}
        pagination={{
          pageSize,
          total: data?.total_count || 0,
        }}
      />
    </div>
  );
};
```

##### Day 4-5: 数据权限配置

**DataPermissionConfig组件** (`studio/pages/permission/components/DataPermissionConfig/DataPermissionConfig.tsx`):
```typescript
// DataPermissionConfig/DataPermissionConfig.tsx
import React, { useState } from 'react';
import { Radio, Select, Form, Input, Button } from '@coze-studio/ui-components';
import { useStyles } from './DataPermissionConfig.styles';

interface PermissionConfig {
  resourceType: string;
  scope: 'ALL' | 'DEPARTMENT' | 'OWN' | 'CUSTOM' | 'NONE';
  customFilter?: any;
}

export const DataPermissionConfig: React.FC = () => {
  const classes = useStyles();
  const [permissions, setPermissions] = useState<PermissionConfig[]>([
    { resourceType: 'bots', scope: 'ALL' },
    { resourceType: 'conversations', scope: 'OWN' },
    { resourceType: 'knowledge', scope: 'DEPARTMENT' },
  ]);

  const handleScopeChange = (index: number, scope: PermissionConfig['scope']) => {
    const newPermissions = [...permissions];
    newPermissions[index].scope = scope;
    setPermissions(newPermissions);
  };

  const resourceTypeOptions = [
    { label: '机器人', value: 'bots' },
    { label: '对话', value: 'conversations' },
    { label: '知识库', value: 'knowledge' },
    { label: '工作流', value: 'workflows' },
    { label: '插件', value: 'plugins' },
  ];

  const scopeOptions = [
    { label: '全部数据', value: 'ALL' },
    { label: '本部门数据', value: 'DEPARTMENT' },
    { label: '仅自己', value: 'OWN' },
    { label: '自定义', value: 'CUSTOM' },
    { label: '无权限', value: 'NONE' },
  ];

  return (
    <div className={classes.container}>
      <h3 className={classes.title}>数据权限配置</h3>

      {permissions.map((perm, index) => (
        <div key={index} className={classes.permissionRow}>
          <span className={classes.resourceType}>
            {resourceTypeOptions.find((opt) => opt.value === perm.resourceType)?.label}
          </span>

          <Radio.Group
            value={perm.scope}
            onChange={(e) => handleScopeChange(index, e.target.value)}
            options={scopeOptions}
          />

          {perm.scope === 'CUSTOM' && (
            <Form className={classes.customFilter}>
              <Input.TextArea
                placeholder="输入自定义过滤条件（JSON格式）"
                rows={3}
              />
            </Form>
          )}
        </div>
      ))}

      <div className={classes.footer}>
        <Button variant="primary">保存</Button>
        <Button variant="outline">取消</Button>
      </div>
    </div>
  );
};
```

---

### Week 5-6: UX优化 + 多语言

#### Week 5: UX优化

**目标**: 优化用户体验，提升交互流畅度

##### Day 1-3: 加载状态优化

**LoadingWrapper组件** (`common/components/LoadingWrapper/LoadingWrapper.tsx`):
```typescript
// LoadingWrapper/LoadingWrapper.tsx
import React from 'react';
import { useStyles } from './LoadingWrapper.styles';

export interface LoadingWrapperProps {
  loading: boolean;
  error?: Error | null;
  empty?: boolean;
  children: React.ReactNode;
}

export const LoadingWrapper: React.FC<LoadingWrapperProps> = ({
  loading,
  error,
  empty,
  children,
}) => {
  const classes = useStyles();

  if (loading) {
    return (
      <div className={classes.container}>
        <div className={classes.spinner} />
        <p className={classes.text}>加载中...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={classes.container}>
        <div className={classes.error}>⚠️</div>
        <p className={classes.text}>加载失败: {error.message}</p>
      </div>
    );
  }

  if (empty) {
    return (
      <div className={classes.container}>
        <div className={classes.empty}>📭</div>
        <p className={classes.text}>暂无数据</p>
      </div>
    );
  }

  return <>{children}</>;
};
```

##### Day 4-5: 错误边界

**ErrorBoundary组件** (`common/components/ErrorBoundary/ErrorBoundary.tsx`):
```typescript
// ErrorBoundary/ErrorBoundary.tsx
import React, { Component, ErrorInfo, ReactNode } from 'react';
import { Button } from '@coze-studio/ui-components';
import { useStyles } from './ErrorBoundary.styles';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ErrorBoundary caught:', error, errorInfo);
  }

  handleReset = () => {
    this.setState({ hasError: false, error: undefined });
    window.location.reload();
  };

  render() {
    if (this.state.hasError) {
      return (
        <div className="error-boundary">
          <h1>出错了</h1>
          <p>{this.state.error?.message}</p>
          <Button onClick={this.handleReset}>重新加载</Button>
        </div>
      );
    }

    return this.props.children;
  }
}
```

#### Week 6: 国际化支持

**目标**: 实现中英文双语支持

##### Day 1-3: i18n配置

**i18n配置** (`common/i18n/config.ts`):
```typescript
// common/i18n/config.ts
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import Backend from 'i18next-http-backend';
import zhCN from './locales/zh-CN.json';
import enUS from './locales/en-US.json';

const resources = {
  'zh-CN': {
    translation: zhCN,
  },
  'en-US': {
    translation: enUS,
  },
};

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'zh-CN',
    lng: 'zh-CN',
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
```

**中文翻译** (`common/i18n/locales/zh-CN.json`):
```json
{
  "common": {
    "confirm": "确认",
    "cancel": "取消",
    "save": "保存",
    "delete": "删除",
    "edit": "编辑",
    "create": "创建",
    "search": "搜索",
    "loading": "加载中...",
    "noData": "暂无数据",
    "error": "出错了"
  },
  "tenant": {
    "title": "租户管理",
    "list": "租户列表",
    "detail": "租户详情",
    "create": "创建租户",
    "edit": "编辑租户",
    "delete": "删除租户",
    "name": "租户名称",
    "type": "租户类型",
    "subscription": "订阅等级",
    "status": "状态",
    "createdAt": "创建时间",
    "types": {
      "individual": "个人",
      "team": "团队",
      "enterprise": "企业"
    },
    "tiers": {
      "free": "免费版",
      "pro": "专业版",
      "enterprise": "企业版"
    }
  },
  "permission": {
    "title": "权限管理",
    "role": "角色",
    "dataScope": "数据权限",
    "fieldPermission": "字段权限",
    "scopes": {
      "ALL": "全部数据",
      "DEPARTMENT": "本部门数据",
      "OWN": "仅自己",
      "CUSTOM": "自定义",
      "NONE": "无权限"
    }
  }
}
```

**英文翻译** (`common/i18n/locales/en-US.json`):
```json
{
  "common": {
    "confirm": "Confirm",
    "cancel": "Cancel",
    "save": "Save",
    "delete": "Delete",
    "edit": "Edit",
    "create": "Create",
    "search": "Search",
    "loading": "Loading...",
    "noData": "No Data",
    "error": "Error"
  },
  "tenant": {
    "title": "Tenant Management",
    "list": "Tenant List",
    "detail": "Tenant Detail",
    "create": "Create Tenant",
    "edit": "Edit Tenant",
    "delete": "Delete Tenant",
    "name": "Tenant Name",
    "type": "Tenant Type",
    "subscription": "Subscription Tier",
    "status": "Status",
    "createdAt": "Created At",
    "types": {
      "individual": "Individual",
      "team": "Team",
      "enterprise": "Enterprise"
    },
    "tiers": {
      "free": "Free",
      "pro": "Pro",
      "enterprise": "Enterprise"
    }
  },
  "permission": {
    "title": "Permission Management",
    "role": "Role",
    "dataScope": "Data Permission",
    "fieldPermission": "Field Permission",
    "scopes": {
      "ALL": "All Data",
      "DEPARTMENT": "Department Data",
      "OWN": "Own Only",
      "CUSTOM": "Custom",
      "NONE": "No Permission"
    }
  }
}
```

##### Day 4-5: 语言切换器

**LanguageSwitcher组件** (`common/components/LanguageSwitcher/LanguageSwitcher.tsx`):
```typescript
// LanguageSwitcher/LanguageSwitcher.tsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import { Select } from '@coze-studio/ui-components';

export const LanguageSwitcher: React.FC = () => {
  const { i18n } = useTranslation();

  const options = [
    { label: '简体中文', value: 'zh-CN' },
    { label: 'English', value: 'en-US' },
  ];

  const handleChange = (lang: string) => {
    i18n.changeLanguage(lang);
  };

  return (
    <Select
      value={i18n.language}
      onChange={handleChange}
      options={options}
      style={{ width: 120 }}
    />
  );
};
```

---

### Week 7-8: 集成测试 + 性能优化

#### Week 7: 组件测试

**目标**: 编写完整的组件测试用例

##### Day 1-5: 组件单元测试

**Button组件测试** (`arch/ui-components/src/components/Button/Button.test.tsx`):
```typescript
// Button.test.tsx
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from './Button';

describe('Button', () => {
  it('renders children correctly', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick when clicked', () => {
    const handleClick = jest.fn();
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
});
```

#### Week 8: 性能优化

**目标**: 优化前端性能，提升加载速度

##### Day 1-3: 代码分割

**路由级代码分割** (`studio/App.tsx`):
```typescript
// App.tsx
import { lazy, Suspense } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { LoadingWrapper } from '@coze-studio/common/components';

// 懒加载页面组件
const TenantList = lazy(() => import('./pages/tenant/TenantList'));
const TenantDetail = lazy(() => import('./pages/tenant/TenantDetail'));
const RoleList = lazy(() => import('./pages/permission/RoleList'));

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Suspense fallback={<LoadingWrapper loading />}>
        <Routes>
          <Route path="/tenants" element={<TenantList />} />
          <Route path="/tenants/:tenantId" element={<TenantDetail />} />
          <Route path="/permissions/roles" element={<RoleList />} />
        </Routes>
      </Suspense>
    </BrowserRouter>
  );
};
```

##### Day 4-5: 构建优化

**Rsbuild配置优化** (`studio/rsbuild.config.ts`):
```typescript
// rsbuild.config.ts
import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';

export default defineConfig({
  plugins: [pluginReact()],
  output: {
    distPath: {
      root: 'dist',
    },
    // 代码分割策略
    splitChunks: {
      strategy: 'split-by-experience',
    },
    // 压缩配置
    minify: 'swc',
    // 目标浏览器
    polyfill: 'entry',
    targets: ['defaults', 'not IE 11'],
  },
  performance: {
    // 移除 console
    removeConsole: process.env.NODE_ENV === 'production',
    // 移除 moment locales
    removeMomentLocale: true,
    // 打包体积分析
    bundleAnalyze: process.env.ANALYZE === 'true',
  },
  // 环境变量
  env: {
    API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:8080',
  },
});
```

---

## 📦 专属开发规范

### 作为前端工程师的特殊要求

#### 1. 组件化原则

✅ **必须**：
- 单一职责：一个组件只做一件事
- Props 清晰：所有输入通过 Props
- 可复用：避免业务逻辑硬编码
- 可测试：纯函数优先

```typescript
// ✅ Good
export const Button: React.FC<ButtonProps> = ({ variant, size, children }) => {
  return <button className={`btn btn-${variant} btn-${size}`}>{children}</button>;
};

// ❌ Bad
export const Button: React.FC = () => {
  const [data, setData] = useState(); // 不应该有状态
  useEffect(() => { fetchData(); }, []); // 不应该有副作用
  return <button>Click</button>;
};
```

#### 2. TypeScript 规范

✅ **必须**：
- 所有组件有明确的 Props 类型
- 避免使用 any
- 使用 interface 定义对象类型
- 使用 type 定义联合类型

```typescript
// ✅ Good
export interface ButtonProps {
  variant?: 'primary' | 'secondary';
  size?: 'sm' | 'md' | 'lg';
  children: React.ReactNode;
}

// ❌ Bad
export const Button = (props: any) => { ... };
```

#### 3. 性能优化

✅ **必须**：
- 使用 React.memo 避免不必要的重渲染
- 使用 useMemo 缓存计算结果
- 使用 useCallback 稳定函数引用

```typescript
// ✅ Good
export const ExpensiveComponent: React.FC<Props> = React.memo(({ data }) => {
  const sortedData = useMemo(() => {
    return data.sort((a, b) => a.createdAt - b.createdAt);
  }, [data]);

  return <div>{sortedData.map(item => <Item key={item.id} />)}</div>;
});

// ❌ Bad
export const ExpensiveComponent: React.FC<Props> = ({ data }) => {
  const sortedData = data.sort((a, b) => a.createdAt - b.createdAt); // 每次渲染都计算
  return <div>{sortedData.map(item => <Item />)}</div>;
};
```

---

## 🤝 协作接口定义

### 与研发 A 协作

**依赖研发 A**:
- API 契约（OpenAPI）
- 数据模型定义
- 错误码定义

**向研发 A 提供**:
- 前端 Mock 数据
- UI/UX 反馈
- API 调用示例

### 与研发 B 协作

**依赖研发 B**:
- 错误码 JSON 文件
- 监控数据展示接口

**向研发 B 提供**:
- 前端性能指标
- 前端错误日志

### 与研发 D 协作

**依赖研发 D**:
- 静态资源部署
- CDN 配置

**向研发 D 提供**:
- 构建产物
- 部署配置

---

## 📅 里程碑和交付物

### Week 2 交付物

- [ ] 设计令牌系统
- [ ] 10+ 基础UI组件
- [ ] 5+ 业务组件
- [ ] 组件 Storybook 文档

### Week 4 交付物

- [ ] 租户管理页面（列表+详情）
- [ ] 权限管理页面（角色+数据权限）
- [ ] 路由配置页面
- [ ] 组件单元测试

### Week 6 交付物

- [ ] UX优化（加载状态、错误处理）
- [ ] 中英文双语支持
- [ ] 语言切换器
- [ ] 组件文档

### Week 8 交付物

- [ ] 完整集成测试
- [ ] 性能优化
- [ ] 构建优化
- [ ] 用户手册

---

## ✅ 质量检查清单

### 代码提交前

- [ ] 组件符合设计令牌系统
- [ ] TypeScript 类型完整
- [ ] 组件有 PropTypes/Interface
- [ ] 单元测试覆盖率 ≥ 70%
- [ ] ESLint 检查通过

### 集成前

- [ ] 与后端 API 集成测试通过
- [ ] 页面响应式适配
- [ ] 多语言切换正常
- [ ] 浏览器兼容性测试

### 发布前

- [ ] 构建成功无警告
- [ ] 性能指标达标（LCP < 2.5s）
- [ ] 无障碍访问测试
- [ ] 用户手册完整

---

## 📊 关键指标

### 代码质量

- TypeScript 覆盖率：100%
- 组件测试覆盖率：≥ 70%
- ESLint 警告数：0
- 组件复用率：≥ 80%

### 性能指标

- First Contentful Paint (FCP)：< 1.5s
- Largest Contentful Paint (LCP)：< 2.5s
- Time to Interactive (TTI)：< 3.5s
- Cumulative Layout Shift (CLS)：< 0.1

### 功能完整性

- Week 2: 组件库 100%
- Week 4: 业务页面 100%
- Week 6: UX优化 100%
- Week 8: 国际化 100%

---

## 📚 参考资料

- [ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md)
- [ZKER-全局一致性检查清单_v1.0.md](../ZKER-全局一致性检查清单_v1.0.md)
- [ZKER-实现差距分析与研发计划_v1.0.md](../ZKER-实现差距分析与研发计划_v1.0.md)

---

**文档状态**: ✅ 已完成
**最后更新**: 2025-01-01
**责任人**: 研发 C - 前端工程师
**评审人**: 技术负责人
