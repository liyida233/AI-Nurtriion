package analytics

import (
	"math"
	"sort"

	"ai-nutrition/backend/internal/models"
)

func EstimateBMR(profile models.UserProfile) float64 {
	if profile.WeightKg <= 0 || profile.HeightCm <= 0 || profile.Age <= 0 {
		return 0
	}
	base := 10*profile.WeightKg + 6.25*profile.HeightCm - 5*float64(profile.Age)
	if profile.Gender == "female" {
		return math.Round(base - 161)
	}
	return math.Round(base + 5)
}

func ActivityFactor(level string) float64 {
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

func ClassifyWeightTrend(records []models.BodyRecord) string {
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

func MovingAverageWeight(records []models.BodyRecord, days int) float64 {
	if len(records) == 0 {
		return 0
	}
	start := len(records) - days
	if start < 0 {
		start = 0
	}
	var total float64
	for _, record := range records[start:] {
		total += record.WeightKg
	}
	return math.Round(total/float64(len(records[start:]))*10) / 10
}

func MacroRatios(protein, carbohydrates, fat float64) (float64, float64, float64) {
	proteinCalories := protein * 4
	carbCalories := carbohydrates * 4
	fatCalories := fat * 9
	total := proteinCalories + carbCalories + fatCalories
	if total == 0 {
		return 0, 0, 0
	}
	return roundPercent(proteinCalories / total * 100), roundPercent(carbCalories / total * 100), roundPercent(fatCalories / total * 100)
}

func NutritionGaps(summary DashboardSummary) []string {
	gaps := []string{}
	if summary.Protein < 420 {
		gaps = append(gaps, "low_weekly_protein")
	}
	if summary.CaloriesIn == 0 {
		gaps = append(gaps, "no_recent_meal_logs")
	}
	if summary.FatRatio > 40 {
		gaps = append(gaps, "high_fat_ratio")
	}
	if summary.CarbohydrateRatio < 30 && summary.CaloriesIn > 0 {
		gaps = append(gaps, "low_carbohydrate_ratio")
	}
	return gaps
}

func MuscleGroupDistribution(sessions []models.WorkoutSession) map[string]int {
	distribution := map[string]int{}
	for _, session := range sessions {
		for _, entry := range session.Entries {
			group := entry.Exercise.MuscleGroup
			if group == "" {
				group = "unknown"
			}
			distribution[group]++
		}
	}
	return distribution
}

func MuscleGroupWarnings(distribution map[string]int) []string {
	if len(distribution) == 0 {
		return []string{"no_recent_workout_logs"}
	}

	warnings := []string{}
	if distribution["lower_body"] == 0 {
		warnings = append(warnings, "lower_body_undertrained")
	}
	if distribution["back"] == 0 {
		warnings = append(warnings, "back_undertrained")
	}

	var maxGroup string
	var maxCount int
	var total int
	for group, count := range distribution {
		total += count
		if count > maxCount {
			maxGroup = group
			maxCount = count
		}
	}
	if total > 0 && float64(maxCount)/float64(total) > 0.6 {
		warnings = append(warnings, "training_overfocused_on_"+maxGroup)
	}
	return warnings
}

func ProgressiveOverloadStatus(sessions []models.WorkoutSession) string {
	if len(sessions) < 2 {
		return "insufficient_data"
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].WorkoutDate.Before(sessions[j].WorkoutDate)
	})

	mid := len(sessions) / 2
	if mid == 0 {
		return "insufficient_data"
	}
	earlyVolume := totalVolume(sessions[:mid])
	recentVolume := totalVolume(sessions[mid:])
	if earlyVolume == 0 && recentVolume > 0 {
		return "improving"
	}
	if earlyVolume == 0 {
		return "insufficient_data"
	}

	change := (recentVolume - earlyVolume) / earlyVolume
	if change > 0.05 {
		return "improving"
	}
	if change < -0.05 {
		return "declining"
	}
	return "stable"
}

func MealQualityScore(summary DashboardSummary) float64 {
	score := 100.0
	if summary.MealCount == 0 {
		return 0
	}
	if summary.Protein < 420 {
		score -= 25
	}
	if summary.FatRatio > 40 {
		score -= 15
	}
	if summary.CarbohydrateRatio < 30 {
		score -= 10
	}
	if summary.CalorieStatus == "surplus" {
		score -= 10
	}
	return Clamp(score, 0, 100)
}

func ClassifyCalorieBalance(balance float64) string {
	if balance < -1000 {
		return "deficit"
	}
	if balance > 1000 {
		return "surplus"
	}
	return "maintenance"
}

func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func roundPercent(value float64) float64 {
	return math.Round(value*10) / 10
}

func totalVolume(sessions []models.WorkoutSession) float64 {
	var volume float64
	for _, session := range sessions {
		for _, entry := range session.Entries {
			volume += float64(entry.Sets*entry.Reps) * entry.WeightKg
		}
	}
	return volume
}
