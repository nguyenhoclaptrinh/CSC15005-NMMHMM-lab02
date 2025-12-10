package serverpkg

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware xác thực token JWT từ header Authorization
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}
		tokenStr := parts[1]

		// Phân tích và xác thực JWT bằng hàm hỗ trợ auth
		token, claims, err := ParseJWT(tokenStr)
		if err != nil || token == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Kiểm tra thời hạn (exp) nếu có
		if expRaw, ok := claims["exp"]; ok {
			switch v := expRaw.(type) {
			case float64:
				if int64(v) < time.Now().Unix() {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
					return
				}
			case int64:
				if v < time.Now().Unix() {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
					return
				}
			}
		}

		// Kiểm tra blacklist theo jti nếu tồn tại
		if jtiRaw, ok := claims["jti"]; ok {
			if jti, ok2 := jtiRaw.(string); ok2 && jti != "" {
				if valid, err := ValidateToken(jti); err == nil && !valid {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
					return
				}
			}
		}

		// Đưa user_id và claims vào context để handler dùng
		if uid, ok := claims["user_id"]; ok {
			c.Set("user_id", uid)
		}
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// CORSMiddleware xử lý CORS (chia sẻ tài nguyên giữa nguồn khác nhau)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Bộ giới hạn tần suất đơn giản (in-memory) theo IP
type rateInfo struct {
	mu       sync.Mutex
	requests int
	window   time.Time
}

var (
	rateMap   = make(map[string]*rateInfo)
	rateMapMu sync.Mutex
	// mức mặc định: 100 yêu cầu mỗi phút
	rateLimitRequests = 100
	rateWindow        = time.Minute
)

// RateLimitMiddleware giới hạn yêu cầu theo IP để ngăn lạm dụng
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = "unknown"
		}

		rateMapMu.Lock()
		ri, ok := rateMap[ip]
		if !ok {
			ri = &rateInfo{requests: 0, window: time.Now().Add(rateWindow)}
			rateMap[ip] = ri
		}
		rateMapMu.Unlock()

		ri.mu.Lock()
		now := time.Now()
		if now.After(ri.window) {
			ri.requests = 0
			ri.window = now.Add(rateWindow)
		}
		ri.requests++
		reqs := ri.requests
		retryAfter := int(ri.window.Sub(now).Seconds())
		ri.mu.Unlock()

		if reqs > rateLimitRequests {
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}

// LoggingMiddleware ghi log các yêu cầu và phản hồi HTTP
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		clientIP := c.ClientIP()
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		size := c.Writer.Size()
		log.Printf("%s %s - %s - %d - %dB - %s", method, path, clientIP, status, size, duration)
	}
}
