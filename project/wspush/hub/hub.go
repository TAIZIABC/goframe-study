// hub/hub.go
// 频道管理中心：管理客户端连接、频道订阅、消息广播
package hub

import (
	"fmt"
	"sync"
	"time"
)

// Client 一个 WebSocket 客户端
type Client struct {
	ID       string
	Channels map[string]bool
	Send     chan []byte
	hub      *Hub
}

// Message 推送消息
type Message struct {
	Channel   string `json:"channel"`
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

// Hub 频道管理中心
type Hub struct {
	clients    map[*Client]bool
	channels   map[string]map[*Client]bool // channel -> clients
	register   chan *Client
	unregister chan *Client
	subscribe  chan subscription
	broadcast  chan Message
	mu         sync.RWMutex
}

type subscription struct {
	client  *Client
	channel string
	action  string // "sub" or "unsub"
}

func New() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		channels:   make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		subscribe:  make(chan subscription),
		broadcast:  make(chan Message, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Printf("  \033[32m+\033[0m 客户端连接 [%s] (在线: %d)\n",
				client.ID, len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				// 从所有频道移除
				for ch := range client.Channels {
					if subs, ok := h.channels[ch]; ok {
						delete(subs, client)
						if len(subs) == 0 {
							delete(h.channels, ch)
						}
					}
				}
				delete(h.clients, client)
				close(client.Send)
				fmt.Printf("  \033[31m-\033[0m 客户端断开 [%s] (在线: %d)\n",
					client.ID, len(h.clients))
			}

		case sub := <-h.subscribe:
			if sub.action == "sub" {
				if _, ok := h.channels[sub.channel]; !ok {
					h.channels[sub.channel] = make(map[*Client]bool)
				}
				h.channels[sub.channel][sub.client] = true
				sub.client.Channels[sub.channel] = true
				fmt.Printf("  📢 [%s] 订阅频道 #%s\n", sub.client.ID, sub.channel)
			} else {
				if subs, ok := h.channels[sub.channel]; ok {
					delete(subs, sub.client)
					if len(subs) == 0 {
						delete(h.channels, sub.channel)
					}
				}
				delete(sub.client.Channels, sub.channel)
				fmt.Printf("  🔕 [%s] 取消订阅 #%s\n", sub.client.ID, sub.channel)
			}

		case msg := <-h.broadcast:
			if subs, ok := h.channels[msg.Channel]; ok {
				data := fmt.Sprintf(`{"channel":"%s","data":"%s","timestamp":"%s"}`,
					msg.Channel, msg.Data, msg.Timestamp)
				for client := range subs {
					select {
					case client.Send <- []byte(data):
					default:
						// 发送缓冲区满，丢弃
					}
				}
				fmt.Printf("  📨 #%s → %d 个客户端: %s\n",
					msg.Channel, len(subs), truncate(msg.Data, 50))
			}
		}
	}
}

// Register 注册客户端
func (h *Hub) Register(c *Client) { h.register <- c }

// Unregister 注销客户端
func (h *Hub) Unregister(c *Client) { h.unregister <- c }

// Subscribe 订阅频道
func (h *Hub) Subscribe(c *Client, channel string) {
	h.subscribe <- subscription{client: c, channel: channel, action: "sub"}
}

// Unsubscribe 取消订阅
func (h *Hub) Unsubscribe(c *Client, channel string) {
	h.subscribe <- subscription{client: c, channel: channel, action: "unsub"}
}

// Publish 向频道推送消息
func (h *Hub) Publish(channel, data string) {
	h.broadcast <- Message{
		Channel:   channel,
		Data:      data,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// Stats 返回统计信息
func (h *Hub) Stats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	channelStats := make(map[string]int)
	for ch, subs := range h.channels {
		channelStats[ch] = len(subs)
	}

	return map[string]interface{}{
		"online_clients": len(h.clients),
		"channels":       channelStats,
	}
}

// NewClient 创建新客户端
func NewClient(id string, h *Hub) *Client {
	return &Client{
		ID:       id,
		Channels: make(map[string]bool),
		Send:     make(chan []byte, 64),
		hub:      h,
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
