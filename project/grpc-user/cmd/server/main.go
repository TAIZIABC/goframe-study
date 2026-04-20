// cmd/server/main.go
// gRPC 服务端入口
package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	"grpc-user/internal/database"
	"grpc-user/internal/interceptor"
	"grpc-user/internal/service"
	pb "grpc-user/proto/userpb"
)

func main() {
	port := flag.Int("port", 50051, "gRPC 服务端口")
	dbPath := flag.String("db", "users.db", "SQLite 数据库路径")
	flag.Parse()

	// 初始化数据库
	if err := database.Init(*dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ 数据库初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 监听端口
	addr := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ 监听失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 gRPC Server，注册拦截器链
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryInterceptor, // 先恢复 panic
			interceptor.LoggingInterceptor,  // 再记日志
		),
	)

	// 注册用户服务
	pb.RegisterUserServiceServer(grpcServer, service.NewUserServiceServer())

	fmt.Println()
	fmt.Println("  ┌──────────────────────────────────────┐")
	fmt.Println("  │     gRPC User Service                │")
	fmt.Println("  ├──────────────────────────────────────┤")
	fmt.Printf("  │  🚀 监听地址:  %s\n", addr)
	fmt.Printf("  │  💾 数据库:    %s\n", *dbPath)
	fmt.Println("  │  🔗 拦截器:   日志 + Panic恢复")
	fmt.Println("  │  按 Ctrl+C 停止")
	fmt.Println("  └──────────────────────────────────────┘")
	fmt.Println()

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ 服务启动失败: %v\n", err)
		os.Exit(1)
	}
}
