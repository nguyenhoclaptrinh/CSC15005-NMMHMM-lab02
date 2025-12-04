package pkg

import "github.com/gin-gonic/gin"

// Chia sẻ ghi chú cho user khác
func ShareNote(c *gin.Context) {
	// Nhận note_id, shared_to_user_id, aes_key_encrypted
	// Lưu vào note_shares table
	c.JSON(200, gin.H{"error": "not implemented"})
}

// Liệt kê các chia sẻ của ghi chú
func ListShares(c *gin.Context) {
	// Lấy note_id
	// Truy vấn note_shares
	// Trả danh sách chia sẻ
	c.JSON(200, gin.H{"error": "not implemented"})
}

// Thu hồi chia sẻ
func RevokeShare(c *gin.Context) {
	// TODO: implement revoke
	c.JSON(200, gin.H{"error": "not implemented"})
}
