// frontend/packages/arch/ui-components/src/components/Radio/Radio.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  transitions,
} from '@coze-studio/common/themes';

interface RadioStyleProps {
  checked: boolean;
  disabled: boolean;
}

export const useStyles = (props: RadioStyleProps) => {
  const { checked, disabled } = props;

  return {
    container: css`
      display: inline-flex;
      align-items: center;
      cursor: ${disabled ? 'not-allowed' : 'pointer'};
      user-select: none;
    `,
    radio: css`
      position: relative;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 16px;
      height: 16px;
      margin-right: ${spacing[2]};
    `,
    input: css`
      position: absolute;
      width: 100%;
      height: 100%;
      opacity: 0;
      cursor: ${disabled ? 'not-allowed' : 'pointer'};
      margin: 0;
    `,
    inner: css`
      width: 100%;
      height: 100%;
      border: 1px solid ${checked ? colors.primary[500] : colors.gray[300]};
      background-color: white;
      border-radius: 50%;
      transition: all ${transitions.base};
      position: relative;

      ${!disabled && `
        &:hover {
          border-color: ${checked ? colors.primary[600] : colors.primary[500]};
        }
      `}

      ${disabled && `
        background-color: ${colors.gray[100]};
        border-color: ${colors.gray[300]};
        cursor: not-allowed;
      `}

      ${checked && `
        &::after {
          content: '';
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          width: 8px;
          height: 8px;
          background-color: ${colors.primary[500]};
          border-radius: 50%;
        }
      `}
    `,
    children: css`
      font-size: ${typography.fontSize.base};
      color: ${disabled ? colors.gray[400] : colors.gray[900]};
    `,
  };
};
