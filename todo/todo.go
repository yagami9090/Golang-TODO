package todo

import (
	"net/http"
	"strconv"
	"todo/entities"

	"github.com/gin-gonic/gin"
)

type Servicer interface {
	Add(entities.Task) error
	Done(uint) error
	List() ([]*entities.Task, error)
}

type Todo struct {
	srv Service
}

func New(srv Service) Todo {
	return Todo{srv: srv}
}

func (todo Todo) Add(c *gin.Context) {
	var task NewTaskTodo
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := todo.srv.Add(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// Business Domain, UseCase, Core
func Add(task NewTaskTodo, repo Repository) error {
	return repo.NewTask(&task).Error
}

func (todo Todo) MarkDone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// if err := todo.repo.TaskDone(uint(id)); err != nil {
	if err := todo.srv.Done(uint(id)); err != nil {
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
	list, err := todo.srv.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, list)
}
