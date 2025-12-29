// frontend/packages/arch/ui-components/src/components/Progress/Progress.tsx

import React, { CSSProperties } from 'react';
import { useStyles } from './Progress.styles';

export type ProgressType = 'line' | 'circle' | 'dashboard';

export interface ProgressProps {
  /**
   * 百分比（0-100）
   */
  percent: number;

  /**
   * 类型
   */
  type?: ProgressType;

  /**
   * 状态
   */
  status?: 'normal' | 'active' | 'success' | 'exception';

  /**
   * 是否显示百分比数字
   */
  showInfo?: boolean;

  /**
   * 进度条线的宽度
   */
  strokeWidth?: number;

  /**
   * 线型进度条的宽度（像素）
   */
  width?: number;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;

  /**
   * 自定义格式化内容
   */
  format?: (percent: number, successPercent: number) => React.ReactNode;
}

export const Progress: React.FC<ProgressProps> = ({
  percent = 0,
  type = 'line',
  status = 'normal',
  showInfo = true,
  strokeWidth = 10,
  width = 120,
  className,
  style,
  format,
}) => {
  const classes = useStyles({ type, status, width, strokeWidth });

  // 确保百分比在 0-100 之间
  const normalizedPercent = Math.max(0, Math.min(100, percent));

  // 格式化显示内容
  const renderInfo = () => {
    if (!showInfo) return null;

    if (format) {
      return <div className={classes.info}>{format(normalizedPercent, 0)}</div>;
    }

    if (status === 'exception') {
      return <div className={classes.info}>✗</div>;
    }

    if (status === 'success') {
      return <div className={classes.info}>✓</div>;
    }

    return <div className={classes.info}>{`${Math.round(normalizedPercent)}%`}</div>;
  };

  // 渲染线型进度条
  if (type === 'line') {
    return (
      <div className={`${classes.container} ${className || ''}`} style={style}>
        <div className={classes.outer}>
          <div className={classes.inner}>
            <div
              className={classes.bg}
              style={{ width: `${normalizedPercent}%` }}
            />
            {status === 'active' && (
              <div className={classes.activeBg} style={{ width: `${normalizedPercent}%` }} />
            )}
          </div>
        </div>
        {showInfo && renderInfo()}
      </div>
    );
  }

  // 渲染圆形进度条
  const radius = 50 - strokeWidth / 2;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (normalizedPercent / 100) * circumference;

  return (
    <div
      className={`${classes.circleContainer} ${className || ''}`}
      style={style}
    >
      <svg width={width} height={width} viewBox="0 0 100 100">
        <circle
          className={classes.circleTrail}
          cx={50}
          cy={50}
          r={radius}
          fill="none"
          strokeWidth={strokeWidth}
        />
        <circle
          className={classes.circlePath}
          cx={50}
          cy={50}
          r={radius}
          fill="none"
          strokeWidth={strokeWidth}
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          style={{
            transform: 'rotate(-90deg)',
            transformOrigin: '50% 50%',
          }}
        />
      </svg>
      {renderInfo()}
    </div>
  );
};

export default Progress;
