-- ================================================================================
-- ZKER 数据库初始化和种子数据脚本
-- ================================================================================
-- 版本: v1.0
-- 生成日期: 2025-01-01
-- 说明: 包含基础数据、配置数据、示例数据
-- ================================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ================================================================================
-- 第一部分: 基础数据初始化
-- ================================================================================

-- 1.1 插入默认租户
INSERT INTO tenants (tenant_id, name, status, subscription_plan, settings, created_at) VALUES
('default', '默认租户', 'active', 'enterprise', '{"lang": "zh-CN", "timezone": "Asia/Shanghai", "currency": "CNY"}', NOW()),
('system', '系统租户', 'active', 'enterprise', '{"lang": "zh-CN", "timezone": "Asia/Shanghai"}', NOW());

-- 1.2 插入系统管理员用户
-- 密码: admin123 (实际部署时需要修改)
INSERT INTO users (user_id, tenant_id, username, email, password_hash, status, email_verified, created_at) VALUES
('admin', 'default', 'admin', 'admin@zker.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', TRUE, NOW()),
('system', 'system', 'system', 'system@zker.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', TRUE, NOW());

-- 1.3 插入系统角色
INSERT INTO roles (tenant_id, role_id, name, code, is_system, created_at) VALUES
('default', 'admin', '系统管理员', 'admin', TRUE, NOW()),
('default', 'developer', '开发者', 'developer', TRUE, NOW()),
('default', 'user', '普通用户', 'user', TRUE, NOW()),
('default', 'guest', '访客', 'guest', TRUE, NOW());

-- 1.4 插入系统权限
INSERT INTO permissions (tenant_id, parent_id, code, name, type, created_at) VALUES
-- 租户管理权限
('default', NULL, 'tenant:manage', '租户管理', 'module', NOW()),
('default', NULL, 'tenant:view', '查看租户', 'operation', NOW()),
('default', NULL, 'tenant:edit', '编辑租户', 'operation', NOW()),
('default', NULL, 'tenant:delete', '删除租户', 'operation', NOW()),

-- 用户管理权限
('default', NULL, 'user:manage', '用户管理', 'module', NOW()),
('default', NULL, 'user:view', '查看用户', 'operation', NOW()),
('default', NULL, 'user:create', '创建用户', 'operation', NOW()),
('default', NULL, 'user:edit', '编辑用户', 'operation', NOW()),
('default', NULL, 'user:delete', '删除用户', 'operation', NOW()),

-- Bot管理权限
('default', NULL, 'bot:manage', 'Bot管理', 'module', NOW()),
('default', NULL, 'bot:view', '查看Bot', 'operation', NOW()),
('default', NULL, 'bot:create', '创建Bot', 'operation', NOW()),
('default', NULL, 'bot:edit', '编辑Bot', 'operation', NOW()),
('default', NULL, 'bot:delete', '删除Bot', 'operation', NOW()),
('default', NULL, 'bot:publish', '发布Bot', 'operation', NOW()),

-- 会话管理权限
('default', NULL, 'conversation:manage', '会话管理', 'module', NOW()),
('default', NULL, 'conversation:view', '查看会话', 'operation', NOW()),
('default', NULL, 'conversation:export', '导出会话', 'operation', NOW()),
('default', NULL, 'conversation:delete', '删除会话', 'operation', NOW()),

-- 知识库管理权限
('default', NULL, 'knowledge:manage', '知识库管理', 'module', NOW()),
('default', NULL, 'knowledge:view', '查看知识库', 'operation', NOW()),
('default', NULL, 'knowledge:create', '创建知识库', 'operation', NOW()),
('default', NULL, 'knowledge:edit', '编辑知识库', 'operation', NOW()),
('default', NULL, 'knowledge:delete', '删除知识库', 'operation', NOW()),

-- 工作流管理权限
('default', NULL, 'workflow:manage', '工作流管理', 'module', NOW()),
('default', NULL, 'workflow:view', '查看工作流', 'operation', NOW()),
('default', NULL, 'workflow:create', '创建工作流', 'operation', NOW()),
('default', NULL, 'workflow:edit', '编辑工作流', 'operation', NOW()),
('default', NULL, 'workflow:execute', '执行工作流', 'operation', NOW()),

-- 系统配置权限
('default', NULL, 'system:config', '系统配置', 'module', NOW()),
('default', NULL, 'system:logs', '系统日志', 'operation', NOW()),
('default', NULL, 'system:monitor', '系统监控', 'operation', NOW());

-- 1.5 为管理员角色分配所有权限
INSERT INTO role_permissions (tenant_id, role_id, permission_id)
SELECT 'default', 'admin', id FROM permissions WHERE tenant_id = 'default';

-- 1.6 分配管理员角色给系统用户
INSERT INTO user_roles (tenant_id, user_id, role_id, granted_at) VALUES
('default', 'admin', 'admin', NOW()),
('default', 'admin', 'developer', NOW());

-- ================================================================================
-- 第二部分: 订阅方案配置
-- ================================================================================

-- 2.1 插入订阅方案
INSERT INTO subscription_plans (plan_id, plan_name, plan_type, description, price, currency, billing_cycle, quotas, features, is_active, is_public, display_order, created_at) VALUES
('free', '免费版', 'free', '适合个人用户和小团队试用', 0.00, 'CNY', 'monthly',
 '{"bots": 1, "conversations": 1000, "messages": 10000, "knowledge_bases": 1, "storage_mb": 1024}',
'["基础对话", "1个Bot", "1000次会话", "1个知识库"]', TRUE, TRUE, 1, NOW()),

('basic', '基础版', 'basic', '适合初创团队和小型企业', 99.00, 'CNY', 'monthly',
'{"bots": 5, "conversations": 10000, "messages": 100000, "knowledge_bases": 5, "storage_mb": 10240, "team_members": 5}',
'["多Bot管理", "团队协作", "优先支持", "API访问"]', TRUE, TRUE, 2, NOW()),

('professional', '专业版', 'professional', '适合中型企业和成长型团队', 299.00, 'CNY', 'monthly',
'{"bots": 20, "conversations": 100000, "messages": 1000000, "knowledge_bases": 20, "storage_mb": 102400, "team_members": 20, "custom_domain": true}',
'["无限Bot", "自定义域名", "高级分析", "专属客服"]', TRUE, TRUE, 3, NOW()),

('enterprise', '企业版', 'enterprise', '适合大型企业和定制化需求', 999.00, 'CNY', 'monthly',
'{"bots": -1, "conversations": -1, "messages": -1, "knowledge_bases": -1, "storage_mb": -1, "team_members": -1, "custom_domain": true, "sla": true}',
'["无限资源", "SLA保障", "专属支持", "定制开发"]', TRUE, TRUE, 4, NOW());

-- ================================================================================
-- 第三部分: 智能路由引擎初始化数据
-- ================================================================================

-- 3.1 插入默认意图
INSERT INTO intents (id, type, name, display_name, description, confidence_threshold, enable_llm, enable_rule, enable_vector, is_active, created_at) VALUES
('intent_greeting', 'greeting', '问候', '用户打招呼/问候', '包括你好、嗨、早上好等问候语', 0.8, TRUE, TRUE, TRUE, TRUE, NOW()),
('intent_question', 'question', '提问', '用户提问/咨询', '用户提出问题寻求答案', 0.7, TRUE, TRUE, TRUE, TRUE, NOW()),
('intent_complaint', 'complaint', '投诉', '投诉/建议', '用户表达不满或提出建议', 0.8, TRUE, TRUE, FALSE, TRUE, NOW()),
('intent_task', 'task', '任务', '任务执行', '用户请求执行特定任务', 0.7, TRUE, TRUE, TRUE, TRUE, NOW()),
('intent_farewell', 'farewell', '告别', '告别/结束', '用户表示对话结束', 0.8, TRUE, TRUE, TRUE, TRUE, NOW()),
('intent_thanks', 'thanks', '感谢', '感谢/致谢', '用户表达感谢', 0.8, TRUE, TRUE, TRUE, TRUE, NOW());

-- 3.2 插入意图训练样本
INSERT INTO intent_samples (intent_type, sample_text, language, source, quality_score, is_active, created_at) VALUES
-- 问候样本
('greeting', '你好', 'zh', 'manual', 1.0, TRUE, NOW()),
('greeting', '嗨', 'zh', 'manual', 1.0, TRUE, NOW()),
('greeting', '早上好', 'zh', 'manual', 1.0, TRUE, NOW()),
('greeting', 'hello', 'en', 'manual', 1.0, TRUE, NOW()),
('greeting', 'hi', 'en', 'manual', 1.0, TRUE, NOW()),

-- 提问样本
('question', '请问怎么使用这个功能？', 'zh', 'manual', 1.0, TRUE, NOW()),
('question', '如何创建一个Bot？', 'zh', 'manual', 1.0, TRUE, NOW()),
('question', '价格是多少？', 'zh', 'manual', 1.0, TRUE, NOW()),
('question', 'What is the price?', 'en', 'manual', 1.0, TRUE, NOW()),

-- 投诉样本
('complaint', '我不满意这个服务', 'zh', 'manual', 1.0, TRUE, NOW()),
('complaint', '你们的系统太慢了', 'zh', 'manual', 1.0, TRUE, NOW()),
('complaint', '我要投诉', 'zh', 'manual', 1.0, TRUE, NOW()),

-- 任务样本
('task', '帮我查一下天气', 'zh', 'manual', 1.0, TRUE, NOW()),
('task', '创建一个新Bot', 'zh', 'manual', 1.0, TRUE, NOW()),
('task', '导出我的数据', 'zh', 'manual', 1.0, TRUE, NOW()),

-- 告别样本
('farewell', '再见', 'zh', 'manual', 1.0, TRUE, NOW()),
('farewell', '拜拜', 'zh', 'manual', 1.0, TRUE, NOW()),
('farewell', 'goodbye', 'en', 'manual', 1.0, TRUE, NOW()),

-- 感谢样本
('thanks', '谢谢', 'zh', 'manual', 1.0, TRUE, NOW()),
('thanks', '非常感谢', 'zh', 'manual', 1.0, TRUE, NOW()),
('thanks', 'thanks', 'en', 'manual', 1.0, TRUE, NOW()),
('thanks', 'thank you', 'en', 'manual', 1.0, TRUE, NOW());

-- 3.3 插入示例服务候选者
INSERT INTO service_candidates (tenant_id, service_id, service_name, service_type, service_config, service_tags, health_status, capacity_limit, current_load, avg_response_time_ms, is_available, created_at) VALUES
('default', 'bot-customer-service', '客服Bot', 'bot',
'{"model": "gpt-4", "temperature": 0.7, "max_tokens": 2000}',
'["customer-service", "priority-high"]', 'healthy', 100, 0, 500, TRUE, NOW()),

('default', 'bot-sales-assistant', '销售助手Bot', 'bot',
'{"model": "gpt-4", "temperature": 0.8, "max_tokens": 1500}',
'["sales", "product-info"]', 'healthy', 50, 0, 600, TRUE, NOW()),

('default', 'workflow-data-analysis', '数据分析工作流', 'workflow',
'{"timeout": 300, "max_steps": 10}',
'["analysis", "data-processing"]', 'healthy', 20, 0, 2000, TRUE, NOW()),

('default', 'knowledge-base-search', '知识库搜索', 'knowledge',
'{"search_method": "hybrid", "top_k": 5}',
'["search", "knowledge"]', 'healthy', 200, 0, 200, TRUE, NOW());

-- ================================================================================
-- 第四部分: 示例组织数据
-- ================================================================================

-- 4.1 插入示例组织
INSERT INTO organizations (tenant_id, organization_id, parent_id, name, code, type, path, level, sort_order, status, created_at) VALUES
('default', 'org_root', NULL, '默认公司', 'ROOT', 'company', '/org_root', 0, 1, 'active', NOW()),
('default', 'org_tech', 'org_root', '技术部', 'TECH', 'department', '/org_root/org_tech', 1, 1, 'active', NOW()),
('default', 'org_sales', 'org_root', '销售部', 'SALES', 'department', '/org_root/org_sales', 1, 2, 'active', NOW()),
('default', 'org_hr', 'org_root', '人事部', 'HR', 'department', '/org_root/org_hr', 1, 3, 'active', NOW());

-- 4.2 插入示例岗位
INSERT INTO positions (tenant_id, organization_id, position_id, name, code, level, category, sort_order, status, created_at) VALUES
('default', 1, 'pos_ceo', 'CEO', 'CEO', 'C1', 'management', 1, 'active', NOW()),
('default', 2, 'pos_cto', 'CTO', 'CTO', 'C2', 'management', 1, 'active', NOW()),
('default', 2, 'pos_senior_dev', '高级工程师', 'SENIOR_DEV', 'L3', 'technical', 2, 'active', NOW()),
('default', 3, 'pos_sales_mgr', '销售经理', 'SALES_MGR', 'M2', 'sales', 1, 'active', NOW()),
('default', 4, 'pos_hr_mgr', '人事经理', 'HR_MGR', 'M2', 'hr', 1, 'active', NOW());

-- ================================================================================
-- 第五部分: 示例Bot数据
-- ================================================================================

-- 5.1 插入示例Bot
INSERT INTO bots (tenant_id, bot_id, bot_alias, name, type, bot_category, visibility, config, capabilities, version_id, version_number, total_conversations, total_messages, status, is_active, owner_id, created_at) VALUES
('default', 'bot_customer_service_demo', '客服Demo', '智能客服演示Bot', 'chatbot', 'customer-service', 'public',
'{"greeting": "您好，我是智能客服，有什么可以帮助您的吗？", "system_prompt": "你是一个专业的客服人员"}',
'["chat", "knowledge_search", "task_handling"]', 'v1.0', 1, 0, 0, 'published', TRUE, 1, NOW()),

('default', 'bot_sales_demo', '销售Demo', '智能销售演示Bot', 'chatbot', 'sales', 'public',
'{"greeting": "您好，我是销售助手，为您介绍我们的产品"}',
'["product_info", "recommendation", "lead_capture"]', 'v1.0', 1, 0, 0, 'published', TRUE, 1, NOW());

-- 5.2 插入示例知识库
INSERT INTO knowledge_bases (tenant_id, knowledge_base_id, name, description, type, rag_config, document_count, total_chunks, status, sync_status, owner_id, created_at) VALUES
('default', 'kb_product_manual', '产品手册', '公司所有产品的使用手册', 'product',
'{"retrieval_strategy": "hybrid", "top_k": 10, "score_threshold": 0.7}', 0, 0, 'active', 'synced', 1, NOW()),

('default', 'kb_faq', '常见问题', '客户常见问题解答', 'faq',
'{"retrieval_strategy": "vector", "top_k": 5, "score_threshold": 0.8}', 0, 0, 'active', 'synced', 1, NOW());

-- ================================================================================
-- 第六部分: 示例工作流数据
-- ================================================================================

-- 6.1 插入示例工作流
INSERT INTO workflow_meta (tenant_id, workflow_id, name, description, mode, content_type, category, tag, latest_version, total_executions, success_rate, avg_execution_time_ms, is_template, status, owner_id, created_at) VALUES
('default', 'wf_data_export', '数据导出', '导出用户数据为Excel文件', 'dsl', 'data_processing', 'data', 'export', 1, 0, 1.0000, 0, TRUE, 'published', 1, NOW()),
('default', 'wf_daily_report', '日报生成', '自动生成每日工作报告', 'dsl', 'automation', 'productivity', 'report', 1, 0, 1.0000, 0, FALSE, 'published', 1, NOW());

-- ================================================================================
-- 第七部分: 示例写作模板
-- ================================================================================

-- 7.1 插入示例写作模板
INSERT INTO writing_templates (tenant_id, template_id, template_name, template_type, description, variables, system_prompt, category, tags, is_active, is_public, owner_id, created_at) VALUES
('default', 'tpl_email_business', '商务邮件', 'email', '标准商务邮件模板',
'{"recipient": {"type": "string", "required": true}, "subject": {"type": "string", "required": true}, "content": {"type": "text", "required": true}}',
'请写一封专业的商务邮件，收件人：{{recipient}}，主题：{{subject}}', 'business', '["email", "business"]', TRUE, TRUE, 1, NOW()),

('default', 'tpl_article_summary', '文章摘要', 'summary', '生成文章摘要',
'{"article": {"type": "text", "required": true}, "length": {"type": "select", "options": ["short", "medium", "long"], "default": "medium"}}',
'请为以下文章生成一个{{length}}长度的摘要', 'content', '["summary", "article"]', TRUE, TRUE, 1, NOW()),

('default', 'tpl_social_media', '社交媒体文案', 'social', '社交媒体营销文案',
'{"product": {"type": "string", "required": true}, "platform": {"type": "select", "options": ["wechat", "weibo", "douyin"], "required": true}, "tone": {"type": "select", "options": ["professional", "casual", "humorous"], "default": "casual"}}',
'请为{{product}}写一条适合{{platform}}平台的{{tone}}风格的营销文案', 'marketing', '["social", "marketing"]', TRUE, TRUE, 1, NOW());

-- ================================================================================
-- 第八部分: 系统配置数据
-- ================================================================================

-- 8.1 插入系统配置
-- 注意：这些配置根据实际情况调整
INSERT INTO tenant_usage (tenant_id, metric_type, metric_name, usage_value, usage_unit, time_window, created_at) VALUES
('default', 'bots', 'Bot数量', 0, 'count', 'current', NOW()),
('default', 'conversations', '会话数量', 0, 'count', 'current', NOW()),
('default', 'messages', '消息数量', 0, 'count', 'current', NOW()),
('default', 'tokens', 'Token使用量', 0, 'tokens', 'current', NOW()),
('default', 'storage', '存储使用量', 0, 'MB', 'current', NOW());

-- ================================================================================
-- 第九部分: 敏感词库
-- ================================================================================

-- 9.1 插入示例敏感词
INSERT INTO sensitive_words (tenant_id, word, word_type, action, replacement, is_active, created_at) VALUES
('default', '暴力', 'illegal', 'mask', '***', TRUE, NOW()),
('default', '涉黄', 'illegal', 'block', NULL, TRUE, NOW()),
('default', '政治敏感', 'illegal', 'block', NULL, TRUE, NOW());

-- ================================================================================
-- 数据验证
-- ================================================================================

-- 验证插入的数据
SELECT '✅ 默认租户' AS item, COUNT(*) AS count FROM tenants WHERE tenant_id = 'default'
UNION ALL
SELECT '✅ 系统用户' AS item, COUNT(*) AS count FROM users WHERE user_id = 'admin'
UNION ALL
SELECT '✅ 系统角色' AS item, COUNT(*) AS count FROM roles WHERE tenant_id = 'default'
UNION ALL
SELECT '✅ 系统权限' AS item, COUNT(*) AS count FROM permissions WHERE tenant_id = 'default'
UNION ALL
SELECT '✅ 订阅方案' AS item, COUNT(*) AS count FROM subscription_plans
UNION ALL
SELECT '✅ 意图定义' AS item, COUNT(*) AS count FROM intents
UNION ALL
SELECT '✅ 意图样本' AS item, COUNT(*) AS count FROM intent_samples
UNION ALL
SELECT '✅ 服务候选者' AS item, COUNT(*) AS count FROM service_candidates
UNION ALL
SELECT '✅ 组织结构' AS item, COUNT(*) AS count FROM organizations WHERE tenant_id = 'default'
UNION ALL
SELECT '✅ 岗位信息' AS item, COUNT(*) AS count FROM positions WHERE tenant_id = 'default'
UNION ALL
SELECT '✅ 示例Bot' AS item, COUNT(*) AS count FROM bots WHERE tenant_id = 'default'
UNION ALL
SELECT '✅ 知识库' AS item, COUNT(*) AS count FROM knowledge_bases WHERE tenant_id = 'default'
UNION ALL
SELECT '✅ 工作流' AS item, COUNT(*) AS count FROM workflow_meta WHERE tenant_id = 'default'
UNION ALL
SELECT '✅ 写作模板' AS item, COUNT(*) AS count FROM writing_templates WHERE tenant_id = 'default';

SET FOREIGN_KEY_CHECKS = 1;

-- ================================================================================
-- 初始化完成提示
-- ================================================================================

SELECT '🎉 数据库初始化完成！' AS message;
SELECT '📝 默认管理员账号: admin / admin123' AS admin_info;
SELECT '⚠️  重要: 请在生产环境中修改默认密码！' AS security_warning;
