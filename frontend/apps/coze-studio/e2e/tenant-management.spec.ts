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
 * 租户管理 E2E 测试套件
 *
 * @description 测试租户管理的完整 CRUD 流程
 *
 * @author 研发C
 * @date 2025-01-04
 */

import { test, expect } from './fixtures/base.fixture';
import { TenantManagementPage, TenantTestData } from './fixtures/tenant.fixture';

/**
 * 测试数据
 *
 * 设计原则：
 * - 每次测试使用唯一数据，避免冲突
 * - 使用时间戳确保唯一性
 */
const generateTenantData = (): TenantTestData => {
  const timestamp = Date.now();
  return {
    tenantName: `E2E测试租户_${timestamp}`,
    contactEmail: `e2e_test_${timestamp}@example.com`,
    contactPhone: '13800138000',
    companyName: 'E2E测试公司',
    description: '这是一个E2E自动化测试创建的租户',
  };
};

/**
 * 租户管理测试套件
 */
test.describe('租户管理', () => {
  let tenantPage: TenantManagementPage;
  let testData: TenantTestData;

  /**
   * 每个测试前初始化
   */
  test.beforeEach(async ({ page }) => {
    tenantPage = new TenantManagementPage(page);
    testData = generateTenantData();

    // 导航到租户管理页面
    await tenantPage.navigateToTenantManagement();
    await tenantPage.waitForTenantListLoaded();
  });

  /**
   * 测试套件1: 租户列表展示
   */
  test.describe('租户列表', () => {
    test('应该正确展示租户列表', async ({ page }) => {
      // 验证列表容器存在
      await expect(page.locator('[data-testid="tenant-list"]')).toBeVisible();

      // 验证表格存在
      await expect(page.locator('table')).toBeVisible();

      // 验证搜索框存在
      await expect(page.locator('input[data-testid="search-input"]')).toBeVisible();

      // 验证新建按钮存在
      await expect(page.locator('[data-testid="btn-create-tenant"]')).toBeVisible();
    });

    test('应该支持搜索租户', async ({ page }) => {
      // 搜索已存在的租户（假设数据中有"测试"关键词的租户）
      await tenantPage.searchTenant('测试');

      // 等待搜索结果
      await page.waitForTimeout(500);

      // 验证搜索框的值
      const searchValue = await page.locator('input[data-testid="search-input"]').inputValue();
      expect(searchValue).toBe('测试');
    });

    test('应该支持分页', async ({ page }) => {
      // 获取当前页码
      const currentPage = await tenantPage.getCurrentPage();
      expect(currentPage).toBe(1);

      // 如果有第二页，点击翻页
      const nextPageButton = page.locator('button[data-testid="page-2"]');
      const hasSecondPage = await nextPageButton.count();

      if (hasSecondPage > 0) {
        await tenantPage.goToPage(2);
        const newPage = await tenantPage.getCurrentPage();
        expect(newPage).toBe(2);
      }
    });
  });

  /**
   * 测试套件2: 创建租户
   */
  test.describe('创建租户', () => {
    test('应该成功创建租户 - 填写所有字段', async ({ page }) => {
      // 点击新建按钮
      await tenantPage.clickCreateTenant();

      // 验证新建弹窗显示
      await expect(page.locator('[data-testid="tenant-create-modal"]')).toBeVisible();

      // 填写表单
      await tenantPage.fillCreateForm(testData);

      // 提交表单
      await tenantPage.submitCreateForm();

      // 等待成功提示
      await tenantPage.waitForToast('租户创建成功');

      // 验证租户出现在列表中
      await page.waitForTimeout(1000);
      const exists = await tenantPage.verifyTenantExists(testData.tenantName);
      expect(exists).toBe(true);
    });

    test('应该成功创建租户 - 只填必填字段', async ({ page }) => {
      // 只提供必填字段
      const minimalData: TenantTestData = {
        tenantName: testData.tenantName,
        contactEmail: testData.contactEmail,
      };

      // 创建租户
      await tenantPage.createTenant(minimalData);

      // 等待成功提示
      await tenantPage.waitForToast('租户创建成功');

      // 验证租户出现在列表中
      const exists = await tenantPage.verifyTenantExists(testData.tenantName);
      expect(exists).toBe(true);
    });

    test('应该验证必填字段', async ({ page }) => {
      // 点击新建按钮
      await tenantPage.clickCreateTenant();

      // 直接提交（不填写任何字段）
      await page.click('[data-testid="btn-submit-create"]');

      // 验证错误提示
      await tenantPage.verifyValidationError('tenant_name', '请输入租户名称');
      await tenantPage.verifyValidationError('contact_email', '请输入联系邮箱');
    });

    test('应该验证邮箱格式', async ({ page }) => {
      // 点击新建按钮
      await tenantPage.clickCreateTenant();

      // 填写无效邮箱
      await page.fill('input[name="tenant_name"]', testData.tenantName);
      await page.fill('input[name="contact_email"]', 'invalid-email');

      // 提交表单
      await page.click('[data-testid="btn-submit-create"]');

      // 验证邮箱格式错误
      await tenantPage.verifyValidationError('contact_email', '请输入有效的邮箱地址');
    });

    test('应该验证租户名称长度', async ({ page }) => {
      // 点击新建按钮
      await tenantPage.clickCreateTenant();

      // 填写过长的租户名称（> 100字符）
      const longName = 'A'.repeat(101);
      await page.fill('input[name="tenant_name"]', longName);
      await page.fill('input[name="contact_email"]', testData.contactEmail);

      // 提交表单
      await page.click('[data-testid="btn-submit-create"]');

      // 验证长度错误
      await tenantPage.verifyValidationError('tenant_name', '租户名称长度应为2-100个字符');
    });

    test('应该可以取消创建', async ({ page }) => {
      // 点击新建按钮
      await tenantPage.clickCreateTenant();

      // 填写表单
      await tenantPage.fillCreateForm(testData);

      // 点击取消按钮
      await page.click('[data-testid="btn-cancel-create"]');

      // 验证弹窗关闭
      await expect(page.locator('[data-testid="tenant-create-modal"]')).toBeHidden();

      // 验证租户未创建
      const exists = await tenantPage.verifyTenantExists(testData.tenantName);
      expect(exists).toBe(false);
    });
  });

  /**
   * 测试套件3: 查看租户详情
   */
  test.describe('查看租户详情', () => {
    test('应该正确展示租户详情', async ({ page }) => {
      // 假设列表中有"测试租户001"，点击查看详情
      const tenantName = '测试租户001';
      const hasTenant = await tenantPage.verifyTenantExists(tenantName);

      if (!hasTenant) {
        test.skip(true, '测试数据不存在，跳过此测试');
      }

      // 点击查看详情
      await tenantPage.clickViewDetail(tenantName);

      // 验证详情抽屉显示
      await expect(page.locator('[data-testid="tenant-detail-drawer"]')).toBeVisible();

      // 验证基本信息显示
      await expect(page.locator('text=租户名称')).toBeVisible();
      await expect(page.locator('text=联系邮箱')).toBeVisible();
      await expect(page.locator('text=状态')).toBeVisible();

      // 验证操作按钮存在
      await expect(page.locator('[data-testid="btn-edit"]')).toBeVisible();

      // 关闭详情抽屉
      await tenantPage.closeDetailDrawer();
    });

    test('应该支持从详情页编辑租户', async ({ page }) => {
      const tenantName = '测试租户001';
      const hasTenant = await tenantPage.verifyTenantExists(tenantName);

      if (!hasTenant) {
        test.skip(true, '测试数据不存在，跳过此测试');
      }

      // 打开详情
      await tenantPage.clickViewDetail(tenantName);

      // 点击编辑按钮
      await page.click('[data-testid="btn-edit"]');

      // 验证编辑弹窗显示
      await expect(page.locator('[data-testid="tenant-edit-modal"]')).toBeVisible();

      // 关闭弹窗
      await page.click('[data-testid="btn-cancel-edit"]');
    });

    test('应该支持刷新详情', async ({ page }) => {
      const tenantName = '测试租户001';
      const hasTenant = await tenantPage.verifyTenantExists(tenantName);

      if (!hasTenant) {
        test.skip(true, '测试数据不存在，跳过此测试');
      }

      // 打开详情
      await tenantPage.clickViewDetail(tenantName);

      // 点击刷新按钮
      await page.click('[data-testid="btn-refresh"]');

      // 等待刷新完成
      await page.waitForTimeout(1000);

      // 验证抽屉仍然打开
      await expect(page.locator('[data-testid="tenant-detail-drawer"]')).toBeVisible();

      // 关闭详情抽屉
      await tenantPage.closeDetailDrawer();
    });
  });

  /**
   * 测试套件4: 编辑租户
   */
  test.describe('编辑租户', () => {
    test('应该成功编辑租户', async ({ page }) => {
      const tenantName = '测试租户001';
      const hasTenant = await tenantPage.verifyTenantExists(tenantName);

      if (!hasTenant) {
        test.skip(true, '测试数据不存在，跳过此测试');
      }

      // 点击编辑
      await tenantPage.clickEdit(tenantName);

      // 验证编辑弹窗显示
      await expect(page.locator('[data-testid="tenant-edit-modal"]')).toBeVisible();

      // 验证数据回显
      const tenantNameValue = await page.locator('input[name="tenant_name"]').inputValue();
      expect(tenantNameValue).toBe(tenantName);

      // 修改描述
      const newDescription = `更新于 ${new Date().toLocaleString()}`;
      await page.fill('textarea[name="description"]', newDescription);

      // 提交修改
      await page.click('[data-testid="btn-submit-edit"]');

      // 等待成功提示
      await tenantPage.waitForToast('租户更新成功');

      // 验证修改生效（打开详情查看）
      await tenantPage.clickViewDetail(tenantName);
      await expect(page.locator(`text=${newDescription}`)).toBeVisible();
    });

    test('应该可以取消编辑', async ({ page }) => {
      const tenantName = '测试租户001';
      const hasTenant = await tenantPage.verifyTenantExists(tenantName);

      if (!hasTenant) {
        test.skip(true, '测试数据不存在，跳过此测试');
      }

      // 点击编辑
      await tenantPage.clickEdit(tenantName);

      // 修改字段
      await page.fill('textarea[name="description"]', '这是一个测试修改');

      // 取消编辑
      await page.click('[data-testid="btn-cancel-edit"]');

      // 验证弹窗关闭
      await expect(page.locator('[data-testid="tenant-edit-modal"]')).toBeHidden();
    });
  });

  /**
   * 测试套件5: 删除租户
   */
  test.describe('删除租户', () => {
    test('应该成功删除租户', async ({ page }) => {
      // 先创建一个测试租户
      await tenantPage.createTenant(testData);
      await page.waitForTimeout(1000);

      // 验证租户存在
      let exists = await tenantPage.verifyTenantExists(testData.tenantName);
      expect(exists).toBe(true);

      // 点击删除
      await tenantPage.clickDelete(testData.tenantName);

      // 验证确认对话框显示
      await expect(page.locator('[data-testid="delete-confirm-modal"]')).toBeVisible();

      // 确认删除
      await tenantPage.confirmDelete();

      // 验证租户已删除
      await page.waitForTimeout(1000);
      exists = await tenantPage.verifyTenantExists(testData.tenantName);
      expect(exists).toBe(false);
    });

    test('应该可以取消删除', async ({ page }) => {
      // 先创建一个测试租户
      await tenantPage.createTenant(testData);
      await page.waitForTimeout(1000);

      // 点击删除
      await tenantPage.clickDelete(testData.tenantName);

      // 取消删除
      await tenantPage.cancelDelete();

      // 验证确认对话框关闭
      await expect(page.locator('[data-testid="delete-confirm-modal"]')).toBeHidden();

      // 验证租户仍然存在
      const exists = await tenantPage.verifyTenantExists(testData.tenantName);
      expect(exists).toBe(true);
    });
  });

  /**
   * 测试套件6: 综合场景
   */
  test.describe('综合场景测试', () => {
    test('完整的租户管理流程', async ({ page }) => {
      // 1. 创建租户
      await tenantPage.createTenant(testData);
      await tenantPage.waitForToast('租户创建成功');

      // 2. 搜索租户
      await tenantPage.searchTenant(testData.tenantName);
      const exists = await tenantPage.verifyTenantExists(testData.tenantName);
      expect(exists).toBe(true);

      // 3. 查看详情
      await tenantPage.clickViewDetail(testData.tenantName);
      await expect(page.locator('[data-testid="tenant-detail-drawer"]')).toBeVisible();

      // 4. 关闭详情
      await tenantPage.closeDetailDrawer();

      // 5. 编辑租户
      await tenantPage.clickEdit(testData.tenantName);
      const newDescription = '更新后的描述';
      await page.fill('textarea[name="description"]', newDescription);
      await page.click('[data-testid="btn-submit-edit"]');
      await tenantPage.waitForToast('租户更新成功');

      // 6. 删除租户
      await tenantPage.clickDelete(testData.tenantName);
      await tenantPage.confirmDelete();

      // 7. 验证删除
      await page.waitForTimeout(1000);
      const deleted = await tenantPage.verifyTenantExists(testData.tenantName);
      expect(deleted).toBe(false);
    });
  });
});
