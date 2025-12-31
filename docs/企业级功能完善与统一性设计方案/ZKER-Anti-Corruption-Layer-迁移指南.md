# ZKER 防腐层（Anti-Corruption Layer）迁移指南

**版本**: v1.0
**日期**: 2025-01-01
**状态**: 正式发布
**适用范围**: 全体开发人员

---

## 📋 目录

- [1. 概述](#1-概述)
- [2. 架构原则](#2-架构原则)
- [3. 目录结构](#3-目录结构)
- [4. 迁移步骤](#4-迁移步骤)
- [5. 示例代码](#5-示例代码)
- [6. 常见问题](#6-常见问题)
- [7. 检查清单](#7-检查清单)

---

## 1. 概述

### 1.1 背景

当前项目中存在**33处跨层依赖**，主要原因是缺少防腐层（Anti-Corruption Layer），导致：

- `bizpkg` 层直接依赖 `api/model` 层（20+文件）
- `domain` 层直接依赖 `api/model` 层（13个文件）

这种跨层依赖违反了 **DDD 分层架构原则**，导致：

1. **依赖倒置原则（DIP）违反**：底层（api/model）影响上层（domain、bizpkg）
2. **耦合度过高**：模型变更需要修改多处代码
3. **测试困难**：无法独立测试domain层
4. **可维护性差**：架构边界不清晰

### 1.2 解决方案

建立 **防腐层（Anti-Corruption Layer）**，隔离外部模型（api/model）和内部领域模型（domain、bizpkg）。

**核心思想**：
```
┌─────────────────────────────────────────────────────────┐
│                     API Layer                            │
│                  (api/model/*)                           │
└─────────────────────────────────────────────────────────┘
                          ↑↓
                    【Adapter Layer】
                (application/adapter/*)
                    防腐层隔离点
                          ↑↓
┌─────────────────────────────────────────────────────────┐
│              Crossdomain Models                          │
│               (crossdomain/model/*)                      │
│           共享的领域模型定义                               │
└─────────────────────────────────────────────────────────┘
                          ↑↓
┌─────────────────────────────────────────────────────────┐
│     Application Layer / Domain Layer / Bizpkg           │
│           使用 crossdomain/model                         │
└─────────────────────────────────────────────────────────┘
```

---

## 2. 架构原则

### 2.1 依赖倒置原则（DIP）

**❌ 错误示例**：bizpkg依赖api/model
```go
// bizpkg/config/base/base.go
import "github.com/coze-dev/coze-studio/backend/api/model/admin/config"

func GetConfig() *config.BasicConfiguration {
    // 直接使用外部模型
}
```

**✅ 正确示例**：bizpkg依赖crossdomain/model
```go
// bizpkg/config/base/base.go
import "github.com/coze-dev/coze-studio/backend/crossdomain/model"

func GetConfig() *model.BasicConfiguration {
    // 使用领域模型
}
```

### 2.2 单一职责原则（SRP）

**Adapter层职责**：仅负责类型转换，不包含业务逻辑。

```go
// ✅ Good: 纯转换函数
func APIConfigToDomain(apiConfig *config.Model) *model.ModelConfig {
    // 仅做数据映射
}

// ❌ Bad: 包含业务逻辑
func APIConfigToDomain(apiConfig *config.Model) *model.ModelConfig {
    // 验证、计算等业务逻辑不应该在这里
    if apiConfig.ModelID <= 0 {
        return nil  // ❌ 验证逻辑应该在domain层
    }
    // ...
}
```

### 2.3 接口隔离原则（ISP）

最小化跨领域接口，只暴露必要的字段和方法。

```go
// ✅ Good: 最小化接口
type BotConfig struct {
    BotID   string
    BotName string
    // 仅包含必要字段
}

// ❌ Bad: 冗余接口
type BotConfig struct {
    BotID         string
    BotName       string
    // 100+ 字段...
}
```

---

## 3. 目录结构

### 3.1 crossdomain/model

通用跨领域模型，在domain、application、bizpkg层共享：

```
backend/crossdomain/model/
├── bot_config.go       # Bot配置模型
├── model_config.go     # LLM模型配置
├── workflow_config.go  # 工作流配置
└── base_config.go      # 基础系统配置
```

### 3.2 application/adapter

适配器层，负责api/model和crossdomain/model之间的转换：

```
backend/application/adapter/
├── config_adapter.go   # 配置相关适配器
├── bot_adapter.go      # Bot相关适配器
├── workflow_adapter.go # 工作流相关适配器
└── helper.go           # 批量转换辅助函数
```

---

## 4. 迁移步骤

### 4.1 渐进式迁移策略

**⚠️ 重要**：不要一次性重构所有代码！采用渐进式迁移：

1. **阶段1**：创建防腐层（已完成）
2. **阶段2**：新代码使用防腐层
3. **阶段3**：逐步重构旧代码
4. **阶段4**：验证和清理

### 4.2 具体步骤

#### Step 1: 识别依赖

查找直接依赖 `api/model` 的文件：

```bash
# Linux/macOS
grep -r "github.com/coze-dev/coze-studio/backend/api/model" backend/bizpkg

# Windows (PowerShell)
Select-String -Path "backend\bizpkg\*.go" -Pattern "github.com/coze-dev/coze-studio/backend/api/model"
```

#### Step 2: 添加适配器

如果适配器不存在，在 `application/adapter` 中添加：

```go
// application/adapter/new_feature_adapter.go
package adapter

import (
    "github.com/coze-dev/coze-studio/backend/api/model/admin/config"
    "github.com/coze-dev/coze-studio/backend/crossdomain/model"
)

// APINewFeatureToDomain 转换函数
func APINewFeatureToDomain(apiFeature *config.NewFeature) *model.NewFeature {
    // 实现转换逻辑
}
```

#### Step 3: 重构调用代码

```go
// ❌ Before: 直接使用API模型
import "github.com/coze-dev/coze-studio/backend/api/model/admin/config"

func ProcessModel(cfg *config.Model) error {
    // ...
}

// ✅ After: 使用领域模型 + 适配器
import (
    "github.com/coze-dev/coze-studio/backend/application/adapter"
    "github.com/coze-dev/coze-studio/backend/crossdomain/model"
)

func ProcessModel(cfg *model.ModelConfig) error {
    // ...
}

// 调用方使用适配器转换
func Handler(apiCfg *config.Model) error {
    domainCfg := adapter.APIConfigToDomainModel(apiCfg)
    return ProcessModel(domainCfg)
}
```

#### Step 4: 编译验证

```bash
cd backend
go build ./...
```

#### Step 5: 测试验证

```bash
go test ./bizpkg/... -v
go test ./domain/... -v
```

---

## 5. 示例代码

### 5.1 Bot配置迁移

**Before（依赖API层）**：
```go
// bizpkg/some_service.go
import "github.com/coze-dev/coze-studio/backend/api/model/app/bot_common"

func GetBotModel(info *bot_common.ModelInfo) (string, error) {
    if info.ModelId == nil {
        return "", errors.New("model_id is nil")
    }
    return fmt.Sprintf("model-%d", *info.ModelId), nil
}
```

**After（使用防腐层）**：
```go
// bizpkg/some_service.go
import "github.com/coze-dev/coze-studio/backend/crossdomain/model"

func GetBotModel(info *model.ModelInfo) (string, error) {
    if info.ModelID == nil {
        return "", errors.New("model_id is nil")
    }
    return fmt.Sprintf("model-%d", *info.ModelID), nil
}

// 调用方在application层使用适配器
// application/bot_service.go
import (
    "github.com/coze-dev/coze-studio/backend/application/adapter"
    apimodel "github.com/coze-dev/coze-studio/backend/api/model/app/bot_common"
)

func (s *BotService) ProcessModel(apiInfo *apimodel.ModelInfo) error {
    // 转换为领域模型
    domainInfo := adapter.APIModelInfoToDomain(apiInfo)

    // 调用bizpkg层
    _, err := bizpkg.GetBotModel(domainInfo)
    return err
}
```

### 5.2 LLM模型配置迁移

**Before**：
```go
// bizpkg/llm/modelbuilder/model_builder.go
import (
    "github.com/coze-dev/coze-studio/backend/api/model/admin/config"
    "github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
)

func NewModelBuilder(modelClass developer_api.ModelClass, cfg *config.Model) (Service, error) {
    // 直接使用API模型
}
```

**After**（使用Wrapper隔离）：
```go
// bizpkg/llm/modelbuilder/adapter.go
import (
    "github.com/coze-dev/coze-studio/backend/api/model/admin/config"
    "github.com/coze-dev/coze-studio/backend/crossdomain/model"
)

// Wrapper隔离外部依赖
type ModelConfigWrapper struct {
    cfg *config.Model
}

func BuildModelWithWrapper(
    ctx interface{},
    modelWrapper *ModelConfigWrapper,
    params *LLMParams,
) (ToolCallingChatModel, error) {
    // 通过wrapper访问，隔离直接依赖
}
```

### 5.3 批量转换示例

```go
// application/adapter/helper.go

func BatchAPIConfigToDomainModel(
    apiConfigs []*config.Model,
) []*model.ModelConfig {
    if apiConfigs == nil {
        return nil
    }

    domainConfigs := make([]*model.ModelConfig, 0, len(apiConfigs))
    for _, apiConfig := range apiConfigs {
        if domainConfig := APIConfigToDomainModel(apiConfig); domainConfig != nil {
            domainConfigs = append(domainConfigs, domainConfig)
        }
    }

    return domainConfigs
}
```

---

## 6. 常见问题

### 6.1 是否需要立即重构所有代码？

**答：不需要！** 采用渐进式迁移：

1. **新代码**：强制使用防腐层
2. **旧代码**：在修改时逐步重构
3. **紧急修复**：可以临时保留旧代码，添加TODO注释

### 6.2 bizpkg/config使用kvstore，如何处理？

**答：持久化层可以保留API模型类型**，因为：

1. kvstore需要序列化/反序列化，使用API模型是合理的
2. bizpkg/config作为配置访问层，不违反分层原则
3. **关键**：不要在bizpkg的业务逻辑中传播API模型

### 6.3 类型别名（Type Alias）是否可以？

**答：不推荐！** 类型别名无法真正隔离依赖：

```go
// ❌ Bad: 类型别名仍暴露外部类型
type BasicConfiguration = config.BasicConfiguration

// ✅ Good: 使用独立的领域模型
type BasicConfiguration struct {
    AdminEmails string
    // ...
}
```

### 6.4 如何处理循环依赖？

**答：通过依赖倒置解决**：

```go
// ❌ 循环依赖
// bizpkg → api/model → bizpkg

// ✅ 正确：依赖抽象
// bizpkg → crossdomain/model
// api/model → crossdomain/model (通过adapter转换)
```

### 6.5 性能开销如何？

**答：开销可忽略**：

1. 适配器只是简单的字段映射，无复杂计算
2. Go编译器会优化简单的转换函数
3. 架构清晰度提升带来的收益远大于微小性能开销

---

## 7. 检查清单

### 7.1 代码审查清单

在提交代码前，检查以下项：

#### ✅ 依赖检查

- [ ] bizpkg层不直接import `api/model`
- [ ] domain层不直接import `api/model`
- [ ] 仅application层可以import `api/model`（通过adapter）

#### ✅ 架构检查

- [ ] 使用crossdomain/model作为共享模型
- [ ] 转换逻辑在application/adapter中
- [ ] 适配器函数只做转换，不包含业务逻辑

#### ✅ 命名检查

- [ ] API模型转Domain：`API{Feature}ToDomain{Feature}`
- [ ] Domain模型转API：`Domain{Feature}ToAPI{Feature}`
- [ ] 批量转换：`BatchAPI{Feature}ToDomain{Feature}`

### 7.2 重构检查清单

修改旧代码时，检查以下项：

- [ ] 删除对 `api/model` 的import
- [ ] 添加对 `crossdomain/model` 的import
- [ ] 在调用链中添加adapter转换
- [ ] 编译通过
- [ ] 测试通过
- [ ] 更新相关文档

---

## 8. 附录

### 8.1 已完成的防腐层

#### crossdomain/model

| 文件 | 说明 | 状态 |
|------|------|------|
| `bot_config.go` | Bot配置模型 | ✅ 完成 |
| `model_config.go` | LLM模型配置 | ✅ 完成 |
| `workflow_config.go` | 工作流配置 | ✅ 完成 |
| `base_config.go` | 基础系统配置 | ✅ 完成 |

#### application/adapter

| 文件 | 说明 | 状态 |
|------|------|------|
| `config_adapter.go` | 配置适配器 | ✅ 完成 |
| `bot_adapter.go` | Bot适配器 | ✅ 完成 |
| `helper.go` | 批量转换辅助函数 | ✅ 完成 |

#### bizpkg Wrapper

| 文件 | 说明 | 状态 |
|------|------|------|
| `adapter.go` | LLM模型构建Wrapper | ✅ 完成 |

### 8.2 待重构的文件

#### bizpkg层（20个文件）

```
bizpkg/config/base/base.go              → 考虑保留（kvstore持久化）
bizpkg/config/config.go                 → 考虑保留（类型别名）
bizpkg/config/knowledge/knowledge.go    → 待重构
bizpkg/config/modelmgr/*.go             → 待重构
bizpkg/llm/modelbuilder/*.go            → 已提供Wrapper，待迁移
```

#### domain层（13个文件）

```
domain/agent/singleagent/entity/*.go    → 待重构
domain/search/entity/*.go               → 待重构
domain/workflow/entity/vo/*.go          → 待重构
domain/plugin/dto/*.go                  → 考虑保留（DTO层）
```

### 8.3 参考资料

- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)
- [实现差距分析与研发计划](./ZKER-实现差距分析与研发计划_v1.0.md)
- [DDD分层架构最佳实践](https://herbertograca.com/2017/09/14/ddd-layered-architecture/)

---

## 更新日志

| 版本 | 日期 | 变更说明 | 作者 |
|------|------|----------|------|
| v1.0 | 2025-01-01 | 初始版本，建立防腐层体系 | AI助手 |

---

**🎯 下一步行动**：

1. ✅ 阅读本迁移指南
2. ✅ 理解防腐层设计原则
3. ✅ 新代码强制使用防腐层
4. 🔄 逐步重构旧代码（优先级：P1）

**⚠️ 重要提醒**：

- 不要一次性大规模重构
- 优先处理新代码和频繁修改的代码
- 保留旧代码的兼容性，添加TODO标记
- 重构后必须运行测试验证

---

**问题反馈**：如有疑问，请联系架构组或创建Issue。
