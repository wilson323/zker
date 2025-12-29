// frontend/packages/arch/ui-components/src/components/TextArea/TextArea.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  borderRadius,
  transitions,
} from '@coze-studio/common/themes';

interface TextAreaStyleProps {
  size: 'sm' | 'md' | 'lg';
  error: boolean;
  block: boolean;
}

export const useStyles = (props: TextAreaStyleProps) => {
  const { size, error, block } = props;

  const sizeStyles = {
    sm: css`
      padding: ${spacing[2]};
      ${typography.fontSize.sm};
    `,
    md: css`
      padding: ${spacing[3]};
      ${typography.fontSize.base};
    `,
    lg: css`
      padding: ${spacing[4]};
      ${typography.fontSize.lg};
    `,
  };

  return {
    textarea: css`
      display: ${block ? 'block' : 'inline-block'};
      width: ${block ? '100%' : 'auto'};
      font-family: ${typography.fontFamily.base};
      border: 1px solid ${error ? colors.error.main : colors.gray[300]};
      background-color: white;
      color: ${colors.gray[900]};
      border-radius: ${borderRadius.base};
      outline: none;
      resize: vertical;
      transition: border-color ${transitions.base};

      &::placeholder {
        color: ${colors.gray[400]};
      }

      &:focus {
        border-color: ${error ? colors.error.main : colors.primary[500]};
        ${!error && `
          outline: 2px solid ${colors.primary[200]};
          outline-offset: -1px;
        `}
      }

      &:disabled {
        background-color: ${colors.gray[50]};
        cursor: not-allowed;
      }

      ${sizeStyles[size]}
    `,
  };
};
