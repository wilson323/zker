// frontend/packages/arch/business-components/src/components/SubscriptionSelector/SubscriptionSelector.styles.ts

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
      width: 100%;
    `,
    label: css`
      display: block;
      margin-bottom: ${spacing[2]};
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[900]};
    `,
    details: css`
      margin-top: ${spacing[4]};
      padding: ${spacing[4]};
      background-color: ${colors.gray[50]};
      border-radius: ${borderRadius.base};
      border: 1px solid ${colors.gray[200]};
    `,
    header: css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: ${spacing[2]};
    `,
    title: css`
      margin: 0;
      font-size: ${typography.fontSize.lg};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
    price: css`
      font-size: ${typography.fontSize.xl};
      font-weight: ${typography.fontWeight.bold};
      color: ${colors.primary[500]};
    `,
    description: css`
      margin: 0 0 ${spacing[3]} 0;
      font-size: ${typography.fontSize.base};
      color: ${colors.gray[600]};
    `,
    features: css`
      margin: 0;
      padding-left: ${spacing[4]};
      list-style: disc;

      li {
        margin-bottom: ${spacing[2]};
        font-size: ${typography.fontSize.sm};
        color: ${colors.gray[700]};
      }
    `,
  };
};
