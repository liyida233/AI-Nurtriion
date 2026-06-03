package workout

import (
	"net/http"
	"time"

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

type workoutEntryRequest struct {
	ExerciseID string  `json:"exerciseId" binding:"required"`
	Sets       int     `json:"sets" binding:"required,min=1"`
	Reps       int     `json:"reps" binding:"required,min=1"`
	WeightKg   float64 `json:"weightKg" binding:"min=0"`
	RestSec    int     `json:"restSec" binding:"min=0"`
}

type createWorkoutRequest struct {
	WorkoutDate string                `json:"workoutDate" binding:"required"`
	Category    string                `json:"category" binding:"required"`
	DurationMin int                   `json:"durationMin" binding:"required,min=1"`
	Notes       string                `json:"notes"`
	Entries     []workoutEntryRequest `json:"entries" binding:"required,min=1"`
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	handler := Handler{db: db, redis: redisClient}
	router.GET("/exercises", handler.ListExercises)
	router.POST("/exercises", handler.CreateExercise)
	router.GET("", handler.List)
	router.POST("", handler.Create)
}

func (h Handler) ListExercises(c *gin.Context) {
	var exercises []models.Exercise
	if err := h.db.Order("name asc").Find(&exercises).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not list exercises")
		return
	}
	httpctx.OK(c, exercises)
}

func (h Handler) CreateExercise(c *gin.Context) {
	var exercise models.Exercise
	if err := c.ShouldBindJSON(&exercise); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	exercise.ID = uuid.NewString()
	if err := h.db.Create(&exercise).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not create exercise")
		return
	}
	httpctx.Created(c, exercise)
}

func (h Handler) List(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var sessions []models.WorkoutSession
	if err := h.db.Preload("Entries.Exercise").Where("user_id = ?", userID).Order("workout_date desc").Limit(50).Find(&sessions).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not list workouts")
		return
	}
	httpctx.OK(c, sessions)
}

func (h Handler) Create(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var req createWorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	workoutDate, err := time.Parse("2006-01-02", req.WorkoutDate)
	if err != nil {
		httpctx.Error(c, http.StatusBadRequest, "workoutDate must be YYYY-MM-DD")
		return
	}

	session := models.WorkoutSession{
		ID:          uuid.NewString(),
		UserID:      userID,
		WorkoutDate: workoutDate,
		Category:    req.Category,
		DurationMin: req.DurationMin,
		Notes:       req.Notes,
	}
	for _, item := range req.Entries {
		session.Entries = append(session.Entries, models.WorkoutEntry{
			ID:         uuid.NewString(),
			ExerciseID: item.ExerciseID,
			Sets:       item.Sets,
			Reps:       item.Reps,
			WeightKg:   item.WeightKg,
			RestSec:    item.RestSec,
		})
	}

	if err := h.db.Create(&session).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not save workout")
		return
	}

	shared.InvalidateDashboard(c.Request.Context(), h.redis, userID)
	httpctx.Created(c, gin.H{"workout": session, "indicators": workoutIndicators(session)})
}

func workoutIndicators(session models.WorkoutSession) gin.H {
	var volume float64
	for _, entry := range session.Entries {
		volume += float64(entry.Sets*entry.Reps) * entry.WeightKg
	}
	return gin.H{"trainingVolume": volume, "entryCount": len(session.Entries)}
}
