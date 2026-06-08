package goal

type GoalRequest struct {
	GoalType     string           `json:"goalType" binding:"required"`
	TargetMetric string           `json:"targetMetric" binding:"required"`
	TargetValue  float64          `json:"targetValue" binding:"required"`
	Deadline     string           `json:"deadline" binding:"required"`
	Priority     string           `json:"priority" binding:"required"`
	Milestones   []MilestoneInput `json:"milestones"`
}

type StatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type MilestoneInput struct {
	Title       string  `json:"title"`
	TargetValue float64 `json:"targetValue"`
	DueDate     string  `json:"dueDate"`
}

type MilestoneRequest struct {
	Completed bool `json:"completed"`
}
