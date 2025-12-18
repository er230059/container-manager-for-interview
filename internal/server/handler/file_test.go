package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"container-manager/internal/application"
	"container-manager/internal/infrastructure/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFileHandler_UploadFile(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a temporary directory for testing file storage
	tempDir, err := os.MkdirTemp("", "test_upload_storage")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Initialize dependencies
	localFileStorage := repository.NewLocalFileStorage(tempDir)
	fileService := application.NewFileService(localFileStorage)
	fileHandler := NewFileHandler(fileService)

	// Setup Gin router
	router := gin.Default()
	// Mock authentication middleware to set userID in context
	router.Use(func(c *gin.Context) {
		c.Set("userID", "1234")
		c.Next()
	})
	router.POST("/upload", fileHandler.UploadFile)

	t.Run("successful upload", func(t *testing.T) {
		// Create a buffer to write our form data to
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		// Create a form file
		filename := "testfile.txt"
		fileContent := "This is a test file content."
		part, err := writer.CreateFormFile("file", filename)
		assert.NoError(t, err)
		_, err = part.Write([]byte(fileContent))
		assert.NoError(t, err)
		writer.Close()

		// Create a request
		req, _ := http.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Record the response
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the file was actually saved
		savedFilePath := filepath.Join(tempDir, "1234", filename)
		_, err = os.Stat(savedFilePath)
		assert.NoError(t, err, "uploaded file should exist")
		content, err := os.ReadFile(savedFilePath)
		assert.NoError(t, err)
		assert.Equal(t, fileContent, string(content))
	})

	t.Run("no file provided", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/upload", nil) // No file in body
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get file from form")
	})
}
