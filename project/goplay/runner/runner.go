// runner/runner.go
// 代码沙箱运行器：临时目录隔离、编译执行、超时控制
package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Result struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
	Duration int64  `json:"duration_ms"`
	Error    string `json:"error,omitempty"`
	Success  bool   `json:"success"`
}

type Config struct {
	Timeout    time.Duration
	MaxOutput  int // 最大输出字节数
}

func DefaultConfig() *Config {
	return &Config{
		Timeout:   10 * time.Second,
		MaxOutput: 64 * 1024, // 64KB
	}
}

// Run 在隔离临时目录中编译并运行 Go 代码
func Run(code string, cfg *Config) *Result {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	start := time.Now()

	// 基本安全检查
	if err := sanitize(code); err != nil {
		return &Result{Error: err.Error(), Duration: time.Since(start).Milliseconds()}
	}

	// 创建临时目录（使用 HOME 下的目录，避免 macOS 系统 temp root 被 Go 拒绝）
	home, _ := os.UserHomeDir()
	baseDir := filepath.Join(home, ".goplay", "sandbox")
	os.MkdirAll(baseDir, 0755)
	tmpDir, err := os.MkdirTemp(baseDir, "run-*")
	if err != nil {
		return &Result{Error: "创建临时目录失败: " + err.Error(), Duration: time.Since(start).Milliseconds()}
	}
	defer os.RemoveAll(tmpDir)

	// 写入 main.go
	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(code), 0644); err != nil {
		return &Result{Error: "写入文件失败: " + err.Error(), Duration: time.Since(start).Milliseconds()}
	}

	// 初始化 go module
	modFile := filepath.Join(tmpDir, "go.mod")
	os.WriteFile(modFile, []byte("module playground\n\ngo 1.23.0\n"), 0644)

	// 带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// 编译
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", "main", ".")
	buildCmd.Dir = tmpDir
	buildCmd.Env = sandboxEnv(tmpDir)

	buildOut, buildErr := buildCmd.CombinedOutput()
	if buildErr != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return &Result{Error: "编译超时", Stderr: string(buildOut), Duration: time.Since(start).Milliseconds()}
		}
		return &Result{
			Stderr:   cleanOutput(string(buildOut), tmpDir),
			ExitCode: 1,
			Duration: time.Since(start).Milliseconds(),
			Error:    "编译失败",
		}
	}

	// 运行
	runCmd := exec.CommandContext(ctx, filepath.Join(tmpDir, "main"))
	runCmd.Dir = tmpDir
	runCmd.Env = sandboxEnv(tmpDir)

	var stdout, stderr strings.Builder
	runCmd.Stdout = &limitWriter{w: &stdout, max: cfg.MaxOutput}
	runCmd.Stderr = &limitWriter{w: &stderr, max: cfg.MaxOutput}

	runErr := runCmd.Run()
	duration := time.Since(start)

	result := &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration.Milliseconds(),
		Success:  true,
	}

	if runErr != nil {
		result.Success = false
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Sprintf("执行超时 (限制 %s)", cfg.Timeout)
			result.ExitCode = -1
		} else if exitErr, ok := runErr.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.Error = runErr.Error()
			result.ExitCode = -1
		}
	}

	return result
}

// sanitize 基本安全检查
func sanitize(code string) error {
	if len(code) > 100*1024 {
		return fmt.Errorf("代码长度超过限制 (最大 100KB)")
	}

	dangerous := []string{
		"os/exec", "syscall", "unsafe",
		"net/http", "plugin",
	}
	for _, pkg := range dangerous {
		if strings.Contains(code, `"`+pkg+`"`) {
			return fmt.Errorf("禁止导入包: %s", pkg)
		}
	}

	return nil
}

// sandboxEnv 构建沙箱环境变量
func sandboxEnv(tmpDir string) []string {
	home, _ := os.UserHomeDir()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(home, "go")
	}
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		// 尝试自动获取
		out, err := exec.Command("go", "env", "GOROOT").Output()
		if err == nil {
			goroot = strings.TrimSpace(string(out))
		}
	}

	tmpSubDir := filepath.Join(tmpDir, "tmp")
	os.MkdirAll(tmpSubDir, 0755)

	return []string{
		"HOME=" + home,
		"TMPDIR=" + tmpSubDir,
		"GOPATH=" + gopath,
		"GOROOT=" + goroot,
		"GOCACHE=" + filepath.Join(tmpSubDir, "cache"),
		"GOMODCACHE=" + filepath.Join(gopath, "pkg", "mod"),
		"PATH=" + os.Getenv("PATH"),
		"GOMEMLIMIT=128MiB",
		"GOFLAGS=-mod=mod",
		"GONOSUMCHECK=*",
		"GONOSUMDB=*",
	}
}

// cleanOutput 清理输出中的临时路径
func cleanOutput(s, tmpDir string) string {
	return strings.ReplaceAll(s, tmpDir+"/", "")
}

// limitWriter 限制写入量
type limitWriter struct {
	w       *strings.Builder
	max     int
	written int
}

func (lw *limitWriter) Write(p []byte) (int, error) {
	remaining := lw.max - lw.written
	if remaining <= 0 {
		return len(p), nil // 静默丢弃
	}
	if len(p) > remaining {
		p = p[:remaining]
		lw.w.Write(p)
		lw.w.WriteString("\n... (输出被截断)")
		lw.written = lw.max
		return len(p), nil
	}
	n, err := lw.w.Write(p)
	lw.written += n
	return n, err
}
