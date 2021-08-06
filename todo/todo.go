package todo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var index int
var tasks map[int]*Task = make(map[int]*Task)

type Task struct {
	Title string
	Done  bool
}

type NewTaskTodo struct {
	Task string `json:"task"`
}

type Serializer interface {
	Decode(io.ReadCloser, interface{}) error
}

type App struct {
	serialize Serializer
}

func NewApp(serialize Serializer) App {
	return App{serialize: serialize}
}

func (app App) AddTask(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var task NewTaskTodo
	if err := app.serialize.Decode(r.Body, &task); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(task)

	New(task.Task)
}

func AddTask(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var task NewTaskTodo
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	New(task.Task)
}

func MarkDone(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index := vars["index"]
	i, err := strconv.Atoi(index)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	tasks[i].Done = true

}

type Encoder interface {
	Encode(io.Writer, interface{}) error
}

func ListTask(rw http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(rw).Encode(List()); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func List() map[int]*Task {
	return tasks
}

func New(task string) {
	defer func() {
		index++
	}()

	tasks[index] = &Task{
		Title: task,
		Done:  false,
	}
}
