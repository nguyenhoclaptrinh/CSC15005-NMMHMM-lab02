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

	// 4. Auth routes
	r.POST("/api/register", serverpkg.Register)
	r.POST("/api/login", serverpkg.Login)
	r.POST("/api/logout", serverpkg.Logout)

	// 5. Notes routes
	notes := r.Group("/api/notes")
	{
		notes.GET("", serverpkg.ListNotes)
		notes.POST("", serverpkg.UploadNote)
		notes.GET("/:id", serverpkg.GetNote)
		notes.DELETE("/:id", serverpkg.DeleteNote)
		notes.POST("/:id/share", serverpkg.ShareNote)
		notes.GET("/:id/share", serverpkg.ListShares)
		notes.DELETE("/:id/share/:share_id", serverpkg.RevokeShare)
	}

	// 6. Temp URL access (may be anonymous) - not implemented

	// 7. Run server
	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}
