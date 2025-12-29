# 13-BotBuilder_提示词版本管理补充完整版

**模块**: 13-BotBuilder扩展
**扩展内容**: 提示词版本管理、A/B测试、智能回滚
**版本**: v2.0
**日期**: 2025-01-03
**优先级**: P1
**工期**: 1.5周

---

## 目录

- [1. 功能概述](#1-功能概述)
- [2. 数据库设计](#2-数据库设计)
- [3. 版本控制服务](#3-版本控制服务)
- [4. 版本对比引擎](#4-版本对比引擎)
- [5. A/B测试引擎](#5-ab测试引擎)
- [6. 智能回滚机制](#6-智能回滚机制)
- [7. API设计](#7-api设计)
- [8. 前端组件](#8-前端组件)
- [9. 实施计划](#9-实施计划)

---

## 1. 功能概述

### 1.1 核心目标

**提示词版本管理系统** 提供Git-like的提示词版本管理能力：

- ✅ **版本控制** - 每次修改自动创建版本，支持分支管理
- ✅ **版本对比** - Diff视图高亮显示差异
- ✅ **A/B测试** - 多版本并行测试，数据驱动决策
- ✅ **智能回滚** - 一键回滚到任意历史版本
- ✅ **效果追踪** - 每个版本的质量指标追踪

### 1.2 应用场景

**场景1：提示词迭代**
```
开发者A修改Bot提示词：
1. 修改系统提示词
2. 点击"保存版本" → 自动创建 v2.0
3. 版本描述："优化了专业术语解释"
4. 系统自动记录：修改时间、修改人、变更内容
```

**场景2：版本对比**
```
开发者想要对比 v1.0 和 v2.0 的差异：
- 打开版本列表
- 选择两个版本进行对比
- 系统显示并排Diff视图（绿色=新增，红色=删除）
- 可逐行查看每个变更点
```

**场景3：A/B测试**
```
产品经理想要测试新提示词效果：
1. 创建 v3.0 版本
2. 启动A/B测试：v2.0 (50%流量) vs v3.0 (50%流量)
3. 运行24小时后查看数据：
   - v2.0: 平均评分4.2, 成功率85%
   - v3.0: 平均评分4.5, 成功率92%
4. 决定：v3.0 效果更好，全量发布
```

**场景4：智能回滚**
```
新版本 v4.0 上线后出现严重问题：
- 用户评分骤降到 2.5
- 大量负面反馈
- 开发者点击"一键回滚"到 v3.0
- 系统立即切换，2分钟内恢复服务
```

---

## 2. 数据库设计

### 2.1 提示词版本表

```sql
CREATE TABLE prompt_versions (
    id VARCHAR(64) PRIMARY KEY COMMENT '版本ID (UUID)',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    version_number VARCHAR(32) NOT NULL COMMENT '版本号 (v1.0, v2.0...)',
    version_name VARCHAR(100) COMMENT '版本名称',
    version_description TEXT COMMENT '版本描述',

    -- 提示词内容
    system_prompt TEXT NOT NULL COMMENT '系统提示词',
    user_prompt_template TEXT COMMENT '用户提示词模板',
    temperature DECIMAL(3,2) COMMENT '温度参数 (0.0-1.0)',
    max_tokens INT COMMENT '最大Token数',
    top_p DECIMAL(3,2) COMMENT 'Top-P采样参数',

    -- 版本关系
    parent_version_id VARCHAR(64) COMMENT '父版本ID',
    is_active BOOLEAN DEFAULT FALSE COMMENT '是否为当前激活版本',
    is_ab_test BOOLEAN DEFAULT FALSE COMMENT '是否参与A/B测试',

    -- 效果统计
    total_conversations INT DEFAULT 0 COMMENT '总对话数',
    avg_rating DECIMAL(3,2) COMMENT '平均评分 (1-5)',
    success_rate DECIMAL(5,2) COMMENT '成功率 (%)',
    avg_response_time_ms INT COMMENT '平均响应时间(毫秒)',
    total_evaluations INT DEFAULT 0 COMMENT '总评价数',

    -- A/B测试数据
    ab_test_traffic_percent INT COMMENT 'A/B测试流量分配 (%)',
    ab_test_start_at DATETIME COMMENT 'A/B测试开始时间',
    ab_test_end_at DATETIME COMMENT 'A/B测试结束时间',

    -- 元数据
    created_by BIGINT COMMENT '创建人ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    activated_at DATETIME COMMENT '激活时间',
    archived_at DATETIME COMMENT '归档时间',
    archived_reason VARCHAR(200) COMMENT '归档原因',

    INDEX idx_bot_id (bot_id),
    INDEX idx_version_number (bot_id, version_number),
    INDEX idx_is_active (bot_id, is_active),
    INDEX idx_ab_test (is_ab_test, ab_test_start_at),
    UNIQUE KEY uk_bot_version (bot_id, version_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='提示词版本表';
```

### 2.2 提示词变更历史表

```sql
CREATE TABLE prompt_change_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '变更ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    version_id_from VARCHAR(64) COMMENT '源版本ID',
    version_id_to VARCHAR(64) NOT NULL COMMENT '目标版本ID',

    -- 变更详情
    change_type ENUM('create', 'update', 'rollback', 'ab_test_winner') NOT NULL COMMENT '变更类型',
    change_summary VARCHAR(200) NOT NULL COMMENT '变更摘要',

    -- 变更内容 (JSON格式)
    field_changes JSON COMMENT '字段变更详情 {"system_prompt":{"from":"...","to":"..."}}',
    diff_summary TEXT COMMENT 'Diff摘要',

    -- 操作信息
    operated_by BIGINT NOT NULL COMMENT '操作人ID',
    operated_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    operation_source ENUM('manual', 'auto_rollback', 'ab_test') COMMENT '操作来源',

    INDEX idx_bot_id (bot_id),
    INDEX idx_version_to (version_id_to),
    INDEX idx_operated_at (operated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='提示词变更历史表';
```

### 2.3 A/B测试实验表

```sql
CREATE TABLE prompt_ab_tests (
    id VARCHAR(64) PRIMARY KEY COMMENT '实验ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    experiment_name VARCHAR(200) NOT NULL COMMENT '实验名称',
    experiment_description TEXT COMMENT '实验描述',

    -- 实验配置
    status ENUM('draft', 'running', 'paused', 'completed', 'cancelled') DEFAULT 'draft',
    start_time DATETIME COMMENT '开始时间',
    end_time DATETIME COMMENT '结束时间',
    min_sample_size INT DEFAULT 100 COMMENT '最小样本量',

    -- 流量分配
    traffic_allocation JSON COMMENT '流量分配 {"v2.0":50,"v3.0":50}',

    -- 实验结果
    winner_version_id VARCHAR(64) COMMENT '获胜版本ID',
    confidence_level DECIMAL(5,2) COMMENT '置信度 (%)',
    statistical_significance BOOLEAN COMMENT '统计显著性',

    -- 创建信息
    created_by BIGINT COMMENT '创建人ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_bot_id (bot_id),
    INDEX idx_status (status),
    INDEX idx_start_time (start_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='A/B测试实验表';
```

### 2.4 提示词评价记录表

```sql
CREATE TABLE prompt_evaluations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评价ID',
    version_id VARCHAR(64) NOT NULL COMMENT '版本ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    user_id BIGINT COMMENT '用户ID',

    -- 评分
    rating INT CHECK (rating >= 1 AND rating <= 5) COMMENT '用户评分 (1-5)',
    is_successful BOOLEAN COMMENT '是否成功 (用户判断)',
    feedback_text TEXT COMMENT '文字反馈',

    -- 性能指标
    response_time_ms INT COMMENT '响应时间(毫秒)',
    tokens_used INT COMMENT '使用的Token数',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '评价时间',

    INDEX idx_version_id (version_id),
    INDEX idx_bot_id (bot_id),
    INDEX idx_rating (rating),
    UNIQUE KEY uk_conversation (conversation_id) -- 每个会话只评价一次
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='提示词评价记录表';
```

---

## 3. 版本控制服务

### 3.1 版本管理核心服务

```go
package prompt

import (
    "context"
    "time"
    "github.com/cloudwego/hertz/pkg/common/hlog"
)

// PromptVersionService 提示词版本管理服务
type PromptVersionService struct {
    versionRepo    repository.PromptVersionRepository
    historyRepo    repository.ChangeHistoryRepository
    evaluationRepo repository.EvaluationRepository
    diffEngine     *DiffEngine
}

// CreateVersion 创建新版本
func (s *PromptVersionService) CreateVersion(
    ctx context.Context,
    req *CreateVersionRequest,
) (*CreateVersionResponse, error) {
    // 1. 获取当前激活版本（作为父版本）
    activeVersion, _ := s.versionRepo.GetActiveVersion(ctx, req.BotID)

    // 2. 生成新版本号
    nextVersionNumber := s.generateNextVersion(ctx, req.BotID)

    // 3. 构建版本对象
    version := &entity.PromptVersion{
        ID:                 uuid.New().String(),
        BotID:              req.BotID,
        VersionNumber:      nextVersionNumber,
        VersionName:        req.VersionName,
        VersionDescription: req.Description,
        SystemPrompt:       req.SystemPrompt,
        UserPromptTemplate: req.UserPromptTemplate,
        Temperature:        req.Temperature,
        MaxTokens:          req.MaxTokens,
        TopP:               req.TopP,
        ParentVersionID:    getPtrID(activeVersion),
        IsActive:           false, // 新版本默认不激活
        CreatedBy:          req.CreatedBy,
    }

    // 4. 保存版本
    if err := s.versionRepo.Create(ctx, version); err != nil {
        hlog.Errorf("Failed to create version: %v", err)
        return nil, err
    }

    // 5. 记录变更历史
    if activeVersion != nil {
        s.recordChangeHistory(ctx, activeVersion, version, req.CreatedBy)
    }

    return &CreateVersionResponse{
        VersionID:     version.ID,
        VersionNumber: nextVersionNumber,
        Message:       "版本创建成功",
    }, nil
}

// ActivateVersion 激活版本
func (s *PromptVersionService) ActivateVersion(
    ctx context.Context,
    versionID string,
    operatedBy int64,
) error {
    // 1. 获取目标版本
    version, err := s.versionRepo.GetByID(ctx, versionID)
    if err != nil {
        return err
    }

    // 2. 获取当前激活版本
    activeVersion, _ := s.versionRepo.GetActiveVersion(ctx, version.BotID)

    // 3. 停用当前版本
    if activeVersion != nil {
        activeVersion.IsActive = false
        s.versionRepo.Update(ctx, activeVersion)
    }

    // 4. 激活新版本
    version.IsActive = true
    version.ActivatedAt = time.Now()
    if err := s.versionRepo.Update(ctx, version); err != nil {
        return err
    }

    // 5. 记录变更历史
    if activeVersion != nil {
        s.historyRepo.Create(ctx, &entity.ChangeHistory{
            BotID:         version.BotID,
            VersionIDFrom: activeVersion.ID,
            VersionIDTo:   version.ID,
            ChangeType:    "update",
            ChangeSummary: fmt.Sprintf("从 %s 激活到 %s", activeVersion.VersionNumber, version.VersionNumber),
            OperatedBy:    operatedBy,
            OperationSource: "manual",
        })
    }

    return nil
}

// RollbackVersion 回滚到指定版本
func (s *PromptVersionService) RollbackVersion(
    ctx context.Context,
    targetVersionID string,
    operatedBy int64,
    reason string,
) error {
    // 1. 获取目标版本
    targetVersion, err := s.versionRepo.GetByID(ctx, targetVersionID)
    if err != nil {
        return err
    }

    // 2. 获取当前版本
    currentVersion, _ := s.versionRepo.GetActiveVersion(ctx, targetVersion.BotID)

    // 3. 创建回滚快照（基于目标版本创建新版本）
    rollbackVersion := &entity.PromptVersion{
        ID:                 uuid.New().String(),
        BotID:              targetVersion.BotID,
        VersionNumber:      s.generateNextVersion(ctx, targetVersion.BotID),
        VersionName:        fmt.Sprintf("回滚到 %s", targetVersion.VersionNumber),
        VersionDescription: reason,
        SystemPrompt:       targetVersion.SystemPrompt,
        UserPromptTemplate: targetVersion.UserPromptTemplate,
        Temperature:        targetVersion.Temperature,
        MaxTokens:          targetVersion.MaxTokens,
        TopP:               targetVersion.TopP,
        ParentVersionID:    getPtrID(currentVersion),
        IsActive:           true, // 回滚版本立即激活
        CreatedBy:          operatedBy,
    }

    // 4. 保存回滚版本
    if err := s.versionRepo.Create(ctx, rollbackVersion); err != nil {
        return err
    }

    // 5. 停用当前版本
    if currentVersion != nil {
        currentVersion.IsActive = false
        s.versionRepo.Update(ctx, currentVersion)
    }

    // 6. 记录变更历史
    s.historyRepo.Create(ctx, &entity.ChangeHistory{
        BotID:         targetVersion.BotID,
        VersionIDFrom: getPtrID(currentVersion),
        VersionIDTo:   rollbackVersion.ID,
        ChangeType:    "rollback",
        ChangeSummary: fmt.Sprintf("回滚到 %s: %s", targetVersion.VersionNumber, reason),
        OperatedBy:    operatedBy,
        OperationSource: "manual",
    })

    hlog.Infof("Rolled back bot %s from %s to %s",
        targetVersion.BotID,
        currentVersion.VersionNumber,
        targetVersion.VersionNumber)

    return nil
}

// generateNextVersion 生成下一个版本号
func (s *PromptVersionService) generateNextVersion(
    ctx context.Context,
    botID string,
) string {
    // 获取最新版本号
    latestVersion, _ := s.versionRepo.GetLatestVersion(ctx, botID)

    if latestVersion == nil {
        return "v1.0"
    }

    // 解析版本号 (v1.0 -> 1.0)
    versionNum := strings.TrimPrefix(latestVersion.VersionNumber, "v")
    parts := strings.Split(versionNum, ".")
    major, _ := strconv.Atoi(parts[0])
    minor, _ := strconv.Atoi(parts[1])

    // 递增次版本号
    minor++

    return fmt.Sprintf("v%d.%d", major, minor)
}

// recordChangeHistory 记录变更历史
func (s *PromptVersionService) recordChangeHistory(
    ctx context.Context,
    fromVersion, toVersion *entity.PromptVersion,
    operatedBy int64,
) {
    // 生成Diff
    diff := s.diffEngine.CompareVersions(fromVersion, toVersion)

    s.historyRepo.Create(ctx, &entity.ChangeHistory{
        BotID:         fromVersion.BotID,
        VersionIDFrom: fromVersion.ID,
        VersionIDTo:   toVersion.ID,
        ChangeType:    "update",
        ChangeSummary: fmt.Sprintf("更新到 %s", toVersion.VersionNumber),
        FieldChanges:  diff.FieldChanges,
        DiffSummary:   diff.Summary,
        OperatedBy:    operatedBy,
        OperationSource: "manual",
    })
}
```

---

## 4. 版本对比引擎

### 4.1 Diff生成服务

```go
package prompt

import (
    "context"
    "encoding/json"
    "strings"
    "github.com/pmezard/go-diff/diffmatchpatch"
)

// DiffEngine 版本对比引擎
type DiffEngine struct {
    dmp *diffmatchpatch.DiffMatchPatch
}

// VersionDiff 版本对比结果
type VersionDiff struct {
    FieldChanges map[string]FieldChange `json:"field_changes"`
    Summary      string                 `json:"summary"`
    HTMLDiff     string                 `json:"html_diff"`
}

type FieldChange struct {
    FieldName string      `json:"field_name"`
    From      interface{} `json:"from"`
    To        interface{} `json:"to"`
    DiffType  string      `json:"diff_type"` // added, removed, modified
}

// NewDiffEngine 创建Diff引擎
func NewDiffEngine() *DiffEngine {
    return &DiffEngine{
        dmp: diffmatchpatch.New(),
    }
}

// CompareVersions 对比两个版本
func (e *DiffEngine) CompareVersions(
    from, to *entity.PromptVersion,
) *VersionDiff {
    diff := &VersionDiff{
        FieldChanges: make(map[string]FieldChange),
    }

    // 1. 对比系统提示词
    if from.SystemPrompt != to.SystemPrompt {
        diff.FieldChanges["system_prompt"] = FieldChange{
            FieldName: "系统提示词",
            From:      e.truncateText(from.SystemPrompt, 100),
            To:        e.truncateText(to.SystemPrompt, 100),
            DiffType:  "modified",
        }
    }

    // 2. 对比用户提示词模板
    if from.UserPromptTemplate != to.UserPromptTemplate {
        diff.FieldChanges["user_prompt_template"] = FieldChange{
            FieldName: "用户提示词模板",
            From:      from.UserPromptTemplate,
            To:        to.UserPromptTemplate,
            DiffType:  "modified",
        }
    }

    // 3. 对比温度参数
    if from.Temperature != to.Temperature {
        diff.FieldChanges["temperature"] = FieldChange{
            FieldName: "温度",
            From:      from.Temperature,
            To:        to.Temperature,
            DiffType:  "modified",
        }
    }

    // 4. 对比最大Token数
    if from.MaxTokens != to.MaxTokens {
        diff.FieldChanges["max_tokens"] = FieldChange{
            FieldName: "最大Token数",
            From:      from.MaxTokens,
            To:        to.MaxTokens,
            DiffType:  "modified",
        }
    }

    // 5. 生成摘要
    changeCount := len(diff.FieldChanges)
    diff.Summary = fmt.Sprintf("本次更新包含 %d 处变更", changeCount)

    return diff
}

// GenerateTextualDiff 生成文本Diff（用于展示）
func (e *DiffEngine) GenerateTextualDiff(
    fromText, toText string,
) string {
    // 执行Diff
    diffs := e.dmp.DiffMain(fromText, toText, false)
    e.dmp.DiffCleanupSemantic(diffs)

    // 生成HTML视图
    var htmlBuilder strings.Builder

    for _, diff := range diffs {
        text := strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)

        switch diff.Type {
        case diffmatchpatch.DiffInsert:
            htmlBuilder.WriteString(fmt.Sprintf(`<span class="diff-added">%s</span>`, text))
        case diffmatchpatch.DiffDelete:
            htmlBuilder.WriteString(fmt.Sprintf(`<span class="diff-removed">%s</span>`, text))
        case diffmatchpatch.DiffEqual:
            htmlBuilder.WriteString(fmt.Sprintf(`<span class="diff-equal">%s</span>`, text))
        }
    }

    return htmlBuilder.String()
}

// GenerateSideBySideDiff 生成并排Diff视图
func (e *DiffEngine) GenerateSideBySideDiff(
    from, to *entity.PromptVersion,
) *SideBySideDiff {
    // 系统提示词并排对比
    systemPromptDiffs := e.generateSideBySide(
        from.SystemPrompt,
        to.SystemPrompt,
    )

    // 用户提示词并排对比
    userPromptDiffs := e.generateSideBySide(
        from.UserPromptTemplate,
        to.UserPromptTemplate,
    )

    return &SideBySideDiff{
        SystemPromptDiff: systemPromptDiffs,
        UserPromptDiff:   userPromptDiffs,
        ParameterChanges: e.extractParameterChanges(from, to),
    }
}

type SideBySideDiff struct {
    SystemPromptDiff *SideBySideTextDiff `json:"system_prompt_diff"`
    UserPromptDiff   *SideBySideTextDiff `json:"user_prompt_diff"`
    ParameterChanges map[string]string   `json:"parameter_changes"`
}

type SideBySideTextDiff struct {
    FromLines []DiffLine `json:"from_lines"`
    ToLines   []DiffLine `json:"to_lines"`
}

type DiffLine struct {
    LineNumber int    `json:"line_number"`
    Content    string `json:"content"`
    ChangeType string `json:"change_type"` // added, removed, unchanged
}

// generateSideBySide 生成并排对比
func (e *DiffEngine) generateSideBySide(
    fromText, toText string,
) *SideBySideTextDiff {
    fromLines := strings.Split(fromText, "\n")
    toLines := strings.Split(toText, "\n")

    diff := e.dmp.DiffMain(fromText, toText, false)
    e.dmp.DiffCleanupSemantic(diffs)

    // 构建并排视图
    result := &SideBySideTextDiff{
        FromLines: make([]DiffLine, 0),
        ToLines:   make([]DiffLine, 0),
    }

    // 这里需要更复杂的算法来映射行号
    // 简化实现：直接显示全文
    for i, line := range fromLines {
        result.FromLines = append(result.FromLines, DiffLine{
            LineNumber: i + 1,
            Content:    line,
            ChangeType: "unchanged",
        })
    }

    for i, line := range toLines {
        result.ToLines = append(result.ToLines, DiffLine{
            LineNumber: i + 1,
            Content:    line,
            ChangeType: "unchanged",
        })
    }

    return result
}

// extractParameterChanges 提取参数变更
func (e *DiffEngine) extractParameterChanges(
    from, to *entity.PromptVersion,
) map[string]string {
    changes := make(map[string]string)

    if from.Temperature != to.Temperature {
        changes["temperature"] = fmt.Sprintf("%.2f → %.2f", from.Temperature, to.Temperature)
    }

    if from.MaxTokens != to.MaxTokens {
        changes["max_tokens"] = fmt.Sprintf("%d → %d", from.MaxTokens, to.MaxTokens)
    }

    if from.TopP != to.TopP {
        changes["top_p"] = fmt.Sprintf("%.2f → %.2f", from.TopP, to.TopP)
    }

    return changes
}

func (e *DiffEngine) truncateText(text string, maxLen int) string {
    runes := []rune(text)
    if len(runes) <= maxLen {
        return text
    }
    return string(runes[:maxLen]) + "..."
}
```

---

## 5. A/B测试引擎

### 5.1 A/B测试服务

```go
package prompt

import (
    "context"
    "hash/fnv"
    "time"
)

// ABTestService A/B测试服务
type ABTestService struct {
    versionRepo  repository.PromptVersionRepository
    testRepo     repository.ABTestRepository
    evalRepo     repository.EvaluationRepository
    statsEngine  *StatisticsEngine
}

// CreateABTest 创建A/B测试实验
func (s *ABTestService) CreateABTest(
    ctx context.Context,
    req *CreateABTestRequest,
) (*CreateABTestResponse, error) {
    // 1. 验证版本
    versionA, _ := s.versionRepo.GetByID(ctx, req.VersionIDA)
    versionB, _ := s.versionRepo.GetByID(ctx, req.VersionIDB)

    if versionA == nil || versionB == nil {
        return nil, errors.New("版本不存在")
    }

    // 2. 创建实验
    experiment := &entity.ABTest{
        ID:                  uuid.New().String(),
        BotID:               versionA.BotID,
        ExperimentName:      req.ExperimentName,
        ExperimentDescription: req.Description,
        Status:              "draft",
        MinSampleSize:       req.MinSampleSize,
        TrafficAllocation: map[string]int{
            versionA.VersionNumber: 50, // 默认50/50
            versionB.VersionNumber: 50,
        },
        CreatedBy: req.CreatedBy,
    }

    if err := s.testRepo.Create(ctx, experiment); err != nil {
        return nil, err
    }

    // 3. 标记版本参与A/B测试
    versionA.IsABTest = true
    versionB.IsABTest = true
    versionA.ABTestTrafficPercent = 50
    versionB.ABTestTrafficPercent = 50

    s.versionRepo.Update(ctx, versionA)
    s.versionRepo.Update(ctx, versionB)

    return &CreateABTestResponse{
        ExperimentID: experiment.ID,
        Message:      "A/B测试实验创建成功",
    }, nil
}

// StartABTest 启动A/B测试
func (s *ABTestService) StartABTest(
    ctx context.Context,
    experimentID string,
) error {
    // 1. 获取实验
    experiment, err := s.testRepo.GetByID(ctx, experimentID)
    if err != nil {
        return err
    }

    // 2. 更新状态
    experiment.Status = "running"
    experiment.StartTime = time.Now()

    if err := s.testRepo.Update(ctx, experiment); err != nil {
        return err
    }

    // 3. 获取参与测试的版本
    versions, _ := s.versionRepo.GetByABTestID(ctx, experimentID)

    for _, version := range versions {
        version.ABTestStartAt = time.Now()
        s.versionRepo.Update(ctx, version)
    }

    hlog.Infof("Started A/B test %s for bot %s", experimentID, experiment.BotID)

    return nil
}

// RouteRequest 路由请求到对应版本（用于A/B测试）
func (s *ABTestService) RouteRequest(
    ctx context.Context,
    botID string,
    userID int64,
    conversationID string,
) (*entity.PromptVersion, error) {
    // 1. 检查是否有运行中的A/B测试
    experiments, _ := s.testRepo.GetRunningTestsByBot(ctx, botID)

    if len(experiments) == 0 {
        // 没有A/B测试，返回当前激活版本
        return s.versionRepo.GetActiveVersion(ctx, botID)
    }

    // 2. 获取最新实验
    experiment := experiments[0]

    // 3. 根据用户ID哈希分配版本（确保同一用户始终访问同一版本）
    hash := fnv.New32()
    hash.Write([]byte(fmt.Sprintf("%d-%s", userID, experiment.ID)))
    hashValue := hash.Sum32()

    // 4. 获取流量分配
    allocation := experiment.TrafficAllocation
    threshold := hashValue % 100

    var selectedVersionNumber string
    cumulative := 0

    for versionNum, percent := range allocation {
        cumulative += percent
        if int(threshold) < cumulative {
            selectedVersionNumber = versionNum
            break
        }
    }

    // 5. 返回对应版本
    version, _ := s.versionRepo.GetByNumber(ctx, botID, selectedVersionNumber)

    hlog.Infof("Routed user %d to version %s (A/B test: %s)",
        userID, selectedVersionNumber, experiment.ID)

    return version, nil
}

// AnalyzeABTest 分析A/B测试结果
func (s *ABTestService) AnalyzeABTest(
    ctx context.Context,
    experimentID string,
) (*ABTestAnalysis, error) {
    // 1. 获取实验
    experiment, err := s.testRepo.GetByID(ctx, experimentID)
    if err != nil {
        return nil, err
    }

    // 2. 获取参与测试的版本
    versions, _ := s.versionRepo.GetByABTestID(ctx, experimentID)

    if len(versions) < 2 {
        return nil, errors.New("参与测试的版本不足2个")
    }

    // 3. 收集各版本的数据
    var versionStats []*VersionStats
    for _, version := range versions {
        evaluations, _ := s.evalRepo.GetByVersionID(ctx, version.ID)

        totalEvals := len(evaluations)
        avgRating := 0.0
        successfulCount := 0

        for _, eval := range evaluations {
            avgRating += float64(eval.Rating)
            if eval.IsSuccessful {
                successfulCount++
            }
        }

        if totalEvals > 0 {
            avgRating = avgRating / float64(totalEvals)
        }

        successRate := 0.0
        if totalEvals > 0 {
            successRate = float64(successfulCount) / float64(totalEvals) * 100
        }

        versionStats = append(versionStats, &VersionStats{
            VersionID:     version.ID,
            VersionNumber: version.VersionNumber,
            TotalSamples:  totalEvals,
            AvgRating:     avgRating,
            SuccessRate:   successRate,
            AvgResponseTime: version.AvgResponseTimeMs,
        })
    }

    // 4. 统计显著性检验
    significance, confidence, winner := s.statsEngine.PerformTTest(
        versionStats[0],
        versionStats[1],
    )

    // 5. 构建分析报告
    return &ABTestAnalysis{
        ExperimentID:         experimentID,
        TotalSamples:         versionStats[0].TotalSamples + versionStats[1].TotalSamples,
        Duration:             time.Since(experiment.StartTime),
        VersionStats:         versionStats,
        Winner:               winner,
        Confidence:           confidence,
        StatisticalSignificance: significance,
        Recommendation:       s.generateRecommendation(versionStats, winner),
    }, nil
}

// CompleteABTest 完成A/B测试（选择获胜版本）
func (s *ABTestService) CompleteABTest(
    ctx context.Context,
    experimentID string,
    winnerVersionID string,
    operatedBy int64,
) error {
    // 1. 获取实验
    experiment, _ := s.testRepo.GetByID(ctx, experimentID)

    // 2. 更新实验状态
    experiment.Status = "completed"
    experiment.EndTime = time.Now()
    experiment.WinnerVersionID = winnerVersionID

    s.testRepo.Update(ctx, experiment)

    // 3. 激活获胜版本
    s.ActivateVersion(ctx, winnerVersionID, operatedBy)

    // 4. 标记其他版本为非A/B测试
    versions, _ := s.versionRepo.GetByABTestID(ctx, experimentID)

    for _, version := range versions {
        if version.ID != winnerVersionID {
            version.IsABTest = false
            version.ABTestEndAt = time.Now()
            s.versionRepo.Update(ctx, version)
        }
    }

    // 5. 记录变更历史
    winnerVersion, _ := s.versionRepo.GetByID(ctx, winnerVersionID)

    s.historyRepo.Create(ctx, &entity.ChangeHistory{
        BotID:         experiment.BotID,
        VersionIDTo:   winnerVersionID,
        ChangeType:    "ab_test_winner",
        ChangeSummary: fmt.Sprintf("A/B测试获胜: %s", winnerVersion.VersionNumber),
        OperatedBy:    operatedBy,
        OperationSource: "ab_test",
    })

    hlog.Infof("Completed A/B test %s, winner: %s", experimentID, winnerVersion.VersionNumber)

    return nil
}
```

### 5.2 统计分析引擎

```go
package prompt

import (
    "math"
)

// StatisticsEngine 统计分析引擎
type StatisticsEngine struct{}

// PerformTTest 执行T检验（统计显著性检验）
func (s *StatisticsEngine) PerformTTest(
    statsA, statsB *VersionStats,
) (significant bool, confidence float64, winner string) {
    // 使用Z检验（大样本情况）

    // 1. 计算样本比例（成功率）
    p1 := statsA.SuccessRate / 100.0
    p2 := statsB.SuccessRate / 100.0
    n1 := statsA.TotalSamples
    n2 := statsB.TotalSamples

    // 2. 计算合并比例
    pPool := (float64(n1)*p1 + float64(n2)*p2) / float64(n1+n2)

    // 3. 计算标准误差
    se := math.Sqrt(pPool * (1 - pPool) * (1.0/float64(n1) + 1.0/float64(n2)))

    // 4. 计算Z分数
    zScore := (p2 - p1) / se

    // 5. 计算P值（双尾检验）
    pValue := 2 * (1 - normalCDF(math.Abs(zScore)))

    // 6. 判断显著性（P < 0.05）
    significant = pValue < 0.05

    // 7. 计算置信度
    confidence = (1 - pValue) * 100

    // 8. 确定获胜者
    if p2 > p1 {
        winner = statsB.VersionNumber
    } else if p1 > p2 {
        winner = statsA.VersionNumber
    }

    return
}

// normalCDF 标准正态分布累积分布函数
func normalCDF(x float64) float64 {
    // 使用近似公式计算
    a1 :=  0.254829592
    a2 := -0.284496736
    a3 :=  1.421413741
    a4 := -1.453152027
    a5 :=  1.061405429
    p  :=  0.3275911

    sign := 1.0
    if x < 0 {
        sign = -1.0
    }
    x = math.Abs(x) / math.Sqrt(2.0)

    t := 1.0 / (1.0 + p*x)
    y := 1.0 - (((((a5*t + a4)*t) + a3)*t + a2)*t + a1)*t*math.Exp(-x*x)

    return 0.5 * (1.0 + sign*y)
}

// GenerateRecommendation 生成建议
func (s *ABTestService) generateRecommendation(
    stats []*VersionStats,
    winner string,
) string {
    if winner == "" {
        return "样本量不足或差异不显著，建议延长测试时间或调整提示词策略"
    }

    var winnerStats *VersionStats
    for _, stat := range stats {
        if stat.VersionNumber == winner {
            winnerStats = stat
            break
        }
    }

    return fmt.Sprintf(
        "版本 %s 表现最佳，评分 %.2f，成功率 %.1f%%，建议全量发布。",
        winner,
        winnerStats.AvgRating,
        winnerStats.SuccessRate,
    )
}
```

---

## 6. 智能回滚机制

### 6.1 自动回滚监控

```go
package prompt

import (
    "context"
    "time"
)

// AutoRollbackMonitor 自动回滚监控器
type AutoRollbackMonitor struct {
    versionRepo    repository.PromptVersionRepository
    evaluationRepo repository.EvaluationRepository
    versionSvc     *PromptVersionService
    alertSvc       *AlertService
}

// MonitorVersionPerformance 监控版本性能
func (m *AutoRollbackMonitor) MonitorVersionPerformance(
    ctx context.Context,
    botID string,
) error {
    // 1. 获取当前激活版本
    activeVersion, err := m.versionRepo.GetActiveVersion(ctx, botID)
    if err != nil {
        return err
    }

    // 2. 获取父版本（作为基准）
    var baselineVersion *entity.PromptVersion
    if activeVersion.ParentVersionID != nil {
        baselineVersion, _ = m.versionRepo.GetByID(ctx, *activeVersion.ParentVersionID)
    }

    if baselineVersion == nil {
        return nil // 没有基准版本，无法比较
    }

    // 3. 收集最近1小时的数据
    oneHourAgo := time.Now().Add(-1 * time.Hour)
    recentEvals, _ := m.evaluationRepo.GetByVersionIDAfter(
        ctx,
        activeVersion.ID,
        oneHourAgo,
    )

    if len(recentEvals) < 10 {
        return nil // 样本量不足
    }

    // 4. 计算当前版本指标
    currentAvgRating := 0.0
    currentSuccessRate := 0.0

    for _, eval := range recentEvals {
        currentAvgRating += float64(eval.Rating)
        if eval.IsSuccessful {
            currentSuccessRate += 1.0
        }
    }

    currentAvgRating /= float64(len(recentEvals))
    currentSuccessRate = currentSuccessRate / float64(len(recentEvals)) * 100

    // 5. 比较基准版本
    baselineRating, _ := m.versionRepo.GetAvgRating(ctx, baselineVersion.ID)
    baselineSuccessRate, _ := m.versionRepo.GetSuccessRate(ctx, baselineVersion.ID)

    // 6. 判断是否需要回滚
    rollbackThreshold := 0.8 // 阈值：评分下降超过20%

    if baselineRating > 0 && currentAvgRating < baselineRating*rollbackThreshold {
        // 评分严重下降，触发自动回滚
        m.triggerAutoRollback(
            ctx,
            activeVersion,
            baselineVersion,
            fmt.Sprintf(
                "评分从 %.2f 下降到 %.2f，触发自动回滚",
                baselineRating,
                currentAvgRating,
            ),
        )

        return nil
    }

    if baselineSuccessRate > 0 && currentSuccessRate < baselineSuccessRate*rollbackThreshold {
        // 成功率严重下降，触发自动回滚
        m.triggerAutoRollback(
            ctx,
            activeVersion,
            baselineVersion,
            fmt.Sprintf(
                "成功率从 %.1f%% 下降到 %.1f%%，触发自动回滚",
                baselineSuccessRate,
                currentSuccessRate,
            ),
        )

        return nil
    }

    return nil
}

// triggerAutoRollback 触发自动回滚
func (m *AutoRollbackMonitor) triggerAutoRollback(
    ctx context.Context,
    currentVersion, rollbackTarget *entity.PromptVersion,
    reason string,
) {
    hlog.Warnf("Auto rollback triggered for bot %s: %s", currentVersion.BotID, reason)

    // 1. 执行回滚
    err := m.versionSvc.RollbackVersion(
        ctx,
        rollbackTarget.ID,
        0, // 系统操作
        reason,
    )

    if err != nil {
        hlog.Errorf("Auto rollback failed: %v", err)
        return
    }

    // 2. 发送告警
    m.alertSvc.SendAlert(ctx, &Alert{
        Type:    "auto_rollback",
        BotID:   currentVersion.BotID,
        Title:   "提示词版本自动回滚",
        Message: reason,
        Severity: "high",
    })
}

// StartMonitoring 启动监控任务
func (m *AutoRollbackMonitor) StartMonitoring(
    ctx context.Context,
    botIDs []string,
) {
    // 每5分钟检查一次
    ticker := time.NewTicker(5 * time.Minute)

    go func() {
        for {
            select {
            case <-ticker.C:
                for _, botID := range botIDs {
                    m.MonitorVersionPerformance(ctx, botID)
                }
            case <-ctx.Done():
                ticker.Stop()
                return
            }
        }
    }()
}
```

---

## 7. API设计

### 7.1 版本管理API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| GET | /api/v1/prompt-versions | 列出版本 | `?bot_id=bot123` | `[{version_number:"v1.0",...}]` |
| POST | /api/v1/prompt-versions | 创建版本 | `{"bot_id":"bot123","system_prompt":"..."}` | `{"version_id":"v456","version_number":"v2.0"}` |
| GET | /api/v1/prompt-versions/:id | 查看版本详情 | - | `{id:"v456",system_prompt:"...",...}` |
| PUT | /api/v1/prompt-versions/:id/activate | 激活版本 | - | `{"success":true}` |
| POST | /api/v1/prompt-versions/:id/rollback | 回滚版本 | `{"reason":"性能下降"}` | `{"success":true,"new_version_id":"v457"}` |
| GET | /api/v1/prompt-versions/:id/stats | 查看版本统计 | - | `{avg_rating:4.5,success_rate:92,...}` |
| GET | /api/v1/prompt-versions/diff | 版本对比 | `?from=v1.0&to=v2.0` | 见下方完整响应 |

**版本对比响应**:
```json
{
  "from_version": {
    "id": "v123",
    "version_number": "v1.0",
    "system_prompt": "你是一个有用的AI助手...",
    "temperature": 0.7
  },
  "to_version": {
    "id": "v456",
    "version_number": "v2.0",
    "system_prompt": "你是一个专业的企业AI助手...",
    "temperature": 0.5
  },
  "field_changes": {
    "system_prompt": {
      "field_name": "系统提示词",
      "from": "你是一个有用的AI助手...",
      "to": "你是一个专业的企业AI助手...",
      "diff_type": "modified"
    },
    "temperature": {
      "field_name": "温度",
      "from": 0.7,
      "to": 0.5,
      "diff_type": "modified"
    }
  },
  "summary": "本次更新包含 2 处变更",
  "html_diff": "<span class='diff-removed'>你是一个有用的</span><span class='diff-added'>你是一个专业的企业</span>AI助手..."
}
```

### 7.2 A/B测试API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| POST | /api/v1/prompt-ab-tests | 创建A/B测试 | `{"experiment_name":"测试v2.0","version_id_a":"v123","version_id_b":"v456"}` | `{"experiment_id":"ab789"}` |
| POST | /api/v1/prompt-ab-tests/:id/start | 启动测试 | - | `{"success":true}` |
| POST | /api/v1/prompt-ab-tests/:id/stop | 停止测试 | - | `{"success":true}` |
| GET | /api/v1/prompt-ab-tests/:id/analysis | 分析结果 | - | 见下方完整响应 |
| POST | /api/v1/prompt-ab-tests/:id/complete | 完成测试 | `{"winner_version_id":"v456"}` | `{"success":true}` |

**A/B测试分析响应**:
```json
{
  "experiment_id": "ab789",
  "total_samples": 500,
  "duration": "24h",
  "version_stats": [
    {
      "version_number": "v1.0",
      "total_samples": 250,
      "avg_rating": 4.2,
      "success_rate": 85.0
    },
    {
      "version_number": "v2.0",
      "total_samples": 250,
      "avg_rating": 4.5,
      "success_rate": 92.0
    }
  ],
  "winner": "v2.0",
  "confidence": 95.0,
  "statistical_significance": true,
  "recommendation": "版本 v2.0 表现最佳，评分 4.5，成功率 92.0%，建议全量发布。"
}
```

---

## 8. 前端组件

### 8.1 版本列表组件

```typescript
// components/prompt/VersionList.tsx
import React, { useEffect, useState } from 'react';
import { Table, Tag, Button, Space, Tooltip } from '@douyinfe/semi-ui';

interface PromptVersion {
  id: string;
  version_number: string;
  version_name: string;
  version_description: string;
  is_active: boolean;
  is_ab_test: boolean;
  avg_rating: number;
  success_rate: number;
  total_evaluations: number;
  created_at: string;
  created_by_name: string;
}

export const VersionList: React.FC<Props> = ({ botId, onActivate, onRollback }) => {
  const [versions, setVersions] = useState<PromptVersion[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchVersions();
  }, [botId]);

  const fetchVersions = async () => {
    const response = await fetch(`/api/v1/prompt-versions?bot_id=${botId}`);
    const data = await response.json();
    setVersions(data);
    setLoading(false);
  };

  const columns = [
    {
      title: '版本号',
      dataKey: 'version_number',
      width: 100,
      render: (version_number: string, record: PromptVersion) => (
        <Space>
          {version_number}
          {record.is_active && <Tag color="green">激活</Tag>}
          {record.is_ab_test && <Tag color="blue">A/B测试中</Tag>}
        </Space>
      ),
    },
    {
      title: '版本名称',
      dataKey: 'version_name',
    },
    {
      title: '描述',
      dataKey: 'version_description',
      width: 300,
      render: (text: string) => (
        <Tooltip content={text}>
          <div className="truncate">{text}</div>
        </Tooltip>
      ),
    },
    {
      title: '平均评分',
      dataKey: 'avg_rating',
      width: 120,
      render: (rating: number) => rating ? rating.toFixed(1) : '-',
    },
    {
      title: '成功率',
      dataKey: 'success_rate',
      width: 100,
      render: (rate: number) => rate ? `${rate.toFixed(1)}%` : '-',
    },
    {
      title: '评价数',
      dataKey: 'total_evaluations',
      width: 100,
    },
    {
      title: '创建时间',
      dataKey: 'created_at',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '创建人',
      dataKey: 'created_by_name',
      width: 120,
    },
    {
      title: '操作',
      width: 200,
      render: (_: any, record: PromptVersion) => (
        <Space>
          {!record.is_active && (
            <Button
              size="small"
              onClick={() => onActivate(record.id)}
            >
              激活
            </Button>
          )}
          {!record.is_active && (
            <Button
              size="small"
              onClick={() => onRollback(record.id)}
            >
              回滚到此版本
            </Button>
          )}
          <Button
            size="small"
            onClick={() => window.open(`/prompt-versions/${record.id}`)}
          >
            查看详情
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <Table
      columns={columns}
      dataSource={versions}
      loading={loading}
      pagination={false}
      rowKey="id"
    />
  );
};
```

### 8.2 版本对比组件

```typescript
// components/prompt/VersionDiffViewer.tsx
import React, { useEffect, useState } from 'react';
import { Select, Button, Card, Descriptions } from '@douyinfe/semi-ui';
import ReactDiffViewer from 'react-diff-viewer';

export const VersionDiffViewer: ReactFC = ({ botId }) => {
  const [versions, setVersions] = useState([]);
  const [fromVersion, setFromVersion] = useState('');
  const [toVersion, setToVersion] = useState('');
  const [diff, setDiff] = useState(null);

  useEffect(() => {
    fetchVersions();
  }, [botId]);

  const fetchVersions = async () => {
    const response = await fetch(`/api/v1/prompt-versions?bot_id=${botId}`);
    const data = await response.json();
    setVersions(data);
    if (data.length >= 2) {
      setFromVersion(data[1].id);
      setToVersion(data[0].id);
    }
  };

  const handleCompare = async () => {
    const response = await fetch(
      `/api/v1/prompt-versions/diff?from=${fromVersion}&to=${toVersion}`
    );
    const data = await response.json();
    setDiff(data);
  };

  return (
    <div className="version-diff-viewer">
      <Card title="版本对比">
        {/* 版本选择 */}
        <div className="version-selectors">
          <Select
            value={fromVersion}
            onChange={setFromVersion}
            placeholder="选择源版本"
            style={{ width: 200, marginRight: 16 }}
          >
            {versions.map((v) => (
              <Select.Option key={v.id} value={v.id}>
                {v.version_number} - {v.version_name}
              </Select.Option>
            ))}
          </Select>

          <span>→</span>

          <Select
            value={toVersion}
            onChange={setToVersion}
            placeholder="选择目标版本"
            style={{ width: 200, marginLeft: 16, marginRight: 16 }}
          >
            {versions.map((v) => (
              <Select.Option key={v.id} value={v.id}>
                {v.version_number} - {v.version_name}
              </Select.Option>
            ))}
          </Select>

          <Button type="primary" onClick={handleCompare}>
            对比
          </Button>
        </div>

        {/* Diff结果 */}
        {diff && (
          <div className="diff-result">
            <Descriptions
              data={[
                { key: '源版本', value: diff.from_version.version_number },
                { key: '目标版本', value: diff.to_version.version_number },
                { key: '变更摘要', value: diff.summary },
              ]}
              row
            />

            {/* 系统提示词Diff */}
            <div className="diff-section">
              <h4>系统提示词</h4>
              <ReactDiffViewer
                oldValue={diff.from_version.system_prompt}
                newValue={diff.to_version.system_prompt}
                splitView={true}
                useDarkTheme={false}
              />
            </div>

            {/* 参数变更 */}
            {Object.keys(diff.field_changes).length > 0 && (
              <div className="parameter-changes">
                <h4>参数变更</h4>
                <Table
                  columns={[
                    { title: '参数', dataKey: 'field_name' },
                    { title: '旧值', dataKey: 'from' },
                    { title: '新值', dataKey: 'to' },
                  ]}
                  dataSource={Object.values(diff.field_changes)}
                  pagination={false}
                />
              </div>
            )}
          </div>
        )}
      </Card>
    </div>
  );
};
```

### 8.3 A/B测试管理组件

```typescript
// components/prompt/ABTestManager.tsx
import React, { useEffect, useState } from 'react';
import { Modal, Form, Button, Card, Progress, Tag } from '@douyinfe/semi-ui';

interface ABTest {
  id: string;
  experiment_name: string;
  status: 'draft' | 'running' | 'completed';
  start_time: string;
  version_stats: Array<{
    version_number: string;
    traffic_percent: number;
    avg_rating: number;
    success_rate: number;
  }>;
}

export const ABTestManager: React.FC<Props> = ({ botId }) => {
  const [tests, setTests] = useState<ABTest[]>([]);
  const [createModalVisible, setCreateModalVisible] = useState(false);

  useEffect(() => {
    fetchABTests();
  }, [botId]);

  const fetchABTests = async () => {
    const response = await fetch(`/api/v1/prompt-ab-tests?bot_id=${botId}`);
    const data = await response.json();
    setTests(data);
  };

  return (
    <div className="ab-test-manager">
      <div className="header">
        <h3>A/B测试实验</h3>
        <Button
          type="primary"
          onClick={() => setCreateModalVisible(true)}
        >
          创建新实验
        </Button>
      </div>

      {/* 实验列表 */}
      <div className="test-list">
        {tests.map((test) => (
          <Card key={test.id} title={test.experiment_name}>
            <div className="test-status">
              <Tag color={
                test.status === 'running' ? 'blue' :
                test.status === 'completed' ? 'green' : 'default'
              }>
                {test.status}
              </Tag>
              <span>开始时间: {new Date(test.start_time).toLocaleString()}</span>
            </div>

            {/* 版本对比 */}
            <div className="version-comparison">
              {test.version_stats.map((stat) => (
                <div key={stat.version_number} className="version-stat">
                  <h4>{stat.version_number} ({stat.traffic_percent}%)</h4>
                  <div>评分: {stat.avg_rating.toFixed(1)}</div>
                  <div>成功率: {stat.success_rate.toFixed(1)}%</div>
                  {test.status === 'running' && (
                    <Progress
                      percent={stat.traffic_percent}
                      showInfo={false}
                      stroke={stat.version_number === 'v2.0' ? 'var(--primary-color)' : '#999'}
                    />
                  )}
                </div>
              ))}
            </div>

            {/* 操作按钮 */}
            {test.status === 'draft' && (
              <Button onClick={() => startTest(test.id)}>启动实验</Button>
            )}
            {test.status === 'running' && (
              <>
                <Button onClick={() => viewAnalysis(test.id)}>查看分析</Button>
                <Button onClick={() => stopTest(test.id)}>停止实验</Button>
              </>
            )}
            {test.status === 'completed' && (
              <Button onClick={() => viewReport(test.id)}>查看报告</Button>
            )}
          </Card>
        ))}
      </div>

      {/* 创建实验Modal */}
      <CreateABTestModal
        visible={createModalVisible}
        botId={botId}
        onClose={() => setCreateModalVisible(false)}
        onSuccess={() => {
          setCreateModalVisible(false);
          fetchABTests();
        }}
      />
    </div>
  );
};
```

---

## 9. 实施计划

### 9.1 开发阶段划分

| 阶段 | 任务 | 工期 | 交付物 |
|------|------|------|--------|
| **第1周** | 数据库设计与创建 | 2天 | 4张表 |
| | 版本控制服务开发 | 3天 | 核心服务 |
| **第2周** | Diff对比引擎 | 2天 | Diff服务 |
| | A/B测试引擎 | 2天 | A/B测试服务 |
| | 智能回滚机制 | 1天 | 监控服务 |
| **第3周** | 前端组件开发 | 3天 | 3个组件 |
| | API集成测试 | 2天 | 测试报告 |

### 9.2 技术风险与缓解

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|----------|
| Diff算法性能 | 中 | 低 | 使用增量Diff + 缓存 |
| A/B测试流量分配不准确 | 高 | 中 | 基于用户ID哈希确保一致性 |
| 自动回滚误判 | 中 | 低 | 增加样本量阈值 + 人工确认 |

### 9.3 成功指标

- ✅ 版本切换时间: **< 2秒**
- ✅ Diff生成速度: **< 500ms** (1000行文本)
- ✅ A/B测试流量分配准确性: **≥ 99%**
- ✅ 自动回滚响应时间: **< 5分钟**

---

**文档结束**
