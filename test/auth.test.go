package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	serverinternal "secure_notes/server/internalpkg"

	"github.com/gin-gonic/gin"
)

// Xác thực: Kiểm tra đăng ký và đăng nhập, bao gồm các trường hợp lỗi.
func TestAuth_RegisterAndLogin_Placeholder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Register handler currently returns 501 NotImplemented
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	serverinternal.Register(c)
	if w.Code != http.StatusNotImplemented {
		t.Fatalf("Register: expected %d, got %d", http.StatusNotImplemented, w.Code)
	}

	// Login handler currently returns 501 NotImplemented
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	serverinternal.Login(c)
	if w.Code != http.StatusNotImplemented {
		t.Fatalf("Login: expected %d, got %d", http.StatusNotImplemented, w.Code)
	}
}
