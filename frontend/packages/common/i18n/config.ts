// frontend/packages/common/i18n/config.ts

/**
 * ZKER 国际化配置
 *
 * 使用i18next实现多语言支持
 */

import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// 导入翻译文件
import zhCN from './locales/zh-CN.json';
import enUS from './locales/en-US.json';

// 资源配置
const resources = {
  'zh-CN': {
    translation: zhCN,
  },
  'en-US': {
    translation: enUS,
  },
};

// 初始化i18next
i18n
  .use(LanguageDetector) // 自动检测用户语言
  .use(initReactI18next) // 绑定react-i18next
  .init({
    resources,
    fallbackLng: 'zh-CN', // 默认语言
    lng: 'zh-CN', // 初始语言
    debug: process.env.NODE_ENV === 'development', // 开发模式启用调试

    interpolation: {
      escapeValue: false, // React已经做了XSS防护
    },

    react: {
      useSuspense: false, // 禁用Suspense
    },
  });

export default i18n;
