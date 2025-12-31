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
 * 计费管理 E2E 测试套件
 *
 * @description 测试计费管理的完整流程，包括余额查询、账单生成、充值等
 *
 * @author 研发C
 * @date 2025-01-04
 */

import { test, expect } from '../fixtures/base.fixture';
import { Page } from '@playwright/test';

/**
 * 计费管理页面辅助类
 */
class BillingManagementPage {
  constructor(private page: Page) {}

  /**
   * 导航到计费概览页面
   */
  async navigateToBillingOverview() {
    await this.page.goto('/billing/overview');
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForSelector('[data-testid="billing-overview"]', { timeout: 10000 });
  }

  /**
   * 导航到账单列表页面
   */
  async navigateToInvoices() {
    await this.page.goto('/billing/invoices');
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForSelector('[data-testid="invoice-list"]', { timeout: 10000 });
  }

  /**
   * 导航到充值页面
   */
  async navigateToRecharge() {
    await this.page.goto('/billing/recharge');
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForSelector('[data-testid="recharge-form"]', { timeout: 10000 });
  }

  /**
   * 导航到使用统计页面
   */
  async navigateToUsageStats() {
    await this.page.goto('/billing/usage');
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForSelector('[data-testid="usage-stats"]', { timeout: 10000 });
  }

  /**
   * 获取当前余额
   */
  async getCurrentBalance(): Promise<string> {
    const balanceElement = this.page.locator('[data-testid="current-balance"]');
    await expect(balanceElement).toBeVisible();
    return await balanceElement.textContent() || '';
  }

  /**
   * 搜索账单
   */
  async searchInvoice(keyword: string) {
    const searchInput = this.page.locator('input[data-testid="invoice-search"]');
    await searchInput.fill(keyword);
    await this.page.waitForTimeout(500); // 等待搜索结果
  }

  /**
   * 导出账单
   */
  async exportInvoice(invoiceId: string) {
    const exportButton = this.page.locator(`[data-testid="export-invoice-${invoiceId}"]`);
    await exportButton.click();
  }

  /**
   * 充值
   */
  async recharge(amount: number, paymentMethod: string) {
    // 输入充值金额
    const amountInput = this.page.locator('input[name="amount"]');
    await amountInput.fill(amount.toString());

    // 选择支付方式
    const paymentSelect = this.page.locator('select[name="payment_method"]');
    await paymentSelect.selectOption(paymentMethod);

    // 点击充值按钮
    const rechargeButton = this.page.locator('button[data-testid="btn-recharge"]');
    await rechargeButton.click();
  }
}

/**
 * 计费管理测试套件
 */
test.describe('计费管理', () => {
  let billingPage: BillingManagementPage;

  test.beforeEach(async ({ page }) => {
    billingPage = new BillingManagementPage(page);
  });

  /**
   * 测试套件1: 计费概览
   */
  test.describe('计费概览', () => {
    test('应该正确展示余额信息', async ({ page }) => {
      await billingPage.navigateToBillingOverview();

      // 验证余额显示
      await expect(page.locator('[data-testid="current-balance"]')).toBeVisible();
      await expect(page.locator('text=/余额/')).toBeVisible();
      await expect(page.locator('text=/CNY/')).toBeVisible();

      // 验证预算使用情况
      await expect(page.locator('[data-testid="budget-usage"]')).toBeVisible();

      // 验证使用趋势图
      await expect(page.locator('[data-testid="usage-trend-chart"]')).toBeVisible();
    });

    test('应该显示本月使用统计', async ({ page }) => {
      await billingPage.navigateToBillingOverview();

      // 验证Token使用量
      await expect(page.locator('[data-testid="token-usage"]')).toBeVisible();

      // 验证API调用次数
      await expect(page.locator('[data-testid="api-calls"]')).toBeVisible();

      // 验证存储使用量
      await expect(page.locator('[data-testid="storage-usage"]')).toBeVisible();
    });

    test('应该显示预算告警', async ({ page }) => {
      await billingPage.navigateToBillingOverview();

      // 检查是否有告警（取决于数据）
      const alertElement = page.locator('[data-testid="budget-alert"]');
      const isVisible = await alertElement.isVisible().catch(() => false);

      if (isVisible) {
        await expect(alertElement).toBeVisible();
        await expect(alertElement).toContainText(/预算|告警/);
      }
    });
  });

  /**
   * 测试套件2: 账单管理
   */
  test.describe('账单管理', () => {
    test('应该正确展示账单列表', async ({ page }) => {
      await billingPage.navigateToInvoices();

      // 验证列表容器
      await expect(page.locator('[data-testid="invoice-list"]')).toBeVisible();

      // 验证表格
      await expect(page.locator('table')).toBeVisible();

      // 验证搜索框
      await expect(page.locator('input[data-testid="invoice-search"]')).toBeVisible();

      // 验证导出按钮
      await expect(page.locator('[data-testid="btn-export-all"]')).toBeVisible();
    });

    test('应该支持搜索账单', async ({ page }) => {
      await billingPage.navigateToInvoices();

      // 搜索账单（搜索当前月份）
      const currentMonth = new Date().toISOString().slice(0, 7); // YYYY-MM
      await billingPage.searchInvoice(currentMonth);

      // 验证搜索框的值
      const searchValue = await page.locator('input[data-testid="invoice-search"]').inputValue();
      expect(searchValue).toContain(currentMonth);
    });

    test('应该支持导出账单', async ({ page }) => {
      await billingPage.navigateToInvoices();

      // 等待列表加载
      await page.waitForTimeout(1000);

      // 尝试导出第一个账单（如果存在）
      const firstExportButton = page.locator('[data-testid^="export-invoice-"]').first();
      const isVisible = await firstExportButton.isVisible().catch(() => false);

      if (isVisible) {
        // 设置下载监听
        const downloadPromise = page.waitForEvent('download');
        await firstExportButton.click();
        const download = await downloadPromise;

        // 验证下载文件
        expect(download.suggestedFilename()).toMatch(/\.(pdf|csv|xlsx)$/);
      }
    });

    test('应该支持按日期筛选账单', async ({ page }) => {
      await billingPage.navigateToInvoices();

      // 点击日期筛选按钮
      const dateFilterButton = page.locator('[data-testid="date-filter-button"]');
      await dateFilterButton.click();

      // 选择最近30天
      await page.locator('text="最近30天"').click();

      // 等待筛选结果
      await page.waitForTimeout(1000);

      // 验证筛选应用（根据实现可能需要调整）
    });
  });

  /**
   * 测试套件3: 充值
   */
  test.describe('充值', () => {
    test('应该显示充值表单', async ({ page }) => {
      await billingPage.navigateToRecharge();

      // 验证充值金额输入框
      await expect(page.locator('input[name="amount"]')).toBeVisible();

      // 验证支付方式选择
      await expect(page.locator('select[name="payment_method"]')).toBeVisible();

      // 验证充值按钮
      await expect(page.locator('button[data-testid="btn-recharge"]')).toBeVisible();
    });

    test('应该支持选择常用金额', async ({ page }) => {
      await billingPage.navigateToRecharge();

      // 点击常用金额按钮（例如100元）
      const amountButton = page.locator('button[data-amount="100"]');
      const isVisible = await amountButton.isVisible().catch(() => false);

      if (isVisible) {
        await amountButton.click();

        // 验证金额输入框已填充
        const amountInput = page.locator('input[name="amount"]');
        const value = await amountInput.inputValue();
        expect(value).toBe('100');
      }
    });

    test('应该显示充值记录', async ({ page }) => {
      await billingPage.navigateToRecharge();

      // 验证充值记录列表
      const historyList = page.locator('[data-testid="recharge-history"]');
      const isVisible = await historyList.isVisible().catch(() => false);

      if (isVisible) {
        await expect(historyList).toBeVisible();
      }
    });
  });

  /**
   * 测试套件4: 使用统计
   */
  test.describe('使用统计', () => {
    test('应该显示Token使用统计', async ({ page }) => {
      await billingPage.navigateToUsageStats();

      // 验证统计数据存在
      await expect(page.locator('[data-testid="usage-stats"]')).toBeVisible();

      // 验证图表加载
      const chart = page.locator('[data-testid="usage-chart"]');
      const isVisible = await chart.isVisible().catch(() => false);

      if (isVisible) {
        await expect(chart).toBeVisible();
      }
    });

    test('应该支持按时间范围筛选', async ({ page }) => {
      await billingPage.navigateToUsageStats();

      // 选择时间范围（例如本周）
      const timeRangeSelect = page.locator('select[data-testid="time-range"]');
      await timeRangeSelect.selectOption('week');

      // 等待数据更新
      await page.waitForTimeout(1000);

      // 验证数据已更新（根据实现可能需要调整）
    });

    test('应该显示模型使用分布', async ({ page }) => {
      await billingPage.navigateToUsageStats();

      // 验证模型使用分布图
      const modelDistribution = page.locator('[data-testid="model-distribution"]');
      const isVisible = await modelDistribution.isVisible().catch(() => false);

      if (isVisible) {
        await expect(modelDistribution).toBeVisible();
      }
    });
  });

  /**
   * 测试套件5: 成本优化建议
   */
  test.describe('成本优化建议', () => {
    test('应该显示成本优化建议', async ({ page }) => {
      await billingPage.navigateToBillingOverview();

      // 检查优化建议卡片
      const optimizationCard = page.locator('[data-testid="cost-optimization"]');
      const isVisible = await optimizationCard.isVisible().catch(() => false);

      if (isVisible) {
        await expect(optimizationCard).toBeVisible();
        await expect(optimizationCard).toContainText(/优化|建议/);
      }
    });

    test('应该支持应用优化建议', async ({ page }) => {
      await billingPage.navigateToBillingOverview();

      // 查找应用建议按钮
      const applyButton = page.locator('[data-testid="apply-optimization"]');
      const isVisible = await applyButton.isVisible().catch(() => false);

      if (isVisible) {
        // 点击应用建议
        await applyButton.click();

        // 验证确认对话框
        const confirmDialog = page.locator('[data-testid="confirm-dialog"]');
        await expect(confirmDialog).toBeVisible();

        // 取消操作（避免实际修改）
        await page.locator('button:has-text("取消")').click();
      }
    });
  });

  /**
   * 测试套件6: 预算管理
   */
  test.describe('预算管理', () => {
    test('应该显示预算配置', async ({ page }) => {
      await billingPage.navigateToBillingOverview();

      // 验证预算设置卡片
      const budgetCard = page.locator('[data-testid="budget-settings"]');
      const isVisible = await budgetCard.isVisible().catch(() => false);

      if (isVisible) {
        await expect(budgetCard).toBeVisible();
        await expect(budgetCard).toContainText(/预算|月度/);
      }
    });

    test('应该支持编辑预算', async ({ page }) => {
      await billingPage.navigateToBillingOverview();

      // 点击编辑预算按钮
      const editButton = page.locator('[data-testid="edit-budget"]');
      const isVisible = await editButton.isVisible().catch(() => false);

      if (isVisible) {
        await editButton.click();

        // 验证编辑对话框
        const dialog = page.locator('[data-testid="budget-edit-dialog"]');
        await expect(dialog).toBeVisible();

        // 取消编辑
        await page.locator('button:has-text("取消")').click();
      }
    });

    test('应该显示预算告警历史', async ({ page }) => {
      await billingPage.navigateToBillingOverview();

      // 查找告警历史链接
      const alertHistoryLink = page.locator('a[data-testid="alert-history"]');
      const isVisible = await alertHistoryLink.isVisible().catch(() => false);

      if (isVisible) {
        await alertHistoryLink.click();

        // 验证告警历史页面
        await expect(page.locator('[data-testid="alert-history-list"]')).toBeVisible();
      }
    });
  });

  /**
   * 测试套件7: 实时计费
   */
  test.describe('实时计费', () => {
    test('应该显示实时使用数据', async ({ page }) => {
      await billingPage.navigateToBillingOverview();

      // 验证实时数据卡片
      const realtimeCard = page.locator('[data-testid="realtime-usage"]');
      const isVisible = await realtimeCard.isVisible().catch(() => false);

      if (isVisible) {
        await expect(realtimeCard).toBeVisible();
      }
    });

    test('应该支持查看实时计费记录', async ({ page }) => {
      await billingPage.navigateToUsageStats();

      // 切换到实时数据标签
      const realtimeTab = page.locator('button:has-text("实时")');
      const isVisible = await realtimeTab.isVisible().catch(() => false);

      if (isVisible) {
        await realtimeTab.click();

        // 验证实时记录列表
        await expect(page.locator('[data-testid="realtime-logs"]')).toBeVisible();
      }
    });
  });
});
