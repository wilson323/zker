# backend/domain/tenant/migration/ddl/003_cleanup_migration_tables.sql
-- 清理脚本3: 删除迁移临时表
-- 执行时机: 清理双写代码后

-- 删除迁移进度表
DROP TABLE IF EXISTS `migration_progress`;

-- 删除检查点表
DROP TABLE IF EXISTS `migration_checkpoint`;

-- 删除失败写入表（可选，建议保留一段时间）
-- DROP TABLE IF EXISTS `failed_writes`;

-- 删除白名单表（可选，灰度发布完成后可删除）
-- DROP TABLE IF EXISTS `whitelist_entries`;

-- 清理完成
SELECT 'Cleanup completed' AS status;
