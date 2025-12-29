// frontend/packages/arch/ui-components/src/components/Alert/Alert.tsx

import React, { CSSProperties } from 'react';
import { useStyles } from './Alert.styles';

export type AlertType = 'success' | 'info' | 'warning' | 'error';

export interface AlertProps {
  /**
   * 类型
   */
  type?: AlertType;

  /**
   * 提示内容
   */
  message: React.ReactNode;

  /**
   * 辅助性文字
   */
  description?: React.ReactNode;

  /**
   * 是否显示关闭图标
   */
  closable?: boolean;

  /**
   * 关闭时的回调
   */
  onClose?: () => void;

  /**
   * 是否显示图标
   */
  showIcon?: boolean;

  /**
   * 自定义图标
   */
  icon?: React.ReactNode;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;
}

const iconMap: Record<AlertType, string> = {
  success: '✓',
  info: 'ℹ',
  warning: '⚠',
  error: '✕',
};

export const Alert: React.FC<AlertProps> = ({
  type = 'info',
  message,
  description,
  closable = false,
  onClose,
  showIcon = false,
  icon,
  className,
  style,
}) => {
  const classes = useStyles({ type });
  const [visible, setVisible] = React.useState(true);

  const handleClose = () => {
    setVisible(false);
    onClose?.();
  };

  if (!visible) {
    return null;
  }

  const renderIcon = () => {
    if (!showIcon) return null;
    return <span className={classes.icon}>{icon || iconMap[type]}</span>;
  };

  return (
    <div className={`${classes.container} ${className || ''}`} style={style}>
      {renderIcon()}
      <div className={classes.content}>
        <div className={classes.message}>{message}</div>
        {description && <div className={classes.description}>{description}</div>}
      </div>
      {closable && (
        <button className={classes.closeButton} onClick={handleClose}>
          ✕
        </button>
      )}
    </div>
  );
};

export default Alert;
