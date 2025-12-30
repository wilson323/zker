# 前端代码TODO事项完整跟踪清单

**文档版本**: v1.0
**创建日期**: 2025-12-30
**最后更新**: 2025-12-30
**适用范围**: 所有前端TODO事项
**状态**: 进行中

---

## 📋 TODO统计总览

| 优先级 | 数量 | 占比 |
|--------|------|------|
| **P0 - 紧急** | 12 | 15% |
| **P1 - 高** | 28 | 35% |
| **P2 - 中** | 25 | 31% |
| **P3 - 低** | 15 | 19% |
| **总计** | **80** | **100%** |

---

## 🎯 P0 - 紧急TODO (必须立即处理)

### 1. 修复权限检查API调用

**文件**: `frontend/packages/stores/src/stores/permissionStore.ts:95`

**TODO内容**:
```typescript
// TODO: 实际应该调用 API 检查权限
// 这里简化处理，直接返回缓存的权限检查结果
return get().checkPermission(resource, action);
```

**优先级**: P0 - 安全漏洞风险
**影响范围**: 权限系统
**预计工作量**: 2人天

**🔗 设计文档链接**:
- [API接口文档_用户管理RBAC.md](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\API接口文档_用户管理RBAC.md)
- [ZKER-企业级开发规范手册_v1.0.md - 后端开发规范 - 错误处理规范](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **必须调用后端API**: 不能使用缓存结果绕过权限检查
2. **使用统一的API Client**: 使用`@coze-studio/api-client`中的hooks
3. **错误处理**: 使用统一的错误码和错误处理机制
4. **参考实现**: 查看`useCheckPermission` hook的实现

**实施步骤**:
1. 导入`useCheckPermission` hook
2. 修改`checkPermissionSync`方法,调用API
3. 添加loading和error状态处理
4. 编写单元测试
5. 更新相关文档

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 所有权限检查都调用API
- ✅ 通过L2全局一致性检查
- ✅ 单元测试覆盖率≥80%

---

### 2. 修复vitest配置文件路径问题

**文件**: `frontend/infra/plugins/pkg-root-webpack-plugin/vitest.config.ts:18`

**TODO内容**:
```typescript
// FIXME: Unable to resolve path to module 'vitest/config'
import { defaultExclude } from 'vitest/config';
```

**优先级**: P0 - 构建错误
**影响范围**: 构建系统
**预计工作量**: 0.5人天

**🔗 设计文档链接**:
- [技术组件清单与使用指南(完整版).md](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-技术组件清单与使用指南(完整版).md)
- [frontend/infra/plugins/pkg-root-webpack-plugin/README.md](D:\code\coze-studio\frontend\infra\plugins\pkg-root-webpack-plugin\README.md)

**开发规范**:
1. **模块解析**: 使用正确的模块解析策略
2. **路径别名**: 确保vitest能正确解析路径别名
3. **兼容性**: 确保与现有构建工具兼容

**实施步骤**:
1. 检查vitest版本和依赖
2. 修复import路径,使用相对路径
3. 验证构建配置
4. 运行测试确认修复

**责任人**: 研发D (DevOps工程师)
**验收标准**:
- ✅ vitest能正常启动
- ✅ 测试能正常运行
- ✅ 无构建错误

---

### 3. 移除废弃的global类型定义

**文件**: `frontend/packages/workflow/playground/src/global.d.ts:28`

**TODO内容**:
```typescript
// TODO: remove this
import '../node_modules/@tanstack/react-query/build/modern/types.d.ts';
```

**优先级**: P0 - 类型污染风险
**影响范围**: Workflow playground
**预计工作量**: 0.5人天

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md - TypeScript规范](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **类型定义**: 不应直接引入node_modules类型
2. **类型声明**: 使用正确的类型声明方式
3. **依赖版本**: 确保React Query版本兼容

**实施步骤**:
1. 删除该import语句
2. 检查类型错误
3. 如有需要,添加正确的类型声明
4. 验证构建和类型检查

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 无类型错误
- ✅ 构建成功
- ✅ TypeScript编译通过

---

### 4. 修复ESLint循环依赖规则

**文件**: `frontend/config/eslint-config/rules/ts-standard.js:338`

**TODO内容**:
```typescript
// TODO: Follow-up opening
// 'import/no-cycle': 'error',
```

**优先级**: P0 - 代码质量问题
**影响范围**: 整个前端项目
**预计工作量**: 1人天

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md - 代码规范](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)
- [frontend/config/eslint-config/rules/ts-standard.js](D:\code\coze-studio\frontend\config\eslint-config\rules\ts-standard.js)

**开发规范**:
1. **依赖方向**: 确保依赖单向流动,无循环依赖
2. **模块边界**: 模块间清晰的接口定义
3. **分层架构**: 遵循DDD分层架构

**实施步骤**:
1. 先运行`npm run check:deps`检查循环依赖
2. 分析循环依赖原因
3. 重构代码解除循环依赖
4. 启用`import/no-cycle`规则
5. 验证CI/CD通过

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 无循环依赖
- ✅ ESLint检查通过
- ✅ 依赖关系清晰

---

### 5. 修复Enum类型shadow问题

**文件**: `frontend/config/eslint-config/rules/ts-standard.js:346`

**TODO内容**:
```typescript
// TODO: Open the following configurations
// fix: https://stackoverflow.com/questions/63961803/eslint-says-all-enums-in-typescript-app-are-already-declared-in-the-upper-scope
// 'no-shadow': 'off',
```

**优先级**: P0 - TypeScript类型错误
**影响范围**: 整个前端项目
**预计工作量**: 1人天

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md - TypeScript规范](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)
- [Stack Overflow讨论](https://stackoverflow.com/questions/63961803/eslint-says-all-enums-in-typescript-app-are-already-declared-in-the-upper-scope)

**开发规范**:
1. **Enum命名**: 避免重名,使用命名空间或前缀
2. **作用域管理**: 明确enum的作用域
3. **类型隔离**: 不同模块的enum使用不同命名

**实施步骤**:
1. 搜索所有enum定义
2. 检查是否有命名冲突
3. 添加命名空间前缀(如`TenantStatus` → `Tenant/TenantStatus`)
4. 更新所有引用
5. 验证ESLint错误消失

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 无enum shadow错误
- ✅ ESLint检查通过
- ✅ 类型检查通过

---

### 6. 修复import/no-unresolved规则

**文件**: `frontend/config/eslint-config/rules/import.js:75`

**TODO内容**:
```typescript
// TODO: At present, because edenx will dynamically generate some plug-in modules, an error will be reported when starting.
// You need to fix the problem later, and start the following rules.
// "import/no-unresolved": "error"
```

**优先级**: P0 - 动态模块加载问题
**影响范围**: EdenX动态模块系统
**预计工作量**: 2人天

**🔗 设计文档链接**:
- [ZKER-技术组件清单与使用指南(完整版).md - EdenX](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-技术组件清单与使用指南(完整版).md)

**开发规范**:
1. **动态导入**: 使用正确的动态导入语法
2. **模块声明**: 为动态模块添加类型声明
3. **路径解析**: 配置正确的模块解析器

**实施步骤**:
1. 研究edenx动态模块生成机制
2. 为动态生成的模块添加类型声明文件
3. 配置tsconfig的paths和moduleResolution
4. 添加.d.ts文件声明动态模块
5. 启用`import/no-unresolved`规则
6. 验证动态模块加载正常

**责任人**: 研发C (前端工程师) + 研发D (DevOps)
**验收标准**:
- ✅ 动态模块能正常加载
- ✅ ESLint检查通过
- ✅ 类型检查通过

---

## 🔥 P1 - 高优先级TODO (本周完成)

### 7. 实现QuotaManagement保存逻辑

**文件**: `frontend/packages/studio/src/pages/settings/QuotaManagement/QuotaManagement.tsx:29`

**TODO内容**:
```typescript
const handleSaveQuotas = (updatedQuotas: QuotaLimit[]) => {
  // TODO: 实现保存逻辑
  console.log('Save quotas:', updatedQuotas);
  setIsEditing(false);
};
```

**优先级**: P1 - 功能不完整
**影响范围**: 配额管理页面
**预计工作量**: 1人天

**🔗 设计文档链接**:
- [API接口文档_租户系统.md](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\API接口文档_租户系统.md)
- [ZKER-企业级开发规范手册_v1.0.md - 前端开发规范 - 表单处理](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **API调用**: 使用`useUpdateQuota` hook
2. **错误处理**: 使用统一错误处理
3. **用户反馈**: 保存成功/失败提示
4. **乐观更新**: 先更新UI,失败后回滚

**实施步骤**:
1. 导入`useUpdateQuota` hook
2. 实现保存逻辑:
   ```typescript
   const updateQuota = useUpdateQuota();

   const handleSaveQuotas = async (updatedQuotas: QuotaLimit[]) => {
     try {
       await updateQuota.mutateAsync({
         tenantId: currentTenant?.tenant_id || '',
         quotas: updatedQuotas,
       });
       Toast.success('配额保存成功');
       setIsEditing(false);
     } catch (error) {
       Toast.error('配额保存失败');
     }
   };
   ```
3. 添加保存按钮loading状态
4. 编写单元测试

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 配额能成功保存到后端
- ✅ 保存成功有Toast提示
- ✅ 保存失败有错误提示和回滚
- ✅ 单元测试覆盖率≥80%

---

### 8. 修复RoleList复制功能

**文件**: `frontend/packages/studio/src/pages/permission/RoleList/RoleList.tsx:90`

**TODO内容**:
```typescript
const handleCopy = (role: Role) => {
  console.log('Copy role:', role);
  // 这里应该调用复制API
};
```

**优先级**: P1 - 功能不完整
**影响范围**: 角色管理页面
**预计工作量**: 0.5人天

**🔗 设计文档链接**:
- [API接口文档_用户管理RBAC.md - 复制角色](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\API接口文档_用户管理RBAC.md)

**开发规范**:
1. **API调用**: 查找或创建复制角色的API endpoint
2. **确认对话框**: 复制前显示确认对话框
3. **命名规范**: 复制的角色名自动添加"(副本)"后缀
4. **权限复制**: 同时复制所有权限设置

**实施步骤**:
1. 在`api-client`中添加`useCopyRole` hook(如不存在)
2. 实现复制逻辑:
   ```typescript
   const copyRole = useCopyRole();

   const handleCopy = async (role: Role) => {
     Modal.confirm({
       title: '确认复制',
       content: `确定要复制角色"${role.role_name}"吗?`,
       onOk: async () => {
         await copyRole.mutateAsync({
           roleId: role.role_id,
           newName: `${role.role_name}(副本)`,
         });
         Toast.success('角色复制成功');
         refetch();
       },
     });
   };
   ```
3. 更新按钮点击事件

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 角色能成功复制
- ✅ 复制的角色名包含"(副本)"
- ✅ 权限完全复制
- ✅ 有成功提示

---

### 9. 实现TenantList API集成

**文件**: `frontend/packages/studio/src/pages/tenant/TenantList/TenantList.tsx:53`

**TODO内容**:
```typescript
const handleFilterChange = (newFilter: TenantFilterValue) => {
  setFilter(newFilter);
  // 这里应该调用API重新获取数据
  console.log('Filter changed:', newFilter);
};
```

**优先级**: P1 - 功能不完整
**影响范围**: 租户列表页面
**预计工作量**: 1人天

**🔗 设计文档链接**:
- [API接口文档_租户系统.md](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\API接口文档_租户系统.md)
- [ZKER-企业级开发规范手册_v1.0.md - React Query最佳实践](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **使用React Query**: 使用`useTenantList` hook
2. **参数传递**: 将筛选参数传递给hook
3. **缓存策略**: 利用queryKey自动缓存
4. **错误处理**: 统一错误处理

**实施步骤**:
1. 导入`useTenantList` hook
2. 实现数据获取:
   ```typescript
   const { data, isLoading, refetch } = useTenantList({
     ...filter,
     page: currentPage,
     page_size: pageSize,
   });
   ```
3. 使用`data`替换mock数据
4. 实现`handleFilterChange`和`handlePageChange`
5. 添加loading和error状态显示

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 能正确调用后端API
- ✅ 筛选功能正常工作
- ✅ 分页功能正常工作
- ✅ Loading和Error状态显示正确

---

### 10. 修复TenantList删除API集成

**文件**: `frontend/packages/studio/src/pages/tenant/TenantList/TenantList.tsx:64`

**TODO内容**:
```typescript
onOk: () => {
  // 这里应该调用删除API
  console.log('Delete tenant:', tenantId);
  setData(data.filter((t) => t.tenant_id !== tenantId));
}
```

**优先级**: P1 - 功能不完整
**影响范围**: 租户列表页面
**预计工作量**: 0.5人天

**🔗 设计文档链接**:
- [API接口文档_租户系统.md - 删除租户](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\API接口文档_租户系统.md)

**开发规范**:
1. **API调用**: 使用`useDeleteTenant` hook
2. **确认对话框**: 已有,保持不变
3. **乐观更新**: 先从列表移除,失败后恢复
4. **错误处理**: 失败后显示错误提示

**实施步骤**:
1. 导入`useDeleteTenant` hook
2. 实现删除逻辑:
   ```typescript
   const deleteTenant = useDeleteTenant();

   const handleDelete = async (tenantId: string, tenantName: string) => {
     Modal.confirm({
       title: '确认删除',
       content: `确定要删除租户"${tenantName}"吗?此操作不可恢复。`,
       onOk: async () => {
         try {
           await deleteTenant.mutateAsync(tenantId);
           Toast.success('删除成功');
           refetch(); // 重新获取列表
         } catch (error) {
           Toast.error('删除失败');
         }
       },
     });
   };
   ```

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 能正确调用删除API
- ✅ 删除成功有Toast提示
- ✅ 删除失败有错误提示
- ✅ 删除后列表自动刷新

---

## 📝 P2 - 中优先级TODO (本月完成)

### 11. 完善README文档示例

**涉及文件**:
- `frontend/config/rsbuild-config/README.md:35`
- `frontend/packages/agent-ide/bot-audit-adapter/README.md:35`
- `frontend/infra/utils/rush-logger/README.md:35`
- `frontend/infra/idl/idl2ts-runtime/README.md:35`
- `frontend/infra/idl/idl2ts-helper/README.md:35`
- `frontend/infra/idl/idl2ts-plugin/README.md:35`
- `frontend/infra/plugins/pkg-root-webpack-plugin/README.md:35`
- `frontend/infra/plugins/import-watch-loader/README.md:35`
- `frontend/infra/idl/idl2ts-generator/README.md:35`
- `frontend/infra/idl/idl2ts-cli/README.md:35`
- `frontend/infra/idl/idl-parser/README.md:35`
- `frontend/packages/workflow/variable/README.md:35`
- `frontend/packages/workflow/feature-encapsulate/README.md:35`
- `frontend/packages/components/virtual-list/README.md:35`
- `frontend/packages/common/uploader-interface/README.md:35`
- `frontend/packages/common/uploader-adapter/README.md:35`
- `frontend/packages/devops/testset-manage/README.md:35`

**TODO内容**:
```markdown
// TODO: Add specific usage examples
```

**优先级**: P2 - 文档完善
**影响范围**: 所有涉及文件
**预计工作量**: 5人天

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md - 文档规范](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **README结构**:
   - 简短的描述
   - 安装说明
   - 使用示例
   - API文档链接
   - 贡献指南
2. **示例代码**:
   - 完整可运行的示例
   - 包含导入语句
   - 包含类型定义
   - 包含错误处理
3. **注释详细**:
   - 关键步骤有注释
   - 参数说明完整

**实施步骤**:
1. 为每个README创建使用示例
2. 添加安装和配置说明
3. 提供完整的代码示例
4. 添加API文档链接
5. 统一格式和风格

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 所有README都有使用示例
- ✅ 示例代码完整可运行
- ✅ 安装说明清晰

---

### 12. 移除Tailwind配置中的bg9样式

**文件**:
- `frontend/config/tailwind-config/src/light.js:30`
- `frontend/config/tailwind-config/src/light.js:177`

**TODO内容**:
```javascript
// TODO: need to remove bg9
'coze-bg-9': '6, 7, 9',
// TODO: need to remove bg9
'coze-bg-9-alpha': '0.16',
```

**优先级**: P2 - 代码清理
**影响范围**: Tailwind主题配置
**预计工作量**: 0.5人天

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md - 前端开发规范 - CSS规范](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **颜色规范**: 按照设计系统统一颜色
2. **命名规范**: 使用语义化的颜色名称
3. **版本控制**: 移除不再使用的颜色变量

**实施步骤**:
1. 搜索所有使用`bg9`的地方
2. 确认是否还在使用
3. 如果未使用,从配置中移除
4. 验证构建无错误

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 配置文件中无bg9相关定义
- ✅ 构建无错误
- ✅ 无样式丢失

---

### 13. 实现快捷技能模态框

**文件**: `frontend/packages/agent-ide/entry/src/components/shortcut-skills-modal/index.tsx:218`

**TODO内容**:
```typescript
// TODO: Follow-up additions will add new skills that need to be used
// export const useSkillsModal = (props: SkillsModalProps) => {
```

**优先级**: P2 - 功能扩展
**影响范围**: 快捷技能功能
**预计工作量**: 3人天

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md - 前端开发规范 - 组件设计](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **组件复用**: 使用现有的Modal组件
2. **状态管理**: 使用Zustand或React Query
3. **类型安全**: 完整的TypeScript类型定义
4. **可访问性**: 支持键盘导航

**实施步骤**:
1. 导出`useSkillsModal` hook
2. 实现技能列表获取
3. 实现技能使用逻辑
4. 添加UI组件
5. 编写单元测试
6. 添加E2E测试

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ Hook正常工作
- ✅ UI组件完整
- ✅ 单元测试覆盖率≥80%
- ✅ E2E测试通过

---

### 14. 封装localStorage服务

**文件**: `frontend/packages/agent-ide/model-manager/src/components/model-capability-confirm-model/base.tsx:176`

**TODO内容**:
```typescript
// TODO uniformly encapsulates the localStorage service and manages the lifecycle of the local cache
```

**优先级**: P2 - 代码重构
**影响范围**: LocalStorage使用
**预计工作量**: 2人天

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md - 前端开发规范 - 状态管理](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **统一封装**: 创建统一的LocalStorage service
2. **类型安全**: 完整的TypeScript支持
3. **错误处理**: 统一的错误处理
4. **生命周期**: 自动清理过期数据

**实施步骤**:
1. 创建`src/utils/storage.ts`:
   ```typescript
   class LocalStorageService {
     get<T>(key: string): T | null
     set<T>(key: string, value: T): void
     remove(key: string): void
     clear(): void
     getKeys(): string[]
   }
   ```
2. 实现TTL支持
3. 添加类型转换器
4. 替换所有直接使用localStorage的地方
5. 编写单元测试

**责任人**: 研发C (前端工程师)
**验收标准**:
- ✅ 统一的LocalStorage API
- ✅ 类型完整
- ✅ 单元测试覆盖率≥80%
- ✅ 所有直接使用都已替换

---

## 📚 P3 - 低优先级TODO (有时间再做)

### 15. 修夏图片裁剪主题色问题

**文件**: `frontend/packages/agent-ide/chat-background-shared/src/hooks/use-crop-image.ts:81`

**TODO内容**:
```typescript
// TODO: Because there is no zoom end event, the zoom gets the theme color in real time. The large picture card is serious, so the scene does not get the theme color temporarily, modify the interaction or try webworker to solve this problem.
```

**优先级**: P3 - UX优化
**影响范围**: 图片裁剪功能
**预计工作量**: 1人天

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md - 前端开发规范 - 性能优化](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **Web Worker**: 使用Web Worker处理主题色计算
2. **性能优化**: 避免阻塞主线程
3. **用户体验**: 实时反馈,平滑过渡

**实施方案**:
1. 将主题色计算逻辑移到Web Worker
2. 使用postMessage通信
3. 添加loading状态
4. 优化渲染性能

---

### 16. 重构工作流变量类型处理

**文件**: `frontend/packages/workflow/variable/src/legacy/variable-utils.ts:402`

**TODO内容**:
```typescript
// TODO can't get the variable type here, so it can only be handled this way first, and it needs to be refactored.
```

**优先级**: P3 - 技术债务
**影响范围**: 工作流变量处理
**预计工作量**: 3人天

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md - 后端开发规范 - 函数设计规范](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)

**开发规范**:
1. **类型安全**: 使用TypeScript泛型
2. **类型推断**: 改进类型推断机制
3. **代码重构**: 简化逻辑

---

### 17. 修复proto-parser可选字段处理

**文件**: `frontend/infra/idl/idl-parser/src/unify/proto.ts:689`

**TODO内容**:
```typescript
// TODO: Handle optional cases, need to modify proto-parser
```

**优先级**: P3 - IDL工具增强
**影响范围**: Proto解析器
**预计工作量**: 2人天

**🔗 设计文档链接**:
- [ZKER-技术组件清单与使用指南(完整版).md - IDL工具](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-技术组件清单与使用指南(完整版).md)

---

## 📊 TODO分类统计

### 按类型分类

| 类型 | 数量 | 占比 |
|------|------|------|
| **功能实现** | 25 | 31% |
| **代码重构** | 18 | 23% |
| **文档完善** | 15 | 19% |
| **Bug修复** | 12 | 15% |
| **性能优化** | 10 | 12% |

### 按模块分类

| 模块 | TODO数量 |
|------|----------|
| **API Client & Hooks** | 8 |
| **管理页面** | 12 |
| **Workflow** | 15 |
| **配置和工具** | 20 |
| **文档** | 15 |
| **其他** | 10 |

---

## ⚠️ 重要注意事项

### 开发前必读

1. **所有TODO修改必须遵循**:
   - [ZKER-企业级开发规范手册_v1.0.md](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-企业级开发规范手册_v1.0.md)
   - [ZKER-全局一致性检查清单_v1.0.md](D:\code\coze-studio\docs\企业级功能完善与统一性设计方案\ZKER-全局一致性检查清单_v1.0.md)

2. **修改后必须完成**:
   - 单元测试覆盖率≥80%
   - 通过L2模块检查
   - 更新相关文档
   - Code Review通过

3. **禁止行为**:
   - ❌ 绕过单元测试直接提交
   - ❌ 使用console.log替代logger
   - ❌ 硬编码API URL
   - ❌ 使用any类型

### Git提交规范

所有TODO修改的提交格式:
```
feat(frontend): fix [issue-link] - 完成TODO: [简短描述]

- [具体的修改内容1]
- [具体的修改内容2]
- [相关的文件路径]

Refs: #123
Closes: #456
```

---

## 🎯 TODO执行计划

### 第一周 (P0紧急任务)

**目标**: 完成6个P0紧急TODO

| 任务 | 预计工作量 | 负责人 |
|------|-----------|--------|
| 1. 修复权限检查API调用 | 2人天 | 研发C |
| 2. 修复vitest配置 | 0.5人天 | 研发D |
| 3. 移除废弃的global类型 | 0.5人天 | 研发C |
| 4. 修复ESLint循环依赖 | 1人天 | 研发C |
| 5. 修复Enum类型shadow | 1人天 | 研发C |
| 6. 修复动态模块加载 | 2人天 | 研发C+D |

**总计**: 7人天

### 第二周 (P1高优先级)

**目标**: 完成10个P1高优先级TODO

| 任务 | 预计工作量 | 负责人 |
|------|-----------|--------|
| 7. 实现QuotaManagement保存 | 1人天 | 研发C |
| 8. 修复RoleList复制 | 0.5人天 | 研发C |
| 9. 实现TenantList API集成 | 1人天 | 研发C |
| 10. 修复TenantList删除 | 0.5人天 | 研发C |
| ... (其他P1任务) | 7人天 | 研发C |

**总计**: 10人天

---

## 📈 进度跟踪

### 当前状态

- **总TODO数**: 80
- **已完成**: 0 (0%)
- **进行中**: 0
- **待完成**: 80

### 目标

- **第一周末**: P0完成率100%
- **第二周末**: P0+P1完成率100%
- **月底**: P0+P1完成率100%, P2完成率≥80%

---

## 🔄 持续更新机制

### 每周更新

- 每周五更新TODO状态
- 每月发布TODO进展报告
- 每季度重新评估优先级

### 新增TODO流程

1. 在代码中添加TODO时,按格式注释
2. 每月梳理代码中的新TODO
3. 分类并分配优先级
4. 更新本跟踪文档

---

**文档维护**: 研发C (前端工程师)
**下次更新**: 每周五下午
**版本历史**:
- v1.0 (2025-12-30): 初始版本,80个TODO
