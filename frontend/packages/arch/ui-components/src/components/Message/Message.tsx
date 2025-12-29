// frontend/packages/arch/ui-components/src/components/Message/Message.tsx

import React, { useEffect, CSSProperties } from 'react';
import { useStyles } from './Message.styles';

export type MessageType = 'success' | 'info' | 'warning' | 'error';

export interface MessageProps {
  /**
   * 类型
   */
  type?: MessageType;

  /**
   * 提示内容
   */
  content: React.ReactNode;

  /**
   * 持续时间（毫秒），0 表示不自动关闭
   */
  duration?: number;

  /**
   * 关闭时的回调
   */
  onClose?: () => void;

  /**
   * 图标
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

const iconMap: Record<MessageType, string> = {
  success: '✓',
  info: 'ℹ',
  warning: '⚠',
  error: '✕',
};

export const Message: React.FC<MessageProps> = ({
  type = 'info',
  content,
  duration = 3000,
  onClose,
  icon,
  className,
  style,
}) => {
  const classes = useStyles({ type });
  const [visible, setVisible] = React.useState(true);

  useEffect(() => {
    if (duration > 0) {
      const timer = setTimeout(() => {
        handleClose();
      }, duration);
      return () => clearTimeout(timer);
    }
  }, [duration]);

  const handleClose = () => {
    setVisible(false);
    onClose?.();
  };

  if (!visible) {
    return null;
  }

  return (
    <div className={`${classes.container} ${className || ''}`} style={style}>
      <span className={classes.icon}>{icon || iconMap[type]}</span>
      <span className={classes.content}>{content}</span>
    </div>
  );
};

export default Message;
