// frontend/apps/coze-studio/src/hooks/useKeyboardNavigation.ts

/**
 * 键盘导航Hook
 *
 * 用途: 为列表、菜单等组件添加键盘导航支持
 */

import { useEffect, useRef } from 'react';

interface UseKeyboardNavigationOptions {
  /**
   * 是否循环导航
   * @default true
   */
  loop?: boolean;

  /**
   * 导航方向
   * @default 'vertical'
   */
  orientation?: 'vertical' | 'horizontal' | 'both';

  /**
   * 项目选择回调
   */
  onSelect?: (index: number) => void;
}

export function useKeyboardNavigation<T extends HTMLElement>(
  itemCount: number,
  options: UseKeyboardNavigationOptions = {}
) {
  const { loop = true, orientation = 'vertical', onSelect } = options;
  const currentIndexRef = useRef(0);

  const handleKeyDown = (e: React.KeyboardEvent<T>) => {
    const prevKeys = orientation === 'vertical' ? ['ArrowUp'] : ['ArrowLeft'];
    const nextKeys = orientation === 'vertical' ? ['ArrowDown'] : ['ArrowRight'];

    if (orientation === 'both') {
      prevKeys.push('ArrowUp', 'ArrowLeft');
      nextKeys.push('ArrowDown', 'ArrowRight');
    }

    if (prevKeys.includes(e.key)) {
      e.preventDefault();
      currentIndexRef.current =
        currentIndexRef.current === 0
          ? loop
            ? itemCount - 1
            : 0
          : currentIndexRef.current - 1;
      onSelect?.(currentIndexRef.current);
    } else if (nextKeys.includes(e.key)) {
      e.preventDefault();
      currentIndexRef.current =
        currentIndexRef.current === itemCount - 1
          ? loop
            ? 0
            : itemCount - 1
          : currentIndexRef.current + 1;
      onSelect?.(currentIndexRef.current);
    } else if (e.key === 'Home') {
      e.preventDefault();
      currentIndexRef.current = 0;
      onSelect?.(0);
    } else if (e.key === 'End') {
      e.preventDefault();
      currentIndexRef.current = itemCount - 1;
      onSelect?.(itemCount - 1);
    } else if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      onSelect?.(currentIndexRef.current);
    }
  };

  return {
    currentIndex: currentIndexRef.current,
    handleKeyDown,
    setCurrentIndex: (index: number) => {
      currentIndexRef.current = index;
    },
  };
}

/**
 * 焦点管理Hook
 *
 * 用途: 管理组件的焦点状态
 */
export function useFocusManagement<T extends HTMLElement>() {
  const ref = useRef<T>(null);
  const previousFocusedElementRef = useRef<HTMLElement | null>(null);

  const trapFocus = () => {
    if (!ref.current) return;

    previousFocusedElementRef.current = document.activeElement as HTMLElement;

    const focusableElements = ref.current.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    );

    const firstElement = focusableElements[0] as HTMLElement;
    const lastElement = focusableElements[
      focusableElements.length - 1
    ] as HTMLElement;

    if (firstElement) {
      firstElement.focus();
    }

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key !== 'Tab') return;

      if (e.shiftKey && document.activeElement === firstElement) {
        e.preventDefault();
        lastElement?.focus();
      } else if (!e.shiftKey && document.activeElement === lastElement) {
        e.preventDefault();
        firstElement?.focus();
      }
    };

    ref.current.addEventListener('keydown', handleKeyDown);
  };

  const restoreFocus = () => {
    if (previousFocusedElementRef.current) {
      previousFocusedElementRef.current.focus();
    }
  };

  return { ref, trapFocus, restoreFocus };
}
