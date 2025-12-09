package serverpkg

import (
	"github.com/gin-gonic/gin"
)

// JWTMiddleware validates JWT tokens from Authorization header
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Extract JWT from Authorization: Bearer <token>
		// TODO: Validate token signature using secret key
		// TODO: Check if token is in blacklist (revoked)
		// TODO: Parse claims (user_id, exp, jti)
		// TODO: Set user_id to context: c.Set("user_id", userID)
		// TODO: Return 401 Unauthorized if invalid
		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Set CORS headers:
		// - Access-Control-Allow-Origin
		// - Access-Control-Allow-Methods (GET, POST, PUT, DELETE)
		// - Access-Control-Allow-Headers (Authorization, Content-Type)
		// - Access-Control-Allow-Credentials
		// TODO: Handle preflight OPTIONS requests
		c.Next()
	}
}

// RateLimitMiddleware limits requests per IP to prevent abuse
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Get client IP: c.ClientIP()
		// TODO: Check rate limit (e.g., 100 req/min per IP)
		// TODO: Use in-memory cache (sync.Map) or Redis for counters
		// TODO: Return 429 Too Many Requests if exceeded
		// TODO: Add Retry-After header
		c.Next()
	}
}

// LoggingMiddleware logs all HTTP requests and responses
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Log request: method, path, IP, user-agent
		// TODO: Record start time
		// TODO: Call c.Next()
		// TODO: Calculate response time
		// TODO: Log response: status code, duration, size
		c.Next()
	}
}
