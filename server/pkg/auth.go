package serverpkg

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ============================================================
// JWT MANAGEMENT
// ============================================================

// JWTSecretKey is the secret key for signing JWT tokens
// TODO: Move to environment variable in production!
var JWTSecretKey = []byte("your-secret-key-change-this-in-production")

// GenerateJWT creates access token and refresh token for authenticated user
func GenerateJWT(userID string, username string) (accessToken string, refreshToken string, err error) {
	// Create access token (JWT) with 15 minutes expiry
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
		"jti":      uuid.New().String(), // JWT ID (unique identifier)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString(JWTSecretKey)
	if err != nil {
		return "", "", err
	}

	// Create refresh token (random 32 bytes)
	refreshTokenBytes := make([]byte, 32)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		return "", "", err
	}
	refreshToken = base64.URLEncoding.EncodeToString(refreshTokenBytes)

	// TODO: Store refresh token in database with 7 days expiry
	// INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES (?, SHA256(?), datetime('now', '+7 days'))

	return accessToken, refreshToken, nil
}

// ParseJWT validates and parses a JWT token
func ParseJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, jwt.ErrTokenInvalidClaims
}

// RefreshToken generates new access token from valid refresh token
func RefreshToken(oldRefreshToken string) (newAccessToken string, err error) {
	// TODO: Query DB: SELECT user_id, username, expires_at FROM refresh_tokens WHERE token_hash = SHA256(?)
	// TODO: Check if refresh token is expired
	// TODO: Generate new access token
	return "", nil
}

// BlacklistToken adds access token to blacklist on logout
func BlacklistToken(jti string, expiresAt time.Time) error {
	// TODO: INSERT INTO token_blacklist (jti, expires_at)
	return nil
}

// ValidateToken checks if token is valid and not blacklisted
func ValidateToken(jti string) (valid bool, err error) {
	// TODO: Query DB: SELECT 1 FROM token_blacklist WHERE jti = ?
	return true, nil
}

// RevokeRefreshToken removes refresh token from database
func RevokeRefreshToken(token string) error {
	// TODO: DELETE FROM refresh_tokens WHERE token_hash = SHA256(?)
	return nil
}

// ============================================================
// PASSWORD HASHING (Argon2id)
// ============================================================

// HashPassword uses Argon2id to hash password with salt
func HashPassword(password string, salt []byte) (hash string, err error) {
	// TODO: Use Argon2id:
	//   - Time cost: 1, Memory: 64MB, Parallelism: 4, Key length: 32 bytes
	//   hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	//   return base64.StdEncoding.EncodeToString(hash), nil
	return "", nil
}

// VerifyPassword compares provided password with stored hash
func VerifyPassword(password string, storedHash string, salt []byte) (valid bool, err error) {
	// TODO: Hash password and compare with constant-time comparison
	return false, nil
}

// GenerateSalt creates a random 16-byte salt
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return salt, err
}

// EncodeSalt converts salt bytes to base64
func EncodeSalt(salt []byte) string {
	return base64.StdEncoding.EncodeToString(salt)
}

// DecodeSalt converts base64 salt string back to bytes
func DecodeSalt(saltStr string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(saltStr)
}

// ============================================================
// AUTH HANDLERS
// ============================================================

// Register handles user registration
// POST /api/auth/register
func Register(c *gin.Context) {
	// TODO: Parse request body: username, password
	// TODO: Validate input (length, complexity)
	// TODO: Check if username exists
	// TODO: Generate salt
	// TODO: Hash password with Argon2id
	// TODO: INSERT INTO users (id, username, password_hash, kdf_salt)
	// TODO: Return 201 Created
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// Login handles user authentication
// POST /api/auth/login
func Login(c *gin.Context) {
	// TODO: Parse request body: username, password_hash (already hashed on client)
	// TODO: Query DB: SELECT id, password_hash, kdf_salt FROM users WHERE username = ?
	// TODO: Verify password hash
	// TODO: Generate JWT tokens
	// TODO: Return access_token, refresh_token
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// Logout handles user logout
// POST /api/auth/logout
func Logout(c *gin.Context) {
	// TODO: Get JWT from Authorization header
	// TODO: Parse JWT to get jti
	// TODO: Blacklist token
	// TODO: Revoke refresh token
	// TODO: Return 200 OK
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// GetSalt returns the KDF salt for a user (used during login)
// GET /api/auth/salt?username=alice
func GetSalt(c *gin.Context) {
	// TODO: Get username from query params
	// TODO: Query DB: SELECT kdf_salt FROM users WHERE username = ?
	// TODO: Return salt in JSON
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
