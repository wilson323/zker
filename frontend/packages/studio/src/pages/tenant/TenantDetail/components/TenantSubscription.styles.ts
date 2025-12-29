// frontend/packages/studio/src/pages/tenant/TenantDetail/components/TenantSubscription.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  typography,
} from '@coze-studio/common/themes';

/**
 * TenantSubscription样式Hook
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
    content: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[4]};
    `,
    section: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[2]};
    `,
    label: css`
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[700]};
    `,
    value: css`
      font-size: ${typography.fontSize.base};
      color: ${colors.gray[900]};
    `,
    tierBadge: css`
      display: inline-block;
      padding: ${spacing[1]} ${spacing[3]};
      color: white;
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
      width: fit-content;
    `,
    switch: css`
      position: relative;
      display: inline-block;
      width: 50px;
      height: 26px;
      cursor: pointer;

      & input {
        opacity: 0;
        width: 0;
        height: 0;

        &:checked + .slider {
          background-color: ${colors.primary[500]};
        }

        &:checked + .slider:before {
          transform: translateX(24px);
        }
      }
    `,
    slider: css`
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-color: ${colors.gray[300]};
      transition: 0.3s;
      border-radius: 26px;

      &:before {
        content: '';
        position: absolute;
        height: 20px;
        width: 20px;
        left: 3px;
        bottom: 3px;
        background-color: white;
        transition: 0.3s;
        border-radius: 50%;
      }
    `,
    quotaSection: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[4]};
      padding-top: ${spacing[4]};
      border-top: 1px solid ${colors.gray[200]};
    `,
    quotaTitle: css`
      margin: 0;
      font-size: ${typography.fontSize.lg};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900];
    `,
    footer: css`
      display: flex;
      justify-content: flex-end;
      padding-top: ${spacing[4]};
      border-top: 1px solid ${colors.gray[200]};
    `,
    upgradeButton: css`
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
