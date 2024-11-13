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

 Date: 12/11/2024 22:15:45
*/

SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_user_token
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_user_token`
(
    `id`            int(11) NOT NULL AUTO_INCREMENT,
    `user_id`       int(11) NOT NULL,
    `access_token`  varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `expires_at`    DATETIME NULL,
    `refresh_token` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `re_expires_at` DATETIME NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `sys_user_token_access_token_key` (`access_token`),
    UNIQUE KEY `sys_user_token_refresh_token_key` (`refresh_token`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET
FOREIGN_KEY_CHECKS = 1;
