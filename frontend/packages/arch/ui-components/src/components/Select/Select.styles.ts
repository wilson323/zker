// frontend/packages/arch/ui-components/src/components/Select/Select.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  borderRadius,
  transitions,
  boxShadow,
  zIndex,
} from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      position: relative;
      width: 100%;
    `,
    trigger: css`
      display: flex;
      align-items: center;
      justify-content: space-between;
      min-height: 36px;
      padding: ${spacing[2]} ${spacing[3]};
      background-color: white;
      border: 1px solid ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      cursor: pointer;
      transition: all ${transitions.base};

      &:hover {
        border-color: ${colors.primary[500]};
      }

      &:focus-within {
        border-color: ${colors.primary[500]};
        outline: 2px solid ${colors.primary[200]};
      }
    `,
    value: css`
      flex: 1;
      color: ${colors.gray[900]};
      font-size: ${typography.fontSize.base};
    `,
    placeholder: css`
      flex: 1;
      color: ${colors.gray[500]};
      font-size: ${typography.fontSize.base};
    `,
    suffix: css`
      display: flex;
      align-items: center;
      gap: ${spacing[2]};
      margin-left: ${spacing[2]};
    `,
    spinner: css`
      width: 14px;
      height: 14px;
      border: 2px solid ${colors.gray[300]};
      border-top-color: ${colors.primary[500]};
      border-radius: 50%;
      animation: spin 0.6s linear infinite;

      @keyframes spin {
        to { transform: rotate(360deg); }
      }
    `,
    clearIcon: css`
      font-size: 16px;
      color: ${colors.gray[500]};
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 16px;
      height: 16px;
      border-radius: 50%;

      &:hover {
        color: ${colors.gray[700]};
        background-color: ${colors.gray[200]};
      }
    `,
    arrow: css`
      font-size: 10px;
      color: ${colors.gray[500]};
      transition: transform ${transitions.base};
    `,
    arrowOpen: css`
      transform: rotate(180deg);
    `,
    dropdown: css`
      position: absolute;
      top: calc(100% + 4px);
      left: 0;
      right: 0;
      max-height: 256px;
      overflow-y: auto;
      background-color: white;
      border: 1px solid ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      box-shadow: ${boxShadow.lg};
      z-index: ${zIndex.dropdown};
    `,
    option: css`
      padding: ${spacing[2]} ${spacing[3]};
      font-size: ${typography.fontSize.base};
      color: ${colors.gray[900]};
      cursor: pointer;
      transition: background-color ${transitions.fast};

      &:hover {
        background-color: ${colors.primary[50]};
      }
    `,
    optionSelected: css`
      background-color: ${colors.primary[100]};
      color: ${colors.primary[700]};
      font-weight: ${typography.fontWeight.medium};

      &:hover {
        background-color: ${colors.primary[100]};
      }
    `,
    optionDisabled: css`
      color: ${colors.gray[400]};
      cursor: not-allowed;

      &:hover {
        background-color: transparent;
      }
    `,
  };
};
