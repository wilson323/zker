// frontend/packages/arch/ui-components/src/components/Checkbox/Checkbox.tsx

import React from 'react';
import { useStyles } from './Checkbox.styles';

export interface CheckboxProps {
  checked?: boolean;
  onChange?: (checked: boolean) => void;
  disabled?: boolean;
  indeterminate?: boolean;
  children?: React.ReactNode;
  className?: string;
}

export const Checkbox: React.FC<CheckboxProps> = ({
  checked = false,
  onChange,
  disabled = false,
  indeterminate = false,
  children,
  className,
}) => {
  const classes = useStyles({ checked, disabled, indeterminate });

  const handleChange = () => {
    if (!disabled) {
      onChange?.(!checked);
    }
  };

  return (
    <label className={`${classes.container} ${className || ''}`}>
      <span className={classes.checkbox}>
        <input
          type="checkbox"
          checked={checked}
          onChange={handleChange}
          disabled={disabled}
          className={classes.input}
        />
        <span className={classes.inner}>
          {indeterminate ? (
            <span className={classes.indeterminateIcon}>-</span>
          ) : checked ? (
            <span className={classes.checkIcon}>✓</span>
          ) : null}
        </span>
      </span>
      {children && <span className={classes.children}>{children}</span>}
    </label>
  );
};
