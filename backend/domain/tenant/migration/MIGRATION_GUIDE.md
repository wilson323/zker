# tenant_id 迁移操作指南

**文档版本**: v1.0
**创建日期**: 2025-01-01
**负责人**: 研发B（后端工程师）
**适用人群**: 运维工程师、DBA、后端开发

---

## 📋 迁移概述

### 目标
将 Coze Studio 从单租户架构迁移到企业级多租户 SaaS 架构，为 10+ 核心业务表添加 `tenant_id` 字段，实现租户数据隔离。

### 涉及表（10张）
- bots, bot_configs, single_agent_draft, published_bots
- conversations, messages
- knowledge_bases, knowledge_chunks
- workflows, workflow_executions

### 迁移策略
**渐进式双写迁移**（5个阶段）：
1. DDL变更（1人天）
2. 数据预迁移（2人天）
3. 双写验证（2人天）
4. 灰度切换（2人天）
5. 清理收尾（1人天）

### 关键约束
- ✅ 零停机迁移
- ✅ 数据一致性 ≥ 99.9%
- ✅ 批处理：1000条/批
- ✅ 事务安全：失败自动回滚

---

## 🚀 阶段 1: DDL 变更（1人天）

### 步骤 1.1: 备份数据库

**⚠️ 必须执行**：

```bash
# 创建备份目录
mkdir -p /backups/migration/$(date +%Y%m%d)

# 全量备份
mysqldump -h127.0.0.1 -uroot -p \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  coze_production > /backups/migration/$(date +%Y%m%d)/backup_before_ddl.sql

# 验证备份文件
ls -lh /backups/migration/$(date +%Y%m%d)/backup_before_ddl.sql
```

---

### 步骤 1.2: 执行 DDL 脚本

**执行时间**：凌晨 2:00-4:00（低峰期）

```bash
# 进入项目目录
cd /path/to/coze-studio

# 执行DDL（约10-30分钟）
mysql -h127.0.0.1 -uroot -p \
  < backend/domain/tenant/migration/ddl/001_add_tenant_id.sql

# 监控执行进度
# 在另一个终端窗口执行：
watch -n 5 "mysql -h127.0.0.1 -uroot -p -e \"SELECT COUNT(*) as count FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND COLUMN_NAME = 'tenant_id';\""
```

---

### 步骤 1.3: 验证 DDL 执行结果

```bash
# 验证字段已添加
mysql -h127.0.0.1 -uroot -p -e "
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_TYPE,
    IS_NULLABLE,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
ORDER BY TABLE_NAME;
"

# 验证索引已创建
mysql -h127.0.0.1 -uroot -p -e "
SHOW INDEX FROM bots WHERE Key_name = 'idx_tenant_id';
SHOW INDEX FROM bots WHERE Key_name = 'uk_tenant_bot';
"

# 预期结果：
# - 10张表都有 tenant_id 字段
# - 所有表都有 idx_tenant_id 索引
# - 部分表有唯一约束索引
```

---

### 步骤 1.4: 业务验证

```bash
# 1. 健康检查
curl http://api-server:8001/health

# 2. 核心功能测试（创建Bot）
curl -X POST http://api-server:8001/api/v1/bots \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test_bot_after_ddl",
    "description": "测试DDL后业务功能"
  }'

# 3. 查询功能测试
curl http://api-server:8001/api/v1/bots

# 预期结果：
# - 健康检查通过
# - 创建Bot成功
# - 查询Bot列表成功
# - 无业务异常
```

---

### 步骤 1.5: 回滚方案（如需回滚）

**触发条件**：
- DDL执行失败
- 业务功能异常
- 性能严重下降

```bash
# 快速回滚（5-10分钟）
mysql -h127.0.0.1 -uroot -p \
  < backend/domain/tenant/migration/ddl/rollback_001_add_tenant_id.sql

# 验证回滚结果
mysql -h127.0.0.1 -uroot -p -e "
SELECT TABLE_NAME, COLUMN_NAME
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE() AND COLUMN_NAME = 'tenant_id';
"

# 预期结果：返回空（tenant_id字段已删除）
```

---

## 🔄 阶段 2: 数据预迁移（2人天）

### 步骤 2.1: 构建迁移工具

```bash
# 进入项目目录
cd /path/to/coze-studio/backend

# 构建迁移工具
go build -o bin/tenant-migrate ./domain/tenant/migration/cmd/migrator

# 验证构建成功
./bin/tenant-migrate --help
```

---

### 步骤 2.2: 执行数据迁移（后台运行）

```bash
# 迁移 bots 表
nohup ./bin/tenant-migrate \
  --table bots \
  --batch-size 1000 \
  --concurrency 10 \
  --default-tenant default_tenant \
  > logs/migrate_bots.log 2>&1 &

# 查看进程
ps aux | grep tenant-migrate

# 查看实时日志
tail -f logs/migrate_bots.log
```

---

### 步骤 2.3: 监控迁移进度

```bash
# API查询进度
curl http://api-server:8001/api/v1/admin/migration/progress/bots

# 返回示例：
{
  "table": "bots",
  "total_records": 500000,
  "migrated_records": 250000,
  "progress_percent": 50.0,
  "status": "in_progress",
  "started_at": "2025-01-01T02:00:00Z",
  "estimated_finish_at": "2025-01-01T04:00:00Z"
}
```

---

### 步骤 2.4: 验证数据迁移

```bash
# 对比记录数
mysql -h127.0.0.1 -uroot -p -e "
SELECT
    COUNT(*) as total,
    SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) as with_tenant_id,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_tenant_id
FROM bots;
"

# 预期结果：
# total = with_tenant_id（100%迁移完成）
# null_tenant_id = 0

# 抽样验证数据正确性
mysql -h127.0.0.1 -uroot -p -e "
SELECT bot_id, tenant_id, created_by
FROM bots
WHERE tenant_id IS NOT NULL
LIMIT 10;
"
```

---

## 🔗 阶段 3: 双写验证（2人天）

### 步骤 3.1: 启用双写

```bash
# 配置双写（通过环境变量或配置文件）
export DUAL_WRITE_ENABLED=true
export DUAL_WRITE_MODE=async

# 重启服务
systemctl restart coze-api-server

# 验证双写已启用
curl http://api-server:8001/api/v1/admin/dual-write/status
```

---

### 步骤 3.2: 数据一致性校验

```bash
# 运行一致性检查工具
./bin/consistency-checker \
  --table bots \
  --sample-size 1000

# 预期结果：
# - 数据一致性 ≥ 99.9%
# - 不一致记录自动标记
```

---

### 步骤 3.3: 监控双写性能

```bash
# 查看Prometheus指标
curl http://api-server:8001/metrics | grep dual_write

# 关键指标：
# - dual_write_total: 双写总次数
# - dual_write_failed_total: 双写失败次数
# - dual_write_duration_seconds: 双写耗时
```

---

## 🚀 阶段 4: 灰度切换（2人天）

### 步骤 4.1: 第一阶段灰度（5%流量）

```bash
# 设置灰度配置
curl -X POST http://api-server:8001/api/v1/admin/gray-release/config \
  -H "Content-Type: application/json" \
  -d '{
    "stage": 1,
    "traffic_percent": 5.0,
    "enabled_tenants": ["tenant_test_1", "tenant_test_2"]
  }'

# 观察24小时
# 监控指标：
# - 错误率 < 0.1%
# - P95延迟 < 基线*1.2
# - 数据一致性 > 99.9%
```

---

### 步骤 4.2-4.4: 后续阶段灰度

```bash
# 第二阶段：20%流量
curl -X POST http://api-server:8001/api/v1/admin/gray-release/config \
  -d '{"stage": 2, "traffic_percent": 20.0}'

# 第三阶段：50%流量
curl -X POST http://api-server:8001/api/v1/admin/gray-release/config \
  -d '{"stage": 3, "traffic_percent": 50.0}'

# 第四阶段：100%流量
curl -X POST http://api-server:8001/api/v1/admin/gray-release/config \
  -d '{"stage": 4, "traffic_percent": 100.0}'
```

---

## 🧹 阶段 5: 清理收尾（1人天）

### 步骤 5.1: 移除双写代码

```bash
# 代码已清理，无需手动操作
# 代码位置：backend/domain/*/service/*.go
```

---

### 步骤 5.2: 设置字段为 NOT NULL

```sql
-- 数据迁移完成后执行
ALTER TABLE bots MODIFY COLUMN tenant_id VARCHAR(36) NOT NULL;
ALTER TABLE bot_configs MODIFY COLUMN tenant_id VARCHAR(36) NOT NULL;
-- ... 其他表类似
```

---

### 步骤 5.3: 验证最终状态

```bash
# 验证所有表都有 tenant_id 且为 NOT NULL
mysql -h127.0.0.1 -uroot -p -e "
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    IS_NULLABLE,
    COLUMN_KEY
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
    'bots', 'bot_configs', 'single_agent_draft', 'published_bots',
    'conversations', 'messages',
    'knowledge_bases', 'knowledge_chunks',
    'workflows', 'workflow_executions'
  )
ORDER BY TABLE_NAME;
"

# 预期结果：
# - 所有表的 IS_NULLABLE = NO
# - 所有表都有索引
```

---

## 📊 监控与告警

### 关键监控指标

**Prometheus 指标**：
```yaml
# 迁移进度
migration_progress_percent{table="bots"}

# 数据一致性
consistency_rate{table="bots"}

# 双写性能
dual_write_duration_seconds{quantile="0.95"}

# 灰度流量
gray_release_traffic_percent{stage="1"}
```

### 告警规则

```yaml
# 迁移失败告警
- alert: MigrationFailed
  expr: migration_status{table="bots"} == 0
  for: 5m
  annotations:
    summary: "表迁移失败"

# 数据一致性告警
- alert: ConsistencyTooLow
  expr: consistency_rate{table="bots"} < 0.999
  for: 10m
  annotations:
    summary: "数据一致性低于99.9%"
```

---

## 🔧 故障排查

### 问题 1: DDL 执行卡住

**现象**：DDL执行超过预期时间

**排查**：
```bash
# 查看DDL进度
mysql -h127.0.0.1 -uroot -p -e "SHOW PROCESSLIST;"

# 查看锁等待
mysql -h127.0.0.1 -uroot -p -e "SHOW ENGINE INNODB STATUS\G" | grep -A 20 "LATEST DETECTED DEADLOCK"
```

**解决**：
- 等待DDL完成（大表可能需要较长时间）
- 或在低峰期重新执行

---

### 问题 2: 数据迁移失败

**现象**：迁移工具报错退出

**排查**：
```bash
# 查看迁移日志
tail -100 logs/migrate_bots.log

# 检查断点
mysql -h127.0.0.1 -uroot -p -e "SELECT * FROM migration_checkpoints WHERE table_name = 'bots';"
```

**解决**：
```bash
# 从断点恢复
./bin/tenant-migrate \
  --table bots \
  --resume-from-checkpoint \
  --batch-size 1000
```

---

### 问题 3: 灰度切换后错误率上升

**现象**：灰度切换后错误率 > 0.1%

**排查**：
```bash
# 查看应用日志
tail -100 /var/log/coze/api-server.log

# 查看Prometheus指标
curl http://api-server:8001/metrics | grep http_requests_total
```

**解决**：
```bash
# 立即回滚
curl -X POST http://api-server:8001/api/v1/admin/gray-release/rollback

# 恢复到上一阶段
curl -X POST http://api-server:8001/api/v1/admin/gray-release/config \
  -d '{"stage": 1, "traffic_percent": 5.0}'
```

---

## ✅ 验收清单

### DDL变更
- [ ] 10张表都已添加 tenant_id 字段
- [ ] 索引创建成功
- [ ] 业务功能正常
- [ ] 备份已完成

### 数据迁移
- [ ] 历史数据100%迁移
- [ ] 迁移成功率100%
- [ ] 数据一致性99.9%+
- [ ] 迁移工具测试通过

### 双写验证
- [ ] 双写适配器正常工作
- [ ] 数据一致性校验通过
- [ ] 失败补偿机制正常
- [ ] 性能影响<10%

### 灰度切换
- [ ] 5%流量灰度成功
- [ ] 20%流量灰度成功
- [ ] 50%流量灰度成功
- [ ] 100%流量切换成功

### 清理收尾
- [ ] 双写代码已移除
- [ ] tenant_id 字段为 NOT NULL
- [ ] 文档已更新
- [ ] 团队已培训

---

## 📞 联系方式

| 角色 | 姓名 | 联系方式 | 职责 |
|-----|------|---------|------|
| 迁移负责人 | 研发B | - | 总体协调、技术决策 |
| DBA | - | - | 数据库操作、性能优化 |
| 运维工程师 | - | - | 监控告警、故障处理 |
| 测试工程师 | - | - | 功能验证、回归测试 |

---

**© 2025 ZKER Project. All rights reserved.**
