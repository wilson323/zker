// frontend/packages/arch/ui-components/src/components/Tag/Tag.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  borderRadius,
} from '@coze-studio/common/themes';

interface TagStyleProps {
  color: 'default' | 'primary' | 'success' | 'warning' | 'error';
}

export const useStyles = (props: TagStyleProps) => {
  const { color } = props;

  const colorMap = {
    default: {
      bg: colors.gray[100],
      text: colors.gray[700],
    },
    primary: {
      bg: colors.primary[100],
      text: colors.primary[700],
    },
    success: {
      bg: '#F6FFED',
      text: '#52C41A',
    },
    warning: {
      bg: '#FFFBE6',
      text: '#FAAD14',
    },
    error: {
      bg: '#FFF1F0',
      text: '#F5222D',
    },
  };

  const colors = colorMap[color];

  return {
    tag: css`
      display: inline-flex;
      align-items: center;
      gap: ${spacing[2]};
      padding: ${spacing[1]} ${spacing[3]};
      background-color: ${colors.bg};
      color: ${colors.text};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
      line-height: 1.5;
    `,
    close: css`
      cursor: pointer;
      font-size: 14px;
      opacity: 0.6;

      &:hover {
        opacity: 1;
      }
    `,
  };
};
