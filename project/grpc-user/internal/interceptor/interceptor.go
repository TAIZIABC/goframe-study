// internal/interceptor/interceptor.go
// gRPC 拦截器：请求日志 + Panic 恢复
package interceptor

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor 请求日志拦截器
func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// 调用实际 handler
	resp, err := handler(ctx, req)

	duration := time.Since(start)
	code := codes.OK
	if err != nil {
		if s, ok := status.FromError(err); ok {
			code = s.Code()
		}
	}

	// 日志着色
	color := "\033[32m" // 绿
	if code != codes.OK {
		color = "\033[31m" // 红
	}

	fmt.Printf("  %s %s%-4s\033[0m %-40s %s\n",
		time.Now().Format("15:04:05"),
		color, code,
		info.FullMethod,
		duration.Round(time.Microsecond),
	)

	return resp, err
}

// RecoveryInterceptor Panic 恢复拦截器
func RecoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  \033[31m[PANIC]\033[0m %s: %v\n%s\n",
				info.FullMethod, r, string(debug.Stack()))
			err = status.Errorf(codes.Internal, "服务内部错误")
		}
	}()
	return handler(ctx, req)
}
