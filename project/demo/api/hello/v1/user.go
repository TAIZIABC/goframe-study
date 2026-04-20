package v1

import "github.com/gogf/gf/v2/frame/g"

type CreateUserReq struct {
	g.Meta `path:"/user" tags:"User" method:"get" summary:"create user"`
	Name   string `json:"name" v:"required#name不能为空" dc:"用户名"`
	Age    int    `json:"age" v:"required#age不能为空" dc:"年龄"`
}

type CreateUserRes struct {
	ID int `json:"id" dc:"用户ID"`
}

type DeleteUserReq struct {
	g.Meta `path:"/user{id}" tags:"User" method:"delete" summary:"delete user"`
	Id     int `json:"id" v:"required#id不能为空" dc:"用户ID"`
}

type DeleteUserRes struct {
	Msg string `json:"msg" dc:"消息"`
}
