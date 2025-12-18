package application

import (
	"container-manager/internal/domain/infrastructure"
	"errors"
	"io"
)

// FileService handles file-related business logic.
type FileService struct {
	fileStorage infrastructure.FileStorage
}

// NewFileService creates a new instance of FileService.
func NewFileService(fs infrastructure.FileStorage) *FileService {
	return &FileService{
		fileStorage: fs,
	}
}

// UploadFile uploads a file for a specific user.
func (s *FileService) UploadFile(userID int64, filename string, fileContent io.Reader) error {
	if filename == "" {
		return errors.New("filename cannot be empty")
	}

	_, err := s.fileStorage.SaveFile(userID, filename, fileContent)
	if err != nil {
		return err
	}

	return nil
}
