package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "todo-api/api/v1"
)

type sTodo struct{}

var todoInstance = &sTodo{}

// Todo 返回 service 单例
func Todo() *sTodo {
	return todoInstance
}

// Create 创建一个 Todo
func (s *sTodo) Create(ctx context.Context, title string, completed bool) (int, error) {
	result, err := g.Model("todos").Ctx(ctx).Data(g.Map{
		"title":     title,
		"completed": completed,
	}).Insert()
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

// List 分页查询 Todo 列表
func (s *sTodo) List(ctx context.Context, page, size int) ([]v1.TodoItem, int, error) {
	model := g.Model("todos").Ctx(ctx)

	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []v1.TodoItem
	err = model.Page(page, size).
		OrderDesc("id").
		Scan(&list)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// Update 更新 Todo
func (s *sTodo) Update(ctx context.Context, id int, title string, completed *bool) error {
	data := g.Map{}
	if title != "" {
		data["title"] = title
	}
	if completed != nil {
		data["completed"] = *completed
	}
	if len(data) == 0 {
		return nil
	}
	_, err := g.Model("todos").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// Delete 删除 Todo
func (s *sTodo) Delete(ctx context.Context, id int) error {
	_, err := g.Model("todos").Ctx(ctx).Where("id", id).Delete()
	return err
}
