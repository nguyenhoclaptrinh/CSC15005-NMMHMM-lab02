package serverinternal

import (
    "github.com/gin-gonic/gin"
    "net/http"
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
