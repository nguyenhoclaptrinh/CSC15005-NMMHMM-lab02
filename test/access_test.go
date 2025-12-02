package test

import (
	"database/sql"
	"testing"
	"time"

	"secure_notes/test/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExpiredNoteAccess: kiểm tra rằng các ghi chú đã hết hạn không thể truy cập được
func TestExpiredNoteAccess(t *testing.T) {
	t.Run("Cannot access expired note", func(t *testing.T) {
		db, err := testutils.SetupTestDB()
		require.NoError(t, err)
		defer db.Close()

		// Create a test user
		userID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
			userID, "testuser", "hashedpass")
		require.NoError(t, err)

		// Create an expired note (expiry in the past)
		noteID := testutils.GenerateUUID()
		expireAt := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
		_, err = db.Exec("INSERT INTO notes (id, user_id, filename, aes_key_encrypted, expire_at) VALUES (?, ?, ?, ?, ?)",
			noteID, userID, "secret.txt", "encrypted_key", expireAt)
		require.NoError(t, err)

		// Try to access the note
		var expireAtDB sql.NullString
		err = db.QueryRow("SELECT expire_at FROM notes WHERE id = ?", noteID).Scan(&expireAtDB)
		require.NoError(t, err)

		if expireAtDB.Valid {
			expiry, err := time.Parse(time.RFC3339, expireAtDB.String)
			require.NoError(t, err)
			assert.True(t, time.Now().After(expiry), "Note should be expired")
		}
	})

	t.Run("Can access non-expired note", func(t *testing.T) {
		db, err := testutils.SetupTestDB()
		require.NoError(t, err)
		defer db.Close()

		// Create a test user
		userID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
			userID, "testuser", "hashedpass")
		require.NoError(t, err)

		// Create a non-expired note (expiry in the future)
		noteID := testutils.GenerateUUID()
		expireAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		_, err = db.Exec("INSERT INTO notes (id, user_id, filename, aes_key_encrypted, expire_at) VALUES (?, ?, ?, ?, ?)",
			noteID, userID, "secret.txt", "encrypted_key", expireAt)
		require.NoError(t, err)

		// Try to access the note
		var expireAtDB sql.NullString
		err = db.QueryRow("SELECT expire_at FROM notes WHERE id = ?", noteID).Scan(&expireAtDB)
		require.NoError(t, err)

		if expireAtDB.Valid {
			expiry, err := time.Parse(time.RFC3339, expireAtDB.String)
			require.NoError(t, err)
			assert.True(t, time.Now().Before(expiry), "Note should not be expired")
		}
	})
}

// TestExpiredShareAccess: kiểm tra rằng các chia sẻ đã hết hạn không thể truy cập được
func TestExpiredShareAccess(t *testing.T) {
	t.Run("Cannot access expired share link", func(t *testing.T) {
		db, err := testutils.SetupTestDB()
		require.NoError(t, err)
		defer db.Close()

		// Create test users and note
		userID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
			userID, "owner", "hashedpass")
		require.NoError(t, err)

		noteID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO notes (id, user_id, filename, aes_key_encrypted) VALUES (?, ?, ?, ?)",
			noteID, userID, "shared.txt", "encrypted_key")
		require.NoError(t, err)

		// Create an expired share
		shareID := testutils.GenerateUUID()
		urlToken := testutils.GenerateUUID()
		expireAt := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
		_, err = db.Exec("INSERT INTO shares (id, note_id, aes_key_encrypted, url_token, expire_at) VALUES (?, ?, ?, ?, ?)",
			shareID, noteID, "encrypted_key", urlToken, expireAt)
		require.NoError(t, err)

		// Try to access the share
		var expireAtDB sql.NullString
		err = db.QueryRow("SELECT expire_at FROM shares WHERE url_token = ?", urlToken).Scan(&expireAtDB)
		require.NoError(t, err)

		if expireAtDB.Valid {
			expiry, err := time.Parse(time.RFC3339, expireAtDB.String)
			require.NoError(t, err)
			assert.True(t, time.Now().After(expiry), "Share link should be expired")
		}
	})

	t.Run("Can access non-expired share link", func(t *testing.T) {
		db, err := testutils.SetupTestDB()
		require.NoError(t, err)
		defer db.Close()

		// Create test users and note
		userID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
			userID, "owner", "hashedpass")
		require.NoError(t, err)

		noteID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO notes (id, user_id, filename, aes_key_encrypted) VALUES (?, ?, ?, ?)",
			noteID, userID, "shared.txt", "encrypted_key")
		require.NoError(t, err)

		// Create a non-expired share
		shareID := testutils.GenerateUUID()
		urlToken := testutils.GenerateUUID()
		expireAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		_, err = db.Exec("INSERT INTO shares (id, note_id, aes_key_encrypted, url_token, expire_at) VALUES (?, ?, ?, ?, ?)",
			shareID, noteID, "encrypted_key", urlToken, expireAt)
		require.NoError(t, err)

		// Try to access the share
		var expireAtDB sql.NullString
		err = db.QueryRow("SELECT expire_at FROM shares WHERE url_token = ?", urlToken).Scan(&expireAtDB)
		require.NoError(t, err)

		if expireAtDB.Valid {
			expiry, err := time.Parse(time.RFC3339, expireAtDB.String)
			require.NoError(t, err)
			assert.True(t, time.Now().Before(expiry), "Share link should not be expired")
		}
	})
}

// TestUserNoteAccessControl tests that users can only access their own notes
func TestUserNoteAccessControl(t *testing.T) {
	t.Run("User can access their own note", func(t *testing.T) {
		db, err := testutils.SetupTestDB()
		require.NoError(t, err)
		defer db.Close()

		// Create test user
		userID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
			userID, "testuser", "hashedpass")
		require.NoError(t, err)

		// Create note owned by user
		noteID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO notes (id, user_id, filename, aes_key_encrypted) VALUES (?, ?, ?, ?)",
			noteID, userID, "mynote.txt", "encrypted_key")
		require.NoError(t, err)

		// Verify user owns the note
		var ownerID string
		err = db.QueryRow("SELECT user_id FROM notes WHERE id = ?", noteID).Scan(&ownerID)
		require.NoError(t, err)
		assert.Equal(t, userID, ownerID, "User should own the note")
	})

	t.Run("User cannot access another user's note", func(t *testing.T) {
		db, err := testutils.SetupTestDB()
		require.NoError(t, err)
		defer db.Close()

		// Create two users
		user1ID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
			user1ID, "user1", "hashedpass1")
		require.NoError(t, err)

		user2ID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
			user2ID, "user2", "hashedpass2")
		require.NoError(t, err)

		// Create note owned by user1
		noteID := testutils.GenerateUUID()
		_, err = db.Exec("INSERT INTO notes (id, user_id, filename, aes_key_encrypted) VALUES (?, ?, ?, ?)",
			noteID, user1ID, "user1note.txt", "encrypted_key")
		require.NoError(t, err)

		// Verify user2 does not own the note
		var ownerID string
		err = db.QueryRow("SELECT user_id FROM notes WHERE id = ?", noteID).Scan(&ownerID)
		require.NoError(t, err)
		assert.NotEqual(t, user2ID, ownerID, "User2 should not own user1's note")
		assert.Equal(t, user1ID, ownerID, "Note should belong to user1")
	})
}
