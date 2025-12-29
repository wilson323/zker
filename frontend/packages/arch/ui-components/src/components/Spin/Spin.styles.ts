// frontend/packages/arch/ui-components/src/components/Spin/Spin.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, zIndex, typography } from '@coze-studio/common/themes';

interface SpinStylesProps {
  size: 'small' | 'default' | 'large';
}

const sizeMap = {
  small: {
    containerSize: '14px',
    dotSize: '3px',
  },
  default: {
    containerSize: '18px',
    dotSize: '4px',
  },
  large: {
    containerSize: '24px',
    dotSize: '6px',
  },
};

export const useStyles = ({ size }: SpinStylesProps) => {
  const sizes = sizeMap[size];

  return {
    container: css`
      display: inline-flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: ${spacing[2]};
    `,

    wrapper: css`
      position: relative;
    `,

    overlay: css`
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: ${spacing[2]};
      background-color: rgba(255, 255, 255, 0.8);
      z-index: ${zIndex.modal};
    `,

    spinner: css`
      display: inline-flex;
      align-items: center;
      gap: ${spacing[1]};
      height: ${sizes.containerSize};
    `,

    dot: css`
      display: inline-block;
      width: ${sizes.dotSize};
      height: ${sizes.dotSize};
      border-radius: 50%;
      background-color: ${colors.primary[500]};
      animation: spin-bounce 1.4s infinite ease-in-out both;

      &:nth-of-type(1) {
        animation-delay: -0.32s;
      }

      &:nth-of-type(2) {
        animation-delay: -0.16s;
      }

      @keyframes spin-bounce {
        0%,
        80%,
        100% {
          transform: scale(0);
        }
        40% {
          transform: scale(1);
        }
      }
    `,

    tip: css`
      color: ${colors.gray[600]};
      font-size: ${typography.fontSize.sm};
      text-align: center;
    `,
  };
};
