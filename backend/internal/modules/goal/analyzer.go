package goal

import (
	"errors"
	"time"
)

func FeasibilityWarning(goalType, metric string, target float64, deadline time.Time) error {
	days := time.Until(deadline).Hours() / 24
	if days < 7 {
		return errors.New("goal deadline should allow at least one week for meaningful progress tracking")
	}
	if goalType == "weight_loss" && metric == "weight_kg" && target > 1.5*(days/7) {
		return errors.New("weight loss target appears too aggressive; use a safer weekly target")
	}
	if goalType == "workout_frequency" && target > 7 {
		return errors.New("workout frequency target should not exceed 7 sessions per week")
	}
	return nil
}
