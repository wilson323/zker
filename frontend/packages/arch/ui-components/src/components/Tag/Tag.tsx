// frontend/packages/arch/ui-components/src/components/Tag/Tag.tsx

import React from 'react';
import { useStyles } from './Tag.styles';

export interface TagProps {
  color?: 'default' | 'primary' | 'success' | 'warning' | 'error';
  children: React.ReactNode;
  closable?: boolean;
  onClose?: () => void;
  className?: string;
}

export const Tag: React.FC<TagProps> = ({
  color = 'default',
  children,
  closable = false,
  onClose,
  className,
}) => {
  const classes = useStyles({ color });

  return (
    <span className={`${classes.tag} ${className || ''}`}>
      {children}
      {closable && (
        <span className={classes.close} onClick={onClose}>
          ×
        </span>
      )}
    </span>
  );
};
