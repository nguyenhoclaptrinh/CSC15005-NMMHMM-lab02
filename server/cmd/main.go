package main

import (
	"log"
	"secure-notes-server/config"
	serverpkg "secure-notes-server/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 2. Init SQLite DB
	db, err := serverpkg.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to init DB:", err)
	}
	defer db.Close()

	// 3. Init Gin router
	r := gin.Default()

	// Global middleware: logging, CORS, rate limit
	r.Use(serverpkg.LoggingMiddleware())
	r.Use(serverpkg.CORSMiddleware())
	r.Use(serverpkg.RateLimitMiddleware())

	// 4. Auth routes (register & login are public)
	r.POST("/api/register", serverpkg.Register)
	r.POST("/api/login", serverpkg.Login)
	// Logout requires valid JWT to blacklist token
	r.POST("/api/logout", serverpkg.JWTMiddleware(), serverpkg.Logout)

	// 5. Notes routes - require authentication
	notes := r.Group("/api/notes")
	notes.Use(serverpkg.JWTMiddleware())
	{
		notes.GET("", serverpkg.ListNotes)
		notes.POST("", serverpkg.UploadNote)
		notes.GET("/:id", serverpkg.GetNote)
		notes.DELETE("/:id", serverpkg.DeleteNote)
		notes.POST("/:id/share", serverpkg.ShareNote)
		notes.GET("/:id/share", serverpkg.ListShares)
		notes.DELETE("/:id/share/:share_id", serverpkg.RevokeShare)
	}

	// Share link endpoints
	// Create and revoke share links require auth
	r.POST("/api/share", serverpkg.JWTMiddleware(), serverpkg.CreateShareLink)
	r.DELETE("/api/share/:id", serverpkg.JWTMiddleware(), serverpkg.RevokeShareLink)
	// Public access to share info/content (may be password-protected)
	r.GET("/api/share/:id/info", serverpkg.GetShareInfo)
	r.GET("/api/share/:id", serverpkg.GetSharedContent)

	// 6. Temp URL access (may be anonymous) - not implemented

	// 7. Run server
	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}
