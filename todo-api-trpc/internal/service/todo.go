// internal/service/todo.go
// -----------------------------------------------
// Todo Service 层 — 业务逻辑
// 对标 GoFrame: internal/service/todo.go
//
// GoFrame 使用内置 ORM 链式调用：
//   g.Model("todos").Ctx(ctx).Data(...).Insert()
//   g.Model("todos").Ctx(ctx).Page(page, size).OrderDesc("id").Scan(&list)
//
// tRPC-Go 没有内置 ORM，使用 database/sql 原生操作
// 功能完全一致，只是操作方式不同
// -----------------------------------------------
package service

import (
	"context"
	"fmt"

	v1 "todo-api-trpc/api/v1"
	"todo-api-trpc/internal/database"
)

// sTodo 服务实例
// 对标 GoFrame: type sTodo struct{}
type sTodo struct{}

// 单例模式
// 对标 GoFrame: var todoInstance = &sTodo{}
var todoInstance = &sTodo{}

// Todo 返回 service 单例
// 对标 GoFrame: func Todo() *sTodo { return todoInstance }
func Todo() *sTodo {
	return todoInstance
}

// Create 创建一个 Todo
// -----------------------------------------------
// 对标 GoFrame:
//
//	func (s *sTodo) Create(ctx context.Context, title string, completed bool) (int, error) {
//	    result, err := g.Model("todos").Ctx(ctx).Data(g.Map{
//	        "title":     title,
//	        "completed": completed,
//	    }).Insert()
//	    id, _ := result.LastInsertId()
//	    return int(id), nil
//	}
//
// tRPC-Go 版使用原生 SQL 代替 GoFrame ORM
// -----------------------------------------------
func (s *sTodo) Create(ctx context.Context, title string, completed bool) (int, error) {
	db := database.GetDB()
	completedInt := 0
	if completed {
		completedInt = 1
	}

	result, err := db.ExecContext(ctx,
		"INSERT INTO `todos` (`title`, `completed`) VALUES (?, ?)",
		title, completedInt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// List 分页查询 Todo 列表
// -----------------------------------------------
// 对标 GoFrame:
//
//	func (s *sTodo) List(ctx context.Context, page, size int) ([]v1.TodoItem, int, error) {
//	    total, err := model.Count()
//	    err = model.Page(page, size).OrderDesc("id").Scan(&list)
//	    return list, total, nil
//	}
//
// GoFrame 的 Page(page, size) 内部自动计算 OFFSET
// 这里手动计算: offset = (page - 1) * size
// -----------------------------------------------
func (s *sTodo) List(ctx context.Context, page, size int) ([]v1.TodoItem, int, error) {
	db := database.GetDB()

	// 查询总数 — 对标 model.Count()
	var total int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM `todos`").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询 — 对标 model.Page(page, size).OrderDesc("id").Scan(&list)
	offset := (page - 1) * size
	rows, err := db.QueryContext(ctx,
		"SELECT `id`, `title`, `completed` FROM `todos` ORDER BY `id` DESC LIMIT ? OFFSET ?",
		size, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []v1.TodoItem
	for rows.Next() {
		var item v1.TodoItem
		var completedInt int
		if err := rows.Scan(&item.Id, &item.Title, &completedInt); err != nil {
			return nil, 0, err
		}
		item.Completed = completedInt == 1
		list = append(list, item)
	}

	if list == nil {
		list = []v1.TodoItem{} // 保证返回空数组而非 null
	}

	return list, total, nil
}

// Update 更新 Todo（支持部分更新）
// -----------------------------------------------
// 对标 GoFrame:
//
//	func (s *sTodo) Update(ctx context.Context, id int, title string, completed *bool) error {
//	    data := g.Map{}
//	    if title != "" { data["title"] = title }
//	    if completed != nil { data["completed"] = *completed }
//	    if len(data) == 0 { return nil }
//	    _, err := g.Model("todos").Where("id", id).Data(data).Update()
//	    return err
//	}
//
// tRPC-Go 版手动拼接 SQL 的 SET 子句实现动态更新
// -----------------------------------------------
func (s *sTodo) Update(ctx context.Context, id int, title string, completed *bool) error {
	db := database.GetDB()

	// 动态构建更新字段 — 对标 data := g.Map{}
	setClauses := []string{}
	args := []interface{}{}

	if title != "" {
		setClauses = append(setClauses, "`title` = ?")
		args = append(args, title)
	}
	if completed != nil {
		setClauses = append(setClauses, "`completed` = ?")
		completedInt := 0
		if *completed {
			completedInt = 1
		}
		args = append(args, completedInt)
	}

	// 没有需要更新的字段 — 对标 if len(data) == 0 { return nil }
	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE `todos` SET %s WHERE `id` = ?",
		joinStrings(setClauses, ", "))
	args = append(args, id)

	_, err := db.ExecContext(ctx, query, args...)
	return err
}

// Delete 删除 Todo
// -----------------------------------------------
// 对标 GoFrame:
//
//	func (s *sTodo) Delete(ctx context.Context, id int) error {
//	    _, err := g.Model("todos").Where("id", id).Delete()
//	    return err
//	}
// -----------------------------------------------
func (s *sTodo) Delete(ctx context.Context, id int) error {
	db := database.GetDB()
	_, err := db.ExecContext(ctx, "DELETE FROM `todos` WHERE `id` = ?", id)
	return err
}

// joinStrings 简单的字符串拼接
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
