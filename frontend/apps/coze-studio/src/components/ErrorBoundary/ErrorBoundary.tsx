// frontend/apps/coze-studio/src/components/ErrorBoundary/ErrorBoundary.tsx

/**
 * 错误边界组件
 *
 * 特性:
 * - 捕获子组件错误
 * - 显示友好错误信息
 * - 提供重试机制
 * - 错误上报
 */

import React, { Component, ErrorInfo, ReactNode } from 'react';

export interface ErrorBoundaryProps {
  children: ReactNode;

  /**
   * 自定义错误展示
   */
  fallback?: ReactNode;

  /**
   * 错误回调
   */
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

export interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // 记录错误到控制台
    console.error('ErrorBoundary caught an error:', error, errorInfo);

    // 调用错误回调
    this.props.onError?.(error, errorInfo);

    // 上报错误到错误监控服务
    // 例如: Sentry, LogRocket等
    // logErrorToService(error, errorInfo);
  }

  handleReset = () => {
    this.setState({ hasError: false, error: undefined });
  };

  render() {
    if (this.state.hasError) {
      // 使用自定义fallback或默认错误UI
      if (this.props.fallback) {
        return this.props.fallback;
      }

      return (
        <div className="error-boundary" role="alert">
          <div className="error-boundary__content">
            <h2 className="error-boundary__title">出错了</h2>
            <p className="error-boundary__message">
              {this.state.error?.message || '页面加载出现问题'}
            </p>
            <div className="error-boundary__actions">
              <button onClick={this.handleReset} className="error-boundary__retry">
                重试
              </button>
              <button
                onClick={() => window.location.reload()}
                className="error-boundary__refresh"
              >
                刷新页面
              </button>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
