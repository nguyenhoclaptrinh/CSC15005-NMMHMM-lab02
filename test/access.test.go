package test

import (
    "testing"

    "net/http"
    "net/http/httptest"
    "time"

    "github.com/gin-gonic/gin"
)

// Giới hạn truy cập: Đảm bảo các liên kết hết hạn không thể truy cập.
func TestAccess_LinkExpiry_Placeholder(t *testing.T) {
    t.Skip("access/link-expiry tests are placeholders until link sharing is implemented")

    // Example future flow:
    // - Create a share with expiry = now + 1s
    // - Attempt access immediately -> allowed
    // - Sleep >1s -> access should be denied
    _ = gin.New
    _ = http.StatusOK
    _ = httptest.NewRecorder
    _ = time.Sleep
}
