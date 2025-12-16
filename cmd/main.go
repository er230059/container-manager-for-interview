package main

import (
	"fmt"
	"leadtek/internal/application"
	"leadtek/internal/infrastructure/persistence"
	"leadtek/internal/server"
	"leadtek/internal/server/handler"
	"leadtek/pkg/config"
	"log"

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

	// Infrastructure Layer
	userRepo := persistence.NewInmemUserRepository()

	// ID Generation
	fmt.Println(cfg.Snowflake.MachineID)
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
