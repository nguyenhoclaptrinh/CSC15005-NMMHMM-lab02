package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"

	serverinternal "secure_notes/server/internalpkg"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// TestKeys holds RSA keypair for testing
type TestKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	PublicPEM  string
}

// GenerateTestRSAKeys creates a fresh RSA keypair for tests.
func GenerateTestRSAKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}

// GenerateTestKeys creates RSA keypair with PEM encoding for tests
func GenerateTestKeys() (*TestKeys, error) {
	priv, pub, err := GenerateTestRSAKeys(2048)
	if err != nil {
		return nil, err
	}

	// Encode public key to PEM
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return &TestKeys{
		PrivateKey: priv,
		PublicKey:  pub,
		PublicPEM:  string(pubPEM),
	}, nil
}

// GenerateUUID generates a simple UUID-like string for testing
func GenerateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// NewInMemoryDB returns an in-memory SQLite DB for tests.
func NewInMemoryDB() (*sql.DB, error) {
	return sql.Open("sqlite3", ":memory:")
}

// SetupTestDB creates an in-memory DB with schema initialized
func SetupTestDB() (*sql.DB, error) {
	db, err := NewInMemoryDB()
	if err != nil {
		return nil, err
	}

	// Create schema
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		public_key TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS notes (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		filename TEXT NOT NULL,
		aes_key_encrypted TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expire_at DATETIME,
		is_deleted INTEGER DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS shares (
		id TEXT PRIMARY KEY,
		note_id TEXT NOT NULL,
		shared_to_user_id TEXT,
		aes_key_encrypted TEXT NOT NULL,
		url_token TEXT UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expire_at DATETIME,
		FOREIGN KEY (note_id) REFERENCES notes(id),
		FOREIGN KEY (shared_to_user_id) REFERENCES users(id)
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
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
// For UNIT TESTS with mock handlers
func NewTestServer() *gin.Engine {
	r := gin.New()

	// Auth endpoints
	r.POST("/api/auth", serverinternal.Register)    // Đăng ký
	r.POST("/api/auth/login", serverinternal.Login) // Đăng nhập

	// Note endpoints
	r.POST("/api/note", serverinternal.UploadNote)       // Upload ghi chú
	r.GET("/api/note", serverinternal.ListNotes)         // Lấy danh sách ghi chú
	r.GET("/api/note/:id", serverinternal.GetNote)       // Lấy ghi chú theo id
	r.DELETE("/api/note/:id", serverinternal.DeleteNote) // Xóa ghi chú
	r.POST("/api/note/share", serverinternal.ShareNote)  // Chia sẻ ghi chú

	// Share endpoints
	r.PATCH("/api/note/:id", serverinternal.RevokeShare) // Hủy share
	r.POST("/api/note/share/", serverinternal.ShareNote) // Tạo share URL (chỉ đích danh)
	r.GET("/api/shares", serverinternal.ListShares)      // Liệt kê shares

	// Public key endpoint
	r.GET("/api/key/:userId", serverinternal.GetPublicKey) // Lấy public key

	return r
}

// NewTestServerWithDB returns a Gin engine with DB injected into context
// For INTEGRATION TESTS with real database
func NewTestServerWithDB(db *sql.DB) *gin.Engine {
	r := gin.New()

	// Middleware: Inject DB vào context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Auth endpoints
	r.POST("/api/auth", serverinternal.Register)
	r.POST("/api/auth/login", serverinternal.Login)

	// Note endpoints
	r.POST("/api/note", serverinternal.UploadNote)
	r.GET("/api/note", serverinternal.ListNotes)
	r.GET("/api/note/:id", serverinternal.GetNote)
	r.DELETE("/api/note/:id", serverinternal.DeleteNote)
	r.POST("/api/note/share", serverinternal.ShareNote)

	// Share endpoints
	r.PATCH("/api/note/:id", serverinternal.RevokeShare)
	r.POST("/api/note/share/", serverinternal.ShareNote)
	r.GET("/api/shares", serverinternal.ListShares)

	// Public key endpoint
	r.GET("/api/key/:userId", serverinternal.GetPublicKey)

	return r
}

// SetupTestDBWithUsers creates DB and seeds test users
// Returns DB and map of usernames to IDs
func SetupTestDBWithUsers() (*sql.DB, map[string]string, error) {
	db, err := SetupTestDB()
	if err != nil {
		return nil, nil, err
	}

	// Generate test users
	aliceID := GenerateUUID()
	bobID := GenerateUUID()

	aliceKeys, _ := GenerateTestKeys()
	bobKeys, _ := GenerateTestKeys()

	// Hash passwords (simplified for testing - trong thực tế dùng bcrypt)
	_, err = db.Exec(`
		INSERT INTO users (id, username, password, public_key) VALUES 
		(?, ?, ?, ?),
		(?, ?, ?, ?)
	`, aliceID, "alice", "$2a$10$alice_hashed_password", aliceKeys.PublicPEM,
		bobID, "bob", "$2a$10$bob_hashed_password", bobKeys.PublicPEM)

	if err != nil {
		return nil, nil, err
	}

	users := map[string]string{
		"alice": aliceID,
		"bob":   bobID,
	}

	return db, users, nil
}
