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

/**
 * Token使用统计页面
 *
 * 展示Token使用概览、模型使用排行、每日使用趋势、Bot使用统计、实时Token记录
 */

import React, { useEffect, useState, useMemo } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Layout, Radio, Spin, Typography } from '@coze-studio/ui-components';
import * as echarts from 'echarts/core';
import { LineChart, BarChart, PieChart } from 'echarts/charts';
import {
  GridComponent,
  TooltipComponent,
  TitleComponent,
  LegendComponent,
} from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import ReactECharts from 'echarts-for-react';

import { useTranslation } from '@coze-studio/i18n';
import { getBillingService } from '@coze-studio/api-client';
import type {
  UsageStatsResponse,
  DailyUsageStatsResponse,
  ModelUsageStatsResponse,
} from '@coze-studio/common/types/billing';

import { useOverviewCards } from './hooks/useOverviewCards';
import { useModelRanking } from './hooks/useModelRanking';
import { useDailyTrend } from './hooks/useDailyTrend';
import { useBotUsage } from './hooks/useBotUsage';
import { useRealtimeRecords } from './hooks/useRealtimeRecords';

// 注册ECharts组件
echarts.use([
  LineChart,
  BarChart,
  PieChart,
  GridComponent,
  TooltipComponent,
  TitleComponent,
  LegendComponent,
  CanvasRenderer,
]);

const { Header, Content } = Layout;
const { Title } = Typography;

/**
 * Token使用统计页面组件
 */
export const TokenUsagePage: React.FC = () => {
  const { tenantId = 'current' } = useParams<{ tenantId: string }>();
  const { t } = useTranslation();
  const billingService = useMemo(() => getBillingService(), []);

  // 时间范围选择器
  const [timeRange, setTimeRange] = useState<7 | 30 | 90>(7);

  // 获取各种统计数据
  const {
    data: usageStats,
    loading: statsLoading,
    error: statsError,
    refetch: refetchStats,
  } = useOverviewCards(billingService, tenantId);

  const {
    data: modelStats,
    loading: modelLoading,
    error: modelError,
    refetch: refetchModel,
  } = useModelRanking(billingService, tenantId);

  const {
    data: dailyStats,
    loading: dailyLoading,
    error: dailyError,
    refetch: refetchDaily,
  } = useDailyTrend(billingService, tenantId, timeRange);

  const {
    data: botUsage,
    loading: botLoading,
    error: botError,
    refetch: refetchBot,
  } = useBotUsage(billingService, tenantId);

  const {
    data: records,
    loading: recordsLoading,
    error: recordsError,
    refetch: refetchRecords,
  } = useRealtimeRecords(billingService, tenantId);

  // 时间范围变化时刷新数据
  useEffect(() => {
    refetchDaily();
  }, [timeRange, refetchDaily]);

  // 全局刷新
  const handleRefresh = () => {
    refetchStats();
    refetchModel();
    refetchDaily();
    refetchBot();
    refetchRecords();
  };

  // 生成成本趋势图表配置
  const costTrendOption = useMemo(() => {
    if (!dailyStats?.daily_stats) {
      return {};
    }

    const dates = dailyStats.daily_stats.map(item => item.date);
    const costs = dailyStats.daily_stats.map(item => item.total_cost);
    const tokens = dailyStats.daily_stats.map(item => item.total_tokens);

    return {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
        },
      },
      legend: {
        data: [t('billing.totalCost'), t('billing.tokenUsage.totalTokens')],
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: dates,
      },
      yAxis: [
        {
          type: 'value',
          name: t('billing.currency'),
          position: 'left',
        },
        {
          type: 'value',
          name: 'Tokens',
          position: 'right',
        },
      ],
      series: [
        {
          name: t('billing.totalCost'),
          type: 'line',
          data: costs,
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(24, 144, 255, 0.3)' },
              { offset: 1, color: 'rgba(24, 144, 255, 0.05)' },
            ]),
          },
          itemStyle: {
            color: '#1890ff',
          },
        },
        {
          name: t('billing.tokenUsage.totalTokens'),
          type: 'line',
          yAxisIndex: 1,
          data: tokens,
          smooth: true,
          itemStyle: {
            color: '#52c41a',
          },
        },
      ],
    };
  }, [dailyStats, t]);

  // 生成模型使用饼图配置
  const modelPieOption = useMemo(() => {
    if (!modelStats?.model_stats) {
      return {};
    }

    const data = modelStats.model_stats.map(item => ({
      name: `${item.model_provider}/${item.model_name}`,
      value: item.total_tokens,
    }));

    return {
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)',
      },
      legend: {
        orient: 'vertical',
        left: 'left',
      },
      series: [
        {
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2,
          },
          label: {
            show: false,
            position: 'center',
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 20,
              fontWeight: 'bold',
            },
          },
          labelLine: {
            show: false,
          },
          data,
        },
      ],
    };
  }, [modelStats]);

  return (
    <Layout style={{ minHeight: '100vh', background: '#f0f2f5' }}>
      <Header
        style={{
          background: '#fff',
          padding: '16px 24px',
          borderBottom: '1px solid #f0f0f0',
        }}
      >
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Title level={3} style={{ margin: 0 }}>
            {t('billing.tokenUsage.title')}
          </Title>
          <Radio.Group
            value={timeRange}
            onChange={e => setTimeRange(e.target.value)}
            optionType="button"
            buttonStyle="solid"
          >
            <Radio.Button value={7}>{t('billing.tokenUsage.last7Days')}</Radio.Button>
            <Radio.Button value={30}>{t('billing.tokenUsage.last30Days')}</Radio.Button>
            <Radio.Button value={90}>{t('billing.tokenUsage.last90Days')}</Radio.Button>
          </Radio.Group>
        </div>
      </Header>

      <Content style={{ padding: '24px' }}>
        {/* 使用概览卡片 */}
        <Spin spinning={statsLoading}>
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
              gap: '16px',
              marginBottom: '24px',
            }}
          >
            <Card>
              <div style={{ fontSize: '14px', color: '#8c8c8c', marginBottom: '8px' }}>
                {t('billing.tokenUsage.totalTokens')}
              </div>
              <div style={{ fontSize: '24px', fontWeight: 'bold' }}>
                {usageStats?.total_tokens?.toLocaleString() || 0}
              </div>
            </Card>
            <Card>
              <div style={{ fontSize: '14px', color: '#8c8c8c', marginBottom: '8px' }}>
                {t('billing.tokenUsage.totalCost')}
              </div>
              <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#1890ff' }}>
                ${usageStats?.total_cost?.toFixed(2) || '0.00'}
              </div>
            </Card>
            <Card>
              <div style={{ fontSize: '14px', color: '#8c8c8c', marginBottom: '8px' }}>
                {t('billing.tokenUsage.totalRequests')}
              </div>
              <div style={{ fontSize: '24px', fontWeight: 'bold' }}>
                {usageStats?.total_requests?.toLocaleString() || 0}
              </div>
            </Card>
            <Card>
              <div style={{ fontSize: '14px', color: '#8c8c8c', marginBottom: '8px' }}>
                {t('billing.tokenUsage.cacheHitRate')}
              </div>
              <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#52c41a' }}>
                {usageStats?.total_requests
                  ? `${((usageStats.cached_requests / usageStats.total_requests) * 100).toFixed(1)}%`
                  : '0%'}
              </div>
            </Card>
            <Card>
              <div style={{ fontSize: '14px', color: '#8c8c8c', marginBottom: '8px' }}>
                {t('billing.tokenUsage.avgResponseTime')}
              </div>
              <div style={{ fontSize: '24px', fontWeight: 'bold' }}>
                {usageStats?.avg_response_time?.toFixed(0) || 0}ms
              </div>
            </Card>
          </div>
        </Spin>

        {/* 成本趋势图表 */}
        <Card
          title={t('billing.tokenUsage.costTrend')}
          style={{ marginBottom: '24px' }}
          extra={<a onClick={handleRefresh}>{t('common.refresh')}</a>}
        >
          <Spin spinning={dailyLoading}>
            <ReactECharts option={costTrendOption} style={{ height: '350px' }} />
          </Spin>
        </Card>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '24px', marginBottom: '24px' }}>
          {/* 模型使用排行 */}
          <Card title={t('billing.tokenUsage.modelRanking')} style={{ height: '500px' }}>
            <Spin spinning={modelLoading}>
              <div style={{ height: '400px' }}>
                <ReactECharts option={modelPieOption} style={{ height: '100%' }} />
              </div>
            </Spin>
          </Card>

          {/* Bot使用统计 */}
          <Card title={t('billing.tokenUsage.botUsage')} style={{ height: '500px' }}>
            <Spin spinning={botLoading}>
              <div style={{ overflow: 'auto', height: '400px' }}>
                {botUsage?.bot_usage?.map(bot => (
                  <div
                    key={bot.bot_id}
                    style={{
                      padding: '12px',
                      borderBottom: '1px solid #f0f0f0',
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                    }}
                  >
                    <div>
                      <div style={{ fontWeight: 'bold', marginBottom: '4px' }}>{bot.bot_name}</div>
                      <div style={{ fontSize: '12px', color: '#8c8c8c' }}>
                        {bot.request_count} {t('billing.tokenUsage.requests')}
                      </div>
                    </div>
                    <div style={{ textAlign: 'right' }}>
                      <div style={{ fontWeight: 'bold', color: '#1890ff' }}>
                        ${bot.total_cost.toFixed(2)}
                      </div>
                      <div style={{ fontSize: '12px', color: '#8c8c8c' }}>
                        {bot.total_tokens.toLocaleString()} {t('billing.tokenUsage.tokens')}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </Spin>
          </Card>
        </div>

        {/* 实时Token记录 */}
        <Card
          title={t('billing.tokenUsage.realtimeRecords')}
          extra={<a onClick={handleRefresh}>{t('common.refresh')}</a>}
        >
          <Spin spinning={recordsLoading}>
            {records?.records?.length ? (
              <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                <thead>
                  <tr style={{ borderBottom: '2px solid #f0f0f0' }}>
                    <th style={{ padding: '12px', textAlign: 'left' }}>
                      {t('billing.tokenUsage.recordTime')}
                    </th>
                    <th style={{ padding: '12px', textAlign: 'left' }}>
                      {t('billing.tokenUsage.user')}
                    </th>
                    <th style={{ padding: '12px', textAlign: 'left' }}>
                      {t('billing.tokenUsage.botName')}
                    </th>
                    <th style={{ padding: '12px', textAlign: 'left' }}>
                      {t('billing.tokenUsage.model')}
                    </th>
                    <th style={{ padding: '12px', textAlign: 'left' }}>
                      {t('billing.tokenUsage.requestType')}
                    </th>
                    <th style={{ padding: '12px', textAlign: 'right' }}>
                      {t('billing.tokenUsage.tokens')}
                    </th>
                    <th style={{ padding: '12px', textAlign: 'right' }}>
                      {t('billing.tokenUsage.cost')}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {records.records.map(record => (
                    <tr key={record.record_id} style={{ borderBottom: '1px solid #f0f0f0' }}>
                      <td style={{ padding: '12px' }}>
                        {new Date(record.created_at).toLocaleString()}
                      </td>
                      <td style={{ padding: '12px' }}>{record.user_id}</td>
                      <td style={{ padding: '12px' }}>{record.bot_id}</td>
                      <td style={{ padding: '12px' }}>
                        {record.model_provider}/{record.model_name}
                      </td>
                      <td style={{ padding: '12px' }}>{record.request_type}</td>
                      <td style={{ padding: '12px', textAlign: 'right' }}>
                        {record.total_tokens.toLocaleString()}
                      </td>
                      <td style={{ padding: '12px', textAlign: 'right' }}>
                        ${record.cost_usd.toFixed(4)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <div style={{ textAlign: 'center', padding: '40px', color: '#8c8c8c' }}>
                {t('common.noData')}
              </div>
            )}
          </Spin>
        </Card>
      </Content>
    </Layout>
  );
};

export default TokenUsagePage;
