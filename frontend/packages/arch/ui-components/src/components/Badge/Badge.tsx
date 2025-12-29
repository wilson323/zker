// frontend/packages/arch/ui-components/src/components/Badge/Badge.tsx

import React from 'react';
import { useStyles } from './Badge.styles';

export interface BadgeProps {
  count?: number;
  showZero?: boolean;
  max?: number;
  dot?: boolean;
  children?: React.ReactNode;
  className?: string;
}

export const Badge: React.FC<BadgeProps> = ({
  count = 0,
  showZero = false,
  max = 99,
  dot = false,
  children,
  className,
}) => {
  const classes = useStyles();

  const displayCount = count > max ? `${max}+` : count;
  const isHidden = count === 0 && !showZero;

  if (!children) {
    return (
      <sup className={`${classes.badge} ${className || ''}`}>
        {dot ? <span className={classes.dot} /> : displayCount}
      </sup>
    );
  }

  return (
    <div className={`${classes.container} ${className || ''}`}>
      {children}
      {!isHidden && (
        <sup className={classes.badge}>
          {dot ? <span className={classes.dot} /> : displayCount}
        </sup>
      )}
    </div>
  );
};
