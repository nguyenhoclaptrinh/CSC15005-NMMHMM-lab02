package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"os"

	serverinternal "secure_notes/server/internalpkg"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// GenerateTestRSAKeys creates a fresh RSA keypair for tests.
func GenerateTestRSAKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}

// NewInMemoryDB returns an in-memory SQLite DB for tests.
func NewInMemoryDB() (*sql.DB, error) {
	return sql.Open("sqlite3", ":memory:")
}

// NewTempStorage creates a temporary directory to act as storage for tests.
func NewTempStorage(prefix string) (string, error) {
	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", err
	}
	return dir, nil
}

// CleanupTempStorage removes the temp storage directory.
func CleanupTempStorage(dir string) error {
	return os.RemoveAll(dir)
}

// NewTestServer returns a Gin engine with basic routes wired to the server handlers.
func NewTestServer() *gin.Engine {
	r := gin.New()
	// register endpoints used in tests
	r.POST("/register", serverinternal.Register)
	r.POST("/login", serverinternal.Login)
	r.POST("/notes", serverinternal.UploadNote)
	r.GET("/notes", serverinternal.ListNotes)
	r.GET("/notes/:id", serverinternal.GetNote)
	r.DELETE("/notes/:id", serverinternal.DeleteNote)
	r.POST("/shares", serverinternal.ShareNote)
	r.GET("/shares", serverinternal.ListShares)
	r.DELETE("/shares/:id", serverinternal.RevokeShare)
	return r
}
