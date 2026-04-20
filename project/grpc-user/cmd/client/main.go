// cmd/client/main.go
// gRPC 客户端调用示例 — 演示完整 CRUD 流程
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "grpc-user/proto/userpb"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "gRPC 服务端地址")
	flag.Parse()

	// 建立连接
	conn, err := grpc.NewClient(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接失败: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println()
	fmt.Println("  ═══════════ gRPC Client 调用示例 ═══════════")

	// 1. 创建用户
	fmt.Println("\n  ── 1. CreateUser ──")
	createResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "张三",
		Email: "zhangsan@example.com",
	})
	if err != nil {
		fmt.Printf("  ✗ CreateUser 失败: %v\n", err)
	} else {
		printUser("  创建成功", createResp.User)
	}

	createResp2, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "李四",
		Email: "lisi@example.com",
	})
	if err != nil {
		fmt.Printf("  ✗ CreateUser 失败: %v\n", err)
	} else {
		printUser("  创建成功", createResp2.User)
	}

	// 2. 查询单个
	fmt.Println("\n  ── 2. GetUser ──")
	if createResp != nil {
		getResp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: createResp.User.Id})
		if err != nil {
			fmt.Printf("  ✗ GetUser 失败: %v\n", err)
		} else {
			printUser("  查询成功", getResp.User)
		}
	}

	// 3. 列表查询
	fmt.Println("\n  ── 3. ListUsers ──")
	listResp, err := client.ListUsers(ctx, &pb.ListUsersRequest{Page: 1, Size: 10})
	if err != nil {
		fmt.Printf("  ✗ ListUsers 失败: %v\n", err)
	} else {
		fmt.Printf("  总数: %d\n", listResp.Total)
		for _, u := range listResp.Users {
			printUser("  ", u)
		}
	}

	// 4. 更新用户
	fmt.Println("\n  ── 4. UpdateUser ──")
	if createResp != nil {
		updateResp, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
			Id:    createResp.User.Id,
			Name:  "张三丰",
			Email: "zhangsanfeng@example.com",
		})
		if err != nil {
			fmt.Printf("  ✗ UpdateUser 失败: %v\n", err)
		} else {
			printUser("  更新成功", updateResp.User)
		}
	}

	// 5. 删除用户
	fmt.Println("\n  ── 5. DeleteUser ──")
	if createResp2 != nil {
		delResp, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: createResp2.User.Id})
		if err != nil {
			fmt.Printf("  ✗ DeleteUser 失败: %v\n", err)
		} else {
			fmt.Printf("  删除成功: %v\n", delResp.Success)
		}
	}

	// 6. 验证删除后列表
	fmt.Println("\n  ── 6. ListUsers (删除后) ──")
	listResp2, err := client.ListUsers(ctx, &pb.ListUsersRequest{Page: 1, Size: 10})
	if err != nil {
		fmt.Printf("  ✗ ListUsers 失败: %v\n", err)
	} else {
		fmt.Printf("  总数: %d\n", listResp2.Total)
		for _, u := range listResp2.Users {
			printUser("  ", u)
		}
	}

	fmt.Println("\n  ═══════════════════════════════════════════")
	fmt.Println()
}

func printUser(prefix string, u *pb.User) {
	fmt.Printf("%s [id=%d] %s <%s>\n", prefix, u.Id, u.Name, u.Email)
}
