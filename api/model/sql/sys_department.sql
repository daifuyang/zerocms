/*
 Navicat MySQL Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50734
 Source Host           : localhost:3306
 Source Schema         : zerocms

 Target Server Type    : MySQL
 Target Server Version : 50734
 File Encoding         : 65001

 Date: 14/11/2024 10:15:10
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_department
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_department`
(
    `id`         int(11)                                NOT NULL AUTO_INCREMENT COMMENT '部门id',
    `parentId`   int(11)                                NOT NULL                    DEFAULT '0' COMMENT '父部门id',
    `ancestors`  varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '组级列表' DEFAULT '',
    `name`       varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部门名称' DEFAULT '',
    `sort`       int(4)                                 NOT NULL                    DEFAULT '0' COMMENT '显示顺序',
    `leader`     varchar(20) COLLATE utf8mb4_unicode_ci                             DEFAULT NULL COMMENT '负责人',
    `phone`      varchar(11) COLLATE utf8mb4_unicode_ci                             DEFAULT NULL COMMENT '联系电话',
    `email`      varchar(50) COLLATE utf8mb4_unicode_ci                             DEFAULT NULL COMMENT '邮箱',
    `status`     tinyint(3)                             NOT NULL                    DEFAULT '1' COMMENT '0=>停用，1=>启用',
    `created_at` datetime                               NOT NULL                    DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime                               NOT NULL                    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `create_by`  int(11)                                NOT NULL                    DEFAULT 0 COMMENT '创建者',
    `update_by`  int(11)                                NOT NULL                    DEFAULT 0 COMMENT '更新者',
    `deleted_at` datetime                                                           DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `sys_department_name_key` (`name`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
