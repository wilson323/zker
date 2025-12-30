# tenant_id 迁移实施指南

**文档版本**: v1.0
**创建日期**: 2025-01-01
**负责人**: 后端架构团队

---

## 📋 迁移概述

### 迁移目标

将现有 Coze Studio 系统从单租户架构升级到多租户 SaaS 架构,核心任务是为所有业务表添加 `tenant_id` 字段,实现租户数据隔离。

### 迁移范围

**涉及表数量**: 24张业务表
**预计工作量**: 8人天
**停机时间**: < 2小时(使用双写方案)

---

## 🎯 迁移策略

### 推荐方案: **双写迁移 + 停机切换**

```
阶段 1: 数据结构准备(无停机)
  ├── 添加 tenant_id 字段(可空)
  └── 创建索引

阶段 2: 数据回填(无停机)
  ├── 基于 space_id 生成租户映射
  └── 分批回填数据,避免长事务

阶段 3: 双写阶段(无停机,2-3天)
  ├── 启用双写适配器
  ├── 新数据同时写旧结构和新结构
  └── 定期数据一致性校验

阶段 4: 停机切换(停机2小时)
  ├── 停止用户写入(只读模式)
  ├── 最后一次增量同步
  ├── 修改 tenant_id 为 NOT NULL
  └── 切换到新架构
```

---

## 📁 迁移文件清单

### SQL 脚本文件

| 文件名 | 说明 | 执行顺序 |
|--------|------|----------|
| `01_add_tenant_id_to_business_tables.sql` | 添加 tenant_id 字段(DDL) | 第1步 |
| `02_backfill_tenant_id_data.sql` | 数据回填(DML) | 第2步 |
| `03_validate_migration.sql` | 迁移验证 | 第3步 |

### 代码文件

| 文件名 | 说明 |
|--------|------|
| `dual_write_adapter.go` | 双写适配器实现 |

---

## 🚀 执行步骤

### 前置准备

#### 1. 环境检查

```bash
# 1. 检查数据库版本(需要 >= 8.0)
mysql --version

# 2. 检查磁盘空间(需要至少 2倍当前数据大小)
df -h

# 3. 检查数据库连接
mysql -h127.0.0.1 -uroot -p -e "SELECT 1"
```

#### 2. 备份数据库

```bash
# 完整备份
mysqldump -h127.0.0.1 -uroot -p \
  --single-transaction \
  --quick \
  --lock-tables=false \
  coze_studio > backup_before_migration_$(date +%Y%m%d_%H%M%S).sql

# 验证备份文件
ls -lh backup_before_migration_*.sql
```

#### 3. 创建迁移数据库用户(可选)

```sql
-- 创建专用迁移用户
CREATE USER IF NOT EXISTS 'migration_user'@'%' IDENTIFIED BY 'secure_password';

-- 授予权限
GRANT SELECT, INSERT, UPDATE, ALTER, CREATE, INDEX ON coze_studio.* TO 'migration_user'@'%';
FLUSH PRIVILEGES;
```

---

### 阶段 1: 数据结构准备(预计30分钟)

#### 执行 `01_add_tenant_id_to_business_tables.sql`

```bash
# 1. 执行 DDL 脚本
mysql -h127.0.0.1 -uroot -p coze_studio < \
  backend/domain/tenant/migration/01_add_tenant_id_to_business_tables.sql

# 2. 验证字段是否添加成功
mysql -h127.0.0.1 -uroot -p coze_studio -e "
SELECT TABLE_NAME, COLUMN_NAME
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = 'coze_studio' AND COLUMN_NAME = 'tenant_id'
ORDER BY TABLE_NAME;"

# 3. 验证索引是否创建
mysql -h127.0.0.1 -uroot -p coze_studio -e "
SHOW INDEX FROM app_draft WHERE Key_name = 'idx_tenant_id';"
```

**预期结果**:
- ✅ 24张表都成功添加 `tenant_id` 字段
- ✅ 每张表都有 `idx_tenant_id` 索引
- ✅ 字段属性为 `VARCHAR(64) NULL`

---

### 阶段 2: 数据回填(预计2-4小时)

#### 执行 `02_backfill_tenant_id_data.sql`

```bash
# 1. 执行数据回填脚本
mysql -h127.0.0.1 -uroot -p coze_studio < \
  backend/domain/tenant/migration/02_backfill_tenant_id_data.sql \
  2>&1 | tee backfill_$(date +%Y%m%d_%H%M%S).log

# 2. 监控执行进度
# 脚本会输出每批次的更新进度,例如:
# app_draft: 更新了 10000 条记录
# app_release_record: 更新了 8567 条记录

# 3. 检查是否有进程卡住
mysql -h127.0.0.1 -uroot -p coze_studio -e "
SHOW PROCESSLIST
WHERE Info LIKE '%UPDATE%'
  AND Time > 300;"
```

**注意事项**:
- ⚠️ 分批处理避免长事务锁定表
- ⚠️ 监控磁盘IO和CPU使用率
- ⚠️ 如果执行时间过长,可以调整 `@batch_size` 参数

---

### 阶段 3: 迁移验证(预计30分钟)

#### 执行 `03_validate_migration.sql`

```bash
# 1. 执行验证脚本
mysql -h127.0.0.1 -uroot -p coze_studio < \
  backend/domain/tenant/migration/03_validate_migration.sql \
  > validation_report_$(date +%Y%m%d_%H%M%S).txt

# 2. 查看验证报告
cat validation_report_*.txt

# 3. 检查关键指标
mysql -h127.0.0.1 -uroot -p coze_studio -e "
-- 检查NULL值数量
SELECT
    TABLE_NAME,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_count
FROM (
    SELECT 'app_draft' as TABLE_NAME, id, tenant_id FROM app_draft
    UNION ALL
    SELECT 'conversation', id, tenant_id FROM conversation
    UNION ALL
    SELECT 'message', id, tenant_id FROM message
) t
GROUP BY TABLE_NAME;"
```

**验证标准**:
- ✅ 所有表的 `tenant_id` NULL 值 = 0
- ✅ 关联表 `tenant_id` 一致性 100%
- ✅ 孤立记录数量 = 0
- ✅ 租户分布合理

**如果验证失败**:

```sql
-- 1. 查找NULL记录
SELECT * FROM app_draft WHERE tenant_id IS NULL LIMIT 10;

-- 2. 手动补充数据
UPDATE app_draft
SET tenant_id = 'tenant_default'
WHERE tenant_id IS NULL;

-- 3. 重新验证
```

---

### 阶段 4: 启用双写适配器(预计2-3天)

#### 4.1 代码集成

```go
// 在 application 层集成双写适配器
// 示例: backend/application/app/app_service.go

package app

import (
    "github.com/coze-studio/backend/domain/tenant/migration"
    "gorm.io/gorm"
)

type AppService struct {
    db               *gorm.DB
    dualWriteAdapter *migration.DualWriteAdapter
}

func NewAppService(db *gorm.DB) *AppService {
    // 创建租户ID解析器
    resolver := migration.SpaceBasedTenantIDResolver(context.Background(), db)

    // 创建双写适配器
    config := migration.DefaultDualWriteConfig
    adapter := migration.NewDualWriteAdapter(db, resolver, config)

    return &AppService{
        db:               db,
        dualWriteAdapter: adapter,
    }
}

// 创建应用时使用双写
func (s *AppService) CreateApp(ctx context.Context, req *CreateAppRequest) error {
    data := map[string]interface{}{
        "id":          generateID(),
        "space_id":    req.SpaceID,
        "name":        req.Name,
        "description": req.Description,
        "created_at":  time.Now().UnixMilli(),
    }

    // 使用双写适配器插入
    return s.dualWriteAdapter.InsertWithTenantID(ctx, "app_draft", data)
}
```

#### 4.2 启用双写

```bash
# 1. 部署新版本代码(包含双写适配器)
make deploy

# 2. 监控双写指标
curl http://api-server:8001/metrics | grep dual_write

# 3. 检查双写失败队列
# 如果失败队列持续增长,说明有问题
```

#### 4.3 双写期间监控

```bash
# 每天执行一次一致性检查
mysql -h127.0.0.1 -uroot -p coze_studio < \
  backend/domain/tenant/migration/03_validate_migration.sql

# 检查双写指标
# - dual_write_success_total 应该持续增长
# - dual_write_failed_total 应该为 0 或很低
# - avg_write_latency 应该 < 100ms
```

---

### 阶段 5: 停机切换(预计2小时)

#### 5.1 停机前检查(T-30分钟)

```bash
# 1. 确认双写已运行2-3天
# 2. 确认数据一致性验证通过
# 3. 通知用户系统维护

# 4. 最后一次备份
mysqldump -h127.0.0.1 -uroot -p coze_studio > backup_before_cutover.sql
```

#### 5.2 执行切换(T=0)

```bash
# T-30分钟: 停止用户写入
mysql -h127.0.0.1 -uroot -p coze_studio -e "SET GLOBAL read_only = ON;"

# T-25分钟: 最后一次增量同步
# (如果双写运行正常,这一步应该很快)

# T-20分钟: 修改 tenant_id 为 NOT NULL
mysql -h127.0.0.1 -uroot -p coze_studio << 'EOF'
ALTER TABLE app_draft MODIFY COLUMN tenant_id VARCHAR(64) NOT NULL;
ALTER TABLE app_release_record MODIFY COLUMN tenant_id VARCHAR(64) NOT NULL;
ALTER TABLE conversation MODIFY COLUMN tenant_id VARCHAR(64) NOT NULL;
ALTER TABLE message MODIFY COLUMN tenant_id VARCHAR(64) NOT NULL;
-- ... 其他表
EOF

# T-15分钟: 验证数据完整性
mysql -h127.0.0.1 -uroot -p coze_studio < \
  backend/domain/tenant/migration/03_validate_migration.sql

# T-10分钟: 切换流量到新系统
# 更新负载均衡配置,重启应用服务

# T-0分钟: 恢复服务
mysql -h127.0.0.1 -uroot -p coze_studio -e "SET GLOBAL read_only = OFF;"

# T+5分钟: 验证核心功能
curl http://api-server:8001/health
curl http://api-server:8001/api/v1/bots
```

---

## ✅ 验收标准

### 功能验收

- [ ] 所有业务表都有 `tenant_id` 字段且为 NOT NULL
- [ ] 所有历史数据都正确关联到租户
- [ ] 新创建的数据自动填充 `tenant_id`
- [ ] 租户间数据隔离正确(跨租户查询返回空)
- [ ] 核心API功能正常(创建Bot、对话、知识库等)

### 性能验收

- [ ] 查询性能没有明显下降(P95 < 200ms)
- [ ] 写入延迟增加 < 20%
- [ ] 数据库连接数正常
- [ ] 磁盘IO使用率 < 70%

### 数据完整性验收

- [ ] 没有孤儿记录
- [ ] 关联关系 100% 一致
- [ ] 没有数据丢失
- [ ] 没有重复数据

---

## 🔄 回滚方案

### 触发条件

- ❌ 数据丢失 > 0.01%
- ❌ 核心功能无法使用
- ❌ 性能严重下降(QPS < 50% 基线)

### 快速回滚(5分钟)

```bash
# 1. 停止新系统
docker compose down

# 2. 从备份恢复
mysql -h127.0.0.1 -uroot -p coze_studio < backup_before_cutover.sql

# 3. 恢复服务
docker compose up -d

# 4. 验证功能
curl http://api-server:8001/health
```

---

## 📊 迁移监控

### 关键指标

```yaml
# Prometheus 指标
migration_tenant_id_fill_percent:
  # tenant_id 填充百分比,目标 100%

migration_dual_write_latency_seconds:
  # 双写延迟,目标 < 0.1s

migration_data_consistency_percent:
  # 数据一致性百分比,目标 100%

migration_orphan_records_count:
  # 孤立记录数,目标 0
```

### 告警规则

```yaml
# 告警规则
alerts:
  - name: TenantIdFillTooSlow
    expr: migration_tenant_id_fill_percent < 99
    for: 30m
    annotations:
      summary: "tenant_id 填充进度过慢"

  - name: DualWriteLatencyHigh
    expr: migration_dual_write_latency_seconds > 0.5
    for: 5m
    annotations:
      summary: "双写延迟过高"

  - name: DataConsistencyLow
    expr: migration_data_consistency_percent < 100
    for: 5m
    annotations:
      summary: "数据一致性检查失败"
```

---

## 📞 联系方式

| 角色 | 姓名 | 联系方式 | 职责 |
|-----|------|---------|------|
| 迁移总负责人 | [待填写] | [待填写] | 总体协调、决策 |
| 数据库负责人 | [待填写] | [待填写] | 数据库迁移 |
| 应用负责人 | [待填写] | [待填写] | 双写适配器集成 |
| 运维负责人 | [待填写] | [待填写] | 基础设施、监控 |

---

## 📚 附录

### A. 常见问题

**Q1: 如果回填脚本执行时间过长怎么办?**

A: 可以调整 `@batch_size` 参数,减小批次大小,减少锁表时间。

**Q2: 双写失败会影响主流程吗?**

A: 不会。双写适配器采用异步双写,双写失败不影响主库写入。

**Q3: 如何验证迁移是否成功?**

A: 执行 `03_validate_migration.sql` 脚本,检查所有验证结果是否为 PASS。

### B. 参考资料

- [ZKER 数据迁移方案 v1.0](../../../docs/企业级功能完善与统一性设计方案/ZKER-数据迁移方案_v1.0.md)
- [ZKER 企业级开发规范手册](../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [MySQL 官方文档 - ALTER TABLE](https://dev.mysql.com/doc/refman/8.0/en/alter-table.html)

---

**文档变更历史**:

| 版本 | 日期 | 变更内容 | 作者 |
|-----|------|---------|------|
| v1.0 | 2025-01-01 | 初始版本 | 后端架构团队 |

---

**© 2025 ZKER Project. All rights reserved.**
