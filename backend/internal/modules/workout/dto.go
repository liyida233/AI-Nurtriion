package workout

type ExerciseRequest struct {
	Name           string `json:"name" binding:"required"`
	Category       string `json:"category"`
	MuscleGroup    string `json:"muscleGroup"`
	Equipment      string `json:"equipment"`
	IntensityLevel string `json:"intensityLevel"`
}

type EntryRequest struct {
	ExerciseID string  `json:"exerciseId" binding:"required"`
	Sets       int     `json:"sets" binding:"required,min=1"`
	Reps       int     `json:"reps" binding:"required,min=1"`
	WeightKg   float64 `json:"weightKg" binding:"min=0"`
	RestSec    int     `json:"restSec" binding:"min=0"`
}

type SessionRequest struct {
	WorkoutDate string         `json:"workoutDate" binding:"required"`
	Category    string         `json:"category" binding:"required"`
	DurationMin int            `json:"durationMin" binding:"required,min=1"`
	Notes       string         `json:"notes"`
	Entries     []EntryRequest `json:"entries" binding:"required,min=1"`
}

type Indicators struct {
	TrainingVolume float64        `json:"trainingVolume"`
	EntryCount     int            `json:"entryCount"`
	MuscleGroups   map[string]int `json:"muscleGroups"`
}
