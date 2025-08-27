package router

import (
	"github.com/gin-gonic/gin"
)

func initRouter(router *gin.Engine) {
	// Set up router and routes
	router.GET("/api/v1/health", main.health)
}
