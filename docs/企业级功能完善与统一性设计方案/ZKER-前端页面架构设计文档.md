# ZKER 前端企业级页面架构设计文档

**版本**: v1.0.0
**最后更新**: 2025-01-01
**状态**: 架构设计完成，待UI实现

---

## 📋 目录

- [设计原则](#设计原则)
- [技术栈](#技术栈)
- [页面组织结构](#页面组织结构)
- [配额管理页面设计](#配额管理页面设计)
- [权限管理页面设计](#权限管理页面设计)
- [状态管理策略](#状态管理策略)
- [UI组件规范](#ui组件规范)
- [路由设计](#路由设计)
- [实施计划](#实施计划)

---

## 🎯 设计原则

### 1. SOLID 原则

**Single Responsibility（单一职责）**:
- 每个页面组件只负责一个功能域
- 拆分大型组件为多个小组件

**Open/Closed（开放封闭）**:
- 组件通过 props 扩展，而非修改内部
- 使用插槽（slots）模式支持自定义

**Liskov Substitution（里氏替换）**:
- 组件接口一致，可替换实现

**Interface Segregation（接口隔离）**:
- Props 接口专一，避免"胖 props"
- 使用 TypeScript 类型严格定义

**Dependency Inversion（依赖倒置）**:
- 依赖抽象（Hooks），不依赖具体实现

### 2. DRY 原则

**复用策略**:
- 共享组件: `@coze-arch/business-components`
- 自定义 Hooks: 已实现 `useTenant`, `useQuota`, `usePermission`
- 工具函数: 统一导出到 `utils/`

### 3. KISS 原则

**简化设计**:
- 避免过度抽象
- 优先使用现有组件库（Semi Design）
- 减少不必要的状态管理

### 4. YAGNI 原则

**避免过度设计**:
- 仅实现当前明确需要的功能
- 不预留"未来可能用到"的代码
- 删除未使用的导入和变量

---

## 🛠️ 技术栈

### 核心框架

| 技术 | 版本 | 用途 |
|------|------|------|
| **React** | 18.3.1 | UI框架 |
| **TypeScript** | 5.8.2 | 类型系统 |
| **React Router** | latest | 路由管理 |
| **Zustand** | latest | 状态管理（轻量） |

### UI 组件库

| 库 | 用途 |
|------|------|
| **@coze-arch/coze-design** | Semi Design 组件库 |
| **@coze-arch/business-components** | 业务组件库 |
| **@douyinfe/semi-illustrations** | 插图资源 |

### 数据层

| 层级 | 实现 | 说明 |
|------|------|------|
| **API 层** | `tenantApi`, `quotaApi`, `permissionApi` | 已完成 |
| **Hooks 层** | `useTenant`, `useQuota`, `usePermission` | 已完成 |
| **组件层** | 页面组件 | 待实现 |

---

## 📁 页面组织结构

### 推荐目录结构

```
frontend/
├── apps/
│   └── coze-studio/
│       └── src/
│           ├── pages/
│           │   ├── tenant/              # 租户管理页面
│           │   │   ├── TenantList.tsx    # 租户列表
│           │   │   ├── TenantDetail.tsx  # 租户详情
│           │   │   └── Subscription.tsx   # 订阅管理
│           │   ├── quota/               # 配额管理页面
│           │   │   ├── QuotaOverview.tsx # 配额概览
│           │   │   ├── QuotaDetail.tsx   # 配额详情
│           │   │   └── QuotaSettings.tsx # 配额设置
│           │   ├── permission/          # 权限管理页面
│           │   │   ├── RoleList.tsx      # 角色列表
│           │   │   ├── RoleDetail.tsx    # 角色详情
│           │   │   ├── Department.tsx    # 部门管理
│           │   │   └── PermissionMatrix.tsx # 权限矩阵
│           │   └── layout/
│           │       ├── EnterpriseLayout.tsx
│           │       └── Sidebar.tsx
│           └── routes/
│               └── enterprise.tsx        # 企业级路由配置
└── packages/
    └── arch/
        └── bot-api/
            ├── tenant-api.ts            # API 层（已完成）
            ├── quota-api.ts
            ├── permission-api.ts
            ├── hooks/                    # Hooks 层（已完成）
            │   ├── useTenant.ts
            │   ├── useQuota.ts
            │   └── usePermission.ts
            └── types/                    # 类型定义（已完成）
                ├── tenant.types.ts
                ├── quota.types.ts
                └── permission.types.ts
```

### 组件层次结构

```
EnterpriseLayout
├── Sidebar (导航菜单)
├── Header (顶部栏)
└── Content (内容区域)
    ├── TenantList / TenantDetail
    ├── QuotaOverview / QuotaDetail
    └── RoleList / RoleDetail / Department
```

---

## 📊 配额管理页面设计

### 页面列表

| 页面 | 路由 | 功能 |
|------|------|------|
| **配额概览** | `/quota/overview` | 所有资源配额总览仪表板 |
| **配额详情** | `/quota/:resourceType` | 单个资源配额详情与趋势 |
| **配额设置** | `/quota/settings` | 配额限制设置（超级管理员） |

### 1. 配额概览页面 (QuotaOverview)

**组件结构**:
```tsx
<QuotaOverview>
  <QuotaSummaryCard />        {/* 配额总览卡片 */}
  <QuotaWarningCard />        {/* 配额警告卡片 */}
  <QuotaTable />              {/* 配额列表表格 */}
  <QuotaUsageChart />          {/* 使用率图表 */}
</QuotaOverview>
```

**核心功能**:
- 显示所有资源类型的配额状态
- 整体使用率仪表板
- 配额警告/告警提示
- 快速跳转到配额详情

**数据获取**:
```typescript
import { useQuotas, useQuotaWarning } from '@coze-arch/bot-api';

function QuotaOverview() {
  const { quotas, overallUsage, loading } = useQuotas('tenant_id');
  const { hasWarning, hasCritical, warnings } = useQuotaWarning(quotas);

  return (
    <div>
      <QuotaSummaryCard overallUsage={overallUsage} />
      {hasWarning && <QuotaWarningCard warnings={warnings} />}
      <QuotaTable quotas={quotas} loading={loading} />
    </div>
  );
}
```

### 2. 配额详情页面 (QuotaDetail)

**组件结构**:
```tsx
<QuotaDetail>
  <QuotaInfoCard />            {/* 配额信息卡片 */}
  <QuotaUsageChart />           {/* 使用量趋势图 */}
  <QuotaHistoryTable />         {/* 历史记录表格 */}
  <QuotaPredictionCard />       {/* 预计耗尽时间 */}
</QuotaDetail>
```

**核心功能**:
- 资源配额详细信息
- 使用量趋势图（折线图）
- 历史记录表格
- 预计耗尽时间提醒
- 快速升级套餐入口

**数据获取**:
```typescript
import { useQuotaUsage } from '@coze-arch/bot-api';

function QuotaDetail({ resourceType }: { resourceType: ResourceType }) {
  const { usage, trend, estimatedExhaustion, loading } = useQuotaUsage(
    'tenant_id',
    resourceType,
  );

  return (
    <div>
      <QuotaInfoCard usage={usage} />
      <QuotaUsageChart usage={usage} trend={trend} />
      {estimatedExhaustion && (
        <QuotaPredictionCard estimatedExhaustion={estimatedExhaustion} />
      )}
    </div>
  );
}
```

### 3. 配额设置页面 (QuotaSettings)

**组件结构**:
```tsx
<QuotaSettings>
  <QuotaLimitForm />            {/* 配额限制表单 */}
  <QuotaResetPeriodSelector />  {/* 重置周期选择 */}
</QuotaSettings>
```

**权限要求**: 仅超级管理员可访问

**核心功能**:
- 修改各资源类型的配额限制
- 设置重置周期
- 批量配置功能

---

## 🔐 权限管理页面设计

### 页面列表

| 页面 | 路由 | 功能 |
|------|------|------|
| **角色列表** | `/permission/roles` | 角色列表与创建 |
| **角色详情** | `/permission/roles/:roleId` | 角色详情与编辑 |
| **部门管理** | `/permission/departments` | 部门树形管理 |
| **权限矩阵** | `/permission/matrix` | 权限矩阵视图 |

### 1. 角色列表页面 (RoleList)

**组件结构**:
```tsx
<RoleList>
  <PageHeader>
    <Button onClick={openCreateModal}>创建角色</Button>
  </PageHeader>
  <RoleTable />               {/* 角色表格 */}
  <CreateRoleModal />         {/* 创建角色弹窗 */}
</RoleList>
```

**核心功能**:
- 角色列表展示（Table）
- 搜索与筛选
- 创建/编辑/删除角色
- 查看角色权限

**数据获取**:
```typescript
import { useRoles, useCreateRole, useDeleteRole } from '@coze-arch/bot-api';

function RoleList() {
  const { roles, loading, refetch } = useRoles('tenant_id');
  const { createRole, loading: creating } = useCreateRole();
  const { deleteRole, loading: deleting } = useDeleteRole();

  // ... 组件实现
}
```

### 2. 角色详情页面 (RoleDetail)

**组件结构**:
```tsx
<RoleDetail>
  <RoleInfoCard />             {/* 角色基本信息 */}
  <PermissionTree />           {/* 功能权限树 */}
  <DataPermissionTable />      {/* 数据权限表格 */}
  <FieldPermissionTable />     {/* 字段权限表格 */}
  <UserListTable />            {/* 拥有此角色的用户列表 */}
</RoleDetail>
```

**核心功能**:
- 角色基本信息编辑
- 功能权限配置（树形选择）
- 数据权限配置（范围选择）
- 字段权限配置（表格编辑）
- 查看拥有此角色的用户

### 3. 部门管理页面 (Department)

**组件结构**:
```tsx
<Department>
  <DepartmentTree />           {/* 部门树 */}
  <MemberTable />              {/* 成员表格 */}
  <AddMemberModal />           {/* 添加成员弹窗 */}
</Department>
```

**核心功能**:
- 部门树形展示
- 创建/编辑/删除部门
- 添加/移除部门成员
- 设置部门负责人

---

## 💾 状态管理策略

### 使用场景

**需要全局状态管理**:
- ✅ 当前租户信息
- ✅ 当前用户权限
- ✅ 全局 Loading 状态

**使用本地 Hooks 即可**:
- ✅ 页面级数据（useQuota, useRole 等）
- ✅ 组件级状态（useState）

### 推荐方案：轻量级 Zustand

```typescript
// stores/enterprise.ts
import { create } from 'zustand';

interface EnterpriseStore {
  currentTenant: Tenant | null;
  setCurrentTenant: (tenant: Tenant) => void;
  permissions: string[];
  setPermissions: (perms: string[]) => void;
}

export const useEnterpriseStore = create<EnterpriseStore>((set) => ({
  currentTenant: null,
  setCurrentTenant: (tenant) => set({ currentTenant: tenant }),
  permissions: [],
  setPermissions: (perms) => set({ permissions: perms }),
}));
```

---

## 🎨 UI组件规范

### Semi Design 组件使用

**表格** (Table):
```tsx
import { Table } from '@coze-arch/coze-design';

<Table
  dataSource={data}
  columns={columns}
  loading={loading}
  pagination={{
    current: page,
    pageSize: pageSize,
    total: total,
    onChange: handlePageChange,
  }}
/>
```

**表单** (Form):
```tsx
import { Form, Input, Select, Button } from '@coze-arch/coze-design';

<Form onSubmit={handleSubmit}>
  <Form.Input field="name" label="名称" required />
  <Form.Select field="type" label="类型" required>
    <Select.Option value="type1">类型1</Select.Option>
  </Form.Select>
  <Button htmlType="submit">提交</Button>
</Form>
```

**弹窗** (Modal):
```tsx
import { Modal } from '@coze-arch/coze-design';

<Modal
  visible={visible}
  onCancel={onCancel}
  onOk={onOk}
  title="标题"
>
  {/* 内容 */}
</Modal>
```

### 样式规范

**使用 CSS Modules**:
```tsx
import s from './QuotaOverview.module.less';

<div className={s.container}>
  <div className={s.card}>
    {/* 内容 */}
  </div>
</div>
```

**使用 Tailwind CSS** (如果已配置):
```tsx
<div className="flex items-center gap-4 p-4 bg-white rounded">
  {/* 内容 */}
</div>
```

---

## 🧭 路由设计

### 路由配置

```typescript
// routes/enterprise.tsx
import { Routes, Route, Navigate } from 'react-router-dom';

export function EnterpriseRoutes() {
  return (
    <Routes>
      {/* 默认重定向到配额概览 */}
      <Route path="/" element={<Navigate to="/quota/overview" replace />} />

      {/* 配额管理 */}
      <Route path="/quota/*" element={<QuotaLayout />}>
        <Route path="overview" element={<QuotaOverview />} />
        <Route path=":resourceType" element={<QuotaDetail />} />
        <Route path="settings" element={<QuotaSettings />} />
      </Route>

      {/* 权限管理 */}
      <Route path="/permission/*" element={<PermissionLayout />}>
        <Route path="roles" element={<RoleList />} />
        <Route path="roles/:roleId" element={<RoleDetail />} />
        <Route path="departments" element={<Department />} />
        <Route path="matrix" element={<PermissionMatrix />} />
      </Route>

      {/* 租户管理（超级管理员） */}
      <Route path="/tenants/*" element={<TenantLayout />}>
        <Route path="" element={<TenantList />} />
        <Route path=":tenantId" element={<TenantDetail />} />
      </Route>
    </Routes>
  );
}
```

---

## 📅 实施计划

### 阶段划分

**阶段一：数据层（已完成）** ✅
- [x] API 层实现（tenantApi, quotaApi, permissionApi）
- [x] Hooks 层实现（useTenant, useQuota, usePermission）
- [x] TypeScript 类型定义

**阶段二：基础页面（待实现）**
- [ ] 配额概览页面
- [ ] 角色列表页面
- [ ] 路由配置

**阶段三：完整功能（待实现）**
- [ ] 配额详情与设置页面
- [ ] 角色详情与权限配置
- [ ] 部门管理页面

**阶段四：优化与完善（待实现）**
- [ ] 性能优化（虚拟滚动、懒加载）
- [ ] 用户体验优化（Loading 状态、错误提示）
- [ ] 响应式设计（移动端适配）

### 工作量估算

| 阶段 | 工作量 | 说明 |
|------|--------|------|
| 阶段一 | 3天 | API + Hooks（已完成） |
| 阶段二 | 5天 | 基础页面 + 路由 |
| 阶段三 | 8天 | 完整功能实现 |
| 阶段四 | 3天 | 优化与完善 |
| **总计** | **19天** | 约4周（1人月） |

---

## 📚 相关文档

- [前端 API 使用指南](../../../../../frontend/packages/arch/bot-api/ENTERPRISE_API.md)
- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)
- [全局一致性检查清单](./ZKER-全局一致性检查清单_v1.0.md)

---

**更新日志**：

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，完整的页面架构设计 | Claude AI |
