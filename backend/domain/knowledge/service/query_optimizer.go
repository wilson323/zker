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

package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// QueryOptimizer 查询优化器
// 通过LLM驱动的查询扩展、同义词替换、拼写纠正等提升检索准确率
type QueryOptimizer struct {
	// LLM客户端(可选,用于高级查询扩展)
	llmClient LLMClient
	// 同义词词典
	synonymDict *SynonymDictionary
	// 拼写纠正器
	spellChecker *SpellChecker
	// 查询缓存
	queryCache *QueryCache
}

// LLMClient LLM客户端接口
type LLMClient interface {
	GenerateQueryVariations(ctx context.Context, query string) ([]string, error)
}

// SynonymDictionary 同义词词典
type SynonymDictionary struct {
	mu     sync.RWMutex
	synonyms map[string][]string
}

// SpellChecker 拼写纠正器
type SpellChecker struct {
	mu       sync.RWMutex
	dict     map[string]bool
	edits    map[string][]string // 编辑距离映射
}

// QueryCache 查询缓存
type QueryCache struct {
	mu    sync.RWMutex
	cache map[string]*CachedQuery
}

// CachedQuery 缓存的查询
type CachedQuery struct {
	Original      string
	Optimized     string
	Expanded      []string
	LastUsed      int64
	UseCount      int
}

// NewQueryOptimizer 创建查询优化器
func NewQueryOptimizer() *QueryOptimizer {
	return &QueryOptimizer{
		synonymDict: NewSynonymDictionary(),
		spellChecker: NewSpellChecker(),
		queryCache:   NewQueryCache(),
	}
}

// WithLLMClient 设置LLM客户端
func (o *QueryOptimizer) WithLLMClient(client LLMClient) *QueryOptimizer {
	o.llmClient = client
	return o
}

// ExpandQuery 扩展查询
// 返回: (优化后的查询, 扩展的查询列表, 错误)
func (o *QueryOptimizer) ExpandQuery(ctx context.Context, query string) (string, []string, error) {
	if query == "" {
		return "", nil, fmt.Errorf("query cannot be empty")
	}

	// 1. 检查缓存
	if cached := o.queryCache.Get(query); cached != nil {
		logs.CtxInfof(ctx, "[QueryOptimizer] cache hit for query: %s", query)
		return cached.Optimized, cached.Expanded, nil
	}

	// 2. 预处理:小写化、去除多余空格
	optimized := o.preprocessQuery(query)

	// 3. 拼写纠正
	corrected, err := o.spellChecker.Correct(optimized)
	if err != nil {
		logs.CtxWarnf(ctx, "[QueryOptimizer] spell correction failed: %v", err)
		corrected = optimized
	}

	// 4. 同义词扩展
	expanded := o.expandWithSynonyms(corrected)

	// 5. LLM查询扩展(如果配置了LLM)
	if o.llmClient != nil {
		llmExpanded, err := o.expandWithLLM(ctx, corrected)
		if err != nil {
			logs.CtxWarnf(ctx, "[QueryOptimizer] LLM expansion failed: %v", err)
		} else {
			expanded = append(expanded, llmExpanded...)
		}
	}

	// 6. 去重
	expanded = o.deduplicateQueries(expanded)

	// 7. 缓存结果
	o.queryCache.Set(query, corrected, expanded)

	logs.CtxInfof(ctx, "[QueryOptimizer] expanded query '%s' to %d variations", query, len(expanded)+1)

	return corrected, expanded, nil
}

// preprocessQuery 预处理查询
func (o *QueryOptimizer) preprocessQuery(query string) string {
	// 小写化
	query = strings.ToLower(query)

	// 去除多余空格
	query = strings.Join(strings.Fields(query), " ")

	// 去除特殊字符(保留中文、英文、数字、空格)
	// query = strings.TrimSpace(regexp.MustCompile(`[^\p{L}\p{N}\s]+`).ReplaceAllString(query, " "))

	return query
}

// expandWithSynonyms 同义词扩展
func (o *QueryOptimizer) expandWithSynonyms(query string) []string {
	words := strings.Fields(query)
	expanded := make([]string, 0)

	// 对每个词查找同义词
	for i, word := range words {
		synonyms := o.synonymDict.GetSynonyms(word)
		if len(synonyms) > 0 {
			// 为每个同义词生成一个新的查询变体
			for _, synonym := range synonyms {
				newWords := make([]string, len(words))
				copy(newWords, words)
				newWords[i] = synonym
				expanded = append(expanded, strings.Join(newWords, " "))
			}
		}
	}

	return expanded
}

// expandWithLLM 使用LLM扩展查询
func (o *QueryOptimizer) expandWithLLM(ctx context.Context, query string) ([]string, error) {
	if o.llmClient == nil {
		return nil, nil
	}

	variations, err := o.llmClient.GenerateQueryVariations(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("LLM query expansion failed: %w", err)
	}

	return variations, nil
}

// deduplicateQueries 去重查询
func (o *QueryOptimizer) deduplicateQueries(queries []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(queries))

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query != "" && !seen[query] {
			seen[query] = true
			result = append(result, query)
		}
	}

	return result
}

// NewSynonymDictionary 创建同义词词典
func NewSynonymDictionary() *SynonymDictionary {
	dict := &SynonymDictionary{
		synonyms: make(map[string][]string),
	}

	// 加载默认同义词
	dict.loadDefaultSynonyms()

	return dict
}

// loadDefaultSynonyms 加载默认同义词
func (d *SynonymDictionary) loadDefaultSynonyms() {
	// IT领域同义词
	d.Add("搜索", []string{"检索", "查找", "查询", "搜寻"})
	d.Add("文档", []string{"资料", "文件", "文献", "文本"})
	d.Add("数据库", []string{"db", "database", "数据存储"})
	d.Add("api", []string{"接口", "interface", "应用程序接口"})
	d.Add("bug", []string{"缺陷", "错误", "故障", "问题"})
	d.Add("优化", []string{"改进", "提升", "改善", "增强"})
	d.Add("配置", []string{"设置", "设定", "参数配置"})

	// 业务领域同义词(根据实际业务调整)
	d.Add("租户", []string{"tenant", "账户", "账号"})
	d.Add("权限", []string{"permission", "授权", "许可"})
	d.Add("角色", []string{"role", "岗位", "职位"})
}

// Add 添加同义词
func (d *SynonymDictionary) Add(word string, synonyms []string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	word = strings.ToLower(word)
	for i, s := range synonyms {
		synonyms[i] = strings.ToLower(s)
	}

	d.synonyms[word] = synonyms

	// 双向绑定:让同义词也能反向查找
	for _, synonym := range synonyms {
		if _, exists := d.synonyms[synonym]; !exists {
			d.synonyms[synonym] = []string{word}
		}
	}
}

// GetSynonyms 获取同义词
func (d *SynonymDictionary) GetSynonyms(word string) []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	word = strings.ToLower(word)
	if synonyms, ok := d.synonyms[word]; ok {
		return synonyms
	}
	return nil
}

// NewSpellChecker 创建拼写纠正器
func NewSpellChecker() *SpellChecker {
	return &SpellChecker{
		dict:  make(map[string]bool),
		edits: make(map[string][]string),
	}
}

// LoadDictionary 加载词典
func (s *SpellChecker) LoadDictionary(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, word := range words {
		word = strings.ToLower(word)
		s.dict[word] = true
	}
}

// Correct 纠正拼写
func (s *SpellChecker) Correct(query string) (string, error) {
	words := strings.Fields(query)
	corrected := make([]string, len(words))

	for i, word := range words {
		if s.dict[word] {
			// 词语在词典中,无需纠正
			corrected[i] = word
		} else {
			// 查找最接近的词
			corrected[i] = s.findClosestWord(word)
		}
	}

	return strings.Join(corrected, " "), nil
}

// findClosestWord 查找最接近的词(基于编辑距离)
func (s *SpellChecker) findClosestWord(word string) string {
	// 如果词典为空,返回原词
	if len(s.dict) == 0 {
		return word
	}

	word = strings.ToLower(word)
	minDist := int(^uint(0) >> 1) // 最大int
	closest := word

	for dictWord := range s.dict {
		dist := editDistance(word, dictWord)
		if dist < minDist {
			minDist = dist
			closest = dictWord
		}
	}

	// 如果编辑距离太大(超过词长的一半),可能是新词,不纠正
	if minDist > len(word)/2 {
		return word
	}

	return closest
}

// editDistance 计算编辑距离(Levenshtein距离)
func editDistance(a, b string) int {
	lenA := len(a)
	lenB := len(b)

	// 创建DP表
	dp := make([][]int, lenA+1)
	for i := range dp {
		dp[i] = make([]int, lenB+1)
	}

	// 初始化
	for i := 0; i <= lenA; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= lenB; j++ {
		dp[0][j] = j
	}

	// 动态规划
	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min(
					dp[i-1][j]+1,   // 删除
					dp[i][j-1]+1,   // 插入
					dp[i-1][j-1]+1, // 替换
				)
			}
		}
	}

	return dp[lenA][lenB]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// NewQueryCache 创建查询缓存
func NewQueryCache() *QueryCache {
	return &QueryCache{
		cache: make(map[string]*CachedQuery),
	}
}

// Get 获取缓存的查询
func (c *QueryCache) Get(query string) *CachedQuery {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query = strings.ToLower(strings.TrimSpace(query))
	if cached, ok := c.cache[query]; ok {
		cached.UseCount++
		return cached
	}
	return nil
}

// Set 设置缓存的查询
func (c *QueryCache) Set(original, optimized string, expanded []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	original = strings.ToLower(strings.TrimSpace(original))
	c.cache[original] = &CachedQuery{
		Original:  original,
		Optimized: optimized,
		Expanded:  expanded,
		LastUsed:  0, // TODO: 设置为当前时间戳
		UseCount:  1,
	}
}

// Clear 清空缓存
func (c *QueryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*CachedQuery)
}

// Size 获取缓存大小
func (c *QueryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.cache)
}
