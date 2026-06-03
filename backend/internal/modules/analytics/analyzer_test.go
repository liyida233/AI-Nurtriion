package analytics

import (
	"testing"
	"time"

	"ai-nutrition/backend/internal/models"
)

func TestEstimateBMR(t *testing.T) {
	profile := models.UserProfile{
		Age:      22,
		Gender:   "male",
		HeightCm: 175,
		WeightKg: 70,
	}

	got := EstimateBMR(profile)
	want := 1689.0
	if got != want {
		t.Fatalf("EstimateBMR() = %v, want %v", got, want)
	}
}

func TestMacroRatios(t *testing.T) {
	protein, carbs, fat := MacroRatios(100, 200, 50)

	if protein != 24.2 || carbs != 48.5 || fat != 27.3 {
		t.Fatalf("MacroRatios() = %.1f %.1f %.1f", protein, carbs, fat)
	}
}

func TestProgressiveOverloadStatus(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	sessions := []models.WorkoutSession{
		sessionWithVolume(base, 3, 10, 30),
		sessionWithVolume(base.AddDate(0, 0, 3), 3, 10, 32),
		sessionWithVolume(base.AddDate(0, 0, 6), 4, 10, 35),
		sessionWithVolume(base.AddDate(0, 0, 9), 4, 10, 36),
	}

	got := ProgressiveOverloadStatus(sessions)
	if got != "improving" {
		t.Fatalf("ProgressiveOverloadStatus() = %s, want improving", got)
	}
}

func TestMuscleGroupWarnings(t *testing.T) {
	distribution := map[string]int{
		"chest": 4,
		"back":  1,
	}

	warnings := MuscleGroupWarnings(distribution)
	if !contains(warnings, "lower_body_undertrained") {
		t.Fatalf("expected lower_body_undertrained warning, got %v", warnings)
	}
	if !contains(warnings, "training_overfocused_on_chest") {
		t.Fatalf("expected chest overfocus warning, got %v", warnings)
	}
}

func TestMealQualityScore(t *testing.T) {
	summary := DashboardSummary{
		MealCount:         10,
		Protein:           300,
		FatRatio:          45,
		CarbohydrateRatio: 25,
		CalorieStatus:     "surplus",
	}

	got := MealQualityScore(summary)
	want := 40.0
	if got != want {
		t.Fatalf("MealQualityScore() = %v, want %v", got, want)
	}
}

func sessionWithVolume(date time.Time, sets, reps int, weight float64) models.WorkoutSession {
	return models.WorkoutSession{
		WorkoutDate: date,
		Entries: []models.WorkoutEntry{
			{Sets: sets, Reps: reps, WeightKg: weight},
		},
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
