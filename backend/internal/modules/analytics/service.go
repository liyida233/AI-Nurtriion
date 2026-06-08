package analytics

import (
	"context"
	"encoding/json"
	"errors"
	"math"
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
	start, end, normalized, err := ResolveRange(period, "", "")
	if err != nil {
		return DashboardSummary{}, err
	}
	return NewService(db, nil).BuildSummary(ctx, userID, normalized, start, end, false)
}

func (s Service) BuildSummary(ctx context.Context, userID string, period string, start, end time.Time, persist bool) (DashboardSummary, error) {
	days := math.Max(1, math.Ceil(end.Sub(start).Hours()/24))

	summary := DashboardSummary{Period: period, GeneratedAt: time.Now(), StartDate: start, EndDate: end, Days: days}

	sessions, err := s.repo.WorkoutSessions(ctx, userID, start, end)
	if err != nil {
		return summary, err
	}
	summary.WorkoutSessions = int64(len(sessions))
	for _, session := range sessions {
		for _, entry := range session.Entries {
			summary.TrainingVolume += float64(entry.Sets*entry.Reps) * entry.WeightKg
		}
	}
	expectedSessions := math.Max(1, days/7*4)
	summary.WorkoutConsistency = Clamp(float64(summary.WorkoutSessions)/expectedSessions*100, 0, 100)
	summary.ProgressiveOverloadStatus = ProgressiveOverloadStatus(sessions)
	summary.MuscleGroupDistribution = MuscleGroupDistribution(sessions)
	summary.MuscleGroupWarnings = MuscleGroupWarnings(summary.MuscleGroupDistribution)

	meals, err := s.repo.Meals(ctx, userID, start, end)
	if err != nil {
		return summary, err
	}
	summary.MealCount = int64(len(meals))
	mealLogDays := map[string]bool{}
	for _, meal := range meals {
		mealLogDays[meal.MealTime.Format("2006-01-02")] = true
		summary.CaloriesIn += meal.TotalCalories
		summary.Protein += meal.TotalProtein
		summary.Carbohydrates += meal.TotalCarbs
		summary.Fat += meal.TotalFat
	}
	summary.MealLogDays = int64(len(mealLogDays))
	summary.ProteinRatio, summary.CarbohydrateRatio, summary.FatRatio = MacroRatios(summary.Protein, summary.Carbohydrates, summary.Fat)

	if profile, err := s.repo.Profile(ctx, userID); err == nil {
		summary.EstimatedBMR = EstimateBMR(profile)
		summary.EstimatedTDEE = summary.EstimatedBMR * ActivityFactor(profile.ActivityLevel)
		summary.CalorieBalance = summary.CaloriesIn - summary.EstimatedTDEE*days
		summary.CalorieStatus = ClassifyCalorieBalance(summary.CalorieBalance)
	}

	bodyRecords, err := s.repo.BodyRecords(ctx, userID, end)
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
		_ = s.repo.SaveSnapshot(ctx, summary.ToSnapshot(userID, start, end))
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

func ResolveRange(period, startDate, endDate string) (time.Time, time.Time, string, error) {
	now := time.Now()
	end := endOfDay(now)
	switch period {
	case "", "weekly":
		return startOfDay(now.AddDate(0, 0, -6)), end, "weekly", nil
	case "daily":
		return startOfDay(now), end, "daily", nil
	case "monthly":
		return startOfDay(now.AddDate(0, -1, 0)), end, "monthly", nil
	case "custom":
		start, err := parseDate(startDate)
		if err != nil {
			return time.Time{}, time.Time{}, "", errors.New("startDate is required for custom period")
		}
		customEnd, err := parseDate(endDate)
		if err != nil {
			return time.Time{}, time.Time{}, "", errors.New("endDate is required for custom period")
		}
		if customEnd.Before(start) {
			return time.Time{}, time.Time{}, "", errors.New("endDate must be after startDate")
		}
		return startOfDay(start), endOfDay(customEnd), "custom", nil
	default:
		return time.Time{}, time.Time{}, "", errors.New("period must be daily, weekly, monthly, or custom")
	}
}

func parseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func endOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), value.Location())
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
