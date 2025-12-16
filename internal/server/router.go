package server

import (
	"leadtek/internal/server/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the routes for the application.
// It accepts handlers as dependencies.
func RegisterRoutes(router *gin.Engine, homeHandler *handler.HomeHandler) {
	router.GET("/", homeHandler.GetHome)
}
