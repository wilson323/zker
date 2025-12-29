// frontend/packages/arch/ui-components/src/components/Button/Button.tsx

import React from 'react';
import { ButtonHTMLAttributes } from 'react';
import { useStyles } from './Button.styles';

/**
 * Button组件Props接口
 */
export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** 按钮变体 */
  variant?: 'primary' | 'secondary' | 'outline' | 'text' | 'danger';
  /** 按钮尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 加载状态 */
  loading?: boolean;
  /** 禁用状态 */
  disabled?: boolean;
  /** 图标 */
  icon?: React.ReactNode;
  /** 块级按钮（占满容器宽度） */
  block?: boolean;
  /** 子元素 */
  children: React.ReactNode;
}

/**
 * Button 按钮
 *
 * 基础按钮组件，支持多种变体、尺寸和状态
 *
 * @example
 * ```tsx
 * <Button variant="primary" size="md" loading>
 *   提交
 * </Button>
 * ```
 */
export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  loading = false,
  disabled = false,
  icon,
  block = false,
  children,
  className,
  onClick,
  ...rest
}) => {
  const classes = useStyles({ variant, size, block });

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    if (loading || disabled) {
      e.preventDefault();
      return;
    }
    onClick?.(e);
  };

  return (
    <button
      className={`${classes.button} ${className || ''}`}
      disabled={disabled || loading}
      onClick={handleClick}
      {...rest}
    >
      {loading && <span className={classes.spinner} />}
      {icon && !loading && <span className={classes.icon}>{icon}</span>}
      {children}
    </button>
  );
};

export default Button;
