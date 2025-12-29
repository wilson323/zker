# 21-用户管理_RBAC细化补充文档

**模块**: 21-用户管理扩展
**扩展内容**: RBAC权限细化
**版本**: v2.0
**日期**: 2025-01-03

---

## 新增功能

### F12 - 数据权限（5级）

1. **全部数据** - 超级管理员
2. **本部门及下级** - 部门经理
3. **本部门** - 部门成员
4. **仅本人** - 普通员工
5. **自定义** - 高级配置

```sql
CREATE TABLE data_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL,
    resource_type VARCHAR(64) NOT NULL,
    permission_scope ENUM('all', 'department_and_below', 'department', 'self', 'custom') NOT NULL,
    custom_org_ids JSON COMMENT '自定义组织ID列表',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### F13 - 字段权限（3级）

1. **visible** - 可见
2. **masked** - 脱敏（如：138****1234）
3. **hidden** - 隐藏

```sql
CREATE TABLE field_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL,
    resource_type VARCHAR(64) NOT NULL,
    field_name VARCHAR(64) NOT NULL,
    permission ENUM('visible', 'masked', 'hidden') NOT NULL,
    mask_rule VARCHAR(128) COMMENT '脱敏规则'
);
```

### F15 - 临时授权

```sql
CREATE TABLE temporary_grants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    grantor_id BIGINT NOT NULL,
    grantee_id BIGINT NOT NULL,
    permission_code VARCHAR(128) NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME COMMENT '结束时间',
    is_revoked BOOLEAN DEFAULT FALSE
);
```

---

## 实施工作量

**工期**: 3周
**成本**: ¥9万
**优先级**: P0
