// frontend/packages/arch/ui-components/src/components/DatePicker/DatePicker.tsx

import React, { useState, CSSProperties } from 'react';
import { useStyles } from './DatePicker.styles';
import { Input } from '../Input';

export type DatePickerType = 'date' | 'datetime' | 'dateRange' | 'month' | 'year';

export interface DatePickerProps {
  /**
   * 选择器类型
   */
  type?: DatePickerType;

  /**
   * 值
   */
  value?: string | string[];

  /**
   * 占位文本
   */
  placeholder?: string;

  /**
   * 是否禁用
   */
  disabled?: boolean;

  /**
   * 日期格式
   */
  format?: string;

  /**
   * 最小日期
   */
  minDate?: string;

  /**
   * 最大日期
   */
  maxDate?: string;

  /**
   * 时区
   */
  timeZone?: string;

  /**
   * 变更回调
   */
  onChange?: (value: string | string[]) => void;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;
}

export const DatePicker: React.FC<DatePickerProps> = ({
  type = 'date',
  value,
  placeholder = '请选择日期',
  disabled = false,
  format = 'YYYY-MM-DD',
  minDate,
  maxDate,
  timeZone,
  onChange,
  className,
  style,
}) => {
  const classes = useStyles();
  const [isOpen, setIsOpen] = useState(false);
  const [currentDate, setCurrentDate] = useState(new Date());
  const containerRef = React.useRef<HTMLDivElement>(null);

  // 点击外部关闭
  React.useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // 格式化日期显示
  const formatDisplayDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return format
      .replace('YYYY', String(year))
      .replace('MM', month)
      .replace('DD', day);
  };

  // 处理日期选择
  const handleDateClick = (date: Date) => {
    const formatted = formatDisplayDate(date);
    onChange?.(formatted);
    setIsOpen(false);
  };

  // 生成日历
  const generateCalendar = () => {
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const startDay = firstDay.getDay();
    const totalDays = lastDay.getDate();

    const days = [];

    // 填充前面的空白
    for (let i = 0; i < startDay; i++) {
      days.push(<div key={`empty-${i}`} className={classes.emptyDay} />);
    }

    // 填充日期
    for (let day = 1; day <= totalDays; day++) {
      const date = new Date(year, month, day);
      const dateStr = formatDisplayDate(date);
      const isSelected = value === dateStr;
      const isToday =
        date.getDate() === new Date().getDate() &&
        date.getMonth() === new Date().getMonth() &&
        date.getFullYear() === new Date().getFullYear();

      days.push(
        <div
          key={day}
          className={`${classes.day} ${isSelected ? classes.selectedDay : ''} ${
            isToday ? classes.today : ''
          }`}
          onClick={() => handleDateClick(date)}
        >
          {day}
        </div>
      );
    }

    return days;
  };

  // 上个月
  const prevMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1));
  };

  // 下个月
  const nextMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1));
  };

  const displayValue = value || '';

  return (
    <div
      ref={containerRef}
      className={`${classes.container} ${className || ''}`}
      style={style}
    >
      <Input
        value={displayValue}
        placeholder={placeholder}
        disabled={disabled}
        readOnly
        onFocus={() => !disabled && setIsOpen(true)}
      />

      {isOpen && (
        <div className={classes.panel}>
          <div className={classes.header}>
            <button className={classes.navButton} onClick={prevMonth}>
              ‹
            </button>
            <span className={classes.title}>
              {currentDate.getFullYear()}年 {currentDate.getMonth() + 1}月
            </span>
            <button className={classes.navButton} onClick={nextMonth}>
              ›
            </button>
          </div>

          <div className={classes.weekHeader}>
            {['日', '一', '二', '三', '四', '五', '六'].map((week) => (
              <div key={week} className={classes.weekDay}>
                {week}
              </div>
            ))}
          </div>

          <div className={classes.calendar}>{generateCalendar()}</div>
        </div>
      )}
    </div>
  );
};

export default DatePicker;
