// frontend/packages/arch/ui-components/src/components/Badge/Badge.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, typography } from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      position: relative;
      display: inline-block;
    `,
    badge: css`
      position: absolute;
      top: 0;
      right: 0;
      transform: translate(50%, -50%);
      background-color: ${colors.error.main};
      color: white;
      border-radius: 10px;
      padding: ${spacing[1]} ${spacing[2]};
      font-size: ${typography.fontSize.xs};
      font-weight: ${typography.fontWeight.bold};
      line-height: 1;
      white-space: nowrap;
      z-index: 1;
    `,
    dot: css`
      display: block;
      width: 6px;
      height: 6px;
      border-radius: 50%;
      background-color: currentColor;
    `,
  };
};
