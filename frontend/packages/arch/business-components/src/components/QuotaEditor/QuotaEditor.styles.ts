// frontend/packages/arch/business-components/src/components/QuotaEditor/QuotaEditor.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, typography } from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[4]};
    `,

    quotaItem: css`
      padding: ${spacing[3]};
      background-color: #fff;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.md};
    `,

    header: css`
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: ${spacing[2]};
    `,

    label: css`
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[800]};
    `,

    unit: css`
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[600]};
    `,

    control: css`
      display: flex;
      align-items: center;
      gap: ${spacing[3]};
      margin-bottom: ${spacing[2]};
    `,

    input: css`
      width: 120px;
    `,

    options: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[2]};
      margin-bottom: ${spacing[2]};
    `,

    checkboxLabel: css`
      display: flex;
      align-items: center;
      gap: ${spacing[1]};
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[800]};
      cursor: pointer;
    `,

    overageFee: css`
      display: flex;
      align-items: center;
      gap: ${spacing[1]};
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[600]};
    `,

    feeInput: css`
      width: 100px;
    `,

    price: css`
      font-size: ${typography.fontSize.sm};
      color: ${colors.primary[500]};
    `,

    total: css`
      padding: ${spacing[3]};
      text-align: right;
      font-size: ${typography.fontSize.xl};
      color: ${colors.primary[500]};
      border-top: 2px solid ${colors.gray[200]};
    `,
  };
};
