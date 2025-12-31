// frontend/apps/coze-studio/src/utils/accessibility.ts

/**
 * 可访问性工具集
 *
 * 用途:
 * - ARIA属性
 * - 键盘导航
 * - 屏幕阅读器支持
 * - 焦点管理
 */

/**
 * 生成唯一的ID
 */
let idCounter = 0;
export function generateId(prefix: string = 'id'): string {
  return `${prefix}-${++idCounter}`;
}

/**
 * 管理焦点
 *
 * 用途: 在模态框、下拉菜单等场景中管理焦点
 */
export function manageFocus(
  element: HTMLElement,
  action: 'trap' | 'restore' | 'set'
): void {
  switch (action) {
    case 'trap':
      // 焦点陷阱: Tab键只在元素内循环
      element.addEventListener('keydown', (e) => {
        if (e.key === 'Tab') {
          const focusableElements = element.querySelectorAll(
            'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
          );
          const firstElement = focusableElements[0] as HTMLElement;
          const lastElement = focusableElements[
            focusableElements.length - 1
          ] as HTMLElement;

          if (e.shiftKey && document.activeElement === firstElement) {
            e.preventDefault();
            lastElement.focus();
          } else if (!e.shiftKey && document.activeElement === lastElement) {
            e.preventDefault();
            firstElement.focus();
          }
        }
      });
      break;

    case 'restore':
      // 恢复之前的焦点
      const previousActiveElement = document.activeElement as HTMLElement;
      if (previousActiveElement) {
        previousActiveElement.focus();
      }
      break;

    case 'set':
      // 设置焦点
      element.focus();
      break;
  }
}

/**
 * 通告屏幕阅读器
 *
 * 用途: 向屏幕阅读器宣告状态变化
 */
export function announceToScreenReader(message: string, priority: 'polite' | 'assertive' = 'polite'): void {
  const announcement = document.createElement('div');
  announcement.setAttribute('role', 'status');
  announcement.setAttribute('aria-live', priority);
  announcement.setAttribute('aria-atomic', 'true');
  announcement.className = 'sr-only';
  announcement.textContent = message;

  document.body.appendChild(announcement);

  // 移除元素
  setTimeout(() => {
    document.body.removeChild(announcement);
  }, 1000);
}

/**
 * 检查颜色对比度
 *
 * 用途: 确保文本颜色符合WCAG标准
 */
export function checkColorContrast(
  foreground: string,
  background: string
): { ratio: number; wcagAA: boolean; wcagAAA: boolean } {
  // 将hex转换为RGB
  const hexToRgb = (hex: string) => {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result
      ? {
          r: parseInt(result[1], 16),
          g: parseInt(result[2], 16),
          b: parseInt(result[3], 16),
        }
      : null;
  };

  const fg = hexToRgb(foreground);
  const bg = hexToRgb(background);

  if (!fg || !bg) {
    return { ratio: 1, wcagAA: false, wcagAAA: false };
  }

  // 计算相对亮度
  const luminance = (r: number, g: number, b: number) => {
    const a = [r, g, b].map((v) => {
      v /= 255;
      return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4);
    });
    return a[0] * 0.2126 + a[1] * 0.7152 + a[2] * 0.0722;
  };

  const lum1 = luminance(fg.r, fg.g, fg.b);
  const lum2 = luminance(bg.r, bg.g, bg.b);
  const brightest = Math.max(lum1, lum2);
  const darkest = Math.min(lum1, lum2);

  const ratio = (brightest + 0.05) / (darkest + 0.05);

  return {
    ratio,
    wcagAA: ratio >= 4.5,
    wcagAAA: ratio >= 7,
  };
}

/**
 * 键盘导航Hook
 *
 * 用途: 添加键盘导航支持
 */
export function useKeyboardNavigation(
  items: HTMLElement[],
  options?: {
    loop?: boolean;
    orientation?: 'horizontal' | 'vertical';
  }
) {
  const { loop = true, orientation = 'vertical' } = options || {};

  const handleKeyDown = (e: KeyboardEvent, currentIndex: number) => {
    const nextKey = orientation === 'vertical' ? 'ArrowDown' : 'ArrowRight';
    const prevKey = orientation === 'vertical' ? 'ArrowUp' : 'ArrowLeft';

    switch (e.key) {
      case nextKey:
        e.preventDefault();
        const nextIndex = (currentIndex + 1) % items.length;
        items[nextIndex].focus();
        break;
      case prevKey:
        e.preventDefault();
        const prevIndex = currentIndex === 0 ? items.length - 1 : currentIndex - 1;
        items[prevIndex].focus();
        break;
      case 'Home':
        e.preventDefault();
        items[0].focus();
        break;
      case 'End':
        e.preventDefault();
        items[items.length - 1].focus();
        break;
    }
  };

  return { handleKeyDown };
}

/**
 * ARIA属性生成器
 *
 * 用途: 为组件生成正确的ARIA属性
 */
export function getAriaProps(config: {
  role?: string;
  label?: string;
  describedBy?: string;
  expanded?: boolean;
  checked?: boolean;
  disabled?: boolean;
  required?: boolean;
  invalid?: boolean;
  live?: 'polite' | 'assertive' | 'off';
}): Record<string, string | boolean> {
  const props: Record<string, string | boolean> = {};

  if (config.role) props['role'] = config.role;
  if (config.label) props['aria-label'] = config.label;
  if (config.describedBy) props['aria-describedby'] = config.describedBy;
  if (config.expanded !== undefined) props['aria-expanded'] = config.expanded;
  if (config.checked !== undefined) props['aria-checked'] = config.checked;
  if (config.disabled !== undefined) props['aria-disabled'] = config.disabled;
  if (config.required !== undefined) props['aria-required'] = config.required;
  if (config.invalid !== undefined) props['aria-invalid'] = config.invalid;
  if (config.live) props['aria-live'] = config.live;

  return props;
}
