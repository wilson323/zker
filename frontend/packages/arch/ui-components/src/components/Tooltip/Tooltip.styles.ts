// frontend/packages/arch/ui-components/src/components/Tooltip/Tooltip.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, zIndex, typography } from '@coze-studio/common/themes';
import { TooltipPlacement } from './Tooltip';

interface TooltipStylesProps {
  placement: TooltipPlacement;
}

export const useStyles = ({ placement }: TooltipStylesProps) => {
  const placementStyles = {
    top: css`
      bottom: calc(100% + 8px);
      left: 50%;
      transform: translateX(-50%);

      .arrow {
        top: 100%;
        left: 50%;
        transform: translateX(-50%);
        border-top-color: ${colors.gray[800]};
        border-right-color: transparent;
        border-bottom-color: transparent;
        border-left-color: transparent;
      }
    `,
    bottom: css`
      top: calc(100% + 8px);
      left: 50%;
      transform: translateX(-50%);

      .arrow {
        bottom: 100%;
        left: 50%;
        transform: translateX(-50%);
        border-bottom-color: ${colors.gray[800]};
        border-right-color: transparent;
        border-top-color: transparent;
        border-left-color: transparent;
      }
    `,
    left: css`
      right: calc(100% + 8px);
      top: 50%;
      transform: translateY(-50%);

      .arrow {
        left: 100%;
        top: 50%;
        transform: translateY(-50%);
        border-left-color: ${colors.gray[800]};
        border-top-color: transparent;
        border-right-color: transparent;
        border-bottom-color: transparent;
      }
    `,
    right: css`
      left: calc(100% + 8px);
      top: 50%;
      transform: translateY(-50%);

      .arrow {
        right: 100%;
        top: 50%;
        transform: translateY(-50%);
        border-right-color: ${colors.gray[800]};
        border-top-color: transparent;
        border-left-color: transparent;
        border-bottom-color: transparent;
      }
    `,
  };

  return {
    wrapper: css`
      position: relative;
      display: inline-block;
    `,

    tooltip: css`
      position: absolute;
      z-index: ${zIndex.tooltip};
      padding: ${spacing[1]} ${spacing[2]};
      background-color: ${colors.gray[800]};
      color: #fff;
      font-size: ${typography.fontSize.sm};
      white-space: nowrap;
      border-radius: ${borderRadius.md};
      pointer-events: none;

      ${placementStyles[placement]}
    `,

    arrow: css`
      position: absolute;
      width: 0;
      height: 0;
      border: 5px solid transparent;
    `,
  };
};
