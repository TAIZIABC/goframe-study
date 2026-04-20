// handler/handler.go
// WebSocket handler + HTTP 推送接口
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"

	"wspush/hub"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clientID uint64

// WSHandler WebSocket 连接处理
func WSHandler(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		id := fmt.Sprintf("client-%d", atomic.AddUint64(&clientID, 1))
		client := hub.NewClient(id, h)
		h.Register(client)

		// 发送欢迎消息
		conn.WriteJSON(map[string]string{
			"type": "connected",
			"id":   id,
		})

		// 读协程：处理客户端发来的订阅/取消订阅命令
		go func() {
			defer func() {
				h.Unregister(client)
				conn.Close()
			}()
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					break
				}
				var cmd struct {
					Action  string `json:"action"`  // subscribe / unsubscribe
					Channel string `json:"channel"`
				}
				if err := json.Unmarshal(msg, &cmd); err != nil {
					continue
				}
				switch cmd.Action {
				case "subscribe":
					if cmd.Channel != "" {
						h.Subscribe(client, cmd.Channel)
					}
				case "unsubscribe":
					if cmd.Channel != "" {
						h.Unsubscribe(client, cmd.Channel)
					}
				}
			}
		}()

		// 写协程：将消息发送给客户端
		go func() {
			defer conn.Close()
			for msg := range client.Send {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					break
				}
			}
		}()
	}
}

// PublishHandler HTTP 推送接口 POST /api/publish
func PublishHandler(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Channel string `json:"channel"`
			Data    string `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "参数解析失败"})
			return
		}
		if req.Channel == "" || req.Data == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "channel 和 data 必填"})
			return
		}

		h.Publish(req.Channel, req.Data)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"channel": req.Channel,
		})
	}
}

// StatsHandler 统计信息 GET /api/stats
func StatsHandler(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.Stats())
	}
}
