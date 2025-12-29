// frontend/packages/studio/src/pages/routing/IntentMatcher/IntentMatcher.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  typography,
  borderRadius,
} from '@coze-studio/common/themes';

export const useStyles = () => {
  return {
    container: css`
      padding: ${spacing[6]};
      background-color: ${colors.gray[50]};
      min-height: 100vh;
    `,
    header: css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: ${spacing[6]};
    `,
    title: css`
      font-size: ${typography.fontSize['2xl']};
      font-weight: ${typography.fontWeight.bold};
      color: ${colors.gray[900]};
      margin: 0;
    `,
    tabs: css`
      display: flex;
      gap: ${spacing[2]};
      margin-bottom: ${spacing[6]};
      border-bottom: 1px solid ${colors.gray[300]};
    `,
    tab: css`
      padding: ${spacing[3]} ${spacing[4]};
      background: none;
      border: none;
      border-bottom: 2px solid transparent;
      cursor: pointer;
      font-size: ${typography.fontSize.base};
      color: ${colors.gray[700]};
      transition: all 0.2s;

      &:hover {
        color: ${colors.primary[500]};
      }
    `,
    activeTab: css`
      padding: ${spacing[3]} ${spacing[4]};
      background: none;
      border: none;
      border-bottom: 2px solid ${colors.primary[500]};
      cursor: pointer;
      font-size: ${typography.fontSize.base};
      color: ${colors.primary[500]};
      font-weight: ${typography.fontWeight.medium};
    `,
    filter: css`
      display: flex;
      gap: ${spacing[3]};
      margin-bottom: ${spacing[4]};
      padding: ${spacing[4]};
      background-color: white;
      border-radius: ${borderRadius.base};
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
    `,
    searchInput: css`
      flex: 1;
      padding: ${spacing[2]} ${spacing[3]};
      border: 1px solid ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      outline: none;

      &:focus {
        border-color: ${colors.primary[500]};
      }
    `,
    typeSelect: css`
      padding: ${spacing[2]} ${spacing[3]};
      border: 1px solid ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      background-color: white;
      cursor: pointer;
      outline: none;

      &:focus {
        border-color: ${colors.primary[500]};
      }
    `,
    table: css`
      background-color: white;
      border-radius: ${borderRadius.base};
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
      overflow: hidden;
    `,
    testTool: css`
      padding: ${spacing[6]};
      background-color: white;
      border-radius: ${borderRadius.base};
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
    `,
    testTitle: css`
      font-size: ${typography.fontSize.xl};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[900]};
      margin-bottom: ${spacing[6]};
    `,
    testInput: css`
      display: flex;
      flex-direction: column;
      gap: ${spacing[3]};
    `,
    testTextarea: css`
      width: 100%;
      padding: ${spacing[3]};
      border: 1px solid ${colors.gray[300]};
      border-radius: ${borderRadius.base};
      font-size: ${typography.fontSize.base};
      font-family: ${typography.fontFamily.base};
      resize: vertical;
      outline: none;

      &:focus {
        border-color: ${colors.primary[500]};
      }
    `,
    testResult: css`
      margin-top: ${spacing[6]};
    `,
    resultEmpty: css`
      padding: ${spacing[6]};
      text-align: center;
      color: ${colors.gray[500]};
      background-color: ${colors.gray[50]};
      border-radius: ${borderRadius.base};
      border: 1px dashed ${colors.gray[300]};
    `,
  };
};
