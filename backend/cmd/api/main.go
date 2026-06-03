package main

import (
	"log"

	"ai-nutrition/backend/internal/cache"
	"ai-nutrition/backend/internal/config"
	"ai-nutrition/backend/internal/database"
	"ai-nutrition/backend/internal/server"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatalf("database seed failed: %v", err)
	}

	redisClient := cache.Connect(cfg)
	app := server.New(cfg, db, redisClient)

	if err := app.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
