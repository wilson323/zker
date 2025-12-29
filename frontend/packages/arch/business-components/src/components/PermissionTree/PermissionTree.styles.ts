// frontend/packages/arch/business-components/src/components/PermissionTree/PermissionTree.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  transitions,
} from '@coze-studio/common/themes';

/**
 * PermissionTree样式Hook
 */
export const useStyles = () => {
  // ==================== 返回样式对象 ====================

  return {
    container: css`
      width: 100%;
      background-color: white;
      border: 1px solid ${colors.gray[200]};
      border-radius: 4px;
      padding: ${spacing[3]};
      max-height: 400px;
      overflow-y: auto;
    `,
    node: css`
      margin-bottom: ${spacing[1]};
    `,
    nodeContent: css`
      display: flex;
      align-items: center;
      padding: ${spacing[1]} 0;
      transition: background-color ${transitions.fast};
      border-radius: 4px;

      &:hover {
        background-color: ${colors.gray[50]};
      }
    `,
    expandIcon: css`
      width: 20px;
      height: 20px;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      margin-right: ${spacing[2]};
      font-size: 10px;
      color: ${colors.gray[400]};
      cursor: pointer;
      transition: transform ${transitions.fast}, color ${transitions.fast};

      &:hover {
        color: ${colors.gray[600]};
      }
    `,
    expanded: css`
      transform: rotate(90deg);
    `,
    checkbox: css`
      width: 16px;
      height: 16px;
      margin-right: ${spacing[2]};
      cursor: pointer;
      accent-color: ${colors.primary[500]};
    `,
    nodeTitle: css`
      flex: 1;
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[700]};
      user-select: none;
    `,
    disabled: css`
      color: ${colors.gray[400]};
      cursor: not-allowed;
    `,
    children: css`
      margin-left: 0;
    `,
  };
};
