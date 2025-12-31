// frontend/apps/coze-studio/e2e/auth.spec.ts

import { test, expect } from '@playwright/test';

/**
 * 认证流程E2E测试
 *
 * 测试场景:
 * 1. 用户登录
 * 2. 用户登出
 * 3. 未登录访问受保护页面重定向
 * 4. Token过期处理
 */

test.describe('认证流程 E2E 测试', () => {
  test('应该能够成功登录', async ({ page }) => {
    await page.goto('http://localhost:8888/login');

    // 填写登录表单
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password123');

    // 点击登录按钮
    await page.click('button[type="submit"]');

    // 等待跳转到首页
    await page.waitForURL('http://localhost:8888');

    // 验证登录成功
    await expect(page).toHaveURL('http://localhost:8888');
    await expect(page.locator('text=欢迎')).toBeVisible();
  });

  test('应该显示登录错误提示', async ({ page }) => {
    await page.goto('http://localhost:8888/login');

    // 填写错误的登录信息
    await page.fill('input[name="username"]', 'wronguser');
    await page.fill('input[name="password"]', 'wrongpass');

    // 点击登录按钮
    await page.click('button[type="submit"]');

    // 等待错误提示
    await page.waitForSelector('text=用户名或密码错误');

    // 验证错误提示
    await expect(page.locator('text=用户名或密码错误')).toBeVisible();
  });

  test('应该能够登出', async ({ page }) => {
    // 先登录
    await page.goto('http://localhost:8888/login');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('http://localhost:8888');

    // 点击登出按钮
    await page.click('button:has-text("登出")');

    // 等待跳转到登录页
    await page.waitForURL('http://localhost:8888/login');

    // 验证登出成功
    await expect(page).toHaveURL('http://localhost:8888/login');
  });

  test('未登录访问受保护页面应该重定向到登录页', async ({ page }) => {
    // 直接访问受保护页面
    await page.goto('http://localhost:8888/tenants');

    // 验证重定向到登录页
    await page.waitForURL('http://localhost:8888/login');
    await expect(page).toHaveURL('http://localhost:8888/login');
  });

  test('应该能够记住登录状态', async ({ page }) => {
    // 登录
    await page.goto('http://localhost:8888/login');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('http://localhost:8888');

    // 刷新页面
    await page.reload();

    // 验证仍然登录
    await expect(page).toHaveURL('http://localhost:8888');
    await expect(page.locator('text=欢迎')).toBeVisible();
  });
});
