// quiz-arena 后端 API 服务
// 前后端分离：仅提供 JSON API
// 题库从 data/questions.json 和 data/coding.json 加载
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 难度等级
const (
	DifficultyEasy   = 1 // 简单
	DifficultyMedium = 2 // 中等
	DifficultyHard   = 3 // 困难
)

// Question 题目结构
type Question struct {
	ID         int      `json:"id"`
	Category   string   `json:"category"`
	Difficulty int      `json:"difficulty"` // 1=简单 2=中等 3=困难
	Title      string   `json:"title"`
	Code       string   `json:"code"`
	Options    []string `json:"options"`
	Answer     int      `json:"-"` // 答案不返回前端
	Explain    string   `json:"-"`
}

// SubmitReq 提交答案请求
type SubmitReq struct {
	ID     int `json:"id"`
	Choice int `json:"choice"`
}

// SubmitResp 提交答案响应
type SubmitResp struct {
	Correct bool   `json:"correct"`
	Answer  int    `json:"answer"`
	Explain string `json:"explain"`
}

// QuestionView 前端展示用
type QuestionView struct {
	ID         int      `json:"id"`
	Category   string   `json:"category"`
	Difficulty int      `json:"difficulty"`
	Title      string   `json:"title"`
	Code       string   `json:"code"`
	Options    []string `json:"options"`
	Total      int      `json:"total"`
	Index      int      `json:"index"`
}

// CategoryInfo 分类统计信息
type CategoryInfo struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// ========== 代码编写题 ==========

// CodingQuestion 编码题结构
type CodingQuestion struct {
	ID         int    `json:"id"`
	Category   string `json:"category"`
	Difficulty int    `json:"difficulty"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Template   string `json:"template"`
	TestCode   string `json:"-"`
	Solution   string `json:"-"`
	Hint       string `json:"hint"`
}

// CodingQuestionView 前端展示用
type CodingQuestionView struct {
	ID         int    `json:"id"`
	Category   string `json:"category"`
	Difficulty int    `json:"difficulty"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Template   string `json:"template"`
	Hint       string `json:"hint"`
	Total      int    `json:"total"`
	Index      int    `json:"index"`
}

// CodingRunReq 提交代码请求
type CodingRunReq struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

// CodingRunResp 执行结果
type CodingRunResp struct {
	Pass     bool   `json:"pass"`
	Output   string `json:"output"`
	Error    string `json:"error"`
	Solution string `json:"solution"`
}

// ========== 题库数据（从 JSON 加载） ==========

var questions []Question
var codingQuestions []CodingQuestion

// ========== JSON 序列化结构体（包含 answer/explain/testCode/solution） ==========

type QuestionJSON struct {
	ID         int      `json:"id"`
	Category   string   `json:"category"`
	Difficulty int      `json:"difficulty"`
	Title      string   `json:"title"`
	Code       string   `json:"code,omitempty"`
	Options    []string `json:"options"`
	Answer     int      `json:"answer"`
	Explain    string   `json:"explain"`
}

type CodingQuestionJSON struct {
	ID         int    `json:"id"`
	Category   string `json:"category"`
	Difficulty int    `json:"difficulty"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Template   string `json:"template"`
	TestCode   string `json:"testCode"`
	Solution   string `json:"solution"`
	Hint       string `json:"hint"`
}

// ========== 题库加载 ==========

// loadQuestionsFromJSON 从 JSON 文件加载题库
func loadQuestionsFromJSON(dataDir string) {
	// 加载选择题
	qPath := filepath.Join(dataDir, "questions.json")
	if data, err := os.ReadFile(qPath); err == nil {
		var qj []QuestionJSON
		if err := json.Unmarshal(data, &qj); err != nil {
			log.Fatalf("❌ 解析 %s 失败: %v", qPath, err)
		}
		questions = make([]Question, len(qj))
		for i, q := range qj {
			questions[i] = Question{
				ID: q.ID, Category: q.Category, Difficulty: q.Difficulty,
				Title: q.Title, Code: q.Code, Options: q.Options,
				Answer: q.Answer, Explain: q.Explain,
			}
		}
		log.Printf("📄 从 %s 加载 %d 道选择题", qPath, len(questions))
	} else {
		log.Fatalf("❌ 找不到题库文件 %s: %v\n   请先运行: go run main.go -export", qPath, err)
	}

	// 加载编码题
	cPath := filepath.Join(dataDir, "coding.json")
	if data, err := os.ReadFile(cPath); err == nil {
		var cj []CodingQuestionJSON
		if err := json.Unmarshal(data, &cj); err != nil {
			log.Fatalf("❌ 解析 %s 失败: %v", cPath, err)
		}
		codingQuestions = make([]CodingQuestion, len(cj))
		for i, q := range cj {
			codingQuestions[i] = CodingQuestion{
				ID: q.ID, Category: q.Category, Difficulty: q.Difficulty,
				Title: q.Title, Desc: q.Desc, Template: q.Template,
				TestCode: q.TestCode, Solution: q.Solution, Hint: q.Hint,
			}
		}
		log.Printf("📄 从 %s 加载 %d 道编码题", cPath, len(codingQuestions))
	} else {
		log.Fatalf("❌ 找不到题库文件 %s: %v\n   请先运行: go run main.go -export", cPath, err)
	}
}

// ========== 选择题 API ==========

// filterQuestions 根据查询条件过滤题目
func filterQuestions(category string, difficulty int, ids []int) []QuestionView {
	result := make([]QuestionView, 0)
	idSet := make(map[int]bool)
	for _, id := range ids {
		idSet[id] = true
	}

	for i, q := range questions {
		if len(ids) > 0 && !idSet[q.ID] {
			continue
		}
		if category != "" && category != "all" && q.Category != category {
			continue
		}
		if difficulty > 0 && q.Difficulty != difficulty {
			continue
		}
		result = append(result, QuestionView{
			ID:         q.ID,
			Category:   q.Category,
			Difficulty: q.Difficulty,
			Title:      q.Title,
			Code:       q.Code,
			Options:    q.Options,
			Total:      len(questions),
			Index:      i,
		})
	}
	return result
}

// parseIDs 解析逗号分隔的 ID 列表
func parseIDs(s string) []int {
	if s == "" {
		return nil
	}
	ids := make([]int, 0)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if i > start {
				if n, err := strconv.Atoi(s[start:i]); err == nil {
					ids = append(ids, n)
				}
			}
			start = i + 1
		}
	}
	return ids
}

// handleQuestions GET /api/questions?category=&difficulty=&ids=
func handleQuestions(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	difficulty, _ := strconv.Atoi(r.URL.Query().Get("difficulty"))
	ids := parseIDs(r.URL.Query().Get("ids"))

	result := filterQuestions(category, difficulty, ids)
	for i := range result {
		result[i].Index = i
		result[i].Total = len(result)
	}
	writeJSON(w, http.StatusOK, result)
}

// handleCategories GET /api/categories
func handleCategories(w http.ResponseWriter, r *http.Request) {
	counts := make(map[string]int)
	for _, q := range questions {
		counts[q.Category]++
	}
	order := []string{
		"Goroutine", "Channel", "Select", "WaitGroup", "Mutex", "Atomic", "Context",
		"切片", "Map", "字符串", "defer", "接口", "错误处理", "指针", "结构体", "泛型", "语法",
		"GoFrame",
		"综合",
	}
	result := make([]CategoryInfo, 0, len(order))
	for _, name := range order {
		if c, ok := counts[name]; ok {
			result = append(result, CategoryInfo{Name: name, Count: c})
		}
	}
	writeJSON(w, http.StatusOK, result)
}

// handleSubmit POST /api/submit
func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SubmitReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var q *Question
	for i := range questions {
		if questions[i].ID == req.ID {
			q = &questions[i]
			break
		}
	}
	if q == nil {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, SubmitResp{
		Correct: req.Choice == q.Answer,
		Answer:  q.Answer,
		Explain: q.Explain,
	})
}

// handleImage GET /api/image
func handleImage(w http.ResponseWriter, r *http.Request) {
	const upstream = "http://127.0.0.1:12150/v2/preview?project=%2FUsers%2Fkingjungle%2FDocuments%2Fwork%2FAppletNew"

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, upstream, nil)
	if err != nil {
		http.Error(w, "build upstream request failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "fetch upstream failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, "upstream status "+resp.Status+": "+string(body), resp.StatusCode)
		return
	}

	if ct := resp.Header.Get("Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	} else {
		w.Header().Set("Content-Type", "image/png")
	}
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		w.Header().Set("Content-Length", cl)
	}
	w.Header().Set("Cache-Control", "public, max-age=5")

	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("copy image body failed: %v", err)
	}
}

// ========== 编码题 API ==========

// handleCodingQuestions GET /api/coding/questions
func handleCodingQuestions(w http.ResponseWriter, r *http.Request) {
	result := make([]CodingQuestionView, 0, len(codingQuestions))
	for i, q := range codingQuestions {
		result = append(result, CodingQuestionView{
			ID: q.ID, Category: q.Category, Difficulty: q.Difficulty,
			Title: q.Title, Desc: q.Desc, Template: q.Template,
			Hint: q.Hint, Total: len(codingQuestions), Index: i,
		})
	}
	writeJSON(w, http.StatusOK, result)
}

// handleCodingRun POST /api/coding/run
func handleCodingRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CodingRunReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Code) == "" {
		writeJSON(w, http.StatusOK, CodingRunResp{Pass: false, Error: "代码不能为空"})
		return
	}

	dangerous := []string{"os.Remove", "os.Exit", "exec.Command", "os/exec", "syscall", "unsafe.Pointer", "net.Dial", "net.Listen", "http.Get", "http.Post"}
	codeLower := strings.ToLower(req.Code)
	for _, d := range dangerous {
		if strings.Contains(codeLower, strings.ToLower(d)) {
			writeJSON(w, http.StatusOK, CodingRunResp{Pass: false, Error: fmt.Sprintf("禁止使用 %s", d)})
			return
		}
	}

	var q *CodingQuestion
	for i := range codingQuestions {
		if codingQuestions[i].ID == req.ID {
			q = &codingQuestions[i]
			break
		}
	}
	if q == nil {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}

	testCode := buildTestCode(q, req.Code)
	output, runErr := runGoCode(testCode, 10*time.Second)

	resp := CodingRunResp{
		Pass:     strings.Contains(output, "ALL_PASS"),
		Output:   output,
		Solution: q.Solution,
	}
	if runErr != nil && !resp.Pass {
		resp.Error = runErr.Error()
	}
	writeJSON(w, http.StatusOK, resp)
}

// ========== 代码解析工具 ==========

func buildTestCode(q *CodingQuestion, userCode string) string {
	testCode := q.TestCode
	if strings.Contains(testCode, "__USER_CODE__") {
		body := extractFuncBody(userCode, q.Template)
		testCode = strings.Replace(testCode, "__USER_CODE__", body, 1)
	}
	if strings.Contains(testCode, "__USER_FIELDS__") {
		testCode = strings.Replace(testCode, "__USER_FIELDS__", extractSection(userCode, "// 在此定义字段", "}"), 1)
		testCode = strings.Replace(testCode, "__USER_INC__", extractMethodBody(userCode, "func (c *Counter) Inc()"), 1)
		testCode = strings.Replace(testCode, "__USER_DEC__", extractMethodBody(userCode, "func (c *Counter) Dec()"), 1)
		testCode = strings.Replace(testCode, "__USER_VALUE__", extractMethodBody(userCode, "func (c *Counter) Value() int"), 1)
	}
	if strings.Contains(testCode, "__USER_FIELDS__") && strings.Contains(testCode, "__USER_FUNCS__") {
		testCode = strings.Replace(testCode, "__USER_FIELDS__", extractSection(userCode, "// 在此定义字段", "}"), 1)
		funcs := extractAllFuncs(userCode)
		testCode = strings.Replace(testCode, "__USER_FUNCS__", funcs, 1)
	}
	return testCode
}

func extractFuncBody(userCode, template string) string {
	lines := strings.Split(userCode, "\n")
	var result []string
	inMain := false
	mainBrace := 0
	skipPackage := false
	skipImport := false
	importBrace := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "package ") {
			skipPackage = true
			continue
		}
		if skipPackage && trimmed == "" {
			skipPackage = false
			continue
		}
		if strings.HasPrefix(trimmed, "import") {
			skipImport = true
			if strings.Contains(trimmed, "(") {
				importBrace++
			}
			if !strings.Contains(trimmed, "(") {
				skipImport = false
			}
			continue
		}
		if skipImport {
			if strings.Contains(trimmed, "(") {
				importBrace++
			}
			if strings.Contains(trimmed, ")") {
				importBrace--
				if importBrace <= 0 {
					skipImport = false
				}
			}
			continue
		}
		if strings.HasPrefix(trimmed, "func main()") {
			inMain = true
			mainBrace = 0
			if strings.Contains(trimmed, "{") {
				mainBrace++
			}
			continue
		}
		if inMain {
			mainBrace += strings.Count(trimmed, "{") - strings.Count(trimmed, "}")
			if mainBrace <= 0 {
				inMain = false
			}
			continue
		}
		result = append(result, line)
	}

	code := strings.TrimSpace(strings.Join(result, "\n"))
	if idx := strings.Index(code, "{"); idx >= 0 {
		if strings.Contains(code[:idx], "func") {
			braceCount := 0
			start := idx + 1
			for i := idx; i < len(code); i++ {
				if code[i] == '{' {
					braceCount++
				}
				if code[i] == '}' {
					braceCount--
					if braceCount == 0 {
						return strings.TrimSpace(code[start:i])
					}
				}
			}
		}
	}
	return code
}

func extractSection(code, startMarker, endMarker string) string {
	idx := strings.Index(code, startMarker)
	if idx < 0 {
		if sIdx := strings.Index(code, "struct {"); sIdx >= 0 {
			braceStart := strings.Index(code[sIdx:], "{")
			if braceStart >= 0 {
				pos := sIdx + braceStart + 1
				depth := 1
				for i := pos; i < len(code); i++ {
					if code[i] == '{' {
						depth++
					}
					if code[i] == '}' {
						depth--
						if depth == 0 {
							return strings.TrimSpace(code[pos:i])
						}
					}
				}
			}
		}
		return ""
	}
	start := idx + len(startMarker)
	rest := code[start:]
	if eIdx := strings.Index(rest, endMarker); eIdx >= 0 {
		return strings.TrimSpace(rest[:eIdx])
	}
	return strings.TrimSpace(rest)
}

func extractMethodBody(code, sig string) string {
	idx := strings.Index(code, sig)
	if idx < 0 {
		return ""
	}
	rest := code[idx+len(sig):]
	bIdx := strings.Index(rest, "{")
	if bIdx < 0 {
		return ""
	}
	depth := 0
	start := bIdx + 1
	for i := bIdx; i < len(rest); i++ {
		if rest[i] == '{' {
			depth++
		}
		if rest[i] == '}' {
			depth--
			if depth == 0 {
				return strings.TrimSpace(rest[start:i])
			}
		}
	}
	return ""
}

func extractAllFuncs(code string) string {
	lines := strings.Split(code, "\n")
	var funcs []string
	inFunc := false
	depth := 0
	var buf []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inFunc && strings.HasPrefix(trimmed, "func") && !strings.HasPrefix(trimmed, "func main") {
			inFunc = true
			depth = 0
			buf = nil
		}
		if inFunc {
			buf = append(buf, line)
			depth += strings.Count(line, "{") - strings.Count(line, "}")
			if depth <= 0 && len(buf) > 1 {
				funcs = append(funcs, strings.Join(buf, "\n"))
				inFunc = false
				buf = nil
			}
		}
	}
	return strings.Join(funcs, "\n\n")
}

// runGoCode 在临时目录执行 Go 代码
func runGoCode(code string, timeout time.Duration) (string, error) {
	tmpDir, err := os.MkdirTemp("", "quiz-arena-*")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module quiz\n\ngo 1.21\n"), 0644)
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644); err != nil {
		return "", fmt.Errorf("写入代码失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "main.go")
	cmd.Dir = tmpDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	output := stdout.String()
	if stderr.Len() > 0 {
		errStr := stderr.String()
		errStr = strings.ReplaceAll(errStr, tmpDir+"/", "")
		errStr = strings.ReplaceAll(errStr, tmpDir, "")
		if output == "" {
			output = errStr
		} else {
			output += "\n" + errStr
		}
	}
	if len(output) > 4000 {
		output = output[:4000] + "\n...(输出过长，已截断)"
	}
	return output, err
}

// ========== 通用 HTTP 工具 ==========

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":      "ok",
		"total":       len(questions),
		"codingTotal": len(codingQuestions),
	})
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// ========== 入口 ==========

func main() {
	addr := flag.String("addr", ":8090", "API 监听地址")
	dataDir := flag.String("data", "data", "题库 JSON 文件目录")
	flag.Parse()

	// 从 JSON 文件加载题库
	loadQuestionsFromJSON(*dataDir)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/categories", handleCategories)
	mux.HandleFunc("/api/questions", handleQuestions)
	mux.HandleFunc("/api/submit", handleSubmit)
	mux.HandleFunc("/api/image", handleImage)
	mux.HandleFunc("/api/coding/questions", handleCodingQuestions)
	mux.HandleFunc("/api/coding/run", handleCodingRun)

	handler := corsMiddleware(logMiddleware(mux))

	log.Printf("🎯 quiz-arena API 启动: http://localhost%s", *addr)
	log.Printf("📚 选择题: %d 道 · 编码题: %d 道", len(questions), len(codingQuestions))
	log.Fatal(http.ListenAndServe(*addr, handler))
}
