// frontend/packages/common/themes/index.ts

/**
 * ZKER 主题系统
 *
 * 提供统一的设计令牌和主题Hook
 * 确保所有UI组件的视觉一致性
 */

// 导出设计令牌
export * from './tokens';

// 导出主题Hook
export { useTheme } from './hooks/useTheme';
export type { Theme } from './hooks/useTheme';
