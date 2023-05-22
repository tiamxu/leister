-- 创建数据库
CREATE DATABASE jenkins DEFAULT CHARACTER SET utf8mb4;
-- 表
CREATE TABLE `item` (
	`code_id` INT ( 255 ) DEFAULT NULL COMMENT '项目id',
	`app_name` VARCHAR ( 255 ) DEFAULT NULL COMMENT '项目名',
	`app_group` VARCHAR ( 255 ) DEFAULT NULL COMMENT '项目组',
	`app_type` VARCHAR ( 255 ) DEFAULT NULL COMMENT '语言类型',
	`ssh_url_to_repo` VARCHAR ( 255 ) DEFAULT NULL COMMENT '项目地址',
	`http_url_to_repo` VARCHAR ( 255 ) DEFAULT NULL  COMMENT '项目地址'
	
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
--插入数据
INSERT INTO `jenkins`.`item` (`code_id`, `app_name`, `app_group`, `app_type`, `ssh_url_to_repo`, `http_url_to_repo`) VALUES (2, 'hello', 'server', 'go', 'git@172.168.1.10:server/hello.git', 'http://172.168.1.10/server/hello.git');