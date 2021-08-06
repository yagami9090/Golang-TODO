package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Ctx struct {
	*gin.Context
}

func (c *Ctx) InternalServerError(err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}

func (c *Ctx) OK() {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type App struct {
	*gin.Engine
}

func New() *App {
	return &App{Engine: gin.Default()}
}

type HandlerFunc func(*Ctx)

func (app *App) PUT(relativePath string, handler HandlerFunc) {
	app.Engine.PUT(relativePath, func(c *gin.Context) {
		handler(&Ctx{Context: c})
	})
}
func (app *App) GET(relativePath string, handler HandlerFunc) {
	app.Engine.GET(relativePath, func(c *gin.Context) {
		handler(&Ctx{Context: c})
	})
}
