package main

import (
	"fmt"
	"log"

	"container-manager/internal/application"
	"container-manager/internal/infrastructure/database/sql"
	"container-manager/internal/infrastructure/repository"
	"container-manager/internal/server"
	"container-manager/internal/server/handler"
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

	userDB := sql.NewUserDatabase(db)
	userRepo := repository.NewUserRepository(userDB)

	// ID Generation
	idNode, err := snowflake.NewNode(cfg.Snowflake.MachineID)
	if err != nil {
		log.Fatalf("failed to create snowflake node: %v", err)
	}

	// Application Layer
	userService := application.NewUserService(userRepo, idNode)

	// Interfaces Layer
	userHandler := handler.NewUserHandler(userService)

	// 2. Setup router and inject handlers
	r := gin.Default()
	server.RegisterRoutes(r, userHandler)

	// 3. Start the server
	address := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", address)
	r.Run(address)
}
