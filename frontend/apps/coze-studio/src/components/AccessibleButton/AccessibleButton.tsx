// frontend/apps/coze-studio/src/components/AccessibleButton/AccessibleButton.tsx

/**
 * 可访问按钮组件
 *
 * 特性:
 * - 完整的ARIA属性
 * - 键盘导航支持
 * - 加载状态
 * - 禁用状态
 * - 图标按钮
 */

import React, { ButtonHTMLAttributes } from 'react';
import { getAriaProps } from '../../utils/accessibility';

export interface AccessibleButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /**
   * 按钮文本
   */
  children?: React.ReactNode;

  /**
   * 按钮类型
   */
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';

  /**
   * 按钮大小
   */
  size?: 'small' | 'medium' | 'large';

  /**
   * 是否加载中
   */
  loading?: boolean;

  /**
   * 是否禁用
   */
  disabled?: boolean;

  /**
   * 是否为图标按钮
   */
  icon?: boolean;

  /**
   * 加载时的文本
   */
  loadingText?: string;

  /**
   * ARIA标签(图标按钮时必需)
   */
  ariaLabel?: string;

  /**
   * 描述元素ID
   */
  describedBy?: string;
}

export const AccessibleButton: React.FC<AccessibleButtonProps> = ({
  children,
  variant = 'primary',
  size = 'medium',
  loading = false,
  disabled = false,
  icon = false,
  loadingText = '加载中...',
  ariaLabel,
  describedBy,
  onClick,
  ...rest
}) => {
  const ariaProps = getAriaProps({
    role: 'button',
    label: ariaLabel || (typeof children === 'string' ? children : undefined),
    describedBy,
    disabled: disabled || loading,
  });

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    if (disabled || loading) {
      e.preventDefault();
      return;
    }

    onClick?.(e);
  };

  return (
    <button
      type="button"
      disabled={disabled || loading}
      aria-busy={loading}
      className={`accessibility-button accessibility-button--${variant} accessibility-button--${size}`}
      onClick={handleClick}
      {...ariaProps}
      {...rest}
    >
      {loading ? (
        <>
          <span className="accessibility-button__spinner" aria-hidden="true" />
          <span className="sr-only">{loadingText}</span>
        </>
      ) : (
        children
      )}
    </button>
  );
};

export default AccessibleButton;
