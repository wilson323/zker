// frontend/packages/common/components/ErrorBoundary/ErrorBoundary.tsx

import React, { Component, ErrorInfo, ReactNode } from 'react';
import { Button } from '@coze-studio/ui-components';
import { useStyles } from './ErrorBoundary.styles';

/**
 * ErrorBoundary状态接口
 */
interface State {
  hasError: boolean;
  error?: Error;
  errorInfo?: ErrorInfo;
}

/**
 * ErrorBoundary组件Props接口
 */
interface Props {
  children: ReactNode;
  /** 自定义错误回调 */
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
  /** 自定义fallback UI */
  fallback?: ReactNode;
}

/**
 * ErrorBoundary 错误边界
 *
 * 捕获子组件树的JavaScript错误，记录错误日志，并显示备用UI
 *
 * @example
 * ```tsx
 * <ErrorBoundary onError={(error, errorInfo) => {
 *   logErrorToService(error, errorInfo);
 * }}>
 *   <App />
 * </ErrorBoundary>
 * ```
 */
export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // 记录错误到日志服务
    console.error('ErrorBoundary caught an error:', error, errorInfo);

    // 调用自定义错误回调
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }

    // 更新state
    this.setState({
      error,
      errorInfo,
    });
  }

  handleReset = () => {
    this.setState({ hasError: false, error: undefined, errorInfo: undefined });
    window.location.reload();
  };

  handleGoHome = () => {
    this.setState({ hasError: false, error: undefined, errorInfo: undefined });
    window.location.href = '/';
  };

  render() {
    const classes = useStyles();

    if (this.state.hasError) {
      // 如果提供了自定义fallback，使用它
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // 默认错误UI
      return (
        <div className={classes.container}>
          <div className={classes.content}>
            <div className={classes.icon}>⚠️</div>
            <h1 className={classes.title}>出错了</h1>
            <p className={classes.description}>
              应用程序遇到了一些问题，请稍后重试
            </p>

            {/* 开发环境显示错误详情 */}
            {process.env.NODE_ENV === 'development' && this.state.error && (
              <details className={classes.details}>
                <summary className={classes.summary}>错误详情</summary>
                <div className={classes.errorInfo}>
                  <p><strong>错误消息:</strong></p>
                  <pre className={classes.pre}>
                    {this.state.error.toString()}
                  </pre>

                  {this.state.errorInfo && (
                    <>
                      <p><strong>组件堆栈:</strong></p>
                      <pre className={classes.pre}>
                        {this.state.errorInfo.componentStack}
                      </pre>
                    </>
                  )}
                </div>
              </details>
            )}

            <div className={classes.actions}>
              <Button variant="primary" onClick={this.handleReset}>
                重新加载
              </Button>
              <Button variant="outline" onClick={this.handleGoHome}>
                返回首页
              </Button>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
