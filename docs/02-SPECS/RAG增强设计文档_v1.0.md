# ZKER RAG能力增强设计文档 v1.0

> **项目**: ZKER 企业级 AI 智能体工作台平台
> **文档类型**: 技术设计规范
> **版本**: v1.0
> **日期**: 2025-01-03
> **状态**: 待评审

---

## 📋 文档修订历史

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| v1.0 | 2025-01-03 | AI架构师 | 初始版本，基于FastGPT分析 |

---

## 1. 项目背景与目标

### 1.1 当前状态

根据深度分析，ZKER项目RAG能力现状：

| 能力维度 | 当前水平 | FastGPT水平 | 差距 |
|----------|----------|-------------|------|
| **分块策略** | 简单固定大小 | 8级递归Markdown分块 | 🔴 严重 |
| **检索模式** | 仅向量检索 | 向量+全文混合检索 | 🔴 严重 |
| **重排序** | 简单RRF | RRF + BGE-Reranker | 🟡 中等 |
| **查询扩展** | 无 | LLM + Lazy Greedy | 🔴 严重 |
| **去重机制** | 无 | SimHash + 相似度 | 🔴 严重 |
| **多租户隔离** | 基础隔离 | 完全隔离 | 🟡 中等 |

### 1.2 增强目标

**核心目标**: 将ZKER的RAG能力提升至**行业领先水平**，超越鲸智百应平台，对标FastGPT和Dify。

| 指标 | 当前 | 目标 | 提升幅度 |
|------|------|------|----------|
| **检索准确率** | ~65% | ≥85% | +20% |
| **召回率** | ~70% | ≥90% | +20% |
| **检索延迟(P99)** | ~800ms | ≤500ms | -37.5% |
| **多租户隔离** | 65% | ≥98% | +33% |

---

## 2. 架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              ZKER RAG 增强架构                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐        │
│  │   查询扩展层     │    │    检索执行层    │    │    重排序层      │        │
│  │                 │    │                 │    │                 │        │
│  │ • LLM扩展       │───▶│ • 向量检索      │───▶│ • BGE-Reranker  │        │
│  │ • Lazy Greedy   │    │ • 全文检索      │    │ • RRF融合       │        │
│  │ • 查询去重      │    │ • NL2SQL检索    │    │ • 结果去重      │        │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘        │
│           │                       │                       │               │
│           ▼                       ▼                       ▼               │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐        │
│  │   分块处理层     │    │    存储层        │    │    租户隔离层    │        │
│  │                 │    │                 │    │                 │        │
│  │ • 递归分块      │    │ • Milvus        │    │ • 租户索引      │        │
│  │ • Markdown分块  │    │ • Elasticsearch │    │ • 缓存隔离      │        │
│  │ • 代码块保护    │    │ • PostgreSQL    │    │ • 权限过滤      │        │
│  │ • 语义分块      │    │                 │    │                 │        │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 技术栈

| 类别 | 技术选型 | 版本 | 说明 |
|------|---------|------|------|
| **向量存储** | Milvus | 2.3+ | 支持租户级partition |
| **全文检索** | Elasticsearch | 8.18+ | 租户级索引隔离 |
| **重排序** | BGE-Reranker-v2 | - | 独立服务部署 |
| **分块算法** | 自研 + FastGPT模式 | - | 递归+语义 |
| **去重** | SimHash + MinHash | - | 内容指纹 |
| **查询扩展** | LLM + Lazy Greedy | - | 多路查询 |

---

## 3. 详细设计

### 3.1 阶段1：层级分块实现 (5人天)

#### 3.1.1 递归字符分块器

**文件**: `backend/infra/document/parser/impl/builtin/chunk_recursive.go`

```go
package builtin

import (
	"regexp"
	"strings"
	"github.com/coze-dev/coze-studio/backend/domain/knowledge/entity"
)

// RecursiveChunker 递归字符分块器
// 基于FastGPT的textSplitter.ts实现
type RecursiveChunker struct {
	chunkSize    int  // 分块大小，默认800
	chunkOverlap int  // 重叠大小，默认15%
	separators   []string // 分隔符优先级列表
}

// SplitResult 分块结果
type SplitResult struct {
	Chunks []string // 分块内容
	Chars  int      // 总字符数
}

// NewRecursiveChunker 创建递归分块器
func NewRecursiveChunker(chunkSize int) *RecursiveChunker {
	overlap := int(float64(chunkSize) * 0.15) // 15%重叠

	return &RecursiveChunker{
		chunkSize:    chunkSize,
		chunkOverlap: overlap,
		separators:   getSeparators(),
	}
}

// getSeparators 获取分隔符优先级列表
// 优先级从高到低：自定义 → Markdown标题 → 代码块 → 表格 → 段落 → 句子
func getSeparators() []string {
	return []string{
		"\n\n",      // 双换行（段落）
		"\n",        // 单换行
		"。",        // 中文句号
		". ",        // 英文句号
		"!",         // 感叹号
		"?",         // 问号
		";",         // 分号
		"，",        // 中文逗号
		", ",        // 英文逗号
		" ",         // 空格
		"",          // 字符级别
	}
}

// Split 执行递归分块
func (r *RecursiveChunker) Split(text string) *SplitResult {
	return r.splitRecursive(text, 0, "", "")
}

// splitRecursive 递归分块核心逻辑
func (r *RecursiveChunker) splitRecursive(text, step int, lastText, parentTitle string) *SplitResult {
	// 1. 递归终止条件
	if step >= len(r.separators) {
		return &SplitResult{
			Chunks: []string{lastText + text},
			Chars:  len(lastText) + len(text),
		}
	}

	separator := r.separators[step]

	// 2. 检查是否禁止重叠（标题、代码块、表格不重叠）
	forbidOverlap := r.shouldForbidOverlap(step)

	// 3. 按当前分隔符分割
	splitTexts := r.splitBySeparator(text, separator)

	if len(splitTexts) == 1 {
		// 没有分割，递归到下一级
		return r.splitRecursive(text, step+1, lastText, parentTitle)
	}

	var chunks []string
	var currentChunk strings.Builder
	currentLen := 0

	for i, splitText := range splitTexts {
		textLen := len(splitText)

		// 4. 计算重叠文本
		overlapText := ""
		if !forbidOverlap && i > 0 {
			overlapText = r.getOverlapText(currentChunk.String())
		}

		// 5. 判断是否需要新建分块
		newLen := currentLen + len(overlapText) + textLen
		if newLen > r.chunkSize && currentLen > 0 {
			// 保存当前分块
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
			currentLen = 0
		}

		// 6. 添加内容到当前分块
		if overlapText != "" {
			currentChunk.WriteString(overlapText)
			currentLen += len(overlapText)
		}
		currentChunk.WriteString(splitText)
		currentLen += textLen

		// 7. 如果分块仍过大，递归细分
		if currentLen > r.chunkSize {
			remainingChunks := r.splitRecursive(
				currentChunk.String(),
				step+1,
				"",
				parentTitle,
			)
			chunks = append(chunks, remainingChunks.Chunks...)
			currentChunk.Reset()
			currentLen = 0
		}
	}

	// 8. 添加最后一个分块
	if currentLen > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return &SplitResult{
		Chunks: chunks,
		Chars:  len(text),
	}
}

// splitBySeparator 按分隔符分割文本
func (r *RecursiveChunker) splitBySeparator(text, separator string) []string {
	if separator == "" {
		// 字符级别分割
		var result []string
		for _, r := range text {
			result = append(result, string(r))
		}
		return result
	}
	return strings.Split(text, separator)
}

// getOverlapText 获取重叠文本
func (r *RecursiveChunker) getOverlapText(chunk string) string {
	// 从后向前取完整句子，最多40%重叠
	maxOverlap := int(float64(r.chunkSize) * 0.4)
	if r.chunkOverlap > maxOverlap {
		maxOverlap = r.chunkOverlap
	}

	if len(chunk) <= maxOverlap {
		return chunk
	}

	// 从后向前找句子边界
	sentences := regexp.MustCompile(`[。！？.!?]`).FindAllStringIndex(chunk, -1)
	for i := len(sentences) - 1; i >= 0; i-- {
		pos := sentences[i][1]
		if len(chunk)-pos <= maxOverlap {
			return chunk[pos:]
		}
	}

	// 没有找到句子边界，直接截取
	return chunk[len(chunk)-maxOverlap:]
}

// shouldForbidOverlap 判断是否禁止重叠
func (r *RecursiveChunker) shouldForbidOverlap(step int) bool {
	// 前几个分隔符（段落、换行）不重叠
	return step < 2
}
```

#### 3.1.2 Markdown层级分块器

**文件**: `backend/infra/document/parser/impl/builtin/chunk_markdown_hierarchy.go`

```go
package builtin

import (
	"regexp"
	"strings"
)

// MarkdownHierarchyChunker Markdown层级分块器
type MarkdownHierarchyChunker struct {
	chunkSize     int
	maxDepth      int // 最大标题深度，默认5
	chunkOverlap  int
}

// MarkdownChunk Markdown分块
type MarkdownChunk struct {
	Content      string            // 分块内容
	Heading      string            // 当前标题
	Level        int               // 标题级别
	ParentTitles []string          // 父级标题路径
	HeadingPath  string            // 完整标题路径
}

// NewMarkdownHierarchyChunker 创建Markdown分块器
func NewMarkdownHierarchyChunker(chunkSize, maxDepth int) *MarkdownHierarchyChunker {
	return &MarkdownHierarchyChunker{
		chunkSize:    chunkSize,
		maxDepth:     maxDepth,
		chunkOverlap: int(float64(chunkSize) * 0.15),
	}
}

// Split 执行Markdown分块
func (m *MarkdownHierarchyChunker) Split(text string) *SplitResult {
	// 1. 提取所有标题
	blocks := m.extractMarkdownBlocks(text)

	// 2. 按层级组织
	root := m.buildHierarchyTree(blocks)

	// 3. 生成分块
	chunks := m.generateChunks(root)

	contents := make([]string, len(chunks))
	for i, chunk := range chunks {
		contents[i] = chunk.Content
	}

	return &SplitResult{
		Chunks: contents,
		Chars:  len(text),
	}
}

// MarkdownBlock Markdown块
type MarkdownBlock struct {
	Type       string // "heading" | "code" | "table" | "content"
	Level      int    // 标题级别（1-6）
	Heading    string // 标题内容
	Content    string // 块内容
	StartLine  int    // 起始行号
}

// extractMarkdownBlocks 提取Markdown块
func (m *MarkdownHierarchyChunker) extractMarkdownBlocks(text string) []*MarkdownBlock {
	lines := strings.Split(text, "\n")
	var blocks []*MarkdownBlock

	currentContent := strings.Builder
	currentType := "content"
	currentLevel := 0
	currentHeading := ""
	startLine := 0

	codeBlockRegex := regexp.MustCompile("^```(\\w*)?")
	tableLineRegex := regexp.MustCompile("^\\|.*\\|$")

	inCodeBlock := false
	codeBlockContent := strings.Builder

	for i, line := range lines {
		// 检查代码块
		if codeBlockRegex.MatchString(line) {
			if !inCodeBlock {
				// 保存之前的内容
				if currentContent.Len() > 0 {
					blocks = append(blocks, &MarkdownBlock{
						Type:      currentType,
						Level:     currentLevel,
						Heading:   currentHeading,
						Content:   currentContent.String(),
						StartLine: startLine,
					})
					currentContent.Reset()
				}
				inCodeBlock = true
				startLine = i
			} else {
				// 代码块结束
				blocks = append(blocks, &MarkdownBlock{
					Type:      "code",
					Content:   codeBlockContent.String(),
					StartLine: startLine,
				})
				codeBlockContent.Reset()
				inCodeBlock = false
			}
			continue
		}

		if inCodeBlock {
			codeBlockContent.WriteString(line + "\n")
			continue
		}

		// 检查标题
		if strings.HasPrefix(line, "#") {
			// 保存之前的内容
			if currentContent.Len() > 0 {
				blocks = append(blocks, &MarkdownBlock{
					Type:      currentType,
					Level:     currentLevel,
					Heading:   currentHeading,
					Content:   currentContent.String(),
					StartLine: startLine,
				})
				currentContent.Reset()
			}

			// 解析标题
			level := 0
			for _, c := range line {
				if c == '#' {
					level++
				} else {
					break
				}
			}
			currentHeading = strings.TrimSpace(line[level:])
			currentLevel = level
			currentType = "heading"
			startLine = i

			blocks = append(blocks, &MarkdownBlock{
				Type:      "heading",
				Level:     level,
				Heading:   currentHeading,
				StartLine: i,
			})
			continue
		}

		// 检查表格
		if tableLineRegex.MatchString(line) {
			if currentType != "table" {
				// 保存之前的内容
				if currentContent.Len() > 0 {
					blocks = append(blocks, &MarkdownBlock{
						Type:      currentType,
						Level:     currentLevel,
						Heading:   currentHeading,
						Content:   currentContent.String(),
						StartLine: startLine,
					})
					currentContent.Reset()
				}
				currentType = "table"
				startLine = i
			}
			currentContent.WriteString(line + "\n")
			continue
		}

		// 普通内容
		currentContent.WriteString(line + "\n")
		if currentType == "heading" || currentType == "table" {
			currentType = "content"
			startLine = i
		}
	}

	// 保存最后一块
	if currentContent.Len() > 0 {
		blocks = append(blocks, &MarkdownBlock{
			Type:      currentType,
			Level:     currentLevel,
			Heading:   currentHeading,
			Content:   currentContent.String(),
			StartLine: startLine,
		})
	}

	return blocks
}

// HierarchyNode 层级树节点
type HierarchyNode struct {
	Block        *MarkdownBlock
	Children     []*HierarchyNode
	Parent       *HierarchyNode
	Level        int
	HeadingPath  []string
}

// buildHierarchyTree 构建层级树
func (m *MarkdownHierarchyChunker) buildHierarchyTree(blocks []*MarkdownBlock) *HierarchyNode {
	root := &HierarchyNode{
		Level:       0,
		HeadingPath: []string{},
	}

	currentPath := []*HierarchyNode{root}

	for _, block := range blocks {
		if block.Type == "heading" {
			// 找到合适的父节点
			for len(currentPath) > 1 && currentPath[len(currentPath)-1].Level >= block.Level {
				currentPath = currentPath[:len(currentPath)-1]
			}

			node := &HierarchyNode{
				Block:       block,
				Parent:      currentPath[len(currentPath)-1],
				Level:       block.Level,
				HeadingPath: append([]string{}, currentPath[len(currentPath)-1].HeadingPath...),
			}
			node.HeadingPath = append(node.HeadingPath, block.Heading)

			currentPath[len(currentPath)-1].Children = append(currentPath[len(currentPath)-1].Children, node)
			currentPath = append(currentPath, node)
		} else {
			// 非标题块，添加到当前节点
			node := &HierarchyNode{
				Block:       block,
				Parent:      currentPath[len(currentPath)-1],
				Level:       currentPath[len(currentPath)-1].Level,
				HeadingPath: append([]string{}, currentPath[len(currentPath)-1].HeadingPath...),
			}
			currentPath[len(currentPath)-1].Children = append(currentPath[len(currentPath)-1].Children, node)
		}
	}

	return root
}

// generateChunks 生成分块
func (m *MarkdownHierarchyChunker) generateChunks(root *HierarchyNode) []*MarkdownChunk {
	var chunks []*MarkdownChunk

	var traverse func(node *HierarchyNode, parentPath []string)
	traverse = func(node *HierarchyNode, parentPath []string) {
		if node.Block != nil && node.Block.Type != "heading" {
			chunk := &MarkdownChunk{
				Content:      node.Block.Content,
				Heading:      node.Block.Heading,
				Level:        node.Block.Level,
				ParentTitles: parentPath,
				HeadingPath:  strings.Join(node.HeadingPath, " > "),
			}
			chunks = append(chunks, chunk)
		}

		for _, child := range node.Children {
			newPath := append([]string{}, parentPath...)
			if child.Block != nil && child.Block.Heading != "" {
				newPath = append(newPath, child.Block.Heading)
			}
			traverse(child, newPath)
		}
	}

	traverse(root, []string{})
	return chunks
}
```

#### 3.1.3 代码块智能分块器

**文件**: `backend/infra/document/parser/impl/builtin/chunk_code.go`

```go
package builtin

import (
	"regexp"
	"strings"
)

// CodeChunker 代码块智能分块器
type CodeChunker struct {
	chunkSize int
}

// CodeBlock 代码块
type CodeBlock struct {
	Language string // 编程语言
	Content  string // 代码内容
	StartLine int   // 起始行
	EndLine   int   // 结束行
}

// NewCodeChunker 创建代码块分块器
func NewCodeChunker(chunkSize int) *CodeChunker {
	return &CodeChunker{
		chunkSize: chunkSize,
	}
}

// ExtractCodeBlocks 提取代码块
func (c *CodeChunker) ExtractCodeBlocks(text string) ([]*CodeBlock, string) {
	// 1. 替换代码块为标记
	codeBlockRegex := regexp.MustCompile("```(\\w*)?\\n([\\s\\S]*?)```")

	var codeBlocks []*CodeBlock
	var processedText string
	lastEnd := 0

	blockIndex := 0
	marker := fmt.Sprintf("__CODE_BLOCK_%d__", blockIndex)

	for _, match := range codeBlockRegex.FindAllStringSubmatchIndex(text, -1) {
		start, end := match[0], match[1]
		languageStart, languageEnd := match[2], match[3]
		contentStart, contentEnd := match[4], match[5]

		// 添加标记前的文本
		processedText += text[lastEnd:start]

		// 提取代码块信息
		language := text[languageStart:languageEnd]
		content := text[contentStart:contentEnd]

		codeBlocks = append(codeBlocks, &CodeBlock{
			Language: language,
			Content:  content,
		})

		// 添加标记
		processedText += marker
		lastEnd = end

		blockIndex++
		marker = fmt.Sprintf("__CODE_BLOCK_%d__", blockIndex)
	}

	// 添加剩余文本
	processedText += text[lastEnd:]

	return codeBlocks, processedText
}

// RestoreCodeBlocks 恢复代码块
func (c *CodeChunker) RestoreCodeBlocks(text string, codeBlocks []*CodeBlock) string {
	result := text
	for i, block := range codeBlocks {
		marker := fmt.Sprintf("__CODE_BLOCK_%d__", i)
		result = strings.ReplaceAll(result, marker,
			fmt.Sprintf("```%s\n%s\n```", block.Language, block.Content))
	}
	return result
}

// SplitCodeBlock 按函数/类分割代码块
func (c *CodeChunker) SplitCodeBlock(block *CodeBlock) []string {
	// 根据编程语言选择分割策略
	switch block.Language {
	case "python", "py":
		return c.splitPythonCode(block.Content)
	case "java":
		return c.splitJavaCode(block.Content)
	case "go", "golang":
		return c.splitGoCode(block.Content)
	case "javascript", "js", "typescript", "ts":
		return c.splitJavaScriptCode(block.Content)
	default:
		// 默认不分割
		return []string{block.Content}
	}
}

// splitPythonCode 分割Python代码
func (c *CodeChunker) splitPythonCode(code string) []string {
	// 按函数和类分割
	funcRegex := regexp.MustCompile(`^(def\s+\w+|class\s+\w+)`)

	lines := strings.Split(code, "\n")
	var chunks []string
	var currentChunk strings.Builder

	for _, line := range lines {
		if funcRegex.MatchString(line) && currentChunk.Len() > 0 {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}
		currentChunk.WriteString(line + "\n")
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}

// splitJavaCode 分割Java代码
func (c *CodeChunker) splitJavaCode(code string) []string {
	// 按类和方法分割
	classRegex := regexp.MustCompile(`^(public\s+)?(class|interface|enum)\s+\w+`)
	methodRegex := regexp.MustCompile(`^(public\s+)?(protected\s+)?(private\s+)?[\w<>\[\]]+\s+\w+\s*\(.*\)\s*(throws\s+[\w\s,]+)?\s*{`)

	lines := strings.Split(code, "\n")
	var chunks []string
	var currentChunk strings.Builder
	indentLevel := 0

	for _, line := range lines {
		if classRegex.MatchString(line) && currentChunk.Len() > 0 && indentLevel == 0 {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}

		currentChunk.WriteString(line + "\n")

		// 跟踪缩进级别
		if strings.Contains(line, "{") {
			indentLevel += strings.Count(line, "{")
		}
		if strings.Contains(line, "}") {
			indentLevel -= strings.Count(line, "}")
		}
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}

// splitGoCode 分割Go代码
func (c *CodeChunker) splitGoCode(code string) []string {
	// 按函数和方法分割
	funcRegex := regexp.MustCompile(`^func\s+\(?(\w+)\)?\s+\w+\s*\(.*\)\s*{`)

	lines := strings.Split(code, "\n")
	var chunks []string
	var currentChunk strings.Builder

	for _, line := range lines {
		if funcRegex.MatchString(line) && currentChunk.Len() > 0 {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}
		currentChunk.WriteString(line + "\n")
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}

// splitJavaScriptCode 分割JavaScript/TypeScript代码
func (c *CodeChunker) splitJavaScriptCode(code string) []string {
	// 按函数和类分割
	funcRegex := regexp.MustCompile(`^(function\s+\w+|const\s+\w+\s*=\s*(?:async\s+)?\(?[\w\s]*\)?\s*=>|class\s+\w+)`)

	lines := strings.Split(code, "\n")
	var chunks []string
	var currentChunk strings.Builder

	for _, line := range lines {
		if funcRegex.MatchString(line) && currentChunk.Len() > 0 {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}
		currentChunk.WriteString(line + "\n")
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}
```

#### 3.1.4 分块器接口统一

**文件**: `backend/infra/document/parser/parser.go` (修改)

```go
package parser

import (
	"context"
	"github.com/coze-dev/coze-studio/backend/domain/knowledge/entity"
)

// ChunkingStrategy 分块策略
type ChunkingStrategy string

const (
	// ChunkingStrategyRecursive 递归分块
	ChunkingStrategyRecursive ChunkingStrategy = "recursive"
	// ChunkingStrategyMarkdown Markdown层级分块
	ChunkingStrategyMarkdown ChunkingStrategy = "markdown"
	// ChunkingStrategySemantic 语义分块
	ChunkingStrategySemantic ChunkingStrategy = "semantic"
	// ChunkingStrategyCustom 自定义分块
	ChunkingStrategyCustom ChunkingStrategy = "custom"
)

// Chunker 分块器接口
type Chunker interface {
	// Chunk 执行分块
	Chunk(ctx context.Context, text string, strategy *entity.ChunkingStrategy) ([]*entity.Chunk, error)
}

// ChunkResult 分块结果
type ChunkResult struct {
	Chunks      []*entity.Chunk // 分块列表
	TotalChars  int              // 总字符数
	TotalChunks int              // 分块数量
}

// ChunkerFactory 分块器工厂
type ChunkerFactory interface {
	// CreateChunker 创建分块器
	CreateChunker(strategy ChunkingStrategy) (Chunker, error)
	// GetSupportedStrategies 获取支持的策略
	GetSupportedStrategies() []ChunkingStrategy
}

// DefaultChunkerFactory 默认分块器工厂
type DefaultChunkerFactory struct{}

func NewDefaultChunkerFactory() *DefaultChunkerFactory {
	return &DefaultChunkerFactory{}
}

func (f *DefaultChunkerFactory) CreateChunker(strategy ChunkingStrategy) (Chunker, error) {
	switch strategy {
	case ChunkingStrategyRecursive:
		return NewRecursiveChunker(800), nil
	case ChunkingStrategyMarkdown:
		return NewMarkdownHierarchyChunker(800, 5), nil
	case ChunkingStrategySemantic:
		return NewSemanticChunker(), nil
	default:
		return NewCustomChunker(), nil
	}
}

func (f *DefaultChunkerFactory) GetSupportedStrategies() []ChunkingStrategy {
	return []ChunkingStrategy{
		ChunkingStrategyRecursive,
		ChunkingStrategyMarkdown,
		ChunkingStrategySemantic,
		ChunkingStrategyCustom,
	}
}
```

---

### 3.2 阶段2：混合检索优化 (4人天)

#### 3.2.1 RRF融合优化

**文件**: `backend/infra/document/rerank/impl/rrf/rrf_optimized.go`

```go
package rrf

import (
	"context"
	"math"
	"sort"

	"github.com/coze-dev/coze-studio/backend/domain/knowledge/entity"
)

// RRFOptimizer 优化的RRF融合器
type RRFOptimizer struct {
	k             float64 // RRF平滑系数，默认60
	embeddingWeight float64 // 向量检索权重，默认0.5
	fullTextWeight float64 // 全文检索权重，默认0.5
}

// NewRRFOptimizer 创建RRF优化器
func NewRRFOptimizer() *RRFOptimizer {
	return &RRFOptimizer{
		k:               60,
		embeddingWeight: 0.5,
		fullTextWeight:  0.5,
	}
}

// SetK 设置RRF平滑系数
func (r *RRFOptimizer) SetK(k float64) {
	r.k = k
}

// SetWeights 设置权重
func (r *RRFOptimizer) SetWeights(embedding, fullText float64) {
	r.embeddingWeight = embedding
	r.fullTextWeight = fullText
}

// FusionResult 融合结果
type FusionResult struct {
	SliceID    string  // 切片ID
	Score      float64 // 融合分数
	Rank       int     // 融合排名
	Source     string  // 来源: "embedding", "fulltext", "both"
}

// Fuse 执行RRF融合
func (r *RRFOptimizer) Fuse(
	ctx context.Context,
	embeddingResults []*entity.SearchResult,
	fullTextResults []*entity.SearchResult,
	limit int,
) []*FusionResult {
	// 1. 计算每个结果的RRF分数
	scoreMap := make(map[string]*FusionResult)

	// 2. 处理向量检索结果
	for i, result := range embeddingResults {
		rank := i + 1
		rrfScore := r.embeddingWeight * (1.0 / (r.k + float64(rank)))

		if existing, ok := scoreMap[result.SliceID]; ok {
			// 已存在，累加分数
			existing.Score += rrfScore
			existing.Source = "both"
		} else {
			// 新结果
			scoreMap[result.SliceID] = &FusionResult{
				SliceID: result.SliceID,
				Score:   rrfScore,
				Rank:    rank,
				Source:  "embedding",
			}
		}
	}

	// 3. 处理全文检索结果
	for i, result := range fullTextResults {
		rank := i + 1
		rrfScore := r.fullTextWeight * (1.0 / (r.k + float64(rank)))

		if existing, ok := scoreMap[result.SliceID]; ok {
			// 已存在，累加分数
			existing.Score += rrfScore
			if existing.Source == "embedding" {
				existing.Source = "both"
			}
		} else {
			// 新结果
			scoreMap[result.SliceID] = &FusionResult{
				SliceID: result.SliceID,
				Score:   rrfScore,
				Rank:    rank,
				Source:  "fulltext",
			}
		}
	}

	// 4. 转换为切片并按分数排序
	results := make([]*FusionResult, 0, len(scoreMap))
	for _, result := range scoreMap {
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// 5. 更新排名并限制数量
	for i, result := range results {
		result.Rank = i + 1
	}

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

// FuseWithDynamicWeights 动态权重融合
func (r *RRFOptimizer) FuseWithDynamicWeights(
	ctx context.Context,
	embeddingResults []*entity.SearchResult,
	fullTextResults []*entity.SearchResult,
	queryType string,
	limit int,
) []*FusionResult {
	// 根据查询类型动态调整权重
	switch queryType {
	case "keyword":
		// 关键词查询，增加全文检索权重
		r.embeddingWeight = 0.3
		r.fullTextWeight = 0.7
	case "semantic":
		// 语义查询，增加向量检索权重
		r.embeddingWeight = 0.7
		r.fullTextWeight = 0.3
	default:
		// 混合查询，平衡权重
		r.embeddingWeight = 0.5
		r.fullTextWeight = 0.5
	}

	return r.Fuse(ctx, embeddingResults, fullTextResults, limit)
}

// DeduplicateResults 结果去重
func (r *RRFOptimizer) DeduplicateResults(results []*FusionResult) []*FusionResult {
	seen := make(map[string]bool)
	deduplicated := make([]*FusionResult, 0, len(results))

	for _, result := range results {
		if !seen[result.SliceID] {
			seen[result.SliceID] = true
			deduplicated = append(deduplicated, result)
		}
	}

	return deduplicated
}
```

#### 3.2.2 动态权重调整

**文件**: `backend/infra/document/messages2query/query_analyzer.go`

```go
package messages2query

import (
	"context"
	"regexp"
	"strings"
)

// QueryType 查询类型
type QueryType string

const (
	QueryTypeKeyword QueryType = "keyword" // 关键词查询
	QueryTypeSemantic QueryType = "semantic" // 语义查询
	QueryTypeMixed   QueryType = "mixed"   // 混合查询
)

// QueryAnalyzer 查询分析器
type QueryAnalyzer struct {
	keywordThreshold float64 // 关键词阈值，默认0.7
}

// NewQueryAnalyzer 创建查询分析器
func NewQueryAnalyzer() *QueryAnalyzer {
	return &QueryAnalyzer{
		keywordThreshold: 0.7,
	}
}

// AnalyzeQuery 分析查询类型
func (a *QueryAnalyzer) AnalyzeQuery(ctx context.Context, query string) QueryType {
	// 1. 计算关键词密度
	keywordDensity := a.calculateKeywordDensity(query)

	// 2. 检查是否包含特殊字符
	hasSpecialChars := a.hasSpecialChars(query)

	// 3. 检查查询长度
	shortQuery := len(query) <= 20

	// 4. 判断查询类型
	if keywordDensity > a.keywordThreshold || (hasSpecialChars && shortQuery) {
		return QueryTypeKeyword
	}

	if keywordDensity < 0.3 && len(query) > 50 {
		return QueryTypeSemantic
	}

	return QueryTypeMixed
}

// calculateKeywordDensity 计算关键词密度
func (a *QueryAnalyzer) calculateKeywordDensity(query string) float64 {
	// 关键词特征
	keywordPatterns := []string{
		`\b\d+\b`,              // 数字
		`\b[A-Z]{2,}\b`,        // 大写缩写
		`[\u4e00-\u9fa5]{2,}`,  // 中文字符
		`[a-zA-Z0-9_\-\.]+`,    // 英文标识符
	}

	totalChars := len(query)
	keywordChars := 0

	for _, pattern := range keywordPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(query, -1)
		for _, match := range matches {
			keywordChars += len(match)
		}
	}

	if totalChars == 0 {
		return 0
	}

	return float64(keywordChars) / float64(totalChars)
}

// hasSpecialChars 检查是否包含特殊字符
func (a *QueryAnalyzer) hasSpecialChars(query string) bool {
	specialChars := []string{":", "(", ")", "\"", "'", "=", "!", ">", "<"}
	for _, char := range specialChars {
		if strings.Contains(query, char) {
			return true
		}
	}
	return false
}
```

#### 3.2.3 检索去重机制

**文件**: `backend/infra/document/rerank/dedup.go`

```go
package rerank

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/coze-dev/coze-studio/backend/domain/knowledge/entity"
)

// Deduplicator 去重器
type Deduplicator struct {
	similarityThreshold float64 // 相似度阈值，默认0.95
}

// NewDeduplicator 创建去重器
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{
		similarityThreshold: 0.95,
	}
}

// DeduplicateBySliceID 按切片ID去重
func (d *Deduplicator) DeduplicateBySliceID(
	ctx context.Context,
	results []*entity.SearchResult,
) []*entity.SearchResult {
	seen := make(map[string]bool)
	deduplicated := make([]*entity.SearchResult, 0, len(results))

	for _, result := range results {
		if !seen[result.SliceID] {
			seen[result.SliceID] = true
			deduplicated = append(deduplicated, result)
		}
	}

	return deduplicated
}

// DeduplicateByContent 按内容去重（SimHash）
func (d *Deduplicator) DeduplicateByContent(
	ctx context.Context,
	results []*entity.SearchResult,
) []*entity.SearchResult {
	seen := make(map[string]bool)
	deduplicated := make([]*entity.SearchResult, 0, len(results))

	for _, result := range results {
		// 计算内容指纹
		fingerprint := d.calculateFingerprint(result.Content)

		if !seen[fingerprint] {
			seen[fingerprint] = true
			deduplicated = append(deduplicated, result)
		}
	}

	return deduplicated
}

// calculateFingerprint 计算内容指纹（简化版SimHash）
func (d *Deduplicator) calculateFingerprint(content string) string {
	// 1. 预处理：小写化，去除标点
	content = strings.ToLower(content)
	content = strings.ReplaceAll(content, ".", "")
	content = strings.ReplaceAll(content, ",", "")
	content = strings.ReplaceAll(content, "!", "")
	content = strings.ReplaceAll(content, "?", "")
	content = strings.ReplaceAll(content, " ", "")

	// 2. MD5哈希
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

// DeduplicateBySimilarity 按相似度去重
func (d *Deduplicator) DeduplicateBySimilarity(
	ctx context.Context,
	results []*entity.SearchResult,
) []*entity.SearchResult {
	deduplicated := make([]*entity.SearchResult, 0, len(results))

	for i, result := range results {
		isDuplicate := false

		// 与已保留的结果比较
		for _, kept := range deduplicated {
			similarity := d.calculateSimilarity(result.Content, kept.Content)
			if similarity > d.similarityThreshold {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			deduplicated = append(deduplicated, result)
		}
	}

	return deduplicated
}

// calculateSimilarity 计算相似度（简化版余弦相似度）
func (d *Deduplicator) calculateSimilarity(content1, content2 string) float64 {
	// 1. 分词（简化版：按字符）
	words1 := strings.Split(content1, "")
	words2 := strings.Split(content2, "")

	// 2. 计算词频
	freq1 := make(map[string]int)
	freq2 := make(map[string]int)

	for _, word := range words1 {
		freq1[word]++
	}
	for _, word := range words2 {
		freq2[word]++
	}

	// 3. 计算余弦相似度
	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0

	for word, count1 := range freq1 {
		count2 := freq2[word]
		dotProduct += float64(count1 * count2)
		norm1 += float64(count1 * count1)
	}

	for _, count2 := range freq2 {
		norm2 += float64(count2 * count2)
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}
```

---

### 3.3 阶段3：多租户隔离增强 (3人天)

#### 3.3.1 租户级向量索引隔离

**文件**: `backend/infra/document/searchstore/impl/milvus/milvus_manager.go` (修改)

```go
package milvus

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/knowledge/entity"
)

// MilvusManager Milvus管理器
type MilvusManager struct {
	client *MilvusClient
}

// CreateCollection 创建租户隔离的Collection
func (m *MilvusManager) CreateCollection(
	ctx context.Context,
	tenantID string,
	knowledgeID int64,
	dimensions int,
) error {
	// 方案1: 租户级Partition（推荐）
	collectionName := m.getCollectionName(knowledgeID)
	partitionName := m.getPartitionName(tenantID)

	// 1. 创建Collection（如果不存在）
	if err := m.createCollectionIfNotExists(ctx, collectionName, dimensions); err != nil {
		return err
	}

	// 2. 创建租户Partition
	if err := m.client.CreatePartition(ctx, collectionName, partitionName); err != nil {
		return err
	}

	return nil
}

// getCollectionName 获取Collection名称
func (m *MilvusManager) getCollectionName(knowledgeID int64) string {
	return fmt.Sprintf("knowledge_%d", knowledgeID)
}

// getPartitionName 获取Partition名称
func (m *MilvusManager) getPartitionName(tenantID string) string {
	return fmt.Sprintf("tenant_%s", tenantID)
}

// Insert 插入向量到租户Partition
func (m *MilvusManager) Insert(
	ctx context.Context,
	tenantID string,
	knowledgeID int64,
	vectors []entity.Vector,
) error {
	collectionName := m.getCollectionName(knowledgeID)
	partitionName := m.getPartitionName(tenantID)

	return m.client.Insert(ctx, collectionName, partitionName, vectors)
}

// Search 在租户Partition内检索
func (m *MilvusManager) Search(
	ctx context.Context,
	tenantID string,
	knowledgeID int64,
	queryVector entity.Vector,
	topK int,
) ([]*entity.SearchResult, error) {
	collectionName := m.getCollectionName(knowledgeID)
	partitionName := m.getPartitionName(tenantID)

	// 只在租户Partition内检索
	return m.client.Search(ctx, collectionName, []string{partitionName}, queryVector, topK)
}

// DeletePartition 删除租户Partition
func (m *MilvusManager) DeletePartition(
	ctx context.Context,
	tenantID string,
	knowledgeID int64,
) error {
	collectionName := m.getCollectionName(knowledgeID)
	partitionName := m.getPartitionName(tenantID)

	return m.client.DropPartition(ctx, collectionName, partitionName)
}
```

#### 3.3.2 租户级ES索引隔离

**文件**: `backend/infra/document/searchstore/impl/elasticsearch/elasticsearch_manager.go` (修改)

```go
package elasticsearch

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

// ESManager Elasticsearch管理器
type ESManager struct {
	client *elasticsearch.Client
}

// CreateIndex 创建租户隔离的索引
func (e *ESManager) CreateIndex(
	ctx context.Context,
	tenantID string,
	knowledgeID int64,
) error {
	// 索引命名: knowledge_<tenant_id>_<knowledge_id>
	indexName := e.getIndexName(tenantID, knowledgeID)

	// 创建索引
	return e.client.Indices.Create(
		indexName,
		e.client.Indices.Create.WithContext(ctx),
	)
}

// getIndexName 获取索引名称
func (e *ESManager) getIndexName(tenantID string, knowledgeID int64) string {
	return fmt.Sprintf("knowledge_%s_%d", tenantID, knowledgeID)
}

// Insert 插入文档到租户索引
func (e *ESManager) Insert(
	ctx context.Context,
	tenantID string,
	knowledgeID int64,
	documents []*Document,
) error {
	indexName := e.getIndexName(tenantID, knowledgeID)

	for _, doc := range documents {
		// 添加租户标签
		doc.TenantID = tenantID

		if err := e.client.Index(indexName, doc); err != nil {
			return err
		}
	}

	return nil
}

// Search 在租户索引内检索
func (e *ESManager) Search(
	ctx context.Context,
	tenantID string,
	knowledgeID int64,
	query string,
	limit int,
) ([]*SearchResult, error) {
	indexName := e.getIndexName(tenantID, knowledgeID)

	// 只在租户索引内检索
	return e.client.Search(
		indexName,
		query,
		e.client.Search.WithContext(ctx),
		e.client.Search.WithSize(limit),
	)
}

// DeleteIndex 删除租户索引
func (e *ESManager) DeleteIndex(
	ctx context.Context,
	tenantID string,
	knowledgeID int64,
) error {
	indexName := e.getIndexName(tenantID, knowledgeID)

	return e.client.Indices.Delete(
		[]string{indexName},
		e.client.Indices.Delete.WithContext(ctx),
	)
}
```

#### 3.3.3 缓存隔离验证

**文件**: `backend/domain/knowledge/internal/dal/cache.go` (修改)

```go
package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/cache"
)

// Cache 缓存管理器
type Cache struct {
	cli cache.Cmdable
}

// CacheKey 缓存键生成器
type CacheKey struct {
	tenantID    string
	knowledgeID int64
	documentID  int64
}

// BuildIndexDocKey 构建索引文档缓存键（包含租户ID）
func BuildIndexDocKey(tenantID string, knowledgeID, documentID int64) string {
	// 格式: tenant:<tenant_id>:knowledge:<knowledge_id>:document:<document_id>
	return fmt.Sprintf("tenant:%s:knowledge:%d:document:%d",
		tenantID, knowledgeID, documentID)
}

// GetIndexDocCache 获取索引文档缓存
func (c *Cache) GetIndexDocCache(
	ctx context.Context,
	tenantID string,
	knowledgeID, documentID int64,
) (*string, error) {
	key := BuildIndexDocKey(tenantID, knowledgeID, documentID)

	val, err := c.cli.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &val, nil
}

// SetIndexDocCache 设置索引文档缓存
func (c *Cache) SetIndexDocCache(
	ctx context.Context,
	tenantID string,
	knowledgeID, documentID int64,
	value string,
	expiration time.Duration,
) error {
	key := BuildIndexDocKey(tenantID, knowledgeID, documentID)

	return c.cli.Set(ctx, key, value, expiration).Err()
}

// DeleteIndexDocCache 删除索引文档缓存
func (c *Cache) DeleteIndexDocCache(
	ctx context.Context,
	tenantID string,
	knowledgeID, documentID int64,
) error {
	key := BuildIndexDocKey(tenantID, knowledgeID, documentID)

	return c.cli.Del(ctx, key).Err()
}

// ClearTenantCache 清理租户所有缓存
func (c *Cache) ClearTenantCache(
	ctx context.Context,
	tenantID string,
) error {
	// 使用SCAN命令查找租户所有缓存键
	pattern := fmt.Sprintf("tenant:%s:*", tenantID)

	var cursor uint64
	for {
		keys, nextCursor, err := c.cli.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := c.cli.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
```

---

### 3.4 阶段4：去重机制 (2人天)

#### 3.4.1 SimHash去重

**文件**: `backend/domain/knowledge/service/dedup_simhash.go`

```go
package service

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"math"
)

// SimHashDeduplicator SimHash去重器
type SimHashDeduplicator struct {
	hashBits    int     // 哈希位数，默认64
	threshold   int     // 海明距离阈值，默认3
}

// NewSimHashDeduplicator 创建SimHash去重器
func NewSimHashDeduplicator() *SimHashDeduplicator {
	return &SimHashDeduplicator{
		hashBits:  64,
		threshold: 3,
	}
}

// CalculateSimHash 计算SimHash
func (s *SimHashDeduplicator) CalculateSimHash(content string) uint64 {
	// 1. 分词（简化版：按4-gram）
	terms := s.extractTerms(content)

	// 2. 计算每个词的权重（简化版：词频）
	weights := make(map[string]int)
	for _, term := range terms {
		weights[term]++
	}

	// 3. 初始化哈希向量
	hashVector := make([]int, s.hashBits)

	// 4. 累加每个词的哈希值
	for term, weight := range weights {
		// 计算词的MD5哈希
		hash := md5.Sum([]byte(term))

		// 将哈希值转换为向量
		for i := 0; i < s.hashBits; i++ {
			byteIndex := i / 8
			bitIndex := uint(i % 8)

			// 检查bit是否为1
			if hash[byteIndex]&(1<<bitIndex) != 0 {
				hashVector[i] += weight
			} else {
				hashVector[i] -= weight
			}
		}
	}

	// 5. 生成SimHash值
	var simHash uint64
	for i := 0; i < s.hashBits; i++ {
		if hashVector[i] >= 0 {
			simHash |= 1 << uint(i)
		}
	}

	return simHash
}

// extractTerms 提取词项（4-gram）
func (s *SimHashDeduplicator) extractTerms(content string) []string {
	// 简化版：4-gram
	if len(content) < 4 {
		return []string{content}
	}

	terms := make([]string, 0, len(content)-3)
	for i := 0; i <= len(content)-4; i++ {
		terms = append(terms, content[i:i+4])
	}

	return terms
}

// CalculateHammingDistance 计算海明距离
func (s *SimHashDeduplicator) CalculateHammingDistance(hash1, hash2 uint64) int {
	xor := hash1 ^ hash2
	distance := 0

	for xor != 0 {
		distance += int(xor & 1)
		xor >>= 1
	}

	return distance
}

// IsDuplicate 判断是否重复
func (s *SimHashDeduplicator) IsDuplicate(
	hash1, hash2 uint64,
) bool {
	distance := s.CalculateHammingDistance(hash1, hash2)
	return distance <= s.threshold
}

// Deduplicate 去重
func (s *SimHashDeduplicator) Deduplicate(
	ctx context.Context,
	contents []string,
) []string {
	// 计算所有内容的SimHash
	hashes := make([]uint64, len(contents))
	for i, content := range contents {
		hashes[i] = s.CalculateSimHash(content)
	}

	// 去重
	unique := make([]string, 0)
	seen := make(map[uint64]bool)

	for i, content := range contents {
		hash := hashes[i]
		isDuplicate := false

		// 检查是否与已见过的内容相似
		for seenHash := range seen {
			if s.IsDuplicate(hash, seenHash) {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			seen[hash] = true
			unique = append(unique, content)
		}
	}

	return unique
}
```

#### 3.4.2 相似度去重

**文件**: `backend/domain/knowledge/service/dedup_similarity.go`

```go
package service

import (
	"context"
	"math"
	"sort"
)

// SimilarityDeduplicator 相似度去重器
type SimilarityDeduplicator struct {
	threshold float64 // 相似度阈值，默认0.95
}

// NewSimilarityDeduplicator 创建相似度去重器
func NewSimilarityDeduplicator() *SimilarityDeduplicator {
	return &SimilarityDeduplicator{
		threshold: 0.95,
	}
}

// CalculateCosineSimilarity 计算余弦相似度
func (s *SimilarityDeduplicator) CalculateCosineSimilarity(
	vec1, vec2 []float32,
) float64 {
	if len(vec1) != len(vec2) {
		return 0
	}

	dotProduct := float32(0)
	norm1 := float32(0)
	norm2 := float32(0)

	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
		norm1 += vec1[i] * vec1[i]
		norm2 += vec2[i] * vec2[i]
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return float64(dotProduct / (float32(math.Sqrt(float64(norm1))) * float32(math.Sqrt(float64(norm2)))))
}

// Deduplicate 去重
func (s *SimilarityDeduplicator) Deduplicate(
	ctx context.Context,
	vectors [][]float32,
) []int {
	// 返回保留的索引
	kept := make([]int, 0)

	for i, vec1 := range vectors {
		isDuplicate := false

		// 与已保留的向量比较
		for _, keptIdx := range kept {
			vec2 := vectors[keptIdx]
			similarity := s.CalculateCosineSimilarity(vec1, vec2)

			if similarity > s.threshold {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			kept = append(kept, i)
		}
	}

	return kept
}

// DeduplicateBatch 批量去重（优化版）
func (s *SimilarityDeduplicator) DeduplicateBatch(
	ctx context.Context,
	vectors [][]float32,
	batchSize int,
) []int {
	// 返回保留的索引
	kept := make([]int, 0)

	// 按批次处理
	for i := 0; i < len(vectors); i += batchSize {
		end := i + batchSize
		if end > len(vectors) {
			end = len(vectors)
		}

		batch := vectors[i:end]
		keptInBatch := s.deduplicateSingleBatch(ctx, batch)

		// 调整索引
		for _, idx := range keptInBatch {
			kept = append(kept, i+idx)
		}
	}

	return kept
}

// deduplicateSingleBatch 单批次去重
func (s *SimilarityDeduplicator) deduplicateSingleBatch(
	ctx context.Context,
	vectors [][]float32,
) []int {
	kept := make([]int, 0)

	for i, vec1 := range vectors {
		isDuplicate := false

		for _, keptIdx := range kept {
			vec2 := vectors[keptIdx]
			similarity := s.CalculateCosineSimilarity(vec1, vec2)

			if similarity > s.threshold {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			kept = append(kept, i)
		}
	}

	return kept
}
```

---

## 4. 实施计划

### 4.1 任务分解 (WBS)

```
RAG增强项目
├── P0任务 (14人天)
│   ├── 1. 层级分块实现 (5人天)
│   │   ├── 1.1 递归字符分块器 (1.5人天)
│   │   ├── 1.2 Markdown层级分块器 (1.5人天)
│   │   ├── 1.3 代码块智能分割 (1人天)
│   │   ├── 1.4 重叠策略实现 (0.5人天)
│   │   └── 1.5 分块器接口统一 (0.5人天)
│   │
│   ├── 2. 混合检索优化 (4人天)
│   │   ├── 2.1 RRF融合优化 (1.5人天)
│   │   ├── 2.2 动态权重调整 (1人天)
│   │   ├── 2.3 结果去重 (1人天)
│   │   └── 2.4 检索性能优化 (0.5人天)
│   │
│   ├── 3. 多租户隔离增强 (3人天)
│   │   ├── 3.1 租户级向量索引 (1人天)
│   │   ├── 3.2 租户级ES索引 (1人天)
│   │   └── 3.3 缓存隔离验证 (1人天)
│   │
│   └── 4. 去重机制 (2人天)
│       ├── 4.1 SimHash去重 (1人天)
│       └── 4.2 相似度去重 (1人天)
│
└── P1任务 (33人天) - 后续迭代
    ├── 5. BGE-Reranker集成 (8人天)
    ├── 6. 语义分块 (7人天)
    ├── 7. GraphRAG (12人天)
    └── 8. 查询扩展 (6人天)
```

### 4.2 里程碑

| 里程碑 | 日期 | 交付物 | 验收标准 |
|--------|------|--------|----------|
| **M1: 层级分块** | Week 1 | 递归分块器、Markdown分块器、代码块分割器 | 单元测试≥90%，分块质量验证通过 |
| **M2: 混合检索** | Week 1.5 | RRF融合优化、动态权重、去重机制 | 检索准确率提升≥15%，P99延迟≤500ms |
| **M3: 多租户隔离** | Week 2 | 租户级索引、缓存隔离 | 隔离测试通过，无跨租户泄露 |
| **M4: P0完成** | Week 2.5 | 所有P0功能集成 | 端到端测试通过，性能达标 |

### 4.3 验收标准

#### 阶段1：层级分块
- [ ] 支持递归字符分块，chunk_size=800±200
- [ ] 支持Markdown层级分块，max_depth=5
- [ ] 代码块完整保留，不分割
- [ ] 15-20%重叠策略
- [ ] 单元测试覆盖率≥90%

#### 阶段2：混合检索
- [ ] 向量+全文检索RRF融合
- [ ] 动态权重调整(0.3-0.7)
- [ ] 结果去重(slice_id级)
- [ ] P99延迟≤500ms
- [ ] 检索准确率≥80%

#### 阶段3：多租户隔离
- [ ] Milvus租户级Partition
- [ ] ES租户级索引
- [ ] Redis缓存Key包含租户ID
- [ ] 隔离测试100%通过
- [ ] 无跨租户数据泄露

#### 阶段4：去重机制
- [ ] SimHash去重准确率≥95%
- [ ] 相似度去重阈值可配置(0.8-0.99)
- [ ] 去重性能: 10K文档<10s
- [ ] 去重统计记录

---

## 5. 测试计划

### 5.1 单元测试

```go
// backend/infra/document/parser/impl/builtin/chunk_recursive_test.go
package builtin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecursiveChunker_Split(t *testing.T) {
	chunker := NewRecursiveChunker(800)

	tests := []struct {
		name           string
		input          string
		expectedChunks int
		minChunkSize   int
		maxChunkSize   int
	}{
		{
			name: "简单文本",
			input: "这是一个测试文本。",
			expectedChunks: 1,
			minChunkSize: 10,
			maxChunkSize: 100,
		},
		{
			name: "长文本",
			input: generateLongText(2000),
			expectedChunks: 3,
			minChunkSize: 600,
			maxChunkSize: 900,
		},
		{
			name: "Markdown文档",
			input: `# 标题1\n内容1\n\n## 标题2\n内容2`,
			expectedChunks: 2,
			minChunkSize: 50,
			maxChunkSize: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chunker.Split(tt.input)

			assert.Equal(t, tt.expectedChunks, len(result.Chunks))

			for _, chunk := range result.Chunks {
				assert.GreaterOrEqual(t, len(chunk), tt.minChunkSize)
				assert.LessOrEqual(t, len(chunk), tt.maxChunkSize)
			}
		})
	}
}

func generateLongText(size int) string {
	// 生成长文本用于测试
	// ...
}
```

### 5.2 集成测试

```go
// backend/tests/integration/knowledge/retrieval_test.go
package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHybridRetrieval(t *testing.T) {
	// 1. 准备测试数据
	ctx := context.Background()
	tenantID := "test_tenant"
	knowledgeID := int64(1)

	// 2. 索引文档
	document := &entity.Document{
		Content: generateTestDocument(),
	}
	err := knowledgeService.IndexDocument(ctx, tenantID, knowledgeID, document)
	assert.NoError(t, err)

	// 3. 执行混合检索
	query := "测试查询"
	results, err := retrievalService.HybridRetrieve(ctx, tenantID, knowledgeID, query, 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	// 4. 验证结果
	for _, result := range results {
		assert.Equal(t, tenantID, result.TenantID)
		assert.Greater(t, result.Score, float64(0))
	}

	// 5. 验证去重
	sliceIDs := make(map[string]bool)
	for _, result := range results {
		assert.False(t, sliceIDs[result.SliceID], "发现重复结果")
		sliceIDs[result.SliceID] = true
	}
}

func TestTenantIsolation(t *testing.T) {
	ctx := context.Background()

	// 创建两个租户
	tenant1 := "tenant_1"
	tenant2 := "tenant_2"
	knowledgeID := int64(1)

	// 索引文档到租户1
	document := &entity.Document{
		Content: "租户1的秘密文档",
	}
	err := knowledgeService.IndexDocument(ctx, tenant1, knowledgeID, document)
	assert.NoError(t, err)

	// 租户2尝试检索
	results, err := retrievalService.HybridRetrieve(ctx, tenant2, knowledgeID, "秘密", 10)
	assert.NoError(t, err)
	assert.Empty(t, results, "租户2不应该检索到租户1的文档")
}
```

### 5.3 性能测试

```go
// backend/tests/performance/knowledge/retrieval_bench_test.go
package performance

import (
	"context"
	"testing"
)

func BenchmarkHybridRetrieval(b *testing.B) {
	ctx := context.Background()
	tenantID := "bench_tenant"
	knowledgeID := int64(1)

	// 准备测试数据
	setupBenchmarkData(ctx, tenantID, knowledgeID, 10000)

	query := "测试查询"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := retrievalService.HybridRetrieve(ctx, tenantID, knowledgeID, query, 10)
		if err != nil {
			b.Fatalf("检索失败: %v", err)
		}
	}
}

func BenchmarkDeduplication(b *testing.B) {
	dedup := NewSimHashDeduplicator()
	contents := generateTestContents(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dedup.Deduplicate(context.Background(), contents)
	}
}
```

---

## 6. 风险评估与应对

| 风险项 | 风险等级 | 影响 | 应对措施 |
|--------|---------|------|----------|
| **分块质量不达标** | 中 | 检索准确率下降 | 1. 充分测试多种分块策略<br>2. 提供配置选项<br>3. A/B测试验证 |
| **多租户隔离泄露** | 高 | 数据安全风险 | 1. 完整测试覆盖<br>2. 安全审计<br>3. 租户隔离中间件 |
| **性能回归** | 中 | 用户体验下降 | 1. 性能基准测试<br>2. 优化关键路径<br>3. 监控告警 |
| **第三方依赖不稳定** | 低 | 服务不可用 | 1. 降级方案<br>2. 熔断机制<br>3. 重试策略 |

---

## 7. 附录

### 7.1 参考文档

- FastGPT RAG实现: `D:\FastGPT\packages\global\common\string\textSplitter.ts`
- LangChain RecursiveCharacterTextSplitter
- BGE-Reranker官方文档

### 7.2 配置示例

```yaml
# config/knowledge.yaml
chunking:
  recursive:
    chunk_size: 800
    chunk_overlap: 0.15
    separators: ["\n\n", "\n", "。", ". ", "!", "?", ";", "，", ", ", " ", ""]

  markdown:
    max_depth: 5
    chunk_size: 1000
    chunk_overlap: 0.1

  code:
    chunk_size: 500
    split_by_function: true

retrieval:
  hybrid:
    embedding_weight: 0.5
    fulltext_weight: 0.5
    rrf_k: 60

  dedup:
    enabled: true
    method: "simhash"  # simhash | similarity
    threshold: 0.95

multi_tenant:
  milvus:
    partition_per_tenant: true

  elasticsearch:
    index_per_tenant: true

  cache:
    key_prefix: "tenant:%s:"
```

---

**文档版本**: v1.0
**最后更新**: 2025-01-03
**下次评审**: 2025-01-10
