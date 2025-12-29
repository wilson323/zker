// frontend/packages/arch/ui-components/src/components/Checkbox/Checkbox.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  borderRadius,
  transitions,
} from '@coze-studio/common/themes';

interface CheckboxStyleProps {
  checked: boolean;
  disabled: boolean;
  indeterminate: boolean;
}

export const useStyles = (props: CheckboxStyleProps) => {
  const { checked, disabled, indeterminate } = props;

  return {
    container: css`
      display: inline-flex;
      align-items: center;
      cursor: ${disabled ? 'not-allowed' : 'pointer'};
      user-select: none;
    `,
    checkbox: css`
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
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100%;
      height: 100%;
      border: 1px solid ${checked || indeterminate ? colors.primary[500] : colors.gray[300]};
      background-color: ${checked || indeterminate ? colors.primary[500] : 'white'};
      border-radius: ${borderRadius.sm};
      transition: all ${transitions.base};

      ${!disabled && `
        &:hover {
          border-color: ${checked || indeterminate ? colors.primary[600] : colors.primary[500]};
        }
      `}

      ${disabled && `
        background-color: ${colors.gray[100]};
        border-color: ${colors.gray[300]};
        cursor: not-allowed;
      `}
    `,
    checkIcon: css`
      color: white;
      font-size: 12px;
      font-weight: ${typography.fontWeight.bold};
      line-height: 1;
    `,
    indeterminateIcon: css`
      color: white;
      font-size: 10px;
      font-weight: ${typography.fontWeight.bold};
      line-height: 1;
    `,
    children: css`
      font-size: ${typography.fontSize.base};
      color: ${disabled ? colors.gray[400] : colors.gray[900]};
    `,
  };
};
