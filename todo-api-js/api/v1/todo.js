/**
 * API 数据模型 & 校验规则
 * -----------------------------------------------
 * 对标 GoFrame: api/v1/model.go + api/v1/todo.go
 *
 * GoFrame 使用 struct tag 定义校验规则：
 *   Title string `v:"required|length:1,200" json:"title"`
 *
 * JS 版本用纯函数实现参数校验，达到相同效果
 * -----------------------------------------------
 */

/**
 * 校验创建 Todo 的请求参数
 * 对标 GoFrame: TodoCreateReq struct
 *   - Title:     string `v:"required|length:1,200"`
 *   - Completed: bool   (默认 false)
 */
function validateCreateReq(body) {
  const errors = [];
  if (!body.title || typeof body.title !== 'string') {
    errors.push('title 是必填项');
  } else if (body.title.length < 1 || body.title.length > 200) {
    errors.push('title 长度需在 1-200 之间');
  }
  if (body.completed !== undefined && typeof body.completed !== 'boolean') {
    errors.push('completed 必须是布尔值');
  }
  return errors;
}

/**
 * 校验列表查询的请求参数
 * 对标 GoFrame: TodoListReq struct
 *   - Page: int `d:"1"  v:"min:1"`
 *   - Size: int `d:"10" v:"max:100"`
 */
function validateListReq(query) {
  const errors = [];
  const page = parseInt(query.page) || 1;
  const size = parseInt(query.size) || 10;

  if (page < 1) {
    errors.push('page 最小值为 1');
  }
  if (size > 100) {
    errors.push('size 最大值为 100');
  }
  return { errors, page, size };
}

/**
 * 校验更新 Todo 的请求参数
 * 对标 GoFrame: TodoUpdateReq struct
 *   - Id:        int    `v:"required" in:"path"`
 *   - Title:     string `v:"length:1,200"`
 *   - Completed: *bool  (指针类型，区分 nil 和 false)
 */
function validateUpdateReq(params, body) {
  const errors = [];
  const id = parseInt(params.id);

  if (!id || id < 1) {
    errors.push('id 是必填项');
  }
  if (body.title !== undefined) {
    if (typeof body.title !== 'string' || body.title.length < 1 || body.title.length > 200) {
      errors.push('title 长度需在 1-200 之间');
    }
  }
  if (body.completed !== undefined && typeof body.completed !== 'boolean') {
    errors.push('completed 必须是布尔值');
  }
  return { errors, id };
}

/**
 * 校验删除 Todo 的请求参数
 * 对标 GoFrame: TodoDeleteReq struct
 *   - Id: int `v:"required" in:"path"`
 */
function validateDeleteReq(params) {
  const errors = [];
  const id = parseInt(params.id);

  if (!id || id < 1) {
    errors.push('id 是必填项');
  }
  return { errors, id };
}

module.exports = {
  validateCreateReq,
  validateListReq,
  validateUpdateReq,
  validateDeleteReq,
};
