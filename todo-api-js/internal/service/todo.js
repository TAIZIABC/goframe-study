/**
 * Todo Service 层 — 业务逻辑
 * -----------------------------------------------
 * 对标 GoFrame: internal/service/todo.go
 *
 * GoFrame 使用单例模式 + ORM：
 *   g.Model("todos").Ctx(ctx).Data(...).Insert()
 *   g.Model("todos").Ctx(ctx).Page(page, size).OrderDesc("id").Scan(&list)
 *
 * JS 版本使用 mysql2 直接执行 SQL，功能完全一致
 * -----------------------------------------------
 */

const { query, execute } = require('../database/db');

/**
 * 单例模式（对标 Go: var todoInstance = &sTodo{}）
 * JS 中 module 本身就是单例，require 时只会执行一次
 */
const todoService = {
  /**
   * Create 创建一个 Todo
   * -----------------------------------------------
   * 对标 Go:
   *   func (s *sTodo) Create(ctx context.Context, title string, completed bool) (int, error) {
   *       result, err := g.Model("todos").Ctx(ctx).Data(g.Map{
   *           "title":     title,
   *           "completed": completed,
   *       }).Insert()
   *       id, _ := result.LastInsertId()
   *       return int(id), nil
   *   }
   * -----------------------------------------------
   */
  async create(title, completed = false) {
    const result = await execute(
      'INSERT INTO `todos` (`title`, `completed`) VALUES (?, ?)',
      [title, completed ? 1 : 0]
    );
    return result.insertId; // 对标 result.LastInsertId()
  },

  /**
   * List 分页查询 Todo 列表
   * -----------------------------------------------
   * 对标 Go:
   *   func (s *sTodo) List(ctx context.Context, page, size int) ([]v1.TodoItem, int, error) {
   *       total, err := model.Count()
   *       err = model.Page(page, size).OrderDesc("id").Scan(&list)
   *       return list, total, nil
   *   }
   * -----------------------------------------------
   */
  async list(page = 1, size = 10) {
    // 查询总数 — 对标 model.Count()
    const countResult = await query('SELECT COUNT(*) AS total FROM `todos`');
    const total = countResult[0].total;

    // 分页查询 — 对标 model.Page(page, size).OrderDesc("id").Scan(&list)
    const offset = (page - 1) * size;
    const list = await query(
      'SELECT `id`, `title`, `completed` FROM `todos` ORDER BY `id` DESC LIMIT ? OFFSET ?',
      [size, offset]
    );

    // completed 字段：MySQL 存储为 TINYINT(1)，需转为布尔值
    const items = list.map((item) => ({
      id: item.id,
      title: item.title,
      completed: item.completed === 1,
    }));

    return { list: items, total };
  },

  /**
   * Update 更新 Todo（支持部分更新）
   * -----------------------------------------------
   * 对标 Go:
   *   func (s *sTodo) Update(ctx context.Context, id int, title string, completed *bool) error {
   *       data := g.Map{}
   *       if title != "" { data["title"] = title }
   *       if completed != nil { data["completed"] = *completed }
   *       if len(data) == 0 { return nil }
   *       _, err := g.Model("todos").Where("id", id).Data(data).Update()
   *       return err
   *   }
   *
   * 注意：Go 版本用 *bool 指针区分"未传值"和"传了false"
   * JS 版本用 undefined 判断实现相同效果
   * -----------------------------------------------
   */
  async update(id, title, completed) {
    // 动态构建更新字段 — 对标 data := g.Map{}
    const fields = [];
    const params = [];

    if (title !== undefined && title !== '') {
      fields.push('`title` = ?');
      params.push(title);
    }
    if (completed !== undefined) {
      fields.push('`completed` = ?');
      params.push(completed ? 1 : 0);
    }

    // 没有需要更新的字段，直接返回 — 对标 if len(data) == 0 { return nil }
    if (fields.length === 0) {
      return;
    }

    params.push(id);
    await execute(
      `UPDATE \`todos\` SET ${fields.join(', ')} WHERE \`id\` = ?`,
      params
    );
  },

  /**
   * Delete 删除 Todo
   * -----------------------------------------------
   * 对标 Go:
   *   func (s *sTodo) Delete(ctx context.Context, id int) error {
   *       _, err := g.Model("todos").Where("id", id).Delete()
   *       return err
   *   }
   * -----------------------------------------------
   */
  async delete(id) {
    await execute('DELETE FROM `todos` WHERE `id` = ?', [id]);
  },
};

module.exports = todoService;
