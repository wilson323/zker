# 性能优化指南

## 概述

本文档提供ZKER前端应用的完整性能优化策略和实施方案。

## 1. 代码分割（Code Splitting）

### 1.1 路由级代码分割

使用React.lazy和Suspense实现路由级懒加载：

```tsx
// studio/App.tsx
import { lazy, Suspense } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { LoadingWrapper } from '@coze-studio/common/components';

// 懒加载页面组件
const TenantList = lazy(() => import('./pages/tenant/TenantList'));
const TenantDetail = lazy(() => import('./pages/tenant/TenantDetail'));
const RoleList = lazy(() => import('./pages/permission/RoleList'));
const RoleDetail = lazy(() => import('./pages/permission/RoleDetail'));

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Suspense fallback={<LoadingWrapper loading />}>
        <Routes>
          <Route path="/tenants" element={<TenantList />} />
          <Route path="/tenants/:tenantId" element={<TenantDetail />} />
          <Route path="/permissions/roles" element={<RoleList />} />
          <Route path="/permissions/roles/:roleId" element={<RoleDetail />} />
        </Routes>
      </Suspense>
    </BrowserRouter>
  );
};
```

### 1.2 组件级代码分割

对于大型组件，使用动态导入：

```tsx
import { lazy, Suspense } from 'react';

const HeavyChart = lazy(() => import('./HeavyChart'));

export const Dashboard = () => {
  return (
    <div>
      <Suspense fallback={<div>加载图表中...</div>}>
        <HeavyChart />
      </Suspense>
    </div>
  );
};
```

## 2. Rsbuild配置优化

### 2.1 构建配置文件

```typescript
// rsbuild.config.ts
import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';

export default defineConfig({
  plugins: [pluginReact()],

  output: {
    distPath: {
      root: 'dist',
    },

    // 代码分割策略
    splitChunks: {
      strategy: 'split-by-experience',

      // 自定义分割配置
      override: {
        chunks: {
          'framework-react': [
            'react',
            'react-dom',
            'react-router-dom',
          ],
          'vendor-ui': [
            '@coze-studio/ui-components',
            '@coze-studio/business-components',
          ],
        },
      },
    },

    // 压缩配置
    minify: 'swc',

    // 目标浏览器
    targets: ['defaults', 'not IE 11'],

    // polyfill
    polyfill: 'entry',
  },

  performance: {
    // 移除console
    removeConsole: process.env.NODE_ENV === 'production',

    // 移除moment locales
    removeMomentLocale: true,

    // 打包体积分析
    bundleAnalyze: process.env.ANALYZE === 'true',

    // chunk大小告警
    chunkSizeWarningLimit: 500 * 1024, // 500KB
  },

  // 环境变量
  env: {
    API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:8080',
    NODE_ENV: process.env.NODE_ENV || 'development',
  },

  // source map
  sourceMap: {
    js: process.env.NODE_ENV === 'production' ? 'hidden-source-map' : true,
  },

  // CSS代码分割
  cssModules: {
    localIdentName: '[local]_[hash:base64:5]',
  },
});
```

### 2.2 生产环境优化

```typescript
// rsbuild.prod.config.ts
import { defineConfig } from '@rsbuild/core';

export default defineConfig({
  output: {
    minify: 'swc', // 使用SWC压缩，速度更快
    cssModules: {
      localIdentName: '[hash:base64:5]',
    },
  },

  performance: {
    gzip: true, // 启用gzip压缩
    brotli: true, // 启用brotli压缩
    removeConsole: true, // 移除console
    cssExtract: true, // 提取CSS到单独文件
  },
});
```

## 3. React性能优化

### 3.1 使用React.memo

```tsx
import { memo } from 'react';

interface ExpensiveComponentProps {
  data: any[];
}

export const ExpensiveComponent = memo<ExpensiveComponentProps>(({ data }) => {
  const sortedData = useMemo(() => {
    return data.sort((a, b) => a.createdAt - b.createdAt);
  }, [data]);

  return (
    <div>
      {sortedData.map(item => (
        <Item key={item.id} data={item} />
      ))}
    </div>
  );
});

// 自定义比较函数
export const MyComponent = memo(({ value }) => {
  return <div>{value}</div>;
}, (prevProps, nextProps) => {
  return prevProps.value === nextProps.value;
});
```

### 3.2 使用useMemo缓存计算结果

```tsx
import { useMemo } from 'react';

export const DataProcessor = ({ items }: { items: Item[] }) => {
  // 缓存昂贵的计算
  const processedData = useMemo(() => {
    return items
      .filter(item => item.active)
      .map(item => ({
        ...item,
        computed: item.value * 2,
      }))
      .sort((a, b) => a.computed - b.computed);
  }, [items]);

  return <DisplayList data={processedData} />;
};
```

### 3.3 使用useCallback稳定函数引用

```tsx
import { useCallback } from 'react';

export const ParentComponent = () => {
  const [data, setData] = useState([]);

  // 稳定的回调函数
  const handleDelete = useCallback((id: string) => {
    setData(prev => prev.filter(item => item.id !== id));
  }, []);

  const handleUpdate = useCallback((id: string, newData: any) => {
    setData(prev => prev.map(item =>
      item.id === id ? { ...item, ...newData } : item
    ));
  }, []);

  return (
    <ChildList
      items={data}
      onDelete={handleDelete}
      onUpdate={handleUpdate}
    />
  );
};
```

### 3.4 虚拟列表优化长列表

```tsx
import { useVirtualizer } from '@tanstack/react-virtual';

export const LongList = ({ items }: { items: Item[] }) => {
  const parentRef = React.useRef<HTMLDivElement>(null);

  const rowVirtualizer = useVirtualizer({
    count: items.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 50, // 每行估计高度
    overscan: 5,
  });

  return (
    <div ref={parentRef} style={{ height: '500px', overflow: 'auto' }}>
      <div style={{ height: `${rowVirtualizer.getTotalSize()}px` }}>
        {rowVirtualizer.getVirtualItems().map((virtualRow) => (
          <div
            key={virtualRow.key}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              height: `${virtualRow.size}px`,
            }}
          >
            <Row item={items[virtualRow.index]} />
          </div>
        ))}
      </div>
    </div>
  );
};
```

## 4. 资源优化

### 4.1 图片优化

```tsx
import Image from 'next/image';

// 使用Next.js Image组件自动优化
export const OptimizedImage = () => {
  return (
    <Image
      src="/logo.png"
      alt="Logo"
      width={200}
      height={100}
      priority // 关键图片预加载
      placeholder="blur" // 加载时显示模糊占位符
    />
  );
};
```

### 4.2 字体优化

```css
/* 使用font-display优化字体加载 */
@font-face {
  font-family: 'Custom Font';
  src: url('./font.woff2') format('woff2');
  font-display: swap; /* 立即显示后备字体 */
  font-weight: 400;
}
```

### 4.3 CDN配置

```typescript
// rsbuild.config.ts
export default defineConfig({
  output: {
    assetPrefix: process.env.CDN_URL, // CDN前缀
  },
});
```

## 5. 缓存策略

### 5.1 Service Worker缓存

```typescript
// service-worker.ts
const CACHE_NAME = 'zker-v1';
const urlsToCache = [
  '/',
  '/main.js',
  '/main.css',
  '/api/tenants',
];

self.addEventListener('install', (event: ExtendableEvent) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(urlsToCache);
    })
  );
});

self.addEventListener('fetch', (event: ExtendableEvent) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      return response || fetch(event.request);
    })
  );
});
```

### 5.2 HTTP缓存头

```nginx
# Nginx配置示例
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
  expires 1y;
  add_header Cache-Control "public, immutable";
}
```

## 6. 性能监控

### 6.1 Web Vitals监控

```tsx
import { useEffect } from 'react';
import { onCLS, onFID, onFCP, onLCP, onTTFB } from 'web-vitals';

export const PerformanceMonitor = () => {
  useEffect(() => {
    onCLS(console.log);
    onFID(console.log);
    onFCP(console.log);
    onLCP(console.log);
    onTTFB(console.log);
  }, []);

  return null;
};
```

### 6.2 性能指标目标

| 指标 | 目标值 | 说明 |
|------|--------|------|
| **FCP** | < 1.5s | First Contentful Paint |
| **LCP** | < 2.5s | Largest Contentful Paint |
| **TTI** | < 3.5s | Time to Interactive |
| **CLS** | < 0.1 | Cumulative Layout Shift |
| **FID** | < 100ms | First Input Delay |
| **TBT** | < 200ms | Total Blocking Time |

## 7. 构建优化命令

```json
// package.json
{
  "scripts": {
    "build": "rsbuild build",
    "build:analyze": "ANALYZE=true rsbuild build",
    "build:prod": "NODE_ENV=production rsbuild build"
  }
}
```

## 8. 最佳实践清单

### 开发阶段
- [ ] 使用React.memo避免不必要的重渲染
- [ ] 使用useMemo缓存计算结果
- [ ] 使用useCallback稳定函数引用
- [ ] 避免在render中创建对象/数组
- [ ] 使用key属性正确渲染列表
- [ ] 懒加载大型组件
- [ ] 使用虚拟列表处理长列表

### 构建阶段
- [ ] 启用代码分割
- [ ] 配置chunk splitting策略
- [ ] 启用压缩（Terser/SWC）
- [ ] 移除未使用的代码（Tree Shaking）
- [ ] 提取CSS到单独文件
- [ ] 使用source map便于调试
- [ ] 启用gzip/brotli压缩

### 部署阶段
- [ ] 配置CDN加速
- [ ] 设置合理的缓存策略
- [ ] 启用Service Worker
- [ ] 使用HTTP/2
- [ ] 优化图片和字体加载

---

**相关文档**：
- [Rsbuild性能优化文档](https://rsbuild.dev/guide/optimization/)
- [React性能优化文档](https://react.dev/learn/render-and-commit)
- [Web Vitals文档](https://web.dev/vitals/)
