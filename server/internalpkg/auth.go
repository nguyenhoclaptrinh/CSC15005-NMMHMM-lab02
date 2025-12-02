package serverinternal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	// Nhận username, password
	// Hash password (bcrypt/argon2)
	// Lưu vào DB
	// Trả user_id hoặc lỗi
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// Đăng nhập user
func Login(c *gin.Context) {
	// Nhận username, password
	// So sánh hash password
	// Tạo JWT token
	// Trả token hoặc lỗi
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// GetPublicKey lấy public key của user
func GetPublicKey(c *gin.Context) {
	// Lấy userId từ URL param
	// Query public_key từ DB
	// Trả public key dưới dạng PEM
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
