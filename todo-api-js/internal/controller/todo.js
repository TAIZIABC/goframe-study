/**
 * Todo Controller 层 — 请求处理
 * -----------------------------------------------
 * 对标 GoFrame: internal/controller/todo.go
 *
 * GoFrame 控制器通过 struct 方法自动绑定路由：
 *   func (c *cTodo) Create(ctx context.Context, req *v1.TodoCreateReq) (res *v1.TodoCreateRes, err error)
 *
 * GoFrame 框架自动完成：
 *   1. 从 HTTP 请求中解析参数到 req struct
 *   2. 根据 v tag 自动校验参数
 *   3. 将返回的 res struct 通过中间件包装为统一 JSON 响应
 *
 * JS 版本使用 Express handler 函数，手动完成上述三步
 * -----------------------------------------------
 */

const todoService = require('../service/todo');
const {
  validateCreateReq,
  validateListReq,
  validateUpdateReq,
  validateDeleteReq,
} = require('../../api/v1/todo');

const todoController = {
  /**
   * Create — 创建 Todo
   * -----------------------------------------------
   * 对标 Go:
   *   func (c *cTodo) Create(ctx context.Context, req *v1.TodoCreateReq) (*v1.TodoCreateRes, error) {
   *       id, err := service.Todo().Create(ctx, req.Title, req.Completed)
   *       return &v1.TodoCreateRes{Id: id}, nil
   *   }
   *
   * 路由: POST /api/v1/todos
   * 对标 Go Meta: `path:"/todos" method:"post"`
   * -----------------------------------------------
   */
  async create(req, res, next) {
    try {
      // 参数校验 — 对标 GoFrame 自动校验 (v tag)
      const errors = validateCreateReq(req.body);
      if (errors.length > 0) {
        return res.json({ code: 51, message: errors.join('; '), data: null });
      }

      const { title, completed = false } = req.body;
      const id = await todoService.create(title, completed);

      // 返回 — 对标 TodoCreateRes{Id: id}
      res.json({ code: 0, message: '', data: { id } });
    } catch (err) {
      next(err);
    }
  },

  /**
   * List — 分页查询 Todo 列表
   * -----------------------------------------------
   * 对标 Go:
   *   func (c *cTodo) List(ctx context.Context, req *v1.TodoListReq) (*v1.TodoListRes, error) {
   *       list, total, err := service.Todo().List(ctx, req.Page, req.Size)
   *       return &v1.TodoListRes{List: list, Total: total}, nil
   *   }
   *
   * 路由: GET /api/v1/todos
   * 对标 Go Meta: `path:"/todos" method:"get"`
   * -----------------------------------------------
   */
  async list(req, res, next) {
    try {
      // 参数解析 + 校验 — 对标 GoFrame 自动解析 query 参数 + d tag 默认值
      const { errors, page, size } = validateListReq(req.query);
      if (errors.length > 0) {
        return res.json({ code: 51, message: errors.join('; '), data: null });
      }

      const data = await todoService.list(page, size);

      // 返回 — 对标 TodoListRes{List: list, Total: total}
      res.json({ code: 0, message: '', data });
    } catch (err) {
      next(err);
    }
  },

  /**
   * Update — 更新 Todo
   * -----------------------------------------------
   * 对标 Go:
   *   func (c *cTodo) Update(ctx context.Context, req *v1.TodoUpdateReq) (*v1.TodoUpdateRes, error) {
   *       return nil, service.Todo().Update(ctx, req.Id, req.Title, req.Completed)
   *   }
   *
   * 路由: PUT /api/v1/todos/:id
   * 对标 Go Meta: `path:"/todos/:id" method:"put"`
   * -----------------------------------------------
   */
  async update(req, res, next) {
    try {
      const { errors, id } = validateUpdateReq(req.params, req.body);
      if (errors.length > 0) {
        return res.json({ code: 51, message: errors.join('; '), data: null });
      }

      // body.completed 可能为 undefined（未传）或 false（显式传值）
      // 这里用 'completed' in req.body 判断，对标 Go 的 *bool 指针
      const completed = 'completed' in req.body ? req.body.completed : undefined;
      await todoService.update(id, req.body.title, completed);

      // 返回空对象 — 对标 TodoUpdateRes{}
      res.json({ code: 0, message: '', data: null });
    } catch (err) {
      next(err);
    }
  },

  /**
   * Delete — 删除 Todo
   * -----------------------------------------------
   * 对标 Go:
   *   func (c *cTodo) Delete(ctx context.Context, req *v1.TodoDeleteReq) (*v1.TodoDeleteRes, error) {
   *       return nil, service.Todo().Delete(ctx, req.Id)
   *   }
   *
   * 路由: DELETE /api/v1/todos/:id
   * 对标 Go Meta: `path:"/todos/:id" method:"delete"`
   * -----------------------------------------------
   */
  async delete(req, res, next) {
    try {
      const { errors, id } = validateDeleteReq(req.params);
      if (errors.length > 0) {
        return res.json({ code: 51, message: errors.join('; '), data: null });
      }

      await todoService.delete(id);

      // 返回空对象 — 对标 TodoDeleteRes{}
      res.json({ code: 0, message: '', data: null });
    } catch (err) {
      next(err);
    }
  },
};

module.exports = todoController;
