package models

import "time"

type User struct {
	ID           string      `gorm:"type:char(36);primaryKey" json:"id"`
	Name         string      `gorm:"size:120;not null" json:"name"`
	Email        string      `gorm:"size:180;uniqueIndex;not null" json:"email"`
	PasswordHash string      `gorm:"size:255;not null" json:"-"`
	Role         string      `gorm:"size:40;not null;default:user" json:"role"`
	Profile      UserProfile `json:"profile,omitempty"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

type UserProfile struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID        string    `gorm:"type:char(36);uniqueIndex;not null" json:"userId"`
	Age           int       `json:"age"`
	Gender        string    `gorm:"size:40" json:"gender"`
	HeightCm      float64   `json:"heightCm"`
	WeightKg      float64   `json:"weightKg"`
	ActivityLevel string    `gorm:"size:60" json:"activityLevel"`
	PrimaryGoal   string    `gorm:"size:80" json:"primaryGoal"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Exercise struct {
	ID             string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name           string    `gorm:"size:140;not null" json:"name"`
	Category       string    `gorm:"size:80" json:"category"`
	MuscleGroup    string    `gorm:"size:80" json:"muscleGroup"`
	Equipment      string    `gorm:"size:80" json:"equipment"`
	IntensityLevel string    `gorm:"size:60" json:"intensityLevel"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type WorkoutSession struct {
	ID          string         `gorm:"type:char(36);primaryKey" json:"id"`
	UserID      string         `gorm:"type:char(36);index;not null" json:"userId"`
	WorkoutDate time.Time      `gorm:"index" json:"workoutDate"`
	Category    string         `gorm:"size:80" json:"category"`
	DurationMin int            `json:"durationMin"`
	Notes       string         `gorm:"type:text" json:"notes"`
	Entries     []WorkoutEntry `gorm:"foreignKey:SessionID" json:"entries"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type WorkoutEntry struct {
	ID         string    `gorm:"type:char(36);primaryKey" json:"id"`
	SessionID  string    `gorm:"type:char(36);index;not null" json:"sessionId"`
	ExerciseID string    `gorm:"type:char(36);index;not null" json:"exerciseId"`
	Exercise   Exercise  `json:"exercise,omitempty"`
	Sets       int       `json:"sets"`
	Reps       int       `json:"reps"`
	WeightKg   float64   `json:"weightKg"`
	RestSec    int       `json:"restSec"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type FoodItem struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name          string    `gorm:"size:140;not null" json:"name"`
	ServingSize   string    `gorm:"size:80" json:"servingSize"`
	Calories      float64   `json:"calories"`
	Protein       float64   `json:"protein"`
	Carbohydrates float64   `json:"carbohydrates"`
	Fat           float64   `json:"fat"`
	Sugar         float64   `json:"sugar"`
	Sodium        float64   `json:"sodium"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type MealLog struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID        string    `gorm:"type:char(36);index;not null" json:"userId"`
	FoodItemID    string    `gorm:"type:char(36);index;not null" json:"foodItemId"`
	FoodItem      FoodItem  `json:"foodItem,omitempty"`
	MealType      string    `gorm:"size:60" json:"mealType"`
	Quantity      float64   `json:"quantity"`
	MealTime      time.Time `gorm:"index" json:"mealTime"`
	TotalCalories float64   `json:"totalCalories"`
	TotalProtein  float64   `json:"totalProtein"`
	TotalCarbs    float64   `json:"totalCarbs"`
	TotalFat      float64   `json:"totalFat"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type BodyRecord struct {
	ID         string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID     string    `gorm:"type:char(36);index;not null" json:"userId"`
	RecordDate time.Time `gorm:"index" json:"recordDate"`
	WeightKg   float64   `json:"weightKg"`
	Note       string    `gorm:"type:text" json:"note"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Goal struct {
	ID           string          `gorm:"type:char(36);primaryKey" json:"id"`
	UserID       string          `gorm:"type:char(36);index;not null" json:"userId"`
	GoalType     string          `gorm:"size:80" json:"goalType"`
	TargetMetric string          `gorm:"size:80" json:"targetMetric"`
	TargetValue  float64         `json:"targetValue"`
	Deadline     time.Time       `json:"deadline"`
	Priority     string          `gorm:"size:40" json:"priority"`
	Status       string          `gorm:"size:40;default:active" json:"status"`
	Milestones   []GoalMilestone `gorm:"foreignKey:GoalID" json:"milestones,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

type GoalMilestone struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	GoalID      string    `gorm:"type:char(36);index;not null" json:"goalId"`
	Title       string    `gorm:"size:140" json:"title"`
	TargetValue float64   `json:"targetValue"`
	DueDate     time.Time `json:"dueDate"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type AnalyticsSnapshot struct {
	ID                 string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID             string    `gorm:"type:char(36);index;not null" json:"userId"`
	PeriodType         string    `gorm:"size:40;index" json:"periodType"`
	StartDate          time.Time `gorm:"index" json:"startDate"`
	EndDate            time.Time `gorm:"index" json:"endDate"`
	WorkoutConsistency float64   `json:"workoutConsistency"`
	TrainingVolume     float64   `json:"trainingVolume"`
	CalorieIntake      float64   `json:"calorieIntake"`
	CalorieBalance     float64   `json:"calorieBalance"`
	GoalAdherence      float64   `json:"goalAdherence"`
	WeightTrend        string    `gorm:"size:80" json:"weightTrend"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type AIRecommendation struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID        string    `gorm:"type:char(36);index;not null" json:"userId"`
	Type          string    `gorm:"size:60" json:"type"`
	PromptContext string    `gorm:"type:text" json:"promptContext"`
	Content       string    `gorm:"type:text" json:"content"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type RecommendationFeedback struct {
	ID               string    `gorm:"type:char(36);primaryKey" json:"id"`
	RecommendationID string    `gorm:"type:char(36);index;not null" json:"recommendationId"`
	Rating           string    `gorm:"size:40" json:"rating"`
	Suitability      string    `gorm:"size:80" json:"suitability"`
	Comment          string    `gorm:"type:text" json:"comment"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type ProgressReport struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID      string    `gorm:"type:char(36);index;not null" json:"userId"`
	PeriodType  string    `gorm:"size:40" json:"periodType"`
	GeneratedAt time.Time `json:"generatedAt"`
	Summary     string    `gorm:"type:text" json:"summary"`
	FileURL     string    `gorm:"size:255" json:"fileUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
