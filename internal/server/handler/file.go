package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"container-manager/internal/application"

	"github.com/gin-gonic/gin"
)

// FileHandler handles file-related HTTP requests.
type FileHandler struct {
	fileService *application.FileService
}

// NewFileHandler creates a new instance of FileHandler.
func NewFileHandler(fs *application.FileService) *FileHandler {
	return &FileHandler{
		fileService: fs,
	}
}

// UploadFile handles file upload requests.
// @Summary Upload file
// @Description Uploads a file to the user's dedicated storage folder.
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "File to upload"
// @Success 200 "OK"
// @Router /files [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	userID, err := strconv.ParseInt(c.GetString("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown user"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to get file from form: %v", err)})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to open file: %v", err)})
		return
	}
	defer openedFile.Close()

	err = h.fileService.UploadFile(userID, file.Filename, openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to upload file: %v", err)})
		return
	}

	c.Status(http.StatusOK)
}
