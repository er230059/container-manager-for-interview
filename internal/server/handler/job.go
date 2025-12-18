package handler

import (
	"net/http"
	"strconv"

	"container-manager/internal/application"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobService application.JobService
}

func NewJobHandler(jobService application.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

func (h *JobHandler) GetJob(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	userID, err := strconv.ParseInt(c.GetString("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown user"})
		return
	}

	job, err := h.jobService.GetJob(c.Request.Context(), int64(userID), jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	response := JobStatusResponse{
		ID:        job.ID,
		Type:      job.Type,
		Status:    string(job.Status),
		Result:    job.Result,
		Error:     job.Error,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}
