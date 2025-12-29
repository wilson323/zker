// frontend/packages/arch/ui-components/src/components/Radio/Radio.tsx

import React from 'react';
import { useStyles } from './Radio.styles';

export interface RadioProps {
  checked?: boolean;
  onChange?: (checked: boolean) => void;
  disabled?: boolean;
  value?: string;
  children?: React.ReactNode;
  className?: string;
}

export const Radio: React.FC<RadioProps> = ({
  checked = false,
  onChange,
  disabled = false,
  children,
  className,
}) => {
  const classes = useStyles({ checked, disabled });

  const handleChange = () => {
    if (!disabled) {
      onChange?.(!checked);
    }
  };

  return (
    <label className={`${classes.container} ${className || ''}`}>
      <span className={classes.radio}>
        <input
          type="radio"
          checked={checked}
          onChange={handleChange}
          disabled={disabled}
          className={classes.input}
        />
        <span className={classes.inner} />
      </span>
      {children && <span className={classes.children}>{children}</span>}
    </label>
  );
};
