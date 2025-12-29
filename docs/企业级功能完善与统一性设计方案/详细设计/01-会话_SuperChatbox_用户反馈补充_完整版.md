# 01-会话_SuperChatbox_用户反馈补充完整版

**模块**: 01-会话_SuperChatbox扩展
**扩展内容**: 用户反馈收集、分析、持续改进
**版本**: v2.0
**日期**: 2025-01-03
**优先级**: P1
**工期**: 1.5周

---

## 目录

- [1. 功能概述](#1-功能概述)
- [2. 数据库设计](#2-数据库设计)
- [3. 反馈收集系统](#3-反馈收集系统)
- [4. 反馈分析引擎](#4-反馈分析引擎)
- [5. 持续改进机制](#5-持续改进机制)
- [6. API设计](#6-api设计)
- [7. 前端组件](#7-前端组件)
- [8. 实施计划](#9-实施计划)

---

## 1. 功能概述

### 1.1 核心目标

**用户反馈系统** 提供完整的用户反馈闭环：

- ✅ **多渠道收集** - 内嵌反馈、邮件调研、定期问卷
- ✅ **智能分类** - AI自动分类反馈主题和情感
- ✅ **数据驱动分析** - 评分趋势、问题Top 10、优化建议
- ✅ **持续改进** - 反馈驱动的产品迭代流程

### 1.2 应用场景

**场景1：单次对话反馈**
```
用户与AI对话结束后，弹出反馈面板：
- 👍 有帮助 (1秒快速反馈)
- 👎 没帮助 (需选择原因：回答不准确/不完整/不相关/其他)
- ⭐⭐⭐⭐⭐ 星级评分
- 💬 文字反馈（可选）
```

**场景2：周期性满意度调研**
```
每月自动发送邮件调研：
- 总体满意度评分
- 各维度评分（准确度/响应速度/易用性/功能完整性）
- 开放性问题反馈
- NPS (净推荐值) 调查
```

**场景3：反馈分析仪表盘**
```
管理员查看反馈分析：
- 平均评分趋势图
- 有助率趋势
- 常见问题Top 10（按频次排序）
- 情感分析分布
- AI生成的优化建议
```

---

## 2. 数据库设计

### 2.1 对话反馈表

```sql
CREATE TABLE conversation_feedback (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '反馈ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT COMMENT '用户ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    message_id VARCHAR(64) COMMENT '消息ID (如果是针对单条消息的反馈)',

    -- 评分
    is_helpful BOOLEAN COMMENT '是否有帮助 (快速反馈)',
    rating TINYINT COMMENT '星级评分 (1-5)',
    nps_score TINYINT COMMENT 'NPS评分 (0-10)',

    -- 分类
    feedback_category ENUM('accurate', 'incomplete', 'irrelevant', 'slow', 'other') COMMENT '反馈分类',
    feedback_reason VARCHAR(100) COMMENT '反馈原因',

    -- 文字反馈
    feedback_text TEXT COMMENT '文字反馈内容',

    -- 自动分析
    sentiment ENUM('positive', 'neutral', 'negative') COMMENT '情感分析 (AI自动)',
    keywords JSON COMMENT '关键词提取 (AI自动)',
    ai_summary TEXT COMMENT 'AI摘要',

    -- 元数据
    bot_id VARCHAR(64) COMMENT 'Bot ID',
    bot_name VARCHAR(100) COMMENT 'Bot名称',
    model_used VARCHAR(100) COMMENT '使用的模型',
    conversation_turns INT COMMENT '对话轮数',
    response_time_ms INT COMMENT '响应时间(毫秒)',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_created (tenant_id, created_at),
    INDEX idx_conversation (conversation_id),
    INDEX idx_bot (bot_id),
    INDEX idx_rating (rating),
    INDEX idx_helpful (is_helpful)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='对话反馈表';
```

### 2.2 反馈统计汇总表

```sql
CREATE TABLE feedback_summary (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    bot_id VARCHAR(64) COMMENT 'Bot ID (NULL表示整体汇总)',
    summary_date DATE NOT NULL COMMENT '汇总日期',
    summary_type ENUM('daily', 'weekly', 'monthly') NOT NULL DEFAULT 'daily',

    -- 统计指标
    total_feedback INT DEFAULT 0 COMMENT '总反馈数',
    helpful_count INT DEFAULT 0 COMMENT '有帮助数量',
    helpful_rate DECIMAL(5,2) COMMENT '有帮助率(%)',
    avg_rating DECIMAL(3,2) COMMENT '平均评分 (1-5)',

    -- 情感分布
    positive_count INT DEFAULT 0,
    neutral_count INT DEFAULT 0,
    negative_count INT DEFAULT 0,

    -- 分类统计
    category_accuracy INT DEFAULT 0 COMMENT '准确性问题',
    category_incomplete INT DEFAULT 0 COMMENT '完整性问题',
    category_irrelevant INT DEFAULT 0 COMMENT '相关性问题',
    category_slow INT DEFAULT 0 COMMENT '速度问题',
    category_other INT DEFAULT 0 COMMENT '其他问题',

    -- NPS统计
    nps_promoters INT DEFAULT 0 COMMENT '推荐者 (9-10分)',
    nps_passives INT DEFAULT 0 COMMENT '中立者 (7-8分)',
    nps_detractors INT DEFAULT 0 COMMENT '贬损者 (0-6分)',
    nps_score INT COMMENT 'NPS分数 (-100~100)',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_tenant_bot_date_type (tenant_id, bot_id, summary_date, summary_type),
    INDEX idx_tenant_date (tenant_id, summary_date),
    INDEX idx_bot_date (bot_id, summary_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='反馈统计汇总表';
```

### 2.3 反馈问题Top榜表

```sql
CREATE TABLE feedback_top_issues (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    bot_id VARCHAR(64) COMMENT 'Bot ID (NULL表示全局)',
    issue_title VARCHAR(200) NOT NULL COMMENT '问题描述',
    issue_category VARCHAR(50) COMMENT '问题分类',

    -- 统计
    mention_count INT DEFAULT 0 COMMENT '提及次数',
    affected_users INT DEFAULT 0 COMMENT '影响用户数',
    first_mentioned_at DATETIME COMMENT '首次提及时间',
    last_mentioned_at DATETIME COMMENT '最后提及时间',

    -- 趋势
    trend ENUM('increasing', 'stable', 'decreasing') COMMENT '趋势',

    -- 优先级
    priority ENUM('low', 'medium', 'high', 'critical') COMMENT '优先级',
    status ENUM('open', 'investigating', 'fixing', 'resolved') DEFAULT 'open',

    -- 关联反馈ID
    related_feedback_ids JSON COMMENT '相关反馈ID列表',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_tenant_bot (tenant_id, bot_id),
    INDEX idx_priority (priority),
    INDEX idx_status (status),
    INDEX idx_mentions (mention_count)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='反馈问题Top榜表';
```

### 2.4 反馈改进任务表

```sql
CREATE TABLE feedback_improvement_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    task_type ENUM('bug_fix', 'feature_request', 'optimization', 'documentation') NOT NULL,

    -- 问题
    issue_summary VARCHAR(200) NOT NULL COMMENT '问题摘要',
    issue_description TEXT COMMENT '问题描述',
    related_feedback_id BIGINT COMMENT '关联反馈ID',

    -- 分析
    root_cause TEXT COMMENT '根因分析',
    proposed_solution TEXT COMMENT '建议方案',

    -- 优先级
    priority ENUM('p0', 'p1', 'p2', 'p3') DEFAULT 'p2' COMMENT '优先级',
    estimated_effort_hours DECIMAL(5,1) COMMENT '预估工时(小时)',

    -- 状态
    status ENUM('backlog', 'planned', 'in_progress', 'testing', 'done', 'cancelled') DEFAULT 'backlog',
    assigned_to BIGINT COMMENT '负责人ID',
    sprint_id VARCHAR(50) COMMENT '迭代ID',

    -- 结果
    completed_at DATETIME COMMENT '完成时间',
    impact_summary TEXT COMMENT '影响总结',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_priority (priority),
    INDEX idx_type (task_type),
    INDEX idx_assigned (assigned_to)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='反馈改进任务表';
```

---

## 3. 反馈收集系统

### 3.1 反馈收集服务

```go
package feedback

import (
    "context"
    "time"
    "github.com/cloudwego/hertz/pkg/common/hlog"
)

// FeedbackCollectionService 反馈收集服务
type FeedbackCollectionService struct {
    feedbackRepo repository.FeedbackRepository
    summaryRepo  repository.FeedbackSummaryRepository
    analyzer     *FeedbackAnalyzer
}

// RecordFeedback 记录用户反馈
func (s *FeedbackCollectionService) RecordFeedback(
    ctx context.Context,
    req *RecordFeedbackRequest,
) (*RecordFeedbackResponse, error) {
    // 1. 构建反馈记录
    feedback := &entity.ConversationFeedback{
        TenantID:        req.TenantID,
        UserID:          req.UserID,
        ConversationID:  req.ConversationID,
        MessageID:       req.MessageID,
        IsHelpful:       req.IsHelpful,
        Rating:          req.Rating,
        NPSScore:        req.NPSScore,
        FeedbackCategory: req.FeedbackCategory,
        FeedbackReason:  req.FeedbackReason,
        FeedbackText:    req.FeedbackText,
        BotID:           req.BotID,
        BotName:         req.BotName,
    }

    // 2. 持久化反馈
    if err := s.feedbackRepo.Create(ctx, feedback); err != nil {
        hlog.Errorf("Failed to record feedback: %v", err)
        return nil, err
    }

    // 3. 异步AI分析
    go s.analyzer.AnalyzeFeedback(context.Background(), feedback)

    // 4. 异步更新汇总统计
    go s.updateSummary(context.Background(), req.TenantID, req.BotID)

    return &RecordFeedbackResponse{
        FeedbackID: feedback.ID,
        Message:    "感谢您的反馈！",
    }, nil
}

// GetFeedbackStats 获取反馈统计
func (s *FeedbackCollectionService) GetFeedbackStats(
    ctx context.Context,
    req *GetFeedbackStatsRequest,
) (*FeedbackStatsResponse, error) {
    // 1. 获取汇总数据
    summary, err := s.summaryRepo.GetByDateRange(
        ctx,
        req.TenantID,
        req.BotID,
        req.StartDate,
        req.EndDate,
    )
    if err != nil {
        return nil, err
    }

    // 2. 计算总体指标
    var totalFeedback, helpfulCount int
    var totalRating float64
    var positive, neutral, negative int

    for _, s := range summary {
        totalFeedback += int(s.TotalFeedback)
        helpfulCount += int(s.HelpfulCount)
        totalRating += float64(s.AvgRating) * float64(s.TotalFeedback)
        positive += int(s.PositiveCount)
        neutral += int(s.NeutralCount)
        negative += int(s.NegativeCount)
    }

    helpfulRate := 0.0
    if totalFeedback > 0 {
        helpfulRate = float64(helpfulCount) / float64(totalFeedback) * 100
        totalRating = totalRating / float64(totalFeedback)
    }

    // 3. 构建响应
    return &FeedbackStatsResponse{
        Period: Period{
            StartDate: req.StartDate,
            EndDate:   req.EndDate,
        },
        TotalFeedback: totalFeedback,
        HelpfulRate:   helpfulRate,
        AvgRating:     totalRating,
        SentimentDistribution: SentimentDistribution{
            Positive: positive,
            Neutral:  neutral,
            Negative: negative,
        },
        DailyBreakdown: s.buildDailyBreakdown(summary),
    }, nil
}
```

### 3.2 邮件调研服务

```go
package feedback

import (
    "context"
    "time"
)

// EmailSurveyService 邮件调研服务
type EmailSurveyService struct {
    feedbackRepo repository.FeedbackRepository
    emailSvc     *EmailService
    scheduler    *CronScheduler
}

// SendMonthlySurvey 发送月度满意度调研
func (s *EmailSurveyService) SendMonthlySurvey(
    ctx context.Context,
    tenantID string,
) error {
    // 1. 获取活跃用户列表（最近一个月有对话的用户）
    activeUsers, err := s.getActiveUsers(ctx, tenantID, 30)
    if err != nil {
        return err
    }

    // 2. 生成调研链接（带唯一token）
    for _, user := range activeUsers {
        surveyToken := generateSurveyToken(user.ID, tenantID)

        surveyLink := fmt.Sprintf(
            "https://coze.zker.com/survey?token=%s",
            surveyToken,
        )

        // 3. 发送调研邮件
        email := &Email{
            To:      user.Email,
            Subject: "📋 zker - 月度满意度调研",
            Template: "monthly_survey",
            Data: map[string]interface{}{
                "UserName":   user.Name,
                "SurveyLink": surveyLink,
                "DueDate":    time.Now().Add(7 * 24 * time.Hour), // 7天后截止
            },
        }

        if err := s.emailSvc.Send(ctx, email); err != nil {
            hlog.Errorf("Failed to send survey email to user %d: %v", user.ID, err)
            continue
        }
    }

    return nil
}

// ScheduleMonthlySurveys 定时调度月度调研
func (s *EmailSurveyService) ScheduleMonthlySurveys() {
    // 每月1号上午10点发送
    s.scheduler.AddFunc("0 10 1 * *", func() {
        // 获取所有启用调研的租户
        tenants := s.getSurveyEnabledTenants()

        for _, tenant := range tenants {
            if err := s.SendMonthlySurvey(context.Background(), tenant.ID); err != nil {
                hlog.Errorf("Failed to send monthly survey for tenant %s: %v", tenant.ID, err)
            }
        }
    })
}
```

---

## 4. 反馈分析引擎

### 4.1 AI分析服务

```go
package feedback

import (
    "context"
    "encoding/json"
    "strings"
)

// FeedbackAnalyzer 反馈分析器
type FeedbackAnalyzer struct {
    llmClient *llm.Client
    repo      repository.FeedbackRepository
}

// AnalyzeFeedback 分析单条反馈
func (a *FeedbackAnalyzer) AnalyzeFeedback(
    ctx context.Context,
    feedback *entity.ConversationFeedback,
) error {
    // 1. 情感分析
    sentiment := a.analyzeSentiment(ctx, feedback.FeedbackText)
    feedback.Sentiment = sentiment

    // 2. 关键词提取
    keywords := a.extractKeywords(ctx, feedback.FeedbackText)
    feedback.Keywords = keywords

    // 3. AI摘要
    summary := a.generateSummary(ctx, feedback.FeedbackText)
    feedback.AISummary = summary

    // 4. 更新数据库
    return a.repo.Update(ctx, feedback)
}

// analyzeSentiment 情感分析
func (a *FeedbackAnalyzer) analyzeSentiment(
    ctx context.Context,
    text string,
) string {
    if text == "" {
        return "neutral"
    }

    prompt := fmt.Sprintf(
        "分析以下用户反馈的情感倾向，只返回 positive/neutral/negative 其中一个词：\n\n%s",
        text,
    )

    response, err := a.llmClient.Chat(ctx, &llm.ChatRequest{
        Model: "gpt-3.5-turbo",
        Messages: []llm.Message{
            {Role: "user", Content: prompt},
        },
        Temperature: 0.1,
    })

    if err != nil {
        hlog.Errorf("Sentiment analysis failed: %v", err)
        return "neutral"
    }

    sentiment := strings.ToLower(strings.TrimSpace(response.Content))
    if sentiment != "positive" && sentiment != "negative" {
        return "neutral"
    }

    return sentiment
}

// extractKeywords 关键词提取
func (a *FeedbackAnalyzer) extractKeywords(
    ctx context.Context,
    text string,
) json.RawMessage {
    if text == "" {
        return json.RawMessage("[]")
    }

    prompt := fmt.Sprintf(
        "从以下用户反馈中提取5个最重要的关键词，以JSON数组格式返回：\n\n%s\n\n示例格式：[\"关键词1\", \"关键词2\", ...]",
        text,
    )

    response, err := a.llmClient.Chat(ctx, &llm.ChatRequest{
        Model: "gpt-3.5-turbo",
        Messages: []llm.Message{
            {Role: "user", Content: prompt},
        },
        Temperature: 0.3,
    })

    if err != nil {
        return json.RawMessage("[]")
    }

    // 解析并验证JSON
    var keywords []string
    if err := json.Unmarshal([]byte(response.Content), &keywords); err != nil {
        hlog.Errorf("Failed to parse keywords: %v", err)
        return json.RawMessage("[]")
    }

    result, _ := json.Marshal(keywords)
    return json.RawMessage(result)
}

// generateSummary 生成AI摘要
func (a *FeedbackAnalyzer) generateSummary(
    ctx context.Context,
    text string,
) string {
    if text == "" {
        return ""
    }

    prompt := fmt.Sprintf(
        "用一句话总结以下用户反馈的核心问题（不超过50字）：\n\n%s",
        text,
    )

    response, err := a.llmClient.Chat(ctx, &llm.ChatRequest{
        Model: "gpt-3.5-turbo",
        Messages: []lll.Message{
            {Role: "user", Content: prompt},
        },
        Temperature: 0.3,
        MaxTokens:   100,
    })

    if err != nil {
        return ""
    }

    return response.Content
}
```

### 4.2 问题Top榜分析

```go
package feedback

import (
    "context"
    "strings"
)

// TopIssuesAnalyzer Top问题分析器
type TopIssuesAnalyzer struct {
    feedbackRepo repository.FeedbackRepository
    issueRepo    repository.TopIssueRepository
    llmClient    *llm.Client
}

// UpdateTopIssues 更新Top问题榜
func (a *TopIssuesAnalyzer) UpdateTopIssues(
    ctx context.Context,
    tenantID string,
    botID string,
) error {
    // 1. 获取最近7天的负面反馈
    negativeFeedbacks, err := a.feedbackRepo.GetNegativeFeedbacks(
        ctx,
        tenantID,
        botID,
        7, // 最近7天
    )
    if err != nil {
        return err
    }

    // 2. 使用LLM提取共性问题
    issues := a.extractCommonIssues(ctx, negativeFeedbacks)

    // 3. 更新或创建问题记录
    for _, issue := range issues {
        existingIssue, _ := a.issueRepo.GetByTitle(ctx, tenantID, botID, issue.Title)

        if existingIssue != nil {
            // 更新现有问题
            existingIssue.MentionCount += issue.MentionCount
            existingIssue.LastMentionedAt = time.Now()
            existingIssue.Trend = a.analyzeTrend(ctx, existingIssue)

            a.issueRepo.Update(ctx, existingIssue)
        } else {
            // 创建新问题
            issue.Priority = a.calculatePriority(issue)
            a.issueRepo.Create(ctx, issue)
        }
    }

    return nil
}

// extractCommonIssues 提取共性问题
func (a *TopIssuesAnalyzer) extractCommonIssues(
    ctx context.Context,
    feedbacks []*entity.ConversationFeedback,
) []*entity.TopIssue {
    // 1. 汇总所有反馈文本
    var feedbackTexts []string
    for _, f := range feedbacks {
        feedbackTexts = append(feedbackTexts, f.FeedbackText)
        if f.AISummary != "" {
            feedbackTexts = append(feedbackTexts, f.AISummary)
        }
    }

    // 2. 使用LLM提取Top 10共性问题
    prompt := fmt.Sprintf(
        "分析以下用户反馈，提取Top 10共性问题，返回JSON数组：\n\n%s\n\n"+
        "格式：[{\"title\":\"问题描述\",\"category\":\"分类\",\"affected_count\":影响用户数}, ...]",
        strings.Join(feedbackTexts, "\n\n"),
    )

    response, err := a.llmClient.Chat(ctx, &llm.ChatRequest{
        Model: "gpt-4",
        Messages: []llm.Message{
            {Role: "user", Content: prompt},
        },
        Temperature: 0.3,
    })

    if err != nil {
        return nil
    }

    // 3. 解析响应
    var issues []struct {
        Title        string `json:"title"`
        Category     string `json:"category"`
        AffectedCount int    `json:"affected_count"`
    }

    if err := json.Unmarshal([]byte(response.Content), &issues); err != nil {
        hlog.Errorf("Failed to parse issues: %v", err)
        return nil
    }

    // 4. 转换为实体
    var result []*entity.TopIssue
    for _, issue := range issues {
        result = append(result, &entity.TopIssue{
            TenantID:       tenantID,
            BotID:          botID,
            IssueTitle:     issue.Title,
            IssueCategory:  issue.Category,
            MentionCount:   issue.AffectedCount,
            AffectedUsers:  issue.AffectedCount,
            FirstMentionedAt: time.Now(),
            LastMentionedAt: time.Now(),
            Status:         "open",
        })
    }

    return result
}

// calculatePriority 计算优先级
func (a *TopIssuesAnalyzer) calculatePriority(
    issue *entity.TopIssue,
) string {
    // 基于影响用户数和趋势计算优先级
    if issue.MentionCount >= 50 {
        return "critical"
    } else if issue.MentionCount >= 20 {
        return "high"
    } else if issue.MentionCount >= 10 {
        return "medium"
    }
    return "low"
}
```

### 4.3 NPS分析

```go
package feedback

import (
    "context"
)

// NPSAnalyzer NPS分析器
type NPSAnalyzer struct {
    feedbackRepo repository.FeedbackRepository
}

// CalculateNPSScore 计算NPS分数
func (a *NPSAnalyzer) CalculateNPSScore(
    ctx context.Context,
    tenantID string,
    startDate, endDate time.Time,
) (*NPSReport, error) {
    // 1. 获取所有NPS评分
    feedbacks, err := a.feedbackRepo.GetNPSScores(
        ctx,
        tenantID,
        startDate,
        endDate,
    )
    if err != nil {
        return nil, err
    }

    // 2. 分类统计
    var promoters, passives, detractors int

    for _, f := range feedbacks {
        if f.NPSScore >= 9 {
            promoters++
        } else if f.NPSScore >= 7 {
            passives++
        } else {
            detractors++
        }
    }

    total := len(feedbacks)
    if total == 0 {
        return &NPSReport{
            TotalResponses: 0,
            NPSScore:       0,
            Promoters:      0,
            Passives:       0,
            Detractors:     0,
        }, nil
    }

    // 3. 计算NPS分数 (-100~100)
    promotersPercent := float64(promoters) / float64(total) * 100
    detractorsPercent := float64(detractors) / float64(total) * 100
    npsScore := int(promotersPercent - detractorsPercent)

    // 4. 评级
    grade := a.getNPSGrade(npsScore)

    return &NPSReport{
        TotalResponses:  total,
        NPSScore:        npsScore,
        Promoters:       promoters,
        Passives:        passives,
        Detractors:      detractors,
        PromotersPercent: promotersPercent,
        DetractorsPercent: detractorsPercent,
        Grade:           grade,
    }, nil
}

// getNPSGrade 获取NPS评级
func (a *NPSAnalyzer) getNPSGrade(score int) string {
    if score >= 70 {
        return "优秀 (Excellent)"
    } else if score >= 50 {
        return "良好 (Good)"
    } else if score >= 20 {
        return "一般 (Fair)"
    } else if score >= 0 {
        return "及格 (Poor)"
    }
    return "不及格 (Terrible)"
}
```

---

## 5. 持续改进机制

### 5.1 改进任务生成服务

```go
package feedback

import (
    "context"
)

// ImprovementTaskGenerator 改进任务生成器
type ImprovementTaskGenerator struct {
    issueRepo  repository.TopIssueRepository
    taskRepo   repository.ImprovementTaskRepository
    llmClient  *llm.Client
}

// GenerateTasksFromIssues 从Top问题生成改进任务
func (g *ImprovementTaskGenerator) GenerateTasksFromIssues(
    ctx context.Context,
    tenantID string,
) error {
    // 1. 获取未解决的高优问题
    issues, err := g.issueRepo.GetUnresolvedIssuesByPriority(
        ctx,
        tenantID,
        "critical", // 只处理critical级别
    )
    if err != nil {
        return err
    }

    // 2. 为每个问题生成改进任务
    for _, issue := range issues {
        // 检查是否已存在关联任务
        existingTask, _ := g.taskRepo.GetByIssueID(ctx, issue.ID)
        if existingTask != nil {
            continue
        }

        // 使用LLM生成任务建议
        taskSuggestion := g.generateTaskSuggestion(ctx, issue)

        // 创建任务
        task := &entity.ImprovementTask{
            TenantID:           tenantID,
            TaskType:          taskSuggestion.TaskType,
            IssueSummary:      issue.IssueTitle,
            IssueDescription:  taskSuggestion.Description,
            RelatedFeedbackID: issue.RelatedFeedbackIDs[0],
            RootCause:         taskSuggestion.RootCause,
            ProposedSolution:  taskSuggestion.Solution,
            Priority:          "p0", // Top问题默认p0
            EstimatedEffortHours: taskSuggestion.EstimatedHours,
            Status:            "backlog",
        }

        if err := g.taskRepo.Create(ctx, task); err != nil {
            hlog.Errorf("Failed to create improvement task: %v", err)
        }
    }

    return nil
}

// generateTaskSuggestion 生成任务建议
func (g *ImprovementTaskGenerator) generateTaskSuggestion(
    ctx context.Context,
    issue *entity.TopIssue,
) *TaskSuggestion {
    prompt := fmt.Sprintf(
        "基于以下用户反馈问题，生成改进任务建议：\n\n"+
        "问题描述：%s\n"+
        "提及次数：%d\n\n"+
        "请以JSON格式返回：\n"+
        "{\n"+
        "  \"task_type\": \"bug_fix|feature_request|optimization\",\n"+
        "  \"description\": \"详细描述\",\n"+
        "  \"root_cause\": \"根因分析\",\n"+
        "  \"solution\": \"建议解决方案\",\n"+
        "  \"estimated_hours\": 8\n"+
        "}",
        issue.IssueTitle,
        issue.MentionCount,
    )

    response, err := g.llmClient.Chat(ctx, &llm.ChatRequest{
        Model: "gpt-4",
        Messages: []llm.Message{
            {Role: "user", Content: prompt},
        },
        Temperature: 0.3,
    })

    if err != nil {
        return &TaskSuggestion{
            TaskType:     "bug_fix",
            Description:  issue.IssueTitle,
            RootCause:    "待分析",
            Solution:     "待确定",
            EstimatedHours: 8.0,
        }
    }

    var suggestion TaskSuggestion
    if err := json.Unmarshal([]byte(response.Content), &suggestion); err != nil {
        hlog.Errorf("Failed to parse task suggestion: %v", err)
        return &TaskSuggestion{
            TaskType:     "bug_fix",
            Description:  issue.IssueTitle,
            EstimatedHours: 8.0,
        }
    }

    return &suggestion
}
```

---

## 6. API设计

### 6.1 反馈记录API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| POST | /api/v1/feedback | 记录反馈 | `{"conversation_id":"c123","is_helpful":true,"rating":5}` | `{"feedback_id":456,"message":"感谢反馈"}` |
| GET | /api/v1/feedback/stats | 查询统计 | `?start_date=2025-01-01&end_date=2025-01-31` | 见下方完整响应 |
| GET | /api/v1/feedback/issues | 获取Top问题 | `?limit=10` | `[{issue_title:"回答不准确",mention_count:50,...}]` |
| GET | /api/v1/feedback/nps | 查询NPS | `?period=monthly` | `{nps_score:65,grade:"优秀",...}` |

**反馈统计响应**:
```json
{
  "period": {
    "start_date": "2025-01-01",
    "end_date": "2025-01-31"
  },
  "total_feedback": 1500,
  "helpful_rate": 85.5,
  "avg_rating": 4.3,
  "sentiment_distribution": {
    "positive": 1100,
    "neutral": 300,
    "negative": 100
  },
  "daily_breakdown": [
    {
      "date": "2025-01-01",
      "total_feedback": 50,
      "helpful_rate": 88.0,
      "avg_rating": 4.5
    }
  ],
  "category_breakdown": {
    "accurate": 80,
    "incomplete": 45,
    "irrelevant": 30,
    "slow": 20,
    "other": 25
  }
}
```

### 6.2 改进任务API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| GET | /api/v1/feedback/tasks | 获取改进任务列表 | `?status=backlog` | `[{id:1,issue_summary:"回答不准确",...}]` |
| POST | /api/v1/feedback/tasks/:id/start | 开始任务 | - | `{"success":true}` |
| POST | /api/v1/feedback/tasks/:id/complete | 完成任务 | `{"impact_summary":"优化后准确率提升20%"}` | `{"success":true}` |

---

## 7. 前端组件

### 7.1 反馈收集组件

```typescript
// components/feedback/FeedbackPanel.tsx
import React, { useState } from 'react';
import { Modal, Radio, Rate, Select, TextArea, Button, Toast } from '@douyinfe/semi-ui';

interface FeedbackPanelProps {
  conversationId: string;
  botId?: string;
  visible: boolean;
  onClose: () => void;
  onSubmit: (feedback: FeedbackData) => Promise<void>;
}

interface FeedbackData {
  isHelpful: boolean | null;
  rating?: number;
  category?: string;
  reason?: string;
  text?: string;
}

export const FeedbackPanel: React.FC<FeedbackPanelProps> = ({
  conversationId,
  botId,
  visible,
  onClose,
  onSubmit,
}) => {
  const [feedback, setFeedback] = useState<FeedbackData>({
    isHelpful: null,
  });
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (feedback.isHelpful === null) {
      Toast.warning('请选择是否有帮助');
      return;
    }

    setSubmitting(true);
    try {
      await onSubmit({
        conversation_id: conversationId,
        bot_id: botId,
        ...feedback,
      });
      Toast.success('感谢您的反馈！');
      onClose();
    } catch (error) {
      Toast.error('提交失败，请稍后重试');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Modal
      title="您的反馈对我们很重要"
      visible={visible}
      onCancel={onClose}
      footer={
        <div>
          <Button onClick={onClose}>取消</Button>
          <Button
            type="primary"
            onClick={handleSubmit}
            loading={submitting}
          >
            提交反馈
          </Button>
        </div>
      }
    >
      <div className="feedback-panel">
        {/* 快速反馈 */}
        <div className="feedback-section">
          <label>这个回答是否有帮助？</label>
          <Radio.Group
            value={feedback.isHelpful}
            onChange={(e) =>
              setFeedback({ ...feedback, isHelpful: e.target.value })
            }
            type="button"
            buttonStyle="solid"
          >
            <Radio value={true}>👍 有帮助</Radio>
            <Radio value={false}>👎 没帮助</Radio>
          </Radio.Group>
        </div>

        {/* 星级评分 */}
        {feedback.isHelpful !== null && (
          <div className="feedback-section">
            <label>请为本次对话打分</label>
            <Rate
              value={feedback.rating}
              onChange={(value) =>
                setFeedback({ ...feedback, rating: value })
              }
            />
          </div>
        )}

        {/* 问题分类 */}
        {feedback.isHelpful === false && (
          <div className="feedback-section">
            <label>请问主要是什么问题？</label>
            <Select
              value={feedback.category}
              onChange={(value) =>
                setFeedback({ ...feedback, category: value })
              }
              placeholder="请选择问题类型"
              style={{ width: '100%' }}
            >
              <Select.Option value="accurate">
                回答不准确
              </Select.Option>
              <Select.Option value="incomplete">
                回答不完整
              </Select.Option>
              <Select.Option value="irrelevant">
                回答不相关
              </Select.Option>
              <Select.Option value="slow">
                响应速度慢
              </Select.Option>
              <Select.Option value="other">
                其他问题
              </Select.Option>
            </Select>
          </div>
        )}

        {/* 文字反馈 */}
        <div className="feedback-section">
          <label>详细说明（可选）</label>
          <TextArea
            value={feedback.text}
            onChange={(e) =>
              setFeedback({ ...feedback, text: e.target.value })
            }
            placeholder="请描述您遇到的问题或建议..."
            maxCount={500}
            rows={4}
          />
        </div>
      </div>
    </Modal>
  );
};
```

### 7.2 反馈分析仪表盘

```typescript
// components/feedback/FeedbackAnalyticsDashboard.tsx
import React, { useEffect, useState } from 'react';
import { Card, Metric, AreaChart, Table, Tag } from '@douyinfe/semi-ui';

interface FeedbackStats {
  total_feedback: number;
  helpful_rate: number;
  avg_rating: number;
  sentiment_distribution: {
    positive: number;
    neutral: number;
    negative: number;
  };
  daily_breakdown: Array<{
    date: string;
    total_feedback: number;
    helpful_rate: number;
    avg_rating: number;
  }>;
  category_breakdown: Record<string, number>;
}

export const FeedbackAnalyticsDashboard: React.FC<Props> = ({
  tenantId,
  dateRange,
  botId,
}) => {
  const [stats, setStats] = useState<FeedbackStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchFeedbackStats();
  }, [tenantId, dateRange, botId]);

  const fetchFeedbackStats = async () => {
    setLoading(true);
    const response = await fetch(
      `/api/v1/feedback/stats?tenant_id=${tenantId}&start_date=${dateRange.start}&end_date=${dateRange.end}&bot_id=${botId || ''}`
    );
    const data = await response.json();
    setStats(data);
    setLoading(false);
  };

  if (loading) return <Spin />;

  return (
    <div className="feedback-analytics-dashboard">
      {/* 关键指标 */}
      <div className="metrics-grid">
        <Card title="总反馈数">
          <Metric value={stats.total_feedback} />
        </Card>
        <Card title="有帮助率">
          <Metric value={stats.helpful_rate} suffix="%" />
        </Card>
        <Card title="平均评分">
          <Metric value={stats.avg_rating} suffix=" / 5" />
        </Card>
        <Card title="负面反馈数">
          <Metric value={stats.sentiment_distribution.negative} />
        </Card>
      </div>

      {/* 情感分布饼图 */}
      <Card title="情感分布">
        <PieChart
          data={[
            { name: '积极', value: stats.sentiment_distribution.positive },
            { name: '中性', value: stats.sentiment_distribution.neutral },
            { name: '消极', value: stats.sentiment_distribution.negative },
          ]}
        />
      </Card>

      {/* 每日趋势图 */}
      <Card title="反馈趋势">
        <AreaChart
          data={stats.daily_breakdown}
          xField="date"
          yField="helpful_rate"
          yAxis={{ label: { formatter: (v) => `${v}%` } }}
        />
      </Card>

      {/* 问题分类表 */}
      <Card title="问题分类统计">
        <Table
          columns={[
            { title: '分类', dataKey: 'category' },
            { title: '数量', dataKey: 'count' },
            { title: '占比', dataKey: 'percent', render: (v) => `${v}%` },
          ]}
          dataSource={Object.entries(stats.category_breakdown).map(
            ([category, count]) => ({
              category,
              count,
              percent: ((count / stats.total_feedback) * 100).toFixed(1),
            })
          )}
          pagination={false}
        />
      </Card>

      {/* Top问题列表 */}
      <TopIssuesList tenantId={tenantId} botId={botId} />
    </div>
  );
};

// Top问题列表子组件
const TopIssuesList: React.FC<{ tenantId: string; botId?: string }> = ({
  tenantId,
  botId,
}) => {
  const [issues, setIssues] = useState([]);

  useEffect(() => {
    fetchTopIssues();
  }, [tenantId, botId]);

  const fetchTopIssues = async () => {
    const response = await fetch(
      `/api/v1/feedback/issues?tenant_id=${tenantId}&bot_id=${botId || ''}&limit=10`
    );
    const data = await response.json();
    setIssues(data);
  };

  return (
    <Card title="Top 10 问题">
      <Table
        columns={[
          { title: '排名', dataKey: 'rank', width: 60 },
          {
            title: '问题描述',
            dataKey: 'issue_title',
          },
          { title: '提及次数', dataKey: 'mention_count' },
          {
            title: '优先级',
            dataKey: 'priority',
            render: (priority) => {
              const colorMap = {
                critical: 'red',
                high: 'orange',
                medium: 'yellow',
                low: 'green',
              };
              return <Tag color={colorMap[priority]}>{priority}</Tag>;
            },
          },
          {
            title: '趋势',
            dataKey: 'trend',
            render: (trend) => {
              const iconMap = {
                increasing: '📈',
                stable: '➡️',
                decreasing: '📉',
              };
              return iconMap[trend];
            },
          },
        ]}
        dataSource={issues.map((issue, index) => ({
          ...issue,
          rank: index + 1,
        }))}
        pagination={false}
      />
    </Card>
  );
};
```

---

## 8. 实施计划

### 8.1 开发阶段划分

| 阶段 | 任务 | 工期 | 交付物 |
|------|------|------|--------|
| **第1周** | 数据库设计与创建 | 2天 | 4张表 |
| | 反馈收集服务开发 | 2天 | Go服务 |
| | AI分析引擎开发 | 1天 | 分析器 |
| **第2周** | 前端组件开发 | 3天 | 3个组件 |
| | API集成测试 | 2天 | 测试报告 |

### 8.2 技术风险与缓解

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|----------|
| AI分析准确性 | 中 | 中 | 定期人工校正，微调提示词 |
| 用户参与度低 | 高 | 中 | 简化反馈流程，提供激励 |
| 邮件调研被屏蔽 | 中 | 中 | 提供站内调研备选方案 |

### 8.3 成功指标

- ✅ 反馈收集率: **> 30%** 的对话有反馈
- ✅ 平均评分: **≥ 4.2/5.0**
- ✅ 有助率: **≥ 85%**
- ✅ NPS分数: **≥ 50 (良好级别)**
- ✅ 问题响应时间: **Top问题 48小时内生成改进任务**

---

**文档结束**
