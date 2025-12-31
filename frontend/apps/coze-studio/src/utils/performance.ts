// frontend/apps/coze-studio/src/utils/performance.ts

/**
 * 性能优化工具集
 *
 * 用途:
 * - 防抖(Debounce)
 * - 节流(Throttle)
 * - 内存优化
 * - 批量更新
 */

/**
 * 防抖函数
 *
 * 用途: 延迟执行函数,在等待时间内多次调用只执行最后一次
 *
 * @example
 * ```tsx
 * const debouncedSearch = useDebounce(search, 300);
 * ```
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout> | null = null;

  return function (this: any, ...args: Parameters<T>) {
    if (timeoutId !== null) {
      clearTimeout(timeoutId);
    }

    timeoutId = setTimeout(() => {
      func.apply(this, args);
    }, wait);
  };
}

/**
 * 节流函数
 *
 * 用途: 限制函数执行频率,在指定时间内只执行一次
 *
 * @example
 * ```tsx
 * const throttledScroll = useThrottle(handleScroll, 100);
 * ```
 */
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle = false;
  let lastResult: ReturnType<T>;

  return function (this: any, ...args: Parameters<T>) {
    if (!inThrottle) {
      inThrottle = true;
      lastResult = func.apply(this, args);

      setTimeout(() => {
        inThrottle = false;
      }, limit);
    }

    return lastResult;
  };
}

/**
 * 批量更新Hook
 *
 * 用途: 批量更新状态,避免多次渲染
 *
 * @example
 * ```tsx
 * const batchUpdates = useBatchUpdates();
 * batchUpdates(() => {
 *   setState1(value1);
 *   setState2(value2);
 * });
 * ```
 */
export function batchUpdates(updates: () => void) {
  // React 18自动批处理,此函数为兼容性保留
  updates();
}

/**
 * 内存优化 - 清理大对象
 *
 * 用途: 释放大对象占用的内存
 */
export function cleanupLargeObject(obj: any): void {
  if (obj && typeof obj === 'object') {
    Object.keys(obj).forEach((key) => {
      delete obj[key];
    });
  }
}

/**
 * 懒加载图片Hook
 *
 * 用途: 使用Intersection Observer懒加载图片
 */
export function useLazyLoadImage(
  src: string,
  options?: IntersectionObserverInit
) {
  const [imageSrc, setImageSrc] = React.useState<string>();
  const imgRef = React.useRef<HTMLImageElement>(null);

  React.useEffect(() => {
    const observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting) {
        setImageSrc(src);
        observer.disconnect();
      }
    }, options);

    if (imgRef.current) {
      observer.observe(imgRef.current);
    }

    return () => observer.disconnect();
  }, [src, options]);

  return { imgRef, imageSrc };
}

/**
 * 虚拟滚动Hook
 *
 * 用途: 渲染大列表时只渲染可见项
 */
export function useVirtualList<T>({
  items,
  itemHeight,
  containerHeight,
  overscan = 3,
}: {
  items: T[];
  itemHeight: number;
  containerHeight: number;
  overscan?: number;
}) {
  const [scrollTop, setScrollTop] = React.useState(0);

  const startIndex = Math.max(0, Math.floor(scrollTop / itemHeight) - overscan);
  const endIndex = Math.min(
    items.length,
    Math.ceil((scrollTop + containerHeight) / itemHeight) + overscan
  );

  const visibleItems = items.slice(startIndex, endIndex);
  const totalHeight = items.length * itemHeight;
  const offsetY = startIndex * itemHeight;

  return {
    visibleItems,
    totalHeight,
    offsetY,
    onScroll: (e: React.UIEvent<HTMLDivElement>) => {
      setScrollTop(e.currentTarget.scrollTop);
    },
  };
}

// React import needed for hooks above
import React from 'react';
