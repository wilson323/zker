// frontend/packages/arch/ui-components/src/components/DatePicker/DatePicker.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, boxShadow, zIndex, typography } from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      position: relative;
      display: inline-block;
    `,

    panel: css`
      position: absolute;
      top: calc(100% + 4px);
      left: 0;
      z-index: ${zIndex.dropdown};
      padding: ${spacing[3]};
      background-color: #fff;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.md};
      box-shadow: ${boxShadow.lg};
      min-width: 280px;
    `,

    header: css`
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: ${spacing[2]};
    `,

    title: css`
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[800]};
    `,

    navButton: css`
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 28px;
      height: 28px;
      padding: 0;
      border: none;
      background: transparent;
      color: ${colors.gray[600]};
      cursor: pointer;
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.xl};

      &:hover {
        background-color: ${colors.gray[100]};
      }
    `,

    weekHeader: css`
      display: grid;
      grid-template-columns: repeat(7, 1fr);
      gap: ${spacing[1]};
      margin-bottom: ${spacing[2]};
    `,

    weekDay: css`
      text-align: center;
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[600]};
      font-weight: ${typography.fontWeight.medium};
    `,

    calendar: css`
      display: grid;
      grid-template-columns: repeat(7, 1fr);
      gap: ${spacing[1]};
    `,

    day: css`
      display: flex;
      align-items: center;
      justify-content: center;
      height: 32px;
      cursor: pointer;
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[800]};
      transition: all 0.2s;

      &:hover {
        background-color: ${colors.primary[50]};
        color: ${colors.primary[500]};
      }
    `,

    emptyDay: css`
      height: 32px;
    `,

    selectedDay: css`
      background-color: ${colors.primary[500]} !important;
      color: #fff !important;
    `,

    today: css`
      border: 1px solid ${colors.primary[500]};
    `,
  };
};
