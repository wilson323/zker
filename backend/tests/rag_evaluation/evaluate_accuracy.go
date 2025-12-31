/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ragevaluation

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// Evaluator RAG准确率评估器
type Evaluator struct {
	searchEngine SearchEngine
	testCases    []*TestCase
	metrics      *EvaluationMetrics
}

// SearchEngine 搜索引擎接口
type SearchEngine interface {
	Search(ctx context.Context, query string, topK int) ([]*schema.Document, error)
}

// TestCase 测试用例
type TestCase struct {
	// 查询内容
	Query string
	// 相关文档ID列表(人工标注)
	RelevantDocIDs []string
	// 相关文档ID到相关性分数的映射(可选)
	RelevanceScores map[string]float64
	// 难度等级: easy, medium, hard
	Difficulty string
	// 分类: keyword, semantic, hybrid
	Category string
}

// EvaluationMetrics 评估指标
type EvaluationMetrics struct {
	// Precision@K:准确率
	PrecisionAtK map[int]float64
	// Recall@K:召回率
	RecallAtK map[int]float64
	// MAP:平均精度均值
	MeanAveragePrecision float64
	// MRR:平均倒数排名
	MeanReciprocalRank float64
	// NDCG:归一化折损累积增益
	NDCG map[int]float64
	// 总体准确率
	OverallAccuracy float64
	// 目标准确率
	TargetAccuracy float64
	// 是否达标
	Achieved bool
}

// EvaluationReport 评估报告
type EvaluationReport struct {
	Timestamp      time.Time
	TestCount      int
	Metrics        *EvaluationMetrics
	PerQueryScores map[string]*QueryScore
	Summary        string
}

// QueryScore 单个查询的得分
type QueryScore struct {
	Query              string
	Precision          float64
	Recall             float64
	AveragePrecision   float64
	ReciprocalRank     float64
	NDCG               float64
	RelevantFound      int
	RelevantTotal      int
	RetrievedRelevant  []string
	RetrievedNonRelevant []string
	MissingRelevant    []string
}

// NewEvaluator 创建评估器
func NewEvaluator(searchEngine SearchEngine, testCases []*TestCase, targetAccuracy float64) *Evaluator {
	return &Evaluator{
		searchEngine: searchEngine,
		testCases:    testCases,
		metrics: &EvaluationMetrics{
			PrecisionAtK:    make(map[int]float64),
			RecallAtK:       make(map[int]float64),
			NDCG:            make(map[int]float64),
			TargetAccuracy:  targetAccuracy,
		},
	}
}

// Evaluate 执行评估
func (e *Evaluator) Evaluate(ctx context.Context, topKs []int) (*EvaluationReport, error) {
	if len(e.testCases) == 0 {
		return nil, fmt.Errorf("no test cases provided")
	}

	logs.CtxInfof(ctx, "[RAGEvaluation] starting evaluation with %d test cases", len(e.testCases))

	report := &EvaluationReport{
		Timestamp:      time.Now(),
		TestCount:      len(e.testCases),
		PerQueryScores: make(map[string]*QueryScore),
	}

	// 初始化指标
	for _, k := range topKs {
		e.metrics.PrecisionAtK[k] = 0.0
		e.metrics.RecallAtK[k] = 0.0
		e.metrics.NDCG[k] = 0.0
	}

	sumAP := 0.0
	sumMRR := 0.0

	// 逐个评估测试用例
	for i, tc := range e.testCases {
		logs.CtxInfof(ctx, "[RAGEvaluation] evaluating test case %d/%d: %s", i+1, len(e.testCases), tc.Query)

		// 执行搜索
		docs, err := e.searchEngine.Search(ctx, tc.Query, 20) // 搜索20条结果用于评估
		if err != nil {
			logs.CtxErrorf(ctx, "[RAGEvaluation] search failed for query '%s': %v", tc.Query, err)
			continue
		}

		// 计算该查询的指标
		queryScore := e.evaluateQuery(tc, docs)

		// 累积指标
		for _, k := range topKs {
			e.metrics.PrecisionAtK[k] += queryScore.Precision
			e.metrics.RecallAtK[k] += queryScore.Recall
			e.metrics.NDCG[k] += queryScore.NDCG
		}

		sumAP += queryScore.AveragePrecision
		sumMRR += queryScore.ReciprocalRank

		// 记录每个查询的得分
		report.PerQueryScores[tc.Query] = queryScore
	}

	// 计算平均指标
	for _, k := range topKs {
		e.metrics.PrecisionAtK[k] /= float64(len(e.testCases))
		e.metrics.RecallAtK[k] /= float64(len(e.testCases))
		e.metrics.NDCG[k] /= float64(len(e.testCases))
	}

	e.metrics.MeanAveragePrecision = sumAP / float64(len(e.testCases))
	e.metrics.MeanReciprocalRank = sumMRR / float64(len(e.testCases))

	// 计算总体准确率 (使用Precision@10)
	e.metrics.OverallAccuracy = e.metrics.PrecisionAtK[10]
	e.metrics.Achieved = e.metrics.OverallAccuracy >= e.metrics.TargetAccuracy

	report.Metrics = e.metrics
	report.Summary = e.generateSummary()

	logs.CtxInfof(ctx, "[RAGEvaluation] evaluation completed:\n%s", report.Summary)

	return report, nil
}

// evaluateQuery 评估单个查询
func (e *Evaluator) evaluateQuery(tc *TestCase, docs []*schema.Document) *QueryScore {
	score := &QueryScore{
		Query:             tc.Query,
		RelevantTotal:     len(tc.RelevantDocIDs),
		RetrievedRelevant: make([]string, 0),
		RetrievedNonRelevant: make([]string, 0),
		MissingRelevant:   make([]string, 0),
	}

	relevantSet := make(map[string]bool)
	for _, id := range tc.RelevantDocIDs {
		relevantSet[id] = true
	}

	// 计算各指标
	retrievedRelevantCount := 0
	var dcg, idcg float64
	averagePrecision := 0.0
	relevantFound := 0
	firstRelevantRank := -1

	for rank, doc := range docs {
		if rank >= 20 {
			break
		}

		docID := doc.ID
		isRelevant := relevantSet[docID]

		// 计算DCG
		if isRelevant {
			relevanceScore := 1.0
			if tc.RelevanceScores != nil {
				if score, ok := tc.RelevanceScores[docID]; ok {
					relevanceScore = score
				}
			}
			dcg += relevanceScore / math.Log2(float64(rank+2))

			retrievedRelevantCount++
			if firstRelevantRank == -1 {
				firstRelevantRank = rank + 1
			}
			relevantFound++
			score.RetrievedRelevant = append(score.RetrievedRelevant, docID)
		} else {
			score.RetrievedNonRelevant = append(score.RetrievedNonRelevant, docID)
		}

		// 计算Average Precision
		if isRelevant {
			relevantFound++
			precisionAtK := float64(relevantFound) / float64(rank+1)
			averagePrecision += precisionAtK
		}
	}

	// 计算IDCG (理想DCG)
	for i, id := range tc.RelevantDocIDs {
		relevanceScore := 1.0
		if tc.RelevanceScores != nil {
			if score, ok := tc.RelevanceScores[id]; ok {
				relevanceScore = score
			}
		}
		idcg += relevanceScore / math.Log2(float64(i+2))
	}

	// 计算NDCG
	ndcg := 0.0
	if idcg > 0 {
		ndcg = dcg / idcg
	}

	// 计算Precision@10和Recall@10
	precisionAt10 := float64(retrievedRelevantCount) / math.Min(10.0, float64(len(docs)))
	recallAt10 := float64(retrievedRelevantCount) / float64(len(tc.RelevantDocIDs))

	// 归一化Average Precision
	if len(tc.RelevantDocIDs) > 0 {
		averagePrecision /= float64(len(tc.RelevantDocIDs))
	}

	// 计算MRR
	mrr := 0.0
	if firstRelevantRank > 0 {
		mrr = 1.0 / float64(firstRelevantRank)
	}

	// 找出缺失的相关文档
	for _, id := range tc.RelevantDocIDs {
		found := false
		for _, doc := range docs {
			if doc.ID == id {
				found = true
				break
			}
		}
		if !found {
			score.MissingRelevant = append(score.MissingRelevant, id)
		}
	}

	// 填充得分
	score.Precision = precisionAt10
	score.Recall = recallAt10
	score.AveragePrecision = averagePrecision
	score.ReciprocalRank = mrr
	score.NDCG = ndcg
	score.RelevantFound = retrievedRelevantCount

	return score
}

// generateSummary 生成评估摘要
func (e *Evaluator) generateSummary() string {
	summary := fmt.Sprintf(`
========================================
RAG检索准确率评估报告
========================================
目标准确率: %.2f%%
实际准确率: %.2f%%
是否达标: %v

关键指标:
- MAP (平均精度均值): %.4f
- MRR (平均倒数排名): %.4f
- NDCG@10: %.4f

Precision@K:
`,
		e.metrics.TargetAccuracy*100,
		e.metrics.OverallAccuracy*100,
		e.metrics.Achieved,
		e.metrics.MeanAveragePrecision,
		e.metrics.MeanReciprocalRank,
		e.metrics.NDCG[10])

	for k, precision := range e.metrics.PrecisionAtK {
		summary += fmt.Sprintf("  - Precision@%d: %.4f\n", k, precision)
	}

	summary += "\nRecall@K:\n"
	for k, recall := range e.metrics.RecallAtK {
		summary += fmt.Sprintf("  - Recall@%d: %.4f\n", k, recall)
	}

	summary += "========================================\n"

	return summary
}

// SaveReport 保存评估报告到文件
func (e *Evaluator) SaveReport(report *EvaluationReport, outputDir string) error {
	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory failed: %w", err)
	}

	// 保存JSON报告
	jsonPath := filepath.Join(outputDir, fmt.Sprintf("rag_evaluation_%d.json", report.Timestamp.Unix()))
	if err := e.saveJSONReport(report, jsonPath); err != nil {
		return fmt.Errorf("save JSON report failed: %w", err)
	}

	// 保存CSV报告
	csvPath := filepath.Join(outputDir, fmt.Sprintf("rag_evaluation_%d.csv", report.Timestamp.Unix()))
	if err := e.saveCSVReport(report, csvPath); err != nil {
		return fmt.Errorf("save CSV report failed: %w", err)
	}

	// 保存文本摘要
	txtPath := filepath.Join(outputDir, fmt.Sprintf("rag_evaluation_summary_%d.txt", report.Timestamp.Unix()))
	if err := e.saveTextReport(report, txtPath); err != nil {
		return fmt.Errorf("save text report failed: %w", err)
	}

	logs.CtxInfof(context.Background(), "[RAGEvaluation] saved evaluation reports to: %s", outputDir)

	return nil
}

// saveJSONReport 保存JSON格式报告
func (e *Evaluator) saveJSONReport(report *EvaluationReport, path string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// saveCSVReport 保存CSV格式报告
func (e *Evaluator) saveCSVReport(report *EvaluationReport, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	headers := []string{
		"Query", "Precision", "Recall", "AP", "MRR", "NDCG",
		"RelevantFound", "RelevantTotal", "RetrievedRelevant",
		"RetrievedNonRelevant", "MissingRelevant",
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// 写入每行数据
	for query, score := range report.PerQueryScores {
		row := []string{
			query,
			fmt.Sprintf("%.4f", score.Precision),
			fmt.Sprintf("%.4f", score.Recall),
			fmt.Sprintf("%.4f", score.AveragePrecision),
			fmt.Sprintf("%.4f", score.ReciprocalRank),
			fmt.Sprintf("%.4f", score.NDCG),
			fmt.Sprintf("%d", score.RelevantFound),
			fmt.Sprintf("%d", score.RelevantTotal),
			fmt.Sprintf("%v", score.RetrievedRelevant),
			fmt.Sprintf("%v", score.RetrievedNonRelevant),
			fmt.Sprintf("%v", score.MissingRelevant),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// saveTextReport 保存文本格式报告
func (e *Evaluator) saveTextReport(report *EvaluationReport, path string) error {
	// 添加详细分析
	detailedSummary := report.Summary + "\n\n详细分析:\n"

	// 找出表现最差的5个查询
	type queryError struct {
		query string
		score float64
	}

	errors := make([]queryError, 0)
	for query, score := range report.PerQueryScores {
		if score.Precision < 0.5 { // 准确率低于50%
			errors = append(errors, queryError{
				query: query,
				score: score.Precision,
			})
		}
	}

	if len(errors) > 0 {
		detailedSummary += fmt.Sprintf("\n表现不佳的查询 (%d 个):\n", len(errors))
		for i, err := range errors {
			if i >= 5 {
				break
			}
			score := report.PerQueryScores[err.query]
			detailedSummary += fmt.Sprintf("  %d. %s\n", i+1, err.query)
			detailedSummary += fmt.Sprintf("     准确率: %.2f%%, 召回: %.2f%%, 缺失相关文档: %d\n",
				score.Precision*100, score.Recall*100, len(score.MissingRelevant))
		}
	}

	return os.WriteFile(path, []byte(detailedSummary), 0644)
}

// LoadTestCases 从文件加载测试用例
func LoadTestCases(jsonPath string) ([]*TestCase, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("read test cases file failed: %w", err)
	}

	var testCases []*TestCase
	if err := json.Unmarshal(data, &testCases); err != nil {
		return nil, fmt.Errorf("parse test cases failed: %w", err)
	}

	return testCases, nil
}

// SaveTestCases 保存测试用例到文件
func SaveTestCases(testCases []*TestCase, jsonPath string) error {
	data, err := json.MarshalIndent(testCases, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal test cases failed: %w", err)
	}

	return os.WriteFile(jsonPath, data, 0644)
}

// GenerateDefaultTestCases 生成默认测试用例
func GenerateDefaultTestCases(knowledgeBase string) []*TestCase {
	// 根据不同的知识库生成测试用例
	testCases := make([]*TestCase, 0)

	switch knowledgeBase {
	case "platform":
		// 平台使用相关
		testCases = append(testCases, &TestCase{
			Query:           "如何创建AI智能体",
			RelevantDocIDs:  []string{"doc1", "doc2"},
			Difficulty:      "easy",
			Category:        "keyword",
		})
		testCases = append(testCases, &TestCase{
			Query:           "ZKER平台的优势是什么",
			RelevantDocIDs:  []string{"doc3", "doc4", "doc5"},
			Difficulty:      "medium",
			Category:        "semantic",
		})
	case "tenant":
		// 租户管理相关
		testCases = append(testCases, &TestCase{
			Query:           "如何创建租户",
			RelevantDocIDs:  []string{"doc10", "doc11"},
			Difficulty:      "easy",
			Category:        "keyword",
		})
	case "permission":
		// 权限管理相关
		testCases = append(testCases, &TestCase{
			Query:           "角色和权限的关系",
			RelevantDocIDs:  []string{"doc20", "doc21", "doc22"},
			Difficulty:      "hard",
			Category:        "hybrid",
		})
	}

	return testCases
}
