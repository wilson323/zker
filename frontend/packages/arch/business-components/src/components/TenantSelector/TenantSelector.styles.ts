// frontend/packages/arch/business-components/src/components/TenantSelector/TenantSelector.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  transitions,
  typography,
} from '@coze-studio/common/themes';

/**
 * TenantSelector样式Hook
 */
export const useStyles = () => {
  // ==================== 返回样式对象 ====================

  return {
    container: css`
      position: relative;
      display: inline-block;
      width: 100%;
    `,
    select: css`
      width: 100%;
      height: 40px;
      padding: ${spacing[2]} ${spacing[3]};
      padding-right: ${spacing[10]};
      border: 1px solid ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      background-color: white;
      color: ${colors.gray[700]};
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.normal};
      cursor: pointer;
      transition: all ${transitions.base};
      outline: none;
      appearance: none;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%238C8C8C' d='M6 9L1 4h10z'/%3E%3C/svg%3E");
      background-repeat: no-repeat;
      background-position: right 12px center;

      &:hover:not(:disabled) {
        border-color: ${colors.primary[500]};
      }

      &:focus {
        border-color: ${colors.primary[500]};
        box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.1);
      }

      &:disabled {
        background-color: ${colors.gray[50]};
        cursor: not-allowed;
        opacity: 0.6;
      }

      &::placeholder {
        color: ${colors.gray[400]};
      }
    `,
    clearButton: css`
      position: absolute;
      right: 36px;
      top: 50%;
      transform: translateY(-50%);
      width: 16px;
      height: 16px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: ${colors.gray[300]};
      border: none;
      border-radius: 50%;
      color: ${colors.gray[600]};
      font-size: 14px;
      cursor: pointer;
      transition: all ${transitions.fast};
      line-height: 1;

      &:hover {
        background: ${colors.gray[400]};
        color: ${colors.gray[800]};
      }

      &:active {
        background: ${colors.gray[500]};
      }
    `,
    loading: css`
      position: absolute;
      right: 12px;
      top: 50%;
      transform: translateY(-50%);
      width: 16px;
      height: 16px;
      border: 2px solid ${colors.gray[200]};
      border-top-color: ${colors.primary[500]};
      border-radius: 50%;
      animation: spin 0.6s linear infinite;

      @keyframes spin {
        to {
          transform: translateY(-50%) rotate(360deg);
        }
      }
    `,
    // 标签样式（用于显示租户类型）
    tagBlue: css`
      color: ${colors.primary[600]};
    `,
    tagGreen: css`
      color: ${colors.success.main};
    `,
    tagPurple: css`
      color: #722ed1;
    `,
    tagGray: css`
      color: ${colors.gray[600]};
    `,
  };
};
