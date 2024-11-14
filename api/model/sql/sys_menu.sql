/*
 Navicat MySQL Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50734
 Source Host           : localhost:3306
 Source Schema         : nextcms

 Target Server Type    : MySQL
 Target Server Version : 50734
 File Encoding         : 65001

 Date: 13/11/2024 20:28:41
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu` (
  `menu_id` int(11) NOT NULL AUTO_INCREMENT,
  `menu_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `order` int(11) NOT NULL DEFAULT '0',
  `path` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `component` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `query` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `is_frame` int(11) NOT NULL DEFAULT '1',
  `is_cache` int(11) NOT NULL DEFAULT '0',
  `menu_type` char(1) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `visible` tinyint(4) NOT NULL DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `perms` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `icon` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_id` int(11) DEFAULT NULL,
  `created_by` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` int(11) DEFAULT NULL,
  `updated_id` int(11) DEFAULT NULL,
  `updated_by` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `updated_at` int(11) DEFAULT NULL,
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`menu_id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
