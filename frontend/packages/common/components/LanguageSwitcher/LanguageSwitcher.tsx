// frontend/packages/common/components/LanguageSwitcher/LanguageSwitcher.tsx

import React from 'react';
import { useTranslation } from 'react-i18next';
import { useStyles } from './LanguageSwitcher.styles';

/**
 * 语言选项接口
 */
interface LanguageOption {
  value: string;
  label: string;
  icon?: string;
}

/**
 * LanguageSwitcher组件Props接口
 */
export interface LanguageSwitcherProps {
  /** 自定义className */
  className?: string;
  /** 是否显示标签 */
  showLabel?: boolean;
}

/**
 * LanguageSwitcher 语言切换器
 *
 * 用于切换应用语言
 *
 * @example
 * ```tsx
 * <LanguageSwitcher showLabel={false} />
 * ```
 */
export const LanguageSwitcher: React.FC<LanguageSwitcherProps> = ({
  className,
  showLabel = false,
}) => {
  const classes = useStyles();
  const { i18n } = useTranslation();

  const languageOptions: LanguageOption[] = [
    { value: 'zh-CN', label: '简体中文', icon: '🇨🇳' },
    { value: 'en-US', label: 'English', icon: '🇺🇸' },
  ];

  const handleChange = (lang: string) => {
    i18n.changeLanguage(lang);
    // 保存到localStorage
    localStorage.setItem('language', lang);
  };

  const currentLang = i18n.language;

  return (
    <div className={`${classes.container} ${className || ''}`}>
      {showLabel && (
        <span className={classes.label}>语言 / Language</span>
      )}
      <select
        className={classes.select}
        value={currentLang}
        onChange={(e) => handleChange(e.target.value)}
      >
        {languageOptions.map((option) => (
          <option key={option.value} value={option.value}>
            {option.icon && `${option.icon} `}{option.label}
          </option>
        ))}
      </select>
    </div>
  );
};

export default LanguageSwitcher;
