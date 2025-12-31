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
 * 租户管理 E2E 测试辅助函数
 *
 * @description 提供租户管理相关的测试辅助方法
 *
 * @author 研发C
 * @date 2025-01-04
 */

import { Page, expect } from '@playwright/test';

/**
 * 租户测试数据
 */
export interface TenantTestData {
  tenantName: string;
  contactEmail: string;
  contactPhone?: string;
  companyName?: string;
  description?: string;
}

/**
 * 租户管理页面辅助类
 *
 * 设计原则：
 * - SOLID: 单一职责，只负责租户管理页面的操作
 * - DRY: 封装可复用的操作
 * - KISS: 保持简单明了
 */
export class TenantManagementPage {
  constructor(private readonly page: Page) {}

  /**
   * 导航到租户管理页面
   */
  async navigateToTenantManagement(): Promise<void> {
    await this.page.goto('/tenant');
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForSelector('[data-testid="tenant-list"]', { state: 'visible' });
  }

  /**
   * 点击"新建租户"按钮
   */
  async clickCreateTenant(): Promise<void> {
    await this.page.click('[data-testid="btn-create-tenant"]');
    await this.page.waitForSelector('[data-testid="tenant-create-modal"]', { state: 'visible' });
  }

  /**
   * 填写新建租户表单
   *
   * @param data - 租户数据
   */
  async fillCreateForm(data: TenantTestData): Promise<void> {
    // 租户名称（必填）
    await this.page.fill('input[name="tenant_name"]', data.tenantName);

    // 联系邮箱（必填）
    await this.page.fill('input[name="contact_email"]', data.contactEmail);

    // 联系电话（可选）
    if (data.contactPhone) {
      await this.page.fill('input[name="contact_phone"]', data.contactPhone);
    }

    // 公司名称（可选）
    if (data.companyName) {
      await this.page.fill('input[name="company_name"]', data.companyName);
    }

    // 描述（可选）
    if (data.description) {
      await this.page.fill('textarea[name="description"]', data.description);
    }
  }

  /**
   * 提交新建租户表单
   */
  async submitCreateForm(): Promise<void> {
    await this.page.click('[data-testid="btn-submit-create"]');

    // 等待请求完成
    await this.page.waitForSelector('[data-testid="tenant-create-modal"]', {
      state: 'hidden',
      timeout: 5000,
    });
  }

  /**
   * 创建租户（完整流程）
   *
   * @param data - 租户数据
   */
  async createTenant(data: TenantTestData): Promise<void> {
    await this.clickCreateTenant();
    await this.fillCreateForm(data);
    await this.submitCreateForm();
  }

  /**
   * 搜索租户
   *
   * @param keyword - 搜索关键词
   */
  async searchTenant(keyword: string): Promise<void> {
    await this.page.fill('input[data-testid="search-input"]', keyword);
    await this.page.press('input[data-testid="search-input"]', 'Enter');
    await this.page.waitForTimeout(500); // 等待搜索结果
  }

  /**
   * 点击"编辑"按钮
   *
   * @param tenantName - 租户名称
   */
  async clickEdit(tenantName: string): Promise<void> {
    const row = this.page.locator(`tr:has-text("${tenantName}")`);
    await row.locator('[data-testid="btn-edit"]').click();
    await this.page.waitForSelector('[data-testid="tenant-edit-modal"]', { state: 'visible' });
  }

  /**
   * 点击"查看详情"按钮
   *
   * @param tenantName - 租户名称
   */
  async clickViewDetail(tenantName: string): Promise<void> {
    const row = this.page.locator(`tr:has-text("${tenantName}")`);
    await row.locator('[data-testid="btn-view"]').click();
    await this.page.waitForSelector('[data-testid="tenant-detail-drawer"]', { state: 'visible' });
  }

  /**
   * 关闭详情抽屉
   */
  async closeDetailDrawer(): Promise<void> {
    await this.page.click('[data-testid="btn-close-drawer"]');
    await this.page.waitForSelector('[data-testid="tenant-detail-drawer"]', {
      state: 'hidden',
    });
  }

  /**
   * 点击"删除"按钮
   *
   * @param tenantName - 租户名称
   */
  async clickDelete(tenantName: string): Promise<void> {
    const row = this.page.locator(`tr:has-text("${tenantName}")`);
    await row.locator('[data-testid="btn-delete"]').click();
  }

  /**
   * 确认删除
   */
  async confirmDelete(): Promise<void> {
    await this.page.click('[data-testid="btn-confirm-delete"]');
    await this.page.waitForSelector('text=删除成功', { state: 'visible', timeout: 5000 });
  }

  /**
   * 取消删除
   */
  async cancelDelete(): Promise<void> {
    await this.page.click('[data-testid="btn-cancel-delete"]');
  }

  /**
   * 获取租户列表中的租户数量
   */
  async getTenantCount(): Promise<number> {
    const rows = await this.page.locator('[data-testid="tenant-list"] tbody tr').count();
    return rows;
  }

  /**
   * 验证租户是否存在于列表中
   *
   * @param tenantName - 租户名称
   */
  async verifyTenantExists(tenantName: string): Promise<boolean> {
    const count = await this.page.locator(`tr:has-text("${tenantName}")`).count();
    return count > 0;
  }

  /**
   * 等待 Toast 消息出现
   *
   * @param message - 消息内容
   */
  async waitForToast(message: string): Promise<void> {
    await this.page.waitForSelector(`text=${message}`, { state: 'visible', timeout: 5000 });
  }

  /**
   * 验证表单验证错误
   *
   * @param fieldName - 字段名称
   * @param errorMessage - 错误消息
   */
  async verifyValidationError(fieldName: string, errorMessage: string): Promise<void> {
    const field = this.page.locator(`input[name="${fieldName}"]`);
    const errorText = await field.evaluate((el) => {
      const parent = el.closest('.semi-form-field');
      return parent?.querySelector('.semi-form-field-msg-error')?.textContent;
    });

    expect(errorText).toContain(errorMessage);
  }

  /**
   * 翻页
   *
   * @param pageNumber - 页码
   */
  async goToPage(pageNumber: number): Promise<void> {
    await this.page.click(`button[data-testid="page-${pageNumber}"]`);
    await this.page.waitForTimeout(500);
  }

  /**
   * 获取当前页码
   */
  async getCurrentPage(): Promise<number> {
    const activePage = this.page.locator('.semi-pagination-item-active');
    const text = await activePage.textContent();
    return parseInt(text || '1', 10);
  }

  /**
   * 等待租户列表加载完成
   */
  async waitForTenantListLoaded(): Promise<void> {
    await this.page.waitForSelector('[data-testid="tenant-list"]', { state: 'visible' });
    await this.page.waitForSelector('.spin-spinning', { state: 'detached' }).catch(() => {
      // 忽略 loading 不存在的情况
    });
  }
}
