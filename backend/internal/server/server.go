package server

import (
	"net/http"

	"ai-nutrition/backend/internal/config"
	"ai-nutrition/backend/internal/middleware"
	"ai-nutrition/backend/internal/modules/analytics"
	"ai-nutrition/backend/internal/modules/auth"
	"ai-nutrition/backend/internal/modules/body"
	"ai-nutrition/backend/internal/modules/goal"
	"ai-nutrition/backend/internal/modules/nutrition"
	"ai-nutrition/backend/internal/modules/profile"
	"ai-nutrition/backend/internal/modules/recommendation"
	"ai-nutrition/backend/internal/modules/workout"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func New(cfg config.Config, db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ai-nutrition-api"})
	})

	auth.RegisterRoutes(api.Group("/auth"), cfg, db, redisClient)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired(cfg))

	profile.RegisterRoutes(protected.Group("/profile"), db, redisClient)
	workout.RegisterRoutes(protected.Group("/workouts"), db, redisClient)
	nutrition.RegisterRoutes(protected.Group("/nutrition"), db, redisClient)
	body.RegisterRoutes(protected.Group("/body-records"), db, redisClient)
	goal.RegisterRoutes(protected.Group("/goals"), db, redisClient)
	analytics.RegisterRoutes(protected.Group("/dashboard"), db, redisClient)
	recommendation.RegisterRoutes(protected.Group("/recommendations"), cfg, db, redisClient)

	return router
}
