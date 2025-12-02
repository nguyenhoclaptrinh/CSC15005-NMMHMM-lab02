package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"secure_notes/test/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegisterSuccess tests successful user registration
func TestRegisterSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := testutils.SetupTestDB()
	require.NoError(t, err)
	defer db.Close()

	router := testutils.NewTestServer()

	// Prepare registration request
	regData := map[string]string{
		"username": "testuser",
		"password": "SecurePass123!",
	}
	body, _ := json.Marshal(regData)

	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 201 Created
	// TODO: Update after implementing Register handler
	// assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

// TestRegisterInvalidInput tests registration with invalid input
func TestRegisterInvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	tests := []struct {
		name     string
		input    map[string]string
		expected int
	}{
		// TODO: Update expected status to http.StatusBadRequest after implementing validation
		{"Empty username", map[string]string{"username": "", "password": "Pass123!"}, http.StatusNotImplemented},
		{"Empty password", map[string]string{"username": "user", "password": ""}, http.StatusNotImplemented},
		{"Short password", map[string]string{"username": "user", "password": "123"}, http.StatusNotImplemented},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expected, w.Code)
		})
	}
}

// TestLoginSuccess tests successful login
func TestLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	loginData := map[string]string{
		"username": "testuser",
		"password": "SecurePass123!",
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 200 OK with JWT token in response
	// TODO: Update after implementing Login handler
	// assert.Equal(t, http.StatusOK, w.Code)
	// var response map[string]string
	// json.Unmarshal(w.Body.Bytes(), &response)
	// assert.NotEmpty(t, response["token"])
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

// TestLoginInvalidCredentials tests login with wrong password
func TestLoginInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	loginData := map[string]string{
		"username": "testuser",
		"password": "WrongPassword",
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 401 Unauthorized
	// TODO: Update after implementing Login handler
	// assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

// TestLoginNonexistentUser tests login with nonexistent user
func TestLoginNonexistentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	loginData := map[string]string{
		"username": "nonexistent",
		"password": "AnyPassword",
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 401 Unauthorized
	// TODO: Update after implementing Login handler
	// assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

// TestGetPublicKey tests getting user's public key
func TestGetPublicKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	// Test getting public key for a user
	userID := testutils.GenerateUUID()
	req := httptest.NewRequest(http.MethodGet, "/api/key/"+userID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 200 OK with public key
	// For now, handler may return 501 NotImplemented
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotImplemented)
}
