package serverpkg

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Register prompts and calls server register endpoint
func Register() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Client-side validation (quick checks to improve UX)
	if err := ValidateInput(username); err != nil {
		fmt.Println("invalid username:", err)
		return
	}
	if err := ValidateInput(password); err != nil {
		fmt.Println("invalid password:", err)
		return
	}

	payload := map[string]string{"username": username, "password": password}
	b, status, err := postJSON("/api/register", payload, false)
	if err != nil {
		LogError("register request failed", err)
		return
	}
	LogInfo(fmt.Sprintf("register status: %d", status))
	fmt.Println(string(b))
}

// Login prompts and calls server login endpoint and stores access token
func Login() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Client-side validation
	if err := ValidateInput(username); err != nil {
		fmt.Println("invalid username:", err)
		return
	}
	if err := ValidateInput(password); err != nil {
		fmt.Println("invalid password:", err)
		return
	}

	payload := map[string]string{"username": username, "password": password}
	b, status, err := postJSON("/api/login", payload, false)
	if err != nil {
		LogError("login request failed", err)
		return
	}
	if status != 200 && status != 201 {
		LogInfo(fmt.Sprintf("login failed: %d", status))
		fmt.Println(string(b))
		return
	}
	var resp map[string]any
	if err := json.Unmarshal(b, &resp); err != nil {
		LogError("invalid login response", err)
		return
	}
	var tokens Tokens
	if raw, ok := resp["access_token"]; ok {
		if tok, ok := raw.(string); ok && tok != "" {
			tokens.AccessToken = tok
		}
	}
	if raw, ok := resp["refresh_token"]; ok {
		if tok, ok := raw.(string); ok && tok != "" {
			tokens.RefreshToken = tok
		}
	}
	if tokens.AccessToken != "" || tokens.RefreshToken != "" {
		if err := SaveTokens(tokens); err != nil {
			LogError("failed to save tokens", err)
			return
		}
		LogInfo("tokens saved")
	}
	fmt.Println(string(b))
}

// UploadNote uploads a file to /api/notes using multipart form
func UploadNote() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("File path: ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)
	if path == "" {
		LogInfo("no file provided")
		return
	}
	// Basic file validation: existence and size limit (50 MB)
	fi, err := os.Stat(path)
	if err != nil {
		LogError("stat file", err)
		return
	}
	const maxSize = 50 * 1024 * 1024 // 50 MB
	if fi.Size() > maxSize {
		LogInfo("file too large (max 50 MB)")
		return
	}

	file, err := os.Open(path)
	if err != nil {
		LogError("open file", err)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		LogError("create form file", err)
		return
	}
	if _, err := io.Copy(fw, file); err != nil {
		LogError("copy file", err)
		return
	}
	// optional title
	fmt.Print("Title (optional): ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)
	if title != "" {
		_ = mw.WriteField("title", title)
	}
	mw.Close()

	respBody, status, err := doRequest(http.MethodPost, apiURL()+"/api/notes", &buf, mw.FormDataContentType(), true)
	if err != nil {
		LogError("upload failed", err)
		return
	}
	LogInfo(fmt.Sprintf("upload status: %d", status))
	fmt.Println(string(respBody))
}

// ListNotes retrieves notes for the authenticated user
func ListNotes() {
	b, status, err := doRequest(http.MethodGet, apiURL()+"/api/notes", nil, "", true)
	if err != nil {
		LogError("list notes failed", err)
		return
	}
	LogInfo(fmt.Sprintf("list notes status: %d", status))
	fmt.Println(string(b))
}

// ShareNote creates a share for a note
func ShareNote() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Note ID: ")
	noteID, _ := reader.ReadString('\n')
	noteID = strings.TrimSpace(noteID)
	fmt.Print("Recipient ID: ")
	recID, _ := reader.ReadString('\n')
	recID = strings.TrimSpace(recID)
	if noteID == "" || recID == "" {
		LogInfo("note ID and recipient are required")
		return
	}
	fmt.Print("Expiry (e.g. 24h) or empty: ")
	expiry, _ := reader.ReadString('\n')
	expiry = strings.TrimSpace(expiry)
	payload := map[string]string{"recipient_id": recID}
	if expiry != "" {
		payload["expiry"] = expiry
	}
	path := "/api/notes/" + noteID + "/share"
	b, status, err := postJSON(path, payload, true)
	if err != nil {
		LogError("share failed", err)
		return
	}
	LogInfo(fmt.Sprintf("share status: %d", status))
	fmt.Println(string(b))
}

// CreateTempURL creates a temporary share URL for a note (server-dependent)
func CreateTempURL() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Note ID: ")
	noteID, _ := reader.ReadString('\n')
	noteID = strings.TrimSpace(noteID)
	if noteID == "" {
		LogInfo("note ID required")
		return
	}
	fmt.Print("Expiry (e.g. 1h) or empty: ")
	expiry, _ := reader.ReadString('\n')
	expiry = strings.TrimSpace(expiry)
	payload := map[string]string{"note_id": noteID}
	if expiry != "" {
		payload["expiry"] = expiry
	}
	b, status, err := postJSON("/api/share/create", payload, true)
	if err != nil {
		LogError("create temp url failed", err)
		return
	}
	LogInfo(fmt.Sprintf("create temp url status: %d", status))
	fmt.Println(string(b))
}

// Logout calls server logout endpoint and clears local token
func Logout() {
	b, status, err := doRequest(http.MethodPost, apiURL()+"/api/auth/logout", nil, "", true)
	if err != nil {
		LogError("logout request failed", err)
		return
	}
	LogInfo(fmt.Sprintf("logout status: %d", status))
	fmt.Println(string(b))
	// Remove saved token regardless of server response
	if err := os.Remove(tokenPath()); err != nil {
		// If file not found, ignore
		if !os.IsNotExist(err) {
			LogError("failed to remove token file", err)
		}
	} else {
		LogInfo("local token removed")
	}
}
