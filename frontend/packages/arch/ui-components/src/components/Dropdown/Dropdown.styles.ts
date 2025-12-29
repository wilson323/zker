// frontend/packages/arch/ui-components/src/components/Dropdown/Dropdown.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, boxShadow, zIndex, typography } from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      position: relative;
      display: inline-block;
    `,

    trigger: css`
      cursor: pointer;
    `,

    menu: css`
      position: absolute;
      top: calc(100% + 4px);
      left: 0;
      min-width: 120px;
      padding: ${spacing[1]} 0;
      background-color: #fff;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.md};
      box-shadow: ${boxShadow.lg};
      z-index: ${zIndex.dropdown};
    `,

    menuItem: css`
      position: relative;
      padding: ${spacing[2]} ${spacing[3]};
      color: ${colors.gray[800]};
      font-size: ${typography.fontSize.base};
      cursor: pointer;
      transition: background-color 0.2s;

      &:hover {
        background-color: ${colors.gray[100]};
      }

      &.disabled {
        color: ${colors.gray[400]};
        cursor: not-allowed;
        pointer-events: none;
      }

      &.danger {
        color: ${colors.error.main};

        &:hover {
          background-color: #fff2f0;
        }
      }

      &.divided {
        border-top: 1px solid ${colors.gray[200]};
      }
    `,

    label: css`
      display: block;
    `,

    submenu: css`
      position: absolute;
      left: 100%;
      top: 0;
      margin-left: 4px;
    `,

    disabled: css`
      opacity: 0.5;
    `,

    danger: css`
      color: ${colors.error.main};
    `,

    divided: css`
      border-top: 1px solid ${colors.gray[200]};
    `,
  };
};
