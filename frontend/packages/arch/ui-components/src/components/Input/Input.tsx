// frontend/packages/arch/ui-components/src/components/Input/Input.tsx

import React from 'react';
import { InputHTMLAttributes, forwardRef } from 'react';
import { useStyles } from './Input.styles';

/**
 * Input组件Props接口
 */
export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  /** 输入框尺寸 */
  size?: 'sm' | 'md' | 'lg';
  /** 错误状态 */
  error?: boolean;
  /** 前缀图标 */
  prefix?: React.ReactNode;
  /** 后缀图标 */
  suffix?: React.ReactNode;
  /** 块级输入框（占满容器宽度） */
  block?: boolean;
}

/**
 * Input 输入框
 *
 * 基础输入框组件，支持多种尺寸和状态
 *
 * @example
 * ```tsx
 * <Input
 *   size="md"
 *   placeholder="请输入用户名"
 *   error={hasError}
 *   prefix={<UserIcon />}
 * />
 * ```
 */
export const Input = forwardRef<HTMLInputElement, InputProps>(({
  size = 'md',
  error = false,
  prefix,
  suffix,
  block = false,
  className,
  disabled,
  ...rest
}, ref) => {
  const classes = useStyles({ size, error, block, hasPrefix: !!prefix, hasSuffix: !!suffix });

  return (
    <div className={`${classes.wrapper} ${className || ''}`}>
      {prefix && <span className={classes.prefix}>{prefix}</span>}
      <input
        ref={ref}
        className={classes.input}
        disabled={disabled}
        {...rest}
      />
      {suffix && <span className={classes.suffix}>{suffix}</span>}
    </div>
  );
});

Input.displayName = 'Input';

export default Input;
