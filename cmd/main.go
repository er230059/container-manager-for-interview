package main

import (
	"fmt"
	"log"

	"container-manager/internal/application"
	containerruntime "container-manager/internal/infrastructure/container_runtime"
	database "container-manager/internal/infrastructure/database/sql"
	"container-manager/internal/infrastructure/repository"
	"container-manager/internal/server"
	"container-manager/internal/server/handler"
	"container-manager/internal/server/middleware"
	"container-manager/pkg/config"
	"container-manager/pkg/postgres"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title Container Manager API
// @version 1.0
// @description This is a sample API for managing containers.

// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// 0. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 1. Create dependencies (Composition Root)

	// Infrastructure Layer - Database
	db, err := postgres.NewClient(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	userDB := database.NewUserDatabase(db)
	containerUserDB := database.NewContainerUserDatabase(db)
	jobDB := database.NewJobDatabase(db)

	// Infrastructure Layer - Container Runtime
	runtime, err := containerruntime.NewDockerContainerRuntime()
	if err != nil {
		log.Fatalf("failed to create container runtime: %v", err)
	}

	// Infrastructure Layer - File Storage
	fileStorage := repository.NewLocalFileStorage(cfg.Storage.BasePath)

	// ID Generation
	idNode, err := snowflake.NewNode(cfg.Snowflake.MachineID)
	if err != nil {
		log.Fatalf("failed to create snowflake node: %v", err)
	}

	// Domain Layer - Repositories
	userRepo := repository.NewUserRepository(userDB)
	containerRepo := repository.NewContainerRepository(runtime, containerUserDB)
	jobRepo := repository.NewJobRepository(jobDB)

	// Application Layer
	userService := application.NewUserService(userRepo, idNode, cfg.Server.JWTSecret)
	fileService := application.NewFileService(fileStorage)
	containerService := application.NewContainerService(containerRepo, jobRepo)
	jobService := application.NewJobService(jobRepo)

	// Handler Layer
	authMiddleware := middleware.NewAuthMiddleware(cfg.Server.JWTSecret)
	userHandler := handler.NewUserHandler(userService)
	containerHandler := handler.NewContainerHandler(containerService)
	fileHandler := handler.NewFileHandler(fileService)
	jobHandler := handler.NewJobHandler(jobService)

	// 2. Setup router and inject handlers
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type", "Accept"}
	r.Use(cors.New(corsConfig))
	server.RegisterRoutes(r, userHandler, containerHandler, fileHandler, jobHandler, authMiddleware)

	// 3. Start the server
	address := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", address)
	r.Run(address)
}
