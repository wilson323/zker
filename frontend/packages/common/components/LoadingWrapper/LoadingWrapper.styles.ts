// frontend/packages/common/components/LoadingWrapper/LoadingWrapper.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  borderRadius,
} from '@coze-studio/common/themes';

/**
 * LoadingWrapper样式Hook
 */
export const useStyles = () => {
  return {
    container: css`
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: ${spacing[8]};
      min-height: 400px;
    `,
    spinner: css`
      width: 48px;
      height: 48px;
      border: 4px solid ${colors.gray[200]};
      border-top-color: ${colors.primary[500]};
      border-radius: 50%;
      animation: spin 0.8s linear infinite;

      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }
    `,
    error: css`
      font-size: 48px;
      margin-bottom: ${spacing[4]};
    `,
    empty: css`
      font-size: 64px;
      margin-bottom: ${spacing[4]};
      opacity: 0.5;
    `,
    title: css`
      margin: 0 0 ${spacing[2]} 0;
      font-size: ${typography.fontSize.lg};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
    text: css`
      margin: 0;
      font-size: ${typography.fontSize.base};
      color: ${colors.gray[500]};
      text-align: center;
    `,
    retryButton: css`
      margin-top: ${spacing[4]};
      padding: ${spacing[2]} ${spacing[4]};
      border: 1px solid ${colors.primary[500]};
      background-color: white;
      color: ${colors.primary[500]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      cursor: pointer;
      outline: none;
      transition: all 0.2s;

      &:hover {
        background-color: ${colors.primary[500]};
        color: white;
      }

      &:active {
        background-color: ${colors.primary[600]};
      }
    `,
  };
};
