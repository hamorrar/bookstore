package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")

	{
		v1.GET("/ping", ping)

		v1.POST("/books", app.createBook)
		v1.GET("/books", app.getAllBooks)
		v1.GET("/books/:id", app.getBook)
		v1.DELETE("/books/:id", app.deleteBook)
	}
	return g
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "version": "v1"})
}
