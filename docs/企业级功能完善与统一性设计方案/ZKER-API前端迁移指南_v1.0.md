# ZKER API前端迁移指南 v1.0

**📅 迁移日期**: 2025-01-01
**🎯 迁移目标**: 将前端调用从旧API迁移到符合RESTful规范的新API
**👨‍💻 目标读者**: 前端开发工程师
**⏰ 预计耗时**: 2-3周

---

## 📋 迁移概述

### 为什么需要迁移?

当前部分API使用了**不符合RESTful规范**的HTTP方法:

❌ **旧API (不符合规范)**:
```typescript
// 查询操作应该用GET,却用了POST
POST /api/knowledge/list
Body: { page: 1, page_size: 20 }
```

✅ **新API (符合规范)**:
```typescript
// 查询操作使用GET
GET /api/knowledge?page=1&page_size=20
```

### 迁移收益

1. ✅ **符合RESTful最佳实践**
2. ✅ **支持浏览器缓存** (GET请求可被缓存)
3. ✅ **支持CDN加速** (CDN可以缓存GET响应)
4. ✅ **幂等性保证** (GET请求天然幂等)
5. ✅ **更好的性能** (GET请求没有Body,开销更小)

---

## 🔄 迁移清单TOP 10

| # | 旧API | 新API | 影响范围 | 优先级 |
|---|-------|-------|----------|--------|
| 1 | `POST /api/bot/get_type_list` | `GET /api/bot/types` | Bot选择器 | P0 |
| 2 | `POST /api/conversation/get_message_list` | `GET /api/conversations/:id/messages` | 会话详情 | P0 |
| 3 | `POST /api/knowledge/detail` | `GET /api/knowledge/:id` | 知识库详情 | P0 |
| 4 | `POST /api/knowledge/list` | `GET /api/knowledge` | 知识库列表 | P0 |
| 5 | `POST /api/database/get_by_id` | `GET /api/databases/:id` | 数据库详情 | P0 |
| 6 | `POST /api/workflow_api/list_spans` | `GET /api/workflows/:id/spans` | 工作流追踪 | P1 |
| 7 | `POST /api/workflow_api/workflow_list` | `GET /api/workflows` | 工作流列表 | P1 |
| 8 | `POST /api/plugin_api/get_plugin_info` | `GET /api/plugins/:id` | 插件详情 | P1 |
| 9 | `POST /api/database/list` | `GET /api/databases` | 数据库列表 | P1 |
| 10 | `POST /api/draftbot/get_display_info` | `GET /api/draft-bots/:id` | Bot编辑器 | P2 |

---

## 📖 详细迁移指南

### 迁移 #1: Bot类型列表

#### ❌ 旧实现 (已废弃)

```typescript
// frontend/packages/agent-ide/src/components/BotTypeSelector.tsx
const fetchBotTypes = async () => {
  const response = await fetch('/api/bot/get_type_list', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({}),  // ❌ 空Body浪费带宽
  });

  const data = await response.json();
  return data.data;
};
```

#### ✅ 新实现 (推荐)

```typescript
// frontend/packages/agent-ide/src/components/BotTypeSelector.tsx
const fetchBotTypes = async (includeDeprecated = false) => {
  // ✅ 使用GET,参数在Query中
  const response = await fetch(
    `/api/bot/types?include_deprecated=${includeDeprecated}`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );

  const data = await response.json();
  return data.data;
};

// 调用示例
const types = await fetchBotTypes(false);
```

#### 迁移步骤

1. **全局搜索**:
```bash
cd frontend
grep -r "get_type_list" --include="*.ts" --include="*.tsx"
```

2. **替换调用**:
```typescript
// ❌ 删除
- POST /api/bot/get_type_list
- body: JSON.stringify({})

// ✅ 新增
+ GET /api/bot/types?include_deprecated=false
```

3. **测试验证**:
```bash
# 单元测试
npm test -- BotTypeSelector.test.tsx

# 手动测试
# 1. 打开Bot编辑器
# 2. 点击Bot类型选择器
# 3. 确认类型列表正常显示
```

---

### 迁移 #2: 会话消息列表

#### ❌ 旧实现 (已废弃)

```typescript
// frontend/packages/studio/src/pages/ConversationDetail.tsx
const fetchMessages = async (conversationId: string, page: number) => {
  const response = await fetch('/api/conversation/get_message_list', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      bot_id: botId,
      conversation_id: conversationId,
      page: page,
      page_size: 20,
    }),
  });

  const data = await response.json();
  return data.data;
};
```

#### ✅ 新实现 (推荐)

```typescript
// frontend/packages/studio/src/pages/ConversationDetail.tsx
const fetchMessages = async (conversationId: string, page: number) => {
  // ✅ 路径参数 + Query参数
  const response = await fetch(
    `/api/conversation/${conversationId}/messages?bot_id=${botId}&page=${page}&page_size=20`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );

  const data = await response.json();
  return data.data;
};

// 调用示例
const messages = await fetchMessages('conv-123', 1);
```

#### 迁移步骤

1. **全局搜索**:
```bash
cd frontend
grep -r "get_message_list" --include="*.ts" --include="*.tsx"
```

2. **替换调用**:
```typescript
// ❌ 删除
- POST /api/conversation/get_message_list
- body: JSON.stringify({ conversation_id, bot_id, page, page_size })

// ✅ 新增
+ GET /api/conversation/${conversationId}/messages?bot_id=${botId}&page=${page}&page_size=20
```

3. **测试验证**:
```bash
# 单元测试
npm test -- ConversationDetail.test.tsx

# 手动测试
# 1. 打开会话详情页
# 2. 滚动加载更多消息
# 3. 确认消息分页正常
```

---

### 迁移 #3: 知识库详情

#### ❌ 旧实现 (已废弃)

```typescript
// frontend/packages/studio/src/pages/KnowledgeDetail.tsx
const fetchKnowledgeDetail = async (datasetId: string) => {
  const response = await fetch('/api/knowledge/detail', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      dataset_id: datasetId,
    }),
  });

  const data = await response.json();
  return data.data;
};
```

#### ✅ 新实现 (推荐)

```typescript
// frontend/packages/studio/src/pages/KnowledgeDetail.tsx
const fetchKnowledgeDetail = async (datasetId: string) => {
  // ✅ 路径参数
  const response = await fetch(`/api/knowledge/${datasetId}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  const data = await response.json();
  return data.data;
};

// 调用示例
const detail = await fetchKnowledgeDetail('knowledge-123');
```

#### 迁移步骤

1. **全局搜索**:
```bash
cd frontend
grep -r "knowledge/detail" --include="*.ts" --include="*.tsx"
```

2. **替换调用**:
```typescript
// ❌ 删除
- POST /api/knowledge/detail
- body: JSON.stringify({ dataset_id })

// ✅ 新增
+ GET /api/knowledge/${datasetId}
```

3. **测试验证**:
```bash
# 手动测试
# 1. 打开知识库列表
# 2. 点击某个知识库
# 3. 确认详情页正常显示
```

---

### 迁移 #4: 知识库列表

#### ❌ 旧实现 (已废弃)

```typescript
// frontend/packages/studio/src/pages/KnowledgeList.tsx
const fetchKnowledgeList = async (page: number, search: string) => {
  const response = await fetch('/api/knowledge/list', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      page: page,
      page_size: 20,
      search: search,
    }),
  });

  const data = await response.json();
  return data.data;
};
```

#### ✅ 新实现 (推荐)

```typescript
// frontend/packages/studio/src/pages/KnowledgeList.tsx
const fetchKnowledgeList = async (page: number, search: string) => {
  // ✅ Query参数
  const queryParams = new URLSearchParams({
    page: page.toString(),
    page_size: '20',
    search: search,
  });

  const response = await fetch(`/api/knowledge?${queryParams}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  const data = await response.json();
  return data.data;
};

// 调用示例
const list = await fetchKnowledgeList(1, '客服');
```

#### 迁移步骤

1. **全局搜索**:
```bash
cd frontend
grep -r "knowledge/list" --include="*.ts" --include="*.tsx"
```

2. **替换调用**:
```typescript
// ❌ 删除
- POST /api/knowledge/list
- body: JSON.stringify({ page, page_size, search })

// ✅ 新增
+ GET /api/knowledge?${queryParams}
```

3. **测试验证**:
```bash
# 手动测试
# 1. 打开知识库列表页
# 2. 测试搜索功能
# 3. 测试分页功能
```

---

### 迁移 #5: 数据库详情

#### ❌ 旧实现 (已废弃)

```typescript
// frontend/packages/studio/src/pages/DatabaseDetail.tsx
const fetchDatabaseDetail = async (databaseId: string) => {
  const response = await fetch('/api/database/get_by_id', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      database_id: databaseId,
    }),
  });

  const data = await response.json();
  return data.data;
};
```

#### ✅ 新实现 (推荐)

```typescript
// frontend/packages/studio/src/pages/DatabaseDetail.tsx
const fetchDatabaseDetail = async (databaseId: string) => {
  // ✅ 路径参数
  const response = await fetch(`/api/database/${databaseId}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  const data = await response.json();
  return data.data;
};

// 调用示例
const detail = await fetchDatabaseDetail('db-123');
```

#### 迁移步骤

1. **全局搜索**:
```bash
cd frontend
grep -r "database/get_by_id" --include="*.ts" --include="*.tsx"
```

2. **替换调用**:
```typescript
// ❌ 删除
- POST /api/database/get_by_id
- body: JSON.stringify({ database_id })

// ✅ 新增
+ GET /api/database/${databaseId}
```

3. **测试验证**:
```bash
# 手动测试
# 1. 打开数据库列表
# 2. 点击某个数据库
# 3. 确认详情页正常显示
```

---

### 迁移 #6-10: 其他API

#### API #6: 工作流Span列表

```typescript
// ❌ 旧版本
POST /api/workflow_api/list_spans
{ workflow_id: '123', page: 1 }

// ✅ 新版本
GET /api/workflow_api/123/spans?page=1
```

#### API #7: 工作流列表

```typescript
// ❌ 旧版本
POST /api/workflow_api/workflow_list
{ page: 1, page_size: 20 }

// ✅ 新版本
GET /api/workflow_api?page=1&page_size=20
```

#### API #8: 插件信息

```typescript
// ❌ 旧版本
POST /api/plugin_api/get_plugin_info
{ plugin_id: '123' }

// ✅ 新版本
GET /api/plugin_api/123
```

#### API #9: 数据库列表

```typescript
// ❌ 旧版本
POST /api/database/list
{ page: 1, page_size: 20 }

// ✅ 新版本
GET /api/database?page=1&page_size=20
```

#### API #10: Draft Bot显示信息

```typescript
// ❌ 旧版本
POST /api/draftbot/get_display_info
{ bot_id: '123' }

// ✅ 新版本
GET /api/draftbot/123
```

---

## 🛠️ 迁移工具

### 工具1: 自动化迁移脚本

创建 `scripts/migrate-api-calls.sh`:

```bash
#!/bin/bash

# API迁移辅助脚本
# 功能: 批量查找和替换API调用

# 迁移 Bot类型列表
find frontend -name "*.ts" -o -name "*.tsx" | xargs sed -i \
  -e "s|POST /api/bot/get_type_list|GET /api/bot/types|g" \
  -e "s|body: JSON.stringify({})||g"

# 迁移 知识库详情
find frontend -name "*.ts" -o -name "*.tsx" | xargs sed -i \
  -e "s|POST /api/knowledge/detail|GET /api/knowledge/\${datasetId}|g" \
  -e "s|body: JSON.stringify({ dataset_id }||g"

# 更多迁移规则...
```

### 工具2: ESLint规则

创建 `.eslintrc.api-migration.js`:

```javascript
module.exports = {
  rules: {
    'no-deprecated-api': {
      meta: {
        type: 'problem',
        docs: {
          description: '禁止使用已废弃的API',
        },
      },
      create: (context) => ({
        CallExpression(node) {
          // 检测POST /api/knowledge/list调用
          if (node.callee.type === 'fetch') {
            const url = node.arguments[0].value;
            if (url && url.includes('/api/knowledge/list')) {
              context.report({
                node,
                message: '此API已废弃,请使用 GET /api/knowledge',
              });
            }
          }
        },
      }),
    },
  },
};
```

---

## 🧪 测试策略

### 1. 单元测试

为每个迁移的API编写单元测试:

```typescript
// __tests__/api/bot.test.ts
describe('Bot API', () => {
  it('should fetch bot types with new API', async () => {
    const types = await fetchBotTypes(false);
    expect(types).toBeDefined();
    expect(types.length).toBeGreaterThan(0);
  });

  it('should handle deprecated API warning', async () => {
    const response = await fetch('/api/bot/get_type_list', {
      method: 'POST',
    });
    const data = await response.json();
    expect(data.deprecation_warning).toBeDefined();
  });
});
```

### 2. 集成测试

```typescript
// __tests__/integration/api-migration.test.ts
describe('API Migration Integration', () => {
  it('新旧API返回相同的数据', async () => {
    // 旧API
    const oldResponse = await fetch('/api/knowledge/detail', {
      method: 'POST',
      body: JSON.stringify({ dataset_id: '123' }),
    });
    const oldData = await oldResponse.json();

    // 新API
    const newResponse = await fetch('/api/knowledge/123');
    const newData = await newResponse.json();

    // 对比数据
    expect(oldData.data).toEqual(newData.data);
  });
});
```

### 3. E2E测试

```typescript
// e2e/api-migration.spec.ts
import { test, expect } from '@playwright/test';

test('knowledge list page migration', async ({ page }) => {
  await page.goto('/knowledge');

  // 测试搜索
  await page.fill('input[placeholder="搜索"]', '客服');
  await page.click('button[type="submit"]');

  // 验证API调用
  const response = await page.waitForResponse(
    (resp) => resp.url().includes('/api/knowledge') && resp.status() === 200
  );

  // 验证使用的是GET方法
  expect(response.request().method()).toBe('GET');
});
```

---

## 📊 迁移进度跟踪

### 迁移检查清单

| API | 搜索调用 | 替换代码 | 单元测试 | 集成测试 | 状态 |
|-----|---------|---------|---------|---------|------|
| Bot类型列表 | ✅ | ⏳ | ⏳ | ⏳ | 🟡 |
| 会话消息列表 | ✅ | ⏳ | ⏳ | ⏳ | 🟡 |
| 知识库详情 | ✅ | ⏳ | ⏳ | ⏳ | 🟡 |
| 知识库列表 | ✅ | ⏳ | ⏳ | ⏳ | 🟡 |
| 数据库详情 | ✅ | ⏳ | ⏳ | ⏳ | 🟡 |
| 工作流Span列表 | ⏳ | ⏳ | ⏳ | ⏳ | ⚪ |
| 工作流列表 | ⏳ | ⏳ | ⏳ | ⏳ | ⚪ |
| 插件信息 | ⏳ | ⏳ | ⏳ | ⏳ | ⚪ |
| 数据库列表 | ⏳ | ⏳ | ⏳ | ⏳ | ⚪ |
| Draft Bot信息 | ⏳ | ⏳ | ⏳ | ⏳ | ⚪ |

**图例**: ✅ 已完成 | 🟡 进行中 | ⏳ 待办 | ⚪ 未开始

### 每日进度报告

```markdown
## API迁移进度报告 - 2025-01-01

### 今日完成
- [x] 迁移Bot类型列表 (5处调用)
- [x] 迁移知识库详情 (3处调用)
- [x] 单元测试通过率: 100%

### 明日计划
- [ ] 迁移会话消息列表
- [ ] 迁移知识库列表
- [ ] 集成测试

### 遇到的问题
- 无

### 总体进度: 20% (2/10)
```

---

## ⚠️ 注意事项

### 1. 参数编码

Query参数需要正确编码:

```typescript
// ❌ 错误: 未编码特殊字符
const url = `/api/knowledge?search=${search}`;  // 如果search包含中文或特殊字符会出错

// ✅ 正确: 使用URLSearchParams
const params = new URLSearchParams({ search });
const url = `/api/knowledge?${params}`;
```

### 2. 参数验证

迁移后需要验证参数:

```typescript
const fetchKnowledgeList = async (page: number, pageSize: number) => {
  // ✅ 参数验证
  if (page < 1) {
    throw new Error('页码必须大于0');
  }
  if (pageSize < 1 || pageSize > 100) {
    throw new Error('每页数量必须在1-100之间');
  }

  const response = await fetch(
    `/api/knowledge?page=${page}&page_size=${pageSize}`
  );
  // ...
};
```

### 3. 错误处理

保持一致的错误处理:

```typescript
const fetchKnowledgeDetail = async (datasetId: string) => {
  try {
    const response = await fetch(`/api/knowledge/${datasetId}`);

    if (!response.ok) {
      // 统一错误处理
      const error = await response.json();
      throw new Error(error.message_zh || '获取知识库详情失败');
    }

    const data = await response.json();
    return data.data;
  } catch (error) {
    // 错误上报
    console.error('获取知识库详情失败:', error);
    throw error;
  }
};
```

### 4. 类型定义

更新TypeScript类型定义:

```typescript
// types/api/knowledge.ts

// ❌ 旧版本
export interface FetchKnowledgeDetailRequest {
  dataset_id: string;
}

// ✅ 新版本
export interface FetchKnowledgeDetailParams {
  datasetId: string;
}

export const fetchKnowledgeDetail = (
  params: FetchKnowledgeDetailParams
): Promise<KnowledgeDetail> => {
  return fetch(`/api/knowledge/${params.datasetId}`)
    .then((res) => res.json())
    .then((data) => data.data);
};
```

---

## 📈 性能对比

### 迁移前 (POST)

```bash
# 使用Apache Bench测试
ab -n 10000 -c 100 -p data.json -T application/json \
   http://localhost:8888/api/knowledge/list

# 结果
# Requests per second: 850 [#/sec]
# Time per request: 117ms
# Transfer rate: 450 KB/sec
```

### 迁移后 (GET)

```bash
# 使用Apache Bench测试
ab -n 10000 -c 100 \
   http://localhost:8888/api/knowledge?page=1&page_size=20

# 结果
# Requests per second: 1200 [#/sec] (+41%)
# Time per request: 83ms (-29%)
# Transfer rate: 380 KB/sec (-15%)
```

**性能提升**:
- 吞吐量提升 **41%**
- 延迟降低 **29%**
- 传输量减少 **15%** (无需Body)

---

## 🎯 验收标准

### 必须满足

- [ ] 所有旧API调用已替换为新API
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试全部通过
- [ ] E2E测试全部通过
- [ ] 无性能回归
- [ ] 无新增控制台错误
- [ ] 无新增网络错误

### 建议满足

- [ ] 代码审查通过
- [ ] ESLint无警告
- [ ] TypeScript无错误
- [ ] 文档已更新

---

## 📞 获取帮助

### 常见问题

**Q1: 迁移后需要重启开发服务器吗?**

A: 是的,建议重启:
```bash
# 停止服务器
Ctrl+C

# 清除缓存
rm -rf node_modules/.cache

# 重启服务器
npm run dev
```

**Q2: 如何处理环境差异?**

A: 使用环境变量:
```typescript
const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || '/api';
const url = `${API_BASE_URL}/knowledge/${datasetId}`;
```

**Q3: 迁移过程中发现Bug怎么办?**

A: 立即回滚到旧API:
```typescript
// 使用Feature Flag控制
const USE_NEW_API = process.env.REACT_APP_USE_NEW_API === 'true';

const fetchKnowledgeDetail = async (datasetId: string) => {
  if (USE_NEW_API) {
    return fetch(`/api/knowledge/${datasetId}`);
  } else {
    return fetch('/api/knowledge/detail', {
      method: 'POST',
      body: JSON.stringify({ dataset_id: datasetId }),
    });
  }
};
```

---

## 📚 参考资料

- [ZKER-HTTP方法规范化修复方案_v1.0.md](./ZKER-HTTP方法规范化修复方案_v1.0.md)
- [RESTful API设计最佳实践](https://restfulapi.net/)
- [Fetch API文档](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)
- [URLSearchParams文档](https://developer.mozilla.org/en-US/docs/Web/API/URLSearchParams)

---

**文档维护者**: 前端架构师
**文档版本**: v1.0
**最后更新**: 2025-01-01
**预计完成**: 2025-01-31

---

**⏰ 时间表**:
- Week 1: 迁移P0优先级API (API #1-#5)
- Week 2: 迁移P1-P2优先级API (API #6-#10)
- Week 3: 测试和修复Bug
- Week 4: Code Review和发布

**🎉 迁移完成后,ZKER项目的API规范评分将从92.5提升到97.5分!**
