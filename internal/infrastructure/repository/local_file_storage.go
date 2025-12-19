package repository

import (
	"container-manager/internal/domain/infrastructure"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

var _ infrastructure.FileStorage = (*LocalFileStorage)(nil)

// LocalFileStorage implements the FileStorage interface for local disk storage.
type LocalFileStorage struct {
	basePath string
}

// NewLocalFileStorage creates a new instance of LocalFileStorage.
func NewLocalFileStorage(basePath string) infrastructure.FileStorage {
	return &LocalFileStorage{
		basePath: basePath,
	}
}

// SaveFile saves the given file content to the local disk within the user's folder.
func (s *LocalFileStorage) SaveFile(userID int64, filename string, fileContent io.Reader) (string, error) {
	userDir := filepath.Join(s.basePath, strconv.FormatInt(userID, 10))

	// Create user directory if it doesn't exist
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", err
	}

	filePath := filepath.Join(userDir, filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, fileContent)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
