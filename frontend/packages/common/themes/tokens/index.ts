// frontend/packages/common/themes/tokens/index.ts

/**
 * ZKER 设计令牌系统
 *
 * 提供统一的设计规范，包括颜色、间距、字体、圆角、阴影等
 * 所有UI组件都应该使用这些令牌，确保视觉一致性
 */

// ==================== 颜色系统 ====================

/**
 * 主色调 (蓝色系)
 * 用于主要操作、链接、强调内容
 */
export const colors = {
  // 主色调
  primary: {
    50: '#E6F7FF',
    100: '#BAE7FF',
    200: '#91D5FF',
    300: '#69C0FF',
    400: '#40A9FF',
    500: '#1890FF',  // 主色
    600: '#096DD9',
    700: '#0050B3',
    800: '#003A8C',
    900: '#002766',
  },

  // 中性色 (灰色系)
  gray: {
    50: '#FAFAFA',
    100: '#F5F5F5',
    200: '#E8E8E8',
    300: '#D9D9D9',
    400: '#BFBFBF',
    500: '#8C8C8C',
    600: '#595959',
    700: '#434343',
    800: '#262626',
    900: '#1F1F1F',
  },

  // 语义色
  success: {
    light: '#95DE64',
    main: '#52C41A',
    dark: '#389E0D',
  },
  warning: {
    light: '#FFD666',
    main: '#FAAD14',
    dark: '#D48806',
  },
  error: {
    light: '#FF7875',
    main: '#F5222D',
    dark: '#CF1322',
  },
  info: {
    light: '#91D5FF',
    main: '#1890FF',
    dark: '#0050B3',
  },
};

// ==================== 间距系统 ====================

/**
 * 间距系统 (基于4的倍数)
 * 用于padding、margin、gap等
 */
export const spacing = {
  0: '0',
  1: '0.25rem',   // 4px
  2: '0.5rem',    // 8px
  3: '0.75rem',   // 12px
  4: '1rem',      // 16px
  5: '1.25rem',   // 20px
  6: '1.5rem',    // 24px
  8: '2rem',      // 32px
  10: '2.5rem',   // 40px
  12: '3rem',     // 48px
  16: '4rem',     // 64px
  20: '5rem',     // 80px
  24: '6rem',     // 96px
};

// ==================== 字体系统 ====================

/**
 * 字体系统
 * 定义字体族、字号、字重、行高
 */
export const typography = {
  fontFamily: {
    base: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
    mono: '"SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace',
  },
  fontSize: {
    xs: '0.75rem',    // 12px
    sm: '0.875rem',   // 14px
    base: '1rem',     // 16px
    lg: '1.125rem',   // 18px
    xl: '1.25rem',    // 20px
    '2xl': '1.5rem',  // 24px
    '3xl': '1.875rem', // 30px
    '4xl': '2.25rem',  // 36px
  },
  fontWeight: {
    normal: 400,
    medium: 500,
    semibold: 600,
    bold: 700,
  },
  lineHeight: {
    tight: 1.25,
    normal: 1.5,
    relaxed: 1.75,
  },
};

// ==================== 圆角系统 ====================

/**
 * 圆角系统
 * 用于border-radius
 */
export const borderRadius = {
  none: '0',
  sm: '0.125rem',   // 2px
  base: '0.25rem',  // 4px
  md: '0.375rem',   // 6px
  lg: '0.5rem',     // 8px
  xl: '0.75rem',    // 12px
  '2xl': '1rem',    // 16px
  full: '9999px',   // 完全圆角
};

// ==================== 阴影系统 ====================

/**
 * 阴影系统
 * 用于box-shadow
 */
export const boxShadow = {
  sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
  base: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)',
  md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
  lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
  xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
  inner: 'inset 0 2px 4px 0 rgba(0, 0, 0, 0.06)',
};

// ==================== 断点系统 ====================

/**
 * 响应式断点系统
 * 用于媒体查询
 */
export const breakpoints = {
  sm: '640px',
  md: '768px',
  lg: '1024px',
  xl: '1280px',
  '2xl': '1536px',
};

// ==================== 过渡动画 ====================

/**
 * 过渡动画
 * 用于transition
 */
export const transitions = {
  fast: '150ms cubic-bezier(0.4, 0, 0.2, 1)',
  base: '200ms cubic-bezier(0.4, 0, 0.2, 1)',
  slow: '300ms cubic-bezier(0.4, 0, 0.2, 1)',
};

// ==================== Z-index系统 ====================

/**
 * Z-index层级系统
 * 用于控制元素堆叠顺序
 */
export const zIndex = {
  dropdown: 1000,
  sticky: 1020,
  fixed: 1030,
  modalBackdrop: 1040,
  modal: 1050,
  popover: 1060,
  tooltip: 1070,
};

// ==================== 导出所有令牌 ====================

/**
 * 完整的设计令牌集合
 */
export const designTokens = {
  colors,
  spacing,
  typography,
  borderRadius,
  boxShadow,
  breakpoints,
  transitions,
  zIndex,
};

export default designTokens;
