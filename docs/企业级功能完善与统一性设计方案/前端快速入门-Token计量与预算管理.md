# Token计量和预算管理 - 前端快速入门

## 快速开始

### 1. 安装依赖

```bash
cd frontend
rush update
```

### 2. 启动开发服务器

```bash
cd apps/coze-studio
npm run dev
```

### 3. 访问页面

- Token使用统计: http://localhost:8888/billing/token-usage
- 预算管理: http://localhost:8888/billing/budget

## API调用示例

### 记录Token使用

```typescript
import { getBillingService } from '@coze-studio/api-client';

async function recordUsage() {
  const service = await getBillingService();
  await service.recordTokenUsage('tenant-123', {
    user_id: 1,
    bot_id: 'bot-uuid',
    model_provider: 'openai',
    model_name: 'gpt-4',
    input_tokens: 1000,
    output_tokens: 500,
    total_tokens: 1500,
    request_type: 'chat',
  });
}
```

### 创建预算

```typescript
async function createBudget() {
  const service = await getBillingService();
  await service.createBudget('tenant-123', {
    budget_type: 'monthly',
    budget_amount: 1000,
    alert_threshold_1: 80,
    alert_threshold_2: 95,
    hard_cap_enabled: true,
    auto_downgrade_enabled: false,
    notification_channels: ['email'],
    notification_recipients: [],
  });
}
```

## 组件使用示例

### Token使用统计卡片

```tsx
import { useOverviewCards } from '@/pages/billing/hooks/useOverviewCards';
import { getBillingService } from '@coze-studio/api-client';

export function MyTokenStats() {
  const service = useMemo(() => getBillingService(), []);
  const { data, loading } = useOverviewCards(service, 'tenant-123');

  if (loading) return <Spin />;
  return (
    <div>
      <p>总Token: {data?.total_tokens}</p>
      <p>总成本: ${data?.total_cost}</p>
    </div>
  );
}
```

### 预算编辑弹窗

```tsx
import { BudgetEditModal } from '@/components/billing/BudgetEditModal';

export function MyBudgetEditor() {
  const [visible, setVisible] = useState(false);
  const [budget, setBudget] = useState(null);

  return (
    <>
      <Button onClick={() => setVisible(true)}>编辑预算</Button>
      <BudgetEditModal
        visible={visible}
        budget={budget}
        onSave={async (data) => {
          // 保存逻辑
          setVisible(false);
        }}
        onCancel={() => setVisible(false)}
      />
    </>
  );
}
```

## 常见问题

### Q: 如何切换语言？
A: 在设置中选择语言，或在代码中调用 `i18n.changeLanguage('en' | 'zh')`

### Q: 如何添加权限控制？
A: 使用 `usePermission` hook:
```tsx
const { hasPermission } = usePermission();
if (!hasPermission('billing:budget:view')) return <NoPermission />;
```

### Q: 如何自定义图表样式？
A: 修改ECharts option配置:
```tsx
const chartOption = {
  // 自定义颜色
  color: ['#1890ff', '#52c41a'],
  // 自定义配置...
};
```

## 更多文档

- [完整集成文档](./前端集成文档-Token计量与预算管理.md)
- [组件清单](./前端组件清单-Token计量与预算管理.md)
- [API文档](./API文档-Token计量与预算管理.md)
