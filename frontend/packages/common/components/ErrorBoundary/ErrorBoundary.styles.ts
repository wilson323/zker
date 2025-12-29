// frontend/packages/common/components/ErrorBoundary/ErrorBoundary.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  borderRadius,
} from '@coze-studio/common/themes';

/**
 * ErrorBoundary样式Hook
 */
export const useStyles = () => {
  return {
    container: css`
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      padding: ${spacing[5]};
      background-color: ${colors.gray[50]};
    `,
    content: css`
      max-width: 600px;
      padding: ${spacing[8]};
      background-color: white;
      border-radius: ${borderRadius.lg};
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
      text-align: center;
    `,
    icon: css`
      font-size: 64px;
      margin-bottom: ${spacing[5]};
    `,
    title: css`
      margin: 0 0 ${spacing[3]} 0;
      font-size: ${typography.fontSize['2xl']};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
    description: css`
      margin: 0 0 ${spacing[6]} 0;
      font-size: ${typography.fontSize.base};
      color: ${colors.gray[500]};
      line-height: ${typography.lineHeight.relaxed};
    `,
    details: css`
      margin: ${spacing[5]} 0;
      text-align: left;
      padding: ${spacing[4]};
      background-color: ${colors.gray[50]};
      border-radius: ${borderRadius.base};
    `,
    summary: css`
      cursor: pointer;
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[700]};
      outline: none;

      &:hover {
        color: ${colors.primary[500]};
      }
    `,
    errorInfo: css`
      margin-top: ${spacing[4]};
    `,
    pre: css`
      margin: ${spacing[2]} 0;
      padding: ${spacing[3]};
      background-color: ${colors.gray[900]};
      color: ${colors.success.light};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.sm};
      font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace;
      overflow-x: auto;
      white-space: pre-wrap;
      word-break: break-all;
    `,
    actions: css`
      display: flex;
      gap: ${spacing[3]};
      justify-content: center;
    `,
  };
};
