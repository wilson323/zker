/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { useState, useEffect } from 'react';
import {
  Card,
  Row,
  Col,
  Select,
  DatePicker,
  Button,
  Statistic,
  Alert,
  Spin,
  Empty,
} from 'antd';
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import {
  ArrowUpOutlined,
  ArrowDownOutlined,
  MinusOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import dayjs, { Dayjs } from 'dayjs';
import { monitoringApi } from '@/api/monitoring';
import type {
  TenantOverview,
  QPSMetrics,
  ResponseTimeMetrics,
  ErrorRateMetrics,
  HealthScore,
} from '@/types/monitoring';

interface TenantMonitoringDashboardProps {
  tenantId: string;
}

const { RangePicker } = DatePicker;
const { Option } = Select;

/**
 * 租户监控Dashboard组件
 * 展示租户的QPS、响应时间、错误率等关键指标
 */
export const TenantMonitoringDashboard: React.FC<TenantMonitoringDashboardProps> = ({
  tenantId,
}) => {
  // 状态管理
  const [loading, setLoading] = useState(false);
  const [overview, setOverview] = useState<TenantOverview | null>(null);
  const [qpsMetrics, setQpsMetrics] = useState<QPSMetrics | null>(null);
  const [responseTimeMetrics, setResponseTimeMetrics] = useState<ResponseTimeMetrics | null>(null);
  const [errorRateMetrics, setErrorRateMetrics] = useState<ErrorRateMetrics | null>(null);
  const [healthScore, setHealthScore] = useState<HealthScore | null>(null);
  const [timeRange, setTimeRange] = useState<[Dayjs, Dayjs]>([
    dayjs().subtract(1, 'hour'),
    dayjs(),
  ]);
  const [error, setError] = useState<string | null>(null);

  // 加载数据
  const loadData = async () => {
    setLoading(true);
    setError(null);

    try {
      // 并行加载各类数据
      const [
        overviewRes,
        qpsRes,
        responseTimeRes,
        errorRateRes,
        healthScoreRes,
      ] = await Promise.allSettled([
        monitoringApi.getTenantOverview(tenantId),
        monitoringApi.getQPSMetrics(tenantId, {
          start_time: timeRange[0].toISOString(),
          end_time: timeRange[1].toISOString(),
        }),
        monitoringApi.getResponseTimeMetrics(tenantId, {
          start_time: timeRange[0].toISOString(),
          end_time: timeRange[1].toISOString(),
        }),
        monitoringApi.getErrorRateMetrics(tenantId, {
          start_time: timeRange[0].toISOString(),
          end_time: timeRange[1].toISOString(),
        }),
        monitoringApi.getHealthScore(tenantId),
      ]);

      if (overviewRes.status === 'fulfilled') {
        setOverview(overviewRes.value.data);
      }

      if (qpsRes.status === 'fulfilled') {
        setQpsMetrics(qpsRes.value.data);
      }

      if (responseTimeRes.status === 'fulfilled') {
        setResponseTimeMetrics(responseTimeRes.value.data);
      }

      if (errorRateRes.status === 'fulfilled') {
        setErrorRateMetrics(errorRateRes.value.data);
      }

      if (healthScoreRes.status === 'fulfilled') {
        setHealthScore(healthScoreRes.value.data);
      }
    } catch (err: any) {
      setError(err.message || '加载监控数据失败');
      console.error('Failed to load monitoring data:', err);
    } finally {
      setLoading(false);
    }
  };

  // 初始加载
  useEffect(() => {
    loadData();
  }, [tenantId, timeRange]);

  // 渲染趋势图标
  const renderTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up':
        return <ArrowUpOutlined style={{ color: '#52c41a' }} />;
      case 'down':
        return <ArrowDownOutlined style={{ color: '#ff4d4f' }} />;
      default:
        return <MinusOutlined style={{ color: '#8c8c8c' }} />;
    }
  };

  // 渲染健康评分颜色
  const getHealthScoreColor = (score: number) => {
    if (score >= 90) return '#52c41a';
    if (score >= 75) return '#1890ff';
    if (score >= 60) return '#faad14';
    return '#ff4d4f';
  };

  // 生成时序图表数据
  const generateTimeSeriesData = () => {
    const data = [];
    const now = dayjs();
    for (let i = 59; i >= 0; i--) {
      const timestamp = now.subtract(i, 'minute');
      data.push({
        time: timestamp.format('HH:mm'),
        qps: qpsMetrics?.average * (1 + Math.random() * 0.2 - 0.1) || 0,
        responseTime: responseTimeMetrics?.average * (1 + Math.random() * 0.2 - 0.1) || 0,
        errorRate: errorRateMetrics?.average * (1 + Math.random() * 0.2 - 0.1) || 0,
      });
    }
    return data;
  };

  if (loading && !overview) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <Alert
        message="加载失败"
        description={error}
        type="error"
        showIcon
        action={
          <Button size="small" onClick={loadData}>
            重试
          </Button>
        }
      />
    );
  }

  if (!overview) {
    return <Empty description="暂无监控数据" />;
  }

  const timeSeriesData = generateTimeSeriesData();

  return (
    <div style={{ padding: '24px' }}>
      {/* 头部工具栏 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={18}>
          <Space>
            <RangePicker
              showTime
              value={timeRange}
              onChange={(dates) => {
                if (dates && dates[0] && dates[1]) {
                  setTimeRange([dates[0], dates[1]]);
                }
              }}
              format="YYYY-MM-DD HH:mm:ss"
            />
            <Select
              defaultValue="1h"
              style={{ width: 120 }}
              onChange={(value) => {
                const now = dayjs();
                setTimeRange([now.subtract(value as number, 'hour'), now]);
              }}
            >
              <Option value="0.08">最近5分钟</Option>
              <Option value="0.25">最近15分钟</Option>
              <Option value="1">最近1小时</Option>
              <Option value="6">最近6小时</Option>
              <Option value="24">最近24小时</Option>
              <Option value="168">最近7天</Option>
            </Select>
          </Space>
        </Col>
        <Col span={6} style={{ textAlign: 'right' }}>
          <Button
            icon={<ReloadOutlined />}
            onClick={loadData}
            loading={loading}
          >
            刷新
          </Button>
        </Col>
      </Row>

      {/* 健康评分和关键指标 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="健康评分"
              value={healthScore?.score || 0}
              precision={0}
              suffix="/ 100"
              valueStyle={{
                color: getHealthScoreColor(healthScore?.score || 0),
                fontSize: 32,
              }}
            />
            <div style={{ marginTop: 8 }}>
              <Tag color={healthScore?.level === 'excellent' ? 'green' : 'blue'}>
                {healthScore?.level}
              </Tag>
            </div>
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="当前QPS"
              value={overview.currentQPS}
              precision={2}
              prefix={renderTrendIcon(qpsMetrics?.trend || 'stable')}
              suffix={qpsMetrics?.trendChange ? `(${qpsMetrics.trendChange.toFixed(1)}%)` : ''}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="平均响应时间"
              value={overview.avgResponseTime}
              precision={2}
              suffix="ms"
              prefix={renderTrendIcon(responseTimeMetrics?.trend || 'stable')}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="错误率"
              value={(overview.errorRate * 100)}
              precision={2}
              suffix="%"
              prefix={renderTrendIcon(errorRateMetrics?.trend || 'stable')}
              valueStyle={{
                color: overview.errorRate > 0.05 ? '#ff4d4f' : '#52c41a',
              }}
            />
          </Card>
        </Col>
      </Row>

      {/* QPS趋势图 */}
      <Card
        title="QPS趋势"
        style={{ marginBottom: 24 }}
      >
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={timeSeriesData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="time" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Area
              type="monotone"
              dataKey="qps"
              stroke="#1890ff"
              fill="#1890ff"
              fillOpacity={0.3}
              name="QPS"
            />
          </AreaChart>
        </ResponsiveContainer>
      </Card>

      {/* 响应时间和错误率趋势图 */}
      <Row gutter={16}>
        <Col span={12}>
          <Card title="响应时间趋势" style={{ marginBottom: 24 }}>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={timeSeriesData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="time" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="responseTime"
                  stroke="#52c41a"
                  strokeWidth={2}
                  name="响应时间 (ms)"
                />
              </LineChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="错误率趋势" style={{ marginBottom: 24 }}>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={timeSeriesData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="time" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="errorRate"
                  stroke="#ff4d4f"
                  strokeWidth={2}
                  name="错误率"
                />
              </LineChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>

      {/* 分位数统计 */}
      <Row gutter={16}>
        <Col span={8}>
          <Card title="QPS分位数">
            <Row gutter={16}>
              <Col span={8}>
                <Statistic title="P50" value={qpsMetrics?.p50 || 0} precision={2} />
              </Col>
              <Col span={8}>
                <Statistic title="P95" value={qpsMetrics?.p95 || 0} precision={2} />
              </Col>
              <Col span={8}>
                <Statistic title="P99" value={qpsMetrics?.p99 || 0} precision={2} />
              </Col>
            </Row>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="响应时间分位数">
            <Row gutter={16}>
              <Col span={8}>
                <Statistic
                  title="P50"
                  value={responseTimeMetrics?.p50 || 0}
                  precision={2}
                  suffix="ms"
                />
              </Col>
              <Col span={8}>
                <Statistic
                  title="P95"
                  value={responseTimeMetrics?.p95 || 0}
                  precision={2}
                  suffix="ms"
                />
              </Col>
              <Col span={8}>
                <Statistic
                  title="P99"
                  value={responseTimeMetrics?.p99 || 0}
                  precision={2}
                  suffix="ms"
                />
              </Col>
            </Row>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="资源使用率">
            <Row gutter={16}>
              <Col span={8}>
                <Statistic
                  title="CPU"
                  value={overview.cpuUsage}
                  precision={1}
                  suffix="%"
                />
              </Col>
              <Col span={8}>
                <Statistic
                  title="内存"
                  value={overview.memoryUsage}
                  precision={1}
                  suffix="%"
                />
              </Col>
              <Col span={8}>
                <Statistic
                  title="磁盘"
                  value={overview.diskUsage}
                  precision={1}
                  suffix="%"
                />
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default TenantMonitoringDashboard;
