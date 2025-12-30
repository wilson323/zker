# RBAC 权限系统开发助手

**版本**: v3.0.0 | **更新**: 2025-01-03

协助 RBAC 权限系统开发，实现 5 级数据权限和 3 级字段权限。

---

## 🎯 使用场景

- 创建角色和权限实体
- 实现数据权限检查（ALL/DEPARTMENT/OWN/CUSTOM/NONE）
- 实现字段权限过滤（hidden/readonly/editable）
- 设计权限中间件

---

## 🔧 权限级别

### 数据权限（5 级）

| 级别 | 说明 | 权限范围 |
|------|------|---------|
| **ALL** | 全部数据 | 可访问所有租户数据（仅管理员） |
| **DEPARTMENT** | 部门数据 | 可访问本部门及下级部门数据 |
| **OWN** | 个人数据 | 仅可访问自己创建的数据 |
| **CUSTOM** | 自定义 | 根据自定义规则过滤 |
| **NONE** | 无权限 | 不可访问任何数据 |

### 字段权限（3 级）

| 级别 | 说明 |
|------|------|
| **hidden** | 隐藏，不可见 |
| **readonly** | 只读，可见但不可修改 |
| **editable** | 可编辑，完全可访问 |

---

## 🔧 代码模板

### 1. 权限检查服务

\`\`\`go
func (s *PermissionService) CheckDataPermission(
    ctx context.Context,
    userID string,
    resourceType string,
    resourceID string,
) error {
    permission, err := s.GetDataPermission(ctx, userID, resourceType)
    if err != nil {
        return err
    }

    switch permission.Level {
    case DataPermissionAll:
        return nil // 管理员，允许访问
    case DataPermissionOwn:
        return s.checkOwnership(ctx, userID, resourceID)
    case DataPermissionNone:
        return errorx.New(errno.PermissionDenied)
    }
    return nil
}
\`\`\`

### 2. 字段过滤

\`\`\`go
func (s *PermissionService) FilterFields(
    ctx context.Context,
    userID string,
    resource interface{},
) error {
    permissions, err := s.GetFieldPermissions(ctx, userID)
    if err != nil {
        return err
    }

    for _, p := range permissions {
        if p.Level == "hidden" {
            // 隐藏字段
        } else if p.Level == "readonly" {
            // 设置为只读
        }
    }
    return nil
}
\`\`\`

---

## 📋 检查清单

- [ ] 角色和权限实体已创建
- [ ] 数据权限检查已实现
- [ ] 字段权限过滤已实现
- [ ] 权限中间件已添加
- [ ] 所有敏感 API 都有权限检查

---

## 📖 相关文档

- [03-DESIGN/rbac/](../../docs/03-DESIGN/rbac/)
- [ZKER-企业级开发规范手册](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

**🎯 目标**: 确保权限系统完整、安全、灵活！
