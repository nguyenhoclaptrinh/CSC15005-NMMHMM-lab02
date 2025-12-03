package main

import (
	"log"
	"secure_notes/server/config"
	serverinternal "secure_notes/server/internalpkg"
	"secure_notes/server/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 2. Init SQLite DB
	db, err := storage.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to init DB:", err)
	}
	defer db.Close()

	// 3. Init Gin router
	r := gin.Default()

	// 4. Auth routes
	r.POST("/api/register", serverinternal.Register)
	r.POST("/api/login", serverinternal.Login)

	// 5. Notes routes
	notes := r.Group("/api/notes")
	{
		notes.GET("", serverinternal.ListNotes)
		notes.POST("", serverinternal.UploadNote)
		notes.GET("/:id", serverinternal.GetNote)
		notes.DELETE("/:id", serverinternal.DeleteNote)
		notes.POST("/:id/share", serverinternal.ShareNote)
		notes.GET("/:id/share", serverinternal.ListShares)
		notes.DELETE("/:id/share/:share_id", serverinternal.RevokeShare)
	}

	// 6. Temp URL access (may be anonymous) - not implemented

	// 7. Run server
	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}
