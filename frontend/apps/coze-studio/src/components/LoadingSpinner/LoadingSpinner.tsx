// frontend/apps/coze-studio/src/components/LoadingSpinner/LoadingSpinner.tsx

/**
 * 加载指示器组件
 *
 * 特性:
 * - 多种尺寸
 * - 全屏/内联模式
 * - 自定义颜色
 * - 可访问性支持
 */

import React from 'react';
import classNames from 'classnames';

export interface LoadingSpinnerProps {
  /**
   * 尺寸
   */
  size?: 'small' | 'medium' | 'large';

  /**
   * 颜色
   */
  color?: 'primary' | 'secondary' | 'white';

  /**
   * 是否全屏
   */
  fullscreen?: boolean;

  /**
   * 加载文本
   */
  text?: string;

  /**
   * 自定义类名
   */
  className?: string;

  /**
   * ARIA标签
   */
  ariaLabel?: string;
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'medium',
  color = 'primary',
  fullscreen = false,
  text,
  className,
  ariaLabel = '加载中',
}) => {
  const spinnerClass = classNames(
    'loading-spinner',
    `loading-spinner--${size}`,
    `loading-spinner--${color}`,
    fullscreen && 'loading-spinner--fullscreen',
    className
  );

  const content = (
    <div className={spinnerClass} role="status" aria-live="polite">
      <div className="loading-spinner__circle" aria-hidden="true">
        <div className="loading-spinner__dot"></div>
        <div className="loading-spinner__dot"></div>
        <div className="loading-spinner__dot"></div>
      </div>
      {text && <p className="loading-spinner__text">{text}</p>}
      <span className="sr-only">{ariaLabel}</span>
    </div>
  );

  if (fullscreen) {
    return <div className="loading-spinner-overlay">{content}</div>;
  }

  return content;
};

export default LoadingSpinner;
