package handler

import (
	"leadtek/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler is the handler for home-related endpoints.
// It depends on an application service.
type HomeHandler struct {
	service *application.GreetingService
}

// NewHomeHandler creates a new HomeHandler.
func NewHomeHandler(service *application.GreetingService) *HomeHandler {
	return &HomeHandler{service: service}
}

// GetHome handles the request for the home page.
func (h *HomeHandler) GetHome(c *gin.Context) {
	msg, err := h.service.GetGreeting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}