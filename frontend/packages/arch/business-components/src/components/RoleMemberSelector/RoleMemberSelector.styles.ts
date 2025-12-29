// frontend/packages/arch/business-components/src/components/RoleMemberSelector/RoleMemberSelector.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, boxShadow, zIndex, typography } from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      position: relative;
    `,

    searchBox: css`
      margin-bottom: ${spacing[2]};
    `,

    dropdown: css`
      position: absolute;
      top: calc(100% + 4px);
      left: 0;
      right: 0;
      max-height: 300px;
      overflow-y: auto;
      background-color: #fff;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.md};
      box-shadow: ${boxShadow.lg};
      z-index: ${zIndex.dropdown};
    `,

    empty: css`
      padding: ${spacing[3]};
      text-align: center;
      color: ${colors.gray[600]};
      font-size: ${typography.fontSize.sm};
    `,

    userItem: css`
      display: flex;
      align-items: center;
      gap: ${spacing[2]};
      padding: ${spacing[2]} ${spacing[3]};
      cursor: pointer;
      transition: background-color 0.2s;

      &:hover {
        background-color: ${colors.gray[100]};
      }
    `,

    avatar: css`
      width: 32px;
      height: 32px;
      border-radius: 50%;
      object-fit: cover;
    `,

    userInfo: css`
      flex: 1;
      min-width: 0;
    `,

    username: css`
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[800]};
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    `,

    email: css`
      font-size: ${typography.fontSize.xs};
      color: ${colors.gray[600]};
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    `,

    memberList: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[2]};
    `,

    memberItem: css`
      display: flex;
      align-items: center;
      gap: ${spacing[2]};
      padding: ${spacing[2]};
      background-color: #fff;
      border: 1px solid ${colors.gray[200]};
      border-radius: ${borderRadius.base};
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
