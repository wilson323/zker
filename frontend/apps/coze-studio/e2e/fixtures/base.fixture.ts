/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * E2E 测试基础 Fixture
 *
 * @description 提供通用的测试辅助方法
 *
 * @author 研发C
 * @date 2025-01-04
 */

import { test as base } from '@playwright/test';

/**
 * 测试基础数据
 */
export interface TestData {
  tenantName: string;
  contactEmail: string;
  contactPhone?: string;
  companyName?: string;
  description?: string;
}

/**
 * 扩展的 test fixture
 *
 * 设计原则：
 * - SOLID: 单一职责，只提供基础测试方法
 * - DRY: 避免重复代码
 * - KISS: 保持简单
 */
export const test = base.extend<{
  /**
   * 导航到指定页面
   */
  navigateToPage: (path: string) => Promise<void>;

  /**
   * 登录系统
   */
  login: (username: string, password: string) => Promise<void>;

  /**
   * 登出系统
   */
  logout: () => Promise<void>;

  /**
   * 等待页面加载完成
   */
  waitForPageReady: () => Promise<void>;

  /**
   * 截图并保存
   */
  takeScreenshot: (name: string) => Promise<void>;
}>({
  /**
   * 导航到指定页面
   *
   * @param path - 页面路径
   */
  navigateToPage: async ({ page }, use) => {
    const navigate = async (path: string) => {
      await page.goto(path);
      await page.waitForLoadState('networkidle');
    };

    await use(navigate);
  },

  /**
   * 登录系统
   *
   * @param username - 用户名
   * @param password - 密码
   */
  login: async ({ page }, use) => {
    const doLogin = async (username: string, password: string) => {
      // 导航到登录页
      await page.goto('/login');

      // 填写登录表单
      await page.fill('input[name="username"]', username);
      await page.fill('input[name="password"]', password);

      // 提交表单
      await page.click('button[type="submit"]');

      // 等待登录成功
      await page.waitForURL('/', { timeout: 10000 });
    };

    await use(doLogin);
  },

  /**
   * 登出系统
   */
  logout: async ({ page }, use) => {
    const doLogout = async () => {
      // 点击用户头像
      await page.click('[data-testid="user-avatar"]');

      // 点击登出按钮
      await page.click('text=登出');

      // 等待跳转到登录页
      await page.waitForURL('/login');
    };

    await use(doLogout);
  },

  /**
   * 等待页面加载完成
   */
  waitForPageReady: async ({ page }, use) => {
    const waitReady = async () => {
      // 等待网络空闲
      await page.waitForLoadState('networkidle');

      // 等待 loading 消失
      await page.waitForSelector('[data-testid="loading"]', { state: 'hidden' }).catch(() => {
        // 忽略找不到 loading 的情况
      });
    };

    await use(waitReady);
  },

  /**
   * 截图并保存
   *
   * @param name - 截图名称
   */
  takeScreenshot: async ({ page }, use) => {
    const screenshot = async (name: string) => {
      await page.screenshot({
        path: `playwright-report/screenshots/${name}.png`,
        fullPage: true,
      });
    };

    await use(screenshot);
  },
});

/**
 * 导出 expect
 */
export const expect = test.expect;
