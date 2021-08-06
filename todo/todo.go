package todo

import (
	"fmt"
	"net/http"
	"strconv"
	"todo/entities"

	"github.com/gin-gonic/gin"
)

var index int
var tasks map[int]*entities.Task = make(map[int]*entities.Task)

type NewTaskTodo struct {
	Topic string `json:"task" xml:"Task" msgpack:"task" yaml:"task"`
}

type Inserter interface {
	Insert(interface{}) error
}

type Repository interface {
	NewTask(*entities.Task) error
}

type Todo struct {
	// db   Inserter
	repo Repository
}

func (todo Todo) Add(c *gin.Context) {
	var task NewTaskTodo
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := todo.repo.NewTask(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
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

func List() map[int]*entities.Task {
	return tasks
}

func New(task string) {
	defer func() {
		index++
	}()

	tasks[index] = &entities.Task{
		Title: task,
		Done:  false,
	}
}
