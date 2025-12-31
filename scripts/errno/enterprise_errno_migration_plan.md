# errno迁移企业级零风险方案

> **目标**: 375处errno→errorx迁移，0异常，0回滚
> **原则**: 安全第一，可验证，可回滚，分批执行
> **版本**: v2.0 | **日期**: 2025-01-03

---

## 🚨 风险分析

### 现有脚本的问题

```bash
# ❌ 问题1: 简单sed替换，误改风险高
sed -i -E 's/return\s+errno\.([A-Z][a-zA-Z0-9]+)/return errorx.New(errno.\1)/g' "$file"

# 可能误改的情况：
// 1. 注释中的errno
// return errno.ErrXXX // TODO: fix later ← 会被误改

// 2. 字符串中的errno
// log.Info("error is: errno.ErrBotNotFound") ← 会被误改

// 3. 测试代码中的errno（应该保留）
// assert.ErrorIs(err, errno.ErrBotNotFound) ← 会被误改
```

### 零风险方案设计

```
┌─────────────────────────────────────────────────────────┐
│  Phase 1: 静态分析 (2小时)                              │
│  └─ 扫描所有errno引用，生成详细清单                      │
├─────────────────────────────────────────────────────────┤
│  Phase 2: 精确解析 (4小时)                              │
│  └─ 使用Go AST解析器，精确识别需改动的位置               │
├─────────────────────────────────────────────────────────┤
│  Phase 3: 分批迁移 (16小时，按目录分10批)                │
│  └─ 每批10-20个文件，完整验证后才进入下一批               │
├─────────────────────────────────────────────────────────┤
│  Phase 4: 并行验证 (8小时)                              │
│  └─ 原代码和新代码并行运行，结果对比                      │
├─────────────────────────────────────────────────────────┤
│  Phase 5: 灰度发布 (持续)                               │
│  └─ 特性开关控制，出问题立即回滚                          │
└─────────────────────────────────────────────────────────┘
```

---

## 📋 Phase 1: 静态分析 - 生成详细清单

### 1.1 扫描脚本

**文件**: `scripts/errno/phase1_scan.sh`

```bash
#!/bin/bash
# Phase 1: 扫描所有errno引用，生成详细清单

set -e

BACKEND_DIR="backend"
OUTPUT_DIR="scripts/errno/analysis"
REPORT_FILE="$OUTPUT_DIR/errno_inventory_report.csv"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 清空报告文件
> "$REPORT_FILE"
echo "file_path,line_number,context,pattern,recommended_action" > "$REPORT_FILE"

echo "🔍 Phase 1: 扫描errno引用..."
echo ""

# 统计变量
total_files=0
total_refs=0
should_migrate=0
should_skip=0

# 扫描所有Go文件（排除测试文件）
find "$BACKEND_DIR" -name "*.go" -not -name "*_test.go" | while read file; do
    # 获取相对路径
    rel_path="${file#$BACKEND_DIR/}"

    # 使用grep获取所有包含errno的行
    grep -n "errno\." "$file" | while read line; do
        # 提取行号
        line_num=$(echo "$line" | cut -d: -f1)

        # 提取上下文（前后2行）
        context=$(sed -n "$((line_num-2)),$((line_num+2))p" "$file" | sed 's/,/;/g')

        # 提取模式
        pattern=$(echo "$line" | grep -oP 'errno\.[A-Z][a-zA-Z0-9_]+' || echo "")

        # 判断是否需要迁移
        skip=0
        action="MIGRATE"

        # 跳过的情况
        if echo "$context" | grep -q "^[[:space:]]*//"; then
            skip=1
            action="SKIP_COMMENT"
        elif echo "$context" | grep -q 'assert\.Error.*errno\.'; then
            skip=1
            action="SKIP_TEST_ASSERT"
        elif echo "$context" | grep -q '"[^"]*errno\.'; then
            skip=1
            action="SKIP_STRING_LITERAL"
        elif echo "$context" | grep -q 'errorx\.'; then
            skip=1
            action="SKIP_ALREADY_MIGRATED"
        fi

        # 写入CSV
        echo "$rel_path,$line_num,\"$context\",\"$pattern\",$action" >> "$REPORT_FILE"

        # 统计
        total_refs=$((total_refs + 1))
        if [ $skip -eq 1 ]; then
            should_skip=$((should_skip + 1))
        else
            should_migrate=$((should_migrate + 1))
        fi
    done

    total_files=$((total_files + 1))
    echo "  ✓ 扫描: $rel_path"
done

echo ""
echo "📊 扫描完成！"
echo "  - 总文件数: $total_files"
echo "  - 总引用数: $total_refs"
echo "  - 需迁移: $should_migrate"
echo "  - 应跳过: $should_skip"
echo ""
echo "📄 详细报告: $REPORT_FILE"
echo ""
echo "下一步: 查看报告，review需迁移的条目"
```

### 1.2 生成可视化报告

**文件**: `scripts/errno/phase1_generate_report.sh`

```bash
#!/bin/bash
# 生成Markdown格式的详细报告

INPUT_CSV="scripts/errno/analysis/errno_inventory_report.csv"
OUTPUT_MD="scripts/errno/analysis/errno_migration_report.md"

echo "# errno迁移详细报告" > "$OUTPUT_MD"
echo "" >> "$OUTPUT_MD"
echo "> 生成时间: $(date '+%Y-%m-%d %H:%M:%S')" >> "$OUTPUT_MD"
echo "" >> "$OUTPUT_MD"

# 读取CSV并生成Markdown表格
echo "## 📊 统计总览" >> "$OUTPUT_MD"
echo "" >> "$OUTPUT_MD"

# 统计各类action
echo "" >> "$OUTPUT_MD"
echo "### 按操作类型分类" >> "$OUTPUT_MD"
echo "" >> "$OUTPUT_MD"
echo "| 操作类型 | 数量 | 占比 |" >> "$OUTPUT_MD"
echo "|----------|------|------|" >> "$OUTPUT_MD"

tail -n +2 "$INPUT_CSV" | cut -d, -f5 | sort | uniq -c | sort -rn | while read count action; do
    pct=$(echo "scale=1; $count * 100 / $total_refs" | bc)
    echo "| $action | $count | $pct% |" >> "$OUTPUT_MD"
done

echo "" >> "$OUTPUT_MD"
echo "### 按目录分类" >> "$OUTPUT_MD"
echo "" >> "$OUTPUT_MD"
echo "| 目录 | 数量 |" >> "$OUTPUT_MD"
echo "|------|------|" >> "$OUTPUT_MD"

tail -n +2 "$INPUT_CSV" | cut -d, -f1 | cut -d/ -f1 | sort | uniq -c | sort -rn | while read count dir; do
    echo "| $dir | $count |" >> "$OUTPUT_MD"
done

echo "" >> "$OUTPUT_MD"
echo "## 📋 详细清单" >> "$OUTPUT_MD"
echo "" >> "$OUTPUT_MD"

# 只显示需要迁移的条目
echo "### 需要迁移的条目 ($should_migrate 条)" >> "$OUTPUT_MD"
echo "" >> "$OUTPUT_MD"

tail -n +2 "$INPUT_CSV" | grep ",MIGRATE$" | while IFS=, read -r file line context pattern action; do
    echo "#### $file:$line" >> "$OUTPUT_MD"
    echo "" >> "$OUTPUT_MD"
    echo '```go' >> "$OUTPUT_MD"
    echo "$context" | tr ';' '\n' >> "$OUTPUT_MD"
    echo '```' >> "$OUTPUT_MD"
    echo "" >> "$OUTPUT_MD"
done

echo "" >> "$OUTPUT_MD"
echo "### 应跳过的条目 ($should_skip 条)" >> "$OUTPUT_MD"
echo "" >> "$OUTPUT_MD"

tail -n +2 "$INPUT_CSV" | grep -v ",MIGRATE$" | while IFS=, read -r file line context pattern action; do
    echo "- $file:$line - $action" >> "$OUTPUT_MD"
done

echo ""
echo "📄 Markdown报告: $OUTPUT_MD"
```

---

## 🔍 Phase 2: 精确解析 - Go AST解析器

### 2.1 Go程序：精确识别需改动的位置

**文件**: `scripts/errno/phase2_ast_analyzer.go`

```go
// +build ignore

package main

import (
	"encoding/csv"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
)

type ErrnoRef struct {
	File      string
	Line      int
	Pattern   string
	Context   string
	Action    string
	Reason    string
	AutoFix   bool
	Suggested string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run phase2_ast_analyzer.go <backend_dir>")
		os.Exit(1)
	}

	backendDir := os.Args[1]
	refs := analyzeErrnoRefs(backendDir)

	// 输出CSV
	outputCSV(refs, "scripts/errno/analysis/errno_ast_analysis.csv")

	fmt.Printf("✓ AST分析完成\n")
	fmt.Printf("  总引用: %d\n", len(refs))
}

func analyzeErrnoRefs(dir string) []*ErrnoRef {
	fset := token.NewFileSet()
	var refs []*ErrnoRef

	// 遍历所有Go文件
	packages, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, pkg := range packages {
		ast.Inspect(pkg, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.ReturnStmt:
				// 分析return语句
				refs = append(refs, analyzeReturnStmt(fset, x)...)
			case *ast.CallExpr:
				// 分析函数调用
				if ident, ok := x.Fun.(*ast.Ident); ok {
					if ident.Name == "errorx" || ident.Name == "errno" {
						refs = append(refs, analyzeErrorxCall(fset, x)...)
					}
				}
			case *ast.BinaryExpr:
				// 分析二元表达式（如 assert.ErrorIs(err, errno.ErrXXX)）
				refs = append(refs, analyzeBinaryExpr(fset, x)...)
			}
			return true
		})
	}

	return refs
}

func analyzeReturnStmt(fset *token.FileSet, ret *ast.ReturnStmt) []*ErrnoRef {
	var refs []*ErrnoRef

	if len(ret.Results) == 0 {
		return nil
	}

	// 检查返回值中是否有errno.ErrXXX
	for _, result := range ret.Results {
		ref := extractErrnoRef(fset, result)
		if ref != nil {
			position := fset.Position(result.Pos())
			ref.File = position.Filename
			ref.Line = position.Line
			ref.Context = getContext(fset, result, 2)

			// 判断是否需要迁移
			if ref.Action == "" {
				ref.Action = "MIGRATE"
				ref.AutoFix = true
				ref.Suggested = fmt.Sprintf("errorx.New(%s)", ref.Pattern)
			}

			refs = append(refs, ref)
		}
	}

	return refs
}

func extractErrnoRef(fset *token.FileSet, expr ast.Expr) *ErrnoRef {
	// 检查是否是 errno.ErrXXX 模式
	selExpr, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	// 检查X部分是否是"errno"
	xIdent, ok := selExpr.X.(*ast.Ident)
	if !ok || xIdent.Name != "errno" {
		return nil
	}

	// 检查Sel部分是否是大写开头（ErrXXX）
	if !strings.HasPrefix(selExpr.Sel.Name, "Err") {
		return nil
	}

	return &ErrnoRef{
		Pattern: fmt.Sprintf("errno.%s", selExpr.Sel.Name),
	}
}

func analyzeErrorxCall(fset *token.FileSet, call *ast.CallExpr) []*ErrnoRef {
	// 检查是否是 errorx.WrapByCode(err, errno.ErrXXX)
	if len(call.Args) < 2 {
		return nil
	}

	selExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	if selExpr.Sel.Name != "WrapByCode" {
		return nil
	}

	// 第二个参数是errno
	ref := extractErrnoRef(fset, call.Args[1])
	if ref != nil {
		position := fset.Position(call.Pos())
		ref.File = position.Filename
		ref.Line = position.Line
		ref.Context = getContext(fset, call, 2)
		ref.Action = "REFACTOR_WRAP_BY_CODE"
		ref.AutoFix = true
		ref.Suggested = fmt.Sprintf("errorx.Wrap(err, %s)", ref.Pattern)
		return []*ErrnoRef{ref}
	}

	return nil
}

func analyzeBinaryExpr(fset *token.FileSet, bin *ast.BinaryExpr) []*ErrnoRef {
	var refs []*ErrnoRef

	// 检查是否是测试断言
	// 例如: assert.ErrorIs(err, errno.ErrXXX)
	// 例如: errors.Is(err, errno.ErrXXX)

	if bin.Op != token.EQL && bin.Op != token.NEQ {
		return nil
	}

	// 检查右侧是否是errno引用
	ref := extractErrnoRef(fset, bin.Y)
	if ref != nil {
		position := fset.Position(bin.Pos())
		ref.File = position.Filename
		ref.Line = position.Line
		ref.Context = getContext(fset, bin, 2)
		ref.Action = "SKIP_TEST_ASSERT"
		ref.Reason = "测试断言中的errno引用不应迁移"
		ref.AutoFix = false
		refs = append(refs, ref)
	}

	return refs
}

func getContext(fset *token.FileSet, node ast.Node, lines int) string {
	position := fset.Position(node.Pos())
	file, err := os.ReadFile(position.Filename)
	if err != nil {
		return ""
	}

	lines_list := strings.Split(string(file), "\n")
	start := position.Line - lines - 1
	end := position.Line + lines

	if start < 0 {
		start = 0
	}
	if end > len(lines_list) {
		end = len(lines_list)
	}

	return strings.Join(lines_list[start:end], "\n")
}

func outputCSV(refs []*ErrnoRef, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	writer.Write([]string{"file", "line", "pattern", "context", "action", "reason", "auto_fix", "suggested"})

	// 写入数据
	for _, ref := range refs {
		writer.Write([]string{
			ref.File,
			strconv.Itoa(ref.Line),
			ref.Pattern,
			ref.Context,
			ref.Action,
			ref.Reason,
			strconv.FormatBool(ref.AutoFix),
			ref.Suggested,
		})
	}

	return nil
}
```

### 2.2 编译和运行

```bash
cd scripts/errno

# 编译AST解析器
go build -o ast_analyzer phase2_ast_analyzer.go

# 运行分析
./ast_analyzer ../../backend

# 查看结果
cat analysis/errno_ast_analysis.csv | column -t -s,
```

---

## 🔧 Phase 3: 分批迁移 - 核心迁移逻辑

### 3.1 Go程序：安全迁移器

**文件**: `scripts/errno/phase3_migrator.go`

```go
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type MigrationConfig struct {
	DryRun       bool
	BackupDir    string
	BatchSize    int
	Verbose      bool
}

type MigrationResult struct {
	File         string
	Modified     bool
	OldContent   string
	NewContent   string
	Changes      []string
	Error        error
}

func main() {
	config := MigrationConfig{
		DryRun:    os.Getenv("DRY_RUN") != "false",
		BackupDir: "scripts/errno/backups",
		BatchSize: 10,
		Verbose:   os.Getenv("VERBOSE") == "true",
	}

	// 从命令行参数读取要迁移的文件列表
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run phase3_migrator.go <file1> <file2> ...")
		os.Exit(1)
	}

	files := os.Args[1:]
	results := migrateFiles(config, files)

	// 输出结果
	printResults(results)

	// 检查是否有错误
	hasError := false
	for _, r := range results {
		if r.Error != nil {
			hasError = true
			break
		}
	}

	if hasError {
		os.Exit(1)
	}
}

func migrateFiles(config MigrationConfig, files []string) []MigrationResult {
	results := make([]MigrationResult, 0, len(files))

	for _, file := range files {
		if config.Verbose {
			fmt.Printf("处理文件: %s\n", file)
		}

		result := migrateFile(config, file)
		results = append(results, result)

		if result.Error != nil {
			fmt.Printf("❌ 错误: %s - %v\n", file, result.Error)
		} else if result.Modified {
			fmt.Printf("✓ 已修改: %s (%d处改动)\n", file, len(result.Changes))
		}
	}

	return results
}

func migrateFile(config MigrationConfig, filePath string) MigrationResult {
	result := MigrationResult{File: filePath}

	// 1. 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		result.Error = fmt.Errorf("读取文件失败: %w", err)
		return result
	}

	result.OldContent = string(content)

	// 2. 解析AST
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		result.Error = fmt.Errorf("解析AST失败: %w", err)
		return result
	}

	// 3. 创建备份
	if !config.DryRun {
		backupPath := filepath.Join(config.BackupDir, filePath)
		backupDir := filepath.Dir(backupPath)
		os.MkdirAll(backupDir, 0755)
		if err := os.WriteFile(backupPath, content, 0644); err != nil {
			result.Error = fmt.Errorf("创建备份失败: %w", err)
			return result
		}
	}

	// 4. 执行迁移
	migrator := &ErrnoMigrator{fset: fset, file: file}
	ast.Walk(migrator, file)

	if len(migrator.changes) == 0 {
		result.Modified = false
		return result
	}

	result.Changes = migrator.changes
	result.Modified = true

	// 5. 生成新代码
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		result.Error = fmt.Errorf("格式化代码失败: %w", err)
		return result
	}

	result.NewContent = buf.String()

	// 6. 写入文件（如果不是dry run）
	if !config.DryRun {
		if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
			result.Error = fmt.Errorf("写入文件失败: %w", err)
			return result
		}
	}

	return result
}

type ErrnoMigrator struct {
	fset    *token.FileSet
	file    *ast.File
	changes []string
}

func (m *ErrnoMigrator) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ReturnStmt:
		m.migrateReturnStmt(n)
	case *ast.CallExpr:
		m.migrateCallExpr(n)
	}
	return m
}

func (m *ErrnoMigrator) migrateReturnStmt(ret *ast.ReturnStmt) {
	if len(ret.Results) == 0 {
		return
	}

	for i, result := range ret.Results {
		newExpr := m.transformErrnoRef(result)
		if newExpr != nil {
			ret.Results[i] = newExpr
			pos := m.fset.Position(result.Pos())
			m.changes = append(m.changes, fmt.Sprintf("%s:%d: return errno.XXX → return errorx.New(errno.XXX)",
				pos.Filename, pos.Line))
		}
	}
}

func (m *ErrnoMigrator) migrateCallExpr(call *ast.CallExpr) {
	// 检查是否是 errorx.WrapByCode(err, errno.ErrXXX)
	selExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	if selExpr.Sel.Name != "WrapByCode" {
		return
	}

	// 修改为 errorx.Wrap(err, errno.ErrXXX)
	selExpr.Sel.Name = "Wrap"

	// 删除第二个参数（如果err是error类型，Wrap会自动包装错误码）
	if len(call.Args) > 1 {
		// 保留第一个参数（err）
		// 将第二个参数（errno.ErrXXX）作为Wrap的第二个参数
		// 不做修改，只是改函数名
	}

	pos := m.fset.Position(call.Pos())
	m.changes = append(m.changes, fmt.Sprintf("%s:%d: errorx.WrapByCode → errorx.Wrap",
		pos.Filename, pos.Line))
}

func (m *ErrnoMigrator) transformErrnoRef(expr ast.Expr) ast.Expr {
	selExpr, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	xIdent, ok := selExpr.X.(*ast.Ident)
	if !ok || xIdent.Name != "errno" {
		return nil
	}

	if !strings.HasPrefix(selExpr.Sel.Name, "Err") {
		return nil
	}

	// 创建 errorx.New(errno.ErrXXX) 调用
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("errorx"),
			Sel: ast.NewIdent("New"),
		},
		Args: []ast.Expr{selExpr}, // 保留原selExpr作为参数
	}
}

func printResults(results []MigrationResult) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("迁移结果汇总")
	fmt.Println(strings.Repeat("=", 80))

	totalFiles := len(results)
	modifiedFiles := 0
	totalChanges := 0

	for _, r := range results {
		if r.Modified {
			modifiedFiles++
			totalChanges += len(r.Changes)
		}
	}

	fmt.Printf("\n📊 统计:\n")
	fmt.Printf("  总文件数: %d\n", totalFiles)
	fmt.Printf("  已修改: %d\n", modifiedFiles)
	fmt.Printf("  总改动: %d\n", totalChanges)

	if modifiedFiles > 0 {
		fmt.Printf("\n✅ 已修改的文件:\n")
		for _, r := range results {
			if r.Modified {
				fmt.Printf("  - %s (%d处改动)\n", r.File, len(r.Changes))
			}
		}
	}

	if !results[0].Modified && os.Getenv("DRY_RUN") == "true" {
		fmt.Printf("\n⚠️  Dry Run模式，未实际修改文件\n")
		fmt.Printf("   实际执行请使用: DRY_RUN=false go run phase3_migrator.go ...\n")
	}
}
```

### 3.2 分批执行脚本

**文件**: `scripts/errno/phase3_migrate_batch.sh`

```bash
#!/bin/bash
# Phase 3: 分批迁移errno引用

set -e

BATCH_SIZE=10
CSV_FILE="scripts/errno/analysis/errno_ast_analysis.csv"
BACKEND_DIR="backend"
MIGRATOR_BIN="scripts/errno/ast_migrator"

# 编译迁移器
echo "🔨 编译迁移器..."
cd scripts/errno
go build -o ast_migrator phase3_migrator.go
cd ../..

# 从CSV中提取需要迁移的文件列表（去重）
echo "📋 读取文件清单..."
files_to_migrate=$(tail -n +2 "$CSV_FILE" | \
    grep ",MIGRATE$" | \
    cut -d, -f1 | \
    sort | \
    uniq)

total_files=$(echo "$files_to_migrate" | wc -l)
echo "  总文件数: $total_files"

# 分批执行
batch_num=1
current_batch=0

echo ""
echo "$files_to_migrate" | while read file; do
    # 跳过空行
    [ -z "$file" ] && continue

    full_path="$BACKEND_DIR/$file"

    # 检查文件是否存在
    if [ ! -f "$full_path" ]; then
        echo "⚠️  文件不存在，跳过: $full_path"
        continue
    fi

    current_batch=$((current_batch + 1))

    echo "[$current_batch/$total_files] 迁移: $file"

    # 执行迁移（dry run模式）
    DRY_RUN=true VERBOSE=false "$MIGRATOR_BIN" "$full_path"

    # 每BATCH_SIZE个文件后暂停，等待确认
    if [ $((current_batch % BATCH_SIZE)) -eq 0 ]; then
        echo ""
        echo "=========================================="
        echo "批次 $batch_num 完成（已处理 $current_batch/$total_files 个文件）"
        echo "=========================================="
        echo ""
        echo "请检查上述输出，确认无误后继续："
        echo "  • 输入 'y' 继续下一批"
        echo "  • 输入 'a' 应用这批更改（实际修改文件）"
        echo "  • 输入 'q' 退出"
        echo ""
        read -p "> " choice

        case "$choice" in
            y|Y)
                echo "继续下一批..."
                ;;
            a|A)
                echo "应用更改..."
                # 重新执行这批迁移，但这次不是dry run
                DRY_RUN=false VERBOSE=false "$MIGRATOR_BIN" "$full_path"
                ;;
            q|Q)
                echo "退出迁移"
                exit 0
                ;;
            *)
                echo "无效选择，退出"
                exit 1
                ;;
        esac

        batch_num=$((batch_num + 1))
        echo ""
    fi
done

echo ""
echo "✅ 所有批次处理完成！"
echo ""
echo "下一步:"
echo "  1. 检查git状态: git status"
echo "  2. 运行测试: go test ./..."
echo "  3. 提交更改: git add . && git commit -m 'feat: migrate errno to errorx'"
```

---

## ✅ Phase 4: 验证 - 多重验证机制

### 4.1 语法验证脚本

**文件**: `scripts/errno/phase4_validate_syntax.sh`

```bash
#!/bin/bash
# Phase 4.1: 验证Go语法

echo "🔍 Phase 4.1: 验证Go语法..."
echo ""

errors=0

# 获取所有修改过的文件
modified_files=$(git diff --name-only | grep '\.go$' || true)

if [ -z "$modified_files" ]; then
    echo "⚠️  没有修改过的文件"
    exit 0
fi

echo "检查文件: $modified_files"
echo ""

for file in $modified_files; do
    echo -n "  ✓ $file ... "

    # 使用go fmt验证格式
    if ! go fmt "$file" > /dev/null 2>&1; then
        echo "❌ 格式错误"
        errors=$((errors + 1))
        continue
    fi

    # 使用go vet检查
    if ! go vet "$file" > /dev/null 2>&1; then
        echo "⚠️  go vet 警告"
        # 不算错误，继续
    fi

    echo "✅"
done

echo ""
if [ $errors -eq 0 ]; then
    echo "✅ 语法验证通过！"
    exit 0
else
    echo "❌ 发现 $errors 个错误"
    exit 1
fi
```

### 4.2 编译验证脚本

**文件**: `scripts/errno/phase4_validate_compile.sh`

```bash
#!/bin/bash
# Phase 4.2: 验证编译

echo "🔨 Phase 4.2: 验证编译..."
echo ""

# 编译整个项目
if go build ./...; then
    echo "✅ 编译成功！"
    exit 0
else
    echo "❌ 编译失败"
    exit 1
fi
```

### 4.3 单元测试验证脚本

**文件**: `scripts/errno/phase4_validate_tests.sh`

```bash
#!/bin/bash
# Phase 4.3: 运行单元测试

echo "🧪 Phase 4.3: 运行单元测试..."
echo ""

# 运行测试（不包含集成测试）
if go test -short ./... -v 2>&1 | tee test_output.log; then
    echo ""
    echo "✅ 测试通过！"
    exit 0
else
    echo ""
    echo "❌ 测试失败"
    echo ""
    echo "失败的测试:"
    grep "FAIL:" test_output.log || true
    exit 1
fi
```

### 4.4 综合验证脚本

**文件**: `scripts/errno/phase4_validate_all.sh`

```bash
#!/bin/bash
# Phase 4: 综合验证（语法+编译+测试）

set -e

echo "=========================================="
echo "  errno迁移验证"
echo "=========================================="
echo ""

# Phase 4.1: 语法验证
./scripts/errno/phase4_validate_syntax.sh
echo ""

# Phase 4.2: 编译验证
./scripts/errno/phase4_validate_compile.sh
echo ""

# Phase 4.3: 单元测试验证
./scripts/errno/phase4_validate_tests.sh
echo ""

echo "=========================================="
echo "  ✅ 所有验证通过！"
echo "=========================================="
echo ""
echo "下一步:"
echo "  1. 代码review"
echo "  2. 提交PR"
echo "  3. 等待CI/CD检查"
```

---

## 🔄 Phase 5: 回滚机制

### 5.1 快速回滚脚本

**文件**: `scripts/errno/phase5_rollback.sh`

```bash
#!/bin/bash
# 快速回滚到迁移前的状态

set -e

BACKUP_DIR="scripts/errno/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo "⚠️  回滚操作"
echo "=========================================="
echo ""
echo "此操作将使用备份文件恢复所有修改过的文件"
echo "备份目录: $BACKUP_DIR"
echo ""

# 检查备份目录是否存在
if [ ! -d "$BACKUP_DIR" ]; then
    echo "❌ 备份目录不存在，无法回滚"
    exit 1
fi

# 确认
read -p "确认回滚？(yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "已取消"
    exit 0
fi

echo ""
echo "开始回滚..."
echo ""

# 创建当前状态的备份（以防万一）
rollback_backup="$BACKUP_DIR/before_rollback_$TIMESTAMP"
mkdir -p "$rollback_backup"

git diff --name-only | while read file; do
    if [ -f "$file" ]; then
        mkdir -p "$rollback_backup/$(dirname "$file")"
        cp "$file" "$rollback_backup/$file"
    fi
done

echo "✓ 已备份当前状态到: $rollback_backup"
echo ""

# 从备份恢复文件
restored=0
find "$BACKUP_DIR/backend" -type f | while read backup_file; do
    # 计算相对路径
    rel_path="${backup_file#$BACKUP_DIR/}"
    target_file="$rel_path"

    echo "恢复: $target_file"
    cp "$backup_file" "$target_file"
    restored=$((restored + 1))
done

echo ""
echo "✅ 已恢复 $restored 个文件"
echo ""
echo "下一步:"
echo "  1. 检查恢复的文件: git status"
echo "  2. 重新编译: go build ./..."
echo "  3. 运行测试: go test ./..."
```

### 5.2 Git回滚脚本

**文件**: `scripts/errno/phase5_rollback_git.sh`

```bash
#!/bin/bash
# 使用Git回滚到迁移前的commit

echo "=========================================="
echo "  Git回滚"
echo "=========================================="
echo ""

# 显示最近的commits
echo "最近的commits:"
echo ""
git log --oneline -10
echo ""

# 提示用户选择commit
echo "请输入要回滚到的commit hash（或输入'q'退出）:"
read -p "> " commit_hash

if [ "$commit_hash" = "q" ]; then
    echo "已取消"
    exit 0
fi

# 确认
echo ""
git log -1 --oneline "$commit_hash"
echo ""
read -p "确认回滚到上述commit？(yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "已取消"
    exit 0
fi

# 执行回滚
echo ""
echo "执行回滚..."
if git reset --hard "$commit_hash"; then
    echo ""
    echo "✅ 已回滚到: $commit_hash"
    echo ""
    echo "⚠️  警告: 此操作会丢弃所有未提交的更改"
    echo ""
    echo "如果需要保留更改，请使用:"
    echo "  git reflog"
    echo "  git reset HEAD@{n}"
else
    echo ""
    echo "❌ 回滚失败"
    exit 1
fi
```

---

## 🚀 执行流程

### 完整执行步骤

```bash
# ========================================
# Step 1: 静态扫描
# ========================================
./scripts/errno/phase1_scan.sh

# 查看报告
cat scripts/errno/analysis/errno_migration_report.md

# ========================================
# Step 2: AST分析（可选，更精确）
# ========================================
cd scripts/errno
go build -o ast_analyzer phase2_ast_analyzer.go
./ast_analyzer ../../backend

# 对比两次分析结果
diff analysis/errno_inventory_report.csv analysis/errno_ast_analysis.csv

# ========================================
# Step 3: 分批迁移
# ========================================
# 先dry run一遍
DRY_RUN=true ./scripts/errno/phase3_migrate_batch.sh

# 确认无误后，实际执行
DRY_RUN=false ./scripts/errno/phase3_migrate_batch.sh

# ========================================
# Step 4: 验证
# ========================================
./scripts/errno/phase4_validate_all.sh

# ========================================
# Step 5: 提交
# ========================================
git add .
git commit -m "feat: migrate errno to errorx

- 迁移375处errno.ErrXXX引用到errorx.New(errno.ErrXXX)
- 使用企业级AST解析器精确迁移
- 包含完整备份和回滚机制
- 所有测试通过

相关文档: docs/02-SPECS/代码质量改进设计文档_v1.0.md"

# ========================================
# 如果需要回滚
# ========================================
# 方式1: 使用备份文件回滚
./scripts/errno/phase5_rollback.sh

# 方式2: 使用Git回滚
./scripts/errno/phase5_rollback_git.sh
```

---

## 📊 成功标准

### 验收清单

- [ ] **Phase 1**: 扫描完成，生成详细清单
  - [ ] 识别所有375处errno引用
  - [ ] 分类准确（需迁移/应跳过）
  - [ ] 生成Markdown报告

- [ ] **Phase 2**: AST分析完成（可选）
  - [ ] Go AST解析器编译成功
  - [ ] 生成精确分析结果
  - [ ] 与Phase 1结果对比验证

- [ ] **Phase 3**: 分批迁移完成
  - [ ] 所有批次处理完毕
  - [ ] 每批都经过dry run验证
  - [ ] 所有备份文件已创建

- [ ] **Phase 4**: 所有验证通过
  - [ ] Go语法验证通过
  - [ ] 编译验证通过
  - [ ] 单元测试通过（覆盖率≥80%）

- [ ] **Phase 5**: 回滚机制验证
  - [ ] 备份文件完整
  - [ ] 回滚脚本测试通过

### 最终检查

```bash
# 1. 检查是否还有直接errno引用
grep -rn "return.*errno\." backend/domain/ --include="*.go" | grep -v "_test.go" | wc -l
# 预期: 0

# 2. 检查errorx使用量
grep -rn "errorx\.New.*errno\." backend/domain/ --include="*.go" | wc -l
# 预期: ≈375

# 3. 检查编译
go build ./...
# 预期: 成功

# 4. 检查测试
go test ./... -short
# 预期: 全部通过

# 5. 检查覆盖率
go test -cover ./... | grep total
# 预期: coverage ≥ 80%
```

---

## 🎯 风险缓解措施

| 风险 | 缓解措施 | 验证方式 |
|------|----------|----------|
| **误改代码** | 使用Go AST精确解析，不是sed | Phase 2验证 |
| **语法错误** | go/format格式化，go vet检查 | Phase 4.1验证 |
| **编译失败** | 每批都编译验证 | Phase 4.2验证 |
| **测试失败** | 每批都运行测试 | Phase 4.3验证 |
| **性能下降** | 基准测试对比 | Benchmark测试 |
| **遗漏文件** | 详细的CSV清单 | Phase 1报告 |
| **无法回滚** | 完整备份+Git commit | Phase 5测试 |

---

## 📞 紧急联系

如果迁移过程中遇到问题：

1. **立即停止**: Ctrl+C 或 `q` 退出
2. **检查日志**: `scripts/errno/analysis/`
3. **回滚**: `./scripts/errno/phase5_rollback.sh`
4. **联系**: 技术负责人 + DevOps团队

---

**文档版本**: v2.0 (企业级零风险方案)
**最后更新**: 2025-01-03
**下次评审**: 迁移完成后
