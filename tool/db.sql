# ************************************************************
# Sequel Ace SQL dump
# 版本号： 20100
#
# https://sequel-ace.com/
# https://github.com/Sequel-Ace/Sequel-Ace
#
# 主机: 127.0.0.1 (MySQL 9.3.0)
# 数据库: egosrc
# 生成时间: 2026-05-31 06:13:54 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE='NO_AUTO_VALUE_ON_ZERO', SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# 转储表 dy_online
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_online`;

CREATE TABLE `dy_online` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `uid` bigint DEFAULT NULL,
  `ip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `active` int DEFAULT NULL,
  `status` smallint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



# 转储表 dy_sys_account
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_account`;

CREATE TABLE `dy_sys_account` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `account` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '登录账号',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '登录密码',
  `truename` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '真实姓名',
  `role_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '权限组',
  `sysgroup` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `login_count` int DEFAULT NULL COMMENT '登录次数',
  `lastlogin` int DEFAULT NULL COMMENT '最后登录时间',
  `lastip` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` smallint NOT NULL DEFAULT '1' COMMENT '用户状态 1正常 2 停用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='OA用户表';

LOCK TABLES `dy_sys_account` WRITE;
/*!40000 ALTER TABLE `dy_sys_account` DISABLE KEYS */;

INSERT INTO `dy_sys_account` (`id`, `account`, `nickname`, `password`, `truename`, `role_id`, `sysgroup`, `login_count`, `lastlogin`, `lastip`, `status`)
VALUES
	(1,'ego','','$2a$10$81LK02ZdAI9Q73OKRkiZl.fQd9aMTO0655or/SWQmpXWyV9e4fB8.','管理员','1,2','1,2',12,1780042523,'',1);

/*!40000 ALTER TABLE `dy_sys_account` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 dy_sys_dict
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_dict`;

CREATE TABLE `dy_sys_dict` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `sys` smallint NOT NULL DEFAULT '0',
  `ptype` int DEFAULT NULL,
  `pname` varchar(35) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `remarks` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

LOCK TABLES `dy_sys_dict` WRITE;
/*!40000 ALTER TABLE `dy_sys_dict` DISABLE KEYS */;

INSERT INTO `dy_sys_dict` (`id`, `sys`, `ptype`, `pname`, `remarks`)
VALUES
	(1,0,3,'TableType','数据表类型'),
	(2,0,3,'FieldType','字段类型'),
	(3,0,3,'YesNo','是否'),
	(4,0,3,'KeyType','字典KEY类型'),
	(5,0,1,'sex','性别'),
	(6,0,1,'color','颜色');

/*!40000 ALTER TABLE `dy_sys_dict` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 dy_sys_dict_value
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_dict_value`;

CREATE TABLE `dy_sys_dict_value` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `dict_id` bigint DEFAULT NULL,
  `keyid` int DEFAULT NULL,
  `defaultval` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `keystr` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `fid` int DEFAULT NULL,
  `state` smallint NOT NULL DEFAULT '1',
  `uptime` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

LOCK TABLES `dy_sys_dict_value` WRITE;
/*!40000 ALTER TABLE `dy_sys_dict_value` DISABLE KEYS */;

INSERT INTO `dy_sys_dict_value` (`id`, `dict_id`, `keyid`, `defaultval`, `keystr`, `fid`, `state`, `uptime`)
VALUES
	(1,1,1,'业务类','',0,1,0),
	(2,1,2,'数据类','',0,1,0),
	(3,2,1,'int',NULL,0,1,0),
	(4,2,2,'smallint',NULL,0,1,0),
	(5,2,3,'tinyint',NULL,0,1,0),
	(6,2,4,'float',NULL,0,1,0),
	(7,2,5,'double',NULL,0,1,0),
	(8,2,6,'varchar',NULL,0,1,0),
	(9,2,7,'text',NULL,0,1,0),
	(10,2,8,'longtext',NULL,0,1,0),
	(11,2,9,'date',NULL,0,1,0),
	(12,2,10,'int64',NULL,0,1,0),
	(13,3,1,'是',NULL,0,1,0),
	(14,3,2,'否',NULL,0,1,0),
	(15,4,1,'整数键',NULL,0,1,0),
	(16,4,2,'字符键',NULL,0,1,0),
	(17,5,1,'男','B1',0,1,0),
	(18,5,2,'女','G',0,1,0),
	(19,6,1,'红','',0,1,0),
	(20,6,2,'蓝','',0,1,0);

/*!40000 ALTER TABLE `dy_sys_dict_value` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 dy_sys_dict_value_svn
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_dict_value_svn`;

CREATE TABLE `dy_sys_dict_value_svn` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `dict_value_id` bigint DEFAULT NULL,
  `oldstr` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `newstr` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `uptime` int DEFAULT NULL,
  `post_uid` bigint DEFAULT NULL,
  `post_user` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `dict_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



# 转储表 dy_sys_group
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_group`;

CREATE TABLE `dy_sys_group` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键 id',
  `gname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组名称',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '图标',
  `indextpl` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '首页模板',
  `state` int NOT NULL DEFAULT '1' COMMENT '状态',
  `ordnum` int NOT NULL COMMENT '排序号',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sys_groups_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

LOCK TABLES `dy_sys_group` WRITE;
/*!40000 ALTER TABLE `dy_sys_group` DISABLE KEYS */;

INSERT INTO `dy_sys_group` (`id`, `gname`, `icon`, `indextpl`, `state`, `ordnum`, `created_at`, `updated_at`, `deleted_at`)
VALUES
	(1,'系统设置','set.png','index/index_system.htm',1,1,'2026-05-28 13:36:12','2026-05-28 13:36:12',NULL),
	(2,'自建应用','zy.png','index/index_system.htm',1,2,'2026-05-28 13:36:33','2026-05-28 13:36:33',NULL);

/*!40000 ALTER TABLE `dy_sys_group` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 dy_sys_menu
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_menu`;

CREATE TABLE `dy_sys_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `mname` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '菜单名',
  `mlink` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '链接',
  `icon` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '图标',
  `pid` bigint DEFAULT NULL COMMENT '父类',
  `ordnum` int DEFAULT NULL COMMENT '排序',
  `status` smallint NOT NULL DEFAULT '1' COMMENT '状态 1 正常',
  `typeid` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单';

LOCK TABLES `dy_sys_menu` WRITE;
/*!40000 ALTER TABLE `dy_sys_menu` DISABLE KEYS */;

INSERT INTO `dy_sys_menu` (`id`, `mname`, `mlink`, `icon`, `pid`, `ordnum`, `status`, `typeid`)
VALUES
	(1,'设置','javascript:;','',0,1,1,1),
	(2,'菜单设置','?module=sysmenu&act=main&do=list','',1,1,1,1),
	(3,'数据字典','?module=sysdict&act=main&do=list','',1,2,1,1),
	(4,'权限管理','?module=sysrole&act=main&do=list','',1,3,1,1),
	(5,'系统模块','?module=sysmodule&act=main&do=list','',1,4,1,1),
	(6,'账号管理','?module=sysaccount&act=main&do=list','',1,5,1,1),
	(7,'系统日志','?module=syslog&act=audit&do=list','',1,6,1,1),
	(8,'系统分组','?module=sysgroup&act=main&do=list','',1,0,1,1);

/*!40000 ALTER TABLE `dy_sys_menu` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 dy_sys_module
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_module`;

CREATE TABLE `dy_sys_module` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `mname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `remarks` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `pid` bigint DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='OA模块组';

LOCK TABLES `dy_sys_module` WRITE;
/*!40000 ALTER TABLE `dy_sys_module` DISABLE KEYS */;

INSERT INTO `dy_sys_module` (`id`, `mname`, `remarks`, `pid`)
VALUES
	(1,'sysgroup','系统组',0),
	(2,'main','',1),
	(3,'sysmodule','系统模块',0),
	(4,'main','',3),
	(5,'sysmenu','系统菜单',0),
	(6,'main','',5),
	(7,'sysrole','系统权限',0),
	(8,'main','',7),
	(9,'sysdict','系统字典',0),
	(10,'main','',9),
	(11,'sysaccount','系统账号',0),
	(12,'main','',11),
	(13,'index','',0),
	(14,'main','',13),
	(15,'file','',0),
	(18,'systable','数据表',0),
	(19,'tab','',18),
	(20,'field','',18),
	(21,'syslog','系统日志',0),
	(22,'audit','审计日志',21);

/*!40000 ALTER TABLE `dy_sys_module` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 dy_sys_operation_log
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_operation_log`;

CREATE TABLE `dy_sys_operation_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作人员账号',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作人员昵称',
  `module` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作模块',
  `action` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作动作',
  `do` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作事件',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作描述(如: 删除菜单[ID:10])',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '请求URL',
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '请求方式(GET/POST/PUT/DELETE)',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作IP',
  `user_agent` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '浏览器UA',
  `param` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '请求参数/提交的表单',
  `data_old` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '修改前的数据(JSON), 仅UPDATE/DELETE有用',
  `data_new` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '修改后的数据(JSON), 仅CREATE/UPDATE有用',
  `status` smallint NOT NULL DEFAULT '1' COMMENT '1 操作成功 2 操作失败',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '错误原因',
  `latency` bigint NOT NULL DEFAULT '0' COMMENT '执行耗时(毫秒)',
  `created_at` datetime(3) NOT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='系统操作/审计日志表';



# 转储表 dy_sys_role
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_role`;

CREATE TABLE `dy_sys_role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '标识唯一',
  `rname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '名称',
  `remarks` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `menucate` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '菜单大类',
  `menulist` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '菜单列表',
  `modulelist` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `sysgroup` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '系统组',
  `rolelist` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sys_role_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

LOCK TABLES `dy_sys_role` WRITE;
/*!40000 ALTER TABLE `dy_sys_role` DISABLE KEYS */;

INSERT INTO `dy_sys_role` (`id`, `name`, `rname`, `remarks`, `menucate`, `menulist`, `modulelist`, `sysgroup`, `rolelist`, `created_at`, `updated_at`, `deleted_at`)
VALUES
	(1,'common','通用组','2','161,148','177,178','14','2','index:main','2026-05-22 10:35:36','2026-05-26 11:36:44',NULL),
	(2,'admin','管理员','11','1','5,6,7,8,9,2,3,4','2,4,6,8,10,12,14,16,19,20,22','1','sysgroup:main,sysmodule:main,sysmenu:main,sysrole:main,sysdict:main,sysaccount:main,index:main,file:2add1,systable:tab,systable:field,syslog:audit','2026-05-22 16:54:39','2026-05-28 13:14:08',NULL);

/*!40000 ALTER TABLE `dy_sys_role` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 dy_sys_token
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dy_sys_token`;

CREATE TABLE `dy_sys_token` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `token` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `uid` int DEFAULT NULL,
  `account` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
