package serverpkg

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================
// SHARE URL APIs - Temporary/Anonymous Share Links
// ============================================================

// CreateShareLink - Tạo Link Chia sẻ
// POST /api/share
// Request: { "content_enc": "base64...", "metadata": { "expires_in": 3600, "max_views": 5, "has_password": true, "access_hash": "sha256..." } }
// Response: { "share_id": "uuid-1234...", "expires_at": "2025-12-31T23:59:00Z" }
func CreateShareLink(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		ContentEnc string `json:"content_enc" binding:"required"`
		Metadata   struct {
			ExpiresIn   int    `json:"expires_in"`   // Seconds
			MaxViews    int    `json:"max_views"`    // 0 = unlimited
			HasPassword bool   `json:"has_password"` // true nếu có password
			AccessHash  string `json:"access_hash"`  // SHA256 hash của password (nếu có)
		} `json:"metadata" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra access_hash nếu has_password = true
	if req.Metadata.HasPassword && req.Metadata.AccessHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access_hash required when has_password is true"})
		return
	}

	db := GetDB()

	// Tính thời gian hết hạn
	var expiresAt *string
	if req.Metadata.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(req.Metadata.ExpiresIn) * time.Second).Format(time.RFC3339)
		expiresAt = &expiry
	}

	// Convert boolean to integer for SQLite
	hasPassword := 0
	if req.Metadata.HasPassword {
		hasPassword = 1
	}

	// Xử lý max_views (NULL nếu unlimited)
	var maxViews interface{}
	if req.Metadata.MaxViews > 0 {
		maxViews = req.Metadata.MaxViews
	} else {
		maxViews = nil
	}

	// Lưu vào shared_links
	query := `
		INSERT INTO shared_links (owner_id, content_enc, expires_at, max_views, has_password, access_hash)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, userID, req.ContentEnc, expiresAt, maxViews, hasPassword, req.Metadata.AccessHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create share link"})
		return
	}

	shareID, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"share_id":   shareID,
		"expires_at": expiresAt,
	})
}

// GetShareInfo - Lấy Thông tin Link
// GET /api/share/:id/info
// Response: { "is_active": true, "requires_password": true, "expires_at": "..." }
func GetShareInfo(c *gin.Context) {
	shareID := c.Param("id")

	db := GetDB()

	query := `
		SELECT expires_at, max_views, current_views, has_password, is_active
		FROM shared_links
		WHERE id = ?
	`

	var expiresAt *string
	var maxViews *int
	var currentViews, hasPassword, isActive int

	err := db.QueryRow(query, shareID).Scan(&expiresAt, &maxViews, &currentViews, &hasPassword, &isActive)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "share link not found"})
		return
	}

	// Kiểm tra trạng thái active
	active := isActive == 1

	// Kiểm tra hết hạn
	if expiresAt != nil {
		expiry, _ := time.Parse(time.RFC3339, *expiresAt)
		if time.Now().After(expiry) {
			active = false
		}
	}

	// Kiểm tra max views
	if maxViews != nil && currentViews >= *maxViews {
		active = false
	}

	c.JSON(http.StatusOK, gin.H{
		"is_active":         active,
		"requires_password": hasPassword == 1,
		"expires_at":        expiresAt,
	})
}

// GetSharedContent - Truy cập & Tải File
// GET /api/share/:id
// Response: { "content_enc": "base64_string" }
func GetSharedContent(c *gin.Context) {
	shareID := c.Param("id")

	// Lấy password hash từ header (nếu có)
	providedHash := c.GetHeader("X-Access-Pass-Hash")

	db := GetDB()

	query := `
		SELECT content_enc, expires_at, max_views, current_views, has_password, access_hash, is_active
		FROM shared_links
		WHERE id = ?
	`

	var contentEnc string
	var expiresAt *string
	var maxViews *int
	var currentViews, hasPassword, isActive int
	var accessHash *string

	err := db.QueryRow(query, shareID).Scan(&contentEnc, &expiresAt, &maxViews, &currentViews, &hasPassword, &accessHash, &isActive)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "share link not found"})
		return
	}

	// Kiểm tra is_active
	if isActive == 0 {
		c.JSON(http.StatusGone, gin.H{"error": "link has been revoked"})
		return
	}

	// Kiểm tra hết hạn
	if expiresAt != nil {
		expiry, _ := time.Parse(time.RFC3339, *expiresAt)
		if time.Now().After(expiry) {
			c.JSON(http.StatusGone, gin.H{"error": "link has expired"})
			return
		}
	}

	// Kiểm tra max views
	if maxViews != nil && currentViews >= *maxViews {
		c.JSON(http.StatusGone, gin.H{"error": "link has reached maximum views"})
		return
	}

	// Kiểm tra password
	if hasPassword == 1 {
		if providedHash == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "password required"})
			return
		}

		if accessHash != nil && *accessHash != providedHash {
			c.JSON(http.StatusForbidden, gin.H{"error": "incorrect password"})
			return
		}
	}

	// Tăng current_views
	db.Exec("UPDATE shared_links SET current_views = current_views + 1, last_accessed_at = datetime('now') WHERE id = ?", shareID)

	c.JSON(http.StatusOK, gin.H{
		"content_enc": contentEnc,
	})
}

// RevokeShareLink - Hủy Chia sẻ
// DELETE /api/share/:id
// Response: { "message": "Link revoked successfully" }
func RevokeShareLink(c *gin.Context) {
	shareID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	db := GetDB()

	// Kiểm tra quyền sở hữu
	var ownerID string
	err := db.QueryRow("SELECT owner_id FROM shared_links WHERE id = ?", shareID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "share link not found"})
		return
	}

	if ownerID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can revoke link"})
		return
	}

	// Đánh dấu is_active = 0
	_, err = db.Exec("UPDATE shared_links SET is_active = 0 WHERE id = ?", shareID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "link revoked successfully",
	})
}