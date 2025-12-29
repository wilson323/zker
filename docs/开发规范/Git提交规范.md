# ZKER Git提交规范

> **文档类型**: 开发规范
> **版本**: v1.0
> **生效日期**: 2025-01-01
> **强制级别**: 🔴 强制执行

---

## 提交信息格式

**格式**: `<type>(<scope>): <subject>`

**示例**:
```bash
feat(bot): 添加Bot创建功能
fix(user): 修复Token过期问题
docs(api): 更新API文档
```

---

## Type类型

| Type | 说明 | 示例 |
|------|------|------|
| `feat` | 新功能 | `feat(tenant): 添加租户注册功能` |
| `fix` | Bug修复 | `fix(auth): 修复登录Token验证失败` |
| `docs` | 文档变更 | `docs(readme): 更新README` |
| `style` | 代码格式 | `style(go): 格式化代码` |
| `refactor` | 重构 | `refactor(saga): 重构Saga协调器` |
| `perf` | 性能优化 | `perf(cache): 优化Redis缓存性能` |
| `test` | 测试 | `test(bot): 添加Bot服务单元测试` |
| `chore` | 构建/工具 | `chore(deps): 升级依赖版本` |

---

## Scope范围

常用Scope:
- `bot` - Bot模块
- `user` - 用户模块
- `tenant` - 租户模块
- `conversation` - 对话模块
- `knowledge` - 知识库模块
- `auth` - 认证授权
- `api` - API层
- `frontend` - 前端
- `infra` - 基础设施

---

## Subject主题

- 使用中文简述
- 不超过50字符
- 首字母小写
- 不以句号结尾

**✅ 正确示例**:
```bash
feat(bot): 添加Bot创建功能
fix(auth): 修复Token过期问题
docs(api): 更新API文档
```

**❌ 错误示例**:
```bash
添加Bot创建功能  # 缺少type和scope
feat(bot): 添加Bot创建功能。  # 不应有句号
Feat(Bot): 添加Bot创建功能  # type不应大写
```

---

## Body正文（可选）

详细描述修改内容，包括：

- **为什么**修改
- **修改了什么**
- **相关Issue**

**示例**:
```bash
feat(bot): 添加Bot创建功能

- 实现Bot基础CRUD
- 添加Bot名称唯一性验证
- 集成知识库自动创建

Closes #123
```

---

## 配置commitlint

**安装**:
```bash
npm install -D @commitlint/cli @commitlint/config-conventional husky
```

**配置文件** `.commitlintrc.js`:
```javascript
module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [2, 'always', [
      'feat', 'fix', 'docs', 'style', 'refactor',
      'perf', 'test', 'chore'
    ]],
    'type-case': [2, 'always', 'lower-case'],
    'type-empty': [2, 'never'],
    'scope-empty': [2, 'never'],
    'subject-case': [2, 'always', 'lower-case'],
    'subject-full-stop': [2, 'never', '.'],
    'subject-max-length': [2, 'always', 50],
  },
};
```

**Git Hook** `package.json`:
```json
{
  "husky": {
    "hooks": {
      "commit-msg": "commitlint -E HUSKY_GIT_PARAMS"
    }
  }
}
```

---

## 提交示例

### ✅ 正确示例

```bash
# 新功能
git commit -m "feat(bot): 添加Bot创建功能"

# Bug修复
git commit -m "fix(auth): 修复Token过期验证失败"

# 文档更新
git commit -m "docs(api): 更新Bot API文档"

# 性能优化
git commit -m "perf(cache): 优化Redis缓存性能"

# 重构
git commit -m "refactor(saga): 重构Saga协调器代码"

# 测试
git commit -m "test(bot): 添加Bot服务单元测试"
```

### ❌ 错误示例

```bash
# 缺少type
git commit -m "添加Bot创建功能"

# type不正确
git commit -m "feature(bot): 添加Bot创建功能"

# subject太长
git commit -m "feat(bot): 实现了一个非常复杂的Bot创建功能包括很多特性"

# subject以句号结尾
git commit -m "feat(bot): 添加Bot创建功能。"
```

---

**文档版本**: v1.0
**最后更新**: 2025-12-30
