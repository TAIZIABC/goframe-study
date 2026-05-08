package main

import (
	"container/list"
	"fmt"
)

// 题目：实现 LRU 缓存
//
// NewLRU(capacity int) *LRU
// Get(key string) (string, bool) — 存在则返回值并标记为最近使用
// Put(key, value string) — 插入/更新，超出容量淘汰最久未使用的
//
// 实现思路：双向链表 + map
// - 链表头部 = 最近使用，尾部 = 最久未使用
// - map 存 key → 链表节点，实现 O(1) 查找

type entry struct {
	key   string
	value string
}

type LRU struct {
	cap   int
	ll    *list.List               // 双向链表，Front = 最近使用
	cache map[string]*list.Element // key → 链表节点
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		cap:   capacity,
		ll:    list.New(),
		cache: make(map[string]*list.Element),
	}
}

func (l *LRU) Get(key string) (string, bool) {
	if elem, ok := l.cache[key]; ok {
		l.ll.MoveToFront(elem) // 标记为最近使用
		return elem.Value.(*entry).value, true
	}
	return "", false
}

func (l *LRU) Put(key, value string) {
	if elem, ok := l.cache[key]; ok {
		// 已存在：更新值并移到最前
		elem.Value.(*entry).value = value
		l.ll.MoveToFront(elem)
		return
	}
	// 容量已满：淘汰最久未使用的（链表尾部）
	if l.ll.Len() >= l.cap {
		oldest := l.ll.Back()
		l.ll.Remove(oldest)
		delete(l.cache, oldest.Value.(*entry).key)
	}
	// 插入新元素到最前
	elem := l.ll.PushFront(&entry{key: key, value: value})
	l.cache[key] = elem
}

func main() {
	c := NewLRU(2)
	c.Put("a", "1")
	c.Put("b", "2")
	fmt.Println(c.Get("a")) // 1 true
	c.Put("c", "3")         // 淘汰 b
	fmt.Println(c.Get("b")) // "" false
	fmt.Println(c.Get("c")) // 3 true
}
