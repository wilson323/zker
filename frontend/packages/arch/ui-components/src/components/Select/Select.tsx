// frontend/packages/arch/ui-components/src/components/Select/Select.tsx

import React, { useState, useRef, useEffect } from 'react';
import { useStyles } from './Select.styles';

export interface SelectOption {
  label: string;
  value: string;
  disabled?: boolean;
}

export interface SelectProps {
  value?: string;
  onChange?: (value: string) => void;
  options: SelectOption[];
  placeholder?: string;
  disabled?: boolean;
  loading?: boolean;
  allowClear?: boolean;
  className?: string;
}

export const Select: React.FC<SelectProps> = ({
  value,
  onChange,
  options,
  placeholder = '请选择',
  disabled = false,
  loading = false,
  allowClear = false,
  className,
}) => {
  const classes = useStyles();
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  const selectedOption = options.find(opt => opt.value === value);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => document.removeEventListener('mousedown', handleClickOutside);
    }
  }, [isOpen]);

  const handleSelect = (optionValue: string) => {
    onChange?.(optionValue);
    setIsOpen(false);
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    onChange?.('');
  };

  return (
    <div ref={containerRef} className={`${classes.container} ${className || ''}`}>
      <div
        className={classes.trigger}
        onClick={() => !disabled && setIsOpen(!isOpen)}
      >
        {selectedOption ? (
          <span className={classes.value}>{selectedOption.label}</span>
        ) : (
          <span className={classes.placeholder}>{placeholder}</span>
        )}

        <div className={classes.suffix}>
          {loading && <span className={classes.spinner} />}
          {!loading && (
            <>
              {allowClear && selectedOption && (
                <span className={classes.clearIcon} onClick={handleClear}>
                  ×
                </span>
              )}
              <span className={`${classes.arrow} ${isOpen ? classes.arrowOpen : ''}`}>
                ▼
              </span>
            </>
          )}
        </div>
      </div>

      {isOpen && !disabled && (
        <div className={classes.dropdown}>
          {options.map(option => (
            <div
              key={option.value}
              className={`${classes.option} ${
                option.value === value ? classes.optionSelected : ''
              } ${
                option.disabled ? classes.optionDisabled : ''
              }`}
              onClick={() => !option.disabled && handleSelect(option.value)}
            >
              {option.label}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
