package todo

import (
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
	TaskDone(uint) error
	List() ([]*entities.Task, error)
}

type Todo struct {
	repo Repository
}

func New(repo Repository) Todo {
	return Todo{repo: repo}
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

func (todo Todo) MarkDone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := todo.repo.TaskDone(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	// tasks[i].Done = true
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (todo Todo) ListTask(c *gin.Context) {
	list, err := todo.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, list)
}
