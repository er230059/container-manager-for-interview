package server

import (
	"container-manager/internal/server/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the routes for the application.
// It accepts handlers as dependencies.
func RegisterRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser)
	}
}
