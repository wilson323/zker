// frontend/packages/arch/ui-components/src/components/Card/Card.tsx

import React, { CSSProperties } from 'react';
import { useStyles } from './Card.styles';

export interface CardProps {
  /**
   * 卡片标题
   */
  title?: React.ReactNode;

  /**
   * 卡片内容
   */
  children: React.ReactNode;

  /**
   * 额外操作区域
   */
  extra?: React.ReactNode;

  /**
   * 是否显示边框
   */
  bordered?: boolean;

  /**
   * 是否可悬浮
   */
  hoverable?: boolean;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;

  /**
   * 点击回调
   */
  onClick?: () => void;
}

export const Card: React.FC<CardProps> = ({
  title,
  children,
  extra,
  bordered = true,
  hoverable = false,
  className,
  style,
  onClick,
}) => {
  const classes = useStyles({ bordered, hoverable });

  return (
    <div
      className={`${classes.card} ${className || ''} ${hoverable ? classes.hoverable : ''}`}
      style={style}
      onClick={onClick}
    >
      {(title || extra) && (
        <div className={classes.header}>
          <div className={classes.title}>{title}</div>
          {extra && <div className={classes.extra}>{extra}</div>}
        </div>
      )}
      <div className={classes.body}>{children}</div>
    </div>
  );
};

export default Card;
