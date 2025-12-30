# Token计量和预算管理 - 前端集成文档

## 概述

本文档描述了Token计量和预算管理功能的前端实现，包括所有组件、API调用、路由配置和使用示例。

## 组件清单

### 1. 类型定义
- **文件**: `frontend/common/types/billing.ts`
- **描述**: 所有Token计量和预算管理相关的TypeScript类型定义

### 2. API客户端服务
- **文件**: `frontend/packages/common/api-client/src/services/billing.service.ts`
- **描述**: 封装所有Token计量和预算管理API调用

### 3. 页面组件
- **Token使用统计页面**: `frontend/apps/coze-studio/src/pages/billing/TokenUsagePage.tsx`
- **预算管理页面**: `frontend/apps/coze-studio/src/pages/billing/BudgetManagementPage.tsx`

### 4. 业务组件
- **预算编辑弹窗**: `frontend/apps/coze-studio/src/components/billing/BudgetEditModal.tsx`

### 5. 自定义Hooks
- **useOverviewCards**: 使用概览卡片数据
- **useModelRanking**: 模型使用排行数据
- **useDailyTrend**: 每日使用趋势数据
- **useBotUsage**: Bot使用统计数据
- **useRealtimeRecords**: 实时Token记录数据

### 6. 国际化文件
- **中文**: `frontend/packages/common/i18n/locales/zh-CN.json`
- **英文**: `frontend/packages/common/i18n/locales/en-US.json`

## 路由配置

### 可用路由

```typescript
// Token使用统计页面
/billing/token-usage

// 预算管理页面
/billing/budget
```

### 路由定义位置

文件: `frontend/apps/coze-studio/src/router/routes.tsx`

```typescript
// 计费管理
{
  path: 'billing/token-usage',
  lazy: async () => {
    const { default: TokenUsagePage } = await import(
      '../pages/billing/TokenUsagePage'
    );
    return { Component: TokenUsagePage };
  },
},
{
  path: 'billing/budget',
  lazy: async () => {
    const { default: BudgetManagementPage } = await import(
      '../pages/billing/BudgetManagementPage'
    );
    return { Component: BudgetManagementPage };
  },
},
```

## API使用示例

### 1. Token计量API

```typescript
import { getBillingService } from '@coze-studio/api-client';

// 获取BillingService实例
const billingService = await getBillingService();

// 记录单次Token使用
const response = await billingService.recordTokenUsage(tenantId, {
  user_id: 123,
  bot_id: 'bot-uuid',
  model_provider: 'openai',
  model_name: 'gpt-4',
  input_tokens: 1000,
  output_tokens: 500,
  total_tokens: 1500,
  request_type: 'chat',
});

// 批量记录Token使用
const batchResponse = await billingService.batchRecordTokenUsage(tenantId, [
  { /* ... */ },
  { /* ... */ },
]);

// 获取使用统计
const stats = await billingService.getUsageStats(tenantId, {
  start_time: Date.now() - 7 * 24 * 60 * 60 * 1000,
  end_time: Date.now(),
  bot_id: 'bot-uuid',
});

// 获取每日使用统计
const dailyStats = await billingService.getDailyUsageStats(tenantId, 7);

// 获取模型使用统计
const modelStats = await billingService.getModelUsageStats(tenantId);
```

### 2. 预算管理API

```typescript
import { getBillingService } from '@coze-studio/api-client';

const billingService = await getBillingService();

// 获取预算设置
const budget = await billingService.getBudget(tenantId);

// 创建预算
const newBudget = await billingService.createBudget(tenantId, {
  budget_type: 'monthly',
  budget_amount: 1000,
  alert_threshold_1: 80,
  alert_threshold_2: 95,
  hard_cap_enabled: true,
  hard_cap_amount: 1200,
  auto_downgrade_enabled: false,
  notification_channels: ['email', 'sms'],
  notification_recipients: [
    {
      recipient_id: 'user-1',
      recipient_type: 'user',
      recipient_address: 'admin@example.com',
      channels: ['email'],
    },
  ],
});

// 更新预算
const updatedBudget = await billingService.updateBudget(tenantId, {
  budget_amount: 1500,
  alert_threshold_1: 85,
});

// 删除预算
await billingService.deleteBudget(tenantId);

// 获取预算使用情况
const budgetUsage = await billingService.getBudgetUsage(tenantId);

// 检查预算状态
const budgetCheck = await billingService.checkBudget(tenantId);

// 获取告警列表
const alerts = await billingService.getBudgetAlerts(tenantId, {
  alert_level: 'warning',
  limit: 20,
});
```

## 组件使用示例

### 1. Token使用统计页面

```tsx
import { TokenUsagePage } from '@/pages/billing/TokenUsagePage';

function App() {
  return <TokenUsagePage />;
}
```

### 2. 预算管理页面

```tsx
import { BudgetManagementPage } from '@/pages/billing/BudgetManagementPage';

function App() {
  return <BudgetManagementPage />;
}
```

### 3. 预算编辑弹窗

```tsx
import { BudgetEditModal } from '@/components/billing/BudgetEditModal';
import type { UpdateBudgetRequest } from '@coze-studio/common/types/billing';

function MyComponent() {
  const [visible, setVisible] = useState(false);
  const [budget, setBudget] = useState<BudgetSettingsDTO | null>(null);

  const handleSave = async (data: UpdateBudgetRequest) => {
    const billingService = await getBillingService();
    await billingService.updateBudget(tenantId, data);
    setVisible(false);
  };

  return (
    <>
      <Button onClick={() => setVisible(true)}>编辑预算</Button>
      <BudgetEditModal
        visible={visible}
        budget={budget}
        onSave={handleSave}
        onCancel={() => setVisible(false)}
      />
    </>
  );
}
```

### 4. 自定义Hooks使用

```tsx
import { useOverviewCards } from '@/pages/billing/hooks/useOverviewCards';
import { getBillingService } from '@coze-studio/api-client';

function MyComponent() {
  const billingService = useMemo(() => getBillingService(), []);
  const tenantId = 'current';

  const { data, loading, error, refetch } = useOverviewCards(billingService, tenantId);

  if (loading) return <Spin />;
  if (error) return <Alert message="加载失败" type="error" />;

  return (
    <div>
      <h1>总Token数: {data?.total_tokens}</h1>
      <h1>总成本: ${data?.total_cost}</h1>
      <Button onClick={refetch}>刷新</Button>
    </div>
  );
}
```

## 权限控制

### 集成RBAC权限系统

```tsx
import { usePermission } from '@coze-studio/auth';

function BudgetManagementPage() {
  const { hasPermission } = usePermission();

  // 页面级权限控制
  if (!hasPermission('billing:budget:view')) {
    return <NoPermission />;
  }

  return (
    <div>
      {/* 编辑按钮权限控制 */}
      {hasPermission('billing:budget:edit') && (
        <Button>编辑预算</Button>
      )}

      {/* 删除按钮权限控制 */}
      {hasPermission('billing:budget:delete') && (
        <Button danger>删除预算</Button>
      )}
    </div>
  );
}
```

## 国际化使用

### 使用翻译

```tsx
import { useTranslation } from '@coze-studio/i18n';

function MyComponent() {
  const { t } = useTranslation();

  return (
    <div>
      <h1>{t('billing.tokenUsage.title')}</h1>
      <p>{t('billing.tokenUsage.totalTokens')}: {totalTokens}</p>
      <p>{t('billing.tokenUsage.totalCost')}: ${totalCost}</p>
    </div>
  );
}
```

### 添加新翻译

在以下文件中添加新的翻译键：
- `frontend/packages/common/i18n/locales/zh-CN.json`
- `frontend/packages/common/i18n/locales/en-US.json`

```json
// zh-CN.json
{
  "billing": {
    "myNewKey": "我的新翻译"
  }
}

// en-US.json
{
  "billing": {
    "myNewKey": "My New Translation"
  }
}
```

## 样式自定义

### 使用Tailwind CSS

```tsx
function MyComponent() {
  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h1 className="text-2xl font-bold text-gray-900">
        预算管理
      </h1>
    </div>
  );
}
```

### 使用Semi Design主题

```tsx
import { ConfigProvider } from '@coze-studio/ui-components';

function App() {
  return (
    <ConfigProvider
      theme={{
        primaryColor: '#1890ff',
        borderRadius: 8,
      }}
    >
      <BudgetManagementPage />
    </ConfigProvider>
  );
}
```

## 测试

### 单元测试示例

```tsx
import { render, screen, waitFor } from '@testing-library/react';
import { TokenUsagePage } from '@/pages/billing/TokenUsagePage';

describe('TokenUsagePage', () => {
  it('should render overview cards', async () => {
    render(<TokenUsagePage />);

    await waitFor(() => {
      expect(screen.getByText('总Token数')).toBeInTheDocument();
      expect(screen.getByText('总成本')).toBeInTheDocument();
    });
  });

  it('should handle time range change', async () => {
    render(<TokenUsagePage />);

    const last30DaysButton = screen.getByText('最近30天');
    fireEvent.click(last30DaysButton);

    await waitFor(() => {
      // 验证数据刷新
    });
  });
});
```

### E2E测试示例

```typescript
describe('Budget Management', () => {
  it('should create budget', () => {
    cy.visit('/billing/budget');
    cy.contains('创建预算').click();
    cy.get('[data-cy="budget-type"]').select('monthly');
    cy.get('[data-cy="budget-amount"]').type('1000');
    cy.contains('保存').click();
    cy.contains('成功').should('be.visible');
  });
});
```

## 性能优化建议

### 1. 使用React.memo

```tsx
export const BudgetEditModal = React.memo<BudgetEditModalProps>(({ visible, budget, onSave, onCancel }) => {
  // ...
}, (prevProps, nextProps) => {
  return prevProps.visible === nextProps.visible &&
         prevProps.budget?.budget_id === nextProps.budget?.budget_id;
});
```

### 2. 使用useMemo缓存计算结果

```tsx
const costTrendOption = useMemo(() => {
  if (!dailyStats?.daily_stats) return {};

  const dates = dailyStats.daily_stats.map(item => item.date);
  const costs = dailyStats.daily_stats.map(item => item.total_cost);

  return { dates, costs };
}, [dailyStats]);
```

### 3. 使用useCallback缓存函数

```tsx
const handleRefresh = useCallback(() => {
  refetchStats();
  refetchModel();
  refetchDaily();
}, [refetchStats, refetchModel, refetchDaily]);
```

### 4. 懒加载ECharts

```tsx
const ReactECharts = lazy(() => import('echarts-for-react'));

function MyComponent() {
  return (
    <Suspense fallback={<Spin />}>
      <ReactECharts option={chartOption} />
    </Suspense>
  );
}
```

## 故障排查

### 常见问题

1. **API调用失败**
   - 检查租户ID是否正确
   - 检查API客户端是否已初始化
   - 检查网络请求是否被拦截器阻止

2. **页面不显示数据**
   - 检查后端API是否正常返回数据
   - 检查浏览器控制台是否有错误
   - 检查数据格式是否符合类型定义

3. **国际化不生效**
   - 检查语言包是否正确导入
   - 检查翻译键是否存在于语言包中
   - 检查i18n provider是否正确配置

### 调试技巧

```tsx
// 开启调试日志
useEffect(() => {
  console.log('Budget data:', budget);
  console.log('Budget usage:', budgetUsage);
}, [budget, budgetUsage]);

// 使用React DevTools检查组件状态
// 使用浏览器Network面板检查API请求
```

## 部署清单

### 前置条件

- [x] 后端API已部署并可访问
- [x] 数据库迁移已完成
- [x] 租户ID已配置

### 构建步骤

```bash
# 1. 安装依赖
cd frontend
rush update

# 2. 构建所有包
rush build

# 3. 构建主应用
cd apps/coze-studio
npm run build
```

### 环境变量配置

```bash
# .env.production
VITE_API_BASE_URL=https://api.example.com
VITE_TENANT_ID=current
```

### 部署验证

- [ ] 访问 `/billing/token-usage` 页面，确认能正常加载
- [ ] 访问 `/billing/budget` 页面，确认能正常加载
- [ ] 测试预算创建、编辑、删除功能
- [ ] 测试数据刷新功能
- [ ] 检查国际化切换是否正常
- [ ] 检查权限控制是否生效

## 总结

本文档提供了Token计量和预算管理功能的完整前端集成指南，包括：

- ✅ 完整的类型定义
- ✅ API客户端封装
- ✅ 页面组件实现
- ✅ 路由配置
- ✅ 国际化支持
- ✅ 权限控制
- ✅ 使用示例
- ✅ 性能优化建议
- ✅ 故障排查指南

所有组件已按照企业级开发规范实现，确保代码质量和可维护性。
