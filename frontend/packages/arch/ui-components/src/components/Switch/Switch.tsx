// frontend/packages/arch/ui-components/src/components/Switch/Switch.tsx

import React from 'react';
import { useStyles } from './Switch.styles';

export interface SwitchProps {
  checked?: boolean;
  onChange?: (checked: boolean) => void;
  disabled?: boolean;
  className?: string;
}

export const Switch: React.FC<SwitchProps> = ({
  checked = false,
  onChange,
  disabled = false,
  className,
}) => {
  const classes = useStyles({ checked, disabled });

  const handleClick = () => {
    if (!disabled) {
      onChange?.(!checked);
    }
  };

  return (
    <button
      type="button"
      className={`${classes.switch} ${className || ''}`}
      onClick={handleClick}
      disabled={disabled}
    >
      <span className={classes.thumb} />
    </button>
  );
};
