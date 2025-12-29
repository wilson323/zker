// frontend/packages/arch/ui-components/src/components/Card/Card.styles.ts

import { css } from '@emotion/react';
import { colors, spacing, borderRadius, transitions, typography, boxShadow } from '@coze-studio/common/themes';

interface CardStylesProps {
  bordered: boolean;
  hoverable: boolean;
}

export const useStyles = ({ bordered, hoverable }: CardStylesProps) => {
  return {
    card: css`
      background-color: #fff;
      border-radius: ${borderRadius.lg};
      transition: all ${transitions.base};

      ${bordered && `border: 1px solid ${colors.gray[200]};`}

      ${hoverable && `
        cursor: pointer;

        &:hover {
          box-shadow: ${boxShadow.md};
        }
      `}
    `,

    header: css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: ${spacing[4]};
      border-bottom: 1px solid ${colors.gray[200]};

      &:not(:last-child) {
        padding-bottom: ${spacing[3]};
      }
    `,

    title: css`
      font-size: ${typography.fontSize.xl};
      font-weight: ${typography.fontWeight.semibold};
      color: ${colors.gray[800]};
    `,

    extra: css`
      font-size: ${typography.fontSize.base};
    `,

    body: css`
      padding: ${spacing[4]};

      &:not(:last-child) {
        padding-top: ${spacing[3]};
      }
    `,

    hoverable: css`
      &:hover {
        border-color: ${colors.primary[400]};
      }
    `,
  };
};
