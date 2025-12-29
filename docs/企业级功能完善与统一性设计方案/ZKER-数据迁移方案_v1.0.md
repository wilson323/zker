# ZKER 智能路由与权限管理系统 - 数据迁移方案

**文档版本**: v1.0
**创建日期**: 2025-01-01
**最后更新**: 2025-01-01
**负责人**: 数据库管理团队
**审批人**: 技术架构委员会

---

## 📋 文档概述

### 迁移目标

本数据迁移方案旨在将现有系统的数据安全、高效地迁移到新的 ZKER 系统架构中，确保：

**核心目标**:
- ✅ **数据完整性**: 零数据丢失，数据一致性 100%
- ✅ **最小停机时间**: 停机时间 < 4 小时（推荐 < 2 小时）
- ✅ **业务连续性**: 迁移过程对用户影响最小化
- ✅ **可回滚性**: 任何阶段出现问题可快速回滚
- ✅ **数据准确性**: 迁移后数据校验通过率 100%

### 迁移范围

**源系统**:
- 当前 Coze 平台数据库（`coze_production`）
- 第三方集成数据（如知识库、对话历史）
- 静态资源文件（用户头像、Bot 头像、附件等）

**目标系统**:
- ZKER 新架构数据库（`zker_production`）
- Milvus 向量数据库
- MinIO 对象存储
- Redis 缓存

**迁移数据类型**:
1. **结构化数据**: 用户、Bot、对话、消息、权限等
2. **非结构化数据**: 文档、图片、音频、视频等
3. **向量数据**: 知识库文档嵌入向量
4. **配置数据**: 系统配置、权限规则、工作流定义等

---

## 🎯 迁移需求分析

### 业务需求

**需求 1: 从单体架构迁移到多租户 SaaS 架构**
```
源系统: 单租户架构，所有数据共享
目标系统: 多租户架构，租户间数据隔离

迁移要点:
- 为每个现有用户分配 tenant_id
- 迁移租户配置（配额、订阅、权限等）
- 更新所有外键关联
```

**需求 2: 从简单权限模型迁移到 RBAC 模型**
```
源系统: 简单的角色-权限映射
目标系统: RBAC（5 级数据权限 + 3 级字段权限）

迁移要点:
- 分析现有权限关系
- 生成新的角色和权限规则
- 迁移用户-角色-权限关联
```

**需求 3: 从单一存储迁移到混合存储**
```
源系统: MySQL 存储所有数据
目标系统: MySQL + Milvus + MinIO + Redis

迁移要点:
- 结构化数据 → MySQL
- 向量数据 → Milvus
- 文件数据 → MinIO
- 热点数据 → Redis
```

### 数据量评估

**数据规模估算** (基于当前系统分析):

| 数据类型 | 记录数 | 数据量 | 增长率 |
|---------|-------|-------|-------|
| 用户 | 100 万 | ~500 MB | 5%/月 |
| Bot | 50 万 | ~5 GB | 10%/月 |
| 对话 | 1000 万 | ~50 GB | 15%/月 |
| 消息 | 1 亿 | ~500 GB | 20%/月 |
| 知识库文档 | 500 万 | ~2 TB | 10%/月 |
| 向量嵌入 | 500 万×1536 维 | ~5 TB | 10%/月 |
| 附件文件 | 200 万 | ~10 TB | 15%/月 |
| **总计** | **~1.2 亿** | **~22 TB** | **~12%/月** |

**迁移时间估算** (基于网络带宽 1 Gbps):
```
理论传输时间 = 22 TB / 1 Gbps = 176,000 秒 ≈ 49 小时

考虑以下因素:
- 数据压缩: 30% 压缩率 → 34 小时
- 并行传输: 4 个并发流 → 8.5 小时
- 转换处理开销: × 1.5 → 13 小时
- 数据校验: × 1.2 → 15.5 小时

实际迁移时间 ≈ 16 小时
```

### 技术约束

**约束 1: 停机时间要求**
```
业务方要求: 停机时间 < 4 小时
推荐方案: 采用滚动迁移 + 停机切换
  - 数据预迁移: 提前 1 周开始，无停机
  - 增量同步: 停机前 2 小时
  - 停机切换: 2 小时（含数据校验）
```

**约束 2: 数据一致性要求**
```
强一致性要求:
  - 用户账户数据
  - 权限数据
  - 财务数据

最终一致性可接受:
  - 对话历史
  - 日志数据
  - 统计数据
```

**约束 3: 存储空间**
```
源系统存储: 30 TB（使用率 75%）
目标系统存储: 50 TB（预留增长空间）

临时空间需求:
  - 备份空间: 30 TB
  - 转换中间数据: 5 TB
  - 总共需求: 85 TB
```

---

## 🔄 迁移策略

### 策略对比

| 迁移策略 | 停机时间 | 复杂度 | 风险 | 适用场景 |
|---------|---------|-------|------|---------|
| **停机迁移** | 长（8-16h） | 低 | 中 | 小规模系统 |
| **滚动迁移** | 中（4-8h） | 中 | 中 | 读多写少系统 |
| **双写迁移** | 短（2-4h） | 高 | 高 | 核心业务系统 |
| **CDC 迁移** | 极短（<1h） | 极高 | 高 | 7×24 业务 |

### 推荐策略: **双写迁移 + 停机切换**

**策略概述**:
```
阶段 1: 数据预迁移（停机前 1 周，无业务影响）
  - 全量数据复制
  - 数据转换和清洗
  - 目标系统数据准备

阶段 2: 双写阶段（停机前 2 天，业务正常运行）
  - 源系统和目标系统同时写入
  - 定期数据一致性校验
  - 性能监控和优化

阶段 3: 增量同步（停机前 2 小时，业务只读）
  - 停止用户写入
  - 最后一次增量同步
  - 数据一致性最终校验

阶段 4: 停机切换（停机 2 小时）
  - 停止源系统服务
  - 切换流量到目标系统
  - 功能验证
  - 源系统进入备用状态

阶段 5: 稳定运行（切换后 1 周）
  - 24 小时监控
  - 性能优化
  - 数据一致性最终验证
```

**优势**:
- ✅ 停机时间短（2 小时）
- ✅ 可随时回滚
- ✅ 数据风险可控
- ✅ 业务影响最小

---

## 🗂️ 数据映射与转换

### 表结构映射

**1. 用户表映射**

源表 (`coze_production.users`):
```sql
CREATE TABLE `users` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `username` VARCHAR(100) NOT NULL UNIQUE,
  `email` VARCHAR(200) NOT NULL UNIQUE,
  `password_hash` VARCHAR(255) NOT NULL,
  `nickname` VARCHAR(100),
  `avatar_url` VARCHAR(500),
  `status` ENUM('active', 'inactive', 'banned') DEFAULT 'active',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

目标表 (`zker_production.users`):
```sql
CREATE TABLE `users` (
  `user_id` VARCHAR(36) PRIMARY KEY,
  `tenant_id` VARCHAR(36) NOT NULL,
  `username` VARCHAR(100) NOT NULL,
  `email` VARCHAR(200) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `nickname` VARCHAR(100),
  `avatar` VARCHAR(500),
  `status` ENUM('active', 'inactive', 'banned') DEFAULT 'active',
  `email_verified` BOOLEAN DEFAULT FALSE,
  `phone` VARCHAR(20),
  `last_login_at` TIMESTAMP NULL,
  `created_by` VARCHAR(36),
  `updated_by` VARCHAR(36),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_tenant_username` (`tenant_id`, `username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**转换脚本**:
```sql
-- 数据迁移 SQL
INSERT INTO zker_production.users (
  user_id,
  tenant_id,
  username,
  email,
  password_hash,
  nickname,
  avatar,
  status,
  email_verified,
  created_by,
  updated_by,
  created_at,
  updated_at
)
SELECT
  UUID() as user_id,
  -- 如果是管理员，分配到系统租户，否则创建个人租户
  IF(
    u.id IN (SELECT user_id FROM admins),
    'system_tenant',
    CONCAT('tenant_', u.id)
  ) as tenant_id,
  u.username,
  u.email,
  u.password_hash,
  u.nickname,
  u.avatar_url as avatar,
  u.status,
  TRUE as email_verified,  -- 源系统假设邮箱已验证
  u.id as created_by,
  u.id as updated_by,
  u.created_at,
  u.updated_at
FROM coze_production.users u;

-- 迁移后处理：为个人用户创建租户
INSERT INTO zker_production.tenants (
  tenant_id,
  tenant_name,
  tenant_type,
  status,
  subscription_tier,
  quota_bots,
  quota_messages,
  created_at
)
SELECT DISTINCT
  u.tenant_id,
  CONCAT(u.username, ' 的个人空间') as tenant_name,
  'individual' as tenant_type,
  'active' as status,
  'free' as subscription_tier,
  10 as quota_bots,
  1000 as quota_messages,
  u.created_at
FROM zker_production.users u
WHERE u.tenant_id NOT IN ('system_tenant', 'public_tenant')
  AND NOT EXISTS (SELECT 1 FROM zker_production.tenants t WHERE t.tenant_id = u.tenant_id);
```

**2. Bot 表映射**

源表 (`coze_production.bots`):
```sql
CREATE TABLE `bots` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `name` VARCHAR(200) NOT NULL,
  `description` TEXT,
  `avatar_url` VARCHAR(500),
  `prompt_template` TEXT,
  `model_config` JSON,
  `status` ENUM('draft', 'published', 'archived') DEFAULT 'draft',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

目标表 (`zker_production.bots`):
```sql
CREATE TABLE `bots` (
  `bot_id` VARCHAR(36) PRIMARY KEY,
  `tenant_id` VARCHAR(36) NOT NULL,
  `created_by` VARCHAR(36) NOT NULL,
  `name` VARCHAR(200) NOT NULL,
  `description` TEXT,
  `avatar` VARCHAR(500),
  `status` ENUM('draft', 'published', 'archived') DEFAULT 'draft',
  `model_config` JSON,
  `prompt_config` JSON,
  `plugin_config` JSON,
  `is_public` BOOLEAN DEFAULT FALSE,
  `version` INT DEFAULT 1,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**转换脚本**:
```sql
-- 数据迁移 SQL
INSERT INTO zker_production.bots (
  bot_id,
  tenant_id,
  created_by,
  name,
  description,
  avatar,
  status,
  model_config,
  prompt_config,
  is_public,
  created_at,
  updated_at
)
SELECT
  UUID() as bot_id,
  u.tenant_id,
  u.user_id as created_by,
  b.name,
  b.description,
  b.avatar_url as avatar,
  b.status,
  b.model_config,
  JSON_OBJECT(
    'system_prompt', b.prompt_template,
    'temperature', 0.7,
    'max_tokens', 2000
  ) as prompt_config,
  FALSE as is_public,
  b.created_at,
  b.updated_at
FROM coze_production.bots b
JOIN zker_production.users u ON b.user_id = CAST(u.created_by AS UNSIGNED)
WHERE u.user_id = CONCAT('user_', b.user_id);
```

**3. 对话表映射**

源表 (`coze_production.conversations`):
```sql
CREATE TABLE `conversations` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `bot_id` BIGINT NOT NULL,
  `title` VARCHAR(200),
  `status` ENUM('active', 'archived') DEFAULT 'active',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_user_id` (`user_id`),
  KEY `idx_bot_id` (`bot_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

目标表 (`zker_production.conversations`):
```sql
CREATE TABLE `conversations` (
  `conv_id` VARCHAR(36) PRIMARY KEY,
  `tenant_id` VARCHAR(36) NOT NULL,
  `bot_id` VARCHAR(36) NOT NULL,
  `created_by` VARCHAR(36) NOT NULL,
  `title` VARCHAR(200),
  `status` ENUM('active', 'archived') DEFAULT 'active',
  `message_count` INT DEFAULT 0,
  `last_message_at` TIMESTAMP NULL,
  `metadata` JSON,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_tenant_bot` (`tenant_id`, `bot_id`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_status_created` (`status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**转换脚本**:
```sql
-- 数据迁移 SQL（分批处理，每批 10 万条）
INSERT INTO zker_production.conversations (
  conv_id,
  tenant_id,
  bot_id,
  created_by,
  title,
  status,
  message_count,
  last_message_at,
  created_at,
  updated_at
)
SELECT
  UUID() as conv_id,
  u.tenant_id,
  b.bot_id,
  u.user_id as created_by,
  c.title,
  c.status,
  (SELECT COUNT(*) FROM coze_production.messages m WHERE m.conv_id = c.id) as message_count,
  (SELECT MAX(created_at) FROM coze_production.messages m WHERE m.conv_id = c.id) as last_message_at,
  c.created_at,
  c.updated_at
FROM coze_production.conversations c
JOIN zker_production.users u ON c.user_id = CAST(SUBSTRING(u.user_id, 6) AS UNSIGNED)
JOIN zker_production.bots b ON c.bot_id = CAST(SUBSTRING(b.bot_id, 5) AS UNSIGNED)
WHERE c.id BETWEEN ? AND ?
ORDER BY c.id
LIMIT 100000;
```

### 字段转换规则

**常见转换规则**:

| 源字段类型 | 目标字段类型 | 转换规则 | 示例 |
|-----------|------------|---------|------|
| BIGINT | VARCHAR(36) | UUID 替换 | 12345 → "a1b2c3d4-e5f6-7890-abcd-ef1234567890" |
| ENUM | VARCHAR | 枚举值转字符串 | 'active' → 'active' |
| JSON | JSON | 字段重命名 + 结构调整 | 见下文 |
| TEXT | TEXT | 无变化 | - |
| TIMESTAMP | TIMESTAMP | 无变化 | - |

**JSON 结构转换示例**:

源系统 `model_config`:
```json
{
  "provider": "openai",
  "model": "gpt-4",
  "apiKey": "sk-xxxxx",
  "temperature": 0.7,
  "maxTokens": 2000,
  "topP": 0.9
}
```

目标系统 `model_config`:
```json
{
  "provider": "openai",
  "model": "gpt-4",
  "credentials": {
    "api_key": "sk-xxxxx",
    "endpoint": "https://api.openai.com/v1"
  },
  "parameters": {
    "temperature": 0.7,
    "max_tokens": 2000,
    "top_p": 0.9,
    "frequency_penalty": 0,
    "presence_penalty": 0
  },
  "fallback": {
    "enabled": true,
    "max_retries": 3,
    "retry_delay": 1000
  }
}
```

**转换函数** (Python):
```python
import json

def transform_model_config(old_config):
  """转换模型配置从旧格式到新格式"""
  if not old_config:
    return None

  old = json.loads(old_config) if isinstance(old_config, str) else old_config

  new_config = {
    'provider': old.get('provider', 'openai'),
    'model': old.get('model', 'gpt-3.5-turbo'),
    'credentials': {
      'api_key': old.get('apiKey', ''),
      'endpoint': get_provider_endpoint(old.get('provider'))
    },
    'parameters': {
      'temperature': old.get('temperature', 0.7),
      'max_tokens': old.get('maxTokens', 2000),
      'top_p': old.get('topP', 0.9),
      'frequency_penalty': old.get('frequencyPenalty', 0),
      'presence_penalty': old.get('presencePenalty', 0)
    },
    'fallback': {
      'enabled': True,
      'max_retries': 3,
      'retry_delay': 1000
    }
  }

  return json.dumps(new_config, ensure_ascii=False)

def get_provider_endpoint(provider):
  """获取不同提供商的 API 端点"""
  endpoints = {
    'openai': 'https://api.openai.com/v1',
    'azure': 'https://YOUR_RESOURCE.openai.azure.com',
    'claude': 'https://api.anthropic.com',
    'qwen': 'https://dashscope.aliyuncs.com/api/v1'
  }
  return endpoints.get(provider, 'https://api.openai.com/v1')
```

---

## 📋 迁移实施步骤

### 阶段 1: 准备阶段（停机前 2 周）

#### 1.1 环境准备

**1.1.1 目标环境搭建**
```bash
# 部署目标系统基础设施
cd zker/docker
cp .env.example .env
# 编辑 .env 配置

# 启动所有服务
docker compose up -d

# 验证服务状态
docker compose ps

# 初始化数据库结构
mysql -h127.0.0.1 -P3306 -uroot -p < docs/企业级功能完善与统一性设计方案/database_schema_ddl.sql

# 初始化基础数据
mysql -h127.0.0.1 -P3306 -uroot -p < docs/企业级功能完善与统一性设计方案/database_init_data.sql
```

**1.1.2 存储空间检查**
```bash
# 检查目标存储空间
df -h

# 检查数据目录
du -sh /var/lib/mysql
du -sh /data/minio
du -sh /data/milvus

# 预估迁移所需空间
python scripts/estimate_storage.py --source-db coze_production
```

**1.1.3 网络带宽测试**
```bash
# 测试源到目标的网络带宽
iperf3 -c target-server -t 60

# 测试数据传输速度
rsync --progress --stats /data/source test-user@target-server:/data/test
```

#### 1.2 数据分析

**1.2.1 数据质量分析**
```sql
-- 检查数据完整性
SELECT
  table_name,
  table_rows,
  data_length,
  index_length,
  (data_length + index_length) as total_size
FROM information_schema.TABLES
WHERE table_schema = 'coze_production'
ORDER BY total_size DESC;

-- 检查孤立记录（外键约束失效）
SELECT
  'messages without conversation' as check_type,
  COUNT(*) as count
FROM coze_production.messages m
WHERE NOT EXISTS (
  SELECT 1 FROM coze_production.conversations c WHERE c.id = m.conv_id
);

-- 检查重复数据
SELECT
  email,
  COUNT(*) as count
FROM coze_production.users
GROUP BY email
HAVING COUNT(*) > 1;

-- 检查 NULL 值
SELECT
  column_name,
  COUNT(*) as null_count
FROM coze_production.users
WHERE username IS NULL
   OR email IS NULL
   OR password_hash IS NULL;
```

**1.2.2 依赖关系分析**
```sql
-- 分析表之间的依赖关系
SELECT
  TABLE_NAME,
  REFERENCED_TABLE_NAME,
  CONSTRAINT_NAME
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = 'coze_production'
  AND REFERENCED_TABLE_NAME IS NOT NULL
ORDER BY TABLE_NAME, REFERENCED_TABLE_NAME;

-- 生成迁移顺序图
# 使用工具: schemaSpy 或 MySQL Workbench
java -jar schemaSpy.jar -t mysql -db coze_production -host localhost -u root -p password -o deps/
```

#### 1.3 迁移工具准备

**1.3.1 迁移脚本开发**
```python
#!/usr/bin/env python3
"""
ZKER 数据迁移工具
支持: MySQL 数据迁移、文件迁移、向量迁移
"""

import argparse
import logging
import sys
from datetime import datetime
from typing import Dict, List, Tuple

import pymysql
import redis
from pymilvus import connections, Collection, FieldSchema, CollectionSchema, DataType

# 配置日志
logging.basicConfig(
  level=logging.INFO,
  format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
  handlers=[
    logging.FileHandler(f'migration_{datetime.now().strftime("%Y%m%d_%H%M%S")}.log'),
    logging.StreamHandler(sys.stdout)
  ]
)
logger = logging.getLogger(__name__)

class DataMigrator:
  """数据迁移器基类"""

  def __init__(self, config: Dict):
    self.config = config
    self.source_db = None
    self.target_db = None
    self.redis_client = None

  def connect_source(self):
    """连接源数据库"""
    self.source_db = pymysql.connect(
      host=self.config['source']['host'],
      port=self.config['source']['port'],
      user=self.config['source']['user'],
      password=self.config['source']['password'],
      database=self.config['source']['database'],
      charset='utf8mb4',
      cursorclass=pymysql.cursors.DictCursor
    )
    logger.info("Connected to source database")

  def connect_target(self):
    """连接目标数据库"""
    self.target_db = pymysql.connect(
      host=self.config['target']['host'],
      port=self.config['target']['port'],
      user=self.config['target']['user'],
      password=self.config['target']['password'],
      database=self.config['target']['database'],
      charset='utf8mb4',
      cursorclass=pymysql.cursors.DictCursor
    )
    logger.info("Connected to target database")

  def connect_redis(self):
    """连接 Redis"""
    self.redis_client = redis.Redis(
      host=self.config['redis']['host'],
      port=self.config['redis']['port'],
      db=self.config['redis']['db'],
      decode_responses=True
    )
    logger.info("Connected to Redis")

  def migrate_table(self, table_name: str, batch_size: int = 10000) -> Tuple[int, int]:
    """
    迁移单张表

    Args:
      table_name: 表名
      batch_size: 批次大小

    Returns:
      (成功记录数, 失败记录数)
    """
    logger.info(f"开始迁移表: {table_name}")

    # 获取总记录数
    with self.source_db.cursor() as cursor:
      cursor.execute(f"SELECT COUNT(*) as total FROM {table_name}")
      total = cursor.fetchone()['total']

    logger.info(f"表 {table_name} 共有 {total} 条记录")

    success_count = 0
    failed_count = 0
    offset = 0

    while offset < total:
      try:
        # 读取源数据
        with self.source_db.cursor() as cursor:
          cursor.execute(f"SELECT * FROM {table_name} LIMIT {batch_size} OFFSET {offset}")
          rows = cursor.fetchall()

        if not rows:
          break

        # 转换并写入目标表
        with self.target_db.cursor() as cursor:
          for row in rows:
            try:
              # 调用转换函数
              transformed = self.transform_row(table_name, row)

              # 插入目标表
              columns = ', '.join(transformed.keys())
              placeholders = ', '.join(['%s'] * len(transformed))
              sql = f"INSERT INTO {table_name} ({columns}) VALUES ({placeholders})"

              cursor.execute(sql, list(transformed.values()))
              success_count += 1

            except Exception as e:
              logger.error(f"记录迁移失败: {row}, 错误: {e}")
              failed_count += 1

          self.target_db.commit()

        offset += batch_size
        logger.info(f"已迁移 {min(offset, total)}/{total} 条记录")

      except Exception as e:
        logger.error(f"批次迁移失败 (offset={offset}): {e}")
        failed_count += len(rows)
        offset += batch_size

    logger.info(f"表 {table_name} 迁移完成: 成功 {success_count}, 失败 {failed_count}")
    return success_count, failed_count

  def transform_row(self, table_name: str, row: Dict) -> Dict:
    """
    转换行数据

    Args:
      table_name: 表名
      row: 原始行数据

    Returns:
      转换后的行数据
    """
    # 根据表名调用不同的转换函数
    transform_funcs = {
      'users': self.transform_user,
      'bots': self.transform_bot,
      'conversations': self.transform_conversation,
      'messages': self.transform_message
    }

    transform_func = transform_funcs.get(table_name)
    if transform_func:
      return transform_func(row)
    else:
      # 默认直接返回
      return row

  def transform_user(self, row: Dict) -> Dict:
    """转换用户数据"""
    return {
      'user_id': f"user_{row['id']}",
      'tenant_id': 'tenant_1',  # 默认租户，后续需要更新
      'username': row['username'],
      'email': row['email'],
      'password_hash': row['password_hash'],
      'nickname': row.get('nickname', ''),
      'avatar': row.get('avatar_url', ''),
      'status': row['status'],
      'email_verified': True,
      'created_by': f"user_{row['id']}",
      'updated_by': f"user_{row['id']}",
      'created_at': row['created_at'],
      'updated_at': row['updated_at']
    }

  def transform_bot(self, row: Dict) -> Dict:
    """转换 Bot 数据"""
    return {
      'bot_id': f"bot_{row['id']}",
      'tenant_id': 'tenant_1',  # 需要根据 user_id 查询
      'created_by': f"user_{row['user_id']}",
      'name': row['name'],
      'description': row.get('description', ''),
      'avatar': row.get('avatar_url', ''),
      'status': row['status'],
      'model_config': row.get('model_config', {}),
      'prompt_config': {'system_prompt': row.get('prompt_template', '')},
      'is_public': False,
      'created_at': row['created_at'],
      'updated_at': row['updated_at']
    }

  def transform_conversation(self, row: Dict) -> Dict:
    """转换对话数据"""
    return {
      'conv_id': f"conv_{row['id']}",
      'tenant_id': 'tenant_1',  # 需要根据 user_id 查询
      'bot_id': f"bot_{row['bot_id']}",
      'created_by': f"user_{row['user_id']}",
      'title': row.get('title', ''),
      'status': row['status'],
      'created_at': row['created_at'],
      'updated_at': row['updated_at']
    }

  def transform_message(self, row: Dict) -> Dict:
    """转换消息数据"""
    return {
      'msg_id': f"msg_{row['id']}",
      'tenant_id': 'tenant_1',
      'conv_id': f"conv_{row['conv_id']}",
      'role': row['role'],
      'content_type': 'text',
      'content': row['content'],
      'created_at': row['created_at']
    }

  def verify_data(self, table_name: str) -> bool:
    """
    验证数据迁移完整性

    Args:
      table_name: 表名

    Returns:
      是否验证通过
    """
    logger.info(f"开始验证表: {table_name}")

    # 对比记录数
    with self.source_db.cursor() as src_cursor:
      src_cursor.execute(f"SELECT COUNT(*) as count FROM {table_name}")
      source_count = src_cursor.fetchone()['count']

    with self.target_db.cursor() as tgt_cursor:
      tgt_cursor.execute(f"SELECT COUNT(*) as count FROM {table_name}")
      target_count = tgt_cursor.fetchone()['count']

    if source_count != target_count:
      logger.error(f"记录数不一致: 源表 {source_count}, 目标表 {target_count}")
      return False

    # 对比数据哈希（抽样）
    with self.source_db.cursor() as src_cursor:
      src_cursor.execute(f"SELECT MD5(GROUP_CONCAT(CONCAT(id, username, email))) as hash FROM {table_name} LIMIT 1000")
      source_hash = src_cursor.fetchone()['hash']

    with self.target_db.cursor() as tgt_cursor:
      tgt_cursor.execute(f"SELECT MD5(GROUP_CONCAT(CONCAT(user_id, username, email))) as hash FROM {table_name} LIMIT 1000")
      target_hash = tgt_cursor.fetchone()['hash']

    if source_hash != target_hash:
      logger.error(f"数据哈希不一致: 源表 {source_hash}, 目标表 {target_hash}")
      return False

    logger.info(f"表 {table_name} 验证通过")
    return True

  def close(self):
    """关闭所有连接"""
    if self.source_db:
      self.source_db.close()
    if self.target_db:
      self.target_db.close()
    if self.redis_client:
      self.redis_client.close()


def main():
  parser = argparse.ArgumentParser(description='ZKER 数据迁移工具')
  parser.add_argument('--config', required=True, help='配置文件路径')
  parser.add_argument('--table', help='指定要迁移的表名，不指定则迁移所有表')
  parser.add_argument('--batch-size', type=int, default=10000, help='批次大小')
  parser.add_argument('--verify', action='store_true', help='迁移后验证数据')

  args = parser.parse_args()

  # 加载配置
  import yaml
  with open(args.config) as f:
    config = yaml.safe_load(f)

  # 创建迁移器
  migrator = DataMigrator(config)

  try:
    # 连接数据库
    migrator.connect_source()
    migrator.connect_target()
    migrator.connect_redis()

    # 执行迁移
    tables = [args.table] if args.table else [
      'users',
      'bots',
      'conversations',
      'messages'
    ]

    for table in tables:
      migrator.migrate_table(table, batch_size=args.batch_size)

      if args.verify:
        if not migrator.verify_data(table):
          logger.error(f"表 {table} 验证失败，停止迁移")
          sys.exit(1)

    logger.info("所有表迁移完成")

  finally:
    migrator.close()


if __name__ == '__main__':
  main()
```

**配置文件示例** (`migration_config.yaml`):
```yaml
source:
  host: source-db.example.com
  port: 3306
  user: root
  password: source_password
  database: coze_production

target:
  host: target-db.example.com
  port: 3306
  user: root
  password: target_password
  database: zker_production

redis:
  host: localhost
  port: 6379
  db: 0

milvus:
  host: localhost
  port: 19530

minio:
  endpoint: localhost:9000
  access_key: minioadmin
  secret_key: minioadmin
  bucket: zker-data
```

**使用示例**:
```bash
# 迁移所有表
python migrate_data.py --config migration_config.yaml --batch-size 10000 --verify

# 迁移单张表
python migrate_data.py --config migration_config.yaml --table users --batch-size 5000
```

---

### 阶段 2: 预迁移阶段（停机前 1 周）

#### 2.1 全量数据迁移

**迁移顺序**（考虑依赖关系）:
```
1. tenants (租户)
2. users (用户)
3. roles (角色)
4. permissions (权限)
5. user_roles (用户-角色关联)
6. bots (Bot)
7. conversations (对话)
8. messages (消息)
9. knowledge_documents (知识库文档)
10. knowledge_chunks (文档分块)
11. workflows (工作流)
12. workflow_executions (工作流执行记录)
```

**执行命令**:
```bash
# 启动全量数据迁移
nohup python migrate_data.py \
  --config migration_config.yaml \
  --batch-size 10000 \
  --verify \
  > migration_full.log 2>&1 &

# 监控迁移进度
tail -f migration_full.log

# 查看迁移统计
grep "已迁移" migration_full.log | tail -20
```

**监控指标**:
```bash
# 检查迁移进度
mysql -h target-db -e "
SELECT
  table_name,
  table_rows,
  ROUND(data_length / 1024 / 1024, 2) as data_mb
FROM information_schema.TABLES
WHERE table_schema = 'zker_production'
ORDER BY table_rows DESC;
"

# 检查源库和目标库记录数差异
python scripts/compare_counts.py \
  --source coze_production \
  --target zker_production
```

#### 2.2 文件迁移

**MinIO 文件迁移脚本**:
```python
#!/usr/bin/env python3
"""
MinIO 文件迁移工具
"""

import os
import hashlib
from minio import Minio
from minio.error import S3Error
import logging

logger = logging.getLogger(__name__)

class FileMigrator:
  """文件迁移器"""

  def __init__(self, source_config: dict, target_config: dict):
    self.source_client = Minio(
      source_config['endpoint'],
      access_key=source_config['access_key'],
      secret_key=source_config['secret_key'],
      secure=False
    )
    self.target_client = Minio(
      target_config['endpoint'],
      access_key=target_config['access_key'],
      secret_key=target_config['secret_key'],
      secure=False
    )
    self.source_bucket = source_config['bucket']
    self.target_bucket = target_config['bucket']

  def migrate_files(self, prefix: str = ''):
    """
    迁移文件

    Args:
      prefix: 文件前缀
    """
    # 确保目标桶存在
    if not self.target_client.bucket_exists(self.target_bucket):
      self.target_client.make_bucket(self.target_bucket)
      logger.info(f"创建目标桶: {self.target_bucket}")

    # 列出源文件
    objects = self.source_client.list_objects(self.source_bucket, prefix=prefix, recursive=True)

    total_count = 0
    success_count = 0
    failed_count = 0
    total_size = 0

    for obj in objects:
      total_count += 1
      total_size += obj.size

      try:
        # 下载源文件
        source_data = self.source_client.get_object(self.source_bucket, obj.object_name)

        # 上传到目标
        self.target_client.put_object(
          self.target_bucket,
          obj.object_name,
          source_data,
          length=obj.size,
          content_type=obj.content_type
        )

        success_count += 1
        logger.info(f"文件迁移成功: {obj.object_name} ({obj.size} bytes)")

        if total_count % 1000 == 0:
          logger.info(f"进度: {success_count}/{total_count}, 失败: {failed_count}")

      except S3Error as e:
        failed_count += 1
        logger.error(f"文件迁移失败: {obj.object_name}, 错误: {e}")

    logger.info(f"文件迁移完成: 总数 {total_count}, 成功 {success_count}, 失败 {failed_count}, 总大小 {total_size / 1024 / 1024:.2f} MB")

  def verify_files(self, prefix: str = '') -> bool:
    """
    验证文件迁移完整性

    Args:
      prefix: 文件前缀

    Returns:
      是否验证通过
    """
    logger.info("开始验证文件完整性...")

    source_objects = self.source_client.list_objects(self.source_bucket, prefix=prefix, recursive=True)
    target_objects = self.target_client.list_objects(self.target_bucket, prefix=prefix, recursive=True)

    source_count = sum(1 for _ in source_objects)
    target_count = sum(1 for _ in target_objects)

    if source_count != target_count:
      logger.error(f"文件数量不一致: 源 {source_count}, 目标 {target_count}")
      return False

    # 抽样验证文件哈希
    source_objects = self.source_client.list_objects(self.source_bucket, prefix=prefix, recursive=True)
    sample_size = min(100, source_count)

    for i, obj in enumerate(source_objects):
      if i >= sample_size:
        break

      # 计算源文件哈希
      source_data = self.source_client.get_object(self.source_bucket, obj.object_name)
      source_hash = hashlib.md5(source_data.read()).hexdigest()

      # 计算目标文件哈希
      target_data = self.target_client.get_object(self.target_bucket, obj.object_name)
      target_hash = hashlib.md5(target_data.read()).hexdigest()

      if source_hash != target_hash:
        logger.error(f"文件哈希不一致: {obj.object_name}")
        return False

    logger.info(f"文件完整性验证通过 (抽样 {sample_size} 个文件)")
    return True


if __name__ == '__main__':
  import yaml

  with open('migration_config.yaml') as f:
    config = yaml.safe_load(f)

  migrator = FileMigrator(
    source_config=config['source_minio'],
    target_config=config['target_minio']
  )

  try:
    migrator.migrate_files()
    if migrator.verify_files():
      logger.info("文件迁移验证通过")
    else:
      logger.error("文件迁移验证失败")
  except Exception as e:
    logger.error(f"文件迁移异常: {e}")
```

#### 2.3 向量数据迁移

**Milvus 向量迁移脚本**:
```python
#!/usr/bin/env python3
"""
Milvus 向量数据迁移工具
"""

import logging
from pymilvus import connections, Collection, FieldSchema, CollectionSchema, DataType, utility
import numpy as np

logger = logging.getLogger(__name__)

class VectorMigrator:
  """向量迁移器"""

  def __init__(self, source_config: dict, target_config: dict):
    connections.connect(alias='source', host=source_config['host'], port=source_config['port'])
    connections.connect(alias='target', host=target_config['host'], port=target_config['port'])

    self.source_alias = 'source'
    self.target_alias = 'target'

  def migrate_collection(self, collection_name: str, batch_size: int = 1000):
    """
    迁移向量集合

    Args:
      collection_name: 集合名称
      batch_size: 批次大小
    """
    logger.info(f"开始迁移向量集合: {collection_name}")

    # 连接源集合
    source_collection = Collection(collection_name, using=self.source_alias)
    source_collection.load()

    # 获取源集合 schema
    source_schema = source_collection.schema

    # 在目标创建集合
    if not utility.has_collection(collection_name, using=self.target_alias):
      target_collection = Collection(
        name=collection_name,
        schema=source_schema,
        using=self.target_alias
      )
      logger.info(f"创建目标集合: {collection_name}")
    else:
      target_collection = Collection(collection_name, using=self.target_alias)

    # 获取总向量数
    total_count = source_collection.num_entities
    logger.info(f"集合 {collection_name} 共有 {total_count} 条向量")

    # 批量迁移
    offset = 0
    success_count = 0

    while offset < total_count:
      try:
        # 从源读取数据
        data = source_collection.fetch(
          ids=list(range(offset, offset + batch_size)),
          output_fields=['*']
        )

        # 插入目标集合
        mr = target_collection.insert(data)
        success_count += mr.insert_count

        offset += batch_size
        logger.info(f"已迁移 {min(offset, total_count)}/{total_count} 条向量")

      except Exception as e:
        logger.error(f"向量迁移失败 (offset={offset}): {e}")
        offset += batch_size

    # 刷新目标集合
    target_collection.flush()

    logger.info(f"向量集合 {collection_name} 迁移完成: 成功 {success_count}")

  def verify_collection(self, collection_name: str) -> bool:
    """
    验证向量迁移完整性

    Args:
      collection_name: 集合名称

    Returns:
      是否验证通过
    """
    logger.info(f"开始验证向量集合: {collection_name}")

    source_collection = Collection(collection_name, using=self.source_alias)
    target_collection = Collection(collection_name, using=self.target_alias)

    source_count = source_collection.num_entities
    target_count = target_collection.num_entities

    if source_count != target_count:
      logger.error(f"向量数量不一致: 源 {source_count}, 目标 {target_count}")
      return False

    logger.info(f"向量集合 {collection_name} 验证通过")
    return True


if __name__ == '__main__':
  import yaml

  with open('migration_config.yaml') as f:
    config = yaml.safe_load(f)

  migrator = VectorMigrator(
    source_config=config['source_milvus'],
    target_config=config['target_milvus']
  )

  try:
    migrator.migrate_collection('knowledge_documents')
    if migrator.verify_collection('knowledge_documents'):
      logger.info("向量迁移验证通过")
  except Exception as e:
    logger.error(f"向量迁移异常: {e}")
```

---

### 阶段 3: 双写阶段（停机前 2 天）

#### 3.1 启用双写

**修改应用层代码**:
```go
// 双写适配器
type DualWriteAdapter struct {
  sourceDB *sql.DB
  targetDB *sql.DB
  enabled  bool
}

func (d *DualWriteAdapter) CreateUser(ctx context.Context, user *User) error {
  // 1. 先写源系统（主）
  err := d.writeToSource(ctx, user)
  if err != nil {
    return err
  }

  // 2. 异步写目标系统（从）
  if d.enabled {
    go func() {
      ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
      defer cancel()

      if err := d.writeToTarget(ctx, user); err != nil {
        log.Printf("写入目标系统失败: %v", err)
        // 记录失败日志，后续补偿
      }
    }()
  }

  return nil
}

func (d *DualWriteAdapter) writeToSource(ctx context.Context, user *User) error {
  _, err := d.sourceDB.ExecContext(ctx,
    "INSERT INTO users (id, username, email) VALUES (?, ?, ?)",
    user.ID, user.Username, user.Email)
  return err
}

func (d *DualWriteAdapter) writeToTarget(ctx context.Context, user *User) error {
  _, err := d.targetDB.ExecContext(ctx,
    "INSERT INTO users (user_id, tenant_id, username, email) VALUES (?, ?, ?, ?)",
    "user_"+user.ID, "tenant_1", user.Username, user.Email)
  return err
}
```

#### 3.2 数据一致性校验

**定期校验脚本**:
```python
#!/usr/bin/env python3
"""
数据一致性校验工具
"""

import hashlib
import logging

logger = logging.getLogger(__name__)

class DataConsistencyChecker:
  """数据一致性检查器"""

  def __init__(self, source_db, target_db):
    self.source_db = source_db
    self.target_db = target_db

  def check_table_consistency(self, table_name: str, key_column: str) -> dict:
    """
    检查表数据一致性

    Args:
      table_name: 表名
      key_column: 主键列名

    Returns:
      检查结果
    """
    logger.info(f"检查表一致性: {table_name}")

    result = {
      'table': table_name,
      'source_count': 0,
      'target_count': 0,
      'missing_in_target': [],
      'missing_in_source': [],
      'data_mismatch': []
    }

    # 获取源表记录数
    with self.source_db.cursor() as cursor:
      cursor.execute(f"SELECT COUNT(*) as count FROM {table_name}")
      result['source_count'] = cursor.fetchone()['count']

    # 获取目标表记录数
    with self.target_db.cursor() as cursor:
      cursor.execute(f"SELECT COUNT(*) as count FROM {table_name}")
      result['target_count'] = cursor.fetchone()['count']

    # 检查缺失记录（抽样 1000 条）
    with self.source_db.cursor() as src_cursor:
      src_cursor.execute(f"SELECT * FROM {table_name} ORDER BY {key_column} LIMIT 1000")
      source_rows = src_cursor.fetchall()

    for row in source_rows:
      key_value = row[key_column]

      # 检查目标表是否存在
      with self.target_db.cursor() as tgt_cursor:
        tgt_cursor.execute(f"SELECT * FROM {table_name} WHERE {key_column} = %s", (key_value,))
        target_row = tgt_cursor.fetchone()

      if not target_row:
        result['missing_in_target'].append(key_value)
      else:
        # 对比数据哈希
        source_hash = self._hash_row(row)
        target_hash = self._hash_row(target_row)

        if source_hash != target_hash:
          result['data_mismatch'].append(key_value)

    logger.info(f"表 {table_name} 一致性检查完成: 源 {result['source_count']}, 目标 {result['target_count']}, 缺失 {len(result['missing_in_target'])}, 不一致 {len(result['data_mismatch'])}")

    return result

  def _hash_row(self, row: dict) -> str:
    """计算行数据哈希"""
    # 排序后拼接
    sorted_values = [str(v) for v in sorted(row.values())]
    row_string = '|'.join(sorted_values)
    return hashlib.md5(row_string.encode()).hexdigest()


if __name__ == '__main__':
  # 定期执行一致性检查
  checker = DataConsistencyChecker(source_db, target_db)

  for table in ['users', 'bots', 'conversations', 'messages']:
    result = checker.check_table_consistency(table, 'id')

    if result['missing_in_target'] or result['data_mismatch']:
      logger.error(f"表 {table} 存在数据不一致，需要补偿")
      # 触发告警
```

---

### 阶段 4: 增量同步与切换（停机 2 小时）

#### 4.1 停机准备

**停机前检查清单**:
```bash
# 1. 确认双写已启用
curl http://api-server:8001/health/dual-write

# 2. 确认数据一致性校验通过
python scripts/check_consistency.py --all-tables

# 3. 备份源系统
mysqldump -h source-db -u root -p coze_production > backup_$(date +%Y%m%d_%H%M%S).sql

# 4. 备份目标系统
mysqldump -h target-db -u root -p zker_production > target_backup_$(date +%Y%m%d_%H%M%S).sql

# 5. 通知用户系统维护
echo "系统将于 22:00-24:00 进行维护，期间无法访问" | sendmail all_users
```

#### 4.2 执行增量同步

**增量同步脚本**:
```python
#!/usr/bin/env python3
"""
增量数据同步工具
"""

import logging
from datetime import datetime, timedelta

logger = logging.getLogger(__name__)

class IncrementalSyncer:
  """增量同步器"""

  def __init__(self, source_db, target_db):
    self.source_db = source_db
    self.target_db = target_db

  def sync_table(self, table_name: str, since: datetime):
    """
    增量同步表

    Args:
      table_name: 表名
      since: 同步起始时间
    """
    logger.info(f"开始增量同步表: {table_name} (自 {since})")

    # 查询增量数据
    with self.source_db.cursor() as cursor:
      cursor.execute(f"""
        SELECT * FROM {table_name}
        WHERE updated_at >= %s
        ORDER BY updated_at
      """, (since,))

      rows = cursor.fetchall()

    if not rows:
      logger.info(f"表 {table_name} 无增量数据")
      return

    # 应用到目标表
    with self.target_db.cursor() as cursor:
      for row in rows:
        try:
          # 使用 INSERT ... ON DUPLICATE KEY UPDATE
          columns = ', '.join(row.keys())
          placeholders = ', '.join(['%s'] * len(row))
          update_clause = ', '.join([f"{k} = VALUES({k})" for k in row.keys() if k != 'id'])

          sql = f"""
            INSERT INTO {table_name} ({columns})
            VALUES ({placeholders})
            ON DUPLICATE KEY UPDATE {update_clause}
          """

          cursor.execute(sql, list(row.values()))

        except Exception as e:
          logger.error(f"增量同步失败: {row}, 错误: {e}")

      self.target_db.commit()

    logger.info(f"表 {table_name} 增量同步完成: {len(rows)} 条记录")


if __name__ == '__main__':
  # 停机前最后增量同步
  sync_time = datetime.now() - timedelta(minutes=30)

  for table in ['users', 'bots', 'conversations', 'messages']:
    syncer.sync_table(table, sync_time)
```

**执行步骤**:
```bash
# T-30 分钟: 停止用户写入（只读模式）
mysql -h source-db -e "SET GLOBAL read_only = ON;"

# T-25 分钟: 最后一次增量同步
python scripts/incremental_sync.py --since "2025-01-01 21:30:00"

# T-15 分钟: 数据一致性最终校验
python scripts/final_verification.py

# T-10 分钟: 切换 DNS
# 1. 更新 DNS 记录，指向新系统
# 2. 降低 DNS TTL（提前 1 天设置 TTL=60）
# 3. 清理所有缓存（CDN、Redis、浏览器缓存）

# T-0 分钟: 流量切换完成，系统恢复
mysql -h target-db -e "SET GLOBAL read_only = OFF;"
```

#### 4.3 切换后验证

**功能验证清单**:
```bash
# 1. 健康检查
curl http://new-api-server:8001/health

# 2. 核心功能测试
python scripts/functional_test.py --url http://new-api-server:8001

# 3. 性能验证
ab -n 1000 -c 100 http://new-api-server:8001/api/v1/bots

# 4. 监控告警
curl http://new-api-server:8001/metrics | prom2json

# 5. 数据抽样查询
mysql -h target-db -e "SELECT COUNT(*) FROM users;"
mysql -h target-db -e "SELECT COUNT(*) FROM conversations;"
mysql -h target-db -e "SELECT COUNT(*) FROM messages;"
```

---

## 🔄 回滚方案

### 回滚触发条件

**必须回滚的情况**:
- ❌ 数据丢失 > 0.01%
- ❌ 核心功能无法使用
- ❌ 系统性能严重下降（QPS < 50% 基线）
- ❌ 出现严重 Bug 导致用户无法使用

**可考虑回滚的情况**:
- ⚠️ 错误率 > 1%
- ⚠️ 响应时间 P95 > 5s
- ⚠️ 数据不一致率 > 0.1%

### 回滚步骤

**快速回滚**（5 分钟内）:
```bash
# 1. 立即停止新系统
docker compose -f zker/docker/docker-compose.yml down

# 2. 恢复源系统服务
mysql -h source-db -e "SET GLOBAL read_only = OFF;"
systemctl start coze-api-server

# 3. 切换 DNS 回源系统
# 更新 DNS A 记录，指向源系统 IP

# 4. 清理 Redis 缓存
redis-cli FLUSHALL

# 5. 验证源系统功能
curl http://source-api:8001/health
```

**完整回滚**（30 分钟内）:
```bash
# 1. 停止双写（如果已启用）
curl -X POST http://api-server:8001/admin/disable-dual-write

# 2. 从备份恢复源系统（如果数据有变化）
mysql -h source-db -u root -p coze_production < backup_before_migration.sql

# 3. 恢复源系统配置
cp /etc/coze/backup/app.conf /etc/coze/app.conf
systemctl restart coze-api-server

# 4. 数据一致性补偿
# 将双写期间的增量数据从目标系统同步回源系统
python scripts/sync_back_to_source.py --since "2025-01-01 22:00:00"

# 5. 全面验证
python scripts/smoke_test.py --source source
```

---

## ✅ 数据验证方案

### 验证层次

**1. 记录数验证**
```sql
-- 对比源表和目标表记录数
SELECT
  'users' as table_name,
  (SELECT COUNT(*) FROM coze_production.users) as source_count,
  (SELECT COUNT(*) FROM zker_production.users) as target_count
UNION ALL
SELECT
  'bots',
  (SELECT COUNT(*) FROM coze_production.bots),
  (SELECT COUNT(*) FROM zker_production.bots)
UNION ALL
SELECT
  'conversations',
  (SELECT COUNT(*) FROM coze_production.conversations),
  (SELECT COUNT(*) FROM zker_production.conversations);
```

**2. 数据哈希验证**
```sql
-- 对比源表和目标表数据哈希
SELECT
  MD5(GROUP_CONCAT(CONCAT(id, username, email ORDER BY id))) as hash
FROM coze_production.users;

SELECT
  MD5(GROUP_CONCAT(CONCAT(user_id, username, email ORDER BY user_id))) as hash
FROM zker_production.users;
```

**3. 业务逻辑验证**
```python
# 验证业务逻辑完整性
def verify_business_logic():
  """验证业务逻辑"""

  # 1. 验证用户-Bot 关联
  with target_db.cursor() as cursor:
    cursor.execute("""
      SELECT COUNT(*) as count
      FROM bots b
      WHERE NOT EXISTS (
        SELECT 1 FROM users u WHERE u.user_id = b.created_by
      )
    """)
    orphan_bots = cursor.fetchone()['count']

    assert orphan_bots == 0, f"存在 {orphan_bots} 个孤立 Bot"

  # 2. 验证对话-消息关联
  with target_db.cursor() as cursor:
    cursor.execute("""
      SELECT COUNT(*) as count
      FROM messages m
      WHERE NOT EXISTS (
        SELECT 1 FROM conversations c WHERE c.conv_id = m.conv_id
      )
    """)
    orphan_messages = cursor.fetchone()['count']

    assert orphan_messages == 0, f"存在 {orphan_messages} 条孤立消息"

  # 3. 验证租户隔离
  with target_db.cursor() as cursor:
    cursor.execute("""
      SELECT tenant_id, COUNT(*) as count
      FROM bots
      GROUP BY tenant_id
      HAVING count > 10000
    """)
    huge_tenants = cursor.fetchall()

    # 检查是否有异常大的租户
    for tenant in huge_tenants:
      logger.warning(f"租户 {tenant['tenant_id']} 有 {tenant['count']} 个 Bots，需要检查")
```

**4. 性能验证**
```bash
# 对比源系统和目标系统性能

# 源系统 QPS
ab -n 10000 -c 100 http://source-api:8001/api/v1/bots

# 目标系统 QPS
ab -n 10000 -c 100 http://target-api:8001/api/v1/bots

# 对比结果（目标系统应 >= 源系统 90%）
```

---

## 🚨 风险评估与应对

### 风险矩阵

| 风险 | 概率 | 影响 | 风险等级 | 应对措施 |
|-----|------|------|---------|---------|
| 数据丢失 | 低 | 极高 | 高 | 多重备份、实时校验、快速回滚 |
| 数据不一致 | 中 | 高 | 高 | 双写验证、定期校验、补偿机制 |
| 停机时间超预期 | 中 | 中 | 中 | 并行迁移、充分预演、应急预案 |
| 性能下降 | 中 | 中 | 中 | 性能测试、容量规划、监控告警 |
| 网络故障 | 低 | 高 | 中 | 冗余网络、断点续传、重试机制 |
| 第三方服务异常 | 低 | 中 | 低 | 降级方案、熔断机制 |

### 关键风险应对

**风险 1: 数据丢失**
```
预防措施:
1. 多重备份
   - 源系统全量备份（停机前）
   - 增量备份（每小时）
   - 目标系统备份（迁移后）

2. 实时校验
   - 每批数据迁移后立即校验
   - 记录级校验（MD5 哈希）
   - 业务逻辑校验

3. 快速回滚
   - 5 分钟内可回滚到源系统
   - 回滚脚本预先准备并测试

应急预案:
- 发现数据丢失立即停止迁移
- 从最近备份恢复
- 重新执行迁移
```

**风险 2: 数据不一致**
```
预防措施:
1. 双写验证
   - 每次双写后记录结果
   - 定期对比源和目标数据
   - 失败自动重试

2. 补偿机制
   - 记录双写失败的数据
   - 后台任务定期补偿
   - 人工介入处理

3. 最终一致性
   - 允许短时间不一致
   - 定期对账任务
   - 自动修复脚本
```

**风险 3: 停机时间超预期**
```
预防措施:
1. 充分预演
   - 停机前 3 天进行全流程演练
   - 记录每个步骤的实际耗时
   - 优化耗时长的步骤

2. 并行迁移
   - 多个表并行迁移
   - 文件和数据并行传输
   - 多线程处理

3. 应急预案
   - 关键步骤超时即跳过
   - 降级为只读模式
   - 分批上线（按用户分批）
```

---

## 📊 迁移监控

### 监控指标

**1. 进度监控**
```bash
# 迁移进度
python scripts/migration_progress.py

# 输出示例:
# users: 1000000/1000000 (100%) ✓
# bots: 450000/500000 (90%) ⏳
# conversations: 8000000/10000000 (80%) ⏳
# messages: 70000000/100000000 (70%) ⏳
# 预计剩余时间: 4 小时 23 分钟
```

**2. 性能监控**
```yaml
# Prometheus 监控指标
- migration_rows_total: 迁移记录总数
- migration_rows_per_second: 迁移速率
- migration_errors_total: 错误总数
- migration_duration_seconds: 迁移耗时
```

**3. 告警规则**
```yaml
# Prometheus 告警规则
groups:
  - name: migration_alerts
    rules:
      - alert: MigrationErrorRateHigh
        expr: rate(migration_errors_total[5m]) > 10
        for: 5m
        annotations:
          summary: "迁移错误率过高"
          description: "5 分钟内错误率 > 10/秒"

      - alert:MigrationStalled
        expr: migration_rows_per_second == 0
        for: 10m
        annotations:
          summary: "迁移已停滞"
          description: "10 分钟内无数据迁移"
```

---

## 📚 附录

### A. 迁移检查清单

**停机前检查** (T-1 周):
- [ ] 目标环境部署完成
- [ ] 所有迁移脚本开发并测试完毕
- [ ] 数据备份完成
- [ ] 网络带宽测试通过
- [ ] 存储空间充足
- [ ] 监控和告警配置完成
- [ ] 回滚方案准备完毕
- [ ] 团队培训完成

**停机前检查** (T-1 天):
- [ ] 双写已启用并验证
- [ ] 数据一致性校验通过
- [ ] 用户通知已发送
- [ ] 运维团队待命
- [ ] 回滚演练完成

**停机期间检查** (T-2h):
- [ ] 源系统进入只读模式
- [ ] 增量同步完成
- [ ] 数据最终校验通过
- [ ] DNS 切换完成
- [ ] 新系统功能验证通过

**停机后检查** (T+1 天):
- [ ] 系统运行稳定
- [ ] 性能指标正常
- [ ] 错误率在正常范围
- [ ] 数据一致性持续监控
- [ ] 用户反馈收集

### B. 常用命令速查

```bash
# 数据库相关
# 导出数据库
mysqldump -h host -u user -p database > backup.sql

# 导入数据库
mysql -h host -u user -p database < backup.sql

# 查看表大小
SELECT
  table_name,
  ROUND(data_length / 1024 / 1024, 2) as data_mb,
  ROUND(index_length / 1024 / 1024, 2) as index_mb
FROM information_schema.TABLES
WHERE table_schema = 'database_name'
ORDER BY data_mb DESC;

# 文件迁移相关
# 同步文件
rsync -avz --progress source/ user@target:/destination/

# 计算文件哈希
md5sum file.txt

# 批量重命名
for file in *.txt; do mv "$file" "prefix_$file"; done
```

### C. 联系方式

| 角色 | 姓名 | 联系方式 | 职责 |
|-----|------|---------|------|
| 迁移总负责人 | [待填写] | [待填写] | 总体协调、决策 |
| 数据库负责人 | [待填写] | [待填写] | 数据库迁移 |
| 应用负责人 | [待填写] | [待填写] | 应用双写、切换 |
| 运维负责人 | [待填写] | [待填写] | 基础设施、监控 |
| 测试负责人 | [待填写] | [待填写] | 功能验证、性能测试 |

---

**文档变更历史**:

| 版本 | 日期 | 变更内容 | 作者 |
|-----|------|---------|------|
| v1.0 | 2025-01-01 | 初始版本 | 数据库管理团队 |

**审批记录**:

| 角色 | 姓名 | 审批意见 | 日期 |
|-----|------|---------|------|
| 技术架构委员会 | [待填写] | [待审批] | [待审批] |
| DBA 负责人 | [待填写] | [待审批] | [待审批] |
| 项目经理 | [待填写] | [待审批] | [待审批] |

---

**© 2025 ZKER Project. All rights reserved.**
