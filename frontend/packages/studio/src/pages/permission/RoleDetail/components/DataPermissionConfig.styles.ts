// frontend/packages/studio/src/pages/permission/RoleDetail/components/DataPermissionConfig.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  typography,
} from '@coze-studio/common/themes';

/**
 * DataPermissionConfig样式Hook
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
      margin: 0 0 ${spacing[2]} 0;
      font-size: ${typography.fontSize.xl};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
    description: css`
      margin: 0 0 ${spacing[5]} 0;
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[500]};
      line-height: ${typography.lineHeight.normal};
    `,
    permissionList: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[5]};
    `,
    permissionItem: css`
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.base};
      padding: ${spacing[4]};
      background-color: ${colors.gray[50]};
    `,
    permissionHeader: css`
      margin-bottom: ${spacing[3]};
    `,
    resourceType: css`
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
    scopeOptions: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[3]};
    `,
    scopeOption: css`
      display: flex;
      align-items: flex-start;
      gap: ${spacing[3]};
      padding: ${spacing[3]};
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.base};
      background-color: white;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        border-color: ${colors.primary[300]};
      }

      &.selected {
        border-color: ${colors.primary[500]};
        background-color: ${colors.primary[50]};
      }

      input[type="radio"] {
        margin-top: 2px;
        accent-color: ${colors.primary[500]};
        cursor: pointer;
      }
    `,
    optionContent: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[1]};
      flex: 1;
    `,
    optionLabel: css`
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[900];
    `,
    optionDescription: css`
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[500]};
    `,
  };
};
