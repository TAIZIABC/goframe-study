-- ============================================
-- GoFrame Todo API - 数据库初始化脚本
-- ============================================

-- 1. 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `todo_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `todo_db`;

-- 2. 创建 todos 表
DROP TABLE IF EXISTS `todos`;
CREATE TABLE `todos` (
  `id`         INT UNSIGNED    NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `title`      VARCHAR(200)    NOT NULL DEFAULT ''     COMMENT '待办事项标题',
  `completed`  TINYINT(1)      NOT NULL DEFAULT 0      COMMENT '是否完成: 0-未完成, 1-已完成',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX `idx_completed` (`completed`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='待办事项表';

-- 3. 插入示例数据
INSERT INTO `todos` (`title`, `completed`) VALUES
  ('学习 Go 语言基础语法', 1),
  ('了解 GoFrame 框架目录结构', 1),
  ('完成 Todo API 的 CRUD 接口', 1),
  ('学习 GoFrame ORM 操作数据库', 0),
  ('学习 GoFrame 中间件机制', 0),
  ('实现 JWT 用户认证系统', 0),
  ('学习 GoFrame 数据校验', 0),
  ('部署项目到服务器', 0);
