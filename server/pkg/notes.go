package serverpkg

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Global database instance
var db *sql.DB

// InitDB initializes the database connection
func InitDB(dbPath string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}

// UploadNote - Tải lên ghi chú mới (đã mã hóa)
// POST /api/notes
// Request: { "title": "...", "content_enc": "base64...", "key_enc": "base64...", "iv_meta": "{...}" }
// Response: { "id": "note_uuid" }
func UploadNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Title      string `json:"title" binding:"required"`
		ContentEnc string `json:"content_enc" binding:"required"`
		KeyEnc     string `json:"key_enc" binding:"required"`
		IVMeta     string `json:"iv_meta" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := GetDB()

	// Lưu note vào database
	query := `
		INSERT INTO notes (user_id, title_enc, content_enc, key_enc, iv_meta)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, userID, req.Title, req.ContentEnc, req.KeyEnc, req.IVMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save note"})
		return
	}

	noteID, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"id": noteID,
	})
}

// ListNotes - Lấy danh sách ghi chú (chỉ metadata)
// GET /api/notes
// Response: [ { "id": "1", "title": "Encrypted...", "created_at": "..." }, ... ]
func ListNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	db := GetDB()

	// Truy vấn notes của user
	query := `
		SELECT id, title_enc, created_at
		FROM notes
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query notes"})
		return
	}
	defer rows.Close()

	var notes []map[string]interface{}
	for rows.Next() {
		var id, titleEnc, createdAt string

		if err := rows.Scan(&id, &titleEnc, &createdAt); err != nil {
			continue
		}

		notes = append(notes, map[string]interface{}{
			"id":         id,
			"title":      titleEnc,
			"created_at": createdAt,
		})
	}

	if notes == nil {
		notes = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, notes)
}

// GetNote - Tải chi tiết nội dung ghi chú
// GET /api/notes/:id
// Response: { "content_enc": "...", "key_enc": "...", "iv_meta": "..." }
func GetNote(c *gin.Context) {
	noteID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	db := GetDB()

	// Lấy thông tin note và kiểm tra quyền truy cập
	query := `
		SELECT user_id, content_enc, key_enc, iv_meta
		FROM notes
		WHERE id = ?
	`
	var ownerID, contentEnc, keyEnc, ivMeta string

	err := db.QueryRow(query, noteID).Scan(&ownerID, &contentEnc, &keyEnc, &ivMeta)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	// Kiểm tra quyền sở hữu
	if ownerID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content_enc": contentEnc,
		"key_enc":     keyEnc,
		"iv_meta":     ivMeta,
	})
}

// DeleteNote - Xóa ghi chú vĩnh viễn
// DELETE /api/notes/:id
func DeleteNote(c *gin.Context) {
	noteID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	db := GetDB()

	// Kiểm tra quyền sở hữu
	var ownerID string
	err := db.QueryRow("SELECT user_id FROM notes WHERE id = ?", noteID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	if ownerID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can delete"})
		return
	}

	// Xóa note khỏi database
	_, err = db.Exec("DELETE FROM notes WHERE id = ?", noteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "note deleted successfully",
	})
}