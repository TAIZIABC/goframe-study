// internal/service/user.go
// UserService gRPC 服务实现
package service

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-user/internal/database"
	pb "grpc-user/proto/userpb"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{}
}

// CreateUser 创建用户
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name 不能为空")
	}

	result, err := database.DB.ExecContext(ctx,
		"INSERT INTO users (name, email) VALUES (?, ?)",
		req.Name, req.Email,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建失败: %v", err)
	}

	id, _ := result.LastInsertId()

	user, err := s.getUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{User: user}, nil
}

// GetUser 查询单个用户
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id 无效")
	}

	user, err := s.getUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{User: user}, nil
}

// ListUsers 分页查询用户列表
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	page := req.Page
	size := req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	// 查总数
	var total int32
	err := database.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * size
	rows, err := database.DB.QueryContext(ctx,
		"SELECT id, name, email, created_at, updated_at FROM users ORDER BY id DESC LIMIT ? OFFSET ?",
		size, offset,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询失败: %v", err)
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		u := &pb.User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, status.Errorf(codes.Internal, "扫描行失败: %v", err)
		}
		users = append(users, u)
	}

	return &pb.ListUsersResponse{Users: users, Total: total}, nil
}

// UpdateUser 更新用户
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id 无效")
	}

	// 动态构建 SQL
	sets := []string{}
	args := []interface{}{}
	if req.Name != "" {
		sets = append(sets, "name = ?")
		args = append(args, req.Name)
	}
	if req.Email != "" {
		sets = append(sets, "email = ?")
		args = append(args, req.Email)
	}
	if len(sets) == 0 {
		return nil, status.Error(codes.InvalidArgument, "没有需要更新的字段")
	}

	sets = append(sets, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, req.Id)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", joinStrings(sets, ", "))
	result, err := database.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	user, err := s.getUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{User: user}, nil
}

// DeleteUser 删除用户
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id 无效")
	}

	result, err := database.DB.ExecContext(ctx, "DELETE FROM users WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "删除失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	return &pb.DeleteUserResponse{Success: true}, nil
}

// getUserByID 内部方法
func (s *UserServiceServer) getUserByID(ctx context.Context, id int64) (*pb.User, error) {
	u := &pb.User{}
	err := database.DB.QueryRowContext(ctx,
		"SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?", id,
	).Scan(&u.Id, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询失败: %v", err)
	}
	return u, nil
}

func joinStrings(strs []string, sep string) string {
	r := ""
	for i, s := range strs {
		if i > 0 {
			r += sep
		}
		r += s
	}
	return r
}
