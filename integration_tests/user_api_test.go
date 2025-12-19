package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAPI_Integration(t *testing.T) {
	setupTestDB(t)

	// Runtime not needed for user tests
	r := setupServer(t, nil)

	t.Run("Register and Login", func(t *testing.T) {
		username := "apiuser"
		password := "password123"

		// 1. Register
		registerPayload := map[string]string{
			"username": username,
			"password": password,
		}
		body, _ := json.Marshal(registerPayload)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var registerResp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &registerResp)
		require.NoError(t, err)
		assert.NotEmpty(t, registerResp["id"])
		assert.Equal(t, username, registerResp["username"])

		// 2. Login
		loginPayload := map[string]string{
			"username": username,
			"password": password,
		}
		body, _ = json.Marshal(loginPayload)
		req, _ = http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var loginResp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &loginResp)
		require.NoError(t, err)
		assert.NotEmpty(t, loginResp["token"])
	})

	t.Run("CreateDuplicateUser", func(t *testing.T) {
		username := "dupapiuser"
		password := "pass"

		payload := map[string]string{
			"username": username,
			"password": password,
		}
		body, _ := json.Marshal(payload)

		// Create first
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Create second
		req2, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		// Expect failure (likely 500 or 409)
		assert.NotEqual(t, http.StatusOK, w2.Code)
	})
}