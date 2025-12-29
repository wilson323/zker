// frontend/packages/studio/src/pages/tenant/TenantList/TenantList.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
} from '@coze-studio/common/themes';

/**
 * TenantList样式Hook
 */
export const useStyles = () => {
  return {
    container: css`
      padding: ${spacing[5]};
      background-color: ${colors.gray[50]};
      min-height: 100vh;
    `,
    header: css`
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: ${spacing[5]};
      padding: ${spacing[4]};
      background-color: white;
      border-radius: 4px;
    `,
    title: css`
      margin: 0;
      font-size: ${typography.fontSize['2xl']};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
  };
};
