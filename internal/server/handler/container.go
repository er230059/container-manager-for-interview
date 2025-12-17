package handler

import (
	"container-manager/internal/application"
	containerruntime "container-manager/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContainerHandler struct {
	service *application.ContainerService
}

func NewContainerHandler(service *application.ContainerService) *ContainerHandler {
	return &ContainerHandler{service: service}
}

type createContainerRequest struct {
	Cmd   []string `json:"cmd"`
	Env   []string `json:"env"`
	Image string   `json:"image" binding:"required"`
}

func (h *ContainerHandler) CreateContainer(c *gin.Context) {
	var req createContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.ParseInt(c.GetString("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown user"})
		return
	}

	opts := containerruntime.ContainerCreateOptions{
		Cmd:   req.Cmd,
		Env:   req.Env,
		Image: req.Image,
	}

	id, err := h.service.CreateContainer(c.Request.Context(), userID, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create container"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ContainerHandler) StartContainer(c *gin.Context) {
	id := c.Param("id")

	err := h.service.StartContainer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start container"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ContainerHandler) StopContainer(c *gin.Context) {
	id := c.Param("id")

	err := h.service.StopContainer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stop container"})
		return
	}

	c.Status(http.StatusOK)
}
