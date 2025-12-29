// frontend/packages/common/components/LoadingWrapper/LoadingWrapper.tsx

import React from 'react';
import { useStyles } from './LoadingWrapper.styles';

/**
 * LoadingWrapper组件Props接口
 */
export interface LoadingWrapperProps {
  /** 加载状态 */
  loading: boolean;
  /** 错误信息 */
  error?: Error | null;
  /** 是否为空 */
  empty?: boolean;
  /** 子元素 */
  children: React.ReactNode;
  /** 自定义加载文字 */
  loadingText?: string;
  /** 自定义空数据文字 */
  emptyText?: string;
  /** 自定义空数据图标 */
  emptyIcon?: React.ReactNode;
}

/**
 * LoadingWrapper 统一加载状态组件
 *
 * 用于统一处理加载中、错误、空数据的显示状态
 *
 * @example
 * ```tsx
 * <LoadingWrapper
 *   loading={isLoading}
 *   error={error}
 *   empty={data.length === 0}
 * >
 *   {data.map(item => <Item key={item.id} {...item} />)}
 * </LoadingWrapper>
 * ```
 */
export const LoadingWrapper: React.FC<LoadingWrapperProps> = ({
  loading,
  error,
  empty = false,
  children,
  loadingText = '加载中...',
  emptyText = '暂无数据',
  emptyIcon,
}) => {
  const classes = useStyles();

  // 加载中状态
  if (loading) {
    return (
      <div className={classes.container}>
        <div className={classes.spinner} />
        <p className={classes.text}>{loadingText}</p>
      </div>
    );
  }

  // 错误状态
  if (error) {
    return (
      <div className={classes.container}>
        <div className={classes.error}>⚠️</div>
        <p className={classes.title}>出错了</p>
        <p className={classes.text}>{error.message}</p>
        <button
          className={classes.retryButton}
          onClick={() => window.location.reload()}
        >
          重新加载
        </button>
      </div>
    );
  }

  // 空数据状态
  if (empty) {
    return (
      <div className={classes.container}>
        <div className={classes.empty}>
          {emptyIcon || '📭'}
        </div>
        <p className={classes.text}>{emptyText}</p>
      </div>
    );
  }

  // 正常状态
  return <>{children}</>;
};

export default LoadingWrapper;
