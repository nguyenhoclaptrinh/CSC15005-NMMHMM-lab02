package serverpkg

import "github.com/gin-gonic/gin"

// Upload ghi chú
func UploadNote(c *gin.Context) {
	// Nhận file (multipart), metadata, aes_key_encrypted
	// Lưu file vào storage/
	// Lưu metadata vào SQLite (notes table)
	// Trả note_id
	c.JSON(200, gin.H{"error": "not implemented"})
}

// Lấy danh sách ghi chú
func ListNotes(c *gin.Context) {
	// Lấy user_id từ JWT
	// Truy vấn notes table
	// Trả danh sách note metadata
	c.JSON(200, gin.H{"error": "not implemented"})
}

// Lấy chi tiết ghi chú
func GetNote(c *gin.Context) {
	// Lấy note_id từ URL
	// Truy vấn SQLite
	// Trả metadata + file đã mã hóa
	c.JSON(200, gin.H{"error": "not implemented"})
}

// Xóa ghi chú
func DeleteNote(c *gin.Context) {
	// Đánh dấu is_deleted = true trong DB
	// Có thể xóa file storage/
	c.JSON(200, gin.H{"error": "not implemented"})
}
