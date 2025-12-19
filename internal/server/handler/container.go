package handler

import (
	"container-manager/internal/application"
	"container-manager/internal/domain/infrastructure"
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

// ListContainers godoc
// @Summary List containers
// @Description Lists all containers belonging to the authenticated user
// @Tags Containers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} ContainerResponse
// @Router /containers [get]
func (h *ContainerHandler) ListContainers(c *gin.Context) {
	userID, err := strconv.ParseInt(c.GetString("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown user"})
		return
	}

	containers, err := h.service.ListContainers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list containers"})
		return
	}

	resp := make([]ContainerResponse, 0, len(containers))
	for _, ct := range containers {
		resp = append(resp, ContainerResponse{
			ID:     ct.ID,
			Image:  ct.Image,
			Cmd:    ct.Cmd,
			Env:    ct.Env,
			Status: string(ct.Status),
		})
	}

	c.JSON(http.StatusOK, resp)
}

// CreateContainer godoc
// @Summary Enqueue a new container creation job
// @Description Enqueues a job to create a new container for the authenticated user.
// @Description The job ID is returned immediately, and the status can be tracked via the /jobs/{id} endpoint.
// @Tags Containers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param container body CreateContainerRequest true "Container creation request"
// @Success 200 {object} CreateContainerResponse
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

	opts := infrastructure.ContainerCreateOptions{
		Cmd:   req.Cmd,
		Env:   req.Env,
		Image: req.Image,
	}

	jobID, err := h.service.CreateContainer(c.Request.Context(), userID, opts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue container creation job"})
		return
	}

	c.JSON(http.StatusOK, CreateContainerResponse{JobID: jobID})

}

// StartContainer godoc
// @Summary Start a container
// @Description Starts a specific container for the authenticated user
// @Tags Containers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Container ID"
// @Success 200 "OK"
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
		} else if err.Error() == "conflict operation" {
			c.JSON(http.StatusConflict, gin.H{"error": "conflict operation"})
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
// @Tags Containers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Container ID"
// @Success 200 "OK"
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
		} else if err.Error() == "conflict operation" {
			c.JSON(http.StatusConflict, gin.H{"error": "conflict operation"})
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
// @Tags Containers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Container ID"
// @Success 200 "OK"
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
		} else if err.Error() == "conflict operation" {
			c.JSON(http.StatusConflict, gin.H{"error": "conflict operation"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove container"})
		return
	}

	c.Status(http.StatusOK)
}
