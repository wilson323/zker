// frontend/packages/arch/ui-components/src/components/Slider/Slider.styles.ts

import { css } from '@emotion/react';
import {
  colors,
  spacing,
  transitions,
} from '@coze-studio/common/themes';

interface SliderStyleProps {
  value: number;
  min: number;
  max: number;
  disabled: boolean;
  isDragging: boolean;
}

export const useStyles = (props: SliderStyleProps) => {
  const { disabled } = props;

  return {
    container: css`
      position: relative;
      width: 100%;
      height: 12px;
      display: flex;
      align-items: center;
      cursor: ${disabled ? 'not-allowed' : 'pointer'};
    `,
    track: css`
      position: absolute;
      width: 100%;
      height: 4px;
      background-color: ${colors.gray[300]};
      border-radius: 2px;
    `,
    fill: css`
      position: absolute;
      height: 4px;
      background-color: ${colors.primary[500]};
      border-radius: 2px;
      transition: width ${transitions.base};
    `,
    thumb: css`
      position: absolute;
      width: 12px;
      height: 12px;
      background-color: white;
      border: 2px solid ${colors.primary[500]};
      border-radius: 50%;
      transform: translate(-50%, -50%);
      top: 50%;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
      transition: left ${transitions.base};
      cursor: ${disabled ? 'not-allowed' : 'grab'};

      ${!disabled && `
        &:hover {
          transform: translate(-50%, -50%) scale(1.2);
        }
      `}
    `,
  };
};
