/**
 * 配置加载模块
 * -----------------------------------------------
 * 对标 GoFrame: manifest/config/config.yaml 自动加载
 * GoFrame 通过 g.Cfg() 全局读取 YAML 配置
 * JS 版本使用 yaml 解析 + 单例模式
 * -----------------------------------------------
 */

// 简化实现：直接用 JS 对象定义配置（生产环境可用 js-yaml 库读取 YAML）
// 支持通过环境变量覆盖，对标 GoFrame 的配置优先级机制
const config = {
  server: {
    port: parseInt(process.env.PORT || '8000'),
  },
  database: {
    host: process.env.DB_HOST || '9.134.232.197',
    port: parseInt(process.env.DB_PORT || '3306'),
    user: process.env.DB_USER || 'root',
    password: process.env.DB_PASSWORD || 'lsj@24625',
    database: process.env.DB_NAME || 'todo_db',
  },
  logger: {
    level: process.env.LOG_LEVEL || 'all',
    stdout: true,
  },
};

module.exports = config;
