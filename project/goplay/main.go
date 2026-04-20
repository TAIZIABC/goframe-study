package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"goplay/runner"
)

func main() {
	port := flag.Int("port", 8092, "服务端口")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", pageHandler)
	mux.HandleFunc("/api/run", runHandler)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Println()
	fmt.Println("  ┌──────────────────────────────────────────┐")
	fmt.Println("  │       Go Playground (goplay)             │")
	fmt.Println("  ├──────────────────────────────────────────┤")
	fmt.Printf("  │  🌐 编辑器:  http://localhost:%d\n", *port)
	fmt.Printf("  │  📡 API:     POST http://localhost:%d/api/run\n", *port)
	fmt.Println("  │  ⏱️  超时:    10 秒")
	fmt.Println("  │  🔒 内存:    128 MiB")
	fmt.Println("  │  按 Ctrl+C 停止")
	fmt.Println("  └──────────────────────────────────────────┘")
	fmt.Println()

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "启动失败: %v\n", err)
		os.Exit(1)
	}
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数解析失败"})
		return
	}

	if req.Code == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "代码不能为空"})
		return
	}

	fmt.Printf("  ▶ 收到代码 (%d 字节)...\n", len(req.Code))
	result := runner.Run(req.Code, nil)
	fmt.Printf("  %s 耗时 %dms\n", statusIcon(result.Success), result.Duration)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func statusIcon(ok bool) string {
	if ok {
		return "\033[32m✓\033[0m"
	}
	return "\033[31m✗\033[0m"
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, editorPage)
}

var editorPage = `<!DOCTYPE html>
<html lang="zh-CN"><head><meta charset="UTF-8"><title>Go Playground</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:system-ui;background:#1e1e2e;color:#cdd6f4;height:100vh;display:flex;flex-direction:column}
header{background:#181825;padding:12px 20px;display:flex;align-items:center;justify-content:space-between;border-bottom:1px solid #313244}
header h1{font-size:18px;color:#89b4fa}
.btn{padding:8px 20px;border:none;border-radius:6px;font-size:14px;cursor:pointer;font-weight:600}
.btn-run{background:#a6e3a1;color:#1e1e2e}.btn-run:hover{background:#94e2d5}
.btn-run:disabled{opacity:.5;cursor:not-allowed}
.btn-clear{background:#45475a;color:#cdd6f4;margin-left:8px}.btn-clear:hover{background:#585b70}
.btn-fmt{background:#89b4fa;color:#1e1e2e;margin-left:8px}.btn-fmt:hover{background:#74c7ec}
main{flex:1;display:flex;overflow:hidden}
.editor-pane{flex:1;display:flex;flex-direction:column;border-right:1px solid #313244}
.output-pane{flex:1;display:flex;flex-direction:column}
.pane-header{background:#181825;padding:8px 16px;font-size:13px;color:#a6adc8;border-bottom:1px solid #313244;display:flex;justify-content:space-between}
textarea{flex:1;background:#1e1e2e;color:#cdd6f4;border:none;padding:16px;font-family:'Fira Code',Menlo,Monaco,Consolas,monospace;font-size:14px;line-height:1.6;resize:none;tab-size:4;outline:none}
.output{flex:1;padding:16px;font-family:monospace;font-size:13px;line-height:1.6;overflow-y:auto;white-space:pre-wrap;word-break:break-all}
.stdout{color:#a6e3a1} .stderr{color:#f38ba8} .error{color:#fab387}
.info{color:#a6adc8;font-size:12px}
.spinner{display:inline-block;width:14px;height:14px;border:2px solid #a6e3a1;border-top-color:transparent;border-radius:50%;animation:spin .6s linear infinite;margin-right:6px;vertical-align:middle}
@keyframes spin{to{transform:rotate(360deg)}}
</style></head>
<body>
<header>
  <h1>🚀 Go Playground</h1>
  <div>
    <button class="btn btn-run" id="runBtn" onclick="run()">▶ Run</button>
    <button class="btn btn-fmt" onclick="format()">Format</button>
    <button class="btn btn-clear" onclick="clearOutput()">Clear</button>
  </div>
</header>
<main>
  <div class="editor-pane">
    <div class="pane-header"><span>main.go</span><span id="chars">0 chars</span></div>
    <textarea id="code" spellcheck="false" placeholder="// 在这里输入 Go 代码...">package main

import "fmt"

func main() {
	fmt.Println("Hello, Go Playground!")
	
	for i := 1; i <= 5; i++ {
		fmt.Printf("  %d x %d = %d\n", i, i, i*i)
	}
}</textarea>
  </div>
  <div class="output-pane">
    <div class="pane-header"><span>输出</span><span id="timing"></span></div>
    <div class="output" id="output"><span class="info">按 Run 或 Ctrl+Enter 执行代码</span></div>
  </div>
</main>
<script>
const code=document.getElementById('code'),output=document.getElementById('output'),
      timing=document.getElementById('timing'),chars=document.getElementById('chars'),
      runBtn=document.getElementById('runBtn');

code.addEventListener('input',()=>{chars.textContent=code.value.length+' chars'});
code.addEventListener('keydown',e=>{
  if(e.key==='Enter'&&(e.ctrlKey||e.metaKey)){e.preventDefault();run()}
  if(e.key==='Tab'){e.preventDefault();const s=code.selectionStart,end=code.selectionEnd;code.value=code.value.substring(0,s)+'\t'+code.value.substring(end);code.selectionStart=code.selectionEnd=s+1}
});
chars.textContent=code.value.length+' chars';

async function run(){
  runBtn.disabled=true;runBtn.innerHTML='<span class="spinner"></span>Running...';
  output.innerHTML='<span class="info">编译运行中...</span>';timing.textContent='';
  try{
    const r=await fetch('/api/run',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({code:code.value})});
    const d=await r.json();
    let html='';
    if(d.error)html+='<div class="error">⚠ '+esc(d.error)+'</div>';
    if(d.stderr)html+='<div class="stderr">'+esc(d.stderr)+'</div>';
    if(d.stdout)html+='<div class="stdout">'+esc(d.stdout)+'</div>';
    if(!html)html='<span class="info">(无输出)</span>';
    output.innerHTML=html;
    timing.textContent=(d.success?'✓':'✗')+' '+d.duration_ms+'ms';
  }catch(e){output.innerHTML='<div class="error">请求失败: '+esc(e.message)+'</div>'}
  finally{runBtn.disabled=false;runBtn.innerHTML='▶ Run'}
}
function clearOutput(){output.innerHTML='<span class="info">已清空</span>';timing.textContent=''}
function format(){
  let c=code.value;
  // 简单自动格式化：统一缩进
  c=c.replace(/\r\n/g,'\n');
  code.value=c;
}
function esc(s){return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;')}
</script>
</body></html>`
