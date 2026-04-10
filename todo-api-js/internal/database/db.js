/**
 * 数据库连接模块
 * -----------------------------------------------
 * 对标 GoFrame: g.Model("todos") / g.DB()
 * GoFrame 通过配置文件自动创建数据库连接池
 * JS 版本使用 mysql2 的连接池，提供类似 g.Model() 的查询能力
 * -----------------------------------------------
 */

const mysql = require('mysql2/promise');
const config = require('../config/config');

let pool = null;

/**
 * 获取数据库连接池（单例）
 * 对标 GoFrame: g.DB() — 全局数据库对象
 */
function getPool() {
  if (!pool) {
    pool = mysql.createPool({
      host: config.database.host,
      port: config.database.port,
      user: config.database.user,
      password: config.database.password,
      database: config.database.database,
      waitForConnections: true,
      connectionLimit: 10,
      queueLimit: 0,
    });
  }
  return pool;
}

/**
 * 执行 SQL 查询的便捷方法
 * 对标 GoFrame: g.Model("table").Ctx(ctx).Where(...).Scan(...)
 */
async function query(sql, params = []) {
  const [rows] = await getPool().execute(sql, params);
  return rows;
}

/**
 * 执行 SQL 并返回结果信息（用于 INSERT/UPDATE/DELETE）
 * 对标 GoFrame: result, err := g.Model("table").Data(...).Insert()
 */
async function execute(sql, params = []) {
  const [result] = await getPool().execute(sql, params);
  return result;
}

module.exports = { getPool, query, execute };
