# ZKER 企业级功能快速开始指南

**版本**: v1.0.0
**最后更新**: 2025-01-01
**适用对象**: 企业管理员、运维人员、开发者

---

## 📋 目录

- [功能概述](#功能概述)
- [5分钟快速部署](#5分钟快速部署)
- [租户管理](#租户管理)
- [配额管理](#配额管理)
- [权限管理](#权限管理)
- [监控告警](#监控告警)
- [常见场景](#常见场景)
- [下一步](#下一步)

---

## 🎯 功能概述

### 什么是企业级功能？

**ZKER 企业级功能** 是一套完整的多租户 SaaS 解决方案，包括：

1. **多租户管理**：支持多个租户隔离，数据安全隔离
2. **配额管理**：灵活的资源配额控制，支持软限制和超量计费
3. **权限管理**：细粒度的 RBAC 权限控制，支持数据权限和字段权限
4. **订阅管理**：支持多种订阅套餐，灵活的计费周期
5. **监控告警**：完整的 Prometheus + Grafana 监控体系

### 核心优势

| 功能 | 传统版本 | 企业级版本 |
|------|---------|-----------|
| **多租户** | ❌ 不支持 | ✅ 完整的租户隔离 |
| **权限控制** | ❌ 简单权限 | ✅ 5级数据权限 + 3级字段权限 |
| **配额限制** | ❌ 无限制 | ✅ 灵活的资源配额 |
| **监控告警** | ❌ 基础监控 | ✅ 企业级监控 + 智能告警 |
| **数据安全** | ⚠️ 基础安全 | ✅ 行业级安全标准 |

---

## ⚡ 5分钟快速部署

### 前置要求

- ✅ Docker 和 Docker Compose 已安装
- ✅ Git 已安装
- ✅ 至少 8GB 内存、4核心 CPU

### 快速部署步骤

#### 1. 启动所有服务

```bash
# 克隆代码（如果还没有）
git clone https://github.com/coze-dev/coze-studio.git
cd coze-studio

# 复制环境变量配置
cp .env.example .env

# 启动所有服务（包括MySQL、Redis、Grafana等）
docker compose -f ./docker/docker-compose.yml up -d

# 等待服务启动（约2-3分钟）
docker compose -f ./docker/docker-compose.yml ps
```

**预期输出**：
```
NAME                    STATUS         PORTS
coze-studio-api         Up             0.0.0.0:8080->8080/tcp
coze-studio-mysql       Up             0.0.0.0:3306->3306/tcp
coze-studio-redis       Up             0.0.0.0:6379->6379/tcp
coze-studio-grafana     Up             0.0.0.0:3000->3000/tcp
coze-studio-prometheus  Up             0.0.0.0:9090->9090/tcp
```

#### 2. 初始化企业级数据

```bash
# 执行数据库迁移
docker compose exec coze-studio-api make migrate

# 初始化系统预置角色
docker compose exec coze-studio-api go run backend/cmd/init_roles/main.go
```

**预期输出**：
```
✅ 创建租户表成功
✅ 创建配额表成功
✅ 创建权限表成功
✅ 创建预置角色成功
- 超级管理员 (super_admin)
- 管理员 (admin)
- 普通用户 (user)
- 客服专员 (customer_service)
- 访客 (guest)
```

#### 3. 访问系统

```bash
# 打开浏览器访问
open http://localhost:8888/sign
```

**注册超级管理员账号**：
- 用户名：`admin`
- 邮箱：`admin@example.com`
- 密码：`Admin123!`

#### 4. 验证部署

```bash
# 检查API健康状态
curl http://localhost:8080/api/health/status

# 预期输出: {"status":"ok"}

# 检查Prometheus指标
curl http://localhost:9090/api/v1/query?query=up

# 访问Grafana监控大盘
open http://localhost:3000
# 默认账号: admin / admin
```

---

## 🏢 租户管理

### 创建第一个租户

**方式一：通过API创建**

```bash
curl -X POST http://localhost:8080/api/v1/tenants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "tenant_name": "示例公司",
    "tenant_type": "enterprise",
    "admin_email": "admin@example.com",
    "admin_password": "SecurePassword123!"
  }'
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant_001",
    "tenant_name": "示例公司",
    "tenant_type": "enterprise",
    "status": "active",
    "subscription_tier": "free",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

**方式二：通过前端界面**

1. 登录后台管理系统：`http://localhost:8888/admin`
2. 导航到 **租户管理** → **租户列表**
3. 点击 **创建租户** 按钮
4. 填写租户信息：
   - 租户名称：`示例公司`
   - 租户类型：`企业版`
   - 管理员邮箱：`admin@example.com`
   - 管理员密码：`SecurePassword123!`
5. 点击 **提交**

### 查看租户详情

**API查询**：
```bash
curl http://localhost:8080/api/v1/tenants/tenant_001 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应字段说明**：

| 字段 | 说明 | 示例 |
|------|------|------|
| `tenant_id` | 租户唯一标识 | `tenant_001` |
| `tenant_name` | 租户名称 | `示例公司` |
| `tenant_type` | 租户类型 | `individual` / `team` / `enterprise` |
| `status` | 租户状态 | `active` / `suspended` / `deleted` |
| `subscription_tier` | 订阅套餐 | `free` / `pro` / `enterprise` |

### 租户隔离验证

**验证数据隔离**：
```bash
# 使用租户A的Token查询Bot列表
curl http://localhost:8080/api/bots \
  -H "X-Tenant-ID: tenant_001" \
  -H "Authorization: Bearer TOKEN_A"

# 使用租户B的Token查询Bot列表
curl http://localhost:8080/api/bots \
  -H "X-Tenant-ID: tenant_002" \
  -H "Authorization: Bearer TOKEN_B"

# 两个租户看到的Bot列表完全不同，数据完全隔离
```

---

## 📊 配额管理

### 查看配额使用情况

**API查询**：
```bash
curl http://localhost:8080/api/v1/tenants/tenant_001/quotas \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例**：
```json
{
  "code": 0,
  "data": {
    "quotas": [
      {
        "resource_type": "bots",
        "limit": 100,
        "used": 15,
        "remaining": 85,
        "unit": "count"
      },
      {
        "resource_type": "users",
        "limit": 1000,
        "used": 150,
        "remaining": 850,
        "unit": "count"
      },
      {
        "resource_type": "api_calls",
        "limit": 1000000,
        "used": 250000,
        "remaining": 750000,
        "unit": "count"
      }
    ],
    "overall_usage": 15.5
  }
}
```

### 配额类型说明

| 资源类型 | 说明 | 免费版限制 | 专业版限制 | 企业版限制 |
|---------|------|-----------|-----------|-----------|
| `bots` | Bot数量 | 10 | 100 | 无限制 |
| `users` | 用户数量 | 100 | 1000 | 无限制 |
| `api_calls` | API调用次数/月 | 10万 | 100万 | 无限制 |
| `storage` | 存储空间(GB) | 10 | 100 | 1000 |
| `knowledge_base` | 知识库数量 | 5 | 50 | 无限制 |

### 检查配额是否超限

**API调用**：
```bash
curl -X POST http://localhost:8080/api/v1/tenants/tenant_001/quotas/bots/check \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1
  }'
```

**响应示例（配额充足）**：
```json
{
  "code": 0,
  "data": {
    "allowed": true,
    "remaining": 85,
    "message": "配额充足"
  }
}
```

**响应示例（配额不足）**：
```json
{
  "code": 108010001,
  "message": "配额不足：Bot数量已达到上限 100，无法创建更多Bot",
  "data": {
    "allowed": false,
    "remaining": 0,
    "current_usage": 100,
    "limit": 100
  }
}
```

### 升级订阅套餐

**通过API升级**：
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions/sub_001/upgrade \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "target_tier": "pro"
  }'
```

---

## 🔐 权限管理

### 创建自定义角色

**API调用**：
```bash
curl -X POST http://localhost:8080/api/v1/roles \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_001" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "role_name": "内容审核员",
    "role_code": "content_reviewer",
    "role_type": "custom",
    "description": "负责审核Bot发布的内容"
  }'
```

### 配置数据权限

**数据权限级别**：

| 级别 | 说明 | 适用场景 |
|------|------|---------|
| `ALL` | 全部数据 | 超级管理员 |
| `DEPARTMENT` | 本部门数据 | 部门经理 |
| `OWN` | 仅自己的数据 | 普通用户 |
| `CUSTOM` | 自定义范围 | 特殊权限需求 |
| `NONE` | 无权限 | 禁止访问 |

**API配置示例**：
```bash
curl -X POST http://localhost:8080/api/v1/roles/role_001/data-permissions \
  -H "Content-Type: application/json" \
  -d '{
    "resource_type": "bots",
    "permission_scope": "DEPARTMENT",
    "department_id": "dept_001"
  }'
```

### 检查用户权限

**API调用**：
```bash
curl -X POST http://localhost:8080/api/v1/permissions/check-data \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_001",
    "resource_type": "bots",
    "resource_id": "bot_123",
    "permission": "read"
  }'
```

**响应示例**：
```json
{
  "code": 0,
  "data": {
    "allowed": true,
    "reason": "用户拥有读取此Bot的权限",
    "permission_type": "DEPARTMENT",
    "matched_roles": ["content_reviewer"]
  }
}
```

---

## 📈 监控告警

### 访问Grafana监控大盘

**启动步骤**：
```bash
# 1. 确保Grafana正在运行
docker compose ps | grep grafana

# 2. 访问Grafana
open http://localhost:3000

# 3. 登录（默认账号）
# 用户名: admin
# 密码: admin
```

**导入监控大盘**：
1. 点击左侧 **+** → **Import**
2. 上传文件：`deploy/monitoring/grafana-dashboard.json`
3. 选择 **Prometheus** 数据源
4. 点击 **Import**

### 核心监控指标

**系统资源指标**：
- CPU使用率（阈值：80%）
- 内存使用率（阈值：85%）
- 磁盘使用率（阈值：85%）

**API性能指标**：
- P95响应时间（阈值：500ms）
- API错误率（阈值：1%）
- API请求量（RPS）

**业务指标**：
- Bot调用成功率（阈值：> 95%）
- 配额使用率（阈值：90%）
- Token使用速率

### 配置告警通知

**步骤**：
1. 在Grafana中，导航到 **Alerting** → **Notification channels**
2. 添加新的通知渠道：
   - 类型：**Webhook**
   - URL：企业微信/钉钉Webhook URL
3. 测试通知是否正常

**告警通知示例**：
```json
{
  "title": "🚨 CPU使用率过高",
  "message": "实例 server-01 CPU使用率持续 > 80% 已5分钟",
  "severity": "warning",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

---

## 🎓 常见场景

### 场景1：为新客户开通服务

**步骤**：
1. 创建租户
2. 分配配额
3. 创建管理员账号
4. 分配角色权限

**完整脚本**：
```bash
#!/bin/bash
# scripts/onboard_new_customer.sh

TENANT_NAME=$1
ADMIN_EMAIL=$2

echo "为客户 $TENANT_NAME 开通服务..."

# 1. 创建租户
TENANT_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/tenants \
  -H "Content-Type: application/json" \
  -d "{
    \"tenant_name\": \"$TENANT_NAME\",
    \"tenant_type\": \"enterprise\",
    \"admin_email\": \"$ADMIN_EMAIL\",
    \"admin_password\": \"Welcome123!\"
  }")

TENANT_ID=$(echo $TENANT_RESPONSE | jq -r '.data.tenant_id')

echo "✅ 租户创建成功: $TENANT_ID"

# 2. 设置配额（企业版默认）
curl -X POST http://localhost:8080/api/v1/tenants/$TENANT_ID/quotas \
  -H "Content-Type: application/json" \
  -d '{
    "bots": 1000,
    "users": 10000,
    "api_calls": 10000000,
    "storage": 1000
  }'

echo "✅ 配额设置完成"

# 3. 分配默认角色
curl -X POST http://localhost:8080/api/v1/users/assign-role \
  -H "Content-Type: application/json" \
  -d "{
    \"user_email\": \"$ADMIN_EMAIL\",
    \"role_code\": \"admin\"
  }"

echo "✅ 角色分配完成"

echo "🎉 客户 $TENANT_NAME 开通完成！"
```

### 场景2：监控配额使用并自动升级

**Python脚本示例**：
```python
#!/usr/bin/env python3
# scripts/monitor_quota_usage.py

import requests
import smtplib
from email.mime.text import MIMEText

API_BASE = "http://localhost:8080"
TOKEN = "YOUR_ADMIN_TOKEN"

def check_quota_usage(tenant_id):
    """检查租户配额使用情况"""
    response = requests.get(
        f"{API_BASE}/api/v1/tenants/{tenant_id}/quotas/usage/details",
        headers={"Authorization": f"Bearer {TOKEN}"}
    )
    data = response.json()

    for quota in data['data']:
        # 如果任一资源使用率 > 90%，发送告警
        if quota['usage_percentage'] > 90:
            send_alert(tenant_id, quota)
            # 自动升级套餐（可选）
            if quota['usage_percentage'] > 95:
                upgrade_subscription(tenant_id)

def send_alert(tenant_id, quota):
    """发送告警邮件"""
    msg = MIMEText(f"""
    租户 {tenant_id} 配额告警：

    资源类型: {quota['resource_type']}
    当前使用: {quota['current_usage']}
    配额限制: {quota['limit']}
    使用率: {quota['usage_percentage']:.1f}%

    请及时升级套餐或清理资源。
    """)

    msg['Subject'] = f"配额告警：{tenant_id}"
    msg['From'] = 'noreply@coze-studio.com'
    msg['To'] = 'admin@coze-studio.com'

    # 发送邮件
    smtplib.SMTP('localhost').send_message(msg)

def upgrade_subscription(tenant_id):
    """自动升级订阅套餐"""
    response = requests.post(
        f"{API_BASE}/api/v1/subscriptions/{tenant_id}/upgrade",
        headers={
            "Authorization": f"Bearer {TOKEN}",
            "Content-Type": "application/json"
        },
        json={"target_tier": "pro"}
    )
    print(f"✅ 租户 {tenant_id} 已自动升级到Pro版")

if __name__ == '__main__':
    # 检查所有活跃租户
    tenants = requests.get(
        f"{API_BASE}/api/v1/tenants",
        headers={"Authorization": f"Bearer {TOKEN}"}
    ).json()['data']['tenants']

    for tenant in tenants:
        if tenant['status'] == 'active':
            check_quota_usage(tenant['tenant_id'])
```

### 场景3：调试权限问题

**调试步骤**：
```bash
# 1. 检查用户的角色
curl http://localhost:8080/api/v1/users/user_001/roles \
  -H "Authorization: Bearer YOUR_TOKEN"

# 2. 检查角色的权限
curl http://localhost:8080/api/v1/roles/role_001 \
  -H "Authorization: Bearer YOUR_TOKEN"

# 3. 检查具体的权限检查结果
curl -X POST http://localhost:8080/api/v1/permissions/check-data \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_001",
    "resource_type": "bots",
    "resource_id": "bot_123",
    "permission": "write"
  }'

# 4. 查看审计日志（如果启用）
curl http://localhost:8080/api/v1/audit/logs \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -G -d "user_id=user_001" \
  -d "action=check_permission"
```

---

## 🚀 下一步

### 深入学习

**必读文档**：
1. [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md) - 完整的开发规范
2. [性能基线文档](./ZKER-性能基线文档.md) - 性能基线和优化建议
3. [故障排查手册](./ZKER-故障排查手册_v1.0.md) - 常见问题排查
4. [监控告警阈值调优指南](./ZKER-监控告警阈值调优指南.md) - 监控调优

**开发指南**：
- [前端页面架构设计文档](./ZKER-前端页面架构设计文档.md)
- [前端API使用指南](../../frontend/packages/arch/bot-api/ENTERPRISE_API.md)
- [数据库设计完整交付清单](./数据库设计完整交付清单.md)

**运维指南**：
- [生产环境压力测试执行指南](./ZKER-生产环境压力测试执行指南.md)
- [灰度发布策略](./ZKER-灰度发布策略_v1.0.md)
- [一键回滚方案](./ZKER-一键回滚方案_v1.0.md)

### 获取帮助

**遇到问题？**

1. **查看FAQ**：[常见问题解答](https://github.com/coze-dev/coze-studio/wiki/9.-FAQ)
2. **提交Issue**：[GitHub Issues](https://github.com/coze-dev/coze-studio/issues)
3. **加入社区**：
   - Discord: [Coze Community](https://discord.gg/sTVN9EVS4B)
   - Telegram: [Coze](https://t.me/+pP9CkPnomDA0Mjgx)

**联系支持团队**：
- 邮箱：support@coze-studio.com
- 企业微信：CozeStudio技术支持群

---

## 📚 附录

### 常用API端点速查

| 功能 | 方法 | 端点 |
|------|------|------|
| **租户管理** | | |
| 创建租户 | POST | `/api/v1/tenants` |
| 获取租户详情 | GET | `/api/v1/tenants/{tenant_id}` |
| 更新租户 | PUT | `/api/v1/tenants/{tenant_id}` |
| 删除租户 | DELETE | `/api/v1/tenants/{tenant_id}` |
| **配额管理** | | |
| 查询配额 | GET | `/api/v1/tenants/{tenant_id}/quotas` |
| 检查配额 | POST | `/api/v1/tenants/{tenant_id}/quotas/{resource_type}/check` |
| 更新配额限制 | PUT | `/api/v1/tenants/{tenant_id}/quotas/{resource_type}/limit` |
| **权限管理** | | |
| 创建角色 | POST | `/api/v1/roles` |
| 获取角色详情 | GET | `/api/v1/roles/{role_id}` |
| 检查数据权限 | POST | `/api/v1/permissions/check-data` |
| 获取权限矩阵 | GET | `/api/v1/permissions/matrix` |

### 环境变量配置

**必需配置**：
```bash
# 数据库配置
MYSQL_HOST=mysql
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_DATABASE=coze_studio

# Redis配置
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT密钥（生产环境必须修改）
JWT_SECRET=your-secret-key-change-in-production

# API密钥（用于模型服务）
MODEL_API_KEY=your-model-api-key
```

**可选配置**：
```bash
# 配额检查模式
QUOTA_ENFORCEMENT_MODE=strict  # strict | permissive | disabled

# 权限检查模式
PERMISSION_CHECK_MODE=strict  # strict | permissive

# 监控配置
PROMETHEUS_ENABLED=true
GRAFANA_ENABLED=true

# 日志级别
LOG_LEVEL=info  # debug | info | warn | error
```

---

**更新日志**：

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，完整的企业级功能快速开始指南 | Claude AI |
