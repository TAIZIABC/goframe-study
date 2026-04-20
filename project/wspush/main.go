package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"wspush/handler"
	"wspush/hub"
)

func main() {
	port := flag.Int("port", 8090, "服务端口")
	flag.Parse()

	h := hub.New()
	go h.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handler.WSHandler(h))
	mux.HandleFunc("/api/publish", handler.PublishHandler(h))
	mux.HandleFunc("/api/stats", handler.StatsHandler(h))
	mux.HandleFunc("/", pageHandler)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Println()
	fmt.Println("  ┌──────────────────────────────────────────┐")
	fmt.Println("  │     WebSocket 消息推送服务 (wspush)      │")
	fmt.Println("  ├──────────────────────────────────────────┤")
	fmt.Printf("  │  🌐 测试页面:  http://localhost:%d\n", *port)
	fmt.Printf("  │  🔌 WebSocket: ws://localhost:%d/ws\n", *port)
	fmt.Printf("  │  📡 推送接口:  POST http://localhost:%d/api/publish\n", *port)
	fmt.Printf("  │  📊 统计接口:  GET  http://localhost:%d/api/stats\n", *port)
	fmt.Println("  │  按 Ctrl+C 停止")
	fmt.Println("  └──────────────────────────────────────────┘")
	fmt.Println()

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "启动失败: %v\n", err)
		os.Exit(1)
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, testPage)
}

var testPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8"><title>WebSocket 推送测试</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:system-ui;background:#1a1a2e;color:#e0e0e0;padding:20px}
h1{color:#00d4ff;margin-bottom:20px;font-size:24px}
.container{display:grid;grid-template-columns:1fr 1fr;gap:20px;max-width:1000px;margin:0 auto}
.panel{background:#16213e;border-radius:12px;padding:20px;border:1px solid #0f3460}
.panel h2{color:#00d4ff;font-size:16px;margin-bottom:12px;border-bottom:1px solid #0f3460;padding-bottom:8px}
input,button{padding:8px 12px;border-radius:6px;border:1px solid #0f3460;background:#0f3460;color:#e0e0e0;font-size:14px}
input{width:100%;margin-bottom:8px}
input:focus{outline:none;border-color:#00d4ff}
button{cursor:pointer;background:#00d4ff;color:#1a1a2e;font-weight:bold;border:none;width:100%;margin-top:4px}
button:hover{background:#00b4d8}
button.red{background:#e74c3c}button.red:hover{background:#c0392b}
.messages{height:300px;overflow-y:auto;background:#0a0a1a;border-radius:8px;padding:10px;font-family:monospace;font-size:13px;margin-top:10px}
.msg{padding:4px 0;border-bottom:1px solid #1a1a2e}
.msg .ch{color:#00d4ff;font-weight:bold}
.msg .time{color:#666;font-size:11px}
.msg .data{color:#4ade80}
.status{display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:6px}
.online{background:#4ade80}.offline{background:#e74c3c}
.tags{display:flex;gap:6px;flex-wrap:wrap;margin-top:8px}
.tag{background:#0f3460;padding:4px 10px;border-radius:12px;font-size:12px;display:flex;align-items:center;gap:4px}
.tag .x{cursor:pointer;color:#e74c3c}
#publish-panel textarea{width:100%;height:60px;background:#0f3460;color:#e0e0e0;border:1px solid #0f3460;border-radius:6px;padding:8px;resize:none;font-size:14px}
</style>
</head>
<body>
<h1>🔌 WebSocket 消息推送测试</h1>
<div class="container">
  <!-- 左侧：连接和订阅 -->
  <div>
    <div class="panel">
      <h2><span id="dot" class="status offline"></span>连接状态: <span id="state">未连接</span></h2>
      <button onclick="connect()">连接</button>
      <button class="red" onclick="disconnect()">断开</button>
      <div style="margin-top:6px;font-size:12px;color:#666">ID: <span id="cid">-</span></div>
    </div>
    <div class="panel" style="margin-top:20px">
      <h2>📢 订阅频道</h2>
      <input id="ch-input" placeholder="输入频道名 (如: news, chat, alerts)" onkeydown="if(event.key==='Enter')sub()">
      <button onclick="sub()">订阅</button>
      <div class="tags" id="channels"></div>
    </div>
    <div class="panel" style="margin-top:20px" id="publish-panel">
      <h2>📨 发送消息 (HTTP)</h2>
      <input id="pub-ch" placeholder="频道名">
      <textarea id="pub-data" placeholder="消息内容"></textarea>
      <button onclick="publish()">推送</button>
    </div>
  </div>
  <!-- 右侧：消息 -->
  <div class="panel">
    <h2>📬 收到的消息</h2>
    <button class="red" onclick="document.getElementById('log').innerHTML=''">清空</button>
    <div class="messages" id="log"></div>
  </div>
</div>
<script>
let ws=null,subs=new Set();
function connect(){
  if(ws)return;
  ws=new WebSocket('ws://'+location.host+'/ws');
  ws.onopen=()=>{setState('已连接',true)};
  ws.onclose=()=>{ws=null;setState('已断开',false)};
  ws.onmessage=(e)=>{
    const d=JSON.parse(e.data);
    if(d.type==='connected'){document.getElementById('cid').textContent=d.id;return}
    addMsg(d);
  };
}
function disconnect(){if(ws){ws.close();ws=null}}
function setState(t,on){
  document.getElementById('state').textContent=t;
  document.getElementById('dot').className='status '+(on?'online':'offline');
}
function sub(){
  const ch=document.getElementById('ch-input').value.trim();
  if(!ch||!ws)return;
  ws.send(JSON.stringify({action:'subscribe',channel:ch}));
  subs.add(ch);renderChannels();
  document.getElementById('ch-input').value='';
}
function unsub(ch){
  if(ws)ws.send(JSON.stringify({action:'unsubscribe',channel:ch}));
  subs.delete(ch);renderChannels();
}
function renderChannels(){
  const el=document.getElementById('channels');
  el.innerHTML=[...subs].map(c=>'<span class="tag">#'+c+' <span class="x" onclick="unsub(\''+c+'\')">✕</span></span>').join('');
}
function addMsg(d){
  const el=document.getElementById('log');
  el.innerHTML='<div class="msg"><span class="ch">#'+d.channel+'</span> <span class="data">'+d.data+'</span> <span class="time">'+d.timestamp+'</span></div>'+el.innerHTML;
}
function publish(){
  const ch=document.getElementById('pub-ch').value.trim();
  const data=document.getElementById('pub-data').value.trim();
  if(!ch||!data)return;
  fetch('/api/publish',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({channel:ch,data:data})})
    .then(r=>r.json()).then(()=>{document.getElementById('pub-data').value=''});
}
connect();
</script>
</body></html>`
