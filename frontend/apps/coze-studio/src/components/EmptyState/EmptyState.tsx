// frontend/apps/coze-studio/src/components/EmptyState/EmptyState.tsx

/**
 * 空状态组件
 *
 * 特性:
 * - 多种预设类型
 * - 自定义图标和图片
 * - 操作按钮
 * - 可访问性支持
 */

import React from 'react';
import classNames from 'classnames';

export interface EmptyStateProps {
  /**
   * 空状态类型
   */
  type?: 'noData' | 'noResults' | 'error' | 'networkError' | 'custom';

  /**
   * 标题
   */
  title?: string;

  /**
   * 描述
   */
  description?: string;

  /**
   * 图片URL
   */
  image?: string;

  /**
   * 图标
   */
  icon?: React.ReactNode;

  /**
   * 操作按钮
   */
  action?: React.ReactNode;

  /**
   * 自定义类名
   */
  className?: string;
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  type = 'noData',
  title,
  description,
  image,
  icon,
  action,
  className,
}) => {
  const titles = {
    noData: '暂无数据',
    noResults: '未找到匹配的结果',
    error: '出错了',
    networkError: '网络连接失败',
    custom: title || '',
  };

  const descriptions = {
    noData: '这里还没有任何内容',
    noResults: '请尝试调整搜索条件',
    error: '请稍后再试',
    networkError: '请检查您的网络连接',
    custom: description || '',
  };

  const displayTitle = title || titles[type];
  const displayDescription = description || descriptions[type];

  const emptyStateClass = classNames(
    'empty-state',
    `empty-state--${type}`,
    className
  );

  return (
    <div className={emptyStateClass} role="status" aria-live="polite">
      {image && (
        <div className="empty-state__image">
          <img src={image} alt="" aria-hidden="true" />
        </div>
      )}

      {icon && <div className="empty-state__icon">{icon}</div>}

      <h3 className="empty-state__title">{displayTitle}</h3>

      {displayDescription && (
        <p className="empty-state__description">{displayDescription}</p>
      )}

      {action && <div className="empty-state__action">{action}</div>}
    </div>
  );
};

export default EmptyState;
