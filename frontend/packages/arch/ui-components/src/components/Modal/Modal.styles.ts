// frontend/packages/arch/ui-components/src/components/Modal/Modal.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  transitions,
  zIndex,
  typography,
} from '@coze-studio/common/themes';

/**
 * Modal样式Props接口
 */
interface ModalStyleProps {
  width: number | string;
}

/**
 * Modal样式Hook
 */
export const useStyles = (props: ModalStyleProps) => {
  const { width } = props;

  // ==================== 返回样式对象 ====================

  return {
    mask: css`
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-color: rgba(0, 0, 0, 0.45);
      z-index: ${zIndex.modalBackdrop};
      display: flex;
      align-items: center;
      justify-content: center;
      padding: ${spacing[4]};
      animation: fadeIn 0.2s ease-out;

      @keyframes fadeIn {
        from {
          opacity: 0;
        }
        to {
          opacity: 1;
        }
      }
    `,
    modal: css`
      position: relative;
      width: ${typeof width === 'number' ? `${width}px` : width};
      max-width: calc(100vw - 32px);
      max-height: calc(100vh - 32px);
      background-color: white;
      border-radius: ${borderRadius.lg};
      box-shadow: ${colors.gray[900]} 0px 4px 12px;
      display: flex;
      flex-direction: column;
      animation: slideIn 0.2s ease-out;

      @keyframes slideIn {
        from {
          transform: translateY(-20px);
          opacity: 0;
        }
        to {
          transform: translateY(0);
          opacity: 1;
        }
      }
    `,
    header: css`
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: ${spacing[4]} ${spacing[5]};
      border-bottom: 1px solid ${colors.gray[200]};
    `,
    title: css`
      margin: 0;
      font-size: ${typography.fontSize.lg};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
      line-height: ${typography.lineHeight.normal};
    `,
    closeButton: css`
      background: none;
      border: none;
      font-size: 24px;
      color: ${colors.gray[400]};
      cursor: pointer;
      padding: 0;
      width: 32px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: ${borderRadius.base};
      transition: all ${transitions.fast};
      line-height: 1;

      &:hover {
        background-color: ${colors.gray[100]};
        color: ${colors.gray[600]};
      }

      &:active {
        background-color: ${colors.gray[200]};
      }
    `,
    body: css`
      flex: 1;
      padding: ${spacing[5]};
      overflow-y: auto;
    `,
    footer: css`
      display: flex;
      justify-content: flex-end;
      gap: ${spacing[3]};
      padding: ${spacing[4]} ${spacing[5]};
      border-top: 1px solid ${colors.gray[200]};
    `,
    footerButton: css`
      padding: ${spacing[2]} ${spacing[4]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      cursor: pointer;
      transition: all ${transitions.fast};
      border: 1px solid transparent;
      min-width: 70px;

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }
    `,
    okButton: css`
      background-color: ${colors.primary[500]};
      color: white;
      border-color: ${colors.primary[500]};

      &:hover:not(:disabled) {
        background-color: ${colors.primary[600]};
        border-color: ${colors.primary[600]};
      }

      &:active:not(:disabled) {
        background-color: ${colors.primary[700]};
        border-color: ${colors.primary[700]};
      }
    `,
    cancelButton: css`
      background-color: white;
      color: ${colors.gray[700]};
      border-color: ${colors.gray[300]};

      &:hover:not(:disabled) {
        color: ${colors.primary[500]};
        border-color: ${colors.primary[500]};
      }

      &:active:not(:disabled) {
        background-color: ${colors.gray[50]};
      }
    `,
    spinner: css`
      display: inline-block;
      width: 14px;
      height: 14px;
      border: 2px solid currentColor;
      border-top-color: transparent;
      border-radius: 50%;
      animation: spin 0.6s linear infinite;
      margin-right: ${spacing[2]};

      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }
    `,
  };
};
