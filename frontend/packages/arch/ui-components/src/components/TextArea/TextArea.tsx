// frontend/packages/arch/ui-components/src/components/TextArea/TextArea.tsx

import React from 'react';
import { useStyles } from './TextArea.styles';

export interface TextAreaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  size?: 'sm' | 'md' | 'lg';
  error?: boolean;
  block?: boolean;
}

export const TextArea: React.FC<TextAreaProps> = ({
  size = 'md',
  error = false,
  block = false,
  className,
  rows = 4,
  ...rest
}) => {
  const classes = useStyles({ size, error, block });

  return (
    <textarea
      className={`${classes.textarea} ${className || ''}`}
      rows={rows}
      {...rest}
    />
  );
};
