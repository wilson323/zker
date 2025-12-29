// frontend/packages/arch/ui-components/src/components/Slider/Slider.tsx

import React, { useState } from 'react';
import { useStyles } from './Slider.styles';

export interface SliderProps {
  min?: number;
  max?: number;
  value?: number;
  onChange?: (value: number) => void;
  disabled?: boolean;
  step?: number;
  className?: string;
}

export const Slider: React.FC<SliderProps> = ({
  min = 0,
  max = 100,
  value = 0,
  onChange,
  disabled = false,
  step = 1,
  className,
}) => {
  const [isDragging, setIsDragging] = useState(false);
  const classes = useStyles({ value, min, max, disabled, isDragging });

  const percentage = ((value - min) / (max - min)) * 100;

  const handleClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (disabled) return;

    const rect = e.currentTarget.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const newValue = Math.round((x / rect.width) * (max - min) + min);
    const steppedValue = Math.round(newValue / step) * step;
    const clampedValue = Math.max(min, Math.min(max, steppedValue));

    onChange?.(clampedValue);
  };

  return (
    <div
      className={`${classes.container} ${className || ''}`}
      onClick={handleClick}
    >
      <div className={classes.track} />
      <div
        className={classes.fill}
        style={{ width: `${percentage}%` }}
      />
      <div
        className={classes.thumb}
        style={{ left: `${percentage}%` }}
      />
    </div>
  );
};
