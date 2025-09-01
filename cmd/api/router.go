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
		v1.PUT("/books/:id", app.updateBook)

		v1.POST("/users", app.createUser)
		v1.GET("/users", app.getAllUsers)
		v1.GET("/users/:id", app.getUser)
		v1.DELETE("/users/:id", app.deleteUser)
		v1.PUT("/users/:id", app.updateUser)

		v1.POST("/orders", app.createOrder)
		v1.GET("/orders", app.getAllOrders)
		v1.GET("/orders/:id", app.getOrder)
		v1.DELETE("/orders/:id", app.deleteOrder)
		v1.PUT("/orders/:id", app.updateOrder)

	}
	return g
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "version": "v1"})
}
