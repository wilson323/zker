# errno迁移脚本使用指南

> **企业级零风险errno迁移方案** - ZKER项目专用
>
> **目标**: 375处`errno.ErrXXX` → `errorx.New(errno.ErrXXX)`，0异常，0回滚

---

## 📋 目录

- [快速开始](#快速开始)
- [执行流程](#执行流程)
- [详细说明](#详细说明)
- [应急处理](#应急处理)
- [FAQ](#faq)

---

## 🚀 快速开始

### 一键执行（推荐）

```bash
# Step 1: 扫描分析
./scripts/errno/phase1_scan.sh

# Step 2: 查看报告
cat scripts/errno/analysis/errno_inventory_report.csv | column -t -s, -o '|'

# Step 3: 分批迁移
./scripts/errno/phase3_migrate_batch.sh

# Step 4: 综合验证
./scripts/errno/phase4_validate_all.sh

# Step 5: 提交代码
git add . && git commit -m "feat: migrate errno to errorx"
```

---

## 📊 执行流程

### Phase 1: 静态扫描

**脚本**: `phase1_scan.sh`

**功能**:
- 扫描所有Go文件中的errno引用
- 生成详细的CSV清单
- 分类：需迁移 / 应跳过

**输出**:
```
scripts/errno/analysis/
├── errno_inventory_report.csv    # 详细清单
└── errno_summary.txt             # 统计汇总
```

**执行**:
```bash
./scripts/errno/phase1_scan.sh
```

**查看结果**:
```bash
# 查看CSV报告（格式化）
cat scripts/errno/analysis/errno_inventory_report.csv | column -t -s, -o '|'

# 查看统计
cat scripts/errno/analysis/errno_summary.txt
```

**CSV字段说明**:
| 字段 | 说明 |
|------|------|
| file_path | 文件相对路径 |
| line_number | 行号 |
| context | 上下文（前后2行） |
| pattern | errno模式（如errno.ErrBotNotFound） |
| recommended_action | 推荐操作（MIGRATE/SKIKE_XXX） |
| reason | 原因说明 |

---

### Phase 2: AST精确分析（可选）

**脚本**: `phase2_ast_analyzer.go`

**功能**:
- 使用Go AST解析器精确分析
- 区分return语句、函数调用、测试断言
- 生成更精确的迁移建议

**执行**:
```bash
cd scripts/errno
go build -o ast_analyzer phase2_ast_analyzer.go
./ast_analyzer ../../backend
```

**输出**:
```
scripts/errno/analysis/errno_ast_analysis.csv
```

**对比分析**:
```bash
# 对比Phase 1和Phase 2的分析结果
diff analysis/errno_inventory_report.csv analysis/errno_ast_analysis.csv
```

---

### Phase 3: 分批迁移

**脚本**: `phase3_migrate_batch.sh`

**功能**:
- 分批迁移（每批10个文件）
- 每批执行前dry run验证
- 交互式确认，出错可暂停

**执行**:
```bash
# 先dry run一遍（不实际修改文件）
./scripts/errno/phase3_migrate_batch.sh

# 确认无误后，实际执行
# 编辑脚本，将 DRY_RUN=true 改为 DRY_RUN=false
./scripts/errno/phase3_migrate_batch.sh
```

**流程**:
1. 编译迁移器（Go程序）
2. 从CSV读取文件清单
3. 分批处理（每批10个文件）
4. 每批执行前dry run验证
5. 等待用户确认（y继续/a应用/q退出）
6. 创建备份文件
7. 执行迁移
8. 进入下一批

**交互示例**:
```
[10/375] 迁移: domain/agent/service/bot_service.go
✓ Dry run完成: domain/agent/service/bot_service.go (3处改动)

==========================================
批次 1 完成（已处理 10/375 个文件）
==========================================

请检查上述输出，确认无误后继续：
  • 输入 'y' 继续下一批
  • 输入 'a' 应用这批更改（实际修改文件）
  • 输入 'q' 退出
```

---

### Phase 4: 综合验证

**脚本**: `phase4_validate_all.sh`

**功能**:
- Go语法验证（go fmt + go vet）
- 编译验证（go build）
- 单元测试验证（go test -short）
- errno使用验证（检查覆盖率）

**执行**:
```bash
./scripts/errno/phase4_validate_all.sh
```

**输出**:
```
==========================================
  errno迁移 - Phase 4: 综合验证
==========================================

==========================================
  Phase 4.1: Go语法验证
==========================================

  [1] bot_service.go ... ✓ OK
  [2] agent_service.go ... ✓ OK
  ...

✅ 语法验证通过（检查了 150 个文件）

==========================================
  Phase 4.2: 编译验证
==========================================

编译整个项目...
✅ 编译成功！

==========================================
  Phase 4.3: 单元测试验证
==========================================

运行单元测试...
✅ 所有测试通过！

==========================================
  Phase 4.4: errno使用验证
==========================================

  1. 直接errno引用: 0 (应为0)
  2. errorx.New使用: 375
  3. errorx.Wrap使用: 42
  4. 错误处理覆盖率: 98% (目标≥95%)

✅ 错误处理覆盖率达标

==========================================
验证总结
==========================================

总耗时: 125秒

检查项结果:
  语法验证:   ✅ 通过
  编译验证:   ✅ 通过
  测试验证:   ✅ 通过
  errno使用:  ✅ 已检查

==========================================
✅ 所有验证通过！
```

---

### Phase 5: 回滚（应急）

**脚本**: `phase5_rollback.sh`

**功能**:
- 从备份文件恢复
- 创建回滚前的备份
- 验证恢复结果

**执行**:
```bash
./scripts/errno/phase5_rollback.sh
```

**流程**:
1. 检查备份目录是否存在
2. 备份当前状态（双重保险）
3. 从备份恢复文件
4. 验证恢复结果（语法+编译）
5. 显示后续操作指引

**示例输出**:
```
==========================================
  errno迁移回滚
==========================================

⚠️  警告：此操作将覆盖当前文件！

此操作将:
  1. 备份当前状态到: scripts/errno/backups/before_rollback_20250103_143052
  2. 从备份目录恢复文件: scripts/errno/backups
  3. 验证恢复结果

备份目录: scripts/errno/backups

✓ 找到 150 个备份文件

请确认回滚操作:
  • 输入 'yes' 继续回滚
  • 输入 'q'   退出

> yes

开始回滚...

✓ 已备份 150 个文件到: scripts/errno/backups/before_rollback_20250103_143052

从备份恢复文件...
✓ 恢复: backend/domain/agent/service/bot_service.go
✓ 恢复: backend/domain/agent/service/agent_service.go
...

✅ 回滚成功！
```

---

## 📚 详细说明

### 文件结构

```
scripts/errno/
├── README.md                           # 本文件
├── enterprise_errno_migration_plan.md  # 完整技术方案
│
├── phase1_scan.sh                      # Phase 1: 静态扫描
├── phase1_generate_report.sh           # Phase 1: 生成Markdown报告
│
├── phase2_ast_analyzer.go              # Phase 2: AST解析器
│
├── phase3_migrator.go                  # Phase 3: 迁移器核心
├── phase3_migrate_batch.sh             # Phase 3: 分批执行脚本
│
├── phase4_validate_syntax.sh           # Phase 4.1: 语法验证
├── phase4_validate_compile.sh          # Phase 4.2: 编译验证
├── phase4_validate_tests.sh            # Phase 4.3: 测试验证
├── phase4_validate_all.sh              # Phase 4: 综合验证
│
├── phase5_rollback.sh                  # Phase 5: 文件回滚
└── phase5_rollback_git.sh              # Phase 5: Git回滚
```

### 输出目录结构

```
scripts/errno/
├── analysis/                           # 分析结果
│   ├── errno_inventory_report.csv     # Phase 1详细清单
│   ├── errno_summary.txt              # Phase 1统计汇总
│   ├── errno_migration_report.md      # Phase 1 Markdown报告
│   └── errno_ast_analysis.csv         # Phase 2 AST分析结果
│
├── backups/                            # 备份文件
│   └── backend/                       # 按原目录结构备份
│       ├── domain/...
│       ├── application/...
│       └── api/...
│
└── logs/                               # 验证日志
    ├── syntax_check.log               # 语法检查日志
    ├── compile_check.log              # 编译检查日志
    └── test_check.log                # 测试检查日志
```

---

## 🚨 应急处理

### 场景1: 迁移过程中出错

**症状**: Phase 3执行过程中出现错误

**处理**:
1. 立即停止（Ctrl+C或输入`q`退出）
2. 查看错误日志
3. 使用phase5回滚

```bash
# 停止迁移（如果还在运行）
Ctrl+C

# 回滚已迁移的文件
./scripts/errno/phase5_rollback.sh
```

### 场景2: 验证失败

**症状**: Phase 4验证不通过

**处理**:
1. 查看具体失败的检查项
2. 阅读详细日志
3. 修复问题或回滚

```bash
# 查看验证日志
cat scripts/errno/logs/syntax_check.log
cat scripts/errno/logs/compile_check.log
cat scripts/errno/logs/test_check.log

# 如果无法修复，回滚
./scripts/errno/phase5_rollback.sh
```

### 场景3: 需要Git回滚

**症状**: 已经提交了PR，需要回滚

**处理**:
1. 使用Git回滚到之前的commit

```bash
# 查看最近的commits
git log --oneline -10

# 回滚到迁移前的commit
git reset --hard <commit-hash>

# 或者创建新的commit回滚
git revert HEAD
```

---

## ❓ FAQ

### Q1: 迁移会修改哪些文件？

**A**: 只修改`backend/domain/`目录下的Go文件（排除测试文件）。

具体范围：
- `backend/domain/**/*.go`
- 排除: `*_test.go`, `gen.go`, `gen/`目录

### Q2: 迁移是自动的吗？

**A**: 半自动。
- **自动部分**: 扫描、分析、代码替换
- **人工部分**: 分批确认、review、提交

### Q3: 如何确保迁移不出错？

**A**: 多重保障：
1. **精确解析**: 使用Go AST，不是简单sed
2. **分批执行**: 每批10个文件，逐步验证
3. **Dry Run**: 先模拟，确认后再实际执行
4. **完整备份**: 每个文件都有备份
5. **多重验证**: 语法+编译+测试
6. **快速回滚**: 出问题可立即回滚

### Q4: 迁移需要多长时间？

**A**: 预计2-4小时
- Phase 1: 10分钟（扫描）
- Phase 2: 20分钟（AST分析，可选）
- Phase 3: 2-3小时（分批迁移，含人工确认）
- Phase 4: 30分钟（验证）
- Phase 5: 应急备用

### Q5: 迁移后代码会有什么变化？

**A**: 变化示例：

**Before**:
```go
if err != nil {
    return nil, errno.ErrBotNotFound
}
```

**After**:
```go
if err != nil {
    return nil, errorx.New(errno.ErrBotNotFound)
}
```

### Q6: 如果迁移失败了怎么办？

**A**: 三种回滚方式：
1. **文件回滚**: `./scripts/errno/phase5_rollback.sh`
2. **Git回滚**: `git reset --hard <commit>`
3. **手动回滚**: 从备份目录手动复制文件

### Q7: 迁移会影响性能吗？

**A**: 影响极小（<1%）。
- `errorx.New`只是包装，不增加额外逻辑
- 运行时性能几乎无差异
- 带来的好处：更好的错误追踪、i18n支持

### Q8: 迁移后需要修改测试代码吗？

**A**: 不需要。
- 测试代码中的`errno.ErrXXX`保持不变
- 断言语句（如`assert.ErrorIs(err, errno.ErrXXX)`）不被修改
- 迁移脚本会自动跳过测试文件

---

## 📞 支持

如有问题，请联系：
- **技术负责人**: [TODO: 填写负责人]
- **DevOps团队**: [TODO: 填写DevOps联系方式]
- **紧急联系**: [TODO: 填写紧急联系方式]

---

## 📄 相关文档

- **完整技术方案**: [enterprise_errno_migration_plan.md](./enterprise_errno_migration_plan.md)
- **代码质量改进文档**: [../../docs/02-SPECS/代码质量改进设计文档_v1.0.md](../../docs/02-SPECS/代码质量改进设计文档_v1.0.md)
- **P0执行计划**: [../../docs/04-IMPLEMENTATION/P0任务详细执行计划_v1.0.md](../../docs/04-IMPLEMENTATION/P0任务详细执行计划_v1.0.md)

---

**文档版本**: v1.0
**最后更新**: 2025-01-03
**维护者**: ZKER架构团队
