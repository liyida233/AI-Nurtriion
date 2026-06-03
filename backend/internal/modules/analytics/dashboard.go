package analytics

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"ai-nutrition/backend/internal/httpctx"
	"ai-nutrition/backend/internal/models"
	"ai-nutrition/backend/internal/modules/shared"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	redis *redis.Client
}

type DashboardSummary struct {
	Period              string  `json:"period"`
	WorkoutSessions     int64   `json:"workoutSessions"`
	TrainingVolume      float64 `json:"trainingVolume"`
	WorkoutConsistency  float64 `json:"workoutConsistency"`
	CaloriesIn          float64 `json:"caloriesIn"`
	Protein             float64 `json:"protein"`
	Carbohydrates       float64 `json:"carbohydrates"`
	Fat                 float64 `json:"fat"`
	EstimatedBMR        float64 `json:"estimatedBmr"`
	EstimatedTDEE       float64 `json:"estimatedTdee"`
	CalorieBalance      float64 `json:"calorieBalance"`
	LatestWeightKg      float64 `json:"latestWeightKg"`
	WeightTrend         string  `json:"weightTrend"`
	ActiveGoals         int64   `json:"activeGoals"`
	CompletedMilestones int64   `json:"completedMilestones"`
	TotalMilestones     int64   `json:"totalMilestones"`
	GoalAdherence       float64 `json:"goalAdherence"`
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	handler := Handler{db: db, redis: redisClient}
	router.GET("/summary", handler.Summary)
}

func (h Handler) Summary(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	if cached, ok := h.readCached(c.Request.Context(), userID); ok {
		httpctx.OK(c, cached)
		return
	}

	summary, err := BuildSummary(c.Request.Context(), h.db, userID, "weekly")
	if err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not build dashboard summary")
		return
	}
	h.writeCached(c.Request.Context(), userID, summary)

	httpctx.OK(c, summary)
}

func BuildSummary(ctx context.Context, db *gorm.DB, userID string, period string) (DashboardSummary, error) {
	now := time.Now()
	start := now.AddDate(0, 0, -7)
	if period == "monthly" {
		start = now.AddDate(0, -1, 0)
	}

	var summary DashboardSummary
	summary.Period = period

	var sessions []models.WorkoutSession
	if err := db.WithContext(ctx).Preload("Entries").Where("user_id = ? AND workout_date >= ?", userID, start).Find(&sessions).Error; err != nil {
		return summary, err
	}
	summary.WorkoutSessions = int64(len(sessions))
	for _, session := range sessions {
		for _, entry := range session.Entries {
			summary.TrainingVolume += float64(entry.Sets*entry.Reps) * entry.WeightKg
		}
	}
	summary.WorkoutConsistency = clamp(float64(summary.WorkoutSessions)/4*100, 0, 100)

	var meals []models.MealLog
	if err := db.WithContext(ctx).Where("user_id = ? AND meal_time >= ?", userID, start).Find(&meals).Error; err != nil {
		return summary, err
	}
	for _, meal := range meals {
		summary.CaloriesIn += meal.TotalCalories
		summary.Protein += meal.TotalProtein
		summary.Carbohydrates += meal.TotalCarbs
		summary.Fat += meal.TotalFat
	}

	var profile models.UserProfile
	if err := db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err == nil {
		summary.EstimatedBMR = estimateBMR(profile)
		summary.EstimatedTDEE = summary.EstimatedBMR * activityFactor(profile.ActivityLevel)
		summary.CalorieBalance = summary.CaloriesIn - summary.EstimatedTDEE*7
	}

	var bodyRecords []models.BodyRecord
	if err := db.WithContext(ctx).Where("user_id = ?", userID).Order("record_date asc").Limit(30).Find(&bodyRecords).Error; err != nil {
		return summary, err
	}
	if len(bodyRecords) > 0 {
		summary.LatestWeightKg = bodyRecords[len(bodyRecords)-1].WeightKg
		summary.WeightTrend = classifyWeightTrend(bodyRecords)
	}

	if err := db.WithContext(ctx).Model(&models.Goal{}).Where("user_id = ? AND status = ?", userID, "active").Count(&summary.ActiveGoals).Error; err != nil {
		return summary, err
	}
	var goals []models.Goal
	if err := db.WithContext(ctx).Preload("Milestones").Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		return summary, err
	}
	for _, goal := range goals {
		for _, milestone := range goal.Milestones {
			summary.TotalMilestones++
			if milestone.Completed {
				summary.CompletedMilestones++
			}
		}
	}
	if summary.TotalMilestones > 0 {
		summary.GoalAdherence = float64(summary.CompletedMilestones) / float64(summary.TotalMilestones) * 100
	}

	snapshot := models.AnalyticsSnapshot{
		ID:                 "", // Filled only when explicitly persisted in future iterations.
		UserID:             userID,
		PeriodType:         period,
		StartDate:          start,
		EndDate:            now,
		WorkoutConsistency: summary.WorkoutConsistency,
		TrainingVolume:     summary.TrainingVolume,
		CalorieIntake:      summary.CaloriesIn,
		CalorieBalance:     summary.CalorieBalance,
		GoalAdherence:      summary.GoalAdherence,
		WeightTrend:        summary.WeightTrend,
	}
	_ = snapshot

	return summary, nil
}

func (h Handler) readCached(ctx context.Context, userID string) (DashboardSummary, bool) {
	var summary DashboardSummary
	if h.redis == nil {
		return summary, false
	}

	value, err := h.redis.Get(ctx, shared.DashboardCacheKey(userID)).Result()
	if err != nil {
		return summary, false
	}
	if err := json.Unmarshal([]byte(value), &summary); err != nil {
		return summary, false
	}
	return summary, true
}

func (h Handler) writeCached(ctx context.Context, userID string, summary DashboardSummary) {
	if h.redis == nil {
		return
	}
	payload, err := json.Marshal(summary)
	if err != nil {
		return
	}
	_ = h.redis.Set(ctx, shared.DashboardCacheKey(userID), payload, 5*time.Minute).Err()
}

func estimateBMR(profile models.UserProfile) float64 {
	if profile.WeightKg <= 0 || profile.HeightCm <= 0 || profile.Age <= 0 {
		return 0
	}
	base := 10*profile.WeightKg + 6.25*profile.HeightCm - 5*float64(profile.Age)
	if profile.Gender == "female" {
		return math.Round(base - 161)
	}
	return math.Round(base + 5)
}

func activityFactor(level string) float64 {
	switch level {
	case "light":
		return 1.375
	case "moderate":
		return 1.55
	case "active":
		return 1.725
	default:
		return 1.2
	}
}

func classifyWeightTrend(records []models.BodyRecord) string {
	if len(records) < 2 {
		return "insufficient_data"
	}
	first := records[0].WeightKg
	last := records[len(records)-1].WeightKg
	delta := last - first
	if math.Abs(delta) < 0.3 {
		return "stable"
	}
	if delta > 0 {
		return "increasing"
	}
	return "decreasing"
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
