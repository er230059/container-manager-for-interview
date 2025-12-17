package main

import (
	"fmt"
	"log"

	"container-manager/internal/application"
	containerruntime "container-manager/internal/infrastructure/container_runtime"
	"container-manager/internal/infrastructure/repository"
	"container-manager/internal/server"
	"container-manager/internal/server/handler"
	"container-manager/internal/server/middleware"
	"container-manager/pkg/config"
	"container-manager/pkg/postgres"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
)

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

	userRepo := repository.NewUserDatabase(db)
	containerRepo := repository.NewContainerUserDatabase(db)

	// Infrastructure Layer - Container Runtime
	runtime, err := containerruntime.NewDockerContainerRuntime()
	if err != nil {
		log.Fatalf("failed to create container runtime: %v", err)
	}

	// ID Generation
	idNode, err := snowflake.NewNode(cfg.Snowflake.MachineID)
	if err != nil {
		log.Fatalf("failed to create snowflake node: %v", err)
	}

	// Application Layer
	userService := application.NewUserService(userRepo, idNode, cfg.Server.JWTSecret)
	containerService := application.NewContainerService(runtime, containerRepo)

	// Handler Layer
	authMiddleware := middleware.NewAuthMiddleware(cfg.Server.JWTSecret)
	userHandler := handler.NewUserHandler(userService)
	containerHandler := handler.NewContainerHandler(containerService)

	// 2. Setup router and inject handlers
	r := gin.Default()
	server.RegisterRoutes(r, userHandler, containerHandler, authMiddleware)

	// 3. Start the server
	address := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", address)
	r.Run(address)
}
