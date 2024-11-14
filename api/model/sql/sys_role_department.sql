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

 Date: 14/11/2024 12:20:18
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_role_dept
-- ----------------------------
CREATE TABLE IF NOT EXISTS `sys_role_department`
(
    `role_id`       int(11) NOT NULL DEFAULT 0 COMMENT '角色ID',
    `department_id` int(11) NOT NULL DEFAULT 0 COMMENT '部门ID',
    PRIMARY KEY (role_id, department_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='角色和部门关联表';

SET FOREIGN_KEY_CHECKS = 1;
