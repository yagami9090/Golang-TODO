package todo

import (
	"net/http"
	"strconv"
	"todo/app"
	"todo/entities"
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

func (todo Todo) Add(c *app.Ctx) {
	var task NewTaskTodo
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := todo.srv.Add(entities.Task{Title: task.Topic}); err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK()
}

func (todo Todo) MarkDone(c *app.Ctx) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := todo.srv.Done(uint(id)); err != nil {
		c.InternalServerError(err)

		return
	}
	c.OK()
}

func (todo Todo) ListTask(c *app.Ctx) {
	list, err := todo.srv.List()
	if err != nil {
		c.InternalServerError(err)
		return
	}
	c.JSON(http.StatusOK, list)
}
