// frontend/packages/arch/ui-components/src/components/Alert/Alert.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, typography } from '@coze-studio/common/themes';
import { AlertType } from './Alert';

interface AlertStylesProps {
  type: AlertType;
}

const colorMap = {
  success: {
    bg: '#f6ffed',
    border: '#b7eb8f',
    text: '#389e0d',
    iconBg: '#d9f7be',
  },
  info: {
    bg: '#e6f7ff',
    border: '#91d5ff',
    text: '#0050b3',
    iconBg: '#bae7ff',
  },
  warning: {
    bg: '#fffbe6',
    border: '#ffe58f',
    text: '#d48806',
    iconBg: '#fff1b8',
  },
  error: {
    bg: '#fff2f0',
    border: '#ffccc7',
    text: '#cf1322',
    iconBg: '#ffccc7',
  },
};

export const useStyles = ({ type }: AlertStylesProps) => {
  const themeColors = colorMap[type];

  return {
    container: css`
      display: flex;
      align-items: flex-start;
      gap: ${spacing[2]};
      padding: ${spacing[3]};
      background-color: ${themeColors.bg};
      border: 1px solid ${themeColors.border};
      border-radius: ${borderRadius.md};
    `,

    icon: css`
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 20px;
      height: 20px;
      color: ${themeColors.text};
      font-size: ${typography.fontSize.base};
    `,

    content: css`
      flex: 1;
    `,

    message: css`
      color: ${themeColors.text};
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
    `,

    description: css`
      margin-top: ${spacing[1]};
      color: ${colors.gray[600]};
      font-size: ${typography.fontSize.sm};
      line-height: 1.5;
    `,

    closeButton: css`
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 20px;
      height: 20px;
      padding: 0;
      border: none;
      background: transparent;
      color: ${colors.gray[400]};
      cursor: pointer;
      font-size: ${typography.fontSize.sm};

      &:hover {
        color: ${colors.gray[800]};
      }
    `,
  };
};
