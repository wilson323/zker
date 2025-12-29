// frontend/packages/common/components/LanguageSwitcher/LanguageSwitcher.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  typography,
} from '@coze-studio/common/themes';

/**
 * LanguageSwitcher样式Hook
 */
export const useStyles = () => {
  return {
    container: css`
      display: flex;
      align-items: center;
      gap: ${spacing[3]};
    `,
    label: css`
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[700]};
    `,
    select: css`
      height: 36px;
      padding: ${spacing[2]} ${spacing[3]};
      border: 1px solid ${colors.gray[300]};
      background-color: white;
      color: ${colors.gray[700]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.sm};
      cursor: pointer;
      outline: none;
      transition: border-color 0.2s;

      &:focus {
        border-color: ${colors.primary[500]};
      }

      &:hover {
        border-color: ${colors.primary[400]};
      }
    `,
  };
};
