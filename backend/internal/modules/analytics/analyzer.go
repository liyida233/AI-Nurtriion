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
	sort.Slice(records, func(i, j int) bool {
		return records[i].RecordDate.Before(records[j].RecordDate)
	})
	end := records[len(records)-1].RecordDate
	startDate := end.AddDate(0, 0, -(days - 1))
	var total float64
	var count float64
	for _, record := range records {
		if record.RecordDate.Before(startDate) || record.RecordDate.After(end) {
			continue
		}
		total += record.WeightKg
		count++
	}
	if count == 0 {
		return 0
	}
	return math.Round(total/count*10) / 10
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
	if summary.Protein < proteinTarget(summary) {
		gaps = append(gaps, "low_protein_intake")
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
	if summary.Protein < proteinTarget(summary) {
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

func proteinTarget(summary DashboardSummary) float64 {
	days := summary.Days
	if summary.MealLogDays > 0 {
		days = float64(summary.MealLogDays)
	}
	if days <= 0 {
		days = 7
	}
	return 60 * days
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
			load := entry.WeightKg
			if load <= 0 {
				load = 1
			}
			volume += float64(entry.Sets*entry.Reps) * load
		}
	}
	return volume
}
