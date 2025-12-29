// frontend/packages/arch/ui-components/src/components/Input/Input.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  transitions,
  typography,
} from '@coze-studio/common/themes';

/**
 * Input样式Props接口
 */
interface InputStyleProps {
  size: 'sm' | 'md' | 'lg';
  error: boolean;
  block: boolean;
  hasPrefix: boolean;
  hasSuffix: boolean;
}

/**
 * Input样式Hook
 */
export const useStyles = (props: InputStyleProps) => {
  const { size, error, block, hasPrefix, hasSuffix } = props;

  // ==================== 尺寸样式 ====================

  const sizeStyles = {
    sm: css`
      height: 32px;
      padding: ${spacing[1]} ${spacing[3]};
      ${typography.fontSize.sm};
    `,
    md: css`
      height: 40px;
      padding: ${spacing[2]} ${spacing[3]};
      ${typography.fontSize.base};
    `,
    lg: css`
      height: 48px;
      padding: ${spacing[3]} ${spacing[4]};
      ${typography.fontSize.lg};
    `,
  };

  // ==================== 边框样式 ====================

  const borderStyles = css`
    border: 1px solid ${error ? colors.error.main : colors.gray[300]};

    &:hover:not(:disabled) {
      border-color: ${error ? colors.error.dark : colors.primary[500]};
    }

    &:focus {
      border-color: ${error ? colors.error.dark : colors.primary[500]};
      box-shadow: 0 0 0 3px ${error ? 'rgba(245, 34, 45, 0.1)' : 'rgba(24, 144, 255, 0.1)'};
    }

    &:disabled {
      background-color: ${colors.gray[50]};
      cursor: not-allowed;
    }
  `;

  // ==================== 返回样式对象 ====================

  return {
    wrapper: css`
      display: ${block ? 'flex' : 'inline-flex'};
      width: ${block ? '100%' : 'auto'};
      align-items: center;
      position: relative;
    `,
    input: css`
      width: 100%;
      border-radius: ${borderRadius.base};
      outline: none;
      transition: all ${transitions.base};
      background-color: white;

      ${hasPrefix && 'padding-left: 40px;'}
      ${hasSuffix && 'padding-right: 40px;'}

      ${sizeStyles[size]}
      ${borderStyles}

      &::placeholder {
        color: ${colors.gray[400]};
      }
    `,
    prefix: css`
      position: absolute;
      left: ${spacing[3]};
      display: flex;
      align-items: center;
      color: ${colors.gray[400]};
      pointer-events: none;
    `,
    suffix: css`
      position: absolute;
      right: ${spacing[3]};
      display: flex;
      align-items: center;
      color: ${colors.gray[400]};
      pointer-events: none;
    `,
  };
};
