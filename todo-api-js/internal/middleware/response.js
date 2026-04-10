/**
 * 统一响应中间件
 * -----------------------------------------------
 * 对标 GoFrame: ghttp.MiddlewareHandlerResponse
 *
 * GoFrame 的 MiddlewareHandlerResponse 会自动将控制器返回值包装为：
 *   { "code": 0, "message": "", "data": { ... } }
 *
 * 在 JS 版本中，我们在 controller 层已经手动包装了响应格式。
 * 这个中间件主要处理：
 *   1. 全局错误捕获（对标 GoFrame 的 err != nil 自动返回错误码）
 *   2. 404 路由未找到
 * -----------------------------------------------
 */

/**
 * 错误处理中间件
 * 当控制器抛出错误时，自动包装为统一错误响应
 * 对标 GoFrame: 当 controller 返回 err 时，中间件自动将 code 设为非0
 */
function errorHandler(err, req, res, next) {
  console.error(`[ERROR] ${req.method} ${req.url} -`, err.message);
  res.status(500).json({
    code: 50,
    message: err.message || 'Internal Server Error',
    data: null,
  });
}

/**
 * 404 处理中间件
 */
function notFoundHandler(req, res) {
  res.status(404).json({
    code: 65,
    message: `Route not found: ${req.method} ${req.url}`,
    data: null,
  });
}

module.exports = { errorHandler, notFoundHandler };
