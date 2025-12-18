package repository

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestLocalFileStorage_SaveFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test_file_storage")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Initialize LocalFileStorage
	storage := NewLocalFileStorage(tempDir)

	userID := int64(123)
	filename := "testfile.txt"
	fileContent := "This is a test file content."
	reader := bytes.NewBufferString(fileContent)

	// Save the file
	savedPath, err := storage.SaveFile(userID, filename, reader)
	if err != nil {
		t.Fatalf("SaveFile failed: %v", err)
	}

	// Verify the saved path
	expectedPath := filepath.Join(tempDir, strconv.FormatInt(userID, 10), filename)
	if savedPath != expectedPath {
		t.Errorf("expected saved path %s, got %s", expectedPath, savedPath)
	}

	// Verify the file exists
	if _, err := os.Stat(savedPath); os.IsNotExist(err) {
		t.Errorf("saved file does not exist at %s", savedPath)
	}

	// Verify the content of the saved file
	readContent, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	if string(readContent) != fileContent {
		t.Errorf("expected file content %q, got %q", fileContent, string(readContent))
	}

	// Test case for directory creation
	userID2 := int64(456)
	filename2 := "anotherfile.txt"
	fileContent2 := "Another test content."
	reader2 := bytes.NewBufferString(fileContent2)

	savedPath2, err := storage.SaveFile(userID2, filename2, reader2)
	if err != nil {
		t.Fatalf("SaveFile for new user failed: %v", err)
	}
	expectedPath2 := filepath.Join(tempDir, strconv.FormatInt(userID2, 10), filename2)
	if savedPath2 != expectedPath2 {
		t.Errorf("expected saved path for new user %s, got %s", expectedPath2, savedPath2)
	}
	if _, err := os.Stat(filepath.Join(tempDir, strconv.FormatInt(userID2, 10))); os.IsNotExist(err) {
		t.Errorf("user directory for %d was not created", userID2)
	}
}
