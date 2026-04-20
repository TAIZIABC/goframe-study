// api/api.go
// HTTP 接口：查看任务状态和执行历史
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"gocron/scheduler"
)

func StartAPI(port int, sched *scheduler.Scheduler) error {
	mux := http.NewServeMux()

	// 任务列表
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tasks": sched.GetTasks(),
		})
	})

	// 任务历史: /api/history?name=xxx&limit=20
	mux.HandleFunc("/api/history", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "缺少 name 参数"})
			return
		}
		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			limit, _ = strconv.Atoi(l)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":    name,
			"history": sched.GetHistory(name, limit),
		})
	})

	// 面板首页
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, dashboardHTML)
	})

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, mux)
}

var dashboardHTML = `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><title>GoCron Dashboard</title>
<style>
  body{font-family:system-ui;margin:0;padding:20px;background:#f5f5f5}
  h1{color:#333;border-bottom:2px solid #4a90d9;padding-bottom:8px}
  table{width:100%;border-collapse:collapse;background:#fff;border-radius:8px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,.1)}
  th{background:#4a90d9;color:#fff;padding:10px;text-align:left}
  td{padding:8px 10px;border-bottom:1px solid #eee}
  .running{color:#e67e22;font-weight:bold} .idle{color:#27ae60}
  .badge{display:inline-block;padding:2px 8px;border-radius:10px;font-size:12px}
  .ok{background:#d4edda;color:#155724} .fail{background:#f8d7da;color:#721c24}
</style></head>
<body><h1>GoCron Dashboard</h1>
<table id="tasks"><tr><th>任务</th><th>Cron</th><th>命令</th><th>状态</th><th>下次执行</th><th>上次执行</th><th>执行/失败</th></tr></table>
<script>
async function refresh(){
  const r=await fetch('/api/tasks');const d=await r.json();
  let html='<tr><th>任务</th><th>Cron</th><th>命令</th><th>状态</th><th>下次执行</th><th>上次执行</th><th>执行/失败</th></tr>';
  for(const t of d.tasks||[]){
    const sc=t.status==='running'?'running':'idle';
    html+=` + "`" + `<tr><td><b>${t.name}</b></td><td><code>${t.cron}</code></td><td><code>${t.command}</code></td>
    <td class="${sc}">${t.status}</td><td>${t.next_run}</td><td>${t.last_run||'-'}</td>
    <td><span class="badge ok">${t.run_count}</span> / <span class="badge fail">${t.fail_count}</span></td></tr>` + "`" + `;
  }
  document.getElementById('tasks').innerHTML=html;
}
refresh();setInterval(refresh,2000);
</script></body></html>`

