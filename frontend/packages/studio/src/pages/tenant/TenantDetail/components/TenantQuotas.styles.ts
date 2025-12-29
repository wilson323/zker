// frontend/packages/studio/src/pages/tenant/TenantDetail/components/TenantQuotas.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  typography,
} from '@coze-studio/common/themes';

/**
 * TenantQuotas样式Hook
 */
export const useStyles = () => {
  return {
    container: css`
      padding: ${spacing[5]};
      background-color: white;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.base};
    `,
    title: css`
      margin: 0 0 ${spacing[5]} 0;
      font-size: ${typography.fontSize.xl};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
    quotaList: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[4]};
    `,
    quotaItem: css`
      padding: ${spacing[4]};
      background-color: ${colors.gray[50]};
      border: 1px solid ${colors.gray[100]};
      border-radius: ${borderRadius.base};
    `,
    footer: css`
      display: flex;
      justify-content: flex-end;
      margin-top: ${spacing[5]};
      padding-top: ${spacing[4]};
      border-top: 1px solid ${colors.gray[200]};
    `,
    adjustButton: css`
      padding: ${spacing[2]} ${spacing[5]};
      border: none;
      background-color: ${colors.primary[500]};
      color: white;
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      cursor: pointer;
      outline: none;
      transition: background-color 0.2s;

      &:hover {
        background-color: ${colors.primary[600]};
      }

      &:active {
        background-color: ${colors.primary[700]};
      }
    `,
  };
};
