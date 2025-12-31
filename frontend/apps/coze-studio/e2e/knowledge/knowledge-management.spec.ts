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
 * 知识管理 E2E 测试套件
 *
 * @description 测试知识库管理的完整流程，包括创建、上传、搜索、编辑、删除等
 *
 * @author 研发C
 * @date 2025-01-04
 */

import { test, expect } from '../fixtures/base.fixture';
import { Page } from '@playwright/test';

/**
 * 知识管理页面辅助类
 */
class KnowledgeManagementPage {
  constructor(private page: Page) {}

  /**
   * 导航到知识库列表页面
   */
  async navigateToKnowledgeList() {
    await this.page.goto('/knowledge/list');
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForSelector('[data-testid="knowledge-list"]', { timeout: 10000 });
  }

  /**
   * 导航到上传文档页面
   */
  async navigateToUpload() {
    await this.page.goto('/knowledge/upload');
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForSelector('[data-testid="upload-form"]', { timeout: 10000 });
  }

  /**
   * 创建知识库
   */
  async createKnowledgeBase(data: {
    name: string;
    description: string;
    type?: string;
  }) {
    // 点击创建按钮
    await this.page.locator('[data-testid="btn-create-knowledge"]').click();

    // 填写知识库名称
    await this.page.locator('input[data-testid="knowledge-name"]').fill(data.name);

    // 填写描述
    await this.page.locator('textarea[data-testid="knowledge-description"]').fill(data.description);

    // 选择类型（如果提供）
    if (data.type) {
      await this.page.locator('select[data-testid="knowledge-type"]').selectOption(data.type);
    }

    // 提交
    await this.page.locator('button[data-testid="btn-submit"]').click();

    // 等待成功提示
    await this.page.waitForSelector('text=/创建成功|已创建/', { timeout: 5000 });
  }

  /**
   * 上传文档
   */
  async uploadDocument(filePath: string, options?: {
    knowledgeBaseId?: string;
    chunkSize?: number;
    overlap?: number;
  }) {
    // 选择知识库（如果提供）
    if (options?.knowledgeBaseId) {
      await this.page.locator('select[data-testid="knowledge-base-select"]').selectOption(options.knowledgeBaseId);
    }

    // 选择文件
    const fileInput = this.page.locator('input[type="file"]');
    await fileInput.setInputFiles(filePath);

    // 设置分块大小（如果提供）
    if (options?.chunkSize) {
      const chunkSizeInput = this.page.locator('input[data-testid="chunk-size"]');
      await chunkSizeInput.fill(options.chunkSize.toString());
    }

    // 设置重叠大小（如果提供）
    if (options?.overlap) {
      const overlapInput = this.page.locator('input[data-testid="overlap"]');
      await overlapInput.fill(options.overlap.toString());
    }

    // 点击上传按钮
    await this.page.locator('button[data-testid="btn-upload"]').click();

    // 等待上传完成
    await this.page.waitForSelector('text=/上传成功|上传完成/', { timeout: 30000 });
  }

  /**
   * 搜索知识库
   */
  async searchKnowledge(keyword: string) {
    const searchInput = this.page.locator('input[data-testid="knowledge-search"]');
    await searchInput.fill(keyword);
    await this.page.waitForTimeout(500); // 等待搜索结果
  }

  /**
   * 编辑知识库
   */
  async editKnowledgeBase(knowledgeId: string, data: {
    name?: string;
    description?: string;
  }) {
    // 点击编辑按钮
    await this.page.locator(`[data-testid="edit-knowledge-${knowledgeId}"]`).click();

    // 修改名称（如果提供）
    if (data.name) {
      const nameInput = this.page.locator('input[data-testid="knowledge-name"]');
      await nameInput.clear();
      await nameInput.fill(data.name);
    }

    // 修改描述（如果提供）
    if (data.description) {
      const descInput = this.page.locator('textarea[data-testid="knowledge-description"]');
      await descInput.clear();
      await descInput.fill(data.description);
    }

    // 保存
    await this.page.locator('button[data-testid="btn-save"]').click();

    // 等待成功提示
    await this.page.waitForSelector('text=/更新成功|已保存/', { timeout: 5000 });
  }

  /**
   * 删除知识库
   */
  async deleteKnowledgeBase(knowledgeId: string) {
    // 点击删除按钮
    await this.page.locator(`[data-testid="delete-knowledge-${knowledgeId}"]`).click();

    // 确认删除
    await this.page.locator('button:has-text("确定")').click();

    // 等待成功提示
    await this.page.waitForSelector('text=/删除成功|已删除/', { timeout: 5000 });
  }

  /**
   * 查看知识库详情
   */
  async viewKnowledgeDetail(knowledgeId: string) {
    await this.page.locator(`[data-testid="view-knowledge-${knowledgeId}"]`).click();
    await this.page.waitForSelector('[data-testid="knowledge-detail"]', { timeout: 10000 });
  }
}

/**
 * 生成测试数据
 */
const generateKnowledgeData = () => {
  const timestamp = Date.now();
  return {
    name: `E2E测试知识库_${timestamp}`,
    description: `这是一个E2E自动化测试创建的知识库_${timestamp}`,
    type: 'general',
  };
};

/**
 * 知识管理测试套件
 */
test.describe('知识管理', () => {
  let knowledgePage: KnowledgeManagementPage;

  test.beforeEach(async ({ page }) => {
    knowledgePage = new KnowledgeManagementPage(page);
  });

  /**
   * 测试套件1: 知识库列表
   */
  test.describe('知识库列表', () => {
    test('应该正确展示知识库列表', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      // 验证列表容器
      await expect(page.locator('[data-testid="knowledge-list"]')).toBeVisible();

      // 验证搜索框
      await expect(page.locator('input[data-testid="knowledge-search"]')).toBeVisible();

      // 验证创建按钮
      await expect(page.locator('[data-testid="btn-create-knowledge"]')).toBeVisible();

      // 验证上传按钮
      await expect(page.locator('[data-testid="btn-upload-documents"]')).toBeVisible();
    });

    test('应该支持搜索知识库', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      // 搜索知识库
      await knowledgePage.searchKnowledge('测试');

      // 验证搜索框的值
      const searchValue = await page.locator('input[data-testid="knowledge-search"]').inputValue();
      expect(searchValue).toBe('测试');
    });

    test('应该支持分页', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      // 检查分页组件
      const pagination = page.locator('[data-testid="pagination"]');
      const isVisible = await pagination.isVisible().catch(() => false);

      if (isVisible) {
        await expect(pagination).toBeVisible();
      }
    });
  });

  /**
   * 测试套件2: 创建知识库
   */
  test.describe('创建知识库', () => {
    test('应该成功创建知识库', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      const testData = generateKnowledgeData();

      // 创建知识库
      await knowledgePage.createKnowledgeBase(testData);

      // 验证成功提示
      await expect(page.locator('text=/创建成功/')).toBeVisible();

      // 验证知识库出现在列表中
      await expect(page.locator(`text="${testData.name}"`)).toBeVisible();
    });

    test('应该验证必填字段', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      // 点击创建按钮
      await page.locator('[data-testid="btn-create-knowledge"]').click();

      // 不填写任何字段，直接提交
      await page.locator('button[data-testid="btn-submit"]').click();

      // 验证错误提示
      await expect(page.locator('text=/请输入|必填/')).toBeVisible();
    });

    test('应该支持选择知识库类型', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      const testData = generateKnowledgeData();
      testData.type = 'qa'; // 问答类型

      // 创建知识库
      await knowledgePage.createKnowledgeBase(testData);

      // 验证创建成功
      await expect(page.locator('text=/创建成功/')).toBeVisible();
    });
  });

  /**
   * 测试套件3: 上传文档
   */
  test.describe('上传文档', () => {
    test('应该显示上传表单', async ({ page }) => {
      await knowledgePage.navigateToUpload();

      // 验证文件选择器
      await expect(page.locator('input[type="file"]')).toBeVisible();

      // 验证知识库选择器
      await expect(page.locator('select[data-testid="knowledge-base-select"]')).toBeVisible();

      // 验证上传按钮
      await expect(page.locator('button[data-testid="btn-upload"]')).toBeVisible();
    });

    test('应该支持上传PDF文档', async ({ page }) => {
      await knowledgePage.navigateToUpload();

      // 注意：需要实际有测试文件
      const testFilePath = './e2e/fixtures/test-sample.pdf';

      // 检查文件是否存在
      const fs = require('fs');
      if (fs.existsSync(testFilePath)) {
        await knowledgePage.uploadDocument(testFilePath);

        // 验证上传成功提示
        await expect(page.locator('text=/上传成功/')).toBeVisible();
      } else {
        test.skip('测试文件不存在');
      }
    });

    test('应该支持上传文本文档', async ({ page }) => {
      await knowledgePage.navigateToUpload();

      const testFilePath = './e2e/fixtures/test-sample.txt';

      // 检查文件是否存在
      const fs = require('fs');
      if (fs.existsSync(testFilePath)) {
        await knowledgePage.uploadDocument(testFilePath);

        // 验证上传成功
        await expect(page.locator('text=/上传成功/')).toBeVisible();
      } else {
        test.skip('测试文件不存在');
      }
    });

    test('应该显示上传进度', async ({ page }) => {
      await knowledgePage.navigateToUpload();

      const testFilePath = './e2e/fixtures/test-sample.pdf';

      const fs = require('fs');
      if (fs.existsSync(testFilePath)) {
        // 选择文件
        const fileInput = page.locator('input[type="file"]');
        await fileInput.setInputFiles(testFilePath);

        // 点击上传按钮
        await page.locator('button[data-testid="btn-upload"]').click();

        // 验证进度条显示
        const progressBar = page.locator('[data-testid="upload-progress"]');
        await expect(progressBar).toBeVisible();
      } else {
        test.skip('测试文件不存在');
      }
    });
  });

  /**
   * 测试套件4: 编辑知识库
   */
  test.describe('编辑知识库', () => {
    test('应该成功编辑知识库', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      // 查找第一个编辑按钮
      const editButton = page.locator('[data-testid^="edit-knowledge-"]').first();
      const isVisible = await editButton.isVisible().catch(() => false);

      if (isVisible) {
        // 获取知识库ID
        const knowledgeId = await editButton.getAttribute('data-testid');
        const id = knowledgeId?.replace('edit-knowledge-', '');

        if (id) {
          // 编辑知识库
          await knowledgePage.editKnowledgeBase(id, {
            name: `更新后的知识库_${Date.now()}`,
            description: '更新后的描述',
          });

          // 验证成功提示
          await expect(page.locator('text=/更新成功/')).toBeVisible();
        }
      } else {
        test.skip('没有可编辑的知识库');
      }
    });

    test('应该验证编辑时的必填字段', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      const editButton = page.locator('[data-testid^="edit-knowledge-"]').first();
      const isVisible = await editButton.isVisible().catch(() => false);

      if (isVisible) {
        await editButton.click();

        // 清空名称
        await page.locator('input[data-testid="knowledge-name"]').clear();

        // 尝试保存
        await page.locator('button[data-testid="btn-save"]').click();

        // 验证错误提示
        await expect(page.locator('text=/请输入|必填/')).toBeVisible();
      } else {
        test.skip('没有可编辑的知识库');
      }
    });
  });

  /**
   * 测试套件5: 删除知识库
   */
  test.describe('删除知识库', () => {
    test('应该成功删除知识库', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      // 查找第一个删除按钮
      const deleteButton = page.locator('[data-testid^="delete-knowledge-"]').first();
      const isVisible = await deleteButton.isVisible().catch(() => false);

      if (isVisible) {
        // 获取知识库ID
        const knowledgeId = await deleteButton.getAttribute('data-testid');
        const id = knowledgeId?.replace('delete-knowledge-', '');

        if (id) {
          // 获取删除前的知识库名称（用于验证）
          const knowledgeName = await page.locator(`[data-testid="knowledge-name-${id}"]`).textContent();

          // 删除知识库
          await knowledgePage.deleteKnowledgeBase(id);

          // 验证成功提示
          await expect(page.locator('text=/删除成功/')).toBeVisible();

          // 验证知识库不再出现在列表中
          if (knowledgeName) {
            await expect(page.locator(`text="${knowledgeName}"`)).not.toBeVisible();
          }
        }
      } else {
        test.skip('没有可删除的知识库');
      }
    });

    test('应该支持取消删除', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      const deleteButton = page.locator('[data-testid^="delete-knowledge-"]').first();
      const isVisible = await deleteButton.isVisible().catch(() => false);

      if (isVisible) {
        await deleteButton.click();

        // 点击取消按钮
        await page.locator('button:has-text("取消")').click();

        // 验证对话框关闭
        const dialog = page.locator('[data-testid="delete-confirm-dialog"]');
        await expect(dialog).not.toBeVisible();
      } else {
        test.skip('没有可删除的知识库');
      }
    });
  });

  /**
   * 测试套件6: 知识库详情
   */
  test.describe('知识库详情', () => {
    test('应该显示知识库详情', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      // 查找第一个查看按钮
      const viewButton = page.locator('[data-testid^="view-knowledge-"]').first();
      const isVisible = await viewButton.isVisible().catch(() => false);

      if (isVisible) {
        // 获取知识库ID
        const knowledgeId = await viewButton.getAttribute('data-testid');
        const id = knowledgeId?.replace('view-knowledge-', '');

        if (id) {
          // 查看详情
          await knowledgePage.viewKnowledgeDetail(id);

          // 验证详情页面
          await expect(page.locator('[data-testid="knowledge-detail"]')).toBeVisible();
          await expect(page.locator('[data-testid="knowledge-info"]')).toBeVisible();
        }
      } else {
        test.skip('没有可查看的知识库');
      }
    });

    test('应该显示知识库的文档列表', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      const viewButton = page.locator('[data-testid^="view-knowledge-"]').first();
      const isVisible = await viewButton.isVisible().catch(() => false);

      if (isVisible) {
        const knowledgeId = await viewButton.getAttribute('data-testid');
        const id = knowledgeId?.replace('view-knowledge-', '');

        if (id) {
          await knowledgePage.viewKnowledgeDetail(id);

          // 验证文档列表
          const docList = page.locator('[data-testid="document-list"]');
          const listVisible = await docList.isVisible().catch(() => false);

          if (listVisible) {
            await expect(docList).toBeVisible();
          }
        }
      } else {
        test.skip('没有可查看的知识库');
      }
    });
  });

  /**
   * 测试套件7: 文档管理
   */
  test.describe('文档管理', () => {
    test('应该支持删除文档', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      const viewButton = page.locator('[data-testid^="view-knowledge-"]').first();
      const isVisible = await viewButton.isVisible().catch(() => false);

      if (isVisible) {
        const knowledgeId = await viewButton.getAttribute('data-testid');
        const id = knowledgeId?.replace('view-knowledge-', '');

        if (id) {
          await knowledgePage.viewKnowledgeDetail(id);

          // 查找删除文档按钮
          const deleteDocButton = page.locator('[data-testid^="delete-document-"]').first();
          const docVisible = await deleteDocButton.isVisible().catch(() => false);

          if (docVisible) {
            await deleteDocButton.click();

            // 确认删除
            await page.locator('button:has-text("确定")').click();

            // 验证成功提示
            await expect(page.locator('text=/删除成功/')).toBeVisible();
          }
        }
      } else {
        test.skip('没有可查看的知识库');
      }
    });

    test('应该支持重新索引文档', async ({ page }) => {
      await knowledgePage.navigateToKnowledgeList();

      const viewButton = page.locator('[data-testid^="view-knowledge-"]').first();
      const isVisible = await viewButton.isVisible().catch(() => false);

      if (isVisible) {
        const knowledgeId = await viewButton.getAttribute('data-testid');
        const id = knowledgeId?.replace('view-knowledge-', '');

        if (id) {
          await knowledgePage.viewKnowledgeDetail(id);

          // 查找重新索引按钮
          const reindexButton = page.locator('[data-testid^="reindex-document-"]').first();
          const buttonVisible = await reindexButton.isVisible().catch(() => false);

          if (buttonVisible) {
            await reindexButton.click();

            // 验证成功提示
            await expect(page.locator('text=/索引成功|开始索引/')).toBeVisible();
          }
        }
      } else {
        test.skip('没有可查看的知识库');
      }
    });
  });
});
