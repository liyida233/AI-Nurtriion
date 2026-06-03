package profile

import (
	"net/http"

	"ai-nutrition/backend/internal/httpctx"
	"ai-nutrition/backend/internal/models"
	"ai-nutrition/backend/internal/modules/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	redis *redis.Client
}

type upsertProfileRequest struct {
	Age           int     `json:"age" binding:"required,min=13,max=100"`
	Gender        string  `json:"gender" binding:"required"`
	HeightCm      float64 `json:"heightCm" binding:"required,min=80,max=250"`
	WeightKg      float64 `json:"weightKg" binding:"required,min=25,max=350"`
	ActivityLevel string  `json:"activityLevel" binding:"required"`
	PrimaryGoal   string  `json:"primaryGoal" binding:"required"`
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	handler := Handler{db: db, redis: redisClient}
	router.GET("", handler.Get)
	router.PUT("", handler.Upsert)
}

func (h Handler) Get(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var profile models.UserProfile
	if err := h.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		httpctx.Error(c, http.StatusNotFound, "profile not found")
		return
	}

	httpctx.OK(c, profile)
}

func (h Handler) Upsert(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var req upsertProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var profile models.UserProfile
	err = h.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		httpctx.Error(c, http.StatusInternalServerError, "could not load profile")
		return
	}
	if err == gorm.ErrRecordNotFound {
		profile.ID = uuid.NewString()
		profile.UserID = userID
	}

	profile.Age = req.Age
	profile.Gender = req.Gender
	profile.HeightCm = req.HeightCm
	profile.WeightKg = req.WeightKg
	profile.ActivityLevel = req.ActivityLevel
	profile.PrimaryGoal = req.PrimaryGoal

	if err := h.db.Save(&profile).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not save profile")
		return
	}

	shared.InvalidateDashboard(c.Request.Context(), h.redis, userID)
	httpctx.OK(c, profile)
}
