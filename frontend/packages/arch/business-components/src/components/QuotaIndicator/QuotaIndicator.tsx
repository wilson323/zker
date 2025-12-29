// frontend/packages/arch/business-components/src/components/QuotaIndicator/QuotaIndicator.tsx

import React from 'react';
import { useStyles } from './QuotaIndicator.styles';

/**
 * QuotaIndicator组件Props接口
 */
export interface QuotaIndicatorProps {
  /** 已使用量 */
  used: number;
  /** 最大值 */
  max: number;
  /** 资源类型名称 */
  resourceType: string;
  /** 是否显示标签 */
  showLabel?: boolean;
  /** 是否显示百分比 */
  showPercentage?: boolean;
  /** 自定义className */
  className?: string;
}

/**
 * QuotaIndicator 配额指示器
 *
 * 用于显示配额使用情况的进度条组件
 *
 * @example
 * ```tsx
 * <QuotaIndicator
 *   used={750}
 *   max={1000}
 *   resourceType="API调用次数"
 *   showLabel={true}
 *   showPercentage={true}
 * />
 * ```
 */
export const QuotaIndicator: React.FC<QuotaIndicatorProps> = ({
  used,
  max,
  resourceType,
  showLabel = true,
  showPercentage = true,
  className,
}) => {
  const classes = useStyles({ used, max });

  const percentage = max > 0 ? (used / max) * 100 : 0;
  const isOverLimit = percentage >= 100;
  const isNearLimit = percentage >= 80 && percentage < 100;

  const getStatusText = () => {
    if (isOverLimit) return '已超限';
    if (isNearLimit) return '即将超限';
    return '正常';
  };

  const getStatusColor = () => {
    if (isOverLimit) return classes.statusError;
    if (isNearLimit) return classes.statusWarning;
    return classes.statusNormal;
  };

  return (
    <div className={`${classes.container} ${className || ''}`}>
      {showLabel && (
        <div className={classes.header}>
          <span className={classes.label}>{resourceType}</span>
          <span className={`${classes.status} ${getStatusColor()}`}>
            {getStatusText()}
          </span>
        </div>
      )}

      <div className={classes.countRow}>
        <span className={classes.count}>
          {used.toLocaleString()} / {max === -1 ? '∞' : max.toLocaleString()}
        </span>
        {showPercentage && max !== -1 && (
          <span className={`${classes.percentage} ${getStatusColor()}`}>
            {percentage.toFixed(1)}%
          </span>
        )}
      </div>

      {max !== -1 && (
        <div className={classes.barContainer}>
          <div
            className={`${classes.bar} ${
              isOverLimit ? classes.overLimit : isNearLimit ? classes.nearLimit : ''
            }`}
            style={{ width: `${Math.min(percentage, 100)}%` }}
          />
        </div>
      )}
    </div>
  );
};

export default QuotaIndicator;
