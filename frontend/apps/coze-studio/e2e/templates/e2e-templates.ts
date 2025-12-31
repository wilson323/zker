/**
 * 企业级E2E测试模板
 * @version 1.0.0
 * @description 基于Playwright的端到端测试标准模板
 */

import { test, expect, Page, BrowserContext } from '@playwright/test'

// =====================================================
// Page Object Model (POM)
// =====================================================

/**
 * 登录页面对象
 */
export class LoginPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto('http://localhost:8888/login')
  }

  async login(username: string, password: string) {
    await this.page.fill('input[name="username"]', username)
    await this.page.fill('input[name="password"]', password)
    await this.page.click('button[type="submit"]')
    await this.page.waitForURL('http://localhost:8888/')
  }

  async expectLoginSuccess() {
    await expect(this.page).toHaveURL(/\/.*/)
    await expect(this.page.locator('text=欢迎')).toBeVisible()
  }
}

/**
 * 租户列表页面对象
 */
export class TenantListPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto('http://localhost:8888/tenant/list')
    // 等待列表加载完成
    await this.page.waitForSelector('.tenant-list-item', { timeout: 5000 })
  }

  async clickCreateButton() {
    await this.page.click('button:has-text("创建租户")')
    await this.page.waitForSelector('h1:has-text("创建租户")')
  }

  async searchTenant(keyword: string) {
    await this.page.fill('input[placeholder="搜索租户"]', keyword)
    await this.page.waitForTimeout(300) // 等待防抖
  }

  async filterByStatus(status: string) {
    await this.page.selectOption('select[name="status"]', status)
    await this.page.waitForTimeout(300)
  }

  async clickDeleteButton(tenantName: string) {
    const row = this.page.locator(`.tenant-list-item:has-text("${tenantName}")`)
    await row.locator('button:has-text("删除")').click()
  }

  async confirmDelete() {
    await this.page.click('button:has-text("确认")')
  }

  async expectTenantVisible(tenantName: string) {
    await expect(this.page.locator(`text=${tenantName}`)).toBeVisible()
  }

  async expectTenantNotVisible(tenantName: string) {
    await expect(this.page.locator(`text=${tenantName}`)).not.toBeVisible()
  }

  async expectSuccessToast(message: string) {
    await expect(this.page.locator(`.toast:has-text("${message}")`)).toBeVisible()
  }
}

/**
 * 租户创建页面对象
 */
export class TenantCreatePage {
  constructor(private page: Page) {}

  async fillForm(data: {
    name: string
    contact: string
    phone: string
    type: string
  }) {
    await this.page.fill('input[name="name"]', data.name)
    await this.page.fill('input[name="contact"]', data.contact)
    await this.page.fill('input[name="phone"]', data.phone)
    await this.page.selectOption('select[name="type"]', data.type)
  }

  async submit() {
    await this.page.click('button:has-text("提交")')
  }

  async expectValidationErrors() {
    await expect(this.page.locator('.error-message')).toBeVisible()
  }
}

/**
 * 租户详情页面对象
 */
export class TenantDetailPage {
  constructor(private page: Page) {}

  async goto(tenantId: string) {
    await this.page.goto(`http://localhost:8888/tenant/${tenantId}`)
    await this.page.waitForSelector('.tenant-detail')
  }

  async clickEditButton() {
    await this.page.click('button:has-text("编辑")')
  }

  async expectBasicInfo(expected: {
    name: string
    contact: string
    status: string
  }) {
    await expect(this.page.locator('.tenant-name')).toHaveText(expected.name)
    await expect(this.page.locator('.tenant-contact')).toHaveText(expected.contact)
    await expect(this.page.locator('.tenant-status')).toHaveText(expected.status)
  }
}

// =====================================================
// 测试固件 (Fixtures)
// =====================================================

/**
 * 测试扩展
 */
export const test = test.extend<{
  loginPage: LoginPage
  tenantListPage: TenantListPage
  authenticatedPage: Page
}>({
  loginPage: async ({ page }, use) => {
    await use(new LoginPage(page))
  },

  tenantListPage: async ({ page }, use) => {
    await use(new TenantListPage(page))
  },

  authenticatedPage: async ({ page, context }, use) => {
    // 登录并保存session
    const loginPage = new LoginPage(page)
    await loginPage.goto()
    await loginPage.login('admin', 'admin123')

    // 保存storage state用于后续测试
    await context.storageState({ path: 'auth.json' })

    await use(page)
  },
})

// =====================================================
// E2E测试套件模板
// =====================================================

test.describe('租户管理 E2E测试', () => {
  test.beforeEach(async ({ page }) => {
    // 每个测试前登录
    const loginPage = new LoginPage(page)
    await loginPage.goto()
    await loginPage.login('admin', 'admin123')
  })

  test.afterEach(async ({ page }) => {
    // 每个测试后清理
    // 例如：删除测试数据
  })

  test('应该成功创建租户', async ({ page }) => {
    // Arrange
    const tenantListPage = new TenantListPage(page)
    const tenantCreatePage = new TenantCreatePage(page)

    // Act
    await tenantListPage.goto()
    await tenantListPage.clickCreateButton()

    await tenantCreatePage.fillForm({
      name: 'E2E测试租户',
      contact: 'e2e@example.com',
      phone: '13800138000',
      type: 'enterprise',
    })
    await tenantCreatePage.submit()

    // Assert
    await tenantListPage.expectSuccessToast('创建成功')
    await tenantListPage.expectTenantVisible('E2E测试租户')
  })

  test('应该搜索租户', async ({ page }) => {
    // Arrange
    const tenantListPage = new TenantListPage(page)
    await tenantListPage.goto()

    // Act
    await tenantListPage.searchTenant('测试')

    // Assert
    await expect(page.locator('.tenant-list-item')).toHaveCount(10)
  })

  test('应该过滤租户状态', async ({ page }) => {
    // Arrange
    const tenantListPage = new TenantListPage(page)
    await tenantListPage.goto()

    // Act
    await tenantListPage.filterByStatus('active')

    // Assert
    const items = page.locator('.tenant-list-item .tenant-status-active')
    await expect(items.first()).toBeVisible()
  })

  test('应该删除租户', async ({ page }) => {
    // Arrange
    const tenantListPage = new TenantListPage(page)
    await tenantListPage.goto()

    // Act
    await tenantListPage.clickDeleteButton('待删除租户')
    await tenantListPage.confirmDelete()

    // Assert
    await tenantListPage.expectSuccessToast('删除成功')
    await tenantListPage.expectTenantNotVisible('待删除租户')
  })

  test('应该查看租户详情', async ({ page }) => {
    // Arrange
    const tenantId = 'tenant-123'
    const tenantDetailPage = new TenantDetailPage(page)

    // Act
    await tenantDetailPage.goto(tenantId)

    // Assert
    await tenantDetailPage.expectBasicInfo({
      name: '测试租户',
      contact: 'test@example.com',
      status: 'active',
    })
  })
})

// =====================================================
// 认证流程测试
// =====================================================

test.describe('认证流程 E2E测试', () => {
  test('应该成功登录', async ({ loginPage }) => {
    // Arrange & Act
    await loginPage.goto()
    await loginPage.login('admin', 'admin123')

    // Assert
    await loginPage.expectLoginSuccess()
  })

  test('应该显示登录错误', async ({ page }) => {
    // Arrange
    await page.goto('http://localhost:8888/login')

    // Act - 使用错误密码
    await page.fill('input[name="username"]', 'admin')
    await page.fill('input[name="password"]', 'wrong_password')
    await page.click('button[type="submit"]')

    // Assert
    await expect(page.locator('.error-message')).toBeVisible()
  })

  test('应该成功登出', async ({ page }) => {
    // Arrange
    const loginPage = new LoginPage(page)
    await loginPage.goto()
    await loginPage.login('admin', 'admin123')

    // Act
    await page.click('button:has-text("登出")')

    // Assert
    await expect(page).toHaveURL('/login')
  })
})

// =====================================================
// 表单验证测试
// =====================================================

test.describe('表单验证 E2E测试', () => {
  test.beforeEach(async ({ page }) => {
    const loginPage = new LoginPage(page)
    await loginPage.goto()
    await loginPage.login('admin', 'admin123')
  })

  test('应该验证必填字段', async ({ page }) => {
    // Arrange
    await page.goto('http://localhost:8888/tenant/create')
    const createPage = new TenantCreatePage(page)

    // Act - 直接提交空表单
    await createPage.submit()

    // Assert
    await createPage.expectValidationErrors()
    await expect(page.locator('text=请输入租户名称')).toBeVisible()
  })

  test('应该验证邮箱格式', async ({ page }) => {
    // Arrange
    await page.goto('http://localhost:8888/tenant/create')
    const createPage = new TenantCreatePage(page)

    // Act - 输入无效邮箱
    await page.fill('input[name="contact"]', 'invalid-email')
    await page.blur('input[name="contact"]')

    // Assert
    await expect(page.locator('text=邮箱格式不正确')).toBeVisible()
  })

  test('应该验证手机号格式', async ({ page }) => {
    // Arrange
    await page.goto('http://localhost:8888/tenant/create')

    // Act - 输入无效手机号
    await page.fill('input[name="phone"]', '123')
    await page.blur('input[name="phone"]')

    // Assert
    await expect(page.locator('text=手机号格式不正确')).toBeVisible()
  })
})

// =====================================================
// 性能测试
// =====================================================

test.describe('性能 E2E测试', () => {
  test('列表加载性能', async ({ page }) => {
    // 记录开始时间
    const startTime = Date.now()

    // Act
    await page.goto('http://localhost:8888/tenant/list')
    await page.waitForSelector('.tenant-list-item')

    // 计算加载时间
    const loadTime = Date.now() - startTime

    // Assert - 页面应在2秒内加载完成
    expect(loadTime).toBeLessThan(2000)
  })

  test('搜索响应性能', async ({ page }) => {
    // Arrange
    await page.goto('http://localhost:8888/tenant/list')

    // Act & 记录时间
    const startTime = Date.now()
    await page.fill('input[placeholder="搜索租户"]', '测试')
    await page.waitForSelector('.tenant-list-item')

    // 计算响应时间
    const responseTime = Date.now() - startTime

    // Assert - 搜索应在500ms内响应
    expect(responseTime).toBeLessThan(500)
  })
})

// =====================================================
// 可访问性测试
// =====================================================

test.describe('可访问性 E2E测试', () => {
  test('应该支持键盘导航', async ({ page }) => {
    // Arrange
    await page.goto('http://localhost:8888/tenant/list')

    // Act - 使用Tab键导航
    await page.keyboard.press('Tab')
    await page.keyboard.press('Tab')

    // Assert
    const focusedElement = await page.evaluate(() => document.activeElement?.tagName)
    expect(focusedElement).toBe('BUTTON')
  })

  test('应该有正确的ARIA标签', async ({ page }) => {
    // Arrange & Act
    await page.goto('http://localhost:8888/tenant/list')

    // Assert
    const searchInput = page.locator('input[placeholder="搜索租户"]')
    await expect(searchInput).toHaveAttribute('aria-label', '搜索租户')
  })

  test('应该支持屏幕阅读器', async ({ page }) => {
    // Arrange & Act
    await page.goto('http://localhost:8888/tenant/list')

    // Assert - 检查语义化HTML
    const main = page.locator('main')
    await expect(main).toBeVisible()

    const heading = page.locator('h1')
    await expect(heading).toBeVisible()
  })
})

// =====================================================
// 跨浏览器测试
// =====================================================

test.describe('跨浏览器 E2E测试', () => {
  test('应该兼容Chrome', async ({ page, browserName }) => {
    test.skip(browserName !== 'chromium')

    await page.goto('http://localhost:8888/tenant/list')
    await expect(page).toHaveTitle(/租户管理/)
  })

  test('应该兼容Firefox', async ({ page, browserName }) => {
    test.skip(browserName !== 'firefox')

    await page.goto('http://localhost:8888/tenant/list')
    await expect(page).toHaveTitle(/租户管理/)
  })

  test('应该兼容Safari', async ({ page, browserName }) => {
    test.skip(browserName !== 'webkit')

    await page.goto('http://localhost:8888/tenant/list')
    await expect(page).toHaveTitle(/租户管理/)
  })
})

// =====================================================
// 并发测试
// =====================================================

test.describe('并发 E2E测试', () => {
  test('多个用户同时创建租户', async ({ browser }) => {
    const contexts = await Promise.all(
      Array(3)
        .fill(0)
        .map(() => browser.newContext())
    )

    const pages = await Promise.all(contexts.map((ctx) => ctx.newPage()))

    // 并发创建租户
    await Promise.all(
      pages.map(async (page, index) => {
        const loginPage = new LoginPage(page)
        await loginPage.goto()
        await loginPage.login('admin', 'admin123')

        const tenantListPage = new TenantListPage(page)
        await tenantListPage.goto()
        await tenantListPage.clickCreateButton()

        const createPage = new TenantCreatePage(page)
        await createPage.fillForm({
          name: `并发测试租户${index}`,
          contact: `concurrent${index}@example.com`,
          phone: '13800138000',
          type: 'enterprise',
        })
        await createPage.submit()

        await tenantListPage.expectSuccessToast('创建成功')
      })
    )

    // 清理
    await Promise.all(contexts.map((ctx) => ctx.close()))
  })
})

// =====================================================
// 导出
// =====================================================

export default {
  LoginPage,
  TenantListPage,
  TenantCreatePage,
  TenantDetailPage,
}
