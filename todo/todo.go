package todo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func AddTask(c *gin.Context) {
	fmt.Println("start add task")
	var task NewTaskTodo
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	New(task.Task)
}

func MarkDone(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	tasks[i].Done = true
}

func ListTask(c *gin.Context) {
	c.JSON(http.StatusOK, List())
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
