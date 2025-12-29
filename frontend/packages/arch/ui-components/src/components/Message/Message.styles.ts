// frontend/packages/arch/ui-components/src/components/Message/Message.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, boxShadow, typography } from '@coze-studio/common/themes';
import { MessageType } from './Message';

interface MessageStylesProps {
  type: MessageType;
}

const colorMap = {
  success: {
    bg: colors.success.main,
    color: '#fff',
  },
  info: {
    bg: colors.info.main,
    color: '#fff',
  },
  warning: {
    bg: colors.warning.main,
    color: '#fff',
  },
  error: {
    bg: colors.error.main,
    color: '#fff',
  },
};

export const useStyles = ({ type }: MessageStylesProps) => {
  const themeColors = colorMap[type];

  return {
    container: css`
      display: inline-flex;
      align-items: center;
      gap: ${spacing[2]};
      padding: ${spacing[2]} ${spacing[3]};
      background-color: ${themeColors.bg};
      color: ${themeColors.color};
      border-radius: ${borderRadius.md};
      box-shadow: ${boxShadow.md};
      font-size: ${typography.fontSize.base};
      animation: message-in 0.3s ease-in-out;
    `,

    icon: css`
      font-size: ${typography.fontSize.lg};
    `,

    content: css`
      line-height: 1.5;
    `,

    '@keyframes message-in': css`
      from {
        opacity: 0;
        transform: translateY(-20px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    `,
  };
};
