// frontend/packages/arch/ui-components/src/components/Switch/Switch.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, transitions } from '@coze-studio/common/themes';

interface SwitchStyleProps {
  checked: boolean;
  disabled: boolean;
}

export const useStyles = (props: SwitchStyleProps) => {
  const { checked, disabled } = props;

  return {
    switch: css`
      position: relative;
      display: inline-flex;
      align-items: center;
      width: 44px;
      height: 22px;
      padding: 0;
      border: none;
      border-radius: 11px;
      background-color: ${checked ? colors.primary[500] : colors.gray[300]};
      cursor: ${disabled ? 'not-allowed' : 'pointer'};
      transition: background-color ${transitions.base};

      ${!disabled && `
        &:hover {
          background-color: ${checked ? colors.primary[600] : colors.gray[400]};
        }
      `}

      ${disabled && `
        opacity: 0.6;
        cursor: not-allowed;
      `}
    `,
    thumb: css`
      width: 18px;
      height: 18px;
      background-color: white;
      border-radius: 50%;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
      transform: ${checked ? 'translateX(22px)' : 'translateX(2px)'};
      transition: transform ${transitions.base};
    `,
  };
};
