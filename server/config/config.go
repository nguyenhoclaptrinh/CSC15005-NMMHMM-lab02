package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Config holds server configuration
type Config struct {
	Port   string
	DBPath string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() (*Config, error) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./server/database/secure_notes.db"
	}
	return &Config{Port: port, DBPath: dbPath}, nil
}

// InitDB initializes the database and creates tables if they don't exist
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("✅ Database initialized successfully")
	return db, nil
}

// createTables creates all necessary tables
func createTables(db *sql.DB) error {
	schema := `
	-- =====================================================
	-- TABLE: users
	-- =====================================================
	CREATE TABLE IF NOT EXISTS users (
		id              TEXT PRIMARY KEY,
		username        VARCHAR(50) NOT NULL UNIQUE,
		password_hash   VARCHAR(255) NOT NULL,
		kdf_salt        VARCHAR(64) NOT NULL,
		created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

	-- =====================================================
	-- TABLE: user_keys
	-- =====================================================
	CREATE TABLE IF NOT EXISTS user_keys (
		user_id         TEXT PRIMARY KEY,
		public_key      TEXT NOT NULL,
		created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- =====================================================
	-- TABLE: notes
	-- =====================================================
	CREATE TABLE IF NOT EXISTS notes (
		id              TEXT PRIMARY KEY,
		user_id         TEXT NOT NULL,
		title_enc       TEXT NOT NULL,
		content_enc     BLOB NOT NULL,
		key_enc         TEXT NOT NULL,
		iv_meta         TEXT NOT NULL,
		created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
	CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);

	-- =====================================================
	-- TABLE: shared_links
	-- =====================================================
	CREATE TABLE IF NOT EXISTS shared_links (
		id                  TEXT PRIMARY KEY,
		owner_id            TEXT NOT NULL,
		content_enc         BLOB NOT NULL,
		sender_public_key   TEXT,
		expires_at          TIMESTAMP,
		max_views           INTEGER,
		current_views       INTEGER DEFAULT 0,
		has_password        BOOLEAN DEFAULT FALSE,
		access_hash         VARCHAR(64),
		is_active           BOOLEAN DEFAULT TRUE,
		created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_accessed_at    TIMESTAMP,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_shared_links_owner ON shared_links(owner_id);
	CREATE INDEX IF NOT EXISTS idx_shared_links_expires ON shared_links(expires_at);
	CREATE INDEX IF NOT EXISTS idx_shared_links_active ON shared_links(is_active);

	-- =====================================================
	-- TABLE: refresh_tokens
	-- =====================================================
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id              TEXT PRIMARY KEY,
		user_id         TEXT NOT NULL,
		token_hash      VARCHAR(64) NOT NULL,
		expires_at      TIMESTAMP NOT NULL,
		is_revoked      BOOLEAN DEFAULT FALSE,
		created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		revoked_at      TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON refresh_tokens(expires_at);

	-- =====================================================
	-- TABLE: token_blacklist
	-- =====================================================
	CREATE TABLE IF NOT EXISTS token_blacklist (
		jti             VARCHAR(36) PRIMARY KEY,
		expires_at      TIMESTAMP NOT NULL,
		blacklisted_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		reason          VARCHAR(50)
	);
	CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires ON token_blacklist(expires_at);
	`

	_, err := db.Exec(schema)
	return err
}

// CleanupExpiredData removes expired tokens and links (run periodically)
func CleanupExpiredData(db *sql.DB) error {
	queries := []string{
		// Xóa shared_links đã hết hạn
		`DELETE FROM shared_links WHERE expires_at < CURRENT_TIMESTAMP AND expires_at IS NOT NULL`,
		
		// Xóa refresh_tokens đã hết hạn
		`DELETE FROM refresh_tokens WHERE expires_at < CURRENT_TIMESTAMP`,
		
		// Xóa JWT blacklist đã hết hạn
		`DELETE FROM token_blacklist WHERE expires_at < CURRENT_TIMESTAMP`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("⚠️  Cleanup warning: %v", err)
		}
	}

	log.Println("🧹 Cleanup completed")
	return nil
}
