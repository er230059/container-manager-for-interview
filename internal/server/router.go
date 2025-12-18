package server

import (
	"container-manager/docs"
	"container-manager/internal/server/handler"
	"container-manager/internal/server/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes sets up the routes for the application.
// It accepts handlers as dependencies.
func RegisterRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	containerHandler *handler.ContainerHandler,
	fileHandler *handler.FileHandler,
	jobHandler *handler.JobHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		containerRoutes.PATCH("/:id/stop", containerHandler.StopContainer)
		containerRoutes.DELETE("/:id", containerHandler.RemoveContainer)
	}

	fileRoutes := router.Group("/files")
	fileRoutes.Use(authMiddleware.Handle())
	{
		fileRoutes.POST("", fileHandler.UploadFile)
	}

	jobRoutes := router.Group("/jobs")
	jobRoutes.Use(authMiddleware.Handle())
	{
		jobRoutes.GET("/:id", jobHandler.GetJob)
	}
}
