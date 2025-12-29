// frontend/packages/arch/ui-components/src/components/Dropdown/Dropdown.tsx

import React, { useRef, useEffect, CSSProperties } from 'react';
import { useStyles } from './Dropdown.styles';

export interface DropdownOption {
  key: string;
  label: React.ReactNode;
  disabled?: boolean;
  danger?: boolean;
  divided?: boolean;
  onClick?: () => void;
  children?: DropdownOption[];
}

export interface DropdownProps {
  /**
   * 下拉菜单项
   */
  menu: DropdownOption[];

  /**
   * 触发元素
   */
  trigger: React.ReactNode;

  /**
   * 触发方式
   */
  triggerType?: 'click' | 'hover';

  /**
   * 是否禁用
   */
  disabled?: boolean;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 菜单样式
   */
  menuStyle?: CSSProperties;
}

export const Dropdown: React.FC<DropdownProps> = ({
  menu,
  trigger,
  triggerType = 'click',
  disabled = false,
  className,
  menuStyle,
}) => {
  const classes = useStyles();
  const [visible, setVisible] = React.useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const timerRef = useRef<NodeJS.Timeout>();

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setVisible(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const handleMouseEnter = () => {
    if (triggerType === 'hover' && !disabled) {
      timerRef.current = setTimeout(() => {
        setVisible(true);
      }, 100);
    }
  };

  const handleMouseLeave = () => {
    if (triggerType === 'hover') {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
      setVisible(false);
    }
  };

  const handleClick = () => {
    if (triggerType === 'click' && !disabled) {
      setVisible(!visible);
    }
  };

  const handleMenuClick = (option: DropdownOption) => {
    if (option.disabled) return;
    option.onClick?.();
    setVisible(false);
  };

  const renderMenu = (items: DropdownOption[], depth = 0) => {
    return (
      <div className={classes.menu} style={menuStyle}>
        {items.map((item) => (
          <div
            key={item.key}
            className={`${classes.menuItem} ${item.disabled ? classes.disabled : ''} ${
              item.danger ? classes.danger : ''
            } ${item.divided ? classes.divided : ''}`}
            onClick={() => handleMenuClick(item)}
          >
            <span className={classes.label}>{item.label}</span>
            {item.children && (
              <div className={classes.submenu}>{renderMenu(item.children, depth + 1)}</div>
            )}
          </div>
        ))}
      </div>
    );
  };

  return (
    <div
      ref={containerRef}
      className={`${classes.container} ${className || ''}`}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      <div className={classes.trigger} onClick={handleClick}>
        {trigger}
      </div>
      {visible && !disabled && renderMenu(menu)}
    </div>
  );
};

export default Dropdown;
