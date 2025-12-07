package serverpkg

import (
	"context"
	"net/http"
    "encoding/base64"
	"encoding/hex"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
    "fmt"
    "strings"
	"time"
	"regexp"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)


type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterResponse struct {
	UserID string `json:"user_id"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	Token   string `json:"token"`
	KdfSalt string `json:"kdf_salt"`
}
type LogoutResponse struct {
	Message string `json:"message"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type UserClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	db               *sql.DB
	jwtSecretKey     = []byte("")
	argonPepper      = ""
	accessTokenExpiry = time.Hour * 24

	argonTime    = uint32(1)
	argonMemory  = uint32(64 * 1024) 
	argonThreads = uint8(4)
	argonKeyLen  = uint32(32)
)
func InitAuth(database *sql.DB, secretKey, pepper string) error {
	if database == nil {
		return fmt.Errorf("database connection is required")
	}
	if len(secretKey) < 32 {
		return fmt.Errorf("JWT secret key must be at least 32 characters")
	}
	db = database
	jwtSecretKey = []byte(secretKey)
	argonPepper = pepper
	return nil
}
func validatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters"
	}
	
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[@#$%^&+=!]`).MatchString(password)
	
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return false, "Password must contain uppercase, lowercase, numbers and special characters (@#$%^&+=!)"
	}
	
	return true, ""
}
func validateUsername(username string) (bool, string) {
	if len(username) < 3 || len(username) > 50 {
		return false, "Username must be between 3 and 50 characters"
	}
	
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return false, "Username can only contain letters, numbers and underscores"
	}
	
	return true, ""
}

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %v", err)
	}
	return salt, nil
}

func serverHashPassword(password string, salt []byte) (string, error) {
	pepperPassword := password + argonPepper
	hash := argon2.IDKey([]byte(pepperPassword), salt, argonTime, argonMemory, argonThreads, argonKeyLen,)
	combined := make([]byte, len(salt)+len(hash))
	copy(combined, salt)
	copy(combined[len(salt):], hash)
	return base64.StdEncoding.EncodeToString(combined), nil
}

func verifyPassword(password, storedHash string, salt []byte) (bool, error) {
	hashedInput, err := serverHashPassword(password, salt)
	if err != nil {
		return false, err
	}
	return hashedInput == storedHash, nil
}
func generateToken(userID, username string) (string, string, error) {
	jti := uuid.New().String()
	expiryTime := time.Now().Add(accessTokenExpiry)
	
	claims := &UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign token: %v", err)
	}
	return tokenString, jti, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func isTokenBlacklisted(jti string) (bool, error) {
	ctx := context.Background()
	query := `SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE jti = $1 AND expires_at > NOW())`
	var exists bool
	err := db.QueryRow(ctx, query, jti).Scan(&exists)
	return exists, err
}

func addToBlacklist(jti string, expiresAt time.Time) error {
	ctx := context.Background()
	query := `INSERT INTO token_blacklist (jti, expires_at) VALUES ($1, $2) 
	          ON CONFLICT (jti) DO UPDATE SET expires_at = EXCLUDED.expires_at`
	_, err := db.Exec(ctx, query, jti, expiresAt)
	return err
}

func Register(c *gin.Context) {
	var req RegisterRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	if valid, msg := validateUsername(req.Username); !valid {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: msg})
		return
	}
	
	if valid, msg := validatePassword(req.Password); !valid {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: msg})
		return
	}
	
	var exists bool
	err := db.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", 
		req.Username).Scan(&exists)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Database error",
		})
		return
	}
	
	if exists {
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: "Username already exists",
		})
		return
	}
	
	kdfSalt, err := generateSalt(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate salt",
		})
		return
	}
	
	passwordHash, err := serverHashPassword(req.Password, kdfSalt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to hash password",
		})
		return
	}
	
	userID := uuid.New()
	
	// Lưu vào database
	_, err = db.Exec(context.Background(),
		`INSERT INTO users (id, username, password_hash, kdf_salt) 
		 VALUES ($1, $2, $3, $4)`,
		userID, req.Username, passwordHash, kdfSalt)
	
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || 
		   strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Username already exists",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create user",
		})
		return
	}
	
	// 201 Created
	c.JSON(http.StatusCreated, RegisterResponse{
		UserID: userID.String(),
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request",
		})
		return
	}
	
	var userID uuid.UUID
	var username string
	var passwordHash string
	var kdfSaltBytes []byte
	
	err := db.QueryRow(context.Background(),
		`SELECT id, username, password_hash, kdf_salt 
		 FROM users WHERE username = $1`, req.Username).Scan(
		&userID, &username, &passwordHash, &kdfSaltBytes)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid credentials",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Database error",
		})
		return
	}
	

	valid, err := verifyPassword(req.Password, passwordHash, kdfSaltBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Authentication error",
		})
		return
	}
	
	if !valid {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid credentials",
		})
		return
	}
	
	// Tạo JWT token
	token, jti, err := generateToken(userID.String(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		KdfSalt: base64.StdEncoding.EncodeToString(kdfSaltBytes),
	})
}
func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Authorization header required",
		})
		return
	}
	
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid authorization format",
		})
		return
	}
	
	tokenString := parts[1]
	
	// Parse token để lấy JTI
	token, claims, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		c.JSON(http.StatusOK, LogoutResponse{
			Message: "Logged out successfully",
		})
		return
	}
	
	// Thêm vào blacklist nếu parse được
	var jti string
	if claimsMap, ok := claims.(jwt.MapClaims); ok {
		if jtiVal, exists := claimsMap["jti"].(string); exists {
			jti = jtiVal
		}
	}
	
	// Nếu có jti, thêm vào blacklist
	if jti != "" {
		var expiryTime time.Time
		if claimsMap, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claimsMap["exp"].(float64); ok {
				expiryTime = time.Unix(int64(exp), 0)
			} else {
				expiryTime = time.Now().Add(accessTokenExpiry)
			}
		} else {
			expiryTime = time.Now().Add(accessTokenExpiry)
		}
		
		_ = addToBlacklist(jti, expiryTime)
	}
	
	// 200 OK
	c.JSON(http.StatusOK, LogoutResponse{
		Message: "Logged out successfully",
	})
}


func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Authorization header required",
			})
			c.Abort()
			return
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid authorization format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}
		
		tokenString := parts[1]
		
		// Parse và validate token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})
		
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}
		
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid token",
			})
			c.Abort()
			return
		}
		
		// Kiểm tra blacklist trong token_blacklist
		if claims, ok := token.Claims.(*UserClaims); ok {
			blacklisted, err := isTokenBlacklisted(claims.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error: "Failed to check token status",
				})
				c.Abort()
				return
			}
			if blacklisted {
				c.JSON(http.StatusUnauthorized, ErrorResponse{
					Error: "Token has been revoked",
				})
				c.Abort()
				return
			}
			
			// Set user info vào context
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("jti", claims.ID)
		}
		
		c.Next()
	}
}