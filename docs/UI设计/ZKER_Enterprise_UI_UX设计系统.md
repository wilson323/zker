# zker Enterprise - UI/UX 界面设计方案

> **文档版本**: v1.0
> **创建日期**: 2025-12-30
> **适用范围**: 前端开发、UI设计、产品团队
> **技术栈**: React 18 + Semi Design + Tailwind CSS + TypeScript

---

## 📋 目录

1. [设计系统基础](#一设计系统基础)
2. [页面布局设计](#二页面布局设计)
3. [组件设计规范](#三组件设计规范)
4. [三大前台页面设计](#四三大前台页面设计)
5. [交互流程设计](#五交互流程设计)
6. [响应式设计](#六响应式设计)
7. [可访问性设计](#七可访问性设计)
8. [设计交付物](#八设计交付物)

---

## 一、设计系统基础

### 1.1 设计原则

```
┌──────────────────────────────────────────────────────────┐
│                  核心设计原则                              │
├──────────────────────────────────────────────────────────┤
│                                                            │
│  ✨ 简洁优先 (Simplicity First)                             │
│     - 减少视觉噪音，突出核心内容                              │
│     - 遵循"少即是多"设计哲学                                  │
│                                                            │
│  🎯 清晰明确 (Clarity)                                      │
│     - 信息层级分明，用户一目了然                              │
│     - 交互反馈及时，状态清晰可见                              │
│                                                            │
│  🚀 高效易用 (Efficiency)                                   │
│     - 最少步骤完成核心任务                                    │
│     - 智能推荐，减少用户思考成本                              │
│                                                            │
│  🎨 一致性 (Consistency)                                    │
│     - 视觉风格统一，组件规范一致                              │
│     - 交互模式一致，降低学习成本                              │
│                                                            │
│  💙 用户为中心 (User-Centric)                               │
│     - 基于真实用户场景设计                                    │
│     - 持续优化，数据驱动迭代                                  │
│                                                            │
└──────────────────────────────────────────────────────────┘
```

### 1.2 色彩系统

#### 1.2.1 品牌色

```css
/* =====================================================
   品牌主色 - 科技蓝
   ===================================================== */

/* 主色（Primary）- 用于主要操作、主要CTA按钮 */
--primary-color: #306EEF;        /* 主题蓝 */
--primary-hover: #255DD8;        /* 悬停态 */
--primary-active: #1F4CB8;       /* 激活态 */
--primary-disabled: #AABFFF;     /* 禁用态 */

/* 辅助色（Secondary）- 用于次要操作 */
--secondary-color: #6B7280;      /* 中性灰 */
--secondary-hover: #4B5563;
--secondary-active: #374151;

/* 成功色（Success）- 用于成功状态 */
--success-color: #10B981;        /* 绿色 */
--success-bg: #D1FAE5;
--success-border: #34D399;

/* 警告色（Warning）- 用于警告状态 */
--warning-color: #F59E0B;        /* 橙色 */
--warning-bg: #FEF3C7;
--warning-border: #FBBF24;

/* 错误色（Error）- 用于错误状态 */
--error-color: #EF4444;          /* 红色 */
--error-bg: #FEE2E2;
--error-border: #F87171;

/* 信息色（Info）- 用于信息提示 */
--info-color: #3B82F6;           /* 蓝色 */
--info-bg: #DBEAFE;
--info-border: #60A5FA;
```

#### 1.2.2 中性色

```css
/* =====================================================
   中性色系统 - 用于文本、背景、边框
   ===================================================== */

/* 文本颜色 */
--text-primary: #111827;      /* 主要文本 - 标题、正文 */
--text-secondary: #6B7280;    /* 次要文本 - 辅助说明 */
--text-tertiary: #9CA3AF;     /* 三级文本 - 占位符 */
--text-quaternary: #D1D5DB;   /* 四级文本 - 禁用文本 */
--text-inverse: #FFFFFF;       /* 反色文本 - 深色背景用 */

/* 背景颜色 */
--bg-primary: #FFFFFF;        /* 主背景 - 页面背景 */
--bg-secondary: #F9FAFB;      /* 次要背景 - 卡片背景 */
--bg-tertiary: #F3F4F6;       /* 三级背景 - 悬浮层背景 */
--bg-overlay: rgba(0, 0, 0, 0.5);  /* 遮罩层背景 */

/* 边框颜色 */
--border-primary: #E5E7EB;    /* 主边框 */
--border-secondary: #D1D5DB;  /* 次要边框 */
--border-focus: #306EEF;       /* 焦点边框 */
```

#### 1.2.3 功能色

```css
/* =====================================================
   功能色 - 用于特定功能标识
   ===================================================== */

/* 链接色 */
--link-color: #306EEF;
--link-hover: #1F4CB8;
--link-active: #1E40AF;
--link-visited: #6366F1;

/* 在线状态 */
--status-online: #10B981;     /* 在线 - 绿色 */
--status-busy: #F59E0B;        /* 忙碌 - 橙色 */
--status-offline: #9CA3AF;     /* 离线 - 灰色 */
--status-away: #6B7280;        /* 离开 - 深灰 */

/* 标签色 */
--tag-blue: #3B82F6;
--tag-green: #10B981;
--tag-orange: #F59E0B;
--tag-red: #EF4444;
--tag-purple: #8B5CF6;
--tag-cyan: #06B6D4;
```

### 1.3 字体系统

#### 1.3.1 字体家族

```css
/* =====================================================
   字体家族 - 中文和英文
   ===================================================== */

/* 字体栈（优先使用系统默认字体，提升性能） */
--font-family-base: -apple-system, BlinkMacSystemFont, "Segoe UI",
                    Roboto, "Helvetica Neue", Arial,
                    "Noto Sans", sans-serif,
                    "Apple Color Emoji", "Segoe UI Emoji",
                    "Segoe UI Symbol", "Noto Color Emoji";

/* 数字字体（用于数据显示） */
--font-family-number: -apple-system, BlinkMacSystemFont, "Segoe UI",
                       Roboto, "Droid Sans", "Helvetica Neue", Arial,
                       sans-serif;

/* 代码字体（等宽字体） */
--font-family-code: "SF Mono", Monaco, "Cascadia Code",
                    "Roboto Mono", Consolas,
                    "Courier New", monospace;
```

#### 1.3.2 字体大小与行高

```css
/* =====================================================
   字体大小规范 - 使用8点栅格系统
   ===================================================== */

/* 字体大小（8点栅格） */
--font-size-xs: 12px;          /* 极小号 - 辅助信息 */
--font-size-sm: 14px;          /* 小号 - 次要文本 */
--font-size-base: 16px;        /* 基础号 - 正文 */
--font-size-lg: 18px;          /* 大号 - 小标题 */
--font-size-xl: 20px;          /* 超大号 - 标题 */
--font-size-2xl: 24px;         /* 2倍大 - 二级标题 */
--font-size-3xl: 30px;         /* 3倍大 - 一级标题 */
--font-size-4xl: 36px;         /* 4倍大 - 特大标题 */

/* 行高（提升可读性） */
--line-height-tight: 1.25;     /* 紧凑 - 标题 */
--line-height-normal: 1.5;     /* 正常 - 正文 */
--line-height-relaxed: 1.75;   /* 宽松 - 长文本 */

/* 字重 */
--font-weight-normal: 400;     /* 常规 */
--font-weight-medium: 500;     /* 中等 */
--font-weight-semibold: 600;   /* 半粗 */
--font-weight-bold: 700;       /* 粗体 */
```

#### 1.3.3 字体使用示例

```typescript
// TypeScript 类型定义
type FontSize =
  | 'xs'    // 12px - 辅助信息、备注
  | 'sm'    // 14px - 次要文本、标签
  | 'base'  // 16px - 正文、表单标签
  | 'lg'    // 18px - 小标题、卡片标题
  | 'xl'    // 20px - 页面标题
  | '2xl'   // 24px - 二级标题
  | '3xl'   // 30px - 一级标题
  | '4xl';  // 36px - 特大标题（首页）

type FontWeight =
  | 'normal'    // 400 - 正文
  | 'medium'    // 500 - 强调
  | 'semibold' // 600 - 小标题
  | 'bold';     // 700 - 大标题

// 使用示例
<Text size="base" weight="normal">
  这是正文文本，用于常规内容展示
</Text>

<Text size="xl" weight="semibold">
  这是页面标题
</Text>
```

### 1.4 间距系统

```css
/* =====================================================
   间距系统 - 基于8点栅格
   ===================================================== */

/* 基础间距单位 */
--spacing-xs: 4px;      /* 0.25rem - 极小间距 */
--spacing-sm: 8px;      /* 0.5rem  - 小间距 */
--spacing-md: 16px;     /* 1rem    - 中等间距 */
--spacing-lg: 24px;     /* 1.5rem  - 大间距 */
--spacing-xl: 32px;     /* 2rem    - 超大间距 */
--spacing-2xl: 48px;    /* 3rem    - 特大间距 */
--spacing-3xl: 64px;    /* 4rem    - 巨大间距 */

/* 布局间距 */
--layout-padding-sm: 16px;   /* 小屏幕内边距 */
--layout-padding-md: 24px;   /* 中等屏幕内边距 */
--layout-padding-lg: 32px;   /* 大屏幕内边距 */
--layout-padding-xl: 48px;   /* 超大屏幕内边距 */

/* 组件间距 */
--component-gap-sm: 8px;     /* 组件内部小间距 */
--component-gap-md: 16px;    /* 组件内部中间距 */
--component-gap-lg: 24px;    /* 组件内部大间距 */
```

### 1.5 圆角系统

```css
/* =====================================================
   圆角系统 - 层级化设计
   ===================================================== */

/* 圆角大小 */
--radius-xs: 2px;      /* 极小圆角 - 标签、徽章 */
--radius-sm: 4px;      /* 小圆角 - 按钮、输入框 */
--radius-md: 6px;      /* 中等圆角 - 卡片 */
--radius-lg: 8px;      /* 大圆角 - 模态框 */
--radius-xl: 12px;     /* 超大圆角 - 大卡片 */
--radius-2xl: 16px;    /* 特大圆角 - 特殊组件 */
--radius-full: 9999px; /* 完全圆角 - 圆形按钮、头像 */

/* 使用场景 */
--radius-button: var(--radius-sm);        /* 按钮圆角 */
--radius-input: var(--radius-sm);         /* 输入框圆角 */
--radius-card: var(--radius-md);          /* 卡片圆角 */
--radius-modal: var(--radius-lg);         /* 模态框圆角 */
--radius-avatar: var(--radius-full);      /* 头像圆角 */
```

### 1.6 阴影系统

```css
/* =====================================================
   阴影系统 - 层级化设计
   ===================================================== */

/* 阴影层级 */
--shadow-xs: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
      /* 极小阴影 - 标签、徽章 */

--shadow-sm: 0 1px 3px 0 rgba(0, 0, 0, 0.1),
            0 1px 2px 0 rgba(0, 0, 0, 0.06);
      /* 小阴影 - 按钮、输入框 */

--shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1),
            0 2px 4px -1px rgba(0, 0, 0, 0.06);
      /* 中等阴影 - 卡片、下拉菜单 */

--shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1),
            0 4px 6px -2px rgba(0, 0, 0, 0.05);
      /* 大阴影 - 模态框、抽屉 */

--shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1),
            0 10px 10px -5px rgba(0, 0, 0, 0.04);
      /* 超大阴影 - 弹出层、提示框 */

--shadow-2xl: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
      /* 特大阴影 - 全屏遮罩层 */

/* 内阴影 */
--shadow-inner: inset 0 2px 4px 0 rgba(0, 0, 0, 0.06);
      /* 内阴影 - 输入框聚焦 */
```

### 1.7 动画与过渡

```css
/* =====================================================
   动画系统 - 平滑自然的过渡效果
   ===================================================== */

/* 过渡时长 */
--duration-fast: 150ms;      /* 快速 - 悬停效果 */
--duration-base: 200ms;      /* 基础 - 淡入淡出 */
--duration-normal: 300ms;    /* 正常 - 展开/收起 */
--duration-slow: 500ms;      /* 慢速 - 复杂动画 */

/* 缓动函数 */
--ease-in: cubic-bezier(0.4, 0, 1, 1);
      /* 加速 */
--ease-out: cubic-bezier(0, 0, 0.2, 1);
      /* 减速 */
--ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);
      /* 加速后减速 */
--ease-bounce: cubic-bezier(0.68, -0.55, 0.265, 1.55);
      /* 弹跳效果 */

/* 动画示例 */
.transition-all {
  transition-property: all;
  transition-duration: var(--duration-base);
  transition-timing-function: var(--ease-in-out);
}

.transition-colors {
  transition-property: color, background-color, border-color;
  transition-duration: var(--duration-fast);
  transition-timing-function: var(--ease-in-out);
}

.transition-transform {
  transition-property: transform;
  transition-duration: var(--duration-normal);
  transition-timing-function: var(--ease-out);
}
```

---

## 二、页面布局设计

### 2.1 整体布局架构

```
┌─────────────────────────────────────────────────────────────┐
│                    顶部导航栏 (Header)                        │
│  [Logo] [导航菜单] [搜索] [通知] [头像]                     │
├──────────┬──────────────────────────────────────────────────┤
│          │                                                   │
│  侧边栏   │            主内容区 (Main Content)              │
│ (Sidebar) │                                                   │
│          │  ┌───────────────────────────────────────────┐  │
│ [一级菜单] │  │                                           │  │
│          │  │           页面内容                         │  │
│ [二级菜单] │  │                                           │  │
│          │  │                                           │  │
│ [功能入口] │  │                                           │  │
│          │  │                                           │  │
│ [用户信息] │  │                                           │  │
│          │  └───────────────────────────────────────────┘  │
│          │                                                   │
├──────────┴──────────────────────────────────────────────────┤
│                    页脚 (Footer) - 可选                      │
└─────────────────────────────────────────────────────────────┘

布局类型：
- 经典布局（上图）：左侧导航 + 右侧内容 - 适用于管理后台
- 顶部布局：顶部导航 + 全宽内容 - 适用于落地页
- 混合布局：顶部一级菜单 + 左侧二级菜单 - 适用于复杂系统
```

### 2.2 顶部导航栏设计

#### 2.2.1 导航栏结构

```typescript
/**
 * 顶部导航栏布局规范
 * 高度：64px (固定)
 * 背景：白色 (#FFFFFF)
 * 边框：底部1px边框 (#E5E7EB)
 * 阴影：小阴影 (shadow-sm)
 */

interface HeaderProps {
  // Logo区域
  logo: {
    src: string;           // Logo图片URL
    text: string;          // 品牌名称
    href: string;          // 首页链接
  };

  // 导航菜单
  navigation: {
    items: NavItem[];      // 菜单项
    activeKey: string;     // 当前激活项
  };

  // 搜索区域
  search?: {
    placeholder: string;   // 占位文本
    onSearch: (query: string) => void;
  };

  // 右侧操作区
  actions: {
    notifications?: number;  // 通知数量（红点）
    messages?: number;        // 消息数量（红点）
    user: {
      name: string;           // 用户名
      avatar: string;         // 头像URL
      menu: MenuItem[];       // 下拉菜单
    };
  };
}

interface NavItem {
  key: string;
  label: string;
  icon?: React.ReactNode;
  href?: string;
  badge?: number | 'dot';     // 徽章数字或红点
  disabled?: boolean;
}

// 使用示例
<Header
  logo={{
    src: '/logo.svg',
    text: 'zker Enterprise',
    href: '/'
  }}
  navigation={{
    items: [
      { key: 'workspace', label: '工作台', icon: <IconDesktop />, href: '/workspace' },
      { key: 'employees', label: '数字员工', icon: <IconRobot />, href: '/employees' },
      { key: 'knowledge', label: '知识库', icon: <IconBook />, href: '/knowledge' },
      { key: 'analytics', label: '数据分析', icon: <IconChart />, href: '/analytics' },
    ],
    activeKey: 'workspace'
  }}
  search={{
    placeholder: '搜索数字员工、知识库...',
    onSearch: (query) => console.log(query)
  }}
  actions={{
    notifications: 5,
    user: {
      name: '张三',
      avatar: '/avatar.jpg',
      menu: [
        { key: 'profile', label: '个人中心' },
        { key: 'settings', label: '系统设置' },
        { key: 'logout', label: '退出登录' }
      ]
    }
  }}
/>
```

#### 2.2.2 导航栏样式规范

```css
/* =====================================================
   顶部导航栏样式
   ===================================================== */

.header {
  /* 布局 */
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 var(--layout-padding-md);

  /* 视觉 */
  background-color: var(--bg-primary);
  border-bottom: 1px solid var(--border-primary);
  box-shadow: var(--shadow-sm);

  /* 定位 */
  position: sticky;
  top: 0;
  z-index: 100;
}

/* Logo区域 */
.header-logo {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  text-decoration: none;
}

.header-logo img {
  width: 32px;
  height: 32px;
}

.header-logo-text {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
}

/* 导航菜单 */
.header-nav {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-sm) var(--spacing-md);
  color: var(--text-secondary);
  text-decoration: none;
  border-radius: var(--radius-sm);
  transition: all var(--duration-fast) var(--ease-in-out);
}

.nav-item:hover {
  background-color: var(--bg-secondary);
  color: var(--text-primary);
}

.nav-item.active {
  background-color: var(--primary-color);
  color: var(--text-inverse);
}

/* 搜索框 */
.header-search {
  width: 320px;
}

/* 右侧操作区 */
.header-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.icon-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-in-out);
}

.icon-button:hover {
  background-color: var(--bg-secondary);
  color: var(--text-primary);
}

/* 用户头像 */
.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
  cursor: pointer;
  border: 2px solid var(--border-primary);
  transition: border-color var(--duration-fast) var(--ease-in-out);
}

.user-avatar:hover {
  border-color: var(--primary-color);
}
```

### 2.3 侧边栏设计

#### 2.3.1 侧边栏结构

```typescript
/**
 * 侧边栏布局规范
 * 宽度：240px (展开) / 72px (收起)
 * 背景：白色 (#FFFFFF)
 * 边框：右侧1px边框 (#E5E7EB)
 */

interface SidebarProps {
  // 侧边栏模式
  mode: 'expanded' | 'collapsed';

  // 菜单数据
  menus: MenuGroup[];

  // 用户信息（底部）
  user?: {
    name: string;
    avatar: string;
    role: string;
  };

  // 折叠状态
  collapsible?: boolean;
}

interface MenuGroup {
  title?: string;           // 分组标题
  items: MenuItem[];        // 菜单项
}

interface MenuItem {
  key: string;
  label: string;
  icon?: React.ReactNode;
  path?: string;
  badge?: number | 'dot';
  disabled?: boolean;
  children?: MenuItem[];   // 子菜单
}

// 使用示例
<Sidebar
  mode="expanded"
  collapsible={true}
  menus={[
    {
      title: '员工前台',
      items: [
        { key: 'chatbox', label: '超级对话', icon: <IconChat />, path: '/chatbox' },
        { key: 'employees', label: '数字员工', icon: <IconRobot />, path: '/employees' },
        { key: 'chatbi', label: '智能问数', icon: <IconBarChart />, path: '/analytics' },
        { key: 'writing', label: '慧笔创作', icon: <IconEdit />, path: '/writing' },
      ]
    },
    {
      title: '管理中台',
      items: [
        { key: 'org', label: '组织管理', icon: <IconOrg />, path: '/admin/org' },
        { key: 'employee-mgmt', label: '员工管理', icon: <IconUser />, path: '/admin/employees' },
        { key: 'knowledge-mgmt', label: '知识管理', icon: <IconBook />, path: '/admin/knowledge' },
        { key: 'analytics', label: '数据看板', icon: <IconChart />, path: '/admin/analytics' },
      ]
    },
    {
      title: '开发后台',
      items: [
        { key: 'workspace', label: '开发工作台', icon: <IconCode />, path: '/dev/workspace' },
        { key: 'plugins', label: '插件开发', icon: <IconPlugin />, path: '/dev/plugins' },
        { key: 'workflows', label: '工作流', icon: <IconFlow />, path: '/dev/workflows' },
      ]
    }
  ]}
  user={{
    name: '张三',
    avatar: '/avatar.jpg',
    role: '管理员'
  }}
/>
```

#### 2.3.2 侧边栏样式规范

```css
/* =====================================================
   侧边栏样式
   ===================================================== */

.sidebar {
  /* 布局 */
  display: flex;
  flex-direction: column;
  width: 240px;
  height: calc(100vh - 64px);
  overflow-y: auto;

  /* 视觉 */
  background-color: var(--bg-primary);
  border-right: 1px solid var(--border-primary);

  /* 过渡动画 */
  transition: width var(--duration-normal) var(--ease-in-out);
}

/* 收起状态 */
.sidebar.collapsed {
  width: 72px;
}

/* 菜单分组 */
.menu-group {
  margin-bottom: var(--spacing-md);
}

.menu-group-title {
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 菜单项 */
.menu-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-md);
  margin: 0 var(--spacing-sm);
  color: var(--text-secondary);
  text-decoration: none;
  border-radius: var(--radius-sm);
  transition: all var(--duration-fast) var(--ease-in-out);
  cursor: pointer;
}

.menu-item:hover {
  background-color: var(--bg-secondary);
  color: var(--text-primary);
}

.menu-item.active {
  background-color: var(--primary-color);
  color: var(--text-inverse);
}

.menu-item-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.menu-item-label {
  flex: 1;
  font-size: var(--font-size-sm);
}

/* 收起状态下的菜单项 */
.sidebar.collapsed .menu-item {
  justify-content: center;
  padding: var(--spacing-sm);
}

.sidebar.collapsed .menu-item-label {
  display: none;
}

/* 子菜单 */
.submenu {
  padding-left: var(--spacing-lg);
}

.submenu-item {
  padding: var(--spacing-xs) var(--spacing-md);
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

/* 用户信息区域 */
.sidebar-user {
  padding: var(--spacing-md);
  border-top: 1px solid var(--border-primary);
  margin-top: auto;
}

.sidebar-user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.sidebar-user-avatar {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
}

.sidebar-user-name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.sidebar-user-role {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
}
```

### 2.4 主内容区设计

```css
/* =====================================================
   主内容区样式
   ===================================================== */

.main-content {
  /* 布局 */
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  /* 背景 */
  background-color: var(--bg-secondary);
}

/* 内容容器 */
.content-container {
  flex: 1;
  overflow-y: auto;
  padding: var(--layout-padding-md);
}

/* 页面头部 */
.page-header {
  margin-bottom: var(--spacing-lg);
}

.page-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.page-subtitle {
  font-size: var(--font-size-base);
  color: var(--text-secondary);
}

/* 页面操作栏 */
.page-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-md);
}

/* 内容卡片 */
.content-card {
  background-color: var(--bg-primary);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-3xl);
  text-align: center;
}

.empty-state-icon {
  font-size: 64px;
  color: var(--text-quaternary);
  margin-bottom: var(--spacing-md);
}

.empty-state-text {
  font-size: var(--font-size-lg);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
}

.empty-state-action {
  margin-top: var(--spacing-md);
}
```

---

## 三、组件设计规范

### 3.1 按钮组件

#### 3.1.1 按钮类型

```typescript
/**
 * 按钮组件规范
 */

type ButtonType =
  | 'primary'     // 主要按钮 - 品牌色背景
  | 'secondary'   // 次要按钮 - 灰色背景
  | 'outline'     // 轮廓按钮 - 透明背景+边框
  | 'text'        // 文本按钮 - 无背景无边框
  | 'danger';     // 危险按钮 - 红色

type ButtonSize =
  | 'small'       // 小号 - 高度28px
  | 'medium'      // 中号 - 高度32px
  | 'large';      // 大号 - 高度40px

interface ButtonProps {
  type?: ButtonType;
  size?: ButtonSize;
  disabled?: boolean;
  loading?: boolean;
  icon?: React.ReactNode;
  iconPosition?: 'left' | 'right';
  block?: boolean;            // 块级按钮（100%宽度）
  danger?: boolean;
  onClick?: () => void;
  children: React.ReactNode;
}

// 使用示例
<Button type="primary" size="medium" loading={false}>
  创建数字员工
</Button>

<Button type="outline" size="medium" icon={<IconPlus />}>
  添加成员
</Button>

<Button type="text" size="small" danger>
  删除
</Button>
```

#### 3.1.2 按钮样式

```css
/* =====================================================
   按钮样式规范
   ===================================================== */

.btn {
  /* 布局 */
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-xs);

  /* 尺寸 */
  height: 32px;
  padding: 0 var(--spacing-md);

  /* 字体 */
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);

  /* 边框 */
  border: 1px solid transparent;
  border-radius: var(--radius-sm);

  /* 过渡 */
  transition: all var(--duration-fast) var(--ease-in-out);

  /* 光标 */
  cursor: pointer;
  user-select: none;
}

/* 按钮尺寸 */
.btn.small {
  height: 28px;
  padding: 0 var(--spacing-sm);
  font-size: var(--font-size-xs);
}

.btn.large {
  height: 40px;
  padding: 0 var(--spacing-lg);
  font-size: var(--font-size-base);
}

/* 按钮类型 */
.btn.primary {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: var(--text-inverse);
}

.btn.primary:hover {
  background-color: var(--primary-hover);
  border-color: var(--primary-hover);
}

.btn.primary:active {
  background-color: var(--primary-active);
  border-color: var(--primary-active);
}

.btn.secondary {
  background-color: var(--bg-secondary);
  border-color: var(--border-primary);
  color: var(--text-primary);
}

.btn.secondary:hover {
  background-color: var(--bg-tertiary);
}

.btn.outline {
  background-color: transparent;
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.btn.outline:hover {
  background-color: var(--primary-color);
  color: var(--text-inverse);
}

.btn.text {
  background-color: transparent;
  border-color: transparent;
  color: var(--primary-color);
}

.btn.text:hover {
  background-color: var(--bg-secondary);
}

.btn.danger {
  background-color: var(--error-color);
  border-color: var(--error-color);
  color: var(--text-inverse);
}

/* 按钮状态 */
.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn.loading {
  position: relative;
  color: transparent;
  pointer-events: none;
}

.btn.loading::after {
  content: '';
  position: absolute;
  width: 16px;
  height: 16px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 块级按钮 */
.btn.block {
  display: flex;
  width: 100%;
}
```

### 3.2 表单组件

#### 3.2.1 输入框

```typescript
/**
 * 输入框组件规范
 */

interface InputProps {
  placeholder?: string;
  value?: string;
  defaultValue?: string;
  disabled?: boolean;
  error?: boolean;
  size?: 'small' | 'medium' | 'large';
  prefix?: React.ReactNode;    // 前缀图标
  suffix?: React.ReactNode;    // 后缀图标
  clearable?: boolean;         // 可清除
  maxLength?: number;
  onChange?: (value: string) => void;
  onPressEnter?: () => void;
}

// 使用示例
<Input
  placeholder="请输入团队名称"
  prefix={<IconSearch />}
  clearable
  maxLength={100}
  onChange={(value) => console.log(value)}
/>
```

```css
/* =====================================================
   输入框样式
   ===================================================== */

.input-wrapper {
  position: relative;
  display: inline-block;
  width: 100%;
}

.input {
  width: 100%;
  height: 32px;
  padding: 0 var(--spacing-sm);

  font-size: var(--font-size-sm);
  color: var(--text-primary);

  background-color: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);

  transition: all var(--duration-fast) var(--ease-in-out);
}

.input:hover {
  border-color: var(--border-secondary);
}

.input:focus {
  outline: none;
  border-color: var(--border-focus);
  box-shadow: 0 0 0 3px rgba(48, 110, 239, 0.1);
}

.input:disabled {
  background-color: var(--bg-secondary);
  color: var(--text-quaternary);
  cursor: not-allowed;
}

.input.error {
  border-color: var(--error-color);
}

.input.error:focus {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

/* 前缀/后缀 */
.input-prefix,
.input-suffix {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  color: var(--text-tertiary);
}

.input-prefix {
  left: var(--spacing-sm);
}

.input-suffix {
  right: var(--spacing-sm);
}

.input-with-prefix {
  padding-left: 32px;
}

.input-with-suffix {
  padding-right: 32px;
}
```

#### 3.2.2 下拉选择器

```typescript
/**
 * 下拉选择器组件规范
 */

interface SelectProps {
  placeholder?: string;
  value?: string | string[];
  options: SelectOption[];
  disabled?: boolean;
  error?: boolean;
  multiple?: boolean;
  clearable?: boolean;
  searchable?: boolean;
  loading?: boolean;
  onChange?: (value: string | string[]) => void;
}

interface SelectOption {
  label: string;
  value: string;
  disabled?: boolean;
  icon?: React.ReactNode;
}

// 使用示例
<Select
  placeholder="请选择角色"
  options={[
    { label: '所有者', value: 'owner' },
    { label: '管理员', value: 'admin' },
    { label: '成员', value: 'member' },
  ]}
  clearable
  onChange={(value) => console.log(value)}
/>
```

### 3.3 卡片组件

```typescript
/**
 * 卡片组件规范
 */

interface CardProps {
  title?: string;
  extra?: React.ReactNode;     // 额外操作（如"更多"链接）
  cover?: string;              // 封面图片
  actions?: React.ReactNode[]; // 操作按钮区
  hoverable?: boolean;         // 悬停效果
  bordered?: boolean;          // 是否显示边框
  shadow?: boolean;            // 是否显示阴影
  children: React.ReactNode;
}

// 使用示例
<Card
  title="数字员工统计"
  extra={<a href="#">更多</a>}
  hoverable
  shadow
>
  <p>卡片内容</p>
</Card>
```

```css
/* =====================================================
   卡片样式
   ===================================================== */

.card {
  background-color: var(--bg-primary);
  border-radius: var(--radius-md);
  overflow: hidden;
  transition: all var(--duration-base) var(--ease-in-out);
}

.card.bordered {
  border: 1px solid var(--border-primary);
}

.card.shadow {
  box-shadow: var(--shadow-md);
}

.card.hoverable:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}

/* 卡片头部 */
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--border-primary);
}

.card-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.card-extra {
  color: var(--text-secondary);
}

/* 卡片内容 */
.card-body {
  padding: var(--spacing-lg);
}

/* 卡片封面 */
.card-cover {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

/* 卡片操作区 */
.card-actions {
  display: flex;
  gap: var(--spacing-sm);
  padding: var(--spacing-md) var(--spacing-lg);
  border-top: 1px solid var(--border-primary);
}

.card-action-item {
  flex: 1;
  text-align: center;
  color: var(--text-secondary);
  cursor: pointer;
  transition: color var(--duration-fast) var(--ease-in-out);
}

.card-action-item:hover {
  color: var(--primary-color);
}
```

### 3.4 模态框组件

```typescript
/**
 * 模态框组件规范
 */

interface ModalProps {
  visible: boolean;
  title?: string;
  width?: number | string;
  centered?: boolean;         // 垂直居中
  closable?: boolean;          // 是否显示关闭按钮
  maskClosable?: boolean;      // 点击遮罩关闭
  footer?: React.ReactNode;   // 底部内容
  onCancel?: () => void;
  onOk?: () => void;
  children: React.ReactNode;
}

// 使用示例
<Modal
  visible={visible}
  title="创建数字员工"
  width={600}
  onCancel={() => setVisible(false)}
  onOk={() => handleSubmit()}
>
  <Form>...</Form>
</Modal>
```

```css
/* =====================================================
   模态框样式
   ===================================================== */

.modal-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--bg-overlay);
  z-index: 1000;
  animation: fadeIn var(--duration-base) var(--ease-in-out);
}

.modal-container {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: var(--bg-primary);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
  z-index: 1001;
  max-height: calc(100vh - 32px);
  display: flex;
  flex-direction: column;
  animation: slideInUp var(--duration-normal) var(--ease-out);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--border-primary);
}

.modal-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.modal-close {
  font-size: 20px;
  color: var(--text-tertiary);
  cursor: pointer;
  transition: color var(--duration-fast) var(--ease-in-out);
}

.modal-close:hover {
  color: var(--text-primary);
}

.modal-body {
  padding: var(--spacing-lg);
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
  padding: var(--spacing-md) var(--spacing-lg);
  border-top: 1px solid var(--border-primary);
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translate(-50%, -48%);
  }
  to {
    opacity: 1;
    transform: translate(-50%, -50%);
  }
}
```

### 3.5 列表与表格

#### 3.5.1 表格组件

```typescript
/**
 * 表格组件规范
 */

interface TableColumn {
  title: string;
  dataIndex: string;
  key: string;
  width?: number;
  align?: 'left' | 'center' | 'right';
  fixed?: 'left' | 'right';
  sorter?: boolean;
  render?: (value: any, record: any) => React.ReactNode;
}

interface TableProps {
  columns: TableColumn[];
  dataSource: any[];
  rowKey: string;
  loading?: boolean;
  pagination?: PaginationConfig;
  onRow?: (record: any) => {
    onClick?: (record: any) => void;
  };
}

// 使用示例
<Table
  columns={[
    { title: '姓名', dataIndex: 'name', key: 'name' },
    { title: '角色', dataIndex: 'role', key: 'role' },
    { title: '操作', key: 'action', render: (_, record) => (
      <Button size="small">编辑</Button>
    )},
  ]}
  dataSource={users}
  rowKey="id"
  pagination={{
    total: 100,
    pageSize: 20,
    current: 1,
  }}
/>
```

```css
/* =====================================================
   表格样式
   ===================================================== */

.table {
  width: 100%;
  border-collapse: collapse;
  background-color: var(--bg-primary);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.table thead {
  background-color: var(--bg-secondary);
}

.table th {
  padding: var(--spacing-sm) var(--spacing-md);
  text-align: left;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-primary);
}

.table td {
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-primary);
}

.table tbody tr:hover {
  background-color: var(--bg-secondary);
}

.table tbody tr.selected {
  background-color: rgba(48, 110, 239, 0.05);
}
```

---

## 四、三大前台页面设计

### 4.1 员工使用前台设计

#### 4.1.1 超级对话页面（Super Chatbox）

```
┌─────────────────────────────────────────────────────────────┐
│  [Logo] [对话] [数字员工] [智能问数] [创作] [搜索] [用户]    │
├─────────────────────────────────────────────────────────────┤
│  侧边栏                │      主对话区        │    数字员工列表  │
│                       │                      │                  │
│  历史对话              │  ┌─────────────────┐  │  [头像] QA助手   │
│  ├ 今天               │  │                 │  │  [头像] 数据分析师│
│  ├ 昨天               │  │  AI: 您好！      │  │  [头像] 智能客服  │
│  ├ 上周               │  │  我是您的AI助理 │  │                  │
│  └ 更早               │  │  有什么可以帮您?│  │  技能推荐         │
│                       │  │                 │  │  • 查询数据      │
│                       │  └─────────────────┘  │  • 生成图表      │
│                       │                      │  • 撰写文档      │
│                       │  [消息列表区域]       │                  │
│                       │                      │                  │
│                       │  ┌─────────────────┐  │                  │
│                       │  │ 输入框          │  │                  │
│                       │  │ [📎][📷][🎤][发送]│  │                  │
│                       │  └─────────────────┘  │                  │
└─────────────────────────────────────────────────────────────┘

页面布局：
- 三栏布局：历史对话（左，可收起）+ 主对话区（中）+ 数字员工列表（右，可收起）
- 默认宽度：左侧240px，中间flex，右侧280px
- 响应式：移动端单栏，可滑动切换
```

**关键交互设计**：

```typescript
/**
 * 超级对话页面交互设计
 */

interface ChatPageDesign {
  // 1. 消息气泡设计
  messageBubble: {
    user: {
      background: 'var(--primary-color)',
      color: 'var(--text-inverse)',
      borderRadius: '16px 16px 4px 16px',  // 右下角直角
      maxWidth: '60%',
      float: 'right',
    },
    ai: {
      background: 'var(--bg-secondary)',
      color: 'var(--text-primary)',
      borderRadius: '16px 16px 16px 4px',  // 左下角直角
      maxWidth: '70%',
      float: 'left',
    }
  },

  // 2. 输入框设计
  inputArea: {
    features: [
      '多模态输入（文字/语音/图片/文件）',
      '智能提示（自动补全/快捷指令）',
      '快捷操作（@员工/#技能/上传文件）',
      '字数统计（限制2000字符）',
      '发送快捷键（Enter发送，Shift+Enter换行）',
    ],
    height: 'auto',  // 自适应高度（最多5行）
    minHeight: '44px',
  },

  // 3. 思考动画
  thinkingAnimation: {
    type: 'typing-indicator',  // 打字机效果
    duration: 3000,  // 最多显示3秒
    dots: 3,  // 3个跳动的点
  },

  // 4. 流式输出
  streamingOutput: {
    enabled: true,
    speed: 20,  // 每秒20个字符
    typingEffect: true,  // 打字机效果
  },
}
```

**消息气泡样式示例**：

```css
/* =====================================================
   消息气泡样式
   ===================================================== */

.message-container {
  display: flex;
  margin-bottom: var(--spacing-md);
}

.message-user {
  justify-content: flex-end;
}

.message-ai {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 70%;
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-base);
  line-height: 1.5;
  word-wrap: break-word;
}

.message-user .message-bubble {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-hover));
  color: var(--text-inverse);
  border-radius: 16px 16px 4px 16px;
}

.message-ai .message-bubble {
  background-color: var(--bg-secondary);
  color: var(--text-primary);
  border-radius: 16px 16px 16px 4px;
  box-shadow: var(--shadow-sm);
}

/* 消息元信息 */
.message-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  margin-top: var(--spacing-xs);
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-full);
}

.message-time {
  font-size: 12px;
  color: var(--text-quaternary);
}

/* 流式输出光标 */
.streaming-cursor::after {
  content: '|';
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  50% {
    opacity: 0;
  }
}
```

#### 4.1.2 数字员工广场页面

```
┌─────────────────────────────────────────────────────────────┐
│  [Logo] [对话] [数字员工] [智能问数] [创作] [搜索] [用户]    │
├─────────────────────────────────────────────────────────────┤
│  页面标题：数字员工广场                                       │
│  [搜索框] [筛选：全部/问答型/操作型/综合型]                 │
├─────────────────────────────────────────────────────────────┤
│  卡片网格布局（3列）                                          │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐              │
│  │ [头像]      │ │ [头像]      │ │ [头像]      │              │
│  │ 客服助手    │ │ 数据分析师  │ │ 智能HR      │              │
│  │ 问答型      │ │ 操作型      │ │ 综合型      │              │
│  │ ⭐ 4.9     │ │ ⭐ 4.8     │ │ ⭐ 4.7     │              │
│  │ 1.2k次使用  │ │ 856次使用   │ │ 643次使用   │              │
│  │ [开始对话]  │ │ [开始对话]  │ │ [开始对话]  │              │
│  └────────────┘ └────────────┘ └────────────┘              │
│  ... 更多卡片                                               │
└─────────────────────────────────────────────────────────────┘

布局说明：
- 搜索栏：固定在页面顶部
- 筛选标签：水平滚动（移动端）
- 卡片网格：响应式3列（桌面）/ 2列（平板）/ 1列（手机）
- 卡片排序：按使用频率、评分、最近使用
```

**数字员工卡片设计**：

```typescript
/**
 * 数字员工卡片组件
 */

interface EmployeeCardProps {
  employee: {
    id: string;
    name: string;
    type: 'qa' | 'operation' | 'comprehensive';
    avatar: string;
    description: string;
    rating: number;
    usageCount: number;
    capabilities: string[];  // 能力标签
  };
  onClick: () => void;
}

// 卡片样式
const cardStyle = {
  container: {
    background: 'var(--bg-primary)',
    borderRadius: 'var(--radius-md)',
    padding: 'var(--spacing-lg)',
    boxShadow: 'var(--shadow-sm)',
    transition: 'all var(--duration-base) var(--ease-in-out)',
    cursor: 'pointer',
    hover: {
      boxShadow: 'var(--shadow-lg)',
      transform: 'translateY(-4px)',
    }
  },
  avatar: {
    width: '64px',
    height: '64px',
    borderRadius: 'var(--radius-full)',
    border: '3px solid var(--border-primary)',
  },
  typeBadge: {
    display: 'inline-block',
    padding: '2px 8px',
    borderRadius: 'var(--radius-full)',
    fontSize: 'var(--font-size-xs)',
    fontWeight: 'var(--font-weight-medium)',
    // 类型颜色
    qa: { background: 'var(--tag-blue)', color: 'white' },
    operation: { background: 'var(--tag-orange)', color: 'white' },
    comprehensive: { background: 'var(--tag-purple)', color: 'white' },
  }
};
```

#### 4.1.3 智能问数页面（ChatBI）

```
┌─────────────────────────────────────────────────────────────┐
│  页面标题：智能问数                                           │
│  [示例问题卡片]                                             │
│  • "上个月哪个产品销售额最高？"                               │
│  • "对比今年和去年同期的增长趋势"                             │
│  • "华东区各门店销售排行"                                     │
├─────────────────────────────────────────────────────────────┤
│  输入区域                                                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 请输入您的数据查询问题...                          │   │
│  │ [📎 上传数据源] [🎤 语音输入] [发送]                │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│  结果展示区                                                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  图表区域                    │  智能解读文字        │   │
│  │  ┌────────────────────┐     │                      │   │
│  │  │  [柱状图]           │     │  根据查询结果，     │   │
│  │  │  [折线图]           │     │  产品A的销售额...     │   │
│  │  │  [数据表格]         │     │  同比增长20%...      │   │
│  │  └────────────────────┘     │                      │   │
│  │                             │  [导出] [分享]       │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘

关键设计点：
1. 示例问题卡片 - 快速引导用户提问
2. 多维度结果展示 - 图表 + 表格 + 文字解读
3. 智能推荐 - 相关问题推荐
4. 历史查询 - 最近的查询记录
```

#### 4.1.4 慧笔创作页面

```
┌─────────────────────────────────────────────────────────────┐
│  页面标题：慧笔创作                                           │
│  [选择模板] [自定义创作]                                     │
├─────────────────────────────────────────────────────────────┤
│  模板选择区（水平滚动）                                      │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐          │
│  │周报  │ │月报  │ │会议纪要│ │工作计划│ │招聘JD│...     │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘          │
├─────────────────────────────────────────────────────────────┤
│  创作区域                                                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 标题：_____________________________                   │   │
│  │                                                     │   │
│  │ [富文本编辑器]                                        │   │
│  │                                                     │   │
│  │                                                     │   │
│  └─────────────────────────────────────────────────────┘   │
│  [优化建议] "可以添加数据支撑，让内容更有说服力"             │
│  [重新生成] [保存] [导出] [分享]                              │
└─────────────────────────────────────────────────────────────┘

创作流程：
1. 选择模板（或自定义）
2. 输入简要要求
3. AI生成初稿
4. 多轮优化（基于用户反馈）
5. 保存/导出/分享
```

---

### 4.2 企业管理中台设计

#### 4.2.1 组织管理页面

```
┌─────────────────────────────────────────────────────────────┐
│  [面包屑] 首页 / 管理中台 / 组织管理                            │
├─────────────────────────────────────────────────────────────┤
│  页面标题：组织管理                                           │
│  [创建组织] [批量导入]                                       │
├─────────────────────────────────────────────────────────────┤
│  左侧：组织树        │  右侧：成员列表                         │
│  ┌──────────────┐   │  ┌───────────────────────────────────┐ │
│  │ 🏢 公司      │   │  │ 搜索: _____________  [搜索]      │ │
│  │   ├ 技术部    │   │  │  筛选: [部门▼] [岗位▼]            │ │
│  │   │   ├ 后端  │   │  │                                   │ │
│  │   │   └ 前端  │   │  │  成员列表：                         │ │
│  │   ├ 产品部    │   │  │  ┌──────────────────────────┐  │ │
│  │   └ 运营部    │   │  │  │ [头像] 张三  产品经理  │  │ │
│  │               │   │  │  │ [编辑] [删除]           │  │ │
│  └──────────────┘   │  │  └──────────────────────────┘  │ │
│  [添加子组织]       │  │                                   │ │
│                    │  │  ... 更多成员                    │ │
│                    │  │                                   │ │
│                    │  │  [批量操作] [导出]               │ │
│                    │  └───────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

布局特点：
- 左右分栏：组织树（左，固定宽度240px）+ 成员列表（右，flex）
- 组织树：可展开/收起，支持拖拽调整
- 成员列表：表格展示，支持排序、筛选、批量操作
```

**组织架构树组件**：

```typescript
/**
 * 组织架构树组件
 */

interface OrgTreeProps {
  data: OrganizationNode[];
  selectedKey?: string;
  onSelect?: (node: OrganizationNode) => void;
  draggable?: boolean;        // 支持拖拽
  expandable?: boolean;      // 可展开/收起
}

interface OrganizationNode {
  id: string;
  name: string;
  type: 'company' | 'department' | 'team';
  parentId: string | null;
  children?: OrganizationNode[];
  memberCount: number;
}

// 节点样式
const treeNodeStyle = {
  node: {
    display: 'flex',
    alignItems: 'center',
    gap: 'var(--spacing-sm)',
    padding: 'var(--spacing-xs) var(--spacing-sm)',
    borderRadius: 'var(--radius-sm)',
    cursor: 'pointer',
    '&:hover': {
      backgroundColor: 'var(--bg-secondary)',
    }
  },
  icon: {
    color: 'var(--text-tertiary)',
  },
  label: {
    fontSize: 'var(--font-size-sm)',
    color: 'var(--text-primary)',
  },
  badge: {
    marginLeft: 'auto',
    fontSize: 'var(--font-size-xs)',
    color: 'var(--text-quaternary)',
  }
};
```

#### 4.2.2 数字员工管理页面

```
┌─────────────────────────────────────────────────────────────┐
│  面包屑 / 管理中台 / 数字员工管理                              │
├─────────────────────────────────────────────────────────────┤
│  Tab: [已发布(15)] [草稿(3)] [已下线(2)]                      │
├─────────────────────────────────────────────────────────────┤
│  工具栏：[创建员工] [批量操作▼] [导出]                         │
├─────────────────────────────────────────────────────────────┤
│  员工列表（表格）                                             │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ 名称   │ 类型 │ 状态 │ 使用次数 │ 满意度 │ 操作      │  │
│  ├───────────────────────────────────────────────────────┤  │
│  │ 客服   │ 问答  │ 已发布│ 1.2k    │ 4.8⭐  │ [查看]   │  │
│  │ 助手   │      │      │         │        │ [编辑]   │  │
│  │        │      │      │         │        │ [下线]   │  │
│  ├───────────────────────────────────────────────────────┤  │
│  │ 数据   │ 操作  │ 已发布│ 856     │ 4.6⭐  │ [查看]   │  │
│  │ 分析师 │      │      │         │        │ [编辑]   │  │
│  └───────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│  <分页器>                                                    │
└─────────────────────────────────────────────────────────────┘

关键功能：
1. Tab切换：已发布 / 草稿 / 已下线
2. 表格操作：查看详情 / 编辑 / 下线 / 删除
3. 状态管理：发布 / 停用 / 下线
4. 使用统计：对话次数、满意度、响应时间
```

#### 4.2.3 数据看板页面

```
┌─────────────────────────────────────────────────────────────┐
│  面包屑 / 管理中台 / 数据看板                                   │
├─────────────────────────────────────────────────────────────┤
│  筛选：[时间范围▼] [组织▼] [数字员工▼]                      │
│  [导出报表] [刷新]                                           │
├─────────────────────────────────────────────────────────────┤
│  核心指标卡片（4列）                                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │ 总用户数  │ │ 活跃用户 │ │ 对话次数 │ │ 满意度   │       │
│  │  1,234   │ │   856    │ │  12.5k   │ │  4.7⭐   │       │
│  │  +12% ↑  │ │  +8% ↑   │ │  +15% ↑  │ │  +0.2↑   │       │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │
├─────────────────────────────────────────────────────────────┤
│  趋势图表（2列）                                             │
│  ┌─────────────────────────────┐ ┌─────────────────────────┐ │
│  │ 用户增长趋势（折线图）        │ │ 使用时长分布（饼图）    │ │
│  │                             │ │                         │ │
│  │     [图表]                   │ │      [图表]            │ │
│  │                             │ │                         │ │
│  └─────────────────────────────┘ └─────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│  数字员工排行榜（表格）                                     │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ 排名 │ 员工名称 │ 类型 │ 使用次数 │ 满意度 │ 趋势    │  │
│  ├───────────────────────────────────────────────────────┤  │
│  │  1   │ 客服助手 │ 问答 │  3.2k   │ 4.9⭐  │ ↑ 15%  │  │
│  │  2   │ 数据分析 │ 操作 │  1.8k   │ 4.7⭐  │ ↑ 8%   │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

设计要点：
1. 核心指标卡片：突出数字 + 趋势箭头 + 百分比
2. 图表联动：点击图表可钻取详细数据
3. 自动刷新：支持设置自动刷新间隔（5分钟/15分钟/30分钟）
4. 导出功能：支持导出Excel/PDF/PNG格式
```

---

### 4.3 开发扩展后台设计

#### 4.3.1 开发工作台页面

```
┌─────────────────────────────────────────────────────────────┐
│  面包屑 / 开发后台 / 工作台                                   │
├─────────────────────────────────────────────────────────────┤
│  左侧：资源导航     │      主开发区                         │
│  ┌──────────────┐   │  ┌─────────────────────────────────┐ │
│  │ 🤖 智能体     │   │  │ Tab: [编辑] [调试] [版本] [设置]  │ │
│  │   ├ 我的bot  │   │  ├─────────────────────────────────┤ │
│  │   ├ 公共bot  │   │  │  编辑器区域                      │ │
│  │   └ 草稿箱    │   │  │  ┌────────────────────────┐    │ │
│  │               │   │  │  │  [提示词编辑器]       │    │ │
│  │ ⚙ 工作流     │   │  │  │                        │    │ │
│  │   ├ 我的流程 │   │  │  │  [工作流可视化编排]    │    │ │
│  │   └ 流程模板 │   │  │  │                        │    │ │
│  │               │   │  │  └────────────────────────┘    │ │
│  │ 📚 知识库     │   │  │                                 │ │
│  │   ├ 我的KB   │   │  │  右侧：属性面板                  │ │
│  │   ├ 公共KB   │   │  │  ┌──────────────────────────┐ │ │
│  │   └ 数据集   │   │  │  │ 模型配置                │ │ │
│  │               │   │  │  │ 知识库绑定              │ │ │
│  │ 🔌 插件      │   │  │  │ 技能绑定                │ │ │
│  │   ├ 我的插件 │   │  │  │ 发布设置                │ │ │
│  │   └ 插件商店 │   │  │  └──────────────────────────┘ │ │
│  │               │   │  │                                 │ │
│  └──────────────┘   │  └─────────────────────────────────┘ │
│                    │  [保存] [调试] [发布] [预览]             │
└─────────────────────────────────────────────────────────────┘

布局特点：
- 三栏布局：资源导航（左）+ 编辑区（中）+ 属性面板（右）
- 编辑区：可切换视图（代码视图 / 可视化视图）
- 属性面板：折叠/展开，支持吸附定位
```

**智能体编辑器设计**：

```typescript
/**
 * 智能体编辑器布局
 */

interface BotEditorLayout {
  // 顶部工具栏
  toolbar: {
    left: [
      '保存', '撤销', '重做', '复制', '粘贴', '删除'
    ],
    center: [
      '代码视图', '可视化视图', '分屏视图'
    ],
    right: [
      '调试', '预览', '发布'
    ]
  },

  // 编辑区域
  editor: {
    // 代码视图
    codeView: {
      language: 'markdown',  // Prompt支持Markdown
      monacoEditor: true,     // 使用Monaco Editor
      lineNumbers: true,
      minimap: true,
      fontSize: 14,
    },

    // 可视化视图（工作流编排）
    visualView: {
      canvas: {
        gridSize: 20,         // 网格大小
        snapToGrid: true,     // 吸附到网格
        zoom: true,           // 支持缩放
        pan: true,            // 支持平移
      },
      nodes: {
        start: { icon: '▶', color: '#10B981' },
        llm: { icon: '🤖', color: '#3B82F6' },
        code: { icon: '💻', color: '#8B5CF6' },
        http: { icon: '🌐', color: '#F59E0B' },
        end: { icon: '⏹', color: '#EF4444' },
      },
      connections: {
        animated: true,       // 连线动画
        smooth: true,          // 平滑曲线
      },
    },
  },

  // 属性面板
  properties: {
    sections: [
      {
        title: '基础信息',
        fields: ['名称', '描述', '头像', '分类'],
      },
      {
        title: '模型配置',
        fields: ['模型', '温度', '最大Token数', '频率惩罚'],
      },
      {
        title: '知识库',
        fields: ['选择知识库', '检索策略', '相似度阈值'],
      },
      {
        title: '技能/插件',
        fields: ['选择插件', '配置参数'],
      },
      {
        title: '发布设置',
        fields: ['发布目录', '授权配置', '版本说明'],
      },
    ],
  },
}
```

**工作流可视化编辑器**：

```typescript
/**
 * 工作流节点设计
 */

interface WorkflowNode {
  id: string;
  type: NodeType;
  position: { x: number; y: number };
  data: NodeData;
}

type NodeType =
  | 'start'        // 开始节点
  | 'end'          // 结束节点
  | 'llm'          // 大模型节点
  | 'code'         // 代码节点
  | 'http'         // HTTP请求节点
  | 'database'     // 数据库节点
  | 'condition'    // 条件分支节点
  | 'loop'         // 循环节点
  | 'parallel'     // 并行节点
  | 'merge';       // 合并节点

// 节点样式
const nodeStyle = {
  container: {
    width: 180,
    padding: 'var(--spacing-md)',
    borderRadius: 'var(--radius-md)',
    boxShadow: 'var(--shadow-md)',
    border: '2px solid var(--border-primary)',
    background: 'var(--bg-primary)',
    cursor: 'move',
    transition: 'all var(--duration-fast) var(--ease-in-out)',
  },
  selected: {
    borderColor: 'var(--primary-color)',
    boxShadow: '0 0 0 3px rgba(48, 110, 239, 0.1)',
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    gap: 'var(--spacing-xs)',
    marginBottom: 'var(--spacing-sm)',
    fontSize: 'var(--font-size-sm)',
    fontWeight: 'var(--font-weight-semibold)',
  },
  icon: {
    width: 20,
    height: 20,
  },
  body: {
    fontSize: 'var(--font-size-xs)',
    color: 'var(--text-secondary)',
  },
  ports: {
    top: { type: 'output', position: 'top' },
    right: { type: 'output', position: 'right' },
    bottom: { type: 'input', position: 'bottom' },
    left: { type: 'input', position: 'left' },
  },
};
```

---

## 五、交互流程设计

### 5.1 核心交互流程

#### 5.1.1 用户登录流程

```
┌────────────────┐
│  登录页面        │
├────────────────┤
│  [Logo]         │
│  [欢迎文案]     │
│  ┌────────────┐ │
│  │邮箱/手机号  │ │
│  ├────────────┤ │
│  │  密码       │ │
│  ├────────────┤ │
│  │  [登录]     │ │
│  └────────────┘ │
│  [忘记密码?]    │
│  [注册账号]     │
└────────────────┘
         ↓
┌────────────────┐
│  加载动画       │
│  [spinner]      │
│  "登录中..."    │
└────────────────┘
         ↓
┌────────────────┐
│  主页面         │
│  [顶部导航]     │
│  [侧边栏]       │
│  [主内容区]     │
└────────────────┘

交互细节：
1. 表单验证：实时验证邮箱格式、密码强度
2. 错误提示：输入框下方显示错误信息
3. 加载状态：登录按钮显示loading动画
4. 记住密码：可选，7天内免登录
5. 自动跳转：登录成功后自动跳转到之前访问的页面
```

#### 5.1.2 创建数字员工流程

```
步骤1：选择类型
┌────────────────────────────────────────────────────┐
│  创建数字员工                                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐           │
│  │ 问答型   │ │ 操作型   │ │ 综合型   │           │
│  │  💬      │ │  ⚙️      │ │  🤖      │           │
│  │  擺知识库 │ │  用技能   │ │ 问答+操作 │           │
│  └──────────┘ └──────────┘ └──────────┘           │
└────────────────────────────────────────────────────┘
         ↓ [下一步]
步骤2：基础信息
┌────────────────────────────────────────────────────┐
│  员工信息                                           │
│  ┌──────────────────────────────────────────┐       │
│  │ 名称: [_________________________]        │       │
│  ├──────────────────────────────────────────┤       │
│  │ 描述: [_________________________]        │       │
│  │        [_________________________]        │       │
│  ├──────────────────────────────────────────┤       │
│  │ 头像: [上传头像]                             │       │
│  └──────────────────────────────────────────┘       │
│  [上一步] [下一步]                                  │
└────────────────────────────────────────────────────┘
         ↓ [下一步]
步骤3：能力配置（根据类型显示不同配置项）
┌────────────────────────────────────────────────────┐
│  能力配置（问答型）                                 │
│  ┌──────────────────────────────────────────┐       │
│  │ 知识库绑定                                   │       │
│  │  [✓] 产品手册 (123篇文档)                  │       │
│  │  [✓] FAQ知识库 (56条问答)                  │       │
│  │  [+ 添加知识库]                             │       │
│  ├──────────────────────────────────────────┤       │
│  │ 模型配置                                     │       │
│  │  模型: [GPT-4 ▼] 温度: [0.7]             │       │
│  ├──────────────────────────────────────────┤       │
│  │ 开场白                                       │       │
│  │  [________________________________]        │       │
│  │  [0/200]                                     │       │
│  ├──────────────────────────────────────────┤       │
│  │ 引导问题                                     │       │
│  │  1. [________________________________]        │       │
│  │  2. [________________________________]        │       │
│  │  3. [________________________________]        │       │
│  │  [+ 添加问题]                                 │       │
│  └──────────────────────────────────────────┘       │
│  [上一步] [保存草稿] [完成]                          │
└────────────────────────────────────────────────────┘

流程特点：
1. 进度指示：顶部显示进度条（1/3, 2/3, 3/3）
2. 表单验证：实时验证，下一步按钮禁用直到必填项完成
3. 保存草稿：随时保存，可稍后继续
4. 完成预览：最后一步可预览效果再提交
```

#### 5.1.3 超级对话交互流程

```
┌─────────────────────────────────────────────────────────────┐
│                    超级对话完整交互流程                        │
└─────────────────────────────────────────────────────────────┘

步骤1：用户输入
┌────────────────────────────────────────────────────┐
│  输入框区域                                         │
│  ┌────────────────────────────────────────────┐   │
│  │ @数据分析助手 上个月销售额最高的产品是...  │   │
│  │ [📎] [📷] [🎤] [发送]                       │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓
步骤2：意图识别
┌────────────────────────────────────────────────────┐
│  系统处理（后台，用户不可见）                         │
│  • NLU意图识别：数据查询                            │
│  • 实体提取：时间=上个月，指标=销售额                │
│  • 技能路由：调用[ChatBI技能]                       │
└────────────────────────────────────────────────────┘
         ↓
步骤3：技能调用
┌────────────────────────────────────────────────────┐
│  [技能指示器]                                       │
│  🔍 正在查询数据...                                 │
│  ◐ 加载中（思考动画）                               │
└────────────────────────────────────────────────────┘
         ↓
步骤4：知识库检索（如果需要）
┌────────────────────────────────────────────────────┐
│  知识库匹配过程                                     │
│  ✓ 检索产品手册知识库（匹配度: 85%）               │
│  ✓ 检索销售数据知识库（匹配度: 92%）               │
│  ✓ 检索FAQ知识库（相关条目: 3条）                  │
└────────────────────────────────────────────────────┘
         ↓
步骤5：AI响应生成
┌────────────────────────────────────────────────────┐
│  消息气泡（AI）                                     │
│  ┌────────────────────────────────────────────┐   │
│  │ [头像] 数据分析师                           │   │
│  │ ├──────────────────────────────────────┤   │
│  │ │ 根据查询结果，上个月产品A的销售额最高，│   │
│  │ │ 达到2,345万元，占总销售额的35%。     │   │
│  │ │                                          │   │
│  │ │ 📊 [查看详细图表]                       │   │
│  │ │ 📄 [导出数据报告]                       │   │
│  │ └──────────────────────────────────────┘   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓
步骤6：用户追问
┌────────────────────────────────────────────────────┐
│  输入框区域                                         │
│  ┌────────────────────────────────────────────┐   │
│  │ 那同比增长多少？                             │   │
│  │ [📎] [📷] [🎤] [发送]                       │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓
步骤7：上下文理解 + 多轮对话
┌────────────────────────────────────────────────────┐
│  消息气泡（AI）                                     │
│  ┌────────────────────────────────────────────┐   │
│  │ [头像] 数据分析师                           │   │
│  │ ├──────────────────────────────────────┤   │
│  │ │ 相比去年同期，产品A销售额增长20%，    │   │
│  │ │ 主要得益于新产品线的推出和市场营销...  │   │
│  │ │                                          │   │
│  │ │ 💡 相关问题推荐：                        │   │
│  │ │ • 各地区的销售占比如何？                 │   │
│  │ │ • 下个月的预测趋势是什么？               │   │
│  │ └──────────────────────────────────────┘   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘

交互特点：
1. 实时打字效果：流式输出，每秒20字符
2. 思考动画：显示"正在思考..."
3. 技能切换可视化：显示调用的技能名称
4. 多轮对话记忆：保持上下文10轮
5. 快捷操作：消息气泡内嵌入操作按钮
6. 相关问题推荐：引导用户继续提问
```

**状态设计**：

```typescript
/**
 * 超级对话状态管理
 */
interface ChatState {
  // 1. 对话状态
  status: 'idle' | 'thinking' | 'streaming' | 'error';

  // 2. 消息类型
  messageType: 'text' | 'image' | 'file' | 'chart' | 'code';

  // 3. 技能调用状态
  skillInvocation: {
    skillName: string;
    status: 'calling' | 'processing' | 'completed' | 'failed';
    progress?: number;  // 0-100
  };

  // 4. 上下文管理
  context: {
    history: Message[];  // 最近10轮对话
    currentEmployee: string;  // 当前对话的员工ID
    knowledgeBaseMatched: string[];  // 匹配的知识库
  };
}
```

---

#### 5.1.4 智能问数（ChatBI）交互流程

```
┌─────────────────────────────────────────────────────────────┐
│                    智能问数完整交互流程                        │
└─────────────────────────────────────────────────────────────┘

步骤1：用户提问
┌────────────────────────────────────────────────────┐
│  页面：智能问数（ChatBI）                            │
│  ┌────────────────────────────────────────────┐   │
│  │ 示例问题：                                   │   │
│  │ • "上个月哪个产品销售额最高？"               │   │
│  │ • "对比今年和去年同期的增长趋势"             │   │
│  │ • "华东区各门店销售排行"                     │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  ┌────────────────────────────────────────────┐   │
│  │ 输入框：请输入您的数据查询问题...          │   │
│  │ [🎤 语音] [📊 数据源] [🔍 查询]            │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓ [用户输入：显示各季度销售趋势]
步骤2：SQL生成（后台）
┌────────────────────────────────────────────────────┐
│  自然语言处理（后台，用户不可见）                     │
│  • 意图识别：时间趋势分析                           │
│  • 实体提取：指标=销售额，维度=季度                 │
│  • SQL生成：                                       │
│    SELECT quarter, SUM(amount) as total            │
│    FROM sales                                     │
│    WHERE year = 2024                              │
│    GROUP BY quarter                               │
│    ORDER BY quarter                               │
└────────────────────────────────────────────────────┘
         ↓
步骤3：数据查询
┌────────────────────────────────────────────────────┐
│  [加载状态]                                         │
│  🔍 正在查询数据库...                               │
│  ◐ 连接数据源                                       │
│  ◐ 执行查询                                         │
│  ✓ 查询完成（234ms）                                │
└────────────────────────────────────────────────────┘
         ↓
步骤4：智能图表推荐
┌────────────────────────────────────────────────────┐
│  图表类型选择（AI自动推荐，用户可切换）               │
│  ┌────────┐ ┌────────┐ ┌────────┐                │
│  │ ✓ 折线图│ │  柱状图 │ │  饼图   │                │
│  │ (推荐)  │ │         │ │         │                │
│  └────────┘ └────────┘ └────────┘                │
│                                                     │
│  选择理由：                                         │
│  "时间趋势数据，折线图最能体现变化趋势"              │
└────────────────────────────────────────────────────┘
         ↓
步骤5：结果展示
┌────────────────────────────────────────────────────┐
│  结果展示区                                         │
│  ┌──────────────────────┬──────────────────────┐ │
│  │ 图表区域              │ 智能解读              │ │
│  │ ┌────────────────┐   │                      │ │
│  │ │  ┌───┐         │   │ 根据查询结果分析：   │ │
│  │ │  │   │ 120M    │   │ • Q4销售额最高，    │ │
│  │ │  └───┘         │   │   达到1.2亿元      │ │
│  │ │      ┌───┐     │   │ • 环比增长15%     │ │
│  │ │      │   │ 95M │   │ • 同比增长23%     │ │
│  │ │      └───┘     │   │                  │ │
│  │ │  Q1  Q2  Q3  Q4│   │ 💡 建议：         │ │
│  │ └────────────────┘   │   Q4增长明显，     │ │
│  │                      │   建议加大营销投入 │ │
│  │ [切换图表]           │                  │ │
│  │ [导出图表]           │ [复制解读]        │ │
│  └──────────────────────┴──────────────────────┘ │
└────────────────────────────────────────────────────┘
         ↓
步骤6：数据钻取
┌────────────────────────────────────────────────────┐
│  用户点击Q4数据柱                                   │
│  ┌────────────────────────────────────────────┐   │
│  │  Q4详细数据                                 │   │
│  │  • 10月: 3,800万                           │   │
│  │  • 11月: 4,200万                           │   │
│  │  • 12月: 4,000万                           │   │
│  │                                             │   │
│  │  [按月钻取] [按地区钻取] [按产品钻取]      │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓
步骤7：导出/分享
┌────────────────────────────────────────────────────┐
│  导出选项                                           │
│  ┌────────────────────────────────────────────┐   │
│  │ 导出格式：                                   │   │
│  │ [📊 Excel] [📈 PNG图片] [📄 PDF报告]        │   │
│  │                                             │   │
│  │ 分享：                                       │   │
│  │ [🔗 复制链接] [✉️ 邮件发送] [💬 生成报告]   │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘

关键交互设计：
1. 自然语言输入：支持语音输入、输入联想
2. 数据源选择：可切换不同数据库/数据集
3. 智能图表推荐：AI根据数据特征推荐最佳图表
4. 实时预览：输入时实时显示SQL预览（高级用户）
5. 多维度钻取：点击图表元素进行数据钻取
6. 智能解读：AI自动生成数据洞察和建议
7. 导出多样性：支持图表、数据、报告多种格式
```

---

#### 5.1.5 慧笔创作交互流程

```
┌─────────────────────────────────────────────────────────────┐
│                    慧笔创作完整交互流程                        │
└─────────────────────────────────────────────────────────────┘

步骤1：选择模板
┌────────────────────────────────────────────────────┐
│  页面：慧笔创作                                      │
│  ┌────────────────────────────────────────────┐   │
│  │ 选择模板（水平滚动）                         │   │
│  │ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐     │   │
│  │ │周报  │ │月报  │ │会议  │ │工作  │ ... │   │
│  │ │      │ │      │ │纪要  │ │计划  │     │   │
│  │ └──────┘ └──────┘ └──────┘ └──────┘     │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  或                                                  │
│  ┌────────────────────────────────────────────┐   │
│  │ [+] 自定义创作                              │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓ [选择：周报模板]
步骤2：输入需求
┌────────────────────────────────────────────────────┐
│  创作需求                                           │
│  ┌────────────────────────────────────────────┐   │
│  │ 标题：____________________________________  │   │
│  ├────────────────────────────────────────────┤   │
│  │ 内容要求：                                  │   │
│  │ ┌────────────────────────────────────────┐ │   │
│  │ │ 本周完成了产品A的上线工作，           │ │   │
│  │ │ 解决了5个bug，参加了3次会议。       │ │   │
│  │ │ 需要生成一份专业的周报。             │ │   │
│  │ │                                        │ │   │
│  │ └────────────────────────────────────────┘ │   │
│  │ 字数统计: 45/500                            │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  选项：                                             │
│  ☑️ 包含工作总结   ☑️ 包含下周计划                 │
│  ☐ 包含数据图表   ☐ 使用正式语气                   │
│                                                     │
│  [上一步] [开始生成]                                │
└────────────────────────────────────────────────────┘
         ↓
步骤3：AI生成初稿
┌────────────────────────────────────────────────────┐
│  [生成状态]                                         │
│  ✍️ 正在生成内容...                                 │
│  ◐ 分析需求（15%）                                  │
│  ◐ 整理信息（40%）                                  │
│  ◐ 撰写初稿（75%）                                  │
│  ✓ 生成完成（100%）                                 │
│                                                     │
│  用时: 2.3秒                                        │
└────────────────────────────────────────────────────┘
         ↓
步骤4：初稿预览与编辑
┌────────────────────────────────────────────────────┐
│  编辑器区域                                         │
│  ┌────────────────────────────────────────────┐   │
│  │ 产品研发周报 - 第48周                       │   │
│  │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │   │
│  │                                             │   │
│  │ 一、本周工作总结                            │   │
│  │ 1. 产品A上线完成                           │   │
│  │    - 完成前端开发，共15个页面             │   │
│  │    - 完成后端接口开发，共28个API          │   │
│  │    - 通过测试验收，Bug率<2%               │   │
│  │                                             │   │
│  │ 2. Bug修复                                 │   │
│  │    - 修复高优先级Bug 3个                  │   │
│  │    - 修复中优先级Bug 2个                  │   │
│  │                                             │   │
│  │ ... （完整内容）                            │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  操作栏：                                           │
│  [📝 继续编辑] [✨ 优化建议] [🔄 重新生成]          │
└────────────────────────────────────────────────────┘
         ↓ [点击"优化建议"]
步骤5：AI优化建议
┌────────────────────────────────────────────────────┐
│  优化建议（右侧面板）                               │
│  ┌────────────────────────────────────────────┐   │
│  │ 💡 优化建议（3条）                         │   │
│  │ ─────────────────────────────────────────  │   │
│  │ 1. 添加数据支撑                            │   │
│  │    当前版本缺少具体数据，建议添加：        │   │
│  │    • 上线后的用户反馈数据                  │   │
│  │    • 性能提升的具体百分比                  │   │
│  │    [✓ 一键应用]                            │   │
│  │                                             │   │
│  │ 2. 增加图表展示                            │   │
│  │    建议添加：                               │   │
│  │    • Bug修复趋势图                         │   │
│  │    • 工作时长分布图                        │   │
│  │    [✓ 一键应用]                            │   │
│  │                                             │   │
│  │ 3. 优化语言表达                            │   │
│  │    将"完成开发"改为"成功交付"，            │   │
│  │    体现专业性和成果导向。                  │   │
│  │    [✓ 一键应用]                            │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓ [应用优化1]
步骤6：多轮优化
┌────────────────────────────────────────────────────┐
│  编辑器区域（已应用优化）                           │
│  ┌────────────────────────────────────────────┐   │
│  │ ...                                         │   │
│  │ 1. 产品A上线完成                           │   │
│  │    - 完成前端开发，共15个页面             │   │
│  │    - 完成后端接口开发，共28个API          │   │
│  │    - 通过测试验收，用户满意度达92% ✨新增 │   │
│  │    - 系统性能提升35%                      │   │
│  │                                             │   │
│  │ [插入图表]                                  │   │
│  │ ┌──────────────────────────────┐          │   │
│  │ │   Bug修复趋势图（第48周）    │          │   │
│  │ │   5│ ████                    │          │   │
│  │ │   3│    ████                 │          │   │
│  │ │   2│       ████              │          │   │
│  │ │    └────────────────          │          │   │
│  │ │      周一 周三 周五           │          │   │
│  │ └──────────────────────────────┘          │   │
│  │                                             │   │
│  │ ...                                         │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  ✅ 已应用1条优化建议                               │
│  [继续优化] [对比编辑] [保存]                       │
└────────────────────────────────────────────────────┘
         ↓
步骤7：对比编辑
┌────────────────────────────────────────────────────┐
│  对比视图（左右分栏）                               │
│  ┌──────────────────┬──────────────────────────┐ │
│  │ 原版本            │ 优化版本                 │ │
│  │ ├────────────────┤├────────────────────────┤ │
│  │ "完成开发"       │ "成功交付" ✨            │ │
│  │                  │                          │ │
│  │ 缺少数据         │ 用户满意度: 92%          │ │
│  │                  │ 性能提升: 35%            │ │
│  └──────────────────┴──────────────────────────┘ │
│                                                     │
│  [全部接受优化] [逐条选择] [手动编辑]               │
└────────────────────────────────────────────────────┘
         ↓
步骤8：保存/导出
┌────────────────────────────────────────────────────┐
│  保存/导出选项                                     │
│  ┌────────────────────────────────────────────┐   │
│  │ 保存：                                       │   │
│  │ 💾 保存到草稿箱                              │   │
│  │ ⭐ 收藏为模板                                 │   │
│  │                                             │   │
│  │ 导出：                                       │   │
│  │ 📄 导出为Word                                │   │
│  │ 📊 导出为PDF                                 │   │
│  │ 📝 复制纯文本                                │   │
│  │                                             │   │
│  │ 分享：                                       │   │
│  │ 🔗 生成分享链接                              │   │
│  │ ✉️ 邮件发送                                  │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘

关键交互设计：
1. 模板驱动：提供常用文档模板，快速开始
2. 结构化输入：引导用户填写关键信息
3. 实时生成：展示生成进度，降低等待焦虑
4. 智能优化：AI主动提供多维度优化建议
5. 一键应用：优化建议可一键应用或部分采纳
6. 对比编辑：原版本和优化版本左右对比
7. 多轮迭代：支持持续优化，直到满意为止
8. 格式多样：支持Word、PDF、纯文本多种导出
```

---

#### 5.1.6 工作流可视化编辑交互流程

```
┌─────────────────────────────────────────────────────────────┐
│                    工作流可视化编辑完整流程                    │
└─────────────────────────────────────────────────────────────┘

步骤1：创建/打开工作流
┌────────────────────────────────────────────────────┐
│  页面：工作流开发                                    │
│  ┌──────────────────────┐  ┌────────────────────┐│
│  │ 我的流程             │  │ 流程模板           ││
│  │ ┌────────────────┐   │  │ ┌──────────────┐  ││
│  │ │ 数据同步流程   │   │  │ │ 审批流模板    │  ││
│  │ │ 客户服务流程   │   │  │ │ 通知流模板    │  ││
│  │ │ [+] 新建流程   │   │  │ │ 数据处理模板  │  ││
│  │ └────────────────┘   │  │ └──────────────┘  ││
│  └──────────────────────┘  └────────────────────┘│
│                                                     │
│  [从空白创建] [从模板创建] [导入流程]               │
└────────────────────────────────────────────────────┘
         ↓ [点击：新建流程]
步骤2：工作流编辑器布局
┌────────────────────────────────────────────────────┐
│  顶部：工具栏                                       │
│  [保存] [调试] [发布] [设置] [⋮ 更多]              │
│  ─────────────────────────────────────────────    │
│  左侧：节点面板         │      画布区域          │
│  ┌──────────────────┐   │  ┌──────────────────┐ │
│  │ 📥 触发器        │   │  │                  │ │
│  │  • Webhook       │   │  │   [开始节点]     │ │
│  │  • 定时任务      │   │  │      ▼           │ │
│  │  • 手动触发      │   │  │  [处理节点]      │ │
│  │                  │   │  │      ▼           │ │
│  │ ⚙ 动作          │   │  │  [条件分支]      │ │
│  │  • HTTP请求      │   │  │     ├─ 是       │ │
│  │  • 数据库操作    │   │  │     └─ 否       │ │
│  │  • 发送通知      │   │  │                  │ │
│  │  • 调用Bot       │   │  │   [结束节点]     │ │
│  │                  │   │  │                  │ │
│  │ 🔄 逻辑控制      │   │  │   [画布网格]     │ │
│  │  • 条件分支      │   │  │                  │ │
│  │  • 循环          │   │  │                  │ │
│  │  • 并行          │   │  │                  │ │
│  │                  │   │  │                  │ │
│  │ 📊 数据处理      │   │  │                  │ │
│  │  • 数据转换      │   │  │                  │ │
│  │  • 数据过滤      │   │  │                  │ │
│  └──────────────────┘   │  └──────────────────┘ │
│                         │                       │
│  [搜索节点...]          │   右侧：属性面板     │
│                         │   ┌──────────────────┐│
│                         │   │ 节点属性          ││
│                         │   │ ──────────────────││
│                         │   │ 名称: HTTP请求    ││
│                         │   │ URL: ____________ ││
│                         │   │ 方法: [POST ▼]   ││
│                         │   │ Headers:          ││
│                         │   │ Body:             ││
│                         │   │                   ││
│                         │   │ 高级设置          ││
│                         │   │ □ 超时重试        ││
│                         │   │ □ 错误处理        ││
│                         │   └──────────────────┘│
└────────────────────────────────────────────────────┘

步骤3：拖拽添加节点
┌────────────────────────────────────────────────────┐
│  用户从左侧面板拖拽"HTTP请求"节点到画布             │
│                                                     │
│  画布状态变化：                                     │
│  ┌────────────────────────────────────────────┐   │
│  │                  [开始节点]                 │   │
│  │                      ▼                      │   │
│  │                  [处理节点]                 │   │
│  │                      ▼                      │   │
│  │              ┌─────────────────┐            │   │
│  │              │  HTTP请求       │ ← 拖拽中   │   │
│  │              │  (虚线轮廓)     │   跟随鼠标  │   │
│  │              └─────────────────┘            │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  智能吸附提示：                                     │
│  💡 松开即可插入到"处理节点"之后                    │
└────────────────────────────────────────────────────┘
         ↓ [松开鼠标，节点放置]
步骤4：连接节点
┌────────────────────────────────────────────────────┐
│  画布状态：节点已添加                               │
│  ┌────────────────────────────────────────────┐   │
│  │              [开始节点]                     │   │
│  │                  ▼                          │   │
│  │              [处理节点]                     │   │
│  │                  ▼                          │   │
│  │              [HTTP请求]                     │   │
│  │              (输出点) ●                     │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  用户从"HTTP请求"的输出点拖拽连线                   │
│  ┌────────────────────────────────────────────┐   │
│  │              [HTTP请求]                     │   │
│  │              (输出点) ●────────────         │   │
│  │                              ↕ 拖拽中       │   │
│  │                              ● (输入点)     │   │
│  │                          [结束节点]         │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  智能连线提示：                                     │
│  💡 连接到"结束节点"将完成工作流                    │
└────────────────────────────────────────────────────┘
         ↓ [松开鼠标，连线完成]
步骤5：配置节点属性
┌────────────────────────────────────────────────────┐
│  右侧属性面板：配置HTTP请求节点                     │
│  ┌────────────────────────────────────────────┐   │
│  │ ⚙ HTTP请求 - 属性配置                      │   │
│  │ ────────────────────────────────────────── │   │
│  │ 基本配置                                     │   │
│  │ 请求名称: [查询用户数据____________]         │   │
│  │ 请求URL:  [https://api.example.com/users]  │   │
│  │ 请求方法: [POST ▼]                          │   │
│  │                                             │   │
│  │ 请求头                                      │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ Content-Type: application/json        │ │   │
│  │ │ Authorization: Bearer {{token}}        │ │   │
│  │ │ [+ 添加请求头]                          │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ 请求体                                      │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ {                                       │ │   │
│  │ │   "userId": "{{input.userId}}"         │ │   │
│  │ │ }                                       │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ 高级设置 ▼                                  │   │
│  │ ☑️ 失败时重试                              │   │
│  │ 重试次数: [3]                               │   │
│  │ 超时时间: [30] 秒                           │   │
│  │                                             │   │
│  │ [取消] [应用]                               │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓ [点击应用]
步骤6：条件分支编辑
┌────────────────────────────────────────────────────┐
│  画布状态：添加条件分支节点                         │
│  ┌────────────────────────────────────────────┐   │
│  │              [处理数据]                     │   │
│  │                  ▼                          │   │
│  │          ┌───────────────┐                 │   │
│  │          │  条件分支     │                 │   │
│  │          │  response.ok │                 │   │
│  │          └───────┬───────┘                 │   │
│  │                  │                         │   │
│  │        ┌─────────┴─────────┐               │   │
│  │        ▼                   ▼               │   │
│  │   [true 分支]          [false 分支]         │   │
│  │   [保存数据]            [记录错误]          │   │
│  │        │                   │               │   │
│  │        └─────────┬─────────┘               │   │
│  │                  ▼                         │   │
│  │              [结束节点]                     │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  条件配置面板：                                     │
│  ┌────────────────────────────────────────────┐   │
│  │ 🔀 条件分支 - 配置                          │   │
│  │ ────────────────────────────────────────── │   │
│  │ 条件表达式                                   │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ {{ response.ok }} == true                 │ │   │
│  │ (支持JavaScript表达式)                    │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ 分支配置                                     │   │
│  │ ✓ true 分支: [保存数据]                    │   │
│  │ ✓ false 分支: [记录错误]                   │   │
│  │                                             │   │
│  │ [+ 添加更多分支]                             │   │
│  │ [取消] [应用]                               │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘

步骤7：调试工作流
┌────────────────────────────────────────────────────┐
│  用户点击[调试]按钮                                │
│  ┌────────────────────────────────────────────┐   │
│  │ 🐞 调试模式                                 │   │
│  │ ────────────────────────────────────────── │   │
│  │ 测试数据输入                                 │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ {                                       │ │   │
│  │ │   "userId": "12345",                    │ │   │
│  │ │   "action": "query"                     │ │   │
│  │ │ }                                       │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ [开始调试]                                  │   │
│  └────────────────────────────────────────────┘   │
│         ↓ [点击开始调试]                         │
│  ┌────────────────────────────────────────────┐   │
│  │ 🐞 调试执行中...                            │   │
│  │ ────────────────────────────────────────── │   │
│  │ 画布状态：                                   │   │
│  │ ✓ [开始节点] (已执行)                       │   │
│  │   ▼                                         │   │
│  │ ✓ [处理数据] (已执行, 125ms)               │   │
│  │   ▼                                         │   │
│  │ ⏳ [条件分支] (执行中...)                    │   │
│  │                                             │   │
│  │ 执行日志：                                   │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ [10:23:45] 开始节点: 开始执行           │ │   │
│  │ │ [10:23:45] 处理数据: 数据解析完成       │ │   │
│  │ │ [10:23:46] 条件分支: 条件判断中...      │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ [停止调试] [查看详情]                       │   │
│  └────────────────────────────────────────────┘   │
│         ↓ [执行完成]                             │
│  ┌────────────────────────────────────────────┐   │
│  │ ✅ 调试完成                                 │   │
│  │ ────────────────────────────────────────── │   │
│  │ 画布状态：                                   │   │
│  │ ✓ [开始节点] (已执行)                       │   │
│  │   ▼                                         │   │
│  │ ✓ [处理数据] (125ms)                       │   │
│  │   ▼                                         │   │
│  │ ✓ [条件分支] (43ms)                        │   │
│  │   ├─ ✓ [true分支] (执行)                   │   │
│  │   └─ ⊘ [false分支] (跳过)                  │   │
│  │       ▼                                     │   │
│  │   ✓ [保存数据] (89ms)                      │   │
│  │       ▼                                     │   │
│  │   ✓ [结束节点] (完成)                       │   │
│  │                                             │   │
│  │ 执行结果：                                   │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ ✅ 所有节点执行成功                      │ │
│  │ │ 总耗时: 257ms                            │ │   │
│  │ │ 输出数据:                                │ │
│  │ │ {                                        │ │   │
│  │ │   "status": "success",                  │ │   │
│  │ │   "saved": true                         │ │   │
│  │ │ }                                        │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ [再次调试] [保存] [发布]                    │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘

关键交互设计：
1. 拖拽式编排：所见即所得，直观易用
2. 智能吸附：节点自动对齐到网格
3. 智能连线：输出点到输入点的自动连线
4. 实时预览：配置即时生效
5. 变量引用：支持{{variable}}语法引用上下文变量
6. 调试模式：单步执行、断点调试、日志查看
7. 版本管理：支持多版本保存和回滚
8. 协作编辑：多人实时协作（类似Figma）
```

---

#### 5.1.7 组织架构管理交互流程

```
┌─────────────────────────────────────────────────────────────┐
│                    组织架构管理完整流程                       │
└─────────────────────────────────────────────────────────────┘

步骤1：进入组织管理
┌────────────────────────────────────────────────────┐
│  页面：组织管理                                      │
│  ──────────────────────────────────────────────   │
│  面包屑：首页 / 管理中台 / 组织管理                 │
│                                                     │
│  [创建组织] [批量导入] [导出] [刷新]               │
└────────────────────────────────────────────────────┘
         ↓ [点击：创建子组织]
步骤2：创建子组织
┌────────────────────────────────────────────────────┐
│  创建组织对话框                                     │
│  ┌────────────────────────────────────────────┐   │
│  │ 创建组织                                     │   │
│  │ ────────────────────────────────────────── │   │
│  │ 组织信息                                     │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ 组织名称: [_________________________]  │ │   │
│  │ ├─────────────────────────────────────────┤ │   │
│  │ │ 上级组织: [公司 ▼]                      │ │   │
│  │ ├─────────────────────────────────────────┤ │   │
│  │ │ 组织类型: [部门 ▼]                      │ │   │
│  │ ├─────────────────────────────────────────┤ │   │
│  │ │ 负责人:   [选择负责人____________]      │ │   │
│  │ ├─────────────────────────────────────────┤ │   │
│  │ │ 描述:     [_________________________]  │ │   │
│  │ │           [_________________________]  │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ 高级设置 ▼                                   │   │
│  │ ☑️ 继承上级组织权限                         │   │
│  │ ☐ 启用部门独立预算                          │   │
│  │                                             │   │
│  │ [取消] [创建]                               │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓ [填写完成后点击创建]
步骤3：添加成员到组织
┌────────────────────────────────────────────────────┐
│  右侧成员列表区域                                   │
│  ┌────────────────────────────────────────────┐   │
│  │ 搜索: [___________] [搜索]                 │   │
│  │ 筛选: [技术部▼] [全部岗位▼]               │   │
│  │                                             │   │
│  │ 成员列表：                                   │   │
│  │ ┌────────────────────────────────────────┐ │   │
│  │ │ [头像] 张三  产品经理  [编辑] [删除]  │ │   │
│  │ ├────────────────────────────────────────┤ │   │
│  │ │ [头像] 李四  前端开发  [编辑] [删除]  │ │   │
│  │ ├────────────────────────────────────────┤ │   │
│  │ │ [头像] 王五  后端开发  [编辑] [删除]  │ │   │
│  │ └────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ [+ 添加成员] [批量导入] [导出]             │   │
│  └────────────────────────────────────────────┘   │
│         ↓ [点击：添加成员]                      │
│  ┌────────────────────────────────────────────┐   │
│  │ 添加成员                                     │   │
│  │ ────────────────────────────────────────── │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ 选择成员：                               │ │   │
│  │ │ ┌─────────────────────────────────────┐ │ │   │
│  │ │ │ ☑ 张三 (zhangsan@example.com)      │ │ │   │
│  │ │ │ ☐ 李四 (lisi@example.com)          │ │ │   │
│  │ │ │ ☐ 王五 (wangwu@example.com)         │ │ │   │
│  │ │ │ [+ 邀请新成员]                       │ │ │   │
│  │ │ └─────────────────────────────────────┘ │ │   │
│  │ ├─────────────────────────────────────────┤ │   │
│  │ │ 分配角色：                               │ │   │
│  │ │ ☑️ 部门管理员   ☑️ 普通成员              │ │   │
│  │ ├─────────────────────────────────────────┤ │   │
│  │ │ 岗位：                                  │ │   │
│  │ │ [产品经理 ▼]                            │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ [取消] [确认添加]                           │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓ [选择成员并确认]
步骤4：权限配置
┌────────────────────────────────────────────────────┐
│  成员操作下拉菜单 → [权限配置]                     │
│  ┌────────────────────────────────────────────┐   │
│  │ 🔐 权限配置 - 张三                          │   │
│  │ ────────────────────────────────────────── │   │
│  │                                             │   │
│  │ 基础权限                                     │   │
│  │ ☑️ 查看组织成员                              │   │
│  │ ☑️ 邀请组织成员                              │   │
│  │ ☑️ 移除组织成员                              │   │
│  │ ☐ 编辑组织信息                              │   │
│  │                                             │   │
│  │ 数字员工权限                                 │   │
│  │ ☑️ 使用所有已发布员工                        │   │
│  │ ☑️ 创建自己的数字员工                        │   │
│  │ ☐ 编辑他人的数字员工                         │   │
│  │ ☐ 删除数字员工                               │   │
│  │                                             │   │
│  │ 数据权限                                     │   │
│  │ ☑️ 查看自己部门数据                          │   │
│  │ ☐ 查看全公司数据                             │   │
│  │ ☑️ 导出数据报表                              │   │
│  │                                             │   │
│  │ 权限模板：[部门管理员 ▼]  [保存为模板]      │   │
│  │                                             │   │
│  │ [取消] [保存]                                │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓ [保存权限]
步骤5：批量操作
┌────────────────────────────────────────────────────┐
│  选中多个成员（勾选框）                             │
│  ┌────────────────────────────────────────────┐   │
│  │ ☑️ [头像] 张三  产品经理                   │   │
│  │ ☑️ [头像] 李四  前端开发                   │   │
│  │ ☑️ [头像] 王五  后端开发                   │   │
│  │                                             │   │
│  │ 批量操作：                                   │   │
│  │ [变更部门] [分配角色] [批量删除] [导出]    │   │
│  └────────────────────────────────────────────┘   │
│         ↓ [点击：变更部门]                       │
│  ┌────────────────────────────────────────────┐   │
│  │ 批量变更部门                                 │   │
│  │ ────────────────────────────────────────── │   │
│  │ 已选择：3名成员                              │   │
│  │                                             │   │
│  │ 目标部门：                                   │   │
│  │ [技术部 ▼] → [产品部 ▼]                    │   │
│  │                                             │   │
│  │ 同时变更岗位：                               │   │
│  │ ☐ 是，指定新岗位 [_____________]            │   │
│  │ ☑️ 否，保持原岗位                            │   │
│  │                                             │   │
│  │ [取消] [确认变更]                            │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓
步骤6：组织架构调整（拖拽）
┌────────────────────────────────────────────────────┐
│  左侧组织树：拖拽节点                               │
│  ┌──────────────────────────┐                    │
│  │ 🏢 公司                 │                    │
│  │   ├ 技术部   │   │   │   ├─ 后端开发          │
│  │   │   ├─ 前端开发       │   │   │   └─ 测试工程          │
│  │   │   └─ 后端开发   │   │   │   ├ 产品部  ←─────────│
│  │   │   └─ 测试工程      │   │   │       │   (拖拽中)    │
│  │   ├ 产品部             │   │   │       ▼               │
│  │   └ 运营部             │   │   │   └ 运营部             │
│  └──────────────────────────┘                    │
│                                                     │
│  拖拽提示：                                         │
│  💡 将"产品部"移动到"技术部"下作为子部门            │
│  ⚠️ 移动后将继承技术部的权限设置                    │
│                                                     │
│  [确认移动] [取消]                                 │
└────────────────────────────────────────────────────┘

关键交互设计：
1. 树形结构：可展开/收起的组织树
2. 拖拽调整：拖拽节点调整组织结构
3. 实时搜索：成员/部门实时搜索过滤
4. 批量操作：多选成员进行批量处理
5. 权限模板：预设权限角色，快速分配
6. 导入导出：支持Excel批量导入成员
7. 审批流程：组织调整需审批（可配置）
```

---

#### 5.1.8 数据分析与看板交互流程

```
┌─────────────────────────────────────────────────────────────┐
│                    数据分析与看板完整流程                     │
└─────────────────────────────────────────────────────────────┘

步骤1：进入数据看板
┌────────────────────────────────────────────────────┐
│  页面：数据看板                                      │
│  ──────────────────────────────────────────────   │
│  筛选：[最近30天▼] [全公司▼] [全部员工▼]         │
│  [导出报表] [刷新] [自定义看板]                    │
└────────────────────────────────────────────────────┘
         ↓ [查看核心指标]
步骤2：查看KPI指标卡片
┌────────────────────────────────────────────────────┐
│  核心指标卡片（4列）                                │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────┐│
│  │ 总用户数  │ │ 活跃用户 │ │ 对话次数 │ │满意度││
│  │  1,234   │ │   856    │ │  12.5k   │ │ 4.7⭐││
│  │  +12% ↑  │ │  +8% ↑   │ │  +15% ↑  │ │+0.2↑││
│  │  📈      │ │  📈      │ │  📈      │ │  📈 ││
│  └──────────┘ └──────────┘ └──────────┘ └──────┘│
│                                                     │
│  用户点击"活跃用户"卡片 → 钻取详情                  │
└────────────────────────────────────────────────────┘
         ↓ [卡片点击]
步骤3：数据钻取
┌────────────────────────────────────────────────────┐
│  活跃用户详细数据                                   │
│  ┌────────────────────────────────────────────┐   │
│  │ ◀ 返回看板    活跃用户详情 - 最近30天      │   │
│  │ ────────────────────────────────────────── │   │
│  │                                             │   │
│  │ 时间趋势图                                   │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │   活跃用户数                              │ │   │
│  │ │ 1000│    ████                           │ │   │
│  │ │  800│  ████                             │ │   │
│  │ │  600│ ████                              │ │   │
│  │ │  400│████                               │ │   │
│  │ │    └────────────────────                │ │   │
│  │ │      W1 W2 W3 W4                        │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ 用户活跃度分布                               │   │
│  │ ┌─────────────┬───────────┬─────────────┐  │   │
│  │ │ 高活跃      │ 中活跃    │ 低活跃      │  │   │
│  │ │ ████ 45%   │ ███ 35%   │ ██ 20%      │  │   │
│  │ └─────────────┴───────────┴─────────────┘  │   │
│  │                                             │   │
│  │ 活跃用户列表（Top 10）                       │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ 排名 │ 用户   │ 部门  │ 活跃度 │ 趋势  │ │   │
│  │ ├─────────────────────────────────────────┤ │   │
│  │ │  1  │ 张三  │ 技术  │ 98%    │ ↑     │ │   │
│  │ │  2  │ 李四  │ 产品  │ 95%    │ ↑     │ │   │
│  │ │  3  │ 王五  │ 运营  │ 92%    │ →     │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ [查看全部用户] [导出数据]                    │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓ [点击趋势图中的"W4"柱]
步骤4：进一步钻取
┌────────────────────────────────────────────────────┐
│  W4（第4周）详细数据                               │
│  ┌────────────────────────────────────────────┐   │
│  │ ◀ 返回    第4周活跃用户详情 (12/23-12/29)   │   │
│  │ ────────────────────────────────────────── │   │
│  │                                             │   │
│  │ 每日活跃趋势                                 │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ 200│ ████                               │ │   │
│  │ │ 150│ ████████                          │ │   │
│  │ │ 100│ ████████████████                  │ │   │
│  │ │  50│ ████████████████████████████████  │ │   │
│  │ │   0└────────────────────────────────    │ │   │
│  │ │     一 二 三 四 五 六 日                 │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ 按部门分布                                   │   │
│  │ ┌─────────────────────────────────────────┐ │   │
│  │ │ 技术部: ████████████ 42%               │ │   │
│  │ │ 产品部: ███████ 28%                    │ │   │
│  │ │ 运营部: █████ 18%                      │ │   │
│  │ │ 其他部: ██ 12%                         │ │   │
│  │ └─────────────────────────────────────────┘ │   │
│  │                                             │   │
│  │ [返回上一级] [导出报表]                     │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
         ↓ [点击：自定义看板]
步骤5：自定义看板
┌────────────────────────────────────────────────────┐
│  自定义看板编辑器                                   │
│  ┌────────────────────────────────────────────┐   │
│  │ 自定义数据看板                               │   │
│  │ ────────────────────────────────────────── │   │
│  │ 左侧：组件库    │    画布区域              │   │
│  │ ┌────────────┐ │  ┌──────────────────────┐│   │
│  │ │ 📊 指标卡  │ │  │ ┌────────┐ ┌────────┐││   │
│  │ │   ├ 总用户 │ │  │ │用户数  │ │活跃度  │││   │
│  │ │   ├ 活跃度 │ │  │ └────────┘ └────────┘││   │
│  │ │   └ 对话数 │ │  │                        ││   │
│  │ ├────────────┤ │  │ ┌──────────────────┐  ││   │
│  │ │ 📈 趋势图  │ │  │ │  用户增长趋势     │  ││   │
│  │ │   ├ 折线图 │ │  │ │    [折线图]       │  ││   │
│  │ │   ├ 柱状图 │ │  │ └──────────────────┘  ││   │
│  │ │   └ 饼图   │ │  │                        ││   │
│  │ ├────────────┤ │  │ ┌──────────────────┐  ││   │
│  │ │ 📋 数据表  │ │  │ │  员工使用排行     │  ││   │
│  │ └────────────┘ │  │ │    [表格]         │  ││   │
│  │                │  │ └──────────────────┘  ││   │
│  │ [搜索组件...] │  │                        ││   │
│  └───────────────┘  └──────────────────────────┘│
│   拖拽组件到画布                                  │   │
│                                                     │
│  右侧：配置面板                                     │
│  ┌────────────────────────────────────────────┐   │
│  │ 组件配置：用户数                             │   │
│  │ ────────────────────────────────────────── │   │
│  │ 标题：[总用户数_________]                   │   │
│  │ 数据源：[用户统计表 ▼]                      │   │
│  │ 指标：[COUNT(user_id) ▼]                   │   │
│  │ 筛选：[WHERE created_at > ▼]                │   │
│  │                                             │   │
│  │ 样式配置                                     │   │
│  │ 主题色：[🔵 蓝色 ▼]                          │   │
│  │ 大小：[中等 ▼]                               │   │
│  │ 显示趋势：☑️                                 │   │
│  │                                             │   │
│  │ [删除组件] [复制] [保存]                     │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  [预览] [保存看板] [发布]                           │
└────────────────────────────────────────────────────┘
         ↓ [配置完成后保存]
步骤6：导出报表
┌────────────────────────────────────────────────────┐
│  用户点击[导出报表]按钮                             │
│  ┌────────────────────────────────────────────┐   │
│  │ 导出数据报表                                 │   │
│  │ ────────────────────────────────────────── │   │
│  │ 导出范围                                     │   │
│  │ ⦿ 当前看板                                   │   │
│  │ ⦿ 自定义范围                                 │   │
│  │                                             │   │
│  │ 时间范围：                                   │   │
│  │ [最近30天 ▼]                                 │   │
│  │                                             │   │
│  │ 导出格式：                                   │   │
│  │ ⦿ Excel (含数据)                             │   │
│  │ ⦿ PDF (含图表)                               │   │
│  │ ⦿ PNG (图片)                                 │   │
│  │                                             │   │
│  │ 包含内容：                                   │   │
│  │ ☑️ KPI指标卡片                               │   │
│  │ ☑️ 趋势图表                                  │   │
│  │ ☑️ 数据表格                                  │   │
│  │ ☑️ 数据洞察                                  │   │
│  │                                             │   │
│  │ [取消] [导出]                                │   │
│  └────────────────────────────────────────────┘   │
│         ↓ [选择PDF并点击导出]                   │
│  ┌────────────────────────────────────────────┐   │
│  │ 📊 正在生成报表...                          │   │
│  │ ◐ 收集数据（35%）                           │   │
│  │ ◐ 生成图表（68%）                           │   │
│  │ ◐ 编译PDF（92%）                            │   │
│  │ ✓ 生成完成（100%）                          │   │
│  │                                             │   │
│  │ 报表已生成：                                 │   │
│  │ 📄 数据看板_2024-12-30.pdf (2.3MB)          │   │
│  │                                             │   │
│  │ [下载] [打开] [分享链接]                     │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘

步骤7：实时监控（自动刷新）
┌────────────────────────────────────────────────────┐
│  用户启用自动刷新功能                               │
│  ┌────────────────────────────────────────────┐   │
│  │ 自动刷新：☑️ 已启用                         │   │
│  │ 刷新间隔：[5分钟 ▼]                         │   │
│  │ 上次刷新：10:23:45                          │   │
│  │ 下次刷新：10:28:45                          │   │
│  │                                             │   │
│  │ 实时数据流：                                 │   │
│  │ [10:25:12] +1 新用户注册                    │   │
│  │ [10:26:34] 对话完成 (用时: 2m15s)           │   │
│  │ [10:27:01] 满意度评分: 5⭐                  │   │
│  │ [10:27:45] 数字员工"客服助手"被调用        │   │
│  └────────────────────────────────────────────┘   │
│                                                     │
│  数据变化时，卡片实时更新动画：                      │
│  ┌──────────┐                                      │
│  │ 对话次数  │ ← 数字滚动动画                      │
│  │ 12,547   │   12,548 → 12,549 → 12,550         │
│  │ +15 ↑    │                                      │
│  └──────────┘                                      │
└────────────────────────────────────────────────────┘

关键交互设计：
1. 下钻分析：点击卡片/图表逐层深入
2. 拖拽式看板：自定义数据看板布局
3. 实时刷新：支持自动刷新和WebSocket推送
4. 多维筛选：时间/组织/员工多维度组合筛选
5. 智能导出：支持Excel/PDF/PNG多种格式
6. 数据动画：数字滚动、图表动态加载
7. 告警设置：可配置数据异常告警
```

---

### 5.2 交互设计原则

```typescript
/**
 * 全局交互设计原则
 */
const interactionPrinciples = {
  // 1. 即时反馈
  immediateFeedback: {
    description: '用户每个操作都应有即时反馈',
    examples: [
      '按钮点击：按下态 + 颜色变化',
      '加载状态：Spinner + 加载文案',
      '操作成功：Toast提示 + 绿色对勾',
      '操作失败：错误提示 + 红色叉号',
    ]
  },

  // 2. 渐进式披露
  progressiveDisclosure: {
    description: '高级功能默认隐藏，按需展示',
    examples: [
      '基础配置默认展开，高级配置折叠',
      '"显示更多"按钮展开完整内容',
      'Hover时显示工具提示',
      '首次使用显示引导浮层',
    ]
  },

  // 3. 防错设计
  errorPrevention: {
    description: '预防错误发生，而非事后补救',
    examples: [
      '表单实时验证，而非提交后验证',
      '删除操作二次确认（重要数据）',
      '不可逆操作醒目提示（红色）',
      '提供撤销功能（Ctrl+Z）',
    ]
  },

  // 4. 一致性
  consistency: {
    description: '相似的交互产生相似的结果',
    examples: [
      '所有页面的保存按钮位置一致',
      '统一的快捷键体系',
      '统一的颜色语义（绿色=成功，红色=危险）',
      '统一的图标含义',
    ]
  },

  // 5. 可恢复性
  recoverability: {
    description: '允许用户撤销或修正操作',
    examples: [
      '删除后30秒内可撤销',
      '表单自动保存草稿',
      '操作历史记录',
      '版本管理和回滚',
    ]
  }
};
```

---

## 六、响应式设计

### 6.1 断点系统

```css
/* =====================================================
   响应式断点 - Tailwind CSS标准
   ===================================================== */

/* 断点定义 */
--breakpoint-sm: 640px;    /* 小屏幕 - 手机竖屏 */
--breakpoint-md: 768px;    /* 中等屏幕 - 手机横屏/小平板 */
--breakpoint-lg: 1024px;   /* 大屏幕 - 平板/小笔记本 */
--breakpoint-xl: 1280px;   /* 超大屏幕 - 笔记本 */
--breakpoint-2xl: 1536px;  /* 2K屏幕 */

/* 响应式容器 */
.container {
  width: 100%;
  margin-right: auto;
  margin-left: auto;
  padding-right: var(--layout-padding-md);
  padding-left: var(--layout-padding-md);
}

@media (min-width: 640px) {
  .container {
    max-width: 640px;
  }
}

@media (min-width: 768px) {
  .container {
    max-width: 768px;
  }
}

@media (min-width: 1024px) {
  .container {
    max-width: 1024px;
  }
}

@media (min-width: 1280px) {
  .container {
    max-width: 1280px;
  }
}

/* 响应式网格 */
.grid {
  display: grid;
  gap: var(--spacing-md);
}

.grid-cols-1 {
  grid-template-columns: repeat(1, minmax(0, 1fr));
}

.grid-cols-2 {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.grid-cols-3 {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

@media (min-width: 640px) {
  .sm\:grid-cols-2 {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (min-width: 1024px) {
  .lg\:grid-cols-3 {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (min-width: 1280px) {
  .xl\:grid-cols-4 {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}
```

### 6.2 移动端适配

```css
/* =====================================================
   移动端特殊样式
   ===================================================== */

/* 侧边栏 - 移动端默认隐藏，点击滑出 */
@media (max-width: 1024px) {
  .sidebar {
    position: fixed;
    left: -240px;
    top: 64px;
    z-index: 999;
    transition: left var(--duration-normal) var(--ease-in-out);
  }

  .sidebar.open {
    left: 0;
  }

  /* 遮罩层 */
  .sidebar-overlay {
    position: fixed;
    top: 64px;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--bg-overlay);
    z-index: 998;
    opacity: 0;
    visibility: hidden;
    transition: all var(--duration-base) var(--ease-in-out);
  }

  .sidebar-overlay.show {
    opacity: 1;
    visibility: visible;
  }
}

/* 导航栏 - 移动端显示汉堡菜单 */
@media (max-width: 1024px) {
  .header-nav {
    display: none;
  }

  .header-hamburger {
    display: flex;
  }
}

/* 表格 - 移动端卡片化 */
@media (max-width: 768px) {
  .table {
    display: block;
  }

  .table thead {
    display: none;
  }

  .table tbody,
  .table tr,
  .table td {
    display: block;
    width: 100%;
  }

  .table tr {
    margin-bottom: var(--spacing-md);
    background-color: var(--bg-primary);
    border: 1px solid var(--border-primary);
    border-radius: var(--radius-md);
  }

  .table td {
    display: flex;
    justify-content: space-between;
    padding: var(--spacing-sm);
    text-align: left;
    border-bottom: 1px solid var(--border-primary);
  }

  .table td::before {
    content: attr(data-label);
    font-weight: var(--font-weight-semibold);
    color: var(--text-secondary);
    margin-right: var(--spacing-md);
  }
}
```

---

## 七、可访问性设计

### 7.1 WCAG 2.1 标准

```css
/* =====================================================
   可访问性设计 - 符合WCAG 2.1 AA级标准
   ===================================================== */

/* 1. 颜色对比度 */
/* 文本与背景对比度至少4.5:1（正常文本）或3:1（大文本） */
.text-primary {
  color: #111827;  /* 对比度: 15.3:1 ✓ */
}

.text-secondary {
  color: #6B7280;  /* 对比度: 4.6:1 ✓ */
}

/* 链接对比度 */
a {
  color: var(--primary-color);  /* 对比度: 4.5:1 ✓ */
}

/* 禁用状态 */
.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 2. 焦点可见性 */
*:focus-visible {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

/* 3. 键盘导航 */
/* 确保所有交互元素可通过键盘访问 */
button:focus,
a:focus,
input:focus,
select:focus,
textarea:focus {
  outline: 2px solid var(--primary-color);
}

/* 4. 屏幕阅读器支持 */
/* 隐藏视觉元素但保留给屏幕阅读器 */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}

/* 5. ARIA属性 */
<button aria-label="关闭对话框" aria-expanded="false">
  <IconClose />
</button>

<input aria-label="邮箱地址" aria-required="true" />

<div role="alert" aria-live="polite">
  操作成功
</div>
```

### 7.2 键盘快捷键

```typescript
/**
 * 全局键盘快捷键定义
 */

const keyboardShortcuts = {
  // 通用快捷键
  'Cmd/Ctrl + K': '打开全局搜索',
  'Cmd/Ctrl + /': '显示快捷键列表',
  'Cmd/Ctrl + B': '加粗',
  'Cmd/Ctrl + I': '斜体',
  'Cmd/Ctrl + Enter': '发送消息',
  'Esc': '关闭弹窗/取消操作',

  // 导航快捷键
  'Option/Alt + ←': '返回上一页',
  'Option/Alt + →': '前进下一页',
  'Cmd/Ctrl + Shift + H': '返回首页',

  // 对话快捷键
  '↑': '上一条消息',
  '↓': '下一条消息',
  'Cmd/Ctrl + ↑': '滚动到顶部',
  'Cmd/Ctrl + ↓': '滚动到底部',

  // 编辑器快捷键
  'Cmd/Ctrl + S': '保存',
  'Cmd/Ctrl + Z': '撤销',
  'Cmd/Ctrl + Shift + Z': '重做',
  'Cmd/Ctrl + /': '注释/取消注释',
};

// 快捷键提示组件
<KeyboardShortcutDialog>
  <h2>键盘快捷键</h2>
  <ShortcutList>
    <ShortcutItem keys={['⌘', 'K']} description="打开全局搜索" />
    <ShortcutItem keys={['⌘', '/']} description="显示快捷键" />
    <ShortcutItem keys={['⌘', 'Enter']} description="发送消息" />
    <ShortcutItem keys={['Esc']} description="取消操作" />
  </ShortcutList>
</KeyboardShortcutDialog>
```

---

## 八、设计交付物

### 8.1 设计文件组织

```
design-system/
├── 📁 tokens/                 # 设计token
│   ├── colors.json            # 颜色定义
│   ├── typography.json        # 字体定义
│   ├── spacing.json           # 间距定义
│   └── shadows.json           # 阴影定义
│
├── 📁 components/            # 组件设计稿
│   ├── buttons/              # 按钮
│   │   ├── primary.figma
│   │   ├── secondary.figma
│   │   └── sizes.figma
│   ├── forms/                # 表单
│   │   ├── inputs.figma
│   │   ├── selects.figma
│   │   └── checkboxes.figma
│   ├── cards/                # 卡片
│   ├── modals/               # 模态框
│   └── tables/               # 表格
│
├── 📁 pages/                 # 页面设计稿
│   ├── employee-frontend/    # 员工前台
│   │   ├── chatbox.figma
│   │   ├── employees.figma
│   │   ├── chatbi.figma
│   │   └── writing.figma
│   ├── admin-backend/        # 管理中台
│   │   ├── organization.figma
│   │   ├── employees-mgmt.figma
│   │   └── analytics.figma
│   └── dev-backend/          # 开发后台
│       ├── workspace.figma
│       ├── plugins.figma
│       └── workflows.figma
│
├── 📁 icons/                 # 图标资源
│   ├── svg/                  # SVG图标
│   └── font/                 # 图标字体
│
└── 📁 illustrations/         # 插画资源
    ├── empty-states/         # 空状态插画
    ├── error-states/         # 错误状态插画
    └── onboarding/           # 引导插画
```

### 8.2 组件代码规范

```typescript
/**
 * 组件代码组织规范
 */

// 1. 组件目录结构
components/
├── team/
│   ├── TeamList.tsx              # 组件主文件
│   ├── TeamList.module.scss      # 组件样式
│   ├── TeamList.test.tsx        # 组件测试
│   ├── index.ts                 # 导出文件
│   └── types.ts                 # 类型定义

// 2. 组件代码模板
import React, { FC } from 'react';
import styles from './TeamList.module.scss';

interface TeamListProps {
  teams: Team[];
  onTeamClick?: (team: Team) => void;
}

export const TeamList: FC<TeamListProps> = ({
  teams,
  onTeamClick,
}) => {
  return (
    <div className={styles.teamList}>
      {teams.map((team) => (
        <div
          key={team.id}
          className={styles.teamCard}
          onClick={() => onTeamClick?.(team)}
          role="button"
          tabIndex={0}
          aria-label={`查看团队${team.name}`}
        >
          {/* 组件内容 */}
        </div>
      ))}
    </div>
  );
};

// 3. 组件样式（CSS Modules）
.teamList {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.teamCard {
  padding: var(--spacing-lg);
  background-color: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-base) var(--ease-in-out);

  &:hover {
    box-shadow: var(--shadow-md);
    transform: translateY(-2px);
  }
}
```

---

## 九、设计实施指南

### 9.1 开发者使用指南

```typescript
/**
 * 设计系统使用示例
 */
import { DesignSystem } from '@/design-system';

// 1. 使用颜色
const primaryColor = DesignSystem.colors.primary;
const primaryHover = DesignSystem.colors.primary.hover;

// 2. 使用间距
const padding = DesignSystem.spacing.md; // 16px

// 3. 使用字体
const textStyle = {
  fontSize: DesignSystem.typography.fontSize.base,
  fontWeight: DesignSystem.typography.fontWeight.medium,
  lineHeight: DesignSystem.typography.lineHeight.normal,
};

// 4. 使用组件
import { Button, Card, Input } from '@/components';

<Button type="primary" size="medium">
  点击我
</Button>

// 5. 响应式布局
import { useBreakpoint } from '@/hooks';

const isMobile = useBreakpoint('sm');
const isTablet = useBreakpoint('md');
const isDesktop = useBreakpoint('lg');

// 6. 暗色模式（可选）
const { theme } = useTheme();
const darkMode = theme.mode === 'dark';
```

### 9.2 设计Token导出

```typescript
/**
 * design-tokens.ts
 * 导出为不同格式：CSS变量、Sass变量、JavaScript对象
 */

// 导出为CSS变量
export function toCSSVariables() {
  return `
:root {
  /* Colors */
  --primary-color: ${colors.primary};
  --primary-hover: ${colors.primary.hover};

  /* Spacing */
  --spacing-sm: ${spacing.sm}px;
  --spacing-md: ${spacing.md}px;

  /* Typography */
  --font-size-base: ${typography.fontSize.base}px;
}
`;
}

// 导出为Sass变量
export function toSassVariables() {
  return `
$primary-color: ${colors.primary};
$spacing-md: ${spacing.md}px;
$font-size-base: ${typography.fontSize.base}px;
`;
}

// 导出为JavaScript对象
export const designTokens = {
  colors,
  spacing,
  typography,
  borderRadius,
  shadows,
  transitions,
};
```

---

## 十、总结

### 10.1 设计系统核心价值

```
┌──────────────────────────────────────────────────────────┐
│            zker Enterprise 设计系统                       │
├──────────────────────────────────────────────────────────┤
│                                                            │
│  ✅ 统一性                                                   │
│     - 统一的视觉语言                                        │
│     - 一致的交互模式                                        │
│     - 规范的组件库                                          │
│                                                            │
│  ✅ 高效性                                                   │
│     - 基于Semi Design，开箱即用                              │
│     - 组件化开发，提升效率                                   │
│     - Design Token，一键换肤                                 │
│                                                            │
│  ✅ 可扩展性                                                 │
│     - 模块化设计，易于扩展                                   │
│     - 插件化架构，灵活定制                                   │
│     - 版本化管理，平滑升级                                   │
│                                                            │
│  ✅ 可访问性                                                 │
│     - 符合WCAG 2.1 AA级标准                                  │
│     - 支持键盘导航                                          │
│     - 屏幕阅读器友好                                        │
│                                                            │
│  ✅ 响应式                                                   │
│     - 移动优先策略                                          │
│     - 多端适配支持                                          │
│     - 平滑断点过渡                                          │
│                                                            │
└──────────────────────────────────────────────────────────┘
```

### 10.2 快速开始

```bash
# 1. 安装依赖
npm install @douyinfe/semi-ui @douyinfe/semi-icons

# 2. 引入样式
import '@douyinfe/semi-ui/dist/css/semi.min.css';
import '@/design-system/tokens/index.css';

# 3. 使用组件
import { Button, Card, Table } from '@douyinfe/semi-ui';

# 4. 自定义主题（可选）
import { DesignSystem } from '@/design-system';

DesignSystem.setTheme({
  primaryColor: '#306EEF',
  borderRadius: '6px',
});
```

---

**文档维护者**: UI/UX设计团队
**最后更新**: 2025-12-30
**文档版本**: v1.0

**相关文档**:
- 企业级功能完善与统一性设计方案
- 统一开发规范与企业级实施注意事项
- 团队管理API完整实现示例
