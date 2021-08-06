package todo

import (
	"net/http"
	"strconv"
	"todo/app"
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

func (todo Todo) Add(c *app.Context) {
	var task NewTaskTodo
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := todo.srv.Add(entities.Task{Title: task.Topic}); err != nil {
		c.InternalError(err)
		return
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

	if err := todo.srv.Done(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
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
