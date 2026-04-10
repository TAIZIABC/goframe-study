// internal/controller/todo.go
package controller

import (
    "context"
    v1 "todo-api/api/v1"
    "todo-api/internal/service"
)

var Todo = &cTodo{}
type cTodo struct{}

func (c *cTodo) Create(ctx context.Context, req *v1.TodoCreateReq) (res *v1.TodoCreateRes, err error) {
    id, err := service.Todo().Create(ctx, req.Title, req.Completed)
    if err != nil {
        return nil, err
    }
    return &v1.TodoCreateRes{Id: id}, nil
}

func (c *cTodo) List(ctx context.Context, req *v1.TodoListReq) (res *v1.TodoListRes, err error) {
    list, total, err := service.Todo().List(ctx, req.Page, req.Size)
    if err != nil {
        return nil, err
    }
    return &v1.TodoListRes{List: list, Total: total}, nil
}

func (c *cTodo) Update(ctx context.Context, req *v1.TodoUpdateReq) (res *v1.TodoUpdateRes, err error) {
    return nil, service.Todo().Update(ctx, req.Id, req.Title, req.Completed)
}

func (c *cTodo) Delete(ctx context.Context, req *v1.TodoDeleteReq) (res *v1.TodoDeleteRes, err error) {
    return nil, service.Todo().Delete(ctx, req.Id)
}