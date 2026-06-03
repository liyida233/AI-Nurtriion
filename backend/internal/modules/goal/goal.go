package goal

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

type createGoalRequest struct {
	GoalType     string  `json:"goalType" binding:"required"`
	TargetMetric string  `json:"targetMetric" binding:"required"`
	TargetValue  float64 `json:"targetValue" binding:"required"`
	Deadline     string  `json:"deadline" binding:"required"`
	Priority     string  `json:"priority" binding:"required"`
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	handler := Handler{db: db, redis: redisClient}
	router.GET("", handler.List)
	router.POST("", handler.Create)
}

func (h Handler) List(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var goals []models.Goal
	if err := h.db.Preload("Milestones").Where("user_id = ?", userID).Order("created_at desc").Find(&goals).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not list goals")
		return
	}
	httpctx.OK(c, goals)
}

func (h Handler) Create(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var req createGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	deadline, err := time.Parse("2006-01-02", req.Deadline)
	if err != nil {
		httpctx.Error(c, http.StatusBadRequest, "deadline must be YYYY-MM-DD")
		return
	}

	if warning := feasibilityWarning(req.GoalType, req.TargetMetric, req.TargetValue, deadline); warning != "" {
		httpctx.Error(c, http.StatusUnprocessableEntity, warning)
		return
	}

	goal := models.Goal{
		ID:           uuid.NewString(),
		UserID:       userID,
		GoalType:     req.GoalType,
		TargetMetric: req.TargetMetric,
		TargetValue:  req.TargetValue,
		Deadline:     deadline,
		Priority:     req.Priority,
		Status:       "active",
	}
	goal.Milestones = buildMilestones(goal)

	if err := h.db.Create(&goal).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not save goal")
		return
	}

	shared.InvalidateDashboard(c.Request.Context(), h.redis, userID)
	httpctx.Created(c, goal)
}

func feasibilityWarning(goalType, metric string, target float64, deadline time.Time) string {
	days := time.Until(deadline).Hours() / 24
	if days < 7 {
		return "goal deadline should allow at least one week for meaningful progress tracking"
	}
	if goalType == "weight_loss" && metric == "weight_kg" && target > 1.5*(days/7) {
		return "weight loss target appears too aggressive; use a safer weekly target"
	}
	if goalType == "workout_frequency" && target > 7 {
		return "workout frequency target should not exceed 7 sessions per week"
	}
	return ""
}

func buildMilestones(goal models.Goal) []models.GoalMilestone {
	now := time.Now()
	totalDays := goal.Deadline.Sub(now).Hours() / 24
	if totalDays <= 0 {
		return nil
	}
	milestoneCount := 3
	milestones := make([]models.GoalMilestone, 0, milestoneCount)
	for i := 1; i <= milestoneCount; i++ {
		ratio := float64(i) / float64(milestoneCount)
		milestones = append(milestones, models.GoalMilestone{
			ID:          uuid.NewString(),
			Title:       "Milestone " + string(rune('0'+i)),
			TargetValue: goal.TargetValue * ratio,
			DueDate:     now.Add(time.Duration(totalDays*ratio) * 24 * time.Hour),
			Completed:   false,
		})
	}
	return milestones
}
