package server

import (
	"container-manager/internal/server/handler"
	"container-manager/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	containerHandler *handler.ContainerHandler,
	jwtSecret string,
) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.POST("/login", userHandler.Login)
	}

	containerRoutes := router.Group("/containers")
	containerRoutes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		containerRoutes.POST("", containerHandler.CreateContainer)
	}
}
