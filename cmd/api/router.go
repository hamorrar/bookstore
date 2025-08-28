package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()
	return g
}

func InitRouter(router *gin.Engine) {
	// Set up router and routes
	router.GET("/api/v1/ping", ping)
	// router.POST("/api/v1/login", auth.Login)
	router.GET("/", homepage)
	// router.POST("/api/v1/Books/toggle", books.Toggle)
	// router.POST("/api/v1/Books/add", books.Add(db))
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "version": "v1"})
}

func homepage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"LoggedIn": "",
		"Username": "",
	})
}
