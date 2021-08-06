package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	*gin.Engine
}

func New(e *gin.Engine) *App {
	return &App{Engine: e}
}

func (app *App) PUT(relativePath string, handler HandlerFunc) {
	app.Engine.PUT(relativePath, func(c *gin.Context) {
		handler(&Context{Context: c})
	})
}

type HandlerFunc func(*Context)

type Context struct {
	*gin.Context
}

func (c *Context) InternalError(err error) {
	c.Context.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}
