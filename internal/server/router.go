package server

import (
	"container-manager/internal/server/handler"
	"container-manager/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the routes for the application.
// It accepts handlers as dependencies.
func RegisterRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	containerHandler *handler.ContainerHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.POST("/login", userHandler.Login)
	}

	containerRoutes := router.Group("/containers")
	containerRoutes.Use(authMiddleware.Handle())
	{
		containerRoutes.POST("", containerHandler.CreateContainer)
		containerRoutes.PATCH("/:id/start", containerHandler.StartContainer)
	}
}
