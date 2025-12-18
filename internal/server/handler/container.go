package handler

import (
	"container-manager/internal/application"
	containerruntime "container-manager/internal/infrastructure/container_runtime"
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

// @BasePath /

// CreateContainer godoc
// @Summary Create a new container
// @Description Creates a new container for the authenticated user
// @Tags containers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param container body CreateContainerRequest true "Container creation request"
// @Success 200 {object} ContainerIDResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /containers [post]
func (h *ContainerHandler) CreateContainer(c *gin.Context) {
	var req CreateContainerRequest
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

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// StartContainer godoc
// @Summary Start a container
// @Description Starts a specific container for the authenticated user
// @Tags containers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Container ID"
// @Success 200 "OK"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /containers/{id}/start [patch]
func (h *ContainerHandler) StartContainer(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseInt(c.GetString("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown user"})
		return
	}

	err = h.service.StartContainer(c.Request.Context(), userID, id)
	if err != nil {
		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		} else if err.Error() == "container not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start container"})
		return
	}

	c.Status(http.StatusOK)
}

// StopContainer godoc
// @Summary Stop a container
// @Description Stops a specific container for the authenticated user
// @Tags containers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Container ID"
// @Success 200 "OK"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /containers/{id}/stop [patch]
func (h *ContainerHandler) StopContainer(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseInt(c.GetString("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown user"})
		return
	}

	err = h.service.StopContainer(c.Request.Context(), userID, id)
	if err != nil {
		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		} else if err.Error() == "container not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stop container"})
		return
	}

	c.Status(http.StatusOK)
}

// RemoveContainer godoc
// @Summary Remove a container
// @Description Removes a specific container for the authenticated user
// @Tags containers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Container ID"
// @Success 200 "OK"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /containers/{id} [delete]
func (h *ContainerHandler) RemoveContainer(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseInt(c.GetString("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown user"})
		return
	}

	err = h.service.RemoveContainer(c.Request.Context(), userID, id)
	if err != nil {
		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		} else if err.Error() == "container not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove container"})
		return
	}

	c.Status(http.StatusOK)
}
