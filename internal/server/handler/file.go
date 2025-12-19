package handler

import (
	"container-manager/internal/application"
	"container-manager/internal/errors"
	"net/http"
	"strconv"

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
		_ = c.Error(err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err))
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		_ = c.Error(errors.InternalServerError.Wrap(err))
		return
	}
	defer openedFile.Close()

	if file.Filename == "" {
		_ = c.Error(errors.BadRequest.New("filename cannot be empty"))
		return
	}

	err = h.fileService.UploadFile(userID, file.Filename, openedFile)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
