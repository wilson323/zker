# ZKER 核心算法实现指南

> **文档类型**: 技术规范 / 算法实现
> **版本**: v1.0
> **生成日期**: 2025-01-01
> **目标读者**: 后端开发人员、架构师
> **使用目标**: 确保核心算法实现标准化、可维护

---

## 📋 目录

1. [文档说明](#文档说明)
2. [智能路由引擎算法](#智能路由引擎算法)
3. [RBAC 权限判断算法](#rbac-权限判断算法)
4. [多租户数据隔离算法](#多租户数据隔离算法)
5. [工作流执行引擎算法](#工作流执行引擎算法)
6. [知识库检索算法](#知识库检索算法)
7. [订阅配额检查算法](#订阅配额检查算法)
8. [数据脱敏算法](#数据脱敏算法)

---

## 文档说明

### 1.1 文档目标

**核心目标**: 提供核心算法的详细伪代码,确保:

✅ **实现标准化**: 所有开发人员对算法有统一理解
✅ **可维护性**: 算法逻辑清晰,易于维护和优化
✅ **可测试性**: 算法步骤明确,便于编写测试用例
✅ **可复用性**: 算法设计为可复用组件

### 1.2 伪代码规范

本文档使用的伪代码规范:

- **变量命名**: camelCase (如 `userInput`, `confidenceScore`)
- **常量命名**: UPPER_SNAKE_CASE (如 `MAX_RETRY_COUNT`)
- **函数命名**: PascalCase 或 camelCase (如 `RouteToService` 或 `routeToService`)
- **注释**: 使用 `//` 单行注释或 `/* */` 多行注释
- **数据结构**: 明确指定类型

### 1.3 算法分类

**按功能分类**:
- **路由决策**: 智能路由引擎相关算法
- **权限控制**: RBAC 权限判断算法
- **数据隔离**: 多租户数据隔离算法
- **执行引擎**: 工作流执行引擎
- **数据检索**: 知识库检索算法
- **配额管理**: 订阅配额检查算法
- **数据安全**: 数据脱敏算法

---

## 智能路由引擎算法

### 2.1 意图识别算法 (Intent Recognition)

**算法名称**: `HybridIntentMatcher`

**功能说明**: 识别用户输入的意图,返回 Top-K 个意图及置信度

**位置**: `backend/service/routing/intent_matcher.go`

#### 2.1.1 输入输出

**输入**:
```go
type IntentMatchInput struct {
    UserInput            string              // 用户输入
    Context              *RoutingContext     // 路由上下文
    HistoricalConversations []Conversation   // 历史会话
}
```

**输出**:
```go
type IntentMatchOutput struct {
    TopIntents          []IntentPrediction   // Top-K 意图及置信度
    ConfidenceScore     float64             // 最高置信度
    MatchedEntities     []Entity            // 识别的实体
    MatchedKeywords     []string            // 匹配的关键词
}
```

**意图预测**:
```go
type IntentPrediction struct {
    IntentName     string  // 意图名称
    Confidence     float64 // 置信度 (0-1)
    MatchMethod    string  // 匹配方法 (rule/similarity/model)
}
```

#### 2.1.2 算法步骤

```text
算法 HybridIntentMatcher:
输入: input (IntentMatchInput)
输出: output (IntentMatchOutput)

1. ========== 预处理阶段 (Preprocessing) ==========

   1.1 文本清洗 (Text Cleaning)
       cleanedText = input.UserInput

       // 去除特殊字符
       cleanedText = RemoveSpecialChars(cleanedText)

       // 转小写
       cleanedText = ToLower(cleanedText)

       // 去除多余空格
       cleanedText = TrimSpaces(cleanedText)

   1.2 分词 (Tokenization)
       tokens = []

       IF IsChineseText(cleanedText):
           // 中文分词: 使用 jieba
           tokens = JiebaCut(cleanedText)
       ELSE:
           // 英文分词: 使用 NLTK
           tokens = NLTKWordTokenize(cleanedText)

   1.3 去除停用词 (Stop Words Removal)
       tokens = RemoveStopWords(tokens, StopWordsList)

2. ========== 特征提取阶段 (Feature Extraction) ==========

   2.1 TF-IDF 特征
       tfidfVector = CalculateTFIDF(tokens, Corpus)

   2.2 词嵌入特征 (Word Embedding) - 可选
       IF UseEmbeddingFeature:
           // 调用 text-embedding-ada-002
           embeddingVector = GetEmbedding(cleanedText)
           // 维度: 1536
       ELSE:
           embeddingVector = null

3. ========== 意图匹配阶段 (Intent Matching) ==========

   candidates = []  // 候选意图列表

   3.1 基于规则的匹配 (Rule-Based Matching)
       FOR EACH intentRule IN IntentRules:

           // 检查规则是否匹配
           isMatch, extractedParams = MatchRule(
               cleanedText,
               intentRule.Pattern
           )

           IF isMatch:
               ADD {
                   IntentName: intentRule.IntentName,
                   Confidence: 1.0,  // 规则匹配置信度为1
                   MatchMethod: "rule",
                   Params: extractedParams
               } TO candidates
       END FOR

   3.2 基于相似度的匹配 (Similarity-Based Matching)
       FOR EACH intentSample IN IntentSamples:

           // 计算余弦相似度
           similarity = CosineSimilarity(
               embeddingVector,
               intentSample.EmbeddingVector
           )

           // 相似度阈值过滤
           IF similarity > SIMILARITY_THRESHOLD:  // 0.7
               ADD {
                   IntentName: intentSample.IntentName,
                   Confidence: similarity,
                   MatchMethod: "similarity"
               } TO candidates
       END FOR

   3.3 基于模型的匹配 (Model-Based Matching) - 可选
       IF UseModelMatching:

           // 调用意图识别模型
           modelResult = IntentModel.Predict(cleanedText)

           FOR EACH intent IN modelResult.Intents:
               IF intent.Confidence > MODEL_CONFIDENCE_THRESHOLD:  // 0.7
                   ADD {
                       IntentName: intent.Name,
                       Confidence: intent.Confidence,
                       MatchMethod: "model"
                   } TO candidates
           END FOR
       END IF

4. ========== 结果聚合阶段 (Aggregation) ==========

   4.1 去重 (Deduplication)
       // 按意图名称去重,保留最高置信度
       uniqueCandidates = DeduplicateByIntentName(candidates)

   4.2 置信度加权 (Confidence Weighting)
       FOR EACH candidate IN uniqueCandidates:

           SWITCH candidate.MatchMethod:
               CASE "rule":
                   finalWeight = 0.3
               CASE "similarity":
                   finalWeight = 0.4
               CASE "model":
                   finalWeight = 0.3

           candidate.FinalConfidence = candidate.Confidence * finalWeight
       END FOR

   4.3 排序 (Sorting)
       uniqueCandidates.SortBy("FinalConfidence", DESC)

   4.4 取 Top-K
       topK = uniqueCandidates[0:K]  // K=3

5. ========== 实体提取阶段 (Entity Extraction) ==========

   5.1 正则表达式提取
       regexEntities = ExtractEntitiesByRegex(cleanedText, EntityPatterns)

   5.2 NER 模型提取 - 可选
       IF UseNERModel:
           nerEntities = NERModel.Extract(cleanedText)
       ELSE:
           nerEntities = []

   5.3 合并实体
       matchedEntities = MergeEntities(regexEntities, nerEntities)

6. ========== 关键词提取阶段 (Keyword Extraction) ==========

   // 提取匹配到的关键词
   matchedKeywords = ExtractMatchedKeywords(
       cleanedText,
       topK,
       IntentKeywords
   )

7. ========== 返回结果 ==========

   RETURN {
       TopIntents: topK,
       ConfidenceScore: topK[0].FinalConfidence,
       MatchedEntities: matchedEntities,
       MatchedKeywords: matchedKeywords
   }
```

#### 2.1.3 复杂度分析

**时间复杂度**:
- 预处理: O(n), n = 输入文本长度
- 特征提取: O(n * m), m = 特征维度
- 规则匹配: O(r), r = 规则数量
- 相似度匹配: O(s * m), s = 意图样本数量
- 模型匹配: O(model_time)
- **总复杂度**: O(n * m + s * m)

**空间复杂度**:
- O(m + s), m = 特征维度, s = 意图样本数量

#### 2.1.4 参数配置

```go
const (
    // 相似度阈值
    SIMILARITY_THRESHOLD = 0.7

    // 模型置信度阈值
    MODEL_CONFIDENCE_THRESHOLD = 0.7

    // Top-K
    DEFAULT_TOP_K = 3

    // 权重分配
    RULE_WEIGHT = 0.3
    SIMILARITY_WEIGHT = 0.4
    MODEL_WEIGHT = 0.3

    // 是否使用词嵌入
    USE_EMBEDDING_FEATURE = true

    // 是否使用模型匹配
    USE_MODEL_MATCHING = false
)
```

#### 2.1.5 异常处理

```go
// 异常处理
func (m *HybridIntentMatcher) Match(input *IntentMatchInput) (*IntentMatchOutput, error) {
    // 1. 输入验证
    if strings.TrimSpace(input.UserInput) == "" {
        // 返回默认意图
        return m.getDefaultIntent(), nil
    }

    // 2. 意图样本为空
    if len(m.intentSamples) == 0 {
        return m.getDefaultIntent(), nil
    }

    // 3. 模型调用失败
    if USE_MODEL_MATCHING {
        result, err := m.modelPredict(input.UserInput)
        if err != nil {
            // 降级到规则+相似度匹配
            log.Warn("Model predict failed, fallback to rule+similarity", "error", err)
            // 继续执行,跳过模型匹配步骤
        }
    }

    // 4. 其他异常
    defer func() {
        if r := recover(); r != nil {
            log.Error("Intent matching panic", "error", r)
        }
    }()

    return output, nil
}
```

#### 2.1.6 性能优化

**优化策略**:

1. **向量检索加速**: 使用 Faiss 加速相似度计算
2. **规则预编译**: 预编译正则表达式
3. **并行计算**: 规则匹配和相似度匹配并行执行
4. **缓存策略**: 缓存常见输入的结果

```go
// 优化示例: 使用 Faiss 加速
func (m *HybridIntentMatcher) similaritySearchFast(embedding []float32) []IntentPrediction {
    // 使用 Faiss IVF 索引
    labels, distances := m.faissIndex.Search(embedding, TOP_K)

    results := make([]IntentPrediction, len(labels))
    for i, label := range labels {
        results[i] = IntentPrediction{
            IntentName: m.intentSamples[label].IntentName,
            Confidence: 1.0 - distances[i], // 距离转相似度
            MatchMethod: "similarity",
        }
    }

    return results
}
```

---

### 2.2 路由决策算法 (Routing Decision)

**算法名称**: `ScoreBasedRouter`

**功能说明**: 基于多因素评分选择最优服务

**位置**: `backend/service/routing/router.go`

#### 2.2.1 输入输出

**输入**:
```go
type RoutingDecisionInput struct {
    Intents            []IntentPrediction   // Top-K 意图
    Context            *RoutingContext      // 路由上下文
    ServiceCandidates  []*Service           // 候选服务
}
```

**输出**:
```go
type RoutingDecisionOutput struct {
    SelectedService    *Service            // 选中的服务
    DecisionID         string              // 决策 ID
    DecisionReason     string              // 决策原因
    RoutingMetadata    *RoutingMetadata    // 路由元数据
}
```

#### 2.2.2 算法步骤

```text
算法 ScoreBasedRouter:
输入: input (RoutingDecisionInput)
输出: output (RoutingDecisionOutput)

1. ========== 候选服务生成 (Candidate Generation) ==========

   candidates = []

   1.1 根据意图匹配服务
       FOR EACH intent IN input.Intents:

           // 查找支持该意图的服务
           matchedServices = ServicesByIntent[intent.IntentName]

           FOR EACH service IN matchedServices:
               ADD service TO candidates
           END FOR
       END FOR

   1.2 过滤不可用服务
       availableCandidates = []

       FOR EACH service IN candidates:

           // 检查服务状态
           IF service.Status != "available":
               CONTINUE
           END IF

           // 检查健康分数
           IF service.HealthScore < HEALTH_SCORE_THRESHOLD:  // 0.5
               CONTINUE
           END IF

           ADD service TO availableCandidates
       END FOR

   // 如果没有可用服务,返回默认服务
   IF LEN(availableCandidates) == 0:
       RETURN {
           SelectedService: DefaultService,
           DecisionReason: "No available service, use default"
       }
   END IF

   1.3 权限过滤
       authorizedCandidates = []

       FOR EACH service IN availableCandidates:

           // 检查用户权限
           IF UserHasPermission(input.Context.UserID, service):
               ADD service TO authorizedCandidates
           END IF
       END FOR

   // 如果没有权限访问的服务,返回权限错误
   IF LEN(authorizedCandidates) == 0:
       RETURN error, "No permission to access any service"
   END IF

2. ========== 服务评分 (Service Scoring) ==========

   FOR EACH service IN authorizedCandidates:

       totalScore = 0.0

       // 2.1 意图匹配度 (权重 0.3)
       intentMatchScore = input.Intents[0].FinalConfidence  // 使用第一意图
       totalScore += 0.3 * intentMatchScore

       // 2.2 服务负载 (权重 0.2)
       loadScore = 1.0 - (service.CurrentLoad / service.MaxLoad)
       totalScore += 0.2 * loadScore

       // 2.3 历史成功率 (权重 0.2)
       successScore = service.SuccessRate  // 最近1小时的成功率
       totalScore += 0.2 * successScore

       // 2.4 地域亲和性 (权重 0.1)
       IF service.Region == input.Context.UserRegion:
           totalScore += 0.1
       END IF

       // 2.5 成本考虑 (权重 0.1)
       // 归一化成本到 [0,1]
       costScore = 1.0 - NormalizeCost(service.CostPerRequest)
       totalScore += 0.1 * costScore

       // 2.6 A/B 测试分流 (权重 0.1) - 如果启用
       IF ABTestEnabled:
           abWeight = GetABTestWeight(service, input.Context.UserID)
           totalScore += 0.1 * abWeight
       END IF

       // 保存最终分数
       service.FinalScore = totalScore

       // 记录评分详情
       service.ScoreDetails = {
           IntentMatch: intentMatchScore,
           Load: loadScore,
           Success: successScore,
           Region: service.Region == input.Context.UserRegion ? 1.0 : 0.0,
           Cost: costScore,
           ABTest: abWeight IF ABTestEnabled ELSE 0.0
       }
   END FOR

3. ========== 选择最优服务 (Select Best Service) ==========

   // 按分数排序
   sortedCandidates = authorizedCandidates.SortBy("FinalScore", DESC)

   selectedService = sortedCandidates[0]

   // 检查分数是否超过阈值
   IF selectedService.FinalScore < SCORE_THRESHOLD:  // 0.5
       // 分数过低,选择默认服务
       selectedService = DefaultService
       decisionReason = "All candidates score below threshold"
   ELSE:
       decisionReason = ExplainRouting(selectedService)
   END IF

4. ========== 记录决策日志 (Log Decision) ==========

   decisionID = GenerateUUID()

   LogRoutingDecision({
       DecisionID: decisionID,
       UserID: input.Context.UserID,
       TenantID: input.Context.TenantID,
       InputIntents: input.Intents,
       Candidates: authorizedCandidates,
       SelectedService: selectedService,
       DecisionReason: decisionReason,
       Timestamp: Now()
   })

5. ========== A/B 测试记录 (Record A/B Test) ==========

   IF ABTestEnabled:
       RecordABTestEvent({
           UserID: input.Context.UserID,
           ServiceID: selectedService.ServiceID,
           ExperimentID: GetCurrentExperimentID(),
           VariantID: GetVariant(selectedService),
           Timestamp: Now()
       })
   END IF

6. ========== 返回结果 ==========

   RETURN {
       SelectedService: selectedService,
       DecisionID: decisionID,
       DecisionReason: decisionReason,
       RoutingMetadata: {
           Score: selectedService.FinalScore,
           ScoreDetails: selectedService.ScoreDetails,
           AlternativeServices: sortedCandidates[1:3]  // 次优的2个服务
       }
   }
```

#### 2.2.3 决策解释 (Explain Routing)

```go
// 决策解释函数
func ExplainRouting(service *Service) string {
    reasons := []string{}

    if service.ScoreDetails.IntentMatch > 0.8 {
        reasons = append(reasons, "高度匹配用户意图")
    }

    if service.ScoreDetails.Load > 0.7 {
        reasons = append(reasons, "服务负载较低")
    }

    if service.ScoreDetails.Success > 0.9 {
        reasons = append(reasons, "历史成功率高")
    }

    if service.ScoreDetails.Region > 0 {
        reasons = append(reasons, "地域亲和性好")
    }

    if service.ScoreDetails.Cost > 0.7 {
        reasons = append(reasons, "成本较低")
    }

    return strings.Join(reasons, "; ")
}
```

#### 2.2.4 复杂度分析

**时间复杂度**:
- O(n + m log m), n = 候选服务数量, m = 排序的服务数量
- 主要是排序的复杂度

**空间复杂度**:
- O(m), m = 候选服务数量

#### 2.2.5 参数配置

```go
const (
    // 健康分数阈值
    HEALTH_SCORE_THRESHOLD = 0.5

    // 最终分数阈值
    SCORE_THRESHOLD = 0.5

    // 权重分配
    INTENT_MATCH_WEIGHT = 0.3
    LOAD_WEIGHT = 0.2
    SUCCESS_WEIGHT = 0.2
    REGION_WEIGHT = 0.1
    COST_WEIGHT = 0.1
    AB_TEST_WEIGHT = 0.1

    // 是否启用 A/B 测试
    AB_TEST_ENABLED = false
)
```

---

### 2.3 负载均衡算法 (Load Balancing)

**算法名称**: `WeightedRoundRobin`

**功能说明**: 加权轮询负载均衡

#### 2.3.1 算法步骤

```text
算法 WeightedRoundRobin:
输入: services (Service列表)
输出: selectedService

1. 初始化
   // 每个服务有两个权重:
   // - currentWeight: 当前权重 (动态)
   // - effectiveWeight: 有效权重 (固定)

2. 选择过程
   totalWeight = SUM(service.effectiveWeight FOR EACH service IN services)
   bestService = null
   maxCurrentWeight = -1

   FOR EACH service IN services:
       // 当前权重 += 有效权重
       service.currentWeight += service.effectiveWeight

       // 找到当前权重最高的服务
       IF service.currentWeight > maxCurrentWeight:
           maxCurrentWeight = service.currentWeight
           bestService = service
       END IF
   END FOR

   // 减少选中服务的当前权重
   bestService.currentWeight -= totalWeight

   RETURN bestService
```

---

## RBAC 权限判断算法

### 3.1 数据权限判断算法 (Row-Level Permission)

**算法名称**: `RowLevelPermissionChecker`

**功能说明**: 检查用户是否有权限访问某行数据

**位置**: `backend/service/rbac/permission_checker.go`

#### 3.1.1 输入输出

**输入**:
```go
type RowLevelPermissionInput struct {
    UserID       string
    ResourceType  string  // 资源类型 (如 "bots", "conversations")
    ResourceID   string  // 资源 ID
    Action       string  // 操作 (READ, WRITE, DELETE)
}
```

**输出**:
```go
type RowLevelPermissionOutput struct {
    Allowed       bool
    Reason        string
    PermissionID  string  // 匹配的权限规则 ID
}
```

#### 3.1.2 算法步骤

```text
算法 RowLevelPermissionChecker:
输入: input (RowLevelPermissionInput)
输出: output (RowLevelPermissionOutput)

1. ========== 获取用户角色 ==========

   userRoles = GetUserRoles(input.UserID)

   // 检查是否为系统管理员
   IF IsSystemAdmin(userRoles):
       RETURN {
           Allowed: true,
           Reason: "System admin has all permissions"
       }
   END IF

2. ========== 获取资源权限规则 ==========

   permissionRules = GetDataPermissions(
       input.ResourceType,
       input.Action,
       userRoles
   )

   // 如果没有规则,默认拒绝
   IF LEN(permissionRules) == 0:
       RETURN {
           Allowed: false,
           Reason: "No permission rule found"
       }
   END IF

3. ========== 检查权限级别 ==========

   // 按优先级排序: 自定义 > 部门 > 组织 > 自己 > 全部
   sortedRules = permissionRules.SortBy("Priority", DESC)

   FOR EACH rule IN sortedRules:

       SWITCH rule.Scope:

           CASE "ALL":
               // 所有数据
               RETURN {
                   Allowed: true,
                   Reason: "Has permission to all data",
                   PermissionID: rule.PermissionID
               }

           CASE "OWN":
               // 仅自己的数据
               resource = GetResource(input.ResourceType, input.ResourceID)

               IF resource.CreatedBy == input.UserID:
                   RETURN {
                       Allowed: true,
                       Reason: "Has permission to own data",
                       PermissionID: rule.PermissionID
                   }
               ELSE:
                   RETURN {
                       Allowed: false,
                       Reason: "No permission to others data",
                       PermissionID: rule.PermissionID
                   }
               END IF

           CASE "DEPARTMENT":
               // 本部门数据
               resource = GetResource(input.ResourceType, input.ResourceID)
               userDept = GetUserDepartment(input.UserID)

               IF resource.DepartmentID == userDept.ID:
                   RETURN {
                       Allowed: true,
                       Reason: "Has permission to department data",
                       PermissionID: rule.PermissionID
                   }
               ELSE:
                   CONTINUE  // 继续检查下一个规则
               END IF

           CASE "ORGANIZATION":
               // 本组织数据
               resource = GetResource(input.ResourceType, input.ResourceID)
               userOrg = GetUserOrganization(input.UserID)

               IF resource.OrganizationID == userOrg.ID:
                   RETURN {
                       Allowed: true,
                       Reason: "Has permission to organization data",
                       PermissionID: rule.PermissionID
                   }
               ELSE:
                   CONTINUE
               END IF

           CASE "CUSTOM":
               // 自定义规则
               IF EvaluateCustomRule(rule, input):
                   RETURN {
                       Allowed: true,
                       Reason: "Custom permission rule passed",
                       PermissionID: rule.PermissionID
                   }
               ELSE:
                   CONTINUE
               END IF

       END SWITCH
   END FOR

4. ========== 所有规则都不匹配 ==========

   RETURN {
       Allowed: false,
       Reason: "No matching permission rule"
   }
```

#### 3.1.3 自定义规则评估

```go
// 自定义规则评估函数
func EvaluateCustomRule(rule *PermissionRule, input *RowLevelPermissionInput) bool {
    // 解析自定义规则表达式
    // 例如: "resource.created_by == user.id OR resource.tags contains 'public'"

    expr, err := ParseExpression(rule.Expression)
    if err != nil {
        log.Error("Failed to parse custom rule", "error", err)
        return false
    }

    // 准备上下文
    context := map[string]interface{}{
        "user": map[string]interface{}{
            "id":     input.UserID,
            "roles":  GetUserRoles(input.UserID),
            "dept":   GetUserDepartment(input.UserID),
            "org":    GetUserOrganization(input.UserID),
        },
        "resource": map[string]interface{}{
            "id":     input.ResourceID,
            "type":   input.ResourceType,
        },
    }

    // 执行表达式
    result, err := expr.Evaluate(context)
    if err != nil {
        log.Error("Failed to evaluate custom rule", "error", err)
        return false
    }

    return result.(bool)
}
```

---

### 3.2 字段权限判断算法 (Field-Level Permission)

**算法名称**: `FieldPermissionChecker`

**功能说明**: 检查用户是否有权限访问某字段

#### 3.2.1 算法步骤

```text
算法 FieldPermissionChecker:
输入: userID, resourceType, fieldName
输出: allowed, reason

1. 获取用户角色
   userRoles = GetUserRoles(userID)

2. 获取字段权限规则
   fieldPermissions = GetFieldPermissions(
       resourceType,
       fieldName,
       userRoles
   )

3. 检查权限
   IF LEN(fieldPermissions) == 0:
       // 没有规则,默认允许
       RETURN true, "No field permission rule, allow by default"
   END IF

4. 检查权限级别
   // 优先级: deny > allow
   hasDeny = false
   hasAllow = false

   FOR EACH permission IN fieldPermissions:
       IF permission.Effect == "deny":
           hasDeny = true
       ELSE IF permission.Effect == "allow":
           hasAllow = true
       END IF
   END FOR

   IF hasDeny:
       RETURN false, "Denied by field permission rule"
   ELSE IF hasAllow:
       RETURN true, "Allowed by field permission rule"
   ELSE:
       RETURN false, "No allow permission rule"
   END IF
```

---

## 多租户数据隔离算法

### 4.1 租户数据查询算法

**算法名称**: `TenantIsolatedQuery`

**功能说明**: 自动添加租户过滤条件,确保数据隔离

**位置**: `backend/pkg/tenant/query_isolation.go`

#### 4.1.1 算法步骤

```text
算法 TenantIsolatedQuery:
输入: query (原始查询), table (表名), context (上下文)
输出: isolatedQuery (添加租户过滤的查询)

1. ========== 获取用户租户 ==========

   user = GetUserFromContext(context)
   tenantID = user.TenantID

   // 检查租户是否有效
   IF !IsTenantActive(tenantID):
       RETURN error, "Tenant is not active or does not exist"
   END IF

2. ========== 检查是否为系统管理员 ==========

   IF user.IsSystemAdmin:
       // 系统管理员可以跨租户查询
       // 但需要记录日志
       LogCrossTenantAccess(user.UserID, query)
       RETURN query, "System admin bypasses tenant isolation"
   END IF

3. ========== 添加租户过滤 ==========

   // 检查表是否有 tenant_id 字段
   IF !TableHasField(table, "tenant_id"):
       // 表没有 tenant_id 字段,可能是系统表
       LogWarning("Table has no tenant_id field", "table", table)
       RETURN query, "Table has no tenant_id field, skip isolation"
   END IF

   // 添加 WHERE 条件
   IF query.HasWHERE():
       query.WHERE.Add(f"{table}.tenant_id = '{tenantID}'")
   ELSE:
       query.WHERE = f"{table}.tenant_id = '{tenantID}'"
   END IF

4. ========== 检查软删除过滤 ==========

   IF TableHasField(table, "deleted_at"):
       // 添加软删除过滤
       query.WHERE.Add(f"{table}.deleted_at IS NULL")
   END IF

5. ========== 返回结果 ==========

   RETURN query, tenantID
```

#### 4.1.2 安全加固

```go
// 安全加固: 强制租户隔离
func TenantIsolationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取用户
        user := GetUserFromContext(c)

        // 获取租户
        tenantID := user.TenantID

        // 检查租户状态
        tenant, err := GetTenant(tenantID)
        if err != nil || tenant.Status != "active" {
            c.JSON(403, gin.H{
                "code": "TARG300",
                "message": "租户不存在或已停用",
            })
            c.Abort()
            return
        }

        // 将租户 ID 存储到上下文
        c.Set("tenant_id", tenantID)

        // 继续处理请求
        c.Next()

        // 记录跨租户访问尝试
        if c.Writer.Status() >= 400 {
            LogCrossTenantAccessAttempt(user.UserID, tenantID, c.Request.URL.Path)
        }
    }
}
```

---

## 工作流执行引擎算法

### 5.1 DAG 执行算法

**算法名称**: `WorkflowDAGExecutor`

**功能说明**: 执行 DAG 工作流,支持节点并行执行

**位置**: `backend/service/workflow/executor.go`

#### 5.1.1 输入输出

**输入**:
```go
type WorkflowExecuteInput struct {
    Workflow     *Workflow
    InputData    map[string]interface{}
    Context      *ExecutionContext
}
```

**输出**:
```go
type WorkflowExecuteOutput struct {
    ExecutionID  string
    Status       string
    Output       map[string]interface{}
    Duration     time.Duration
    NodeResults  map[string]*NodeResult
}
```

#### 5.1.2 算法步骤

```text
算法 WorkflowDAGExecutor:
输入: input (WorkflowExecuteInput)
输出: output (WorkflowExecuteOutput)

1. ========== 初始化 ==========

   executionID = GenerateUUID()

   execution = {
       ExecutionID: executionID,
       WorkflowID: input.Workflow.ID,
       Status: "running",
       StartTime: Now(),
       InputData: input.InputData,
       NodeResults: {}
   }

   SaveExecution(execution)

2. ========== 构建 DAG ==========

   dag = BuildDAG(input.Workflow.Nodes, input.Workflow.Edges)

   // 检查是否有环
   IF HasCycle(dag):
       execution.Status = "failed"
       execution.Error = "Workflow contains cycle"
       SaveExecution(execution)
       RETURN error, "Workflow contains cycle"
   END IF

3. ========== 拓扑排序 ==========

   try:
       sortedNodes = TopologicalSort(dag)
   except CycleDetected:
       execution.Status = "failed"
       execution.Error = "Cycle detected in workflow"
       SaveExecution(execution)
       RETURN error, "Cycle detected"

4. ========== 执行节点 ==========

   completedNodes = []  // 已完成的节点
   failedNodes = []     // 失败的节点
   nodeResults = {}     // 节点执行结果

   // 创建并行执行池
   pool = NewWorkerPool(MAX_CONCURRENT_NODES)  // 最大并发数

   FOR EACH node IN sortedNodes:

       // 4.1 检查前置节点是否完成
       IF NOT AllDependenciesCompleted(node, completedNodes):
           CONTINUE  // 等待前置节点完成
       END IF

       // 4.2 准备节点输入
       nodeInput = PrepareNodeInput(
           node,
           input.InputData,
           nodeResults
       )

       // 4.3 提交到执行池
       pool.Submit(func() {

           // 4.4 执行节点
           result = ExecuteNode(node, nodeInput, input.Context)

           // 4.5 记录执行结果
           nodeResults[node.ID] = result

           // 4.6 记录执行日志
           LogNodeExecution(executionID, node, result)

           // 4.7 检查是否失败
           IF result.Status == "failed":
               failedNodes.ADD(node)

               // 根据错误处理策略决定是否继续
               IF node.ErrorStrategy == "stop":
                   // 停止整个工作流
                   pool.StopAll()
                   CancelPendingNodes()
               ELSE IF node.ErrorStrategy == "continue":
                   // 继续执行其他节点
                   CONTINUE
               ELSE IF node.ErrorStrategy == "retry":
                   // 重试
                   FOR retryCount IN 1 TO node.MaxRetries:
                       result = ExecuteNode(node, nodeInput, input.Context)
                       IF result.Status == "success":
                           BREAK
                       END IF
                       WAIT(REtryDelay)  // 等待后重试
                   END FOR

                   // 重试仍失败
                   IF result.Status == "failed":
                       failedNodes.ADD(node)
                   END IF
               END IF
           ELSE:
               // 成功
               completedNodes.ADD(node)
           END IF
       })
   END FOR

   // 等待所有节点完成
   pool.WaitAll()

5. ========== 收集输出 ==========

   outputData = CollectWorkflowOutput(
       input.Workflow.OutputNodes,
       nodeResults
   )

6. ========== 更新执行状态 ==========

   execution.Status = "completed"
   execution.Output = outputData
   execution.NodeResults = nodeResults
   execution.EndTime = Now()
   execution.Duration = execution.EndTime - execution.StartTime

   SaveExecution(execution)

7. ========== 返回结果 ==========

   RETURN {
       ExecutionID: executionID,
       Status: execution.Status,
       Output: outputData,
       Duration: execution.Duration,
       NodeResults: nodeResults
   }
```

#### 5.1.3 节点输入准备

```go
// 准备节点输入
func PrepareNodeInput(node *Node, globalInput map[string]interface{}, nodeResults map[string]*NodeResult) map[string]interface{} {
    nodeInput := make(map[string]interface{})

    // 1. 添加全局输入
    for key, value := range globalInput {
        nodeInput[key] = value
    }

    // 2. 添加前置节点的输出
    for _, edge := range node.InputEdges {
        sourceNodeID := edge.Source
        sourceResult := nodeResults[sourceNodeID]

        if sourceResult != nil && sourceResult.Status == "success" {
            // 使用输出变量的别名
            if edge.Alias != "" {
                nodeInput[edge.Alias] = sourceResult.Output
            } else {
                nodeInput[sourceNodeID] = sourceResult.Output
            }
        }
    }

    return nodeInput
}
```

#### 5.1.4 并行优化

```go
// 使用工作池实现并行执行
type WorkerPool struct {
    maxWorkers int
    taskChan   chan func()
    wg         sync.WaitGroup
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
    pool := &WorkerPool{
        maxWorkers: maxWorkers,
        taskChan:   make(chan func(), maxWorkers*2),
    }

    // 启动工作协程
    for i := 0; i < maxWorkers; i++ {
        pool.wg.Add(1)
        go pool.worker()
    }

    return pool
}

func (p *WorkerPool) worker() {
    defer p.wg.Done()

    for task := range p.taskChan {
        task()
    }
}

func (p *WorkerPool) Submit(task func()) {
    p.taskChan <- task
}

func (p *WorkerPool) WaitAll() {
    close(p.taskChan)
    p.wg.Wait()
}
```

---

## 知识库检索算法

### 6.1 混合检索算法

**算法名称**: `HybridKnowledgeRetriever`

**功能说明**: 结合向量检索和关键词检索,返回最相关的文档

**位置**: `backend/service/knowledge/retriever.go`

#### 6.1.1 输入输出

**输入**:
```go
type KnowledgeRetrieveInput struct {
    Query              string
    KnowledgeBaseID    string
    TopK              int
    ScoreThreshold     float64
}
```

**输出**:
```go
type KnowledgeRetrieveOutput struct {
    Chunks            []*KnowledgeChunk
    Scores            []float64
    TotalRetrieved    int
    RetrievalTime     time.Duration
}
```

#### 6.1.2 算法步骤

```text
算法 HybridKnowledgeRetriever:
输入: input (KnowledgeRetrieveInput)
输出: output (KnowledgeRetrieveOutput)

1. ========== 查询预处理 ==========

   cleanedQuery = PreprocessText(input.Query)
   queryEmbedding = GetEmbedding(cleanedQuery)  // 1536 维向量

2. ========== 向量检索 (Semantic Search) ==========

   // 2.1 Milvus 向量检索
   milvusResult = MilvusSearch({
       Collection: input.KnowledgeBaseID,
       Vector: queryEmbedding,
       TopK: input.TopK * 2,  // 召回更多,后续重排序
       MetricType: "COSINE"
   })

   // 2.2 计算向量相似度分数
   vectorResults = []
   FOR EACH hit IN milvusResult.Hits:
       chunk = GetChunk(hit.ID)
       chunk.VectorScore = hit.Distance  // 余弦距离 (0-1)
       ADD chunk TO vectorResults
   END FOR

3. ========== 关键词检索 (Keyword Search) ==========

   // 3.1 Elasticsearch 全文检索
   esResult = ElasticsearchSearch({
       Index: input.KnowledgeBaseID,
       Query: cleanedQuery,
       TopK: input.TopK * 2,
       Analyzer: "ik_smart"  // 中文分词
   })

   // 3.2 计算 BM25 分数
   keywordResults = []
   FOR EACH hit IN esResult.Hits:
       chunk = GetChunk(hit.ID)
       chunk.KeywordScore = NormalizeBM25(hit.Score)  // 归一化到 [0,1]
       ADD chunk TO keywordResults
   END FOR

4. ========== 结果融合 (Result Fusion) ==========

   // 4.1 合并结果
   allResults = Merge(vectorResults, keywordResults)

   // 4.2 去重
   uniqueResults = DeduplicateByID(allResults)

   // 4.3 重排序 (Reciprocal Rank Fusion, RRF)
   FOR EACH chunk IN uniqueResults:

       // RRF 算法
       vectorRank := GetRank(chunk.ID, vectorResults)  // 向量检索排名
       keywordRank := GetRank(chunk.ID, keywordResults)  // 关键词检索排名

       // RRF 分数
       rrfScore := 0.0

       IF vectorRank > 0:
           rrfScore += 0.6 / (K + vectorRank)  // K=60
       END IF

       IF keywordRank > 0:
           rrfScore += 0.4 / (K + keywordRank)
       END IF

       chunk.FinalScore = rrfScore
   END FOR

   // 4.4 排序
   sortedResults = uniqueResults.SortBy("FinalScore", DESC)

   // 4.5 过滤阈值
   filteredResults = []
   FOR EACH chunk IN sortedResults:
       IF chunk.FinalScore >= input.ScoreThreshold:
           ADD chunk TO filteredResults
       ELSE:
           BREAK  // 后续分数更低,无需继续
       END IF
   END FOR

   // 4.6 取 Top-K
   topResults = filteredResults[0:input.TopK]

5. ========== 返回结果 ==========

   RETURN {
       Chunks: topResults,
       Scores: topResults.FinalScore,
       TotalRetrieved: LEN(topResults),
       RetrievalTime: Now() - startTime
   }
```

#### 6.1.3 RRF 算法说明

**Reciprocal Rank Fusion (RRF)**:
- 一种融合多个排序结果的算法
- 对每个排名位置计算倒数和
- 公式: `RRF = Σ (weight / (K + rank))`
- K 通常取 60

**优势**:
- 不依赖具体分数值
- 对异常值鲁棒
- 计算简单高效

---

## 订阅配额检查算法

### 7.1 多维配额检查算法

**算法名称**: `MultiDimensionalQuotaChecker`

**功能说明**: 检查用户操作是否超过租户配额

**位置**: `backend/service/subscription/quota_checker.go`

#### 7.1.1 输入输出

**输入**:
```go
type QuotaCheckInput struct {
    TenantID      string
    ResourceType  string
    Action        string  // CREATE, READ, UPDATE, DELETE
    Quantity      int
}
```

**输出**:
```go
type QuotaCheckOutput struct {
    Allowed          bool
    Reason           string
    RemainingQuota   int
    CurrentUsage     int
    Limit            int
}
```

#### 7.1.2 算法步骤

```text
算法 MultiDimensionalQuotaChecker:
输入: input (QuotaCheckInput)
输出: output (QuotaCheckOutput)

1. ========== 获取租户订阅 ==========

   subscription = GetActiveSubscription(input.TenantID)

   IF subscription == nil:
       RETURN {
           Allowed: false,
           Reason: "No active subscription"
       }
   END IF

2. ========== 获取订阅方案配额 ==========

   plan = GetSubscriptionPlan(subscription.PlanID)
   quotas = plan.Quotas[input.ResourceType]

   IF quotas == nil:
       // 没有限制
       RETURN {
           Allowed: true,
           Reason: "No quota limit for this resource"
       }
   END IF

3. ========== 获取当前使用量 ==========

   usage = GetTenantUsage(input.TenantID, input.ResourceType)

4. ========== 检查维度配额 ==========

   FOR EACH dimension IN quotas.Dimensions:

       // 4.1 时间维度配额
       IF dimension.Type == "time":

           currentUsage := 0

           SWITCH dimension.Period:
               CASE "hourly":
                   currentUsage = usage.LastHourCount
               CASE "daily":
                   currentUsage = usage.LastDayCount
               CASE "monthly":
                   currentUsage = usage.LastMonthCount
           END SWITCH

           IF currentUsage + input.Quantity > dimension.Limit:
               RETURN {
                   Allowed: false,
                   Reason: f"Exceeds {dimension.Period} quota: {currentUsage + input.Quantity}/{dimension.Limit}",
                   RemainingQuota: dimension.Limit - currentUsage,
                   CurrentUsage: currentUsage,
                   Limit: dimension.Limit
               }
           END IF

       // 4.2 总量维度配额
       ELSE IF dimension.Type == "total":

           currentUsage := usage.TotalCount

           IF currentUsage + input.Quantity > dimension.Limit:
               RETURN {
                   Allowed: false,
                   Reason: f"Exceeds total quota: {currentUsage + input.Quantity}/{dimension.Limit}",
                   RemainingQuota: dimension.Limit - currentUsage,
                   CurrentUsage: currentUsage,
                   Limit: dimension.Limit
               }
           END IF

       // 4.3 并发维度配额
       ELSE IF dimension.Type == "concurrent":

           currentUsage := usage.ConcurrentCount

           IF currentUsage + input.Quantity > dimension.Limit:
               RETURN {
                   Allowed: false,
                   Reason: f"Exceeds concurrent quota: {currentUsage + input.Quantity}/{dimension.Limit}",
                   RemainingQuota: dimension.Limit - currentUsage,
                   CurrentUsage: currentUsage,
                   Limit: dimension.Limit
               }
           END IF
       END IF
   END FOR

5. ========== 检查配额升级 ==========

   IF plan.AllowQuotaExceed AND usage.CanUpgrade:
       // 允许超额,但需要付费
       extraCost := CalculateExtraCost(input.Quantity, plan)

       RETURN {
           Allowed: true,
           Reason: f"Quota exceeded, extra cost: {extraCost}",
           RemainingQuota: 0,
           CurrentUsage: usage.TotalCount,
           Limit: -1  // -1 表示无限制,但需要付费
       }
   END IF

6. ========== 配额检查通过 ==========

   RETURN {
       Allowed: true,
       Reason: "Quota check passed",
       RemainingQuota: quotas.Limit - usage.TotalCount,
       CurrentUsage: usage.TotalCount,
       Limit: quotas.Limit
   }

7. ========== 记录配额使用 ==========

   IF output.Allowed:
       RecordQuotaUsage(input.TenantID, input.ResourceType, input.Quantity)
   END IF
```

#### 7.1.3 配额维度

**配额维度类型**:

1. **时间维度**: 按时间窗口限制
   - 每小时 (hourly)
   - 每天 (daily)
   - 每月 (monthly)

2. **总量维度**: 总数限制
   - 总数量 (total)

3. **并发维度**: 并发数限制
   - 并发数 (concurrent)

**配额检查顺序**: 按优先级检查,任何一维超限即拒绝

---

## 数据脱敏算法

### 8.1 敏感数据脱敏算法

**算法名称**: `SensitiveDataMasker`

**功能说明**: 根据用户角色和数据类型脱敏敏感数据

**位置**: `backend/pkg/masking/masker.go`

#### 8.1.1 输入输出

**输入**:
```go
type DataMaskInput struct {
    Data        string
    DataType    string
    UserRole    string
    Context     map[string]interface{}
}
```

**输出**:
```go
type DataMaskOutput struct {
    MaskedData  string
    OriginalData string
    MaskLevel   string  // none, partial, full
}
```

#### 8.1.2 算法步骤

```text
算法 SensitiveDataMasker:
输入: input (DataMaskInput)
输出: output (DataMaskOutput)

1. ========== 检查是否需要脱敏 ==========

   // 1.1 管理员不脱敏
   IF input.UserRole IN ["admin", "system_admin"]:
       RETURN {
           MaskedData: input.Data,
           OriginalData: input.Data,
           MaskLevel: "none"
       }

   // 1.2 数据所有者不脱敏
   IF IsDataOwner(input.Context, input.Data):
       RETURN {
           MaskedData: input.Data,
           OriginalData: input.Data,
           MaskLevel: "none"
       }

   // 1.3 检查是否为敏感数据
   IF !IsSensitiveData(input.DataType):
       RETURN {
           MaskedData: input.Data,
           OriginalData: input.Data,
           MaskLevel: "none"
       }

2. ========== 根据数据类型脱敏 ==========

   maskedData := ""

   SWITCH input.DataType:

       CASE "phone":
           // 手机号: 138****1234
           maskedData = input.Data[0:3] + "****" + input.Data[7:11]

       CASE "email":
           // 邮箱: u***@example.com
           parts = Split(input.Data, "@")
           username = parts[0]
           domain = parts[1]
           maskedUsername = username[0] + "***"
           maskedData = maskedUsername + "@" + domain

       CASE "id_card":
           // 身份证: 110101********1234
           maskedData = input.Data[0:6] + "********" + input.Data[14:18]

       CASE "bank_card":
           // 银行卡: 6222********1234
           maskedData = input.Data[0:4] + "****" + "*" * 4 + input.Data[12:16]

       CASE "name":
           // 姓名: 王**
           maskedData = input.Data[0] + "**"

       CASE "address":
           // 地址: 北京市朝阳区****
           maskedData = input.Data[0:6] + "****"

       DEFAULT:
           // 未知类型,完全脱敏
           maskedData = "****"

   END SWITCH

3. ========== 根据用户角色调整脱敏级别 ==========

   // 审计人员可以看到部分信息
   IF input.UserRole == "auditor":

       SWITCH input.DataType:
           CASE "phone":
               maskedData = input.Data[0:3] + "****" + input.Data[9:11]  // 显示最后2位
           CASE "email":
               maskedData = input.Data[0:2] + "***@" + parts[1]  // 显示前2位
           CASE "id_card":
               maskedData = input.Data[0:6] + "****" + input.Data[16:18]  // 显示最后2位
       END SWITCH
   END IF

4. ========== 返回结果 ==========

   RETURN {
       MaskedData: maskedData,
       OriginalData: input.Data,
       MaskLevel: GetMaskLevel(input.UserRole)
   }
```

#### 8.1.3 脱敏规则表

| 数据类型 | 管理员 | 所有者 | 普通用户 | 审计人员 |
|---------|--------|--------|---------|---------|
| **phone** | 不脱敏 | 不脱敏 | 138\*\*\*\*1234 | 138\*\*\*\*34 |
| **email** | 不脱敏 | 不脱敏 | u\*\*\*@example.com | u\*\*\*@example.com |
| **id_card** | 不脱敏 | 不脱敏 | 110101\*\*\*\*\*\*\*1234 | 110101\*\*\*\*\*\*\*34 |
| **bank_card** | 不脱敏 | 不脱敏 | 6222\*\*\*\*\*\*\*1234 | 6222\*\*\*\*\*\*\*1234 |
| **name** | 不脱敏 | 不脱敏 | 王\*\* | 王\*\* |
| **address** | 不脱敏 | 不脱敏 | 北京市朝阳区\*\*\*\* | 北京市朝阳区\*\*\*\* |

---

## 附录

### A. 算法复杂度汇总

| 算法名称 | 时间复杂度 | 空间复杂度 | 备注 |
|---------|-----------|-----------|------|
| **意图识别** | O(n*m + s*m) | O(m + s) | n=输入长度, m=特征维度, s=样本数 |
| **路由决策** | O(n + n log n) | O(n) | n=候选服务数 |
| **RBAC权限** | O(r) | O(r) | r=规则数 |
| **租户隔离** | O(1) | O(1) | 添加WHERE条件 |
| **工作流执行** | O(V + E) | O(V) | V=节点数, E=边数 |
| **知识检索** | O(k*m + k*m) | O(k) | k=top_k, m=特征维度 |
| **配额检查** | O(d) | O(1) | d=维度数 |
| **数据脱敏** | O(1) | O(1) | 字符串操作 |

### B. 相关文档

- [ZKER 统一错误码定义规范](./ZKER-统一错误码定义规范.md)
- [ZKER 技术组件清单与使用指南](./ZKER-技术组件清单与使用指南(完整版).md)
- [可执行性文档缺失分析与补充方案](./可执行性文档缺失分析与补充方案_v1.0.md)
- [OpenAPI 3.0 规范](./openapi/zker-api-v1-core-modules.yaml)

---

**文档状态**: ✅ 初版完成
**维护责任**: ZKER 后端架构团队
**审核状态**: 待审核
