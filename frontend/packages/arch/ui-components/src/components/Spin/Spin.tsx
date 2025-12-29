// frontend/packages/arch/ui-components/src/components/Spin/Spin.tsx

import React, { CSSProperties } from 'react';
import { useStyles } from './Spin.styles';

export interface SpinProps {
  /**
   * 加载中状态
   */
  spinning?: boolean;

  /**
   * 尺寸
   */
  size?: 'small' | 'default' | 'large';

  /**
   * 提示文本
   */
  tip?: string;

  /**
   * 延迟显示加载中（毫秒）
   */
  delay?: number;

  /**
   * 包装器样式
   */
  wrapperClassName?: string;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;

  /**
   * 子元素（当 spinning 为 false 时显示）
   */
  children?: React.ReactNode;
}

export const Spin: React.FC<SpinProps> = ({
  spinning = true,
  size = 'default',
  tip,
  delay = 0,
  wrapperClassName,
  className,
  style,
  children,
}) => {
  const classes = useStyles({ size });
  const [showSpin, setShowSpin] = React.useState(!delay);

  React.useEffect(() => {
    if (delay && spinning) {
      const timer = setTimeout(() => {
        setShowSpin(true);
      }, delay);
      return () => clearTimeout(timer);
    } else {
      setShowSpin(spinning);
    }
  }, [spinning, delay]);

  if (children) {
    return (
      <div className={`${classes.wrapper} ${wrapperClassName || ''}`}>
        {showSpin && (
          <div className={`${classes.overlay} ${className || ''}`} style={style}>
            <span className={classes.spinner}>
              <span className={classes.dot} />
              <span className={classes.dot} />
              <span className={classes.dot} />
            </span>
            {tip && <div className={classes.tip}>{tip}</div>}
          </div>
        )}
        {children}
      </div>
    );
  }

  if (!showSpin) {
    return null;
  }

  return (
    <div className={`${classes.container} ${className || ''}`} style={style}>
      <span className={classes.spinner}>
        <span className={classes.dot} />
        <span className={classes.dot} />
        <span className={classes.dot} />
      </span>
      {tip && <div className={classes.tip}>{tip}</div>}
    </div>
  );
};

export default Spin;
