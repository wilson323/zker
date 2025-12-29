# ZKER 智能路由与权限管理系统 - 灰度发布策略

**文档版本**: v1.0
**创建日期**: 2025-01-01
**最后更新**: 2025-01-01
**负责人**: 发布管理团队
**审批人**: 技术架构委员会

---

## 📋 文档概述

### 灰度发布目标

本灰度发布策略旨在实现新版本系统的平滑、安全上线，确保：

**核心目标**:
- ✅ **风险可控**: 逐步放量，快速发现问题，降低影响范围
- ✅ **用户体验**: 最小化对用户的干扰和中断
- ✅ **数据驱动**: 基于监控指标和用户反馈做出决策
- ✅ **快速回滚**: 任何阶段出现问题可立即回滚
- ✅ **平滑过渡**: 从 0% 到 100% 流量平滑切换

### 灰度发布适用场景

**适合灰度发布的场景**:
- ✅ 核心业务逻辑变更
- ✅ API 接口变更
- ✅ 数据库架构变更
- ✅ 性能优化变更
- ✅ UI/UX 重大调整
- ✅ 第三方服务集成

**不适合灰度发布的场景**:
- ❌ 紧急安全补丁（直接全量发布）
- ❌ 配置修正（无风险，直接发布）
- ❌ 文档更新（不涉及代码）

---

## 🎯 灰度发布策略

### 策略对比

| 策略 | 优点 | 缺点 | 适用场景 |
|-----|------|------|---------|
| **按用户 ID** | 实现简单，可回溯 | 用户分散，不易测试 | 内部用户、早期采用者 |
| **按用户属性** | 目标明确，便于对比 | 需要用户画像数据 | 特定用户群体测试 |
| **按地域** | 隔离性好，易监控 | 地域差异可能影响结果 | 多地域系统 |
| **按流量比例** | 平滑放量，易控制 | 无法针对特定用户 | 全量发布前最后阶段 |
| **按功能开关** | 灵活性高，风险低 | 代码复杂度增加 | 功能模块级灰度 |

### 推荐策略: **组合策略**

**阶段 1: 内部用户灰度** (1-2 天)
```
策略: 按用户 ID（内部员工白名单）
流量: 100% 内部用户
目的: 验证核心功能无重大缺陷
```

**阶段 2: 早期采用者灰度** (2-3 天)
```
策略: 按用户属性（注册时间 > 6 个月，活跃度高）
流量: 5% 全站用户
目的: 验证用户体验和性能
```

**阶段 3: 地域灰度** (3-5 天)
```
策略: 按地域（选择非核心市场）
流量: 20% 全站用户（特定地域 100%）
目的: 验证地域特定功能和性能
```

**阶段 4: 流量比例灰度** (5-7 天)
```
策略: 按流量比例（逐步放量）
流量: 50% → 80% → 100%
目的: 验证系统稳定性和性能
```

---

## 🔄 灰度发布阶段

### 阶段 0: 准备阶段（发布前 1 周）

#### 0.1 发布前检查

**代码检查**:
```bash
# 1. 代码审查通过
# 所有 PR 已 review 并 approved

# 2. 单元测试通过
rush test --coverage

# 3. 集成测试通过
rush test:integration

# 4. E2E 测试通过
rush test:e2e

# 5. 性能基准测试通过
npm run perf:baseline
```

**文档检查**:
```bash
# 1. CHANGELOG 更新
grep "## [Unreleased]" CHANGELOG.md

# 2. API 文档更新
# OpenAPI 规范已更新并提交

# 3. 数据库迁移脚本准备
ls docs/企业级功能完善与统一性设计方案/migrations/*.sql

# 4. 回滚脚本准备
ls scripts/rollback/*.sh
```

**环境检查**:
```bash
# 1. 预发布环境部署完成
ssh pre-deploy-server "docker ps | grep zker"

# 2. 数据库备份完成
ssh pre-db-server "ls -lh /backups/latest/"

# 3. 监控和告警配置完成
curl http://pre-prometheus:9090/api/v1/alerts | jq '.data.alerts[] | select(.labels.alertname="ZKERCanaryAlert")'

# 4. 回滚演练完成
./scripts/rollback_drill.sh --environment pre
```

#### 0.2 灰度配置

**功能开关配置**:
```yaml
# config/features.yaml
features:
  intelligent_routing:
    enabled: true
    rollout_strategy: "user_id"
    whitelist:
      - "user_employee_*"
      - "user_beta_tester_*"
    percentage: 0  # 初始 0%

  rbac_permission:
    enabled: true
    rollout_strategy: "user_attribute"
    attribute: "registration_date"
    condition: "> 2024-01-01"
    percentage: 5  # 5% 用户

  multi_tenant:
    enabled: false  # 暂不启用
    rollout_strategy: "none"
```

**路由规则配置** (Nginx):
```nginx
# nginx-canary.conf
upstream zker_stable {
    server 10.0.1.10:8001 weight=100;
    server 10.0.1.11:8001 weight=100;
    server 10.0.1.12:8001 weight=100;
}

upstream zker_canary {
    server 10.0.2.10:8001 weight=100;
    server 10.0.2.11:8001 weight=100;
}

map $http_x_user_id $is_canary_user {
    default 0;
    ~^user_employee_ 1;  # 内部员工
    ~^user_beta_ 1;      # Beta 测试者
}

server {
    listen 80;
    server_name api.zker.com;

    location / {
        if ($is_canary_user = 1) {
            proxy_pass http://zker_canary;
            break;
        }
        proxy_pass http://zker_stable;
    }
}
```

**流量分配配置** (Istio VirtualService):
```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: zker-api
spec:
  http:
  - match:
    - headers:
        x-user-group:
          exact: employee
    route:
    - destination:
        host: zker-canary
      weight: 100  # 100% 员工流量到灰度版本

  - route:
    - destination:
        host: zker-stable
      weight: 95  # 95% 流量到稳定版本
    - destination:
        host: zker-canary
      weight: 5   # 5% 流量到灰度版本
```

---

### 阶段 1: 内部用户灰度（1-2 天）

#### 1.1 发布准备

**时间**: T-1 天 18:00（非工作时间）

**发布前检查清单**:
```bash
# 1. 确认所有团队成员在线
slack notify "发布团队，今晚 20:00 开始内部灰度发布，请确认在线"

# 2. 确认监控面板就绪
open https://grafana.zker.com/d/canary-monitor

# 3. 确认告警通道畅通
pagerduty test --team "release-team"

# 4. 确认回滚准备就绪
./scripts/rollback_check.sh --dry-run
```

#### 1.2 执行发布

**步骤 1: 数据库迁移** (T-30 分钟)
```bash
# 1. 备份数据库
mysqldump -h pre-db -u root -p zker_pre > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. 执行迁移脚本（非阻塞）
mysql -h pre-db -u root -p zker_pre < migrations/add_tenant_id_column.sql

# 3. 验证迁移结果
mysql -h pre-db -u root -p zker_pre -e "DESCRIBE users;"

# 4. 填充默认数据
mysql -h pre-db -u root -p zker_pre < migrations/fill_default_tenant_id.sql
```

**步骤 2: 部署灰度版本** (T-15 分钟)
```bash
# 1. 构建 Docker 镜像
docker build -t zker/api-server:v2.0.0-canary .

# 2. 推送到镜像仓库
docker push zker/api-server:v2.0.0-canary

# 3. 部署到灰度节点
kubectl apply -f k8s/canary-deployment.yaml

# 4. 等待 Pod 就绪
kubectl rollout status deployment/zker-canary

# 5. 验证部署
kubectl get pods -l app=zker,version=canary
```

**步骤 3: 启用功能开关** (T-5 分钟)
```bash
# 1. 更新功能开关配置
kubectl edit configmap feature-flags

# 2. 重启灰度 Pod 以加载新配置
kubectl rollout restart deployment/zker-canary

# 3. 验证功能开关生效
curl -H "X-User-ID: user_employee_001" http://canary-api:8001/api/v1/features | jq .
```

**步骤 4: 切换内部用户流量** (T-0 分钟)
```bash
# 1. 更新 Nginx 配置
kubectl apply -f k8s/nginx-canary-rule.yaml

# 2. 重新加载 Nginx
kubectl exec -it $(kubectl get pod -l app=nginx -o jsonpath='{.items[0].metadata.name}') -- nginx -s reload

# 3. 验证路由规则
curl -H "X-User-ID: user_employee_001" http://api.zker.com/api/v1/health | jq '.version'
# 预期输出: "2.0.0-canary"
```

#### 1.3 监控与验证

**关键指标监控**:
```yaml
# Prometheus 查询
# 1. 错误率
rate(http_requests_total{status=~"5..", version="canary"}[5m]) / rate(http_requests_total{version="canary"}[5m])

# 2. 响应时间
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket{version="canary"}[5m]))

# 3. QPS
rate(http_requests_total{version="canary"}[5m])

# 4. CPU 使用率
rate(process_cpu_seconds_total{version="canary"}[5m])

# 5. 内存使用率
process_resident_memory_bytes{version="canary"}
```

**功能验证清单**:
```bash
# 1. 用户登录
./scripts/test_api.sh --endpoint http://canary-api:8001 --test login

# 2. 创建 Bot
./scripts/test_api.sh --endpoint http://canary-api:8001 --test create_bot

# 3. 发送对话
./scripts/test_api.sh --endpoint http://canary-api:8001 --test send_message

# 4. 知识库检索
./scripts/test_api.sh --endpoint http://canary-api:8001 --test knowledge_search

# 5. 权限校验
./scripts/test_api.sh --endpoint http://canary-api:8001 --test permission_check
```

**告警规则**:
```yaml
groups:
  - name: canary_alerts
    rules:
      - alert: CanaryErrorRateHigh
        expr: |
          rate(http_requests_total{version="canary",status=~"5.."}[5m])
          /
          rate(http_requests_total{version="canary"}[5m]) > 0.01
        for: 5m
        annotations:
          summary: "灰度版本错误率过高"
          description: "5 分钟内错误率 > 1%"

      - alert: CanaryResponseTimeSlow
        expr: |
          histogram_quantile(0.95,
            rate(http_request_duration_seconds_bucket{version="canary"}[5m])
          ) > 2
        for: 5m
        annotations:
          summary: "灰度版本响应时间过慢"
          description: "P95 响应时间 > 2s"

      - alert: CanaryQPSDrop
        expr: |
          rate(http_requests_total{version="canary"}[5m]) < 10
        for: 10m
        annotations:
          summary: "灰度版本 QPS 异常下降"
          description: "QPS < 10，可能存在服务异常"
```

#### 1.4 决策点

**继续下一阶段的条件** (全部满足):
- ✅ 错误率 < 0.1%
- ✅ P95 响应时间 < 稳定版本 × 1.2
- ✅ 无 P0 级 Bug
- ✅ 核心功能验证通过
- ✅ 无性能严重退化

**回滚条件** (任一满足):
- ❌ 错误率 > 1%
- ❌ P95 响应时间 > 稳定版本 × 2
- ❌ 发现 P0 级 Bug
- ❌ 核心功能不可用
- ❌ 数据不一致或丢失

**延长观察条件**:
- ⚠️ 错误率 0.1%-1%
- ⚠️ P95 响应时间介于稳定版本 × 1.2-2
- ⚠️ 发现 P1 级 Bug（可快速修复）

---

### 阶段 2: 早期采用者灰度（2-3 天）

#### 2.1 用户筛选

**筛选条件**:
```sql
-- 选择早期采用者
SELECT
  user_id,
  username,
  email,
  registration_date,
  login_count_last_30d,
  bot_count
FROM users
WHERE
  -- 注册时间 > 6 个月
  registration_date < DATE_SUB(CURDATE(), INTERVAL 6 MONTH)
  -- 最近 30 天活跃
  AND last_login_at > DATE_SUB(NOW(), INTERVAL 30 DAY)
  -- 登录次数 > 20 次（高活跃）
  AND login_count_last_30d > 20
  -- 创建了 Bot（深度用户）
  AND bot_count > 0
  -- 不是内部员工（避免重复）
  AND user_id NOT LIKE 'user_employee_%'
  -- 最近 7 天无严重投诉
  AND NOT EXISTS (
    SELECT 1 FROM user_complaints c
    WHERE c.user_id = users.user_id
      AND c.complaint_date > DATE_SUB(NOW(), INTERVAL 7 DAY)
      AND c.severity = 'high'
  )
ORDER BY
  login_count_last_30d DESC,
  bot_count DESC
LIMIT 50000;  -- 5 万用户（约占 5%）
```

**导入白名单**:
```bash
# 1. 导出用户 ID 列表
mysql -h pre-db -u root -p zker_pre -e "
  SELECT user_id FROM users WHERE is_early_adopter = 1
" > early_adopters.txt

# 2. 导入 Redis
cat early_adopters.txt | while read user_id; do
  redis-cli SADD canary:whitelist $user_id
done

# 3. 验证白名单
redis-cli SCARD canary:whitelist
redis-cli SRANDMEMBER canary:whitelist 10
```

#### 2.2 流量切换

**更新路由规则**:
```nginx
# nginx-canary.conf（阶段 2）
map $http_x_user_id $is_canary_user {
    default 0;

    # 内部员工（100%）
    ~^user_employee_ 1;

    # Beta 测试者（100%）
    ~^user_beta_ 1;

    # 早期采用者（从 Redis 动态查询）
    default redis_canary_check;
}

# Redis 查询函数（通过 Lua 脚本实现）
location / {
    rewrite_by_lua_block {
        local redis = require "resty.redis"
        local red = redis:new()

        red.connect(red, "redis", 6379)

        local user_id = ngx.var.http_x_user_id
        local is_canary, err = red:sismember("canary:whitelist", user_id)

        if is_canary == 1 then
            ngx.var.upstream_name = "zker_canary"
        else
            ngx.var.upstream_name = "zker_stable"
        end
    }

    proxy_pass http://$upstream_name;
}
```

**Istio VirtualService 更新**:
```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: zker-api
spec:
  http:
  # 内部员工 → 灰度版本
  - match:
    - headers:
        x-user-group:
          exact: employee
    route:
    - destination:
        host: zker-canary
      weight: 100

  # Beta 测试者 → 灰度版本
  - match:
    - headers:
        x-user-group:
          exact: beta_tester
    route:
    - destination:
        host: zker-canary
      weight: 100

  # 早期采用者（基于 HTTP Header 标识）
  - match:
    - headers:
        x-early-adopter:
          exact: "true"
    route:
    - destination:
        host: zker-canary
      weight: 100

  # 其他用户 → 稳定版本
  - route:
    - destination:
        host: zker-stable
      weight: 100
```

#### 2.3 监控与对比

**A/B 对比分析**:
```sql
-- 对比稳定版本和灰度版本的用户行为
SELECT
  version,
  COUNT(DISTINCT user_id) as unique_users,
  COUNT(*) as total_requests,
  AVG(response_time) as avg_response_time,
  PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY response_time) as p95_response_time,
  SUM(CASE WHEN status_code >= 500 THEN 1 ELSE 0 END) / COUNT(*) as error_rate,
  SUM(CASE WHEN is_bot_created THEN 1 ELSE 0 END) as bots_created,
  SUM(CASE WHEN is_message_sent THEN 1 ELSE 0 END) as messages_sent
FROM api_logs
WHERE created_at >= NOW() - INTERVAL 1 DAY
  AND version IN ('stable', 'canary')
GROUP BY version;
```

**用户反馈收集**:
```javascript
// 前端反馈组件
const CanaryFeedback = () => {
  const [feedback, setFeedback] = useState({
    rating: 0,
    comments: '',
    issues: []
  });

  const submitFeedback = async () => {
    await fetch('/api/v1/feedback/canary', {
      method: 'POST',
      body: JSON.stringify({
        version: 'canary',
        ...feedback
      })
    });
  };

  return (
    <Modal visible={isCanaryVersion()}>
      <h3>您正在使用新版本</h3>
      <Rate value={feedback.rating} onChange={val => setFeedback({...feedback, rating: val})} />
      <TextArea
        placeholder="请告诉我们您的使用体验..."
        value={feedback.comments}
        onChange={e => setFeedback({...feedback, comments: e.target.value})}
      />
      <Checkbox.Group
        options={['性能问题', '功能异常', '界面问题', '其他']}
        onChange={vals => setFeedback({...feedback, issues: vals})}
      />
      <Button onClick={submitFeedback}>提交反馈</Button>
    </Modal>
  );
};
```

**反馈分析 Dashboard** (Grafana):
```json
{
  "dashboard": {
    "title": "灰度版本用户反馈",
    "panels": [
      {
        "title": "反馈数量趋势",
        "targets": [
          {
            "expr": "rate(canary_feedback_total[1h])"
          }
        ]
      },
      {
        "title": "平均评分",
        "targets": [
          {
            "expr": "avg(canary_feedback_rating)"
          }
        ]
      },
      {
        "title": "问题分类统计",
        "targets": [
          {
            "expr": "count by (issue_type) (canary_feedback_total{issue_type!=\"\"})"
          }
        ]
      },
      {
        "title": "负面反馈趋势",
        "targets": [
          {
            "expr": "rate(canary_feedback_total{rating<3}[1h])"
          }
        ]
      }
    ]
  }
}
```

#### 2.4 决策点

**继续下一阶段的条件** (全部满足):
- ✅ 错误率 < 0.05%
- ✅ P95 响应时间 ≤ 稳定版本
- ✅ 用户反馈评分 ≥ 4.0/5.0
- ✅ 负面反馈率 < 5%
- ✅ 核心功能正常

**回滚条件** (任一满足):
- ❌ 错误率 > 0.5%
- ❌ P95 响应时间 > 稳定版本 × 1.5
- ❌ 用户反馈评分 < 3.0/5.0
- ❌ 负面反馈率 > 20%
- ❌ 出现数据安全问题

**延长观察条件**:
- ⚠️ 错误率 0.05%-0.5%
- ⚠️ 用户反馈评分 3.0-4.0
- ⚠️ 发现性能优化空间

---

### 阶段 3: 地域灰度（3-5 天）

#### 3.1 地域选择

**选择原则**:
1. 非核心市场（业务影响小）
2. 用户量适中（便于监控）
3. 基础设施完善（支持灰度部署）
4. 法律法规兼容（无特殊合规要求）

**推荐地域**:
```
第一批（20% 流量）:
  - 亚太地区: 新加坡、马来西亚
  - 原因: 时区相近，基础设施完善，用户量适中

第二批（50% 流量）:
  - 欧洲: 德国、荷兰
  - 原因: GDPR 合规，用户接受度高

第三批（100% 流量）:
  - 北美: 美国（非核心州）、加拿大
  - 原因: 用户基数大，最后验证
```

#### 3.2 地域路由配置

**基于 GeoIP 的路由**:
```nginx
# nginx-geoip.conf
geoip_country /usr/share/GeoIP/GeoIP.dat;

map $geoip_country_code $is_canary_region {
    default 0;
    SG 1;  # 新加坡
    MY 1;  # 马来西亚
    DE 1;  # 德国
    NL 1;  # 荷兰
}

map $is_canary_region $upstream_name {
    1 zker_canary;
    default zker_stable;
}

server {
    listen 80;
    server_name api.zker.com;

    location / {
        proxy_pass http://$upstream_name;
    }
}
```

**基于 Cloudflare Workers 的路由**:
```javascript
// cloudflare-worker.js
const CANARY_REGIONS = ['SG', 'MY', 'DE', 'NL'];
const STABLE_ORIGIN = 'https://api-stable.zker.com';
const CANARY_ORIGIN = 'https://api-canary.zker.com';

addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
});

async function handleRequest(request) {
  const country = request.cf.country;

  const origin = CANARY_REGIONS.includes(country)
    ? CANARY_ORIGIN
    : STABLE_ORIGIN;

  // 添加地域标识 Header
  const modifiedRequest = new Request(request, {
    headers: {
      ...request.headers,
      'X-User-Country': country,
      'X-Is-Canary': CANARY_REGIONS.includes(country).toString()
    }
  });

  return fetch(origin, modifiedRequest);
}
```

#### 3.3 地域监控

**地域特定指标**:
```yaml
# Prometheus 查询（按地域分组）
# 1. 各地域 QPS
sum by (country_code) (rate(http_requests_total{version="canary"}[5m]))

# 2. 各地域错误率
sum by (country_code) (rate(http_requests_total{version="canary",status=~"5.."}[5m]))
/
sum by (country_code) (rate(http_requests_total{version="canary"}[5m]))

# 3. 各地域响应时间
histogram_quantile(0.95,
  sum by (country_code, le) (
    rate(http_request_duration_seconds_bucket{version="canary"}[5m])
  )
)

# 4. 各地域活跃用户数
count by (country_code) (
  rate(user_login_total{version="canary"}[5m]) > 0
)
```

**Grafana Dashboard**:
```json
{
  "dashboard": {
    "title": "地域灰度监控",
    "panels": [
      {
        "title": "各国 QPS 对比",
        "type": "graph",
        "targets": [
          {
            "expr": "sum by (country_code) (rate(http_requests_total{version=\"canary\"}[5m]))"
          }
        ]
      },
      {
        "title": "各国错误率对比",
        "type": "heatmap",
        "targets": [
          {
            "expr": "sum by (country_code) (rate(http_requests_total{version=\"canary\",status=~\"5..\"}[5m])) / sum by (country_code) (rate(http_requests_total{version=\"canary\"}[5m]))"
          }
        ]
      },
      {
        "title": "各国响应时间对比",
        "type": "table",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, sum by (country_code, le) (rate(http_request_duration_seconds_bucket{version=\"canary\"}[5m])))"
          }
        ]
      }
    ]
  }
}
```

#### 3.4 决策点

**继续下一阶段的条件** (全部满足):
- ✅ 所有灰度地域错误率 < 0.05%
- ✅ 所有灰度地域 P95 响应时间 ≤ 稳定版本
- ✅ 无地域特定 Bug
- ✅ 用户反馈评分 ≥ 4.2/5.0

**回滚条件** (任一满足):
- ❌ 任何地域错误率 > 1%
- ❌ 发现地域特定严重 Bug
- ❌ 出现合规性问题

**部分回滚条件**:
- ⚠️ 特定地域出现问题
- ⚠️ 回滚该地域到稳定版本，其他地域继续

---

### 阶段 4: 流量比例灰度（5-7 天）

#### 4.1 逐步放量

**放量计划**:
```
Day 1: 50% 流量
Day 2: 50% 流量（观察日）
Day 3: 80% 流量
Day 4: 80% 流量（观察日）
Day 5: 100% 流量
Day 6-7: 100% 流量（最终观察）
```

**Istio VirtualService 配置**:
```yaml
# 50% 流量
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: zker-api
spec:
  http:
  - route:
    - destination:
        host: zker-stable
      weight: 50
    - destination:
        host: zker-canary
      weight: 50

---
# 80% 流量
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: zker-api
spec:
  http:
  - route:
    - destination:
        host: zker-stable
      weight: 20
    - destination:
        host: zker-canary
      weight: 80

---
# 100% 流量
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: zker-api
spec:
  http:
  - route:
    - destination:
        host: zker-canary
      weight: 100
```

**自动放量脚本**:
```bash
#!/bin/bash
# auto-rollout.sh

STAGES=(
  "50:1d"
  "50:1d"
  "80:1d"
  "80:1d"
  "100:2d"
)

for stage in "${STAGES[@]}"; do
  PERCENTAGE=$(echo $stage | cut -d: -f1)
  DURATION=$(echo $stage | cut -d: -f2)

  echo "切换到 ${PERCENTAGE}% 灰度流量，持续 ${DURATION}"

  # 更新 Istio VirtualService
  kubectl apply -f k8s/virtualservice-${PERCENTAGE}.yaml

  # 等待观察期
  sleep ${DURATION}

  # 检查健康状态
  if ! ./scripts/check_canary_health.sh; then
    echo "健康检查失败，停止自动放量"
    exit 1
  fi

  echo "${PERCENTAGE}% 流量阶段完成"
done

echo "灰度发布完成！"
```

#### 4.2 全量监控

**核心指标大盘**:
```json
{
  "dashboard": {
    "title": "灰度发布核心指标",
    "panels": [
      {
        "title": "总 QPS",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{version=\"canary\"}[5m]))",
            "legendFormat": "灰度版本"
          },
          {
            "expr": "sum(rate(http_requests_total{version=\"stable\"}[5m]))",
            "legendFormat": "稳定版本"
          }
        ]
      },
      {
        "title": "错误率对比",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{version=\"canary\",status=~\"5..\"}[5m])) / sum(rate(http_requests_total{version=\"canary\"}[5m]))",
            "legendFormat": "灰度版本"
          },
          {
            "expr": "sum(rate(http_requests_total{version=\"stable\",status=~\"5..\"}[5m])) / sum(rate(http_requests_total{version=\"stable\"}[5m]))",
            "legendFormat": "稳定版本"
          }
        ]
      },
      {
        "title": "P95 响应时间对比",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{version=\"canary\"}[5m])))",
            "legendFormat": "灰度版本"
          },
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{version=\"stable\"}[5m])))",
            "legendFormat": "稳定版本"
          }
        ]
      },
      {
        "title": "CPU 使用率对比",
        "targets": [
          {
            "expr": "avg(rate(process_cpu_seconds_total{version=\"canary\"}[5m])) * 100",
            "legendFormat": "灰度版本"
          },
          {
            "expr": "avg(rate(process_cpu_seconds_total{version=\"stable\"}[5m])) * 100",
            "legendFormat": "稳定版本"
          }
        ]
      },
      {
        "title": "内存使用率对比",
        "targets": [
          {
            "expr": "avg(process_resident_memory_bytes{version=\"canary\"}) / 1024 / 1024 / 1024",
            "legendFormat": "灰度版本 (GB)"
          },
          {
            "expr": "avg(process_resident_memory_bytes{version=\"stable\"}) / 1024 / 1024 / 1024",
            "legendFormat": "稳定版本 (GB)"
          }
        ]
      }
    ]
  }
}
```

#### 4.3 决策点

**完成灰度发布的条件** (全部满足):
- ✅ 100% 流量运行 48 小时无异常
- ✅ 错误率 ≤ 稳定版本
- ✅ P95 响应时间 ≤ 稳定版本
- ✅ 无 P0/P1 级 Bug
- ✅ 用户反馈评分 ≥ 4.5/5.0

**回滚条件** (任一满足):
- ❌ 错误率 > 稳定版本 × 2
- ❌ P95 响应时间 > 稳定版本 × 1.5
- ❌ 发现严重 Bug
- ❌ 出现数据安全问题

---

## 🔄 回滚方案

### 回滚触发条件

**自动回滚** (无需人工干预):
```yaml
# Prometheus 告警规则（自动触发回滚）
groups:
  - name: auto_rollback_rules
    rules:
      - alert: AutoRollbackCriticalError
        expr: |
          rate(http_requests_total{version="canary",status="500"}[2m]) > 10
        for: 2m
        annotations:
          summary: "自动回滚：严重错误"
          description: "2 分钟内 500 错误 > 10/秒"

      - alert: AutoRollbackDataLoss
        expr: |
          abs(
            rate(database_transactions_total{version="canary"}[5m]) -
            rate(database_transactions_total{version="stable"}[5m])
          ) > 100
        for: 5m
        annotations:
          summary: "自动回滚：数据不一致"
          description: "5 分钟内数据库事务数差异 > 100/秒"
```

**人工回滚** (需要决策):
- ⚠️ 错误率持续上升但未达到自动回滚阈值
- ⚠️ 用户反馈评分 < 3.0/5.0
- ⚠️ 性能严重下降（QPS < 基线 50%）
- ⚠️ 发现 P1 级 Bug（可快速修复）

### 回滚步骤

**步骤 1: 立即切换流量** (1 分钟内)
```bash
# 1. 更新 Istio VirtualService
kubectl apply -f k8s/virtualservice-rollback.yaml

# 验证流量已切换
kubectl get virtualservice zker-api -o yaml | grep weight

# 2. 或者通过命令行快速回滚
kubectl patch virtualservice zker-api --type json -p '
  [
    {
      "op": "replace",
      "path": "/spec/http/0/route/0/weight",
      "value": 100
    },
    {
      "op": "replace",
      "path": "/spec/http/0/route/1/weight",
      "value": 0
    }
  ]
'
```

**步骤 2: 停止灰度版本** (2 分钟内)
```bash
# 1. 缩容灰度 Pod 到 0
kubectl scale deployment zker-canary --replicas=0

# 2. 验证灰度 Pod 已停止
kubectl get pods -l app=zker,version=canary

# 3. （可选）删除灰度部署
kubectl delete deployment zker-canary
```

**步骤 3: 数据库回滚** (5 分钟内)
```bash
# 1. 停止应用服务（如果数据库有破坏性变更）
kubectl scale deployment zker-stable --replicas=0

# 2. 执行回滚脚本
mysql -h db -u root -p zker_production < migrations/rollback/add_tenant_id_column.sql

# 3. 验证回滚结果
mysql -h db -u root -p zker_production -e "DESCRIBE users;"

# 4. 重启应用服务
kubectl scale deployment zker-stable --replicas=10

# 5. 验证服务恢复
curl http://api.zker.com/health
```

**步骤 4: 通知与复盘** (30 分钟内)
```bash
# 1. 发送回滚通知
slack notify "灰度版本 v2.0.0-canary 已回滚，原因：[填写原因]"

# 2. 创建复盘 Issue
jira create --summary "灰度发布回滚复盘 v2.0.0-canary" \
  --description "回滚原因、影响范围、后续改进措施"

# 3. 安排复盘会议
calendar schedule --meeting "灰度发布回滚复盘" --attendees @release-team --duration 1h
```

### 回滚验证

**验证清单**:
```bash
# 1. 流量已全部切回稳定版本
kubectl get virtualservice zker-api -o yaml | grep -A 10 "weight: 100"

# 2. 灰度 Pod 已停止
kubectl get pods -l app=zker,version=canary
# 预期输出: No resources found.

# 3. 稳定版本健康
curl http://api.zker.com/health | jq '.status'
# 预期输出: "healthy"

# 4. 核心功能正常
./scripts/smoke_test.sh --endpoint http://api.zker.com

# 5. 监控指标恢复
# 错误率下降
# QPS 恢复正常
# 响应时间恢复基线
```

---

## 📊 监控与告警

### 监控体系

**4 层监控金字塔**:
```
1. 基础设施层（Infrastructure）
   - 服务器 CPU、内存、磁盘、网络
   - 数据库连接数、慢查询
   - Redis 命中率、内存使用

2. 平台层（Platform）
   - Kubernetes Pod 状态
   - 容器资源使用
   - 服务网格指标

3. 应用层（Application）
   - API QPS、错误率、响应时间
   - 业务指标（Bot 创建数、对话数）
   - 自定义业务指标

4. 用户体验层（User Experience）
   - 页面加载时间
   - 功能使用率
   - 用户反馈评分
```

### 核心监控指标

**黄金指标** (Golden Signals):
```yaml
# 1. 延迟 (Latency)
指标: HTTP 请求响应时间
查询: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
告警阈值: P95 > 2s

# 2. 流量 (Traffic)
指标: HTTP 请求 QPS
查询: rate(http_requests_total[5m])
告警阈值: QPS < 10 (持续 10 分钟)

# 3. 错误 (Errors)
指标: HTTP 5xx 错误率
查询: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])
告警阈值: 错误率 > 1%

# 4. 饱和度 (Saturation)
指标: CPU 使用率、内存使用率
查询: (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m]))) * 100
告警阈值: CPU > 80%, 内存 > 85%
```

### Grafana Dashboard

**灰度发布总览 Dashboard**:
```json
{
  "dashboard": {
    "title": "ZKER 灰度发布监控",
    "tags": ["canary", "release"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "流量分配",
        "type": "stat",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{version=\"canary\"}[5m])) / sum(rate(http_requests_total[5m])) * 100",
            "legendFormat": "灰度流量占比"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "unit": "percent",
            "min": 0,
            "max": 100,
            "thresholds": {
              "steps": [
                {"color": "green", "value": 0},
                {"color": "yellow", "value": 50},
                {"color": "red", "value": 90}
              ]
            }
          }
        }
      },
      {
        "id": 2,
        "title": "QPS 对比",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{version=\"canary\"}[5m]))",
            "legendFormat": "灰度版本"
          },
          {
            "expr": "sum(rate(http_requests_total{version=\"stable\"}[5m]))",
            "legendFormat": "稳定版本"
          }
        ]
      },
      {
        "id": 3,
        "title": "错误率对比",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{version=\"canary\",status=~\"5..\"}[5m])) / sum(rate(http_requests_total{version=\"canary\"}[5m]))",
            "legendFormat": "灰度版本"
          },
          {
            "expr": "sum(rate(http_requests_total{version=\"stable\",status=~\"5..\"}[5m])) / sum(rate(http_requests_total{version=\"stable\"}[5m]))",
            "legendFormat": "稳定版本"
          }
        ],
        "alert": {
          "conditions": [
            {
              "evaluator": {
                "params": [0.01],
                "type": "gt"
              },
              "operator": {
                "type": "and"
              },
              "query": {
                "params": ["A", "5m", "now"]
              },
              "reducer": {
                "params": [],
                "type": "avg"
              },
              "type": "query"
            }
          ]
        }
      },
      {
        "id": 4,
        "title": "P95 响应时间对比",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{version=\"canary\"}[5m])))",
            "legendFormat": "灰度版本"
          },
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{version=\"stable\"}[5m])))",
            "legendFormat": "稳定版本"
          }
        ]
      },
      {
        "id": 5,
        "title": "资源使用率对比",
        "type": "graph",
        "targets": [
          {
            "expr": "avg(rate(process_cpu_seconds_total{version=\"canary\"}[5m])) * 100",
            "legendFormat": "灰度版本 CPU (%)"
          },
          {
            "expr": "avg(rate(process_cpu_seconds_total{version=\"stable\"}[5m])) * 100",
            "legendFormat": "稳定版本 CPU (%)"
          }
        ]
      },
      {
        "id": 6,
        "title": "用户反馈评分",
        "type": "gauge",
        "targets": [
          {
            "expr": "avg(canary_feedback_rating)",
            "legendFormat": "平均评分"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "unit": "short",
            "min": 1,
            "max": 5,
            "thresholds": {
              "steps": [
                {"color": "red", "value": 0},
                {"color": "yellow", "value": 3},
                {"color": "green", "value": 4}
              ]
            }
          }
        }
      },
      {
        "id": 7,
        "title": "业务指标对比",
        "type": "table",
        "targets": [
          {
            "expr": "sum by (version) (rate(bots_created_total[1h]))",
            "legendFormat": "Bot 创建数/小时"
          },
          {
            "expr": "sum by (version) (rate(messages_sent_total[1h]))",
            "legendFormat": "消息发送数/小时"
          },
          {
            "expr": "sum by (version) (rate(knowledge_searches_total[1h]))",
            "legendFormat": "知识库检索数/小时"
          }
        ],
        "transformations": [
          {
            "id": "organize",
            "options": {
              "excludeByName": {},
              "indexByName": {},
              "renameByName": {}
            }
          }
        ]
      }
    ]
  }
}
```

### 告警配置

**PagerDuty 集成**:
```yaml
# Prometheus AlertManager 配置
route:
  receiver: 'pagerduty'
  group_by: ['alertname', 'version']
  group_wait: 10s
  group_interval: 5m
  repeat_interval: 3h

receivers:
  - name: 'pagerduty'
    pagerduty_configs:
      - service_key: '<PAGERDUTY_SERVICE_KEY>'
        description: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'

inhibit_rules:
  # 如果灰度版本已经回滚，抑制灰度版本的其他告警
  - source_match:
      alertname: 'CanaryRolledBack'
    target_match_re:
      version: 'canary'
    equal: ['alertname']
```

**Slack 集成**:
```yaml
receivers:
  - name: 'slack'
    slack_configs:
      - api_url: '<SLACK_WEBHOOK_URL>'
        channel: '#release-ops'
        title: '🚨 灰度发布告警'
        text: |
          *告警名称*: {{ .CommonLabels.alertname }}
          *严重程度*: {{ .CommonLabels.severity }}
          *版本*: {{ .CommonLabels.version }}
          *摘要*: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}
          *描述*: {{ range .Alerts }}{{ .Annotations.description }}{{ end }}
          *查看详情*: {{ .ExternalURL }}
        actions:
          - type: button
            text: '查看 Grafana'
            url: '{{ .ExternalURL }}'
          - type: button
            text: '执行回滚'
            url: 'https://zker.com/rollback?alert={{ .GroupLabels.alertname }}'
```

---

## 📚 附录

### A. 灰度发布检查清单

**发布前检查**:
- [ ] 代码审查完成
- [ ] 所有测试通过
- [ ] CHANGELOG 更新
- [ ] 数据库迁移脚本准备
- [ ] 回滚脚本准备
- [ ] 监控 Dashboard 配置
- [ ] 告警规则配置
- [ ] 灰度环境部署
- [ ] 功能开关配置
- [ ] 团队培训完成

**阶段 1 检查** (内部用户灰度):
- [ ] 数据库迁移完成
- [ ] 灰度版本部署成功
- [ ] 功能开关启用
- [ ] 内部用户流量切换
- [ ] 核心功能验证通过
- [ ] 监控指标正常
- [ ] 无 P0 级 Bug

**阶段 2 检查** (早期采用者灰度):
- [ ] 用户白名单配置
- [ ] 流量切换成功
- [ ] A/B 对比分析完成
- [ ] 用户反馈收集正常
- [ ] 无严重性能问题
- [ ] 无数据安全问题

**阶段 3 检查** (地域灰度):
- [ ] 地域路由配置
- [ ] 各地域监控正常
- [ ] 无地域特定 Bug
- [ ] 合规性验证通过

**阶段 4 检查** (流量比例灰度):
- [ ] 50% 流量运行正常
- [ ] 80% 流量运行正常
- [ ] 100% 流量运行 48 小时
- [ ] 核心指标达标
- [ ] 用户反馈评分 ≥ 4.5

**完成检查**:
- [ ] 灰度版本成为稳定版本
- [ ] 旧版本下线
- [ ] 功能开关移除
- [ ] 文档更新
- [ ] 复盘会议完成

### B. 常用命令速查

```bash
# 灰度发布相关
# 查看当前灰度流量占比
kubectl get virtualservice zker-api -o yaml | grep -A 5 "weight:"

# 切换流量（快速回滚）
kubectl patch virtualservice zker-api --type json -p '[{"op": "replace", "path": "/spec/http/0/route/0/weight", "value": 100}]'

# 查看灰度 Pod 状态
kubectl get pods -l app=zker,version=canary

# 查看灰度版本日志
kubectl logs -f -l app=zker,version=canary

# 监控相关
# 查询灰度版本错误率
curl -s 'http://prometheus:9090/api/v1/query?query=rate(http_requests_total{version="canary",status=~"5.."}[5m])' | jq

# 查询灰度版本 P95 响应时间
curl -s 'http://prometheus:9090/api/v1/query?query=histogram_quantile(0.95,rate(http_request_duration_seconds_bucket{version="canary"}[5m]))' | jq

# 功能开关相关
# 查看当前功能开关状态
curl http://api.zker.com/api/v1/features | jq

# 启用功能开关
curl -X POST http://api.zker.com/api/v1/features/intelligent_routing/enable -d '{"percentage": 50}'

# 禁用功能开关
curl -X POST http://api.zker.com/api/v1/features/intelligent_routing/disable
```

### C. 联系方式

| 角色 | 姓名 | 联系方式 | 职责 |
|-----|------|---------|------|
| 发布总负责人 | [待填写] | [待填写] | 总体协调、最终决策 |
| 开发负责人 | [待填写] | [待填写] | 代码质量、Bug 修复 |
| 测试负责人 | [待填写] | [待填写] | 功能验证、用户反馈 |
| 运维负责人 | [待填写] | [待填写] | 部署、监控、回滚 |
| 产品负责人 | [待填写] | [待填写] | 用户体验、业务影响 |

---

**文档变更历史**:

| 版本 | 日期 | 变更内容 | 作者 |
|-----|------|---------|------|
| v1.0 | 2025-01-01 | 初始版本 | 发布管理团队 |

**审批记录**:

| 角色 | 姓名 | 审批意见 | 日期 |
|-----|------|---------|------|
| 技术架构委员会 | [待填写] | [待审批] | [待审批] |
| 发布负责人 | [待填写] | [待审批] | [待审批] |
| 项目经理 | [待填写] | [待审批] | [待审批] |

---

**© 2025 ZKER Project. All rights reserved.**
