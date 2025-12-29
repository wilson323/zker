// frontend/packages/studio/src/pages/tenant/TenantDetail/components/TenantBasicInfo.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  typography,
} from '@coze-studio/common/themes';

/**
 * TenantBasicInfo样式Hook
 */
export const useStyles = () => {
  return {
    container: css`
      padding: ${spacing[5]};
      background-color: white;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.base};
    `,
    title: css`
      margin: 0 0 ${spacing[5]} 0;
      font-size: ${typography.fontSize.xl};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
    form: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[4]};
    `,
    formRow: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[2]};
    `,
    label: css`
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[700]};
    `,
    readonlyInput: css`
      height: 40px;
      padding: ${spacing[2]} ${spacing[3]};
      border: 1px solid ${colors.gray[200]};
      background-color: ${colors.gray[50]};
      color: ${colors.gray[400]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      cursor: not-allowed;
    `,
    select: css`
      height: 40px;
      padding: ${spacing[2]} ${spacing[3]};
      border: 1px solid ${colors.gray[300]};
      background-color: white;
      color: ${colors.gray[700]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      cursor: pointer;
      outline: none;
      transition: border-color 0.2s;

      &:focus {
        border-color: ${colors.primary[500]};
      }
    `,
    textarea: css`
      padding: ${spacing[2]} ${spacing[3]};
      border: 1px solid ${colors.gray[300]};
      background-color: white;
      color: ${colors.gray[700]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      font-family: inherit;
      resize: vertical;
      outline: none;
      transition: border-color 0.2s;

      &:focus {
        border-color: ${colors.primary[500]};
      }
    `,
    footer: css`
      display: flex;
      justify-content: flex-end;
      margin-top: ${spacing[4]};
    `,
    submitButton: css`
      padding: ${spacing[2]} ${spacing[5]};
      border: none;
      background-color: ${colors.primary[500]};
      color: white;
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      cursor: pointer;
      outline: none;
      transition: background-color 0.2s;

      &:hover {
        background-color: ${colors.primary[600]};
      }

      &:active {
        background-color: ${colors.primary[700]};
      }
    `,
  };
};
