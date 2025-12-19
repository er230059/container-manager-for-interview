package handler

import (
	"net/http"
	"strconv"

	"container-manager/internal/application"
	"container-manager/internal/errors"

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

// Job godoc
// @Summary Get job
// @Description Get job by a job ID
// @Tags Jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Job ID"
// @Success 200 {object} GetJobResponse
// @Router /jobs/{id} [get]
func (h *JobHandler) GetJob(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		_ = c.Error(errors.BadRequest.New("job ID is required"))
		return
	}

	userID, err := strconv.ParseInt(c.GetString("userID"), 10, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}

	job, err := h.jobService.GetJob(c.Request.Context(), int64(userID), jobID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := GetJobResponse{
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
