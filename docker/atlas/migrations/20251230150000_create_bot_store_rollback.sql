-- Bot商店表回滚脚本
-- 回滚时间: 2025-12-30
-- 描述: 删除Bot商店相关的表

-- 删除bot_store_items表
DROP TABLE IF EXISTS bot_store_items;

-- 删除bot_store_categories表
DROP TABLE IF EXISTS bot_store_categories;
