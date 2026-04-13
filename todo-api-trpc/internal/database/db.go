// internal/database/db.go
// -----------------------------------------------
// 数据库连接模块
// 对标 GoFrame: g.DB() / g.Model("todos")
//
// GoFrame 通过配置文件自动创建数据库连接池：
//   database:
//     default:
//       link: "mysql:root:xxx@tcp(...)/todo_db"
//   使用: g.Model("todos") 即可操作
//
// tRPC-Go 没有内置 ORM，需要手动管理数据库连接
// 这里使用 database/sql + go-sql-driver/mysql
// -----------------------------------------------
package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

// Config 数据库配置
// 对标 GoFrame config.yaml 中的 database 配置段
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// DefaultConfig 返回默认数据库配置（从环境变量读取）
// 对标 GoFrame: manifest/config/config.yaml 中的 database.default.link
func DefaultConfig() *Config {
	port := 3306
	if v := os.Getenv("DB_PORT"); v != "" {
		if p, err := fmt.Sscanf(v, "%d", &port); p == 0 || err != nil {
			port = 3306
		}
	}
	return &Config{
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     port,
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_NAME", "todo_db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Init 初始化数据库连接
// 对标 GoFrame: 框架启动时自动根据配置初始化数据库
func Init(cfg *Config) error {
	var initErr error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			initErr = fmt.Errorf("open database failed: %w", err)
			return
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)

		if err = db.Ping(); err != nil {
			initErr = fmt.Errorf("ping database failed: %w", err)
			return
		}
	})
	return initErr
}

// GetDB 获取数据库连接（单例）
// 对标 GoFrame: g.DB()
func GetDB() *sql.DB {
	return db
}
