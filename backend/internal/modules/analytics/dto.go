package analytics

import "time"

type DashboardSummary struct {
	Period                    string         `json:"period"`
	GeneratedAt               time.Time      `json:"generatedAt"`
	StartDate                 time.Time      `json:"startDate"`
	EndDate                   time.Time      `json:"endDate"`
	Days                      float64        `json:"days"`
	WorkoutSessions           int64          `json:"workoutSessions"`
	TrainingVolume            float64        `json:"trainingVolume"`
	WorkoutConsistency        float64        `json:"workoutConsistency"`
	ProgressiveOverloadStatus string         `json:"progressiveOverloadStatus"`
	MuscleGroupDistribution   map[string]int `json:"muscleGroupDistribution"`
	MuscleGroupWarnings       []string       `json:"muscleGroupWarnings"`
	CaloriesIn                float64        `json:"caloriesIn"`
	Protein                   float64        `json:"protein"`
	Carbohydrates             float64        `json:"carbohydrates"`
	Fat                       float64        `json:"fat"`
	MealCount                 int64          `json:"mealCount"`
	MealLogDays               int64          `json:"mealLogDays"`
	MealQualityScore          float64        `json:"mealQualityScore"`
	ProteinRatio              float64        `json:"proteinRatio"`
	CarbohydrateRatio         float64        `json:"carbohydrateRatio"`
	FatRatio                  float64        `json:"fatRatio"`
	NutritionGaps             []string       `json:"nutritionGaps"`
	EstimatedBMR              float64        `json:"estimatedBmr"`
	EstimatedTDEE             float64        `json:"estimatedTdee"`
	CalorieBalance            float64        `json:"calorieBalance"`
	CalorieStatus             string         `json:"calorieStatus"`
	LatestWeightKg            float64        `json:"latestWeightKg"`
	WeightMovingAverage7DayKg float64        `json:"weightMovingAverage7DayKg"`
	WeightTrend               string         `json:"weightTrend"`
	ActiveGoals               int64          `json:"activeGoals"`
	CompletedMilestones       int64          `json:"completedMilestones"`
	TotalMilestones           int64          `json:"totalMilestones"`
	GoalAdherence             float64        `json:"goalAdherence"`
}
