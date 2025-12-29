// frontend/packages/arch/business-components/src/components/DataScopeSelector/DataScopeSelector.styles.ts

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
    header: css`
      margin-bottom: ${spacing[4]};
    `,
    title: css`
      margin: 0 0 ${spacing[1]} 0;
      font-size: ${typography.fontSize.lg};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
    subtitle: css`
      margin: 0;
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[600]};
    `,
    options: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[3]};
    `,
    option: css`
      display: flex;
      gap: ${spacing[3]};
      padding: ${spacing[4]};
      background-color: white;
      border: 1px solid ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        border-color: ${colors.primary[500]};
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      }
    `,
    optionSelected: css`
      border-color: ${colors.primary[500]};
      background-color: ${colors.primary[50]};
    `,
    optionContent: css`
      flex: 1;
    `,
    optionHeader: css`
      display: flex;
      align-items: center;
      gap: ${spacing[2]};
      margin-bottom: ${spacing[1]};
    `,
    icon: css`
      font-size: 20px;
    `,
    label: css`
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[900]};
    `,
    description: css`
      margin: 0;
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[600]};
      line-height: 1.5;
    `,
  };
};
