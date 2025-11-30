package config

import (
    "os"
)

// Config holds minimal server configuration used by main.
type Config struct {
    Port   string
    DBPath string
}

// LoadConfig loads configuration from environment variables with sensible defaults.
func LoadConfig() (*Config, error) {
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "secure_notes.db"
    }
    return &Config{Port: port, DBPath: dbPath}, nil
}
