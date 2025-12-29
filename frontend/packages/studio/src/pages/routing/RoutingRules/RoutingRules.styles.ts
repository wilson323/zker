// frontend/packages/studio/src/pages/routing/RoutingRules/RoutingRules.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  borderRadius,
} from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      padding: ${spacing[6]};
      background-color: ${colors.gray[50]};
      min-height: 100vh;
    `,
    header: css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: ${spacing[6]};
    `,
    title: css`
      font-size: ${typography.fontSize['2xl']};
      font-weight: ${typography.fontWeight.bold};
      color: ${colors.gray[900]};
      margin: 0;
    `,
    filter: css`
      display: flex;
      gap: ${spacing[3]};
      margin-bottom: ${spacing[4]};
      padding: ${spacing[4]};
      background-color: white;
      border-radius: ${borderRadius.base};
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
    `,
    searchInput: css`
      flex: 1;
      padding: ${spacing[2]} ${spacing[3]};
      border: 1px solid ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      outline: none;

      &:focus {
        border-color: ${colors.primary[500]};
      }
    `,
    table: css`
      background-color: white;
      border-radius: ${borderRadius.base};
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
      overflow: hidden;
    `,
  };
};
