package todo

type NewTaskTodo struct {
	Topic string `json:"task" xml:"Task" msgpack:"task" yaml:"task"`
}
