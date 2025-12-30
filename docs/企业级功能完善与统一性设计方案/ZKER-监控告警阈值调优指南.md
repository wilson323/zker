# ZKER 监控告警阈值调优指南

**版本**: v1.0.0
**最后更新**: 2025-01-01
**状态**: 持续优化

---

## 📋 目录

- [调优概述](#调优概述)
- [数据收集方法](#数据收集方法)
- [阈值计算方法](#阈值计算方法)
- [分级告警策略](#分级告警策略)
- [动态调整机制](#动态调整机制)
- [实战案例分析](#实战案例分析)
- [调优最佳实践](#调优最佳实践)

---

## 🎯 调优概述

### 调优目标

**监控告警阈值调优** 的核心目标是：

1. **减少误报**：避免"狼来了"效应，让团队对告警保持敏感度
2. **减少漏报**：确保真正的问题能被及时发现
3. **提高响应速度**：通过合理的分级，让关键问题得到优先处理
4. **适应业务变化**：随着业务增长，阈值需要动态调整
5. **平衡运维成本**：避免告警过多导致运维人员疲劳

### 调优原则

**SMART原则**：
- **Specific（具体）**：每个告警都有明确的含义和应对措施
- **Measurable（可衡量）**：阈值基于数据，而非主观判断
- **Achievable（可实现）**：阈值设定要切合实际，避免过高或过低
- **Relevant（相关）**：告警要与业务目标相关
- **Time-bound（有时限）**：定期review和调整阈值

**告警金字塔**：
```
        /\
       /  \  紧急告警（电话/PagerDuty）
      /____\
     /      \  警告告警（即时通讯）
    /________\
   /          \  信息告警（邮件/日志）
  /____________\
```

### 调优时机

**需要调整阈值的信号**：
- ❌ 误报率 > 50%（超过一半的告警是误报）
- ❌ 漏报率 > 10%（有故障但未触发告警）
- ❌ 告警疲劳（团队成员开始忽略告警）
- ❌ 业务增长 > 20%（业务规模变化较大）
- ❌ 系统架构变更（新增服务、扩容等）

**定期review周期**：
- **每周**：review误报和漏报情况
- **每月**：分析告警趋势，调整阈值
- **每季度**：全面review告警策略

---

## 📊 数据收集方法

### 1. 历史数据分析

**Prometheus查询历史数据**：
```promql
# 查询过去7天的CPU使用率
100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[7d])) by (instance)

# 查询过去30天的API P95响应时间
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[30d]))

# 查询过去7天的错误率
(rate(http_requests_total{status=~"5.."}[7d]) / rate(http_requests_total[7d])) * 100
```

**导出数据到CSV**：
```bash
# 使用Prometheus API导出数据
curl -G 'http://prometheus.coze-studio.com/api/v1/query_range' \
  -d 'query=100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance))' \
  -d 'start=2025-01-01T00:00:00Z' \
  -d 'end=2025-01-08T00:00:00Z' \
  -d 'step=3600' > cpu_usage_7days.json

# 转换为CSV（使用jq）
jq -r '.data.result[] | [.metric.instance, (.values | map(join(",")) | join("\n"))] | join(",")' cpu_usage_7days.json > cpu_usage_7days.csv
```

### 2. 百分位数计算

**为什么使用百分位数**：
- **平均值**容易受极端值影响
- **百分位数**更能反映真实体验

**常用百分位数**：
```
P50 (中位数): 50%的用户体验到的性能
P95: 95%的用户体验到的性能（SLA常用）
P99: 99%的用户体验到的性能（关键业务）
P99.9: 99.9%的用户体验到的性能（极度关键）
```

**计算方法**：
```python
import pandas as pd
import numpy as np

# 读取CSV数据
df = pd.read_csv('api_response_time_7days.csv')

# 计算百分位数
p50 = df['response_time'].quantile(0.50)
p95 = df['response_time'].quantile(0.95)
p99 = df['response_time'].quantile(0.99)
p999 = df['response_time'].quantile(0.999)

print(f"P50: {p50}ms")
print(f"P95: {p95}ms")
print(f"P99: {p99}ms")
print(f"P99.9: {p999}ms")
```

### 3. 季节性分析

**识别业务周期**：
```
时间维度        峰值时段      最低时段      调整策略
-------------------------------------------------
小时           9-11点, 14-16点  凌晨2-6点   分时段阈值
星期           工作日        周末         分星期阈值
月份           月末、季末    月中         月度动态阈值
特殊事件       双11、618    平时         临时阈值
```

**Python代码示例**：
```python
# 分析不同时段的流量模式
df['hour'] = pd.to_datetime(df['timestamp']).dt.hour
hourly_avg = df.groupby('hour')['request_count'].mean()

# 找出高峰时段
peak_hours = hourly_avg[hourly_avg > hourly_avg.quantile(0.75)].index.tolist()
print(f"高峰时段: {peak_hours}")

# 为高峰时段设置更高的阈值
peak_threshold = hourly_avg.max() * 1.2
normal_threshold = hourly_avg.mean() * 1.2
```

### 4. 异常值检测

**使用3-sigma规则**：
```python
# 计算3-sigma边界
mean = df['cpu_usage'].mean()
std = df['cpu_usage'].std()

upper_bound = mean + 3 * std
lower_bound = mean - 3 * std

# 识别异常值
outliers = df[(df['cpu_usage'] > upper_bound) | (df['cpu_usage'] < lower_bound)]

print(f"异常值数量: {len(outliers)}")
print(f"异常值占比: {len(outliers) / len(df) * 100:.2f}%")
```

**使用IQR（四分位距）**：
```python
# 计算IQR
Q1 = df['response_time'].quantile(0.25)
Q3 = df['response_time'].quantile(0.75)
IQR = Q3 - Q1

# 定义异常值边界
lower_bound = Q1 - 1.5 * IQR
upper_bound = Q3 + 1.5 * IQR

# 识别异常值
outliers = df[(df['response_time'] < lower_bound) | (df['response_time'] > upper_bound)]

# 评估异常值是否需要关注
if len(outliers) / len(df) > 0.05:  # 超过5%是异常值
    print("⚠️ 系统存在较多异常值，需要调查原因")
else:
    print("✅ 异常值在合理范围内")
```

---

## 🧮 阈值计算方法

### 方法1：基于统计分布

**正态分布法**：
```python
# 假设数据符合正态分布
import scipy.stats as stats

mean = df['cpu_usage'].mean()
std = df['cpu_usage'].std()

# 计算不同置信度下的阈值
confidence_levels = {
    'warning': stats.norm.ppf(0.80),   # 80%分位数
  'critical': stats.norm.ppf(0.95),   # 95%分位数
  'emergency': stats.norm.ppf(0.99),  # 99%分位数
}

thresholds = {}
for level, z_score in confidence_levels.items():
    thresholds[level] = mean + z_score * std

print(thresholds)
# 输出: {'warning': 75%, 'critical': 85%, 'emergency': 95%}
```

**百分位数法（推荐）**：
```python
# 直接使用历史数据的百分位数作为阈值
thresholds = {
    'warning': df['cpu_usage'].quantile(0.90),    # P90
    'critical': df['cpu_usage'].quantile(0.95),   # P95
    'emergency': df['cpu_usage'].quantile(0.99),  # P99
}

print(thresholds)
# 输出: {'warning': 78%, 'critical': 85%, 'emergency': 94%}
```

### 方法2：基于业务目标

**SLA驱动法**：
```python
# 假设SLA要求：99%的请求在500ms内响应
sla_p99 = 500  # ms
sla_target = 0.99

# 计算当前性能是否满足SLA
current_p99 = df['response_time'].quantile(0.99)

if current_p99 > sla_p99:
    print(f"❌ 当前性能不满足SLA: P99={current_p99}ms > SLA={sla_p99}ms")
    # 计算需要优化的程度
    optimization_needed = (current_p99 - sla_p99) / sla_p99 * 100
    print(f"需要优化: {optimization_needed:.1f}%")
else:
    print(f"✅ 当前性能满足SLA: P99={current_p99}ms")

# 设置告警阈值（在SLA基础上留有余量）
alert_thresholds = {
    'warning': sla_p99 * 0.7,      # 350ms（提前预警）
    'critical': sla_p99 * 0.9,     # 450ms（即将突破）
    'emergency': sla_p99 * 1.1,    # 550ms（已经突破）
}
```

### 方法3：基于机器学习

**使用LSTM预测异常**：
```python
from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import LSTM, Dense
import numpy as np

# 准备时间序列数据
def prepare_data(data, look_back=60):
    X, y = [], []
    for i in range(len(data) - look_back):
        X.append(data[i:i+look_back])
        y.append(data[i+look_back])
    return np.array(X), np.array(y)

# 训练LSTM模型
model = Sequential([
    LSTM(50, input_shape=(look_back, 1)),
    Dense(1)
])

model.compile(optimizer='adam', loss='mse')

# 假设df['cpu_usage']是时间序列
X, y = prepare_data(df['cpu_usage'].values)
model.fit(X, y, epochs=100, batch_size=32)

# 预测下一个值
last_60_points = df['cpu_usage'].values[-60:].reshape(1, 60, 1)
prediction = model.predict(last_60_points)

# 计算预测误差
actual = df['cpu_usage'].values[-1]
error = abs(prediction[0][0] - actual)

# 设置告警阈值（基于预测误差）
if error > 3 * std:  # 超过3倍标准差
    print("🚨 检测到异常！")
```

### 方法4：动态阈值算法

**使用EWMA（指数加权移动平均）**：
```python
# 计算EWMA
def calculate_ewma(data, alpha=0.3):
    ewma = [data[0]]
    for i in range(1, len(data)):
        ewma.append(alpha * data[i] + (1 - alpha) * ewma[-1])
    return ewma

df['ewma'] = calculate_ewma(df['cpu_usage'].values)

# 计算动态阈值（基于EWMA）
df['upper_threshold'] = df['ewma'] + 2 * df['cpu_usage'].std()
df['lower_threshold'] = df['ewma'] - 2 * df['cpu_usage'].std()

# 检测异常
anomalies = df[(df['cpu_usage'] > df['upper_threshold']) |
              (df['cpu_usage'] < df['lower_threshold'])]

print(f"检测到 {len(anomalies)} 个异常点")
```

**Prometheus动态阈值查询**：
```promql
# 基于过去7天的平均值 + 2倍标准差
threshold = (
  avg_over_time(cpu_usage[7d]) +
  2 * stddev_over_time(cpu_usage[7d])
)

# 判断当前值是否超过阈值
cpu_usage > threshold
```

---

## 🚨 分级告警策略

### 告警级别定义

| 级别 | 名称 | 触发条件 | 通知方式 | 响应时间 | 处理优先级 |
|------|------|---------|---------|---------|-----------|
| **P0** | 紧急 | 系统不可用或即将崩溃 | 电话 + PagerDuty + 所有渠道 | 5分钟 | 最高 |
| **P1** | 严重 | 核心功能受损，影响大量用户 | 即时通讯（企业微信/钉钉）+ 邮件 | 15分钟 | 高 |
| **P2** | 警告 | 性能下降，但仍可用 | 邮件 + 日志记录 | 1小时 | 中 |
| **P3** | 信息 | 指标异常，但影响不大 | 日志记录 | 4小时 | 低 |

### CPU使用率告警策略

**示例**：
```yaml
# alerts.yml
groups:
  - name: cpu_alerts
    interval: 30s
    rules:
      # P3: 信息告警 - CPU使用率偏高
      - alert: CPUUsageHigh
        expr: 100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) > 70
        for: 10m
        labels:
          severity: info
          priority: P3
        annotations:
          summary: "CPU使用率偏高"
          description: "实例 {{ $labels.instance }} CPU使用率持续 > 70% 已10分钟"
          runbook: "https://docs.coze-studio.com/runbooks/cpu-high"

      # P2: 警告告警 - CPU使用率高
      - alert: CPUUsageVeryHigh
        expr: 100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) > 80
        for: 5m
        labels:
          severity: warning
          priority: P2
        annotations:
          summary: "CPU使用率高"
          description: "实例 {{ $labels.instance }} CPU使用率持续 > 80% 已5分钟"
          runbook: "https://docs.coze-studio.com/runbooks/cpu-very-high"

      # P1: 严重告警 - CPU使用率非常高
      - alert: CPUUsageCritical
        expr: 100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) > 90
        for: 2m
        labels:
          severity: critical
          priority: P1
        annotations:
          summary: "CPU使用率严重告警"
          description: "实例 {{ $labels.instance }} CPU使用率持续 > 90% 已2分钟"
          runbook: "https://docs.coze-studio.com/runbooks/cpu-critical"

      # P0: 紧急告警 - CPU即将饱和
      - alert: CPUUsageEmergency
        expr: 100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) > 95
        for: 1m
        labels:
          severity: emergency
          priority: P0
        annotations:
          summary: "CPU使用率紧急告警"
          description: "实例 {{ $labels.instance }} CPU使用率持续 > 95% 已1分钟，即将饱和！"
          runbook: "https://docs.coze-studio.com/runbooks/cpu-emergency"
```

**阈值设置逻辑**：
```
警告（70%）：基于P90，预留30%余量
严重（80%）：基于P95，预留20%余量
紧急（90%）：基于P99，预留10%余量
危急（95%）：基于P99.9，立即处理
```

### API响应时间告警策略

**示例**：
```yaml
groups:
  - name: api_performance_alerts
    interval: 30s
    rules:
      # P3: 信息告警 - P95响应时间偏高
      - alert: APIResponseTimeHigh
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 0.5
        for: 10m
        labels:
          severity: info
          priority: P3
        annotations:
          summary: "API P95响应时间偏高"
          description: "API {{ $labels.endpoint }} P95响应时间 > 500ms 已10分钟"

      # P2: 警告告警 - P95响应时间高
      - alert: APIResponseTimeVeryHigh
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1.0
        for: 5m
        labels:
          severity: warning
          priority: P2
        annotations:
          summary: "API P95响应时间高"
          description: "API {{ $labels.endpoint }} P95响应时间 > 1秒 已5分钟"

      # P1: 严重告警 - P95响应时间非常高
      - alert: APIResponseTimeCritical
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 2.0
        for: 2m
        labels:
          severity: critical
          priority: P1
        annotations:
          summary: "API P95响应时间严重告警"
          description: "API {{ $labels.endpoint }} P95响应时间 > 2秒 已2分钟"

      # P0: 紧急告警 - P95响应时间极慢
      - alert: APIResponseTimeEmergency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 5.0
        for: 1m
        labels:
          severity: emergency
          priority: P0
        annotations:
          summary: "API P95响应时间紧急告警"
          description: "API {{ $labels.endpoint }} P95响应时间 > 5秒 已1分钟，用户体验极差！"
```

### 业务指标告警策略

**Bot调用失败率**：
```yaml
groups:
  - name: business_alerts
    interval: 30s
    rules:
      # P2: 警告告警 - Bot失败率偏高
      - alert: BotFailureRateHigh
        expr: (sum(rate(bot_invocation_total{status="failed"}[5m])) by (tenant_id) / sum(rate(bot_invocation_total[5m])) by (tenant_id)) > 0.05
        for: 10m
        labels:
          severity: warning
          priority: P2
        annotations:
          summary: "Bot调用失败率偏高"
          description: "租户 {{ $labels.tenant_id }} Bot调用失败率 > 5% 已10分钟"

      # P1: 严重告警 - Bot失败率高
      - alert: BotFailureRateCritical
        expr: (sum(rate(bot_invocation_total{status="failed"}[5m])) by (tenant_id) / sum(rate(bot_invocation_total[5m])) by (tenant_id)) > 0.10
        for: 5m
        labels:
          severity: critical
          priority: P1
        annotations:
          summary: "Bot调用失败率严重"
          description: "租户 {{ $labels.tenant_id }} Bot调用失败率 > 10% 已5分钟"

      # P0: 紧急告警 - Bot失败率极高
      - alert: BotFailureRateEmergency
        expr: (sum(rate(bot_invocation_total{status="failed"}[5m])) by (tenant_id) / sum(rate(bot_invocation_total[5m])) by (tenant_id)) > 0.20
        for: 2m
        labels:
          severity: emergency
          priority: P0
        annotations:
          summary: "Bot调用失败率紧急"
          description: "租户 {{ $labels.tenant_id }} Bot调用失败率 > 20% 已2分钟，服务严重受损！"
```

---

## 🔄 动态调整机制

### 1. 自动调优脚本

**基于历史数据自动调整阈值**：
```python
#!/usr/bin/env python3
# scripts/auto_adjust_alert_thresholds.py

import prometheus_api_client as prom
import yaml
import json

# 从Prometheus获取历史数据
def get_metric_data(metric_name, days=7):
    prom_client = prom.PrometheusConnect(url="http://prometheus.coze-studio.com")

    # 查询过去7天的数据
    query = f'{metric_name}[{days}d]'
    result = prom_client.custom_query(query)

    # 提取时间序列数据
    data = []
    for series in result:
        for value in series['values']:
            data.append(float(value[1]))

    return data

# 计算阈值
def calculate_thresholds(data):
    import numpy as np

    thresholds = {
        'info': np.percentile(data, 90),
        'warning': np.percentile(data, 95),
        'critical': np.percentile(data, 99),
        'emergency': np.percentile(data, 99.9),
    }

    return thresholds

# 更新告警规则文件
def update_alert_rules(rule_file, metric_name, new_thresholds):
    with open(rule_file, 'r') as f:
        alerts = yaml.safe_load(f)

    # 更新阈值
    for group in alerts['groups']:
        for rule in group['rules']:
            if metric_name in rule['expr']:
                # 根据告警级别更新阈值
                if rule['labels']['severity'] == 'info':
                    rule['expr'] = rule['expr'].replace('> 70', f'> {new_thresholds["info"]:.0f}')
                elif rule['labels']['severity'] == 'warning':
                    rule['expr'] = rule['expr'].replace('> 80', f'> {new_thresholds["warning"]:.0f}')
                elif rule['labels']['severity'] == 'critical':
                    rule['expr'] = rule['expr'].replace('> 90', f'> {new_thresholds["critical"]:.0f}')
                elif rule['labels']['severity'] == 'emergency':
                    rule['expr'] = rule['expr'].replace('> 95', f'> {new_thresholds["emergency"]:.0f}')

    # 写回文件
    with open(rule_file, 'w') as f:
        yaml.dump(alerts, f)

# 主函数
def main():
    metrics = [
        '100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance)',
        'histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))',
    ]

    for metric in metrics:
        # 获取历史数据
        data = get_metric_data(metric, days=7)

        # 计算阈值
        thresholds = calculate_thresholds(data)

        # 更新告警规则
        update_alert_rules('deploy/monitoring/alerts.yml', metric, thresholds)

        print(f"✅ 已更新 {metric} 的阈值: {thresholds}")

if __name__ == '__main__':
    main()
```

### 2. A/B测试验证

**对比新旧阈值的效果**：
```python
# 对新旧阈值进行A/B测试
old_threshold = 80
new_threshold = 85

# 统计旧阈值的误报率
false_positives_old = len(df[(df['cpu_usage'] > old_threshold) & (df['is_anomaly'] == False)])
total_alerts_old = len(df[df['cpu_usage'] > old_threshold])
false_positive_rate_old = false_positives_old / total_alerts_old

# 统计新阈值的误报率
false_positives_new = len(df[(df['cpu_usage'] > new_threshold) & (df['is_anomaly'] == False)])
total_alerts_new = len(df[df['cpu_usage'] > new_threshold])
false_positive_rate_new = false_positives_new / total_alerts_new

print(f"旧阈值误报率: {false_positive_rate_old:.2%}")
print(f"新阈值误报率: {false_positive_rate_new:.2%}")

# 如果新阈值误报率更低，则采用新阈值
if false_positive_rate_new < false_positive_rate_old:
    print("✅ 新阈值更优，建议采用")
else:
    print("❌ 新阈值不如旧阈值，保持原阈值")
```

### 3. 季节性动态调整

**根据时段自动调整**：
```promql
# 白天工作时间（9-18点）使用较高阈值
# 夜间低峰期（18-9点）使用较低阈值

# 获取当前小时
current_hour = hour()

# 根据时段设置不同阈值
threshold = (
  (current_hour >= 9 and current_hour <= 18) * 80 +  # 白天: 80%
  (current_hour < 9 or current_hour > 18) * 70       # 夜间: 70%
)

# 判断CPU使用率是否超过阈值
cpu_usage > threshold
```

**完整Prometheus规则**：
```yaml
groups:
  - name: dynamic_cpu_alerts
    interval: 30s
    rules:
      - alert: CPUUsageAdaptive
        # 动态阈值：白天80%，夜间70%
        expr: |
          100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) > (
            (hour() >= 9 and hour() <= 18) * 80 +
            (hour() < 9 or hour() > 18) * 70
          )
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "CPU使用率超过动态阈值"
          description: "实例 {{ $labels.instance }} CPU使用率超过 {{ $value }}%（当前时段阈值: {{ $value }}%）"
```

---

## 📈 实战案例分析

### 案例1：误报率过高的CPU告警

**问题**：
- 原始阈值：CPU使用率 > 70%
- 误报率：60%（100个告警中有60个是误报）
- 团队反应：开始忽略告警

**分析过程**：
1. 收集过去30天的CPU使用率数据
2. 绘制分布图，发现CPU使用率呈双峰分布：
   - 第一个峰值：60-65%（正常业务处理）
   - 第二个峰值：85-90%（定时任务批处理）
3. 发现定时任务在凌晨2点运行，持续1小时

**解决方案**：
```yaml
# 方案1：分时段设置阈值
- alert: CPUUsageHigh
  expr: |
    100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) > (
      # 白天工作时间: 75%
      (hour() >= 9 and hour() <= 18) * 75 +
      # 夜间低峰期: 65%
      (hour() < 9 or hour() > 18) * 65 +
      # 凌晨批处理时段: 90%
      (hour() >= 2 and hour() <= 3) * 90
    )
  for: 10m

# 方案2：排除批处理时段
- alert: CPUUsageHigh
  expr: 100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) > 75
  for: 10m
  # 在2-3点之间不告警
  inhibit_rules:
    - source_match:
        hour: '2|3'
      target_match:
        alertname: CPUUsageHigh
      equal: ['instance']
```

**效果**：
- 误报率从60%降低到15%
- 团队重新对告警保持敏感

### 案例2：漏报的数据库连接池告警

**问题**：
- 原始阈值：连接数 > 150（最大200）
- 漏报情况：有3次连接池耗尽，但未触发告警
- 原因分析：连接数在1分钟内从100飙升到200，但`for: 5m`条件未满足

**解决方案**：
```yaml
# 方案1：使用速率检测
- alert: DBConnectionPoolRapidIncrease
  expr: |
    # 检测连接数增长率
    (db_connections_active - db_connections_active offset 1m) / db_connections_active offset 1m > 0.5
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "数据库连接数快速增长"
    description: "1分钟内连接数增长超过50%"

# 方案2：预测性告警
- alert: DBConnectionPoolWillExhaust
  expr: |
    # 预测2分钟后会耗尽
    predict_linear(db_connections_active[5m], 2*60) > 180
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "数据库连接池即将耗尽"
    description: "预测2分钟后连接数将达到 {{ $value }}，接近上限200"
```

**效果**：
- 漏报率从30%降低到5%
- 提前2-3分钟发现连接池增长趋势

### 案例3：业务量增长导致的阈值失效

**问题**：
- 原始阈值：API错误率 > 0.1%（1/1000）
- 业务量增长：从1万RPS增长到10万RPS
- 结果：错误率绝对值增加了10倍，但仍低于阈值

**解决方案**：
```yaml
# 方案1：绝对值+相对值双重条件
- alert: APIErrorRateHigh
  expr: |
    # 绝对值: 错误率 > 0.1%
    (rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])) > 0.001
    and
    # 相对值: 错误数环比增长 > 50%
    (
      rate(http_requests_total{status=~"5.."}[5m]) -
      rate(http_requests_total{status=~"5.."}[5m] offset 5m)
    ) / rate(http_requests_total{status=~"5.."}[5m] offset 5m) > 0.5
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "API错误率过高"
    description: "错误率 > 0.1% 且环比增长 > 50%"

# 方案2：动态阈值（基于历史平均值）
- alert: APIErrorRateAnomaly
  expr: |
    # 当前错误率
    (rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]))
    >
    # 过去7天平均值 + 3倍标准差
    (
      avg_over_time(rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])[7d]) +
      3 * stddev_over_time(rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])[7d])
    )
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "API错误率异常"
    description: "错误率超出历史正常范围"
```

**效果**：
- 成功捕获了业务量增长后的异常错误率
- 减少了因业务增长导致的误报

---

## 💡 调优最佳实践

### 1. 告警精炼原则

**DRY原则（Don't Repeat Yourself）**：
```yaml
# ❌ Bad: 重复的告警规则
- alert: CPUHighInstance1
  expr: cpu_usage{instance="instance1"} > 80

- alert: CPUHighInstance2
  expr: cpu_usage{instance="instance2"} > 80

- alert: CPUHighInstance3
  expr: cpu_usage{instance="instance3"} > 80

# ✅ Good: 使用by标签聚合
- alert: CPUHigh
  expr: avg(cpu_usage) by (instance) > 80
```

**告警抑制（减少冗余告警）**：
```yaml
# 高优先级告警会抑制低优先级告警
inhibit_rules:
  # 如果CPU紧急告警触发，抑制CPU警告告警
  - source_match:
      severity: 'emergency'
      alertname: 'CPUUsageEmergency'
    target_match:
      severity: 'warning'
      alertname: 'CPUUsageHigh'
    equal: ['instance']

  # 如果服务完全不可用，抑制所有该服务的其他告警
  - source_match:
      alertname: 'ServiceDown'
    target_match_re:
      service: '.*'
    equal: ['instance']
```

### 2. 告警通知优化

**减少告警风暴**：
```yaml
# 使用group_by和group_wait聚合告警
receivers:
  - name: 'wechat'
    wechat_configs:
      - corp_id: 'xxx'
        api_secret: 'xxx'
        to_party: 'ops-team'
        # 聚合配置
        group_by: ['alertname', 'cluster']
        group_wait: 30s      # 等待30秒聚合同组告警
        group_interval: 5m   # 同组告警5分钟内只发送一次
        repeat_interval: 4h  # 重复告警每4小时发送一次
```

**告警分级通知**：
```yaml
# 根据严重级别选择不同通知渠道
receivers:
  # P0紧急告警：电话 + 即时通讯
  - name: 'emergency'
    pagerduty_configs:
      - service_key: 'xxx'
    wechat_configs:
      - corp_id: 'xxx'
        to_user: 'ops-lead'

  # P1严重告警：即时通讯
  - name: 'critical'
    wechat_configs:
      - corp_id: 'xxx'
        to_party: 'ops-team'

  # P2警告告警：邮件
  - name: 'warning'
    email_configs:
      - to: 'ops-team@coze-studio.com'

  # P3信息告警：仅日志
  - name: 'info'
    webhook_configs:
      - url: 'http://log-server.com/api/logs'
```

### 3. 持续改进流程

**每周Review Checklist**：
- [ ] 统计误报数量和占比
- [ ] 分析漏报原因（如果有的话）
- [ ] 识别告警疲劳（是否有团队成员忽略告警）
- [ ] 回顾告警响应时间是否符合预期
- [ ] 根据上述分析调整阈值

**每月优化任务**：
- [ ] 分析所有告警的分布情况
- [ ] 识别最频繁的告警（可能是误报）
- [ ] 评估新增告警的必要性
- [ ] 更新Runbook文档
- [ ] 团队分享告警处理经验

**季度全面review**：
- [ ] 重新审视所有告警规则的必要性
- [ ] 基于业务增长调整阈值
- [ ] 优化告警通知策略
- [ ] 评估告警系统的整体效果
- [ ] 制定下个季度的优化计划

### 4. 避免常见陷阱

**陷阱1：阈值设定过低**：
```yaml
# ❌ Bad: 过于敏感的阈值
- alert: CPUUsageHigh
  expr: cpu_usage > 50  # 50%就告警，太低了
  for: 1m

# ✅ Good: 合理的阈值
- alert: CPUUsageHigh
  expr: cpu_usage > 80  # 80%才告警，基于P95
  for: 5m
```

**陷阱2：持续时间过短**：
```yaml
# ❌ Bad: 持续时间太短，容易误报
- alert: APIErrorRate
  expr: error_rate > 0.01
  for: 10s  # 10秒太短了

# ✅ Good: 合理的持续时间
- alert: APIErrorRate
  expr: error_rate > 0.01
  for: 5m   # 5分钟更合理
```

**陷阱3：忽略业务上下文**：
```yaml
# ❌ Bad: 忽略定时任务
- alert: CPUUsageHigh
  expr: cpu_usage > 80
  for: 5m
  # 凌晨2点的批处理任务也会触发

# ✅ Good: 考虑业务上下文
- alert: CPUUsageHigh
  expr: cpu_usage > 80
  for: 5m
  # 排除批处理时段
  inhibit_rules:
    - source_match:
        hour: '2|3'
```

---

## 📚 相关文档

- [Prometheus告警规则说明](../../deploy/monitoring/README.md)
- [Grafana监控大盘使用指南](../../deploy/monitoring/GRAFANA_GUIDE.md)
- [性能基线文档](./ZKER-性能基线文档.md)
- [故障排查手册](./ZKER-故障排查手册_v1.0.md)

---

**更新日志**：

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，完整的监控告警阈值调优指南 | Claude AI |
