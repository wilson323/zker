-- MySQL主从复制初始化脚本

-- 创建复制用户
CREATE USER IF NOT EXISTS 'repl'@'%' IDENTIFIED BY 'repl_password';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';

-- 创建监控用户
CREATE USER IF NOT EXISTS 'proxysql'@'%' IDENTIFIED BY 'proxysql_password';
GRANT REPLICATION CLIENT ON *.* TO 'proxysql'@'%';
GRANT SELECT ON performance_schema.* TO 'proxysql'@'%';

-- 创建应用用户
CREATE USER IF NOT EXISTS 'zker_app'@'%' IDENTIFIED BY 'app_password';
GRANT ALL PRIVILEGES ON zker.* TO 'zker_app'@'%';

FLUSH PRIVILEGES;

-- 查看Master状态
SHOW MASTER STATUS;

-- 记录File和Position，用于配置Slave
-- 示例: mysql-bin.000001, Position 154
