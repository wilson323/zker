// frontend/packages/arch/business-components/src/components/TenantStats/TenantStats.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, boxShadow, typography } from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
      gap: ${spacing[3]};
    `,

    statCard: css`
      padding: ${spacing[3]};
      background-color: #fff;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.md};
      transition: box-shadow 0.2s;

      &:hover {
        box-shadow: ${boxShadow.md};
      }
    `,

    header: css`
      display: flex;
      align-items: center;
      gap: ${spacing[2]};
      margin-bottom: ${spacing[2]};
    `,

    icon: css`
      font-size: 24px;
    `,

    label: css`
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[600]};
    `,

    value: css`
      font-size: ${typography.fontSize['2xl']};
      font-weight: ${typography.fontWeight.bold};
      color: ${colors.gray[800]};
      margin-bottom: ${spacing[1]};
    `,

    active: css`
      display: flex;
      align-items: center;
      gap: ${spacing[1]};
      margin-bottom: ${spacing[2]};
      font-size: ${typography.fontSize.sm};
    `,

    activeLabel: css`
      color: ${colors.gray[600]};
    `,

    activeValue: css`
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.success.main};
    `,

    percentage: css`
      color: ${colors.gray[400]};
    `,

    bar: css`
      height: 4px;
      background-color: ${colors.gray[200]};
      border-radius: 2px;
      overflow: hidden;
    `,

    barFill: css`
      height: 100%;
      transition: width 0.3s ease-in-out;
    `,
  };
};
