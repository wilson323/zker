// frontend/apps/coze-studio/rsbuild.prod.config.ts

import { defineConfig } from '@rsbuild/core';

export default defineConfig({
  output: {
    minify: 'swc',
    cssModules: {
      localIdentName: '[hash:base64:5]',
    },
  },

  performance: {
    gzip: true,
    brotli: true,
    removeConsole: true,
    removeMomentLocale: true,
    chunkSizeWarningLimit: 500 * 1024,
  },

  source: {
    sourceMap: {
      js: 'hidden-source-map',
    },
  },
});
