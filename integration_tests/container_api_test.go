package integration_tests

import (
	"bytes"
	containerruntime "container-manager/internal/infrastructure/container_runtime"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainerAPI_Integration(t *testing.T) {
	setupTestDB(t)

	// Use Real Docker Runtime
	runtime, err := containerruntime.NewDockerContainerRuntime()
	require.NoError(t, err, "Docker must be available for integration tests")

	r := setupServer(t, runtime)

	t.Run("CreateContainer Success Flow", func(t *testing.T) {
		// 1. Register User
		username := "containeruser"
		password := "password123"

		authHeader := registerAndLogin(t, r, username, password)

		// 2. Create Container via API
		// Use a lightweight image to speed up tests
		payload := map[string]interface{}{
			"image": "alpine:latest",
			"cmd":   []string{"echo", "hello"},
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/containers", bytes.NewBuffer(body))
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		jobID := resp["job_id"]
		require.NotEmpty(t, jobID, "Job ID should be returned")

		// 3. Poll Job Status via API
		var containerID string
		assert.Eventually(t, func() bool {
			reqJob, _ := http.NewRequest("GET", "/jobs/"+jobID, nil)
			reqJob.Header.Set("Authorization", authHeader)
			wJob := httptest.NewRecorder()
			r.ServeHTTP(wJob, reqJob)

			if wJob.Code != http.StatusOK {
				return false
			}

			var jobMap map[string]interface{}
			_ = json.Unmarshal(wJob.Body.Bytes(), &jobMap)

			status, ok := jobMap["status"].(string)
			if !ok {
				return false
			}

			if status == "completed" {
				// Extract container ID from result
				if result, ok := jobMap["result"].(map[string]interface{}); ok {
					if cid, ok := result["container_id"].(string); ok {
						containerID = cid
						return true
					}
				}
				return true
			}
			if status == "failed" {
				t.Logf("Job failed: %v", jobMap["error"])
				return false
			}
			return false
		}, 10*time.Second, 1*time.Second, "Job should complete (pulling image might take time)")

		require.NotEmpty(t, containerID, "Container ID should be in job result")

		// 4. Verify ListContainers
		reqList, _ := http.NewRequest("GET", "/containers", nil)
		reqList.Header.Set("Authorization", authHeader)
		wList := httptest.NewRecorder()
		r.ServeHTTP(wList, reqList)

		require.Equal(t, http.StatusOK, wList.Code)
		var containers []map[string]interface{}
		err = json.Unmarshal(wList.Body.Bytes(), &containers)
		require.NoError(t, err)

		assert.NotEmpty(t, containers)
		found := false
		for _, c := range containers {
			if c["id"] == containerID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created container should be in list")

		// 5. Cleanup
		_ = runtime.Remove(context.Background(), containerID)
	})
}

// Helper to register and return Bearer token
func registerAndLogin(t *testing.T, r http.Handler, username, password string) string {
	// Register
	regPayload := map[string]string{"username": username, "password": password}
	regBody, _ := json.Marshal(regPayload)
	reqReg, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(regBody))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)
	require.Equal(t, http.StatusOK, wReg.Code)

	// Login
	loginBody, _ := json.Marshal(regPayload)
	reqLogin, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)
	require.Equal(t, http.StatusOK, wLogin.Code)

	var loginResp map[string]interface{}
	_ = json.Unmarshal(wLogin.Body.Bytes(), &loginResp)
	token, ok := loginResp["token"].(string)
	require.True(t, ok, "Token not found in login response")

	return "Bearer " + token
}
