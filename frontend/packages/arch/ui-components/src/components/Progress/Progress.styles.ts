// frontend/packages/arch/ui-components/src/components/Progress/Progress.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, typography } from '@coze-studio/common/themes';
import { ProgressType } from './Progress';

interface ProgressStylesProps {
  type: ProgressType;
  status: 'normal' | 'active' | 'success' | 'exception';
  width: number;
  strokeWidth: number;
}

const statusColorMap = {
  normal: colors.primary[500],
  active: colors.primary[500],
  success: colors.success.main,
  exception: colors.error.main,
};

export const useStyles = ({ type, status, width, strokeWidth }: ProgressStylesProps) => {
  const statusColor = statusColorMap[status];

  return {
    container: css`
      display: flex;
      align-items: center;
      gap: ${spacing[2]};
    `,

    outer: css`
      flex: 1;
    `,

    inner: css`
      position: relative;
      display: inline-block;
      width: 100%;
      height: ${strokeWidth}px;
      background-color: ${colors.gray[200]};
      border-radius: ${borderRadius.base};
      overflow: hidden;
    `,

    bg: css`
      position: absolute;
      top: 0;
      left: 0;
      height: 100%;
      background-color: ${statusColor};
      border-radius: ${borderRadius.base};
      transition: width 0.3s ease-in-out;
    `,

    activeBg: css`
      position: absolute;
      top: 0;
      left: 0;
      height: 100%;
      background: linear-gradient(
        to right,
        transparent 0%,
        rgba(255, 255, 255, 0.3) 50%,
        transparent 100%
      );
      animation: progress-active 2.4s cubic-bezier(0.23, 1, 0.32, 1) infinite;
    `,

    info: css`
      display: inline-block;
      min-width: 2em;
      text-align: right;
      color: ${colors.gray[800]};
      font-size: ${typography.fontSize.sm};
      font-weight: ${typography.fontWeight.medium};
    `,

    circleContainer: css`
      position: relative;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: ${width}px;
      height: ${width}px;
    `,

    circleTrail: css`
      stroke: ${colors.gray[200]};
    `,

    circlePath: css`
      stroke: ${statusColor};
      transition: stroke-dashoffset 0.3s ease-in-out;
    `,

    '@keyframes progress-active': css`
      0% {
        transform: translateX(-100%);
        width: 100%;
      }
      100% {
        transform: translateX(100%);
        width: 100%;
      }
    `,
  };
};
