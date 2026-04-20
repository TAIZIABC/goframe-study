# dbdiff — MySQL 数据库结构对比工具

连接两个 MySQL 数据库，对比所有表的结构差异（新增/删除表、字段差异、索引差异），自动生成同步 SQL。

## 使用方式

```bash
# 基本用法：对比开发库和测试库
go run main.go \
  -source 'root:123@tcp(localhost:3306)/dev_db' \
  -target 'root:123@tcp(localhost:3306)/test_db'

# 输出到文件
go run main.go \
  -source '...' -target '...' \
  -out sync.sql

# 不生成 DROP TABLE（安全模式）
go run main.go \
  -source '...' -target '...' \
  -no-drop -out sync.sql
```

## 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-source` | (必填) | 源库 DSN（期望的结构），如开发环境 |
| `-target` | (必填) | 目标库 DSN（待同步），如测试环境 |
| `-out` | - | SQL 输出文件路径（不指定则打印到终端） |
| `-no-drop` | `false` | 不生成 DROP TABLE 语句 |

## 输出示例

```
  🔍 连接源库...
  ✓ 源库 [dev_db]: 12 张表
  🔍 连接目标库...
  ✓ 目标库 [test_db]: 10 张表

  🔄 对比结构差异...

  发现 8 项差异:
  ─────────────────────────────────────────────────────────────────
    新增表: 2
    新增字段: 3
    字段变更: 2
    新增索引: 1
  ─────────────────────────────────────────────────────────────────
  + 新增表    orders
  + 新增表    payments
  + 新增字段  users.avatar  →  varchar(500) NULL
  ~ 字段变更  users.email
      源: varchar(200) NOT NULL
      目: varchar(100) NOT NULL
  ...

  生成 8 条 SQL 语句:

  CREATE TABLE IF NOT EXISTS `orders` (...);
  ALTER TABLE `users` ADD COLUMN `avatar` varchar(500) AFTER `email`;
  ALTER TABLE `users` MODIFY COLUMN `email` varchar(200) NOT NULL;
  ...
```

## 项目结构

```
dbdiff/
├── main.go               # CLI 入口
├── schema/
│   └── schema.go          # MySQL 表结构读取（information_schema）
├── diff/
│   └── diff.go            # 结构对比引擎（表/字段/索引级）
└── sqlgen/
    └── sqlgen.go           # SQL 语句生成（CREATE/ALTER/DROP）
```

## 检测的差异类型

| 差异 | 生成 SQL |
|------|---------|
| 新增表 | `CREATE TABLE ...` |
| 删除表 | `DROP TABLE ...`（可用 `-no-drop` 禁止） |
| 新增字段 | `ALTER TABLE ... ADD COLUMN ... AFTER ...` |
| 删除字段 | `ALTER TABLE ... DROP COLUMN ...` |
| 字段变更（类型/默认值/注释等） | `ALTER TABLE ... MODIFY COLUMN ...` |
| 新增索引 | `ALTER TABLE ... ADD INDEX ...` |
| 删除索引 | `ALTER TABLE ... DROP INDEX ...` |
| 索引变更 | 先 DROP 再 ADD |

## 功能特点

- 通过 `information_schema` 读取表结构，无需额外权限
- 检测字段的**类型、可空、默认值、Extra、注释**差异
- 检测索引的**列组合、唯一性、类型**差异
- ADD COLUMN 自动带 `AFTER` 子句保持字段顺序
- `-no-drop` 安全模式避免误删表
- 彩色差异摘要（绿色新增/红色删除/黄色变更）
- SQL 可直接导出到文件用于生产执行
