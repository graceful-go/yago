package todo

type TodoStatus = int64

var (
	TodoStatusNone TodoStatus = 0
	TodoStatusDone TodoStatus = 1
)

type TodoItem struct {
	Title      string     `json:"title"`
	Status     TodoStatus `json:"status"`
	CreateTime int64      `json:"createTime"`
	DoneTime   int64      `json:"doneTime"`
	Stared     bool       `json:"stared"`
}
