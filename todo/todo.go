package todo

import (
	"encoding/json"
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
	Decode(io.Reader, interface{}) error
	Encode(io.Writer, interface{}) error
}

type JSONSerializer struct{}

func (j JSONSerializer) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
func (j JSONSerializer) Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func NewJSONSerializer() JSONSerializer {
	return JSONSerializer{}
}

type App struct {
	serialize Serializer
}

func NewApp(serialize Serializer) *App {
	return &App{
		serialize: serialize,
	}
}

func (app *App) AddTask(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var task NewTaskTodo
	if err := app.serialize.Decode(r.Body, &task); err != nil {
		// if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

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
