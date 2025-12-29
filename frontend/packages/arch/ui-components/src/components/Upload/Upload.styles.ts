// frontend/packages/arch/ui-components/src/components/Upload/Upload.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, typography } from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      display: inline-block;
    `,

    uploadButton: css`
      padding: ${spacing[2]} ${spacing[3]};
      background-color: ${colors.primary[500]};
      color: #fff;
      border: none;
      border-radius: ${borderRadius.base};
      cursor: pointer;
      font-size: ${typography.fontSize.base};
      transition: background-color 0.2s;

      &:hover:not(:disabled) {
        background-color: ${colors.primary[600]};
      }

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }
    `,

    dragger: css`
      position: relative;
      width: 100%;
      padding: ${spacing[5]};
      background-color: #fff;
      border: 2px dashed ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        border-color: ${colors.primary[500]};
      }
    `,

    dragging: css`
      border-color: ${colors.primary[500]};
      background-color: ${colors.primary[50]};
    `,

    dragContent: css`
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: ${spacing[2]};
    `,

    dragIcon: css`
      font-size: 48px;
    `,

    dragText: css`
      color: ${colors.gray[600]};
      font-size: ${typography.fontSize.base};
    `,

    list: css`
      margin-top: ${spacing[3]};
    `,

    listItem: css`
      position: relative;
      display: flex;
      align-items: center;
      gap: ${spacing[2]};
      padding: ${spacing[2]};
      background-color: #fff;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.base};
      margin-bottom: ${spacing[2]};
    `,

    fileInfo: css`
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: ${spacing[2]};
    `,

    fileName: css`
      color: ${colors.gray[800]};
      font-size: ${typography.fontSize.sm};
    `,

    fileStatus: css`
      color: ${colors.gray[600]};
      font-size: ${typography.fontSize.sm};
    `,

    progressBar: css`
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 2px;
      background-color: ${colors.gray[200]};
    `,

    progressFill: css`
      height: 100%;
      background-color: ${colors.primary[500]};
      transition: width 0.3s;
    `,

    removeButton: css`
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 20px;
      height: 20px;
      padding: 0;
      border: none;
      background: transparent;
      color: ${colors.gray[400]};
      cursor: pointer;
      font-size: ${typography.fontSize.sm};

      &:hover {
        color: ${colors.error.main};
      }
    `,
  };
};
