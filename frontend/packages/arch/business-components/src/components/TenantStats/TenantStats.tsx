// frontend/packages/arch/business-components/src/components/TenantStats/TenantStats.tsx

import React, { CSSProperties } from 'react';
import { useStyles } from './TenantStats.styles';
import { TenantStatsDTO } from '@coze-studio/api-client';

export interface TenantStatsProps {
  /**
   * 租户统计数据
   */
  stats: TenantStatsDTO;

  /**
   * 是否显示图表
   */
  showChart?: boolean;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;
}

export const TenantStats: React.FC<TenantStatsProps> = ({
  stats,
  showChart = true,
  className,
  style,
}) => {
  const classes = useStyles();

  const statItems = [
    {
      label: '机器人总数',
      value: stats.total_bots,
      active: stats.active_bots,
      icon: '🤖',
      color: '#1890ff',
    },
    {
      label: '用户总数',
      value: stats.total_users,
      active: stats.active_users,
      icon: '👥',
      color: '#52c41a',
    },
    {
      label: '消息总数',
      value: stats.total_messages,
      icon: '💬',
      color: '#722ed1',
    },
    {
      label: 'Token使用量',
      value: stats.total_tokens_used,
      icon: '📊',
      color: '#fa8c16',
    },
    {
      label: '存储使用量',
      value: `${stats.storage_used_mb} MB`,
      icon: '💾',
      color: '#13c2c2',
    },
  ];

  const calculatePercentage = (active: number, total: number): number => {
    return total > 0 ? Math.round((active / total) * 100) : 0;
  };

  return (
    <div className={`${classes.container} ${className || ''}`} style={style}>
      {statItems.map((item) => (
        <div key={item.label} className={classes.statCard}>
          <div className={classes.header}>
            <span className={classes.icon}>{item.icon}</span>
            <span className={classes.label}>{item.label}</span>
          </div>

          <div className={classes.value}>{item.value}</div>

          {item.active !== undefined && (
            <div className={classes.active}>
              <span className={classes.activeLabel}>活跃：</span>
              <span className={classes.activeValue}>{item.active}</span>
              <span className={classes.percentage}>
                ({calculatePercentage(item.active, item.value as number)}%)
              </span>
            </div>
          )}

          {showChart && (
            <div className={classes.bar}>
              <div
                className={classes.barFill}
                style={{
                  width:
                    item.active !== undefined
                      ? `${calculatePercentage(item.active, item.value as number)}%`
                      : '100%',
                  backgroundColor: item.color,
                }}
              />
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default TenantStats;
