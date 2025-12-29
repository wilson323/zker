// frontend/packages/arch/ui-components/src/components/Form/FormItem.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
} from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    formItem: css`
      margin-bottom: ${spacing[4]};
    `,
    label: css`
      display: block;
      margin-bottom: ${spacing[2]};
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[900]};
    `,
    required: css`
      color: ${colors.error.main};
      margin-left: ${spacing[1]};
    `,
    content: css`
      position: relative;
    `,
    error: css`
      margin-top: ${spacing[1]};
      font-size: ${typography.fontSize.sm};
      color: ${colors.error.main};
    `,
  };
};
