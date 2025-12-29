// frontend/packages/arch/ui-components/src/components/Form/FormItem.tsx

import React from 'react';
import { useFormContext } from './Form';
import { useStyles } from './FormItem.styles';

export interface FormItemProps {
  name: string;
  label?: string;
  required?: boolean;
  children: React.ReactElement;
  className?: string;
}

export const FormItem: React.FC<FormItemProps> = ({
  name,
  label,
  required = false,
  children,
  className,
}) => {
  const classes = useStyles();
  const { values, errors, touched, setFieldValue } = useFormContext();

  const value = values[name];
  const error = touched[name] ? errors[name] : '';

  const handleChange = (newValue: any) => {
    setFieldValue(name, newValue);
  };

  return (
    <div className={`${classes.formItem} ${className || ''}`}>
      {label && (
        <label className={classes.label}>
          {label}
          {required && <span className={classes.required}>*</span>}
        </label>
      )}

      <div className={classes.content}>
        {React.cloneElement(children, {
          value,
          onChange: handleChange,
          error: !!error,
        } as any)}

        {error && <div className={classes.error}>{error}</div>}
      </div>
    </div>
  );
};
