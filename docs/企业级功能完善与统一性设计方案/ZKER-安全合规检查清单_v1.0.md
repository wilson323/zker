# ZKER 企业级安全合规检查清单 v1.0

**项目名称**: ZKER 企业级功能完善
**版本**: v1.0
**更新日期**: 2025-12-30
**适用范围**: 所有开发人员、测试人员、DevOps工程师

---

## 📋 使用说明

### 检查频率

- **代码提交前**: 必须完成个人检查 (L1)
- **合并到主分支前**: 必须完成模块检查 (L2)
- **发布前**: 必须完成集成检查 (L3)
- **定期审计**: 每季度完成发布检查 (L4)

### 检查等级

```
L1: 个人检查   ← 开发者自检
L2: 模块检查   ← 模块负责人审查
L3: 集成检查   ← QA团队测试
L4: 发布检查   ← 安全团队审计
```

---

## 1. SQL注入防护检查

### 1.1 编码规范检查

- [ ] **禁止字符串拼接SQL**
  ```go
  // ❌ 错误
  sql := fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", userID)

  // ✅ 正确
  db.Where("id = ?", userID).First(&user)
  ```

- [ ] **使用参数化查询**
  ```go
  // ✅ 使用GORM
  db.Where("name = ? AND age > ?", name, age).Find(&users)

  // ✅ 使用.Raw()时也要参数化
  db.Raw("SELECT * FROM bots WHERE tenant_id = ?", tenantID).Scan(&bots)
  ```

- [ ] **输入验证**
  ```go
  // ✅ 验证输入
  if !isValidUUID(userID) {
      return errors.New("invalid user_id")
  }
  ```

### 1.2 代码审查清单

- [ ] 所有`.Raw()`调用都有安全审查注释
- [ ] 所有用户输入都经过验证
- [ ] 所有表名、字段名都使用白名单验证
- [ ] 没有直接拼接用户输入到SQL

### 1.3 测试要求

- [ ] 单元测试覆盖所有数据库查询
- [ ] 集成测试包含SQL注入攻击测试
- [ ] 渗透测试通过SQL注入检测

**验证命令**:
```bash
# 搜索潜在SQL注入
grep -rn "\.Raw(" backend/ --include="*.go"
grep -rn "fmt\.Sprintf.*SELECT" backend/ --include="*.go"

# 运行测试
cd backend && go test ./... -run TestSQLInjection -v
```

---

## 2. XSS防护检查

### 2.1 后端防护

- [ ] **输出转义**
  ```go
  // ✅ 转义HTML特殊字符
  import "html"

  escaped := html.EscapeString(userInput)
  ```

- [ ] **Content-Type头**
  ```go
  // ✅ 强制JSON响应
  ctx.SetContentType("application/json; charset=utf-8")

  // ❌ 禁止返回text/html
  ```

- [ ] **响应清理中间件**
  ```go
  // ✅ 自动清理所有响应
  r.Use(middleware.SanitizeResponseMiddleware())
  ```

### 2.2 前端防护

- [ ] **DOMPurify清理**
  ```typescript
  // ✅ 清理HTML
  import DOMPurify from 'dompurify';

  const clean = DOMPurify.sanitize(dirtyHTML);
  ```

- [ ] **React自动转义**
  ```tsx
  // ✅ React默认转义
  const message = '<script>alert("XSS")</script>';
  return <div>{message}</div>; // 自动转义

  // ⚠️ dangerouslySetInnerHTML必须先清理
  const clean = DOMPurify.sanitize(html);
  return <div dangerouslySetInnerHTML={{ __html: clean }} />;
  ```

- [ ] **避免innerHTML**
  ```typescript
  // ❌ 错误
  element.innerHTML = userInput;

  // ✅ 正确
  element.textContent = userInput;

  // 或者使用DOMPurify
  element.innerHTML = DOMPurify.sanitize(userInput);
  ```

### 2.3 内容安全策略(CSP)

- [ ] **CSP头配置**
  ```go
  // ✅ 添加CSP头
  c.SetHeader("Content-Security-Policy",
      "default-src 'self'; "+
      "script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
      "style-src 'self' 'unsafe-inline';")
  ```

- [ ] **不允许内联脚本**(除特殊情况)
  ```html
  <!-- ❌ 禁止内联脚本 -->
  <script>alert('XSS')</script>

  <!-- ✅ 使用外部脚本 -->
  <script src="/static/js/app.js"></script>
  ```

### 2.4 测试要求

- [ ] XSS攻击测试用例通过
- [ ] DOMPurify配置正确
- [ ] CSP头配置正确
- [ ] 所有用户输入都经过清理

**验证方法**:
```bash
# 使用OWASP ZAP扫描
zap-cli quick-scan --self-contained http://localhost:8888

# 手动测试
# 在Bot描述字段输入:
<script>alert(document.cookie)</script>
# 预期: 不应该弹出alert,应该看到转义后的文本
```

---

## 3. CSRF防护检查

### 3.1 后端实现

- [ ] **CSRF Token生成**
  ```go
  // ✅ 使用加密安全的随机数生成器
  token, err := security.GenerateCSRFToken(userID)
  ```

- [ ] **CSRF Token验证**
  ```go
  // ✅ 验证Token
  if err := security.ValidateCSRFToken(token, userID); err != nil {
      return err
  }
  ```

- [ ] **Token有效期**
  - [ ] Token有效期 ≤ 24小时
  - [ ] Token一次性使用
  - [ ] Token绑定用户ID

### 3.2 前端实现

- [ ] **Token获取**
  ```typescript
  // ✅ 从服务器获取Token
  const csrfToken = await getCSRFToken();
  ```

- [ ] **Token发送**
  ```typescript
  // ✅ 添加到请求头
  fetch('/api/bots', {
    method: 'POST',
    headers: {
      'X-CSRF-Token': csrfToken,
    },
  });
  ```

- [ ] **Cookie配置**
  ```go
  // ✅ 设置SameSite属性
  http.SetCookie(c, &http.Cookie{
      Name:     "session_id",
      SameSite: http.SameSiteStrictMode,
      Secure:   true,
      HttpOnly: true,
  })
  ```

### 3.3 豁免规则

- [ ] GET/HEAD/OPTIONS请求不需要CSRF Token
- [ ] 公开API(如登录、注册)不需要CSRF Token
- [ ] 其他API必须验证CSRF Token

### 3.4 测试要求

- [ ] CSRF攻击测试用例通过
- [ ] Token验证逻辑测试通过
- [ ] Token过期测试通过

**验证方法**:
```html
<!-- CSRF攻击测试页面 -->
<html>
<body>
  <form id="evil-form" action="http://localhost:8888/api/bot/123" method="POST">
    <input type="hidden" name="action" value="delete">
  </form>
  <script>
    document.getElementById('evil-form').submit();
  </script>
</body>
</html>

# 预期: 请求被拒绝,返回403 Forbidden
```

---

## 4. 权限校验检查

### 4.1 中间件检查

- [ ] **所有写操作都有权限校验**
  ```go
  // ✅ 在路由层添加权限中间件
  r.PUT("/api/bots/:bot_id",
      middleware.RequirePermission(middleware.PermissionCheckConfig{
          ResourceType: entity.ResourceTypeBots,
          Action:       "update",
          GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
      }),
      handler.UpdateBot)
  ```

- [ ] **租户隔离**
  ```go
  // ✅ 强制租户隔离
  filter.TenantID = getTenantID(ctx)
  ```

- [ ] **资源所有权验证**
  ```go
  // ✅ 验证用户是否拥有资源
  if bot.CreatorID != userID {
      return errno.ErrPermissionDenied
  }
  ```

### 4.2 Service层检查

- [ ] **双重校验**(防御深度)
  ```go
  // ✅ Service层也进行权限校验
  func (s *BotService) UpdateBot(ctx context.Context, botID string, req *UpdateBotRequest) error {
      // 1. 检查Bot所有权
      bot, err := s.repo.GetByID(ctx, botID)
      if err != nil {
          return err
      }

      // 2. 验证权限
      if !s.hasPermission(ctx, bot) {
          return errno.ErrPermissionDenied
      }

      // 3. 更新
      return s.repo.Update(ctx, bot)
  }
  ```

### 4.3 权限测试

- [ ] 正常权限测试通过
- [ ] 越权访问测试被拒绝
- [ ] 租户隔离测试通过
- [ ] 角色权限测试通过

**验证方法**:
```bash
# 测试越权访问
# 1. 用户A创建Bot
BOT_ID=$(curl -X POST http://localhost:8888/api/bots \
  -H "X-User-ID: user-a" \
  -d '{"name": "Bot A"}' | jq '.bot_id')

# 2. 用户B尝试修改用户A的Bot
curl -X PUT http://localhost:8888/api/bots/$BOT_ID \
  -H "X-User-ID: user-b" \
  -d '{"name": "Hacked"}'

# 预期: 返回403 Forbidden
```

---

## 5. 敏感数据保护检查

### 5.1 密码安全

- [ ] **密码哈希算法**
  ```go
  // ✅ 使用Argon2id
  hash := argon2.IDKey(password, salt, time, memory, threads, keyLen)
  ```

- [ ] **参数配置**
  - [ ] 内存 ≥ 64MB
  - [ ] 迭代次数 ≥ 3
  - [ ] 并行线程 ≥ 4
  - [ ] 盐值长度 ≥ 16字节

- [ ] **禁止**:
  - [ ] ❌ 明文存储密码
  - [ ] ❌ 使用MD5
  - [ ] ❌ 使用SHA-256
  - [ ] ❌ 使用可逆加密

### 5.2 数据脱敏

- [ ] **日志脱敏**
  ```go
  // ✅ 脱敏敏感数据
  log := maskingSvc.SanitizeForLog(requestData)
  logs.Infof(ctx, "request: %s", log)
  ```

- [ ] **响应脱敏**
  ```go
  // ✅ 不返回敏感字段
  type UserResponse struct {
      UserID   string `json:"user_id"`
      Username string `json:"username"`
      // ❌ 不返回Password字段
  }
  ```

- [ ] **数据库加密**(可选)
  ```sql
  -- 敏感字段使用AES加密
  ALTER TABLE users MODIFY COLUMN api_key VARBINARY(255);
  ```

### 5.3 传输安全

- [ ] **强制HTTPS**
  ```nginx
  # Nginx配置
  server {
      listen 80;
      server_name example.com;
      return 301 https://$server_name$request_uri;
  }
  ```

- [ ] **TLS配置**
  - [ ] 使用TLS 1.3
  - [ ] 禁用TLS 1.0/1.1
  - [ ] 使用强加密套件

- [ ] **证书验证**
  - [ ] 使用有效的SSL证书
  - [ ] 证书不过期
  - [ ] 启用证书固定(Certificate Pinning)

### 5.4 测试要求

- [ ] 密码哈希算法测试通过
- [ ] 数据脱敏测试通过
- [ ] 传输加密测试通过
- [ ] 敏感数据不泄露到日志

**验证方法**:
```bash
# 检查日志中是否有明文密码
grep -r "password" backend/logs/ | grep -v "*****"

# 应该没有结果,或者都是脱敏后的***
```

---

## 6. GDPR合规检查

### 6.1 数据最小化

- [ ] **只收集必要数据**
  ```go
  // ✅ 只收集必要字段
  type User struct {
      UserID   string `json:"user_id"`
      Email    string `json:"email"`
      Username string `json:"username"`
      // ❌ 不收集不必要的字段
  }
  ```

- [ ] **数据保留期限**
  - [ ] 设定数据保留期限
  - [ ] 过期数据自动删除
  - [ ] 提供数据导出功能

### 6.2 用户权利

- [ ] **被遗忘权**
  ```go
  // ✅ 支持删除用户数据
  func (s *UserService) DeleteUser(ctx context.Context, userID string) error
  ```

- [ ] **数据导出**
  ```go
  // ✅ 支持导出用户数据
  func (s *UserService) ExportUserData(ctx context.Context, userID string) error
  ```

- [ ] **数据访问**
  ```go
  // ✅ 用户可以查看自己的数据
  func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error)
  ```

### 6.3 知情同意

- [ ] **同意记录**
  ```sql
  CREATE TABLE user_consents (
      user_id VARCHAR(36),
      consent_type VARCHAR(50),
      granted BOOLEAN,
      granted_at TIMESTAMP,
      revoked_at TIMESTAMP
  );
  ```

- [ ] **同意撤回**
  ```go
  // ✅ 用户可以撤回同意
  func (s *UserService) RevokeConsent(ctx context.Context, userID, consentType string) error
  ```

### 6.4 审计日志

- [ ] **完整审计日志**
  - [ ] 记录所有数据访问
  - [ ] 记录所有数据修改
  - [ ] 记录所有数据导出
  - [ ] 日志不可篡改(SHA-256签名)

- [ ] **日志查询**
  ```go
  // ✅ 支持查询用户数据访问记录
  func (s *AuditService) GetUserLogs(ctx context.Context, userID string) ([]*AuditLog, error)
  ```

### 6.5 测试要求

- [ ] 数据最小化测试通过
- [ ] 用户权利测试通过
- [ ] 同意管理测试通过
- [ ] 审计日志测试通过

---

## 7. 审计日志检查

### 7.1 日志完整性

- [ ] **记录所有敏感操作**
  ```go
  // ✅ 记录敏感操作
  auditSvc.LogOperation(ctx, &entity.AuditLog{
      Action:   entity.AuditActionUserDelete,
      Resource: entity.AuditResourceUser,
      ResourceID: userID,
  })
  ```

- [ ] **日志字段完整**
  - [ ] user_id
  - [ ] tenant_id
  - [ ] action
  - [ ] resource
  - [ ] resource_id
  - [ ] request_ip
  - [ ] created_at
  - [ ] status

- [ ] **敏感数据脱敏**
  ```go
  // ✅ 脱敏敏感数据
  log.RequestData = maskingSvc.SanitizeForLog(requestData)
  ```

### 7.2 日志安全

- [ ] **防篡改**
  ```go
  // ✅ SHA-256签名
  log.Signature = generateSignature(log)
  ```

- [ ] **访问控制**
  - [ ] 只有管理员可以查看审计日志
  - [ ] 审计日志操作也有审计记录

- [ ] **长期保存**
  - [ ] 定期归档到对象存储
  - [ ] 归档数据加密存储

### 7.3 日志查询

- [ ] **多条件查询**
  ```go
  // ✅ 支持多条件过滤
  func (s *AuditService) QueryLogs(ctx context.Context, filter *LogFilter) ([]*AuditLog, int64, error)
  ```

- [ ] **日志导出**
  ```go
  // ✅ 支持导出为CSV/Excel/JSON
  func (s *AuditService) ExportLogs(ctx context.Context, filter *LogFilter, format string) error
  ```

### 7.4 测试要求

- [ ] 审计日志完整性测试通过
- [ ] 审计日志防篡改测试通过
- [ ] 审计日志查询测试通过
- [ ] 审计日志导出测试通过

---

## 8. 安全响应头检查

### 8.1 必需的安全头

- [ ] **X-Content-Type-Options: nosniff**
  ```go
  c.SetHeader("X-Content-Type-Options", "nosniff")
  ```

- [ ] **X-Frame-Options: DENY**
  ```go
  c.SetHeader("X-Frame-Options", "DENY")
  ```

- [ ] **X-XSS-Protection: 1; mode=block**
  ```go
  c.SetHeader("X-XSS-Protection", "1; mode=block")
  ```

- [ ] **Content-Security-Policy**
  ```go
  c.SetHeader("Content-Security-Policy",
      "default-src 'self'; "+
      "script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
      "style-src 'self' 'unsafe-inline';")
  ```

- [ ] **Strict-Transport-Security** (仅HTTPS)
  ```go
  if c.Request.URI().Scheme() == "https" {
      c.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
  }
  ```

### 8.2 推荐的安全头

- [ ] **Referrer-Policy: strict-origin-when-cross-origin**
- [ ] **Permissions-Policy** (旧称 Feature-Policy)
- [ ] **Cross-Origin-Opener-Policy**
- [ ] **Cross-Origin-Resource-Policy**

### 8.3 验证方法

```bash
# 检查响应头
curl -I https://api.zker.com/api/bots

# 预期输出包含:
# X-Content-Type-Options: nosniff
# X-Frame-Options: DENY
# X-XSS-Protection: 1; mode=block
# Content-Security-Policy: default-src 'self'; ...
```

---

## 9. 依赖安全检查

### 9.1 依赖扫描

- [ ] **定期扫描漏洞**
  ```bash
  # Go依赖漏洞扫描
  go list -json -m all | nancy sleuth

  # 或使用
  go mod vulnerability
  ```

- [ ] **及时更新依赖**
  ```bash
  # 更新依赖
  go get -u ./...
  go mod tidy
  ```

### 9.2 许可证合规

- [ ] **检查依赖许可证**
  ```bash
  go-licenses check ./...
  go-licenses save ./... --save_path=licenses
  ```

- [ ] **禁止的许可证**
  - [ ] ❌ GPL
  - [ ] ❌ AGPL
  - [ ] ✅ MIT, Apache 2.0, BSD

### 9.3 供应链安全

- [ ] **使用go.sum验证依赖**
  ```bash
  # 确保go.sum文件提交到版本控制
  git add go.sum
  ```

- [ ] **使用Go Checksum Database**
  ```bash
  # GOSUMDB默认启用
  export GOSUMDB="sum.golang.org"
  ```

---

## 10. CI/CD安全检查

### 10.1 自动化安全测试

- [ ] **集成到CI/CD**
  ```yaml
  # .github/workflows/security.yml
  name: Security Scan
  on: [push, pull_request]

  jobs:
    security:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v3

        # Go安全扫描
        - name: Run Gosec
          uses: securego/gosec@master
          with:
            args: './...'

        # 依赖漏洞扫描
        - name: Run Nancy
          uses: sonatype-nft-community/nancy-github-action@main

        # Docker镜像扫描
        - name: Run Trivy
          uses: aquasecurity/trivy-action@master
          with:
            scan-type: 'fs'
            scan-ref: '.'
            format: 'sarif'
            output: 'trivy-results.sarif'
  ```

### 10.2 代码审查检查点

- [ ] **PR模板包含安全检查**
  ```markdown
  ## 安全检查
  - [ ] 所有用户输入都经过验证
  - [ ] 所有SQL查询都使用参数化
  - [ ] 所有敏感数据都经过脱敏
  - [ ] 没有硬编码的密钥/密码
  - [ ] 权限校验正确实现
  ```

- [ ] **自动化安全测试通过**
  - [ ] gosec扫描通过
  - [ ] 依赖扫描通过
  - [ ] 单元测试覆盖率 ≥ 80%

### 10.3 部署前检查

- [ ] **环境变量检查**
  ```bash
  # 确保生产环境不使用默认密钥
  if [ "$SECRET_KEY" = "default-secret-key" ]; then
      echo "ERROR: Default secret key detected!"
      exit 1
  fi
  ```

- [ ] **配置验证**
  ```bash
  # 验证必需的环境变量
  required_vars=("DB_HOST" "DB_PASSWORD" "SECRET_KEY")
  for var in "${required_vars[@]}"; do
      if [ -z "${!var}" ]; then
          echo "ERROR: $var is not set!"
          exit 1
      fi
  done
  ```

---

## 11. 安全培训检查

### 11.1 开发培训

- [ ] **新员工入职培训**
  - [ ] 安全编码规范培训
  - [ ] OWASP Top 10培训
  - [ ] GDPR合规培训

- [ ] **定期安全培训**
  - [ ] 每季度安全意识培训
  - [ ] 最新漏洞案例分析
  - [ ] 安全工具使用培训

### 11.2 安全意识

- [ ] **钓鱼邮件演练**
  - [ ] 每季度进行钓鱼邮件测试
  - [ ] 统计点击率
  - [ ] 提供安全培训

- [ ] **密码安全**
  - [ ] 禁止密码复用
  - [ ] 强制密码复杂度
  - [ ] 定期更换密码

### 11.3 应急响应

- [ ] **应急响应计划**
  - [ ] 安全事件分类
  - [ ] 响应流程
  - [ ] 责任人
  - [ ] 通知机制

- [ ] **应急演练**
  - [ ] 每半年进行一次应急演练
  - [ ] 演练报告
  - [ ] 改进措施

---

## 12. 定期安全审计

### 12.1 内部审计

- [ ] **月度安全扫描**
  - [ ] gosec扫描
  - [ ] 依赖漏洞扫描
  - [ ] Docker镜像扫描

- [ ] **季度渗透测试**
  - [ ] 内部团队或第三方
  - [ ] 覆盖所有关键功能
  - [ ] 修复所有高危漏洞

### 12.2 外部审计

- [ ] **年度安全审计**
  - [ ] 聘请第三方安全公司
  - [ ] 全面安全评估
  - [ ] 合规性审计

- [ ] **认证审计**
  - [ ] 等保三级认证
  - [ ] GDPR合规认证
  - [ ] SOC 2认证

---

## 13. 安全评分

### 13.1 评分标准

| 维度 | 权重 | 得分 | 加权得分 |
|------|------|------|---------|
| SQL注入防护 | 15% | ___ | ___ |
| XSS防护 | 15% | ___ | ___ |
| CSRF防护 | 15% | ___ | ___ |
| 权限校验 | 15% | ___ | ___ |
| 敏感数据保护 | 10% | ___ | ___ |
| GDPR合规 | 10% | ___ | ___ |
| 审计日志 | 10% | ___ | ___ |
| 安全响应头 | 5% | ___ | ___ |
| 依赖安全 | 3% | ___ | ___ |
| CI/CD安全 | 2% | ___ | ___ |
| **总分** | **100%** | **___** | **___** |

### 13.2 评级标准

- **90-100分**: 🟢 优秀
- **80-89分**: 🟢 良好
- **70-79分**: 🟡 中等
- **60-69分**: 🟡 及格
- **<60分**: 🔴 不合格

### 13.3 达标要求

- **发布要求**: ≥ 90分
- **等保三级**: ≥ 85分
- **GDPR合规**: ≥ 90分
- **SOC 2**: ≥ 90分

---

## 14. 快速检查命令

### 14.1 后端检查

```bash
# SQL注入检查
grep -rn "\.Raw(" backend/ --include="*.go"
grep -rn "fmt\.Sprintf.*SELECT" backend/ --include="*.go"

# 密码安全检查
grep -rn "bcrypt\|argon2" backend/ --include="*.go"

# 敏感数据日志检查
grep -rn "password.*log\|token.*log" backend/ --include="*.go" -i

# 权限校验检查
grep -rn "CheckPermission\|RequirePermission" backend/api --include="*.go"

# 安全扫描
cd backend
gosec ./...
staticcheck ./...
go vet ./...
```

### 14.2 前端检查

```bash
# XSS检查
grep -rn "dangerouslySetInnerHTML" frontend/ --include="*.tsx" --include="*.jsx"

# DOMPurify使用检查
grep -rn "DOMPurify\|sanitize" frontend/ --include="*.tsx" --include="*.ts"

# CSRF Token检查
grep -rn "X-CSRF-Token\|csrf" frontend/ --include="*.ts"
```

### 14.3 配置检查

```bash
# 环境变量检查
grep -rn "SECRET_KEY\|DB_PASSWORD" backend/conf/ --include=".env*"

# Nginx配置检查
grep -rn "ssl\|https" docker/nginx/

# Docker配置检查
grep -rn "USER\|privilege" docker/
```

---

## 15. 常见安全问题

### 15.1 硬编码密钥

❌ **错误**:
```go
const API_KEY = "sk-1234567890abcdef"
```

✅ **正确**:
```go
apiKey := os.Getenv("API_KEY")
if apiKey == "" {
    return errors.New("API_KEY is required")
}
```

### 15.2 SQL注入

❌ **错误**:
```go
sql := fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", userID)
db.Raw(sql).Scan(&user)
```

✅ **正确**:
```go
db.Where("id = ?", userID).First(&user)
```

### 15.3 XSS漏洞

❌ **错误**:
```tsx
<div dangerouslySetInnerHTML={{ __html: userInput }} />
```

✅ **正确**:
```tsx
import DOMPurify from 'dompurify';

const clean = DOMPurify.sanitize(userInput);
<div dangerouslySetInnerHTML={{ __html: clean }} />
```

### 15.4 CSRF漏洞

❌ **错误**:
```go
// 没有CSRF Token验证
r.POST("/api/bots", handler.CreateBot)
```

✅ **正确**:
```go
r.POST("/api/bots",
    middleware.CSRFMiddleware(),
    handler.CreateBot)
```

---

## 16. 应急响应清单

### 16.1 安全事件分类

| 等级 | 描述 | 响应时间 |
|------|------|---------|
| 🔴 P0 | 严重漏洞(如数据泄露) | 1小时内 |
| 🟡 P1 | 高危漏洞(如SQL注入) | 4小时内 |
| 🟢 P2 | 中危漏洞(如XSS) | 24小时内 |
| 🔵 P3 | 低危问题(如配置优化) | 1周内 |

### 16.2 应急响应流程

```
1. 发现 → 2. 报告 → 3. 评估 → 4. 修复 → 5. 验证 → 6. 复盘
```

**详细步骤**:

1. **发现**: 监控告警、用户报告、安全扫描
2. **报告**: 立即报告给安全团队和CTO
3. **评估**: 评估影响范围和严重程度
4. **修复**: 制定修复方案并实施
5. **验证**: 验证修复有效性
6. **复盘**: 编写事故报告,改进流程

---

## 附录A: 安全工具清单

### A.1 静态分析工具

| 工具 | 用途 | 链接 |
|------|------|------|
| gosec | Go安全扫描 | https://github.com/securego/gosec |
| staticcheck | Go静态分析 | https://staticcheck.io/ |
| golangci-lint | Go综合检查 | https://golangci-lint.run/ |

### A.2 动态测试工具

| 工具 | 用途 | 链接 |
|------|------|------|
| OWASP ZAP | Web应用安全测试 | https://www.zaproxy.org/ |
| Burp Suite | Web应用安全测试 | https://portswigger.net/burp |

### A.3 依赖扫描工具

| 工具 | 用途 | 链接 |
|------|------|------|
| nancy | Go依赖漏洞扫描 | https://github.com/sonatype-nft-community/nancy |
| Trivy | Docker镜像扫描 | https://aquasecurity.github.io/trivy/ |

### A.4 前端安全工具

| 工具 | 用途 | 链接 |
|------|------|------|
| DOMPurify | HTML清理 | https://github.com/cure53/DOMPurify |
| eslint-plugin-security | ESLint安全规则 | https://github.com/nodesecurity/eslint-plugin-security |

---

**检查清单结束**

**版本**: v1.0
**更新日期**: 2025-12-30
**下次审查**: 2025-01-30
**维护人员**: 安全团队
