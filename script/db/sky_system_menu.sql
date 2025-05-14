/*
 Navicat Premium Dump SQL

 Source Server         : MySQL@127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 80012 (8.0.12)
 Source Host           : 127.0.0.1:3306
 Source Schema         : sky admin pro

 Target Server Type    : MySQL
 Target Server Version : 80012 (8.0.12)
 File Encoding         : 65001

 Date: 13/05/2025 15:55:48
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sky_system_menu
-- ----------------------------
DROP TABLE IF EXISTS `sky_system_menu`;
CREATE TABLE `sky_system_menu`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sort` int(11) NULL DEFAULT NULL,
  `creator` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `modifier` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `belong_dept` int(11) NULL DEFAULT NULL,
  `remark` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `status` tinyint(1) NOT NULL,
  `create_time` datetime NULL DEFAULT NULL,
  `update_time` datetime NULL DEFAULT NULL,
  `type` smallint(6) NOT NULL,
  `icon` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `title` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `permission` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `component` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `api` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `path` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `redirect` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `isHide` tinyint(1) NOT NULL,
  `isLink` varchar(520) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NULL DEFAULT NULL,
  `isKeepAlive` tinyint(1) NOT NULL,
  `isFull` tinyint(1) NOT NULL,
  `isAffix` tinyint(1) NOT NULL,
  `parent_id` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `sky_system_menu_parent_id_8aef2e6a`(`parent_id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 270 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_520_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sky_system_menu
-- ----------------------------
INSERT INTO `sky_system_menu` VALUES (1, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:10:47', '2024-10-25 00:10:50', 1, 'HomeFilled', '仪表盘', 'dashboard', 'dashboard', NULL, NULL, NULL, '/dashboard', '/dashboard/analysis', 0, NULL, 0, 0, 0, NULL);
INSERT INTO `sky_system_menu` VALUES (2, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:14:59', '2024-10-25 00:15:02', 2, 'DataAnalysis', '分析页', 'dashboard:analysis', 'analysis', 'dashboard/analysis/index', NULL, NULL, '/dashboard/analysis', NULL, 0, NULL, 1, 0, 1, 1);
INSERT INTO `sky_system_menu` VALUES (3, 2, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:15:04', '2024-10-25 00:15:09', 2, 'DataLine', '控制台', 'dashboard:console', 'console', 'dashboard/console/index', NULL, NULL, '/dashboard/console', NULL, 0, NULL, 1, 0, 1, 1);
INSERT INTO `sky_system_menu` VALUES (4, 2, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:08', '2024-10-25 00:20:10', 1, 'Setting', '系统管理', 'system', 'system', NULL, NULL, NULL, '/system', '/system/user', 0, NULL, 0, 0, 0, NULL);
INSERT INTO `sky_system_menu` VALUES (5, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:13', '2024-10-25 00:20:15', 2, 'UserFilled', '用户管理', 'system:user', 'user', 'system/user/index', NULL, NULL, '/system/user', NULL, 0, NULL, 1, 0, 0, 4);
INSERT INTO `sky_system_menu` VALUES (6, 2, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:17', '2024-10-25 00:20:32', 2, 'Avatar', '角色管理', 'system:role', 'role', 'system/role/index', NULL, NULL, '/system/role', NULL, 0, NULL, 1, 0, 0, 4);
INSERT INTO `sky_system_menu` VALUES (7, 3, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:18', '2024-10-25 00:20:34', 2, 'Menu', '菜单管理', 'system:menu', 'menu', 'system/menu/index', NULL, NULL, '/system/menu', NULL, 0, NULL, 1, 0, 0, 4);
INSERT INTO `sky_system_menu` VALUES (8, 4, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:21', '2024-10-25 00:20:35', 2, 'Histogram', '部门管理', 'system:dept', 'dept', 'system/dept/index', NULL, NULL, '/system/dept', NULL, 0, NULL, 1, 0, 0, 4);
INSERT INTO `sky_system_menu` VALUES (9, 5, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:22', '2024-10-25 00:20:37', 2, 'Postcard', '岗位管理', 'system:post', 'post', 'system/post/index', NULL, NULL, '/system/post', NULL, 0, NULL, 1, 0, 0, 4);
INSERT INTO `sky_system_menu` VALUES (10, 3, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:24', '2024-10-25 00:20:39', 1, 'FirstAidKit', '系统工具', 'tools', 'tools', NULL, NULL, NULL, '/tools', '/tools/dict', 0, NULL, 0, 0, 0, NULL);
INSERT INTO `sky_system_menu` VALUES (11, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:25', '2024-10-25 00:20:41', 2, 'Management', '数据字典', 'tools:dict', 'dict', 'tools/dict/index', NULL, NULL, '/tools/dict', NULL, 0, NULL, 1, 0, 0, 10);
INSERT INTO `sky_system_menu` VALUES (12, 2, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:29', '2024-10-25 00:20:42', 2, 'List', '字典信息', 'tools:dicType', 'dicType', 'tools/dicType/index', NULL, NULL, '/tools/dicType', NULL, 0, NULL, 1, 0, 0, 10);
INSERT INTO `sky_system_menu` VALUES (13, 3, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-25 00:20:29', '2024-10-25 00:20:42', 2, 'Monitor', '系统信息', 'tools:config', 'config', 'tools/config/index', NULL, NULL, '/tools/config', NULL, 0, NULL, 1, 0, 0, 10);
INSERT INTO `sky_system_menu` VALUES (14, 4, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-28 20:20:26', '2024-10-25 00:20:42', 2, 'Coffee', '存储桶', 'tools:bucket', 'bucket', 'tools/bucket/index', NULL, NULL, '/tools/bucket', NULL, 0, NULL, 1, 0, 0, 10);
INSERT INTO `sky_system_menu` VALUES (15, 5, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-28 20:20:26', '2024-10-25 00:20:42', 2, 'Pear', '存储桶详情', 'tools:bucket:detail', 'bucketDetail', 'tools/bucket/detail', NULL, NULL, '/tools/bucket/detail/:params', NULL, 1, NULL, 1, 0, 0, 10);
INSERT INTO `sky_system_menu` VALUES (16, 6, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-10-28 20:20:26', '2024-10-25 00:20:42', 2, 'Upload', '文件上传', 'tools:upload', 'upload', 'tools/upload/index', NULL, NULL, '/tools/upload', NULL, 0, NULL, 1, 0, 0, 10);
INSERT INTO `sky_system_menu` VALUES (19, 4, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-06 16:16:10', '2024-11-06 16:16:10', 1, 'Reading', '日志管理', 'logs', 'logs', NULL, NULL, NULL, '/logs', '/logs/login_log', 0, NULL, 0, 0, 0, NULL);
INSERT INTO `sky_system_menu` VALUES (25, 3, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-06 16:35:34', '2024-11-09 20:31:13', 2, 'CircleClose', '异常日志', 'logs:error_log', 'error_log', 'logs/error_log/index', NULL, NULL, '/logs/error_log', NULL, 0, NULL, 1, 0, 0, 19);
INSERT INTO `sky_system_menu` VALUES (24, 2, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-06 16:35:34', '2024-11-06 16:35:34', 2, 'Notebook', '操作日志', 'logs:operate_log', 'opreate_log', 'logs/operate_log/index', NULL, NULL, '/logs/operate_log', NULL, 0, NULL, 1, 0, 0, 19);
INSERT INTO `sky_system_menu` VALUES (23, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-06 16:32:52', '2024-11-09 20:31:02', 2, 'Calendar', '登录日志', 'logs:login_log', 'login_log', 'logs/login_log/index', NULL, NULL, '/logs/login_log', NULL, 0, NULL, 1, 0, 0, 19);
INSERT INTO `sky_system_menu` VALUES (258, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:errorLog:batchDelete', 'button', NULL, '/sky/logs/errorLog/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 25);
INSERT INTO `sky_system_menu` VALUES (257, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:errorLog:delete', 'button', NULL, '/sky/logs/errorLog/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 25);
INSERT INTO `sky_system_menu` VALUES (255, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:operateLog:view', 'button', NULL, '/sky/logs/operateLog', '1', NULL, NULL, 1, NULL, 0, 0, 0, 24);
INSERT INTO `sky_system_menu` VALUES (256, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:errorLog:export', 'button', NULL, '/sky/logs/errorLog/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 25);
INSERT INTO `sky_system_menu` VALUES (253, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:operateLog:delete', 'button', NULL, '/sky/logs/operateLog/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 24);
INSERT INTO `sky_system_menu` VALUES (254, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:operateLog:batchDelete', 'button', NULL, '/sky/logs/operateLog/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 24);
INSERT INTO `sky_system_menu` VALUES (252, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:operateLog:export', 'button', NULL, '/sky/logs/operateLog/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 24);
INSERT INTO `sky_system_menu` VALUES (251, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:loginLog:view', 'button', NULL, '/sky/logs/loginLog', '1', NULL, NULL, 1, NULL, 0, 0, 0, 23);
INSERT INTO `sky_system_menu` VALUES (250, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:loginLog:batchDelete', 'button', NULL, '/sky/logs/loginLog/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 23);
INSERT INTO `sky_system_menu` VALUES (249, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:loginLog:delete', 'button', NULL, '/sky/logs/loginLog/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 23);
INSERT INTO `sky_system_menu` VALUES (247, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:file:view', 'button', NULL, '/sky/tools/file', '1', NULL, NULL, 1, NULL, 0, 0, 0, 15);
INSERT INTO `sky_system_menu` VALUES (248, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:loginLog:export', 'button', NULL, '/sky/logs/loginLog/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 23);
INSERT INTO `sky_system_menu` VALUES (246, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:file:update', 'button', NULL, '/sky/tools/file/update/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 15);
INSERT INTO `sky_system_menu` VALUES (245, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:file:batchDelete', 'button', NULL, '/sky/tools/file/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 15);
INSERT INTO `sky_system_menu` VALUES (244, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:file:delete', 'button', NULL, '/sky/tools/file/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 15);
INSERT INTO `sky_system_menu` VALUES (242, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Upload', '批量上传', 'system:file:batchUpload', 'button', NULL, '/sky/tools/uploadFile/batch/', '2', NULL, NULL, 1, NULL, 0, 0, 0, 15);
INSERT INTO `sky_system_menu` VALUES (243, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:file:export', 'button', NULL, '/sky/tools/file/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 15);
INSERT INTO `sky_system_menu` VALUES (241, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Upload', '上传文件', 'system:file:upload', 'button', NULL, '/sky/tools/uploadFile/', '2', NULL, NULL, 1, NULL, 0, 0, 0, 15);
INSERT INTO `sky_system_menu` VALUES (239, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:bucket:update', 'button', NULL, '/sky/tools/bucket/update/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 14);
INSERT INTO `sky_system_menu` VALUES (240, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:bucket:view', 'button', NULL, '/sky/tools/bucket', '1', NULL, NULL, 1, NULL, 0, 0, 0, 14);
INSERT INTO `sky_system_menu` VALUES (238, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:bucket:batchDelete', 'button', NULL, '/sky/tools/bucket/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 14);
INSERT INTO `sky_system_menu` VALUES (237, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:bucket:delete', 'button', NULL, '/sky/tools/bucket/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 14);
INSERT INTO `sky_system_menu` VALUES (235, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Upload', '导入', 'system:bucket:import', 'button', NULL, '/sky/tools/bucket/import', '2', NULL, NULL, 1, NULL, 0, 0, 0, 14);
INSERT INTO `sky_system_menu` VALUES (236, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:bucket:export', 'button', NULL, '/sky/tools/bucket/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 14);
INSERT INTO `sky_system_menu` VALUES (234, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Plus', '新增', 'system:bucket:add', 'button', NULL, '/sky/tools/bucket/add', '2', NULL, NULL, 1, NULL, 0, 0, 0, 14);
INSERT INTO `sky_system_menu` VALUES (233, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:config:view', 'button', NULL, '/sky/tools/config', '1', NULL, NULL, 1, NULL, 0, 0, 0, 13);
INSERT INTO `sky_system_menu` VALUES (232, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:config:update', 'button', NULL, '/sky/tools/config/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 13);
INSERT INTO `sky_system_menu` VALUES (231, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:dicType:view', 'button', NULL, '/sky/tools/dicType', '1', NULL, NULL, 1, NULL, 0, 0, 0, 12);
INSERT INTO `sky_system_menu` VALUES (230, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:dicType:update', 'button', NULL, '/sky/tools/dicType/update/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 12);
INSERT INTO `sky_system_menu` VALUES (229, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:dicType:batchDelete', 'button', NULL, '/sky/tools/dicType/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 12);
INSERT INTO `sky_system_menu` VALUES (227, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:dicType:export', 'button', NULL, '/sky/tools/dicType/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 12);
INSERT INTO `sky_system_menu` VALUES (228, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:dicType:delete', 'button', NULL, '/sky/tools/dicType/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 12);
INSERT INTO `sky_system_menu` VALUES (226, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Upload', '导入', 'system:dicType:import', 'button', NULL, '/sky/tools/dicType/import', '2', NULL, NULL, 1, NULL, 0, 0, 0, 12);
INSERT INTO `sky_system_menu` VALUES (225, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Plus', '新增', 'system:dicType:add', 'button', NULL, '/sky/tools/dicType/add', '2', NULL, NULL, 1, NULL, 0, 0, 0, 12);
INSERT INTO `sky_system_menu` VALUES (224, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:dict:view', 'button', NULL, '/sky/tools/dict', '1', NULL, NULL, 1, NULL, 0, 0, 0, 11);
INSERT INTO `sky_system_menu` VALUES (223, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:dict:update', 'button', NULL, '/sky/tools/dict/update/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 11);
INSERT INTO `sky_system_menu` VALUES (222, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:dict:batchDelete', 'button', NULL, '/sky/tools/dict/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 11);
INSERT INTO `sky_system_menu` VALUES (221, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:dict:delete', 'button', NULL, '/sky/tools/dict/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 11);
INSERT INTO `sky_system_menu` VALUES (220, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:dict:export', 'button', NULL, '/sky/tools/dict/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 11);
INSERT INTO `sky_system_menu` VALUES (219, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Upload', '导入', 'system:dict:import', 'button', NULL, '/sky/tools/dict/import', '2', NULL, NULL, 1, NULL, 0, 0, 0, 11);
INSERT INTO `sky_system_menu` VALUES (217, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:post:view', 'button', NULL, '/sky/system/post', '1', NULL, NULL, 1, NULL, 0, 0, 0, 9);
INSERT INTO `sky_system_menu` VALUES (218, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Plus', '新增', 'system:dict:add', 'button', NULL, '/sky/tools/dict/add', '2', NULL, NULL, 1, NULL, 0, 0, 0, 11);
INSERT INTO `sky_system_menu` VALUES (216, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:post:update', 'button', NULL, '/sky/system/post/update/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 9);
INSERT INTO `sky_system_menu` VALUES (215, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:post:batchDelete', 'button', NULL, '/sky/system/post/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 9);
INSERT INTO `sky_system_menu` VALUES (214, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:post:delete', 'button', NULL, '/sky/system/post/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 9);
INSERT INTO `sky_system_menu` VALUES (213, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:post:export', 'button', NULL, '/sky/system/post/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 9);
INSERT INTO `sky_system_menu` VALUES (212, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Upload', '导入', 'system:post:import', 'button', NULL, '/sky/system/post/import', '2', NULL, NULL, 1, NULL, 0, 0, 0, 9);
INSERT INTO `sky_system_menu` VALUES (210, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:dept:view', 'button', NULL, '/sky/system/dept', '1', NULL, NULL, 1, NULL, 0, 0, 0, 8);
INSERT INTO `sky_system_menu` VALUES (211, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Plus', '新增', 'system:post:add', 'button', NULL, '/sky/system/post/add', '2', NULL, NULL, 1, NULL, 0, 0, 0, 9);
INSERT INTO `sky_system_menu` VALUES (209, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:dept:update', 'button', NULL, '/sky/system/dept/update/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 8);
INSERT INTO `sky_system_menu` VALUES (208, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:dept:batchDelete', 'button', NULL, '/sky/system/dept/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 8);
INSERT INTO `sky_system_menu` VALUES (207, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:dept:delete', 'button', NULL, '/sky/system/dept/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 8);
INSERT INTO `sky_system_menu` VALUES (206, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:dept:export', 'button', NULL, '/sky/system/dept/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 8);
INSERT INTO `sky_system_menu` VALUES (204, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Plus', '新增', 'system:dept:add', 'button', NULL, '/sky/system/dept/add', '2', NULL, NULL, 1, NULL, 0, 0, 0, 8);
INSERT INTO `sky_system_menu` VALUES (205, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Upload', '导入', 'system:dept:import', 'button', NULL, '/sky/system/dept/import', '2', NULL, NULL, 1, NULL, 0, 0, 0, 8);
INSERT INTO `sky_system_menu` VALUES (203, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:menu:view', 'button', NULL, '/sky/system/menu', '1', NULL, NULL, 1, NULL, 0, 0, 0, 7);
INSERT INTO `sky_system_menu` VALUES (202, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:menu:update', 'button', NULL, '/sky/system/menu/update/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 7);
INSERT INTO `sky_system_menu` VALUES (201, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:menu:batchDelete', 'button', NULL, '/sky/system/menu/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 7);
INSERT INTO `sky_system_menu` VALUES (200, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:menu:delete', 'button', NULL, '/sky/system/menu/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 7);
INSERT INTO `sky_system_menu` VALUES (199, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:menu:export', 'button', NULL, '/sky/system/menu/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 7);
INSERT INTO `sky_system_menu` VALUES (198, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Upload', '导入', 'system:menu:import', 'button', NULL, '/sky/system/menu/import', '2', NULL, NULL, 1, NULL, 0, 0, 0, 7);
INSERT INTO `sky_system_menu` VALUES (196, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:role:view', 'button', NULL, '/sky/system/role', '1', NULL, NULL, 1, NULL, 0, 0, 0, 6);
INSERT INTO `sky_system_menu` VALUES (197, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Plus', '新增', 'system:menu:add', 'button', NULL, '/sky/system/menu/add', '2', NULL, NULL, 1, NULL, 0, 0, 0, 7);
INSERT INTO `sky_system_menu` VALUES (195, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '分配权限', 'system:role:rolePermission', 'button', NULL, '/sky/system/role/saveRolePermission/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 6);
INSERT INTO `sky_system_menu` VALUES (194, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '分配菜单', 'system:role:roleMenu', 'button', NULL, '/sky/system/role/saveRoleMenu/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 6);
INSERT INTO `sky_system_menu` VALUES (193, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:role:update', 'button', NULL, '/sky/system/role/update/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 6);
INSERT INTO `sky_system_menu` VALUES (192, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:role:batchDelete', 'button', NULL, '/sky/system/role/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 6);
INSERT INTO `sky_system_menu` VALUES (191, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:role:delete', 'button', NULL, '/sky/system/role/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 6);
INSERT INTO `sky_system_menu` VALUES (190, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:role:export', 'button', NULL, '/sky/system/role/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 6);
INSERT INTO `sky_system_menu` VALUES (189, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Upload', '导入', 'system:role:import', 'button', NULL, '/sky/system/role/import', '2', NULL, NULL, 1, NULL, 0, 0, 0, 6);
INSERT INTO `sky_system_menu` VALUES (188, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Plus', '新增', 'system:role:add', 'button', NULL, '/sky/system/role/add', '2', NULL, NULL, 1, NULL, 0, 0, 0, 6);
INSERT INTO `sky_system_menu` VALUES (187, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:user:view', 'button', NULL, '/sky/system/user', '1', NULL, NULL, 1, NULL, 0, 0, 0, 5);
INSERT INTO `sky_system_menu` VALUES (186, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Edit', '修改', 'system:user:update', 'button', NULL, '/sky/system/user/update/', '3', NULL, NULL, 1, NULL, 0, 0, 0, 5);
INSERT INTO `sky_system_menu` VALUES (185, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '批量删除', 'system:user:batchDelete', 'button', NULL, '/sky/system/user/batch/del', '2', NULL, NULL, 1, NULL, 0, 0, 0, 5);
INSERT INTO `sky_system_menu` VALUES (184, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Delete', '删除', 'system:user:delete', 'button', NULL, '/sky/system/user/del/', '4', NULL, NULL, 1, NULL, 0, 0, 0, 5);
INSERT INTO `sky_system_menu` VALUES (183, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'download', '导出', 'system:user:export', 'button', NULL, '/sky/system/user/export', '2', NULL, NULL, 1, NULL, 0, 0, 0, 5);
INSERT INTO `sky_system_menu` VALUES (182, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'Plus', '新增', 'system:user:add', 'button', NULL, '/sky/system/user/add', '2', NULL, NULL, 1, NULL, 0, 0, 0, 5);
INSERT INTO `sky_system_menu` VALUES (259, 1, 'XiaoYu', 'XiaoYu', 3, NULL, 1, '2024-11-10 00:28:08', '2024-11-10 00:28:08', 3, 'View', '查询', 'system:errorLog:view', 'button', NULL, '/sky/logs/errorLog', '1', NULL, NULL, 1, NULL, 0, 0, 0, 25);
INSERT INTO `sky_system_menu` VALUES (260, 5, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-11 16:06:24', '2024-11-11 16:06:24', 1, 'DataBoard', '数据管理', 'data', 'data', NULL, NULL, NULL, '/data', '/data/database', 0, NULL, 0, 0, 0, NULL);
INSERT INTO `sky_system_menu` VALUES (261, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-11 16:07:28', '2024-11-11 16:12:09', 2, 'Coin', '数据库管理', 'data:database', 'database', 'data/database/index', NULL, NULL, '/data/database', NULL, 0, NULL, 1, 0, 0, 260);
INSERT INTO `sky_system_menu` VALUES (262, 2, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-11 16:08:53', '2024-11-11 16:12:19', 2, 'Monitor', 'Redis监控', 'data:redis', 'redis', 'data/redis/index', NULL, NULL, '/data/redis', NULL, 0, NULL, 1, 0, 0, 260);
INSERT INTO `sky_system_menu` VALUES (263, 3, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-11 16:10:50', '2024-11-11 16:12:30', 2, 'Pouring', '数据缓存', 'data:cache', 'cache', 'data/cache/index', NULL, NULL, '/data/cache', NULL, 0, NULL, 1, 0, 0, 260);
INSERT INTO `sky_system_menu` VALUES (264, 6, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-15 10:29:52', '2024-11-15 10:29:52', 1, 'CloseBold', '异常页面', 'error', 'error', NULL, NULL, NULL, '/error/404', NULL, 0, NULL, 0, 0, 0, NULL);
INSERT INTO `sky_system_menu` VALUES (265, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-15 10:30:39', '2024-11-15 10:30:39', 2, 'QuestionFilled', '403', 'error:403', '403', 'error/403/index', NULL, NULL, '/error/403', NULL, 0, NULL, 1, 0, 0, 264);
INSERT INTO `sky_system_menu` VALUES (266, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-15 10:31:25', '2024-11-15 10:31:25', 2, 'CircleCloseFilled', '404', 'error:404', '404', 'error/404/index', NULL, NULL, '/error/404', NULL, 0, NULL, 1, 0, 0, 264);
INSERT INTO `sky_system_menu` VALUES (267, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-15 10:31:56', '2024-11-15 10:31:56', 2, 'WarningFilled', '500', 'error:500', '500', 'error/500/index', NULL, NULL, '/error/500', NULL, 0, NULL, 1, 0, 0, 264);
INSERT INTO `sky_system_menu` VALUES (268, 7, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-15 10:38:47', '2024-11-15 10:38:47', 1, 'Link', '外部连接', 'link', 'link', NULL, NULL, NULL, '/link/ElementPlus', NULL, 0, NULL, 0, 0, 0, NULL);
INSERT INTO `sky_system_menu` VALUES (269, 1, 'SkyAdmin', 'SkyAdmin', 2, NULL, 1, '2024-11-15 10:39:33', '2024-11-15 10:39:33', 2, 'Eleme', 'ElementPlus', 'link:ElementPlus', 'ElementPlus', 'link/ElementPlus/index', NULL, NULL, '/link/ElementPlus', NULL, 0, NULL, 1, 0, 0, 268);

SET FOREIGN_KEY_CHECKS = 1;
