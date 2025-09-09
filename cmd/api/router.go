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

		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)

		v1.GET("/books/:id", app.getBook)
		v1.GET("/books", app.getPageOfBooks)

		v1.GET("/orders", app.getAllOrders)
		v1.GET("/orders/:id", app.getOrder)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())

	{
		authGroup.GET("/users/all", app.getAllUsers)
		authGroup.GET("/users/:id", app.getUser)
		authGroup.PUT("/users/:id", app.updateUser)
		authGroup.DELETE("/users/:id", app.deleteUser)

		authGroup.GET("/books/all", app.getAllBooks)
		authGroup.POST("/books", app.createBook)
		authGroup.PUT("/books/:id", app.updateBook)
		authGroup.DELETE("/books/:id", app.deleteBook)

		authGroup.POST("/orders", app.createOrder)
		authGroup.PUT("/orders/:id", app.updateOrder)
		authGroup.DELETE("/orders/:id", app.deleteOrder)
	}

	return g
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "version": "v1"})
}
