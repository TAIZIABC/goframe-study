package v1

// TodoItem 是列表返回的单个条目
type TodoItem struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
