// frontend/packages/arch/ui-components/src/components/Button/Button.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  transitions,
  typography,
} from '@coze-studio/common/themes';

/**
 * Button样式Props接口
 */
interface ButtonStyleProps {
  variant: 'primary' | 'secondary' | 'outline' | 'text' | 'danger';
  size: 'sm' | 'md' | 'lg';
  block: boolean;
}

/**
 * Button样式Hook
 */
export const useStyles = (props: ButtonStyleProps) => {
  const { variant, size, block } = props;

  // ==================== 尺寸样式 ====================

  const sizeStyles = {
    sm: css`
      padding: ${spacing[2]} ${spacing[3]};
      ${typography.fontSize.sm};
    `,
    md: css`
      padding: ${spacing[3]} ${spacing[4]};
      ${typography.fontSize.base};
    `,
    lg: css`
      padding: ${spacing[4]} ${spacing[6]};
      ${typography.fontSize.lg};
    `,
  };

  // ==================== 变体样式 ====================

  const variantStyles = {
    primary: css`
      background-color: ${colors.primary[500]};
      color: white;
      border: 1px solid ${colors.primary[500]};

      &:hover:not(:disabled) {
        background-color: ${colors.primary[600]};
        border-color: ${colors.primary[600]};
      }

      &:active:not(:disabled) {
        background-color: ${colors.primary[700]};
        border-color: ${colors.primary[700]};
      }
    `,
    secondary: css`
      background-color: ${colors.gray[100]};
      color: ${colors.gray[800]};
      border: 1px solid ${colors.gray[200]};

      &:hover:not(:disabled) {
        background-color: ${colors.gray[200]};
        border-color: ${colors.gray[300]};
      }

      &:active:not(:disabled) {
        background-color: ${colors.gray[300]};
      }
    `,
    outline: css`
      background-color: transparent;
      border: 1px solid ${colors.gray[300]};
      color: ${colors.gray[700]};

      &:hover:not(:disabled) {
        border-color: ${colors.primary[500]};
        color: ${colors.primary[500]};
        background-color: ${colors.primary[50]};
      }

      &:active:not(:disabled) {
        background-color: ${colors.primary[100]};
      }
    `,
    text: css`
      background-color: transparent;
      border: 1px solid transparent;
      color: ${colors.primary[500]};

      &:hover:not(:disabled) {
        background-color: ${colors.primary[50]};
      }

      &:active:not(:disabled) {
        background-color: ${colors.primary[100]};
      }
    `,
    danger: css`
      background-color: ${colors.error.main};
      color: white;
      border: 1px solid ${colors.error.main};

      &:hover:not(:disabled) {
        background-color: ${colors.error.dark};
        border-color: ${colors.error.dark};
      }

      &:active:not(:disabled) {
        background-color: ${colors.error.dark};
        opacity: 0.8;
      }
    `,
  };

  // ==================== 返回样式对象 ====================

  return {
    button: css`
      display: ${block ? 'block' : 'inline-flex'};
      width: ${block ? '100%' : 'auto'};
      align-items: center;
      justify-content: center;
      gap: ${spacing[2]};
      border: none;
      border-radius: ${borderRadius.base};
      font-weight: ${typography.fontWeight.medium};
      cursor: pointer;
      transition: all ${transitions.base};
      outline: none;
      white-space: nowrap;

      &:focus-visible {
        outline: 2px solid ${colors.primary[400]};
        outline-offset: 2px;
      }

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      ${sizeStyles[size]}
      ${variantStyles[variant]}
    `,
    spinner: css`
      width: 1em;
      height: 1em;
      border: 2px solid currentColor;
      border-top-color: transparent;
      border-radius: 50%;
      animation: spin 0.6s linear infinite;

      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }
    `,
    icon: css`
      display: inline-flex;
      align-items: center;
      justify-content: center;
    `,
  };
};
