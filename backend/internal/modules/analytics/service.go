package analytics

import (
	"context"
	"encoding/json"
	"time"

	"ai-nutrition/backend/internal/models"
	"ai-nutrition/backend/internal/modules/shared"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	repo  Repository
	redis *redis.Client
}

func NewService(db *gorm.DB, redisClient *redis.Client) Service {
	return Service{db: db, repo: NewRepository(db), redis: redisClient}
}

func BuildSummary(ctx context.Context, db *gorm.DB, userID string, period string) (DashboardSummary, error) {
	return NewService(db, nil).BuildSummary(ctx, userID, period, false)
}

func (s Service) BuildSummary(ctx context.Context, userID string, period string, persist bool) (DashboardSummary, error) {
	now := time.Now()
	start := now.AddDate(0, 0, -7)
	if period == "monthly" {
		start = now.AddDate(0, -1, 0)
	}
	if period != "weekly" && period != "monthly" {
		period = "weekly"
	}

	summary := DashboardSummary{Period: period, GeneratedAt: now}

	sessions, err := s.repo.WorkoutSessions(ctx, userID, start)
	if err != nil {
		return summary, err
	}
	summary.WorkoutSessions = int64(len(sessions))
	for _, session := range sessions {
		for _, entry := range session.Entries {
			summary.TrainingVolume += float64(entry.Sets*entry.Reps) * entry.WeightKg
		}
	}
	summary.WorkoutConsistency = Clamp(float64(summary.WorkoutSessions)/4*100, 0, 100)
	summary.ProgressiveOverloadStatus = ProgressiveOverloadStatus(sessions)
	summary.MuscleGroupDistribution = MuscleGroupDistribution(sessions)
	summary.MuscleGroupWarnings = MuscleGroupWarnings(summary.MuscleGroupDistribution)

	meals, err := s.repo.Meals(ctx, userID, start)
	if err != nil {
		return summary, err
	}
	summary.MealCount = int64(len(meals))
	for _, meal := range meals {
		summary.CaloriesIn += meal.TotalCalories
		summary.Protein += meal.TotalProtein
		summary.Carbohydrates += meal.TotalCarbs
		summary.Fat += meal.TotalFat
	}
	summary.ProteinRatio, summary.CarbohydrateRatio, summary.FatRatio = MacroRatios(summary.Protein, summary.Carbohydrates, summary.Fat)

	if profile, err := s.repo.Profile(ctx, userID); err == nil {
		summary.EstimatedBMR = EstimateBMR(profile)
		summary.EstimatedTDEE = summary.EstimatedBMR * ActivityFactor(profile.ActivityLevel)
		days := 7.0
		if period == "monthly" {
			days = 30
		}
		summary.CalorieBalance = summary.CaloriesIn - summary.EstimatedTDEE*days
		summary.CalorieStatus = ClassifyCalorieBalance(summary.CalorieBalance)
	}

	bodyRecords, err := s.repo.BodyRecords(ctx, userID)
	if err != nil {
		return summary, err
	}
	if len(bodyRecords) > 0 {
		summary.LatestWeightKg = bodyRecords[len(bodyRecords)-1].WeightKg
		summary.WeightMovingAverage7DayKg = MovingAverageWeight(bodyRecords, 7)
		summary.WeightTrend = ClassifyWeightTrend(bodyRecords)
	}

	activeGoals, err := s.repo.ActiveGoalCount(ctx, userID)
	if err != nil {
		return summary, err
	}
	summary.ActiveGoals = activeGoals

	goals, err := s.repo.Goals(ctx, userID)
	if err != nil {
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

	summary.NutritionGaps = NutritionGaps(summary)
	summary.MealQualityScore = MealQualityScore(summary)

	if persist {
		_ = s.repo.SaveSnapshot(ctx, summary.ToSnapshot(userID, start, now))
	}

	return summary, nil
}

func (s Service) CachedSummary(ctx context.Context, userID, period string) (DashboardSummary, bool) {
	var summary DashboardSummary
	if s.redis == nil || period != "weekly" {
		return summary, false
	}

	value, err := s.redis.Get(ctx, shared.DashboardCacheKey(userID)).Result()
	if err != nil {
		return summary, false
	}
	if err := json.Unmarshal([]byte(value), &summary); err != nil {
		return summary, false
	}
	return summary, true
}

func (s Service) CacheSummary(ctx context.Context, userID string, summary DashboardSummary) {
	if s.redis == nil || summary.Period != "weekly" {
		return
	}
	payload, err := json.Marshal(summary)
	if err != nil {
		return
	}
	_ = s.redis.Set(ctx, shared.DashboardCacheKey(userID), payload, 5*time.Minute).Err()
}

func (s Service) ListSnapshots(ctx context.Context, userID string) ([]models.AnalyticsSnapshot, error) {
	return s.repo.ListSnapshots(ctx, userID)
}

func (s DashboardSummary) ToSnapshot(userID string, start, end time.Time) *models.AnalyticsSnapshot {
	return &models.AnalyticsSnapshot{
		ID:                 uuid.NewString(),
		UserID:             userID,
		PeriodType:         s.Period,
		StartDate:          start,
		EndDate:            end,
		WorkoutConsistency: s.WorkoutConsistency,
		TrainingVolume:     s.TrainingVolume,
		CalorieIntake:      s.CaloriesIn,
		CalorieBalance:     s.CalorieBalance,
		GoalAdherence:      s.GoalAdherence,
		WeightTrend:        s.WeightTrend,
	}
}
