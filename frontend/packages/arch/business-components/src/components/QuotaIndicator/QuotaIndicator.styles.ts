// frontend/packages/arch/business-components/src/components/QuotaIndicator/QuotaIndicator.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  borderRadius,
  typography,
} from '@coze-studio/common/themes';

/**
 * QuotaIndicator样式Props接口
 */
interface QuotaIndicatorStyleProps {
  used: number;
  max: number;
}

/**
 * QuotaIndicator样式Hook
 */
export const useStyles = (props: QuotaIndicatorStyleProps) => {
  const { used, max } = props;

  const percentage = max > 0 ? (used / max) * 100 : 0;
  const isOverLimit = percentage >= 100;
  const isNearLimit = percentage >= 80 && percentage < 100;

  // ==================== 返回样式对象 ====================

  return {
    container: css`
      width: 100%;
    `,
    header: css`
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: ${spacing[2]};
    `,
    label: css`
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
      color: ${colors.gray[700]};
    `,
    status: css`
      font-size: ${typography.fontSize.xs};
      font-weight: ${typography.fontWeight.medium};
      padding: ${spacing[1]} ${spacing[2]};
      border-radius: ${borderRadius.sm};
    `,
    statusNormal: css`
      color: ${colors.success.main};
      background-color: ${colors.success.light};
    `,
    statusWarning: css`
      color: ${colors.warning.dark};
      background-color: ${colors.warning.light};
    `,
    statusError: css`
      color: ${colors.error.dark};
      background-color: ${colors.error.light};
    `,
    countRow: css`
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: ${spacing[2]};
    `,
    count: css`
      font-size: ${typography.fontSize.sm};
      color: ${colors.gray[600]};
    `,
    percentage: css`
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.semibold};
    `,
    barContainer: css`
      width: 100%;
      height: 8px;
      background-color: ${colors.gray[100]};
      border-radius: ${borderRadius.full};
      overflow: hidden;
    `,
    bar: css`
      height: 100%;
      background-color: ${colors.success.main};
      border-radius: ${borderRadius.full};
      transition: width 0.3s ease, background-color 0.3s ease;

      ${isNearLimit && css`
        background-color: ${colors.warning.main};
      `}

      ${isOverLimit && css`
        background-color: ${colors.error.main};
      `}
    `,
    nearLimit: css`
      background-color: ${colors.warning.main};
    `,
    overLimit: css`
      background-color: ${colors.error.main};
    `,
  };
};
