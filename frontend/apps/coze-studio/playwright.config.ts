/**
 * 企业级Playwright配置
 * @version 1.0.0
 * @description E2E测试配置，支持多浏览器、并发、视频录制
 */

import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  // 测试目录
  testDir: './e2e',

  // 超时配置
  timeout: 30 * 1000,
  expect: {
    timeout: 5000,
  },

  // 完全并行运行测试
  fullyParallel: true,

  // CI环境允许失败重试
  retries: process.env.CI ? 2 : 0,

  // 并发worker数量
  workers: process.env.CI ? 1 : undefined,

  // 报告器配置
  reporter: [
    ['html', { outputFolder: 'playwright-report', open: 'never' }],
    ['json', { outputFile: 'test-results/test-results.json' }],
    ['junit', { outputFile: 'test-results/test-results.xml' }],
    ['list'],
  ],

  // 共享配置
  use: {
    // 基础URL
    baseURL: 'http://localhost:8888',

    // 追踪配置（失败时保留）
    trace: 'retain-on-failure',

    // 截图配置
    screenshot: 'only-on-failure',

    // 视频录制
    video: 'retain-on-failure',

    // 浏览器上下文选项
    viewport: { width: 1280, height: 720 },
    ignoreHTTPSErrors: true,
    actionTimeout: 10000,
    navigationTimeout: 30000,
  },

  // 测试项目配置（不同浏览器）
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },

    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },

    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },

    // 移动端测试
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
    },
  ],

  // 开发服务器配置
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:8888',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },

  // 输出目录
  outputDir: 'test-results',

  // 依赖测试前的全局设置
  globalSetup: './e2e/global-setup.ts',

  // 测试后的全局清理
  globalTeardown: './e2e/global-teardown.ts',
})
