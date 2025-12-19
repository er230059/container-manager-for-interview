package application

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"container-manager/internal/application/mocks"

	"go.uber.org/mock/gomock"
)

func TestFileService_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	fileService := NewFileService(mockFileStorage)

	userID := int64(1000)
	filename := "test_file.txt"
	fileContent := "hello world"
	reader := bytes.NewBufferString(fileContent)
	expectedPath := "/path/to/uploaded/file"

	t.Run("success", func(t *testing.T) {
		mockFileStorage.EXPECT().SaveFile(userID, filename, gomock.Any()).Return(expectedPath, nil)

		err := fileService.UploadFile(userID, filename, reader)
		if err != nil {
			t.Errorf("UploadFile returned an error: %v", err)
		}
	})

	t.Run("empty filename", func(t *testing.T) {
		err := fileService.UploadFile(userID, "", reader)
		if err == nil {
			t.Error("UploadFile did not return an error for empty filename")
		}
		if err.Error() != "filename cannot be empty" {
			t.Errorf("UploadFile returned wrong error for empty filename: got %q, want %q", err.Error(), "bad request")
		}
	})

	t.Run("fileStorage SaveFile error", func(t *testing.T) {
		mockError := errors.New("storage error")
		mockFileStorage.EXPECT().SaveFile(userID, filename, gomock.Any()).Return("", mockError)

		err := fileService.UploadFile(userID, filename, reader)
		if err == nil {
			t.Error("UploadFile did not return an error when SaveFile fails")
		}
		if err.Error() != "storage error" {
			t.Errorf("UploadFile returned wrong error when SaveFile fails: got %q, want %q", err.Error(), "failed to save file: storage error")
		}
	})

	t.Run("gomock.Any matches io.Reader", func(t *testing.T) {
		// Ensure gomock.Any() correctly matches io.Reader.
		// A new reader is needed because the previous one might have been consumed.
		newReader := bytes.NewBufferString(fileContent)
		mockFileStorage.EXPECT().SaveFile(userID, filename, gomock.Any()).DoAndReturn(func(_ int64, _ string, r io.Reader) (string, error) {
			readBytes, err := io.ReadAll(r)
			if err != nil {
				return "", err
			}
			if string(readBytes) != fileContent {
				t.Errorf("mock received wrong content: got %q, want %q", string(readBytes), fileContent)
			}
			return expectedPath, nil
		})

		err := fileService.UploadFile(userID, filename, newReader)
		if err != nil {
			t.Errorf("UploadFile returned an error: %v", err)
		}
	})
}
