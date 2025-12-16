package main

import (
	"leadtek/internal/application"
	"leadtek/internal/infrastructure/persistence"
	"leadtek/internal/server"
	"leadtek/internal/server/handler"
	"log"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Create dependencies (Composition Root)

	// Infrastructure Layer
	userRepo := persistence.NewInmemUserRepository()

	// ID Generation
	idNode, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("failed to create snowflake node: %v", err)
	}

	// Application Layer
	userService := application.NewUserService(userRepo, idNode)

	// Interfaces Layer
	userHandler := handler.NewUserHandler(userService)

	// 2. Setup router and inject handlers
	r := gin.Default()
	server.RegisterRoutes(r, userHandler) // Pass handlers to router

	// 3. Start the server
	r.Run()
}
