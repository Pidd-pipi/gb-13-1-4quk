-- campusbooks 校园二手书交易平台 MySQL 初始化脚本
-- 由 mysql 官方镜像 /docker-entrypoint-initdb.d 首次启动时自动执行（幂等）。
-- 数据表由后端 GORM AutoMigrate 在启动时自动创建，此处只保证库与账号存在。

CREATE DATABASE IF NOT EXISTS campusbooks_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'campusbooks_user'@'%' IDENTIFIED BY 'campusbooks_pwd';
GRANT ALL PRIVILEGES ON campusbooks_db.* TO 'campusbooks_user'@'%';
FLUSH PRIVILEGES;
