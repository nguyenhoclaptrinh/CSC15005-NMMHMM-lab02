package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

const defaultAPIURL = "http://localhost:8080"

func apiURL() string {
	if v := os.Getenv("API_URL"); v != "" {
		return v
	}
	return defaultAPIURL
}

func tokenPath() string {
	if v := os.Getenv("TOKEN_PATH"); v != "" {
		return v
	}
	return ".client_token"
}

func saveToken(token string) error {
	// For backward compatibility, save raw token if caller used saveToken.
	// New code prefers SaveTokens which writes JSON.
	return os.WriteFile(tokenPath(), []byte(token), 0600)
}

func loadToken() (string, error) {
	b, err := os.ReadFile(tokenPath())
	if err != nil {
		return "", err
	}
	// Try parsing as JSON {"access_token": "...", "refresh_token": "..."}
	var t Tokens
	if err := json.Unmarshal(b, &t); err == nil {
		if t.AccessToken != "" {
			return t.AccessToken, nil
		}
	}
	// Fallback: raw token string
	return string(bytes.TrimSpace(b)), nil
}

// IsLoggedIn reports whether a token file exists (quick check used by CLI)
func IsLoggedIn() bool {
	if _, err := os.Stat(tokenPath()); err == nil {
		return true
	}
	return false
}

// Tokens represents stored auth tokens on the client.
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// SaveTokens writes access and refresh tokens to disk as JSON (0600).
func SaveTokens(t Tokens) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(tokenPath(), b, 0600)
}

// LoadTokens reads stored tokens. Returns error if file missing or unreadable.
// If the token file is a raw string (legacy), it will be returned as AccessToken.
func LoadTokens() (Tokens, error) {
	var t Tokens
	b, err := os.ReadFile(tokenPath())
	if err != nil {
		return t, err
	}
	if err := json.Unmarshal(b, &t); err == nil {
		if t.AccessToken != "" || t.RefreshToken != "" {
			return t, nil
		}
	}
	// Fallback: raw token string
	t.AccessToken = string(bytes.TrimSpace(b))
	return t, nil
}

// doRequest performs an HTTP request and returns the response body and status code.
func doRequest(method, url string, body io.Reader, contentType string, withAuth bool) ([]byte, int, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if withAuth {
		tok, err := loadToken()
		if err == nil && tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return b, resp.StatusCode, err
}

// postJSON helper
func postJSON(path string, payload interface{}, withAuth bool) ([]byte, int, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}
	return doRequest(http.MethodPost, apiURL()+path, bytes.NewReader(b), "application/json", withAuth)
}
