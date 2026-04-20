# gRPC User Service

> Go gRPC 微服务模板项目，实现用户 CRUD，SQLite 存储，集成日志和恢复拦截器。

## 项目结构

```
grpc-user/
├── proto/
│   ├── user.proto              # Protobuf 服务定义
│   └── userpb/
│       ├── user.pb.go          # 生成的消息代码
│       └── user_grpc.pb.go     # 生成的 gRPC 桩代码
├── cmd/
│   ├── server/main.go          # 服务端入口
│   └── client/main.go          # 客户端调用示例
├── internal/
│   ├── database/db.go          # SQLite 初始化
│   ├── service/user.go         # UserService CRUD 实现
│   └── interceptor/interceptor.go  # 日志 + Panic恢复 拦截器
├── go.mod
└── README.md
```

## 快速启动

```bash
# 1. 安装依赖
go mod tidy

# 2. 启动服务端
go run cmd/server/main.go

# 3. 另一个终端，运行客户端
go run cmd/client/main.go
```

## 重新生成 protobuf 代码

```bash
protoc --proto_path=proto \
  --go_out=proto/userpb --go_opt=paths=source_relative \
  --go-grpc_out=proto/userpb --go-grpc_opt=paths=source_relative \
  user.proto
```

## gRPC 接口

| RPC 方法 | 说明 |
|----------|------|
| `CreateUser` | 创建用户 |
| `GetUser` | 查询单个用户 |
| `ListUsers` | 分页查询用户列表 |
| `UpdateUser` | 更新用户 |
| `DeleteUser` | 删除用户 |
