// frontend/packages/studio/src/pages/tenant/TenantList/components/TenantFilter.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
} from '@coze-studio/common/themes';

/**
 * TenantFilter样式Hook
 */
export const useStyles = () => {
  return {
    container: css`
      padding: ${spacing[4]};
      background-color: white;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.base};
      margin-bottom: ${spacing[4]};
    `,
    filterGroup: css`
      display: flex;
      gap: ${spacing[4]};
      flex-wrap: wrap;
    `,
    filterItem: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[2]};
      min-width: 200px;
    `,
    label: css`
      font-size: 14px;
      font-weight: 500;
      color: ${colors.gray[700]};
    `,
    select: css`
      height: 40px;
      padding: ${spacing[2]} ${spacing[3]};
      border: 1px solid ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      background-color: white;
      color: ${colors.gray[700]};
      font-size: 14px;
      cursor: pointer;
      outline: none;
      transition: border-color 0.2s;

      &:focus {
        border-color: ${colors.primary[500]};
      }
    `,
    resetButton: css`
      height: 40px;
      padding: ${spacing[2]} ${spacing[4]};
      border: 1px solid ${colors.gray[300]};
      background-color: white;
      color: ${colors.gray[700]};
      border-radius: ${borderRadius.base};
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      outline: none;
      transition: all 0.2s;

      &:hover {
        border-color: ${colors.primary[500]};
        color: ${colors.primary[500]};
      }

      &:active {
        background-color: ${colors.gray[50]};
      }
    `,
  };
};
