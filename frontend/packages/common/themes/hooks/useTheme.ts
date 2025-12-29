// frontend/packages/common/themes/hooks/useTheme.ts

import { useMemo } from 'react';
import {
  colors,
  spacing,
  typography,
  borderRadius,
  boxShadow,
  breakpoints,
  transitions,
  zIndex,
} from '../tokens';

/**
 * 主题接口定义
 */
export interface Theme {
  colors: typeof colors;
  spacing: typeof spacing;
  typography: typeof typography;
  borderRadius: typeof borderRadius;
  boxShadow: typeof boxShadow;
  breakpoints: typeof breakpoints;
  transitions: typeof transitions;
  zIndex: typeof zIndex;
}

/**
 * useTheme Hook
 *
 * 提供统一的设计令牌访问
 * 在组件中使用此Hook获取主题配置
 *
 * @example
 * ```tsx
 * import { useTheme } from '@coze-studio/common/themes';
 *
 * export const MyComponent = () => {
 *   const theme = useTheme();
 *   return (
 *     <div style={{ color: theme.colors.primary[500] }}>
 *       Hello
 *     </div>
 *   );
 * };
 * ```
 */
export const useTheme = (): Theme => {
  return useMemo(() => ({
    colors,
    spacing,
    typography,
    borderRadius,
    boxShadow,
    breakpoints,
    transitions,
    zIndex,
  }), []);
};

export default useTheme;
