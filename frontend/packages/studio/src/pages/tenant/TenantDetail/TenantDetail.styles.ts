// frontend/packages/studio/src/pages/tenant/TenantDetail/TenantDetail.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  typography,
} from '@coze-studio/common/themes';

/**
 * TenantDetail样式Hook
 */
export const useStyles = () => {
  return {
    container: css`
      padding: ${spacing[5]};
      background-color: ${colors.gray[50]};
      min-height: 100vh;
    `,
    header: css`
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: ${spacing[5]};
      padding: ${spacing[4]};
      background-color: white;
      border-radius: ${borderRadius.base};
    `,
    headerLeft: css`
      display: flex;
      align-items: center;
      gap: ${spacing[4]};
    `,
    backButton: css`
      padding: ${spacing[2]} ${spacing[3]};
      border: 1px solid ${colors.gray[300]};
      background-color: white;
      color: ${colors.gray[700]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.sm};
      cursor: pointer;
      outline: none;
      transition: all 0.2s;

      &:hover {
        border-color: ${colors.primary[500]};
        color: ${colors.primary[500]};
      }
    `,
    title: css`
      margin: 0;
      font-size: ${typography.fontSize['2xl']};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
    `,
    headerRight: css`
      display: flex;
      align-items: center;
      gap: ${spacing[3]};
    `,
    status: css`
      padding: ${spacing[1]} ${spacing[3]};
      border-radius: ${borderRadius.sm};
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
    `,
    statusActive: css`
      background-color: ${colors.success.light};
      color: ${colors.success.main};
    `,
    tabs: css`
      display: flex;
      gap: ${spacing[2]};
      margin-bottom: ${spacing[4]};
    `,
    tab: css`
      padding: ${spacing[3]} ${spacing[5]};
      border: 1px solid ${colors.gray[200]};
      background-color: white;
      color: ${colors.gray[700]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      font-weight: ${typography.fontWeight.medium};
      cursor: pointer;
      outline: none;
      transition: all 0.2s;

      &:hover {
        border-color: ${colors.primary[500]};
        color: ${colors.primary[500]};
      }
    `,
    activeTab: css`
      background-color: ${colors.primary[500]};
      color: white;
      border-color: ${colors.primary[500]};

      &:hover {
        background-color: ${colors.primary[600]};
        border-color: ${colors.primary[600]};
        color: white;
      }
    `,
    tabContent: css`
      /* 标签页内容区域 */
    `,
  };
};
