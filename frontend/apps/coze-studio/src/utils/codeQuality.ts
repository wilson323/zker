// frontend/apps/coze-studio/src/utils/codeQuality.ts

/**
 * 代码质量工具集
 *
 * 用途:
 * - 代码检查
 * - 性能监控
 * - 内存泄漏检测
 * - 最佳实践验证
 */

/**
 * 检查内存泄漏
 *
 * 用途: 检测事件监听器、定时器等是否正确清理
 */
export function checkMemoryLeaks() {
  const warnings: string[] = [];

  // 检查过多的定时器
  if ((window as any).__timers && (window as any).__timers.length > 50) {
    warnings.push(`过多定时器: ${(window as any).__timers.length}个`);
  }

  // 检查过多的事件监听器
  const elements = document.querySelectorAll('*');
  let totalListeners = 0;

  elements.forEach((element) => {
    const listeners = (element as any).eventListeners;
    if (listeners) {
      totalListeners += listeners.length;
    }
  });

  if (totalListeners > 500) {
    warnings.push(`过多事件监听器: ${totalListeners}个`);
  }

  return {
    hasLeaks: warnings.length > 0,
    warnings,
  };
}

/**
 * 性能监控
 *
 * 用途: 监控组件渲染性能
 */
export function monitorPerformance(componentName: string) {
  const startTime = performance.now();

  return {
    end: () => {
      const endTime = performance.now();
      const duration = endTime - startTime;

      if (duration > 100) {
        console.warn(
          `[Performance Warning] ${componentName} 渲染耗时: ${duration.toFixed(2)}ms`
        );
      }

      return {
        duration,
        startTime,
        endTime,
      };
    },
  };
}

/**
 * 代码复杂度检查
 *
 * 用途: 检查函数复杂度
 */
export function checkComplexity(fn: Function): {
  complexity: number;
  level: 'low' | 'medium' | 'high';
  suggestion?: string;
} {
  const fnString = fn.toString();

  // 计算圈复杂度(简化版)
  const patterns = [
    /if/g,
    /else/g,
    /for/g,
    /while/g,
    /switch/g,
    /case/g,
    /catch/g,
    /&&/g,
    /\|\|/g,
    /\?/g,
  ];

  let complexity = 1;
  patterns.forEach((pattern) => {
    const matches = fnString.match(pattern);
    if (matches) {
      complexity += matches.length;
    }
  });

  if (complexity <= 10) {
    return { complexity, level: 'low' };
  } else if (complexity <= 20) {
    return {
      complexity,
      level: 'medium',
      suggestion: '建议拆分函数以降低复杂度',
    };
  } else {
    return {
      complexity,
      level: 'high',
      suggestion: '复杂度过高,必须拆分函数',
    };
  }
}

/**
 * 组件最佳实践检查
 *
 * 用途: 检查React组件是否符合最佳实践
 */
export function checkReactBestPractices(component: React.ComponentType): {
  passed: boolean;
  issues: string[];
} {
  const issues: string[] = [];
  const componentString = component.toString();

  // 检查是否使用类组件(应该使用函数组件)
  if (componentString.includes('class ') && componentString.includes('extends ')) {
    issues.push('使用类组件,建议改用函数组件 + Hooks');
  }

  // 检查是否使用any类型
  if (componentString.includes(': any') || componentString.includes('<any>')) {
    issues.push('使用了any类型,应该使用具体类型');
  }

  // 检查是否缺少PropTypes/TypeScript
  if (!componentString.includes('propTypes') && !componentString.includes('interface')) {
    issues.push('缺少PropTypes或TypeScript类型定义');
  }

  // 检查是否有console.log
  if (componentString.includes('console.log')) {
    issues.push('包含console.log,应该在生产环境中移除');
  }

  return {
    passed: issues.length === 0,
    issues,
  };
}

/**
 * 依赖检查
 *
 * 用途: 检查依赖使用是否合理
 */
export function checkDependencies(
  deps: React.DependencyList,
  effectName: string
): {
  valid: boolean;
  issues: string[];
} {
  const issues: string[] = [];

  // 检查是否有过多的依赖
  if (deps.length > 10) {
    issues.push(
      `${effectName}: 依赖项过多(${deps.length}个),可能需要拆分effect`
    );
  }

  // 检查是否有函数依赖
  const functionDeps = deps.filter(
    (dep) => typeof dep === 'function'
  );
  if (functionDeps.length > 0) {
    issues.push(
      `${effectName}: 依赖数组包含函数,可能导致无限循环,考虑使用useCallback`
    );
  }

  // 检查是否有对象依赖
  const objectDeps = deps.filter((dep) => typeof dep === 'object');
  if (objectDeps.length > 0) {
    issues.push(
      `${effectName}: 依赖数组包含对象,可能导致无限循环,考虑使用useMemo`
    );
  }

  return {
    valid: issues.length === 0,
    issues,
  };
}

/**
 * Bundle分析
 *
 * 用途: 分析打包产物大小
 */
export function analyzeBundle(bundleStats: any) {
  const { assets } = bundleStats;

  const largeAssets = assets.filter((asset: any) => asset.size > 100000); // 100KB

  return {
    totalSize: assets.reduce((sum: number, asset: any) => sum + asset.size, 0),
    largeAssets,
    compressionRatio: 0.7, // 假设gzip压缩后减少30%
  };
}
