package main

import "fmt"

// 题目：实现事件发射器（观察者模式）
//
// EventEmitter 支持：
//   On(event string, handler func(data any)) — 注册监听
//   Off(event string) — 移除所有该事件的监听
//   Emit(event string, data any) — 触发事件，依次调用所有 handler
//   Once(event string, handler func(data any)) — 只触发一次后自动移除
//
// 示例：
//   em := NewEventEmitter()
//   em.On("login", func(d any) { fmt.Println("user login:", d) })
//   em.Emit("login", "Alice")  // 输出: user login: Alice

type EventEmitter struct {
	listeners map[string][]func(data any)
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string][]func(data any)),
	}
}

func (e *EventEmitter) On(event string, handler func(data any)) {
	e.listeners[event] = append(e.listeners[event], handler)
}

func (e *EventEmitter) Off(event string) {
	delete(e.listeners, event)
}

func (e *EventEmitter) Emit(event string, data any) {
	for _, handler := range e.listeners[event] {
		handler(data)
	}
}

func (e *EventEmitter) Once(event string, handler func(data any)) {
	e.On(event, func(data any) {
		handler(data)
		e.Off(event) // 执行一次后移除该事件所有监听
	})
}

func main() {
	em := NewEventEmitter()
	em.On("login", func(d any) { fmt.Println("handler1:", d) })
	em.On("login", func(d any) { fmt.Println("handler2:", d) })
	em.Once("logout", func(d any) { fmt.Println("bye:", d) })

	em.Emit("login", "Alice")
	// handler1: Alice
	// handler2: Alice

	em.Emit("logout", "Bob")
	// bye: Bob

	em.Emit("logout", "Charlie")
	// （无输出，Once 已移除）
}
