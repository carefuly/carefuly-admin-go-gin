/*
 Navicat Premium Dump SQL

 Source Server         : MySQL@192.168.66.102
 Source Server Type    : MySQL
 Source Server Version : 80404 (8.4.4)
 Source Host           : 192.168.66.102:3306
 Source Schema         : carefuly-admin-go-gin

 Target Server Type    : MySQL
 Target Server Version : 80404 (8.4.4)
 File Encoding         : 65001

 Date: 14/05/2025 10:33:20
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for careful_system_menu
-- ----------------------------
DROP TABLE IF EXISTS `careful_system_menu`;
CREATE TABLE `careful_system_menu`  (
  `id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `sort` bigint NULL DEFAULT 1 COMMENT '显示排序',
  `version` bigint NULL DEFAULT 1 COMMENT '版本号',
  `creator` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人',
  `modifier` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '修改人',
  `belong_dept` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '数据归属部门',
  `create_time` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '修改时间',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
  `menuType` tinyint NULL DEFAULT NULL COMMENT '菜单类型',
  `icon` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'HomeFilled' COMMENT '菜单图标',
  `title` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单标题',
  `permission` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '权限标识',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '组件名称',
  `component` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '组件地址',
  `api` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '接口地址',
  `method` tinyint NULL DEFAULT NULL COMMENT '接口请求方法',
  `path` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '路由地址',
  `redirect` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '重定向地址',
  `isHide` tinyint NULL DEFAULT 0 COMMENT '是否隐藏',
  `isLink` tinyint NULL DEFAULT 0 COMMENT '是否外链',
  `linkUrl` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '外链地址',
  `isKeepAlive` tinyint NULL DEFAULT 0 COMMENT '是否页面缓存',
  `isFull` tinyint NULL DEFAULT 0 COMMENT '是否缓存全屏',
  `isAffix` tinyint NULL DEFAULT 0 COMMENT '是否缓存固定路由',
  `parent_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '上级菜单',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uni_menu_title_unique`(`menuType` ASC, `title` ASC, `parent_id` ASC) USING BTREE,
  INDEX `idx_careful_system_menu_creator`(`creator` ASC) USING BTREE,
  INDEX `idx_careful_system_menu_modifier`(`modifier` ASC) USING BTREE,
  INDEX `idx_careful_system_menu_belong_dept`(`belong_dept` ASC) USING BTREE,
  INDEX `idx_title`(`title` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of careful_system_menu
-- ----------------------------
INSERT INTO `careful_system_menu` VALUES ('0DDDDB7C-936C-4A9C-8F8C-9BDD1EA534FB', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Menu', '菜单管理', 'system:menu', 'menu', 'system/menu/index', NULL, NULL, '/system/menu', NULL, 0, 0, NULL, 1, 0, 0, '91EB6373-BC7F-437E-AECF-2B45C89F805E');
INSERT INTO `careful_system_menu` VALUES ('164D9139-7E5E-4D90-94E9-8774C7BD5DF9', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Coffee', '存储桶', 'tools:bucket', 'bucket', 'tools/bucket/index', NULL, NULL, '/tools/bucket', NULL, 0, 0, NULL, 1, 0, 0, 'B6716CC6-524E-4028-8BD6-2E2A4EAA1AED');
INSERT INTO `careful_system_menu` VALUES ('2BC480DF-51E5-4EA8-A00D-CCAD68BB3089', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Histogram', '部门管理', 'system:dept', 'dept', 'system/dept/index', NULL, NULL, '/system/dept', NULL, 0, 0, NULL, 1, 0, 0, '91EB6373-BC7F-437E-AECF-2B45C89F805E');
INSERT INTO `careful_system_menu` VALUES ('37278C2F-7E96-4759-B194-4BF2EBD55AE1', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 1, 'Reading', '日志管理', 'logs', 'logs', NULL, NULL, NULL, '/logs', '/logs/login_log', 0, 0, NULL, 0, 0, 0, NULL);
INSERT INTO `careful_system_menu` VALUES ('3BB96B01-5D21-48E6-B1B5-2C42BF7B4586', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'CircleClose', '异常日志', 'logs:error_log', 'error_log', 'logs/error_log/index', NULL, NULL, '/logs/error_log', NULL, 0, 0, NULL, 1, 0, 0, '37278C2F-7E96-4759-B194-4BF2EBD55AE1');
INSERT INTO `careful_system_menu` VALUES ('50B8A4BB-18C1-4D02-9552-0035EFE1808C', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 1, 'Setting', '系统管理', 'system', 'system', NULL, NULL, NULL, '/system', '/system/user', 0, 0, NULL, 0, 0, 0, NULL);
INSERT INTO `careful_system_menu` VALUES ('541A6AE6-4728-4939-80B7-060B6AF67E0B', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Management', '数据字典', 'tools:dict', 'dict', 'tools/dict/index', NULL, NULL, '/tools/dict', NULL, 0, 0, NULL, 1, 0, 0, 'B6716CC6-524E-4028-8BD6-2E2A4EAA1AED');
INSERT INTO `careful_system_menu` VALUES ('56D8A8C2-D7AF-44E1-B011-367DC06939DB', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'DataAnalysis', '分析页', 'dashboard:analysis', 'analysis', 'dashboard/analysis/index', NULL, NULL, '/dashboard/analysis', NULL, 0, 0, NULL, 1, 0, 1, '89B1621F-41ED-4052-AC4F-A2037CE259E5');
INSERT INTO `careful_system_menu` VALUES ('56F8EDD7-C68E-4AFE-AE41-EE4AC458D0BC', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Avatar', '角色管理', 'system:role', 'role', 'system/role/index', NULL, NULL, '/system/role', NULL, 0, 0, NULL, 1, 0, 0, '91EB6373-BC7F-437E-AECF-2B45C89F805E');
INSERT INTO `careful_system_menu` VALUES ('57790752-0DD2-4F9D-B0EA-3D7436B47C13', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Postcard', '岗位管理', 'system:post', 'post', 'system/post/index', NULL, NULL, '/system/post', NULL, 0, 0, NULL, 1, 0, 0, '91EB6373-BC7F-437E-AECF-2B45C89F805E');
INSERT INTO `careful_system_menu` VALUES ('5A36E968-EB7A-4249-BD64-DB1F3A91166F', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'List', '字典信息', 'tools:dicType', 'dicType', 'tools/dicType/index', NULL, NULL, '/tools/dicType', NULL, 0, 0, NULL, 1, 0, 0, 'B6716CC6-524E-4028-8BD6-2E2A4EAA1AED');
INSERT INTO `careful_system_menu` VALUES ('5C3762DE-FB83-4FF9-868B-041CC984A049', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'DataLine', '控制台', 'dashboard:console', 'console', 'dashboard/console/index', NULL, NULL, '/dashboard/console', NULL, 0, 0, NULL, 1, 0, 1, '89B1621F-41ED-4052-AC4F-A2037CE259E5');
INSERT INTO `careful_system_menu` VALUES ('66D7AE24-EB0C-4D0E-9D08-D095F1C328D9', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Notebook', '操作日志', 'logs:operate_log', 'opreate_log', 'logs/operate_log/index', NULL, NULL, '/logs/operate_log', NULL, 0, 0, NULL, 1, 0, 0, '37278C2F-7E96-4759-B194-4BF2EBD55AE1');
INSERT INTO `careful_system_menu` VALUES ('89B1621F-41ED-4052-AC4F-A2037CE259E5', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 1, 'HomeFilled', '仪表盘', 'dashboard', 'dashboard', NULL, NULL, NULL, '/dashboard', '/dashboard/analysis', 0, 0, NULL, 0, 0, 0, NULL);
INSERT INTO `careful_system_menu` VALUES ('8D5C09C5-7966-4862-AAF6-859F2C562D26', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Calendar', '登录日志', 'logs:login_log', 'login_log', 'logs/login_log/index', NULL, NULL, '/logs/login_log', NULL, 0, 0, NULL, 1, 0, 0, '37278C2F-7E96-4759-B194-4BF2EBD55AE1');
INSERT INTO `careful_system_menu` VALUES ('91EB6373-BC7F-437E-AECF-2B45C89F805E', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'UserFilled', '用户管理', 'system:user', 'user', 'system/user/index', NULL, NULL, '/system/user', NULL, 0, 0, NULL, 1, 0, 0, '91EB6373-BC7F-437E-AECF-2B45C89F805E');
INSERT INTO `careful_system_menu` VALUES ('91F31C20-5751-4698-91A6-444EB9B0D0A4', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Pear', '存储桶详情', 'tools:bucket:detail', 'bucketDetail', 'tools/bucket/detail', NULL, NULL, '/tools/bucket/detail/:params', NULL, 1, 0, NULL, 1, 0, 0, 'B6716CC6-524E-4028-8BD6-2E2A4EAA1AED');
INSERT INTO `careful_system_menu` VALUES ('B6716CC6-524E-4028-8BD6-2E2A4EAA1AED', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 1, 'FirstAidKit', '系统工具', 'tools', 'tools', NULL, NULL, NULL, '/tools', '/tools/dict', 0, 0, NULL, 0, 0, 0, NULL);
INSERT INTO `careful_system_menu` VALUES ('CB252391-3945-482C-A730-997927AD4CC0', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Monitor', '系统信息', 'tools:config', 'config', 'tools/config/index', NULL, NULL, '/tools/config', NULL, 0, 0, NULL, 1, 0, 0, 'B6716CC6-524E-4028-8BD6-2E2A4EAA1AED');
INSERT INTO `careful_system_menu` VALUES ('D6539824-9EA4-42B4-819D-DA08724F6393', 1, 1, 'careful', 'careful', NULL, '2025-05-13 16:24:10.856', '2025-05-13 16:24:28.451', NULL, 2, 'Upload', '文件上传', 'tools:upload', 'upload', 'tools/upload/index', NULL, NULL, '/tools/upload', NULL, 0, 0, NULL, 1, 0, 0, 'B6716CC6-524E-4028-8BD6-2E2A4EAA1AED');

SET FOREIGN_KEY_CHECKS = 1;
