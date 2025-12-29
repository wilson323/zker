// frontend/packages/arch/ui-components/src/components/Tooltip/Tooltip.tsx

import React, { CSSProperties } from 'react';
import { useStyles } from './Tooltip.styles';

export type TooltipPlacement = 'top' | 'bottom' | 'left' | 'right';

export interface TooltipProps {
  /**
   * 提示内容
   */
  title: React.ReactNode;

  /**
   * 位置
   */
  placement?: TooltipPlacement;

  /**
   * 子元素
   */
  children: React.ReactElement;

  /**
   * 延迟显示（毫秒）
   */
  delay?: number;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;
}

export const Tooltip: React.FC<TooltipProps> = ({
  title,
  placement = 'top',
  children,
  delay = 0,
  className,
  style,
}) => {
  const classes = useStyles({ placement });
  const [visible, setVisible] = React.useState(false);
  const timerRef = React.useRef<NodeJS.Timeout>();

  const handleMouseEnter = () => {
    if (delay) {
      timerRef.current = setTimeout(() => {
        setVisible(true);
      }, delay);
    } else {
      setVisible(true);
    }
  };

  const handleMouseLeave = () => {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
    }
    setVisible(false);
  };

  return (
    <div
      className={classes.wrapper}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      {visible && (
        <div className={`${classes.tooltip} ${className || ''}`} style={style}>
          {title}
          <span className={classes.arrow} />
        </div>
      )}
      {children}
    </div>
  );
};

export default Tooltip;
