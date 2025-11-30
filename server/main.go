package main

import (
    "secure_notes/internal/server/config"
    "secure_notes/internal/server/handlers"
    "secure_notes/internal/server/storage"
    "github.com/gin-gonic/gin"
    "log"
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
    r.POST("/api/register", handlers.Register)
    r.POST("/api/login", handlers.Login)
    r.POST("/api/refresh-token", handlers.RefreshToken)

    // 5. Notes routes
    notes := r.Group("/api/notes")
    notes.Use(handlers.AuthMiddleware)
    {
        notes.GET("", handlers.ListNotes)
        notes.POST("", handlers.UploadNote)
        notes.GET("/:id", handlers.GetNote)
        notes.DELETE("/:id", handlers.DeleteNote)
        notes.POST("/:id/share", handlers.ShareNote)
        notes.GET("/:id/share", handlers.ListShares)
        notes.DELETE("/:id/share/:share_id", handlers.RevokeShare)
        notes.POST("/:id/temp-url", handlers.CreateTempURL)
    }

    // 6. Temp URL access (may be anonymous)
    r.GET("/api/temp-url/:token", handlers.AccessTempURL)
    r.DELETE("/api/temp-url/:token", handlers.RevokeTempURL)

    // 7. Run server
    log.Println("Server running on port", cfg.Port)
    r.Run(":" + cfg.Port)
}
