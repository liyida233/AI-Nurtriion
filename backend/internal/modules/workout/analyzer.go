package workout

import "ai-nutrition/backend/internal/models"

func Analyze(session models.WorkoutSession) Indicators {
	indicators := Indicators{
		MuscleGroups: map[string]int{},
	}

	for _, entry := range session.Entries {
		indicators.EntryCount++
		indicators.TrainingVolume += float64(entry.Sets*entry.Reps) * entry.WeightKg
		if entry.Exercise.MuscleGroup != "" {
			indicators.MuscleGroups[entry.Exercise.MuscleGroup]++
		}
	}

	return indicators
}
