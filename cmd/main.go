package main

import (
	"leadtek/internal/application"
	"leadtek/internal/infrastructure/persistence"
	"leadtek/internal/server"
	"leadtek/internal/server/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Create dependencies (Composition Root)
	// Infrastructure Layer
	greetingRepo := persistence.NewInmemGreetingRepository()

	// Application Layer
	greetingService := application.NewGreetingService(greetingRepo)

	// Interfaces Layer
	homeHandler := handler.NewHomeHandler(greetingService)

	// 2. Setup router and inject handlers
	r := gin.Default()
	server.RegisterRoutes(r, homeHandler) // Pass handler to router

	// 3. Start the server
	r.Run()
}
