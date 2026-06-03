package goal

type GoalRequest struct {
	GoalType     string  `json:"goalType" binding:"required"`
	TargetMetric string  `json:"targetMetric" binding:"required"`
	TargetValue  float64 `json:"targetValue" binding:"required"`
	Deadline     string  `json:"deadline" binding:"required"`
	Priority     string  `json:"priority" binding:"required"`
}

type StatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type MilestoneRequest struct {
	Completed bool `json:"completed"`
}
