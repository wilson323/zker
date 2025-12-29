# ZKER 智能路由与权限管理系统 - 一键回滚方案

**文档版本**: v1.0
**创建日期**: 2025-01-01
**最后更新**: 2025-01-01
**负责人**: 运维团队
**审批人**: 技术架构委员会

---

## 📋 文档概述

### 回滚目标

本一键回滚方案旨在实现系统发布后的快速、安全回滚，确保：

**核心目标**:
- ✅ **快速响应**: 检测到问题后 5 分钟内完成回滚
- ✅ **自动化**: 一键触发，全自动执行，无需人工干预
- ✅ **安全可靠**: 回滚过程数据不丢失，服务不中断
- ✅ **可追溯**: 完整记录回滚原因、过程、结果
- ✅ **可验证**: 回滚后自动验证系统功能正常

### 回滚类型

**1. 应用回滚** (Application Rollback)
```
范围: 应用代码变更
回滚对象: Docker 镜像、Kubernetes Deployment
典型场景: 代码 Bug、性能下降、功能异常
回滚时间: < 2 分钟
```

**2. 数据库回滚** (Database Rollback)
```
范围: 数据库架构变更
回滚对象: DDL 脚本、数据迁移
典型场景: 数据库错误、数据丢失、性能问题
回滚时间: < 10 分钟
```

**3. 配置回滚** (Configuration Rollback)
```
范围: 系统配置变更
回滚对象: ConfigMap、环境变量、功能开关
典型场景: 配置错误、参数调整不当
回滚时间: < 1 分钟
```

**4. 基础设施回滚** (Infrastructure Rollback)
```
范围: 基础设施变更
回滚对象: Kubernetes 版本、网络配置、存储配置
典型场景: K8s 升级失败、网络故障、存储问题
回滚时间: < 30 分钟
```

---

## 🎯 回滚触发条件

### 自动触发条件

**严重故障** (立即回滚):
```yaml
auto_rollback_triggers:
  critical:
    # 1. 服务完全不可用
    - name: "ServiceDown"
      condition: "up == 0"
      duration: "30s"
      action: "immediate_rollback"

    # 2. 错误率 > 10%
    - name: "HighErrorRate"
      condition: "rate(http_requests_total{status=~\"5..\"}[1m]) / rate(http_requests_total[1m]) > 0.1"
      duration: "30s"
      action: "immediate_rollback"

    # 3. P95 响应时间 > 10s
    - name: "SlowResponse"
      condition: "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[1m])) > 10"
      duration: "30s"
      action: "immediate_rollback"

    # 4. 数据丢失
    - name: "DataLoss"
      condition: "abs(database_transactions_delta) > 100"
      duration: "10s"
      action: "immediate_rollback"

    # 5. 数据库连接失败
    - name: "DatabaseConnectionFailure"
      condition: "mysql_up == 0"
      duration: "10s"
      action: "immediate_rollback"
```

**高危故障** (5 分钟内回滚):
```yaml
  high:
    # 1. 错误率 > 1%
    - name: "ElevatedErrorRate"
      condition: "rate(http_requests_total{status=~\"5..\"}[5m]) / rate(http_requests_total[5m]) > 0.01"
      duration: "5m"
      action: "rollback_5min"

    # 2. P95 响应时间 > 5s
    - name: "DegradedPerformance"
      condition: "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 5"
      duration: "5m"
      action: "rollback_5min"

    # 3. QPS 下降 > 50%
    - name: "TrafficDrop"
      condition: "rate(http_requests_total[5m]) < baseline_qps * 0.5"
      duration: "5m"
      action: "rollback_5min"

    # 4. CPU 使用率 > 95%
    - name: "HighCPUUsage"
      condition: "avg(rate(process_cpu_seconds_total[5m])) * 100 > 95"
      duration: "5m"
      action: "rollback_5min"

    # 5. 内存使用率 > 95%
    - name: "HighMemoryUsage"
      condition: "process_resident_memory_bytes / node_memory_MemTotal_bytes > 0.95"
      duration: "5m"
      action: "rollback_5min"
```

**中等故障** (人工决策):
```yaml
  medium:
    # 1. 错误率上升但未达高危阈值
    - name: "GradualErrorIncrease"
      condition: "rate(http_requests_total{status=~\"5..\"}[10m]) / rate(http_requests_total[10m]) > 0.005"
      duration: "10m"
      action: "notify_and_wait"

    # 2. 用户反馈评分 < 3.0
    - name: "PoorUserFeedback"
      condition: "avg(canary_feedback_rating) < 3.0"
      duration: "30m"
      action: "notify_and_wait"

    # 3. 性能轻微下降
    - name: "PerformanceDegradation"
      condition: "p95_response_time > baseline_p95 * 1.5"
      duration: "10m"
      action: "notify_and_wait"
```

### 人工触发条件

**运维团队触发**:
```bash
# 场景 1: 监控发现异常趋势
./scripts/rollback.sh --reason "监控发现错误率持续上升，预计 10 分钟内达到阈值"

# 场景 2: 用户严重投诉
./scripts/rollback.sh --reason "收到大量用户投诉功能异常，影响业务"

# 场景 3: 第三方服务故障
./scripts/rollback.sh --reason "依赖的第三方 AI 服务异常，回滚到兼容版本"

# 场景 4: 安全问题
./scripts/rollback.sh --reason "发现安全漏洞，需要紧急回滚修复" --severity critical
```

**开发团队触发**:
```bash
# 场景 1: 发现严重 Bug
./scripts/rollback.sh --reason "测试发现 P0 级 Bug，影响核心功能" --bug-id ZKER-1234

# 场景 2: 数据一致性问题
./scripts/rollback.sh --reason "发现数据一致性问题，可能导致数据丢失"

# 场景 3: 性能严重下降
./scripts/rollback.sh --reason "性能测试发现性能下降 > 50%，不满足上线要求"
```

---

## 🤖 自动化回滚系统

### 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                     回滚控制系统                              │
│                                                               │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐   │
│  │ Prometheus   │───▶│ AlertManager │───▶│  Webhook     │   │
│  │ (监控告警)    │    │ (告警路由)    │    │  (触发器)    │   │
│  └──────────────┘    └──────────────┘    └──────────────┘   │
│                                                    │         │
│                                                    ▼         │
│  ┌──────────────────────────────────────────────────────┐  │
│  │            回滚决策引擎 (Rollout Controller)         │  │
│  │                                                       │  │
│  │  1. 接收告警                                          │  │
│  │  2. 评估严重程度                                        │  │
│  │  3. 决定回滚策略                                        │  │
│  │  4. 执行回滚操作                                        │  │
│  │  5. 验证回滚结果                                        │  │
│  └──────────────────────────────────────────────────────┘  │
│                          │                                 │
│                          ▼                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              回滚执行引擎 (ArgoCD/Flux)               │  │
│  │                                                       │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │ 应用回滚    │  │ 数据库回滚  │  │ 配置回滚    │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
│                          │                                 │
│                          ▼                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              通知系统 (Notification)                 │  │
│  │                                                       │  │
│  │  ┌───────────┐  ┌───────────┐  ┌───────────┐       │  │
│  │  │ Slack     │  │ PagerDuty │  │ Email     │       │  │
│  │  └───────────┘  └───────────┘  └───────────┘       │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 回滚决策引擎

**决策逻辑** (Python 实现):
```python
#!/usr/bin/env python3
"""
回滚决策引擎
根据告警严重程度和历史数据自动决定回滚策略
"""

import logging
from enum import Enum
from typing import Dict, List, Optional

logger = logging.getLogger(__name__)


class RollbackSeverity(Enum):
    """回滚严重程度"""
    CRITICAL = "critical"  # 立即回滚
    HIGH = "high"          # 5 分钟内回滚
    MEDIUM = "medium"      # 需要人工决策
    LOW = "low"            # 记录观察


class RollbackDecisionEngine:
    """回滚决策引擎"""

    def __init__(self, config: Dict):
        self.config = config
        self.rollback_history = []
        self.baseline_metrics = self._load_baseline_metrics()

    def evaluate_alert(self, alert: Dict) -> RollbackSeverity:
        """
        评估告警严重程度

        Args:
          alert: 告警信息

        Returns:
          回滚严重程度
        """
        alert_name = alert.get('alertname', '')
        alert_value = alert.get('value', 0)
        alert_duration = alert.get('duration', '0s')

        logger.info(f"评估告警: {alert_name}, 值: {alert_value}, 持续时间: {alert_duration}")

        # 严重故障 → 立即回滚
        if alert_name in [
            'ServiceDown',
            'HighErrorRate',
            'SlowResponse',
            'DataLoss',
            'DatabaseConnectionFailure'
        ]:
            logger.warning(f"检测到严重故障: {alert_name}，触发立即回滚")
            return RollbackSeverity.CRITICAL

        # 高危故障 → 5 分钟内回滚
        if alert_name in [
            'ElevatedErrorRate',
            'DegradedPerformance',
            'TrafficDrop',
            'HighCPUUsage',
            'HighMemoryUsage'
        ]:
            # 检查是否持续超过阈值时间
            if self._parse_duration(alert_duration) >= 300:  # 5 分钟
                logger.warning(f"检测到高危故障: {alert_name}，触发 5 分钟回滚")
                return RollbackSeverity.HIGH
            else:
                logger.info(f"检测到高危故障: {alert_name}，尚未达到持续时间阈值，继续观察")

        # 中等故障 → 人工决策
        if alert_name in [
            'GradualErrorIncrease',
            'PoorUserFeedback',
            'PerformanceDegradation'
        ]:
            logger.info(f"检测到中等故障: {alert_name}，需要人工决策")
            return RollbackSeverity.MEDIUM

        return RollbackSeverity.LOW

    def should_rollback(self, severity: RollbackSeverity) -> bool:
        """
        判断是否应该回滚

        Args:
          severity: 严重程度

        Returns:
          是否回滚
        """
        # 严重和高危故障自动回滚
        if severity in [RollbackSeverity.CRITICAL, RollbackSeverity.HIGH]:
            return True

        # 中等故障需要人工确认
        if severity == RollbackSeverity.MEDIUM:
            # 检查是否在维护窗口
            if self._is_maintenance_window():
                logger.info("当前在维护窗口，允许人工决策回滚")
                return self._get_human_approval()

        return False

    def execute_rollback(self, severity: RollbackSeverity, alert: Dict):
        """
        执行回滚操作

        Args:
          severity: 严重程度
          alert: 告警信息
        """
        logger.info(f"开始执行回滚，严重程度: {severity.value}")

        # 1. 记录回滚原因
        rollback_record = {
            'timestamp': self._get_current_time(),
            'severity': severity.value,
            'alert': alert,
            'decision': 'automatic' if severity != RollbackSeverity.MEDIUM else 'manual'
        }
        self.rollback_history.append(rollback_record)

        # 2. 通知相关人员
        self._notify_rollback_start(rollback_record)

        # 3. 执行回滚操作
        try:
            if severity == RollbackSeverity.CRITICAL:
                self._execute_immediate_rollback()
            elif severity == RollbackSeverity.HIGH:
                self._execute_gradual_rollback()
            else:
                self._execute_manual_rollback()

            # 4. 验证回滚结果
            if self._verify_rollback():
                logger.info("回滚成功")
                self._notify_rollback_success(rollback_record)
            else:
                logger.error("回滚验证失败")
                self._notify_rollback_failure(rollback_record)

        except Exception as e:
            logger.error(f"回滚执行失败: {e}")
            self._notify_rollback_error(rollback_record, str(e))

    def _execute_immediate_rollback(self):
        """执行立即回滚"""
        logger.info("执行立即回滚...")

        # 1. 立即切换流量到上一版本
        self._switch_traffic(to_version='previous')

        # 2. 停止当前版本 Pod
        self._scale_down_pods(replicas=0)

        # 3. 验证服务恢复
        self._verify_service_health()

    def _execute_gradual_rollback(self):
        """执行渐进式回滚"""
        logger.info("执行渐进式回滚...")

        # 1. 逐步切换流量（50% → 100%）
        for percentage in [50, 100]:
            logger.info(f"切换 {percentage}% 流量到上一版本")
            self._switch_traffic(to_version='previous', percentage=percentage)

            # 等待观察
            import time
            time.sleep(60)  # 1 分钟观察期

            # 验证健康状态
            if not self._verify_service_health():
                logger.warning(f"{percentage}% 流量切换后健康检查失败，继续切换")

        # 2. 停止当前版本 Pod
        self._scale_down_pods(replicas=0)

    def _execute_manual_rollback(self):
        """执行人工回滚（需要确认）"""
        logger.info("等待人工确认回滚...")

        # 发送通知，请求人工确认
        self._request_manual_approval()

        # 等待人工操作
        #（实际实现中，这里可以暂停，等待外部 API 调用确认）

    def _switch_traffic(self, to_version: str, percentage: int = 100):
        """切换流量"""
        logger.info(f"切换 {percentage}% 流量到 {to_version} 版本")

        if to_version == 'previous':
            # 更新 Istio VirtualService
            import subprocess
            subprocess.run([
                'kubectl', 'patch', 'virtualservice', 'zker-api',
                '--type=json',
                '-p', f'[{{"op": "replace", "path": "/spec/http/0/route/0/weight", "value": {percentage}}}]'
            ])

    def _scale_down_pods(self, replicas: int):
        """缩容 Pod"""
        logger.info(f"缩容当前版本 Pod 到 {replicas}")

        import subprocess
        subprocess.run([
            'kubectl', 'scale', 'deployment', 'zker-canary',
            f'--replicas={replicas}'
        ])

    def _verify_service_health(self) -> bool:
        """验证服务健康"""
        logger.info("验证服务健康状态...")

        # 1. 检查 Pod 状态
        result = subprocess.run([
            'kubectl', 'get', 'pods', '-l', 'app=zker,version=stable'
        ], capture_output=True, text=True)

        if 'Running' not in result.stdout:
            logger.error("Pod 状态异常")
            return False

        # 2. 检查健康检查端点
        import requests
        try:
            response = requests.get('http://api.zker.com/health', timeout=5)
            if response.status_code != 200:
                logger.error(f"健康检查失败: {response.status_code}")
                return False
        except Exception as e:
            logger.error(f"健康检查请求失败: {e}")
            return False

        # 3. 检查核心指标
        #（实际实现中，这里应该查询 Prometheus，检查错误率、响应时间等）

        return True

    def _notify_rollback_start(self, record: Dict):
        """通知回滚开始"""
        logger.info(f"发送回滚开始通知: {record}")
        # 实际实现中，这里应该调用 Slack、PagerDuty 等

    def _notify_rollback_success(self, record: Dict):
        """通知回滚成功"""
        logger.info(f"发送回滚成功通知: {record}")

    def _notify_rollback_failure(self, record: Dict):
        """通知回滚失败"""
        logger.error(f"发送回滚失败通知: {record}")

    def _notify_rollback_error(self, record: Dict, error: str):
        """通知回滚异常"""
        logger.error(f"发送回滚异常通知: {record}, 错误: {error}")

    def _request_manual_approval(self):
        """请求人工确认"""
        logger.info("发送人工确认请求...")
        # 实际实现中，这里应该发送 Slack 消息或 PagerDuty 通知

    def _load_baseline_metrics(self) -> Dict:
        """加载基线指标"""
        # 实际实现中，这里应该从 Prometheus 或配置文件加载
        return {
            'baseline_qps': 1000,
            'baseline_p95': 1.5,
            'baseline_error_rate': 0.001
        }

    def _parse_duration(self, duration: str) -> int:
        """解析持续时间字符串为秒数"""
        # "5m" → 300, "30s" → 30
        if duration.endswith('s'):
            return int(duration[:-1])
        elif duration.endswith('m'):
            return int(duration[:-1]) * 60
        elif duration.endswith('h'):
            return int(duration[:-1]) * 3600
        return 0

    def _is_maintenance_window(self) -> bool:
        """检查是否在维护窗口"""
        # 实际实现中，这里应该检查当前时间是否在预定义的维护窗口
        return True

    def _get_human_approval(self) -> bool:
        """获取人工确认"""
        # 实际实现中，这里应该等待人工操作确认
        return False

    def _get_current_time(self) -> str:
        """获取当前时间"""
        from datetime import datetime
        return datetime.now().isoformat()


def main():
    """主函数"""
    import argparse
    import json

    parser = argparse.ArgumentParser(description='回滚决策引擎')
    parser.add_argument('--alert', required=True, help='告警信息 (JSON 格式)')
    parser.add_argument('--dry-run', action='store_true', help='演练模式，不执行实际回滚')

    args = parser.parse_args()

    # 解析告警信息
    alert = json.loads(args.alert)

    # 创建决策引擎
    engine = RollbackDecisionEngine(config={})

    # 评估告警
    severity = engine.evaluate_alert(alert)

    # 决定是否回滚
    if engine.should_rollback(severity):
        if not args.dry_run:
            engine.execute_rollback(severity, alert)
        else:
            logger.info(f"演练模式: 将执行 {severity.value} 级别回滚")


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    main()
```

### 一键回滚 CLI 工具

**完整实现** (Bash):
```bash
#!/bin/bash
# rollback.sh - 一键回滚工具

set -euo pipefail

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/../config/rollback_config.yaml"
LOG_FILE="/var/log/zker/rollback_$(date +%Y%m%d_%H%M%S).log"
ROLLBACK_RECORD_FILE="/var/lib/zker/rollback_history.json"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $@" | tee -a "$LOG_FILE"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $@" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $@" | tee -a "$LOG_FILE"
}

# 加载配置
load_config() {
    if [[ ! -f "$CONFIG_FILE" ]]; then
        log_error "配置文件不存在: $CONFIG_FILE"
        exit 1
    fi

    # 解析 YAML 配置（需要 yq 工具）
    if command -v yq &> /dev/null; then
        KUBERNETES_CONTEXT=$(yq e '.kubernetes.context' "$CONFIG_FILE")
        ROLLBACK_TIMEOUT=$(yq e '.rollback.timeout' "$CONFIG_FILE")
        NOTIFICATION_ENABLED=$(yq e '.notification.enabled' "$CONFIG_FILE")
    else
        log_warn "yq 工具未安装，使用默认配置"
        KUBERNETES_CONTEXT="default"
        ROLLBACK_TIMEOUT=300  # 5 分钟
        NOTIFICATION_ENABLED=true
    fi

    log_info "配置加载成功: K8S_CONTEXT=$KUBERNETES_CONTEXT, TIMEOUT=$ROLLBACK_TIMEOUT"
}

# 记录回滚历史
record_rollback() {
    local severity=$1
    local reason=$2
    local result=$3

    local record=$(cat <<EOF
{
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "severity": "$severity",
  "reason": "$reason",
  "result": "$result",
  "triggered_by": "$(whoami)",
  "hostname": "$(hostname)"
}
EOF
)

    echo "$record" >> "$ROLLBACK_RECORD_FILE"
    log_info "回滚历史已记录: $ROLLBACK_RECORD_FILE"
}

# 发送通知
send_notification() {
    local status=$1
    local message=$2

    if [[ "$NOTIFICATION_ENABLED" != "true" ]]; then
        return
    fi

    log_info "发送通知: $status - $message"

    # Slack 通知
    if command -v slack &> /dev/null; then
        slack send "$message" --channel '#release-ops'
    fi

    # PagerDuty 通知
    if [[ "$status" == "critical" ]]; then
        if command -v pagerduty &> /dev/null; then
            pagerduty trigger --summary "系统回滚" --desc "$message"
        fi
    fi
}

# 应用回滚
rollback_application() {
    local from_version=$1
    local to_version=$2

    log_info "开始应用回滚: $from_version → $to_version"

    # 1. 检查目标版本是否存在
    log_info "检查目标版本: $to_version"
    if ! kubectl get deployment "zker-$to_version" &> /dev/null; then
        log_error "目标版本不存在: zker-$to_version"
        return 1
    fi

    # 2. 更新 Istio VirtualService
    log_info "切换流量到版本: $to_version"
    kubectl patch virtualservice zker-api --type=json -p "[
        {\"op\": \"replace\", \"path\": \"/spec/http/0/route/0/destination/name\", \"value\": \"zker-$to_version\"},
        {\"op\": \"replace\", \"path\": \"/spec/http/0/route/0/weight\", \"value\": 100}
    ]" || {
        log_error "流量切换失败"
        return 1
    }

    # 3. 等待流量切换生效
    log_info "等待流量切换生效..."
    sleep 10

    # 4. 验证新版本健康
    log_info "验证版本 $to_version 健康状态"
    if ! verify_application_health "$to_version"; then
        log_error "健康验证失败"
        return 1
    }

    # 5. 缩容旧版本
    log_info "缩容旧版本: $from_version"
    kubectl scale deployment "zker-$from_version" --replicas=0 || {
        log_warn "旧版本缩容失败（非致命）"
    }

    log_info "应用回滚完成"
    return 0
}

# 数据库回滚
rollback_database() {
    local migration_name=$1

    log_info "开始数据库回滚: $migration_name"

    # 1. 检查回滚脚本是否存在
    local rollback_script="${SCRIPT_DIR}/../migrations/rollback/${migration_name}.sql"
    if [[ ! -f "$rollback_script" ]]; then
        log_error "回滚脚本不存在: $rollback_script"
        return 1
    fi

    # 2. 备份当前数据库
    log_info "备份数据库..."
    local backup_file="/backups/pre_rollback_$(date +%Y%m%d_%H%M%S).sql"
    if ! mysqldump -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" > "$backup_file"; then
        log_error "数据库备份失败"
        return 1
    fi
    log_info "数据库备份完成: $backup_file"

    # 3. 执行回滚脚本
    log_info "执行数据库回滚脚本..."
    if ! mysql -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$rollback_script"; then
        log_error "数据库回滚失败"
        return 1
    fi

    # 4. 验证回滚结果
    log_info "验证数据库回滚结果..."
    if ! verify_database_schema; then
        log_error "数据库验证失败"
        return 1
    fi

    log_info "数据库回滚完成"
    return 0
}

# 配置回滚
rollback_config() {
    local config_type=$1  # env, configmap, feature-flag
    local config_name=$2

    log_info "开始配置回滚: $config_type/$config_name"

    case "$config_type" in
        env)
            # 环境变量回滚
            kubectl set env deployment/"$DEPLOYMENT_NAME" --from=env-backup/"$config_name" || {
                log_error "环境变量回滚失败"
                return 1
            }
            ;;
        configmap)
            # ConfigMap 回滚
            kubectl rollout undo deployment/"$DEPLOYMENT_NAME" --from-configmap="$config_name" || {
                log_error "ConfigMap 回滚失败"
                return 1
            }
            ;;
        feature-flag)
            # 功能开关回滚
            curl -X POST "http://feature-flag-server/api/v1/flags/$config_name/disable" || {
                log_error "功能开关回滚失败"
                return 1
            }
            ;;
        *)
            log_error "未知配置类型: $config_type"
            return 1
            ;;
    esac

    # 重启 Pod 以加载新配置
    log_info "重启 Pod 以加载新配置..."
    kubectl rollout restart deployment/"$DEPLOYMENT_NAME"

    log_info "配置回滚完成"
    return 0
}

# 验证应用健康
verify_application_health() {
    local version=$1
    local timeout=60
    local start_time=$(date +%s)

    log_info "验证应用健康状态 (超时: ${timeout}s)..."

    while true; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))

        if [[ $elapsed -gt $timeout ]]; then
            log_error "健康验证超时"
            return 1
        fi

        # 检查 Pod 状态
        local ready_pods=$(kubectl get pods -l "app=zker,version=$version" -o jsonpath='{.items[*].status.containerStatuses[*].ready}' | wc -w)
        local total_pods=$(kubectl get pods -l "app=zker,version=$version" -o jsonpath='{.items[*].status.containerStatuses[*].ready}' | wc -w)

        if [[ $ready_pods -eq $total_pods ]] && [[ $total_pods -gt 0 ]]; then
            log_info "Pod 状态正常: $ready_pods/$total_pods ready"
        else
            log_warn "Pod 状态异常: $ready_pods/$total_pods ready，等待中..."
            sleep 5
            continue
        fi

        # 检查健康检查端点
        local health_status=$(curl -s -o /dev/null -w "%{http_code}" http://api.zker.com/health || echo "000")
        if [[ "$health_status" == "200" ]]; then
            log_info "健康检查通过"
        else
            log_warn "健康检查失败: HTTP $health_status，等待中..."
            sleep 5
            continue
        fi

        # 所有检查通过
        return 0
    done
}

# 验证数据库模式
verify_database_schema() {
    log_info "验证数据库模式..."

    # 检查关键表是否存在
    local required_tables=("users" "bots" "conversations" "messages")
    for table in "${required_tables[@]}"; do
        if ! mysql -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES LIKE '$table'" | grep -q "$table"; then
            log_error "表不存在: $table"
            return 1
        fi
    done

    log_info "数据库模式验证通过"
    return 0
}

# 主回滚流程
main() {
    local severity=""
    local reason=""
    local rollback_type=""
    local dry_run=false

    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            --severity)
                severity="$2"
                shift 2
                ;;
            --reason)
                reason="$2"
                shift 2
                ;;
            --type)
                rollback_type="$2"
                shift 2
                ;;
            --dry-run)
                dry_run=true
                shift
                ;;
            --help)
                echo "Usage: $0 --severity <critical|high|medium> --reason <reason> --type <application|database|config>"
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                exit 1
                ;;
        esac
    done

    # 参数验证
    if [[ -z "$severity" ]] || [[ -z "$reason" ]] || [[ -z "$rollback_type" ]]; then
        log_error "缺少必要参数"
        exit 1
    fi

    log_info "=========================================="
    log_info "开始一键回滚流程"
    log_info "严重程度: $severity"
    log_info "回滚原因: $reason"
    log_info "回滚类型: $rollback_type"
    log_info "演练模式: $dry_run"
    log_info "=========================================="

    # 加载配置
    load_config

    # 发送回滚开始通知
    send_notification "warning" "回滚开始: $severity - $reason"

    # 记录回滚开始
    record_rollback "$severity" "$reason" "started"

    # 执行回滚
    local result="failed"
    if [[ "$dry_run" == "true" ]]; then
        log_warn "演练模式: 跳过实际回滚操作"
        result="dry_run"
    else
        case "$rollback_type" in
            application)
                if rollback_application "canary" "stable"; then
                    result="success"
                fi
                ;;
            database)
                if rollback_database "add_tenant_id_column"; then
                    result="success"
                fi
                ;;
            config)
                if rollback_config "feature-flag" "intelligent_routing"; then
                    result="success"
                fi
                ;;
            *)
                log_error "未知回滚类型: $rollback_type"
                ;;
        esac
    fi

    # 记录回滚结果
    record_rollback "$severity" "$reason" "$result"

    # 发送回滚结果通知
    if [[ "$result" == "success" ]]; then
        send_notification "info" "回滚成功: $reason"
        log_info "回滚完成 ✓"
    elif [[ "$result" == "dry_run" ]]; then
        send_notification "info" "回滚演练: $reason"
        log_info "回滚演练完成"
    else
        send_notification "critical" "回滚失败: $reason"
        log_error "回滚失败 ✗"
        exit 1
    fi
}

# 执行主函数
main "$@"
```

**使用示例**:
```bash
# 场景 1: 严重故障，立即回滚应用
./rollback.sh \
  --severity critical \
  --reason "错误率 > 10%，服务不可用" \
  --type application

# 场景 2: 数据库迁移失败，回滚数据库
./rollback.sh \
  --severity high \
  --reason "数据库迁移脚本执行失败" \
  --type database

# 场景 3: 功能开关异常，回滚配置
./rollback.sh \
  --severity medium \
  --reason "智能路由功能导致性能下降" \
  --type config

# 场景 4: 演练回滚流程
./rollback.sh \
  --severity critical \
  --reason "定期回滚演练" \
  --type application \
  --dry-run
```

---

## ✅ 回滚验证

### 验证层次

**1. 基础设施验证**
```bash
# 1.1 Pod 状态验证
kubectl get pods -l app=zker,version=stable
# 预期: 所有 Pod 状态为 Running

# 1.2 服务端点验证
kubectl get endpoints zker-api
# 预期: 至少有 3 个可用端点

# 1.3 资源使用验证
kubectl top pods -l app=zker,version=stable
# 预期: CPU < 70%, 内存 < 80%
```

**2. 应用功能验证**
```bash
# 2.1 健康检查
curl -f http://api.zker.com/health || echo "健康检查失败"

# 2.2 核心功能测试
./scripts/smoke_test.sh --endpoint http://api.zker.com
# 测试项目:
# - 用户登录
# - Bot 列表查询
# - 创建对话
# - 发送消息
# - 知识库检索

# 2.3 API 响应时间测试
ab -n 1000 -c 100 http://api.zker.com/api/v1/bots
# 预期: P95 < 2s
```

**3. 数据一致性验证**
```sql
-- 3.1 数据完整性检查
SELECT
  'users' as table_name,
  COUNT(*) as row_count,
  COUNT(DISTINCT user_id) as unique_users
FROM users
UNION ALL
SELECT
  'bots',
  COUNT(*),
  COUNT(DISTINCT bot_id)
FROM bots;

-- 3.2 外键约束检查
SELECT
  COUNT(*) as orphaned_bots
FROM bots b
WHERE NOT EXISTS (
  SELECT 1 FROM users u WHERE u.user_id = b.created_by
);
-- 预期: 0 条孤立记录

-- 3.3 数据哈希对比
--（与回滚前备份对比）
```

**4. 性能验证**
```yaml
# 4.1 关键指标对比
指标:
  - QPS: 应恢复到基线 ± 10%
  - 错误率: < 0.1%
  - P95 响应时间: < 2s
  - CPU 使用率: < 70%
  - 内存使用率: < 80%

# 4.2 Prometheus 查询
# QPS
rate(http_requests_total{version="stable"}[5m])

# 错误率
rate(http_requests_total{version="stable",status=~"5.."}[5m]) / rate(http_requests_total{version="stable"}[5m])

# P95 响应时间
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket{version="stable"}[5m]))
```

---

## 🔄 回滚后恢复

### 数据补偿

**场景**: 双写期间的增量数据补偿

```python
#!/usr/bin/env python3
"""
数据补偿工具
将回滚期间双写到新版本的数据同步回旧版本
"""

import pymysql
import logging

logger = logging.getLogger(__name__)


class DataCompensation:
    """数据补偿器"""

    def __init__(self, source_db_config, target_db_config):
        self.source_db = pymysql.connect(**source_db_config)
        self.target_db = pymysql.connect(**target_db_config)

    def compensate_incremental_data(self, since_timestamp):
        """
        补偿增量数据

        Args:
          since_timestamp: 起始时间戳
        """
        logger.info(f"开始补偿增量数据（起始时间: {since_timestamp}）...")

        # 1. 补偿用户数据
        self._compensate_table('users', since_timestamp)

        # 2. 补偿 Bot 数据
        self._compensate_table('bots', since_timestamp)

        # 3. 补偿对话数据
        self._compensate_table('conversations', since_timestamp)

        # 4. 补偿消息数据
        self._compensate_table('messages', since_timestamp)

        logger.info("增量数据补偿完成")

    def _compensate_table(self, table_name, since_timestamp):
        """
        补偿单张表的数据

        Args:
          table_name: 表名
          since_timestamp: 起始时间戳
        """
        logger.info(f"补偿表: {table_name}")

        # 从源库查询增量数据
        with self.source_db.cursor() as cursor:
            cursor.execute(f"""
                SELECT * FROM {table_name}
                WHERE updated_at > %s
                ORDER BY updated_at
            """, (since_timestamp,))

            rows = cursor.fetchall()

        # 写入目标库
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
                    logger.error(f"数据补偿失败: {row}, 错误: {e}")

            self.target_db.commit()

        logger.info(f"表 {table_name} 补偿完成: {len(rows)} 条记录")


if __name__ == '__main__':
    # 补偿回滚前 30 分钟的数据
    from datetime import datetime, timedelta

    since = datetime.now() - timedelta(minutes=30)

    compensator = DataCompensation(
        source_db_config={
            'host': 'canary-db',
            'user': 'root',
            'password': 'password',
            'database': 'zker_canary'
        },
        target_db_config={
            'host': 'stable-db',
            'user': 'root',
            'password': 'password',
            'database': 'zker_stable'
        }
    )

    compensator.compensate_incremental_data(since)
```

### 修复后重新发布

**流程**:
```
1. 分析回滚原因
   - 查看日志
   - 分析监控数据
   - 收集用户反馈

2. 修复问题
   - 修复代码 Bug
   - 优化配置参数
   - 更新数据库脚本

3. 测试验证
   - 单元测试
   - 集成测试
   - 性能测试

4. 重新发布
   - 按照灰度发布流程重新上线
   - 缩短灰度周期（如果问题已明确修复）
   - 加强监控和告警
```

---

## 📚 附录

### A. 回滚演练计划

**演练目的**:
- 验证回滚流程的有效性
- 测试自动化回滚系统
- 培训运维团队
- 发现潜在问题

**演练频率**:
- 自动回滚系统: 每月 1 次
- 手动回滚流程: 每季度 1 次
- 全流程演练: 每半年 1 次

**演练检查清单**:
```yaml
演练前准备:
  - [ ] 选择演练时间窗口（低峰期）
  - [ ] 准备演练场景脚本
  - [ ] 通知相关人员
  - [ ] 准备回滚演练环境

演练执行:
  - [ ] 模拟告警触发
  - [ ] 验证自动回滚响应
  - [ ] 记录回滚耗时
  - [ ] 验证回滚结果

演练后总结:
  - [ ] 记录演练报告
  - [ ] 识别改进点
  - [ ] 更新回滚文档
  - [ ] 培训团队成员
```

**演练场景示例**:
```bash
# 场景 1: 应用回滚演练
./scripts/rollback_drill.sh --scenario application --severity critical

# 场景 2: 数据库回滚演练
./scripts/rollback_drill.sh --scenario database --severity high

# 场景 3: 配置回滚演练
./scripts/rollback_drill.sh --scenario config --severity medium

# 场景 4: 组合回滚演练
./scripts/rollback_drill.sh --scenario combined --severity critical
```

---

**文档变更历史**:

| 版本 | 日期 | 变更内容 | 作者 |
|-----|------|---------|------|
| v1.0 | 2025-01-01 | 初始版本 | 运维团队 |

**审批记录**:

| 角色 | 姓名 | 审批意见 | 日期 |
|-----|------|---------|------|
| 技术架构委员会 | [待填写] | [待审批] | [待审批] |
| 运维负责人 | [待填写] | [待审批] | [待审批] |
| 项目经理 | [待填写] | [待审批] | [待审批] |

---

**© 2025 ZKER Project. All rights reserved.**
