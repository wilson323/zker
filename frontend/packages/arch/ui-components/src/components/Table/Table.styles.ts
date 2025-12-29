// frontend/packages/arch/ui-components/src/components/Table/Table.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  transitions,
  typography,
} from '@coze-studio/common/themes';

/**
 * Table样式Props接口
 */
interface TableStyleProps {
  bordered: boolean;
}

/**
 * Table样式Hook
 */
export const useStyles = (props: TableStyleProps) => {
  const { bordered } = props;

  // ==================== 返回样式对象 ====================

  return {
    container: css`
      width: 100%;
    `,
    tableWrapper: css`
      width: 100%;
      overflow-x: auto;
      background-color: white;
      border-radius: ${borderRadius.base};
      ${bordered ? `border: 1px solid ${colors.gray[200]};` : ''}
    `,
    table: css`
      width: 100%;
      border-collapse: collapse;
      border-spacing: 0;
    `,
    thead: css`
      background-color: ${colors.gray[50];
    `,
    th: css`
      padding: ${spacing[3]} ${spacing[4]};
      text-align: left;
      font-weight: ${typography.fontWeight.semibold};
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[700]};
      border-bottom: 1px solid ${colors.gray[200]};
      white-space: nowrap;

      &:first-child {
        border-left: none;
      }

      &:last-child {
        border-right: none;
      }
    `,
    tr: css`
      transition: background-color ${transitions.fast};

      &:hover {
        background-color: ${colors.gray[50];
      }
    `,
    td: css`
      padding: ${spacing[3]} ${spacing[4]};
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[700]};
      border-bottom: 1px solid ${colors.gray[100]};

      &:first-child {
        border-left: none;
      }

      &:last-child {
        border-right: none;
      }
    `,
    loadingCell: css`
      padding: ${spacing[8]} !important;
      text-align: center;
      border-bottom: none !important;
    `,
    spinner: css`
      display: inline-block;
      width: 24px;
      height: 24px;
      border: 3px solid ${colors.gray[200]};
      border-top-color: ${colors.primary[500]};
      border-radius: 50%;
      animation: spin 0.6s linear infinite;
      margin-bottom: ${spacing[2]};

      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }
    `,
    loadingText: css`
      display: block;
      margin-top: ${spacing[2]};
      color: ${colors.gray[500]};
      font-size: ${typography.fontSize.sm};
    `,
    emptyCell: css`
      padding: ${spacing[8]} !important;
      text-align: center;
      color: ${colors.gray[400]};
      font-size: ${typography.fontSize.base};
      border-bottom: none !important;
    `,
    pagination: css`
      display: flex;
      align-items: center;
      justify-content: center;
      gap: ${spacing[2]};
      padding: ${spacing[4]};
      background-color: white;
      ${bordered ? `border: 1px solid ${colors.gray[200]};` : ''}
      border-top: none;
      border-radius: 0 0 ${borderRadius.base} ${borderRadius.base};
    `,
    pageButton: css`
      min-width: 32px;
      height: 32px;
      padding: 0 ${spacing[2]};
      border: 1px solid ${colors.gray[200]};
      background-color: white;
      color: ${colors.gray[700]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.sm};
      cursor: pointer;
      transition: all ${transitions.fast};

      &:hover:not(:disabled) {
        border-color: ${colors.primary[500]};
        color: ${colors.primary[500]};
      }

      &:active:not(:disabled) {
        background-color: ${colors.primary[50]};
      }

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }
    `,
    activePage: css`
      background-color: ${colors.primary[500]} !important;
      color: white !important;
      border-color: ${colors.primary[500]} !important;

      &:hover {
        background-color: ${colors.primary[600]} !important;
        border-color: ${colors.primary[600]} !important;
      }
    `,
    ellipsis: css`
      padding: 0 ${spacing[2]};
      color: ${colors.gray[400]};
      font-size: ${typography.fontSize.sm};
    `,
  };
};
