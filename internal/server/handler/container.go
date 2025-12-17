package handler

import (
	"container-manager/internal/application"
	containerruntime "container-manager/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContainerHandler struct {
	service *application.ContainerService
}

func NewContainerHandler(service *application.ContainerService) *ContainerHandler {
	return &ContainerHandler{service: service}
}

type createContainerRequest struct {
	Image string   `json:"image" binding:"required"`
	Env   []string `json:"env"`
}

func (h *ContainerHandler) CreateContainer(c *gin.Context) {
	var req createContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	opts := containerruntime.ContainerCreateOptions{
		Image: req.Image,
		Env:   req.Env,
	}

	id, err := h.service.CreateContainer(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create container"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
