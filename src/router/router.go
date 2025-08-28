package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	// Set up router and routes
	router.GET("/api/v1/ping", ping)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "version": "v1"})
}
