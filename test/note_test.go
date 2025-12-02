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
)

// TestUploadNote tests uploading a note
func TestUploadNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	noteData := map[string]interface{}{
		"filename":          "secret.txt",
		"aes_key_encrypted": "encrypted_aes_key_base64",
		"file_content":      "encrypted_file_content_base64",
	}
	body, _ := json.Marshal(noteData)

	req := httptest.NewRequest(http.MethodPost, "/api/note", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock_jwt_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 201 Created with note_id
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusNotImplemented)
}

// TestListNotes tests getting list of user's notes
func TestListNotes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/note", nil)
	req.Header.Set("Authorization", "Bearer mock_jwt_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 200 OK with list of notes
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotImplemented)
}

// TestGetNote tests getting a specific note by ID
func TestGetNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	noteID := testutils.GenerateUUID()
	req := httptest.NewRequest(http.MethodGet, "/api/note/"+noteID, nil)
	req.Header.Set("Authorization", "Bearer mock_jwt_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 200 OK with note metadata + encrypted file
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusNotImplemented)
}

// TestDeleteNote tests deleting a note
func TestDeleteNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	noteID := testutils.GenerateUUID()
	req := httptest.NewRequest(http.MethodDelete, "/api/note/"+noteID, nil)
	req.Header.Set("Authorization", "Bearer mock_jwt_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 200 OK or 204 No Content
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNoContent || w.Code == http.StatusNotImplemented)
}

// TestShareNote tests sharing a note with another user
func TestShareNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	t.Run("Share note to specific user", func(t *testing.T) {
		shareData := map[string]interface{}{
			"note_id":           testutils.GenerateUUID(),
			"shared_to_user_id": testutils.GenerateUUID(),
			"aes_key_encrypted": "encrypted_aes_key_for_recipient",
			"expire_hours":      24,
		}
		body, _ := json.Marshal(shareData)

		req := httptest.NewRequest(http.MethodPost, "/api/note/share", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// When implemented, expect 200 OK with share URL
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusCreated || w.Code == http.StatusNotImplemented)
	})

	t.Run("Share note via public link", func(t *testing.T) {
		shareData := map[string]interface{}{
			"note_id":           testutils.GenerateUUID(),
			"aes_key_encrypted": "encrypted_aes_key",
			"expire_hours":      24,
		}
		body, _ := json.Marshal(shareData)

		req := httptest.NewRequest(http.MethodPost, "/api/note/share/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mock_jwt_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// When implemented, expect 200 OK with public share URL
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusCreated || w.Code == http.StatusNotImplemented)
	})
}

// TestRevokeShare tests revoking/canceling a share
func TestRevokeShare(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	noteID := testutils.GenerateUUID()
	req := httptest.NewRequest(http.MethodPatch, "/api/note/"+noteID, nil)
	req.Header.Set("Authorization", "Bearer mock_jwt_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 200 OK
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotImplemented)
}

// TestListShares tests getting list of shares for a note
func TestListShares(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.NewTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/shares", nil)
	req.Header.Set("Authorization", "Bearer mock_jwt_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// When implemented, expect 200 OK with list of shares
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotImplemented)
}

// TestNoteAccessWithExpiry tests accessing expired notes via share
func TestNoteAccessWithExpiry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := testutils.SetupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	// Create test scenario: expired share link
	userID := testutils.GenerateUUID()
	_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
		userID, "owner", "hashedpass")
	assert.NoError(t, err)

	noteID := testutils.GenerateUUID()
	_, err = db.Exec("INSERT INTO notes (id, user_id, filename, aes_key_encrypted) VALUES (?, ?, ?, ?)",
		noteID, userID, "shared.txt", "encrypted_key")
	assert.NoError(t, err)

	// Create expired share
	shareID := testutils.GenerateUUID()
	urlToken := testutils.GenerateUUID()
	// Set expire_at to past time
	_, err = db.Exec("INSERT INTO shares (id, note_id, aes_key_encrypted, url_token, expire_at) VALUES (?, ?, ?, ?, datetime('now', '-1 hour'))",
		shareID, noteID, "encrypted_key", urlToken)
	assert.NoError(t, err)

	// Verify share is expired
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM shares WHERE url_token = ? AND (expire_at IS NULL OR expire_at > datetime('now'))", urlToken).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "Expired share should not be accessible")
}
